import {createFromIconfontCN} from '@ant-design/icons';

/**
 * 图标组件(生产环境需使用本地资源)
 */
const Icon = createFromIconfontCN({
    scriptUrl: '//at.alicdn.com/t/c/font_4443635_kbv1qz4vie.js', // 在 iconfont.cn 上生成
});

/**
 * 图标数据
 */
const Doc = "icon-doc";//文档
const FeedBack = "icon-feedback";//信息
const Customer = "icon-customer";//客服
const Notice = "icon-notice";//通知
const Bottom = "icon-bottom";//向下
const Search = "icon-search";//搜索
const Light = "icon-light";//亮色
const Dark = "icon-dark";
const Order = "icon-order";
const Setting = "icon-set";
const PayLog = "icon-pay";
const Right = "icon-right";
const Out = "icon-out";
/**
 * 图标数组
 */
const Icons = [
    Doc,
    FeedBack,
    Customer,
    Notice,
    Bottom,
    Search,
    Light,
    Dark,
    Order,
    Setting,
    PayLog,
    Right,
    Out
]

export {
    Icon,
    Icons,
    Doc,
    FeedBack,
    Customer,
    Notice,
    Bottom,
    Search,
    Light,
    Dark,
    Order,
    Setting,
    PayLog,
    Right,
    Out
}