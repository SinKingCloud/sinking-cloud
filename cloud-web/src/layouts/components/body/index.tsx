import React from "react";
import {useModel} from "@@/exports";
import {ConfigProvider, Layout} from "antd";
import {createStyles} from "antd-style";
import {Spin} from "@/components";

const useStyles = createStyles(({token}): any => {
    return {
        body: {
            width: "100%",
            margin: "0 auto",
            lineHeight: "80vh"
        },
    };
});

export type BodyProps = {
    loading?: boolean;//是否加载状态
    style?: any;//样式
    className?: any;//样式名
};
/**
 * 页面主体部分
 * @param props
 * @constructor
 */
const Body: React.FC<BodyProps> = (props) => {
    const {children, loading, style, className} = props;
    const {styles: {body}} = useStyles();
    const systemTheme = useModel("theme");
    return <ConfigProvider theme={systemTheme?.themes}>
        {(loading && <Spin spinning={true} size="large" className={body}></Spin>) ||
            <Layout style={style} className={className}>
                {children}
            </Layout>}
    </ConfigProvider>
}

export default Body