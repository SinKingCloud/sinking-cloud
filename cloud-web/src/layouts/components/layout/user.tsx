import React, {useCallback, useEffect, useMemo, useRef, useState} from "react";
import {history, Outlet, useLocation, useModel, useSelectedRoutes} from "umi";
import {getCurrentMenus, getFirstMenuWithoutChildren, getParentList, historyPush} from "@/utils/route";
import {createStyles} from 'antd-style';
import Settings from "@/../config/defaultSettings";
import defaultSettings from "@/../config/defaultSettings";
import {Icon, Layout, ProModal, type ProModalRef, type SinKingRef, Title, useTheme} from "sinking-antd";
import {getPath} from "@/utils/assets";
import {deleteHeader} from "@/utils/auth";
import {logout} from "@/service/common/login";
import {getMessageInfo, getMessageList} from "@/service/user/message";
import useSocket, {MESSAGE_REFRESH} from "@/utils/socket";
import {App, Avatar, Button, Col, Empty, Grid, Popover, Row, Spin, Tooltip, Typography} from "antd";
import {ago} from "@/utils/time";

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
        messageTrigger: css`
            position: relative;
            color: ${token.colorPrimaryHover} !important;

            &:hover {
                color: ${token.colorPrimaryHover} !important;
            }

            &[data-count]:after {
                content: attr(data-count);
                position: absolute;
                top: 4px;
                right: 3px;
                min-width: 13px;
                height: 13px;
                padding: 0 3px;
                border: 1px solid ${isDarkMode ? token.colorBgContainer : "#fff"};
                border-radius: 999px;
                background: ${token.colorError};
                color: #fff;
                font-size: 9px;
                font-weight: 500;
                line-height: 11px;
                text-align: center;
                transform: translate(42%, -42%);
                pointer-events: none;
                box-shadow: 0 1px 3px rgba(0, 0, 0, 0.12);
            }
        `,
        messageBox: css`
            &.ant-popover {
                max-width: calc(100vw - 24px);
            }

            .ant-popover-content,
            .ant-popover-inner,
            .ant-popover-container,
            .ant-spin-nested-loading,
            .ant-spin-container {
                max-width: 100%;
            }

            .ant-popover-container {
                width: 310px;
                max-width: calc(100vw - 24px);
                overflow: hidden;
                border-radius: ${token.borderRadiusLG}px;
                box-shadow: ${token.boxShadowSecondary};
            }

            .ant-popover-inner,
            .ant-popover-container {
                padding: 0 !important;
            }

            .ant-popover-arrow:before {
                background: ${isDarkMode ? token.colorBgElevated : token.colorBgContainer} !important;
            }

            @media (max-width: 575px) {
                &.ant-popover {
                    left: 12px !important;
                    right: 12px !important;
                    width: auto !important;
                    max-width: none;
                }

                .ant-popover-content,
                .ant-popover-inner,
                .ant-popover-container,
                .ant-spin-nested-loading,
                .ant-spin-container {
                    width: 100%;
                }

                .ant-popover-container {
                    max-width: none;
                }

                .ant-popover-arrow {
                    display: none;
                }
            }
        `,
        messageList: css`
            background: ${isDarkMode ? token.colorBgElevated : token.colorBgContainer};

            .message-head {
                display: flex;
                align-items: center;
                justify-content: space-between;
                gap: 12px;
                padding: 13px 14px 11px;
                border-bottom: 1px solid ${isDarkMode ? "rgba(255,255,255,0.07)" : token.colorBorderSecondary};
            }

            .message-head-title {
                color: ${token.colorTextHeading};
                font-size: 14px;
                font-weight: 700;
                line-height: 22px;
            }

            .message-head-link {
                color: ${token.colorPrimary};
                font-size: 12px;
                line-height: 18px;
                white-space: nowrap;
                cursor: pointer;
            }

            .message-head-link:hover {
                color: ${token.colorPrimaryHover};
            }

            .message-body {
                max-height: 348px;
                overflow: auto;
            }

            .ant-empty {
                margin: 34px 0;
            }

        `,
        messageItem: css`
            padding: 10px 14px;
            border-radius: 0;
            cursor: pointer;
            position: relative;
            transition: background-color 0.2s ease;

            &:hover {
                background: ${isDarkMode ? "rgba(255,255,255,0.035)" : token.colorFillQuaternary};
            }

            .item-title {
                display: flex;
                align-items: center;
                gap: 8px;
                min-width: 0;
                color: ${token.colorText};
                font-size: 13px;
                font-weight: 400;
                line-height: 20px;
            }

            .item-title span:last-child {
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
                flex: 1;
            }

            .unread-dot {
                width: 5px;
                height: 5px;
                border-radius: 50%;
                flex: 0 0 auto;
                background: ${token.colorPrimary};
            }

            .item-time {
                margin-top: 5px;
                padding-left: 13px;
                font-size: 12px;
                line-height: 17px;
                color: ${token.colorTextTertiary};
            }
        `,
        messageDetail: css`
            padding: 0 4px;

            .detail-header {
                text-align: center;
                padding: 16px 0 12px;
                margin-bottom: 20px;
                border-bottom: 1px solid ${isDarkMode ? "rgba(255,255,255,0.08)" : "rgba(0,0,0,0.06)"};
            }

            .detail-header h2 {
                margin: 0;
                color: ${token.colorTextHeading};
                font-size: 20px;
                font-weight: 600;
                line-height: 1.4;
            }

            .detail-header .meta {
                display: flex;
                align-items: center;
                justify-content: center;
                flex-wrap: wrap;
                gap: 8px 18px;
                margin-top: 10px;
                color: ${token.colorTextTertiary};
                font-size: 12px;
            }

            .detail-content {
                padding: 0 4px 12px;
                color: ${token.colorText};
                font-size: 14px;
                line-height: 1.8;
                white-space: pre-wrap;
                word-break: break-word;
            }

            .detail-content p {
                margin: 0 0 12px;
            }

            .detail-content p:last-child {
                margin-bottom: 0;
            }

            .detail-content img {
                display: block;
                max-width: 100%;
                margin: 12px 0;
                border-radius: ${token.borderRadius}px;
            }

            .detail-content a {
                color: ${token.colorPrimary};
            }
        `,
        messageDetailSpin: css`
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 40px 0;
        `,
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
    const screens = Grid.useBreakpoint();
    const isMobile = Boolean((screens.xs || screens.sm) && !screens.md);
    const {message, modal} = App.useApp();
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
            messageTrigger,
            messageBox,
            messageList,
            messageItem,
            messageDetail,
            messageDetailSpin,
            menuItemLabel,
            menuItemLeadIcon
        }
    } = useRightTopStyles();

    /**
     * 消息部分
     */
    const detailModalRef = useRef<ProModalRef>({} as ProModalRef);
    const [notice, setNotice] = useState<{ open: boolean; loading: boolean; list: any[]; total: number }>({
        open: false,
        loading: false,
        list: [],
        total: 0,
    });
    const [detailLoading, setDetailLoading] = useState(false);
    const [detail, setDetail] = useState<any>();
    const noticeCount = notice.total > 99 ? "99+" : String(notice.total);
    const refreshNotice = useCallback(async (showLoading = false) => {
        if (showLoading) {
            setNotice((data) => ({...data, loading: true}));
        }
        try {
            const resp = await getMessageList({
                body: {page: 1, page_size: 5, status: "0"},
            });
            if (resp?.code == 200) {
                setNotice((data) => ({
                    ...data,
                    list: resp?.data?.list || [],
                    total: resp?.data?.total || 0,
                }));
            }
        } finally {
            if (showLoading) {
                setNotice((data) => ({...data, loading: false}));
            }
        }
    }, []);
    const closeNoticePanel = () => setNotice((data) => ({...data, open: false}));
    const openNoticePage = () => {
        closeNoticePanel();
        historyPush("user.person.message");
    };
    const openNotice = async (item: any) => {
        closeNoticePanel();
        setDetail(item);
        detailModalRef.current?.show();
        setDetailLoading(true);
        const resp = await getMessageInfo({body: {id: item.id}});
        setDetailLoading(false);
        if (resp?.code != 200) {
            message.error(resp?.message || "获取消息失败");
            detailModalRef.current?.hide();
            return;
        }
        const info = resp?.data || item;
        setDetail(info);
        setNotice((data) => ({
            ...data,
            total: Math.max(0, data.total - 1),
            list: data.list.filter((x) => x.id != item.id),
        }));
    };
    useEffect(() => {
        void refreshNotice(true);
    }, [refreshNotice]);
    useSocket({
        onMessage: (data) => {
            if (data?.type == MESSAGE_REFRESH) {
                void refreshNotice();
            }
        },
    });

    /**
     * 退出登录
     */
    const confirmOutLogin = () => {
        modal.confirm({
            title: "退出登录",
            content: "确定退出当前账号吗？",
            okText: "确 定",
            okButtonProps: {danger: true, type: "default"},
            cancelText: "取 消",
            mask: {
                closable: true,
            },
            onOk: async () => {
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
            },
        } as any);
    }

    return <>
        <Tooltip title={theme?.getModeName(theme?.mode as any)}>
            <Icon type={theme?.isDarkMode() ? 'icon-dark' : (theme?.isAutoMode() ? 'icon-auto' : 'icon-light')}
                  className={icon}
                  onClick={() => {
                      theme?.toggle?.();
                  }}/>
        </Tooltip>
        <Popover
            rootClassName={messageBox}
            placement="bottomRight"
            trigger={["hover", "click"]}
            open={isMobile ? false : notice.open}
            onOpenChange={(open) => {
                if (isMobile) {
                    return;
                }
                setNotice((data) => ({...data, open}));
            }}
            autoAdjustOverflow={true}
            content={(
                <Spin spinning={notice.loading}>
                    <div className={messageList}>
                        <div className="message-head">
                            <Title>
                                <span className="message-head-title">站内消息</span>
                            </Title>
                            <span className="message-head-link" onClick={openNoticePage}>
                                查看全部
                            </span>
                        </div>
                        <div className="message-body">
                            {notice.list.length <= 0 ? (
                                <Empty image={Empty.PRESENTED_IMAGE_SIMPLE}/>
                            ) : notice.list.map((item) => (
                                <div key={item.id} className={messageItem} onClick={() => openNotice(item)}>
                                    <div className="item-title">
                                        <span className="unread-dot"/>
                                        <span>{item.title || "未命名消息"}</span>
                                    </div>
                                    <div className="item-time">
                                        <Tooltip title={item.create_time || "-"}>
                                            <span>
                                                {item.create_time ? (ago(item.create_time) || item.create_time) : "-"}
                                            </span>
                                        </Tooltip>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                </Spin>
            )}
        >
            <Icon
                type="icon-notice"
                className={`${icon} ${messageTrigger}`}
                data-count={notice.total > 0 ? noticeCount : undefined}
                onClick={(event) => {
                    if (!isMobile) {
                        return;
                    }
                    event.preventDefault();
                    event.stopPropagation();
                    openNoticePage();
                }}
            />
        </Popover>
        <ProModal
            ref={detailModalRef}
            title={<Title>消息详情</Title>}
            width={560}
            onCancel={() => detailModalRef.current?.hide()}
            afterClose={() => setDetail(undefined)}
            modalProps={{
                footer: null,
                destroyOnHidden: true,
            } as any}
        >
            {detailLoading ? (
                <div className={messageDetailSpin}>
                    <Spin/>
                </div>
            ) : (
                <div className={messageDetail}>
                    <div className="detail-header">
                        <h2>{detail?.title || "消息详情"}</h2>
                        <div className="meta">
                            <span>创建时间：{detail?.create_time || "-"}</span>
                        </div>
                    </div>
                    <div
                        className="detail-content"
                        dangerouslySetInnerHTML={{__html: detail?.content || "暂无消息内容"}}
                    />
                </div>
            )}
        </ProModal>
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
                         <li className={menuItem} onClick={() => historyPush("user.person.info")}>
                             <div className={menuItemLabel}><Icon type="SettingOutlined" className={menuItemLeadIcon}/>账户管理
                             </div>
                             <Icon type="RightOutlined"/>
                         </li>
                         <li className={menuItem} onClick={() => historyPush("user.person.session")}>
                             <div className={menuItemLabel}><Icon type="DeploymentUnitOutlined"
                                                                  className={menuItemLeadIcon}/>会话管理
                             </div>
                             <Icon type="RightOutlined"/>
                         </li>
                         <li className={menuItem} onClick={() => historyPush("user.person.certificate")}>
                             <div className={menuItemLabel}><Icon type="SafetyCertificateOutlined"
                                                                  className={menuItemLeadIcon}/>实名认证
                             </div>
                             <Icon type="RightOutlined"/>
                         </li>
                         <li className={menuItem} onClick={confirmOutLogin}>
                             <div className={menuItemLabel}><Icon type="LogoutOutlined" className={menuItemLeadIcon}/>退出登录
                             </div>
                             <Icon type="RightOutlined"/>
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
        }
    };
});

