// @ts-ignore
import {request} from 'umi';

/** 获取通知列表 GET /admin/web/notice/list */
export async function getNoticeList(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/web/notice/list', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 获取通知列表 GET /admin/web/notice/info */
export async function getNoticeInfo(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/web/notice/info', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 更新通知 POST /admin/web/notice/update */
export async function updateNotice(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/web/notice/update', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 更新通知 POST /admin/web/notice/delete */
export async function deleteNotice(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/web/notice/delete', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 更新通知 POST /admin/web/notice/create */
export async function createNotice(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/web/notice/create', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
