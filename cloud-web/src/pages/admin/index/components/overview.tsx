import React from "react";
import {Card, Tag} from "antd";
import {
    CloudServerOutlined,
    CreditCardOutlined,
    DatabaseOutlined,
    GatewayOutlined,
    LockOutlined,
    ScheduleOutlined,
    SendOutlined,
    TeamOutlined,
} from "@ant-design/icons";
import {createStyles} from "antd-style";

type ServiceItem = {
    name: string;
    code: string;
    desc: string;
    count: number;
    status: string;
    scope: string;
    tone: string;
    icon: React.ReactNode;
};

const useStyles: any = createStyles(({token, isDarkMode}: any): any => ({
    card: {
        borderRadius: token.borderRadiusLG,
        ".ant-card-body": {
            padding: 0,
        },
    },
    header: {
        display: "flex",
        alignItems: "center",
        justifyContent: "space-between",
        gap: "16px",
        padding: "18px 20px",
        borderBottom: `1px solid ${isDarkMode ? "rgba(255,255,255,0.07)" : token.colorBorderSecondary}`,
        ".overview-title": {
            display: "block",
            color: token.colorTextHeading,
            fontSize: "15px",
            fontWeight: 600,
            lineHeight: "22px",
        },
        ".overview-desc": {
            display: "block",
            marginTop: "3px",
            color: token.colorTextTertiary,
            fontSize: "12px",
            lineHeight: "18px",
        },
        ".overview-health": {
            display: "flex",
            alignItems: "center",
            gap: "10px",
            flexShrink: 0,
        },
        ".health-dot": {
            position: "relative",
            width: "8px",
            height: "8px",
            borderRadius: "50%",
            background: token.colorSuccess,
            boxShadow: `0 0 0 4px ${isDarkMode ? "rgba(82,196,26,0.12)" : "rgba(82,196,26,0.14)"}`,
        },
        "@media (max-width: 576px)": {
            alignItems: "flex-start",
            flexDirection: "column",
            padding: "16px",
        },
    },
    summary: {
        display: "grid",
        gridTemplateColumns: "300px minmax(0, 1fr)",
        minHeight: "210px",
        "@media (max-width: 960px)": {
            gridTemplateColumns: "1fr",
        },
    },
    profile: {
        display: "flex",
        flexDirection: "column",
        justifyContent: "flex-start",
        gap: "16px",
        padding: "20px",
        borderRight: `1px solid ${isDarkMode ? "rgba(255,255,255,0.07)" : token.colorBorderSecondary}`,
        background: isDarkMode ? "rgba(255,255,255,0.012)" : token.colorFillAlter,
        ".profile-title": {
            display: "block",
            color: token.colorTextHeading,
            fontWeight: 600,
            fontSize: "14px",
            lineHeight: "22px",
        },
        ".profile-desc": {
            display: "block",
            marginTop: "5px",
            color: token.colorTextTertiary,
            fontSize: "12px",
            lineHeight: "18px",
        },
        ".stack-list": {
            display: "flex",
            flexDirection: "column",
            gap: "10px",
        },
        ".stack-item": {
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            gap: "12px",
            padding: "11px 12px",
            border: `1px solid ${isDarkMode ? "rgba(255,255,255,0.07)" : token.colorBorderSecondary}`,
            borderRadius: token.borderRadius,
            background: token.colorBgContainer,
        },
        ".stack-name": {
            display: "block",
            color: token.colorTextHeading,
            fontSize: "13px",
            fontWeight: 600,
            lineHeight: "20px",
        },
        ".stack-desc": {
            display: "block",
            color: token.colorTextTertiary,
            fontSize: "12px",
            lineHeight: "18px",
        },
        ".stack-tag": {
            flexShrink: 0,
        },
        ".ant-tag": {
            marginInlineEnd: 0,
        },
        "@media (max-width: 960px)": {
            borderRight: "none",
            borderBottom: `1px solid ${isDarkMode ? "rgba(255,255,255,0.07)" : token.colorBorderSecondary}`,
        },
        "@media (max-width: 576px)": {
            padding: "16px",
        },
    },
    serviceList: {
        display: "grid",
        alignContent: "stretch",
        height: "100%",
        gridTemplateColumns: "repeat(4, minmax(0, 1fr))",
        gridAutoRows: "1fr",
        "@media (max-width: 1180px)": {
            gridTemplateColumns: "repeat(2, minmax(0, 1fr))",
        },
        "@media (max-width: 576px)": {
            gridTemplateColumns: "1fr",
        },
    },
    servicePanel: {
        minWidth: 0,
        height: "100%",
    },
    service: {
        position: "relative",
        display: "flex",
        flexDirection: "column",
        justifyContent: "space-between",
        gap: "12px",
        minHeight: "116px",
        padding: "15px 16px 13px",
        borderRight: `1px solid ${isDarkMode ? "rgba(255,255,255,0.07)" : token.colorBorderSecondary}`,
        borderBottom: `1px solid ${isDarkMode ? "rgba(255,255,255,0.07)" : token.colorBorderSecondary}`,
        background: token.colorBgContainer,
        overflow: "hidden",
        transition: "background-color 0.2s ease",
        "&:hover": {
            background: isDarkMode ? "rgba(255,255,255,0.03)" : token.colorFillQuaternary,
        },
        "&:before": {
            position: "absolute",
            right: "-66px",
            bottom: "-70px",
            width: "184px",
            height: "184px",
            borderRadius: "50%",
            background: `radial-gradient(circle, var(--service-color) 0%, ${
                isDarkMode ? "rgba(255,255,255,0.08)" : "rgba(255,255,255,0.68)"
            } 48%, transparent 78%)`,
            content: "\"\"",
            filter: "blur(24px)",
            opacity: isDarkMode ? 0.22 : 0.24,
            transform: "rotate(-18deg)",
        },
        "&:nth-child(4n)": {
            borderRight: "none",
        },
        "&:nth-last-child(-n + 4)": {
            borderBottom: "none",
        },
        ".service-watermark": {
            position: "absolute",
            right: "-25px",
            bottom: "-29px",
            color: "var(--service-color)",
            fontSize: "138px",
            lineHeight: 1,
            opacity: isDarkMode ? 0.09 : 0.075,
            pointerEvents: "none",
            transform: "rotate(-24deg)",
            filter: "blur(7.5px)",
        },
        ".service-main": {
            position: "relative",
            zIndex: 1,
            minWidth: 0,
        },
        ".service-icon": {
            width: "28px",
            height: "28px",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            flexShrink: 0,
            color: "var(--service-color)",
            fontSize: "20px",
        },
        ".service-top": {
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            gap: "10px",
            minWidth: 0,
        },
        ".service-title": {
            display: "flex",
            alignItems: "center",
            gap: "9px",
            minWidth: 0,
        },
        ".service-title-content": {
            display: "block",
            minWidth: 0,
        },
        ".service-name": {
            display: "block",
            color: token.colorTextHeading,
            fontSize: "13px",
            fontWeight: 600,
            lineHeight: "20px",
            overflow: "hidden",
            textOverflow: "ellipsis",
            whiteSpace: "nowrap",
        },
        ".service-status": {
            display: "inline-flex",
            alignItems: "center",
            gap: "5px",
            color: token.colorSuccess,
            fontSize: "12px",
            lineHeight: "18px",
            flexShrink: 0,
        },
        ".service-code": {
            display: "block",
            marginTop: "1px",
            color: token.colorTextQuaternary,
            fontSize: "12px",
            lineHeight: "18px",
            overflow: "hidden",
            textOverflow: "ellipsis",
            whiteSpace: "nowrap",
        },
        ".service-desc": {
            display: "-webkit-box",
            marginTop: "10px",
            color: token.colorTextTertiary,
            fontSize: "12px",
            lineHeight: "18px",
            overflow: "hidden",
            WebkitBoxOrient: "vertical",
            WebkitLineClamp: 2,
        },
        ".service-footer": {
            position: "relative",
            zIndex: 1,
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            gap: "10px",
            minWidth: 0,
            paddingTop: "9px",
            borderTop: `1px solid ${isDarkMode ? "rgba(255,255,255,0.06)" : token.colorBorderSecondary}`,
        },
        ".service-meta": {
            display: "inline-flex",
            alignItems: "center",
            gap: "8px",
            flexShrink: 0,
        },
        ".service-count": {
            color: token.colorTextSecondary,
            fontSize: "12px",
            lineHeight: "18px",
            whiteSpace: "nowrap",
        },
        ".service-scope": {
            color: token.colorTextQuaternary,
            fontSize: "12px",
            lineHeight: "18px",
            whiteSpace: "nowrap",
        },
        ".status-dot": {
            width: "6px",
            height: "6px",
            borderRadius: "50%",
            background: token.colorSuccess,
        },
        "@media (max-width: 1180px)": {
            "&:nth-child(2n)": {
                borderRight: "none",
            },
            "&:nth-last-child(-n + 4)": {
                borderBottom: `1px solid ${isDarkMode ? "rgba(255,255,255,0.07)" : token.colorBorderSecondary}`,
            },
            "&:nth-last-child(-n + 2)": {
                borderBottom: "none",
            },
        },
        "@media (max-width: 576px)": {
            minHeight: "112px",
            borderRight: "none",
            "&:nth-last-child(-n + 2)": {
                borderBottom: `1px solid ${isDarkMode ? "rgba(255,255,255,0.07)" : token.colorBorderSecondary}`,
            },
            "&:last-child": {
                borderBottom: "none",
            },
        },
    },
}));

