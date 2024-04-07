// @ts-ignore
import {request} from 'umi';

/** 挂机统计 GET /open/user/uin/count */
export async function getUinCount(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/uin/count', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 挂机统计 GET /open/user/uin/self */
export async function getSelfInfo(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/uin/self', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 活动信息 GET /open/user/uin/activity */
export async function getActivityInfo(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/uin/activity', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 开通等级 POST /open/user/uin/shop */
export async function shopLevel(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/uin/shop', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

export async function getAuthInfo(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/uin/auth', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}
