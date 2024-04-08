import {Table as AntdTable} from 'antd';

const Table = ({children, ...props}) => {
    return (
        <AntdTable {...props}>
            {children}
        </AntdTable>
    );
};

export default Table;