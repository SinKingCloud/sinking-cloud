import {matchRoutes} from 'umi';
import defaultSettings from "../config/defaultSettings";
import {initDefaultSetting} from "@/utils/init";

/**
 * 初始化defaultSetting
 */
await initDefaultSetting();

/**
 * 路由跳转事件
 * @param props 传递参数
 */
export function onRouteChange(props: any) {
    const route = matchRoutes(props?.clientRoutes, props?.location?.pathname)?.pop()?.route;
    if (route && route?.title) {
        document.title = route?.title + " - " + (defaultSettings?.name || defaultSettings?.title);
    }
}