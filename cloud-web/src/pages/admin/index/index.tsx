import React from "react";
import {Body} from "sinking-antd";
import {Card} from "antd";
import {createStyles} from "antd-style";
import Info from "./components/info";
import Link from "./components/link";
import Overview from "./components/overview";

const useStyles: any = createStyles((): any => ({
    card: {
        ".ant-card-body": {
            padding: 0,
        },
    },
}));

export default (): React.ReactNode => {
    const {styles} = useStyles();

    return (
        <Body>
            <Card className={styles.card}>
                <Info/>
                <Link/>
            </Card>
            <Overview/>
        </Body>
    );
};
