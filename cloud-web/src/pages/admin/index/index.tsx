import React from 'react';
import {Body, Title} from "sinking-antd";
import {Avatar, Button, Card, Col, Row, Typography} from "antd";
import {createStyles} from 'antd-style';
import {useModel} from "umi";
import {historyPush} from "@/utils/route";
import {AdminPageCode} from "@/utils/permission";

type QuickItem = {
    title: string;
    desc: string;
    route: string;
    tone: string;
    bg: string;
    icon: React.ReactNode;
};

type NoteItem = {
    title: string;
    content: string;
    time: string;
};

const useStyles: any = createStyles(({token, isDarkMode}: any): any => ({
    page: {
        display: 'flex',
        flexDirection: 'column',
        gap: '12px',
    },
    overviewCard: {
        position: 'relative',
        overflow: 'hidden',
        minHeight: '146px',
        borderRadius: token.borderRadiusLG,
        border: 'none',
        background: token.colorPrimaryBg,
        boxShadow: isDarkMode ? 'none' : '0 12px 32px rgba(22, 119, 255, 0.10)',
        '&::before': {
            content: '""',
            position: 'absolute',
            top: 0,
            left: 0,
            width: '100%',
            height: '100%',
            backgroundImage: `
                linear-gradient(${isDarkMode ? 'rgba(255,255,255,0.04)' : 'rgba(22,119,255,0.07)'} 1px, transparent 1px),
                linear-gradient(90deg, ${isDarkMode ? 'rgba(255,255,255,0.04)' : 'rgba(22,119,255,0.07)'} 1px, transparent 1px)
            `,
            backgroundSize: '22px 22px',
            opacity: isDarkMode ? 0.42 : 0.58,
            pointerEvents: 'none',
        },
        '.ant-card-body': {
            position: 'relative',
            zIndex: 1,
            height: '100%',
            padding: '12px 24px',
        },
        '@media (max-width: 520px)': {
            minHeight: 'auto',
            '.ant-card-body': {
                padding: '16px',
            },
        },
    },
    overview: {
        position: 'relative',
        display: 'flex',
        alignItems: 'center',
        minHeight: '128px',
        '.overview-main': {
            position: 'relative',
            zIndex: 1,
            flex: 1,
            minWidth: 0,
            display: 'flex',
            alignItems: 'center',
            gap: '18px',
        },
        '.overview-info': {
            minWidth: 0,
        },
        '.overview-avatar-wrap': {
            position: 'relative',
            width: '140px',
            height: '140px',
            flexShrink: 0,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            color: token.colorPrimary,
        },
        '.overview-avatar-effect': {
            position: 'absolute',
            inset: 0,
            width: '100%',
            height: '100%',
            opacity: isDarkMode ? 0.64 : 0.86,
        },
        '.overview-avatar': {
            position: 'relative',
            zIndex: 1,
            flexShrink: 0,
            border: `4px solid ${isDarkMode ? 'rgba(255,255,255,0.12)' : token.colorBgContainer}`,
            background: token.colorPrimary,
            boxShadow: isDarkMode ? 'none' : '0 12px 26px rgba(22, 119, 255, 0.18)',
        },
        '.overview-avatar img': {
            width: '100%',
            height: '100%',
            objectFit: 'cover',
        },
        '.overview-name': {
            fontSize: '22px',
            fontWeight: 600,
            color: token.colorTextHeading,
            lineHeight: '30px',
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            whiteSpace: 'nowrap',
        },
        '.overview-desc': {
            display: 'block',
            marginTop: '6px',
            fontSize: '14px',
            color: token.colorTextSecondary,
            lineHeight: '22px',
        },
        '.overview-visual': {
            position: 'absolute',
            zIndex: 0,
            top: '50%',
            right: '-18px',
            width: '390px',
            opacity: isDarkMode ? 0.5 : 0.66,
            color: token.colorPrimary,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            transform: 'translateY(-50%)',
            pointerEvents: 'none',
        },
        '.overview-visual svg': {
            width: '390px',
            height: '164px',
        },
        '@media (max-width: 900px)': {
            '.overview-visual': {
                right: '-42px',
                width: '260px',
                opacity: isDarkMode ? 0.38 : 0.5,
            },
            '.overview-visual svg': {
                width: '260px',
                height: '110px',
            },
        },
        '@media (max-width: 520px)': {
            minHeight: '92px',
            '.overview-main': {
                alignItems: 'center',
                gap: '12px',
            },
            '.overview-avatar-wrap': {
                width: '96px',
                height: '96px',
            },
            '.overview-name': {
                fontSize: '18px',
                lineHeight: '26px',
            },
            '.overview-desc': {
                marginTop: '4px',
                fontSize: '13px',
                lineHeight: '20px',
            },
            '.overview-visual': {
                right: '-54px',
                width: '190px',
                opacity: isDarkMode ? 0.3 : 0.4,
            },
            '.overview-visual svg': {
                width: '190px',
                height: '80px',
            },
        },
    },
    contentCard: {
        height: '100%',
        borderRadius: token.borderRadiusLG,
        border: `1px solid ${isDarkMode ? 'rgba(255,255,255,0.08)' : token.colorBorderSecondary}`,
        '.ant-card-head': {
            minHeight: '46px',
            borderBottom: `1px solid ${isDarkMode ? 'rgba(255,255,255,0.06)' : token.colorBorderSecondary}`,
        },
        '.ant-card-head-title': {
            padding: '12px 0',
        },
        '.ant-card-body': {
            padding: '14px',
        },
    },
    quickGrid: {
        display: 'grid',
        gridTemplateColumns: 'repeat(3, minmax(0, 1fr))',
        gap: '10px',
        '@media (max-width: 1080px)': {
            gridTemplateColumns: 'repeat(2, minmax(0, 1fr))',
        },
        '@media (max-width: 576px)': {
            gridTemplateColumns: '1fr',
        },
    },
    quickItem: {
        display: 'flex',
        alignItems: 'center',
        gap: '12px',
        minHeight: '76px',
        padding: '12px',
        borderRadius: token.borderRadius,
        border: `1px solid ${isDarkMode ? 'rgba(255,255,255,0.07)' : token.colorBorderSecondary}`,
        background: isDarkMode ? 'rgba(255,255,255,0.02)' : token.colorBgContainer,
        cursor: 'pointer',
        transition: 'border-color .18s ease, background .18s ease, box-shadow .18s ease',
        '&:hover': {
            borderColor: token.colorPrimaryBorder,
            background: isDarkMode ? 'rgba(255,255,255,0.04)' : token.colorFillQuaternary,
            boxShadow: isDarkMode ? 'none' : '0 6px 16px rgba(0, 0, 0, 0.04)',
        },
        '.quick-icon': {
            width: '38px',
            height: '38px',
            borderRadius: token.borderRadius,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            flexShrink: 0,
        },
        '.quick-main': {
            flex: 1,
            minWidth: 0,
        },
        '.quick-title': {
            display: 'block',
            fontSize: '13px',
            fontWeight: 600,
            color: token.colorTextHeading,
            lineHeight: '20px',
        },
        '.quick-desc': {
            display: 'block',
            marginTop: '4px',
            fontSize: '12px',
            color: token.colorTextTertiary,
            lineHeight: '18px',
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            whiteSpace: 'nowrap',
        },
        '.quick-arrow': {
            color: token.colorTextQuaternary,
            flexShrink: 0,
            transition: 'color .18s ease, transform .18s ease',
        },
        '&:hover .quick-arrow': {
            color: token.colorPrimary,
            transform: 'translateX(2px)',
        },
    },
    noteList: {
        display: 'flex',
        flexDirection: 'column',
    },
    noteItem: {
        display: 'flex',
        gap: '10px',
        padding: '11px 0',
        borderBottom: `1px solid ${isDarkMode ? 'rgba(255,255,255,0.06)' : token.colorBorderSecondary}`,
        '&:first-child': {
            paddingTop: 0,
        },
        '&:last-child': {
            borderBottom: 'none',
            paddingBottom: 0,
        },
        '.note-mark': {
            width: '24px',
            height: '24px',
            borderRadius: token.borderRadiusSM,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            flexShrink: 0,
            color: token.colorPrimary,
            background: isDarkMode ? 'rgba(22,119,255,0.14)' : token.colorFillTertiary,
        },
        '.note-main': {
            flex: 1,
            minWidth: 0,
        },
        '.note-top': {
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            gap: '10px',
        },
        '.note-title': {
            minWidth: 0,
            fontSize: '13px',
            fontWeight: 600,
            color: token.colorTextHeading,
            lineHeight: '20px',
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            whiteSpace: 'nowrap',
        },
        '.note-time': {
            flexShrink: 0,
            fontSize: '12px',
            color: token.colorTextTertiary,
            lineHeight: '18px',
        },
        '.note-content': {
            marginTop: '4px',
            fontSize: '12px',
            color: token.colorTextSecondary,
            lineHeight: '18px',
            display: '-webkit-box',
            WebkitLineClamp: 2,
            WebkitBoxOrient: 'vertical',
            overflow: 'hidden',
        },
    },
}));

