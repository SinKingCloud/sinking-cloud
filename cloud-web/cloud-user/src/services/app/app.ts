// @ts-ignore
import {request} from 'umi';

/** 获取应用列表 GET /user/app/list */
export async function getAppList(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/user/app/list', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 添加应用 POST /user/app/create */
export async function createApp(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/user/app/create', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 更新应用 POST /user/app/update */
export async function updateApp(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/user/app/update', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 删除应用 POST /user/app/delete */
export async function deleteApp(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/user/app/delete', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
