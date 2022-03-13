// @ts-ignore
import {request} from 'umi';

/** 获取接口列表 GET /admin/systems/access/list */
export async function getAccessList(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/access/list', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 添加权限 POST /admin/systems/access/create */
export async function createAccess(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/access/create', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 更新权限 POST /admin/systems/access/update */
export async function updateAccess(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/access/update', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 删除权限 POST /admin/systems/access/delete */
export async function deleteAccess(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/access/delete', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
