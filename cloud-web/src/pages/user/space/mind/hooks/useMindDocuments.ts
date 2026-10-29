import {useCallback, useEffect, useRef, useState} from "react";
import {App, Form} from "antd";
import type {MindValue} from "@/pages/components/mind";
import {
    copyMindData,
    createMindData,
    deleteMindData,
    getMindDataInfo,
    getMindDataList,
    updateMindData,
} from "@/service/user/mind";
import {ALL_CATEGORY_ID} from "../constants";
import type {MindDocument, MindFormValues, MindSortType, PageResult} from "../types";
import {createBlankValue, getResponseData, toMindDocument} from "../utils";

export type UseMindDocumentsOptions = {
    activeCategoryId: string;
    onOpen: (record: MindDocument, value: MindValue) => void;
    onUpdated: (record: MindDocument) => void;
    onDeleted: (record: MindDocument) => void;
};

export const useMindDocuments = ({
    activeCategoryId,
    onOpen,
    onUpdated,
    onDeleted,
}: UseMindDocumentsOptions) => {
    const {message, modal} = App.useApp();
    const requestRef = useRef(0);
    const infoRequestRef = useRef(0);
    const [form] = Form.useForm<MindFormValues>();
    const [documents, setDocuments] = useState<MindDocument[]>([]);
    const [documentTotal, setDocumentTotal] = useState(0);
    const [documentsLoading, setDocumentsLoading] = useState(true);
    const [reloadKey, setReloadKey] = useState(0);
    const [keywordInput, setKeywordInput] = useState("");
    const [keyword, setKeyword] = useState("");
    const [sortType, setSortType] = useState<MindSortType>("update-desc");
    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const [documentOpeningId, setDocumentOpeningId] = useState<string>();
    const [createOpen, setCreateOpen] = useState(false);
    const [editingDocument, setEditingDocument] = useState<MindDocument>();
    const [formSubmitting, setFormSubmitting] = useState(false);

    useEffect(() => {
        const timer = window.setTimeout(() => {
            setKeyword(keywordInput.trim());
            setPage(1);
        }, 350);
        return () => window.clearTimeout(timer);
    }, [keywordInput]);

    useEffect(() => {
        setPage(1);
    }, [activeCategoryId, sortType]);

    useEffect(() => {
        const requestId = ++requestRef.current;
        setDocumentsLoading(true);
        const body: Record<string, any> = {
            page,
            page_size: pageSize,
            order_by_field: sortType === "create-desc" ? "create_time" : "update_time",
            order_by_type: "desc",
        };
        if (activeCategoryId !== ALL_CATEGORY_ID) body.type_id = activeCategoryId;
        if (keyword) body.keyword = keyword;
        getMindDataList({body}).then((response) => {
            if (requestId !== requestRef.current) return;
            const result = getResponseData<PageResult<any>>(response, "获取导图失败");
            setDocuments((Array.isArray(result?.list) ? result?.list : []).map(toMindDocument).filter((item) => item.id));
            setDocumentTotal(Number(result?.total || 0));
        }).catch((error) => {
            if (requestId === requestRef.current) {
                setDocuments([]);
                setDocumentTotal(0);
                message.error(error instanceof Error ? error.message : "获取导图失败");
            }
        }).finally(() => {
            if (requestId === requestRef.current) {
                setDocumentsLoading(false);
            }
        });
    }, [activeCategoryId, keyword, message, page, pageSize, reloadKey, sortType]);

    const reloadDocuments = useCallback(() => {
        setReloadKey((value) => value + 1);
    }, []);

    const moveCategoryToAll = useCallback((categoryId: string) => {
        setDocuments((items) => items.map((item) => item.categoryId === categoryId ? {
            ...item,
            categoryId: ALL_CATEGORY_ID,
        } : item));
    }, []);

    const updateDocumentTime = useCallback((documentId: string, updateTime: string) => {
        setDocuments((items) => items.map((item) => item.id === documentId ? {
            ...item,
            updateTime,
        } : item));
    }, []);

    const searchDocuments = useCallback((value: string) => {
        setKeyword(value);
        setPage(1);
    }, []);

    const changePagination = useCallback((nextPage: number, nextPageSize: number) => {
        setPageSize(nextPageSize);
        setPage(nextPage);
    }, []);

    const openDocument = useCallback(async (record: MindDocument) => {
        if (!record?.id) return;
        const requestId = ++infoRequestRef.current;
        setDocumentOpeningId(record?.id);
        try {
            const response = await getMindDataInfo({body: {id: record?.id}});
            if (requestId !== infoRequestRef.current) return;
            const info = toMindDocument(getResponseData<any>(response, "获取导图失败"));
            if (!info?.id || !info?.value) throw new Error("导图内容为空");
            onOpen(info, info?.value);
        } catch (error) {
            if (requestId === infoRequestRef.current) {
                message.error(error instanceof Error ? error.message : "获取导图失败");
            }
        } finally {
            if (requestId === infoRequestRef.current) {
                setDocumentOpeningId(undefined);
            }
        }
    }, [message, onOpen]);

    const cancelOpenDocument = useCallback(() => {
        infoRequestRef.current += 1;
        setDocumentOpeningId(undefined);
    }, []);

    const openCreate = useCallback(() => {
        setEditingDocument(undefined);
        form.setFieldsValue({
            title: "",
            description: "",
            categoryId: activeCategoryId,
        });
        setCreateOpen(true);
    }, [activeCategoryId, form]);

    const openEditDocument = useCallback((record: MindDocument) => {
        if (!record?.id) return;
        setEditingDocument(record);
        form.setFieldsValue({
            title: record?.title,
            description: record?.description,
            categoryId: record?.categoryId,
        });
        setCreateOpen(true);
    }, [form]);

    const closeDocumentForm = useCallback(() => {
        setCreateOpen(false);
        setEditingDocument(undefined);
    }, []);

    const saveDocumentInfo = async () => {
        try {
            const values = await form.validateFields();
            setFormSubmitting(true);
            if (editingDocument) {
                const response = await updateMindData({body: {
                    id: Number(editingDocument?.id),
                    type_id: values.categoryId,
                    title: values.title,
                    describe: values.description || "",
                }});
                getResponseData(response, "更新导图失败");
                const nextRecord: MindDocument = {
                    ...editingDocument,
                    categoryId: values.categoryId,
                    title: values.title.trim(),
                    description: values.description?.trim() || "",
                };
                setDocuments((items) => items.map((item) => item.id === nextRecord.id ? nextRecord : item));
                onUpdated(nextRecord);
                setCreateOpen(false);
                message.success("导图信息已更新");
                reloadDocuments();
                return;
            }
            const value = createBlankValue(values.title.trim());
            const response = await createMindData({body: {
                type_id: Number(values.categoryId),
                title: values.title,
                describe: values.description || "",
                value: JSON.stringify(value),
            }});
            const record = {...toMindDocument(getResponseData<any>(response, "创建导图失败")), value};
            if (!record?.id) throw new Error("创建导图失败");
            setCreateOpen(false);
            onOpen(record, value);
            message.success("导图已创建");
            reloadDocuments();
        } catch (error) {
            if ((error as any)?.errorFields) return;
            message.error(error instanceof Error ? error.message : "保存导图信息失败");
        } finally {
            setFormSubmitting(false);
        }
    };

    const duplicateDocument = (record: MindDocument) => {
        if (!record?.id) return;
        modal.confirm({
            title: "复制导图",
            content: `确定复制${record?.title}吗？`,
            okText: "复制",
            onOk: async () => {
                try {
                    const response = await copyMindData({body: {id: Number(record?.id)}});
                    getResponseData(response, "复制导图失败");
                    message.success("导图已复制");
                    setPage(1);
                    reloadDocuments();
                } catch (error) {
                    message.error(error instanceof Error ? error.message : "复制导图失败");
                }
            },
        });
    };

    const deleteDocument = (record: MindDocument) => {
        if (!record?.id) return;
        modal.confirm({
            title: "删除导图",
            content: `确定删除${record?.title}吗？`,
            okText: "删除",
            okButtonProps: {danger: true},
            onOk: async () => {
                try {
                    const response = await deleteMindData({body: {id: Number(record?.id)}});
                    getResponseData(response, "删除导图失败");
                    setDocuments((items) => items.filter((item) => item.id !== record?.id));
                    onDeleted(record);
                    if (documents.length <= 1 && page > 1) {
                        setPage((current) => Math.max(1, current - 1));
                    } else {
                        reloadDocuments();
                    }
                    message.success("导图已删除");
                } catch (error) {
                    message.error(error instanceof Error ? error.message : "删除导图失败");
                }
            },
        });
    };

    return {
        form,
        documents,
        documentTotal,
        documentsLoading,
        keywordInput,
        sortType,
        page,
        pageSize,
        documentOpeningId,
        createOpen,
        editingDocument,
        formSubmitting,
        setKeywordInput,
        searchDocuments,
        setSortType,
        changePagination,
        openDocument,
        cancelOpenDocument,
        openCreate,
        openEditDocument,
        closeDocumentForm,
        saveDocumentInfo,
        duplicateDocument,
        deleteDocument,
        reloadDocuments,
        moveCategoryToAll,
        updateDocumentTime,
    };
};

export default useMindDocuments;
