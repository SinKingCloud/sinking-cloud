import {getWebInfo} from "@/service/web";
import defaultSettings from "../../config/defaultSettings";

/**
 * 初始化defaultSetting
 */
export async function initDefaultSetting() {
    const res = await getWebInfo();
    Object.keys(res?.data || {}).forEach((k) => {
        defaultSettings[k] = res?.data[k];
    });
    return defaultSettings;
}