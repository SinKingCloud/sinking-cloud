import React from "react";
import {Button, Drawer, Empty, Segmented, Spin, Tooltip} from "antd";
import {EyeOutlined} from "@ant-design/icons";
import {VERSION_AUTO, VERSION_MANUAL} from "../constants";
import type {MindDocument, MindVersionRecord, MindVersionSaveType} from "../types";
import {getTimeText} from "../utils";

export type VersionProps = {
    styles: Record<string, string>;
    open: boolean;
    document?: MindDocument;
    filter: MindVersionSaveType;
    items: MindVersionRecord[];
    loading: boolean;
    loadFailed: boolean;
    onClose: () => void;
    onLoadMore: (documentId: string, filter: MindVersionSaveType) => void | Promise<void>;
    onFilterChange: (filter: MindVersionSaveType) => void;
    onPreview: (record: MindVersionRecord) => void | Promise<void>;
};

const Version: React.FC<VersionProps> = ({
    styles,
    open,
    document,
    filter,
    items,
    loading,
    loadFailed,
    onClose,
    onLoadMore,
    onFilterChange,
    onPreview,
}) => (
    <Drawer
        title="版本记录"
        rootClassName={styles.versionDrawer}
        open={open}
        onClose={onClose}
        destroyOnHidden
    >
        <div
            className={styles.versionList}
            onScroll={(event) => {
                const element = event.currentTarget;
                if (element.scrollHeight - element.scrollTop - element.clientHeight <= 96 && document) {
                    void onLoadMore(document.id, filter);
                }
            }}
        >
            <div className="version-filter">
                <Segmented
                    value={filter}
                    options={[
                        {label: "手动保存", value: VERSION_MANUAL},
                        {label: "自动保存", value: VERSION_AUTO},
                    ]}
                    onChange={(value) => onFilterChange(value as MindVersionSaveType)}
                />
            </div>

            {loading && items.length <= 0 ? (
                <div className="empty"><Spin/></div>
            ) : items.length <= 0 ? (
                <div className="empty">
                    <Empty image={Empty.PRESENTED_IMAGE_SIMPLE}>
                        {loadFailed && document && (
                            <Button size="small" onClick={() => void onLoadMore(document.id, filter)}>
                                重新加载
                            </Button>
                        )}
                    </Empty>
                </div>
            ) : (
                <ol className="timeline" aria-label="历史版本">
                    {items.map((item) => {
                        const isCurrent = item?.current;
                        return (
                            <li
                                className={`version${isCurrent ? " current" : ""}`}
                                key={item?.id}
                                aria-current={isCurrent ? "true" : undefined}
                            >
                                <div className="axis" aria-hidden="true">
                                    <span className="dot"/>
                                    <span className="line"/>
                                </div>
                                <div className="body" onClick={() => void onPreview(item)}>
                                    <div className="head">
                                        <div className="meta">
                                            <span className="snapshot-version">
                                                {item?.saveType === VERSION_AUTO ? "自动保存" : "手动保存"}
                                            </span>
                                            {isCurrent && <span className="current-label">当前版本</span>}
                                        </div>
                                        <Tooltip title={item?.createTime}>
                                            <span className="snapshot-time">{getTimeText(item?.createTime)}</span>
                                        </Tooltip>
                                    </div>
                                    <div className="snapshot-text">{item?.remark}</div>
                                    <div className="snapshot-extra">
                                        <span/>
                                        <Button
                                            className="preview-button"
                                            type="link"
                                            size="small"
                                            icon={<EyeOutlined/>}
                                            aria-label={`预览 ${item?.createTime} 的版本`}
                                            onClick={(event) => {
                                                event.stopPropagation();
                                                void onPreview(item);
                                            }}
                                        >
                                            预览
                                        </Button>
                                    </div>
                                </div>
                            </li>
                        );
                    })}
                </ol>
            )}
            {loading && items.length > 0 && (
                <div style={{display: "flex", justifyContent: "center", padding: 12}}><Spin/></div>
            )}
            {loadFailed && items.length > 0 && document && (
                <div style={{display: "flex", justifyContent: "center", padding: 12}}>
                    <Button size="small" onClick={() => void onLoadMore(document.id, filter)}>
                        重新加载
                    </Button>
                </div>
            )}
        </div>
    </Drawer>
);

export default Version;
