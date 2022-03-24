import React from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import {Alert, Card} from 'antd';

export default (): React.ReactNode => {
  return (
    <PageContainer title={false}>
      <Card>
        <Alert
          message="文档"
          type="success"
          showIcon
          banner
          style={{
            margin: -12,
            marginBottom: 24,
          }}
        />
      </Card>
    </PageContainer>
  );
};
