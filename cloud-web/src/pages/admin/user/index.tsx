import React, {useRef, useState} from "react";
import {Body, ProModal, ProModalRef, ProTable, ProTableRef, Title} from "sinking-antd";
import {App, Avatar, Button, Dropdown, Form, Result, Select, Space, Typography} from "antd";
import {UserOutlined} from "@ant-design/icons";
import {createStyles} from "antd-style";
import {getData} from "@/utils/page";
import useControl from "@/utils/control";
import {AdminControlCode} from "@/utils/permission";
import useEnum from "@/utils/enum";
import {getUserList, updateUser} from "@/service/admin/user";
import {
    buildOptions,
    buildValueEnum,
    dateRangeTransform,
    emptyTableData,
    normalizeIds,
} from "@/utils/page";
import Session, {type SessionRef} from "./components/session";
import Info, {type InfoRef} from "./components/info";
import Certificate, {type CertificateRef} from "./components/certificate";

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
    const {message} = App.useApp();
    const {styles} = useStyles();
    const [enumData, enumLoading] = useEnum(["user.user", "user.real_auth"]) as any;
    const control = useControl();
    const hasListControl = control.hasControl(AdminControlCode.WEB_USER_LIST);
    const hasUpdateControl = control.hasControl(AdminControlCode.WEB_USER_UPDATE);
    const hasSessionControl = control.hasControl(AdminControlCode.WEB_USER_SESSION);
    const hasCertificateControl = control.hasControl(AdminControlCode.WEB_USER_CERTIFICATE);
    const tableRef = useRef<ProTableRef>({} as ProTableRef);
    const infoRef = useRef<InfoRef>({} as InfoRef);
    const sessionRef = useRef<SessionRef>({} as SessionRef);
    const certificateRef = useRef<CertificateRef>({} as CertificateRef);
    const batchModalRef = useRef<ProModalRef>({} as ProModalRef);
    const [batchForm] = Form.useForm();
    const [batchIds, setBatchIds] = useState<number[]>([]);
    const [batchSubmitting, setBatchSubmitting] = useState(false);
    const userStatus = enumData?.["user.user"]?.status || {};
    const realAuthStatus = enumData?.["user.real_auth"]?.status || {};
    const realAuthType = enumData?.["user.real_auth"]?.type || {};
    const statusOptions = buildOptions(userStatus, true);
    const realAuthStatusOptions = buildOptions(realAuthStatus, true);
    const realAuthTypeOptions = buildOptions(realAuthType, true);
    const statusValueEnum = buildValueEnum(userStatus, {"0": "green", "1": "red"});
    const realAuthInfoValueEnum = Object.entries(realAuthType || {}).reduce((data: any, [key, value]) => {
        data[key] = {
            text: key === "0" ? (realAuthStatus?.[key] || value) : value,
            color: key === "0" ? "default" : key === "1" ? "blue" : "purple",
        };
        return data;
    }, {});

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

    const openSession = (record: any) => {
        if (!hasSessionControl) {
            message?.warning("暂无会话查看权限");
            return;
        }
        sessionRef.current?.open(record);
    };

    const openEdit = (record: any) => {
        if (!hasUpdateControl) {
            message?.warning("暂无修改权限");
            return;
        }
        infoRef.current?.openEdit(record);
    };

    const openBatchEdit = (ids: any[] = []) => {
        if (!hasUpdateControl) {
            message?.warning("暂无修改权限");
            return;
        }
        const nextIds = normalizeIds(ids);
        if (nextIds.length <= 0) {
            message?.warning("请选择要修改的用户");
            return;
        }
        setBatchIds(nextIds);
        batchForm.resetFields();
        batchModalRef.current?.show();
    };

    const openCertificate = (record: any) => {
        if (!hasCertificateControl) {
            message?.warning("暂无实名查看权限");
            return;
        }
        certificateRef.current?.open(record);
    };

    const refreshTable = () => {
        tableRef.current?.refreshTableData?.();
    };

    const handleBatchSubmit = async (values: any) => {
        if (!hasUpdateControl) {
            message?.warning("暂无修改权限");
            return;
        }
        if (batchIds.length <= 0) {
            message?.warning("请选择要修改的用户");
            return;
        }
        setBatchSubmitting(true);
        await updateUser({
            body: {
                ids: batchIds,
                status: values?.status?.toString(),
            },
            onSuccess: (res: any) => {
                batchModalRef.current?.hide();
                tableRef.current?.clearSelectedRows?.();
                refreshTable();
                message?.success(res?.message || "操作成功");
            },
            onFail: (res: any) => {
                message?.error(res?.message || "操作失败");
            },
            onFinally: () => setBatchSubmitting(false),
        });
    };

    const renderRealAuthInfo = (record: any) => {
        const realAuthStatusValue = Number(record?.real_auth_status);
        if (realAuthStatusValue <= 0) {
            return realAuthStatus?.[record?.real_auth_status] || realAuthType?.["0"] || "未实名";
        }
        return realAuthType?.[record?.real_auth_type] || realAuthStatus?.[record?.real_auth_status] || "-";
    };

    const columns: any[] = [
        {
            title: "ID",
            dataIndex: "id",
            tip: "用户ID",
            sorter: true,
            hideInSearch: true,
        },
        {
            title: "头像",
            dataIndex: "avatar",
            tip: "用户头像",
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
            title: "用户账号",
            dataIndex: "account",
            tip: "用户账号",
            valueType: "text",
            hideInTable: true,
        },
        {
            title: "用户昵称",
            dataIndex: "nick_name",
            tip: "用户昵称",
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
            title: "实名信息",
            dataIndex: "real_auth_info",
            tip: "用户实名信息",
            valueEnum: realAuthInfoValueEnum,
            render: (_: any, record: any) => renderRealAuthInfo(record),
            transform: (value: any) => {
                if (value === undefined || value === null || value === "") {
                    return {};
                }
                const realAuthTypeValue = Number(value);
                if (realAuthTypeValue > 0) {
                    return {
                        real_auth_status: "1",
                        real_auth_type: String(realAuthTypeValue),
                    };
                }
                return {
                    real_auth_status: "0",
                };
            },
        },
        {
            title: "用户状态",
            dataIndex: "status",
            tip: "用户状态",
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
        ...(hasUpdateControl || hasCertificateControl || hasSessionControl ? [{
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
                        onClick: () => openEdit(record),
                    } : null,
                    hasCertificateControl ? {
                        key: "certificate",
                        label: "实名",
                        onClick: () => openCertificate(record),
                    } : null,
                    hasSessionControl ? {
                        key: "session",
                        label: "会话",
                        onClick: () => openSession(record),
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
                    title={<Title>用户列表</Title>}
                    rowKey="id"
                    columns={columns}
                    defaultPage={1}
                    defaultPageSize={20}
                    request={(params, sort) => {
                        if (!hasListControl) {
                            return Promise.resolve(emptyTableData);
                        }
                        return getData(params, sort, getUserList);
                    }}
                    rowSelection={hasUpdateControl ? {
                        fixed: true,
                        rightExtra: (
                            <Button key="edit" type="primary" ghost
                                    onClick={() => openBatchEdit(tableRef.current?.getSelectedRowKeys?.() as any)}>
                                批量编辑
                            </Button>
                        ),
                    } : false}
                    paginationAffix={true}
                    selectionAffix={true}
                />
                <Info
                    ref={infoRef}
                    control={control}
                    statusOptions={statusOptions}
                    realAuthStatusOptions={realAuthStatusOptions}
                    realAuthTypeOptions={realAuthTypeOptions}
                    onSuccess={refreshTable}
                />
                <Session ref={sessionRef} control={control}/>
                <Certificate ref={certificateRef} control={control}/>
                <ProModal
                    ref={batchModalRef}
                    title={<Title>批量编辑</Title>}
                    onOk={batchForm.submit}
                    width={320}
                    modalProps={{
                        confirmLoading: batchSubmitting,
                        okText: "提交",
                        cancelText: "取消",
                        forceRender: true,
                        mask: {
                            closable: true,
                        },
                        afterClose: () => {
                            batchForm.resetFields();
                            setBatchIds([]);
                            setBatchSubmitting(false);
                        },
                    } as any}>
                    <Form form={batchForm} layout="vertical" onFinish={handleBatchSubmit}>
                        <Form.Item name="status" label="用户状态" rules={[{required: true, message: "请选择用户状态"}]}>
                            <Select placeholder="请选择用户状态" options={statusOptions as any}/>
                        </Form.Item>
                    </Form>
                </ProModal>
            </>}
        </Body>
    );
};
