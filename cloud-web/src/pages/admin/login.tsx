import {
    UserOutlined,
    LoginOutlined,
    SwapOutlined
} from '@ant-design/icons';
import {App, Button, Checkbox, Form, Input, Tabs, Avatar, Spin} from 'antd';
import React, {useRef, useState, useEffect} from 'react';
import {Body, Icon, useTheme} from 'sinking-antd';
import {useModel} from "umi";
import {login, logout} from "@/service/common/login";
import {setLoginToken, deleteHeader} from "@/utils/auth";
import Captcha, {CaptchaRef} from "@/pages/components/captcha";
import {sendSms} from "@/service/common/sms";
import {sendMail} from "@/service/common/mail";
import Settings from "@/../config/defaultSettings";
import {historyPush} from "@/utils/route";
import {createStyles} from 'antd-style';
import {getPath} from "@/utils/assets";

export default () => {
    const {message, modal} = App.useApp();
    const captcha = useRef<CaptchaRef>({});
    const web = useModel("web");
    const user = useModel("user");
    const theme = useTheme();
    const isCompact = theme?.isCompactTheme() || false;
    const backgroundUrl = getPath("/images/background.svg");

    const useStyles: any = React.useMemo(() =>
            createStyles(({css, responsive, token, isDarkMode}: any) => ({
                body: css`
                    padding: 0;
                `,
                container: css`
                    display: flex;
                    position: absolute;
                    top: 0;
                    left: 0;
                    right: 0;
                    bottom: 0;
                    width: 100%;
                    height: 100%;
                    flex-direction: column;
                    justify-content: flex-start;
                    align-items: center;
                    padding-top: ${isCompact ? 100 : 80}px;
                    margin: 0;
                    background-color: ${isDarkMode ? token.colorBgContainer : '#ffffff'};
                    background-image: url('${backgroundUrl}');
                    background-repeat: no-repeat;
                    background-position: center top;
                    background-size: min(1361px, 100%) auto;
                    overflow: auto;

                    ${responsive.sm} {
                        padding-left: 24px;
                        padding-right: 24px;
                        background-size: 160% auto;
                        background-position: center 0;
                    }
                `,
                rightPanel: css`
                    width: 100%;
                    max-width: ${isCompact ? 320 : 440}px;
                    display: flex;
                    flex-direction: column;
                    align-items: center;

                    ${responsive.md} {
                        max-width: 100%;
                    }
                `,
                loginContent: css`
                    max-width: ${isCompact ? 300 : 360}px;
                    width: 100%;
                    position: relative;
                    z-index: 1;

                    ${responsive.md} {
                        max-width: 360px;
                    }

                    ${responsive.sm} {
                        max-width: 100%;
                    }
                `,
                mobileHeader: css`
                    text-align: center;
                    margin-bottom: ${token.margin}px;

                    .mobile-logo-title {
                        display: flex;
                        align-items: center;
                        justify-content: center;
                        gap: ${token.marginSM}px;
                        margin-bottom: ${token.marginXS}px;
                        height: ${isCompact ? 40 : 48}px;
                    }

                    img {
                        width: ${isCompact ? 40 : 48}px;
                        height: ${isCompact ? 40 : 48}px;
                        flex-shrink: 0;
                    }

                    .mobile-title {
                        font-size: ${isCompact ? 40 : 48}px;
                        font-weight: 700;
                        color: ${isDarkMode ? token.colorText : token.colorTextHeading};
                        line-height: ${isCompact ? 40 : 48}px;
                        height: ${isCompact ? 40 : 48}px;
                        display: flex;
                        align-items: center;
                    }

                    .mobile-subtitle {
                        font-size: ${token.fontSizeLG}px;
                        color: ${isDarkMode ? token.colorTextSecondary : token.colorTextTertiary};
                        line-height: 1.5;
                        margin-top: ${token.margin + 2}px;
                    }

                    ${responsive.sm} {
                        margin-bottom: ${token.marginSM}px;
                        padding-top: ${token.margin}px;

                        .mobile-logo-title {
                            height: 40px;
                            margin-bottom: ${token.marginXXS}px;
                        }

                        .mobile-title {
                            font-size: 40px;
                            line-height: 40px;
                            height: 40px;
                        }

                        img {
                            width: 40px;
                            height: 40px;
                        }
                    }
                `,
                loginTitle: css`
                    font-size: ${isCompact ? token.fontSizeHeading3 : token.fontSizeHeading2}px;
                    font-weight: 700;
                    color: ${isDarkMode ? token.colorText : token.colorTextHeading};
                    margin-bottom: ${token.marginXXS}px;
                    text-align: center;

                    ${responsive.sm} {
                        font-size: ${token.fontSizeHeading3}px;
                    }
                `,
                loginSubtitle: css`
                    font-size: ${isCompact ? token.fontSize : token.fontSizeSM}px;
                    color: ${isDarkMode ? token.colorTextSecondary : token.colorTextTertiary};
                    margin-bottom: ${token.marginSM}px;
                    text-align: center;

                    ${responsive.sm} {
                        font-size: ${token.fontSize}px;
                        margin-bottom: ${token.marginXS}px;
                    }
                `,
                tabsWrapper: css`
                    .ant-tabs-nav {
                        margin-bottom: ${isCompact ? 10 : 0}px !important;
                        height: ${isCompact ? 48 : 60}px !important;

                        &::before {
                            border: none !important;
                        }
                    }

                    .ant-tabs-tab {
                        padding: 0;
                        height: ${isCompact ? 48 : 60}px !important;
                        display: flex;
                        align-items: center;
                        justify-content: center;

                        &:hover {
                            color: ${token.colorPrimaryHover};
                        }
                    }

                    .ant-tabs-tab-btn {
                        font-size: ${token.fontSizeHeading4}px;
                        font-weight: 600;
                        color: ${isDarkMode ? token.colorTextTertiary : '#979797'};
                        line-height: 1;
                        transform: scale(0.8);
                        transform-origin: center;
                        transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1),
                        color 0.3s ease,
                        font-weight 0.3s ease;
                        will-change: transform;
                        -webkit-font-smoothing: antialiased;
                        -moz-osx-font-smoothing: grayscale;
                        text-rendering: optimizeLegibility;
                    }

                    .ant-tabs-tab-active .ant-tabs-tab-btn {
                        font-weight: 700;
                        color: ${isDarkMode ? token.colorText : token.colorTextHeading};
                        transform: scale(1);
                    }

                    .ant-tabs-nav-wrap {
                        &::after {
                            content: '';
                            position: absolute;
                            bottom: 0;
                            left: 0;
                            right: 0;
                            height: 2px;
                            background: ${isDarkMode ? token.colorBorderSecondary : '#f0f0f0'};
                        }
                    }

                    .ant-tabs-ink-bar {
                        position: absolute;
                        bottom: 0;
                        height: 4px !important;
                        width: 55px !important;
                        border-radius: ${token.borderRadiusSM}px !important;
                        background: ${token.colorPrimary} !important;
                        z-index: 1;
                        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
                    }

                    .ant-tabs-tabpane {
                        margin-top: ${token.margin}px;
                    }
                `,
                formWrapper: css`
                    .ant-form-item {
                        margin-bottom: ${isCompact ? 20 : 25}px;
                    }

                    .ant-input-affix-wrapper,
                    .ant-input {
                        height: ${isCompact ? 40 : 45}px;
                        font-size: ${token.fontSize}px;
                        background: ${isDarkMode ? token.colorBgElevated : 'rgb(255, 255, 255)'};
                        border: ${isDarkMode ? `1px solid ${token.colorBorder}` : 'none'};
                        border-radius: ${token.borderRadiusLG}px;
                        box-shadow: ${isDarkMode ? `0 2px 8px rgba(0,0,0,0.15)` : 'rgba(0, 0, 0, 0.05) 0 7px 25px 0'};
                        padding: 0 ${token.paddingLG}px;
                        transition: box-shadow 0.3s ease;

                        &:hover {
                            box-shadow: ${isDarkMode ? `0 4px 12px rgba(0,0,0,0.25)` : 'rgba(0, 0, 0, 0.08) 0 10px 30px 0'};
                        }

                        &:focus-within {
                            box-shadow: ${isDarkMode
                                    ? `0 4px 12px rgba(0,0,0,0.3), 0 0 0 2px ${token.colorPrimaryBg}`
                                    : `rgba(0, 0, 0, 0.08) 0 10px 30px 0, 0 0 0 2px ${token.colorPrimaryBg}`};
                        }

                        &.ant-input-status-error,
                        &.ant-input-affix-wrapper-status-error {
                            border-color: ${token.colorError} !important;
                            box-shadow: ${isDarkMode
                                    ? `0 2px 8px rgba(0,0,0,0.15)`
                                    : 'rgba(0, 0, 0, 0.05) 0 7px 25px 0'} !important;

                            .ant-input {
                                border-color: transparent !important;
                            }
                        }

                        .ant-input {
                            background: transparent;
                            border: none;
                            box-shadow: none;
                            height: auto;
                            padding: 0;
                        }

                        input::placeholder {
                            color: ${isDarkMode ? token.colorTextTertiary : 'darkgray'};
                        }

                        .anticon {
                            color: ${isDarkMode ? token.colorTextTertiary : 'rgba(0,0,0,0.45)'};
                            font-size: ${token.fontSize}px;
                        }
                    }

                    .ant-btn {
                        border-radius: ${token.borderRadiusLG}px;
                        font-weight: 500;
                    }

                    .ant-btn-primary {
                        height: ${isCompact ? 40 : 45}px;
                        font-size: ${token.fontSizeLG}px;
                        border-radius: ${token.borderRadiusLG}px;
                        background: ${token.colorPrimary};
                        border-color: ${token.colorPrimary};
                        box-shadow: ${isDarkMode ? `0 4px 12px ${token.colorPrimary}40` : `${token.colorPrimary}80 0 6px 27px 0`};
                        font-weight: 600;

                        &:hover {
                            background: ${token.colorPrimaryHover} !important;
                            border-color: ${token.colorPrimaryHover} !important;
                            box-shadow: ${isDarkMode ? `0 6px 16px ${token.colorPrimaryHover}40` : `${token.colorPrimaryHover}80 0 8px 32px 0`};
                        }
                    }

                    .ant-btn-default {
                        height: ${isCompact ? 40 : 52}px;
                        border-radius: ${token.borderRadiusLG}px;
                        background: ${isDarkMode ? token.colorBgElevated : '#ffffff'} !important;
                        border: 1px solid ${isDarkMode ? token.colorBorder : 'transparent'} !important;
                        box-shadow: ${isDarkMode ? `0 2px 8px rgba(0,0,0,0.15)` : 'rgba(0, 0, 0, 0.05) 0 4px 12px 0'} !important;

                        &:hover,
                        &:focus,
                        &:focus-visible {
                            background: ${isDarkMode ? token.colorBgTextHover : token.colorBgTextHover} !important;
                            border-color: ${isDarkMode ? token.colorBorder : 'transparent'} !important;
                            box-shadow: ${isDarkMode ? `0 4px 12px rgba(0,0,0,0.25)` : 'rgba(0, 0, 0, 0.08) 0 6px 16px 0'} !important;
                        }
                    }
                `,
                extraRow: css`
                    display: flex;
                    justify-content: space-between;
                    align-items: center;
                    margin-bottom: ${isCompact ? token.marginLG : token.marginXL}px;

                    .ant-checkbox-wrapper {
                        font-size: ${token.fontSizeSM}px;
                        color: ${isDarkMode ? token.colorTextSecondary : 'rgba(0, 0, 0, 0.5)'};
                    }

                    a {
                        font-size: ${token.fontSizeSM}px;
                        color: ${isDarkMode ? token.colorTextSecondary : 'rgba(0, 0, 0, 0.45)'};
                        transition: all 0.3s ease;

                        &:hover {
                            color: ${token.colorPrimary};
                        }
                    }
                `,
                loggedInCard: css`
                    text-align: center;
                    padding: ${token.paddingXL}px;
                    background: ${isDarkMode ? token.colorBgElevated : 'rgb(255, 255, 255)'};
                    border-radius: ${token.borderRadiusLG}px;
                    border: ${isDarkMode ? `1px solid ${token.colorBorder}` : 'none'};
                    margin-bottom: ${token.marginXL}px;
                    box-shadow: ${isDarkMode ? '0 2px 8px rgba(0,0,0,0.2)' : 'rgba(0, 0, 0, 0.08) 0 7px 25px 0'};

                    .user-avatar {
                        margin-bottom: ${token.margin}px;
                        border: 3px solid ${token.colorPrimary};
                        box-shadow: 0 4px 12px ${token.colorPrimary}40;
                    }

                    .user-name {
                        font-size: ${token.fontSizeLG}px;
                        font-weight: 600;
                        color: ${isDarkMode ? token.colorText : token.colorTextHeading};
                        margin-bottom: ${token.marginXS}px;
                    }

                    .user-contact {
                        font-size: ${token.fontSizeSM}px;
                        color: ${isDarkMode ? token.colorTextSecondary : token.colorTextTertiary};
                        margin-bottom: ${token.marginLG}px;
                    }
                `,
                quickActions: css`
                    display: flex;
                    flex-direction: column;
                    gap: ${isCompact ? token.marginXS : token.marginSM}px;

                    .ant-btn {
                        height: ${isCompact ? 40 : 50}px;
                        font-size: ${token.fontSize}px;
                        border-radius: ${token.borderRadiusLG}px;
                        font-weight: 600;
                    }

                    .ant-btn-primary {
                        background: ${token.colorPrimary};
                        border-color: ${token.colorPrimary};
                        box-shadow: ${isDarkMode ? `0 4px 12px ${token.colorPrimary}40` : `${token.colorPrimary}80 0 6px 27px 0`};

                        &:hover {
                            background: ${token.colorPrimaryHover} !important;
                            border-color: ${token.colorPrimaryHover} !important;
                            box-shadow: ${isDarkMode ? `0 6px 16px ${token.colorPrimaryHover}40` : `${token.colorPrimaryHover}80 0 8px 32px 0`};
                        }
                    }

                    .ant-btn-default {
                        background: ${isDarkMode ? token.colorBgElevated : '#ffffff'};
                        border: ${isDarkMode ? `1px solid ${token.colorBorder}` : '#ffffff'};
                        box-shadow: ${isDarkMode ? `0 2px 8px rgba(0,0,0,0.15)` : 'rgba(0, 0, 0, 0.05) 0 4px 12px 0'};

                        &:hover {
                            background: ${isDarkMode ? token.colorBgTextHover : token.colorBgTextHover};
                            box-shadow: ${isDarkMode ? `0 4px 12px rgba(0,0,0,0.25)` : 'rgba(0, 0, 0, 0.08) 0 6px 16px 0'};
                        }
                    }
                `,
                reloginLink: css`
                    text-align: center;
                    margin-top: ${token.marginLG}px;

                    a {
                        font-size: ${token.fontSizeSM}px;
                        color: ${token.colorPrimary};

                        &:hover {
                            color: ${token.colorPrimaryHover};
                        }
                    }
                `,
                sendSmsBtn: css`
                    font-size: ${token.fontSizeSM}px;
                    color: ${token.colorPrimary};
                    cursor: pointer;
                    user-select: none;

                    &:hover {
                        color: ${token.colorPrimaryHover};
                    }
                `,
                sendSmsDisabled: css`
                    font-size: ${token.fontSizeSM}px;
                    color: ${isDarkMode ? token.colorTextTertiary : 'gray'};
                    cursor: default;
                `,
            })),
        [isCompact]
    );

    const {styles} = useStyles();
    const [form] = Form.useForm();

    const [type, setType] = useState<string>('account');
    const [isLoading, setIsLoading] = useState(false);
    const [isCodeSendLoading, setIsCodeSendLoading] = useState(false);
    const [countdown, setCountdown] = useState(0);
    const [showRelogin, setShowRelogin] = useState(false);
    const [checkingLogin, setCheckingLogin] = useState(true);

    /**
     * 初始化时检查登录状态
     */
    useEffect(() => {
        setCheckingLogin(true);
        user?.refreshSystemUser?.((d: any) => {
            setCheckingLogin(false);
            if (!d) {
                setShowRelogin(true);
            }
        });
    }, []);

    /**
     * 验证账号（手机号或邮箱）
     * @param value 账号
     */
    const validateAccount = (value: string): { valid: boolean; type: 'phone' | 'email' | 'unknown' } => {
        const phoneRegex = /^1[3-9]\d{9}$/;
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (phoneRegex.test(value)) {
            return {valid: true, type: 'phone'};
        }
        if (emailRegex.test(value)) {
            return {valid: true, type: 'email'};
        }
        return {valid: false, type: 'unknown'};
    };

    /**
     * 获取验证码倒计时
     */
    const startCountdown = () => {
        let timeLeft = 60;
        setCountdown(timeLeft);
        const timer = setInterval(() => {
            timeLeft--;
            if (timeLeft <= 0) {
                clearInterval(timer);
                setCountdown(0);
            } else {
                setCountdown(timeLeft);
            }
        }, 1000);
    };

    // 统一登录处理
    const handleLogin = async (values: any, captchaRes: any, loginType: string) => {
        const body: any = {
            account: values.account,
            is_admin: true,
            token: captchaRes?.token,
            captcha_x: captchaRes?.x,
            captcha_y: captchaRes?.y
        };

        // 根据登录类型添加对应字段
        if (loginType === "account") {
            body.password = values.password;
        } else {
            body.code = values.code;
        }

        await login({
            body,
            onSuccess: (r: any) => {
                setLoginToken(r?.data);
                user?.refreshSystemUser?.(() => {
                    message?.success(r?.message);
                    historyPush("admin.index");
                    setIsLoading(false);
                });
            },
            onFail: (r: any) => {
                if (loginType === "code") {
                    modal.error({
                        title: '登录失败',
                        content: r?.message,
                        okText: "确认"
                    });
                } else {
                    message?.error(r?.message || "请求失败");
                }
                setIsLoading(false);
            }
        });
    };

    /**
     * 发送验证码（自动识别手机号或邮箱）
     */
    const handleSendCode = () => {
        const account = form.getFieldValue("account" as any);
        const result = validateAccount(account);
        if (!result.valid) {
            message?.error("请输入正确的手机号或邮箱");
            return;
        }
        if (countdown > 0) return;

        const sendFn = result.type === 'email' ? sendMail : sendSms;
        const targetValue = result.type === 'email' ? {mail: account} : {phone: account};

        captcha?.current?.Show?.(async (res) => {
            setIsCodeSendLoading(true);
            await sendFn({
                body: {
                    token: res?.token,
                    captcha_x: res?.x,
                    captcha_y: res?.y,
                    is_admin: true,
                    ...targetValue
                },
                onSuccess: (r) => {
                    message?.success(r?.message);
                    startCountdown();
                },
                onFail: (r) => message?.error(r?.message),
                onFinally: () => setIsCodeSendLoading(false)
            });
        }, () => {
            setIsCodeSendLoading(false);
        });
    };

    /**
     * 提交登陆表单
     * @param values 表单数据
     */
    const handleSubmit = async (values: any) => {
        // 验证输入
        if (type === "account") {
            if (!values.account || values.account.length === 0) {
                message?.error("请输入账号或手机号");
                return;
            }
            if (!values.password || values.password.length <= 5) {
                message?.error("请输入正确的密码");
                return;
            }
        } else {
            const result = validateAccount(values.account);
            if (!result.valid) {
                message?.error("请输入正确的手机号或邮箱");
                return;
            }
            if (!values.code || values.code.length !== 6) {
                message?.error("请输入正确的验证码");
                return;
            }
        }

        // 显示验证码并执行登录
        captcha?.current?.Show?.(async (res) => {
            setIsLoading(true);
            await handleLogin(values, res, type);
        });
    };

    /**
     * 退出登录
     */
    const outLogin = async () => {
        message.loading({content: "正在退出登录...", key: "outLogin"});
        await logout({
            body: {is_admin: true},
            onSuccess: (r) => {
                message.success(r?.message || "退出登录成功");
                deleteHeader();
                user?.setSystem(undefined);
                message.destroy("outLogin");
                setShowRelogin(true);
            },
            onFail: (r) => {
                message.error(r?.message || "退出登录失败");
                message.destroy("outLogin");
            }
        });
    }

    const items: any = [
        {key: 'account', label: '密码登录'},
        {key: 'code', label: '验证登录'}
    ];

    return (
        <Body space={false} animation={false} className={styles.body}>
            <div className={styles.container}>
                <Captcha ref={captcha}/>
                <div className={styles.rightPanel}>
                    <div className={styles.mobileHeader}>
                        <div className="mobile-logo-title">
                            {web?.info?.logo ? (
                                <img src={web?.info?.logo || ((Settings?.basePath || '/') + 'logo.svg')} alt="logo"/>
                            ) : (
                                <Icon type="icon-logo" style={{
                                    fontSize: '48px',
                                    color: theme?.themes?.token?.colorPrimary || '#1890ff'
                                }}/>
                            )}
                            <div className="mobile-title">{web?.info?.name || Settings?.title}</div>
                        </div>
                        <div className="mobile-subtitle">后台管理平台</div>
                    </div>
                    <div className={styles.loginContent}>
                        {checkingLogin ? (
                            <div style={{textAlign: 'center', padding: '80px 0'}}>
                                <Spin size="large"/>
                            </div>
                        ) : user?.system?.id && !showRelogin ? (
                            <>
                                <div className={styles.loginTitle}>欢迎回来</div>
                                <div className={styles.loginSubtitle}>检测到您已登录</div>

                                <div className={styles.loggedInCard}>
                                    <Avatar
                                        size={80}
                                        src={user?.system?.avatar || false}
                                        icon={<UserOutlined/>}
                                        className="user-avatar"
                                    />
                                    <div className="user-name">
                                        {user?.system?.nick_name || '管理员' + (user?.system?.phone || user?.system?.email || user?.system?.account)?.toString().slice(0, 4)}
                                    </div>
                                    <div className="user-contact">
                                        {user?.system?.phone || user?.system?.email || user?.system?.account || '未绑定联系方式'}
                                    </div>
                                </div>
                                <div className={styles.quickActions}>
                                    <Button type="primary" icon={<LoginOutlined/>}
                                            onClick={() => historyPush("admin.index")} block>
                                        快捷登录
                                    </Button>
                                    <Button icon={<SwapOutlined/>} onClick={() => setShowRelogin(true)} block>
                                        切换账号
                                    </Button>
                                </div>
                                <div className={styles.reloginLink}>
                                    <a onClick={outLogin}>
                                        退出当前账号
                                    </a>
                                </div>
                            </>
                        ) : (
                            <>
                                <Tabs items={items} activeKey={type} onChange={setType} size="large" centered
                                      className={styles.tabsWrapper}/>
                                <Form form={form} onFinish={handleSubmit} className={styles.formWrapper}>
                                    {type === 'account' && (
                                        <>
                                            <Form.Item name='account'
                                                       rules={[{required: true, message: '请输入账号或手机号或邮箱'}]}>
                                                <Input
                                                    placeholder='请输入账号或手机号或邮箱'
                                                    allowClear
                                                />
                                            </Form.Item>
                                            <Form.Item name='password'
                                                       rules={[{required: true, message: '请输入密码'}]}>
                                                <Input.Password
                                                    placeholder='请输入密码'
                                                />
                                            </Form.Item>
                                            <div className={styles.extraRow}>
                                                <Checkbox defaultChecked>下次自动登录</Checkbox>
                                                <a onClick={() => {
                                                    message?.warning("请使用验证登录后修改密码");
                                                    setType("code");
                                                }}>忘记密码</a>
                                            </div>
                                        </>
                                    )}
                                    {type === 'code' && (
                                        <>
                                            <Form.Item name='account'
                                                       rules={[{required: true, message: '请输入手机号或邮箱'}]}>
                                                <Input
                                                    placeholder='请输入手机号或邮箱'
                                                    allowClear
                                                />
                                            </Form.Item>
                                            <Form.Item name='code' rules={[{required: true, message: '请输入验证码'}]}>
                                                <Input
                                                    placeholder='请输入验证码'
                                                    allowClear
                                                    suffix={
                                                        countdown > 0 ? (
                                                            <span
                                                                className={styles.sendSmsDisabled}>{countdown}s后重新发送</span>
                                                        ) : (
                                                            <span
                                                                className={styles.sendSmsBtn}
                                                                onClick={() => {
                                                                    const account = form.getFieldValue('account' as any);
                                                                    if (!account) {
                                                                        message?.error('请先输入手机号或邮箱');
                                                                        return;
                                                                    }
                                                                    handleSendCode();
                                                                }}
                                                            >
                                                        {isCodeSendLoading ? '发送中...' : '发送验证码'}
                                                    </span>
                                                        )
                                                    }
                                                />
                                            </Form.Item>
                                            <div className={styles.extraRow}>
                                                <Checkbox defaultChecked>下次自动登录</Checkbox>
                                            </div>
                                        </>
                                    )}
                                    <Form.Item>
                                        <Button
                                            type='primary'
                                            loading={isLoading}
                                            htmlType='submit'
                                            block>
                                            登 录
                                        </Button>
                                    </Form.Item>
                                </Form>
                            </>
                        )}
                    </div>
                </div>
            </div>
        </Body>
    );
};
