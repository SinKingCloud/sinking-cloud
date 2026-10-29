import React, {forwardRef, useImperativeHandle, useLayoutEffect, useRef, useState} from "react";
import {App, Form, Input} from "antd";
import {useModel} from "umi";
import {ProModal, ProModalRef, Title} from "sinking-antd";
import {updateAccount} from "@/service/user/person";

export interface AccountProps {
    account?: string;
    onSuccess?: (data: any) => void;
}

export interface AccountRef {
    show: () => void;
}

const Account = forwardRef<AccountRef, AccountProps>((props, ref) => {
    const {account, onSuccess} = props;
    const {message} = App.useApp();
    const user = useModel("user");
    const [form] = Form.useForm();
    const [formOpen, setFormOpen] = useState(false);
    const mountedForm = useRef<any>(null);
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const [loading, setLoading] = useState(false);

    /**
     * 打开账号修改弹窗，已设置账号时阻止重复修改。
     */
    useImperativeHandle(ref, () => ({
        show: () => {
            if (account) {
                message?.warning("账号已设置，无法重复修改");
                return;
            }
            modalRef.current?.show();
        },
    }));

    /**
     * 提交登录账号并同步刷新当前前台用户信息。
     */
    const handleFinish = async (values: any) => {
        setLoading(true);
        await updateAccount({
            body: {
                account: values?.account,
            },
            onSuccess: (res: any) => {
                message?.success(res?.message || "账号修改成功");
                user?.refreshWebUser?.(() => undefined);
                onSuccess?.(values);
                modalRef.current?.hide();
            },
            onFail: (res: any) => message?.error(res?.message || "账号修改失败"),
            onFinally: () => setLoading(false),
        });
    };

    useLayoutEffect(() => {
        if (!formOpen || !mountedForm.current) return;
        form.resetFields();
    }, [form, formOpen]);

    return (
        <ProModal
            ref={modalRef}
            afterOpenChange={setFormOpen}
            title={<Title>修改账号</Title>}
            width={350}
            okText="确认"
            onOk={form.submit}
            onCancel={() => modalRef.current?.hide()}
            modalProps={{
                confirmLoading: loading,
                okButtonProps: {style: {minWidth: 72}},
                cancelButtonProps: {style: {minWidth: 72}},
                style: {
                    top: 100,
                    paddingBottom: 100,
                },
                destroyOnHidden: true,
            } as any}>
            <Form ref={mountedForm} form={form} clearOnDestroy layout="vertical" onFinish={handleFinish}>
                <Form.Item
                    label="登录账号"
                    name="account"
                    rules={[
                        {required: true, message: "请输入账号"},
                        {pattern: /^[a-zA-Z0-9]{5,20}$/, message: "账号为5-20位字母或数字"},
                    ]}>
                    <Input placeholder="请输入新的登录账号" maxLength={20} showCount autoFocus/>
                </Form.Item>
            </Form>
        </ProModal>
    );
});

export default Account;