const ModuleIcon = ({type}: { type: 'user' | 'order' | 'master' | 'cron' | 'config' | 'driver' }) => {
    const paths = {
        user: (
            <>
                <circle cx="16" cy="12" r="5" stroke="currentColor" strokeWidth="2"/>
                <path d="M7 26C8.5 20 23.5 20 25 26" stroke="currentColor" strokeWidth="2" strokeLinecap="round"/>
            </>
        ),
        order: (
            <>
                <rect x="8" y="6" width="16" height="20" rx="4" stroke="currentColor" strokeWidth="2"/>
                <path d="M12 12H20M12 17H20M12 22H17" stroke="currentColor" strokeWidth="2" strokeLinecap="round"/>
            </>
        ),
        master: (
            <>
                <circle cx="12" cy="13" r="4" stroke="currentColor" strokeWidth="2"/>
                <circle cx="22" cy="15" r="3" stroke="currentColor" strokeWidth="2"/>
                <path d="M6 26C7.2 21.5 17 21.5 18 26M18 25C18.8 22.4 25.5 22.4 26 25" stroke="currentColor" strokeWidth="2" strokeLinecap="round"/>
            </>
        ),
        cron: (
            <>
                <circle cx="16" cy="16" r="10" stroke="currentColor" strokeWidth="2"/>
                <path d="M16 10V16L20 20" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
            </>
        ),
        config: (
            <>
                <path d="M16 11A5 5 0 1 1 16 21A5 5 0 0 1 16 11Z" stroke="currentColor" strokeWidth="2"/>
                <path d="M16 5V8M16 24V27M5 16H8M24 16H27M8.2 8.2L10.4 10.4M21.6 21.6L23.8 23.8M23.8 8.2L21.6 10.4M10.4 21.6L8.2 23.8" stroke="currentColor" strokeWidth="2" strokeLinecap="round"/>
            </>
        ),
        driver: (
            <>
                <rect x="7" y="8" width="18" height="16" rx="4" stroke="currentColor" strokeWidth="2"/>
                <path d="M11 13H21M11 18H17" stroke="currentColor" strokeWidth="2" strokeLinecap="round"/>
                <circle cx="21" cy="19" r="1.6" fill="currentColor"/>
            </>
        ),
    };
    return (
        <svg width="24" height="24" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
            {paths[type]}
        </svg>
    );
};

