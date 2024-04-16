import Layout from "@/layouts/components";
import {getUserMenuItems} from "@/utils/route";
import defaultSettings from "@/../config/defaultSettings";
import {Doc} from "@/components/icon";
import {useModel} from "umi";
import {useEffect, useState} from "react";
import {deleteHeader} from "@/utils/auth";
import {historyPush} from "@/utils/route";
import request from "@/utils/request";

/**
 * 中间件
 * @param ctx context
 * @param next 执行函数
 */
const check = async (ctx, next) => {
    await next();
    if (ctx.res.code == 503 || ctx.res.code == 403) {
        deleteHeader();
        historyPush("user.login");
    }
}
request.use(check);

/**
 * 用户系统
 */
export default () => {

    /**
     * 全局信息
     */
    const web = useModel("web");//网站信息
    const user = useModel("user");//用户信息

    /**
     * 初始化用户信息
     */
    const [loading, setLoading] = useState(true);
    const initUser = () => {
        setLoading(true);
        user?.getWebUser()?.then((u) => {
            user?.setWeb(u);
        })?.finally(() => {
            setLoading(false);
        });
    }

    useEffect(() => {
        initUser();
    }, []);

    return (
        <Layout loading={loading}
                menus={getUserMenuItems()}
                footer={<>©{new Date().getFullYear()} All Right Revered {web?.info?.name || defaultSettings?.title}</>}
                headerLeft={<>左侧内容</>}
                headerRight={<>右侧内容</>}
                menuBottomBtnIcon={Doc}
                menuBottomBtnText={"开发文档"}
                onMenuBottomBtnClick={() => {
                    console.log("点击了底部按钮")
                }}
                menuCollapsedWidth={60}
                menuUnCollapsedWidth={210}
                onLogoClick={() => {
                    console.log("点击了logo")
                }}
                onMenuClick={(item) => {
                    console.log("点击了菜单", item)
                }}
                collapsedLogo={() => {
                    return <img src="/logo.svg" width="39"
                                alt={defaultSettings?.title}/>;
                }}
                unCollapsedLogo={(isLight) => {
                    return <img src={isLight ? "/logo_light.png" : "/logo_dark.png"} width="130"
                                alt={defaultSettings?.title}/>;
                }}
                onMenuBtnClick={(state) => {
                    console.log("点击了菜单按钮", state)
                }}
        />
    );
}