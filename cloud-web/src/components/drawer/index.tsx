import {Drawer as AntdDrawer} from 'antd';

const Drawer = ({children, ...props}) => {
    return (
        <AntdDrawer {...props}>
            {children}
        </AntdDrawer>
    );
};

export default Drawer;