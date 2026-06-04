import React, {forwardRef, useImperativeHandle, useRef, useState} from "react";
import {App, Button, Space, Typography} from "antd";
import {ProModal, ProModalRef, ProTable, ProTableRef, Title} from "sinking-antd";
import {createStyles} from "antd-style";
import {getData} from "@/utils/page";
import {getUserSessionList, logoutUserSession} from "@/service/admin/user";
import {dateRangeTransform, emptyTableData} from "@/utils/page";
import {AdminControlCode} from "@/utils/permission";

export type SessionRef = {
    open: (record: any) => void;
};

type SessionProps = {
    control: any;
};

const useStyles: any = createStyles(({css}: any) => ({
    timeColumn: css`
        min-width: 180px;
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

const Session = forwardRef<SessionRef, SessionProps>(({control}, ref) => {
    const {message, modal} = App.useApp();
    const {styles} = useStyles();
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const tableRef = useRef<ProTableRef>({} as ProTableRef);
    const [userInfo, setUserInfo] = useState<any>();
    const hasSessionControl = control.hasControl(AdminControlCode.WEB_USER_SESSION);
    const hasLogoutControl = control.hasControl(AdminControlCode.WEB_USER_LOGOUT);

    useImperativeHandle(ref, () => ({
        open: (record: any) => {
            if (!hasSessionControl) {
                message?.warning("暂无会话查看权限");
                return;
            }
            setUserInfo(record || {});
            modalRef.current?.show();
            setTimeout(() => tableRef.current?.refreshTableData?.(), 0);
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
            okType: "danger",
            cancelText: "取 消",
            mask: {
                closable: true,
            },
            onOk: async () => {
                await logoutUserSession({
                    body: {
                        id: Number(userInfo?.id),
                        tokens: nextTokens,
                    },
                    onSuccess: (res: any) => {
                        tableRef.current?.refreshTableData?.();
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

    const renderTimeRows = (record: any) => (
        <Space orientation="vertical" size={2} className={styles.timeColumn}>
            <div className={styles.timeRow}>
                <Typography.Text type="secondary" className={styles.labelText}>登录:</Typography.Text>
                <Typography.Text className={styles.smallText}>{record?.create_time || "-"}</Typography.Text>
            </div>
            <div className={styles.timeRow}>
                <Typography.Text type="secondary" className={styles.labelText}>过期:</Typography.Text>
                <Typography.Text className={styles.smallText}>{record?.expire_time || "-"}</Typography.Text>
            </div>
        </Space>
    );

    const columns: any[] = [
        {
            title: "登录设备",
            dataIndex: "device",
            tip: "登录设备",
            valueType: "text",
            copyable: true,
        },
        {
            title: "登录IP",
            dataIndex: "ip",
            tip: "登录IP",
            valueType: "text",
            copyable: true,
        },
        {
            title: "登录地点",
            dataIndex: "location",
            tip: "登录地点",
            valueType: "text",
            copyable: true,
        },
        {
            title: "登录时间",
            dataIndex: "create_time",
            tip: "登录时间",
            valueType: "dateRange",
            sorter: true,
            hideInTable: true,
            transform: dateRangeTransform("create_time"),
        },
        {
            title: "过期时间",
            dataIndex: "expire_time",
            tip: "会话过期时间",
            valueType: "dateRange",
            sorter: true,
            hideInTable: true,
            transform: dateRangeTransform("expire_time"),
        },
        {
            title: "相关时间",
            dataIndex: "time_info",
            tip: "登录时间/过期时间",
            hideInSearch: true,
            render: (_: any, record: any) => renderTimeRows(record),
        },
        ...(hasLogoutControl ? [{
            title: "操作",
            valueType: "option",
            key: "option",
            width: 100,
            hideInSearch: true,
            render: (_: any, record: any) => (
                <Button size="small" danger onClick={() => logoutSession([record?.token])}>
                    注销
                </Button>
            ),
        }] : []),
    ];

    return (
        <ProModal
            ref={modalRef}
            title={<Title>登录会话</Title>}
            width={910}
            modalProps={{
                footer: null,
                mask: {
                    closable: true,
                },
                destroyOnHidden: false,
                afterClose: () => {
                    setUserInfo(undefined);
                    tableRef.current?.clearSelectedRows?.();
                },
            } as any}>
            <ProTable
                ref={tableRef}
                rowKey="token"
                columns={columns}
                defaultPage={1}
                defaultPageSize={10}
                rowSelection={hasLogoutControl ? {
                    fixed: true,
                    rightExtra: [
                        <Button key="logout" type="primary" danger ghost
                                onClick={() => logoutSession(tableRef.current?.getSelectedRowKeys?.() as any)}>
                            批量注销
                        </Button>,
                    ],
                } : false}
                request={(params, sort) => {
                    if (!hasSessionControl || !userInfo?.id) {
                        return Promise.resolve(emptyTableData);
                    }
                    return getData({...params, id: Number(userInfo.id)}, sort, getUserSessionList);
                }}
            />
        </ProModal>
    );
});

Session.displayName = "Session";

export default Session;
