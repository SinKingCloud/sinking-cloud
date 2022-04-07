import React, {useEffect, useRef, useState} from 'react';
import {Avatar, Button, Dropdown, Form, Input, Menu, message, Modal, Select, Table, Upload} from 'antd';
import {FooterToolbar, PageContainer} from '@ant-design/pro-layout';
import ProTable from '@ant-design/pro-table';
import {getPageControls, getParams} from "@/util/page";
import {DownOutlined, LoadingOutlined, PlusOutlined} from "@ant-design/icons";
import {getUserList, updateUser} from "@/services/web/user";
import {getUploadUrl} from "@/services/common/upload";

export default (): React.ReactNode => {
  /**
   * 控件权限
   */
  const [controls, setControls] = useState({});
  useEffect(() => {
    getPageControls().then((d) => {
      setControls(d)
    });
  }, []);

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
    if (values.id != undefined) {
      values.ids = [values.id];
      values.status = values.status.toString();
      values.phone = values.phone.toString();
      delete values.id;
    }
    // @ts-ignore
    if (uploadFileList.length > 0) {
      // @ts-ignore
      values.avatar = uploadFileList[0]?.url || ""
    } else {
      values.avatar = ""
    }
    setIsModalVisible(false)
    updateUser(values).then((r) => {
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
    updateUser(values).then((r) => {
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
      title: '头像',
      tip: '用户头像',
      hideInSearch: true,
      render: (text: string, record: any) => {
        return <>
          <Avatar src={record?.avatar}
                  style={{backgroundColor: "green"}}>{record?.nick_name || record?.account}</Avatar>
        </>;
      }
    },
    {
      title: '昵称',
      tip: '用户昵称',
      dataIndex: 'nick_name',
    },
    {
      title: '账号',
      dataIndex: 'account',
      tip: '用户账号',
      copyable: true,
    },
    {
      title: '手机号',
      dataIndex: 'phone',
      tip: '绑定手机',
      copyable: true,
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      tip: '绑定邮箱',
      copyable: true,
    },
    {
      title: '余额',
      dataIndex: 'money',
      tip: '账户余额',
      render: (text: string, record: any) => {
        return record?.money.toFixed(2) + "元"
      }
    },
    {
      title: '登陆IP',
      dataIndex: 'login_ip',
      tip: '上次登录IP',
      copyable: true,
    },
    {
      title: '登陆时间',
      valueType: 'dateTime',
      dataIndex: 'login_time',
      tip: '上次登陆时间',
      sorter: true,
      hideInSearch: true,
      editable: false,
    },
    {
      title: '修改时间',
      valueType: 'dateTime',
      dataIndex: 'update_time',
      tip: '上次修改时间',
      sorter: true,
      hideInSearch: true,
      editable: false,
    },
    {
      title: '注册时间',
      valueType: 'dateTime',
      dataIndex: 'create_time',
      tip: '用户注册时间',
      sorter: true,
      hideInSearch: true,
      editable: false,
    },
    {
      title: '登陆时间',
      valueType: 'dateTimeRange',
      dataIndex: 'login_time',
      tip: '上次登陆时间',
      hideInTable: true,
      search: {
        transform: (value: any) => {
          return {
            login_time_start: value[0],
            login_time_end: value[1],
          };
        },
      },
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
      title: '注册时间',
      valueType: 'dateTimeRange',
      dataIndex: 'create_time',
      tip: '用户注册时间',
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
      title: '实名状态',
      dataIndex: 'real_auth',
      tip: '用户实名状态',
      sorter: true,
      valueEnum: {
        1: {
          text: '已实名',
          status: 'success',
        },
        0: {
          text: '未实名',
          status: 'error',
        },
      },
    },
    {
      title: '账号状态',
      dataIndex: 'status',
      tip: '用户状态',
      sorter: true,
      valueEnum: {
        0: {
          text: '正常',
          status: 'success',
        },
        1: {
          text: '封禁',
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
        return <Dropdown overlay={<Menu>
          <Menu.Item key={"edit"} hidden={controls["web.user:edit"] == undefined}>
            <a style={{fontSize: "small"}} onClick={() => {
              form.setFieldsValue(record);
              if (record?.avatar) {
                // @ts-ignore
                setUploadFileList([{uid: '-1', name: 'image.png', status: 'done', url: record?.avatar}]);
              }
              setIsModalVisible(true);
            }}>
              编 辑
            </a>
          </Menu.Item>
        </Menu>} trigger={['click']} placement="bottomCenter" arrow={true}>
          <Button size={"small"} style={{fontSize: "small"}} onClick={e => e.preventDefault()}>操
            作 <DownOutlined/></Button>
        </Dropdown>;
      }
    },
  ];

  /**
   * 图片上传
   */
  const [uploadLoading, setUploadLoading] = useState(false);
  const [uploadFileList, setUploadFileList] = useState();
  const beforeUpload = (file: any) => {
    const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png';
    if (!isJpgOrPng) {
      message.error('请上传图片格式文件');
      return;
    }
    const isLt2M = file.size / 1024 / 1024 < 2;
    if (!isLt2M) {
      message.error('图片不能大于2MB');
      return;
    }
    return isJpgOrPng && isLt2M;
  }
  const handleChange = (info: any) => {
    setUploadFileList(info.fileList);
    if (info.file.status === 'uploading') {
      setUploadLoading(true);
      return;
    }
    if (info.file.status === 'done') {
      setUploadLoading(false);
      const response = info.file.response;
      if (response?.code != 200) {
        message.error(response?.message || "上传文件错误")
      } else {
        info.file.url = info.file.response.data.url;
        setUploadFileList(info.fileList);
      }
      return;
    }
  };

  return (
    <PageContainer title={false}>
      <Modal key={"form"} destroyOnClose={true} forceRender={true}
             title={form.getFieldValue("id") == undefined ? "新 建" : "编 辑"}
             visible={isModalVisible} onOk={form.submit} okText={"确 认"} onCancel={() => {
        // @ts-ignore
        setUploadFileList([]);
        setIsModalVisible(false);
        form.resetFields();
      }}>
        <Form form={form} name="control-hooks" onFinish={onFormFinish} labelAlign="right" labelCol={{span: 4}}
              wrapperCol={{span: 18}}>
          <Form.Item name={"id"} label="ID" hidden={true}>
            <Input placeholder="请输入ID"/>
          </Form.Item>
          <Form.Item name={"avatar"} label="头像">
            <Upload
              name="file"
              listType="picture-card"
              className="avatar-uploader"
              fileList={uploadFileList}
              showUploadList={true}
              beforeUpload={beforeUpload}
              onChange={handleChange}
              action={getUploadUrl()}
              maxCount={1}
              onPreview={(file) => {
                // @ts-ignore
                if (file?.url) {
                  window.open(file?.url);
                  return;
                }
                // @ts-ignore
                window.open(file?.response.data.url);
              }}
            >
              <div>
                {uploadLoading ? <LoadingOutlined/> : <PlusOutlined/>}
                <div style={{marginTop: 8}}>上传</div>
              </div>
            </Upload>
          </Form.Item>
          <Form.Item name={"account"} label="账号" rules={[{required: true}]}
                     hidden={form.getFieldValue("account") != undefined}>
            <Input placeholder="请输入用户账号"/>
          </Form.Item>
          <Form.Item name={"phone"} label="手机" rules={[{required: true}]}>
            <Input type={"number"} placeholder="请输入用户手机"/>
          </Form.Item>
          <Form.Item name={"status"} label="状态" rules={[{required: true}]}>
            <Select placeholder="请选择用户状态">
              <Select.Option value={0}>正常</Select.Option>
              <Select.Option value={1}>禁止</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name={"password"} label="密码">
            <Input placeholder="不修改则留空"/>
          </Form.Item>
          <Form.Item name={"nick_name"} label="昵称">
            <Input placeholder="请输入用户昵称"/>
          </Form.Item>
          <Form.Item name={"email"} label="邮箱">
            <Input placeholder="请输入用户邮箱"/>
          </Form.Item>

        </Form>
      </Modal>

      <Modal key={"edit"} width={350} destroyOnClose={true} forceRender={true} title="批量编辑" visible={isModalEditVisible}
             onOk={edit.submit} okText={"确 认"} onCancel={() => {
        setIsModalEditVisible(false);
        edit.resetFields();
      }}>
        <Form form={edit} name="control-hooks"  onFinish={onEditFinish} labelAlign="right" labelCol={{span: 6}}
              wrapperCol={{span: 16}}>
          <Form.Item name={"status"} label="状态">
            <Select placeholder="请选择用户状态">
              <Select.Option value={0}>正常</Select.Option>
              <Select.Option value={1}>封禁</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>

      <ProTable
        // @ts-ignore
        columns={columns}
        defaultSize={"small"}
        form={{layout: "vertical"}}
        headerTitle={'用户列表'}
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
          const data = await getUserList(getParams(params, sort));
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
      />
      {selectedRowsState?.length > 0 && (
        <FooterToolbar
          extra={
            <div>
              已选择&nbsp;&nbsp;<a style={{fontWeight: 600}}>{selectedRowsState.length}</a>&nbsp;&nbsp;项
            </div>
          }>
          <Button type="primary" hidden={controls["web.user:edit"] == undefined}
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
