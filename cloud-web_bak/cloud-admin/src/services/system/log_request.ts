// @ts-ignore
import {request} from 'umi';

/** 获取请求记录列表 GET /admin/systems/log_request/list */
export async function getLogRequestList(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/log_request/list', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 获取请求记录详情 GET /admin/systems/log_request/info */
export async function getLogRequestInfo(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/log_request/info', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 删除请求记录 POST /admin/systems/log_request/delete */
export async function deleteLogRequest(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/log_request/delete', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
