import React from "react";
import {Form, Space} from "antd";
import {createStyles} from "antd-style";
import {BaseForm, FormField, Loading} from "./common";
import {useSystemConfig} from "../hooks";
import {uiConfig} from "@/service/admin/config";

type UiConfigFormProps = {
    title: string;
    prefix: string;
};

const layoutOptions = [
    {value: "top", label: "上下布局"},
    {value: "left", label: "左右布局"},
];

const themeOptions = [
    {value: "light", label: "亮色模式"},
    {value: "dark", label: "暗色模式"},
];

const switchOptions = [
    {value: "1", label: "开启"},
    {value: "0", label: "关闭"},
];

const radiusOptions = Array.from({length: 16}, (_, index) => ({
    value: String(index),
    label: `${index}px${index === 0 ? " (无圆角)" : ""}`,
}));

const useStyles: any = createStyles(({css, token}: any) => ({
    wrapper: css`
        width: 100%;
        max-width: 720px;
    `,
    title: css`
        margin-bottom: ${token.marginSM}px;
        font-size: ${token.fontSizeLG}px;
        font-weight: 600;
        color: ${token.colorText};
    `,
}));

const UiConfigForm = ({title, prefix}: UiConfigFormProps): React.ReactNode => {
    const {styles} = useStyles();
    const {form, dataLoading, submitLoading, saveConfig, resetForm} = useSystemConfig({
        request: uiConfig,
    });

    return (
        <div>
            <div className={styles.title}>{title}</div>
            {dataLoading ? <Loading/> : (
                <Form form={form} layout="vertical">
                    <BaseForm form={form} submitLoading={submitLoading} onFinish={saveConfig}
                              onReset={() => resetForm()}>
                        <FormField type="select" name={`${prefix}.compact`} label="紧凑模式"
                                   tooltip="控制界面元素间距，开启后界面更紧凑"
                                   placeholder="请选择紧凑模式是否开启"
                                   options={switchOptions} required/>
                        <FormField type="select" name={`${prefix}.layout`} label="界面布局"
                                   tooltip="选择界面的整体布局方式，上下布局或左右布局"
                                   placeholder="请选择界面布局方式"
                                   options={layoutOptions} required/>
                        <FormField type="select" name={`${prefix}.theme`} label="菜单主题"
                                   tooltip="选择菜单的主题模式，亮色或暗色"
                                   placeholder="请选择菜单主题模式"
                                   options={themeOptions} required/>
                        <FormField type="select" name={`${prefix}.watermark`} label="界面水印"
                                   tooltip="控制界面是否显示用户昵称水印效果"
                                   placeholder="请选择是否显示界面水印"
                                   options={switchOptions} required/>
                        <FormField type="select" name={`${prefix}.radius`} label="主题圆角"
                                   tooltip="调整界面元素的圆角大小，数值越大圆角越明显"
                                   placeholder="请选择主题圆角大小"
                                   options={radiusOptions} required/>
                        <FormField type="color" name={`${prefix}.color`} label="主题颜色"
                                   tooltip="选择界面的主题颜色，影响整体视觉风格"
                                   format="rgb" defaultColor="rgb(0, 83, 251)"/>
                    </BaseForm>
                </Form>
            )}
        </div>
    );
};

const UiView: React.FC = () => {
    const {styles} = useStyles();

    return (
        <Space className={styles.wrapper} orientation="vertical" size={16}>
            <UiConfigForm title="后台界面" prefix="ui.admin"/>
            <UiConfigForm title="前台界面" prefix="ui.user"/>
        </Space>
    );
};

export default UiView;
