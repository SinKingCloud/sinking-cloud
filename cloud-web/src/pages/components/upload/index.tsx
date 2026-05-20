import {forwardRef, useImperativeHandle} from "react";
import defaultSettings from "@/../config/defaultSettings";
import {getHeaders} from "@/utils/auth";
import type {API} from "../../../../typings";

const DEFAULT_CHUNK_SIZE = 8 * 1024 * 1024;
const MIN_CHUNK_SIZE = 5 * 1024 * 1024;
const DEFAULT_HASH_CHUNK_SIZE = 2 * 1024 * 1024;
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

class Md5 {
    private a = 0x67452301 | 0;
    private b = 0xefcdab89 | 0;
    private c = 0x98badcfe | 0;
    private d = 0x10325476 | 0;
    private buffer = new Uint8Array(64);
    private bufferLength = 0;
    private bytesHashed = 0;

    update(data: Uint8Array): Md5 {
        let position = 0;
        this.bytesHashed += data.length;

        if (this.bufferLength > 0) {
            const available = 64 - this.bufferLength;
            const used = Math.min(available, data.length);
            this.buffer.set(data.subarray(0, used), this.bufferLength);
            this.bufferLength += used;
            position += used;
            if (this.bufferLength === 64) {
                this.processBlock(this.buffer, 0);
                this.bufferLength = 0;
            }
        }

        while (position + 64 <= data.length) {
            this.processBlock(data, position);
            position += 64;
        }

        if (position < data.length) {
            this.buffer.set(data.subarray(position), 0);
            this.bufferLength = data.length - position;
        }

        return this;
    }

    digest(): string {
        const bitLength = this.bytesHashed * 8;
        const paddingLength = this.bufferLength < 56 ? 56 - this.bufferLength : 120 - this.bufferLength;
        const padding = new Uint8Array(paddingLength + 8);
        const low = bitLength >>> 0;
        const high = Math.floor(bitLength / 0x100000000) >>> 0;

        padding[0] = 0x80;
        for (let i = 0; i < 4; i++) {
            padding[paddingLength + i] = (low >>> (8 * i)) & 0xff;
            padding[paddingLength + 4 + i] = (high >>> (8 * i)) & 0xff;
        }
        this.update(padding);

        return [this.a, this.b, this.c, this.d].map(wordToHex).join("");
    }

