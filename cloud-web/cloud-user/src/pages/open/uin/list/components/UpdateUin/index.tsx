import React, {useState} from 'react';
import {Form, Input, message, Modal} from "antd";
import {changeUinPassword, uinLogin} from "@/services/open/uin/uin";
// @ts-ignore
import CryptoJS from 'crypto-js';
// @ts-ignore
import Script from 'react-load-script';
import {getUrlQuery} from "@/util/string";

export type UpdateUinProps = {
  uin?: string;
  visible?: boolean;
  onOk?: () => void;
  onCancel?: () => void;
};

const UpdateUin: React.FC<UpdateUinProps> = (props) => {
  const {uin, visible, onOk, onCancel} = props;
  const [step, setStep] = useState("normal");
  /**
   * 表单处理
   */
  const [form] = Form.useForm();
  const [confirmLoading, setConfirmLoading] = useState(false);
  const captchaLogin = (url: string) => {
    const obj = getUrlQuery(url)
    //需要滑块验证码
    let top = document.body.scrollTop || document.documentElement.scrollTop
    scrollBy(0, -top);
    //截取sid
    const sid = obj["sid"];
    const cap_cd = obj["cap_cd"];
    // @ts-ignore
    sinking.setAid(16).setData(initUin, sid, cap_cd).init(function (t: any) {
      //验证成功
      uinLogin({uin: initUin, step: "captcha", ticket: t?.ticket || "ticket"}).then((r) => {
        loginHandle(r);
      }).catch(() => {
        setConfirmLoading(false);
      });
    });
  }
  const loginHandle = (r: any) => {
    if (r.code == 200) {
      //登陆成功
      message.success(r.message || "更新状态成功");
      onOk?.();
    } else if (r.code == 400) {
      captchaLogin(r.data);
    } else if (r.code == 401) {
      setConfirmLoading(false);
      //需要短信验证码
      setStep("sms");
    } else {
      setStep("normal");
      setConfirmLoading(false);
      message.error(r.message || "更新状态失败");
      return
    }
  }
  const onFinish = async (values: any) => {
    if (step == "normal") {
      setConfirmLoading(true);
      if ((values.pwd?.length || 0) > 5) {
        await changeUinPassword({uin: initUin, pwd: CryptoJS.MD5(values.pwd).toString()})
      }
      //发起登陆请求
      uinLogin({uin: initUin, step: "normal"}).then((r) => {
        loginHandle(r)
      });
    } else {
      if ((values.sms?.length || 0) != 6) {
        message.error("请输入正确的验证码")
        return;
      }
      setConfirmLoading(true);
      uinLogin({uin: initUin, step: "sms", code: values.sms}).then((r) => {
        loginHandle(r);
      }).catch(() => {
        setConfirmLoading(false);
      });
    }
  }
  /**
   * 监听数据
   */
  const [initUin, setInitUin] = useState(uin);
  if (initUin != uin) {
    setStep("normal");
    setInitUin(uin);
    setConfirmLoading(false);
  }
  return (
    <Modal title={"更新挂机"} key={"UpdateUin"} width={370} destroyOnClose={true} forceRender={true} maskClosable={false}
           visible={visible} onOk={form.submit} okText={"登陆"}
           onCancel={onCancel} confirmLoading={confirmLoading}>
      <Script url={"https://cdn.open-api.cn/app/uin/captcha/captcha.js"}/>
      <Form name="control-hooks" form={form} onFinish={onFinish} labelAlign="right"
            labelCol={{span: 6}}>
        <Form.Item name={"pwd"} hidden={step != "normal"} label="挂机密码">
          <Input placeholder="不更新则留空"/>
        </Form.Item>
        <Form.Item name={"sms"} hidden={step != "sms"} label="短信验证码">
          <Input placeholder="请输入短信验证码" type={"number"}/>
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default UpdateUin;
