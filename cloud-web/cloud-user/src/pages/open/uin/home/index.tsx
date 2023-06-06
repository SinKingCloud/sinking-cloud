import React, {useEffect, useRef, useState} from 'react';
import {PageContainer, RouteContext} from '@ant-design/pro-layout';
import {
  Button,
  Col,
  Descriptions,
  List,
  message,
  Modal,
  Popover,
  Row,
  Skeleton,
  Statistic, Table,
  Tag,
  Tooltip
} from 'antd';
import styles from './style.less';
import {
  AreaChartOutlined,
  ArrowUpOutlined,
  BarChartOutlined,
  ExclamationCircleOutlined,
  LineChartOutlined, NotificationOutlined, QuestionCircleOutlined,
  StockOutlined
} from "@ant-design/icons";
import {getParams} from "@/util/page";
import ProTable from "@ant-design/pro-table";
import {getLevelList} from "@/services/open/uin/uin";
import ProCard, {StatisticCard} from "@ant-design/pro-card";
import RcResizeObserver from 'rc-resize-observer';
import {getSelfInfo, getUinCount, shopLevel} from "@/services/open/uin/home";
import {getNoticeList} from "@/services/common/notice";
import {ago} from "@/util/time";
import NoticeInfo from "@/components/Other/Notice/NoticeInfo";
import ProDescriptions from '@ant-design/pro-descriptions';


