import defaultSettings from "../../../config/defaultSettings";
import {getRandStr} from "@/util/string";

/**
 * 获取上传地址
 */
export function getCaptchaUrl(): any {
  let tempToken = getRandStr(16);
  return {
    url: defaultSettings.api + "/auth/verify/code?token=" + tempToken,
    token: tempToken,
  }
}
