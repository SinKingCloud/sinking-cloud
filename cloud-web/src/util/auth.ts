/**
 * 设置登录token
 * @param type
 * @param token
 */
export function setLoginToken(type: string, token: string) {
  localStorage.setItem('sinking-token', token);
  localStorage.setItem('sinking-type', type);
}

/**
 * 获取登录token
 * @returns {string}
 */
export function getLoginToken() {
  return localStorage.getItem('sinking-token');
}

/**
 * 获取登录token
 * @returns {string}
 */
export function getLoginType() {
  return localStorage.getItem('sinking-type');
}

/**
 * 返回请求header
 * @returns {{}}
 */
export function getHeaders() {
  return {
    'sinking-token': getLoginToken(),
    'sinking-type': getLoginType()
  };
}

/**
 * 判断是否后台
 */
export function isAdmin(): boolean {
  return localStorage.getItem('sinking-system') == "admin"
}

/**
 * 设置当前登陆系统
 * @param isAdmin 是否为后台
 */
export function setSystem(isAdmin: boolean) {
  localStorage.setItem('sinking-system', String(isAdmin));
}
