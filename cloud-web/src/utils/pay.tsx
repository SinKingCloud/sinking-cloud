/** 设置支付跳转地址 */
export function setPayJumpUrl(url: string = "") {
    localStorage.setItem("payJumpUrl", url || window.location.href);
}

/** 获取支付跳转地址 */
export function getPayJumpUrl(): string {
    return localStorage.getItem("payJumpUrl") || "";
}
