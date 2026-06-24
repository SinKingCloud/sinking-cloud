import React, {useRef, useState} from "react";
import {App, Button, Card, Empty, Row, Space, Tag, Timeline, Tooltip, Typography} from "antd";
import {createStyles} from "antd-style";
import {Body, Icon, ProModal, ProModalRef, ProTable, ProTableRef, Title} from "sinking-antd";
import useEnum from "@/utils/enum";
import {buildValueEnum, dateRangeTransform, getData} from "@/utils/page";
import {formatMoney} from "@/utils/string";
import {getOrderList, getOrderLogList} from "@/service/user/pay";

const DEFAULT_PAGE_SIZE = 10;
const LOG_PAGE_SIZE = 20;
const typeColors: Record<string, string> = {"0": "blue", "1": "purple"};
const statusColors: Record<string, string> = {"0": "orange", "1": "gold", "2": "green", "3": "red"};
const payStatusColors: Record<string, string> = {"0": "orange", "1": "green"};

const enumText = (data: any, value: any) => (
    data?.[String(value)] || data?.[value] || (value === undefined || value === null || value === "" ? "-" : String(value))
);

const useStyles = createStyles(({token, css, isDarkMode}: any) => {
    const compact = token.sizeStep <= 3;
    const radius = token.borderRadiusLG;
    return {
        page: css`
            display: flex;
            flex-direction: column;
            flex-wrap: nowrap;
            gap: ${compact ? 12 : 14}px;
            max-width: 1120px;
            margin: 0 auto;
        `,
        hero: css`
            position: relative;
            overflow: hidden;
            isolation: isolate;
            border-radius: ${radius}px;
            padding: ${compact ? "16px 20px" : "18px 24px"};
            background:
                radial-gradient(circle at 18% 16%, ${isDarkMode ? "rgba(255,255,255,0.11)" : `${token.colorPrimary}26`} 0%, transparent 34%),
                radial-gradient(circle at 82% 42%, ${isDarkMode ? "rgba(255,255,255,0.08)" : "rgba(19,194,194,0.18)"} 0%, transparent 36%),
                linear-gradient(135deg, ${token.colorBgContainer} 0%, ${isDarkMode ? "rgba(255,255,255,0.018)" : token.colorFillQuaternary} 100%);
            box-shadow: ${isDarkMode ? "0 14px 38px rgba(0,0,0,0.24)" : "0 16px 36px rgba(15,23,42,0.065)"};

            &::before {
                content: "";
                position: absolute;
                right: -130px;
                bottom: -160px;
                width: 460px;
                height: 360px;
                border-radius: 50%;
                background:
                    radial-gradient(circle, ${token.colorPrimaryBgHover} 0%, transparent 46%),
                    conic-gradient(from 20deg, transparent 0 14%, ${token.colorPrimaryBorder} 22%, transparent 38%, ${token.colorPrimaryBgHover} 56%, transparent 76%, ${token.colorPrimaryBg} 90%, transparent 100%);
                filter: blur(52px);
                opacity: ${isDarkMode ? 0.16 : 0.32};
                pointer-events: none;
                z-index: 0;
            }

            .hero-watermark {
                position: absolute;
                right: -22px;
                bottom: -84px;
                z-index: 0;
                color: ${token.colorPrimary};
                font-size: ${compact ? 190 : 220}px;
                line-height: 1;
                opacity: ${isDarkMode ? 0.11 : 0.08};
                filter: blur(13px);
                transform: rotate(-18deg);
                pointer-events: none;
            }
        `,
        heroInner: css`
            position: relative;
            z-index: 1;
            display: grid;
            grid-template-columns: minmax(0, 1fr) minmax(230px, 280px);
            align-items: center;
            gap: ${compact ? 12 : 16}px;

            @media (max-width: 760px) {
                grid-template-columns: 1fr;
            }
        `,
        heroTitle: css`
            margin: 0;
            color: ${token.colorText};
            font-size: ${compact ? 26 : 28}px;
            font-weight: 700;
            line-height: ${compact ? 34 : 36}px;
            letter-spacing: 0;
        `,
        heroDesc: css`
            max-width: 600px;
            margin-top: 6px;
            color: ${token.colorTextSecondary};
            font-size: 13px;
            line-height: 22px;
        `,
        heroStats: css`
            display: grid;
            grid-template-columns: repeat(2, minmax(0, 1fr));
            gap: 8px;

            .stat {
                min-width: 0;
                padding: ${compact ? "9px 10px" : "10px 12px"};
                border-radius: ${radius}px;
                background: ${isDarkMode ? "rgba(255,255,255,0.040)" : "rgba(255,255,255,0.64)"};
                box-shadow: inset 0 0 0 1px ${isDarkMode ? "rgba(255,255,255,0.060)" : "rgba(5,5,5,0.055)"};
                backdrop-filter: blur(12px);
            }

            .label {
                color: ${token.colorTextTertiary};
                font-size: 11px;
                line-height: 16px;
            }

            .value {
                display: block;
                margin-top: 2px;
                color: ${token.colorText};
                font-size: 14px;
                font-weight: 600;
                line-height: 22px;
            }
        `,
        table: css`
            white-space: normal !important;

            .ant-table,
            .ant-table-container,
            .ant-table-content {
                display: block;
                background: transparent;
            }

            .ant-table-content {
                overflow: visible !important;
            }

            .ant-table colgroup,
            .ant-table-thead {
                display: none;
            }

            .ant-table-tbody {
                display: grid;
                grid-template-columns: 1fr;
                gap: ${compact ? 8 : 10}px;
                min-height: 180px;
                padding: ${compact ? 8 : 10}px;
                background: transparent;
            }

            .ant-table-tbody > tr {
                display: block;
                background: transparent !important;
            }

            .ant-table-tbody > tr.ant-table-measure-row,
            .ant-table-tbody > tr[aria-hidden="true"],
            .ant-table-tbody > tr.hidden-row {
                display: none !important;
            }

            .ant-table-tbody > tr > td {
                display: block;
                padding: 0 !important;
                border: 0 !important;
                background: transparent !important;
            }

            .ant-table-tbody > tr:hover > td {
                background: transparent !important;
            }

            .ant-table-placeholder {
                grid-column: 1 / -1;
            }

            .ant-table-placeholder > td {
                padding: 32px 12px !important;
                border-radius: ${radius}px;
                background: ${token.colorBgContainer} !important;
            }
        `,
        orderCard: css`
            display: grid;
            grid-template-columns: minmax(190px, 1.1fr) minmax(260px, 1.4fr) minmax(170px, 0.8fr) auto;
            align-items: center;
            gap: ${compact ? 12 : 14}px;
            min-width: 0;
            padding: ${compact ? "11px 12px" : "12px 14px"};
            border-radius: ${radius}px;
            background: ${token.colorBgContainer};
            box-shadow:
                inset 0 0 0 1px ${isDarkMode ? "rgba(255,255,255,0.065)" : "rgba(5,5,5,0.06)"},
                ${isDarkMode ? "0 4px 12px rgba(0,0,0,0.10)" : "0 4px 14px rgba(15,23,42,0.025)"};
            transition: background-color 0.18s ease, box-shadow 0.18s ease;

            &:hover {
                background: ${isDarkMode ? "rgba(255,255,255,0.018)" : token.colorFillAlter};
                box-shadow: inset 0 0 0 1px ${isDarkMode ? "rgba(255,255,255,0.095)" : "rgba(5,5,5,0.085)"};
            }

            @media (max-width: 980px) {
                grid-template-columns: 1fr;
                align-items: stretch;
                gap: ${compact ? 8 : 10}px;
            }
        `,
        orderMain: css`
            min-width: 0;

            .name {
                display: flex;
                align-items: center;
                gap: 8px;
                min-width: 0;
            }

            .dot {
                flex: none;
                width: 6px;
                height: 6px;
                border-radius: 999px;
                background: ${token.colorPrimary};
                box-shadow: 0 0 0 3px ${token.colorPrimaryBg};
            }

            .title {
                min-width: 0;
                color: ${token.colorText};
                font-size: ${compact ? 13 : 14}px;
                font-weight: 600;
                line-height: ${compact ? 20 : 21}px;
            }

            .id {
                display: block;
                margin-top: 2px;
                color: ${token.colorTextTertiary};
                font-size: ${compact ? 11 : 10}px;
                line-height: 16px;
            }

            .id .ant-typography {
                color: inherit;
                font-size: inherit;
                line-height: inherit;
            }
        `,
        infoGrid: css`
            min-width: 0;
            display: grid;
            grid-template-columns: repeat(3, minmax(0, 1fr));
            gap: ${compact ? 10 : 12}px;

            @media (max-width: 620px) {
                grid-template-columns: 1fr;
                gap: 6px;
            }
        `,
        infoItem: css`
            min-width: 0;
            padding-left: ${compact ? 9 : 10}px;
            border-left: 1px solid ${isDarkMode ? "rgba(255,255,255,0.065)" : "rgba(5,5,5,0.055)"};

            .label {
                color: ${token.colorTextTertiary};
                font-size: 11px;
                line-height: 14px;
            }

            .value {
                display: block;
                margin-top: 1px;
                color: ${token.colorTextSecondary};
                font-size: 12px;
                line-height: 17px;
            }

            @media (max-width: 620px) {
                display: grid;
                grid-template-columns: 62px minmax(0, 1fr);
                align-items: center;
                gap: 8px;
                padding-left: 0;
                border-left: 0;

                .value {
                    margin-top: 0;
                    text-align: right;
                }
            }
        `,
        statusStack: css`
            display: flex;
            flex-direction: column;
            align-items: flex-start;
            gap: 5px;

            @media (max-width: 980px) {
                flex-direction: row;
                flex-wrap: wrap;
            }
        `,
        operate: css`
            display: flex;
            justify-content: flex-end;
            white-space: nowrap;
        `,
        modalLog: css`
            .ant-timeline {
                margin-top: 4px;
            }

            .log-card {
                padding: 10px 12px;
                border-radius: ${radius}px;
                background: ${isDarkMode ? "rgba(255,255,255,0.026)" : token.colorFillQuaternary};
            }

            .log-head {
                display: flex;
                align-items: center;
                justify-content: space-between;
                gap: 10px;
                margin-bottom: 6px;
            }

            .log-title {
                color: ${token.colorText};
                font-size: 13px;
                font-weight: 600;
                line-height: 20px;
            }

            .log-time {
                color: ${token.colorTextTertiary};
                font-size: 11px;
                line-height: 16px;
            }

            .log-content {
                color: ${token.colorTextSecondary};
                font-size: 12px;
                line-height: 20px;
            }
        `,
    };
});

