import {Breadcrumb as AntdBreadcrumb} from 'antd';

const Breadcrumb = ({children, ...props}) => {
    return (
        <AntdBreadcrumb {...props}>
            {children}
        </AntdBreadcrumb>
    );
};

export default Breadcrumb;