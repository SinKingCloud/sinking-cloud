import React, {useEffect, useState} from 'react';
import {Form, Input, message, Modal, Select} from "antd";
import {getSelfInfo} from "@/services/open/uin/home";
// @ts-ignore
import CryptoJS from 'crypto-js';
import {getServerList, uinBuy} from "@/services/open/uin/uin";

export type BuyProps = {
  uin?: string;
  visible?: boolean;
  onCancel?: () => void;
  onSuccess?: (uin: string, type: number, month: number, isNewUin: boolean) => void;
};

const Buy: React.FC<BuyProps> = (props) => {
  const {uin, visible, onCancel, onSuccess} = props;
  /**
   * 初始化数据
   */
  const [servers, setServers] = useState({})
  const [selfInfo, setSelfInfo] = useState({});
  const init = () => {
    getSelfInfo().then((r) => {
      if (r.code == 200) {
        setSelfInfo(r.data);
      }
    });
    getServerList().then((r) => {
      if (r.code != 200) {
        message.error(r.message || "请求失败").then();
      } else {
        const map = {};
        r.data.list.forEach((item: any) => {
          map[item?.key] = item;
        });
        setServers(map);
      }
    });
  }
  useEffect(() => {
    init();
  }, []);
  /**
   * 处理表单
   */
  const [form] = Form.useForm();
  const onFinish = (values: any) => {
    let isNewUin = false
    if (initUin != "" && (initUin?.length || 0) >= 5) {
      //续费
      values.uin = initUin;
      delete values.pwd;
      delete values.server_key;
    } else {
      isNewUin = true;
      if (values.uin == undefined || values.uin == "") {
        message.error("请输入挂机账号");
        return
      }
      if (values.pwd == undefined || values.pwd == "") {
        message.error("请输入挂机密码");
        return
      } else {
        values.pwd = CryptoJS.MD5(values.pwd).toString();
      }
      if (values.server_key == undefined || values.server_key == "") {
        message.error("请选择挂机地区");
        return
      }
    }
    //关闭弹窗
    onCancel?.();
    const key = "uinBuy";
    message.loading({content: '正在开通挂机服务', key, duration: 60}).then();
    uinBuy(values).then((r) => {
      if (r.code == 200) {
        form.resetFields();
        message.success({content: r.message, key, duration: 2}).then();
        onSuccess?.(values.uin, values.type, values.month, isNewUin);
      } else {
        message.error({content: r.message, key, duration: 2}).then();
      }
    });
  }
  /**
   * 处理价格
   */
  const [option, setOption] = useState(<></>);
  const onChange = (e: number) => {
    let type = "";
    if (e == 1) {
      type = "player_";
    } else if (e == 2) {
      type = "sign_"
    }
    // @ts-ignore
    const obj = selfInfo.level;
    const ele = (<><Select.Option
      value={1}>月度(1个月){obj.hasOwnProperty(type + "month_price") && ",价格" + (obj[type + "month_price"].toFixed("2") || "0.00") + "元"}</Select.Option>
      <Select.Option
        value={3}>季度(3个月){obj.hasOwnProperty(type + "quarter_price") && ",价格" + (obj[type + "quarter_price"].toFixed("2") || "0.00") + "元"}</Select.Option>
      <Select.Option
        value={6}>半年度(6个月){obj.hasOwnProperty(type + "half_price") && ",价格" + (obj[type + "half_price"].toFixed("2") || "0.00") + "元"}</Select.Option>
      <Select.Option
        value={12}>年度(12个月){obj.hasOwnProperty(type + "year_price") && ",价格" + (obj[type + "year_price"].toFixed("2") || "0.00") + "元"}</Select.Option></>);
    setOption(ele);
  }
  /**
   * 监听数据
   */
  const [initUin, setInitUin] = useState(uin);
  if (initUin != uin) {
    setInitUin(uin);
    onChange(0);
  }
  return (
    <Modal title={(initUin?.length || 0) >= 5 ? "挂机续费" : "挂机开通"} key={"Buy"} width={350} destroyOnClose={true}
           forceRender={true}
           visible={visible} onOk={form.submit} okText={(initUin?.length || 0) >= 5 ? "续费" : "开通"}
           onCancel={onCancel}>
      <Form form={form} name="control-hooks" onFinish={onFinish}
            labelAlign="right"
            labelCol={{span: 6}}
            wrapperCol={{span: 16}}>
        <Form.Item name={"uin"} hidden={(initUin?.length || 0) >= 5} label="挂机账号">
          <Input placeholder="请输入挂机账号"/>
        </Form.Item>
        <Form.Item name={"pwd"} hidden={(initUin?.length || 0) >= 5} label="挂机密码">
          <Input placeholder="请输入挂机密码"/>
        </Form.Item>
        <Form.Item name={"server_key"} hidden={(initUin?.length || 0) >= 5} label="挂机地区">
          <Select placeholder="请选择挂机地区">
            {Object.keys(servers).map((k) => {
              return <Select.Option key={k} value={k}>{servers?.[k]?.name}</Select.Option>;
            })}
          </Select>
        </Form.Item>
        <Form.Item name={"type"} label="套装类型" rules={[{required: true}]} initialValue={0}>
          <Select placeholder="请选择套装" onChange={onChange}>
            <Select.Option value={0}>全套</Select.Option>
            <Select.Option value={1}>加速</Select.Option>
            <Select.Option value={2}>扩展</Select.Option>
          </Select>
        </Form.Item>
        <Form.Item name={"month"} label="开通时长" rules={[{required: true}]} initialValue={1}>
          <Select placeholder="请选择时长">
            {option}
          </Select>
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default Buy;
