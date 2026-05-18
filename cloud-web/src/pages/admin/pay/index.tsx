import React, {useRef, useState} from 'react';
import {Body, ProModal, ProModalRef, ProTable, ProTableRef, Title} from "sinking-antd";
import {App, Button, Col, Dropdown, Form, Input, Result, Row, Select, Spin} from "antd";
import {getData} from "@/utils/page";
import useControl from "@/utils/control";
import useEnum from "@/utils/enum";
import {
    createPayChannel,
    deletePayChannel,
    getPayChannelInfo,
    getPayChannelList,
    updatePayChannel
} from "@/service/admin/pay";

export default (): React.ReactNode => {
    const {message, modal} = App.useApp();
    const [enumData, enumLoading] = useEnum(["pay.channel"]) as any;
    const control = useControl();
    const hasListControl = control.hasControl("admin.system.pay:list");
    const hasCreateControl = control.hasControl("admin.system.pay:create");
    const hasUpdateControl = control.hasControl("admin.system.pay:update");
    const hasDeleteControl = control.hasControl("admin.system.pay:delete");
    const channelEnums = enumData?.["pay.channel"] || {};
    const channelTypes = channelEnums?.type || {};
    const channelPlugins = channelEnums?.plugin || {};
    const channelStatus = channelEnums?.status || {};

    const buildOptions = (data: any, numberValue = false) => Object.entries(data || {}).map(([key, value]) => ({
        label: value,
        value: numberValue ? Number(key) : key,
    }));
    const buildValueEnum = (data: any, colorMap?: any) => Object.fromEntries(Object.entries(data || {}).map(([key, value]) => {
        const item: any = {text: value};
        if (colorMap?.[key]) {
            item.color = colorMap[key];
        }
        return [key, item];
    }));

    const typeOptions = buildOptions(channelTypes);
    const pluginOptions = buildOptions(channelPlugins, true);
    const statusOptions = buildOptions(channelStatus, true);
    const typeValueEnum = buildValueEnum(channelTypes, {alipay: "blue", wxpay: "green"});
    const pluginValueEnum = buildValueEnum(channelPlugins);
    const statusValueEnum = buildValueEnum(channelStatus, {"0": "green", "1": "red"});

    /**
     * 通道表单
     */
    const tableRef = useRef<ProTableRef>({} as ProTableRef);
    const channelModalRef = useRef<ProModalRef>({} as ProModalRef);
    const [channelForm] = Form.useForm();
    const [channelScene, setChannelScene] = useState<any>({
        action: "create",
        id: undefined,
        title: "新增通道",
    });
    const [channelSubmitting, setChannelSubmitting] = useState(false);
    const [channelInfoLoading, setChannelInfoLoading] = useState(false);
    const [configFields, setConfigFields] = useState<string[]>([]);
    const isEditMode = channelScene?.action === "edit";

    const openCreate = () => {
        if (!hasCreateControl) {
            message?.warning("暂无创建权限");
            return;
        }
        const defaultPlugin = Number(pluginOptions?.[0]?.value || 0);
        setChannelScene({
            action: "create",
            id: undefined,
            title: "新增通道",
        });
        channelForm.resetFields();
        channelForm.setFieldsValue({
            plugin: defaultPlugin,
            status: 0,
        });
        setConfigFields([]);
        setChannelInfoLoading(false);
        channelModalRef.current?.show();
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
        setChannelScene({
            action: "edit",
            id,
            title: "编辑通道",
        });
        channelForm.resetFields();
        setConfigFields([]);
        setChannelInfoLoading(true);
        channelModalRef.current?.show();
        getPayChannelInfo({
            body: {id},
            onSuccess: (res: any) => {
                const info = res?.data || {};
                channelForm.setFieldsValue({
                    id: info?.id,
                    name: info?.name,
                    status: Number(info?.status),
                });
                const normalizeConfig = (config: any) => {
                    if (typeof config === "string") {
                        try {
                            return JSON.parse(config || "{}");
                        } catch (_) {
                            return {};
                        }
                    }
                    return config || {};
                };
                const nextConfig = normalizeConfig(info?.config);
                channelForm.resetFields(["config"] as any);
                setConfigFields(Object.keys(nextConfig));
                channelForm.setFieldsValue({
                    config: nextConfig,
                });
            },
            onFail: (res: any) => {
                message?.error(res?.message || "获取通道详情失败");
                channelModalRef.current?.hide();
            },
        }).finally(() => {
            setChannelInfoLoading(false)
        });
    };

    const handleChannelSubmit = (values: any) => {
        if (isEditMode && !hasUpdateControl) {
            message?.warning("暂无修改权限");
            return;
        }
        if (!isEditMode && !hasCreateControl) {
            message?.warning("暂无创建权限");
            return;
        }
        const body: any = {
            name: values?.name,
        };
        if (isEditMode) {
            body.id = channelScene?.id || values?.id;
            body.status = values?.status?.toString();
            if (configFields.length > 0) {
                body.config = JSON.stringify(values?.config || {});
            }
        } else {
            body.type = values?.type;
            body.plugin = Number(values?.plugin);
            body.status = Number(values?.status);
        }
        setChannelSubmitting(true);
        (isEditMode ? updatePayChannel : createPayChannel)({
            body,
            onSuccess: (res: any) => {
                channelModalRef.current?.hide();
                tableRef.current?.refreshTableData?.();
                tableRef.current?.clearSelectedRows?.();
                message?.success(res?.message || "操作成功");
            },
            onFail: (res: any) => {
                message?.error(res?.message || "操作失败");
            },
            onFinally: () => setChannelSubmitting(false),
        });
    };

    /**
     * 删除通道
     */
    const onDelete = (id: any) => {
        if (!hasDeleteControl) {
            message?.warning("暂无删除权限");
            return;
        }
        modal.confirm({
            title: "删除支付通道",
            content: `确定删除该支付通道吗？`,
            okText: "确 定",
            okType: "danger",
            cancelText: "取 消",
            maskClosable: true,
            onOk: async () => {
                await deletePayChannel({
                    body: {
                        id: id
                    },
                    onSuccess: (res) => {
                        tableRef.current?.refreshTableData?.();
                        message?.success(res?.message || "删除成功");
                    },
                    onFail: (res) => {
                        message?.error(res?.message || "删除失败");
                        return;
                    }
                });
            },
        } as any);
    };

    /**
     * 表格
     */
    const columns: any[] = [
        {
            title: "ID",
            dataIndex: "id",
            tip: "通道ID",
            width: 80,
            sorter: true,
            hideInSearch: true,
        },
        {
            title: "通道名称",
            dataIndex: "name",
            tip: "通道名称",
            valueType: "text",
            copyable: true,
        },
        {
            title: "通道类型",
            dataIndex: "type",
            tip: "通道类型",
            valueEnum: typeValueEnum,
            hideInSearch: true,
        },
        {
            title: "使用插件",
            dataIndex: "plugin",
            tip: "使用插件",
            valueEnum: pluginValueEnum,
        },
        {
            title: "通道状态",
            dataIndex: "status",
            tip: "通道状态",
            valueEnum: statusValueEnum,
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
                        onClick: () => onDelete(record?.id),
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
                    title={<Title>支付管理</Title>}
                    rowKey="id"
                    columns={columns}
                    defaultPage={1}
                    defaultPageSize={10}
                    extra={hasCreateControl ? (
                        <Button type="primary" onClick={openCreate}>新增</Button>
                    ) : undefined}
                    request={(params, sort) => {
                        if (!hasListControl) {
                            return Promise.resolve({
                                data: [],
                                success: true,
                                total: 0,
                            });
                        }
                        return getData(params, sort, getPayChannelList);
                    }}
                    paginationAffix={true}
                    selectionAffix={true}
                />

                <ProModal
                    ref={channelModalRef}
                    title={<Title>{channelScene?.title}</Title>}
                    onOk={channelForm.submit}
                    width={420}
                    modalProps={{
                        confirmLoading: channelSubmitting || channelInfoLoading,
                        forceRender: true,
                        okText: "提交",
                        cancelText: "取消",
                        footer: channelInfoLoading ? null : undefined,
                        afterClose: () => {
                            channelForm.resetFields();
                            setChannelSubmitting(false);
                            setChannelInfoLoading(false);
                            setConfigFields([]);
                        },
                    } as any}>
                    {channelInfoLoading ? (
                        <div style={{height: 240, display: "flex", alignItems: "center", justifyContent: "center"}}>
                            <Spin/>
                        </div>
                    ) : (
                        <Form form={channelForm} layout="vertical" onFinish={handleChannelSubmit}>
                            <Form.Item name="id" hidden>
                                <Input/>
                            </Form.Item>
                            <Row gutter={16}>
                                <Col xs={24} sm={12}>
                                    <Form.Item name="name" label="通道名称"
                                               rules={[{required: true, message: "请输入通道名称"}]}>
                                        <Input placeholder="请输入通道名称" maxLength={100} showCount/>
                                    </Form.Item>
                                </Col>
                                <Col xs={24} sm={12}>
                                    <Form.Item name="status" label="通道状态"
                                               rules={[{required: true, message: "请选择通道状态"}]}>
                                        <Select placeholder="请选择通道状态" options={statusOptions as any}/>
                                    </Form.Item>
                                </Col>
                            </Row>
                            {!isEditMode && (
                                <Row gutter={16}>
                                    <Col xs={24} sm={12}>
                                        <Form.Item name="type" label="通道类型"
                                                   rules={[{required: true, message: "请选择通道类型"}]}>
                                            <Select placeholder="请选择通道类型" options={typeOptions as any}/>
                                        </Form.Item>
                                    </Col>
                                    <Col xs={24} sm={12}>
                                        <Form.Item name="plugin" label="使用插件"
                                                   rules={[{required: true, message: "请选择使用插件"}]}>
                                            <Select placeholder="请选择使用插件" options={pluginOptions as any}/>
                                        </Form.Item>
                                    </Col>
                                </Row>
                            )}
                            {isEditMode && configFields.map((field) => (
                                <Form.Item key={field}
                                           name={["config", field] as any}
                                           label={field}
                                           preserve={false}
                                           rules={[{required: true, message: `请输入${field}`}]}>
                                    <Input placeholder={`请输入${field}`}/>
                                </Form.Item>
                            ))}
                        </Form>
                    )}
                </ProModal>
            </>}
        </Body>
    );
};
