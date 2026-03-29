import React from "react";
import {Outlet} from "umi";
import {App} from "antd";


/**
 * 首页系统
 */
const SKLayout: React.FC = () => {
    return <App>
        <Outlet/>
    </App>;
}

export default SKLayout;