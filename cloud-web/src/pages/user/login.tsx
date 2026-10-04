import {App, Avatar, Button, Checkbox, Form, Input, Spin, Tabs} from 'antd';
import React, {useEffect, useId, useRef, useState} from 'react';
import {Body, Icon, useTheme} from 'sinking-antd';
import {useModel} from "umi";
import {login, logout} from "@/service/common/login";
import {deleteHeader, setLoginToken} from "@/utils/auth";
import Captcha, {CaptchaRef} from "@/pages/components/captcha";
import {sendSms} from "@/service/common/sms";
import {sendMail} from "@/service/common/mail";
import Settings from "@/../config/defaultSettings";
import {historyPush} from "@/utils/route";
import {createStyles} from 'antd-style';

export default () => {
    const {message, modal} = App.useApp();
    const captcha = useRef<CaptchaRef>({});
    const brandPanelRef = useRef<HTMLDivElement>(null);
    const brandVisualRef = useRef<HTMLDivElement>(null);
    const [brandVisualRunning, setBrandVisualRunning] = useState(false);
    const brandVisualId = `login-orbit-${useId().replace(/:/g, "")}`;
    const web = useModel("web");
    const user = useModel("user");
    const theme = useTheme();
    const isCompact = theme?.isCompactTheme() || false;

    const useStyles: any = React.useMemo(() =>
            createStyles(({css, responsive, token, isDarkMode}: any) => ({
                background: css`
                    position: absolute;
                    inset: 0;
                    z-index: -1;
                    pointer-events: none;
                    overflow: hidden;
                `,
                backgroundSvg: css`
                    position: absolute;
                    left: 50%;
                    top: var(--login-bg-top, 0px);
                    width: var(--login-bg-width, 100%);
                    height: auto;
                    max-width: none;
                    transform: translate(-50%, var(--login-bg-shift, 0%));
                    display: block;
                `,
                body: css`
                    padding: 0;
                `,
                container: css`
                    display: flex;
                    position: absolute;
                    flex: 1;
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
                        height: 100%;
                        overflow-y: auto;
                        overflow-x: hidden;
                        top: 0;
                        left: 0;
                    }

                    ${responsive.sm} {
                        padding-top: ${isCompact ? 48 : 48}px;
                        padding-left: 24px;
                        padding-right: 24px;
                    }
                `,
                leftPanel: css`
                    flex: 1;
                    min-width: 0;
                    display: grid;
                    grid-template-rows: auto minmax(0, 1fr) auto;
                    padding: clamp(28px, 4vw, ${isCompact ? 56 : 72}px);
                    position: relative;
                    overflow: hidden;
                    isolation: isolate;
                    background: radial-gradient(ellipse at 95% 45%,
                        color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 9 : 5}%, transparent),
                        color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 3 : 2}%, transparent) 45%, transparent 78%),
                        radial-gradient(ellipse at 25% 85%,
                            color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 7 : 4}%, transparent), transparent 65%),
                        color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 3 : 2}%, ${token.colorBgContainer});

                    &::before {
                        content: '';
                        position: absolute;
                        inset: 0;
                        z-index: -1;
                        pointer-events: none;
                        background-image: linear-gradient(color-mix(in srgb, ${token.colorPrimary} 5%, transparent) 1px, transparent 1px),
                            linear-gradient(90deg, color-mix(in srgb, ${token.colorPrimary} 5%, transparent) 1px, transparent 1px);
                        background-size: 56px 56px;
                        mask-image: radial-gradient(ellipse at 45% 55%, black, transparent 72%);
                    }

                    ${responsive.md} {
                        display: none;
                    }
                `,
                brandGlow: css`
                    display: none;
                    position: absolute;
                    top: 0;
                    left: 0;
                    z-index: -1;
                    width: 620px;
                    height: 620px;
                    border-radius: 50%;
                    pointer-events: none;
                    opacity: 0;
                    background: radial-gradient(circle,
                        color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 18 : 10}%, transparent),
                        color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 7 : 4}%, transparent) 35%, transparent 70%);
                    transform: translate3d(var(--brand-glow-x, 0px), var(--brand-glow-y, 0px), 0) translate(-50%, -50%);
                    transition: opacity .45s ease;

                    [data-hovered="true"] > & { opacity: 1; }
                    @media (hover: hover) and (pointer: fine) and (prefers-reduced-motion: no-preference) { display: block; }
                `,
                brandLogo: css`
                    display: inline-flex;
                    align-items: center;
                    justify-self: start;
                    gap: 10px;
                    padding: 0;
                    border: 0;
                    background: transparent;
                    color: ${token.colorTextHeading};
                    font: inherit;
                    font-size: ${isCompact ? 16 : 18}px;
                    font-weight: 650;
                    cursor: pointer;

                    img { width: 32px; height: 32px; object-fit: contain; }
                    .anticon { color: ${token.colorPrimary}; font-size: 32px; }
                    &:focus-visible { outline: 2px solid ${token.colorPrimary}; outline-offset: 6px; }
                `,
                brandContent: css`
                    min-height: 0;
                    max-width: ${isCompact ? 540 : 620}px;
                    width: 100%;
                    margin: auto;
                    padding-top: 24px;
                    animation: login-brand-enter .8s ease-out both;

                    @keyframes login-brand-enter {
                        from { opacity: 0; transform: translateY(16px); }
                        to { opacity: 1; transform: none; }
                    }
                    @media (prefers-reduced-motion: reduce) { animation: none; }
                    @media (max-height: 650px) { padding-top: 16px; }
                `,
                brandEyebrow: css`
                    display: flex;
                    align-items: center;
                    gap: 10px;
                    margin-bottom: 18px;
                    color: ${token.colorPrimary};
                    font-size: 10px;
                    font-weight: 500;
                    letter-spacing: 2.5px;

                    &::before { content: ''; width: 20px; height: 1px; background: currentColor; }
                `,
                brandTitle: css`
                    margin: 0;
                    color: ${token.colorTextHeading};
                    font-size: clamp(28px, 3.8vw, ${isCompact ? 48 : 56}px);
                    font-weight: 700;
                    line-height: 1.3;
                    letter-spacing: -1.5px;
                `,
                brandSubtitle: css`
                    margin: 16px 0 0;
                    font-size: ${isCompact ? 13 : 14}px;
                    color: ${token.colorTextSecondary};
                    line-height: 1.8;
                `,
                brandVisual: css`
                    --login-primary: ${token.colorPrimary};
                    --login-bright: ${token.colorPrimaryHover};
                    --login-deep: ${token.colorPrimaryActive};
                    --login-surface: ${isDarkMode ? token.colorBgElevated : token.colorBgContainer};
                    --login-glass: color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 18 : 5}%, var(--login-surface));
                    --login-soft: color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 38 : 18}%, var(--login-surface));
                    --login-edge: color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 68 : 38}%, var(--login-surface));
                    --login-highlight: ${isDarkMode ? `color-mix(in srgb, ${token.colorPrimaryHover} 65%, ${token.colorText})` : token.colorBgContainer};
                    position: relative;
                    pointer-events: none;
                    user-select: none;

                    svg { display: block; width: 100%; height: 100%; overflow: visible; }
                    .login-halo, .login-ring, .login-ring-inner { transform-box: view-box; }
                    .login-core { animation: login-core-float 8s ease-in-out infinite; }
                    .login-halo { animation: login-halo-pulse 8s ease-in-out infinite; transform-origin: 320px 210px; }
                    .login-ring { animation: login-ring-turn 48s linear infinite; transform-origin: 320px 210px; }
                    .login-ring-inner { animation: login-ring-turn 36s linear infinite reverse; transform-origin: 320px 210px; }
                    .login-signal { animation: login-signal-flow 14s linear infinite; }
                    .login-signal-secondary { animation-duration: 19s; animation-delay: -8s; }
                    .login-node { animation: login-node-light 6s ease-in-out infinite; }
                    .login-node-secondary { animation-delay: -3s; }
                    &[data-running="false"] * { animation-play-state: paused !important; }

                    @keyframes login-core-float {
                        0%, 100% { transform: translateY(3px); }
                        50% { transform: translateY(-6px); }
                    }
                    @keyframes login-ring-turn { to { transform: rotate(360deg); } }
                    @keyframes login-signal-flow { to { stroke-dashoffset: -1000; } }
                    @keyframes login-halo-pulse {
                        0%, 100% { opacity: .5; transform: scale(.96); }
                        50% { opacity: .9; transform: scale(1.04); }
                    }
                    @keyframes login-node-light {
                        0%, 100% { opacity: .35; }
                        50% { opacity: 1; }
                    }
                    @media (prefers-reduced-motion: reduce) {
                        * { animation: none !important; }
                    }
                    width: 100%;
                    height: clamp(210px, 50vh, 480px);
                    margin: 14px auto 0;
                    aspect-ratio: auto;
                    transform: perspective(1000px) translate3d(0, 0, 0) rotateX(0deg) rotateY(0deg);
                    transition: transform .65s cubic-bezier(.22, 1, .36, 1);

                    @media (hover: hover) and (pointer: fine) and (prefers-reduced-motion: no-preference) {
                        [data-hovered="true"] & {
                            transform: perspective(1000px)
                                translate3d(var(--brand-shift-x, 0px), var(--brand-shift-y, 0px), 0)
                                rotateX(var(--brand-rotate-x, 0deg)) rotateY(var(--brand-rotate-y, 0deg));
                            transition-duration: .25s;
                        }
                    }
                    @media (prefers-reduced-motion: reduce) { transform: none; transition: none; }

                    @media (max-height: 650px) { height: 34vh; margin-top: 0; }
                `,
                brandFooter: css`
                    display: flex;
                    align-items: center;
                    flex-wrap: wrap;
                    gap: 12px 24px;
                    color: ${token.colorTextTertiary};
                    font-size: 11px;

                    span { display: inline-flex; align-items: center; gap: 8px; }
                    .anticon { color: ${token.colorPrimary}; font-size: 12px; }
                `,
                rightPanel: css`
                    width: ${isCompact ? 440 : 500}px;
                    position: relative;
                    isolation: isolate;
                    --login-bg-width: 150%;
                    --login-bg-top: 50%;
                    --login-bg-shift: -50%;
                    display: flex;
                    flex-direction: column;
                    justify-content: center;
                    align-items: center;
                    padding: ${isCompact ? token.paddingXL * 2 : token.paddingXL * 2.5}px;
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
                        width: ${isCompact ? 420 : 480}px;
                        padding: ${token.paddingXL * 2}px;
                    }

                    ${responsive.md} {
                        width: 100%;
                        min-height: 100%;
                        padding: ${token.paddingXL * 2}px ${token.paddingXL}px;
                        box-shadow: none;
                        justify-content: flex-start;
                        --login-bg-width: min(1361px, 160%);
                        animation: none;
                    }

                    ${responsive.sm} {
                        padding: ${token.paddingXL}px ${token.paddingLG}px;
                        min-height: 100%;
                        --login-bg-width: 100%;
                        --login-bg-top: 0px;
                        --login-bg-shift: 0%;
                    }
                `,
                loginContent: css`
                    max-width: ${isCompact ? 350 : 450}px;
                    width: 100%;
                    position: relative;
                    z-index: 1;

                    ${responsive.md} {
                        max-width: ${isCompact ? 360 : 460}px;
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
                        margin-top: ${isCompact ? token.marginLG : token.marginXL}px;
                        margin-bottom: ${isCompact ? token.marginLG : token.marginXXL}px;
                        position: relative;
                        z-index: 1;

                        .mobile-logo-title {
                            display: flex;
                            align-items: center;
                            justify-content: center;
                            gap: ${token.marginSM}px;
                            margin-bottom: ${token.marginLG}px;
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
                            width: 50px;
                            height: 50px;
                        }
                    }
                `,
                loginTitle: css`
                    font-size: ${token.fontSizeHeading2}px;
                    font-weight: 700;
                    color: ${isDarkMode ? token.colorText : token.colorTextHeading};
                    margin-bottom: ${token.marginXXS}px;
                    text-align: left;

                    ${responsive.md} {
                        display: none;
                    }
                `,
                loginSubtitle: css`
                    font-size: ${isCompact ? token.fontSize : token.fontSizeSM}px;
                    color: ${isDarkMode ? token.colorTextSecondary : token.colorTextTertiary};
                    margin-bottom: ${token.marginXL}px;
                    text-align: left;

                    ${responsive.md} {
                        display: none;
                    }
                `,
                tabsWrapper: css`
                    .ant-tabs-nav {
                        margin-bottom: ${token.marginSM}px !important;
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
                        bottom: 5px !important;
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
                        border: 1px solid ${isDarkMode ? token.colorBorder : 'transparent'};
                        border-radius: ${token.borderRadiusLG}px;
                        box-shadow: ${isDarkMode ? `0 2px 8px rgba(0,0,0,0.15)` : 'rgba(0, 0, 0, 0.05) 0 7px 25px 0'};
                        padding: 0 ${token.paddingLG}px;
                        transition: box-shadow 0.3s ease;

                        &:hover {
                            box-shadow: ${isDarkMode ? `0 4px 12px rgba(0,0,0,0.25)` : 'rgba(0, 0, 0, 0.08) 0 10px 30px 0'};
                        }

                        &:focus-within {
                            box-shadow: ${isDarkMode
                                    ? `0 4px 12px rgba(0,0,0,0.3)`
                                    : `rgba(0, 0, 0, 0.08) 0 10px 30px 0`};
                            outline: none;
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

                            &:focus,
                            &:focus-within,
                            &.ant-input-affix-wrapper-focused {
                                box-shadow: ${isDarkMode
                                        ? `0 4px 12px rgba(0,0,0,0.3)`
                                        : `rgba(0, 0, 0, 0.08) 0 10px 30px 0`} !important;
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
                    margin-bottom: ${isCompact ? token.marginSM : token.marginLG}px;

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
                    margin-bottom: 0 !important;

                    .ant-checkbox-wrapper {
                        font-size: ${token.fontSizeSM}px;
                        color: ${isDarkMode ? token.colorTextSecondary : 'rgba(0, 0, 0, 0.45)'};
                    }
                `,
                loggedInCard: css`
                    text-align: center;
                    padding: 33px 32px 15px 32px;
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
                    gap: ${token.marginSM}px;

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
    const [isRead, setIsRead] = useState(true);
    const [countdown, setCountdown] = useState(0);
    const [showRelogin, setShowRelogin] = useState(false);
    const [checkingLogin, setCheckingLogin] = useState(true);

    useEffect(() => {
        let visible = false;
        const update = () => setBrandVisualRunning(visible && !document.hidden);
        const observer = new IntersectionObserver(([entry]) => {
            visible = entry.isIntersecting;
            update();
        });
        if (brandVisualRef.current) observer.observe(brandVisualRef.current);
        document.addEventListener("visibilitychange", update);
        return () => {
            observer.disconnect();
            document.removeEventListener("visibilitychange", update);
        };
    }, []);

    useEffect(() => {
        const panel = brandPanelRef.current;
        if (!panel) return;
        const hoverMotion = window.matchMedia("(hover: hover) and (pointer: fine) and (prefers-reduced-motion: no-preference)");
        let frame = 0;
        let pointer: {x: number; y: number} | null = null;
        const reset = () => {
            if (frame) window.cancelAnimationFrame(frame);
            frame = 0;
            pointer = null;
            panel.removeAttribute("data-hovered");
        };
        const update = () => {
            frame = 0;
            if (!pointer) return;
            const rect = panel.getBoundingClientRect();
            if (!rect.width || !rect.height) return reset();
            const x = pointer.x - rect.left;
            const y = pointer.y - rect.top;
            const horizontal = Math.max(-1, Math.min(1, x / rect.width * 2 - 1));
            const vertical = Math.max(-1, Math.min(1, y / rect.height * 2 - 1));
            panel.style.setProperty("--brand-glow-x", `${x}px`);
            panel.style.setProperty("--brand-glow-y", `${y}px`);
            panel.style.setProperty("--brand-shift-x", `${(horizontal * 6).toFixed(2)}px`);
            panel.style.setProperty("--brand-shift-y", `${(vertical * 4).toFixed(2)}px`);
            panel.style.setProperty("--brand-rotate-x", `${(-vertical * 2).toFixed(2)}deg`);
            panel.style.setProperty("--brand-rotate-y", `${(horizontal * 3).toFixed(2)}deg`);
            panel.setAttribute("data-hovered", "true");
        };
        const onPointerMove = (event: PointerEvent) => {
            if (!hoverMotion.matches || event.pointerType !== "mouse") return reset();
            pointer = {x: event.clientX, y: event.clientY};
            if (!frame) frame = window.requestAnimationFrame(update);
        };
        panel.addEventListener("pointermove", onPointerMove, {passive: true});
        panel.addEventListener("pointerleave", reset);
        panel.addEventListener("pointercancel", reset);
        window.addEventListener("scroll", reset, {passive: true, capture: true});
        window.addEventListener("resize", reset);
        window.addEventListener("blur", reset);
        document.addEventListener("visibilitychange", reset);
        hoverMotion.addEventListener("change", reset);
        return () => {
            reset();
            panel.removeEventListener("pointermove", onPointerMove);
            panel.removeEventListener("pointerleave", reset);
            panel.removeEventListener("pointercancel", reset);
            window.removeEventListener("scroll", reset, true);
            window.removeEventListener("resize", reset);
            window.removeEventListener("blur", reset);
            document.removeEventListener("visibilitychange", reset);
            hoverMotion.removeEventListener("change", reset);
        };
    }, []);

    /**
     * 初始化时检查登录状态
     */
    useEffect(() => {
        setCheckingLogin(true);
        user?.refreshWebUser?.((d: any) => {
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

    /**
     * 统一登录处理
     * @param values 表单数据
     * @param captchaRes 验证码
     * @param loginType 登陆类型
     */
    const handleLogin = async (values: any, captchaRes: any, loginType: string) => {
        const body: any = {
            account: values.account,
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

        if (!captcha?.current?.Show) {
            return;
        }
        setIsCodeSendLoading(true);
        captcha.current.Show(
            async (res) => {
                await sendFn({
                    body: {
                        token: res?.token,
                        captcha_x: res?.x,
                        captcha_y: res?.y,
                        ...targetValue
                    },
                    onSuccess: (r) => {
                        message?.success(r?.message);
                        startCountdown();
                    },
                    onFail: (r) => message?.error(r?.message),
                    onFinally: () => setIsCodeSendLoading(false)
                });
            },
            () => setIsCodeSendLoading(false)
        );
    };

    /**
     * 表单提交
     * @param values 表单数据
     */
    const handleSubmit = async (values: any) => {
        if (!isRead) {
            message?.error("请阅读并同意《用户使用条款》协议");
            return;
        }

        // 验证输入
        if (type === "account") {
            if (!values.account || values.account.length === 0) {
                message?.error("请输入账号或手机号或邮箱");
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
        if (!captcha?.current?.Show) {
            return;
        }
        setIsLoading(true);
        captcha.current.Show(
            async (res) => {
                await handleLogin(values, res, type);
            },
            () => setIsLoading(false)
        );
    };

    /**
     * 退出登录
     */
    const outLogin = async () => {
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
    }

    const items: any = [
        {key: 'account', label: '密码登录'},
        {key: 'code', label: '验证登录'}
    ];

    return (
        <Body space={false} animation={false} className={styles.body}>
            <div className={styles.container}>
                <Captcha ref={captcha}/>
                <div ref={brandPanelRef} className={styles.leftPanel}>
                    <div className={styles.brandGlow} aria-hidden="true"/>
                    <button type="button" className={styles.brandLogo} onClick={() => historyPush("index.index")}
                            aria-label="返回首页">
                        {web?.info?.logo ? <img src={web.info.logo} alt=""/> : <Icon type="icon-logo"/>}
                        <span>{web?.info?.name || Settings?.title}</span>
                    </button>
                    <div className={styles.brandContent}>
                        <div className={styles.brandEyebrow}>YOUR SPACE. YOUR IDEAS.</div>
                        <h1 className={styles.brandTitle}>让想法成为现实</h1>
                        <p className={styles.brandSubtitle}>连接你的工作空间，继续创造更多可能。</p>
                        <div ref={brandVisualRef} className={styles.brandVisual} data-running={brandVisualRunning} aria-hidden="true">
                            <svg viewBox="0 0 640 420" fill="none">
                                <defs>
                                    <radialGradient id={`${brandVisualId}-aura`}>
                                        <stop stopColor="var(--login-primary)" stopOpacity=".15"/>
                                        <stop offset="1" stopColor="var(--login-primary)" stopOpacity="0"/>
                                    </radialGradient>
                                    <radialGradient id={`${brandVisualId}-shadow`}>
                                        <stop stopColor="var(--login-primary)" stopOpacity=".15"/>
                                        <stop offset="1" stopColor="var(--login-primary)" stopOpacity="0"/>
                                    </radialGradient>
                                    <linearGradient id={`${brandVisualId}-ring`} x1="201" y1="76" x2="448" y2="338" gradientUnits="userSpaceOnUse">
                                        <stop stopColor="var(--login-highlight)"/>
                                        <stop offset=".28" stopColor="var(--login-soft)"/>
                                        <stop offset=".6" stopColor="var(--login-primary)" stopOpacity=".75"/>
                                        <stop offset="1" stopColor="var(--login-highlight)"/>
                                    </linearGradient>
                                    <linearGradient id={`${brandVisualId}-rim`} x1="194" y1="95" x2="405" y2="317" gradientUnits="userSpaceOnUse">
                                        <stop stopColor="var(--login-highlight)"/>
                                        <stop offset=".45" stopColor="var(--login-primary)"/>
                                        <stop offset="1" stopColor="var(--login-highlight)"/>
                                    </linearGradient>
                                    <radialGradient id={`${brandVisualId}-core`} cx=".28" cy=".23" r=".85">
                                        <stop stopColor="var(--login-highlight)"/>
                                        <stop offset=".3" stopColor="var(--login-glass)"/>
                                        <stop offset=".65" stopColor="var(--login-soft)"/>
                                        <stop offset="1" stopColor="var(--login-primary)" stopOpacity=".75"/>
                                    </radialGradient>
                                    <radialGradient id={`${brandVisualId}-reflection`} cx=".32" cy=".22" r=".68">
                                        <stop stopColor="var(--login-highlight)" stopOpacity=".8"/>
                                        <stop offset=".6" stopColor="var(--login-highlight)" stopOpacity="0"/>
                                    </radialGradient>
                                    <linearGradient id={`${brandVisualId}-orbit`} x1="98" y1="138" x2="541" y2="301" gradientUnits="userSpaceOnUse">
                                        <stop stopColor="var(--login-edge)" stopOpacity=".25"/>
                                        <stop offset=".45" stopColor="var(--login-primary)" stopOpacity=".6"/>
                                        <stop offset="1" stopColor="var(--login-edge)" stopOpacity=".35"/>
                                    </linearGradient>
                                    <path id={`${brandVisualId}-track`} pathLength="1000" d="M91 258C65 208 145 132 269 96C394 60 526 79 551 129C577 179 498 254 373 290C249 326 117 308 91 258Z"/>
                                    <path id={`${brandVisualId}-track-secondary`} pathLength="1000" d="M93 138C116 95 242 112 372 175C502 238 588 323 565 366C542 409 416 392 286 329C156 266 70 181 93 138Z"
                                          transform="translate(0 -35)"/>
                                </defs>

                                <ellipse className="login-halo" cx="320" cy="210" rx="260" ry="178" fill={`url(#${brandVisualId}-aura)`}/>
                                <ellipse cx="320" cy="377" rx="184" ry="24" fill={`url(#${brandVisualId}-shadow)`}/>
                                <g stroke="var(--login-edge)" strokeWidth=".7" opacity=".4">
                                    <path d="M91 349H173M467 349H549M320 43V64M320 354V366"/>
                                    <path d="M101 344V354M111 346V352M121 346V352M519 346V352M529 346V352M539 344V354"/>
                                    <circle cx="80" cy="349" r="2"/><circle cx="560" cy="349" r="2"/>
                                </g>

                                <use href={`#${brandVisualId}-track`} stroke={`url(#${brandVisualId}-orbit)`} strokeWidth="1"/>
                                <use href={`#${brandVisualId}-track-secondary`} stroke="var(--login-edge)" strokeWidth=".8" strokeOpacity=".45"/>

                                <g transform="translate(320 210) rotate(-28) scale(.82 1) translate(-320 -210)">
                                    <circle cx="326" cy="217" r="125" stroke="var(--login-edge)" strokeWidth="18" strokeOpacity=".18"/>
                                    <circle cx="320" cy="210" r="125" stroke={`url(#${brandVisualId}-ring)`} strokeWidth="18" strokeOpacity=".8"/>
                                    <circle cx="320" cy="210" r="135" stroke={`url(#${brandVisualId}-rim)`} strokeWidth="1.2"/>
                                    <circle cx="320" cy="210" r="115" stroke="var(--login-highlight)" strokeWidth="1" strokeOpacity=".8"/>
                                    <g className="login-ring">
                                        <circle cx="320" cy="210" r="144" stroke="var(--login-edge)" strokeWidth="2" strokeDasharray="1 12" strokeOpacity=".5"/>
                                        <path d="M320 66A144 144 0 0 1 453 155M320 354A144 144 0 0 1 187 265"
                                              stroke="var(--login-primary)" strokeWidth="3" strokeLinecap="round" strokeOpacity=".8"/>
                                        <circle cx="320" cy="66" r="4" fill="var(--login-bright)"/>
                                        <circle cx="320" cy="354" r="3" fill="var(--login-highlight)" stroke="var(--login-primary)"/>
                                    </g>
                                    <g className="login-ring-inner" stroke="var(--login-edge)" strokeWidth=".8" strokeOpacity=".6">
                                        <path d="M213 210A107 107 0 0 1 320 103M427 210A107 107 0 0 1 320 317"/>
                                        <path d="M210 201V219M311 100H329M430 201V219M311 320H329"/>
                                    </g>
                                </g>

                                <g className="login-core">
                                    <circle cx="320" cy="210" r="70" fill={`url(#${brandVisualId}-core)`} stroke="var(--login-edge)" strokeOpacity=".6"/>
                                    <circle cx="320" cy="210" r="69" fill={`url(#${brandVisualId}-reflection)`}/>
                                    <g stroke="var(--login-highlight)" strokeWidth=".9" strokeOpacity=".7">
                                        <ellipse cx="320" cy="210" rx="34" ry="69" transform="rotate(-28 320 210)"/>
                                        <ellipse cx="320" cy="210" rx="69" ry="24" transform="rotate(-28 320 210)"/>
                                        <path d="M271 163C295 149 330 147 348 159M352 266C332 279 306 281 289 272"/>
                                    </g>
                                    <path d="M315 195L332 199L326 216L309 212Z" fill="var(--login-primary)" fillOpacity=".85"/>
                                    <path d="M315 195L322 184L339 189L332 199Z" fill="var(--login-highlight)" fillOpacity=".9"/>
                                    <path d="M332 199L339 189L333 207L326 216Z" fill="var(--login-deep)" fillOpacity=".8"/>
                                    <circle cx="291" cy="177" r="4" fill="var(--login-highlight)" fillOpacity=".8"/>
                                </g>

                                <path d="M551 129C577 179 498 254 373 290C249 326 117 308 91 258" stroke={`url(#${brandVisualId}-orbit)`} strokeWidth="1.3"/>
                                <path d="M93 103C70 146 156 231 286 294C416 357 542 374 565 331" stroke="var(--login-edge)" strokeWidth="1" strokeOpacity=".5"/>
                                <use className="login-signal" href={`#${brandVisualId}-track`} stroke="var(--login-bright)"
                                     strokeWidth="4" strokeDasharray="2 998" strokeLinecap="round"/>
                                <use className="login-signal" href={`#${brandVisualId}-track`} stroke="var(--login-primary)"
                                     strokeWidth="10" strokeDasharray="1 999" strokeLinecap="round" strokeOpacity=".1"/>
                                <use className="login-signal login-signal-secondary" href={`#${brandVisualId}-track-secondary`}
                                     stroke="var(--login-primary)" strokeWidth="3" strokeDasharray="2 998" strokeLinecap="round"/>
                                <g className="login-node">
                                    <circle cx="139" cy="276" r="10" fill="var(--login-primary)" fillOpacity=".06"/>
                                    <circle cx="139" cy="276" r="3" fill="var(--login-bright)"/>
                                </g>
                                <g className="login-node login-node-secondary">
                                    <circle cx="512" cy="114" r="11" fill="var(--login-primary)" fillOpacity=".06"/>
                                    <circle cx="512" cy="114" r="3" fill="var(--login-bright)"/>
                                </g>
                            </svg>
                        </div>
                    </div>
                    <div className={styles.brandFooter}>
                        <span><Icon type="CheckOutlined"/>灵感记录</span>
                        <span><Icon type="CheckOutlined"/>便捷管理</span>
                        <span><Icon type="CheckOutlined"/>安全守护</span>
                    </div>
                </div>
                <div className={styles.rightPanel}>
                    <div className={styles.background} aria-hidden="true">
                        <svg className={styles.backgroundSvg} viewBox="0 0 1361 609" focusable="false">
                            <g stroke="none" strokeWidth="1" fill="none" fillRule="evenodd">
                                <g transform="translate(-79.000000, -82.000000)">
                                    <g transform="translate(77.000000, 73.000000)">
                                        <g opacity="0.8" transform="translate(74.901416, 569.699158) rotate(-7.000000) translate(-74.901416, -569.699158) translate(4.901416, 525.199158)">
                                            <ellipse fill="#CFDAE6" opacity="0.25" cx="63.5748792" cy="32.468367" rx="21.7830479" ry="21.766008"/>
                                            <ellipse fill="#CFDAE6" opacity="0.599999964" cx="5.98746479" cy="13.8668601" rx="5.2173913" ry="5.21330997"/>
                                            <path d="M38.1354514,88.3520215 C43.8984227,88.3520215 48.570234,83.6838647 48.570234,77.9254015 C48.570234,72.1669383 43.8984227,67.4987816 38.1354514,67.4987816 C32.3724801,67.4987816 27.7006688,72.1669383 27.7006688,77.9254015 C27.7006688,83.6838647 32.3724801,88.3520215 38.1354514,88.3520215 Z" fill="#CFDAE6" opacity="0.45"/>
                                            <path d="M64.2775582,33.1704963 L119.185836,16.5654915" stroke="#CFDAE6" strokeWidth="1.73913043" strokeLinecap="round" strokeLinejoin="round"/>
                                            <path d="M42.1431708,26.5002681 L7.71190162,14.5640702" stroke="#E0B4B7" strokeWidth="0.702678964" opacity="0.7" strokeLinecap="round" strokeLinejoin="round" strokeDasharray="1.405357899873153,2.108036953469981"/>
                                            <path d="M63.9262187,33.521561 L43.6721326,69.3250951" stroke="#BACAD9" strokeWidth="0.702678964" strokeLinecap="round" strokeLinejoin="round" strokeDasharray="1.405357899873153,2.108036953469981"/>
                                            <g transform="translate(126.850922, 13.543654) rotate(30.000000) translate(-126.850922, -13.543654) translate(117.285705, 4.381889)" fill="#CFDAE6">
                                                <ellipse opacity="0.45" cx="9.13482653" cy="9.12768076" rx="9.13482653" ry="9.12768076"/>
                                                <path d="M18.2696531,18.2553615 C18.2696531,13.2142826 14.1798519,9.12768076 9.13482653,9.12768076 C4.08980114,9.12768076 0,13.2142826 0,18.2553615 L18.2696531,18.2553615 Z" transform="translate(9.134827, 13.691521) scale(-1, -1) translate(-9.134827, -13.691521) "/>
                                            </g>
                                        </g>
                                        <g transform="translate(216.294700, 123.725600) rotate(-5.000000) translate(-216.294700, -123.725600) translate(106.294700, 35.225600)">
                                            <ellipse fill="#CFDAE6" opacity="0.25" cx="29.1176471" cy="29.1402439" rx="29.1176471" ry="29.1402439"/>
                                            <ellipse fill="#CFDAE6" opacity="0.3" cx="29.1176471" cy="29.1402439" rx="21.5686275" ry="21.5853659"/>
                                            <ellipse stroke="#CFDAE6" opacity="0.4" cx="179.019608" cy="138.146341" rx="23.7254902" ry="23.7439024"/>
                                            <ellipse fill="#BACAD9" opacity="0.5" cx="29.1176471" cy="29.1402439" rx="10.7843137" ry="10.7926829"/>
                                            <path d="M29.1176471,39.9329268 L29.1176471,18.347561 C23.1616351,18.347561 18.3333333,23.1796097 18.3333333,29.1402439 C18.3333333,35.1008781 23.1616351,39.9329268 29.1176471,39.9329268 Z" fill="#BACAD9"/>
                                            <g opacity="0.45" transform="translate(172.000000, 131.000000)" fill="#E6A1A6">
                                                <ellipse cx="7.01960784" cy="7.14634146" rx="6.47058824" ry="6.47560976"/>
                                                <path d="M0.549019608,13.6219512 C4.12262681,13.6219512 7.01960784,10.722722 7.01960784,7.14634146 C7.01960784,3.56996095 4.12262681,0.670731707 0.549019608,0.670731707 L0.549019608,13.6219512 Z" transform="translate(3.784314, 7.146341) scale(-1, 1) translate(-3.784314, -7.146341) "/>
                                            </g>
                                            <ellipse fill="#CFDAE6" cx="218.382353" cy="138.685976" rx="1.61764706" ry="1.61890244"/>
                                            <ellipse fill="#E0B4B7" opacity="0.35" cx="179.558824" cy="175.381098" rx="1.61764706" ry="1.61890244"/>
                                            <ellipse fill="#E0B4B7" opacity="0.35" cx="180.098039" cy="102.530488" rx="2.15686275" ry="2.15853659"/>
                                            <path d="M28.9985381,29.9671598 L171.151018,132.876024" stroke="#CFDAE6" opacity="0.8"/>
                                        </g>
                                        <g opacity="0.799999952" transform="translate(1054.100635, 36.659317) rotate(-11.000000) translate(-1054.100635, -36.659317) translate(1026.600635, 4.659317)">
                                            <ellipse stroke="#CFDAE6" strokeWidth="0.941176471" cx="43.8135593" cy="32" rx="11.1864407" ry="11.2941176"/>
                                            <g transform="translate(34.596774, 23.111111)" fill="#BACAD9">
                                                <ellipse opacity="0.45" cx="9.18534718" cy="8.88888889" rx="8.47457627" ry="8.55614973"/>
                                                <path d="M9.18534718,17.4450386 C13.8657264,17.4450386 17.6599235,13.6143199 17.6599235,8.88888889 C17.6599235,4.16345787 13.8657264,0.332739156 9.18534718,0.332739156 L9.18534718,17.4450386 Z"/>
                                            </g>
                                            <path d="M34.6597385,24.809694 L5.71666084,4.76878945" stroke="#CFDAE6" strokeWidth="0.941176471"/>
                                            <ellipse stroke="#CFDAE6" strokeWidth="0.941176471" cx="3.26271186" cy="3.29411765" rx="3.26271186" ry="3.29411765"/>
                                            <ellipse fill="#F7E1AD" cx="2.79661017" cy="61.1764706" rx="2.79661017" ry="2.82352941"/>
                                            <path d="M34.6312443,39.2922712 L5.06366663,59.785082" stroke="#CFDAE6" strokeWidth="0.941176471"/>
                                        </g>
                                        <g opacity="0.33" transform="translate(1282.537219, 446.502867) rotate(-10.000000) translate(-1282.537219, -446.502867) translate(1142.537219, 327.502867)">
                                            <g transform="translate(141.333539, 104.502742) rotate(275.000000) translate(-141.333539, -104.502742) translate(129.333539, 92.502742)" fill="#BACAD9">
                                                <circle opacity="0.45" cx="11.6666667" cy="11.6666667" r="11.6666667"/>
                                                <path d="M23.3333333,23.3333333 C23.3333333,16.8900113 18.1099887,11.6666667 11.6666667,11.6666667 C5.22334459,11.6666667 0,16.8900113 0,23.3333333 L23.3333333,23.3333333 Z" transform="translate(11.666667, 17.500000) scale(-1, -1) translate(-11.666667, -17.500000) "/>
                                            </g>
                                            <circle fill="#CFDAE6" cx="201.833333" cy="87.5" r="5.83333333"/>
                                            <path d="M143.5,88.8126685 L155.070501,17.6038544" stroke="#BACAD9" strokeWidth="1.16666667"/>
                                            <path d="M17.5,37.3333333 L127.466252,97.6449735" stroke="#BACAD9" strokeWidth="1.16666667"/>
                                            <polyline stroke="#CFDAE6" strokeWidth="1.16666667" points="143.902597 120.302281 174.935455 231.571342 38.5 147.510847 126.366941 110.833333"/>
                                            <path d="M159.833333,99.7453842 L195.416667,89.25" stroke="#E0B4B7" strokeWidth="1.16666667" opacity="0.6"/>
                                            <path d="M205.333333,82.1372105 L238.719406,36.1666667" stroke="#BACAD9" strokeWidth="1.16666667"/>
                                            <path d="M266.723424,132.231988 L207.083333,90.4166667" stroke="#CFDAE6" strokeWidth="1.16666667"/>
                                            <circle fill="#C1D1E0" cx="156.916667" cy="8.75" r="8.75"/>
                                            <circle fill="#C1D1E0" cx="39.0833333" cy="148.75" r="5.25"/>
                                            <circle fillOpacity="0.6" fill="#D1DEED" cx="8.75" cy="33.25" r="8.75"/>
                                            <circle fillOpacity="0.6" fill="#D1DEED" cx="243.833333" cy="30.3333333" r="5.83333333"/>
                                            <circle fill="#E0B4B7" cx="175.583333" cy="232.75" r="5.25"/>
                                        </g>
                                    </g>
                                </g>
                            </g>
                        </svg>
                    </div>
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
                                    color: theme?.themes?.token?.colorPrimary || '#1890ff'
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
                                        src={user?.web?.avatar || false}
                                        icon={<Icon type="UserOutlined"/>}
                                        className="user-avatar"
                                    />
                                    <div className="user-name">
                                        {user?.web?.nick_name || '用户' + (user?.web?.phone || user?.web?.email || user?.web?.account)?.toString().slice(0, 4)}
                                    </div>
                                    <div className="user-contact">
                                        {user?.web?.phone || user?.web?.email || user?.web?.account || '未绑定联系方式'}
                                    </div>
                                </div>
                                <div className={styles.quickActions}>
                                    <Button type="primary" icon={<Icon type="LoginOutlined"/>}
                                            onClick={() => historyPush("user.index")} block>
                                        快捷登录
                                    </Button>
                                    <Button icon={<Icon type="SwapOutlined"/>} onClick={() => setShowRelogin(true)} block>
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
                                <div className={styles.loginTitle}>欢迎登录</div>
                                <div className={styles.loginSubtitle}>请输入您的账号信息</div>
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
                                    <Form.Item style={{marginBottom: 12}}>
                                        <Button type='primary' loading={isLoading} htmlType='submit' block>
                                            {type === 'code' ? '登录 / 注册' : '登 录'}
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
