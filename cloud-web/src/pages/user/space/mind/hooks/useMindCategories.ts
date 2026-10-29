import type {UIEvent} from "react";
import {useCallback, useEffect, useMemo, useRef, useState} from "react";
import {App, Form} from "antd";
import {createMindType, deleteMindType, getMindTypeList, updateMindType,} from "@/service/user/mind";
import {getEnum} from "@/utils/enum";
import {ALL_CATEGORY_ID, CATEGORY_PAGE_SIZE} from "../constants";
import type {CategoryForm, CursorState, MindCategory, PageResult} from "../types";
import {getResponseData, toMindCategory} from "../utils";

export type UseMindCategoriesOptions = {
    onDeleted: (categoryId: string) => void;
};

export type MindCategoryOption = {
    label: string;
    value: string;
};

const initialCursor = (): CursorState => ({hasMore: true, nextId: "", nextLastId: ""});

const initialPresetCategories = (): MindCategory[] => [
    {id: ALL_CATEGORY_ID, name: "全部类型", preset: true},
];

export const useMindCategories = ({onDeleted}: UseMindCategoriesOptions) => {
    const {message, modal} = App.useApp();
    const loadingRef = useRef(false);
    const presetInitializedRef = useRef(false);
    const customInitializedRef = useRef(false);
    const cursorRef = useRef<CursorState>(initialCursor());
    const [form] = Form.useForm<CategoryForm>();
    const [presetCategories, setPresetCategories] = useState<MindCategory[]>(initialPresetCategories);
    const [customCategories, setCustomCategories] = useState<MindCategory[]>([]);
    const [loading, setLoading] = useState(true);
    const [loadFailed, setLoadFailed] = useState(false);
    const [activeCategoryId, setActiveCategoryId] = useState(ALL_CATEGORY_ID);
    const [formOpen, setFormOpen] = useState(false);
    const [editingCategory, setEditingCategory] = useState<MindCategory>();
    const [submitting, setSubmitting] = useState(false);
    const categories = useMemo(
        () => [...presetCategories, ...customCategories],
        [customCategories, presetCategories],
    );
    const baseCategoryOptions = useMemo<MindCategoryOption[]>(
        () => categories.map((item) => ({label: item.name, value: item.id})),
        [categories],
    );

    useEffect(() => {
        let active = true;
        getEnum("space.mind").then((data: any) => {
            if (!active) return;
            const types = data?.type || {};
            const items = Object.entries(types)
                .map(([id, name]) => ({
                    id,
                    name: String(name),
                    preset: true,
                }))
                .sort((left, right) => Number(left.id) - Number(right.id));
            setPresetCategories(items.length > 0 ? items : initialPresetCategories());
        }).catch((error) => {
            if (active) {
                message.error(error instanceof Error ? error.message : "获取预设分类失败");
            }
        }).finally(() => {
            if (active) {
                presetInitializedRef.current = true;
                if (customInitializedRef.current) setLoading(false);
            }
        });
        return () => {
            active = false;
        };
    }, [message]);

    const loadMoreCategories = useCallback(async (reset = false) => {
        if (loadingRef.current) return;
        const cursor = reset ? initialCursor() : cursorRef.current;
        if (!cursor.hasMore) return;
        if (reset) {
            customInitializedRef.current = false;
            setLoading(true);
        }
        loadingRef.current = true;
        setLoadFailed(false);
        try {
            const body: Record<string, any> = {
                page_size: CATEGORY_PAGE_SIZE,
                order_by_field: "id",
                order_by_type: "desc",
            };
            if (cursor.nextId) body.cursor_id = cursor.nextId;
            if (cursor.nextLastId) body.cursor_last_id = cursor.nextLastId;
            const response = await getMindTypeList({body});
            const result = getResponseData<PageResult<any>>(response, "获取分类失败");
            const nextItems = (Array.isArray(result?.list) ? result?.list : []).map(toMindCategory).filter((item) => item.id);
            setCustomCategories((current) => {
                const source = reset ? nextItems : [...current, ...nextItems];
                return Array.from(new Map(source.map((item) => [item.id, item])).values());
            });
            cursorRef.current = {
                hasMore: Boolean(result?.has_more),
                nextId: result?.next_cursor_id || "",
                nextLastId: result?.next_cursor_last_id || "",
            };
        } catch (error) {
            cursorRef.current = {...cursorRef.current, hasMore: false};
            setLoadFailed(true);
            message.error(error instanceof Error ? error.message : "获取分类失败");
        } finally {
            loadingRef.current = false;
            if (reset) {
                customInitializedRef.current = true;
                if (presetInitializedRef.current) setLoading(false);
            }
        }
    }, [message]);

    useEffect(() => {
        void loadMoreCategories(true);
    }, [loadMoreCategories]);

    const retryCategories = useCallback(() => {
        cursorRef.current = {...cursorRef.current, hasMore: true};
        void loadMoreCategories(customCategories.length <= 0);
    }, [customCategories.length, loadMoreCategories]);

    const handleCategoryPopupScroll = useCallback((event: UIEvent<HTMLDivElement>) => {
        const element = event.currentTarget;
        if (element.scrollHeight - element.scrollTop - element.clientHeight <= 72) {
            void loadMoreCategories();
        }
    }, [loadMoreCategories]);

    const getCategorySelectOptions = useCallback((editingCategoryId?: string): MindCategoryOption[] => {
        const options = [...baseCategoryOptions];
        if (editingCategoryId !== undefined && !options.some((item) => item.value === editingCategoryId)) {
            options.push({label: `当前分类（${editingCategoryId}）`, value: editingCategoryId});
        }
        return options;
    }, [baseCategoryOptions]);

    const openCreateCategory = () => {
        setEditingCategory(undefined);
        form.setFieldsValue({name: ""});
        setFormOpen(true);
    };

    const openEditCategory = (category: MindCategory) => {
        if (!category?.id || category?.preset) return;
        setEditingCategory(category);
        form.setFieldsValue({name: category?.name});
        setFormOpen(true);
    };

    const closeCategoryForm = () => {
        setFormOpen(false);
        setEditingCategory(undefined);
    };

    const saveCategory = async () => {
        try {
            const values = await form.validateFields();
            setSubmitting(true);
            if (editingCategory) {
                const response = await updateMindType({body: {
                    id: Number(editingCategory?.id),
                    name: values.name,
                }});
                getResponseData(response, "更新分类失败");
                setCustomCategories((items) => items.map((item) => item.id === editingCategory?.id ? {
                    ...item,
                    name: values.name.trim(),
                } : item));
                setFormOpen(false);
                message.success("分类已更新");
                return;
            }
            const response = await createMindType({body: {name: values.name}});
            const category = toMindCategory(getResponseData<any>(response, "创建分类失败"));
            if (!category?.id) throw new Error("创建分类失败");
            setCustomCategories((items) => [category, ...items.filter((item) => item.id !== category?.id)]);
            setActiveCategoryId(category?.id);
            setFormOpen(false);
            message.success("分类已创建");
        } catch (error) {
            if ((error as any)?.errorFields) return;
            message.error(error instanceof Error ? error.message : "保存分类失败");
        } finally {
            setSubmitting(false);
        }
    };

    const deleteCategory = (category: MindCategory) => {
        if (!category?.id || category?.preset) return;
        modal.confirm({
            title: "删除分类",
            content: `删除${category?.name}后，该分类下的导图会归入全部导图。`,
            okText: "删除",
            okButtonProps: {danger: true},
            onOk: async () => {
                try {
                    const response = await deleteMindType({body: {id: Number(category?.id)}});
                    getResponseData(response, "删除分类失败");
                    setCustomCategories((items) => items.filter((item) => item.id !== category?.id));
                    if (activeCategoryId === category?.id) {
                        setActiveCategoryId(ALL_CATEGORY_ID);
                    }
                    onDeleted(category?.id);
                    message.success("分类已删除");
                } catch (error) {
                    message.error(error instanceof Error ? error.message : "删除分类失败");
                }
            },
        });
    };

    return {
        categoryForm: form,
        categories,
        categoryLoading: loading,
        categoryLoadFailed: loadFailed,
        activeCategoryId,
        setActiveCategoryId,
        categoryOpen: formOpen,
        categoryRecord: editingCategory,
        categorySubmitting: submitting,
        loadMoreCategories,
        retryCategories,
        handleCategoryPopupScroll,
        getCategorySelectOptions,
        openCreateCategory,
        openEditCategory,
        closeCategoryForm,
        saveCategory,
        deleteCategory,
    };
};

export default useMindCategories;
