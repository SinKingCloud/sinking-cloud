import {forwardRef, useImperativeHandle} from "react";
import defaultSettings from "@/../config/defaultSettings";
import {getHeaders} from "@/utils/auth";
import {getFileMd5} from "./hash";
import type {API} from "../../../../typings";

const DEFAULT_CHUNK_SIZE = 8 * 1024 * 1024;
const MIN_CHUNK_SIZE = 5 * 1024 * 1024;
const AUTO_CHUNK_SIZE_LIMIT = 2 * 1024 * 1024;
const SIMPLE_UPLOAD_URL = "/driver/upload/simple";
const CHUNK_UPLOAD_URL = "/driver/upload/chunk";
const CHUNK_RESUME_STORAGE_KEY = "upload-task";
const LEGACY_CHUNK_RESUME_STORAGE_PREFIX = "upload-chunk:";

type SceneId = number | string;
type MaybeSceneId = SceneId | null | undefined;
type UploadFormData = any;

export type UploadMode = "simple" | "chunk";
export type UploadProgressPhase = "hash" | "init" | "status" | "upload" | "finish" | "complete" | "pause";
export type UploadTaskMode = UploadMode | "auto";
export type UploadTaskStatus = "waiting" | "uploading" | "success" | "error" | "aborted";

export interface UploadStartInfo {
    mode: UploadMode;
    file: File;
    sceneId: SceneId;
    total: number;
    chunkSize?: number;
    chunkCount?: number;
}

export interface UploadProgressInfo extends UploadStartInfo {
    phase: UploadProgressPhase;
    loaded: number;
    percent: number;
    hash?: string;
    chunkIndex?: number;
    uploadId?: string;
}

export interface UploadResult extends UploadStartInfo {
    data?: any;
    url?: string;
    hash?: string;
    uploadId?: string;
    response?: API.Response;
}

interface UploadCallbacks {
    onStart?: (info: UploadStartInfo) => void;
    onProgress?: (info: UploadProgressInfo) => void;
    onComplete?: (result: UploadResult) => void;
    onError?: (error: Error) => void;
}

interface BaseUploadOptions extends UploadCallbacks {
    sceneId: MaybeSceneId;
    file?: File | null;
    url?: string | undefined;
    headers?: Record<string, string> | undefined;
    timeout?: number | undefined;
    withCredentials?: boolean | undefined;
    signal?: AbortSignal | undefined;
}

interface SimpleUploadOptions extends BaseUploadOptions {
}

interface ChunkUploadOptions extends BaseUploadOptions {
    hash?: string;
    uploadId?: string;
    resume?: boolean;
    resumeKey?: string;
    chunkSize?: number;
    hashChunkSize?: number;
}

interface SelectUploadFileOptions {
    accept?: string;
    signal?: AbortSignal;
}

interface UploadFileSelection {
    file?: File;
    release: () => void;
}

type WakeSimpleUploadOptions = Omit<SimpleUploadOptions, "file"> & {
    file?: File;
    accept?: string;
};

type WakeChunkUploadOptions = Omit<ChunkUploadOptions, "file"> & {
    file?: File;
    accept?: string;
};

type WakeUploadOptions = WakeChunkUploadOptions & {
    mode?: UploadMode;
};

export interface UploadTaskSnapshot {
    id: string;
    mode: UploadTaskMode;
    status: UploadTaskStatus;
    phase?: UploadProgressPhase;
    fileName?: string;
    fileSize?: number;
    sceneId?: SceneId;
    percent: number;
    loaded: number;
    total: number;
    chunkIndex?: number;
    chunkCount?: number;
    hash?: string;
    uploadId?: string;
    url?: string;
    error?: string;
    createTime: number;
    updateTime: number;
}

export interface UploadTaskHandle {
    id: string;
    promise: Promise<UploadResult | undefined>;
    abort: () => void;
    snapshot: () => UploadTaskSnapshot | undefined;
}

export type UploadTaskListener = (tasks: UploadTaskSnapshot[]) => void;

export interface UploadRef {
    upload?: (options: WakeUploadOptions) => UploadTaskHandle;
    auto?: (options: WakeUploadOptions) => UploadTaskHandle;
    simple?: (options: WakeSimpleUploadOptions) => UploadTaskHandle;
    chunk?: (options: WakeChunkUploadOptions) => UploadTaskHandle;
    abort?: (id: string) => void;
    terminate?: (id: string) => void;
    pause?: (id?: string) => void;
    list?: () => UploadTaskSnapshot[];
    getTask?: (id: string) => UploadTaskSnapshot | undefined;
    setConcurrency?: (value: number) => void;
    getConcurrency?: () => number;
    subscribe?: (listener: UploadTaskListener) => () => void;
}

