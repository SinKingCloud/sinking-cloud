import React, {useEffect, useState} from "react";
import {App, Form} from "antd";
import {BaseForm, FormField, Loading} from "./common";
import {useSystemConfig} from "../hooks";
import {payConfig} from "@/service/admin/config";
import {getEnum} from "@/utils/enum";
import {buildOptions} from "@/utils/page";

const PayView: React.FC = () => {
    const {message} = App.useApp();
    const [channelOptions, setChannelOptions] = useState<any[]>([]);
    const [channelLoaded, setChannelLoaded] = useState(false);
    const {form, dataLoading, submitLoading, saveConfig, resetForm} = useSystemConfig({
        request: payConfig,
        format: (values: any) => ({
            min: Number(values?.min || 0),
            max: Number(values?.max || 0),
            balance: Number(values?.balance || 0),
            alipay: Number(values?.alipay || 0),
            wxpay: Number(values?.wxpay || 0),
            expire: Number(values?.expire || 5),
        }),
        normalize: (data: any) => ({
            min: Number(data?.min || 0),
            max: Number(data?.max || 0),
            balance: Number(data?.balance || 0),
            alipay: Number(data?.alipay || 0),
            wxpay: Number(data?.wxpay || 0),
            expire: Number(data?.expire || 5),
        }),
    });

    useEffect(() => {
        if (channelLoaded) {
            return;
        }
        getEnum("pay.channel").then((data: any) => {
            setChannelOptions([
                {label: "关闭", value: 0},
                ...buildOptions(data?.name, true),
            ]);
            setChannelLoaded(true);
        }).catch((err: any) => {
            message?.error(err?.message || "加载支付通道失败");
        });
    }, [channelLoaded]);

    if (dataLoading) {
        return <Loading/>;
    }

    return (
        <Form form={form} layout="vertical">
            <BaseForm form={form} submitLoading={submitLoading} onFinish={saveConfig} onReset={() => resetForm()}>
                <FormField type="number" name="min" label="最小充值金额" min={0} required/>
                <FormField type="number" name="max" label="最大充值金额" min={0} required/>
                <FormField type="select" name="balance" label="账户余额支付"
                           options={[{label: "关闭", value: 0}, {label: "开启", value: 1},]}
                           required/>
                <FormField type="select" name="alipay" label="支付宝支付通道" options={channelOptions} required/>
                <FormField type="select" name="wxpay" label="微信支付通道" options={channelOptions} required/>
                <FormField type="number" name="expire" label="订单有效期(分钟)" min={5} required/>
            </BaseForm>
        </Form>
    );
};

export default PayView;
