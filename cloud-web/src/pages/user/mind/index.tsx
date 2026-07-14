import React, {useEffect, useMemo, useRef, useState} from "react";
import {
    App,
    Button,
    Card,
    Col,
    Drawer,
    Dropdown,
    Empty,
    Form,
    Input,
    Modal,
    Pagination,
    Row,
    Segmented,
    Select,
    Tooltip,
} from "antd";
import {
    ApartmentOutlined,
    CloseOutlined,
    CopyOutlined,
    DeleteOutlined,
    DownloadOutlined,
    DownOutlined,
    EditOutlined,
    EyeOutlined,
    FileImageOutlined,
    FileMarkdownOutlined,
    FileTextOutlined,
    HistoryOutlined,
    LeftOutlined,
    MoreOutlined,
    PlusOutlined,
    RightOutlined,
    SaveOutlined,
    SearchOutlined,
    ShareAltOutlined,
    SortDescendingOutlined,
} from "@ant-design/icons";
import {createStyles} from "antd-style";
import {Body} from "sinking-antd";
import Mind, {
    getMindStats,
    type MindRef,
    type MindSaveReason,
    type MindStats,
    type MindValue,
} from "@/pages/components/mind";
import {ago} from "@/utils/time";

type MindCategory = {
    id: string;
    name: string;
};

type MindDocument = {
    id: string;
    categoryId: string;
    title: string;
    description: string;
    tags: string[];
    value: MindValue;
    nodeCount: number;
    depth: number;
    version: number;
    shared?: boolean;
    createTime: string;
    updateTime: string;
};

type MindVersionRecord = {
    id: string;
    mindId: string;
    version: number;
    saveType: MindVersionSaveType;
    remark: string;
    nodeCount: number;
    depth: number;
    value: MindValue;
    createTime: string;
};

type MindVersionSaveType = "manual" | "auto";

type CreateMindForm = {
    title: string;
    description?: string;
    categoryId: string;
    templateId: string;
};

type CategoryForm = {
    name: string;
};

type MindSortType = "update-desc" | "create-desc";

const mindSortOptions: Array<{label: string; value: MindSortType}> = [
    {label: "最近更新", value: "update-desc"},
    {label: "最近创建", value: "create-desc"},
];

const ALL_CATEGORY_ID = "all";

const previewColors = [
    "#1677ff", "#4096ff", "#0958d9", "#69b1ff",
    "#13c2c2", "#36cfc9", "#08979c", "#5cdbd3",
    "#52c41a", "#73d13d", "#389e0d", "#95de64",
    "#faad14", "#ffc53d", "#d48806", "#ffd666",
    "#fa8c16", "#ffa940", "#d46b08", "#ffc069",
    "#f5222d", "#ff4d4f", "#cf1322", "#ff7875",
    "#eb2f96", "#f759ab", "#c41d7f", "#ff85c0",
    "#722ed1", "#9254de", "#531dab", "#b37feb",
    "#2f54eb", "#597ef7", "#1d39c4", "#85a5ff",
    "#fa541c", "#ff7a45", "#d4380d", "#ff9c6e",
    "#a0d911", "#bae637", "#7cb305", "#d3f261",
];

const getStableColor = (value?: string | number) => {
    const text = String(value || "mind");
    let hash = 0;
    for (let index = 0; index < text.length; index += 1) {
        hash = (hash * 31 + text.charCodeAt(index)) >>> 0;
    }
    return previewColors[hash % previewColors.length];
};

const demoCategoryNames = [
    "需求分析", "用户研究", "竞品调研", "版本规划", "增长策略", "运营活动",
    "交互设计", "视觉规范", "前端架构", "后端服务", "数据库设计", "接口文档",
    "测试方案", "发布流程", "监控告警", "性能优化", "安全合规", "故障复盘",
    "项目管理", "会议记录", "读书笔记", "课程学习", "灵感收集", "写作素材",
    "旅行计划", "健康管理", "理财规划", "年度目标", "团队建设", "客户方案",
];

const initialCategories: MindCategory[] = [
    {id: "product", name: "产品构想"},
    {id: "tech", name: "技术方案"},
    {id: "knowledge", name: "知识整理"},
    {id: "personal", name: "个人计划"},
    {id: "server", name: "服务器管理"},
    ...demoCategoryNames.slice(0, 5).map((name, index) => ({
        id: `demo-category-${index + 1}`,
        name,
    })),
];

const padTime = (value: number) => String(value).padStart(2, "0");

const formatDateTime = (date: Date) => (
    `${date.getFullYear()}-${padTime(date.getMonth() + 1)}-${padTime(date.getDate())} ${padTime(date.getHours())}:${padTime(date.getMinutes())}:${padTime(date.getSeconds())}`
);

const before = (days: number, hours = 0) => {
    const date = new Date();
    date.setDate(date.getDate() - days);
    date.setHours(date.getHours() - hours);
    return formatDateTime(date);
};

const getTimeText = (value: string) => ago(value) || value;

const getDateTimeValue = (value: string) => new Date(value.replace(/-/g, "/")).getTime();

const createValue = (nodes: MindValue["nodes"]): MindValue => ({nodes});

const cloneValue = (value: MindValue): MindValue => JSON.parse(JSON.stringify(value));

