// @ts-ignore
import {request} from 'umi';

/** oauth列表 GET /open/user/oauth/list */
export async function getOauthList(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/oauth/web/list', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 添加oauth POST /open/user/oauth/create */
export async function createOauth(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/oauth/web/create', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 更新oauth POST /open/user/oauth/update */
export async function updateOauth(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/oauth/web/update', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 删除oauth POST /open/user/oauth/delete */
export async function deleteOauth(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/oauth/web/delete', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
