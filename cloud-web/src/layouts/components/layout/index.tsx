import React, {useEffect, useMemo, useState} from "react";
import {useLocation, useModel, Outlet, history, useSelectedRoutes} from "umi";
import {historyPush, getCurrentMenus, getFirstMenuWithoutChildren, getParentList} from "@/utils/route";
import {createStyles} from 'antd-style';
import Settings from "@/../config/defaultSettings";
import defaultSettings from "@/../config/defaultSettings";
import {Icon, Layout} from "sinking-antd";
import {App} from "antd";


/**
 * 首页系统
 */
const SKLayout: React.FC = () => {
    /**
     * 全局信息
     */
    const user = useModel("user");

    /**
     * 初始化数据
     */
    useEffect(() => {
        user?.refreshWebUser();
    }, []);

    return <App>
        <Outlet/>
    </App>;
}

export default SKLayout;