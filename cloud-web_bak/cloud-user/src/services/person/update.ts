import {request} from "@@/plugin-request/request";

/** 修改资料 POST /admin/person/update/info */
export async function updateInfo(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/user/person/update/info', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 修改密码 POST /admin/person/update/password */
export async function updatePassword(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/user/person/update/password', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 修改手机 POST /admin/person/update/phone */
export async function updatePhone(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/user/person/update/phone', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
