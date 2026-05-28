import React, {useRef, useState} from "react";
import {Body, ProModal, ProModalRef, ProTable, ProTableRef, Title} from "sinking-antd";
import {App, Avatar, Button, Col, Dropdown, Form, Input, Result, Row, Select, Space, Typography} from "antd";
import {UserOutlined} from "@ant-design/icons";
import {getData} from "@/utils/page";
import {AVATAR_SCENE_ID, getUploadUrl} from "@/utils/upload";
import useControl from "@/utils/control";
import useEnum from "@/utils/enum";
import {Upload, type UploadRef} from "@/pages/components/upload";
import {
    createMaster,
    deleteMaster,
    getMasterList,
    updateMaster,
} from "@/service/admin/master";
import MasterSession, {type MasterSessionRef} from "./components/session";
import MasterRole, {type MasterRoleRef} from "./components/role";
import {
    buildOptions,
    buildValueEnum,
    dateRangeTransform,
    emptyTableData,
    normalizeIds,
} from "./components/common";
import {createStyles} from "antd-style";

const useStyles: any = createStyles(({token}: any): any => ({
    avatarContainer: {
        display: "flex",
        alignItems: "center",
        gap: 12,
        padding: 12,
        marginBottom: 16,
        borderRadius: token.borderRadius,
        background: token.colorFillAlter,
        border: `1px solid ${token.colorBorderSecondary}`,
    },
    avatarInfo: {
        minWidth: 0,
    },
    avatarUpload: {
        position: "relative",
        cursor: "pointer",
        borderRadius: "50%",
        overflow: "hidden",
        flex: "none",
        "&:hover .upload-hint": {
            opacity: 1,
        },
    },
    uploadHint: {
        position: "absolute",
        left: 0,
        right: 0,
        bottom: 0,
        padding: "4px 0",
        color: "#fff",
        fontSize: 12,
        lineHeight: "16px",
        textAlign: "center",
        background: "rgba(0, 0, 0, 0.6)",
        opacity: 0,
        transition: "opacity 0.2s",
    },
    avatarTitle: {
        fontSize: 14,
        fontWeight: 500,
    },
    avatarDesc: {
        marginTop: 4,
        fontSize: 12,
        color: token.colorTextTertiary,
    },
    form: {
        ".ant-form-item": {
            marginBottom: 16,
        },
    },
    infoColumn: {
        width: "100%",
        minWidth: 180,
    },
    infoRow: {
        display: "flex",
        alignItems: "center",
        gap: 4,
        minWidth: 0,
    },
    labelText: {
        flex: "none",
        minWidth: 40,
        fontSize: 12,
    },
    smallText: {
        flex: 1,
        minWidth: 0,
        fontSize: 12,
        whiteSpace: "nowrap",
    },
}));

