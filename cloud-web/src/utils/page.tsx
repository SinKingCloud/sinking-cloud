/**
 * table查询数据格式化
 * @param params 查询参数
 * @param sort 排序参数
 */
export function getParams(params: any, sort: any): any {
    function getOrder(obj: any): any {
        for (const key in obj) {
            const temp: any = {
                order_by_field: key
            };
            if (obj[key] == "ascend") {
                temp.order_by_type = "asc";
            } else {
                temp.order_by_type = "desc";
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
 * 获取分页
 * @param params 参数
 * @param sort 排序
 * @param request 请求
 */
export async function getData(params: any, sort: any, request: any): Promise<any> {
    const fetchParams = getParams(params, sort)
    const data = await request({
        body: {
            ...fetchParams
        }
    });
    return {
        data: data?.data?.list === undefined || data?.data?.list === null || data?.data?.list?.length <= 0 ? [] : data?.data?.list,
        success: data?.code === 200,
        total: data?.data?.total,
    };
}

/** 空表格数据 */
export const emptyTableData = {
    data: [],
    success: true,
    total: 0,
};

/**
 * 生成日期范围搜索字段转换函数。
 * @param name 字段前缀
 */
export const dateRangeTransform = (name: string) => (value: any) => ({
    [`${name}_start`]: value?.[0]?.format ? value[0].format("YYYY-MM-DD HH:mm:ss") : value?.[0],
    [`${name}_end`]: value?.[1]?.format ? value[1].format("YYYY-MM-DD HH:mm:ss") : value?.[1],
});

/** 规范化 ID 数组 */
export const normalizeIds = (ids: any[] = []) => Array.from(new Set(
    (ids || []).map((item) => Number(item)).filter((item) => item > 0)
));

/**
 * 将枚举对象转换为 Select options。
 * @param data 枚举数据
 * @param numberValue value 是否转换为数字
 * @param labelRender 自定义 label
 */
export const buildOptions = (
    data: any,
    numberValue = false,
    labelRender?: (value: any, key: string) => any,
) => Object.entries(data || {}).map(([key, value]) => ({
    label: labelRender ? labelRender(value, key) : value,
    value: numberValue ? Number(key) : key,
}));

/**
 * 将枚举对象转换为 ProTable 的 valueEnum。
 * @param data 枚举数据
 * @param colors valueEnum 颜色映射
 */
export const buildValueEnum = (data: any, colors: Record<string, string> = {}) => Object.fromEntries(
    Object.entries(data || {}).map(([key, value]) => {
        const item: any = {text: value};
        if (colors[key]) {
            item.color = colors[key];
        }
        return [key, item];
    })
);
