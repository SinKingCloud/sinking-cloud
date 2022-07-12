import type {FC} from 'react';
import {Avatar, Card, Col, Skeleton, Row, Statistic, Button, List} from 'antd';

import {Link} from 'umi';
import {PageContainer} from '@ant-design/pro-layout';
import styles from './style.less';

// @ts-ignore
import Marquee from 'react-fast-marquee';
import React, {useEffect, useState} from "react";
import {useModel} from "@@/plugin-model/useModel";
import ProCard from "@ant-design/pro-card";
import {NotificationOutlined, PlusOutlined, QqOutlined} from "@ant-design/icons";
import NoticeInfo from "@/components/Other/Notice/NoticeInfo";
import {getNoticeList} from "@/services/common/notice";
import {ago} from "@/util/time";

const PageHeaderContent: FC<{ currentUser: Partial<any> }> = ({currentUser}) => {
  const loading = currentUser && Object.keys(currentUser).length;
  if (!loading) {
    return <Skeleton avatar paragraph={{rows: 1}} active/>;
  }
  return (
    <div className={styles.pageHeaderContent}>
      <div className={styles.avatar}>
        <Avatar size="large" src={currentUser?.avatar}/>
      </div>
      <div className={styles.content}>
        <div className={styles.contentTitle}>
          你好，
          {currentUser?.nick_name}！
        </div>
        <div style={{maxWidth: "450px"}}>
          尊敬的 <b>{currentUser?.nick_name}</b> ,欢迎回来，我们已等候多时！
        </div>
      </div>
    </div>
  );
};

const ExtraContent: FC<{ currentUser: Partial<any> }> = ({currentUser}) => (
  <div className={styles.extraContent}>
    <div className={styles.statItem}>
      <Statistic title="绑定手机" value={currentUser?.phone}/>
    </div>
    <div className={styles.statItem}>
      <Statistic title="账户余额" prefix={"￥"} value={parseFloat(currentUser?.money || 0).toFixed(2)}/>
    </div>
  </div>
);

