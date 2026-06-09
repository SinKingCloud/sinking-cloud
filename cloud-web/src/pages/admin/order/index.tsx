import React, {forwardRef, useImperativeHandle, useRef, useState} from "react";
import {
    App,
    Button,
    Descriptions,
    Dropdown,
    Empty,
    Form,
    Input,
    InputNumber,
    Pagination,
    Result,
    Segmented,
    Select,
    Spin,
    Tag,
    Typography,
} from "antd";
import {Body, ProModal, ProModalRef, ProTable, ProTableRef, Title, useTheme} from "sinking-antd";
import {createStyles} from "antd-style";
import {buildOptions, buildValueEnum, dateRangeTransform, emptyTableData, getData} from "@/utils/page";
import useControl from "@/utils/control";
import useEnum from "@/utils/enum";
import {AdminControlCode} from "@/utils/permission";
import {formatMoney} from "@/utils/string";
import {
    clearOrder,
    deliverOrder,
    getOrderInfo,
    getOrderList,
    getOrderLogInfo,
    getOrderLogList,
    refoundOrder,
} from "@/service/admin/pay";

type OrderActionRef = {
    open: (record?: any) => void;
};

type DetailRef = {
    open: (id: any) => void;
};

type OrderActionProps = {
    onSuccess?: () => void;
};

type DetailProps = {
    orderEnums: any;
};

type LogProps = {
    logEnums: any;
};

type OrderOperateProps = OrderActionProps & {
    type: "deliver" | "refound";
};

const ORDER_STATUS = {
    UNPAY: 0,
    PAYED: 1,
    REFUNDED: 2,
};
const ORDER_TYPE = {
    RECHARGE: 0,
    BUY: 1,
};
const PAY_STATUS = {
    UNPAY: 0,
    PAYED: 1,
};

const useStyles: any = createStyles(({css, token}: any, props: any = {}) => {
    const isCompactMode = props?.isCompactMode || false;
    const modalGap = isCompactMode ? 10 : 12;
    const cardPadding = isCompactMode ? 10 : 12;

    return {
        orderNo: css`
            min-width: 180px;
        `,
        mainInfo: css`
            min-width: 220px;

            .name {
                max-width: 260px;
                overflow: hidden;
                color: ${token.colorText};
                font-weight: 500;
                line-height: 22px;
                text-overflow: ellipsis;
                white-space: nowrap;
            }

            .sub {
                margin-top: 2px;
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 18px;
            }
        `,
        moneyInfo: css`
            min-width: 150px;
            display: grid;
            gap: 2px;

            .row {
                display: flex;
                align-items: center;
                justify-content: space-between;
                gap: 8px;
            }

            .label {
                flex: none;
                color: ${token.colorTextTertiary};
                font-size: 12px;
            }

            .value {
                font-variant-numeric: tabular-nums;
            }
        `,
        timeInfo: css`
            min-width: 180px;
            display: grid;
            gap: 2px;
            font-size: 12px;
            line-height: 18px;

            .label {
                color: ${token.colorTextTertiary};
            }
        `,
        modalContent: css`
            display: flex;
            flex-direction: column;
            gap: ${modalGap}px;
        `,
        summary: css`
            padding: ${cardPadding}px;
            border: 1px solid ${token.colorBorderSecondary};
            border-radius: ${token.borderRadius}px;
            background: ${token.colorFillQuaternary};

            .title {
                color: ${token.colorText};
                font-weight: 500;
                line-height: 22px;
            }

            .meta {
                margin-top: 4px;
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 18px;
                word-break: break-all;
            }
        `,
        infoLoading: css`
            min-height: 220px;
            display: flex;
            align-items: center;
            justify-content: center;
        `,
        detailGrid: css`
            .ant-descriptions-item-label {
                width: 108px;
            }
        `,
        operateForm: css`
            .ant-form-item {
                margin-bottom: ${isCompactMode ? 12 : 16}px;
            }

            .tips {
                margin-bottom: ${isCompactMode ? 10 : 12}px;
                padding: ${cardPadding}px;
                color: ${token.colorTextSecondary};
                font-size: 13px;
                line-height: 20px;
                border: 1px solid ${token.colorBorderSecondary};
                border-radius: ${token.borderRadius}px;
                background: ${token.colorFillQuaternary};
            }
        `,
        logModal: css`
            .panel {
                display: flex;
                flex-direction: column;
                gap: ${modalGap}px;
            }

            .log-filter {
                display: grid;
                grid-template-columns: repeat(3, minmax(0, 1fr));
                gap: 8px;
            }

            .log-list {
                min-height: 260px;
            }

            .log-card {
                padding: ${cardPadding}px;
                border: 1px solid ${token.colorBorderSecondary};
                border-radius: ${token.borderRadius}px;
                background: ${token.colorBgContainer};
            }

            .log-card + .log-card {
                margin-top: 8px;
            }

            .log-head,
            .log-foot {
                display: flex;
                align-items: center;
                justify-content: space-between;
                gap: 8px;
            }

            .log-content {
                margin-top: 6px;
                color: ${token.colorTextSecondary};
                font-size: 13px;
                line-height: 20px;
                word-break: break-all;
            }

            .log-foot {
                margin-top: 8px;
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 18px;
            }

            .pager {
                display: flex;
                align-items: center;
                justify-content: space-between;
                gap: 8px;
                color: ${token.colorTextTertiary};
                font-size: 12px;
            }

            @media (max-width: 575px) {
                .log-filter {
                    grid-template-columns: 1fr;
                }

                .log-foot,
                .pager {
                    align-items: flex-start;
                    flex-direction: column;
                }
            }
        `,
    };
});