type RequestFormOptions = {
    url: string;
    formData: UploadFormData;
    headers?: Record<string, string> | undefined;
    timeout?: number | undefined;
    withCredentials?: boolean | undefined;
    signal?: AbortSignal | any;
    onUploadProgress?: ((event: any) => void) | undefined;
};

type TaskResolve = (result: UploadResult | undefined) => void;
type TaskReject = (error: Error) => void;

type InternalUploadTask = {
    snapshot: UploadTaskSnapshot;
    options: WakeUploadOptions;
    controller: AbortController;
    promise: Promise<UploadResult | undefined>;
    resolve: TaskResolve;
    reject: TaskReject;
};

type ResumeCache = {
    version: number;
    tasks: Record<string, {
        uploadId: string;
        updateTime: number;
    }>;
    hashes: Record<string, {
        hash: string;
        updateTime: number;
    }>;
};

let legacyResumeCacheMigrated = false;

const percent = (loaded: number, total: number) => {
    if (!total || total <= 0) {
        return 0;
    }
    return Math.min(100, Number(((loaded / total) * 100).toFixed(2)));
};

const appendFormValue = (formData: UploadFormData, key: string, value?: string | number | null) => {
    if (value === undefined || value === null) {
        return;
    }
    formData.append(key, value.toString());
};

const createFormData = (): UploadFormData => {
    return new FormData();
};

const buildUrl = (url: string) => {
    if (/^https?:\/\//i.test(url)) {
        return url;
    }
    const prefix = defaultSettings?.gateway || "";
    if (!prefix) {
        return url;
    }
    return `${prefix.replace(/\/$/, "")}/${url.replace(/^\//, "")}`;
};

const buildResultUrl = (url?: string) => {
    if (!url || /^https?:\/\//i.test(url)) {
        return url;
    }
    const nextUrl = url.startsWith("/preview/") ? `/driver${url}` : url;
    return buildUrl(nextUrl);
};

const buildHeaders = (headers?: Record<string, string>) => {
    const data: Record<string, any> = {
        ...getHeaders(),
        ...(headers || {}),
    };
    return Object.fromEntries(Object.entries(data).filter(([key, value]) => {
        return key.toLowerCase() !== "content-type" && value !== undefined && value !== null && value !== "";
    }));
};

const requestForm = (options: RequestFormOptions): Promise<API.Response> => {
    const {
        url,
        formData,
        headers,
        timeout = 0,
        withCredentials = false,
        signal,
        onUploadProgress,
    } = options;

    return new Promise((resolve, reject) => {
        if (signal?.aborted) {
            reject(new Error("上传已取消"));
            return;
        }

        const xhr = new XMLHttpRequest();

        const abort: EventListener = () => {
            xhr.abort();
            cleanup();
            reject(new Error("上传已取消"));
        };

        const cleanup = () => {
            signal?.removeEventListener("abort", abort);
        };

        xhr.open("POST", buildUrl(url), true);
        xhr.timeout = timeout;
        xhr.withCredentials = withCredentials;
        Object.entries(buildHeaders(headers)).forEach(([key, value]) => {
            xhr.setRequestHeader(key, value as string);
        });
        xhr.upload.onprogress = onUploadProgress || null;
        xhr.onload = () => {
            cleanup();
            let response: API.Response;
            try {
                response = JSON.parse(xhr.responseText || "{}");
            } catch (_) {
                response = {
                    code: xhr.status,
                    message: xhr.responseText || xhr.statusText,
                };
            }

            if (xhr.status < 200 || xhr.status >= 300) {
                reject(new Error(response?.message || "上传失败"));
                return;
            }
            if (response?.code !== 200) {
                reject(new Error(response?.message || "上传失败"));
                return;
            }
            resolve(response);
        };
        xhr.onerror = () => {
            cleanup();
            reject(new Error("上传失败"));
        };
        xhr.ontimeout = () => {
            cleanup();
            reject(new Error("上传超时"));
        };
        xhr.onabort = () => {
            cleanup();
            reject(new Error("上传已取消"));
        };
        signal?.addEventListener("abort", abort, {once: true});
        xhr.send(formData);
    });
};

const throwIfUploadCanceled = (signal?: AbortSignal) => {
    if (signal?.aborted) {
        throw new Error("上传已取消");
    }
};

