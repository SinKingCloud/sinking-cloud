import React, {useEffect, useMemo, useState} from "react";
import {Body, Icon, useTheme} from "sinking-antd";
import {Menu, Result} from "antd";
import {createStyles} from "antd-style";
import useControl from "@/utils/control";
import {AdminControlCode, AdminPageCode} from "@/utils/permission";
import WebView from "./components/web";
import UiView from "./components/ui";
import AliView from "./components/ali";
import PayView from "./components/pay";

const useStyles = createStyles(({css, token, isDarkMode}, {isCompactMode}: {isCompactMode: boolean}) => ({
    page: css`
        width: 100%;
        max-width: ${isCompactMode ? 1280 : 1450}px;
        margin: 0 auto;
    `,
    workspace: css`
        container-name: setting-workspace;
        container-type: inline-size;
        min-width: 0;
        overflow: hidden;
        border-radius: ${token.borderRadiusLG}px;
        background: ${token.colorBgContainer};
        box-shadow: ${token.boxShadowTertiary};
    `,
    hero: css`
        position: relative;
        min-height: ${isCompactMode ? 104 : 118}px;
        display: flex;
        align-items: center;
        justify-content: space-between;
        overflow: hidden;
        isolation: isolate;
        background: linear-gradient(
            118deg,
            color-mix(in srgb, ${token.colorPrimary}, ${token.colorBgContainer} ${isDarkMode ? 94 : 93}%) 0%,
            ${token.colorBgContainer} 58%,
            color-mix(in srgb, #13c2c2, ${token.colorBgContainer} ${isDarkMode ? 97 : 96}%) 100%
        );

        .hero-copy {
            position: relative;
            z-index: 2;
            min-width: 0;
            padding: ${isCompactMode ? "14px 18px" : "18px 22px"};
        }

        .eyebrow {
            display: flex;
            align-items: center;
            gap: 7px;
            color: color-mix(in srgb, ${token.colorPrimary}, ${token.colorTextTertiary} 62%);
            font-size: 9px;
            font-weight: 600;
            line-height: 16px;
        }

        .status-dot {
            width: 5px;
            height: 5px;
            flex: none;
            border-radius: 50%;
            background: color-mix(in srgb, ${token.colorPrimary}, #fff 26%);
            box-shadow: 0 0 0 4px color-mix(in srgb, ${token.colorPrimary}, transparent 90%);
        }

        h1 {
            margin: 3px 0 0;
            color: ${token.colorTextHeading};
            font-size: ${isCompactMode ? 19 : 20}px;
            font-weight: 600;
            line-height: ${isCompactMode ? 26 : 28}px;
        }

        .hero-visual {
            position: absolute;
            z-index: 1;
            inset: 0 82px 0 auto;
            width: 52%;
            opacity: ${isDarkMode ? .12 : .25};
            pointer-events: none;
        }

        svg { display: block; width: 100%; height: 100%; color: ${token.colorPrimary}; }
        .flow-line { fill: none; stroke: currentColor; stroke-width: 1; stroke-linecap: round; stroke-linejoin: round; stroke-opacity: .16; }
        .flow-primary { stroke-width: 1.6; stroke-opacity: .5; }
        .flow-dash { stroke-dasharray: 3 10; stroke-opacity: .22; }
        .node-halo { fill: currentColor; fill-opacity: .055; stroke: currentColor; stroke-opacity: .13; }
        .node-core { fill: currentColor; fill-opacity: .52; }

        @container setting-workspace (max-width: 600px) {
            min-height: ${isCompactMode ? 104 : 116}px;
            .hero-copy { padding: 16px 14px; }
            h1 { font-size: 18px; line-height: 25px; }
            .hero-visual { right: 42px; width: 68%; opacity: ${isDarkMode ? .07 : .14}; }
            .flow-detail { display: none; }
        }

        @container setting-workspace (max-width: 430px) {
            .hero-visual { right: 16px; width: 76%; opacity: ${isDarkMode ? .045 : .09}; }
        }

        @container setting-workspace (max-width: 360px) {
            .hero-visual { display: none; }
        }
    `,
    settingsMain: css`
        display: flex;
        width: 100%;
        min-width: 0;
        padding: ${isCompactMode ? 16 : token.paddingLG}px;
        box-sizing: border-box;

        .ant-menu-light.ant-menu-inline .ant-menu-item::after {
            top: 19%;
            right: 6px;
            height: 60%;
            border-top-left-radius: 15px;
            border-bottom-left-radius: 15px;
        }

        @container setting-workspace (max-width: 760px) {
            flex-direction: column;
            padding: ${isCompactMode ? 10 : 12}px;
        }
    `,
    leftMenu: css`
        width: ${isCompactMode ? 200 : 224}px;
        flex: none;
        min-width: 0;

        .ant-menu-item {
            border-radius: ${token.borderRadius}px;
            font-size: ${token.fontSize}px;
            font-weight: 600;
        }

        @container setting-workspace (max-width: 760px) {
            width: 100%;
            margin-bottom: 10px;
            overflow: hidden;
        }
    `,
    menuScroll: css`
        min-width: 0;

        @container setting-workspace (max-width: 760px) {
            width: 100%;
            overflow-x: auto;
            overflow-y: hidden;
            overscroll-behavior-x: contain;
            scrollbar-width: none;
            touch-action: pan-x;
            -webkit-overflow-scrolling: touch;
            &::-webkit-scrollbar { display: none; }
        }
    `,
    menu: css`
        border-inline-end: none !important;

        @container setting-workspace (max-width: 760px) {
            display: flex !important;
            flex-direction: row;
            flex-wrap: nowrap;
            gap: 6px;
            width: max-content;
            min-width: 100%;
            padding: 0;
            background: transparent;

            .ant-menu-item {
                display: flex;
                flex: 0 0 auto;
                align-items: center;
                justify-content: center;
                width: auto !important;
                min-width: ${isCompactMode ? 92 : 96}px;
                height: ${isCompactMode ? 32 : 34}px;
                min-height: ${isCompactMode ? 32 : 34}px;
                margin: 0 !important;
                padding-inline: ${isCompactMode ? 10 : 12}px !important;
                background: ${token.colorFillQuaternary};
                color: ${token.colorTextSecondary};
                font-weight: 400;
                line-height: ${isCompactMode ? 32 : 34}px;
                white-space: nowrap;
            }

            .ant-menu-item::after { display: none !important; }
            .ant-menu-item-selected {
                background: ${token.colorPrimaryBg} !important;
                color: ${token.colorPrimary} !important;
                font-weight: 600;
            }
        }
    `,
    right: css`
        min-width: 0;
        flex: 1;
        margin-left: 1px;
        padding: ${isCompactMode ? "6px 28px" : "8px 40px"};
        border-left: 1px solid ${token.colorSplit};

        @container setting-workspace (max-width: 760px) {
            margin-left: 0;
            padding: ${isCompactMode ? 8 : 10}px;
            border-left: none;
        }
    `,
    title: css`
        margin-bottom: ${isCompactMode ? 16 : 20}px;
        color: ${token.colorTextHeading};
        font-size: ${token.fontSizeHeading4}px;
        font-weight: 600;
        line-height: ${isCompactMode ? 24 : 28}px;
    `,
}));

