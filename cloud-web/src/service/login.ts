import {API} from "@/../typings";
import {get, post} from "@/utils/request";

/** 登录接口 POST /auth/web/login */
export async function webLogin(params: API.RequestParams = {}) {
    return post("/auth/web/login", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/**
 * 注销登录
 */
export async function outLogin(params: API.RequestParams = {}) {
    return get("/auth/web/out_login", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/**
 * 获取短信验证码
 */
export async function getSms(params:API.RequestParams = {}) {
    return post("/auth/verify/sms", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}