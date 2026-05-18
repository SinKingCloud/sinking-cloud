import {API} from "../../../typings";
import {get, post} from "@/utils/request";


/** 定时任务列表 */
export async function getCronList(params: API.RequestParams = {}) {
    return get("/admin/cron/list", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 定时任务日志 */
export async function getCronLog(params: API.RequestParams = {}) {
    return get("/admin/cron/log", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 定时任务详情 */
export async function getCronInfo(params: API.RequestParams = {}) {
    return get("/admin/cron/info", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 执行定时任务 */
export async function runCron(params: API.RequestParams = {}) {
    return post("/admin/cron/run", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 创建定时任务 */
export async function createCron(params: API.RequestParams = {}) {
    return post("/admin/cron/create", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 删除定时任务 */
export async function deleteCron(params: API.RequestParams = {}) {
    return post("/admin/cron/delete", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 修改定时任务 */
export async function updateCron(params: API.RequestParams = {}) {
    return post("/admin/cron/update", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}