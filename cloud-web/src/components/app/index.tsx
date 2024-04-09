import {App as AntdApp} from 'antd';

const App = ({children, ...props}) => {
    return (
        <AntdApp {...props}>
            {children}
        </AntdApp>
    );
};

export default AntdApp;