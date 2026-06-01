import React, {useMemo, useState} from "react";
import {Body} from "sinking-antd";
import {Card, Result, Tabs} from "antd";
import {createStyles} from "antd-style";
import useControl from "@/utils/control";
import {AdminControlCode} from "@/utils/permission";
import FileManage from "./components/file";
import SceneManage from "./components/scene";
import ChannelManage from "./components/channel";

const driverTabs: Array<{
    key: string;
    label: string;
    listControl: string;
    component: (props: { control: any }) => React.ReactNode;
}> = [
    {
        key: "file",
        label: "文件管理",
        listControl: AdminControlCode.SYSTEM_DRIVER_FILE_LIST,
        component: (props) => <FileManage {...props}/>,
    },
    {
        key: "scene",
        label: "场景管理",
        listControl: AdminControlCode.SYSTEM_DRIVER_SCENE_LIST,
        component: (props) => <SceneManage {...props}/>,
    },
    {
        key: "channel",
        label: "通道管理",
        listControl: AdminControlCode.SYSTEM_DRIVER_CHANNEL_LIST,
        component: (props) => <ChannelManage {...props}/>,
    },
];

const useStyles: any = createStyles(({css, token}: any) => ({
    navCard: css`
        .ant-card-body {
            padding: 0 15px;
        }
    `,
    tabs: css`
        .ant-tabs-nav {
            margin: 0;
        }

        .ant-tabs-nav::before {
            display: none;
        }

        .ant-tabs-tab {
            margin: 0;
            padding: 15px 0;
            color: ${token.colorTextSecondary};
        }

        .ant-tabs-tab + .ant-tabs-tab {
            margin-left: 18px;
        }

        .ant-tabs-tab-btn {
            font-size: ${token.fontSize}px;
            font-weight: 400;
            line-height: 20px;
        }

        .ant-tabs-ink-bar {
            height: 3px !important;
            border-radius: ${token.borderRadius}px;
        }

        .ant-tabs-content-holder {
            display: none;
        }
    `,
}));

export default (): React.ReactNode => {
    const {styles} = useStyles();
    const control = useControl();
    const [activeKey, setActiveKey] = useState<string>("file");
    const visibleTabs = useMemo(() => {
        return driverTabs.filter((item) => control.hasControl(item.listControl));
    }, [control.controls]);
    const activeTab = visibleTabs.find((item) => item.key === activeKey) || visibleTabs[0];
    const currentKey = activeTab?.key || visibleTabs[0]?.key;

    const onTabChange = (value: string) => {
        if (value === currentKey) {
            return;
        }
        setActiveKey(value);
    };

    return (
        <Body loading={control.loading}>
            {!control.loading && visibleTabs.length <= 0 && (
                <Result status="403" title="403" subTitle="您无权限访问此页面,请联系管理员"/>
            )}
            {!control.loading && activeTab && (
                <>
                    <Card className={styles.navCard} variant="borderless">
                        <Tabs
                            className={styles.tabs}
                            activeKey={currentKey}
                            onChange={onTabChange}
                            items={visibleTabs.map((item) => ({
                                key: item.key,
                                label: item.label,
                            })) as any}
                        />
                    </Card>
                    {activeTab.component({control})}
                </>
            )}
        </Body>
    );
};
