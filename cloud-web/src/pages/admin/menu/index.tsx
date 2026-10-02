import React, {useEffect, useMemo, useRef, useState} from 'react';
import type {MenuProps} from "antd";
import {
    App,
    AutoComplete,
    Button,
    Card,
    Col,
    Descriptions,
    Dropdown,
    Empty,
    Form,
    Input,
    InputNumber,
    Menu,
    Popover,
    Result,
    Row,
    Select,
    Space,
    Spin,
    Tag,
    Typography,
} from "antd";
import {Body, getAllIconNames, Icon, ProModal, ProModalRef, ProTable, ProTableRef, Title} from "sinking-antd";

import {createStyles} from "antd-style";
import {createMenu, deleteMenu, getMenuControls, getMenuInfo, getMenuList, updateMenu} from "@/service/admin/menu";
import useEnum from "@/utils/enum";
import {buildOptions, emptyTableData, getData} from "@/utils/page";
import useControl, {clearControlCache} from "@/utils/control";
import {AdminControlCode, AdminPageCode} from "@/utils/permission";

const useStyles: any = createStyles(({css, token}: any) => ({
    menuCard: css`
        .ant-card-body {
            padding: ${token.paddingSM}px;
        }
    `,
    infoCard: css`
        .ant-card-body {
            padding: 17px 20px;
        }
    `,
    infoExtra: css`
        display: flex;
        align-items: center;
        gap: ${token.marginXXS}px;
    `,
    infoActionBtn: css`
        height: ${token.controlHeightSM}px !important;
        padding-inline: ${token.paddingXS}px !important;
        border-radius: ${token.borderRadiusSM}px !important;
        font-size: ${token.fontSizeSM}px !important;
        line-height: ${token.lineHeightSM} !important;
        display: inline-flex !important;
        align-items: center !important;
    `,
    cardIconBtn: css`
        width: ${token.controlHeightSM}px;
        height: ${token.controlHeightSM}px;
        padding: 0;
        border-radius: ${token.borderRadiusSM}px;
        display: inline-flex;
        align-items: center;
        justify-content: center;

        &,
        &:hover,
        &:focus,
        &:active {
            background: transparent !important;
        }

        .ant-btn-icon {
            font-size: ${token.fontSizeSM}px;
        }
    `,
    menu: css`
        border-inline-end: 0 !important;
        background: transparent;

        .ant-menu,
        .ant-menu-sub {
            background: transparent !important;
        }

        .ant-menu-item,
        .ant-menu-submenu-title {
            padding-right: ${token.paddingLG}px !important;
        }
    `,
    addMenu: css`
        width: 100%;
        margin-top: ${token.marginSM}px;
        border-style: dashed;
        font-size: ${token.fontSizeSM}px;

        .ant-btn-icon {
            font-size: ${token.fontSizeSM}px;
        }
    `,
    fullWidth: css`
        max-width: 260px;
    `,
    row: css`
        display: flex;
        align-items: center;
        gap: 4px;
    `,
    smallText: css`
        font-size: 12px;
    `,
    permissionPopover: css`
        font-size: 12px;
        line-height: 1.6;

        .ant-tag {
            font-size: 12px;
            line-height: 18px;
        }
    `,
    hidden: css`
        display: none;
    `,
    formLoading: css`
        min-height: 220px;
        display: flex;
        align-items: center;
        justify-content: center;
    `,
    iconLoading: css`
        min-height: 48px;
        display: flex;
        align-items: center;
        justify-content: center;
    `,
}));

