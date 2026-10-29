import React, {forwardRef, useEffect, useImperativeHandle, useLayoutEffect, useRef, useState} from "react";
import {App, Button, Col, Form, Input, Row} from "antd";
import {useModel} from "umi";
import {ProModal, ProModalRef, Title} from "sinking-antd";
import Captcha, {type CaptchaRef} from "@/pages/components/captcha";
import {sendSms} from "@/service/common/sms";
import {updatePhone} from "@/service/admin/person";

export interface PhoneProps {
    phone?: string;
    ensureVerify?: () => Promise<boolean>;
    onSuccess?: (data: any) => void;
}

export interface PhoneRef {
    show: () => void;
}

const Phone = forwardRef<PhoneRef, PhoneProps>((props, ref) => {
    const {ensureVerify, onSuccess} = props;
    const {message} = App.useApp();
    const user = useModel("user");
    const [form] = Form.useForm();
    const [formOpen, setFormOpen] = useState(false);
    const mountedForm = useRef<any>(null);
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const captchaRef = useRef<CaptchaRef>({});
    const timerRef = useRef<number | undefined>(undefined);
    const [loading, setLoading] = useState(false);
    const [smsLoading, setSmsLoading] = useState(false);
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
     * 启动短信验证码 60 秒倒计时。
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
     * 通过前置安全验证后打开手机号修改弹窗。
     */
    useImperativeHandle(ref, () => ({
        show: async () => {
            const pass = await ensureVerify?.();
            if (pass === false) {
                return;
            }
            clearTimer();
            modalRef.current?.show();
        },
    }));

    /**
     * 校验新手机号并发送短信验证码。
     */
    const handleSendSms = () => {
        const nextPhone = form.getFieldValue("phone" as any);
        if (!nextPhone) {
            message?.error("请先输入新手机号");
            return;
        }
        if (!/^1[3-9]\d{9}$/.test(nextPhone)) {
            message?.error("请输入正确的手机号格式");
            return;
        }
        if (!captchaRef.current?.Show) {
            return;
        }
        setSmsLoading(true);
        captchaRef.current.Show(
            async (captchaRes) => {
                await sendSms({
                    body: {
                        phone: nextPhone,
                        token: captchaRes?.token,
                        captcha_x: captchaRes?.x,
                        captcha_y: captchaRes?.y,
                        is_admin: true,
                    },
                    onSuccess: (res: any) => {
                        message?.success(res?.message || "验证码发送成功");
                        startCountdown();
                    },
                    onFail: (res: any) => message?.error(res?.message || "验证码发送失败"),
                    onFinally: () => setSmsLoading(false),
                });
            },
            () => setSmsLoading(false)
        );
    };

    /**
     * 提交新手机号和验证码并刷新当前后台用户信息。
     */
    const handleFinish = async (values: any) => {
        setLoading(true);
        await updatePhone({
            body: {
                phone: values?.phone,
                code: values?.code,
            },
            onSuccess: (res: any) => {
                message?.success(res?.message || "手机号修改成功");
                user?.refreshSystemUser?.(() => undefined);
                onSuccess?.(values);
                modalRef.current?.hide();
            },
            onFail: (res: any) => message?.error(res?.message || "手机号修改失败"),
            onFinally: () => setLoading(false),
        });
    };

    useLayoutEffect(() => {
        if (formOpen && mountedForm.current) form.resetFields();
    }, [form, formOpen]);

    return (
        <>
            <Captcha ref={captchaRef}/>
            <ProModal
                ref={modalRef}
                afterOpenChange={setFormOpen}
                title={<Title>修改手机</Title>}
                width={350}
                okText="确认"
                onOk={form.submit}
                onCancel={() => modalRef.current?.hide()}
                afterClose={() => {
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
                } as any}>
                <Form ref={mountedForm} form={form} clearOnDestroy layout="vertical" onFinish={handleFinish}>
                    <Form.Item
                        label="手机号"
                        name="phone"
                        rules={[
                            {required: true, message: "请输入手机号"},
                            {pattern: /^1[3-9]\d{9}$/, message: "请输入正确的手机号格式"},
                        ]}>
                        <Input placeholder="请输入新的手机号" maxLength={11} autoFocus/>
                    </Form.Item>
                    <Form.Item
                        label="验证码"
                        name="code"
                        rules={[
                            {required: true, message: "请输入验证码"},
                            {len: 6, message: "验证码为6位数字"},
                        ]}>
                        <Row gutter={8}>
                            <Col span={14}>
                                <Input placeholder="请输入验证码" maxLength={6}/>
                            </Col>
                            <Col span={10}>
                                <Button loading={smsLoading} disabled={countdown > 0} onClick={handleSendSms} block>
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

export default Phone;
