import {useCallback, useEffect, useRef, useState} from "react";
import {App, Form} from "antd";
import {createMindShare, deleteMindShare, getMindShareList, updateMindShare,} from "@/service/user/mind";
import {copyText} from "@/utils/string";
import {useModel} from "umi";
import {SHARE_ENABLE, SHARE_PAGE_SIZE} from "../constants";
import type {CursorState, MindDocument, MindShareRecord, PageResult, ShareFormValues,} from "../types";
import {getResponseData, toMindShare} from "../utils";

export type UseMindSharesOptions = {
    onChanged: () => void;
};

const initialCursor = (): CursorState => ({hasMore: true, nextId: "", nextLastId: ""});

export const useMindShares = ({onChanged}: UseMindSharesOptions) => {
    const {message, modal} = App.useApp();
    const web = useModel("web");
    const loadingRef = useRef(false);
    const requestRef = useRef(0);
    const actionRequestRef = useRef(0);
    const documentRef = useRef<MindDocument | undefined>(undefined);
    const cursorRef = useRef<CursorState>(initialCursor());
    const [form] = Form.useForm<ShareFormValues>();
    const [open, setOpen] = useState(false);
    const [document, setDocument] = useState<MindDocument>();
    const [items, setItems] = useState<MindShareRecord[]>([]);
    const [loading, setLoading] = useState(false);
    const [loadFailed, setLoadFailed] = useState(false);
    const [formOpen, setFormOpen] = useState(false);
    const [editingShare, setEditingShare] = useState<MindShareRecord>();
    const [submitting, setSubmitting] = useState(false);

    useEffect(() => {
        documentRef.current = document;
    }, [document]);

    const loadShares = useCallback(async (dataId: string, reset = false) => {
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
                page_size: SHARE_PAGE_SIZE,
                order_by_field: "id",
                order_by_type: "desc",
            };
            if (cursor.nextId) body.cursor_id = cursor.nextId;
            if (cursor.nextLastId) body.cursor_last_id = cursor.nextLastId;
            const response = await getMindShareList({body});
            if (requestId !== requestRef.current) return;
            const result = getResponseData<PageResult<any>>(response, "获取分享列表失败");
            const nextItems = (Array.isArray(result?.list) ? result?.list : []).map(toMindShare).filter((item) => item.id);
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
                message.error(error instanceof Error ? error.message : "获取分享列表失败");
            }
        } finally {
            if (requestId === requestRef.current) {
                loadingRef.current = false;
                setLoading(false);
            }
        }
    }, [message]);

    const openShareManager = (record: MindDocument) => {
        if (!record?.id) return;
        actionRequestRef.current += 1;
        documentRef.current = record;
        setSubmitting(false);
        setDocument(record);
        setItems([]);
        setLoadFailed(false);
        setOpen(true);
        cursorRef.current = initialCursor();
        void loadShares(record?.id, true);
    };

    const closeShareManager = () => {
        requestRef.current += 1;
        actionRequestRef.current += 1;
        documentRef.current = undefined;
        loadingRef.current = false;
        setLoading(false);
        setSubmitting(false);
        setFormOpen(false);
        setOpen(false);
        setDocument(undefined);
        setItems([]);
    };

    const openCreateShare = () => {
        setEditingShare(undefined);
        form.setFieldsValue({expireType: "permanent", expireTime: undefined, remark: ""});
        setFormOpen(true);
    };

    const openEditShare = (record: MindShareRecord) => {
        if (!record?.id) return;
        setEditingShare(record);
        form.setFieldsValue({
            status: record?.status,
            expireType: record?.expireTime ? "custom" : "permanent",
            expireTime: record?.expireTime || undefined,
            remark: record?.remark,
        });
        setFormOpen(true);
    };

    const copyShareLink = useCallback(async (record: MindShareRecord) => {
        if (!record?.code) return;
        const siteUrl = String(web?.info?.url || "").trim();
        if (!siteUrl) {
            message.error("网站地址未配置");
            return;
        }
        try {
            const baseUrl = new URL(siteUrl);
            baseUrl.search = "";
            baseUrl.hash = "";
            baseUrl.pathname = `${baseUrl.pathname.replace(/\/+$/, "")}/`;
            const url = new URL(`mind/${encodeURIComponent(record?.code)}`, baseUrl);
            if (await copyText(url.toString())) {
                message.success("分享链接已复制");
                return;
            }
        } catch {
            message.error("网站地址配置错误");
            return;
        }
        message.error("复制链接失败");
    }, [message, web?.info?.url]);

    const closeShareForm = () => {
        setFormOpen(false);
        setEditingShare(undefined);
    };

    const saveShare = async () => {
        const currentDocument = documentRef.current;
        if (!currentDocument) return;
        const share = editingShare;
        const requestId = ++actionRequestRef.current;
        try {
            const values = await form.validateFields();
            setSubmitting(true);
            let expireTime: string | undefined = "";
            let expireDate: Date | undefined;
            if (values.expireType === "day") {
                expireDate = new Date();
                expireDate.setDate(expireDate.getDate() + 1);
            } else if (values.expireType === "threeDays") {
                expireDate = new Date();
                expireDate.setDate(expireDate.getDate() + 3);
            } else if (values.expireType === "week") {
                expireDate = new Date();
                expireDate.setDate(expireDate.getDate() + 7);
            } else if (values.expireType === "month") {
                expireDate = new Date();
                const currentDate = expireDate.getDate();
                expireDate.setDate(1);
                expireDate.setMonth(expireDate.getMonth() + 1);
                const lastDate = new Date(
                    expireDate.getFullYear(),
                    expireDate.getMonth() + 1,
                    0,
                ).getDate();
                expireDate.setDate(Math.min(currentDate, lastDate));
            } else if (values.expireType === "custom") {
                expireTime = values.expireTime || "";
                if (share?.expireTime === expireTime) expireTime = undefined;
            }
            if (expireDate) {
                const month = String(expireDate.getMonth() + 1).padStart(2, "0");
                const date = String(expireDate.getDate()).padStart(2, "0");
                const hour = String(expireDate.getHours()).padStart(2, "0");
                const minute = String(expireDate.getMinutes()).padStart(2, "0");
                const second = String(expireDate.getSeconds()).padStart(2, "0");
                expireTime = `${expireDate.getFullYear()}-${month}-${date} ${hour}:${minute}:${second}`;
            }
            if (share) {
                const body: Record<string, any> = {
                    data_id: Number(currentDocument?.id),
                    id: Number(share?.id),
                    status: String(values.status),
                    remark: values.remark || "",
                };
                if (expireTime !== undefined) body.expire_time = expireTime;
                const response = await updateMindShare({body});
                getResponseData(response, "更新分享失败");
            } else {
                const response = await createMindShare({body: {
                    data_id: Number(currentDocument?.id),
                    status: SHARE_ENABLE,
                    expire_time: expireTime || "",
                    remark: values.remark || "",
                }});
                getResponseData(response, "创建分享失败");
            }
            if (requestId !== actionRequestRef.current
                || documentRef.current?.id !== currentDocument?.id) return;
            message.success(share ? "分享已更新" : "分享已创建");
            setFormOpen(false);
            await loadShares(currentDocument?.id, true);
            onChanged();
        } catch (error) {
            if ((error as any)?.errorFields) return;
            if (requestId === actionRequestRef.current) {
                message.error(error instanceof Error ? error.message : "保存分享失败");
            }
        } finally {
            if (requestId === actionRequestRef.current) {
                setSubmitting(false);
            }
        }
    };

    const removeShare = (record: MindShareRecord) => {
        const currentDocument = documentRef.current;
        if (!currentDocument?.id || !record?.id) return;
        modal.confirm({
            title: "删除分享",
            content: `确定删除分享码 ${record?.code} 吗？`,
            okText: "删除",
            okButtonProps: {danger: true},
            onOk: async () => {
                try {
                    const response = await deleteMindShare({body: {
                        data_id: Number(currentDocument?.id),
                        id: Number(record?.id),
                    }});
                    getResponseData(response, "删除分享失败");
                    if (documentRef.current?.id === currentDocument?.id) {
                        await loadShares(currentDocument?.id, true);
                    }
                    onChanged();
                    message.success("分享已删除");
                } catch (error) {
                    message.error(error instanceof Error ? error.message : "删除分享失败");
                }
            },
        });
    };

    return {
        shareForm: form,
        shareOpen: open,
        shareDocument: document,
        shareItems: items,
        shareLoading: loading,
        shareLoadFailed: loadFailed,
        shareFormOpen: formOpen,
        editingShare,
        shareSubmitting: submitting,
        loadShares,
        copyShareLink,
        openShareManager,
        closeShareManager,
        openCreateShare,
        openEditShare,
        closeShareForm,
        saveShare,
        removeShare,
    };
};

export default useMindShares;
