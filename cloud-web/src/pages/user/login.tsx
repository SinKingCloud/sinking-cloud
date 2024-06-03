import React, {useEffect, useRef, useState} from "react";
import {createStyles, useResponsive} from "antd-style";
import {App, Button, Checkbox, Col, Form, Input, Row, Tabs} from "antd";
import defaultSettings from "../../../config/defaultSettings";
import {Variant} from "antd/es/form/hooks/useVariants";
import {historyPush} from "@/utils/route";
import {LeftOutlined} from "@ant-design/icons";
import {getSms, webLogin} from "@/service/login";
import {setLoginToken} from "@/utils/auth";
import {useModel} from "umi";
import {Captcha, Body} from "@/components";
import {CaptchaRef} from "@/components/captcha";

const useStyles = createStyles(({token, css, responsive}): any => {
    return {
        body: {
            padding: "0px",
            backgroundColor: "white",
        },
        main: css`
            height: 100vh;
            width: 100%;
            user-select: none;

            ${responsive.sm} {
                height: 100%;
            }
        `,
        left: {
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
        },
        right: {
            backgroundImage: "url(" + (defaultSettings?.basePath || "/") + "login.png)",
            backgroundRepeat: "no-repeat",
            backgroundSize: "cover"
        },
        form: css`
            width: 450px;
            height: 600px;
            margin: 50px 50px;

            ${responsive.sm} {
                width: 100%;
                height: 100%;
                margin: 0;
                padding: 50px 10px;
            }
        `,
        logoContainer: css`
            ${responsive.sm} {
                text-align: center;
            }
        `,
        logo: css`
            width: 150px;
            cursor: pointer;
        `,
        login: css`
            padding: 70px 50px;
            width: 450px;
            margin: auto;

            ${responsive.sm} {
                padding: 70px 10px;
                width: 350px;
                height: 100%;
            }
        `,
        tab: css`
            .ant-tabs-tab-btn {
                font-size: 15px;
                font-weight: bold;
                color: #979797;
            }

            .ant-tabs-ink-bar {
                height: 5px !important;
                border-radius: 5px;
                width: 35px !important;
            }

            .ant-tabs-ink-bar-animated {
                transition: width 0.1s, left 0.1s, right 0.1s !important;
            }

            .ant-tabs-tabpane {
                margin-top: 15px;
            }

            .ant-tabs-tab.ant-tabs-tab-active .ant-tabs-tab-btn {
                font-size: 17px;
            }

            .ant-tabs-nav {
                height: 50px;
            }

            .ant-tabs-nav::before {
                border-bottom: none;
            }
        `,
        input: css`
            height: 50px;
            font-size: 14px;
            background: rgb(255, 255, 255);
            border: medium;
            border-radius: 8px;
            box-shadow: rgba(0, 0, 0, 0.05) 0 7px 25px 0;
            padding: 0 20px;

            input::placeholder {
                color: darkgray;
            }
        `,
        item: css`
            height: 55px;
        `,
        btn: css`
            height: 45px;
            border-radius: 6px;
            font-size: 15px;
            box-shadow: rgba(79, 122, 253, 0.5) 0 6px 27px 0;
        `,
        gray: css`
            margin: 0 5px;
            color: rgba(0, 0, 0, 0.35);
        `,
        leftBtn: {
            float: "left",
            color: "rgba(0, 0, 0, 0.35)"
        },
        rightBtn: {
            float: "right",
            color: "rgba(0, 0, 0, 0.35)"
        },
        sendSmsBtn: css`
            font-size: 12px;
            color: ${token?.colorPrimaryText};
            cursor: pointer;
        `,
        leftBtn2: {
            float: "left",
            color: "rgba(0, 0, 0, 0.5)"
        },
        rightBtn2: {
            float: "right",
            color: "rgba(0, 0, 0, 0.35)",
            marginTop: "21px",
            cursor: "pointer"
        },
        findPwdTitle: {
            overflow: "hidden",
            marginTop: "-10px"
        },
        sms: {
            fontSize: "12px",
            color: "gray",
            cursor: "pointer",
        }
    };
});

