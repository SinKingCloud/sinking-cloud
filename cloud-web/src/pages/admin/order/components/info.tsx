import React, {forwardRef, useEffect, useImperativeHandle, useMemo, useRef, useState} from "react";

import {App, Button, Empty, Pagination, Segmented, Space, Spin, Tag, Tooltip, Typography} from "antd";
import {Icon, ProModal, ProModalRef, Title, useTheme} from "sinking-antd";
import {createStyles} from "antd-style";
import {getOrderInfo, getOrderLogInfo, getOrderLogList} from "@/service/admin/pay";
import {formatMoney} from "@/utils/string";
import PayType from "./pay-type";
import Refund, {type RefundRef} from "./refund";

type ActivePanel = "info" | "log";

type LogDetailRef = {
    open: (id: any) => void;
    reset: () => void;
};

export type InfoRef = {
    open: (record: any, panel?: ActivePanel) => void;
};

type InfoProps = {
    orderEnums: any;
    logEnums: any;
    payTypes?: any;
    channelNames?: any;
    hasDeliverControl?: boolean;
    canDeliver?: (record: any) => boolean;
    onDeliver?: (record: any) => void;
    hasRefundControl?: boolean;
    canRefund?: (record: any) => boolean;
    onRefundSuccess?: () => void;
    typeColors?: Record<string, string>;
    statusColors?: Record<string, string>;
    payStatusColors?: Record<string, string>;
};

const PAGE_SIZE = 8;
const LOG_TYPE_COLORS = ["blue", "green", "orange", "red", "purple"];

