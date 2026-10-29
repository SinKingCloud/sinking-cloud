import React, {forwardRef, useEffect, useImperativeHandle, useRef, useState} from "react";
import {App, Grid} from "antd";
import {createStyles} from "antd-style";
import type {ModalTableProps, ModalTableRef} from "sinking-antd";
import {Icon, ModalTable, Title} from "sinking-antd";
import {getCronLog} from "@/service/admin/cron";
import {AdminControlCode} from "@/utils/permission";

export type CronLogRef = {
    open: (record: any) => void;
};

type CronLogProps = {
    control: any;
};

type LogCursor = {
    id: string;
    lastId: string;
};

const pageSize = 10;

const useStyles = createStyles(({css, token}) => ({
    footer: css`
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 12px;
        min-height: ${token.controlHeight}px;
        text-align: start;

        .log-summary {
            display: flex;
            align-items: center;
            gap: 8px;
            color: ${token.colorTextSecondary};
            font-size: ${token.fontSizeSM + 1}px;
            line-height: 20px;
            white-space: nowrap;
        }

        .log-current, .log-capacity {
            display: inline-flex;
            align-items: baseline;
            gap: 4px;
        }

        .log-current strong {
            font-weight: 500;
            font-variant-numeric: tabular-nums;
        }

        .log-capacity {
            color: ${token.colorTextSecondary};
            font-size: ${token.fontSizeSM}px;
            font-variant-numeric: tabular-nums;
        }

        .log-summary-divider {
            color: ${token.colorTextQuaternary};
        }

        @media (max-width: 575px) {
            gap: 8px;

            .log-summary { flex-direction: column; align-items: flex-start; gap: 0; }
            .log-summary-divider { display: none; }
        }
    `,
    navigation: css`
        flex: none;
        display: flex;
        align-items: center;
        gap: 8px;
    `,
    pageButton: css`
        appearance: none;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        gap: 6px;
        min-width: 86px;
        height: ${token.controlHeightSM + 8}px;
        margin: 0;
        padding: 0 10px;
        border: 0;
        border-radius: ${token.borderRadius}px;
        outline: none;
        color: ${token.colorTextSecondary};
        background: ${token.colorFillQuaternary};
        box-shadow: none;
        font: inherit;
        font-size: ${token.fontSizeSM + 1}px;
        line-height: 1;
        white-space: nowrap;
        cursor: pointer;
        user-select: none;
        touch-action: manipulation;

        .anticon { font-size: ${token.fontSizeSM}px; }

        &:not(:disabled):hover,
        &:not(:disabled):active { color: ${token.colorText}; }

        &:focus-visible {
            text-decoration: underline;
            text-underline-offset: 4px;
        }

        &:disabled { color: ${token.colorTextQuaternary}; cursor: not-allowed; }

        @media (max-width: 575px) {
            min-width: 78px;
            padding-inline: 8px;
        }
    `,
}));

