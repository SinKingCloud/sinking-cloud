import React, {useEffect} from "react";
import Layout from "@/layouts/components";
import {getAllMenuDataWithChildren, getUserMenuItems} from "@/utils/route";
import defaultSettings from "@/../config/defaultSettings";
import {
    Doc,
    FeedBack,
    Customer,
    Notice,
    Bottom,
    Search,
    Light,
    Dark,
    Setting,
    Right,
    PayLog,
    Order, Out, Exit, Add
} from "@/components/icon";
import {Icon} from "@/components";
import {useModel} from "umi";
import {useState} from "react";
import {deleteHeader, getLoginToken} from "@/utils/auth";
import {historyPush} from "@/utils/route";
import request from "@/utils/request";
import {App, Avatar, Col, Divider, Empty, Input, Popover, Row, Tag, Tooltip} from "antd";
import {createStyles, useResponsive} from "antd-style";
import {outLogin} from "@/service/login";
import {pinyin} from "pinyin-pro";

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
    const useStyles = createStyles(({css, token}): any => {
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
            content: {
                width: "450px"
            },
            popover: css`
                user-select: none;

                .ant-popover-inner {
                    box-shadow: 0 2px 8px 0 rgba(0, 0, 0, .15) !important;
                    border-radius: 5px !important;
                }

                .ant-popover-content {
                    margin-left: -25px !important;
                }

                .ant-popover-inner-content {
                    margin-top: -20px !important;
                    padding: 10px 5px 0 5px;
                }
            `,
            panelSearch: {
                padding: "0"
            },
            panelList: {
                padding: "5px 0",
            },
            line: {
                height: "40px",
                lineHeight: "40px",
            },
            title: css`
                font-size: 13px;
                font-weight: bolder;
                float: left;

                div {
                    width: 5px;
                    background-color: ${token?.colorPrimary};
                    border-radius: 0 5px 5px 0;
                    height: 20px;
                    float: left;
                    margin-top: 14.7%;
                    margin-right: 7px;
                }
            `,
            divider: {
                margin: "0"
            },
            btn: {
                margin: "5px",
                fontSize: "10px",
                color: "gray",
                border: "1px solid gray",
                padding: "0 5px",
                transition: "background-color 0.3s ease",
                cursor: "pointer",
                ":hover": {
                    backgroundColor: "rgba(0, 0, 0, 0.1)",
                }
            },
            empty: {
                margin: "0 0 10px 0"
            },
            icon: {
                marginTop: "5px",
                float: "right",
                fontSize: "14px",
                padding: "6px",
                borderRadius: "5px",
                cursor: "pointer",
                transition: "background-color 0.3s ease",
                ":hover": {
                    backgroundColor: "rgba(0, 0, 0, 0.1)",
                }
            },

            menu: css`
                max-height: 300px;
                overflow-y: auto;


                h4 {
                    color: rgb(119, 119, 119);
                    margin-bottom: 0;
                    margin-top: 5px;
                    margin-left: 5px;
                }

                ul {
                    list-style: none;
                    padding: 0;
                    margin: 0 0 10px 0;

                    .ant-empty {
                        margin: 10px !important;
                    }

                    li {
                        padding: 3px 0 3px 10px;
                        cursor: pointer;
                        font-size: 11px;
                        color: rgb(119, 119, 119);
                        margin: 3px 10px 3px 0;
                        transition: background-color 0.5s ease;
                        border-radius: 5px;
                        display: flex;
                        justify-content: space-between;

                        :hover {
                            background-color: rgba(0, 0, 0, 0.1);
                        }

                        .anticon {
                            margin-right: 10px;
                            font-size: 8px;
                        }
                    }
                }

                ::-webkit-scrollbar {
                    width: 5px !important; /* 纵向滚动条 */
                    background-color: transparent;
                }

                ::-webkit-scrollbar-track {
                    background-color: transparent;
                    box-shadow: inset 0 0 6px rgba(0, 0, 0, 0);
                }

                ::-webkit-scrollbar-thumb {
                    background-color: rgba(210, 210, 210, 0.11);
                    border-radius: 10px;
                    box-shadow: inset 0 0 6px rgba(0, 0, 0, 0);
                }
            `,
        }
    });
    const {
        styles: {
            content,
            panelSearch,
            panelList,
            search,
            popover,
            line,
            title,
            divider,
            btn,
            icon,
            empty,
            menu
        }
    } = useStyles();
    const {mobile} = useResponsive();
    /**
     * 显示控制
     */
    const [open, setOpen] = useState(false);
    const [value, setValue] = useState("");
    const handleChange = (e: any) => {
        setValue(e?.target?.value);
    }
    /**
     * 快捷功能操作
     */
    const [edit, setEdit] = useState(false);
    const editToggle = () => {
        setEdit(!edit);
    }
    const [list, setList] = useState<any>([{name: "test"}]);
    /**
     * 系统功能操作
     */
    const [menus, setMenus] = useState<any>(getAllMenuDataWithChildren({
        "label": "系统首页",
        children: getUserMenuItems()
    }));
    const Menu: React.FC<{ num?: any, index?: any }> = (props) => {
        return <>{menus.map((x, d) => {
            if (d % (props?.num || 0) == (props?.index || 0)) {
                let liNum = 0
                return <div key={"menu-ul-" + d}>
                    <h4>{x?.label}</h4>
                    <ul>{x?.children?.map((a, b) => {
                        if (value && a?.label && !a?.label?.includes(value) && !pinyin(a?.label, {
                            toneType: 'none',
                            separator: ''
                        })?.includes(value) && !pinyin(a?.label, {
                            pattern: 'first',
                            toneType: 'none',
                            separator: ''
                        })?.includes(value)) {
                            return;
                        }
                        liNum++
                        return <Tooltip placement={"right"} key={"menu-tooltip-" + d + "-" + b}
                                        title={edit ? "添加到快捷功能" : ""}>
                            <li key={"menu-li-" + d + "-" + b}>{a?.label}<Icon type={edit ? Add : Right}/></li>
                        </Tooltip>
                    })}
                        {liNum <= 0 && <Empty key={"menu-empty-" + d}
                                              image={Empty.PRESENTED_IMAGE_SIMPLE}/>}
                    </ul>
                </div>
            }
        })}</>;
    }


    return <>
        {!mobile &&
            <Popover autoAdjustOverflow={false}
                     placement="bottomLeft"
                     trigger="click"
                     open={open}
                     onOpenChange={setOpen}
                     overlayClassName={popover}
                     content={<div className={content}>
                         <div className={panelList}>
                             <div className={line}>
                                 <div className={title}>
                                     <div/>
                                     功能列表
                                 </div>
                                 <Tooltip title={edit ? "保存功能并退出" : "编辑或删除功能"}>
                                     <Icon type={edit ? Exit : Setting} className={icon} onClick={editToggle}/>
                                 </Tooltip>
                             </div>
                             <div className={menu}>
                                 <Row align="top">
                                     <Col span={8}>
                                         <Menu num={3} index={0}/>
                                     </Col>
                                     <Col span={8}>
                                         <Menu num={3} index={1}/>
                                     </Col>
                                     <Col span={8}>
                                         <Menu num={3} index={2}/>
                                     </Col>
                                 </Row>
                             </div>
                             <Divider className={divider}/>
                         </div>
                         <div className={panelSearch}>
                             <div className={line}>
                                 <div className={title}>
                                     <div/>
                                     快捷功能
                                 </div>
                             </div>
                             {((list?.length || 0) <= 0 &&
                                 <Empty className={empty}
                                        image={Empty.PRESENTED_IMAGE_SIMPLE}/>) || new Array(10).fill(null).map((_, index) => (
                                 <Tag size={"small"} className={btn} key={index} closable={edit}>系统首页</Tag>
                             ))}
                         </div>
                     </div>}>
                <Input className={search} value={value} prefix={<Icon type={Search}/>} placeholder="搜索菜单"
                       onChange={handleChange}/>
            </Popover>
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
    const {message} = App.useApp();
    /**
     * 样式
     */
    const useStyles = createStyles(({css, token}): any => {
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
            content: {
                width: "250px",
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
            popover: css`
                .ant-popover-inner {
                    padding: 0;
                    box-shadow: 0 2px 8px 0 rgba(0, 0, 0, .15) !important;
                    border-radius: 5px !important;
                }

                .ant-popover-arrow:before {
                    background-color: ${token?.colorPrimaryBorderHover} !important;
                }
            `,
            popBg: {
                height: "70px",
                backgroundColor: token?.colorPrimaryBorderHover,
                overflow: "hidden",
                backgroundImage: "url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAIgAAACGBAMAAAD0nt8RAAAAD1BMVEVHcEz///////////////8T4DEaAAAABXRSTlMADAYJA8T7L0gAAALASURBVGjezVpbcoMwDCS2DxBhDgC0BwhNDoDb3v9MfaQQMLb1cqfVVz6YjbS7ks2IptnFCOD7RhcDfIafVRgGvqPVJ/IZmoLcD4YqFbuAgAIkrCCKeqYaICsGnP8WxP0bkE05F71hK4EoWvC0YHh9EwN0v93Fr9frs3aefP/PjdA9Pfo3F8xuxRk74LlaTBpLaPOATaRA0A+lBPAB6rBUNy2KVXwmDOF8YwSs3g1Ij8zYVgPi0PYjlGOppJVATmi9BgcZUJDH1PKokzIamwdGVkGDPGEeB2S+jU9lT2/4KFASip4edxgtfpxDidJiIpvmOsoTYfQEI8UPmR3GOFJuOLHGe0o97YYT8fYKNE4jSnYP7mUpH2x292SW0vKtIyQlNAeM4pEzpTQ0E3BAXOpJA4mYqZTcCx+BCRKOOg5JDE+m5AskjVECcfGzJoNRsokFanRkSgpxJlNSiAudEgmIoYPM+AVWAzKQMTzlzUAOwqCkq0FJV8FqBa9NFWzCsBrt2BLbhEEJ1KDE16DE16Ckq0FJV4GSrNc4lGS9xqEk6zUOJVmvcSgB/UAq2GSqAWJq2IQxpQsgpsZcY4h89Jpbf0xSr5kngBeuaSOv3f/9xuQ2LUnP4tanFWlZqfiMID1nHnQZry/Kv/FB3PF9giTz9foyju/vc6Spl7TQHW5IqDaAPBphKhn/hBogbqoAwpv7WRNKue2kkzI/ZoSpzPIjtTWZ1+03BkgfzRPJKdQupmjlo98v1hoVp9AyiGbFKdRSN0Wy1x56C6FLJKsuhiYzYZtleL0iu2zQtoROXQwqM3nlMmmLQWRm7BkHfSJ5mXkL6aAuJiszd5FsBc1Lkpm/ATbqYpIyS1bRcSo35U5dseLfT0rpXt3qE6GtUTjcKj4SWFNRfcPh9Kvs9bLRN7qY1B+k/HCrTeSL24saozF4MR+gwScpYNNOxQAAAABJRU5ErkJggg==)",
                backgroundRepeat: "no-repeat",
                backgroundPosition: "right",
                borderTopLeftRadius: "5px",
                borderTopRightRadius: "5px",
                lineHeight: "70px"
            },
            avatar: {
                height: "40px",
                width: "40px",
                marginLeft: "20px",
                marginRight: "10px"
            },
            popBgNickName: {
                color: "white",
                fontSize: "12px",
                letterSpacing: "1px",
                lineHeight: "100%",
                marginLeft: "12px",
            },
            popBgNickNameTop: {
                fontSize: 15,
                fontWeight: "bolder",
                marginTop: 18,
                marginBottom: 10
            },
            regTime: {
                fontSize: "11px",
            },
            menu: {
                listStyle: "none",
                padding: 0,
                margin: 0,
                userSelect: "none",
                "li:last-of-type": {
                    borderTop: "0.5px solid rgba(189, 189, 189, 0.2)",
                    borderRadius: "0px 0px 5px 5px",
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
                padding: "0px 20px",
                transition: "background-color 0.3s ease",
                color: "rgb(115, 115, 115)",
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
        };
    });
    const {
        styles: {
            icon,
            img,
            nickname,
            content,
            profile,
            pop,
            popBg,
            avatar,
            popBgNickName,
            popBgNickNameTop,
            regTime,
            popover,
            menu,
            menuItem
        }
    } = useStyles();

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
        <Popover className={profile} overlayClassName={popover}
                 autoAdjustOverflow={false} placement="bottomRight"
                 content={<div className={content}>
                     <Row className={popBg}>
                         <Col span={6}><Avatar className={avatar} src={user?.web?.avatar}/></Col>
                         <Col span={14} className={popBgNickName}>
                             <div className={popBgNickNameTop}>{user?.web?.nick_name || "未设置"}</div>
                             <div className={regTime}>+86 {user?.web?.phone || "未设置"}</div>
                         </Col>
                     </Row>
                     <ul className={menu}>
                         <li className={menuItem}>
                             <div><Icon type={Setting}/>账号设置</div>
                             <Icon type={Right}/>
                         </li>
                         <li className={menuItem}>
                             <div><Icon type={PayLog}/>财务管理</div>
                             <Icon type={Right}/>
                         </li>
                         <li className={menuItem}>
                             <div><Icon type={Order}/>我的订单</div>
                             <Icon type={Right}/>
                         </li>
                         <li className={menuItem} onClick={async () => {
                             message?.loading({content: '正在退出登录', duration: 600000, key: "outLogin"});
                             await outLogin({
                                 onSuccess: (res) => {
                                     if (res?.code == 200) {
                                         message?.success(res?.message);
                                         deleteHeader();
                                         historyPush("user.login");
                                     }
                                 },
                                 onFail: (res) => {
                                     message?.error(res?.message);
                                 },
                                 onFinally: () => {
                                     message?.destroy("outLogin");
                                 }
                             })
                         }}>
                             <div><Icon type={Out}/>退出登录</div>
                             <Icon type={Right}/>
                         </li>
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