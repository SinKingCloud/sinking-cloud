import {useState} from "react";
import {theme} from "antd";

// 获取默认风格
const getLightTheme = (): any => {
    return {
        mode: "light",
        token: {
            colorPrimary: "rgba(7,53,237,1)",
            colorInfo: "rgba(7,53,237,1)",
        },
        algorithm: [theme.defaultAlgorithm, theme.compactAlgorithm],
        components: {
            Menu: {
                itemMarginBlock: 0,
                itemMarginInline: 0,
                itemBorderRadius: 0,
                activeBarWidth: 4,
                itemHeight: 45,
                subMenuItemBg: "rgba(255, 255, 255, 0)",
            },
            Layout: {
                headerBg: "white",
            }
        }
    }
}

// 获取暗色风格
const getDarkTheme = (): any => {
    let themes = getLightTheme()
    themes.mode = "dark"
    themes.components.Layout.headerBg = "#141414"
    themes.algorithm = [theme.darkAlgorithm, theme.compactAlgorithm]
    return themes
}

export default () => {
    const [themes, setThemes] = useState<any>(getLightTheme());
    const setLightTheme = () => {
        setThemes(getLightTheme());
    }
    const setDarkTheme = () => {
        setThemes(getDarkTheme());
    }
    const toggle = () => {
        if ((themes?.mode || 'light') == 'light') {
            setDarkTheme();
        } else {
            setLightTheme();
        }
    }
    return {
        themes,
        toggle,
        setLightTheme,
        setDarkTheme
    };
};
