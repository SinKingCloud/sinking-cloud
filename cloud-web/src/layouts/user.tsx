import React, {useEffect} from "react";
import Layout from "@/layouts/components";
import {getUserMenuItems} from "@/utils/route";
import defaultSettings from "@/../config/defaultSettings";
import {Doc, Icon, FeedBack, Customer, Notice, Bottom, Search} from "@/components/icon";
import {useModel} from "umi";
import {useState} from "react";
import {deleteHeader, getLoginToken} from "@/utils/auth";
import {historyPush} from "@/utils/route";
import request from "@/utils/request";
import {Avatar, Input, Popover, Tooltip} from "antd";
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
 * 左侧部分组件
 * @constructor
 */
const LeftTop: React.FC = () => {
    /**
     * 样式
     */
    const useStyles = createStyles(({css}): any => {
        return {
            search: css`
                margin-left: 5px;
                height: 30px;
                line-height: 22px;
                width: 90px;
                transition: background-color, width 0.3s ease !important;
                border-radius: 15px;

                :hover {
                    border-color: #d9d9d9 !important;
                    width: 180px;
                }
            `
        };
    });
    const {styles: {search}} = useStyles();
    const {mobile} = useResponsive();

    /**
     * 显示控制
     */
    const [open, setOpen] = useState(false);
    const [value, setValue] = useState("");

    return <>
        {!mobile &&
            <>
                <Popover autoAdjustOverflow={false}
                         placement="bottomLeft"
                         content={<div style={{width: "400px", height: "300px"}}>
                             测试内容
                         </div>} trigger="focus" open={open} onOpenChange={setOpen}>
                    <Input className={search} value={value} prefix={<Icon type={Search}/>} placeholder="搜索菜单"
                           onChange={(e) => {
                               setValue(e?.target?.value);
                               setOpen(e?.target?.value != "")
                           }}/>
                </Popover>
            </>
        }
    </>
}

/**
 * 右侧部分组件
 * @constructor
 */
const RightTop: React.FC = () => {
    /**
     * 全局数据
     */
    const user = useModel("user");//用户信息

    /**
     * 样式
     */
    const useStyles = createStyles(({css}): any => {
        return {
            icon: {
                fontSize: "17px",
                padding: "5px",
                marginRight: "5px",
                cursor: "pointer",
                borderRadius: "5px",
                transition: "background-color 0.5s ease",
                ":hover": {
                    backgroundColor: "rgba(0, 0, 0, 0.1)",
                }
            },
            img: {
                marginBottom: "5px",
                cursor: "pointer",
            },
            nickname: {
                marginLeft: "3px",
                fontSize: "13px",
                color: "rgb(122,122,122)",
                fontWeight: "bold",
                cursor: "pointer",
            },
            profile: css`
                margin-left: 10px;
                margin-right: 20px;

                .anticon {
                    margin-left: 2px;
                    font-size: 10px;
                }
            `
        };
    });
    const {styles: {icon, img, nickname, profile}} = useStyles();

    return <>
        <Tooltip title="我要反馈">
            <Icon type={FeedBack} className={icon}/>
        </Tooltip>
        <Tooltip title="在线客服">
            <Icon type={Customer} className={icon}/>
        </Tooltip>
        <Tooltip title="站内消息">
            <Icon type={Notice} className={icon}/>
        </Tooltip>
        <Popover className={profile} autoAdjustOverflow={false} placement="bottomRight"
                 content={<div style={{width: "250px", height: "300px"}}>
                     测试内容
                 </div>}>
            <Avatar className={img} src={user?.web?.avatar}/>
            <span className={nickname}>{user?.web?.nick_name || "未设置"}</span>
            <Icon type={Bottom}/>
        </Popover>
    </>
}

/**
 * 用户系统
 */
export default () => {

    /**
     * 全局信息
     */
    const web = useModel("web");//网站信息
    const user = useModel("user");//用户信息

    /**
     * 初始化用户信息
     */
    const [loading, setLoading] = useState(true);
    const initUser = () => {
        setLoading(true);
        if (getLoginToken() == "") {
            historyPush("user.login");
            return;
        }
        user?.getWebUser()?.then((u: any) => {
            user?.setWeb(u);
        })?.finally(() => {
            setLoading(false);
        });
    }
    /**
     * 搜索菜单栏
     */
    useEffect(() => {
        initUser();
    }, []);

    return (
        <Layout loading={loading}
                menus={getUserMenuItems()}
                footer={<>©{new Date().getFullYear()} All Right Revered {web?.info?.name || defaultSettings?.title}</>}
                headerLeft={<LeftTop/>}
                headerRight={<RightTop/>}
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