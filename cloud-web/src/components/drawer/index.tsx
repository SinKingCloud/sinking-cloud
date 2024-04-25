import React from 'react';
import {Drawer as AntdDrawer} from 'antd';

const Drawer = ({...props}) => {
    return (
        <AntdDrawer {...props}>
            {props?.children}
        </AntdDrawer>
    );
};

export default Drawer;