const useStyles: any = createStyles(({css, token}: any, props: any = {}): any => {
    const isCompactMode = props?.isCompactMode || false;
    const dotSize = isCompactMode ? 8 : 9;
    const axisTop = isCompactMode ? 13 : 15;

    return {
        modalBody: css`
            display: flex;
            flex-direction: column;
            gap: ${isCompactMode ? 8 : 10}px;
        `,
        segmentWrap: css`
            flex: none;
            display: flex;
            justify-content: center;
            padding: ${isCompactMode ? "1px 0 3px" : "2px 0 4px"};
        `,
        panelContent: css`
            min-height: 0;
            display: flex;
            flex-direction: column;
            gap: ${isCompactMode ? 8 : 10}px;
            padding: 0 2px 2px;
        `,
        loading: css`
            min-height: 300px;
            display: flex;
            align-items: center;
            justify-content: center;
        `,
        summary: css`
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 12px;
            padding: 2px 2px 4px;

            .title {
                flex: 1;
                min-width: 0;
            }

            .titleLine {
                display: flex;
                align-items: center;
                gap: 8px;
                min-width: 0;
            }

            .mark {
                flex: none;
                width: 3px;
                height: 16px;
                border-radius: 99px;
                background: ${token.colorPrimary};
            }

            .name {
                min-width: 0;
                overflow: hidden;
                color: ${token.colorText};
                font-size: 15px;
                font-weight: 600;
                line-height: 22px;
                text-overflow: ellipsis;
                white-space: nowrap;
            }

            .id {
                display: block;
                margin-top: 2px;
                padding-left: 11px;
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 18px;
            }

            .actions {
                flex: none;
                display: flex;
                align-items: center;
                justify-content: flex-end;
                gap: 6px;
            }

            .ant-tag {
                margin-inline-end: 0;
            }

            @media (max-width: 575px) {
                align-items: flex-start;
                flex-direction: column;

                .actions {
                    justify-content: flex-start;
                }
            }
        `,
        metrics: css`
            display: grid;
            grid-template-columns: repeat(3, minmax(0, 1fr));
            gap: 8px;

            @media (max-width: 575px) {
                grid-template-columns: repeat(3, minmax(0, 1fr));
            }
        `,
        metric: css`
            min-width: 0;
            padding: 10px;
            border: 1px solid ${token.colorBorderSecondary};
            border-radius: ${token.borderRadius}px;
            background: ${token.colorBgContainer};

            .label {
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 18px;
            }

            .value {
                margin-top: 2px;
                overflow: hidden;
                font-size: 15px;
                font-variant-numeric: tabular-nums;
                line-height: 22px;
                text-overflow: ellipsis;
                white-space: nowrap;
            }
        `,
        sectionGrid: css`
            display: grid;
            grid-template-columns: repeat(2, minmax(0, 1fr));
            gap: 8px;

            @media (max-width: 575px) {
                grid-template-columns: minmax(0, 1fr);
            }
        `,
        section: css`
            min-width: 0;
            padding: 10px 12px;
            border: 1px solid ${token.colorBorderSecondary};
            border-radius: ${token.borderRadius}px;
            background: ${token.colorBgContainer};

            .sectionTitle {
                margin-bottom: 8px;
                color: ${token.colorText};
                font-size: 13px;
                font-weight: 500;
                line-height: 20px;
            }

            .infoList {
                display: flex;
                flex-direction: column;
                gap: 6px;
            }

            .infoItem {
                display: flex;
                align-items: center;
                gap: 8px;
                min-width: 0;
                min-height: 20px;
            }

            .label {
                flex: none;
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 20px;
            }

            .value {
                display: flex;
                align-items: center;
                flex: 1;
                min-width: 0;
                color: ${token.colorText};
                font-size: 12px;
                line-height: 20px;
                word-break: break-all;
            }

            .value .ant-typography {
                font-size: 12px;
                line-height: 20px;
            }

            .value .ant-tag {
                margin-inline-end: 0;
                line-height: 18px;
            }
        `,
        timeSection: css`
            grid-column: 1 / -1;

            && > .infoList {
                display: grid;
                grid-template-columns: repeat(2, minmax(0, 1fr));
                gap: 6px 16px;
            }

            @media (max-width: 430px) {
                && > .infoList {
                    grid-template-columns: minmax(0, 1fr);
                }
            }
        `,
        logPanel: css`
            display: flex;
            flex-direction: column;
            gap: ${isCompactMode ? 8 : 10}px;
            min-height: 0;
        `,
        timelineShell: css`
        `,
        timeline: css`
            display: flex;
            flex-direction: column;
            padding: 2px 0 0;
        `,
        timelineItem: css`
            position: relative;
            display: grid;
            grid-template-columns: ${isCompactMode ? 20 : 24}px minmax(0, 1fr);
            gap: ${isCompactMode ? 8 : 10}px;
            padding-bottom: ${isCompactMode ? 8 : 10}px;

            &:last-child {
                padding-bottom: 0;

                .line {
                    display: none;
                }
            }

            .axis {
                position: relative;
                display: flex;
                justify-content: center;
                align-items: flex-start;
                min-height: 100%;
                padding-top: ${axisTop}px;
            }

            .dot {
                position: relative;
                z-index: 1;
                width: ${dotSize}px;
                height: ${dotSize}px;
                border: 2px solid ${token.colorPrimary};
                border-radius: 50%;
                background: ${token.colorBgContainer};
            }

            .line {
                position: absolute;
                top: 0;
                bottom: -${isCompactMode ? 10 : 12}px;
                left: 50%;
                width: 1px;
                transform: translateX(-50%);
                background: ${token.colorBorderSecondary};
            }

            .body {
                position: relative;
                min-width: 0;
                padding: ${isCompactMode ? "8px 10px" : "10px 12px"};
                border-radius: ${token.borderRadius}px;
                background: ${token.colorFillQuaternary};
                cursor: pointer;
                transition: background ${token.motionDurationSlow};
            }

            .body:before {
                content: "";
                position: absolute;
                top: ${axisTop - 1}px;
                left: -6px;
                width: 0;
                height: 0;
                border-top: 6px solid transparent;
                border-bottom: 6px solid transparent;
                border-right: 6px solid ${token.colorFillQuaternary};
                transition: border-right-color ${token.motionDurationSlow};
            }

            .body:hover {
                background: ${token.colorFillSecondary};
            }

            .body:hover:before {
                border-right-color: ${token.colorFillSecondary};
            }

            .head {
                display: flex;
                align-items: center;
                justify-content: space-between;
                gap: 8px;
                min-width: 0;
                margin-bottom: ${isCompactMode ? 5 : 6}px;
            }

            .meta {
                display: flex;
                align-items: center;
                flex-wrap: wrap;
                gap: 6px;
                min-width: 0;
            }

            .ant-tag {
                margin-inline-end: 0;
            }

            .time {
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 18px;
            }

            .time {
                flex: none;
                font-variant-numeric: tabular-nums;
                white-space: nowrap;
            }

            .text {
                margin-top: ${isCompactMode ? 5 : 6}px;
                color: ${token.colorTextSecondary};
                font-size: 12px;
                line-height: 18px;
                word-break: break-all;
            }

            .extra {
                display: flex;
                align-items: center;
                justify-content: space-between;
                gap: 8px;
                margin-top: ${isCompactMode ? 5 : 6}px;
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 18px;
            }

            .operator {
                display: flex;
                align-items: center;
                gap: 4px;
                flex: 1;
                min-width: 0;
                overflow: hidden;
                white-space: nowrap;
            }

            .operator-type {
                flex: none;
                line-height: 18px;
                white-space: nowrap;
            }

            .operator-id {
                position: relative;
                display: inline-block;
                flex: none;
                padding-right: 20px;
                height: 18px;
                line-height: 18px;
                white-space: nowrap;
            }

            .operator-id-text {
                color: ${token.colorTextTertiary};
                font-size: 12px;
                height: 18px;
                line-height: 18px;
                white-space: nowrap;
            }

            .operator-copy {
                position: absolute;
                top: 1px;
                right: 0;
                width: 16px;
                height: 16px;
                overflow: hidden;
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 1;
            }

            .operator-copy .ant-typography-copy {
                display: inline-flex;
                align-items: center;
                justify-content: center;
                position: absolute;
                inset: 0;
                width: 16px;
                height: 16px;
                margin-inline-start: 0;
                line-height: 1;
            }

            .operator-copy .ant-typography-copy svg {
                display: block;
            }

            .operator-copy .ant-typography-actions {
                display: block;
                width: 16px;
                height: 16px;
            }

            .detail {
                flex: none;
                padding: 0;
                height: auto;
                font-size: 12px;
            }

            @media (max-width: 575px) {
                grid-template-columns: ${isCompactMode ? 16 : 18}px minmax(0, 1fr);
                gap: 8px;

                .head {
                    align-items: flex-start;
                    flex-direction: column;
                    gap: 4px;
                }

                .extra {
                    align-items: flex-end;
                }

                .time {
                    white-space: normal;
                }
            }
        `,
        empty: css`
            display: flex;
            align-items: center;
            justify-content: center;
            padding: ${isCompactMode ? "24px 0" : "28px 0"};
        `,
        logDetail: css`
            border: 1px solid ${token.colorBorderSecondary};
            border-radius: ${token.borderRadius}px;
            background: ${token.colorBgContainer};

            .placeholder {
                min-height: 240px;
                display: flex;
                align-items: center;
                justify-content: center;
                color: ${token.colorTextTertiary};
                font-size: 12px;
            }

            .head {
                display: flex;
                align-items: flex-start;
                justify-content: space-between;
                gap: 12px;
                padding: 10px 12px;
                border-bottom: 1px solid ${token.colorBorderSecondary};
                background: ${token.colorFillQuaternary};
            }

            .headMain {
                min-width: 0;
            }

            .headMeta {
                display: inline-flex;
                align-items: center;
                flex: none;
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 18px;
                white-space: nowrap;
            }

            .content {
                margin-top: 6px;
                color: ${token.colorText};
                font-size: 13px;
                font-weight: 500;
                line-height: 20px;
                word-break: break-all;
            }

            .info {
                display: grid;
                grid-template-columns: 1.1fr 1fr 1.2fr;
                gap: 8px;
                padding: 10px 12px;
                border-bottom: 1px solid ${token.colorBorderSecondary};
            }

            .infoRow {
                min-width: 0;
                font-size: 12px;
                line-height: 20px;
            }

            .label {
                display: block;
                margin-bottom: 2px;
                color: ${token.colorTextTertiary};
            }

            .value {
                display: block;
                min-width: 0;
                color: ${token.colorText};
                word-break: break-all;
            }

            .value .ant-typography {
                display: inline-flex;
                align-items: center;
                max-width: 100%;
                font-size: 12px;
                line-height: 20px;
                white-space: nowrap;
            }

            .headMeta .ant-typography-copy,
            .value .ant-typography-copy {
                flex: none;
                display: inline-flex;
                align-items: center;
                justify-content: center;
                height: 16px;
                margin-inline-start: 4px;
                line-height: 1;
                vertical-align: top;
            }

            .headMeta .ant-typography-copy svg,
            .value .ant-typography-copy svg {
                display: block;
            }

            .changeHeader,
            .changeRow {
                display: grid;
                grid-template-columns: 120px repeat(2, minmax(0, 1fr));
            }

            .changeHeader {
                border-bottom: 1px solid ${token.colorBorderSecondary};
                background: ${token.colorFillQuaternary};
                color: ${token.colorTextSecondary};
                font-size: 12px;
                font-weight: 500;
                line-height: 18px;
            }

            .changeHeader > div,
            .changeCell {
                min-width: 0;
                padding: 8px 10px;
            }

            .changeRow {
                border-bottom: 1px solid ${token.colorBorderSecondary};
            }

            .changeRow:last-child {
                border-bottom: none;
            }

            .changeRow.changed {
                background: ${token.colorWarningBg};
            }

            .changeKey {
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 18px;
                word-break: break-all;
            }

            .changeValue {
                color: ${token.colorTextSecondary};
                font-size: 12px;
                line-height: 18px;
                white-space: pre-wrap;
                word-break: break-all;
            }

            .changed .changeValue.after {
                color: ${token.colorError};
                font-weight: 500;
            }

            @media (max-width: 700px) {
                .head {
                    flex-direction: column;
                    gap: 6px;
                }

                .info {
                    grid-template-columns: minmax(0, 1fr);
                }

                .changeHeader,
                .changeRow {
                    grid-template-columns: 92px repeat(2, minmax(0, 1fr));
                }
            }
        `,
        pagination: css`
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 8px;
            color: ${token.colorTextTertiary};
            font-size: 12px;
        `,
    };
});

