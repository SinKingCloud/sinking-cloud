// @ts-ignore
import {request} from 'umi';

/** 获取用户订单列表 GET /admin/web/user_order/list */
export async function getUserOrderList(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/web/user_order/list', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 删除订单记录 POST /admin/web/user_order/delete */
export async function deleteUserOrder(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/web/user_order/delete', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
