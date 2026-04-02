import React from "react";
import {Body} from "sinking-antd";
import {Button, Result} from "antd";
import {createStyles} from "antd-style";
import {historyPush} from "@/utils/route";
import {useNavigate} from 'umi';

const useStyles = createStyles((): any => {
    return {
        body: {
            width: "100%",
            height: "100vh",
            paddingTop: "10vh",
        },
    };
});

export default () => {
    const {styles: {body}} = useStyles();
    const navigate = useNavigate();
    return (
        <Body className={body}>
            <Result
                status="500"
                title="500"
                subTitle="服务器连接失败,请返回首页稍后重试"
                extra={<Button type="primary" onClick={() => {
                    if (window?.history?.length > 1) {
                        navigate(-1);
                    } else {
                        historyPush("/");
                    }
                }}>返回上一页</Button>}
            />
        </Body>
    );
};