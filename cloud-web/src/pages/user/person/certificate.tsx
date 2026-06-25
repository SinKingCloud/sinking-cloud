import React, {useEffect, useMemo, useRef, useState} from "react";
import {App, Button, Card, Form, Input, QRCode, Row} from "antd";
import {createStyles} from "antd-style";
import {Body, Icon, useTheme} from "sinking-antd";
import {certificate} from "@/service/user/person";
import Script from "@/pages/components/script";
import {getPath} from "@/utils/assets";

const REAL_AUTH_TYPE = {
    PERSON: 1,
    COMPANY: 2,
};

const CERTIFY_CACHE_KEY = "user.real.auth.certify";
const CERTIFY_POLL_INTERVAL = 3000;
const CERTIFY_POLL_MAX_COUNT = 60;
const ALIYUN_META_SCRIPT = getPath("/aliyun/face.js");

const AUTH_TYPE_OPTIONS = [
    {
        label: "个人认证",
        desc: "适用于个人身份信息核验",
        value: REAL_AUTH_TYPE.PERSON,
        icon: "UserOutlined",
    },
    {
        label: "企业认证",
        desc: "适用于企业主体和经办人核验",
        value: REAL_AUTH_TYPE.COMPANY,
        icon: "BankOutlined",
    },
];

const STEP_ITEMS = [
    {title: "选择类型", desc: "确认认证主体", icon: "SafetyCertificateOutlined"},
    {title: "填写资料", desc: "录入证件信息", icon: "EditOutlined"},
    {title: "在线核验", desc: "完成实名校验", icon: "AuditOutlined"},
    {title: "认证完成", desc: "资料正式生效", icon: "CheckCircleOutlined"},
];

