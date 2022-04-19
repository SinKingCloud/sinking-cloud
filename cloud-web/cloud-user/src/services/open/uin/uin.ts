// @ts-ignore
import {request} from 'umi';

/** 挂机列表 GET /open/user/uin/list */
export async function getUinList(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/uin/list', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 挂机地区列表 GET /open/user/uin/servers */
export async function getServerList(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/uin/servers', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 项目列表 GET /open/user/uin/project */
export async function getProject(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/uin/project', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 等级列表 GET /open/user/uin/levels */
export async function getLevelList(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/uin/levels', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 挂机信息 POST /open/user/uin/operate/info */
export async function getUinInfo(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/uin/operate/info', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 修改功能 POST /open/user/uin/operate/setting */
export async function changeUinSetting(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/uin/operate/setting', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 删除挂机 POST /open/user/uin/operate/delete */
export async function deleteUin(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/uin/operate/delete', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 提交挂机 POST /open/user/uin/operate/run */
export async function uinRun(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/uin/operate/run', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
