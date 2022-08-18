import React, {useEffect, useState} from 'react';
import {Button, Form, Input, message, Modal, Select, Typography} from "antd";
import {changeUinPassword, changeUinPhone, getAreaList, uinLogin, uinSms} from "@/services/open/uin/uin";
// @ts-ignore
import CryptoJS from 'crypto-js';
// @ts-ignore
import Script from 'react-load-script';
import {getUrlQuery, implode} from "@/util/string";
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
    window.location.href = 'mqqapi://forward/url?src_type=web&style=default&=1&version=1&url_prefix=' + window.btoa(qrurl);
  }
  const qrCodeHandleOk = () => {
    setIsQrCodeConfirmLoading(true);
    uinLogin({uin: initUin, step: "qrcode", login_type: "pwd"}).then((r) => {
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
  /**
   * 短信验证处理
   */
  const [smsStep, setSmsStep] = useState("1");//短信验证步骤
  const [smsHash, setSmsHash] = useState("");//短信验证票据
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
          uinSms({uin: initUin, step: "1", rand: t?.randstr, ticket: t?.ticket, hash: smsHash}).then((r) => {
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
          message.warn("请完成滑块验证")
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
        hash: smsHash,
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
        hash: smsHash,
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
  const loginHandle = (r: any) => {
    if (r?.code == 200) {
      setConfirmLoading(false);
      setIsQrCodeConfirmLoading(false);
      setIsQrCodeModalVisible(false);
      //登陆成功
      message.success(r?.message || "更新状态成功");
      onOk?.();
    } else if (r?.code == 400) {
      captchaLogin(r?.data);
    } else if (r?.code == 401) {
      setConfirmLoading(false);
      setIsQrCodeConfirmLoading(false);
      //需要短信验证码
      setStep("sms");
      setQrCode(r?.data);
    } else if (r?.code == 399) {
      setConfirmLoading(false);
      setIsQrCodeConfirmLoading(false);
      //开始走验证短信步骤
      setSmsStep("1");
      setSmsHash(r?.data);
      showSmsModal();
    } else {
      setStep("normal");
      setQrCode("");
      setConfirmLoading(false);
      setIsQrCodeConfirmLoading(false);
      message.error(r?.message || "更新状态失败");
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
        if ((values.phone?.length || 0) >= 5) {
          await changeUinPhone({uin: initUin, phone: values.phone, phone_area: values?.phone_area || ""})
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
  /**
   * 区号数据
   */
  const [areaCode, setAreaCode] = useState<any>([]);
  const getAreaCode = () => {
    getAreaList().then((r) => {
      if (r?.code == 200) {
        setAreaCode(r?.data || []);
      }
    });
  };
  useEffect(() => {
    getAreaCode();
  }, []);
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
               setInitLoginType(formRef?.current?.getFieldValue("login_type"))
               onCancel?.();
             }} confirmLoading={confirmLoading}>
        <Script url={"https://cdn.yunhuoke.cn/app/uin/captcha/captcha.js"}/>
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
          <Form.Item name={"phone_area"}
                     hidden={step != "normal" || (initLoginType != "phone" && form.getFieldValue("login_type") != "phone")}
                     label="手机区号">
            <Select
              placeholder="留空则使用上次登陆的手机区号"
              showSearch
              optionFilterProp="children"
              filterOption={(input, option) => {
                return implode('', option!.children).includes(input.toLowerCase());
              }}
            >
              {areaCode.map((item: any) => {
                return <Select.Option key={item?.area_cn + item?.code}
                                      value={item?.code}>{item?.area_cn}+{item?.code}</Select.Option>
              })}
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
