import React from 'react';
import {Alert} from 'antd';

export default (): React.ReactNode => {
  return (
    <Alert
      message="测试文本"
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
