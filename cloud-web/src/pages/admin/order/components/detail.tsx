import React, {forwardRef, useImperativeHandle, useRef, useState} from "react";
import {App, Spin, Tag, Typography} from "antd";
import {ProModal, ProModalRef, Title} from "sinking-antd";
import {createStyles} from "antd-style";
import {getOrderInfo} from "@/service/admin/pay";
import {formatMoney} from "@/utils/string";

export type DetailRef = {
    open: (id: any) => void;
};

type DetailProps = {
    orderEnums: any;
    typeColors?: Record<string, string>;
    statusColors?: Record<string, string>;
    payStatusColors?: Record<string, string>;
};

const useStyles: any = createStyles(({css, token}: any): any => ({
    loading: css`
        min-height: 260px;
        display: flex;
        align-items: center;
        justify-content: center;
    `,
    content: css`
        display: grid;
        gap: 12px;
    `,
    summary: css`
        padding: 14px;
        border: 1px solid ${token.colorBorderSecondary};
        border-radius: ${token.borderRadius}px;
        background: ${token.colorFillQuaternary};

        .top {
            display: flex;
            align-items: flex-start;
            justify-content: space-between;
            gap: 12px;
            min-width: 0;
        }

        .title {
            flex: 1;
            min-width: 0;
        }

        .name {
            overflow: hidden;
            color: ${token.colorText};
            font-size: 15px;
            font-weight: 600;
            line-height: 22px;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .id {
            margin-top: 4px;
            color: ${token.colorTextTertiary};
            font-size: 12px;
            line-height: 18px;
        }

        .tags {
            display: flex;
            align-items: center;
            justify-content: flex-end;
            flex-wrap: wrap;
            gap: 6px;
        }

        .ant-tag {
            margin-inline-end: 0;
        }
    `,
    metrics: css`
        display: grid;
        grid-template-columns: repeat(3, minmax(0, 1fr));
        gap: 8px;

        @media (max-width: 575px) {
            grid-template-columns: repeat(2, minmax(0, 1fr));
        }
    `,
    metric: css`
        min-width: 0;
        padding: 12px;
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
            overflow: hidden;
            color: ${token.colorText};
            font-size: 16px;
            font-weight: 600;
            font-variant-numeric: tabular-nums;
            line-height: 24px;
            text-overflow: ellipsis;
            white-space: nowrap;
        }
    `,
    section: css`
        padding: 12px;
        border: 1px solid ${token.colorBorderSecondary};
        border-radius: ${token.borderRadius}px;
        background: ${token.colorBgContainer};

        .sectionTitle {
            margin-bottom: 10px;
            color: ${token.colorText};
            font-size: 13px;
            font-weight: 600;
            line-height: 20px;
        }

        .grid {
            display: grid;
            grid-template-columns: repeat(2, minmax(0, 1fr));
            gap: 10px 14px;
        }

        .item {
            min-width: 0;
        }

        .label {
            color: ${token.colorTextTertiary};
            font-size: 12px;
            line-height: 18px;
        }

        .value {
            margin-top: 2px;
            min-width: 0;
            color: ${token.colorText};
            font-size: 13px;
            line-height: 20px;
            word-break: break-all;
        }

        @media (max-width: 575px) {
            .grid {
                grid-template-columns: minmax(0, 1fr);
            }
        }
    `,
}));

const enumText = (data: any, value: any) => data?.[String(value)] || data?.[value] || "-";

const InfoItem = ({label, value, copy}: any) => (
    <div className="item">
        <div className="label">{label}</div>
        <div className="value">
            <Typography.Text copyable={copy && value ? {text: String(value)} : false}>
                {value || "-"}
            </Typography.Text>
        </div>
    </div>
);

const Detail = forwardRef<DetailRef, DetailProps>(({orderEnums, typeColors, statusColors, payStatusColors}, ref) => {
    const {message} = App.useApp();
    const {styles} = useStyles();
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const [info, setInfo] = useState<any>();
    const [loading, setLoading] = useState(false);

    useImperativeHandle(ref, () => ({
        open: (id: any) => {
            const orderId = Number(id);
            if (!orderId) {
                message?.warning("订单ID错误");
                return;
            }
            setInfo(undefined);
            setLoading(true);
            modalRef.current?.show();
            void getOrderInfo({
                body: {id: orderId},
                onSuccess: (res: any) => {
                    setInfo(res?.data || {});
                },
                onFail: (res: any) => {
                    message?.error(res?.message || "获取订单详情失败");
                    modalRef.current?.hide();
                },
            })
                .finally(() => {
                    setLoading(false);
                })
                .catch(() => undefined);
        },
    }));

    return (
        <ProModal
            ref={modalRef}
            title={<Title>订单详情</Title>}
            width={560}
            modalProps={{
                footer: null,
                mask: {
                    closable: true,
                },
                afterClose: () => {
                    setInfo(undefined);
                    setLoading(false);
                },
            } as any}>
            {loading ? (
                <div className={styles.loading}>
                    <Spin/>
                </div>
            ) : (
                <div className={styles.content}>
                    <div className={styles.summary}>
                        <div className="top">
                            <div className="title">
                                <div className="name">{info?.name || "订单"}</div>
                                <Typography.Text
                                    className="id"
                                    copyable={info?.id ? {text: String(info.id)} : false}>
                                    ID {info?.id || "-"}
                                </Typography.Text>
                            </div>
                            <div className="tags">
                                <Tag color={typeColors?.[String(info?.type)] || "default"}>{enumText(orderEnums?.type, info?.type)}</Tag>
                                <Tag color={statusColors?.[String(info?.status)] || "default"}>
                                    {enumText(orderEnums?.status, info?.status)}
                                </Tag>
                                <Tag color={payStatusColors?.[String(info?.pay_status)] || "default"}>
                                    {enumText(orderEnums?.pay_status, info?.pay_status)}
                                </Tag>
                            </div>
                        </div>
                    </div>
                    <div className={styles.metrics}>
                        <div className={styles.metric}>
                            <div className="label">订单金额</div>
                            <div className="value">{formatMoney(info?.amount)}</div>
                        </div>
                        <div className={styles.metric}>
                            <div className="label">支付金额</div>
                            <div className="value">{formatMoney(info?.pay_amount)}</div>
                        </div>
                        <div className={styles.metric}>
                            <div className="label">退款金额</div>
                            <div className="value">{formatMoney(info?.refund_amount)}</div>
                        </div>
                    </div>
                    <div className={styles.section}>
                        <div className="sectionTitle">支付信息</div>
                        <div className="grid">
                            <InfoItem label="用户ID" value={info?.user_id} copy/>
                            <InfoItem label="通道ID" value={info?.channel_id} copy/>
                            <InfoItem label="支付方式" value={info?.pay_type}/>
                            <InfoItem label="第三方支付订单号" value={info?.trade_no} copy/>
                        </div>
                    </div>
                    <div className={styles.section}>
                        <div className="sectionTitle">时间信息</div>
                        <div className="grid">
                            <InfoItem label="创建时间" value={info?.create_time}/>
                            <InfoItem label="过期时间" value={info?.expired_time}/>
                            <InfoItem label="支付时间" value={info?.pay_time}/>
                            <InfoItem label="完成时间" value={info?.complete_time}/>
                            <InfoItem label="更新时间" value={info?.update_time}/>
                        </div>
                    </div>
                </div>
            )}
        </ProModal>
    );
});

Detail.displayName = "Detail";

export default Detail;
