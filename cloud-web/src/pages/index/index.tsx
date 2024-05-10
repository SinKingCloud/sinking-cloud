import React from "react";
import {Body} from "@/layouts/components";
import {App, Button, Form} from "antd";

export default () => {
    const {message} = App.useApp();
    return (
        <Body breadCrumb={false}>
            <Button onClick={() => {
                message?.error("succ")
            }}>test</Button>
        </Body>
    );
};