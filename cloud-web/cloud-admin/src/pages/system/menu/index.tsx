import React, {useEffect, useRef, useState} from 'react';
import {getPageControls, getParams} from "@/util/page";
import ProTable from "@ant-design/pro-table";
import {PageContainer} from "@ant-design/pro-layout";
import {createMenu, deleteMenu, getMenuList, updateMenu} from "@/services/system/menu";
import {getTreeData, toTreeData} from "@/util/tree";
import {Button, Dropdown, Form, Input, InputNumber, Menu, message, Modal, Select, Space, Tag} from "antd";
import * as allIcons from "@ant-design/icons";
import {getAccessList} from "@/services/system/access";

export default (): React.ReactNode => {
  /**
   * 获取控件权限机权限列表
   */
  const [controls, setControls] = useState({});
  const [access, setAccess] = useState([]);
  const [accessMap, setAccessMap] = useState({});
  useEffect(() => {
    getPageControls().then((d) => {
      setControls(d)
    });
    getAccessList({page_size: 1000}).then((r) => {
      if (r.code != 200) {
        message.error(r.message || "请求失败").then()
      } else {
        // @ts-ignore
        const children = [];
        const map = {};
        r.data.list.forEach((item: any) => {
          map[item?.id] = item;
          children.push(<Select.Option key={item?.id} title={item?.id + item?.name + item?.method}
                                       value={item?.id}>{item?.name}-{item?.method}</Select.Option>);
        })
        // @ts-ignore
        setAccess(children)
        setAccessMap(map)
      }
    })
  }, []);

  /**
   * table表格渲染
   */
  const actionRef = useRef();
  const columns = [
    {
      title: '名称',
      dataIndex: 'title',
      tip: '名称',
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
      valueEnum: {
        0: {
          text: '菜单',
        },
        1: {
          text: '控件',
        },
      },
      render: (text: string, record: any) => {
        if (record?.children != undefined && record?.children.length > 0) {
          let isPage = true;
          record?.children.forEach((item: any) => {
            if (item.type == 0) {
              isPage = false
            }
          })
          if (isPage) {
            return <Space size={0}>
              <Tag color="red">
                <allIcons.GoldOutlined/>
                页面</Tag>
            </Space>;
          }
        }
        return <Space size={0}>
          {record?.type == 0 ? <Tag color="blue">
            <allIcons.PaperClipOutlined/>
            菜单</Tag> : <Tag color="green">
            <allIcons.GoldOutlined/>
            控件</Tag>}
        </Space>;
      }
    },
    {
      title: '标识',
      dataIndex: 'key',
      tip: '标识KEY',
    },
    {
      title: '关联权限',
      dataIndex: 'accesses',
      tip: '该菜单关联的权限',
      hideInSearch: true,
      render: (text: string, record: any) => {
        let tempRecord = record;
        const node: any[] = [];
        if (tempRecord.accesses != undefined && tempRecord.accesses != "") {
          let arr = Array.isArray(tempRecord?.accesses) ? tempRecord?.accesses : JSON.parse(tempRecord.accesses);
          if (arr != undefined && arr.length > 0) {
            arr.forEach((item: any) => {
              const info = accessMap[item];
              const color = {
                GET: "#984eff",
                POST: "#ff8847",
                PUT: "#4ec4ff",
                DELETE: "#d246ff",
              };
              node.push(<Tag color={color[info.method] || "red"} key={tempRecord?.id + "-" + item}>{info.name}</Tag>);
            });
          }
        }
        if (node.length <= 0) {
          node.push(<Tag color="default" key={tempRecord?.id + "-0"}>无关联</Tag>);
        }
        return node;
      }
    },
    {
      title: '目录',
      dataIndex: 'path',
      tip: '地址栏目录',
    },
    {
      title: '排序',
      dataIndex: 'sort',
      tip: '排序',
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
                if (!Array.isArray(record?.accesses)) {
                  let arr = [];
                  if (record?.accesses != undefined && record?.accesses != "") {
                    arr = JSON.parse(record?.accesses);
                  }
                  record.accesses = arr;
                }
                //record.accesses = Array.isArray(record?.accesses) ? record?.accesses : JSON.parse(record?.accesses);
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
    values.accesses = JSON.stringify(values.accesses != undefined && values.accesses.length > 0 ? values.accesses : []);
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
              <Select.Option value={1}>控件</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name={"title"} label="标题" rules={[{required: true}]}>
            <Input placeholder="请输入标题"/>
          </Form.Item>
          <Form.Item name={"key"} label="标识" rules={[{required: true}]}>
            <Input placeholder="请输入标识"/>
          </Form.Item>
          <Form.Item name={"accesses"} label="权限">
            <Select mode="multiple" allowClear placeholder="请输入绑定权限"
                    filterOption={(inputValue: string, option?: any) => {
                      return option?.title.indexOf(inputValue) != -1
                    }}>
              {access}
            </Select>
          </Form.Item>
          <Form.Item name={"path"} label="目录"
                     hidden={form.getFieldValue("id") != undefined && form.getFieldValue("type") == 1}>
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
        form={{layout: "vertical",autoFocusFirstInput: false}}
        headerTitle={'菜单管理'}
        actionRef={actionRef}
        rowKey={'id'}
        // @ts-ignore
        columns={columns}
        indentSize={20}
        request={async (params, sort) => {
          params.pageSize = 1000;
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
