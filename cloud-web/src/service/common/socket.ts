import {API} from "../../../typings";
import {get} from "@/utils/request";

/** 获取前台用户 WebSocket 连接 token GET /auth/socket */
export async function getSocketToken(params: API.RequestParams = {}) {
    return get("/auth/socket", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}
