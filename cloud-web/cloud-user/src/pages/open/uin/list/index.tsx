import React, {useRef, useState} from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import ProTable from '@ant-design/pro-table';
import {getParams} from "@/util/page";
import {Avatar, Button, Dropdown, Menu, message, Tag} from "antd";
import {getUinList} from "@/services/open/uin/uin";
import {DownOutlined} from "@ant-design/icons";
import UpdateUin from "@/pages/open/uin/list/components/UpdateUin";

export default (): React.ReactNode => {
  /**
   * 表单处理
   */
  const actionRef = useRef();
  const ref = useRef();
  /**
   * 挂机地区列表
   */
  const [servers, setServers] = useState({})

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

            }}>
              执行挂机
            </a>
          </Menu.Item>
          <Menu.Item key={"order"}>
            <a style={{fontSize: "small"}} onClick={() => {

            }}>
              订单记录
            </a>
          </Menu.Item>
          <Menu.Item key={"log"}>
            <a style={{fontSize: "small"}} onClick={() => {
              setUinVisible(true)
            }}>
              操作记录
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
  const [pwd, setPwd] = useState("");
  const [server, setServer] = useState("");
  const showUinModal = (uin: string, server: string, pwd: string = "") => {
    setUin(uin);
    setServer(server);
    setPwd(pwd)
    setUinVisible(true);
  }
  const hideUinModal = () => {
    setUin("");
    setServer("");
    setPwd("")
    setUinVisible(false);
  }

  return (
    <PageContainer title={false}>
      <UpdateUin visible={uinVisible} uin={uin} pwd={pwd} server={server} onOk={() => {
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

          }}>
            添加
          </Button>,
        ]}
      />
    </PageContainer>
  );
};
