import React, {useEffect, useState} from "react";
import {useModel} from "umi";
import {App, Button, Card, Empty, Pagination, Spin, Tag, Tooltip, Typography} from "antd";
import {createStyles} from "antd-style";
import {Body, Icon} from "sinking-antd";
import {getSessionList, logoutSession} from "@/service/user/person";
import {ago} from "@/utils/time";

const DEFAULT_PAGE_SIZE = 8;

const useStyles = createStyles(({token, css, isDarkMode}: any) => {
    const compact = token.sizeStep <= 3;
    const radius = token.borderRadiusLG;
    const paginationColor = isDarkMode ? "rgba(255,255,255,0.65)" : "rgba(0,0,0,0.65)";
    const paginationBg = isDarkMode ? "rgba(255, 255, 255, 0.1) !important" : "rgba(0, 0, 0, 0.05) !important";

    return {
        page: css`
            display: flex;
            flex-direction: column;
            gap: ${compact ? 12 : 14}px;
            max-width: 1120px;
            margin: 0 auto;
        `,
        hero: css`
            position: relative;
            overflow: hidden;
            isolation: isolate;
            min-height: ${compact ? 120 : 130}px;
            border-radius: ${radius}px;
            padding: ${compact ? "16px 20px" : "18px 24px"};
            background:
                radial-gradient(circle at 18% 16%, ${isDarkMode ? "rgba(255,255,255,0.11)" : `${token.colorPrimary}26`} 0%, transparent 34%),
                radial-gradient(circle at 82% 42%, ${isDarkMode ? "rgba(255,255,255,0.08)" : "rgba(19,194,194,0.20)"} 0%, transparent 36%),
                radial-gradient(circle at 64% 100%, ${isDarkMode ? "rgba(255,255,255,0.06)" : "rgba(114,46,209,0.13)"} 0%, transparent 32%),
                linear-gradient(135deg, ${token.colorBgContainer} 0%, ${isDarkMode ? "rgba(255,255,255,0.018)" : token.colorFillQuaternary} 100%);
            box-shadow: ${isDarkMode ? "0 14px 38px rgba(0,0,0,0.24)" : "0 16px 36px rgba(15,23,42,0.065)"};

            &::before {
                content: "";
                position: absolute;
                right: -130px;
                bottom: -160px;
                width: 460px;
                height: 360px;
                border-radius: 50%;
                background:
                    radial-gradient(circle, ${token.colorPrimaryBgHover} 0%, transparent 46%),
                    conic-gradient(from 20deg, transparent 0 14%, ${token.colorPrimaryBorder} 22%, transparent 38%, ${token.colorPrimaryBgHover} 56%, transparent 76%, ${token.colorPrimaryBg} 90%, transparent 100%);
                filter: blur(52px);
                opacity: ${isDarkMode ? 0.16 : 0.36};
                pointer-events: none;
                z-index: 0;
            }

            &::after {
                content: "";
                position: absolute;
                inset: 0;
                z-index: 0;
                opacity: ${isDarkMode ? 0.16 : 0.26};
                pointer-events: none;
                background:
                    linear-gradient(90deg, ${isDarkMode ? "rgba(255,255,255,0.034)" : "rgba(22,119,255,0.040)"} 1px, transparent 1px),
                    linear-gradient(0deg, ${isDarkMode ? "rgba(255,255,255,0.026)" : "rgba(19,194,194,0.032)"} 1px, transparent 1px);
                background-size: 34px 34px;
                mask-image: linear-gradient(110deg, transparent 0%, #000 18%, #000 52%, transparent 88%);
            }

            .hero-watermark {
                position: absolute;
                right: -30px;
                bottom: -92px;
                z-index: 0;
                color: ${token.colorPrimary};
                font-size: ${compact ? 210 : 250}px;
                line-height: 1;
                opacity: ${isDarkMode ? 0.11 : 0.085};
                filter: blur(14px);
                transform: rotate(-18deg);
                pointer-events: none;
            }

            @media (max-width: 900px) {
                min-height: auto;
                padding: ${compact ? "18px" : "20px"};
            }

            @media (max-width: 620px) {
                padding: ${compact ? "15px 16px" : "16px 18px"};

                .hero-watermark {
                    right: -42px;
                    bottom: -58px;
                    font-size: ${compact ? 160 : 176}px;
                    opacity: ${isDarkMode ? 0.08 : 0.06};
                    filter: blur(12px);
                }
            }
        `,
        heroInner: css`
            position: relative;
            z-index: 1;
            display: grid;
            grid-template-columns: minmax(0, 1fr) minmax(230px, 260px);
            align-items: center;
            gap: ${compact ? 12 : 16}px;

            @media (max-width: 900px) {
                grid-template-columns: 1fr;
                align-items: stretch;
                gap: ${compact ? 14 : 16}px;
            }

            @media (max-width: 620px) {
                gap: ${compact ? 10 : 12}px;
            }
        `,
        heroContent: css`
            min-width: 0;

            @media (max-width: 900px) {
                width: 100%;
            }
        `,
        heroTitle: css`
            margin: 0;
            color: ${token.colorText};
            font-size: ${compact ? 26 : 28}px;
            font-weight: 700;
            line-height: ${compact ? 34 : 36}px;
            letter-spacing: 0;

            @media (max-width: 900px) {
                font-size: ${compact ? 23 : 24}px;
                line-height: ${compact ? 31 : 32}px;
            }

            @media (max-width: 620px) {
                font-size: ${compact ? 22 : 23}px;
                line-height: ${compact ? 30 : 31}px;
            }
        `,
        heroDesc: css`
            max-width: 600px;
            margin-top: 6px;
            color: ${token.colorTextSecondary};
            font-size: 13px;
            line-height: 22px;

            @media (max-width: 620px) {
                font-size: 12px;
                line-height: 20px;
            }
        `,
        heroVisual: css`
            position: relative;
            display: flex;
            align-items: center;
            justify-content: flex-end;
            min-width: 0;
            min-height: auto;

            @media (max-width: 900px) {
                justify-content: stretch;
                width: 100%;
            }

        `,
        visualCard: css`
            position: relative;
            box-sizing: border-box;
            overflow: hidden;
            width: 100%;
            padding: ${compact ? 10 : 11}px ${compact ? 12 : 14}px;
            border-radius: ${radius}px;
            background:
                linear-gradient(135deg, ${isDarkMode ? "rgba(255,255,255,0.052)" : "rgba(255,255,255,0.76)"} 0%, ${isDarkMode ? "rgba(255,255,255,0.026)" : "rgba(255,255,255,0.54)"} 100%);
            box-shadow:
                inset 0 0 0 1px ${isDarkMode ? "rgba(255,255,255,0.070)" : "rgba(22,119,255,0.105)"},
                ${isDarkMode ? "0 10px 22px rgba(0,0,0,0.16)" : "0 12px 24px rgba(22,119,255,0.08)"};
            backdrop-filter: blur(14px);

            &::before {
                content: "";
                position: absolute;
                right: -36px;
                top: -44px;
                width: 112px;
                height: 112px;
                border-radius: 999px;
                background: radial-gradient(circle, ${token.colorPrimaryBgHover} 0%, transparent 66%);
                opacity: ${isDarkMode ? 0.18 : 0.55};
                pointer-events: none;
            }

            .visual-head {
                position: relative;
                z-index: 1;
                display: flex;
                align-items: center;
                justify-content: space-between;
                gap: 8px;
                margin-bottom: ${compact ? 7 : 8}px;
            }

            .visual-title {
                color: ${token.colorText};
                font-size: 12px;
                font-weight: 600;
                line-height: 18px;
            }

            .visual-status {
                display: inline-flex;
                align-items: center;
                gap: 4px;
                padding: 2px 7px;
                border-radius: 999px;
                background: ${token.colorSuccessBg};
                color: ${token.colorSuccess};
                font-size: 11px;
                line-height: 16px;
            }

            .visual-status::before {
                content: "";
                width: 5px;
                height: 5px;
                border-radius: 999px;
                background: ${token.colorSuccess};
            }

            .visual-device {
                position: relative;
                z-index: 1;
                display: block;
                padding-bottom: ${compact ? 7 : 8}px;
                color: ${token.colorText};
                font-size: 13px;
                font-weight: 600;
                line-height: 20px;
                box-shadow: inset 0 -1px 0 ${isDarkMode ? "rgba(255,255,255,0.060)" : "rgba(5,5,5,0.055)"};
            }

            .visual-list {
                position: relative;
                z-index: 1;
                display: grid;
                gap: ${compact ? 3 : 4}px;
                margin-top: ${compact ? 6 : 8}px;
            }

            .visual-row {
                display: grid;
                grid-template-columns: 34px minmax(0, 1fr);
                align-items: center;
                gap: 8px;
                min-width: 0;
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 18px;
            }

            .visual-row .value {
                display: block;
                overflow: hidden;
                min-width: 0;
                color: ${token.colorTextSecondary};
                font-size: 11px;
                text-align: right;
                text-overflow: ellipsis;
                white-space: nowrap;
            }

            @media (max-width: 900px) {
                max-width: none;
                padding: ${compact ? "12px 14px" : "14px 16px"};

                .visual-list {
                    grid-template-columns: repeat(3, minmax(0, 1fr));
                    gap: 8px;
                }

                .visual-row {
                    display: block;
                    min-width: 0;
                }

                .visual-row .value {
                    margin-top: 2px;
                    text-align: left;
                }
            }

            @media (max-width: 620px) {
                padding: ${compact ? "10px 11px" : "11px 12px"};

                .visual-head {
                    margin-bottom: 7px;
                }

                .visual-device {
                    padding-bottom: 7px;
                    font-size: 12px;
                    line-height: 18px;
                }

                .visual-list {
                    grid-template-columns: repeat(3, minmax(0, 1fr));
                    gap: 6px 8px;
                }

                .visual-row {
                    display: block;
                }

                .visual-row .value {
                    margin-top: 1px;
                    text-align: left;
                }
            }

            @media (max-width: 420px) {
                .visual-list {
                    grid-template-columns: 1fr;
                    gap: 3px;
                }

                .visual-row {
                    display: grid;
                    grid-template-columns: 34px minmax(0, 1fr);
                }

                .visual-row .value {
                    margin-top: 0;
                    text-align: right;
                }
            }
        `,
        listCard: css`
            overflow: hidden;
            border-radius: ${radius}px;

            .ant-card-head {
                min-height: ${compact ? 46 : 50}px;
                padding: 0 ${compact ? 14 : 16}px;
                border-bottom-color: ${isDarkMode ? "rgba(255,255,255,0.07)" : "rgba(5,5,5,0.06)"};
                background:
                    radial-gradient(circle at 0% 0%, ${isDarkMode ? "rgba(255,255,255,0.05)" : `${token.colorPrimary}0d`} 0%, transparent 34%),
                    ${token.colorBgContainer};
            }

            .ant-card-head-title {
                padding: 0;
            }

            .ant-card-extra {
                padding: 0;
            }

            .ant-card-body {
                padding: 0;
            }
        `,
        panelTitle: css`
            display: inline-flex;
            align-items: center;
            gap: 8px;
            color: ${token.colorText};
            font-size: ${compact ? 14 : 15}px;
            font-weight: 600;

            .anticon {
                color: ${token.colorPrimary};
            }
        `,
        panelExtra: css`
            display: flex;
            align-items: center;
            justify-content: flex-end;
            gap: 10px;
            color: ${token.colorTextTertiary};
            font-size: 12px;
        `,
        listBody: css`
            min-height: 180px;
            padding: ${compact ? 10 : 12}px;
            background: ${isDarkMode ? "rgba(255,255,255,0.012)" : token.colorFillQuaternary};
        `,
        sessionGrid: css`
            display: grid;
            grid-template-columns: repeat(2, minmax(0, 1fr));
            gap: ${compact ? 10 : 12}px;

            @media (max-width: 900px) {
                grid-template-columns: 1fr;
            }
        `,
        sessionCard: css`
            position: relative;
            display: flex;
            flex-direction: column;
            gap: ${compact ? 10 : 12}px;
            min-width: 0;
            padding: ${compact ? "12px" : "14px"};
            border: 1px solid ${token.colorBorderSecondary};
            border-radius: ${radius}px;
            background: ${token.colorBgContainer};
            transition: border-color 0.18s ease, background-color 0.18s ease, box-shadow 0.18s ease;

            &:hover {
                border-color: ${token.colorPrimaryBorder};
                background: ${isDarkMode ? "rgba(255,255,255,0.018)" : token.colorFillAlter};
                box-shadow: ${isDarkMode ? "none" : "0 8px 18px rgba(15,23,42,0.04)"};
            }
        `,
        deviceMain: css`
            min-width: 0;
            display: flex;
            align-items: flex-start;
            justify-content: space-between;
            gap: ${compact ? 10 : 12}px;
            padding-bottom: ${compact ? 9 : 10}px;
            border-bottom: 1px solid ${token.colorBorderSecondary};

            .name-row {
                display: flex;
                align-items: center;
                justify-content: space-between;
                min-width: 0;
                gap: 6px;
            }

            .device-content {
                min-width: 0;
                flex: 1;
            }
        `,
        deviceName: css`
            min-width: 0;
            color: ${token.colorText};
            font-size: 13px;
            font-weight: 600;
        `,
        deviceSub: css`
            margin-top: 4px;
            color: ${token.colorTextTertiary};
            font-size: 12px;
        `,
        sessionInfo: css`
            min-width: 0;
            display: grid;
            grid-template-columns: repeat(3, minmax(0, 1fr));
            gap: 1px;
            overflow: hidden;
            border-radius: ${radius}px;
            background: ${token.colorBorderSecondary};

            @media (max-width: 620px) {
                grid-template-columns: repeat(2, minmax(0, 1fr));
            }

            @media (max-width: 430px) {
                grid-template-columns: 1fr;
            }
        `,
        sessionMeta: css`
            min-width: 0;
            padding: ${compact ? "8px 9px" : "9px 10px"};
            background: ${isDarkMode ? "rgba(255,255,255,0.018)" : token.colorFillQuaternary};

            .label {
                color: ${token.colorTextTertiary};
                font-size: 11px;
                line-height: 15px;
            }

            .value {
                display: block;
                color: ${token.colorTextSecondary};
                font-size: 12px;
                line-height: 17px;
                margin-top: 2px;
            }
        `,
        sessionOperate: css`
            display: flex;
            align-items: center;
            justify-content: flex-end;
            gap: ${compact ? 8 : 10}px;
            flex: none;
            white-space: nowrap;
        `,
        pagerCard: css`
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: ${compact ? 10 : 12}px;
            border-radius: ${radius}px;
            background: ${token.colorBgContainer};

            .ant-card-body {
                width: 100%;
                display: flex;
                align-items: center;
                justify-content: center;
                padding: ${compact ? 12 : 15}px;
            }

            .ant-pagination {
                color: ${paginationColor};
                font-size: 12px;
            }

            .ant-pagination-item,
            .ant-pagination-item-active {
                background-color: transparent !important;
            }

            .ant-pagination-item a {
                color: ${paginationColor};
            }

            .ant-pagination-options {
                display: block;

                .ant-select {
                    height: var(--ant-pagination-item-size-actual);
                    line-height: calc(var(--ant-pagination-item-size-actual) - 2px);
                }

                .ant-select-selector {
                    color: ${paginationColor} !important;
                    border-color: transparent !important;
                    border-radius: ${token.borderRadius}px !important;
                    background-color: ${paginationBg};
                }

                .ant-select-selection-item {
                    color: ${paginationColor};
                    font-size: 12px;
                    line-height: calc(var(--ant-pagination-item-size-actual) - 2px);
                }

                .ant-select-arrow {
                    color: ${paginationColor};
                }

                .ant-select-content {
                    color: ${paginationColor};
                    font-size: 12px;
                    transition: background-color 0.3s ease;
                    border-color: transparent !important;
                    border-radius: ${token.borderRadius}px !important;
                }

                .ant-pagination-options-quick-jumper input {
                    color: ${paginationColor};
                    font-size: 12px;
                    border-color: transparent;
                    border-radius: ${token.borderRadius}px !important;
                    background-color: ${paginationBg};
                }
            }

            @media (max-width: 620px) {
                .ant-card-body {
                    padding: ${compact ? 10 : 12}px;
                }

                .ant-pagination {
                    width: 100%;
                    justify-content: center;
                }

                .ant-pagination-total-text,
                .ant-pagination-options-quick-jumper {
                    display: none;
                }
            }
        `,
        empty: css`
            display: flex;
            align-items: center;
            justify-content: center;
            min-height: 240px;
        `,
        "@media (max-width: 620px)": {
            deviceMain: css`
                align-items: stretch;
                flex-direction: column;

                .name-row {
                    align-items: flex-start;
                }
            `,
            sessionOperate: css`
                justify-content: flex-start;
            `,
            sessionInfo: css`
                grid-template-columns: 1fr;
            `,
            pagerCard: css`
                .ant-pagination {
                    flex-wrap: wrap;
                    row-gap: 6px;
                }

                .ant-pagination-options-quick-jumper {
                    display: none;
                }
            `,
            panelExtra: css`
                .count {
                    display: none;
                }
            `,
        },
    };
});

