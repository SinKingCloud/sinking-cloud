import React from "react";
import {createStyles} from "antd-style";
import {Icon} from "@/components";
import {Information, KeFu, Notice, Search} from "@/components/icon";
import {App, Avatar, Dropdown, Input, message, Space} from "antd";
import {DownOutlined} from "@ant-design/icons";
import {useModel} from "@@/exports";
import {outLogin} from "@/service/login";
import {history} from 'umi'

const useStyles = createStyles((token: any): any => {
    return {
        right: {
            float: "right",
            zIndex: 2
        },
        left: {
            float: "left",
            zIndex: 2
        },
        icon: {
            fontSize: "20px",
            marginRight: "15px",
            cursor: "pointer",
        },
        img: {
            marginBottom: "10px",
        },
        drop: {
            marginLeft: "8px",
            fontSize: "14px",
            color: "gray",
            marginRight: "20px",
        },
    };
});

export type HeaderProps = {
    right?: any;//右侧
    left?: any;//左侧
};

const Header: React.FC<HeaderProps> = (props) => {
    const user = useModel("user")
    const {styles: {right, left, icon, img, drop}} = useStyles();
    return <>
        {props?.left && <div className={left}>
            <Input placeholder="搜索菜单" prefix={<Icon type={Search}/>}/>
        </div>}
        {props?.right && <div className={right}>
            <Icon type={Information} className={icon}/>
            <Icon type={KeFu} className={icon}/>
            <Icon type={Notice} className={icon}/>
            <Avatar size="large" className={img} src={user?.web?.avatar}/>
            <Dropdown className={drop} arrow={true} placement={"bottomRight"} menu={{
                items: [{
                    key: '1',
                    label: (
                        <a  onClick={async () => {
                            await outLogin({
                                onSuccess: (r) => {
                                    if(r?.code === 200){
                                        message?.success(r?.message)
                                        localStorage.removeItem("sinking-token")
                                    }
                                },
                                onFail:(r)=>{
                                    message?.error(r?.message)
                                }
                            })
                            history.push("/user/login")
                        }}>退出登录</a>
                    )
                }]
            }}>
                <a onClick={(e) => e.preventDefault()}>
                    <Space>
                        {user?.web?.nick_name || "未设置"}
                        <DownOutlined/>
                    </Space>
                </a>
            </Dropdown>
        </div>}
    </>
}

export default Header;