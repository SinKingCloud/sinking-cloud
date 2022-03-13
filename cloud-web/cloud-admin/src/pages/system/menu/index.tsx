import React, {useEffect, useRef, useState} from 'react';
import {getPageControls, getParams} from "@/util/page";
import {checkMobile} from "@/util/device";
import ProTable from "@ant-design/pro-table";
import {PageContainer} from "@ant-design/pro-layout";
import {createMenu, deleteMenu, getMenuList, updateMenu} from "@/services/system/menu";
import {getTreeData, toTreeData} from "@/util/tree";
import {Button, Dropdown, Form, Input, InputNumber, Menu, message, Modal, Select, Space, Tag} from "antd";
import * as allIcons from "@ant-design/icons";

export default (): React.ReactNode => {
  /**
   * 控件权限
   */
  const [controls, setControls] = useState({});
  useEffect(() => {
    getPageControls().then((d) => {
      setControls(d)
    })
  }, []);

  /**
   * table表格渲染
   */
  const actionRef = useRef();
  const hidden = !checkMobile();
  const columns = [
    {
      title: '名称',
      dataIndex: 'title',
      tip: '名称',
      ellipsis: hidden,
      render: (test: string, record: any) => {
        return (
          <div>
            {record?.type == 0 && React.createElement(allIcons[record?.icon || "RightOutlined"])} {record?.title}
          </div>
        );
      }
    },
    {
      title: '类型',
      dataIndex: 'type',
      tip: '控件类型',
      ellipsis: hidden,
      valueEnum: {
        0: {
          text: '菜单',
        },
        1: {
          text: '控件',
        },
      },
      render: (text: string, record: any) => {
        return <Space size={0}>
          {record?.type == 0 ? <Tag color="blue">
            <allIcons.PaperClipOutlined/>
            菜单</Tag> : <Tag color="green">
            <allIcons.GoldOutlined/>
            按钮</Tag>}
        </Space>;
      }
    },
    {
      title: '标识',
      dataIndex: 'key',
      tip: '标识KEY',
      ellipsis: hidden,
    },
    {
      title: '目录',
      dataIndex: 'path',
      tip: '地址栏目录',
      ellipsis: hidden,
    },
    {
      title: '排序',
      dataIndex: 'sort',
      tip: '排序',
      ellipsis: hidden,
    },
    {
      title: '状态',
      dataIndex: 'status',
      tip: '接口状态',
      valueEnum: {
        0: {
          text: '正常显示',
          status: 'success',
        },
        1: {
          text: '禁止显示',
          status: 'error',
        },
      },
    }, {
      title: '操作',
      dataIndex: 'Option',
      hideInSearch: true,
      editable: false,
      render: (test: string, record: any) => {
        return <Dropdown overlay={(
          <Menu>
            {record?.type == 0 && <Menu.Item key={"add"} hidden={controls["system.menu:add"] == undefined}>
              <a style={{fontSize: "small"}} onClick={() => {
                form.setFieldsValue({pid: record?.id, type: record.pid == 0 ? 0 : 1})
                setIsFormVisible(true)
              }}>
                添 加
              </a>
            </Menu.Item>}
            <Menu.Item key={"edit"} hidden={controls["system.menu:edit"] == undefined}>
              <a style={{fontSize: "small"}} onClick={() => {
                form.setFieldsValue(record)
                setIsFormVisible(true)
              }}>
                编 辑
              </a>
            </Menu.Item>
            <Menu.Item key={"delete"} hidden={controls["system.menu:delete"] == undefined}>
              <a style={{fontSize: "small"}} onClick={() => {
                deleteSubmit(getIds(getTreeData([record], "children")))
              }}>
                删 除
              </a>
            </Menu.Item>
          </Menu>
        )} trigger={['click']} placement="bottomCenter" arrow={true}>
          <Button size={"small"} style={{fontSize: "small"}} onClick={e => e.preventDefault()}>操
            作 <allIcons.DownOutlined/></Button>
        </Dropdown>;
      }
    },
  ];

  /**
   * 编辑表单
   */
  const [form] = Form.useForm();//新建编辑表单
  const [isFormVisible, setIsFormVisible] = useState(false);
  const onFormFinish = async (values: any) => {
    setIsFormVisible(false);
    let api = createMenu;
    if (values.id != undefined) {
      api = updateMenu;
      values.ids = [values.id];
      values.status = values.status.toString();
      values.sort = values.sort.toString();
      delete values.id;
    }
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
   * 获取数据id
   * @param data
   */
  const getIds = (data: any[]) => {
    let ids: any[] = [];
    data.map((k: any) => {
      ids.push(k.id);
    });
    return ids;
  }

  /**
   * 删除数据
   * @param ids id列表
   */
  const deleteSubmit = (ids: any[]) => {
    Modal.confirm({
      title: '您确定要删除该数据吗?',
      icon: <allIcons.ExclamationCircleOutlined/>,
      content: '连带该子级数据也会删除且不可恢复',
      okType: 'danger',
      onOk() {
        deleteMenu({ids: ids}).then((r) => {
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
      <Modal key={"editForm"} destroyOnClose={true} forceRender={true}
             title={form.getFieldValue("id") == undefined ? "新 建" : "编 辑"}
             visible={isFormVisible} onOk={form.submit} okText={"确 认"} onCancel={() => {
        setIsFormVisible(false);
        form.resetFields();
      }}>
        <Form form={form} name="control-hooks" onFinish={onFormFinish} labelAlign="right" labelCol={{span: 4}}
              wrapperCol={{span: 18}}>
          <Form.Item name={"id"} label="ID" hidden={true}>
            <Input placeholder="请输入ID"/>
          </Form.Item>
          <Form.Item name={"pid"} label="PID" hidden={true}>
            <Input placeholder="请输入PID"/>
          </Form.Item>
          <Form.Item name={"type"} label="类型" rules={[{required: true}]}>
            <Select placeholder="请输入添加的类型"
                    disabled={(form.getFieldValue("pid") == undefined || form.getFieldValue("id") != undefined) || form.getFieldValue("pid") == 0}>
              <Select.Option value={0}>菜单</Select.Option>
              <Select.Option value={1}>按钮</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name={"title"} label="标题" rules={[{required: true}]}>
            <Input placeholder="请输入标题"/>
          </Form.Item>
          <Form.Item name={"key"} label="标识" rules={[{required: true}]}>
            <Input placeholder="请输入标识"/>
          </Form.Item>
          <Form.Item name={"path"} label="目录">
            <Input placeholder="请输入目录"/>
          </Form.Item>
          <Form.Item name={"icon"} label="图标">
            <Input placeholder="请输入图标"/>
          </Form.Item>
          <Form.Item name={"sort"} label="排序">
            <InputNumber placeholder="请输入排序数值"/>
          </Form.Item>
          <Form.Item name={"status"} label="状态" rules={[{required: true}]}>
            <Select placeholder="请输入状态">
              <Select.Option value={0}>正常显示</Select.Option>
              <Select.Option value={1}>禁止显示</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
      <ProTable
        defaultSize={"small"}
        form={{layout: "vertical"}}
        headerTitle={'菜单管理'}
        actionRef={actionRef}
        rowKey={'id'}
        // @ts-ignore
        columns={columns}
        indentSize={20}
        request={async (params, sort) => {
          params.pageSize = 1999;
          //params.type = 0;
          const response = await getMenuList(getParams(params, sort));
          let data = response.data.list === undefined || response.data.list === null || response.data.list.length <= 0 ? [] : response.data.list;
          return {
            data: toTreeData(data, 0, "id", "pid", "children", (item: any) => {
              if (item.children == undefined || item.children.length <= 0) {
                delete item.children;
              }
              return item;
            }),
            success: response.code === 200,
            total: response.data.total,
          };
        }}
        toolBarRender={() => [
          <Button key="primary" type="primary" onClick={() => {
            form.setFieldsValue({pid: 0, type: 0})
            setIsFormVisible(true)
          }} hidden={controls["system.menu:add"] == undefined}>
            {(controls["system.menu:add"])?.title || "新建"}
          </Button>,
        ]}
        search={false}
        pagination={false}/>
    </PageContainer>
  );
};
