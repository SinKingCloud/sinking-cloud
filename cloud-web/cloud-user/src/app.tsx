// @ts-ignore
import React from 'react';
import type {Settings as LayoutSettings} from '@ant-design/pro-layout';
import {PageLoading} from '@ant-design/pro-layout';
// @ts-ignore
import type {RunTimeLayoutConfig} from 'umi';
// @ts-ignore
import {history, Link, RequestConfig} from 'umi';
import RightContent from '@/components/RightContent';
import Footer from '@/components/Footer';
import {queryCurrentUser} from './services/person/info';
import * as allIcons from '@ant-design/icons';
import {message} from 'antd';
import defaultSettings from '../config/defaultSettings';
import * as auth from './util/auth';

/**
 * 开发环境
 */
const isDev = process.env.NODE_ENV === 'development';

/**
 * 无需登陆的页面
 */
const allowPaths = ['/user/login'];//无需登陆页面


/** 获取用户信息展示一个 loading */
export const initialStateConfig = {
  // @ts-ignore
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

// ProLayout 支持的api https://procomponents.ant.design/components/layout
// @ts-ignore
export const layout: RunTimeLayoutConfig = ({initialState}) => {
  return {
    onPageChange: (location: { pathname: string; }) => {
      // 如果没有登录，重定向到 login
      if (!initialState?.currentUser && allowPaths.indexOf(location?.pathname as string) < 0) {
        history.push(allowPaths[0]);
      }
    },
    rightContentRender: () => <RightContent/>,
    disableContentMargin: false,
    waterMarkProps: {
      content: initialState?.currentUser?.nick_name,
    },
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
  if (ctx.res.code == 503 || ctx.res.code == 403) {
    history.push(allowPaths[0]);
    return message.warning(ctx.res.message);
  }
  if (ctx.res.code == 401) {
    message.warning(ctx.res.message);
  }
};
export const request: RequestConfig = {
  middlewares: [requestMiddleware],
  errorHandler,
};
