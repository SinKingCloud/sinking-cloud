import React, {useRef, useState} from "react";
import {App, Button, Empty, Timeline, Typography} from "antd";
import {createStyles} from "antd-style";
import {Body, ProModal, ProModalRef, ProTable, ProTableRef, Title} from "sinking-antd";
import useEnum from "@/utils/enum";
import {buildValueEnum, dateRangeTransform, getData} from "@/utils/page";
import {formatMoney} from "@/utils/string";
import {getOrderList, getOrderLogList} from "@/service/user/pay";

const LOG_PAGE_SIZE = 20;
const typeColors: Record<string, string> = {"0": "blue", "1": "purple"};
const statusColors: Record<string, string> = {"0": "orange", "1": "gold", "2": "green", "3": "red"};
const payStatusColors: Record<string, string> = {"0": "orange", "1": "green"};
const payTypeColors: Record<string, string> = {alipay: "blue", wxpay: "green", wechat: "green", balance: "gold"};

const enumText = (data: any, value: any) => (
    data?.[String(value)] || data?.[value] || (value === undefined || value === null || value === "" ? "-" : String(value))
);

const hasValue = (value: any) => value !== undefined && value !== null && value !== "";

const useStyles: any = createStyles(({token}: any) => ({
    page: {
        width: "100%",
        maxWidth: 1440,
        margin: "0 auto",
    },
    money: {
        color: token.colorPrimary,
        fontVariantNumeric: "tabular-nums",
    },
    payMoney: {
        color: token.colorSuccess,
        fontVariantNumeric: "tabular-nums",
    },
    refundMoney: {
        color: token.colorError,
        fontVariantNumeric: "tabular-nums",
    },
    logModal: {
        ".log-summary": {
            marginBottom: 16,
            padding: "10px 12px",
            borderRadius: token.borderRadiusLG,
            background: token.colorFillQuaternary,
        },
        ".log-title": {
            display: "block",
            color: token.colorText,
            fontSize: 13,
            fontWeight: 600,
            lineHeight: "20px",
        },
        ".log-id": {
            display: "inline-flex",
            marginTop: 2,
            color: token.colorTextTertiary,
            fontSize: 12,
            lineHeight: "18px",
        },
        ".log-id .ant-typography": {
            color: "inherit",
            fontSize: "inherit",
            lineHeight: "inherit",
        },
        ".log-card": {
            padding: "10px 12px",
            borderRadius: token.borderRadiusLG,
            background: token.colorFillQuaternary,
        },
        ".log-head": {
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            gap: 12,
            marginBottom: 6,
        },
        ".log-type": {
            color: token.colorText,
            fontSize: 13,
            fontWeight: 600,
            lineHeight: "20px",
        },
        ".log-time": {
            flex: "none",
            color: token.colorTextTertiary,
            fontSize: 12,
            lineHeight: "18px",
        },
        ".log-content": {
            color: token.colorTextSecondary,
            fontSize: 12,
            lineHeight: "20px",
            wordBreak: "break-word",
        },
        ".empty-wrap": {
            display: "grid",
            placeItems: "center",
            minHeight: 180,
        },
    },
}));

