import React, {useState} from 'react';
import {Button, Form, Input, message, Modal, Select} from "antd";
import {changeUinPassword, changeUinPhone, uinLogin} from "@/services/open/uin/uin";
// @ts-ignore
import CryptoJS from 'crypto-js';
// @ts-ignore
import Script from 'react-load-script';
import {getUrlQuery} from "@/util/string";
// @ts-ignore
import QRCode from "qrcode.react/lib";
import {checkMobile} from "@/util/device";

export type UpdateUinProps = {
  formRef?: any;
  uin?: string;
  visible?: boolean;
  onOk?: () => void;
  onCancel?: () => void;
};

const UpdateUin: React.FC<UpdateUinProps> = (props) => {
  const {uin, formRef, visible, onOk, onCancel} = props;
  const [step, setStep] = useState("normal");
  const [qrCode, setQrCode] = useState("");
  const [isQrCodeModalVisible, setIsQrCodeModalVisible] = useState(false);
  const [isQrCodeConfirmLoading, setIsQrCodeConfirmLoading] = useState(false);
  const showQrCodeModal = () => {
    setIsQrCodeModalVisible(true);
  };
  const jumpQQ = (qrurl: string) => {
    const ua = window.navigator.userAgent.toLowerCase();
    const is_ios = ua.indexOf('iphone') > -1 || ua.indexOf('ipad') > -1;
    let schemacallback = '';
    if (is_ios) {
      schemacallback = 'weixin://';
    } else if (ua.indexOf('ucbrowser') > -1) {
      schemacallback = 'ucweb://';
    } else if (ua.indexOf('meizu') > -1) {
      schemacallback = 'mzbrowser://';
    } else if (ua.indexOf('liebaofast') > -1) {
      schemacallback = 'lb://';
    } else if (ua.indexOf('baidubrowser') > -1) {
      schemacallback = 'bdbrowser://';
    } else if (ua.indexOf('baiduboxapp') > -1) {
      schemacallback = 'bdapp://';
    } else if (ua.indexOf('mqqbrowser') > -1) {
      schemacallback = 'mqqbrowser://';
    } else if (ua.indexOf('qihoobrowser') > -1) {
      schemacallback = 'qihoobrowser://';
    } else if (ua.indexOf('chrome') > -1) {
      schemacallback = 'googlechrome://';
    } else if (ua.indexOf('sogoumobilebrowser') > -1) {
      schemacallback = 'SogouMSE://';
    } else if (ua.indexOf('xiaomi') > -1) {
      schemacallback = 'miuibrowser://';
    } else {
      schemacallback = 'googlechrome://';
    }
    if (is_ios) {
      message.warn('跳转到QQ登录后请手动返回当前浏览器');
      window.location.href = 'wtloginmqq3://ptlogin/qlogin?qrcode=' + encodeURIComponent(qrurl) + '&schemacallback=' + encodeURIComponent(schemacallback);
    } else {
      window.location.href = 'wtloginmqq://ptlogin/qlogin?qrcode=' + encodeURIComponent(qrurl) + '&schemacallback=' + encodeURIComponent(schemacallback);
    }
  }
  const qrCodeHandleOk = () => {
    setIsQrCodeConfirmLoading(true);
    uinLogin({uin: initUin, step: "qrcode", login_type: "pwd"}).then((r) => {
      if (r?.code != 200) {
        message.warn("请使用手机QQ扫描二维码并授权登录")
      }
      loginHandle(r);
    }).catch(() => {
      setIsQrCodeConfirmLoading(false);
    });
  };
  const qrCodeHandleCancel = () => {
    setIsQrCodeModalVisible(false);
  };
  /**
   * 表单处理
   */
  const [form] = Form.useForm();
  const [confirmLoading, setConfirmLoading] = useState(false);
  const captchaLogin = (url: string) => {
    const obj = getUrlQuery(url);
    //需要滑块验证码
    let top = document.body.scrollTop || document.documentElement.scrollTop
    scrollBy(0, -top);
    //截取sid
    const sid = obj["sid"];
    const cap_cd = obj["cap_cd"];
    // @ts-ignore
    sinking.setAid(16).setData(initUin, sid, cap_cd).init(function (t: any) {
      //验证成功
      uinLogin({uin: initUin, step: "captcha", ticket: t?.ticket || "ticket", login_type: "pwd"}).then((r) => {
        loginHandle(r);
      }).catch(() => {
        setConfirmLoading(false);
      });
    });
  }
  const loginHandle = (r: any) => {
    if (r.code == 200) {
      setConfirmLoading(false);
      setIsQrCodeConfirmLoading(false);
      setIsQrCodeModalVisible(false);
      //登陆成功
      message.success(r.message || "更新状态成功");
      onOk?.();
    } else if (r.code == 400) {
      captchaLogin(r.data);
    } else if (r.code == 401) {
      setConfirmLoading(false);
      setIsQrCodeConfirmLoading(false);
      //需要短信验证码
      setStep("sms");
      setQrCode(r?.data);
    } else if (r.code == 402) {
      setConfirmLoading(false);
      setIsQrCodeConfirmLoading(false);
      //必须要扫码验证
      setQrCode(r?.data);
      showQrCodeModal();
    } else {
      setStep("normal");
      setConfirmLoading(false);
      setIsQrCodeConfirmLoading(false);
      message.error(r.message || "更新状态失败");
      return
    }
  }
  const onFinish = async (values: any) => {
    if (step == "normal") {
      setConfirmLoading(true);
      if (values.login_type == "pwd") {
        if ((values.pwd?.length || 0) > 5) {
          await changeUinPassword({uin: initUin, pwd: CryptoJS.MD5(values.pwd).toString()})
        }
      } else if (values.login_type == "phone") {
        if ((values.phone?.length || 0) == 11) {
          await changeUinPhone({uin: initUin, phone: values.phone})
        }
      }
      //发起登陆请求
      if (values.login_type == "phone") {
        //需要滑块验证码
        let top = document.body.scrollTop || document.documentElement.scrollTop
        scrollBy(0, -top);
        // @ts-ignore
        sinking.setAid(2081081773).init(function (t: any) {
          //验证成功
          uinLogin({
            uin: initUin,
            step: "normal",
            ticket: t?.ticket || "ticket",
            rand: t?.randstr || "rand", login_type: values.login_type
          }).then((r) => {
            loginHandle(r);
          }).catch(() => {
            setConfirmLoading(false);
          });
        });
      } else {
        uinLogin({uin: initUin, step: "normal", login_type: values.login_type}).then((r) => {
          loginHandle(r)
        });
      }
    } else {
      if ((values.sms?.length || 0) != 6) {
        message.error("请输入正确的验证码")
        return;
      }
      setConfirmLoading(true);
      uinLogin({uin: initUin, step: "sms", code: values.sms, login_type: values.login_type}).then((r) => {
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
    setQrCode("");
    setInitUin(uin);
    setConfirmLoading(false);
  }
  const [initLoginType, setInitLoginType] = useState(formRef?.current?.getFieldValue("login_type"));
  return (
    <>
      <Modal title="请使用QQ扫描二维码" key={"qrCodeVerify"
      } width={370} destroyOnClose={true} maskClosable={false} visible={isQrCodeModalVisible}
             onOk={qrCodeHandleOk} onCancel={qrCodeHandleCancel} confirmLoading={isQrCodeConfirmLoading}
             okText={"我已扫码"}>
        <QRCode
          style={{margin: "0px auto"}}
          value={qrCode}
          size={150}
          fgColor="#000000"
        />
        <div style={{textAlign: "center", marginTop: "20px", marginBottom: "-10px"}} hidden={!checkMobile()}>
          <Button type="primary" onClick={() => {
            jumpQQ(qrCode);
          }}>跳转QQ快捷登录</Button>
        </div>
      </Modal>
      <Modal title={"更新挂机"} key={"UpdateUin"} width={370} destroyOnClose={true} forceRender={true} maskClosable={false}
             visible={visible} onOk={form.submit} okText={"登陆"}
             onCancel={() => {
               formRef?.current?.resetFields();
               setInitLoginType(formRef?.current?.getFieldValue("login_type"))
               onCancel?.();
             }} confirmLoading={confirmLoading}>
        <Script url={"https://cdn.open-api.cn/app/uin/captcha/captcha.js"}/>
        <Form name="control-hooks" form={form} ref={formRef} onFinish={onFinish} labelAlign="right"
              labelCol={{span: 6}}>
          <Form.Item name={"login_type"} label="登陆方式" rules={[{required: true}]}>
            <Select placeholder="请选择登陆方式" onChange={(e) => {
              setInitLoginType(e);
            }}>
              <Select.Option value={"pwd"}>密码登陆</Select.Option>
              <Select.Option value={"phone"}>手机号登陆</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name={"phone"}
                     hidden={step != "normal" || (initLoginType != "phone" && form.getFieldValue("login_type") != "phone")}
                     label="登陆手机">
            <Input placeholder="留空则使用上次登陆的手机"/>
          </Form.Item>
          <Form.Item name={"pwd"}
                     hidden={step != "normal" || (initLoginType != "pwd" && form.getFieldValue("login_type") != "pwd")}
                     label="挂机密码">
            <Input placeholder="留空则使用上次更新的密码"/>
          </Form.Item>
          <Form.Item name={"sms"} hidden={step != "sms"} label="短信验证码">
            <Input placeholder="请输入短信验证码" type={"number"}/>
          </Form.Item>
          <div hidden={initLoginType == "phone" || qrCode == ""}
               style={{float: "right", marginTop: "-10px", marginBottom: "-10px"}}>
            <a onClick={showQrCodeModal}>切换扫码验证</a></div>
        </Form>
      </Modal>
    </>
  );
};

export default UpdateUin;
