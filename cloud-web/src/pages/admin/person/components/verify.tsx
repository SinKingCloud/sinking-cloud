import React, {forwardRef, useEffect, useImperativeHandle, useRef, useState} from "react";
import {App, Button, Form, Input, Steps, Typography} from "antd";
import {CheckCircleFilled, LockOutlined, MailOutlined, MobileOutlined} from "@ant-design/icons";
import {createStyles} from "antd-style";
import {ProModal, ProModalRef, Title} from "sinking-antd";
import Captcha, {type CaptchaRef} from "@/pages/components/captcha";
import {sendMail} from "@/service/common/mail";
import {sendSms} from "@/service/common/sms";
import {verifySecurity} from "@/service/admin/person";

export type VerifyInfo = {
    status?: boolean;
    phone?: string;
    mail?: string;
    password?: boolean;
};

export type VerifyRef = {
    ensure: () => Promise<boolean>;
};

const useStyles: any = createStyles(({token, responsive}: any): any => ({
    verifyContent: {
        display: "grid",
        gap: 16,
        padding: "8px 8px 12px",
    },
    verifySteps: {
        padding: "8px 12px 0",

        [responsive.mobile]: {
            padding: "4px 0 0",
        },

        ".ant-steps-item": {
            minWidth: 0,
            overflow: "visible",
        },

        ".ant-steps-item-title": {
            display: "block",
            maxWidth: "100%",
            paddingInlineEnd: "0 !important",
            fontSize: "12px !important",
            lineHeight: "18px !important",
            textAlign: "center",
            whiteSpace: "nowrap",
        },

        ".ant-steps-item-icon": {
            display: "inline-flex !important",
            alignItems: "center",
            justifyContent: "center",
            flex: "none",
            float: "none !important",
            marginInlineStart: "0 !important",
            marginInlineEnd: "0 !important",
            marginTop: "0 !important",
            width: "20px !important",
            height: "20px !important",
            lineHeight: "1 !important",
            fontSize: "11px !important",
        },

        ".ant-steps-item-container": {
            display: "block !important",
            minWidth: 0,
            textAlign: "center",
        },

        ".ant-steps-icon": {
            display: "inline-flex",
            alignItems: "center",
            justifyContent: "center",
            lineHeight: "1 !important",
        },

        ".ant-steps-item-content": {
            display: "block !important",
            minWidth: 0,
            width: "auto !important",
            minHeight: "0 !important",
            marginTop: "4px !important",
            textAlign: "center",
        },

        ".ant-steps-item-tail": {
            top: "10px !important",
        },
    },
    panel: {
        display: "grid",
        alignContent: "start",
        gap: 16,
        padding: "4px 12px 0px",

        [responsive.mobile]: {
            paddingInline: 0,
        },
    },
    choiceHead: {
        padding: "0 4px",
    },
    choiceTitle: {
        color: token.colorTextSecondary,
        fontSize: 13,
        fontWeight: 600,
        lineHeight: "20px",
    },
    choiceDesc: {
        color: token.colorTextTertiary,
        fontSize: 12,
        lineHeight: "18px",
        marginTop: 4,
    },
    methodList: {
        display: "grid",
        gap: 8,
    },
    methodItem: {
        position: "relative",
        display: "flex",
        alignItems: "center",
        gap: 12,
        minHeight: 42,
        padding: "8px 12px",
        borderRadius: token.borderRadius,
        background: token.colorFillQuaternary,
        cursor: "pointer",
        transition: "none",

        "&:hover": {
            background: token.colorFillSecondary,
        },
    },
    methodActive: {
        background: `${token.colorPrimaryBg} !important`,
    },
    methodDisabled: {
        display: "none",
        opacity: 0.45,
        cursor: "not-allowed",

        "&:hover": {
            background: token.colorFillQuaternary,
        },
    },
    methodIcon: {
        width: 24,
        height: 24,
        flex: "none",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        borderRadius: token.borderRadiusSM,
        color: token.colorTextTertiary,
        background: "transparent",
        fontSize: 15,
    },
    methodIconActive: {
        color: token.colorTextSecondary,
    },
    methodBody: {
        flex: 1,
        minWidth: 0,
    },
    methodTitle: {
        color: token.colorText,
        fontSize: 13,
        fontWeight: 500,
        lineHeight: "20px",
    },
    methodValue: {
        color: token.colorTextTertiary,
        fontSize: 12,
        lineHeight: "18px",
        overflow: "hidden",
        textOverflow: "ellipsis",
        whiteSpace: "nowrap",
    },
    methodCheck: {
        width: 22,
        height: 22,
        flex: "none",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        color: token.colorPrimary,
        opacity: 0,
        transition: "none",
        fontSize: 14,
    },
    methodCheckActive: {
        opacity: 1,
    },
    targetBox: {
        display: "flex",
        alignItems: "center",
        justifyContent: "space-between",
        gap: 16,
        padding: "12px",
        borderRadius: token.borderRadius,
        background: token.colorFillQuaternary,

        [responsive.mobile]: {
            display: "block",
        },
    },
    targetLabel: {
        color: token.colorTextTertiary,
        fontSize: 12,
        lineHeight: "18px",
    },
    targetValue: {
        color: token.colorText,
        fontSize: 13,
        fontWeight: 500,
        lineHeight: "20px",
        overflow: "hidden",
        textOverflow: "ellipsis",
        whiteSpace: "nowrap",
    },
    targetRight: {
        minWidth: 0,
        textAlign: "right",

        [responsive.mobile]: {
            marginTop: 8,
            textAlign: "left",
        },
    },
    changeMethod: {
        height: "auto",
        padding: 0,
        fontSize: 12,
    },
    form: {
        ".ant-form-item": {
            marginBottom: 0,
        },

        ".ant-form-item-label": {
            paddingBottom: 4,
        },
    },
    codeRow: {
        display: "flex",
        alignItems: "flex-end",
        gap: 12,

        [responsive.mobile]: {
            display: "block",
        },
    },
    codeInputWrap: {
        flex: 1,
        minWidth: 0,
    },
    codeInput: {
        height: 32,
        borderRadius: token.borderRadius,
    },
    resendButton: {
        width: 96,
        height: 32,
        flex: "none",

        [responsive.mobile]: {
            width: "100%",
            marginTop: 12,
        },
    },
}));

