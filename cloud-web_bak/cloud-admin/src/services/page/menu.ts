// @ts-ignore
import {request} from 'umi';

/** 获取菜单 GET /admin/page/menu */
// eslint-disable-next-line @typescript-eslint/ban-types
export async function getMenus(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/page/menu', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}


/** 获取控件列表 GET /admin/page/control */
// eslint-disable-next-line @typescript-eslint/ban-types
export async function getControls(params?: {}, options?: { [key: string]: any }) {
  return request<API.Response>('/admin/page/control', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}