export default function Index(): React.ReactNode {
    const {message, modal} = App.useApp();
    const user = useModel("user");
    const {styles} = useStyles();
    const [loading, setLoading] = useState(false);
    const [logoutLoading, setLogoutLoading] = useState(false);
    const [list, setList] = useState<any[]>([]);
    const [total, setTotal] = useState(0);
    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);

    const currentSession = user?.web?.session || {};

    const isExpired = (expireTime?: string) => {
        if (!expireTime) {
            return false;
        }
        const time = new Date(expireTime.replace(/-/g, "/")).getTime();
        return Number.isFinite(time) && time <= Date.now();
    };

    const fetchList = async (nextPage = page, nextPageSize = pageSize) => {
        setLoading(true);
        await getSessionList({
            body: {
                page: nextPage,
                page_size: nextPageSize,
                order_by_field: "create_time",
                order_by_type: "desc",
            },
            onSuccess: (res: any) => {
                const data = res?.data || {};
                setList(data?.list || []);
                setTotal(data?.total || 0);
                setPage(data?.page || nextPage);
                setPageSize(data?.page_size || nextPageSize);
            },
            onFail: (res: any) => message.error(res?.message || "获取会话列表失败"),
        }).finally(() => {
            setLoading(false);
        });
    };

    useEffect(() => {
        void fetchList(1, DEFAULT_PAGE_SIZE);
    }, []);

    const handleLogout = (tokens: string[]) => {
        const nextTokens = Array.from(new Set((tokens || []).filter(Boolean)));
        if (!nextTokens.length) {
            message.warning("请选择要注销的会话");
            return;
        }
        modal.confirm({
            title: "注销会话",
            content: nextTokens.length > 1 ? `确定注销选中的 ${nextTokens.length} 个会话吗？` : "确定注销该会话吗？",
            okText: "注销",
            okButtonProps: {danger: true},
            cancelText: "取消",
            mask: {closable: true},
            onOk: async () => {
                setLogoutLoading(true);
                await logoutSession({
                    body: {tokens: nextTokens},
                    onSuccess: (res: any) => {
                        message.success(res?.message || "注销成功");
                        void fetchList(page, pageSize);
                    },
                    onFail: (res: any) => message.error(res?.message || "注销失败"),
                }).finally(() => {
                    setLogoutLoading(false);
                });
            },
        });
    };

    return (
        <Body>
            <div className={styles.page}>
                <Card className={styles.hero} variant="borderless" styles={{body: {padding: 0}}}>
                    <Icon className="hero-watermark" type="LaptopOutlined"/>
                    <div className={styles.heroInner}>
                        <div className={styles.heroContent}>
                            <h1 className={styles.heroTitle}>会话管理</h1>
                            <div className={styles.heroDesc}>集中审阅登录终端，及时注销异常或不再使用的访问会话。</div>
                        </div>
                        <div className={styles.heroVisual}>
                            <div className={styles.visualCard}>
                                <div className="visual-head">
                                    <span className="visual-title">当前会话</span>
                                    <span className="visual-status">在线</span>
                                </div>
                                <Typography.Text className="visual-device" ellipsis={{tooltip: currentSession?.device || "未知设备"}}>
                                    {currentSession?.device || "未知设备"}
                                </Typography.Text>
                                <div className="visual-list">
                                    <div className="visual-row">
                                        <span>IP</span>
                                        <Typography.Text className="value" ellipsis={{tooltip: currentSession?.ip || "-"}}>
                                            {currentSession?.ip || "-"}
                                        </Typography.Text>
                                    </div>
                                    <div className="visual-row">
                                        <span>地点</span>
                                        <Typography.Text className="value" ellipsis={{tooltip: currentSession?.location || "-"}}>
                                            {currentSession?.location || "-"}
                                        </Typography.Text>
                                    </div>
                                    <div className="visual-row">
                                        <span>时间</span>
                                        <Tooltip title={currentSession?.login_time || "-"}>
                                            <Typography.Text className="value" ellipsis>
                                                {currentSession?.login_time ? (ago(currentSession.login_time) || currentSession.login_time) : "-"}
                                            </Typography.Text>
                                        </Tooltip>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </Card>

                <Card
                    className={styles.listCard}
                    variant="borderless"
                    title={(
                        <div className={styles.panelTitle}>
                            <Icon type="LaptopOutlined"/>
                            <span>登录会话</span>
                        </div>
                    )}
                    extra={<div className={styles.panelExtra}><span className="count">共 {total} 条</span></div>}
                >
                    <Spin spinning={loading}>
                        <div className={styles.listBody}>
                            {list.length ? (
                                <div className={styles.sessionGrid}>
                                    {list.map((item) => {
                                        const token = String(item?.token || "");
                                        const expired = isExpired(item?.expire_time);
                                        return (
                                            <div className={styles.sessionCard} key={token || item?.id}>
                                                <div className={styles.deviceMain}>
                                                    <div className="device-content">
                                                        <div className="name-row">
                                                            <Typography.Text className={styles.deviceName} ellipsis={{tooltip: item?.device || "未知设备"}}>
                                                                {item?.device || "未知设备"}
                                                            </Typography.Text>
                                                            <Tag color={expired ? "default" : "success"}>{expired ? "已过期" : "有效"}</Tag>
                                                        </div>
                                                        <div className={styles.deviceSub}>
                                                            <Typography.Text type="secondary" ellipsis={{tooltip: item?.location || "-"}}>
                                                                {item?.location || "-"}
                                                            </Typography.Text>
                                                        </div>
                                                    </div>
                                                    <div className={styles.sessionOperate}>
                                                        <Button
                                                            danger
                                                            type="text"
                                                            size="small"
                                                            loading={logoutLoading}
                                                            onClick={() => handleLogout([token])}
                                                        >
                                                            注销
                                                        </Button>
                                                    </div>
                                                </div>
                                                <div className={styles.sessionInfo}>
                                                    <div className={styles.sessionMeta}>
                                                        <div className="label">IP 地址</div>
                                                        <Typography.Text className="value" copyable={!!item?.ip} ellipsis={{tooltip: item?.ip || "-"}}>
                                                            {item?.ip || "-"}
                                                        </Typography.Text>
                                                    </div>
                                                    <div className={styles.sessionMeta}>
                                                        <div className="label">登录时间</div>
                                                        <Typography.Text className="value" ellipsis={{tooltip: item?.create_time || "-"}}>
                                                            {item?.create_time || "-"}
                                                        </Typography.Text>
                                                    </div>
                                                    <div className={styles.sessionMeta}>
                                                        <div className="label">过期时间</div>
                                                        <Typography.Text className="value" ellipsis={{tooltip: item?.expire_time || "-"}}>
                                                            {item?.expire_time || "-"}
                                                        </Typography.Text>
                                                    </div>
                                                </div>
                                            </div>
                                        );
                                    })}
                                </div>
                            ) : (
                                <div className={styles.empty}>
                                    <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="暂无会话"/>
                                </div>
                            )}
                        </div>
                    </Spin>
                </Card>

                {total > 0 && (
                    <Card className={styles.pagerCard} variant="borderless">
                        <Pagination
                            align="center"
                            size="small"
                            current={page}
                            pageSize={pageSize}
                            total={total}
                            showLessItems
                            showQuickJumper
                            pageSizeOptions={[10, 20, 50, 100]}
                            showSizeChanger={{
                                showSearch: false,
                                variant: "filled",
                                size: "small",
                            }}
                            showTotal={(count, range) => `第 ${range[0]}-${range[1]} 条 / 共 ${count} 条`}
                            onChange={(nextPage, nextPageSize) => void fetchList(nextPage, nextPageSize)}
                        />
                    </Card>
                )}
            </div>
        </Body>
    );
}
