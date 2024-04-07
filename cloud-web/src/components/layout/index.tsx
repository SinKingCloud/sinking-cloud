import {Layout as AntdLayout} from 'antd';

const Layout = ({children, ...props}) => {
    return (
        <AntdLayout {...props}>
            {children}
        </AntdLayout>
    );
};

const Sider = ({children, ...props}) => {
    return (
        <AntdLayout.Sider {...props}>
            {children}
        </AntdLayout.Sider>
    );
};

const Header = ({children, ...props}) => {
    return (
        <AntdLayout.Header {...props}>
            {children}
        </AntdLayout.Header>
    );
};

const Content = ({children, ...props}) => {
    return (
        <AntdLayout.Content {...props}>
            {children}
        </AntdLayout.Content>
    );
};

const Footer = ({children, ...props}) => {
    return (
        <AntdLayout.Footer {...props}>
            {children}
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