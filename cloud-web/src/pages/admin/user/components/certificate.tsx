import React, {forwardRef, useEffect, useImperativeHandle, useRef, useState} from "react";
import {App, Empty, Image, Spin, Typography} from "antd";
import {createStyles} from "antd-style";
import {Icon, ProModal, ProModalRef, Title, useTheme} from "sinking-antd";

import {copyText} from "@/utils/string";
import {getUserCertificate} from "@/service/admin/user";
import {AdminControlCode} from "@/utils/permission";

export type CertificateRef = {
    open: (record: any) => void;
};

type CertificateProps = {
    control: any;
};

const useStyles: any = createStyles(({css, token}: any, props: any = {}) => {
    const isCompactMode = props?.isCompactMode || false;
    const contentGap = isCompactMode ? 5 : 8;
    const sectionGap = isCompactMode ? 4 : 6;
    const infoPadding = isCompactMode ? "4px 8px" : "8px 10px";
    const cardWidth = isCompactMode ? 144 : 180;
    const cardHeight = isCompactMode ? 88 : 116;
    const mobileCardHeight = isCompactMode ? 78 : 96;
    const mediaPadding = isCompactMode ? "5px 8px" : "10px 12px";
    const playSize = isCompactMode ? 30 : 44;
    const playFontSize = isCompactMode ? 19 : 28;

    return {
        loading: css`
            min-height: ${isCompactMode ? 160 : 220}px;
            display: flex;
            align-items: center;
            justify-content: center;
        `,
        content: css`
            display: flex;
            flex-direction: column;
            gap: ${contentGap}px;
        `,
        section: css`
            display: flex;
            flex-direction: column;
            gap: ${sectionGap}px;
        `,
        sectionHead: css`
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: ${token.marginSM}px;
            padding-top: ${isCompactMode ? 0 : 2}px;
        `,
        sectionTitle: css`
            color: ${token.colorText};
            font-size: ${token.fontSizeSM}px;
            font-weight: 600;
            line-height: ${isCompactMode ? 18 : 20}px;
        `,
        sectionExtra: css`
            color: ${token.colorTextTertiary};
            font-size: ${token.fontSizeSM}px;
        `,
        infoPanel: css`
            display: grid;
            grid-template-columns: repeat(2, minmax(0, 1fr));
            overflow: hidden;
            border: 1px solid ${token.colorBorderSecondary};
            border-radius: ${token.borderRadius}px;
            background: ${token.colorBgContainer};

            @media (max-width: 575px) {
                grid-template-columns: 1fr;
            }
        `,
        infoItem: css`
            min-width: 0;
            padding: ${infoPadding};
            border-right: 1px solid ${token.colorBorderSecondary};
            border-bottom: 1px solid ${token.colorBorderSecondary};

            &:nth-of-type(2n) {
                border-right: none;
            }

            &:nth-last-of-type(-n + 2) {
                border-bottom: none;
            }

            @media (max-width: 575px) {
                border-right: none;

                &:nth-last-of-type(-n + 2) {
                    border-bottom: 1px solid ${token.colorBorderSecondary};
                }

                &:last-child {
                    border-bottom: none;
                }
            }
        `,
        label: css`
            margin-bottom: ${isCompactMode ? 2 : 6}px;
            color: ${token.colorTextTertiary};
            font-size: ${token.fontSizeSM}px;
            line-height: ${isCompactMode ? 16 : 18}px;
        `,
        value: css`
            max-width: 100%;
            color: ${token.colorText};
            font-size: ${token.fontSize}px;
            line-height: ${isCompactMode ? 18 : 22}px;
            word-break: break-all;
        `,
        mediaGrid: css`
            display: grid;
            grid-template-columns: repeat(auto-fill, ${cardWidth}px);
            align-items: start;
            gap: ${isCompactMode ? 8 : token.margin}px;

            @media (max-width: 575px) {
                grid-template-columns: repeat(2, minmax(0, 1fr));
                gap: ${isCompactMode ? 6 : token.marginSM}px;
            }
        `,
        mediaCard: css`
            width: ${cardWidth}px;
            overflow: hidden;
            border: 1px solid ${token.colorBorderSecondary};
            border-radius: ${token.borderRadiusLG}px;
            background: ${token.colorBgContainer};
            transition: border-color 0.2s, box-shadow 0.2s;

            &:hover {
                border-color: ${token.colorPrimaryBorder};
                box-shadow: ${token.boxShadowTertiary};
            }

            @media (max-width: 575px) {
                width: 100%;
                min-width: 0;
            }
        `,
        mediaThumb: css`
            width: 100%;
            height: ${cardHeight}px;
            position: relative;
            display: flex;
            align-items: center;
            justify-content: center;
            overflow: hidden;
            background: ${token.colorFillQuaternary};

            & .ant-image {
                width: 100%;
                height: 100%;
                display: block;
            }

            & .ant-image-img {
                width: 100%;
                height: 100%;
                display: block;
                object-fit: cover;
            }

            @media (max-width: 575px) {
                height: ${mobileCardHeight}px;
            }
        `,
        image: css`
            width: 100% !important;
            height: 100% !important;
            object-fit: cover;
            cursor: pointer;
        `,
        videoThumb: css`
            color: #fff;
            cursor: pointer;
            background: #111827;
        `,
        videoFrame: css`
            width: 100%;
            height: 100%;
            display: block;
            object-fit: contain;
            background: #111827;
        `,
        playMask: css`
            position: absolute;
            inset: 0;
            display: flex;
            align-items: center;
            justify-content: center;
            background: rgba(0, 0, 0, 0.18);
        `,
        playButton: css`
            width: ${playSize}px;
            height: ${playSize}px;
            display: flex;
            align-items: center;
            justify-content: center;
            border-radius: 50%;
            color: #fff;
            font-size: ${playFontSize}px;
            background: rgba(0, 0, 0, 0.42);
            box-shadow: 0 8px 24px rgba(0, 0, 0, 0.24);
        `,
        mediaMeta: css`
            padding: ${mediaPadding};
            border-top: 1px solid ${token.colorBorderSecondary};
        `,
        mediaTitleRow: css`
            min-width: 0;
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: ${token.marginXS}px;
        `,
        mediaName: css`
            min-width: 0;
            color: ${token.colorText};
            font-size: ${token.fontSizeSM}px;
            font-weight: 500;
            line-height: ${isCompactMode ? 16 : 20}px;
            overflow: hidden;
            white-space: nowrap;
            text-overflow: ellipsis;
        `,
        mediaActions: css`
            flex: none;
            display: flex;
            align-items: center;
            gap: ${isCompactMode ? 4 : token.marginXS}px;

            a {
                font-size: ${isCompactMode ? 10 : 11}px;
                line-height: 16px;
            }
        `,
        mediaEmpty: css`
            padding: ${isCompactMode ? token.paddingSM : token.paddingLG}px 0;
            border: 1px solid ${token.colorBorderSecondary};
            border-radius: ${token.borderRadius}px;
        `,
    };
});