const ArrowIcon = () => (
    <svg width="16" height="16" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M6 4L10 8L6 12" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round"/>
    </svg>
);

const NoteIcon = () => (
    <svg width="16" height="16" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M4 2.5H10.5L13 5V13.5H4V2.5Z" stroke="currentColor" strokeWidth="1.5" strokeLinejoin="round"/>
        <path d="M10.5 2.5V5H13M6 8H11M6 10.5H10" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/>
    </svg>
);

const AvatarEffect = () => (
    <svg className="overview-avatar-effect" viewBox="0 0 128 128" fill="none" xmlns="http://www.w3.org/2000/svg">
        <defs>
            <radialGradient id="avatarGlow" cx="50%" cy="50%" r="50%">
                <stop offset="0" stopColor="currentColor" stopOpacity=".12"/>
                <stop offset=".58" stopColor="currentColor" stopOpacity=".04"/>
                <stop offset="1" stopColor="currentColor" stopOpacity="0"/>
            </radialGradient>
            <linearGradient id="avatarTrace" x1="14" y1="64" x2="114" y2="64" gradientUnits="userSpaceOnUse">
                <stop stopColor="currentColor" stopOpacity="0"/>
                <stop offset=".5" stopColor="currentColor" stopOpacity=".24"/>
                <stop offset="1" stopColor="currentColor" stopOpacity="0"/>
            </linearGradient>
            <filter id="avatarSpotBlur" x="-30%" y="-30%" width="160%" height="160%">
                <feGaussianBlur stdDeviation="1.8"/>
            </filter>
        </defs>
        <circle cx="64" cy="64" r="60" fill="url(#avatarGlow)"/>
        <circle cx="64" cy="64" r="52" stroke="currentColor" strokeOpacity=".10"/>
        <circle cx="64" cy="64" r="42" stroke="currentColor" strokeOpacity=".08" strokeDasharray="2 10"/>
        <g stroke="url(#avatarTrace)" strokeLinecap="round">
            <path d="M14 44C32 48 42 54 48 62" strokeWidth="1.5"/>
            <path d="M114 44C96 48 86 54 80 62" strokeWidth="1.5"/>
            <path d="M16 84C34 80 44 74 50 66" strokeWidth="1.5"/>
            <path d="M112 84C94 80 84 74 78 66" strokeWidth="1.5"/>
            <path d="M64 8C64 26 64 34 64 42" strokeWidth="1.2" strokeOpacity=".72"/>
        </g>
        <g filter="url(#avatarSpotBlur)" fill="currentColor">
            <circle cx="0" cy="0" r="3.4" fillOpacity=".72">
                <animateTransform attributeName="transform" type="translate" values="14 44; 30 48; 42 56; 48 62" keyTimes="0;0.4;0.78;1" dur="5.8s" repeatCount="indefinite"/>
                <animate attributeName="opacity" values="0;.8;.55;0" keyTimes="0;0.18;0.82;1" dur="5.8s" repeatCount="indefinite"/>
            </circle>
            <circle cx="0" cy="0" r="3.2" fillOpacity=".62">
                <animateTransform attributeName="transform" type="translate" values="114 44; 98 48; 86 56; 80 62" keyTimes="0;0.4;0.78;1" dur="6.4s" repeatCount="indefinite"/>
                <animate attributeName="opacity" values="0;.72;.48;0" keyTimes="0;0.18;0.82;1" dur="6.4s" repeatCount="indefinite"/>
            </circle>
            <circle cx="0" cy="0" r="3.2" fillOpacity=".58">
                <animateTransform attributeName="transform" type="translate" values="16 84; 32 80; 44 72; 50 66" keyTimes="0;0.42;0.78;1" dur="6.8s" repeatCount="indefinite"/>
                <animate attributeName="opacity" values="0;.68;.46;0" keyTimes="0;0.18;0.82;1" dur="6.8s" repeatCount="indefinite"/>
            </circle>
            <circle cx="0" cy="0" r="3.4" fillOpacity=".68">
                <animateTransform attributeName="transform" type="translate" values="112 84; 96 80; 84 72; 78 66" keyTimes="0;0.38;0.78;1" dur="6.1s" repeatCount="indefinite"/>
                <animate attributeName="opacity" values="0;.76;.5;0" keyTimes="0;0.18;0.82;1" dur="6.1s" repeatCount="indefinite"/>
            </circle>
            <circle cx="0" cy="0" r="3" fillOpacity=".55">
                <animateTransform attributeName="transform" type="translate" values="64 8; 64 22; 64 34; 64 42" keyTimes="0;0.42;0.78;1" dur="7.2s" repeatCount="indefinite"/>
                <animate attributeName="opacity" values="0;.62;.42;0" keyTimes="0;0.18;0.82;1" dur="7.2s" repeatCount="indefinite"/>
            </circle>
        </g>
        <circle cx="64" cy="64" r="34" fill="currentColor" fillOpacity=".035"/>
        <circle cx="64" cy="64" r="44" stroke="currentColor" strokeOpacity=".075"/>
    </svg>
);

