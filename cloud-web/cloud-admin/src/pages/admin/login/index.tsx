import {LockOutlined, MobileOutlined, SafetyCertificateOutlined, UserOutlined} from '@ant-design/icons';
import {Col, Form, Image, Input, message, Row, Tabs} from 'antd';
import React, {useEffect, useState} from 'react';
import {LoginForm, ProFormCaptcha, ProFormText} from '@ant-design/pro-form';
import Footer from '@/components/Footer';
import {login} from '@/services/admin/login';
import defaultSettings from '@/../config/defaultSettings';
import styles from './index.less';
import {history, useModel} from "umi";
import {setLoginToken} from "@/util/auth";
import {getCaptchaUrl} from "@/services/common/captcha";

const Login: React.FC = () => {
  /**
   * 登陆方式
   */
  const [type, setType] = useState<string>('account');
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
    setToken(cap.token)
    setCaptchaUrl(cap.url)
  }
  useEffect(changeCaptcha, [])
  const {initialState, setInitialState} = useModel('@@initialState');
  /**
   * 获取当前用户信息
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
   * 提交登录
   * @param values 登陆表单
   */
  const handleSubmit = async (values: AdminAPI.LoginParams) => {
    try {
      values.token = token
      values.device = "web"
      // 登录
      const response = await login(values);
      if (response.code == 200) {
        message.success(response.message)
        setLoginToken("web", response.data)
        await fetchUserInfo();
        /** 此方法会跳转到 redirect 参数所在的位置 */
        if (!history) return;
        const {query} = history.location;
        const {redirect} = query as { redirect: string };
        history.push(redirect || '/');
        return;
      } else {
        message.error(response.message)
        changeCaptcha()
      }
    } catch (error) {
      console.log(error)
      message.error('登录失败，请重试！');
    }
  };
  // @ts-ignore
  return (
    <div className={styles.container}>
      <div className={styles.content}>
        <LoginForm
          logo={<img alt="logo" src="/logo.svg"/>}
          title={defaultSettings.title}
          subTitle="后台管理系统"
          onFinish={async (values) => {
            await handleSubmit(values as AdminAPI.LoginParams);
          }}
        >
          <Tabs activeKey={type} onChange={setType} hidden={true}>
            <Tabs.TabPane
              key="account"
              tab='账户密码登录'
            />
            <Tabs.TabPane
              key="mobile"
              tab='手机号登录'
            />
          </Tabs>
          {type === 'account' && (
            <>
              <ProFormText
                name="account"
                fieldProps={{
                  size: 'large',
                  prefix: <UserOutlined className={styles.prefixIcon}/>,
                }}
                placeholder='用户名'
                rules={[
                  {
                    required: true,
                    message: '请输入用户名!',
                  },
                ]}
              />
              <ProFormText.Password
                name="password"
                fieldProps={{
                  size: 'large',
                  prefix: <LockOutlined className={styles.prefixIcon}/>,
                }}
                placeholder='密码'
                rules={[
                  {
                    required: true,
                    message: '请输入密码！',
                  },
                ]}
              />
            </>
          )}
          {type === 'mobile' && (
            <>
              <ProFormText
                fieldProps={{
                  size: 'large',
                  prefix: <MobileOutlined className={styles.prefixIcon}/>,
                }}
                name="mobile"
                placeholder='手机号'
                rules={[
                  {
                    required: true,
                    message: '请输入手机号！',
                  },
                  {
                    pattern: /^1\d{10}$/,
                    message: '手机号格式错误！',
                  },
                ]}
              />
              <ProFormCaptcha
                fieldProps={{
                  size: 'large',
                  prefix: <LockOutlined className={styles.prefixIcon}/>,
                }}
                captchaProps={{
                  size: 'large',
                }}
                placeholder='请输入验证码'
                captchaTextRender={(timing, count) => {
                  if (timing) {
                    return `${count} 获取验证码`;
                  }
                  return '获取验证码';
                }}
                name="code"
                rules={[
                  {
                    required: true,
                    message: '请输入验证码！',
                  },
                ]}
                onGetCaptcha={async (phone) => {
                  message.success(phone + '获取验证码成功！验证码为：1234');
                }}
              />
            </>
          )}
          <Row>
            <Col span={16}>
              <Form.Item
                name="captcha"
                rules={[{required: true, message: '请输入图片验证码!'}]}
              >
                <Input
                  prefix={<SafetyCertificateOutlined className="site-form-item-icon"/>}
                  placeholder="图形验证码"
                  type="number"
                  size="large"
                />
              </Form.Item>
            </Col>
            <Col span={6} offset={1}>
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
        </LoginForm>
      </div>
      <Footer/>
    </div>
  );
};

export default Login;
