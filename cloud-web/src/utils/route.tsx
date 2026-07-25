import {history} from 'umi';
import route, {admin, adminPath, index, indexName, user, userPath,} from '../../config/routes'
import React from "react"
import {Icon} from "sinking-antd";

/**
 * 递归获取完整路径
 * @param routes
 * @param name
 * @param params
 * @param currentPath
 */
function getPathByName(routes: any, name: any, params: any = {}, currentPath = ''): any {
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
const routeCache: any = {}

/**
 * 根据name获取路径
 * @param name 标识
 * @param params 参数
 */
export function getFullPathByName(name: any, params: any = {}): string {
    if (Object.keys(params).length <= 0 && routeCache?.[name] != undefined) {
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
export function historyPush(name: any, params = {}) {
    history?.push(getFullPathByName(name, params) || name || "/");
}

/**
 * 递归获取菜单
 * @param routes
 * @param parentPath
 * @param hideMenu
 */
export function getMenuItems(routes: any, parentPath = '/', hideMenu = true) {
    return (routes || []).map((route: any) => {
        if (route?.path === "") {
            return undefined;
        }
        if (hideMenu && route?.hideInMenu) {
            return undefined;
        }
        const menuItem: any = {
            label: route?.title,
            path: route?.path,
            name: route?.name
        };
        if (route?.icon) {
            if (typeof route?.icon === 'string' && route?.icon != "0") {
                menuItem.icon = <Icon type={route?.icon}/>;
            }
        }
        if (parentPath != '/' && !route?.path?.startsWith('/')) {
            menuItem.key = `${parentPath}/${route.path}`;
        } else {
            if (parentPath == '/' && route.path != '/') {
                menuItem.key = `/${route.path}`;
            } else {
                menuItem.key = `${route.path}`;
            }
        }
        if (route?.routes) {
            const children = getMenuItems(route?.routes, menuItem.key, hideMenu).filter(Boolean);
            if (children.length > 0) {
                menuItem.children = children;
            } else {
                menuItem.children = []
            }
        }
        return menuItem;
    }).filter(Boolean);
}


const parentCache: any = {}


/**
 * 根据name查找父级菜单
 * @param data 菜单
 * @param name 标识
 * @param cacheKey 缓存来源标识
 */
export function getParentList(data: any[], name: string, cacheKey = "default"): any[] {
    if (!name) {
        return [];
    }
    if (!data || data.length <= 0) {
        return [];
    }
    const cacheName = `${cacheKey}:${name}`;
    if (parentCache?.[cacheName] != undefined) {
        return parentCache[cacheName];
    }
    let parents: any[] = [];

    function findParent(data: any, name: string): boolean {
        if (data?.children && data?.children?.length > 0) {
            if (data?.name === name) {
                parents.push(data);
                return true;
            }
        } else {
            for (const item of data) {
                if (item?.name === name) {
                    parents.push(item);
                    return true;
                }
                if (item?.children) {
                    if (findParent(item?.children, name)) {
                        parents.push(item);
                        return true;
                    }
                }
            }

        }
        return false;
    }

    findParent(data, name);
    parents = parents.reverse()
    if (parents.length > 0) {
        parentCache[cacheName] = parents
    }
    return parents;
}

/**
 * 当前访问系统path
 */
export function getCurrentPath(pathName: any): any {
    const mode = String(pathName || "").split("/").filter(Boolean)[0] || indexName;
    if (mode == userPath) {
        return '/' + userPath;
    }
    if (mode == adminPath) {
        return '/' + adminPath;
    }
    return '/' + indexName;
}

/**
 * 当前访问系统路由
 */
export function getCurrentMenus(pathName: any, hideMenu = true): any {
    const mode = String(pathName || "").split("/").filter(Boolean)[0] || indexName;
    if (mode == userPath) {
        return getMenuItems(user, '/' + userPath, hideMenu);
    }
    if (mode == adminPath) {
        return getMenuItems(admin, '/' + adminPath, hideMenu);
    }
    return getMenuItems(index, '/', hideMenu);
}


/**
 * 获取第一个菜单(排除children)
 * @param menu
 */
export function getFirstMenuWithoutChildren(menu: any[]): any | null {
    for (const item of menu) {
        if (!item.hasOwnProperty('children')) {
            return item; // 如果对象没有 children 属性，则返回该对象
        } else {
            const result = getFirstMenuWithoutChildren(item.children);
            if (result !== null) {
                return result;
            }
        }
    }
    return null;
}

/**
 * 清理路由缓存
 */
export function clearRouteCache() {
    Object.keys(parentCache).forEach((key) => {
        delete parentCache[key];
    });
    Object.keys(routeCache).forEach((key) => {
        delete routeCache[key];
    });
}
