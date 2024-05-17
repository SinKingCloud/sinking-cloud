import React, {useEffect} from "react";
import Layout from "@/layouts/components";
import {getUserMenuItems} from "@/utils/route";
import defaultSettings from "@/../config/defaultSettings";
import {Doc, FeedBack, Customer, Notice, Bottom, Search, Light, Dark} from "@/components/icon";
import {Icon} from "@/components";
import {useModel} from "umi";
import {useState} from "react";
import {deleteHeader, getLoginToken} from "@/utils/auth";
import {historyPush} from "@/utils/route";
import request from "@/utils/request";
import {App, Avatar, Button, Col, Input, Popover, Row, Tooltip} from "antd";
import {createStyles, useResponsive} from "antd-style";
import {outLogin} from "@/service/login";

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
                font-size: 12px;

                :hover {
                    border-color: #d9d9d9 !important;
                    width: 180px;
                }
            `,
         tab:{
             padding:"4px",
                 boxSizing:"border-box",
                 marginRight:"40px",
                 color: "rgb(122,122,122)",
                 border:"1px solid rgb(122,122,122)",
                 borderRadius:"5px",
                 cursor: "pointer",
         }
        };
    });
    const {styles: {search,tab}} = useStyles();
    const {mobile} = useResponsive();

    /**
     * 显示控制
     */
    const [open, setOpen] = useState(false);
    const [value, setValue] = useState("");
    const [menu, setMenu] = useState<any>(getUserMenuItems())
    const userRouter = getUserMenuItems()
    const handleChange = (e: any) => {
        const filterUserRouter = userRouter.filter((item: any) => {
            return item?.label?.includes(e?.target?.value)
        })
        setMenu(filterUserRouter)
        if(!e?.target?.value){
            setMenu(getUserMenuItems())
        }
        setValue(e?.target?.value);
        setOpen(e?.target?.value != "")
    }
    return <>
        {!mobile &&
            <>
                <Popover autoAdjustOverflow={false}
                         placement="bottomLeft"
                         content={<div style={{width: "400px", height: "300px"}} onClick={()=>{setOpen(false)}}>
                             <div style={{width: "100%"}}>
                                 <p style={{fontSize: 16, fontWeight: "bold", marginBottom: 8}}>最近访问</p>
                                 <div style={{width: "100%", display: "flex", flex: "wrap"}}>
                                     {menu.map((item: any) => (
                                         <Button className={tab} key={item?.key} onClick={() => {
                                             historyPush(item?.name)
                                         }}>{item?.label}</Button>
                                     ))}
                                 </div>
                             </div>
                         </div>} trigger="click" open={open} onOpenChange={setOpen}>
                    <Input className={search} value={value} prefix={<Icon type={Search}/>} placeholder="搜索菜单"
                           onChange={handleChange}/>
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
    const theme = useModel("theme");//主题信息
    const {message} = App.useApp()
    /**
     * 样式
     */
    const useStyles = createStyles(({css,token}): any => {
        return {
            icon: {
                fontSize: "17px",
                padding: "7px",
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
            },
            nickname: {
                marginLeft: "3px",
                fontSize: "13px",
                color: "rgb(122,122,122)",
                fontWeight: "bold",
            },
            profile: css`
                margin-left: 10px;
                margin-right: 20px;

                .anticon {
                    margin-left: 2px;
                    font-size: 10px;
                }
            `,
            pop: css`
                display: initial;
                padding: 11px 5px;
                border-radius: 10px;
                transition: background-color 0.5s ease;
                cursor: pointer;
            `,
            box:css`
                .ant-popover-inner{
                    box-shadow: 0 2px 8px 0 rgba(0,0,0,.15) !important;
                }
                .ant-popover-arrow:before {
                    background-color: rgba(85, 126, 253, 1) !important;
                }
            `,
            content_top: {
                height: "70px",
                width: "100%",
                backgroundColor: token?.colorPrimaryBorderHover,
                borderBottom:"1px solid gray",
                overflow: "hidden",
                backgroundImage: "url(https://www.xiaoke.cn/xk/static/bg_blue.a89f3e00.png)",
                backgroundRepeat: "no-repeat",
                backgroundSize: "100% 100%",
                backgroundPosition: "center",
                borderTopLeftRadius: "5px",
                borderTopRightRadius: "5px",
                lineHeight:"70px"
            },
            ava:{
                height: "40px",
                width: "40px",
                marginLeft: "20px",
                marginRight:"10px"
            },
            top_text:{
                color:"#fff",
                fontSize: "12px",
                letterSpacing: "1px",
                lineHeight:"100%",
                marginLeft: "12px",
            },
            text: {
                cursor: "pointer",
                borderRadius: "5px",
                letterSpacing: "1px",
                height: "40px",
                lineHeight:"40px",
                color:"#000",
                fontSize: "14px",
                paddingLeft: "10px",
                boxSizing: "border-box",
                transition: "background-color 0.5s ease",
                ":hover": {
                    backgroundColor: "rgba(0, 0, 0, 0.1)",
                }
            }
        };
    });
    const {styles: {icon, img, nickname, profile, pop,content_top,ava,text,top_text,box}} = useStyles();

    return <>
        <Tooltip title={theme?.isDarkTheme() ? "暗夜模式" : "日间模式"}>
            <Icon type={theme?.isDarkTheme() ? Dark : Light} className={icon} onClick={() => {
                theme?.toggle?.();
            }}/>
        </Tooltip>
        <Tooltip title="我要反馈">
            <Icon type={FeedBack} className={icon}/>
        </Tooltip>
        <Tooltip title="在线客服">
            <Icon type={Customer} className={icon}/>
        </Tooltip>
        <Tooltip title="站内消息">
            <Icon type={Notice} className={icon}/>
        </Tooltip>
        <Popover className={profile} overlayInnerStyle={{padding: 0}} overlayClassName={box} autoAdjustOverflow={false} placement="bottomRight"
                 content={<div style={{width: "220px"}}>
                     <Row className={content_top}>
                         <Col span={6}><Avatar className={ava} src={user?.web?.avatar}/></Col>
                         <Col span={8} className={top_text}>
                             <div style={{
                                 fontSize: 13,
                                 marginTop: 18,
                                 marginBottom: 10
                             }}>{user?.web?.nick_name || "未设置"}</div>
                             <div>测试内容...</div>
                         </Col>
                     </Row>
                     <ul style={{listStyle: "none", padding: 0, margin: 0,borderBottom:"1px solid gray"}}>
                         <li className={text}>账号设置</li>
                         <li className={text}>系统公告</li>
                         <li className={text}>联系我们</li>
                     </ul>
                     <ul style={{listStyle: "none", padding: 0, margin: 0}}>
                         <li className={text} onClick={ async() => {
                             message?.loading({content: '正在退出登录', duration: 600000, key: "outLogin"});
                             await outLogin({
                                 onSuccess: (res) => {
                                     if (res?.code == 200) {
                                         message?.success(res?.message)
                                         localStorage.removeItem("sinking-token")
                                         historyPush("user.login")
                                     }
                                 },
                                 onFail: (res) => {
                                     message?.error(res?.message)
                                 },
                                 onFinally:()=>{
                                     message?.destroy("outLogin")
                                 }
                             })
                         }}>退出登录</li>
                     </ul>
                 </div>}>
            <div className={pop}>
                <Avatar className={img} src={user?.web?.avatar}>
                    {(user?.web?.nick_name || "未设置")[0]}
                </Avatar>
                <span className={nickname}>{user?.web?.nick_name || "未设置"}</span>
                <Icon type={Bottom}/>
            </div>
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