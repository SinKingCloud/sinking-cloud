import React, {useRef, useState} from "react";
import {App, Button, Dropdown, Spin, Tag, Tooltip} from "antd";
import {createStyles} from "antd-style";
import type {PageTableProps, PageTableRef, ProModalRef} from "sinking-antd";
import {Body, Icon, PageTable, ProModal, Title, useTheme} from "sinking-antd";
import {getParams} from "@/utils/page";
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
const statusText: Record<string, string> = {"0": "未读", "1": "已读"};

type UserMessage = {
    id: number | string;
    title: string;
    content?: string;
    status: number;
    read_time?: string;
    create_time?: string;
    update_time?: string;
};

const useStyles = createStyles(({css, token, isDarkMode}: any, {isCompactMode}: {isCompactMode: boolean}): any => ({
    page: css`
        width: 100%;
        max-width: 1450px;
        margin: 0 auto;
    `,
    header: css`
        position: relative;
        display: flex;
        align-items: center;
        min-height: ${isCompactMode ? 159 : 185}px;
        padding: ${isCompactMode ? 20 : 26}px;
        overflow: hidden;
        background: ${isDarkMode
            ? `linear-gradient(135deg, color-mix(in srgb, ${token.colorPrimary} 9%, ${token.colorBgContainer}) 0%, ${token.colorBgContainer} 100%)`
            : token.colorFillAlter};
        isolation: isolate;

        @media (max-width: 900px) {
            min-height: auto;
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
    messageTitle: css`
        display: inline-flex;
        align-items: center;
        gap: 10px;
        width: 100%;
        min-width: 0;
        padding: 0;
        border: 0;
        background: transparent;
        color: ${token.colorTextSecondary};
        font: inherit;
        line-height: 22px;
        text-align: left;
        cursor: pointer;

        &.unread {
            color: ${token.colorText};
        }

        &:hover {
            color: ${token.colorPrimary};
        }

        &:focus-visible {
            outline: 2px solid ${token.colorPrimaryBorder};
            outline-offset: 4px;
            border-radius: ${token.borderRadiusSM}px;
        }

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
            background: transparent;
        }

        .unread-dot.unread {
            background: ${token.colorPrimary};
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
    const theme = useTheme();
    const {styles} = useStyles({isCompactMode: theme?.isCompactTheme?.() || false});
    const tableRef = useRef<PageTableRef<UserMessage> | null>(null);
    const detailModalRef = useRef<ProModalRef | null>(null);
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
            title: "全部已读",
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

    const columns: PageTableProps<UserMessage>["columns"] = [
        {
            title: "标题",
            dataIndex: "title",
            width: 320,
            ellipsis: {showTitle: false},
            render: (_: any, record: UserMessage) => (
                <Tooltip title={record.title || "未命名消息"}>
                    <button
                        type="button"
                        className={`${styles.messageTitle} ${isUnread(record.status) ? "unread" : ""}`}
                        onClick={() => openDetail(record)}
                    >
                        <span className={`unread-dot ${isUnread(record.status) ? "unread" : ""}`} aria-hidden="true"/>
                        <span>{record.title || "未命名消息"}</span>
                    </button>
                </Tooltip>
            ),
        },
        {
            title: "创建于",
            dataIndex: "create_time",
            width: 120,
            sorter: true,
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
            sorter: true,
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
            key: "option",
            width: 80,
            align: "center",
            fixed: "right",
            className: "action-cell",
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
                    <Dropdown menu={{items}} trigger={["click"]} placement="bottom" arrow={{pointAtCenter: true}}>
                        <Button size="small" aria-label="消息操作">操作</Button>
                    </Dropdown>
                );
            },
        },
    ];

    return (
        <Body>
            <PageTable<UserMessage>
                ref={tableRef}
                ariaLabel="站内消息"
                rootClassName={styles.page}
                rowKey="id"
                columns={columns}
                scroll={{x: 760}}
                defaultPage={1}
                defaultPageSize={10}
                defaultSort={{id: "descend"}}
                hero={<section className={styles.header}>
                    <div className={styles.headerDecor} aria-hidden="true">
                        <span className="header-ambient"/>
                        <span className="header-shine"/>
                        <span className="header-watermark"><Icon type="MailOutlined"/></span>
                    </div>
                    <div className={styles.headerMain}>
                        <h1>站内消息</h1>
                        <p>账户通知、业务提醒与系统消息集中呈现，重要动态及时掌握。</p>
                    </div>
                </section>}
                search={{
                    ariaLabel: "搜索消息",
                    placeholder: "搜索消息标题或内容",
                    maxLength: 100,
                    fields: [
                        {value: "keyword", label: "全部", placeholder: "搜索消息标题或内容"},
                        {value: "title", label: "标题", placeholder: "搜索消息标题"},
                        {value: "content", label: "内容", placeholder: "搜索消息内容"},
                    ],
                }}
                filters={[
                    {
                        type: "select",
                        name: "status",
                        label: "消息状态筛选",
                        icon: "MailOutlined",
                        allLabel: "全部消息",
                        options: [
                            {value: String(MESSAGE_STATUS_UNREAD), label: "未读消息"},
                            {value: String(MESSAGE_STATUS_READ), label: "已读消息"},
                        ],
                        onChange: () => tableRef.current?.clearSelectedRows(),
                    },
                ]}
                toolbar={{
                    refresh: {ariaLabel: "刷新消息列表"},
                    actions: [
                        {
                            key: "operate",
                            type: "button",
                            label: "操作消息",
                            icon: "SettingOutlined",
                            menu: {
                                items: [
                                    {key: "read-all", label: "全部已读", onClick: markAllRead},
                                    {key: "clear", label: "清空消息", danger: true, onClick: clearAllMessages},
                                ],
                            },
                        },
                    ],
                }}
                rowSelection={{
                    fixed: true,
                    actions: (keys) => [
                        {
                            key: "read",
                            type: "button",
                            label: "批量已读",
                            icon: "CheckOutlined",
                            onClick: () => markRead(keys),
                        },
                        {
                            key: "delete",
                            type: "button",
                            label: "批量删除",
                            icon: "DeleteOutlined",
                            danger: true,
                            onClick: () => removeMessages(keys),
                        },
                    ],
                }}
                request={async (params, sort) => {
                    const response = await getMessageList({body: getParams(params, sort)});
                    if (response?.code !== 200) {
                        throw new Error(response?.message || "获取站内消息列表失败");
                    }
                    return {data: response.data?.list, total: response.data?.total ?? 0, success: true};
                }}
                onRequestError={(error) => message.error(error.message)}
                pagination={{unit: "条"}}
                paginationAffix={{offsetBottom: 15}}
            />
            <ProModal
                ref={detailModalRef}
                title={<Title>消息详情</Title>}
                width={560}
                onCancel={() => detailModalRef.current?.hide()}
                afterClose={() => setDetail(undefined)}
                modalProps={{
                    footer: null,
                    destroyOnHidden: true,
                    mask: {closable: true},
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
