import React from "react";
//@ts-ignore
import Script from 'react-load-script';
export type CaptchaProps = {

};
export interface CaptchaRef {
    show?: (callback: (res: any) => void) => void;
}
const Captcha = React.forwardRef<CaptchaRef, CaptchaProps>((props, ref) => {

    /**
     * 显示验证码
     * @param callback 验证回调
     * @param onSuccess
     * @param onFail
     */
        // @ts-ignore
    const show = (callback: (res: any) => void = undefined,onSuccess: (res: any) => void = undefined,onFail: (res: any) => void = undefined) => {
            // @ts-ignore
            const captcha = new TencentCaptcha('2046626881', (res) => {
                callback?.(res);
                if(res.set == 0){
                    onSuccess?.(res);
                    onFail?.(res);
                }
            });
            captcha?.show();
        }
    React.useImperativeHandle(ref, () => ({
        show: show,
    }));
    return (
        <Script url={"https://ssl.captcha.qq.com/TCaptcha.js"}/>
    );
});
export default Captcha;