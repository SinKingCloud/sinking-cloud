import React, {useEffect, useRef, useState} from "react";
import {Drawer} from "antd";
import {useModel} from "umi";
import {Body, Icon, useTheme} from "sinking-antd";
import Settings from "@/../config/defaultSettings";
import {deleteHeader} from "@/utils/auth";
import {logout} from "@/service/common/login";
import {historyPush} from "@/utils/route";
import HomeVisual from "./components/visual";
import {createStyles} from "antd-style";

const useStyles = createStyles(({css, token, isDarkMode}) => {
    const ink = token.colorTextHeading;
    const muted = token.colorTextSecondary;
    const line = token.colorBorderSecondary;
    const paper = token.colorBgContainer;
    const accent = token.colorPrimary;
    const wash = `color-mix(in srgb, ${accent} 2%, ${paper})`;

    return {
        body: css`padding: 0 !important;`,
        page: css`
            position: relative; overflow: clip; width: 100%; color: ${ink}; background: ${paper}; -webkit-font-smoothing: antialiased;
            *, *::before, *::after { box-sizing: border-box; }
            h1, h2, h3, p { margin: 0; }
            button, a { -webkit-tap-highlight-color: transparent; }
            button { font-family: inherit; cursor: pointer; }
            button:focus-visible, a:focus-visible { outline: 2px solid ${accent}; outline-offset: 5px; }
            .home-reveal { transition: opacity .45s ease, transform .45s ease; }
            &[data-motion="true"] .home-reveal:not([data-visible="true"]) { opacity: 0; transform: translateY(14px); }
            @media (prefers-reduced-motion: reduce) { *, *::before, *::after { animation: none !important; transition: none !important; scroll-behavior: auto !important; } }
        `,
        container: css`
            width: min(100%, 1450px); margin-inline: auto; padding-inline: clamp(24px, 5vw, 80px);
            @media (max-width: 575px) { padding-inline: 22px; }
        `,
        header: css`
            position: fixed; inset: 0 0 auto; z-index: 50; height: 76px; border-bottom: 1px solid transparent;
            background: transparent; transition: background .25s, border-color .25s;
            &[data-scrolled="true"] {
                background: color-mix(in srgb, ${paper} 97%, transparent);
                border-color: color-mix(in srgb, ${line} 60%, transparent);
                @supports ((backdrop-filter: blur(1px)) or (-webkit-backdrop-filter: blur(1px))) {
                    background: color-mix(in srgb, ${paper} ${isDarkMode ? 68 : 64}%, transparent);
                    -webkit-backdrop-filter: blur(20px) saturate(150%);
                    backdrop-filter: blur(20px) saturate(150%);
                }
            }
            @media (max-width: 767px) { height: 64px; }
        `,
        headerInner: css`display: flex; align-items: center; justify-content: space-between; gap: 34px; height: 100%;`,
        logo: css`
            display: inline-flex; align-items: center; flex-shrink: 0; gap: 10px; border: 0; padding: 0; color: ${ink}; background: transparent; font-size: 18px; font-weight: 650; letter-spacing: -.5px;
            img { width: 32px; height: 32px; object-fit: contain; }
            .anticon { color: ${accent}; font-size: 31px; }
            @media (max-width: 575px) { font-size: 16px; img { width: 28px; height: 28px; } }
        `,
        nav: css`
            display: flex; align-items: center; justify-content: center; gap: 30px; margin-left: auto;
            button { position: relative; padding: 9px 0; border: 0; color: ${muted}; background: transparent; font-size: 13px; transition: color .2s;
                &:hover, &[aria-current="true"] { color: ${ink}; }
                &[aria-current="true"]::after { position: absolute; bottom: 0; left: calc(50% - 2px); width: 4px; height: 4px; border-radius: 50%; background: ${accent}; content: ""; }
            }
            @media (max-width: 991px) { gap: 22px; }
            @media (max-width: 767px) { display: none; }
        `,
        headerActions: css`
            display: flex; align-items: center; gap: 20px;
            @media (max-width: 991px) { gap: 12px; }
            @media (max-width: 767px) { display: none; }
        `,
        textButton: css`
            display: inline-flex; align-items: center; justify-content: center; gap: 10px; border: 0; padding: 0; color: ${accent}; background: transparent; font-size: 13px; transition: opacity .2s;
            &:hover { opacity: .7; }
        `,
        themeButton: css`
            display: grid; place-items: center; width: 32px; height: 32px; border: 0; border-radius: 50%; color: ${muted}; background: transparent; font-size: 18px;
            &:hover { color: ${accent}; background: ${wash}; }
        `,
        primaryButton: css`
            position: relative; overflow: hidden; display: inline-flex; align-items: center; justify-content: center; gap: 22px; height: 48px; border: 1px solid ${accent}; border-radius: 7px; padding: 0 25px;
            color: ${token.colorTextLightSolid}; background: ${accent}; font-size: 14px; font-weight: 500; transition: background .18s, border-color .18s, box-shadow .2s;
            &::before { position: absolute; inset: -30% auto -30% -65%; width: 45%; background: linear-gradient(90deg, transparent, rgba(255,255,255,.23), transparent); transform: skewX(-22deg); transition: left .5s ease; content: ""; }
            &:hover { background: ${token.colorPrimaryHover}; border-color: ${token.colorPrimaryHover}; box-shadow: 0 6px 20px -9px color-mix(in srgb, ${accent} 55%, transparent); &::before { left: 125%; } }
            &:active { background: ${token.colorPrimaryActive}; border-color: ${token.colorPrimaryActive}; }
            .anticon { font-size: 13px; }
            &[data-small="true"] { height: 36px; gap: 12px; padding: 0 17px; font-size: 12px; border-radius: 6px; }
        `,
        secondaryButton: css`
            display: inline-flex; align-items: center; justify-content: center; gap: 14px; height: 48px; border: 1px solid ${line}; border-radius: 8px; padding: 0 24px;
            color: ${ink}; background: ${paper}; font-size: 14px; transition: border-color .18s, color .18s;
            &:hover { color: ${accent}; border-color: ${accent}; }
            .anticon { font-size: 12px; }
        `,
        mobileToggle: css`
            display: none;
            @media (max-width: 767px) { display: grid; place-items: center; width: 34px; height: 34px; border: 0; padding: 0; color: ${ink}; background: transparent; font-size: 20px; }
        `,
        mobileMenuBody: css`
            display: flex; flex-direction: column; padding: 20px 24px 28px !important; background: ${paper};
            nav { display: grid; margin-bottom: 30px; }
            nav button { display: flex; align-items: center; justify-content: space-between; border: 0; border-bottom: 1px solid ${line}; padding: 20px 0; color: ${ink}; background: transparent; font-size: 16px; .anticon { font-size: 12px; color: ${muted}; } }
        `,
        mobileTheme: css`
            margin-top: auto; margin-bottom: 24px;
            > span { display: block; margin-bottom: 12px; color: ${muted}; font-size: 12px; }
            > div { display: flex; gap: 8px; }
            button { display: inline-flex; align-items: center; justify-content: center; flex: 1; gap: 7px; height: 38px; border: 1px solid ${line}; border-radius: 6px; color: ${muted}; background: transparent; font-size: 13px;
                &[aria-pressed="true"] { color: ${accent}; border-color: ${accent}; background: ${wash}; }
            }
        `,
        mobileActions: css`display: grid; gap: 12px;`,
        spotlight: css`
            display: none; position: absolute; z-index: -1; top: 0; left: 0; width: 440px; height: 440px; pointer-events: none; opacity: 0;
            background: radial-gradient(circle, color-mix(in srgb, ${accent} ${isDarkMode ? 18 : 11}%, transparent), color-mix(in srgb, ${accent} 4%, transparent) 38%, transparent 70%);
            transform: translate3d(var(--spotlight-x, 0px), var(--spotlight-y, 0px), 0) translate(-50%, -50%);
            transition: opacity .4s ease;
            [data-spotlight="hero"] > & { width: 780px; height: 780px; }
            [data-highlighted="true"] > & { opacity: 1; }
            @media (hover: hover) and (pointer: fine) and (prefers-reduced-motion: no-preference) { display: block; }
        `,
        hero: css`
            position: relative; isolation: isolate; overflow: hidden; display: flex; flex-direction: column; min-height: 100vh; min-height: 100svh; padding-top: 76px;
            background: linear-gradient(115deg, ${paper} 22%, color-mix(in srgb, ${accent} ${isDarkMode ? 10 : 4}%, ${paper}) 75%, color-mix(in srgb, ${accent} ${isDarkMode ? 13 : 6}%, ${paper}));
            &::before {
                position: absolute; z-index: -1; inset: 0; pointer-events: none;
                background-image: linear-gradient(color-mix(in srgb, ${accent} 5%, transparent) 1px, transparent 1px), linear-gradient(90deg, color-mix(in srgb, ${accent} 5%, transparent) 1px, transparent 1px);
                background-size: 72px 72px;
                mask-image: radial-gradient(ellipse at 76% 46%, #000 12%, transparent 68%);
                -webkit-mask-image: radial-gradient(ellipse at 76% 46%, #000 12%, transparent 68%);
                content: "";
            }
            &::after { position: absolute; z-index: -1; inset: 10% -10% 0 40%; pointer-events: none; background: radial-gradient(ellipse at center, color-mix(in srgb, ${accent} ${isDarkMode ? 11 : 7}%, transparent), transparent 62%); content: ""; }
            @media (max-width: 767px) { padding-top: 64px; }
        `,
        heroInner: css`
            position: relative; display: grid; flex: 1; grid-template-columns: .92fr 1.08fr; align-items: center; gap: 22px;
            padding-top: 35px; padding-bottom: 40px;
            @media (max-width: 1100px) { grid-template-columns: 1fr 1fr; gap: 0; }
            @media (max-width: 767px) { grid-template-columns: 1fr; gap: 0; padding-top: 40px; padding-bottom: 5px; }
        `,
        heroCopy: css`
            position: relative; z-index: 2; min-width: 0; padding-bottom: 10px;
            h1 { margin-top: 27px; color: ${ink}; font-size: clamp(46px, 5.25vw, 78px); line-height: 1.23; font-weight: 700; letter-spacing: -3.4px; white-space: nowrap;
                span { display: block; color: ${accent}; background: linear-gradient(115deg, ${token.colorPrimaryActive} 8%, ${accent} 58%, ${token.colorPrimaryHover}); background-clip: text; -webkit-background-clip: text; -webkit-text-fill-color: transparent; }
            }
            > p { max-width: 390px; margin-top: 25px; color: ${muted}; font-size: 15px; line-height: 1.9; }
            @media (max-width: 1100px) { h1 { font-size: clamp(42px, 5.2vw, 60px); letter-spacing: -2.8px; } }
            @media (max-width: 767px) { padding-bottom: 0; text-align: center; h1 { margin-top: 22px; font-size: clamp(40px, 8.5vw, 58px); letter-spacing: -2px; } > p { max-width: 360px; margin: 19px auto 0; font-size: 13px; } }
        `,
        heroEyebrow: css`
            display: inline-flex; align-items: center; gap: 10px; color: ${accent}; font-size: 10px; font-weight: 500; letter-spacing: 1.6px;
            > span { width: 6px; height: 6px; border-radius: 50%; background: ${accent}; box-shadow: 0 0 0 5px color-mix(in srgb, ${accent} ${isDarkMode ? 10 : 6}%, transparent); }
            @media (max-width: 575px) { font-size: 9px; letter-spacing: 1.5px; }
        `,
        heroActions: css`
            display: flex; gap: 12px; margin-top: 34px;
            button { height: 50px; }
            @media (max-width: 767px) { justify-content: center; margin-top: 25px; }
            @media (max-width: 575px) { button { height: 44px; font-size: 13px; padding-inline: 20px; gap: 16px; } }
        `,
        heroNotes: css`
            display: flex; align-items: center; gap: 18px; margin-top: 29px; color: ${muted}; font-size: 11px;
            span { display: inline-flex; align-items: center; gap: 6px; }
            .anticon { color: ${accent}; font-size: 11px; }
            @media (max-width: 767px) { justify-content: center; gap: 20px; margin-top: 21px; font-size: 10px; }
        `,
        heroArt: css`
            position: relative; min-width: 0; width: calc(100% + 90px); margin-left: -30px; margin-right: -60px;
            @media (max-width: 1100px) { width: calc(100% + 60px); margin-inline: -20px -40px; }
            @media (max-width: 767px) { width: min(115%, 500px); max-width: none; justify-self: center; margin: -14px 0 -15px; }
        `,
        visual: css`
            width: 100%; max-height: min(650px, calc(100svh - 250px)); aspect-ratio: 1.05;
            @media (max-width: 767px) { max-height: 350px; }
        `,
        artCaption: css`
            position: absolute; bottom: 3%; left: 17%; right: 13%; display: flex; justify-content: space-between; align-items: center; gap: 15px; padding-top: 13px; border-top: 1px solid color-mix(in srgb, ${accent} 12%, transparent); color: ${muted}; font-size: 8px; letter-spacing: 1.5px;
            > span:last-child { color: ${accent}; font-size: 9px; letter-spacing: .7px; }
            @media (max-width: 767px) { display: none; }
        `,
        heroBottom: css`
            position: relative; display: flex; align-items: center; justify-content: space-between; gap: 25px; min-height: 86px; color: ${muted}; font-size: 10px;
            &::before { position: absolute; top: 0; left: clamp(24px, 5vw, 80px); right: clamp(24px, 5vw, 80px); height: 1px; background: ${line}; content: ""; }
            > span { font-size: 8px; letter-spacing: 1.8px; }
            > div { display: flex; gap: 42px; span { display: flex; align-items: center; gap: 13px; } b { color: ${accent}; font-size: 9px; font-weight: 400; font-family: monospace; } }
            button { display: grid; place-items: center; width: 32px; height: 32px; border: 1px solid ${line}; border-radius: 50%; color: ${muted}; background: transparent; transition: border-color .2s, color .2s; &:hover { color: ${accent}; border-color: ${accent}; } }
            @media (max-width: 767px) { min-height: 61px; > div { gap: 22px; span { gap: 7px; } } > span { display: none; } font-size: 10px; }
            @media (max-width: 360px) { > div { gap: 14px; } }
        `,
        section: css`
            padding-block: 112px;
            @media (max-width: 991px) { padding-block: 80px; }
            @media (max-width: 575px) { padding-block: 64px; }
        `,
        eyebrow: css`display: inline-block; color: ${muted}; font-size: 9px; line-height: 1.6; letter-spacing: 1.8px; font-weight: 500;`,
        sectionHeading: css`
            display: flex; align-items: flex-end; justify-content: space-between; gap: 35px; margin-bottom: 44px;
            h2 { margin-top: 18px; font-size: clamp(30px, 3vw, 42px); font-weight: 550; line-height: 1.4; letter-spacing: -1.3px; span { color: ${muted}; } }
            > p { max-width: 320px; padding-bottom: 4px; color: ${muted}; font-size: 13px; line-height: 1.9; }
            @media (max-width: 767px) { display: block; margin-bottom: 30px; h2 { font-size: 30px; letter-spacing: -.8px; } > p { margin-top: 20px; font-size: 12px; } }
        `,
        featureGrid: css`
            display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 18px;
            article { position: relative; overflow: hidden; isolation: isolate; padding: 32px 32px 35px; border: 1px solid ${line}; border-radius: 10px; background: ${paper}; transition: border-color .25s, transform .25s; }
            @media (hover: hover) and (pointer: fine) { article:hover { border-color: ${token.colorPrimaryBorder}; transform: translateY(-4px); } }
            h3 { margin-top: 24px; font-size: 16px; font-weight: 550; }
            p { max-width: 320px; margin-top: 11px; color: ${muted}; font-size: 12px; line-height: 1.9; }
            @media (max-width: 991px) { article { padding: 26px; } }
            @media (max-width: 767px) { grid-template-columns: repeat(2, minmax(0, 1fr)); article { padding: 24px; } }
            @media (max-width: 420px) { article { padding: 22px 18px; } h3 { margin-top: 19px; font-size: 15px; } p { font-size: 11px; } }
        `,
        featureTop: css`
            display: flex; align-items: center; justify-content: space-between;
            svg { width: 30px; height: 30px; color: ${accent}; } span { color: ${muted}; font-size: 9px; opacity: .5; font-variant-numeric: tabular-nums; }
        `,
        productsSection: css`background: ${wash};`,
        productWorkspace: css`
            display: grid; grid-template-columns: 215px minmax(0, 1fr); border: 1px solid ${line}; border-radius: 12px; overflow: hidden; background: ${paper};
            @media (max-width: 991px) { grid-template-columns: 175px minmax(0, 1fr); }
            @media (max-width: 767px) { grid-template-columns: 1fr; }
        `,
        productMenu: css`
            display: flex; flex-direction: column; justify-content: center; gap: 8px; padding: 22px 14px; border-right: 1px solid ${line}; scrollbar-width: thin;
            button { display: flex; align-items: center; gap: 12px; min-width: 0; height: 60px; padding: 0 17px; border: 0; border-radius: 6px; color: ${muted}; background: transparent; font-size: 13px; white-space: nowrap; transition: background .2s, color .2s;
                &[aria-pressed="true"] { color: ${accent}; background: ${wash}; }
                &:hover { color: ${accent}; } > svg { width: 25px; height: 25px; } > .anticon { margin-left: auto; font-size: 12px; }
            }
            @media (max-width: 991px) { padding-inline: 10px; button { padding-inline: 12px; gap: 9px; font-size: 12px; } }
            @media (max-width: 767px) { flex-direction: row; justify-content: flex-start; overflow-x: auto; padding: 10px; border-right: 0; border-bottom: 1px solid ${line}; button { flex: none; height: 42px; min-width: 120px; > svg { width: 22px; height: 22px; } > .anticon { display: none; } } }
        `,
        productDetail: css`
            display: grid; grid-template-columns: 1fr 1fr; gap: 36px; padding: 48px;
            @media (max-width: 1199px) { gap: 25px; padding: 34px; }
            @media (min-width: 768px) and (max-width: 991px) { grid-template-columns: 1fr; }
            @media (max-width: 767px) { grid-template-columns: 1fr; padding: 30px; gap: 30px; }
        `,
        productCopy: css`
            h3 { margin-top: 17px; font-size: 30px; font-weight: 550; }
            p { max-width: 340px; margin-top: 17px; color: ${muted}; font-size: 13px; line-height: 1.9; }
            ul { display: grid; grid-template-columns: repeat(2, 1fr); gap: 14px 18px; padding: 0; margin: 25px 0 28px; list-style: none; }
            li { display: flex; align-items: center; gap: 8px; color: ${muted}; font-size: 12px; .anticon { color: ${accent}; font-size: 10px; } }
            @media (max-width: 575px) { h3 { font-size: 26px; } p { font-size: 12px; } li { font-size: 11px; } }
        `,
        productDiagram: css`
            position: relative; align-self: center; min-height: 280px; border-radius: 50%; background: radial-gradient(ellipse at center, ${token.colorPrimaryBg}, transparent 68%);
            &::before, &::after { position: absolute; top: 50%; left: 15%; width: 70%; height: 1px; background: ${token.colorPrimaryBorder}; content: ""; }
            &::before { transform: rotate(32deg); } &::after { transform: rotate(-32deg); }
            @media (max-width: 767px) { min-height: 230px; width: min(100%, 400px); margin-inline: auto; }
        `,
        diagramCenter: css`
            position: absolute; z-index: 2; top: 50%; left: 50%; display: grid; place-items: center; width: 94px; height: 94px; border: 1px solid ${token.colorPrimaryBorder}; border-radius: 22px;
            color: ${accent}; background: ${paper}; transform: translate(-50%, -50%); box-shadow: 0 12px 24px -16px color-mix(in srgb, ${accent} 20%, transparent);
            svg { width: 43px; height: 43px; }
            @media (max-width: 575px) { width: 74px; height: 74px; border-radius: 18px; svg { width: 36px; height: 36px; } }
        `,
        diagramNode: css`
            position: absolute; z-index: 1; display: flex; align-items: center; justify-content: center; gap: 7px; min-width: 110px; height: 37px; padding: 0 12px; border: 1px solid ${line}; border-radius: 6px; color: ${muted}; background: ${paper}; font-size: 11px;
            i { width: 4px; height: 4px; border-radius: 50%; background: ${accent}; opacity: .6; }
            &[data-position="0"] { left: 0; top: 14%; } &[data-position="1"] { right: 0; top: 14%; }
            &[data-position="2"] { left: 0; bottom: 14%; } &[data-position="3"] { right: 0; bottom: 14%; }
            @media (max-width: 575px) { min-width: 100px; padding-inline: 8px; font-size: 10px; }
        `,
        manifesto: css`
            text-align: center;
            h2 { margin-top: 20px; font-size: clamp(30px, 3.5vw, 48px); font-weight: 550; line-height: 1.5; letter-spacing: -1.5px; span { color: ${muted}; } }
            p { margin-top: 22px; color: ${muted}; font-size: 13px; line-height: 2; }
            @media (max-width: 767px) { h2 { font-size: 30px; span { display: block; } } p { font-size: 12px; } }
        `,
        advantages: css`
            display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 32px; margin-top: 62px;
            article { padding-top: 23px; border-top: 1px solid ${line}; } span { color: ${accent}; font-size: 10px; }
            h3 { margin-top: 18px; font-size: 16px; font-weight: 550; } p { margin-top: 10px; color: ${muted}; font-size: 12px; line-height: 1.85; }
            @media (max-width: 767px) { grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 30px; margin-top: 40px; }
        `,
        testimonials: css`padding-top: 0;`,
        testimonialGrid: css`
            display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 35px;
            article { display: flex; flex-direction: column; padding: 24px 0 0; border-top: 1px solid ${line}; }
            blockquote { margin: 3px 0 26px; color: ${muted}; font-size: 13px; line-height: 2; }
            @media (max-width: 767px) { grid-template-columns: 1fr; gap: 28px; }
        `,
        quoteMark: css`height: 44px; color: ${accent}; opacity: .5; font-size: 54px; font-family: Georgia, serif; line-height: 1;`,
        author: css`
            display: flex; align-items: center; gap: 11px; margin-top: auto;
            > span { display: grid; place-items: center; width: 34px; height: 34px; border: 1px solid ${line}; border-radius: 50%; color: ${muted}; background: ${wash}; font-size: 12px; }
            strong { display: block; font-size: 12px; font-weight: 500; } small { display: block; margin-top: 3px; color: ${muted}; font-size: 10px; }
        `,
        cta: css`
            display: flex; align-items: center; justify-content: space-between; gap: 35px; padding: 50px 0; border-top: 1px solid ${line}; border-bottom: 1px solid ${line};
            h2 { margin-top: 16px; font-size: 38px; font-weight: 550; letter-spacing: -1px; span { color: ${accent}; } } p { margin-top: 17px; color: ${muted}; font-size: 13px; }
            @media (max-width: 767px) { flex-direction: column; align-items: flex-start; padding-block: 34px; gap: 28px; h2 { font-size: 30px; } p { font-size: 12px; } }
        `,
        ctaActions: css`
            display: flex; flex-shrink: 0; gap: 12px;
            @media (max-width: 991px) { flex-direction: column; }
            @media (max-width: 767px) { flex-direction: row; }
            @media (max-width: 575px) { width: 100%; button { flex: 1; padding-inline: 15px; height: 44px; font-size: 13px; gap: 10px; } }
        `,
        footer: css`
            display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 18px; padding-block: 30px; color: ${muted}; font-size: 10px;
            > div { display: flex; flex-wrap: wrap; align-items: center; gap: 12px; } a { color: inherit; text-decoration: none; &:hover { color: ${accent}; } }
            @media (max-width: 575px) { padding-block: 25px; gap: 12px; font-size: 9px; }
        `,
        backToTop: css`
            position: fixed; z-index: 40; right: 24px; bottom: 24px; display: grid; place-items: center; width: 38px; height: 38px; border: 1px solid ${line}; border-radius: 8px; color: ${muted}; background: ${paper}; opacity: 0; visibility: hidden; transition: opacity .2s, color .2s;
            &[data-visible="true"] { opacity: 1; visibility: visible; } &:hover { color: ${accent}; }
            @media (max-width: 575px) { right: 18px; bottom: 18px; width: 34px; height: 34px; }
        `,
    };
});

