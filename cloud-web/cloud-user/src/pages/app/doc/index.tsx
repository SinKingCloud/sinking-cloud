import React from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import {Card, Typography} from 'antd';

export default (): React.ReactNode => {
  return (
    <PageContainer title={false}>
      <Card>
        <Typography.Text strong>
          页面正在建设中
        </Typography.Text>
      </Card>
    </PageContainer>
  );
};
