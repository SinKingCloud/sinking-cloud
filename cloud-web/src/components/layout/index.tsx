import React from 'react';
import {Layout as AntdLayout} from 'antd';

const Layout = ({...props}) => {
    return (
        <AntdLayout {...props}>
            {props?.children}
        </AntdLayout>
    );
};

const Sider = ({...props}) => {
    return (
        <AntdLayout.Sider {...props}>
            {props?.children}
        </AntdLayout.Sider>
    );
};

const Header = ({...props}) => {
    return (
        <AntdLayout.Header {...props}>
            {props?.children}
        </AntdLayout.Header>
    );
};

const Content = ({...props}) => {
    return (
        <AntdLayout.Content {...props}>
            {props?.children}
        </AntdLayout.Content>
    );
};

const Footer = ({...props}) => {
    return (
        <AntdLayout.Footer {...props}>
            {props?.children}
        </AntdLayout.Footer>
    );
};

export {
    Layout,
    Sider,
    Header,
    Content,
    Footer
}