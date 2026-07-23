import React, {forwardRef, useImperativeHandle, useRef, useState} from "react";
import {App, Empty, Pagination, Segmented, Space, Spin, Tag, Tooltip, Typography} from "antd";
import {createStyles} from "antd-style";
import {ProModal, ProModalRef, Title} from "sinking-antd";
import {getOrderLogList} from "@/service/user/pay";
import {formatMoney} from "@/utils/string";
import PayType, {getEnumText} from "./pay-type";

const LOG_PAGE_SIZE = 7;
type ActivePanel = "info" | "log";

const moneyItems = [
    {field: "amount", label: "订单金额", color: "#1677ff"},
    {field: "pay_amount", label: "实付金额", color: "#52c41a"},
    {field: "refund_amount", label: "退款金额", color: "#ff4d4f"},
];

const hasValue = (value: any) => value !== undefined && value !== null && value !== "";

const enumTag = (data: any, colors: Record<string, string>, value: any) => (
    <Tag color={colors?.[String(value)] || "default"}>{getEnumText(data, value)}</Tag>
);

const isExpiredTime = (value: any) => {
    if (!value) {
        return false;
    }
    const time = new Date(String(value).replace(/-/g, "/")).getTime();
    return Number.isFinite(time) && time > 0 && time < Date.now();
};

