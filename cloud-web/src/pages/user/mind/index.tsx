import React, {useEffect, useMemo, useRef, useState} from "react";
import {App, Button, Card, Drawer, Dropdown, Form, Input, Modal, Select, Tooltip,} from "antd";
import {
    ApartmentOutlined,
    CloseOutlined,
    CopyOutlined,
    DeleteOutlined,
    DownloadOutlined,
    EditOutlined,
    FileImageOutlined,
    FileMarkdownOutlined,
    FileTextOutlined,
    HistoryOutlined,
    MoreOutlined,
    PlusOutlined,
    SaveOutlined,
    ShareAltOutlined,
    StarFilled,
    StarOutlined,
} from "@ant-design/icons";
import {createStyles} from "antd-style";
import {Body, ProTable, type ProTableRef, Title} from "sinking-antd";
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
    favorite?: boolean;
    shared?: boolean;
    createTime: string;
    updateTime: string;
};

type MindVersionRecord = {
    id: string;
    mindId: string;
    version: number;
    remark: string;
    nodeCount: number;
    depth: number;
    value: MindValue;
    createTime: string;
};

type CreateMindForm = {
    title: string;
    description?: string;
    categoryId: string;
    templateId: string;
};

type CategoryForm = {
    name: string;
};

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

const initialCategories: MindCategory[] = [
    {id: "product", name: "产品构想"},
    {id: "tech", name: "技术方案"},
    {id: "knowledge", name: "知识整理"},
    {id: "personal", name: "个人计划"},
    {id: "server", name: "服务器管理"},
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

const initialDocuments: MindDocument[] = [
    withStats({
        id: "mind-product-roadmap",
        categoryId: "product",
        title: "云上豁者产品路线图",
        description: "个人工作台、博客主页和工具插件的产品结构。",
        tags: ["工作台", "博客", "插件"],
        favorite: true,
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
        favorite: true,
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
];

const createVersionRecord = (record: MindDocument, remark = "保存快照"): MindVersionRecord => ({
    id: `${record.id}-v${record.version}-${Date.now().toString(36)}`,
    mindId: record.id,
    version: record.version,
    remark,
    nodeCount: record.nodeCount,
    depth: record.depth,
    value: cloneValue(record.value),
    createTime: record.updateTime,
});

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
        min-height: var(--mind-page-height);

        @supports not (height: 100dvh) {
            --mind-page-height: calc(100vh - 177px);
        }

        @media (max-width: 640px) {
            min-height: auto;
        }
    `,
    library: css`
        display: grid;
        grid-template-columns: 238px minmax(0, 1fr);
        gap: 18px;
        min-width: 0;

        @media (max-width: 900px) {
            grid-template-columns: 1fr;
            gap: 12px;
        }
    `,
    categoryPanel: css`
        position: sticky;
        top: 0;
        display: flex;
        flex-direction: column;
        min-width: 0;
        max-height: var(--mind-page-height);
        border: 0;
        border-radius: ${token.borderRadiusLG}px;
        background: ${token.colorBgContainer};
        align-self: start;
        overflow: hidden;

        > .ant-card-head {
            min-height: 48px;
            padding: 0 12px;
            border-bottom: 1px solid ${token.colorBorderSecondary};
        }

        > .ant-card-head .ant-card-head-title {
            padding: 11px 0;
        }

        > .ant-card-head .ant-card-extra {
            padding: 0;
        }

        > .ant-card-body {
            display: flex;
            flex-direction: column;
            flex: 1;
            min-height: 0;
            padding: 7px;
            overflow: hidden;
        }

        @media (max-width: 900px) {
            position: static;
            max-height: min(360px, var(--mind-page-height));
        }

        @media (max-width: 640px) {
            height: auto;
            max-height: none;

            > .ant-card-head {
                min-height: 44px;
                padding: 0 10px;
            }

            > .ant-card-head .ant-card-head-title {
                padding: 9px 0;
            }

            > .ant-card-body {
                padding: 7px 8px;
            }
        }
    `,
    categoryScroll: css`
        flex: 1;
        min-height: 0;
        overflow-y: auto;
        padding-right: 2px;
        scrollbar-width: thin;

        @media (max-width: 640px) {
            overflow-x: auto;
            overflow-y: hidden;
            padding-right: 0;
            scrollbar-width: none;

            &::-webkit-scrollbar {
                display: none;
            }
        }
    `,
    categoryFooter: css`
        flex: none;
        padding-top: 7px;

        @media (max-width: 640px) {
            display: none;
        }
    `,
    categoryGrid: css`
        display: flex;
        flex-direction: column;
        gap: 3px;

        @media (max-width: 640px) {
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
        min-height: 38px;
        padding: 0 5px 0 10px;
        overflow: hidden;
        border: 0;
        border-radius: ${token.borderRadius}px;
        background: transparent;
        cursor: pointer;
        transition: border-color 0.16s ease, background 0.16s ease, color 0.16s ease;

        &:hover,
        &.active {
            background: ${isDarkMode ? "rgba(255,255,255,0.04)" : token.colorFillQuaternary};
        }

        &.active {
            color: ${token.colorTextHeading};
            background: ${isDarkMode ? "color-mix(in srgb, var(--accent), transparent 86%)" : "color-mix(in srgb, var(--accent), #ffffff 92%)"};
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
            box-shadow: 0 0 0 3px color-mix(in srgb, var(--accent), transparent 86%);
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
        }

        .action .ant-btn {
            color: ${token.colorTextTertiary};
        }

        @media (max-width: 640px) {
            flex: 0 0 auto;
            min-height: 34px;
            padding: 0 8px 0 9px;

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
        background: transparent;
        color: ${token.colorTextTertiary};
        cursor: pointer;
        font-size: 13px;
        line-height: 20px;
        text-align: left;
        transition: border-color 0.16s ease, background 0.16s ease, color 0.16s ease;

        &.mobile {
            display: none;
        }

        &:hover {
            border-color: ${token.colorPrimaryBorder};
            background: ${isDarkMode ? "rgba(255,255,255,0.035)" : token.colorFillQuaternary};
            color: ${token.colorPrimary};
        }

        .icon {
            display: inline-flex;
            align-items: center;
            justify-content: center;
            width: 18px;
            height: 18px;
            border-radius: ${token.borderRadiusSM}px;
            background: ${isDarkMode ? "rgba(255,255,255,0.04)" : token.colorFillTertiary};
        }

        @media (max-width: 640px) {
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
    content: css`
        min-width: 0;
    `,
    documentTable: css`
        .ant-table {
            background: transparent;
            white-space: normal !important;
        }

        .ant-table-container,
        .ant-table-content,
        .ant-table-tbody {
            background: transparent;
        }

        .ant-table-tbody > tr > td.ant-table-cell {
            padding: 0 !important;
            border-bottom: 0 !important;
            background: transparent !important;
        }

        .ant-table-tbody > tr:hover > td.ant-table-cell {
            background: transparent !important;
        }
    `,
    tableHeader: css`
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 16px;
        width: 100%;
        min-width: 0;

        > *:first-child {
            flex: none;
        }

        @media (max-width: 640px) {
            align-items: stretch;
            flex-direction: column;
            gap: 10px;
        }
    `,
    searchBox: css`
        width: 320px;

        &.ant-input-search > .ant-input,
        &.ant-input-search > .ant-input-affix-wrapper {
            border-top-right-radius: 0 !important;
            border-bottom-right-radius: 0 !important;
        }

        &.ant-input-search > .ant-btn,
        &.ant-input-search > .ant-input-search-button,
        &.ant-input-search > .ant-btn-compact-last-item {
            border-top-left-radius: 0 !important;
            border-bottom-left-radius: 0 !important;
        }

        @media (max-width: 820px) {
            width: 260px;
        }

        @media (max-width: 640px) {
            width: 100%;
        }
    `,
    documentGrid: css`
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(208px, 1fr));
        gap: 12px;

        @media (max-width: 640px) {
            grid-template-columns: 1fr;
            gap: 10px;
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
        transition: border-color 0.16s ease, transform 0.16s ease, box-shadow 0.16s ease;

        &.ant-card > .ant-card-body {
            padding: 0 !important;
        }

        &:hover {
            border-color: ${token.colorPrimary};
            box-shadow: ${isDarkMode ? "none" : "0 10px 24px rgba(15, 23, 42, 0.06)"};
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

        .preview-root {
            filter: ${isDarkMode ? "none" : "drop-shadow(0 4px 8px rgba(15, 23, 42, 0.06))"};
        }

        .preview-dot {
            fill: color-mix(in srgb, var(--accent), ${token.colorBgContainer} 18%);
        }

        .body {
            padding: 9px 10px 10px;
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
            margin-right: -4px;
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

            > .ant-card-body {
                display: grid;
                grid-template-columns: 108px minmax(0, 1fr);
            }

            .preview {
                height: auto;
                min-height: 104px;
                border-right: 1px solid ${token.colorBorderSecondary};
                border-bottom: 0;
            }

            .body {
                min-width: 0;
                padding: 9px 10px;
            }

            .head {
                gap: 6px;
                min-height: 22px;
            }

            .title {
                font-size: 13px;
                line-height: 20px;
            }

            .desc {
                -webkit-line-clamp: 2;
                min-height: auto;
                margin-top: 2px;
            }

            .meta {
                flex-wrap: wrap;
                justify-content: flex-start;
                gap: 3px 10px;
                margin-top: 6px;
            }
        }

        @media (max-width: 420px) {
            > .ant-card-body {
                grid-template-columns: 92px minmax(0, 1fr);
            }

            .preview {
                min-height: 96px;
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
        box-shadow: ${isDarkMode ? "none" : "0 14px 34px rgba(18,30,46,0.055)"};
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
    versionList: css`
        display: flex;
        flex-direction: column;
        gap: 8px;

        .version {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 12px;
            padding: 10px 12px;
            border: 1px solid ${token.colorBorderSecondary};
            border-radius: ${token.borderRadius}px;
            background: ${token.colorBgContainer};
        }

        .main {
            min-width: 0;
        }

        .name {
            color: ${token.colorTextHeading};
            font-size: 13px;
            line-height: 20px;
            font-weight: 600;
        }

        .desc {
            margin-top: 2px;
            color: ${token.colorTextTertiary};
            font-size: 12px;
            line-height: 18px;
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
    const pageRef = useRef<HTMLDivElement>(null);
    const editorRef = useRef<MindRef>(null);
    const tableRef = useRef<ProTableRef>(null);
    const [form] = Form.useForm<CreateMindForm>();
    const [categoryForm] = Form.useForm<CategoryForm>();
    const selectedTemplateId = Form.useWatch("templateId", form);
    const [categories, setCategories] = useState<MindCategory[]>(initialCategories);
    const [documents, setDocuments] = useState<MindDocument[]>(initialDocuments);
    const [activeCategoryId, setActiveCategoryId] = useState(ALL_CATEGORY_ID);
    const [keyword, setKeyword] = useState("");
    const [selectedId, setSelectedId] = useState<string>();
    const [editorValue, setEditorValue] = useState<MindValue>(initialDocuments[0].value);
    const [editorRevision, setEditorRevision] = useState(0);
    const [createOpen, setCreateOpen] = useState(false);
    const [categoryOpen, setCategoryOpen] = useState(false);
    const [categoryRecord, setCategoryRecord] = useState<MindCategory>();
    const [versionOpen, setVersionOpen] = useState(false);
    const [versions, setVersions] = useState<Record<string, MindVersionRecord[]>>(initialVersions);
    const [pageHeight, setPageHeight] = useState<number>();

    const selectedDocument = documents.find((item) => item.id === selectedId);
    const versionItems = selectedId ? versions[selectedId] || [] : [];

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

    const filteredDocuments = useMemo(() => {
        const search = keyword.trim().toLowerCase();
        return documents
            .filter((item) => activeCategoryId === ALL_CATEGORY_ID || item.categoryId === activeCategoryId)
            .filter((item) => !search || `${item.title} ${item.description} ${item.tags.join(" ")}`.toLowerCase().includes(search))
            .sort((a, b) => new Date(b.updateTime.replace(/-/g, "/")).getTime() - new Date(a.updateTime.replace(/-/g, "/")).getTime());
    }, [activeCategoryId, documents, keyword]);

    const currentCategory = categories.find((item) => item.id === activeCategoryId);
    const categoryName = activeCategoryId === ALL_CATEGORY_ID ? "全部导图" : currentCategory?.name || "导图";
    const tableKey = `${activeCategoryId}-${keyword.trim()}`;

    useEffect(() => {
        tableRef.current?.refreshTableData?.();
    }, [documents]);

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
        setCreateOpen(false);
        openDocument(record);
        message.success("导图已创建");
    };

    const persistDocument = async (value?: MindValue, stats?: MindStats, reason: MindSaveReason = "manual") => {
        if (!selectedId) {
            return;
        }
        const nextValue = value || editorRef.current?.getValue() || editorValue;
        const nextStats = stats || getMindStats(nextValue);
        const now = formatDateTime(new Date());
        setDocuments((items) => items.map((item) => item.id === selectedId ? {
            ...item,
            value: cloneValue(nextValue),
            nodeCount: nextStats.nodeCount,
            depth: nextStats.depth,
            version: reason === "auto" ? item.version : item.version + 1,
            updateTime: now,
        } : item));
        setEditorValue(cloneValue(nextValue));
        if (reason !== "auto") {
            const nextVersion = (selectedDocument?.version || 0) + 1;
            setVersions((items) => ({
                ...items,
                [selectedId]: [{
                    id: `${selectedId}-v${nextVersion}-${Date.now().toString(36)}`,
                    mindId: selectedId,
                    version: nextVersion,
                    remark: reason === "shortcut" ? "快捷键保存" : "手动保存",
                    nodeCount: nextStats.nodeCount,
                    depth: nextStats.depth,
                    value: cloneValue(nextValue),
                    createTime: now,
                }, ...(items[selectedId] || [])].slice(0, 20),
            }));
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
            favorite: false,
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
                if (record.id === selectedId) {
                    setSelectedId(undefined);
                    setVersionOpen(false);
                }
                message.success("导图已删除");
            },
        });
    };

    const toggleFavorite = (record: MindDocument) => {
        setDocuments((items) => items.map((item) => item.id === record.id ? {...item, favorite: !item.favorite} : item));
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
        const nextValue = cloneValue(record.value);
        const stats = getMindStats(nextValue);
        const nextVersion = selectedDocument.version + 1;
        setDocuments((items) => items.map((item) => item.id === selectedDocument.id ? {
            ...item,
            value: nextValue,
            nodeCount: stats.nodeCount,
            depth: stats.depth,
            version: nextVersion,
            updateTime: now,
        } : item));
        setVersions((items) => ({
            ...items,
            [selectedDocument.id]: [{
                id: `${selectedDocument.id}-v${nextVersion}-${Date.now().toString(36)}`,
                mindId: selectedDocument.id,
                version: nextVersion,
                remark: `从 v${record.version} 恢复`,
                nodeCount: stats.nodeCount,
                depth: stats.depth,
                value: nextValue,
                createTime: now,
            }, ...(items[selectedDocument.id] || [])].slice(0, 20),
        }));
        setEditorValue(cloneValue(nextValue));
        setEditorRevision((value) => value + 1);
        message.success("版本已恢复");
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

    const categoryMap = useMemo(() => new Map(categories.map((item) => [item.id, item])), [categories]);

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
    const closeEditor = () => {
        setSelectedId(undefined);
        setVersionOpen(false);
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
                        <button className={styles.editorToolButton} type="button" onClick={() => setVersionOpen(true)}>
                            <HistoryOutlined/>
                        </button>
                    </Tooltip>
                    <Tooltip title={selectedDocument.shared ? "关闭分享" : "开启分享"}>
                        <button className={styles.editorToolButton} type="button" onClick={() => toggleShare(selectedDocument)}>
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
                    <Card
                        className={styles.categoryPanel}
                        title={<Title>导图分类</Title>}
                    >
                        <div className={styles.categoryScroll}>
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
                        <div className={styles.categoryFooter}>
                            <button className={styles.categoryAddCard} type="button" onClick={openCreateCategory}>
                                <span className="icon"><PlusOutlined/></span>
                                <span>添加分类</span>
                            </button>
                        </div>
                    </Card>

                    <div className={styles.content}>
                        <ProTable
                            key={tableKey}
                            ref={tableRef}
                            rowKey="id"
                            title={(
                                <div className={styles.tableHeader}>
                                    <Title>{categoryName}</Title>
                                    <Input.Search
                                        allowClear
                                        className={styles.searchBox}
                                        placeholder="搜索标题、描述或标签"
                                        enterButton="搜索"
                                        value={keyword}
                                        onChange={(event) => setKeyword(event.target.value)}
                                        onSearch={(value) => setKeyword(value)}
                                    />
                                </div>
                            )}
                            search={false}
                            pageHidden={filteredDocuments.length <= 0}
                            defaultPage={1}
                            defaultPageSize={12}
                            tableProps={{
                                showHeader: false,
                                className: styles.documentTable,
                                style: {whiteSpace: "normal"},
                            } as any}
                            request={async (params) => {
                                const current = Number(params?.current || 1);
                                const pageSize = Number(params?.pageSize || 12);
                                const start = (current - 1) * pageSize;
                                const list = filteredDocuments.slice(start, start + pageSize);
                                return {
                                    success: true,
                                    data: list.length <= 0 ? [] : [{
                                        id: "mind-document-card-list",
                                        list,
                                    }],
                                    total: filteredDocuments.length,
                                };
                            }}
                            columns={[{
                                title: "导图",
                                dataIndex: "list",
                                render: (_: any, record: {list: MindDocument[]}) => (
                                    <div className={styles.documentGrid}>
                                        {record.list.map((item) => {
                                            const category = categoryMap.get(item.categoryId);
                                            return (
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
                                                                <Tooltip title={item.favorite ? "取消收藏" : "收藏"}>
                                                                    <Button
                                                                        type="text"
                                                                        size="small"
                                                                        icon={item.favorite ? <StarFilled/> : <StarOutlined/>}
                                                                        onClick={() => toggleFavorite(item)}
                                                                    />
                                                                </Tooltip>
                                                                <Dropdown trigger={["click"]} placement="bottomRight" menu={{items: documentMenu(item) as any}}>
                                                                    <Button type="text" size="small" icon={<MoreOutlined/>}/>
                                                                </Dropdown>
                                                            </span>
                                                        </div>
                                                        <div className="desc">{item.description || "暂无描述"}</div>
                                                        <div className="meta">
                                                            <span>{category?.name || "未分类"}</span>
                                                            <span>{item.nodeCount} 节点</span>
                                                            <Tooltip title={item.updateTime}>
                                                                <span>{getTimeText(item.updateTime)}</span>
                                                            </Tooltip>
                                                        </div>
                                                    </div>
                                                </Card>
                                            );
                                        })}
                                    </div>
                                ),
                            } as any]}
                        />
                    </div>
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
                open={versionOpen}
                onClose={() => setVersionOpen(false)}
                destroyOnHidden
            >
                {versionItems.length <= 0 ? (
                    <div className={styles.versionList}>
                        <div className="version">
                            <div className="main">
                                <div className="name">暂无版本记录</div>
                                <div className="desc">手动保存后会生成版本快照。</div>
                            </div>
                        </div>
                    </div>
                ) : (
                    <div className={styles.versionList}>
                        {versionItems.map((item) => (
                            <div className="version" key={item.id}>
                                <div className="main">
                                    <div className="name">版本 v{item.version}</div>
                                    <div className="desc">{item.remark} · {item.nodeCount} 节点 · {item.depth} 层 · {item.createTime}</div>
                                </div>
                                <Button size="small" onClick={() => restoreVersion(item)}>恢复</Button>
                            </div>
                        ))}
                    </div>
                )}
            </Drawer>

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
