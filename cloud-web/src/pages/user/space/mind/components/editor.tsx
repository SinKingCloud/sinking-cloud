import React from "react";
import {App, Drawer, Dropdown, Spin, Tooltip} from "antd";
import {
    ApartmentOutlined,
    CloseOutlined,
    CopyOutlined,
    DownloadOutlined,
    FileImageOutlined,
    FileMarkdownOutlined,
    FileTextOutlined,
    FullscreenExitOutlined,
    FullscreenOutlined,
    HistoryOutlined,
    SaveOutlined,
    ShareAltOutlined,
} from "@ant-design/icons";
import Mind, {type MindProps, type MindRef, type MindValue} from "@/pages/components/mind";
import type {MindDocument} from "../types";

export type EditorProps = {
    styles: Record<string, string>;
    open: boolean;
    loading: boolean;
    document?: MindDocument;
    editorRef: React.RefObject<MindRef | null>;
    value: MindValue;
    revision: number;
    closing: boolean;
    onExportPng: () => void | Promise<void>;
    onExportJpg: () => void | Promise<void>;
    onCopyImage: () => void | Promise<void>;
    onClose: () => void | Promise<void>;
    onSave: () => void;
    onArrange: () => void;
    onOpenVersionHistory: () => void;
    onOpenShare: (document: MindDocument) => void;
    onDataChange: (value: MindValue) => void;
    onDirty: () => void;
    onPersist: NonNullable<MindProps["onSave"]>;
};