const useStyles = createStyles(({css, token, isDarkMode, responsive}: any, props: any = {}) => {
    const isCompactMode = props?.isCompactMode || false;
    const panelBg = isDarkMode
        ? `linear-gradient(135deg, color-mix(in srgb, ${token.colorPrimary} 8%, ${token.colorBgContainer}) 0%, ${token.colorBgContainer} 100%)`
        : token.colorFillAlter;

    return {
        page: css`
            width: 100%;
            max-width: 1440px;
            margin: 0 auto;
            display: grid;
            gap: ${isCompactMode ? 10 : 12}px;

            ${responsive.mobile} {
                min-height: auto;
                padding: 0;
            }

            @media (max-width: 900px) {
                min-height: auto;
            }
        `,
        shell: css`
            width: 100%;
            overflow: hidden;
            border: none;
            border-radius: ${token.borderRadiusLG}px;
            background: ${token.colorBgContainer};
            box-shadow: ${isDarkMode ? "none" : "0 16px 42px rgba(15, 23, 42, 0.055)"};
        `,
        hero: css`
            position: relative;
            display: grid;
            grid-template-columns: minmax(0, 1fr) 260px;
            gap: ${isCompactMode ? 16 : 20}px;
            padding: ${isCompactMode ? "20px" : "26px"};
            overflow: hidden;
            border-bottom: 1px solid ${token.colorBorderSecondary};
            background: ${panelBg};
            isolation: isolate;

            &::before {
                content: '';
                position: absolute;
                inset: -120px -160px;
                background: radial-gradient(circle at 16% 20%, ${token.colorPrimaryBgHover} 0%, transparent 34%),
                radial-gradient(circle at 82% 28%, ${token.colorPrimaryBg} 0%, transparent 32%),
                radial-gradient(circle at 70% 92%, ${token.colorPrimaryBorder} 0%, transparent 36%);
                opacity: ${isDarkMode ? 0.32 : 0.68};
                filter: blur(2px);
                transform: rotate(-6deg);
                z-index: -2;
            }

            &::after {
                content: '';
                position: absolute;
                right: -190px;
                bottom: -220px;
                width: 560px;
                height: 560px;
                border-radius: 50%;
                background: radial-gradient(circle, ${token.colorPrimaryBgHover} 0%, transparent 48%),
                conic-gradient(from 36deg, transparent 0 14%, ${token.colorPrimaryBorder} 22%, transparent 38%, ${token.colorPrimaryBgHover} 56%, transparent 76%, ${token.colorPrimaryBg} 88%, transparent 100%);
                filter: blur(58px);
                opacity: ${isDarkMode ? 0.16 : 0.36};
                z-index: -1;
            }

            ${responsive.mobile} {
                grid-template-columns: 1fr;
                gap: ${isCompactMode ? 14 : 16}px;
                padding: ${isCompactMode ? 18 : 20}px;
            }

            @media (max-width: 900px) {
                grid-template-columns: 1fr;
                gap: ${isCompactMode ? 14 : 16}px;
                padding: ${isCompactMode ? "18px" : "20px"};
            }
        `,
        heroMain: css`
            position: relative;
            z-index: 1;
            min-width: 0;
            display: flex;
            flex-direction: column;
            justify-content: center;
            gap: ${isCompactMode ? 8 : 10}px;

            .hero-title {
                min-width: 0;
                color: ${token.colorText};
                font-size: ${isCompactMode ? 26 : 28}px;
                font-weight: 700;
                line-height: ${isCompactMode ? 34 : 36}px;
                letter-spacing: 0;
            }

            .desc {
                max-width: 600px;
                color: ${token.colorTextSecondary};
                font-size: 13px;
                line-height: 22px;
            }

            ${responsive.mobile} {
                .hero-title {
                    font-size: ${isCompactMode ? 23 : 24}px;
                    line-height: ${isCompactMode ? 31 : 32}px;
                }
            }

            @media (max-width: 560px) {
                .hero-title {
                    font-size: ${isCompactMode ? 22 : 23}px;
                    line-height: ${isCompactMode ? 30 : 31}px;
                }

                .desc {
                    font-size: 12px;
                    line-height: 20px;
                }
            }
        `,
        heroVisual: css`
            position: relative;
            z-index: 1;
            min-height: ${isCompactMode ? 118 : 132}px;
            padding: ${isCompactMode ? "16px" : "18px"};
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 12px;
            overflow: hidden;
            border-radius: ${token.borderRadiusLG}px;
            background: ${isDarkMode ? "rgba(255,255,255,0.035)" : "rgba(255,255,255,0.72)"};
            box-shadow: ${isDarkMode ? "none" : "0 12px 28px rgba(15, 23, 42, 0.045)"};

            .visual-info {
                min-width: 0;
            }

            .visual-label {
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 18px;
            }

            .visual-value {
                margin-top: 5px;
                color: ${token.colorText};
                font-size: ${isCompactMode ? 20 : 22}px;
                font-weight: 700;
                line-height: 28px;
            }

            .visual-desc {
                margin-top: 5px;
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 18px;
            }

            .visual-icon {
                position: relative;
                width: ${isCompactMode ? 82 : 88}px;
                height: ${isCompactMode ? 68 : 72}px;
                flex: none;
                display: flex;
                align-items: center;
                justify-content: center;
                color: ${token.colorPrimary};
            }

            .visual-icon svg {
                width: 100%;
                height: 100%;
                display: block;
            }

            .visual-icon::before {
                content: '';
                position: absolute;
                inset: -34px;
                border-radius: 50%;
                background: radial-gradient(circle, ${token.colorPrimaryBgHover}, transparent 68%);
                filter: blur(24px);
                opacity: ${isDarkMode ? 0.18 : 0.32};
                z-index: -1;
            }

            @media (max-width: 900px) {
                min-height: auto;
                padding: ${isCompactMode ? "12px 14px" : "14px 16px"};
            }

            @media (max-width: 560px) {
                .visual-icon {
                    width: ${isCompactMode ? 66 : 70}px;
                    height: ${isCompactMode ? 54 : 58}px;
                }

                .visual-value {
                    font-size: ${isCompactMode ? 18 : 20}px;
                    line-height: 26px;
                }
            }
        `,
        workspace: css`
            display: grid;
            grid-template-columns: 320px minmax(0, 1fr);

            ${responsive.mobile} {
                grid-template-columns: 1fr;
            }

            @media (max-width: 900px) {
                grid-template-columns: 1fr;
            }
        `,
        guide: css`
            padding: ${isCompactMode ? "16px" : "20px"};
            border-right: 1px solid ${token.colorBorderSecondary};
            background: ${isDarkMode ? "rgba(255,255,255,0.01)" : token.colorBgContainer};

            .guide-title {
                color: ${token.colorText};
                font-size: 14px;
                font-weight: 600;
                line-height: 22px;
            }

            .guide-desc {
                margin-top: 4px;
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 18px;
            }

            ${responsive.mobile} {
                padding: ${isCompactMode ? 14 : 16}px;
                border-right: none;
                border-bottom: 1px solid ${token.colorBorderSecondary};
            }

            @media (max-width: 900px) {
                padding: ${isCompactMode ? 14 : 16}px;
                border-right: none;
                border-bottom: 1px solid ${token.colorBorderSecondary};
            }
        `,
        stepList: css`
            display: grid;
            gap: ${isCompactMode ? 8 : 10}px;
            margin-top: ${isCompactMode ? 14 : 16}px;

            ${responsive.mobile} {
                grid-template-columns: repeat(4, minmax(0, 1fr));
                gap: 6px;
            }

            @media (max-width: 900px) {
                grid-template-columns: repeat(4, minmax(0, 1fr));
                gap: 6px;
            }
        `,
        stepItem: css`
            position: relative;
            display: grid;
            grid-template-columns: ${isCompactMode ? 28 : 30}px minmax(0, 1fr);
            gap: ${isCompactMode ? 8 : 9}px;
            min-width: 0;
            padding: ${isCompactMode ? "8px 10px 8px 11px" : "9px 10px 9px 12px"};
            border-radius: ${token.borderRadius}px;
            border: 1px solid transparent;
            background: transparent;

            &::after {
                content: '';
                position: absolute;
                left: ${isCompactMode ? 25 : 27}px;
                top: ${isCompactMode ? 36 : 38}px;
                bottom: ${isCompactMode ? -10 : -12}px;
                width: 1px;
                background: ${token.colorBorderSecondary};
            }

            &:last-child::after {
                display: none;
            }

            .step-icon {
                position: relative;
                z-index: 1;
                width: ${isCompactMode ? 28 : 30}px;
                height: ${isCompactMode ? 28 : 30}px;
                display: flex;
                align-items: center;
                justify-content: center;
                border: 1px solid ${token.colorBorderSecondary};
                border-radius: 50%;
                color: ${token.colorTextTertiary};
                background: ${token.colorBgContainer};
                font-size: ${isCompactMode ? 13 : 14}px;
            }

            .step-title {
                color: ${token.colorTextSecondary};
                font-size: ${isCompactMode ? 12 : 12}px;
                line-height: 18px;
            }

            .step-desc {
                margin-top: 1px;
                color: ${token.colorTextTertiary};
                font-size: 10px;
                line-height: 14px;
            }

            &.active {
                background: ${token.colorPrimaryBg};
            }

            &.active .step-icon {
                color: #fff;
                border-color: ${token.colorPrimary};
                background: ${token.colorPrimary};
                box-shadow: 0 4px 12px ${token.colorPrimaryBgHover};
            }

            &.done .step-icon {
                color: ${token.colorTextTertiary};
                border-color: ${token.colorBorderSecondary};
                background: ${token.colorBgContainer};
            }

            &.active .step-title {
                color: ${token.colorPrimary};
                font-weight: 600;
            }

            ${responsive.mobile} {
                display: flex;
                align-items: center;
                flex-direction: column;
                gap: 5px;
                padding: ${isCompactMode ? "8px 10px" : "9px 11px"};
                border-color: ${token.colorBorderSecondary};
                background: ${token.colorBgContainer};
                text-align: center;

                &::after {
                    display: none;
                }

                .step-icon {
                    width: ${isCompactMode ? 26 : 28}px;
                    height: ${isCompactMode ? 26 : 28}px;
                    font-size: ${isCompactMode ? 12 : 13}px;
                }

                .step-title {
                    font-size: 11px;
                    line-height: 15px;
                }

                .step-desc {
                    display: none;
                }
            }

            @media (max-width: 900px) {
                display: flex;
                align-items: center;
                flex-direction: column;
                gap: 5px;
                padding: ${isCompactMode ? "8px 9px" : "9px 10px"};
                border-color: ${token.colorBorderSecondary};
                background: ${token.colorBgContainer};
                text-align: center;

                &::after {
                    display: none;
                }

                .step-icon {
                    width: ${isCompactMode ? 26 : 28}px;
                    height: ${isCompactMode ? 26 : 28}px;
                    font-size: ${isCompactMode ? 12 : 13}px;
                }

                .step-title {
                    font-size: 11px;
                    line-height: 15px;
                }

                .step-desc {
                    display: none;
                }
            }

            @media (max-width: 430px) {
                padding: ${isCompactMode ? "7px 4px" : "8px 5px"};

                .step-icon {
                    width: ${isCompactMode ? 24 : 26}px;
                    height: ${isCompactMode ? 24 : 26}px;
                }

                .step-title {
                    font-size: 10px;
                    line-height: 14px;
                }
            }
        `,
        content: css`
            min-width: 0;
            padding: ${isCompactMode ? "18px" : "22px"};

            ${responsive.mobile} {
                padding: ${isCompactMode ? 14 : 16}px;
            }

            @media (max-width: 900px) {
                padding: ${isCompactMode ? 16 : 18}px;
            }

            @media (max-width: 560px) {
                padding: ${isCompactMode ? 14 : 16}px;
            }
        `,
        sectionHead: css`
            display: flex;
            align-items: flex-start;
            justify-content: space-between;
            gap: 12px;
            margin-bottom: ${isCompactMode ? 14 : 16}px;

            .section-title {
                color: ${token.colorText};
                font-size: ${isCompactMode ? 19 : 20}px;
                font-weight: 700;
                line-height: ${isCompactMode ? 28 : 30}px;
            }

            .section-desc {
                margin-top: 4px;
                color: ${token.colorTextTertiary};
                font-size: 13px;
                line-height: 20px;
            }

            ${responsive.mobile} {
                flex-direction: column;
            }

            @media (max-width: 560px) {
                margin-bottom: ${isCompactMode ? 12 : 14}px;

                .section-title {
                    font-size: ${isCompactMode ? 18 : 19}px;
                    line-height: ${isCompactMode ? 26 : 28}px;
                }

                .section-desc {
                    font-size: 12px;
                    line-height: 19px;
                }
            }
        `,
        formContent: css`
            width: min(370px, 100%);
            margin: 0 auto;

            ${responsive.mobile} {
                width: 100%;
            }
        `,
        typeGrid: css`
            display: grid;
            grid-template-columns: 1fr;
            gap: 10px;
        `,
        typeItem: css`
            width: 100%;
            min-height: auto;
            padding: ${isCompactMode ? "12px" : "13px 14px"};
            display: grid;
            grid-template-columns: ${isCompactMode ? 34 : 36}px minmax(0, 1fr) 18px;
            align-items: center;
            gap: ${isCompactMode ? 10 : 12}px;
            border: 1px solid ${token.colorBorderSecondary};
            border-radius: ${token.borderRadius}px;
            background: ${token.colorBgContainer};
            text-align: left;
            cursor: pointer;
            transition: border-color .2s ease, background .2s ease, box-shadow .2s ease;

            .type-top {
                display: contents;
            }

            .type-icon {
                width: ${isCompactMode ? 34 : 36}px;
                height: ${isCompactMode ? 34 : 36}px;
                display: flex;
                align-items: center;
                justify-content: center;
                border-radius: 50%;
                color: ${token.colorPrimary};
                background: ${token.colorPrimaryBg};
                font-size: ${isCompactMode ? 16 : 17}px;
            }

            .type-check {
                width: 18px;
                height: 18px;
                flex: none;
                grid-column: 3;
                grid-row: 1;
                border-radius: 50%;
                border: 1px solid ${token.colorBorder};
                background: ${token.colorBgContainer};
            }

            .title {
                color: ${token.colorText};
                font-size: ${isCompactMode ? 14 : 14}px;
                line-height: 22px;
                font-weight: 600;
            }

            .desc {
                margin-top: 2px;
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 18px;
            }

            &:hover {
                border-color: ${token.colorPrimaryBorder};
                box-shadow: ${isDarkMode ? "none" : "0 10px 24px rgba(15, 23, 42, 0.055)"};
            }

            &.active {
                border-color: ${token.colorPrimaryBorder};
                background: ${token.colorPrimaryBg};
            }

            &.active .type-check {
                border-color: ${token.colorPrimary};
                background: radial-gradient(circle, ${token.colorPrimary} 0 42%, transparent 48%);
            }
        `,
        actionBar: css`
            display: flex;
            align-items: center;
            justify-content: flex-end;
            gap: 8px;
            margin-top: ${isCompactMode ? 14 : 16}px;

            ${responsive.mobile} {
                display: grid;
                grid-template-columns: 1fr;
            }
        `,
        formActionBar: css`
            width: min(370px, 100%);
            margin: ${isCompactMode ? 14 : 16}px auto 0;
            padding-top: 0;
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 10px;

            &.single {
                justify-content: flex-end;
            }

            ${responsive.mobile} {
                width: 100%;
                display: flex;
                align-items: center;
                justify-content: space-between;
                gap: 8px;
            }
        `,
        verifyPanel: css`
            position: relative;
            width: min(370px, 100%);
            margin: 0 auto;
            display: grid;
            justify-items: center;
            gap: ${isCompactMode ? 12 : 14}px;
            padding: ${isCompactMode ? "4px 0" : "6px 0"};

            &::after {
                content: '认证过程由第三方实名服务完成';
                width: 100%;
                margin-top: ${isCompactMode ? 2 : 4}px;
                padding: ${isCompactMode ? "8px 10px" : "9px 12px"};
                border-radius: ${token.borderRadius}px;
                color: ${token.colorTextTertiary};
                background: ${isDarkMode ? "rgba(255,255,255,0.014)" : token.colorFillQuaternary};
                font-size: 12px;
                line-height: 18px;
                text-align: center;
            }
        `,
        verifyHead: css`
            display: flex;
            flex-direction: column;
            align-items: flex-start;
            width: 100%;
            text-align: left;
            gap: 4px;

            .title {
                color: ${token.colorText};
                font-size: ${isCompactMode ? 19 : 20}px;
                line-height: ${isCompactMode ? 27 : 28}px;
                font-weight: 700;
            }

            .desc {
                max-width: 300px;
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 19px;
            }
        `,
        qrCard: css`
            display: flex;
            align-items: center;
            justify-content: center;
            padding: ${isCompactMode ? "6px" : "8px"};
            border-radius: ${token.borderRadius}px;
            background: ${token.colorBgContainer};
        `,
        verifyActions: css`
            display: grid;
            gap: ${isCompactMode ? 7 : 8}px;
            margin-top: 2px;
            width: 100%;

            .verify-extra {
                min-width: 0;
                display: flex;
                align-items: center;
                justify-content: space-between;
                gap: 8px;
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 18px;
            }

            .verify-extra .ant-btn {
                height: 20px;
                padding: 0;
                font-size: 12px;
            }

            @media (max-width: 430px) {
                .verify-extra {
                    justify-content: center;
                    flex-direction: column;
                    gap: 3px;
                }
            }
        `,
        pending: css`
            min-height: ${isCompactMode ? 260 : 280}px;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            text-align: center;
            padding: ${isCompactMode ? 20 : 24}px;
            border: 1px solid ${token.colorBorderSecondary};
            border-radius: ${token.borderRadius}px;
            background: ${isDarkMode ? "rgba(255,255,255,0.014)" : token.colorFillQuaternary};

            .pending-icon {
                width: 54px;
                height: 54px;
                margin-bottom: 14px;
                display: flex;
                align-items: center;
                justify-content: center;
                border-radius: 18px;
                color: ${token.colorPrimary};
                background: ${token.colorPrimaryBg};
                font-size: 26px;
            }

            .title {
                color: ${token.colorText};
                font-size: ${isCompactMode ? 18 : 20}px;
                font-weight: 700;
                line-height: 28px;
            }

            .desc {
                max-width: 440px;
                margin-top: 8px;
                color: ${token.colorTextTertiary};
                font-size: 13px;
                line-height: 22px;
            }
        `,
        done: css`
            display: grid;
            gap: ${isCompactMode ? 14 : 16}px;
        `,
        doneHead: css`
            padding: ${isCompactMode ? "14px" : "16px"};
            display: grid;
            gap: 10px;
            border-radius: ${token.borderRadius}px;
            border: 1px solid ${token.colorBorderSecondary};
            background: ${token.colorBgContainer};

            .done-icon {
                width: ${isCompactMode ? 22 : 24}px;
                height: ${isCompactMode ? 22 : 24}px;
                flex: none;
                display: flex;
                align-items: center;
                justify-content: center;
                color: ${token.colorSuccess};
                font-size: ${isCompactMode ? 20 : 22}px;
            }

            .done-row {
                display: flex;
                align-items: center;
                justify-content: space-between;
                gap: 12px;
            }

            .done-main {
                min-width: 0;
                display: flex;
                align-items: center;
                gap: 10px;
            }

            .done-title {
                color: ${token.colorText};
                font-size: ${isCompactMode ? 16 : 17}px;
                font-weight: 600;
                line-height: 24px;
            }

            .done-type {
                flex: none;
                padding: 2px 8px;
                border-radius: ${token.borderRadiusSM}px;
                color: ${token.colorSuccess};
                background: ${token.colorSuccessBg};
                font-size: 12px;
                line-height: 18px;
            }

            .done-desc {
                color: ${token.colorTextSecondary};
                font-size: 13px;
                line-height: 20px;
            }

            ${responsive.mobile} {
                .done-row {
                    align-items: center;
                    flex-direction: row;
                }
            }
        `,
        infoGrid: css`
            display: grid;
            grid-template-columns: repeat(2, minmax(0, 1fr));
            gap: 10px;

            ${responsive.mobile} {
                grid-template-columns: 1fr;
            }
        `,
        infoItem: css`
            min-width: 0;
            padding: ${isCompactMode ? "11px 12px" : "12px 14px"};
            border: 1px solid ${token.colorBorderSecondary};
            border-radius: ${token.borderRadius}px;
            background: ${token.colorBgContainer};

            .label {
                margin-bottom: 5px;
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 18px;
            }

            .value {
                color: ${token.colorText};
                font-size: 13px;
                line-height: 20px;
                word-break: break-all;
            }
        `,
    };
});

