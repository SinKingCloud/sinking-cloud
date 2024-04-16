import request from "@/utils/request";
import {API} from "@/../typings";

/*****  获取网站用户信息  *****/
export async function getWebUserInfo(params?: {}, options?: { [key: string]: any }) {
    return request<API.Response>('/user/person/info', {
        method: 'GET',
        params: {
            ...params,
        },
        ...(options || {}),
    });
}

/*****  获取系统用户信息  *****/
export async function getSystemUserInfo(params?: {}, options?: { [key: string]: any }) {
    return request<API.Response>('/admin/person/info', {
        method: 'GET',
        params: {
            ...params,
        },
        ...(options || {}),
    });
}