const WelcomeGraphic = () => (
    <svg viewBox="0 0 420 176" fill="none" xmlns="http://www.w3.org/2000/svg">
        <defs>
            <linearGradient id="cyberLine" x1="40" y1="42" x2="398" y2="132" gradientUnits="userSpaceOnUse">
                <stop stopColor="currentColor" stopOpacity=".16"/>
                <stop offset=".45" stopColor="currentColor" stopOpacity=".88"/>
                <stop offset=".68" stopColor="#13c2c2" stopOpacity=".9"/>
                <stop offset="1" stopColor="#597ef7" stopOpacity=".48"/>
            </linearGradient>
            <linearGradient id="cyberScan" x1="0" y1="0" x2="92" y2="0" gradientUnits="userSpaceOnUse">
                <stop stopColor="#13c2c2" stopOpacity="0"/>
                <stop offset=".46" stopColor="#13c2c2" stopOpacity=".68"/>
                <stop offset="1" stopColor="currentColor" stopOpacity="0"/>
            </linearGradient>
            <radialGradient id="cyberGlowCore" cx="50%" cy="50%" r="50%">
                <stop stopColor="#13c2c2" stopOpacity=".74"/>
                <stop offset=".42" stopColor="currentColor" stopOpacity=".3"/>
                <stop offset="1" stopColor="currentColor" stopOpacity="0"/>
            </radialGradient>
            <filter id="cyberGlow" x="-18%" y="-36%" width="136%" height="172%">
                <feDropShadow dx="0" dy="0" stdDeviation="3" floodColor="#13c2c2" floodOpacity=".44"/>
                <feDropShadow dx="0" dy="0" stdDeviation="8" floodColor="#1677ff" floodOpacity=".22"/>
            </filter>
            <filter id="cyberSoftBlur" x="-20%" y="-40%" width="140%" height="180%">
                <feGaussianBlur stdDeviation="8" result="blur"/>
                <feMerge>
                    <feMergeNode in="blur"/>
                    <feMergeNode in="SourceGraphic"/>
                </feMerge>
            </filter>
            <clipPath id="cyberClip">
                <path d="M226 36H318L364 88L318 140H226L180 88L226 36Z"/>
            </clipPath>
        </defs>
        <g transform="translate(8 0)">
            <g clipPath="url(#cyberClip)" opacity=".72">
                <g stroke="currentColor" strokeOpacity=".13">
                    <path d="M180 62H364"/>
                    <path d="M180 88H364"/>
                    <path d="M180 114H364"/>
                    <path d="M226 36V140"/>
                    <path d="M272 36V140"/>
                    <path d="M318 36V140"/>
                </g>
                <rect x="130" y="36" width="64" height="104" fill="url(#cyberScan)" opacity=".72">
                    <animate attributeName="x" values="130;364;130" dur="5.8s" repeatCount="indefinite"/>
                </rect>
            </g>

            <g stroke="url(#cyberLine)" strokeLinecap="round" filter="url(#cyberGlow)">
                <path d="M32 58H112L180 88H232" strokeWidth="3" strokeDasharray="42 210">
                    <animate attributeName="stroke-dashoffset" values="0;-252" dur="3.8s" repeatCount="indefinite"/>
                </path>
                <path d="M34 118H114L180 88H232" strokeWidth="3" strokeDasharray="36 216">
                    <animate attributeName="stroke-dashoffset" values="0;-252" dur="4.5s" repeatCount="indefinite"/>
                </path>
                <path d="M312 88H352L398 48" strokeWidth="3.2" strokeDasharray="42 200">
                    <animate attributeName="stroke-dashoffset" values="0;-242" dur="4.1s" repeatCount="indefinite"/>
                </path>
                <path d="M312 88H352L400 128" strokeWidth="3.2" strokeDasharray="42 200">
                    <animate attributeName="stroke-dashoffset" values="0;-242" dur="4.7s" repeatCount="indefinite"/>
                </path>
            </g>

            <g filter="url(#cyberSoftBlur)">
                <circle cx="272" cy="88" r="46" fill="url(#cyberGlowCore)">
                    <animate attributeName="opacity" values=".42;.72;.42" dur="4.2s" repeatCount="indefinite"/>
                </circle>
            </g>

            <g filter="url(#cyberGlow)">
                <circle cx="272" cy="88" r="14" fill="currentColor" fillOpacity=".46"/>
                <circle cx="272" cy="88" r="28" stroke="#13c2c2" strokeOpacity=".66" strokeWidth="2.2" strokeDasharray="18 12">
                    <animateTransform attributeName="transform" type="rotate" from="0 272 88" to="360 272 88" dur="7.5s" repeatCount="indefinite"/>
                </circle>
                <circle cx="272" cy="88" r="40" stroke="currentColor" strokeOpacity=".28" strokeWidth="1.6" strokeDasharray="5 12">
                    <animateTransform attributeName="transform" type="rotate" from="360 272 88" to="0 272 88" dur="12s" repeatCount="indefinite"/>
                </circle>
            </g>

            <g fill="currentColor">
                <rect x="222" y="58" width="56" height="5" rx="2.5" fillOpacity=".42"/>
                <rect x="218" y="78" width="42" height="6" rx="3" fillOpacity=".30"/>
                <rect x="284" y="78" width="32" height="6" rx="3" fill="#13c2c2" fillOpacity=".52"/>
                <rect x="222" y="100" width="88" height="6" rx="3" fillOpacity=".24"/>
                <rect x="238" y="118" width="56" height="5" rx="2.5" fill="#597ef7" fillOpacity=".38"/>
            </g>

            <g filter="url(#cyberGlow)">
                <circle cx="112" cy="58" r="4.5" fill="#13c2c2">
                    <animate attributeName="opacity" values=".32;1;.32" dur="3.2s" repeatCount="indefinite"/>
                </circle>
                <circle cx="114" cy="118" r="4.5" fill="#13c2c2">
                    <animate attributeName="opacity" values=".22;.9;.22" dur="3.8s" repeatCount="indefinite"/>
                </circle>
                <circle cx="398" cy="48" r="5" fill="#597ef7">
                    <animate attributeName="opacity" values=".24;.95;.24" dur="4s" repeatCount="indefinite"/>
                </circle>
                <circle cx="400" cy="128" r="5" fill="#597ef7">
                    <animate attributeName="opacity" values=".28;1;.28" dur="4.4s" repeatCount="indefinite"/>
                </circle>
                <circle cx="180" cy="88" r="4" fill="currentColor" fillOpacity=".72"/>
                <circle cx="364" cy="88" r="4" fill="currentColor" fillOpacity=".58"/>
            </g>

            <g stroke="currentColor" strokeOpacity=".20" strokeLinecap="round">
                <path d="M54 34H142"/>
                <path d="M52 144H154"/>
                <path d="M324 28H392"/>
                <path d="M334 150H404"/>
            </g>
        </g>
    </svg>
);

