import {API} from "../../../typings";
import {get, post} from "@/utils/request";

/** 员工列表 */
export async function getMasterList(params: API.RequestParams = {}) {
    return get("/admin/master/list", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 创建员工 */
export async function createMaster(params: API.RequestParams = {}) {
    return post("/admin/master/create", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 修改员工 */
export async function updateMaster(params: API.RequestParams = {}) {
    return post("/admin/master/update", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 删除员工 */
export async function deleteMaster(params: API.RequestParams = {}) {
    return post("/admin/master/delete", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 分配员工角色 */
export async function updateMasterRole(params: API.RequestParams = {}) {
    return post("/admin/master/role", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 员工会话列表 */
export async function getMasterSessionList(params: API.RequestParams = {}) {
    return get("/admin/master/session", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 员工操作日志 */
export async function getMasterLogList(params: API.RequestParams = {}) {
    return get("/admin/master/log", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 注销员工会话 */
export async function logoutMasterSession(params: API.RequestParams = {}) {
    return post("/admin/master/logout", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}
