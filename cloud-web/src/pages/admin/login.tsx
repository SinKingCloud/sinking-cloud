import {
    UserOutlined,
    SafetyCertificateOutlined,
    ThunderboltOutlined,
    CustomerServiceOutlined,
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
import Settings from "@/../config/defaultSettings";
import {historyPush} from "@/utils/route";
import {createStyles} from 'antd-style';

const useStyles: any = createStyles(({css, responsive, token, isDarkMode}: any): any => {
    return {
        container: css`
            display: flex;
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            width: 100%;
            height: 100%;
            margin: 0;
            padding: 0;
            background: ${isDarkMode ? token.colorBgContainer : '#ffffff'};
            animation: fadeIn 0.6s ease-out;

            @keyframes fadeIn {
                from {
                    opacity: 0;
                }
                to {
                    opacity: 1;
                }
            }

            ${responsive.md} {
                flex-direction: column;
                position: fixed;
                width: 100vw;
                height: 100vh;
                overflow-y: auto;
                overflow-x: hidden;
                top: 0;
                left: 0;
            }
        `,
        leftPanel: css`
            flex: 1;
            background: linear-gradient(135deg, ${token.colorPrimary} 0%, ${token.colorPrimaryHover} 100%);
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            padding: ${token.paddingXL * 2}px;
            position: relative;
            overflow: hidden;
            animation: slideInLeft 0.8s ease-out;

            @keyframes slideInLeft {
                from {
                    opacity: 0;
                    transform: translateX(-50px);
                }
                to {
                    opacity: 1;
                    transform: translateX(0);
                }
            }

            @media (prefers-reduced-motion: reduce) {
                animation: none;
            }

            &::before {
                content: '';
                position: absolute;
                top: 0;
                left: 0;
                right: 0;
                bottom: 0;
                background-image: radial-gradient(circle at 20% 30%, rgba(255, 255, 255, 0.2) 0%, transparent 50%),
                radial-gradient(circle at 80% 70%, rgba(255, 255, 255, 0.15) 0%, transparent 50%),
                radial-gradient(circle at 40% 80%, rgba(255, 255, 255, 0.1) 0%, transparent 50%);
                z-index: 1;
                animation: breathe 4s ease-in-out infinite;
            }

            &::after {
                content: '';
                position: absolute;
                top: 0;
                left: 0;
                right: 0;
                bottom: 0;
                background-image: linear-gradient(45deg, transparent 48%, rgba(255, 255, 255, 0.03) 49%, rgba(255, 255, 255, 0.03) 51%, transparent 52%),
                linear-gradient(-45deg, transparent 48%, rgba(255, 255, 255, 0.03) 49%, rgba(255, 255, 255, 0.03) 51%, transparent 52%);
                background-size: 60px 60px;
                z-index: 0;
                animation: slide 5s linear infinite;
            }

            @keyframes breathe {
                0%, 100% {
                    opacity: 1;
                    transform: scale(1);
                }
                50% {
                    opacity: 0.8;
                    transform: scale(1.05);
                }
            }

            @keyframes slide {
                0% {
                    background-position: 0 0;
                }
                100% {
                    background-position: 60px 60px;
                }
            }

            ${responsive.md} {
                display: none;
            }
        `,
        brandContent: css`
            position: relative;
            z-index: 2;
            text-align: center;
            color: white;
            max-width: 500px;
            width: 100%;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
        `,
        brandLogo: css`
            font-size: ${token.fontSizeHeading1 * 2.5}px;
            color: white;
            margin-bottom: ${token.marginXS}px;
            transition: all 0.4s ease;

            &:hover {
                transform: scale(1.1) rotate(5deg);
            }

            img {
                width: 80px;
                height: 80px;
                border-radius: ${token.borderRadiusLG}px;
                transition: all 0.4s ease;
            }

            ${responsive.md} {
                font-size: ${token.fontSizeHeading1 * 2}px;
                margin-bottom: ${token.marginXS}px;

                img {
                    width: 64px;
                    height: 64px;
                }
            }

            ${responsive.sm} {
                font-size: ${token.fontSizeHeading1 * 1.5}px;
                margin-bottom: ${token.marginXXS}px;

                img {
                    width: 56px;
                    height: 56px;
                }
            }
        `,
        brandTitle: css`
            font-size: ${token.fontSizeHeading1 * 1.5}px;
            font-weight: 700;
            margin-bottom: ${token.marginXS}px;
            color: white;
            letter-spacing: 1px;

            ${responsive.md} {
                font-size: ${token.fontSizeHeading2}px;
                margin-bottom: ${token.marginXXS}px;
            }

            ${responsive.sm} {
                font-size: ${token.fontSizeHeading3}px;
                margin-bottom: ${token.marginXXS}px;
            }
        `,
        brandSubtitle: css`
            font-size: ${token.fontSize}px;
            color: rgba(255, 255, 255, 0.9);
            margin-bottom: ${token.marginLG}px;
            line-height: ${token.lineHeight};

            ${responsive.md} {
                font-size: ${token.fontSizeSM}px;
                margin-bottom: ${token.margin}px;
            }

            ${responsive.sm} {
                font-size: ${token.fontSizeSM}px;
                margin-bottom: ${token.marginSM}px;
            }
        `,
        brandFeatures: css`
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: ${token.marginSM}px;
            margin-top: ${token.marginLG}px;

            ${responsive.md} {
                grid-template-columns: 1fr 1fr;
                gap: ${token.marginSM}px;
                margin-top: ${token.margin}px;
            }

            ${responsive.sm} {
                grid-template-columns: 1fr;
                gap: ${token.marginXS}px;
                margin-top: ${token.marginSM}px;
            }
        `,
        featureItem: css`
            display: flex;
            align-items: center;
            gap: ${token.marginXS}px;
            text-align: left;
            padding: ${token.paddingXS}px ${token.paddingSM}px;
            border-radius: ${token.borderRadius}px;
            background: rgba(255, 255, 255, 0.05);
            backdrop-filter: blur(10px);
            transition: all 0.3s ease;
            cursor: pointer;

            &:hover {
                background: rgba(255, 255, 255, 0.15);
                transform: translateX(8px);
                box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);

                .feature-icon {
                    transform: scale(1.2);
                }
            }

            .feature-icon {
                font-size: ${token.fontSizeLG}px;
                color: white;
                opacity: 0.9;
                transition: all 0.3s ease;
            }

            .feature-text {
                font-size: ${token.fontSizeSM}px;
                color: rgba(255, 255, 255, 0.85);
                line-height: ${token.lineHeight};
            }

            ${responsive.md} {
                .feature-icon {
                    font-size: ${token.fontSize}px;
                }

                .feature-text {
                    font-size: ${token.fontSizeSM}px;
                }
            }
        `,
        rightPanel: css`
            width: 520px;
            background: ${isDarkMode ? token.colorBgContainer : '#ffffff'};
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            padding: ${token.paddingXL * 2.5}px;
            box-shadow: ${isDarkMode ? '0 0 40px rgba(0,0,0,0.5)' : '-2px 0 40px rgba(0,0,0,0.08)'};
            animation: slideInRight 0.8s ease-out;

            @keyframes slideInRight {
                from {
                    opacity: 0;
                    transform: translateX(50px);
                }
                to {
                    opacity: 1;
                    transform: translateX(0);
                }
            }

            @media (prefers-reduced-motion: reduce) {
                animation: none;
            }

            ${responsive.lg} {
                width: 460px;
                padding: ${token.paddingXL * 2}px;
            }

            ${responsive.md} {
                width: 100%;
                min-height: 100vh;
                padding: ${token.paddingXL * 2}px ${token.paddingXL}px;
                box-shadow: none;
                justify-content: flex-start;
                background: ${isDarkMode ? token.colorBgContainer : '#ffffff'};
                animation: none;
                position: relative;
            }

            ${responsive.sm} {
                padding: ${token.paddingXL}px ${token.paddingLG}px;
                min-height: 100vh;
            }
        `,
        loginContent: css`
            max-width: 420px;
            width: 100%;
            position: relative;
            z-index: 1;

            ${responsive.md} {
                max-width: 450px;
            }

            ${responsive.sm} {
                max-width: 100%;
            }
        `,
        mobileHeader: css`
            display: none;

            ${responsive.md} {
                display: block;
                text-align: center;
                margin-top: ${token.marginXL}px;
                margin-bottom: ${token.marginXXL}px;
                position: relative;
                z-index: 1;

                .mobile-logo-title {
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    gap: ${token.marginSM}px;
                    margin-bottom: ${token.marginLG}px;
                    height: 48px;
                }

                img {
                    width: 48px;
                    height: 48px;
                    flex-shrink: 0;
                }

                .mobile-title {
                    font-size: 48px;
                    font-weight: 700;
                    color: ${isDarkMode ? token.colorText : token.colorTextHeading};
                    line-height: 48px;
                    height: 48px;
                    display: flex;
                    align-items: center;
                }

                .mobile-subtitle {
                    font-size: ${token.fontSizeSM}px;
                    color: ${isDarkMode ? token.colorTextSecondary : token.colorTextTertiary};
                    line-height: 1.5;
                }
            }

            ${responsive.sm} {
                margin-top: ${token.marginLG}px;
                margin-bottom: ${token.marginXL}px;

                .mobile-logo-title {
                    gap: ${token.marginXS}px;
                    margin-bottom: ${token.margin}px;
                    height: 40px;
                }

                img {
                    width: 40px;
                    height: 40px;
                }

                .mobile-title {
                    font-size: 40px;
                    line-height: 40px;
                    height: 40px;
                }
            }
        `,
        loginTitle: css`
            font-size: ${token.fontSizeHeading2}px;
            font-weight: 700;
            color: ${isDarkMode ? token.colorText : token.colorTextHeading};
            margin-bottom: ${token.marginXS}px;

            ${responsive.md} {
                display: none;
            }
        `,
        loginSubtitle: css`
            font-size: ${token.fontSizeSM}px;
            color: ${isDarkMode ? token.colorTextSecondary : token.colorTextTertiary};
            margin-bottom: ${token.marginXL}px;

            ${responsive.md} {
                display: none;
            }
        `,
        tabsWrapper: css`
            margin-bottom: ${token.marginLG}px;

            .ant-tabs-nav {
                margin-bottom: 0 !important;
                height: 60px !important;

                &::before {
                    border: none !important;
                }
            }

            .ant-tabs-tab {
                padding: 0;
                height: 60px !important;
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
                margin-bottom: ${token.marginXL}px;
            }

            .ant-input-affix-wrapper,
            .ant-input {
                height: 52px;
                font-size: ${token.fontSize}px;
                background: ${isDarkMode ? token.colorBgElevated : 'rgb(255, 255, 255)'};
                border: ${isDarkMode ? `1px solid ${token.colorBorder}` : 'none'};
                border-radius: ${token.borderRadiusLG}px;
                box-shadow: ${isDarkMode ? `0 2px 8px rgba(0,0,0,0.15)` : 'rgba(0, 0, 0, 0.05) 0 7px 25px 0'};
                padding: 0 ${token.paddingLG}px;
                transition: all 0.3s ease;

                &:hover {
                    box-shadow: ${isDarkMode ? `0 4px 12px rgba(0,0,0,0.25)` : 'rgba(0, 0, 0, 0.08) 0 10px 30px 0'};
                }

                &:focus-within {
                    box-shadow: ${isDarkMode
                            ? `0 4px 12px rgba(0,0,0,0.3), 0 0 0 2px ${token.colorPrimaryBg}`
                            : `rgba(0, 0, 0, 0.08) 0 10px 30px 0, 0 0 0 2px ${token.colorPrimaryBg}`};
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
                height: 50px;
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
                height: 52px;
                border-radius: ${token.borderRadiusLG}px;
                background: ${isDarkMode ? token.colorBgElevated : '#ffffff'};
                border: ${isDarkMode ? `1px solid ${token.colorBorder}` : 'none'};
                box-shadow: ${isDarkMode ? `0 2px 8px rgba(0,0,0,0.15)` : 'rgba(0, 0, 0, 0.05) 0 4px 12px 0'};

                &:hover {
                    background: ${isDarkMode ? token.colorBgTextHover : token.colorBgTextHover};
                    box-shadow: ${isDarkMode ? `0 4px 12px rgba(0,0,0,0.25)` : 'rgba(0, 0, 0, 0.08) 0 6px 16px 0'};
                }
            }
        `,
        phonePrefix: css`
            color: ${isDarkMode ? token.colorTextTertiary : 'rgba(0, 0, 0, 0.35)'};
            margin-right: ${token.marginXS}px;
        `,
        divider: css`
            color: ${isDarkMode ? token.colorTextTertiary : 'rgba(0, 0, 0, 0.35)'};
            margin: 0 ${token.marginXS}px;
        `,
        extraRow: css`
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: ${token.marginLG}px;

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
        agreement: css`
            .ant-checkbox-wrapper {
                font-size: ${token.fontSizeSM}px;
                color: ${isDarkMode ? token.colorTextSecondary : token.colorTextTertiary};

                a {
                    color: ${token.colorPrimary};

                    &:hover {
                        color: ${token.colorPrimaryHover};
                    }
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

            .user-phone {
                font-size: ${token.fontSizeSM}px;
                color: ${isDarkMode ? token.colorTextSecondary : token.colorTextTertiary};
                margin-bottom: ${token.marginLG}px;
            }
        `,
        quickActions: css`
            display: flex;
            flex-direction: column;
            gap: ${token.marginSM}px;

            .ant-btn {
                height: 50px;
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
                border: ${isDarkMode ? `1px solid ${token.colorBorder}` : 'none'};
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
    };
});

export default () => {
    const {styles} = useStyles();
    const {message, modal} = App.useApp();
    const captcha = useRef<CaptchaRef>({});
    const web = useModel("web");
    const user = useModel("user");
    const themeModel = useTheme();
    const [form] = Form.useForm();

    const [type, setType] = useState<string>('account');
    const [isLoading, setIsLoading] = useState(false);
    const [isSmsSendLoading, setIsSmsSendLoading] = useState(false);
    const [isRead, setIsRead] = useState(true);
    const [countdown, setCountdown] = useState(0);
    const [showRelogin, setShowRelogin] = useState(false);
    const [checkingLogin, setCheckingLogin] = useState(true);

    // 初始化时检查登录状态
    useEffect(() => {
        setCheckingLogin(true);
        user?.refreshWebUser?.((d: any) => {
            setCheckingLogin(false);
            if (!d) {
                setShowRelogin(true);
            }
        });
    }, []);

    // 验证手机号
    const validatePhone = (phone: string): boolean => {
        if (!/^1[3-9]\d{9}$/.test(phone)) {
            message?.error("请输入正确的手机号");
            return false;
        }
        return true;
    };

    // 获取验证码倒计时
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
            phone: values.phone,
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
                user?.refreshWebUser(() => {
                    message?.success(r?.message);
                    historyPush("user.index");
                    setIsLoading(false);
                });
            },
            onFail: (r: any) => {
                if (loginType === "phone") {
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

    // 发送短信验证码
    const handleSendSms = () => {
        const phone = form.getFieldValue("phone" as any);
        if (!validatePhone(phone)) return;
        if (countdown > 0) return;

        captcha?.current?.Show?.(async (res) => {
            setIsSmsSendLoading(true);
            await sendSms({
                body: {
                    token: res?.token,
                    captcha_x: res?.x,
                    captcha_y: res?.y,
                    phone: phone
                },
                onSuccess: (r) => {
                    message?.success(r?.message);
                    startCountdown();
                },
                onFail: (r) => message?.error(r?.message),
                onFinally: () => setIsSmsSendLoading(false)
            });
        }, () => {
            setIsSmsSendLoading(false);
        });
    };

    // 表单提交
    const handleSubmit = async (values: any) => {
        if (!isRead) {
            message?.error("请阅读并同意《用户使用条款》协议");
            return;
        }
        // 验证输入
        if (!validatePhone(values.phone)) return;

        // 验证密码或验证码
        if (type === "account" && (!values.password || values.password.length <= 5)) {
            message?.error("请输入正确的密码");
            return;
        }
        if (type === "phone" && (!values.code || values.code.length !== 6)) {
            message?.error("请输入正确的手机验证码");
            return;
        }

        // 显示验证码并执行登录
        captcha?.current?.Show?.(async (res) => {
            setIsLoading(true);
            await handleLogin(values, res, type);
        });
    };

    const items: any = [
        {key: 'account', label: '密码登录'},
        {key: 'phone', label: '短信登录'}
    ];

    return (
        <Body
            space={false}
            animation={false}
            className=""
            style={{
                padding: 0,
                margin: 0,
                height: '100%',
                overflow: 'hidden',
                background: 'transparent'
            }}>
            <div className={styles.container}>
                <Captcha ref={captcha}/>

                {/* 左侧品牌区域 */}
                <div className={styles.leftPanel}>
                    <div className={styles.brandContent}>
                        <div className={styles.brandLogo}>
                            {web?.info?.logo ? (
                                <img
                                    src={web?.info?.logo || ((Settings?.basePath || '/') + 'logo.svg')}
                                    alt="logo"
                                />
                            ) : (
                                <Icon type="icon-logo"/>
                            )}
                        </div>

                        <div className={styles.brandTitle}>
                            {web?.info?.name || Settings?.title}
                        </div>

                        <div className={styles.brandSubtitle}>
                            为您提供专业、安全、高效的一站式网络服务平台
                        </div>

                        <div className={styles.brandFeatures}>
                            <div className={styles.featureItem}>
                                <SafetyCertificateOutlined className="feature-icon"/>
                                <span className="feature-text">安全可靠</span>
                            </div>
                            <div className={styles.featureItem}>
                                <ThunderboltOutlined className="feature-icon"/>
                                <span className="feature-text">快速响应</span>
                            </div>
                            <div className={styles.featureItem}>
                                <UserOutlined className="feature-icon"/>
                                <span className="feature-text">专业团队</span>
                            </div>
                            <div className={styles.featureItem}>
                                <CustomerServiceOutlined className="feature-icon"/>
                                <span className="feature-text">7×24服务</span>
                            </div>
                        </div>
                    </div>
                </div>

                {/* 右侧登录区域 */}
                <div className={styles.rightPanel}>
                    {/* 移动端头部 */}
                    <div className={styles.mobileHeader}>
                        <div className="mobile-logo-title">
                            {web?.info?.logo ? (
                                <img
                                    src={web?.info?.logo || ((Settings?.basePath || '/') + 'logo.svg')}
                                    alt="logo"
                                />
                            ) : (
                                <Icon type="icon-logo" style={{
                                    fontSize: '48px',
                                    color: themeModel?.themes?.token?.colorPrimary || '#1890ff'
                                }}/>
                            )}
                            <div className="mobile-title">{web?.info?.name || Settings?.title}</div>
                        </div>
                        <div className="mobile-subtitle">为您提供专业、安全、高效的一站式网络服务平台</div>
                    </div>

                    <div className={styles.loginContent}>
                        {checkingLogin ? (
                            <div style={{textAlign: 'center', padding: '80px 0'}}>
                                <Spin size="large"/>
                            </div>
                        ) : user?.web?.id && !showRelogin ? (
                            <>
                                <div className={styles.loginTitle}>欢迎回来</div>
                                <div className={styles.loginSubtitle}>检测到您已登录</div>

                                <div className={styles.loggedInCard}>
                                    <Avatar
                                        size={80}
                                        src={user?.web?.avatar}
                                        icon={<UserOutlined/>}
                                        className="user-avatar"
                                    />
                                    <div className="user-name">
                                        {user?.web?.nick_name || '用户' + user?.web?.phone?.toString().slice(-4)}
                                    </div>
                                    <div className="user-phone">
                                        {user?.web?.phone || '未绑定手机'}
                                    </div>
                                </div>

                                <div className={styles.quickActions}>
                                    <Button
                                        type="primary"
                                        icon={<LoginOutlined/>}
                                        onClick={() => historyPush("user.index")}
                                        block
                                    >
                                        快捷登录
                                    </Button>
                                    <Button
                                        icon={<SwapOutlined/>}
                                        onClick={() => setShowRelogin(true)}
                                        block
                                    >
                                        切换账号
                                    </Button>
                                </div>

                                <div className={styles.reloginLink}>
                                    <a onClick={async () => {
                                        message.loading({content: "正在退出登录...", key: "outLogin"});
                                        await logout({
                                            onSuccess: (r) => {
                                                message.success(r?.message || "退出登录成功");
                                                deleteHeader();
                                                user?.setWeb(undefined);
                                                message.destroy("outLogin");
                                                setShowRelogin(true);
                                            },
                                            onFail: (r) => {
                                                message.error(r?.message || "退出登录失败");
                                                message.destroy("outLogin");
                                            }
                                        });
                                    }}>
                                        退出当前账号
                                    </a>
                                </div>
                            </>
                        ) : (
                            <>
                                <div className={styles.loginTitle}>欢迎登录</div>
                                <div className={styles.loginSubtitle}>请输入您的账号信息</div>

                                <Tabs
                                    items={items}
                                    activeKey={type}
                                    onChange={setType}
                                    size="large"
                                    centered
                                    className={styles.tabsWrapper}
                                />

                                <Form form={form} onFinish={handleSubmit} className={styles.formWrapper}>
                                    {type === 'account' && (
                                        <>
                                            <Form.Item name='phone' rules={[{required: true, message: '请输入手机号'}]}>
                                                <Input
                                                    prefix={
                                                        <div>
                                                            <span className={styles.phonePrefix}>+86</span>
                                                            <span className={styles.divider}>|</span>
                                                        </div>
                                                    }
                                                    placeholder='请输入账号或手机号'
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
                                                    message?.warning("请使用短信验证登陆后修改密码");
                                                    setType("phone");
                                                }}>忘记密码</a>
                                            </div>
                                        </>
                                    )}

                                    {type === 'phone' && (
                                        <>
                                            <Form.Item name='phone' rules={[{required: true, message: '请输入手机号'}]}>
                                                <Input
                                                    prefix={
                                                        <div>
                                                            <span className={styles.phonePrefix}>+86</span>
                                                            <span className={styles.divider}>|</span>
                                                        </div>
                                                    }
                                                    placeholder='请输入手机号'
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
                                                                    const phone = form.getFieldValue('phone' as any);
                                                                    if (!phone) {
                                                                        message?.error('请先输入手机号');
                                                                        return;
                                                                    }
                                                                    handleSendSms();
                                                                }}
                                                            >
                                                        {isSmsSendLoading ? '发送中...' : '发送验证码'}
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
                                            block
                                        >
                                            {type === 'phone' ? '登录 / 注册' : '登 录'}
                                        </Button>
                                    </Form.Item>

                                    <Form.Item style={{marginBottom: 0}} className={styles.agreement}>
                                        <Checkbox checked={isRead} onChange={(e) => setIsRead(e.target.checked)}>
                                            我已阅读并同意 <a href=''>《用户使用条款》</a>
                                        </Checkbox>
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