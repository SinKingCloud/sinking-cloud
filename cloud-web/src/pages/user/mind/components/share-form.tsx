import React from "react";
import type {FormInstance} from "antd";
import {DatePicker, Form, Input, Modal, Segmented, Select} from "antd";
import {Title} from "sinking-antd";
import {SHARE_DISABLE, SHARE_ENABLE} from "../constants";
import type {MindShareRecord, ShareFormValues} from "../types";

export type ShareFormProps = {
    styles: Record<string, string>;
    form: FormInstance<ShareFormValues>;
    open: boolean;
    editingShare?: MindShareRecord;
    submitting: boolean;
    onSave: () => void | Promise<void>;
    onCancel: () => void;
};

const ShareForm: React.FC<ShareFormProps> = ({
    styles,
    form,
    open,
    editingShare,
    submitting,
    onSave,
    onCancel,
}) => {
    const expireType = Form.useWatch("expireType", form);

    return (
        <Modal
            title={<Title>{editingShare ? "编辑分享" : "新建分享"}</Title>}
            width={350}
            zIndex={1100}
            open={open}
            confirmLoading={submitting}
            closable={!submitting}
            mask={{closable: !submitting}}
            keyboard={!submitting}
            cancelButtonProps={{disabled: submitting}}
            onOk={() => void onSave()}
            onCancel={onCancel}
            afterClose={() => form.resetFields()}
            okText="保存"
            destroyOnHidden
            forceRender
        >
            <Form form={form} layout="vertical" className={styles.mindForm} requiredMark={false}>
                {editingShare && (
                    <Form.Item name="status" label="分享状态" rules={[{required: true, message: "请选择分享状态"}]}>
                        <Select options={[
                            {label: "开启", value: SHARE_ENABLE},
                            {label: "关闭", value: SHARE_DISABLE},
                        ]}/>
                    </Form.Item>
                )}
                <Form.Item
                    name="expireType"
                    label="有效期"
                    rules={[{required: true, message: "请选择有效期"}]}
                >
                    <Segmented
                        block
                        size="small"
                        className="share-expire-segmented"
                        options={[
                            {label: "永久", value: "permanent"},
                            {label: "1天", value: "day"},
                            {label: "3天", value: "threeDays"},
                            {label: "7天", value: "week"},
                            {label: "1个月", value: "month"},
                            {label: "自定义", value: "custom"},
                        ]}
                    />
                </Form.Item>
                {expireType === "custom" && (
                    <Form.Item
                        name="expireTime"
                        label="到期时间"
                        preserve={false}
                        getValueProps={() => ({})}
                        getValueFromEvent={(_, value) => value || undefined}
                        rules={[{
                            validator: async (_, value) => {
                                if (!value) throw new Error("请选择到期时间");
                                if (value === editingShare?.expireTime) return;
                                const timestamp = new Date(value.replace(" ", "T")).getTime();
                                if (!Number.isFinite(timestamp) || timestamp <= Date.now()) {
                                    throw new Error("到期时间必须晚于当前时间");
                                }
                            },
                        }]}
                    >
                        <DatePicker
                            showTime
                            format="YYYY-MM-DD HH:mm:ss"
                            placeholder={editingShare?.expireTime
                                ? `当前：${editingShare.expireTime}`
                                : "选择到期时间"}
                            disabledDate={(current) => current.valueOf() < new Date().setHours(0, 0, 0, 0)}
                            style={{width: "100%"}}
                        />
                    </Form.Item>
                )}
                <Form.Item name="remark" label="备注" rules={[{max: 60, message: "备注最多 60 个字符"}]}>
                    <Input maxLength={60} placeholder="用于区分不同分享链接（可选）"/>
                </Form.Item>
            </Form>
        </Modal>
    );
};

export default ShareForm;
