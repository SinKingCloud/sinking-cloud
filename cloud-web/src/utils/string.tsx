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
 * 获取字符串get参数
 * @param url
 */
export function getUrlQuery(url: string) {
    const noMatchData = {};
    if (!url) return noMatchData;
    const reg = /(?:\?|&)(?:(?:([^=&]+?)=([^&#]+))|([^&#]+))/gi;
    const obj = {};
    // @ts-ignore
    const fun = (match, $1, $2 = '', $3) => {
        if ($3) {
            obj[$3] = '';
        } else {
            obj[$1] = encodeURIComponent($2);
        }
    }
    // @ts-ignore
    url.replace(reg, fun);
    return obj
}

/**
 * 获取访问字符
 * @param url
 */
export function getQueryString(url: string = ""): string {
    if (url == "") {
        url = window.location.href;
    }
    const firstIndex = decodeURI(url).indexOf("?");
    if (firstIndex < 0) {
        return ""
    }
    return url.substring(firstIndex + 1);
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