export default (): React.ReactNode => {
    const {message, modal} = App.useApp();
    const {styles} = useStyles();
    const [enumData, enumLoading] = useEnum(["admin.menu", "admin.permission"]) as any;
    const control = useControl();
    const hasUpdateControl = control.hasControl(AdminControlCode.SYSTEM_MENU_UPDATE);
    const hasCreateControl = control.hasControl(AdminControlCode.SYSTEM_MENU_CREATE);
    const hasDeleteControl = control.hasControl(AdminControlCode.SYSTEM_MENU_DELETE);
    const hasListControl = control.hasControl(AdminControlCode.SYSTEM_MENU_LIST);
    const hasControlsControl = control.hasControl(AdminControlCode.SYSTEM_MENU_CONTROLS);
    const menuTypes = enumData?.["admin.menu"]?.types || {};
    const menuStatus = enumData?.["admin.menu"]?.status || {};
    const permissionCodes = enumData?.["admin.permission"]?.codes || {};
    const menuTypeOptions = buildOptions(menuTypes, true).filter((item: any) => item.value !== 2);
    const statusOptions = buildOptions(menuStatus, true);
    const permissionOptions = buildOptions(permissionCodes, false, (value, key) => `${value}（${key}）`);
    const statusColor: any = {
        "0": {text: menuStatus?.["0"], color: "success"},
        "1": {text: menuStatus?.["1"], color: "error"},
    };
    /**
     * 加载图标
     */
    const iconLoadedRef = useRef(false);
    const iconMountedRef = useRef(true);
    const noIconOption = {
        value: "无图标",
        label: "无图标",
    };
    const [iconOptions, setIconOptions] = useState<any[]>([]);
    const [iconLoading, setIconLoading] = useState(false);

    const loadIconOptions = () => {
        if (iconLoadedRef.current || iconLoading) {
            return;
        }
        setIconLoading(true);
        getAllIconNames().then((icons) => {
            if (!iconMountedRef.current) {
                return;
            }
            setIconOptions([
                noIconOption,
                ...(icons || []).map((item) => ({
                    value: item,
                    label: (
                        <Space size={6}>
                            <Icon type={item}/>
                            <span>{item}</span>
                        </Space>
                    ),
                })),
            ]);
            iconLoadedRef.current = true;
        }).catch(() => {
            if (iconMountedRef.current) {
                setIconOptions([noIconOption]);
            }
        }).finally(() => {
            if (iconMountedRef.current) {
                setIconLoading(false);
            }
        });
    };

    /**
     * 左侧菜单数据
     */
    const [menus, setMenus] = useState<any>([]);
    const [selectedKeys, setSelectedKeys] = useState<string[]>([]);
    const [openKeys, setOpenKeys] = useState<string[]>([]);
    const [loading, setLoading] = useState(true);
    const selectedId = Number(selectedKeys[0] || 0);
    const menuInfoRequestRef = useRef(0);
    const [menuInfo, setMenuInfo] = useState<any>();
    const [infoLoading, setInfoLoading] = useState(true);
    const loadMenus = (resetSelected?: boolean) => {
        if (!hasListControl) {
            setMenus([]);
            setMenuInfo(undefined);
            setSelectedKeys([]);
            setOpenKeys([]);
            setLoading(false);
            setInfoLoading(false);
            return;
        }
        const shouldResetSelected = resetSelected === true;
        setLoading(true);
        getMenuList({
            onSuccess: (res) => {
                const list = res?.data || [];
                const getOpenKeys = (list: any[] = []): string[] => {
                    const keys: string[] = [];
                    const loop = (items: any[]) => {
                        items.forEach((item: any) => {
                            if (item.children && item.children.length > 0) {
                                keys.push(String(item?.id));
                                loop(item.children);
                            }
                        });
                    };
                    loop(list);
                    return keys;
                };
                const buildMenuItems = (list: any[] = []): MenuProps["items"] => {
                    return list.map((item) => ({
                        key: String(item?.id),
                        icon: item?.icon && item.icon !== "0" ? <Icon type={item.icon}/> : undefined,
                        label: item.name,
                        onTitleClick: item.children && item.children.length > 0 ? ({key}: any) => {
                            setSelectedKeys([String(key)]);
                        } : undefined,
                        children: item.children && item.children.length > 0 ? buildMenuItems(item.children) : undefined,
                    }));
                };
                setMenus(buildMenuItems(list));
                setOpenKeys(getOpenKeys(list));
                if (list.length > 0 && (selectedKeys.length <= 0 || shouldResetSelected)) {
                    setSelectedKeys([String(list[0]?.id)]);
                } else if (shouldResetSelected) {
                    setSelectedKeys([]);
                }
            },
            onFail: (res) => {
                setMenus([]);
                setMenuInfo(undefined);
                setSelectedKeys([]);
                setOpenKeys([]);
                setInfoLoading(false);
                message?.error(res?.message || "获取数据失败");
            }
        }).finally(() => {
            setLoading(false);
        })
    };


    /**
     * 右侧菜单信息
     */
    const loadMenuInfo = (id?: number) => {
        const requestId = ++menuInfoRequestRef.current;
        if (!hasListControl) {
            setMenuInfo(undefined);
            setInfoLoading(false);
            return;
        }
        if (!id) {
            setMenuInfo(undefined);
            setInfoLoading(false);
            return;
        }
        setInfoLoading(true);
        getMenuInfo({
            body: {
                id: id
            },
            onFail: (res) => {
                if (requestId === menuInfoRequestRef.current) {
                    message?.error(res?.message || "获取数据失败");
                }
            },
            onSuccess: (res) => {
                if (requestId === menuInfoRequestRef.current) {
                    setMenuInfo(res?.data || {});
                }
            }
        }).finally(() => {
            if (requestId === menuInfoRequestRef.current) {
                setInfoLoading(false);
            }
        })
    };

    /**
     * 表格操作
     */
    const actionRef = useRef<ProTableRef>({} as ProTableRef);
    const controlColumns: any[] = [
        {
            title: '控件信息',
            dataIndex: 'name',
            tip: '控件名称/控件标识',
            width: 260,
            render: (_: any, record: any) => (
                <Space orientation="vertical" size={2} className={styles.fullWidth}>
                    <div className={styles.row}>
                        <Typography.Text strong className={styles.smallText} ellipsis={{tooltip: record?.name}}>
                            {record?.name || "-"}
                        </Typography.Text>
                    </div>
                    <div className={styles.row}>
                        <Typography.Text type="secondary" className={styles.smallText}
                                         ellipsis={{tooltip: record?.code}}>
                            {record?.code || "-"}
                        </Typography.Text>
                    </div>
                </Space>
            ),
        },
        {
            title: '权限名称',
            dataIndex: 'permission_codes',
            tip: '绑定权限',
            width: 100,
            hideInSearch: true,
            render: (_: any, record: any) => {
                const codes = record?.permission_codes || [];
                if (codes.length <= 0) {
                    return "-";
                }
                return (
                    <Popover
                        content={(
                            <Space orientation="vertical" size={4} className={styles.permissionPopover}>
                                {codes.map((code: string) => (
                                    <div key={code}>
                                        {permissionCodes?.[code] || code}
                                        <Tag style={{marginLeft: 8}}>{code}</Tag>
                                    </div>
                                ))}
                            </Space>
                        )}>
                        <span>查看</span>
                    </Popover>
                );
            },
        },
        {
            title: '状态',
            dataIndex: 'status',
            tip: '控件状态',
            width: 100,
            sorter: true,
            valueEnum: statusColor,
        },
        {
            title: '相关时间',
            dataIndex: 'create_time',
            tip: '创建时间/更新时间',
            width: 220,
            hideInSearch: true,
            sorter: true,
            render: (_: any, record: any) => (
                <Space orientation="vertical" size={2}>
                    <div>
                        <Typography.Text type="secondary">创建于 </Typography.Text>
                        <Typography.Text>{record?.create_time || "-"}</Typography.Text>
                    </div>
                    <div>
                        <Typography.Text type="secondary">更新于 </Typography.Text>
                        <Typography.Text>{record?.update_time || "-"}</Typography.Text>
                    </div>
                </Space>
            ),
        },
        ...(hasUpdateControl || hasDeleteControl ? [{
            title: '操作',
            width: 100,
            fixed: 'right',
            hideInSearch: true,
            render: (_: any, record: any) => {
                const items = [
                    hasUpdateControl ? {
                        key: 'edit',
                        label: '编辑',
                        onClick: () => openEdit(record?.id, "control"),
                    } : null,
                    hasDeleteControl ? {
                        key: 'delete',
                        label: '删除',
                        danger: true,
                        onClick: () => deleteMenuOrControls([record?.id]),
                    } : null
                ].filter(Boolean) as any;
                if (items.length <= 0) {
                    return "-";
                }
                return (
                    <Dropdown
                        menu={{items}}
                        trigger={['hover']}
                        arrow
                        placement="bottom">
                        <Button size="small">操作</Button>
                    </Dropdown>
                );
            }
        }] : []),
    ] as any;

    /**
     * 创建/编辑菜单或控件
     */
    const formModalRef = useRef<ProModalRef>({} as ProModalRef);
    const [form] = Form.useForm();
    const [formLoading, setFormLoading] = useState(false);
    const [formSubmitting, setFormSubmitting] = useState(false);
    const [formScene, setFormScene] = useState<any>({
        action: "create",
        target: "menu",
        title: "添加菜单",
        parentId: 0,
        type: 0,
        typeEditable: true,
    });
    const isEditMode = formScene?.action === "edit";
    const isControlForm = formScene?.target === "control";
    const showTypeField = !isEditMode && formScene?.typeEditable;
    const showCodeField = !isEditMode;
    const codeOptions = useMemo(() => {
        const codes = Object.values(isControlForm ? AdminControlCode : AdminPageCode);
        return codes.map((code) => ({
            value: code,
            label: code,
        }));
    }, [isControlForm]);
    const openCreateRootMenu = () => {
        if (!hasCreateControl) {
            message?.warning("暂无创建权限");
            return;
        }
        loadIconOptions();
        setFormScene({
            action: "create",
            target: "menu",
            title: "添加菜单",
            id: undefined,
            parentId: 0,
            type: 0,
            typeEditable: true,
        });
        form.resetFields();
        form.setFieldsValue({
            parent_id: 0,
            type: 0,
            sort: 0,
            status: 0,
            permission_codes: [],
        });
        setFormLoading(false);
        formModalRef.current?.show();
    };
    const openCreateChildMenu = () => {
        if (!hasCreateControl) {
            message?.warning("暂无创建权限");
            return;
        }
        if (!menuInfo?.id || menuInfo?.type !== 0) {
            message?.warning("请选择目录后再新增菜单");
            return;
        }
        loadIconOptions();
        setFormScene({
            action: "create",
            target: "menu",
            title: "新增菜单",
            id: undefined,
            parentId: menuInfo.id,
            type: 1,
            typeEditable: true,
        });
        form.resetFields();
        form.setFieldsValue({
            parent_id: menuInfo.id,
            type: 1,
            sort: 0,
            status: 0,
            permission_codes: [],
        });
        setFormLoading(false);
        formModalRef.current?.show();
    };
    const openCreateControl = () => {
        if (!hasCreateControl) {
            message?.warning("暂无创建权限");
            return;
        }
        if (!menuInfo?.id || menuInfo?.type !== 1) {
            message?.warning("请选择页面后再添加控件");
            return;
        }
        setFormScene({
            action: "create",
            target: "control",
            title: "添加控件",
            id: undefined,
            parentId: menuInfo.id,
            type: 2,
            typeEditable: false,
        });
        form.resetFields();
        form.setFieldsValue({
            parent_id: menuInfo.id,
            type: 2,
            status: 0,
            permission_codes: [],
        });
        setFormLoading(false);
        formModalRef.current?.show();
    };
    const openEdit = async (id: number, target: "menu" | "control") => {
        if (!hasUpdateControl) {
            message?.warning("暂无修改权限");
            return;
        }
        if (!id) {
            return;
        }
        if (target !== "control") {
            loadIconOptions();
        }
        setFormScene({
            action: "edit",
            target,
            title: target === "control" ? "编辑控件" : "编辑菜单",
            id,
            type: target === "control" ? 2 : menuInfo?.type,
            typeEditable: false,
        });
        form.resetFields();
        setFormLoading(true);
        formModalRef.current?.show();
        await getMenuInfo({
            body: {id},
            onSuccess: (res: any) => {
                const data = res?.data || {};
                const nextTarget = data?.type === 2 ? "control" : "menu";
                setFormScene((prev: any) => ({
                    ...prev,
                    target: nextTarget,
                    title: nextTarget === "control" ? "编辑控件" : "编辑菜单",
                    id: data?.id || id,
                    type: data?.type,
                }));
                const formValues: any = {
                    ...data,
                    permission_codes: data?.permission_codes || [],
                };
                if (nextTarget === "control") {
                    delete formValues.icon;
                } else {
                    formValues.icon = data?.icon === "0" ? "无图标" : data?.icon;
                }
                form.setFieldsValue(formValues);
            },
            onFail: (res: any) => {
                formModalRef.current?.hide();
                message.error(res?.message || "获取数据失败");
            },
            onFinally: () => setFormLoading(false),
        });
    };
    const handleSubmit = async (values: any) => {
        if (isEditMode && !hasUpdateControl) {
            message?.warning("暂无修改权限");
            return;
        }
        if (!isEditMode && !hasCreateControl) {
            message?.warning("暂无创建权限");
            return;
        }
        const type = isEditMode ? formScene?.type : (formScene?.typeEditable ? values?.type : formScene?.type);
        const target = type === 2 ? "control" : "menu";
        const params: any = {
            ...values,
            type,
        };
        if (params.icon === "无图标") {
            params.icon = "0";
        }
        if (target === "control") {
            delete params.icon;
        }
        if (!isEditMode) {
            params.parent_id = formScene?.parentId || 0;
            params.sort = target === "control" ? 0 : (values?.sort || 0);
            params.permission_codes = target === "control" ? (values?.permission_codes || []) : [];
        }
        if (isEditMode) {
            params.id = formScene?.id || values?.id;
            params.status = params.status?.toString();
            params.sort = params.sort?.toString();
            delete params.code;
            delete params.type;
            delete params.parent_id;
            if (target !== "control") {
                delete params.permission_codes;
            }
            if (target === "control") {
                delete params.sort;
            }
        }
        setFormSubmitting(true);
        await (isEditMode ? updateMenu : createMenu)({
            body: params,
            onSuccess: (res: any) => {
                message.success(res?.message || "操作成功");
                formModalRef.current?.hide();
                if (target === "control") {
                    if (menuInfo?.code) {
                        clearControlCache(menuInfo.code);
                    }
                    actionRef.current?.clearSelectedRows?.();
                    if (hasControlsControl) {
                        actionRef.current?.refreshTableData?.();
                    }
                } else {
                    if (hasListControl) {
                        loadMenus();
                    }
                    if (isEditMode && selectedId) {
                        loadMenuInfo(selectedId);
                    }
                }
            },
            onFail: (res: any) => {
                message.error(res?.message || "操作失败");
            },
            onFinally: () => setFormSubmitting(false),
        });
    };
    const deleteMenuOrControls = (ids: number[], refreshMenu = false) => {
        if (!hasDeleteControl) {
            message?.warning("暂无删除权限");
            return;
        }
        if (ids.length <= 0) {
            message?.warning("请选择要删除的数据");
            return;
        }
        modal.confirm({
            title: "确认删除",
            icon: <Icon type="ExclamationCircleOutlined"/>,
            content: "删除后数据不可恢复，确定要删除吗？",
            mask: {closable: true},
            okButtonProps: {danger: true, type: "default"},
            onOk: async () => {
                await deleteMenu({
                    body: {
                        ids
                    },
                    onSuccess: (res: any) => {
                        message.success(res?.message || "删除成功");
                        if (!refreshMenu && menuInfo?.code) {
                            clearControlCache(menuInfo.code);
                        }
                        actionRef.current?.clearSelectedRows?.();
                        if (hasControlsControl) {
                            actionRef.current?.refreshTableData?.();
                        }
                        if (refreshMenu) {
                            loadMenus(true);
                        }
                    },
                    onFail: (res: any) => {
                        message.error(res?.message || "删除失败");
                    },
                });
            },
        } as any);
    };

    useEffect(() => {
        return () => {
            iconMountedRef.current = false;
        };
    }, []);

    useEffect(() => {
        if (control.loading) {
            return;
        }
        loadMenus();
    }, [control.loading, hasListControl]);

    useEffect(() => {
        if (control.loading) {
            return;
        }
        loadMenuInfo(selectedId);
    }, [control.loading, hasListControl, selectedId]);

    const showMenuInfoActions = Boolean(menuInfo?.id && (hasCreateControl || hasUpdateControl || hasDeleteControl));
    const showControlTable = Boolean(hasControlsControl && menuInfo?.id && menuInfo?.type === 1);

    return (
        <Body loading={control?.loading}>
            {!hasListControl && <Result
                status="403"
                title="403"
                subTitle="您无权限访问此页面,请联系管理员"
            />}
            {hasListControl && <>
                <Row gutter={[10, 10]}>
                    <Col xs={24} sm={24} md={10} lg={7} xl={5} xxl={4} xxxl={3}>
                        <Card className={styles.menuCard} title={<Title>系统菜单</Title>} variant={"borderless"}
                              extra={
                                  hasListControl ? (
                                      <Button className={styles.cardIconBtn}
                                              type="text"
                                              size="small"
                                              title="刷新"
                                              icon={<Icon type="ReloadOutlined"/>}
                                              onClick={() => loadMenus()}/>
                                  ) : null
                              }>
                            <Spin spinning={loading}>
                                {hasListControl && menus.length > 0 ? (
                                    <Menu className={styles.menu}
                                          mode="inline"
                                          items={menus}
                                          selectedKeys={selectedKeys}
                                          openKeys={openKeys}
                                          onClick={({key}) => setSelectedKeys([key])}/>
                                ) : (
                                    <Empty image={Empty.PRESENTED_IMAGE_SIMPLE}/>
                                )}
                                {hasCreateControl && (
                                    <Button className={styles.addMenu} icon={<Icon type="PlusOutlined"/>}
                                            onClick={openCreateRootMenu}>
                                        添加菜单
                                    </Button>
                                )}
                            </Spin>
                        </Card>
                    </Col>
                    <Col xs={24} sm={24} md={14} lg={17} xl={19} xxl={20} xxxl={21}>
                        <Row gutter={[10, 10]}>
                            <Col span={24}>
                                <Card className={styles.infoCard} title={<Title>菜单信息</Title>} variant={"borderless"}
                                      extra={(
                                          <Space size={4} className={styles.infoExtra}>
                                              {showMenuInfoActions && (
                                                  <Dropdown
                                                      menu={{
                                                          items: [
                                                              hasCreateControl && menuInfo?.type === 0 ? {
                                                                  key: "create",
                                                                  label: "新增",
                                                                  onClick: openCreateChildMenu,
                                                              } : null,
                                                              hasUpdateControl ? {
                                                                  key: "edit",
                                                                  label: "编辑",
                                                                  onClick: () => openEdit(menuInfo?.id, "menu"),
                                                              } : null,
                                                              hasDeleteControl ? {
                                                                  key: "delete",
                                                                  label: "删除",
                                                                  danger: true,
                                                                  onClick: () => deleteMenuOrControls([menuInfo?.id], true),
                                                              } : null,
                                                          ].filter(Boolean) as any,
                                                      }}
                                                      trigger={["hover"]}
                                                      arrow={{pointAtCenter: true}}
                                                      placement="bottom">
                                                      <Button className={styles.infoActionBtn}
                                                              size="small">操作</Button>
                                                  </Dropdown>
                                              )}
                                              {hasListControl && (
                                                  <Button className={styles.cardIconBtn}
                                                          type="text"
                                                          size="small"
                                                          title="刷新"
                                                          icon={<Icon type="ReloadOutlined"/>}
                                                          onClick={() => loadMenuInfo(selectedId)}/>
                                              )}
                                          </Space>
                                      )}
                                      loading={infoLoading || enumLoading}>
                                    {menuInfo?.id ? (
                                        <Descriptions
                                            column={{xs: 1, sm: 1, md: 1, lg: 2, xl: 2, xxl: 2}} size={"small"}>
                                            <Descriptions.Item label="菜单名称">{menuInfo?.name}</Descriptions.Item>
                                            <Descriptions.Item label="菜单标识">
                                                <Typography.Text copyable>{menuInfo?.code}</Typography.Text>
                                            </Descriptions.Item>
                                            <Descriptions.Item
                                                label="菜单类型">{menuTypes?.[menuInfo?.type] || "-"}</Descriptions.Item>
                                            <Descriptions.Item label="菜单状态">
                                                {statusColor?.[String(menuInfo?.status)]?.text ? (
                                                    <Tag color={statusColor?.[String(menuInfo?.status)]?.color}>
                                                        {statusColor?.[String(menuInfo?.status)]?.text}
                                                    </Tag>
                                                ) : "-"}
                                            </Descriptions.Item>
                                            <Descriptions.Item
                                                label="菜单图标">{menuInfo?.icon && menuInfo?.icon !== "0" ? menuInfo.icon : "-"}</Descriptions.Item>
                                            <Descriptions.Item label="菜单排序">{menuInfo?.sort}</Descriptions.Item>
                                            <Descriptions.Item
                                                label="创建时间">{menuInfo?.create_time || "-"}</Descriptions.Item>
                                            <Descriptions.Item
                                                label="更新时间">{menuInfo?.update_time || "-"}</Descriptions.Item>
                                        </Descriptions>
                                    ) : (
                                        <Empty image={Empty.PRESENTED_IMAGE_SIMPLE}/>
                                    )}
                                </Card>
                            </Col>
                            {showControlTable && (
                                <Col span={24} className={infoLoading || enumLoading ? styles?.hidden : ""}>
                                    <ProTable
                                        key={selectedId}
                                        ref={actionRef}
                                        search={false}
                                        extraRefreshBtn={true}
                                        title={<Title>控件管理</Title>}
                                        rowKey={'id'}
                                        columns={controlColumns}
                                        request={(params, sort) => {
                                            if (!hasControlsControl || !selectedId) {
                                                return Promise.resolve(emptyTableData);
                                            }
                                            return getData({...params, id: selectedId}, sort, getMenuControls)
                                        }}
                                        defaultPage={1}
                                        defaultPageSize={10}
                                        extra={hasCreateControl ? (
                                            <Button type="primary" className={styles.infoActionBtn}
                                                    onClick={openCreateControl}>新增</Button>
                                        ) : undefined}
                                        rowSelection={hasDeleteControl ? {
                                            fixed: true,
                                            rightExtra: (
                                                <Button key="delete" danger ghost
                                                        onClick={() => deleteMenuOrControls(actionRef.current?.getSelectedRows?.().map((row: any) => row.id) || [])}>
                                                    批量删除
                                                </Button>
                                            )
                                        } : undefined}
                                    />
                                </Col>
                            )}
                        </Row>
                    </Col>
                </Row>
                <ProModal
                    ref={formModalRef}
                    title={<Title>{formScene?.title}</Title>}
                    width={420}
                    onOk={form.submit}
                    modalProps={{
                        mask: {closable: true},
                        confirmLoading: formSubmitting,
                        forceRender: true,
                        okText: "提交",
                        cancelText: "取消",
                        okButtonProps: {
                            disabled: formLoading,
                        },
                        style: {
                            top: 100,
                            paddingBottom: 100,
                        },
                        afterClose: () => {
                            form.resetFields();
                            setFormLoading(false);
                            setFormSubmitting(false);
                        },
                    } as any}>
                    <div className={formLoading ? styles.formLoading : styles.hidden}>
                        <Spin/>
                    </div>
                    <Form form={form} onFinish={handleSubmit} layout={"vertical"}
                          className={formLoading ? styles.hidden : ""}>
                        <Form.Item name="id" hidden>
                            <Input/>
                        </Form.Item>
                        <Form.Item name="parent_id" hidden>
                            <Input/>
                        </Form.Item>
                        {isControlForm ? (
                            <>
                                <Row gutter={16}>
                                    <Col xs={24} sm={isEditMode ? 12 : 24}>
                                        <Form.Item name="name" label="名称"
                                                   rules={[{required: true, message: "请输入名称"}]}>
                                            <Input placeholder="请输入名称" maxLength={100} showCount/>
                                        </Form.Item>
                                    </Col>
                                    {isEditMode && (
                                        <Col xs={24} sm={12}>
                                            <Form.Item name="status" label="状态"
                                                       rules={[{required: true, message: "请选择状态"}]}>
                                                <Select placeholder="请选择状态" options={statusOptions}/>
                                            </Form.Item>
                                        </Col>
                                    )}
                                </Row>
                                {!isEditMode && (
                                    <Row gutter={16}>
                                        <Col xs={24} sm={12}>
                                            <Form.Item name="code" label="标识"
                                                       rules={[{required: true, message: "请输入标识"}]}>
                                                <AutoComplete placeholder="请选择或输入标识"
                                                              allowClear
                                                              maxLength={100}
                                                              options={codeOptions}
                                                              showSearch={{
                                                                  filterOption: (inputValue: string, option: any) =>
                                                                      String(option?.value || "").toLowerCase().includes(inputValue.toLowerCase())
                                                              }}/>
                                            </Form.Item>
                                        </Col>
                                        <Col xs={24} sm={12}>
                                            <Form.Item name="status" label="状态"
                                                       rules={[{required: true, message: "请选择状态"}]}>
                                                <Select placeholder="请选择状态" options={statusOptions}/>
                                            </Form.Item>
                                        </Col>
                                    </Row>
                                )}
                                <Form.Item name="permission_codes" label="绑定权限"
                                           rules={[{required: true, message: "请选择绑定权限"}]}>
                                    <Select placeholder="请选择绑定权限"
                                            mode="multiple"
                                            allowClear
                                            showSearch
                                            options={permissionOptions}/>
                                </Form.Item>
                            </>
                        ) : (
                            <>
                                <Row gutter={16}>
                                    <Col xs={24} sm={showCodeField ? 12 : 24}>
                                        <Form.Item name="name" label="名称"
                                                   rules={[{required: true, message: "请输入名称"}]}>
                                            <Input placeholder="请输入名称" maxLength={100} showCount/>
                                        </Form.Item>
                                    </Col>
                                    {showCodeField && (
                                        <Col xs={24} sm={12}>
                                            <Form.Item name="code" label="标识"
                                                       rules={[{required: true, message: "请输入标识"}]}>
                                                <AutoComplete placeholder="请选择或输入标识"
                                                              allowClear
                                                              maxLength={100}
                                                              options={codeOptions}
                                                              showSearch={{
                                                                  filterOption: (inputValue: string, option: any) =>
                                                                      String(option?.value || "").toLowerCase().includes(inputValue.toLowerCase())
                                                              }}/>
                                            </Form.Item>
                                        </Col>
                                    )}
                                </Row>
                                <Row gutter={16}>
                                    <Col xs={24} sm={12}>
                                        <Form.Item name="icon" label="图标">
                                            <AutoComplete placeholder="请选择或输入图标"
                                                          allowClear
                                                          maxLength={30}
                                                          options={iconOptions}
                                                          notFoundContent={iconLoading ? (
                                                              <div className={styles.iconLoading}>
                                                                  <Spin size="small"/>
                                                              </div>
                                                          ) : null}
                                                          showSearch={{
                                                              filterOption: (inputValue: string, option: any) =>
                                                                  String(option?.value || "").toLowerCase().includes(inputValue.toLowerCase())
                                                          }
                                                          }/>
                                        </Form.Item>
                                    </Col>
                                    {showTypeField && (
                                        <Col xs={24} sm={12}>
                                            <Form.Item name="type" label="类型"
                                                       rules={[{required: true, message: "请选择类型"}]}>
                                                <Select placeholder="请选择类型" options={menuTypeOptions}/>
                                            </Form.Item>
                                        </Col>
                                    )}
                                    {!showTypeField && (
                                        <Col xs={24} sm={12}>
                                            <Form.Item name="status" label="状态"
                                                       rules={[{required: true, message: "请选择状态"}]}>
                                                <Select placeholder="请选择状态" options={statusOptions}/>
                                            </Form.Item>
                                        </Col>
                                    )}
                                </Row>
                                {showTypeField ? (
                                    <Row gutter={16}>
                                        <Col xs={24} sm={12}>
                                            <Form.Item name="status" label="状态"
                                                       rules={[{required: true, message: "请选择状态"}]}>
                                                <Select placeholder="请选择状态" options={statusOptions}/>
                                            </Form.Item>
                                        </Col>
                                        <Col xs={24} sm={12}>
                                            <Form.Item name="sort" label="排序"
                                                       rules={[{required: true, message: "请输入排序"}]}>
                                                <InputNumber min={0} precision={0} style={{width: "100%"}}/>
                                            </Form.Item>
                                        </Col>
                                    </Row>
                                ) : (
                                    <Form.Item name="sort" label="排序"
                                               rules={[{required: true, message: "请输入排序"}]}>
                                        <InputNumber min={0} precision={0} style={{width: "100%"}}/>
                                    </Form.Item>
                                )}
                            </>
                        )}
                    </Form>
                </ProModal>
            </>}
        </Body>
    );
};
