import {useState} from "react";
import {API} from "@/../typings";
import {getPersonInfo as getWebUserInfo} from "@/service/user/person";
import {getPersonInfo as getSystemUserInfo} from "@/service/admin/person";

export default () => {
    const [web, setWeb] = useState<API.WebUser>();//网站用户
    const [system, setSystem] = useState<API.SystemUser>();//网站用户
    /**
     * 获取网站用户信息
     */
    const getWebUser = async () => {
        const resp = await getWebUserInfo();
        if (resp?.code != 200) {
            return undefined;
        }
        return resp?.data || undefined;
    }

    /**
     * 刷新网站用户信息
     */
    const refreshWebUser = (callback: (d: any) => void) => {
        getWebUser().then((d) => {
            if (d) {
                setWeb(d);
            }
            callback?.(d);
        });
    }

    /**
     * 获取网站用户信息
     */
    const getSystemUser = async () => {
        const resp = await getSystemUserInfo();
        if (resp?.code != 200) {
            return undefined;
        }
        return resp?.data || undefined;
    }

    /**
     * 刷新网站用户信息
     */
    const refreshSystemUser = (callback: (d: any) => void) => {
        getSystemUser().then((d) => {
            if (d) {
                setSystem(d);
            }
            callback?.(d);
        });
    }

    return {
        web,
        setWeb,
        getWebUser,
        refreshWebUser,
        system,
        setSystem,
        getSystemUser,
        refreshSystemUser,
    };
};