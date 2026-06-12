import React, {forwardRef, useEffect, useImperativeHandle, useRef, useState} from "react";
import {App, Button, DatePicker, Empty, Pagination, Segmented, Select, Space, Spin, Tag, Typography} from "antd";
import {RollbackOutlined, TruckOutlined} from "@ant-design/icons";
import {ProModal, ProModalRef, Title} from "sinking-antd";
import {createStyles} from "antd-style";
import {getOrderInfo, getOrderLogInfo, getOrderLogList} from "@/service/admin/pay";
import {formatMoney} from "@/utils/string";
import PayType from "./pay-type";
import Refund, {type RefundRef} from "./refund";

type ActivePanel = "info" | "log";

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

const ALL_VALUE = "__all__";
const PAGE_SIZE = 8;
const LOG_TYPE_COLORS = ["blue", "green", "orange", "red", "purple"];

const useStyles: any = createStyles(({css, token}: any): any => ({
    modalBody: css`
        display: flex;
        flex-direction: column;
        gap: 10px;
        min-height: 360px;
        max-height: calc(100dvh - 160px);
        overflow: hidden;
    `,
    segmentWrap: css`
        flex: none;
        display: flex;
        justify-content: center;
        padding: 2px 0 4px;
    `,
    panelContent: css`
        flex: 1;
        min-height: 0;
        overflow: auto;
        display: flex;
        flex-direction: column;
        gap: 10px;
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

        .infoList {
            display: grid;
            grid-template-columns: repeat(2, minmax(0, 1fr));
            gap: 6px 16px;
        }

        @media (max-width: 430px) {
            .infoList {
                grid-template-columns: minmax(0, 1fr);
            }
        }
    `,
    logPanel: css`
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 8px;
        min-height: 0;
    `,
    toolbar: css`
        display: grid;
        grid-template-columns: 128px 128px minmax(0, 1fr) auto;
        gap: 8px;

        .dateRange {
            width: 100%;
        }

        @media (max-width: 575px) {
            grid-template-columns: repeat(2, minmax(0, 1fr));

            .dateRange {
                grid-column: 1 / -1;
            }
        }
    `,
    logContent: css`
        flex: 1;
        display: grid;
        grid-template-columns: minmax(0, 1fr) 260px;
        gap: 8px;
        min-height: 0;

        @media (max-width: 700px) {
            grid-template-columns: minmax(0, 1fr);
        }
    `,
    timelineShell: css`
        min-height: 260px;
        overflow: auto;
        border: 1px solid ${token.colorBorderSecondary};
        border-radius: ${token.borderRadius}px;
        background: ${token.colorBgContainer};
    `,
    timeline: css`
        display: flex;
        flex-direction: column;
        padding: 4px 0;
    `,
    timelineItem: css`
        position: relative;
        display: grid;
        grid-template-columns: 18px minmax(0, 1fr);
        gap: 8px;
        padding: 7px 10px;
        cursor: pointer;

        &:last-child {
            .line {
                display: none;
            }
        }

        .axis {
            position: relative;
            display: flex;
            justify-content: center;
            padding-top: 7px;
        }

        .dot {
            position: relative;
            z-index: 1;
            width: 8px;
            height: 8px;
            border: 2px solid ${token.colorPrimary};
            border-radius: 50%;
            background: ${token.colorBgContainer};
        }

        .line {
            position: absolute;
            top: 18px;
            bottom: -8px;
            left: 50%;
            width: 1px;
            transform: translateX(-50%);
            background: ${token.colorBorderSecondary};
        }

        .body {
            min-width: 0;
            padding: 7px 9px;
            border-radius: ${token.borderRadius}px;
            transition: background ${token.motionDurationSlow};
        }

        &.active .body,
        .body:hover {
            background: ${token.colorFillQuaternary};
        }

        .head {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 8px;
            min-width: 0;
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

        .logId,
        .time,
        .operator {
            color: ${token.colorTextTertiary};
            font-size: 12px;
            line-height: 18px;
        }

        .time {
            font-variant-numeric: tabular-nums;
            white-space: nowrap;
        }

        .text {
            margin-top: 4px;
            color: ${token.colorText};
            font-size: 13px;
            line-height: 20px;
            word-break: break-all;
        }

        .operator {
            margin-top: 2px;
        }
    `,
    empty: css`
        min-height: 240px;
        display: flex;
        align-items: center;
        justify-content: center;
    `,
    logDetail: css`
        min-height: 260px;
        overflow: auto;
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
            padding: 10px;
            border-bottom: 1px solid ${token.colorBorderSecondary};
            background: ${token.colorFillQuaternary};
        }

        .content {
            margin-top: 8px;
            color: ${token.colorText};
            font-size: 13px;
            line-height: 20px;
            word-break: break-all;
        }

        .info {
            padding: 10px;
            border-bottom: 1px solid ${token.colorBorderSecondary};
        }

        .infoRow {
            display: flex;
            gap: 8px;
            font-size: 12px;
            line-height: 20px;
        }

        .label {
            flex: none;
            color: ${token.colorTextTertiary};
        }

        .value {
            min-width: 0;
            word-break: break-all;
        }

        .json {
            padding: 10px;
            border-bottom: 1px solid ${token.colorBorderSecondary};
        }

        .json:last-child {
            border-bottom: none;
        }

        .jsonTitle {
            margin-bottom: 6px;
            color: ${token.colorText};
            font-size: 13px;
            font-weight: 500;
            line-height: 20px;
        }

        .jsonText {
            max-height: 170px;
            margin-bottom: 0 !important;
            overflow: auto;
            color: ${token.colorTextSecondary};
            font-size: 12px;
            line-height: 18px;
            white-space: pre-wrap;
            word-break: break-all;
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
}));

const enumText = (data: any, value: any) => data?.[String(value)] || data?.[value] || "-";
const textValue = (value: any) => (value === undefined || value === null || value === "" ? "-" : String(value));
const isExpiredTime = (value: any) => {
    if (!value) {
        return false;
    }
    const time = new Date(String(value).replace(/-/g, "/")).getTime();
    return Number.isFinite(time) && time > 0 && time < Date.now();
};
const buildOptions = (data: any, allLabel = "全部") => [
    {label: allLabel, value: ALL_VALUE},
    ...Object.entries(data || {}).map(([value, label]) => ({value, label: String(label)})),
];
const buildTypeColors = (data: any) => Object.keys(data || {}).reduce((colors: Record<string, string>, key, index) => {
    colors[String(key)] = LOG_TYPE_COLORS[index % LOG_TYPE_COLORS.length];
    return colors;
}, {});
const formatJson = (value: any) => {
    if (!value) {
        return "-";
    }
    try {
        return JSON.stringify(JSON.parse(String(value)), null, 2);
    } catch (_) {
        return String(value);
    }
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
    return (
        <Space size={4}>
            <Typography.Text>{text}</Typography.Text>
            <Tag color={expired ? "red" : "green"}>{expired ? "已过期" : "未过期"}</Tag>
        </Space>
    );
};

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
    const {styles} = useStyles();
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const refundRef = useRef<RefundRef>({} as RefundRef);
    const [activePanel, setActivePanel] = useState<ActivePanel>("info");
    const [orderInfo, setOrderInfo] = useState<any>();
    const [info, setInfo] = useState<any>();
    const [infoLoading, setInfoLoading] = useState(false);
    const [logType, setLogType] = useState(ALL_VALUE);
    const [operatorType, setOperatorType] = useState(ALL_VALUE);
    const [range, setRange] = useState<any>();
    const [page, setPage] = useState(1);
    const [total, setTotal] = useState(0);
    const [logs, setLogs] = useState<any[]>([]);
    const [logLoading, setLogLoading] = useState(false);
    const [activeLogId, setActiveLogId] = useState<any>();
    const [logDetail, setLogDetail] = useState<any>();
    const [logDetailLoading, setLogDetailLoading] = useState(false);
    const logTypes = logEnums?.type || {};
    const operatorTypes = logEnums?.operator_type || {};
    const logTypeOptions = buildOptions(logTypes, "全部类型");
    const operatorOptions = buildOptions(operatorTypes, "全部操作人");
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
                message?.error(res?.message || "获取订单详情失败");
                modalRef.current?.hide();
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
        if (logType !== ALL_VALUE) {
            body.type = logType;
        }
        if (operatorType !== ALL_VALUE) {
            body.operator_type = operatorType;
        }
        if (range?.[0] && range?.[1]) {
            body.create_time_start = range[0]?.startOf
                ? range[0].startOf("day").format("YYYY-MM-DD HH:mm:ss")
                : range[0];
            body.create_time_end = range[1]?.endOf
                ? range[1].endOf("day").format("YYYY-MM-DD HH:mm:ss")
                : range[1];
        }
        void getOrderLogList({
            body,
            onSuccess: (res: any) => {
                const nextList = res?.data?.list || [];
                setLogs(nextList);
                setTotal(Number(res?.data?.total || 0));
                if (nextList.length > 0) {
                    loadLogDetail(nextList[0]?.id);
                } else {
                    setActiveLogId(undefined);
                    setLogDetail(undefined);
                }
            },
            onFail: (res: any) => message?.error(res?.message || "获取订单日志失败"),
        })
            .finally(() => {
                setLogLoading(false);
            })
            .catch(() => undefined);
    };

    const loadLogDetail = (id: any) => {
        const logId = Number(id);
        if (!logId) {
            message?.warning("日志ID错误");
            return;
        }
        setActiveLogId(logId);
        setLogDetail(undefined);
        setLogDetailLoading(true);
        void getOrderLogInfo({
            body: {id: logId},
            onSuccess: (res: any) => {
                setLogDetail(res?.data || {});
            },
            onFail: (res: any) => message?.error(res?.message || "获取日志详情失败"),
        })
            .finally(() => {
                setLogDetailLoading(false);
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
            setLogType(ALL_VALUE);
            setOperatorType(ALL_VALUE);
            setRange(undefined);
            setPage(1);
            setTotal(0);
            setLogs([]);
            setActiveLogId(undefined);
            setLogDetail(undefined);
            modalRef.current?.show();
            loadInfo(id);
        },
    }));

    useEffect(() => {
        if (activePanel !== "log" || !orderId) {
            return;
        }
        loadLogs();
    }, [activePanel, orderId, logType, operatorType, range, page]);

    const resetLogFilter = () => {
        setLogType(ALL_VALUE);
        setOperatorType(ALL_VALUE);
        setRange(undefined);
        setPage(1);
        setActiveLogId(undefined);
        setLogDetail(undefined);
    };

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
                            copyable={data?.id ? {text: String(data.id)} : false}>
                            订单ID {data?.id || "-"}
                        </Typography.Text>
                    </div>
                    {(hasDeliverControl || hasRefundControl) && (
                        <div className="actions">
                            {hasDeliverControl && (
                                <Button
                                    type="primary"
                                    size="small"
                                    icon={<TruckOutlined/>}
                                    disabled={!canDeliver?.(data)}
                                    onClick={() => onDeliver?.(data)}>
                                    发货
                                </Button>
                            )}
                            {hasRefundControl && (
                                <Button
                                    danger
                                    size="small"
                                    icon={<RollbackOutlined/>}
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
                                value={<Tag color={typeColors?.[String(data?.type)] || "default"}>{enumText(orderEnums?.type, data?.type)}</Tag>}
                            />
                            <InfoRow
                                label="订单状态"
                                value={<Tag color={statusColors?.[String(data?.status)] || "default"}>{enumText(orderEnums?.status, data?.status)}</Tag>}
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
                                value={<Tag color={payStatusColors?.[String(data?.pay_status)] || "default"}>{enumText(orderEnums?.pay_status, data?.pay_status)}</Tag>}
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

    const renderJson = (title: string, value: any) => {
        const text = formatJson(value);
        return (
            <div className="json">
                <div className="jsonTitle">{title}</div>
                <Typography.Paragraph
                    className="jsonText"
                    copyable={text !== "-" ? {text} : false}>
                    {text}
                </Typography.Paragraph>
            </div>
        );
    };

    const renderLogDetail = () => {
        if (logDetailLoading) {
            return (
                <div className="placeholder">
                    <Spin/>
                </div>
            );
        }
        if (!logDetail) {
            return <div className="placeholder">选择左侧日志查看详情</div>;
        }
        return (
            <>
                <div className="head">
                    <Space size={6} wrap>
                        <Tag color={logTypeColor(logDetail?.type)}>{enumText(logTypes, logDetail?.type)}</Tag>
                        <Typography.Text type="secondary">#{logDetail?.id || "-"}</Typography.Text>
                    </Space>
                    <div className="content">{logDetail?.content || "-"}</div>
                </div>
                <div className="info">
                    <div className="infoRow">
                        <span className="label">操作人</span>
                        <span className="value">
                            {enumText(operatorTypes, logDetail?.operator_type)} #{logDetail?.operator_id || 0}
                        </span>
                    </div>
                    <div className="infoRow">
                        <span className="label">创建时间</span>
                        <span className="value">{logDetail?.create_time || "-"}</span>
                    </div>
                </div>
                {renderJson("变更前", logDetail?.before_info)}
                {renderJson("变更后", logDetail?.after_info)}
            </>
        );
    };

    const renderLog = () => (
        <div className={styles.logPanel}>
            <div className={styles.summary}>
                <div className="title">
                    <div className="name">{orderInfo?.name || info?.name || "订单日志"}</div>
                    <Typography.Text
                        className="id"
                        copyable={orderId ? {text: String(orderId)} : false}>
                        订单ID {orderId || "-"}
                    </Typography.Text>
                </div>
            </div>
            <div className={styles.toolbar}>
                <Select
                    value={logType}
                    options={logTypeOptions}
                    onChange={(value) => {
                        setLogType(value);
                        setPage(1);
                        setActiveLogId(undefined);
                        setLogDetail(undefined);
                    }}
                />
                <Select
                    value={operatorType}
                    options={operatorOptions}
                    onChange={(value) => {
                        setOperatorType(value);
                        setPage(1);
                        setActiveLogId(undefined);
                        setLogDetail(undefined);
                    }}
                />
                <DatePicker.RangePicker
                    className="dateRange"
                    value={range}
                    onChange={(value) => {
                        setRange(value);
                        setPage(1);
                        setActiveLogId(undefined);
                        setLogDetail(undefined);
                    }}
                />
                <Button onClick={resetLogFilter}>重置</Button>
            </div>
            <Spin spinning={logLoading}>
                <div className={styles.logContent}>
                    <div className={styles.timelineShell}>
                        {logs.length > 0 ? (
                            <div className={styles.timeline}>
                                {logs.map((item: any) => (
                                    <div
                                        className={`${styles.timelineItem} ${Number(activeLogId) === Number(item?.id) ? "active" : ""}`}
                                        key={item?.id}
                                        onClick={() => loadLogDetail(item?.id)}>
                                        <div className="axis">
                                            <span className="dot"/>
                                            <span className="line"/>
                                        </div>
                                        <div className="body">
                                            <div className="head">
                                                <div className="meta">
                                                    <Tag color={logTypeColor(item?.type)}>
                                                        {enumText(logTypes, item?.type)}
                                                    </Tag>
                                                    <Typography.Text className="logId">#{item?.id || "-"}</Typography.Text>
                                                    <Typography.Text className="time">{item?.create_time || "-"}</Typography.Text>
                                                </div>
                                            </div>
                                            <div className="text">{item?.content || "-"}</div>
                                            <div className="operator">
                                                {enumText(operatorTypes, item?.operator_type)}
                                                {item?.operator_id ? ` #${item.operator_id}` : ""}
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
                    <div className={styles.logDetail}>
                        {renderLogDetail()}
                    </div>
                </div>
            </Spin>
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
                            setActiveLogId(undefined);
                            setLogDetail(undefined);
                        }}
                    />
                </div>
            )}
        </div>
    );

    return (
        <ProModal
            ref={modalRef}
            title={<Title>订单信息</Title>}
            width={600}
            modalProps={{
                footer: null,
                forceRender: true,
                mask: {
                    closable: true,
                },
                styles: {
                    body: {
                        maxHeight: "calc(100dvh - 150px)",
                        overflow: "hidden",
                        padding: "2px 6px 10px",
                    },
                },
                afterClose: () => {
                    setActivePanel("info");
                    setOrderInfo(undefined);
                    setInfo(undefined);
                    setInfoLoading(false);
                    setLogType(ALL_VALUE);
                    setOperatorType(ALL_VALUE);
                    setRange(undefined);
                    setPage(1);
                    setTotal(0);
                    setLogs([]);
                    setLogLoading(false);
                    setActiveLogId(undefined);
                    setLogDetail(undefined);
                    setLogDetailLoading(false);
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
            </div>
        </ProModal>
    );
});

Info.displayName = "Info";

export default Info;
