import {useEffect, useMemo, useState} from "react";
import {useSelectedRoutes} from "umi";
import {getControl} from "@/service/admin/page";

export type ControlItem = {
    code?: string;
    name?: string;
    icon?: string;
};

type UseControlResult = {
    code: string;
    controls: ControlItem[];
    loading: boolean;
    error: Error | null;
    hasControl: (controlCode?: string) => boolean;
    getControlItem: (controlCode?: string) => ControlItem | undefined;
    refresh: () => void;
};

const controlCache: Record<string, ControlItem[]> = {};

/**
 * 清理页面控件缓存；传入页面标识时只清理该页面，不传则清空全部缓存。
 */
export function clearControlCache(code?: string) {
    if (code) {
        delete controlCache[code];
        return;
    }
    Object.keys(controlCache).forEach((key) => {
        delete controlCache[key];
    });
}

/**
 * 获取当前页面可用控件，并提供控件存在性判断、控件详情读取和刷新能力。
 */
export function useControl(code?: string): UseControlResult {
    const match = useSelectedRoutes();
    const currentCode = useMemo(() => {
        return code || match?.at(-1)?.route?.name || "";
    }, [code, match]);
    const [controls, setControls] = useState<ControlItem[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<Error | null>(null);
    const [reloadKey, setReloadKey] = useState(0);

    useEffect(() => {
        if (!currentCode) {
            setControls([]);
            setLoading(false);
            setError(null);
            return;
        }

        const cached = controlCache[currentCode];
        if (cached) {
            setControls(cached);
            setLoading(false);
            setError(null);
            return;
        }

        let mounted = true;
        setControls([]);
        setLoading(true);
        setError(null);

        getControl({
            body: {
                code: currentCode,
            },
            onSuccess: (res) => {
                if (!mounted) {
                    return;
                }
                const data = Array.isArray(res?.data) ? res.data : [];
                controlCache[currentCode] = data;
                setControls(data);
            },
            onFail: (res) => {
                if (!mounted) {
                    return;
                }
                setControls([]);
                setError(new Error(res?.message || "获取页面控件失败"));
            }
        }).finally(() => {
            if (mounted) {
                setLoading(false);
            }
        });

        return () => {
            mounted = false;
        };
    }, [currentCode, reloadKey]);

    const controlMap = useMemo(() => {
        const map = new Map<string, ControlItem>();
        controls.forEach((item) => {
            if (item?.code) {
                map.set(item.code, item);
            }
        });
        return map;
    }, [controls]);

    return {
        code: currentCode,
        controls,
        loading,
        error,
        hasControl: (controlCode?: string) => {
            if (!controlCode) {
                return false;
            }
            return controlMap.has(controlCode || "");
        },
        getControlItem: (controlCode?: string) => {
            if (!controlCode) {
                return undefined;
            }
            return controlMap.get(controlCode);
        },
        refresh: () => {
            if (currentCode) {
                delete controlCache[currentCode];
            }
            setReloadKey((value) => value + 1);
        }
    };
}

export default useControl;
