import React, {useState} from 'react';
import {Alert, Button, Form, Input, message, Modal, Typography} from "antd";
import {changeUinPassword, uinLogin, uinSms} from "@/services/open/uin/uin";
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
  const [uinPhone, setUinPhone] = useState("");
  const [isQrCodeModalVisible, setIsQrCodeModalVisible] = useState(false);
  const [isQrCodeConfirmLoading, setIsQrCodeConfirmLoading] = useState(false);
  const showQrCodeModal = () => {
    setIsQrCodeModalVisible(true);
  };
  const jumpQQ = (qrurl: string) => {
    window.location.href = 'mqqapi://forward/url?src_type=web&style=default&=1&version=1&url_prefix=' + window.btoa(qrurl);
  }
  const qrCodeHandleOk = () => {
    setIsQrCodeConfirmLoading(true);
    uinLogin({uin: initUin, step: "normal"}).then((r) => {
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
      if ((t?.ticket || "") == "" || (t?.randstr || "") == "") {
        setConfirmLoading(false);
        message.warn("请完成滑块验证");
        return;
      }
      //验证成功
      uinLogin({uin: initUin, step: "captcha", ticket: t?.ticket || "ticket"}).then((r) => {
        loginHandle(r);
      }).catch(() => {
        setConfirmLoading(false);
      });
    });
  }
  /**
   * 短信验证处理
   */
  const [smsStep, setSmsStep] = useState("1");//短信验证步骤
  const [smsTicket, setSmsTicket] = useState("");//短信验证ticket
  const [smsRandStr, setSmsRandStr] = useState("");//短信验证rand_str
  const [smsData, setSmsData] = useState<any>({});//返回data
  const [smsPhone, setSmsPhone] = useState("");//手机号输入框
  const [isSmsModalVisible, setIsSmsModalVisible] = useState(false);//短信验证弹窗
  const [isSmsConfirmLoading, setIsSmsConfirmLoading] = useState(false);//短信验证按钮
  const showSmsModal = () => {
    setIsSmsModalVisible(true);
  };
  const hideSmsModal = () => {
    setIsSmsModalVisible(false);
  };
  const smsHandleOk = () => {
    setIsSmsConfirmLoading(true);
    //步骤一:初始化滑块并获取绑定手机
    if (smsStep == "1") {
      //初始化滑块
      // @ts-ignore
      sinking.setAid(2090581062).init(function (t: any) {
        if ((t?.ticket || "") != "" && (t?.randstr || "") != "") {
          setSmsTicket(t?.ticket);
          setSmsRandStr(t?.randstr);
          //获取绑定手机
          uinSms({uin: initUin, step: "1", rand: t?.randstr, ticket: t?.ticket}).then((r) => {
            setIsSmsConfirmLoading(false);
            if (r?.code != 200) {
              message.error(r?.message || "请求失败");
            } else {
              setSmsData(r?.data);
              setSmsStep("2");
            }
          });
        } else {
          setIsSmsConfirmLoading(false);
          message.warn("请完成滑块验证");
        }
      });
    }
    //步骤二:填写手机号并发起验证请求
    if (smsStep == "2") {
      const tempPhone = smsData?.phone || "";
      if (smsPhone == "" || smsPhone.length != tempPhone.length) {
        setIsSmsConfirmLoading(false);
        message.error("请填写完整手机号码");
        return;
      }
      if (!/^[0-9]+.?[0-9]*$/.test(smsPhone)) {
        message.error("手机号码必须为数字");
        return;
      }
      //验证手机号是否正确
      for (let i = 0; i < tempPhone.length; i++) {
        if (tempPhone.charAt(i) != '*' && smsPhone.charAt(i) != tempPhone.charAt(i)) {
          setIsSmsConfirmLoading(false);
          message.error("请填写正确的手机号");
          return;
        }
      }
      //发起请求
      uinSms({
        uin: initUin,
        step: "2",
        rand: smsRandStr,
        ticket: smsTicket,
        phone: smsPhone
      }).then((r) => {
        setIsSmsConfirmLoading(false);
        if (r?.code != 200) {
          message.error(r?.message || "请求失败");
        } else {
          setSmsData(r?.data);
          setSmsStep("3");
        }
      });
    }
    //步骤三:验证是否发起短信
    if (smsStep == "3") {
      //发起请求
      uinSms({
        uin: initUin,
        step: "3",
        rand: smsRandStr,
        ticket: smsTicket,
        phone: smsPhone
      }).then((r) => {
        setIsSmsConfirmLoading(false);
        if (r?.code != 200) {
          message.error(r?.message || "请求失败");
        } else {
          setSmsData({});
          setSmsStep("1");
          smsHandleCancel();
          form?.submit();
        }
      });
    }
  };
  const smsHandleCancel = () => {
    hideSmsModal()
  };
  //登陆处理
  const loginHandle = (r: any, reset: boolean = true) => {
    if (r?.code == 200) {
      setConfirmLoading(false);
      setIsQrCodeConfirmLoading(false);
      setIsQrCodeModalVisible(false);
      //登陆成功
      message.success(r?.message || "更新状态成功");
      form?.resetFields();
      onOk?.();
    } else if (r?.code == 400) {
      captchaLogin(r?.data);
    } else if (r?.code == 399) {
      setConfirmLoading(false);
      setIsQrCodeConfirmLoading(false);
      //开始走验证短信步骤
      setSmsStep("1");
      showSmsModal();
    } else if (r?.code == 398) {
      setConfirmLoading(false);
      setIsQrCodeConfirmLoading(false);
      //展示选择验证方式弹窗
      if ((r?.data?.qrcode || "") != "") {
        setQrCode(r?.data?.qrcode);
      }
      if ((r?.data?.phone || "") != "") {
        setUinPhone(r?.data?.phone);
      }
      setStep("sms");
    } else {
      if (reset) {
        setStep("normal");
        setQrCode("");
      }
      setConfirmLoading(false);
      setIsQrCodeConfirmLoading(false);
      message.error(r?.message || "更新状态失败");
      return
    }
  }
  const onFinish = async (values: any) => {
    if (step == "normal") {
      setConfirmLoading(true);
      let body;
      let change;
      if ((values.pwd?.length || 0) > 5) {
        body = {uin: initUin, pwd: CryptoJS.MD5(values.pwd).toString()};
        change = changeUinPassword;
      }
      if (change != undefined && body != undefined) {
        const rep = await change(body);
        if (rep?.code != 200) {
          loginHandle(rep);
          return;
        }
      }
      //发起登陆请求
      uinLogin({uin: initUin, step: "normal"}).then((r) => {
        loginHandle(r, false)
      });
    } else {
      if ((values.sms?.length || 0) != 6) {
        message.error("请输入正确的验证码")
        return;
      }
      setConfirmLoading(true);
      uinLogin({uin: initUin, step: "sms", code: values.sms}).then((r) => {
        loginHandle(r, false);
      }).catch(() => {
        setConfirmLoading(false);
      });
    }
  }
  /**
   * 获取验证码
   */
  const [getSmsLoading, setIsGetSmsLoading] = useState(false);
  const getSms = () => {
    const key = "getSms"
    message.loading({content: '正在获取验证码', key, duration: 60}).then();
    setIsGetSmsLoading(true);
    uinLogin({uin: initUin, step: "get_sms"}).then((r) => {
      setIsGetSmsLoading(false);
      if (r?.code == 200) {
        message.success({content: r?.message, key, duration: 2}).then();
      } else {
        message.error({content: r?.message || "获取验证码失败", key, duration: 2}).then();
      }
    });
  }
  /**
   * 监听数据
   */
  const [initUin, setInitUin] = useState(uin);
  if (initUin != uin) {
    setStep("normal");
    setQrCode("");
    setUinPhone("");
    setInitUin(uin);
    setConfirmLoading(false);
  }
  return (
    <>
      <Modal title="身份验证" key={"smsVerify"
      } width={370} destroyOnClose={true} maskClosable={false} visible={isSmsModalVisible}
             onOk={smsHandleOk} onCancel={smsHandleCancel} confirmLoading={isSmsConfirmLoading}
             okText={smsStep == "1" ? "开始验证" : "下一步"}>
        {smsStep == "1" && <>
          <p>
            当前登录存在安全风险，为保障您的账号安全，请使用常用设备或者经过身份验证后登陆
          </p>
        </>}
        {smsStep == "2" && <>
          <p>
            请输入 {smsData?.phone} 的完整手机号验证身份
          </p>
          <div>
            <Input placeholder="请输入完整的手机号" onChange={(e) => {
              setSmsPhone(e?.target?.value);
            }}/>
          </div>
        </>}
        {smsStep == "3" && <>
          <p>
            请使用手机号 {smsPhone}
          </p>
          <p>
            编辑短信 {smsData?.sms || "1"}
          </p>
          <div>
            发送至 <Typography.Text copyable>{smsData?.phone || "106988881700511"}</Typography.Text>
          </div>
          <br/>
          <p>
            发送完毕请回到此页面点击下一步
          </p>
        </>}
      </Modal>
      <Modal title="请使用手机扫描二维码" key={"qrCodeVerify"
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
          }}>跳转手机快捷扫码</Button>
        </div>
      </Modal>
      <Modal title={"更新挂机"} key={"UpdateUin"} width={370} destroyOnClose={true} forceRender={true} maskClosable={false}
             visible={visible} onOk={form.submit} okText={"登陆"}
             onCancel={() => {
               formRef?.current?.resetFields();
               onCancel?.();
             }} confirmLoading={confirmLoading}>
        <Script url={"https://cdn.yunhuoke.cn/app/uin/captcha/captcha.js"}/>
        <Alert message={<>
          <div>请使用验证码验证或手机扫码验证设备锁</div>
          <div hidden={uinPhone == ""}>当前密保手机:{uinPhone}</div>
        </>} type="info"
               style={{marginTop: "-10px", marginBottom: "20px", display: step == "sms" ? "flex" : "none"}} showIcon/>
        <Form name="control-hooks" form={form} ref={formRef} onFinish={onFinish}
              labelAlign="right"
              labelCol={{span: 6}} layout={"horizontal"}>
          <Form.Item name={"pwd"}
                     hidden={step != "normal"}
                     label="挂机密码">
            <Input placeholder="留空则使用上次更新的密码"/>
          </Form.Item>
          <Form.Item hidden={step != "sms"} label="短信验证码" style={{marginBottom: 0}}>
            <Form.Item name={"sms"} style={{display: 'inline-block', width: 'calc(70% - 8px)'}}>
              <Input placeholder="请输入短信验证码" type={"number"} onChange={(e) => {
                form?.setFieldsValue({sms: e.target.value})
              }} style={{maxWidth: "155px"}}/>
            </Form.Item>
            <Form.Item style={{display: 'inline-block', width: 'calc(30% - 8px)', margin: '0 8px'}}>
              <Button type="primary" style={{marginLeft: "10px"}} size={"small"} ghost loading={getSmsLoading}
                      onClick={getSms}>
                获取
              </Button>
            </Form.Item>
          </Form.Item>
          <div hidden={qrCode == "" || step != "sms"}
               style={{float: "right", marginTop: "-10px", marginBottom: "-10px"}}>
            <a onClick={showQrCodeModal}>切换扫码验证</a></div>
        </Form>
      </Modal>
    </>
  );
};

export default UpdateUin;