const getDefaultType = (info?: VerifyInfo) => {
    if (info?.phone) {
        return "phone";
    }
    if (info?.mail) {
        return "mail";
    }
    if (info?.password) {
        return "password";
    }
    return "phone";
};

/**
 * 根据验证方式获取对应的绑定目标。
 */
const getTargetValue = (info: VerifyInfo, type: string) => {
    if (type === "password") {
        return info?.password ? "账户密码" : "";
    }
    return type === "mail" ? info?.mail : info?.phone;
};

/**
 * 将验证方式转换为界面展示名称。
 */
const getTypeName = (type: string) => {
    if (type === "password") {
        return "密码";
    }
    return type === "mail" ? "邮箱" : "手机";
};

const Verify = forwardRef<VerifyRef>((_, ref) => {
    const {message} = App.useApp();
    const {styles, cx} = useStyles();
    const [form] = Form.useForm();
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const captchaRef = useRef<CaptchaRef>({});
    const resolverRef = useRef<((value: boolean) => void) | null>(null);
    const timerRef = useRef<number | undefined>(undefined);
    const [info, setInfo] = useState<VerifyInfo>({});
    const [verifyType, setVerifyType] = useState<string>("phone");
    const [currentStep, setCurrentStep] = useState(0);
    const [primaryLoading, setPrimaryLoading] = useState(false);
    const [resendLoading, setResendLoading] = useState(false);
    const [countdown, setCountdown] = useState(0);

    /**
     * 清理验证码重发倒计时。
     */
    const clearCountdown = () => {
        if (timerRef.current) {
            window.clearInterval(timerRef.current);
            timerRef.current = undefined;
        }
        setCountdown(0);
    };

    /**
     * 结束验证流程并将结果回传给调用方。
     */
    const resolveAndClose = (value: boolean) => {
        resolverRef.current?.(value);
        resolverRef.current = null;
        modalRef.current?.hide();
    };

    /**
     * 启动验证码重发倒计时。
     */
    const startCountdown = () => {
        clearCountdown();
        let seconds = 60;
        setCountdown(seconds);
        timerRef.current = window.setInterval(() => {
            seconds -= 1;
            if (seconds <= 0) {
                clearCountdown();
                return;
            }
            setCountdown(seconds);
        }, 1000);
    };

    /**
     * 组件卸载时释放验证码倒计时。
     */
    useEffect(() => () => {
        clearCountdown();
    }, []);

    /**
     * 查询当前账号是否仍需要安全验证。
     */
    const loadVerifyInfo = async () => {
        const res = await verifySecurity({
            body: {
                action: "get",
            },
        });
        if (res?.code !== 200) {
            message?.error(res?.message || "获取验证状态失败");
            return undefined;
        }
        return res?.data || {};
    };

    /**
     * 暴露安全验证入口；无需验证时直接通过，需要验证时打开流程弹窗。
     */
    useImperativeHandle(ref, () => ({
        ensure: async () => {
            const nextInfo = await loadVerifyInfo();
            if (!nextInfo) {
                return false;
            }
            if (nextInfo?.status) {
                return true;
            }
            if (!nextInfo?.phone && !nextInfo?.mail && !nextInfo?.password) {
                message?.warning("当前账号未绑定手机号、邮箱或密码，无法完成安全验证");
                return false;
            }
            const nextType = getDefaultType(nextInfo);
            setInfo(nextInfo);
            setVerifyType(nextType);
            setCurrentStep(0);
            setPrimaryLoading(false);
            setResendLoading(false);
            clearCountdown();
            form.resetFields();
            modalRef.current?.show();
            return new Promise<boolean>((resolve) => {
                resolverRef.current = resolve;
            });
        },
    }));

    /**
     * 切换手机或邮箱验证方式，并重置当前验证码输入。
     */
    const selectType = (type: string) => {
        if (!getTargetValue(info, type)) {
            return;
        }
        setVerifyType(type);
        clearCountdown();
        form.resetFields(["code"] as any);
    };

    /**
     * 发送当前验证方式对应的验证码。
     */
    const sendCode = (nextStep: boolean, inline = false) => {
        if (countdown > 0 || primaryLoading || resendLoading) {
            return;
        }
        const target = getTargetValue(info, verifyType);
        if (!target) {
            message?.warning("当前验证方式不可用");
            return;
        }
        const setLoading = inline ? setResendLoading : setPrimaryLoading;
        const sendFn = verifyType === "mail" ? sendMail : sendSms;
        setLoading(true);
        if (!captchaRef.current?.Show) {
            setLoading(false);
            return;
        }
        captchaRef.current?.Show?.(
            async (res: any) => {
                await sendFn({
                    body: {
                        is_admin: true,
                        token: res?.token,
                        captcha_x: res?.x,
                        captcha_y: res?.y,
                    },
                    onSuccess: (r: any) => {
                        message?.success(r?.message || "验证码已发送");
                        startCountdown();
                        if (nextStep) {
                            form.resetFields(["code"] as any);
                            setCurrentStep(1);
                        }
                    },
                    onFail: (r: any) => message?.error(r?.message || "验证码发送失败"),
                    onFinally: () => setLoading(false),
                });
            },
            () => setLoading(false)
        );
    };

    /**
     * 提交验证码，验证通过后直接回到原修改流程。
     */
    const handleSubmit = async (values: any) => {
        setPrimaryLoading(true);
        await verifySecurity({
            body: {
                action: "verify",
                type: verifyType,
                code: values?.code,
                password: values?.password,
            },
            onSuccess: () => {
                setPrimaryLoading(false);
                clearCountdown();
                resolveAndClose(true);
            },
            onFail: (res: any) => message?.error(res?.message || "安全验证失败"),
            onFinally: () => setPrimaryLoading(false),
        });
    };

    /**
     * 根据当前步骤处理弹窗主按钮动作。
     */
    const handleOk = () => {
        if (currentStep === 0) {
            sendCode(true);
            return;
        }
        if (currentStep === 1) {
            form.submit();
            return;
        }
    };

    const isPasswordVerify = verifyType === "password";
    const okText = currentStep === 0 ? (isPasswordVerify ? "继续" : "发送") : "验证";
    const targetValue = getTargetValue(info, verifyType);

    /**
     * 渲染单个可选验证方式。
     */
    const renderMethod = (type: string, icon: React.ReactNode, title: string, value?: string) => {
        const active = verifyType === type;
        const disabled = !value;
        return (
            <div
                className={cx(styles.methodItem, active && styles.methodActive, disabled && styles.methodDisabled)}
                onClick={() => selectType(type)}>
                <span className={cx(styles.methodIcon, active && styles.methodIconActive)}>{icon}</span>
                <div className={styles.methodBody}>
                    <div className={styles.methodTitle}>{title}</div>
                    <div className={styles.methodValue}>{value || "未绑定"}</div>
                </div>
                <span className={cx(styles.methodCheck, active && styles.methodCheckActive)}>
                    <CheckCircleFilled/>
                </span>
            </div>
        );
    };

    return (
        <>
            <Captcha ref={captchaRef}/>
            <ProModal
                ref={modalRef}
                title={<Title>安全验证</Title>}
                width={420}
                okText={okText}
                onOk={handleOk}
                onCancel={() => resolveAndClose(false)}
                afterClose={() => {
                    setPrimaryLoading(false);
                    setResendLoading(false);
                    setCurrentStep(0);
                    clearCountdown();
                    form.resetFields();
                    if (resolverRef.current) {
                        resolverRef.current(false);
                        resolverRef.current = null;
                    }
                }}
                modalProps={{
                    confirmLoading: primaryLoading,
                    cancelText: "取消",
                    okButtonProps: {style: {minWidth: 72}},
                    cancelButtonProps: {style: {minWidth: 72}},
                    mask: {closable: false},
                    style: {
                        top: 100,
                        paddingBottom: 100,
                    },
                    destroyOnHidden: true,
                    forceRender: true,
                } as any}>
                <Form form={form} layout="vertical" className={styles.form} onFinish={handleSubmit}>
                    <div className={styles.verifyContent}>
                        <Steps
                            className={styles.verifySteps}
                            current={currentStep}
                            titlePlacement="vertical"
                            responsive={false}
                            size="small"
                            items={[
                                {title: "验证方式"},
                                {title: "安全验证"},
                                {title: "验证通过"},
                            ]}
                        />
                        {currentStep === 0 && (
                            <div className={styles.panel}>
                                <div className={styles.choiceHead}>
                                    <div className={styles.choiceTitle}>请选择验证方式</div>
                                    <div className={styles.choiceDesc}>
                                        {info?.password && !info?.phone && !info?.mail
                                            ? "当前账号未绑定手机或邮箱，可使用登录密码完成验证。"
                                            : "使用已绑定的手机或邮箱接收验证码。"}
                                    </div>
                                </div>
                                <div className={styles.methodList}>
                                    {renderMethod("phone", <MobileOutlined/>, "手机验证", info?.phone)}
                                    {renderMethod("mail", <MailOutlined/>, "邮箱验证", info?.mail)}
                                    {renderMethod("password", <LockOutlined/>, "密码验证", info?.password ? "登录密码" : "")}
                                </div>
                            </div>
                        )}
                        {currentStep === 1 && (
                            <div className={styles.panel}>
                                {!isPasswordVerify && <div className={styles.targetBox}>
                                    <div>
                                        <div className={styles.targetLabel}>验证码已发送至</div>
                                        <div className={styles.targetValue}>{getTypeName(verifyType)}验证</div>
                                    </div>
                                    <div className={styles.targetRight}>
                                        <Typography.Text type="secondary" ellipsis>
                                            {targetValue || "-"}
                                        </Typography.Text>
                                        <div>
                                            <Button
                                                type="link"
                                                className={styles.changeMethod}
                                                onClick={() => {
                                                    clearCountdown();
                                                    form.resetFields(["code"] as any);
                                                    setCurrentStep(0);
                                                }}>
                                                更换方式
                                            </Button>
                                        </div>
                                    </div>
                                </div>}
                                <div className={styles.codeRow}>
                                    <div className={styles.codeInputWrap}>
                                        {isPasswordVerify ? <Form.Item
                                                name="password"
                                                label="登录密码"
                                                rules={[
                                                    {required: true, message: "请输入登录密码"},
                                                    {min: 6, max: 20, message: "密码长度为 6-20 个字符"},
                                                ]}>
                                                <Input
                                                    className={styles.codeInput}
                                                    maxLength={20}
                                                    showCount
                                                    placeholder="请输入当前账号登录密码"
                                                    autoFocus
                                                />
                                            </Form.Item> :
                                            <Form.Item
                                                name="code"
                                                label="验证码"
                                                rules={[
                                                    {required: true, message: "请输入验证码"},
                                                    {pattern: /^\d{6}$/, message: "验证码为 6 位数字"},
                                                ]}>
                                                <Input
                                                    className={styles.codeInput}
                                                    maxLength={6}
                                                    placeholder="请输入 6 位验证码"
                                                    autoFocus
                                                />
                                            </Form.Item>}
                                    </div>
                                    {!isPasswordVerify && <Button
                                        className={styles.resendButton}
                                        loading={resendLoading}
                                        disabled={countdown > 0}
                                        onClick={() => sendCode(false, true)}
                                    >
                                        {countdown > 0 ? `${countdown}s 后重发` : "重新发送"}
                                    </Button>}
                                </div>
                            </div>
                        )}
                    </div>
                </Form>
            </ProModal>
        </>
    );
});

export default Verify;
