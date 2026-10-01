import React, {forwardRef, useEffect, useImperativeHandle, useMemo, useRef, useState} from "react";
import {App, Button, ConfigProvider, Form, InputNumber, Radio, Spin} from "antd";
import {createStyles} from "antd-style";
import type {ProModalRef} from "sinking-antd";
import {Icon, ProModal, Title} from "sinking-antd";
import useEnum from "@/utils/enum";
import {formatMoney} from "@/utils/string";
import {setPayJumpUrl} from "@/utils/pay";
import {getPayConfig, rechargeBalance, submitPay} from "@/service/user/pay";

export interface RechargeRef {
    show: () => void;
}

const quickAmounts = [10, 30, 50, 100, 200];
const payIcon: Record<string, {icon: string; color?: string}> = {
    alipay: {icon: "AlipayOutlined", color: "#1677ff"},
    wxpay: {icon: "WechatOutlined", color: "#07c160"},
};

const useStyles = createStyles(({css, token}) => ({
    form: css`
        padding-top: 4px;
        .ant-form-item { margin-bottom: 20px; }
    `,
    amounts: css`
        display: grid;
        grid-template-columns: repeat(3, minmax(0, 1fr));
        gap: 8px;
        width: 100%;

        .ant-radio-button-wrapper {
            position: relative;
            display: flex;
            align-items: center;
            justify-content: center;
            min-width: 0;
            height: 52px;
            padding: 0 8px;
            border: 1px solid ${token.colorBorderSecondary};
            border-radius: ${token.borderRadius}px;
            background: ${token.colorFillAlter};
            color: color-mix(in srgb, ${token.colorTextSecondary} 60%, ${token.colorText});
            font-variant-numeric: tabular-nums;
            box-shadow: none;
        }
        .ant-radio-button-wrapper::before { display: none; }
        .amount-option { display: inline-flex; align-items: baseline; gap: 4px; }
        .amount-value { font-size: 18px; font-weight: 600; line-height: 24px; }
        .amount-unit { color: ${token.colorTextTertiary}; font-size: 12px; }
        .custom { font-size: 14px; }
        .ant-radio-button-wrapper:not(.ant-radio-button-wrapper-disabled):hover { border-color: ${token.colorPrimaryBorder}; }
        .ant-radio-button-wrapper-checked:not(.ant-radio-button-wrapper-disabled) { border-color: ${token.colorPrimary}; color: color-mix(in srgb, ${token.colorTextSecondary} 60%, ${token.colorText}); background: ${token.colorPrimaryBg}; box-shadow: none; }
        .ant-radio-button-wrapper-disabled, .ant-radio-button-wrapper-disabled .amount-unit { color: ${token.colorTextDisabled}; }
    `,
    customAmount: css`
        width: 100%;
        height: 40px;
        margin-top: 10px;
        .ant-input-number-input { height: 38px; font-size: 16px; font-variant-numeric: tabular-nums; }
        .ant-input-number-prefix { color: ${token.colorTextTertiary}; }
    `,
    hint: css`
        margin-top: 8px;
        color: ${token.colorTextTertiary};
        font-size: 12px;
        line-height: 20px;
    `,
    payOptions: css`
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
        gap: 8px;
        width: 100%;
        .ant-radio-wrapper { display: flex; flex-direction: row-reverse; align-items: center; justify-content: space-between; gap: 12px; min-width: 0; min-height: 46px; margin: 0; padding: 10px 12px; border: 1px solid ${token.colorBorderSecondary}; border-radius: ${token.borderRadius}px; background: ${token.colorBgContainer}; }
        .ant-radio-wrapper > span:last-child { display: flex; align-items: center; flex: 1; min-width: 0; padding: 0; }
        .ant-radio-wrapper:not(.ant-radio-wrapper-disabled):hover { border-color: ${token.colorPrimaryBorder}; background: ${token.colorFillAlter}; }
        .ant-radio-wrapper-checked:not(.ant-radio-wrapper-disabled) { border-color: ${token.colorPrimaryBorder}; background: ${token.colorBgContainer}; }
        .pay-label { display: flex; align-items: center; gap: 8px; color: ${token.colorTextSecondary}; font-size: 14px; line-height: 22px; }
        .pay-label .anticon { flex: none; width: 20px; height: 20px; font-size: 20px; line-height: 1; }
        .pay-label .anticon svg { display: block; }
        .ant-radio-wrapper-checked .pay-label { color: ${token.colorText}; }
        .ant-radio-wrapper-disabled .pay-label { color: ${token.colorTextDisabled}; }
    `,
    checkout: css`
        padding-top: 4px;
        .ant-form-item-row { flex-wrap: nowrap; align-items: center; }
        .ant-form-item-label { padding: 0; text-align: start; }
        .checkout-value { display: inline-flex; align-items: baseline; gap: 5px; color: ${token.colorText}; font-size: 24px; font-weight: 600; line-height: 32px; font-variant-numeric: tabular-nums; white-space: nowrap; }
        .checkout-currency { color: ${token.colorTextSecondary}; font-size: 14px; font-weight: 400; }
    `,
    submit: css`
        width: 100%;
        height: 40px;
        font-size: 14px;
    `,
}));

