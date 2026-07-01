import React, {useEffect, useMemo, useState} from "react";
import {App, Button, Card, Col, Form, InputNumber, Radio, Row, Typography} from "antd";
import {createStyles} from "antd-style";
import {Body, Icon, Title, useTheme} from "sinking-antd";
import useEnum from "@/utils/enum";
import {formatMoney} from "@/utils/string";
import {setPayJumpUrl} from "@/utils/pay";
import {getBalanceInfo, getPayConfig, rechargeBalance, submitPay} from "@/service/user/pay";

const quickAmounts = [10, 30, 50, 100, 200];
const payIcon: Record<string, { icon: string; color?: string }> = {
    alipay: {icon: "AlipayOutlined", color: "#1677ff"},
    wxpay: {icon: "WechatOutlined", color: "#07c160"},
    balance: {icon: "RedEnvelopeOutlined", color: "#faad14"},
};

const useStyles = createStyles(({token, css, isDarkMode}: any, props: any = {}) => {
    const compact = props?.isCompactMode || false;
    const radius = token.borderRadiusLG;
    const gap = compact ? 12 : 14;
    const cardPadding = compact ? 16 : 18;
    const borderColor = isDarkMode ? "rgba(255,255,255,0.08)" : "rgba(15,23,42,0.06)";
    const surface = isDarkMode ? "rgba(255,255,255,0.04)" : "rgba(255,255,255,0.72)";
    const primarySoft = isDarkMode ? `color-mix(in srgb, ${token.colorPrimary} 10%, transparent)` : `color-mix(in srgb, ${token.colorPrimary} 7%, ${token.colorBgContainer})`;
    const primaryHover = isDarkMode ? "rgba(255,255,255,0.035)" : `color-mix(in srgb, ${token.colorPrimary} 4%, ${token.colorBgContainer})`;
    const glow = isDarkMode ? `color-mix(in srgb, ${token.colorPrimary} 28%, transparent)` : `${token.colorPrimary}18`;
    const heroBg = isDarkMode
        ? `linear-gradient(135deg, color-mix(in srgb, ${token.colorPrimary} 12%, ${token.colorBgContainer}) 0%, ${token.colorBgContainer} 100%)`
        : `linear-gradient(135deg, ${token.colorPrimaryBg} 0%, ${token.colorBgContainer} 58%, ${token.colorPrimaryBgHover} 100%)`;

    return {
        page: css`
            width: 100%;
            max-width: 1440px;
            margin: 0 auto;
        `,
        hero: css`
            position: relative;
            overflow: hidden;
            isolation: isolate;
            min-height: ${compact ? 120 : 130}px;
            border-radius: ${radius}px;
            padding: ${compact ? "16px 18px" : "18px 22px"};
            background: ${heroBg};

            &::before {
                content: "";
                position: absolute;
                right: -92px;
                top: -118px;
                width: 300px;
                height: 300px;
                border-radius: 999px;
                background: radial-gradient(circle, ${token.colorPrimaryBgHover} 0%, transparent 64%);
                filter: blur(30px);
                opacity: ${isDarkMode ? 0.16 : 0.40};
                pointer-events: none;
                z-index: 0;
            }

            .hero-ambient,
            .hero-watermark {
                position: absolute;
                pointer-events: none;
            }

            .hero-ambient {
                inset: 0;
                z-index: 0;
                opacity: ${isDarkMode ? 0.34 : 0.50};
                background:
                    radial-gradient(circle at 18% 16%, ${glow} 0%, transparent 38%),
                    radial-gradient(circle at 82% 78%, ${isDarkMode ? "rgba(34,211,238,0.12)" : token.colorPrimaryBgHover} 0%, transparent 42%);
                filter: blur(3px);
                transform: rotate(-7deg) translate3d(-18px, -12px, 0);
                mask-image: linear-gradient(120deg, rgba(0, 0, 0, 0.78), rgba(0, 0, 0, 0.34) 58%, transparent 100%);
            }

            .hero-watermark {
                right: -72px;
                bottom: -102px;
                z-index: 1;
                color: ${token.colorPrimary};
                font-size: ${compact ? 230 : 280}px;
                line-height: 1;
                opacity: ${isDarkMode ? 0.14 : 0.11};
                transform: rotate(-24deg);
                filter: blur(18px);
            }

            @media (max-width: 620px) {
                .hero-watermark {
                    right: -58px;
                    bottom: -78px;
                    font-size: ${compact ? 180 : 210}px;
                    opacity: ${isDarkMode ? 0.11 : 0.08};
                }
            }
        `,
        heroInner: css`
            position: relative;
            z-index: 1;
            display: grid;
            grid-template-columns: minmax(0, 1fr) minmax(240px, 280px);
            align-items: center;
            gap: ${compact ? 14 : 18}px;

            @media (max-width: 900px) {
                grid-template-columns: 1fr;
            }
        `,
        heroTitle: css`
            margin: 0;
            color: ${token.colorText};
            font-size: ${compact ? 24 : 27}px;
            line-height: ${compact ? "32px" : "35px"};
            font-weight: 700;
            letter-spacing: 0;

            @media (max-width: 620px) {
                font-size: ${compact ? 22 : 23}px;
                line-height: ${compact ? "30px" : "31px"};
            }
        `,
        heroDesc: css`
            max-width: 610px;
            margin-top: 7px;
            color: ${token.colorTextSecondary};
            font-size: ${compact ? 12 : 13}px;
            line-height: ${compact ? "21px" : "22px"};
        `,
        wallet: css`
            position: relative;
            overflow: hidden;
            padding: ${compact ? "12px 13px" : "13px 14px"};
            border: 1px solid ${borderColor};
            border-radius: ${radius}px;
            background: ${surface};
            backdrop-filter: blur(16px);

            .wallet-label {
                display: flex;
                align-items: center;
                justify-content: space-between;
                gap: 10px;
                color: ${token.colorTextSecondary};
                font-size: 12px;
                line-height: 20px;
            }

            .wallet-refresh {
                display: inline-flex;
                align-items: center;
                justify-content: center;
                width: 22px;
                height: 22px;
                margin: -2px;
                padding: 0;
                border: 0;
                border-radius: ${token.borderRadiusSM}px;
                background: transparent;
                color: ${token.colorTextTertiary};
                cursor: pointer;
                transition: color 0.18s ease, background-color 0.18s ease;
            }

            .wallet-refresh:hover {
                color: ${token.colorPrimary};
                background: ${isDarkMode ? "rgba(255,255,255,0.045)" : token.colorFillQuaternary};
            }

            .wallet-refresh:disabled {
                cursor: not-allowed;
                color: ${token.colorTextDisabled};
                background: transparent;
            }

            .wallet-refresh.loading {
                color: ${token.colorPrimary};
            }

            .wallet-refresh.loading .refresh-icon {
                animation: walletRefreshSpin 0.8s linear infinite;
            }

            @keyframes walletRefreshSpin {
                to {
                    transform: rotate(360deg);
                }
            }

            .money {
                display: block;
                margin-top: 4px;
                color: ${token.colorPrimary};
                font-size: ${compact ? 27 : 31}px;
                line-height: ${compact ? "35px" : "39px"};
                font-weight: 700;
                letter-spacing: 0;
                font-variant-numeric: tabular-nums;
            }

            .wallet-meta {
                display: flex;
                align-items: center;
                justify-content: space-between;
                gap: 12px;
                margin-top: ${compact ? 9 : 10}px;
                padding-top: ${compact ? 9 : 10}px;
                border-top: 1px solid ${token.colorBorderSecondary};
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 19px;
            }

            .wallet-meta strong {
                color: ${token.colorText};
                font-weight: 600;
                font-variant-numeric: tabular-nums;
            }
        `,
        workspace: css`
            width: 100%;

            .ant-col {
                min-width: 0;
            }
        `,
        workspaceLeft: css`
            @media (min-width: 992px) {
                padding-right: ${gap / 2}px;
            }
        `,
        workspaceRight: css`
            @media (min-width: 992px) {
                padding-left: ${gap / 2}px;
            }
        `,
        panel: css`
            overflow: hidden;
            border: 1px solid ${borderColor};
            border-radius: ${radius}px;
            background:
                linear-gradient(180deg, ${token.colorBgContainer} 0%, ${isDarkMode ? "rgba(255,255,255,0.025)" : "rgba(255,255,255,0.86)"} 100%);

            > .ant-card-body {
                padding: ${cardPadding}px;
            }

            .ant-form-item {
                margin-bottom: 0;
            }

            .ant-form-item:last-child {
                margin-bottom: 0;
            }
        `,
        block: css`
            & + & {
                margin-top: ${compact ? 16 : 18}px;
                padding-top: ${compact ? 16 : 18}px;
            }
        `,
        blockHead: css`
            display: flex;
            align-items: baseline;
            justify-content: space-between;
            gap: 12px;
            margin-bottom: ${compact ? 10 : 12}px;

            .title {
                color: ${token.colorText};
                font-size: ${compact ? 14 : 15}px;
                line-height: ${compact ? "21px" : "22px"};
                font-weight: 600;
            }

            .desc {
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 20px;
                text-align: right;
            }

            @media (max-width: 620px) {
                align-items: flex-start;
                flex-direction: column;
                gap: 3px;

                .desc {
                    text-align: left;
                }
            }
        `,
        amountGrid: css`
            display: grid;
            grid-template-columns: repeat(3, minmax(0, 1fr));
            gap: ${compact ? 9 : 10}px;
            max-width: 96%;
            margin: 0 auto;

            @media (max-width: 620px) {
                max-width: none;
                margin: 0;
                grid-template-columns: repeat(2, minmax(0, 1fr));
            }
        `,
        amountItem: css`
            position: relative;
            width: 100%;
            height: ${compact ? 58 : 62}px;
            padding: ${compact ? "8px 10px" : "9px 11px"};
            border: 1px solid ${token.colorBorderSecondary};
            border-radius: ${token.borderRadius}px;
            background: ${token.colorBgContainer};
            color: ${token.colorText};
            letter-spacing: 0;
            cursor: pointer;
            text-align: left;
            transition: border-color 0.18s ease, background-color 0.18s ease, color 0.18s ease;

            .amount-main {
                display: flex;
                align-items: center;
                height: 100%;
                color: ${token.colorText};
                font-size: ${compact ? 17 : 19}px;
                line-height: ${compact ? "23px" : "25px"};
                font-weight: 700;
                letter-spacing: 0;
                font-variant-numeric: tabular-nums;
                transition: color 0.18s ease;
            }

            &:hover {
                border-color: ${token.colorPrimaryBorder};
                background: ${primaryHover};
            }

            &.active {
                border-color: ${token.colorPrimaryBorder};
                background: ${primarySoft};
            }

            &.active .amount-main {
                color: ${token.colorPrimary};
            }

            &.active::after {
                content: "";
                position: absolute;
                top: 10px;
                right: 10px;
                width: 6px;
                height: 6px;
                border-radius: 999px;
                background: ${token.colorPrimary};
            }

            &:disabled {
                cursor: not-allowed;
                color: ${token.colorTextDisabled};
                background: ${token.colorBgContainerDisabled};
            }
        `,
        customAmountItem: css`
            position: relative;
            height: ${compact ? 58 : 62}px;
            padding: ${compact ? "7px 9px" : "8px 10px"};
            border: 1px solid ${token.colorBorderSecondary};
            border-radius: ${token.borderRadius}px;
            background: ${token.colorBgContainer};
            transition: border-color 0.18s ease, background-color 0.18s ease;

            .ant-input-number {
                width: 100%;
                height: 100%;
                border: 0;
                background: transparent;
                box-shadow: none;
                outline: none;
            }

            .ant-input-number:hover,
            .ant-input-number-outlined,
            .ant-input-number-outlined:hover,
            .ant-input-number-outlined:focus,
            .ant-input-number-outlined:focus-within,
            .ant-input-number-focused,
            .ant-input-number:focus,
            .ant-input-number:focus-within {
                border: 0;
                background: transparent;
                box-shadow: none;
                outline: none;
            }

            &:hover,
            &.active,
            &:focus-within {
                border-color: ${token.colorPrimaryBorder};
                background: ${primarySoft};
            }

            .ant-input-number-input {
                height: ${compact ? 40 : 44}px;
                padding: 0;
                color: ${token.colorText};
                font-size: ${compact ? 17 : 19}px;
                font-weight: 700;
                text-align: left;
                font-variant-numeric: tabular-nums;
            }

            .ant-input-number-input:focus {
                box-shadow: none;
                outline: none;
            }

            .ant-input-number-handler-wrap {
                display: none;
            }

            &.active .ant-input-number-input,
            &:focus-within .ant-input-number-input {
                color: ${token.colorPrimary};
            }
        `,
        payOptions: css`
            display: grid;
            grid-template-columns: repeat(3, minmax(0, 1fr));
            gap: ${compact ? 8 : 10}px;
            max-width: 96%;
            margin: 0 auto;

            .ant-radio-wrapper {
                display: block;
                width: 100%;
                margin: 0;
            }

            .ant-radio {
                display: none;
            }

            .ant-radio-wrapper > span:last-child {
                display: block;
                padding-inline-start: 0;
            }

            @media (max-width: 620px) {
                max-width: none;
                margin: 0;
                gap: ${compact ? 7 : 8}px;
                grid-template-columns: repeat(2, minmax(0, 1fr));
            }

            @media (max-width: 360px) {
                grid-template-columns: 1fr;
            }
        `,
        payOption: css`
            position: relative;
            display: flex;
            align-items: center;
            gap: 9px;
            height: ${compact ? 58 : 62}px;
            padding: ${compact ? "8px 10px" : "9px 11px"};
            border: 1px solid ${token.colorBorderSecondary};
            border-radius: ${token.borderRadius}px;
            background: ${token.colorBgContainer};
            transition: border-color 0.18s ease, background-color 0.18s ease;

            .icon {
                display: inline-flex;
                align-items: center;
                justify-content: center;
                flex: 0 0 auto;
                width: 20px;
                height: 20px;
                font-size: 16px;
                color: ${token.colorTextTertiary};
                transition: color 0.18s ease;
            }

            .name {
                min-width: 0;
                color: ${token.colorText};
                font-size: 13px;
                line-height: 20px;
                font-weight: 600;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
                transition: color 0.18s ease;
            }

            .check {
                display: none;
                margin-left: auto;
                color: ${token.colorPrimary};
                font-size: 13px;
            }

            &:hover {
                border-color: ${token.colorPrimaryBorder};
                background: ${primaryHover};
            }

            &.active {
                border-color: ${token.colorPrimaryBorder};
                background: ${primarySoft};
            }

            &.active .name {
                color: ${token.colorPrimary};
            }

            &.active .check {
                display: inline-flex;
            }

            @media (max-width: 620px) {
                gap: 7px;
                height: ${compact ? 52 : 54}px;
                padding: ${compact ? "7px 9px" : "8px 10px"};
            }
        `,
        checkout: css`
            position: sticky;
            top: ${gap}px;
            overflow: hidden;
            border: 1px solid ${borderColor};
            border-radius: ${radius}px;
            background:
                linear-gradient(180deg, ${token.colorBgContainer} 0%, ${isDarkMode ? "rgba(255,255,255,0.025)" : "rgba(255,255,255,0.88)"} 100%);

            @media (max-width: 930px) {
                position: static;
            }

            > .ant-card-body {
                padding: 0;
            }
        `,
        checkoutBody: css`
            padding: ${cardPadding}px;
        `,
        totalPanel: css`
            display: grid;
            gap: 0;
            padding-bottom: ${compact ? 9 : 11}px;
            border-bottom: 1px solid ${token.colorBorderSecondary};

            .label {
                color: ${token.colorTextSecondary};
                font-size: 12px;
                line-height: 20px;
            }

            .amount {
                color: ${token.colorPrimary};
                font-size: ${compact ? 32 : 36}px;
                line-height: ${compact ? "36px" : "40px"};
                font-weight: 760;
                letter-spacing: 0;
                text-align: right;
                font-variant-numeric: tabular-nums;
            }
        `,
        checkoutRows: css`
            display: grid;
            gap: ${compact ? 9 : 10}px;
            padding: ${compact ? "14px 0" : "16px 0"};
            border-bottom: 1px solid ${token.colorBorderSecondary};
        `,
        checkoutRow: css`
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 14px;

            .label {
                flex: 0 0 auto;
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 20px;
            }

            .value {
                min-width: 0;
                color: ${token.colorText};
                font-size: 13px;
                line-height: 20px;
                text-align: right;
            }
        `,
        checkoutNotice: css`
            margin-top: ${compact ? 12 : 14}px;
            padding: ${compact ? "9px 10px" : "10px 11px"};
            border-radius: ${token.borderRadius}px;
            color: ${token.colorTextSecondary};
            font-size: 12px;
            line-height: 20px;
            background: ${isDarkMode ? "rgba(255,255,255,0.035)" : token.colorFillQuaternary};
        `,
        actionBar: css`
            display: grid;
            gap: ${compact ? 9 : 10}px;
            margin-top: ${compact ? 14 : 16}px;
        `,
        payButton: css`
            height: ${compact ? 38 : 40}px;
            font-weight: 600;
        `,
    };
});

