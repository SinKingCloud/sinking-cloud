import React, {useEffect, useState} from 'react';
import {Avatar, Drawer, List, message, Spin, Switch, Tag, Tooltip} from "antd";
import {changeUinSetting, getProject, getServerList, getUinInfo} from "@/services/open/uin/uin";
import {EditOutlined} from "@ant-design/icons";

export type EditUinProps = {
  uin?: string;
  visible?: boolean;
  onCancel?: (e: any) => void;
};

const EditUin: React.FC<EditUinProps> = (props) => {
  const {uin, visible, onCancel} = props;
  const [u, setU] = useState(uin);
  const [info, setInfo] = useState(null);
  const [speedArr, setSpeedArr] = useState([]);
  const [loading, setLoading] = useState(false);
  /**
   * 挂机地区列表
   */
  const [servers, setServers] = useState({})
  /**
   * 初始化数据
   */
  const init = () => {
    setLoading(true);
    getUinInfo({uin: uin}).then((r) => {
      setLoading(false);
      if (r.code == 200) {
        setInfo(r.data);
      }
    });
  };
  useEffect(() => {
    getProject().then((r) => {
      if (r.code == 200) {
        const getArr = (obj: any): any => {
          let temp = new Array();
          Object.keys(obj).forEach((k: any) => {
            let tempObj = obj[k];
            tempObj.key = k;
            temp.push(tempObj);
          });
          temp = temp.sort((a, b) => {
            return a.sort - b.sort
          });
          return temp;
        }
        setSpeedArr(getArr(Object.assign(r.data?.speed, r.data?.sign)));
      }
    });
    getServerList().then((r) => {
      if (r.code != 200) {
        message.error(r.message || "请求失败").then();
      } else {
        const map = {};
        r.data.list.forEach((item: any) => {
          map[item?.key] = item;
        });
        setServers(map);
      }
    });
  }, []);
  if (uin != "" && u != uin) {
    setU(uin);
    init();
  }
  return (
    <>
      <Drawer title="功能管理" key="EditUin" placement="right" onClose={onCancel} visible={visible}>
        <Spin spinning={loading}>
          <div style={{display: !loading ? "block" : "none", marginTop: "-20px"}}>
            <div style={{textAlign: "center", borderBottom: "1px solid #e1e1e1", padding: "10px"}}>
              <div>
                <Avatar style={{width: "80px", height: "80px"}} src={// @ts-ignore
                  "https://q1.qlogo.cn/g?b=qq&nk=" + info?.uin.uin + "&s=100&t=20190225"} size="large"/>
              </div>
              <div style={{marginTop: "5px"}}>
                <Tag color={"processing"}>{
                  // @ts-ignore
                  info?.uin.uin}</Tag>
              </div>
              <div style={{marginTop: "5px", fontSize: "8px", color: "rgb(140 140 140)"}}>
                <span>挂机地区:<Tooltip title="点击修改挂机地区"><a>{
                  // @ts-ignore
                  servers[info?.uin.server_key]?.name || "默认"}<EditOutlined
                  style={{fontSize: "11px"}}/></a></Tooltip></span>
                <span style={{marginLeft: "30px"}}>挂机时间:<Tooltip title="点击修改挂机时间"><a>{
                  // @ts-ignore
                  (Array(2).join("0") + parseInt(String(info?.uin?.timing / 60))).slice(-2) + ":" + (Array(2).join("0") + parseInt(String(info?.uin?.timing % 60))).slice(-2)
                }<EditOutlined
                  style={{fontSize: "13px"}}/></a></Tooltip></span>
              </div>
            </div>
            {uin != null && <List
              itemLayout="horizontal"
              // @ts-ignore
              dataSource={speedArr}
              renderItem={item => (
                <List.Item actions={[
                  <>
                    {// @ts-ignore
                      item?.edit == "true" &&
                      // @ts-ignore
                      <Switch key={info?.uin.uin + item?.key} defaultChecked={(info?.setting[item?.key] || "0") == "1"}
                        // @ts-ignore
                              onChange={(checked) => {
                                let form = {
                                  uin: uin,
                                  setting: {},
                                };
                                // @ts-ignore
                                form["setting"][item?.key] = checked ? "1" : "0";
                                changeUinSetting(form).then((r) => {
                                  if (r.code != 200) {
                                    message.error(r.message)
                                  }
                                })
                              }}/>}
                    {// @ts-ignore
                      item?.edit == "false" && "附加项目"}
                  </>
                ]}>
                  <List.Item.Meta
                    // @ts-ignore
                    avatar={<Avatar src={item?.logo} style={{borderRadius: "0px"}}/>}
                    // @ts-ignore
                    title={<><a>{item?.name}</a>&nbsp;<Tag
                      color={"success"}>剩余:{// @ts-ignore
                      info?.uin[item?.type == "speed" ? "player_day" : "sign_day"]}天</Tag></>}
                    // @ts-ignore
                    description={<span style={{fontSize: "7px"}}>{item?.describe}</span>}
                  />
                </List.Item>
              )}
            />}
          </div>
        </Spin>
      </Drawer>
    </>
  )
    ;
};

export default EditUin;
