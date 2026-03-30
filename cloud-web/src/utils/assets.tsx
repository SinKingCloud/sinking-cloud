import Settings from "@/../config/defaultSettings";

/**
 * 处理资源路径，正确拼接 basePath
 */
export function getPath(path: string) {
    const basePath = (Settings?.basePath || '').replace(/\/$/, ''); // 去掉结尾的斜杠
    const srcPath = path.replace(/^\//, ''); // 去掉开头的斜杠
    return basePath ? `${basePath}/${srcPath}` : `/${srcPath}`;
}