const readBlobBuffer = async (blob: Blob, signal?: AbortSignal): Promise<ArrayBuffer> => {
    try {
        throwIfUploadCanceled(signal);
        let buffer: ArrayBuffer;
        if (typeof blob.arrayBuffer === "function") {
            buffer = await blob.arrayBuffer();
        } else {
            buffer = await new Promise<ArrayBuffer>((resolve, reject) => {
                const reader = new FileReader();
                reader.onload = () => resolve(reader.result as ArrayBuffer);
                reader.onerror = () => reject(reader.error || new Error("读取文件失败"));
                reader.readAsArrayBuffer(blob);
            });
        }
        throwIfUploadCanceled(signal);
        return buffer;
    } catch (error) {
        if (signal?.aborted || (error instanceof Error && error.message === "上传已取消")) {
            throw new Error("上传已取消");
        }
        throw new Error("读取所选文件失败，请确认图库原图已下载完成后重试");
    }
};

// Detach lazy file-provider data before XMLHttpRequest serializes the multipart body.
const materializeBlob = async (blob: Blob, signal?: AbortSignal) => {
    return new Blob([await readBlobBuffer(blob, signal)], {type: blob.type});
};

const getSceneId = (sceneId: MaybeSceneId) => {
    if (sceneId === undefined || sceneId === null || sceneId === "") {
        return undefined;
    }
    return sceneId;
};

const normalizeChunkSize = (chunkSize?: number) => {
    if (!chunkSize || chunkSize <= 0) {
        return DEFAULT_CHUNK_SIZE;
    }
    return Math.max(MIN_CHUNK_SIZE, chunkSize);
};

const createStartInfo = (mode: UploadMode, file: File, sceneId: SceneId, chunkSize?: number): UploadStartInfo => {
    const count = chunkSize ? Math.max(1, Math.ceil(file.size / chunkSize)) : undefined;
    return {
        mode,
        file,
        sceneId,
        total: file.size,
        chunkSize,
        chunkCount: count,
    };
};

const emitProgress = (
    info: UploadStartInfo,
    phase: UploadProgressPhase,
    loaded: number,
    options: UploadCallbacks,
    extra?: Partial<UploadProgressInfo>,
) => {
    options.onProgress?.({
        ...info,
        ...extra,
        phase,
        loaded,
        percent: phase === "complete" ? 100 : percent(loaded, info.total),
    });
};

const buildChunkResumeKey = (
    sceneId: SceneId,
    hash: string,
    file: File,
    chunkSize: number,
    customKey?: string,
) => {
    if (customKey) {
        return customKey;
    }
    const parts = [sceneId, hash, file.name, file.size, chunkSize].map((item) => encodeURIComponent(String(item)));
    return parts.join(":");
};

const buildChunkHashKey = (
    sceneId: SceneId,
    file: File,
    chunkSize: number,
) => {
    return [sceneId, file.name, file.size, chunkSize].map((item) => encodeURIComponent(String(item))).join(":");
};

const migrateLegacyResumeCache = (cache: ResumeCache): ResumeCache => {
    if (legacyResumeCacheMigrated || typeof window === "undefined" || !window.localStorage) {
        return cache;
    }
    legacyResumeCacheMigrated = true;
    try {
        const legacyKeys: string[] = [];
        for (let index = 0; index < window.localStorage.length; index++) {
            const key = window.localStorage.key(index);
            if (key?.startsWith(LEGACY_CHUNK_RESUME_STORAGE_PREFIX)) {
                legacyKeys.push(key as any);
            }
        }
        legacyKeys.forEach((key) => {
            const uploadId = window.localStorage.getItem(key);
            const resumeKey = key.slice(LEGACY_CHUNK_RESUME_STORAGE_PREFIX.length);
            if (uploadId && resumeKey && !cache.tasks[resumeKey]) {
                cache.tasks[resumeKey] = {
                    uploadId,
                    updateTime: Date.now(),
                };
            }
            window.localStorage.removeItem(key);
        });
        if (legacyKeys.length > 0) {
            window.localStorage.setItem(CHUNK_RESUME_STORAGE_KEY, JSON.stringify(cache));
        }
    } catch (_) {
    }
    return cache;
};

