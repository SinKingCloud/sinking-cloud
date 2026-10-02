import React, {forwardRef, useImperativeHandle, useRef, useState} from "react";
import {App, Button, Grid, Typography} from "antd";
import {createStyles} from "antd-style";
import type {ModalTableProps, ModalTableRef} from "sinking-antd";
import {ModalTable, Title} from "sinking-antd";
import {emptyTableData, getData} from "@/utils/page";
import {getMasterSessionList, logoutMasterSession} from "@/service/admin/master";
import {AdminControlCode} from "@/utils/permission";

export type MasterSessionRef = {
    open: (record: any) => void;
};

type MasterSessionProps = {
    control: any;
};

const useStyles = createStyles(({css, token}) => ({
    table: css`
        && .ant-table-thead > tr > th {
            font-size: ${token.fontSizeSM + 1}px;
        }

        && .ant-table-tbody > tr > td,
        && .ant-table-tbody .ant-typography,
        && .ant-table-tbody .ant-btn {
            font-size: ${token.fontSizeSM + 1}px;
        }

        .session-source {
            display: flex;
            flex-direction: column;
            align-items: flex-start;
            gap: 2px;
            min-width: 0;
        }

        .session-source .ant-typography {
            line-height: 1.5;
            overflow-wrap: anywhere;
        }

        && .session-location.ant-typography {
            font-size: ${token.fontSizeSM}px;
        }

        .session-time {
            font-size: ${token.fontSizeSM}px;
            color: ${token.colorTextSecondary};
            font-variant-numeric: tabular-nums;
            white-space: nowrap;
        }

        .session-time-row {
            display: flex;
            align-items: center;
            gap: 4px;
            line-height: 18px;
        }

        .session-time-label {
            flex: none;
            color: ${token.colorTextTertiary};
        }

        .ant-typography-copy {
            margin-inline-start: 6px;
            font-size: ${token.fontSizeSM}px;
            color: ${token.colorTextTertiary};
        }

        .ant-typography-copy:hover { color: ${token.colorPrimary}; }
    `,
}));

const MasterSession = forwardRef<MasterSessionRef, MasterSessionProps>(({control}, ref) => {
    const {message, modal} = App.useApp();
    const {styles} = useStyles();
    const screens = Grid.useBreakpoint();
    const tableRef = useRef<ModalTableRef | null>(null);
    const [userInfo, setUserInfo] = useState<any>();
    const hasSessionControl = control.hasControl(AdminControlCode.SYSTEM_MASTER_SESSION);
    const hasLogoutControl = control.hasControl(AdminControlCode.SYSTEM_MASTER_LOGOUT);

    useImperativeHandle(ref, () => ({
        open: (record: any) => {
            if (!hasSessionControl) {
                message?.warning("暂无会话查看权限");
                return;
            }
            setUserInfo(record || {});
            tableRef.current?.show();
        },
    }));

    const logoutSession = (tokens: any[] = []) => {
        if (!hasLogoutControl) {
            message?.warning("暂无注销会话权限");
            return;
        }
        const nextTokens = Array.from(new Set((tokens || []).map((item) => String(item || "")).filter(Boolean)));
        if (nextTokens.length <= 0) {
            message?.warning("请选择要注销的会话");
            return;
        }
        modal.confirm({
            title: "注销会话",
            content: `确定注销选中的 ${nextTokens.length} 个会话吗？`,
            okText: "确 定",
            okButtonProps: {danger: true, type: "default"},
            cancelText: "取 消",
            mask: {
                closable: true,
            },
            onOk: async () => {
                await logoutMasterSession({
                    body: {
                        id: Number(userInfo?.id),
                        tokens: nextTokens,
                    },
                    onSuccess: (res: any) => {
                        tableRef.current?.reload();
                        tableRef.current?.clearSelectedRows?.();
                        message?.success(res?.message || "注销成功");
                    },
                    onFail: (res: any) => {
                        message?.error(res?.message || "注销失败");
                    },
                });
            },
        } as any);
    };

    const columns: ModalTableProps["columns"] = [
        {
            title: "登录设备",
            dataIndex: "device",
            width: 180,
            render: (value: string) => (
                <Typography.Text copyable={value ? {text: value} : false}>{value || "-"}</Typography.Text>
            ),
        },
        {
            title: "登录来源",
            dataIndex: "ip",
            width: 190,
            render: (value: string, record: any) => (
                <div className="session-source">
                    <Typography.Text className="session-ip" copyable={value ? {text: value} : false}>
                        {value || "-"}
                    </Typography.Text>
                    <Typography.Text className="session-location" type="secondary"
                                     copyable={record?.location ? {text: record.location} : false}>
                        {record?.location || "-"}
                    </Typography.Text>
                </div>
            ),
        },
        {
            title: "相关时间",
            dataIndex: "create_time",
            width: 205,
            sorter: true,
            showSorterTooltip: {title: "按登录时间排序"},
            render: (_: unknown, record: any) => (
                <div className="session-source session-time">
                    <div className="session-time-row">
                        <span className="session-time-label">登录于:</span>
                        <span>{record?.create_time || "-"}</span>
                    </div>
                    <div className="session-time-row">
                        <span className="session-time-label">过期于:</span>
                        <span>{record?.expire_time || "-"}</span>
                    </div>
                </div>
            ),
        },
        ...(hasLogoutControl ? [{
            title: "操作",
            key: "option",
            width: 80,
            fixed: "right" as const,
            render: (_: any, record: any) => (
                <Button size="small" danger onClick={() => logoutSession([record?.token])}>
                    注销
                </Button>
            ),
        }] : []),
    ];

    return (
        <ModalTable
            ref={tableRef}
            className={styles.table}
            title={<Title>登录会话</Title>}
            width={780}
            minHeight={440}
            scrollMode="modal"
            modalProps={{
                mask: {closable: true},
                style: {
                    top: screens.md ? 100 : 24,
                    paddingBottom: screens.md ? 100 : 24,
                    maxWidth: "calc(100vw - 24px)",
                },
                afterClose: () => {
                    setUserInfo(undefined);
                },
            }}
            rowKey="token"
            columns={columns}
            scroll={{x: hasLogoutControl ? 700 : 576}}
            defaultPageSize={7}
            defaultSort={{create_time: "descend"}}
            params={{id: Number(userInfo?.id)}}
            requestEnabled={hasSessionControl && !!userInfo?.id}
            request={(params, sort) => {
                if (!hasSessionControl || !userInfo?.id) {
                    return Promise.resolve(emptyTableData);
                }
                return getData({...params}, sort, getMasterSessionList);
            }}
            search={{
                placeholder: "搜索登录设备",
                ariaLabel: "搜索会话",
                maxLength: 100,
                fields: [
                    {value: "device", label: "设备", placeholder: "搜索登录设备"},
                    {value: "ip", label: "IP", placeholder: "搜索登录 IP", maxLength: 50},
                    {value: "location", label: "地点", placeholder: "搜索登录地点"},
                ],
            }}
            filters={[{type: "date", name: "create_time", label: "登录时间筛选"}]}
            toolbar={{refresh: {ariaLabel: "刷新会话列表"}}}
            rowSelection={hasLogoutControl ? {
                fixed: true,
                columnWidth: 44,
                actions: (keys) => [{
                    key: "logout",
                    type: "button",
                    label: "批量注销",
                    icon: "LogoutOutlined",
                    danger: true,
                    onClick: () => logoutSession(keys),
                }],
            } : false}
            pagination={{unit: "条"}}/>
    );
});

MasterSession.displayName = "MasterSession";

export default MasterSession;
