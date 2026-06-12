import React, {forwardRef, useImperativeHandle} from "react";
import {App} from "antd";
import {deliverOrder} from "@/service/admin/pay";

export type DeliverRef = {
    open: (record: any) => void;
};

type DeliverProps = {
    hasControl?: boolean;
    canDeliver?: (record: any) => boolean;
    onSuccess?: () => void;
};

const Deliver = forwardRef<DeliverRef, DeliverProps>(({hasControl, canDeliver, onSuccess}, ref) => {
    const {message, modal} = App.useApp();

    useImperativeHandle(ref, () => ({
        open: (record: any) => {
            if (!hasControl) {
                message?.warning("暂无订单发货权限");
                return;
            }
            if (canDeliver && !canDeliver(record)) {
                message?.warning("当前订单不可发货");
                return;
            }
            const id = Number(record?.id || record);
            if (!id) {
                message?.warning("订单ID错误");
                return;
            }
            modal.confirm({
                title: "订单发货",
                content: "确定更改状态并发货吗?",
                okText: "确定",
                cancelText: "取消",
                mask: {
                    closable: true,
                },
                onOk: async () => {
                    await deliverOrder({
                        body: {id},
                        onSuccess: (res: any) => {
                            onSuccess?.();
                            message?.success(res?.message || "发货成功");
                        },
                        onFail: (res: any) => {
                            message?.error(res?.message || "发货失败");
                        },
                    });
                },
            } as any);
        },
    }));

    return null;
});

Deliver.displayName = "Deliver";

export default Deliver;
