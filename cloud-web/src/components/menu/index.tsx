import {Menu as AntdMenu} from 'antd';

const Menu = ({children, ...props}) => {
    return (
        <AntdMenu {...props}>
            {children}
        </AntdMenu>
    );
};

export default Menu;