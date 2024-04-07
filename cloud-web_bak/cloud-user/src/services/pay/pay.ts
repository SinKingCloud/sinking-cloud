// @ts-ignore
import {request} from 'umi';

/** 消费记录 GET /user/pay/log */
export async function getPayLog(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/user/pay/log', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 账户充值 POST /user/pay/recharge */
export async function recharge(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/user/pay/recharge', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
