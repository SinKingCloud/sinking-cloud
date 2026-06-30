import {API} from "../../../typings";
import {get, post} from "@/utils/request";

/** 用户消息列表 */
export async function getMessageList(params: API.RequestParams = {}) {
    return get("/user/message/list", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 用户消息详情 */
export async function getMessageInfo(params: API.RequestParams = {}) {
    return get("/user/message/info", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 用户消息已读 */
export async function readMessage(params: API.RequestParams = {}) {
    return post("/user/message/read", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 删除用户消息 */
export async function deleteMessage(params: API.RequestParams = {}) {
    return post("/user/message/delete", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 清空用户消息 */
export async function clearMessage(params: API.RequestParams = {}) {
    return post("/user/message/clear", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}
