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
        component: './index/index/welcome',
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
            component: './index/index/welcome',
          },
          {
            name: "挂机列表",
            path: 'list',
            component: './index/index/welcome',
          },
          {
            redirect: '/open/uin/home',
          },
        ],
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
        name: "应用列表",
        path: 'list',
        component: './app/list',
      },
      {
        name: "开发文档",
        path: 'doc',
        component: './app/doc',
      },
      {
        redirect: '/app/list',
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