const enumText = (data: any, value: any) => data?.[String(value)] || data?.[value] || "-";
const textValue = (value: any) => (value === undefined || value === null || value === "" ? "-" : String(value));
const isExpiredTime = (value: any) => {
    if (!value) {
        return false;
    }
    const time = new Date(String(value).replace(/-/g, "/")).getTime();
    return Number.isFinite(time) && time > 0 && time < Date.now();
};
const buildTypeColors = (data: any) => Object.keys(data || {}).reduce((colors: Record<string, string>, key, index) => {
    colors[String(key)] = LOG_TYPE_COLORS[index % LOG_TYPE_COLORS.length];
    return colors;
}, {});
const parseJsonValue = (value: any) => {
    if (!value) {
        return undefined;
    }
    if (typeof value === "object") {
        return value;
    }
    try {
        return JSON.parse(String(value));
    } catch (_) {
        return undefined;
    }
};

const CHANGE_FIELD_LABELS: Record<string, string> = {
    id: "订单编号",
    order_id: "订单编号",
    user_id: "用户编号",
    channel_id: "通道编号",
    name: "订单名称",
    type: "订单类型",
    extra: "扩展信息",
    amount: "订单金额",
    status: "订单状态",
    trade_no: "外部单号",
    pay_amount: "实付金额",
    refund_amount: "退款金额",
    pay_status: "支付状态",
    pay_type: "支付方式",
    pay_time: "支付时间",
    expired_time: "过期时间",
    complete_time: "完成时间",
    create_time: "创建时间",
    update_time: "更新时间",
    operator_id: "操作账号",
    operator_type: "操作类型",
    content: "日志内容",
    before_info: "变更前值",
    after_info: "变更后值",
};
const MONEY_CHANGE_FIELDS = ["amount", "pay_amount", "refund_amount"];

