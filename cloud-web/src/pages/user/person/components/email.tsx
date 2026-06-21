import React, {forwardRef, useEffect, useImperativeHandle, useRef, useState} from "react";
import {App, Button, Col, Form, Input, Row} from "antd";
import {useModel} from "umi";
import {ProModal, ProModalRef, Title} from "sinking-antd";
import Captcha, {type CaptchaRef} from "@/pages/components/captcha";
import {sendMail} from "@/service/common/mail";
import {updateEmail} from "@/service/user/person";

export interface EmailProps {
    email?: string;
    ensureVerify?: () => Promise<boolean>;
    onSuccess?: (data: any) => void;
}

export interface EmailRef {
    show: () => void;
}

const Email = forwardRef<EmailRef, EmailProps>((props, ref) => {
    const {ensureVerify, onSuccess} = props;
    const {message} = App.useApp();
    const user = useModel("user");
    const [form] = Form.useForm();
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const captchaRef = useRef<CaptchaRef>({});
    const timerRef = useRef<number | undefined>(undefined);
    const [loading, setLoading] = useState(false);
    const [mailLoading, setMailLoading] = useState(false);
    const [countdown, setCountdown] = useState(0);

    /**
     * 清理验证码倒计时定时器。
     */
    const clearTimer = () => {
        if (timerRef.current) {
            window.clearInterval(timerRef.current);
            timerRef.current = undefined;
        }
        setCountdown(0);
    };

    /**
     * 启动邮箱验证码 60 秒倒计时。
     */
    const startCountdown = () => {
        clearTimer();
        let time = 60;
        setCountdown(time);
        timerRef.current = window.setInterval(() => {
            time -= 1;
            if (time <= 0) {
                clearTimer();
                return;
            }
            setCountdown(time);
        }, 1000);
    };

    /**
     * 组件卸载时释放倒计时定时器。
     */
    useEffect(() => () => {
        if (timerRef.current) {
            window.clearInterval(timerRef.current);
        }
    }, []);

    /**
     * 通过前置安全验证后打开邮箱修改弹窗。
     */
    useImperativeHandle(ref, () => ({
        show: async () => {
            const pass = await ensureVerify?.();
            if (pass === false) {
                return;
            }
            form.resetFields();
            clearTimer();
            modalRef.current?.show();
        },
    }));

    /**
     * 校验新邮箱并发送邮箱验证码。
     */
    const handleSendMail = () => {
        const mail = form.getFieldValue("mail" as any);
        if (!mail) {
            message?.error("请先输入新邮箱");
            return;
        }
        if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(mail)) {
            message?.error("请输入正确的邮箱地址");
            return;
        }
        if (!captchaRef.current?.Show) {
            return;
        }
        setMailLoading(true);
        captchaRef.current.Show(
            async (captchaRes) => {
                await sendMail({
                    body: {
                        mail,
                        token: captchaRes?.token,
                        captcha_x: captchaRes?.x,
                        captcha_y: captchaRes?.y,
                    },
                    onSuccess: (res: any) => {
                        message?.success(res?.message || "验证码发送成功");
                        startCountdown();
                    },
                    onFail: (res: any) => message?.error(res?.message || "验证码发送失败"),
                    onFinally: () => setMailLoading(false),
                });
            },
            () => setMailLoading(false)
        );
    };

    /**
     * 提交新邮箱和验证码并刷新当前前台用户信息。
     */
    const handleFinish = async (values: any) => {
        setLoading(true);
        await updateEmail({
            body: {
                mail: values?.mail,
                code: values?.code,
            },
            onSuccess: (res: any) => {
                message?.success(res?.message || "邮箱修改成功");
                user?.refreshWebUser?.(() => undefined);
                onSuccess?.(values);
                modalRef.current?.hide();
            },
            onFail: (res: any) => message?.error(res?.message || "邮箱修改失败"),
            onFinally: () => setLoading(false),
        });
    };

    return (
        <>
            <Captcha ref={captchaRef}/>
            <ProModal
                ref={modalRef}
                title={<Title>修改邮箱</Title>}
                width={350}
                okText="确认"
                onOk={form.submit}
                onCancel={() => modalRef.current?.hide()}
                afterClose={() => {
                    form.resetFields();
                    clearTimer();
                }}
                modalProps={{
                    confirmLoading: loading,
                    okButtonProps: {style: {minWidth: 72}},
                    cancelButtonProps: {style: {minWidth: 72}},
                    style: {
                        top: 100,
                        paddingBottom: 100,
                    },
                    destroyOnHidden: true,
                    forceRender: true,
                } as any}>
                <Form form={form} layout="vertical" onFinish={handleFinish}>
                    <Form.Item
                        label="邮箱地址"
                        name="mail"
                        rules={[
                            {required: true, message: "请输入邮箱"},
                            {type: "email", message: "请输入正确的邮箱地址"},
                        ]}>
                        <Input placeholder="请输入新的邮箱" maxLength={100} autoFocus/>
                    </Form.Item>
                    <Form.Item
                        label="验证码"
                        name="code"
                        rules={[
                            {required: true, message: "请输入邮箱验证码"},
                            {len: 6, message: "验证码为6位数字"},
                        ]}>
                        <Row gutter={8}>
                            <Col span={14}>
                                <Input placeholder="请输入邮箱验证码" maxLength={6}/>
                            </Col>
                            <Col span={10}>
                                <Button loading={mailLoading} disabled={countdown > 0} onClick={handleSendMail} block>
                                    {countdown > 0 ? `${countdown}s后重新获取` : "获取验证码"}
                                </Button>
                            </Col>
                        </Row>
                    </Form.Item>
                </Form>
            </ProModal>
        </>
    );
});

export default Email;
