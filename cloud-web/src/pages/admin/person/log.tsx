import React from "react";
import {Space, Tag, Typography} from "antd";
import {createStyles} from "antd-style";
import {Body, ProTable, Title} from "sinking-antd";
import {buildValueEnum, dateRangeTransform, getData} from "@/utils/page";
import useEnum from "@/utils/enum";
import {getLogList} from "@/service/admin/person";

const useStyles = createStyles(({css, token}) => ({
    typeTag: css`
        && {
            margin-inline-end: 0;
            color: ${token.colorPrimaryText};
            background: ${token.colorPrimaryBg};
        }
    `,
}));

export default (): React.ReactNode => {
    const {styles} = useStyles();
    const [enumData, enumLoading] = useEnum("admin.user_log");
    const typeData = enumData?.type || {};

    const columns: any[] = [
        {
            title: "日志ID / 请求ID",
            dataIndex: "id",
            width: 300,
            sorter: true,
            hideInSearch: true,
            render: (value: number, record: any) => (
                <Space orientation="vertical" size={2}>
                    <Typography.Text copyable={value ? {text: String(value)} : false}>
                        {value || "-"}
                    </Typography.Text>
                    <Typography.Text type="secondary"
                                     copyable={record?.request_id ? {text: record.request_id} : false}>
                        {record?.request_id || "-"}
                    </Typography.Text>
                </Space>
            ),
        },
        {
            title: "操作类型",
            dataIndex: "type",
            sorter: true,
            valueEnum: buildValueEnum(typeData),
            render: (value: number) => (
                <Tag className={styles.typeTag} variant="filled">
                    {typeData[value] || "-"}
                </Tag>
            ),
        },
        {
            title: "操作信息",
            dataIndex: "info",
            hideInSearch: true,
            render: (_: any, record: any) => (
                <Space orientation="vertical" size={2}>
                    <Typography.Text>{record?.title || "-"}</Typography.Text>
                    <Typography.Text type="secondary">{record?.content || "-"}</Typography.Text>
                </Space>
            ),
        },
        {
            title: "请求来源",
            dataIndex: "source",
            hideInSearch: true,
            render: (_: any, record: any) => (
                <Space orientation="vertical" size={2}>
                    <Typography.Text copyable={record?.request_ip ? {text: record.request_ip} : false}>
                        {record?.request_ip || "-"}
                    </Typography.Text>
                    <Typography.Text type="secondary">
                        {record?.location || "-"}
                    </Typography.Text>
                </Space>
            ),
        },
        {
            title: "操作时间",
            dataIndex: "create_time",
            valueType: "dateRange",
            sorter: true,
            transform: dateRangeTransform("create_time"),
            render: (_: any, record: any) => record?.create_time || "-",
        },
        {
            title: "请求ID",
            dataIndex: "request_id",
            hideInTable: true,
        },
        {
            title: "请求IP",
            dataIndex: "request_ip",
            hideInTable: true,
            props: {maxLength: 50},
        },
        {
            title: "请求地点",
            dataIndex: "location",
            hideInTable: true,
            props: {maxLength: 100},
        },
        {
            title: "操作标题",
            dataIndex: "title",
            hideInTable: true,
        },
        {
            title: "操作内容",
            dataIndex: "content",
            hideInTable: true,
        },
    ];

    return (
        <Body loading={enumLoading}>
            <ProTable
                rowKey="id"
                extraRefreshBtn={true}
                title={<Title>操作日志</Title>}
                columns={columns}
                defaultPage={1}
                defaultPageSize={10}
                request={(params, sort) => getData(params, sort, getLogList)}
                paginationAffix={true}/>
        </Body>
    );
}
