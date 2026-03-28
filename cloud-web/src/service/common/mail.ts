import {post} from "@/utils/request";
import {API} from "../../../typings";

/** 发送邮箱验证码 POST /auth/mail */
export async function sendMail(params: API.RequestParams = {}) {
    return post("/auth/mail", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}