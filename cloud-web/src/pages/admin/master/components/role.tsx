import React, {forwardRef, useImperativeHandle, useRef, useState} from "react";
import {App, Form, Select} from "antd";
import {ProModal, ProModalRef, Title} from "sinking-antd";
import useEnum from "@/utils/enum";
import {updateMasterRole} from "@/service/admin/master";
import {buildOptions, normalizeIds} from "@/utils/page";
import {AdminControlCode} from "@/utils/permission";

export type RoleRef = {
    open: (record: any) => void;
};

type RoleProps = {
    control: any;
    onSuccess?: () => void;
};

const Role = forwardRef<RoleRef, RoleProps>(({control, onSuccess}, ref) => {
    const {message} = App.useApp();
    const [enumData] = useEnum(["admin.role"]) as any;
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const [form] = Form.useForm();
    const [roleUser, setRoleUser] = useState<any>();
    const [submitting, setSubmitting] = useState(false);
    const hasRoleControl = control.hasControl(AdminControlCode.SYSTEM_MASTER_ROLE);
    const roleNames = enumData?.["admin.role"]?.names || {};
    const roleOptions = buildOptions(roleNames, true);

    useImperativeHandle(ref, () => ({
        open: (record: any) => {
            if (!hasRoleControl) {
                message?.warning("暂无分配角色权限");
                return;
            }
            setRoleUser(record || {});
            form.resetFields();
            form.setFieldsValue({
                role_ids: normalizeIds(record?.role_ids || []),
            });
            modalRef.current?.show();
        },
    }));

    const handleSubmit = async (values: any) => {
        if (!hasRoleControl) {
            message?.warning("暂无分配角色权限");
            return;
        }
        setSubmitting(true);
        await updateMasterRole({
            body: {
                id: Number(roleUser?.id),
                role_ids: normalizeIds(values?.role_ids || []),
            },
            onSuccess: (res: any) => {
                modalRef.current?.hide();
                onSuccess?.();
                message?.success(res?.message || "分配成功");
            },
            onFail: (res: any) => {
                message?.error(res?.message || "分配失败");
            },
            onFinally: () => setSubmitting(false),
        });
    };

    return (
        <ProModal
            ref={modalRef}
            title={<Title>编辑角色</Title>}
            onOk={form.submit}
            width={340}
            modalProps={{
                confirmLoading: submitting,
                style: {
                    top: 100,
                    paddingBottom: 100,
                },
                mask: {
                    closable: true,
                },
                okText: "提交",
                cancelText: "取消",
                afterClose: () => {
                    form.resetFields();
                    setRoleUser(undefined);
                    setSubmitting(false);
                },
            } as any}
        >
            <Form form={form} layout="vertical" onFinish={handleSubmit}>
                <Form.Item name="role_ids" label="角色" rules={[{required: true, message: "请选择角色"}]}>
                    <Select mode="multiple" placeholder="请选择角色" options={roleOptions as any}/>
                </Form.Item>
            </Form>
        </ProModal>
    );
});

Role.displayName = "Role";

export default Role;
