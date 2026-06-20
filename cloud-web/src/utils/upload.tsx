export const SYSTEM_FILE_SCENE_ID = 1;

export const getUploadUrl = (result: any) => {
    if (result?.url) {
        return result.url;
    }
    if (typeof result?.data === "string") {
        return result.data;
    }
    return result?.data?.url || result?.data?.path || result?.data?.src;
};
