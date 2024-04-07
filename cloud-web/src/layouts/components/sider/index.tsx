import {useModel, useLocation, history} from 'umi';
import React, {useEffect, useState} from 'react';
import {ConfigProvider, theme} from 'antd';
import {createStyles, useResponsive} from "antd-style";
import {Menu, Layout, Icon} from "@/components";
import {getAllMenuItems} from "@/utils/route";

const useStyles = createStyles(({token}): any => {
    return {
        left: {
            position: "relative",
            background: token?.colorBgContainer,
            height: "100%",
        },
        menu: {
            zIndex: 1,
            overflow: "auto",
            overflowX: "hidden",
            borderRight: "none !important",
            marginBottom: "50px",
            fontWeight: "bold",
            ":focus-visible": {
                outline: "none !important"
            },
            "::-webkit-scrollbar": {
                width: 0
            }
        },
        menuTop: {
            zIndex: 2,
            height: "100px",
            lineHeight: "100px !important",
            width: "100%",
            padding: "10px !important",
            textAlign: "center",
            overflow: "hidden",
            position: "sticky",
            cursor: "pointer",
            top: 0,
        },
        menuBottom: {
            background: token?.colorBgContainer,
            padding: "0px !important",
            borderRadius: "0px",
            height: "50px",
            textAlign: "center",
            lineHeight: "50px !important",
            borderTop: "0.5px solid " + token?.colorBorder + " !important",
            fontWeight: "bolder",
            fontSize: 14,
            color: token?.colorTextSecondary,
            cursor: "pointer",
            overflow: "hidden",
            position: "absolute",
            width: "100%",
            bottom: "0",
            zIndex: 2,
            transition: "background 0.2s !important",
        },
        logoClose: {
            margin: "0px -10px",
        },
    };
});

export type SiderProps = {
    collapsed?: boolean;//菜单展开状态
    menus?: any;//菜单列表
    onMenuClick?: (item: any) => void;//点击菜单回调
    onLogoClick?: () => void;//点击logo回调
    collapsedLogo?: (isLight: boolean) => any;//折叠时logo
    unCollapsedLogo?: (isLight: boolean) => any;//展开时logo
    menuBottomBtnIcon?: string;//底部按钮图标
    menuBottomBtnText?: string;//底部按钮文字
    onMenuBottomBtnClick?: () => void;//点击底部按钮回调
};

const Sider: React.FC<SiderProps> = (props) => {
    const {
        collapsed,
        onMenuClick,
        menus,
        onLogoClick,
        collapsedLogo,
        unCollapsedLogo,
        menuBottomBtnIcon = null,
        menuBottomBtnText = null,
        onMenuBottomBtnClick,
    } = props;
    const {mobile} = useResponsive();
    const systemTheme = useModel("theme");
    const {styles: {left, menu, menuTop, menuBottom, logoClose}} = useStyles();
    const [selectedKeys, setSelectedKeys] = useState([]);
    const location = useLocation();
    useEffect(() => {
        setSelectedKeys([location?.pathname]);
    }, [location]);

    /**
     * 只展开一个子集
     */
    const [openKeys, setOpenKeys] = useState<string[]>([]);
    let parentKey: string[] = [];
    const findParent = (arr: any[], path: string, parent: string[] = []): string[] => {
        let a = parent;
        for (const k in arr) {
            if (parentKey.length == 0) {
                if (arr[k]?.key === path) {
                    parentKey = parent;
                } else {
                    a = [];
                }
                if (arr[k]?.children && arr[k]?.children!.length > 0) {
                    if (findParent(arr[k].children!, path, [arr[k].key, ...parent]).length > 0) {
                        parent = parent.concat(findParent(arr[k].children!, path, [arr[k].key, ...parent]))
                    }
                }
            }
        }
        return a
    }
    const onOpenChange: any = (keys) => {
        const openKey = keys.find((key) => openKeys.indexOf(key) === -1);
        if (openKey) {
            const key = findParent(getAllMenuItems(), openKey, []) || [];
            setOpenKeys([openKey, key, ...parentKey]);
            parentKey = [];
        } else {
            setOpenKeys(keys);
        }
    };
    const {token} = theme?.useToken();
    return (
        <Layout.Layout className={left}>
            <Layout.Header className={menuTop} onClick={() => {
                onLogoClick?.();
            }}>
                {(mobile || (!mobile && !collapsed)) &&
                    unCollapsedLogo?.(systemTheme?.themes?.mode == "light")
                }
                {!mobile && collapsed &&
                    collapsedLogo?.(systemTheme?.themes?.mode == "light")
                }
            </Layout.Header>
            <Layout.Content>
                <ConfigProvider theme={{
                    components: {
                        Menu: {
                            itemSelectedColor: token?.colorPrimaryText,
                            itemColor: token?.colorTextSecondary,
                            itemHoverColor: token?.colorTextSecondary,
                            fontSize: 13,
                        }
                    }
                }}>
                    <Menu selectedKeys={selectedKeys}
                          mode="inline"
                          items={menus} className={menu}
                          openKeys={openKeys}
                          onOpenChange={onOpenChange}
                          onClick={(item) => {
                              history?.push(item?.key);
                              setSelectedKeys([item?.key]);
                              onMenuClick?.(item);
                          }}/>
                </ConfigProvider>
            </Layout.Content>
            {(menuBottomBtnIcon || menuBottomBtnText) &&
                <Layout.Footer className={menuBottom} onClick={onMenuBottomBtnClick}>
                    {menuBottomBtnIcon && <Icon type={menuBottomBtnIcon}/>}
                    {(mobile || (!mobile && !collapsed)) && menuBottomBtnText}
                </Layout.Footer>}
        </Layout.Layout>)
}

export default Sider