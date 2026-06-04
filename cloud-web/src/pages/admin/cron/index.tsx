import React, {useRef, useState} from 'react';
import {Body, ProModal, ProModalRef, ProTable, ProTableRef, Title} from "sinking-antd";
import {App, Button, Col, Dropdown, Form, Input, Result, Row, Select, Spin} from "antd";
import {createStyles} from "antd-style";
import defaultSettings from "@/../config/defaultSettings";
import {buildOptions, buildValueEnum, dateRangeTransform, emptyTableData, getData, normalizeIds} from "@/utils/page";
import useEnum from "@/utils/enum";
import useControl from "@/utils/control";
import {AdminControlCode} from "@/utils/permission";
import AceEditor from "@/pages/components/ace-editor";
import {
    createCron,
    deleteCron,
    getCronLog,
    getCronInfo,
    getCronList,
    runCron,
    updateCron
} from "@/service/admin/cron";

const AcePath = defaultSettings?.basePath + "ace/" || "/ace/";

const useStyles: any = createStyles(({token}: any) => ({
    ace: {
        ".ace_editor": {
            borderRadius: token.borderRadius + "px !important",
        },
    },
    formLoading: {
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        height: 260,
        width: "100%",
    },
}));

export default (): React.ReactNode => {
    const {message, modal} = App.useApp();
    const {styles} = useStyles();
    const [enumData, enumLoading] = useEnum(["cron.task"]) as any;
    const control = useControl();
    const hasListControl = control.hasControl(AdminControlCode.SYSTEM_CRON_LIST);
    const hasCreateControl = control.hasControl(AdminControlCode.SYSTEM_CRON_CREATE);
    const hasUpdateControl = control.hasControl(AdminControlCode.SYSTEM_CRON_UPDATE);
    const hasDeleteControl = control.hasControl(AdminControlCode.SYSTEM_CRON_DELETE);
    const hasRunControl = control.hasControl(AdminControlCode.SYSTEM_CRON_RUN);
    const hasLogControl = control.hasControl(AdminControlCode.SYSTEM_CRON_LOG);
    const taskEnums = enumData?.["cron.task"] || {};
    const taskTypes = taskEnums?.types || {};
    const taskStatus = taskEnums?.status || {};
    const taskConcurrents = taskEnums?.concurrents || {};

    const typeOptions = buildOptions(taskTypes, true);
    const statusOptions = buildOptions(taskStatus, true);
    const concurrentOptions = buildOptions(taskConcurrents, true);
    const typeValueEnum = buildValueEnum(taskTypes, {"0": "blue", "1": "purple"});
    const statusValueEnum = buildValueEnum(taskStatus, {"0": "green", "1": "red"});
    const concurrentValueEnum = buildValueEnum(taskConcurrents, {"0": "green", "1": "orange"});

    /**
     * 新增/编辑表单
     */
    const tableRef = useRef<ProTableRef>({} as ProTableRef);
    const formModalRef = useRef<ProModalRef>({} as ProModalRef);
    const [form] = Form.useForm();
    const [formScene, setFormScene] = useState<any>({
        action: "create",
        id: undefined,
        title: "新增任务",
    });
    const [formSubmitting, setFormSubmitting] = useState(false);
    const [formInfoLoading, setFormInfoLoading] = useState(false);
    const [scriptContent, setScriptContent] = useState("");
    const [taskType, setTaskType] = useState<number>(0);
    const isEditMode = formScene?.action === "edit";
    const isScriptTask = Number(taskType) === 0;

    const openCreate = () => {
        if (!hasCreateControl) {
            message?.warning("暂无创建权限");
            return;
        }
        setFormScene({
            action: "create",
            id: undefined,
            title: "新增任务",
        });
        setTaskType(0);
        setScriptContent("");
        form.resetFields();
        form.setFieldsValue({
            type: 0,
            status: 0,
            concurrent: 0,
        });
        setFormInfoLoading(false);
        formModalRef.current?.show();
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
            title: "编辑任务",
        });
        setTaskType(Number(record?.type || 0));
        setScriptContent("");
        form.resetFields();
        form.setFieldsValue({
            name: record?.name,
            type: Number(record?.type),
            spec: record?.spec,
            status: Number(record?.status),
            concurrent: Number(record?.concurrent),
        });
        setFormInfoLoading(true);
        formModalRef.current?.show();
        getCronInfo({
            body: {id},
            onSuccess: (res: any) => {
                const info = res?.data || {};
                const type = Number(info?.type || 0);
                setTaskType(type);
                setScriptContent(type === 0 ? (info?.script || "") : "");
                form.setFieldsValue({
                    name: info?.name,
                    type,
                    spec: info?.spec,
                    script: type === 0 ? undefined : info?.script,
                    status: Number(info?.status),
                    concurrent: Number(info?.concurrent),
                });
            },
            onFail: (res: any) => {
                message?.error(res?.message || "读取任务失败");
                formModalRef.current?.hide();
            },
        }).finally(() => {
            setFormInfoLoading(false)
        });
    };

    const handleSubmit = (values: any) => {
        if (isEditMode && !hasUpdateControl) {
            message?.warning("暂无修改权限");
            return;
        }
        if (!isEditMode && !hasCreateControl) {
            message?.warning("暂无创建权限");
            return;
        }
        const type = Number(values?.type);
        const script = type === 0 ? scriptContent : values?.script;
        if (!String(script || "").trim()) {
            message?.warning(type === 0 ? "请输入任务内容" : "请输入调用方法");
            return;
        }
        const body: any = {
            name: values?.name,
            spec: values?.spec,
            script,
        };
        if (isEditMode) {
            body.id = formScene?.id;
            body.type = String(type);
            body.status = values?.status?.toString();
            body.concurrent = values?.concurrent?.toString();
        } else {
            body.type = type;
            body.status = Number(values?.status);
            body.concurrent = Number(values?.concurrent);
        }
        setFormSubmitting(true);
        (isEditMode ? updateCron : createCron)({
            body,
            onSuccess: (res: any) => {
                formModalRef.current?.hide();
                tableRef.current?.refreshTableData?.();
                tableRef.current?.clearSelectedRows?.();
                message?.success(res?.message || "操作成功");
            },
            onFail: (res: any) => {
                message?.error(res?.message || "操作失败");
            },
        }).finally(() => {
            setFormSubmitting(false)
        });
    };

    /**
     * 删除任务
     * @param ids 任务ID
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
            title: "删除任务",
            content: `确定删除选中的 ${nextIds.length} 个任务吗？`,
            okText: "确 定",
            okType: "danger",
            cancelText: "取 消",
            maskClosable: true,
            onOk: () => deleteCron({
                body: {ids: nextIds},
                onSuccess: (res: any) => {
                    tableRef.current?.refreshTableData?.();
                    tableRef.current?.clearSelectedRows?.();
                    message?.success(res?.message || "删除成功");
                },
                onFail: (res: any) => {
                    message?.error(res?.message || "删除失败");
                },
            }),
        } as any);
    };
    /**
     * 执行任务
     * @param id 任务ID
     */
    const onRun = async (id: any) => {
        if (!hasRunControl) {
            message?.warning("暂无执行权限");
            return;
        }
        const taskId = Number(id);
        if (!taskId) {
            return;
        }
        await runCron({
            body: {id: taskId},
            onSuccess: (res: any) => {
                message?.success(res?.message || "执行成功");
            },
            onFail: (res: any) => {
                message?.error(res?.message || "执行失败");
            },
        });
    };

    /**
     * 任务日志
     */
    const logModalRef = useRef<ProModalRef>({} as ProModalRef);
    const logTableRef = useRef<ProTableRef>({} as ProTableRef);
    const [logScene, setLogScene] = useState<any>({
        id: undefined,
        title: "任务日志",
    });

    const openLog = (record: any) => {
        if (!hasLogControl) {
            message?.warning("暂无任务日志权限");
            return;
        }
        const id = Number(record?.id);
        if (!id) {
            return;
        }
        setLogScene({
            id,
            title: `任务日志`,
        });
        logModalRef.current?.show();
    };

    const logColumns: any[] = [
        {
            title: "ID",
            dataIndex: "id",
            width: 80,
            sorter: true,
        },
        {
            title: "日志内容",
            dataIndex: "content",
        },
        {
            title: "创建时间",
            dataIndex: "create_time",
            sorter: true,
        },
    ] as any;

    /**
     * 表格
     */
    const columns: any[] = [
        {
            title: "ID",
            dataIndex: "id",
            tip: "任务ID",
            width: 80,
            sorter: true,
            hideInSearch: true,
        },
        {
            title: "任务名称",
            dataIndex: "name",
            tip: "任务名称",
            valueType: "text",
        },
        {
            title: "任务类型",
            dataIndex: "type",
            tip: "任务类型",
            valueEnum: typeValueEnum,
        },
        {
            title: "任务表达式",
            dataIndex: "spec",
            tip: "任务表达式",
            valueType: "text",
            copyable: true,
            hideInSearch: true,
        },
        {
            title: "允许并发",
            dataIndex: "concurrent",
            tip: "允许并发",
            valueEnum: concurrentValueEnum,
        },
        {
            title: "任务状态",
            dataIndex: "status",
            tip: "任务状态",
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
        ...(hasRunControl || hasLogControl || hasUpdateControl || hasDeleteControl ? [{
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
                    hasLogControl ? {
                        key: "log",
                        label: "日志",
                        onClick: () => openLog(record),
                    } : null,
                    hasRunControl ? {
                        key: "run",
                        label: "执行",
                        onClick: () => onRun(record?.id),
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
                    title={<Title>任务管理</Title>}
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
                            return Promise.resolve(emptyTableData);
                        }
                        return getData(params, sort, getCronList);
                    }}
                    paginationAffix={true}
                    selectionAffix={true}
                />

                <ProModal
                    ref={formModalRef}
                    title={<Title>{formScene?.title}</Title>}
                    onOk={form.submit}
                    width={700}
                    modalProps={{
                        confirmLoading: formSubmitting || formInfoLoading,
                        forceRender: true,
                        okText: "提交",
                        cancelText: "取消",
                        footer: formInfoLoading ? null : undefined,
                        style: {top: 30},
                        afterClose: () => {
                            form.resetFields();
                            setFormSubmitting(false);
                            setFormInfoLoading(false);
                            setScriptContent("");
                            setTaskType(0);
                        },
                    } as any}>
                    <Form form={form} layout="vertical" onFinish={handleSubmit}>
                        {formInfoLoading ? (
                            <div className={styles.formLoading}>
                                <Spin/>
                            </div>
                        ) : (
                            <>
                                <Row gutter={16}>
                                    <Col xs={24} sm={12}>
                                        <Form.Item name="name" label="任务名称"
                                                   rules={[{required: true, message: "请输入任务名称"}]}>
                                            <Input placeholder="请输入任务名称" maxLength={100} showCount/>
                                        </Form.Item>
                                    </Col>
                                    <Col xs={24} sm={12}>
                                        <Form.Item name="spec" label="任务表达式"
                                                   rules={[{required: true, message: "请输入任务表达式"}]}>
                                            <Input placeholder="请输入任务表达式，例如 */5 * * * *"/>
                                        </Form.Item>
                                    </Col>
                                </Row>
                                <Row gutter={16}>
                                    <Col xs={24} sm={8}>
                                        <Form.Item name="type" label="任务类型"
                                                   rules={[{required: true, message: "请选择任务类型"}]}>
                                            <Select placeholder="请选择任务类型"
                                                    options={typeOptions as any}
                                                    onChange={(value) => {
                                                        const nextType = Number(value);
                                                        setTaskType(nextType);
                                                        setScriptContent("");
                                                        form.setFieldsValue({script: undefined});
                                                    }}/>
                                        </Form.Item>
                                    </Col>
                                    <Col xs={24} sm={8}>
                                        <Form.Item name="status" label="任务状态"
                                                   rules={[{required: true, message: "请选择任务状态"}]}>
                                            <Select placeholder="请选择任务状态" options={statusOptions as any}/>
                                        </Form.Item>
                                    </Col>
                                    <Col xs={24} sm={8}>
                                        <Form.Item name="concurrent" label="允许并发"
                                                   rules={[{required: true, message: "请选择并发配置"}]}>
                                            <Select placeholder="请选择并发配置" options={concurrentOptions as any}/>
                                        </Form.Item>
                                    </Col>
                                </Row>
                                {isScriptTask ? (
                                    <Form.Item label="任务内容" required>
                                        <AceEditor
                                            value={scriptContent}
                                            mode="sh"
                                            showPrintMargin={false}
                                            theme="monokai"
                                            width="100%"
                                            height={420}
                                            acePath={AcePath}
                                            onChange={(value: string) => setScriptContent(value)}
                                            className={styles?.ace}
                                        />
                                    </Form.Item>
                                ) : (
                                    <Form.Item name="script" label="调用方法"
                                               rules={[{required: true, message: "请输入调用方法"}]}>
                                        <Input placeholder="请输入调用方法，例如 service.function" maxLength={200}/>
                                    </Form.Item>
                                )}
                            </>
                        )}
                    </Form>
                </ProModal>

                <ProModal
                    ref={logModalRef}
                    title={<Title>{logScene?.title}</Title>}
                    width={760}
                    modalProps={{
                        forceRender: true,
                        footer: null,
                        mask: {
                            closable: true,
                        },
                        afterClose: () => {
                            logTableRef.current?.clearSelectedRows?.();
                            setLogScene({
                                id: undefined,
                                title: "任务日志",
                            });
                        },
                    } as any}>
                    <ProTable
                        key={logScene?.id || 0}
                        ref={logTableRef}
                        search={false}
                        extraRefreshBtn={false}
                        title={false}
                        rowKey="id"
                        columns={logColumns}
                        defaultPage={1}
                        defaultPageSize={10}
                        request={(params, sort) => {
                            if (!hasLogControl || !logScene?.id) {
                                return Promise.resolve(emptyTableData);
                            }
                            return getData({...params, id: logScene.id}, sort, getCronLog);
                        }}
                        paginationAffix={true}
                    />
                </ProModal>
            </>}
        </Body>
    );
};