export default forwardRef<RechargeRef, {onRefreshBalance: () => Promise<void>}>(({onRefreshBalance}, ref) => {
    const {message, modal} = App.useApp();
    const [form] = Form.useForm();
    const {styles} = useStyles();
    const modalRef = useRef<ProModalRef>(null);
    const [enumData, enumLoading] = useEnum(["pay.pay"]);
    const [config, setConfig] = useState<any>({});
    const [loading, setLoading] = useState(true);
    const [submitting, setSubmitting] = useState(false);
    const [money, setMoney] = useState<number>(quickAmounts[0]);
    const [customAmount, setCustomAmount] = useState(false);
    const [payType, setPayType] = useState<string>();
    const payTypes = enumData["pay.pay"]?.type || {};

    const enabledPayOptions = useMemo(() => Object.entries(payTypes)
        .filter(([value]) => value !== "balance" && config?.[value])
        .map(([value, label]) => ({value, label: String(label)})), [config, payTypes]);
    const selectedPayType = payType || enabledPayOptions[0]?.value;
    const min = Number(config?.min || 0);
    const max = Number(config?.max || 0);
    const amountRange = min > 0 && max > 0 ? `${formatMoney(min)} - ${formatMoney(max)} 元`
        : min > 0 ? `不少于 ${formatMoney(min)} 元` : max > 0 ? `不超过 ${formatMoney(max)} 元` : "";

    useImperativeHandle(ref, () => ({
        show: () => {
            setLoading(true);
            modalRef.current?.show();
            getPayConfig({
                onSuccess: (response: any) => {
                    const nextConfig = response?.data || {};
                    setConfig(nextConfig);
                    const minAmount = Number(nextConfig.min || 0);
                    const maxAmount = Number(nextConfig.max || 0);
                    if ((minAmount > 0 && money < minAmount) || (maxAmount > 0 && money > maxAmount)) {
                        const nextMoney = quickAmounts.find((amount) => amount >= minAmount && (!maxAmount || amount <= maxAmount));
                        setMoney(nextMoney ?? Math.max(1, minAmount));
                        setCustomAmount(nextMoney === undefined);
                    }
                },
                onFail: (response: any) => { setConfig({}); message.error(response?.message || "获取支付配置失败"); },
            }).finally(() => setLoading(false));
        },
    }));

    useEffect(() => {
        if (enabledPayOptions.length > 0 && (!payType || !enabledPayOptions.some((item) => item.value === payType))) {
            const nextPayType = enabledPayOptions[0].value;
            setPayType(nextPayType);
            form.setFieldValue("type", nextPayType);
        }
    }, [enabledPayOptions, form, payType]);

    const handleSubmit = async (values: any) => {
        const amount = Number(money || 0);
        const min = Number(config?.min || 0);
        const max = Number(config?.max || 0);
        const type = values?.type || selectedPayType;
        if (amount <= 0) {
            message.warning("请输入充值金额");
            return;
        }
        if (min > 0 && amount < min) {
            message.warning(`充值金额不能小于 ${formatMoney(min)} 元`);
            return;
        }
        if (max > 0 && amount > max) {
            message.warning(`充值金额不能大于 ${formatMoney(max)} 元`);
            return;
        }
        if (!type) {
            message.warning("请选择支付方式");
            return;
        }
        const openInCurrentPage = typeof navigator !== "undefined" && /iPhone|iPad|iPod|Macintosh/.test(navigator.userAgent);
        let payWindow: Window | null = null;
        if (!openInCurrentPage && typeof window !== "undefined") payWindow = window.open("about:blank", "_blank");
        const closePayWindow = () => {
            if (payWindow && !payWindow.closed) payWindow.close();
        };
        setSubmitting(true);
        try {
            const createRes = await rechargeBalance({
                body: {money: amount, type},
                onFail: (response: any) => message.error(response?.message || "创建充值订单失败"),
            });
            if (createRes?.code !== 200) {
                closePayWindow();
                return;
            }
            const submitRes = await submitPay({
                body: {id: Number(createRes?.data)},
                onFail: (response: any) => message.error(response?.message || "获取支付链接失败"),
            });
            if (submitRes?.code !== 200) {
                closePayWindow();
                return;
            }
            const payUrl = submitRes?.data;
            if (!payUrl) {
                closePayWindow();
                message.error("获取支付链接失败");
                return;
            }
            setPayJumpUrl(window.location.href);
            modalRef.current?.hide();
            if (payWindow && !payWindow.closed) payWindow.location.href = payUrl;
            else if (openInCurrentPage) {
                window.location.href = payUrl;
                return;
            } else window.open(payUrl, "_blank");
            modal.confirm({
                title: "支付确认",
                content: "支付完成后可刷新余额信息，支付窗口关闭后也可以重新打开支付页面。",
                okText: "我已支付",
                cancelText: "重新支付",
                mask: {closable: true},
                onOk: onRefreshBalance,
                onCancel: () => {
                    if (typeof window !== "undefined" && payUrl) window.open(payUrl, "_blank");
                },
            });
        } finally {
            setSubmitting(false);
        }
    };

    const pageLoading = loading || enumLoading;
    const payDisabled = pageLoading || enabledPayOptions.length <= 0;

    return <ProModal ref={modalRef} title={<Title>在线充值</Title>} width={460}
        modalProps={{footer: null, forceRender: true, mask: {closable: !submitting}, closable: !submitting, keyboard: !submitting} as any}>
        <Spin spinning={pageLoading}>
            <Form className={styles.form} form={form} layout="vertical" requiredMark onFinish={handleSubmit}
                style={{visibility: pageLoading ? "hidden" : "visible"}}>
                <Form.Item label="充值金额" required>
                    <ConfigProvider wave={{disabled: true}}>
                        <Radio.Group className={styles.amounts} value={customAmount ? "custom" : money} disabled={pageLoading || submitting}
                            onChange={(event) => {
                                setCustomAmount(event.target.value === "custom");
                                if (event.target.value !== "custom") setMoney(event.target.value);
                            }}>
                            {quickAmounts.map((amount) => <Radio.Button key={amount} value={amount}
                                disabled={(min > 0 && amount < min) || (max > 0 && amount > max)}>
                                <span className="amount-option"><span className="amount-value">{amount}</span><span className="amount-unit">元</span></span>
                            </Radio.Button>)}
                            <Radio.Button className="custom" value="custom">其他金额</Radio.Button>
                        </Radio.Group>
                    </ConfigProvider>
                    {customAmount && <InputNumber className={styles.customAmount} aria-label="充值金额" value={money || undefined}
                        min={1} precision={0} step={1} controls={false} prefix="¥" placeholder="输入充值金额" autoFocus
                        disabled={pageLoading || submitting} onChange={(value) => setMoney(Number(value || 0))}/>}
                    {amountRange && <div className={styles.hint}>单次充值 {amountRange}</div>}
                </Form.Item>
                <Form.Item label="支付方式" name="type" rules={[{required: true, message: "请选择支付方式"}]}
                    extra={!pageLoading && enabledPayOptions.length === 0 ? "暂未开启可用支付方式" : undefined}>
                    <Radio.Group className={styles.payOptions} disabled={payDisabled || submitting} onChange={(event) => setPayType(event.target.value)}>
                        {enabledPayOptions.map((item) => <Radio key={item.value} value={item.value}>
                            <span className="pay-label"><Icon type={payIcon[item.value]?.icon || "CreditCardOutlined"} style={{color: payIcon[item.value]?.color}}/>{item.label}</span>
                        </Radio>)}
                    </Radio.Group>
                </Form.Item>
                <Form.Item className={styles.checkout} label="应付金额" layout="horizontal" colon={false}
                    labelCol={{flex: "auto"}} wrapperCol={{flex: "none"}}>
                    <span className="checkout-value"><span className="checkout-currency">¥</span>{formatMoney(money)}</span>
                </Form.Item>
                <Button className={styles.submit} type="primary" htmlType="submit" loading={submitting} disabled={payDisabled}>确认支付</Button>
            </Form>
        </Spin>
    </ProModal>;
});