const readResumeCache = (): ResumeCache => {
    if (typeof window === "undefined" || !window.localStorage) {
        return {version: 1, tasks: {}, hashes: {}};
    }
    try {
        const text = window.localStorage.getItem(CHUNK_RESUME_STORAGE_KEY);
        if (!text) {
            return migrateLegacyResumeCache({version: 1, tasks: {}, hashes: {}});
        }
        const data = JSON.parse(text);
        return migrateLegacyResumeCache({
            version: Number(data?.version || 1),
            tasks: typeof data?.tasks === "object" && data?.tasks ? data.tasks : {},
            hashes: typeof data?.hashes === "object" && data?.hashes ? data.hashes : {},
        });
    } catch (_) {
        return migrateLegacyResumeCache({version: 1, tasks: {}, hashes: {}});
    }
};

const writeResumeCache = (cache: ResumeCache) => {
    if (typeof window === "undefined" || !window.localStorage) {
        return;
    }
    try {
        window.localStorage.setItem(CHUNK_RESUME_STORAGE_KEY, JSON.stringify(cache));
    } catch (_) {
    }
};

const getResumeUploadId = (key: string) => {
    const task = readResumeCache().tasks[key];
    return typeof task?.uploadId === "string" && task.uploadId ? task.uploadId : undefined;
};

const setResumeUploadId = (key: string, uploadId: string) => {
    if (!uploadId) {
        return;
    }
    const cache = readResumeCache();
    cache.tasks[key] = {
        uploadId,
        updateTime: Date.now(),
    };
    writeResumeCache(cache);
};

const getCachedFileHash = (key: string) => {
    const item = readResumeCache().hashes[key];
    return typeof item?.hash === "string" && item.hash ? item.hash : undefined;
};

const setCachedFileHash = (key: string, hash: string) => {
    if (!hash) {
        return;
    }
    const cache = readResumeCache();
    cache.hashes[key] = {
        hash,
        updateTime: Date.now(),
    };
    writeResumeCache(cache);
};

const removeResumeUploadId = (key: string) => {
    const cache = readResumeCache();
    if (!cache.tasks[key]) {
        return;
    }
    delete cache.tasks[key];
    writeResumeCache(cache);
};

export const clearUploadResumeCache = () => {
    if (typeof window === "undefined" || !window.localStorage) {
        return;
    }
    try {
        window.localStorage.removeItem(CHUNK_RESUME_STORAGE_KEY);
    } catch (_) {
    }
};

const normalizeUploadedParts = (parts: any): number[] => {
    if (!Array.isArray(parts)) {
        return [];
    }
    return [...new Set(parts.map((item) => Number(item)).filter((item) => Number.isInteger(item) && item > 0))];
};

const getPartSize = (partNumber: number, fileSize: number, chunkSize: number) => {
    const start = (partNumber - 1) * chunkSize;
    if (start >= fileSize) {
        return 0;
    }
    return Math.max(0, Math.min(chunkSize, fileSize - start));
};

const sumUploadedBytes = (parts: number[], fileSize: number, chunkSize: number) => {
    return parts.reduce((total, partNumber) => total + getPartSize(partNumber, fileSize, chunkSize), 0);
};

const toError = (error: unknown) => error instanceof Error ? error : new Error((error as any)?.message || "上传失败");

const isUploadCanceled = (error: unknown) => toError(error).message === "上传已取消";

const rejectUploadError = (error: unknown, onError?: (error: Error) => void): Promise<never> => {
    const nextError = toError(error);
    onError?.(nextError);
    return Promise.reject(nextError);
};

const getRequestSignal = (options: BaseUploadOptions) => options.signal;

const createTaskId = () => {
    const random = Math.random().toString(36).slice(2, 10);
    return `upload_${Date.now()}_${random}`;
};

const cloneTask = (task: InternalUploadTask): UploadTaskSnapshot => ({...task.snapshot});

class UploadManager {
    private tasks = new Map<string, InternalUploadTask>();
    private queue: string[] = [];
    private activeCount = 0;
    private concurrency = 2;
    private listeners = new Set<UploadTaskListener>();

    list = () => Array.from(this.tasks.values()).map(cloneTask);

    getTask = (id: string) => {
        const task = this.tasks.get(id);
        return task ? cloneTask(task) : undefined;
    };

    subscribe = (listener: UploadTaskListener) => {
        this.listeners.add(listener);
        listener(this.list());
        return () => {
            this.listeners.delete(listener);
        };
    };

    setConcurrency = (value: number) => {
        this.concurrency = Math.max(1, Math.floor(Number(value) || 1));
        this.runQueue();
    };

    getConcurrency = () => this.concurrency;

    auto = (options: WakeUploadOptions) => this.createTask({...options, mode: options.mode});

    upload = (options: WakeUploadOptions) => this.auto(options);

    simple = (options: WakeSimpleUploadOptions) => this.createTask({
        ...options,
        mode: "simple",
    });

