import React, {useEffect, useMemo, useState} from "react";
import {useLocation, useModel, Outlet, history, useSelectedRoutes} from "umi";
import {historyPush, getCurrentMenus, getFirstMenuWithoutChildren, getParentList} from "@/utils/route";
import {createStyles} from 'antd-style';
import Settings from "@/../config/defaultSettings";
import defaultSettings from "@/../config/defaultSettings";
import {Icon, Layout} from "sinking-antd";


/**
 * 样式信息
 */
const useSKLayoutStyles = createStyles(({token}): any => {
    return {
        menuColor: {
            ">.ant-menu-item>.ant-menu-item-icon,>.ant-menu-submenu>.ant-menu-submenu-title>.ant-menu-item-icon": {
                color: token?.colorPrimaryHover + " !important",
            },
        },
        collapsedImg: {
            width: "35px",
        },
        collapsedLogo: {
            fontSize: "27px",
        },
        unCollapsed: {
            overflow: "hidden",
            position: "absolute",
            display: "inline-flex",
            ">span": {
                fontSize: "27px",
            },
            ">img": {
                width: "35px",
                float: "left"
            },
            ">div": {
                fontSize: "25px",
                marginLeft: "5px",
                fontWeight: "bolder",
                float: "left",
                lineHeight: "35px",
                whiteSpace: "nowrap",
            }
        },
    };
});

/**
 * 用户系统
 */
const SKLayout: React.FC = () => {
    /**
     * 全局信息
     */
    const user = useModel("user");
    const web = useModel("web");
    const location = useLocation();
    const match = useSelectedRoutes();
    const [loading, setLoading] = useState(true);
    const {styles: {menuColor, collapsedImg, unCollapsed, collapsedLogo}} = useSKLayoutStyles();

    /**
     * 计算面包屑数据
     */
    const breadCrumbItems = useMemo(() => {
        if (!location?.pathname) return [];
        const items = getParentList(getCurrentMenus(location?.pathname, false), match?.at(-1)?.route?.name);
        const temp = [{
            title: "首页",
            onClick: () => {
                historyPush(getFirstMenuWithoutChildren(getCurrentMenus(location?.pathname))?.name || "");
            },
        }];
        const handleItemClick = (x: any) => {
            if (x?.children && x?.children?.length > 0) {
                historyPush(getFirstMenuWithoutChildren(x?.children)?.name || "");
            } else {
                historyPush(x?.name);
            }
        };
        items.forEach((x) => {
            temp.push({
                title: x?.label,
                onClick: () => handleItemClick(x),
            });
        });
        return temp;
    }, [location?.pathname, match]);

    /**
     * 初始化数据
     */
    useEffect(() => {
        setLoading(true)
        user?.refreshWebUser(() => {
            setLoading(false)
        });
    }, []);

    return (
        <Layout loading={loading || !user?.web}
                pathname={location?.pathname}
                matchedRoutes={match || []}
                onNavigate={(path) => history.push(path)}
                breadCrumbItems={breadCrumbItems}
                hideBreadCrumb={match?.at(-1)?.route?.hideBreadCrumb}
                waterMark={web?.info?.ui?.water_mark ? [user?.web?.nick_name, user?.web?.email] : ""}
                menus={getCurrentMenus(location?.pathname, true)}
                menuClassNames={menuColor}
                layout={web?.info?.ui?.layout != "left" ? "horizontal" : "inline"}
                flowLayout={web?.info?.ui?.layout != "left"}
                menuTheme={web?.info?.ui?.theme == "dark" ? "dark" : "light"}
                footer={<>©{new Date().getFullYear()} All Right
                    Revered {web?.info?.name || Settings?.title}</>}
                headerHidden={false}
                headerFixed={false}
                menuCollapsedWidth={60}
                menuUnCollapsedWidth={210}
                collapsedLogo={() => {
                    return web?.info?.logo ? <img
                            src={web?.info?.logo || (Settings?.basePath || "/") + "logo.svg"}
                            alt={web?.info?.name} className={collapsedImg}/> :
                        <Icon type="icon-logo" style={{color: web?.info?.ui?.color || defaultSettings?.color}}
                              className={collapsedLogo}/>
                }}
                unCollapsedLogo={() => {
                    return (
                        <div className={unCollapsed}>
                            {web?.info?.logo ? <img
                                    src={web?.info?.logo || (Settings?.basePath || "/") + "logo.svg"}
                                    alt={web?.info?.name || Settings?.title}/> :
                                <Icon type="icon-logo"
                                      style={{color: web?.info?.ui?.color || defaultSettings?.color}}/>}
                            <div
                                style={{color: web?.info?.ui?.color ? web?.info?.ui?.color : defaultSettings?.color}}>
                                {web?.info?.name || Settings?.title}
                            </div>
                        </div>)
                }}>
            <Outlet/>
        </Layout>
    );
}


export default SKLayout;