export default (): React.ReactNode => {
    const {message, modal} = App.useApp();
    const theme = useTheme();
    const [form] = Form.useForm();
    const pollTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
    const pollCountRef = useRef(0);
    const pollStoppedRef = useRef(false);
    const pollVersionRef = useRef(0);
    const [loading, setLoading] = useState(true);
    const [submitLoading, setSubmitLoading] = useState(false);
    const [verifyLoading, setVerifyLoading] = useState(false);
    const [metaReady, setMetaReady] = useState(false);
    const [info, setInfo] = useState<any>({});
    const [authType, setAuthType] = useState(REAL_AUTH_TYPE.PERSON);
    const [formStep, setFormStep] = useState(0);
    const [certify, setCertify] = useState<any>();

    const isCompactMode = theme?.isCompactTheme?.() || false;
    const isAuthed = !!info?.status;
    const currentAuthTypeName = info?.company_name ? "企业实名" : "个人实名";
    const currentStep = isAuthed ? 3 : certify?.id ? 2 : formStep;
    const {styles} = useStyles({isCompactMode});

    const authedItems = useMemo(() => {
        const items = [
            {label: "认证类型", value: currentAuthTypeName},
            {label: "认证姓名", value: info?.card_name || "-"},
            {label: "证件号码", value: info?.card_no || "-"},
            {label: "认证时间", value: info?.auth_time || "-"},
        ];
        if (info?.company_name || info?.company_no) {
            items.splice(1, 0, {label: "企业名称", value: info?.company_name || "-"});
            items.splice(2, 0, {label: "企业编码", value: info?.company_no || "-"});
        }
        return items;
    }, [info, currentAuthTypeName]);

    const stopCertifyPoll = () => {
        pollStoppedRef.current = true;
        pollVersionRef.current += 1;
        if (pollTimerRef.current) {
            clearTimeout(pollTimerRef.current);
            pollTimerRef.current = null;
        }
    };

    const requestInfo = async (options: { showLoading?: boolean; silent?: boolean } = {}) => {
        if (options?.showLoading) {
            setLoading(true);
        }
        return new Promise<any>((resolve, reject) => {
            void certificate({
                body: {action: "info"},
                onSuccess: (res: any) => {
                    const data = res?.data || {};
                    setInfo(data);
                    if (data?.status) {
                        localStorage.removeItem(CERTIFY_CACHE_KEY);
                    }
                    resolve(data);
                },
                onFail: (res: any) => {
                    if (!options?.silent) {
                        message?.error(res?.message || "获取实名信息失败");
                    }
                    reject(res);
                },
                onFinally: () => {
                    if (options?.showLoading) {
                        setLoading(false);
                    }
                },
            });
        });
    };

    const requestVerify = async (certifyId?: string) => {
        if (!certifyId) {
            return Promise.reject({message: "认证ID不能为空"});
        }
        return new Promise<any>((resolve, reject) => {
            void certificate({
                body: {
                    action: "verify",
                    id: certifyId,
                },
                onSuccess: (res: any) => resolve(res?.data || {}),
                onFail: (res: any) => reject(res),
            });
        });
    };

    const loadInfo = async () => {
        try {
            const data = await requestInfo({showLoading: true});
            if (data?.status) {
                stopCertifyPoll();
                setCertify(undefined);
            }
        } catch (_) {
        }
    };

    const rollbackCertify = (content: string) => {
        stopCertifyPoll();
        setCertify(undefined);
        setFormStep(1);
        localStorage.removeItem(CERTIFY_CACHE_KEY);
        message?.error(content);
    };

    const startCertifyPoll = (certifyId?: string) => {
        if (!certifyId) {
            return;
        }
        stopCertifyPoll();
        const version = pollVersionRef.current + 1;
        pollVersionRef.current = version;
        pollStoppedRef.current = false;
        pollCountRef.current = 0;

        const poll = async () => {
            if (pollStoppedRef.current || pollVersionRef.current !== version) {
                return;
            }
            pollCountRef.current += 1;
            try {
                const data = await requestVerify(certifyId);
                if (pollStoppedRef.current || pollVersionRef.current !== version) {
                    return;
                }
                if (!data?.status) {
                    if (pollCountRef.current >= CERTIFY_POLL_MAX_COUNT) {
                        stopCertifyPoll();
                        message?.warning("认证状态查询已暂停，可稍后重新进入页面查看结果");
                        return;
                    }
                    if (data?.break === false) {
                        pollTimerRef.current = setTimeout(poll, CERTIFY_POLL_INTERVAL);
                        return;
                    }
                    rollbackCertify(data?.message || "认证状态查询失败");
                    return;
                }
                stopCertifyPoll();
                setCertify(undefined);
                localStorage.removeItem(CERTIFY_CACHE_KEY);
                try {
                    await requestInfo({silent: true});
                } catch (_) {
                }
                return;
            } catch (res: any) {
                if (pollStoppedRef.current || pollVersionRef.current !== version) {
                    return;
                }
                if (pollCountRef.current >= CERTIFY_POLL_MAX_COUNT) {
                    stopCertifyPoll();
                    message?.warning("认证状态查询已暂停，可稍后重新进入页面查看结果");
                    return;
                }
                rollbackCertify(res?.message || "认证状态查询失败");
            }
        };

        pollTimerRef.current = setTimeout(poll, CERTIFY_POLL_INTERVAL);
    };

    useEffect(() => {
        const cache = localStorage.getItem(CERTIFY_CACHE_KEY);
        if (cache) {
            try {
                const data = JSON.parse(cache);
                setCertify(data);
                if (data?.id) {
                    startCertifyPoll(data?.id);
                }
            } catch (_) {
                localStorage.removeItem(CERTIFY_CACHE_KEY);
            }
        }
        void loadInfo();
        return () => stopCertifyPoll();
    }, []);

    const openAuthUrl = (url?: string) => {
        if (!url) {
            message?.warning("未获取到认证链接");
            return;
        }
        window.open(url, "_blank", "noopener,noreferrer");
    };

    const handleSubmit = async () => {
        if (!metaReady || typeof window.getMetaInfo !== "function") {
            message?.warning("实名认证环境初始化中，请稍后再试");
            return;
        }
        const values = await form.validateFields();
        const isCompany = authType === REAL_AUTH_TYPE.COMPANY;
        const env = window.getMetaInfo();
        setSubmitLoading(true);
        await certificate({
            body: {
                action: "init",
                type: String(authType),
                card_name: values?.card_name,
                card_no: values?.card_no,
                company_name: isCompany ? values?.company_name : "",
                company_no: isCompany ? values?.company_no : "",
                env: typeof env === "string" ? env : JSON.stringify(env),
            },
            onSuccess: (res: any) => {
                const data = res?.data || {};
                setCertify(data);
                localStorage.setItem(CERTIFY_CACHE_KEY, JSON.stringify(data));
                if (data?.id) {
                    startCertifyPoll(data?.id);
                }
                if (data?.url) {
                    modal.confirm({
                        title: "继续实名认证",
                        content: "是否打开认证页面完成核验？",
                        okText: "确认",
                        cancelText: "取消",
                        mask: {closable: true},
                        onOk: () => openAuthUrl(data?.url),
                    } as any);
                }
            },
            onFail: (res: any) => {
                message?.error(res?.message || "发起实名失败");
            },
            onFinally: () => {
                setSubmitLoading(false);
            },
        });
    };

    const handleVerify = async () => {
        if (!certify?.id) {
            message?.warning("请先发起实名认证");
            return;
        }
        setVerifyLoading(true);
        try {
            const data = await requestVerify(certify?.id);
            if (!data?.status) {
                if (data?.break === false) {
                    message?.warning(data?.message || "认证未完成");
                    return;
                }
                rollbackCertify(data?.message || "实名验证失败");
                return;
            }
            message?.success(data?.message || "实名验证成功");
            stopCertifyPoll();
            setCertify(undefined);
            localStorage.removeItem(CERTIFY_CACHE_KEY);
            await loadInfo();
        } catch (res: any) {
            rollbackCertify(res?.message || "实名验证失败");
        } finally {
            setVerifyLoading(false);
        }
    };

    const resetCertify = () => {
        stopCertifyPoll();
        setCertify(undefined);
        setFormStep(0);
        form.resetFields();
        form.setFieldsValue({type: authType});
        localStorage.removeItem(CERTIFY_CACHE_KEY);
    };

    const backToCertifyForm = () => {
        stopCertifyPoll();
        setCertify(undefined);
        setFormStep(1);
        localStorage.removeItem(CERTIFY_CACHE_KEY);
    };

    return (
        <Body loading={loading}>
            <Script
                src={ALIYUN_META_SCRIPT}
                cache={true}
                timeout={10000}
                retryCount={2}
                onLoad={() => setMetaReady(true)}
                onError={() => setMetaReady(false)}
            />
            <Row className={styles.page}>
                <Card className={styles.shell} variant="borderless" styles={{body: {padding: 0}}}>
                    <div className={styles.hero}>
                        <div className={styles.heroMain}>
                            <div className="hero-title">实名认证</div>
                            <div className="desc">
                                依据真实身份完成账户核验，用于提升账户安全等级，并满足资金、订单等关键业务的合规校验要求。
                            </div>
                        </div>
                        <div className={styles.heroVisual}>
                            <div className="visual-info">
                                <div className="visual-label">当前认证</div>
                                <div className="visual-value">{isAuthed ? currentAuthTypeName : "待完成"}</div>
                                <div className="visual-desc">{isAuthed ? "实名信息已生效" : "按流程完成在线核验"}</div>
                            </div>
                            <div className="visual-icon">
                                <svg viewBox="0 0 96 76" fill="none" aria-hidden="true">
                                    <defs>
                                        <linearGradient id="realAuthCard" x1="14" y1="14" x2="74" y2="60" gradientUnits="userSpaceOnUse">
                                            <stop stopColor="currentColor" stopOpacity="0.12"/>
                                            <stop offset="1" stopColor="currentColor" stopOpacity="0.035"/>
                                        </linearGradient>
                                        <linearGradient id="realAuthSecurityShield" x1="57" y1="22" x2="84" y2="58" gradientUnits="userSpaceOnUse">
                                            <stop stopColor="currentColor" stopOpacity="0.28"/>
                                            <stop offset="0.58" stopColor="currentColor" stopOpacity="0.13"/>
                                            <stop offset="1" stopColor="currentColor" stopOpacity="0.05"/>
                                        </linearGradient>
                                        <linearGradient id="realAuthSecurityCore" x1="61" y1="28" x2="80" y2="54" gradientUnits="userSpaceOnUse">
                                            <stop stopColor="currentColor" stopOpacity="0.08"/>
                                            <stop offset="1" stopColor="currentColor" stopOpacity="0.02"/>
                                        </linearGradient>
                                        <filter id="realAuthShadow" x="4" y="8" width="86" height="62" filterUnits="userSpaceOnUse" colorInterpolationFilters="sRGB">
                                            <feDropShadow dx="0" dy="6" stdDeviation="5" floodColor="#172033" floodOpacity="0.12"/>
                                        </filter>
                                    </defs>
                                    <g filter="url(#realAuthShadow)">
                                        <rect x="10" y="13" width="62" height="48" rx="9" fill="url(#realAuthCard)" stroke="currentColor" strokeOpacity="0.42" strokeWidth="1.4"/>
                                        <rect x="15" y="18" width="52" height="38" rx="6" stroke="currentColor" strokeOpacity="0.12"/>
                                    </g>
                                    <path d="M19 28H43" stroke="currentColor" strokeWidth="2" strokeLinecap="round" opacity="0.45"/>
                                    <path d="M19 39H52" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" opacity="0.22"/>
                                    <path d="M19 49H38" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" opacity="0.14"/>
                                    <path
                                        d="M49 22C54 28 54 36 49 42M56 20C62 28 62 38 56 46M43 24C47 29 47 35 43 40"
                                        stroke="currentColor"
                                        strokeOpacity="0.18"
                                        strokeWidth="1.2"
                                        strokeLinecap="round"
                                    />
                                    <g>
                                        <path
                                            d="M70 20.5L84 26.2V37.4C84 48.6 78.2 56.1 70 59.5C61.8 56.1 56 48.6 56 37.4V26.2L70 20.5Z"
                                            fill="url(#realAuthSecurityShield)"
                                            stroke="currentColor"
                                            strokeOpacity="0.46"
                                            strokeWidth="1.6"
                                            strokeLinejoin="round"
                                        />
                                        <path
                                            d="M70 26.8L78.2 30.4V38.1C78.2 45.1 75.2 49.8 70 52.4C64.8 49.8 61.8 45.1 61.8 38.1V30.4L70 26.8Z"
                                            fill="url(#realAuthSecurityCore)"
                                            stroke="currentColor"
                                            strokeOpacity="0.16"
                                            strokeWidth="1"
                                            strokeLinejoin="round"
                                        />
                                        <path d="M64.2 39.8C64.2 36.2 66.6 33.7 70 33.7C73.4 33.7 75.8 36.2 75.8 39.8" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" opacity="0.58"/>
                                        <path d="M66.7 40.2C66.7 38.1 68 36.7 70 36.7C72 36.7 73.3 38.1 73.3 40.2" stroke="currentColor" strokeWidth="1.35" strokeLinecap="round" opacity="0.5"/>
                                        <path d="M69.9 39.8C68 41.7 68.2 44.8 70.8 46.5" stroke="currentColor" strokeWidth="1.35" strokeLinecap="round" opacity="0.64"/>
                                        <path d="M73.7 43.3C72.9 46.4 70.8 48.4 68.1 49.3" stroke="currentColor" strokeWidth="1.35" strokeLinecap="round" opacity="0.42"/>
                                        <path d="M65.2 44.3C66 47.4 68.1 49.4 70.9 50.5" stroke="currentColor" strokeWidth="1.35" strokeLinecap="round" opacity="0.32"/>
                                        <circle cx="70" cy="40" r="1.7" fill="currentColor" opacity="0.64"/>
                                    </g>
                                </svg>
                            </div>
                        </div>
                    </div>
                    <div className={styles.workspace}>
                        <div className={styles.guide}>
                            <div className="guide-title">认证流程</div>
                            <div className="guide-desc">按步骤完成后，系统会自动同步认证结果。</div>
                            <div className={styles.stepList}>
                                {STEP_ITEMS.map((item, index) => (
                                    <div
                                        key={item.title}
                                        className={`${styles.stepItem} ${currentStep === index ? "active" : ""} ${currentStep > index ? "done" : ""}`}
                                    >
                                        <div className="step-icon">
                                            <Icon type={item.icon}/>
                                        </div>
                                        <div>
                                            <div className="step-title">{item.title}</div>
                                            <div className="step-desc">{item.desc}</div>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </div>
                        <div className={styles.content}>
                            {isAuthed ? (
                                <div className={styles.done}>
                                    <div className={styles.doneHead}>
                                        <div className="done-row">
                                            <div className="done-main">
                                                <div className="done-icon">
                                                    <Icon type="CheckCircleFilled"/>
                                                </div>
                                                <div className="done-title">实名认证已完成</div>
                                            </div>
                                            <div className="done-type">{currentAuthTypeName}</div>
                                        </div>
                                        <div className="done-desc">认证结果已生效，当前账户已完成真实身份核验。</div>
                                    </div>
                                    <div className={styles.infoGrid}>
                                        {authedItems.map((item) => (
                                            <div className={styles.infoItem} key={item.label}>
                                                <div className="label">{item.label}</div>
                                                <div className="value">{item.value}</div>
                                            </div>
                                        ))}
                                    </div>
                                </div>
                            ) : (
                                <Form form={form} layout="vertical" initialValues={{type: REAL_AUTH_TYPE.PERSON}}>
                                    {!certify?.id ? (
                                        <>
                                            <Form.Item name="type" hidden>
                                                <Input/>
                                            </Form.Item>
                                            {formStep === 0 ? (
                                                <div className={styles.formContent}>
                                                    <div className={styles.sectionHead}>
                                                        <div>
                                                            <div className="section-title">选择认证类型</div>
                                                            <div
                                                                className="section-desc">请选择本次需要核验的主体类型。
                                                            </div>
                                                        </div>
                                                    </div>
                                                    <div className={styles.typeGrid}>
                                                        {AUTH_TYPE_OPTIONS.map((item) => (
                                                            <button
                                                                key={item.value}
                                                                type="button"
                                                                className={`${styles.typeItem} ${authType === item.value ? "active" : ""}`}
                                                                onClick={() => {
                                                                    setAuthType(item.value);
                                                                    form.setFieldsValue({type: item.value});
                                                                }}
                                                            >
                                                                <div className="type-top">
                                                                    <div className="type-icon">
                                                                        <Icon type={item.icon}/>
                                                                    </div>
                                                                    <span className="type-check"/>
                                                                </div>
                                                                <div>
                                                                    <div className="title">{item.label}</div>
                                                                    <div className="desc">{item.desc}</div>
                                                                </div>
                                                            </button>
                                                        ))}
                                                    </div>
                                                    <div className={`${styles.formActionBar} single`}>
                                                        <Button type="primary" onClick={() => setFormStep(1)}>
                                                            下一步
                                                        </Button>
                                                    </div>
                                                </div>
                                            ) : (
                                                <div className={styles.formContent}>
                                                    <div className={styles.sectionHead}>
                                                        <div>
                                                            <div className="section-title">填写认证资料</div>
                                                            <div className="section-desc">
                                                                请填写与证件一致的信息，提交后进入在线核验。
                                                            </div>
                                                        </div>
                                                    </div>
                                                    {authType === REAL_AUTH_TYPE.COMPANY && (
                                                        <>
                                                            <Form.Item
                                                                name="company_name"
                                                                label="企业名称"
                                                                rules={[{required: true, message: "请输入企业名称"}]}
                                                            >
                                                                <Input placeholder="请输入企业名称" maxLength={100}/>
                                                            </Form.Item>
                                                            <Form.Item
                                                                name="company_no"
                                                                label="企业编码"
                                                                rules={[{required: true, message: "请输入企业编码"}]}
                                                            >
                                                                <Input placeholder="请输入统一社会信用代码"
                                                                       maxLength={50}/>
                                                            </Form.Item>
                                                        </>
                                                    )}
                                                    <Form.Item
                                                        name="card_name"
                                                        label="身份证姓名"
                                                        rules={[{required: true, message: "请输入身份证姓名"}]}
                                                    >
                                                        <Input placeholder="请输入身份证姓名" maxLength={20}/>
                                                    </Form.Item>
                                                    <Form.Item
                                                        name="card_no"
                                                        label="身份证号码"
                                                        rules={[
                                                            {required: true, message: "请输入身份证号码"},
                                                            {min: 15, max: 18, message: "身份证号码长度为15-18位"},
                                                        ]}
                                                    >
                                                        <Input placeholder="请输入身份证号码" maxLength={18}/>
                                                    </Form.Item>
                                                    <div className={styles.formActionBar}>
                                                        <Button onClick={() => setFormStep(0)}>
                                                            上一步
                                                        </Button>
                                                        <Button type="primary" loading={submitLoading}
                                                                onClick={handleSubmit}>
                                                            下一步
                                                        </Button>
                                                    </div>
                                                </div>
                                            )}
                                        </>
                                    ) : (
                                        <>
                                            {certify?.url ? (
                                                <div className={styles.verifyPanel}>
                                                    <div className={styles.verifyHead}>
                                                        <div className="title">手机扫码认证</div>
                                                        <div className="desc">
                                                            扫码完成核验，结果将自动同步。
                                                        </div>
                                                    </div>
                                                    <div className={styles.qrCard}>
                                                        <QRCode
                                                            value={certify?.url}
                                                            size={isCompactMode ? 168 : 180}
                                                            bordered={false}
                                                        />
                                                    </div>
                                                    <div className={styles.verifyActions}>
                                                        <Button block type="primary"
                                                                icon={<Icon type="ExportOutlined"/>}
                                                                onClick={() => openAuthUrl(certify?.url)}>
                                                            打开认证页面
                                                        </Button>
                                                        <div className="verify-extra">
                                                            <span>认证完成后将自动同步结果</span>
                                                            <Button type="link" onClick={backToCertifyForm}>
                                                                返回上一步
                                                            </Button>
                                                        </div>
                                                    </div>
                                                </div>
                                            ) : (
                                                <div className={styles.pending}>
                                                    <div className="pending-icon">
                                                        <Icon type="LinkOutlined"/>
                                                    </div>
                                                    <div className="title">认证链接已生成</div>
                                                    <div className="desc">
                                                        请完成认证后返回当前页面查询认证结果。
                                                    </div>
                                                </div>
                                            )}
                                            {!certify?.url && (
                                                <div className={styles.actionBar}>
                                                    <Button onClick={resetCertify}>重新填写</Button>
                                                    <Button type="primary" onClick={handleVerify}
                                                            loading={verifyLoading}>
                                                        查询结果
                                                    </Button>
                                                </div>
                                            )}
                                        </>
                                    )}
                                </Form>
                            )}
                        </div>
                    </div>
                </Card>
            </Row>
        </Body>
    );
};

declare global {
    interface Window {
        getMetaInfo?: () => unknown;
    }
}
