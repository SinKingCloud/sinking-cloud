import {request} from "@@/plugin-request/request";

/** 修改资料 POST /admin/person/update */
export async function updateInfo(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/person/update', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 修改密码 POST /admin/person/password */
export async function updatePassword(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/person/password', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
