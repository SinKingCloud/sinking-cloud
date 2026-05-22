import React, {useEffect, useRef, useState} from 'react';
import {useModel} from 'umi';
import {Body, Icon} from 'sinking-antd';
import {ago} from '@/utils/time';
import {
    Avatar,
    Button,
    Card,
    Col,
    Row,
} from 'antd';
import {createStyles} from 'antd-style';
import Info, {InfoRef} from '@/pages/admin/person/components/info';
import Account, {AccountRef} from '@/pages/admin/person/components/account';
import Password, {PasswordRef} from '@/pages/admin/person/components/password';
import AvatarUpload, {AvatarRef} from '@/pages/admin/person/components/avatar';
import PhoneComponent, {PhoneRef} from '@/pages/admin/person/components/phone';
import EmailComponent, {EmailRef} from '@/pages/admin/person/components/email';
import Verify, {VerifyRef} from '@/pages/admin/person/components/verify';
import {getPersonInfo} from '@/service/admin/person';

const useStyles = createStyles(({css, token, responsive}) => {
    return {
        userCard: css`
            background: ${token.colorBgContainer};
            border-radius: ${token.borderRadius * 2}px;
            box-shadow: 0 4px 20px ${token.boxShadowSecondary};

            .user-header {
                background: linear-gradient(135deg, ${token.colorPrimary}, ${token.colorPrimaryActive});
                border-radius: ${token.borderRadius}px ${token.borderRadius}px 0 0;
                padding: 20px;
                color: white;
                position: relative;
                overflow: hidden;

                &::before {
                    content: '';
                    position: absolute;
                    top: 0;
                    left: 0;
                    right: 0;
                    bottom: 0;
                    background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 60 60"><defs><pattern id="texture" width="60" height="60" patternUnits="userSpaceOnUse"><circle cx="15" cy="15" r="2" fill="rgba(255,255,255,0.1)"/><circle cx="45" cy="15" r="1.5" fill="rgba(255,255,255,0.08)"/><circle cx="30" cy="30" r="1.8" fill="rgba(255,255,255,0.12)"/><circle cx="15" cy="45" r="1.2" fill="rgba(255,255,255,0.06)"/><circle cx="45" cy="45" r="2.2" fill="rgba(255,255,255,0.09)"/><rect x="0" y="0" width="1" height="1" fill="rgba(255,255,255,0.03)"/><rect x="20" y="10" width="1" height="1" fill="rgba(255,255,255,0.04)"/><rect x="40" y="25" width="1" height="1" fill="rgba(255,255,255,0.03)"/><rect x="10" y="35" width="1" height="1" fill="rgba(255,255,255,0.05)"/><rect x="50" y="50" width="1" height="1" fill="rgba(255,255,255,0.04)"/></pattern></defs><rect width="100%" height="100%" fill="url(%23texture)"/></svg>');
                    opacity: 0.6;
                    z-index: 1;
                }

                &::after {
                    content: '';
                    position: absolute;
                    top: 0;
                    left: 0;
                    right: 0;
                    bottom: 0;
                    background: linear-gradient(135deg,
                    rgba(255, 255, 255, 0.15) 0%,
                    transparent 40%,
                    transparent 60%,
                    rgba(255, 255, 255, 0.08) 100%
                    );
                    pointer-events: none;
                    z-index: 2;
                }

                .user-content {
                    position: relative;
                    z-index: 3;
                    padding-right: 120px;

                    ${responsive.mobile} {
                        padding-right: 92px;
                    }
                }

                .user-name {
                    font-size: 22px;
                    font-weight: 700;
                    margin: 0 0 6px 0;
                    text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
                    overflow: hidden;
                    text-overflow: ellipsis;
                    white-space: nowrap;
                }

                .user-id {
                    font-size: 13px;
                    opacity: 0.9;
                    margin: 0 0 12px 0;
                }

                .user-stats {
                    display: flex;
                    gap: 20px;

                    ${responsive.mobile} {
                        gap: 12px;
                    }

                    .stat-item {
                        min-width: 0;

                        .stat-label {
                            font-size: 11px;
                            opacity: 0.8;
                            margin-bottom: 3px;
                        }

                        .stat-value {
                            font-size: 14px;
                            font-weight: 600;
                            overflow: hidden;
                            text-overflow: ellipsis;
                            white-space: nowrap;
                            max-width: 180px;

                            ${responsive.mobile} {
                                max-width: 105px;
                            }
                        }
                    }
                }
            }

            .avatar-section {
                position: absolute;
                top: 50%;
                right: 20px;
                transform: translateY(-50%);
                z-index: 3;

                .avatar-wrapper {
                    position: relative;

                    .user-avatar {
                        border: 4px solid ${token.colorBgContainer};
                        box-shadow: 0 4px 16px ${token.boxShadowSecondary};
                        transition: all 0.3s ease;

                        &:hover {
                            transform: scale(1.05);
                            box-shadow: 0 6px 20px ${token.boxShadow};
                        }
                    }
                }
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

            .user-header {
                background: linear-gradient(135deg, ${token.colorPrimary}, ${token.colorPrimaryActive}) !important;
            }
        `,

        bottomCard: css`
            border-radius: 0 0 ${token.borderRadius}px ${token.borderRadius}px;
            border-top: none;
            margin-top: 0;

            .ant-card-body {
                padding: 16px;
            }
        `,

        settingSection: css`
            padding: 0;
        `,

        avatarWrapper: css`
            position: relative;
            display: inline-block;
        `,

        avatarStyle: css`
            cursor: pointer;
        `,

        sessionCard: css`
            margin-top: 15px;
            border-radius: ${token.borderRadius}px;

            .ant-card-body {
                padding: 14px 16px;
            }

            .session-wrap {
                display: flex;
                align-items: center;
                gap: 18px;

                ${responsive.mobile} {
                    display: block;
                }
            }

            .session-head {
                display: flex;
                align-items: center;
                gap: 10px;
                flex: none;
                width: 150px;

                ${responsive.mobile} {
                    width: auto;
                    margin-bottom: 12px;
                }
            }

            .session-icon {
                width: 36px;
                height: 36px;
                border-radius: 50%;
                display: flex;
                align-items: center;
                justify-content: center;
                font-size: 17px;
                color: ${token.colorPrimary};
                background: ${token.colorPrimaryBg};
            }

            .session-title {
                font-size: 14px;
                font-weight: 600;
                color: ${token.colorText};
                line-height: 20px;
            }

            .session-desc {
                font-size: 11px;
                color: ${token.colorTextTertiary};
                line-height: 16px;
            }

            .session-grid {
                flex: 1;
                min-width: 0;
                display: grid;
                grid-template-columns: repeat(5, minmax(0, 1fr));
                gap: 10px;

                ${responsive.mobile} {
                    grid-template-columns: 1fr;
                }
            }

            .session-item {
                min-width: 0;
                padding: 9px 10px;
                border-radius: ${token.borderRadius}px;
                border: 1px solid ${token.colorBorderSecondary};
                background: ${token.colorFillQuaternary};
            }

            .session-label {
                display: block;
                margin-bottom: 3px;
                font-size: 12px;
                color: ${token.colorTextTertiary};
            }

            .session-value {
                display: block;
                font-size: 13px;
                color: ${token.colorTextSecondary};
                overflow: hidden;
                white-space: nowrap;
                text-overflow: ellipsis;
            }
        `,
    };
});

