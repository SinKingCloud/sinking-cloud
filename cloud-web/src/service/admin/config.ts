import {API} from "../../../typings";
import {post} from "@/utils/request";

/** 网站配置 */
export async function webConfig(params: API.RequestParams = {}) {
    return post("/admin/config/web", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 界面配置 */
export async function uiConfig(params: API.RequestParams = {}) {
    return post("/admin/config/ui", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 阿里云配置 */
export async function aliyunConfig(params: API.RequestParams = {}) {
    return post("/admin/config/aliyun", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 支付配置 */
export async function payConfig(params: API.RequestParams = {}) {
    return post("/admin/pay/config", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}
