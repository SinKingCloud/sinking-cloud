import React, {useRef, useState} from "react";
import {useModel} from "umi";
import {App, Badge, Button, Card, Tooltip, Typography} from "antd";
import {createStyles} from "antd-style";
import type {PageTableProps, PageTableRef} from "sinking-antd";
import {Body, Icon, PageTable, useTheme} from "sinking-antd";
import {getSessionList, logoutSession} from "@/service/user/person";
import {ago} from "@/utils/time";
import {getParams} from "@/utils/page";

type UserSession = {
    token: string;
    device?: string;
    ip?: string;
    location?: string;
    create_time?: string;
    expire_time?: string;
    update_time?: string;
};

const useStyles = createStyles(({css, token, isDarkMode}, {isCompactMode}: {isCompactMode: boolean}) => ({
    page: css`
        width: 100%;
        max-width: 1450px;
        margin: 0 auto;
    `,
    header: css`
        position: relative;
        display: grid;
        grid-template-columns: minmax(0, 1fr) 300px;
        align-items: center;
        gap: ${isCompactMode ? 16 : 20}px;
        min-height: ${isCompactMode ? 159 : 185}px;
        padding: ${isCompactMode ? 20 : 26}px;
        overflow: hidden;
        background: ${isDarkMode
            ? `linear-gradient(135deg, color-mix(in srgb, ${token.colorPrimary} 9%, ${token.colorBgContainer}) 0%, ${token.colorBgContainer} 100%)`
            : token.colorFillAlter};
        isolation: isolate;

        @media (max-width: 900px) {
            grid-template-columns: 1fr;
            min-height: auto;
            gap: ${isCompactMode ? 14 : 16}px;
            padding: ${isCompactMode ? 18 : 20}px;
        }
    `,
    headerMain: css`
        position: relative;
        z-index: 1;
        min-width: 0;
        display: flex;
        flex-direction: column;
        gap: ${isCompactMode ? 8 : 10}px;

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

        @media (max-width: 560px) {
            h1 {
                font-size: ${isCompactMode ? 22 : 23}px;
                line-height: ${isCompactMode ? 30 : 31}px;
            }

            p {
                font-size: 12px;
                line-height: 20px;
            }
        }
    `,
    headerDecor: css`
        position: absolute;
        z-index: 0;
        inset: 0;
        overflow: hidden;
        pointer-events: none;

        .header-ambient,
        .header-shine,
        .header-watermark {
            position: absolute;
            pointer-events: none;
        }

        .header-ambient {
            inset: 0;
            z-index: 0;
        }

        .header-ambient::before {
            content: '';
            position: absolute;
            inset: -78px;
            background:
                linear-gradient(135deg,
                    ${isDarkMode ? `color-mix(in srgb, ${token.colorPrimary} 42%, transparent)` : token.colorPrimaryBgHover} 0%,
                    ${isDarkMode ? `color-mix(in srgb, ${token.colorPrimary} 22%, transparent)` : token.colorPrimaryBg} 44%,
                    transparent 76%),
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
            z-index: 1;
            opacity: ${isDarkMode ? 0.20 : 0.34};
            background: linear-gradient(110deg, transparent 16%,
                ${isDarkMode ? "rgba(34,211,238,0.18)" : "rgba(255,255,255,0.70)"} 58%, transparent 96%);
            mix-blend-mode: ${isDarkMode ? "screen" : "normal"};
        }

        .header-watermark {
            right: -78px;
            bottom: -104px;
            z-index: 1;
            color: ${token.colorPrimary};
            font-size: 320px;
            line-height: 1;
            opacity: ${isDarkMode ? 0.17 : 0.14};
            transform: rotate(-24deg);
            filter: blur(20px);
        }

        @media (max-width: 900px) {
            .header-watermark {
                right: -62px;
                bottom: -90px;
                font-size: 260px;
            }
        }

        @media (max-width: 560px) {
            .header-ambient::before {
                inset: -54px;
            }

            .header-watermark {
                right: -48px;
                bottom: -68px;
                font-size: 210px;
                filter: blur(18px);
            }
        }
    `,
    table: css`
        .session-source {
            display: flex;
            flex-direction: column;
            align-items: flex-start;
            gap: 2px;
            min-width: 0;
        }

        .session-source .ant-typography,
        .session-device {
            max-width: 100%;
            font-size: ${token.fontSizeSM + 1}px;
            line-height: 20px;
        }

        .session-location.ant-typography,
        .session-time {
            font-size: ${token.fontSizeSM}px;
            color: ${token.colorTextSecondary};
        }

        .session-time-row {
            display: flex;
            align-items: center;
            gap: 4px;
            line-height: 20px;
            white-space: nowrap;
            font-variant-numeric: tabular-nums;
        }

        .session-time-label {
            color: ${token.colorTextTertiary};
        }

        .ant-typography-copy {
            margin-inline-start: 6px;
            font-size: ${token.fontSizeSM}px;
            color: ${token.colorTextTertiary};
        }

        .ant-typography-copy:hover {
            color: ${token.colorPrimary};
        }
    `,
    currentSession: css`
        position: relative;
        z-index: 1;
        container: current-session-card / inline-size;
        width: 100%;
        overflow: hidden;
        background: ${isDarkMode ? "rgba(255,255,255,0.035)" : "rgba(255,255,255,0.72)"};
        font-size: ${token.fontSizeSM}px;
        line-height: 18px;

        &.ant-card:not(.ant-card-bordered) {
            border: none;
            border-radius: ${token.borderRadiusLG}px;
            box-shadow: ${isDarkMode ? "none" : "0 12px 28px rgba(15, 23, 42, 0.045)"};
        }

        .ant-card-body {
            padding: ${token.paddingSM + 4}px;
        }

        .current-summary {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 14px;
        }

        .current-symbol {
            display: flex;
            align-items: center;
            justify-content: center;
            flex: none;
            width: ${isCompactMode ? 56 : 64}px;
            height: ${isCompactMode ? 45 : 51}px;
            color: ${token.colorPrimary};
        }

        .current-symbol svg {
            display: block;
            width: 100%;
            height: 100%;
        }

        .current-symbol .connection-surface {
            fill: ${token.colorBgContainer};
        }

        .current-identity {
            flex: 1;
            min-width: 0;
        }

        .current-label {
            display: inline-flex;
            align-items: center;
            gap: 6px;
            color: ${token.colorTextTertiary};
            line-height: 16px;
        }

        .current-device.ant-typography {
            display: block;
            width: 100%;
            margin-top: 4px;
            color: ${token.colorText};
            font-size: ${token.fontSize + 1}px;
            font-weight: 600;
            line-height: 22px;
        }

        .current-source {
            display: grid;
            grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
            gap: 12px;
            margin-top: 12px;
        }

        .current-source-item {
            display: flex;
            align-items: center;
            gap: 6px;
            min-width: 0;
        }

        .current-source-item .anticon {
            flex: none;
            color: ${token.colorTextTertiary};
        }

        .current-source .ant-typography {
            flex: 1;
            min-width: 0;
            color: ${token.colorTextSecondary};
            font-size: inherit;
            line-height: inherit;
        }

        .current-location {
            display: flex;
            align-items: center;
            gap: 6px;
            margin-top: 8px;
            color: ${token.colorTextTertiary};
            line-height: 16px;
        }

        .current-location .anticon {
            flex: none;
        }

        .current-location span:last-child {
            min-width: 0;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        @container current-session-card (min-width: 480px) {
            .ant-card-body {
                display: grid;
                grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
                align-items: center;
                gap: 4px 28px;
            }

            .current-summary {
                grid-row: 1 / 3;
                gap: 12px;
            }

            .current-symbol {
                order: -1;
            }

            .current-source {
                margin-top: 0;
            }

            .current-location {
                margin-top: 0;
            }
        }
    `,
}));

