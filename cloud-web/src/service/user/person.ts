import {API} from "../../../typings";
import {get, post} from "@/utils/request";

/** 获取个人信息 GET /user/person/info */
export async function getPersonInfo(params: API.RequestParams = {}) {
    return get("/user/person/info", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 修改个人信息 POST /user/person/update */
export async function updatePersonInfo(params: API.RequestParams = {}) {
    return post("/user/person/update", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 安全验证 POST /user/security/verify */
export async function verifySecurity(params: API.RequestParams = {}) {
    return post("/user/security/verify", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 修改密码 POST /user/security/password */
export async function updatePassword(params: API.RequestParams = {}) {
    return post("/user/security/password", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 修改邮箱 POST /user/security/email */
export async function updateEmail(params: API.RequestParams = {}) {
    return post("/user/security/email", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 修改手机 POST /user/security/phone */
export async function updatePhone(params: API.RequestParams = {}) {
    return post("/user/security/phone", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 修改账号 POST /user/security/account */
export async function updateAccount(params: API.RequestParams = {}) {
    return post("/user/security/account", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 会话列表 GET /user/security/session */
export async function getSessionList(params: API.RequestParams = {}) {
    return get("/user/security/session", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 注销会话 POST /user/security/logout */
export async function logoutSession(params: API.RequestParams = {}) {
    return post("/user/security/logout", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 实名认证 POST /user/security/certificate */
export async function certificate(params: API.RequestParams = {}) {
    return post("/user/security/certificate", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}
