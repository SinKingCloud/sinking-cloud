import React, {forwardRef, useEffect, useImperativeHandle, useRef, useState} from "react";
import {App, Button, Form, Progress, Select, Typography} from "antd";
import {
    CheckCircleOutlined,
    CloudUploadOutlined,
    FileOutlined,
    PauseCircleOutlined,
    ReloadOutlined,
    UploadOutlined,
} from "@ant-design/icons";
import {ProModal, ProModalRef, Title} from "sinking-antd";
import {createStyles} from "antd-style";
import {
    upload as uploadManager,
    type UploadProgressInfo,
    type UploadProgressPhase,
    type UploadResult,
    type UploadStartInfo,
    type UploadTaskHandle,
} from "@/pages/components/upload";
import {formatSize} from "@/utils/string";

type ProgressStatus = "normal" | "active" | "exception" | "success";

type UploadProgressState = {
    status: ProgressStatus;
    phase?: UploadProgressPhase | "ready";
    mode?: UploadStartInfo["mode"];
    percent: number;
    loaded: number;
    total: number;
    fileName?: string;
    fileSize?: number;
    chunkIndex?: number;
    chunkCount?: number;
    url?: string;
    error?: string;
};

export type FileUploadModalRef = {
    open: () => void;
};

type FileUploadModalProps = {
    sceneOptions: any[];
    onUploaded?: () => void;
};

const phaseText: Record<string, string> = {
    ready: "等待上传",
    hash: "计算哈希",
    init: "初始化",
    status: "检查断点",
    upload: "上传中",
    finish: "合并文件",
    complete: "上传完成",
    pause: "已暂停",
};

const modeText: Record<UploadStartInfo["mode"], string> = {
    simple: "普通上传",
    chunk: "分片上传",
};

const getUploadPercent = (info: UploadProgressInfo) => {
    if (info.phase === "complete") return 100;
    if (info.mode === "simple") return info.percent;
    if (info.phase === "hash") return Number(Math.min(8, info.percent * 0.08).toFixed(2));
    if (info.phase === "init" || info.phase === "status") return 10;
    if (info.phase === "upload") return Number(Math.min(96, 10 + info.percent * 0.86).toFixed(2));
    if (info.phase === "finish") return 98;
    return info.percent;
};

const useStyles: any = createStyles(({css, token}: any) => ({
    form: css`
        .ant-form-item:last-child {
            margin-bottom: 0;
        }
    `,
    input: css`
        position: fixed;
        top: -9999px;
        left: -9999px;
        width: 1px;
        height: 1px;
        opacity: 0;
        pointer-events: none;
    `,
    panel: css`
        padding: ${token.paddingLG}px;
        border: 1px dashed ${token.colorBorder};
        border-radius: ${token.borderRadiusLG}px;
        background: ${token.colorFillAlter};
    `,
    panelHead: css`
        display: flex;
        align-items: center;
        gap: ${token.margin}px;
    `,
    iconBox: css`
        width: 48px;
        height: 48px;
        border-radius: ${token.borderRadiusLG}px;
        background: ${token.colorPrimaryBg};
        color: ${token.colorPrimary};
        display: flex;
        align-items: center;
        justify-content: center;
        flex: none;
    `,
    cloudIcon: css`
        font-size: 24px;
        line-height: 1;
    `,
    main: css`
        min-width: 0;
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 2px;
    `,
    title: css`
        color: ${token.colorText};
        font-size: ${token.fontSizeLG}px;
        line-height: 24px;
    `,
    desc: css`
        color: ${token.colorTextTertiary};
        font-size: ${token.fontSizeSM}px;
        line-height: 20px;
    `,
    file: css`
        margin-top: ${token.margin}px;
        padding: ${token.paddingXS}px ${token.paddingSM}px;
        border-radius: ${token.borderRadius}px;
        background: ${token.colorBgContainer};
        display: flex;
        align-items: center;
        gap: ${token.marginSM}px;
    `,
    fileIcon: css`
        color: ${token.colorTextTertiary};
        flex: none;
    `,
    fileName: css`
        max-width: 280px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    `,
    progress: css`
        margin-top: ${token.margin}px;
        padding: ${token.paddingSM}px ${token.padding}px;
        border: 1px solid ${token.colorBorderSecondary};
        border-radius: ${token.borderRadius}px;
        background: ${token.colorFillTertiary};
    `,
    meta: css`
        display: grid;
        grid-template-columns: repeat(3, minmax(0, 1fr));
        gap: ${token.marginSM}px;
        margin-bottom: ${token.marginSM}px;
    `,
    metaItem: css`
        min-width: 0;
        display: flex;
        flex-direction: column;
        gap: 2px;
    `,
    metaLabel: css`
        color: ${token.colorTextTertiary};
        font-size: ${token.fontSizeSM}px;
        line-height: 20px;
    `,
    metaValue: css`
        color: ${token.colorText};
        font-size: ${token.fontSize}px;
        line-height: 22px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    `,
    progressText: css`
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: ${token.marginSM}px;
        margin-bottom: ${token.marginXXS}px;
    `,
    result: css`
        display: block;
        margin-top: ${token.marginXS}px;
        word-break: break-all;
    `,
    footer: css`
        display: flex;
        justify-content: flex-end;
    `,
}));

