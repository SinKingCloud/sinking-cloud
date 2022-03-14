import React, {useEffect, useRef, useState} from 'react';
import {Button, Dropdown, Form, Input, Menu, message, Modal, Select, Table} from 'antd';
import {FooterToolbar, PageContainer} from '@ant-design/pro-layout';
import ProTable from '@ant-design/pro-table';
import {createAccess, deleteAccess, getAccessList, updateAccess} from "@/services/system/access";
import {getPageControls, getParams} from "@/util/page";
import {DownOutlined, ExclamationCircleOutlined} from "@ant-design/icons";

export default (): React.ReactNode => {
  /**
   * 控件权限
   */
  const [controls, setControls] = useState({});
  useEffect(() => {
    getPageControls().then((d) => {
      setControls(d)
    })
  }, [])

  /**
   * 表单处理
   */
  const actionRef = useRef();
  const ref = useRef();
  const [selectedRowsState, setSelectedRows] = useState([]);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [isModalEditVisible, setIsModalEditVisible] = useState(false);

  /**
   * 获取选中项
   */
  const getSelectedIds = () => {
    let ids: any[] = [];
    selectedRowsState.map((k: any) => {
      ids.push(k.id);
    });
    return ids;
  };

  /**
   * 表单
   */
  const [form] = Form.useForm();//新建编辑表单
  const [edit] = Form.useForm();//批量修改表单
  /**
   * 新建编辑提交表单
   * @param values 表单项
   */
  const onFormFinish = async (values: any) => {
    let api = createAccess;
    if (values.id != undefined) {
      api = updateAccess;
      values.ids = [values.id];
      values.status = values.status.toString();
      delete values.id;
    }
    setIsModalVisible(false)
    api(values).then((r) => {
      if (r.code != 200) {
        message.error(r.message || "请求失败").then()
      } else {
        // @ts-ignore
        actionRef.current.reload()
        form.resetFields();
        message.success(r.message).then()
      }
    })
  }

  /**
   * 批量修改提交表单
   * @param values 表单项
   */
  const onEditFinish = async (values: any) => {
    setIsModalEditVisible(false);
    values.ids = getSelectedIds()
    if (values.status != undefined && values.status >= 0) {
      values.status = values.status.toString();
    }
    updateAccess(values).then((r) => {
      if (r.code != 200) {
        message.error(r.message || "请求失败").then()
      } else {
        // @ts-ignore
        actionRef.current.reload()
        // @ts-ignore
        actionRef.current.clearSelected();
        edit.resetFields();
        message.success(r.message).then()
      }
    })
  }

  /**
   * 删除数据
   * @param ids id列表
   */
  const deleteSubmit = (ids: any[]) => {
    Modal.confirm({
      title: '您确定要删除该数据吗?',
      icon: <ExclamationCircleOutlined/>,
      content: '删除后该数据不可恢复',
      okType: 'danger',
      onOk() {
        deleteAccess({ids: ids}).then((r) => {
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

  /**
   * 批量删除
   */
  const modalDelete = () => {
    deleteSubmit(getSelectedIds())
  }

  /**
   * table表格渲染
   */
  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      tip: '记录ID',
      hideInSearch: true,
      sorter: true,
    },
    {
      title: '名称',
      dataIndex: 'name',
      tip: '接口名称',
    },
    {
      title: '描述',
      dataIndex: 'describe',
      tip: '接口描述',
    },
    {
      title: '地址',
      dataIndex: 'url',
      tip: '接口地址',
      copyable: true,
    },
    {
      title: '方式',
      dataIndex: 'method',
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
      title: '修改时间',
      valueType: 'dateTime',
      dataIndex: 'update_time',
      tip: '上次修改时间',
      sorter: true,
      hideInSearch: true,
    },
    {
      title: '添加时间',
      valueType: 'dateTime',
      dataIndex: 'create_time',
      tip: '上次修改时间',
      sorter: true,
      hideInSearch: true,
    },
    {
      title: '修改时间',
      valueType: 'dateRange',
      dataIndex: 'update_time',
      tip: '上次修改时间',
      hideInTable: true,
      search: {
        transform: (value: any) => {
          return {
            update_time_start: value[0],
            update_time_end: value[1],
          };
        },
      },
    },
    {
      title: '添加时间',
      valueType: 'dateRange',
      dataIndex: 'create_time',
      tip: '添加记录时间',
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
    {
      title: '状态',
      dataIndex: 'status',
      tip: '接口状态',
      sorter: true,
      valueEnum: {
        0: {
          text: '允许访问',
          status: 'success',
        },
        1: {
          text: '禁止访问',
          status: 'error',
        },
      },
    },
    {
      title: '操作',
      dataIndex: 'Option',
      hideInSearch: true,
      editable: false,
      render: (text: string, record: any) => {
        return <Dropdown overlay={(
          <Menu>
            <Menu.Item key={"edit"} hidden={controls["system.access:edit"] == undefined}>
              <a style={{fontSize: "small"}} onClick={() => {
                form.setFieldsValue(record)
                setIsModalVisible(true);
              }}>
                编 辑
              </a>
            </Menu.Item>
            <Menu.Item key={"delete"} hidden={controls["system.access:delete"] == undefined}>
              <a style={{fontSize: "small"}} onClick={() => {
                // @ts-ignore
                deleteSubmit([record.id])
              }}>
                删 除
              </a>
            </Menu.Item>
          </Menu>
        )} trigger={['click']} placement="bottomCenter" arrow={true}>
          <Button size={"small"} style={{fontSize: "small"}} onClick={e => e.preventDefault()}>操
            作 <DownOutlined/></Button>
        </Dropdown>;
      }
    },
  ];

  return (
    <PageContainer title={false}>
      <Modal key={"form"} destroyOnClose={true} forceRender={true}
             title={form.getFieldValue("id") == undefined ? "新 建" : "编 辑"}
             visible={isModalVisible} onOk={form.submit} okText={"确 认"} onCancel={() => {
        setIsModalVisible(false);
        form.resetFields();
      }}>
        <Form form={form} name="control-hooks" onFinish={onFormFinish} labelAlign="right" labelCol={{span: 4}}
              wrapperCol={{span: 18}}>
          <Form.Item name={"id"} label="ID" hidden={true}>
            <Input placeholder="请输入ID"/>
          </Form.Item>
          <Form.Item name={"name"} label="名称" rules={[{required: true}]}>
            <Input placeholder="请输入名称"/>
          </Form.Item>
          <Form.Item name={"describe"} label="描述" rules={[{required: true}]}>
            <Input.TextArea placeholder="请输入描述"/>
          </Form.Item>
          <Form.Item name={"url"} label="请求地址" rules={[{required: true}]}>
            <Input placeholder="请输入地址"/>
          </Form.Item>
          <Form.Item name={"method"} label="请求方式" rules={[{required: true}]}>
            <Select placeholder="请输入请求方式">
              <Select.Option value="GET">GET</Select.Option>
              <Select.Option value="POST">POST</Select.Option>
              <Select.Option value="UPDATE">UPDATE</Select.Option>
              <Select.Option value="PUT">PUT</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name={"status"} label="状态" rules={[{required: true}]}>
            <Select placeholder="请输入接口状态">
              <Select.Option value={0}>正常请求</Select.Option>
              <Select.Option value={1}>禁止请求</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>

      <Modal key={"edit"} width={350} destroyOnClose={true} forceRender={true} title="批量编辑" visible={isModalEditVisible}
             onOk={edit.submit} okText={"确 认"} onCancel={() => {
        setIsModalEditVisible(false);
        edit.resetFields();
      }}>
        <Form form={edit} name="control-hooks" onFinish={onEditFinish} labelAlign="right" labelCol={{span: 6}}
              wrapperCol={{span: 16}}>
          <Form.Item name={"method"} label="请求方式">
            <Select placeholder="请输入请求方式">
              <Select.Option value="GET">GET</Select.Option>
              <Select.Option value="POST">POST</Select.Option>
              <Select.Option value="UPDATE">UPDATE</Select.Option>
              <Select.Option value="PUT">PUT</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name={"status"} label="状态">
            <Select placeholder="请输入接口状态">
              <Select.Option value={0}>正常请求</Select.Option>
              <Select.Option value={1}>禁止请求</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>

      <ProTable
        // @ts-ignore
        columns={columns}
        defaultSize={"small"}
        form={{layout: "vertical"}}
        headerTitle={'接口列表'}
        actionRef={actionRef}
        formRef={ref}
        rowKey={'id'}
        options={{
          density: true,
          fullScreen: true,
          setting: true,
        }}

        rowSelection={{
          selections: [Table.SELECTION_ALL, Table.SELECTION_INVERT, Table.SELECTION_NONE],
          onChange: (_, selectedRows: any) => {
            setSelectedRows(selectedRows);
          },
        }}
        request={async (params, sort) => {
          const data = await getAccessList(getParams(params, sort));
          return {
            data: data.data.list === undefined || data.data.list === null || data.data.list.length <= 0 ? [] : data.data.list,
            success: data.code === 200,
            total: data.data.total,
          };
        }}
        search={{
          defaultCollapsed: true,
          labelWidth: 'auto',
        }}
        toolBarRender={() => [
          <Button key="primary" type="primary" onClick={() => {
            setIsModalVisible(true)
          }} hidden={controls["system.access:add"] == undefined}>
            {(controls["system.access:add"])?.title || "新建"}
          </Button>,
        ]}/>
      {selectedRowsState?.length > 0 && (
        <FooterToolbar
          extra={
            <div>
              已选择&nbsp;&nbsp;<a style={{fontWeight: 600}}>{selectedRowsState.length}</a>&nbsp;&nbsp;项
            </div>
          }>
          <Button danger={true} onClick={modalDelete} hidden={controls["system.access:delete"] == undefined}>
            批量删除
          </Button>
          <Button type="primary" hidden={controls["system.access:edit"] == undefined}
                  onClick={async () => {
                    setIsModalEditVisible(true);
                  }}>
            批量编辑
          </Button>
        </FooterToolbar>
      )}
    </PageContainer>
  );
};
