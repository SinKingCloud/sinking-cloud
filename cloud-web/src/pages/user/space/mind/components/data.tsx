import React from "react";
import type {MenuProps} from "antd";
import {Button, Card, Col, Dropdown, Empty, Pagination, Row, Spin, Tooltip,} from "antd";
import {
    CopyOutlined,
    DeleteOutlined,
    EditOutlined,
    FileTextOutlined,
    MoreOutlined,
    ShareAltOutlined,
} from "@ant-design/icons";
import type {MindDocument} from "../types";
import {getStableColor, getTimeText} from "../utils";

const MindPreviewSvg: React.FC = () => (
    <svg className="mind-preview" viewBox="0 0 320 170" aria-hidden="true" focusable="false">
        <rect width="320" height="170" className="preview-surface"/>
        <path
            className="grid-line"
            d="M40 0V170M80 0V170M120 0V170M160 0V170M200 0V170M240 0V170M280 0V170M0 34H320M0 68H320M0 102H320M0 136H320"
        />
        <path className="preview-link" d="M159 85C188 85 184 47 216 47"/>
        <path className="preview-link" d="M159 85C190 85 184 123 218 123"/>
        <path className="preview-link" d="M161 85C132 85 136 47 104 47"/>
        <path className="preview-link" d="M161 85C130 85 136 123 102 123"/>
        <rect className="preview-root" x="126" y="70" width="68" height="30" rx="15"/>
        <rect className="preview-node" x="216" y="35" width="62" height="24" rx="12"/>
        <rect className="preview-node" x="218" y="111" width="62" height="24" rx="12"/>
        <rect className="preview-node" x="42" y="35" width="62" height="24" rx="12"/>
        <rect className="preview-node" x="40" y="111" width="62" height="24" rx="12"/>
        <circle className="preview-dot" cx="160" cy="85" r="4"/>
    </svg>
);

type DataProps = {
    styles: Record<string, string>;
    documents: MindDocument[];
    loading: boolean;
    total: number;
    page: number;
    pageSize: number;
    openingDocumentId?: string;
    onOpen: (document: MindDocument) => void | Promise<void>;
    onEdit: (document: MindDocument) => void;
    onCopy: (document: MindDocument) => void | Promise<void>;
    onShare: (document: MindDocument) => void;
    onDelete: (document: MindDocument) => void;
    onPaginationChange: (page: number, pageSize: number) => void;
};

const Data: React.FC<DataProps> = ({
    styles,
    documents,
    loading,
    total,
    page,
    pageSize,
    openingDocumentId,
    onOpen,
    onEdit,
    onCopy,
    onShare,
    onDelete,
    onPaginationChange,
}) => {
    const getMenuItems = (document: MindDocument): MenuProps["items"] => [
        {
            key: "open",
            icon: <EditOutlined/>,
            label: "打开导图",
            onClick: ({domEvent}: any) => {
                domEvent.stopPropagation();
                void onOpen(document);
            },
        },
        {
            key: "info",
            icon: <FileTextOutlined/>,
            label: "编辑信息",
            onClick: ({domEvent}: any) => {
                domEvent.stopPropagation();
                onEdit(document);
            },
        },
        {
            key: "copy",
            icon: <CopyOutlined/>,
            label: "复制导图",
            onClick: ({domEvent}: any) => {
                domEvent.stopPropagation();
                void onCopy(document);
            },
        },
        {
            key: "share",
            icon: <ShareAltOutlined/>,
            label: "分享管理",
            onClick: ({domEvent}: any) => {
                domEvent.stopPropagation();
                onShare(document);
            },
        },
        {type: "divider" as const},
        {
            key: "delete",
            icon: <DeleteOutlined/>,
            label: "删除导图",
            danger: true,
            onClick: ({domEvent}: any) => {
                domEvent.stopPropagation();
                onDelete(document);
            },
        },
    ];

    return (
        <Row>
        <Col className={styles.documentDataColumn} span={24}>
            <div className={`${styles.documentDataCard} ${loading || documents.length <= 0 ? "empty" : ""}`}>
                {loading ? (
                    <Spin size="large"/>
                ) : documents.length > 0 ? (
                    <div className={styles.documentGrid}>
                        {documents.map((item) => (
                            <Card
                                key={item?.id}
                                className={`${styles.documentCard}${openingDocumentId === item?.id ? " loading" : ""}`}
                                style={{"--accent": getStableColor(item?.id || item?.title)} as React.CSSProperties}
                                styles={{body: {padding: 0}}}
                                onClick={() => void onOpen(item)}
                            >
                                <div className="preview">
                                    <MindPreviewSvg/>
                                </div>
                                <div className="body">
                                    <div className="head">
                                        <button
                                            className="title"
                                            type="button"
                                            onClick={(event) => {
                                                event.stopPropagation();
                                                void onOpen(item);
                                            }}
                                        >
                                            {item?.title}
                                        </button>
                                        <span className="actions" onClick={(event) => event.stopPropagation()}>
                                            <Dropdown
                                                trigger={["click"]}
                                                placement="bottomRight"
                                                arrow={{pointAtCenter: true}}
                                                classNames={{root: styles.actionDropdown}}
                                                menu={{items: getMenuItems(item)}}
                                            >
                                                <Button
                                                    type="text"
                                                    size="small"
                                                    icon={<MoreOutlined/>}
                                                    aria-label={`${item?.title}导图操作`}
                                                />
                                            </Dropdown>
                                        </span>
                                    </div>
                                    <div className="desc">{item?.description || "暂无描述"}</div>
                                    <div className="meta">
                                        <span>{item?.shared ? "已分享" : "仅自己可见"}</span>
                                        <Tooltip title={item?.updateTime}>
                                            <span>{getTimeText(item?.updateTime)}</span>
                                        </Tooltip>
                                    </div>
                                </div>
                            </Card>
                        ))}
                    </div>
                ) : <Empty image={Empty.PRESENTED_IMAGE_SIMPLE}/>}
            </div>
        </Col>
        {total > 0 && (
            <Col className={styles.documentPaginationColumn} span={24}>
                <div className={styles.documentPaginationCard}>
                    <Pagination
                        current={page}
                        pageSize={pageSize}
                        total={total}
                        align="center"
                        size="small"
                        responsive
                        showQuickJumper
                        showLessItems
                        showSizeChanger={{
                            showSearch: false,
                            variant: "filled",
                            size: "small",
                            className: styles.pageSizeSelect,
                            classNames: {
                                popup: {
                                    root: styles.pageSizeDropdown,
                                },
                            },
                        } as any}
                        pageSizeOptions={[10, 20, 50, 100]}
                        showTotal={(nextTotal, range) => `第 ${range[0]}-${range[1]} 条 / 共 ${nextTotal} 条`}
                        onChange={(nextPage, nextPageSize) => {
                            const sizeChanged = nextPageSize !== pageSize;
                            onPaginationChange(sizeChanged ? 1 : nextPage, nextPageSize);
                        }}
                    />
                </div>
            </Col>
        )}
        </Row>
    );
};

export default Data;
