import React from "react";
import {Body} from "sinking-antd";
import Info from "./components/info";
import Link from "./components/link";
import Overview from "./components/overview";

export default (): React.ReactNode => {
    return (
        <Body>
            <Info/>
            <Link/>
            <Overview/>
        </Body>
    );
};
