import {extend} from 'umi-request';
import defaultSettings from "@/../config/defaultSettings";
import {getHeaders} from "@/utils/auth";
import {historyPush} from "@/utils/route";

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

export default request;
