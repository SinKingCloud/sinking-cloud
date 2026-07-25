import React from "react";
import type {FormInstance} from "antd";
import {Form, Input, Modal} from "antd";
import {Title} from "sinking-antd";
import type {CategoryForm as CategoryFormValues, MindCategory} from "../types";

export type CategoryFormProps = {
    styles: Record<string, string>;
    form: FormInstance<CategoryFormValues>;
    open: boolean;
    category?: MindCategory;
    submitting: boolean;
    onSave: () => void | Promise<void>;
    onCancel: () => void;
};

const CategoryForm: React.FC<CategoryFormProps> = ({
    styles,
    form,
    open,
    category,
    submitting,
    onSave,
    onCancel,
}) => (
    <Modal
        title={<Title>{category ? "修改分类" : "新建分类"}</Title>}
        width={350}
        open={open}
        confirmLoading={submitting}
        onOk={() => void onSave()}
        onCancel={onCancel}
        afterClose={() => form.resetFields()}
        okText="保存"
        destroyOnHidden
        forceRender
    >
        <Form form={form} layout="vertical" className={styles.mindForm} requiredMark={false}>
            <Form.Item
                name="name"
                label="分类名称"
                rules={[
                    {required: true, whitespace: true, message: "请输入分类名称"},
                    {max: 24, message: "分类名称最多 24 个字符"},
                ]}
            >
                <Input maxLength={24} placeholder="请输入分类名称" autoFocus/>
            </Form.Item>
        </Form>
    </Modal>
);

export default CategoryForm;
