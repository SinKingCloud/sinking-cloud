import React, {useRef} from "react";
import {Body, ProTable, ProTableRef, Title} from "sinking-antd";
import {App, Avatar, Button, Dropdown, Result, Space, Typography} from "antd";
import {UserOutlined} from "@ant-design/icons";
import {getData} from "@/utils/page";
import useControl from "@/utils/control";
import useEnum from "@/utils/enum";
import {
    deleteMaster,
    getMasterList,
} from "@/service/admin/master";
import Info, {type InfoRef} from "./components/info";
import MasterSession, {type MasterSessionRef} from "./components/session";
import Role, {type RoleRef} from "./components/role";
import {
    buildOptions,
    buildValueEnum,
    dateRangeTransform,
    emptyTableData,
    normalizeIds,
} from "@/utils/page";
import {createStyles} from "antd-style";

const useStyles: any = createStyles((): any => ({
    infoColumn: {
        width: "100%",
        minWidth: 180,
    },
    infoRow: {
        display: "flex",
        alignItems: "center",
        gap: 4,
        minWidth: 0,
    },
    labelText: {
        flex: "none",
        minWidth: 40,
        fontSize: 12,
    },
    smallText: {
        flex: 1,
        minWidth: 0,
        fontSize: 12,
        whiteSpace: "nowrap",
    },
}));

