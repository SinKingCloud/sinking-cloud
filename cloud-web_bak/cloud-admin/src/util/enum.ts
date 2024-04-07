import {getEnum} from "@/services/page/enum";
import {implode} from "@/util/string";

/**
 * 枚举数据
 */
export enum Enums {
  WebNotice = "web.notice",//公告
}

/**
 * 获取枚举数据
 * @param key 标识key
 */
export async function getEnums(key: any = []): Promise<any> {
  const response = await getEnum({keys: implode("|", key)});
  if (response.code == 200) {
    return response?.data || {};
  }
  return {};
}
