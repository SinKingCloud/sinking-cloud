import {extend} from 'umi-request';
import defaultSettings from "@/../config/defaultSettings";
import {getHeaders} from "@/utils/auth";
import {historyPush} from "@/utils/route";
import {API} from "../../typings";

/**
 * request对象
 */
const request = extend({
    prefix: defaultSettings?.gateway || '',
    timeout: 60000,
    errorHandler: errorHandle,
    headers: getHeaders() || {},
});

/**
 * 异常处理
 */
function errorHandle(): void {
    return historyPush("error");
}

/**
 * 中间件
 * @param ctx context
 * @param next 执行函数
 */
const check = async (ctx: any, next: any) => {
    await next();
    if (ctx.res.code == 401) {
        historyPush("notAllowed");//无权限页面;
    }
}
request.use(check);

/**
 * 发送post请求
 * @param url 请求地址
 * @param body 请求载荷
 * @param onSuccess code200回调
 * @param onFail code500回调
 * @param onFinally 最后执行
 * @param onError 错误捕获
 * @param options 其他选项
 */
export async function post(url: string = "", body: any = {}, onSuccess: (res) => void = undefined, onFail: (res) => void = undefined, onFinally: () => void = undefined, onError: (error) => void = undefined, options?: {
    [key: string]: any
}) {
    let res = {};
    await request<API.Response>(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        data: body,
        ...(options || {}),
    }).then((r) => {
        res = r;
        if (r?.code == 200) {
            onSuccess?.(r);
        } else {
            onFail?.(r);
        }
    }).catch((error) => {
        onError?.(error);
    }).finally(() => {
        onFinally?.();
    });
    return res;
}

/**
 * 发送post请求
 * @param url 请求地址
 * @param params 请求参数
 * @param onSuccess code200回调
 * @param onFail code500回调
 * @param onFinally 最后执行
 * @param onError 错误捕获
 * @param options 其他选项
 */
export async function get(url: string = "", params?: {}, onSuccess: (res) => void = undefined, onFail: (res) => void = undefined, onFinally: () => void = undefined, onError: (error) => void = undefined, options?: {
    [key: string]: any
}) {
    let res = {};
    await request<API.Response>(url, {
        method: 'GET',
        params: {
            ...params,
        },
        ...(options || {}),
    }).then((r) => {
        res = r;
        if (r?.code == 200) {
            onSuccess?.(r);
        } else {
            onFail?.(r);
        }
    }).catch((error) => {
        onError?.(error);
    }).finally(() => {
        onFinally?.();
    });
    return res;
}

export default request;
