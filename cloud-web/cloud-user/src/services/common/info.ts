// @ts-ignore
import {request} from "umi";

/** 获取当前站点信息 GET /auth/web/info */
export async function queryCurrentWeb(options?: { [key: string]: any }) {
  return request<{
    data: API.WebInfo;
  }>('/auth/web/info', {
    method: 'GET',
    ...(options || {}),
  });
}
