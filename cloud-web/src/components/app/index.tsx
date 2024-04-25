import {App as AntdApp} from 'antd';

const App = ({...props}) => {
    return (
        <AntdApp {...props}>
            {props?.children}
        </AntdApp>
    );
};

export default App;