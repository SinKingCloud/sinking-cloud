import {API} from "../../../typings";
import {get, post} from "@/utils/request";

/** 支付通道列表 */
export async function getPayChannelList(params: API.RequestParams = {}) {
    return get("/admin/pay/channel/list", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 支付通道详情 */
export async function getPayChannelInfo(params: API.RequestParams = {}) {
    return get("/admin/pay/channel/info", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 创建支付通道 */
export async function createPayChannel(params: API.RequestParams = {}) {
    return post("/admin/pay/channel/create", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 删除支付通道 */
export async function deletePayChannel(params: API.RequestParams = {}) {
    return post("/admin/pay/channel/delete", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 修改支付通道 */
export async function updatePayChannel(params: API.RequestParams = {}) {
    return post("/admin/pay/channel/update", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 用户余额详情 */
export async function getUserBalanceInfo(params: API.RequestParams = {}) {
    return get("/admin/pay/balance/info", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 用户余额操作 */
export async function operateUserBalance(params: API.RequestParams = {}) {
    return post("/admin/pay/balance/operate", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 用户余额明细 */
export async function getUserBalanceLogList(params: API.RequestParams = {}) {
    return get("/admin/pay/balance/log", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}
