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

/** 修改挂机地区 POST /open/user/uin/operate/server */
export async function changeUinArea(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/uin/operate/server', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 修改挂机时间 POST /open/user/uin/operate/timing */
export async function changeUinTiming(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/uin/operate/timing', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 修改挂机密码 POST /open/user/uin/operate/password */
export async function changeUinPassword(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/uin/operate/password', {
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

/** 修改挂机状态 POST /open/user/uin/operate/status */
export async function uinStatus(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/uin/operate/status', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 开通挂机 POST /open/user/uin/buy */
export async function uinBuy(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/uin/buy', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 登陆账号 POST /open/user/uin/operate/login */
export async function uinLogin(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/uin/operate/login', {
    method: 'POST',
    timeout: 600000,
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 扫码验证 POST /open/user/uin/operate/qr_login */
export async function uinQrLogin(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/uin/operate/qr_login', {
    method: 'POST',
    timeout: 600000,
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 短信验证 POST /open/user/uin/operate/sms */
export async function uinSms(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/uin/operate/sms', {
    method: 'POST',
    timeout: 600000,
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 重置挂机设备 POST /open/user/uin/operate/reset */
export async function uinReset(body: any, options?: { [key: string]: any }) {
  return request<API.Response>('/open/user/uin/operate/reset', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
