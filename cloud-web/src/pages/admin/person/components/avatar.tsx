import React, {forwardRef, useImperativeHandle, useRef, useState} from "react";
import {App, Avatar, Spin} from "antd";
import {Icon} from "sinking-antd";

import {createStyles} from "antd-style";
import {useModel} from "umi";
import {Upload, type UploadRef} from "@/pages/components/upload";
import {updatePersonInfo} from "@/service/admin/person";
import {SYSTEM_FILE_SCENE_ID, getUploadUrl} from "@/utils/upload";

export interface AvatarProps {
    avatar?: string;
    nickName?: string;
    onSuccess?: (avatarUrl: string) => void;
    showAvatar?: boolean;
    size?: number;
}

export interface AvatarRef {
    triggerUpload: () => void;
}

const useStyles: any = createStyles(({token}: any): any => ({
    avatarUpload: {
        position: "relative",
        display: "inline-block",
        cursor: "pointer",
        "&:hover .avatar-overlay": {
            opacity: 1,
        },
    },
    avatarContainer: {
        position: "relative",
        display: "inline-block",
        cursor: "pointer",
    },
    avatarStyle: {
        cursor: "pointer",
        border: `3px solid ${token.colorBgContainer}`,
    },
    avatarOverlay: {
        position: "absolute",
        inset: 0,
        borderRadius: "50%",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        opacity: 0,
        color: "#fff",
        fontSize: 14,
        pointerEvents: "none",
        transition: "opacity 0.25s ease",
        backgroundColor: "rgba(0, 0, 0, 0.52)",
    },
    loadingContainer: {
        display: "inline-block",
    },
}));

const AvatarUpload = forwardRef<AvatarRef, AvatarProps>((props, ref) => {
    const {avatar, nickName, onSuccess, showAvatar = true, size = 34} = props;
    const {message} = App.useApp();
    const {styles} = useStyles();
    const user = useModel("user");
    const uploadRef = useRef<UploadRef>({} as UploadRef);
    const [loading, setLoading] = useState(false);

    /**
     * 选择并上传头像，上传成功后保存到当前后台账号资料。
     */
    const triggerUpload = () => {
        if (loading) {
            return;
        }
        setLoading(true);
        const task = uploadRef.current?.upload?.({
            sceneId: SYSTEM_FILE_SCENE_ID,
            accept: "image/*",
            onComplete: async (result: any) => {
                const url = getUploadUrl(result);
                if (!url) {
                    setLoading(false);
                    message?.warning("头像上传成功，但未返回地址");
                    return;
                }
                await updatePersonInfo({
                    body: {avatar: url},
                    onSuccess: (res: any) => {
                        message?.success(res?.message || "头像已更新");
                        user?.refreshSystemUser?.(() => undefined);
                        onSuccess?.(url);
                    },
                    onFail: (res: any) => message?.error(res?.message || "头像更新失败"),
                    onFinally: () => setLoading(false),
                });
            },
            onError: (error: Error) => {
                setLoading(false);
                message?.error(error?.message || "头像上传失败");
            },
        } as any);
        if (!task) {
            setLoading(false);
            return;
        }
        task.promise.then((result) => {
            if (!result) {
                setLoading(false);
            }
        }).catch(() => setLoading(false));
    };

    /**
     * 暴露头像上传触发方法给设置项的“修改”按钮。
     */
    useImperativeHandle(ref, () => ({
        triggerUpload,
    }));

    return (
        <div className={styles.avatarUpload}>
            <Upload ref={uploadRef}/>
            {showAvatar && (
                <div className={styles.avatarContainer} onClick={triggerUpload}>
                    <Spin spinning={loading}>
                        <Avatar
                            size={size}
                            src={avatar}
                            icon={<Icon type="UserOutlined"/>}
                            className={styles.avatarStyle}>
                            {nickName || "用户"}
                        </Avatar>
                    </Spin>
                    <div className={`${styles.avatarOverlay} avatar-overlay`}>
                        <Icon type="CameraOutlined"/>
                    </div>
                </div>
            )}
            {loading && !showAvatar && (
                <div className={styles.loadingContainer}>
                    <Spin spinning={loading} size="small"/>
                </div>
            )}
        </div>
    );
});

export default AvatarUpload;
