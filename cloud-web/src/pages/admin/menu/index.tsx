import React from 'react';
import {Body, Title} from "sinking-antd";
import {Card, Col, Row} from "antd";

export default (): React.ReactNode => {
    return (
        <Body>
            <Row gutter={[10, 10]}>
                <Col xs={24} sm={12} md={10} lg={8} xl={6} xxl={5} xxxl={4}>
                    <Card style={{minHeight: "100%"}} title={<Title>菜单系统</Title>} variant={"borderless"}>
                        测试
                    </Card>
                </Col>
                <Col xs={24} sm={12} md={14} lg={16} xl={18} xxl={19} xxxl={20}>
                    <Row gutter={[10, 10]}>
                        <Col span={24}>
                            <Card title={<Title>菜单信息</Title>} variant={"borderless"}>
                                测试
                            </Card>
                        </Col>
                        <Col span={24}>
                            <Card title={<Title>控件管理</Title>} variant={"borderless"}>
                                测试
                            </Card>
                        </Col>
                    </Row>
                </Col>
            </Row>

        </Body>
    );
};
