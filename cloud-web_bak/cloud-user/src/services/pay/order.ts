// @ts-ignore
import {request} from 'umi';

/** 订单记录 GET /user/pay/order */
export async function getPayOrder(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/user/pay/order', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}
