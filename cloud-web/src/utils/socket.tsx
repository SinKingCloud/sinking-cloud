import {useEffect, useRef} from "react";
import defaultSettings from "@/../config/defaultSettings";
import {getLoginToken} from "@/utils/auth";
import {getSocketToken as requestSocketToken} from "@/service/common/socket";

export type SocketMessage<T = any> = {
    id?: string;
    type?: string;
    data?: T;
    time?: number;
};

export type SocketListener<T = any> = (message: SocketMessage<T>) => void;

export type SocketOptions<T = any> = {
    onOpen?: () => void;
    onMessage?: SocketListener<T>;
    onError?: (error: Error) => void;
    onClose?: () => void;
};

type SocketOptionsRef = {
    current: SocketOptions | undefined;
};

const DEFAULT_SOCKET_PATH = "/socket";

// 模块级单例状态：多个 useSocket 调用共享同一个真实 WebSocket 连接。
let sharedSocket: WebSocket | undefined;
let connecting = false;
let reconnectTimer: ReturnType<typeof setTimeout> | undefined;
let reconnectAttempts = 0;
let version = 0;
let activeCount = 0;

// 存 ref 而不是直接存回调，保证回调更新时不需要重复注册监听器。
const eventListeners = new Set<SocketOptionsRef>();

// 统一分发事件，不把消息写入 React state，避免 socket 消息触发页面重渲染。
function emit(type: keyof SocketOptions, data?: any) {
    eventListeners.forEach(({current}) => {
        (current?.[type] as ((value?: any) => void) | undefined)?.(data);
    });
}

function clearReconnectTimer() {
    if (reconnectTimer) {
        clearTimeout(reconnectTimer);
        reconnectTimer = undefined;
    }
}

function getSocketUrl(path: string = DEFAULT_SOCKET_PATH) {
    const gateway = defaultSettings?.gateway || "";
    const browserOrigin = typeof window !== "undefined" ? window.location.origin : "http://localhost";
    const target = new URL(gateway || "/", browserOrigin);
    const socketPath = path.startsWith("/") ? path.slice(1) : path;
    const gatewayPath = target.pathname.endsWith("/") ? target.pathname : `${target.pathname}/`;

    target.pathname = `${gatewayPath}${socketPath}`.replace(/\/{2,}/g, "/");
    target.search = "";
    target.hash = "";
    if (target.protocol === "https:") {
        target.protocol = "wss:";
    } else if (target.protocol === "http:") {
        target.protocol = "ws:";
    }
    return target.toString();
}

async function readMessageData(data: any) {
    if (typeof data === "string") {
        return data;
    }
    if (data instanceof ArrayBuffer) {
        return new TextDecoder().decode(data);
    }
    if (typeof Blob !== "undefined" && data instanceof Blob) {
        return data.text();
    }
    return data;
}

function normalizeMessage(payload: any): SocketMessage {
    if (typeof payload !== "string") {
        return {type: "message", data: payload, time: Date.now()};
    }
    try {
        const parsed = JSON.parse(payload);
        if (parsed && typeof parsed === "object") {
            return parsed;
        }
        return {type: "message", data: parsed, time: Date.now()};
    } catch {
        return {type: "message", data: payload, time: Date.now()};
    }
}

function closeSocket() {
    version += 1;
    clearReconnectTimer();
    connecting = false;
    const current = sharedSocket;
    sharedSocket = undefined;
    if (typeof WebSocket !== "undefined" && current && current.readyState !== WebSocket.CLOSED && current.readyState !== WebSocket.CLOSING) {
        current.close(1000, "client close");
    }
}

function scheduleReconnect(connectVersion: number) {
    if (activeCount <= 0 || !getLoginToken()) {
        return;
    }
    // 指数退避重连，最多 30 秒，避免异常时频繁请求 token 和重建连接。
    reconnectAttempts += 1;
    const delay = Math.min(1200 * Math.pow(2, reconnectAttempts - 1), 30000);
    clearReconnectTimer();
    reconnectTimer = setTimeout(() => {
        if (connectVersion === version && activeCount > 0) {
            void connectSocket();
        }
    }, delay);
}

async function connectSocket() {
    if (typeof WebSocket === "undefined") {
        return;
    }
    if (!getLoginToken()) {
        closeSocket();
        return;
    }
    // 已有连接或正在连接时直接复用，保证同一时刻只有一个 socket。
    if (connecting || sharedSocket && sharedSocket.readyState !== WebSocket.CLOSED) {
        return;
    }

    clearReconnectTimer();
    version += 1;
    const connectVersion = version;
    connecting = true;

    try {
        // 后端 token 为一次性授权，创建连接前实时获取。
        const res = await requestSocketToken();
        const token = String(res?.data || "");
        if (connectVersion !== version || activeCount <= 0) {
            return;
        }
        if (res?.code != 200 || !token) {
            emit("onError", new Error(res?.message || "WebSocket token 无效"));
            scheduleReconnect(connectVersion);
            return;
        }

        const socketUrl = getSocketUrl();
        const ws: any = new WebSocket(socketUrl, token);
        ws.binaryType = "arraybuffer";
        sharedSocket = ws;

        ws.onopen = () => {
            if (connectVersion !== version) {
                connecting = false;
                ws.close();
                return;
            }
            connecting = false;
            reconnectAttempts = 0;
            emit("onOpen");
        };

        ws.onmessage = (event) => {
            void readMessageData(event.data).then((payload) => {
                emit("onMessage", normalizeMessage(payload));
            });
        };

        ws.onerror = () => {
            emit("onError", new Error("WebSocket 连接异常"));
        };

        ws.onclose = () => {
            connecting = false;
            if (sharedSocket === ws) {
                sharedSocket = undefined;
            }
            emit("onClose");
            if (connectVersion !== version) {
                return;
            }
            scheduleReconnect(connectVersion);
        };
    } catch (err: any) {
        emit("onError", err instanceof Error ? err : new Error("WebSocket 连接失败"));
        if (connectVersion === version) {
            scheduleReconnect(connectVersion);
        }
    } finally {
        if (connectVersion === version && !sharedSocket) {
            connecting = false;
        }
    }
}

/**
 * 前台用户 WebSocket Hook。
 * 多处调用只会共享一个 WebSocket 连接，每个调用方注册的事件都会被触发。
 */
export function useSocket<T = any>(options?: SocketOptions<T>) {
    const optionsRef = useRef(options);
    optionsRef.current = options;

    useEffect(() => {
        if (!options) {
            return undefined;
        }
        eventListeners.add(optionsRef as SocketOptionsRef);
        return () => {
            eventListeners.delete(optionsRef as SocketOptionsRef);
        };
    }, [!!options]);

    useEffect(() => {
        if (!getLoginToken()) {
            closeSocket();
            return undefined;
        }
        // 记录活跃调用方数量，最后一个调用方卸载时关闭共享连接。
        activeCount += 1;
        void connectSocket();
        return () => {
            activeCount = Math.max(0, activeCount - 1);
            if (activeCount <= 0) {
                closeSocket();
            }
        };
    }, []);
}

export default useSocket;
