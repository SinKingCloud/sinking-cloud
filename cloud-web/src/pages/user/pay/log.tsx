import React, {useEffect, useRef, useState} from "react";
import {App, Button, Card, Col, DatePicker, Popover, Row, Segmented, Tag, Typography} from "antd";
import {createStyles} from "antd-style";
import {Body, Icon, ProTable, ProTableRef, Title, useTheme} from "sinking-antd";
import useEnum from "@/utils/enum";
import {getData} from "@/utils/page";
import {formatMoney} from "@/utils/string";
import {getBalanceInfo, getBalanceLogList} from "@/service/user/pay";

type LogTimeRange = "all" | "three_days" | "week" | "month" | "custom";
type SelectOption<T = string> = {
    label: React.ReactNode;
    value: T;
};

const ALL_VALUE = "all";
const typeColors: Record<string, string> = {
    "0": "blue",
    "1": "green",
    "2": "red",
    "3": "gold",
    "4": "purple",
};
const timeOptions: SelectOption<LogTimeRange>[] = [
    {label: "全部", value: "all"},
    {label: "最近3天", value: "three_days"},
    {label: "最近7天", value: "week"},
    {label: "最近一个月", value: "month"},
    {label: "自定义", value: "custom"},
];

const useStyles = createStyles(({token, css, isDarkMode}: any, props: any = {}) => {
    const compact = props?.isCompactMode || false;
    const radius = token.borderRadiusLG;
    return {
        layout: css`
            width: 100%;
            max-width: 1120px;
            margin: 0 auto !important;
        `,
        overviewCol: css`
            min-width: 0;

            @media (min-width: 992px) {
                flex: 0 0 360px;
                max-width: 360px;
                padding-right: ${compact ? 12 : 14}px;
            }

            @media (min-width: 1200px) {
                flex-basis: 400px;
                max-width: 400px;
            }

            @media (max-width: 991px) {
                flex: 0 0 100%;
                max-width: 100%;
                padding-right: 0;
            }
        `,
        recordCol: css`
            min-width: 0;

            @media (min-width: 992px) {
                flex: 1 1 0;
                max-width: none;
            }

            @media (max-width: 991px) {
                flex: 0 0 100%;
                max-width: 100%;
            }
        `,
        overview: css`
            width: 100%;
            overflow: hidden;
            border-radius: ${radius}px;
            background: ${token.colorBgContainer};
            box-shadow: ${isDarkMode ? "none" : "0 10px 28px rgba(15, 23, 42, 0.04)"};

        `,
        balancePanel: css`
            position: relative;
            width: 100%;
            min-height: ${compact ? 270 : 292}px;
            padding: ${compact ? 18 : 22}px;
            display: grid;
            align-content: start;
            gap: ${compact ? 16 : 18}px;
            overflow: hidden;
            isolation: isolate;
            background: ${isDarkMode ? "rgba(255,255,255,0.012)" : token.colorFillAlter};
            box-shadow: ${isDarkMode ? "none" : "inset 0 -1px 0 rgba(15,23,42,0.06)"};

            &::before {
                content: "";
                position: absolute;
                inset: -78px;
                z-index: 0;
                background: radial-gradient(circle at 18% 12%, ${isDarkMode ? "rgba(255,255,255,0.12)" : token.colorPrimaryBgHover} 0%, transparent 38%),
                radial-gradient(circle at 82% 78%, ${isDarkMode ? "rgba(255,255,255,0.10)" : token.colorPrimaryBg} 0%, transparent 42%),
                radial-gradient(circle at 64% 24%, ${isDarkMode ? "rgba(255,255,255,0.08)" : token.colorPrimaryBorder} 0%, transparent 34%),
                linear-gradient(90deg, ${isDarkMode ? "rgba(255,255,255,0.03)" : token.colorPrimaryBg} 1px, transparent 1px),
                linear-gradient(0deg, ${isDarkMode ? "rgba(255,255,255,0.026)" : token.colorPrimaryBg} 1px, transparent 1px);
                background-size: 100% 100%, 100% 100%, 100% 100%, 34px 34px, 34px 34px;
                opacity: ${isDarkMode ? 0.42 : 0.58};
                filter: blur(3px);
                transform: rotate(-7deg) translate3d(-18px, -12px, 0);
                mask-image: linear-gradient(135deg, rgba(0, 0, 0, 0.82), rgba(0, 0, 0, 0.34) 58%, transparent 100%);
            }

            &::after {
                content: "";
                position: absolute;
                right: -176px;
                bottom: -154px;
                z-index: 0;
                width: 500px;
                height: 500px;
                border-radius: 50%;
                background: radial-gradient(circle, ${isDarkMode ? "rgba(255,255,255,0.10)" : token.colorPrimaryBgHover} 0%, transparent 46%),
                conic-gradient(from 0deg, transparent 0 12%, ${isDarkMode ? "rgba(255,255,255,0.12)" : token.colorPrimaryBgHover} 20%, transparent 34%, ${isDarkMode ? "rgba(255,255,255,0.10)" : token.colorPrimaryBg} 52%, transparent 72%, ${isDarkMode ? "rgba(255,255,255,0.09)" : token.colorPrimaryBorder} 86%, transparent 100%);
                filter: blur(42px);
                opacity: ${isDarkMode ? 0.20 : 0.30};
                transform: rotate(18deg) scale(1.02);
            }

            @media (max-width: 991px) {
                min-height: auto;
                box-shadow: ${isDarkMode ? "none" : "inset 0 -1px 0 rgba(15,23,42,0.06)"};
            }

            @media (max-width: 620px) {
                padding: ${compact ? 14 : 16}px;
                gap: 14px;

                &::before {
                    inset: -54px;
                    background-size: 100% 100%, 100% 100%, 100% 100%, 28px 28px, 28px 28px;
                }

                &::after {
                    right: -116px;
                    bottom: -118px;
                    width: 330px;
                    height: 330px;
                    filter: blur(36px);
                }

            }
        `,
        balanceMain: css`
            position: relative;
            z-index: 2;
            display: grid;
            gap: ${compact ? 14 : 16}px;
            min-width: 0;

            .head {
                min-width: 0;
            }

            .title {
                margin: 0;
                color: ${token.colorText};
                font-size: ${compact ? 24 : 26}px;
                font-weight: 700;
                line-height: ${compact ? 32 : 34}px;
            }

            .subtitle {
                margin-top: 5px;
                color: ${token.colorTextSecondary};
                font-size: 12px;
                line-height: 20px;
            }

            .balance-card {
                min-width: 0;
                padding: ${compact ? "13px 14px" : "15px 16px"};
                border-radius: ${token.borderRadius}px;
                background: rgba(255, 255, 255, ${isDarkMode ? 0.025 : 0.58});
            }

            .balance-content {
                min-width: 0;
                flex: 1;
            }

            .label {
                color: ${token.colorTextTertiary};
                font-size: 13px;
                line-height: 20px;
            }

            .money {
                display: block;
                margin-top: 2px;
                color: ${token.colorText};
                font-size: ${compact ? 30 : 34}px;
                font-weight: 700;
                line-height: ${compact ? 38 : 42}px;
                font-variant-numeric: tabular-nums;
            }

            .desc {
                margin-top: 4px;
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 18px;
            }

            @media (max-width: 620px) {
                gap: 12px;

                .title {
                    font-size: ${compact ? 22 : 23}px;
                    line-height: ${compact ? 30 : 31}px;
                }

                .money {
                    font-size: ${compact ? 26 : 28}px;
                    line-height: ${compact ? 34 : 36}px;
                }
            }
        `,
        balanceGrid: css`
            position: relative;
            z-index: 2;
            display: grid;
            grid-template-columns: 1fr;
            gap: 8px;

            @media (max-width: 920px) {
                grid-template-columns: repeat(2, minmax(0, 1fr));
            }

            @media (max-width: 620px) {
                gap: 6px;
            }

            @media (max-width: 430px) {
                grid-template-columns: 1fr;
            }
        `,
        balanceItem: css`
            min-width: 0;
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 12px;
            padding: ${compact ? "9px 11px" : "8px 10px"};
            border-radius: ${token.borderRadius}px;
            background: rgba(255, 255, 255, ${isDarkMode ? 0.02 : 0.56});

            .label {
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 18px;
            }

            .value {
                display: block;
                color: ${token.colorText};
                font-size: 13px;
                line-height: 20px;
                text-align: right;
                white-space: nowrap;
                overflow: hidden;
                text-overflow: ellipsis;
            }

            .money {
                color: ${token.colorSuccess};
            }

            .frozen {
                color: ${token.colorWarning};
            }

            .recharge {
                color: ${token.colorPrimary};
            }

            .consume {
                color: ${token.colorError};
            }

            @media (max-width: 620px) {
                padding: 7px 8px;

                .label {
                    margin-bottom: 2px;
                    font-size: ${compact ? 12 : 11}px;
                }

                .value {
                    font-size: ${compact ? 13 : 12}px;
                    line-height: ${compact ? 19 : 18}px;
                }
            }
        `,
        notePanel: css`
            position: relative;
            z-index: 2;
            display: grid;
            gap: ${compact ? 7 : 8}px;
            padding-top: ${compact ? 2 : 4}px;

            .note-title {
                color: ${token.colorTextSecondary};
                font-size: 12px;
                font-weight: 600;
                line-height: 18px;
            }

            .note-list {
                display: grid;
                gap: 6px;
            }

            .note-item {
                display: grid;
                grid-template-columns: 3px minmax(0, 1fr);
                gap: 8px;
                align-items: start;
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 18px;
            }

            .note-item::before {
                content: "";
                width: 3px;
                height: 3px;
                margin-top: 8px;
                border-radius: 999px;
                background: ${token.colorTextTertiary};
                opacity: 0.68;
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
                padding: ${compact ? "8px 0" : "10px 0"};
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

            .ant-table-pagination {
                justify-content: center !important;
            }
        `,
        recordCard: css`
            min-width: 0;
            border-radius: ${radius}px;

            > .ant-card-body {
                padding: ${compact ? "10px" : "12px"};
            }

            > .ant-card-head {
                min-height: ${compact ? 38 : 42}px;
                padding: 0 ${compact ? 10 : 12}px;
            }

            > .ant-card-head .ant-card-head-title {
                padding: ${compact ? "8px 0" : "10px 0"};
            }

            .ant-row {
                margin: 0 !important;
            }

            .ant-row > .ant-col:nth-child(2) {
                padding: 0 !important;
            }

            .ant-row > .ant-col:nth-child(2) .ant-card {
                background: transparent;
                box-shadow: none;
            }

            .ant-row > .ant-col:nth-child(2) .ant-card-head {
                display: none;
            }

            .ant-row > .ant-col:nth-child(2) .ant-card-body {
                padding: 0 !important;
            }
        `,
        recordFilter: css`
            min-width: 0;
            margin-bottom: ${compact ? 8 : 10}px;

            .filter-row {
                position: relative;
                min-width: 0;
            }

            .filter-row + .filter-row {
                margin-top: ${compact ? 5 : 9}px;
            }

            .filter-options {
                min-width: 0;
                padding: 1px 0;
                overflow-x: auto;
                scrollbar-width: none;
            }

            .filter-options::-webkit-scrollbar {
                display: none;
            }

            .filter-scroll-button {
                position: absolute;
                top: 50%;
                z-index: 2;
                width: ${compact ? 22 : 26}px;
                height: ${compact ? 22 : 26}px;
                min-width: ${compact ? 22 : 26}px;
                padding: 0;
                color: ${token.colorTextQuaternary};
                border: none;
                background: transparent;
                box-shadow: none;
                opacity: 0;
                pointer-events: none;
                transform: translateY(-50%);
                transition: none;
            }

            .filter-row:hover .filter-scroll-button,
            .filter-scroll-button:hover {
                opacity: 1;
                pointer-events: auto;
                color: ${token.colorTextSecondary};
                background: transparent !important;
            }

            .filter-scroll-button.left {
                left: ${compact ? -14 : -16}px;
            }

            .filter-scroll-button.right {
                right: ${compact ? -14 : -16}px;
            }

            .ant-segmented {
                max-width: 100%;
                padding: ${compact ? 0 : 2}px;
                background: ${token.colorFillQuaternary};
                border-radius: ${token.borderRadius}px;
            }

            .ant-segmented-group {
                flex-wrap: nowrap;
                gap: ${compact ? 0 : 2}px;
            }

            .ant-segmented-item {
                border-radius: ${token.borderRadiusSM}px;
                color: ${token.colorTextSecondary};
                font-size: ${compact ? 12 : 13}px;
                min-height: ${compact ? 25 : 31}px;
                line-height: ${compact ? "25px" : "31px"};
            }

            .ant-segmented-item-label {
                min-height: ${compact ? 25 : 31}px;
                line-height: ${compact ? "25px" : "31px"};
                padding: 0 ${compact ? 10 : 13}px;
            }

            .ant-segmented-item-selected {
                color: ${token.colorPrimary};
                background: ${isDarkMode ? "rgba(255,255,255,0.08)" : token.colorBgContainer};
            }

            .custom-filter-label {
                display: inline-flex;
            }

            .custom-picker-panel {
                position: relative;
                display: block;
                width: ${compact ? 276 : 300}px;
            }

            .custom-picker-group {
                display: flex;
                gap: ${compact ? 9 : 10}px;
            }

            .custom-picker-panel .ant-picker {
                flex: 1;
                min-width: 0;
                width: 100%;
            }

            @media (max-width: 575px) {
                .filter-scroll-button {
                    display: none;
                }

                .ant-segmented-item {
                    min-height: ${compact ? 25 : 30}px;
                    line-height: ${compact ? "25px" : "30px"};
                }

                .ant-segmented-item-label {
                    min-height: ${compact ? 25 : 30}px;
                    line-height: ${compact ? "25px" : "30px"};
                    padding: 0 ${compact ? 8 : 11}px;
                }
            }
        `,
        logCard: css`
            display: grid;
            grid-template-columns: minmax(180px, 1fr) minmax(280px, 1.5fr) minmax(120px, auto);
            align-items: center;
            gap: ${compact ? 12 : 14}px;
            min-width: 0;
            padding: ${compact ? "11px 12px" : "12px 14px"};
            border-radius: ${radius}px;
            background: ${token.colorBgContainer};
            box-shadow: inset 0 0 0 1px ${isDarkMode ? "rgba(255,255,255,0.065)" : "rgba(5,5,5,0.06)"},
            ${isDarkMode ? "0 4px 12px rgba(0,0,0,0.10)" : "0 4px 14px rgba(15,23,42,0.025)"};
            transition: background-color 0.18s ease, box-shadow 0.18s ease;

            &:hover {
                background: ${isDarkMode ? "rgba(255,255,255,0.018)" : token.colorFillAlter};
                box-shadow: inset 0 0 0 1px ${isDarkMode ? "rgba(255,255,255,0.095)" : "rgba(5,5,5,0.085)"};
            }

            @media (max-width: 1320px) and (min-width: 901px) {
                grid-template-columns: 1fr;
                align-items: stretch;
                gap: ${compact ? 8 : 9}px;
            }

            @media (max-width: 900px) {
                grid-template-columns: 1fr;
                align-items: stretch;
                gap: ${compact ? 9 : 11}px;
                min-width: 0;
            }
        `,
        logMain: css`
            min-width: 0;

            @media (max-width: 1320px) and (min-width: 901px) {
                grid-column: 1 / -1;
                grid-row: 1;
            }

            .type {
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

            .desc {
                display: block;
                margin-top: 3px;
                color: ${token.colorTextTertiary};
                font-size: 11px;
                line-height: 16px;
            }
        `,
        changeGrid: css`
            min-width: 0;
            display: grid;
            grid-template-columns: repeat(2, minmax(0, 1fr));
            gap: ${compact ? 10 : 12}px;

            @media (max-width: 1320px) and (min-width: 901px) {
                grid-column: 1 / -1;
                grid-row: 2;
                padding-top: 0;
                grid-template-columns: repeat(2, minmax(0, 1fr));
                gap: ${compact ? 8 : 10}px;
            }

            @media (max-width: 620px) {
                grid-template-columns: 1fr;
                gap: ${compact ? 6 : 7}px;
            }
        `,
        changeItem: css`
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

            @media (max-width: 900px) {
                padding-left: 0;
                border-left: 0;
            }

            @media (max-width: 620px) {
                display: grid;
                grid-template-columns: 62px minmax(0, 1fr);
                align-items: center;
                gap: 8px;

                .value {
                    margin-top: 0;
                    text-align: right;
                }
            }
        `,
        amountBox: css`
            min-width: 0;
            text-align: right;
            font-variant-numeric: tabular-nums;

            @media (max-width: 1320px) and (min-width: 901px) {
                display: grid;
                grid-template-columns: minmax(0, 1fr) auto;
                align-items: end;
                gap: 8px;
                grid-column: 1 / -1;
                grid-row: 3;
                padding-top: ${compact ? 2 : 3}px;
                text-align: left;

                .amount {
                    grid-column: 2;
                    text-align: right;
                }

                .time {
                    grid-column: 1;
                    grid-row: 1;
                    margin-top: 0;
                    text-align: left;
                }
            }

            .amount {
                min-width: 0;
                color: ${token.colorText};
                font-size: ${compact ? 15 : 16}px;
                font-weight: 600;
                line-height: 24px;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
            }

            .amount.income {
                color: ${token.colorSuccess};
            }

            .amount.outcome {
                color: ${token.colorError};
            }

            .time {
                display: block;
                margin-top: 1px;
                color: ${token.colorTextTertiary};
                font-size: 11px;
                line-height: 16px;
            }

            @media (max-width: 900px) {
                display: grid;
                grid-template-columns: minmax(0, 1fr) auto;
                align-items: end;
                gap: 8px;
                padding-top: ${compact ? 8 : 9}px;
                text-align: left;

                .amount {
                    grid-column: 2;
                    line-height: 20px;
                    text-align: right;
                }

                .time {
                    grid-column: 1;
                    grid-row: 1;
                    margin-top: 0;
                    text-align: left;
                }
            }
        `,
    };
});

