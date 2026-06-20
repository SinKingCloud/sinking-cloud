import React, {useEffect, useState} from "react";
import {useLocation, useModel, history} from "umi";
import {ConfigProvider, Spin} from "antd";
import {getCurrentPath, historyPush} from "@/utils/route";
import {adminPath, userPath} from "../../config/routes";
import {Admin, User, Index} from "@/layouts/components";
import Title from "./components/title";
import {getHeaders} from "@/utils/auth";
import request from "@/utils/request";
import defaultSettings from "../../config/defaultSettings";
import {setIconfontUrl, Theme, useTheme} from "sinking-antd";
import {createStyles} from "antd-style";
import zhCN from 'antd/locale/zh_CN';

/**
 * 中间件
 * @param ctx context
 * @param next 执行函数
 */
const check = async (ctx: any, next: any) => {
    ctx.req.options.headers = getHeaders();
    await next();
    const path = getCurrentPath((history as any)?.location?.pathname || "/")?.substring(1);
    if (path == userPath || path == adminPath) {
        if ((ctx.res.code == 403)) {
            historyPush(path == userPath ? "user.login" : "admin.login");//登陆页面
            return
        }
        if (ctx.res.code == 401) {
            historyPush("notAllowed");
            return
        }
    }
}
request.use(check);

/**
 * 样式信息
 */
const useStyles = createStyles((): any => {
    return {
        load: {
            width: "100%",
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            height: "80vh",
        },
    };
});

const ProLayout = () => {
    const {styles: {load}} = useStyles();
    const theme = useTheme();
    const web = useModel("web");
    const location = useLocation();
    const [systemPath, setSystemPath] = useState("");
    const [isLoginPage, setIsLoginPage] = useState(false);
    const currentSystemPath = getCurrentPath(location?.pathname)?.substring(1);
    const ui = currentSystemPath == adminPath ? web?.info?.ui?.admin : web?.info?.ui?.user;

    // 初始化 iconfont 地址
    useEffect(() => {
        if (defaultSettings?.iconfontUrl) {
            setIconfontUrl(defaultSettings.iconfontUrl);
        }
    }, []);

    /**
     * 监听路由数据
     */
    useEffect(() => {
        const path = currentSystemPath;
        let loginPage = false
        if (location?.pathname?.startsWith("/" + path + "/login")) {
            setIsLoginPage(true);
            loginPage = true;
        }
        if (loginPage != isLoginPage) {
            setIsLoginPage(loginPage);
        }
        if (path != systemPath) {
            setSystemPath(path);
        }
    }, [location?.pathname]);

    // 初始化主题
    useEffect(() => {
        if (ui) {
            if (ui.color) {
                theme?.setColor(ui.color);
            }
            if (ui.radius >= 0) {
                theme?.setRadius(ui.radius <= 15 ? ui.radius : 0);
            }
            if (ui.compact) {
                theme?.setCompactTheme();
            } else {
                theme?.setDefaultTheme();
            }
        }
    }, [ui]);

    if (!web?.info || !ui) {
        return <Spin spinning={true} size="large" className={load}/>;
    }
    if (!isLoginPage && systemPath == userPath) {
        return <User/>
    }
    if (!isLoginPage && systemPath == adminPath) {
        return <Admin/>
    }
    return <Index/>;
}

export default () => {
    return (
        <ConfigProvider locale={zhCN}>
            <Theme>
                <Title/>
                <ProLayout/>
            </Theme>
        </ConfigProvider>
    );
}