const enumText = (data: any, value: any) => data?.[String(value)] || data?.[value] || "-";
const numberValue = (value: any) => Number(value || 0);
const canDeliver = (record: any) => numberValue(record?.status) === ORDER_STATUS.UNPAY;
const canRefound = (record: any) => (
    numberValue(record?.type) !== ORDER_TYPE.RECHARGE
    && numberValue(record?.pay_status) === PAY_STATUS.PAYED
    && numberValue(record?.pay_amount) > numberValue(record?.refund_amount)
);

const moneyRows = (record: any) => [
    ["订单", record?.amount],
    ["支付", record?.pay_amount],
    ["退款", record?.refund_amount],
];

const TimeRows = ({record}: any) => (
    <>
        <div><span className="label">创建: </span>{record?.create_time || "-"}</div>
        <div><span className="label">支付: </span>{record?.pay_time || "-"}</div>
        <div><span className="label">完成: </span>{record?.complete_time || "-"}</div>
    </>
);

export default (): React.ReactNode => {
    const {message} = App.useApp();
    const theme = useTheme();
    const isCompactMode = theme?.isCompactTheme?.() || false;
    const {styles} = useStyles({isCompactMode});
    const control = useControl();
    const [enumData, enumLoading] = useEnum(["pay.order", "pay.order_log"]) as any;
    const tableRef = useRef<ProTableRef>({} as ProTableRef);
    const detailRef = useRef<DetailRef>({} as DetailRef);
    const deliverRef = useRef<OrderActionRef>({} as OrderActionRef);
    const refoundRef = useRef<OrderActionRef>({} as OrderActionRef);
    const clearRef = useRef<OrderActionRef>({} as OrderActionRef);
    const logRef = useRef<DetailRef>({} as DetailRef);

    const hasListControl = control.hasControl(AdminControlCode.WEB_ORDER_LIST);
    const hasInfoControl = control.hasControl(AdminControlCode.WEB_ORDER_INFO);
    const hasDeliverControl = control.hasControl(AdminControlCode.WEB_ORDER_DELIVER);
    const hasRefoundControl = control.hasControl(AdminControlCode.WEB_ORDER_REFOUND);
    const hasClearControl = control.hasControl(AdminControlCode.WEB_ORDER_CLEAR);

    const orderEnums = enumData?.["pay.order"] || {};
    const logEnums = enumData?.["pay.order_log"] || {};
    const orderType = orderEnums?.type || {};
    const orderStatus = orderEnums?.status || {};
    const payStatus = orderEnums?.pay_status || {};
    const typeValueEnum = buildValueEnum(orderType, {"0": "blue", "1": "purple"});
    const statusValueEnum = buildValueEnum(orderStatus, {"0": "orange", "1": "green", "2": "red"});
    const payStatusValueEnum = buildValueEnum(payStatus, {"0": "orange", "1": "green"});
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

    const columns: any[] = [
        {
            title: "ID",
            dataIndex: "id",
            tip: "订单ID",
            width: 80,
            sorter: true,
            hideInSearch: true,
        },
        {
            title: "订单信息",
            dataIndex: "trade_no",
            tip: "订单号",
            valueType: "text",
            copyable: true,
            render: (_: any, record: any) => (
                <div className={styles.mainInfo}>
                    <div className="name">{record?.name || "-"}</div>
                    <Typography.Text className={styles.orderNo} copyable>{record?.trade_no || "-"}</Typography.Text>
                    <div className="sub">外部订单: {record?.out_trade_no || "-"}</div>
                </div>
            ),
        },
        {
            title: "用户ID",
            dataIndex: "user_id",
            tip: "用户ID",
            valueType: "text",
            width: 110,
            sorter: true,
            copyable: true,
        },
        {
            title: "通道ID",
            dataIndex: "channel_id",
            tip: "通道ID",
            valueType: "text",
            width: 110,
            sorter: true,
            copyable: true,
        },
        {
            title: "订单类型",
            dataIndex: "type",
            tip: "订单类型",
            valueEnum: typeValueEnum,
            sorter: true,
        },
        {
            title: "订单状态",
            dataIndex: "status",
            tip: "订单状态",
            valueEnum: statusValueEnum,
            sorter: true,
        },
        {
            title: "支付状态",
            dataIndex: "pay_status",
            tip: "支付状态",
            valueEnum: payStatusValueEnum,
            sorter: true,
        },
        {
            title: "金额信息",
            dataIndex: "amount",
            tip: "订单/支付/退款金额",
            hideInSearch: true,
            render: (_: any, record: any) => (
                <div className={styles.moneyInfo}>
                    {moneyRows(record).map(([label, value]) => (
                        <div className="row" key={label}>
                            <span className="label">{label}</span>
                            <Typography.Text className="value">{formatMoney(value)}</Typography.Text>
                        </div>
                    ))}
                </div>
            ),
        },
        {
            title: "外部订单号",
            dataIndex: "out_trade_no",
            tip: "外部订单号",
            valueType: "text",
            copyable: true,
            hideInTable: true,
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
            hideInSearch: true,
            render: (_: any, record: any) => (
                <div className={styles.timeInfo}>
                    <TimeRows record={record}/>
                </div>
            ),
        },
        ...(hasInfoControl || hasDeliverControl || hasRefoundControl ? [{
            title: "操作",
            valueType: "option",
            key: "option",
            width: 128,
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
                    defaultPageSize={10}
                    extra={hasClearControl ? <Button onClick={openClear}>清理订单</Button> : undefined}
                    request={(params, sort) => {
                        if (!hasListControl) {
                            return Promise.resolve(emptyTableData);
                        }
                        return getData(params, sort, getOrderList);
                    }}
                    paginationAffix={true}
                    selectionAffix={true}
                />
                <OrderDetail ref={detailRef} orderEnums={orderEnums}/>
                <OrderOperate ref={deliverRef} type="deliver" onSuccess={refresh}/>
                <OrderOperate ref={refoundRef} type="refound" onSuccess={refresh}/>
                <OrderClear ref={clearRef} onSuccess={refresh}/>
                <OrderLog ref={logRef} logEnums={logEnums}/>
            </>}
        </Body>
    );
};

const OrderSummary = ({record, orderEnums}: any) => (
    <div>
        <div className="title">{record?.name || "订单"}</div>
        <div className="meta">订单号: {record?.trade_no || "-"}</div>
        <div className="meta">
            {enumText(orderEnums?.type, record?.type)} / {enumText(orderEnums?.status, record?.status)} / {formatMoney(record?.amount)}
        </div>
    </div>
);

const OrderDetail = forwardRef<DetailRef, DetailProps>(({orderEnums}, ref) => {
    const {message} = App.useApp();
    const theme = useTheme();
    const {styles} = useStyles({isCompactMode: theme?.isCompactTheme?.() || false});
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
            getOrderInfo({
                body: {id: orderId},
                onSuccess: (res: any) => setInfo(res?.data || {}),
                onFail: (res: any) => {
                    message?.error(res?.message || "获取订单详情失败");
                    modalRef.current?.hide();
                },
            }).finally(() => setLoading(false));
        },
    }));

    return (
        <ProModal
            ref={modalRef}
            title={<Title>订单详情</Title>}
            width={560}
            modalProps={{
                footer: null,
                mask: {closable: true},
                afterClose: () => {
                    setInfo(undefined);
                    setLoading(false);
                },
            } as any}>
            {loading ? (
                <div className={styles.infoLoading}>
                    <Spin/>
                </div>
            ) : (
                <div className={styles.modalContent}>
                    <div className={styles.summary}>
                        <OrderSummary record={info} orderEnums={orderEnums}/>
                    </div>
                    <Descriptions className={styles.detailGrid} column={1} size="small" bordered>
                        <Descriptions.Item label="订单ID">{info?.id || "-"}</Descriptions.Item>
                        <Descriptions.Item label="用户ID">{info?.user_id || "-"}</Descriptions.Item>
                        <Descriptions.Item label="通道ID">{info?.channel_id || "-"}</Descriptions.Item>
                        <Descriptions.Item label="订单类型">{enumText(orderEnums?.type, info?.type)}</Descriptions.Item>
                        <Descriptions.Item label="订单状态">{enumText(orderEnums?.status, info?.status)}</Descriptions.Item>
                        <Descriptions.Item label="支付状态">{enumText(orderEnums?.pay_status, info?.pay_status)}</Descriptions.Item>
                        <Descriptions.Item label="订单金额">{formatMoney(info?.amount)}</Descriptions.Item>
                        <Descriptions.Item label="支付金额">{formatMoney(info?.pay_amount)}</Descriptions.Item>
                        <Descriptions.Item label="退款金额">{formatMoney(info?.refund_amount)}</Descriptions.Item>
                        <Descriptions.Item label="订单号">{info?.trade_no || "-"}</Descriptions.Item>
                        <Descriptions.Item label="外部订单号">{info?.out_trade_no || "-"}</Descriptions.Item>
                        <Descriptions.Item label="创建时间">{info?.create_time || "-"}</Descriptions.Item>
                        <Descriptions.Item label="支付时间">{info?.pay_time || "-"}</Descriptions.Item>
                        <Descriptions.Item label="过期时间">{info?.expired_time || "-"}</Descriptions.Item>
                        <Descriptions.Item label="完成时间">{info?.complete_time || "-"}</Descriptions.Item>
                        <Descriptions.Item label="更新时间">{info?.update_time || "-"}</Descriptions.Item>
                    </Descriptions>
                </div>
            )}
        </ProModal>
    );
});

