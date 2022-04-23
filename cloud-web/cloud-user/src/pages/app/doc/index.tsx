import React from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import {Card, Typography} from 'antd';

export default (): React.ReactNode => {
  return (
    <PageContainer title={false}>
      <Card>
        <Typography.Text strong>
          请移步至<a href={"https://doc.open-api.cn"} target={"_blank"}>https://doc.open-api.cn</a>查看开发文档
        </Typography.Text>
      </Card>
    </PageContainer>
  );
};
