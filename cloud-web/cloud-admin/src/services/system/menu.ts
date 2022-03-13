// @ts-ignore
import {request} from 'umi';

/** 获取接口列表 GET /admin/systems/menu/list */
export async function getMenuList(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/menu/list', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 添加权限 POST /admin/systems/menu/create */
export async function createMenu(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/menu/create', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 更新权限 POST /admin/systems/menu/update */
export async function updateMenu(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/menu/update', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 删除权限 POST /admin/systems/menu/delete */
export async function deleteMenu(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/menu/delete', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
