/**
 * 获取随机字符
 * @param length 长度
 */
export function getRandStr(length = 16) {
    let str = 'abcdefghijklmnopqrstuvwxyz';
    str += str.toUpperCase();
    str += '0123456789'
    let _str = '';
    for (let i = 0; i < length; i++) {
        const rand = Math.floor(Math.random() * str.length);
        _str += str[rand];
    }
    return _str
}

/**
 * 规范化配置，兼容字符串和对象两种格式。
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

/**
 * 格式化金额。
 */
export const formatMoney = (value: any) => {
    const money = Number(value || 0);
    return Number.isFinite(money) ? money.toFixed(2) : "0.00";
};

/**
 * 复制文本。
 */
export const copyText = async (text?: string) => {
    if (!text) {
        return false;
    }
    try {
        if (navigator?.clipboard?.writeText) {
            await navigator.clipboard.writeText(text);
            return true;
        }
        const textarea = document.createElement("textarea");
        textarea.value = text;
        textarea.style.position = "fixed";
        textarea.style.opacity = "0";
        document.body.appendChild(textarea);
        textarea.select();
        const successful = document.execCommand("copy");
        document.body.removeChild(textarea);
        return successful;
    } catch {
        return false;
    }
};
