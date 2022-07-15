// @ts-ignore
import {request} from 'umi';

/** 获取授权列表 GET /open/user/auth/list */
export async function getAuthList(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/auth/list', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 下载源码 POST /open/user/auth/download */
export async function Download(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/auth/download', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 获取商品列表 GET /open/user/auth/shop/list */
export async function getShopList(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/auth/shop/list', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 购买商品 GET /open/user/auth/shop/buy */
export async function shopBuy(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/auth/shop/buy', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
