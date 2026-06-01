import React, {useRef, useState} from "react";
import {App, Button, Space, Typography} from "antd";
import {createStyles} from "antd-style";
import {Body, ProTable, ProTableRef, Title} from "sinking-antd";
import {dateRangeTransform, getData} from "@/utils/page";
import {getSessionList, logoutSession} from "@/service/admin/person";

const useStyles: any = createStyles(({css}: any) => ({
    timeColumn: css`
        min-width: 190px;
    `,
    timeRow: css`
        display: flex;
        align-items: center;
        gap: 4px;
        white-space: nowrap;
    `,
    labelText: css`
        flex: none;
        min-width: 40px;
        font-size: 12px;
    `,
    smallText: css`
        flex: 1;
        min-width: 0;
        font-size: 12px;
    `,
}));

export default (): React.ReactNode => {
    const {message, modal} = App.useApp();
    const {styles} = useStyles();
    const tableRef = useRef<ProTableRef>({} as ProTableRef);
    const [selectedTokens, setSelectedTokens] = useState<string[]>([]);
    const [logoutLoading, setLogoutLoading] = useState(false);

    const handleLogout = (tokens: any[] = []) => {
        const nextTokens = Array.from(new Set((tokens || []).map((item) => String(item || "")).filter(Boolean)));
        if (nextTokens.length <= 0) {
            message?.warning("请选择要注销的会话");
            return;
        }
        modal.confirm({
            title: "注销会话",
            content: `确定注销选中的 ${nextTokens.length} 个会话吗？`,
            okText: "确 定",
            okType: "danger",
            cancelText: "取 消",
            mask: {
                closable: true,
            },
            onOk: async () => {
                setLogoutLoading(true);
                await logoutSession({
                    body: {
                        tokens: nextTokens,
                    },
                    onSuccess: async (res: any) => {
                        message?.success(res?.message || "注销成功");
                        tableRef.current?.clearSelectedRows?.();
                        setSelectedTokens([]);
                        tableRef.current?.refreshTableData?.();
                    },
                    onFail: (res: any) => {
                        message?.error(res?.message || "注销失败");
                    },
                    onFinally: () => {
                        setLogoutLoading(false);
                    },
                });
            },
        } as any);
    };

    const columns: any[] = [
        {
            title: "登录设备",
            dataIndex: "device",
            valueType: "text",
            copyable: true,
        },
        {
            title: "登录IP",
            dataIndex: "ip",
            valueType: "text",
            copyable: true,
        },
        {
            title: "登录地点",
            dataIndex: "location",
            valueType: "text",
            copyable: true,
        },
        {
            title: "登录时间",
            dataIndex: "create_time",
            valueType: "dateRange",
            sorter: true,
            hideInTable: true,
            transform: dateRangeTransform("create_time"),
        },
        {
            title: "过期时间",
            dataIndex: "expire_time",
            valueType: "dateRange",
            sorter: true,
            hideInTable: true,
            transform: dateRangeTransform("expire_time"),
        },
        {
            title: "相关时间",
            dataIndex: "time_info",
            hideInSearch: true,
            render: (_: any, record: any) => {
                return <Space orientation="vertical" size={2} className={styles.timeColumn}>
                    <div className={styles.timeRow}>
                        <Typography.Text type="secondary" className={styles.labelText}>登录:</Typography.Text>
                        <Typography.Text className={styles.smallText}>{record?.create_time || "-"}</Typography.Text>
                    </div>
                    <div className={styles.timeRow}>
                        <Typography.Text type="secondary" className={styles.labelText}>过期:</Typography.Text>
                        <Typography.Text className={styles.smallText}>{record?.expire_time || "-"}</Typography.Text>
                    </div>
                </Space>
            },
        },
        {
            title: "操作",
            valueType: "option",
            key: "option",
            width: 110,
            hideInSearch: true,
            render: (_: any, record: any) => (
                <Button size="small" danger onClick={() => handleLogout([record?.token])}>
                    注销
                </Button>
            ),
        },
    ];

    return (
        <Body>
            <ProTable
                ref={tableRef}
                rowKey="token"
                extraRefreshBtn={true}
                title={<Title>会话管理</Title>}
                columns={columns}
                defaultPage={1}
                defaultPageSize={10}
                rowSelection={{
                    fixed: true,
                    onChange: (keys: React.Key[]) => {
                        setSelectedTokens((keys || []).map((item) => String(item || "")).filter(Boolean));
                    },
                    rightExtra: [
                        <Button
                            key="logout"
                            type="primary"
                            danger
                            ghost
                            loading={logoutLoading}
                            onClick={() => handleLogout(selectedTokens)}>
                            批量注销
                        </Button>,
                    ],
                }}
                request={async (params: any, sort: any) => {
                    return getData(params, sort, getSessionList);
                }}
                paginationAffix={true}
                selectionAffix={true}
            />
        </Body>
    );
}
