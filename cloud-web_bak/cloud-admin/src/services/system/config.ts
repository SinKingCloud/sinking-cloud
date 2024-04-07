// @ts-ignore
import {request} from 'umi';

/** 获取配置列表 GET /admin/systems/config/list */
export async function getConfigList(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/config/list', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 添加配置 POST /admin/systems/config/create */
export async function createConfig(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/config/create', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 更新配置 POST /admin/systems/config/update */
export async function updateConfig(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/config/update', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 批量更新配置 POST /admin/systems/config/update */
export async function updateConfigs(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/config/updates', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 删除配置 POST /admin/systems/config/delete */
export async function deleteConfig(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/config/delete', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
