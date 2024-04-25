import React from "react";
import Layout from "@/layouts/components";
import {getAdminMenuItems, historyPush} from "@/utils/route";
import defaultSettings from "../../config/defaultSettings";
import {deleteHeader} from "@/utils/auth";
import request from "@/utils/request";


/**
 * 中间件
 * @param ctx context
 * @param next 执行函数
 */
const check = async (ctx: any, next: any) => {
    await next();
    if (ctx.res.code == 503 || ctx.res.code == 403) {
        deleteHeader();
        historyPush("admin.login");
    }
}
request.use(check);

export default () => {
    return (
        <Layout menus={getAdminMenuItems()}
                footer={<>©{new Date().getFullYear()} All Right Revered {defaultSettings?.title}</>}
                headerLeft={<>左侧内容</>}
                headerRight={<>右侧内容</>}
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