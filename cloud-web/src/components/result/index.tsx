import React from 'react';
import {Result as AntdResult} from 'antd';

const Result = ({...props}) => {
    return (
        <AntdResult {...props}>
            {props?.children}
        </AntdResult>
    );
};

export default Result;