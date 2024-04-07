// @ts-ignore
import {request} from 'umi';

/** 获取附件记录列表 GET /admin/systems/file/list */
export async function getFileList(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/file/list', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 删除附件记录 POST /admin/systems/file/delete */
export async function deleteFile(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/systems/file/delete', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