type HomeIconType = "deploy" | "shield" | "speed" | "module" | "support" | "global" | "monitor" | "api" | "safe" | "ops";

const HomeIcon = ({type}: {type: HomeIconType}) => {
    const line = {stroke: "currentColor", strokeWidth: 1.8, strokeLinecap: "round" as const, strokeLinejoin: "round" as const};
    let content: React.ReactNode;
    switch (type) {
        case "deploy":
            content = <><path d="M24 6 38 14v16L24 38 10 30V14Z"/><path d="m10 14 14 8 14-8M24 22v16M17 42h14M24 38v4"/></>;
            break;
        case "shield":
            content = <><path d="m24 5 14 5v10.5C38 31 32.5 38.5 24 43c-8.5-4.5-14-12-14-22.5V10Z"/><path d="m18 24 4 4 8-10"/></>;
            break;
        case "speed":
            content = <><path d="M9 34a15 15 0 1 1 30 0M14 34h20M24 34l8-13M15 21h5M28 21h5M24 13v5"/><circle cx="24" cy="34" r="3"/></>;
            break;
        case "module":
            content = <><rect x="8" y="8" width="13" height="13" rx="2"/><rect x="27" y="8" width="13" height="13" rx="2"/><rect x="8" y="27" width="13" height="13" rx="2"/><rect x="27" y="27" width="13" height="13" rx="2"/></>;
            break;
        case "support":
            content = <><circle cx="24" cy="18" r="8"/><path d="M10 41c2.5-7 7.2-11 14-11s11.5 4 14 11M9 23v-3a15 15 0 0 1 30 0v6a5 5 0 0 1-5 5h-4"/></>;
            break;
        case "global":
            content = <><circle cx="24" cy="24" r="17"/><path d="M7 24h34M24 7c5 4.5 7.5 10.2 7.5 17S29 36.5 24 41M24 7c-5 4.5-7.5 10.2-7.5 17S19 36.5 24 41"/></>;
            break;
        case "monitor":
            content = <><rect x="7" y="9" width="34" height="24" rx="3"/><path d="M17 39h14M24 33v6m-9-15 5-5 5 5 7-9"/></>;
            break;
        case "api":
            content = <><rect x="7" y="10" width="34" height="28" rx="3"/><path d="m17 20-5 4 5 4m14-8 5 4-5 4m-5-12-4 16"/></>;
            break;
        case "safe":
            content = <><path d="M12 22v-5a12 12 0 0 1 24 0v5"/><rect x="9" y="22" width="30" height="20" rx="3"/><path d="M24 30v5"/></>;
            break;
        case "ops":
            content = <><path d="M24 7v6M24 35v6M7 24h6M35 24h6m-23-12 4.2 4.2m15.6 15.6L36 36m0-24-4.2 4.2M16.2 31.8 12 36"/><circle cx="24" cy="24" r="9"/><path d="m21 24 2 2 5-5"/></>;
            break;
    }
    return <svg viewBox="0 0 48 48" fill="none" aria-hidden="true" {...line}>{content}</svg>;
};

