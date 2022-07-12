/**
 * 获取随机字符
 * @param length 长度
 */
export function getRandStr(length = 16) {
  let str = 'abcdefghijklmnopqrstuvwxyz';
  str += str.toUpperCase();
  str += '0123456789'
  let _str = '';
  for (let i = 0; i < length; i++) {
    const rand = Math.floor(Math.random() * str.length);
    _str += str[rand];
  }
  return _str
}

/**
 * 对象分割成字符串
 * @param glue 分割字符
 * @param pieces 数据
 */
export function implode(glue: string, pieces: any) {
  let i = '',
    retVal = '',
    tGlue = '';
  if (arguments.length === 1) {
    pieces = glue;
    glue = '';
  }
  if (typeof pieces === 'object') {
    if (Object.prototype.toString.call(pieces) === '[object Array]') {
      return pieces.join(glue);
    }
    for (i in pieces) {
      retVal += tGlue + pieces[i];
      tGlue = glue;
    }
    return retVal;
  }
  return pieces;
}
