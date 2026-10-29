import type {API} from "@/../typings";
import type {MindValue} from "@/pages/components/mind";
import {ago} from "@/utils/time";
import {VERSION_AUTO, VERSION_MANUAL} from "./constants";
import type {MindCategory, MindDocument, MindShareRecord, MindVersionRecord} from "./types";

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

export const getStableColor = (value?: string | number) => {
    const text = String(value || "mind");
    let hash = 0;
    for (let index = 0; index < text.length; index += 1) {
        hash = (hash * 31 + text.charCodeAt(index)) >>> 0;
    }
    return previewColors[hash % previewColors.length];
};

export const getTimeText = (value: string) => ago(value) || value;

export const cloneValue = (value: MindValue): MindValue => JSON.parse(JSON.stringify(value));

export const getMindValueSignature = (value: MindValue) => JSON.stringify(value);

export const createBlankValue = (title: string): MindValue => ({
    nodes: [{id: "root", title, root: true, color: "#1677ff"}],
});

export const parseMindValue = (value: unknown): MindValue => {
    const data = typeof value === "string" ? JSON.parse(value) : value;
    if (!data || typeof data !== "object" || !Array.isArray((data as MindValue)?.nodes)) {
        throw new Error("导图内容格式错误");
    }
    return data as MindValue;
};

export const getResponseData = <T, >(response: API.Response, fallback: string): T => {
    if (response?.code !== 200) {
        throw new Error(response?.message || fallback);
    }
    return response?.data as T;
};

const toId = (value: unknown) => String(value ?? "");

export const toMindDocument = (data: any): MindDocument => ({
    id: toId(data?.id),
    categoryId: toId(data?.type_id ?? 0),
    title: data?.title || "未命名导图",
    description: data?.describe || "",
    value: data?.value ? parseMindValue(data?.value) : undefined,
    shared: Boolean(data?.shared),
    createTime: data?.create_time || "",
    updateTime: data?.update_time || "",
});

export const toMindCategory = (data: any): MindCategory => ({
    id: toId(data?.id),
    name: data?.name || "未命名分类",
    preset: false,
    createTime: data?.create_time || "",
    updateTime: data?.update_time || "",
});

export const toMindVersion = (data: any): MindVersionRecord => ({
    id: toId(data?.id),
    mindId: toId(data?.data_id),
    saveType: Number(data?.type) === VERSION_AUTO ? VERSION_AUTO : VERSION_MANUAL,
    remark: data?.remark || (Number(data?.type) === VERSION_AUTO ? "自动保存" : "手动保存"),
    current: Boolean(data?.current),
    value: data?.value ? parseMindValue(data?.value) : undefined,
    createTime: data?.create_time || "",
    updateTime: data?.update_time || "",
});

export const toMindShare = (data: any): MindShareRecord => ({
    id: toId(data?.id),
    mindId: toId(data?.data_id),
    code: data?.code || "",
    status: Number(data?.status),
    expireTime: data?.expire_time && !String(data?.expire_time).startsWith("0001-") ? data?.expire_time : "",
    remark: data?.remark || "",
    createTime: data?.create_time || "",
    updateTime: data?.update_time || "",
});

export const downloadDataUrl = (dataUrl: string, filename: string) => {
    const link = document.createElement("a");
    link.href = dataUrl;
    link.download = filename;
    link.click();
};
