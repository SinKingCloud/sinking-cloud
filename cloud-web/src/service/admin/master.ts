import {API} from "../../../typings";
import {get, post} from "@/utils/request";

const request = (method: any, url: string, params: API.RequestParams = {}) => {
    return method(url, params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
};

/** 员工列表 */
export async function getMasterList(params: API.RequestParams = {}) {
    return request(get, "/admin/master/list", params);
}

/** 创建员工 */
export async function createMaster(params: API.RequestParams = {}) {
    return request(post, "/admin/master/create", params);
}

/** 修改员工 */
export async function updateMaster(params: API.RequestParams = {}) {
    return request(post, "/admin/master/update", params);
}

/** 删除员工 */
export async function deleteMaster(params: API.RequestParams = {}) {
    return request(post, "/admin/master/delete", params);
}

/** 分配员工角色 */
export async function updateMasterRole(params: API.RequestParams = {}) {
    return request(post, "/admin/master/role", params);
}

/** 员工会话列表 */
export async function getMasterSessionList(params: API.RequestParams = {}) {
    return request(get, "/admin/master/session", params);
}

/** 注销员工会话 */
export async function logoutMasterSession(params: API.RequestParams = {}) {
    return request(post, "/admin/master/logout", params);
}
