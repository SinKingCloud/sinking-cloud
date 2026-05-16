import {API} from "../../../typings";
import {get, post} from "@/utils/request";

/** 角色列表 */
export async function getRoleList(params: API.RequestParams = {}) {
    return get("/admin/role/list", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 角色详情 */
export async function getRoleInfo(params: API.RequestParams = {}) {
    return get("/admin/role/info", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 角色菜单树 */
export async function getRoleMenu(params: API.RequestParams = {}) {
    return get("/admin/role/menu", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 创建角色 */
export async function createRole(params: API.RequestParams = {}) {
    return post("/admin/role/create", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 修改角色 */
export async function updateRole(params: API.RequestParams = {}) {
    return post("/admin/role/update", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 删除角色 */
export async function deleteRole(params: API.RequestParams = {}) {
    return post("/admin/role/delete", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}
