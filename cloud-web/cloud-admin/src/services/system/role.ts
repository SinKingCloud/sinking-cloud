// @ts-ignore
import {request} from 'umi';

/** 获取角色列表 GET /admin/systems/role/list */
export async function getRoleList(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/role/list', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 添加角色 POST /admin/systems/role/create */
export async function createRole(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/role/create', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 更新角色 POST /admin/systems/role/update */
export async function updateRole(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/role/update', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 删除角色 POST /admin/systems/role/delete */
export async function deleteRole(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/role/delete', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
