import React, {forwardRef, useImperativeHandle, useRef, useState} from "react";
import {App, Form, InputNumber, Tag, Typography} from "antd";
import {ProModal, ProModalRef, Title} from "sinking-antd";
import {createStyles} from "antd-style";
import {deliverOrder, refoundOrder} from "@/service/admin/pay";
import {formatMoney} from "@/utils/string";

export type OperateRef = {
    open: (record: any) => void;
};

type OperateProps = {
    type: "deliver" | "refound";
    orderEnums: any;
    typeColors?: Record<string, string>;
    onSuccess?: () => void;
};

const useStyles: any = createStyles(({css, token}: any): any => ({
    panel: css`
        display: grid;
        gap: 12px;
    `,
    summary: css`
        padding: 14px;
        border: 1px solid ${token.colorBorderSecondary};
        border-radius: ${token.borderRadius}px;
        background: ${token.colorFillQuaternary};

        .name {
            overflow: hidden;
            color: ${token.colorText};
            font-size: 15px;
            font-weight: 600;
            line-height: 22px;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .meta {
            margin-top: 6px;
            display: flex;
            align-items: center;
            flex-wrap: wrap;
            gap: 6px 10px;
            min-width: 0;
        }

        .id {
            color: ${token.colorTextTertiary};
            font-size: 12px;
            line-height: 18px;
        }

        .ant-tag {
            margin-inline-end: 0;
        }
    `,
    amountGrid: css`
        display: grid;
        grid-template-columns: repeat(3, minmax(0, 1fr));
        gap: 8px;

        @media (max-width: 575px) {
            grid-template-columns: repeat(2, minmax(0, 1fr));
        }
    `,
    amountItem: css`
        padding: 11px 12px;
        border: 1px solid ${token.colorBorderSecondary};
        border-radius: ${token.borderRadius}px;
        background: ${token.colorBgContainer};

        .label {
            color: ${token.colorTextTertiary};
            font-size: 12px;
            line-height: 18px;
        }

        .value {
            margin-top: 4px;
            color: ${token.colorText};
            font-size: 15px;
            font-weight: 600;
            font-variant-numeric: tabular-nums;
            line-height: 22px;
        }
    `,
    notice: css`
        padding: 10px 12px;
        border: 1px solid ${token.colorBorderSecondary};
        border-radius: ${token.borderRadius}px;
        background: ${token.colorBgContainer};
        color: ${token.colorTextSecondary};
        font-size: 13px;
        line-height: 20px;
    `,
    form: css`
        padding: 12px;
        border: 1px solid ${token.colorBorderSecondary};
        border-radius: ${token.borderRadius}px;
        background: ${token.colorBgContainer};

        .ant-form-item {
            margin-bottom: 0;
        }
    `,
}));

const enumText = (data: any, value: any) => data?.[String(value)] || data?.[value] || "-";
const numberValue = (value: any) => Number(value || 0);

const Operate = forwardRef<OperateRef, OperateProps>(({type, orderEnums, typeColors, onSuccess}, ref) => {
    const {message} = App.useApp();
    const {styles} = useStyles();
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const [form] = Form.useForm();
    const [record, setRecord] = useState<any>();
    const [submitting, setSubmitting] = useState(false);
    const isRefound = type === "refound";
    const title = isRefound ? "订单退款" : "订单发货";
    const refundableMoney = Math.max(0, numberValue(record?.pay_amount) - numberValue(record?.refund_amount));

    useImperativeHandle(ref, () => ({
        open: (data: any) => {
            const nextRecord = data || {};
            setRecord(nextRecord);
            form.resetFields();
            form.setFieldsValue({
                money: isRefound ? Math.max(0, numberValue(nextRecord?.pay_amount) - numberValue(nextRecord?.refund_amount)) : undefined,
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
        await (isRefound ? refoundOrder : deliverOrder)({
            body: isRefound ? {id, money: Number(values?.money || 0)} : {id},
            onSuccess: (res: any) => {
                modalRef.current?.hide();
                onSuccess?.();
                message?.success(res?.message || "操作成功");
            },
            onFail: (res: any) => {
                message?.error(res?.message || "操作失败");
            },
        }).finally(() => {
            setSubmitting(false);
        });
    };

    return (
        <ProModal
            ref={modalRef}
            title={<Title>{title}</Title>}
            onOk={form.submit}
            width={460}
            modalProps={{
                confirmLoading: submitting,
                okText: isRefound ? "确认退款" : "确认发货",
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
            <div className={styles.panel}>
                <div className={styles.summary}>
                    <div className="name">{record?.name || "订单"}</div>
                    <div className="meta">
                        <Typography.Text
                            className="id"
                            copyable={record?.id ? {text: String(record.id)} : false}>
                            ID {record?.id || "-"}
                        </Typography.Text>
                        <Tag color={typeColors?.[String(record?.type)] || "default"}>{enumText(orderEnums?.type, record?.type)}</Tag>
                    </div>
                </div>
                <div className={styles.amountGrid}>
                    <div className={styles.amountItem}>
                        <div className="label">订单金额</div>
                        <div className="value">{formatMoney(record?.amount)}</div>
                    </div>
                    <div className={styles.amountItem}>
                        <div className="label">支付金额</div>
                        <div className="value">{formatMoney(record?.pay_amount)}</div>
                    </div>
                    <div className={styles.amountItem}>
                        <div className="label">可退金额</div>
                        <div className="value">{formatMoney(refundableMoney)}</div>
                    </div>
                </div>
                <div className={styles.notice}>
                    {isRefound
                        ? "退款成功后会记录退款金额，并返还到用户可用余额。输入 0 时按当前可退款金额处理。"
                        : "发货会触发订单业务回调，业务处理成功后订单状态会更新为已到账。"}
                </div>
                {isRefound && (
                    <Form className={styles.form} form={form} layout="vertical" onFinish={submit}>
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
                )}
                {!isRefound && <Form form={form} onFinish={submit}/>}
            </div>
        </ProModal>
    );
});

Operate.displayName = "Operate";

export default Operate;
