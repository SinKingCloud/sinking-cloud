import {matchRoutes} from 'umi';
import defaultSettings from "../config/defaultSettings";

export async function getInitialState(): Promise<{
    name?: string;
}> {
    return {
        name: "测试",
    }
}

/**
 * 路由跳转事件
 * @param clientRoutes 路由
 * @param location location对象
 */
export function onRouteChange({clientRoutes, location}) {
    const route = matchRoutes(clientRoutes, location.pathname)?.pop()?.route;
    if (route && route?.title) {
        document.title = route?.title + " - " + defaultSettings?.title;
    }
}