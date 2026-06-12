import React, {forwardRef, useImperativeHandle, useRef, useState} from "react";
import {App, Form, InputNumber} from "antd";
import {ProModal, ProModalRef, Title} from "sinking-antd";
import {refoundOrder} from "@/service/admin/pay";

export type RefundRef = {
    open: (record: any) => void;
};

type RefundProps = {
    onSuccess?: () => void;
};

const Refund = forwardRef<RefundRef, RefundProps>(({onSuccess}, ref) => {
    const {message} = App.useApp();
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const [form] = Form.useForm();
    const [record, setRecord] = useState<any>();
    const [submitting, setSubmitting] = useState(false);
    const refundableMoney = Math.max(0, Number(record?.pay_amount || 0) - Number(record?.refund_amount || 0));

    useImperativeHandle(ref, () => ({
        open: (data: any) => {
            const nextRecord = data || {};
            setRecord(nextRecord);
            form.resetFields();
            form.setFieldsValue({
                money: Math.max(0, Number(nextRecord?.pay_amount || 0) - Number(nextRecord?.refund_amount || 0)),
            });
            modalRef.current?.show();
        },
    }));

    const submit = async (values: any) => {
        const id = Number(record?.id);
        if (!id) {
            message?.warning("订单ID错误");
            return;
        }
        setSubmitting(true);
        await refoundOrder({
            body: {id, money: Number(values?.money || 0)},
            onSuccess: (res: any) => {
                modalRef.current?.hide();
                onSuccess?.();
                message?.success(res?.message || "退款成功");
            },
            onFail: (res: any) => {
                message?.error(res?.message || "退款失败");
            },
        }).finally(() => {
            setSubmitting(false);
        });
    };

    return (
        <ProModal
            ref={modalRef}
            title={<Title>订单退款</Title>}
            onOk={form.submit}
            width={320}
            modalProps={{
                confirmLoading: submitting,
                okText: "提交",
                cancelText: "取消",
                forceRender: true,
                mask: {
                    closable: true,
                },
                afterClose: () => {
                    setRecord(undefined);
                    setSubmitting(false);
                    form.resetFields();
                },
            } as any}>
            <Form form={form} layout="vertical" onFinish={submit}>
                <Form.Item
                    name="money"
                    label="退款金额"
                    rules={[{required: true, message: "请输入退款金额"}]}>
                    <InputNumber
                        min={0}
                        max={refundableMoney || undefined}
                        precision={2}
                        step={1}
                        style={{width: "100%"}}
                        placeholder="请输入退款金额"
                    />
                </Form.Item>
            </Form>
        </ProModal>
    );
});

Refund.displayName = "Refund";

export default Refund;
