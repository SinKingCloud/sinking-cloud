export default [
  {
    path: '/user',
    layout: false,
    routes: [
      {
        path: '/user',
        routes: [
          {
            name: '账户登录',
            path: '/user/login',
            component: './user/login',
          },
          {
            redirect: '/user/login',
          },
        ],
      },
      {
        component: './404',
      },
    ],
  },
  {
    name: "数据概览",
    icon: "HomeOutlined",
    path: '/index',
    routes: [
      {
        name: "控制台",
        path: 'index',
        component: './index/index',
      },
      {
        redirect: '/index/index',
      },
    ],
  },
  {
    name: "能力管理",
    icon: "AppstoreAddOutlined",
    path: '/open',
    routes: [
      {
        name: "云端挂机",
        path: 'uin',
        routes: [
          {
            name: "能力概览",
            path: 'home',
            component: './open/uin/home',
          },
          {
            name: "挂机列表",
            path: 'list',
            component: './open/uin/list',
          },
          {
            redirect: '/open/uin/home',
          },
        ],
      },
      {
        name: "全站登陆",
        path: 'oauth',
        routes: [
          {
            name: "接入管理",
            path: 'web',
            component: './open/oauth/web',
          },
          {
            redirect: '/open/oauth/web',
          },
        ],
      },
      {
        name: "开发文档",
        path: 'doc',
        component: './app/doc',
      },
      {
        redirect: '/index/index',
      },
    ],
  },
  {
    name: "应用管理",
    icon: "SettingOutlined",
    path: '/app',
    routes: [
      {
        name: "接入管理",
        path: 'open',
        component: './app/open',
      },
      {
        name: "授权管理",
        path: 'auth',
        component: './app/auth',
      },
      {
        redirect: '/app/open',
      },
    ],
  },
  {
    name: "财务管理",
    icon: "DollarCircleOutlined",
    path: '/pay',
    routes: [
      {
        name: "账户充值",
        path: 'recharge',
        component: './pay/recharge',
      },
      {
        name: "订单记录",
        path: 'order',
        component: './pay/order',
      },
      {
        name: "消费明细",
        path: 'log',
        component: './pay/log',
      },
      {
        redirect: '/pay/recharge',
      },
    ],
  },
  {
    name: "账户管理",
    icon: "UserOutlined",
    path: '/person',
    routes: [
      {
        name: "资料管理",
        path: 'setting',
        component: './person/setting',
      },
      {
        name: "实名认证",
        path: 'real',
        component: './person/real',
      },
      {
        name: "操作日志",
        path: 'log',
        component: './person/log',
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
