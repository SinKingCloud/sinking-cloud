import React, {useState} from 'react';
import {PageContainer} from "@ant-design/pro-layout";
import ProCard from "@ant-design/pro-card";
import {Alert, Avatar, Button, Col, Form, Image, Input, message, Modal, Row, Space, Tag, Upload} from "antd";
import ProForm, {ProFormDigit, ProFormText} from "@ant-design/pro-form";
import {LockOutlined, UploadOutlined} from "@ant-design/icons";
import {useModel} from "@@/plugin-model/useModel";
import {updateInfo, updatePassword} from "@/services/person/update";
import {getUploadUrl} from "@/services/common/upload";

export default (): React.ReactNode => {
  const {initialState, setInitialState} = useModel('@@initialState');
  // @ts-ignore
  const {currentUser} = initialState;
  const [isUploadAvatarLoading, setIsUploadAvatarLoading] = useState(false);
  /**
   * 获取当前用户信息
   */
  const fetchUserInfo = async () => {
    const userInfo = await initialState?.fetchUserInfo?.();
    if (userInfo) {
      await setInitialState((s) => ({
        ...s,
        currentUser: userInfo,
      }));
    }
  };
  /**
   * 修改密码modal
   */
  const [isPasswordModalVisible, setIsPasswordModalVisible] = useState(false);
  const [isPasswordModalLoading, setIsPasswordModalLoading] = useState(false);
  const [passwordForm] = Form.useForm();

  return (
    <PageContainer header={{title: false, breadcrumb: {},}}>
      <Modal key={"edit"} width={350} destroyOnClose={true} forceRender={true} title="修改密码"
             visible={isPasswordModalVisible} confirmLoading={isPasswordModalLoading}
             onOk={passwordForm.submit} okText={"确 认"} onCancel={() => {
        setIsPasswordModalVisible(false);
        passwordForm.resetFields();
      }}>
        <Form form={passwordForm} name="control-hooks" onFinish={async (values: any) => {
          setIsPasswordModalLoading(true);
          updatePassword(values).then(async (r) => {
            setIsPasswordModalLoading(false);
            if (r.code != 200) {
              message.error(r.message || "请求失败").then()
            } else {
              setIsPasswordModalVisible(false);
              message.success(r.message).then(() => {
                passwordForm.resetFields();
              });
            }
          })
        }} labelAlign="right" labelCol={{span: 6}}
              wrapperCol={{span: 16}}>
          <Form.Item name={"old_password"} label="旧密码" rules={[{required: true, message: '请输入旧密码',},]}>
            <Input placeholder="请输入旧密码"/>
          </Form.Item>
          <Form.Item name={"password"} label="新密码" rules={[{required: true, message: '请输入新密码',},]}>
            <Input placeholder="请输入新密码"/>
          </Form.Item>
        </Form>
      </Modal>
      <Row>
        <Col xs={{span: 24, offset: 0}} lg={{span: 14, offset: 5}}>
          <ProCard title="账户设置" headerBordered>
            <Row>
              <Col lg={{span: 7, offset: 0}} style={{textAlign: "center", margin: "0px auto 20px auto"}}>
                <Space direction="vertical" size={"small"}>
                  <Avatar style={{width: "150px", height: "150px"}} src={
                    <Image
                      style={{width: "150px", height: "150px"}}
                      src={currentUser?.avatar || "https://aliyun_id_photo_bucket.oss.aliyuncs.com/default_handsome.jpg"}/>
                  }/>
                  <Alert message={currentUser?.nick_name} style={{fontWeight: "bolder", fontSize: "20px"}} type="info"/>
                  <div>
                    <Tag color="green">{currentUser?.role_name}</Tag>
                    <Tag color="cyan">ID:{currentUser?.id}</Tag>
                  </div>
                  <div>
                    <Upload name="file" showUploadList={false} action={getUploadUrl()} beforeUpload={(file: any) => {
                      const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png';
                      if (!isJpgOrPng) {
                        message.error('请上传图片格式文件');
                        return false;
                      }
                      const isLt2M = file.size / 1024 / 1024 < 2;
                      if (!isLt2M) {
                        message.error('图片不能大于2MB');
                        return false;
                      }
                      return isJpgOrPng && isLt2M;
                    }} onChange={(info: any) => {
                      if (info.file.status == "uploading") {
                        setIsUploadAvatarLoading(true);
                      } else if (info.file.status === 'done') {
                        const url = info.file.response.data.url;
                        updateInfo({avatar: url}).then(async (r) => {
                          setIsUploadAvatarLoading(false);
                          if (r.code != 200) {
                            message.error(r.message || "请求失败").then()
                          } else {
                            await fetchUserInfo();
                            message.success(r.message).then();
                          }
                        })
                      } else if (info.file.status === 'error') {
                        setIsUploadAvatarLoading(false);
                        message.error(`文件上传失败`);
                      }
                    }}>
                      <Button loading={isUploadAvatarLoading}>
                        <UploadOutlined/> 上传头像
                      </Button>
                    </Upload>
                  </div>
                </Space>
              </Col>
              <Col lg={{span: 13, offset: 1}}>
                <ProForm onFinish={async (values) => {
                  values.phone = values.phone.toString();
                  await updateInfo(values).then(async (r) => {
                    if (r.code != 200) {
                      message.error(r.message || "请求失败").then()
                    } else {
                      await fetchUserInfo();
                      setIsPasswordModalVisible(false);
                      message.success(r.message).then();
                    }
                  })
                }}>
                  <ProFormText width="md" label="账号">
                    <Row gutter={6} wrap={false}>
                      <Col flex={5}>
                        <Input type="text" value={currentUser?.account} disabled={true}/>
                      </Col>
                      <Col flex={7}>
                        <Button onClick={() => {
                          passwordForm.resetFields();
                          setIsPasswordModalVisible(true);
                        }}>
                          <LockOutlined/>修改密码
                        </Button>
                      </Col>
                    </Row>
                  </ProFormText>
                  <ProFormText
                    width="md"
                    name="nick_name"
                    label="昵称"
                    initialValue={currentUser?.nick_name}
                    tooltip="账户昵称"
                    placeholder="请输入账户昵称"
                    rules={[{required: true, message: '请输入账户昵称',},]}
                  />
                  <ProFormDigit
                    width="md"
                    name="phone"
                    label="手机号"
                    initialValue={currentUser?.phone}
                    tooltip="账户绑定手机号"
                    placeholder="请输入绑定手机号"
                    rules={[{required: true, message: '请输入绑定手机号',},]}
                  />
                  <ProFormText
                    width="md"
                    name="email"
                    label="邮箱"
                    initialValue={currentUser?.email}
                    tooltip="账户绑定邮箱"
                    placeholder="请输入绑定邮箱"
                    rules={[{required: true, message: '请输入绑定邮箱!',},]}
                  />
                  <ProFormText
                    width="md"
                    name="login_time"
                    initialValue={currentUser?.login_time}
                    label="登陆时间"
                    tooltip="上次登陆时间"
                    disabled={true}
                  />
                </ProForm>
              </Col>
            </Row>
          </ProCard>
        </Col>
      </Row>
    </PageContainer>
  );
};