export default (): React.ReactNode => {
    const {message, modal} = App.useApp();
    const {styles} = useStyles();
    const [enumData, enumLoading] = useEnum(["admin.user", "admin.role"]) as any;
    const control = useControl();
    const hasListControl = control.hasControl("admin.system.master:list");
    const hasCreateControl = control.hasControl("admin.system.master:create");
    const hasUpdateControl = control.hasControl("admin.system.master:update");
    const hasDeleteControl = control.hasControl("admin.system.master:delete");
    const hasRoleControl = control.hasControl("admin.system.master:role");
    const hasSessionControl = control.hasControl("admin.system.master:session");
    const tableRef = useRef<ProTableRef>({} as ProTableRef);
    const formModalRef = useRef<ProModalRef>({} as ProModalRef);
    const roleRef = useRef<MasterRoleRef>({} as MasterRoleRef);
    const sessionRef = useRef<MasterSessionRef>({} as MasterSessionRef);
    const avatarUploadRef = useRef<UploadRef>({} as UploadRef);
    const [form] = Form.useForm();
    const avatarValue = Form.useWatch("avatar" as any, form);
    const [formSubmitting, setFormSubmitting] = useState(false);
    const [avatarUploading, setAvatarUploading] = useState(false);
    const [formScene, setFormScene] = useState<any>({
        action: "create",
        id: undefined,
        title: "新增账户",
    });
    const isCreateMode = formScene?.action === "create";
    const userStatus = enumData?.["admin.user"]?.status || {};
    const roleNames = enumData?.["admin.role"]?.names || {};
    const statusOptions = buildOptions(userStatus);
    const statusValueEnum = buildValueEnum(userStatus, {"0": "green", "1": "red"});

    const uploadAvatar = () => {
        if (avatarUploading) {
            return;
        }
        setAvatarUploading(true);
        const task = avatarUploadRef.current?.upload?.({
            sceneId: AVATAR_SCENE_ID,
            accept: "image/*",
            onComplete: async (result: any) => {
                const url = getUploadUrl(result);
                if (!url) {
                    setAvatarUploading(false);
                    message?.warning("头像上传成功，但未返回地址");
                    return;
                }
                form.setFieldValue("avatar" as any, url);
                if (isCreateMode || !formScene?.id) {
                    message?.success("头像上传成功");
                    return;
                }
                try {
                    await updateMaster({
                        body: {
                            id: Number(formScene.id),
                            avatar: url,
                        },
                        onSuccess: (res: any) => {
                            tableRef.current?.refreshTableData?.();
                            message?.success(res?.message || "头像保存成功");
                        },
                        onFail: (res: any) => {
                            message?.error(res?.message || "头像保存失败");
                        },
                        onFinally: () => setAvatarUploading(false),
                    });
                } catch (_) {
                    setAvatarUploading(false);
                }
            },
            onError: (error: Error) => {
                message?.error(error?.message || "头像上传失败");
            },
        } as any);
        if (!task) {
            setAvatarUploading(false);
            return;
        }
        task.promise.then((result) => {
            if (!result || isCreateMode || !formScene?.id) {
                setAvatarUploading(false);
            }
        }).catch(() => setAvatarUploading(false));
    };

    const openCreate = () => {
        if (!hasCreateControl) {
            message?.warning("暂无创建权限");
            return;
        }
        setFormScene({
            action: "create",
            id: undefined,
            title: "新增账户",
        });
        form.resetFields();
        form.setFieldsValue({status: 0});
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
            title: "编辑账户",
        });
        form.resetFields();
        form.setFieldsValue({
            account: record?.account,
            phone: record?.phone,
            email: record?.email,
            avatar: record?.avatar,
            nick_name: record?.nick_name,
            status: Number(record?.status),
        });
        formModalRef.current?.show();
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
        const body: any = {
            account: values?.account,
            phone: values?.phone,
            email: values?.email,
            password: values?.password,
            avatar: values?.avatar,
            nick_name: values?.nick_name,
        };
        if (isCreateMode) {
            body.status = Number(values?.status);
        } else {
            body.id = Number(formScene?.id);
            body.status = values?.status?.toString();
        }
        await (isCreateMode ? createMaster : updateMaster)({
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
            onFinally: () => setFormSubmitting(false),
        });
    };

    const onDelete = (id: any) => {
        if (!hasDeleteControl) {
            message?.warning("暂无删除权限");
            return;
        }
        const nextId = Number(id);
        if (!nextId) {
            message?.warning("请选择要删除的员工");
            return;
        }
        modal.confirm({
            title: "删除员工",
            content: "确定删除该员工吗？",
            okText: "确 定",
            okType: "danger",
            cancelText: "取 消",
            mask: {
                closable: true,
            },
            onOk: async () => {
                await deleteMaster({
                    body: {id: nextId},
                    onSuccess: (res: any) => {
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

    const roleRender = (ids: any[] = []) => {
        const names = normalizeIds(ids).map((id) => roleNames?.[id]).filter(Boolean);
        if (names.length <= 0) {
            return "-";
        }
        return (
            <Space size={4} wrap>
                {names.map((name) => (
                    <Typography.Text key={name}>{name}</Typography.Text>
                ))}
            </Space>
        );
    };

    const renderRows = (items: { label: string; value: any }[], copyable = true) => (
        <Space orientation="vertical" size={2} className={styles.infoColumn}>
            {items.map((item) => (
                <div className={styles.infoRow} key={item.label}>
                    <Typography.Text type="secondary" className={styles.labelText}>{item.label}</Typography.Text>
                    <Typography.Text
                        className={styles.smallText}
                        copyable={copyable && item.value ? {text: item.value} : false}>
                        {item.value || "-"}
                    </Typography.Text>
                </div>
            ))}
        </Space>
    );

    const columns: any[] = [
        {
            title: "ID",
            dataIndex: "id",
            tip: "员工ID",
            sorter: true,
            hideInSearch: true,
        },
        {
            title: "头像",
            dataIndex: "avatar",
            tip: "员工头像",
            hideInSearch: true,
            render: (text: any) => <Avatar src={text || undefined} icon={<UserOutlined/>}/>,
        },
        {
            title: "账号信息",
            dataIndex: "account_info",
            tip: "账号/昵称",
            hideInSearch: true,
            render: (_: any, record: any) => renderRows([
                {label: "账号:", value: record?.account},
                {label: "昵称:", value: record?.nick_name},
            ]),
        },
        {
            title: "账号",
            dataIndex: "account",
            tip: "账号",
            valueType: "text",
            hideInTable: true,
        },
        {
            title: "昵称",
            dataIndex: "nick_name",
            tip: "昵称",
            valueType: "text",
            hideInTable: true,
        },
        {
            title: "联系方式",
            dataIndex: "contact_info",
            tip: "手机号/邮箱",
            hideInSearch: true,
            render: (_: any, record: any) => renderRows([
                {label: "手机:", value: record?.phone},
                {label: "邮箱:", value: record?.email},
            ]),
        },
        {
            title: "手机号",
            dataIndex: "phone",
            tip: "手机号",
            valueType: "text",
            hideInTable: true,
        },
        {
            title: "邮箱",
            dataIndex: "email",
            tip: "邮箱",
            valueType: "text",
            hideInTable: true,
        },
        {
            title: "角色",
            dataIndex: "role_ids",
            tip: "员工拥有的角色",
            hideInSearch: true,
            render: (_: any, record: any) => roleRender(record?.role_ids || []),
        },
        {
            title: "状态",
            dataIndex: "status",
            tip: "员工状态",
            valueEnum: statusValueEnum,
        },
        {
            title: "创建时间",
            valueType: "dateRange",
            dataIndex: "create_time",
            tip: "创建时间",
            hideInTable: true,
            sorter: true,
            transform: dateRangeTransform("create_time"),
        },
        {
            title: "更新时间",
            valueType: "dateRange",
            dataIndex: "update_time",
            tip: "更新时间",
            hideInTable: true,
            sorter: true,
            transform: dateRangeTransform("update_time"),
        },
        {
            title: "相关时间",
            dataIndex: "time_info",
            tip: "创建时间/更新时间",
            hideInSearch: true,
            render: (_: any, record: any) => renderRows([
                {label: "创建于:", value: record?.create_time},
                {label: "更新于:", value: record?.update_time},
            ], false),
        },
        ...(hasUpdateControl || hasDeleteControl || hasRoleControl || hasSessionControl ? [{
            title: "操作",
            valueType: "option",
            key: "option",
            fixed: "right",
            hideInSearch: true,
            render: (_: any, record: any) => {
                const items = [
                    hasUpdateControl ? {
                        key: "edit",
                        label: "编辑",
                        onClick: () => openEdit(record),
                    } : null,
                    hasRoleControl ? {
                        key: "role",
                        label: "角色",
                        onClick: () => roleRef.current?.open(record),
                    } : null,
                    hasSessionControl ? {
                        key: "session",
                        label: "会话",
                        onClick: () => sessionRef.current?.open(record),
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

    return (
        <Body loading={control.loading || enumLoading}>
            {!hasListControl && <Result status="403" title="403" subTitle="您无权限访问此页面,请联系管理员"/>}
            {hasListControl && <>
                <ProTable
                    ref={tableRef}
                    extraRefreshBtn={true}
                    title={<Title>员工列表</Title>}
                    rowKey="id"
                    columns={columns}
                    defaultPage={1}
                    defaultPageSize={20}
                    extra={hasCreateControl ? (
                        <Button type="primary" onClick={openCreate}>
                            新增
                        </Button>
                    ) : null}
                    request={(params, sort) => {
                        if (!hasListControl) {
                            return Promise.resolve(emptyTableData);
                        }
                        return getData(params, sort, getMasterList);
                    }}
                    paginationAffix={true}
                    selectionAffix={true}
                />

                <ProModal
                    ref={formModalRef}
                    title={<Title>{formScene?.title}</Title>}
                    onOk={form.submit}
                    width={460}
                    modalProps={{
                        confirmLoading: formSubmitting,
                        okText: "提交",
                        cancelText: "取消",
                        forceRender: true,
                        afterClose: () => {
                            form.resetFields();
                            setFormSubmitting(false);
                        },
                    } as any}
                >
                    <Form form={form} layout="vertical" onFinish={handleSubmit} className={styles.form}>
                        <div className={styles.avatarContainer}>
                            <Upload ref={avatarUploadRef}/>
                            <div className={styles.avatarUpload} onClick={uploadAvatar}>
                                <Avatar size={72} src={(avatarValue || undefined) as any} icon={<UserOutlined/>}/>
                                <div className={`${styles.uploadHint} upload-hint`}>上传</div>
                            </div>
                            <div className={styles.avatarInfo}>
                                <div className={styles.avatarTitle}>头像</div>
                                <div className={styles.avatarDesc}>点击头像上传图片</div>
                            </div>
                        </div>
                        <Form.Item name="avatar" hidden>
                            <Input/>
                        </Form.Item>
                        <Row gutter={16}>
                            <Col span={12}>
                                <Form.Item name="account" label="账号"
                                           rules={[{required: isCreateMode, message: "请输入账号"}]}>
                                    <Input placeholder="请输入账号" maxLength={32}/>
                                </Form.Item>
                            </Col>
                            <Col span={12}>
                                <Form.Item name="password" label="密码"
                                           rules={[{required: isCreateMode, message: "请输入登录密码"}]}>
                                    <Input placeholder={isCreateMode ? "请输入登录密码" : "留空则不修改密码"}
                                           maxLength={64}/>
                                </Form.Item>
                            </Col>
                        </Row>
                        <Row gutter={16}>
                            <Col span={12}>
                                <Form.Item name="status" label="状态"
                                           rules={[{required: isCreateMode, message: "请选择状态"}]}>
                                    <Select placeholder="请选择状态" options={statusOptions as any}/>
                                </Form.Item>
                            </Col>
                            <Col span={12}>
                                <Form.Item name="nick_name" label="昵称"
                                           rules={[{max: 20, message: "昵称不能超过20个字符"}]}>
                                    <Input placeholder="请输入昵称" maxLength={20} showCount/>
                                </Form.Item>
                            </Col>
                        </Row>
                        <Row gutter={16}>
                            <Col span={12}>
                                <Form.Item name="phone" label="手机"
                                           rules={[{pattern: /^1[3-9]\d{9}$/, message: "请输入正确的手机号"}]}>
                                    <Input placeholder="请输入手机号码" maxLength={11}/>
                                </Form.Item>
                            </Col>
                            <Col span={12}>
                                <Form.Item name="email" label="邮箱"
                                           rules={[{type: "email", message: "请输入正确的邮箱地址"}]}>
                                    <Input placeholder="请输入邮箱地址" maxLength={80}/>
                                </Form.Item>
                            </Col>
                        </Row>
                    </Form>
                </ProModal>

                <MasterRole ref={roleRef} control={control} onSuccess={() => tableRef.current?.refreshTableData?.()}/>
                <MasterSession ref={sessionRef} control={control}/>
            </>}
        </Body>
    );
};
