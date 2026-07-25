import type {MindSortType, MindVersionSaveType} from "./types";

export const mindSortOptions: Array<{label: string; value: MindSortType}> = [
    {label: "最近更新", value: "update-desc"},
    {label: "最近创建", value: "create-desc"},
];

export const ALL_CATEGORY_ID = "0";
export const CATEGORY_PAGE_SIZE = 20;
export const VERSION_PAGE_SIZE = 30;
export const SHARE_PAGE_SIZE = 20;
export const VERSION_MANUAL: MindVersionSaveType = 0;
export const VERSION_AUTO: MindVersionSaveType = 1;
export const SHARE_ENABLE = 0;
export const SHARE_DISABLE = 1;
