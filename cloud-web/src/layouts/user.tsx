import React, {useEffect} from "react";
import Layout from "@/layouts/components";
import {getUserMenuItems} from "@/utils/route";
import defaultSettings from "@/../config/defaultSettings";
import {Doc, Icon, Information, KeFu, Notice} from "@/components/icon";
import {useModel, history} from "umi";
import {useState} from "react";
import {deleteHeader, getLoginToken} from "@/utils/auth";
import {historyPush} from "@/utils/route";
import request from "@/utils/request";
import {Avatar, Dropdown, message, Select, Space} from "antd";
import {outLogin} from "@/service/login";
import {DownOutlined} from "@ant-design/icons";
import {createStyles, useResponsive} from "antd-style";

/**
 * 中间件
 * @param ctx context
 * @param next 执行函数
 */
const check = async (ctx: any, next: any) => {
    await next();
    if (ctx.res.code == 503 || ctx.res.code == 403) {
        deleteHeader();
        historyPush("user.login");
    }
}
request.use(check);
/**
 * 样式
 */
const useStyles = createStyles((): any => {
    return {
        icon: {
            fontSize: "20px",
            marginRight: "15px",
            cursor: "pointer",
        },
        img: {
            marginBottom: "10px",
        },
        drop: {
            marginLeft: "8px",
            fontSize: "14px",
            color: "gray",
            marginRight: "20px",
        },
        list: {
            backgroundColor: "#fff"
        },
        listItem: {
            cursor: "pointer"
        }
    };
});
/**
 * 用户系统
 */
export default () => {

    /**
     * 全局信息
     */
    const {styles: {icon, img, drop}} = useStyles();
    const web = useModel("web");//网站信息
    const user = useModel("user");//用户信息
    const {mobile} = useResponsive()
    /**
     * 初始化用户信息
     */
    const [loading, setLoading] = useState(true);
    const initUser = () => {
        setLoading(true);
        if (getLoginToken() == "") {
            historyPush("user.login")
        }
        user?.getWebUser()?.then((u: any) => {
            user?.setWeb(u);
        })?.finally(() => {
            setLoading(false);
        });
    }

    const [menu, setMenu] = useState<any>()
    const initMenu = ()=>{
        const userMenu = getUserMenuItems();
        setMenu(userMenu)
    }
    /**
     * 搜索菜单栏
     */
    useEffect(() => {
        initUser();
        initMenu()
    }, []);

    return (
        <Layout loading={loading}
                menus={menu}
                footer={<>©{new Date().getFullYear()} All Right Revered {web?.info?.name || defaultSettings?.title}</>}
                headerLeft={<>
                    {!mobile &&
                        <Select
                            showSearch
                            style={{width: 200}}
                            placeholder="请选择菜单"
                            optionFilterProp="label"
                            filterOption={(input, option) => {
                                return option?.props?.label?.toUpperCase().indexOf(input.toUpperCase()) >= 0
                            }}
                            options={menu}
                        />
                    }
                </>}
                headerRight={
                    <>
                        <Icon type={Information} className={icon}/>
                        <Icon type={KeFu} className={icon}/>
                        <Icon type={Notice} className={icon}/>
                        <Avatar size="large" className={img} src={user?.web?.avatar}/>
                        <Dropdown className={drop} trigger={["click"]} overlayStyle={{width: "max-content"}}
                                  arrow={true}
                                  placement={"bottomRight"} menu={{
                            items: [{
                                key: '1',
                                label: (
                                    <a onClick={async () => {
                                        await outLogin({
                                            onSuccess: (r) => {
                                                if (r?.code === 200) {
                                                    message?.success(r?.message)
                                                    localStorage.removeItem("sinking-token")
                                                }
                                            },
                                            onFail: (r) => {
                                                message?.error(r?.message)
                                            }
                                        })
                                        history.push("/user/login")
                                    }}>退出登录</a>
                                )
                            }]
                        }}>
                            <a onClick={(e) => e.preventDefault()}>
                                <Space>
                                    {user?.web?.nick_name || "未设置"}
                                    <DownOutlined/>
                                </Space>
                            </a>
                        </Dropdown>
                    </>}
                menuBottomBtnIcon={Doc}
                menuBottomBtnText={"开发文档"}
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
                    return <img src={(defaultSettings?.basePath || "/") + "logo.svg"} width="39"
                                alt={defaultSettings?.title}/>;
                }}
                unCollapsedLogo={(isLight) => {
                    return <img
                        src={(defaultSettings?.basePath || "/") + (isLight ? "logo_light.png" : "logo_dark.png")}
                        width="130"
                        alt={defaultSettings?.title}/>;
                }}
                onMenuBtnClick={(state) => {
                    console.log("点击了菜单按钮", state)
                }}
        />
    );
}