import React, {useRef, useState} from 'react';
import {Body, ProModal, ProModalRef, ProTable, ProTableRef, Title} from "sinking-antd";
import {App, Button, Dropdown, Empty, Form, Input, Result, Select, Space, Spin, Tag, Tree, Typography} from "antd";
import {buildOptions, buildValueEnum, dateRangeTransform, emptyTableData, getData, normalizeIds} from "@/utils/page";
import useControl from "@/utils/control";
import {AdminControlCode} from "@/utils/permission";
import useEnum, {clearEnumCache} from "@/utils/enum";
import {
    createRole,
    deleteRole,
    getRoleInfo,
    getRoleList,
    getRoleMenu,
    updateRole
} from "@/service/admin/role";

export default (): React.ReactNode => {
    const {message, modal} = App.useApp();
    const [enumData, enumLoading] = useEnum(["admin.role", "admin.menu"]) as any;
    const control = useControl();
    const hasListControl = control.hasControl(AdminControlCode.SYSTEM_ROLE_LIST);
    const hasCreateControl = control.hasControl(AdminControlCode.SYSTEM_ROLE_CREATE);
    const hasUpdateControl = control.hasControl(AdminControlCode.SYSTEM_ROLE_UPDATE);
    const hasDeleteControl = control.hasControl(AdminControlCode.SYSTEM_ROLE_DELETE);
    const roleStatus = enumData?.["admin.role"]?.status || {};
    const menuTypes = enumData?.["admin.menu"]?.types || {};
    const statusOptions = buildOptions(roleStatus, true);
    const statusValueEnum = buildValueEnum(roleStatus, {"0": "green", "1": "red"});

    /**
     * 新增/编辑表单
     */
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const tableRef = useRef<ProTableRef>({} as ProTableRef);
    const [form] = Form.useForm();
    const [formSubmitting, setFormSubmitting] = useState(false);
    const [infoLoading, setInfoLoading] = useState(false);
    const [infoLoaded, setInfoLoaded] = useState(false);
    const [formScene, setFormScene] = useState<any>({
        action: "create",
        id: undefined,
        title: "新增角色",
    });
    const isCreateMode = formScene?.action === "create";
    const canShowMenuTree = isCreateMode ? hasCreateControl : hasUpdateControl;

    /**
     * 菜单树
     */
    const [roleMenus, setRoleMenus] = useState<any[]>([]);
    const [menuLoaded, setMenuLoaded] = useState(false);
    const [menuLoading, setMenuLoading] = useState(false);
    const [checkedMenuIds, setCheckedMenuIds] = useState<any[]>([]);

    const buildTreeData = (list: any[] = []): any[] => {
        return (list || []).map((item) => ({
            key: Number(item?.id),
            title: (
                <Space size={6}>
                    <span>{item?.name}</span>
                    <Tag color={Number(item?.type) === 2 ? "blue" : "default"}>
                        {menuTypes?.[Number(item?.type)] || "菜单"}
                    </Tag>
                    <Typography.Text type="secondary">{item?.code}</Typography.Text>
                </Space>
            ),
            children: buildTreeData(item?.children || []),
        }));
    };

    const getDescendantIds = (item: any): number[] => {
        const ids: number[] = [];
        const loop = (list: any[] = []) => {
            (list || []).forEach((child) => {
                const id = Number(child?.id);
                if (id > 0) {
                    ids.push(id);
                }
                loop(child?.children || []);
            });
        };
        loop(item?.children || []);
        return ids;
    };

    const getTreeCheckedKeys = (ids: any[] = [], tree: any[] = roleMenus) => {
        const selected = new Set(normalizeIds(ids));
        const result = new Set(selected);
        const loop = (list: any[] = []) => {
            (list || []).forEach((item) => {
                const id = Number(item?.id);
                const children = item?.children || [];
                if (id > 0 && children.length > 0) {
                    const descendantIds = getDescendantIds(item);
                    const hasSelectedDescendant = descendantIds.some((itemId) => selected.has(itemId));
                    const allDescendantsSelected = descendantIds.length > 0 && descendantIds.every((itemId) => selected.has(itemId));
                    if (hasSelectedDescendant && !allDescendantsSelected) {
                        result.delete(id);
                    }
                    loop(children);
                }
            });
        };
        loop(tree);
        return Array.from(result);
    };

    const filterDetachedPageMenuIds = (ids: any[] = [], tree: any[] = roleMenus) => {
        const selected = new Set(normalizeIds(ids));
        const loop = (list: any[] = []) => {
            (list || []).forEach((item) => {
                const id = Number(item?.id);
                const type = Number(item?.type);
                const children = item?.children || [];
                if (id > 0 && type === 1 && children.length > 0 && selected.has(id)) {
                    const descendantIds = getDescendantIds(item);
                    const hasSelectedDescendant = descendantIds.some((itemId) => selected.has(itemId));
                    if (!hasSelectedDescendant) {
                        selected.delete(id);
                    }
                }
                loop(children);
            });
        };
        loop(tree);
        return Array.from(selected);
    };

    const collectMenuIdsWithParents = (ids: any[] = []) => {
        const selected = new Set(normalizeIds(ids));
        const parentMap = new Map<number, number>();
        const loop = (list: any[] = [], parentId = 0) => {
            (list || []).forEach((item) => {
                const id = Number(item?.id);
                if (id > 0 && parentId > 0) {
                    parentMap.set(id, parentId);
                }
                loop(item?.children || [], id);
            });
        };
        loop(roleMenus);
        Array.from(selected).forEach((id) => {
            let parentId = parentMap.get(id);
            while (parentId && parentId > 0) {
                selected.add(parentId);
                parentId = parentMap.get(parentId);
            }
        });
        return Array.from(selected);
    };

    const loadRoleMenus = (): Promise<any[]> => {
        if (!hasCreateControl && !hasUpdateControl) {
            return Promise.resolve(roleMenus);
        }
        if (menuLoaded) {
            return Promise.resolve(roleMenus);
        }
        let nextMenus: any[] = [];
        setMenuLoading(true);
        return getRoleMenu({
            onSuccess: (res: any) => {
                nextMenus = Array.isArray(res?.data) ? res.data : [];
                setRoleMenus(nextMenus);
                setMenuLoaded(true);
            },
            onFail: (res: any) => {
                setRoleMenus([]);
                message?.error(res?.message || "获取角色菜单失败");
            },
            onFinally: () => setMenuLoading(false),
        }).then(() => nextMenus);
    };

    const onMenuCheck = (keys: any) => {
        const nextKeys = normalizeIds(Array.isArray(keys) ? keys : keys?.checked || []);
        setCheckedMenuIds(nextKeys);
        form.setFieldValue("menu_ids" as any, collectMenuIdsWithParents(nextKeys));
    };

    const openCreate = async () => {
        if (!hasCreateControl) {
            message?.warning("暂无创建权限");
            return;
        }
        setFormScene({
            action: "create",
            id: undefined,
            title: "新增角色",
        });
        setInfoLoaded(false);
        setCheckedMenuIds([]);
        form.resetFields();
        form.setFieldsValue({
            status: 0,
            menu_ids: [],
        });
        modalRef.current?.show();
        await loadRoleMenus();
    };

    const loadEditInfo = (id: number) => {
        if (!hasUpdateControl) {
            return;
        }
        setInfoLoading(true);
        loadRoleMenus().then((menus) => {
            getRoleInfo({
                body: {id},
                onSuccess: (res: any) => {
                    const info = res?.data || {};
                    const menuIds = filterDetachedPageMenuIds(info?.menu_ids || [], menus);
                    form.setFieldsValue({
                        name: info?.name,
                        status: Number(info?.status),
                        menu_ids: menuIds,
                    });
                    setCheckedMenuIds(getTreeCheckedKeys(menuIds, menus));
                    setInfoLoaded(true);
                },
                onFail: (res: any) => {
                    message?.error(res?.message || "获取角色详情失败");
                },
            }).finally(() => {
                setInfoLoading(false)
            });
        });
    };

    const openEdit = (record: any) => {
        if (!hasUpdateControl) {
            message?.warning("暂无修改权限");
            return;
        }
        const id = Number(record?.id);
        if (!id) {
            return;
        }
        setFormScene({
            action: "edit",
            id,
            title: "编辑角色",
        });
        setInfoLoaded(false);
        setCheckedMenuIds([]);
        form.resetFields();
        form.setFieldsValue({
            name: record?.name,
            status: Number(record?.status),
        });
        modalRef.current?.show();
        loadEditInfo(id);
    };

    const handleSubmit = async (values: any) => {
        if (isCreateMode && !hasCreateControl) {
            message?.warning("暂无创建权限");
            return;
        }
        if (!isCreateMode && !hasUpdateControl) {
            message?.warning("暂无修改权限");
            return;
        }
        const menuIds = collectMenuIdsWithParents(values?.menu_ids || []);
        const body = isCreateMode ? {
            name: values?.name,
            status: Number(values?.status),
            menu_ids: menuIds,
        } : {
            id: formScene?.id,
            name: values?.name,
            status: values?.status?.toString(),
            ...(infoLoaded ? {menu_ids: menuIds} : {}),
        };
        setFormSubmitting(true);
        await (isCreateMode ? createRole : updateRole)({
            body,
            onSuccess: (res: any) => {
                clearEnumCache("admin.role");
                modalRef.current?.hide();
                tableRef.current?.refreshTableData?.();
                tableRef.current?.clearSelectedRows?.();
                message?.success(res?.message || "操作成功");
            },
            onFail: (res: any) => {
                message?.error(res?.message || "操作失败");
            },
            onFinally: () => setFormSubmitting(false),
        });
    };

    /**
     * 删除角色
     */
    const onDelete = (ids: any[] = []) => {
        if (!hasDeleteControl) {
            message?.warning("暂无删除权限");
            return;
        }
        const nextIds = normalizeIds(ids);
        if (nextIds.length <= 0) {
            message?.warning("请选择要删除的数据");
            return;
        }
        modal.confirm({
            title: "删除角色",
            content: `确定删除选中的 ${nextIds.length} 个角色吗？`,
            okText: "确 定",
            okType: "danger",
            cancelText: "取 消",
            onOk: async () => {
                await deleteRole({
                    body: {
                        ids: nextIds,
                    },
                    onSuccess: (res: any) => {
                        clearEnumCache("admin.role");
                        tableRef.current?.refreshTableData?.();
                        tableRef.current?.clearSelectedRows?.();
                        message?.success(res?.message || "删除成功");
                    },
                    onFail: (res: any) => {
                        message?.error(res?.message || "删除失败");
                    }
                });
            }
        } as any);
    };

    /**
     * 表格
     */
    const columns: any[] = [
        {
            title: "ID",
            dataIndex: "id",
            tip: "角色ID",
            width: 80,
            sorter: true,
            hideInSearch: true,
        },
        {
            title: "角色名称",
            dataIndex: "name",
            tip: "角色名称",
            valueType: "text",
            copyable: true,
        },
        {
            title: "角色状态",
            dataIndex: "status",
            tip: "角色状态",
            valueEnum: statusValueEnum,
        },
        {
            title: "创建时间",
            valueType: "dateRange",
            dataIndex: "create_time",
            tip: "创建时间",
            sorter: true,
            transform: dateRangeTransform("create_time"),
        },
        {
            title: "更新时间",
            valueType: "dateRange",
            dataIndex: "update_time",
            tip: "更新时间",
            sorter: true,
            transform: dateRangeTransform("update_time"),
        },
        ...(hasUpdateControl || hasDeleteControl ? [{
            title: "操作",
            valueType: "option",
            key: "option",
            width: 100,
            hideInSearch: true,
            render: (_: any, record: any) => {
                const items = [
                    hasUpdateControl ? {
                        key: "edit",
                        label: "编辑",
                        onClick: () => openEdit(record),
                    } : null,
                    hasDeleteControl ? {
                        key: "delete",
                        label: "删除",
                        danger: true,
                        onClick: () => onDelete([record?.id]),
                    } : null,
                ].filter(Boolean) as any;
                if (items.length <= 0) {
                    return "-";
                }
                return (
                    <Dropdown menu={{items}} trigger={["click"]} placement="bottom" arrow>
                        <Button size="small">操作</Button>
                    </Dropdown>
                );
            },
        }] : []),
    ] as any;

    return (
        <Body loading={control?.loading || enumLoading}>
            {!hasListControl && <Result
                status="403"
                title="403"
                subTitle="您无权限访问此页面,请联系管理员"
            />}
            {hasListControl && <>
                <ProTable
                    ref={tableRef}
                    extraRefreshBtn={true}
                    title={<Title>角色管理</Title>}
                    rowKey="id"
                    columns={columns}
                    defaultPage={1}
                    defaultPageSize={10}
                    extra={hasCreateControl ? (
                        <Button type="primary" onClick={openCreate}>新增</Button>
                    ) : undefined}
                    rowSelection={hasDeleteControl ? {
                        fixed: true,
                        rightExtra: (
                            <Button key="delete" type="primary" danger ghost
                                    onClick={() => onDelete(tableRef.current?.getSelectedRowKeys?.() as any)}>
                                批量删除
                            </Button>
                        )
                    } : false}
                    request={(params, sort) => {
                        if (!hasListControl) {
                            return Promise.resolve(emptyTableData);
                        }
                        return getData(params, sort, getRoleList);
                    }}
                    paginationAffix={true}
                    selectionAffix={true}
                />

                <ProModal
                    ref={modalRef}
                    title={<Title>{formScene?.title}</Title>}
                    onOk={form.submit}
                    width={560}
                    modalProps={{
                        confirmLoading: formSubmitting || menuLoading || infoLoading,
                        forceRender: true,
                        okText: "提交",
                        cancelText: "取消",
                        afterClose: () => {
                            form.resetFields();
                            setFormSubmitting(false);
                            setInfoLoading(false);
                            setInfoLoaded(false);
                            setCheckedMenuIds([]);
                        },
                    } as any}>
                    <Form form={form} layout="vertical" onFinish={handleSubmit}>
                        <Form.Item name="name" label="角色名称" rules={[{required: true, message: "请输入角色名称"}]}>
                            <Input placeholder="请输入角色名称" maxLength={30} showCount/>
                        </Form.Item>
                        <Form.Item name="status" label="角色状态" rules={[{required: true, message: "请选择角色状态"}]}>
                            <Select placeholder="请选择角色状态" options={statusOptions as any}/>
                        </Form.Item>
                        {canShowMenuTree && (
                            <Form.Item name="menu_ids" label="菜单项">
                                <Spin spinning={menuLoading || infoLoading}>
                                    {roleMenus.length > 0 ? (
                                        <Tree
                                            checkable
                                            defaultExpandAll
                                            checkedKeys={checkedMenuIds}
                                            treeData={buildTreeData(roleMenus)}
                                            onCheck={onMenuCheck as any}
                                        />
                                    ) : (
                                        <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="暂无菜单项"/>
                                    )}
                                </Spin>
                            </Form.Item>
                        )}
                    </Form>
                </ProModal>
            </>}
        </Body>
    );
};
