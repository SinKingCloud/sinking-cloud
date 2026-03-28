import {API} from "../../../typings";
import {get} from "@/utils/request";

/** 获取个人信息 GET /user/person/info */
export async function getPersonInfo(params: API.RequestParams = {}) {
    return get("/admin/person/info", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}