const toProgress = (info: UploadProgressInfo): UploadProgressState => ({
    status: "active",
    phase: info.phase,
    mode: info.mode,
    percent: getUploadPercent(info),
    loaded: info.loaded,
    total: info.total,
    fileName: info.file?.name,
    fileSize: info.file?.size,
    chunkIndex: info.chunkIndex,
    chunkCount: info.chunkCount,
});

const toCompleteProgress = (result: UploadResult): UploadProgressState => ({
    status: "success",
    phase: "complete",
    mode: result.mode,
    percent: 100,
    loaded: result.total,
    total: result.total,
    fileName: result.file?.name,
    fileSize: result.file?.size,
    chunkCount: result.chunkCount,
    url: result.url,
});

const FileUploadModal = forwardRef<FileUploadModalRef, FileUploadModalProps>(({sceneOptions, onUploaded}, ref) => {
    const {styles} = useStyles();
    const {message} = App.useApp();
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const inputRef = useRef<HTMLInputElement | null>(null);
    const fileRef = useRef<File | null>(null);
    const taskRef = useRef<UploadTaskHandle | undefined>(undefined);
    const frameRef = useRef<number | undefined>(undefined);
    const progressRef = useRef<UploadProgressInfo | undefined>(undefined);
    const [form] = Form.useForm();
    const [loading, setLoading] = useState(false);
    const [paused, setPaused] = useState(false);
    const [progress, setProgress] = useState<UploadProgressState>();

    const clearFrame = () => {
        if (frameRef.current) {
            cancelAnimationFrame(frameRef.current);
            frameRef.current = undefined;
        }
    };

    const reset = () => {
        clearFrame();
        form.resetFields();
        setLoading(false);
        setPaused(false);
        setProgress(undefined);
        fileRef.current = null;
        taskRef.current = undefined;
        progressRef.current = undefined;
        if (inputRef.current) {
            inputRef.current.value = "";
        }
    };

    useEffect(() => clearFrame, []);

    useImperativeHandle(ref, () => ({
        open: () => {
            reset();
            modalRef.current?.show();
        },
    }));

    const flushProgress = () => {
        frameRef.current = undefined;
        if (progressRef.current) {
            setProgress(toProgress(progressRef.current));
        }
    };

    const onProgress = (info: UploadProgressInfo) => {
        progressRef.current = info;
        if (!frameRef.current) {
            frameRef.current = requestAnimationFrame(flushProgress);
        }
    };

    const onComplete = (result: UploadResult) => {
        clearFrame();
        progressRef.current = undefined;
        setProgress(toCompleteProgress(result));
        onUploaded?.();
        message?.success("上传成功");
    };

    const onError = (error: Error) => {
        clearFrame();
        progressRef.current = undefined;
        if (error?.message === "上传已取消") {
            setPaused(true);
            setProgress((prev) => prev ? {...prev, status: "normal", phase: "pause"} : prev);
            return;
        }
        setProgress((prev) => ({
            status: "exception",
            phase: prev?.phase || "upload",
            mode: prev?.mode,
            percent: prev?.percent || 0,
            loaded: prev?.loaded || 0,
            total: prev?.total || 0,
            fileName: prev?.fileName,
            fileSize: prev?.fileSize,
            chunkIndex: prev?.chunkIndex,
            chunkCount: prev?.chunkCount,
            error: error?.message || "上传失败",
        }));
        message?.error(error?.message || "上传失败");
    };

    const startUpload = (file?: File | null) => {
        if (!file) return;
        const sceneId = form.getFieldValue("scene_id" as any);
        if (sceneId === undefined || sceneId === null || sceneId === "") {
            message?.warning("请选择上传场景");
            return;
        }

        const lastTask = taskRef.current?.snapshot?.();
        const cachedHash = lastTask?.fileName === file.name && lastTask?.fileSize === file.size ? lastTask?.hash : undefined;

        setLoading(true);
        setPaused(false);
        setProgress(undefined);
        progressRef.current = undefined;

        const task = uploadManager.upload({
            sceneId,
            file,
            hash: cachedHash,
            onStart: (info: UploadStartInfo) => setProgress({
                status: "active",
                phase: "ready",
                mode: info.mode,
                percent: 0,
                loaded: 0,
                total: info.total,
                fileName: info.file?.name,
                fileSize: info.file?.size,
                chunkCount: info.chunkCount,
            }),
            onProgress,
            onComplete,
            onError,
        } as any);

        taskRef.current = task;
        task.promise.catch(() => undefined).finally(() => setLoading(false));
    };

    const chooseFile = () => {
        if (loading) return;
        form.validateFields(["scene_id"] as any).then(() => {
            if (paused && fileRef.current) {
                startUpload(fileRef.current);
                return;
            }
            if (inputRef.current) {
                inputRef.current.value = "";
                inputRef.current.click();
            }
        }).catch(() => undefined);
    };

    const pauseUpload = () => {
        const taskId = taskRef.current?.id;
        if (taskId) {
            uploadManager.abort(taskId);
        }
        setLoading(false);
        setPaused(true);
        setProgress((prev) => prev ? {...prev, status: "normal", phase: "pause"} : prev);
        message?.info("上传已暂停");
    };

    const uploadMeta = progress ? [
        ["上传方式", progress.mode ? modeText[progress.mode] : "自动识别"],
        ["当前阶段", phaseText[progress.phase || "ready"] || "-"],
        ["文件大小", formatSize(progress.total)],
    ] : [];

    return (
        <ProModal
            ref={modalRef}
            title={<Title>上传文件</Title>}
            width={400}
            modalProps={{
                forceRender: true,
                destroyOnHidden: false,
                style: {
                    top: 100,
                    paddingBottom: 100,
                },
                mask: {
                    closable: !loading,
                },
                afterClose: reset,
                footer: (
                    <div className={styles.footer}>
                        <Button disabled={loading} onClick={() => modalRef.current?.hide()}>关闭</Button>
                    </div>
                ),
            } as any}
        >
            <Form form={form} layout="vertical" className={styles.form}>
                <input
                    ref={inputRef}
                    className={styles.input}
                    type="file"
                    onChange={(event) => {
                        fileRef.current = event.target.files?.[0] || null;
                        startUpload(fileRef.current);
                    }}
                />

                <Form.Item name="scene_id" label="上传场景" rules={[{required: true, message: "请选择上传场景"}]}>
                    <Select placeholder="请选择上传场景" options={sceneOptions as any} disabled={loading}/>
                </Form.Item>

                <div className={styles.panel}>
                    <div className={styles.panelHead}>
                        <div className={styles.iconBox}>
                            <CloudUploadOutlined className={styles.cloudIcon}/>
                        </div>
                        <div className={styles.main}>
                            <Typography.Text className={styles.title}>上传文件</Typography.Text>
                            <Typography.Text className={styles.desc}>上传并获取外链</Typography.Text>
                        </div>
                        <Button
                            type="primary"
                            icon={paused ? <ReloadOutlined/> : <UploadOutlined/>}
                            loading={loading}
                            disabled={loading}
                            onClick={chooseFile}
                        >
                            {paused ? "继续上传" : progress?.status === "success" ? "重新上传" : "选择文件"}
                        </Button>
                    </div>

                    {progress && (
                        <div className={styles.file}>
                            {progress.status === "success" ? (
                                <CheckCircleOutlined className={styles.fileIcon}/>
                            ) : (
                                <FileOutlined className={styles.fileIcon}/>
                            )}
                            <div className={styles.main}>
                                <Typography.Text
                                    className={styles.fileName}>{progress.fileName || "-"}</Typography.Text>
                                <Typography.Text type="secondary">
                                    {formatSize(progress.fileSize || progress.total)}
                                </Typography.Text>
                            </div>
                        </div>
                    )}
                </div>

                {progress && (
                    <div className={styles.progress}>
                        <div className={styles.meta}>
                            {uploadMeta.map(([label, value]) => (
                                <div className={styles.metaItem} key={label}>
                                    <span className={styles.metaLabel}>{label}</span>
                                    <span className={styles.metaValue}>{value}</span>
                                </div>
                            ))}
                        </div>
                        <div className={styles.progressText}>
                            <Typography.Text type="secondary">
                                {formatSize(progress.loaded)} / {formatSize(progress.total)}
                            </Typography.Text>
                            <Typography.Text type="secondary">
                                {progress.chunkIndex && progress.chunkCount ? `分片 ${progress.chunkIndex}/${progress.chunkCount}` : " "}
                            </Typography.Text>
                        </div>
                        <Progress
                            percent={progress.percent}
                            status={progress.status === "active" ? "active" : progress.status}
                            format={(value) => `${Math.floor(Number(value) || 0)}%`}
                        />
                        {progress.error && (
                            <Typography.Text type="danger" className={styles.result}>{progress.error}</Typography.Text>
                        )}
                        {progress.url && (
                            <Typography.Text copyable className={styles.result}>{progress.url}</Typography.Text>
                        )}
                    </div>
                )}

                {loading && (
                    <Button block icon={<PauseCircleOutlined/>} onClick={pauseUpload}>
                        暂停上传
                    </Button>
                )}
            </Form>
        </ProModal>
    );
});

export default FileUploadModal;
