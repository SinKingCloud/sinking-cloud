import {API} from "../../../typings";
import {get, post} from "@/utils/request";

/** 存储通道列表 */
export async function getChannelList(params: API.RequestParams = {}) {
    return get("/admin/driver/channel/list", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 存储通道详情 */
export async function getChannelInfo(params: API.RequestParams = {}) {
    return get("/admin/driver/channel/info", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 创建存储通道 */
export async function createChannel(params: API.RequestParams = {}) {
    return post("/admin/driver/channel/create", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 修改存储通道 */
export async function updateChannel(params: API.RequestParams = {}) {
    return post("/admin/driver/channel/update", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 删除存储通道 */
export async function deleteChannel(params: API.RequestParams = {}) {
    return post("/admin/driver/channel/delete", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 存储场景列表 */
export async function getSceneList(params: API.RequestParams = {}) {
    return get("/admin/driver/scene/list", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 创建存储场景 */
export async function createScene(params: API.RequestParams = {}) {
    return post("/admin/driver/scene/create", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 修改存储场景 */
export async function updateScene(params: API.RequestParams = {}) {
    return post("/admin/driver/scene/update", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 删除存储场景 */
export async function deleteScene(params: API.RequestParams = {}) {
    return post("/admin/driver/scene/delete", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 文件列表 */
export async function getFileList(params: API.RequestParams = {}) {
    return get("/admin/driver/file/list", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 修改文件 */
export async function updateFile(params: API.RequestParams = {}) {
    return post("/admin/driver/file/update", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 删除文件 */
export async function deleteFile(params: API.RequestParams = {}) {
    return post("/admin/driver/file/delete", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}