const getChangeFieldLabel = (field: string) => CHANGE_FIELD_LABELS[field] || field;

const getChangeEnumData = (field: string, enums: any) => {
    if (field === "status") {
        return enums?.orderEnums?.status;
    }
    if (field === "pay_status") {
        return enums?.orderEnums?.pay_status;
    }
    if (field === "pay_type") {
        return enums?.payTypes;
    }
    if (field === "channel_id") {
        return enums?.channelNames;
    }
    return undefined;
};

const formatChangeValue = (field: string, value: any, enums?: any) => {
    if (value === undefined || value === null || value === "") {
        return "-";
    }
    const enumData = getChangeEnumData(field, enums);
    if (enumData) {
        return enumText(enumData, value);
    }
    if (MONEY_CHANGE_FIELDS.includes(field)) {
        return formatMoney(value);
    }
    if (typeof value === "string") {
        const parsed = parseJsonValue(value);
        if (parsed !== undefined) {
            return JSON.stringify(parsed, null, 2);
        }
        return value;
    }
    if (typeof value === "object") {
        return JSON.stringify(value, null, 2);
    }
    return String(value);
};

const buildChangeRows = (beforeValue: any, afterValue: any, enums?: any) => {
    const before = parseJsonValue(beforeValue) || {};
    const after = parseJsonValue(afterValue) || {};
    const keys = Array.from(new Set([...Object.keys(before), ...Object.keys(after)]));
    return keys.map((key) => {
        const beforeRaw = before?.[key];
        const afterRaw = after?.[key];
        const beforeText = formatChangeValue(key, beforeRaw, enums);
        const afterText = formatChangeValue(key, afterRaw, enums);
        return {
            key,
            beforeText,
            afterText,
            changed: formatChangeValue("", beforeRaw) !== formatChangeValue("", afterRaw),
        };
    });
};

