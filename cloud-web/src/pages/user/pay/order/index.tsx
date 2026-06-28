import React, {useRef} from "react";
import {Button, Dropdown, Space, Tag, Typography} from "antd";
import {createStyles} from "antd-style";
import {Body, ProTable, Title} from "sinking-antd";
import useEnum from "@/utils/enum";
import {buildValueEnum, dateRangeTransform, getData} from "@/utils/page";
import {formatMoney} from "@/utils/string";
import {getOrderList} from "@/service/user/pay";
import Info, {type InfoRef} from "./components/info";
import PayType, {getEnumText} from "./components/pay-type";

const typeColors: Record<string, string> = {"0": "purple", "1": "cyan"};
const statusColors: Record<string, string> = {"0": "default", "1": "orange", "2": "green", "3": "red"};
const payStatusColors: Record<string, string> = {"0": "orange", "1": "green"};
const payTypeColors: Record<string, string> = {alipay: "blue", wxpay: "green", wechat: "green", balance: "gold"};
const logTypeColors: Record<string, string> = {"0": "blue", "1": "green", "2": "red"};
const moneyItems = [
    {field: "amount", label: "订单金额", color: "#1677ff"},
    {field: "pay_amount", label: "实付金额", color: "#52c41a"},
    {field: "refund_amount", label: "退款金额", color: "#ff4d4f"},
];

const hasValue = (value: any) => value !== undefined && value !== null && value !== "";

const enumTag = (data: any, colors: Record<string, string>, value: any) => (
    <Tag color={colors?.[String(value)] || "default"}>{getEnumText(data, value)}</Tag>
);

