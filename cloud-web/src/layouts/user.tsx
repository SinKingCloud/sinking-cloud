import Layout from "@/layouts/components";
import {getUserMenuItems} from "@/utils/route";
import defaultSettings from "@/../config/defaultSettings";
import {Doc} from "@/components/icon";

export default () => {
    return (
        <Layout menus={getUserMenuItems()}
                footer={<>©{new Date().getFullYear()} All Right Revered {defaultSettings?.title}</>}
                headerLeft={<>左侧内容</>}
                headerRight={<>右侧内容</>}
                menuBottomBtnIcon={Doc}
                menuBottomBtnText={"测试按钮"}
                onMenuBottomBtnClick={() => {
                    console.log("点击了底部按钮")
                }}
                menuCollapsedWidth={60}
                menuUnCollapsedWidth={210}
                onLogoClick={() => {
                    console.log("点击了logo")
                }}
                onMenuClick={(item) => {
                    console.log("点击了菜单", item)
                }}
                collapsedLogo={() => {
                    return <img src="/logo.svg" width="39"
                                alt={defaultSettings?.title}/>;
                }}
                unCollapsedLogo={(isLight) => {
                    return <img src={isLight ? "/logo_light.png" : "/logo_dark.png"} width="130"
                                alt={defaultSettings?.title}/>;
                }}
                onMenuBtnClick={(state) => {
                    console.log("点击了菜单按钮", state)
                }}
        />
    );
}