const InfoRow = ({label, value, copy}: any) => {
    const text = textValue(value);
    return (
        <div className="infoItem">
            <div className="label">{label}</div>
            <div className="value">
                {React.isValidElement(value) ? value : (
                    <Typography.Text copyable={copy && text !== "-" ? {text} : false}>{text}</Typography.Text>
                )}
            </div>
        </div>
    );
};

const ExpiredTime = ({value}: any) => {
    const text = textValue(value);
    if (text === "-") {
        return <Typography.Text>{text}</Typography.Text>;
    }
    const expired = isExpiredTime(value);
    if (expired) {
        return (
            <Tooltip title={text}>
                <Tag color="red">已过期</Tag>
            </Tooltip>
        );
    }
    return (
        <Space size={4}>
            <Typography.Text>{text}</Typography.Text>
            <Tag color="green">未过期</Tag>
        </Space>
    );
};

const ChangeCompare = ({before, after, orderEnums, payTypes, channelNames}: any) => {
    const rows = buildChangeRows(before, after, {orderEnums, payTypes, channelNames});
    if (rows.length === 0) {
        return <div className="placeholder">暂无变更信息</div>;
    }
    return (
        <div className="changeCompare">
            <div className="changeHeader">
                <div>字段名称</div>
                <div>变更前</div>
                <div>变更后</div>
            </div>
            {rows.map((item) => (
                <div className={`changeRow ${item.changed ? "changed" : ""}`} key={item.key}>
                    <div className="changeCell changeKey">{getChangeFieldLabel(item.key)}</div>
                    <div className="changeCell changeValue">{item.beforeText}</div>
                    <div className="changeCell changeValue after">{item.afterText}</div>
                </div>
            ))}
        </div>
    );
};

const LogDetail = forwardRef<LogDetailRef, any>(({
                                                     logTypes,
                                                     operatorTypes,
                                                     logTypeColor,
                                                     orderEnums,
                                                     payTypes,
                                                     channelNames,
                                                 }, ref) => {
    const {message} = App.useApp();
    const {styles} = useStyles();
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const [loading, setLoading] = useState(false);
    const [detail, setDetail] = useState<any>();

    useImperativeHandle(ref, () => ({
        open: (id: any) => {
            const logId = Number(id);
            if (!logId) {
                message?.warning("日志ID错误");
                return;
            }
            setDetail(undefined);
            setLoading(true);
            modalRef.current?.show();
            void getOrderLogInfo({
                body: {id: logId},
                onSuccess: (res: any) => setDetail(res?.data || {}),
                onFail: (res: any) => {
                    modalRef.current?.hide();
                    message?.error(res?.message || "获取日志详情失败");
                },
            })
                .finally(() => setLoading(false))
                .catch(() => undefined);
        },
        reset: () => {
            modalRef.current?.hide();
            setLoading(false);
            setDetail(undefined);
        },
    }));

    return (
        <ProModal
            ref={modalRef}
            title={<Title>日志详情</Title>}
            width={720}
            modalProps={{
                footer: null,

                style: {
                    top: 100,
                    paddingBottom: 100,
                },
                mask: {
                    closable: true,
                },
                styles: {
                    body: {
                        padding: "4px 6px 10px",
                    },
                },
                afterClose: () => {
                    setLoading(false);
                    setDetail(undefined);
                },
            } as any}>
            <div className={styles.logDetail}>
                {loading ? (
                    <div className="placeholder">
                        <Spin/>
                    </div>
                ) : detail ? (
                    <>
                        <div className="head">
                            <div className="headMain">
                                <Space size={6} wrap>
                                    <Tag color={logTypeColor(detail?.type)}>{enumText(logTypes, detail?.type)}</Tag>
                                </Space>
                                <div className="content">{detail?.content || "-"}</div>
                            </div>
                            <Typography.Text className="headMeta"
                                             copyable={detail?.id ? {text: String(detail?.id)} : false}>
                                日志ID {detail?.id || "-"}
                            </Typography.Text>
                        </div>
                        <div className="info">
                            <div className="infoRow">
                                <span className="label">订单ID</span>
                                <span className="value">
                                    <Typography.Text
                                        copyable={detail?.order_id ? {text: String(detail?.order_id)} : false}>
                                        {detail?.order_id || "-"}
                                    </Typography.Text>
                                </span>
                            </div>
                            <div className="infoRow">
                                <span className="label">操作人</span>
                                <span className="value">
                                    {enumText(operatorTypes, detail?.operator_type)}
                                    {detail?.operator_id ? (
                                        <>
                                            {" "}
                                            <Typography.Text copyable={{text: String(detail?.operator_id)}}>
                                                账号ID：{detail?.operator_id}
                                            </Typography.Text>
                                        </>
                                    ) : null}
                                </span>
                            </div>
                            <div className="infoRow">
                                <span className="label">创建时间</span>
                                <span className="value">{detail?.create_time || "-"}</span>
                            </div>
                        </div>
                        <ChangeCompare
                            before={detail?.before_info}
                            after={detail?.after_info}
                            orderEnums={orderEnums}
                            payTypes={payTypes}
                            channelNames={channelNames}
                        />
                    </>
                ) : (
                    <div className="placeholder">暂无日志详情</div>
                )}
            </div>
        </ProModal>
    );
});

