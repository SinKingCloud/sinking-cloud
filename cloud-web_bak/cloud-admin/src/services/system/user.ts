// @ts-ignore
import {request} from 'umi';

/** 获取用户列表 GET /admin/systems/user/list */
export async function getUserList(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/user/list', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 添加用户 POST /admin/systems/user/create */
export async function createUser(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/user/create', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 更新用户 POST /admin/systems/user/update */
export async function updateUser(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/user/update', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 删除用户 POST /admin/systems/user/delete */
export async function deleteUser(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/user/delete', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