    private processBlock(data: Uint8Array, offset: number) {
        const x: number[] = [];
        for (let i = 0; i < 16; i++) {
            const j = offset + i * 4;
            x[i] = (data[j] | (data[j + 1] << 8) | (data[j + 2] << 16) | (data[j + 3] << 24)) | 0;
        }

        let a = this.a;
        let b = this.b;
        let c = this.c;
        let d = this.d;

        a = ff(a, b, c, d, x[0], 7, -680876936);
        d = ff(d, a, b, c, x[1], 12, -389564586);
        c = ff(c, d, a, b, x[2], 17, 606105819);
        b = ff(b, c, d, a, x[3], 22, -1044525330);
        a = ff(a, b, c, d, x[4], 7, -176418897);
        d = ff(d, a, b, c, x[5], 12, 1200080426);
        c = ff(c, d, a, b, x[6], 17, -1473231341);
        b = ff(b, c, d, a, x[7], 22, -45705983);
        a = ff(a, b, c, d, x[8], 7, 1770035416);
        d = ff(d, a, b, c, x[9], 12, -1958414417);
        c = ff(c, d, a, b, x[10], 17, -42063);
        b = ff(b, c, d, a, x[11], 22, -1990404162);
        a = ff(a, b, c, d, x[12], 7, 1804603682);
        d = ff(d, a, b, c, x[13], 12, -40341101);
        c = ff(c, d, a, b, x[14], 17, -1502002290);
        b = ff(b, c, d, a, x[15], 22, 1236535329);

        a = gg(a, b, c, d, x[1], 5, -165796510);
        d = gg(d, a, b, c, x[6], 9, -1069501632);
        c = gg(c, d, a, b, x[11], 14, 643717713);
        b = gg(b, c, d, a, x[0], 20, -373897302);
        a = gg(a, b, c, d, x[5], 5, -701558691);
        d = gg(d, a, b, c, x[10], 9, 38016083);
        c = gg(c, d, a, b, x[15], 14, -660478335);
        b = gg(b, c, d, a, x[4], 20, -405537848);
        a = gg(a, b, c, d, x[9], 5, 568446438);
        d = gg(d, a, b, c, x[14], 9, -1019803690);
        c = gg(c, d, a, b, x[3], 14, -187363961);
        b = gg(b, c, d, a, x[8], 20, 1163531501);
        a = gg(a, b, c, d, x[13], 5, -1444681467);
        d = gg(d, a, b, c, x[2], 9, -51403784);
        c = gg(c, d, a, b, x[7], 14, 1735328473);
        b = gg(b, c, d, a, x[12], 20, -1926607734);

        a = hh(a, b, c, d, x[5], 4, -378558);
        d = hh(d, a, b, c, x[8], 11, -2022574463);
        c = hh(c, d, a, b, x[11], 16, 1839030562);
        b = hh(b, c, d, a, x[14], 23, -35309556);
        a = hh(a, b, c, d, x[1], 4, -1530992060);
        d = hh(d, a, b, c, x[4], 11, 1272893353);
        c = hh(c, d, a, b, x[7], 16, -155497632);
        b = hh(b, c, d, a, x[10], 23, -1094730640);
        a = hh(a, b, c, d, x[13], 4, 681279174);
        d = hh(d, a, b, c, x[0], 11, -358537222);
        c = hh(c, d, a, b, x[3], 16, -722521979);
        b = hh(b, c, d, a, x[6], 23, 76029189);
        a = hh(a, b, c, d, x[9], 4, -640364487);
        d = hh(d, a, b, c, x[12], 11, -421815835);
        c = hh(c, d, a, b, x[15], 16, 530742520);
        b = hh(b, c, d, a, x[2], 23, -995338651);

        a = ii(a, b, c, d, x[0], 6, -198630844);
        d = ii(d, a, b, c, x[7], 10, 1126891415);
        c = ii(c, d, a, b, x[14], 15, -1416354905);
        b = ii(b, c, d, a, x[5], 21, -57434055);
        a = ii(a, b, c, d, x[12], 6, 1700485571);
        d = ii(d, a, b, c, x[3], 10, -1894986606);
        c = ii(c, d, a, b, x[10], 15, -1051523);
        b = ii(b, c, d, a, x[1], 21, -2054922799);
        a = ii(a, b, c, d, x[8], 6, 1873313359);
        d = ii(d, a, b, c, x[15], 10, -30611744);
        c = ii(c, d, a, b, x[6], 15, -1560198380);
        b = ii(b, c, d, a, x[13], 21, 1309151649);
        a = ii(a, b, c, d, x[4], 6, -145523070);
        d = ii(d, a, b, c, x[11], 10, -1120210379);
        c = ii(c, d, a, b, x[2], 15, 718787259);
        b = ii(b, c, d, a, x[9], 21, -343485551);

        this.a = add32(this.a, a);
        this.b = add32(this.b, b);
        this.c = add32(this.c, c);
        this.d = add32(this.d, d);
    }
}

const add32 = (a: number, b: number) => (a + b) | 0;

const rotateLeft = (value: number, bits: number) => (value << bits) | (value >>> (32 - bits));

const common = (q: number, a: number, b: number, x: number, s: number, t: number) => {
    return add32(rotateLeft(add32(add32(a, q), add32(x, t)), s), b);
};

const ff = (a: number, b: number, c: number, d: number, x: number, s: number, t: number) => {
    return common((b & c) | (~b & d), a, b, x, s, t);
};

const gg = (a: number, b: number, c: number, d: number, x: number, s: number, t: number) => {
    return common((b & d) | (c & ~d), a, b, x, s, t);
};

const hh = (a: number, b: number, c: number, d: number, x: number, s: number, t: number) => {
    return common(b ^ c ^ d, a, b, x, s, t);
};

const ii = (a: number, b: number, c: number, d: number, x: number, s: number, t: number) => {
    return common(c ^ (b | ~d), a, b, x, s, t);
};

const wordToHex = (word: number) => {
    let output = "";
    for (let i = 0; i < 4; i++) {
        output += `0${((word >>> (i * 8)) & 0xff).toString(16)}`.slice(-2);
    }
    return output;
};

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
    return Object.fromEntries(Object.entries(data).filter(([, value]) => value !== undefined && value !== null && value !== ""));
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

const readBlob = async (blob: Blob): Promise<Uint8Array> => {
    if (typeof blob.arrayBuffer === "function") {
        return new Uint8Array(await blob.arrayBuffer());
    }
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = () => resolve(new Uint8Array(reader.result as ArrayBuffer));
        reader.onerror = () => reject(reader.error || new Error("读取文件失败"));
        reader.readAsArrayBuffer(blob);
    });
};

const yieldToBrowser = () => new Promise<void>((resolve) => {
    if (typeof requestAnimationFrame === "function") {
        requestAnimationFrame(() => resolve());
        return;
    }
    setTimeout(resolve, 0);
});

