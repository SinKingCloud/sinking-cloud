import React from 'react';
import {Space as AntdSpace} from 'antd';

const Space = ({...props}) => {
    return (
        <AntdSpace {...props}>
            {props?.children}
        </AntdSpace>
    );
};

export default Space;