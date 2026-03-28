import {API} from "@/../typings";
import {get} from "@/utils/request";

/*****  获取当前网站信息 GET /info  *****/
export async function getInfo(params: API.RequestParams = {}) {
    return get("/info", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}