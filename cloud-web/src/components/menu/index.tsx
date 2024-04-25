import React from 'react';
import {Menu as AntdMenu} from 'antd';

const Menu = ({...props}) => {
    return (
        <AntdMenu {...props}>
            {props?.children}
        </AntdMenu>
    );
};

export default Menu;