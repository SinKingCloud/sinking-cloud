import {Result as AntdResult} from 'antd';

const Result = ({children, ...props}) => {
    return (
        <AntdResult {...props}>
            {children}
        </AntdResult>
    );
};

export default Result;