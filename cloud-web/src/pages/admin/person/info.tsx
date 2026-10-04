import React, {useRef} from 'react';
import {useModel} from 'umi';
import {Body, Icon, useTheme} from 'sinking-antd';
import {Avatar, Button, Card, Col, Row,} from 'antd';
import {createStyles} from 'antd-style';
import {ago} from '@/utils/time';
import Info, {InfoRef} from '@/pages/admin/person/components/info';
import Account, {AccountRef} from '@/pages/admin/person/components/account';
import Password, {PasswordRef} from '@/pages/admin/person/components/password';
import AvatarUpload, {AvatarRef} from '@/pages/admin/person/components/avatar';
import PhoneComponent, {PhoneRef} from '@/pages/admin/person/components/phone';
import EmailComponent, {EmailRef} from '@/pages/admin/person/components/email';
import Verify, {VerifyRef} from '@/pages/admin/person/components/verify';

const useStyles = createStyles(({css, token, responsive, isDarkMode}: any, props: any = {}) => {
    const isCompactMode = props?.isCompactMode || false;
    const pagePadding = isCompactMode ? 12 : 16;
    const profilePadding = isCompactMode ? 24 : 30;
    const sessionPadding = isCompactMode ? "14px 18px 16px" : "16px 24px 18px";
    const avatarSize = isCompactMode ? 80 : 92;
    const mobileAvatarSize = isCompactMode ? 64 : 68;

    return {
        pageShell: css`
            padding-bottom: ${pagePadding}px;
        `,
        userCard: css`
            border: 1px solid ${token.colorBorderSecondary};
            border-radius: ${token.borderRadius}px;
            overflow: hidden;
            background: ${token.colorBgContainer};

            .profile-shell {
                position: relative;
                overflow: hidden;
                background:
                    linear-gradient(142deg, ${token.colorPrimary} 0%, ${token.colorPrimaryHover} 46%, ${token.colorPrimaryActive} 100%);
                isolation: isolate;

                &::before {
                    content: '';
                    position: absolute;
                    inset: -32px;
                    z-index: -2;
                    background:
                        linear-gradient(118deg, transparent 0%, rgba(255, 255, 255, ${isDarkMode ? 0.04 : 0.08}) 38%, transparent 72%),
                        repeating-linear-gradient(135deg, rgba(255, 255, 255, ${isDarkMode ? 0.025 : 0.05}) 0 1px, transparent 1px 22px);
                    opacity: ${isDarkMode ? 0.46 : 0.58};
                    animation: profileTextureFlow 18s ease-in-out infinite;
                    pointer-events: none;
                }

                &::after {
                    content: '';
                    position: absolute;
                    inset: 0;
                    z-index: -1;
                    background:
                        linear-gradient(90deg, rgba(0, 0, 0, ${isDarkMode ? 0.16 : 0.08}), transparent 38%, rgba(0, 0, 0, ${isDarkMode ? 0.18 : 0.10})),
                        linear-gradient(180deg, rgba(0, 0, 0, ${isDarkMode ? 0.05 : 0.00}), rgba(0, 0, 0, ${isDarkMode ? 0.22 : 0.14}));
                    pointer-events: none;
                }
            }

            @keyframes profileTextureFlow {
                0%, 100% {
                    transform: translate3d(-10px, 0, 0);
                    background-position: 0 0, 0 0;
                }
                50% {
                    transform: translate3d(10px, 0, 0);
                    background-position: 42px 0, 24px 24px;
                }
            }

            .profile-grid {
                position: relative;
                z-index: 1;
                display: flex;
                flex-direction: column;
            }

            .profile-main {
                position: relative;
                min-width: 0;
                padding: ${profilePadding}px;
                display: flex;
                flex-direction: column;
                justify-content: center;
                gap: ${isCompactMode ? 14 : 18}px;
                color: #fff;

                ${responsive.mobile} {
                    padding: ${isCompactMode ? "16px 14px 14px" : "20px 16px 18px"};
                    gap: ${isCompactMode ? 12 : 15}px;
                }
            }

            .profile-main::before {
                content: '';
                position: absolute;
                inset: 0;
                background: linear-gradient(135deg, rgba(255, 255, 255, ${isDarkMode ? 0.03 : 0.06}), transparent 54%);
                pointer-events: none;
            }

            .profile-main-top {
                position: relative;
                z-index: 1;
                display: flex;
                align-items: center;
                gap: ${isCompactMode ? 16 : 20}px;
                min-width: 0;
            }

            .avatar-section {
                min-width: 0;

                .user-avatar {
                    width: ${avatarSize}px !important;
                    height: ${avatarSize}px !important;
                    border: 3px solid rgba(255, 255, 255, 0.92);
                    box-shadow: 0 12px 28px rgba(0, 0, 0, 0.22);
                    cursor: pointer;

                    ${responsive.mobile} {
                        width: ${mobileAvatarSize}px !important;
                        height: ${mobileAvatarSize}px !important;
                    }
                }
            }

            .user-content {
                min-width: 0;
                flex: 1;
            }

            .user-title-line {
                display: flex;
                align-items: center;
                gap: 10px;
                min-width: 0;
                margin-bottom: ${isCompactMode ? 10 : 12}px;
                flex-wrap: wrap;
            }

            .user-name {
                margin: 0;
                min-width: 0;
                color: #fff;
                font-size: ${isCompactMode ? 26 : 30}px;
                font-weight: 800;
                line-height: ${isCompactMode ? 34 : 38}px;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;

                ${responsive.mobile} {
                    font-size: ${isCompactMode ? 22 : 24}px;
                    line-height: ${isCompactMode ? 30 : 32}px;
                }
            }

            .user-id {
                flex: none;
                display: inline-flex;
                align-items: center;
                height: 22px;
                padding: 0 8px;
                border-radius: ${token.borderRadius}px;
                color: rgba(255, 255, 255, 0.82);
                background: rgba(255, 255, 255, 0.14);
                font-size: 12px;
                line-height: 22px;
                backdrop-filter: blur(8px);
            }

            .user-contact {
                display: flex;
                flex-wrap: wrap;
                gap: ${isCompactMode ? 6 : 8}px;
                min-width: 0;

                ${responsive.mobile} {
                    gap: 5px;
                }
            }

            .contact-item {
                display: inline-flex;
                align-items: center;
                gap: 6px;
                flex: 0 1 auto;
                min-width: 0;
                max-width: 100%;
                height: ${isCompactMode ? 28 : 30}px;
                padding: 0 11px;
                border-radius: ${token.borderRadius}px;
                color: rgba(255, 255, 255, 0.88);
                background: rgba(255, 255, 255, 0.13);
                font-size: 12px;
                backdrop-filter: blur(8px);

                ${responsive.mobile} {
                    gap: 4px;
                    height: 24px;
                    padding: 0 8px;
                    font-size: 11px;
                }

                .anticon {
                    flex: none;
                    color: rgba(255, 255, 255, 0.68);
                    font-size: 12px;

                    ${responsive.mobile} {
                        font-size: 11px;
                    }
                }
            }

            .contact-value {
                min-width: 0;
                max-width: 100%;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
            }

            .profile-session {
                min-width: 0;
                display: flex;
                flex-direction: column;
                gap: ${isCompactMode ? 10 : 12}px;
                padding: ${sessionPadding};
                color: #fff;
                background: rgba(255, 255, 255, ${isDarkMode ? 0.06 : 0.08});
                box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.12);
                backdrop-filter: blur(12px);

                ${responsive.mobile} {
                    padding: ${isCompactMode ? "12px 14px 14px" : "14px 16px 16px"};
                }
            }

            .session-head {
                display: flex;
                align-items: center;
                justify-content: space-between;
                gap: 12px;
                min-width: 0;
            }

            .session-title-wrap {
                display: flex;
                align-items: center;
                gap: 8px;
                min-width: 0;
            }

            .session-dot {
                width: 8px;
                height: 8px;
                flex: none;
                border-radius: 50%;
                background: #79f2c0;
                box-shadow: 0 0 0 4px rgba(121, 242, 192, 0.16);
            }

            .session-title {
                min-width: 0;
                color: rgba(255, 255, 255, 0.94);
                font-size: 13px;
                font-weight: 600;
                line-height: 20px;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
            }

            .session-state {
                flex: none;
                height: 22px;
                padding: 0 8px;
                border-radius: ${token.borderRadius}px;
                color: rgba(255, 255, 255, 0.86);
                background: rgba(255, 255, 255, 0.12);
                font-size: 12px;
                line-height: 22px;
            }

            .session-meta {
                display: grid;
                grid-template-columns: repeat(2, minmax(0, 1fr));
                gap: ${isCompactMode ? 10 : 12}px ${isCompactMode ? 14 : 18}px;
                min-width: 0;

                ${responsive.mobile} {
                    grid-template-columns: 1fr;
                }
            }

            .session-meta-item {
                min-width: 0;
                display: flex;
                align-items: flex-start;
                gap: 10px;
                padding: ${isCompactMode ? 6 : 8}px;
                border-radius: ${token.borderRadius}px;
                background: rgba(255, 255, 255, ${isDarkMode ? 0.05 : 0.07});
            }

            .session-meta-icon {
                width: ${isCompactMode ? 26 : 28}px;
                height: ${isCompactMode ? 26 : 28}px;
                flex: none;
                display: flex;
                align-items: center;
                justify-content: center;
                border-radius: 50%;
                background: rgba(255, 255, 255, 0.12);
                color: rgba(255, 255, 255, 0.8);
                font-size: 13px;
            }

            .session-meta-body {
                min-width: 0;
            }

            .session-meta-label {
                display: block;
                margin-bottom: 2px;
                color: rgba(255, 255, 255, 0.66);
                font-size: 12px;
                line-height: 18px;
            }

            .session-meta-value {
                display: block;
                min-width: 0;
                color: rgba(255, 255, 255, 0.94);
                font-size: 13px;
                font-weight: 500;
                line-height: 20px;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
            }
        `,

        settingItem: css`
            padding: 12px;
            background: ${token.colorBgContainer};
            border-radius: ${token.borderRadius}px;
            border: 1px solid ${token.colorBorderSecondary};
            transition: all 0.2s ease;

            &:hover {
                border-color: ${token.colorPrimary};
                box-shadow: 0 2px 8px ${token.boxShadowSecondary};
            }

            ${responsive.mobile} {
                padding: 10px;
            }

            .item-content {
                display: flex;
                align-items: center;
                justify-content: space-between;

                .item-left {
                    display: flex;
                    align-items: center;
                    flex: 1;
                    min-width: 0;

                    .item-icon {
                        font-size: 17px;
                        color: ${token.colorPrimaryBorderHover};
                        margin: 0 15px 0 7px;

                        ${responsive.mobile} {
                            font-size: 18px;
                        }
                    }

                    .item-info {
                        flex: 1;
                        min-width: 0;

                        .item-title {
                            font-size: 14px;
                            font-weight: 400;
                            color: ${token.colorTextSecondary};
                            margin-bottom: 2px;

                            ${responsive.mobile} {
                                font-size: 13px;
                            }
                        }

                        .item-desc {
                            font-size: 11px;
                            color: ${token.colorTextTertiary};
                            line-height: 1.3;
                            overflow: hidden;
                            text-overflow: ellipsis;
                            white-space: nowrap;
                            max-width: 200px;

                            ${responsive.mobile} {
                                font-size: 10px;
                                max-width: 150px;
                            }
                        }
                    }
                }

                .item-right {
                    display: flex;
                    align-items: center;
                    gap: 8px;

                    .item-value {
                        color: ${token.colorTextSecondary};
                        font-size: 13px;
                        max-width: 200px;
                        overflow: hidden;
                        text-overflow: ellipsis;
                        white-space: nowrap;
                        display: flex;
                        align-items: center;
                    }

                    .edit-btn {
                        color: ${token.colorPrimary};
                        padding: 2px 6px;
                        border-radius: ${token.borderRadius / 2}px;
                        font-size: 13px;

                        &:hover {
                            background: ${token.colorPrimaryBg};
                        }
                    }
                }
            }
        `,

        topCard: css`
            border-radius: ${token.borderRadius}px ${token.borderRadius}px 0 0;
            border-bottom: none;
            margin-bottom: 0;
        `,

        bottomCard: css`
            border: 1px solid ${token.colorBorderSecondary};
            border-radius: 0 0 ${token.borderRadius}px ${token.borderRadius}px;
            border-top: none;
            margin-top: 0;
            background: ${token.colorBgContainer};
            box-shadow: ${isDarkMode ? "none" : "0 12px 28px rgba(15, 23, 42, 0.05)"};

            .ant-card-body {
                padding: ${isCompactMode ? 12 : 16}px;
            }
        `,

        settingSection: css`
            padding: 0;
        `,

        avatarStyle: css`
            cursor: pointer;
        `,

    };
});

