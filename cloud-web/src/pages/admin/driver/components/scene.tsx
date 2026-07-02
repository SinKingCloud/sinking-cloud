import React, {useRef, useState} from "react";
import {App, Button, Dropdown, Form, Input, InputNumber, Select} from "antd";
import {ProModal, ProModalRef, ProTable, ProTableRef, Title} from "sinking-antd";
import {createStyles} from "antd-style";
import {getData} from "@/utils/page";
import useEnum, {clearEnumCache} from "@/utils/enum";
import {
    createScene,
    deleteScene,
    getSceneList,
    updateScene,
} from "@/service/admin/driver";
import {buildOptions, buildValueEnum, dateRangeTransform} from "@/utils/page";
import {formatSizeFromKb} from "@/utils/string";
import {AdminControlCode} from "@/utils/permission";

type SceneManageProps = {
    control: any;
};

const useStyles: any = createStyles(({css}: any) => ({
    fullInput: css`
        width: 100%;
    `,
}));

const SceneManage: React.FC<SceneManageProps> = ({control}) => {
    const {styles} = useStyles();
    const {message, modal} = App.useApp();
    const [enumData] = useEnum(["driver.channel", "driver.scene"]) as any;
    const tableRef = useRef<ProTableRef>({} as ProTableRef);
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const [form] = Form.useForm();
    const [submitLoading, setSubmitLoading] = useState(false);
    const [scene, setScene] = useState<any>({
        action: "create",
        id: undefined,
        title: "新增场景",
    });
    const isEditMode = scene?.action === "edit";
    const hasListControl = control.hasControl(AdminControlCode.SYSTEM_DRIVER_SCENE_LIST);
    const hasCreateControl = control.hasControl(AdminControlCode.SYSTEM_DRIVER_SCENE_CREATE);
    const hasUpdateControl = control.hasControl(AdminControlCode.SYSTEM_DRIVER_SCENE_UPDATE);
    const hasDeleteControl = control.hasControl(AdminControlCode.SYSTEM_DRIVER_SCENE_DELETE);
    const sceneEnums = enumData?.["driver.scene"] || {};
    const channelEnums = enumData?.["driver.channel"] || {};
    const statusOptions = buildOptions(sceneEnums?.status, true);
    const statusValueEnum = buildValueEnum(sceneEnums?.status, {"0": "green", "1": "red"});
    const loginOptions = buildOptions(sceneEnums?.login, true);
    const loginValueEnum = buildValueEnum(sceneEnums?.login, {"0": "blue", "1": "green"});
    const channelOptions = buildOptions(channelEnums?.name, true);
    const channelValueEnum = buildValueEnum(channelEnums?.name);

    const openCreate = () => {
        if (!hasCreateControl) {
            message?.warning("暂无创建权限");
            return;
        }
        setScene({
            action: "create",
            id: undefined,
            title: "新增场景",
        });
        form.resetFields();
        form.setFieldsValue({
            allow_ext: "*",
            max_size: 0,
            login: Number(loginOptions?.[0]?.value || 0),
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
            title: "编辑场景",
        });
        form.resetFields();
        form.setFieldsValue({
            id,
            name: record?.name,
            max_size: Number(record?.max_size || 0),
            allow_ext: record?.allow_ext,
            channel_id: Number(record?.channel_id || 0) || undefined,
            login: Number(record?.login),
            remark: record?.remark,
            status: Number(record?.status),
        });
        modalRef.current?.show();
    };

    const onDelete = (id: any) => {
        if (!hasDeleteControl) {
            message?.warning("暂无删除权限");
            return;
        }
        modal.confirm({
            title: "删除场景",
            content: "确定删除该场景吗？",
            okText: "确 定",
            okButtonProps: {danger: true, type: "default"},
            cancelText: "取 消",
            mask: {
                closable: true,
            },
            onOk: async () => {
                await deleteScene({
                    body: {id: Number(id)},
                    onSuccess: (res: any) => {
                        clearEnumCache("driver.scene");
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
            max_size: isEditMode ? values?.max_size?.toString() : Number(values?.max_size || 0),
            allow_ext: values?.allow_ext,
            login: isEditMode ? values?.login?.toString() : Number(values?.login),
            remark: values?.remark,
            status: isEditMode ? values?.status?.toString() : Number(values?.status),
        };
        if (isEditMode) {
            body.id = scene?.id || values?.id;
        } else {
            body.channel_id = Number(values?.channel_id);
        }
        setSubmitLoading(true);
        await (isEditMode ? updateScene : createScene)({
            body,
            onSuccess: (res: any) => {
                clearEnumCache("driver.scene");
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
            tip: "场景ID",
            width: 80,
            sorter: true,
            hideInSearch: true,
        },
        {
            title: "场景名称",
            dataIndex: "name",
            tip: "场景名称",
            formItemProps: {
                label: "场景名称",
            },
            valueType: "text",
            copyable: true,
        },
        {
            title: "通道",
            dataIndex: "channel_id",
            tip: "场景使用的存储通道",
            formItemProps: {
                label: "存储通道",
            },
            valueEnum: channelValueEnum,
        },
        {
            title: "最大限制(KB)",
            dataIndex: "max_size",
            tip: "场景允许上传的最大文件大小,单位KB",
            hideInSearch: true,
            render: (text: any) => {
                if (text > 0) {
                    return formatSizeFromKb(text);
                }
                return "无限制";
            },
        },
        {
            title: "登录限制",
            dataIndex: "login",
            tip: "场景上传是否需要登录",
            formItemProps: {
                label: "登录限制",
            },
            valueEnum: loginValueEnum,
        },
        {
            title: "允许类型",
            dataIndex: "allow_ext",
            tip: "场景允许上传的文件类型",
            valueType: "text",
            hideInSearch: true,
        },
        {
            title: "状态",
            dataIndex: "status",
            tip: "场景状态",
            formItemProps: {
                label: "场景状态",
            },
            valueEnum: statusValueEnum,
        },
        {
            title: "备注",
            dataIndex: "remark",
            tip: "场景备注",
            formItemProps: {
                label: "场景备注",
            },
            valueType: "text",
        },
        {
            title: "创建时间",
            valueType: "dateRange",
            dataIndex: "create_time",
            tip: "场景创建时间",
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
            tip: "场景更新时间",
            formItemProps: {
                label: "更新时间",
            },
            sorter: true,
            transform: dateRangeTransform("update_time"),
        },
        ...(hasUpdateControl || hasDeleteControl ? [{
            title: "操作",
            tip: "场景操作",
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
                title={<Title>场景管理</Title>}
                rowKey="id"
                columns={columns}
                defaultPage={1}
                defaultPageSize={10}
                extra={hasCreateControl ? (
                    <Button type="primary" onClick={openCreate}>新增</Button>
                ) : undefined}
                request={(params, sort) => getData(params, sort, getSceneList)}
                paginationAffix={true}
                selectionAffix={true}
            />

            <ProModal
                ref={modalRef}
                title={<Title>{scene?.title}</Title>}
                onOk={form.submit}
                width={420}
                modalProps={{
                    confirmLoading: submitLoading,
                    okText: "提交",
                    cancelText: "取消",
                    style: {
                        top: 100,
                        paddingBottom: 100,
                    },
                    afterClose: () => {
                        form.resetFields();
                        setSubmitLoading(false);
                    },
                } as any}
            >
                <Form form={form} layout="vertical" onFinish={onSubmit}>
                    <Form.Item name="id" hidden>
                        <Input/>
                    </Form.Item>
                    <Form.Item name="name" label="场景名称" rules={[{required: true, message: "请输入场景名称"}]}>
                        <Input placeholder="请输入场景名称" maxLength={100} showCount/>
                    </Form.Item>
                    {!isEditMode && (
                        <Form.Item name="channel_id" label="使用通道"
                                   rules={[{required: true, message: "请选择使用通道"}]}>
                            <Select placeholder="请选择使用通道" options={channelOptions as any}/>
                        </Form.Item>
                    )}
                    <Form.Item name="login" label="登录限制" rules={[{required: true, message: "请选择登录限制"}]}>
                        <Select placeholder="请选择登录限制" options={loginOptions as any}/>
                    </Form.Item>
                    <Form.Item name="status" label="场景状态" rules={[{required: true, message: "请选择场景状态"}]}>
                        <Select placeholder="请选择场景状态" options={statusOptions as any}/>
                    </Form.Item>
                    <Form.Item name="max_size" label="最大限制(KB)"
                               rules={[{required: true, message: "请输入最大限制,0为不限制"}]}>
                        <InputNumber className={styles.fullInput} min={0} precision={0}
                                     placeholder="请输入最大限制,0为不限制"/>
                    </Form.Item>
                    <Form.Item name="allow_ext" label="允许类型"
                               rules={[{required: true, message: "请输入限制类型,*为不限制"}]}>
                        <Input placeholder="请输入允许类型, 如 jpg,png 或 *"/>
                    </Form.Item>
                    <Form.Item name="remark" label="备注">
                        <Input.TextArea placeholder="请输入备注" rows={3}/>
                    </Form.Item>
                </Form>
            </ProModal>
        </>
    );
};

export default SceneManage;
