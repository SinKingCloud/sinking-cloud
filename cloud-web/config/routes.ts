/**
 * 用户系统路由
 */
export const userPath = "user";
export const user = [
    {
        path: "index",
        title: "数据概览",
        name: "user.index",
        icon: 'icon-home',
        hideInMenu: false,
        component: "@/pages/user/index",
    },
];
export const adminPath = "admin";
export const admin = [
    {
        path: "index",
        component: "@/pages/admin/index",
        title: "网站概览",
        name: "admin.index",
        icon: 'icon-home',
    },
    {
        path: "system",
        title: "系统管理",
        name: "admin.system",
        icon: 'icon-setting',
        routes: [
            {
                path: "user",
                component: "@/pages/admin/index",
                title: "用户管理",
                name: "admin.system.user",
                icon: '•',
            },
            {
                path: "role",
                component: "@/pages/admin/index",
                title: "角色管理",
                name: "admin.system.role",
                icon: '•',
            },
            {
                path: "permission",
                component: "@/pages/admin/index",
                title: "权限管理",
                name: "admin.system.permission",
                icon: '•',
            },
            {
                path: "menu",
                component: "@/pages/admin/index",
                title: "菜单管理",
                name: "admin.system.menu",
                icon: '•',
            },
            {
                path: "task",
                component: "@/pages/admin/index",
                title: "任务管理",
                name: "admin.system.task",
                icon: '•',
            },
            {
                path: "storage",
                component: "@/pages/admin/index",
                title: "存储管理",
                name: "admin.system.storage",
                icon: '•',
            },
            {
                path: "config",
                component: "@/pages/admin/index",
                title: "系统设置",
                name: "admin.system.config",
                icon: '•',
            },
        ]
    },
    {
        path: "person",
        title: "账户管理",
        name: "admin.person",
        icon: 'icon-my',
        routes: [
            {
                path: "setting",
                component: "@/pages/admin/index",
                title: "账户管理",
                name: "admin.person.setting",
                icon: '•',
            },
            {
                path: "log",
                component: "@/pages/admin/index",
                title: "操作日志",
                name: "admin.person.log",
                icon: '•',
            },
        ]
    },
]
/**
 * 首页系统路由
 */
export const indexPath = "index";
export const index = [
    {
        path: "index",
        component: "@/pages/index/index",
        title: "系统首页",
        name: "index.index",
        icon: 'FileTextOutlined',
        hideInMenu: false,
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
        name: "redirect." + indexPath,
        title: index[0]?.title,
        routes: [
            {
                path: "/",
                component: index[0]?.component,
                title: index[0]?.title,
                name: index[0]?.name,
                hideInMenu: true,
            },
            {
                path: '/' + userPath + '/login',
                component: '@/pages/user/login',
                name: "user.login",
                title: "帐号登录",
                hideInMenu: true,
            },
            {
                path: '/' + adminPath + '/login',
                component: '@/pages/admin/login',
                name: "admin.login",
                title: "后台登录",
                hideInMenu: true,
            },
        ],
    },
    {
        path: indexPath,
        name: indexPath,
        title: "系统首页",
        routes: index,
    },
    /**
     * 用户路由
     */
    {
        path: "/" + userPath,
        name: "redirect." + userPath,
        redirect: '/' + userPath + '/' + (user[0]?.path || 'index'),
        hideInMenu: true,
    },
    {
        path: userPath,
        name: userPath,
        title: "用户系统",
        routes: user,
    },
    /**
     * admin路由
     */
    {
        path: "/" + adminPath,
        name: "redirect." + adminPath,
        redirect: '/' + adminPath + '/' + (admin[0]?.path || 'index'),
        hideInMenu: true,
    },
    {
        path: adminPath,
        name: adminPath,
        title: "后台网站",
        routes: admin,
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