/**
 * 客服弹窗组件样式
 */
const useCustomerStyles: any = createStyles(({token, isDarkMode}: any): any => {
    return {
        container: {
            position: "relative" as const,
            overflow: "hidden",
            borderRadius: token.borderRadiusLG + "px",
            display: "flex",
            flexDirection: "column" as const,
        },
        closeBtn: {
            position: "absolute" as const,
            top: "16px",
            right: "16px",
            width: "32px",
            height: "32px",
            borderRadius: "50%",
            background: "rgba(255, 255, 255, 0.2)",
            border: "none",
            color: "white",
            fontSize: "16px",
            cursor: "pointer",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            transition: "all 0.3s ease",
            zIndex: 10,
            backdropFilter: "blur(10px)",
            "&:hover": {
                background: "rgba(255, 255, 255, 0.3)",
                transform: "rotate(90deg)",
            }
        },
        header: {
            textAlign: "center" as const,
            padding: "32px 28px 28px",
            background: `linear-gradient(135deg, ${token.colorPrimary} 0%, ${token.colorPrimaryHover} 100%)`,
            position: "relative" as const,
            overflow: "hidden",
            "&::before": {
                content: '""',
                position: "absolute",
                top: "-50%",
                right: "-20%",
                width: "160px",
                height: "160px",
                background: "rgba(255, 255, 255, 0.1)",
                borderRadius: "50%",
            },
            "&::after": {
                content: '""',
                position: "absolute",
                bottom: "-30%",
                left: "-10%",
                width: "120px",
                height: "120px",
                background: "rgba(255, 255, 255, 0.08)",
                borderRadius: "50%",
            },
            ".service-title": {
                fontSize: "18px",
                fontWeight: "700",
                color: "white",
                marginBottom: "6px",
                position: "relative",
                zIndex: 1,
                letterSpacing: "0.5px",
            },
            ".service-subtitle": {
                fontSize: "12px",
                color: "rgba(255, 255, 255, 0.9)",
                position: "relative",
                zIndex: 1,
            }
        },
        body: {
            padding: "24px",
            paddingBottom: "0",
            background: isDarkMode ? token.colorBgContainer : "#f8f9fa",
            minHeight: "160px",
            flex: 1,
        },
        contactBox: {
            background: isDarkMode ? token.colorBgElevated : "white",
            borderRadius: token.borderRadiusLG + "px",
            padding: "24px",
            border: `1px solid ${isDarkMode ? "rgba(255,255,255,0.1)" : "rgba(0,0,0,0.08)"}`,
            boxShadow: isDarkMode ? "0 6px 20px rgba(0,0,0,0.4)" : "0 6px 20px rgba(0,0,0,0.08)",
            marginBottom: "24px",
            position: "relative" as const,
            display: "flex",
            alignItems: "center",
            gap: "20px",
            "&::before": {
                content: '""',
                position: "absolute",
                left: "24px",
                top: "24px",
                width: "4px",
                height: "24px",
                background: `linear-gradient(180deg, ${token.colorPrimary}, ${token.colorPrimaryHover})`,
                borderRadius: "2px",
            }
        },
        contactLeft: {
            flex: 1,
            paddingLeft: "16px",
        },
        contactLabel: {
            fontSize: "12px",
            color: isDarkMode ? token.colorTextTertiary : token.colorTextSecondary,
            fontWeight: "600",
            textTransform: "uppercase" as const,
            letterSpacing: "0.5px",
            marginBottom: "12px",
        },
        contactContent: {
            fontSize: "15px",
            lineHeight: "1.8",
            color: isDarkMode ? token.colorText : token.colorTextHeading,
            whiteSpace: "pre-wrap" as const,
            wordBreak: "break-word" as const,
            fontWeight: "500",
            ".ant-typography": {
                marginBottom: 0,
                fontSize: "15px",
                color: "inherit",
                lineHeight: "1.8",
            }
        },
        copyBtn: {
            height: "32px",
            borderRadius: token.borderRadius + "px",
            fontWeight: "600",
            fontSize: "13px",
            paddingLeft: "20px",
            paddingRight: "20px",
            transition: "all 0.3s ease",
        },
        footer: {
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            gap: "8px",
            padding: "14px",
            margin: "0 -24px 0",
            background: isDarkMode ? token.colorBgElevated : token.colorPrimaryBg,
            ".footer-text": {
                fontSize: "12px",
                color: isDarkMode ? token.colorTextSecondary : token.colorTextTertiary,
                margin: 0,
                fontWeight: "500",
            }
        }
    };
});