const OrderOperate = forwardRef<OrderActionRef, OrderOperateProps>(({type, onSuccess}, ref) => {
    const {message} = App.useApp();
    const theme = useTheme();
    const {styles} = useStyles({isCompactMode: theme?.isCompactTheme?.() || false});
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
            setRecord(data || {});
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
            body: isRefound ? {id, money: values?.money || 0} : {id},
            onSuccess: (res: any) => {
                modalRef.current?.hide();
                onSuccess?.();
                message?.success(res?.message || "操作成功");
            },
            onFail: (res: any) => message?.error(res?.message || "操作失败"),
        }).finally(() => setSubmitting(false));
    };

    return (
        <ProModal
            ref={modalRef}
            title={<Title>{title}</Title>}
            onOk={form.submit}
            width={420}
            modalProps={{
                confirmLoading: submitting,
                okText: "提交",
                cancelText: "取消",
                forceRender: true,
                mask: {closable: true},
                afterClose: () => {
                    setRecord(undefined);
                    setSubmitting(false);
                    form.resetFields();
                },
            } as any}>
            <Form className={styles.operateForm} form={form} layout="vertical" onFinish={submit}>
                <div className="tips">
                    {isRefound
                        ? `可退款金额 ${formatMoney(refundableMoney)}，不填或填 0 将按后端规则执行全额可退款金额。`
                        : "发货会触发订单回调，成功后订单会进入已支付/已完成状态。"}
                </div>
                <Form.Item label="订单号">
                    <Input value={record?.trade_no || "-"} disabled/>
                </Form.Item>
                {isRefound && (
                    <Form.Item name="money" label="退款金额">
                        <InputNumber
                            min={0}
                            max={refundableMoney || undefined}
                            precision={2}
                            controls={false}
                            style={{width: "100%"}}
                            placeholder="请输入退款金额，0 表示按后端规则退款"
                        />
                    </Form.Item>
                )}
            </Form>
        </ProModal>
    );
});

