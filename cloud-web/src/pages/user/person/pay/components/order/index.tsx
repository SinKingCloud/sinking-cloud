import React, {useRef} from "react";
import {App, Button, Dropdown, Tag, Typography} from "antd";
import {createStyles} from "antd-style";
import type {PageTableProps} from "sinking-antd";
import {PageTable} from "sinking-antd";
import useEnum from "@/utils/enum";
import {getParams} from "@/utils/page";
import {formatMoney} from "@/utils/string";
import {getOrderList} from "@/service/user/pay";
import Info, {type InfoRef} from "./components/info";
import PayType, {getEnumText} from "./components/pay-type";

const typeColors: Record<string, string> = {"0": "default", "1": "cyan"};
const statusColors: Record<string, string> = {"0": "orange", "1": "blue", "2": "green", "3": "default"};
const payStatusColors: Record<string, string> = {"0": "default", "1": "green"};
const logTypeColors: Record<string, string> = {"0": "blue", "1": "green", "2": "red"};
const useStyles = createStyles(({css, token}) => ({
    table: css`
        .order-cell { display: flex; flex-direction: column; align-items: flex-start; gap: 3px; min-width: 0; color: ${token.colorTextSecondary}; font-size: 12px; line-height: 20px; }
        .order-name.ant-typography { max-width: 100%; color: ${token.colorText}; font-size: 13px; line-height: 20px; }
        .order-id.ant-typography { max-width: 100%; color: ${token.colorTextSecondary}; font-size: 12px; line-height: 20px; }
        .order-label { color: ${token.colorTextTertiary}; }
        .order-value { display: flex; align-items: center; gap: 6px; white-space: nowrap; font-variant-numeric: tabular-nums; }
        .order-amount { color: ${token.colorText}; }
        .order-cell .ant-tag { margin-inline-end: 0; font-size: 12px; line-height: 18px; }
        .ant-typography-copy { margin-inline-start: 6px; color: ${token.colorTextTertiary}; font-size: 12px; }
        .ant-typography-copy:hover { color: ${token.colorPrimary}; }
    `,
}));

export default ({hero}: {hero: React.ReactElement}): React.ReactNode => {
    const {message} = App.useApp();
    const {styles} = useStyles();
    const infoRef = useRef<InfoRef>(null);
    const [enumData, enumLoading] = useEnum(["pay.order", "pay.order_log", "pay.pay"]);
    const orderEnums = enumData["pay.order"] || {};
    const logEnums = enumData["pay.order_log"] || {};
    const orderType = orderEnums.type || {};
    const orderStatus = orderEnums.status || {};
    const payStatus = orderEnums.pay_status || {};
    const logTypes = logEnums.type || {};
    const operatorTypes = logEnums.operator_type || {};
    const payTypes = enumData["pay.pay"]?.type || {};
    const filterFields = [
        {name: "type", label: "订单类型", options: orderType},
        {name: "status", label: "订单状态", options: orderStatus},
        {name: "pay_status", label: "支付状态", options: payStatus},
        {name: "pay_type", label: "支付方式", options: payTypes},
    ];
    const dateFields = [{name: "create_time", label: "创建时间"}, {name: "pay_time", label: "支付时间"}, {name: "complete_time", label: "完成时间"}];

    const columns: PageTableProps<any>["columns"] = [
        {
            title: "订单信息", dataIndex: "id", width: 260, sorter: true,
            render: (_: unknown, record: any) => <div className="order-cell">
                <Typography.Text className="order-name" ellipsis={{tooltip: record.name || "-"}}>{record.name || "-"}</Typography.Text>
                <Typography.Text className="order-id" copyable={record.id ? {text: String(record.id)} : false}>
                    {record.id ? String(record.id) : "-"}
                </Typography.Text>
            </div>,
        },
        {
            title: "金额", dataIndex: "amount", width: 155,
            render: (_: unknown, record: any) => <div className="order-cell">
                <div className="order-value"><span className="order-label">应付</span><span className="order-amount">¥ {formatMoney(record.amount)}</span></div>
                <div className="order-value"><span className="order-label">实付</span><span>¥ {formatMoney(record.pay_amount)}</span></div>
                {Number(record.refund_amount) > 0 && <div className="order-value"><span className="order-label">退款</span><span>¥ {formatMoney(record.refund_amount)}</span></div>}
            </div>,
        },
        {
            title: "支付方式", dataIndex: "pay_type", width: 120,
            render: (_: unknown, record: any) => <div className="order-cell">
                <PayType data={payTypes} value={record.pay_type}/>
                <span className="order-label">{getEnumText(orderType, record.type)}</span>
            </div>,
        },
        {
            title: "状态", dataIndex: "status", width: 110,
            render: (_: unknown, record: any) => <div className="order-cell">
                <Tag color={statusColors[String(record.status)] || "default"}>{getEnumText(orderStatus, record.status)}</Tag>
                <span>{getEnumText(payStatus, record.pay_status)}</span>
            </div>,
        },
        {
            title: "相关时间", dataIndex: "create_time", width: 225, sorter: true,
            render: (_: unknown, record: any) => <div className="order-cell">
                {dateFields.map((item) => <div className="order-value" key={item.name}>
                    <span className="order-label">{item.label.slice(0, 2)}于:</span>
                    <span>{record[item.name] && !String(record[item.name]).startsWith("0001-01-01") ? record[item.name] : "-"}</span>
                </div>)}
            </div>,
        },
        {
            title: "操作", key: "option", width: 80, align: "center", fixed: "right", className: "action-cell",
            render: (_: unknown, record: any) => <Dropdown trigger={["click"]} placement="bottom" arrow={{pointAtCenter: true}} menu={{items: [
                {key: "info", label: "详情", onClick: () => infoRef.current?.open(record, "info")},
                {key: "log", label: "日志", onClick: () => infoRef.current?.open(record, "log")},
            ]}}>
                <Button size="small" aria-label="订单操作">操作</Button>
            </Dropdown>,
        },
    ];

    return <>
        <PageTable
            ariaLabel="订单记录"
            rootClassName={styles.table}
            hero={hero}
            search={{name: "id", ariaLabel: "搜索订单编号", placeholder: "搜索订单编号", maxLength: 19}}
            filters={[
                ...filterFields.map((item) => ({
                    type: "select" as const, name: item.name, label: item.label, valueEnum: item.options,
                    allLabel: `全部${item.label.slice(2)}`,
                    showSelectedLabel: false,
                    disabled: enumLoading,
                })),
                ...dateFields.map((item) => ({
                    type: "date" as const, name: item.name, label: item.label, showSelectedLabel: false,
                })),
            ]}
            toolbar={{refresh: {ariaLabel: "刷新订单记录"}}}
            rowKey="id"
            columns={columns}
            scroll={{x: 950}}
            defaultPageSize={10}
            loading={enumLoading || undefined}
            requestEnabled={!enumLoading}
            request={async (params, sort) => {
                const response = await getOrderList({body: getParams(params, sort)});
                if (response?.code !== 200) throw new Error(response?.message || "获取订单记录失败");
                return {data: response.data?.list, total: response.data?.total ?? 0, success: true};
            }}
            onRequestError={(error) => message.error(error.message)}
            pagination={{unit: "条"}}
            paginationAffix={{offsetBottom: 15}}
        />
        <Info ref={infoRef} orderType={orderType} orderStatus={orderStatus} payStatus={payStatus}
              logTypes={logTypes} operatorTypes={operatorTypes} payTypes={payTypes} typeColors={typeColors}
              statusColors={statusColors} payStatusColors={payStatusColors} logTypeColors={logTypeColors}/>
    </>;
};