const quickItems: QuickItem[] = [
    {
        title: '用户管理',
        desc: '用户资料、实名、会话',
        route: AdminPageCode.WEB_USER,
        tone: '#1677ff',
        bg: 'rgba(22,119,255,0.1)',
        icon: <ModuleIcon type="user"/>,
    },
    {
        title: '订单管理',
        desc: '订单详情、发货、退款',
        route: AdminPageCode.WEB_ORDER,
        tone: '#fa8c16',
        bg: 'rgba(250,140,22,0.12)',
        icon: <ModuleIcon type="order"/>,
    },
    {
        title: '员工管理',
        desc: '员工资料、角色、会话',
        route: AdminPageCode.SYSTEM_MASTER,
        tone: '#13c2c2',
        bg: 'rgba(19,194,194,0.1)',
        icon: <ModuleIcon type="master"/>,
    },
    {
        title: '任务管理',
        desc: '定时任务、执行日志',
        route: AdminPageCode.SYSTEM_CRON,
        tone: '#722ed1',
        bg: 'rgba(114,46,209,0.1)',
        icon: <ModuleIcon type="cron"/>,
    },
    {
        title: '系统设置',
        desc: '站点、界面、支付配置',
        route: AdminPageCode.SYSTEM_CONFIG,
        tone: '#52c41a',
        bg: 'rgba(82,196,26,0.12)',
        icon: <ModuleIcon type="config"/>,
    },
    {
        title: '存储管理',
        desc: '文件、场景、存储通道',
        route: AdminPageCode.SYSTEM_DRIVER,
        tone: '#eb2f96',
        bg: 'rgba(235,47,150,0.1)',
        icon: <ModuleIcon type="driver"/>,
    },
];

