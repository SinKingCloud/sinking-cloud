import {request} from "@@/plugin-request/request";

/** 获取支付链接 GET /auth/pay/submit */
export async function paySubmit(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/auth/pay/submit', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}
