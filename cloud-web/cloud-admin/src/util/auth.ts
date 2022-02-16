/**
 * 设置登录token
 * @param type
 * @param token
 */
export function setLoginToken(type: string, token: string) {
  localStorage.setItem('sinking-token', token);
  localStorage.setItem('sinking-device', type);
}

/**
 * 获取登录token
 * @returns {string}
 */
export function getLoginToken() {
  return localStorage.getItem('sinking-token');
}

/**
 * 获取登录方式
 * @returns {string}
 */
export function getLoginType() {
  return localStorage.getItem('sinking-device');
}

/**
 * 返回请求header
 * @returns {{}}
 */
export function getHeaders() {
  return {
    'sinking-token': getLoginToken(),
    'sinking-device': getLoginType()
  };
}
