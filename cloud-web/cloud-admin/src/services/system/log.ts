// @ts-ignore
import {request} from 'umi';

/** 获取操作日志记录列表 GET /admin/systems/log/list */
export async function getLogList(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/log/list', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 删除操作记录 POST /admin/systems/log/delete */
export async function deleteLog(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/log/delete', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
