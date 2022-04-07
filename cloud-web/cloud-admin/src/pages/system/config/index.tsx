import React, {useEffect, useState} from 'react';
import {PageContainer} from "@ant-design/pro-layout";
import ProCard from '@ant-design/pro-card';
import ProForm, {ProFormText} from "@ant-design/pro-form";
import {getConfigList, updateConfigs} from "@/services/system/config";
import {message} from "antd";
import {toTreeData} from "@/util/tree";

export default (): React.ReactNode => {
  const [tab, setTab] = useState('');
  const [isCardLoading, setIsCardLoading] = useState(false);
  const [configTrees, setConfigTrees] = useState([]);
  /**
   * 初始化表单值
   */
  const getConfigs = async () => {
    setIsCardLoading(true);
    return await getConfigList({page_size: 1000}).then((r): any => {
      if (r.code != 200) {
        message.error(r.message || "请求失败").then(() => {
          return [];
        })
      } else {
        return r.data.list;
      }
    });
  }
  /**
   * 提交表单
   */
  const onFinish = async (values: any) => {
    let data = {configs: []}
    Object.keys(values).forEach((key) => {
      // @ts-ignore
      data.configs.push({key: key, value: values[key]})
    })
    await updateConfigs(data).then((r) => {
      if (r.code != 200) {
        message.error(r.message || "请求失败").then();
      } else {
        message.success(r.message || "修改成功").then();
      }
    })
  }
  /**
   * 初始化数据
   */
  // @ts-ignore
  useEffect(async () => {
    setIsCardLoading(true)
    const data = toTreeData(await getConfigs(), 0, "id", "pid", "children")
    setConfigTrees(data);
    setTab(data[0]?.id.toString())
    setIsCardLoading(false);
  }, []);

  return (
    <PageContainer title={false}>
      <ProCard
        tabs={{
          activeKey: tab,
          onChange: (key) => {
            setTab(key);
          },
        }}
        loading={isCardLoading}
      >
        {configTrees.map((tab: any) => {
          return (<ProCard.TabPane key={tab?.id} tab={tab?.name}>
            <ProForm key={"form" + tab?.id} onFinish={onFinish}>
              {tab?.children.map((configs: any) => {
                return (<ProFormText
                  allowClear={false}
                  width="md"
                  key={configs?.id}
                  name={configs?.key}
                  initialValue={configs?.value}
                  label={configs?.name}
                  tooltip={configs?.describe}
                  placeholder={"请输入" + configs?.name}
                  rules={[{required: true, message: "请输入" + configs?.name,},]}
                />);
              })}
            </ProForm>
          </ProCard.TabPane>);
        })}
      </ProCard>
    </PageContainer>
  );
};