const navigation = [
    {id: "home", name: "首页"},
    {id: "features", name: "平台优势"},
    {id: "products", name: "产品能力"},
    {id: "testimonials", name: "用户评价"},
];

const features: {icon: HomeIconType; title: string; description: string}[] = [
    {icon: "deploy", title: "快速部署", description: "分钟级快速部署，无需复杂配置，开箱即用，让您的项目快速上线运营。"},
    {icon: "shield", title: "安全防护", description: "多重安全防护体系，数据加密传输，符合行业安全标准，保障业务安全。"},
    {icon: "speed", title: "高性能", description: "采用先进的技术架构，毫秒级响应速度，支持高并发访问，稳定可靠。"},
    {icon: "module", title: "弹性扩展", description: "灵活的扩展机制，根据业务需求自动调整资源，应对流量高峰。"},
    {icon: "support", title: "专业支持", description: "7×24小时专业技术支持，快速响应您的需求，助力业务持续发展。"},
    {icon: "global", title: "全球加速", description: "全球CDN加速节点，智能路由分发，无论用户在哪里都能获得极速体验。"},
];

const products: {icon: HomeIconType; title: string; description: string; features: string[]}[] = [
    {icon: "monitor", title: "智能监控", description: "实时监控系统运行状态，智能告警，让您随时掌握系统健康状况。", features: ["实时数据监控", "智能告警通知", "可视化报表", "性能分析"]},
    {icon: "api", title: "API 服务", description: "提供丰富的API接口，文档完善，便于集成，助力快速开发。", features: ["RESTful API", "完善文档", "SDK支持", "高并发处理"]},
    {icon: "safe", title: "数据安全", description: "企业级数据安全方案，多重备份机制，确保数据万无一失。", features: ["数据加密", "定期备份", "灾难恢复", "权限管理"]},
    {icon: "ops", title: "自动化运维", description: "自动化部署与运维工具，提升运维效率，降低人工成本。", features: ["自动部署", "日志管理", "版本控制", "回滚机制"]},
];

