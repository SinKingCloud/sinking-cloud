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

//配置合并
Settings = Object.assign(Settings, {
  title: '沉沦云网络',
  logo: '/logo.svg',
  pwa: false,
  footer: '沉沦云科技开发',
  loadingTitle: '沉沦云平台',
  api: 'http://127.0.0.1:1000',
  base: '/',
});

export default Settings;
