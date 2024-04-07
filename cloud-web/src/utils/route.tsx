import {history} from 'umi';
import route, {user, userPath, admin, adminPath} from '../../config/routes'
import {Icon} from "@/components";
import React from "react";

/**
 * 递归获取完整路径
 * @param routes
 * @param name
 * @param params
 * @param currentPath
 */
function getPathByName(routes, name, params = {}, currentPath = ''): any {
    for (const route of routes) {
        let finalPath = currentPath + (currentPath && !currentPath.endsWith('/') ? '/' : '') + (route.path.startsWith('/') ? route.path.slice(1) : route.path);
        if (params && Object.keys(params).length > 0) {
            for (const key in params) {
                if (route.path.includes(`:${key}`)) {
                    finalPath = finalPath.replace(`:${key}`, params[key]);
                }
            }
        }
        if (route.name === name) {
            return finalPath.startsWith('/') ? finalPath : `/${finalPath}`;
        }
        if (route.routes) {
            const foundPath = getPathByName(route.routes, name, params, finalPath);
            if (foundPath) {
                return foundPath;
            }
        }
    }
    return null;
}

/**
 * 路由缓存
 */
const routeCache = {}

/**
 * 根据name获取路径
 * @param name 标识
 * @param params 参数
 */
export function getFullPathByName(name, params = {}): string {
    if (Object.keys(params) <= 0 && routeCache[name]) {
        return routeCache[name];
    }
    routeCache[name] = getPathByName(route, name, params, '');
    return routeCache[name];
}

/**
 * 跳转页面
 * @param name 标识
 * @param params 参数
 */
export function historyPush(name, params = {}) {
    history.push(getFullPathByName(name, params));
}

/**
 * 递归获取菜单
 * @param routes
 * @param parentPath
 */
function getMenuItems(routes, parentPath = '/') {
    return routes.map((route) => {
        if (route?.hideInMenu) {
            return;
        }
        const menuItem: any = {
            label: route?.title,
            path: route?.path,
        };
        if (route?.icon) {
            menuItem.icon = <Icon type={route?.icon}/>;
        }
        if (parentPath != '/' && !route.path.startsWith('/')) {
            menuItem.key = `${parentPath}/${route.path}`;
        } else {
            if (parentPath == '/' && route.path != '/') {
                menuItem.key = `/${route.path}`;
            } else {
                menuItem.key = `${route.path}`;
            }
        }
        if (route.routes) {
            const children = getMenuItems(route.routes, menuItem.key);
            if (children.length > 0) {
                menuItem.children = children;
            }
        }
        return menuItem;
    });
}

/**
 * 获取user菜单
 */
export function getUserMenuItems() {
    return getMenuItems(user, '/' + userPath);
}

/**
 * 获取admin菜单
 */
export function getAdminMenuItems() {
    return getMenuItems(admin, '/' + adminPath);
}

/**
 * 获取菜单
 */
export function getAllMenuItems() {
    return getMenuItems(route)
}