// @ts-ignore
import {request} from 'umi';

/** 登录接口 POST /auth/system/login */
export async function login(body: AdminAPI.LoginParams, options?: { [key: string]: any }) {
  return request<API.Response>('/auth/system/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