const PayOption = ({value, label, active, className}: { value: string; label: string; active: boolean; className: string }) => {
    const meta = payIcon[value] || {icon: "CreditCardOutlined"};
    return (
        <div className={`${className} ${active ? "active" : ""}`}>
            <span className="icon" style={{color: meta.color}}>
                <Icon type={meta.icon}/>
            </span>
            <span className="name">{label}</span>
            <span className="check">
                <Icon type="CheckOutlined"/>
            </span>
        </div>
    );
};

export default (): React.ReactNode => {
    const {message, modal} = App.useApp();
    const [form] = Form.useForm();
    const theme = useTheme();
    const isCompactMode = theme?.isCompactTheme?.() || false;
    const {styles} = useStyles({isCompactMode});
    const [enumData, enumLoading] = useEnum(["pay.pay"]) as any;
    const [config, setConfig] = useState<any>({});
    const [balance, setBalance] = useState<any>({});
    const [loading, setLoading] = useState(true);
    const [balanceLoading, setBalanceLoading] = useState(false);
    const [submitting, setSubmitting] = useState(false);
    const [money, setMoney] = useState<number>(quickAmounts[0]);
    const [customMoney, setCustomMoney] = useState<number>();
    const [payType, setPayType] = useState<string>();
    const payTypes = enumData?.["pay.pay"]?.type || {};

    const enabledPayOptions = useMemo(() => Object.entries(payTypes)
        .filter(([value]) => value !== "balance" && config?.[value])
        .map(([value, label]) => ({value, label: String(label)})), [config, payTypes]);

    const amountRange = useMemo(() => {
        const min = Number(config?.min || 0);
        const max = Number(config?.max || 0);
        if (min > 0 && max > 0) {
            return `${formatMoney(min)} - ${formatMoney(max)} 元`;
        }
        if (min > 0) {
            return `不少于 ${formatMoney(min)} 元`;
        }
        if (max > 0) {
            return `不超过 ${formatMoney(max)} 元`;
        }
        return "不限制";
    }, [config]);
    const currentPayText = enabledPayOptions.find((item) => item.value === payType)?.label || "-";

    const refreshBalance = async () => {
        setBalanceLoading(true);
        try {
            const res = await getBalanceInfo({
                onFail: (failRes: any) => message?.error(failRes?.message || "获取余额信息失败"),
            });
            if (res?.code === 200) {
                setBalance(res?.data || {});
            }
        } finally {
            setBalanceLoading(false);
        }
    };

    const loadData = async () => {
        setLoading(true);
        try {
            const [configRes, balanceRes] = await Promise.all([
                getPayConfig({
                    onFail: (res: any) => message?.error(res?.message || "获取支付配置失败"),
                }),
                getBalanceInfo({
                    onFail: (res: any) => message?.error(res?.message || "获取余额信息失败"),
                }),
            ]);
            if (configRes?.code === 200) {
                setConfig(configRes?.data || {});
            }
            if (balanceRes?.code === 200) {
                setBalance(balanceRes?.data || {});
            }
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        void loadData();
    }, []);

    useEffect(() => {
        if (enabledPayOptions.length > 0 && (!payType || !enabledPayOptions.some((item) => item.value === payType))) {
            const nextPayType = enabledPayOptions[0].value;
            setPayType(nextPayType);
            form.setFieldValue("type", nextPayType);
        }
    }, [enabledPayOptions, form, payType]);

    const handleSubmit = async (values: any) => {
        const amount = Number(money || 0);
        const min = Number(config?.min || 0);
        const max = Number(config?.max || 0);
        const type = values?.type || payType;
        if (amount <= 0) {
            message?.warning("请输入充值金额");
            return;
        }
        if (min > 0 && amount < min) {
            message?.warning(`充值金额不能小于 ${formatMoney(min)} 元`);
            return;
        }
        if (max > 0 && amount > max) {
            message?.warning(`充值金额不能大于 ${formatMoney(max)} 元`);
            return;
        }
        if (!type) {
            message?.warning("请选择支付方式");
            return;
        }
        const openInCurrentPage = typeof navigator !== "undefined" && /iPhone|iPad|iPod|Macintosh/.test(navigator.userAgent);
        let payWindow: Window | null = null;
        if (!openInCurrentPage && typeof window !== "undefined") {
            payWindow = window.open("about:blank", "_blank");
        }
        const closePayWindow = () => {
            if (payWindow && !payWindow.closed) {
                payWindow.close();
            }
        };
        setSubmitting(true);
        try {
            const createRes = await rechargeBalance({
                body: {
                    money: amount,
                    type,
                },
                onFail: (res: any) => message?.error(res?.message || "创建充值订单失败"),
            });
            if (createRes?.code !== 200) {
                closePayWindow();
                return;
            }
            const submitRes = await submitPay({
                body: {id: Number(createRes?.data)},
                onFail: (res: any) => message?.error(res?.message || "获取支付链接失败"),
            });
            if (submitRes?.code !== 200) {
                closePayWindow();
                return;
            }
            const payUrl = submitRes?.data;
            if (!payUrl) {
                closePayWindow();
                message?.error("获取支付链接失败");
                return;
            }
            setPayJumpUrl(window.location.href);
            if (payWindow && !payWindow.closed) {
                payWindow.location.href = payUrl;
            } else if (openInCurrentPage) {
                window.location.href = payUrl;
                return;
            } else {
                window.open(payUrl, "_blank");
            }
            modal.confirm({
                title: "支付确认",
                content: "支付完成后可刷新余额信息，支付窗口关闭后也可以重新打开支付页面。",
                okText: "我已支付",
                cancelText: "重新支付",
                mask: {
                    closable: true,
                },
                onOk: refreshBalance,
                onCancel: () => {
                    if (typeof window !== "undefined" && payUrl) {
                        window.open(payUrl, "_blank");
                    }
                },
            } as any);
        } finally {
            setSubmitting(false);
        }
    };

    const pageLoading = loading || enumLoading;
    const formDisabled = pageLoading;
    const payDisabled = pageLoading || enabledPayOptions.length <= 0;

    return (
        <Body loading={pageLoading}>
            <Row className={styles.page} gutter={[0, isCompactMode ? 12 : 14]}>
                <Col span={24}>
                    <Card className={styles.hero} variant="borderless" styles={{body: {padding: 0}}}>
                        <span className="hero-ambient"/>
                        <span className="hero-watermark">
                            <Icon type="CreditCardOutlined"/>
                        </span>
                        <div className={styles.heroInner}>
                            <div>
                                <h1 className={styles.heroTitle}>在线充值</h1>
                                <div className={styles.heroDesc}>
                                    选择充值金额与支付渠道，系统将创建充值订单并引导至支付页面完成付款。
                                </div>
                            </div>
                            <div className={styles.wallet}>
                                <div className="wallet-label">
                                    <span>账户可用余额</span>
                                    <button
                                        type="button"
                                        className={`wallet-refresh ${balanceLoading ? "loading" : ""}`}
                                        disabled={loading || balanceLoading}
                                        onClick={refreshBalance}>
                                        <Icon type="ReloadOutlined" className="refresh-icon"/>
                                    </button>
                                </div>
                                <Typography.Text className="money" ellipsis={{tooltip: formatMoney(balance?.money)}}>
                                    {formatMoney(balance?.money)}
                                </Typography.Text>
                                <div className="wallet-meta">
                                    <span>冻结金额</span>
                                    <strong>{formatMoney(balance?.frozen)}</strong>
                                </div>
                            </div>
                        </div>
                    </Card>
                </Col>
                <Col span={24}>
                    <Row className={styles.workspace} gutter={[0, isCompactMode ? 12 : 14]} align="top">
                        <Col className={styles.workspaceLeft} xs={24} md={24} lg={16} xl={17}>
                            <Card className={styles.panel} title={<Title>在线充值</Title>} variant="borderless">
                                <Form form={form} layout="vertical" onFinish={handleSubmit}>
                                    <div className={styles.block}>
                                        <div className={styles.blockHead}>
                                            <span className="title">充值金额</span>
                                            <span className="desc">选择本次充值额度</span>
                                        </div>
                                        <div className={styles.amountGrid}>
                                            {quickAmounts.map((item) => (
                                                <button
                                                    key={item}
                                                    type="button"
                                                    className={`${styles.amountItem} ${customMoney === undefined && Number(money) === item ? "active" : ""}`}
                                                    disabled={formDisabled}
                                                    onClick={() => {
                                                        setCustomMoney(undefined);
                                                        setMoney(item);
                                                    }}>
                                                    <span className="amount-main">{item} 元</span>
                                                </button>
                                            ))}
                                            <div className={`${styles.customAmountItem} ${customMoney !== undefined ? "active" : ""}`}>
                                                <InputNumber
                                                    value={customMoney}
                                                    min={1}
                                                    precision={0}
                                                    step={1}
                                                    placeholder="自定义"
                                                    disabled={formDisabled}
                                                    onChange={(value) => {
                                                        const nextMoney = Number(value || 0);
                                                        setCustomMoney(value === null ? undefined : nextMoney);
                                                        setMoney(nextMoney);
                                                    }}
                                                />
                                            </div>
                                        </div>
                                    </div>
                                    <div className={styles.block}>
                                        <div className={styles.blockHead}>
                                            <span className="title">支付方式</span>
                                            <span className="desc">支付成功后自动同步到账</span>
                                        </div>
                                        <Form.Item name="type" rules={[{required: true, message: "请选择支付方式"}]}>
                                            <Radio.Group
                                                style={{width: "100%"}}
                                                disabled={formDisabled}
                                                value={payType}
                                                onChange={(event) => setPayType(event?.target?.value)}>
                                                <div className={styles.payOptions}>
                                                    {enabledPayOptions.map((item) => (
                                                        <Radio key={item.value} value={item.value}>
                                                            <PayOption
                                                                className={styles.payOption}
                                                                value={item.value}
                                                                label={item.label}
                                                                active={payType === item.value}
                                                            />
                                                        </Radio>
                                                    ))}
                                                </div>
                                            </Radio.Group>
                                        </Form.Item>
                                        {enabledPayOptions.length <= 0 && !pageLoading && (
                                            <Typography.Text type="secondary">暂未开启可用支付方式，请稍后再试。</Typography.Text>
                                        )}
                                    </div>
                                </Form>
                            </Card>
                        </Col>
                        <Col className={styles.workspaceRight} xs={24} md={24} lg={8} xl={7}>
                            <Card className={styles.checkout} title={<Title>订单确认</Title>} variant="borderless">
                                <div className={styles.checkoutBody}>
                                    <div className={styles.totalPanel}>
                                        <span className="label">应付金额</span>
                                        <Typography.Text className="amount" ellipsis>
                                            {formatMoney(money)}
                                        </Typography.Text>
                                    </div>
                                    <div className={styles.checkoutRows}>
                                        <div className={styles.checkoutRow}>
                                            <span className="label">支付方式</span>
                                            <Typography.Text className="value" ellipsis={{tooltip: currentPayText}}>
                                                {currentPayText}
                                            </Typography.Text>
                                        </div>
                                        <div className={styles.checkoutRow}>
                                            <span className="label">充值范围</span>
                                            <Typography.Text className="value" ellipsis={{tooltip: amountRange}}>
                                                {amountRange}
                                            </Typography.Text>
                                        </div>
                                        <div className={styles.checkoutRow}>
                                            <span className="label">有效时间</span>
                                            <span className="value">{Number(config?.expire || 0) || 5} 分钟</span>
                                        </div>
                                        <div className={styles.checkoutRow}>
                                            <span className="label">到账方式</span>
                                            <span className="value">支付完成后自动入账</span>
                                        </div>
                                    </div>
                                    <div className={styles.checkoutNotice}>
                                        提交后将在新标签页打开支付页面，充值状态可在订单管理和财务明细中查看。
                                    </div>
                                    <div className={styles.actionBar}>
                                        <Button className={styles.payButton} type="primary" block loading={submitting} disabled={payDisabled} onClick={() => form.submit()}>
                                            确认支付
                                        </Button>
                                    </div>
                                </div>
                            </Card>
                        </Col>
                    </Row>
                </Col>
            </Row>
        </Body>
    );
};
