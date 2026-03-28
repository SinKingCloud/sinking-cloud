import React, {useEffect, useState} from "react";
import {Body, Icon, useTheme} from "sinking-antd";
import {useModel} from "umi";
import {Row, Col, Statistic} from "antd";
import {createStyles} from 'antd-style';
import {
    RocketOutlined,
    CheckCircleOutlined,
    StarFilled,
    MenuOutlined,
    LoginOutlined,
    CloseOutlined,
    LogoutOutlined,
    HomeOutlined,
} from "@ant-design/icons";
import Settings from "@/../config/defaultSettings";
import {deleteHeader} from "@/utils/auth";
import {logout} from "@/service/common/login";
import {historyPush} from "@/utils/route";

/**
 * 处理资源路径，正确拼接 basePath
 */
const getAssetPath = (path: string): string => {
    const basePath = (Settings?.basePath || '').replace(/\/$/, ''); // 去掉结尾的斜杠
    const srcPath = path.replace(/^\//, ''); // 去掉开头的斜杠
    return basePath ? `${basePath}/${srcPath}` : `/${srcPath}`;
};

// 优化后的动画元素配置 - 减少数量提升性能
const ANIMATION_ELEMENTS = {
    // 减少粒子数量 4 -> 2
    particles: Array.from({length: 2}, (_, i) => ({
        id: i,
        left: `${30 + i * 40}%`,
        top: `${40 + i * 20}%`,
        duration: `${20 + i * 4}s`,
        delay: `${i * 3}s`,
        tx: `${(i === 0 ? 100 : -100)}px`,
        ty: `${(i === 0 ? 80 : -80)}px`
    })),
    // 减少光球数量 5 -> 2
    glowOrbs: Array.from({length: 2}, (_, i) => ({
        id: i,
        left: `${25 + i * 50}%`,
        top: `${35 + i * 15}%`,
        size: `${250 + i * 50}px`,
        duration: `${18 + i * 4}s`,
        delay: `${i * 3}s`,
        tx: `${(i === 0 ? 60 : -60)}px`,
        ty: `${(i === 0 ? 50 : -50)}px`
    }))
};

const useStyles: any = createStyles(({token, css, isDarkMode}) => ({
    title: css`
        margin: 0;
        padding: 0;
        font-weight: 600;
        line-height: 1.2;
        color: ${token.colorText};
    `,
    h1: css`
        font-size: 56px;
        font-weight: 800;
        margin-bottom: 24px;
        line-height: 1.1;

        @media (max-width: 768px) {
            font-size: 32px;
            margin-bottom: 16px;
        }
    `,
    h2: css`
        font-size: 48px;
        font-weight: 800;
        margin-bottom: 16px;

        @media (max-width: 768px) {
            font-size: 32px;
            margin-bottom: 12px;
        }
    `,
    h3: css`
        font-size: 24px;
        font-weight: 600;
        margin-bottom: 48px;

        @media (max-width: 768px) {
            font-size: 18px;
            margin-bottom: 40px;
        }
    `,
    h4: css`
        font-size: 18px;
        font-weight: 600;
        margin-bottom: 10px;

        @media (max-width: 768px) {
            font-size: 16px;
            margin-bottom: 8px;
        }
    `,
    paragraph: css`
        margin: 0;
        padding: 0;
        line-height: 1.6;
        color: ${token.colorTextSecondary};
    `,
    divider: css`
        height: 1px;
        background: ${token.colorBorderSecondary};
        border: none;
        margin: 24px 0;

        @media (max-width: 768px) {
            margin: 20px 0;
        }
    `,
    pageWrapper: css`
        background: ${token.colorBgLayout};
        width: 100%;
        min-height: 100%;
        overflow-x: hidden;
        margin: 0;
        padding: 0;
        position: relative;
        isolation: isolate;
        overscroll-behavior-y: none;
        -webkit-overflow-scrolling: touch;
    `,
    noPadding: css`
        padding: 0 !important;
    `,
    navbar: css`
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        z-index: 9999;
        padding: 0;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        pointer-events: auto;
        isolation: isolate;
        -webkit-transform: translateZ(0);
        transform: translateZ(0);

        &::after {
            content: '';
            position: absolute;
            bottom: 0;
            left: 0;
            right: 0;
            height: 1px;
            background: ${token.colorBorderSecondary};
            opacity: 0;
            transition: opacity 0.3s ease;
        }

        &.scrolled::after {
            opacity: 1;
        }
    `,
    contentWrapper: css`
        padding-top: 0;
        position: relative;
        z-index: 1;
    `,
    navWrapper: css`
        background: ${isDarkMode
                ? 'rgba(10, 10, 10, 0.5)'
                : 'rgba(255, 255, 255, 0.55)'};
        backdrop-filter: blur(20px) saturate(160%);
        -webkit-backdrop-filter: blur(20px) saturate(160%);
        transition: background 0.3s ease, box-shadow 0.3s ease, border-color 0.3s ease;
        border-bottom: 1px solid ${isDarkMode
                ? 'rgba(255, 255, 255, 0.08)'
                : 'transparent'};
        position: relative;
        z-index: 1;
        transform: translateZ(0);
        will-change: transform;

        .scrolled & {
            background: ${isDarkMode
                    ? 'rgba(10, 10, 10, 0.75)'
                    : 'rgba(255, 255, 255, 0.8)'};
            backdrop-filter: blur(24px) saturate(170%);
            -webkit-backdrop-filter: blur(24px) saturate(170%);
            box-shadow: ${isDarkMode
                    ? '0 1px 0 rgba(255, 255, 255, 0.05), 0 4px 16px rgba(0, 0, 0, 0.4)'
                    : '0 1px 0 rgba(0, 0, 0, 0.05)'};
            border-bottom-color: ${isDarkMode
                    ? 'rgba(255, 255, 255, 0.1)'
                    : 'transparent'};
        }
    `,
    themeToggle: css`
        width: 36px;
        height: 36px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 18px;
        color: ${token.colorText};
        cursor: pointer;
        border-radius: ${token.borderRadius}px;
        background: transparent;
        border: none;
        transition: all 0.2s ease;
        transform-origin: center;
        backface-visibility: hidden;
        -webkit-font-smoothing: antialiased;

        svg {
            transform-origin: center;
            backface-visibility: hidden;
        }

        &:hover {
            background: ${token.colorFillQuaternary};
        }

        &:active {
            transform: scale(0.95);
        }
    `,
    navContainer: css`
        max-width: 1400px;
        margin: 0 auto;
        padding: 0 max(24px, env(safe-area-inset-right));
        height: 64px;
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 32px;

        @media (max-width: 768px) {
            height: 56px;
            padding: 0 16px;
            gap: 16px;
        }
    `,
    logo: css`
        display: flex;
        align-items: center;
        gap: 12px;
        cursor: pointer;
        transition: opacity 0.2s ease;
        flex-shrink: 0;
        min-width: 200px;

        &:hover {
            opacity: 0.8;
        }
    `,
    logoIcon: css`
        width: 32px;
        height: 32px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 32px;
        color: ${token.colorPrimary};
        transition: all 0.2s ease;
        flex-shrink: 0;

        img {
            width: 32px;
            height: 32px;
            object-fit: contain;
            border-radius: ${token.borderRadius}px;
        }

        @media (max-width: 768px) {
            width: 28px;
            height: 28px;
            font-size: 28px;

            img {
                width: 28px;
                height: 28px;
                border-radius: ${token.borderRadiusSM}px;
            }
        }
    `,
    logoText: css`
        font-size: 16px;
        font-weight: 600;
        color: ${token.colorText};
        letter-spacing: -0.01em;

        @media (max-width: 768px) {
            font-size: 15px;
        }
    `,
    navMenu: css`
        display: flex;
        align-items: center;
        gap: 8px;
        flex: 1;
        justify-content: center;

        @media (max-width: 968px) {
            display: none;
        }
    `,
    navLink: css`
        font-size: 14px;
        color: ${token.colorTextSecondary};
        cursor: pointer;
        transition: all 0.2s ease;
        font-weight: 500;
        position: relative;
        padding: 8px 16px;
        border-radius: ${token.borderRadius}px;
        white-space: nowrap;
        display: flex;
        align-items: center;
        justify-content: center;

        &:hover {
            color: ${token.colorText};
            background: ${token.colorFillQuaternary};
        }

        &:active {
            transform: scale(0.98);
        }
    `,
    navActions: css`
        display: flex;
        align-items: center;
        gap: 8px;
        flex-shrink: 0;
        min-width: 200px;
        justify-content: flex-end;

        @media (max-width: 968px) {
            display: none;
        }
    `,
    navBtn: css`
        height: 36px;
        padding: 0 16px;
        font-size: 14px;
        font-weight: 500;
        border-radius: ${token.borderRadiusLG}px;
        border: none;
        cursor: pointer;
        transition: all 0.2s ease;
        white-space: nowrap;
        display: inline-flex;
        align-items: center;
        gap: 6px;
        position: relative;
        letter-spacing: 0.2px;

        &.btn-secondary {
            background: ${token.colorFillQuaternary};
            color: ${token.colorText};

            &:hover {
                background: ${token.colorFillTertiary};
            }

            &:active {
                transform: scale(0.96);
            }
        }

        &.btn-primary {
            background: ${token.colorPrimary};
            color: ${token.colorWhite};

            &:hover {
                background: ${token.colorPrimaryHover};
            }

            &:active {
                transform: scale(0.96);
            }
        }
    `,
    mobileMenuBtn: css`
        display: none;
        width: 36px;
        height: 36px;
        align-items: center;
        justify-content: center;
        font-size: 18px;
        color: ${token.colorText};
        cursor: pointer;
        border-radius: ${token.borderRadius}px;
        background: transparent;
        border: none;
        transition: all 0.2s ease;
        transform-origin: center;
        backface-visibility: hidden;
        -webkit-font-smoothing: antialiased;

        svg {
            transform-origin: center;
            backface-visibility: hidden;
        }

        &:hover {
            background: ${token.colorFillQuaternary};
        }

        &:active {
            transform: scale(0.95);
        }

        @media (max-width: 968px) {
            display: flex;
        }
    `,
    mobileMenuMask: css`
        position: fixed;
        inset: 0;
        background: ${isDarkMode ? 'rgba(0, 0, 0, 0.65)' : 'rgba(0, 0, 0, 0.45)'};
        backdrop-filter: blur(4px);
        -webkit-backdrop-filter: blur(4px);
        z-index: 9998;
        opacity: 0;
        pointer-events: none;
        transition: opacity 0.3s ease;

        &.open {
            opacity: 1;
            pointer-events: auto;
        }
    `,
    mobileMenuDrawer: css`
        position: fixed;
        top: 0;
        right: 0;
        bottom: 0;
        width: 300px;
        background: ${isDarkMode
                ? 'rgba(18, 18, 18, 0.6)'
                : 'rgba(255, 255, 255, 0.65)'};
        backdrop-filter: blur(28px) saturate(170%);
        -webkit-backdrop-filter: blur(28px) saturate(170%);
        box-shadow: ${isDarkMode
                ? '-2px 0 50px rgba(0, 0, 0, 0.7)'
                : '-2px 0 50px rgba(0, 0, 0, 0.2)'};
        border-left: 1px solid ${isDarkMode
                ? 'rgba(255, 255, 255, 0.12)'
                : 'rgba(0, 0, 0, 0.1)'};
        z-index: 9999;
        display: flex;
        flex-direction: column;
        transform: translateX(100%);
        transition: transform 0.4s cubic-bezier(0.32, 0.72, 0, 1);

        &.open {
            transform: translateX(0);
        }
    `,
    mobileMenuHeader: css`
        padding: 20px 20px 16px;
        display: flex;
        align-items: center;
        justify-content: flex-end;
        flex-shrink: 0;
    `,
    mobileMenuClose: css`
        width: 28px;
        height: 28px;
        border-radius: 50%;
        border: none;
        background: ${isDarkMode
                ? 'rgba(255, 255, 255, 0.1)'
                : 'rgba(0, 0, 0, 0.08)'};
        color: ${token.colorTextSecondary};
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: all 0.2s ease;
        font-size: 13px;
        transform-origin: center;
        backface-visibility: hidden;
        -webkit-font-smoothing: antialiased;

        svg {
            transform-origin: center;
            backface-visibility: hidden;
        }

        &:hover {
            background: ${isDarkMode
                    ? 'rgba(255, 255, 255, 0.15)'
                    : 'rgba(0, 0, 0, 0.12)'};
            color: ${token.colorText};
        }

        &:active {
            transform: scale(0.9);
        }
    `,
    mobileMenuBody: css`
        padding: 0 20px 32px;
        flex: 1;
        display: flex;
        flex-direction: column;
        overflow-y: auto;

        /* 隐藏滚动条但保持功能 */

        &::-webkit-scrollbar {
            width: 0;
            display: none;
        }

        scrollbar-width: none;
    `,
    mobileMenuNav: css`
        flex: 1;
        padding-bottom: 24px;
        margin-bottom: 0;
    `,
    mobileMenuItem: css`
        padding: 15px 8px;
        font-size: 15px;
        color: ${token.colorText};
        cursor: pointer;
        transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
        font-weight: 400;
        margin-bottom: 0;
        background: transparent;
        border: none;
        border-bottom: 0.5px solid ${isDarkMode
                ? 'rgba(255, 255, 255, 0.08)'
                : 'rgba(0, 0, 0, 0.08)'};
        position: relative;
        display: flex;
        align-items: center;
        justify-content: space-between;
        letter-spacing: 0.2px;

        &::after {
            content: '›';
            font-size: 18px;
            color: ${token.colorTextTertiary};
            transition: all 0.2s ease;
            opacity: 0.25;
            font-weight: 300;
        }

        &:hover {
            background: transparent;
            color: ${token.colorPrimary};
            padding-left: 6px;

            &::after {
                opacity: 0.8;
                transform: translateX(-2px);
                color: ${token.colorPrimary};
            }
        }

        &:active {
            opacity: 0.4;
        }

        &:last-child {
            border-bottom: none;
        }
    `,
    mobileThemeSection: css`
        padding: 16px 8px 20px;
        margin-bottom: 16px;
        border-top: 0.5px solid ${isDarkMode
                ? 'rgba(255, 255, 255, 0.08)'
                : 'rgba(0, 0, 0, 0.08)'};
    `,
    mobileThemeLabel: css`
        font-size: 11px;
        color: ${token.colorTextTertiary};
        padding: 0 4px 10px;
        font-weight: 500;
        letter-spacing: 0.3px;
    `,
    mobileThemeButtons: css`
        display: flex;
        gap: 4px;
        padding: 2px;
        background: ${token.colorFillQuaternary};
        border-radius: ${token.borderRadiusSM}px;
    `,
    mobileThemeBtn: css`
        flex: 1;
        height: 28px;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 3px;
        font-size: 11px;
        font-weight: 500;
        border-radius: ${token.borderRadiusSM}px;
        border: none;
        background: transparent;
        color: ${token.colorTextSecondary};
        cursor: pointer;
        transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);

        &.active {
            background: ${token.colorBgContainer};
            color: ${token.colorText};
            box-shadow: ${token.boxShadowTertiary};
        }

        &:hover:not(.active) {
            color: ${token.colorText};
        }

        &:active {
            transform: scale(0.95);
        }

        .anticon {
            font-size: 13px;
            transform-origin: center;
            backface-visibility: hidden;
            -webkit-font-smoothing: antialiased;
        }
    `,
    mobileActions: css`
        margin-top: auto;
        padding: 20px 8px 0;
        border-top: 0.5px solid ${isDarkMode
                ? 'rgba(255, 255, 255, 0.08)'
                : 'rgba(0, 0, 0, 0.08)'};
        display: flex;
        flex-direction: column;
        gap: 12px;
    `,
    mobileActionItem: css`
        padding: 0 20px;
        height: 50px;
        font-size: 15px;
        font-weight: 500;
        border-radius: ${token.borderRadiusLG}px;
        border: 1px solid ${token.colorBorder};
        cursor: pointer;
        transition: all 0.2s ease;
        background: ${token.colorBgContainer};
        color: ${token.colorText};
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 8px;
        letter-spacing: 0.2px;
        width: 100%;

        .anticon {
            font-size: 16px;
        }

        &:hover {
            background: ${token.colorBgTextHover};
            border-color: ${token.colorPrimary};
            color: ${token.colorPrimary};
        }

        &:active {
            transform: scale(0.98);
        }
    `,
    mobileBtnPrimary: css`
        background: ${token.colorPrimary};
        border-color: transparent;
        color: ${token.colorWhite} !important;
        font-weight: 600;

        .anticon {
            color: ${token.colorWhite} !important;
        }

        &:hover {
            background: ${token.colorPrimaryHover};
            border-color: transparent;
            color: ${token.colorWhite} !important;
        }
    `,
    heroSection: css`
        position: relative;
        width: 100%;
        min-height: 100vh;
        height: 100vh;
        display: flex;
        align-items: center;
        justify-content: center;
        margin: 0;
        padding: 0;
        background: ${isDarkMode
                ? `radial-gradient(ellipse 80% 50% at 50% 50%, ${token.colorBgLayout} 0%, ${token.colorBgContainer} 100%)`
                : `radial-gradient(ellipse 80% 50% at 50% 50%, ${token.colorBgContainer} 0%, ${token.colorBgLayout} 100%)`};
        overflow: hidden;

        @media (max-width: 768px) {
            min-height: 100vh;
            height: 100vh;
        }

        &::before {
            content: '';
            position: absolute;
            width: 1200px;
            height: 1200px;
            top: -600px;
            left: 50%;
            transform: translateX(-50%);
            background: radial-gradient(circle, ${token.colorPrimary}${isDarkMode ? '18' : '12'} 0%, transparent 60%);
            border-radius: 50%;
            opacity: 0.7;
        }

        &::after {
            content: '';
            position: absolute;
            inset: 0;
            background: ${isDarkMode
                    ? `radial-gradient(circle at 50% 30%, ${token.colorPrimary}08, transparent 50%)`
                    : `radial-gradient(circle at 50% 30%, ${token.colorPrimary}12, transparent 50%)`};
            opacity: 0.6;
        }
    `,
    heroSpotlight: css`
        position: absolute;
        inset: 0;
        background: radial-gradient(
                ellipse 50% 35% at 50% 40%,
                ${token.colorPrimary}${isDarkMode ? '08' : '06'} 0%,
                transparent 65%
        );
        opacity: 0.9;
    `,
    heroGrid: css`
        position: absolute;
        inset: 0;
        background-image: linear-gradient(${isDarkMode ? `${token.colorPrimary}05` : `${token.colorPrimary}03`} 1px, transparent 1px),
        linear-gradient(90deg, ${isDarkMode ? `${token.colorPrimary}05` : `${token.colorPrimary}03`} 1px, transparent 1px);
        background-size: 80px 80px;
        mask-image: radial-gradient(ellipse 80% 60% at 50% 40%, black 0%, transparent 70%);
        -webkit-mask-image: radial-gradient(ellipse 80% 60% at 50% 40%, black 0%, transparent 70%);
        animation: gridMove 40s linear infinite;
        opacity: ${isDarkMode ? '0.4' : '0.3'};

        @keyframes gridMove {
            0% {
                transform: translate(0, 0);
            }
            100% {
                transform: translate(80px, 80px);
            }
        }
    `,
    heroScanLine: css`
        position: absolute;
        inset: 0;
        overflow: hidden;
        pointer-events: none;

        &::before {
            content: '';
            position: absolute;
            left: 0;
            right: 0;
            height: 2px;
            background: linear-gradient(
                    90deg,
                    transparent 0%,
                    ${isDarkMode ? `${token.colorPrimary}80` : `${token.colorPrimary}60`} 50%,
                    transparent 100%
            );
            box-shadow: 0 0 20px ${isDarkMode ? `${token.colorPrimary}60` : `${token.colorPrimary}40`};
            animation: scanVertical 8s ease-in-out infinite;
        }

        &::after {
            content: '';
            position: absolute;
            top: 0;
            bottom: 0;
            width: 2px;
            background: linear-gradient(
                    180deg,
                    transparent 0%,
                    ${isDarkMode ? `${token.colorPrimary}80` : `${token.colorPrimary}60`} 50%,
                    transparent 100%
            );
            box-shadow: 0 0 20px ${isDarkMode ? `${token.colorPrimary}60` : `${token.colorPrimary}40`};
            animation: scanHorizontal 10s ease-in-out infinite;
        }

        @keyframes scanVertical {
            0%, 100% {
                top: -10%;
                opacity: 0;
            }
            10% {
                opacity: 1;
            }
            90% {
                opacity: 1;
            }
            100% {
                top: 110%;
                opacity: 0;
            }
        }

        @keyframes scanHorizontal {
            0%, 100% {
                left: -10%;
                opacity: 0;
            }
            10% {
                opacity: 1;
            }
            90% {
                opacity: 1;
            }
            100% {
                left: 110%;
                opacity: 0;
            }
        }
    `,
    heroCircuits: css`
        position: absolute;
        inset: 0;
        pointer-events: none;
        overflow: hidden;
        opacity: ${isDarkMode ? '0.15' : '0.1'};

        &::before,
        &::after {
            content: '';
            position: absolute;
            border: 1px solid ${token.colorPrimary};
            border-radius: 50%;
        }

        &::before {
            width: 600px;
            height: 600px;
            top: 10%;
            left: -10%;
            animation: circuitPulse 15s ease-in-out infinite;
        }

        &::after {
            width: 800px;
            height: 800px;
            bottom: -20%;
            right: -15%;
            animation: circuitPulse 20s ease-in-out infinite reverse;
        }

        @keyframes circuitPulse {
            0%, 100% {
                transform: scale(1);
                opacity: 0.3;
            }
            50% {
                transform: scale(1.1);
                opacity: 0.6;
            }
        }
    `,
    heroDots: css`
        position: absolute;
        inset: 0;
        background-image: radial-gradient(circle at 15% 25%, ${isDarkMode ? `${token.colorText}40` : `${token.colorPrimary}50`} 2px, transparent 2px),
        radial-gradient(circle at 85% 30%, ${isDarkMode ? `${token.colorText}35` : `${token.colorPrimary}45`} 1.5px, transparent 1.5px),
        radial-gradient(circle at 40% 75%, ${isDarkMode ? `${token.colorText}30` : `${token.colorPrimary}40`} 1.5px, transparent 1.5px),
        radial-gradient(circle at 70% 80%, ${isDarkMode ? `${token.colorText}45` : `${token.colorPrimary}55`} 2px, transparent 2px),
        radial-gradient(circle at 25% 60%, ${isDarkMode ? `${token.colorText}35` : `${token.colorPrimary}45`} 1.5px, transparent 1.5px);
        background-size: 100% 100%;
        animation: dotsGlow 4s ease-in-out infinite;
        pointer-events: none;

        @keyframes dotsGlow {
            0%, 100% {
                opacity: 0.4;
            }
            50% {
                opacity: 0.8;
            }
        }
    `,
    heroDataStream: css`
        position: absolute;
        inset: 0;
        pointer-events: none;
        overflow: hidden;

        &::before {
            content: '';
            position: absolute;
            width: 1px;
            height: 100px;
            left: 20%;
            background: linear-gradient(
                    180deg,
                    transparent 0%,
                    ${isDarkMode ? `${token.colorPrimary}80` : `${token.colorPrimary}60`} 50%,
                    transparent 100%
            );
            box-shadow: 0 0 10px ${isDarkMode ? `${token.colorPrimary}80` : `${token.colorPrimary}60`};
            animation: dataFlow1 5s linear infinite;
        }

        &::after {
            content: '';
            position: absolute;
            width: 1px;
            height: 80px;
            left: 75%;
            background: linear-gradient(
                    180deg,
                    transparent 0%,
                    ${isDarkMode ? `${token.colorPrimary}60` : `${token.colorPrimary}50`} 50%,
                    transparent 100%
            );
            box-shadow: 0 0 8px ${isDarkMode ? `${token.colorPrimary}60` : `${token.colorPrimary}50`};
            animation: dataFlow2 6s linear infinite 2s;
        }

        @keyframes dataFlow1 {
            0% {
                top: -100px;
                opacity: 0;
            }
            10% {
                opacity: 1;
            }
            90% {
                opacity: 1;
            }
            100% {
                top: 100%;
                opacity: 0;
            }
        }

        @keyframes dataFlow2 {
            0% {
                top: -80px;
                opacity: 0;
            }
            10% {
                opacity: 1;
            }
            90% {
                opacity: 1;
            }
            100% {
                top: 100%;
                opacity: 0;
            }
        }
    `,
    heroTechLines: css`
        position: absolute;
        inset: 0;
        pointer-events: none;

        &::before {
            content: '';
            position: absolute;
            left: 10%;
            top: 30%;
            width: 150px;
            height: 1px;
            background: linear-gradient(
                    90deg,
                    transparent 0%,
                    ${isDarkMode ? `${token.colorPrimary}50` : `${token.colorPrimary}40`} 50%,
                    transparent 100%
            );
            animation: techLineExpand 3s ease-in-out infinite;
        }

        &::after {
            content: '';
            position: absolute;
            right: 15%;
            bottom: 35%;
            width: 200px;
            height: 1px;
            background: linear-gradient(
                    90deg,
                    transparent 0%,
                    ${isDarkMode ? `${token.colorPrimary}40` : `${token.colorPrimary}30`} 50%,
                    transparent 100%
            );
            animation: techLineExpand 4s ease-in-out infinite 1.5s;
        }

        @keyframes techLineExpand {
            0%, 100% {
                transform: scaleX(0.5);
                opacity: 0.3;
            }
            50% {
                transform: scaleX(1);
                opacity: 0.8;
            }
        }
    `,
    heroParticles: css`
        position: absolute;
        inset: 0;
        overflow: hidden;
        pointer-events: none;
    `,
    particle: css`
        position: absolute;
        width: 3px;
        height: 3px;
        background: ${token.colorPrimary};
        border-radius: 50%;
        box-shadow: 0 0 4px ${token.colorPrimary};
        animation: particleFloat var(--duration) cubic-bezier(0.4, 0, 0.6, 1) infinite;

        @keyframes particleFloat {
            0% {
                opacity: 0;
                transform: translate3d(0, 0, 0) scale(0.6);
            }
            15% {
                opacity: 0.6;
            }
            50% {
                opacity: 0.8;
                transform: translate3d(calc(var(--tx) / 2), calc(var(--ty) / 2), 0) scale(1);
            }
            85% {
                opacity: 0.6;
            }
            100% {
                opacity: 0;
                transform: translate3d(var(--tx), var(--ty), 0) scale(0.6);
            }
        }
    `,
    heroOrbs: css`
        position: absolute;
        inset: 0;
        overflow: hidden;
        pointer-events: none;
    `,
    orb: css`
        position: absolute;
        border-radius: 50%;
        background: radial-gradient(circle at 30% 30%,
        ${token.colorPrimary}${isDarkMode ? '20' : '15'},
        ${token.colorPrimary}${isDarkMode ? '10' : '08'},
        transparent
        );
        filter: blur(30px);
        animation: orbFloat var(--duration) ease-in-out infinite;
        animation-delay: var(--delay);

        @keyframes orbFloat {
            0%, 100% {
                transform: translate3d(0, 0, 0) scale(1);
                opacity: ${isDarkMode ? '0.3' : '0.25'};
            }
            50% {
                transform: translate3d(var(--tx), var(--ty), 0) scale(1.1);
                opacity: ${isDarkMode ? '0.5' : '0.45'};
            }
        }
    `,
    heroGradient: css`
        position: absolute;
        inset: 0;
        background: radial-gradient(circle at 30% 20%, ${token.colorPrimary}${isDarkMode ? '10' : '08'} 0%, transparent 45%),
        radial-gradient(circle at 70% 80%, ${token.colorPrimaryHover}${isDarkMode ? '08' : '06'} 0%, transparent 45%);
        opacity: 0.6;
    `,
    heroGlowOrbs: css`
        position: absolute;
        inset: 0;
        pointer-events: none;
        overflow: hidden;
    `,
    glowOrb: css`
        position: absolute;
        border-radius: 50%;
        filter: blur(25px);
        background: radial-gradient(circle at center,
        ${token.colorPrimary}50 0%,
        ${token.colorPrimary}25 40%,
        transparent 70%
        );
        animation: glowFloat var(--duration) ease-in-out infinite;
        animation-delay: var(--delay);

        @keyframes glowFloat {
            0%, 100% {
                transform: translate(0, 0) scale(1);
                opacity: 0.4;
            }
            50% {
                transform: translate(var(--tx), var(--ty)) scale(1.2);
                opacity: 0.7;
            }
        }
    `,
    heroContent: css`
        position: relative;
        z-index: 1;
        text-align: center;
        max-width: 1100px;
        width: 100%;
        margin: 0 auto;
        padding: 0 20px;
        transform: translateY(-32px);

        @media (max-width: 768px) {
            padding: 0 16px;
            transform: translateY(-20px);
        }
    `,
    heroTitle: css`
        font-size: 80px !important;
        font-weight: 900 !important;
        margin-bottom: 32px !important;
        line-height: 1.05 !important;
        color: ${token.colorText} !important;
        letter-spacing: -0.04em;
        animation: titleFadeIn 1.2s cubic-bezier(0.16, 1, 0.3, 1) forwards;
        will-change: opacity, transform;

        @keyframes titleFadeIn {
            from {
                opacity: 0;
                transform: translateY(40px) scale(0.95);
            }
            to {
                opacity: 1;
                transform: translateY(0) scale(1);
            }
        }

        @media (max-width: 768px) {
            font-size: 48px !important;
            margin-bottom: 20px !important;
            letter-spacing: -0.02em;
        }
    `,
    heroSubtitle: css`
        font-size: 28px !important;
        color: ${token.colorText} !important;
        margin-bottom: 56px !important;
        font-weight: 600;
        opacity: 0;
        line-height: 1.35 !important;
        letter-spacing: -0.015em;
        animation: subtitleFadeIn 1.2s cubic-bezier(0.16, 1, 0.3, 1) 0.3s forwards;
        text-shadow: 0 2px 12px ${isDarkMode ? 'rgba(0, 0, 0, 0.3)' : 'rgba(255, 255, 255, 0.8)'};

        @keyframes subtitleFadeIn {
            from {
                opacity: 0;
                transform: translateY(30px);
                filter: blur(10px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
                filter: blur(0);
            }
        }

        @media (max-width: 768px) {
            font-size: 20px !important;
            margin-bottom: 44px !important;
        }
    `,
    ctaButtons: css`
        opacity: 0;
        animation: buttonFadeIn 0.9s cubic-bezier(0.16, 1, 0.3, 1) 0.5s forwards;
        display: flex;
        gap: 16px;
        justify-content: center;
        align-items: center;

        @keyframes buttonFadeIn {
            from {
                opacity: 0;
                transform: translateY(20px) scale(0.95);
            }
            to {
                opacity: 1;
                transform: translateY(0) scale(1);
            }
        }

        @media (max-width: 768px) {
            flex-direction: column;
            gap: 12px;
            width: 100%;
        }
    `,
    btnPrimary: css`
        height: 56px;
        padding: 0 40px;
        font-size: 16px;
        font-weight: 600;
        border-radius: ${token.borderRadiusLG}px;
        border: none;
        background: ${token.colorPrimary};
        color: ${token.colorWhite};
        cursor: pointer;
        transition: all 0.2s ease;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        gap: 8px;
        letter-spacing: 0.3px;

        &:hover {
            background: ${token.colorPrimaryHover};
            transform: translateY(-2px);
            box-shadow: 0 8px 20px ${token.colorPrimary}50;
        }

        &:active {
            transform: translateY(0);
        }

        @media (max-width: 768px) {
            height: 52px;
            padding: 0 36px;
            font-size: 15px;
            width: 100%;
            max-width: 320px;
        }
    `,
    btnSecondary: css`
        height: 56px;
        padding: 0 40px;
        font-size: 16px;
        font-weight: 600;
        border-radius: ${token.borderRadiusLG}px;
        border: 1px solid ${token.colorBorder};
        background: ${token.colorBgContainer};
        color: ${token.colorText};
        cursor: pointer;
        transition: all 0.2s ease;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        gap: 8px;
        letter-spacing: 0.3px;

        &:hover {
            background: ${token.colorBgTextHover};
            border-color: ${token.colorPrimary};
            color: ${token.colorPrimary};
            transform: translateY(-2px);
            box-shadow: 0 8px 20px ${token.colorFillQuaternary};
        }

        &:active {
            transform: translateY(0);
        }

        @media (max-width: 768px) {
            height: 52px;
            padding: 0 36px;
            font-size: 15px;
            width: 100%;
            max-width: 320px;
        }
    `,
    section: css`
        max-width: 1400px;
        margin: 0 auto;
        padding: 80px 20px;
        width: 100%;

        @media (max-width: 1200px) {
            padding: 70px 16px;
        }

        @media (max-width: 768px) {
            padding: 60px 16px;
        }
    `,
    sectionTitle: css`
        text-align: center;
        font-size: 48px !important;
        font-weight: 800 !important;
        margin-bottom: 16px !important;
        background: ${isDarkMode
                ? `linear-gradient(135deg, ${token.colorText} 0%, ${token.colorTextSecondary} 50%, ${token.colorTextTertiary} 100%)`
                : `linear-gradient(135deg, ${token.colorText} 0%, ${token.colorTextSecondary} 50%, ${token.colorTextTertiary} 100%)`};
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
        letter-spacing: -0.02em;
        line-height: 1.2 !important;

        @media (max-width: 768px) {
            font-size: 32px !important;
            margin-bottom: 12px !important;
            letter-spacing: -0.01em;
        }
    `,
    sectionSubtitle: css`
        text-align: center;
        font-size: 16px !important;
        color: ${token.colorTextSecondary} !important;
        margin-bottom: 48px !important;
        max-width: 700px;
        margin-left: auto !important;
        margin-right: auto !important;
        font-weight: 400;
        line-height: 1.6 !important;

        @media (max-width: 768px) {
            font-size: 14px !important;
            margin-bottom: 40px !important;
        }
    `,
    statsSection: css`
        background: ${isDarkMode
                ? `linear-gradient(135deg, ${token.colorBgElevated} 0%, ${token.colorBgLayout} 100%)`
                : `linear-gradient(135deg, ${token.colorBgLayout} 0%, ${token.colorBgContainer} 100%)`};
        padding: 60px 20px;
        margin: 0;
        width: 100%;
        position: relative;
        overflow: hidden;

        &::before {
            content: '';
            position: absolute;
            inset: 0;
            background: ${isDarkMode
                    ? `linear-gradient(135deg, ${token.colorPrimary}15 0%, transparent 50%, ${token.colorPrimaryHover}10 100%)`
                    : `linear-gradient(135deg, ${token.colorPrimary}05 0%, transparent 50%, ${token.colorPrimaryHover}03 100%)`};
        }

        @media (max-width: 1200px) {
            padding: 50px 16px;
        }

        @media (max-width: 768px) {
            padding: 40px 16px;
        }
    `,
    statsContainer: css`
        max-width: 1200px;
        margin: 0 auto;
        position: relative;
        z-index: 1;
    `,
    statCard: css`
        text-align: center;
        padding: 32px 16px;
        background: ${isDarkMode
                ? 'rgba(255, 255, 255, 0.03)'
                : 'rgba(255, 255, 255, 0.8)'};
        border-radius: ${token.borderRadiusLG}px;
        border: 1px solid ${isDarkMode
                ? 'rgba(255, 255, 255, 0.1)'
                : 'rgba(0, 0, 0, 0.06)'};
        backdrop-filter: blur(10px);
        -webkit-backdrop-filter: blur(10px);
        transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
        position: relative;
        overflow: hidden;

        &::before {
            content: '';
            position: absolute;
            inset: 0;
            background: ${token.colorPrimary};
            opacity: 0;
            transition: opacity 0.4s ease;
        }

        &:hover {
            transform: translateY(-8px) scale(1.02);
            border-color: ${token.colorPrimary};
            box-shadow: ${isDarkMode
                    ? `0 12px 32px ${token.colorPrimary}60`
                    : `0 12px 32px ${token.colorPrimary}40`};

            &::before {
                opacity: 0.05;
            }

            .ant-statistic-content {
                transform: scale(1.05);
            }
        }

        .ant-statistic-title {
            color: ${token.colorText} !important;
            font-size: 12px !important;
            margin-bottom: 8px;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.3px;
            position: relative;
            z-index: 1;
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
        }

        .ant-statistic-content {
            color: ${token.colorPrimary} !important;
            font-size: 36px !important;
            font-weight: 800 !important;
            transition: transform 0.4s ease;
            position: relative;
            z-index: 1;
            word-break: break-all;
        }

        .ant-statistic-content-value {
            color: ${token.colorPrimary} !important;
        }

        @media (max-width: 768px) {
            padding: 20px 8px;

            .ant-statistic-title {
                font-size: 10px !important;
                margin-bottom: 6px;
                letter-spacing: 0.2px;
                line-height: 1.2 !important;
            }

            .ant-statistic-content {
                font-size: 24px !important;
            }
        }
    `,
    featureCard: css`
        height: 100%;
        padding: 40px 32px;
        border-radius: ${token.borderRadiusLG * 2.5}px;
        transition: all 0.6s cubic-bezier(0.34, 1.56, 0.64, 1);
        border: 1px solid ${isDarkMode ? `${token.colorPrimary}15` : token.colorBorderSecondary};
        background: ${isDarkMode
                ? `linear-gradient(145deg, ${token.colorBgContainer}f5 0%, rgba(10, 10, 10, 0.8) 100%)`
                : `linear-gradient(145deg, ${token.colorBgContainer} 0%, rgba(248, 248, 248, 0.95) 100%)`};
        position: relative;
        overflow: hidden;
        backdrop-filter: blur(12px);
        box-shadow: ${isDarkMode
                ? `0 6px 20px rgba(0, 0, 0, 0.35), inset 0 1px 0 rgba(255, 255, 255, 0.08)`
                : `0 4px 12px rgba(0, 0, 0, 0.06), inset 0 1px 0 rgba(255, 255, 255, 0.8)`};

        &::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            height: 120px;
            background: ${isDarkMode
                    ? `radial-gradient(ellipse 100% 80% at 50% 0%, ${token.colorPrimary}12 0%, transparent 70%)`
                    : `radial-gradient(ellipse 100% 80% at 50% 0%, ${token.colorPrimary}08 0%, transparent 70%)`};
            opacity: 0;
            transition: opacity 0.6s ease;
            z-index: 0;
            pointer-events: none;
        }

        > * {
            position: relative;
            z-index: 1;
        }

        &:hover {
            transform: translateY(-16px) scale(1.03);
            box-shadow: ${isDarkMode
                    ? `0 24px 70px rgba(0, 0, 0, 0.7), 0 0 0 1px ${token.colorPrimary}60, 0 0 70px ${token.colorPrimary}25, inset 0 2px 0 rgba(255, 255, 255, 0.12)`
                    : `0 24px 70px rgba(0, 0, 0, 0.15), 0 0 0 1px ${token.colorPrimary}40, 0 0 50px ${token.colorPrimary}18, inset 0 2px 0 rgba(255, 255, 255, 1)`};
            border-color: ${token.colorPrimary};

            &::before {
                opacity: 1;
            }

            > div:first-child {
                transform: translateY(-12px) scale(1.15) rotate(8deg);
                box-shadow: ${isDarkMode
                        ? `0 24px 60px ${token.colorPrimary}60, 0 0 40px ${token.colorPrimary}40, inset 0 2px 0 rgba(255,255,255,0.2)`
                        : `0 20px 50px ${token.colorPrimary}35, 0 0 30px ${token.colorPrimary}25, inset 0 2px 0 rgba(255,255,255,1)`};
                border-color: ${token.colorPrimary};
                background: ${isDarkMode
                        ? `linear-gradient(135deg, ${token.colorPrimary}40 0%, ${token.colorPrimaryHover}20 50%, ${token.colorPrimary}15 100%)`
                        : `linear-gradient(135deg, ${token.colorPrimary}20 0%, ${token.colorPrimaryHover}12 50%, ${token.colorPrimary}08 100%)`};

                &::before {
                    opacity: 1;
                    animation-play-state: running;
                }

                img {
                    transform: scale(1.25) rotate(-10deg);
                    filter: ${isDarkMode
                            ? `brightness(1.6) contrast(1.25) saturate(1.4) drop-shadow(0 6px 16px ${token.colorPrimary})`
                            : `brightness(1.15) contrast(1.12) saturate(1.3) drop-shadow(0 4px 12px ${token.colorPrimary}70)`};
                }
            }

            h4 {
                color: ${token.colorPrimary};
                transform: translateY(-2px);
            }
        }

        @media (max-width: 768px) {
            padding: 32px 26px;
            border-radius: ${token.borderRadiusLG * 2}px;
        }
    `,
    featureIcon: css`
        width: 80px;
        height: 80px;
        display: flex;
        align-items: center;
        justify-content: center;
        background: ${isDarkMode
                ? `linear-gradient(135deg, ${token.colorPrimary}25 0%, ${token.colorPrimaryHover}12 50%, ${token.colorPrimary}08 100%)`
                : `linear-gradient(135deg, ${token.colorPrimary}12 0%, ${token.colorPrimaryHover}08 50%, ${token.colorPrimary}05 100%)`};
        border-radius: ${token.borderRadiusLG * 2}px;
        margin: 0 auto 24px;
        transition: all 0.6s cubic-bezier(0.34, 1.56, 0.64, 1);
        box-shadow: ${isDarkMode
                ? `0 10px 30px ${token.colorPrimary}30, inset 0 1px 0 rgba(255,255,255,0.1)`
                : `0 6px 18px ${token.colorPrimary}20, inset 0 1px 0 rgba(255,255,255,0.8)`};
        border: 1px solid ${isDarkMode ? `${token.colorPrimary}25` : `${token.colorPrimary}15`};
        position: relative;
        overflow: hidden;

        &::before {
            content: '';
            position: absolute;
            inset: -50%;
            background: ${isDarkMode
                    ? `conic-gradient(from 180deg, transparent, ${token.colorPrimary}15, transparent)`
                    : `conic-gradient(from 180deg, transparent, ${token.colorPrimary}10, transparent)`};
            opacity: 0;
            transition: all 0.6s ease;
            animation: iconRotate 4s linear infinite;
            animation-play-state: paused;
        }

        @keyframes iconRotate {
            0% {
                transform: rotate(0deg);
            }
            100% {
                transform: rotate(360deg);
            }
        }

        img {
            width: 40px;
            height: 40px;
            object-fit: contain;
            filter: ${isDarkMode
                    ? `brightness(1.3) contrast(1.1) drop-shadow(0 2px 4px ${token.colorPrimary}40)`
                    : `brightness(0.95) contrast(1.05) drop-shadow(0 1px 2px rgba(0,0,0,0.1))`};
            transition: all 0.6s cubic-bezier(0.34, 1.56, 0.64, 1);
        }


        @media (max-width: 768px) {
            width: 64px;
            height: 64px;
            margin-bottom: 18px;

            img {
                width: 32px;
                height: 32px;
            }
        }
    `,
    featureTitle: css`
        font-size: 18px !important;
        font-weight: 600 !important;
        margin-bottom: 10px !important;
        text-align: center;
        color: ${token.colorText};
        transition: all 0.3s ease;

        @media (max-width: 768px) {
            font-size: 16px !important;
            margin-bottom: 8px !important;
        }
    `,
    featureDescription: css`
        color: ${token.colorTextSecondary} !important;
        line-height: 1.6 !important;
        text-align: center;
        margin-bottom: 0 !important;
        font-size: 14px;

        @media (max-width: 768px) {
            font-size: 13px !important;
        }
    `,
    productCard: css`
        height: 100%;
        padding: 36px 30px;
        border-radius: ${token.borderRadiusLG * 2.2}px;
        transition: all 0.6s cubic-bezier(0.34, 1.56, 0.64, 1);
        border: 1px solid ${isDarkMode ? `${token.colorPrimary}12` : token.colorBorderSecondary};
        background: ${isDarkMode
                ? `linear-gradient(145deg, ${token.colorBgContainer}f8 0%, rgba(15, 15, 15, 0.85) 100%)`
                : `linear-gradient(145deg, ${token.colorBgContainer} 0%, rgba(252, 252, 252, 0.9) 100%)`};
        position: relative;
        overflow: hidden;
        backdrop-filter: blur(10px);
        box-shadow: ${isDarkMode
                ? `0 8px 24px rgba(0, 0, 0, 0.4), inset 0 1px 0 rgba(255, 255, 255, 0.1)`
                : `0 6px 16px rgba(0, 0, 0, 0.06), inset 0 1px 0 rgba(255, 255, 255, 0.9)`};

        &::before {
            content: '';
            position: absolute;
            top: -2px;
            left: -2px;
            right: -2px;
            height: 140px;
            background: ${isDarkMode
                    ? `radial-gradient(ellipse 120% 100% at 50% 0%, ${token.colorPrimary}15 0%, transparent 65%)`
                    : `radial-gradient(ellipse 120% 100% at 50% 0%, ${token.colorPrimary}10 0%, transparent 65%)`};
            opacity: 0;
            transition: opacity 0.6s ease;
            z-index: 0;
            pointer-events: none;
        }

        > *:not(.product-icon-bg) {
            position: relative;
            z-index: 1;
        }

        &:hover {
            transform: translateY(-14px) scale(1.03);
            box-shadow: ${isDarkMode
                    ? `0 22px 60px rgba(0, 0, 0, 0.6), 0 0 0 1px ${token.colorPrimary}55, 0 0 60px ${token.colorPrimary}22, inset 0 2px 0 rgba(255, 255, 255, 0.15)`
                    : `0 22px 60px rgba(0, 0, 0, 0.12), 0 0 0 1px ${token.colorPrimary}35, 0 0 45px ${token.colorPrimary}15, inset 0 2px 0 rgba(255, 255, 255, 1)`};
            border-color: ${token.colorPrimary};

            &::before {
                opacity: 1;
            }

            .product-icon-bg {
                transform: scale(1.08) rotate(-3deg);
                opacity: ${isDarkMode ? '0.25' : '0.22'};
            }

            h4 {
                color: ${token.colorPrimary};
            }
        }

        @media (max-width: 768px) {
            padding: 28px 24px;
            border-radius: ${token.borderRadiusLG * 1.8}px;
        }
    `,
    productIconBg: css`
        position: absolute;
        top: -20px;
        right: -20px;
        width: 160px;
        height: 160px;
        opacity: ${isDarkMode ? '0.18' : '0.15'};
        transition: all 0.6s cubic-bezier(0.34, 1.56, 0.64, 1);
        z-index: 0;
        pointer-events: none;

        img {
            width: 100%;
            height: 100%;
            object-fit: contain;
        }

        @media (max-width: 768px) {
            width: 120px;
            height: 120px;
            top: -15px;
            right: -15px;
        }
    `,
    productTitle: css`
        font-size: 20px !important;
        font-weight: 700 !important;
        margin-bottom: 12px !important;
        color: ${token.colorText};
        transition: all 0.3s ease;
        position: relative;
        z-index: 1;

        @media (max-width: 768px) {
            font-size: 18px !important;
            margin-bottom: 10px !important;
        }
    `,
    productDesc: css`
        color: ${token.colorTextSecondary} !important;
        margin-bottom: 20px !important;
        min-height: 52px;
        font-size: 15px;
        line-height: 1.7;
        font-weight: 400;
        position: relative;
        z-index: 1;

        @media (max-width: 768px) {
            font-size: 14px !important;
            min-height: auto;
            margin-bottom: 16px !important;
        }
    `,
    productFeatures: css`
        margin-top: 20px;
        padding-top: 20px;
        border-top: 1px solid ${token.colorBorderSecondary};
        position: relative;
        z-index: 1;

        @media (max-width: 768px) {
            margin-top: 16px;
            padding-top: 16px;
        }
    `,
    productFeatureItem: css`
        display: flex;
        align-items: center;
        margin-bottom: 8px;
        color: ${token.colorTextSecondary};
        font-size: 13px;
        padding: 6px 10px;
        border-radius: ${token.borderRadius}px;
        background: ${isDarkMode ? 'rgba(255, 255, 255, 0.02)' : 'rgba(0, 0, 0, 0.02)'};
        transition: all 0.3s ease;

        &:last-child {
            margin-bottom: 0;
        }

        &:hover {
            background: ${isDarkMode ? `${token.colorPrimary}10` : `${token.colorPrimary}08`};
            color: ${token.colorText};
            transform: translateX(4px);
        }

        .anticon {
            color: ${token.colorSuccess};
            margin-right: 8px;
            font-size: 13px;
            transition: all 0.3s ease;
        }

        &:hover .anticon {
            color: ${token.colorPrimary};
            transform: scale(1.15);
        }

        @media (max-width: 768px) {
            font-size: 12px;
            margin-bottom: 6px;
            padding: 5px 8px;

            .anticon {
                font-size: 12px;
                margin-right: 6px;
            }
        }
    `,
    testimonialCard: css`
        height: 100%;
        padding: 28px 24px;
        border-radius: ${token.borderRadiusLG}px;
        border: 1px solid ${token.colorBorderSecondary};
        background: ${token.colorBgContainer};
        transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
        position: relative;
        overflow: hidden;

        &::before {
            content: '';
            position: absolute;
            top: -2px;
            left: -2px;
            right: -2px;
            bottom: -2px;
            background: linear-gradient(135deg, ${token.colorPrimary}15, ${token.colorPrimaryHover}10);
            border-radius: ${token.borderRadiusLG}px;
            opacity: 0;
            transition: opacity 0.4s ease;
            z-index: 0;
            pointer-events: none;
        }

        > * {
            position: relative;
            z-index: 1;
        }

        &:hover {
            transform: translateY(-8px) scale(1.02);
            box-shadow: ${isDarkMode
                    ? `0 16px 32px rgba(0, 0, 0, 0.5), 0 0 40px ${token.colorPrimary}25`
                    : `0 16px 32px rgba(0, 0, 0, 0.12), 0 0 30px ${token.colorPrimary}20`};

            &::before {
                opacity: 1;
            }

            > div:first-child .anticon {
                transform: scale(1.2) rotate(15deg);
            }

            > div:last-child > div:first-child {
                transform: scale(1.1);
                box-shadow: ${isDarkMode
                        ? `0 6px 16px ${token.colorPrimary}50`
                        : `0 4px 12px ${token.colorPrimary}30`};
            }
        }

        @media (max-width: 768px) {
            padding: 24px 20px;
        }
    `,
    testimonialStars: css`
        margin-bottom: 16px;
        display: flex;
        gap: 4px;
        justify-content: center;

        .anticon {
            color: ${token.colorWarning};
            font-size: 16px;
            transition: all 0.3s ease;
        }
    `,
    testimonialContent: css`
        font-size: 15px;
        color: ${token.colorText};
        line-height: 1.8;
        margin-bottom: 24px;
        position: relative;
        padding: 16px;
        background: ${isDarkMode
                ? 'rgba(255, 255, 255, 0.02)'
                : 'rgba(0, 0, 0, 0.02)'};
        border-radius: ${token.borderRadius}px;
        border-left: 3px solid ${token.colorPrimary};

        &::before {
            content: '"';
            position: absolute;
            left: 8px;
            top: 8px;
            font-size: 32px;
            color: ${token.colorPrimary};
            opacity: 0.2;
            font-family: Georgia, serif;
            line-height: 1;
        }

        @media (max-width: 768px) {
            font-size: 14px;
            margin-bottom: 20px;
            padding: 14px;
        }
    `,
    testimonialAuthor: css`
        display: flex;
        align-items: center;
        gap: 12px;
    `,
    testimonialAvatar: css`
        width: 48px;
        height: 48px;
        border-radius: 50%;
        background: ${isDarkMode
                ? `linear-gradient(135deg, ${token.colorPrimary} 0%, ${token.colorPrimaryHover} 100%)`
                : token.colorPrimaryBg};
        display: flex;
        align-items: center;
        justify-content: center;
        color: ${isDarkMode ? '#fff' : token.colorPrimary};
        font-weight: 700;
        font-size: 18px;
        flex-shrink: 0;
        box-shadow: ${isDarkMode
                ? `0 4px 12px ${token.colorPrimary}40`
                : `0 2px 8px ${token.colorPrimary}20`};
        transition: all 0.3s ease;
    `,
    testimonialInfo: css`
        flex: 1;
    `,
    testimonialName: css`
        font-weight: 600;
        margin-bottom: 2px;
        color: ${token.colorText};
        font-size: 14px;
    `,
    testimonialRole: css`
        color: ${token.colorTextTertiary};
        font-size: 12px;
    `,
    advantageSection: css`
        background: ${isDarkMode ? token.colorBgLayout : token.colorBgContainer};
        padding: 60px 20px;
        width: 100%;
        margin: 0;

        @media (max-width: 1200px) {
            padding: 50px 16px;
        }

        @media (max-width: 768px) {
            padding: 40px 16px;
        }
    `,
    advantageCard: css`
        text-align: center;
        padding: 40px 24px;
        border-radius: ${token.borderRadiusLG}px;
        background: ${token.colorBgContainer};
        transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
        height: 100%;
        border: 1px solid ${token.colorBorderSecondary};
        position: relative;
        overflow: hidden;

        &::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            height: 3px;
            background: linear-gradient(90deg, ${token.colorPrimary}, ${token.colorPrimaryHover});
            transform: scaleX(0);
            transition: transform 0.4s ease;
        }

        &:hover {
            background: ${isDarkMode ? token.colorBgElevated : token.colorBgContainer};
            box-shadow: ${isDarkMode
                    ? `0 12px 24px rgba(0, 0, 0, 0.4), 0 0 30px ${token.colorPrimary}20`
                    : `0 12px 24px rgba(0, 0, 0, 0.1), 0 0 20px ${token.colorPrimary}15`};
            border-color: ${token.colorPrimary}40;
            transform: translateY(-6px);

            &::before {
                transform: scaleX(1);
            }

            .advantage-number {
                transform: scale(1.15);
            }
        }

        @media (max-width: 768px) {
            padding: 28px 16px;
        }
    `,
    advantageNumber: css`
        font-size: 56px;
        font-weight: 800;
        background: linear-gradient(135deg, ${token.colorPrimary} 0%, ${token.colorPrimaryHover} 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
        margin-bottom: 16px;
        transition: transform 0.4s ease;

        @media (max-width: 768px) {
            font-size: 42px;
            margin-bottom: 12px;
        }
    `,
    advantageTitle: css`
        font-size: 19px;
        font-weight: 700;
        margin-bottom: 10px;
        color: ${token.colorText};

        @media (max-width: 768px) {
            font-size: 17px;
            margin-bottom: 8px;
        }
    `,
    advantageDesc: css`
        color: ${token.colorTextSecondary};
        font-size: 14px;
        line-height: 1.6;

        @media (max-width: 768px) {
            font-size: 13px;
        }
    `,
    ctaSection: css`
        background: ${isDarkMode
                ? `linear-gradient(135deg, ${token.colorPrimary}15 0%, ${token.colorPrimary}08 100%)`
                : `linear-gradient(135deg, ${token.colorPrimary}05 0%, ${token.colorPrimary}10 100%)`};
        width: 100vw;
        margin-left: calc(-50vw + 50%);
        padding: 80px 0;
        position: relative;
    `,
    ctaCard: css`
        background: transparent;
        border: none;
        border-radius: 0;
        text-align: center;
        padding: 0 20px;
        max-width: 1400px;
        margin: 0 auto;

        @media (max-width: 768px) {
            padding: 0 16px;
        }
    `,
    ctaTitle: css`
        margin-bottom: 16px !important;
        color: ${token.colorText} !important;

        @media (max-width: 768px) {
            font-size: 24px !important;
        }
    `,
    ctaDescription: css`
        font-size: 18px;
        margin-bottom: 40px;
        color: ${token.colorTextSecondary} !important;

        @media (max-width: 768px) {
            font-size: 16px;
            margin-bottom: 32px;
        }
    `,
    footer: css`
        background: ${isDarkMode ? token.colorBgElevated : token.colorBgContainer};
        border-top: 1px solid ${token.colorBorderSecondary};
        padding: 48px 20px 32px;
        text-align: center;
        width: 100%;
        margin: 0;

        @media (max-width: 1200px) {
            padding: 40px 16px 28px;
        }

        @media (max-width: 768px) {
            padding: 36px 16px 24px;
        }
    `,
    footerContainer: css`
        max-width: 1200px;
        margin: 0 auto;
    `,
    footerInfo: css`
        display: flex;
        flex-direction: column;
        gap: 12px;
        align-items: center;

        @media (max-width: 768px) {
            gap: 10px;
        }
    `,
    footerCopyright: css`
        font-size: 13px;
        color: ${token.colorTextTertiary};

        @media (max-width: 768px) {
            font-size: 12px;
        }
    `,
    footerBeian: css`
        font-size: 13px;
        color: ${token.colorTextTertiary};
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 8px;

        a {
            color: ${token.colorTextTertiary};
            transition: color 0.2s ease;
            text-decoration: none;

            &:hover {
                color: ${token.colorPrimary};
                text-decoration: none;
            }
        }

        @media (max-width: 768px) {
            font-size: 12px;
            flex-direction: column;
            gap: 6px;
        }
    `,
    backToTop: css`
        position: fixed;
        right: 24px;
        bottom: 24px;
        z-index: 1000;
        width: 48px;
        height: 48px;
        border-radius: 50%;
        border: none;
        background: ${token.colorPrimary};
        color: white;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 16px;
        box-shadow: ${isDarkMode
                ? `0 4px 12px ${token.colorPrimary}40, 0 2px 8px rgba(0, 0, 0, 0.4)`
                : `${token.colorPrimary}80 0 6px 27px 0, 0 2px 8px rgba(0, 0, 0, 0.15)`};
        backdrop-filter: blur(12px);
        -webkit-backdrop-filter: blur(12px);
        border: ${isDarkMode
                ? `1px solid ${token.colorBorder}`
                : 'none'};
        transition: all 0.3s ease;
        opacity: 0;
        transform: translateY(20px);
        pointer-events: none;
        font-weight: 600;

        &.visible {
            opacity: 1;
            transform: translateY(0);
            pointer-events: auto;
        }

        &:hover {
            background: ${token.colorPrimaryHover};
            transform: translateY(-4px) scale(1.05);
            box-shadow: ${isDarkMode
                    ? `0 6px 16px ${token.colorPrimaryHover}40, 0 4px 12px rgba(0, 0, 0, 0.5)`
                    : `${token.colorPrimaryHover}80 0 8px 32px 0, 0 4px 12px rgba(0, 0, 0, 0.2)`};
            color: white;
        }

        &:active {
            transform: translateY(-2px) scale(1.02);
        }

        @media (max-width: 768px) {
            right: 16px;
            bottom: 16px;
            width: 44px;
            height: 44px;
            font-size: 14px;
        }
    `,
}));

export default () => {
    const {styles} = useStyles();
    const {info} = useModel("web");
    const theme = useTheme();
    const user = useModel("user");

    // 统计数据 - 写死避免动画卡顿
    const stats = {
        userCount: 50000,
        serviceCount: 200,
        apiCount: 1000000,
        uptime: 99.9
    };

    // UI 状态
    const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
    const [scrolled, setScrolled] = useState(false);
    const [showBackToTop, setShowBackToTop] = useState(false);

    // 滚动监听 - 监听 window 滚动（AdaptiveScroll 已做事件桥接）
    useEffect(() => {
        let ticking = false;

        const handleScroll = () => {
            if (!ticking) {
                window.requestAnimationFrame(() => {
                    const scrollY = window.scrollY;
                    setScrolled(scrollY > 50);
                    setShowBackToTop(scrollY > 300);
                    ticking = false;
                });
                ticking = true;
            }
        };

        window.addEventListener('scroll', handleScroll, {passive: true});

        return () => {
            window.removeEventListener('scroll', handleScroll);
        };
    }, []);

    const features = [
        {
            icon: getAssetPath("images/target.svg"),
            title: "快速部署",
            description: "分钟级快速部署，无需复杂配置，开箱即用，让您的项目快速上线运营。",
        },
        {
            icon: getAssetPath("images/警示.svg"),
            title: "安全防护",
            description: "多重安全防护体系，数据加密传输，符合行业安全标准，保障业务安全。",
        },
        {
            icon: getAssetPath("images/profit.svg"),
            title: "高性能",
            description: "采用先进的技术架构，毫秒级响应速度，支持高并发访问，稳定可靠。",
        },
        {
            icon: getAssetPath("images/模块.svg"),
            title: "弹性扩展",
            description: "灵活的扩展机制，根据业务需求自动调整资源，应对流量高峰。",
        },
        {
            icon: getAssetPath("images/user.svg"),
            title: "专业支持",
            description: "7×24小时专业技术支持，快速响应您的需求，助力业务持续发展。",
        },
        {
            icon: getAssetPath("images/导航.svg"),
            title: "全球加速",
            description: "全球CDN加速节点，智能路由分发，无论用户在哪里都能获得极速体验。",
        }
    ];

    const products = [
        {
            icon: getAssetPath("images/定位.svg"),
            title: "智能监控",
            description: "实时监控系统运行状态，智能告警，让您随时掌握系统健康状况。",
            features: ["实时数据监控", "智能告警通知", "可视化报表", "性能分析"]
        },
        {
            icon: getAssetPath("images/设置.svg"),
            title: "API 服务",
            description: "提供丰富的API接口，文档完善，便于集成，助力快速开发。",
            features: ["RESTful API", "完善文档", "SDK支持", "高并发处理"]
        },
        {
            icon: getAssetPath("images/警示.svg"),
            title: "数据安全",
            description: "企业级数据安全方案，多重备份机制，确保数据万无一失。",
            features: ["数据加密", "定期备份", "灾难恢复", "权限管理"]
        },
        {
            icon: getAssetPath("images/指南.svg"),
            title: "自动化运维",
            description: "自动化部署与运维工具，提升运维效率，降低人工成本。",
            features: ["自动部署", "日志管理", "版本控制", "回滚机制"]
        }
    ];

    const testimonials = [
        {
            content: "使用这个平台后，我们的开发效率提升了3倍，服务稳定性也大幅提高。团队可以更专注于业务开发，而不是运维工作。强烈推荐！",
            name: "张伟",
            role: "技术总监",
            avatar: "张"
        },
        {
            content: "非常专业的服务，技术支持响应及时，功能强大且易用。帮助我们快速搭建了业务系统，节省了大量时间和成本。",
            name: "李娜",
            role: "产品经理",
            avatar: "李"
        },
        {
            content: "安全性和稳定性都非常出色，API文档清晰易懂，集成起来很顺畅。是我用过最好的服务平台之一。",
            name: "王强",
            role: "架构师",
            avatar: "王"
        }
    ];

    const advantages = [
        {number: "01", title: "简单易用", desc: "直观的操作界面，降低使用门槛"},
        {number: "02", title: "灵活定制", desc: "支持个性化配置，满足不同需求"},
        {number: "03", title: "持续创新", desc: "不断更新迭代，引领技术潮流"},
        {number: "04", title: "生态完善", desc: "丰富的生态系统，无缝对接"}
    ];

    const handleGetStarted = () => {
        if (user?.web) {
            historyPush("user.index");
        } else {
            historyPush("login");
        }
    };

    const handleLogout = async () => {
        try {
            await logout({
                onSuccess: () => {
                    deleteHeader();
                    user?.setWeb(undefined);
                },
                onFail: () => {
                }
            });
        } catch {
        }
    };

    const scrollToSection = (id: string) => {
        const element = document.getElementById(id);
        if (element) {
            const offsetTop = element.offsetTop - 80;
            window.scrollTo({
                top: offsetTop,
                behavior: 'smooth'
            });
        }
        setMobileMenuOpen(false);
    };

    const scrollToTop = () => {
        window.scrollTo({
            top: 0,
            behavior: 'smooth'
        });
    };

    const navItems = [
        {label: '首页', id: 'home'},
        {label: '特性', id: 'features'},
        {label: '产品', id: 'products'},
        {label: '评价', id: 'testimonials'}
    ];

    return (
        <Body space={false} animation={false} className={styles.noPadding}>
            {/* 固定的 Navigation */}
            <nav className={`${styles.navbar} ${scrolled ? 'scrolled' : ''}`}>
                <div className={styles.navWrapper}>
                    <div className={styles.navContainer}>
                        <div className={styles.logo} onClick={() => scrollToSection('home')}>
                            <div className={styles.logoIcon}>
                                {info?.logo ? (
                                    <img
                                        src={info?.logo || getAssetPath('logo.svg')}
                                        alt="logo"
                                    />
                                ) : (
                                    <Icon type="icon-logo"/>
                                )}
                            </div>
                            <span className={styles.logoText}>{info?.name || Settings?.title || "云上豁者"}</span>
                        </div>

                        <div className={styles.navMenu}>
                            {navItems.map((item) => (
                                <div
                                    key={item.id}
                                    className={styles.navLink}
                                    onClick={() => scrollToSection(item.id)}
                                >
                                    {item.label}
                                </div>
                            ))}
                        </div>

                        <div className={styles.navActions}>
                            <button
                                className={styles.themeToggle}
                                onClick={() => theme?.toggle?.()}
                                title={theme?.getModeName(theme?.mode as any)}
                            >
                                <Icon
                                    type={theme?.isDarkMode() ? "icon-dark" : (theme?.isAutoMode() ? 'icon-auto' : 'icon-light')}/>
                            </button>
                            {user?.web ? (
                                <>
                                    <button
                                        className={`${styles.navBtn} btn-secondary`}
                                        onClick={() => historyPush("user.index")}
                                    >
                                        <HomeOutlined/>
                                        进入系统
                                    </button>
                                    <button
                                        className={`${styles.navBtn} btn-primary`}
                                        onClick={handleLogout}
                                    >
                                        <LogoutOutlined/>
                                        退出
                                    </button>
                                </>
                            ) : (
                                <>
                                    <button
                                        className={`${styles.navBtn} btn-secondary`}
                                        onClick={() => historyPush("login")}
                                    >
                                        <LoginOutlined/>
                                        登录
                                    </button>
                                    <button
                                        className={`${styles.navBtn} btn-primary`}
                                        onClick={handleGetStarted}
                                    >
                                        开始使用
                                    </button>
                                </>
                            )}
                        </div>

                        <button
                            className={styles.mobileMenuBtn}
                            onClick={() => setMobileMenuOpen(true)}
                        >
                            <MenuOutlined/>
                        </button>
                    </div>
                </div>
            </nav>

            {/* Mobile Menu - 预渲染，纯 CSS 动画 */}
            <div
                className={`${styles.mobileMenuMask} ${mobileMenuOpen ? 'open' : ''}`}
                onClick={() => setMobileMenuOpen(false)}
            />
            <div className={`${styles.mobileMenuDrawer} ${mobileMenuOpen ? 'open' : ''}`}>
                <div className={styles.mobileMenuHeader}>
                    <button
                        className={styles.mobileMenuClose}
                        onClick={() => setMobileMenuOpen(false)}
                    >
                        <CloseOutlined/>
                    </button>
                </div>

                <div className={styles.mobileMenuBody}>
                    <div className={styles.mobileMenuNav}>
                        {navItems.map((item, _) => (
                            <div
                                key={item.id}
                                className={styles.mobileMenuItem}
                                onClick={() => scrollToSection(item.id)}
                            >
                                {item.label}
                            </div>
                        ))}
                    </div>

                    <div className={styles.mobileThemeSection}>
                        <div className={styles.mobileThemeLabel}>主题模式</div>
                        <div className={styles.mobileThemeButtons}>
                            <button
                                className={`${styles.mobileThemeBtn} ${theme?.isLightMode() ? 'active' : ''}`}
                                onClick={() => theme?.setLightMode?.()}
                            >
                                <Icon type={"icon-light"}/>
                                亮色
                            </button>
                            <button
                                className={`${styles.mobileThemeBtn} ${theme?.isDarkMode() ? 'active' : ''}`}
                                onClick={() => theme?.setDarkMode?.()}
                            >
                                <Icon type={"icon-dark"}/>
                                暗色
                            </button>
                            <button
                                className={`${styles.mobileThemeBtn} ${theme?.isAutoMode() ? 'active' : ''}`}
                                onClick={() => theme?.setAutoMode?.()}
                            >
                                <Icon type={"icon-auto"}/>
                                自动
                            </button>
                        </div>
                    </div>

                    <div className={styles.mobileActions}>
                        {user?.web ? (
                            <>
                                <button
                                    className={`${styles.mobileActionItem} ${styles.mobileBtnPrimary}`}
                                    onClick={() => {
                                        historyPush("user.index");
                                        setMobileMenuOpen(false);
                                    }}
                                >
                                    <HomeOutlined/>
                                    <span>进入系统</span>
                                </button>
                                <button
                                    className={styles.mobileActionItem}
                                    onClick={() => {
                                        handleLogout();
                                        setMobileMenuOpen(false);
                                    }}
                                >
                                    <LogoutOutlined/>
                                    <span>退出登录</span>
                                </button>
                            </>
                        ) : (
                            <>
                                <button
                                    className={`${styles.mobileActionItem} ${styles.mobileBtnPrimary}`}
                                    onClick={() => {
                                        handleGetStarted();
                                        setMobileMenuOpen(false);
                                    }}
                                >
                                    <RocketOutlined/>
                                    <span>免费开始</span>
                                </button>
                                <button
                                    className={styles.mobileActionItem}
                                    onClick={() => {
                                        historyPush("login");
                                        setMobileMenuOpen(false);
                                    }}
                                >
                                    <LoginOutlined/>
                                    <span>登录账号</span>
                                </button>
                            </>
                        )}
                    </div>
                </div>
            </div>

            {/* 页面内容 */}
            <div className={styles.pageWrapper}>
                <div className={styles.contentWrapper}>
                    {/* Hero Section - 科技感背景动画 */}
                    <div id="home" className={styles.heroSection}>
                        {/* 基础背景层 */}
                        <div className={styles.heroGradient}></div>
                        <div className={styles.heroSpotlight}></div>

                        {/* 科技感动画层 - 轻量级纯CSS动画 */}
                        <div className={styles.heroGrid}></div>
                        <div className={styles.heroScanLine}></div>
                        <div className={styles.heroCircuits}></div>
                        <div className={styles.heroDots}></div>
                        <div className={styles.heroDataStream}></div>
                        <div className={styles.heroTechLines}></div>

                        {/* 发光球 - 已减少到 2 个 */}
                        <div className={styles.heroGlowOrbs}>
                            {ANIMATION_ELEMENTS.glowOrbs.map((orb) => (
                                <div
                                    key={orb.id}
                                    className={styles.glowOrb}
                                    style={{
                                        left: orb.left,
                                        top: orb.top,
                                        width: orb.size,
                                        height: orb.size,
                                        '--duration': orb.duration,
                                        '--delay': orb.delay,
                                        '--tx': orb.tx,
                                        '--ty': orb.ty
                                    } as React.CSSProperties}
                                />
                            ))}
                        </div>

                        {/* 粒子 - 已减少到 2 个 */}
                        <div className={styles.heroParticles}>
                            {ANIMATION_ELEMENTS.particles.map((particle) => (
                                <div
                                    key={particle.id}
                                    className={styles.particle}
                                    style={{
                                        left: particle.left,
                                        top: particle.top,
                                        '--duration': particle.duration,
                                        '--delay': particle.delay,
                                        '--tx': particle.tx,
                                        '--ty': particle.ty
                                    } as React.CSSProperties}
                                />
                            ))}
                        </div>

                        <div className={styles.heroContent}>
                            <h1 className={`${styles.title} ${styles.h1} ${styles.heroTitle}`}>
                                {info?.name || Settings?.title}
                            </h1>
                            <h3 className={`${styles.title} ${styles.h3} ${styles.heroSubtitle}`}>
                                {(info as any)?.describe || "赋能企业数字化转型，构建高效业务系统"}
                            </h3>

                            <div className={styles.ctaButtons}>
                                <button
                                    className={styles.btnPrimary}
                                    onClick={handleGetStarted}
                                >
                                    {user?.web ? (
                                        <>
                                            <HomeOutlined/>
                                            进入系统
                                        </>
                                    ) : (
                                        <>
                                            <RocketOutlined/>
                                            开始体验
                                        </>
                                    )}
                                </button>
                                {!user?.web && (
                                    <button
                                        className={styles.btnSecondary}
                                        onClick={() => historyPush("login")}
                                    >
                                        已有账号？登录
                                    </button>
                                )}
                            </div>
                        </div>
                    </div>

                    {/* Stats Section */}
                    <div className={styles.statsSection}>
                        <div className={styles.statsContainer}>
                            <Row gutter={[16, 24]}>
                                <Col xs={12} sm={12} md={6}>
                                    <div className={styles.statCard}>
                                        <Statistic
                                            title="服务用户"
                                            value={stats.userCount}
                                            suffix="+"
                                        />
                                    </div>
                                </Col>
                                <Col xs={12} sm={12} md={6}>
                                    <div className={styles.statCard}>
                                        <Statistic
                                            title="功能服务"
                                            value={stats.serviceCount}
                                            suffix="+"
                                        />
                                    </div>
                                </Col>
                                <Col xs={12} sm={12} md={6}>
                                    <div className={styles.statCard}>
                                        <Statistic
                                            title="API 调用量"
                                            value={stats.apiCount}
                                            suffix="+"
                                        />
                                    </div>
                                </Col>
                                <Col xs={12} sm={12} md={6}>
                                    <div className={styles.statCard}>
                                        <Statistic
                                            title="稳定性"
                                            value={stats.uptime}
                                            suffix="%"
                                            precision={1}
                                        />
                                    </div>
                                </Col>
                            </Row>
                        </div>
                    </div>

                    {/* Features Section */}
                    <div id="features" className={styles.section}>
                        <h2 className={`${styles.title} ${styles.h2} ${styles.sectionTitle}`}>
                            核心特性
                        </h2>
                        <p className={`${styles.paragraph} ${styles.sectionSubtitle}`}>
                            为您提供企业级的服务能力，助力业务快速发展
                        </p>

                        <Row gutter={[24, 24]}>
                            {features.map((feature, index) => (
                                <Col xs={24} sm={12} md={8} lg={8} key={index}>
                                    <div className={styles.featureCard}>
                                        <div className={styles.featureIcon}>
                                            <img src={feature.icon} alt={feature.title}/>
                                        </div>
                                        <h4 className={`${styles.title} ${styles.h4} ${styles.featureTitle}`}>
                                            {feature.title}
                                        </h4>
                                        <p className={`${styles.paragraph} ${styles.featureDescription}`}>
                                            {feature.description}
                                        </p>
                                    </div>
                                </Col>
                            ))}
                        </Row>
                    </div>

                    {/* Products Section */}
                    <div id="products" className={styles.advantageSection}>
                        <div className={styles.section}>
                            <h2 className={`${styles.title} ${styles.h2} ${styles.sectionTitle}`}>
                                产品服务
                            </h2>
                            <p className={`${styles.paragraph} ${styles.sectionSubtitle}`}>
                                丰富的产品矩阵，满足不同场景需求
                            </p>

                            <Row gutter={[24, 24]}>
                                {products.map((product, index) => (
                                    <Col xs={24} sm={12} md={12} lg={6} key={index}>
                                        <div className={styles.productCard}>
                                            <div className={`${styles.productIconBg} product-icon-bg`}>
                                                <img src={product.icon} alt={product.title}/>
                                            </div>
                                            <h4 className={`${styles.title} ${styles.h4} ${styles.productTitle}`}>
                                                {product.title}
                                            </h4>
                                            <p className={`${styles.paragraph} ${styles.productDesc}`}>
                                                {product.description}
                                            </p>
                                            <div className={styles.productFeatures}>
                                                {product.features.map((feature, idx) => (
                                                    <div key={idx} className={styles.productFeatureItem}>
                                                        <CheckCircleOutlined/>
                                                        <span>{feature}</span>
                                                    </div>
                                                ))}
                                            </div>
                                        </div>
                                    </Col>
                                ))}
                            </Row>
                        </div>
                    </div>

                    {/* Advantages Section */}
                    <div className={styles.section}>
                        <h2 className={`${styles.title} ${styles.h2} ${styles.sectionTitle}`}>
                            为什么选择我们
                        </h2>
                        <p className={`${styles.paragraph} ${styles.sectionSubtitle}`}>
                            专业、可靠、创新，是我们的核心价值
                        </p>

                        <Row gutter={[16, 16]}>
                            {advantages.map((advantage, index) => (
                                <Col xs={12} sm={12} md={6} key={index}>
                                    <div className={styles.advantageCard}>
                                        <div className={`${styles.advantageNumber} advantage-number`}>
                                            {advantage.number}
                                        </div>
                                        <div className={styles.advantageTitle}>
                                            {advantage.title}
                                        </div>
                                        <div className={styles.advantageDesc}>
                                            {advantage.desc}
                                        </div>
                                    </div>
                                </Col>
                            ))}
                        </Row>
                    </div>

                    {/* Testimonials Section */}
                    <div id="testimonials" className={styles.advantageSection}>
                        <div className={styles.section}>
                            <h2 className={`${styles.title} ${styles.h2} ${styles.sectionTitle}`}>
                                用户评价
                            </h2>
                            <p className={`${styles.paragraph} ${styles.sectionSubtitle}`}>
                                听听用户怎么说
                            </p>

                            <Row gutter={[24, 24]}>
                                {testimonials.map((testimonial, index) => (
                                    <Col xs={24} sm={24} md={8} key={index}>
                                        <div className={styles.testimonialCard}>
                                            <div className={styles.testimonialStars}>
                                                {[1, 2, 3, 4, 5].map((star) => (
                                                    <StarFilled key={star}/>
                                                ))}
                                            </div>
                                            <div className={styles.testimonialContent}>
                                                "{testimonial.content}"
                                            </div>
                                            <div className={styles.testimonialAuthor}>
                                                <div className={styles.testimonialAvatar}>
                                                    {testimonial.avatar}
                                                </div>
                                                <div className={styles.testimonialInfo}>
                                                    <div className={styles.testimonialName}>
                                                        {testimonial.name}
                                                    </div>
                                                    <div className={styles.testimonialRole}>
                                                        {testimonial.role}
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    </Col>
                                ))}
                            </Row>
                        </div>
                    </div>

                    {/* CTA Section */}
                    <div className={styles.ctaSection}>
                        <div className={styles.ctaCard}>
                            <h2 className={`${styles.title} ${styles.h2} ${styles.ctaTitle}`}>
                                准备好开始了吗？
                            </h2>
                            <p className={`${styles.paragraph} ${styles.ctaDescription}`}>
                                立即注册，开启您的数字化之旅
                            </p>
                            <div className={styles.ctaButtons}>
                                <button
                                    className={styles.btnPrimary}
                                    onClick={() => user?.web ? historyPush("user.index") : historyPush("login")}
                                >
                                    {user?.web ? (
                                        <>
                                            <HomeOutlined/>
                                            进入控制台
                                        </>
                                    ) : (
                                        <>
                                            <RocketOutlined/>
                                            立即开始
                                        </>
                                    )}
                                </button>
                                <button
                                    className={styles.btnSecondary}
                                    onClick={() => window.location.href = `mailto:${info?.contact || 'contact@example.com'}`}
                                >
                                    咨询合作
                                </button>
                            </div>
                        </div>
                    </div>

                    {/* Footer */}
                    <footer className={styles.footer}>
                        <div className={styles.footerContainer}>

                            <div className={styles.footerInfo}>
                                <div className={styles.footerCopyright}>
                                    © {new Date().getFullYear()} {info?.name || Settings?.title || "云上豁者"}. All
                                    rights
                                    reserved.
                                </div>
                                <div className={styles.footerBeian}>
                                    {(info as any)?.icp && (
                                        <a
                                            href="https://beian.miit.gov.cn"
                                            target="_blank"
                                            rel="noopener noreferrer"
                                        >
                                            {(info as any)?.icp}
                                        </a>
                                    )}
                                    {(info as any)?.police_beian && (
                                        <>
                                            {(info as any)?.icp && <span>|</span>}
                                            <a
                                                href="http://www.beian.gov.cn/portal/registerSystemInfo"
                                                target="_blank"
                                                rel="noopener noreferrer"
                                            >
                                                {(info as any)?.police_beian}
                                            </a>
                                        </>
                                    )}
                                </div>
                            </div>
                        </div>
                    </footer>

                    {/* Back to Top Button */}
                    <button
                        className={`${styles.backToTop} ${showBackToTop ? 'visible' : ''}`}
                        onClick={scrollToTop}
                        title="返回顶部"
                    >
                        <Icon type="UpOutlined"/>
                    </button>
                </div>
            </div>
        </Body>
    );
};