LogDetail.displayName = "LogDetail";

const Info = forwardRef<InfoRef, InfoProps>(({
                                                 orderEnums,
                                                 logEnums,
                                                 payTypes,
                                                 channelNames,
                                                 hasDeliverControl,
                                                 canDeliver,
                                                 onDeliver,
                                                 hasRefundControl,
                                                 canRefund,
                                                 onRefundSuccess,
                                                 typeColors,
                                                 statusColors,
                                                 payStatusColors,
                                             }, ref) => {
    const {message} = App.useApp();
    const theme = useTheme();
    const isCompactMode = theme?.isCompactTheme?.() || false;
    const {styles} = useStyles(useMemo(() => ({isCompactMode}), [isCompactMode]));
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const refundRef = useRef<RefundRef>({} as RefundRef);
    const logDetailRef = useRef<LogDetailRef>({} as LogDetailRef);
    const [activePanel, setActivePanel] = useState<ActivePanel>("info");
    const [orderInfo, setOrderInfo] = useState<any>();
    const [info, setInfo] = useState<any>();
    const [infoLoading, setInfoLoading] = useState(false);
    const [page, setPage] = useState(1);
    const [total, setTotal] = useState(0);
    const [logs, setLogs] = useState<any[]>([]);
    const [logLoading, setLogLoading] = useState(false);
    const logTypes = logEnums?.type || {};
    const operatorTypes = logEnums?.operator_type || {};
    const logTypeColors = buildTypeColors(logTypes);
    const logTypeColor = (value: any) => logTypeColors?.[String(value)] || "default";
    const orderId = Number(orderInfo?.id);

    const loadInfo = (id: number) => {
        if (!id) {
            return;
        }
        setInfoLoading(true);
        void getOrderInfo({
            body: {id},
            onSuccess: (res: any) => {
                const nextInfo = res?.data || {};
                setInfo(nextInfo);
                setOrderInfo((prev: any) => ({...(prev || {}), ...nextInfo}));
            },
            onFail: (res: any) => {
                modalRef.current?.hide();
                message?.error(res?.message || "获取订单详情失败");
            },
        })
            .finally(() => {
                setInfoLoading(false);
            })
            .catch(() => undefined);
    };

    const loadLogs = () => {
        if (!orderId) {
            return;
        }
        setLogLoading(true);
        const body: any = {
            order_id: orderId,
            page,
            page_size: PAGE_SIZE,
            order_by_field: "id",
            order_by_type: "desc",
        };
        void getOrderLogList({
            body,
            onSuccess: (res: any) => {
                setLogs(Array.isArray(res?.data?.list) ? res?.data?.list?.filter(Boolean) : []);
                setTotal(Number(res?.data?.total || 0));
            },
            onFail: (res: any) => message?.error(res?.message || "获取订单日志失败"),
        })
            .finally(() => {
                setLogLoading(false);
            })
            .catch(() => undefined);
    };

    useImperativeHandle(ref, () => ({
        open: (record: any, panel: ActivePanel = "info") => {
            const id = Number(record?.id || record);
            if (!id) {
                message?.warning("订单ID错误");
                return;
            }
            setActivePanel(panel);
            setOrderInfo(typeof record === "object" ? record : {id});
            setInfo(undefined);
            setPage(1);
            setTotal(0);
            setLogs([]);
            setLogLoading(panel === "log");
            modalRef.current?.show();
            loadInfo(id);
        },
    }));

    useEffect(() => {
        if (activePanel !== "log" || !orderId) {
            return;
        }
        loadLogs();
    }, [activePanel, orderId, page]);

    const renderInfo = () => {
        if (infoLoading) {
            return (
                <div className={styles.loading}>
                    <Spin/>
                </div>
            );
        }
        const data = info || orderInfo || {};
        return (
            <>
                <div className={styles.summary}>
                    <div className="title">
                        <div className="titleLine">
                            <span className="mark"/>
                            <div className="name">{data?.name || "订单"}</div>
                        </div>
                        <Typography.Text
                            className="id"
                            copyable={data?.id ? {text: String(data?.id)} : false}>
                            订单ID {data?.id || "-"}
                        </Typography.Text>
                    </div>
                    {(hasDeliverControl || hasRefundControl) && (
                        <div className="actions">
                            {hasDeliverControl && (
                                <Button
                                    type="primary"
                                    size="small"
                                    icon={<Icon type="TruckOutlined"/>}
                                    disabled={!canDeliver?.(data)}
                                    onClick={() => onDeliver?.(data)}>
                                    发货
                                </Button>
                            )}
                            {hasRefundControl && (
                                <Button
                                    danger
                                    size="small"
                                    icon={<Icon type="RollbackOutlined"/>}
                                    disabled={!canRefund?.(data)}
                                    onClick={() => refundRef.current?.open(data)}>
                                    退款
                                </Button>
                            )}
                        </div>
                    )}
                </div>
                <div className={styles.metrics}>
                    <div className={styles.metric}>
                        <div className="label">订单金额</div>
                        <div className="value" style={{color: "#1677ff"}}>{formatMoney(data?.amount)}</div>
                    </div>
                    <div className={styles.metric}>
                        <div className="label">实付金额</div>
                        <div className="value" style={{color: "#52c41a"}}>{formatMoney(data?.pay_amount)}</div>
                    </div>
                    <div className={styles.metric}>
                        <div className="label">退款金额</div>
                        <div className="value" style={{color: "#ff4d4f"}}>{formatMoney(data?.refund_amount)}</div>
                    </div>
                </div>
                <div className={styles.sectionGrid}>
                    <div className={styles.section}>
                        <div className="sectionTitle">订单信息</div>
                        <div className="infoList">
                            <InfoRow label="订单ID" value={data?.id} copy/>
                            <InfoRow label="外部ID" value={data?.trade_no} copy/>
                            <InfoRow label="所属用户" value={data?.user_id} copy/>
                            <InfoRow
                                label="订单类型"
                                value={<Tag
                                    color={typeColors?.[String(data?.type)] || "default"}>{enumText(orderEnums?.type, data?.type)}</Tag>}
                            />
                            <InfoRow
                                label="订单状态"
                                value={<Tag
                                    color={statusColors?.[String(data?.status)] || "default"}>{enumText(orderEnums?.status, data?.status)}</Tag>}
                            />
                        </div>
                    </div>
                    <div className={styles.section}>
                        <div className="sectionTitle">支付信息</div>
                        <div className="infoList">
                            <InfoRow label="支付方式" value={<PayType data={payTypes} value={data?.pay_type}/>}/>
                            <InfoRow label="通道名称" value={enumText(channelNames, data?.channel_id)}/>
                            <InfoRow label="通道ID" value={data?.channel_id} copy/>
                            <InfoRow
                                label="支付状态"
                                value={<Tag
                                    color={payStatusColors?.[String(data?.pay_status)] || "default"}>{enumText(orderEnums?.pay_status, data?.pay_status)}</Tag>}
                            />
                            <InfoRow label="过期时间" value={<ExpiredTime value={data?.expired_time}/>}/>
                        </div>
                    </div>
                    <div className={`${styles.section} ${styles.timeSection}`}>
                        <div className="sectionTitle">时间信息</div>
                        <div className="infoList">
                            <InfoRow label="创建时间" value={data?.create_time}/>
                            <InfoRow label="支付时间" value={data?.pay_time}/>
                            <InfoRow label="完成时间" value={data?.complete_time}/>
                            <InfoRow label="更新时间" value={data?.update_time}/>
                        </div>
                    </div>
                </div>
            </>
        );
    };

    const renderLog = () => {
        if (logLoading) {
            return (
                <div className={styles.loading}>
                    <Spin/>
                </div>
            );
        }
        return (
            <div className={styles.logPanel}>
                <div className={styles.timelineShell}>
                    {logs.length > 0 ? (
                        <div className={styles.timeline}>
                            {logs.map((item: any) => (
                                <div className={styles.timelineItem} key={item?.id}>
                                    <div className="axis">
                                        <span className="dot"/>
                                        <span className="line"/>
                                    </div>
                                    <div className="body" onClick={() => logDetailRef.current?.open(item?.id)}>
                                        <div className="head">
                                            <div className="meta">
                                                <Tag color={logTypeColor(item?.type)}>
                                                    {enumText(logTypes, item?.type)}
                                                </Tag>
                                            </div>
                                            <Typography.Text
                                                className="time">
                                                {item?.create_time || "-"}
                                            </Typography.Text>
                                        </div>
                                        <div className="text">{item?.content || "-"}</div>
                                        <div className="extra">
                                            <span className="operator">
                                                <span className="operator-type">操作人：{enumText(operatorTypes, item?.operator_type)}</span>
                                                {item?.operator_id ? (
                                                    <span className="operator-id" onClick={(event) => event.stopPropagation()}>
                                                        <span className="operator-id-text">账号ID：{item?.operator_id}</span>
                                                        <Typography.Text
                                                            className="operator-copy"
                                                            copyable={{text: String(item?.operator_id), tooltips: ["复制", "复制成功"]}}
                                                        />
                                                    </span>
                                                ) : null}
                                            </span>
                                            <Button
                                                className="detail"
                                                size="small"
                                                type="link"
                                                onClick={(event) => {
                                                    event.stopPropagation();
                                                    logDetailRef.current?.open(item?.id);
                                                }}>
                                                查看详情
                                            </Button>
                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>
                    ) : (
                        <div className={styles.empty}>
                            <Empty image={Empty.PRESENTED_IMAGE_SIMPLE}/>
                        </div>
                    )}
                </div>
                {total > 0 && (
                    <div className={styles.pagination}>
                        <span>共 {total} 条</span>
                        <Pagination
                            size="small"
                            current={page}
                            pageSize={PAGE_SIZE}
                            total={total}
                            showSizeChanger={false}
                            onChange={(nextPage) => {
                                setPage(nextPage);
                            }}
                        />
                    </div>
                )}
            </div>
        );
    };

    return (
        <ProModal
            ref={modalRef}
            title={<Title>订单信息</Title>}
            width={600}
            modalProps={{
                footer: null,

                style: {
                    top: 100,
                    paddingBottom: 100,
                },
                mask: {
                    closable: true,
                },
                styles: {
                    body: {
                        padding: "2px 6px 10px",
                    },
                },
                afterClose: () => {
                    setActivePanel("info");
                    setOrderInfo(undefined);
                    setInfo(undefined);
                    setInfoLoading(false);
                    setPage(1);
                    setTotal(0);
                    setLogs([]);
                    setLogLoading(false);
                    logDetailRef.current?.reset();
                },
            } as any}>
            <div className={styles.modalBody}>
                <div className={styles.segmentWrap}>
                    <Segmented
                        value={activePanel}
                        options={[
                            {label: "详情", value: "info"},
                            {label: "日志", value: "log"},
                        ]}
                        onChange={(value) => setActivePanel(value as ActivePanel)}
                    />
                </div>
                <div className={styles.panelContent}>
                    {activePanel === "log" ? renderLog() : renderInfo()}
                </div>
                <Refund
                    ref={refundRef}
                    onSuccess={onRefundSuccess}
                />
                <LogDetail
                    ref={logDetailRef}
                    logTypes={logTypes}
                    operatorTypes={operatorTypes}
                    logTypeColor={logTypeColor}
                    orderEnums={orderEnums}
                    payTypes={payTypes}
                    channelNames={channelNames}
                />
            </div>
        </ProModal>
    );
});

Info.displayName = "Info";

export default Info;
