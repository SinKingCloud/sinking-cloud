// @ts-ignore
import {request} from 'umi';

/** 挂机列表 GET /open/user/uin/list */
export async function getUinList(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/uin/list', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 挂机地区列表 GET /open/user/uin/servers */
export async function getServerList(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/uin/servers', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}