    chunk = (options: WakeChunkUploadOptions) => this.createTask({
        ...options,
        mode: "chunk",
    });

    abort = (id: string) => {
        const task = this.tasks.get(id);
        if (!task || task.snapshot.status === "success" || task.snapshot.status === "error" || task.snapshot.status === "aborted") {
            return;
        }
        task.controller.abort();
        this.queue = this.queue.filter((item) => item !== id);
        this.updateTask(task, {
            status: "aborted",
            phase: "pause",
            error: "上传已取消",
        });
        task.resolve(undefined);
        this.runQueue();
    };

    terminate = this.abort;

    pause = (id?: string) => {
        if (id) {
            this.abort(id);
            return;
        }
        Array.from(this.tasks.values()).forEach((task) => {
            if (task.snapshot.status === "uploading") {
                this.abort(task.snapshot.id);
            }
        });
    };

    private createTask = (options: WakeUploadOptions): UploadTaskHandle => {
        const id = createTaskId();
        const now = Date.now();
        let resolveTask: TaskResolve = () => undefined;
        let rejectTask: TaskReject = () => undefined;
        const promise = new Promise<UploadResult | undefined>((resolve, reject) => {
            resolveTask = resolve;
            rejectTask = reject;
        });
        const file = options.file;
        const task: InternalUploadTask = {
            options,
            controller: new AbortController(),
            promise,
            resolve: resolveTask,
            reject: rejectTask,
            snapshot: {
                id,
                mode: options.mode || "auto",
                status: "waiting",
                fileName: file?.name,
                fileSize: file?.size,
                sceneId: options.sceneId === undefined || options.sceneId === null ? undefined : options.sceneId,
                percent: 0,
                loaded: 0,
                total: file?.size || 0,
                createTime: now,
                updateTime: now,
            },
        };
        this.tasks.set(id, task);
        this.queue.push(id);
        this.emit();
        this.runQueue();
        return {
            id,
            promise,
            abort: () => this.abort(id),
            snapshot: () => this.getTask(id),
        };
    };

    private runQueue = () => {
        while (this.activeCount < this.concurrency && this.queue.length > 0) {
            const id = this.queue.shift();
            if (!id) {
                continue;
            }
            const task = this.tasks.get(id);
            if (!task || task.snapshot.status !== "waiting") {
                continue;
            }
            this.runTask(task);
        }
    };

    private runTask = (task: InternalUploadTask) => {
        this.activeCount += 1;
        this.updateTask(task, {status: "uploading"});
        const options = this.wrapCallbacks(task);
        task.snapshot.mode = options.mode || "auto";
        const runner = options.mode === "chunk" ? wakeChunkUpload : options.mode === "simple" ? wakeSimpleUpload : wakeUpload;
        runner({
            ...options,
            signal: task.controller.signal,
        } as any).then((result) => {
            if (!result) {
                this.updateTask(task, {
                    status: "aborted",
                    phase: "pause",
                    error: "未选择文件",
                });
                task.resolve(undefined);
                return;
            }
            if (task.snapshot.status !== "aborted") {
                this.updateTask(task, {
                    status: "success",
                    phase: "complete",
                    percent: 100,
                    loaded: result?.total || task.snapshot.total,
                    total: result?.total || task.snapshot.total,
                    hash: result?.hash || task.snapshot.hash,
                    uploadId: result?.uploadId,
                    url: result?.url,
                });
                task.resolve(result);
            }
        }).catch((error) => {
            const nextError = toError(error);
            if (isUploadCanceled(nextError)) {
                this.updateTask(task, {
                    status: "aborted",
                    phase: "pause",
                    error: nextError.message,
                });
                task.resolve(undefined);
                return;
            }
            this.updateTask(task, {
                status: "error",
                error: nextError.message,
            });
            task.reject(nextError);
        }).finally(() => {
            this.activeCount = Math.max(0, this.activeCount - 1);
            this.runQueue();
        });
    };

    private wrapCallbacks = (task: InternalUploadTask): WakeUploadOptions => {
        const options = task.options;
        return {
            ...options,
            onStart: (info) => {
                this.updateTask(task, {
                    fileName: info.file?.name,
                    fileSize: info.file?.size,
                    sceneId: info.sceneId,
                    total: info.total,
                    chunkCount: info.chunkCount,
                    mode: info.mode,
                });
                options.onStart?.(info);
            },
            onProgress: (info) => {
                this.updateTask(task, {
                    phase: info.phase,
                    loaded: info.loaded,
                    total: info.total,
                    percent: info.percent,
                    chunkIndex: info.chunkIndex,
                    chunkCount: info.chunkCount,
                    hash: info.hash,
                    uploadId: info.uploadId,
                });
                options.onProgress?.(info);
            },
            onComplete: (result) => {
                this.updateTask(task, {
                    status: "success",
                    phase: "complete",
                    percent: 100,
                    loaded: result.total,
                    total: result.total,
                    hash: result.hash,
                    uploadId: result.uploadId,
                    url: result.url,
                });
                options.onComplete?.(result);
            },
            onError: (error) => {
                options.onError?.(error);
            },
        };
    };

