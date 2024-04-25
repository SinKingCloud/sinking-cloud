import React from 'react';
import {Breadcrumb as AntdBreadcrumb} from 'antd';

const Breadcrumb = ({...props}) => {
    return (
        <AntdBreadcrumb {...props}>
            {props?.children}
        </AntdBreadcrumb>
    );
};

export default Breadcrumb;