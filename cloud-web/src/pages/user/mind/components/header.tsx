import React from "react";
import {Dropdown, Input} from "antd";
import {DownOutlined, SearchOutlined, SortDescendingOutlined,} from "@ant-design/icons";
import {mindSortOptions} from "../constants";
import type {MindSortType} from "../types";

const MindCreateHeroGraphic: React.FC = () => (
    <svg className="hero-graphic" viewBox="0 0 390 154" aria-hidden="true" focusable="false">
        <ellipse className="hero-orbit" cx="239" cy="77" rx="126" ry="61"/>
        <path className="hero-link" d="M219 77C184 77 181 40 143 40M219 77C181 77 180 116 139 116"/>
        <path className="hero-link secondary" d="M267 77C302 77 303 38 341 38M267 77C304 77 306 115 345 115"/>
        <path className="hero-link" d="M243 59V25M243 95V130"/>
        <rect className="hero-node" x="91" y="29" width="57" height="22" rx="7"/>
        <rect className="hero-node" x="79" y="105" width="65" height="22" rx="7"/>
        <rect className="hero-node secondary" x="336" y="27" width="48" height="22" rx="7"/>
        <rect className="hero-node secondary" x="340" y="104" width="45" height="22" rx="7"/>
        <rect className="hero-node" x="219" y="8" width="48" height="20" rx="7"/>
        <rect className="hero-node secondary" x="217" y="127" width="52" height="20" rx="7"/>
        <rect className="hero-core" x="217" y="57" width="52" height="40" rx="12"/>
        <circle className="hero-dot" cx="243" cy="77" r="4"/>
        <circle className="hero-dot" cx="113" cy="40" r="2.5"/>
        <circle className="hero-dot" cx="362" cy="115" r="2.5"/>
    </svg>
);

type HeaderProps = {
    styles: Record<string, string>;
    keyword: string;
    sortType: MindSortType;
    onKeywordChange: (value: string) => void;
    onSearch: (value: string) => void;
    onSortChange: (sortType: MindSortType) => void;
    onCreate: () => void;
};

const Header: React.FC<HeaderProps> = ({
    styles,
    keyword,
    sortType,
    onKeywordChange,
    onSearch,
    onSortChange,
    onCreate,
}) => {
    const activeSortLabel = mindSortOptions.find((item) => item.value === sortType)?.label || "最近更新";

    return (
        <div className={styles.contentPanel}>
            <section className={styles.createHero} aria-labelledby="mind-create-title">
                <div className="hero-copy">
                    <div className="eyebrow">
                        <span className="status-dot"/>
                        <span>Mind Studio</span>
                    </div>
                    <div className="hero-title" id="mind-create-title">让灵感快速长成一张清晰的导图</div>
                    <div className="hero-desc">把零散想法梳理成清晰结构，让每一步都更容易推进。</div>
                </div>
                <div className="hero-visual">
                    <MindCreateHeroGraphic/>
                </div>
                <button className="create-button" type="button" onClick={onCreate}>
                    <span>新建导图</span>
                </button>
            </section>
            <div className={styles.searchCard}>
                <Input
                    allowClear
                    className={styles.searchBox}
                    prefix={<SearchOutlined/>}
                    placeholder="搜索标题或描述"
                    value={keyword}
                    onChange={(event) => onKeywordChange(event.target.value)}
                    onPressEnter={(event) => onSearch(event.currentTarget.value.trim())}
                />
                <Dropdown
                    trigger={["click"]}
                    placement="bottomRight"
                    classNames={{root: styles.sortDropdown}}
                    menu={{
                        selectable: true,
                        selectedKeys: [sortType],
                        items: mindSortOptions.map((item) => ({key: item.value, label: item.label})),
                        onClick: ({key}) => onSortChange(key as MindSortType),
                    }}
                >
                    <button
                        className={styles.sortTrigger}
                        type="button"
                        aria-label={`排序方式：${activeSortLabel}`}
                    >
                        <SortDescendingOutlined className="marker"/>
                        <span className="value">{activeSortLabel}</span>
                        <DownOutlined className="arrow"/>
                    </button>
                </Dropdown>
            </div>
        </div>
    );
};

export default Header;
