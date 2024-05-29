import React from "react";
import Body from "@/layouts/components/body";
import Footer from "@/layouts/components/footer";
import Header from "@/layouts/components/header";
import Sider from "@/layouts/components/sider";
import 'dayjs/locale/zh-cn';
import {Theme} from "@/components";

import SkLayout, {LayoutProps} from "@/layouts/components/sinking";

const Layout: React.FC<LayoutProps> = (props) => {
    return (<Theme>
        <SkLayout {...props}/>
    </Theme>);
}

export {
    Body,
    Footer,
    Header,
    Sider,
}

export default Layout;