const OrderClear = forwardRef<OrderActionRef, OrderActionProps>(({onSuccess}, ref) => {
    const {message} = App.useApp();
    const theme = useTheme();
    const {styles} = useStyles({isCompactMode: theme?.isCompactTheme?.() || false});
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const [form] = Form.useForm();
    const [submitting, setSubmitting] = useState(false);

    useImperativeHandle(ref, () => ({
        open: () => {
            form.resetFields();
            form.setFieldsValue({day: 30});
            modalRef.current?.show();
        },
    }));

    const submit = async (values: any) => {
        setSubmitting(true);
        await clearOrder({
            body: {day: Number(values?.day || 1)},
            onSuccess: (res: any) => {
                modalRef.current?.hide();
                onSuccess?.();
                message?.success(res?.message || "清理成功");
            },
            onFail: (res: any) => message?.error(res?.message || "清理失败"),
        }).finally(() => setSubmitting(false));
    };

    return (
        <ProModal
            ref={modalRef}
            title={<Title>清理订单</Title>}
            onOk={form.submit}
            width={420}
            modalProps={{
                confirmLoading: submitting,
                okText: "清理",
                cancelText: "取消",
                forceRender: true,
                mask: {closable: true},
                afterClose: () => {
                    form.resetFields();
                    setSubmitting(false);
                },
            } as any}>
            <Form className={styles.operateForm} form={form} layout="vertical" onFinish={submit}>
                <div className="tips">仅清理超过保留天数的过期未支付订单，并同步删除相关订单日志。</div>
                <Form.Item name="day" label="保留天数" rules={[{required: true, message: "请输入保留天数"}]}>
                    <InputNumber min={1} precision={0} style={{width: "100%"}} placeholder="请输入保留天数"/>
                </Form.Item>
            </Form>
        </ProModal>
    );
});

