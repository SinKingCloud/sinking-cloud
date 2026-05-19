import React from "react";
import {Form} from "antd";
import {BaseForm, FormField, Loading} from "./common";
import {useSystemConfig} from "../hooks";
import {webConfig} from "@/service/admin/config";

const WebView: React.FC = () => {
    const {form, dataLoading, submitLoading, saveConfig, resetForm} = useSystemConfig({
        request: webConfig,
    });

    if (dataLoading) {
        return <Loading/>;
    }

    return (
        <Form form={form} layout="vertical">
            <BaseForm form={form} submitLoading={submitLoading} onFinish={saveConfig} onReset={() => resetForm()}>
                <FormField type="input" name="web.url" label="网站地址" placeholder="请输入网站地址" required/>
                <FormField type="input" name="web.api" label="网站API地址" placeholder="请输入网站API地址" required/>
                <FormField type="input" name="web.name" label="网站名称" placeholder="请输入网站名称" required/>
                <FormField type="input" name="web.title" label="网站标题" placeholder="请输入网站标题" required/>
                <FormField type="input" name="web.keywords" label="网站关键词" placeholder="请输入网站关键词"/>
                <FormField type="textarea" name="web.describe" label="网站描述" placeholder="请输入网站描述" rows={4}/>
                <FormField type="input" name="web.contact" label="网站联系人" placeholder="请输入网站联系人"/>
                <FormField type="input" name="web.logo" label="网站logo" placeholder="请输入网站logo"/>
            </BaseForm>
        </Form>
    );
};

export default WebView;
