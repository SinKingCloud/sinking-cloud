import React, {forwardRef, useImperativeHandle, useRef, useState} from "react";
import {App, Form, Input} from "antd";
import {useModel} from "umi";
import {ProModal, ProModalRef, Title} from "sinking-antd";
import {updatePersonInfo} from "@/service/admin/person";

export interface InfoProps {
    nickName?: string;
    onSuccess?: (data: any) => void;
}

export interface InfoRef {
    show: () => void;
}

const Info = forwardRef<InfoRef, InfoProps>((props, ref) => {
    const {nickName, onSuccess} = props;
    const {message} = App.useApp();
    const user = useModel("user");
    const [form] = Form.useForm();
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const [loading, setLoading] = useState(false);

    /**
     * 打开昵称编辑弹窗并回填当前昵称。
     */
    useImperativeHandle(ref, () => ({
        show: () => {
            form.setFieldsValue({
                nick_name: nickName || "",
            });
            modalRef.current?.show();
        },
    }));

    /**
     * 保存基础信息并刷新当前后台用户信息。
     */
    const handleFinish = async (values: any) => {
        setLoading(true);
        await updatePersonInfo({
            body: values,
            onSuccess: (res: any) => {
                message?.success(res?.message || "修改成功");
                user?.refreshSystemUser?.(() => undefined);
                onSuccess?.(values);
                modalRef.current?.hide();
            },
            onFail: (res: any) => message?.error(res?.message || "修改失败"),
            onFinally: () => setLoading(false),
        });
    };

    return (
        <ProModal
            ref={modalRef}
            title={<Title>修改信息</Title>}
            width={350}
            okText="确认"
            onOk={form.submit}
            onCancel={() => modalRef.current?.hide()}
            afterClose={() => form.resetFields()}
            modalProps={{
                confirmLoading: loading,
                destroyOnHidden: true,
                forceRender: true,
            } as any}>
            <Form form={form} layout="vertical" onFinish={handleFinish}>
                <Form.Item
                    label="账号昵称"
                    name="nick_name"
                    rules={[
                        {required: true, message: "请输入昵称"},
                        {min: 1, max: 15, message: "昵称长度为1-15个字符"},
                    ]}>
                    <Input placeholder="请输入新的昵称" maxLength={15} showCount autoFocus/>
                </Form.Item>
            </Form>
        </ProModal>
    );
});

export default Info;