type ConfigItem = {
    key: string;
    title: React.ReactNode;
    icon: string;
    control: string;
    component: React.ReactNode;
};

const configItems: ConfigItem[] = [
    {
        key: "web",
        title: "网站设置",
        icon: "GlobalOutlined",
        control: AdminControlCode.SYSTEM_CONFIG_WEB,
        component: <WebView/>,
    },
    {
        key: "ui",
        title: "界面设置",
        icon: "LayoutOutlined",
        control: AdminControlCode.SYSTEM_CONFIG_UI,
        component: <UiView/>,
    },
    {
        key: "aliyun",
        title: "阿里云设置",
        icon: "CloudOutlined",
        control: AdminControlCode.SYSTEM_CONFIG_ALIYUN,
        component: <AliView/>,
    },
    {
        key: "pay",
        title: "支付设置",
        icon: "WalletOutlined",
        control: AdminControlCode.SYSTEM_CONFIG_PAY,
        component: <PayView/>,
    },
];

export default (): React.ReactNode => {
    const theme = useTheme();
    const {styles} = useStyles({isCompactMode: Boolean(theme?.isCompactTheme?.())});
    const control = useControl(AdminPageCode.SYSTEM_CONFIG);
    const visibleItems = useMemo(() => {
        return configItems.filter((item) => control.hasControl(item.control));
    }, [control.controls]);
    const [selectKey, setSelectKey] = useState<string>("web");
    const activeItem = visibleItems.find((item) => item.key === selectKey) || visibleItems[0];
    const activeKey = activeItem?.key;

    useEffect(() => {
        if (!control.loading && activeKey && activeKey !== selectKey) {
            setSelectKey(activeKey);
        }
    }, [control.loading, activeKey, selectKey]);

    return (
        <Body loading={control.loading}>
            {!control.loading && visibleItems.length <= 0 && (
                <Result status="403" title="403" subTitle="您无权限访问此页面,请联系管理员"/>
            )}
            {!control.loading && activeItem && (
                <div className={styles.page}>
                    <section className={styles.workspace}>
                        <header className={styles.hero}>
                            <div className="hero-copy">
                                <div className="eyebrow"><span className="status-dot"/>SYSTEM SETTINGS</div>
                                <h1>系统设置</h1>
                            </div>
                            <div className="hero-visual">
                                <svg viewBox="0 0 480 132" aria-hidden="true" focusable="false">
                                    <path className="flow-line flow-detail" d="M34 108C116 87 139 29 229 35S346 107 468 52"/>
                                    <path className="flow-line flow-dash" d="M60 30C143 28 174 84 255 86S371 38 454 72"/>
                                    <path className="flow-line flow-primary" d="M26 88C106 84 143 49 215 51S324 94 457 58"/>
                                    <g className="flow-node" transform="translate(154 57)"><circle className="node-halo" r="13"/><circle className="node-core" r="3.4"/></g>
                                    <g className="flow-node" transform="translate(282 78)"><circle className="node-halo" r="10"/><circle className="node-core" r="2.8"/></g>
                                    <g className="flow-node" transform="translate(390 67)"><circle className="node-halo" r="14"/><circle className="node-core" r="3.6"/></g>
                                </svg>
                            </div>
                        </header>
                        <div className={styles.settingsMain}>
                            <div className={styles.leftMenu}>
                                <div className={styles.menuScroll}>
                                    <Menu
                                        className={styles.menu}
                                        mode="inline"
                                        selectedKeys={[activeKey]}
                                        onClick={({key}) => setSelectKey(key)}
                                        items={visibleItems.map((item) => ({
                                            key: item.key,
                                            label: item.title,
                                            icon: <Icon type={item.icon}/>,
                                        }))}
                                    />
                                </div>
                            </div>
                            <section className={styles.right}>
                                <div className={styles.title}>{activeItem.title}</div>
                                {activeItem.component}
                            </section>
                        </div>
                    </section>
                </div>
            )}
        </Body>
    );
};
