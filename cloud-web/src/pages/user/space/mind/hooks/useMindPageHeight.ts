import type {CSSProperties} from "react";
import {useEffect, useRef, useState} from "react";

export const useMindPageHeight = (dependency?: unknown) => {
    const pageRef = useRef<HTMLDivElement | null>(null);
    const [pageHeight, setPageHeight] = useState<number>();

    useEffect(() => {
        // 按页面真实位置测量当前视口可用高度，避免 Header、面包屑高度变化导致左侧滚动区失准。
        let frame = 0;
        const updateHeight = () => {
            cancelAnimationFrame(frame);
            frame = requestAnimationFrame(() => {
                const page = pageRef.current;
                if (!page) return;
                const top = page.getBoundingClientRect().top;
                const nextHeight = Math.max(320, Math.floor(window.innerHeight - top - 18));
                setPageHeight((value) => value === nextHeight ? value : nextHeight);
            });
        };
        updateHeight();
        window.addEventListener("resize", updateHeight);
        window.addEventListener("orientationchange", updateHeight);
        return () => {
            cancelAnimationFrame(frame);
            window.removeEventListener("resize", updateHeight);
            window.removeEventListener("orientationchange", updateHeight);
        };
    }, [dependency]);

    const pageStyle = pageHeight
        ? {"--mind-page-height": `${pageHeight}px`} as CSSProperties
        : undefined;

    return {pageRef, pageStyle};
};

export default useMindPageHeight;
