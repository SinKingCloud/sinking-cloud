import React, {useEffect, useRef, useState} from "react";
import {Button, Spin, Tooltip} from "antd";
import {createStyles} from "antd-style";
import {useParams} from "umi";
import Mind, {type MindRef, type MindValue} from "@/pages/components/mind";
import {getMindShare} from "@/service/user/mind";
import {getResponseData, getTimeText, parseMindValue} from "@/pages/user/mind/utils";
import {historyPush} from "@/utils/route";

type SharedMindData = {
    title: string;
    describe: string;
    value: MindValue;
    updateTime: string;
};

const useStyles = createStyles(({css, token}: any) => ({
    page: css`
        position: fixed;
        inset: 0;
        display: flex;
        flex-direction: column;
        min-width: 0;
        min-height: 100vh;
        min-height: 100dvh;
        overflow: hidden;
        background: ${token.colorBgLayout};
    `,
    header: css`
        display: flex;
        flex: none;
        align-items: center;
        justify-content: space-between;
        gap: 20px;
        min-width: 0;
        min-height: 42px;
        padding: 4px 14px;
        border-bottom: 1px solid ${token.colorBorderSecondary};
        background: ${token.colorBgContainer};

        .share-document-info {
            display: flex;
            flex: 1;
            align-items: center;
            gap: 10px;
            min-width: 0;
        }

        h1 {
            flex: 0 1 auto;
            max-width: min(42vw, 420px);
            margin: 0;
            overflow: hidden;
            color: ${token.colorTextHeading};
            font-size: 15px;
            font-weight: 600;
            line-height: 18px;
            letter-spacing: 0;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        p {
            flex: 1;
            min-width: 0;
            margin: 0;
            overflow: hidden;
            color: ${token.colorTextTertiary};
            font-size: 11px;
            line-height: 15px;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .share-update-time {
            flex: none;
            color: ${token.colorTextTertiary};
            font-size: 11px;
            line-height: 18px;
            white-space: nowrap;
        }

        @media (max-width: 640px) {
            min-height: 40px;
            padding: 4px 10px;

            p {
                display: none;
            }

            .share-update-time {
                display: none;
            }
        }
    `,
    canvas: css`
        flex: 1;
        min-width: 0;
        min-height: 0;
    `,
    state: css`
        display: flex;
        flex: 1;
        align-items: center;
        justify-content: center;
        min-width: 0;
        min-height: 0;
        padding: 24px;
    `,
    unavailable: css`
        width: min(100%, 360px);
        text-align: center;

        h1 {
            margin: 0;
            color: ${token.colorTextHeading};
            font-size: 18px;
            font-weight: 600;
            line-height: 26px;
            letter-spacing: 0;
        }

        p {
            margin: 6px 0 20px;
            color: ${token.colorTextSecondary};
            font-size: 13px;
            line-height: 21px;
            letter-spacing: 0;
        }
    `,
}));

const MindSharePage: React.FC = () => {
    const {styles} = useStyles();
    const params = useParams<"code">();
    const mindRef = useRef<MindRef | null>(null);
    const code = String(params.code || "").trim();
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");
    const [data, setData] = useState<SharedMindData>();

    useEffect(() => {
        let active = true;
        setData(undefined);
        setError("");
        if (code.length !== 8) {
            setLoading(false);
            setError("分享不可用");
            return () => {
                active = false;
            };
        }
        setLoading(true);
        void getMindShare({body: {code}}).then((response) => {
            if (!active) return;
            const result = getResponseData<any>(response, "获取分享内容失败");
            setData({
                title: result?.title || "未命名导图",
                describe: result?.describe || "",
                value: parseMindValue(result?.value),
                updateTime: result?.update_time || "",
            });
        }).catch(() => {
            if (active) {
                setError("分享不可用");
            }
        }).finally(() => {
            if (active) setLoading(false);
        });
        return () => {
            active = false;
        };
    }, [code]);

    useEffect(() => {
        if (!data) return;
        let innerFrame = 0;
        const frame = requestAnimationFrame(() => {
            innerFrame = requestAnimationFrame(() => mindRef.current?.fitView());
        });
        return () => {
            cancelAnimationFrame(frame);
            cancelAnimationFrame(innerFrame);
        };
    }, [data]);

    return (
        <div className={styles.page}>
            {loading ? (
                <div className={styles.state}><Spin size="large"/></div>
            ) : error || !data ? (
                <div className={styles.state}>
                    <div className={styles.unavailable}>
                        <h1>分享不可用</h1>
                        <p>链接可能已过期、被关闭或不存在</p>
                        <Button type="primary" onClick={() => historyPush("/")}>返回首页</Button>
                    </div>
                </div>
            ) : (
                <>
                    <header className={styles.header}>
                        <div className="share-document-info">
                            <h1>{data.title}</h1>
                            {data.describe && <p>{data.describe}</p>}
                        </div>
                        {data.updateTime && (
                            <Tooltip title={data.updateTime}>
                                <span className="share-update-time">更新于 {getTimeText(data.updateTime)}</span>
                            </Tooltip>
                        )}
                    </header>
                    <main className={styles.canvas}>
                        <Mind
                            key={code}
                            ref={mindRef}
                            defaultValue={data.value}
                            height="100%"
                            readOnly
                            showControls
                            showMiniMap
                        />
                    </main>
                </>
            )}
        </div>
    );
};

export default MindSharePage;
