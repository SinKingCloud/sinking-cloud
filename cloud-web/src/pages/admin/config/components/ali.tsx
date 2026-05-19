import React from "react";
import {Form} from "antd";
import {BaseForm, FormField, Loading} from "./common";
import {useSystemConfig} from "../hooks";
import {aliyunConfig} from "@/service/admin/config";

const AliView: React.FC = () => {
    const {form, dataLoading, submitLoading, saveConfig, resetForm} = useSystemConfig({
        request: aliyunConfig,
    });

    if (dataLoading) {
        return <Loading/>;
    }

    return (
        <Form form={form} layout="vertical">
            <BaseForm form={form} submitLoading={submitLoading} onFinish={saveConfig} onReset={() => resetForm()}>
                <FormField type="input" name="ali.key" label="阿里云key" placeholder="请输入阿里云key"/>
                <FormField type="input" name="ali.secret" label="阿里云secret" placeholder="请输入阿里云secret"/>
                <FormField type="input" name="ali.sms.sign" label="阿里云sms sign" placeholder="请输入阿里云sms sign"/>
                <FormField type="input" name="ali.sms.template" label="阿里云sms template" placeholder="请输入阿里云sms template"/>
                <FormField type="input" name="ali.sms.var" label="阿里云sms变量名称" placeholder="请输入阿里云sms变量名称"/>
                <FormField type="input" name="ali.mail.account" label="阿里云mail 发信地址" placeholder="请输入阿里云mail 发信地址"/>
                <FormField type="input" name="ali.face.scene" label="阿里云实名认证场景ID" placeholder="请输入阿里云实名认证场景ID"/>
                <FormField type="input" name="ali.face.domain" label="阿里云实名认证场景资料储存oss域名" placeholder="请输入阿里云实名认证场景资料储存oss域名"/>
            </BaseForm>
        </Form>
    );
};

export default AliView;
