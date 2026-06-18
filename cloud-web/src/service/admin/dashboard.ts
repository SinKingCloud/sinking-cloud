import {API} from "../../../typings";
import {get} from "@/utils/request";

/** 获取服务数量 */
export async function getServiceCount(params: API.RequestParams = {}) {
    return get("/admin/dashboard/service", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}
