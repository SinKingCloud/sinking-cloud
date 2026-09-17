import {forwardRef, memo, useMemo, useState} from "react";
import {Input, InputNumber, Select, Typography} from "antd";
import {createStyles} from "antd-style";
import {Icon, useTheme} from "sinking-antd";
import {describeSchedule, parseSchedule, ScheduleMode} from "../utils";

const scheduleOptions = [
    {label: "每秒", value: "second"},
    {label: "每分", value: "minute"},
    {label: "每小时", value: "hour"},
    {label: "每天", value: "day"},
    {label: "每周", value: "week"},
    {label: "每月", value: "month"},
    {label: "自定义", value: "custom"},
];

const weekdayOptions = [
    {label: "周一", value: 1},
    {label: "周二", value: 2},
    {label: "周三", value: 3},
    {label: "周四", value: 4},
    {label: "周五", value: 5},
    {label: "周六", value: 6},
    {label: "周日", value: 0},
];

const defaultExpressions: Partial<Record<ScheduleMode, string>> = {
    second: "*/5 * * * * *",
    minute: "0 */5 * * * *",
    hour: "0 0 */1 * * *",
    day: "0 0 2 * * *",
    week: "0 0 2 * * 1",
    month: "0 0 2 1 * *",
};

const useStyles = createStyles(({css, token, isDarkMode}, {compact}: {compact: boolean}) => ({
    schedule: css`
        min-width: 0;
        container-type: inline-size;
    `,
    editor: css`
        min-width: 0;
        display: flex;
        align-items: center;
        gap: 8px;

        .cycle-type { width: 132px; flex: none; }
        .cycle-value { min-width: 0; flex: 1; }
        .schedule-number { width: 96px; flex: none; }
        .schedule-minute { width: 88px; flex: none; }
        .schedule-time { width: 140px; flex: none; color-scheme: ${isDarkMode ? "dark" : "light"}; }
        .schedule-weekday { width: 104px; flex: none; }

        @container (max-width: 520px) {
            flex-direction: column;
            align-items: stretch;

            .cycle-type { width: 100%; }
            .schedule-number, .schedule-minute { width: 64px; }
            .schedule-time { width: 124px; }
            .schedule-weekday { width: 88px; }
        }
    `,
    inlineControl: css`
        min-width: 0;
        min-height: ${token.controlHeight}px;
        display: flex;
        align-items: center;
        gap: 8px;
        white-space: nowrap;

        > span { flex: none; }

        @container (max-width: 520px) {
            flex-wrap: wrap;
            gap: 6px;
        }
    `,
    preview: css`
        min-width: 0;
        min-height: ${compact ? 32 : 38}px;
        margin-top: ${compact ? 6 : 8}px;
        padding: ${compact ? "6px 9px" : "8px 11px"};
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: ${compact ? 10 : 12}px;
        border-radius: ${token.borderRadius}px;
        background: ${token.colorFillQuaternary};

        .schedule-description.ant-typography {
            min-width: 0;
            flex: 1;
            display: inline-flex;
            align-items: center;
            gap: 6px;
            margin: 0;
            color: ${token.colorTextSecondary};
            font-size: ${compact ? 11 : 12}px;
            line-height: 20px;
        }

        .schedule-description .anticon {
            flex: none;
            color: ${token.colorTextTertiary};
        }

        .schedule-expression.ant-typography {
            min-width: 0;
            flex: 0 1 auto;
            display: inline-flex;
            align-items: center;
            gap: 6px;
            margin: 0;
            padding-inline-start: ${compact ? 10 : 12}px;
            border-inline-start: 1px solid ${token.colorSplit};
            color: ${token.colorTextTertiary};
            font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
            font-size: ${compact ? 11 : 12}px;
            line-height: 20px;
        }

        .expression-value {
            min-width: 0;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .ant-typography-copy {
            flex: none;
            margin-inline-start: 0;
            color: ${token.colorTextTertiary};
        }

        .ant-typography-copy:hover { color: ${token.colorPrimary}; }

        @container (max-width: 400px) {
            flex-wrap: wrap;
            gap: 4px;

            .schedule-description.ant-typography { flex-basis: 100%; }
            .schedule-expression.ant-typography { padding-inline-start: 0; border: 0; }
        }
    `,
}));

type ScheduleProps = {
    id?: string;
    value?: string;
    onChange?: (value: string) => void;
};

