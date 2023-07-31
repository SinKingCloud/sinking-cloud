import {request} from "@@/plugin-request/request";

/** 获取推送令牌 GET /auth/web/push */
export async function genPushTicket(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/auth/system/push', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}
