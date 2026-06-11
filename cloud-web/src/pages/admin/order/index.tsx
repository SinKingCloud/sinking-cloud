import React, {useRef} from "react";
import {App, Button, Dropdown, Result, Space, Tag, Typography} from "antd";
import {AlipayOutlined, RedEnvelopeOutlined, WechatOutlined} from "@ant-design/icons";
import {Body, ProTable, ProTableRef, Title} from "sinking-antd";
import {createStyles} from "antd-style";
import {buildValueEnum, dateRangeTransform, emptyTableData, getData} from "@/utils/page";
import useControl from "@/utils/control";
import useEnum from "@/utils/enum";
import {AdminControlCode} from "@/utils/permission";
import {formatMoney} from "@/utils/string";
import {getOrderList} from "@/service/admin/pay";
import Detail, {type DetailRef} from "./components/detail";
import Operate, {type OperateRef} from "./components/operate";
import Clear, {type ClearRef} from "./components/clear";
import Log, {type LogRef} from "./components/log";

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
    payTypeText: {
        fontSize: 12,
        lineHeight: "18px",
    },
    walletIcon: {
        color: "#faad14",
        fontSize: 14,
    },
    alipayIcon: {
        color: "#1677ff",
        fontSize: 14,
    },
    wechatIcon: {
        color: "#07c160",
        fontSize: 14,
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
    const detailRef = useRef<DetailRef>({} as DetailRef);
    const deliverRef = useRef<OperateRef>({} as OperateRef);
    const refoundRef = useRef<OperateRef>({} as OperateRef);
    const clearRef = useRef<ClearRef>({} as ClearRef);
    const logRef = useRef<LogRef>({} as LogRef);

    const hasListControl = control.hasControl(AdminControlCode.WEB_ORDER_LIST);
    const hasInfoControl = control.hasControl(AdminControlCode.WEB_ORDER_INFO);
    const hasDeliverControl = control.hasControl(AdminControlCode.WEB_ORDER_DELIVER);
    const hasRefoundControl = control.hasControl(AdminControlCode.WEB_ORDER_REFOUND);
    const hasClearControl = control.hasControl(AdminControlCode.WEB_ORDER_CLEAR);
    const hasOperateControl = hasInfoControl || hasDeliverControl || hasRefoundControl;

    const orderEnums = enumData?.["pay.order"] || {};
    const logEnums = enumData?.["pay.order_log"] || {};
    const channelEnums = enumData?.["pay.channel"] || {};
    const payEnums = enumData?.["pay.pay"] || {};
    const orderType = orderEnums?.type || {};
    const orderStatus = orderEnums?.status || {};
    const payStatus = orderEnums?.pay_status || {};
    const channelNames = channelEnums?.name || {};
    const payTypes = payEnums?.type || {};
    const typeColors: Record<string, string> = {"0": "blue", "1": "purple"};
    const statusColors: Record<string, string> = {"0": "orange", "1": "blue", "2": "green", "3": "red"};
    const payStatusColors: Record<string, string> = {"0": "orange", "1": "green"};
    const typeColor = (value: any) => typeColors?.[String(value)] || "default";
    const statusColor = (value: any) => statusColors?.[String(value)] || "default";
    const payStatusColor = (value: any) => payStatusColors?.[String(value)] || "default";
    const typeValueEnum = buildValueEnum(orderType, typeColors);
    const statusValueEnum = buildValueEnum(orderStatus, statusColors);
    const payStatusValueEnum = buildValueEnum(payStatus, payStatusColors);
    const canDeliver = (record: any) => [0, 1].includes(numberValue(record?.status));
    const canRefound = (record: any) => (
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

    const renderRows = (items: { label: string; value: any; copy?: any; className?: string; style?: React.CSSProperties }[]) => (
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

    const renderPayType = (value: any) => {
        const text = enumText(payTypes, value);
        const key = String(value || "").toLowerCase();
        const label = String(text || "").toLowerCase();
        if (key === "alipay" || key === "1" || label.includes("支付宝")) {
            return (
                <Space size={4}>
                    <AlipayOutlined className={styles.alipayIcon}/>
                    <Typography.Text className={styles.payTypeText} style={{color: "#1677ff"}}>{text}</Typography.Text>
                </Space>
            );
        }
        if (key === "wxpay" || key === "wechat" || key === "2" || label.includes("微信")) {
            return (
                <Space size={4}>
                    <WechatOutlined className={styles.wechatIcon}/>
                    <Typography.Text className={styles.payTypeText} style={{color: "#07c160"}}>{text}</Typography.Text>
                </Space>
            );
        }
        if (key === "balance" || key === "0" || label.includes("余额")) {
            return (
                <Space size={4}>
                    <RedEnvelopeOutlined className={styles.walletIcon}/>
                    <Typography.Text className={styles.payTypeText} style={{color: "#faad14"}}>{text}</Typography.Text>
                </Space>
            );
        }
        return <Typography.Text className={styles.payTypeText}>{text}</Typography.Text>;
    };

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
        {label: "支付方式:", value: renderPayType(record?.pay_type)},
        {label: "支付通道:", value: hasValue(record?.channel_id) ? enumText(channelNames, record?.channel_id) : "-"},
        {label: "通道ID:", value: record?.channel_id, copy: record?.channel_id},
    ]);

    const renderStatusCell = (record: any) => renderRows([
        {label: "订单状态:", value: <Tag color={statusColor(record?.status)}>{enumText(orderStatus, record?.status)}</Tag>},
        {label: "支付状态:", value: <Tag color={payStatusColor(record?.pay_status)}>{enumText(payStatus, record?.pay_status)}</Tag>},
    ]);

    const renderMoneyCell = (record: any) => renderRows([
        {label: "订单金额:", value: formatMoney(record?.amount), className: styles.moneyText, style: {color: "#1677ff"}},
        {label: "实付金额:", value: formatMoney(record?.pay_amount), className: styles.moneyText, style: {color: "#52c41a"}},
        {label: "退款金额:", value: formatMoney(record?.refund_amount), className: styles.moneyText, style: {color: "#ff4d4f"}},
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
        ...(hasOperateControl ? [{
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
                        onClick: () => detailRef.current?.open(record?.id),
                    } : null,
                    hasInfoControl ? {
                        key: "log",
                        label: "日志",
                        onClick: () => logRef.current?.open(record),
                    } : null,
                    hasDeliverControl ? {
                        key: "deliver",
                        label: "发货",
                        disabled: !canDeliver(record),
                        onClick: () => deliverRef.current?.open(record),
                    } : null,
                    hasRefoundControl ? {
                        key: "refound",
                        label: "退款",
                        disabled: !canRefound(record),
                        onClick: () => refoundRef.current?.open(record),
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
                <Detail
                    ref={detailRef}
                    orderEnums={orderEnums}
                    typeColors={typeColors}
                    statusColors={statusColors}
                    payStatusColors={payStatusColors}
                />
                <Operate ref={deliverRef} type="deliver" orderEnums={orderEnums} typeColors={typeColors} onSuccess={refresh}/>
                <Operate ref={refoundRef} type="refound" orderEnums={orderEnums} typeColors={typeColors} onSuccess={refresh}/>
                <Clear ref={clearRef} onSuccess={refresh}/>
                <Log ref={logRef} logEnums={logEnums}/>
            </>}
        </Body>
    );
};
