import React, {useEffect, useMemo, useState} from "react";
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
import {Title} from "sinking-antd";
import {getServiceCount} from "@/service/admin/dashboard";

type ServiceItem = {
    name: string;
    code: string;
    desc: string;
    count: number;
    scope: string;
    tone: string;
    icon: React.ReactNode;
};

const useStyles: any = createStyles(({css, token, isDarkMode}: any): any => ({
    card: {
        borderRadius: token.borderRadiusLG,
        overflow: "hidden",
        ".ant-card-head": {
            marginBottom: 0,
            borderBottom: "none",
            boxShadow: `inset 0 -1px 0 ${isDarkMode ? "rgba(255,255,255,0.06)" : "rgba(5,5,5,0.06)"}`,
        },
        ".ant-card-body": {
            padding: 0,
            overflow: "hidden",
            borderRadius: `0 0 ${token.borderRadiusLG}px ${token.borderRadiusLG}px`,
        },
    },
    extra: {
        display: "flex",
        alignItems: "center",
        gap: "10px",
        flexShrink: 0,
        ".health-dot": {
            position: "relative",
            width: "8px",
            height: "8px",
            borderRadius: "50%",
            background: "var(--overview-status-color)",
            boxShadow: `0 0 0 4px var(--overview-status-shadow)`,
        },
    },
    summary: {
        display: "grid",
        gridTemplateColumns: "300px minmax(0, 1fr)",
        minHeight: "210px",
        overflow: "hidden",
        borderRadius: `0 0 ${token.borderRadiusLG}px ${token.borderRadiusLG}px`,
        "@media (max-width: 960px)": {
            gridTemplateColumns: "1fr",
        },
    },
    profile: {
        position: "relative",
        display: "flex",
        flexDirection: "column",
        justifyContent: "flex-start",
        gap: "16px",
        padding: "20px",
        borderRight: `1px solid ${isDarkMode ? "rgba(255,255,255,0.07)" : token.colorBorderSecondary}`,
        borderBottomLeftRadius: token.borderRadiusLG,
        background: isDarkMode ? "rgba(255,255,255,0.012)" : token.colorFillAlter,
        overflow: "hidden",
        isolation: "isolate",
        transform: "perspective(920px) rotateX(var(--profile-rotate-x, 0deg)) rotateY(var(--profile-rotate-y, 0deg)) scale(var(--profile-scale, 1))",
        transformOrigin: "50% 50%",
        transformStyle: "preserve-3d",
        transition: "transform var(--profile-transform-duration, 0.18s) ease-out, box-shadow 0.32s ease, background-color 0.32s ease",
        willChange: "transform",
        "&:before": {
            position: "absolute",
            top: "-58px",
            right: "-72px",
            width: "172px",
            height: "172px",
            borderRadius: "50%",
            background: `radial-gradient(circle, ${token.colorPrimary} 0%, ${
                isDarkMode ? "rgba(255,255,255,0.08)" : "rgba(255,255,255,0.7)"
            } 48%, transparent 78%)`,
            content: "\"\"",
            filter: "blur(36px)",
            opacity: isDarkMode ? 0.09 : 0.12,
            transform: "rotate(-18deg)",
        },
        "&:after": {
            position: "absolute",
            left: "-68px",
            bottom: "-72px",
            width: "178px",
            height: "178px",
            borderRadius: "50%",
            background: `radial-gradient(circle, #13c2c2 0%, ${
                isDarkMode ? "rgba(255,255,255,0.07)" : "rgba(255,255,255,0.72)"
            } 50%, transparent 78%)`,
            content: "\"\"",
            filter: "blur(38px)",
            opacity: isDarkMode ? 0.08 : 0.1,
            transform: "rotate(18deg)",
        },
        ".profile-head, .stack-list": {
            position: "relative",
            zIndex: 2,
            transform: "translateZ(22px)",
        },
        ".profile-ambient": {
            position: "absolute",
            inset: 0,
            zIndex: 0,
            pointerEvents: "none",
        },
        ".profile-shine": {
            position: "absolute",
            inset: 0,
            zIndex: 1,
            opacity: 0,
            pointerEvents: "none",
            background: `radial-gradient(circle at var(--profile-shine-x, 50%) var(--profile-shine-y, 50%), ${
                isDarkMode ? "rgba(255,255,255,0.16)" : "rgba(255,255,255,0.78)"
            } 0%, ${isDarkMode ? "rgba(255,255,255,0.06)" : "rgba(22,119,255,0.10)"} 18%, transparent 48%)`,
            transition: "opacity 0.24s ease",
            mixBlendMode: isDarkMode ? "screen" : "normal",
        },
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
            background: isDarkMode ? "rgba(20,20,20,0.58)" : "rgba(255,255,255,0.78)",
            backdropFilter: "blur(10px)",
            transform: "translateZ(0)",
            transition: "transform 0.26s ease, border-color 0.26s ease, box-shadow 0.26s ease, background-color 0.26s ease",
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
            borderBottomLeftRadius: 0,
            borderBottom: `1px solid ${isDarkMode ? "rgba(255,255,255,0.06)" : "rgba(5,5,5,0.06)"}`,
        },
        "@media (max-width: 576px)": {
            padding: "16px",
        },
        "@media (hover: hover)": {
            "&:hover": {
                "--profile-scale": "1.012",
                boxShadow: isDarkMode ? "inset 0 1px 0 rgba(255,255,255,0.04)" : "0 14px 34px rgba(22,119,255,0.07)",
            },
            "&:hover .profile-shine": {
                opacity: isDarkMode ? 0.32 : 0.44,
            },
            ".stack-item:hover": {
                transform: "perspective(640px) translateZ(14px) translateY(-1px)",
                borderColor: isDarkMode ? "rgba(255,255,255,0.12)" : "rgba(22,119,255,0.18)",
                boxShadow: isDarkMode ? "0 10px 22px rgba(0,0,0,0.22)" : "0 10px 22px rgba(22,119,255,0.08)",
                background: isDarkMode ? "rgba(24,24,24,0.66)" : "rgba(255,255,255,0.86)",
            },
        },
    },
    profileMotion: css`
        .profile-ambient::before {
            content: '';
            position: absolute;
            inset: -52px;
            background:
                radial-gradient(circle at 24% 18%, ${isDarkMode ? "rgba(255,255,255,0.10)" : "rgba(22,119,255,0.16)"} 0%, transparent 36%),
                radial-gradient(circle at 80% 72%, ${isDarkMode ? "rgba(255,255,255,0.08)" : "rgba(19,194,194,0.13)"} 0%, transparent 38%),
                linear-gradient(112deg, transparent 0 40%, ${isDarkMode ? "rgba(255,255,255,0.11)" : "rgba(22,119,255,0.16)"} 48%, transparent 58%) -220px 0 / 210px 100% no-repeat,
                linear-gradient(90deg, ${isDarkMode ? "rgba(255,255,255,0.024)" : "rgba(22,119,255,0.036)"} 1px, transparent 1px),
                linear-gradient(0deg, ${isDarkMode ? "rgba(255,255,255,0.02)" : "rgba(19,194,194,0.03)"} 1px, transparent 1px);
            background-size: 100% 100%, 100% 100%, 210px 100%, 32px 32px, 32px 32px;
            opacity: ${isDarkMode ? 0.32 : 0.38};
            filter: blur(2.5px);
            transform: rotate(-7deg) translate3d(-18px, -12px, 0);
            animation: overviewProfileScan 6.4s cubic-bezier(0.45, 0, 0.25, 1) infinite;
            -webkit-mask-image: linear-gradient(135deg, rgba(0, 0, 0, 0.82), rgba(0, 0, 0, 0.34) 58%, transparent 100%);
            mask-image: linear-gradient(135deg, rgba(0, 0, 0, 0.82), rgba(0, 0, 0, 0.34) 58%, transparent 100%);
        }

        .profile-ambient::after {
            content: '';
            position: absolute;
            right: -126px;
            bottom: -104px;
            width: 360px;
            height: 360px;
            border-radius: 50%;
            background:
                radial-gradient(circle, ${isDarkMode ? "rgba(255,255,255,0.08)" : "rgba(22,119,255,0.13)"} 0%, transparent 46%),
                conic-gradient(from 0deg, transparent 0 12%, ${isDarkMode ? "rgba(255,255,255,0.10)" : "rgba(22,119,255,0.16)"} 20%, transparent 34%, ${isDarkMode ? "rgba(255,255,255,0.08)" : "rgba(19,194,194,0.13)"} 52%, transparent 72%, ${isDarkMode ? "rgba(255,255,255,0.08)" : "rgba(47,84,235,0.12)"} 86%, transparent 100%);
            filter: blur(36px);
            opacity: ${isDarkMode ? 0.13 : 0.16};
            transform: rotate(0deg) scale(1);
            animation: overviewProfileOrbit 8s linear infinite;
        }

        &:hover .profile-ambient::before {
            opacity: ${isDarkMode ? 0.42 : 0.48};
            filter: blur(3.5px);
        }

        &:hover .profile-ambient::after {
            opacity: ${isDarkMode ? 0.2 : 0.24};
            filter: blur(42px);
        }

        @keyframes overviewProfileScan {
            0% {
                background-position: 0 0, 0 0, -220px 0, 0 0, 0 0;
            }
            50% {
                background-position: 0 0, 0 0, 54% 0, 16px 16px, -16px -16px;
            }
            100% {
                background-position: 0 0, 0 0, calc(100% + 220px) 0, 32px 32px, -32px -32px;
            }
        }

        @keyframes overviewProfileOrbit {
            0% {
                transform: rotate(0deg) scale(0.94);
            }
            50% {
                transform: rotate(180deg) scale(1.04);
            }
            100% {
                transform: rotate(360deg) scale(0.94);
            }
        }
    `,
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
        "&:last-child": {
            borderBottomRightRadius: token.borderRadiusLG,
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
            color: "var(--service-status-color)",
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
            background: "var(--service-status-color)",
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
        "@media (max-width: 960px)": {
            "&:nth-last-child(2)": {
                borderBottomLeftRadius: token.borderRadiusLG,
            },
        },
        "@media (max-width: 576px)": {
            minHeight: "112px",
            borderRight: "none",
            "&:nth-last-child(-n + 2)": {
                borderBottom: `1px solid ${isDarkMode ? "rgba(255,255,255,0.07)" : token.colorBorderSecondary}`,
            },
            "&:nth-last-child(2)": {
                borderBottomLeftRadius: 0,
            },
            "&:last-child": {
                borderBottom: "none",
                borderBottomLeftRadius: token.borderRadiusLG,
                borderBottomRightRadius: token.borderRadiusLG,
            },
        },
    },
}));

