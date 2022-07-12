// @ts-ignore
import {request} from 'umi';

/** 获取枚举信息 GET /admin/page/enum */
export async function getEnum(params?:{},options?: { [key: string]: any }) {
  return request<API.Response>('/admin/page/enum', {
    method: 'GET',
    params:{
      ...params,
    },
    ...(options || {}),
  });
}
