import {extend} from 'umi-request';
import defaultSettings from "@/../config/defaultSettings";
import {message} from "@/components";

/**
 * request对象
 */
export const request = extend({
    prefix: defaultSettings?.gateway || '',
    timeout: 60000,
    errorHandler: errorHandle,
});

/**
 * 异常处理
 * @param error 错误信息
 */
function errorHandle(error): any {
    const {response} = error;
    if (response && response.code) {
        return message.error({content: `请求错误`, key: "errorHandler", duration: 3});
    }
    if (!response) {
        return message.error({content: `您的网络发生异常`, key: "errorHandler", duration: 3});
    }
    return response;
}

/**
 * 中间件
 */
request.use(async (ctx, next) => {
    await next();
});
