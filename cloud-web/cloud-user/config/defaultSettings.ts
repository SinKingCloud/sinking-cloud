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
  footer: 'All Right Revered',
  loadingTitle: '沉沦云平台',
  api: '',
  base: '/',
};
/**
 * 开发环境
 */
const isDev = process.env.NODE_ENV === 'development';
//配置合并
Settings = Object.assign(Settings, {
  title: '云开放平台',
  logo: '/logo.svg',
  pwa: false,
  footer: '沉沦云科技(运营邮箱:lszcl@vip.qq.com)',
  loadingTitle: '沉沦云开放平台',
  api: isDev ? 'http://127.0.0.1:1000' : 'https://gateway.open-api.cn',
  base: '/',
});

export default Settings;
