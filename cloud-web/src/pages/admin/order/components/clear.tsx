import React, {forwardRef, useImperativeHandle, useLayoutEffect, useRef, useState} from "react";
import {App, Form, InputNumber} from "antd";
import {ProModal, ProModalRef, Title} from "sinking-antd";
import {clearOrder} from "@/service/admin/pay";

export type ClearRef = {
    open: () => void;
};

type ClearProps = {
    onSuccess?: () => void;
};

const Clear = forwardRef<ClearRef, ClearProps>(({onSuccess}, ref) => {
    const {message} = App.useApp();
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const [form] = Form.useForm();
    const [initialValues, setInitialValues] = useState<any>({});
    const [formOpen, setFormOpen] = useState(false);
    const mountedForm = useRef<any>(null);
    const [submitting, setSubmitting] = useState(false);

    useImperativeHandle(ref, () => ({
        open: () => {
            setInitialValues({day: 30});
            modalRef.current?.show();
        },
    }));

    const submit = async (values: any) => {
        setSubmitting(true);
        await clearOrder({
            body: {day: Number(values?.day || 1)},
            onSuccess: (res: any) => {
                modalRef.current?.hide();
                onSuccess?.();
                message?.success(res?.message || "清理成功");
            },
            onFail: (res: any) => {
                message?.error(res?.message || "清理失败");
            },
        }).finally(() => {
            setSubmitting(false);
        });
    };

    useLayoutEffect(() => {
        if (!formOpen || !mountedForm.current) return;
        form.resetFields();
        form.setFieldsValue(initialValues);
    }, [form, formOpen, initialValues]);

    return (
        <ProModal
            ref={modalRef}
            afterOpenChange={setFormOpen}
            title={<Title>清理订单</Title>}
            onOk={form.submit}
            width={320}
            modalProps={{
                confirmLoading: submitting,
                okText: "提交",
                cancelText: "取消",
                destroyOnHidden: true,
                style: {
                    top: 100,
                    paddingBottom: 100,
                },
                mask: {
                    closable: true,
                },
                afterClose: () => {
                    setSubmitting(false);
                },
            } as any}>
            <Form ref={mountedForm} form={form} initialValues={initialValues} clearOnDestroy layout="vertical" onFinish={submit}>
                <Form.Item name="day" label="保留天数" rules={[{required: true, message: "请输入保留天数"}]}>
                    <InputNumber min={1} precision={0} step={1} style={{width: "100%"}} placeholder="请输入保留天数"/>
                </Form.Item>
            </Form>
        </ProModal>
    );
});

Clear.displayName = "Clear";

export default Clear;