const useStyles: any = createStyles(({token}: any): any => ({
    modalBody: {
        display: "flex",
        flexDirection: "column",
        gap: 10,
    },
    segmentWrap: {
        flex: "none",
        display: "flex",
        justifyContent: "center",
        padding: "2px 0 4px",
    },
    panelContent: {
        minHeight: 0,
        display: "flex",
        flexDirection: "column",
        gap: 10,
        padding: "0 2px 2px",
    },
    metrics: {
        display: "grid",
        gridTemplateColumns: "repeat(3, minmax(0, 1fr))",
        gap: 8,
    },
    metric: {
        minWidth: 0,
        padding: 10,
        border: `1px solid ${token.colorBorderSecondary}`,
        borderRadius: token.borderRadius,
        background: token.colorBgContainer,
        ".label": {
            color: token.colorTextTertiary,
            fontSize: 12,
            lineHeight: "18px",
        },
        ".value": {
            marginTop: 2,
            overflow: "hidden",
            fontSize: 15,
            fontVariantNumeric: "tabular-nums",
            lineHeight: "22px",
            textOverflow: "ellipsis",
            whiteSpace: "nowrap",
        },
    },
    sectionGrid: {
        display: "grid",
        gridTemplateColumns: "repeat(2, minmax(0, 1fr))",
        gap: 8,
        "@media (max-width: 575px)": {
            gridTemplateColumns: "minmax(0, 1fr)",
        },
    },
    section: {
        minWidth: 0,
        padding: "10px 12px",
        border: `1px solid ${token.colorBorderSecondary}`,
        borderRadius: token.borderRadius,
        background: token.colorBgContainer,
        ".sectionTitle": {
            marginBottom: 8,
            color: token.colorText,
            fontSize: 13,
            fontWeight: 500,
            lineHeight: "20px",
        },
        ".infoList": {
            display: "flex",
            flexDirection: "column",
            gap: 6,
        },
        ".infoItem": {
            display: "flex",
            alignItems: "center",
            gap: 8,
            minWidth: 0,
            minHeight: 20,
        },
        ".label": {
            flex: "none",
            color: token.colorTextTertiary,
            fontSize: 12,
            lineHeight: "20px",
        },
        ".value": {
            display: "flex",
            alignItems: "center",
            flex: 1,
            minWidth: 0,
            color: token.colorText,
            fontSize: 12,
            lineHeight: "20px",
            wordBreak: "break-all",
        },
        ".value .ant-typography": {
            fontSize: 12,
            lineHeight: "20px",
        },
        ".value .ant-tag": {
            marginInlineEnd: 0,
            lineHeight: "18px",
        },
    },
    timeSection: {
        gridColumn: "1 / -1",
        "&& > .infoList": {
            display: "grid",
            gridTemplateColumns: "repeat(2, minmax(0, 1fr))",
            gap: "6px 16px",
        },
        "@media (max-width: 430px)": {
            "&& > .infoList": {
                gridTemplateColumns: "minmax(0, 1fr)",
            },
        },
    },
    logModal: {
        ".timeline-shell": {
            minHeight: 0,
        },
        ".timeline": {
            display: "flex",
            flexDirection: "column",
            paddingTop: 2,
        },
        ".timeline-item": {
            position: "relative",
            display: "grid",
            gridTemplateColumns: "22px minmax(0, 1fr)",
            gap: 10,
            paddingBottom: 10,
        },
        ".timeline-item:last-child": {
            paddingBottom: 0,
        },
        ".timeline-item:last-child .line": {
            display: "none",
        },
        ".axis": {
            position: "relative",
            display: "flex",
            alignItems: "flex-start",
            justifyContent: "center",
            minHeight: "100%",
            paddingTop: 13,
        },
        ".dot": {
            position: "relative",
            zIndex: 1,
            width: 8,
            height: 8,
            border: `2px solid ${token.colorPrimary}`,
            borderRadius: "50%",
            background: token.colorBgContainer,
        },
        ".line": {
            position: "absolute",
            top: 0,
            bottom: -12,
            left: "50%",
            width: 1,
            transform: "translateX(-50%)",
            background: token.colorBorderSecondary,
        },
        ".log-card": {
            position: "relative",
            minWidth: 0,
            padding: "10px 12px",
            borderRadius: token.borderRadius,
            background: token.colorFillQuaternary,
        },
        ".log-card:before": {
            content: "\"\"",
            position: "absolute",
            top: 11,
            left: -6,
            width: 0,
            height: 0,
            borderTop: "6px solid transparent",
            borderBottom: "6px solid transparent",
            borderRight: `6px solid ${token.colorFillQuaternary}`,
        },
        ".log-head": {
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            gap: 8,
            minWidth: 0,
            marginBottom: 6,
        },
        ".log-meta": {
            display: "flex",
            alignItems: "center",
            flexWrap: "wrap",
            gap: 6,
            minWidth: 0,
        },
        ".ant-tag": {
            marginInlineEnd: 0,
        },
        ".log-time": {
            flex: "none",
            color: token.colorTextTertiary,
            fontSize: 12,
            lineHeight: "18px",
            fontVariantNumeric: "tabular-nums",
            whiteSpace: "nowrap",
        },
        ".log-content": {
            color: token.colorTextSecondary,
            fontSize: 12,
            lineHeight: "18px",
            wordBreak: "break-word",
        },
        ".log-extra": {
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            gap: 8,
            marginTop: 6,
            color: token.colorTextTertiary,
            fontSize: 12,
            lineHeight: "18px",
        },
        ".log-operator": {
            display: "flex",
            alignItems: "center",
            gap: 4,
            flex: 1,
            minWidth: 0,
            overflow: "hidden",
            whiteSpace: "nowrap",
        },
        ".operator-type": {
            flex: "none",
            lineHeight: "18px",
            whiteSpace: "nowrap",
        },
        ".operator-id": {
            position: "relative",
            display: "inline-block",
            flex: "none",
            paddingRight: 20,
            height: 18,
            lineHeight: "18px",
            whiteSpace: "nowrap",
        },
        ".operator-id-text": {
            color: token.colorTextTertiary,
            fontSize: 12,
            height: 18,
            lineHeight: "18px",
            whiteSpace: "nowrap",
        },
        ".operator-copy": {
            position: "absolute",
            top: 1,
            right: 0,
            width: 16,
            height: 16,
            overflow: "hidden",
            color: token.colorTextTertiary,
            fontSize: 12,
            lineHeight: 1,
        },
        ".operator-copy .ant-typography-copy": {
            display: "inline-flex",
            alignItems: "center",
            justifyContent: "center",
            position: "absolute",
            inset: 0,
            width: 16,
            height: 16,
            marginInlineStart: 0,
            lineHeight: 1,
        },
        ".operator-copy .ant-typography-copy svg": {
            display: "block",
        },
        ".operator-copy .ant-typography-actions": {
            display: "block",
            width: 16,
            height: 16,
        },
        ".log-pagination": {
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            gap: 12,
            marginTop: 12,
            color: token.colorTextTertiary,
            fontSize: 12,
            lineHeight: "18px",
        },
        ".empty-wrap": {
            display: "grid",
            placeItems: "center",
            minHeight: 80,
        },
        ".loading-wrap": {
            display: "grid",
            placeItems: "center",
            minHeight: 80,
        },
        "@media (max-width: 575px)": {
            ".timeline-item": {
                gridTemplateColumns: "18px minmax(0, 1fr)",
                gap: 8,
            },
            ".log-head": {
                alignItems: "flex-start",
                flexDirection: "column",
                gap: 4,
            },
            ".log-time": {
                whiteSpace: "normal",
            },
            ".log-extra": {
                alignItems: "flex-end",
            },
            ".log-pagination": {
                justifyContent: "space-between",
            },
        },
    },
}));

export type InfoRef = {
    open: (record: any, panel?: ActivePanel) => void;
};