const Schedule = forwardRef<HTMLDivElement, ScheduleProps>(({id, value, onChange}, ref) => {
    const theme = useTheme();
    const {styles} = useStyles({compact: theme?.isCompactTheme?.() || false});
    const [customValue, setCustomValue] = useState<string>();
    const schedule = useMemo(() => parseSchedule(value), [value]);
    const mode = customValue !== undefined && customValue === value ? "custom" : schedule.mode;
    const description = describeSchedule(schedule);
    const timeValue = `${String(schedule.hour ?? 0).padStart(2, "0")}:${String(schedule.minute ?? 0).padStart(2, "0")}`;
    const emit = (expression: string) => onChange?.(expression);

    const changeMode = (nextMode: ScheduleMode) => {
        setCustomValue(nextMode === "custom" ? value : undefined);
        const expression = defaultExpressions[nextMode];
        if (expression) emit(expression);
    };

    const changeTime = (value: string) => {
        if (!/^\d{2}:\d{2}$/.test(value)) return;
        const [hour, minute] = value.split(":").map(Number);
        if (mode === "day") {
            emit(`0 ${minute} ${hour} * * *`);
        } else if (mode === "week") {
            emit(`0 ${minute} ${hour} * * ${schedule.weekday ?? 1}`);
        } else {
            emit(`0 ${minute} ${hour} ${schedule.day ?? 1} * *`);
        }
    };

    return (
        <div id={id} ref={ref} className={styles.schedule}>
            <div className={styles.editor}>
                <Select className="cycle-type" aria-label="执行周期" value={mode}
                        options={scheduleOptions} onChange={changeMode}/>
                <div className="cycle-value">
                    {mode === "second" && (
                        <div className={styles.inlineControl}>
                            <span>每</span>
                            <InputNumber className="schedule-number" aria-label="间隔秒数" min={1} max={59}
                                         precision={0} value={schedule.interval ?? 5}
                                         onChange={(interval) => emit(`*/${interval ?? 1} * * * * *`)}/>
                            <span>秒</span>
                        </div>
                    )}
                    {mode === "minute" && (
                        <div className={styles.inlineControl}>
                            <span>每</span>
                            <InputNumber className="schedule-number" aria-label="间隔分钟数" min={1} max={59}
                                         precision={0} value={schedule.interval ?? 5}
                                         onChange={(interval) => emit(`0 */${interval ?? 1} * * * *`)}/>
                            <span>分钟</span>
                        </div>
                    )}
                    {mode === "hour" && (
                        <div className={styles.inlineControl}>
                            <span>每</span>
                            <InputNumber className="schedule-number" aria-label="间隔小时数" min={1} max={23}
                                         precision={0} value={schedule.interval ?? 1}
                                         onChange={(interval) => emit(`0 ${schedule.minute ?? 0} */${interval ?? 1} * * *`)}/>
                            <span>小时</span>
                            <InputNumber className="schedule-minute" aria-label="执行分钟" min={0} max={59}
                                         precision={0} value={schedule.minute ?? 0}
                                         onChange={(minute) => emit(`0 ${minute ?? 0} */${schedule.interval ?? 1} * * *`)}/>
                            <span>分钟</span>
                        </div>
                    )}
                    {(mode === "day" || mode === "week" || mode === "month") && (
                        <div className={styles.inlineControl}>
                            <span>{mode === "day" ? "每天" : mode === "week" ? "每周" : "每月"}</span>
                            {mode === "week" && (
                                <Select className="schedule-weekday" aria-label="执行星期" value={schedule.weekday ?? 1}
                                        options={weekdayOptions}
                                        onChange={(weekday) => emit(`0 ${schedule.minute ?? 0} ${schedule.hour ?? 0} * * ${weekday}`)}/>
                            )}
                            {mode === "month" && (
                                <>
                                    <InputNumber className="schedule-number" aria-label="每月执行日期" min={1} max={31}
                                                 precision={0} value={schedule.day ?? 1}
                                                 onChange={(day) => emit(`0 ${schedule.minute ?? 0} ${schedule.hour ?? 0} ${day ?? 1} * *`)}/>
                                    <span>日</span>
                                </>
                            )}
                            <Input className="schedule-time" aria-label="执行时间" type="time" step={60}
                                   value={timeValue} onChange={(event) => changeTime(event.target.value)}/>
                        </div>
                    )}
                    {mode === "custom" && (
                        <Input aria-label="自定义表达式" value={value || ""} placeholder="秒 分 时 日 月 星期"
                               onChange={(event) => {
                                   setCustomValue(event.target.value);
                                   emit(event.target.value);
                               }}/>
                    )}
                </div>
            </div>
            <div className={styles.preview}>
                <Typography.Text className="schedule-description">
                    <Icon type="ClockCircleOutlined"/>
                    {description}
                </Typography.Text>
                <Typography.Text className="schedule-expression" copyable={value ? {text: value} : false}>
                    <span className="expression-value" title={value}>{value || "-"}</span>
                </Typography.Text>
            </div>
        </div>
    );
});

Schedule.displayName = "Schedule";

export default memo(Schedule);
