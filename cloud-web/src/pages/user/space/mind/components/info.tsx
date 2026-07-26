import React from "react";
import type {FormInstance} from "antd";
import {Form, Input, Modal, Select, Spin} from "antd";
import {Title} from "sinking-antd";
import type {MindDocument, MindFormValues} from "../types";

type CategoryOption = {
    label: string;
    value: string;
};

export type InfoProps = {
    styles: Record<string, string>;
    form: FormInstance<MindFormValues>;
    open: boolean;
    editingDocument?: MindDocument;
    submitting: boolean;
    categoryOptions: CategoryOption[];
    categoryLoading: boolean;
    onCategoryPopupScroll: React.UIEventHandler<HTMLDivElement>;
    onSave: () => void | Promise<void>;
    onCancel: () => void;
};

const Info: React.FC<InfoProps> = ({
    styles,
    form,
    open,
    editingDocument,
    submitting,
    categoryOptions,
    categoryLoading,
    onCategoryPopupScroll,
    onSave,
    onCancel,
}) => (
    <Modal
        title={<Title>{editingDocument ? "编辑信息" : "新建导图"}</Title>}
        width={420}
        open={open}
        confirmLoading={submitting}
        onOk={() => void onSave()}
        onCancel={onCancel}
        afterClose={() => form.resetFields()}
        okText={editingDocument ? "保存" : "创建"}
        destroyOnHidden
        forceRender
    >
        <Form form={form} layout="vertical" className={styles.mindForm} requiredMark={false}>
            <Form.Item
                name="title"
                label="导图名称"
                rules={[
                    {required: true, whitespace: true, message: "请输入导图名称"},
                    {max: 60, message: "导图名称最多 60 个字符"},
                ]}
            >
                <Input maxLength={60} placeholder="请输入导图名称" autoFocus/>
            </Form.Item>
            <Form.Item
                name="categoryId"
                label="所属分类"
                rules={[{required: true, message: "请选择所属分类"}]}
            >
                <Select
                    showSearch
                    optionFilterProp="label"
                    placeholder="请选择分类"
                    options={categoryOptions}
                    onPopupScroll={onCategoryPopupScroll}
                    notFoundContent={categoryLoading ? <Spin size="small"/> : undefined}
                />
            </Form.Item>
            <Form.Item name="description" label="导图描述" rules={[{max: 500, message: "导图描述最多 500 个字符"}]}>
                <Input.TextArea rows={4} maxLength={500} placeholder="补充一段简短描述（可选）"/>
            </Form.Item>
        </Form>
    </Modal>
);

export default Info;
