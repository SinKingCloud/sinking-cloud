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
        redirect: '/web/user',
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
