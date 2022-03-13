// @ts-ignore
import {request} from 'umi';

/** 获取菜单 GET /admin/page/menu */
export async function getMenus(options?: { [key: string]: any }) {
  return request<any>('/admin/page/menu', {
    method: 'GET',
    ...(options || {}),
  });
}


/** 获取控件列表 GET /admin/page/control */
export async function getControls(params?:{},options?: { [key: string]: any }) {
  return request<any>('/admin/page/control', {
    method: 'GET',
    params:{
      ...params,
    },
    ...(options || {}),
  });
}