    private updateTask = (task: InternalUploadTask, patch: Partial<UploadTaskSnapshot>) => {
        if ((Object.keys(patch) as (keyof UploadTaskSnapshot)[]).every((key) => Object.is(task.snapshot[key], patch[key]))) {
            return;
        }
        task.snapshot = {
            ...task.snapshot,
            ...patch,
            updateTime: Date.now(),
        };
        this.emit();
    };

    private emit = () => {
        if (this.listeners.size === 0) {
            return;
        }
        const tasks = this.list();
        this.listeners.forEach((listener) => listener(tasks));
    };
}

export const upload = new UploadManager();

/**
 * 简单上传。
 */
async function simpleUpload(options: SimpleUploadOptions): Promise<UploadResult> {
    const sceneId = getSceneId(options.sceneId);
    if (sceneId === undefined) {
        return rejectUploadError(new Error("上传场景不能为空"), options.onError);
    }
    const file = options.file;
    if (!file) {
        return rejectUploadError(new Error("上传文件不能为空"), options.onError);
    }

    try {
        const signal = getRequestSignal(options);
        const startInfo = createStartInfo("simple", file, sceneId);

        options.onStart?.(startInfo);
        const formData = createFormData();
        appendFormValue(formData, "scene_id", sceneId);
        formData.append("file", await materializeBlob(file, signal), file.name);

        const response = await requestForm({
            url: options.url || SIMPLE_UPLOAD_URL,
            formData,
            headers: options.headers,
            timeout: options.timeout,
            withCredentials: options.withCredentials,
            signal,
            onUploadProgress: (event) => {
                emitProgress(
                    startInfo,
                    "upload",
                    Math.min(event.loaded || 0, file.size),
                    options,
                );
            },
        });
        const result: UploadResult = {
            ...startInfo,
            data: response?.data,
            url: buildResultUrl(response?.data as string),
            response,
        };
        emitProgress(startInfo, "complete", file.size, options);
        options.onComplete?.(result);
        return result;
    } catch (error) {
        return rejectUploadError(error, options.onError);
    }
}

/**
 * 分片上传。
 */
