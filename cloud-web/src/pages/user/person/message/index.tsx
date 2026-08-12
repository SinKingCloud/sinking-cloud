import React, {useEffect, useRef, useState} from "react";
import {App, Button, Card, Dropdown, Space, Spin, Tabs, Tag, Tooltip} from "antd";
import {createStyles} from "antd-style";
import {Body, ProModal, ProModalRef, ProTable, ProTableRef, Title} from "sinking-antd";
import {getData} from "@/utils/page";
import {ago} from "@/utils/time";
import {
    clearMessage,
    deleteMessage,
    getMessageInfo,
    getMessageList,
    readAllMessage,
    readMessage,
} from "@/service/user/message";

const MESSAGE_STATUS_UNREAD = 0;
const MESSAGE_STATUS_READ = 1;
const MESSAGE_TAB_ALL = "all";
const MESSAGE_TAB_UNREAD = "unread";
const MESSAGE_TAB_READ = "read";
const statusText: Record<string, string> = {"0": "未读", "1": "已读"};
const messageTabs = [
    {key: MESSAGE_TAB_ALL, label: "全部消息"},
    {key: MESSAGE_TAB_UNREAD, label: "未读消息"},
    {key: MESSAGE_TAB_READ, label: "已读消息"},
];

type UserMessage = {
    id: number | string;
    title: string;
    content?: string;
    status: number;
    read_time?: string;
    create_time?: string;
    update_time?: string;
};

const useStyles = createStyles(({css, token, isDarkMode}: any): any => ({
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
    messageTitle: css`
        display: flex;
        align-items: center;
        gap: 8px;
        min-width: 0;
        cursor: pointer;

        span:last-child {
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .unread-dot {
            width: 7px;
            height: 7px;
            border-radius: 50%;
            flex: 0 0 auto;
            background: ${token.colorTextQuaternary};
        }

        .unread-dot.unread {
            background: ${token.colorPrimary};
            box-shadow: 0 0 0 3px ${isDarkMode ? "rgba(22,119,255,0.16)" : token.colorPrimaryBg};
        }
    `,
    statusTag: css`
        font-size: 12px;
        line-height: 18px;
        margin-inline-end: 0;
    `,
    timeText: css`
        color: ${token.colorTextSecondary};
        font-size: 12px;
    `,
    modalBody: css`
        padding: 0 4px;
    `,
    detailHeader: css`
        text-align: center;
        padding: 16px 0 12px;
        margin-bottom: 20px;
        border-bottom: 1px solid ${isDarkMode ? "rgba(255,255,255,0.08)" : "rgba(0,0,0,0.06)"};

        h2 {
            margin: 0;
            color: ${token.colorTextHeading};
            font-size: 20px;
            font-weight: 600;
            line-height: 1.4;
        }

        .meta {
            display: flex;
            align-items: center;
            justify-content: center;
            flex-wrap: wrap;
            gap: 8px 18px;
            margin-top: 10px;
            color: ${token.colorTextTertiary};
            font-size: 12px;
        }
    `,
    detailContent: css`
        padding: 0 4px 12px;
        color: ${token.colorText};
        font-size: 14px;
        line-height: 1.8;
        white-space: pre-wrap;
        word-break: break-word;

        p {
            margin: 0 0 12px;

            &:last-child {
                margin-bottom: 0;
            }
        }

        img {
            display: block;
            max-width: 100%;
            margin: 12px 0;
            border-radius: ${token.borderRadius}px;
        }

        a {
            color: ${token.colorPrimary};
        }
    `,
    spinContainer: css`
        display: flex;
        justify-content: center;
        align-items: center;
        padding: 40px 0;
    `,
}));

const isUnread = (value: any) => Number(value) === MESSAGE_STATUS_UNREAD;

const getStatusText = (value: any) => {
    if (value === undefined || value === null || value === "") {
        return "-";
    }
    return statusText[String(value)] || (isUnread(value) ? "未读" : "已读");
};

const padTime = (value: number) => String(value).padStart(2, "0");
const formatDateTime = (value: Date) => `${value.getFullYear()}-${padTime(value.getMonth() + 1)}-${padTime(value.getDate())} ${padTime(value.getHours())}:${padTime(value.getMinutes())}:${padTime(value.getSeconds())}`;
const parseDateTime = (value: string) => {
    const date = new Date(value.includes("T") ? value : value.replace(/-/g, "/"));
    return Number.isFinite(date.getTime()) ? date : undefined;
};

const normalizeTime = (value: any) => {
    if (!value) {
        return "";
    }
    if (value instanceof Date) {
        return formatDateTime(value);
    }
    if (typeof value?.format === "function") {
        return value.format("YYYY-MM-DD HH:mm:ss");
    }
    const time = String(value);
    const date = parseDateTime(time);
    return date ? formatDateTime(date) : time;
};

