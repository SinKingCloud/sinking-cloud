import {LockOutlined, SafetyCertificateOutlined, UserOutlined} from '@ant-design/icons';
import {Alert, Avatar, Button, Checkbox, Col, Form, Image, Input, message, Row, Space, Spin, Tabs} from 'antd';
import React, {useEffect, useState} from 'react';
import Footer from '@/components/Footer';
import styles from './index.less';
// @ts-ignore
import {history, useModel} from "umi";
import {getCaptchaUrl} from "@/services/common/captcha";
import defaultSettings from "../../../../config/defaultSettings";
import {login, outLogin, ticketAuth} from "@/services/user/login";
import {setLoginToken} from "@/util/auth";
import {sendLoginSms} from "@/services/common/sms";

const Login: React.FC = () => {
  /**
   * 登陆方式
   */
  const [type, setType] = useState<string>('account');
  const [isLoading, setIsLoading] = useState(false);
  const [isSmsSendLoading, setIsSmsSendLoading] = useState(false);
  const [isRead, setIsRead] = useState(true);
  /**
   * 表单
   */
  const [form] = Form.useForm();
  /**
   * 验证码信息
   */
  const [token, setToken] = useState<string>("");
  const [captchaUrl, setCaptchaUrl] = useState<string>("");
  /**
   * 加载验证码
   */
  const changeCaptcha = () => {
    let cap = getCaptchaUrl();
    setToken(cap.token);
    setCaptchaUrl(cap.url);
    form.setFieldsValue({captcha: ""});
  }
  const {initialState, setInitialState} = useModel('@@initialState');
  /**
   * 获取当前用户信息
   */
  const fetchUserInfo = async () => {
    const userInfo = await initialState?.fetchUserInfo?.();
    if (userInfo) {
      await setInitialState((s:any) => ({
        ...s,
        currentUser: userInfo,
      }));
    }
  };

  /**
   * 退出登录，并且将当前的 url 保存
   */
  const loginOut = async () => {
    await outLogin();
  };

  const [sendCodeDisabled, setSendCodeDisabled] = useState(false);
  const getCode = (e: any) => {
    let time = 60;
    const timer = setInterval(() => {
      setSendCodeDisabled(true);
      e.target.innerHTML = `${time}秒后重新获取`;
      // eslint-disable-next-line no-plusplus
      time--;
      if (time <= 0) {
        setSendCodeDisabled(false);
        e.target.innerHTML = ' 获取验证码';
        time = 0;
        clearInterval(timer);
      }
    }, 1000);
  };

  const [redirectUrl, setRedirectUrl] = useState("/");
  const [ticket, setTicket] = useState("");
  const [loading, setLoading] = useState(true);
  /**
   * 初始化数据
   */
  useEffect(() => {
    setLoading(true);
    changeCaptcha();
    const {query} = history.location;
    const {redirect, ticket} = query as { redirect: string, ticket: string };
    setRedirectUrl(redirect);
    setTicket(ticket);
    fetchUserInfo().then(() => {
      setLoading(false);
    });
  }, []);

  /**
   * 授权登录
   */
  const ticketLogin = async () => {
    //授权ticket在跳转
    await ticketAuth({ticket: ticket, redirect_url: redirectUrl}).then((r) => {
      if (r.code == 200) {
        window.location.href = redirectUrl || "/";
      } else {
        message.error(r.message);
      }
    });
  }
  return (
    <Spin size="large" spinning={loading}>
      <div className={styles.container}>
        <div className={styles.content}>
          <div className={styles.top}>
            <div className={styles.header}>
              <img alt='logo' className={styles.logo} src="/logo.svg"/>
              <span style={{fontSize: "30px", fontWeight: "bolder"}}>{defaultSettings.title}</span>
            </div>
            <div className={styles.desc}>
              {defaultSettings.loadingTitle}
            </div>
          </div>
          <div className={styles.main}>
            {!initialState?.currentUser && (<>
              <Tabs activeKey={type} onChange={setType}>
                <Tabs.TabPane
                  key="account"
                  tab='密码登录'
                />
                <Tabs.TabPane
                  key="mobile"
                  tab='手机登录'
                />
              </Tabs>
              <Form form={form} onFinish={(values: any) => {
                if (!isRead) {
                  message.error("请阅读并同意《用户使用条款》协议")
                  return;
                }
                values.token = token;
                values.device = "web";
                if (type == "account") {
                  if (values.account == undefined || values.account.length < 5) {
                    message.error("请输入正确的账户名");
                    return;
                  }
                  if (values.password == undefined || values.password?.length <= 5) {
                    message.error("请输入正确的密码");
                    return;
                  }
                } else {
                  if (!/^[1][3,4,5,6,7,8,9][0-9]{9}$/.test(values.phone)) {
                    message.error("请输入正确的手机号码");
                    return;
                  }
                  if (values.code == undefined || values.code?.length != 6) {
                    message.error("请输入正确的短信验证码");
                    return;
                  }
                }
                if (values.captcha == undefined || values.captcha?.length != 4) {
                  message.error("请输入图形验证码");
                  return;
                }
                setIsLoading(true);
                login(values).then(async (r) => {
                  changeCaptcha();
                  if (r.code != 200) {
                    setIsLoading(false);
                    message.error(r.message || "请求失败").then()
                  } else {
                    setLoginToken("web", r.data);
                    await fetchUserInfo();
                    message.success(r.message);
                    /** 判断是否授权登陆 */
                    if (redirectUrl?.substring(0, 4) == "http") {
                      await ticketLogin();
                      setIsLoading(false);
                    } else {
                      setIsLoading(false);
                      history.push(redirectUrl || '/');
                    }
                  }
                });
              }}>
                {type === 'account' && (
                  <>
                    <Form.Item name='account'>
                      <Input prefix={<UserOutlined className='site-form-item-icon'/>} size={'large'}
                             placeholder='请输入账户'/>
                    </Form.Item>
                    <Form.Item name='password'>
                      <Input type={"password"} prefix={<LockOutlined className='site-form-item-icon'/>} size={'large'}
                             placeholder='请输入密码'/>
                    </Form.Item></>
                )}
                {type === 'mobile' && (
                  <>
                    <Form.Item name='phone'>
                      <Input prefix={<UserOutlined className='site-form-item-icon'/>} size={'large'}
                             placeholder='请输入手机号'/>
                    </Form.Item>
                    <Form.Item style={{marginBottom: "0px"}}>
                      <Row gutter={4} wrap={false}>
                        <Col flex='7'>
                          <Form.Item name='code'>
                            <Input prefix={<LockOutlined className='site-form-item-icon'/>} placeholder='短信验证码'
                                   size={'large'}/>
                          </Form.Item>
                        </Col>
                        <Col>
                          <Button size={'large'}
                                  loading={isSmsSendLoading}
                                  onClick={(e) => {
                                    let phone = form.getFieldValue("phone");
                                    if (!/^[1][3,4,5,6,7,8,9][0-9]{9}$/.test(phone)) {
                                      message.error("请输入正确的手机号码")
                                      return;
                                    }
                                    let captcha = form.getFieldValue("captcha");
                                    if (captcha == undefined || captcha.length != 4) {
                                      message.error("请输入正确的图形验证码")
                                      return;
                                    }
                                    setIsSmsSendLoading(true);
                                    sendLoginSms({phone: phone, token: token, captcha: captcha}).then((r) => {
                                      setIsSmsSendLoading(false);
                                      changeCaptcha();
                                      if (r.code != 200) {
                                        message.error(r.message || "请求失败").then()
                                      } else {
                                        getCode(e);
                                        form.setFieldsValue({captcha: ""});
                                        message.success(r.message);
                                      }
                                    });
                                  }}
                                  disabled={sendCodeDisabled}>获取验证码</Button>
                        </Col>
                      </Row>
                    </Form.Item>
                  </>
                )}
                <Row>
                  <Col span={16}>
                    <Form.Item
                      name="captcha"
                    >
                      <Input
                        prefix={<SafetyCertificateOutlined className="site-form-item-icon"/>}
                        placeholder="图形验证码"
                        type="number"
                        size="large"
                      />
                    </Form.Item>
                  </Col>
                  <Col span={7} offset={1}>
                    <Image
                      height={40}
                      preview={false}
                      src={captchaUrl}
                      onClick={() => {
                        changeCaptcha()
                      }}
                    />
                  </Col>
                </Row>
                <Form.Item hidden={type != "account"} style={{marginTop: "-10px"}}>
                  <a onClick={() => {
                    setType("mobile");
                    message.warning("使用手机号登陆自动创建账户")
                  }}>注册帐号</a>
                  <a style={{float: 'right',}} onClick={() => {
                    setType("mobile");
                    message.warning("请使用验证码登陆")
                  }}>忘记密码</a>
                </Form.Item>
                <Form.Item>
                  <Button type='primary' loading={isLoading} htmlType='submit' size={'large'} block>
                    登 录
                  </Button>
                </Form.Item>
                <Row>
                  <Form.Item name="read">
                    <Checkbox checked={isRead} onChange={(e) => {
                      setIsRead(e.target.checked);
                    }}>我已阅读并同意 <a href=''>《用户使用条款》</a>协议</Checkbox>
                  </Form.Item>
                </Row>
              </Form>
            </>)}
            {initialState?.currentUser && (<>
              <div className={styles.auth}>
                <Space direction="vertical" size="middle" style={{display: 'flex'}}>
                  <div>
                    <Avatar style={{width: "100px", height: "100px"}} src={
                      <Image
                        style={{width: "100px", height: "100px"}}
                        src={initialState?.currentUser?.avatar || "https://aliyun_id_photo_bucket.oss.aliyuncs.com/default_handsome.jpg"}/>
                    }/>
                    <Alert message={initialState?.currentUser?.nick_name}
                           className={styles.nick} type="info"/>
                  </div>
                  <Button type='primary' size={'large'} block onClick={async () => {
                    /** 判断是否授权登陆 */
                    if (redirectUrl?.substring(0, 4) == "http") {
                      await ticketLogin();
                      setIsLoading(false);
                    } else {
                      setIsLoading(false);
                      history.push(redirectUrl || '/');
                    }
                  }}>
                    快捷登陆
                  </Button>
                  <Button type='default' size={'large'} block onClick={() => {
                    setInitialState((s: any) => ({...s, currentUser: undefined})).then(() => {
                      loginOut().then(() => {
                        message.success("您已注销登录");
                      });
                    });
                  }}>
                    注销登录
                  </Button>
                </Space>
              </div>
            </>)}
          </div>
        </div>
        <Footer/>
      </div>
    </Spin>
  );
};

export default Login;
