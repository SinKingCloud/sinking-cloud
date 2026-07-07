// noinspection TypeScriptValidateTypes,JSUnusedLocalSymbols,JSUnusedGlobalSymbols,SpellCheckingInspection,TypeScriptUnresolvedReference,TypeScriptRedundantGenericType,JSDeprecatedSymbols

import React, {
    forwardRef,
    useCallback,
    useEffect,
    useImperativeHandle,
    useLayoutEffect,
    useMemo,
    useRef,
    useState,
} from "react";
import {App, Popover, Tooltip} from "antd";
import {
    ApartmentOutlined,
    BgColorsOutlined,
    BoldOutlined,
    BorderOutlined,
    BoxPlotOutlined,
    BranchesOutlined,
    ClearOutlined,
    DeleteOutlined,
    EditOutlined,
    FontColorsOutlined,
    ForkOutlined,
    ItalicOutlined,
    PlusOutlined,
    StrikethroughOutlined,
    UnderlineOutlined,
    UndoOutlined,
} from "@ant-design/icons";
import {createStyles} from "antd-style";
import {
    applyNodeChanges,
    Background,
    BackgroundVariant,
    Controls,
    type Edge,
    Handle,
    MarkerType,
    MiniMap,
    type Node,
    type NodeChange,
    type NodeProps,
    Position,
    ReactFlow,
    type ReactFlowInstance,
    ReactFlowProvider,
} from "@xyflow/react";
import "@xyflow/react/dist/style.css";
import {useTheme} from "sinking-antd";

// ==================== 对外类型 ====================

/** 节点在画布中的坐标。 */
export type MindPoint = {
    /** 横向坐标。 */
    x: number;
    /** 纵向坐标。 */
    y: number;
};

/** 标题或描述的文字样式。 */
export type MindTextStyle = {
    /** 文字颜色，支持 CSS 颜色值。 */
    color?: string | undefined;
    /** 字号，单位 px。 */
    fontSize?: number | undefined;
    /** 字重，例如 600、bold、normal。 */
    fontWeight?: React.CSSProperties["fontWeight"] | undefined;
    /** 字体样式，例如 italic、normal。 */
    fontStyle?: React.CSSProperties["fontStyle"] | undefined;
    /** 文字装饰，例如 underline、line-through。 */
    textDecoration?: React.CSSProperties["textDecoration"] | undefined;
};

/** 节点整体和节点内文字的样式。 */
export type MindNodeStyle = {
    /** 节点填充色。 */
    background?: string | undefined;
    /** 标题默认颜色，未设置 titleStyle.color 时生效。 */
    color?: string | undefined;
    /** 节点边框颜色。 */
    borderColor?: string | undefined;
    /** 节点圆角，单位 px。 */
    borderRadius?: number | undefined;
    /** 标题默认字号，未设置 titleStyle.fontSize 时生效。 */
    fontSize?: number | undefined;
    /** 节点阴影。 */
    boxShadow?: string | undefined;
    /** 标题文字样式。 */
    titleStyle?: MindTextStyle | undefined;
    /** 描述文字样式。 */
    descriptionStyle?: MindTextStyle | undefined;
};

/** 思维导图节点数据。 */
export type MindNode = {
    /** 节点 ID，需要全局唯一。 */
    id: string;
    /** 节点标题纯文本。 */
    title: string;
    /** 节点标题富文本 HTML，组件会做基础清洗。 */
    titleHtml?: string;
    /** 节点描述纯文本。 */
    description?: string;
    /** 节点描述富文本 HTML，组件会做基础清洗。 */
    descriptionHtml?: string;
    /** 父节点 ID，不传表示根节点或一级节点。 */
    parentId?: string;
    /** 节点主题色。 */
    color?: string;
    /** 父节点连到当前节点这条分支线的颜色，未设置时跟随节点颜色。 */
    lineColor?: string;
    /** 是否为根节点。 */
    root?: boolean;
    /** 节点位置，不传时由组件布局生成。 */
    position?: MindPoint;
    /** 节点样式。 */
    style?: MindNodeStyle;
};

/** 更新节点时允许传入的字段，color / lineColor 传 null 表示恢复默认。 */
export type MindNodeUpdate = Partial<Pick<MindNode, "title" | "titleHtml" | "description" | "descriptionHtml" | "style">> & {
    /** 节点主题色，传 null 恢复默认色。 */
    color?: string | null;
    /** 分支线颜色，传 null 表示跟随节点色。 */
    lineColor?: string | null;
};

/** 思维导图完整数据。 */
export type MindValue = {
    /** 节点列表。 */
    nodes: MindNode[];
};

/** 思维导图统计信息。 */
export type MindStats = {
    /** 节点总数。 */
    nodeCount: number;
    /** 最大层级深度。 */
    depth: number;
};

/** Mind 组件通过 ref 暴露的命令式 API。 */
export type MindRef = {
    /** 获取当前思维导图数据。 */
    getValue: () => MindValue;
    /** getValue 的别名，便于业务语义化调用。 */
    getData: () => MindValue;
    /** 获取当前节点数量和层级深度。 */
    getStats: () => MindStats;
    /** 新增子节点，不传 parentId 时使用当前选中节点。 */
    addChild: (parentId?: string) => void;
    /** 新增同级节点，不传 nodeId 时使用当前选中节点。 */
    addSibling: (nodeId?: string) => void;
    /** 进入节点编辑模式，不传 nodeId 时编辑当前选中节点。 */
    edit: (nodeId?: string) => void;
    /** 更新指定节点数据。 */
    updateNode: (id: string, data: MindNodeUpdate) => void;
    /** 设置节点主题色。 */
    setNodeColor: (id: string, color: string) => void;
    /** 设置父节点连到当前节点的分支线颜色。 */
    setLineColor: (id: string, color: string) => void;
    /** 设置节点填充色。 */
    setBackgroundColor: (id: string, color: string) => void;
    /** 设置节点边框颜色。 */
    setBorderColor: (id: string, color: string) => void;
    /** 合并更新节点样式。 */
    setNodeStyle: (id: string, style: MindNodeStyle) => void;
    /** 清除节点主题色并恢复默认色。 */
    clearNodeColor: (id: string) => void;
    /** 清除分支线颜色并恢复跟随节点色。 */
    clearLineColor: (id: string) => void;
    /** 清除节点填充色。 */
    clearBackgroundColor: (id: string) => void;
    /** 清除节点边框颜色。 */
    clearBorderColor: (id: string) => void;
    /** 撤销最近一次数据变更，成功返回 true。 */
    undo: () => boolean;
    /** 当前是否可以撤销。 */
    canUndo: () => boolean;
    /** 删除当前选中的节点，根节点不会被删除。 */
    removeSelected: () => void;
    /** 按树结构重新整理布局。 */
    layout: () => void;
    /** layout 的别名。 */
    arrange: () => void;
    /** 适配视图到当前画布内容。 */
    fitView: () => void;
    /** 放大画布。 */
    zoomIn: () => void;
    /** 缩小画布。 */
    zoomOut: () => void;
    /** 重置为初始 value / defaultValue。 */
    reset: () => void;
    /** 当前是否全屏。 */
    isFullscreen: () => boolean;
    /** 进入全屏。 */
    openFullscreen: () => Promise<void>;
    /** 退出全屏。 */
    exitFullscreen: () => Promise<void>;
    /** 切换全屏状态。 */
    toggleFullscreen: () => Promise<void>;
    /** 当前是否显示左下角控制器。 */
    isControlsVisible: () => boolean;
    /** 显示左下角控制器。 */
    showControls: () => void;
    /** 隐藏左下角控制器。 */
    hideControls: () => void;
    /** 切换左下角控制器显示状态。 */
    toggleControls: () => void;
    /** 当前是否显示右下角缩略图。 */
    isMiniMapVisible: () => boolean;
    /** 显示右下角缩略图。 */
    showMiniMap: () => void;
    /** 隐藏右下角缩略图。 */
    hideMiniMap: () => void;
    /** 切换右下角缩略图显示状态。 */
    toggleMiniMap: () => void;
};

/** Mind 组件属性。 */
export type MindProps = {
    /** 受控数据，传入后由外部管理数据更新。 */
    value?: MindValue;
    /** 非受控初始数据。 */
    defaultValue?: MindValue;
    /** 外层 className。 */
    className?: string;
    /** 外层样式。 */
    style?: React.CSSProperties;
    /** 组件高度，默认 100%。 */
    height?: number | string;
    /** 组件最小高度，默认 0。 */
    minHeight?: number | string;
    /** 只读模式，禁用编辑、拖拽和节点操作。 */
    readOnly?: boolean;
    /** 受控显示左下角控制器。 */
    showControls?: boolean;
    /** 非受控默认是否显示左下角控制器，默认 true。 */
    defaultShowControls?: boolean;
    /** 受控显示右下角缩略图。 */
    showMiniMap?: boolean;
    /** 非受控默认是否显示右下角缩略图，默认 true。 */
    defaultShowMiniMap?: boolean;
    /** 新建节点默认标题。 */
    newNodeTitle?: string;
    /** 自定义节点 ID 生成器。 */
    createId?: () => string;
    /** 数据变更回调。 */
    onChange?: (value: MindValue, stats: MindStats) => void;
    /** onChange 的别名，便于业务语义化监听。 */
    onDataChange?: (value: MindValue, stats: MindStats) => void;
    /** 左下角控制器显示状态变化回调。 */
    onControlsChange?: (visible: boolean) => void;
    /** 全屏状态变化回调。 */
    onFullscreenChange?: (fullscreen: boolean) => void;
    /** 右下角缩略图显示状态变化回调。 */
    onMiniMapChange?: (visible: boolean) => void;
};

// ==================== 内部类型 ====================

type FlowNodeData = {
    title: string;
    titleHtml?: string;
    description?: string;
    descriptionHtml?: string;
    parentId?: string;
    color: string;
    lineColor: string;
    root?: boolean;
    style?: MindNodeStyle;
};

type FlowNode = Node<FlowNodeData, "mind">;

type MindNodeViewProps = NodeProps<FlowNode> & {
    styles: ReturnType<typeof useStyles>["styles"];
    readOnly?: boolean;
    editingId?: string;
    onStartEdit: (id: string) => void;
    onUpdateNode: (id: string, data: MindNodeUpdate) => void;
    onRemoveNode: (id: string) => void;
    onEditDone: () => void;
    onTextTargetChange: (target: TextTarget) => void;
    onRichSelectionChange: () => void;
    onShouldKeepEdit: () => boolean;
    onTextSelectionStart: () => void;
};

type ContextMenu = {
    x: number;
    y: number;
    nodeId?: string;
    position?: MindPoint;
};

type MindContextValue = Omit<MindNodeViewProps, keyof NodeProps<FlowNode>>;

type TextTarget = "title" | "description";

const MindContext = React.createContext<MindContextValue | null>(null);

// ==================== 常量配置 ====================

const rootId = "root";
const nodeWidth = 196;
const nodeMinHeight = 46;
const titleLineHeight = 20;
const descriptionLineHeight = 18;
const defaultFontSize = 13;
const boldFontWeight = 600;
const defaultDescriptionFontSize = 12;
const defaultBorderRadius = 6;
const minTextFontSize = 11;
const maxTextFontSize = 22;
const levelGap = 250;
const rowGap = 28;
const palette = ["#1677ff", "#13c2c2", "#52c41a", "#faad14", "#eb2f96", "#722ed1", "#fa541c"];
const stylePalette = ["#ffffff", "#000000", "#8c8c8c", "#f5222d", "#fa8c16", "#fadb14", "#52c41a", "#13c2c2", "#1677ff", "#722ed1", "#eb2f96"];
const nodeShadow = "0 8px 20px rgba(15,23,42,0.14)";
const reactFlowNoDragClass = "no" + "drag";
const reactFlowNoPanClass = "no" + "pan";
const reactFlowNoWheelClass = "no" + "wheel";
const reactFlowLockClassName = [reactFlowNoDragClass, reactFlowNoPanClass, reactFlowNoWheelClass].join(" ");

// ==================== 富文本处理 ====================
const normalizeTitle = (title?: string) => (title || "").replace(/\s+/g, " ").trim();

const normalizeDescription = (description?: string) => {
    const value = description?.trim();
    return value || undefined;
};

const escapeHtml = (value?: string) => (value || "")
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#39;");

const textToRichHtml = (value?: string) => {
    const text = value || "";
    return text ? escapeHtml(text).replace(/\n/g, "<br>") : undefined;
};

const getRichStyle = (element: HTMLElement) => {
    const styles: string[] = [];
    const color = element.style.color || element.getAttribute("color");
    const fontSize = element.style.fontSize;
    const fontWeight = element.style.fontWeight;
    const fontStyle = element.style.fontStyle;
    const textDecoration = element.style.textDecoration || element.style.textDecorationLine;
    if (color) {
        styles.push(`color:${color}`);
    }
    if (/^\d+(\.\d+)?px$/.test(fontSize)) {
        styles.push(`font-size:${fontSize}`);
    }
    if (fontWeight === "normal" || fontWeight === "bold" || /^[1-9]00$/.test(fontWeight)) {
        styles.push(`font-weight:${fontWeight}`);
    }
    if (fontStyle === "normal" || fontStyle === "italic") {
        styles.push(`font-style:${fontStyle}`);
    }
    if (textDecoration) {
        const values = textDecoration.split(/\s+/).filter((item) => ["none", "underline", "line-through"].includes(item));
        if (values.length > 0) {
            styles.push(`text-decoration:${Array.from(new Set(values)).join(" ")}`);
        }
    }
    return styles.join(";");
};

