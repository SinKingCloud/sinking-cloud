import React, {useEffect, useRef, useState} from 'react';
import {Button, DatePicker, Descriptions, Form, message, Modal, Spin, Typography} from 'antd';
import {PageContainer} from '@ant-design/pro-layout';
import ProTable from '@ant-design/pro-table';
import {getPageControls, getParams} from "@/util/page";
import {deleteLogRequest, getLogRequestInfo, getLogRequestList} from "@/services/system/log_request";
import {ExclamationCircleOutlined} from "@ant-design/icons";

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
      editable: false,
    },
    {
      title: '请求ID',
      dataIndex: 'request_id',
      tip: '该请求的唯一ID',
      sorter: true,
    },
    {
      title: '请求IP',
      dataIndex: 'request_ip',
      tip: '请求者IP',
      copyable: true,
    },
    {
      title: '请求地址',
      dataIndex: 'request_url',
      tip: '请求地址',
      copyable: true,
    },
    {
      title: '请求方式',
      dataIndex: 'request_method',
      tip: '请求方式',
      valueEnum: {
        GET: {
          text: 'GET',
          status: 'success',
        },
        POST: {
          text: 'POST',
          status: 'error',
        },
        UPDATE: {
          text: 'UPDATE',
          status: 'warning',
        },
        PUT: {
          text: 'PUT',
          status: 'default',
        },
      },
    },
    {
      title: '请求时间',
      valueType: 'dateTime',
      dataIndex: 'request_time',
      tip: '发起请求时间',
      sorter: true,
      hideInSearch: true,
      editable: false,
    },
    {
      title: '响应时间',
      valueType: 'dateTime',
      dataIndex: 'response_time',
      tip: '请求响应时间',
      hideInSearch: true,
      editable: false,
      render: (text: string, record: any) => {
        return record?.response_time + "ms";
      }
    },
    {
      title: '响应码',
      dataIndex: 'response_code',
      tip: 'HTTP响应码',
    },
    {
      title: '服务组',
      dataIndex: 'group_name',
      tip: '微服务服务组名称',
    },
    {
      title: '服务名',
      dataIndex: 'app_name',
      tip: '微服务服务名称',
    },
    {
      title: '服务地址',
      dataIndex: 'app_addr',
      tip: '微服务节点服务地址',
    },
    {
      title: '请求时间',
      valueType: 'dateTimeRange',
      dataIndex: 'update_time',
      tip: '请求接口时间',
      hideInTable: true,
      search: {
        transform: (value: any) => {
          return {
            request_time_start: value[0],
            request_time_end: value[1],
          };
        },
      },
    },
    {
      title: '操作',
      dataIndex: 'Option',
      hideInSearch: true,
      editable: false,
      render: (text: string, record: any) => {
        return <Button type="default" style={{fontSize: "10px"}} size="small" onClick={() => {
          setIsInfoModalVisible(true)
          setIsInfoLoadingVisible(true)
          setLogRequestInfo({});
          getLogRequestInfo({id: record?.id}).then((r) => {
            if (r.code != 200) {
              setIsInfoModalVisible(false)
              message.error(r.message || "请求失败").then()
            } else {
              setLogRequestInfo(r.data);
            }
            setIsInfoLoadingVisible(false)
          })
        }
        }>
          详细报文
        </Button>;
      }
    },
  ];

  /**
   * 弹窗处理
   */
  const [logRequestInfo, setLogRequestInfo] = useState({});
  const [isInfoModalVisible, setIsInfoModalVisible] = useState(false);
  const [isInfoLoadingVisible, setIsInfoLoadingVisible] = useState(false);

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
        deleteLogRequest({
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
        <Form form={deleteForm} name="control-hooks" onFinish={onDeleteFinish} labelAlign="right"
              labelCol={{span: 6}}
              wrapperCol={{span: 24}}>
          <Form.Item name="range-time-picker" label="" rules={[{required: true}]}>
            <DatePicker.RangePicker showTime format="YYYY-MM-DD HH:mm:ss"/>
          </Form.Item>
        </Form>
      </Modal>
      <Modal key={"info"} destroyOnClose={true} forceRender={true} title="详细报文" visible={isInfoModalVisible}
             footer={null}
             onCancel={() => {
               setIsInfoModalVisible(false)
             }}>
        <Spin spinning={isInfoLoadingVisible}>
          <div hidden={isInfoLoadingVisible}>
            <Descriptions column={1} size={"small"}>
              <Descriptions.Item label="请求ID">
                <Typography.Text copyable>{
                  // @ts-ignore
                  logRequestInfo?.request_id || "加载中..."
                }</Typography.Text>
              </Descriptions.Item>
              <Descriptions.Item label="服务组">
                <Typography.Text copyable>{
                  // @ts-ignore
                  logRequestInfo?.group_name || "加载中..."
                }</Typography.Text>
              </Descriptions.Item>
              <Descriptions.Item label="服务名">
                <Typography.Text copyable>{
                  // @ts-ignore
                  logRequestInfo?.app_name || "加载中..."
                }</Typography.Text>
              </Descriptions.Item>
              <Descriptions.Item label="服务地址">
                <Typography.Text copyable>{
                  // @ts-ignore
                  logRequestInfo?.app_addr || "加载中..."
                }</Typography.Text>
              </Descriptions.Item>
              <Descriptions.Item label="请求IP">
                <Typography.Text copyable>{
                  // @ts-ignore
                  logRequestInfo?.request_ip || "加载中..."
                }</Typography.Text>
              </Descriptions.Item>
              <Descriptions.Item label="请求地址">
                <Typography.Text copyable>{
                  // @ts-ignore
                  logRequestInfo?.request_url || "加载中..."
                }</Typography.Text>
              </Descriptions.Item>
              <Descriptions.Item label="请求类型">
                <Typography.Text copyable>{
                  // @ts-ignore
                  logRequestInfo?.request_method || "加载中..."
                }</Typography.Text>
              </Descriptions.Item>
              <Descriptions.Item label="请求时间">
                <Typography.Text>{
                  // @ts-ignore
                  logRequestInfo?.request_time || "加载中..."
                }</Typography.Text>
              </Descriptions.Item>
              <Descriptions.Item label="响应码">
                <Typography.Text>{
                  // @ts-ignore
                  logRequestInfo?.response_code || "加载中..."
                }</Typography.Text>
              </Descriptions.Item>
              <Descriptions.Item label="响应时间">
                <Typography.Text>{
                  // @ts-ignore
                  logRequestInfo?.response_time + "ms" || "加载中..."
                }</Typography.Text>
              </Descriptions.Item>
            </Descriptions>

            <Descriptions column={1} size={"small"}>
              <Descriptions.Item>
                <Typography.Paragraph>
                  请求体：
                  <pre>
                  {
                    // @ts-ignore
                    logRequestInfo?.request_header
                  }
                    <br/><br/>
                    {
                      // @ts-ignore
                      logRequestInfo?.request_body
                    }
                </pre>
                </Typography.Paragraph>
              </Descriptions.Item>
              <Descriptions.Item>
                <Typography.Paragraph>
                  响应体：
                  <pre>
                  {
                    // @ts-ignore
                    logRequestInfo?.response_header
                  }
                    <br/><br/>
                    {
                      // @ts-ignore
                      logRequestInfo?.response_body
                    }
                </pre>
                </Typography.Paragraph>
              </Descriptions.Item>
            </Descriptions>
          </div>
        </Spin>
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
          const data = await getLogRequestList(getParams(params, sort));
          if (data.code != 200) {
            message.error(data.message);
          }
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
          }} hidden={controls["system.log_request:delete"] == undefined}>
            {(controls["system.log_request:delete"])?.title || "清理"}
          </Button>,
        ]}
      />
    </PageContainer>
  );
};