export default (): React.ReactNode => {
  /**
   * tab操作
   */
  const [tabStatus, setTabStatus] = useState({
    operationKey: 'buy',
    tabActiveKey: 'detail',
  });
  const onTabChange = (tabActiveKey: any) => {
    setTabStatus({...tabStatus, tabActiveKey});
  };
  /**
   * 获取公告信息
   */
  const [noticeLoading, setNoticeLoading] = useState(true);
  const [noticeData, setNoticeData] = useState<any>([]);
  const getNoticeData = async () => {
    setNoticeLoading(true);
    const d = await getNoticeList({page: 1, page_size: 7, place: 'open.uin.home'});
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
  const [pageLoading, setPageLoading] = useState(true);
  const [selfInfo, setSelfInfo] = useState<any>({});
  const [uinCount, setUinCount] = useState<any>({});
  const init = async () => {
    setPageLoading(true);
    getNoticeData();
    await getSelfInfo().then((r) => {
      if (r.code == 200) {
        setSelfInfo(r.data);
      }
    });
    await getUinCount().then((r) => {
      if (r.code == 200) {
        setUinCount(r.data);
      }
    });
    setPageLoading(false);
  }
  useEffect(() => {
    init().then();
  }, []);
  const ref = useRef();
  /**
   * table表格渲染
   */
  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      tip: '等级ID',
      hideInSearch: true,
      hideInTable: true,
    },
    {
      title: '等级名称',
      dataIndex: 'name',
      tip: '折扣等级名称',
      render: (text: string, record: any) => {
        return <Tooltip title={record?.describe}>{record?.name}</Tooltip>;
      }
    },
    {
      title: '使用门槛',
      dataIndex: 'price',
      tip: '满足其中一个等级门槛(满足每月消费金额/直接按月开通等级)',
      render: (text: string, record: any) => {
        return <>
          月消:{record?.month_spend.toFixed(2)}元<br/>
          直购:{record?.price.toFixed(2)}元
        </>;
      }
    },
    {
      title: '全套价格',
      dataIndex: 'price',
      tip: '挂机全套套餐下单价格',
      render: (text: string, record: any) => {
        return <div style={{maxWidth: "250px"}}>
          <Tag color="green">月:{record?.month_price.toFixed(2) || "0.00"}￥</Tag>
          <Tag color="yellow">季:{record?.quarter_price.toFixed(2) || "0.00"}￥</Tag>
          <Tag color="blue">半年:{record?.half_price.toFixed(2) || "0.00"}￥</Tag>
          <br/>
          <Tag color="red">年:{record?.year_price.toFixed(2) || "0.00"}￥</Tag>
          <Tag color="purple">永久:{record?.forever_price.toFixed(2) || "0.00"}￥</Tag>
        </div>;
      }
    },
    {
      title: '单套价格',
      dataIndex: 'price',
      tip: '挂机单套套餐下单价格',
      render: (text: string, record: any) => {
        return <div style={{maxWidth: "450px"}}>
          挂机：
          <Tag color="green">月:{record?.player_month_price.toFixed(2) || "0.00"}￥</Tag>
          <Tag color="yellow">季:{record?.player_quarter_price.toFixed(2) || "0.00"}￥</Tag>
          <Tag color="blue">半年:{record?.player_half_price.toFixed(2) || "0.00"}￥</Tag>
          <Tag color="red">年:{record?.player_year_price.toFixed(2) || "0.00"}￥</Tag>
          <Tag color="purple">永久:{record?.player_forever_price.toFixed(2) || "0.00"}￥</Tag>
          <br/>
          扩展：
          <Tag color="green">月:{record?.sign_month_price.toFixed(2) || "0.00"}￥</Tag>
          <Tag color="yellow">季:{record?.sign_quarter_price.toFixed(2) || "0.00"}￥</Tag>
          <Tag color="blue">半年:{record?.sign_half_price.toFixed(2) || "0.00"}￥</Tag>
          <Tag color="red">年:{record?.sign_year_price.toFixed(2) || "0.00"}￥</Tag>
          <Tag color="purple">永久:{record?.sign_forever_price.toFixed(2) || "0.00"}￥</Tag>
        </div>;
      }
    },
    {
      title: '操作',
      dataIndex: 'Option',
      hideInSearch: true,
      editable: false,
      render: (text: string, record: any) => {
        // @ts-ignore
        const btn = selfInfo?.level.id == record?.id
        return <Button type={btn ? "primary" : "default"} style={{fontSize: "10px"}} size="small" onClick={() => {
          if (record?.price <= 0) {
            message.warning("免费等级无需开通,系统自动结转");
            return;
          }
          Modal.confirm({
            title: '确定要花费' + record?.price.toFixed(2) + '元开通1个月' + record?.name + '等级吗?',
            icon: <ExclamationCircleOutlined/>,
            content: '开通不同等级则会清空已开通等级,同等级则顺延1月到期时间',
            okType: 'primary',
            onOk() {
              const key = "shop";
              message.loading({content: '正在开通' + record?.name + '等级', key, duration: 60}).then();
              shopLevel({level_id: record?.id || 0, month: 1}).then((r) => {
                if (r.code == 200) {
                  init().then(() => {
                    message.success({content: '开通成功', key, duration: 2}).then();
                  });
                } else {
                  message.error({content: r.message, key, duration: 2}).then();
                }
              });
            },
          });
        }}>
          {btn ? "续期" : "开通"}
        </Button>;
      }
    },
  ];

  const description = (
    <Skeleton active loading={pageLoading}>
      <RouteContext.Consumer>
        {({isMobile}) => (
          <Descriptions size="small" className={styles.headerList} column={isMobile ? 1 : 2}>
            <Descriptions.Item label="折扣等级">
              <Tooltip title={"每月自动根据上月消费开通对应折扣体系,支持直接升级"}>
                <a>{ // @ts-ignore
                  selfInfo?.level?.name || "默认"}<ArrowUpOutlined/></a>
              </Tooltip>
            </Descriptions.Item>
            <Descriptions.Item label="到期时间">
              <Tooltip title={"折扣等级到期时间"}>
                <a>{ // @ts-ignore
                  selfInfo?.expire || "长期"}</a>
              </Tooltip>
            </Descriptions.Item>
            <Descriptions.Item label="正常挂机">
              <Tooltip title={"正常挂机账号数量"}>
                <span style={{color: "#09ae09"}}>{ // @ts-ignore
                  uinCount?.uin_num_all_normal || 0}</span>
              </Tooltip>
            </Descriptions.Item>
            <Descriptions.Item label="异常挂机">
              <Tooltip title={"异常挂机账号数量"}>
                <span style={{color: "#de6300"}}>{ // @ts-ignore
                  uinCount?.uin_num_all_un_normal || 0}</span>
              </Tooltip>
            </Descriptions.Item>
          </Descriptions>
        )}
      </RouteContext.Consumer>
    </Skeleton>
  );
  const extra = (
    <Skeleton active loading={pageLoading}>
      <div className={styles.moreInfo}>
        <Statistic title="挂机账号" value={// @ts-ignore
          ((uinCount?.uin_num_all_normal || 0) + (uinCount?.uin_num_all_un_normal || 0)) + "个"}/>
        <Statistic title="折扣等级" value={ // @ts-ignore
          selfInfo?.level?.name || "默认"}/>
      </div>
    </Skeleton>
  );
  const [responsive, setResponsive] = useState(false);

  return (
    <PageContainer
      title={false}
      content={description}
      loading={pageLoading}
      extraContent={extra}
      className={styles.pageHeader}
      tabActiveKey={tabStatus.tabActiveKey}
      onTabChange={onTabChange}
      tabList={[
        {
          key: 'detail',
          tab: '能力概览',
        },
        {
          key: 'buy',
          tab: '折扣等级',
        },
      ]}
    >
      <div className={styles.main}>
        <NoticeInfo id={noticeId} visible={noticeVisible} onClose={hideNoticeInfoModal}/>
        {tabStatus.tabActiveKey == "detail" && (
          <Row gutter={8} style={{paddingBottom: "20px"}}>
            <Col xs={{span: 24, offset: 0}} lg={{span: 16, offset: 0}}>
              <Row gutter={8} style={{paddingBottom: "20px"}}>
                <Col xs={{span: 24, offset: 0}} lg={{span: 24, offset: 0}}>
                  <RcResizeObserver
                    key="resize-observer"
                    onResize={(offset: any) => {
                      setResponsive(offset.width < 596);
                    }}
                  >
                    <StatisticCard.Group direction={responsive ? 'column' : 'row'}>
                      <StatisticCard
                        statistic={{
                          title: '正常挂机',
                          // @ts-ignore
                          value: uinCount?.uin_num_all_normal || 0,
                          icon: (
                            <AreaChartOutlined style={{fontSize: "30px"}}/>
                          ),
                        }}
                      />
                      <StatisticCard
                        statistic={{
                          title: '异常挂机',
                          // @ts-ignore
                          value: uinCount?.uin_num_all_un_normal || 0,
                          icon: (
                            <BarChartOutlined style={{fontSize: "30px"}}/>
                          ),
                        }}
                      />
                      <StatisticCard
                        statistic={{
                          title: '新增正常',
                          // @ts-ignore
                          value: uinCount?.uin_num_today_normal || 0,
                          icon: (
                            <LineChartOutlined style={{fontSize: "30px"}}/>
                          ),
                        }}
                      />
                      <StatisticCard
                        statistic={{
                          title: '新增异常',
                          // @ts-ignore
                          value: uinCount?.uin_num_today_un_normal || 0,
                          icon: (
                            <StockOutlined style={{fontSize: "30px"}}/>
                          ),
                        }}
                      />
                    </StatisticCard.Group>
                  </RcResizeObserver>
                </Col>
              </Row>
              <Row gutter={8}>
                <Col xs={{span: 24, offset: 0}} lg={{span: 14, offset: 0}}>
                  <Row gutter={8} style={{paddingBottom: "20px"}}>
                    <Col xs={{span: 24, offset: 0}} lg={{span: 24, offset: 0}}>
                      <ProCard title="产品价格" bordered headerBordered={true}>
                        <ProDescriptions
                          style={{marginTop: "-10px", marginBottom: "-10px"}}
                          loading={pageLoading}
                          column={{xs: 1, sm: 1, md: 1}}
                          labelStyle={{fontSize: "12px", lineHeight: "25px"}}
                          contentStyle={{fontSize: "12px", lineHeight: "25px"}}
                        >
                          <ProDescriptions.Item label="账户余额">{selfInfo?.money + '元' || "获取失败"}</ProDescriptions.Item>
                          <ProDescriptions.Item label="折扣等级">{selfInfo?.level?.name || "获取失败"}</ProDescriptions.Item>
                          <ProDescriptions.Item label="到期时间">
                            <Popover title={"折扣说明"}
                                     content={<div style={{fontSize: "10px"}}>折扣等级会每月初根据您站点月销<br/>情况自动评估,也可以直接购买<br/>提升等级,详情请访问官网
                                     </div>}>
                              {selfInfo?.expire || "获取失败"} &nbsp;<QuestionCircleOutlined/>
                            </Popover>
                          </ProDescriptions.Item>
                          <ProDescriptions.Item>
                            <Table style={{maxWidth: "95%"}} dataSource={[
                              {
                                key: '1',
                                name: '全套',
                                month: selfInfo?.level?.month_price || 0,
                                quarter: selfInfo?.level?.quarter_price || 0,
                                half: selfInfo?.level?.half_price || 0,
                                year: selfInfo?.level?.year_price || 0,
                                forever: selfInfo?.level?.forever_price || 0,
                              },
                              {
                                key: '2',
                                name: '加速',
                                age: 42,
                                month: selfInfo?.level?.player_month_price || 0,
                                quarter: selfInfo?.level?.player_quarter_price || 0,
                                half: selfInfo?.level?.player_half_price || 0,
                                year: selfInfo?.level?.player_year_price || 0,
                                forever: selfInfo?.level?.player_forever_price || 0,
                              },
                              {
                                key: '3',
                                name: '扩展',
                                month: selfInfo?.level?.sign_month_price || 0,
                                quarter: selfInfo?.level?.sign_quarter_price || 0,
                                half: selfInfo?.level?.sign_half_price || 0,
                                year: selfInfo?.level?.sign_year_price || 0,
                                forever: selfInfo?.level?.sign_forever_price || 0,
                              },
                            ]} columns={[
                              {
                                title: '名称',
                                dataIndex: 'name',
                                key: 'name',
                                fixed: 'left',
                                width: 50,
                              },
                              {
                                title: '月度',
                                dataIndex: 'month',
                                key: 'month',
                                width: 50,
                              },
                              {
                                title: '季度',
                                dataIndex: 'quarter',
                                key: 'quarter',
                                width: 50,
                              },
                              {
                                title: '半年',
                                dataIndex: 'half',
                                key: 'half',
                                width: 50,
                              },
                              {
                                title: '一年',
                                dataIndex: 'year',
                                key: 'year',
                                width: 50,
                              },
                              {
                                title: '永久',
                                dataIndex: 'forever',
                                key: 'forever',
                                width: 50,
                              },
                            ]} pagination={false} bordered={true} scroll={{y: 500}} size={"small"}/>
                          </ProDescriptions.Item>
                        </ProDescriptions>
                      </ProCard>

                    </Col>
                  </Row>
                </Col>
                <Col xs={{span: 24, offset: 0}} lg={{span: 10, offset: 0}}>
                  <Row gutter={8} style={{paddingBottom: "20px"}}>
                    <Col xs={{span: 24, offset: 0}} lg={{span: 24, offset: 0}}>
                      <ProCard title="产品权限" bordered headerBordered={true}>
                        <ProDescriptions
                          style={{marginTop: "-10px", marginBottom: "-15px"}}
                          loading={pageLoading}
                          column={{xs: 1, sm: 1, md: 1}}
                          labelStyle={{fontSize: "12px", lineHeight: "25px"}}
                          contentStyle={{justifyContent: "flex-end", fontSize: "12px", lineHeight: "25px"}}
                        >
                          <ProDescriptions.Item label="主动扫码"
                          >
                            <Button type="primary" ghost size="small"
                                    style={{fontSize: "10px", height: "20px", padding: "0px 5px"}}>开通</Button>

                          </ProDescriptions.Item>
                          <ProDescriptions.Item label="Cookie提取">
                            <Button type="primary" ghost size="small"
                                    style={{fontSize: "10px", height: "20px", padding: "0px 5px"}}>开通</Button>
                          </ProDescriptions.Item>
                          <ProDescriptions.Item label="授权登录">
                            <Button type="primary" ghost size="small"
                                    style={{fontSize: "10px", height: "20px", padding: "0px 5px"}}>开通</Button>
                          </ProDescriptions.Item>
                        </ProDescriptions>
                      </ProCard>
                    </Col>
                  </Row>
                </Col>
              </Row>
            </Col>
            <Col xs={{span: 24, offset: 0}} lg={{span: 8, offset: 0}}>
              <Row gutter={8} style={{paddingBottom: "20px"}}>
                <Col xs={{span: 24, offset: 0}} lg={{span: 24, offset: 0}}>
                  <ProCard title="系统公告" bordered headerBordered={true}
                           loading={noticeLoading}
                           bodyStyle={{padding: "0 20px 0 20px"}}>
                    <List
                      loading={noticeLoading}
                      itemLayout="horizontal"
                      dataSource={noticeData}
                      renderItem={(item: any) => (
                        <Skeleton avatar title={false} loading={noticeLoading} key={item?.id} active>
                          <List.Item
                            actions={[<a key="list-loadmore-edit" onClick={() => {
                              showNoticeInfoModal(item?.id || 0);
                            }}>查看</a>]}
                            key={item?.id}
                          >
                            <List.Item.Meta
                              key={"notice-" + item?.id}
                              avatar={<NotificationOutlined style={{fontSize: "15px", lineHeight: "50px"}}/>}
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
            </Col>
          </Row>
        )}
        {tabStatus.tabActiveKey == "buy" && (
          <Row gutter={8}>
            <Col xs={{span: 24, offset: 0}} lg={{span: 24, offset: 0}}>
              <ProTable
                actionRef={ref}
                defaultSize={"small"}
                form={{layout: "vertical"}}
                headerTitle={'等级列表'}
                rowKey={'id'}
                options={{
                  density: true,
                  fullScreen: true,
                  setting: true,
                }}
                // @ts-ignore
                columns={columns}
                request={async (params, sort) => {
                  params.order_type = 0;
                  const data = await getLevelList(getParams(params, sort));
                  if (data.code != 200) {
                    message.error(data.message);
                  }
                  return {
                    data: data.data.list === undefined || data.data.list === null || data.data.list.length <= 0 ? [] : data.data.list,
                    success: data.code === 200,
                    total: data.data.total,
                  };
                }}
                pagination={{defaultPageSize: 100}}
                search={false}
              />
            </Col>
          </Row>
        )}
      </div>
    </PageContainer>
  );
};
