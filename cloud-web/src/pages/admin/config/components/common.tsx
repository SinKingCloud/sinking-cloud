import React, {useRef, useState} from "react";
import {App, Button, ColorPicker, Flex, Form, Image, Input, InputNumber, Select, Spin} from "antd";
import {Icon, useTheme} from "sinking-antd";

import {createStyles} from "antd-style";
import {Upload, type UploadRef} from "@/pages/components/upload";
import {getUploadUrl, SYSTEM_FILE_SCENE_ID} from "@/utils/upload";

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
    uploadBox: css`
        position: relative;
        width: 80px;
        height: 80px;
        display: flex;
        align-items: center;
        justify-content: center;
        border: 1px dashed ${token.colorBorder};
        border-radius: ${token.borderRadius}px;
        background-color: ${token.colorFillAlter};
        cursor: pointer;
        transition: border-color 0.3s, background-color 0.3s;
        overflow: hidden;

        &:hover {
            border-color: ${token.colorPrimary};
            background-color: ${token.colorFillSecondary};
        }
    `,
    uploadContent: css`
        height: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        color: ${token.colorTextSecondary};

        .anticon {
            margin-bottom: 8px;
            font-size: 20px;
        }
    `,
    uploadText: css`
        font-size: ${token.fontSizeSM}px;
        color: ${token.colorTextSecondary};
    `,
    imagePreview: css`
        width: 100%;
        height: 100%;
        object-fit: cover;
    `,
    uploadMask: css`
        position: absolute;
        inset: 0;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 10px;
        opacity: 0;
        background-color: rgba(0, 0, 0, 0.48);
        transition: opacity 0.3s;

        &:hover {
            opacity: 1;
        }

        .anticon {
            padding: 5px;
            color: #fff;
            font-size: 16px;
            cursor: pointer;
            transition: transform 0.3s;

            &:hover {
                transform: scale(1.18);
            }
        }
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
    type: "input" | "textarea" | "number" | "select" | "color" | "upload";
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
    accept?: string;
    sceneId?: number | string;
    uploadText?: string;
};

type UploadFieldProps = {
    name: string;
    label: string;
    tooltip?: string;
    rules?: any[];
    style?: any;
    accept?: string;
    sceneId?: number | string;
    uploadText?: string;
};

const imageAccept = ".svg,.svgz,image/svg+xml,image/png,image/jpeg,image/webp,image/gif";

const HiddenValue: React.FC<{ value?: string; onChange?: (value: string) => void }> = () => null;

const UploadField = ({
                         name,
                         label,
                         tooltip,
                         rules,
                         style,
                         accept = imageAccept,
                         sceneId = SYSTEM_FILE_SCENE_ID,
                         uploadText = "上传",
                     }: UploadFieldProps): React.ReactNode => {
    const {message, modal} = App.useApp();
    const {styles} = useStyles();
    const uploadRef = useRef<UploadRef>({} as UploadRef);
    const [loading, setLoading] = useState(false);
    const [previewImage, setPreviewImage] = useState("");

    return (
        <Form.Item label={label} tooltip={tooltip || label} style={style}>
            <Upload ref={uploadRef}/>
            <Form.Item name={name as any} rules={rules} noStyle>
                <HiddenValue/>
            </Form.Item>
            <Form.Item shouldUpdate={(prevValues, currentValues) => prevValues?.[name] !== currentValues?.[name]} noStyle>
                {({getFieldValue, setFieldValue}: any) => {
                    const value = getFieldValue(name as any);
                    const uploadFile = () => {
                        if (loading) {
                            return;
                        }
                        setLoading(true);
                        const task = uploadRef.current?.upload?.({
                            sceneId,
                            accept,
                            onComplete: (result: any) => {
                                const url = getUploadUrl(result);
                                if (!url) {
                                    message?.warning(`${label}上传成功，但未返回地址`);
                                    setLoading(false);
                                    return;
                                }
                                setFieldValue(name as any, url);
                                message?.success(`${label}上传成功`);
                                setLoading(false);
                            },
                            onError: (error: Error) => {
                                message?.error(error?.message || `${label}上传失败`);
                                setLoading(false);
                            },
                        } as any);
                        if (!task) {
                            setLoading(false);
                            return;
                        }
                        task.promise.then((result: any) => {
                            if (!result) {
                                setLoading(false);
                            }
                        }).catch(() => setLoading(false));
                    };
                    const removeFile = (event: React.MouseEvent) => {
                        event.stopPropagation();
                        modal.confirm({
                            title: "确认删除",
                            content: `确定要移除${label}吗？`,
                            okText: "确定",
                            cancelText: "取消",
                            mask: {
                                closable: true,
                            },
                            onOk: () => setFieldValue(name as any, ""),
                        } as any);
                    };
                    const previewFile = (event: React.MouseEvent) => {
                        event.stopPropagation();
                        setPreviewImage(value);
                    };
                    return (
                        <>
                            <div className={styles.uploadBox} onClick={() => !value && uploadFile()}>
                                {loading ? (
                                    <Spin/>
                                ) : value ? (
                                    <>
                                        <img src={value} alt={label} className={styles.imagePreview}/>
                                        <div className={styles.uploadMask}>
                                            <Icon type="EyeOutlined" onClick={previewFile}/>
                                            <Icon type="DeleteOutlined" onClick={removeFile}/>
                                        </div>
                                    </>
                                ) : (
                                    <div className={styles.uploadContent}>
                                        <Icon type="PlusOutlined"/>
                                        <div className={styles.uploadText}>{uploadText}</div>
                                    </div>
                                )}
                            </div>
                            <Image
                                src={previewImage}
                                style={{display: "none"}}
                                preview={{
                                    open: !!previewImage,
                                    src: previewImage,
                                    onOpenChange: (open: boolean) => {
                                        if (!open) {
                                            setPreviewImage("");
                                        }
                                    },
                                }}
                            />
                        </>
                    );
                }}
            </Form.Item>
        </Form.Item>
    );
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
                              maxWidth,
                              format = "rgb",
                              accept,
                              sceneId,
                              uploadText,
                          }: FormFieldProps): React.ReactNode => {
    const theme = useTheme();
    const inputAction = type === "select" ? "请选择" : type === "upload" ? "请上传" : "请输入";
    const itemRules = rules || (required ? [{required: true, message: `${inputAction}${label}`}] : undefined);
    const itemTooltip = tooltip || label;
    const itemStyle = {maxWidth: maxWidth ?? (theme?.isCompactTheme?.() ? 300 : 350), width: "100%"};
    if (type === "upload") {
        return (
            <UploadField name={name} label={label} tooltip={itemTooltip} rules={itemRules} style={itemStyle}
                         accept={accept} sceneId={sceneId} uploadText={uploadText}/>
        );
    }
    if (type === "textarea") {
        return (
            <Form.Item name={name as any} label={label} tooltip={itemTooltip} rules={itemRules} style={itemStyle}>
                <Input.TextArea placeholder={placeholder || `请输入${label}`} rows={rows}/>
            </Form.Item>
        );
    }
    if (type === "number") {
        return (
            <Form.Item name={name as any} label={label} tooltip={itemTooltip} rules={itemRules} style={itemStyle}>
                <InputNumber min={min} precision={precision} style={{width: "100%"}}
                             placeholder={placeholder || `请输入${label}`}/>
            </Form.Item>
        );
    }
    if (type === "select") {
        return (
            <Form.Item name={name as any} label={label} tooltip={itemTooltip} rules={itemRules} style={itemStyle}>
                <Select placeholder={placeholder || `请选择${label}`} options={options || []}/>
            </Form.Item>
        );
    }
    if (type === "color") {
        return (
            <Form.Item label={label} tooltip={itemTooltip} style={itemStyle}>
                <Flex gap={8} align="center" wrap>
                    <Form.Item name={name as any} rules={itemRules} noStyle
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
        <Form.Item name={name as any} label={label} tooltip={itemTooltip} rules={itemRules} style={itemStyle}>
            <Input placeholder={placeholder || `请输入${label}`}/>
        </Form.Item>
    );
};