async function chunkUpload(options: ChunkUploadOptions): Promise<UploadResult> {
    const sceneId = getSceneId(options.sceneId);
    if (sceneId === undefined) {
        return rejectUploadError(new Error("上传场景不能为空"), options.onError);
    }
    const file = options.file;
    if (!file || file.size <= 0) {
        return rejectUploadError(new Error("上传文件不能为空"), options.onError);
    }

    try {
        const signal = getRequestSignal(options);
        const chunkSize = normalizeChunkSize(options.chunkSize);
        const startInfo = createStartInfo("chunk", file, sceneId, chunkSize);

        options.onStart?.(startInfo);
        const hashKey = buildChunkHashKey(sceneId, file, chunkSize);
        const cachedHash = options.hash || getCachedFileHash(hashKey);
        const hash = cachedHash || await getFileMd5(file, (loaded, total) => {
            emitProgress(startInfo, "hash", loaded, options, {total});
        }, options.hashChunkSize, signal);
        setCachedFileHash(hashKey, hash);
        emitProgress(startInfo, "hash", file.size, options, {hash});

        const resume = options.resume !== false;
        const resumeKey = buildChunkResumeKey(sceneId, hash, file, chunkSize, options.resumeKey);
        let uploadId = options.uploadId || (resume ? getResumeUploadId(resumeKey) : undefined);
        let uploadedParts: number[] = [];
        let uploadedBytes = 0;
        let initResponse: API.Response | undefined;
        let initData: any = {};

        if (uploadId) {
            try {
                emitProgress(startInfo, "status", 0, options, {uploadId});
                const statusFormData = createFormData();
                appendFormValue(statusFormData, "scene_id", sceneId);
                appendFormValue(statusFormData, "step", "status");
                appendFormValue(statusFormData, "hash", hash);
                appendFormValue(statusFormData, "upload_id", uploadId);

                const statusResponse = await requestForm({
                    url: options.url || CHUNK_UPLOAD_URL,
                    formData: statusFormData,
                    headers: options.headers,
                    timeout: options.timeout,
                    withCredentials: options.withCredentials,
                    signal,
                });
                const statusData = statusResponse?.data || {};
                uploadedParts = normalizeUploadedParts(statusData?.uploaded_parts);
                uploadedBytes = sumUploadedBytes(uploadedParts, file.size, chunkSize);
                emitProgress(startInfo, "upload", Math.min(uploadedBytes, file.size), options, {uploadId});
            } catch (error) {
                if (isUploadCanceled(error)) {
                    return rejectUploadError(error);
                }
                uploadId = undefined;
                uploadedParts = [];
                uploadedBytes = 0;
                removeResumeUploadId(resumeKey);
            }
        }

        if (!uploadId) {
            emitProgress(startInfo, "init", 0, options);
            const initFormData = createFormData();
            appendFormValue(initFormData, "scene_id", sceneId);
            appendFormValue(initFormData, "step", "init");
            appendFormValue(initFormData, "hash", hash);
            appendFormValue(initFormData, "name", file.name);
            appendFormValue(initFormData, "total_size", file.size);
            appendFormValue(initFormData, "chunk_size", chunkSize);

            initResponse = await requestForm({
                url: options.url || CHUNK_UPLOAD_URL,
                formData: initFormData,
                headers: options.headers,
                timeout: options.timeout,
                withCredentials: options.withCredentials,
                signal,
            });
            initData = initResponse?.data || {};
        }

        if (initData?.done) {
            const result: UploadResult = {
                ...startInfo,
                data: initData,
                url: buildResultUrl(initData?.url),
                hash,
                uploadId: initData?.upload_id,
                response: initResponse,
            };
            removeResumeUploadId(resumeKey);
            emitProgress(startInfo, "complete", file.size, options);
            options.onComplete?.(result);
            return result;
        }

        if (!uploadId) {
            uploadId = initData?.upload_id;
        }
        if (!uploadId) {
            return rejectUploadError(new Error("上传任务ID不能为空"), options.onError);
        }
        if (resume) {
            setResumeUploadId(resumeKey, uploadId);
        }

        const chunkCount = startInfo.chunkCount || Math.ceil(file.size / chunkSize);
        const uploadedPartSet = new Set(uploadedParts);
        for (let index = 0; index < chunkCount; index++) {
            const partNumber = index + 1;
            const start = index * chunkSize;
            const end = Math.min(start + chunkSize, file.size);
            if (uploadedPartSet.has(partNumber)) {
                emitProgress(startInfo, "upload", Math.min(uploadedBytes, file.size), options, {
                    uploadId,
                    chunkIndex: partNumber,
                });
                continue;
            }
            const chunk = await materializeBlob(file.slice(start, end, file.type), signal);
            const formData = createFormData();
            appendFormValue(formData, "scene_id", sceneId);
            appendFormValue(formData, "step", "upload");
            appendFormValue(formData, "hash", hash);
            appendFormValue(formData, "upload_id", uploadId);
            appendFormValue(formData, "part_number", partNumber);
            formData.append("file", chunk, file.name);

            await requestForm({
                url: options.url || CHUNK_UPLOAD_URL,
                formData,
                headers: options.headers,
                timeout: options.timeout,
                withCredentials: options.withCredentials,
                signal,
                onUploadProgress: (event) => {
                    const loaded = Math.min(uploadedBytes + Math.min(event.loaded || 0, end - start), file.size);
                    emitProgress(startInfo, "upload", loaded, options, {
                        uploadId,
                        chunkIndex: partNumber,
                    });
                },
            });
            uploadedBytes += end - start;
            uploadedPartSet.add(partNumber);
            emitProgress(startInfo, "upload", Math.min(uploadedBytes, file.size), options, {
                uploadId,
                chunkIndex: partNumber,
            });
        }

        emitProgress(startInfo, "finish", file.size, options, {uploadId});
        const finishFormData = createFormData();
        appendFormValue(finishFormData, "scene_id", sceneId);
        appendFormValue(finishFormData, "step", "finish");
        appendFormValue(finishFormData, "hash", hash);
        appendFormValue(finishFormData, "name", file.name);
        appendFormValue(finishFormData, "upload_id", uploadId);

        const finishResponse = await requestForm({
            url: options.url || CHUNK_UPLOAD_URL,
            formData: finishFormData,
            headers: options.headers,
            timeout: options.timeout,
            withCredentials: options.withCredentials,
            signal,
        });
        const finishData = finishResponse?.data || {};
        const result: UploadResult = {
            ...startInfo,
            data: finishData,
            url: buildResultUrl(finishData?.url),
            hash,
            uploadId,
            response: finishResponse,
        };
        removeResumeUploadId(resumeKey);
        emitProgress(startInfo, "complete", file.size, options, {uploadId});
        options.onComplete?.(result);
        return result;
    } catch (error) {
        return rejectUploadError(error, options.onError);
    }
}