/**
 * 客服弹窗组件
 */
const CustomerModal: React.FC<{
    open: boolean;
    onCancel: () => void;
    contact?: string;
}> = ({open, onCancel, contact}) => {
    const {styles} = useCustomerStyles();
    const {message} = App.useApp();

    const handleCopyContact = async () => {
        if (!contact) return;
        try {
            if (navigator.clipboard && window.isSecureContext) {
                await navigator.clipboard.writeText(contact);
                message.success("联系方式已复制到剪贴板");
                return;
            }
            const textarea = document.createElement("textarea");
            textarea.value = contact;
            textarea.style.position = "fixed";
            textarea.style.opacity = "0";
            textarea.style.top = "0";
            textarea.style.left = "0";
            document.body.appendChild(textarea);
            textarea.focus();
            textarea.select();
            textarea.setSelectionRange(0, contact.length);
            const successful = document.execCommand("copy");
            document.body.removeChild(textarea);
            if (successful) {
                message.success("联系方式已复制到剪贴板");
            } else {
                message.error("复制失败，请手动复制");
            }
        } catch {
            message.error("复制失败，请手动复制");
        }
    };

    return (
        <ProModal
            width={420}
            onCancel={onCancel}
            modalProps={{
                open,
                title: null,
                closable: false,
                footer: null,
                centered: false,
                mask: {
                    closable: true,
                },
                style: {
                    top: 40,
                },
                styles: {
                    container: {
                        padding: 0,
                    },
                    body: {
                        padding: 0,
                        borderRadius: styles.container.borderRadius,
                    },
                    content: {
                        padding: 0,
                        borderRadius: styles.container.borderRadius,
                    },
                }
            } as any}
        >
            <div className={styles.container}>
                <button className={styles.closeBtn} onClick={onCancel}>
                    <Icon type="CloseOutlined"/>
                </button>

                <div className={styles.header}>
                    <div className="service-title">客户服务中心</div>
                    <div className="service-subtitle">我们随时为您提供专业的服务支持</div>
                </div>

                <div className={styles.body}>
                    <div className={styles.contactBox}>
                        <div className={styles.contactLeft}>
                            <div className={styles.contactLabel}>联系方式</div>
                            <div className={styles.contactContent}>
                                <Typography.Paragraph copyable={false}>
                                    {contact || "暂未设置联系方式"}
                                </Typography.Paragraph>
                            </div>
                        </div>
                        {contact && (
                            <Button
                                type="primary"
                                className={styles.copyBtn}
                                onClick={handleCopyContact}
                                ghost
                                size="small"
                            >
                                联系
                            </Button>
                        )}
                    </div>

                    <div className={styles.footer}>
                        <span className="footer-text">服务时间: 周一至周日 9:00 - 22:00</span>
                    </div>
                </div>
            </div>
        </ProModal>
    );
};

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
    const [customerVisible, setCustomerVisible] = useState(false);
    const layoutRef = useRef<SinKingRef | null>(null);
    const {styles: {menuColor, collapsedImg, collapsedLogo, unCollapsed}} = useSKLayoutStyles();

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
        user?.refreshWebUser((data: any) => {
            setLoading(false);
            if (!data) {
                historyPush("user.login");
            }
        });
    }, []);

    return <>
        <CustomerModal
            open={customerVisible}
            onCancel={() => setCustomerVisible(false)}
            contact={web?.info?.contact}
        />
        <Layout ref={layoutRef}
                loading={loading || !user?.web}
                pathname={location?.pathname}
                matchedRoutes={match || []}
                onNavigate={(path) => history.push(path)}
                breadCrumbItems={breadCrumbItems}
                hideBreadCrumb={match?.at(-1)?.route?.hideBreadCrumb}
                waterMark={web?.info?.ui?.user?.watermark ? [user?.web?.nick_name, user?.web?.email] : ""}
                menus={getCurrentMenus(location?.pathname, true)}
                menuClassNames={menuColor}
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
                menuBottomBtnIcon={web?.info?.contact ? <Icon type="CustomerServiceOutlined"/> : false}
                menuBottomBtnText={web?.info?.contact ? <>&nbsp;&nbsp;客户服务</> : false}
                onMenuBottomBtnClick={() => {
                    setCustomerVisible(true);
                }}
                collapsedLogo={() => {
                    return web?.info?.logo ? <img
                            src={web?.info?.logo || getPath("logo.svg")}
                            alt={web?.info?.name} className={collapsedImg}/> :
                        <Icon type="icon-logo" style={{color: web?.info?.ui?.user?.color || defaultSettings?.color}}
                              className={collapsedLogo}/>
                }}
                unCollapsedLogo={() => {
                    return <div className={unCollapsed}>
                        {web?.info?.logo ? <img
                                src={web?.info?.logo || getPath("logo.svg")}
                                alt={web?.info?.name || Settings?.title}/> :
                            <Icon type="icon-logo"
                                  style={{color: web?.info?.ui?.user?.color || defaultSettings?.color}}/>}
                        <div
                            style={{color: web?.info?.ui?.user?.color ? web?.info?.ui?.user?.color : defaultSettings?.color}}>
                            {web?.info?.name || Settings?.title}
                        </div>
                    </div>
                }}>
            <Outlet context={{layoutRef}}/>
        </Layout>
    </>;
}


export default SKLayout;
