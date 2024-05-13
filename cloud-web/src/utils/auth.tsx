import defaultSettings from "../../config/defaultSettings";

/**
 * 设置登录token
 * @param device 设备
 * @param token 登录token
 */
export function setLoginToken(device: string, token: string) {
    localStorage.setItem('sinking-token', token);
    localStorage.setItem('sinking-device', device);
}

/**
 * 获取登录token
 * @returns {string}
 */
export function getLoginToken() {
    return localStorage.getItem('sinking-token') || "";
}

/**
 * 获取登录方式
 * @returns {string}
 */
export function getLoginType() {
    return localStorage.getItem('sinking-device') || "";
}

/**
 * 删除token
 */
export function deleteHeader() {
    localStorage.removeItem('sinking-token');
    localStorage.removeItem('sinking-device');
}

/**
 * 返回请求header
 * @returns {{}}
 */
export function getHeaders() {
    return {
        'sinking-token': getLoginToken(),
        'sinking-device': getLoginType(),
        'sinking-version': defaultSettings?.version || 'none',
    };
}