const PayType = ({data, value}: { data?: any; value: any }) => {
    const text = enumText(data, value);
    const key = String(value || "").toLowerCase();
    if (key === "alipay") {
        return <Space size={4}><Icon type="AlipayOutlined" style={{color: "#1677ff", fontSize: 13}}/><span style={{color: "#1677ff"}}>{text}</span></Space>;
    }
    if (key === "wxpay" || key === "wechat") {
        return <Space size={4}><Icon type="WechatOutlined" style={{color: "#07c160", fontSize: 13}}/><span style={{color: "#07c160"}}>{text}</span></Space>;
    }
    if (key === "balance") {
        return <Space size={4}><Icon type="RedEnvelopeOutlined" style={{color: "#faad14", fontSize: 13}}/><span style={{color: "#faad14"}}>{text}</span></Space>;
    }
    return <>{text}</>;
};

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

    const renderInfo = (label: string, value: any, copy?: any) => (
        <div className={styles.infoItem}>
            <div className="label">{label}</div>
            <Typography.Text className="value" copyable={copy ? {text: String(copy)} : false} ellipsis={{tooltip: value || "-"}}>
                {value || "-"}
            </Typography.Text>
        </div>
    );

    const renderOrder = (_: any, record: any) => (
        <div className={styles.orderCard}>
            <div className={styles.orderMain}>
                <div className="name">
                    <span className="dot"/>
                    <Typography.Text className="title" ellipsis={{tooltip: record?.name || "-"}}>
                        {record?.name || "-"}
                    </Typography.Text>
                </div>
                <div className="id">
                    <Typography.Text copyable={record?.id ? {text: String(record.id)} : false}>
                        订单ID {record?.id || "-"}
                    </Typography.Text>
                </div>
            </div>
            <div className={styles.infoGrid}>
                {renderInfo("支付方式", <PayType data={payTypes} value={record?.pay_type}/>)}
                {renderInfo("订单金额", formatMoney(record?.amount))}
                {renderInfo("实付金额", formatMoney(record?.pay_amount))}
            </div>
            <div className={styles.statusStack}>
                <Tag color={typeColors?.[String(record?.type)] || "default"}>{enumText(orderType, record?.type)}</Tag>
                <Tag color={statusColors?.[String(record?.status)] || "default"}>{enumText(orderStatus, record?.status)}</Tag>
                <Tag color={payStatusColors?.[String(record?.pay_status)] || "default"}>{enumText(payStatus, record?.pay_status)}</Tag>
            </div>
            <div className={styles.operate}>
                <Button type="link" size="small" onClick={() => loadLogs(record)}>
                    日志
                </Button>
            </div>
        </div>
    );

    return (
        <Body loading={enumLoading}>
            <Row className={styles.page}>
                <Card className={styles.hero} variant="borderless" styles={{body: {padding: 0}}}>
                    <Icon className="hero-watermark" type="ShoppingCartOutlined"/>
                    <div className={styles.heroInner}>
                        <div>
                            <h1 className={styles.heroTitle}>我的订单</h1>
                            <div className={styles.heroDesc}>集中查看账户下的支付订单、到账状态和操作记录。</div>
                        </div>
                        <div className={styles.heroStats}>
                            <div className="stat">
                                <div className="label">订单类型</div>
                                <Typography.Text className="value" ellipsis>在线充值 / 在线下单</Typography.Text>
                            </div>
                            <div className="stat">
                                <div className="label">订单状态</div>
                                <Typography.Text className="value" ellipsis>支付、到账、退款</Typography.Text>
                            </div>
                        </div>
                    </div>
                </Card>
                <ProTable
                    ref={tableRef}
                    rowKey="id"
                    title={<Title>订单记录</Title>}
                    extraRefreshBtn={true}
                    search={true}
                    defaultPage={1}
                    defaultPageSize={DEFAULT_PAGE_SIZE}
                    paginationAffix={true}
                    columns={[
                        {title: "订单ID", dataIndex: "id", valueType: "text", hideInTable: true},
                        {title: "订单类型", dataIndex: "type", valueEnum: buildValueEnum(orderType, typeColors), hideInTable: true},
                        {title: "订单状态", dataIndex: "status", valueEnum: buildValueEnum(orderStatus, statusColors), hideInTable: true},
                        {title: "支付状态", dataIndex: "pay_status", valueEnum: buildValueEnum(payStatus, payStatusColors), hideInTable: true},
                        {title: "支付方式", dataIndex: "pay_type", valueEnum: buildValueEnum(payTypes), hideInTable: true},
                        {title: "创建时间", dataIndex: "create_time", valueType: "dateRange", transform: dateRangeTransform("create_time"), hideInTable: true},
                        {title: "支付时间", dataIndex: "pay_time", valueType: "dateRange", transform: dateRangeTransform("pay_time"), hideInTable: true},
                        {title: "订单记录", dataIndex: "id", hideInSearch: true, render: renderOrder},
                    ]}
                    request={(params: any, sort: any) => getData(params, sort, getOrderList)}
                    {...({
                        tableProps: {
                            className: styles.table,
                            showHeader: false,
                            rowClassName: (record: any) => record?.id ? "" : "hidden-row",
                        },
                    } as any)}
                />
            </Row>
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
                <Body loading={logLoading} className={styles.modalLog}>
                    <Typography.Text type="secondary">
                        订单ID：{currentOrder?.id || "-"}
                    </Typography.Text>
                    {logs.length > 0 ? (
                        <Timeline
                            style={{marginTop: 16}}
                            items={logs.map((item: any) => ({
                                children: (
                                    <div className="log-card">
                                        <div className="log-head">
                                            <span className="log-title">{enumText(logEnums?.type, item?.type)}</span>
                                            <Tooltip title={item?.create_time || "-"}>
                                                <span className="log-time">{item?.create_time || "-"}</span>
                                            </Tooltip>
                                        </div>
                                        <div className="log-content">{item?.content || "-"}</div>
                                    </div>
                                ),
                            }))}
                        />
                    ) : (
                        <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="暂无日志"/>
                    )}
                </Body>
            </ProModal>
        </Body>
    );
};