const sanitizeRichHtml = (html?: string) => {
    if (!html) {
        return undefined;
    }
    if (typeof document === "undefined") {
        return escapeHtml(html.replace(/<br\s*\/?>/gi, "\n").replace(/<[^>]+>/g, ""));
    }
    const template = document.createElement("template");
    template.innerHTML = html;
    const sanitizeNode = (node: globalThis.Node): string => {
        if (node.nodeType === globalThis.Node.TEXT_NODE) {
            return escapeHtml(node.textContent || "");
        }
        if (!(node instanceof HTMLElement)) {
            return "";
        }
        const children = Array.from(node.childNodes).map(sanitizeNode).join("");
        const tag = node.tagName.toLowerCase();
        if (tag === "br") {
            return "<br>";
        }
        if (tag === "div" || tag === "p") {
            return children ? `${children}<br>` : "<br>";
        }
        if (tag === "b" || tag === "strong") {
            return children ? `<strong>${children}</strong>` : "";
        }
        if (tag === "i" || tag === "em") {
            return children ? `<em>${children}</em>` : "";
        }
        if (tag === "u") {
            return children ? `<u>${children}</u>` : "";
        }
        if (tag === "s" || tag === "strike") {
            return children ? `<s>${children}</s>` : "";
        }
        if (tag === "span" || tag === "font") {
            const style = getRichStyle(node);
            return style && children ? `<span style="${style}">${children}</span>` : children;
        }
        return children;
    };
    const content = Array.from(template.content.childNodes).map(sanitizeNode).join("").replace(/(<br>)+$/g, "");
    return content || undefined;
};

const getRichText = (html?: string) => {
    if (!html) {
        return "";
    }
    if (typeof document === "undefined") {
        return html.replace(/<br\s*\/?>/gi, "\n").replace(/<[^>]+>/g, "");
    }
    const template = document.createElement("template");
    template.innerHTML = (sanitizeRichHtml(html) || "").replace(/<br>/g, "\n");
    return (template.content.textContent || "").replace(/\u00a0/g, " ");
};

const normalizeRichHtml = (html?: string) => {
    const value = sanitizeRichHtml(html);
    return getRichText(value).trim() ? value : undefined;
};

const hasStyledRichHtml = (html?: string) => /<(strong|em|u|s|span)\b/i.test(html || "");

const hasField = <T extends object>(data: T, key: keyof T) => Object.prototype.hasOwnProperty.call(data, key);

// ==================== 富文本选区样式 ====================

const setElementStyle = (element: HTMLElement, style: React.CSSProperties) => {
    Object.entries(style).forEach(([key, value]) => {
        if (value === undefined || value === null || value === "") {
            return;
        }
        const nextValue = key === "fontSize" && typeof value === "number" ? `${value}px` : String(value);
        if (key === "fontSize") {
            element.style.setProperty("font-size", nextValue);
        } else if (key === "fontWeight") {
            element.style.setProperty("font-weight", nextValue);
        } else if (key === "fontStyle") {
            element.style.setProperty("font-style", nextValue);
        } else if (key === "textDecoration") {
            element.style.setProperty("text-decoration-line", nextValue);
        } else if (key === "borderColor") {
            element.style.setProperty("border-color", nextValue);
        } else if (key === "boxShadow") {
            element.style.setProperty("box-shadow", nextValue);
        } else {
            element.style.setProperty(key.replace(/[A-Z]/g, (match) => `-${match.toLowerCase()}`), nextValue);
        }
    });
};

const unwrapElement = (element: HTMLElement) => {
    const parent = element.parentNode;
    if (!parent) {
        return;
    }
    while (element.firstChild) {
        parent.insertBefore(element.firstChild, element);
    }
    parent.removeChild(element);
};

const clearEmptyRichElement = (element: HTMLElement) => {
    const tag = element.tagName.toLowerCase();
    const style = element.getAttribute("style")?.trim();
    if (!style) {
        element.removeAttribute("style");
    }
    if ((tag === "span" || tag === "font") && !element.getAttribute("style") && !element.getAttribute("color")) {
        unwrapElement(element);
    }
};

const clearRichFragmentStyle = (fragment: DocumentFragment, style: React.CSSProperties) => {
    const targets = {
        color: hasField(style, "color"),
        fontSize: hasField(style, "fontSize"),
        fontWeight: hasField(style, "fontWeight"),
        fontStyle: hasField(style, "fontStyle"),
        textDecoration: hasField(style, "textDecoration"),
    };
    const elements = Array.from(fragment.querySelectorAll("*")).reverse() as HTMLElement[];

    // 先清理选区内部同类型样式，避免外层取消后内层样式继续生效。
    elements.forEach((element) => {
        if (!element.parentNode) {
            return;
        }
        const tag = element.tagName.toLowerCase();
        if (targets.fontWeight && (tag === "strong" || tag === "b")) {
            unwrapElement(element);
            return;
        }
        if (targets.fontStyle && (tag === "em" || tag === "i")) {
            unwrapElement(element);
            return;
        }
        if (targets.textDecoration && (tag === "u" || tag === "s" || tag === "strike")) {
            unwrapElement(element);
            return;
        }
        if (targets.color) {
            element.style.removeProperty("color");
            element.removeAttribute("color");
        }
        if (targets.fontSize) {
            element.style.removeProperty("font-size");
        }
        if (targets.fontWeight) {
            element.style.removeProperty("font-weight");
        }
        if (targets.fontStyle) {
            element.style.removeProperty("font-style");
        }
        if (targets.textDecoration) {
            element.style.removeProperty("text-decoration");
            element.style.removeProperty("text-decoration-line");
        }
        clearEmptyRichElement(element);
    });
};

const getRichEditorBaseStyle = (editor: HTMLElement): React.CSSProperties => {
    const computed = window.getComputedStyle(editor);
    const fontSize = Number.parseFloat(computed.fontSize);
    return {
        color: computed.color,
        fontSize: Number.isFinite(fontSize) ? Math.round(fontSize) : undefined,
        fontWeight: computed.fontWeight,
        fontStyle: computed.fontStyle === "italic" ? "italic" : "normal",
        textDecoration: computed.textDecorationLine === "none" ? "none" : computed.textDecorationLine,
    };
};

const getAppliedRichStyle = (style: React.CSSProperties, editor: HTMLElement) => {
    const next = {...style};
    if (next.color === "inherit") {
        next.color = getRichEditorBaseStyle(editor).color;
    }
    return next;
};

