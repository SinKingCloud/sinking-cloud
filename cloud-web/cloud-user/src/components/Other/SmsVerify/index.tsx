import React, {useEffect, useState} from 'react';
import {Button, Col, Form, Image, Input, message, Row} from "antd";
import {FormInstance} from "antd/lib/form/hooks/useForm";
import {getCaptchaUrl} from "@/services/common/captcha";
import {sendSms} from "@/services/common/sms";

export type SmsVerifyProps = {
  onFinish: (value?: any) => Promise<boolean | void>;
  form: FormInstance<any>;
  phone?: string;
  topNodes?: React.ReactNode;
  bottomNodes?: React.ReactNode;
};

const SmsVerify: React.FC<SmsVerifyProps> = (props) => {
  const {onFinish, form, phone, bottomNodes, topNodes} = props;
  /**
   * 验证码信息
   */
  const [token, setToken] = useState<string>("");
  const [captchaUrl, setCaptchaUrl] = useState<string>("");
  const [isCaptchaHidden, setIsCaptchaHidden] = useState<boolean>(false);
  /**
   * 加载验证码
   */
  const changeCaptcha = () => {
    let cap = getCaptchaUrl();
    setToken(cap.token);
    setCaptchaUrl(cap.url);
  }

  const [sendCodeDisabled, setSendCodeDisabled] = useState(false);
  const [isSmsSendLoading, setIsSmsSendLoading] = useState(false);
  const getCode = (e: any) => {
    let time = 60;
    const timer = setInterval(() => {
      setIsCaptchaHidden(true)
      setSendCodeDisabled(true);
      e.target.innerHTML = `${time}秒后重新获取`;
      // eslint-disable-next-line no-plusplus
      time--;
      if (time <= 0) {
        getCaptchaUrl();
        setIsCaptchaHidden(false);
        setSendCodeDisabled(false);
        e.target.innerHTML = ' 获取验证码';
        time = 0;
        clearInterval(timer);
      }
    }, 1000);
  };

  useEffect(() => {
    changeCaptcha();
  }, []);
  return (
    <>
      <Form form={form} name="control-hooks" onFinish={async (values: any) => {
        values.token = token;
        await onFinish(values);
      }} labelAlign="right"
            labelCol={{span: 6}}>
        {topNodes}
        <Form.Item name={"code"} label="短信验证码" rules={[{required: true, message: '请输入短信验证码',},]}>
          <Row gutter={6} wrap={false}>
            <Col flex={22}>
              <Input placeholder="请输入短信验证码"/>
            </Col>
            <Col flex={2}>
              <Button disabled={sendCodeDisabled} loading={isSmsSendLoading} onClick={(e) => {
                if (phone != undefined && phone != "" && !/^[1][3,4,5,6,7,8,9][0-9]{9}$/.test(phone)) {
                  message.error("请输入正确的手机号码")
                  return;
                }
                let captcha = form?.getFieldValue("captcha");
                if (captcha == undefined || captcha.length != 4) {
                  message.error("请输入正确的图形验证码")
                  return;
                }
                setIsSmsSendLoading(true);
                sendSms({phone: phone?.toString() || "", token: token, captcha: captcha}).then((r) => {
                  setIsSmsSendLoading(false);
                  changeCaptcha();
                  if (r.code != 200) {
                    message.error(r.message || "请求失败").then()
                  } else {
                    getCode(e);
                    message.success(r.message);
                  }
                });
              }}>
                获取验证码
              </Button>
            </Col>
          </Row>
        </Form.Item>
        <Form.Item name={"captcha"} label="图形验证码" hidden={isCaptchaHidden}>
          <Row gutter={6} wrap={false}>
            <Col flex={20}>
              <Input placeholder="请输入图形验证码"/>
            </Col>
            <Col flex={4}>
              <Image
                height={30}
                width={70}
                preview={false}
                src={captchaUrl}
                onClick={() => {
                  changeCaptcha()
                }}
              />
            </Col>
          </Row>
        </Form.Item>
        {bottomNodes}
      </Form>
    </>
  );
};

export default SmsVerify;
