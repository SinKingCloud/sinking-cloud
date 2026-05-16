import React, {useRef, useState} from 'react';
import {Body, ProModal, ProModalRef, ProTable, ProTableRef, Title} from "sinking-antd";
import {App, Button, Dropdown, Form, Input, Result} from "antd";
import {getData} from "@/utils/page";
import useControl from "@/utils/control";
import {
    createPermission,
    deletePermission,
    getPermissionList,
    updatePermission
} from "@/service/admin/permission";

export default (): React.ReactNode => {
    const {message, modal} = App.useApp();
    const control = useControl();
    const hasListControl = control.hasControl("admin.system.permission:list");
    const hasCreateControl = control.hasControl("admin.system.permission:create");
    const hasUpdateControl = control.hasControl("admin.system.permission:update");
    const hasDeleteControl = control.hasControl("admin.system.permission:delete");

    /**
     * 新增/编辑表单
     */
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const tableRef = useRef<ProTableRef>({} as ProTableRef);
    const [form] = Form.useForm();
    const [formSubmitting, setFormSubmitting] = useState(false);
    const [formScene, setFormScene] = useState<any>({
        action: "create",
        ids: [],
        title: "新增权限",
    });
    const isCreateMode = formScene?.action === "create";

    const openCreate = () => {
        if (!hasCreateControl) {
            message?.warning("暂无创建权限");
            return;
        }
        setFormScene({
            action: "create",
            ids: [],
            title: "新增权限",
        });
        form.resetFields();
        modalRef.current?.show();
    };

    const openEdit = (record: any) => {
        if (!hasUpdateControl) {
            message?.warning("暂无修改权限");
            return;
        }
        if (!record?.id) {
            return;
        }
        setFormScene({
            action: "edit",
            ids: [record.id],
            title: "编辑权限",
        });
        form.resetFields();
        form.setFieldsValue({
            code: record?.code,
            name: record?.name,
            path: record?.path,
        });
        modalRef.current?.show();
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
        setFormSubmitting(true);
        const request = isCreateMode ? createPermission : updatePermission;
        await request({
            body: isCreateMode ? {
                code: values?.code,
                name: values?.name,
                path: values?.path,
            } : {
                ids: formScene?.ids || [],
                name: values?.name,
                path: values?.path,
            },
            onSuccess: (res: any) => {
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
     * 删除权限
     */
    const onDelete = (ids: any[] = []) => {
        if (!hasDeleteControl) {
            message?.warning("暂无删除权限");
            return;
        }
        const nextIds = (ids || []).filter(Boolean);
        if (nextIds.length <= 0) {
            message?.warning("请选择要删除的数据");
            return;
        }
        modal.confirm({
            title: "删除权限",
            content: `确定删除选中的 ${nextIds.length} 条权限吗？`,
            okText: "确 定",
            okType: "danger",
            cancelText: "取 消",
            maskClosable: true,
            onOk: async () => {
                await deletePermission({
                    body: {
                        ids: nextIds,
                    },
                    onSuccess: (res: any) => {
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
            tip: "权限ID",
            width: 80,
            sorter: true,
            hideInSearch: true,
        },
        {
            title: "权限标识",
            dataIndex: "code",
            tip: "权限唯一标识",
            valueType: "text",
            copyable: true,
        },
        {
            title: "权限名称",
            dataIndex: "name",
            tip: "权限显示名称",
            valueType: "text",
            copyable: true,
        },
        {
            title: "请求路径",
            dataIndex: "path",
            tip: "权限绑定的接口路径",
            valueType: "text",
            copyable: true,
        },
        {
            title: "创建时间",
            valueType: "dateRange",
            dataIndex: "create_time",
            tip: "创建时间",
            sorter: true,
            transform: (value: any) => ({
                create_time_start: value?.[0]?.format ? value[0].format("YYYY-MM-DD HH:mm:ss") : value?.[0],
                create_time_end: value?.[1]?.format ? value[1].format("YYYY-MM-DD HH:mm:ss") : value?.[1],
            }),
        },
        {
            title: "更新时间",
            valueType: "dateRange",
            dataIndex: "update_time",
            tip: "更新时间",
            sorter: true,
            transform: (value: any) => ({
                update_time_start: value?.[0]?.format ? value[0].format("YYYY-MM-DD HH:mm:ss") : value?.[0],
                update_time_end: value?.[1]?.format ? value[1].format("YYYY-MM-DD HH:mm:ss") : value?.[1],
            }),
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
        <Body loading={control?.loading}>
            {!hasListControl && <Result
                status="403"
                title="403"
                subTitle="您无权限访问此页面,请联系管理员"
            />}
            {hasListControl && <>
                <ProTable
                    ref={tableRef}
                    extraRefreshBtn={true}
                    title={<Title>权限列表</Title>}
                    rowKey="id"
                    columns={columns}
                    defaultPage={1}
                    defaultPageSize={10}
                    extra={hasCreateControl ? (
                        <Button type="primary" onClick={openCreate}>新增</Button>
                    ) : undefined}
                    rowSelection={hasDeleteControl ? {
                        rightExtra: (
                            <Button key="delete" type="primary" danger ghost
                                    onClick={() => onDelete(tableRef.current?.getSelectedRowKeys?.() as any)}>
                                批量删除
                            </Button>
                        )
                    } : false}
                    request={(params, sort) => {
                        if (!hasListControl) {
                            return Promise.resolve({
                                data: [],
                                success: true,
                                total: 0,
                            });
                        }
                        return getData(params, sort, getPermissionList);
                    }}
                    paginationAffix={true}
                    selectionAffix={true}
                />

                <ProModal
                    ref={modalRef}
                    title={<Title>{formScene?.title}</Title>}
                    onOk={form.submit}
                    width={430}
                    modalProps={{
                        confirmLoading: formSubmitting,
                        forceRender: true,
                        okText: "提交",
                        cancelText: "取消",
                        afterClose: () => {
                            form.resetFields();
                            setFormSubmitting(false);
                        },
                    } as any}>
                    <Form form={form} layout="vertical" onFinish={handleSubmit}>
                        {isCreateMode ? (
                            <Form.Item name="code" label="权限标识"
                                       rules={[{required: true, message: "请输入权限标识"}]}>
                                <Input placeholder="请输入权限标识" maxLength={50} showCount/>
                            </Form.Item>
                        ) : (
                            <Form.Item name="code" label="权限标识" hidden>
                                <Input/>
                            </Form.Item>
                        )}
                        <Form.Item name="name" label="权限名称" rules={[{required: true, message: "请输入权限名称"}]}>
                            <Input placeholder="请输入权限名称" maxLength={isCreateMode ? 20 : 100} showCount/>
                        </Form.Item>
                        <Form.Item name="path" label="请求路径" rules={[{required: true, message: "请输入请求路径"}]}>
                            <Input placeholder="请输入请求路径" maxLength={isCreateMode ? 100 : 200} showCount/>
                        </Form.Item>
                    </Form>
                </ProModal>
            </>}
        </Body>
    );
};
