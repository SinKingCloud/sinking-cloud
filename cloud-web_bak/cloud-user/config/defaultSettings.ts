let Settings = {
  layout: 'side',
  contentWidth: 'Fluid',
  splitMenus: false,
  fixedHeader: true,
  fixSiderbar: true,
  primaryColor: '#0051eb',
  navTheme: 'light',
  title: '沉沦云网络',
  logo: 'https://gw.alipayobjects.com/zos/rmsportal/KDpgvguMpGfqaHPjicRK.svg',
  pwa: false,
  iconfontUrl: '',
  api: '',
  ws: '',
  base: '/',
  publicPath: '/',
};
/**
 * 开发环境
 */
const isDev = process.env.NODE_ENV === 'development';
//配置合并
Settings = Object.assign(Settings, {
  title: '云开放平台',
  logo: '/logo.svg',
  api: isDev ? 'https://gateway.clwl.online' : 'https://gateway.clwl.online',
  ws: isDev ? 'ws://127.0.0.1:1000/ws/server' : 'wss://gateway.clwl.online/ws/server',
  base: '/admin/',
  publicPath: '/admin/',
});

export default Settings;
