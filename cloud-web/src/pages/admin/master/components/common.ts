export const AVATAR_SCENE_ID = 1;

export const emptyTableData = {
    data: [],
    success: true,
    total: 0,
};

export const dateRangeTransform = (name: string) => (value: any) => ({
    [`${name}_start`]: value?.[0]?.format ? value[0].format("YYYY-MM-DD HH:mm:ss") : value?.[0],
    [`${name}_end`]: value?.[1]?.format ? value[1].format("YYYY-MM-DD HH:mm:ss") : value?.[1],
});

export const normalizeIds = (ids: any[] = []) => Array.from(new Set(
    (ids || []).map((item) => Number(item)).filter((item) => item > 0)
));

export const buildOptions = (data: any) => Object.entries(data || {}).map(([key, value]) => ({
    label: value,
    value: Number(key),
}));

export const buildValueEnum = (data: any, colors: Record<string, string> = {}) => Object.fromEntries(
    Object.entries(data || {}).map(([key, value]) => {
        const item: any = {text: value};
        if (colors[key]) {
            item.color = colors[key];
        }
        return [key, item];
    })
);

export const getUploadUrl = (result: any) => {
    if (result?.url) {
        return result.url;
    }
    if (typeof result?.data === "string") {
        return result.data;
    }
    return result?.data?.url || result?.data?.path || result?.data?.src;
};
