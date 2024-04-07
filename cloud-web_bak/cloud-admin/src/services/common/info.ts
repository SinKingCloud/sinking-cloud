// @ts-ignore
import {request} from "umi";

/** 获取当前站点信息 GET /auth/system/info */
export async function queryCurrentWeb(options?: { [key: string]: any }) {
  return request<{
    data: API.WebInfo;
  }>('/auth/system/info', {
    method: 'GET',
    ...(options || {}),
  });
}