export default function Index(): React.ReactNode {
    const {message, modal} = App.useApp();
    const user = useModel("user");
    const theme = useTheme();
    const {styles} = useStyles({isCompactMode: theme?.isCompactTheme?.() || false});
    const tableRef = useRef<PageTableRef<UserSession> | null>(null);
    const [logoutLoading, setLogoutLoading] = useState(false);
    const currentSession = user?.web?.session || {};

    const handleLogout = (tokens: React.Key[]) => {
        const nextTokens = Array.from(new Set((tokens || []).map((token) => String(token || "")).filter(Boolean)));
        if (!nextTokens.length) {
            message.warning("请选择要注销的会话");
            return;
        }
        modal.confirm({
            title: "注销会话",
            content: nextTokens.length > 1 ? `确定注销选中的 ${nextTokens.length} 个会话吗？` : "确定注销该会话吗？",
            okText: "注销",
            okButtonProps: {danger: true, type: "default"},
            cancelText: "取消",
            mask: {closable: true},
            onOk: async () => {
                setLogoutLoading(true);
                await logoutSession({
                    body: {tokens: nextTokens},
                    onSuccess: (res: any) => {
                        message.success(res?.message || "注销成功");
                        tableRef.current?.reload();
                        tableRef.current?.clearSelectedRows();
                    },
                    onFail: (res: any) => message.error(res?.message || "注销失败"),
                }).finally(() => setLogoutLoading(false));
            },
        });
    };

    const columns: PageTableProps<UserSession>["columns"] = [
        {
            title: "登录设备",
            dataIndex: "device",
            width: 260,
            ellipsis: {showTitle: false},
            render: (value: string) => (
                <Typography.Text className="session-device" ellipsis={{tooltip: value || "未知设备"}}>
                    {value || "未知设备"}
                </Typography.Text>
            ),
        },
        {
            title: "登录来源",
            dataIndex: "ip",
            width: 220,
            render: (value: string, record: UserSession) => (
                <div className="session-source">
                    <Typography.Text copyable={value ? {text: value} : false} ellipsis={{tooltip: value || "-"}}>
                        {value || "-"}
                    </Typography.Text>
                    <Typography.Text className="session-location" type="secondary"
                                     ellipsis={{tooltip: record.location || "-"}}>
                        {record.location || "-"}
                    </Typography.Text>
                </div>
            ),
        },
        {
            title: "相关时间",
            dataIndex: "create_time",
            width: 220,
            render: (_: unknown, record: UserSession) => (
                <div className="session-source session-time">
                    <div className="session-time-row">
                        <span className="session-time-label">登录于:</span>
                        <span>{record.create_time || "-"}</span>
                    </div>
                    <div className="session-time-row">
                        <span className="session-time-label">过期于:</span>
                        <span>{record.expire_time || "-"}</span>
                    </div>
                </div>
            ),
        },
        {
            title: "操作",
            key: "option",
            width: 80,
            fixed: "right",
            className: "action-cell",
            render: (_: unknown, record: UserSession) => (
                <Button size="small" danger loading={logoutLoading} onClick={() => handleLogout([record.token])}>
                    注销
                </Button>
            ),
        },
    ];

    return (
        <Body>
            <PageTable<UserSession>
                ref={tableRef}
                ariaLabel="会话管理"
                rootClassName={`${styles.page} ${styles.table}`}
                rowKey="token"
                columns={columns}
                scroll={{x: 824}}
                defaultPageSize={10}
                hero={<section className={styles.header}>
                    <div className={styles.headerDecor} aria-hidden="true">
                        <span className="header-ambient"/>
                        <span className="header-shine"/>
                        <span className="header-watermark"><Icon type="LaptopOutlined"/></span>
                    </div>
                    <div className={styles.headerMain}>
                        <h1>会话管理</h1>
                        <p>查看账户的登录设备与访问记录，及时管理已登录会话，保障账户访问安全。</p>
                    </div>
                    <Card size="small" variant="borderless" className={styles.currentSession} aria-label="当前会话">
                        <div className="current-summary">
                            <div className="current-identity">
                                <span className="current-label"><Badge status="success"/>当前会话</span>
                                <Typography.Text className="current-device" ellipsis={{tooltip: currentSession.device || "未知设备"}}>
                                    {currentSession.device || "未知设备"}
                                </Typography.Text>
                            </div>
                            <div className="current-symbol">
                                <svg viewBox="0 0 80 64" fill="none" aria-hidden="true">
                                    {/iOS|iPhone|iPad|Android|Mobile|Phone/i.test(currentSession.device || "") ? <>
                                        <rect x="23" y="5" width="31" height="53" rx="7" fill="currentColor" fillOpacity="0.06" stroke="currentColor" strokeOpacity="0.42" strokeWidth="1.5"/>
                                        <rect x="27" y="12" width="23" height="37" rx="3" stroke="currentColor" strokeOpacity="0.14"/>
                                        <path d="M34 9H43M35 53H42" stroke="currentColor" strokeOpacity="0.44" strokeWidth="1.5" strokeLinecap="round"/>
                                        <path d="M32 22H43M32 28H40" stroke="currentColor" strokeOpacity="0.22" strokeWidth="1.5" strokeLinecap="round"/>
                                    </> : <>
                                        <rect x="9" y="11" width="53" height="36" rx="5" fill="currentColor" fillOpacity="0.06" stroke="currentColor" strokeOpacity="0.42" strokeWidth="1.5"/>
                                        <rect x="13" y="15" width="45" height="28" rx="2.5" stroke="currentColor" strokeOpacity="0.14"/>
                                        <path d="M5 48H66L62 53H9L5 48Z" fill="currentColor" fillOpacity="0.1" stroke="currentColor" strokeOpacity="0.38" strokeWidth="1.5" strokeLinejoin="round"/>
                                        <path d="M29 48H41M20 25H43M20 32H35" stroke="currentColor" strokeOpacity="0.25" strokeWidth="1.5" strokeLinecap="round"/>
                                    </>}
                                    <circle className="connection-surface" cx="63" cy="45" r="12"/>
                                    <circle cx="63" cy="45" r="12" fill="currentColor" fillOpacity="0.08" stroke="currentColor" strokeOpacity="0.38" strokeWidth="1.4"/>
                                    <path d="M61 46.5L60 47.5A2.8 2.8 0 0 1 56 43.5L59 40.5A2.8 2.8 0 0 1 63 40.5M65 43.5L66 42.5A2.8 2.8 0 0 1 70 46.5L67 49.5A2.8 2.8 0 0 1 63 49.5M60.5 47.5L65.5 42.5" stroke="currentColor" strokeOpacity="0.64" strokeWidth="1.5" strokeLinecap="round"/>
                                </svg>
                            </div>
                        </div>
                        <div className="current-source">
                            <div className="current-source-item">
                                <Icon type="GlobalOutlined"/>
                                <Typography.Text ellipsis={{tooltip: `登录 IP: ${currentSession.ip || "-"}`}}>
                                    {currentSession.ip || "-"}
                                </Typography.Text>
                            </div>
                            <Tooltip title={`登录时间: ${currentSession.login_time || "-"}`}>
                                <div className="current-source-item">
                                    <Icon type="ClockCircleOutlined"/>
                                    <Typography.Text ellipsis>
                                        {currentSession.login_time ? (ago(currentSession.login_time) || currentSession.login_time) : "-"}登录
                                    </Typography.Text>
                                </div>
                            </Tooltip>
                        </div>
                        <Tooltip title={`登录地点: ${currentSession.location || "-"}`}>
                            <div className="current-location">
                                <Icon type="EnvironmentOutlined"/>
                                <span>{currentSession.location || "-"}</span>
                            </div>
                        </Tooltip>
                    </Card>
                </section>}
                search={{
                    ariaLabel: "搜索会话",
                    placeholder: "搜索设备、IP 或登录地点",
                    maxLength: 100,
                    fields: [
                        {value: "keyword", label: "全部", placeholder: "搜索设备、IP 或登录地点"},
                        {value: "device", label: "设备", placeholder: "搜索登录设备"},
                        {value: "ip", label: "IP", placeholder: "搜索登录 IP", maxLength: 50},
                        {value: "location", label: "地点", placeholder: "搜索登录地点"},
                    ],
                }}
                filters={[{type: "date", name: "create_time", label: "登录时间筛选"}]}
                toolbar={{refresh: {ariaLabel: "刷新会话列表"}}}
                rowSelection={{
                    fixed: true,
                    columnWidth: 44,
                    actions: (keys) => [{
                        key: "logout",
                        type: "button",
                        label: "批量注销",
                        icon: "LogoutOutlined",
                        danger: true,
                        disabled: logoutLoading,
                        onClick: () => handleLogout(keys),
                    }],
                }}
                request={async (params, sort) => {
                    const response = await getSessionList({body: getParams(params, sort)});
                    if (response?.code !== 200) {
                        throw new Error(response?.message || "获取会话列表失败");
                    }
                    return {data: response.data?.list, total: response.data?.total ?? 0, success: true};
                }}
                onRequestError={(error) => message.error(error.message)}
                pagination={{unit: "条"}}
                paginationAffix={{offsetBottom: 15}}
            />
        </Body>
    );
}
