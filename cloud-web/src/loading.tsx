import React from "react";
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

export default () => {
    const {styles: {body}} = useStyles();
    return (
        <Spin size="large" className={body}/>
    );
};