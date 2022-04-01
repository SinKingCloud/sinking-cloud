import React from 'react';
import {Modal} from "antd";

export type BuyProps = {
  uin?: string;
  server?: string;
  visible?: boolean;
  onOk?: (e: React.MouseEvent<HTMLElement>) => void;
  onCancel?: (e: React.MouseEvent<HTMLElement>) => void;
};

const Buy: React.FC<BuyProps> = (props) => {
  const {uin, visible, onOk, onCancel} = props;
  return (
    <Modal title={uin ? "挂机续期" : "挂机开通"} key={"Buy"} footer={null} destroyOnClose={true} forceRender={true}
           visible={visible} onOk={onOk}
           onCancel={onCancel}>

    </Modal>
  );
};

export default Buy;
