import React, {forwardRef, useImperativeHandle, useLayoutEffect, useRef, useState} from "react";
import {App, Form, Select} from "antd";
import {ProModal, ProModalRef, Title} from "sinking-antd";
import {updateUser} from "@/service/admin/user";
import {buildOptions, normalizeIds} from "@/utils/page";
import {AdminControlCode} from "@/utils/permission";
import useEnum from "@/utils/enum";

export type BatchRef = {
    open: (ids?: any[]) => void;
};

type BatchProps = {
    control: any;
    onSuccess?: () => void;
    onClearSelected?: () => void;
};

const Batch = forwardRef<BatchRef, BatchProps>(({
    control,
    onSuccess,
    onClearSelected,
}, ref) => {
    const {message} = App.useApp();
    const [enumData, enumLoading] = useEnum(["user.user"]) as any;
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const [form] = Form.useForm();
    const [initialValues, setInitialValues] = useState<any>({});
    const [formOpen, setFormOpen] = useState(false);
    const mountedForm = useRef<any>(null);
    const [ids, setIds] = useState<number[]>([]);
    const [submitting, setSubmitting] = useState(false);
    const hasUpdateControl = control.hasControl(AdminControlCode.WEB_USER_UPDATE);
    const statusOptions = buildOptions(enumData?.["user.user"]?.status || {}, true);

    const open = (nextIds: any[] = []) => {
        if (!hasUpdateControl) {
            message?.warning("暂无修改权限");
            return;
        }
        const userIds = normalizeIds(nextIds);
        if (userIds.length <= 0) {
            message?.warning("请选择要修改的用户");
            return;
        }
        setIds(userIds);
        setInitialValues({});
        modalRef.current?.show();
    };

    useImperativeHandle(ref, () => ({
        open,
    }));

    const handleSubmit = async (values: any) => {
        if (!hasUpdateControl) {
            message?.warning("暂无修改权限");
            return;
        }
        if (ids.length <= 0) {
            message?.warning("请选择要修改的用户");
            return;
        }
        setSubmitting(true);
        await updateUser({
            body: {
                ids,
                status: values?.status?.toString(),
            },
            onSuccess: (res: any) => {
                modalRef.current?.hide();
                onClearSelected?.();
                onSuccess?.();
                message?.success(res?.message || "操作成功");
            },
            onFail: (res: any) => {
                message?.error(res?.message || "操作失败");
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
            title={<Title>批量编辑</Title>}
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
                    setIds([]);
                    setSubmitting(false);
                },
            } as any}>
            <Form ref={mountedForm} form={form} initialValues={initialValues} clearOnDestroy layout="vertical" onFinish={handleSubmit}>
                <Form.Item name="status" label="用户状态" rules={[{required: true, message: "请选择用户状态"}]}>
                    <Select loading={enumLoading} placeholder="请选择用户状态" options={statusOptions as any}/>
                </Form.Item>
            </Form>
        </ProModal>
    );
});

Batch.displayName = "Batch";

export default Batch;
