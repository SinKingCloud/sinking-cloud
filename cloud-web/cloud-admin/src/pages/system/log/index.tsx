import React, {useEffect, useRef, useState} from 'react';
import {Button, DatePicker, Form, message, Modal} from 'antd';
import {PageContainer} from '@ant-design/pro-layout';
import ProTable from '@ant-design/pro-table';
import {getPageControls, getParams} from "@/util/page";
import {ExclamationCircleOutlined} from "@ant-design/icons";
import {deleteLog, getLogList} from "@/services/system/log";

export default (): React.ReactNode => {
  /**
   * 控件权限
   */
  const [controls, setControls] = useState({});
  useEffect(() => {
    getPageControls().then((d) => {
      setControls(d)
    });
  }, [])

  /**
   * 表单处理
   */
  const actionRef = useRef();
  const ref = useRef();

  /**
   * table表格渲染
   */
  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      tip: '记录ID',
      hideInSearch: true,
      hideInTable: true,
      sorter: true,
    },
    {
      title: '请求ID',
      dataIndex: 'request_id',
      tip: '该请求的唯一ID',
      copyable: true,
    },
    {
      title: '用户账号',
      dataIndex: 'account',
      tip: '用户账号',
      copyable: true,
    },
    {
      title: '用户手机',
      dataIndex: 'phone',
      tip: '用户手机',
      copyable: true,
    },
    {
      title: '请求IP',
      dataIndex: 'request_ip',
      tip: '请求者IP',
      copyable: true,
    },
    {
      title: '操作类型',
      dataIndex: 'type',
      tip: '操作类型',
      valueEnum: {
        0: {
          text: '查看',
          status: 'success',
        },
        1: {
          text: '新增',
          status: 'warning',
        },
        2: {
          text: '删除',
          status: 'error',
        },
        3: {
          text: '修改',
          status: 'default',
        },
      },
    },
    {
      title: '标题',
      dataIndex: 'title',
      tip: '标题',
    },
    {
      title: '内容',
      dataIndex: 'content',
      tip: '内容',
    },

    {
      title: '操作时间',
      valueType: 'dateTime',
      dataIndex: 'create_time',
      tip: '操作时间',
      sorter: true,
      hideInSearch: true,
      editable: false,
    },
    {
      title: '操作时间',
      valueType: 'dateTimeRange',
      dataIndex: 'create_time',
      tip: '请求接口时间',
      hideInTable: true,
      search: {
        transform: (value: any) => {
          return {
            create_time_start: value[0],
            create_time_end: value[1],
          };
        },
      },
    },
  ];

  /**
   * 弹窗处理
   */
  const [deleteForm] = Form.useForm();//删除表单
  const [isDeleteModalVisible, setIsDeleteModalVisible] = useState(false);
  const onDeleteFinish = async (values: any) => {
    setIsDeleteModalVisible(false);
    Modal.confirm({
      title: '您确定要删除该时间区间数据吗?',
      icon: <ExclamationCircleOutlined/>,
      content: '删除后该数据不可恢复',
      okType: 'danger',
      onOk() {
        deleteForm.resetFields();
        const rangeTimeValue = values['range-time-picker'];
        deleteLog({
          request_time_start: rangeTimeValue[0].format('YYYY-MM-DD HH:mm:ss'),
          request_time_end: rangeTimeValue[1].format('YYYY-MM-DD HH:mm:ss')
        }).then((r) => {
          if (r.code != 200) {
            message.error(r.message || "请求失败").then()
          } else {
            // @ts-ignore
            actionRef.current.reloadAndRest()
            message.success(r.message).then()
          }
        })
      },
    });
  }

  return (
    <PageContainer title={false}>
      <Modal key={"delete"} width={350} destroyOnClose={true} forceRender={true} title="清理数据"
             visible={isDeleteModalVisible}
             onOk={deleteForm.submit} okText={"确 认"} onCancel={() => {
        setIsDeleteModalVisible(false);
        deleteForm.resetFields();
      }}>
        <Form form={deleteForm} layout="vertical" name="control-hooks" onFinish={onDeleteFinish} labelAlign="right"
              labelCol={{span: 6}}
              wrapperCol={{span: 24}}>
          <Form.Item name="range-time-picker" label="" rules={[{required: true}]}>
            <DatePicker.RangePicker showTime format="YYYY-MM-DD HH:mm:ss"/>
          </Form.Item>
        </Form>
      </Modal>
      <ProTable
        defaultSize={"small"}
        form={{layout: "vertical"}}
        headerTitle={'请求日志'}
        actionRef={actionRef}
        formRef={ref}
        rowKey={'id'}
        options={{
          density: true,
          fullScreen: true,
          setting: true,
        }}
        // @ts-ignore
        columns={columns}
        request={async (params, sort) => {
          const data = await getLogList(getParams(params, sort));
          return {
            data: data.data.list === undefined || data.data.list === null || data.data.list.length <= 0 ? [] : data.data.list,
            success: data.code === 200,
            total: data.data.total,
          };
        }}
        search={{
          defaultCollapsed: true,
          labelWidth: 20,
        }}
        toolBarRender={() => [
          <Button key="primary" type="primary" onClick={() => {
            setIsDeleteModalVisible(true)
          }} hidden={controls["system.log:delete"] == undefined}>
            {(controls["system.log:delete"])?.title || "清理"}
          </Button>,
        ]}
      />
    </PageContainer>
  );
};
