// @ts-ignore
import {request} from 'umi';

/** 获取实名信息 GET /user/person/log */
export async function getLogList(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/user/person/log', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}
