/**
 * 后台系统路由
 */
export const adminPath = "admin";
export const admin = [
    {
        path: 'test',
        component: '@/pages/user/index',
        icon: 'icon-doc',
        title: "测试二级",
        name: "testtwo11",
        hideInMenu: false,
    },
    {
        path: 'admin',
        component: '@/pages/user/docs',
        icon: 'icon-doc',
        title: "测试二级",
        name: "testtwo22",
        hideInMenu: false,
    },
];

/**
 * 用户系统路由
 */
export const userPath = "user";
export const user = [
    {
        path: "index",
        component: "@/pages/user/index",
        title: "系统首页",
        name: "user.index",
        icon: 'icon-doc',
        hideInMenu: false,
    },
    {
        path: "docs",
        component: "@/pages/user/docs",
        title: "开发文档",
        name: "doc",
        icon: 'icon-doc',
        hideInMenu: false,
    },
    {
        path: "test",
        title: "测试菜单",
        name: "doc12",
        icon: 'icon-doc',
        hideInMenu: false,
        routes: [
            {
                path: "index",
                component: "@/pages/user/index",
                title: "系统首页",
                name: "user.index12",
                hideInMenu: false,
            },
            {
                path: "docs",
                title: "开发文档",
                name: "doc1234s",
                hideInMenu: false,
                routes: [
                    {
                        path: "index",
                        component: "@/pages/user/index",
                        title: "系统首页",
                        name: "user.index124",
                        hideInMenu: false,
                    },
                    {
                        path: "docs",
                        component: "@/pages/user/docs",
                        title: "开发文档",
                        name: "doc1234",
                        hideInMenu: false,
                    },
                ]
            },
        ]
    },
];

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
        icon: 'icon-doc',
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
        component: index[0]?.component,
        layout: false,
    },
    {
        path: indexPath,
        name: indexPath,
        title: "系统首页",
        icon: 'icon-doc',
        routes: index,
        layout: false,
    },
    /**
     * 后台路由
     */
    {
        path: "/" + adminPath,
        name: "redirect." + adminPath,
        redirect: '/' + adminPath + '/' + (admin[0]?.path || 'index'),
        layout: false,
        hideInMenu: true,
    },
    {
        path: adminPath,
        name: adminPath,
        title: "后台管理",
        icon: 'icon-doc',
        routes: admin,
        layout: false,
        component: '@/layouts/admin',
    },
    /**
     * 用户路由
     */
    {
        path: "/" + userPath,
        name: "redirect." + userPath,
        redirect: '/' + userPath + '/' + (user[0]?.path || 'index'),
        layout: false,
        hideInMenu: true,
    },
    {
        path: userPath,
        name: userPath,
        title: "用户系统",
        icon: 'icon-doc',
        layout: false,
        component: '@/layouts/user',
        routes: user,
    },
    /**
     * 404
     */
    {
        path: '/*',
        component: '@/pages/404.tsx',
        name: "notFound",
        title: "页面不存在",
        layout: false,
        hideInMenu: true,
    }
];