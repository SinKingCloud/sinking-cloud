import React from 'react';
import {Table as AntdTable} from 'antd';

const Table = ({...props}) => {
    return (
        <AntdTable {...props}>
            {props?.children}
        </AntdTable>
    );
};

export default Table;