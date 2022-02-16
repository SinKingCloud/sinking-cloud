import React from 'react';
import type {Settings as LayoutSettings} from '@ant-design/pro-layout';
import {PageLoading} from '@ant-design/pro-layout';
import type {RunTimeLayoutConfig} from 'umi';
import {history, Link, RequestConfig} from 'umi';
import RightContent from '@/components/RightContent';
import Footer from '@/components/Footer';
import {queryCurrentUser} from './services/user/info';
import * as allIcons from '@ant-design/icons';
import {message} from 'antd';
import defaultSettings from '../config/defaultSettings';
import * as auth from './util/auth';
import {MenuDataItem} from "@umijs/route-utils";

/**
 * 开发环境
 */
const isDev = process.env.NODE_ENV === 'development';

/**
 * 无需登陆的页面
 */
const allowPaths = ['/admin/login'];//无需登陆页面

/** 获取用户信息比较慢的时候会展示一个 loading */
export const initialStateConfig = {
  loading: <PageLoading/>,
};

/**
 * 初始化数据
 * @see  https://umijs.org/zh-CN/plugins/plugin-initial-state
 * */
export async function getInitialState(): Promise<{
  settings?: Partial<LayoutSettings>;
  currentUser?: API.UserInfo;
  fetchUserInfo?: () => Promise<API.UserInfo | undefined>;
}> {
  const fetchUserInfo = async () => {
    try {
      const response = await queryCurrentUser();
      // @ts-ignore
      if (response.code != 200) {
        if (auth.getLoginToken() != null && auth.getLoginType() != null) {
          // @ts-ignore
          message.warning(response.message)
        }
        auth.deleteHeader()
        history.push(allowPaths[0]);
        return undefined;
      }
      return response.data;
    } catch (error) {
      history.push(allowPaths[0]);
      auth.deleteHeader()
    }
    return undefined;
  };
  // 如果是登录页面，不执行
  if (allowPaths.indexOf(history.location.pathname) < 0) {
    const currentUser = await fetchUserInfo();
    return {
      fetchUserInfo,
      currentUser,
      settings: {},
    };
  }
  return {
    fetchUserInfo,
    settings: {},
  };
}

/**
 * 菜单图标映射
 * @param menus 菜单
 * @param iconType 图标类型
 */
const fixMenuItemIcon = (menus: MenuDataItem[], iconType = 'Outlined'): MenuDataItem[] => {
  menus.forEach((item) => {
    const {icon, children} = item;
    if (typeof icon === 'string') {
      let fixIconName = icon.slice(0, 1).toLocaleUpperCase() + icon.slice(1) + iconType;
      item.icon = React.createElement(allIcons[fixIconName] || allIcons[icon]);
    }
    children && children.length > 0 ? (item.children = fixMenuItemIcon(children)) : null;
  });
  return menus;
};

// ProLayout 支持的api https://procomponents.ant.design/components/layout
// @ts-ignore
export const layout: RunTimeLayoutConfig = ({initialState}) => {
  return {
    menu: {
      locale: false,
      params: {
        userId: initialState?.currentUser?.account,
      },
      request: async (params: Record<string, any>, defaultMenuData: MenuDataItem[]) => {
        return [{
          path: '/index',
          name: '用户首页',
          icon: 'smile',
          routes: [
            {
              path: 'index',
              name: '示例首页12',
            },
          ],
        }];
      },
    },
    onPageChange: () => {
      const {location} = history;
      // 如果没有登录，重定向到 login
      if (!initialState?.currentUser && allowPaths.indexOf(location.pathname) < 0) {
        history.push(allowPaths[0]);
      }
    },
    rightContentRender: () => <RightContent/>,
    disableContentMargin: false,
    waterMarkProps: {
      content: initialState?.currentUser?.nick_name,
    },
    menuDataRender: fixMenuItemIcon,
    footerRender: () => <Footer/>,
    links: isDev
      ? [
        <Link to="/umi/plugin/openapi" target="_blank">
          <allIcons.LinkOutlined/>
          <span>OpenAPI 文档</span>
        </Link>,
        <Link to="/~docs">
          <allIcons.BookOutlined/>
          <span>业务组件文档</span>
        </Link>,
      ]
      : [],
    menuHeaderRender: undefined,
    // 自定义 403 页面
    // unAccessible: <div>unAccessible</div>,
    ...initialState?.settings,
  };
};

/**
 * 拦截器错误处理
 * @param error 错误信息
 */
const errorHandler = (error: { response: RequestConfig; }) => {
  const {response} = error;
  if (response && response.code) {
    return message.error(`请求错误`);
  }
  if (!response) {
    return message.error('您的网络发生异常');
  }
  if (response.code == 503 || response.code == 403) {
    history.push(allowPaths[0]);
    return message.warning(response.message);
  }
  if (response.code == 401) {
    message.warning(response.message);
  }
  return response;
};
/**
 * 拦截器
 * @param ctx 请求对象
 * @param next 执行
 */
const requestMiddleware = async (ctx: any, next: () => void) => {
  ctx.req.url = defaultSettings.api + ctx.req.url;
  ctx.req.options.headers = auth.getHeaders();
  await next();
};
export const request: RequestConfig = {
  errorHandler,
  middlewares: [requestMiddleware],
};