export default (): React.ReactNode => {
    const {message} = App.useApp();
    const {styles} = useStyles();
    const tableRef = useRef<ProTableRef>({} as ProTableRef);
    const logModalRef = useRef<ProModalRef>({} as ProModalRef);
    const [enumData, enumLoading] = useEnum(["pay.order", "pay.order_log", "pay.pay"]) as any;
    const [currentOrder, setCurrentOrder] = useState<any>();
    const [logs, setLogs] = useState<any[]>([]);
    const [logLoading, setLogLoading] = useState(false);

    const orderEnums = enumData?.["pay.order"] || {};
    const logEnums = enumData?.["pay.order_log"] || {};
    const payEnums = enumData?.["pay.pay"] || {};
    const orderType = orderEnums?.type || {};
    const orderStatus = orderEnums?.status || {};
    const payStatus = orderEnums?.pay_status || {};
    const payTypes = payEnums?.type || {};

    const formatTime = (value: any) => {
        const text = value === undefined || value === null ? "" : String(value);
        return !text || text.startsWith("0001-01-01") ? "-" : text;
    };

    const loadLogs = async (record: any) => {
        setCurrentOrder(record);
        setLogs([]);
        logModalRef.current?.show();
        setLogLoading(true);
        await getOrderLogList({
            body: {
                order_id: Number(record?.id),
                page: 1,
                page_size: LOG_PAGE_SIZE,
                order_by_field: "id",
                order_by_type: "desc",
            },
            onSuccess: (res: any) => setLogs(res?.data?.list || []),
            onFail: (res: any) => message?.error(res?.message || "获取订单日志失败"),
        }).finally(() => {
            setLogLoading(false);
        });
    };

    const columns: any[] = [
        {
            title: "订单ID",
            dataIndex: "id",
            tip: "订单ID",
            valueType: "text",
            sorter: true,
            copyable: true,
            width: 180,
        },
        {
            title: "订单名称",
            dataIndex: "name",
            tip: "订单名称",
            width: 180,
            hideInSearch: true,
            ellipsis: true,
        },
        {
            title: "订单类型",
            dataIndex: "type",
            tip: "订单类型",
            valueEnum: buildValueEnum(orderType, typeColors),
            width: 110,
        },
        {
            title: "订单状态",
            dataIndex: "status",
            tip: "订单状态",
            valueEnum: buildValueEnum(orderStatus, statusColors),
            width: 110,
        },
        {
            title: "支付状态",
            dataIndex: "pay_status",
            tip: "支付状态",
            valueEnum: buildValueEnum(payStatus, payStatusColors),
            width: 110,
        },
        {
            title: "支付方式",
            dataIndex: "pay_type",
            tip: "支付方式",
            valueEnum: buildValueEnum(payTypes, payTypeColors),
            width: 110,
        },
        {
            title: "订单金额",
            dataIndex: "amount",
            tip: "订单金额",
            width: 110,
            hideInSearch: true,
            render: (value: any) => <Typography.Text className={styles.money}>{formatMoney(value)}</Typography.Text>,
        },
        {
            title: "实付金额",
            dataIndex: "pay_amount",
            tip: "实付金额",
            width: 110,
            hideInSearch: true,
            render: (value: any) => <Typography.Text className={styles.payMoney}>{formatMoney(value)}</Typography.Text>,
        },
        {
            title: "退款金额",
            dataIndex: "refund_amount",
            tip: "退款金额",
            width: 110,
            hideInSearch: true,
            render: (value: any) => <Typography.Text className={styles.refundMoney}>{formatMoney(value)}</Typography.Text>,
        },
        {
            title: "支付时间",
            valueType: "dateRange",
            dataIndex: "pay_time",
            tip: "支付时间",
            sorter: true,
            hideInTable: true,
            transform: dateRangeTransform("pay_time"),
        },
        {
            title: "完成时间",
            valueType: "dateRange",
            dataIndex: "complete_time",
            tip: "完成时间",
            sorter: true,
            hideInTable: true,
            transform: dateRangeTransform("complete_time"),
        },
        {
            title: "创建时间",
            valueType: "dateRange",
            dataIndex: "create_time",
            tip: "创建时间",
            sorter: true,
            hideInTable: true,
            transform: dateRangeTransform("create_time"),
        },
        {
            title: "创建时间",
            dataIndex: "create_time",
            tip: "创建时间",
            width: 170,
            hideInSearch: true,
            render: (value: any) => formatTime(value),
        },
        {
            title: "支付时间",
            dataIndex: "pay_time",
            tip: "支付时间",
            width: 170,
            hideInSearch: true,
            render: (value: any) => formatTime(value),
        },
        {
            title: "完成时间",
            dataIndex: "complete_time",
            tip: "完成时间",
            width: 170,
            hideInSearch: true,
            render: (value: any) => formatTime(value),
        },
        {
            title: "操作",
            valueType: "option",
            key: "option",
            width: 80,
            fixed: "right",
            hideInSearch: true,
            render: (_: any, record: any) => (
                <Button size="small" onClick={() => loadLogs(record)}>
                    日志
                </Button>
            ),
        },
    ];

    return (
        <Body loading={enumLoading}>
            <div className={styles.page}>
                <ProTable
                    ref={tableRef}
                    rowKey="id"
                    title={<Title>订单管理</Title>}
                    extraRefreshBtn={true}
                    search={true}
                    columns={columns}
                    defaultPage={1}
                    defaultPageSize={20}
                    request={(params: any, sort: any) => getData(params, sort, getOrderList)}
                    paginationAffix={true}
                    {...({
                        tableProps: {
                            scroll: {x: 1800},
                        },
                    } as any)}
                />
            </div>
            <ProModal
                ref={logModalRef}
                title={<Title>订单日志</Title>}
                width={560}
                modalProps={{
                    footer: null,
                    style: {top: 100, paddingBottom: 100},
                    mask: {closable: true},
                    destroyOnHidden: true,
                } as any}>
                <Body loading={logLoading} className={styles.logModal}>
                    <div className="log-summary">
                        <Typography.Text className="log-title" ellipsis={{tooltip: currentOrder?.name || "-"}}>
                            {currentOrder?.name || "-"}
                        </Typography.Text>
                        <div className="log-id">
                            <Typography.Text copyable={hasValue(currentOrder?.id) ? {text: String(currentOrder.id)} : false}>
                                订单ID：{currentOrder?.id || "-"}
                            </Typography.Text>
                        </div>
                    </div>
                    {logs.length > 0 ? (
                        <Timeline
                            items={logs.map((item: any) => ({
                                children: (
                                    <div className="log-card">
                                        <div className="log-head">
                                            <span className="log-type">{enumText(logEnums?.type, item?.type)}</span>
                                            <span className="log-time">{formatTime(item?.create_time)}</span>
                                        </div>
                                        <div className="log-content">{item?.content || "-"}</div>
                                    </div>
                                ),
                            }))}
                        />
                    ) : (
                        <div className="empty-wrap">
                            <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="暂无日志"/>
                        </div>
                    )}
                </Body>
            </ProModal>
        </Body>
    );
};
