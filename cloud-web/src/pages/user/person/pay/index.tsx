import React, {useEffect, useRef, useState} from "react";
import {App, Button, Card, ConfigProvider, Segmented, Spin, Tooltip, Typography} from "antd";
import {createStyles} from "antd-style";
import {Body, Icon, useTheme} from "sinking-antd";
import {useLocation} from "umi";
import {getBalanceInfo} from "@/service/user/pay";
import {formatMoney} from "@/utils/string";
import Balance from "./components/balance";
import Order from "./components/order";
import Recharge, {type RechargeRef} from "./components/recharge";

type PayTab = "recharge" | "balance" | "order";
const tabAlias: Record<string, PayTab> = {
    balance: "balance", log: "balance", recharge: "recharge", charge: "recharge",
    topup: "recharge", order: "order", orders: "order",
};

const getTabFromLocation = (location: any): PayTab => {
    const search = new URLSearchParams(location?.search || "");
    const hash = String(location?.hash || "").replace(/^#/, "");
    const queryTab = search.get("tab") || search.get("pay_tab") || hash;
    const target = tabAlias[String(queryTab || "").toLowerCase()];
    if (target) return target;
    const paths = String(location?.pathname || "").split("/").filter(Boolean);
    const payIndex = paths.lastIndexOf("pay");
    return tabAlias[String(paths?.[payIndex + 1] || "").toLowerCase()] || "balance";
};

const useStyles = createStyles(({css, token, isDarkMode}, {isCompactMode}: {isCompactMode: boolean}) => ({
    page: css`
        width: 100%;
        min-width: 0;
        max-width: 1450px;
        margin: 0 auto;
    `,
    hero: css`
        position: relative;
        overflow: hidden;
        background: ${isDarkMode
            ? `linear-gradient(135deg, color-mix(in srgb, ${token.colorPrimary} 9%, ${token.colorBgContainer}) 0%, ${token.colorBgContainer} 100%)`
            : token.colorFillAlter};
        isolation: isolate;
    `,
    header: css`
        position: relative;
        z-index: 1;
        display: grid;
        grid-template-columns: minmax(0, 1fr) ${isCompactMode ? 320 : 360}px;
        align-items: center;
        gap: ${isCompactMode ? 16 : 20}px;
        min-height: ${isCompactMode ? 159 : 185}px;
        padding: 0 ${isCompactMode ? 20 : 26}px;

        @media (max-width: 900px) {
            grid-template-columns: minmax(0, 1fr);
            min-height: auto;
            gap: ${isCompactMode ? 14 : 16}px;
            padding: ${isCompactMode ? 18 : 20}px ${isCompactMode ? 18 : 20}px ${isCompactMode ? 48 : 52}px;
        }
    `,
    headerMain: css`
        position: relative;
        z-index: 1;
        display: flex;
        flex-direction: column;
        gap: ${isCompactMode ? 8 : 10}px;
        min-width: 0;

        h1 {
            margin: 0;
            color: ${token.colorText};
            font-size: ${isCompactMode ? 26 : 28}px;
            font-weight: 700;
            line-height: ${isCompactMode ? 34 : 36}px;
            letter-spacing: 0;
        }
        p {
            max-width: 600px;
            margin: 0;
            color: ${token.colorTextSecondary};
            font-size: 13px;
            line-height: 22px;
        }
        @media (min-width: 901px) and (max-width: 1100px) {
            padding-bottom: ${isCompactMode ? 40 : 44}px;
        }
        @media (max-width: 560px) {
            h1 { font-size: ${isCompactMode ? 22 : 23}px; line-height: ${isCompactMode ? 30 : 31}px; }
            p { font-size: 12px; line-height: 20px; }
        }
    `,
    headerDecor: css`
        position: absolute;
        z-index: 0;
        inset: 0;
        overflow: hidden;
        pointer-events: none;

        .header-ambient, .header-shine, .header-watermark { position: absolute; pointer-events: none; }
        .header-ambient { inset: 0; }
        .header-ambient::before {
            content: '';
            position: absolute;
            inset: -78px;
            background:
                linear-gradient(135deg,
                    ${isDarkMode ? `color-mix(in srgb, ${token.colorPrimary} 42%, transparent)` : token.colorPrimaryBgHover} 0%,
                    ${isDarkMode ? `color-mix(in srgb, ${token.colorPrimary} 22%, transparent)` : token.colorPrimaryBg} 44%, transparent 76%),
                linear-gradient(30deg, transparent 18%,
                    ${isDarkMode ? "rgba(34,211,238,0.16)" : token.colorPrimaryBg} 72%, transparent 100%);
            opacity: ${isDarkMode ? 0.42 : 0.58};
            filter: blur(24px);
            transform: rotate(-7deg) translate3d(-18px, -12px, 0);
            -webkit-mask-image: linear-gradient(135deg, rgba(0,0,0,0.82), rgba(0,0,0,0.34) 58%, transparent 100%);
            mask-image: linear-gradient(135deg, rgba(0,0,0,0.82), rgba(0,0,0,0.34) 58%, transparent 100%);
        }
        .header-shine {
            inset: 0;
            opacity: ${isDarkMode ? 0.20 : 0.34};
            background: linear-gradient(110deg, transparent 16%,
                ${isDarkMode ? "rgba(34,211,238,0.18)" : "rgba(255,255,255,0.70)"} 58%, transparent 96%);
            mix-blend-mode: ${isDarkMode ? "screen" : "normal"};
        }
        .header-watermark {
            right: -78px;
            bottom: -104px;
            color: ${token.colorPrimary};
            font-size: 320px;
            line-height: 1;
            opacity: ${isDarkMode ? 0.17 : 0.14};
            transform: rotate(-24deg);
            filter: blur(20px);
        }
        @media (max-width: 900px) {
            .header-watermark { right: -62px; bottom: -90px; font-size: 260px; }
        }
        @media (max-width: 560px) {
            .header-ambient::before { inset: -54px; }
            .header-watermark { right: -48px; bottom: -68px; font-size: 210px; filter: blur(18px); }
        }
    `,
    balanceArea: css`
        position: relative;
        z-index: 1;
        width: 100%;
        min-width: 0;

        .balance-loading { position: absolute; inset: 0; display: grid; place-items: center; }
    `,
    balance: css`
        width: 100%;
        min-width: 0;
        overflow: hidden;
        background: ${isDarkMode ? "rgba(255,255,255,0.035)" : "rgba(255,255,255,0.72)"};

        &.ant-card:not(.ant-card-bordered) {
            border: none;
            border-radius: ${token.borderRadiusLG}px;
            box-shadow: ${isDarkMode ? "none" : "0 12px 28px rgba(15, 23, 42, 0.045)"};
        }

        .ant-card-body { padding: 0; }
        .balance-main { display: grid; grid-template-columns: minmax(0, 1fr) auto; align-items: center; gap: ${isCompactMode ? 12 : 16}px; padding: ${isCompactMode ? 12 : 18}px; }
        .balance-info { min-width: 0; }
        .balance-label { display: inline-flex; align-items: center; gap: ${isCompactMode ? 4 : 6}px; color: ${token.colorTextTertiary}; font-size: 12px; line-height: ${isCompactMode ? 16 : 18}px; }
        .balance-label .anticon { color: ${token.colorTextTertiary}; font-size: 11px; }
        .balance-value { display: flex; align-items: baseline; gap: 4px; min-width: 0; margin-top: ${isCompactMode ? 2 : 4}px; font-variant-numeric: tabular-nums; }
        .balance-value > span:first-child { flex: none; color: ${token.colorTextSecondary}; font-size: ${isCompactMode ? 12 : 13}px; }
        .balance-value .ant-typography { min-width: 0; color: ${token.colorText}; font-size: ${isCompactMode ? 20 : 24}px; font-weight: 600; line-height: ${isCompactMode ? 26 : 30}px; }
        .balance-recharge {
            justify-self: end;
            height: ${isCompactMode ? 30 : 36}px;
            padding-inline: ${isCompactMode ? 12 : 14}px;
            border-radius: ${token.borderRadiusLG}px;
            font-size: ${isCompactMode ? 12 : 13}px;
            font-weight: 500;
            transition: none;
            animation: none;

            &, &:not(:disabled):not(.ant-btn-disabled):is(:hover, :active) {
                border: 0;
                box-shadow: none;
                transform: none;
            }
        }
        .balance-summary {
            display: grid;
            grid-template-columns: repeat(3, minmax(0, 1fr));
            padding: ${isCompactMode ? 8 : 12}px ${isCompactMode ? 6 : 8}px;
            background:
                linear-gradient(110deg,
                    color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 12 : 6}%, transparent),
                    color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 5 : 2}%, transparent)),
                ${isDarkMode ? "rgba(255,255,255,0.035)" : "rgba(255,255,255,0.58)"};
            -webkit-backdrop-filter: blur(14px);
            backdrop-filter: blur(14px);
        }
        .balance-summary-item { position: relative; min-width: 0; padding-inline: ${isCompactMode ? 4 : 6}px; text-align: center; }
        .balance-summary-item + .balance-summary-item::before {
            content: '';
            position: absolute;
            inset-inline-start: 0;
            top: 50%;
            width: 1px;
            height: ${isCompactMode ? 18 : 22}px;
            background: color-mix(in srgb, ${token.colorTextSecondary} 18%, transparent);
            transform: translateY(-50%);
        }
        .balance-summary-label { display: block; margin-top: ${isCompactMode ? 2 : 3}px; color: ${token.colorTextTertiary}; font-size: 12px; line-height: ${isCompactMode ? 16 : 18}px; }
        .balance-summary-value.ant-typography { display: block; color: ${token.colorTextSecondary}; font-size: ${isCompactMode ? 12 : 13}px; font-weight: 500; line-height: ${isCompactMode ? 18 : 20}px; font-variant-numeric: tabular-nums; }
        @media (max-width: 380px) {
            .balance-main { gap: ${isCompactMode ? 6 : 8}px; padding: ${isCompactMode ? 10 : 14}px; }
            .balance-summary-item { padding-inline: ${isCompactMode ? 2 : 3}px; }
        }
    `,
    tabs: css`
        position: absolute;
        inset-inline-start: 0;
        bottom: 0;
        z-index: 2;
        max-width: 100%;
        font-size: 13px;

        &.ant-segmented {
            display: block;
            padding: 0 ${isCompactMode ? 10 : 12}px;
            border-radius: 0;
            background: transparent;
        }
        .ant-segmented-group { gap: 0; }
        .ant-segmented-item {
            border-radius: ${token.borderRadiusLG}px ${token.borderRadiusLG}px 0 0;
            background: transparent;
            box-shadow: none;
            color: ${token.colorTextSecondary};
            font-weight: 500;
        }
        .ant-segmented-item::after, .ant-segmented-thumb { display: none; }
        .ant-segmented-item-focused { outline-offset: -2px; }
        .ant-segmented-item-label {
            display: flex;
            align-items: center;
            justify-content: center;
            min-width: ${isCompactMode ? 88 : 96}px;
            min-height: ${isCompactMode ? 34 : 38}px;
            padding: 0 18px;
            line-height: 20px;
        }
        .ant-segmented-item-selected, .ant-segmented-item-selected-text {
            background: ${token.colorBgContainer};
            color: ${token.colorPrimary};
        }
        @container data-table-workspace (max-width: 680px) {
            &.ant-segmented { padding: 0 ${isCompactMode ? 8 : 11}px; }
        }
        @container data-table-workspace (max-width: 430px) {
            &.ant-segmented { padding: 0 ${isCompactMode ? 8 : 10}px; }
        }
    `,
}));

export default (): React.ReactNode => {
    const {message} = App.useApp();
    const theme = useTheme();
    const {styles} = useStyles({isCompactMode: theme?.isCompactTheme?.() || false});
    const location = useLocation();
    const rechargeRef = useRef<RechargeRef>(null);
    const [tab, setTab] = useState<"balance" | "order">(() => getTabFromLocation(location) === "order" ? "order" : "balance");
    const [balance, setBalance] = useState<any>();
    const [balanceLoading, setBalanceLoading] = useState(true);

    const refreshBalance = async () => {
        setBalanceLoading(true);
        try {
            const response = await getBalanceInfo();
            if (response?.code !== 200) {
                message.error(response?.message || "获取余额信息失败");
                return;
            }
            setBalance(response.data || {});
        } finally {
            setBalanceLoading(false);
        }
    };

    useEffect(() => { void refreshBalance(); }, []);

    useEffect(() => {
        const target = getTabFromLocation(location);
        setTab(target === "order" ? "order" : "balance");
        if (target === "recharge") rechargeRef.current?.show();
    }, [location?.pathname, location?.search, location?.hash]);

    const totalMoney = Number(balance?.money || 0) + Number(balance?.frozen || 0);
    const hero = <section className={styles.hero}>
        <div className={styles.headerDecor} aria-hidden="true">
            <span className="header-ambient"/>
            <span className="header-shine"/>
            <span className="header-watermark"><Icon type="WalletOutlined"/></span>
        </div>
        <div className={styles.header}>
            <div className={styles.headerMain}>
                <h1>财务管理</h1>
                <p>查看账户余额与收支明细，管理充值和支付订单。</p>
            </div>
            <div className={styles.balanceArea} aria-busy={balanceLoading}>
                <Card className={styles.balance} size="small" variant="borderless" aria-label="账户资金概览"
                      style={{visibility: balanceLoading ? "hidden" : "visible"}}>
                    <div className="balance-main">
                        <div className="balance-info">
                            <Tooltip title={balance ? `账户总余额 ¥ ${formatMoney(totalMoney)}` : "账户总余额"}>
                                <span className="balance-label">可用余额<Icon type="InfoCircleOutlined"/></span>
                            </Tooltip>
                            <div className="balance-value"><span>¥</span><Typography.Text ellipsis={{tooltip: balance ? formatMoney(balance.money) : "-"}}>{balance ? formatMoney(balance.money) : "-"}</Typography.Text></div>
                        </div>
                        <ConfigProvider wave={{disabled: true}}>
                            <Button className="balance-recharge" color="primary" variant="solid"
                                    onClick={() => rechargeRef.current?.show()}>在线充值</Button>
                        </ConfigProvider>
                    </div>
                    <div className="balance-summary">
                        {[
                            {field: "frozen", label: "冻结金额"},
                            {field: "total_recharge", label: "累计充值"},
                            {field: "total_consume", label: "累计消费"},
                        ].map((item) => <div className="balance-summary-item" key={item.field}>
                            <Typography.Text className="balance-summary-value" ellipsis={{tooltip: balance ? formatMoney(balance[item.field]) : "-"}}>{balance ? `¥ ${formatMoney(balance[item.field])}` : "-"}</Typography.Text>
                            <span className="balance-summary-label">{item.label}</span>
                        </div>)}
                    </div>
                </Card>
                {balanceLoading && <div className="balance-loading" role="status" aria-label="正在加载余额"><Spin/></div>}
            </div>
        </div>
        <Segmented className={styles.tabs} aria-label="财务记录" value={tab}
                   onChange={(value) => setTab(value as "balance" | "order")}
                   options={[{value: "balance", label: "余额明细"}, {value: "order", label: "订单记录"}]}/>
    </section>;

    return <Body>
        <div className={styles.page}>
            {tab === "balance" ? <Balance hero={hero}/> : <Order hero={hero}/>}
            <Recharge ref={rechargeRef} onRefreshBalance={refreshBalance}/>
        </div>
    </Body>;
};
