import {API} from "@/../typings";
import {post} from "@/utils/request";

/** 登录接口 POST /auth/web/login */
export async function webLogin(params: API.RequestParams = {}) {
    return post("/auth/web/login", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}