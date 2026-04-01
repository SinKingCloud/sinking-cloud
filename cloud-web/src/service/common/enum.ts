import {API} from "../../../typings";
import {get} from "@/utils/request";

/** 获取枚举数据 */
export async function getEnum(params: API.RequestParams = {}) {
    return get("/enum", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}
