import React, {forwardRef, useImperativeHandle, useRef, useState} from "react";
import {App, Grid, Tag, Typography} from "antd";
import {createStyles} from "antd-style";
import type {ModalTableProps, ModalTableRef} from "sinking-antd";
import {ModalTable, Title} from "sinking-antd";
import {emptyTableData, getParams} from "@/utils/page";
import {getMasterLogList} from "@/service/admin/master";
import {AdminControlCode} from "@/utils/permission";

export type MasterLogRef = {
    open: (record: any) => void;
};

type MasterLogProps = {
    control: any;
    typeData: Record<string, string>;
};

const useStyles = createStyles(({css, token}) => ({
    table: css`
        && .ant-table-thead > tr > th,
        && .ant-table-tbody > tr > td,
        && .ant-table-tbody .ant-typography {
            font-size: ${token.fontSizeSM + 1}px;
        }

        .log-info {
            display: flex;
            flex-direction: column;
            align-items: flex-start;
            gap: 2px;
            min-width: 0;
        }

        .log-info .ant-typography {
            line-height: 1.5;
            overflow-wrap: anywhere;
        }

        && .log-location.ant-typography,
        && .log-request.ant-typography {
            font-size: ${token.fontSizeSM}px;
        }

        .log-time {
            color: ${token.colorTextSecondary};
            font-variant-numeric: tabular-nums;
            white-space: nowrap;
        }

        .ant-typography-copy {
            margin-inline-start: 6px;
            font-size: ${token.fontSizeSM}px;
            color: ${token.colorTextTertiary};
        }

        .ant-typography-copy:hover { color: ${token.colorPrimary}; }
    `,
    typeTag: css`
        && {
            margin-inline-end: 0;
            color: ${token.colorPrimaryText};
            background: ${token.colorPrimaryBg};
        }
    `,
}));

const MasterLog = forwardRef<MasterLogRef, MasterLogProps>(({control, typeData}, ref) => {
    const {message} = App.useApp();
    const {styles} = useStyles();
    const screens = Grid.useBreakpoint();
    const tableRef = useRef<ModalTableRef | null>(null);
    const [userInfo, setUserInfo] = useState<any>();
    const hasLogControl = control.hasControl(AdminControlCode.SYSTEM_MASTER_LOG);

    useImperativeHandle(ref, () => ({
        open: (record: any) => {
            if (!hasLogControl) {
                message?.warning("暂无日志查看权限");
                return;
            }
            setUserInfo(record || {});
            tableRef.current?.show();
        },
    }));

    const columns: ModalTableProps["columns"] = [
        {
            title: "日志ID / 请求ID",
            dataIndex: "id",
            width: 300,
            sorter: true,
            render: (value: number, record: any) => (
                <div className="log-info">
                    <Typography.Text copyable={value ? {text: String(value)} : false}>
                        {value || "-"}
                    </Typography.Text>
                    <Typography.Text className="log-request" type="secondary"
                                     copyable={record?.request_id ? {text: record.request_id} : false}>
                        {record?.request_id || "-"}
                    </Typography.Text>
                </div>
            ),
        },
        {
            title: "操作类型",
            dataIndex: "type",
            width: 120,
            render: (value: number) => (
                <Tag className={styles.typeTag} variant="filled">
                    {typeData[value] || "-"}
                </Tag>
            ),
        },
        {
            title: "操作信息",
            dataIndex: "title",
            render: (value: string, record: any) => (
                <div className="log-info">
                    <Typography.Text>{value || "-"}</Typography.Text>
                    <Typography.Text type="secondary">{record?.content || "-"}</Typography.Text>
                </div>
            ),
        },
        {
            title: "请求来源",
            dataIndex: "request_ip",
            render: (value: string, record: any) => (
                <div className="log-info">
                    <Typography.Text copyable={value ? {text: value} : false}>{value || "-"}</Typography.Text>
                    <Typography.Text className="log-location" type="secondary">
                        {record?.location || "-"}
                    </Typography.Text>
                </div>
            ),
        },
        {
            title: "操作时间",
            dataIndex: "create_time",
            sorter: true,
            render: (value: string) => <span className="log-time">{value || "-"}</span>,
        },
    ];

    return (
        <ModalTable
            ref={tableRef}
            className={styles.table}
            title={<Title>操作日志</Title>}
            width={1000}
            minHeight={440}
            scrollMode="modal"
            modalProps={{
                mask: {closable: true},
                style: {
                    top: screens.md ? 100 : 24,
                    paddingBottom: screens.md ? 100 : 24,
                    maxWidth: "calc(100vw - 24px)",
                },
                afterClose: () => setUserInfo(undefined),
            }}
            rowKey="id"
            columns={columns}
            scroll={{x: 950}}
            defaultPageSize={7}
            params={{id: userInfo?.id}}
            requestEnabled={hasLogControl && !!userInfo?.id}
            request={async (params, sort) => {
                if (!hasLogControl || !userInfo?.id) {
                    return emptyTableData;
                }
                const res = await getMasterLogList({body: getParams({...params}, sort)});
                if (res?.code !== 200) {
                    throw new Error(res?.message || "获取操作日志失败");
                }
                return {data: res.data?.list || [], total: res.data?.total || 0, success: true};
            }}
            search={{
                placeholder: "搜索操作日志",
                ariaLabel: "搜索操作日志",
                maxLength: 100,
                fields: [
                    {value: "title", label: "标题", placeholder: "搜索操作标题"},
                    {value: "content", label: "内容", placeholder: "搜索操作内容"},
                    {value: "request_ip", label: "IP", placeholder: "搜索请求 IP", maxLength: 50},
                    {value: "location", label: "地点", placeholder: "搜索请求地点"},
                    {value: "request_id", label: "请求ID", placeholder: "搜索请求 ID"},
                ],
            }}
            filters={[
                {type: "select", name: "type", label: "操作类型筛选", allLabel: "全部类型", valueEnum: typeData},
                {type: "date", name: "create_time", label: "操作时间筛选"},
            ]}
            toolbar={{refresh: {ariaLabel: "刷新操作日志列表"}}}
            pagination={{unit: "条"}}/>
    );
});

MasterLog.displayName = "MasterLog";

export default MasterLog;
