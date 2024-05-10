import {get} from "@/utils/request";
import {API} from "@/../typings";

/*****  获取网站用户信息  *****/
export async function getWebInfo(params: API.RequestParams = {}) {
    return get("/auth/web/info", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}
