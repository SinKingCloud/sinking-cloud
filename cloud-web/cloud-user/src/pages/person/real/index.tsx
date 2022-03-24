// @ts-ignore
import React, {useEffect, useState} from 'react';
import {Alert, Button, Col, Descriptions, Form, Input, message, Result, Row, Select, Space, Steps} from 'antd';
import {PageContainer} from "@ant-design/pro-layout";
import ProCard from "@ant-design/pro-card";
import {authInit, checkAuth, getAuthInfo} from "@/services/person/auth";
import {checkMobile} from "@/util/device";
// @ts-ignore
import QRCode from "qrcode.react/lib/index";
import {history} from "umi";
import {useModel} from "@@/plugin-model/useModel";


export default (): React.ReactNode => {
  const [current, setCurrent] = useState(0);
  const [isLoading, setIsLoading] = useState(true)
  const [isBtnLoading, setIsBtnLoading] = useState(false)
  const [authInfo, setAuthInfo] = useState({card_no: "", card_name: ""});
  const [authData, setAuthData] = useState({id: "", url: ""});
  const [formHidden, setFormHidden] = useState(false);
  const [qrcodeHidden, setQrcodeHidden] = useState(true);
  const {initialState, setInitialState} = useModel('@@initialState');
  /**
   * 刷新当前用户信息
   */
  const fetchUserInfo = async () => {
    const userInfo = await initialState?.fetchUserInfo?.();
    if (userInfo) {
      await setInitialState((s) => ({
        ...s,
        currentUser: userInfo,
      }));
    }
  };

  /**
   * 初始化信息
   */
  const initAuthInfo = () => {
    getAuthInfo().then((r) => {
      if (r.code != 200) {
        message.error(r.message || "请求失败").then()
      } else {
        setAuthInfo(r.data)
        if (r.data?.real_auth == 1) {
          setCurrent(2)
        }
      }
      setIsLoading(false)
    });
  }
  /**
   * 显示实名信息
   */
  useEffect(() => {
    initAuthInfo();
  }, []);

  /**
   * 递归判断状态
   */
  const checkStatus = (id: string) => {
    if (id != "" && localStorage.getItem('internal') == id && location.pathname == "/person/real" && current == 1) {
      checkAuth({certifyId: id}).then(async (r) => {
        if (r.code != 200) {
          setTimeout(() => {
            checkStatus(id);
          }, 500)
        } else {
          localStorage.removeItem('internal');
          initAuthInfo();
          await fetchUserInfo();
        }
      });
    } else {
      localStorage.removeItem('internal');
    }
  }

  /**
   * 步骤控制
   */
  const next = () => {
    setCurrent(current + 1);
  };
  const prev = () => {
    setCurrent(current - 1);
  };
  return (
    <PageContainer title={false}>
      <ProCard loading={isLoading}>
        <Steps current={current} direction="horizontal">
          <Steps.Step key="0" title="隐私说明"/>
          <Steps.Step key="1" title="验证身份"/>
          <Steps.Step key="2" title="验证结果"/>
        </Steps>
        <div style={{marginTop: "30px"}}>
          <Row>
            <Col xs={{span: 24, offset: 0}} lg={{span: 10, offset: 7}}>
              {current == 0 && (
                <>
                  <Alert message={"本站采用阿里云金融级安全验证,保证您的信息不被泄露"} type="info" showIcon/>
                  <Space align={"center"} direction={"vertical"} style={{marginTop: "20px"}}>
                    <Descriptions title="隐私及安全说明" column={1}>
                      <Descriptions.Item>
                        根据中华人民共和国《网络安全法》,本站将验证您的真人信息。请及时完成实名验证以使用完整的服务
                      </Descriptions.Item>
                      <Descriptions.Item>
                        保护用户隐私是本站的重点原则，本站通过技术手段、提供隐私保护服务功能、强化内部管理等办法充分保护用户的个人资料安全。
                      </Descriptions.Item>
                      <Descriptions.Item>
                        本站严格保护您个人信息的安全。我们使用各种安全技术和程序来保护您的个人信息不被未经授权的访问、使用或泄漏。
                      </Descriptions.Item>
                      <Descriptions.Item>
                        本站在法律要求或符合本站的相关服务条款、软件许可使用协议约定的情况下透露您的个人信息，或者有充分理由相信必须这样做才能：
                        <br/>
                        (a) 满足法律或行政法规的明文规定，或者符合本站适用的法律程序；
                        <br/>
                        (b) 符合本真相关服务条款、软件许可使用协议的约定；
                        <br/>
                        (c) 保护本站权利或财产，以及在紧急情况下保护本站员工、本站产品或服务的用户或大众的个人安全。
                      </Descriptions.Item>
                      <Descriptions.Item>
                        本站不会未经您的允许将这些信息与第三方共享，本声明已经列出的上述情况除外。
                      </Descriptions.Item>
                    </Descriptions>
                  </Space>
                  <div style={{textAlign: "center", marginTop: "20px"}}>
                    <Button type="primary" onClick={() => next()}>
                      我已知晓，下一步
                    </Button>
                  </div>
                </>
              )}
              {current == 1 && (
                <>
                  <Alert message={"请填写您的身份信息，确认提交后将会进入支付宝进行人脸核验"} type="info" showIcon/>
                  <Row>
                    <Col xs={{span: 24, offset: 0}} lg={{span: 16, offset: 4}}>
                      <div style={{marginTop: "20px"}}>
                        <div hidden={qrcodeHidden} style={{textAlign: "center"}}>
                          <QRCode
                            style={{margin: "0px auto"}}
                            value={authData?.url}
                            size={150}
                            fgColor="#000000"
                          />
                          <div style={{marginTop: "10px", color: "#8c8c8c"}}>
                            请打开支付宝APP扫码完成验证
                          </div>
                          <div style={{marginTop: "30px"}}>
                            <Button type="primary" onClick={() => {
                              localStorage.removeItem('internal');
                              setFormHidden(false)
                              setQrcodeHidden(true);
                            }}>
                              重新提交
                            </Button>
                          </div>
                        </div>
                        <Form hidden={formHidden} layout={"vertical"} onFinish={(values) => {
                          setIsBtnLoading(true)
                          authInit(values).then((r) => {
                            setIsBtnLoading(false)
                            if (r.code != 200) {
                              message.error(r.message || "请求失败").then()
                            } else {
                              setAuthData(r.data);
                              localStorage.setItem('internal', r.data.id);
                              checkStatus(r.data.id);
                              //显示二维码
                              setFormHidden(true)
                              setQrcodeHidden(false);
                              //移动端设备唤醒支付宝
                              if (checkMobile()) {
                                window.open(r.data.url);
                              }
                            }
                          })
                        }}>
                          <Form.Item name="card_type" label="证件类型" initialValue={"IDENTITY_CARD"}
                                     rules={[{required: true}]}>
                            <Select>
                              <Select.Option value="IDENTITY_CARD">居民身份证</Select.Option>
                              <Select.Option value="HOME_VISIT_PERMIT_HK_MC">港澳通行证</Select.Option>
                              <Select.Option value="HOME_VISIT_PERMIT_TAIWAN">台湾通行证</Select.Option>
                              <Select.Option value="RESIDENCE_PERMIT_HK_MC">港澳居住证</Select.Option>
                              <Select.Option value="RESIDENCE_PERMIT_TAIWAN">台湾居住证</Select.Option>
                            </Select>
                          </Form.Item>
                          <Form.Item name="card_name" label="身份证姓名" rules={[{required: true}]}>
                            <Input placeholder="请输入您的身份证姓名"/>
                          </Form.Item>
                          <Form.Item name="card_no" label="身份证号码" rules={[{required: true}]}>
                            <Input placeholder="请输入您的身份证号码"/>
                          </Form.Item>
                          <Form.Item style={{textAlign: "center",}}>
                            <Button style={{margin: '0 8px'}} onClick={() => prev()}>
                              上一步
                            </Button>
                            <Button type="primary" loading={isBtnLoading} htmlType="submit">
                              提交认证
                            </Button>
                          </Form.Item>
                        </Form>
                      </div>
                    </Col>
                  </Row>
                </>
              )}
              {current == 2 && (
                <>
                  <Result
                    status="success"
                    title="您已完成实名认证"
                    subTitle="您现在可以使用本站提供的全部服务"
                    extra={[
                      <Button type="primary" key="console" onClick={() => {
                        history.goBack();
                      }}>
                        返回上一页
                      </Button>,
                    ]}
                  />
                  <Space align={"center"} direction={"vertical"}
                         style={{padding: "10px", borderRadius: "10px", backgroundColor: "#e7e7e7"}}>
                    <Descriptions title="实名信息" column={1}>
                      <Descriptions.Item label="身份证姓名">
                        {authInfo?.card_name}
                      </Descriptions.Item>
                      <Descriptions.Item label="身份证号码">
                        {authInfo?.card_no}
                      </Descriptions.Item>
                    </Descriptions>
                  </Space>
                </>
              )}
            </Col>
          </Row>
        </div>
      </ProCard>
    </PageContainer>
  );
};
