import React, {useEffect, useMemo, useState} from "react";
import {Body, Title} from "sinking-antd";
import {Card, Menu, Result} from "antd";
import {createStyles} from "antd-style";
import useControl from "@/utils/control";
import WebView from "./components/web";
import UiView from "./components/ui";
import AliView from "./components/ali";
import PayView from "./components/pay";

type ConfigItem = {
    key: string;
    title: React.ReactNode;
    control: string;
    component: React.ReactNode;
};

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

const configItems: ConfigItem[] = [
    {
        key: "web",
        title: "网站设置",
        control: "admin.system.config:web",
        component: <WebView/>,
    },
    {
        key: "ui",
        title: "界面设置",
        control: "admin.system.config:ui",
        component: <UiView/>,
    },
    {
        key: "aliyun",
        title: "阿里云设置",
        control: "admin.system.config:aliyun",
        component: <AliView/>,
    },
    {
        key: "pay",
        title: "支付设置",
        control: "admin.system.config:pay",
        component: <PayView/>,
    },
];

export default (): React.ReactNode => {
    const {styles} = useStyles();
    const control = useControl("admin.system.config");
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
                <Card title={<Title>系统设置</Title>} variant="borderless">
                    <div className={styles.main}>
                        <div className={styles.leftMenu}>
                            <Menu
                                className={styles.menu}
                                mode="inline"
                                selectedKeys={[activeKey]}
                                onClick={({key}) => setSelectKey(key)}
                                items={visibleItems.map((item) => ({
                                    key: item.key,
                                    label: item.title,
                                }))}
                            />
                        </div>
                        <div className={styles.right}>
                            <div className={styles.title}>{activeItem.title}</div>
                            {activeItem.component}
                        </div>
                    </div>
                </Card>
            )}
        </Body>
    );
};
