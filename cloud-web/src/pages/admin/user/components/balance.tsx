import React, {forwardRef, useEffect, useImperativeHandle, useRef, useState} from "react";
import {
    App,
    Button,
    DatePicker,
    Empty,
    Form,
    Input,
    InputNumber,
    Pagination,
    Popover,
    Segmented,
    Spin,
} from "antd";
import type {FormInstance} from "antd";

import {Icon, ProModal, ProModalRef, Title, useTheme} from "sinking-antd";
import {createStyles} from "antd-style";
import {getUserBalanceInfo, getUserBalanceLogList, operateUserBalance} from "@/service/admin/pay";
import useEnum from "@/utils/enum";
import {AdminControlCode} from "@/utils/permission";
import {formatMoney} from "@/utils/string";

// ==================== 类型定义 ====================

/** 财务弹窗对外暴露的方法 */
export type BalanceRef = {
    open: (record: any) => void;
};

type BalanceProps = {
    control: any;
    onSuccess?: () => void;
};

type ActivePanel = "overview" | "operate" | "log";
type OperateAccount = "money" | "frozen";
type OperationType = "add" | "deduct";
type LogTimeRange = "all" | "three_days" | "week" | "month" | "custom";

type ControlState = {
    info: boolean;
    operate: boolean;
    log: boolean;
};

type SelectOption<T = string> = {
    label: string;
    value: T;
};

type Styles = Record<string, string>;

type PanelConfig = {
    label: string;
    value: ActivePanel;
    visible: boolean;
    render: () => React.ReactNode;
};

type FilterScrollerProps = {
    options: SelectOption<any>[];
    value: string;
    styles: Styles;
    onChange: (value: string) => void;
    renderItem?: (item: SelectOption, value: string, node: React.ReactNode) => React.ReactNode;
};

type BalanceOverviewProps = {
    balanceInfo: any;
    loading: boolean;
    canOperate: boolean;
    styles: Styles;
    onOperate: (account: OperateAccount) => void;
};

type BalanceOperateProps = {
    account: OperateAccount;
    balanceInfo: any;
    form: FormInstance;
    styles: Styles;
    onAccountChange: (account: OperateAccount) => void;
    onSubmit: (values: any) => void;
};

type BalanceLogProps = {
    userId?: number | string;
    hasControl: boolean;
    typeEnum: Record<string, string>;
    styles: Styles;
    message: any;
};

// ==================== 常量配置 ====================

const ALL_VALUE = "all";
const PAGE_SIZE = 5;
const DEFAULT_OPERATION_TYPE: OperationType = "add";

const ACCOUNT_OPTIONS: SelectOption<OperateAccount>[] = [
    {label: "可用余额", value: "money"},
    {label: "冻结金额", value: "frozen"},
];

const LOG_TIME_OPTIONS: SelectOption<LogTimeRange>[] = [
    {label: "全部", value: "all"},
    {label: "最近3天", value: "three_days"},
    {label: "最近一周", value: "week"},
    {label: "最近一个月", value: "month"},
    {label: "自定义", value: "custom"},
];

const BALANCE_ITEMS: { label: string; field: string; account: OperateAccount }[] = [
    {label: "可用余额", field: "money", account: "money"},
    {label: "冻结金额", field: "frozen", account: "frozen"},
];

const STAT_ITEMS = [
    {label: "累计充值", field: "total_recharge"},
    {label: "累计消费", field: "total_consume"},
];

// ==================== 样式 ====================