const services: ServiceItem[] = [
    {
        name: "网关服务",
        code: "gateway",
        desc: "统一入口、鉴权转发与路由治理",
        count: 1,
        status: "运行中",
        scope: "入口",
        tone: "#1677ff",
        icon: <GatewayOutlined/>,
    },
    {
        name: "授权服务",
        code: "auth",
        desc: "验证码、登录授权与身份校验",
        count: 1,
        status: "运行中",
        scope: "鉴权",
        tone: "#722ed1",
        icon: <LockOutlined/>,
    },
    {
        name: "推送服务",
        code: "socket",
        desc: "站内连接、实时消息与事件推送",
        count: 1,
        status: "运行中",
        scope: "推送",
        tone: "#13c2c2",
        icon: <SendOutlined/>,
    },
    {
        name: "定时任务",
        code: "cron",
        desc: "任务编排、定时执行与运行日志",
        count: 1,
        status: "运行中",
        scope: "调度",
        tone: "#faad14",
        icon: <ScheduleOutlined/>,
    },
    {
        name: "支付服务",
        code: "pay",
        desc: "订单、余额、支付通道与回调",
        count: 1,
        status: "运行中",
        scope: "交易",
        tone: "#fa541c",
        icon: <CreditCardOutlined/>,
    },
    {
        name: "存储服务",
        code: "driver",
        desc: "文件场景、存储通道与资源预览",
        count: 1,
        status: "运行中",
        scope: "存储",
        tone: "#2f54eb",
        icon: <DatabaseOutlined/>,
    },
    {
        name: "用户服务",
        code: "user",
        desc: "前台用户、实名信息与会话能力",
        count: 1,
        status: "运行中",
        scope: "用户",
        tone: "#52c41a",
        icon: <TeamOutlined/>,
    },
    {
        name: "后台服务",
        code: "admin",
        desc: "后台员工、权限菜单与系统配置",
        count: 1,
        status: "运行中",
        scope: "后台",
        tone: "#eb2f96",
        icon: <CloudServerOutlined/>,
    },
];

