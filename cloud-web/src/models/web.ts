import {useEffect, useState} from "react";
import {API} from "@/../typings";
import {getWebInfo} from "@/service/web";
import defaultSettings from "@/../config/defaultSettings";

export default () => {
    const [info, setInfo] = useState<API.WebInfo>({
        name: defaultSettings?.title,
    });
    /**
     * 获取站点信息
     */
    const getInfo = async () => {
        const resp = await getWebInfo();
        return resp?.data || {};
    }

    /**
     * 刷新站点信息
     */
    const refreshInfo = () => {
        getInfo().then((d) => {
            setInfo(d);
        });
    }

    /**
     * 初始化信息
     */
    useEffect(() => {
        refreshInfo();
    }, []);

    return {
        info,
        setInfo,
        getInfo,
        refreshInfo
    };
};