const useStyles: any = createStyles(({css, token}: any, props: any = {}): any => {
    const isCompactMode = props?.isCompactMode || false;

    return {
        // 基础弹窗
        modalBody: {
            display: "flex",
            flexDirection: "column",
            minHeight: 0,
            padding: 0,
        },
        panelHeader: {
            display: "flex",
            justifyContent: "center",
            margin: isCompactMode ? "4px 0 12px" : `5px 0 15px 0`,
            "@media (max-width: 575px)": {
                margin: isCompactMode ? "5px 0 12px" : `${token.marginXXS}px 0 ${token.marginSM}px`,
            },
        },
        segment: {
            width: "100%",
            "& .ant-segmented-group": {
                width: "100%",
            },
            "& .ant-segmented-item": {
                flex: 1,
                textAlign: "center",
            },
        },
        loading: {
            minHeight: 168,
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
        },

        // 概览模块
        overview: css`
            padding-top: ${isCompactMode ? 4 : 2}px;

            .balance-grid {
                display: grid;
                grid-template-columns: repeat(2, minmax(0, 1fr));
                gap: ${isCompactMode ? 10 : token.marginSM}px;
            }

            .balance-card {
                min-width: 0;
                padding: ${isCompactMode ? 10 : token.paddingSM}px;
                border: 1px solid ${token.colorBorderSecondary};
                border-radius: ${token.borderRadius}px;
                background: ${token.colorFillQuaternary};
            }

            .head,
            .overview-footer {
                display: flex;
                align-items: center;
                justify-content: space-between;
                gap: ${token.marginXS}px;
            }

            .head {
                margin-bottom: ${token.marginXXS}px;
            }

            .label {
                color: ${token.colorTextTertiary};
                font-size: ${token.fontSizeSM}px;
                line-height: 18px;
            }

            .action {
                flex: none;
                height: ${token.controlHeightXS}px;
                padding: 0 2px;
                font-size: ${token.fontSizeSM}px;
                line-height: ${token.controlHeightXS}px;
            }

            .balance-card .value {
                color: ${token.colorText};
                font-size: 22px;
                font-weight: 600;
                line-height: 30px;
            }

            .value {
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
                font-variant-numeric: tabular-nums;
            }

            .stats {
                display: grid;
                grid-template-columns: repeat(2, minmax(0, 1fr));
                margin-top: ${isCompactMode ? 10 : token.marginSM}px;
                padding: ${isCompactMode ? "9px 0" : `${token.paddingSM}px 0`};
                border: 1px solid ${token.colorBorderSecondary};
                border-radius: ${token.borderRadius}px;
            }

            .stats .item {
                min-width: 0;
                padding: 0 ${token.paddingSM}px;
            }

            .stats .item + .item {
                border-left: 1px solid ${token.colorBorderSecondary};
            }

            .stats .label {
                margin-bottom: 2px;
            }

            .stats .value {
                color: ${token.colorTextSecondary};
                font-size: ${token.fontSize}px;
                line-height: 22px;
            }

            .overview-footer {
                margin-top: ${isCompactMode ? 10 : token.marginSM}px;
                color: ${token.colorTextTertiary};
                font-size: ${token.fontSizeSM}px;
                line-height: 20px;
            }

            @media (max-width: 575px) {
                .balance-grid {
                    gap: ${isCompactMode ? 8 : 9}px;
                }

                .balance-card {
                    padding: ${isCompactMode ? 8 : 9}px;
                }

                .balance-card .value {
                    font-size: ${isCompactMode ? 16 : 17}px;
                    line-height: 24px;
                }
            }
        `,

        // 操作模块
        operate: css`
            display: flex;
            flex-direction: column;
            gap: ${isCompactMode ? 12 : 14}px;

            .ant-form-item {
                margin-bottom: 0;
            }

            .account-switch {
                display: grid;
                grid-template-columns: repeat(2, minmax(0, 1fr));
                gap: ${isCompactMode ? 8 : 10}px;
            }

            .account-option {
                min-width: 0;
                padding: ${isCompactMode ? "9px 10px" : "10px 12px"};
                text-align: left;
                border: 1px solid ${token.colorBorderSecondary};
                border-radius: ${token.borderRadius}px;
                background: ${token.colorBgContainer};
                cursor: pointer;
                font: inherit;
            }

            .account-option.active {
                border-color: ${token.colorPrimary};
                background: ${token.colorPrimaryBg};
                box-shadow: inset 0 0 0 1px ${token.colorPrimaryBorder};
            }

            .account-head,
            .field-head {
                display: flex;
                align-items: center;
                justify-content: space-between;
                gap: ${token.marginXS}px;
            }

            .account-label,
            .field-label {
                color: ${token.colorTextSecondary};
                font-size: ${token.fontSizeSM}px;
                line-height: 18px;
            }

            .account-current {
                color: ${token.colorTextTertiary};
                font-size: 12px;
                line-height: 18px;
            }

            .account-value {
                min-width: 0;
                margin-top: ${isCompactMode ? 3 : 4}px;
                color: ${token.colorText};
                font-size: ${isCompactMode ? 16 : 17}px;
                font-weight: 600;
                line-height: ${isCompactMode ? 23 : 24}px;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
                font-variant-numeric: tabular-nums;
            }

            .form-grid {
                display: grid;
                grid-template-columns: repeat(2, minmax(0, 1fr));
                gap: ${isCompactMode ? 11 : 12}px;
            }

            .form-field {
                min-width: 0;
            }

            .form-field.full {
                grid-column: 1 / -1;
            }

            .field-head {
                margin-bottom: ${isCompactMode ? 5 : 6}px;
            }

            .required {
                margin-right: 3px;
                color: ${token.colorError};
            }

            .ant-input-number {
                width: 100%;
            }

            .textarea {
                resize: none;
                border-radius: ${token.borderRadius}px;
            }

            @media (max-width: 575px) {
                .form-grid {
                    grid-template-columns: 1fr;
                }

                .account-switch {
                    gap: ${isCompactMode ? 7 : 8}px;
                }

                .account-option {
                    padding: ${isCompactMode ? "8px 9px" : "9px 10px"};
                }

                .account-value {
                    font-size: ${isCompactMode ? 15 : 16}px;
                    line-height: 22px;
                }
            }
        `,

        // 明细模块
        logFilter: {
            flex: "none",
            display: "flex",
            flexDirection: "column",
            gap: 6,
            marginBottom: 10,
        },
        logPanel: {
            display: "flex",
            flexDirection: "column",
            minHeight: 0,
        },
        customPickerPanel: {
            position: "relative",
            display: "block",
            width: isCompactMode ? 276 : 300,
            "& .ant-picker": {
                width: "100%",
            },
        },
        customPickerGroup: {
            display: "flex",
            gap: isCompactMode ? 7 : 8,
            "& .ant-picker": {
                flex: 1,
                minWidth: 0,
            },
        },
        filterRow: {
            position: "relative",
            overflow: "visible",
            "&:hover .scrollButton": {
                opacity: 1,
            },
        },
        filterScrollButton: {
            position: "absolute",
            top: 0,
            zIndex: 1,
            width: 24,
            height: 24,
            minWidth: 24,
            padding: 0,
            color: token.colorTextQuaternary,
            border: "none",
            background: "transparent",
            boxShadow: "none",
            opacity: 0,
            transition: "none",
            "& *": {
                transition: "none",
            },
            "&:hover, &:focus, &:active": {
                opacity: 1,
                color: token.colorTextSecondary,
                background: "transparent !important",
            },
            "&.left": {
                left: -12,
            },
            "&.right": {
                right: -12,
            },
            "@media (max-width: 575px)": {
                display: "none",
            },
        },
        filterScroll: {
            display: "flex",
            gap: 6,
            overflowX: "auto",
            padding: isCompactMode ? "0 1px 2px" : "0 0 2px",
            scrollbarWidth: "none",
            "&::-webkit-scrollbar": {
                display: "none",
            },
            "& .pill": {
                flex: "none",
                height: isCompactMode ? 25 : 24,
                padding: "0 9px",
                border: "none",
                borderRadius: 999,
                color: token.colorTextSecondary,
                background: token.colorFillQuaternary,
                fontSize: 12,
                lineHeight: isCompactMode ? "25px" : "24px",
                whiteSpace: "nowrap",
                cursor: "pointer",
                transition: "none",
            },
            "& .pill.active": {
                color: token.colorPrimary,
                background: token.colorPrimaryBg,
            },
        },
        logListWrap: {
            display: "flex",
            flexDirection: "column",
            minHeight: 0,
            paddingRight: 1,
        },
        logList: {
            display: "flex",
            flexDirection: "column",
            gap: 0,
            overflow: "hidden",
            background: token.colorBgContainer,
            borderTop: `1px solid ${token.colorBorderSecondary}`,
            borderBottom: `1px solid ${token.colorBorderSecondary}`,
        },
        logEmpty: {
            height: 260,
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
        },
        logCard: {
            position: "relative",
            minWidth: 0,
            padding: isCompactMode ? "8px 2px 8px 12px" : "9px 2px 9px 14px",
            borderBottom: `1px solid ${token.colorBorderSecondary}`,
            background: "transparent",
            transition: "background-color 0.15s",
            "&:before": {
                content: "''",
                position: "absolute",
                top: 10,
                bottom: 10,
                left: 0,
                width: 3,
                borderRadius: 3,
                background: token.colorSuccess,
            },
            "&.outcome:before": {
                background: token.colorError,
            },
            "&:last-child": {
                borderBottom: "none",
            },
            "&:hover": {
                background: token.colorFillQuaternary,
            },
            "& .head": {
                display: "flex",
                alignItems: "center",
                justifyContent: "space-between",
                gap: 10,
                minWidth: 0,
            },
            "& .title": {
                display: "flex",
                alignItems: "center",
                gap: 6,
                minWidth: 0,
            },
            "& .marker": {
                display: "none",
            },
            "& .type": {
                minWidth: 0,
                color: token.colorText,
                fontSize: 13,
                fontWeight: 500,
                lineHeight: "19px",
                overflow: "hidden",
                textOverflow: "ellipsis",
                whiteSpace: "nowrap",
            },
            "& .amount": {
                flex: "none",
                color: token.colorSuccess,
                fontSize: 15,
                fontWeight: 600,
                lineHeight: "20px",
                fontVariantNumeric: "tabular-nums",
                whiteSpace: "nowrap",
            },
            "& .amount.outcome": {
                color: token.colorError,
            },
            "& .meta": {
                display: "flex",
                alignItems: "center",
                justifyContent: "space-between",
                gap: 8,
                minWidth: 0,
                marginTop: 1,
            },
            "& .desc": {
                flex: 1,
                minWidth: 0,
                color: token.colorTextSecondary,
                fontSize: 12,
                lineHeight: "17px",
                overflow: "hidden",
                textOverflow: "ellipsis",
                whiteSpace: "nowrap",
            },
            "& .time": {
                flex: "none",
                color: token.colorTextTertiary,
                fontSize: 12,
                lineHeight: "17px",
                fontVariantNumeric: "tabular-nums",
                whiteSpace: "nowrap",
            },
            "& .change": {
                display: "flex",
                alignItems: "center",
                justifyContent: "space-between",
                gap: 5,
                minWidth: 0,
                marginTop: 5,
                color: token.colorTextTertiary,
                fontSize: 12,
                lineHeight: "17px",
                fontVariantNumeric: "tabular-nums",
            },
            "& .change-values": {
                display: "flex",
                flexWrap: "wrap",
                alignItems: "center",
                gap: 5,
                minWidth: 0,
            },
            "& .change-label": {
                flex: "none",
            },
            "& .change-value": {
                minWidth: 0,
                color: token.colorTextSecondary,
                overflow: "hidden",
                textOverflow: "ellipsis",
                whiteSpace: "nowrap",
            },
            "& .change-part": {
                display: "flex",
                alignItems: "center",
                gap: 4,
                minWidth: 0,
                maxWidth: "100%",
                padding: "1px 6px",
                borderRadius: 4,
                background: token.colorFillQuaternary,
            },
            "@media (max-width: 575px)": {
                padding: "8px 0 8px 12px",
                "& .amount": {
                    fontSize: 13,
                    lineHeight: "18px",
                },
                "& .meta": {
                    gap: 6,
                },
                "& .change": {
                    gap: 4,
                },
                "& .change-values": {
                    gap: 4,
                },
                "& .change-part": {
                    padding: "1px 5px",
                },
                "& .time": {
                    maxWidth: "42%",
                    overflow: "hidden",
                    textOverflow: "ellipsis",
                    whiteSpace: "nowrap",
                },
            },
        },
        pagination: {
            flex: "none",
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            gap: isCompactMode ? 10 : 8,
            marginTop: isCompactMode ? 14 : 12,
            color: token.colorTextTertiary,
            fontSize: 12,
            lineHeight: "20px",
            whiteSpace: "nowrap",
            "@media (max-width: 575px)": {
                alignItems: "center",
                flexDirection: "row",
                gap: 6,
                marginTop: isCompactMode ? 10 : 10,
                "& > span": {
                    flex: "none",
                },
                "& .ant-pagination": {
                    flex: "none",
                    display: "flex",
                    flexWrap: "nowrap",
                    alignItems: "center",
                },
                "& .ant-pagination-prev, & .ant-pagination-next, & .ant-pagination-item": {
                    marginInlineEnd: 2,
                },
            },
        },
    };
});

