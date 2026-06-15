import React, {forwardRef, useImperativeHandle, useRef, useState} from "react";
import {App, Avatar, Col, Form, Input, Row, Select} from "antd";
import {UserOutlined} from "@ant-design/icons";
import {ProModal, ProModalRef, Title} from "sinking-antd";
import {createStyles} from "antd-style";
import {AVATAR_SCENE_ID, getUploadUrl} from "@/utils/upload";
import {Upload, type UploadRef} from "@/pages/components/upload";
import {createMaster, updateMaster} from "@/service/admin/master";
import {AdminControlCode} from "@/utils/permission";

export type InfoRef = {
    openCreate: () => void;
    openEdit: (record: any) => void;
};

type InfoProps = {
    control: any;
    statusOptions: any[];
    onSuccess?: () => void;
};

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
}));

const Info = forwardRef<InfoRef, InfoProps>(({control, statusOptions, onSuccess}, ref) => {
    const {message} = App.useApp();
    const {styles} = useStyles();
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
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
    const hasCreateControl = control.hasControl(AdminControlCode.SYSTEM_MASTER_CREATE);
    const hasUpdateControl = control.hasControl(AdminControlCode.SYSTEM_MASTER_UPDATE);
    const isCreateMode = formScene?.action === "create";

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
        modalRef.current?.show();
    };

    useImperativeHandle(ref, () => ({
        openCreate,
        openEdit,
    }));

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
                            onSuccess?.();
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
                modalRef.current?.hide();
                onSuccess?.();
                message?.success(res?.message || "操作成功");
            },
            onFail: (res: any) => {
                message?.error(res?.message || "操作失败");
            },
            onFinally: () => setFormSubmitting(false),
        });
    };

    return (
        <ProModal
            ref={modalRef}
            title={<Title>{formScene?.title}</Title>}
            onOk={form.submit}
            width={460}
            modalProps={{
                confirmLoading: formSubmitting,
                okText: "提交",
                cancelText: "取消",
                forceRender: true,
                style: {
                    top: 100,
                    paddingBottom: 100,
                },
                afterClose: () => {
                    form.resetFields();
                    setFormSubmitting(false);
                    setAvatarUploading(false);
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
    );
});

Info.displayName = "Info";

export default Info;
