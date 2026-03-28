import React from "react";
import {App} from "antd";

/**
 * 挂载静态方法
 * @param container 容器
 */
export function rootContainer(container: React.ReactElement) {
    return (
        <App>
            {container}
        </App>
    )
}