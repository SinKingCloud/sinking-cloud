import {API} from "../../../typings";
import {get, post} from "@/utils/request";

/** 获取个人信息 */
export async function getPersonInfo(params: API.RequestParams = {}) {
    return get("/admin/person/info", params);
}

/** 更新个人资料 */
export async function updatePersonInfo(params: API.RequestParams = {}) {
    return post("/admin/person/update", params);
}

/** 安全验证信息/提交验证 */
export async function verifySecurity(params: API.RequestParams = {}) {
    return post("/admin/security/verify", params);
}

/** 修改账号 */
export async function updateAccount(params: API.RequestParams = {}) {
    return post("/admin/security/account", params);
}

/** 修改手机号 */
export async function updatePhone(params: API.RequestParams = {}) {
    return post("/admin/security/phone", params);
}

/** 修改邮箱 */
export async function updateEmail(params: API.RequestParams = {}) {
    return post("/admin/security/email", params);
}

/** 修改密码 */
export async function updatePassword(params: API.RequestParams = {}) {
    return post("/admin/security/password", params);
}

/** 会话列表 */
export async function getSessionList(params: API.RequestParams = {}) {
    return get("/admin/security/session", params);
}

/** 注销会话 */
export async function logoutSession(params: API.RequestParams = {}) {
    return post("/admin/security/logout", params);
}