const CronLog = forwardRef<CronLogRef, CronLogProps>(({control}, ref) => {
    const {message} = App.useApp();
    const {styles} = useStyles();
    const screens = Grid.useBreakpoint();
    const tableRef = useRef<ModalTableRef | null>(null);
    const requestRef = useRef(0);
    const loadingRef = useRef(false);
    const cursorsRef = useRef<LogCursor[]>([{id: "", lastId: ""}]);
    const [taskId, setTaskId] = useState<number>();
    const [loading, setLoading] = useState(false);
    const [loadFailed, setLoadFailed] = useState(false);
    const [items, setItems] = useState<any[]>([]);
    const [sort, setSort] = useState<Record<string, "ascend" | "descend">>({});
    const [page, setPage] = useState({index: 0, hasMore: false, nextId: "", nextLastId: ""});
    const hasLogControl = control.hasControl(AdminControlCode.SYSTEM_CRON_LOG);
    const retry = loadFailed && items.length === 0;

    useEffect(() => () => {
        requestRef.current += 1;
    }, []);

    const loadLog = async (cursor: LogCursor = {id: "", lastId: ""}, index = 0, nextSort = sort) => {
        if (!hasLogControl || !taskId) {
            return;
        }
        const requestId = ++requestRef.current;
        loadingRef.current = true;
        setLoading(true);
        setLoadFailed(false);
        const [field, order] = Object.entries(nextSort)[0] || ["id", "descend"];
        try {
            const res = await getCronLog({
                body: {
                    id: taskId,
                    page_size: pageSize,
                    order_by_field: field,
                    order_by_type: order === "ascend" ? "asc" : "desc",
                    ...(cursor.id ? {cursor_id: cursor.id} : {}),
                    ...(cursor.lastId ? {cursor_last_id: cursor.lastId} : {}),
                },
            });
            if (requestId !== requestRef.current) {
                return;
            }
            if (res?.code !== 200) {
                throw new Error(res?.message || "获取任务日志失败");
            }
            const data = res?.data || {};
            setItems(Array.isArray(data?.list) ? data?.list?.filter(Boolean) : []);
            setSort(nextSort);
            cursorsRef.current = [...cursorsRef.current.slice(0, index), cursor];
            setPage({
                index,
                hasMore: Boolean(data?.has_more),
                nextId: data?.next_cursor_id || "",
                nextLastId: data?.next_cursor_last_id || "",
            });
        } catch (error) {
            if (requestId === requestRef.current) {
                setLoadFailed(true);
                message?.error(error instanceof Error ? error.message : "获取任务日志失败");
            }
        } finally {
            if (requestId === requestRef.current) {
                loadingRef.current = false;
                setLoading(false);
            }
        }
    };

    useImperativeHandle(ref, () => ({
        open: (record: any) => {
            if (!hasLogControl) {
                message?.warning("暂无任务日志权限");
                return;
            }
            const id = Number(record?.id);
            if (!id) {
                return;
            }
            requestRef.current += 1;
            loadingRef.current = true;
            cursorsRef.current = [{id: "", lastId: ""}];
            setLoading(true);
            setLoadFailed(false);
            setItems([]);
            setSort({});
            setPage({index: 0, hasMore: false, nextId: "", nextLastId: ""});
            setTaskId(id);
            tableRef.current?.show();
        },
    }));

    const columns: ModalTableProps["columns"] = [
        {
            title: "ID",
            dataIndex: "id",
            width: 180,
            sorter: true,
            sortOrder: sort.id || null,
        },
        {
            title: "日志内容",
            dataIndex: "content",
            render: (value: string) => value || "-",
        },
        {
            title: "创建时间",
            dataIndex: "create_time",
            width: 180,
            sorter: true,
            sortOrder: sort.create_time || null,
        },
    ];

    return (
        <ModalTable
            ref={tableRef}
            title={<Title>任务日志</Title>}
            width={760}
            minHeight={440}
            scrollMode="modal"
            onOpenChange={(open) => {
                if (!open) {
                    requestRef.current += 1;
                    loadingRef.current = false;
                    setLoading(false);
                    setTaskId(undefined);
                }
            }}
            modalProps={{
                style: {
                    top: screens.md ? 100 : 24,
                    paddingBottom: screens.md ? 100 : 24,
                    maxWidth: "calc(100vw - 24px)",
                },
                mask: {closable: true},
                footer: (
                    <div className={styles.footer}>
                        <div className="log-summary" role="status" aria-live="polite">
                            <span className="log-current">第 <strong>{page.index + 1}</strong> 页</span>
                            <span className="log-summary-divider" aria-hidden="true">·</span>
                            <span className="log-capacity">每页 {pageSize} 条</span>
                        </div>
                        <nav className={styles.navigation} aria-label="任务日志分页">
                            <button type="button" className={styles.pageButton}
                                    disabled={loading || page.index === 0} onClick={() => {
                                if (!loadingRef.current && page.index > 0) {
                                    void loadLog(cursorsRef.current[page.index - 1], page.index - 1);
                                }
                            }}><Icon type="LeftOutlined"/><span>上一页</span></button>
                            <button type="button" className={styles.pageButton}
                                    disabled={loading || (!page.hasMore && !retry)} onClick={() => {
                                if (loadingRef.current) {
                                    return;
                                }
                                if (retry) {
                                    void loadLog();
                                } else if (page.hasMore) {
                                    void loadLog({id: page.nextId, lastId: page.nextLastId}, page.index + 1);
                                }
                            }}><span>{retry ? "重试" : "下一页"}</span>{!retry && <Icon type="RightOutlined"/>}</button>
                        </nav>
                    </div>
                ),
            }}
            search={false}
            rowKey="id"
            columns={columns}
            scroll={{x: 700}}
            dataSource={items}
            loading={loading}
            requestEnabled={hasLogControl && !!taskId}
            onReload={() => loadLog()}
            onSortChange={(field, order) => {
                void loadLog({id: "", lastId: ""}, 0, field && order ? {[field]: order} : {});
            }}
            toolbar={false}
            pagination={false}/>
    );
});

CronLog.displayName = "CronLog";

export default CronLog;