export default () => {
    const {
        styles: {
            userCard, settingItem, topCard, bottomCard, settingSection,
            avatarWrapper, avatarStyle, sessionCard,
        },
    } = useStyles();
    const user = useModel('user');
    const infoRef = useRef<InfoRef>({} as InfoRef);
    const accountRef = useRef<AccountRef>({} as AccountRef);
    const passwordRef = useRef<PasswordRef>({} as PasswordRef);
    const avatarRef = useRef<AvatarRef>({} as AvatarRef);
    const phoneRef = useRef<PhoneRef>({} as PhoneRef);
    const emailRef = useRef<EmailRef>({} as EmailRef);
    const verifyRef = useRef<VerifyRef>({} as VerifyRef);
    const [person, setPerson] = useState<any>();
    const [loading, setLoading] = useState(true);

    /**
     * 拉取后台个人资料和当前会话信息。
     */
    const refreshPerson = async () => {
        setLoading(true);
        await getPersonInfo({
            onSuccess: (res: any) => {
                setPerson(res?.data || undefined);
            },
            onFinally: () => {
                setLoading(false);
            },
        });
    };

    /**
     * 页面初始化时加载个人资料。
     */
    useEffect(() => {
        refreshPerson();
    }, []);

    /**
     * 调用独立安全验证流程，供敏感操作前置确认。
     */
    const ensureVerify = () => verifyRef.current?.ensure?.() || Promise.resolve(false);
    const currentPerson: any = person || user?.system || {};
    const currentSession = currentPerson?.session || {};

    /**
     * 子组件保存成功后刷新页面资料。
     */
    const handleSuccess = () => {
        refreshPerson();
    };

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
                onSuccess={handleSuccess}
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
        <Body loading={loading}>
            <Row justify="center">
                <Col xs={24} sm={24} md={24} lg={20} xl={18} xxl={16}>
                    <Row gutter={[15, 0]}>
                        <Col span={24}>
                            <div className={`${userCard} ${topCard}`}>
                                <div className="user-header">
                                    <div className="user-content">
                                        <h2 className="user-name">{currentPerson?.nick_name || ''}</h2>
                                        <p className="user-id">UID: {currentPerson?.id || ''}</p>
                                        <div className="user-stats">
                                            <div className="stat-item">
                                                <div className="stat-label">最近登录IP</div>
                                                <div className="stat-value">{currentSession?.ip || '-'}</div>
                                            </div>
                                            <div className="stat-item">
                                                <div className="stat-label">最近登录时间</div>
                                                <div className="stat-value">
                                                    {currentSession?.login_time ? ago(currentSession.login_time) : '-'}
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                    <div className="avatar-section">
                                        <div className="avatar-wrapper">
                                            <div className={avatarWrapper}>
                                                <Avatar
                                                    size={85}
                                                    src={currentPerson?.avatar}
                                                    icon={<Icon type="UserOutlined"/>}
                                                    className={`user-avatar ${avatarStyle}`}
                                                />
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
                                            <Col xs={24} sm={24} md={12} lg={12} xl={12} xxl={12} key={index}>
                                                {renderSettingItem(item.icon, item.title, item.desc, item.value, item.onEdit)}
                                            </Col>
                                        ))}
                                    </Row>
                                </div>
                            </Card>
                        </Col>
                        <Col span={24}>
                            <Card className={sessionCard} variant="borderless">
                                <div className="session-wrap">
                                    <div className="session-head">
                                        <div className="session-icon">
                                            <Icon type="CloudServerOutlined"/>
                                        </div>
                                        <div>
                                            <div className="session-title">当前会话</div>
                                            <div className="session-desc">本次后台登录信息</div>
                                        </div>
                                    </div>
                                    <div className="session-grid">
                                        <div className="session-item">
                                            <span className="session-label">登录设备</span>
                                            <span className="session-value">{currentSession?.device || '-'}</span>
                                        </div>
                                        <div className="session-item">
                                            <span className="session-label">登录IP</span>
                                            <span className="session-value">{currentSession?.ip || '-'}</span>
                                        </div>
                                        <div className="session-item">
                                            <span className="session-label">登录地点</span>
                                            <span className="session-value">{currentSession?.location || '-'}</span>
                                        </div>
                                        <div className="session-item">
                                            <span className="session-label">登录时间</span>
                                            <span className="session-value">{currentSession?.login_time || '-'}</span>
                                        </div>
                                        <div className="session-item">
                                            <span className="session-label">过期时间</span>
                                            <span className="session-value">{currentSession?.expire_time || '-'}</span>
                                        </div>
                                    </div>
                                </div>
                            </Card>
                        </Col>
                    </Row>
                </Col>
            </Row>
            <Info ref={infoRef} nickName={currentPerson?.nick_name} onSuccess={handleSuccess}/>
            <Account ref={accountRef} account={currentPerson?.account} onSuccess={handleSuccess}/>
            <Password ref={passwordRef} ensureVerify={ensureVerify} onSuccess={handleSuccess}/>
            <PhoneComponent ref={phoneRef} phone={currentPerson?.phone} ensureVerify={ensureVerify} onSuccess={handleSuccess}/>
            <EmailComponent ref={emailRef} email={currentPerson?.email} ensureVerify={ensureVerify} onSuccess={handleSuccess}/>
            <Verify ref={verifyRef}/>
        </Body>
    );
};
