import React from 'react';
import {ThemeProvider} from 'antd-style';
import {useModel} from "umi";

export type ThemeProps = {
    theme?: any;//主题
    children?: any;//子内容
};

const Theme: React.FC<ThemeProps> = ({...props}) => {
    const systemTheme = useModel("theme");
    return (
        <ThemeProvider theme={props?.theme ? props?.theme : systemTheme?.themes}
                       appearance={props?.theme ? props?.theme?.mode : systemTheme?.themes?.mode}>
            {props?.children}
        </ThemeProvider>
    );
};

export default Theme;