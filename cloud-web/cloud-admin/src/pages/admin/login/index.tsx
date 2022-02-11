import {LockOutlined, MobileOutlined, SafetyCertificateOutlined, UserOutlined} from '@ant-design/icons';
import {Alert, Col, Form, Image, Input, message, Row, Tabs} from 'antd';
import React, {useState} from 'react';
import {LoginForm, ProFormCaptcha, ProFormCheckbox, ProFormText} from '@ant-design/pro-form';
import {history, useModel} from 'umi';
import Footer from '@/components/Footer';
import {login} from '@/services/ant-design-pro/api';
import {getFakeCaptcha} from '@/services/ant-design-pro/login';

import styles from './index.less';

const LoginMessage: React.FC<{
  content: string;
}> = ({content}) => (
  <Alert
    style={{
      marginBottom: 24,
    }}
    message={content}
    type="error"
    showIcon
  />
);

const Login: React.FC = () => {
  const [userLoginState, setUserLoginState] = useState<API.LoginResult>({});
  const [type, setType] = useState<string>('account');
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
  const handleSubmit = async (values: API.LoginParams) => {
    try {
      // 登录
      const msg = await login({...values, type});
      if (msg.status === 'ok') {
        message.success('登录成功！');
        await fetchUserInfo();
        /** 此方法会跳转到 redirect 参数所在的位置 */
        if (!history) return;
        const {query} = history.location;
        const {redirect} = query as { redirect: string };
        history.push(redirect || '/');
        return;
      }
      console.log(msg);
      // 如果失败去设置用户错误信息
      setUserLoginState(msg);
    } catch (error) {
      message.error('登录失败，请重试！');
    }
  };
  const {status, type: loginType} = userLoginState;

  /**
   * 切换tab状态
   * @param tab 栏目
   */
  const handleTab = (tab: string) => {
    setUserLoginState({status: ""})
    setType(tab)
  };
  return (
    <div className={styles.container}>
      <div className={styles.content}>
        <LoginForm
          logo={<img alt="logo" src="/logo.svg"/>}
          title="SinKing Admin"
          subTitle="后台管理系统登陆"
          initialValues={{
            autoLogin: true,
          }}
          onFinish={async (values) => {
            await handleSubmit(values as API.LoginParams);
          }}
        >
          <Tabs activeKey={type} onChange={handleTab}>
            <Tabs.TabPane
              key="account"
              tab='账户密码登录'
            />
            <Tabs.TabPane
              key="mobile"
              tab='手机号登录'
            />
          </Tabs>

          {status === 'error' && loginType === 'account' && (
            <LoginMessage
              content='账户或密码错误'
            />
          )}
          {type === 'account' && (
            <>
              <ProFormText
                name="username"
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
          {status === 'error' && loginType === 'mobile' && <LoginMessage content="验证码错误"/>}
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
                name="captcha"
                rules={[
                  {
                    required: true,
                    message: '请输入验证码！',
                  },
                ]}
                onGetCaptcha={async (phone) => {
                  const result = await getFakeCaptcha({
                    phone,
                  });
                  if (result === false) {
                    return;
                  }
                  message.success('获取验证码成功！验证码为：1234');
                }}
              />
            </>
          )}
          <Row>
            <Col span={16}>
              <Form.Item
                name="password2"
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
                src="http://127.0.0.1:1003/verify/code"
                onClick={() => {
                  message.success("点击了验证码");
                }}
                onLoad={(e) => {
                  console.log(e)
                }}
              />
            </Col>
          </Row>
          <div style={{marginBottom: 24,}}>
            <ProFormCheckbox noStyle name="autoLogin">
              自动登录
            </ProFormCheckbox>
          </div>
        </LoginForm>
      </div>
      <Footer/>
    </div>
  );
};

export default Login;
