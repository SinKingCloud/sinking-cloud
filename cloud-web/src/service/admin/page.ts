import {API} from "../../../typings";
import {get} from "@/utils/request";

/** 获取菜单 */
export async function getMenus(params: API.RequestParams = {}) {
    return get("/admin/page/menu", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 获取控件 */
export async function getControl(params: API.RequestParams = {}) {
    return get("/admin/page/control", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}