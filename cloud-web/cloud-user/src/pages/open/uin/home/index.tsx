import React, {useRef, useState} from 'react';
import {PageContainer, RouteContext} from '@ant-design/pro-layout';
import {Button, Col, Descriptions, Empty, message, Row, Statistic, Tag, Tooltip} from 'antd';
import styles from './style.less';
import {
  AreaChartOutlined,
  ArrowUpOutlined,
  BarChartOutlined,
  LineChartOutlined,
  StockOutlined
} from "@ant-design/icons";
import {getParams} from "@/util/page";
import ProTable from "@ant-design/pro-table";
import {getLevelList} from "@/services/open/uin/uin";
import ProCard, {StatisticCard} from "@ant-design/pro-card";
import RcResizeObserver from 'rc-resize-observer';


export default (): React.ReactNode => {
  /**
   * tab操作
   */
  const [tabStatus, seTabStatus] = useState({
    operationKey: 'buy',
    tabActiveKey: 'detail',
  });
  const onTabChange = (tabActiveKey: any) => {
    seTabStatus({...tabStatus, tabActiveKey});
  };
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
        return <div style={{maxWidth: "350px"}}>
          <Tag color="green">月:3.00￥</Tag>
          <Tag color="yellow">季:3.00￥</Tag>
          <Tag color="blue">半年:3.00￥</Tag>
          <Tag color="red">年:3.00￥</Tag>
        </div>;
      }
    },
    {
      title: '单套价格',
      dataIndex: 'price',
      tip: '挂机单套套餐下单价格',
      render: (text: string, record: any) => {
        return <div style={{maxWidth: "350px"}}>
          挂机：
          <Tag color="green">月:3.00￥</Tag>
          <Tag color="yellow">季:3.00￥</Tag>
          <Tag color="blue">半年:3.00￥</Tag>
          <Tag color="red">年:3.00￥</Tag>
          <br/>
          扩展：
          <Tag color="green">月:3.00￥</Tag>
          <Tag color="yellow">季:3.00￥</Tag>
          <Tag color="blue">半年:3.00￥</Tag>
          <Tag color="red">年:3.00￥</Tag>
        </div>;
      }
    },
    {
      title: '操作',
      dataIndex: 'Option',
      hideInSearch: true,
      editable: false,
      render: (text: string, record: any) => {
        return <Button type="default" style={{fontSize: "10px"}} size="small" onClick={() => {

        }}>
          开通
        </Button>;
      }
    },
  ];
  const description = (
    <RouteContext.Consumer>
      {({isMobile}) => (
        <Descriptions size="small" className={styles.headerList} column={isMobile ? 1 : 2}>
          <Descriptions.Item label="折扣等级">
            <Tooltip title={"每月自动根据上月消费开通对应折扣体系,支持直接升级"}>
              <a>T7<ArrowUpOutlined/></a>
            </Tooltip>
          </Descriptions.Item>
          <Descriptions.Item label="到期时间">
            <Tooltip title={"折扣等级到期时间"}>
              <a>2020-02-09 12:00:00</a>
            </Tooltip>
          </Descriptions.Item>
          <Descriptions.Item label="正常挂机">
            <Tooltip title={"正常挂机账号数量"}>
              <span style={{color: "#09ae09"}}>1000</span>
            </Tooltip>
          </Descriptions.Item>
          <Descriptions.Item label="异常挂机">
            <Tooltip title={"异常挂机账号数量"}>
              <span style={{color: "#de6300"}}>500</span>
            </Tooltip>
          </Descriptions.Item>
        </Descriptions>
      )}
    </RouteContext.Consumer>
  );
  const extra = (
    <div className={styles.moreInfo}>
      <Statistic title="挂机账号" value="200个"/>
      <Statistic title="折扣等级" value={"T7"}/>
    </div>
  );
  const [responsive, setResponsive] = useState(false);
  return (
    <PageContainer
      title={false}
      content={description}
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
                          value: 2176,
                          icon: (
                            <AreaChartOutlined style={{fontSize: "30px"}}/>
                          ),
                        }}
                      />
                      <StatisticCard
                        statistic={{
                          title: '异常挂机',
                          value: 475,
                          icon: (
                            <BarChartOutlined style={{fontSize: "30px"}}/>
                          ),
                        }}
                      />
                      <StatisticCard
                        statistic={{
                          title: '新增正常',
                          value: 87,
                          icon: (
                            <LineChartOutlined style={{fontSize: "30px"}}/>
                          ),
                        }}
                      />
                      <StatisticCard
                        statistic={{
                          title: '新增异常',
                          value: 1754,
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
                <Col xs={{span: 24, offset: 0}} lg={{span: 12, offset: 0}}>
                  <Row gutter={8} style={{paddingBottom: "20px"}}>
                    <Col xs={{span: 24, offset: 0}} lg={{span: 24, offset: 0}}>
                      <ProCard title="能力概览" bordered headerBordered={true}>
                        <Empty/>
                      </ProCard>
                    </Col>
                  </Row>
                </Col>
                <Col xs={{span: 24, offset: 0}} lg={{span: 12, offset: 0}}>
                  <Row gutter={8} style={{paddingBottom: "20px"}}>
                    <Col xs={{span: 24, offset: 0}} lg={{span: 24, offset: 0}}>
                      <ProCard title="系统公告" bordered headerBordered={true}>
                        <Empty/>
                      </ProCard>
                    </Col>
                  </Row>
                </Col>
              </Row>
            </Col>
            <Col xs={{span: 24, offset: 0}} lg={{span: 8, offset: 0}}>
              <Row gutter={8} style={{paddingBottom: "20px"}}>
                <Col xs={{span: 24, offset: 0}} lg={{span: 24, offset: 0}}>
                  <ProCard title="更新日志" bordered headerBordered={true}>
                    <Empty/>
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
