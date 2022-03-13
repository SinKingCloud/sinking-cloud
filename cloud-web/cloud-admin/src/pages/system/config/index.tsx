import React from 'react';
import {Alert} from 'antd';

export default (): React.ReactNode => {
  return (
    <Alert
      message="配置列表"
      type="success"
      showIcon
      banner
      style={{
        margin: -12,
        marginBottom: 48,
      }}
    />
  );
};
