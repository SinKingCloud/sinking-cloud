import React from 'react';
import {Spin as AntdSpin} from 'antd';

const Spin = ({...props}) => {
    return (
        <AntdSpin {...props}>
            {props?.children}
        </AntdSpin>
    );
};

export default Spin;