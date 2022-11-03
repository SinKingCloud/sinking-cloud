import React, {useEffect, useState} from 'react';
import {Avatar, Badge, Drawer, List, message, Modal, Select, Spin, Switch, Tag, TimePicker, Tooltip} from "antd";
import {
  changeUinArea,
  changeUinSetting,
  changeUinTiming,
  getProject,
  getServerList,
  getUinInfo
} from "@/services/open/uin/uin";
import moment from "moment";
import {ExclamationCircleOutlined, InfoCircleOutlined} from "@ant-design/icons";

export type EditUinProps = {
  uin?: string;
  visible?: boolean;
  onCancel?: (e: any) => void;
};

const EditUin: React.FC<EditUinProps> = (props) => {
  const {uin, visible, onCancel} = props;
  const [u, setU] = useState(uin);
  const [serverKey, setServerKey] = useState("");
  const [timing, setTiming] = useState("00:00");
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
        setServerKey(r.data.uin.server_key);
        setTiming(((Array(2).join("0") + parseInt(String(r.data.uin?.timing / 60))).slice(-2) + ":" + (Array(2).join("0") + parseInt(String(r.data.uin?.timing % 60))).slice(-2)).toString())
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
    getServerList({page_size: 999}).then((r) => {
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
  /**
   * 修改挂机地区
   */
  const changeArea = (uin: string, serverKey: string) => {
    Modal.confirm({
      title: '确定要更改挂机地区吗?',
      icon: <ExclamationCircleOutlined/>,
      content: '更改挂机地区可能会影响挂机稳定',
      okType: 'primary',
      onOk() {
        const key = "changeArea";
        message.loading({content: '正在修改挂机地区', key, duration: 60}).then();
        changeUinArea({uin: uin, server_key: serverKey}).then((r) => {
          if (r.code == 200) {
            setServerKey(serverKey);
            message.success({content: r.message, key, duration: 2}).then();
          } else {
            message.error({content: r.message, key, duration: 2}).then();
          }
        });
      },
    });
  }
  /**
   * 修改挂机时间
   */
  const changeTiming = (uin: string, timing: string) => {
    Modal.confirm({
      title: '确定要修改挂机时间为' + timing + '吗?',
      icon: <ExclamationCircleOutlined/>,
      content: '下次挂机将会在每日' + timing + "时刻挂机",
      okType: 'primary',
      onOk() {
        const key = "changeTiming";
        message.loading({content: '正在修改挂机时间', key, duration: 60}).then();
        let tempArr = timing.split(":");
        changeUinTiming({uin: uin, timing: parseInt(tempArr[0]) * 60 + parseInt(tempArr[1])}).then((r) => {
          if (r.code == 200) {
            setTiming(timing);
            message.success({content: r.message, key, duration: 2}).then();
          } else {
            message.error({content: r.message, key, duration: 2}).then();
          }
        });
      },
    });
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
                {// @ts-ignore
                  info?.uin.uin_status == 0 &&
                  <Tag color={"success"}>正常</Tag>
                }
                {// @ts-ignore
                  info?.uin.uin_status != 0 &&
                  <Tooltip title="需要更新状态后才能继续挂机"><Tag color={"red"}>异常</Tag></Tooltip>
                }
              </div>
              <div style={{marginTop: "5px", fontSize: "8px", color: "rgb(140 140 140)"}}>
                <span>挂机地区:<Tooltip title="点击修改挂机地区">
                   <Select key={"serverKey"} placeholder="请选择" style={{fontSize: "8px"}} value={serverKey}
                           bordered={false} onChange={(e) => {
                     // @ts-ignore
                     changeArea(info?.uin.uin, e);
                   }}>
                    {Object.keys(servers).map((k) => {
                      return <Select.Option key={k} value={k}>{servers[k]?.name || "默认"}</Select.Option>;
                    })}
                    </Select>
                </Tooltip>
                </span>
                <span>挂机时间:<Tooltip title="点击修改挂机时间">
                  <TimePicker key={"timing"} style={{maxWidth: "80px"}} bordered={false} allowClear={false}
                              value={moment(timing, 'HH:mm')}
                              format={'HH:mm'} onChange={(time: any, timeString: string) => {
                    // @ts-ignore
                    changeTiming(info?.uin.uin, timeString);
                  }}/>
                </Tooltip>
                </span>
              </div>
            </div>
            {uin != null && <List
              itemLayout="horizontal"
              // @ts-ignore
              dataSource={speedArr}
              renderItem={item => {
                // @ts-ignore
                let days = info?.uin[item?.type == "speed" ? "player_day" : "sign_day"]
                return <>
                  <List.Item actions={[
                    <>
                      {// @ts-ignore
                        item?.edit == "true" &&
                        // @ts-ignore
                        <Switch key={info?.uin.uin + item?.key}
                          // @ts-ignore
                                defaultChecked={(info?.setting[item?.key] || "0") == "1"}
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
                      avatar={<Badge  color={info?.job_status[item?.key] == "1" ? "green" : (info?.job_status[item?.key] == "2" ? "blue" : "red")}
                        dot>
                        <Avatar src={
                          // @ts-ignore
                          item?.logo} style={{borderRadius: "0px"}}/>
                      </Badge>}
                      // @ts-ignore
                      title={<><a>{item?.name}</a>&nbsp;{days <= 0 &&
                        <Tag color={"red"}>天数不足</Tag>}
                        {days > 0 && <Tag color={"success"}>剩余:{days}天</Tag>}</>}
                      // @ts-ignore
                      description={<><InfoCircleOutlined style={{fontSize: "12px"}}/>&nbsp;<span style={{fontSize: "10px"}}>{item?.describe}</span></>}
                    />
                  </List.Item>
                </>
              }}
            />}
          </div>
        </Spin>
      </Drawer>
    </>
  )
    ;
};

export default EditUin;
