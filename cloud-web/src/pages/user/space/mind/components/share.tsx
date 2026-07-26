import React, {useEffect, useRef} from "react";
import type {MenuProps} from "antd";
import {Button, Drawer, Dropdown, Empty, Spin, Tooltip} from "antd";
import {DeleteOutlined, EditOutlined, LinkOutlined, MoreOutlined, PlusOutlined} from "@ant-design/icons";
import {SHARE_ENABLE} from "../constants";
import type {MindDocument, MindShareRecord} from "../types";
import {getTimeText} from "../utils";

export type ShareProps = {
    styles: Record<string, string>;
    open: boolean;
    document?: MindDocument;
    items: MindShareRecord[];
    loading: boolean;
    loadFailed: boolean;
    onClose: () => void;
    onCreate: () => void;
    onLoadMore: (documentId: string) => void | Promise<void>;
    onCopyLink: (record: MindShareRecord) => void | Promise<void>;
    onEdit: (record: MindShareRecord) => void;
    onRemove: (record: MindShareRecord) => void;
};

const Share: React.FC<ShareProps> = ({
    styles,
    open,
    document,
    items,
    loading,
    loadFailed,
    onClose,
    onCreate,
    onLoadMore,
    onCopyLink,
    onEdit,
    onRemove,
}) => {
    const listRef = useRef<HTMLDivElement | null>(null);
    const loadMoreRef = useRef<HTMLDivElement | null>(null);
    const documentId = document?.id;

    useEffect(() => {
        const root = listRef.current;
        const target = loadMoreRef.current;
        if (!open || !documentId || !root || !target || typeof IntersectionObserver === "undefined") {
            return;
        }
        const observer = new IntersectionObserver(([entry]) => {
            if (entry?.isIntersecting) {
                void onLoadMore(documentId);
            }
        }, {
            root,
            rootMargin: "0px 0px 120px",
        });
        observer.observe(target);
        return () => observer.disconnect();
    }, [documentId, items.length, onLoadMore, open]);

    const getMenuItems = (record: MindShareRecord): MenuProps["items"] => [
        {
            key: "copy",
            icon: <LinkOutlined/>,
            label: "复制链接",
            onClick: () => void onCopyLink(record),
        },
        {
            key: "edit",
            icon: <EditOutlined/>,
            label: "编辑分享",
            onClick: () => onEdit(record),
        },
        {type: "divider" as const},
        {
            key: "delete",
            icon: <DeleteOutlined/>,
            label: "删除分享",
            danger: true,
            onClick: () => onRemove(record),
        },
    ];

    return (
        <Drawer
            title={(
                <div className="share-title">
                    <span className="share-title-text">分享管理</span>
                    {document?.title && (
                        <Tooltip title={document.title}>
                            <span className="share-title-document">{document.title}</span>
                        </Tooltip>
                    )}
                </div>
            )}
            extra={(
                <Button
                    className="share-create-button"
                    type="text"
                    size="small"
                    icon={<PlusOutlined/>}
                    disabled={!document}
                    onClick={onCreate}
                >
                    新建分享
                </Button>
            )}
            rootClassName={styles.shareDrawer}
            size={420}
            open={open}
            onClose={onClose}
            destroyOnHidden
        >
            <div className={styles.sharePanel}>
                <div
                    ref={listRef}
                    className="share-list"
                    onScroll={(event) => {
                        const element = event.currentTarget;
                        if (element.scrollHeight - element.scrollTop - element.clientHeight <= 96 && documentId) {
                            void onLoadMore(documentId);
                        }
                    }}
                >
                    {loading && items.length <= 0 ? (
                        <div className="share-loading"><Spin size="large"/></div>
                    ) : items.length <= 0 ? (
                        <div className="share-empty">
                            <Empty image={Empty.PRESENTED_IMAGE_SIMPLE}>
                                {loadFailed && documentId && (
                                    <Button size="small" onClick={() => void onLoadMore(documentId)}>
                                        重新加载
                                    </Button>
                                )}
                            </Empty>
                        </div>
                    ) : (
                        <>
                            <div className="share-items">
                                {items.map((item) => {
                                    const expireTimestamp = item.expireTime
                                        ? new Date(item.expireTime.replace(" ", "T")).getTime()
                                        : Number.POSITIVE_INFINITY;
                                    const expired = Number.isFinite(expireTimestamp) && expireTimestamp <= Date.now();
                                    const enabled = item.status === SHARE_ENABLE && !expired;
                                    const state = expired ? "expired" : enabled ? "enabled" : "disabled";
                                    const stateText = expired ? "已过期" : enabled ? "已开启" : "已关闭";
                                    return (
                                        <article className="share-item" data-state={state} key={item.id}>
                                            <div className="share-item-head">
                                                <div className="share-code-line">
                                                    <span className="share-code">{item.code}</span>
                                                    <span className="share-state"><i/>{stateText}</span>
                                                </div>
                                                <Dropdown
                                                    trigger={["click"]}
                                                    placement="bottomRight"
                                                    arrow={{pointAtCenter: true}}
                                                    classNames={{root: styles.actionDropdown}}
                                                    menu={{items: getMenuItems(item)}}
                                                >
                                                    <Button
                                                        className="share-action-button"
                                                        type="text"
                                                        size="small"
                                                        icon={<MoreOutlined/>}
                                                        aria-label={`管理分享 ${item.code}`}
                                                    />
                                                </Dropdown>
                                            </div>
                                            <div className={`share-remark${item.remark ? "" : " empty"}`}>
                                                {item.remark || "未填写备注"}
                                            </div>
                                            <div className="share-meta">
                                                <Tooltip title={item.expireTime || "永久有效"}>
                                                    <span>{item.expireTime ? `有效至 ${item.expireTime}` : "永久有效"}</span>
                                                </Tooltip>
                                                <i/>
                                                <Tooltip title={item.createTime}>
                                                    <span>创建于 {getTimeText(item.createTime)}</span>
                                                </Tooltip>
                                            </div>
                                        </article>
                                    );
                                })}
                            </div>
                            <div ref={loadMoreRef} className="share-sentinel" aria-hidden="true"/>
                            {loading && (
                                <div className="share-load-more"><Spin size="small"/></div>
                            )}
                            {loadFailed && documentId && (
                                <div className="share-load-more">
                                    <Button size="small" onClick={() => void onLoadMore(documentId)}>
                                        重新加载
                                    </Button>
                                </div>
                            )}
                        </>
                    )}
                </div>
            </div>
        </Drawer>
    );
};

export default Share;
