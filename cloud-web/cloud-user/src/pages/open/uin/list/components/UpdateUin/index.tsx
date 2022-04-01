import React, {useEffect, useState} from 'react';
import {Form, message, Modal} from "antd";
import {getServerList} from "@/services/open/uin/uin";

export type UpdateUinProps = {
  uin?: string;
  pwd?: string;
  server?: string;
  visible?: boolean;
  onOk?: (e: React.MouseEvent<HTMLElement>) => void;
  onCancel?: (e: React.MouseEvent<HTMLElement>) => void;
  onRender?: (servers: any) => void;
};

const UpdateUin: React.FC<UpdateUinProps> = (props) => {
  const {uin, visible, onOk, onCancel, onRender} = props;
  /**
   * 挂机地区列表
   */
  const [, setServers] = useState({})
  useEffect(() => {
    //初始化地区列表
    getServerList().then((r) => {
      if (r.code != 200) {
        message.error(r.message || "请求失败").then()
      } else {
        const map = {};
        r.data.list.forEach((item: any) => {
          map[item?.key] = item;
        });
        setServers(map);
        onRender?.(map)
      }
    })
  }, [])
  return (
    <Modal title={uin ? "更新挂机" : "添加挂机"} key={"UpdateUin"} footer={null} destroyOnClose={true} forceRender={true}
           visible={visible} onOk={onOk}
           onCancel={onCancel}>
      <Form name="control-hooks" onFinish={async (values: any) => {

      }} labelAlign="right"
            labelCol={{span: 6}}>

      </Form>
    </Modal>
  );
};

export default UpdateUin;
