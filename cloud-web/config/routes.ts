/**
 * 用户系统路由
 */
import {AdminPageCode} from "../src/utils/permission";

export const userPath = "user";
export const user = [
    {
        path: "index",
        title: "系统首页",
        name: "user.index",
        icon: 'icon-home',
        hideInMenu: false,
        component: "@/pages/user/index",
    },
    {
        path: "mind",
        title: "思维导图",
        name: "user.mind",
        icon: 'icon-order',
        hideInMenu: false,
        routes: [
            {
                path: "index",
                component: "@/pages/user/mind",
                title: "我的导图",
                name: "user.mind.index",
                icon: '•',
            },
            {
                path: "",
                redirect: "index",
            },
        ]
    },
    {
        path: "person",
        title: "账户管理",
        name: "user.person",
        icon: 'icon-account',
        hideInMenu: false,
        routes: [
            {
                path: "info",
                component: "@/pages/user/person/info",
                title: "账户信息",
                name: "user.person.info",
                icon: '•',
            },
            {
                path: "pay",
                component: "@/pages/user/person/pay",
                title: "财务管理",
                name: "user.person.pay",
                icon: '•',
            },
            {
                path: "certificate",
                component: "@/pages/user/person/certificate",
                title: "实名认证",
                name: "user.person.certificate",
                icon: '•',
            },
            {
                path: "session",
                component: "@/pages/user/person/session",
                title: "会话管理",
                name: "user.person.session",
                icon: '•',
            },
            {
                path: "message",
                component: "@/pages/user/person/message",
                title: "站内消息",
                name: "user.person.message",
                icon: '•',
            },
            {
                path: "",
                redirect: "info",
            },
        ]
    },
];
export const adminPath = "admin";
export const admin = [
    {
        path: "index",
        component: "@/pages/admin/index",
        title: "系统首页",
        icon: "icon-home",
        name: AdminPageCode.INDEX,
    },
    {
        path: "web",
        title: "网站管理",
        name: AdminPageCode.WEB,
        routes: [
            {
                path: "user",
                component: "@/pages/admin/user",
                title: "用户管理",
                name: AdminPageCode.WEB_USER,
            },
            {
                path: "order",
                component: "@/pages/admin/order",
                title: "订单管理",
                name: AdminPageCode.WEB_ORDER,
            },
            {
                path: "",
                redirect: "user",
            },
        ]
    },
    {
        path: "system",
        title: "系统管理",
        name: AdminPageCode.SYSTEM,
        routes: [
            {
                path: "master",
                component: "@/pages/admin/master",
                title: "员工管理",
                name: AdminPageCode.SYSTEM_MASTER,
            },
            {
                path: "role",
                component: "@/pages/admin/role",
                title: "角色管理",
                name: AdminPageCode.SYSTEM_ROLE,
            },
            {
                path: "permission",
                component: "@/pages/admin/permission",
                title: "权限管理",
                name: AdminPageCode.SYSTEM_PERMISSION,
            },
            {
                path: "menu",
                component: "@/pages/admin/menu",
                title: "菜单管理",
                name: AdminPageCode.SYSTEM_MENU,
            },
            {
                path: "cron",
                component: "@/pages/admin/cron",
                title: "任务管理",
                name: AdminPageCode.SYSTEM_CRON,
            },
            {
                path: "driver",
                component: "@/pages/admin/driver",
                title: "存储管理",
                name: AdminPageCode.SYSTEM_DRIVER,
            },
            {
                path: "pay",
                component: "@/pages/admin/pay",
                title: "支付管理",
                name: AdminPageCode.SYSTEM_PAY,
            },
            {
                path: "config",
                component: "@/pages/admin/config",
                title: "系统设置",
                name: AdminPageCode.SYSTEM_CONFIG,
            },
            {
                path: "",
                redirect: "master",
            },
        ]
    },
    {
        path: "person",
        title: "账户管理",
        name: "admin.person",
        routes: [
            {
                path: "info",
                component: "@/pages/admin/person/info",
                title: "账户信息",
                name: "admin.person.info",
            },
            {
                path: "session",
                component: "@/pages/admin/person/session",
                title: "会话管理",
                name: "admin.person.session",
            },
            {
                path: "",
                redirect: "info",
            },
        ]
    },
]
/**
 * 首页系统路由
 */
export const indexName = "index";
export const index = [
    {
        path: "/",
        component: "@/pages/index/index",
        title: "系统首页",
        name: "index.index",
        icon: 'FileTextOutlined',
        hideInMenu: false,
    },
    {
        path: "mind/:code",
        component: "@/pages/index/mind",
        name: "index.mind",
        title: "导图分享",
        hideInMenu: true,
    },
];
/**
 * 系统路由
 */
export default [
    /**
     * 首页路由
     */
    {
        path: "/",
        name: indexName,
        title: index[0]?.title,
        routes: index,
    },
    /**
     * 用户路由
     */
    {
        path: userPath,
        name: userPath,
        title: "用户系统",
        routes: [
            {
                path: "login",
                component: "@/pages/user/login",
                name: "user.login",
                title: "帐号登录",
                hideInMenu: true,
            },
            ...user,
            {
                path: "",
                redirect: user[0]?.path || "index",
            },
        ],
    },
    /**
     * admin路由
     */
    {
        path: adminPath,
        name: adminPath,
        title: "后台网站",
        routes: [
            {
                path: "login",
                component: "@/pages/admin/login",
                name: "admin.login",
                title: "后台登录",
                hideInMenu: true,
            },
            ...admin,
            {
                path: "",
                redirect: admin[0]?.path || "index",
            },
        ],
    },
    /**
     * 支付结果
     */
    {
        path: "pay",
        component: "@/pages/pay.tsx",
        name: "pay",
        title: "支付结果",
        hideInMenu: true,
    },
    /**
     * 500
     */
    {
        path: '/500',
        component: '@/pages/500.tsx',
        name: "error",
        title: "服务器错误",
        hideInMenu: true,
    },
    /**
     * 403
     */
    {
        path: '/403',
        component: '@/pages/403.tsx',
        name: "notAllowed",
        title: "无权限",
        hideInMenu: true,
    },
    /**
     * 404
     */
    {
        path: '/*',
        component: '@/pages/404.tsx',
        name: "notFound",
        title: "页面不存在",
        hideInMenu: true,
    }
];
