import {QqOutlined, QuestionCircleOutlined} from '@ant-design/icons';
import {Modal, Space, Typography} from 'antd';
import React from 'react';
import {useModel} from 'umi';
import Avatar from './AvatarDropdown';
import styles from './index.less';
const GlobalHeaderRight: React.FC = () => {
  const {initialState} = useModel('@@initialState');

  if (!initialState || !initialState.settings) {
    return null;
  }

  const {navTheme, layout} = initialState.settings;
  let className = styles.right;

  if ((navTheme === 'dark' && layout === 'top') || layout === 'mix') {
    className = `${styles.right}  ${styles.dark}`;
  }
  return (
    <Space className={className}>
      <span
        className={styles.action}
        onClick={() => {
          Modal.info({
            title: '7*24小时客服为您解决问题',
            content: (<><QqOutlined/><Typography.Text
              copyable={{tooltips: ['复制', '复制成功']}}>{initialState?.currentWeb?.contact}</Typography.Text></>),
          });
        }}
      >
        <QuestionCircleOutlined/>
      </span>
      <Avatar menu={true}/>
    </Space>
  );
};
export default GlobalHeaderRight;
