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
import {getQueryString} from "@/util/string";
import {queryCurrentWeb} from "@/services/common/info";

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
  currentWeb?: API.WebInfo;
  fetchUserInfo?: () => Promise<API.UserInfo | undefined>;
  fetchWebInfo?: () => Promise<API.WebInfo | undefined>;
}> {
  const fetchUserInfo = async () => {
    const {redirect, ticket} = history.location.query as { redirect: string, ticket: string };
    if (redirect) {
      localStorage.setItem("redirect", redirect || "");
    }
    if (ticket) {
      localStorage.setItem("ticket", ticket || "");
    }
    const redirectUrl = allowPaths[0] + getQueryString()
    try {
      const response = await queryCurrentUser();
      // @ts-ignore
      if (response.code != 200) {
        if (auth.getLoginToken() != null && auth.getLoginType() != null) {
          // @ts-ignore
          message.warning(response.message)
        }
        auth.deleteHeader()
        history.push(redirectUrl);
        return undefined;
      }
      return response.data;
    } catch (error) {
      history.push(redirectUrl);
      auth.deleteHeader()
    }
    return undefined;
  };
  const fetchWebInfo = async () => {
    const response = await queryCurrentWeb();
    // @ts-ignore
    if (response.code != 200) {
      return undefined;
    }
    return response.data;
  };
  const currentWeb = await fetchWebInfo();
  // 如果是登录页面，不执行
  if (allowPaths.indexOf(history.location.pathname) < 0) {
    const currentUser = await fetchUserInfo();
    return {
      fetchUserInfo,
      currentUser,
      fetchWebInfo,
      currentWeb,
      settings: {},
    };
  }
  return {
    fetchUserInfo,
    fetchWebInfo,
    currentWeb,
    settings: {},
  };
}

// ProLayout 支持的api https://procomponents.ant.design/components/layout
// @ts-ignore
export const layout: RunTimeLayoutConfig = ({initialState}) => {
  return {
    title: initialState?.currentWeb?.name || defaultSettings?.title,
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
    return message.error({content: `请求错误`, key: "errorHandler", duration: 3});
  }
  if (!response) {
    return message.error({content: `您的网络发生异常`, key: "errorHandler", duration: 3});
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
    let token = auth.getLoginToken();
    auth.deleteHeader();
    if (token != null && token != "") {
      return message.error({content: ctx.res.message || "登陆超时", key: "requestMiddleware", duration: 3});
    } else {
      return false;
    }
  }
  if (ctx.res.code == 401) {
    return message.error({content: ctx.res.message || "权限不足", key: "requestMiddleware", duration: 3});
  }
};
export const request: RequestConfig = {
  middlewares: [requestMiddleware],
  errorHandler,
};
