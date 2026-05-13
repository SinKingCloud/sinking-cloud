import {API} from "../../../typings";
import {get, post} from "@/utils/request";

/** 菜单列表 */
export async function getMenuList(params: API.RequestParams = {}) {
    return get("/admin/menu/list", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 菜单或控件详情 */
export async function getMenuInfo(params: API.RequestParams = {}) {
    return get("/admin/menu/info", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 查看页面控件列表 */
export async function getMenuControls(params: API.RequestParams = {}) {
    return get("/admin/menu/controls", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 创建页面或控件 */
export async function createMenu(params: API.RequestParams = {}) {
    return post("/admin/menu/create", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 删除页面或控件 */
export async function deleteMenu(params: API.RequestParams = {}) {
    return post("/admin/menu/delete", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 修改页面或控件 */
export async function updateMenu(params: API.RequestParams = {}) {
    return post("/admin/menu/update", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}
