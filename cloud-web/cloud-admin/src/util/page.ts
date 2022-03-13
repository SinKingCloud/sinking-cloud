/**
 * table查询数据格式化
 * @param params 查询参数
 * @param sort 排序参数
 */
import {getControls} from "@/services/page/menu";

export function getParams(params: any, sort: any): any {
  function getOrder(obj: any): any {
    for (let key in obj) {
      let temp = {
        order_by_field: key
      };
      if (obj[key] == "ascend") {
        temp["order_by_type"] = "asc"
      } else {
        temp["order_by_type"] = "desc"
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
export async function getPageControls(): Promise<any> {
  let key = location.pathname.replaceAll("/", ".").substring(1);
  let response = await getControls({key: key});
  if(response.code==200){
    let data = {};
    response.data.forEach((k:any)=>{
      data[k.key] = k;
    })
    return data;
  }
  return {};
}
