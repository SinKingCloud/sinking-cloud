import React, {forwardRef, useImperativeHandle, useLayoutEffect, useRef, useState} from "react";
import {App, Form, Input} from "antd";
import {useModel} from "umi";
import {ProModal, ProModalRef, Title} from "sinking-antd";
import {updatePassword} from "@/service/user/person";

export interface PasswordProps {
    ensureVerify?: () => Promise<boolean>;
    onSuccess?: (data: any) => void;
}

export interface PasswordRef {
    show: () => void;
}

const Password = forwardRef<PasswordRef, PasswordProps>((props, ref) => {
    const {ensureVerify, onSuccess} = props;
    const {message} = App.useApp();
    const user = useModel("user");
    const [form] = Form.useForm();
    const [formOpen, setFormOpen] = useState(false);
    const mountedForm = useRef<any>(null);
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const [loading, setLoading] = useState(false);

    /**
     * 通过前置安全验证后打开密码修改弹窗。
     */
    useImperativeHandle(ref, () => ({
        show: async () => {
            const pass = await ensureVerify?.();
            if (pass === false) {
                return;
            }
            modalRef.current?.show();
        },
    }));

    /**
     * 提交新密码并同步刷新当前前台用户信息。
     */
    const handleFinish = async (values: any) => {
        setLoading(true);
        await updatePassword({
            body: {
                password: values?.password,
            },
            onSuccess: (res: any) => {
                message?.success(res?.message || "密码修改成功");
                user?.refreshWebUser?.(() => undefined);
                onSuccess?.(values);
                modalRef.current?.hide();
            },
            onFail: (res: any) => message?.error(res?.message || "密码修改失败"),
            onFinally: () => setLoading(false),
        });
    };

    useLayoutEffect(() => {
        if (!formOpen || !mountedForm.current) return;
        form.resetFields();
    }, [form, formOpen]);

    return (
        <ProModal
            ref={modalRef}
            afterOpenChange={setFormOpen}
            title={<Title>修改密码</Title>}
            width={350}
            okText="确认"
            onOk={form.submit}
            onCancel={() => modalRef.current?.hide()}
            modalProps={{
                confirmLoading: loading,
                okButtonProps: {style: {minWidth: 72}},
                cancelButtonProps: {style: {minWidth: 72}},
                style: {
                    top: 100,
                    paddingBottom: 100,
                },
                destroyOnHidden: true,
            } as any}>
            <Form ref={mountedForm} form={form} clearOnDestroy layout="vertical" onFinish={handleFinish}>
                <Form.Item
                    label="新账户密码"
                    name="password"
                    rules={[
                        {required: true, message: "请输入新密码"},
                        {min: 6, max: 20, message: "密码长度为6-20个字符"},
                    ]}>
                    <Input.Password placeholder="请输入新密码" maxLength={20} showCount autoFocus/>
                </Form.Item>
                <Form.Item
                    label="确认新密码"
                    name="confirm_password"
                    dependencies={["password" as any]}
                    rules={[
                        {required: true, message: "请再次输入新密码"},
                        ({getFieldValue}: any) => ({
                            validator(_: any, value: string) {
                                if (!value || getFieldValue("password") === value) {
                                    return Promise.resolve();
                                }
                                return Promise.reject(new Error("两次输入的密码不一致"));
                            },
                        }),
                    ]}>
                    <Input.Password placeholder="请再次输入新密码" maxLength={20} showCount/>
                </Form.Item>
            </Form>
        </ProModal>
    );
});

export default Password;
