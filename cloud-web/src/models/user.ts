import {useState} from "react";
import {API} from "@/../typings";
import {getSystemUserInfo, getWebUserInfo} from "@/service/user";

export default () => {
    const [web, setWeb] = useState<API.WebUser>();//网站用户
    const [system, setSystem] = useState<API.SystemUser>();//后台用户

    /**
     * 获取网站用户信息
     */
    const getWebUser = async () => {
        const resp = await getWebUserInfo();
        return resp?.data || undefined;
    }

    /**
     * 刷新网站用户信息
     */
    const refreshWebUser = () => {
        getWebUser().then((d) => {
            setWeb(d);
        });
    }

    /**
     * 获取系统用户信息
     */
    const getSystemUser = async () => {
        const resp = await getSystemUserInfo();
        return resp?.data || undefined;
    }

    /**
     * 刷新系统用户信息
     */
    const refreshSystemUser = () => {
        getSystemUser().then((d) => {
            setSystem(d);
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
        refreshSystemUser
    };
};