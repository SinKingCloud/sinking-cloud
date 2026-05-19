import {useEffect, useState} from "react";
import {App, Form} from "antd";
import {useModel} from "umi";

type UseSystemConfigOptions = {
    request: any;
    saveRequest?: any;
    enabled?: boolean;
    normalize?: (data: any) => any;
    format?: (values: any) => Record<string, any>;
    getBody?: any;
};

export const useSystemConfig = ({
                                    request,
                                    saveRequest,
                                    enabled = true,
                                    normalize,
                                    format,
                                    getBody,
                                }: UseSystemConfigOptions) => {
    const {message} = App.useApp();
    const webModel = useModel("web");
    const [form] = Form.useForm();
    const [dataLoading, setDataLoading] = useState(false);
    const [submitLoading, setSubmitLoading] = useState(false);

    const loadConfig = () => {
        if (!enabled) {
            return;
        }
        setDataLoading(true);
        request({
            body: getBody === undefined ? {action: "get"} : getBody,
            onSuccess: (res: any) => {
                const data = normalize ? normalize(res?.data || {}) : res?.data || {};
                form.setFieldsValue(data);
            },
            onFail: (res: any) => {
                message?.error(res?.message || "加载配置失败");
            },
            onFinally: () => setDataLoading(false),
        });
    };

    const saveConfig = (values: any) => {
        if (!enabled) {
            message?.warning("暂无操作权限");
            return;
        }
        setSubmitLoading(true);
        const data: any = format ? format(values || {}) : Object.fromEntries(Object.entries(values || {}).map(([key, value]) => [
            key,
            value === undefined || value === null ? "" : String(value),
        ]));
        const requestBody = saveRequest ? data : {
            action: "set",
            data,
        };
        (saveRequest || request)({
            body: requestBody,
            onSuccess: (res: any) => {
                message?.success(res?.message || "配置保存成功");
                loadConfig();
                webModel?.refreshInfo?.();
            },
            onFail: (res: any) => {
                message?.error(res?.message || "配置保存失败");
            },
            onFinally: () => setSubmitLoading(false),
        });
    };

    const resetForm = (names?: any[]) => {
        form.resetFields(names);
    };

    useEffect(() => {
        loadConfig();
    }, [enabled]);

    return {
        form,
        dataLoading,
        submitLoading,
        loadConfig,
        saveConfig,
        resetForm,
    };
};

export default useSystemConfig;
