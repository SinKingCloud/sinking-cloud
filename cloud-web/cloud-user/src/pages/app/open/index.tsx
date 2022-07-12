import React, {useEffect, useState} from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import {Badge, Button, Card, Form, Input, List, message, Modal, Pagination, Select, Tooltip, Typography} from 'antd';
import styles from './style.less';
import {
  CodeSandboxOutlined,
  DeleteOutlined,
  EditOutlined,
  ExclamationCircleOutlined,
  EyeOutlined,
  PlusOutlined
} from "@ant-design/icons";
import {createApp, deleteApp, getAppList, updateApp} from "@/services/app/app";
import {getRandStr} from "@/util/string";

export default (): React.ReactNode => {
  const nullData = {};
  const [appList, setAppList] = useState([]);
  const [page, setPage] = useState({currentPage: 1, total: 0});
  const [pageLoading, setPageLoading] = useState(false);
  const getApps = (p: number = 1) => {
    setPageLoading(true);
    getAppList({page_size: 11, page: p}).then((r) => {
      setPageLoading(false);
      if (r.code != 200) {
        message.error(r.message || "获取应用列表失败");
      } else {
        setAppList(r.data.list);
        setPage({currentPage: p, total: r.data.total});
      }
    })
  }
  useEffect(() => {
    getApps();
  }, []);
  const content = (
    <div className={styles.pageHeaderContent}>
      <p>
        本平台提供了十分强大的开放能力，当您需要使用平台开放能力时，
        您需要再此页面创建应用并获取密匙以此完成对接。
      </p>
      <div className={styles.contentLink}>
        <a onClick={() => {
          setIsModalVisible(true);
        }}>
          <img alt="" src="https://gw.alipayobjects.com/zos/rmsportal/MjEImQtenlyueSmVEfUD.svg"/>
          创建应用
        </a>
        <a>
          <img alt="" src="https://gw.alipayobjects.com/zos/rmsportal/NbuDUAuBlIApFuDvWiND.svg"/>
          获取密匙
        </a>
        <a>
          <img alt="" src="https://gw.alipayobjects.com/zos/rmsportal/ohOEPSYdDTNnyMbGuyLb.svg"/>
          产品接入
        </a>
      </div>
    </div>
  );
  const extraContent = (
    <div className={styles.extraImg}>
      <img
        alt="图片加载失败"
        src="https://gw.alipayobjects.com/zos/rmsportal/RzwpdLnhmvDJToTdfDPe.png"
      />
    </div>
  );
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [form] = Form.useForm();//新建编辑表单
  const formOnFinish = async (values: any) => {
    let api = createApp;
    if (values?.id || 0 > 0) {
      api = updateApp;
      values.status = values?.status.toString() || "0";
    }
    setIsModalVisible(false);
    form.resetFields();
    api(values).then((r) => {
      if (r.code != 200) {
        message.error(r.message || "请求失败").then()
      } else {
        getApps(page.currentPage)
      }
    });
  }
  const [isViewModalVisible, setIsViewModalVisible] = useState(false);
  const [viewForm] = Form.useForm();//密匙表单
  return (
    <PageContainer loading={pageLoading} content={content} extraContent={extraContent}>
      <Modal key={"form"} destroyOnClose={true} width={"400px"} forceRender={true}
             title={form.getFieldValue("id") == undefined ? "添加应用" : "编辑应用"}
             visible={isModalVisible} onOk={form.submit} okText={"确 认"} onCancel={() => {
        setIsModalVisible(false);
        form.resetFields();
      }}>
        <Form form={form} name="control-hooks" onFinish={formOnFinish} labelAlign="right" labelCol={{span: 6}}
              wrapperCol={{span: 16}}>
          <Form.Item name={"id"} label="ID" hidden={true}>
            <Input placeholder="请输入ID"/>
          </Form.Item>
          <Form.Item name={"name"} label="应用名称" rules={[{required: true}]}>
            <Input placeholder="请输入应用名称"/>
          </Form.Item>
          <Form.Item name={"describe"} label="应用描述" rules={[{required: true}]}>
            <Input.TextArea placeholder="请输入应用描述"/>
          </Form.Item>
          <Form.Item name={"status"} label="应用状态" rules={[{required: true}]}>
            <Select placeholder="请输入应用状态">
              <Select.Option value={0}>正常</Select.Option>
              <Select.Option value={2}>禁用</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
      <Modal key={"view"} width={"400px"} destroyOnClose={true} forceRender={true} title={"应用密匙"}
             visible={isViewModalVisible}
             okText={"重 置"} onOk={() => {
        Modal.confirm({
          title: '确定要重置密匙吗?',
          icon: <ExclamationCircleOutlined/>,
          content: '重置后该应用需使用新的密匙重新接入',
          okType: 'danger',
          onOk() {
            updateApp({id: viewForm.getFieldValue("id"), key: getRandStr(32).toUpperCase()}).then((r) => {
              setIsViewModalVisible(false);
              form.resetFields();
              if (r.code != 200) {
                message.error(r.message || "请求失败").then()
              } else {
                getApps(page.currentPage)
              }
            });
          },
        });
      }} onCancel={() => {
        setIsViewModalVisible(false);
        form.resetFields();
      }}>
        <Form form={viewForm} name="control-hooks2" labelAlign="right" labelCol={{span: 6}}
              wrapperCol={{span: 18}}>
          <Form.Item name={"id"} label="应用ID" rules={[{required: true}]}>
            <Input placeholder="请输入应用ID" disabled/>
          </Form.Item>
          <Form.Item name={"key"} label="应用密匙" rules={[{required: true}]}>
            <Input placeholder="请输入应用密匙" disabled/>
          </Form.Item>
        </Form>
      </Modal>
      <div className={styles.cardList}>
        <List
          loading={false}
          grid={{
            gutter: 16,
            xs: 1,
            sm: 2,
            md: 3,
            lg: 3,
            xl: 4,
            xxl: 4,
          }}
          dataSource={[nullData, ...appList]}
          renderItem={(item: any) => {
            if (item && item?.id) {
              return (
                <List.Item key={item?.id}>
                  <Card
                    hoverable
                    className={styles.card}
                    actions={[
                      <Tooltip title={"查看密匙"}><a key="view" onClick={() => {
                        viewForm.setFieldsValue(item);
                        setIsViewModalVisible(true);
                      }}><EyeOutlined/></a></Tooltip>,
                      <Tooltip title={"编辑应用"}><a key="edit" onClick={() => {
                        form.setFieldsValue(item);
                        setIsModalVisible(true);
                      }}><EditOutlined/></a></Tooltip>,
                      <Tooltip title={"删除应用"}><a key="delete" onClick={() => {
                        Modal.confirm({
                          title: '您确定要删除该应用吗?',
                          icon: <ExclamationCircleOutlined/>,
                          content: '删除后该应用所关联的数据将会转移至主应用',
                          okType: 'danger',
                          onOk() {
                            Modal.confirm({
                              title: '请再次确认是否删除?',
                              icon: <ExclamationCircleOutlined/>,
                              content: '删除后将不可恢复',
                              okType: 'danger',
                              onOk() {
                                deleteApp({id: item?.id}).then((r) => {
                                  if (r.code != 200) {
                                    message.error(r.message || "请求失败").then()
                                  } else {
                                    getApps(page.currentPage)
                                    message.success(r.message).then()
                                  }
                                })
                              },
                            });
                          },
                        });
                      }}><DeleteOutlined/></a></Tooltip>
                    ]}>
                    <Card.Meta
                      avatar={<Tooltip title={"应用状态:" + (item?.status == 0 ? "正常" : "禁用")}><Badge
                        color={item?.status == 0 ? "green" : "red"}
                        dot={true}><CodeSandboxOutlined
                        className={styles.cardAvatar}/></Badge></Tooltip>}
                      title={<a>{item?.name}</a>}
                      description={
                        <Typography.Paragraph
                          className={styles.item}
                          ellipsis={{
                            rows: 3,
                          }}
                        >
                          {item?.describe}
                        </Typography.Paragraph>
                      }
                    />
                  </Card>
                </List.Item>
              );
            }
            return (
              <List.Item>
                <Button type="dashed" className={styles.newButton} onClick={() => {
                  setIsModalVisible(true);
                }}>
                  <PlusOutlined/> 添加应用
                </Button>
              </List.Item>
            );
          }}
        />
      </div>
      <div hidden={page.total <= 11}>
        <Pagination style={{textAlign: "center", marginTop: "20px"}} responsive={true} pageSize={10}
                    defaultCurrent={page?.currentPage || 1}
                    onChange={(p) => {
                      setPage({currentPage: p, total: page.total});
                      getApps(p);
                    }}
                    total={page?.total || 0}/>
      </div>
    </PageContainer>
  );
};
