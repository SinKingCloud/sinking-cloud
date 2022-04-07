// @ts-ignore
import {request} from 'umi';

/** 获取用户列表 GET /admin/web/user/list */
export async function getUserList(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/web/user/list', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 更新用户 POST /admin/web/user/update */
export async function updateUser(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/web/user/update', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

