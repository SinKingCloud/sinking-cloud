import {Space as AntdSpace} from 'antd';

const Space = ({children, ...props}) => {
    return (
        <AntdSpace {...props}>
            {children}
        </AntdSpace>
    );
};

export default Space;