const services: ServiceItem[] = [
    {
        name: "网关服务",
        code: "gateway",
        desc: "统一入口、鉴权转发与路由治理",
        count: 0,
        scope: "入口",
        tone: "#1677ff",
        icon: <GatewayOutlined/>,
    },
    {
        name: "授权服务",
        code: "auth",
        desc: "验证码、登录授权与身份校验",
        count: 0,
        scope: "鉴权",
        tone: "#722ed1",
        icon: <LockOutlined/>,
    },
    {
        name: "推送服务",
        code: "socket",
        desc: "站内连接、实时消息与事件推送",
        count: 0,
        scope: "推送",
        tone: "#13c2c2",
        icon: <SendOutlined/>,
    },
    {
        name: "定时任务",
        code: "cron",
        desc: "任务编排、定时执行与运行日志",
        count: 0,
        scope: "调度",
        tone: "#faad14",
        icon: <ScheduleOutlined/>,
    },
    {
        name: "支付服务",
        code: "pay",
        desc: "订单、余额、支付通道与回调",
        count: 0,
        scope: "交易",
        tone: "#fa541c",
        icon: <CreditCardOutlined/>,
    },
    {
        name: "存储服务",
        code: "driver",
        desc: "文件场景、存储通道与资源预览",
        count: 0,
        scope: "存储",
        tone: "#2f54eb",
        icon: <DatabaseOutlined/>,
    },
    {
        name: "用户服务",
        code: "user",
        desc: "前台用户、实名信息与会话能力",
        count: 0,
        scope: "用户",
        tone: "#52c41a",
        icon: <TeamOutlined/>,
    },
    {
        name: "后台服务",
        code: "admin",
        desc: "后台员工、权限菜单与系统配置",
        count: 0,
        scope: "后台",
        tone: "#eb2f96",
        icon: <CloudServerOutlined/>,
    },
];

