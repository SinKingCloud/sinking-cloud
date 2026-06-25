import React from "react";
import {Icon} from "sinking-antd";

import {Space, Typography} from "antd";
import {createStyles} from "antd-style";

type PayTypeProps = {
    data?: any;
    value: any;
};

const useStyles: any = createStyles((): any => {
    return {
        text: {
            fontSize: 12,
            lineHeight: "18px",
        },
        walletIcon: {
            color: "#faad14",
            fontSize: 14,
        },
        alipayIcon: {
            color: "#1677ff",
            fontSize: 14,
        },
        wechatIcon: {
            color: "#07c160",
            fontSize: 14,
        },
    }
});

const enumText = (data: any, value: any) => (
    data?.[String(value)] || data?.[value] || (value === undefined || value === null || value === "" ? "-" : String(value))
);

const PayType = ({data, value}: PayTypeProps): React.ReactNode => {
    const {styles} = useStyles();
    const text = enumText(data, value);
    const key = String(value || "").toLowerCase();
    const label = String(text || "").toLowerCase();
    if (key === "alipay" || key === "1" || label.includes("支付宝")) {
        return (
            <Space size={4}>
                <Icon type="AlipayOutlined" className={styles.alipayIcon}/>
                <Typography.Text className={styles.text} style={{color: "#1677ff"}}>{text}</Typography.Text>
            </Space>
        );
    }
    if (key === "wxpay" || key === "wechat" || key === "2" || label.includes("微信")) {
        return (
            <Space size={4}>
                <Icon type="WechatOutlined" className={styles.wechatIcon}/>
                <Typography.Text className={styles.text} style={{color: "#07c160"}}>{text}</Typography.Text>
            </Space>
        );
    }
    if (key === "balance" || key === "0" || label.includes("余额")) {
        return (
            <Space size={4}>
                <Icon type="RedEnvelopeOutlined" className={styles.walletIcon}/>
                <Typography.Text className={styles.text} style={{color: "#faad14"}}>{text}</Typography.Text>
            </Space>
        );
    }
    return <Typography.Text className={styles.text}>{text}</Typography.Text>;
};

export default PayType;
