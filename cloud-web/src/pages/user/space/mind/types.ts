import type {MindSaveReason, MindValue} from "@/pages/components/mind";

export type MindCategory = {
    id: string;
    name: string;
    preset: boolean;
    createTime?: string;
    updateTime?: string;
};

export type MindDocument = {
    id: string;
    categoryId: string;
    title: string;
    description: string;
    value?: MindValue;
    shared: boolean;
    createTime: string;
    updateTime: string;
};

export type MindVersionSaveType = 0 | 1;

export type MindVersionRecord = {
    id: string;
    mindId: string;
    saveType: MindVersionSaveType;
    remark: string;
    current: boolean;
    value?: MindValue;
    nodeCount?: number;
    depth?: number;
    createTime: string;
    updateTime: string;
};

export type MindFormValues = {
    title: string;
    description?: string;
    categoryId: string;
};

export type CategoryForm = {
    name: string;
};

export type MindShareRecord = {
    id: string;
    mindId: string;
    code: string;
    status: number;
    expireTime?: string;
    remark: string;
    createTime: string;
    updateTime: string;
};

export type ShareFormValues = {
    status: number;
    expireType: "permanent" | "day" | "threeDays" | "week" | "month" | "custom";
    expireTime?: string;
    remark?: string;
};

export type PageResult<T> = {
    total?: number;
    page?: number;
    page_size?: number;
    list?: T[];
    has_more?: boolean;
    next_cursor_id?: string;
    next_cursor_last_id?: string;
};

export type CursorState = {
    hasMore: boolean;
    nextId: string;
    nextLastId: string;
};

export type MindSaveTask = {
    dataId: string;
    value: MindValue;
    signature: string;
    saveType: MindVersionSaveType;
    reason: MindSaveReason;
    waiters: Array<{
        resolve: () => void;
        reject: (reason?: unknown) => void;
    }>;
};

export type MindSortType = "update-desc" | "create-desc";