export default (): React.ReactNode => {
    const {styles} = useStyles();
    const [serviceCount, setServiceCount] = useState<Record<string, number>>({});
    const serviceList = useMemo(() => services.map((item) => ({
        ...item,
        count: serviceCount?.[item.code] ?? item.count,
    })), [serviceCount]);
    const overviewStatus = useMemo(() => {
        const normalCount = serviceList.filter((item) => item.count > 0).length;
        if (normalCount === serviceList.length) {
            return {
                color: "success",
                text: "服务正常",
                dotColor: "#52c41a",
                dotShadow: "rgba(82,196,26,0.14)",
            };
        }
        if (normalCount === 0) {
            return {
                color: "error",
                text: "服务异常",
                dotColor: "#ff4d4f",
                dotShadow: "rgba(255,77,79,0.14)",
            };
        }
        return {
            color: "warning",
            text: "部分异常",
            dotColor: "#faad14",
            dotShadow: "rgba(250,173,20,0.16)",
        };
    }, [serviceList]);

    useEffect(() => {
        void getServiceCount({
            onSuccess: (res) => {
                setServiceCount(res?.data || {});
            },
        });
    }, []);

    const handleProfileMouseMove = (event: React.MouseEvent<HTMLDivElement>) => {
        const target = event.currentTarget;
        const rect = target.getBoundingClientRect();
        const x = (event.clientX - rect.left) / rect.width;
        const y = (event.clientY - rect.top) / rect.height;
        target.style.setProperty("--profile-rotate-x", `${((0.5 - y) * 12).toFixed(2)}deg`);
        target.style.setProperty("--profile-rotate-y", `${((x - 0.5) * 15).toFixed(2)}deg`);
        target.style.setProperty("--profile-shine-x", `${(x * 100).toFixed(1)}%`);
        target.style.setProperty("--profile-shine-y", `${(y * 100).toFixed(1)}%`);
        target.style.setProperty("--profile-transform-duration", "0.08s");
    };

    const handleProfileMouseLeave = (event: React.MouseEvent<HTMLDivElement>) => {
        const target = event.currentTarget;
        target.style.setProperty("--profile-transform-duration", "0.55s");
        target.style.setProperty("--profile-shine-x", "50%");
        target.style.setProperty("--profile-shine-y", "50%");
        requestAnimationFrame(() => {
            target.style.setProperty("--profile-rotate-x", "0deg");
            target.style.setProperty("--profile-rotate-y", "0deg");
        });
    };

    return (
        <Card
            className={styles.card}
            variant="borderless"
            title={<Title>服务概览</Title>}
            extra={
                <div
                    className={styles.extra}
                    style={{
                        "--overview-status-color": overviewStatus.dotColor,
                        "--overview-status-shadow": overviewStatus.dotShadow,
                    } as React.CSSProperties}
                >
                    <span className="health-dot"/>
                    <Tag color={overviewStatus.color}>{overviewStatus.text}</Tag>
                </div>
            }
        >
            <div className={styles.summary}>
                <div
                    className={`${styles.profile} ${styles.profileMotion}`}
                    onMouseMove={handleProfileMouseMove}
                    onMouseLeave={handleProfileMouseLeave}
                >
                    <span className="profile-ambient"/>
                    <span className="profile-shine"/>
                    <div className="profile-head">
                        <span className="profile-title">Sinking Cloud</span>
                        <span className="profile-desc">基于sinking技术栈的微服务中后台系统</span>
                    </div>
                    <div className="stack-list">
                        <div className="stack-item">
                            <span>
                                <span className="stack-name">sinking-antd</span>
                                <span className="stack-desc">界面组件体系</span>
                            </span>
                            <Tag className="stack-tag">Ui</Tag>
                        </div>
                        <div className="stack-item">
                            <span>
                                <span className="stack-name">sinking-go</span>
                                <span className="stack-desc">后端服务框架</span>
                            </span>
                            <Tag color="processing" className="stack-tag">Go</Tag>
                        </div>
                        <div className="stack-item">
                            <span>
                                <span className="stack-name">sinking-consule</span>
                                <span className="stack-desc">服务注册发现与配置管理</span>
                            </span>
                            <Tag color="success" className="stack-tag">Consul</Tag>
                        </div>
                    </div>
                </div>
                <div className={styles.servicePanel}>
                    <div className={styles.serviceList}>
                        {serviceList.map((item) => (
                            <div
                                className={styles.service}
                                key={item.code}
                                style={{
                                    "--service-color": item.tone,
                                    "--service-status-color": item.count > 0 ? "#52c41a" : "#ff4d4f",
                                } as React.CSSProperties}
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
                                            {item.count > 0 ? "健康" : "异常"}
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