export default () => {
    const theme = useTheme();
    const isCompactMode = theme?.isCompactTheme?.() || false;
    const {
        styles: {
            pageShell, userCard, settingItem, topCard, bottomCard, settingSection,
            avatarStyle,
        },
    } = useStyles({isCompactMode});
    const user = useModel('user');
    const infoRef = useRef<InfoRef>({} as InfoRef);
    const accountRef = useRef<AccountRef>({} as AccountRef);
    const passwordRef = useRef<PasswordRef>({} as PasswordRef);
    const avatarRef = useRef<AvatarRef>({} as AvatarRef);
    const phoneRef = useRef<PhoneRef>({} as PhoneRef);
    const emailRef = useRef<EmailRef>({} as EmailRef);
    const verifyRef = useRef<VerifyRef>({} as VerifyRef);
    /**
     * 调用独立安全验证流程，供敏感操作前置确认。
     */
    const ensureVerify = () => verifyRef.current?.ensure?.() || Promise.resolve(false);
    const currentPerson: any = user?.system || {};
    const currentSession = currentPerson?.session || {};
    const profileContacts = [
        {
            icon: <Icon type="IdcardOutlined"/>,
            value: currentPerson?.account,
        },
        {
            icon: <Icon type="PhoneOutlined"/>,
            value: currentPerson?.phone,
        },
        {
            icon: <Icon type="MailOutlined"/>,
            value: currentPerson?.email,
        },
    ].filter((item) => item.value !== undefined && item.value !== null && item.value !== "");
    const sessionSummary = [
        {label: '登录设备', value: currentSession?.device || '-', icon: <Icon type="LaptopOutlined"/>},
        {label: '登录IP', value: currentSession?.ip || '-', icon: <Icon type="GlobalOutlined"/>},
        {label: '登录地点', value: currentSession?.location || '-', icon: <Icon type="EnvironmentOutlined"/>},
        {
            label: '登录时间',
            value: currentSession?.login_time ? ago(currentSession.login_time) : '-',
            icon: <Icon type="ClockCircleOutlined"/>,
        },
    ];

    /**
     * 账户设置项配置，集中维护展示内容和操作入口。
     */
    const settingItems = [
        {
            icon: <Icon type="CameraFilled"/>,
            title: '账号头像',
            desc: '设置您的个人头像',
            value: <AvatarUpload
                ref={avatarRef}
                avatar={currentPerson?.avatar}
                nickName={currentPerson?.nick_name}
            />,
            onEdit: () => avatarRef.current?.triggerUpload(),
        },
        {
            icon: <Icon type="EditFilled"/>,
            title: '账号昵称',
            desc: '设置您的显示昵称',
            value: currentPerson?.nick_name || '未设置',
            onEdit: () => infoRef.current?.show(),
        },
        {
            icon: <Icon type="IdcardFilled"/>,
            title: '登录账号',
            desc: '账号仅支持设置一次',
            value: currentPerson?.account || '未设置',
            onEdit: () => accountRef.current?.show(),
        },
        {
            icon: <Icon type="PhoneFilled"/>,
            title: '绑定手机',
            desc: '用于登录和安全验证',
            value: currentPerson?.phone || '未绑定',
            onEdit: () => phoneRef.current?.show(),
        },
        {
            icon: <Icon type="MailFilled"/>,
            title: '绑定邮箱',
            desc: '用于通知和安全验证',
            value: currentPerson?.email || '未绑定',
            onEdit: () => emailRef.current?.show(),
        },
        {
            icon: <Icon type="SafetyCertificateFilled"/>,
            title: '登录密码',
            desc: '用于账户登录验证',
            value: currentPerson?.password ? '已设置' : '未设置',
            onEdit: () => passwordRef.current?.show(),
        },
    ];

    /**
     * 渲染单个账户设置项。
     */
    const renderSettingItem = (
        icon: React.ReactNode,
        title: string,
        desc: string,
        value: string | React.ReactNode,
        onEdit?: () => void,
    ) => {
        return (
            <div className={settingItem}>
                <div className="item-content">
                    <div className="item-left">
                        <div className="item-icon">
                            {icon}
                        </div>
                        <div className="item-info">
                            <div className="item-title">{title}</div>
                            <div className="item-desc">{desc}</div>
                        </div>
                    </div>
                    <div className="item-right">
                        <div className="item-value">
                            {value}
                        </div>
                        <Button
                            type="text"
                            className="edit-btn"
                            onClick={onEdit}
                        >
                            修改
                        </Button>
                    </div>
                </div>
            </div>
        );
    };

    return (
        <Body loading={!user?.system}>
            <div className={pageShell}>
                <Row justify="center">
                    <Col xs={24} sm={24} md={22} lg={20} xl={18} xxl={16} xxxl={14}>
                        <Row gutter={[15, 0]}>
                            <Col span={24}>
                                <div className={`${userCard} ${topCard}`}>
                                    <div className="profile-shell">
                                        <div className="profile-grid">
                                            <div className="profile-main">
                                                <div className="profile-main-top">
                                                    <div className="avatar-section">
                                                        <Avatar
                                                            size={85}
                                                            src={currentPerson?.avatar}
                                                            icon={<Icon type="UserOutlined"/>}
                                                            className={`user-avatar ${avatarStyle}`}
                                                            onClick={() => avatarRef.current?.triggerUpload()}
                                                        />
                                                    </div>
                                                    <div className="user-content">
                                                        <div className="user-title-line">
                                                            {currentPerson?.nick_name && (
                                                                <h2 className="user-name">{currentPerson.nick_name}</h2>
                                                            )}
                                                            <span className="user-id">UID {currentPerson?.id || '-'}</span>
                                                        </div>
                                                        {profileContacts.length > 0 && (
                                                            <div className="user-contact">
                                                                {profileContacts.map((item, index) => (
                                                                    <span className="contact-item" key={index}>
                                                                        {item.icon}
                                                                        <span className="contact-value">{item.value}</span>
                                                                    </span>
                                                                ))}
                                                            </div>
                                                        )}
                                                    </div>
                                                </div>
                                            </div>
                                            <div className="profile-session">
                                                <div className="session-head">
                                                    <div className="session-title-wrap">
                                                        <span className="session-dot"/>
                                                        <span className="session-title">当前会话</span>
                                                    </div>
                                                    <span className="session-state">在线</span>
                                                </div>
                                                <div className="session-meta">
                                                    {sessionSummary.map((item) => (
                                                        <div className="session-meta-item" key={item.label}>
                                                            <span className="session-meta-icon">{item.icon}</span>
                                                            <span className="session-meta-body">
                                                                <span className="session-meta-label">{item.label}</span>
                                                                <span className="session-meta-value">{item.value}</span>
                                                            </span>
                                                        </div>
                                                    ))}
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </Col>
                            <Col span={24}>
                                <Card className={bottomCard} variant="borderless">
                                    <div className={settingSection}>
                                        <Row gutter={[16, 16]}>
                                            {settingItems.map((item, index) => (
                                                <Col xs={24} sm={24} md={24} lg={12} xl={12} xxl={12} xxxl={12}
                                                     key={index}>
                                                    {renderSettingItem(item.icon, item.title, item.desc, item.value, item.onEdit)}
                                                </Col>
                                            ))}
                                        </Row>
                                    </div>
                                </Card>
                            </Col>
                        </Row>
                    </Col>
                </Row>
            </div>
            <Info ref={infoRef} nickName={currentPerson?.nick_name}/>
            <Account ref={accountRef} account={currentPerson?.account}/>
            <Password ref={passwordRef} ensureVerify={ensureVerify}/>
            <PhoneComponent ref={phoneRef} phone={currentPerson?.phone} ensureVerify={ensureVerify}/>
            <EmailComponent ref={emailRef} email={currentPerson?.email} ensureVerify={ensureVerify}/>
            <Verify ref={verifyRef}/>
        </Body>
    );
};
