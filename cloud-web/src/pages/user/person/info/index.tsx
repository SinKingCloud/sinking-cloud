import React, {useRef} from "react";
import {useModel} from "umi";
import {Avatar, Button, Card, Col, Row, Tag, Typography} from "antd";
import {createStyles} from "antd-style";
import {Body, Icon, Title, useTheme} from "sinking-antd";
import {ago} from "@/utils/time";
import Info, {InfoRef} from "./components/info";
import Account, {AccountRef} from "./components/account";
import Password, {PasswordRef} from "./components/password";
import AvatarUpload, {AvatarRef} from "./components/avatar";
import PhoneComponent, {PhoneRef} from "./components/phone";
import EmailComponent, {EmailRef} from "./components/email";
import Verify, {VerifyRef} from "./components/verify";
import {historyPush} from "@/utils/route";

type SettingItem = {
    icon: React.ReactNode;
    title: string;
    desc?: string;
    value: React.ReactNode;
    disabled?: boolean;
    onClick?: () => void;
    readonly?: boolean;
};

const useStyles = createStyles(({css, token, responsive, isDarkMode}: any, props: any = {}) => {
    const isCompactMode = props?.isCompactMode || false;
    const panelPadding = isCompactMode ? 14 : 16;
    const avatarSize = isCompactMode ? 76 : 80;
    const panelBg = isDarkMode
        ? `linear-gradient(135deg, color-mix(in srgb, ${token.colorPrimary} 9%, ${token.colorBgContainer}) 0%, ${token.colorBgContainer} 100%)`
        : token.colorFillAlter;
    const primaryGlow = isDarkMode ? `color-mix(in srgb, ${token.colorPrimary} 42%, transparent)` : token.colorPrimaryBgHover;
    const primaryMist = isDarkMode ? `color-mix(in srgb, ${token.colorPrimary} 22%, transparent)` : token.colorPrimaryBg;
    const cyanGlow = isDarkMode ? "rgba(34, 211, 238, 0.16)" : token.colorPrimaryBg;
    const violetGlow = isDarkMode ? "rgba(139, 92, 246, 0.14)" : token.colorPrimaryBg;
    const shineGlow = isDarkMode ? "rgba(34, 211, 238, 0.18)" : "rgba(255,255,255,0.70)";

    return {
        pageRow: css`
            width: 100%;
            max-width: 1450px;
            margin: 0 auto;
        `,
        overview: css`
            width: 100%;
            overflow: hidden;
            border-radius: ${token.borderRadiusLG}px;
            background: ${token.colorBgContainer};
            box-shadow: ${isDarkMode ? "none" : "0 10px 28px rgba(15, 23, 42, 0.04)"};

            .ant-row {
                align-items: stretch;
            }

            .ant-col {
                display: flex;
            }
        `,
        profilePanel: css`
            position: relative;
            width: 100%;
            height: 100%;
            min-height: 214px;
            padding: ${panelPadding}px;
            display: grid;
            align-content: center;
            gap: 25px;
            box-shadow: ${isDarkMode ? "none" : "inset -1px 0 0 rgba(15,23,42,0.06)"};
            background: ${panelBg};
            isolation: isolate;
            overflow: hidden;

            &::before {
                content: '';
                position: absolute;
                top: -58px;
                right: -72px;
                width: 172px;
                height: 172px;
                border-radius: 50%;
                background: radial-gradient(circle, ${token.colorPrimary} 0%, ${primaryMist} 48%, transparent 78%);
                filter: blur(36px);
                opacity: ${isDarkMode ? 0.09 : 0.12};
                transform: rotate(-18deg);
                z-index: -1;
            }

            &::after {
                content: '';
                position: absolute;
                left: -68px;
                bottom: -72px;
                width: 178px;
                height: 178px;
                border-radius: 50%;
                background: radial-gradient(circle, ${token.colorPrimaryHover} 0%, ${cyanGlow} 50%, transparent 78%);
                filter: blur(38px);
                opacity: ${isDarkMode ? 0.08 : 0.1};
                transform: rotate(18deg);
                z-index: -1;
            }

            .profile-ambient,
            .profile-shine,
            .profile-watermark {
                position: absolute;
                pointer-events: none;
            }

            .profile-ambient {
                inset: 0;
                z-index: 0;
            }

            .profile-ambient::before {
                content: '';
                position: absolute;
                inset: -78px;
                background:
                    radial-gradient(circle at 22% 12%, ${primaryGlow} 0%, transparent 38%),
                    radial-gradient(circle at 82% 78%, ${cyanGlow} 0%, transparent 42%),
                    radial-gradient(circle at 64% 24%, ${violetGlow} 0%, transparent 34%);
                opacity: ${isDarkMode ? 0.42 : 0.58};
                filter: blur(3px);
                transform: rotate(-7deg) translate3d(-18px, -12px, 0);
                -webkit-mask-image: linear-gradient(135deg, rgba(0, 0, 0, 0.82), rgba(0, 0, 0, 0.34) 58%, transparent 100%);
                mask-image: linear-gradient(135deg, rgba(0, 0, 0, 0.82), rgba(0, 0, 0, 0.34) 58%, transparent 100%);
            }

            .profile-ambient::after {
                content: '';
                position: absolute;
                right: -176px;
                bottom: -154px;
                width: 500px;
                height: 500px;
                border-radius: 50%;
                background:
                    radial-gradient(circle, ${primaryMist} 0%, transparent 46%),
                    conic-gradient(from 0deg, transparent 0 12%, ${primaryGlow} 20%, transparent 34%, ${cyanGlow} 52%, transparent 72%, ${violetGlow} 86%, transparent 100%);
                filter: blur(42px);
                opacity: ${isDarkMode ? 0.22 : 0.32};
                transform: rotate(18deg) scale(1.02);
            }

            .profile-shine {
                inset: 0;
                z-index: 1;
                opacity: ${isDarkMode ? 0.20 : 0.34};
                background: radial-gradient(circle at 68% 40%, ${shineGlow} 0%, ${primaryMist} 20%, transparent 52%);
                mix-blend-mode: ${isDarkMode ? "screen" : "normal"};
            }

            .profile-watermark {
                right: -78px;
                bottom: -104px;
                z-index: 1;
                color: ${token.colorPrimary};
                font-size: 320px;
                line-height: 1;
                opacity: ${isDarkMode ? 0.17 : 0.14};
                transform: rotate(-24deg);
                filter: blur(20px);
            }

            @media (max-width: 991px) {
                box-shadow: ${isDarkMode ? "none" : "inset 0 -1px 0 rgba(15,23,42,0.06)"};
                min-height: auto;

                .profile-ambient::after {
                    right: -148px;
                    bottom: -150px;
                    width: 430px;
                    height: 430px;
                }

                .profile-watermark {
                    right: -62px;
                    bottom: -90px;
                    font-size: 260px;
                }
            }

            ${responsive.mobile} {
                padding: ${isCompactMode ? 12 : 14}px;
                gap: 25px;

                .profile-ambient::before {
                    inset: -54px;
                    background-size: 100% 100%, 100% 100%, 100% 100%;
                }

                .profile-ambient::after {
                    right: -116px;
                    bottom: -118px;
                    width: 330px;
                    height: 330px;
                    filter: blur(36px);
                }

                .profile-watermark {
                    right: -48px;
                    bottom: -68px;
                    font-size: 210px;
                    filter: blur(18px);
                }
            }
        `,
        profileMain: css`
            position: relative;
            z-index: 2;
            display: flex;
            align-items: center;
            gap: ${isCompactMode ? 14 : 16}px;
            min-width: 0;
            align-self: center;

            .avatar-wrap {
                position: relative;
                flex: none;
                cursor: pointer;
            }

            .avatar {
                width: ${avatarSize}px !important;
                height: ${avatarSize}px !important;
                border: 3px solid ${isDarkMode ? "rgba(255,255,255,0.08)" : token.colorBgContainer};
                box-shadow: 0 8px 20px rgba(15, 23, 42, ${isDarkMode ? 0.26 : 0.12});
            }

            .camera {
                position: absolute;
                right: -1px;
                bottom: -1px;
                width: 26px;
                height: 26px;
                display: flex;
                align-items: center;
                justify-content: center;
                border-radius: 50%;
                color: #fff;
                background: ${token.colorPrimary};
                border: 2px solid ${isDarkMode ? "rgba(255,255,255,0.08)" : token.colorBgContainer};
                font-size: 12px;
            }

            .identity {
                min-width: 0;
                flex: 1;
            }

            .name-line {
                display: flex;
                align-items: center;
                flex-wrap: wrap;
                gap: 8px;
                min-width: 0;
                margin-bottom: 6px;
            }

            .name {
                min-width: 0;
                margin: 0;
                color: ${token.colorText};
                font-size: ${isCompactMode ? 24 : 25}px;
                font-weight: 700;
                line-height: ${isCompactMode ? 32 : 33}px;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
            }

            .sub {
                color: ${token.colorTextTertiary};
                font-size: 13px;
                line-height: 20px;
            }

            ${responsive.mobile} {
                gap: 12px;

                .avatar {
                    width: ${isCompactMode ? 68 : 68}px !important;
                    height: ${isCompactMode ? 68 : 68}px !important;
                }

                .camera {
                    width: 24px;
                    height: 24px;
                    font-size: 11px;
                }

                .name {
                    font-size: ${isCompactMode ? 22 : 22}px;
                    line-height: ${isCompactMode ? 30 : 30}px;
                }
            }
        `,
        metaGrid: css`
            position: relative;
            z-index: 2;
            display: grid;
            grid-template-columns: repeat(3, minmax(0, 1fr));
            gap: 8px;
            margin-top: 0;

            ${responsive.mobile} {
                grid-template-columns: repeat(3, minmax(0, 1fr));
                gap: 6px;
                margin-top: 0;
            }
        `,
        metaItem: css`
            min-width: 0;
            padding: ${isCompactMode ? "9px 11px" : "8px 10px"};
            border: none;
            border-radius: ${token.borderRadius}px;
            background: rgba(255, 255, 255, ${isDarkMode ? 0.02 : 0.56});

            .label {
                display: flex;
                align-items: center;
                gap: 6px;
                margin-bottom: 4px;
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 18px;
            }

            .value {
                display: block;
                color: ${token.colorText};
                font-size: 13px;
                line-height: 20px;
                white-space: nowrap;
                overflow: hidden;
                text-overflow: ellipsis;
            }

            ${responsive.mobile} {
                padding: 7px 8px;

                .label {
                    margin-bottom: 2px;
                    font-size: ${isCompactMode ? 12 : 11}px;
                }

                .value {
                    font-size: ${isCompactMode ? 13 : 12}px;
                    line-height: ${isCompactMode ? 19 : 18}px;
                }
            }
        `,
        securityPanel: css`
            width: 100%;
            height: 100%;
            min-height: 214px;
            padding: ${panelPadding}px;
            display: flex;
            flex-direction: column;
            gap: ${isCompactMode ? 12 : 14}px;
            background: ${token.colorBgContainer};

            .panel-head {
                display: flex;
                align-items: center;
                justify-content: space-between;
                gap: 12px;
            }

            .title {
                color: ${token.colorText};
                font-size: 15px;
                font-weight: 700;
                line-height: 24px;
            }

            .desc {
                margin-top: 2px;
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 18px;
            }

            .security-main {
                display: flex;
                align-items: center;
                gap: 14px;
            }

            .security-ring {
                width: ${isCompactMode ? 66 : 72}px;
                height: ${isCompactMode ? 66 : 72}px;
                flex: none;
                position: relative;
                display: flex;
                align-items: center;
                justify-content: center;
                border-radius: 50%;
                background: ${token.colorBgContainer};
                box-shadow: 0 6px 16px rgba(15, 23, 42, ${isDarkMode ? 0.20 : 0.08});
            }

            .security-svg {
                position: absolute;
                inset: 0;
                width: 100%;
                height: 100%;
                pointer-events: none;
            }

            .security-track {
                fill: none;
                stroke: ${token.colorFillSecondary};
                stroke-width: 8;
            }

            .security-progress {
                fill: none;
                stroke: ${token.colorPrimary};
                stroke-width: 8;
                stroke-linecap: round;
                opacity: ${isDarkMode ? 0.78 : 0.58};
                transform: rotate(-90deg);
                transform-origin: 50% 50%;
                transition: stroke-dashoffset 0.2s ease;
            }

            .ring-content {
                position: relative;
                z-index: 1;
                display: flex;
                flex-direction: column;
                align-items: center;
                justify-content: center;
                line-height: 1;
            }

            .ring-value {
                color: ${token.colorText};
                font-size: ${isCompactMode ? 17 : 18}px;
                font-weight: 700;
                line-height: 20px;
            }

            .ring-label {
                margin-top: 3px;
                color: ${token.colorTextTertiary};
                font-size: 11px;
                line-height: 14px;
            }

            .security-list {
                flex: 1;
                min-width: 0;
                display: grid;
                gap: 8px;
            }

            .security-item {
                width: 100%;
                display: flex;
                align-items: center;
                justify-content: space-between;
                gap: 10px;
                min-width: 0;
                padding: 7px 9px;
                border: none;
                border-radius: ${token.borderRadius}px;
                background: ${token.colorFillQuaternary};
                cursor: pointer;
                transition: background-color 0.2s ease, color 0.2s ease;
                text-align: left;
            }

            .security-item:hover {
                background: ${token.colorFillSecondary};
            }

            .security-item:focus-visible {
                outline: 2px solid ${token.colorPrimaryBorder};
                outline-offset: 2px;
            }

            .security-name {
                min-width: 0;
                color: ${token.colorTextSecondary};
                font-size: 13px;
            }

            @media (max-width: 991px) {
                min-height: auto;
            }

            ${responsive.mobile} {
                padding: ${isCompactMode ? 12 : 14}px;
                gap: ${isCompactMode ? 12 : 14}px;

                .panel-head {
                    gap: 10px;
                }

                .security-main {
                    gap: 12px;
                }

                .security-ring {
                    width: ${isCompactMode ? 60 : 64}px;
                    height: ${isCompactMode ? 60 : 64}px;
                }

                .security-list {
                    gap: 6px;
                }

                .security-item {
                    padding: 7px 9px;
                }
            }
        `,
        settingGrid: css`
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(min(260px, 100%), 1fr));
            gap: ${isCompactMode ? 8 : 10}px;
            padding: ${isCompactMode ? 10 : 12}px;

            @media (max-width: 720px) {
                grid-template-columns: 1fr;
            }
        `,
        flatCard: css`
            overflow: hidden;

            .ant-card-body {
                overflow: hidden;
                border-radius: 0 0 ${token.borderRadiusLG}px ${token.borderRadiusLG}px;
            }
        `,
        settingItem: css`
            min-width: 0;
            display: grid;
            grid-template-columns: ${isCompactMode ? 24 : 24}px minmax(0, 1fr) minmax(100px, 38%);
            align-items: center;
            column-gap: ${isCompactMode ? 9 : 10}px;
            row-gap: 2px;
            min-height: ${isCompactMode ? 56 : 58}px;
            padding: ${isCompactMode ? "10px 12px" : "10px 12px"};
            border-radius: ${token.borderRadius}px;
            border: 1px solid ${token.colorBorderSecondary};
            background: ${token.colorBgContainer};
            transition: border-color 0.2s ease;

            &:hover {
                border-color: ${token.colorPrimaryBorder};
            }

            &:hover .item-icon {
                color: ${token.colorPrimary};
            }

            .item-icon {
                width: ${isCompactMode ? 24 : 24}px;
                height: ${isCompactMode ? 24 : 24}px;
                display: flex;
                align-items: center;
                justify-content: center;
                color: ${token.colorPrimary};
                font-size: ${isCompactMode ? 15 : 15}px;
                opacity: ${isDarkMode ? 0.82 : 0.76};
                transition: opacity 0.2s ease, color 0.2s ease;
            }

            .item-main {
                min-width: 0;
            }

            .item-action {
                min-width: 0;
                width: 100%;
                display: flex;
                align-items: center;
                justify-content: flex-end;
                gap: ${isCompactMode ? 6 : 8}px;
                white-space: nowrap;
            }

            .item-title {
                min-width: 0;
                color: ${token.colorText};
                font-size: ${isCompactMode ? 13 : 13}px;
                font-weight: 400;
                line-height: ${isCompactMode ? 20 : 20}px;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
            }

            .item-value {
                min-width: 0;
                flex: 1 1 auto;
                width: auto;
                min-height: 24px;
                display: flex;
                align-items: center;
                justify-content: flex-end;
                color: ${token.colorTextSecondary};
                font-size: ${isCompactMode ? 12 : 12}px;
                line-height: ${isCompactMode ? 20 : 20}px;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
            }

            .item-value .ant-typography {
                min-width: 0;
                max-width: 100%;
                display: block;
                color: inherit;
                font-size: inherit;
                line-height: inherit;
            }

            .item-desc {
                min-width: 0;
                color: ${token.colorTextTertiary};
                font-size: ${isCompactMode ? 11 : 11}px;
                line-height: ${isCompactMode ? 16 : 16}px;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
                max-width: 200px;
            }

            .edit-btn {
                flex: none;
                height: ${isCompactMode ? 22 : 22}px;
                min-width: 34px;
                padding-inline: 0;
                font-size: ${isCompactMode ? 12 : 12}px;
                white-space: nowrap;
            }

            ${responsive.mobile} {
                grid-template-columns: ${isCompactMode ? 24 : 24}px minmax(0, 1fr) minmax(88px, 42%);
                column-gap: ${isCompactMode ? 8 : 9}px;

                .item-desc {
                    max-width: 150px;
                }

                .item-action {
                    gap: ${isCompactMode ? 5 : 6}px;
                }
            }

            @media (max-width: 430px) {
                grid-template-columns: ${isCompactMode ? 24 : 24}px minmax(0, 1fr) minmax(80px, 44%);
                column-gap: ${isCompactMode ? 7 : 8}px;

                .item-desc {
                    max-width: none;
                }
            }
        `,
        sessionCard: css`
            min-width: 0;

            .session-link {
                height: 22px;
                padding-inline: 0;
                font-size: ${isCompactMode ? 12 : 12}px;
            }

            .session-list {
                display: grid;
                overflow: hidden;
                gap: ${isCompactMode ? 8 : 10}px;
                padding: ${isCompactMode ? 10 : 12}px;
            }

            .session-item {
                display: flex;
                align-items: center;
                gap: ${isCompactMode ? 10 : 10}px;
                min-width: 0;
                padding: ${isCompactMode ? "11px 12px" : "11px 12px"};
                border-radius: ${token.borderRadius}px;
                border: 1px solid ${token.colorBorderSecondary};
                background: ${token.colorBgContainer};
                transition: border-color 0.2s ease;
            }

            .session-item:hover {
                border-color: ${token.colorPrimaryBorder};
            }

            .session-item:hover .session-icon {
                color: ${token.colorPrimary};
            }

            .session-icon {
                width: ${isCompactMode ? 24 : 24}px;
                height: ${isCompactMode ? 24 : 24}px;
                flex: none;
                display: flex;
                align-items: center;
                justify-content: center;
                color: ${token.colorPrimary};
                font-size: ${isCompactMode ? 15 : 15}px;
                opacity: ${isDarkMode ? 0.82 : 0.76};
                transition: opacity 0.2s ease, color 0.2s ease;
            }

            .session-label {
                color: ${token.colorTextTertiary};
                font-size: ${isCompactMode ? 11 : 11}px;
                line-height: ${isCompactMode ? 16 : 16}px;
            }

            .session-value {
                display: block;
                color: ${token.colorText};
                font-size: ${isCompactMode ? 12 : 12}px;
                line-height: ${isCompactMode ? 20 : 20}px;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
            }
        `,
    };
});

