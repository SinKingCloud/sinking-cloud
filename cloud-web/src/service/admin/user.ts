import {API} from "../../../typings";
import {get, post} from "@/utils/request";

/** 用户列表 */
export async function getUserList(params: API.RequestParams = {}) {
    return get("/admin/user/list", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 修改用户 */
export async function updateUser(params: API.RequestParams = {}) {
    return post("/admin/user/update", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 用户实名信息 */
export async function getUserCertificate(params: API.RequestParams = {}) {
    return get("/admin/user/certificate", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 用户会话列表 */
export async function getUserSessionList(params: API.RequestParams = {}) {
    return get("/admin/user/session", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 注销用户会话 */
export async function logoutUserSession(params: API.RequestParams = {}) {
    return post("/admin/user/logout", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}