export default () => {
    /**
     * 样式表
     */const {
        styles: {
            body, main, right, left, form, logoContainer, logo, login, tab, input, btn, gray, item,
            leftBtn, rightBtn, leftBtn2, rightBtn2, findPwdTitle, sendSmsBtn, sms
        }
    } = useStyles();
    /**
     * 相关组件
     */
    const {message, modal} = App.useApp();
    const {mobile} = useResponsive();
    const user = useModel("user");
    const theme = useModel("theme");
    /**
     * 找回密码panel
     */
    const [loginPanel, setLoginPanel] = useState(true);
    const showLoginPanel = () => {
        setLoginPanel(true);
    }
    const showFindPwdPanel = () => {
        setLoginPanel(false);
    }
    /**
     * 显示验证码
     * @param callback 验证回调
     */
    const captcha = useRef<CaptchaRef>({});
    /**
     * 密码登录
     * @param values 登录表单
     */
    const [pwdLoginBtnLoading, setPwdLoginBtnLoading] = useState(false);

    /**
     * 验证码登录
     */
    const [smsForm] = Form.useForm();
    /**
     * 登录处理
     */
    const checkJump = () => {
        user?.refreshWebUser(() => {
            historyPush("user.index");
        });
    }
    const loginOnFinish = (values: any) => {
        setPwdLoginBtnLoading(true);
        captcha?.current?.Show?.(async (res) => {
            await webLogin({
                body: {
                    token: res?.randstr,
                    captcha: res?.ticket,
                    device: mobile ? "mobile" : "pc",
                    ...values
                },
                onSuccess: (r) => {
                    setLoginToken(mobile ? "mobile" : "pc", r?.data);
                    message?.success(r?.message);
                    checkJump();
                },
                onFail: (r) => {
                    modal.error({
                        title: '登录失败',
                        content: r?.message,
                        okText: "确认"
                    })
                }
            });
            setPwdLoginBtnLoading(false);
        }, () => {
            message?.warning("请完成验证码认证");
            setPwdLoginBtnLoading(false);
        });
    }

    /**
     * 验证码处理
     */
    const [countdown, setCountdown] = useState(0); // 倒计时秒数
    const startCountdown = () => {
        let timeLeft = 60;
        setCountdown(timeLeft);
        const timer = setInterval(() => {
            if (timeLeft <= 0) {
                clearInterval(timer);
                setCountdown(0);
            } else {
                setCountdown(timeLeft - 1);
            }
            timeLeft -= 1;
        }, 1000);
    };
    const getCode = (phone: any) => {
        if (phone == undefined || phone == "") {
            return message?.info("请输入手机号");
        }
        if (countdown > 0) return
        captcha?.current?.Show?.(async (res) => {
            setCountdown(-1);
            await getSms({
                body: {
                    token: res?.randstr,
                    captcha: res?.ticket,
                    phone: phone
                },
                onSuccess: (r) => {
                    message?.success(r?.message)
                    startCountdown()
                },
                onFail: (r) => {
                    message?.error(r?.message)
                }
            })
        }, () => {
            message?.warning("请完成验证码认证")
        })
    }
    const SendSms: React.FC<{ form?: any }> = (props) => {
        return <>
            {countdown <= 0 && <a className={sendSmsBtn}
                                  onClick={() => getCode(props?.form?.getFieldValue("phone"))}>
                {countdown < 0 ? '发送中...' : '发送验证码'}
            </a>}
            {countdown > 0 && <span
                className={countdown > 0 ? sms : sendSmsBtn}>{countdown > 0 && `${countdown}s后重新发送`}</span>}
        </>;
    }

    useEffect(() => {
        checkJump();
    }, []);

    return (
        <Body className={body} themes={theme?.getLightTheme()} loading={false} breadCrumb={false} space={false}>
            <Captcha ref={captcha}/>
            <Row className={main}>
                <Col xs={24} sm={24} md={11} className={left}>
                    <div className={form}>
                        <div className={logoContainer}>
                            <img src={(defaultSettings?.basePath || "/") + "logo_light.png"} className={logo}
                                 alt={defaultSettings?.title} onClick={() => {
                                historyPush("/")
                            }}/>
                        </div>
                        <div className={login}>
                            {(loginPanel && <Tabs rootClassName={tab} defaultActiveKey="1" items={[
                                {
                                    key: '1',
                                    label: '密码登录',
                                    forceRender: true,
                                    children: <Form name="pwd_login" onFinish={loginOnFinish}>
                                        <Form.Item name="account" className={item}
                                                   rules={[{required: true, message: '请输入账号或手机号'}]}>
                                            <Input prefix={<div>
                                                <span className={gray}>+86</span>
                                                <span className={gray}>|</span>
                                            </div>} allowClear placeholder="请输入账号或手机号" className={input}
                                                   variant={0 as Variant}/>
                                        </Form.Item>
                                        <Form.Item name="password" className={item}
                                                   rules={[{required: true, message: '请输入密码'}]}>
                                            <Input.Password placeholder="请输入密码" className={input}
                                                            variant={0 as Variant}/>
                                        </Form.Item>
                                        <Form.Item>
                                            <Checkbox className={leftBtn} defaultChecked>下次自动登录</Checkbox>
                                            <a className={rightBtn} onClick={showFindPwdPanel}>忘记密码</a>
                                        </Form.Item>
                                        <Form.Item>
                                            <Button type="primary" htmlType="submit" block className={btn}
                                                    loading={pwdLoginBtnLoading}>
                                                登录
                                            </Button>
                                        </Form.Item>
                                    </Form>,
                                },
                                {
                                    key: '2',
                                    label: '验证码登录',
                                    forceRender: true,
                                    children:
                                        <Form name="phone_login" form={smsForm} onFinish={loginOnFinish}>
                                            <Form.Item name="phone" className={item}
                                                       rules={[{required: true, message: '请输入手机号'}]}>
                                                <Input allowClear placeholder="请输入手机号" className={input}
                                                       variant={0 as Variant}
                                                       prefix={<div>
                                                           <span className={gray}>+86</span>
                                                           <span className={gray}>|</span>
                                                       </div>}/>
                                            </Form.Item>
                                            <Form.Item name="code" className={item}
                                                       rules={[{required: true, message: '请输入验证码'}]}>
                                                <Input allowClear placeholder="请输入验证码" className={input}
                                                       suffix={<SendSms
                                                           form={smsForm}/>}
                                                       variant={0 as Variant}/>
                                            </Form.Item>
                                            <Form.Item>
                                                <Checkbox className={leftBtn} defaultChecked>下次自动登录</Checkbox>
                                            </Form.Item>
                                            <Form.Item>
                                                <Button type="primary" htmlType="submit" block className={btn}
                                                        loading={pwdLoginBtnLoading}>
                                                    登录/注册
                                                </Button>
                                            </Form.Item>
                                        </Form>,
                                },
                            ]}/>) || <>
                                <div className={findPwdTitle}>
                                    <h2 className={leftBtn2}>找回密码</h2>
                                    <h4 className={rightBtn2} onClick={showLoginPanel}><LeftOutlined/>返回登录</h4>
                                </div>
                                <Form name="phone_login">
                                    <Form.Item name="phone" className={item}
                                               rules={[{required: true, message: '请输入手机号'}]}>
                                        <Input allowClear placeholder="请输入手机号" className={input}
                                               variant={0 as Variant}
                                               prefix={<div>
                                                   <span className={gray}>+86</span>
                                                   <span className={gray}>|</span>
                                               </div>}/>
                                    </Form.Item>
                                    <Form.Item name="code" className={item}
                                               rules={[{required: true, message: '请输入验证码'}]}>
                                        <Input allowClear placeholder="请输入验证码" className={input}
                                               suffix={<SendSms form={smsForm}/>}
                                               variant={0 as Variant}/>
                                    </Form.Item>
                                    <Form.Item name="password" className={item}
                                               rules={[{required: true, message: '请输入新密码'}]}>
                                        <Input.Password placeholder="新密码(8-20位字母或数字符号)" className={input}
                                                        ariant={0 as Variant}/>
                                    </Form.Item>
                                    <Form.Item>
                                        <Checkbox className={leftBtn} defaultChecked>下次自动登录</Checkbox>
                                    </Form.Item>
                                    <Form.Item>
                                        <Button type="primary" htmlType="submit" block className={btn}>
                                            登录
                                        </Button>
                                    </Form.Item>
                                </Form>,
                            </>
                            }
                        </div>
                    </div>
                </Col>
                <Col xs={0} sm={0} md={13} className={right}/>
            </Row>
        </Body>
    );
};
