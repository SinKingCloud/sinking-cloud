import 'umi/typings';

declare namespace API {

    /**
     * 快捷请求结构
     */
    type RequestParams = {
        body?: any,
        onSuccess?: (res) => void,
        onFail?: (res) => void,
        onFinally?: () => void
    }

    /**
     * 通用响应
     */
    type Response = {
        code?: number;
        message?: string;
        data?: any;
        request_id?: string;
    }
    /**
     * 站点信息
     */
    type WebInfo = {
        name?: string;
        contact?: string;
    }

    /**
     * 网站用户信息
     */
    type WebUser = {
        id?: number;
        account?: string;
        money?: string;
        nick_name?: string;
        avatar?: string;
        email?: string;
        phone?: string;
        login_ip?: string;
        login_time?: string;
    }

    /**
     * 系统用户信息
     */
    type SystemUser = {
        id?: number;
        account?: string;
        real_auth?: string;
        nick_name?: string;
        avatar?: string;
        email?: string;
        phone?: string;
        login_ip?: string;
        login_time?: string;
    }
}