export default (): React.ReactNode => {
    const {styles} = useStyles();

    return (
        <Card className={styles.card} variant="borderless">
            <div className={styles.header}>
                <div>
                    <span className="overview-title">服务概览</span>
                    <span className="overview-desc">面向个人业务沉淀的微服务后台基础系统</span>
                </div>
                <div className="overview-health">
                    <span className="health-dot"/>
                    <Tag color="success">服务正常</Tag>
                </div>
            </div>
            <div className={styles.summary}>
                <div className={styles.profile}>
                    <div>
                        <span className="profile-title">Sinking Cloud</span>
                        <span className="profile-desc">基于sinking技术栈的微服务中后台系统</span>
                    </div>
                    <div className="stack-list">
                        <div className="stack-item">
                            <span>
                                <span className="stack-name">sinking-antd</span>
                                <span className="stack-desc">后台界面组件体系</span>
                            </span>
                            <Tag className="stack-tag">UI</Tag>
                        </div>
                        <div className="stack-item">
                            <span>
                                <span className="stack-name">sinking-go</span>
                                <span className="stack-desc">后端业务服务框架</span>
                            </span>
                            <Tag color="processing" className="stack-tag">Go</Tag>
                        </div>
                        <div className="stack-item">
                            <span>
                                <span className="stack-name">sinking-consule</span>
                                <span className="stack-desc">服务注册发现与节点治理</span>
                            </span>
                            <Tag color="success" className="stack-tag">Consul</Tag>
                        </div>
                    </div>
                </div>
                <div className={styles.servicePanel}>
                    <div className={styles.serviceList}>
                        {services.map((item) => (
                            <div
                                className={styles.service}
                                key={item.code}
                                style={{"--service-color": item.tone} as React.CSSProperties}
                            >
                                <span className="service-watermark">{item.icon}</span>
                                <div className="service-main">
                                    <div className="service-top">
                                        <span className="service-title">
                                            <span className="service-icon">
                                                {item.icon}
                                            </span>
                                            <span className="service-title-content">
                                                <span className="service-name">{item.name}</span>
                                                <span className="service-code">{item.code}</span>
                                            </span>
                                        </span>
                                        <span className="service-status">
                                            <span className="status-dot"/>
                                            {item.status}
                                        </span>
                                    </div>
                                    <span className="service-desc">{item.desc}</span>
                                </div>
                                <div className="service-footer">
                                    <span className="service-code">{item.code}</span>
                                    <span className="service-meta">
                                        <span className="service-count">{item.count} 实例</span>
                                        <span className="service-scope">{item.scope}</span>
                                    </span>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </Card>
    );
};
