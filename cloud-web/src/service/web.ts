import request from "@/utils/request";
import {API} from "@/../typings";

/*****  获取网站信息  *****/
export async function getWebInfo(params?: {}, options?: { [key: string]: any }) {
    return request<API.Response>('/auth/web/info', {
        method: 'GET',
        params: {
            ...params,
        },
        ...(options || {}),
    });
}