const formatAgoTime = (value: any) => {
    const time = normalizeTime(value);
    if (!time) {
        return "-";
    }
    return ago(time) || time;
};

export default (): React.ReactNode => {
    const {message, modal} = App.useApp();
    const {styles} = useStyles();
    const tableRef = useRef<ProTableRef>({} as ProTableRef);
    const detailModalRef = useRef<ProModalRef>({} as ProModalRef);
    const statusReadyRef = useRef(false);
    const [activeStatus, setActiveStatus] = useState(MESSAGE_TAB_ALL);
    const [detailLoading, setDetailLoading] = useState(false);
    const [detail, setDetail] = useState<UserMessage>();

    const refreshTable = (clearSelected = false) => {
        tableRef.current?.refreshTableData?.();
        if (clearSelected) {
            tableRef.current?.clearSelectedRows?.();
        }
    };

    const openDetail = async (record: UserMessage) => {
        setDetail(record);
        detailModalRef.current?.show();
        setDetailLoading(true);
        const resp = await getMessageInfo({body: {id: record?.id}});
        setDetailLoading(false);
        if (resp?.code !== 200) {
            message?.error(resp?.message || "获取站内消息详情失败");
            detailModalRef.current?.hide();
            return;
        }
        setDetail(resp?.data || record);
        if (isUnread(record?.status)) {
            refreshTable();
        }
    };

    const markRead = async (ids: any[]) => {
        const targetIds = Array.from(new Set((ids || []).map((item) => Number(item)).filter((item) => item > 0)));
        if (targetIds.length <= 0) {
            message?.warning("请选择消息");
            return;
        }
        modal.confirm({
            title: "批量已读",
            content: `确定将选中的 ${targetIds.length} 条站内消息标记为已读吗？`,
            okText: "确定",
            cancelText: "取消",
            onOk: async () => {
                const resp = await readMessage({body: {ids: targetIds}});
                if (resp?.code !== 200) {
                    message?.error(resp?.message || "操作失败");
                    return;
                }
                message?.success(resp?.message || "操作成功");
                refreshTable(true);
            },
        } as any);
    };

    const markAllRead = () => {
        modal.confirm({
            title: "一键已读",
            content: "确定将所有站内消息标记为已读吗？",
            okText: "确定",
            cancelText: "取消",
            onOk: async () => {
                const resp = await readAllMessage();
                if (resp?.code !== 200) {
                    message?.error(resp?.message || "操作失败");
                    return;
                }
                message?.success(resp?.message || "操作成功");
                refreshTable(true);
            },
        } as any);
    };

    const removeMessages = (ids: any[]) => {
        const targetIds = Array.from(new Set((ids || []).map((item) => Number(item)).filter((item) => item > 0)));
        if (targetIds.length <= 0) {
            message?.warning("请选择消息");
            return;
        }
        modal.confirm({
            title: "删除站内消息",
            content: `确定删除选中的 ${targetIds.length} 条站内消息吗？`,
            okText: "删除",
            okButtonProps: {danger: true, type: "default"},
            cancelText: "取消",
            onOk: async () => {
                const resp = await deleteMessage({body: {ids: targetIds}});
                if (resp?.code !== 200) {
                    message?.error(resp?.message || "删除失败");
                    return;
                }
                message?.success(resp?.message || "删除成功");
                refreshTable(true);
            },
        } as any);
    };

    const clearAllMessages = () => {
        modal.confirm({
            title: "清空站内消息",
            content: "确定清空所有站内消息吗？",
            okText: "清空",
            okButtonProps: {danger: true, type: "default"},
            cancelText: "取消",
            onOk: async () => {
                const resp = await clearMessage();
                if (resp?.code !== 200) {
                    message?.error(resp?.message || "清空失败");
                    return;
                }
                message?.success(resp?.message || "清空成功");
                refreshTable(true);
            },
        } as any);
    };

    useEffect(() => {
        if (!statusReadyRef.current) {
            statusReadyRef.current = true;
            return;
        }
        refreshTable(true);
    }, [activeStatus]);

    const columns = [
        {
            title: "标题",
            dataIndex: "title",
            width: 220,
            render: (_: any, record: UserMessage) => (
                <div className={styles.messageTitle} onClick={() => openDetail(record)}>
                    <span className={`unread-dot ${isUnread(record?.status) ? "unread" : ""}`}/>
                    <span>{record?.title || "未命名消息"}</span>
                </div>
            ),
        },
        {
            title: "创建于",
            dataIndex: "create_time",
            width: 120,
            render: (value: any, record: UserMessage) => {
                const time = normalizeTime(record?.create_time || value);
                return (
                    <Tooltip title={time || "-"}>
                        <span className={styles.timeText}>{formatAgoTime(time)}</span>
                    </Tooltip>
                );
            },
        },
        {
            title: "阅读于",
            dataIndex: "read_time",
            width: 120,
            render: (value: any, record: UserMessage) => {
                const time = normalizeTime(record?.read_time || value);
                return (
                    <Tooltip title={time || "-"}>
                        <span className={styles.timeText}>{formatAgoTime(time)}</span>
                    </Tooltip>
                );
            },
        },
        {
            title: "状态",
            dataIndex: "status",
            width: 96,
            render: (_: any, record: UserMessage) => (
                <Tag className={styles.statusTag} color={isUnread(record?.status) ? "blue" : "default"}>
                    {getStatusText(record?.status)}
                </Tag>
            ),
        },
        {
            title: "操作",
            valueType: "option",
            key: "option",
            width: 86,
            render: (_: any, record: UserMessage) => {
                const items = [
                    {
                        key: "info",
                        label: "详情",
                        onClick: () => openDetail(record),
                    },
                    {
                        key: "delete",
                        danger: true,
                        label: "删除",
                        onClick: () => removeMessages([record?.id]),
                    },
                ];
                return (
                    <Dropdown menu={{items}} trigger={["click"]} placement="bottom" arrow>
                        <Button size="small">操作</Button>
                    </Dropdown>
                );
            },
        },
    ] as any;

    return (
        <Body>
            <div className={styles.page}>
                <Card className={styles.navCard} variant="borderless">
                    <Tabs
                        className={styles.tabs}
                        activeKey={activeStatus}
                        onChange={(value) => {
                            if (value !== activeStatus) {
                                setActiveStatus(value);
                            }
                        }}
                        items={messageTabs as any}
                    />
                </Card>
                <ProTable
                    ref={tableRef}
                    rowKey="id"
                    title={<Title>站内消息</Title>}
                    extra={(
                        <Dropdown
                            trigger={["click"]}
                            placement="bottom"
                            arrow
                            menu={{
                                items: [
                                    {
                                        key: "read-all",
                                        label: "一键已读",
                                        onClick: markAllRead,
                                    },
                                    {
                                        key: "clear",
                                        danger: true,
                                        label: "清空消息",
                                        onClick: clearAllMessages,
                                    },
                                ],
                            }}
                        >
                            <Button size="small" type="primary">操作</Button>
                        </Dropdown>
                    )}
                    extraRefreshBtn={true}
                    search={false}
                    columns={columns}
                    defaultPage={1}
                    defaultPageSize={10}
                    rowSelection={{
                        fixed: true,
                        rightExtra: (
                            <Space wrap>
                                <Button
                                    key="read"
                                    type="primary"
                                    ghost
                                    onClick={() => markRead(tableRef.current?.getSelectedRowKeys?.() as any)}
                                >
                                    批量已读
                                </Button>
                                <Button
                                    key="delete"
                                    type="primary"
                                    danger
                                    ghost
                                    onClick={() => removeMessages(tableRef.current?.getSelectedRowKeys?.() as any)}
                                >
                                    批量删除
                                </Button>
                            </Space>
                        ),
                    }}
                    request={(params: any) => {
                        const status = activeStatus === MESSAGE_TAB_UNREAD
                            ? String(MESSAGE_STATUS_UNREAD)
                            : activeStatus === MESSAGE_TAB_READ
                                ? String(MESSAGE_STATUS_READ)
                                : undefined;
                        return getData({
                            ...(status !== undefined ? {status} : {}),
                            ...params,
                        }, {}, getMessageList);
                    }}
                    paginationAffix={true}
                    selectionAffix={true}
                    {...({
                        tableProps: {
                            scroll: {x: "max-content"},
                            rowClassName: (record: UserMessage) => isUnread(record?.status) ? "unread-row" : "read-row",
                        },
                    } as any)}
                />
            </div>
            <ProModal
                ref={detailModalRef}
                title={<Title>消息详情</Title>}
                width={560}
                onCancel={() => detailModalRef.current?.hide()}
                afterClose={() => setDetail(undefined)}
                modalProps={{
                    footer: null,
                    destroyOnHidden: true,
                } as any}
            >
                {detailLoading ? (
                    <div className={styles.spinContainer}>
                        <Spin/>
                    </div>
                ) : (
                    <div className={styles.modalBody}>
                        <div className={styles.detailHeader}>
                            <h2>{detail?.title || "消息详情"}</h2>
                            <div className="meta">
                                <span>创建时间：{detail?.create_time || "-"}</span>
                            </div>
                        </div>
                        <div
                            className={styles.detailContent}
                            dangerouslySetInnerHTML={{__html: detail?.content || "暂无消息内容"}}
                        />
                    </div>
                )}
            </ProModal>
        </Body>
    );
};