const getFileMd5 = async (
    file: File,
    onProgress?: (loaded: number, total: number) => void,
    hashChunkSize: number = DEFAULT_HASH_CHUNK_SIZE,
) => {
    const md5 = new Md5();
    const chunkSize = Math.max(1, hashChunkSize);
    let loaded = 0;

    while (loaded < file.size) {
        const end = Math.min(loaded + chunkSize, file.size);
        const chunk = await readBlob(file.slice(loaded, end));
        md5.update(chunk);
        loaded = end;
        onProgress?.(loaded, file.size);
        await yieldToBrowser();
    }

    if (file.size <= 0) {
        onProgress?.(0, 0);
    }

    return md5.digest();
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
        task.snapshot = {
            ...task.snapshot,
            ...patch,
            updateTime: Date.now(),
        };
        this.emit();
    };

    private emit = () => {
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
        const startInfo = createStartInfo("simple", file, sceneId);

        options.onStart?.(startInfo);
        const formData = createFormData();
        appendFormValue(formData, "scene_id", sceneId);
        formData.append("file", file, file.name);

        const response = await requestForm({
            url: options.url || SIMPLE_UPLOAD_URL,
            formData,
            headers: options.headers,
            timeout: options.timeout,
            withCredentials: options.withCredentials,
            signal: getRequestSignal(options),
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
        const chunkSize = normalizeChunkSize(options.chunkSize);
        const startInfo = createStartInfo("chunk", file, sceneId, chunkSize);

        options.onStart?.(startInfo);
        const hashKey = buildChunkHashKey(sceneId, file, chunkSize);
        const cachedHash = options.hash || getCachedFileHash(hashKey);
        const hash = cachedHash || await getFileMd5(file, (loaded, total) => {
            emitProgress(startInfo, "hash", loaded, options, {total});
        }, options.hashChunkSize);
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
                    signal: getRequestSignal(options),
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
                signal: getRequestSignal(options),
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
            const chunk = file.slice(start, end);
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
                signal: getRequestSignal(options),
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
            signal: getRequestSignal(options),
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

/**
 * 唤醒文件选择。
 */
function selectUploadFile(options: SelectUploadFileOptions = {}): Promise<File | undefined> {
    return new Promise((resolve, reject) => {
        if (typeof document === "undefined") {
            reject(new Error("当前环境不支持文件选择"));
            return;
        }

        const input = document.createElement("input");
        let settled = false;

        const cleanup = () => {
            window.removeEventListener("focus", onFocus);
            input.remove();
        };

        const done = (file?: File) => {
            if (settled) {
                return;
            }
            settled = true;
            cleanup();
            resolve(file);
        };

        const onFocus = () => {
            setTimeout(() => {
                if (!settled && (!input.files || input.files.length <= 0)) {
                    done();
                }
            }, 300);
        };

        input.type = "file";
        input.style.display = "none";
        if (options.accept) {
            input.accept = options.accept;
        }
        input.onchange = () => done(input.files?.[0]);
        document.body.appendChild(input);
        setTimeout(() => window.addEventListener("focus", onFocus, {once: true}), 0);
        input.click();
    });
}

async function wakeSimpleUpload(options: WakeSimpleUploadOptions): Promise<UploadResult | undefined> {
    const {accept, file: selectedFile, ...rest} = options;
    let file = selectedFile;
    if (!file) {
        file = await selectUploadFile({accept});
    }
    if (!file) {
        return undefined;
    }
    return simpleUpload({
        ...rest,
        file,
    });
}

async function wakeChunkUpload(options: WakeChunkUploadOptions): Promise<UploadResult | undefined> {
    const {accept, file: selectedFile, ...rest} = options;
    let file = selectedFile;
    if (!file) {
        file = await selectUploadFile({accept});
    }
    if (!file) {
        return undefined;
    }
    return chunkUpload({
        ...rest,
        file,
    });
}

async function wakeUpload(options: WakeUploadOptions): Promise<UploadResult | undefined> {
    const {accept, file: selectedFile, mode, ...rest} = options;
    let file = selectedFile;
    if (!file) {
        file = await selectUploadFile({accept});
    }
    if (!file) {
        return undefined;
    }
    const uploadMode = mode || (file.size > AUTO_CHUNK_SIZE_LIMIT ? "chunk" : "simple");
    if (uploadMode === "chunk") {
        return chunkUpload({
            ...rest,
            file,
        });
    }
    return simpleUpload({
        ...rest,
        file,
    } as any);
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
