import React, {useEffect, useMemo, useState} from "react";
import {useLocation, useModel, Outlet, history, useSelectedRoutes} from "umi";
import {historyPush, getCurrentMenus, getFirstMenuWithoutChildren, getParentList} from "@/utils/route";
import {createStyles} from 'antd-style';
import Settings from "@/../config/defaultSettings";
import defaultSettings from "@/../config/defaultSettings";
import {Icon, Layout, useTheme} from "sinking-antd";
import {getPath} from "@/utils/assets";
import {deleteHeader} from "@/utils/auth";
import {logout} from "@/service/common/login";
import {App, Avatar, Col, Popover, Row, Tooltip} from "antd";

/**
 * 样式
 */
const useRightTopStyles = createStyles(({css, token, isDarkMode}: any): any => {
    return {
        img: {
            marginBottom: "5px",
        },
        nickname: {
            marginLeft: "3px",
            fontSize: "13px",
            color: isDarkMode ? token.colorTextSecondary : "rgb(150,150,150)",
            fontWeight: "bold",
        },
        bottomIconDark: {
            color: isDarkMode ? token.colorTextSecondary : "rgb(150,150,150)",
        },
        pop: css`
            margin-left: 10px;
            margin-right: 20px;
            display: initial;
            padding: 11px 5px;
            border-radius: 10px;
            transition: background-color 0.3s ease;
            cursor: pointer;

            .anticon {
                margin-left: 2px;
                font-size: 10px;
            }
        `,
        box: css`
            .ant-popover-container {
                padding: 0 !important;
                box-shadow: 0 0 5px 0 rgba(0, 0, 0, 0.15) !important;
                width: 210px;
            }

            .ant-popover-arrow:before {
                background-color: ${token?.colorPrimary} !important;
            }
        `,
        content_top: css`
            height: 70px;
            width: 100%;
            background-color: ${token?.colorPrimary};
            overflow: hidden;
            background-image: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAIgAAACGBAMAAAD0nt8RAAAAD1BMVEVHcEz///////////////8T4DEaAAAABXRSTlMADAYJA8T7L0gAAALASURBVGjezVpbcoMwDCS2DxBhDgC0BwhNDoDb3v9MfaQQMLb1cqfVVz6YjbS7ks2IptnFCOD7RhcDfIafVRgGvqPVJ/IZmoLcD4YqFbuAgAIkrCCKeqYaICsGnP8WxP0bkE05F71hK4EoWvC0YHh9EwN0v93Fr9frs3aefP/PjdA9Pfo3F8xuxRk74LlaTBpLaPOATaRA0A+lBPAB6rBUNy2KVXwmDOF8YwSs3g1Ij8zYVgPi0PYjlGOppJVATmi9BgcZUJDH1PKokzIamwdGVkGDPGEeB2S+jU9lT2/4KFASip4edxgtfpxDidJiIpvmOsoTYfQEI8UPmR3GOFJuOLHGe0o97YYT8fYKNE4jSnYP7mUpH2x292SW0vKtIyQlNAeM4pEzpTQ0E3BAXOpJA4mYqZTcCx+BCRKOOg5JDE+m5AskjVECcfGzJoNRsokFanRkSgpxJlNSiAudEgmIoYPM+AVWAzKQMTzlzUAOwqCkq0FJV8FqBa9NFWzCsBrt2BLbhEEJ1KDE16DE16Ckq0FJV4GSrNc4lGS9xqEk6zUOJVmvcSgB/UAq2GSqAWJq2IQxpQsgpsZcY4h89Jpbf0xSr5kngBeuaSOv3f/9xuQ2LUnP4tanFWlZqfiMID1nHnQZry/Kv/FB3PF9giTz9foyju/vc6Spl7TQHW5IqDaAPBphKhn/hBogbqoAwpv7WRNKue2kkzI/ZoSpzPIjtTWZ1+03BkgfzRPJKdQupmjlo98v1hoVp9AyiGbFKdRSN0Wy1x56C6FLJKsuhiYzYZtleL0iu2zQtoROXQwqM3nlMmmLQWRm7BkHfSJ5mXkL6aAuJiszd5FsBc1Lkpm/ATbqYpIyS1bRcSo35U5dseLfT0rpXt3qE6GtUTjcKj4SWFNRfcPh9Kvs9bLRN7qY1B+k/HCrTeSL24saozF4MR+gwScpYNNOxQAAAABJRU5ErkJggg==);
            background-repeat: no-repeat;
            background-position: right;
            border-top-left-radius: ${token?.borderRadius}px;
            border-top-right-radius: ${token?.borderRadius}px;
            line-height: 70px;
        `,
        ava: {
            height: "40px",
            width: "40px",
            marginLeft: "10px",
        },
        top_text: {
            color: "#fff",
            fontSize: "12px",
            letterSpacing: "1px",
            lineHeight: "100%",
            marginLeft: "8px",
            width: "auto"
        },
        userName: {
            fontSize: "13px",
            marginTop: "18px",
            marginBottom: "10px",
        },
        menuItemLabel: {
            display: "flex",
        },
        menuItemLeadIcon: {
            fontSize: "11px",
        },
        menu: {
            listStyle: "none",
            padding: 0,
            margin: 0,
            userSelect: "none",
            "li:last-of-type": {
                borderTop: "0.5px solid rgba(189, 189, 189, 0.2)",
                borderRadius: "0px 0px " + token.borderRadius + "px " + token.borderRadius + "px",
                height: "45px",
                lineHeight: "45px",
            }
        },
        menuItem: {
            cursor: "pointer",
            letterSpacing: "1px",
            height: "40px",
            lineHeight: "40px",
            fontSize: "12px",
            padding: "0px 15px",
            transition: "background-color 0.3s ease",
            color: isDarkMode ? token.colorTextSecondary : "rgba(0,0,0,0.65)",
            display: "flex",
            justifyContent: "space-between",
            ":hover": {
                backgroundColor: "rgba(0, 0, 0, 0.03)",
            },
            ".anticon": {
                fontSize: "11px",
            },
            "div>.anticon": {
                fontSize: "12.5px",
                marginRight: "7px"
            }
        },
        icon: {
            fontSize: "17px",
            padding: "7px",
            marginRight: "5px",
            cursor: "pointer",
            borderRadius: "5px",
            transition: "background-color 0.3s ease",
            ":hover": {
                backgroundColor: "rgba(0, 0, 0, 0.1)",
            },
            color: isDarkMode ? token.colorTextSecondary : "rgb(150,150,150)"

        },
    };
});

