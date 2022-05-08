export default [
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
    routes: [
      {
        path: 'index',
        component: './index/index/welcome',
      },
      {
        redirect: '/index/index',
      },
    ],
  },
  {
    path: '/system',
    routes: [
      {
        path: 'access',
        component: './system/access',
      },
      {
        path: 'role',
        component: './system/role',
      },
      {
        path: 'menu',
        component: './system/menu',
      },
      {
        path: 'config',
        component: './system/config',
      },
      {
        path: 'user',
        component: './system/user',
      },
      {
        path: 'log',
        component: './system/log',
      },
      {
        path: 'log_request',
        component: './system/log_request',
      },
      {
        path: 'file',
        component: './system/file',
      },
      {
        redirect: '/system/user',
      },
    ],
  },
  {
    path: '/web',
    routes: [
      {
        path: 'user',
        component: './web/user',
      },
      {
        path: 'order',
        component: './web/user_order',
      },
      {
        path: 'notice',
        component: './web/notice',
      },
      {
        redirect: '/web/user',
      },
    ],
  },
  {
    name: "账号管理",
    path: '/person',
    routes: [
      {
        name: "账号资料",
        path: 'setting',
        component: './person/setting',
      },
      {
        redirect: '/person/setting',
      },
    ],
  },
  {
    path: '/',
    redirect: '/index/index',
  },
  {
    component: './404',
  },
];
