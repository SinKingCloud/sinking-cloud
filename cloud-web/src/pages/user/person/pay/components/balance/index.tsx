import React, {useMemo} from "react";
import {App, Tag, Typography} from "antd";
import {createStyles} from "antd-style";
import type {PageTableProps} from "sinking-antd";
import {PageTable, useTheme} from "sinking-antd";
import useEnum from "@/utils/enum";
import {getParams} from "@/utils/page";
import {formatMoney} from "@/utils/string";
import {getBalanceLogList} from "@/service/user/pay";

type BalanceLog = {
    id: number | string;
    type: number;
    money: number;
    money_before: number;
    money_after: number;
    frozen_before: number;
    frozen_after: number;
    describe?: string;
    create_time?: string | null;
};

const typeColors: Record<string, string> = {"0": "green", "1": "red", "2": "red", "3": "green", "4": "gold", "5": "blue", "6": "green"};

const useStyles = createStyles(({css, token}, {isCompactMode}: {isCompactMode: boolean}) => ({
    text: css`
        color: ${token.colorTextSecondary};
        font-size: ${isCompactMode ? 12 : 13}px;
        line-height: 20px;
        font-variant-numeric: tabular-nums;
        white-space: nowrap;
        .ant-typography-copy { margin-inline-start: 6px; font-size: 12px; }

        @media (hover: hover) and (pointer: fine) {
            .ant-typography-copy { opacity: 0; pointer-events: none; }
            &:hover .ant-typography-copy, &:focus-within .ant-typography-copy { opacity: 1; pointer-events: auto; }
        }
    `,
    transaction: css`
        display: grid;
        gap: 2px;
        min-width: 0;
        color: ${token.colorText};
        font-size: ${isCompactMode ? 12 : 13}px;
        line-height: 20px;
        .ant-tag { justify-self: start; margin-inline-end: 0; font-size: 12px; line-height: 18px; }
    `,
    describe: css`
        display: block;
        color: ${token.colorTextSecondary};
        font-size: 12px;
        line-height: 20px;
    `,
    amount: css`
        display: inline-flex;
        align-items: baseline;
        justify-content: flex-start;
        gap: 4px;
        color: ${token.colorTextSecondary};
        font-size: ${isCompactMode ? 16 : 17}px;
        font-weight: 500;
        line-height: 22px;
        font-variant-numeric: lining-nums tabular-nums;
        white-space: nowrap;
        .amount-currency { color: ${token.colorTextTertiary}; font-size: 12px; font-weight: 400; line-height: 20px; }
    `,
    balance: css`
        display: grid;
        justify-items: start;
        gap: 2px;
        min-width: 0;
        color: ${token.colorTextSecondary};
        font-size: ${isCompactMode ? 12 : 13}px;
        line-height: 20px;
        font-variant-numeric: tabular-nums;
        text-align: left;
        white-space: nowrap;

        .balance-line { display: flex; align-items: baseline; justify-content: flex-start; gap: 8px; max-width: 100%; }
        .balance-label { flex: none; color: ${token.colorTextTertiary}; font-size: 12px; }
        .balance-change { display: flex; align-items: baseline; justify-content: flex-start; gap: 8px; min-width: 0; }
        .balance-before, .balance-after { min-width: 0; overflow: hidden; text-overflow: ellipsis; }
        .balance-before { color: ${token.colorTextTertiary}; }
        .balance-arrow { flex: none; color: ${token.colorTextTertiary}; font-size: 12px; }
    `,
}));

export default ({hero}: {hero: React.ReactElement}): React.ReactNode => {
    const {message} = App.useApp();
    const theme = useTheme();
    const {styles} = useStyles(useMemo(() => ({isCompactMode: theme?.isCompactTheme?.() || false}), [theme?.isCompactTheme?.()]));
    const [enumData, enumLoading] = useEnum(["pay.balance_log"]);
    const typeEnum = enumData?.["pay.balance_log"]?.type || {};

    const columns: PageTableProps<BalanceLog>["columns"] = [
        {
            title: "明细ID", dataIndex: "id", width: 180, sorter: true,
            render: (value: number | string) => <Typography.Text className={styles.text} copyable={{text: String(value)}}>{value}</Typography.Text>,
        },
        {
            title: "交易信息", dataIndex: "type", width: 220, ellipsis: {showTitle: false},
            render: (value: number, record: BalanceLog) => <div className={styles.transaction}>
                <Tag color={typeColors[String(value)] || "default"}>{typeEnum[String(value)] || String(value)}</Tag>
                {record?.describe && <Typography.Text className={styles.describe} ellipsis={{tooltip: record?.describe}}>{record?.describe}</Typography.Text>}
            </div>,
        },
        {
            title: "变动金额", dataIndex: "money", width: 140, align: "left", sorter: true,
            render: (value: number) => <span className={styles.amount}>
                <span>{Number(value) > 0 ? "+" : Number(value) < 0 ? "−" : ""}{formatMoney(Math.abs(Number(value)))}</span>
                <span className="amount-currency">元</span>
            </span>,
        },
        {
            title: "余额变化", dataIndex: "money_after", width: 260, align: "left",
            render: (_: number, record: BalanceLog) => <div className={styles.balance}>
                {[
                    {label: "可用", before: record?.money_before, after: record?.money_after},
                    {label: "冻结", before: record?.frozen_before, after: record?.frozen_after},
                ].map((item) => {
                    const before = `¥ ${formatMoney(item.before)}`;
                    const after = `¥ ${formatMoney(item.after)}`;
                    return <div className="balance-line" key={item.label}>
                        <span className="balance-label">{item.label}</span>
                        <span className="balance-change">
                            <span className="balance-before" title={`变动前 ${before}`}>{before}</span>
                            <span className="balance-arrow" aria-hidden="true">→</span>
                            <span className="balance-after" title={`变动后 ${after}`}>{after}</span>
                        </span>
                    </div>;
                })}
            </div>,
        },
        {
            title: "记录时间", dataIndex: "create_time", width: 190, sorter: true,
            render: (value: string | null) => <span className={styles.text}>{value || "-"}</span>,
        },
    ];

    return <PageTable<BalanceLog>
        ariaLabel="余额明细"
        hero={hero}
        rowKey="id"
        columns={columns}
        scroll={{x: 990}}
        defaultPageSize={10}
        loading={enumLoading || undefined}
        requestEnabled={!enumLoading}
        search={{name: "describe", ariaLabel: "搜索交易描述", placeholder: "搜索交易描述", maxLength: 100}}
        filters={[
            {type: "select", name: "type", label: "资金类型", allLabel: "全部类型", valueEnum: typeEnum, disabled: enumLoading},
            {type: "date", name: "create_time", label: "记录时间"},
        ]}
        toolbar={{refresh: {ariaLabel: "刷新余额明细"}}}
        request={async (params, sort) => {
            const response = await getBalanceLogList({body: getParams(params, sort)});
            if (response?.code !== 200) throw new Error(response?.message || "获取资金记录失败");
            return {data: Array.isArray(response?.data?.list) ? response?.data?.list?.filter(Boolean) : [], total: response?.data?.total ?? 0, success: true};
        }}
        onRequestError={(error) => message.error(error?.message || "获取资金记录失败")}
        pagination={{unit: "条"}}
        paginationAffix={{offsetBottom: 15}}
    />;
};
