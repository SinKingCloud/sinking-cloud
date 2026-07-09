import React, {useMemo, useRef, useState} from "react";
import {App, Button, Card, Space, Tooltip, Typography} from "antd";
import {
    ApartmentOutlined,
    CopyOutlined,
    DeleteOutlined,
    EyeInvisibleOutlined,
    EyeOutlined,
    FileImageOutlined,
    ForkOutlined,
    FullscreenExitOutlined,
    FullscreenOutlined,
    PlusOutlined,
    ReloadOutlined,
    SaveOutlined,
    ThunderboltOutlined,
    UndoOutlined,
    ZoomInOutlined,
} from "@ant-design/icons";
import {createStyles} from "antd-style";
import Mind, {createMindBenchmarkValue, getMindStats, type MindRef, type MindValue,} from "@/pages/components/mind";
import {Body, Title} from "sinking-antd";

const useStyles = createStyles(({css, token}: any) => ({
    page: css`
        display: grid;
        grid-template-columns: 260px minmax(0, 1fr);
        gap: ${token.margin}px;

        @media (max-width: 900px) {
            grid-template-columns: 1fr;
        }
    `,
    sideCard: css`
        position: sticky;
        top: ${token.margin}px;
        align-self: start;

        @media (max-width: 900px) {
            position: static;
        }
    `,
    statGrid: css`
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: 8px;
        margin-bottom: ${token.margin}px;
    `,
    statItem: css`
        padding: 10px 8px;
        border-radius: ${token.borderRadius}px;
        border: 1px solid ${token.colorBorderSecondary};
        background: ${token.colorBgContainer};
        text-align: center;

        strong {
            display: block;
            color: ${token.colorTextHeading};
            font-size: 16px;
            line-height: 22px;
            font-weight: 600;
        }

        span {
            color: ${token.colorTextTertiary};
            font-size: 12px;
            line-height: 18px;
        }
    `,
    canvasCard: css`
        min-width: 0;

        .ant-card-body {
            padding: 0;
            overflow: hidden;
        }
    `,
    toolbar: css`
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: ${token.marginSM}px;
        padding: ${token.paddingSM}px ${token.padding}px;
        border-bottom: 1px solid ${token.colorBorderSecondary};
        background: ${token.colorBgContainer};

        @media (max-width: 700px) {
            align-items: flex-start;
            flex-direction: column;
        }
    `,
    toolbarHint: css`
        max-width: 420px;
    `,
}));

