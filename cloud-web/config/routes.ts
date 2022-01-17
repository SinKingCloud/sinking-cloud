export default [
  {
    path: '/user',
    layout: false,
    routes: [
      {
        path: '/user',
        routes: [
          {
            name: 'login',
            path: '/user/login',
            component: './user/login',
          },
        ],
      },
      {
        component: './404',
      },
    ],
  },
  {
    path: '/admin',
    layout: false,
    routes: [
      {
        path: '/admin',
        routes: [
          {
            name: 'login',
            path: '/admin/login',
            component: './admin/login',
          },
        ],
      },
      {
        component: './404',
      },
    ],
  },
  {
    path: '/index',
    name: '用户首页',
    icon: 'smile',
    routes: [
      {
        path: 'index',
        name: '示例首页1',
        icon: 'smile',
        component: './user/index/welcome',
      },
      {
        path: 'index2',
        name: '示例首页2',
        icon: 'smile',
        component: './user/index/admin',
      },
      {
        component: './404',
      },
    ],
  },
  {
    path: '/test',
    name: '二级菜单',
    icon: 'crown',
    routes: [
      {
        path: 'index',
        name: '欢迎页面',
        icon: 'smile',
        component: './user/test',
      },
      {
        component: './404',
      },
    ],
  },
  {
    path: '/',
    redirect: '/index',
  },
  {
    component: './404',
  },
];