type InfoProps = {
    orderType: any;
    orderStatus: any;
    payStatus: any;
    logTypes: any;
    operatorTypes: any;
    payTypes: any;
    typeColors: Record<string, string>;
    statusColors: Record<string, string>;
    payStatusColors: Record<string, string>;
    logTypeColors: Record<string, string>;
};

const Info = forwardRef<InfoRef, InfoProps>((props, ref): React.ReactNode => {
    const {message} = App.useApp();
    const {styles} = useStyles();
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const [data, setData] = useState<any>();
    const [logs, setLogs] = useState<any[]>([]);
    const [logTotal, setLogTotal] = useState(0);
    const [logPage, setLogPage] = useState(1);
    const [logLoading, setLogLoading] = useState(false);
    const [logLoaded, setLogLoaded] = useState(false);
    const [activePanel, setActivePanel] = useState<ActivePanel>("info");

    const reset = () => {
        setActivePanel("info");
        setData(undefined);
        setLogs([]);
        setLogTotal(0);
        setLogPage(1);
        setLogLoaded(false);
        setLogLoading(false);
    };

    useImperativeHandle(ref, () => ({
        open: (record: any, panel: ActivePanel = "info") => {
            reset();
            setData(record);
            setActivePanel(panel);
            modalRef.current?.show();
            if (panel === "log") {
                void loadLogs(1, record);
            }
        },
    }));

    const formatTime = (value: any) => {
        const text = value === undefined || value === null ? "" : String(value);
        return !text || text.startsWith("0001-01-01") ? "-" : text;
    };

    const loadLogs = async (page = 1, record: any = data) => {
        if (!record?.id) {
            return;
        }
        setLogPage(page);
        setLogLoading(true);
        await getOrderLogList({
            body: {
                order_id: Number(record.id),
                page,
                page_size: LOG_PAGE_SIZE,
                order_by_field: "id",
                order_by_type: "desc",
            },
            onSuccess: (res: any) => {
                setLogs(res?.data?.list || []);
                setLogTotal(Number(res?.data?.total || 0));
                setLogLoaded(true);
            },
            onFail: (res: any) => message?.error(res?.message || "获取订单日志失败"),
        }).finally(() => {
            setLogLoading(false);
        });
    };

    const changePanel = (value: ActivePanel) => {
        setActivePanel(value);
        if (value === "log" && !logLoaded && !logLoading) {
            void loadLogs(1);
        }
    };

    const InfoRow = ({label, value, copy}: { label: string; value: any; copy?: any }) => {
        const isElement = React.isValidElement(value);
        if (!hasValue(value) && !isElement) {
            return null;
        }
        const text = hasValue(value) ? String(value) : "-";
        return (
            <div className="infoItem">
                <div className="label">{label}</div>
                <div className="value">
                    {isElement ? value : (
                        <Typography.Text copyable={copy ? {text: String(copy)} : false}>{text}</Typography.Text>
                    )}
                </div>
            </div>
        );
    };

    const renderInfoRows = (items: { label: string; value: any; copy?: any }[]) => (
        items.map((item) => <InfoRow key={item.label} {...item}/>)
    );

    const Metric = ({label, value, color}: { label: string; value: any; color: string }) => (
        <div className={styles.metric}>
            <div className="label">{label}</div>
            <div className="value" style={{color}}>{formatMoney(value)}</div>
        </div>
    );

    const ExpiredTime = ({value}: { value: any }) => {
        const text = formatTime(value);
        if (text === "-") {
            return <Typography.Text>{text}</Typography.Text>;
        }
        if (isExpiredTime(value)) {
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

    const renderInfo = () => {
        const orderInfoItems = [
            {label: "订单名称", value: data?.name},
            {label: "订单类型", value: hasValue(data?.type) ? enumTag(props.orderType, props.typeColors, data?.type) : undefined},
            {label: "订单ID", value: data?.id, copy: data?.id},
            {label: "外部ID", value: data?.trade_no, copy: data?.trade_no},
            {label: "所属用户", value: data?.user_id, copy: data?.user_id},
        ];
        const payInfoItems = [
            {label: "支付方式", value: hasValue(data?.pay_type) ? <PayType data={props.payTypes} value={data?.pay_type}/> : undefined},
            {label: "订单状态", value: hasValue(data?.status) ? enumTag(props.orderStatus, props.statusColors, data?.status) : undefined},
            {label: "支付状态", value: hasValue(data?.pay_status) ? enumTag(props.payStatus, props.payStatusColors, data?.pay_status) : undefined},
            {label: "通道名称", value: data?.channel_name},
            {label: "通道ID", value: data?.channel_id, copy: data?.channel_id},
        ];
        const timeInfoItems = [
            {label: "创建时间", value: formatTime(data?.create_time)},
            {label: "支付时间", value: formatTime(data?.pay_time)},
            {label: "完成时间", value: formatTime(data?.complete_time)},
            {label: "过期时间", value: <ExpiredTime value={data?.expired_time}/>},
        ];

        return (
            <>
                <div className={styles.metrics}>
                    {moneyItems.map((item) => (
                        <Metric
                            key={item.field}
                            label={item.label}
                            value={data?.[item.field]}
                            color={item.color}
                        />
                    ))}
                </div>
                <div className={styles.sectionGrid}>
                    <div className={styles.section}>
                        <div className="sectionTitle">订单信息</div>
                        <div className="infoList">{renderInfoRows(orderInfoItems)}</div>
                    </div>
                    <div className={styles.section}>
                        <div className="sectionTitle">支付信息</div>
                        <div className="infoList">{renderInfoRows(payInfoItems)}</div>
                    </div>
                    <div className={`${styles.section} ${styles.timeSection}`}>
                        <div className="sectionTitle">时间信息</div>
                        <div className="infoList">{renderInfoRows(timeInfoItems)}</div>
                    </div>
                </div>
            </>
        );
    };

    const renderLogExtra = (item: any) => {
        if (!hasValue(item?.operator_type) && !hasValue(item?.operator_id)) {
            return null;
        }
        return (
            <div className="log-extra">
                <span className="log-operator">
                    {hasValue(item?.operator_type) ? (
                        <span className="operator-type">操作人：{getEnumText(props.operatorTypes, item?.operator_type)}</span>
                    ) : null}
                    {hasValue(item?.operator_id) ? (
                        <span className="operator-id">
                            <span className="operator-id-text">账号ID：{item.operator_id}</span>
                            <Typography.Text
                                className="operator-copy"
                                copyable={{text: String(item.operator_id), tooltips: ["复制", "复制成功"]}}
                            />
                        </span>
                    ) : null}
                </span>
            </div>
        );
    };

    const renderLogList = () => {
        if (logs.length === 0) {
            return (
                <div className="empty-wrap">
                    <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="暂无日志"/>
                </div>
            );
        }
        return (
            <div className="timeline">
                {logs.map((item: any) => (
                    <div className="timeline-item" key={item?.id || `${item?.type}-${item?.create_time}`}>
                        <div className="axis">
                            <span className="dot"/>
                            <span className="line"/>
                        </div>
                        <div className="log-card">
                            <div className="log-head">
                                <div className="log-meta">
                                    {hasValue(item?.type) ? enumTag(props.logTypes, props.logTypeColors, item?.type) : null}
                                </div>
                                <Typography.Text className="log-time">{formatTime(item?.create_time)}</Typography.Text>
                            </div>
                            {hasValue(item?.content) ? <div className="log-content">{item.content}</div> : null}
                            {renderLogExtra(item)}
                        </div>
                    </div>
                ))}
            </div>
        );
    };

    return (
        <ProModal
            ref={modalRef}
            title={<Title>订单详情</Title>}
            width={600}
            modalProps={{
                footer: null,
                forceRender: true,
                style: {top: 100, paddingBottom: 100},
                mask: {closable: true},
                styles: {
                    body: {
                        padding: "2px 6px 10px",
                    },
                },
                destroyOnHidden: true,
                afterClose: reset,
            } as any}>
            <div className={styles.modalBody}>
                <div className={styles.segmentWrap}>
                    <Segmented
                        value={activePanel}
                        options={[
                            {label: "详情", value: "info"},
                            {label: "日志", value: "log"},
                        ]}
                        onChange={(value) => changePanel(value as ActivePanel)}
                    />
                </div>
                <div className={styles.panelContent}>
                    {activePanel === "log" ? (
                        <div className={styles.logModal}>
                            {logLoading ? (
                                <div className="loading-wrap">
                                    <Spin/>
                                </div>
                            ) : (
                                <>
                                    <div className="timeline-shell">
                                        {renderLogList()}
                                    </div>
                                    {logTotal > LOG_PAGE_SIZE ? (
                                        <div className="log-pagination">
                                            <span>共 {logTotal} 条</span>
                                            <Pagination
                                                size="small"
                                                current={logPage}
                                                pageSize={LOG_PAGE_SIZE}
                                                total={logTotal}
                                                showSizeChanger={false}
                                                onChange={(page) => void loadLogs(page)}
                                            />
                                        </div>
                                    ) : null}
                                </>
                            )}
                        </div>
                    ) : renderInfo()}
                </div>
            </div>
        </ProModal>
    );
});

export default Info;