export default (): React.ReactNode => {
    const {message} = App.useApp();
    const theme = useTheme();
    const isCompactMode = theme?.isCompactTheme?.() || false;
    const {styles} = useStyles({isCompactMode});
    const tableRef = useRef<ProTableRef>({} as ProTableRef);
    const typeScrollRef = useRef<HTMLDivElement | null>(null);
    const timeScrollRef = useRef<HTMLDivElement | null>(null);
    const filtersRef = useRef<any>({});
    const customRangeRef = useRef<any[]>([]);
    const [enumData, enumLoading] = useEnum(["pay.balance_log"]) as any;
    const [info, setInfo] = useState<any>();
    const [infoLoading, setInfoLoading] = useState(true);
    const [filterType, setFilterType] = useState(ALL_VALUE);
    const [timeRange, setTimeRange] = useState<LogTimeRange>("all");
    const [customOpen, setCustomOpen] = useState(false);
    const [customRange, setCustomRange] = useState<any>();
    const typeEnum = enumData?.["pay.balance_log"]?.type || {};
    const typeOptions: SelectOption[] = [
        {label: "全部", value: ALL_VALUE},
        ...Object.entries(typeEnum || {}).map(([value, label]) => ({value, label: String(label)})),
    ];

    const loadInfo = async () => {
        setInfoLoading(true);
        await getBalanceInfo({
            onSuccess: (res: any) => setInfo(res?.data || {}),
            onFail: (res: any) => message?.error(res?.message || "获取余额信息失败"),
        }).finally(() => {
            setInfoLoading(false);
        });
    };

    useEffect(() => {
        void loadInfo();
    }, []);

    const getFilters = (type: string, range: LogTimeRange, custom: any = customRange) => {
        const next: any = type === ALL_VALUE ? {} : {type};
        if (range === "all") {
            return next;
        }
        if (range === "custom") {
            if (!custom?.[0] && !custom?.[1]) {
                return next;
            }
            if (custom?.[0]) {
                next.create_time_start = custom[0]?.startOf
                    ? custom[0].startOf("day").format("YYYY-MM-DD HH:mm:ss")
                    : custom[0];
            }
            if (custom?.[1]) {
                next.create_time_end = custom[1]?.endOf
                    ? custom[1].endOf("day").format("YYYY-MM-DD HH:mm:ss")
                    : custom[1];
            }
            return next;
        }
        const days = range === "three_days" ? 3 : range === "week" ? 7 : 30;
        const end = new Date();
        const start = new Date(end.getTime());
        start.setDate(start.getDate() - days + 1);
        const startMonth = String(start.getMonth() + 1).padStart(2, "0");
        const startDay = String(start.getDate()).padStart(2, "0");
        const endMonth = String(end.getMonth() + 1).padStart(2, "0");
        const endDay = String(end.getDate()).padStart(2, "0");
        return {
            ...next,
            create_time_start: `${start.getFullYear()}-${startMonth}-${startDay} 00:00:00`,
            create_time_end: `${end.getFullYear()}-${endMonth}-${endDay} 23:59:59`,
        };
    };

    const updateFilters = (next: any) => {
        filtersRef.current = next;
        tableRef.current?.refreshTableData?.();
    };

    const changeType = (value: string) => {
        setFilterType(value);
        updateFilters(getFilters(value, timeRange));
    };

    const changeTimeRange = (value: LogTimeRange) => {
        if (value === "custom") {
            setCustomOpen(true);
            return;
        }
        setTimeRange(value);
        customRangeRef.current = [];
        setCustomRange(undefined);
        setCustomOpen(false);
        updateFilters(getFilters(filterType, value, undefined));
    };

    const applyCustomTime = (value: any, index: 0 | 1) => {
        const next = [...(customRangeRef.current || [])];
        if (!value) {
            next[index] = undefined;
            customRangeRef.current = next;
            setCustomRange(next);
            if (next?.[0] || next?.[1]) {
                setTimeRange("custom");
                updateFilters(getFilters(filterType, "custom", next));
            } else {
                setTimeRange("all");
                updateFilters(getFilters(filterType, "all", undefined));
            }
            return;
        }
        next[index] = value;
        const range = next?.[0] && next?.[1] && next[0]?.isAfter?.(next[1])
            ? [next[1], next[0]]
            : next;
        customRangeRef.current = range;
        setCustomRange(range);
        setTimeRange("custom");
        updateFilters(getFilters(filterType, "custom", range));
        if (range?.[0] && range?.[1]) {
            setCustomOpen(false);
        }
    };

    const scrollFilter = (target: "type" | "time", distance: number) => {
        const ref = target === "type" ? typeScrollRef : timeScrollRef;
        ref.current?.scrollBy({left: distance});
    };

    const renderChange = (label: string, before: any, after: any) => (
        <div className={styles.changeItem}>
            <div className="label">{label}</div>
            <Typography.Text className="value" ellipsis={{tooltip: `${formatMoney(before)} → ${formatMoney(after)}`}}>
                {formatMoney(before)} → {formatMoney(after)}
            </Typography.Text>
        </div>
    );

    const renderLog = (_: any, record: any) => {
        const money = Number(record?.money || 0);
        const amountClassName = money > 0 ? "income" : money < 0 ? "outcome" : "";
        const typeText = typeEnum?.[String(record?.type)] || typeEnum?.[record?.type] || (
            record?.type === undefined || record?.type === null || record?.type === "" ? "-" : String(record?.type)
        );
        return (
            <div className={styles.logCard}>
                <div className={styles.logMain}>
                    <div className="type">
                        <span className="dot"/>
                        <Tag color={typeColors?.[String(record?.type)] || "default"}>
                            {typeText}
                        </Tag>
                    </div>
                    <Typography.Text className="desc" ellipsis={{tooltip: record?.describe || "暂无描述"}}>
                        {record?.describe || "暂无描述"}
                    </Typography.Text>
                </div>
                <div className={styles.changeGrid}>
                    {renderChange("可用余额", record?.money_before, record?.money_after)}
                    {renderChange("冻结金额", record?.frozen_before, record?.frozen_after)}
                </div>
                <div className={styles.amountBox}>
                    <div className={`amount ${amountClassName}`}>
                        {money > 0 ? "+" : ""}{formatMoney(money)}
                    </div>
                    <Typography.Text className="time" ellipsis={{tooltip: record?.create_time || "-"}}>
                        {record?.create_time || "-"}
                    </Typography.Text>
                </div>
            </div>
        );
    };

    const balanceItems = [
        {label: "可用余额", value: info?.money, className: "money"},
        {label: "冻结金额", value: info?.frozen, className: "frozen"},
        {label: "累计充值", value: info?.total_recharge, className: "recharge"},
        {label: "累计消费", value: info?.total_consume, className: "consume"},
    ];
    const totalMoney = Number(info?.money || 0) + Number(info?.frozen || 0);
    const loading = enumLoading || infoLoading;
    const customPickerPanel = (
        <div className="custom-picker-panel">
            <div className="custom-picker-group">
                <DatePicker
                    inputReadOnly
                    placeholder="开始日期"
                    value={customRange?.[0]}
                    getPopupContainer={() => document.body}
                    onChange={(value: any) => applyCustomTime(value, 0)}
                />
                <DatePicker
                    inputReadOnly
                    placeholder="结束日期"
                    value={customRange?.[1]}
                    getPopupContainer={() => document.body}
                    onChange={(value: any) => applyCustomTime(value, 1)}
                />
            </div>
        </div>
    );
    const customTimeOptions = timeOptions.map((item) => item.value === "custom" ? {
        ...item,
        label: (
            <Popover
                arrow
                content={customPickerPanel}
                open={customOpen}
                placement="bottomLeft"
                trigger="click"
                onOpenChange={(open) => {
                    setCustomOpen(open);
                }}>
                <span className="custom-filter-label">自定义</span>
            </Popover>
        ),
    } : item);

    return (
        <Body loading={loading}>
            {!loading && (
                <Row className={styles.layout} gutter={[0, 14]}>
                    <Col className={styles.overviewCol} span={24}>
                        <Card className={styles.overview} variant="borderless" styles={{body: {padding: 0}}}>
                            <div className={styles.balancePanel}>
                                <div className={styles.balanceMain}>
                                    <div className="head">
                                        <h1 className="title">财务明细</h1>
                                        <div
                                            className="subtitle">聚合账户余额、冻结资金和资金流水，便于快速核对账户变化。
                                        </div>
                                    </div>
                                    <div className="balance-card">
                                        <div className="balance-content">
                                            <div className="label">账户共计余额</div>
                                            <Typography.Text className="money"
                                                             ellipsis={{tooltip: formatMoney(totalMoney)}}>
                                                {formatMoney(totalMoney)}
                                            </Typography.Text>
                                            <div className="desc">该余额包含可用余额和冻结余额。</div>
                                        </div>
                                    </div>
                                </div>
                                <div className={styles.balanceGrid}>
                                    {balanceItems.map((item) => (
                                        <div className={styles.balanceItem} key={item.label}>
                                            <div className="label">{item.label}</div>
                                            <Typography.Text className={`value ${item.className}`}
                                                             ellipsis={{tooltip: formatMoney(item.value)}}>
                                                {formatMoney(item.value)}
                                            </Typography.Text>
                                        </div>
                                    ))}
                                </div>
                                <div className={styles.notePanel}>
                                    <div className="note-title">记账说明</div>
                                    <div className="note-list">
                                        <div className="note-item">可用余额用于支付、充值和业务结算。</div>
                                        <div className="note-item">冻结金额来自待处理或风控中的资金。</div>
                                        <div className="note-item">资金流水按发生时间倒序展示，便于核对。</div>
                                    </div>
                                </div>
                            </div>
                        </Card>
                    </Col>
                    <Col className={styles.recordCol} span={24}>
                        <Card
                            className={styles.recordCard}
                            title={<Title>资金记录</Title>}
                            extra={(
                                <Button
                                    color="default"
                                    icon={<Icon type="RedoOutlined" style={{transform: "rotate(260deg)"}}/>}
                                    variant="link"
                                    onClick={() => tableRef.current?.refreshTableData?.()}
                                />
                            )}
                            variant="borderless">
                            <div className={styles.recordFilter}>
                                <div className="filter-row">
                                    <Button
                                        className="filter-scroll-button left"
                                        htmlType="button"
                                        icon={<Icon type="CaretLeftFilled"/>}
                                        size="small"
                                        type="text"
                                        onClick={() => scrollFilter("type", -180)}
                                    />
                                    <div ref={typeScrollRef} className="filter-options">
                                        <Segmented
                                            options={typeOptions}
                                            value={filterType}
                                            onChange={(value) => changeType(String(value))}
                                        />
                                    </div>
                                    <Button
                                        className="filter-scroll-button right"
                                        htmlType="button"
                                        icon={<Icon type="CaretRightFilled"/>}
                                        size="small"
                                        type="text"
                                        onClick={() => scrollFilter("type", 180)}
                                    />
                                </div>
                                <div className="filter-row">
                                    <Button
                                        className="filter-scroll-button left"
                                        htmlType="button"
                                        icon={<Icon type="CaretLeftFilled"/>}
                                        size="small"
                                        type="text"
                                        onClick={() => scrollFilter("time", -180)}
                                    />
                                    <div ref={timeScrollRef} className="filter-options">
                                        <Segmented
                                            options={customTimeOptions}
                                            value={timeRange}
                                            onChange={(value) => changeTimeRange(value as LogTimeRange)}
                                        />
                                    </div>
                                    <Button
                                        className="filter-scroll-button right"
                                        htmlType="button"
                                        icon={<Icon type="CaretRightFilled"/>}
                                        size="small"
                                        type="text"
                                        onClick={() => scrollFilter("time", 180)}
                                    />
                                </div>
                            </div>
                            <ProTable
                                ref={tableRef}
                                rowKey="id"
                                title={false}
                                extraRefreshBtn={false}
                                search={false}
                                defaultPage={1}
                                defaultPageSize={10}
                                pageInTable={true}
                                paginationProps={{placement: ["bottomCenter"]} as any}
                                columns={[
                                    {title: "资金记录", dataIndex: "id", hideInSearch: true, render: renderLog},
                                ]}
                                request={(params: any, sort: any) => getData({
                                    order_by_field: "id",
                                    order_by_type: "desc",
                                    ...filtersRef.current,
                                    ...params,
                                }, sort, getBalanceLogList)}
                                {...({
                                    tableProps: {
                                        className: styles.table,
                                        showHeader: false,
                                        rowClassName: (record: any) => record?.id ? "" : "hidden-row",
                                    },
                                } as any)}
                            />
                        </Card>
                    </Col>
                </Row>
            )}
        </Body>
    );
};
