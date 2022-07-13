import React, {useEffect, useState} from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import {
  Badge,
  Button,
  Card,
  Descriptions, Empty,
  List,
  message,
  Modal,
  Pagination, Spin,
  Tag,
  Timeline,
  Tooltip,
  Typography
} from 'antd';
import styles from './style.less';
import {
  CloudDownloadOutlined, CloudSyncOutlined, CodeSandboxOutlined, ExclamationCircleOutlined,
  EyeOutlined,
  PlusOutlined,
  ShoppingCartOutlined, SlackOutlined,
} from "@ant-design/icons";
import {Download, getAuthList, getShopList} from "@/services/app/auth";
import {getUpdateLog} from "@/services/app/app";

export default (): React.ReactNode => {
  const nullData = {};
  const [pageLoading, setPageLoading] = useState(false);
  /**
   * 我的授权
   */
  const [appList, setAppList] = useState([]);
  const [page, setPage] = useState({currentPage: 1, total: 0});
  const getApps = (p: number = 1) => {
    setPageLoading(true);
    getAuthList({page_size: 11, page: p}).then((r) => {
      setPageLoading(false);
      if (r.code != 200) {
        message.error(r.message || "获取应用列表失败");
      } else {
        setAppList(r.data.list);
        setPage({currentPage: p, total: r.data.total});
      }
    })
  }
  /**
   * 应用商店
   */
  const [shopList, setShopList] = useState([]);
  const [shopPage, setShopPage] = useState({currentPage: 1, total: 0});
  const getShops = (p: number = 1) => {
    setPageLoading(true);
    getShopList({page_size: 11, page: p}).then((r) => {
      setPageLoading(false);
      if (r.code != 200) {
        message.error(r.message || "获取商品列表失败");
      } else {
        setShopList(r.data.list);
        setShopPage({currentPage: p, total: r.data.total});
      }
    });
  }

  /**
   * 初始化
   */
  useEffect(() => {
    getShops();
  }, []);

  /**
   * tab操作
   */
  const [tabStatus, setTabStatus] = useState({
    operationKey: 'detail',
    tabActiveKey: 'shop',
  });
  const onTabChange = (tabActiveKey: any) => {
    if (tabActiveKey == 'detail') {
      getApps();
    } else if (tabActiveKey == 'shop') {
      getShops();
    }
    setTabStatus({...tabStatus, tabActiveKey});
  };
  const content = (
    <div className={styles.pageHeaderContent}>
      <p>
        本平台研发了众多产品,在此页面您可购买使用许可并下载程序源代码。
      </p>
      <div className={styles.contentLink}>
        <a>
          <img alt="" src="https://gw.alipayobjects.com/zos/rmsportal/MjEImQtenlyueSmVEfUD.svg"/>
          购买应用
        </a>
        <a>
          <img alt="" src="https://gw.alipayobjects.com/zos/rmsportal/NbuDUAuBlIApFuDvWiND.svg"/>
          安装应用
        </a>
        <a>
          <img alt="" src="https://gw.alipayobjects.com/zos/rmsportal/ohOEPSYdDTNnyMbGuyLb.svg"/>
          产品上线
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
  /**
   * 程序详情弹窗
   */
  const [isViewModalVisible, setIsViewModalVisible] = useState(false);
  const [authInfo, setAuthInfo] = useState<any>({});
  /**
   * 程序详情弹窗
   */
  const [isInfoModalVisible, setIsInfoModalVisible] = useState(false);
  const [shopInfo, setShopInfo] = useState<any>({});
  /**
   * 更新日志弹窗
   */
  const [isUpdateModalVisible, setIsUpdateModalVisible] = useState(false);
  const [updateLogLoading, setUpdateLogLoading] = useState(false);
  const [updateLog, setUpdateLog] = useState<any>([]);
  const showAppUpdate = (id: number) => {
    setIsUpdateModalVisible(true);
    setUpdateLogLoading(true);
    getUpdateLog({auth_id: id}).then((r) => {
      if (r?.code == 200) {
        setUpdateLog(r?.data?.list || []);
      }
      setUpdateLogLoading(false);
    });
  };
  /**
   * 下载程序
   */
  const downloadApp = (id: number) => {
    Modal.confirm({
      title: '确定要下载程序吗?',
      icon: <ExclamationCircleOutlined/>,
      content: '系统将生成您的专属指纹程序',
      okType: 'primary',
      onOk() {
        const key = "downloadApp";
        message.loading({
          content: '正在构建程序',
          key,
          duration: 600,
        }).then();
        Download({id: id}).then((r) => {
          if (r?.code == 200) {
            message.destroy(key);
            Modal.success({
              content: '构建程序成功,请点击按钮开始下载',
              okText: '下载',
              onOk() {
                window.location.href = r?.data;
              }
            });
          } else {
            message.error({content: r?.message || "请求错误", key, duration: 2}).then();
          }
        });
      },
    });
  };
  return (
    <PageContainer loading={pageLoading} content={content} extraContent={extraContent}
                   tabActiveKey={tabStatus.tabActiveKey}
                   onTabChange={onTabChange}
                   tabList={[
                     {
                       key: 'shop',
                       tab: '应用商店',
                     },
                     {
                       key: 'detail',
                       tab: '我的授权',
                     },
                   ]}>
      <Modal key={"updates"} width={"450px"} destroyOnClose={true} forceRender={true} footer={null} title={"更新日志"}
             visible={isUpdateModalVisible}
             onCancel={() => {
               setIsUpdateModalVisible(false);
             }}>
        <div hidden={!updateLogLoading} style={{textAlign: "center"}}>
          <Spin spinning={updateLogLoading}/>
        </div>
        <div hidden={updateLogLoading}>
          <Timeline>
            {updateLog?.map((item: any) => {
              return <Timeline.Item key={item?.id}>
                <div>
                  <span style={{fontWeight: "bolder", color: "#676767"}}>
                    {item?.create_time}
                  </span>
                  &nbsp;&nbsp;&nbsp;
                  <Tag>no.{item?.version_number}</Tag>
                </div>
                <div style={{marginTop: "7px"}}>
                  <Tag color={"processing"}>{item?.title}</Tag>
                  <Tag color={"success"}>版本:{item?.version}</Tag>
                </div>
                <div style={{fontSize: "10px", marginTop: "7px"}} dangerouslySetInnerHTML={{
                  __html: item?.content,
                }}>
                </div>
              </Timeline.Item>;
            })}
          </Timeline>
        </div>
      </Modal>
      {tabStatus.tabActiveKey == "shop" && (<>
        <Modal key={"info"} width={"450px"} destroyOnClose={true} forceRender={true} footer={null} title={"应用详情"}
               visible={isInfoModalVisible}
               bodyStyle={{padding: "5px"}}
               onCancel={() => {
                 setIsInfoModalVisible(false);
               }}>
          <Descriptions
            bordered
            column={1}
          >
            <Descriptions.Item label="应用名称">{shopInfo?.name || "未知应用"}(ID:{shopInfo?.id || 0})</Descriptions.Item>
            <Descriptions.Item label="应用价格">{shopInfo?.price?.toFixed(2) || "0.00"}元</Descriptions.Item>
            <Descriptions.Item label="应用描述">{shopInfo?.describe || "无描述"}</Descriptions.Item>
            <Descriptions.Item label="应用状态">
              {shopInfo?.status == 0 && "正常"}{shopInfo?.status == 1 && "禁售"}
            </Descriptions.Item>
            <Descriptions.Item label="累计下载">{shopInfo?.release_down_num || 0}次</Descriptions.Item>
          </Descriptions>
        </Modal>
        <div className={styles.cardList}>
          {shopPage.total > 0 && <List
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
            dataSource={[...shopList]}
            renderItem={(item: any) => {
              if (item && item?.id) {
                return (
                  <List.Item key={item?.id}>
                    <Card
                      hoverable
                      className={styles.card}
                      bodyStyle={{padding: "24px 24px 0 24px"}}
                      actions={[
                        <Tooltip title={"查看详情"}><a key="info" onClick={() => {
                          setShopInfo(item);
                          setIsInfoModalVisible(true);
                        }}><EyeOutlined/></a></Tooltip>,
                        <Tooltip title={"更新日志"}><a key="update" onClick={() => {
                          showAppUpdate(item?.id);
                        }}><CloudSyncOutlined/></a></Tooltip>,
                        <Tooltip title={"购买应用"}><a key="buy" onClick={() => {
                          if (item?.status != 0) {
                            message.error("该应用以暂停销售,详情请咨询客服")
                          } else {

                          }
                        }}><ShoppingCartOutlined/></a></Tooltip>
                      ]}>
                      <Card.Meta
                        avatar={<Tooltip title={"应用状态:" + (item?.status == 0 ? "正常" : "禁售")}><Badge
                          color={item?.status == 0 ? "green" : "red"}
                          dot={true}><SlackOutlined
                          className={styles.cardAvatar}/></Badge></Tooltip>}
                        title={<a>{item?.name}</a>}
                        description={
                          <Typography.Paragraph
                            className={styles.item}
                            ellipsis={{
                              rows: 1,
                            }}
                          >
                            <Tag color="processing">应用价格: {item?.price.toFixed(2)}元</Tag>
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
                    let tabActiveKey = "shop";
                    setTabStatus({...tabStatus, tabActiveKey});
                  }}>
                    <PlusOutlined/> 添加授权
                  </Button>
                </List.Item>
              );
            }}
          />}
          {shopPage.total <= 0 && <Empty image={Empty.PRESENTED_IMAGE_SIMPLE}/>}
        </div>
        <div hidden={shopPage.total <= 11}>
          <Pagination style={{textAlign: "center", marginTop: "20px"}} responsive={true} pageSize={10}
                      defaultCurrent={shopPage?.currentPage || 1}
                      onChange={(p) => {
                        setShopPage({currentPage: p, total: shopPage.total});
                        getShops(p);
                      }}
                      total={shopPage?.total || 0}/>
        </div>
      </>)}
      {tabStatus.tabActiveKey == "detail" && (<>
        <Modal key={"view"} width={"450px"} destroyOnClose={true} forceRender={true} footer={null} title={"授权详情"}
               visible={isViewModalVisible}
               bodyStyle={{padding: "5px"}}
               onCancel={() => {
                 setIsViewModalVisible(false);
               }}>
          <Descriptions
            bordered
            column={1}
          >
            <Descriptions.Item label="应用名称">{authInfo?.name || "未知应用"}(ID:{authInfo?.auth_id || 0})</Descriptions.Item>
            <Descriptions.Item label="应用价格">{authInfo?.price?.toFixed(2) || "0.00"}元</Descriptions.Item>
            <Descriptions.Item label="授权域名">{authInfo?.domain || "未知"}</Descriptions.Item>
            <Descriptions.Item label="许可证代码">{authInfo?.code || "未知"}</Descriptions.Item>
            <Descriptions.Item label="授权时间">{authInfo?.create_time || "未知"}</Descriptions.Item>
            <Descriptions.Item label="绑定IP">{authInfo?.ip || "未绑定"}</Descriptions.Item>
          </Descriptions>
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
                      bodyStyle={{padding: "24px 24px 0 24px"}}
                      actions={[
                        <Tooltip title={"查看详情"}><a key="view" onClick={() => {
                          setAuthInfo(item);
                          setIsViewModalVisible(true);
                        }}><EyeOutlined/></a></Tooltip>,
                        <Tooltip title={"更新日志"}><a key="update" onClick={() => {
                          showAppUpdate(item?.auth_id);
                        }}><CloudSyncOutlined/></a></Tooltip>,
                        <Tooltip title={"下载程序"}><a key="download" onClick={() => {
                          downloadApp(item?.id);
                        }}><CloudDownloadOutlined/></a></Tooltip>
                      ]}>
                      <Card.Meta
                        avatar={<Tooltip title={"授权状态:" + (item?.status == 0 ? "正常" : "禁用")}><Badge
                          color={item?.status == 0 ? "green" : "red"}
                          dot={true}><CodeSandboxOutlined
                          className={styles.cardAvatar}/></Badge></Tooltip>}
                        title={<a>{item?.name}</a>}
                        description={
                          <Typography.Paragraph
                            className={styles.item}
                            ellipsis={{
                              rows: 1,
                            }}
                          >
                            <Tag color="green">授权域名: {item?.domain}</Tag>
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
                    let tabActiveKey = "shop";
                    setTabStatus({...tabStatus, tabActiveKey});
                  }}>
                    <PlusOutlined/> 添加授权
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
      </>)}
    </PageContainer>
  );
};
