import React, {useState} from 'react';
import {PageContainer} from "@ant-design/pro-layout";
import ProCard from '@ant-design/pro-card';

export default (): React.ReactNode => {
  const [tab, setTab] = useState('tab1');
  return (
    <PageContainer title={false}>
      <ProCard
        tabs={{
          activeKey: tab,
          onChange: (key) => {
            setTab(key);
          },
        }}
      >
        <ProCard.TabPane key="tab1" tab="产品一">
          内容一
        </ProCard.TabPane>
        <ProCard.TabPane key="tab2" tab="产品二">
          内容二
        </ProCard.TabPane>
        <ProCard.TabPane key="tab3" tab="产品一">
          内容一
        </ProCard.TabPane>
        <ProCard.TabPane key="tab4" tab="产品一">
          内容一
        </ProCard.TabPane>
        <ProCard.TabPane key="tab5" tab="产品一">
          内容一
        </ProCard.TabPane>
        <ProCard.TabPane key="tab6" tab="产品一">
          内容一
        </ProCard.TabPane>
      </ProCard>
    </PageContainer>
  );
};
