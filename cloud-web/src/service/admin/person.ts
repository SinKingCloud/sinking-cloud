import {API} from "../../../typings";
import {get, post} from "@/utils/request";

/** 获取个人信息 */
export async function getPersonInfo(params: API.RequestParams = {}) {
    return get("/admin/person/info", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 更新个人资料 */
export async function updatePersonInfo(params: API.RequestParams = {}) {
    return post("/admin/person/update", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 安全验证信息/提交验证 */
export async function verifySecurity(params: API.RequestParams = {}) {
    return post("/admin/security/verify", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 修改账号 */
export async function updateAccount(params: API.RequestParams = {}) {
    return post("/admin/security/account", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 修改手机号 */
export async function updatePhone(params: API.RequestParams = {}) {
    return post("/admin/security/phone", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 修改邮箱 */
export async function updateEmail(params: API.RequestParams = {}) {
    return post("/admin/security/email", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 修改密码 */
export async function updatePassword(params: API.RequestParams = {}) {
    return post("/admin/security/password", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 会话列表 */
export async function getSessionList(params: API.RequestParams = {}) {
    return get("/admin/security/session", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 操作日志 */
export async function getLogList(params: API.RequestParams = {}) {
    return get("/admin/person/log", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 注销会话 */
export async function logoutSession(params: API.RequestParams = {}) {
    return post("/admin/security/logout", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}
