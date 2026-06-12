import React, {useRef} from "react";
import {App, Button, Dropdown, Result, Space, Tag, Typography} from "antd";
import {Body, ProTable, ProTableRef, Title} from "sinking-antd";
import {createStyles} from "antd-style";
import {buildValueEnum, dateRangeTransform, emptyTableData, getData} from "@/utils/page";
import useControl from "@/utils/control";
import useEnum from "@/utils/enum";
import {AdminControlCode} from "@/utils/permission";
import {formatMoney} from "@/utils/string";
import {getOrderList} from "@/service/admin/pay";
import Info, {type InfoRef} from "./components/info";
import Refund, {type RefundRef} from "./components/refund";
import Clear, {type ClearRef} from "./components/clear";
import Deliver, {type DeliverRef} from "./components/deliver";
import PayType from "./components/pay-type";

const useStyles: any = createStyles((): any => ({
    infoColumn: {
        width: "100%",
    },
    infoRow: {
        display: "flex",
        alignItems: "center",
        gap: 4,
        minWidth: 0,
    },
    labelText: {
        flex: "none",
        fontSize: 12,
        lineHeight: "18px",
    },
    smallText: {
        flex: 1,
        minWidth: 0,
        fontSize: 12,
        lineHeight: "18px",
        whiteSpace: "nowrap",
    },
    strongText: {
        flex: 1,
        minWidth: 0,
        fontSize: 12,
        fontWeight: 500,
        lineHeight: "18px",
        whiteSpace: "nowrap",
    },
    moneyText: {
        flex: 1,
        minWidth: 0,
        fontSize: 12,
        fontVariantNumeric: "tabular-nums",
        lineHeight: "18px",
        whiteSpace: "nowrap",
    },
    tagText: {
        flex: 1,
        minWidth: 0,
        display: "flex",
        alignItems: "center",
    },
}));

const enumText = (data: any, value: any) => (
    data?.[String(value)] || data?.[value] || (value === undefined || value === null || value === "" ? "-" : String(value))
);
const numberValue = (value: any) => Number(value || 0);
const hasValue = (value: any) => value !== undefined && value !== null && value !== "";

