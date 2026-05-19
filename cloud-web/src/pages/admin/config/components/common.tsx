import React from "react";
import {Button, ColorPicker, Flex, Form, Input, InputNumber, Select, Spin} from "antd";
import {createStyles} from "antd-style";

const useStyles: any = createStyles(({css, token}: any) => ({
    actions: css`
        display: flex;
        gap: ${token.marginSM}px;
        margin-top: ${token.marginLG}px;
    `,
    loading: css`
        min-height: 240px;
        display: flex;
        align-items: center;
        justify-content: center;
    `,
}));

type BaseFormProps = {
    form: any;
    children: React.ReactNode;
    submitLoading?: boolean;
    onFinish: (values: any) => void;
    onReset?: () => void;
};

type FormFieldProps = {
    type: "input" | "textarea" | "number" | "select" | "color";
    name: string;
    label: string;
    tooltip?: string;
    placeholder?: string;
    required?: boolean;
    rows?: number;
    min?: number;
    precision?: number;
    options?: any[];
    rules?: any[];
    maxWidth?: string | number;
    format?: "hex" | "rgb";
    defaultColor?: string;
};

export const Loading = (): React.ReactNode => {
    const {styles} = useStyles();
    return (
        <div className={styles.loading}>
            <Spin/>
        </div>
    );
};

export const BaseForm = ({
                             form,
                             children,
                             submitLoading,
                             onFinish,
                             onReset,
                         }: BaseFormProps): React.ReactNode => {
    const {styles} = useStyles();

    const onSubmit = async () => {
        try {
            const values = await form.validateFields();
            onFinish(values);
        } catch (_) {
        }
    };

    return (
        <div>
            {children}
            <div className={styles.actions}>
                <Button onClick={onReset}>
                    重置
                </Button>
                <Button type="primary" onClick={onSubmit} loading={submitLoading}>
                    保存
                </Button>
            </div>
        </div>
    );
};

export const FormField = ({
                              type,
                              name,
                              label,
                              tooltip,
                              placeholder,
                              required,
                              rows = 4,
                              min = 0,
                              precision,
                              options,
                              rules,
                              maxWidth = "400px",
                              format = "rgb",
                          }: FormFieldProps): React.ReactNode => {
    const itemRules = rules || (required ? [{required: true, message: `${type === "select" ? "请选择" : "请输入"}${label}`}] : undefined);
    const itemStyle = {maxWidth, width: "100%"};
    if (type === "textarea") {
        return (
            <Form.Item name={name} label={label} tooltip={tooltip} rules={itemRules} style={itemStyle}>
                <Input.TextArea placeholder={placeholder || `请输入${label}`} rows={rows}/>
            </Form.Item>
        );
    }
    if (type === "number") {
        return (
            <Form.Item name={name} label={label} tooltip={tooltip} rules={itemRules} style={itemStyle}>
                <InputNumber min={min} precision={precision} style={{width: "100%"}}
                             placeholder={placeholder || `请输入${label}`}/>
            </Form.Item>
        );
    }
    if (type === "select") {
        return (
            <Form.Item name={name} label={label} tooltip={tooltip} rules={itemRules} style={itemStyle}>
                <Select placeholder={placeholder || `请选择${label}`} options={options || []}/>
            </Form.Item>
        );
    }
    if (type === "color") {
        return (
            <Form.Item label={label} tooltip={tooltip} style={itemStyle}>
                <Flex gap={8} align="center">
                    <Form.Item name={name} rules={itemRules} noStyle
                               getValueFromEvent={(_: any, value: string) => value}>
                        <ColorPicker showText format={format} defaultFormat={format}/>
                    </Form.Item>
                    <Form.Item shouldUpdate noStyle>
                        {({setFieldValue}: any) => (
                            <Button onClick={() => setFieldValue(name as any, "")}>
                                清空
                            </Button>
                        )}
                    </Form.Item>
                </Flex>
            </Form.Item>
        );
    }
    return (
        <Form.Item name={name} label={label} tooltip={tooltip} rules={itemRules} style={itemStyle}>
            <Input placeholder={placeholder || `请输入${label}`}/>
        </Form.Item>
    );
};