const Editor: React.FC<EditorProps> = ({
    styles,
    open,
    loading,
    document,
    editorRef,
    value,
    revision,
    closing,
    onExportPng,
    onExportJpg,
    onCopyImage,
    onClose,
    onSave,
    onArrange,
    onOpenVersionHistory,
    onOpenShare,
    onDataChange,
    onDirty,
    onPersist,
}) => {
    const {message} = App.useApp();
    const [fullscreen, setFullscreen] = React.useState(false);

    React.useEffect(() => {
        const fullscreenDocument = window.document as Document & {
            webkitFullscreenElement?: Element;
        };
        const handleFullscreenChange = () => {
            setFullscreen(
                fullscreenDocument.fullscreenElement === fullscreenDocument.documentElement ||
                fullscreenDocument.webkitFullscreenElement === fullscreenDocument.documentElement,
            );
            requestAnimationFrame(() => editorRef.current?.fitView());
        };
        fullscreenDocument.addEventListener("fullscreenchange", handleFullscreenChange);
        fullscreenDocument.addEventListener("webkitfullscreenchange", handleFullscreenChange);
        return () => {
            fullscreenDocument.removeEventListener("fullscreenchange", handleFullscreenChange);
            fullscreenDocument.removeEventListener("webkitfullscreenchange", handleFullscreenChange);
        };
    }, [editorRef]);

    const toggleFullscreen = async () => {
        const fullscreenDocument = window.document as Document & {
            webkitExitFullscreen?: () => Promise<void> | void;
            webkitFullscreenElement?: Element;
        };
        const element = fullscreenDocument.documentElement as HTMLElement & {
            webkitRequestFullscreen?: () => Promise<void> | void;
        };
        const active = fullscreenDocument.fullscreenElement === element ||
            fullscreenDocument.webkitFullscreenElement === element;
        try {
            if (active) {
                if (fullscreenDocument.exitFullscreen) {
                    await fullscreenDocument.exitFullscreen();
                } else {
                    await fullscreenDocument.webkitExitFullscreen?.();
                }
                return;
            }
            if (element.requestFullscreen) {
                await element.requestFullscreen();
            } else if (element.webkitRequestFullscreen) {
                await element.webkitRequestFullscreen();
            } else {
                message.warning("当前浏览器不支持全屏");
            }
        } catch {
            message.error("全屏切换失败");
        }
    };

    const closeEditor = async () => {
        const fullscreenDocument = window.document as Document & {
            webkitExitFullscreen?: () => Promise<void> | void;
            webkitFullscreenElement?: Element;
        };
        const active = fullscreenDocument.fullscreenElement === fullscreenDocument.documentElement ||
            fullscreenDocument.webkitFullscreenElement === fullscreenDocument.documentElement;
        if (active) {
            try {
                if (fullscreenDocument.exitFullscreen) {
                    await fullscreenDocument.exitFullscreen();
                } else {
                    await fullscreenDocument.webkitExitFullscreen?.();
                }
            } catch {
                // 退出全屏失败不阻止关闭编辑器。
            }
        }
        await onClose();
    };

    React.useEffect(() => {
        if (!open || loading || !document) {
            return;
        }
        let frame = 0;
        let innerFrame = 0;
        frame = requestAnimationFrame(() => {
            innerFrame = requestAnimationFrame(() => editorRef.current?.fitView());
        });
        return () => {
            cancelAnimationFrame(frame);
            cancelAnimationFrame(innerFrame);
        };
    }, [document?.id, editorRef, loading, open, revision]);

    const exportItems = [
        {key: "png", icon: <FileImageOutlined/>, label: "导出 PNG", onClick: () => void onExportPng()},
        {key: "jpg", icon: <FileImageOutlined/>, label: "导出 JPG", onClick: () => void onExportJpg()},
        {key: "json", icon: <FileTextOutlined/>, label: "导出 JSON", onClick: () => editorRef.current?.downloadJson()},
        {
            key: "markdown",
            icon: <FileMarkdownOutlined/>,
            label: "导出 Markdown",
            onClick: () => editorRef.current?.downloadMarkdown(),
        },
        {type: "divider" as const},
        {key: "copy", icon: <CopyOutlined/>, label: "复制图片", onClick: () => void onCopyImage()},
    ];

    return (
        <Drawer
            open={open}
            rootClassName={styles.editorDrawer}
            placement="right"
            size="100%"
            closable={false}
            mask={false}
            destroyOnHidden
            onClose={() => void closeEditor()}
            afterOpenChange={(nextOpen) => {
                if (!nextOpen) {
                    return;
                }
                requestAnimationFrame(() => {
                    requestAnimationFrame(() => editorRef.current?.fitView());
                });
            }}
        >
            <div className={`${styles.editor} ${styles.editorFullscreen}`} aria-busy={closing || loading}>
                {loading ? (
                    <div className={styles.editorLoading}>
                        <Spin size="large"/>
                    </div>
                ) : document ? (
                    <>
                        <div className={styles.editorToolbar} style={closing ? {pointerEvents: "none"} : undefined}>
                            <div className="actions">
                                <Tooltip title="保存">
                                    <button className={styles.editorToolButton} type="button" onClick={onSave}>
                                        <SaveOutlined/>
                                    </button>
                                </Tooltip>
                                <Tooltip title="整理布局">
                                    <button className={styles.editorToolButton} type="button" onClick={onArrange}>
                                        <ApartmentOutlined/>
                                    </button>
                                </Tooltip>
                                <Tooltip title={fullscreen ? "退出全屏" : "全屏"}>
                                    <button
                                        className={styles.editorToolButton}
                                        type="button"
                                        aria-label={fullscreen ? "退出全屏" : "全屏"}
                                        onClick={() => void toggleFullscreen()}
                                    >
                                        {fullscreen ? <FullscreenExitOutlined/> : <FullscreenOutlined/>}
                                    </button>
                                </Tooltip>
                                <Dropdown trigger={["click"]} placement="bottomRight" menu={{items: exportItems as any}}>
                                    <button className={styles.editorToolButton} type="button" title="导出">
                                        <DownloadOutlined/>
                                    </button>
                                </Dropdown>
                                <Tooltip title="版本记录">
                                    <button
                                        className={styles.editorToolButton}
                                        type="button"
                                        aria-label="版本记录"
                                        onClick={onOpenVersionHistory}
                                    >
                                        <HistoryOutlined/>
                                    </button>
                                </Tooltip>
                                <Tooltip title="分享管理">
                                    <button
                                        className={styles.editorToolButton}
                                        type="button"
                                        onClick={() => onOpenShare(document)}
                                    >
                                        <ShareAltOutlined/>
                                    </button>
                                </Tooltip>
                            </div>
                        </div>
                        <div className={styles.canvas} style={closing ? {pointerEvents: "none"} : undefined}>
                            <Mind
                                key={`${document.id}-${revision}`}
                                ref={editorRef}
                                defaultValue={value}
                                height="100%"
                                autoSave
                                autoSaveDelay={5000}
                                onDataChange={onDataChange}
                                onDirty={onDirty}
                                onSave={onPersist}
                                readOnly={closing}
                            />
                        </div>
                    </>
                ) : null}
                <button
                    className={styles.editorClose}
                    type="button"
                    aria-label="关闭导图"
                    disabled={closing}
                    onClick={() => void closeEditor()}
                >
                    <CloseOutlined/>
                </button>
            </div>
        </Drawer>
    );
};

export default Editor;
