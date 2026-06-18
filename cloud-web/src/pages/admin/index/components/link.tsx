import React from "react";
import {Card} from "antd";
import {RightOutlined, SettingOutlined, ShoppingCartOutlined, TeamOutlined, UserOutlined} from "@ant-design/icons";
import {createStyles} from "antd-style";
import {historyPush} from "@/utils/route";
import {AdminPageCode} from "@/utils/permission";

type QuickItem = {
    title: string;
    desc: string;
    route: string;
    icon: React.ReactNode;
};

const useStyles: any = createStyles(({token, isDarkMode}: any): any => ({
    card: {
        ".ant-card-body": {
            padding: 0,
        },
    },
    grid: {
        display: "grid",
        gridTemplateColumns: "repeat(4, minmax(0, 1fr))",
        gap: 0,
        overflow: "hidden",
        borderRadius: token.borderRadiusLG,
        "@media (max-width: 1080px)": {
            gridTemplateColumns: "repeat(2, minmax(0, 1fr))",
        },
        "@media (max-width: 576px)": {
            gridTemplateColumns: "1fr",
        },
    },
    item: {
        display: "flex",
        alignItems: "center",
        gap: "10px",
        minHeight: "72px",
        padding: "12px 14px",
        borderRight: `1px solid ${isDarkMode ? "rgba(255,255,255,0.07)" : token.colorBorderSecondary}`,
        background: isDarkMode ? "rgba(255,255,255,0.015)" : token.colorBgContainer,
        cursor: "pointer",
        transition: "background .18s ease",
        "&:last-child": {
            borderRight: "none",
        },
        "&:hover": {
            background: isDarkMode ? "rgba(255,255,255,0.04)" : token.colorFillQuaternary,
        },
        ".quick-icon": {
            width: "28px",
            height: "28px",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            flexShrink: 0,
            color: token.colorPrimary,
            fontSize: "21px",
        },
        ".quick-main": {
            flex: 1,
            minWidth: 0,
        },
        ".quick-title": {
            display: "block",
            fontSize: "13px",
            fontWeight: 600,
            color: token.colorTextHeading,
            lineHeight: "20px",
        },
        ".quick-desc": {
            display: "block",
            marginTop: "4px",
            fontSize: "12px",
            color: token.colorTextTertiary,
            lineHeight: "18px",
            overflow: "hidden",
            textOverflow: "ellipsis",
            whiteSpace: "nowrap",
        },
        ".quick-arrow": {
            color: token.colorTextQuaternary,
            flexShrink: 0,
            transition: "color .18s ease, transform .18s ease",
        },
        "&:hover .quick-arrow": {
            color: token.colorPrimary,
            transform: "translateX(2px)",
        },
        "@media (max-width: 1080px)": {
            "&:nth-child(2n)": {
                borderRight: "none",
            },
            "&:nth-child(-n + 2)": {
                borderBottom: `1px solid ${isDarkMode ? "rgba(255,255,255,0.07)" : token.colorBorderSecondary}`,
            },
        },
        "@media (max-width: 576px)": {
            borderRight: "none",
            borderBottom: `1px solid ${isDarkMode ? "rgba(255,255,255,0.07)" : token.colorBorderSecondary}`,
            "&:last-child": {
                borderBottom: "none",
            },
        },
    },
}));

const quickItems: QuickItem[] = [
    {
        title: "用户管理",
        desc: "用户资料、实名、会话",
        route: AdminPageCode.WEB_USER,
        icon: <UserOutlined/>,
    },
    {
        title: "订单管理",
        desc: "订单详情、发货、退款",
        route: AdminPageCode.WEB_ORDER,
        icon: <ShoppingCartOutlined/>,
    },
    {
        title: "员工管理",
        desc: "员工资料、角色、会话",
        route: AdminPageCode.SYSTEM_MASTER,
        icon: <TeamOutlined/>,
    },
    {
        title: "系统设置",
        desc: "站点、界面、支付配置",
        route: AdminPageCode.SYSTEM_CONFIG,
        icon: <SettingOutlined/>,
    },
];

export default (): React.ReactNode => {
    const {styles} = useStyles();

    return (
        <Card className={styles.card} variant="borderless">
            <div className={styles.grid}>
                {quickItems.map((item) => (
                    <div
                        className={styles.item}
                        key={item.title}
                        onClick={() => historyPush(item.route)}
                    >
                        <div className="quick-icon">
                            {item.icon}
                        </div>
                        <div className="quick-main">
                            <span className="quick-title">{item.title}</span>
                            <span className="quick-desc">{item.desc}</span>
                        </div>
                        <span className="quick-arrow">
                            <RightOutlined/>
                        </span>
                    </div>
                ))}
            </div>
        </Card>
    );
};
