import React, {useEffect, useRef, useState} from "react";
import {Button, Dropdown, Spin} from "antd";
import {
    DeleteOutlined,
    EditOutlined,
    LeftOutlined,
    MoreOutlined,
    PlusOutlined,
    RightOutlined,
} from "@ant-design/icons";
import type {MindCategory} from "../types";
import {getStableColor} from "../utils";

type CategoryProps = {
    styles: Record<string, string>;
    categories: MindCategory[];
    activeCategoryId: string;
    loading: boolean;
    loadFailed: boolean;
    onSelect: (categoryId: string) => void;
    onLoadMore: () => void;
    onCreate: () => void;
    onEdit: (category: MindCategory) => void;
    onDelete: (category: MindCategory) => void;
    onRetry: () => void;
};

const Category: React.FC<CategoryProps> = ({
    styles,
    categories,
    activeCategoryId,
    loading,
    loadFailed,
    onSelect,
    onLoadMore,
    onCreate,
    onEdit,
    onDelete,
    onRetry,
}) => {
    const scrollRef = useRef<HTMLDivElement | null>(null);
    const loadMoreRef = useRef(onLoadMore);
    const [scrollState, setScrollState] = useState({canLeft: false, canRight: false});

    useEffect(() => {
        loadMoreRef.current = onLoadMore;
    }, [onLoadMore]);

    useEffect(() => {
        const element = scrollRef.current;
        if (!element) return;

        const updateScrollState = (loadMore = false) => {
            const maxScrollLeft = Math.max(0, element.scrollWidth - element.clientWidth);
            const nextState = {
                canLeft: element.scrollLeft > 1,
                canRight: element.scrollLeft < maxScrollLeft - 1,
            };
            setScrollState((current) => current.canLeft === nextState.canLeft
                && current.canRight === nextState.canRight ? current : nextState);
            if (loadMore) {
                const horizontal = element.scrollWidth > element.clientWidth + 1;
                const remaining = horizontal
                    ? element.scrollWidth - element.scrollLeft - element.clientWidth
                    : element.scrollHeight - element.scrollTop - element.clientHeight;
                if (remaining <= 96) {
                    loadMoreRef.current();
                }
            }
        };
        const handleWheel = (event: WheelEvent) => {
            const maxScrollLeft = Math.max(0, element.scrollWidth - element.clientWidth);
            if (maxScrollLeft <= 0) return;
            const delta = Math.abs(event.deltaX) > Math.abs(event.deltaY) ? event.deltaX : event.deltaY;
            const canScroll = delta < 0 ? element.scrollLeft > 0 : element.scrollLeft < maxScrollLeft;
            if (!canScroll) return;
            event.preventDefault();
            element.scrollLeft += delta;
        };

        const handleScroll = () => updateScrollState(true);
        updateScrollState();
        element.addEventListener("scroll", handleScroll, {passive: true});
        element.addEventListener("wheel", handleWheel, {passive: false});
        const observer = new ResizeObserver(() => updateScrollState());
        observer.observe(element);
        if (element.firstElementChild) {
            observer.observe(element.firstElementChild);
        }
        return () => {
            element.removeEventListener("scroll", handleScroll);
            element.removeEventListener("wheel", handleWheel);
            observer.disconnect();
        };
    }, [categories]);

    useEffect(() => {
        if (loading) return;
        const frame = requestAnimationFrame(() => {
            const element = scrollRef.current;
            if (!element) return;
            const horizontal = element.scrollWidth > element.clientWidth + 1;
            const hasOverflow = horizontal
                ? element.scrollWidth > element.clientWidth + 1
                : element.scrollHeight > element.clientHeight + 1;
            if (!hasOverflow) {
                loadMoreRef.current();
            }
        });
        return () => cancelAnimationFrame(frame);
    }, [categories.length, loading]);

    const scrollCategories = (direction: -1 | 1) => {
        const element = scrollRef.current;
        if (!element) return;
        element.scrollBy({
            left: direction * Math.max(180, Math.round(element.clientWidth * 0.68)),
            behavior: "smooth",
        });
    };

    return (
        <aside className={styles.categoryPanel} aria-label="导图分类">
            <div className={styles.categoryInner}>
                {loading ? (
                    <div className={styles.categoryLoading}>
                        <Spin/>
                    </div>
                ) : (
                    <>
                        <div className={styles.categoryRail}>
                            <button
                                className={styles.categoryScrollButton}
                                type="button"
                                aria-label="向左浏览分类"
                                disabled={!scrollState.canLeft}
                                onClick={() => scrollCategories(-1)}
                            >
                                <LeftOutlined/>
                            </button>
                            <div ref={scrollRef} className={styles.categoryScroll}>
                                <div className={styles.categoryGrid}>
                                    {categories.map((item) => (
                                        <div
                                            key={item?.id}
                                            className={`${styles.categoryCard} ${activeCategoryId === item?.id ? "active" : ""}`}
                                            style={{"--accent": getStableColor(item?.id || item?.name)} as React.CSSProperties}
                                            onClick={() => onSelect(item?.id)}
                                        >
                                            <div className="row">
                                                <span className="name-wrap">
                                                    <span className="color"/>
                                                    <span className="name">{item?.name}</span>
                                                </span>
                                                <span className="right">
                                                    {!item?.preset && (
                                                        <Dropdown
                                                            trigger={["click"]}
                                                            placement="bottomRight"
                                                            arrow={{pointAtCenter: true}}
                                                            classNames={{root: styles.actionDropdown}}
                                                            menu={{
                                                                items: [
                                                                    {
                                                                        key: "edit",
                                                                        icon: <EditOutlined/>,
                                                                        label: "修改分类",
                                                                        onClick: ({domEvent}: any) => {
                                                                            domEvent.stopPropagation();
                                                                            onEdit(item);
                                                                        },
                                                                    },
                                                                    {
                                                                        key: "delete",
                                                                        icon: <DeleteOutlined/>,
                                                                        label: "删除分类",
                                                                        danger: true,
                                                                        onClick: ({domEvent}: any) => {
                                                                            domEvent.stopPropagation();
                                                                            onDelete(item);
                                                                        },
                                                                    },
                                                                ] as any,
                                                            }}
                                                        >
                                                            <span
                                                                className="action"
                                                                onClick={(event) => event.stopPropagation()}
                                                            >
                                                                <Button
                                                                    type="text"
                                                                    size="small"
                                                                    icon={<MoreOutlined/>}
                                                                    aria-label={`${item?.name}分类操作`}
                                                                />
                                                            </span>
                                                        </Dropdown>
                                                    )}
                                                </span>
                                            </div>
                                        </div>
                                    ))}
                                    <button
                                        className={`${styles.categoryAddCard} mobile`}
                                        type="button"
                                        onClick={onCreate}
                                    >
                                        <span className="icon"><PlusOutlined/></span>
                                        <span>添加分类</span>
                                    </button>
                                    {loadFailed && (
                                        <Button type="link" size="small" onClick={onRetry}>
                                            重新加载
                                        </Button>
                                    )}
                                </div>
                            </div>
                            <button
                                className={styles.categoryScrollButton}
                                type="button"
                                aria-label="向右浏览分类"
                                disabled={!scrollState.canRight}
                                onClick={() => scrollCategories(1)}
                            >
                                <RightOutlined/>
                            </button>
                        </div>
                        <div className={styles.categoryFooter}>
                            <button className={styles.categoryAddCard} type="button" onClick={onCreate}>
                                <span className="icon"><PlusOutlined/></span>
                                <span>添加分类</span>
                            </button>
                        </div>
                    </>
                )}
            </div>
        </aside>
    );
};

export default Category;