const toHexColor = (value?: string) => {
    if (!value) {
        return "#1677ff";
    }
    if (/^#[0-9a-f]{6}$/i.test(value)) {
        return value;
    }
    if (/^#[0-9a-f]{3}$/i.test(value)) {
        return `#${value.slice(1).split("").map((item) => item + item).join("")}`;
    }
    const match = value.match(/rgba?\((\d+),\s*(\d+),\s*(\d+)/i);
    if (!match) {
        return "#1677ff";
    }
    return `#${[match[1], match[2], match[3]].map((item) => Math.max(0, Math.min(255, Number(item))).toString(16).padStart(2, "0")).join("")}`;
};

const createPlainRichFragment = (text: string) => {
    const fragment = document.createDocumentFragment();
    text.split("\n").forEach((line, index) => {
        if (index > 0) {
            fragment.appendChild(document.createElement("br"));
        }
        fragment.appendChild(document.createTextNode(line));
    });
    return fragment;
};

const getRangeText = (range: Range, node: globalThis.Node) => {
    const text = node.textContent || "";
    if (node === range.startContainer && node === range.endContainer) {
        return text.slice(range.startOffset, range.endOffset);
    }
    if (node === range.startContainer) {
        return text.slice(range.startOffset);
    }
    if (node === range.endContainer) {
        return text.slice(0, range.endOffset);
    }
    return text;
};

const getRichRangeStyle = (range: Range, editor: HTMLElement): MindTextStyle => {
    const nodes: Text[] = [];
    const walker = document.createTreeWalker(editor, globalThis.NodeFilter.SHOW_TEXT, {
        acceptNode: (node) => {
            if (!range.intersectsNode(node) || !getRangeText(range, node).trim()) {
                return globalThis.NodeFilter.FILTER_REJECT;
            }
            return globalThis.NodeFilter.FILTER_ACCEPT;
        },
    });
    while (walker.nextNode()) {
        nodes.push(walker.currentNode as Text);
    }
    if (nodes.length <= 0) {
        const element = range.startContainer instanceof HTMLElement ? range.startContainer : range.startContainer.parentElement;
        const computed = window.getComputedStyle(element && editor.contains(element) ? element : editor);
        const fontSize = Number.parseFloat(computed.fontSize);
        return {
            color: computed.color,
            fontSize: Number.isFinite(fontSize) ? Math.round(fontSize) : undefined,
            fontWeight: computed.fontWeight,
            fontStyle: computed.fontStyle === "italic" ? "italic" : "normal",
            textDecoration: computed.textDecorationLine === "none" ? "none" : computed.textDecorationLine,
        };
    }

    const colors = new Set<string>();
    const fontSizes = new Set<number>();
    const bold = new Set<boolean>();
    const italic = new Set<boolean>();
    const underline = new Set<boolean>();
    const strike = new Set<boolean>();

    nodes.forEach((node) => {
        const element = node.parentElement || editor;
        const computed = window.getComputedStyle(element);
        const fontSize = Number.parseFloat(computed.fontSize);
        const fontWeight = computed.fontWeight;
        const decoration = computed.textDecorationLine || "";
        colors.add(computed.color);
        if (Number.isFinite(fontSize)) {
            fontSizes.add(Math.round(fontSize));
        }
        bold.add(fontWeight === "bold" || Number(fontWeight) >= 600);
        italic.add(computed.fontStyle === "italic");
        underline.add(decoration.includes("underline"));
        strike.add(decoration.includes("line-through"));
    });

    const textDecoration = [
        underline.size === 1 && underline.has(true) ? "underline" : "",
        strike.size === 1 && strike.has(true) ? "line-through" : "",
    ].filter(Boolean).join(" ");

    return {
        color: colors.size === 1 ? Array.from(colors)[0] : undefined,
        fontSize: fontSizes.size === 1 ? Array.from(fontSizes)[0] : undefined,
        fontWeight: bold.size === 1 && bold.has(true) ? 600 : "normal",
        fontStyle: italic.size === 1 && italic.has(true) ? "italic" : "normal",
        textDecoration: textDecoration || "none",
    };
};

const cleanMindTextStyle = (style?: MindTextStyle) => {
    if (!style) {
        return undefined;
    }
    const next: MindTextStyle = {...style};
    (Object.keys(next) as (keyof MindTextStyle)[]).forEach((key) => {
        if (next[key] === undefined || next[key] === "") {
            delete next[key];
        }
    });
    return Object.keys(next).length > 0 ? next : undefined;
};

const cleanMindNodeStyle = (style?: MindNodeStyle) => {
    if (!style) {
        return undefined;
    }
    const next: MindNodeStyle = {
        ...style,
        titleStyle: cleanMindTextStyle(style.titleStyle),
        descriptionStyle: cleanMindTextStyle(style.descriptionStyle),
    };
    (Object.keys(next) as (keyof MindNodeStyle)[]).forEach((key) => {
        if (next[key] === undefined || next[key] === "") {
            delete next[key];
        }
    });
    return Object.keys(next).length > 0 ? next : undefined;
};

const cloneMindNodeStyle = (style?: MindNodeStyle) => cleanMindNodeStyle(style ? {
    ...style,
    titleStyle: style.titleStyle ? {...style.titleStyle} : undefined,
    descriptionStyle: style.descriptionStyle ? {...style.descriptionStyle} : undefined,
} : undefined);

const mergeMindNodeStyle = (style?: MindNodeStyle, patch?: MindNodeStyle) => {
    if (!patch) {
        return cleanMindNodeStyle(style);
    }
    const next: MindNodeStyle = {...style, ...patch};
    if (hasField(patch, "titleStyle")) {
        next.titleStyle = patch.titleStyle ? {...style?.titleStyle, ...patch.titleStyle} : undefined;
    }
    if (hasField(patch, "descriptionStyle")) {
        next.descriptionStyle = patch.descriptionStyle ? {...style?.descriptionStyle, ...patch.descriptionStyle} : undefined;
    }
    return cleanMindNodeStyle(next);
};

const isInputTarget = (target: EventTarget | null) => {
    const element = target as HTMLElement | null;
    return element instanceof HTMLInputElement ||
        element instanceof HTMLTextAreaElement ||
        element instanceof HTMLSelectElement ||
        Boolean(element?.isContentEditable);
};

const isMindInteractiveTarget = (target: EventTarget | null) => {
    const element = target as HTMLElement | null;
    return Boolean(element?.closest("[data-mind-toolbar='true'], [data-mind-color-picker='true'], [data-mind-edit-content='true'], [data-mind-rich-editor='true'], .ant-popover"));
};

const isMindColorTarget = (target: EventTarget | null) => {
    const element = target as HTMLElement | null;
    return Boolean(element?.closest("[data-mind-color-picker='true']"));
};

const isMindColorDropdownTarget = (target: EventTarget | null) => {
    const element = target as HTMLElement | null;
    return Boolean(element?.closest(".ant-popover"));
};

// ==================== 数据转换 ====================

/** 创建并规范化一个思维导图节点，适合外部生成初始 value 时使用。 */
export const createMindNode = (data: MindNode): MindNode => {
    const color = data.color || palette[0];
    const titleHtml = normalizeRichHtml(data.titleHtml);
    const descriptionHtml = normalizeRichHtml(data.descriptionHtml);
    return {
        ...data,
        title: normalizeTitle(data.title || getRichText(titleHtml)),
        titleHtml,
        description: normalizeDescription(data.description || getRichText(descriptionHtml)),
        descriptionHtml,
        color,
        lineColor: data.lineColor || color,
        style: cleanMindNodeStyle(data.style),
    };
};

const getDefaultId = () => `node_${Date.now()}_${Math.random().toString(36).slice(2, 7)}`;

const getTextLines = (text?: string, lineChars = 14) => {
    if (!text) {
        return 0;
    }
    return text.split(/\r?\n/).reduce((total, line) => {
        return total + Math.max(1, Math.ceil((line.trim().length || 1) / lineChars));
    }, 0);
};

const getTextStyleKey = (target: TextTarget) => target === "title" ? "titleStyle" : "descriptionStyle";

const getToolbarTextStyle = (style: MindNodeStyle | undefined, target: TextTarget): MindTextStyle => {
    if (target === "title") {
        return {
            color: style?.color,
            fontSize: style?.fontSize || defaultFontSize,
            ...style?.titleStyle,
        };
    }
    return {...style?.descriptionStyle};
};

const getTitleLineHeight = (style?: MindNodeStyle) => {
    const fontSize = style?.titleStyle?.fontSize || style?.fontSize || defaultFontSize;
    return Math.max(titleLineHeight, Math.ceil(fontSize * 1.45));
};

const getDescriptionLineHeight = (style?: MindNodeStyle) => {
    const fontSize = style?.descriptionStyle?.fontSize || defaultDescriptionFontSize;
    return Math.max(descriptionLineHeight, Math.ceil(fontSize * 1.45));
};

const getNodeBoxStyle = (style?: MindNodeStyle): React.CSSProperties => ({
    background: style?.background,
    borderColor: style?.borderColor,
    borderRadius: style?.borderRadius,
    boxShadow: style?.boxShadow,
});

const getTitleTextStyle = (style?: MindNodeStyle): React.CSSProperties => ({
    color: style?.titleStyle?.color || style?.color,
    fontSize: style?.titleStyle?.fontSize || style?.fontSize || defaultFontSize,
    fontWeight: style?.titleStyle?.fontWeight,
    fontStyle: style?.titleStyle?.fontStyle,
    textDecoration: style?.titleStyle?.textDecoration,
    lineHeight: `${getTitleLineHeight(style)}px`,
});

const getDescriptionTextStyle = (style?: MindNodeStyle): React.CSSProperties => ({
    ...style?.descriptionStyle,
    lineHeight: `${getDescriptionLineHeight(style)}px`,
});

const toggleTextDecoration = (value: React.CSSProperties["textDecoration"], decoration: string) => {
    const items = new Set(String(value || "").split(/\s+/).filter((item) => item && item !== "none"));
    if (items.has(decoration)) {
        items.delete(decoration);
    } else {
        items.add(decoration);
    }
    return items.size > 0 ? Array.from(items).join(" ") : undefined;
};

const hasTextDecoration = (value: React.CSSProperties["textDecoration"], decoration: string) => {
    return String(value || "").split(/\s+/).includes(decoration);
};

const getMindNodeHeight = (node: Pick<MindNode, "title" | "description"> & { style?: MindNodeStyle }) => {
    const descriptionLines = getTextLines(node.description, 18);
    const descriptionHeight = descriptionLines > 0 ? 4 + descriptionLines * getDescriptionLineHeight(node.style) : 0;
    return Math.max(nodeMinHeight, 16 + getTitleLineHeight(node.style) + descriptionHeight);
};

// ==================== 画布节点转换 ====================

const cloneFlowNodes = (nodes: FlowNode[]) => nodes.map((node) => ({
    ...node,
    data: {
        ...node.data,
        style: cloneMindNodeStyle(node.data.style),
    },
    position: {...node.position},
}));

const toFlowNode = (node: MindNode): FlowNode => {
    const titleHtml = normalizeRichHtml(node.titleHtml);
    const descriptionHtml = normalizeRichHtml(node.descriptionHtml);
    const title = normalizeTitle(node.title || getRichText(titleHtml));
    const description = normalizeDescription(node.description || getRichText(descriptionHtml));
    return {
        id: node.id,
        type: "mind",
        position: node.position || {x: 0, y: 0},
        width: nodeWidth,
        height: getMindNodeHeight({title, description, style: node.style}),
        data: {
            title,
            titleHtml,
            description,
            descriptionHtml,
            parentId: node.parentId,
            color: node.color || palette[0],
            lineColor: node.lineColor || node.color || palette[0],
            root: node.root,
            style: cleanMindNodeStyle(node.style),
        },
    };
};

const toMindValue = (nodes: FlowNode[]): MindValue => ({
    nodes: nodes.map((node) => ({
        id: node.id,
        title: node.data.title,
        titleHtml: node.data.titleHtml,
        description: node.data.description,
        descriptionHtml: node.data.descriptionHtml,
        parentId: node.data.parentId,
        color: node.data.color,
        lineColor: node.data.lineColor,
        root: node.data.root,
        position: node.position,
        style: cleanMindNodeStyle(node.data.style),
    })),
});

// ==================== 树结构计算 ====================
const getFlowChildrenMap = (nodes: FlowNode[]) => {
    const map = new Map<string, FlowNode[]>();
    nodes.forEach((node) => {
        if (!node.data.parentId) {
            return;
        }
        map.set(node.data.parentId, (map.get(node.data.parentId) || []).concat(node));
    });
    return map;
};

const getDepth = <T extends { id: string; parentId?: string }>(nodes: T[], id = rootId): number => {
    const children = nodes.filter((node) => node.parentId === id);
    if (children.length <= 0) {
        return nodes.some((node) => node.id === id) ? 1 : 0;
    }
    return 1 + Math.max(...children.map((node) => getDepth(nodes, node.id)));
};

/** 统计思维导图节点数量和最大层级深度。 */
export const getMindStats = (value: MindValue): MindStats => ({
    nodeCount: value.nodes.length,
    depth: getDepth(value.nodes),
});

const getFlowStats = (nodes: FlowNode[]) => getMindStats(toMindValue(nodes));

// ==================== 自动布局 ====================
const layoutNodes = (nodes: FlowNode[], layoutRootId = rootId): FlowNode[] => {
    const childrenMap = getFlowChildrenMap(nodes);
    const nodeMap = new Map(nodes.map((node) => [node.id, node]));
    const positionMap = new Map<string, MindPoint>();

    const sortNodes = (items: FlowNode[]) => [...items].sort((a, b) => a.position.y - b.position.y || a.position.x - b.position.x);

    const getFlowNodeHeight = (id: string) => {
        const node = nodeMap.get(id);
        if (!node) {
            return nodeMinHeight;
        }
        return getMindNodeHeight({
            title: node.data.title,
            description: node.data.description,
            style: node.data.style,
        });
    };

    const getChildrenHeight = (children: FlowNode[]) => {
        if (children.length <= 0) {
            return 0;
        }
        return children.reduce((total, child, index) => {
            return total + getSubtreeHeight(child.id) + (index > 0 ? rowGap : 0);
        }, 0);
    };

    const getSubtreeHeight = (id: string): number => {
        const children = sortNodes(childrenMap.get(id) || []);
        if (children.length <= 0) {
            return getFlowNodeHeight(id);
        }
        return Math.max(getFlowNodeHeight(id), getChildrenHeight(children));
    };

    const placeNode = (id: string, depth: number, top: number): { nextTop: number; center: number } => {
        const children = sortNodes(childrenMap.get(id) || []);
        const subtreeHeight = getSubtreeHeight(id);
        if (children.length <= 0) {
            const center = top + subtreeHeight / 2;
            positionMap.set(id, {x: depth * levelGap, y: center});
            return {nextTop: top + subtreeHeight + rowGap, center};
        }

        const childrenHeight = getChildrenHeight(children);
        let cursor = top + Math.max(0, (subtreeHeight - childrenHeight) / 2);
        const centers = children.map((child) => {
            const placed = placeNode(child.id, depth + 1, cursor);
            cursor = placed.nextTop;
            return placed.center;
        });
        const center = (centers[0] + centers[centers.length - 1]) / 2;
        positionMap.set(id, {x: depth * levelGap, y: center});
        return {nextTop: top + subtreeHeight + rowGap, center};
    };

    const rootHeight = getSubtreeHeight(layoutRootId);
    placeNode(layoutRootId, 0, -rootHeight / 2);

    return nodes.map((node) => {
        const point = positionMap.get(node.id);
        const height = getMindNodeHeight({
            title: node.data.title,
            description: node.data.description,
            style: node.data.style,
        });
        if (!point) {
            return node;
        }
        return {
            ...node,
            height,
            position: {
                x: point.x - nodeWidth / 2,
                y: point.y - height / 2,
            },
        };
    });
};

const normalizeValue = (value?: MindValue) => {
    const sourceNodes = value?.nodes || [];
    const flowNodes = sourceNodes.map(toFlowNode);
    return sourceNodes.some((node) => !node.position) ? layoutNodes(flowNodes) : flowNodes;
};

const selectOnly = (nodes: FlowNode[], ids?: string | string[]) => {
    const selectedIds = new Set(Array.isArray(ids) ? ids : (ids ? [ids] : []));
    return nodes.map((node) => ({...node, selected: selectedIds.has(node.id)}));
};

const getDescendantIds = (nodes: FlowNode[], ids: string[]) => {
    const childrenMap = getFlowChildrenMap(nodes);
    const result = new Set(ids);
    const walk = (id: string) => {
        (childrenMap.get(id) || []).forEach((child) => {
            result.add(child.id);
            walk(child.id);
        });
    };
    ids.forEach(walk);
    return result;
};

const toEdges = (nodes: FlowNode[]): Edge[] => nodes
    .filter((node) => Boolean(node.data.parentId))
    .map((node) => {
        const lineColor = node.data.lineColor || node.data.color;
        return {
            id: `${node.data.parentId}-${node.id}`,
            source: node.data.parentId!,
            target: node.id,
            type: "default",
            sourceHandle: "right",
            targetHandle: "left",
            markerEnd: {
                type: MarkerType.ArrowClosed,
                color: lineColor,
                width: 14,
                height: 14,
            },
            style: {
                stroke: lineColor,
                strokeWidth: 1.4,
                opacity: 0.72,
            },
        };
    });

// ==================== 组件样式 ====================
const useStyles = createStyles(({css, token, isDarkMode}: any) => ({
    editor: css`
        position: relative;
        display: flex;
        min-width: 0;
        width: 100%;
        height: 100%;
        overflow: hidden;
        outline: none;
        background: ${token.colorBgLayout};

        .react-flow {
            flex: 1;
            min-width: 0;
            width: 100%;
            height: 100%;
        }

        .react-flow__pane {
            cursor: grab;
        }

        .react-flow__pane:active {
            cursor: grabbing;
        }

        .react-flow__attribution {
            display: none;
        }

        .react-flow__controls {
            overflow: hidden;
            border: 1px solid ${token.colorBorderSecondary};
            border-radius: ${token.borderRadius}px;
            background: ${token.colorBgContainer};
            box-shadow: ${isDarkMode ? "0 6px 16px rgba(0,0,0,0.22)" : "0 6px 16px rgba(0,0,0,0.06)"};
        }

        .react-flow__controls-button {
            width: 30px;
            height: 30px;
            border-bottom: 1px solid ${token.colorBorderSecondary};
            background: ${token.colorBgContainer};
            color: ${token.colorTextSecondary};
        }

        .react-flow__controls-button:hover {
            background: ${token.colorFillTertiary};
            color: ${token.colorPrimary};
        }

        .react-flow__controls-button svg {
            width: 14px;
            height: 14px;
        }

        .react-flow__minimap {
            overflow: hidden;
            border: 1px solid ${token.colorBorderSecondary};
            border-radius: ${token.borderRadius}px;
            background: ${token.colorBgContainer};
            box-shadow: ${isDarkMode ? "0 6px 16px rgba(0,0,0,0.22)" : "0 6px 16px rgba(0,0,0,0.06)"};
        }

        .react-flow__selection {
            border: 1px solid ${token.colorPrimary};
            background: ${isDarkMode ? "rgba(22,119,255,0.18)" : token.colorPrimaryBg};
        }
    `,
    fullscreen: css`
        width: 100vw !important;
        height: 100vh !important;
        min-height: 0 !important;
    `,
    node: css`
        position: relative;
        display: flex;
        align-items: center;
        width: ${nodeWidth}px;
        min-height: ${nodeMinHeight}px;
        padding: 8px 12px;
        border: 1px solid ${token.colorBorderSecondary};
        border-radius: ${token.borderRadius}px;
        background: ${token.colorBgContainer};
        color: ${token.colorTextHeading};
        box-shadow: ${isDarkMode ? "0 4px 12px rgba(0,0,0,0.18)" : "0 4px 12px rgba(0,0,0,0.04)"};
        font-size: 13px;
        line-height: 20px;
        font-weight: 400;
        cursor: grab;
        user-select: none;

        &:hover {
            border-color: ${token.colorPrimaryBorderHover};
        }
    `,
    nodeRoot: css`
        border-color: ${token.colorPrimaryBorder};
        background: ${token.colorPrimaryBg};
        color: ${token.colorPrimary};
        font-weight: 600;
    `,
    nodeEditing: css`
        cursor: text;
        user-select: text;
    `,
    nodeContent: css`
        width: 100%;
        min-width: 0;
    `,
    nodeTitle: css`
        min-width: 0;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
    `,
    nodeDescription: css`
        margin-top: 4px;
        color: ${token.colorTextSecondary};
        font-size: 12px;
        line-height: ${descriptionLineHeight}px;
        font-weight: 400;
        white-space: pre-wrap;
        word-break: break-word;
        overflow-wrap: anywhere;
    `,
    nodeEditContent: css`
        display: grid;
        width: 100%;
        flex: 1;
        gap: 6px;
        min-width: 0;
        user-select: text;
    `,
    nodeTitleInput: css`
        width: 100%;
        min-height: ${titleLineHeight}px;
        box-sizing: border-box;
        padding: 0;
        border: 0;
        border-radius: 0;
        outline: none;
        background: transparent;
        color: inherit;
        font-size: 13px;
        line-height: ${titleLineHeight}px;
        white-space: nowrap;
        overflow: hidden;
        cursor: text;
        user-select: text;

        &:focus {
            box-shadow: none;
        }

        &:empty:before {
            content: attr(data-placeholder);
            color: ${token.colorTextQuaternary};
            pointer-events: none;
        }
    `,
    nodeDescriptionInput: css`
        width: 100%;
        min-height: 38px;
        box-sizing: border-box;
        padding: 0;
        border: 0;
        border-radius: 0;
        outline: none;
        background: transparent;
        color: ${token.colorTextSecondary};
        font-size: 12px;
        line-height: ${descriptionLineHeight}px;
        white-space: pre-wrap;
        word-break: break-word;
        overflow-wrap: anywhere;
        cursor: text;
        user-select: text;

        &:focus {
            box-shadow: none;
        }

        &:empty:before {
            content: attr(data-placeholder);
            color: ${token.colorTextQuaternary};
            pointer-events: none;
        }
    `,
    styleToolbar: css`
        position: absolute;
        top: 12px;
        left: 50%;
        z-index: 12;
        display: flex;
        align-items: center;
        gap: 4px;
        flex-wrap: wrap;
        justify-content: center;
        max-width: calc(100% - 32px);
        padding: 5px;
        border: 1px solid ${token.colorBorderSecondary};
        border-radius: ${token.borderRadiusLG}px;
        background: ${isDarkMode ? "rgba(20,20,20,0.94)" : "rgba(255,255,255,0.96)"};
        box-shadow: ${token.boxShadowSecondary};
        backdrop-filter: blur(12px);
        transform: translateX(-50%);
    `,
    styleToolbarGroup: css`
        display: flex;
        align-items: center;
        gap: 2px;
        min-width: 0;
        height: 30px;
        padding: 2px;
        border-radius: ${token.borderRadius}px;
        background: ${token.colorFillQuaternary};
    `,
    styleToolbarTooltip: css`
        display: inline-flex;
        align-items: center;
        line-height: 0;
    `,
    styleToolbarDivider: css`
        width: 1px;
        height: 24px;
        margin: 0 2px;
        background: ${token.colorBorderSecondary};
    `,
    styleToolbarButton: css`
        display: inline-flex;
        align-items: center;
        justify-content: center;
        width: 26px;
        height: 26px;
        padding: 0;
        border: 0;
        border-radius: ${token.borderRadiusSM}px;
        background: transparent;
        color: ${token.colorTextSecondary};
        font-size: 14px;
        line-height: 1;
        cursor: pointer;

        &:hover {
            color: ${token.colorPrimary};
            background: ${token.colorFillTertiary};
        }

        &:disabled {
            color: ${token.colorTextDisabled};
            cursor: not-allowed;
            background: transparent;
        }
    `,
    styleToolbarButtonActive: css`
        color: ${token.colorPrimary};
        background: ${token.colorPrimaryBg};
    `,
    styleToolbarButtonDisabled: css`
        color: ${token.colorTextDisabled};
        cursor: not-allowed;
        background: transparent;

        &:hover {
            color: ${token.colorTextDisabled};
            background: transparent;
        }
    `,
    styleToolbarValue: css`
        min-width: 22px;
        color: ${token.colorTextTertiary};
        font-size: 12px;
        line-height: 20px;
        text-align: center;
    `,
    styleToolbarColorButton: css`
        position: relative;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        width: 26px;
        height: 26px;
        padding: 0;
        border: 0;
        border-radius: ${token.borderRadiusSM}px;
        background: transparent;
        color: ${token.colorTextSecondary};
        font-size: 14px;
        line-height: 1;
        cursor: pointer;

        &:hover {
            color: ${token.colorPrimary};
            background: ${token.colorFillTertiary};
        }
    `,
    styleToolbarColorDot: css`
        position: absolute;
        left: 7px;
        right: 7px;
        bottom: 4px;
        height: 2px;
        border-radius: 2px;
        background: var(--mind-toolbar-color);
        box-shadow: 0 0 0 1px ${token.colorBgElevated};
    `,
    colorPanel: css`
        width: 188px;
        padding: 2px;
    `,
    colorPanelGrid: css`
        display: grid;
        grid-template-columns: repeat(6, 22px);
        gap: 7px;
        margin-bottom: 10px;
    `,
    colorSwatch: css`
        width: 22px;
        height: 22px;
        padding: 0;
        border: 1px solid ${token.colorBorderSecondary};
        border-radius: ${token.borderRadiusSM}px;
        background: var(--mind-swatch-color);
        cursor: pointer;

        &:hover {
            border-color: ${token.colorPrimary};
            box-shadow: 0 0 0 2px ${token.colorPrimaryBg};
        }
    `,
    colorPanelFooter: css`
        display: flex;
        align-items: center;
        gap: 8px;
    `,
    colorInput: css`
        width: 32px;
        height: 24px;
        padding: 0;
        border: 1px solid ${token.colorBorderSecondary};
        border-radius: ${token.borderRadiusSM}px;
        background: transparent;
        cursor: pointer;
    `,
    colorClearButton: css`
        flex: 1;
        height: 24px;
        padding: 0 8px;
        border: 0;
        border-radius: ${token.borderRadiusSM}px;
        background: ${token.colorFillTertiary};
        color: ${token.colorTextSecondary};
        font-size: 12px;
        cursor: pointer;

        &:hover {
            color: ${token.colorPrimary};
            background: ${token.colorPrimaryBg};
        }
    `,
    handle: css`
        width: 7px;
        height: 7px;
        border: 1px solid ${token.colorBgContainer};
        background: var(--mind-node-color);
        opacity: 0;
    `,
    contextMenu: css`
        position: fixed;
        z-index: 1000;
        min-width: 184px;
        padding: 4px;
        border: 1px solid ${token.colorBorderSecondary};
        border-radius: ${token.borderRadius}px;
        background: ${token.colorBgElevated};
        box-shadow: ${token.boxShadowSecondary};
    `,
    contextMenuItem: css`
        display: flex;
        align-items: center;
        gap: 8px;
        width: 100%;
        height: 32px;
        padding: 0 10px;
        border: 0;
        border-radius: ${token.borderRadiusSM}px;
        background: transparent;
        color: ${token.colorText};
        font-size: 13px;
        line-height: 32px;
        text-align: left;
        cursor: pointer;

        &:hover {
            background: ${token.colorFillTertiary};
        }

        &:disabled {
            color: ${token.colorTextDisabled};
            cursor: not-allowed;
            background: transparent;
        }
    `,
    contextMenuDanger: css`
        color: ${token.colorError};

        &:hover {
            background: ${token.colorErrorBg};
        }
    `,
    contextMenuDivider: css`
        height: 1px;
        margin: 4px;
        background: ${token.colorBorderSecondary};
    `,
}));

// ==================== 节点渲染 ====================

const MindNodeView: React.FC<MindNodeViewProps> = ({
    id,
    data,
    styles,
    readOnly,
    editingId,
    onStartEdit,
    onUpdateNode,
    onRemoveNode,
    onEditDone,
    onTextTargetChange,
    onRichSelectionChange,
    onShouldKeepEdit,
    onTextSelectionStart,
}) => {
    const titleRef = useRef<HTMLDivElement | null>(null);
    const descriptionRef = useRef<HTMLDivElement | null>(null);
    const editing = editingId === id;
    const titleTextStyle = getTitleTextStyle(data.style);
    const descriptionTextStyle = getDescriptionTextStyle(data.style);
    const titleHtml = data.titleHtml || textToRichHtml(data.title) || "";
    const descriptionHtml = data.descriptionHtml || textToRichHtml(data.description) || "";

    const setEditorContent = useCallback(() => {
        if (titleRef.current) {
            titleRef.current.innerHTML = titleHtml;
        }
        if (descriptionRef.current) {
            descriptionRef.current.innerHTML = descriptionHtml;
        }
    }, [descriptionHtml, titleHtml]);

    useEffect(() => {
        if (!editing) {
            return;
        }
        setEditorContent();
        window.requestAnimationFrame(() => {
            const input = titleRef.current;
            if (!input) {
                return;
            }
            input.focus();
            const range = document.createRange();
            range.selectNodeContents(input);
            range.collapse(false);
            const selection = window.getSelection();
            selection?.removeAllRanges();
            selection?.addRange(range);
            onTextTargetChange("title");
            onRichSelectionChange();
        });
    }, [editing, onRichSelectionChange, onTextTargetChange, setEditorContent]);

    const finishEdit = useCallback((cancel = false) => {
        if (!cancel) {
            const nextTitleHtml = normalizeRichHtml(titleRef.current?.innerHTML);
            const nextDescriptionHtml = normalizeRichHtml(descriptionRef.current?.innerHTML);
            const title = normalizeTitle(getRichText(nextTitleHtml));
            const description = normalizeDescription(getRichText(nextDescriptionHtml));
            if (!title && !data.root) {
                onRemoveNode(id);
                onEditDone();
                return;
            }
            const finalTitle = title || data.title;
            const titleHtmlValue = hasStyledRichHtml(nextTitleHtml) ? nextTitleHtml : undefined;
            const descriptionHtmlValue = hasStyledRichHtml(nextDescriptionHtml) ? nextDescriptionHtml : undefined;
            if (
                finalTitle !== data.title ||
                titleHtmlValue !== data.titleHtml ||
                description !== data.description ||
                descriptionHtmlValue !== data.descriptionHtml
            ) {
                onUpdateNode(id, {
                    title: finalTitle,
                    titleHtml: titleHtmlValue,
                    description,
                    descriptionHtml: descriptionHtmlValue,
                });
            }
        } else {
            setEditorContent();
        }
        onEditDone();
    }, [data.description, data.descriptionHtml, data.root, data.title, data.titleHtml, id, onEditDone, onRemoveNode, onUpdateNode, setEditorContent]);

    const stopInputEvent = useCallback((event: React.SyntheticEvent) => {
        event.stopPropagation();
    }, []);

    const focusEditorEnd = useCallback((target: TextTarget) => {
        const editor = target === "title" ? titleRef.current : descriptionRef.current;
        if (!editor) {
            return;
        }
        editor.focus();
        const range = document.createRange();
        range.selectNodeContents(editor);
        range.collapse(false);
        const selection = window.getSelection();
        selection?.removeAllRanges();
        selection?.addRange(range);
        onTextTargetChange(target);
        onRichSelectionChange();
    }, [onRichSelectionChange, onTextTargetChange]);

    const startInputSelection = useCallback((event: React.SyntheticEvent) => {
        onTextSelectionStart();
        event.stopPropagation();
    }, [onTextSelectionStart]);

    const onEditContentMouseDown = useCallback((event: React.MouseEvent<HTMLDivElement> | React.PointerEvent<HTMLDivElement>) => {
        startInputSelection(event);
        const element = event.target as HTMLElement | null;
        if (element?.closest("[data-mind-rich-editor='true']")) {
            return;
        }
        event.preventDefault();
        window.requestAnimationFrame(() => focusEditorEnd("description"));
    }, [focusEditorEnd, startInputSelection]);

    const startCardEdit = useCallback((event: React.SyntheticEvent) => {
        if (editing || isInputTarget(event.target) || readOnly) {
            return;
        }
        event.preventDefault();
        event.stopPropagation();
        onStartEdit(id);
    }, [editing, id, onStartEdit, readOnly]);

    const onEditorFocus = useCallback((target: TextTarget) => {
        onTextTargetChange(target);
        onRichSelectionChange();
    }, [onRichSelectionChange, onTextTargetChange]);

    const onEditorPaste = useCallback((event: React.ClipboardEvent<HTMLDivElement>) => {
        event.preventDefault();
        const selection = window.getSelection();
        if (selection && selection.rangeCount > 0) {
            const range = selection.getRangeAt(0);
            range.deleteContents();
            range.insertNode(createPlainRichFragment(event.clipboardData.getData("text/plain")));
        }
        onRichSelectionChange();
    }, [onRichSelectionChange]);

    const onTitleKeyDown = useCallback((event: React.KeyboardEvent<HTMLDivElement>) => {
        event.stopPropagation();
        if (event.key === "Enter") {
            event.preventDefault();
            finishEdit();
        }
        if (event.key === "Escape") {
            event.preventDefault();
            finishEdit(true);
        }
    }, [finishEdit]);

    const onDescriptionKeyDown = useCallback((event: React.KeyboardEvent<HTMLDivElement>) => {
        event.stopPropagation();
        if (event.key === "Enter" && (event.metaKey || event.ctrlKey)) {
            event.preventDefault();
            finishEdit();
        }
        if (event.key === "Escape") {
            event.preventDefault();
            finishEdit(true);
        }
    }, [finishEdit]);

    return (
        <div
            className={`${styles.node} ${data.root ? styles.nodeRoot : ""} ${editing ? styles.nodeEditing : ""}`}
            style={{
                "--mind-node-color": data.color,
                ...getNodeBoxStyle(data.style),
            } as React.CSSProperties}
            onMouseDownCapture={(event) => {
                if (event.detail >= 2) {
                    startCardEdit(event);
                }
            }}
            onDoubleClickCapture={startCardEdit}
        >
            <Handle id="left" type="target" position={Position.Left} className={styles.handle}/>
            {editing ? (
                <div
                    data-mind-edit-content="true"
                    className={`${styles.nodeEditContent} nodrag nopan nowheel`}
                    onBlur={(event) => {
                        const target = event.currentTarget;
                        window.setTimeout(() => {
                            const activeElement = document.activeElement as HTMLElement | null;
                            if (
                                !target.contains(activeElement) &&
                                !activeElement?.closest("[data-mind-toolbar='true'], [data-mind-color-picker='true'], .ant-color-picker, .ant-color-picker-dropdown, .ant-popover") &&
                                !onShouldKeepEdit()
                            ) {
                                finishEdit();
                            }
                        }, 0);
                    }}
                    onPointerDown={onEditContentMouseDown}
                    onMouseDown={onEditContentMouseDown}
                    onMouseUp={stopInputEvent}
                    onClick={stopInputEvent}
                    onDoubleClick={stopInputEvent}
                    onContextMenu={stopInputEvent}
                >
                    <div
                        ref={titleRef}
                        className={styles.nodeTitleInput}
                        style={titleTextStyle}
                        contentEditable
                        suppressContentEditableWarning
                        data-mind-rich-editor="true"
                        data-mind-editor-id={id}
                        data-mind-target="title"
                        data-placeholder="标题"
                        onFocus={() => onEditorFocus("title")}
                        onMouseUp={onRichSelectionChange}
                        onKeyUp={onRichSelectionChange}
                        onPaste={onEditorPaste}
                        onKeyDown={onTitleKeyDown}
                    />
                    <div
                        ref={descriptionRef}
                        className={styles.nodeDescriptionInput}
                        style={descriptionTextStyle}
                        contentEditable
                        suppressContentEditableWarning
                        data-mind-rich-editor="true"
                        data-mind-editor-id={id}
                        data-mind-target="description"
                        data-placeholder="添加描述"
                        onFocus={() => onEditorFocus("description")}
                        onMouseUp={onRichSelectionChange}
                        onKeyUp={onRichSelectionChange}
                        onPaste={onEditorPaste}
                        onKeyDown={onDescriptionKeyDown}
                    />
                </div>
            ) : (
                <div className={styles.nodeContent}>
                    <div className={styles.nodeTitle} style={titleTextStyle} dangerouslySetInnerHTML={{__html: titleHtml}}/>
                    {(data.description || data.descriptionHtml) && (
                        <div className={styles.nodeDescription} style={descriptionTextStyle} dangerouslySetInnerHTML={{__html: descriptionHtml}}/>
                    )}
                </div>
            )}
            <Handle id="right" type="source" position={Position.Right} className={styles.handle}/>
        </div>
    );
};

const MindNodeRenderer: React.FC<NodeProps<FlowNode>> = (props) => {
    const context = React.useContext(MindContext);
    if (!context) {
        return null;
    }
    return <MindNodeView {...props} {...context}/>;
};

const mindNodeTypes = {
    mind: MindNodeRenderer,
};

// ==================== 主组件 ====================

/** 思维导图组件，支持受控/非受控数据、富文本节点、快捷键、全屏和命令式 API。 */
const Mind = forwardRef<MindRef, MindProps>((props, ref) => {
    const {
        value,
        defaultValue,
        className,
        style,
        height = "100%",
        minHeight = 0,
        readOnly,
        showControls,
        defaultShowControls = true,
        showMiniMap,
        defaultShowMiniMap = true,
        newNodeTitle = "新节点",
        createId = getDefaultId,
        onChange,
        onDataChange,
        onControlsChange,
        onFullscreenChange,
        onMiniMapChange,
    } = props;
    const {styles} = useStyles();
    const {message} = App.useApp();
    const theme = useTheme();
    const rootRef = useRef<HTMLDivElement | null>(null);
    const [flow, setFlow] = useState<ReactFlowInstance<FlowNode, Edge> | null>(null);
    const [nodes, setNodes] = useState<FlowNode[]>(() => normalizeValue(value || defaultValue));
    const [editingId, setEditingId] = useState<string>();
    const [contextMenu, setContextMenu] = useState<ContextMenu | null>(null);
    const [textTarget, setTextTarget] = useState<TextTarget>("title");
    const [richSelectionActive, setRichSelectionActive] = useState(false);
    const [richSelectionStyle, setRichSelectionStyle] = useState<MindTextStyle>({});
    const [toolbarPopupOpen, setToolbarPopupOpen] = useState(false);
    const [toolbarColorOpen, setToolbarColorOpen] = useState<string>();
    const [fullscreen, setFullscreen] = useState(false);
    const [controlsVisible, setControlsVisible] = useState(defaultShowControls);
    const [miniMapVisible, setMiniMapVisible] = useState(defaultShowMiniMap);
    const initialValueRef = useRef<MindValue | undefined>(value || defaultValue);
    const contextMenuRef = useRef<HTMLDivElement | null>(null);
    const richSelectionRef = useRef<Range | null>(null);
    const toolbarPopupOpenRef = useRef(false);
    const toolbarMouseDownRef = useRef(false);
    const editOpeningRef = useRef(false);
    const editOpeningTimerRef = useRef<number | undefined>(undefined);
    const textSelectingRef = useRef(false);
    const textSelectingTimerRef = useRef<number | undefined>(undefined);
    const nodesRef = useRef(nodes);
    const selectedIdsRef = useRef<string[]>([]);
    const historyRef = useRef<FlowNode[][]>([]);
    const pendingEmitRef = useRef(false);
    const activeRef = useRef(false);
    const controlsVisibleRef = useRef(showControls ?? defaultShowControls);
    const fullscreenRef = useRef(false);
    const miniMapVisibleRef = useRef(showMiniMap ?? defaultShowMiniMap);
    const isDark = Boolean(theme?.isDarkTheme?.() || theme?.isDarkMode?.());
    const edges = useMemo(() => toEdges(nodes), [nodes]);
    const selectedNode = useMemo(() => nodes.find((node) => selectedIdsRef.current.includes(node.id)), [nodes]);
    const currentShowControls = showControls ?? controlsVisible;
    const currentShowMiniMap = showMiniMap ?? miniMapVisible;

    // ----- 基础状态同步 -----
    useEffect(() => {
        nodesRef.current = nodes;
    }, [nodes]);

    useEffect(() => {
        toolbarPopupOpenRef.current = toolbarPopupOpen;
    }, [toolbarPopupOpen]);

    useEffect(() => {
        fullscreenRef.current = fullscreen;
    }, [fullscreen]);

    useEffect(() => {
        controlsVisibleRef.current = currentShowControls;
    }, [currentShowControls]);

    useEffect(() => {
        miniMapVisibleRef.current = currentShowMiniMap;
    }, [currentShowMiniMap]);

    useEffect(() => {
        if (value) {
            if (!initialValueRef.current) {
                initialValueRef.current = value;
            }
            const nextNodes = normalizeValue(value);
            selectedIdsRef.current = [];
            setNodes(nextNodes);
        }
    }, [value]);

    useEffect(() => {
        if (!value && defaultValue && !initialValueRef.current) {
            initialValueRef.current = defaultValue;
            setNodes(normalizeValue(defaultValue));
        }
    }, [defaultValue, value]);

    // ----- 数据提交与历史快照 -----
    const emitChange = useCallback((nextNodes: FlowNode[]) => {
        const nextValue = toMindValue(nextNodes);
        const stats = getMindStats(nextValue);
        onChange?.(nextValue, stats);
        onDataChange?.(nextValue, stats);
    }, [onChange, onDataChange]);

    useEffect(() => {
        if (!pendingEmitRef.current) {
            return;
        }
        pendingEmitRef.current = false;
        emitChange(nodes);
    }, [emitChange, nodes]);

    const pushHistory = useCallback((items: FlowNode[]) => {
        historyRef.current = historyRef.current.concat([cloneFlowNodes(items)]).slice(-20);
    }, []);

    const updateNodes = useCallback((updater: FlowNode[] | ((items: FlowNode[]) => FlowNode[]), emit = true, history = emit) => {
        setNodes((items) => {
            const nextNodes = typeof updater === "function" ? updater(items) : updater;
            if (nextNodes === items) {
                return items;
            }
            if (history) {
                pushHistory(items);
            }
            if (emit) {
                pendingEmitRef.current = true;
            }
            return nextNodes;
        });
    }, [pushHistory]);

    const selectNodes = useCallback((ids?: string | string[]) => {
        const selectedIds = Array.isArray(ids) ? ids : (ids ? [ids] : []);
        selectedIdsRef.current = selectedIds;
        setNodes((items) => selectOnly(items, selectedIds));
    }, []);

    const fitView = useCallback(() => {
        flow?.fitView({padding: 0.22, duration: 260});
    }, [flow]);

    const zoomIn = useCallback(() => {
        flow?.zoomIn({duration: 180});
    }, [flow]);

    const zoomOut = useCallback(() => {
        flow?.zoomOut({duration: 180});
    }, [flow]);

    // ----- 节点操作 -----
    const updateNode = useCallback((id: string, data: MindNodeUpdate) => {
        if (readOnly) {
            return;
        }
        updateNodes((items) => {
            let changed = false;
            const nextItems = items.map((node) => {
                if (node.id !== id) {
                    return node;
                }
                const hasTitle = hasField(data, "title");
                const hasTitleHtml = hasField(data, "titleHtml");
                const hasDescription = hasField(data, "description");
                const hasDescriptionHtml = hasField(data, "descriptionHtml");
                const hasStyle = hasField(data, "style");
                const nextTitle = hasTitle ? normalizeTitle(data.title) : node.data.title;
                const title = node.data.root && !nextTitle ? node.data.title : nextTitle;
                const titleHtml = hasTitleHtml ? normalizeRichHtml(data.titleHtml) : node.data.titleHtml;
                const description = hasDescription ? normalizeDescription(data.description) : node.data.description;
                const descriptionHtml = hasDescriptionHtml ? normalizeRichHtml(data.descriptionHtml) : node.data.descriptionHtml;
                const color = data.color === null ? palette[0] : (data.color || node.data.color);
                const lineColor = data.lineColor === null ? color : (data.lineColor || node.data.lineColor || color);
                const nextStyle = hasStyle ? mergeMindNodeStyle(node.data.style, data.style) : node.data.style;
                if (
                    title === node.data.title &&
                    titleHtml === node.data.titleHtml &&
                    description === node.data.description &&
                    descriptionHtml === node.data.descriptionHtml &&
                    color === node.data.color &&
                    lineColor === node.data.lineColor &&
                    nextStyle === node.data.style
                ) {
                    return node;
                }
                changed = true;
                return {
                    ...node,
                    height: getMindNodeHeight({title, description, style: nextStyle}),
                    data: {
                        ...node.data,
                        title,
                        titleHtml,
                        description,
                        descriptionHtml,
                        color,
                        lineColor,
                        style: nextStyle,
                    },
                };
            });
            if (!changed) {
                return items;
            }
            return nextItems;
        });
    }, [readOnly, updateNodes]);

    const setLineColor = useCallback((id: string, color: string) => {
        updateNode(id, {lineColor: color});
    }, [updateNode]);

    const setNodeColor = useCallback((id: string, color: string) => {
        updateNode(id, {color, style: {color}});
    }, [updateNode]);

    const setNodeStyle = useCallback((id: string, style: MindNodeStyle) => {
        updateNode(id, {style});
    }, [updateNode]);

    const setBackgroundColor = useCallback((id: string, color: string) => {
        setNodeStyle(id, {background: color});
    }, [setNodeStyle]);

    const setBorderColor = useCallback((id: string, color: string) => {
        setNodeStyle(id, {borderColor: color});
    }, [setNodeStyle]);

    const clearNodeColor = useCallback((id: string) => {
        updateNode(id, {color: null, style: {color: undefined}});
    }, [updateNode]);

    const clearLineColor = useCallback((id: string) => {
        updateNode(id, {lineColor: null});
    }, [updateNode]);

    const clearBackgroundColor = useCallback((id: string) => {
        setNodeStyle(id, {background: undefined});
    }, [setNodeStyle]);

    const clearBorderColor = useCallback((id: string) => {
        setNodeStyle(id, {borderColor: undefined});
    }, [setNodeStyle]);

    const startEdit = useCallback((id?: string) => {
        if (readOnly) {
            return;
        }
        const targetId = id || selectedNode?.id;
        if (!targetId) {
            message.warning("请先选择一个节点");
            return;
        }
        window.clearTimeout(editOpeningTimerRef.current);
        editOpeningRef.current = true;
        editOpeningTimerRef.current = window.setTimeout(() => {
            editOpeningRef.current = false;
        }, 260);
        setEditingId(targetId);
        selectNodes(targetId);
    }, [message, readOnly, selectNodes, selectedNode?.id]);

    const addChild = useCallback((parentId?: string, position?: MindPoint) => {
        if (readOnly) {
            return;
        }
        const items = nodesRef.current;
        const parent = items.find((item) => item.id === (parentId || selectedNode?.id)) || items.find((item) => item.id === rootId);
        if (!parent) {
            return;
        }
        const id = createId();
        const nextNode = toFlowNode({
            id,
            title: newNodeTitle,
            parentId: parent.id,
            color: parent.data.color,
            lineColor: parent.data.lineColor || parent.data.color,
            position: {
                x: position?.x ?? parent.position.x + levelGap,
                y: position?.y ?? parent.position.y + rowGap,
            },
        });
        selectedIdsRef.current = [id];
        updateNodes((current) => selectOnly(position ? current.concat(nextNode) : layoutNodes(current.concat(nextNode)), id));
        window.requestAnimationFrame(fitView);
    }, [createId, fitView, newNodeTitle, readOnly, selectedNode?.id, updateNodes]);

    const addSibling = useCallback((nodeId?: string) => {
        if (readOnly) {
            return;
        }
        const items = nodesRef.current;
        const target = items.find((item) => item.id === (nodeId || selectedNode?.id));
        if (!target) {
            message.warning("请先选择一个节点");
            return;
        }
        if (target.id === rootId) {
            message.warning("根节点不能新增同级节点");
            return;
        }
        const id = createId();
        const nextNode = toFlowNode({
            id,
            title: newNodeTitle,
            parentId: target.data.parentId,
            color: target.data.color,
            lineColor: target.data.lineColor || target.data.color,
            position: {
                x: target.position.x,
                y: target.position.y + rowGap,
            },
        });
        selectedIdsRef.current = [id];
        updateNodes((current) => selectOnly(layoutNodes(current.concat(nextNode)), id));
        window.requestAnimationFrame(fitView);
    }, [createId, fitView, message, newNodeTitle, readOnly, selectedNode?.id, updateNodes]);

    const removeNodes = useCallback((ids: string[]) => {
        if (readOnly) {
            return;
        }
        if (ids.length <= 0) {
            message.warning("请先选择要删除的节点");
            return;
        }
        const targetIds = ids.filter((id) => id !== rootId);
        if (targetIds.length <= 0) {
            message.warning("根节点不能删除");
            return;
        }
        const removeIds = getDescendantIds(nodesRef.current, targetIds);
        selectedIdsRef.current = [];
        setEditingId(undefined);
        updateNodes((items) => layoutNodes(items.filter((node) => !removeIds.has(node.id))));
    }, [message, readOnly, updateNodes]);

    const removeSelected = useCallback(() => {
        removeNodes([...selectedIdsRef.current]);
    }, [removeNodes]);

    const arrangeMind = useCallback(() => {
        updateNodes((items) => layoutNodes(items));
        window.requestAnimationFrame(fitView);
    }, [fitView, updateNodes]);

    const resetMind = useCallback(() => {
        const nextNodes = normalizeValue(initialValueRef.current);
        selectedIdsRef.current = [];
        updateNodes(nextNodes);
        window.requestAnimationFrame(fitView);
    }, [fitView, updateNodes]);

    const undoMind = useCallback(() => {
        const snapshot = historyRef.current.pop();
        if (!snapshot) {
            return false;
        }
        selectedIdsRef.current = [];
        setEditingId(undefined);
        updateNodes(selectOnly(cloneFlowNodes(snapshot)), true, false);
        window.requestAnimationFrame(fitView);
        return true;
    }, [fitView, updateNodes]);

    // ----- 视图与全屏控制 -----
    const openFullscreen = useCallback(async () => {
        const element = rootRef.current;
        if (!element || document.fullscreenElement === element) {
            return;
        }
        await element.requestFullscreen?.();
    }, []);

    const exitFullscreen = useCallback(async () => {
        if (document.fullscreenElement === rootRef.current) {
            await document.exitFullscreen?.();
        }
    }, []);

    const toggleFullscreen = useCallback(async () => {
        if (document.fullscreenElement === rootRef.current) {
            await exitFullscreen();
            return;
        }
        await openFullscreen();
    }, [exitFullscreen, openFullscreen]);

    useEffect(() => {
        const handleFullscreenChange = () => {
            const nextFullscreen = document.fullscreenElement === rootRef.current;
            setFullscreen(nextFullscreen);
            onFullscreenChange?.(nextFullscreen);
            window.requestAnimationFrame(fitView);
        };
        document.addEventListener("fullscreenchange", handleFullscreenChange);
        return () => {
            document.removeEventListener("fullscreenchange", handleFullscreenChange);
        };
    }, [fitView, onFullscreenChange]);

    // ----- 左下角控制器 -----
    const changeControlsVisible = useCallback((visible: boolean) => {
        controlsVisibleRef.current = visible;
        if (showControls === undefined) {
            setControlsVisible(visible);
        }
        onControlsChange?.(visible);
    }, [onControlsChange, showControls]);

    const showControlsAction = useCallback(() => {
        changeControlsVisible(true);
    }, [changeControlsVisible]);

    const hideControlsAction = useCallback(() => {
        changeControlsVisible(false);
    }, [changeControlsVisible]);

    const toggleControls = useCallback(() => {
        changeControlsVisible(!controlsVisibleRef.current);
    }, [changeControlsVisible]);

    // ----- 右下角缩略图 -----
    const changeMiniMapVisible = useCallback((visible: boolean) => {
        miniMapVisibleRef.current = visible;
        if (showMiniMap === undefined) {
            setMiniMapVisible(visible);
        }
        onMiniMapChange?.(visible);
    }, [onMiniMapChange, showMiniMap]);

    const showMiniMapAction = useCallback(() => {
        changeMiniMapVisible(true);
    }, [changeMiniMapVisible]);

    const hideMiniMapAction = useCallback(() => {
        changeMiniMapVisible(false);
    }, [changeMiniMapVisible]);

    const toggleMiniMap = useCallback(() => {
        changeMiniMapVisible(!miniMapVisibleRef.current);
    }, [changeMiniMapVisible]);

    // ----- 编辑保护、右键菜单与快捷键 -----
    const closeContextMenu = useCallback(() => {
        setContextMenu(null);
    }, []);

    const finishTextSelectionGuard = useCallback(() => {
        window.clearTimeout(textSelectingTimerRef.current);
        textSelectingTimerRef.current = window.setTimeout(() => {
            textSelectingRef.current = false;
        }, 220);
    }, []);

    const startTextSelectionGuard = useCallback(() => {
        window.clearTimeout(textSelectingTimerRef.current);
        textSelectingRef.current = true;
        document.removeEventListener("pointerup", finishTextSelectionGuard, true);
        document.removeEventListener("mouseup", finishTextSelectionGuard, true);
        document.addEventListener("pointerup", finishTextSelectionGuard, {capture: true, once: true});
        document.addEventListener("mouseup", finishTextSelectionGuard, {capture: true, once: true});
    }, [finishTextSelectionGuard]);

    useEffect(() => {
        return () => {
            window.clearTimeout(editOpeningTimerRef.current);
            window.clearTimeout(textSelectingTimerRef.current);
            document.removeEventListener("pointerup", finishTextSelectionGuard, true);
            document.removeEventListener("mouseup", finishTextSelectionGuard, true);
        };
    }, [finishTextSelectionGuard]);

    const startToolbarInteraction = useCallback(() => {
        toolbarMouseDownRef.current = true;
        window.setTimeout(() => {
            toolbarMouseDownRef.current = false;
        }, 180);
    }, []);

    const changeToolbarPopupOpen = useCallback((open: boolean, key?: string, keepTextToolbar = false) => {
        toolbarPopupOpenRef.current = open;
        setToolbarColorOpen(open ? key : undefined);
        if (keepTextToolbar) {
            setToolbarPopupOpen(open);
        } else if (!open) {
            setToolbarPopupOpen(false);
        }
    }, []);

    const shouldKeepEdit = useCallback(() => {
        return editOpeningRef.current ||
            textSelectingRef.current ||
            toolbarMouseDownRef.current ||
            toolbarPopupOpenRef.current ||
            Boolean(document.querySelector(".ant-popover"));
    }, []);

    const getRichSelectionEditor = useCallback((node?: globalThis.Node | null) => {
        const element = node instanceof HTMLElement ? node : node?.parentElement;
        return element?.closest?.("[data-mind-rich-editor='true']") as HTMLElement | null;
    }, []);

    const getRichEditor = useCallback((nodeId: string, target: TextTarget) => {
        return Array.from(rootRef.current?.querySelectorAll("[data-mind-rich-editor='true']") || [])
            .find((item) => {
                const element = item as HTMLElement;
                return element.dataset.mindEditorId === nodeId && element.dataset.mindTarget === target;
            }) as HTMLElement | undefined || null;
    }, []);

    const saveRichSelection = useCallback(() => {
        const selection = window.getSelection();
        if (!selection || selection.rangeCount <= 0) {
            if (!toolbarPopupOpenRef.current) {
                setRichSelectionActive(false);
            }
            return;
        }
        const range = selection.getRangeAt(0);
        const editor = getRichSelectionEditor(range.commonAncestorContainer);
        if (!editor || !rootRef.current?.contains(editor)) {
            if (!toolbarPopupOpenRef.current) {
                setRichSelectionActive(false);
            }
            return;
        }
        richSelectionRef.current = range.cloneRange();
        const target = editor.dataset.mindTarget as TextTarget | undefined;
        if (target) {
            setTextTarget(target);
        }
        setRichSelectionActive(!range.collapsed && Boolean(range.toString().trim()));
        setRichSelectionStyle(getRichRangeStyle(range, editor));
    }, [getRichSelectionEditor]);

    useEffect(() => {
        if (!editingId) {
            return;
        }
        const handleSelectionChange = () => {
            saveRichSelection();
        };
        document.addEventListener("selectionchange", handleSelectionChange);
        return () => {
            document.removeEventListener("selectionchange", handleSelectionChange);
        };
    }, [editingId, saveRichSelection]);

    const restoreRichSelection = useCallback(() => {
        const range = richSelectionRef.current;
        if (!range) {
            return false;
        }
        const editor = getRichSelectionEditor(range.commonAncestorContainer);
        if (!editor || !rootRef.current?.contains(editor)) {
            return false;
        }
        editor.focus({preventScroll: true});
        const selection = window.getSelection();
        selection?.removeAllRanges();
        selection?.addRange(range);
        return true;
    }, [getRichSelectionEditor]);

    const selectRichRange = useCallback((range: Range, editor: HTMLElement) => {
        editor.focus({preventScroll: true});
        const selection = window.getSelection();
        selection?.removeAllRanges();
        selection?.addRange(range);
        richSelectionRef.current = range.cloneRange();
        const target = editor.dataset.mindTarget as TextTarget | undefined;
        if (target) {
            setTextTarget(target);
        }
        setRichSelectionActive(!range.collapsed && Boolean(range.toString().trim()));
        setRichSelectionStyle(getRichRangeStyle(range, editor));
    }, []);

    const applyRichSelectionStyle = useCallback((style: React.CSSProperties) => {
        if (!editingId) {
            return false;
        }
        if (!restoreRichSelection()) {
            return false;
        }
        const selection = window.getSelection();
        if (!selection || selection.rangeCount <= 0) {
            return false;
        }
        const range = selection.getRangeAt(0);
        const editor = getRichSelectionEditor(range.commonAncestorContainer);
        if (!editor || range.collapsed) {
            return false;
        }
        const nextStyle = getAppliedRichStyle(style, editor);
        const span = document.createElement("span");
        setElementStyle(span, nextStyle);
        const fragment = range.extractContents();
        clearRichFragmentStyle(fragment, nextStyle);
        span.appendChild(fragment);
        range.insertNode(span);
        editor.normalize();
        const nextRange = document.createRange();
        nextRange.selectNodeContents(span);
        selectRichRange(nextRange, editor);
        return true;
    }, [editingId, getRichSelectionEditor, restoreRichSelection, selectRichRange]);

    const runRichCommand = useCallback((command: "underline" | "strikeThrough") => {
        if (!editingId) {
            return false;
        }
        if (!restoreRichSelection()) {
            return false;
        }
        const selection = window.getSelection();
        if (!selection || selection.rangeCount <= 0) {
            return false;
        }
        const range = selection.getRangeAt(0);
        const editor = getRichSelectionEditor(range.commonAncestorContainer);
        if (!editor || range.collapsed) {
            return false;
        }
        document.execCommand(command, false);
        editor.normalize();
        if (selection.rangeCount > 0) {
            selectRichRange(selection.getRangeAt(0).cloneRange(), editor);
        }
        return true;
    }, [editingId, getRichSelectionEditor, restoreRichSelection, selectRichRange]);

    const clearRichSelectionStyle = useCallback(() => {
        if (!editingId) {
            return false;
        }
        if (!restoreRichSelection()) {
            return false;
        }
        const selection = window.getSelection();
        if (!selection || selection.rangeCount <= 0) {
            return false;
        }
        const range = selection.getRangeAt(0);
        const editor = getRichSelectionEditor(range.commonAncestorContainer);
        if (!editor || range.collapsed) {
            return false;
        }
        const span = document.createElement("span");
        setElementStyle(span, getRichEditorBaseStyle(editor));
        span.appendChild(createPlainRichFragment(range.toString()));
        range.deleteContents();
        range.insertNode(span);
        editor.normalize();
        const nextRange = document.createRange();
        nextRange.selectNodeContents(span);
        selectRichRange(nextRange, editor);
        return true;
    }, [editingId, getRichSelectionEditor, restoreRichSelection, selectRichRange]);

    const runRichFontSize = useCallback((fontSize: number) => {
        return applyRichSelectionStyle({fontSize});
    }, [applyRichSelectionStyle]);

    const flushEditingNode = useCallback(() => {
        if (!editingId) {
            return;
        }
        const node = nodesRef.current.find((item) => item.id === editingId);
        if (!node) {
            setEditingId(undefined);
            return;
        }
        const nextTitleHtml = normalizeRichHtml(getRichEditor(editingId, "title")?.innerHTML);
        const nextDescriptionHtml = normalizeRichHtml(getRichEditor(editingId, "description")?.innerHTML);
        const title = normalizeTitle(getRichText(nextTitleHtml));
        const description = normalizeDescription(getRichText(nextDescriptionHtml));
        if (!title && !node.data.root) {
            removeNodes([editingId]);
            return;
        }
        updateNode(editingId, {
            title: title || node.data.title,
            titleHtml: hasStyledRichHtml(nextTitleHtml) ? nextTitleHtml : undefined,
            description,
            descriptionHtml: hasStyledRichHtml(nextDescriptionHtml) ? nextDescriptionHtml : undefined,
        });
        setEditingId(undefined);
    }, [editingId, getRichEditor, removeNodes, updateNode]);

    const focusEditor = useCallback(() => {
        activeRef.current = true;
        rootRef.current?.focus({preventScroll: true});
    }, []);

    useLayoutEffect(() => {
        if (!contextMenu) {
            return;
        }
        const element = contextMenuRef.current;
        if (!element) {
            return;
        }
        const rect = element.getBoundingClientRect();
        const margin = 8;
        const maxX = Math.max(margin, window.innerWidth - rect.width - margin);
        const maxY = Math.max(margin, window.innerHeight - rect.height - margin);
        const x = Math.min(Math.max(margin, contextMenu.x), maxX);
        const y = Math.min(Math.max(margin, contextMenu.y), maxY);
        if (x !== contextMenu.x || y !== contextMenu.y) {
            setContextMenu((item) => item ? {...item, x, y} : item);
        }
    }, [contextMenu]);

    useEffect(() => {
        const handleClick = (event: MouseEvent) => {
            activeRef.current = Boolean(rootRef.current?.contains(event.target as globalThis.Node));
            closeContextMenu();
        };
        const handleKeyDown = (event: KeyboardEvent) => {
            const target = event.target as HTMLElement | null;
            const isInput = isInputTarget(target);
            const selectedIds = [...selectedIdsRef.current];
            const active = activeRef.current || Boolean(rootRef.current?.contains(document.activeElement));
            if (event.key === "Escape") {
                closeContextMenu();
            }
            if (!active || isInput) {
                return;
            }
            if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "z" && !event.shiftKey) {
                event.preventDefault();
                closeContextMenu();
                undoMind();
                return;
            }
            if (!readOnly && selectedIds.length > 0 && (event.key === "Delete" || event.key === "Backspace")) {
                event.preventDefault();
                closeContextMenu();
                removeNodes(selectedIds);
            }
        };
        document.addEventListener("click", handleClick);
        document.addEventListener("keydown", handleKeyDown);
        return () => {
            document.removeEventListener("click", handleClick);
            document.removeEventListener("keydown", handleKeyDown);
        };
    }, [closeContextMenu, readOnly, removeNodes, undoMind]);

    // ----- React Flow 节点上下文与事件 -----
    const nodeContext = useMemo<MindContextValue>(() => ({
        styles,
        readOnly,
        editingId,
        onStartEdit: startEdit,
        onUpdateNode: updateNode,
        onRemoveNode: (id) => removeNodes([id]),
        onEditDone: () => setEditingId(undefined),
        onTextTargetChange: setTextTarget,
        onRichSelectionChange: saveRichSelection,
        onShouldKeepEdit: shouldKeepEdit,
        onTextSelectionStart: startTextSelectionGuard,
    }), [editingId, readOnly, removeNodes, saveRichSelection, shouldKeepEdit, startEdit, startTextSelectionGuard, styles, updateNode]);

    const onNodesChange = useCallback((changes: NodeChange<FlowNode>[]) => {
        const shouldEmit = changes.some((change) => change.type !== "select" && change.type !== "dimensions");
        const shouldHistory = shouldEmit && changes.every((change) => change.type !== "position");
        updateNodes((items) => {
            const nextNodes = applyNodeChanges(changes, items) as FlowNode[];
            if (changes.some((change) => change.type === "select")) {
                selectedIdsRef.current = nextNodes.filter((node) => node.selected).map((node) => node.id);
            }
            return nextNodes;
        }, shouldEmit, shouldHistory);
    }, [updateNodes]);

    const openNodeContextMenu = useCallback((event: React.MouseEvent, node: FlowNode) => {
        if (readOnly) {
            return;
        }
        event.preventDefault();
        event.stopPropagation();
        if (editingId && editingId !== node.id) {
            flushEditingNode();
        }
        focusEditor();
        selectNodes(node.id);
        setContextMenu({x: event.clientX, y: event.clientY, nodeId: node.id});
    }, [editingId, flushEditingNode, focusEditor, readOnly, selectNodes]);

    const openPaneContextMenu = useCallback((event: MouseEvent | React.MouseEvent) => {
        if (readOnly) {
            return;
        }
        event.preventDefault();
        event.stopPropagation();
        flushEditingNode();
        focusEditor();
        const position = flow?.screenToFlowPosition({x: event.clientX, y: event.clientY});
        selectNodes();
        setContextMenu({x: event.clientX, y: event.clientY, position});
    }, [flow, flushEditingNode, focusEditor, readOnly, selectNodes]);

    const onContextAction = useCallback((key: string) => {
        const nodeId = contextMenu?.nodeId;
        if (key === "child") {
            addChild(nodeId || rootId, nodeId ? undefined : contextMenu?.position);
        } else if (key === "sibling") {
            addSibling(nodeId);
        } else if (key === "edit") {
            startEdit(nodeId);
        } else if (key === "delete" && nodeId) {
            removeNodes([nodeId]);
        } else if (key === "undo") {
            undoMind();
        }
        closeContextMenu();
    }, [addChild, addSibling, closeContextMenu, contextMenu, removeNodes, startEdit, undoMind]);

    // ----- 工具栏状态与动作 -----
    const onToolbarColor = useCallback((type: "background" | "border" | "line", color: string) => {
        const nodeId = selectedNode?.id;
        if (!nodeId) {
            return;
        }
        if (type === "background") {
            setBackgroundColor(nodeId, color);
        } else if (type === "border") {
            setBorderColor(nodeId, color);
        } else {
            setLineColor(nodeId, color);
        }
    }, [selectedNode?.id, setBackgroundColor, setBorderColor, setLineColor]);

    const onToolbarClearColor = useCallback((type: "background" | "border" | "line") => {
        const nodeId = selectedNode?.id;
        if (!nodeId) {
            return;
        }
        if (type === "background") {
            clearBackgroundColor(nodeId);
        } else if (type === "border") {
            clearBorderColor(nodeId);
        } else {
            clearLineColor(nodeId);
        }
    }, [clearBackgroundColor, clearBorderColor, clearLineColor, selectedNode?.id]);

    const updateSelectedTextStyle = useCallback((style?: MindTextStyle) => {
        const nodeId = selectedNode?.id;
        if (!nodeId) {
            return;
        }
        setNodeStyle(nodeId, {[getTextStyleKey(textTarget)]: style} as MindNodeStyle);
    }, [selectedNode?.id, setNodeStyle, textTarget]);

    const selectedStyle = selectedNode?.data.style || {};
    const toolbarTextStyle = editingId ? {
        ...getToolbarTextStyle(selectedStyle, textTarget),
        ...richSelectionStyle,
    } : getToolbarTextStyle(selectedStyle, textTarget);
    const toolbarTextColor = toolbarTextStyle.color || (textTarget === "title" ? (selectedNode?.data.color || palette[0]) : (isDark ? "#d9d9d9" : "#595959"));
    const toolbarTextFontSize = toolbarTextStyle.fontSize || (textTarget === "title" ? defaultFontSize : defaultDescriptionFontSize);
    const toolbarTextBold = toolbarTextStyle.fontWeight !== "normal" && (toolbarTextStyle.fontWeight === "bold" || Number(toolbarTextStyle.fontWeight) >= 600);
    const toolbarTextItalic = toolbarTextStyle.fontStyle === "italic";
    const toolbarTextUnderline = toolbarTextStyle.textDecoration !== "none" && hasTextDecoration(toolbarTextStyle.textDecoration, "underline");
    const toolbarTextStrike = toolbarTextStyle.textDecoration !== "none" && hasTextDecoration(toolbarTextStyle.textDecoration, "line-through");

    const changeSelectedTextFontSize = useCallback((step: number) => {
        const fontSize = Math.min(maxTextFontSize, Math.max(minTextFontSize, toolbarTextFontSize + step));
        if (editingId) {
            runRichFontSize(fontSize);
            return;
        }
        if (!runRichFontSize(fontSize)) {
            updateSelectedTextStyle({fontSize});
        }
    }, [editingId, runRichFontSize, toolbarTextFontSize, updateSelectedTextStyle]);

    const setSelectedTextColor = useCallback((color: string) => {
        if (editingId) {
            applyRichSelectionStyle({color});
            return;
        }
        if (!applyRichSelectionStyle({color})) {
            updateSelectedTextStyle({color});
        }
    }, [applyRichSelectionStyle, editingId, updateSelectedTextStyle]);

    const clearSelectedTextColor = useCallback(() => {
        if (editingId) {
            applyRichSelectionStyle({color: "inherit"});
            return;
        }
        updateSelectedTextStyle({color: undefined});
    }, [applyRichSelectionStyle, editingId, updateSelectedTextStyle]);

    const clearSelectedTextStyle = useCallback(() => {
        if (clearRichSelectionStyle()) {
            return;
        }
        const nodeId = selectedNode?.id;
        if (!nodeId) {
            return;
        }
        setNodeStyle(nodeId, textTarget === "title" ? {
            color: undefined,
            fontSize: undefined,
            titleStyle: undefined,
        } : {
            descriptionStyle: undefined,
        });
    }, [clearRichSelectionStyle, selectedNode?.id, setNodeStyle, textTarget]);

    const toggleSelectedTextBold = useCallback(() => {
        const fontWeight = toolbarTextBold ? "normal" : boldFontWeight;
        if (editingId) {
            applyRichSelectionStyle({fontWeight});
            return;
        }
        if (!applyRichSelectionStyle({fontWeight})) {
            updateSelectedTextStyle({fontWeight: toolbarTextBold ? undefined : boldFontWeight});
        }
    }, [applyRichSelectionStyle, editingId, toolbarTextBold, updateSelectedTextStyle]);

    const toggleSelectedTextItalic = useCallback(() => {
        const fontStyle = toolbarTextItalic ? "normal" : "italic";
        if (editingId) {
            applyRichSelectionStyle({fontStyle});
            return;
        }
        if (!applyRichSelectionStyle({fontStyle})) {
            updateSelectedTextStyle({fontStyle: toolbarTextItalic ? undefined : "italic"});
        }
    }, [applyRichSelectionStyle, editingId, toolbarTextItalic, updateSelectedTextStyle]);

    const toggleSelectedTextDecoration = useCallback((decoration: string) => {
        const textDecoration = toggleTextDecoration(toolbarTextStyle.textDecoration, decoration);
        if (editingId) {
            if (runRichCommand(decoration === "underline" ? "underline" : "strikeThrough")) {
                return;
            }
            applyRichSelectionStyle({textDecoration: textDecoration || "none"});
            return;
        }
        if (!applyRichSelectionStyle({textDecoration: textDecoration || "none"})) {
            updateSelectedTextStyle({textDecoration});
        }
    }, [applyRichSelectionStyle, editingId, runRichCommand, toolbarTextStyle.textDecoration, updateSelectedTextStyle]);

    const changeSelectedRadius = useCallback((step: number) => {
        const nodeId = selectedNode?.id;
        if (!nodeId) {
            return;
        }
        const borderRadius = Math.min(18, Math.max(0, (selectedNode.data.style?.borderRadius ?? defaultBorderRadius) + step));
        setNodeStyle(nodeId, {borderRadius});
    }, [selectedNode, setNodeStyle]);

    const toggleSelectedShadow = useCallback(() => {
        const nodeId = selectedNode?.id;
        if (!nodeId) {
            return;
        }
        setNodeStyle(nodeId, {
            boxShadow: selectedNode.data.style?.boxShadow ? undefined : nodeShadow,
        });
    }, [selectedNode, setNodeStyle]);

    const toolbarFillColor = selectedStyle.background || (isDark ? "#1f1f1f" : "#ffffff");
    const toolbarBorderColor = selectedStyle.borderColor || (isDark ? "#434343" : "#d9d9d9");
    const toolbarLineColor = selectedNode?.data.lineColor || selectedNode?.data.color || palette[0];
    const toolbarRadius = selectedStyle.borderRadius ?? defaultBorderRadius;
    const toolbarShadow = Boolean(selectedStyle.boxShadow);
    const showTextToolbar = Boolean(selectedNode && !readOnly && editingId && (richSelectionActive || toolbarPopupOpen));
    const showCardToolbar = Boolean(selectedNode && !readOnly && !showTextToolbar);

    // ----- 工具栏渲染工具 -----
    const renderToolbarColor = useCallback((
        label: string,
        icon: React.ReactNode,
        value: string,
        onChange: (color: string) => void,
        onClear: () => void,
        keepTextToolbar = false,
    ) => {
        const open = toolbarColorOpen === label;
        const changeOpen = (nextOpen: boolean) => {
            startToolbarInteraction();
            if (nextOpen) {
                saveRichSelection();
            }
            changeToolbarPopupOpen(nextOpen, label, keepTextToolbar);
        };
        const selectColor = (color: string, close = true) => {
            startToolbarInteraction();
            onChange(color);
            if (close) {
                changeToolbarPopupOpen(false, label, keepTextToolbar);
            }
        };
        const content = (
            <div
                data-mind-color-picker="true"
                className={styles.colorPanel}
                onMouseDown={(event) => event.stopPropagation()}
                onClick={(event) => event.stopPropagation()}
            >
                <div className={styles.colorPanelGrid}>
                    {stylePalette.map((color) => (
                        <button
                            key={color}
                            type="button"
                            className={styles.colorSwatch}
                            style={{"--mind-swatch-color": color} as React.CSSProperties}
                            onMouseDown={(event) => {
                                event.preventDefault();
                                event.stopPropagation();
                            }}
                            onClick={() => selectColor(color)}
                        />
                    ))}
                </div>
                <div className={styles.colorPanelFooter}>
                    <input
                        type="color"
                        className={styles.colorInput}
                        value={toHexColor(value)}
                        onMouseDown={(event) => {
                            startToolbarInteraction();
                            event.stopPropagation();
                        }}
                        onChange={(event) => selectColor(event.target.value, false)}
                    />
                    <button
                        type="button"
                        className={styles.colorClearButton}
                        onMouseDown={(event) => {
                            event.preventDefault();
                            event.stopPropagation();
                        }}
                        onClick={() => {
                            startToolbarInteraction();
                            onClear();
                            changeToolbarPopupOpen(false, label, keepTextToolbar);
                        }}
                    >
                        清除颜色
                    </button>
                </div>
            </div>
        );
        return (
            <Popover
                open={open}
                trigger="click"
                placement="bottom"
                content={content}
                destroyOnHidden
                getPopupContainer={() => rootRef.current || document.body}
                onOpenChange={changeOpen}
            >
                <span data-mind-color-picker="true" className={styles.styleToolbarTooltip}>
                    <Tooltip title={label}>
                        <button
                            type="button"
                            tabIndex={-1}
                            className={styles.styleToolbarColorButton}
                            style={{"--mind-toolbar-color": value} as React.CSSProperties}
                            onMouseDown={() => {
                                startToolbarInteraction();
                                saveRichSelection();
                            }}
                        >
                            {icon}
                            <span className={styles.styleToolbarColorDot}/>
                        </button>
                    </Tooltip>
                </span>
            </Popover>
        );
    }, [
        changeToolbarPopupOpen,
        saveRichSelection,
        startToolbarInteraction,
        styles.colorClearButton,
        styles.colorInput,
        styles.colorPanel,
        styles.colorPanelFooter,
        styles.colorPanelGrid,
        styles.colorSwatch,
        styles.styleToolbarColorButton,
        styles.styleToolbarColorDot,
        styles.styleToolbarTooltip,
        toolbarColorOpen,
    ]);

    const renderToolbarButton = useCallback((
        title: string,
        icon: React.ReactNode,
        onClick: () => void,
        disabled = false,
        active = false,
    ) => (
        <Tooltip title={title}>
            <span className={styles.styleToolbarTooltip}>
                <button
                    type="button"
                    aria-disabled={disabled}
                    tabIndex={-1}
                    className={`${styles.styleToolbarButton} ${active ? styles.styleToolbarButtonActive : ""} ${disabled ? styles.styleToolbarButtonDisabled : ""}`}
                    onClick={() => {
                        if (!disabled) {
                            onClick();
                        }
                    }}
                >
                    {icon}
                </button>
            </span>
        </Tooltip>
    ), [styles.styleToolbarButton, styles.styleToolbarButtonActive, styles.styleToolbarButtonDisabled, styles.styleToolbarTooltip]);

    // ----- 命令式 API -----
    useImperativeHandle(ref, () => ({
        getValue: () => toMindValue(nodesRef.current),
        getData: () => toMindValue(nodesRef.current),
        getStats: () => getFlowStats(nodesRef.current),
        addChild,
        addSibling,
        edit: startEdit,
        updateNode,
        setNodeColor,
        setLineColor,
        setBackgroundColor,
        setBorderColor,
        setNodeStyle,
        clearNodeColor,
        clearLineColor,
        clearBackgroundColor,
        clearBorderColor,
        undo: undoMind,
        canUndo: () => historyRef.current.length > 0,
        removeSelected,
        layout: arrangeMind,
        arrange: arrangeMind,
        fitView,
        zoomIn,
        zoomOut,
        reset: resetMind,
        isFullscreen: () => fullscreenRef.current,
        openFullscreen,
        exitFullscreen,
        toggleFullscreen,
        isControlsVisible: () => controlsVisibleRef.current,
        showControls: showControlsAction,
        hideControls: hideControlsAction,
        toggleControls,
        isMiniMapVisible: () => miniMapVisibleRef.current,
        showMiniMap: showMiniMapAction,
        hideMiniMap: hideMiniMapAction,
        toggleMiniMap,
    }), [
        addChild,
        addSibling,
        arrangeMind,
        clearBackgroundColor,
        clearBorderColor,
        clearLineColor,
        clearNodeColor,
        exitFullscreen,
        fitView,
        hideControlsAction,
        hideMiniMapAction,
        openFullscreen,
        removeSelected,
        resetMind,
        setBackgroundColor,
        setBorderColor,
        setLineColor,
        setNodeColor,
        setNodeStyle,
        showControlsAction,
        showMiniMapAction,
        startEdit,
        toggleControls,
        toggleMiniMap,
        toggleFullscreen,
        undoMind,
        updateNode,
        zoomIn,
        zoomOut,
    ]);

    // ----- 渲染 -----
    return (
        <div
            ref={rootRef}
            className={`${styles.editor} ${fullscreen ? styles.fullscreen : ""} ${className || ""}`}
            tabIndex={0}
            onMouseDownCapture={(event) => {
                if (!isInputTarget(event.target) && !isMindInteractiveTarget(event.target)) {
                    focusEditor();
                }
            }}
            style={{
                ...style,
                height: fullscreen ? "100vh" : height,
                minHeight: fullscreen ? 0 : minHeight,
            }}
        >
            <MindContext.Provider value={nodeContext}>
                <ReactFlowProvider>
                    <ReactFlow
                        nodes={nodes}
                        edges={edges}
                        nodeTypes={mindNodeTypes}
                        onInit={setFlow}
                        onNodesChange={onNodesChange}
                        onSelectionChange={({nodes: selectedNodes}) => {
                            selectedIdsRef.current = selectedNodes.map((node) => node.id);
                        }}
                        onNodeDragStart={() => {
                            if (!readOnly) {
                                pushHistory(nodesRef.current);
                            }
                        }}
                        onNodeClick={(_, node) => {
                            if (shouldKeepEdit()) {
                                return;
                            }
                            if (editingId && editingId !== node.id) {
                                flushEditingNode();
                            }
                            focusEditor();
                            selectNodes(node.id);
                        }}
                        onNodeDoubleClick={(event, node) => {
                            focusEditor();
                            event.preventDefault();
                            event.stopPropagation();
                            startEdit(node.id);
                        }}
                        onPaneClick={() => {
                            if (shouldKeepEdit()) {
                                return;
                            }
                            flushEditingNode();
                            focusEditor();
                            selectNodes();
                        }}
                        onNodeContextMenu={openNodeContextMenu}
                        onPaneContextMenu={openPaneContextMenu}
                        fitView
                        fitViewOptions={{padding: 0.22}}
                        defaultEdgeOptions={{type: "default"}}
                        nodesDraggable={!readOnly}
                        nodesConnectable={false}
                        deleteKeyCode={null}
                        elementsSelectable={!readOnly}
                        panOnDrag
                        selectionOnDrag={!readOnly}
                        zoomOnDoubleClick={false}
                        snapToGrid
                        snapGrid={[8, 8]}
                        minZoom={0.25}
                        maxZoom={2}
                        proOptions={{hideAttribution: true}}
                    >
                        <Background
                            variant={BackgroundVariant.Dots}
                            gap={24}
                            size={1}
                            color={isDark ? "rgba(255,255,255,0.11)" : "rgba(15,23,42,0.12)"}
                        />
                        {currentShowControls && <Controls showInteractive={false}/>}
                        {currentShowMiniMap && (
                            <MiniMap
                                pannable
                                zoomable
                                nodeColor={(node) => (node as FlowNode).data.color}
                                maskColor={isDark ? "rgba(0,0,0,0.35)" : "rgba(255,255,255,0.62)"}
                            />
                        )}
                    </ReactFlow>
                </ReactFlowProvider>
            </MindContext.Provider>
            {showTextToolbar && (
                <div
                    data-mind-toolbar="true"
                    className={`${styles.styleToolbar} nodrag nopan nowheel`}
                    onMouseDownCapture={(event) => {
                        startToolbarInteraction();
                        saveRichSelection();
                        if (isMindColorDropdownTarget(event.target)) {
                            return;
                        }
                        if (isMindColorTarget(event.target)) {
                            return;
                        }
                        event.preventDefault();
                        event.stopPropagation();
                    }}
                    onMouseDown={(event) => event.stopPropagation()}
                    onClick={(event) => event.stopPropagation()}
                    onDoubleClick={(event) => event.stopPropagation()}
                    onContextMenu={(event) => {
                        event.preventDefault();
                        event.stopPropagation();
                    }}
                >
                    <div className={styles.styleToolbarGroup}>
                        {renderToolbarButton("减小字号", <span>-</span>, () => changeSelectedTextFontSize(-1), toolbarTextFontSize <= minTextFontSize)}
                        <Tooltip title="字号">
                            <span className={styles.styleToolbarValue}>{toolbarTextFontSize}</span>
                        </Tooltip>
                        {renderToolbarButton("增大字号", <span>+</span>, () => changeSelectedTextFontSize(1), toolbarTextFontSize >= maxTextFontSize)}
                        {renderToolbarButton("加粗", <BoldOutlined/>, toggleSelectedTextBold, false, toolbarTextBold)}
                        {renderToolbarButton("斜体", <ItalicOutlined/>, toggleSelectedTextItalic, false, toolbarTextItalic)}
                        {renderToolbarButton("下划线", <UnderlineOutlined/>, () => toggleSelectedTextDecoration("underline"), false, toolbarTextUnderline)}
                        {renderToolbarButton("删除线", <StrikethroughOutlined/>, () => toggleSelectedTextDecoration("line-through"), false, toolbarTextStrike)}
                        {renderToolbarColor("文字颜色", <FontColorsOutlined/>, toolbarTextColor, setSelectedTextColor, clearSelectedTextColor, true)}
                        {renderToolbarButton("清除文字样式", <ClearOutlined/>, clearSelectedTextStyle)}
                    </div>
                </div>
            )}
            {showCardToolbar && selectedNode && (
                <div
                    data-mind-toolbar="true"
                    className={`${styles.styleToolbar} nodrag nopan nowheel`}
                    onMouseDownCapture={(event) => {
                        startToolbarInteraction();
                        if (isMindColorDropdownTarget(event.target)) {
                            return;
                        }
                        if (isMindColorTarget(event.target)) {
                            return;
                        }
                        event.stopPropagation();
                    }}
                    onMouseDown={(event) => event.stopPropagation()}
                    onClick={(event) => event.stopPropagation()}
                    onDoubleClick={(event) => event.stopPropagation()}
                    onContextMenu={(event) => {
                        event.preventDefault();
                        event.stopPropagation();
                    }}
                >
                    <div className={styles.styleToolbarGroup}>
                        {renderToolbarButton("新增子节点", <PlusOutlined/>, () => addChild(selectedNode.id))}
                        {renderToolbarButton("新增同级节点", <ForkOutlined/>, () => addSibling(selectedNode.id), selectedNode.id === rootId)}
                        {renderToolbarButton("编辑节点", <EditOutlined/>, () => startEdit(selectedNode.id))}
                        {renderToolbarButton("删除节点", <DeleteOutlined/>, () => removeNodes([selectedNode.id]), selectedNode.id === rootId)}
                        {renderToolbarButton("撤销", <UndoOutlined/>, undoMind, historyRef.current.length <= 0)}
                        {renderToolbarButton("整理布局", <ApartmentOutlined/>, arrangeMind)}
                    </div>
                    <div className={styles.styleToolbarDivider}/>
                    <div className={styles.styleToolbarGroup}>
                        {renderToolbarColor("填充颜色", <BgColorsOutlined/>, toolbarFillColor, (color) => onToolbarColor("background", color), () => onToolbarClearColor("background"))}
                        {renderToolbarColor("边框颜色", <BorderOutlined/>, toolbarBorderColor, (color) => onToolbarColor("border", color), () => onToolbarClearColor("border"))}
                        {!selectedNode.data.root && renderToolbarColor("分支颜色", <BranchesOutlined/>, toolbarLineColor, (color) => onToolbarColor("line", color), () => onToolbarClearColor("line"))}
                        {renderToolbarButton("减小圆角", <span>-</span>, () => changeSelectedRadius(-2), toolbarRadius <= 0)}
                        <Tooltip title="节点圆角">
                            <span className={styles.styleToolbarValue}>{toolbarRadius}</span>
                        </Tooltip>
                        {renderToolbarButton("增大圆角", <span>+</span>, () => changeSelectedRadius(2), toolbarRadius >= 18)}
                        {renderToolbarButton("节点阴影", <BoxPlotOutlined/>, toggleSelectedShadow, false, toolbarShadow)}
                    </div>
                </div>
            )}
            {contextMenu && (
                <div
                    ref={contextMenuRef}
                    className={styles.contextMenu}
                    style={{left: contextMenu.x, top: contextMenu.y}}
                    onMouseDown={(event) => event.stopPropagation()}
                    onClick={(event) => event.stopPropagation()}
                    onContextMenu={(event) => {
                        event.preventDefault();
                        event.stopPropagation();
                    }}
                >
                    <button
                        type="button"
                        className={styles.contextMenuItem}
                        onClick={() => onContextAction("child")}
                    >
                        <PlusOutlined/>
                        <span>{contextMenu.nodeId ? "新增子节点" : "新增一级节点"}</span>
                    </button>
                    <button
                        type="button"
                        disabled={historyRef.current.length <= 0}
                        className={styles.contextMenuItem}
                        onClick={() => onContextAction("undo")}
                    >
                        <UndoOutlined/>
                        <span>撤销</span>
                    </button>
                    {contextMenu.nodeId && (
                        <>
                            <button
                                type="button"
                                disabled={contextMenu.nodeId === rootId}
                                className={styles.contextMenuItem}
                                onClick={() => onContextAction("sibling")}
                            >
                                <ForkOutlined/>
                                <span>新增同级节点</span>
                            </button>
                            <button
                                type="button"
                                className={styles.contextMenuItem}
                                onClick={() => onContextAction("edit")}
                            >
                                <EditOutlined/>
                                <span>编辑节点</span>
                            </button>
                            <div className={styles.contextMenuDivider}/>
                            <button
                                type="button"
                                disabled={contextMenu.nodeId === rootId}
                                className={`${styles.contextMenuItem} ${styles.contextMenuDanger}`}
                                onClick={() => onContextAction("delete")}
                            >
                                <DeleteOutlined/>
                                <span>删除节点</span>
                            </button>
                        </>
                    )}
                </div>
            )}
        </div>
    );
});

export default Mind;