// ==================== 工具函数 ====================

const getLogFilters = (type: string, range: LogTimeRange, custom?: any) => {
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

// ==================== 通用筛选条 ====================

const FilterScroller = ({options, value, onChange, renderItem, styles}: FilterScrollerProps) => {
    const scrollRef = useRef<HTMLDivElement | null>(null);
    const scroll = (distance: number) => scrollRef.current?.scrollBy({left: distance});

    return (
        <div className={styles.filterRow}>
            <Button
                className={`${styles.filterScrollButton} scrollButton left`}
                htmlType="button"
                icon={<Icon type="CaretLeftFilled"/>}
                size="small"
                title="向左滚动"
                type="text"
                onClick={() => scroll(-160)}
            />
            <div ref={scrollRef} className={styles.filterScroll}>
                {options.map((item) => {
                    const itemValue = String(item.value);
                    const node = (
                        <button
                            className={`pill ${value === itemValue ? "active" : ""}`}
                            type="button"
                            onClick={() => onChange(itemValue)}>
                            {String(item.label)}
                        </button>
                    );
                    return (
                        <React.Fragment key={itemValue}>
                            {renderItem ? renderItem(item, itemValue, node) : node}
                        </React.Fragment>
                    );
                })}
            </div>
            <Button
                className={`${styles.filterScrollButton} scrollButton right`}
                htmlType="button"
                icon={<Icon type="CaretRightFilled"/>}
                size="small"
                title="向右滚动"
                type="text"
                onClick={() => scroll(160)}
            />
        </div>
    );
};

// ==================== 主弹窗 ====================

const Balance = forwardRef<BalanceRef, BalanceProps>(({control, onSuccess}, ref) => {
    const {message} = App.useApp();
    const theme = useTheme();
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const [form] = Form.useForm();
    const [enumData] = useEnum(["pay.balance_log"]) as any;
    const [userInfo, setUserInfo] = useState<any>();
    const [balanceInfo, setBalanceInfo] = useState<any>();
    const [activePanel, setActivePanel] = useState<ActivePanel>("overview");
    const [operateAccount, setOperateAccount] = useState<OperateAccount>("money");
    const [loading, setLoading] = useState(false);
    const [submitting, setSubmitting] = useState(false);

    const isCompactMode = theme?.isCompactTheme?.() || false;
    const {styles} = useStyles({isCompactMode});
    const permissions: ControlState = {
        info: control.hasControl(AdminControlCode.WEB_USER_BALANCE_INFO),
        operate: control.hasControl(AdminControlCode.WEB_USER_BALANCE_OPERATE),
        log: control.hasControl(AdminControlCode.WEB_USER_BALANCE_LOG),
    };
    const balanceLogType = enumData?.["pay.balance_log"]?.type || {};

    const resetOperateForm = (clear = false) => {
        if (clear) {
            form.resetFields();
        }
        form.setFieldsValue({operation_type: DEFAULT_OPERATION_TYPE});
    };

    const loadBalanceInfo = (userId: number) => {
        setLoading(true);
        return getUserBalanceInfo({
            body: {user_id: userId},
            onSuccess: (res: any) => setBalanceInfo(res?.data || {}),
            onFail: (res: any) => message?.error(res?.message || "获取余额详情失败"),
        }).finally(() => setLoading(false));
    };

    const reloadBalanceInfo = (userId: number) => {
        loadBalanceInfo(userId).catch(() => undefined);
    };

    // 只有操作表单挂载后才触发表单方法，避免 useForm 未连接警告。
    useEffect(() => {
        if (activePanel === "operate") {
            resetOperateForm(true);
        }
    }, [activePanel]);

    const open = (record: any) => {
        const userId = Number(record?.id);
        if (!defaultPanel) {
            message?.warning("暂无财务权限");
            return;
        }
        if (!userId) {
            message?.warning("用户ID错误");
            return;
        }
        setUserInfo(record || {});
        setBalanceInfo(undefined);
        setActivePanel(defaultPanel);
        setOperateAccount("money");
        modalRef.current?.show();
        if (permissions.info || permissions.operate) {
            reloadBalanceInfo(userId);
        }
    };

    useImperativeHandle(ref, () => ({open}));

    const changePanel = (value: any) => {
        const nextPanel = value as ActivePanel;
        setActivePanel(nextPanel);
        if ((nextPanel === "overview" || nextPanel === "operate") && !balanceInfo && userInfo?.id) {
            reloadBalanceInfo(Number(userInfo.id));
        }
    };

    const changeOperateAccount = (account: OperateAccount) => {
        setOperateAccount(account);
        if (activePanel === "operate") {
            resetOperateForm();
        }
    };

    const showOperate = (account: OperateAccount) => {
        if (!permissions.operate) {
            message?.warning("暂无余额操作权限");
            return;
        }
        changeOperateAccount(account);
        setActivePanel("operate");
    };

    const submitOperate = async (values: any) => {
        const userId = Number(userInfo?.id);
        const amount = Math.abs(Number(values?.money || 0));
        if (!permissions.operate) {
            message?.warning("暂无余额操作权限");
            return;
        }
        if (!userId) {
            message?.warning("用户ID错误");
            return;
        }
        setSubmitting(true);
        await operateUserBalance({
            body: {
                user_id: userId,
                money: values?.operation_type === "deduct" ? -amount : amount,
                is_operate_frozen: operateAccount === "frozen",
                describe: values?.describe,
            },
            onSuccess: (res: any) => {
                resetOperateForm(true);
                reloadBalanceInfo(userId);
                onSuccess?.();
                message?.success(res?.message || "操作成功");
            },
            onFail: (res: any) => message?.error(res?.message || "操作失败"),
        }).finally(() => setSubmitting(false));
    };

    const panels: PanelConfig[] = [
        {
            label: "概览",
            value: "overview" as ActivePanel,
            visible: permissions.info,
            render: () => (
                <BalanceOverview
                    balanceInfo={balanceInfo}
                    loading={loading}
                    canOperate={permissions.operate}
                    styles={styles}
                    onOperate={showOperate}
                />
            ),
        },
        {
            label: "操作",
            value: "operate" as ActivePanel,
            visible: permissions.operate,
            render: () => (
                <BalanceOperate
                    account={operateAccount}
                    balanceInfo={balanceInfo}
                    form={form}
                    styles={styles}
                    onAccountChange={changeOperateAccount}
                    onSubmit={submitOperate}
                />
            ),
        },
        {
            label: "明细",
            value: "log" as ActivePanel,
            visible: permissions.log,
            render: () => (
                <BalanceLog
                    userId={userInfo?.id}
                    hasControl={permissions.log}
                    typeEnum={balanceLogType}
                    styles={styles}
                    message={message}
                />
            ),
        },
    ].filter((item) => item.visible);
    const defaultPanel = panels[0]?.value;
    const panelOptions = panels.map(({label, value}) => ({label, value}));
    const activePanelRender = panels.find((item) => item.value === activePanel)?.render;

    return (
        <ProModal
            ref={modalRef}
            title={<Title>财务信息</Title>}
            onOk={() => activePanel === "operate" ? form.submit() : modalRef.current?.hide()}
            width={500}
            modalProps={{
                loading: loading && activePanel === "operate",
                confirmLoading: submitting,
                footer: activePanel === "operate" ? undefined : null,
                okText: "提交",
                cancelText: "取消",
                forceRender: true,
                style: {
                    top: 100,
                    paddingBottom: 100,
                },
                mask: {
                    closable: true,
                },
                styles: {
                    body: {
                        padding: isCompactMode ? "4px 8px 10px" : "2px 4px 8px",
                    },
                },
                afterClose: () => {
                    setUserInfo(undefined);
                    setBalanceInfo(undefined);
                    setActivePanel("overview");
                    setOperateAccount("money");
                    setLoading(false);
                    setSubmitting(false);
                },
            } as any}>
            <div className={styles.modalBody}>
                <div className={styles.panelHeader}>
                    <Segmented className={styles.segment} value={activePanel} options={panelOptions}
                               onChange={changePanel}/>
                </div>
                {activePanelRender?.()}
            </div>
        </ProModal>
    );
});

// ==================== 概览模块 ====================

const BalanceOverview = ({balanceInfo, loading, canOperate, styles, onOperate}: BalanceOverviewProps) => {
    if (loading) {
        return (
            <div className={styles.loading}>
                <Spin/>
            </div>
        );
    }

    return (
        <div className={styles.overview}>
            <div className="balance-grid">
                {BALANCE_ITEMS.map((item) => (
                    <div className="balance-card" key={item.field}>
                        <div className="head">
                            <div className="label">{item.label}</div>
                            {canOperate && (
                                <Button type="link" size="small" className="action"
                                        onClick={() => onOperate(item.account)}>
                                    操作
                                </Button>
                            )}
                        </div>
                        <div className="value">{formatMoney(balanceInfo?.[item.field])}</div>
                    </div>
                ))}
            </div>
            <div className="stats">
                {STAT_ITEMS.map((item) => (
                    <div className="item" key={item.field}>
                        <div className="label">{item.label}</div>
                        <div className="value">{formatMoney(balanceInfo?.[item.field])}</div>
                    </div>
                ))}
            </div>
            <div className="overview-footer">
                <span>最后操作</span>
                <span className="value">{balanceInfo?.last_operate_time || "-"}</span>
            </div>
        </div>
    );
};

// ==================== 操作模块 ====================

const BalanceOperate = ({account, balanceInfo, form, styles, onAccountChange, onSubmit}: BalanceOperateProps) => {
    const accountInfo = account === "frozen"
        ? {label: "冻结金额", value: balanceInfo?.frozen}
        : {label: "可用余额", value: balanceInfo?.money};
    const operationOptions: SelectOption<OperationType>[] = account === "frozen" ? [
        {label: "冻结", value: "add"},
        {label: "释放", value: "deduct"},
    ] : [
        {label: "加款", value: "add"},
        {label: "扣除", value: "deduct"},
    ];

    return (
        <Form
            className={styles.operate}
            form={form}
            layout="vertical"
            onFinish={onSubmit}>
            <div className="account-switch">
                {ACCOUNT_OPTIONS.map((item) => (
                    <button
                        key={item.value}
                        className={`account-option ${account === item.value ? "active" : ""}`}
                        type="button"
                        onClick={() => onAccountChange(item.value)}>
                        <div className="account-head">
                            <span className="account-label">{item.label}</span>
                            {account === item.value && <span className="account-current">当前</span>}
                        </div>
                        <div className="account-value">
                            {formatMoney(item.value === "frozen" ? balanceInfo?.frozen : balanceInfo?.money)}
                        </div>
                    </button>
                ))}
            </div>
            <div className="form-grid">
                <div className="form-field">
                    <div className="field-head">
                        <span className="field-label">操作类型</span>
                    </div>
                    <Form.Item name="operation_type" initialValue={DEFAULT_OPERATION_TYPE} noStyle>
                        <Segmented block className={styles.segment} options={operationOptions}/>
                    </Form.Item>
                </div>
                <div className="form-field">
                    <div className="field-head">
                        <span className="field-label"><span className="required">*</span>操作金额</span>
                    </div>
                    <Form.Item name="money" rules={[
                        {required: true, message: "请输入操作金额"},
                        {
                            validator: (_: any, value: any) => {
                                if (value === undefined || value === null || value === "") {
                                    return Promise.resolve();
                                }
                                return Number(value) >= 0
                                    ? Promise.resolve()
                                    : Promise.reject(new Error("操作金额不能小于0"));
                            },
                        },
                    ]}>
                        <InputNumber
                            min={0}
                            precision={2}
                            step={1}
                            placeholder={`请输入${accountInfo.label}操作金额`}
                        />
                    </Form.Item>
                </div>
                <div className="form-field full">
                    <div className="field-head">
                        <span className="field-label">操作描述</span>
                    </div>
                    <Form.Item name="describe">
                        <Input.TextArea
                            className="textarea"
                            placeholder="可填写操作说明"
                            rows={3}
                            maxLength={120}
                        />
                    </Form.Item>
                </div>
            </div>
        </Form>
    );
};

// ==================== 明细模块 ====================

const BalanceLog = ({userId, hasControl, typeEnum, styles, message}: BalanceLogProps) => {
    const [filters, setFilters] = useState<any>({});
    const [type, setType] = useState<string>(ALL_VALUE);
    const [timeRange, setTimeRange] = useState<LogTimeRange>("all");
    const [customOpen, setCustomOpen] = useState(false);
    const [customRange, setCustomRange] = useState<any>();
    const customRangeRef = useRef<any[]>([]);
    const [list, setList] = useState<any[]>([]);
    const [page, setPage] = useState(1);
    const [total, setTotal] = useState(0);
    const [loading, setLoading] = useState(false);
    const typeOptions: SelectOption[] = [
        {label: "全部", value: ALL_VALUE},
        ...Object.entries(typeEnum || {}).map(([value, label]) => ({value, label: String(label)})),
    ];

    const changeType = (value: string) => {
        setType(value);
        setPage(1);
        setFilters(getLogFilters(value, timeRange, customRangeRef.current));
    };

    const changeTimeRange = (value: LogTimeRange) => {
        setPage(1);
        if (value === "custom") {
            setCustomOpen(true);
            return;
        }
        setTimeRange(value);
        setCustomOpen(false);
        customRangeRef.current = [];
        setCustomRange(undefined);
        setFilters(getLogFilters(type, value, undefined));
    };

    const applyCustomTime = (value: any, index: 0 | 1) => {
        const next = [...(customRangeRef.current || [])];
        setPage(1);
        if (!value) {
            next[index] = undefined;
            customRangeRef.current = next;
            setCustomRange(next);
            if (next?.[0] || next?.[1]) {
                setTimeRange("custom");
                setFilters(getLogFilters(type, "custom", next));
            } else {
                setTimeRange("all");
                setFilters(getLogFilters(type, "all", undefined));
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
        setFilters(getLogFilters(type, "custom", range));
        if (range?.[0] && range?.[1]) {
            setCustomOpen(false);
        }
    };

    useEffect(() => {
        if (!hasControl || !userId) {
            return;
        }
        let disposed = false;
        setLoading(true);
        const request = getUserBalanceLogList({
            body: {
                user_id: Number(userId),
                page,
                page_size: PAGE_SIZE,
                order_by_field: "id",
                order_by_type: "desc",
                ...filters,
            },
            onFail: (res: any) => message?.error(res?.message || "获取余额明细失败"),
        }).then((res: any) => {
            if (disposed) {
                return;
            }
            setList(res?.data?.list || []);
            setTotal(Number(res?.data?.total || 0));
        }).finally(() => {
            if (!disposed) {
                setLoading(false);
            }
        });
        request.catch(() => undefined);
        return () => {
            disposed = true;
        };
    }, [filters, hasControl, message, page, userId]);

    const customPickerPanel = (
        <div className={styles.customPickerPanel}>
            <div className={styles.customPickerGroup}>
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

    return (
        <div className={`${styles.logPanel} ${list.length > 0 ? "" : "empty"}`}>
            <div className={styles.logFilter}>
                <FilterScroller options={typeOptions} value={type} styles={styles} onChange={changeType}/>
                <FilterScroller
                    options={LOG_TIME_OPTIONS}
                    value={timeRange}
                    styles={styles}
                    onChange={(value: string) => changeTimeRange(value as LogTimeRange)}
                    renderItem={(_: any, itemValue: string, node: React.ReactNode) => itemValue === "custom" ? (
                        <Popover
                            arrow
                            content={customPickerPanel}
                            open={customOpen}
                            placement="bottomLeft"
                            trigger="click"
                            onOpenChange={(open) => setCustomOpen(open)}>
                            {node}
                        </Popover>
                    ) : node}
                />
            </div>
            {list.length > 0 ? (
                <div className={styles.logListWrap}>
                    <Spin spinning={loading}>
                        <div className={styles.logList}>
                            {list.map((item: any) => {
                                const money = Number(item?.money || 0);
                                const amountText = `${money > 0 ? "+" : ""}${formatMoney(item?.money)}`;
                                return (
                                    <div
                                        key={item?.id || `${item?.create_time}-${item?.money}`}
                                        className={`${styles.logCard} ${money >= 0 ? "income" : "outcome"}`}>
                                        <div className="head">
                                            <div className="title">
                                                <span className="marker"/>
                                                <div
                                                    className="type">{typeEnum?.[item?.type] || item?.type || "-"}</div>
                                            </div>
                                            <div className={`amount ${money >= 0 ? "income" : "outcome"}`}>
                                                {amountText}
                                            </div>
                                        </div>
                                        <div className="meta">
                                            <span className="desc">{item?.describe || "暂无描述"}</span>
                                        </div>
                                        <div className="change">
                                            <span className="change-values">
                                                <span className="change-part">
                                                    <span className="change-label">余额</span>
                                                    <span className="change-value">
                                                        {formatMoney(item?.money_before)} -&gt; {formatMoney(item?.money_after)}
                                                    </span>
                                                </span>
                                                <span className="change-part">
                                                    <span className="change-label">冻结</span>
                                                    <span className="change-value">
                                                        {formatMoney(item?.frozen_before)} -&gt; {formatMoney(item?.frozen_after)}
                                                    </span>
                                                </span>
                                            </span>
                                            <span className="time">{item?.create_time || "-"}</span>
                                        </div>
                                    </div>
                                );
                            })}
                        </div>
                    </Spin>
                </div>
            ) : (
                <Spin spinning={loading}>
                    <div className={styles.logEmpty}>
                        <Empty image={Empty.PRESENTED_IMAGE_SIMPLE}/>
                    </div>
                </Spin>
            )}
            {total > 0 && (
                <div className={styles.pagination}>
                    <span>共 {total} 条</span>
                    <Pagination
                        size="small"
                        current={page}
                        pageSize={PAGE_SIZE}
                        total={total}
                        showSizeChanger={false}
                        onChange={setPage}
                    />
                </div>
            )}
        </div>
    );
};

Balance.displayName = "Balance";

export default Balance;
