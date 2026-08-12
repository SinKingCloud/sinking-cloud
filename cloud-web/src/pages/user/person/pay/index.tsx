import React, {useEffect, useState} from "react";
import {Card, Tabs} from "antd";
import {createStyles} from "antd-style";
import {Body} from "sinking-antd";
import {useLocation} from "umi";
import Balance from "./components/balance";
import Order from "./components/order";
import Recharge from "./components/recharge";

type PayTab = "recharge" | "balance" | "order";
const tabAlias: Record<string, PayTab> = {
    balance: "balance",
    log: "balance",
    recharge: "recharge",
    charge: "recharge",
    topup: "recharge",
    order: "order",
    orders: "order",
};

const getTabFromLocation = (location: any): PayTab => {
    const search = new URLSearchParams(location?.search || "");
    const hash = String(location?.hash || "").replace(/^#/, "");
    const queryTab = search.get("tab") || search.get("pay_tab") || hash;
    const target = tabAlias[String(queryTab || "").toLowerCase()];
    if (target) {
        return target;
    }
    const paths = String(location?.pathname || "").split("/").filter(Boolean);
    const payIndex = paths.lastIndexOf("pay");
    return tabAlias[String(paths?.[payIndex + 1] || "").toLowerCase()] || "balance";
};

const RechargePanel = React.memo((props: {embedded?: boolean}): React.ReactElement => <Recharge {...props}/>);

const useStyles = createStyles(({css, token}: any): any => ({
    page: css`
        width: 100%;
        max-width: 1450px;
        margin: 0 auto;
    `,
    navCard: css`
        margin-bottom: 12px;

        .ant-card-body {
            padding: 0 15px;
        }
    `,
    tabs: css`
        width: 100%;

        .ant-tabs-nav {
            margin: 0;
        }

        .ant-tabs-nav::before {
            display: none;
        }

        .ant-tabs-nav-list {
            gap: 18px;
        }

        .ant-tabs-tab {
            margin: 0 !important;
            padding: 15px 0;
            color: ${token.colorTextSecondary};
        }

        .ant-tabs-tab + .ant-tabs-tab {
            margin-left: 0;
        }

        .ant-tabs-tab-btn {
            font-size: ${token.fontSize}px;
            font-weight: 400;
            line-height: 20px;
        }

        .ant-tabs-ink-bar {
            height: 3px !important;
            border-radius: ${token.borderRadius}px;
        }

        .ant-tabs-content-holder {
            display: none;
        }
    `,
    tabPanel: css`
        display: none;

        &.active {
            display: block;
        }
    `,
}));

export default (): React.ReactNode => {
    const {styles} = useStyles();
    const location = useLocation();
    const [tab, setTab] = useState<PayTab>(() => getTabFromLocation(location));
    const [rechargeMounted, setRechargeMounted] = useState(() => getTabFromLocation(location) === "recharge");

    useEffect(() => {
        const nextTab = getTabFromLocation(location);
        setTab((current) => current === nextTab ? current : nextTab);
        if (nextTab === "recharge") {
            setRechargeMounted(true);
        }
    }, [location?.pathname, location?.search, location?.hash]);

    const changeTab = (value: string) => {
        const nextTab = value as PayTab;
        setTab(nextTab);
        if (nextTab === "recharge") {
            setRechargeMounted(true);
        }
    };

    return (
        <Body>
            <div className={styles.page}>
                <Card className={styles.navCard} variant="borderless">
                    <Tabs
                        className={styles.tabs}
                        activeKey={tab}
                        onChange={changeTab}
                        items={[
                            {key: "balance", label: "余额明细"},
                            {key: "recharge", label: "在线充值"},
                            {key: "order", label: "订单记录"},
                        ] as any}
                    />
                </Card>
                {tab === "balance" && <Balance embedded/>}
                <div className={`${styles.tabPanel} ${tab === "recharge" ? "active" : ""}`}>
                    {rechargeMounted && <RechargePanel embedded/>}
                </div>
                {tab === "order" && <Order embedded/>}
            </div>
        </Body>
    );
};