export default () => {
    const theme = useTheme();
    const isCompactMode = theme?.isCompactTheme?.() || false;
    const {styles} = useStyles({isCompactMode});
    const pageGutter = isCompactMode ? 12 : 14;
    const user = useModel("user");
    const infoRef = useRef<InfoRef>({} as InfoRef);
    const accountRef = useRef<AccountRef>({} as AccountRef);
    const passwordRef = useRef<PasswordRef>({} as PasswordRef);
    const avatarRef = useRef<AvatarRef>({} as AvatarRef);
    const phoneRef = useRef<PhoneRef>({} as PhoneRef);
    const emailRef = useRef<EmailRef>({} as EmailRef);
    const verifyRef = useRef<VerifyRef>({} as VerifyRef);
    const ensureVerify = () => verifyRef.current?.ensure?.() || Promise.resolve(false);
    const currentPerson: any = user?.web || {};
    const currentSession = currentPerson?.session || {};
    const displayName = currentPerson?.nick_name || currentPerson?.account || currentPerson?.phone || currentPerson?.email || "用户";
    const hour = new Date().getHours();
    const greeting = hour < 6 ? "夜深了" : hour < 12 ? "上午好" : hour < 18 ? "下午好" : "晚上好";
    const completedSecurityCount = [currentPerson?.phone, currentPerson?.email, currentPerson?.password].filter(Boolean).length;
    const securityPercent = Math.round((completedSecurityCount / 3) * 100);
    const securityDesc = completedSecurityCount >= 3 ? "基础安全配置已完善" : `已完成 ${completedSecurityCount}/3 项基础配置`;
    const securityRingRadius = 30;
    const securityRingLength = Math.PI * 2 * securityRingRadius;
    const securityRingOffset = securityRingLength * (1 - securityPercent / 100);

    const metaItems = [
        {label: "用户ID", value: currentPerson?.id || "-", icon: <Icon type="FieldNumberOutlined"/>},
        {label: "登录账号", value: currentPerson?.account || "未设置", icon: <Icon type="IdcardOutlined"/>},
        {label: "最近登录", value: currentSession?.login_time ? ago(currentSession.login_time) : "-", icon: <Icon type="ClockCircleOutlined"/>},
    ];

    const securityItems = [
        {label: "安全手机", value: currentPerson?.phone ? "已绑定" : "未绑定", done: !!currentPerson?.phone, onClick: () => phoneRef.current?.show()},
        {label: "安全邮箱", value: currentPerson?.email ? "已绑定" : "未绑定", done: !!currentPerson?.email, onClick: () => emailRef.current?.show()},
        {label: "账户密码", value: currentPerson?.password ? "已设置" : "未设置", done: !!currentPerson?.password, onClick: () => passwordRef.current?.show()},
    ];

    const sessionItems = [
        {label: "登录设备", value: currentSession?.device || "-", icon: <Icon type="LaptopOutlined"/>},
        {label: "登录IP", value: currentSession?.ip || "-", icon: <Icon type="GlobalOutlined"/>},
        {label: "登录地点", value: currentSession?.location || "-", icon: <Icon type="EnvironmentOutlined"/>},
    ];

    const profileSettings: SettingItem[] = [
        {
            icon: <Icon type="CameraOutlined"/>,
            title: "账号头像",
            desc: "设置您的个人头像",
            value: <AvatarUpload
                ref={avatarRef}
                avatar={currentPerson?.avatar}
                nickName={currentPerson?.nick_name}
                size={36}
            />,
            onClick: () => avatarRef.current?.triggerUpload(),
        },
        {
            icon: <Icon type="EditOutlined"/>,
            title: "账号昵称",
            desc: "设置您的显示昵称",
            value: currentPerson?.nick_name || "未设置",
            onClick: () => infoRef.current?.show(),
        },
        {
            icon: <Icon type="IdcardOutlined"/>,
            title: "登录账号",
            desc: "账号仅支持设置一次",
            value: currentPerson?.account || "未设置",
            disabled: !!currentPerson?.account,
            onClick: () => accountRef.current?.show(),
        },
    ];

    const securitySettings: SettingItem[] = [
        {
            icon: <Icon type="PhoneOutlined"/>,
            title: "绑定手机",
            desc: "用于登录和安全验证",
            value: currentPerson?.phone || "未绑定",
            onClick: () => phoneRef.current?.show(),
        },
        {
            icon: <Icon type="MailOutlined"/>,
            title: "绑定邮箱",
            desc: "用于通知和安全验证",
            value: currentPerson?.email || "未绑定",
            onClick: () => emailRef.current?.show(),
        },
        {
            icon: <Icon type="SafetyCertificateOutlined"/>,
            title: "登录密码",
            desc: "用于账户登录验证",
            value: currentPerson?.password ? "已设置" : "未设置",
            onClick: () => passwordRef.current?.show(),
        },
    ];

    const renderText = (value: React.ReactNode) => {
        if (typeof value !== "string" && typeof value !== "number") {
            return value;
        }
        return (
            <Typography.Text ellipsis={{tooltip: String(value)}}>{value}</Typography.Text>
        );
    };

    const renderSettingItem = (item: SettingItem) => (
        <div className={styles.settingItem} key={item.title}>
            <span className="item-icon">{item.icon}</span>
            <div className="item-main">
                <div className="item-title">{item.title}</div>
                <div className="item-desc">{item.desc || "账户资料信息"}</div>
            </div>
            <div className="item-action">
                <div className="item-value">{renderText(item.value)}</div>
                {!item.readonly && (
                    <Button
                        className="edit-btn"
                        type="link"
                        size="small"
                        disabled={item.disabled}
                        onClick={item.onClick}
                    >
                        修改
                    </Button>
                )}
            </div>
        </div>
    );

    return (
        <Body loading={!user?.web}>
            <Row className={styles.pageRow} gutter={[0, pageGutter]}>
                <Col span={24}>
                    <Card
                        className={styles.overview}
                        variant="borderless"
                        styles={{body: {padding: 0}}}
                    >
                        <Row gutter={0}>
                            <Col xs={24} lg={15}>
                                <div className={styles.profilePanel}>
                                    <span className="profile-ambient"/>
                                    <span className="profile-shine"/>
                                    <span className="profile-watermark">
                                        <Icon type="UserOutlined"/>
                                    </span>
                                    <div className={styles.profileMain}>
                                        <div className="avatar-wrap" onClick={() => avatarRef.current?.triggerUpload()}>
                                            <Avatar
                                                src={currentPerson?.avatar}
                                                icon={<Icon type="UserOutlined"/>}
                                                className="avatar"
                                            />
                                            <span className="camera">
                                                <Icon type="CameraOutlined"/>
                                            </span>
                                        </div>
                                        <div className="identity">
                                            <div className="name-line">
                                                <h1 className="name">{displayName}</h1>
                                            </div>
                                            <div className="sub">{greeting}，{displayName}，欢迎回来</div>
                                        </div>
                                    </div>
                                    <div className={styles.metaGrid}>
                                        {metaItems.map((item) => (
                                            <div className={styles.metaItem} key={item.label}>
                                                <span className="label">{item.icon}{item.label}</span>
                                                <Typography.Text className="value" ellipsis={{tooltip: String(item.value)}}>
                                                    {item.value}
                                                </Typography.Text>
                                            </div>
                                        ))}
                                    </div>
                                </div>
                            </Col>
                            <Col xs={24} lg={9}>
                                <div className={styles.securityPanel}>
                                    <div className="panel-head">
                                        <div>
                                            <div className="title">账户安全</div>
                                            <div className="desc">{securityDesc}</div>
                                        </div>
                                        <div className="security-ring">
                                            <svg className="security-svg" viewBox="0 0 72 72" aria-hidden="true">
                                                <circle className="security-track" cx="36" cy="36" r={securityRingRadius}/>
                                                <circle
                                                    className="security-progress"
                                                    cx="36"
                                                    cy="36"
                                                    r={securityRingRadius}
                                                    strokeDasharray={securityRingLength}
                                                    strokeDashoffset={securityRingOffset}
                                                />
                                            </svg>
                                            <span className="ring-content">
                                                <span className="ring-value">{securityPercent}</span>
                                                <span className="ring-label">安全</span>
                                            </span>
                                        </div>
                                    </div>
                                    <div className="security-main">
                                        <div className="security-list">
                                            {securityItems.map((item) => (
                                                <button className="security-item" key={item.label} type="button" onClick={item.onClick}>
                                                    <span className="security-name">{item.label}</span>
                                                    <Tag color={item.done ? "success" : "warning"}>{item.value}</Tag>
                                                </button>
                                            ))}
                                        </div>
                                    </div>
                                </div>
                            </Col>
                        </Row>
                    </Card>
                </Col>
                <Col span={24}>
                    <Row gutter={[pageGutter, pageGutter]}>
                        <Col xs={24} lg={15}>
                            <Card
                                className={styles.flatCard}
                                title={<Title>账户设置</Title>}
                                variant="borderless"
                                styles={{body: {padding: 0}}}
                            >
                                <div className={styles.settingGrid}>
                                    {[...profileSettings, ...securitySettings].map(renderSettingItem)}
                                </div>
                            </Card>
                        </Col>
                        <Col xs={24} lg={9}>
                            <Card
                                className={`${styles.flatCard} ${styles.sessionCard}`}
                                title={<Title>当前会话</Title>}
                                extra={(
                                    <Button
                                        className="session-link"
                                        type="link"
                                        size="small"
                                        onClick={() => historyPush("user.person.session")}
                                    >
                                        全部
                                    </Button>
                                )}
                                variant="borderless"
                                styles={{body: {padding: 0}}}
                            >
                                <div className="session-list">
                                    {sessionItems.map((item) => (
                                        <div className="session-item" key={item.label}>
                                            <span className="session-icon">{item.icon}</span>
                                            <div>
                                                <div className="session-label">{item.label}</div>
                                                <Typography.Text className="session-value" ellipsis={{tooltip: String(item.value)}}>
                                                    {item.value}
                                                </Typography.Text>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            </Card>
                        </Col>
                    </Row>
                </Col>
            </Row>

            <Info ref={infoRef} nickName={currentPerson?.nick_name}/>
            <Account ref={accountRef} account={currentPerson?.account}/>
            <Password ref={passwordRef} ensureVerify={ensureVerify}/>
            <PhoneComponent ref={phoneRef} phone={currentPerson?.phone} ensureVerify={ensureVerify}/>
            <EmailComponent ref={emailRef} email={currentPerson?.email} ensureVerify={ensureVerify}/>
            <Verify ref={verifyRef}/>
        </Body>
    );
};