const TestMindPage: React.FC = () => {
    const {styles} = useStyles();
    const {message} = App.useApp();
    const editorRef = useRef<MindRef>(null);
    const [value, setValue] = useState<MindValue>({nodes: []});
    const [controlsVisible, setControlsVisible] = useState(true);
    const [fullscreen, setFullscreen] = useState(false);
    const [miniMapVisible, setMiniMapVisible] = useState(true);
    const [lastSave, setLastSave] = useState("暂未保存");
    const stats = useMemo(() => getMindStats(value), [value]);

    const resetMind = () => {
        editorRef.current?.reset();
    };

    const copyJson = async () => {
        try {
            await navigator.clipboard?.writeText(editorRef.current?.exportJson() || JSON.stringify(value, null, 2));
            message.success("JSON 已复制");
        } catch (error) {
            message.error("复制失败");
        }
    };

    const downloadUrl = (url: string, filename: string) => {
        const link = document.createElement("a");
        link.href = url;
        link.download = filename;
        link.click();
    };

    const exportPng = async () => {
        try {
            const url = await editorRef.current?.exportPng();
            if (url) {
                downloadUrl(url, `mind-${Date.now()}.png`);
                message.success("图片已导出");
            }
        } catch (error) {
            message.error(error instanceof Error ? error.message : "导出失败");
        }
    };

    const exportJpg = async () => {
        try {
            const url = await editorRef.current?.exportJpg();
            if (url) {
                downloadUrl(url, `mind-${Date.now()}.jpg`);
                message.success("JPG 已导出");
            }
        } catch (error) {
            message.error(error instanceof Error ? error.message : "导出失败");
        }
    };

    const copyImage = async () => {
        try {
            await editorRef.current?.copyImage();
            message.success("图片已复制");
        } catch (error) {
            const messageText = error instanceof Error ? error.message : "";
            if (/不支持复制|浏览器不允许|剪贴板权限|HTTPS|localhost/.test(messageText)) {
                const url = await editorRef.current?.exportPng();
                if (url) {
                    downloadUrl(url, `mind-${Date.now()}.png`);
                    message.warning("当前浏览器环境不支持复制，已自动下载 PNG");
                    return;
                }
            }
            message.error(error instanceof Error ? error.message : "复制失败");
        }
    };

    const loadBenchmark = () => {
        const nextValue = createMindBenchmarkValue(120);
        editorRef.current?.importJson(nextValue);
        setValue(nextValue);
        message.success("已加载 120 节点测试数据");
    };

    return (
        <Body>
            <div className={styles.page}>
                <Card className={styles.sideCard} title={<Title>思维导图</Title>} variant="borderless">
                    <div className={styles.statGrid}>
                        <div className={styles.statItem}>
                            <strong>{stats.nodeCount}</strong>
                            <span>节点</span>
                        </div>
                        <div className={styles.statItem}>
                            <strong>{stats.depth}</strong>
                            <span>层级</span>
                        </div>
                    </div>
                    <Space wrap>
                        <Button icon={<PlusOutlined/>} type="primary" onClick={() => editorRef.current?.addChild()}>
                            子节点
                        </Button>
                        <Button icon={<ForkOutlined/>} onClick={() => editorRef.current?.addSibling()}>
                            同级
                        </Button>
                        <Button danger icon={<DeleteOutlined/>} onClick={() => editorRef.current?.removeSelected()}>
                            删除
                        </Button>
                        <Button icon={<UndoOutlined/>} onClick={() => editorRef.current?.undo()}>
                            撤销
                        </Button>
                        <Button icon={<ApartmentOutlined/>} onClick={() => editorRef.current?.layout()}>
                            整理
                        </Button>
                        <Button
                            icon={fullscreen ? <FullscreenExitOutlined/> : <FullscreenOutlined/>}
                            onClick={() => editorRef.current?.toggleFullscreen()}
                        >
                            {fullscreen ? "退出全屏" : "全屏"}
                        </Button>
                        <Button
                            icon={<ZoomInOutlined/>}
                            onClick={() => editorRef.current?.toggleControls()}
                        >
                            {controlsVisible ? "隐藏缩放" : "显示缩放"}
                        </Button>
                        <Button
                            icon={miniMapVisible ? <EyeInvisibleOutlined/> : <EyeOutlined/>}
                            onClick={() => editorRef.current?.toggleMiniMap()}
                        >
                            {miniMapVisible ? "隐藏缩略图" : "显示缩略图"}
                        </Button>
                        <Button icon={<SaveOutlined/>} onClick={() => editorRef.current?.save()}>
                            保存
                        </Button>
                        <Button icon={<SaveOutlined/>} onClick={copyJson}>
                            复制JSON
                        </Button>
                        <Button icon={<SaveOutlined/>} onClick={() => editorRef.current?.downloadJson()}>
                            导出JSON
                        </Button>
                        <Button icon={<SaveOutlined/>} onClick={() => editorRef.current?.downloadMarkdown()}>
                            导出MD
                        </Button>
                        <Button icon={<FileImageOutlined/>} onClick={exportPng}>
                            PNG
                        </Button>
                        <Button icon={<FileImageOutlined/>} onClick={exportJpg}>
                            JPG
                        </Button>
                        <Button icon={<CopyOutlined/>} onClick={copyImage}>
                            复制图片
                        </Button>
                        <Button icon={<ReloadOutlined/>} onClick={resetMind}>
                            重置
                        </Button>
                        <Button icon={<ThunderboltOutlined/>} onClick={loadBenchmark}>
                            120节点
                        </Button>
                    </Space>
                </Card>
                <Card
                    className={styles.canvasCard}
                    title={<Title>测试页面</Title>}
                    variant="borderless"
                    extra={(
                        <Tooltip title="双击节点编辑，选中节点可改分支颜色，Delete 删除，Ctrl/Cmd + Z 撤销">
                            <Typography.Text ellipsis className={styles.toolbarHint} type="secondary">
                                React Flow
                            </Typography.Text>
                        </Tooltip>
                    )}
                >
                    <div className={styles.toolbar}>
                        <Typography.Text strong>思维导图编辑器</Typography.Text>
                        <Typography.Text type="secondary">自动保存：{lastSave}</Typography.Text>
                    </div>
                    <Mind
                        ref={editorRef}
                        defaultValue={value}
                        height="calc(100vh - 218px)"
                        minHeight={560}
                        autoSave
                        autoSaveDelay={4000}
                        onChange={setValue}
                        onSave={(nextValue, nextStats, reason) => {
                            setValue(nextValue);
                            setLastSave(`${reason === "auto" ? "自动" : "手动"}保存于 ${new Date().toLocaleTimeString()}`);
                            if (reason !== "auto") {
                                message.success(`导图已保存，${nextStats.nodeCount} 个节点`);
                            }
                        }}
                        onControlsChange={setControlsVisible}
                        onFullscreenChange={setFullscreen}
                        onMiniMapChange={setMiniMapVisible}
                    />
                </Card>
            </div>
        </Body>
    );
};

export default TestMindPage;
