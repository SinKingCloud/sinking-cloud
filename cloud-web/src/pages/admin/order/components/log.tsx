import React, {forwardRef, useEffect, useImperativeHandle, useRef, useState} from "react";
import {App, Button, DatePicker, Empty, Pagination, Select, Space, Spin, Tag, Typography} from "antd";
import {ProModal, ProModalRef, Title} from "sinking-antd";
import {createStyles} from "antd-style";
import {getOrderLogInfo, getOrderLogList} from "@/service/admin/pay";

export type LogRef = {
    open: (record: any) => void;
};

type LogProps = {
    logEnums: any;
};

const ALL_VALUE = "__all__";
const PAGE_SIZE = 8;
const LOG_TYPE_COLORS = ["blue", "green", "orange", "red", "purple"];

const useStyles: any = createStyles(({css, token}: any): any => ({
    panel: css`
        display: grid;
        gap: 12px;
    `,
    orderHeader: css`
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
    `,
    filter: css`
        display: flex;
        align-items: center;
        gap: 8px;
        flex-wrap: wrap;

        .filterSelect {
            width: 128px;
        }

        .filterDate {
            width: 240px;
        }

        @media (max-width: 575px) {
            display: grid;
            grid-template-columns: repeat(2, minmax(0, 1fr));

            .filterSelect,
            .filterDate,
            .reset {
                width: 100%;
            }

            .filterDate {
                grid-column: 1 / -1;
            }
        }
    `,
    listWrap: css`
        min-height: 260px;
    `,
    list: css`
        display: grid;
        gap: 8px;
    `,
    item: css`
        position: relative;
        overflow: hidden;
        padding: 12px 12px 12px 16px;
        border: 1px solid ${token.colorBorderSecondary};
        border-radius: ${token.borderRadius}px;
        background: ${token.colorBgContainer};

        &:before {
            content: "";
            position: absolute;
            top: 0;
            left: 0;
            bottom: 0;
            width: 3px;
            background: ${token.colorPrimary};
        }

        .head {
            display: flex;
            align-items: flex-start;
            justify-content: space-between;
            gap: 10px;
            min-width: 0;
        }

        .main {
            flex: 1;
            min-width: 0;
        }

        .tags {
            display: flex;
            align-items: center;
            flex-wrap: wrap;
            gap: 6px;
            min-width: 0;
        }

        .ant-tag {
            margin-inline-end: 0;
        }

        .content {
            margin-top: 7px;
            color: ${token.colorText};
            font-size: 13px;
            line-height: 20px;
            word-break: break-all;
        }

        .meta {
            margin-top: 8px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 8px;
            color: ${token.colorTextTertiary};
            font-size: 12px;
            line-height: 18px;
        }

        .operator {
            min-width: 0;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .time {
            flex: none;
            font-variant-numeric: tabular-nums;
        }
    `,
    empty: css`
        min-height: 220px;
        display: flex;
        align-items: center;
        justify-content: center;
        border: 1px dashed ${token.colorBorderSecondary};
        border-radius: ${token.borderRadius}px;
        background: ${token.colorBgContainer};
    `,
    pagination: css`
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 10px;
        color: ${token.colorTextTertiary};
        font-size: 12px;

        @media (max-width: 575px) {
            .ant-pagination {
                flex: none;
            }
        }
    `,
    detailLoading: css`
        min-height: 220px;
        display: flex;
        align-items: center;
        justify-content: center;
    `,
    detail: css`
        display: grid;
        gap: 12px;
    `,
    detailSummary: css`
        padding: 14px;
        border: 1px solid ${token.colorBorderSecondary};
        border-radius: ${token.borderRadius}px;
        background: ${token.colorFillQuaternary};

        .content {
            margin-top: 8px;
            color: ${token.colorText};
            font-size: 13px;
            line-height: 20px;
            word-break: break-all;
        }
    `,
    detailGrid: css`
        display: grid;
        grid-template-columns: repeat(2, minmax(0, 1fr));
        gap: 8px;

        @media (max-width: 575px) {
            grid-template-columns: minmax(0, 1fr);
        }
    `,
    detailInfo: css`
        padding: 10px 12px;
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
            color: ${token.colorText};
            font-size: 13px;
            line-height: 20px;
            word-break: break-all;
        }
    `,
    jsonBlock: css`
        padding: 12px;
        border: 1px solid ${token.colorBorderSecondary};
        border-radius: ${token.borderRadius}px;
        background: ${token.colorBgContainer};

        .title {
            margin-bottom: 8px;
            color: ${token.colorText};
            font-size: 13px;
            font-weight: 600;
            line-height: 20px;
        }

        .text {
            margin-bottom: 0 !important;
            color: ${token.colorTextSecondary};
            font-size: 12px;
            line-height: 18px;
            white-space: pre-wrap;
            word-break: break-all;
        }
    `,
}));

