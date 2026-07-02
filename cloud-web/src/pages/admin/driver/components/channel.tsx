import React, {useRef, useState} from "react";
import {App, Button, Dropdown, Form, Input, Select, Spin} from "antd";
import {ProModal, ProModalRef, ProTable, ProTableRef, Title} from "sinking-antd";
import {createStyles} from "antd-style";
import {getData} from "@/utils/page";
import useEnum, {clearEnumCache} from "@/utils/enum";
import {
    createChannel,
    deleteChannel,
    getChannelInfo,
    getChannelList,
    updateChannel,
} from "@/service/admin/driver";
import {buildOptions, buildValueEnum, dateRangeTransform} from "@/utils/page";
import {normalizeConfig} from "@/utils/string";
import {AdminControlCode} from "@/utils/permission";

type ChannelManageProps = {
    control: any;
};

const useStyles: any = createStyles(({css}: any) => ({
    infoLoading: css`
        height: 240px;
        display: flex;
        align-items: center;
        justify-content: center;
    `,
}));

const ChannelManage: React.FC<ChannelManageProps> = ({control}) => {
    const {styles} = useStyles();
    const {message, modal} = App.useApp();
    const [enumData] = useEnum(["driver.channel"]) as any;
    const tableRef = useRef<ProTableRef>({} as ProTableRef);
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const [form] = Form.useForm();
    const [scene, setScene] = useState<any>({
        action: "create",
        id: undefined,
        title: "新增通道",
    });
    const [submitLoading, setSubmitLoading] = useState(false);
    const [infoLoading, setInfoLoading] = useState(false);
    const [configFields, setConfigFields] = useState<string[]>([]);
    const isEditMode = scene?.action === "edit";
    const hasListControl = control.hasControl(AdminControlCode.SYSTEM_DRIVER_CHANNEL_LIST);
    const hasCreateControl = control.hasControl(AdminControlCode.SYSTEM_DRIVER_CHANNEL_CREATE);
    const hasUpdateControl = control.hasControl(AdminControlCode.SYSTEM_DRIVER_CHANNEL_UPDATE);
    const hasDeleteControl = control.hasControl(AdminControlCode.SYSTEM_DRIVER_CHANNEL_DELETE);
    const channelEnums = enumData?.["driver.channel"] || {};
    const typeOptions = buildOptions(channelEnums?.type, true);
    const statusOptions = buildOptions(channelEnums?.status, true);
    const typeValueEnum = buildValueEnum(channelEnums?.type);
    const statusValueEnum = buildValueEnum(channelEnums?.status, {"0": "green", "1": "red"});

    const openCreate = () => {
        if (!hasCreateControl) {
            message?.warning("暂无创建权限");
            return;
        }
        setScene({
            action: "create",
            id: undefined,
            title: "新增通道",
        });
        form.resetFields();
        setConfigFields([]);
        setInfoLoading(false);
        form.setFieldsValue({
            type: Number(typeOptions?.[0]?.value || 0),
            status: Number(statusOptions?.[0]?.value || 0),
        });
        modalRef.current?.show();
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
        setScene({
            action: "edit",
            id,
            title: "编辑通道",
        });
        form.resetFields();
        setConfigFields([]);
        setInfoLoading(true);
        modalRef.current?.show();
        getChannelInfo({
            body: {id},
            onSuccess: (res: any) => {
                const info = res?.data || {};
                const config = normalizeConfig(info?.config);
                form.setFieldsValue({
                    id: info?.id,
                    name: info?.name,
                    status: Number(info?.status),
                    config,
                });
                setConfigFields(Object.keys(config));
            },
            onFail: (res: any) => {
                modalRef.current?.hide();
                message?.error(res?.message || "获取通道详情失败");
            },
        }).finally(() => {
            setInfoLoading(false);
        });
    };

    const onDelete = (id: any) => {
        if (!hasDeleteControl) {
            message?.warning("暂无删除权限");
            return;
        }
        modal.confirm({
            title: "删除通道",
            content: "确定删除该存储通道吗？",
            okText: "确 定",
            okButtonProps: {danger: true, type: "default"},
            cancelText: "取 消",
            mask: {
                closable: true,
            },
            onOk: async () => {
                await deleteChannel({
                    body: {id: Number(id)},
                    onSuccess: (res: any) => {
                        clearEnumCache("driver.channel");
                        tableRef.current?.refreshTableData?.();
                        message?.success(res?.message || "删除成功");
                    },
                    onFail: (res: any) => {
                        message?.error(res?.message || "删除失败");
                    },
                });
            },
        } as any);
    };

    const onSubmit = async (values: any) => {
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
            body.id = scene?.id || values?.id;
            body.status = values?.status?.toString();
            if (configFields.length > 0) {
                body.config = JSON.stringify(values?.config || {});
            }
        } else {
            body.type = Number(values?.type);
            body.status = Number(values?.status);
        }
        setSubmitLoading(true);
        await (isEditMode ? updateChannel : createChannel)({
            body,
            onSuccess: (res: any) => {
                clearEnumCache("driver.channel");
                modalRef.current?.hide();
                tableRef.current?.refreshTableData?.();
                tableRef.current?.clearSelectedRows?.();
                message?.success(res?.message || "操作成功");
            },
            onFail: (res: any) => {
                message?.error(res?.message || "操作失败");
            },
            onFinally: () => setSubmitLoading(false),
        });
    };

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
            tip: "存储通道名称",
            formItemProps: {
                label: "通道名称",
            },
            valueType: "text",
            copyable: true,
        },
        {
            title: "通道类型",
            dataIndex: "type",
            tip: "存储通道类型",
            formItemProps: {
                label: "通道类型",
            },
            valueEnum: typeValueEnum,
        },
        {
            title: "通道状态",
            dataIndex: "status",
            tip: "存储通道状态",
            formItemProps: {
                label: "通道状态",
            },
            valueEnum: statusValueEnum,
        },
        {
            title: "创建时间",
            valueType: "dateRange",
            dataIndex: "create_time",
            tip: "通道创建时间",
            formItemProps: {
                label: "创建时间",
            },
            sorter: true,
            transform: dateRangeTransform("create_time"),
        },
        {
            title: "更新时间",
            valueType: "dateRange",
            dataIndex: "update_time",
            tip: "通道更新时间",
            formItemProps: {
                label: "更新时间",
            },
            sorter: true,
            transform: dateRangeTransform("update_time"),
        },
        ...(hasUpdateControl || hasDeleteControl ? [{
            title: "操作",
            tip: "通道操作",
            valueType: "option",
            key: "option",
            width: 100,
            fixed: "right",
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
    ];

    if (!hasListControl) {
        return null;
    }

    return (
        <>
            <ProTable
                ref={tableRef}
                extraRefreshBtn={true}
                title={<Title>通道管理</Title>}
                rowKey="id"
                columns={columns}
                defaultPage={1}
                defaultPageSize={10}
                extra={hasCreateControl ? (
                    <Button type="primary" onClick={openCreate}>新增</Button>
                ) : undefined}
                request={(params, sort) => getData(params, sort, getChannelList)}
                paginationAffix={true}
                selectionAffix={true}
            />

            <ProModal
                ref={modalRef}
                title={<Title>{scene?.title}</Title>}
                onOk={form.submit}
                width={420}
                modalProps={{
                    confirmLoading: submitLoading || infoLoading,
                    forceRender: true,
                    okText: "提交",
                    cancelText: "取消",
                    footer: infoLoading ? null : undefined,
                    style: {
                        top: 100,
                        paddingBottom: 100,
                    },
                    afterClose: () => {
                        form.resetFields();
                        setSubmitLoading(false);
                        setInfoLoading(false);
                        setConfigFields([]);
                    },
                } as any}
            >
                {infoLoading ? (
                    <div className={styles.infoLoading}>
                        <Spin/>
                    </div>
                ) : (
                    <Form form={form} layout="vertical" onFinish={onSubmit}>
                        <Form.Item name="id" hidden>
                            <Input/>
                        </Form.Item>
                        <Form.Item name="name" label="通道名称" rules={[{required: true, message: "请输入通道名称"}]}>
                            <Input placeholder="请输入通道名称" maxLength={100} showCount/>
                        </Form.Item>
                        {!isEditMode && (
                            <Form.Item name="type" label="通道类型"
                                       rules={[{required: true, message: "请选择通道类型"}]}>
                                <Select placeholder="请选择通道类型" options={typeOptions as any}/>
                            </Form.Item>
                        )}
                        <Form.Item name="status" label="通道状态" rules={[{required: true, message: "请选择通道状态"}]}>
                            <Select placeholder="请选择通道状态" options={statusOptions as any}/>
                        </Form.Item>
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
        </>
    );
};

export default ChannelManage;