export default (): React.ReactNode => {
    const {message, modal} = App.useApp();
    const {styles} = useStyles();
    const [enumData, enumLoading] = useEnum(["admin.user", "admin.role"]) as any;
    const control = useControl();
    const hasListControl = control.hasControl("admin.system.master:list");
    const hasCreateControl = control.hasControl("admin.system.master:create");
    const hasUpdateControl = control.hasControl("admin.system.master:update");
    const hasDeleteControl = control.hasControl("admin.system.master:delete");
    const hasRoleControl = control.hasControl("admin.system.master:role");
    const hasSessionControl = control.hasControl("admin.system.master:session");
    const tableRef = useRef<ProTableRef>({} as ProTableRef);
    const infoRef = useRef<InfoRef>({} as InfoRef);
    const roleRef = useRef<RoleRef>({} as RoleRef);
    const sessionRef = useRef<MasterSessionRef>({} as MasterSessionRef);
    const userStatus = enumData?.["admin.user"]?.status || {};
    const roleNames = enumData?.["admin.role"]?.names || {};
    const statusOptions = buildOptions(userStatus, true);
    const statusValueEnum = buildValueEnum(userStatus, {"0": "green", "1": "red"});

    const onDelete = (id: any) => {
        if (!hasDeleteControl) {
            message?.warning("暂无删除权限");
            return;
        }
        const nextId = Number(id);
        if (!nextId) {
            message?.warning("请选择要删除的员工");
            return;
        }
        modal.confirm({
            title: "删除员工",
            content: "确定删除该员工吗？",
            okText: "确 定",
            okType: "danger",
            cancelText: "取 消",
            mask: {
                closable: true,
            },
            onOk: async () => {
                await deleteMaster({
                    body: {id: nextId},
                    onSuccess: (res: any) => {
                        tableRef.current?.refreshTableData?.();
                        tableRef.current?.clearSelectedRows?.();
                        message?.success(res?.message || "删除成功");
                    },
                    onFail: (res: any) => {
                        message?.error(res?.message || "删除失败");
                    },
                });
            },
        } as any);
    };

    const roleRender = (ids: any[] = []) => {
        const names = normalizeIds(ids).map((id) => roleNames?.[id]).filter(Boolean);
        if (names.length <= 0) {
            return "-";
        }
        return (
            <Space size={4} wrap>
                {names.map((name) => (
                    <Typography.Text key={name}>{name}</Typography.Text>
                ))}
            </Space>
        );
    };

    const renderRows = (items: { label: string; value: any }[], copyable = true) => (
        <Space orientation="vertical" size={2} className={styles.infoColumn}>
            {items.map((item) => (
                <div className={styles.infoRow} key={item.label}>
                    <Typography.Text type="secondary" className={styles.labelText}>{item.label}</Typography.Text>
                    <Typography.Text
                        className={styles.smallText}
                        copyable={copyable && item.value ? {text: item.value} : false}>
                        {item.value || "-"}
                    </Typography.Text>
                </div>
            ))}
        </Space>
    );

    const columns: any[] = [
        {
            title: "ID",
            dataIndex: "id",
            tip: "员工ID",
            sorter: true,
            hideInSearch: true,
        },
        {
            title: "头像",
            dataIndex: "avatar",
            tip: "员工头像",
            hideInSearch: true,
            render: (text: any) => <Avatar src={text || undefined} icon={<UserOutlined/>}/>,
        },
        {
            title: "账号信息",
            dataIndex: "account_info",
            tip: "账号/昵称",
            hideInSearch: true,
            render: (_: any, record: any) => renderRows([
                {label: "账号:", value: record?.account},
                {label: "昵称:", value: record?.nick_name},
            ]),
        },
        {
            title: "员工账号",
            dataIndex: "account",
            tip: "员工账号",
            valueType: "text",
            hideInTable: true,
        },
        {
            title: "员工昵称",
            dataIndex: "nick_name",
            tip: "员工昵称",
            valueType: "text",
            hideInTable: true,
        },
        {
            title: "联系方式",
            dataIndex: "contact_info",
            tip: "手机号/邮箱",
            hideInSearch: true,
            render: (_: any, record: any) => renderRows([
                {label: "手机:", value: record?.phone},
                {label: "邮箱:", value: record?.email},
            ]),
        },
        {
            title: "手机号码",
            dataIndex: "phone",
            tip: "手机号码",
            valueType: "text",
            hideInTable: true,
        },
        {
            title: "邮箱地址",
            dataIndex: "email",
            tip: "邮箱地址",
            valueType: "text",
            hideInTable: true,
        },
        {
            title: "角色",
            dataIndex: "role_ids",
            tip: "员工拥有的角色",
            hideInSearch: true,
            render: (_: any, record: any) => roleRender(record?.role_ids || []),
        },
        {
            title: "员工状态",
            dataIndex: "status",
            tip: "员工状态",
            valueEnum: statusValueEnum,
        },
        {
            title: "创建时间",
            valueType: "dateRange",
            dataIndex: "create_time",
            tip: "创建时间",
            hideInTable: true,
            sorter: true,
            transform: dateRangeTransform("create_time"),
        },
        {
            title: "更新时间",
            valueType: "dateRange",
            dataIndex: "update_time",
            tip: "更新时间",
            hideInTable: true,
            sorter: true,
            transform: dateRangeTransform("update_time"),
        },
        {
            title: "相关时间",
            dataIndex: "time_info",
            tip: "创建时间/更新时间",
            hideInSearch: true,
            render: (_: any, record: any) => renderRows([
                {label: "创建于:", value: record?.create_time},
                {label: "更新于:", value: record?.update_time},
            ], false),
        },
        ...(hasUpdateControl || hasDeleteControl || hasRoleControl || hasSessionControl ? [{
            title: "操作",
            valueType: "option",
            key: "option",
            fixed: "right",
            hideInSearch: true,
            render: (_: any, record: any) => {
                const items = [
                    hasUpdateControl ? {
                        key: "edit",
                        label: "编辑",
                        onClick: () => infoRef.current?.openEdit(record),
                    } : null,
                    hasRoleControl ? {
                        key: "role",
                        label: "角色",
                        onClick: () => roleRef.current?.open(record),
                    } : null,
                    hasSessionControl ? {
                        key: "session",
                        label: "会话",
                        onClick: () => sessionRef.current?.open(record),
                    } : null,
                    hasDeleteControl ? {
                        key: "delete",
                        label: "删除",
                        danger: true,
                        onClick: () => onDelete(record?.id),
                    } : null,
                ].filter(Boolean) as any;
                if (items.length <= 0) {
                    return "-";
                }
                return (
                    <Dropdown menu={{items}} trigger={["click"]} placement="bottom" arrow>
                        <Button size="small">操作</Button>
                    </Dropdown>
                );
            },
        }] : []),
    ];

    return (
        <Body loading={control.loading || enumLoading}>
            {!hasListControl && <Result status="403" title="403" subTitle="您无权限访问此页面,请联系管理员"/>}
            {hasListControl && <>
                <ProTable
                    ref={tableRef}
                    extraRefreshBtn={true}
                    title={<Title>员工列表</Title>}
                    rowKey="id"
                    columns={columns}
                    defaultPage={1}
                    defaultPageSize={20}
                    extra={hasCreateControl ? (
                        <Button type="primary" onClick={() => infoRef.current?.openCreate()}>
                            新增
                        </Button>
                    ) : null}
                    request={(params, sort) => {
                        if (!hasListControl) {
                            return Promise.resolve(emptyTableData);
                        }
                        return getData(params, sort, getMasterList);
                    }}
                    paginationAffix={true}
                    selectionAffix={true}
                />

                <Info
                    ref={infoRef}
                    control={control}
                    statusOptions={statusOptions}
                    onSuccess={() => {
                        tableRef.current?.refreshTableData?.();
                        tableRef.current?.clearSelectedRows?.();
                    }}
                />
                <Role ref={roleRef} control={control} onSuccess={() => tableRef.current?.refreshTableData?.()}/>
                <MasterSession ref={sessionRef} control={control}/>
            </>}
        </Body>
    );
};
