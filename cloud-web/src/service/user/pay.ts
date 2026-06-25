import {API} from "../../../typings";
import {get, post} from "@/utils/request";

/** 支付配置 */
export async function getPayConfig(params: API.RequestParams = {}) {
    return get("/pay/config", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 获取支付链接 */
export async function submitPay(params: API.RequestParams = {}) {
    return post("/pay/submit", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 用户余额详情 */
export async function getBalanceInfo(params: API.RequestParams = {}) {
    return get("/user/pay/balance/info", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 用户余额明细 */
export async function getBalanceLogList(params: API.RequestParams = {}) {
    return get("/user/pay/balance/log", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 在线充值 */
export async function rechargeBalance(params: API.RequestParams = {}) {
    return post("/user/pay/balance/recharge", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 用户订单列表 */
export async function getOrderList(params: API.RequestParams = {}) {
    return get("/user/pay/order/list", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 用户订单日志 */
export async function getOrderLogList(params: API.RequestParams = {}) {
    return get("/user/pay/order/log", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}
