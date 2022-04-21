import React from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import {Card, Typography} from 'antd';


export default (): React.ReactNode => {
  return (
    <PageContainer title={false}>
      <Card>
        <Typography.Text strong>
          大道至简，简约但不简单。
        </Typography.Text>
      </Card>
    </PageContainer>
  );
};