const Certificate = forwardRef<CertificateRef, CertificateProps>(({control}, ref) => {
    const {message} = App.useApp();
    const theme = useTheme();
    const isCompactMode = theme?.isCompactTheme?.() || false;
    const {styles} = useStyles({isCompactMode});
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const videoRef = useRef<HTMLVideoElement | null>(null);
    const [info, setInfo] = useState<any>();
    const [loading, setLoading] = useState(false);
    const hasCertificateControl = control.hasControl(AdminControlCode.WEB_USER_CERTIFICATE);

    useEffect(() => {
        const onFullscreenChange = () => {
            if (!document.fullscreenElement && videoRef.current) {
                videoRef.current.pause();
                videoRef.current.controls = false;
            }
        };
        document.addEventListener("fullscreenchange", onFullscreenChange);
        return () => document.removeEventListener("fullscreenchange", onFullscreenChange);
    }, []);

    useImperativeHandle(ref, () => ({
        open: (record: any) => {
            if (!hasCertificateControl) {
                message?.warning("暂无实名查看权限");
                return;
            }
            const id = Number(record?.id);
            if (!id) {
                return;
            }
            setInfo(undefined);
            setLoading(true);
            modalRef.current?.show();
            void getUserCertificate({
                body: {id},
                onSuccess: (res: any) => {
                    setInfo(res?.data || {});
                },
                onFail: (res: any) => {
                    message?.error(res?.message || "获取实名信息失败");
                },
            }).finally(() => {
                setLoading(false);
            });
        },
    }));

    const pics = Array.isArray(info?.pics) ? info.pics.filter(Boolean) : [];
    const hasCompanyInfo = !!(info?.company_name || info?.company_no);
    const hasMedia = pics.length > 0 || !!info?.video;
    const infoItems = [
        {label: "认证类型", value: hasCompanyInfo ? "企业实名" : "个人实名"},
        {label: "认证时间", value: info?.auth_time},
        {label: "姓名", value: info?.card_name, copyable: true},
        {label: "证件号", value: info?.card_no, copyable: true},
        hasCompanyInfo ? {label: "企业名称", value: info?.company_name, copyable: true} : null,
        hasCompanyInfo ? {label: "企业证件", value: info?.company_no, copyable: true} : null,
    ].filter(Boolean) as any[];

    const renderValue = (item: any) => (
        <Typography.Text className={styles.value} copyable={item?.copyable && item?.value ? {text: item.value} : false}>
            {item?.value || "-"}
        </Typography.Text>
    );

    const playFullscreenVideo = async () => {
        const video = videoRef.current;
        if (!video) {
            return;
        }
        video.controls = true;
        video.muted = false;
        const webkitVideo = video as any;
        let fullscreenFailed = false;
        try {
            if (video.requestFullscreen) {
                await video.requestFullscreen();
            } else if (webkitVideo.webkitEnterFullscreen) {
                webkitVideo.webkitEnterFullscreen();
            } else {
                fullscreenFailed = true;
            }
        } catch {
            fullscreenFailed = true;
        }
        try {
            await video.play();
        } catch {
            // 忽略自动播放限制，保留原生控制条给用户手动播放。
        }
        if (fullscreenFailed) {
            message?.info("当前浏览器未进入全屏，可使用播放器全屏按钮");
        }
    };

    const openLink = (url?: string) => {
        if (!url) {
            return;
        }
        window.open(url, "_blank", "noopener,noreferrer");
    };

    const copyLink = async (url?: string) => {
        if (!url) {
            return;
        }
        if (await copyText(url)) {
            message?.success("复制成功");
            return;
        }
        message?.error("复制失败");
    };

    return (
        <ProModal
            ref={modalRef}
            title={<Title>实名信息</Title>}
            width={540}
            modalProps={{
                footer: null,
                style: {
                    top: 100,
                    paddingBottom: 100,
                },
                mask: {
                    closable: true,
                },
                destroyOnHidden: false,
                styles: {
                    body: {
                        padding: isCompactMode ? 6 : undefined,
                    },
                },
                afterClose: () => {
                    if (videoRef.current) {
                        videoRef.current.pause();
                        videoRef.current.controls = false;
                    }
                    setInfo(undefined);
                    setLoading(false);
                },
            } as any}>
            {loading ? (
                <div className={styles.loading}>
                    <Spin/>
                </div>
            ) : info?.status ? (
                <div className={styles.content}>
                    <div className={styles.section}>
                        <div className={styles.sectionHead}>
                            <div className={styles.sectionTitle}>基础资料</div>
                        </div>
                        <div className={styles.infoPanel}>
                            {infoItems.map((item) => (
                                <div className={styles.infoItem} key={item.label}>
                                    <div className={styles.label}>{item.label}</div>
                                    {renderValue(item)}
                                </div>
                            ))}
                        </div>
                    </div>

                    <div className={styles.section}>
                        <div className={styles.sectionHead}>
                            <div className={styles.sectionTitle}>留存资料</div>
                            {hasMedia &&
                                <div className={styles.sectionExtra}>{pics.length + (info?.video ? 1 : 0)} 项</div>}
                        </div>
                        {hasMedia ? (
                            <div className={styles.mediaGrid}>
                                <Image.PreviewGroup>
                                    {pics.map((url: string, index: number) => (
                                        <div className={styles.mediaCard} key={`${url}-${index}`}>
                                            <div className={styles.mediaThumb}>
                                                <Image src={url} className={styles.image}/>
                                            </div>
                                            <div className={styles.mediaMeta}>
                                                <div className={styles.mediaTitleRow}>
                                                    <div className={styles.mediaName}>证件影像 {index + 1}</div>
                                                    <div className={styles.mediaActions}>
                                                        <Typography.Link
                                                            onClick={() => openLink(url)}>下载</Typography.Link>
                                                        <Typography.Link
                                                            onClick={() => copyLink(url)}>复制</Typography.Link>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    ))}
                                </Image.PreviewGroup>
                                {info?.video && (
                                    <div className={styles.mediaCard}>
                                        <div className={`${styles.mediaThumb} ${styles.videoThumb}`}
                                             onClick={playFullscreenVideo}>
                                            <video
                                                ref={videoRef}
                                                src={info.video}
                                                className={styles.videoFrame}
                                                muted
                                                playsInline
                                                preload="metadata"
                                            />
                                            <div className={styles.playMask}>
                                                <div className={styles.playButton}>
                                                    <Icon type="PlayCircleOutlined"/>
                                                </div>
                                            </div>
                                        </div>
                                        <div className={styles.mediaMeta}>
                                            <div className={styles.mediaTitleRow}>
                                                <div className={styles.mediaName}>认证视频</div>
                                                <div className={styles.mediaActions}>
                                                    <Typography.Link
                                                        onClick={() => openLink(info.video)}>下载</Typography.Link>
                                                    <Typography.Link
                                                        onClick={() => copyLink(info.video)}>复制</Typography.Link>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                )}
                            </div>
                        ) : (
                            <div className={styles.mediaEmpty}>
                                <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="暂无留存资料"/>
                            </div>
                        )}
                    </div>
                </div>
            ) : (
                <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="暂无实名信息"/>
            )}
        </ProModal>
    );
});

Certificate.displayName = "Certificate";

export default Certificate;
