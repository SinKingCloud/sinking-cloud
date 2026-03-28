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
        avatar?: string;
        nick_name?: string;
        account?: string;
        phone?: string;
        email?: string;
        real_auth?: boolean;
        session?: Session;
    }

    /**
     * 用户信息
     */
    type SystemUser = {
        id?: number;
        avatar?: string;
        nick_name?: string;
        account?: string;
        phone?: string;
        email?: string;
        session?: Session;
    }

    /**
     * 会话信息
     */
    type Session = {
        device?: string;
        ip?: string;
        location?: string;
        login_time?: string;
        expire_time?: string;
    }

    /**
     * 网站信息
     */
    type WebInfo = {
        name?: string;
        title?: string;
        keywords?: string;
        describe?: string;
        contact?: string;
        url?: string;
        logo?: string;
        ui?: Ui,
    }

    /**
     * 界面信息
     */
    type Ui = {
        water_mark?: string;
        layout?: string;
        theme?: string;
        color?: string;
        compact?: boolean;
        radius?: number;
    }
}

