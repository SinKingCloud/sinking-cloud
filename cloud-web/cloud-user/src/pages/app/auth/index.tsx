import React, {useEffect, useState} from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import {
  Badge,
  Button,
  Card,
  Descriptions,
  List,
  message,
  Modal,
  Pagination,
  Tag,
  Tooltip,
  Typography
} from 'antd';
import styles from './style.less';
import {
  CloudDownloadOutlined, ExclamationCircleOutlined,
  EyeOutlined,
  PlusOutlined,
  SafetyCertificateOutlined,
} from "@ant-design/icons";
import {Download, getAuthList} from "@/services/app/auth";

export default (): React.ReactNode => {
  const nullData = {};
  const [appList, setAppList] = useState([]);
  const [page, setPage] = useState({currentPage: 1, total: 0});
  const [pageLoading, setPageLoading] = useState(false);
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
  useEffect(() => {
    //getApps();
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
   * 详情弹窗
   */
  const [isViewModalVisible, setIsViewModalVisible] = useState(false);
  const [authInfo, setAuthInfo] = useState<any>({});
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
            <Descriptions.Item label="程序名称">{authInfo?.name || "未知应用"}(ID:{authInfo?.auth_id || 0})</Descriptions.Item>
            <Descriptions.Item label="程序价格">{authInfo?.price?.toFixed(2) || "0.00"}元</Descriptions.Item>
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
                        <Tooltip title={"下载程序"}><a key="download" onClick={() => {
                          downloadApp(item?.id);
                        }}><CloudDownloadOutlined/></a></Tooltip>
                      ]}>
                      <Card.Meta
                        avatar={<Tooltip title={"授权状态:" + (item?.status == 0 ? "正常" : "禁用")}><Badge
                          color={item?.status == 0 ? "green" : "red"}
                          dot={true}><SafetyCertificateOutlined
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