/** Keep the picker mounted while consumers read file-provider-backed selections. */
function selectUploadFile(options: SelectUploadFileOptions = {}): Promise<UploadFileSelection> {
    return new Promise((resolve, reject) => {
        if (typeof document === "undefined") {
            reject(new Error("当前环境不支持文件选择"));
            return;
        }

        const input = document.createElement("input");
        let settled = false;
        let released = false;
        let focusTimer: number | undefined;

        const cleanupHandlers = () => {
            window.removeEventListener("focus", onFocus);
            input.removeEventListener("cancel", onCancel);
            options.signal?.removeEventListener("abort", onAbort);
            input.onchange = null;
            if (focusTimer !== undefined) {
                window.clearTimeout(focusTimer);
                focusTimer = undefined;
            }
        };

        const release = () => {
            if (released) {
                return;
            }
            released = true;
            cleanupHandlers();
            input.remove();
        };

        const done = (file?: File) => {
            if (settled) {
                return;
            }
            settled = true;
            cleanupHandlers();
            resolve({file, release});
        };

        const onFocus = () => {
            focusTimer = window.setTimeout(() => {
                if (!settled && (!input.files || input.files.length <= 0)) {
                    done();
                }
            }, 2000);
        };

        const onCancel = () => done();

        const onAbort = () => {
            if (settled) {
                return;
            }
            settled = true;
            release();
            reject(new Error("上传已取消"));
        };

        input.type = "file";
        input.style.display = "none";
        if (options.accept) {
            input.accept = options.accept;
        }
        input.onchange = () => done(input.files?.[0]);
        input.addEventListener("cancel", onCancel, {once: true});
        options.signal?.addEventListener("abort", onAbort, {once: true});
        document.body.appendChild(input);
        if (options.signal?.aborted) {
            onAbort();
            return;
        }
        window.addEventListener("focus", onFocus, {once: true});
        input.click();
    });
}

async function withUploadFile<T>(
    selectedFile: File | undefined,
    options: SelectUploadFileOptions,
    handler: (file: File) => Promise<T>,
): Promise<T | undefined> {
    const selection = selectedFile
        ? {file: selectedFile, release: () => undefined}
        : await selectUploadFile(options);
    try {
        if (!selection.file) {
            return undefined;
        }
        return await handler(selection.file);
    } finally {
        selection.release();
    }
}

async function wakeSimpleUpload(options: WakeSimpleUploadOptions): Promise<UploadResult | undefined> {
    const {accept, file: selectedFile, ...rest} = options;
    return withUploadFile(selectedFile, {accept, signal: rest.signal}, (file) => simpleUpload({...rest, file}));
}

async function wakeChunkUpload(options: WakeChunkUploadOptions): Promise<UploadResult | undefined> {
    const {accept, file: selectedFile, ...rest} = options;
    return withUploadFile(selectedFile, {accept, signal: rest.signal}, (file) => chunkUpload({...rest, file}));
}

async function wakeUpload(options: WakeUploadOptions): Promise<UploadResult | undefined> {
    const {accept, file: selectedFile, mode, ...rest} = options;
    return withUploadFile(selectedFile, {accept, signal: rest.signal}, async (file) => {
        const uploadMode = mode || (file.size > AUTO_CHUNK_SIZE_LIMIT ? "chunk" : "simple");
        if (uploadMode === "chunk") {
            return chunkUpload({...rest, file});
        }
        return simpleUpload({...rest, file} as any);
    });
}

export const Upload = forwardRef<UploadRef>((_, ref): any => {
    useImperativeHandle(ref, () => ({
        upload: upload.upload,
        auto: upload.auto,
        simple: upload.simple,
        chunk: upload.chunk,
        abort: upload.abort,
        terminate: upload.terminate,
        pause: upload.pause,
        list: upload.list,
        getTask: upload.getTask,
        setConcurrency: upload.setConcurrency,
        getConcurrency: upload.getConcurrency,
        subscribe: upload.subscribe,
    }));
    return null;
});
