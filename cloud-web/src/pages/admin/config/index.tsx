import React, {useEffect, useMemo, useState} from "react";
import {Body, Title} from "sinking-antd";
import {Card, Menu, Result} from "antd";
import {createStyles} from "antd-style";
import useControl from "@/utils/control";
import WebView from "./components/web";
import UiView from "./components/ui";
import AliView from "./components/ali";
import PayView from "./components/pay";

type ConfigKey = "web" | "ui" | "aliyun" | "pay";

const useStyles: any = createStyles(({css, responsive, isDarkMode, token}: any) => ({
    main: css`
        display: flex;
        width: 100%;
        height: 100%;
        padding-bottom: 16px;

        .ant-menu-light.ant-menu-inline .ant-menu-item::after {
            top: 19%;
            right: 6px;
            border-top-left-radius: 15px;
            border-bottom-left-radius: 15px;
            height: 60%;
        }

        ${responsive.md} {
            flex-direction: column;
        }
    `,
    leftMenu: css`
        width: 224px;

        .ant-menu-item {
            font-weight: bold;
            font-size: ${token.fontSize}px;
            border-radius: ${token.borderRadius}px;
        }

        ${responsive.md} {
            width: 100%;
            border: none;
            margin-bottom: 10px;
        }
    `,
    right: css`
        flex: 1;
        padding: 8px 40px;
        margin-left: 1px;
        border-left: ${isDarkMode ? "1px solid rgb(50, 50, 50)" : "1px solid #f0f0f0"};

        ${responsive.md} {
            padding: 10px;
            border-left: none;
        }
    `,
    title: css`
        margin-bottom: 20px;
        font-weight: bolder;
        font-size: ${token.fontSizeHeading4}px;
        line-height: 28px;
    `,
    menu: css`
        border-inline-end: none !important;
    `,
}));

const menuMap: Record<ConfigKey, React.ReactNode> = {
    web: "网站设置",
    ui: "界面设置",
    aliyun: "阿里云设置",
    pay: "支付设置",
};

const controlMap: Record<ConfigKey, string> = {
    web: "admin.system.config:web",
    ui: "admin.system.config:ui",
    aliyun: "admin.system.config:aliyun",
    pay: "admin.system.config:pay",
};

const componentMap: Record<ConfigKey, React.ReactNode> = {
    web: <WebView/>,
    ui: <UiView/>,
    aliyun: <AliView/>,
    pay: <PayView/>,
};

export default (): React.ReactNode => {
    const {styles} = useStyles();
    const control = useControl("admin.system.config");
    const menuKeys = useMemo(() => {
        return (Object.keys(menuMap) as ConfigKey[]).filter((key) => control.hasControl(controlMap[key]));
    }, [control.controls]);
    const [selectKey, setSelectKey] = useState<ConfigKey>("web");
    const activeKey = menuKeys.includes(selectKey) ? selectKey : menuKeys[0];

    useEffect(() => {
        if (!control.loading && activeKey && activeKey !== selectKey) {
            setSelectKey(activeKey);
        }
    }, [control.loading, activeKey, selectKey]);

    return (
        <Body loading={control.loading}>
            {!control.loading && menuKeys.length <= 0 && (
                <Result status="403" title="403" subTitle="您无权限访问此页面,请联系管理员"/>
            )}
            {!control.loading && menuKeys.length > 0 && (
                <Card title={<Title>系统设置</Title>} variant="borderless">
                    <div className={styles.main}>
                        <div className={styles.leftMenu}>
                            <Menu
                                className={styles.menu}
                                mode="inline"
                                selectedKeys={[activeKey]}
                                onClick={({key}) => setSelectKey(key as ConfigKey)}
                                items={menuKeys.map((key) => ({
                                    key,
                                    label: menuMap[key],
                                }))}
                            />
                        </div>
                        <div className={styles.right}>
                            <div className={styles.title}>{menuMap[activeKey]}</div>
                            {componentMap[activeKey]}
                        </div>
                    </div>
                </Card>
            )}
        </Body>
    );
};
