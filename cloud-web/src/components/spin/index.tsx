import {Spin as AntdSpin} from 'antd';

const Spin = ({children, ...props}) => {
    return (
        <AntdSpin {...props}>
            {children}
        </AntdSpin>
    );
};

export default Spin;