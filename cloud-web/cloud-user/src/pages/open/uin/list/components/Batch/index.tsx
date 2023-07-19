import React, {useEffect, useState} from 'react';
import {Row, Col, Form, List, message, Modal, Select, Tabs, Button, TimePicker, Avatar, Spin} from "antd";
import {
  changeUinArea, changeUinTiming, deleteUin,
  getServerList,
  uinRun,
  uinStatus,
} from "@/services/open/uin/uin";
import moment from "moment";
import {ExclamationCircleOutlined, LoadingOutlined} from "@ant-design/icons";
import styles from './index.less';

export interface BatchRef {
  Show: (uin: any) => void;
}

export interface BatchProps {
  actionRef?: any;
  visible: boolean;
  onCancel: () => void;
}

const Batch = React.forwardRef<BatchRef, BatchProps>((props, ref) => {
  const {visible, actionRef, onCancel} = props;

  /**
   * 批量操作
   */
  const [selectedRowsState, setSelectedRows] = useState<any>([]);
  const [activeKey, setActiveKey] = useState("1");
  const onTabChange = (key: string) => {
    setActiveKey(key);
  };
  const updateRowsState = (index: any, newValue: any) => {
    const newArray = selectedRowsState.map((item: any, i: any) => {
      if (i === index) {
        return newValue; // 更新指定索引的元素
      }
      return item; // 保持其他元素不变
    });
    setSelectedRows(newArray); // 使用新的数组副本更新state
  };

  const [serverKey, setServerKey] = useState("");
  const [loading, setLoading] = useState(false)
  const [form] = Form.useForm();
  const [timing, setTiming] = useState("");
  const changeServerAndTiming = () => {
    Modal.confirm({
      title: '确定要保存设置吗',
      icon: <ExclamationCircleOutlined/>,
      content: '您的设置将立即生效',
      okType: 'primary',
      onOk() {
        setLoading(true)
        const rows = selectedRowsState
        const values = form?.getFieldsValue();
        if (values != undefined && values?.server != undefined && values?.server != "") {
          setTimeout(async () => {
            for (let i = 0; i < rows?.length; i++) {
              const row = rows[i]
              // @ts-ignore`
              row["job" + activeKey + "_status"] = 1;
              const r = await changeUinArea({uin: row?.uin, server_key: values?.server})
              if (r?.code == 200) {
                // @ts-ignore
                row["job" + activeKey + "_status"] = 2;
                setServerKey(values?.server);
                setLoading(false)
              } else {
                // @ts-ignore
                row["job" + activeKey + "_status"] = 3;
                setLoading(false)
              }
              updateRowsState(i, row);
            }
            // @ts-ignore
            actionRef?.current?.reload()
          }, 0)
        }
        if (values != undefined && values?.time != undefined) {
          const time = values?.time?._d
          const formattedTime = moment(time).format('HH:mm');
          const tempArr = formattedTime.split(":");
          setTimeout(async () => {
            for (let i = 0; i < rows?.length; i++) {
              const row = rows[i]
              // @ts-ignore`
              row["job" + activeKey + "_status"] = 1;
              const r = await changeUinTiming({uin: row?.uin, timing: parseInt(tempArr[0]) * 60 + parseInt(tempArr[1])})
              if (r?.code == 200) {
                // @ts-ignore
                row["job" + activeKey + "_status"] = 2;
                setTiming(formattedTime);
                setLoading(false)
              } else {
                // @ts-ignore
                row["job" + activeKey + "_status"] = 3;
                setLoading(false)
              }
              updateRowsState(i, row);
            }
            // @ts-ignore
            actionRef?.current?.reload()
          }, 0)
        }
      },
    });
  }
  const suspendAfk = () => {
    Modal.confirm({
      title: '确定要暂停挂机吗?',
      icon: <ExclamationCircleOutlined/>,
      content: '您的设置将会立即生效',
      okType: 'primary',
      onOk() {
        setLoading(true)
        setTimeout(async () => {
          const rows = selectedRowsState
          for (let i = 0; i < rows?.length; i++) {
            const row = rows[i]
            // @ts-ignore`
            row["job" + activeKey + "_status"] = 1;
            const r = await uinStatus({uin: row?.uin, status: 1})
            if (r?.code == 200) {
              // @ts-ignore
              row["job" + activeKey + "_status"] = 2;
              setLoading(false)
            } else {
              // @ts-ignore
              row["job" + activeKey + "_status"] = 3;
              setLoading(false)
            }
            updateRowsState(i, row);
          }
          // @ts-ignore
          actionRef?.current?.reload()
        }, 0)
      },
    });
  }

  const restoreAfk = () => {
    Modal.confirm({
      title: '确定要恢复挂机吗?',
      icon: <ExclamationCircleOutlined/>,
      content: '您的设置将会立即生效',
      okType: 'primary',
      onOk() {
        setLoading(true)
        setTimeout(async () => {
          const rows = selectedRowsState
          for (let i = 0; i < rows?.length; i++) {
            const row = rows[i]
            // @ts-ignore`
            row["job" + activeKey + "_status"] = 1;
            const r = await uinStatus({uin: row?.uin, status: 0})
            if (r?.code == 200) {
              // @ts-ignore
              row["job" + activeKey + "_status"] = 2;
              setLoading(false)
            } else {
              // @ts-ignore
              row["job" + activeKey + "_status"] = 3;
              setLoading(false)
            }
            updateRowsState(i, row);
          }
          // @ts-ignore
          actionRef?.current?.reload()
        }, 0)
      },
    });
  }
  const deleteAfk = () => {
    Modal.confirm({
      title: '确定要删除挂机吗?',
      icon: <ExclamationCircleOutlined/>,
      content: '删除挂机后将不可恢复',
      okType: 'danger',
      onOk() {
        setLoading(true)
        setTimeout(async () => {
          const rows = selectedRowsState
          for (let i = 0; i < rows?.length; i++) {
            const row = rows[i]
            // @ts-ignore`
            row["job" + activeKey + "_status"] = 1;
            const r = await deleteUin({uin: row?.uin})
            if (r?.code == 200) {
              // @ts-ignore
              row["job" + activeKey + "_status"] = 2;
              setLoading(false)
            } else {
              // @ts-ignore
              row["job" + activeKey + "_status"] = 3;
              setLoading(false)
            }
            updateRowsState(i, row);
          }
          // @ts-ignore
          actionRef?.current?.reload()
        }, 0)
      },
    });
  };
  const executeAfk = () => {
    Modal.confirm({
      title: '确定要执行挂机吗?',
      icon: <ExclamationCircleOutlined/>,
      content: '提交挂机后将会立即执行',
      okType: 'primary',
      onOk() {
        setLoading(true)
        setTimeout(async () => {
          const rows = selectedRowsState
          for (let i = 0; i < rows?.length; i++) {
            const row = rows[i]
            // @ts-ignore`
            row["job" + activeKey + "_status"] = 1;
            const r = await uinRun({uin: row?.uin});
            if (r?.code == 200) {
              // @ts-ignore
              row["job" + activeKey + "_status"] = 2;
              setLoading(false)
            } else {
              // @ts-ignore
              row["job" + activeKey + "_status"] = 3;
              setLoading(false)
            }
            updateRowsState(i, row);
          }
          // @ts-ignore
          actionRef?.current?.reload()
        }, 0)
      },
    });
  };
  const antIcon = <LoadingOutlined style={{fontSize: 14}} spin/>
  const cancelBtnClick = () => {
    selectedRowsState?.map((item: any, index: any) => {
      form?.resetFields()
      item["job1_status"] = 0;
      item["job2_status"] = 0;
      item["job3_status"] = 0;
      item["job4_status"] = 0;
      item["job5_status"] = 0;
      updateRowsState(index, item)
    })
    onCancel()
  }

  const renderItem = (item: any) => (
    <List.Item
      key={
        // @ts-ignore
        "renderItem_" + activeKey + "_" + item?.id || ""}
      actions={[
        <>
          <Spin spinning={loading} indicator={antIcon}/>
          {(item?.["job" + activeKey + "_status"] || 0) == "0" && "待执行"}
          {item?.["job" + activeKey + "_status"] == 1 && "正在执行"}
          {item?.["job" + activeKey + "_status"] == 2 && "执行成功"}
          {item?.["job" + activeKey + "_status"] == 3 && "执行失败"}
        </>
      ]}
    >
      <List.Item.Meta
        avatar={<Avatar size={30}
                        src={// @ts-ignore
                          "https://q1.qlogo.cn/g?b=qq&nk=" + item?.uin + "&s=100&t=20190225"}/>}
        title={
          <div>
            {// @ts-ignore
              item?.uin || ""}
          </div>
        }
      />
    </List.Item>
  );
  /**
   * 挂机地区列表
   */
  const [servers, setServers] = useState({})
  /**
   * 初始化数据
   */
  const Show = (arr: any) => {
    setSelectedRows(arr)
  };
  useEffect(() => {
    getServerList({page_size: 999}).then((r) => {
      if (r.code != 200) {
        message.error(r.message || "请求失败").then();
      } else {
        const map = {};
        r?.data?.list?.forEach((item: any) => {
          map[item?.key] = item;
        });
        setServers(map);
      }
    });
  }, []);

  // 将 Test 方法绑定到 ref 对象的 current 属性上
  React.useImperativeHandle(ref, () => ({
    Show: Show,
  }));

  return (
    <>
      <Modal key={"Batch"} destroyOnClose={true} forceRender={true} title={"批量操作"}
             visible={visible} onCancel={cancelBtnClick} footer={<>
        {
          <Button onClick={cancelBtnClick}>
            取消
          </Button>
        }
        {
          activeKey == "1" && <Button type="primary"
                                      onClick={executeAfk}>
            开始执行
          </Button>
        }
        {
          activeKey == "2" && <Button type="primary"
                                      onClick={async () => {
                                        suspendAfk()
                                      }}>
            确认暂停
          </Button>
        }
        {
          activeKey == "3" && <Button type="primary"
                                      onClick={async () => {
                                        restoreAfk()
                                      }}>
            确认恢复
          </Button>
        }
        {
          activeKey == "4" && <Button type="primary"
                                      onClick={async () => {
                                        deleteAfk()
                                      }}>
            确认删除
          </Button>
        }
        {
          activeKey == "5" && <Button type="primary"
                                      onClick={async () => {
                                        changeServerAndTiming()
                                      }}>
            保存设置
          </Button>
        }
      </>}>
        <div style={{marginTop: "-30px"}} className={styles.operationModel}>
          <Tabs defaultActiveKey="1" activeKey={activeKey} onChange={onTabChange}>
            <Tabs.TabPane tab="执行挂机" key="1">
            </Tabs.TabPane>
            <Tabs.TabPane tab="暂停挂机" key="2">
            </Tabs.TabPane>
            <Tabs.TabPane tab="恢复挂机" key="3">
            </Tabs.TabPane>
            <Tabs.TabPane tab="删除挂机" key="4">
            </Tabs.TabPane>
            <Tabs.TabPane tab="修改信息" key="5">
              <Form
                form={form}
                layout="vertical"
              >
                <Row>
                  <Col span={12}>
                    <Form.Item name="server">
                      <Select key={"serverKey"} placeholder="请选择挂机地区" style={{width: "95%"}}
                              value={serverKey}
                              bordered={true}>
                        {Object.keys(servers).map((k) => {
                          return <Select.Option key={k} value={k}>{servers[k]?.name || "默认"}</Select.Option>;
                        })}
                      </Select>
                    </Form.Item>
                  </Col>
                  <Col span={12}>
                    <Form.Item name="time">
                      <TimePicker key={"timing"} bordered={true} allowClear={false}
                                  suffixIcon={false}
                                  popupClassName={styles['time']}
                                  placeholder="请选择挂机时间"
                                  style={{textAlign: "center", width: "100%"}}
                                  value={moment(timing, 'HH:mm')}
                                  format={'HH:mm'}/>
                    </Form.Item>
                  </Col>
                </Row>
              </Form>
            </Tabs.TabPane>
          </Tabs>
          <List
            dataSource={selectedRowsState}
            renderItem={renderItem}
          />
        </div>
      </Modal>
    </>
  );
});
export default Batch;
