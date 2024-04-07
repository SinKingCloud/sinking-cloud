import React, {useEffect, useState} from 'react';
import {PageContainer} from "@ant-design/pro-layout";
import ProCard from '@ant-design/pro-card';
import ProForm, {ProFormDigit, ProFormSelect, ProFormText, ProFormTextArea} from "@ant-design/pro-form";
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
        return [];
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
            <ProForm key={"form" + tab?.id} onFinish={onFinish} autoFocusFirstInput={false}>
              {tab?.children.map((configs: any) => {
                let control;
                const cons = (configs?.control || "").split("|");
                let values: any = {};
                if (cons.length > 1) {
                  for (let i = 1; i < cons.length; i++) {
                    const temp = cons[i].split(":");
                    if (temp.length == 2) {
                      values[temp[0]] = temp[1];
                    }
                  }
                }
                switch (cons[0]) {
                  case "select":
                    control = <ProFormSelect
                      allowClear={false}
                      width="md"
                      key={configs?.id}
                      name={configs?.key}
                      initialValue={configs?.value}
                      label={configs?.name}
                      tooltip={configs?.describe}
                      valueEnum={values}
                      placeholder={"请选择" + configs?.name}
                      rules={[{required: true, message: "请选择" + configs?.name,},]}
                    />
                    break;
                  case "number":
                    control = <ProFormDigit
                      width="md"
                      key={configs?.id}
                      name={configs?.key}
                      initialValue={configs?.value}
                      label={configs?.name}
                      tooltip={configs?.describe}
                      min={values?.min || 0}
                      max={values?.max || 100}
                      placeholder={"请输入" + configs?.name}
                      rules={[{required: true, message: "请输入" + configs?.name,},]}
                    />
                    break
                  case "text_area":
                    control = <ProFormTextArea
                      allowClear={false}
                      width="md"
                      key={configs?.id}
                      name={configs?.key}
                      initialValue={configs?.value}
                      label={configs?.name}
                      tooltip={configs?.describe}
                      placeholder={"请输入" + configs?.name}
                      rules={[{required: true, message: "请输入" + configs?.name,},]}
                    />
                    break
                  default:
                    control = <ProFormText
                      allowClear={false}
                      width="md"
                      key={configs?.id}
                      name={configs?.key}
                      initialValue={configs?.value}
                      label={configs?.name}
                      tooltip={configs?.describe}
                      placeholder={"请输入" + configs?.name}
                      rules={[{required: true, message: "请输入" + configs?.name,},]}
                    />
                    break;
                }
                return control;
              })}
            </ProForm>
          </ProCard.TabPane>);
        })}
      </ProCard>
    </PageContainer>
  );
};
