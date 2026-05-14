import React from 'react';
import {Body, Title} from "sinking-antd";
import {Card, Col, Row} from "antd";

export default (): React.ReactNode => {
    return (
        <Body>
            <Row gutter={[10, 10]}>
                <Col xs={24} sm={24} md={10} lg={7} xl={6} xxl={3} xxxl={2}>
                    <Card style={{minHeight: "100%"}} title={<Title>系统菜单</Title>} variant={"borderless"}>
                        测试
                    </Card>
                </Col>
                <Col xs={24} sm={24} md={14} lg={17} xl={18} xxl={21} xxxl={22}>
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
