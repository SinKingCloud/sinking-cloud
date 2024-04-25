import React from 'react';
import {Button as AntdButton} from 'antd';

const Button = ({...props}) => {
    return (
        <AntdButton {...props}>
            {props?.children}
        </AntdButton>
    );
};

export default Button;