const enumText = (data: any, value: any) => data?.[String(value)] || data?.[value] || "-";
const buildOptions = (data: any) => [
    {label: "全部", value: ALL_VALUE},
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

const DetailInfo = ({label, value, copy}: any) => {
    const text = value === undefined || value === null || value === "" ? "-" : String(value);
    return (
        <div>
            <div className="label">{label}</div>
            <div className="value">
                <Typography.Text copyable={copy && text !== "-" ? {text} : false}>{text}</Typography.Text>
            </div>
        </div>
    );
};

const Log = forwardRef<LogRef, LogProps>(({logEnums}, ref) => {
    const {message} = App.useApp();
    const {styles} = useStyles();
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const detailRef = useRef<ProModalRef>({} as ProModalRef);
    const [orderInfo, setOrderInfo] = useState<any>();
    const [type, setType] = useState(ALL_VALUE);
    const [operatorType, setOperatorType] = useState(ALL_VALUE);
    const [range, setRange] = useState<any>();
    const [page, setPage] = useState(1);
    const [total, setTotal] = useState(0);
    const [list, setList] = useState<any[]>([]);
    const [loading, setLoading] = useState(false);
    const [detail, setDetail] = useState<any>();
    const [detailLoading, setDetailLoading] = useState(false);
    const logTypes = logEnums?.type || {};
    const operatorTypes = logEnums?.operator_type || {};
    const typeOptions = buildOptions(logTypes);
    const operatorOptions = buildOptions(operatorTypes);
    const typeColors = buildTypeColors(logTypes);
    const typeColor = (value: any) => typeColors?.[String(value)] || "default";

    const loadList = () => {
        const orderId = Number(orderInfo?.id);
        if (!orderId) {
            return;
        }
        setLoading(true);
        const body: any = {
            order_id: orderId,
            page,
            page_size: PAGE_SIZE,
            order_by_field: "id",
            order_by_type: "desc",
        };
        if (type !== ALL_VALUE) {
            body.type = type;
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
                setList(res?.data?.list || []);
                setTotal(Number(res?.data?.total || 0));
            },
            onFail: (res: any) => message?.error(res?.message || "获取订单日志失败"),
        })
            .finally(() => {
                setLoading(false);
            })
            .catch(() => undefined);
    };

    useImperativeHandle(ref, () => ({
        open: (record: any) => {
            setOrderInfo(record || {});
            setType(ALL_VALUE);
            setOperatorType(ALL_VALUE);
            setRange(undefined);
            setPage(1);
            setList([]);
            setTotal(0);
            modalRef.current?.show();
        },
    }));

    useEffect(() => {
        if (!orderInfo?.id) {
            return;
        }
        loadList();
    }, [orderInfo?.id, type, operatorType, range, page]);

    const openDetail = (id: any) => {
        const logId = Number(id);
        if (!logId) {
            message?.warning("日志ID错误");
            return;
        }
        setDetail(undefined);
        setDetailLoading(true);
        detailRef.current?.show();
        void getOrderLogInfo({
            body: {id: logId},
            onSuccess: (res: any) => {
                setDetail(res?.data || {});
            },
            onFail: (res: any) => {
                message?.error(res?.message || "获取日志详情失败");
                detailRef.current?.hide();
            },
        })
            .finally(() => {
                setDetailLoading(false);
            })
            .catch(() => undefined);
    };

    const resetFilter = () => {
        setType(ALL_VALUE);
        setOperatorType(ALL_VALUE);
        setRange(undefined);
        setPage(1);
    };

    return (
        <>
            <ProModal
                ref={modalRef}
                title={<Title>订单日志</Title>}
                width={640}
                modalProps={{
                    footer: null,
                    mask: {
                        closable: true,
                    },
                    afterClose: () => {
                        setOrderInfo(undefined);
                        setList([]);
                        setTotal(0);
                        setLoading(false);
                    },
                } as any}>
                <div className={styles.panel}>
                    <div className={styles.orderHeader}>
                        <div className="top">
                            <div className="title">
                                <div className="name">{orderInfo?.name || "订单日志"}</div>
                                <Typography.Text
                                    className="id"
                                    copyable={orderInfo?.id ? {text: String(orderInfo.id)} : false}>
                                    ID {orderInfo?.id || "-"}
                                </Typography.Text>
                            </div>
                        </div>
                    </div>
                    <div className={styles.filter}>
                        <Select
                            className="filterSelect"
                            size="middle"
                            value={type}
                            options={typeOptions}
                            onChange={(value) => {
                                setType(value);
                                setPage(1);
                            }}
                        />
                        <Select
                            className="filterSelect"
                            size="middle"
                            value={operatorType}
                            options={operatorOptions}
                            onChange={(value) => {
                                setOperatorType(value);
                                setPage(1);
                            }}
                        />
                        <DatePicker.RangePicker
                            className="filterDate"
                            value={range}
                            onChange={(value) => {
                                setRange(value);
                                setPage(1);
                            }}
                        />
                        <Button className="reset" onClick={resetFilter}>重置</Button>
                    </div>
                    <Spin spinning={loading}>
                        <div className={styles.listWrap}>
                            {list.length > 0 ? (
                                <div className={styles.list}>
                                    {list.map((item: any) => (
                                        <div className={styles.item} key={item?.id}>
                                            <div className="head">
                                                <div className="main">
                                                    <div className="tags">
                                                        <Tag color={typeColor(item?.type)}>
                                                            {enumText(logTypes, item?.type)}
                                                        </Tag>
                                                        <Typography.Text type="secondary">#{item?.id || "-"}</Typography.Text>
                                                    </div>
                                                    <div className="content">{item?.content || "-"}</div>
                                                </div>
                                                <Button size="small" type="link" onClick={() => openDetail(item?.id)}>
                                                    详情
                                                </Button>
                                            </div>
                                            <div className="meta">
                                                <span className="operator">
                                                    {enumText(operatorTypes, item?.operator_type)}
                                                    {item?.operator_id ? ` #${item.operator_id}` : ""}
                                                </span>
                                                <span className="time">{item?.create_time || "-"}</span>
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
                                onChange={setPage}
                            />
                        </div>
                    )}
                </div>
            </ProModal>
            <ProModal
                ref={detailRef}
                title={<Title>日志详情</Title>}
                width={560}
                modalProps={{
                    footer: null,
                    mask: {
                        closable: true,
                    },
                    afterClose: () => {
                        setDetail(undefined);
                        setDetailLoading(false);
                    },
                } as any}>
                {detailLoading ? (
                    <div className={styles.detailLoading}>
                        <Spin/>
                    </div>
                ) : (
                    <div className={styles.detail}>
                        <div className={styles.detailSummary}>
                            <Space size={6} wrap>
                                <Tag color={typeColor(detail?.type)}>{enumText(logTypes, detail?.type)}</Tag>
                                <Typography.Text type="secondary">#{detail?.id || "-"}</Typography.Text>
                            </Space>
                            <div className="content">{detail?.content || "-"}</div>
                        </div>
                        <div className={`${styles.detailGrid} ${styles.detailInfo}`}>
                            <DetailInfo label="订单ID" value={detail?.order_id} copy/>
                            <DetailInfo label="操作人" value={`${enumText(operatorTypes, detail?.operator_type)} #${detail?.operator_id || 0}`}/>
                            <DetailInfo label="创建时间" value={detail?.create_time}/>
                        </div>
                        <div className={styles.jsonBlock}>
                            <div className="title">变更前</div>
                            <Typography.Paragraph
                                className="text"
                                copyable={detail?.before_info ? {text: formatJson(detail?.before_info)} : false}>
                                {formatJson(detail?.before_info)}
                            </Typography.Paragraph>
                        </div>
                        <div className={styles.jsonBlock}>
                            <div className="title">变更后</div>
                            <Typography.Paragraph
                                className="text"
                                copyable={detail?.after_info ? {text: formatJson(detail?.after_info)} : false}>
                                {formatJson(detail?.after_info)}
                            </Typography.Paragraph>
                        </div>
                    </div>
                )}
            </ProModal>
        </>
    );
});

Log.displayName = "Log";

export default Log;
