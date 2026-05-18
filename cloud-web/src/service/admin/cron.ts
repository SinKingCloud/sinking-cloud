import {API} from "../../../typings";
import {get, post} from "@/utils/request";


/** 定时任务列表 */
export async function getTaskList(params: API.RequestParams = {}) {
    return get("/admin/task/list", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 定时任务日志 */
export async function getTaskLog(params: API.RequestParams = {}) {
    return get("/admin/task/log", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 定时任务详情 */
export async function getTaskInfo(params: API.RequestParams = {}) {
    return get("/admin/task/info", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 执行定时任务 */
export async function runTask(params: API.RequestParams = {}) {
    return post("/admin/task/run", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 创建定时任务 */
export async function createTask(params: API.RequestParams = {}) {
    return post("/admin/task/create", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 删除定时任务 */
export async function deleteTask(params: API.RequestParams = {}) {
    return post("/admin/task/delete", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 修改定时任务 */
export async function updateTask(params: API.RequestParams = {}) {
    return post("/admin/task/update", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}