const OrderLog = forwardRef<DetailRef, LogProps>(({logEnums}, ref) => {
    const {message} = App.useApp();
    const theme = useTheme();
    const {styles} = useStyles({isCompactMode: theme?.isCompactTheme?.() || false});
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const [filters, setFilters] = useState<any>({});
    const [order, setOrder] = useState<any>();
    const [detail, setDetail] = useState<any>();
    const [active, setActive] = useState<"list" | "detail">("list");
    const [list, setList] = useState<any[]>([]);
    const [page, setPage] = useState(1);
    const [total, setTotal] = useState(0);
    const [loading, setLoading] = useState(false);
    const logTypes = logEnums?.type || {};
    const operatorTypes = logEnums?.operator_type || {};

    const loadList = (nextPage = page, nextFilters = filters, nextOrder = order) => {
        const orderId = Number(nextOrder?.id || nextOrder);
        if (!orderId) {
            return;
        }
        setLoading(true);
        getOrderLogList({
            body: {
                order_id: orderId,
                page: nextPage,
                page_size: 6,
                order_by_field: "id",
                order_by_type: "desc",
                ...nextFilters,
            },
            onSuccess: (res: any) => {
                setList(res?.data?.list || []);
                setTotal(Number(res?.data?.total || 0));
            },
            onFail: (res: any) => message?.error(res?.message || "获取订单日志失败"),
        }).finally(() => setLoading(false));
    };

    useImperativeHandle(ref, () => ({
        open: (record: any) => {
            setOrder(record || {});
            setFilters({});
            setPage(1);
            setActive("list");
            setDetail(undefined);
            modalRef.current?.show();
            setTimeout(() => loadList(1, {}, record), 0);
        },
    }));

    const changeFilter = (values: any) => {
        const nextFilters = Object.fromEntries(Object.entries(values || {}).filter(([, value]) => value !== undefined && value !== ""));
        setFilters(nextFilters);
        setPage(1);
        loadList(1, nextFilters);
    };

    const openDetail = (id: any) => {
        const logId = Number(id);
        if (!logId) {
            return;
        }
        setLoading(true);
        getOrderLogInfo({
            body: {id: logId},
            onSuccess: (res: any) => {
                setDetail(res?.data || {});
                setActive("detail");
            },
            onFail: (res: any) => message?.error(res?.message || "获取日志详情失败"),
        }).finally(() => setLoading(false));
    };

    return (
        <ProModal
            ref={modalRef}
            title={<Title>订单日志</Title>}
            width={560}
            modalProps={{
                footer: null,
                mask: {closable: true},
                afterClose: () => {
                    setOrder(undefined);
                    setDetail(undefined);
                    setList([]);
                    setFilters({});
                    setPage(1);
                    setTotal(0);
                },
            } as any}>
            <div className={styles.logModal}>
                <Segmented
                    block
                    value={active}
                    options={[
                        {label: "日志列表", value: "list"},
                        {label: "日志详情", value: "detail", disabled: !detail},
                    ]}
                    onChange={(value: any) => setActive(value)}
                />
                {active === "list" ? (
                    <div className="panel">
                        <div className={styles.summary}>
                            <div className="title">{order?.name || "订单日志"}</div>
                            <div className="meta">订单号: {order?.trade_no || "-"}</div>
                        </div>
                        <Form layout="vertical" onFinish={changeFilter}>
                            <div className="log-filter">
                                <Form.Item name="type" label="日志类型">
                                    <Select allowClear placeholder="全部" options={buildOptions(logTypes, true) as any}/>
                                </Form.Item>
                                <Form.Item name="operator_type" label="操作人类型">
                                    <Select allowClear placeholder="全部" options={buildOptions(operatorTypes, true) as any}/>
                                </Form.Item>
                                <Form.Item label=" ">
                                    <Button htmlType="submit" type="primary" block>筛选</Button>
                                </Form.Item>
                            </div>
                        </Form>
                        <Spin spinning={loading}>
                            <div className="log-list">
                                {list.length <= 0 ? (
                                    <Empty image={Empty.PRESENTED_IMAGE_SIMPLE}/>
                                ) : list.map((item) => (
                                    <div className="log-card" key={item?.id}>
                                        <div className="log-head">
                                            <Tag color="blue">{enumText(logTypes, item?.type)}</Tag>
                                            <Typography.Text type="secondary">{item?.create_time || "-"}</Typography.Text>
                                        </div>
                                        <div className="log-content">{item?.content || "-"}</div>
                                        <div className="log-foot">
                                            <span>{enumText(operatorTypes, item?.operator_type)} #{item?.operator_id || 0}</span>
                                            <Button type="link" size="small" onClick={() => openDetail(item?.id)}>详情</Button>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </Spin>
                        {total > 0 && (
                            <div className="pager">
                                <span>共 {total} 条</span>
                                <Pagination
                                    size="small"
                                    current={page}
                                    pageSize={6}
                                    total={total}
                                    showSizeChanger={false}
                                    onChange={(nextPage) => {
                                        setPage(nextPage);
                                        loadList(nextPage);
                                    }}
                                />
                            </div>
                        )}
                    </div>
                ) : (
                    <Spin spinning={loading}>
                        <Descriptions column={1} size="small" bordered>
                            <Descriptions.Item label="日志ID">{detail?.id || "-"}</Descriptions.Item>
                            <Descriptions.Item label="订单ID">{detail?.order_id || "-"}</Descriptions.Item>
                            <Descriptions.Item label="日志类型">{enumText(logTypes, detail?.type)}</Descriptions.Item>
                            <Descriptions.Item label="操作人">{enumText(operatorTypes, detail?.operator_type)} #{detail?.operator_id || 0}</Descriptions.Item>
                            <Descriptions.Item label="日志内容">{detail?.content || "-"}</Descriptions.Item>
                            <Descriptions.Item label="创建时间">{detail?.create_time || "-"}</Descriptions.Item>
                        </Descriptions>
                    </Spin>
                )}
            </div>
        </ProModal>
    );
});

OrderDetail.displayName = "OrderDetail";
OrderOperate.displayName = "OrderOperate";
OrderClear.displayName = "OrderClear";
OrderLog.displayName = "OrderLog";