export default (): React.ReactNode => {
    const {message} = App.useApp();
    const {styles} = useStyles();
    const control = useControl();
    const [enumData, enumLoading] = useEnum(["pay.order", "pay.order_log", "pay.channel", "pay.pay"]) as any;
    const tableRef = useRef<ProTableRef>({} as ProTableRef);
    const infoRef = useRef<InfoRef>({} as InfoRef);
    const deliverRef = useRef<DeliverRef>({} as DeliverRef);
    const refundRef = useRef<RefundRef>({} as RefundRef);
    const clearRef = useRef<ClearRef>({} as ClearRef);

    const hasListControl = control.hasControl(AdminControlCode.WEB_ORDER_LIST);
    const hasInfoControl = control.hasControl(AdminControlCode.WEB_ORDER_INFO);
    const hasDeliverControl = control.hasControl(AdminControlCode.WEB_ORDER_DELIVER);
    const hasRefundControl = control.hasControl(AdminControlCode.WEB_ORDER_REFOUND);
    const hasClearControl = control.hasControl(AdminControlCode.WEB_ORDER_CLEAR);
    const hasActionControl = hasInfoControl || hasDeliverControl || hasRefundControl;

    const orderEnums = enumData?.["pay.order"] || {};
    const logEnums = enumData?.["pay.order_log"] || {};
    const channelEnums = enumData?.["pay.channel"] || {};
    const payEnums = enumData?.["pay.pay"] || {};
    const orderType = orderEnums?.type || {};
    const orderStatus = orderEnums?.status || {};
    const payStatus = orderEnums?.pay_status || {};
    const channelNames = channelEnums?.name || {};
    const payTypes = payEnums?.type || {};
    const typeColors: Record<string, string> = {"0": "purple", "1": "cyan"};
    const statusColors: Record<string, string> = {"0": "default", "1": "orange", "2": "green", "3": "red"};
    const payStatusColors: Record<string, string> = {"0": "orange", "1": "green"};
    const typeColor = (value: any) => typeColors?.[String(value)] || "default";
    const statusColor = (value: any) => statusColors?.[String(value)] || "default";
    const payStatusColor = (value: any) => payStatusColors?.[String(value)] || "default";
    const typeValueEnum = buildValueEnum(orderType, typeColors);
    const statusValueEnum = buildValueEnum(orderStatus, statusColors);
    const payStatusValueEnum = buildValueEnum(payStatus, payStatusColors);
    const canDeliver = (record: any) => [0, 1].includes(numberValue(record?.status));
    const canRefund = (record: any) => (
        numberValue(record?.type) !== 0
        && numberValue(record?.pay_status) === 1
        && numberValue(record?.pay_amount) > numberValue(record?.refund_amount)
    );

    const refresh = () => {
        tableRef.current?.refreshTableData?.();
        tableRef.current?.clearSelectedRows?.();
    };

    const openClear = () => {
        if (!hasClearControl) {
            message?.warning("暂无订单清理权限");
            return;
        }
        clearRef.current?.open();
    };

    const openDeliver = (record: any) => deliverRef.current?.open(record);

    const openRefund = (record: any) => {
        if (!hasRefundControl) {
            message?.warning("暂无订单退款权限");
            return;
        }
        if (!canRefund(record)) {
            message?.warning("当前订单不可退款");
            return;
        }
        refundRef.current?.open(record);
    };

    const renderRows = (items: {
        label: string;
        value: any;
        copy?: any;
        className?: string;
        style?: React.CSSProperties
    }[]) => (
        <Space orientation="vertical" size={2} className={styles.infoColumn}>
            {items.map((item) => (
                <div className={styles.infoRow} key={item.label}>
                    <Typography.Text type="secondary" className={styles.labelText}>{item.label}</Typography.Text>
                    {React.isValidElement(item.value) ? (
                        <span className={styles.tagText}>{item.value}</span>
                    ) : (
                        <Typography.Text
                            className={item.className || styles.smallText}
                            style={item.style}
                            copyable={item.copy ? {text: String(item.copy)} : false}>
                            {hasValue(item.value) ? item.value : "-"}
                        </Typography.Text>
                    )}
                </div>
            ))}
        </Space>
    );

    const renderOrderCell = (record: any) => renderRows([
        {label: "订单名称:", value: record?.name, className: styles.strongText},
        {label: "订单类型:", value: <Tag color={typeColor(record?.type)}>{enumText(orderType, record?.type)}</Tag>},
        {label: "所属用户:", value: record?.user_id, copy: record?.user_id},
    ]);

    const renderNoCell = (record: any) => renderRows([
        {label: "订单:", value: record?.id ? `${record.id}` : "-", copy: record?.id, className: styles.strongText},
        {label: "外部:", value: record?.trade_no, copy: record?.trade_no},
    ]);

    const renderPayCell = (record: any) => renderRows([
        {label: "支付方式:", value: <PayType data={payTypes} value={record?.pay_type}/>},
        {label: "支付通道:", value: hasValue(record?.channel_id) ? enumText(channelNames, record?.channel_id) : "-"},
        {label: "通道ID:", value: record?.channel_id, copy: record?.channel_id},
    ]);

    const renderStatusCell = (record: any) => renderRows([
        {
            label: "订单状态:",
            value: <Tag color={statusColor(record?.status)}>{enumText(orderStatus, record?.status)}</Tag>
        },
        {
            label: "支付状态:",
            value: <Tag color={payStatusColor(record?.pay_status)}>{enumText(payStatus, record?.pay_status)}</Tag>
        },
    ]);

    const renderMoneyCell = (record: any) => renderRows([
        {
            label: "订单金额:",
            value: formatMoney(record?.amount),
            className: styles.moneyText,
            style: {color: "#1677ff"}
        },
        {
            label: "实付金额:",
            value: formatMoney(record?.pay_amount),
            className: styles.moneyText,
            style: {color: "#52c41a"}
        },
        {
            label: "退款金额:",
            value: formatMoney(record?.refund_amount),
            className: styles.moneyText,
            style: {color: "#ff4d4f"}
        },
    ]);

    const renderTimeCell = (record: any) => renderRows([
        {label: "创建:", value: record?.create_time},
        {label: "支付:", value: record?.pay_time},
        {label: "完成:", value: record?.complete_time},
    ]);

    const columns: any[] = [
        {
            title: "编号信息",
            dataIndex: "no_info",
            tip: "订单编号/外部编号",
            width: 250,
            hideInSearch: true,
            render: (_: any, record: any) => renderNoCell(record),
        },
        {
            title: "订单ID",
            dataIndex: "id",
            tip: "订单ID",
            valueType: "text",
            hideInTable: true,
        },
        {
            title: "外部ID",
            dataIndex: "trade_no",
            tip: "外部ID",
            valueType: "text",
            hideInTable: true,
        },
        {
            title: "订单信息",
            dataIndex: "id",
            tip: "订单名称/订单类型/所属用户",
            sorter: true,
            width: 240,
            hideInSearch: true,
            render: (_: any, record: any) => renderOrderCell(record),
        },
        {
            title: "支付信息",
            dataIndex: "pay_info",
            tip: "支付方式/支付通道/通道ID",
            width: 220,
            hideInSearch: true,
            render: (_: any, record: any) => renderPayCell(record),
        },
        {
            title: "用户ID",
            dataIndex: "user_id",
            tip: "用户ID",
            valueType: "text",
            hideInTable: true,
        },
        {
            title: "通道ID",
            dataIndex: "channel_id",
            tip: "通道ID",
            valueType: "text",
            hideInTable: true,
        },
        {
            title: "状态信息",
            dataIndex: "status_info",
            tip: "订单状态/支付状态",
            width: 210,
            hideInSearch: true,
            render: (_: any, record: any) => renderStatusCell(record),
        },
        {
            title: "订单类型",
            dataIndex: "type",
            tip: "订单类型",
            valueEnum: typeValueEnum,
            hideInTable: true,
        },
        {
            title: "订单状态",
            dataIndex: "status",
            tip: "订单状态",
            valueEnum: statusValueEnum,
            hideInTable: true,
        },
        {
            title: "支付状态",
            dataIndex: "pay_status",
            tip: "支付状态",
            valueEnum: payStatusValueEnum,
            hideInTable: true,
        },
        {
            title: "支付方式",
            dataIndex: "pay_type",
            tip: "支付方式",
            valueEnum: buildValueEnum(payTypes),
            hideInTable: true,
        },
        {
            title: "金额信息",
            dataIndex: "amount",
            tip: "订单/支付/退款金额",
            width: 190,
            hideInSearch: true,
            render: (_: any, record: any) => renderMoneyCell(record),
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
            title: "更新时间",
            valueType: "dateRange",
            dataIndex: "update_time",
            tip: "更新时间",
            sorter: true,
            hideInTable: true,
            transform: dateRangeTransform("update_time"),
        },
        {
            title: "相关时间",
            dataIndex: "time_info",
            tip: "创建/支付/完成时间",
            width: 220,
            hideInSearch: true,
            render: (_: any, record: any) => renderTimeCell(record),
        },
        ...(hasActionControl ? [{
            title: "操作",
            valueType: "option",
            key: "option",
            width: 88,
            fixed: "right",
            hideInSearch: true,
            render: (_: any, record: any) => {
                const items = [
                    hasInfoControl ? {
                        key: "info",
                        label: "详情",
                        onClick: () => infoRef.current?.open(record, "info"),
                    } : null,
                    hasInfoControl ? {
                        key: "log",
                        label: "日志",
                        onClick: () => infoRef.current?.open(record, "log"),
                    } : null,
                    hasDeliverControl ? {
                        key: "deliver",
                        label: "发货",
                        disabled: !canDeliver(record),
                        onClick: () => openDeliver(record),
                    } : null,
                    hasRefundControl ? {
                        key: "refund",
                        label: "退款",
                        disabled: !canRefund(record),
                        onClick: () => openRefund(record),
                    } : null,
                ].filter(Boolean) as any[];
                if (items.length <= 0) {
                    return "-";
                }
                return (
                    <Dropdown menu={{items}} trigger={["click"]} placement="bottom" arrow>
                        <Button size="small">操作</Button>
                    </Dropdown>
                );
            },
        }] : []),
    ];

    return (
        <Body loading={control?.loading || enumLoading}>
            {!hasListControl && <Result
                status="403"
                title="403"
                subTitle="您无权限访问此页面,请联系管理员"
            />}
            {hasListControl && <>
                <ProTable
                    ref={tableRef}
                    extraRefreshBtn={true}
                    title={<Title>订单管理</Title>}
                    rowKey="id"
                    columns={columns}
                    defaultPage={1}
                    defaultPageSize={20}
                    extra={hasClearControl ? (
                        <Button type="primary" onClick={openClear}>
                            清理
                        </Button>
                    ) : undefined}
                    request={(params, sort) => {
                        if (!hasListControl) {
                            return Promise.resolve(emptyTableData);
                        }
                        return getData(params, sort, getOrderList);
                    }}
                    paginationAffix={true}
                    selectionAffix={true}
                />
                <Info
                    ref={infoRef}
                    orderEnums={orderEnums}
                    logEnums={logEnums}
                    payTypes={payTypes}
                    channelNames={channelNames}
                    hasDeliverControl={hasDeliverControl}
                    canDeliver={canDeliver}
                    onDeliver={openDeliver}
                    hasRefundControl={hasRefundControl}
                    canRefund={canRefund}
                    onRefundSuccess={refresh}
                    typeColors={typeColors}
                    statusColors={statusColors}
                    payStatusColors={payStatusColors}
                />
                <Deliver
                    ref={deliverRef}
                    hasControl={hasDeliverControl}
                    canDeliver={canDeliver}
                    onSuccess={refresh}
                />
                <Refund ref={refundRef} orderEnums={orderEnums} typeColors={typeColors} onSuccess={refresh}/>
                <Clear ref={clearRef} onSuccess={refresh}/>
            </>}
        </Body>
    );
};
