/**
 * 将枚举对象转换为 Select options。
 */
export const buildOptions = (data: any, numberValue = false) => Object.entries(data || {}).map(([key, value]) => ({
    label: value,
    value: numberValue ? Number(key) : key,
}));

/**
 * 将枚举对象转换为 ProTable 的 valueEnum。
 */
export const buildValueEnum = (data: any, colorMap?: any) => Object.fromEntries(Object.entries(data || {}).map(([key, value]) => {
    const item: any = {text: value};
    if (colorMap?.[key]) {
        item.color = colorMap[key];
    }
    return [key, item];
}));

/**
 * 生成日期范围搜索字段转换函数。
 */
export const dateRangeTransform = (prefix: string) => (value: any) => ({
    [`${prefix}_start`]: value?.[0]?.format ? value[0].format("YYYY-MM-DD HH:mm:ss") : value?.[0],
    [`${prefix}_end`]: value?.[1]?.format ? value[1].format("YYYY-MM-DD HH:mm:ss") : value?.[1],
});

/**
 * 规范化通道配置，兼容字符串和对象两种格式。
 */
export const normalizeConfig = (config: any) => {
    if (typeof config === "string") {
        try {
            return JSON.parse(config || "{}");
        } catch (_) {
            return {};
        }
    }
    return config || {};
};

/**
 * 格式化文件大小。
 */
export const formatSize = (size: any) => {
    const value = Number(size || 0);
    if (value <= 0) {
        return "0 B";
    }
    const units = ["B", "KB", "MB", "GB", "TB"];
    let nextValue = value;
    let index = 0;
    while (nextValue >= 1024 && index < units.length - 1) {
        nextValue = nextValue / 1024;
        index += 1;
    }
    return `${nextValue.toFixed(index === 0 ? 0 : 2)} ${units[index]}`;
};

/**
 * 格式化 KB 单位的大小。
 */
export const formatSizeFromKb = (size: any) => {
    const value = Number(size || 0);
    if (value <= 0) {
        return "0 KB";
    }
    return `${value} KB`;
};