const notes: NoteItem[] = [
    {title: '补单检查', content: '核对发货中订单，确认是否存在需要手动补单的记录。', time: '今天 09:30'},
    {title: '上线前确认', content: '检查支付通道、存储通道和任务脚本配置。', time: '昨天 18:10'},
    {title: '后台优化', content: '继续统一弹窗间距、列表密度和移动端显示效果。', time: '06-15 14:20'},
];

const getGreeting = () => {
    const hour = new Date().getHours();
    if (hour < 12) {
        return '早上好';
    }
    if (hour < 18) {
        return '下午好';
    }
    return '晚上好';
};

export default (): React.ReactNode => {
    const {styles} = useStyles();
    const user = useModel("user");
    const nickName = user?.system?.nick_name || '管理员';

    return (
        <Body>
            <div className={styles.page}>
                <Card className={styles.overviewCard} variant="borderless">
                    <div className={styles.overview}>
                        <div className="overview-main">
                            <div className="overview-avatar-wrap">
                                <AvatarEffect/>
                                <Avatar
                                    size={{xs: 64, sm: 76, md: 88, lg: 88, xl: 88, xxl: 88, xxxl: 88}}
                                    src={user?.system?.avatar || undefined}
                                    alt={nickName}
                                    className="overview-avatar"
                                >
                                    {nickName?.slice(0, 1)}
                                </Avatar>
                            </div>
                            <div className="overview-info">
                                <div className="overview-name">{getGreeting()}，{nickName}</div>
                                <Typography.Text className="overview-desc">
                                    欢迎回来，今天也把重要的事情稳稳推进。
                                </Typography.Text>
                            </div>
                        </div>
                        <div className="overview-visual">
                            <WelcomeGraphic/>
                        </div>
                    </div>
                </Card>

                <Row gutter={[10, 10]}>
                    <Col xs={24} xl={16}>
                        <Card className={styles.contentCard} title={<Title>快捷访问</Title>} variant="borderless">
                            <div className={styles.quickGrid}>
                                {quickItems.map((item) => (
                                    <div
                                        className={styles.quickItem}
                                        key={item.title}
                                        onClick={() => historyPush(item.route)}
                                    >
                                        <div className="quick-icon" style={{color: item.tone, background: item.bg}}>
                                            {item.icon}
                                        </div>
                                        <div className="quick-main">
                                            <span className="quick-title">{item.title}</span>
                                            <span className="quick-desc">{item.desc}</span>
                                        </div>
                                        <span className="quick-arrow">
                                            <ArrowIcon/>
                                        </span>
                                    </div>
                                ))}
                            </div>
                        </Card>
                    </Col>

                    <Col xs={24} xl={8}>
                        <Card
                            className={styles.contentCard}
                            title={<Title>便签事项</Title>}
                            variant="borderless"
                            extra={<Button type="link" size="small">查看全部</Button>}
                        >
                            <div className={styles.noteList}>
                                {notes.map((item) => (
                                    <div className={styles.noteItem} key={item.title}>
                                        <span className="note-mark">
                                            <NoteIcon/>
                                        </span>
                                        <div className="note-main">
                                            <div className="note-top">
                                                <Typography.Text className="note-title">{item.title}</Typography.Text>
                                                <span className="note-time">{item.time}</span>
                                            </div>
                                            <div className="note-content">{item.content}</div>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </Card>
                    </Col>
                </Row>
            </div>
        </Body>
    );
};
