import React, {useState} from "react";
import {Modal, Skeleton, Spin} from "antd";
import {ClockCircleOutlined, EyeOutlined, UserOutlined} from "@ant-design/icons";
import {ago} from "@/util/time";
import {getNoticeInfo} from "@/services/common/notice";

export type NoticeInfoProps = {
  id?: number;
  visible?: boolean;
  onClose?: (value?: any) => void;
};

export interface NoticeInfoRef {
  Show?: (id: number | undefined) => void;
}

const NoticeInfo = React.forwardRef<NoticeInfoRef, NoticeInfoProps>((props, ref) => {
  const {id, visible, onClose} = props;
  const [loading, setLoading] = useState(false);
  const [noticeData, setNoticeData] = useState<any>({});
  const handleCancel = () => {
    onClose?.(noticeData);
  };
  // eslint-disable-next-line @typescript-eslint/no-shadow
  const getNoticeData = (id2: number | undefined) => {
    let temp = id
    if ((id2 || 0) > 0) {
      temp = id2
    }
    setLoading(true);
    getNoticeInfo({id: temp}).then((r) => {
      if (r?.code == 200) {
        setNoticeData(r?.data);
      }
      setLoading(false);
    });
  };
  // 将 Test 方法绑定到 ref 对象的 current 属性上
  React.useImperativeHandle(ref, () => ({
    Show: getNoticeData,
  }));
  return (
    <>
      <Modal visible={visible} footer={null} width={550} bodyStyle={{padding: "10px 10px 5px 10px"}}
             onCancel={handleCancel} maskStyle={{backgroundColor: "#72727245 !important"}}>
        <Spin spinning={loading}>
          <Skeleton title={false} loading={loading} active>
            <div style={{margin: "10px"}}>
              <div style={{
                textAlign: "center",
                fontSize: "20px",
                fontWeight: "bolder",
                color: "#545454"
              }}>{noticeData?.title}</div>
              <div style={{textAlign: "center", fontSize: "13px", color: "#8d8d8d", margin: "10px 0"}}>
                <UserOutlined/> 管理员&nbsp;&nbsp;&nbsp;
                <ClockCircleOutlined/> {ago(noticeData?.create_time || '0000-00-00 00:00:00')}&nbsp;&nbsp;&nbsp;
                <EyeOutlined/> {noticeData?.look_num || 0}次
              </div>
              <div style={{margin: "10px 0"}} dangerouslySetInnerHTML={{__html: noticeData?.content || ""}}/>
            </div>
          </Skeleton>
        </Spin>
      </Modal>
    </>
  );
});


export default NoticeInfo;
