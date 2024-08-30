import React, {useState} from 'react';
import {Alert, Button, message, Modal, Spin} from "antd";
// @ts-ignore
import {history, useModel} from "umi";
// @ts-ignore
import CryptoJS from 'crypto-js';
// @ts-ignore
import Script from 'react-load-script';
// @ts-ignore
import QRCode from "qrcode.react/lib";

import styles from './style.less';
import {uinQrLogin} from "@/services/open/uin/uin";
import {checkMobile} from "@/util/device";
import {qqJumpUrl} from "@/util/string";

export type QrLoginProps = {
  uin?: string;
  visible?: boolean;
  onOk?: () => void;
  onCancel?: () => void;
};

const QrLogin: React.FC<QrLoginProps> = (props) => {
  const {uin, visible, onOk, onCancel} = props;
  const [qrCode, setQrCode] = useState("");
  const [qrcodeLoading, setQrcodeLoading] = useState(true);
  const [qrcodeMessage, setQrcodeMessage] = useState("请刷新二维码");
  const [qrcodeTips, setQrcodeTips] = useState("点击刷新");
  /**
   * 监听数据
   */
  const [initUin, setInitUin] = useState(uin);
  if (initUin != uin) {
    setInitUin(uin);
    setQrcodeMessage("请刷新二维码")
  }
  const qrCodeHandleOk = () => {
    onOk?.();
  };
  const qrCodeHandleCancel = () => {
    localStorage.removeItem("qrcode_id");
    setQrcodeMessage("请刷新二维码");
    setQrcodeLoading(true)
    setQrcodeTips("点击刷新");
    setQrCode("");
    onCancel?.();
  };
  const queryQrCodeStatus = (id?: string) => {
    setTimeout(() => {
      if (localStorage.getItem("qrcode_id") != id || history.location.pathname != "/open/uin/list") {
        return;
      }
      uinQrLogin({uin: initUin, step: "query"}).then(async (r: any) => {
        if (r?.code == 200) {
          setQrcodeMessage(r?.data?.message || "请扫描二维码");
          if (r?.data?.status == 500 || r?.data?.status == 501 || r?.data?.status == 504 || r?.data?.status == 505 || r?.data?.status == 506) {
            localStorage.removeItem("qrcode_id");
          } else if (r?.data?.status == 200) {
            localStorage.removeItem("qrcode_id");
            qrCodeHandleCancel()
            message.success("更新状态成功");
          } else if (r?.data?.status == 502 || r?.data?.status == 503) {
            queryQrCodeStatus(id);
          }
        } else {
          setQrcodeMessage(r?.message || "查询状态失败");
          setQrcodeLoading(true)
          setQrcodeTips("点击刷新");
        }
      });
    }, 1500);
  }
  const fetchQrCode = () => {
    localStorage.removeItem("qrcode_id");
    setQrcodeLoading(true)
    setQrcodeTips("");
    setQrcodeMessage("正在生成二维码");
    uinQrLogin({uin: initUin, step: "code"}).then((r) => {
      if (r?.code == 200) {
        setQrcodeMessage("请扫描二维码");
        setQrcodeLoading(false);
        const id = r?.data
        setQrCode(id);
        localStorage.setItem("qrcode_id", id);
        queryQrCodeStatus(id);
      } else {
        setQrcodeMessage(r?.message || "获取二维码失败");
        setQrcodeLoading(true)
        setQrcodeTips("点击刷新");
      }
    });
  };
  return (
    <>
      <Modal title="更新挂机" key={"qrLogin"
      } width={370} visible={visible} destroyOnClose={true} maskClosable={false}
             onCancel={qrCodeHandleCancel} onOk={qrCodeHandleOk} footer={null}>
        <Alert message="请使用手机QQ扫码下方二维码" type="info" showIcon/>
        <div style={{
          border: "1px solid #e1e1e1",
          backgroundColor: "#fbfbfb",
          borderRadius: "10px",
          width: "125px",
          height: "125px",
          margin: "15px auto",
          position: "relative",
          cursor: "pointer"
        }} onClick={fetchQrCode}>
          <div className={styles.border_corner + " " + styles.border_corner_right_bottom}></div>
          <div className={styles.border_corner + " " + styles.border_corner_right_top}></div>
          <div className={styles.border_corner + " " + styles.border_corner_left_bottom}></div>
          <div className={styles.border_corner + " " + styles.border_corner_left_top}></div>
          <Spin spinning={qrcodeLoading || qrCode == ""} tip={qrcodeTips}>
            <QRCode
              style={{margin: "14px"}}
              value={qrCode}
              size={95}
              fgColor="#000000"
            />
          </Spin>
        </div>
        <div style={{textAlign: "center", margin: "15px auto -5px auto", color: "#8c8c8c"}}>
          {qrcodeMessage || "请扫描二维码"}
        </div>
        <div hidden={!checkMobile() || qrCode == ""} style={{textAlign: "center", marginTop: "10px"}}>
          <Button type='primary' style={{maxWidth: "200px", marginTop: "10px"}} htmlType='button' block
                  onClick={() => {
                    qqJumpUrl(qrCode);
                    message.warning("验证完毕请回到此页面");
                  }}>
            快捷登陆
          </Button>
        </div>
      </Modal>
    </>
  );
};

export default QrLogin;
