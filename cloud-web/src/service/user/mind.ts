import {API} from "../../../typings";
import {get, post} from "@/utils/request";

/** 获取导图列表 */
export async function getMindDataList(params: API.RequestParams = {}) {
    return get("/user/space/mind/data/list", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 获取导图详情 */
export async function getMindDataInfo(params: API.RequestParams = {}) {
    return get("/user/space/mind/data/info", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 创建导图 */
export async function createMindData(params: API.RequestParams = {}) {
    return post("/user/space/mind/data/create", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 更新导图信息 */
export async function updateMindData(params: API.RequestParams = {}) {
    return post("/user/space/mind/data/update", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 删除导图 */
export async function deleteMindData(params: API.RequestParams = {}) {
    return post("/user/space/mind/data/delete", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 复制导图 */
export async function copyMindData(params: API.RequestParams = {}) {
    return post("/user/space/mind/data/copy", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 保存导图内容 */
export async function saveMindData(params: API.RequestParams = {}) {
    return post("/user/space/mind/data/save", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 获取自定义分类 */
export async function getMindTypeList(params: API.RequestParams = {}) {
    return get("/user/space/mind/type/list", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 创建自定义分类 */
export async function createMindType(params: API.RequestParams = {}) {
    return post("/user/space/mind/type/create", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 更新自定义分类 */
export async function updateMindType(params: API.RequestParams = {}) {
    return post("/user/space/mind/type/update", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 删除自定义分类 */
export async function deleteMindType(params: API.RequestParams = {}) {
    return post("/user/space/mind/type/delete", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 获取版本列表 */
export async function getMindVersionList(params: API.RequestParams = {}) {
    return get("/user/space/mind/version/list", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 获取版本详情 */
export async function getMindVersionInfo(params: API.RequestParams = {}) {
    return get("/user/space/mind/version/info", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 恢复版本 */
export async function restoreMindVersion(params: API.RequestParams = {}) {
    return post("/user/space/mind/version/restore", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 获取公开导图分享 */
export async function getMindShare(params: API.RequestParams = {}) {
    return get("/space/mind/share", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 获取分享列表 */
export async function getMindShareList(params: API.RequestParams = {}) {
    return get("/user/space/mind/share/list", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 创建分享 */
export async function createMindShare(params: API.RequestParams = {}) {
    return post("/user/space/mind/share/create", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 更新分享 */
export async function updateMindShare(params: API.RequestParams = {}) {
    return post("/user/space/mind/share/update", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 删除分享 */
export async function deleteMindShare(params: API.RequestParams = {}) {
    return post("/user/space/mind/share/delete", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}
