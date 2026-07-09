import React from "react";
import {Space, Typography} from "antd";
import {Icon} from "sinking-antd";

export const getEnumText = (data: any, value: any) => (
    data?.[String(value)] || data?.[value] || (value === undefined || value === null || value === "" ? "-" : String(value))
);

const PayType = ({data, value}: { data?: any; value: any }): React.ReactNode => {
    const text = getEnumText(data, value);
    const key = String(value || "").toLowerCase();
    const label = String(text || "").toLowerCase();
    if (key === "alipay" || key === "1" || label.includes("支付宝")) {
        return (
            <Space size={4}>
                <Icon type="AlipayOutlined" style={{color: "#1677ff", fontSize: 14}}/>
                <Typography.Text style={{color: "#1677ff", fontSize: 12, lineHeight: "18px"}}>{text}</Typography.Text>
            </Space>
        );
    }
    if (key === "wxpay" || key === "wechat" || key === "2" || label.includes("微信")) {
        return (
            <Space size={4}>
                <Icon type="WechatOutlined" style={{color: "#07c160", fontSize: 14}}/>
                <Typography.Text style={{color: "#07c160", fontSize: 12, lineHeight: "18px"}}>{text}</Typography.Text>
            </Space>
        );
    }
    if (key === "balance" || key === "0" || label.includes("余额")) {
        return (
            <Space size={4}>
                <Icon type="RedEnvelopeOutlined" style={{color: "#faad14", fontSize: 14}}/>
                <Typography.Text style={{color: "#faad14", fontSize: 12, lineHeight: "18px"}}>{text}</Typography.Text>
            </Space>
        );
    }
    return <Typography.Text style={{fontSize: 12, lineHeight: "18px"}}>{text}</Typography.Text>;
};

export default PayType;