export default (): React.ReactNode => {
  const {initialState, setInitialState} = useModel('@@initialState');
  // @ts-ignore
  const {currentUser} = initialState;
  /**
   * 获取当前用户信息
   */
  const fetchUserInfo = async () => {
    const userInfo = await initialState?.fetchUserInfo?.();
    if (userInfo) {
      await setInitialState((s: any) => ({
        ...s,
        currentUser: userInfo,
      }));
    }
  };
  const [userLoading, setUserLoading] = useState(false);
  const getUserData = () => {
    setUserLoading(true);
    fetchUserInfo().then(() => {
      setUserLoading(false)
    });
  };
  /**
   * 获取公告信息
   */
  const [noticeLoading, setNoticeLoading] = useState(true);
  const [noticeData, setNoticeData] = useState<any>([]);
  const getNoticeData = async () => {
    setNoticeLoading(true);
    let d = await getNoticeList({page: 1, page_size: 7, place: "index.index"});
    if (d?.code == 200) {
      setNoticeData(d?.data?.list);
    }
    setNoticeLoading(false);
  };
  /**
   * 查看公告
   */
  const [noticeId, setNoticeId] = useState(0);
  const [noticeVisible, setNoticeVisible] = useState(false);
  const showNoticeInfoModal = (id: number) => {
    setNoticeId(id);
    setNoticeVisible(true);
  }
  const hideNoticeInfoModal = () => {
    setNoticeId(0);
    setNoticeVisible(false);
  }
  /**
   * 初始化数据
   */
  useEffect(() => {
    getNoticeData();
    getUserData();
  }, []);
  return (
    <PageContainer title={false}>
      <NoticeInfo id={noticeId} visible={noticeVisible} onClose={hideNoticeInfoModal}/>
      <Row gutter={24}>
        <Col xl={24} md={24}>

          <Card style={{marginBottom: "20px"}} loading={userLoading}>
            <Row gutter={24}>
              <Col>
                <PageHeaderContent
                  currentUser={currentUser}
                />
              </Col>
              <Col style={{margin: "10px auto 0 auto"}}>
                <ExtraContent currentUser={currentUser}/>
              </Col>
            </Row>
          </Card>
        </Col>
        <Col xl={16} lg={24} md={24} sm={24} xs={24}>
          <Row>
            <Col xl={24} lg={24} md={24} sm={24} xs={24}>
              <ProCard
                className={styles.projectList}
                headerBordered
                style={{marginBottom: 24}}
                title="我的应用"
                bordered={false}
                extra={<Link to="/app/open">查看全部</Link>}
                loading={false}
                bodyStyle={{padding: 0}}
              >
                <Card.Grid className={styles.projectGrid} style={{padding: 0}}>
                  <Card bodyStyle={{padding: 0}} bordered={false}>
                    <Card.Meta
                      style={{padding: "24px"}}
                      title={
                        <div className={styles.cardTitle}>
                          <Avatar size="small" src="https://q4.qlogo.cn/headimg_dl?dst_uin=2361935468&spec=100"/>
                          <Link to={""}>2361935468 <QqOutlined/></Link>
                        </div>
                      }
                      description={<>
                        <Row gutter={24} style={{fontSize: "10px"}}>
                          <Col xl={12} lg={12} md={12} sm={12} xs={12}>
                            挂机地区: <span style={{color: "blue"}}>上海</span>
                          </Col>
                          <Col xl={12} lg={12} md={12} sm={12} xs={12}>
                            挂机状态: <span style={{color: "green"}}>正常</span>
                          </Col>
                          <Col xl={12} lg={12} md={12} sm={12} xs={12}>
                            挂机剩余: <span style={{color: "blue"}}>30天</span>
                          </Col>
                          <Col xl={12} lg={12} md={12} sm={12} xs={12}>
                            扩展剩余: <span style={{color: "blue"}}>30天</span>
                          </Col>
                        </Row>
                      </>}
                    />
                    <div className={styles.projectItemContent}>
                      <Link to="" title={"查看详情"}>添加于2天前</Link>
                      <span className={styles.datetime} title="查看详情">
                        <Link to="">查看</Link>
                      </span>
                    </div>
                  </Card>
                </Card.Grid>

                <Card.Grid className={styles.projectGrid} style={{padding: 0}}>
                  <Card bodyStyle={{padding: 0}} bordered={false}>
                    <Card.Meta
                      style={{padding: "24px"}}
                      title={
                        <div className={styles.cardTitle}>
                          <Avatar size="small" src="https://q4.qlogo.cn/headimg_dl?dst_uin=7758521&spec=100"/>
                          <Link to={""}>7758521 <QqOutlined/></Link>
                        </div>
                      }
                      description={<>
                        <Row gutter={24} style={{fontSize: "10px"}}>
                          <Col xl={12} lg={12} md={12} sm={12} xs={12}>
                            挂机地区: <span style={{color: "blue"}}>北京</span>
                          </Col>
                          <Col xl={12} lg={12} md={12} sm={12} xs={12}>
                            挂机状态: <span style={{color: "red"}}>异常</span>
                          </Col>
                          <Col xl={12} lg={12} md={12} sm={12} xs={12}>
                            挂机剩余: <span style={{color: "blue"}}>30天</span>
                          </Col>
                          <Col xl={12} lg={12} md={12} sm={12} xs={12}>
                            扩展剩余: <span style={{color: "blue"}}>30天</span>
                          </Col>
                        </Row>
                      </>}
                    />
                    <div className={styles.projectItemContent}>
                      <Link to="" title={"查看详情"}>添加于1周前</Link>
                      <span className={styles.datetime} title="查看详情">
                        <Link to="">查看</Link>
                      </span>
                    </div>
                  </Card>
                </Card.Grid>

                <Card.Grid className={styles.projectGrid} style={{padding: 0}}>
                  <Card bodyStyle={{padding: 0}} bordered={false}>
                    <Card.Meta
                      style={{padding: "24px"}}
                      title={
                        <div className={styles.cardTitle}>
                          <Avatar size="small" src="https://q4.qlogo.cn/headimg_dl?dst_uin=2961466671&spec=100"/>
                          <Link to={""}>2961466671 <QqOutlined/></Link>
                        </div>
                      }
                      description={<>
                        <Row gutter={24} style={{fontSize: "10px"}}>
                          <Col xl={12} lg={12} md={12} sm={12} xs={12}>
                            挂机地区: <span style={{color: "blue"}}>安徽</span>
                          </Col>
                          <Col xl={12} lg={12} md={12} sm={12} xs={12}>
                            挂机状态: <span style={{color: "green"}}>正常</span>
                          </Col>
                          <Col xl={12} lg={12} md={12} sm={12} xs={12}>
                            挂机剩余: <span style={{color: "blue"}}>30天</span>
                          </Col>
                          <Col xl={12} lg={12} md={12} sm={12} xs={12}>
                            扩展剩余: <span style={{color: "blue"}}>30天</span>
                          </Col>
                        </Row>
                      </>}
                    />
                    <div className={styles.projectItemContent}>
                      <Link to="" title={"查看详情"}>添加于2个月前</Link>
                      <span className={styles.datetime} title="查看详情">
                        <Link to="">查看</Link>
                      </span>
                    </div>
                  </Card>
                </Card.Grid>

                <Card.Grid className={styles.projectGrid} style={{padding: 0}}>
                  <Card bodyStyle={{padding: 0}} bordered={false}>
                    <Card.Meta
                      style={{padding: "24px"}}
                      title={
                        <div className={styles.cardTitle}>
                          <Avatar size="small" src="https://q4.qlogo.cn/headimg_dl?dst_uin=5211314&spec=100"/>
                          <Link to={""}>5211314 <QqOutlined/></Link>
                        </div>
                      }
                      description={<>
                        <Row gutter={24} style={{fontSize: "10px"}}>
                          <Col xl={12} lg={12} md={12} sm={12} xs={12}>
                            挂机地区: <span style={{color: "blue"}}>安徽</span>
                          </Col>
                          <Col xl={12} lg={12} md={12} sm={12} xs={12}>
                            挂机状态: <span style={{color: "green"}}>正常</span>
                          </Col>
                          <Col xl={12} lg={12} md={12} sm={12} xs={12}>
                            挂机剩余: <span style={{color: "blue"}}>30天</span>
                          </Col>
                          <Col xl={12} lg={12} md={12} sm={12} xs={12}>
                            扩展剩余: <span style={{color: "blue"}}>30天</span>
                          </Col>
                        </Row>
                      </>}
                    />
                    <div className={styles.projectItemContent}>
                      <Link to="" title={"查看详情"}>添加于2个月前</Link>
                      <span className={styles.datetime} title="查看详情">
                        <Link to="">查看</Link>
                      </span>
                    </div>
                  </Card>
                </Card.Grid>
                <Card.Grid className={styles.projectGrid} style={{padding: "0px"}}>
                  <Card bodyStyle={{padding: 0}} bordered={false}>
                    <Button type="dashed" className={styles.newButton}>
                      <PlusOutlined/> 立即添加
                    </Button>
                  </Card>
                </Card.Grid>
              </ProCard>
            </Col>
          </Row>
        </Col>
        <Col xl={8} lg={24} md={24} sm={24} xs={24}>
          <ProCard
            style={{marginBottom: 24}}
            headerBordered
            title="系统公告"
            loading={noticeLoading}
            bodyStyle={{padding: "0 20px 0 20px"}}
          >
            <List
              loading={noticeLoading}
              itemLayout="horizontal"
              dataSource={noticeData}
              renderItem={(item: any) => (
                <Skeleton avatar title={false} loading={noticeLoading} active>
                  <List.Item
                    actions={[<a key="list-loadmore-edit" onClick={() => {
                      showNoticeInfoModal(item?.id || 0);
                    }}>查看</a>]}
                  >
                    <List.Item.Meta
                      key={"notice-" + item?.id}
                      avatar={<NotificationOutlined style={{fontSize: "20px", lineHeight: "50px"}}/>}
                      title={item?.title}
                      description={<span
                        style={{fontSize: "10px"}}>发布于 {ago(item?.create_time)} ,共 {item?.look_num} 次浏览</span>}
                    />
                  </List.Item>
                </Skeleton>
              )}
            />
          </ProCard>
        </Col>
      </Row>
    </PageContainer>
  );
};

