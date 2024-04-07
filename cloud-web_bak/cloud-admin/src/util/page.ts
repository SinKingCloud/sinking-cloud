/**
 * table查询数据格式化
 * @param params 查询参数
 * @param sort 排序参数
 */
import {getControls} from "@/services/page/menu";
import defaultSettings from "../../config/defaultSettings";

export function getParams(params: any, sort: any): any {
  function getOrder(obj: any): any {
    for (const key in obj) {
      const temp = {};
      if (obj[key] == "ascend") {
        // @ts-ignore
        temp.order_by_type = "asc"
      } else {
        // @ts-ignore
        temp.order_by_type = "desc"
      }
      return temp;
    }
    return {};
  }

  params.page = params.current;
  params.page_size = params.pageSize;
  delete params.current;
  delete params.pageSize;
  return Object.assign(params, getOrder(sort));
}

/**
 * 获取页面控件
 */
export async function getPageControls(path: string = ""): Promise<any> {
  // eslint-disable-next-line no-param-reassign
  path = path == "" ? location.pathname : path;
  const key = path.replaceAll("/", ".").substring(defaultSettings.base.length);
  const response = await getControls({key: key});
  if (response.code == 200) {
    const data = {};
    response.data.forEach((k: any) => {
      data[k.key] = k;
    })
    return data;
  }
  return {};
}


