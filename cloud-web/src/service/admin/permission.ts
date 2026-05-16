import {API} from "../../../typings";
import {get, post} from "@/utils/request";

/** 权限列表 */
export async function getPermissionList(params: API.RequestParams = {}) {
    return get("/admin/permission/list", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 创建权限 */
export async function createPermission(params: API.RequestParams = {}) {
    return post("/admin/permission/create", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 修改权限 */
export async function updatePermission(params: API.RequestParams = {}) {
    return post("/admin/permission/update", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 删除权限 */
export async function deletePermission(params: API.RequestParams = {}) {
    return post("/admin/permission/delete", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}