const useStyles: any = createStyles((): any => ({
    page: {
        width: "100%",
        maxWidth: 1440,
        margin: "0 auto",
    },
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

export default (): React.ReactNode => {
    const {styles} = useStyles();
    const infoRef = useRef<InfoRef>({} as InfoRef);
    const [enumData, enumLoading] = useEnum(["pay.order", "pay.order_log", "pay.pay"]) as any;
    const orderEnums = enumData?.["pay.order"] || {};
    const logEnums = enumData?.["pay.order_log"] || {};
    const payEnums = enumData?.["pay.pay"] || {};
    const orderType = orderEnums?.type || {};
    const orderStatus = orderEnums?.status || {};
    const payStatus = orderEnums?.pay_status || {};
    const logTypes = logEnums?.type || {};
    const operatorTypes = logEnums?.operator_type || {};
    const payTypes = payEnums?.type || {};

    const formatTime = (value: any) => {
        const text = value === undefined || value === null ? "" : String(value);
        return !text || text.startsWith("0001-01-01") ? "-" : text;
    };

    const renderRows = (items: {
        label: string;
        value: any;
        copy?: any;
        className?: string;
        style?: React.CSSProperties;
    }[]) => (
        <Space orientation="vertical" size={2} className={styles.infoColumn}>
            {items.map((item) => (
                <div className={styles.infoRow} key={item.label}>
                    {item.label ? (
                        <Typography.Text type="secondary" className={styles.labelText}>{item.label}</Typography.Text>
                    ) : null}
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
        {label: "订单类型:", value: enumTag(orderType, typeColors, record?.type)},
        {label: "订单名称:", value: record?.name, className: styles.strongText},
    ]);

    const renderStatusCell = (record: any) => renderRows([
        {label: "订单状态:", value: enumTag(orderStatus, statusColors, record?.status)},
        {label: "支付状态:", value: enumTag(payStatus, payStatusColors, record?.pay_status)},
    ]);

    const renderMoneyCell = (record: any) => renderRows(
        moneyItems.map((item) => ({
            label: `${item.label}:`,
            value: formatMoney(record?.[item.field]),
            className: styles.moneyText,
            style: {color: item.color},
        })),
    );

    const renderTimeCell = (record: any) => renderRows([
        {label: "创建:", value: formatTime(record?.create_time)},
        {label: "支付:", value: formatTime(record?.pay_time)},
        {label: "完成:", value: formatTime(record?.complete_time)},
    ]);

    const columns: any[] = [
        {
            title: "编号信息",
            dataIndex: "id",
            tip: "订单ID",
            width: 165,
            hideInSearch: true,
            render: (_: any, record: any) => renderRows([
                {label: "", value: record?.id ? `${record.id}` : "-", copy: record?.id, className: styles.strongText},
            ]),
        },
        {
            title: "订单ID",
            dataIndex: "id",
            tip: "订单ID",
            valueType: "text",
            hideInTable: true,
        },
        {
            title: "订单类型",
            dataIndex: "type",
            tip: "订单类型",
            valueEnum: buildValueEnum(orderType, typeColors),
            hideInTable: true,
        },
        {
            title: "订单状态",
            dataIndex: "status",
            tip: "订单状态",
            valueEnum: buildValueEnum(orderStatus, statusColors),
            hideInTable: true,
        },
        {
            title: "支付状态",
            dataIndex: "pay_status",
            tip: "支付状态",
            valueEnum: buildValueEnum(payStatus, payStatusColors),
            hideInTable: true,
        },
        {
            title: "支付方式",
            dataIndex: "pay_type",
            tip: "支付方式",
            valueEnum: buildValueEnum(payTypes, payTypeColors),
            hideInTable: true,
        },
        {
            title: "订单信息",
            dataIndex: "order_info",
            tip: "订单名称/订单类型",
            width: 220,
            hideInSearch: true,
            render: (_: any, record: any) => renderOrderCell(record),
        },
        {
            title: "支付信息",
            dataIndex: "pay_info",
            tip: "支付方式",
            width: 88,
            hideInSearch: true,
            render: (_: any, record: any) => <PayType data={payTypes} value={record?.pay_type}/>,
        },
        {
            title: "状态信息",
            dataIndex: "status_info",
            tip: "订单状态/支付状态",
            width: 145,
            hideInSearch: true,
            render: (_: any, record: any) => renderStatusCell(record),
        },
        {
            title: "金额信息",
            dataIndex: "amount",
            tip: "订单/支付/退款金额",
            width: 145,
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
            title: "相关时间",
            dataIndex: "time_info",
            tip: "创建/支付/完成时间",
            width: 220,
            hideInSearch: true,
            render: (_: any, record: any) => renderTimeCell(record),
        },
        {
            title: "操作",
            valueType: "option",
            key: "option",
            width: 88,
            fixed: "right",
            hideInSearch: true,
            render: (_: any, record: any) => {
                const items = [
                    {
                        key: "info",
                        label: "详情",
                        onClick: () => infoRef.current?.open(record, "info"),
                    },
                    {
                        key: "log",
                        label: "日志",
                        onClick: () => infoRef.current?.open(record, "log"),
                    },
                ];
                return (
                    <Dropdown menu={{items}} trigger={["click"]} placement="bottom" arrow>
                        <Button size="small">操作</Button>
                    </Dropdown>
                );
            },
        },
    ];

    return (
        <Body loading={enumLoading}>
            <div className={styles.page}>
                <ProTable
                    rowKey="id"
                    title={<Title>订单管理</Title>}
                    extraRefreshBtn={true}
                    search={true}
                    columns={columns}
                    defaultPage={1}
                    defaultPageSize={10}
                    request={(params: any, sort: any) => getData(params, sort, getOrderList)}
                    paginationAffix={true}
                    {...({
                        tableProps: {
                            scroll: {x: "max-content"},
                        },
                    } as any)}
                />
            </div>
            <Info
                ref={infoRef}
                orderType={orderType}
                orderStatus={orderStatus}
                payStatus={payStatus}
                logTypes={logTypes}
                operatorTypes={operatorTypes}
                payTypes={payTypes}
                typeColors={typeColors}
                statusColors={statusColors}
                payStatusColors={payStatusColors}
                logTypeColors={logTypeColors}
            />
        </Body>
    );
};
