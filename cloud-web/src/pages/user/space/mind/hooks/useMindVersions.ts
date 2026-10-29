import {useCallback, useEffect, useMemo, useRef, useState} from "react";
import {App} from "antd";
import {getMindStats} from "@/pages/components/mind";
import {getMindVersionInfo, getMindVersionList} from "@/service/user/mind";
import {VERSION_MANUAL, VERSION_PAGE_SIZE} from "../constants";
import type {CursorState, MindDocument, MindVersionRecord, MindVersionSaveType, PageResult,} from "../types";
import {getResponseData, toMindVersion} from "../utils";

type UseMindVersionsOptions = {
    selectedDocumentRef: {current: MindDocument | undefined};
};

const initialCursor = (): CursorState => ({hasMore: true, nextId: "", nextLastId: ""});

export const useMindVersions = ({selectedDocumentRef}: UseMindVersionsOptions) => {
    const {message} = App.useApp();
    const loadingRef = useRef(false);
    const requestRef = useRef(0);
    const infoRequestRef = useRef(0);
    const openRef = useRef(false);
    const filterRef = useRef<MindVersionSaveType>(VERSION_MANUAL);
    const cursorRef = useRef<CursorState>(initialCursor());
    const [open, setOpen] = useState(false);
    const [filter, setFilter] = useState<MindVersionSaveType>(VERSION_MANUAL);
    const [previewOpen, setPreviewOpen] = useState(false);
    const [preview, setPreview] = useState<MindVersionRecord>();
    const [items, setItems] = useState<MindVersionRecord[]>([]);
    const [loading, setLoading] = useState(false);
    const [loadFailed, setLoadFailed] = useState(false);
    const currentRecord = useMemo(() => items.find((item) => item.current), [items]);

    useEffect(() => {
        openRef.current = open;
        filterRef.current = filter;
    }, [filter, open]);

    const loadMore = useCallback(async (dataId: string, type: MindVersionSaveType, reset = false) => {
        if (!reset && loadingRef.current) return;
        const cursor = reset ? initialCursor() : cursorRef.current;
        if (!cursor.hasMore) return;
        const requestId = ++requestRef.current;
        loadingRef.current = true;
        setLoading(true);
        setLoadFailed(false);
        try {
            const body: Record<string, any> = {
                data_id: dataId,
                type: String(type),
                page_size: VERSION_PAGE_SIZE,
                order_by_field: "id",
                order_by_type: "desc",
            };
            if (cursor.nextId) body.cursor_id = cursor.nextId;
            if (cursor.nextLastId) body.cursor_last_id = cursor.nextLastId;
            const response = await getMindVersionList({body});
            if (requestId !== requestRef.current) return;
            const result = getResponseData<PageResult<any>>(response, "获取版本记录失败");
            const nextItems = (Array.isArray(result?.list) ? result?.list : []).map(toMindVersion).filter((item) => item.id);
            setItems((current) => {
                const source = reset ? nextItems : [...current, ...nextItems];
                return Array.from(new Map(source.map((item) => [item.id, item])).values());
            });
            cursorRef.current = {
                hasMore: Boolean(result?.has_more),
                nextId: result?.next_cursor_id || "",
                nextLastId: result?.next_cursor_last_id || "",
            };
        } catch (error) {
            if (requestId === requestRef.current) {
                setLoadFailed(true);
                message.error(error instanceof Error ? error.message : "获取版本记录失败");
            }
        } finally {
            if (requestId === requestRef.current) {
                loadingRef.current = false;
                setLoading(false);
            }
        }
    }, [message]);

    const openHistory = useCallback(() => {
        const document = selectedDocumentRef.current;
        if (!document) return;
        const type = currentRecord?.saveType ?? VERSION_MANUAL;
        filterRef.current = type;
        openRef.current = true;
        setFilter(type);
        setItems([]);
        setLoadFailed(false);
        infoRequestRef.current += 1;
        setPreviewOpen(false);
        setPreview(undefined);
        setOpen(true);
        cursorRef.current = initialCursor();
        void loadMore(document?.id, type, true);
    }, [currentRecord?.saveType, loadMore, selectedDocumentRef]);

    const changeFilter = useCallback((type: MindVersionSaveType) => {
        const document = selectedDocumentRef.current;
        if (!document) return;
        filterRef.current = type;
        setFilter(type);
        setItems([]);
        setLoadFailed(false);
        cursorRef.current = initialCursor();
        void loadMore(document?.id, type, true);
    }, [loadMore, selectedDocumentRef]);

    const openPreview = useCallback(async (record: MindVersionRecord) => {
        const document = selectedDocumentRef.current;
        if (!document?.id || !record?.id) return;
        const requestId = ++infoRequestRef.current;
        setPreview(undefined);
        setPreviewOpen(true);
        try {
            const response = await getMindVersionInfo({body: {data_id: document?.id, id: record?.id}});
            if (requestId !== infoRequestRef.current
                || selectedDocumentRef.current?.id !== document?.id) return;
            const info = toMindVersion(getResponseData<any>(response, "获取版本详情失败"));
            if (!info?.value) throw new Error("版本内容为空");
            const stats = getMindStats(info?.value);
            setPreview({...info, nodeCount: stats.nodeCount, depth: stats.depth});
        } catch (error) {
            if (requestId === infoRequestRef.current) {
                setPreviewOpen(false);
                setPreview(undefined);
                message.error(error instanceof Error ? error.message : "获取版本详情失败");
            }
        }
    }, [message, selectedDocumentRef]);

    const closePreview = useCallback(() => {
        infoRequestRef.current += 1;
        setPreviewOpen(false);
        setPreview(undefined);
    }, []);

    const close = useCallback(() => {
        openRef.current = false;
        requestRef.current += 1;
        infoRequestRef.current += 1;
        loadingRef.current = false;
        setLoading(false);
        setOpen(false);
        setPreviewOpen(false);
        setPreview(undefined);
    }, []);

    const reset = useCallback(() => {
        close();
        setItems([]);
        setLoadFailed(false);
        cursorRef.current = initialCursor();
    }, [close]);

    const recordSavedVersion = useCallback((record: MindVersionRecord, type: MindVersionSaveType) => {
        if (!openRef.current || !record?.id) return;
        setItems((current) => {
            const previous = current
                .filter((item) => item.id !== record?.id)
                .map((item) => ({...item, current: false}));
            return filterRef.current === type ? [record, ...previous] : previous;
        });
    }, []);

    return {
        open,
        filter,
        previewOpen,
        previewLoading: previewOpen && !preview,
        preview,
        items,
        loading,
        loadFailed,
        loadMore,
        openHistory,
        changeFilter,
        openPreview,
        closePreview,
        close,
        reset,
        recordSavedVersion,
    };
};

export default useMindVersions;