/**
 * 右侧部分组件
 * @constructor
 */
const RightTop: React.FC = () => {
    /**
     * 全局数据
     */
    const user = useModel("user");//用户信息
    const web = useModel("web");//站点信息
    const theme = useTheme();//主题信息
    const {message} = App.useApp();
    const {
        styles: {
            img,
            nickname,
            bottomIconDark,
            pop,
            content_top,
            ava,
            top_text,
            userName,
            box,
            menu,
            menuItem,
            icon,
            menuItemLabel,
            menuItemLeadIcon
        }
    } = useRightTopStyles();

    /**
     * 退出登录
     */
    const outLogin = async () => {
        message?.loading({content: '正在退出登录', duration: 600000, key: "outLogin"});
        await logout({
            onSuccess: (r) => {
                message?.success(r?.message || "退出登录成功")
                deleteHeader()
                user?.setWeb(undefined);
                message?.destroy("outLogin")
                historyPush("/");
            },
            onFail: (r) => {
                message?.error(r?.message || "退出登录失败")
            }
        })
    }

    return <>
        <Tooltip title={theme?.getModeName(theme?.mode as any)}>
            <Icon type={theme?.isDarkMode() ? 'icon-dark' : (theme?.isAutoMode() ? 'icon-auto' : 'icon-light')}
                  className={icon}
                  onClick={() => {
                      theme?.toggle?.();
                  }}/>
        </Tooltip>
        <Popover rootClassName={box} autoAdjustOverflow={false}
                 placement="bottomRight"
                 content={<>
                     <Row className={content_top}>
                         <Col span={6}><Avatar className={ava} src={user?.web?.avatar} alt={user?.web?.nick_name}>
                             {(user?.web?.nick_name || "未设置")}
                         </Avatar></Col>
                         <Col span={16} className={top_text}>
                             <div className={userName}>{user?.web?.nick_name || "未设置"}</div>
                             <div>{(user?.web?.phone || user?.web?.email || user?.web?.account)?.toString() || "未绑定"}</div>
                         </Col>
                     </Row>
                     <ul className={menu}>
                         <li className={menuItem} onClick={() => historyPush("/")}>
                             <div className={menuItemLabel}><Icon type="SettingOutlined" className={menuItemLeadIcon}/>账号管理
                             </div>
                             <Icon type="RightOutlined"></Icon>
                         </li>
                         <li className={menuItem} onClick={() => historyPush("/")}>
                             <div className={menuItemLabel}><Icon type="UnorderedListOutlined"
                                                                  className={menuItemLeadIcon}/>操作日志
                             </div>
                             <Icon type="RightOutlined"></Icon>
                         </li>
                         <li className={menuItem} onClick={outLogin}>
                             <div className={menuItemLabel}><Icon type="LogoutOutlined" className={menuItemLeadIcon}/>退出登录
                             </div>
                             <Icon type="RightOutlined"></Icon>
                         </li>
                     </ul>
                 </>}>
            <div className={pop}>
                <Avatar className={img} src={user?.web?.avatar} alt={user?.web?.nick_name}>
                    {(user?.web?.nick_name || "未设置")}
                </Avatar>
                <span className={nickname}>{user?.web?.nick_name || "未设置"}</span>
                <Icon className={web?.info?.ui?.user?.theme == "dark" ? bottomIconDark : ""} type="DownOutlined"/>
            </div>
        </Popover>
    </>
}

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
        menuHeader: {
            ".ant-layout-header": {
                height: "120px !important"
            }
        },
        collapsedLogo: {
            width: "39px",
        },
        unCollapsedLogo: {
            width: "160px",
        }
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
    const {styles: {menuColor, menuHeader, collapsedLogo, unCollapsedLogo}} = useSKLayoutStyles();

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
        user?.refreshWebUser((data) => {
            setLoading(false);
            if (!data) {
                historyPush("user.login");
            }
        });
    }, []);

    return (<>
        <Layout loading={loading || !user?.web}
                pathname={location?.pathname}
                matchedRoutes={match || []}
                onNavigate={(path) => history.push(path)}
                breadCrumbItems={breadCrumbItems}
                hideBreadCrumb={match?.at(-1)?.route?.hideBreadCrumb}
                waterMark={web?.info?.ui?.user?.watermark ? [user?.web?.nick_name, user?.web?.email] : ""}
                menus={getCurrentMenus(location?.pathname, true)}
                menuClassNames={menuColor}
                menuHeaderClassNames={menuHeader}
                layout={web?.info?.ui?.user?.layout != "left" ? "horizontal" : "inline"}
                flowLayout={web?.info?.ui?.user?.layout != "left"}
                menuTheme={web?.info?.ui?.user?.theme == "dark" ? "dark" : "light"}
                footer={<>©{new Date().getFullYear()} All Right
                    Revered {web?.info?.name || Settings?.title}</>}
                headerHidden={false}
                headerFixed={false}
                headerRight={<RightTop/>}
                menuCollapsedWidth={60}
                menuUnCollapsedWidth={210}
                menuBottomBtnIcon={web?.info?.contact ? <Icon type={"icon-customer"}/> : false}
                menuBottomBtnText={web?.info?.contact ? <>&nbsp;&nbsp;客户服务</> : false}
                collapsedLogo={() => {
                    return <img src={getPath("logo.svg")} className={collapsedLogo}
                                alt={defaultSettings?.title}/>;
                }}
                unCollapsedLogo={(isLight) => {
                    return <img
                        src={getPath((isLight ? "logo_light.png" : "logo_dark.png"))}
                        className={unCollapsedLogo}
                        alt={defaultSettings?.title}/>;
                }}>
            <Outlet/>
        </Layout>
    </>);
}


export default SKLayout;