const MindPreviewSvg: React.FC = () => (
    <svg className="mind-preview" viewBox="0 0 320 170" aria-hidden="true" focusable="false">
        <rect width="320" height="170" className="preview-surface"/>
        <path className="grid-line" d="M40 0V170M80 0V170M120 0V170M160 0V170M200 0V170M240 0V170M280 0V170M0 34H320M0 68H320M0 102H320M0 136H320"/>
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

const MindCreateHeroGraphic: React.FC = () => (
    <svg className="hero-graphic" viewBox="0 0 390 154" aria-hidden="true" focusable="false">
        <ellipse className="hero-orbit" cx="239" cy="77" rx="126" ry="61"/>
        <path className="hero-link" d="M219 77C184 77 181 40 143 40M219 77C181 77 180 116 139 116"/>
        <path className="hero-link secondary" d="M267 77C302 77 303 38 341 38M267 77C304 77 306 115 345 115"/>
        <path className="hero-link" d="M243 59V25M243 95V130"/>
        <rect className="hero-node" x="91" y="29" width="57" height="22" rx="7"/>
        <rect className="hero-node" x="79" y="105" width="65" height="22" rx="7"/>
        <rect className="hero-node secondary" x="336" y="27" width="48" height="22" rx="7"/>
        <rect className="hero-node secondary" x="340" y="104" width="45" height="22" rx="7"/>
        <rect className="hero-node" x="219" y="8" width="48" height="20" rx="7"/>
        <rect className="hero-node secondary" x="217" y="127" width="52" height="20" rx="7"/>
        <rect className="hero-core" x="217" y="57" width="52" height="40" rx="12"/>
        <circle className="hero-dot" cx="243" cy="77" r="4"/>
        <circle className="hero-dot" cx="113" cy="40" r="2.5"/>
        <circle className="hero-dot" cx="362" cy="115" r="2.5"/>
    </svg>
);

const withStats = (data: Omit<MindDocument, "nodeCount" | "depth">): MindDocument => {
    const stats = getMindStats(data.value);
    return {
        ...data,
        nodeCount: stats.nodeCount,
        depth: stats.depth,
    };
};

const defaultValue = (title: string, color = "#1677ff") => createValue([
    {id: "root", title, root: true, color},
    {id: "goal", parentId: "root", title: "目标", color: "#13c2c2"},
    {id: "plan", parentId: "root", title: "计划", color: "#52c41a"},
    {id: "next", parentId: "plan", title: "下一步", color: "#faad14"},
]);

const templateItems = [
    {
        id: "blank",
        name: "空白导图",
        description: "中心主题和基础分支",
        color: "#1677ff",
        value: defaultValue("中心主题", "#1677ff"),
    },
    {
        id: "product",
        name: "产品路线",
        description: "适合拆功能和排优先级",
        color: "#13c2c2",
        value: createValue([
            {id: "root", title: "产品路线", root: true, color: "#13c2c2"},
            {id: "position", parentId: "root", title: "定位", color: "#1677ff"},
            {id: "user", parentId: "root", title: "用户", color: "#52c41a"},
            {id: "feature", parentId: "root", title: "核心功能", color: "#faad14"},
            {id: "growth", parentId: "root", title: "增长", color: "#eb2f96"},
        ]),
    },
    {
        id: "tech",
        name: "技术方案",
        description: "适合梳理模块和链路",
        color: "#722ed1",
        value: createValue([
            {id: "root", title: "技术方案", root: true, color: "#722ed1"},
            {id: "api", parentId: "root", title: "接口", color: "#1677ff"},
            {id: "data", parentId: "root", title: "数据", color: "#13c2c2"},
            {id: "queue", parentId: "root", title: "消息队列", color: "#52c41a"},
            {id: "ops", parentId: "root", title: "部署运维", color: "#faad14"},
        ]),
    },
];

const demoDocumentTopics = [
    "产品需求拆解", "季度路线规划", "用户访谈总结", "核心流程优化", "增长实验计划",
    "技术架构升级", "接口治理方案", "数据库迁移", "性能排查清单", "安全建设计划",
    "项目里程碑", "版本发布检查", "线上问题复盘", "团队协作流程", "知识库结构",
    "个人学习计划", "年度目标拆解", "读书笔记整理", "旅行行程规划", "内容选题规划",
];

const demoDocuments: MindDocument[] = Array.from({length: 7}, (_, index) => {
    const category = initialCategories[index % initialCategories.length];
    const topic = demoDocumentTopics[index % demoDocumentTopics.length];
    const sequence = index + 1;
    return withStats({
        id: `mind-demo-${sequence}`,
        categoryId: category.id,
        title: `${topic} ${String(sequence).padStart(2, "0")}`,
        description: `用于验证「${category.name}」分类下的列表滚动、搜索、排序与分页效果。`,
        tags: [category.name, "演示数据", `第${sequence}组`],
        version: (index % 6) + 1,
        shared: index % 7 === 0,
        createTime: before((index % 28) + 2, index % 8),
        updateTime: before(index % 16, index % 12),
        value: defaultValue(topic, getStableColor(category.id)),
    });
});

const initialDocuments: MindDocument[] = [
    withStats({
        id: "mind-product-roadmap",
        categoryId: "product",
        title: "云上豁者产品路线图",
        description: "个人工作台、博客主页和工具插件的产品结构。",
        tags: ["工作台", "博客", "插件"],
        shared: true,
        version: 8,
        createTime: before(9),
        updateTime: before(0, 2),
        value: createValue([
            {id: "root", title: "云上豁者", root: true, color: "#1677ff"},
            {id: "workspace", parentId: "root", title: "个人工作台", description: "工具统一入口", color: "#13c2c2"},
            {id: "mind", parentId: "workspace", title: "思维导图", color: "#52c41a"},
            {id: "doc", parentId: "workspace", title: "文档笔记", color: "#faad14"},
            {id: "server", parentId: "workspace", title: "远程连接", color: "#722ed1"},
            {id: "blog", parentId: "root", title: "个人主页", description: "模板和内容系统", color: "#eb2f96"},
            {id: "market", parentId: "blog", title: "模板市场", color: "#fa541c"},
        ]),
    }),
    withStats({
        id: "mind-ticket-message",
        categoryId: "tech",
        title: "工单消息通知链路",
        description: "梳理站内消息、WebSocket 刷新事件和异步队列。",
        tags: ["消息", "WebSocket"],
        version: 4,
        createTime: before(7),
        updateTime: before(1, 3),
        value: createValue([
            {id: "root", title: "工单处理流程", root: true, color: "#13c2c2"},
            {id: "create", parentId: "root", title: "用户发起", color: "#1677ff"},
            {id: "assign", parentId: "root", title: "分配处理人", color: "#52c41a"},
            {id: "handle", parentId: "assign", title: "处理中", description: "记录进度和附件", color: "#faad14"},
            {id: "notify", parentId: "handle", title: "站内消息通知", color: "#eb2f96"},
            {id: "done", parentId: "root", title: "完结归档", color: "#722ed1"},
        ]),
    }),
    withStats({
        id: "mind-deploy",
        categoryId: "tech",
        title: "三节点高可用部署",
        description: "数据库、缓存、消息队列和应用层拆分。",
        tags: ["PostgreSQL", "Redis", "NATS"],
        version: 5,
        createTime: before(5),
        updateTime: before(2),
        value: createValue([
            {id: "root", title: "三节点部署", root: true, color: "#722ed1"},
            {id: "db", parentId: "root", title: "PostgreSQL", color: "#1677ff"},
            {id: "redis", parentId: "root", title: "Redis", color: "#13c2c2"},
            {id: "nats", parentId: "root", title: "NATS", color: "#52c41a"},
            {id: "app", parentId: "root", title: "K8s 应用层", color: "#faad14"},
            {id: "obs", parentId: "app", title: "监控日志", color: "#eb2f96"},
        ]),
    }),
    withStats({
        id: "mind-blog-template",
        categoryId: "knowledge",
        title: "个人主页模板引擎",
        description: "模板包、变量、KV 存储和静态资源缓存。",
        tags: ["模板", "个人主页"],
        version: 2,
        createTime: before(4),
        updateTime: before(3, 4),
        value: createValue([
            {id: "root", title: "模板引擎", root: true, color: "#52c41a"},
            {id: "page", parentId: "root", title: "页面渲染", color: "#1677ff"},
            {id: "asset", parentId: "root", title: "资源缓存", color: "#13c2c2"},
            {id: "kv", parentId: "root", title: "KV 数据", color: "#faad14"},
            {id: "market", parentId: "root", title: "模板市场", color: "#eb2f96"},
        ]),
    }),
    withStats({
        id: "mind-week-plan",
        categoryId: "personal",
        title: "本周开发计划",
        description: "把脑子里的东西拆小，按优先级推进。",
        tags: ["计划", "复盘"],
        version: 1,
        createTime: before(1),
        updateTime: before(1),
        value: defaultValue("本周计划", "#eb2f96"),
    }),
    ...demoDocuments,
];

const MAX_VERSION_RECORDS_PER_TYPE = 20;
// Keep frequent draft saves in one rolling snapshot instead of flooding version history.
const AUTO_VERSION_SNAPSHOT_INTERVAL = 60 * 1000;

const createVersionRecord = (
    record: MindDocument,
    remark = "保存快照",
    saveType: MindVersionSaveType = "manual",
): MindVersionRecord => ({
    id: `${record.id}-v${record.version}-${Date.now().toString(36)}`,
    mindId: record.id,
    version: record.version,
    saveType,
    remark,
    nodeCount: record.nodeCount,
    depth: record.depth,
    value: cloneValue(record.value),
    createTime: record.updateTime,
});

const limitVersionRecords = (records: MindVersionRecord[]) => {
    const counts: Record<MindVersionSaveType, number> = {manual: 0, auto: 0};
    return records.filter((record) => {
        counts[record.saveType] += 1;
        return counts[record.saveType] <= MAX_VERSION_RECORDS_PER_TYPE;
    });
};

const getMindValueSignature = (value: MindValue) => JSON.stringify(value);

const initialVersions = initialDocuments.reduce<Record<string, MindVersionRecord[]>>((data, item) => ({
    ...data,
    [item.id]: [createVersionRecord(item, "初始版本")],
}), {});

const downloadDataUrl = (dataUrl: string, filename: string) => {
    const link = document.createElement("a");
    link.href = dataUrl;
    link.download = filename;
    link.click();
};

const useStyles = createStyles(({css, token, isDarkMode}: any) => ({
    page: css`
        --mind-page-height: calc(100dvh - 177px);
        --mind-library-min-height: calc(var(--mind-page-height) - 63px);
        min-height: var(--mind-library-min-height);

        @supports not (height: 100dvh) {
            --mind-page-height: calc(100vh - 177px);
        }
    `,
    library: css`
        container-name: mind-library;
        container-type: inline-size;
        min-width: 0;
        min-height: inherit;
    `,
    libraryCard: css`
        min-height: var(--mind-library-min-height);
        overflow: hidden;
        border-radius: ${token.borderRadiusLG}px;
        background: ${token.colorBgContainer};

        > .ant-card-body {
            display: grid;
            grid-template-columns: 238px minmax(0, 1fr);
            grid-template-rows: auto auto minmax(0, 1fr);
            min-width: 0;
            min-height: var(--mind-library-min-height);
            padding: 0 !important;
        }

        @container mind-library (max-width: 780px) {
            > .ant-card-body {
                grid-template-columns: 1fr;
                grid-template-rows: auto auto auto minmax(0, 1fr);
            }
        }
    `,
    categoryPanel: css`
        position: relative;
        display: block;
        grid-column: 1;
        grid-row: 2 / span 2;
        width: auto;
        min-width: 0;
        min-height: 0;
        max-height: none;
        padding: 0;
        border-radius: 0 0 0 ${token.borderRadiusLG}px;
        background: ${token.colorBgContainer};
        align-self: stretch;
        overflow: hidden;

        @container mind-library (max-width: 780px) {
            grid-column: 1;
            grid-row: 2;
            width: 100%;
            max-height: none;
            border-radius: 0;
            align-self: stretch;
        }

        @media (max-width: 640px) {
            width: 100%;
            height: auto;
            max-height: none;
        }
    `,
    categoryInner: css`
        position: absolute;
        inset: 0;
        display: flex;
        flex-direction: column;
        min-width: 0;
        min-height: 0;
        padding: 10px 4px 10px 10px;

        @container mind-library (max-width: 780px) {
            position: static;
        }

        @media (max-width: 640px) {
            padding: 9px;
        }
    `,
    categoryRail: css`
        display: flex;
        flex: 1;
        flex-direction: column;
        min-width: 0;
        min-height: 0;

        @container mind-library (max-width: 780px) {
            flex-direction: row;
            align-items: center;
            gap: 6px;
            width: 100%;
        }
    `,
    categoryScrollButton: css`
        display: none;
        flex: 0 0 auto;
        align-items: center;
        justify-content: center;
        width: 28px;
        height: 34px;
        padding: 0;
        border: 0;
        border-radius: ${token.borderRadius}px;
        outline: none;
        background: ${token.colorFillQuaternary};
        color: ${token.colorTextSecondary};
        cursor: pointer;
        transition: background 0.16s ease, color 0.16s ease, opacity 0.16s ease;

        &:hover:not(:disabled),
        &:focus-visible:not(:disabled) {
            background: ${token.colorFillTertiary};
            color: ${token.colorPrimary};
        }

        &:focus-visible {
            box-shadow: 0 0 0 2px ${token.colorPrimaryBg};
        }

        &:disabled {
            opacity: 0.3;
            cursor: default;
        }

        .anticon {
            font-size: 11px;
        }

        @container mind-library (max-width: 780px) {
            display: inline-flex;
        }

        @media (max-width: 640px), (hover: none) and (pointer: coarse) {
            display: none;
        }
    `,
    categoryFooter: css`
        flex: none;
        padding-top: 8px;

        @container mind-library (max-width: 780px) {
            display: none;
        }
    `,
    categoryAddCard: css`
        display: flex;
        align-items: center;
        gap: 8px;
        width: 100%;
        min-width: 0;
        min-height: 38px;
        padding: 0 10px;
        border: 1px dashed ${token.colorBorder};
        border-radius: ${token.borderRadius}px;
        outline: none;
        appearance: none;
        background: transparent;
        color: ${token.colorTextTertiary};
        cursor: pointer;
        font: inherit;
        font-size: 12px;
        line-height: 20px;
        text-align: left;
        transition: border-color 0.16s ease, background 0.16s ease, color 0.16s ease;

        &.mobile {
            display: none;
        }

        &:hover,
        &:focus-visible {
            border-color: ${token.colorPrimaryBorder};
            background: ${token.colorFillQuaternary};
            color: ${token.colorPrimary};
        }

        &:focus-visible {
            box-shadow: 0 0 0 2px ${token.colorPrimaryBg};
        }

        .icon {
            display: inline-flex;
            align-items: center;
            justify-content: center;
            width: 16px;
            height: 16px;
            border-radius: ${token.borderRadiusSM}px;
            background: ${token.colorFillTertiary};
            font-size: 10px;
        }

        @container mind-library (max-width: 780px) {
            min-height: 34px;
            padding: 0 9px;

            &.mobile {
                display: flex;
                flex: 0 0 auto;
                width: auto;
                white-space: nowrap;
            }
        }
    `,
    categoryScroll: css`
        flex: 1;
        min-height: 0;
        overflow-y: scroll;
        padding: 0 2px 0 0;
        scrollbar-color: ${token.colorBorderSecondary} transparent;
        scrollbar-width: thin;

        &::-webkit-scrollbar {
            width: 6px;
        }

        &::-webkit-scrollbar-track {
            background: transparent;
        }

        &::-webkit-scrollbar-thumb {
            border-radius: ${token.borderRadiusSM}px;
            background: ${token.colorBorderSecondary};
        }

        &::-webkit-scrollbar-thumb:hover {
            background: ${token.colorTextQuaternary};
        }

        @container mind-library (max-width: 780px) {
            overflow-x: auto;
            overflow-y: hidden;
            padding: 0;
            scrollbar-width: none;

            &::-webkit-scrollbar {
                display: none;
            }
        }
    `,
    categoryGrid: css`
        display: flex;
        flex-direction: column;
        gap: 4px;

        @container mind-library (max-width: 780px) {
            flex-direction: row;
            gap: 6px;
            min-width: max-content;
        }
    `,
    categoryCard: css`
        position: relative;
        display: flex;
        align-items: center;
        min-width: 0;
        min-height: 40px;
        padding: 0 7px 0 11px;
        overflow: hidden;
        border: 0;
        border-radius: ${token.borderRadius}px;
        background: transparent;
        cursor: pointer;
        transition: border-color 0.16s ease, background 0.16s ease, color 0.16s ease;

        &:hover,
        &.active {
            background: ${token.colorFillQuaternary};
        }

        &.active {
            color: ${token.colorTextHeading};
            background: ${token.colorPrimaryBg};
        }

        .row {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 8px;
            width: 100%;
            min-width: 0;
        }

        .right {
            display: flex;
            align-items: center;
            gap: 3px;
            flex: none;
        }

        .name-wrap {
            display: inline-flex;
            align-items: center;
            gap: 8px;
            min-width: 0;
        }

        .color {
            flex: 0 0 auto;
            width: 7px;
            height: 7px;
            border-radius: 50%;
            background: color-mix(in srgb, var(--accent), ${token.colorTextTertiary} 8%);
            opacity: 0.62;
        }

        .action {
            display: inline-flex;
            margin-right: -3px;
            opacity: 0;
            transition: opacity 0.16s ease;
        }

        &:hover .action,
        &.active .action {
            opacity: 1;
        }

        .name {
            min-width: 0;
            overflow: hidden;
            color: ${token.colorTextSecondary};
            font-size: 13px;
            font-weight: 400;
            line-height: 20px;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        &.active .name {
            color: ${token.colorTextHeading};
            font-weight: 600;
        }

        &.active .color {
            background: var(--accent);
            box-shadow: 0 0 0 3px color-mix(in srgb, var(--accent), transparent 78%);
            opacity: 1;
        }

        .action .ant-btn {
            color: ${token.colorTextTertiary};
        }

        @container mind-library (max-width: 780px) {
            flex: 0 0 auto;
            min-height: 34px;
            padding: 0 8px 0 9px;
            background: ${token.colorFillQuaternary};

            .name-wrap {
                gap: 6px;
            }

            .color {
                width: 7px;
                height: 7px;
            }

            .action {
                opacity: 1;
            }
        }
    `,
    content: css`
        display: contents;
        min-width: 0;

        > .ant-row {
            display: flex;
            flex-direction: column;
            flex-wrap: nowrap;
            grid-column: 2;
            grid-row: 3;
            height: 100%;
            row-gap: 0 !important;
            column-gap: 0 !important;
            margin: 0 !important;
            min-width: 0;
        }

        > .ant-row > .ant-col {
            flex: none;
            width: 100%;
            max-width: 100%;
            padding: 0 !important;
        }

        @container mind-library (max-width: 780px) {
            display: contents;

            > .ant-row {
                grid-column: 1;
                grid-row: 4;
            }
        }
    `,
    contentPanel: css`
        display: contents;

        @container mind-library (max-width: 780px) {
            display: contents;
        }
    `,
    documentDataCard: css`
        display: block;
        flex: 1;
        width: 100%;
        background: ${token.colorBgContainer};
        padding: 0 6px 24px;

        &.empty {
            display: flex;
            align-items: center;
            justify-content: center;
            min-height: 240px;
            padding: 24px;
        }

        @media (max-width: 640px) {
            padding: 0 11px 12px;
        }
    `,
    documentDataColumn: css`
        display: flex;
        flex: 1 1 auto !important;
        min-height: 0;
    `,
    documentPaginationCard: css`
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 12px 15px;
        background: transparent;

        .ant-pagination {
            width: 100%;
            justify-content: center;
        }

        .ant-pagination-total-text {
            color: ${token.colorTextTertiary};
            font-size: 12px;
        }

        .ant-pagination-options {
            margin-inline-start: 10px;
        }

        .ant-pagination-item,
        .ant-pagination-prev .ant-pagination-item-link,
        .ant-pagination-next .ant-pagination-item-link {
            min-width: 24px;
            height: 24px;
            margin-inline-end: 4px;
            border: 0;
            border-radius: ${token.borderRadiusSM}px;
            background: transparent;
            color: ${token.colorTextSecondary};
            font-size: 11px;
            line-height: 24px;
        }

        .ant-pagination-item a {
            color: inherit;
        }

        .ant-pagination-item:hover,
        .ant-pagination-prev:hover .ant-pagination-item-link,
        .ant-pagination-next:hover .ant-pagination-item-link {
            background: ${token.colorFillTertiary};
            color: ${token.colorTextSecondary};
        }

        .ant-pagination-item-active,
        .ant-pagination-item-active:hover {
            background: ${token.colorBgContainer};
            box-shadow: inset 0 0 0 1px ${token.colorPrimaryBorder};
            color: ${token.colorPrimary};
        }

        .ant-pagination-item-active a,
        .ant-pagination-item-active:hover a {
            color: ${token.colorPrimary};
        }

        .ant-pagination-options-size-changer.ant-select .ant-select-selector {
            height: 24px !important;
            padding-inline: 7px 22px !important;
            border: 0 !important;
            border-radius: ${token.borderRadiusSM}px !important;
            background: ${token.colorBgContainer} !important;
            box-shadow: inset 0 0 0 1px ${token.colorBorderSecondary} !important;
        }

        .ant-pagination-options-size-changer .ant-select-selection-item,
        .ant-pagination-options-quick-jumper {
            color: ${token.colorTextTertiary};
            font-size: 11px;
            line-height: 24px;
        }

        .ant-pagination-options-quick-jumper input {
            width: 32px;
            height: 24px;
            margin-inline: 4px;
            border-radius: ${token.borderRadiusSM}px;
            color: ${token.colorTextSecondary};
            font-size: 11px;
        }

        @container mind-library (max-width: 780px) {
            padding: 12px;
        }
    `,
    documentPaginationColumn: css`
        margin-top: auto;
    `,
    pageSizeSelect: css`
        && {
            height: 24px;
            font-size: 11px;
        }

        && .ant-select-selector {
            height: 24px !important;
            padding-inline: 7px 22px !important;
            border: 0 !important;
            border-radius: ${token.borderRadiusSM}px !important;
            background: ${token.colorBgContainer} !important;
            box-shadow: inset 0 0 0 1px ${token.colorBorderSecondary} !important;
        }

        && .ant-select-selection-item {
            color: ${token.colorTextSecondary} !important;
            font-size: 11px !important;
            font-weight: 400 !important;
            line-height: 24px !important;
        }

        && .ant-select-arrow {
            color: ${token.colorTextQuaternary};
            font-size: 9px;
        }
    `,
    pageSizeDropdown: css`
        && {
            padding: 2px;
            border-radius: ${token.borderRadius}px;
        }

        && .ant-select-item {
            min-height: 26px;
            padding: 4px 8px;
            border-radius: ${token.borderRadiusSM}px;
            color: ${token.colorTextSecondary};
            font-size: 11px;
            line-height: 18px;
        }

        && .ant-select-item-option-content {
            font-size: 11px !important;
            font-weight: 400;
        }

        && .ant-select-item-option-selected:not(.ant-select-item-option-disabled) {
            background: ${token.colorFillTertiary};
            color: ${token.colorTextSecondary};
            font-weight: 400;
        }
    `,
    searchCard: css`
        grid-column: 2;
        grid-row: 2;
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 16px;
        min-width: 0;
        padding: 16px 12px 12px;
        background: ${token.colorBgContainer};

        @media (max-width: 640px) {
            gap: 8px;
            padding: 12px 12px 10px;
        }

        @container mind-library (max-width: 780px) {
            grid-column: 1;
            grid-row: 3;
        }
    `,
    searchBox: css`
        flex: 0 1 300px;
        width: 300px;
        max-width: 100%;
        min-width: 0;

        &.ant-input-affix-wrapper {
            height: 36px;
            padding: 0 11px;
            border: 1px solid transparent;
            border-radius: ${token.borderRadiusSM}px;
            background: color-mix(in srgb, ${token.colorFillTertiary} 45%, ${token.colorFillQuaternary});
            box-shadow: none;
            transition: border-color 0.16s ease, background-color 0.16s ease, box-shadow 0.16s ease;
        }

        &.ant-input-affix-wrapper:hover {
            background: color-mix(in srgb, ${token.colorFillTertiary} 65%, ${token.colorFillQuaternary});
        }

        &.ant-input-affix-wrapper-focused {
            border-color: ${token.colorPrimaryBorder};
            background: ${token.colorBgContainer};
            box-shadow: 0 0 0 3px ${token.colorPrimaryBg};
        }

        .ant-input-prefix {
            margin-right: 8px;
            color: ${token.colorTextTertiary};
            font-size: 14px;
        }

        .ant-input {
            background: transparent;
            color: color-mix(in srgb, ${token.colorTextSecondary}, ${token.colorBgContainer} 24%);
            font-size: 11px;
        }

        .ant-input::placeholder {
            color: ${token.colorTextQuaternary};
            font-size: 11px;
        }

        @media (max-width: 820px) {
            flex-basis: 260px;
            width: 260px;
        }

        @media (max-width: 640px) {
            flex: 1 1 auto;
            width: auto;
        }
    `,
    sortTrigger: css`
        display: inline-flex;
        flex: none;
        align-items: center;
        justify-content: center;
        gap: 5px;
        height: 32px;
        padding: 0 8px;
        border: 1px solid transparent;
        border-radius: ${token.borderRadiusSM}px;
        outline: none;
        appearance: none;
        background: ${token.colorFillQuaternary};
        color: ${token.colorTextSecondary};
        cursor: pointer;
        font: inherit;
        line-height: 1;
        transition: background-color 0.16s ease, color 0.16s ease;

        .marker {
            color: ${token.colorTextQuaternary};
            font-size: 11px;
        }

        .value {
            color: ${token.colorTextTertiary};
            font-size: 11px;
            font-weight: 400;
            white-space: nowrap;
            transition: color 0.16s ease;
        }

        .arrow {
            color: ${token.colorTextQuaternary};
            font-size: 9px;
            transition: color 0.16s ease, transform 0.16s ease;
        }

        &:hover,
        &.ant-dropdown-open {
            background: ${token.colorFillSecondary};
        }

        &:hover .value,
        &.ant-dropdown-open .value {
            color: ${token.colorTextSecondary};
        }

        &.ant-dropdown-open .arrow {
            color: ${token.colorTextTertiary};
            transform: rotate(180deg);
        }

        &:focus-visible {
            border-color: ${token.colorPrimaryBorder};
            background: ${token.colorBgContainer};
            box-shadow: 0 0 0 3px ${token.colorPrimaryBg};
        }

        @media (max-width: 640px) {
            flex: 0 0 auto;
        }
    `,
    sortDropdown: css`
        && {
            min-width: 100px;
            padding: 0;
            border-radius: ${token.borderRadius}px;
        }

        && .ant-dropdown-menu {
            padding: 2px !important;
            border-radius: ${token.borderRadius}px !important;
            background: ${token.colorBgElevated};
            box-shadow: ${token.boxShadowSecondary};
        }

        && .ant-dropdown-menu .ant-dropdown-menu-item {
            min-height: 26px !important;
            margin: 0 !important;
            padding: 0 8px !important;
            border-radius: ${token.borderRadiusSM}px !important;
            color: ${token.colorTextTertiary};
            font-size: 11px !important;
            font-weight: 400;
            line-height: 26px !important;
        }

        && .ant-dropdown-menu .ant-dropdown-menu-item .ant-dropdown-menu-title-content {
            font-size: 11px !important;
            line-height: 26px !important;
        }

        && .ant-dropdown-menu .ant-dropdown-menu-item:hover {
            background: ${token.colorFillQuaternary};
            color: ${token.colorTextSecondary};
        }

        && .ant-dropdown-menu .ant-dropdown-menu-item-selected,
        && .ant-dropdown-menu .ant-dropdown-menu-item-selected:hover {
            background: color-mix(in srgb, ${token.colorPrimary}, ${token.colorBgElevated} 94%);
            color: ${token.colorTextSecondary};
        }

        && .ant-dropdown-menu .ant-dropdown-menu-item .ant-dropdown-menu-item-icon {
            margin-inline-end: 4px !important;
            color: ${token.colorTextQuaternary};
            font-size: 9px !important;
        }
    `,
    documentGrid: css`
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(208px, 1fr));
        gap: 14px;
        padding: 5px 6px 8px;

        @media (max-width: 640px) {
            grid-template-columns: 1fr;
            gap: 10px;
            padding: 3px 1px 6px;
        }
    `,
    createHero: css`
        grid-column: 1 / -1;
        grid-row: 1;
        position: relative;
        display: flex;
        align-items: center;
        justify-content: space-between;
        min-height: 128px;
        margin-bottom: 0;
        overflow: hidden;
        border: 0;
        border-radius: ${token.borderRadiusLG}px ${token.borderRadiusLG}px 0 0;
        background: ${isDarkMode
            ? "linear-gradient(118deg, rgba(22,119,255,0.1) 0%, rgba(255,255,255,0.035) 58%, rgba(114,46,209,0.08) 100%)"
            : "linear-gradient(118deg, #f7faff 0%, #eef4ff 58%, #f7f4ff 100%)"};
        box-shadow: none;

        &::before {
            position: absolute;
            inset: 0;
            background-image:
                linear-gradient(rgba(72, 132, 202, 0.055) 1px, transparent 1px),
                linear-gradient(90deg, rgba(72, 132, 202, 0.055) 1px, transparent 1px);
            background-size: 24px 24px;
            content: "";
            -webkit-mask-image: linear-gradient(90deg, transparent 16%, rgba(0, 0, 0, 0.18) 55%, #000 100%);
            mask-image: linear-gradient(90deg, transparent 16%, rgba(0, 0, 0, 0.18) 55%, #000 100%);
            pointer-events: none;
        }

        &::after {
            position: absolute;
            top: 50%;
            right: 20%;
            width: 150px;
            height: 150px;
            border-radius: 50%;
            background: color-mix(in srgb, ${token.colorPrimary}, transparent 90%);
            content: "";
            filter: blur(48px);
            pointer-events: none;
            transform: translateY(-50%);
        }

        .hero-copy {
            position: relative;
            z-index: 2;
            flex: 1;
            min-width: 0;
            max-width: 68%;
            padding: 18px 24px 19px;
        }

        .eyebrow {
            display: inline-flex;
            align-items: center;
            gap: 7px;
            color: color-mix(in srgb, ${token.colorPrimary}, ${token.colorTextTertiary} 62%);
            font-size: 10px;
            font-weight: 600;
            line-height: 17px;
            letter-spacing: 0.09em;
            text-transform: uppercase;
        }

        .status-dot {
            width: 5px;
            height: 5px;
            border-radius: 50%;
            background: color-mix(in srgb, ${token.colorPrimary}, #ffffff 28%);
            box-shadow: 0 0 0 4px color-mix(in srgb, ${token.colorPrimary}, transparent 91%);
        }

        .hero-title {
            margin-top: 3px;
            color: color-mix(in srgb, ${token.colorTextHeading}, transparent 8%);
            font-size: 19px;
            font-weight: 600;
            line-height: 27px;
            letter-spacing: 0.005em;
        }

        .hero-desc {
            max-width: 480px;
            margin-top: 2px;
            color: ${token.colorTextTertiary};
            font-size: 12px;
            line-height: 19px;
        }

        .create-button {
            position: relative;
            z-index: 3;
            display: inline-flex;
            align-items: center;
            justify-content: center;
            flex: none;
            box-sizing: border-box;
            width: auto;
            min-width: 88px;
            height: 36px;
            margin: 0 24px 0 0;
            padding: 0 17px;
            border: 0;
            border-radius: ${token.borderRadius}px;
            outline: none;
            appearance: none;
            background: ${token.colorPrimary};
            color: #ffffff;
            cursor: pointer;
            font: inherit;
            font-size: 13px;
            font-weight: 500;
            line-height: 1;
            white-space: nowrap;
            transform: none;
            transition: background-color 0.18s ease;
        }

        .create-button:hover,
        .create-button:focus-visible {
            background: ${token.colorPrimaryHover};
            color: #ffffff;
            transform: none;
        }

        .create-button:focus-visible {
            box-shadow: 0 0 0 3px ${token.colorPrimaryBg};
        }

        .create-button:active {
            background: ${token.colorPrimaryActive};
            transform: none;
        }

        .hero-visual {
            position: absolute;
            top: 0;
            right: 106px;
            z-index: 1;
            width: 38%;
            height: 100%;
            opacity: 0.34;
            pointer-events: none;
        }

        .hero-graphic {
            display: block;
            width: 100%;
            height: 100%;
            overflow: visible;
        }

        .hero-orbit {
            fill: none;
            stroke: color-mix(in srgb, ${token.colorPrimary}, transparent 82%);
            stroke-dasharray: 5 7;
            stroke-width: 1;
        }

        .hero-link {
            fill: none;
            stroke: color-mix(in srgb, ${token.colorPrimary}, transparent 62%);
            stroke-linecap: round;
            stroke-width: 1.3;
        }

        .hero-link.secondary {
            stroke: rgba(114, 46, 209, 0.3);
        }

        .hero-node {
            fill: ${isDarkMode ? "rgba(255,255,255,0.045)" : "rgba(255,255,255,0.58)"};
            stroke: color-mix(in srgb, ${token.colorPrimary}, transparent 58%);
            stroke-width: 1;
        }

        .hero-node.secondary {
            stroke: rgba(114, 46, 209, 0.34);
        }

        .hero-core {
            fill: color-mix(in srgb, ${token.colorPrimary}, transparent 42%);
            stroke: color-mix(in srgb, ${token.colorPrimary}, #ffffff 48%);
            stroke-width: 1;
        }

        .hero-dot {
            fill: color-mix(in srgb, ${token.colorPrimary}, #ffffff 22%);
        }

        @media (max-width: 720px) {
            min-height: 126px;

            .hero-copy {
                max-width: calc(100% - 118px);
                padding: 16px 18px;
            }

            .hero-title {
                font-size: 17px;
                line-height: 24px;
            }

            .hero-desc {
                display: -webkit-box;
                overflow: hidden;
                -webkit-line-clamp: 2;
                -webkit-box-orient: vertical;
            }

            .create-button {
                height: 36px;
                margin-right: 16px;
            }

            .hero-visual {
                right: 76px;
                width: 50%;
                opacity: 0.2;
            }
        }

        @media (max-width: 420px) {
            min-height: 132px;

            .hero-copy {
                max-width: calc(100% - 104px);
                padding: 15px 8px 15px 16px;
            }

            .eyebrow {
                letter-spacing: 0.06em;
            }

            .hero-title {
                font-size: 15px;
                line-height: 22px;
            }

            .hero-desc {
                font-size: 11px;
                line-height: 17px;
            }

            .create-button {
                margin-right: 12px;
                padding: 0 14px;
            }

            .hero-visual {
                right: 58px;
                width: 64%;
                opacity: 0.13;
            }
        }
    `,
    documentCard: css`
        position: relative;
        min-width: 0;
        min-height: 184px;
        overflow: hidden;
        border-radius: ${token.borderRadiusLG}px;
        background: ${token.colorBgContainer};
        cursor: pointer;
        transition: border-color 0.16s ease, transform 0.16s ease;

        &.ant-card > .ant-card-body {
            padding: 0 !important;
        }

        &:hover {
            border-color: ${token.colorPrimary};
        }

        .preview {
            position: relative;
            height: 104px;
            overflow: hidden;
            border-bottom: 1px solid ${token.colorBorderSecondary};
            background: ${isDarkMode ? "rgba(255,255,255,0.012)" : "#fbfcff"};
        }

        .mind-preview {
            display: block;
            width: 100%;
            height: 100%;
        }

        .preview-surface {
            fill: ${isDarkMode ? "rgba(255,255,255,0.012)" : "#fbfcff"};
        }

        .grid-line {
            fill: none;
            stroke: ${isDarkMode ? "rgba(255,255,255,0.035)" : "rgba(15, 23, 42, 0.035)"};
            stroke-width: 1;
        }

        .preview-link {
            fill: none;
            stroke: color-mix(in srgb, var(--accent), ${token.colorTextTertiary} 28%);
            stroke-linecap: round;
            stroke-width: 1.8;
        }

        .preview-root,
        .preview-node {
            fill: ${token.colorBgContainer};
            stroke: color-mix(in srgb, var(--accent), ${token.colorBorderSecondary} 46%);
            stroke-width: 1.2;
        }

        .preview-dot {
            fill: color-mix(in srgb, var(--accent), ${token.colorBgContainer} 18%);
        }

        .body {
            padding: 11px 13px 12px;
        }

        .head {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 8px;
            min-height: 24px;
        }

        .actions {
            display: flex;
            align-items: center;
            gap: 2px;
            flex: none;
            margin-right: 0;
        }

        .title {
            display: block;
            flex: 1 1 auto;
            min-width: 0;
            padding: 0;
            border: 0;
            background: transparent;
            color: ${token.colorTextHeading};
            cursor: pointer;
            font-size: 14px;
            font-weight: 600;
            line-height: 21px;
            text-align: left;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .title:hover {
            color: ${token.colorPrimary};
        }

        .desc {
            display: -webkit-box;
            min-height: 18px;
            margin-top: 3px;
            overflow: hidden;
            color: ${token.colorTextTertiary};
            font-size: 12px;
            line-height: 18px;
            -webkit-line-clamp: 1;
            -webkit-box-orient: vertical;
        }

        .meta {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 8px;
            margin-top: 7px;
            color: ${token.colorTextTertiary};
            font-size: 12px;
            line-height: 18px;
        }

        .meta span {
            min-width: 0;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        @media (max-width: 640px) {
            min-height: 0;
            border-color: ${token.colorBorderSecondary};

            > .ant-card-body {
                display: block;
            }

            .preview {
                height: 124px;
                min-height: 0;
                border-right: 0;
                border-bottom: 1px solid ${token.colorBorderSecondary};
            }

            .mind-preview {
                transform: scale(1.08);
            }

            .body {
                min-width: 0;
                padding: 12px 14px 13px;
            }

            .head {
                gap: 8px;
                min-height: 24px;
            }

            .title {
                font-size: 14px;
                line-height: 21px;
            }

            .desc {
                -webkit-line-clamp: 1;
                min-height: 18px;
                margin-top: 3px;
            }

            .meta {
                flex-wrap: nowrap;
                justify-content: space-between;
                gap: 12px;
                margin-top: 8px;
            }
        }

        @media (max-width: 420px) {
            .preview {
                height: 112px;
            }

            .body {
                padding: 11px 12px 12px;
            }
        }

    `,
    editor: css`
        position: relative;
        height: var(--mind-page-height);
        min-height: min(680px, var(--mind-page-height));
        overflow: hidden;
        border: 1px solid ${token.colorBorderSecondary};
        border-radius: ${token.borderRadiusLG}px;
        background: ${token.colorBgContainer};
    `,
    editorDrawer: css`
        .ant-drawer-content-wrapper {
            width: 100vw !important;
            max-width: 100vw;
        }

        .ant-drawer-content {
            height: 100dvh;
            overflow: hidden;
            border-radius: 0 !important;
        }

        .ant-drawer-body {
            height: 100%;
            padding: 0;
            overflow: hidden;
            border-radius: 0 !important;
        }

        @supports not (height: 100dvh) {
            .ant-drawer-content {
                height: 100vh;
            }
        }
    `,
    editorToolbar: css`
        position: absolute;
        top: 12px;
        left: 12px;
        z-index: 20;
        display: flex;
        align-items: center;
        min-width: 0;
        max-width: calc(100vw - 68px);

        .actions {
            display: flex;
            align-items: center;
            gap: 6px;
            min-width: 0;
            padding: 5px;
            border-radius: ${token.borderRadiusLG}px;
            background: ${isDarkMode ? "rgba(30,30,30,0.84)" : "rgba(255,255,255,0.9)"};
            box-shadow: ${isDarkMode ? "0 8px 22px rgba(0,0,0,0.3)" : "0 8px 22px rgba(15,23,42,0.08)"};
            backdrop-filter: blur(10px);
        }

        @media (max-width: 720px) {
            .actions {
                max-width: calc(100vw - 68px);
                overflow-x: auto;
                scrollbar-width: none;
            }

            .actions::-webkit-scrollbar {
                display: none;
            }
        }
    `,
    editorToolButton: css`
        display: inline-flex;
        align-items: center;
        justify-content: center;
        flex: 0 0 auto;
        width: 30px;
        height: 30px;
        padding: 0;
        border: 0;
        border-radius: ${token.borderRadius}px;
        outline: none;
        appearance: none;
        color: ${token.colorTextSecondary};
        background: transparent;
        cursor: pointer;
        transition: color 0.16s ease, background 0.16s ease, transform 0.16s ease;

        .anticon {
            font-size: 14px;
            line-height: 1;
        }

        &:hover,
        &:focus-visible {
            color: ${token.colorPrimary};
            background: ${isDarkMode ? "rgba(255,255,255,0.08)" : token.colorFillTertiary};
        }

        &:active {
            transform: scale(0.96);
        }
    `,
    editorClose: css`
        position: absolute;
        top: 12px;
        right: 12px;
        z-index: 20;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        width: 32px;
        height: 32px;
        padding: 0;
        border: 0;
        border-radius: 50%;
        outline: none;
        appearance: none;
        color: ${isDarkMode ? "rgba(255,255,255,0.86)" : token.colorTextSecondary};
        background: ${isDarkMode ? "rgba(42,42,42,0.86)" : "rgba(255,255,255,0.86)"};
        box-shadow: ${isDarkMode ? "0 8px 22px rgba(0,0,0,0.32)" : "0 8px 22px rgba(15,23,42,0.08)"};
        backdrop-filter: blur(10px);
        cursor: pointer;
        transition: color 0.16s ease, background 0.16s ease, transform 0.16s ease;

        .anticon {
            font-size: 14px;
            line-height: 1;
        }

        &:hover,
        &:focus-visible {
            color: ${token.colorError};
            background: ${isDarkMode ? "rgba(255,77,79,0.16)" : token.colorErrorBg};
        }

        &:active {
            transform: scale(0.96);
        }
    `,
    editorFullscreen: css`
        width: 100vw;
        height: 100dvh !important;
        min-height: 0 !important;
        border: 0;
        border-radius: 0;
        box-shadow: none;

        @supports not (height: 100dvh) {
            height: 100vh !important;
        }
    `,
    canvas: css`
        position: absolute;
        inset: 0;
        width: 100%;
        height: 100%;
        min-width: 0;
        min-height: 0;
        background: ${isDarkMode ? "rgba(255,255,255,0.01)" : "#fbfcfe"};
    `,
    versionDrawer: css`
        .ant-drawer-content-wrapper {
            width: min(420px, 100vw) !important;
            max-width: 100vw;
            transition-property: transform;
        }

        .ant-drawer-section {
            overflow: hidden;
            border-radius: ${token.borderRadiusLG}px 0 0 ${token.borderRadiusLG}px;
        }

        .ant-drawer-header {
            min-height: 56px;
            padding: 14px 18px;
            border-bottom-color: ${token.colorBorderSecondary};
        }

        .ant-drawer-title {
            color: ${token.colorTextHeading};
            font-size: 15px;
            font-weight: 600;
        }

        .ant-drawer-body {
            padding: 0;
            overflow-y: auto;
        }

        @media (max-width: 640px) {
            .ant-drawer-section {
                border-radius: 0;
            }

            .ant-drawer-header {
                padding-inline: 16px;
            }

            .ant-drawer-body {
                padding-bottom: env(safe-area-inset-bottom);
            }
        }
    `,
    versionList: css`
        min-height: 100%;
        background: ${token.colorBgContainer};

        .version-filter {
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 14px 18px 12px;
        }

        .version-filter .ant-segmented {
            width: auto;
            padding: 2px;
            background: ${token.colorFillTertiary};
        }

        .version-filter .ant-segmented-item-label {
            min-height: 28px;
            padding-inline: 12px;
            font-size: 12px;
            line-height: 28px;
        }

        .timeline {
            display: flex;
            flex-direction: column;
            margin: 0;
            padding: 2px 18px 18px;
            list-style: none;
        }

        .version {
            position: relative;
            display: grid;
            grid-template-columns: 24px minmax(0, 1fr);
            gap: 10px;
            min-width: 0;
            padding-bottom: 10px;
        }

        .version:last-child {
            padding-bottom: 0;
        }

        .version:last-child .line {
            display: none;
        }

        .version.current .dot {
            background: ${token.colorPrimary};
        }

        .axis {
            position: relative;
            display: flex;
            align-items: flex-start;
            align-self: stretch;
            justify-content: center;
            padding-top: 15px;
        }

        .dot {
            position: relative;
            z-index: 1;
            width: 9px;
            height: 9px;
            border: 2px solid ${token.colorPrimary};
            border-radius: 50%;
            background: ${token.colorBgContainer};
        }

        .line {
            position: absolute;
            top: 0;
            bottom: -12px;
            left: 50%;
            width: 1px;
            transform: translateX(-50%);
            background: ${token.colorBorderSecondary};
        }

        .body {
            position: relative;
            box-sizing: border-box;
            width: 100%;
            min-width: 0;
            padding: 10px 12px;
            border-radius: ${token.borderRadius}px;
            background: ${token.colorFillQuaternary};
            cursor: pointer;
            transition: background ${token.motionDurationSlow};
        }

        .body::before {
            position: absolute;
            top: 14px;
            left: -6px;
            width: 0;
            height: 0;
            border-top: 6px solid transparent;
            border-right: 6px solid ${token.colorFillQuaternary};
            border-bottom: 6px solid transparent;
            content: "";
            transition: border-right-color ${token.motionDurationSlow};
        }

        @media (hover: hover) and (pointer: fine) {
            .body:hover {
                background: ${token.colorFillSecondary};
            }

            .body:hover::before {
                border-right-color: ${token.colorFillSecondary};
            }
        }

        .head {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 8px;
            min-width: 0;
            min-height: 20px;
            margin-bottom: 6px;
        }

        .meta {
            display: flex;
            flex: 1 1 auto;
            align-items: center;
            flex-wrap: nowrap;
            gap: 7px;
            min-width: 0;
            overflow: hidden;
        }

        .snapshot-version {
            flex: none;
            padding-inline: 6px;
            border-radius: ${token.borderRadiusSM}px;
            background: ${token.colorPrimaryBg};
            color: ${token.colorPrimary};
            font-size: 11px;
            font-weight: 600;
            line-height: 20px;
        }

        .current-label {
            flex: none;
            color: ${token.colorPrimary};
            font-size: 11px;
            line-height: 20px;
            white-space: nowrap;
        }

        .snapshot-time {
            flex: 0 1 auto;
            min-width: 0;
            max-width: 45%;
            overflow: hidden;
            color: ${token.colorTextTertiary};
            font-size: 11px;
            line-height: 18px;
            font-variant-numeric: tabular-nums;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .snapshot-text {
            min-height: 18px;
            overflow: hidden;
            color: ${token.colorTextSecondary};
            font-size: 12px;
            line-height: 18px;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .snapshot-extra {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 8px;
            min-height: 28px;
            margin-top: 6px;
        }

        .snapshot-stats {
            display: flex;
            flex: 1 1 auto;
            align-items: center;
            flex-wrap: nowrap;
            gap: 5px 10px;
            min-width: 0;
            overflow: hidden;
            color: ${token.colorTextTertiary};
            font-size: 11px;
            line-height: 18px;
            white-space: nowrap;
        }

        .preview-button.ant-btn {
            flex: none;
            height: 28px;
            min-width: 48px;
            padding-inline: 4px;
            font-size: 12px;
        }

        .preview-button.ant-btn:focus-visible {
            outline: 1px solid ${token.colorPrimaryBorder};
            outline-offset: -2px;
            box-shadow: none;
        }

        .empty {
            display: flex;
            align-items: center;
            justify-content: center;
            min-height: 240px;
            padding: 24px;
        }

        @media (max-width: 640px) {
            .version-filter {
                padding-inline: 16px;
            }

            .timeline {
                padding-inline: 16px;
            }

            .version {
                grid-template-columns: 18px minmax(0, 1fr);
                gap: 8px;
            }
        }
    `,
    versionPreviewModal: css`
        .ant-modal {
            width: min(900px, calc(100vw - 32px)) !important;
            max-width: calc(100vw - 32px);
        }

        @media (max-width: 640px), (max-width: 960px) and (max-height: 520px) {
            .ant-modal-wrap {
                box-sizing: border-box;
                overflow: hidden;
                padding:
                    max(16px, env(safe-area-inset-top))
                    max(12px, env(safe-area-inset-right))
                    max(16px, env(safe-area-inset-bottom))
                    max(12px, env(safe-area-inset-left));
            }

            .ant-modal {
                top: 0;
                width: 100% !important;
                max-width: 900px;
                margin: 0 auto;
                padding: 0;
            }

            .ant-modal-container {
                display: flex;
                flex-direction: column;
                width: 100%;
                height: min(
                    620px,
                    calc(
                        100vh
                        - max(16px, env(safe-area-inset-top))
                        - max(16px, env(safe-area-inset-bottom))
                    )
                );
                height: min(
                    620px,
                    calc(
                        100dvh
                        - max(16px, env(safe-area-inset-top))
                        - max(16px, env(safe-area-inset-bottom))
                    )
                );
                max-height: 100%;
                min-height: 0;
                border-radius: ${token.borderRadiusLG}px;
                overflow: hidden;
            }

            .ant-modal-header {
                flex: none;
            }

            .ant-modal-body {
                display: flex;
                flex: 1;
                min-width: 0;
                min-height: 0;
                overflow: hidden;
            }

            .ant-modal-footer {
                flex: none;
            }
        }
    `,
    versionPreview: css`
        display: flex;
        flex-direction: column;
        width: 100%;
        min-width: 0;

        .preview-meta {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 16px;
            margin-bottom: 12px;
        }

        .preview-remark {
            min-width: 0;
            max-width: 100%;
            overflow: hidden;
            color: ${token.colorTextHeading};
            font-size: 13px;
            font-weight: 500;
            line-height: 20px;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .preview-details {
            display: flex;
            align-items: center;
            justify-content: flex-end;
            flex-wrap: wrap;
            gap: 4px 12px;
            color: ${token.colorTextTertiary};
            font-size: 11px;
            line-height: 20px;
            text-align: right;
        }

        .preview-canvas {
            height: clamp(300px, 56dvh, 520px);
            min-height: 0;
            overflow: hidden;
            border: 1px solid ${token.colorBorderSecondary};
            border-radius: ${token.borderRadiusLG}px;
            background: ${token.colorBgLayout};
        }

        @media (max-width: 640px), (max-width: 960px) and (max-height: 520px) {
            flex: 1;
            min-height: 0;

            .preview-meta {
                flex-direction: column;
                align-items: flex-start;
                gap: 5px;
                margin-bottom: 10px;
            }

            .preview-details {
                width: 100%;
                max-width: 100%;
                justify-content: flex-start;
                text-align: left;
            }

            .preview-canvas {
                flex: 1;
                height: auto;
                min-height: 0;
            }
        }

        @media (min-width: 480px) and (max-width: 960px) and (max-height: 520px) {
            .preview-meta {
                flex-direction: row;
                align-items: center;
                gap: 16px;
                margin-bottom: 8px;
            }

            .preview-details {
                width: auto;
                justify-content: flex-end;
                text-align: right;
            }
        }
    `,
    templateGrid: css`
        display: grid;
        grid-template-columns: repeat(3, 1fr);
        gap: 8px;

        @media (max-width: 640px) {
            grid-template-columns: 1fr;
        }

        .template {
            padding: 10px;
            border: 1px solid ${token.colorBorderSecondary};
            border-radius: ${token.borderRadius}px;
            background: ${token.colorBgContainer};
            cursor: pointer;
        }

        .template:hover,
        .template.active {
            border-color: color-mix(in srgb, var(--accent), ${token.colorBorderSecondary} 42%);
            background: ${isDarkMode ? "rgba(255,255,255,0.04)" : "color-mix(in srgb, var(--accent), #ffffff 94%)"};
        }

        .name {
            color: ${token.colorTextHeading};
            font-size: 13px;
            font-weight: 600;
            line-height: 20px;
        }

        .desc {
            color: ${token.colorTextTertiary};
            font-size: 12px;
            line-height: 18px;
        }
    `,
}));

const MindPage: React.FC = () => {
    const {styles} = useStyles();
    const {message, modal} = App.useApp();
    const pageRef = useRef<HTMLDivElement | null>(null);
    const categoryScrollRef = useRef<HTMLDivElement | null>(null);
    const editorRef = useRef<MindRef | null>(null);
    const versionPreviewRef = useRef<MindRef | null>(null);
    const versionPreviewCanvasRef = useRef<HTMLDivElement | null>(null);
    const versionSequenceRef = useRef<Record<string, number>>(Object.fromEntries(
        initialDocuments.map((record) => [record.id, record.version]),
    ));
    const lastVersionSignatureRef = useRef<Record<string, string>>(Object.fromEntries(
        initialDocuments.map((record) => [record.id, getMindValueSignature(record.value)]),
    ));
    const [form] = Form.useForm<CreateMindForm>();
    const [categoryForm] = Form.useForm<CategoryForm>();
    const selectedTemplateId = Form.useWatch("templateId", form);
    const [categories, setCategories] = useState<MindCategory[]>(initialCategories);
    const [documents, setDocuments] = useState<MindDocument[]>(initialDocuments);
    const [activeCategoryId, setActiveCategoryId] = useState(ALL_CATEGORY_ID);
    const [keyword, setKeyword] = useState("");
    const [sortType, setSortType] = useState<MindSortType>("update-desc");
    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(12);
    const [selectedId, setSelectedId] = useState<string>();
    const [editorValue, setEditorValue] = useState<MindValue>(initialDocuments[0].value);
    const [editorRevision, setEditorRevision] = useState(0);
    const [createOpen, setCreateOpen] = useState(false);
    const [categoryOpen, setCategoryOpen] = useState(false);
    const [categoryRecord, setCategoryRecord] = useState<MindCategory>();
    const [versionOpen, setVersionOpen] = useState(false);
    const [versionFilter, setVersionFilter] = useState<MindVersionSaveType>("manual");
    const [versionPreview, setVersionPreview] = useState<MindVersionRecord>();
    const [versions, setVersions] = useState<Record<string, MindVersionRecord[]>>(initialVersions);
    const [pageHeight, setPageHeight] = useState<number>();
    const [categoryScrollState, setCategoryScrollState] = useState({canLeft: false, canRight: false});

    const selectedDocument = documents.find((item) => item.id === selectedId);
    const versionItems = selectedId ? versions[selectedId] || [] : [];
    const filteredVersionItems = versionItems.filter((item) => item.saveType === versionFilter);
    const currentVersionRecord = useMemo(() => {
        if (!selectedDocument) {
            return undefined;
        }
        const currentSignature = getMindValueSignature(editorValue);
        return versionItems.find((item) => item.version === selectedDocument.version
            && getMindValueSignature(item.value) === currentSignature);
    }, [editorValue, selectedDocument, versionItems]);

    useEffect(() => {
        const element = versionPreviewCanvasRef.current;
        if (!versionPreview || !element || typeof ResizeObserver === "undefined") {
            return;
        }
        let frame = 0;
        const observer = new ResizeObserver(() => {
            cancelAnimationFrame(frame);
            frame = requestAnimationFrame(() => versionPreviewRef.current?.fitView());
        });
        observer.observe(element);
        return () => {
            cancelAnimationFrame(frame);
            observer.disconnect();
        };
    }, [versionPreview]);

    const takeNextVersion = (mindId: string, fallback = 0) => {
        const nextVersion = Math.max(versionSequenceRef.current[mindId] || 0, fallback) + 1;
        versionSequenceRef.current[mindId] = nextVersion;
        return nextVersion;
    };

    useEffect(() => {
        // 按页面真实位置测量当前视口可用高度，避免 Header、面包屑高度变化导致左侧滚动区失准。
        let frame = 0;
        const updateHeight = () => {
            cancelAnimationFrame(frame);
            frame = requestAnimationFrame(() => {
                const page = pageRef.current;
                if (!page) return;
                const top = page.getBoundingClientRect().top;
                const nextHeight = Math.max(320, Math.floor(window.innerHeight - top - 18));
                setPageHeight((value) => value === nextHeight ? value : nextHeight);
            });
        };
        updateHeight();
        window.addEventListener("resize", updateHeight);
        window.addEventListener("orientationchange", updateHeight);
        return () => {
            cancelAnimationFrame(frame);
            window.removeEventListener("resize", updateHeight);
            window.removeEventListener("orientationchange", updateHeight);
        };
    }, [selectedId]);

    useEffect(() => {
        const element = categoryScrollRef.current;
        if (!element) return;

        const updateScrollState = () => {
            const maxScrollLeft = Math.max(0, element.scrollWidth - element.clientWidth);
            const nextState = {
                canLeft: element.scrollLeft > 1,
                canRight: element.scrollLeft < maxScrollLeft - 1,
            };
            setCategoryScrollState((current) => current.canLeft === nextState.canLeft
                && current.canRight === nextState.canRight ? current : nextState);
        };
        const handleWheel = (event: WheelEvent) => {
            const maxScrollLeft = Math.max(0, element.scrollWidth - element.clientWidth);
            if (maxScrollLeft <= 0) return;
            const delta = Math.abs(event.deltaX) > Math.abs(event.deltaY) ? event.deltaX : event.deltaY;
            const canScroll = delta < 0 ? element.scrollLeft > 0 : element.scrollLeft < maxScrollLeft;
            if (!canScroll) return;
            event.preventDefault();
            element.scrollLeft += delta;
        };

        updateScrollState();
        element.addEventListener("scroll", updateScrollState, {passive: true});
        element.addEventListener("wheel", handleWheel, {passive: false});
        const observer = new ResizeObserver(updateScrollState);
        observer.observe(element);
        if (element.firstElementChild) {
            observer.observe(element.firstElementChild);
        }
        return () => {
            element.removeEventListener("scroll", updateScrollState);
            element.removeEventListener("wheel", handleWheel);
            observer.disconnect();
        };
    }, [categories]);

    const scrollCategories = (direction: -1 | 1) => {
        const element = categoryScrollRef.current;
        if (!element) return;
        element.scrollBy({
            left: direction * Math.max(180, Math.round(element.clientWidth * 0.68)),
            behavior: "smooth",
        });
    };

    const filteredDocuments = useMemo(() => {
        const search = keyword.trim().toLowerCase();
        return documents
            .filter((item) => activeCategoryId === ALL_CATEGORY_ID || item.categoryId === activeCategoryId)
            .filter((item) => !search || `${item.title} ${item.description} ${item.tags.join(" ")}`.toLowerCase().includes(search))
            .sort((a, b) => {
                const field = sortType === "create-desc" ? "createTime" : "updateTime";
                return new Date(b[field].replace(/-/g, "/")).getTime() - new Date(a[field].replace(/-/g, "/")).getTime();
            });
    }, [activeCategoryId, documents, keyword, sortType]);

    const activeSortLabel = mindSortOptions.find((item) => item.value === sortType)?.label || "最近更新";
    const pagedDocuments = useMemo(() => {
        const start = (page - 1) * pageSize;
        return filteredDocuments.slice(start, start + pageSize);
    }, [filteredDocuments, page, pageSize]);

    useEffect(() => {
        setPage(1);
    }, [activeCategoryId, keyword, sortType]);

    useEffect(() => {
        const lastPage = Math.max(1, Math.ceil(filteredDocuments.length / pageSize));
        setPage((current) => Math.min(current, lastPage));
    }, [filteredDocuments.length, pageSize]);

    const openDocument = (record: MindDocument) => {
        setSelectedId(record.id);
        setEditorValue(cloneValue(record.value));
        setEditorRevision((value) => value + 1);
    };

    const openCreate = () => {
        form.setFieldsValue({
            title: "",
            description: "",
            categoryId: activeCategoryId === ALL_CATEGORY_ID ? (categories[0]?.id || "") : activeCategoryId,
            templateId: "blank",
        });
        setCreateOpen(true);
    };

    const openCreateCategory = () => {
        setCategoryRecord(undefined);
        categoryForm.setFieldsValue({
            name: "",
        });
        setCategoryOpen(true);
    };

    const openEditCategory = (category: MindCategory) => {
        setCategoryRecord(category);
        categoryForm.setFieldsValue({
            name: category.name,
        });
        setCategoryOpen(true);
    };

    const saveCategory = async () => {
        const values = await categoryForm.validateFields();
        if (categoryRecord) {
            setCategories((items) => items.map((item) => item.id === categoryRecord.id ? {
                ...item,
                name: values.name,
            } : item));
            setCategoryOpen(false);
            message.success("分类已更新");
            return;
        }
        const category: MindCategory = {
            id: `category-${Date.now().toString(36)}`,
            name: values.name,
        };
        setCategories((items) => [...items, category]);
        setActiveCategoryId(category.id);
        setCategoryOpen(false);
        message.success("分类已创建");
    };

    const deleteCategory = (category: MindCategory) => {
        if (categories.length <= 1) {
            message.warning("至少保留一个分类");
            return;
        }
        const target = categories.find((item) => item.id !== category.id);
        if (!target) {
            return;
        }
        modal.confirm({
            title: "删除分类",
            content: `删除「${category.name}」后，该分类下的导图会移动到「${target.name}」。`,
            okText: "删除",
            okButtonProps: {danger: true},
            onOk: () => {
                setCategories((items) => items.filter((item) => item.id !== category.id));
                setDocuments((items) => items.map((item) => item.categoryId === category.id ? {
                    ...item,
                    categoryId: target.id,
                } : item));
                if (activeCategoryId === category.id) {
                    setActiveCategoryId(target.id);
                }
                setCategoryOpen(false);
                message.success("分类已删除");
            },
        });
    };

    const createDocument = async () => {
        const values = await form.validateFields();
        const template = templateItems.find((item) => item.id === values.templateId) || templateItems[0];
        const now = formatDateTime(new Date());
        const value = cloneValue(template.value);
        value.nodes = value.nodes.map((node) => node.root ? {...node, title: values.title, color: template.color} : node);
        const record = withStats({
            id: `mind-${Date.now().toString(36)}`,
            categoryId: values.categoryId,
            title: values.title,
            description: values.description || "",
            tags: [categories.find((item) => item.id === values.categoryId)?.name || "导图"],
            value,
            version: 1,
            createTime: now,
            updateTime: now,
        });
        setDocuments((items) => [record, ...items]);
        setVersions((items) => ({
            ...items,
            [record.id]: [createVersionRecord(record, "初始版本")],
        }));
        versionSequenceRef.current[record.id] = record.version;
        lastVersionSignatureRef.current[record.id] = getMindValueSignature(record.value);
        setCreateOpen(false);
        openDocument(record);
        message.success("导图已创建");
    };

    const persistDocument = async (value?: MindValue, stats?: MindStats, reason: MindSaveReason = "manual") => {
        if (!selectedId || !selectedDocument) {
            return;
        }
        const nextValue = value || editorRef.current?.getValue() || editorValue;
        const nextStats = stats || getMindStats(nextValue);
        const nowDate = new Date();
        const now = formatDateTime(nowDate);
        const saveType: MindVersionSaveType = reason === "auto" ? "auto" : "manual";
        const signature = getMindValueSignature(nextValue);
        const shouldWriteVersion = saveType === "manual"
            || lastVersionSignatureRef.current[selectedId] !== signature;
        const latestVersion = versionItems[0];
        const rollingAutoVersion = saveType === "auto"
            && latestVersion?.saveType === "auto"
            && latestVersion.remark === "自动保存"
            && nowDate.getTime() - getDateTimeValue(latestVersion.createTime) < AUTO_VERSION_SNAPSHOT_INTERVAL
            ? latestVersion
            : undefined;
        const shouldAppendVersion = shouldWriteVersion && !rollingAutoVersion;
        const nextVersion = shouldAppendVersion
            ? takeNextVersion(selectedId, selectedDocument.version)
            : (rollingAutoVersion?.version
                || Math.max(versionSequenceRef.current[selectedId] || 0, selectedDocument.version));

        setDocuments((items) => items.map((item) => item.id === selectedId ? {
            ...item,
            value: cloneValue(nextValue),
            nodeCount: nextStats.nodeCount,
            depth: nextStats.depth,
            version: nextVersion,
            updateTime: now,
        } : item));
        setEditorValue(cloneValue(nextValue));

        if (shouldWriteVersion) {
            const nextRecord: MindVersionRecord = {
                id: rollingAutoVersion?.id
                    || `${selectedId}-${saveType}-v${nextVersion}-${Date.now().toString(36)}`,
                mindId: selectedId,
                version: nextVersion,
                saveType,
                remark: reason === "auto" ? "自动保存" : (reason === "shortcut" ? "快捷键保存" : "手动保存"),
                nodeCount: nextStats.nodeCount,
                depth: nextStats.depth,
                value: cloneValue(nextValue),
                createTime: now,
            };
            setVersions((items) => ({
                ...items,
                [selectedId]: rollingAutoVersion
                    ? (items[selectedId] || []).map((item) => item.id === rollingAutoVersion.id ? nextRecord : item)
                    : limitVersionRecords([nextRecord, ...(items[selectedId] || [])]),
            }));
            lastVersionSignatureRef.current[selectedId] = signature;
        }

        if (reason !== "auto") {
            message.success("导图已保存");
        }
    };

    const duplicateDocument = (record: MindDocument) => {
        const now = formatDateTime(new Date());
        const copy = {
            ...record,
            id: `mind-${Date.now().toString(36)}`,
            title: `${record.title} 副本`,
            value: cloneValue(record.value),
            shared: false,
            version: 1,
            createTime: now,
            updateTime: now,
        };
        setDocuments((items) => [copy, ...items]);
        setVersions((items) => ({
            ...items,
            [copy.id]: [createVersionRecord(copy, "复制创建")],
        }));
        versionSequenceRef.current[copy.id] = copy.version;
        lastVersionSignatureRef.current[copy.id] = getMindValueSignature(copy.value);
        message.success("导图已复制");
    };

    const deleteDocument = (record: MindDocument) => {
        modal.confirm({
            title: "删除导图",
            content: `确定删除「${record.title}」吗？`,
            okText: "删除",
            okButtonProps: {danger: true},
            onOk: () => {
                setDocuments((items) => items.filter((item) => item.id !== record.id));
                setVersions((items) => {
                    const nextItems = {...items};
                    delete nextItems[record.id];
                    return nextItems;
                });
                delete versionSequenceRef.current[record.id];
                delete lastVersionSignatureRef.current[record.id];
                if (record.id === selectedId) {
                    setSelectedId(undefined);
                    setVersionOpen(false);
                    setVersionPreview(undefined);
                }
                message.success("导图已删除");
            },
        });
    };

    const toggleShare = (record: MindDocument) => {
        setDocuments((items) => items.map((item) => item.id === record.id ? {...item, shared: !item.shared} : item));
        message.success(record.shared ? "分享已关闭" : "分享已开启");
    };

    const saveEditor = () => {
        const task = editorRef.current?.save("manual");
        if (task) {
            void task.catch((error) => message.error(error instanceof Error ? error.message : "导图保存失败"));
        }
    };

    const arrangeEditor = () => {
        editorRef.current?.layout();
    };

    const exportEditorImage = async (type: "png" | "jpg") => {
        const editor = editorRef.current;
        if (!editor || !selectedDocument) {
            return;
        }
        const dataUrl = type === "png" ? await editor.exportPng() : await editor.exportJpg();
        downloadDataUrl(dataUrl, `${selectedDocument.title}.${type}`);
        message.success(`${type.toUpperCase()} 已导出`);
    };

    const copyEditorImage = () => {
        const task = editorRef.current?.copyImage();
        if (task) {
            void task.catch((error) => message.error(error instanceof Error ? error.message : "复制图片失败"));
        }
    };

    const restoreVersion = (record: MindVersionRecord) => {
        if (!selectedDocument) {
            return;
        }
        const now = formatDateTime(new Date());
        const currentValue = cloneValue(editorRef.current?.getValue() || selectedDocument.value);
        const currentStats = getMindStats(currentValue);
        const nextValue = cloneValue(record.value);
        const nextStats = getMindStats(nextValue);
        const currentSignature = getMindValueSignature(currentValue);
        const targetSignature = getMindValueSignature(nextValue);
        if (currentSignature === targetSignature) {
            message.info("当前内容已与该版本一致");
            return;
        }
        const latestVersion = versionItems[0];
        const shouldBackupCurrent = !latestVersion
            || getMindValueSignature(latestVersion.value) !== currentSignature;
        const backupVersion = shouldBackupCurrent
            ? takeNextVersion(selectedDocument.id, selectedDocument.version)
            : selectedDocument.version;
        const restoredVersion = takeNextVersion(selectedDocument.id, backupVersion);
        const backupRecord: MindVersionRecord | undefined = shouldBackupCurrent ? {
            id: `${selectedDocument.id}-auto-v${backupVersion}-${Date.now().toString(36)}`,
            mindId: selectedDocument.id,
            version: backupVersion,
            saveType: "auto",
            remark: "恢复前备份",
            nodeCount: currentStats.nodeCount,
            depth: currentStats.depth,
            value: currentValue,
            createTime: now,
        } : undefined;
        const restoredRecord: MindVersionRecord = {
            id: `${selectedDocument.id}-manual-v${restoredVersion}-${Date.now().toString(36)}`,
            mindId: selectedDocument.id,
            version: restoredVersion,
            saveType: "manual",
            remark: `从 v${record.version} 恢复`,
            nodeCount: nextStats.nodeCount,
            depth: nextStats.depth,
            value: nextValue,
            createTime: now,
        };

        setDocuments((items) => items.map((item) => item.id === selectedDocument.id ? {
            ...item,
            value: nextValue,
            nodeCount: nextStats.nodeCount,
            depth: nextStats.depth,
            version: restoredVersion,
            updateTime: now,
        } : item));
        setVersions((items) => ({
            ...items,
            [selectedDocument.id]: limitVersionRecords([
                restoredRecord,
                ...(backupRecord ? [backupRecord] : []),
                ...(items[selectedDocument.id] || []),
            ]),
        }));
        lastVersionSignatureRef.current[selectedDocument.id] = getMindValueSignature(nextValue);
        setEditorValue(cloneValue(nextValue));
        setEditorRevision((value) => value + 1);
        setVersionPreview(undefined);
        setVersionOpen(false);
        message.success(backupRecord ? "版本已恢复，恢复前内容已自动备份" : "版本已恢复");
    };

    const confirmRestoreVersion = (record: MindVersionRecord) => {
        if (!selectedDocument) {
            return;
        }
        const currentValue = editorRef.current?.getValue() || selectedDocument.value;
        const currentSignature = getMindValueSignature(currentValue);
        if (currentSignature === getMindValueSignature(record.value)) {
            message.info("当前内容已与该版本一致");
            return;
        }
        const latestVersion = versionItems[0];
        const shouldBackupCurrent = !latestVersion
            || getMindValueSignature(latestVersion.value) !== currentSignature;
        modal.confirm({
            title: `恢复到 v${record.version}？`,
            content: shouldBackupCurrent
                ? `恢复前会自动备份当前内容，然后使用 ${record.createTime} 保存的版本替换画布。`
                : `当前内容已保存在版本记录中，将使用 ${record.createTime} 保存的版本替换画布。`,
            okText: "确认",
            cancelText: "取消",
            onOk: () => restoreVersion(record),
        });
    };

    const documentMenu = (record: MindDocument) => [
        {key: "open", icon: <EditOutlined/>, label: "打开导图", onClick: ({domEvent}: any) => {
                domEvent.stopPropagation();
                openDocument(record);
            }},
        {key: "copy", icon: <CopyOutlined/>, label: "复制导图", onClick: ({domEvent}: any) => {
                domEvent.stopPropagation();
                duplicateDocument(record);
            }},
        {key: "share", icon: <ShareAltOutlined/>, label: record.shared ? "关闭分享" : "开启分享", onClick: ({domEvent}: any) => {
                domEvent.stopPropagation();
                toggleShare(record);
            }},
        {type: "divider" as const},
        {key: "delete", icon: <DeleteOutlined/>, label: "删除导图", danger: true, onClick: ({domEvent}: any) => {
                domEvent.stopPropagation();
                deleteDocument(record);
            }},
    ];

    const categoryItems = [
        {id: ALL_CATEGORY_ID, name: "全部导图", accent: getStableColor(ALL_CATEGORY_ID)},
        ...categories.map((item) => ({
            ...item,
            accent: getStableColor(item.id || item.name),
        })),
    ];
    const editorExportItems = [
        {key: "png", icon: <FileImageOutlined/>, label: "导出 PNG", onClick: () => {
                void exportEditorImage("png").catch((error) => message.error(error instanceof Error ? error.message : "PNG 导出失败"));
            }},
        {key: "jpg", icon: <FileImageOutlined/>, label: "导出 JPG", onClick: () => {
                void exportEditorImage("jpg").catch((error) => message.error(error instanceof Error ? error.message : "JPG 导出失败"));
            }},
        {key: "json", icon: <FileTextOutlined/>, label: "导出 JSON", onClick: () => editorRef.current?.downloadJson()},
        {key: "markdown", icon: <FileMarkdownOutlined/>, label: "导出 Markdown", onClick: () => editorRef.current?.downloadMarkdown()},
        {type: "divider" as const},
        {key: "copy", icon: <CopyOutlined/>, label: "复制图片", onClick: copyEditorImage},
    ];
    const pageStyle = pageHeight ? {"--mind-page-height": `${pageHeight}px`} as React.CSSProperties : undefined;
    const openVersionHistory = () => {
        setVersionFilter(currentVersionRecord?.saveType || "manual");
        setVersionPreview(undefined);
        setVersionOpen(true);
    };
    const closeEditor = () => {
        setSelectedId(undefined);
        setVersionOpen(false);
        setVersionPreview(undefined);
    };
    const editorContent = selectedDocument ? (
        <div className={`${styles.editor} ${styles.editorFullscreen}`}>
            <div className={styles.editorToolbar}>
                <div className="actions">
                    <Tooltip title="保存">
                        <button className={styles.editorToolButton} type="button" onClick={saveEditor}>
                            <SaveOutlined/>
                        </button>
                    </Tooltip>
                    <Tooltip title="整理布局">
                        <button className={styles.editorToolButton} type="button" onClick={arrangeEditor}>
                            <ApartmentOutlined/>
                        </button>
                    </Tooltip>
                    <Dropdown trigger={["click"]} placement="bottomRight" menu={{items: editorExportItems as any}}>
                        <button className={styles.editorToolButton} type="button" title="导出">
                            <DownloadOutlined/>
                        </button>
                    </Dropdown>
                    <Tooltip title="版本记录">
                        <button
                            className={styles.editorToolButton}
                            type="button"
                            aria-label="版本记录"
                            onClick={openVersionHistory}
                        >
                            <HistoryOutlined/>
                        </button>
                    </Tooltip>
                    <Tooltip title={selectedDocument.shared ? "关闭分享" : "开启分享"}>
                        <button className={styles.editorToolButton} type="button" onClick={() => toggleShare(selectedDocument as any)}>
                            <ShareAltOutlined/>
                        </button>
                    </Tooltip>
                </div>
            </div>
            <button className={styles.editorClose} type="button" onClick={closeEditor}>
                <CloseOutlined/>
            </button>
            <div className={styles.canvas}>
                <Mind
                    key={`${selectedDocument.id}-${editorRevision}`}
                    ref={editorRef}
                    defaultValue={editorValue}
                    height="100%"
                    autoSave
                    autoSaveDelay={4000}
                    onDataChange={(value) => {
                        setEditorValue(value);
                    }}
                    onSave={persistDocument}
                />
            </div>
        </div>
    ) : null;

    return (
        <Body space={false}>
            <div ref={pageRef} className={styles.page} style={pageStyle}>
                <div className={styles.library}>
                    <Card className={styles.libraryCard}>
                    <aside className={styles.categoryPanel} aria-label="导图分类">
                        <div className={styles.categoryInner}>
                        <div className={styles.categoryRail}>
                        <button
                            className={styles.categoryScrollButton}
                            type="button"
                            aria-label="向左浏览分类"
                            disabled={!categoryScrollState.canLeft}
                            onClick={() => scrollCategories(-1)}
                        >
                            <LeftOutlined/>
                        </button>
                        <div ref={categoryScrollRef} className={styles.categoryScroll}>
                            <div className={styles.categoryGrid}>
                                {categoryItems.map((item) => (
                                    <div
                                        key={item.id}
                                        className={`${styles.categoryCard} ${activeCategoryId === item.id ? "active" : ""}`}
                                        style={{"--accent": item.accent} as React.CSSProperties}
                                        onClick={() => setActiveCategoryId(item.id)}
                                    >
                                        <div className="row">
                                            <span className="name-wrap">
                                                <span className="color"/>
                                                <span className="name">{item.name}</span>
                                            </span>
                                            <span className="right">
                                                {item.id !== ALL_CATEGORY_ID && (
                                                    <Dropdown
                                                        trigger={["click"]}
                                                        placement="bottomRight"
                                                        menu={{
                                                            items: [
                                                                {key: "edit", icon: <EditOutlined/>, label: "修改分类", onClick: ({domEvent}: any) => {
                                                                        domEvent.stopPropagation();
                                                                        openEditCategory(item);
                                                                    }},
                                                                {key: "delete", icon: <DeleteOutlined/>, label: "删除分类", danger: true, onClick: ({domEvent}: any) => {
                                                                        domEvent.stopPropagation();
                                                                        deleteCategory(item);
                                                                    }},
                                                            ] as any,
                                                        }}
                                                    >
                                                        <span className="action" onClick={(event) => event.stopPropagation()}>
                                                            <Button type="text" size="small" icon={<MoreOutlined/>}/>
                                                        </span>
                                                    </Dropdown>
                                                )}
                                            </span>
                                        </div>
                                    </div>
                                ))}
                                <button className={`${styles.categoryAddCard} mobile`} type="button" onClick={openCreateCategory}>
                                    <span className="icon"><PlusOutlined/></span>
                                    <span>添加分类</span>
                                </button>
                            </div>
                        </div>
                        <button
                            className={styles.categoryScrollButton}
                            type="button"
                            aria-label="向右浏览分类"
                            disabled={!categoryScrollState.canRight}
                            onClick={() => scrollCategories(1)}
                        >
                            <RightOutlined/>
                        </button>
                        </div>
                        <div className={styles.categoryFooter}>
                            <button className={styles.categoryAddCard} type="button" onClick={openCreateCategory}>
                                <span className="icon"><PlusOutlined/></span>
                                <span>添加分类</span>
                            </button>
                        </div>
                        </div>
                    </aside>

                    <div className={styles.content}>
                        <div className={styles.contentPanel}>
                        <section className={styles.createHero} aria-labelledby="mind-create-title">
                            <div className="hero-copy">
                                <div className="eyebrow">
                                    <span className="status-dot"/>
                                    <span>Mind Studio</span>
                                </div>
                                <div className="hero-title" id="mind-create-title">让灵感快速长成一张清晰的导图</div>
                                <div className="hero-desc">从空白画布自由创作，或选择模板快速梳理想法与计划。</div>
                            </div>
                            <div className="hero-visual">
                                <MindCreateHeroGraphic/>
                            </div>
                            <button className="create-button" type="button" onClick={openCreate}>
                                <span>新建导图</span>
                            </button>
                        </section>
                        <div className={styles.searchCard}>
                            <Input
                                allowClear
                                className={styles.searchBox}
                                prefix={<SearchOutlined/>}
                                placeholder="搜索标题、描述或标签"
                                value={keyword}
                                onChange={(event) => setKeyword(event.target.value)}
                                onPressEnter={(event) => setKeyword(event.currentTarget.value)}
                            />
                            <Dropdown
                                trigger={["click"]}
                                placement="bottomRight"
                                classNames={{root: styles.sortDropdown}}
                                menu={{
                                    selectable: true,
                                    selectedKeys: [sortType],
                                    items: mindSortOptions.map((item) => ({key: item.value, label: item.label})),
                                    onClick: ({key}) => setSortType(key as MindSortType),
                                }}
                            >
                                <button
                                    className={styles.sortTrigger}
                                    type="button"
                                    aria-label={`排序方式：${activeSortLabel}`}
                                >
                                    <SortDescendingOutlined className="marker"/>
                                    <span className="value">{activeSortLabel}</span>
                                    <DownOutlined className="arrow"/>
                                </button>
                            </Dropdown>
                        </div>
                        </div>
                        <Row>
                            <Col className={styles.documentDataColumn} span={24}>
                                <div className={`${styles.documentDataCard} ${pagedDocuments.length <= 0 ? "empty" : ""}`}>
                                    {pagedDocuments.length > 0 ? (
                                        <div className={styles.documentGrid}>
                                            {pagedDocuments.map((item) => (
                                                <Card
                                                    key={item.id}
                                                    className={styles.documentCard}
                                                    style={{"--accent": getStableColor(item.categoryId || item.id || item.title)} as React.CSSProperties}
                                                    styles={{body: {padding: 0}}}
                                                    onClick={() => openDocument(item)}
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
                                                                    openDocument(item);
                                                                }}
                                                            >
                                                                {item.title}
                                                            </button>
                                                            <span className="actions" onClick={(event) => event.stopPropagation()}>
                                                                <Dropdown trigger={["click"]} placement="bottomRight" menu={{items: documentMenu(item) as any}}>
                                                                    <Button type="text" size="small" icon={<MoreOutlined/>}/>
                                                                </Dropdown>
                                                            </span>
                                                        </div>
                                                        <div className="desc">{item.description || "暂无描述"}</div>
                                                        <div className="meta">
                                                            <span>{item.shared ? "已分享" : "仅自己可见"}</span>
                                                            <Tooltip title={item.updateTime}>
                                                                <span>{getTimeText(item.updateTime)}</span>
                                                            </Tooltip>
                                                        </div>
                                                    </div>
                                                </Card>
                                            ))}
                                        </div>
                                    ) : <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="暂无导图"/>}
                                </div>
                            </Col>
                            {filteredDocuments.length > 0 && (
                                <Col className={styles.documentPaginationColumn} span={24}>
                                    <div className={styles.documentPaginationCard}>
                                        <Pagination
                                            current={page}
                                            pageSize={pageSize}
                                            total={filteredDocuments.length}
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
                                            showTotal={(total, range) => `第 ${range[0]}-${range[1]} 条 / 共 ${total} 条`}
                                            onChange={(nextPage, nextPageSize) => {
                                                const sizeChanged = nextPageSize !== pageSize;
                                                setPageSize(nextPageSize);
                                                setPage(sizeChanged ? 1 : nextPage);
                                            }}
                                        />
                                    </div>
                                </Col>
                            )}
                        </Row>
                    </div>
                    </Card>
                </div>
            </div>

            <Drawer
                open={Boolean(selectedDocument)}
                rootClassName={styles.editorDrawer}
                placement="right"
                size="100%"
                closable={false}
                mask={false}
                destroyOnHidden
                onClose={closeEditor}
                afterOpenChange={(open) => {
                    if (!open) {
                        return;
                    }
                    requestAnimationFrame(() => {
                        requestAnimationFrame(() => editorRef.current?.fitView());
                    });
                }}
            >
                {editorContent}
            </Drawer>

            <Drawer
                title="版本记录"
                rootClassName={styles.versionDrawer}
                open={versionOpen}
                onClose={() => {
                    setVersionOpen(false);
                    setVersionPreview(undefined);
                }}
                destroyOnHidden
            >
                <div className={styles.versionList}>
                    <div className="version-filter">
                        <Segmented
                            value={versionFilter}
                            options={[
                                {label: "手动保存", value: "manual"},
                                {label: "自动保存", value: "auto"},
                            ]}
                            onChange={(value) => setVersionFilter(value as MindVersionSaveType)}
                        />
                    </div>

                    {filteredVersionItems.length <= 0 ? (
                        <div className="empty">
                            <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="暂无数据"/>
                        </div>
                    ) : (
                        <ol className="timeline" aria-label="历史版本">
                            {filteredVersionItems.map((item) => {
                                const isCurrent = item.id === currentVersionRecord?.id;
                                return (
                                    <li
                                        className={`version${isCurrent ? " current" : ""}`}
                                        key={item.id}
                                        aria-current={isCurrent ? "true" : undefined}
                                    >
                                        <div className="axis" aria-hidden="true">
                                            <span className="dot"/>
                                            <span className="line"/>
                                        </div>
                                        <div
                                            className="body"
                                            onClick={() => setVersionPreview(item)}
                                        >
                                            <div className="head">
                                                <div className="meta">
                                                    <span className="snapshot-version">v{item.version}</span>
                                                    {isCurrent && <span className="current-label">当前版本</span>}
                                                </div>
                                                <Tooltip title={item.createTime}>
                                                    <span className="snapshot-time">{getTimeText(item.createTime)}</span>
                                                </Tooltip>
                                            </div>
                                            <div className="snapshot-text">{item.remark}</div>
                                            <div className="snapshot-extra">
                                                <div className="snapshot-stats">
                                                    <span>{item.nodeCount} 个节点</span>
                                                    <span>{item.depth} 层</span>
                                                </div>
                                                <Button
                                                    className="preview-button"
                                                    type="link"
                                                    size="small"
                                                    icon={<EyeOutlined/>}
                                                    aria-label={`预览 v${item.version}`}
                                                    onClick={(event) => {
                                                        event.stopPropagation();
                                                        setVersionPreview(item);
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
                </div>
            </Drawer>

            <Modal
                title={versionPreview ? `预览 v${versionPreview.version}` : "版本预览"}
                rootClassName={styles.versionPreviewModal}
                open={Boolean(versionPreview)}
                width={900}
                centered
                okText="恢复"
                cancelText="返回"
                onOk={() => {
                    if (versionPreview) {
                        confirmRestoreVersion(versionPreview);
                    }
                }}
                onCancel={() => setVersionPreview(undefined)}
                afterOpenChange={(open) => {
                    if (!open) {
                        return;
                    }
                    requestAnimationFrame(() => {
                        requestAnimationFrame(() => versionPreviewRef.current?.fitView());
                    });
                }}
                destroyOnHidden
            >
                {versionPreview && (
                    <div className={styles.versionPreview}>
                        <div className="preview-meta">
                            <div className="preview-remark">{versionPreview.remark}</div>
                            <div className="preview-details">
                                <span>{versionPreview.createTime}</span>
                                <span>{versionPreview.nodeCount} 个节点</span>
                                <span>{versionPreview.depth} 层</span>
                            </div>
                        </div>
                        <div ref={versionPreviewCanvasRef} className="preview-canvas">
                            <Mind
                                key={versionPreview.id}
                                ref={versionPreviewRef}
                                defaultValue={cloneValue(versionPreview.value)}
                                height="100%"
                                readOnly
                                showControls
                                showMiniMap={false}
                                shortcuts={false}
                            />
                        </div>
                    </div>
                )}
            </Modal>

            <Modal
                title="新建导图"
                open={createOpen}
                onOk={createDocument}
                onCancel={() => setCreateOpen(false)}
                okText="创建"
                destroyOnHidden
            >
                <Form form={form} layout="vertical">
                    <Form.Item name="title" label="导图名称" rules={[{required: true, message: "请输入导图名称"}]}>
                        <Input maxLength={60} autoFocus/>
                    </Form.Item>
                    <Form.Item name="description" label="描述">
                        <Input.TextArea rows={3} maxLength={200}/>
                    </Form.Item>
                    <Form.Item name="categoryId" label="所属分类" rules={[{required: true, message: "请选择所属分类"}]}>
                        <Select options={categories.map((item) => ({label: item.name, value: item.id}))}/>
                    </Form.Item>
                    <Form.Item name="templateId" label="模板" rules={[{required: true, message: "请选择模板"}]}>
                        <Select
                            options={templateItems.map((item) => ({
                                label: `${item.name} - ${item.description}`,
                                value: item.id,
                            }))}
                        />
                    </Form.Item>
                    <div className={styles.templateGrid}>
                        {templateItems.map((item) => (
                            <div
                                className={`template ${selectedTemplateId === item.id ? "active" : ""}`}
                                key={item.id}
                                style={{"--accent": item.color} as React.CSSProperties}
                                onClick={() => form.setFieldValue("templateId", item.id)}
                            >
                                <div className="name">{item.name}</div>
                                <div className="desc">{item.description}</div>
                            </div>
                        ))}
                    </div>
                </Form>
            </Modal>

            <Modal
                title={categoryRecord ? "修改分类" : "新建分类"}
                open={categoryOpen}
                onOk={saveCategory}
                onCancel={() => setCategoryOpen(false)}
                okText="保存"
                destroyOnHidden
            >
                <Form form={categoryForm} layout="vertical">
	                    <Form.Item name="name" label="分类名称" rules={[{required: true, message: "请输入分类名称"}]}>
	                        <Input maxLength={24} autoFocus/>
	                    </Form.Item>
	                </Form>
            </Modal>
        </Body>
    );
};

export default MindPage;
