import React, {useRef, useState} from "react";
import {App, Button, Checkbox, Dropdown, Form, Select, Space} from "antd";
import {UploadOutlined} from "@ant-design/icons";
import {ProModal, ProModalRef, ProTable, ProTableRef, Title} from "sinking-antd";
import {getData} from "@/utils/page";
import useEnum, {clearEnumCache} from "@/utils/enum";
import {
    deleteFile,
    getFileList,
    updateFile,
} from "@/service/admin/driver";
import {buildOptions, buildValueEnum, dateRangeTransform, formatSize} from "./common";
import FileUploadModal, {type FileUploadModalRef} from "./upload";

type FileManageProps = {
    control: any;
};

const FileManage: React.FC<FileManageProps> = ({control}) => {
    const {message, modal} = App.useApp();
    const [enumData] = useEnum(["driver.channel", "driver.scene", "driver.file"]) as any;
    const tableRef = useRef<ProTableRef>({} as ProTableRef);
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const fileUploadRef = useRef<FileUploadModalRef>({} as FileUploadModalRef);
    const [form] = Form.useForm();
    const [submitLoading, setSubmitLoading] = useState(false);
    const [editIds, setEditIds] = useState<any[]>([]);
    const hasListControl = control.hasControl("admin.system.driver.file:list");
    const hasUpdateControl = control.hasControl("admin.system.driver.file:update");
    const hasDeleteControl = control.hasControl("admin.system.driver.file:delete");
    const fileEnums = enumData?.["driver.file"] || {};
    const sceneEnums = enumData?.["driver.scene"] || {};
    const channelEnums = enumData?.["driver.channel"] || {};
    const statusOptions = buildOptions(fileEnums?.status, true);
    const statusValueEnum = buildValueEnum(fileEnums?.status, {"0": "green", "1": "red"});
    const sceneValueEnum = buildValueEnum(sceneEnums?.name);
    const channelValueEnum = buildValueEnum(channelEnums?.name);
    const sceneOptions = buildOptions(sceneEnums?.name, true);

    const openUpload = () => fileUploadRef.current?.open();

    const onUploaded = () => {
        clearEnumCache("driver.file");
        tableRef.current?.refreshTableData?.();
    };

    const openEdit = (ids: any[] = [], status?: any) => {
        if (!hasUpdateControl) {
            message?.warning("暂无修改权限");
            return;
        }
        const nextIds = Array.from(new Set((ids || []).map((item) => Number(item)).filter((item) => item > 0)));
        if (nextIds.length <= 0) {
            message?.warning("请选择要修改的文件");
            return;
        }
        setEditIds(nextIds);
        form.resetFields();
        form.setFieldsValue({
            status: status === undefined ? undefined : Number(status),
        });
        modalRef.current?.show();
    };

    const onDelete = (ids: any[] = []) => {
        if (!hasDeleteControl) {
            message?.warning("暂无删除权限");
            return;
        }
        const nextIds = Array.from(new Set((ids || []).map((item) => Number(item)).filter((item) => item > 0)));
        if (nextIds.length <= 0) {
            message?.warning("请选择要删除的文件");
            return;
        }
        let clear = false;
        modal.confirm({
            title: "删除文件",
            content: (
                <Space orientation="vertical">
                    <span>{`确定删除选中的 ${nextIds.length} 个文件吗？`}</span>
                    <Checkbox onChange={(e) => {
                        clear = e.target.checked;
                    }}>
                        同时清理文件对象
                    </Checkbox>
                </Space>
            ),
            okText: "确 定",
            okType: "danger",
            cancelText: "取 消",
            mask: {
                closable: true,
            },
            onOk: async () => {
                await deleteFile({
                    body: {
                        ids: nextIds,
                        clear,
                    },
                    onSuccess: (res: any) => {
                        clearEnumCache("driver.file");
                        tableRef.current?.refreshTableData?.();
                        tableRef.current?.clearSelectedRows?.();
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
        if (!hasUpdateControl) {
            message?.warning("暂无修改权限");
            return;
        }
        setSubmitLoading(true);
        await updateFile({
            body: {
                ids: editIds,
                status: values?.status?.toString(),
            },
            onSuccess: (res: any) => {
                modalRef.current?.hide();
                tableRef.current?.refreshTableData?.();
                tableRef.current?.clearSelectedRows?.();
                message?.success(res?.message || "更新成功");
            },
            onFail: (res: any) => {
                message?.error(res?.message || "更新失败");
            },
            onFinally: () => setSubmitLoading(false),
        });
    };

    const columns: any[] = [
        {
            title: "ID",
            dataIndex: "id",
            tip: "文件ID",
            width: 80,
            sorter: true,
            hideInSearch: true,
        },
        {
            title: "文件名称",
            dataIndex: "name",
            tip: "文件名称",
            formItemProps: {
                label: "文件名称",
            },
            valueType: "text",
            copyable: true,
        },
        {
            title: "文件后缀",
            dataIndex: "ext",
            tip: "文件后缀",
            formItemProps: {
                label: "文件后缀",
            },
            valueType: "text",
            width: 100,
        },
        {
            title: "文件大小",
            dataIndex: "size",
            tip: "文件大小",
            hideInSearch: true,
            render: (text: any) => formatSize(text),
        },
        {
            title: "通道",
            dataIndex: "channel_id",
            tip: "文件所属存储通道",
            formItemProps: {
                label: "存储通道",
            },
            valueEnum: channelValueEnum,
        },
        {
            title: "场景",
            dataIndex: "scene_id",
            tip: "文件所属存储场景",
            formItemProps: {
                label: "存储场景",
            },
            valueEnum: sceneValueEnum,
        },
        {
            title: "哈希",
            dataIndex: "hash",
            tip: "文件哈希值",
            formItemProps: {
                label: "文件哈希",
            },
            valueType: "text",
            copyable: true,
        },
        {
            title: "路径",
            dataIndex: "path",
            tip: "文件存储路径",
            formItemProps: {
                label: "文件路径",
            },
            valueType: "text",
            copyable: true,
        },
        {
            title: "状态",
            dataIndex: "status",
            tip: "文件状态",
            formItemProps: {
                label: "文件状态",
            },
            valueEnum: statusValueEnum,
        },
        {
            title: "创建时间",
            valueType: "dateRange",
            dataIndex: "create_time",
            tip: "文件创建时间",
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
            tip: "文件更新时间",
            formItemProps: {
                label: "更新时间",
            },
            sorter: true,
            transform: dateRangeTransform("update_time"),
        },
        ...(hasUpdateControl || hasDeleteControl ? [{
            title: "操作",
            tip: "文件操作",
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
                        onClick: () => openEdit([record?.id], record?.status),
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
    ];

    if (!hasListControl) {
        return null;
    }

    return (
        <>
            <ProTable
                ref={tableRef}
                extraRefreshBtn={true}
                title={<Title>文件管理</Title>}
                rowKey="id"
                columns={columns}
                defaultPage={1}
                defaultPageSize={10}
                extra={(
                    <Button type="primary" onClick={openUpload}>
                        上传
                    </Button>
                )}
                rowSelection={(hasUpdateControl || hasDeleteControl) ? {
                    fixed: true,
                    rightExtra: [
                        hasDeleteControl ? (
                            <Button key="delete" type="primary" danger ghost
                                    onClick={() => onDelete(tableRef.current?.getSelectedRowKeys?.() as any)}>
                                批量删除
                            </Button>
                        ) : null,
                        hasUpdateControl ? (
                            <Button key="edit" type="primary" ghost
                                    onClick={() => openEdit(tableRef.current?.getSelectedRowKeys?.() as any)}>
                                批量编辑
                            </Button>
                        ) : null,
                    ].filter(Boolean),
                } : false}
                request={(params, sort) => getData(params, sort, getFileList)}
                paginationAffix={true}
                selectionAffix={true}
            />

            <ProModal
                ref={modalRef}
                title={<Title>编辑文件</Title>}
                onOk={form.submit}
                width={320}
                modalProps={{
                    confirmLoading: submitLoading,
                    okText: "提交",
                    cancelText: "取消",
                    afterClose: () => {
                        form.resetFields();
                        setEditIds([]);
                        setSubmitLoading(false);
                    },
                } as any}>
                <Form form={form} layout="vertical" onFinish={onSubmit}>
                    <Form.Item name="status" label="文件状态" rules={[{required: true, message: "请选择文件状态"}]}>
                        <Select placeholder="请选择文件状态" options={statusOptions as any}/>
                    </Form.Item>
                </Form>
            </ProModal>

            <FileUploadModal
                ref={fileUploadRef}
                sceneOptions={sceneOptions}
                onUploaded={onUploaded}
            />
        </>
    );
};

export default FileManage;