const advantages = [
    {title: "简单易用", describe: "直观的操作界面，降低使用门槛"},
    {title: "灵活定制", describe: "支持个性化配置，满足不同需求"},
    {title: "持续创新", describe: "不断更新迭代，引领技术潮流"},
    {title: "生态完善", describe: "丰富的生态系统，无缝对接"},
];

const testimonials = [
    {content: "使用这个平台后，我们的开发效率提升了3倍，服务稳定性也大幅提高。团队可以更专注于业务开发，而不是运维工作。强烈推荐！", name: "张伟", role: "技术总监"},
    {content: "非常专业的服务，技术支持响应及时，功能强大且易用。帮助我们快速搭建了业务系统，节省了大量时间和成本。", name: "李娜", role: "产品经理"},
    {content: "安全性和稳定性都非常出色，API文档清晰易懂，集成起来很顺畅。是我用过最好的服务平台之一。", name: "王强", role: "架构师"},
];

export default () => {
    const {styles, cx} = useStyles();
    const {info} = useModel("web");
    const user = useModel("user");
    const theme = useTheme();
    const pageRef = useRef<HTMLDivElement>(null);
    const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
    const [scrolled, setScrolled] = useState(false);
    const [showBackToTop, setShowBackToTop] = useState(false);
    const [activeSection, setActiveSection] = useState("home");
    const [productIndex, setProductIndex] = useState(0);
    const product = products[productIndex];
    const site = info as typeof info & {logo?: string; describe?: string; icp?: string; police_beian?: string};
    const brand = site?.name || Settings.title || "云上豁者";
    const themeIcon = theme?.isDarkMode() ? "icon-dark" : theme?.isAutoMode() ? "icon-auto" : "icon-light";
    const themeName = theme?.getModeName(theme.mode) || "跟随系统";

    useEffect(() => {
        let frame = 0;
        const updateScroll = () => {
            frame = 0;
            setScrolled(window.scrollY > 30);
            setShowBackToTop(window.scrollY > 450);
            let current = "home";
            navigation.forEach(item => {
                const element = document.getElementById(item.id);
                if (element && element.getBoundingClientRect().top <= window.innerHeight * .3) current = item.id;
            });
            setActiveSection(current);
        };
        const onScroll = () => {
            if (!frame) frame = window.requestAnimationFrame(updateScroll);
        };
        const onResize = () => {
            if (window.innerWidth >= 768) setMobileMenuOpen(false);
            onScroll();
        };
        updateScroll();
        window.addEventListener("scroll", onScroll, {passive: true});
        window.addEventListener("resize", onResize);
        return () => {
            window.removeEventListener("scroll", onScroll);
            window.removeEventListener("resize", onResize);
            if (frame) window.cancelAnimationFrame(frame);
        };
    }, []);

    useEffect(() => {
        const page = pageRef.current;
        if (!page) return;
        const reducedMotion = window.matchMedia("(prefers-reduced-motion: reduce)");
        const revealObserver = new IntersectionObserver(entries => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    entry.target.setAttribute("data-visible", "true");
                    revealObserver.unobserve(entry.target);
                }
            });
        }, {threshold: .08});
        if (!reducedMotion.matches) {
            page.setAttribute("data-motion", "true");
            page.querySelectorAll(".home-reveal").forEach(element => revealObserver.observe(element));
        }
        const onMotionChange = () => {
            if (reducedMotion.matches) page.removeAttribute("data-motion");
        };
        reducedMotion.addEventListener("change", onMotionChange);
        return () => {
            revealObserver.disconnect();
            reducedMotion.removeEventListener("change", onMotionChange);
        };
    }, []);

    useEffect(() => {
        const page = pageRef.current;
        if (!page) return;
        const hoverMotion = window.matchMedia("(hover: hover) and (pointer: fine) and (prefers-reduced-motion: no-preference)");
        let target: HTMLElement | null = null;
        let frame = 0;
        let x = 0;
        let y = 0;
        const reset = () => {
            if (frame) window.cancelAnimationFrame(frame);
            frame = 0;
            target?.removeAttribute("data-highlighted");
            target = null;
        };
        const update = () => {
            frame = 0;
            if (!target) return;
            const rect = target.getBoundingClientRect();
            target.style.setProperty("--spotlight-x", `${x - rect.left}px`);
            target.style.setProperty("--spotlight-y", `${y - rect.top}px`);
            target.setAttribute("data-highlighted", "true");
        };
        const onPointerMove = (event: PointerEvent) => {
            if (!hoverMotion.matches || event.pointerType !== "mouse") {
                reset();
                return;
            }
            const next = event.target instanceof Element ? event.target.closest<HTMLElement>("[data-spotlight]") : null;
            if (next !== target) {
                reset();
                target = next;
            }
            if (!target) return;
            x = event.clientX;
            y = event.clientY;
            if (!frame) frame = window.requestAnimationFrame(update);
        };
        page.addEventListener("pointermove", onPointerMove, {passive: true});
        page.addEventListener("pointerleave", reset);
        page.addEventListener("pointercancel", reset);
        window.addEventListener("scroll", reset, {passive: true});
        window.addEventListener("resize", reset);
        window.addEventListener("blur", reset);
        document.addEventListener("visibilitychange", reset);
        hoverMotion.addEventListener("change", reset);
        return () => {
            reset();
            page.removeEventListener("pointermove", onPointerMove);
            page.removeEventListener("pointerleave", reset);
            page.removeEventListener("pointercancel", reset);
            window.removeEventListener("scroll", reset);
            window.removeEventListener("resize", reset);
            window.removeEventListener("blur", reset);
            document.removeEventListener("visibilitychange", reset);
            hoverMotion.removeEventListener("change", reset);
        };
    }, []);

    const scrollToSection = (id: string) => {
        const element = document.getElementById(id);
        if (element) {
            window.scrollTo({
                top: Math.max(0, element.getBoundingClientRect().top + window.scrollY - (window.innerWidth < 768 ? 64 : 76)),
                behavior: window.matchMedia("(prefers-reduced-motion: reduce)").matches ? "auto" : "smooth",
            });
        }
        setMobileMenuOpen(false);
    };

    const handleGetStarted = () => {
        setMobileMenuOpen(false);
        historyPush(user?.web ? "user.index" : "user.login");
    };

    const handleLogout = async () => {
        setMobileMenuOpen(false);
        try {
            await logout({
                onSuccess: () => {
                    deleteHeader();
                    user?.setWeb(undefined);
                },
                onFail: () => {},
            });
        } catch {}
    };

    return (
        <Body space={false} animation={false} className={styles.body}>
            <div ref={pageRef} className={styles.page}>
                <header className={styles.header} data-scrolled={scrolled}>
                    <div className={cx(styles.container, styles.headerInner)}>
                        <button className={styles.logo} onClick={() => scrollToSection("home")} aria-label="返回首页顶部">
                            {site?.logo ? <img src={site.logo} alt=""/> : <Icon type="icon-logo"/>}
                            <span>{brand}</span>
                        </button>
                        <nav className={styles.nav} aria-label="首页导航">
                            {navigation.map(item => <button key={item.id} onClick={() => scrollToSection(item.id)} aria-current={activeSection === item.id}>{item.name}</button>)}
                        </nav>
                        <div className={styles.headerActions}>
                            <button className={styles.themeButton} onClick={() => theme?.toggle()} aria-label={themeName} title={themeName}><Icon type={themeIcon}/></button>
                            <button className={styles.textButton} onClick={user?.web ? handleLogout : () => historyPush("user.login")}>{user?.web ? "退出" : "登录"}</button>
                            <button className={styles.primaryButton} data-small="true" onClick={handleGetStarted}>{user?.web ? "进入系统" : "开始使用"}<Icon type="ArrowRightOutlined"/></button>
                        </div>
                        <button className={styles.mobileToggle} onClick={() => setMobileMenuOpen(true)} aria-label="打开导航菜单" aria-expanded={mobileMenuOpen}><Icon type="MenuOutlined"/></button>
                    </div>
                </header>

                <Drawer title={brand} size={360} open={mobileMenuOpen} onClose={() => setMobileMenuOpen(false)} classNames={{body: styles.mobileMenuBody}}>
                    <nav aria-label="移动端首页导航">
                        {navigation.map(item => <button key={item.id} onClick={() => scrollToSection(item.id)}>{item.name}<Icon type="ArrowRightOutlined"/></button>)}
                    </nav>
                    <div className={styles.mobileTheme}>
                        <span>主题模式</span>
                        <div>
                            <button aria-pressed={theme?.isLightMode()} onClick={() => theme?.setLightMode()}><Icon type="icon-light"/>亮色</button>
                            <button aria-pressed={theme?.isDarkMode()} onClick={() => theme?.setDarkMode()}><Icon type="icon-dark"/>暗色</button>
                            <button aria-pressed={theme?.isAutoMode()} onClick={() => theme?.setAutoMode()}><Icon type="icon-auto"/>自动</button>
                        </div>
                    </div>
                    <div className={styles.mobileActions}>
                        <button className={styles.primaryButton} onClick={handleGetStarted}>{user?.web ? "进入系统" : "开始使用"}<Icon type="ArrowRightOutlined"/></button>
                        {user?.web && <button className={styles.secondaryButton} onClick={handleLogout}>退出登录</button>}
                    </div>
                </Drawer>

                <main>
                    <section id="home" className={styles.hero} data-spotlight="hero">
                        <div className={styles.spotlight} aria-hidden="true"/>
                        <div className={cx(styles.container, styles.heroInner)}>
                            <div className={styles.heroCopy}>
                                <div className={styles.heroEyebrow}><span/>A NEW WAY TO CONNECT</div>
                                <h1>连接你的世界<span>让可能发生</span></h1>
                                <p>{site?.describe || "连接灵感、服务与日常。以简单的体验，开启数字世界的更多可能。"}</p>
                                <div className={styles.heroActions}>
                                    <button className={styles.primaryButton} onClick={handleGetStarted}>{user?.web ? "进入工作空间" : "开始体验"}<Icon type="ArrowRightOutlined"/></button>
                                    <button className={styles.secondaryButton} onClick={() => scrollToSection("products")}>探索产品<Icon type="ArrowDownOutlined"/></button>
                                </div>
                                <div className={styles.heroNotes}><span><Icon type="CheckOutlined"/>轻盈体验</span><span><Icon type="CheckOutlined"/>安全可靠</span><span><Icon type="CheckOutlined"/>持续进化</span></div>
                            </div>
                            <div className={styles.heroArt} aria-hidden="true">
                                <HomeVisual className={styles.visual}/>
                                <div className={styles.artCaption}><span>IDEAS. CONNECTED.</span><span>构建无限可能</span></div>
                            </div>
                        </div>
                        <div className={cx(styles.container, styles.heroBottom)}>
                            <span>TECHNOLOGY, MADE HUMAN.</span>
                            <div><span><b>01</b>高效连接</span><span><b>02</b>安全守护</span><span><b>03</b>自由扩展</span></div>
                            <button onClick={() => scrollToSection("features")} aria-label="探索平台优势"><Icon type="ArrowDownOutlined"/></button>
                        </div>
                    </section>

                    <section id="features" className={cx(styles.container, styles.section)}>
                        <div className={cx(styles.sectionHeading, "home-reveal")}>
                            <div><span className={styles.eyebrow}>01 / BUILT TO GO FURTHER</span><h2>强大的内核。<br/><span>轻盈的体验。</span></h2></div>
                            <p>每一个顺畅的瞬间，背后都有可靠的技术。<br/>从快速开始，到从容进阶。</p>
                        </div>
                        <div className={styles.featureGrid}>
                            {features.map((feature, index) => <article key={feature.title} className="home-reveal" data-spotlight="card"><div className={styles.spotlight} aria-hidden="true"/><div className={styles.featureTop}><HomeIcon type={feature.icon}/><span>0{index + 1}</span></div><h3>{feature.title}</h3><p>{feature.description}</p></article>)}
                        </div>
                    </section>

                    <section id="products" className={styles.productsSection}>
                        <div className={cx(styles.container, styles.section)}>
                            <div className={cx(styles.sectionHeading, "home-reveal")}>
                                <div><span className={styles.eyebrow}>02 / AN OPEN ECOSYSTEM</span><h2>让能力相连，<span>让场景延伸。</span></h2></div>
                                <p>让每个场景都有合适的工具，<br/>让每一步都更从容。</p>
                            </div>
                            <div className={cx(styles.productWorkspace, "home-reveal")}>
                                <div className={styles.productMenu} aria-label="选择产品能力">
                                    {products.map((item, index) => <button key={item.title} onClick={() => setProductIndex(index)} aria-pressed={productIndex === index} aria-controls="home-product-detail"><HomeIcon type={item.icon}/><span>{item.title}</span><Icon type="ArrowRightOutlined"/></button>)}
                                </div>
                                <div id="home-product-detail" className={styles.productDetail}>
                                    <div className={styles.productCopy}><span className={styles.eyebrow}>0{productIndex + 1} / CONNECTED SERVICES</span><h3>{product.title}</h3><p>{product.description}</p><ul>{product.features.map(feature => <li key={feature}><Icon type="CheckOutlined"/>{feature}</li>)}</ul><button className={styles.textButton} onClick={handleGetStarted}>开始使用<Icon type="ArrowRightOutlined"/></button></div>
                                    <div className={styles.productDiagram} aria-hidden="true"><div className={styles.diagramCenter}><HomeIcon type={product.icon}/></div>{product.features.map((feature, index) => <div key={feature} className={styles.diagramNode} data-position={index}><i/> {feature}</div>)}</div>
                                </div>
                            </div>
                        </div>
                    </section>

                    <section className={cx(styles.container, styles.section)}>
                        <div className={cx(styles.manifesto, "home-reveal")}>
                            <span className={styles.eyebrow}>LESS FRICTION. MORE FOCUS.</span>
                            <h2>好的工具，<span>让你忘记工具。</span></h2>
                            <p>直观的体验、可靠的服务，与持续生长的能力。<br/>陪你把日常的小想法，变成下一步的大可能。</p>
                        </div>
                        <div className={styles.advantages}>
                            {advantages.map((item, index) => <article key={item.title} className="home-reveal"><span>0{index + 1}</span><h3>{item.title}</h3><p>{item.describe}</p></article>)}
                        </div>
                    </section>

                    <section id="testimonials" className={cx(styles.container, styles.section, styles.testimonials)}>
                        <div className={cx(styles.sectionHeading, "home-reveal")}>
                            <div><span className={styles.eyebrow}>IN GOOD COMPANY</span><h2>体验，由使用者诉说。</h2></div>
                            <p>好的体验，体现在每一次使用中。</p>
                        </div>
                        <div className={styles.testimonialGrid}>
                            {testimonials.map(item => <article key={item.name} className="home-reveal">
                                <span className={styles.quoteMark} aria-hidden="true">“</span>
                                <blockquote>{item.content}</blockquote>
                                <div className={styles.author}><span>{item.name.slice(0, 1)}</span><div><strong>{item.name}</strong><small>{item.role}</small></div></div>
                            </article>)}
                        </div>
                    </section>

                    <div className={styles.container}>
                        <section className={cx(styles.cta, "home-reveal")}>
                            <div><span className={styles.eyebrow}>MAKE ROOM FOR WHAT'S NEXT</span><h2>下一步，<span>一起向前。</span></h2><p>让想法有处安放，让可能自在生长。</p></div>
                            <div className={styles.ctaActions}>
                                <button className={styles.primaryButton} onClick={handleGetStarted}>{user?.web ? "进入控制台" : "立即开始"}<Icon type="ArrowRightOutlined"/></button>
                                <button className={styles.secondaryButton} onClick={() => window.location.href = `mailto:${site?.contact || Settings.contact}`}>咨询合作</button>
                            </div>
                        </section>
                        <footer className={styles.footer}>
                            <span>© {new Date().getFullYear()} {brand}. All rights reserved.</span>
                            <div>
                                {site?.icp && <a href="https://beian.miit.gov.cn" target="_blank" rel="noopener noreferrer">{site.icp}</a>}
                                {site?.police_beian && <a href="http://www.beian.gov.cn/portal/registerSystemInfo" target="_blank" rel="noopener noreferrer">{site.police_beian}</a>}
                            </div>
                        </footer>
                    </div>
                </main>
                <button className={styles.backToTop} data-visible={showBackToTop} onClick={() => scrollToSection("home")} aria-label="返回顶部" title="返回顶部"><Icon type="ArrowUpOutlined"/></button>
            </div>
        </Body>
    );
};
