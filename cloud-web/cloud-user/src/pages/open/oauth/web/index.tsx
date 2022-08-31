import React, {useEffect, useRef, useState} from 'react';
import {Button, Dropdown, Form, Input, Menu, message, Modal, Select} from 'antd';
import {PageContainer} from '@ant-design/pro-layout';
import ProTable from '@ant-design/pro-table';
import {getParams} from "@/util/page";
import {DownOutlined, ExclamationCircleOutlined} from "@ant-design/icons";
import {createOauth, deleteOauth, getOauthList, updateOauth} from "@/services/open/oauth/oauth";
import {getAppList} from "@/services/app/app";

export default (): React.ReactNode => {
  /**
   * 表单处理
   */
  const actionRef = useRef();
  const ref = useRef();
  const [isModalVisible, setIsModalVisible] = useState(false);

  /**
   * 应用列表
   */
  const [apps, setApps] = useState({});
  useEffect(() => {
    getAppList({page_size: 999, page: 1}).then((r: any) => {
      if (r.code != 200) {
        message.error(r.message || "请求失败").then()
      } else {
        const map = {};
        r.data.list.forEach((item: any) => {
          map[item?.id] = item;
        });
        setApps(map);
      }
    });
  }, []);
  /**
   * 表单
   */
  const [form] = Form.useForm();//新建编辑表单
  /**
   * 新建编辑提交表单
   * @param values 表单项
   */
  const onFormFinish = async (values: any) => {
    let api = createOauth;
    if (values.id != undefined) {
      api = updateOauth;
      values.status = values.status.toString();
      values.app_id = values.app_id.toString();
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
   * 删除数据
   * @param id id列表
   */
  const deleteSubmit = (id: number) => {
    Modal.confirm({
      title: '您确定要删除该数据吗?',
      icon: <ExclamationCircleOutlined/>,
      content: '删除后该数据不可恢复',
      okType: 'danger',
      onOk() {
        deleteOauth({id: id}).then((r) => {
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
   * table表格渲染
   */
  const columns = [
    {
      title: 'OauthID',
      dataIndex: 'id',
      tip: 'Oauth网站ID',
      hideInSearch: true,
      sorter: true,
    },
    {
      title: '名称',
      dataIndex: 'name',
      tip: '接口名称',
    },
    {
      title: '域名',
      dataIndex: 'domain',
      tip: '接入地址',
    },
    {
      title: '所属应用',
      dataIndex: 'app_id',
      tip: '账号所属应用',
      render: (text: string, record: any) => {
        return apps[record?.app_id]?.name || "主应用";
      },
      renderFormItem: () => {
        const list = new Array();
        Object.keys(apps).forEach((key: any) => {
          list.push(<Select.Option key={apps[key]?.id} value={apps[key]?.id}>{apps[key]?.name}</Select.Option>);
        });
        return <Select placeholder="请选择所属应用" allowClear={true}>
          {list}
        </Select>;
      }
    },
    {
      title: '状态',
      dataIndex: 'status',
      tip: '接入状态',
      sorter: true,
      valueEnum: {
        0: {
          text: '正常',
          status: 'success',
        },
        1: {
          text: '暂停',
          status: 'warning',
        },
        2: {
          text: '封禁',
          status: 'error',
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
      valueType: 'dateTimeRange',
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
      valueType: 'dateTimeRange',
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
      title: '操作',
      dataIndex: 'Option',
      hideInSearch: true,
      editable: false,
      render: (text: string, record: any) => {
        return <Dropdown overlay={(
          <Menu>
            <Menu.Item key={"edit"}>
              <a style={{fontSize: "small"}} onClick={() => {
                form.setFieldsValue(record)
                setIsModalVisible(true);
              }}>
                编 辑
              </a>
            </Menu.Item>
            <Menu.Item key={"delete"}>
              <a style={{fontSize: "small"}} onClick={() => {
                // @ts-ignore
                deleteSubmit(record?.id)
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
      <Modal key={"form"} destroyOnClose={true} width={380} forceRender={true}
             title={form.getFieldValue("id") == undefined ? "新 建" : "编 辑"}
             visible={isModalVisible} onOk={form.submit} okText={"确 认"} onCancel={() => {
        setIsModalVisible(false);
        form.resetFields();
      }}>
        <Form form={form} name="control-hooks" onFinish={onFormFinish} labelAlign="right" labelCol={{span: 6}}
              wrapperCol={{span: 17}}>
          <Form.Item name={"id"} label="ID" hidden={true}>
            <Input placeholder="请输入ID"/>
          </Form.Item>
          <Form.Item name={"name"} label="网站名称" rules={[{required: true}]}>
            <Input placeholder="请输入名称"/>
          </Form.Item>
          <Form.Item name={"domain"} label="网站域名" rules={[{required: true}]}>
            <Input placeholder="请输入域名"/>
          </Form.Item>
          <Form.Item name={"app_id"} label="绑定应用" rules={[{required: true}]}>
            <Select placeholder="请选择绑定应用">
              {
                Object.values(apps)?.map((k: any) => {
                  return <Select.Option key={k?.id} value={k?.id}>{k?.name}</Select.Option>
                })
              }
            </Select>
          </Form.Item>
          <Form.Item name={"status"} label="接入状态" rules={[{required: true}]}>
            <Select placeholder="请选择接入状态">
              <Select.Option value={0}>正常</Select.Option>
              <Select.Option value={1}>暂停</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
      <ProTable
        // @ts-ignore
        columns={columns}
        defaultSize={"small"}
        form={{layout: "vertical",autoFocusFirstInput: false}}
        headerTitle={'接口列表'}
        actionRef={actionRef}
        formRef={ref}
        rowKey={'id'}
        options={{
          density: true,
          fullScreen: true,
          setting: true,
        }}
        request={async (params, sort) => {
          const data = await getOauthList(getParams(params, sort));
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
          labelWidth: 'auto',
        }}
        toolBarRender={() => [
          <Button key="primary" type="primary" onClick={() => {
            setIsModalVisible(true)
          }}>
            添加
          </Button>,
        ]}/>
    </PageContainer>
  );
};
