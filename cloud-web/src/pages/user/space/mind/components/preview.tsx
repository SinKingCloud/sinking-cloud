import React, {useEffect, useRef} from "react";
import {Modal, Spin} from "antd";
import {Title} from "sinking-antd";
import Mind, {type MindRef} from "@/pages/components/mind";
import type {MindVersionRecord} from "../types";
import {cloneValue} from "../utils";

export type PreviewProps = {
    styles: Record<string, string>;
    open: boolean;
    loading: boolean;
    version?: MindVersionRecord;
    onRestore: (record: MindVersionRecord) => void;
    onClose: () => void;
};

const Preview: React.FC<PreviewProps> = ({styles, open, loading, version, onRestore, onClose}) => {
    const mindRef = useRef<MindRef | null>(null);
    const canvasRef = useRef<HTMLDivElement | null>(null);

    useEffect(() => {
        const element = canvasRef.current;
        if (!version || !element || typeof ResizeObserver === "undefined") {
            return;
        }
        let frame = 0;
        let innerFrame = 0;
        const fitView = () => {
            cancelAnimationFrame(frame);
            cancelAnimationFrame(innerFrame);
            frame = requestAnimationFrame(() => {
                innerFrame = requestAnimationFrame(() => mindRef.current?.fitView());
            });
        };
        const observer = new ResizeObserver(() => {
            fitView();
        });
        observer.observe(element);
        fitView();
        return () => {
            cancelAnimationFrame(frame);
            cancelAnimationFrame(innerFrame);
            observer.disconnect();
        };
    }, [version]);

    return (
        <Modal
            title={<Title>版本预览</Title>}
            rootClassName={styles.versionPreviewModal}
            open={open}
            width={900}
            centered
            okText="恢复"
            cancelText="返回"
            okButtonProps={{disabled: loading || !version}}
            onOk={() => {
                if (version) {
                    onRestore(version);
                }
            }}
            onCancel={onClose}
            afterOpenChange={(open) => {
                if (!open) {
                    return;
                }
                requestAnimationFrame(() => {
                    requestAnimationFrame(() => mindRef.current?.fitView());
                });
            }}
            destroyOnHidden
        >
            {loading ? (
                <div className={styles.versionPreviewLoading}>
                    <Spin size="large"/>
                </div>
            ) : version ? (
                <div className={styles.versionPreview}>
                    <div className="preview-meta">
                        <div className="preview-remark">{version.remark}</div>
                        <div className="preview-details">
                            <span>{version.createTime}</span>
                            <span>{version.nodeCount} 个节点</span>
                            <span>{version.depth} 层</span>
                        </div>
                    </div>
                    <div ref={canvasRef} className="preview-canvas">
                        <Mind
                            key={version.id}
                            ref={mindRef}
                            defaultValue={cloneValue(version.value!)}
                            height="100%"
                            readOnly
                            showControls
                            showMiniMap={false}
                            shortcuts={false}
                        />
                    </div>
                </div>
            ) : null}
        </Modal>
    );
};

export default Preview;
