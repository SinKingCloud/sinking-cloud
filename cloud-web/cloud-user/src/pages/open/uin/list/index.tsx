import React, {useEffect, useRef, useState} from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import ProTable from '@ant-design/pro-table';
import {getParams} from "@/util/page";
import {Avatar, Button, Dropdown, Menu, message, Modal, Select, Tag} from "antd";
import {getUinList, uinRun} from "@/services/open/uin/uin";
import {DownOutlined, ExclamationCircleOutlined} from "@ant-design/icons";
import UpdateUin from "@/pages/open/uin/list/components/UpdateUin";
import {getAppList} from "@/services/app/app";

export default (): React.ReactNode => {
  /**
   * 表单处理
   */
  const actionRef = useRef();
  const ref = useRef();
  /**
   * 挂机地区列表
   */
  const [servers, setServers] = useState({});
  /**
   * 应用列表
   */
  const [apps, setApps] = useState({});
  useEffect(() => {
    getAppList({page_size: 9999, page: 1}).then((r: any) => {
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
   * 执行挂机
   * @param uin 账号
   */
  const uinPlayer = (uin: string) => {
    Modal.confirm({
      title: '确定要执行挂机吗?',
      icon: <ExclamationCircleOutlined/>,
      content: '提交挂机后将会立即执行',
      okType: 'danger',
      onOk() {
        const key = "uinPlayer";
        message.loading({content: '正在提交挂机指令', key, duration: 60}).then();
        uinRun({uin: uin}).then((r) => {
          if (r.code == 200) {
            message.success({content: r.message, key, duration: 2}).then();
          } else {
            message.error({content: r.message, key, duration: 2}).then();
          }
        });
      },
    });
  };
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
      title: '头像',
      tip: 'QQ头像',
      hideInSearch: true,
      render: (text: string, record: any) => {
        return <>
          <Avatar src={"https://q1.qlogo.cn/g?b=qq&nk=" + record?.uin + "&s=100&t=20190225"}
                  style={{backgroundColor: "green"}}>{record?.uin}</Avatar>
        </>;
      }
    },
    {
      title: '账号',
      dataIndex: 'uin',
      tip: '挂机账号',
      copyable: true,
    },
    {
      title: 'QQ状态',
      dataIndex: 'uin_status',
      tip: 'QQ挂机状态,异常状态需要重新登陆才能恢复',
      valueEnum: {
        0: {
          text: '正常',
          status: 'success',
        },
        1: {
          text: '异常',
          status: 'error',
        },
      },
    },
    {
      title: '挂机剩余天数',
      dataIndex: 'player_day',
      tip: '剩余挂机到期时间',
      hideInSearch: true,
      render: (text: string, record: any) => {
        if (record?.player_day <= 0 && record?.sign_day <= 0) {
          return (<Tag color="error">未下单或已到期</Tag>);
        }
        return (
          <>
            {record?.player_day > 0 && <Tag color="processing">加速:{record?.player_day + "天"}</Tag>}
            {record?.player_day <= 0 && <Tag color="warning">加速:无</Tag>}
            {record?.sign_day > 0 && <Tag color="success">扩展:{record?.sign_day + "天"}</Tag>}
            {record?.sign_day <= 0 && <Tag color="warning">扩展:无</Tag>}
          </>
        );
      }
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
      title: '挂机地区',
      dataIndex: 'server_key',
      tip: '账号挂机地区',
      hideInSearch: true,
      render: (text: string, record: any) => {
        return servers[record?.server_key]?.name || record?.server_key;
      }
    },
    {
      title: '挂机时间',
      dataIndex: 'timing',
      tip: '定时挂机时间',
      hideInSearch: true,
      render: (text: string, record: any) => {
        return (Array(2).join("0") + parseInt(String(record?.timing / 60))).slice(-2) + ":" + (Array(2).join("0") + parseInt(String(record?.timing % 60))).slice(-2);
      }
    },
    {
      title: '已挂机天数',
      dataIndex: 'sign_day',
      tip: '在系统累计已挂机天数',
      hideInSearch: true,
      render: (text: string, record: any) => {
        return record?.run_day + "天";
      }
    },
    {
      title: '功能状态',
      dataIndex: 'status',
      tip: 'QQ的功能状态',
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
      title: '创建时间',
      valueType: 'dateTime',
      dataIndex: 'create_time',
      tip: '创建时间',
      hideInSearch: true,
    },
    {
      title: '运行时间',
      valueType: 'dateTime',
      dataIndex: 'run_time',
      tip: '上次挂机时间',
      hideInSearch: true,
    },
    {
      title: '创建时间',
      valueType: 'dateTimeRange',
      dataIndex: 'create_time',
      tip: '创建时间',
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
        let hide = record?.player_day <= 0 && record?.sign_day <= 0;
        return <Dropdown overlay={<Menu>
          <Menu.Item key={"password"} hidden={hide}>
            <a style={{fontSize: "small"}} onClick={() => {
              showUinModal(record?.uin, record?.sever_key);
            }}>
              更新挂机
            </a>
          </Menu.Item>
          <Menu.Item key={"buy"}>
            <a style={{fontSize: "small"}} onClick={() => {

            }}>
              挂机续期
            </a>
          </Menu.Item>
          <Menu.Item key={"edit"} hidden={hide}>
            <a style={{fontSize: "small"}} onClick={() => {

            }}>
              编辑功能
            </a>
          </Menu.Item>
          <Menu.Item key={"delete"}>
            <a style={{fontSize: "small"}} onClick={() => {

            }}>
              删除挂机
            </a>
          </Menu.Item>
          <Menu.Item key={"run"} hidden={hide}>
            <a style={{fontSize: "small"}} onClick={() => {
              uinPlayer(record?.uin);
            }}>
              执行挂机
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
   * 挂机账号添加更新
   */
  const [uinVisible, setUinVisible] = useState(false);
  const [uin, setUin] = useState("");
  const showUinModal = (uin: string = "", server: string = "", pwd: string = "") => {
    setUin(uin);
    setUinVisible(true);
  }
  const hideUinModal = () => {
    setUin("");
    setUinVisible(false);
  }

  return (
    <PageContainer title={false}>
      <UpdateUin visible={uinVisible} uin={uin} onOk={() => {
        hideUinModal()
      }} onCancel={() => {
        hideUinModal()
      }} onRender={(s) => {
        setServers(s);
      }}/>
      <ProTable
        defaultSize={"small"}
        form={{layout: "vertical"}}
        headerTitle={'挂机列表'}
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
          const data = await getUinList(getParams(params, sort));
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
            showUinModal();
          }}>
            添加
          </Button>,
        ]}
      />
    </PageContainer>
  );
};
