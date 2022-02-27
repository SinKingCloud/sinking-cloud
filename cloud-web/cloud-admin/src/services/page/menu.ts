// @ts-ignore
import {request} from 'umi';

/** 获取菜单 GET /admin/page/menu */
export async function getMenus(options?: { [key: string]: any }) {
  return request<any>('/admin/page/menu', {
    method: 'GET',
    ...(options || {}),
  });
}
