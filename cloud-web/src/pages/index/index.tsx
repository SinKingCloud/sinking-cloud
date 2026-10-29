import React, {useEffect, useId, useRef, useState} from "react";
import {Drawer} from "antd";
import {useModel} from "umi";
import {Body, Icon, useTheme} from "sinking-antd";
import Settings from "@/../config/defaultSettings";
import {deleteHeader, getLoginToken} from "@/utils/auth";
import {logout} from "@/service/common/login";
import {historyPush} from "@/utils/route";
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
            &[data-motion="true"] .home-reveal:not([data-visible="true"]) { opacity: 0; translate: 0 24px; }
            &[data-motion="true"] .home-reveal[data-visible="true"] { animation: home-reveal-rise .8s cubic-bezier(.22, 1, .36, 1) var(--reveal-delay, 0s) backwards; }
            &[data-motion="true"] .home-reveal:focus-within { opacity: 1; translate: none; animation: none; }
            .home-reveal[data-step="1"] { --reveal-delay: .08s; }
            .home-reveal[data-step="2"] { --reveal-delay: .16s; }
            .home-reveal[data-step="3"] { --reveal-delay: .24s; }
            @keyframes home-enter-fade {
                from { opacity: 0; }
                to { opacity: 1; }
            }
            @keyframes home-enter-rise {
                from { opacity: 0; transform: translate3d(0, 20px, 0); }
                to { opacity: 1; transform: none; }
            }
            @keyframes home-reveal-rise {
                from { opacity: 0; translate: 0 24px; }
                to { opacity: 1; translate: none; }
            }
            @keyframes home-enter-visual {
                from { opacity: 0; transform: translate3d(0, 24px, 0) scale(.97); }
                to { opacity: 1; transform: none; }
            }
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
        headerInner: css`
            display: flex; align-items: center; justify-content: space-between; gap: 34px; height: 100%;
            animation: home-enter-fade .5s ease-out backwards;
            &:focus-within { animation: none; }
        `,
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
            position: relative; overflow: hidden; display: inline-flex; align-items: center; justify-content: center; gap: 22px; height: 48px; border: 1px solid ${accent}; border-radius: ${token.borderRadius}px; padding: 0 25px;
            color: ${token.colorTextLightSolid}; background: ${accent}; font-size: 14px; font-weight: 500; transition: background .18s, border-color .18s, box-shadow .2s;
            &::before { position: absolute; inset: -30% auto -30% -65%; width: 45%; background: linear-gradient(90deg, transparent, rgba(255,255,255,.23), transparent); transform: skewX(-22deg); transition: left .5s ease; content: ""; }
            &:hover { background: ${token.colorPrimaryHover}; border-color: ${token.colorPrimaryHover}; box-shadow: 0 6px 20px -9px color-mix(in srgb, ${accent} 55%, transparent); &::before { left: 125%; } }
            &:active { background: ${token.colorPrimaryActive}; border-color: ${token.colorPrimaryActive}; }
            .anticon { font-size: 13px; }
            &[data-small="true"] { height: 36px; gap: 12px; padding: 0 17px; font-size: 12px; }
        `,
        secondaryButton: css`
            display: inline-flex; align-items: center; justify-content: center; gap: 14px; height: 48px; border: 1px solid ${line}; border-radius: ${token.borderRadius}px; padding: 0 24px;
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
            button { display: inline-flex; align-items: center; justify-content: center; flex: 1; gap: 7px; height: 38px; border: 1px solid ${line}; border-radius: ${token.borderRadius}px; color: ${muted}; background: transparent; font-size: 13px;
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
            > * { animation: home-enter-rise .7s cubic-bezier(.22, 1, .36, 1) backwards; }
            > :focus-within { animation: none; }
            > :nth-child(1) { animation-delay: .05s; }
            > :nth-child(2) { animation-delay: .12s; }
            > :nth-child(3) { animation-delay: .22s; }
            > :nth-child(4) { animation-delay: .32s; }
            > :nth-child(5) { animation-delay: .4s; }
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
            animation: home-enter-visual 1s cubic-bezier(.22, 1, .36, 1) .15s backwards;
            @media (max-width: 1100px) { width: calc(100% + 60px); margin-inline: -20px -40px; }
            @media (max-width: 767px) { width: min(115%, 500px); max-width: none; justify-self: center; margin: -14px 0 -15px; }
        `,
        heroVisual: css`
            --core-surface: ${isDarkMode ? token.colorBgElevated : "white"};
            --core-tint: color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 12 : 5}%, var(--core-surface));
            --core-soft: color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 24 : 14}%, var(--core-surface));
            --core-light: color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 40 : 28}%, var(--core-surface));
            --core-border: color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 60 : 42}%, var(${isDarkMode ? "--core-highlight" : "--core-surface"}));
            --core-line: color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 35 : 58}%, var(${isDarkMode ? "--core-highlight" : "--core-surface"}));
            --core-highlight: ${isDarkMode ? `color-mix(in srgb, ${token.colorPrimaryHover} 70%, ${token.colorText})` : "white"};
            --core-primary: ${token.colorPrimary};
            --core-bright: ${token.colorPrimaryHover};
            --core-deep: ${token.colorPrimaryActive};
            position: relative;
            width: 100%;
            max-height: min(650px, calc(100svh - 250px));
            aspect-ratio: 1.05;
            pointer-events: none;
            user-select: none;

            svg { display: block; width: 100%; height: 100%; overflow: visible; }
            .core-top { animation: home-core-top 7s ease-in-out infinite; }
            .core-middle { animation: home-core-middle 7s ease-in-out infinite; }
            .core-scan { animation: home-core-scan 7s ease-in-out infinite; }
            .core-pulse { animation: home-core-pulse 5s ease-in-out infinite; }
            .core-signal { animation: home-core-signal 5s linear infinite; }
            .core-signal-secondary { animation-duration: 6.5s; animation-delay: -3s; }
            .core-node { animation: home-core-node 5s ease-in-out infinite; }
            &[data-running="false"] * { animation-play-state: paused !important; }

            @keyframes home-core-top {
                0%, 100% { transform: translateY(0); }
                50% { transform: translateY(-12px); }
            }
            @keyframes home-core-middle {
                0%, 100% { transform: translateY(0); }
                50% { transform: translateY(7px); }
            }
            @keyframes home-core-scan {
                0%, 15% { transform: translateY(-110px); opacity: 0; }
                25% { opacity: .8; }
                70% { opacity: .8; }
                85%, 100% { transform: translateY(110px); opacity: 0; }
            }
            @keyframes home-core-pulse {
                0%, 100% { opacity: .35; }
                50% { opacity: .85; }
            }
            @keyframes home-core-signal {
                from { stroke-dashoffset: 340; }
                to { stroke-dashoffset: -340; }
            }
            @keyframes home-core-node {
                0%, 20%, 80%, 100% { opacity: .45; }
                45%, 60% { opacity: 1; }
            }
            @media (prefers-reduced-motion: reduce) {
                * { animation: none !important; }
                .core-scan { opacity: 0; }
            }
            @media (max-width: 767px) { max-height: 350px; }
        `,
        artCaption: css`
            position: absolute; bottom: 3%; left: 17%; right: 13%; display: flex; justify-content: space-between; align-items: center; gap: 15px; padding-top: 13px; border-top: 1px solid color-mix(in srgb, ${accent} 12%, transparent); color: ${muted}; font-size: 8px; letter-spacing: 1.5px;
            > span:last-child { color: ${accent}; font-size: 9px; letter-spacing: .7px; }
            @media (max-width: 767px) { display: none; }
        `,
        heroBottom: css`
            position: relative; display: flex; align-items: center; justify-content: space-between; gap: 25px; min-height: 86px; color: ${muted}; font-size: 10px;
            animation: home-enter-fade .7s ease-out .45s backwards;
            &:focus-within { animation: none; }
            &::before { position: absolute; top: 0; left: clamp(24px, 5vw, 80px); right: clamp(24px, 5vw, 80px); height: 1px; background: ${line}; content: ""; }
            > span { font-size: 8px; letter-spacing: 1.8px; }
            > div { display: flex; gap: 42px; span { display: flex; align-items: center; gap: 13px; } b { color: ${accent}; font-size: 9px; font-weight: 400; font-family: monospace; } }
            button { display: grid; place-items: center; width: 32px; height: 32px; border: 1px solid ${line}; border-radius: 50%; color: ${muted}; background: transparent; transition: border-color .2s, color .2s; &:hover { color: ${accent}; border-color: ${accent}; } }
            @media (max-width: 767px) { min-height: 61px; > div { gap: 22px; span { gap: 7px; } } > span { display: none; } font-size: 10px; }
            @media (max-width: 360px) { > div { gap: 14px; } }
        `,
        section: css`
            padding-block: 96px;
            @media (max-width: 991px) { padding-block: 80px; }
            @media (max-width: 575px) { padding-block: 64px; }
        `,
        eyebrow: css`
            display: inline-flex; align-items: center; gap: 10px; color: ${accent}; font-size: 9px; line-height: 1.6; letter-spacing: 1.8px; font-weight: 600;
            &::before { width: 18px; height: 1px; background: currentColor; opacity: .55; content: ""; }
        `,
        sectionHeading: css`
            display: flex; align-items: flex-end; justify-content: space-between; gap: 35px; margin-bottom: 44px;
            h2 { margin-top: 18px; font-size: clamp(30px, 3vw, 42px); font-weight: 650; line-height: 1.4; letter-spacing: -1.3px; span { color: ${muted}; } }
            > p { max-width: 320px; padding-bottom: 4px; color: ${muted}; font-size: 13px; line-height: 1.9; }
            @media (max-width: 767px) { display: block; margin-bottom: 30px; h2 { font-size: 30px; letter-spacing: -.8px; } > p { margin-top: 20px; font-size: 12px; } }
        `,
        featuresSection: css`
            position: relative; isolation: isolate;
            background: linear-gradient(180deg, ${paper}, color-mix(in srgb, ${accent} ${isDarkMode ? 5 : 2.5}%, ${paper}));
            h2 { white-space: nowrap; }
            &::before { position: absolute; z-index: -1; inset: 12% 0 0; background: radial-gradient(ellipse at 12% 30%, color-mix(in srgb, ${accent} 7%, transparent), transparent 55%); pointer-events: none; content: ""; }
            @media (max-width: 767px) { h2 { font-size: clamp(22px, 6vw, 30px); } }
        `,
        featureGrid: css`
            display: grid; grid-template-columns: repeat(12, minmax(0, 1fr)); gap: 18px;
            article {
                position: relative; grid-column: span 3; min-width: 0; overflow: hidden; isolation: isolate; padding: 28px 25px; border: 1px solid color-mix(in srgb, ${accent} 10%, ${line}); border-radius: ${token.borderRadiusLG}px;
                background: linear-gradient(145deg, ${paper}, color-mix(in srgb, ${accent} ${isDarkMode ? 7 : 3}%, ${paper})); transition: border-color .25s, transform .25s, box-shadow .25s;
                &:nth-child(-n+2) { grid-column: span 6; min-height: 242px; padding: 32px; background: linear-gradient(120deg, ${paper} 15%, color-mix(in srgb, ${accent} ${isDarkMode ? 15 : 8}%, ${paper})); }
                &:nth-child(-n+2) h3 { position: relative; margin-top: 24px; font-size: 21px; }
                &:nth-child(-n+2) p { position: relative; max-width: 65%; }
            }
            @media (hover: hover) and (pointer: fine) { article:hover { border-color: ${token.colorPrimaryBorder}; transform: translateY(-4px); box-shadow: 0 16px 38px -28px color-mix(in srgb, ${accent} 45%, transparent); } }
            h3 { margin-top: 25px; font-size: 16px; font-weight: 600; }
            p { margin-top: 12px; color: ${muted}; font-size: 12px; line-height: 1.9; }
            @media (max-width: 991px) { article { grid-column: span 6; padding: 26px; } article:nth-child(-n+2) { padding: 28px; } }
            @media (max-width: 575px) {
                gap: 12px; article { padding: 24px 18px; }
                article:nth-child(-n+2) { grid-column: span 12; min-height: 215px; padding: 26px; }
                article:nth-child(-n+2) p { max-width: 68%; }
                h3 { margin-top: 20px; font-size: 15px; } p { font-size: 11px; }
            }
        `,
        featureTop: css`
            position: relative; display: flex; align-items: center; justify-content: space-between;
            > svg { width: 45px; height: 45px; padding: 8px; border: 1px solid color-mix(in srgb, ${accent} 12%, transparent); border-radius: ${token.borderRadius}px; color: ${accent}; background: color-mix(in srgb, ${accent} 5%, ${paper}); }
            > span { color: ${accent}; font-size: 10px; opacity: .5; font-family: monospace; font-variant-numeric: tabular-nums; }
        `,
        featureArt: css`
            position: absolute; right: -24px; bottom: -24px; display: grid; place-items: center; width: 210px; height: 210px; border: 1px solid color-mix(in srgb, ${accent} 13%, transparent); border-radius: 50%; color: ${accent}; pointer-events: none; transform: rotate(-18deg);
            &::before, &::after { position: absolute; inset: 25px; border: 1px solid color-mix(in srgb, ${accent} 18%, transparent); border-radius: inherit; content: ""; }
            &::after { inset: 48px; background: linear-gradient(135deg, color-mix(in srgb, ${paper} 85%, transparent), color-mix(in srgb, ${accent} 12%, transparent)); box-shadow: 0 12px 24px -18px ${accent}; }
            > svg { position: relative; z-index: 1; width: 62px; height: 62px; opacity: .5; transition: transform .6s cubic-bezier(.22, 1, .36, 1); }
            article:hover & > svg { transform: translateY(-6px) rotate(18deg); }
            @media (max-width: 1199px) { width: 180px; height: 180px; right: -40px; }
            @media (max-width: 575px) { width: 175px; height: 175px; right: -32px; bottom: -24px; > svg { width: 48px; height: 48px; } }
        `,
        productsSection: css`
            position: relative; isolation: isolate; overflow: hidden;
            background: linear-gradient(140deg, color-mix(in srgb, ${accent} ${isDarkMode ? 10 : 6}%, ${paper}), color-mix(in srgb, ${accent} 2%, ${paper}) 70%);
            &::before { position: absolute; z-index: -1; inset: 0; pointer-events: none; background-image: radial-gradient(color-mix(in srgb, ${accent} 14%, transparent) .8px, transparent .8px); background-size: 22px 22px; mask-image: linear-gradient(90deg, transparent 30%, #000); -webkit-mask-image: linear-gradient(90deg, transparent 30%, #000); content: ""; }
            @media (max-width: 767px) { h2 span { display: block; } }
        `,
        productWorkspace: css`
            border: 1px solid color-mix(in srgb, ${accent} 14%, ${line}); border-radius: ${token.borderRadiusLG}px; overflow: hidden; background: ${paper};
            box-shadow: 0 24px 70px -48px color-mix(in srgb, ${accent} 40%, transparent);
        `,
        productMenu: css`
            display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 8px; padding: 12px; border-bottom: 1px solid ${line}; background: color-mix(in srgb, ${accent} 2%, ${paper}); scrollbar-width: thin;
            button { display: flex; align-items: center; justify-content: center; gap: 11px; min-width: 0; height: 56px; padding: 0 18px; border: 1px solid transparent; border-radius: ${token.borderRadius}px; color: ${muted}; background: transparent; font-size: 13px; white-space: nowrap; transition: background .2s, color .2s, border-color .2s;
                &[aria-pressed="true"] { color: ${accent}; border-color: color-mix(in srgb, ${accent} 18%, transparent); background: ${paper}; box-shadow: 0 3px 8px -6px color-mix(in srgb, ${accent} 40%, transparent); }
                &:hover { color: ${accent}; } > svg { width: 24px; height: 24px; } > .anticon { font-size: 11px; opacity: .6; }
            }
            @media (max-width: 767px) { display: flex; overflow-x: auto; padding: 9px; gap: 6px; button { flex: none; height: 44px; padding-inline: 14px; font-size: 12px; > svg { width: 21px; height: 21px; } > .anticon { display: none; } } }
        `,
        productDetail: css`
            display: grid; grid-template-columns: .9fr 1.1fr; align-items: center; gap: 48px; padding: 46px; animation: home-enter-rise .5s cubic-bezier(.22, 1, .36, 1) backwards;
            > * { min-width: 0; }
            @media (max-width: 1199px) { gap: 28px; padding: 34px; }
            @media (max-width: 767px) { grid-template-columns: 1fr; padding: 28px; gap: 26px; }
            @media (max-width: 380px) { padding: 22px; }
        `,
        productCopy: css`
            h3 { margin-top: 17px; font-size: 30px; font-weight: 650; }
            p { max-width: 340px; margin-top: 17px; color: ${muted}; font-size: 13px; line-height: 1.9; }
            ul { display: grid; grid-template-columns: repeat(2, 1fr); gap: 14px 18px; padding: 0; margin: 25px 0 28px; list-style: none; }
            li { display: flex; align-items: center; gap: 8px; color: ${muted}; font-size: 12px; .anticon { color: ${accent}; font-size: 10px; } }
            @media (max-width: 575px) { h3 { font-size: 26px; } p { font-size: 12px; } li { font-size: 11px; } }
        `,
        productVisual: css`
            --product-accent: ${token.colorPrimary};
            --product-bright: ${token.colorPrimaryHover};
            --product-line: color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 36 : 23}%, ${token.colorBgContainer});
            --product-paper: ${token.colorBgContainer};
            position: relative;
            width: 100%;
            min-width: 0;
            align-self: center;
            container: product-visual / inline-size;
            isolation: isolate;

            &, * { box-sizing: border-box; }
            &[data-running="false"] *, &[data-running="false"] *::before, &[data-running="false"] *::after { animation-play-state: paused !important; }

            @keyframes home-product-signal {
                from { stroke-dashoffset: 105; }
                to { stroke-dashoffset: 0; }
            }
            @keyframes home-product-orbit {
                from { transform: rotate(0deg); }
                to { transform: rotate(360deg); }
            }
            @keyframes home-product-breathe {
                0%, 100% { opacity: .4; transform: scale(.95); }
                50% { opacity: .8; transform: scale(1.04); }
            }
            @media (prefers-reduced-motion: reduce) {
                *, *::before, *::after { animation: none !important; }
                .product-signal { display: none; }
            }
        `,
        productNetwork: css`
            position: relative;
            display: grid;
            grid-template-columns: repeat(2, minmax(0, 1fr));
            grid-template-rows: 58px minmax(158px, 1fr) 58px;
            column-gap: 24px;
            min-height: 310px;
            padding: 18px 8px;

            &::before {
                position: absolute; z-index: -1; inset: 0;
                background-image: linear-gradient(color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 9 : 5}%, transparent) 1px, transparent 1px), linear-gradient(90deg, color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 9 : 5}%, transparent) 1px, transparent 1px);
                background-size: 25px 25px;
                mask-image: radial-gradient(ellipse at center, black 20%, transparent 74%);
                -webkit-mask-image: radial-gradient(ellipse at center, black 20%, transparent 74%);
                content: "";
            }
            &::after {
                position: absolute; z-index: -1; inset: 12% 0;
                background: radial-gradient(ellipse at center, color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 18 : 9}%, transparent), transparent 68%);
                content: "";
            }
            @container product-visual (max-width: 330px) {
                min-height: 282px;
                grid-template-rows: 62px minmax(134px, 1fr) 62px;
                column-gap: 18px;
                padding: 12px 4px;
            }
        `,
        productConnections: css`
            position: absolute; z-index: 0; inset: 18px 8px;
            width: calc(100% - 16px); height: calc(100% - 36px);
            overflow: visible; pointer-events: none;
            .product-signal {
                stroke: var(--product-accent);
                stroke-width: 2;
                stroke-linecap: round;
                stroke-dasharray: 5 100;
                animation: home-product-signal 5.5s linear infinite;
            }
            @container product-visual (max-width: 330px) {
                inset: 12px 4px;
                width: calc(100% - 8px); height: calc(100% - 24px);
            }
        `,
        productNode: css`
            position: relative; z-index: 2;
            display: flex; align-items: center; justify-content: center; gap: 9px;
            min-width: 0; padding: 9px;
            border: 1px solid color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 24 : 14}%, ${token.colorBorderSecondary});
            border-radius: ${token.borderRadius}px;
            color: ${token.colorTextSecondary};
            background: linear-gradient(135deg, ${token.colorBgContainer}, color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 6 : 2}%, ${token.colorBgContainer}));
            box-shadow: 0 5px 12px -9px color-mix(in srgb, ${token.colorPrimary} 30%, transparent), inset 0 1px 0 color-mix(in srgb, ${token.colorTextLightSolid} ${isDarkMode ? 8 : 70}%, transparent);
            font-size: 11px; line-height: 1.5;

            > span { min-width: 0; overflow-wrap: anywhere; }
            > svg { flex: none; width: 18px; height: 18px; color: var(--product-accent); opacity: .85; }
            &::after {
                position: absolute; bottom: -3px; left: calc(50% - 2px); width: 5px; height: 5px;
                border: 1px solid var(--product-accent); border-radius: 50%; background: var(--product-paper); content: "";
            }
            &[data-position="2"]::after, &[data-position="3"]::after { top: -3px; bottom: auto; }
            @container product-visual (max-width: 330px) {
                flex-direction: column; gap: 4px; padding: 7px 5px;
                font-size: 10px; text-align: center;
                > svg { width: 16px; height: 16px; }
            }
        `,
        productCenter: css`
            position: relative; z-index: 1;
            display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 13px;
            grid-column: 1 / -1; grid-row: 2;
            min-width: 0;

            > span { position: relative; z-index: 1; max-width: 100%; padding: 3px 10px; color: ${token.colorTextSecondary}; background: color-mix(in srgb, ${token.colorBgContainer} 85%, transparent); border-radius: ${token.borderRadiusSM}px; font-size: 10px; font-weight: 500; letter-spacing: 1px; overflow-wrap: anywhere; text-align: center; }
            @container product-visual (max-width: 330px) { gap: 9px; }
        `,
        productOrbit: css`
            position: absolute; z-index: -1; top: 50%; left: 50%; width: 174px; height: 174px;
            border: 1px solid color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 18 : 10}%, transparent);
            border-radius: 50%; transform: translate(-50%, -50%); pointer-events: none;
            &::before {
                position: absolute; inset: 15px;
                border: 1px dashed color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 32 : 22}%, transparent);
                border-radius: 50%; animation: home-product-orbit 70s linear infinite; content: "";
            }
            &::after {
                position: absolute; inset: -15px; border-radius: 50%;
                background: radial-gradient(circle, color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 17 : 10}%, transparent), transparent 65%);
                animation: home-product-breathe 6s ease-in-out infinite; content: "";
            }
            > i { position: absolute; top: 26px; right: 21px; width: 5px; height: 5px; border-radius: 50%; background: var(--product-accent); box-shadow: 0 0 0 4px color-mix(in srgb, ${token.colorPrimary} 8%, transparent); opacity: .6; }
            @container product-visual (max-width: 330px) { width: 145px; height: 145px; > i { top: 20px; right: 17px; } }
        `,
        productCore: css`
            position: relative; display: grid; place-items: center; width: 76px; height: 76px;
            color: var(--product-accent);
            border: 1px solid color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 55 : 32}%, ${token.colorBgContainer});
            border-radius: ${token.borderRadiusLG}px;
            background: linear-gradient(145deg, ${token.colorBgContainer} 15%, color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 20 : 9}%, ${token.colorBgContainer}));
            box-shadow: 0 14px 24px -17px color-mix(in srgb, ${token.colorPrimary} 65%, transparent), inset 0 1px 0 color-mix(in srgb, ${token.colorTextLightSolid} ${isDarkMode ? 16 : 90}%, transparent);
            &::before { position: absolute; z-index: -1; inset: 5px -5px -6px 5px; border: 1px solid var(--product-line); border-radius: inherit; background: color-mix(in srgb, ${token.colorPrimary} ${isDarkMode ? 8 : 4}%, ${token.colorBgContainer}); transform: rotate(6deg); content: ""; }
            &::after { position: absolute; top: 7px; right: 7px; width: 3px; height: 3px; border-radius: 50%; background: var(--product-bright); content: ""; }
            > svg { width: 37px; height: 37px; }
            @container product-visual (max-width: 330px) { width: 65px; height: 65px; > svg { width: 32px; height: 32px; } }
        `,
        experienceLayout: css`
            display: grid; grid-template-columns: 1fr 1.1fr; align-items: center; gap: 80px;
            @media (max-width: 991px) { gap: 42px; }
            @media (max-width: 767px) { grid-template-columns: 1fr; gap: 40px; }
        `,
        manifesto: css`
            h2 { margin-top: 20px; font-size: clamp(26px, 3.2vw, 44px); font-weight: 650; line-height: 1.5; letter-spacing: -1.5px; white-space: nowrap; span { color: ${muted}; } }
            p { margin-top: 22px; color: ${muted}; font-size: 13px; line-height: 2; }
            @media (max-width: 767px) { h2 { font-size: clamp(22px, 6vw, 30px); letter-spacing: -.8px; } p { font-size: 12px; } }
        `,
        advantages: css`
            display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 30px;
            article { position: relative; padding: 23px 0 0 20px; border-left: 1px solid color-mix(in srgb, ${accent} 20%, ${line}); }
            article::before { position: absolute; top: 0; left: -2px; width: 3px; height: 24px; border-radius: ${token.borderRadiusXS}px; background: ${accent}; opacity: .5; content: ""; transition: height .4s; }
            article:hover::before { height: 58px; }
            span { color: ${accent}; font-size: 11px; font-family: monospace; }
            h3 { margin-top: 18px; font-size: 16px; font-weight: 600; } p { margin-top: 10px; color: ${muted}; font-size: 12px; line-height: 1.85; }
            @media (max-width: 575px) { gap: 28px 18px; article { padding-left: 16px; } }
        `,
        testimonials: css`
            position: relative; isolation: isolate; background: linear-gradient(180deg, color-mix(in srgb, ${accent} ${isDarkMode ? 6 : 3}%, ${paper}), ${paper});
            &::before { position: absolute; z-index: -1; inset: 0; background: radial-gradient(ellipse at 50% 20%, color-mix(in srgb, ${accent} 7%, transparent), transparent 60%); pointer-events: none; content: ""; }
        `,
        testimonialGrid: css`
            display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 20px;
            article { position: relative; isolation: isolate; overflow: hidden; display: flex; flex-direction: column; padding: 30px; border: 1px solid color-mix(in srgb, ${accent} 10%, ${line}); border-radius: ${token.borderRadiusLG}px; background: ${paper}; transition: border-color .25s, transform .25s; }
            article:nth-child(2) { background: linear-gradient(145deg, ${paper}, color-mix(in srgb, ${accent} ${isDarkMode ? 10 : 5}%, ${paper})); }
            @media (hover: hover) and (pointer: fine) { article:hover { transform: translateY(-4px); border-color: ${token.colorPrimaryBorder}; } }
            blockquote { margin: 3px 0 26px; color: ${muted}; font-size: 13px; line-height: 2; }
            @media (max-width: 991px) { article { padding: 25px; } blockquote { font-size: 12px; } }
            @media (max-width: 767px) { grid-template-columns: 1fr; gap: 16px; }
        `,
        quoteMark: css`height: 44px; color: ${accent}; opacity: .45; font-size: 54px; font-family: Georgia, serif; line-height: 1;`,
        author: css`
            display: flex; align-items: center; gap: 11px; margin-top: auto;
            > span { display: grid; place-items: center; width: 36px; height: 36px; border: 1px solid color-mix(in srgb, ${accent} 15%, transparent); border-radius: ${token.borderRadius}px; color: ${accent}; background: color-mix(in srgb, ${accent} 7%, ${paper}); font-size: 12px; }
            strong { display: block; font-size: 12px; font-weight: 500; } small { display: block; margin-top: 3px; color: ${muted}; font-size: 10px; }
        `,
        cta: css`
            position: relative; isolation: isolate; overflow: hidden; display: flex; align-items: center; justify-content: space-between; gap: 35px; padding: 56px; border: 1px solid color-mix(in srgb, ${accent} 18%, ${line}); border-radius: ${token.borderRadiusLG}px;
            background: linear-gradient(120deg, color-mix(in srgb, ${accent} ${isDarkMode ? 12 : 6}%, ${paper}), color-mix(in srgb, ${accent} ${isDarkMode ? 20 : 12}%, ${paper}));
            &::before { position: absolute; z-index: -1; top: -190px; right: 5%; width: 520px; height: 520px; border: 1px solid color-mix(in srgb, ${accent} 12%, transparent); border-radius: 50%; box-shadow: 0 0 0 54px color-mix(in srgb, ${accent} 3%, transparent), 0 0 0 110px color-mix(in srgb, ${accent} 3%, transparent); pointer-events: none; content: ""; }
            h2 { margin-top: 16px; font-size: 38px; font-weight: 650; letter-spacing: -1px; span { color: ${accent}; } } p { margin-top: 17px; color: ${muted}; font-size: 13px; }
            @media (max-width: 991px) { padding: 42px; }
            @media (max-width: 767px) { flex-direction: column; align-items: flex-start; padding: 32px; gap: 28px; h2 { font-size: 30px; } p { font-size: 12px; } }
            @media (max-width: 575px) { padding: 30px 24px; h2 { font-size: 28px; } }
        `,
        ctaActions: css`
            display: flex; flex-shrink: 0; gap: 12px;
            @media (max-width: 991px) { flex-direction: column; }
            @media (max-width: 767px) { flex-direction: row; }
            @media (max-width: 575px) { width: 100%; button { flex: 1; padding-inline: 15px; height: 44px; font-size: 13px; gap: 10px; } }
        `,
        footer: css`
            display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 18px; padding-block: 30px; color: ${muted}; font-size: 10px; text-align: center;
            > div { display: flex; flex-wrap: wrap; align-items: center; justify-content: center; gap: 12px; &:empty { display: none; } } a { color: inherit; text-decoration: none; &:hover { color: ${accent}; } }
            @media (max-width: 575px) { padding-block: 25px; gap: 12px; font-size: 9px; }
        `,
        backToTop: css`
            position: fixed; z-index: 40; right: 24px; bottom: 24px; display: grid; place-items: center; width: 38px; height: 38px; border: 1px solid ${line}; border-radius: ${token.borderRadius}px; color: ${muted}; background: ${paper}; opacity: 0; visibility: hidden; transition: opacity .2s, color .2s;
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

const products: {icon: "monitor" | "api" | "safe" | "ops"; title: string; description: string; features: string[]}[] = [
    {icon: "monitor", title: "智能监控", description: "实时监控系统运行状态，智能告警，让您随时掌握系统健康状况。", features: ["实时数据监控", "智能告警通知", "可视化报表", "性能分析"]},
    {icon: "api", title: "API 服务", description: "提供丰富的API接口，文档完善，便于集成，助力快速开发。", features: ["RESTful API", "完善文档", "SDK支持", "高并发处理"]},
    {icon: "safe", title: "数据安全", description: "企业级数据安全方案，多重备份机制，确保数据万无一失。", features: ["数据加密", "定期备份", "灾难恢复", "权限管理"]},
    {icon: "ops", title: "自动化运维", description: "自动化部署与运维工具，提升运维效率，降低人工成本。", features: ["自动部署", "日志管理", "版本控制", "回滚机制"]},
];

const productConnections = [
    "M92 54V83Q92 95 104 95H176Q200 95 200 119V137",
    "M308 54V83Q308 95 296 95H224Q200 95 200 119V137",
    "M92 220V193Q92 181 104 181H176Q200 181 200 157V137",
    "M308 220V193Q308 181 296 181H224Q200 181 200 157V137",
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
    const heroVisualRef = useRef<HTMLDivElement>(null);
    const productVisualRef = useRef<HTMLDivElement>(null);
    const [heroVisualRunning, setHeroVisualRunning] = useState(false);
    const [productVisualRunning, setProductVisualRunning] = useState(false);
    const heroVisualId = `home-core-${useId().replace(/:/g, "")}`;
    const productVisualTitleId = `product-visual-${useId().replace(/:/g, "")}`;
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
        const token = getLoginToken();
        if (!token) return;
        let active = true;
        user?.getWebUser().then((data: any) => {
            if (active && getLoginToken() === token) {
                user?.setWeb(data);
            }
        });
        return () => {
            active = false;
        };
    }, []);

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
        }, {threshold: 0, rootMargin: "0px 0px -48px 0px"});
        const onMotionChange = () => {
            if (reducedMotion.matches) {
                page.removeAttribute("data-motion");
                revealObserver.disconnect();
                page.querySelectorAll(".home-reveal").forEach(element => element.setAttribute("data-visible", "true"));
                return;
            }
            page.setAttribute("data-motion", "true");
            page.querySelectorAll('.home-reveal:not([data-visible="true"])').forEach(element => revealObserver.observe(element));
        };
        onMotionChange();
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

    useEffect(() => {
        let visible = true;
        const update = () => setHeroVisualRunning(visible && !document.hidden);
        const observer = new IntersectionObserver(([entry]) => {
            visible = entry.isIntersecting;
            update();
        }, {threshold: 0});
        if (heroVisualRef.current) observer.observe(heroVisualRef.current);
        document.addEventListener("visibilitychange", update);
        update();
        return () => {
            observer.disconnect();
            document.removeEventListener("visibilitychange", update);
        };
    }, []);

    useEffect(() => {
        let inViewport = false;
        const reducedMotion = window.matchMedia("(prefers-reduced-motion: reduce)");
        const update = () => setProductVisualRunning(inViewport && !document.hidden && !reducedMotion.matches);
        const observer = typeof IntersectionObserver === "undefined" ? undefined : new IntersectionObserver(([entry]) => {
            inViewport = entry.isIntersecting;
            update();
        }, {threshold: 0});
        if (productVisualRef.current) observer?.observe(productVisualRef.current);
        document.addEventListener("visibilitychange", update);
        reducedMotion.addEventListener("change", update);
        update();
        return () => {
            observer?.disconnect();
            document.removeEventListener("visibilitychange", update);
            reducedMotion.removeEventListener("change", update);
        };
    }, [product.icon]);

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
                                <h1>让想法成为现实<span>让可能发生</span></h1>
                                <p>{brand}</p>
                                <div className={styles.heroActions}>
                                    <button className={styles.primaryButton} onClick={handleGetStarted}>开始体验<Icon type="ArrowRightOutlined"/></button>
                                    <button className={styles.secondaryButton} onClick={() => scrollToSection("products")}>探索产品<Icon type="ArrowDownOutlined"/></button>
                                </div>
                                <div className={styles.heroNotes}><span><Icon type="CheckOutlined"/>轻盈体验</span><span><Icon type="CheckOutlined"/>安全可靠</span><span><Icon type="CheckOutlined"/>持续进化</span></div>
                            </div>
                            <div className={styles.heroArt} aria-hidden="true">
                                <div ref={heroVisualRef} className={styles.heroVisual} data-running={heroVisualRunning} aria-hidden="true">
                                    <svg viewBox="0 0 640 600" fill="none">
                                        <defs>
                                            <linearGradient id={`${heroVisualId}-glass`} x1="146" y1="160" x2="447" y2="319" gradientUnits="userSpaceOnUse">
                                                <stop stopColor="var(--core-tint)" stopOpacity=".95"/>
                                                <stop offset=".46" stopColor="var(--core-soft)" stopOpacity=".72"/>
                                                <stop offset="1" stopColor="var(--core-border)" stopOpacity=".5"/>
                                            </linearGradient>
                                            <linearGradient id={`${heroVisualId}-base`} x1="167" y1="365" x2="457" y2="534" gradientUnits="userSpaceOnUse">
                                                <stop stopColor="var(--core-tint)" stopOpacity=".92"/>
                                                <stop offset=".5" stopColor="var(--core-soft)" stopOpacity=".95"/>
                                                <stop offset="1" stopColor="var(--core-light)" stopOpacity=".75"/>
                                            </linearGradient>
                                            <linearGradient id={`${heroVisualId}-rim`} x1="111" y1="149" x2="515" y2="232" gradientUnits="userSpaceOnUse">
                                                <stop stopColor="var(--core-border)"/>
                                                <stop offset=".46" stopColor="var(--core-highlight)"/>
                                                <stop offset="1" stopColor="var(--core-primary)"/>
                                            </linearGradient>
                                            <linearGradient id={`${heroVisualId}-chip`} x1="259" y1="251" x2="381" y2="305" gradientUnits="userSpaceOnUse">
                                                <stop stopColor="var(--core-bright)"/>
                                                <stop offset=".5" stopColor="var(--core-primary)"/>
                                                <stop offset="1" stopColor="var(--core-deep)"/>
                                            </linearGradient>
                                            <linearGradient id={`${heroVisualId}-beam`} x1="0" y1="229" x2="0" y2="426" gradientUnits="userSpaceOnUse">
                                                <stop stopColor="var(--core-primary)" stopOpacity="0"/>
                                                <stop offset=".45" stopColor="var(--core-bright)" stopOpacity=".14"/>
                                                <stop offset="1" stopColor="var(--core-primary)" stopOpacity="0"/>
                                            </linearGradient>
                                            <radialGradient id={`${heroVisualId}-shadow`}>
                                                <stop stopColor="var(--core-primary)" stopOpacity=".18"/>
                                                <stop offset="1" stopColor="var(--core-primary)" stopOpacity="0"/>
                                            </radialGradient>
                                            <radialGradient id={`${heroVisualId}-light`}>
                                                <stop stopColor="var(--core-light)" stopOpacity=".55"/>
                                                <stop offset="1" stopColor="var(--core-soft)" stopOpacity="0"/>
                                            </radialGradient>
                                            <clipPath id={`${heroVisualId}-top-clip`}>
                                                <path d="M320 64L523 169L320 274L117 169Z"/>
                                            </clipPath>
                                        </defs>

                                        <ellipse cx="325" cy="517" rx="286" ry="72" fill={`url(#${heroVisualId}-shadow)`}/>
                                        <ellipse cx="320" cy="272" rx="275" ry="255" fill={`url(#${heroVisualId}-light)`}/>
                                        <g stroke="var(--core-border)" strokeWidth=".8" opacity=".45">
                                            <path d="M39 429L320 284L602 429L320 575Z" strokeDasharray="3 6"/>
                                            <path d="M70 413L43 427L70 441M570 413L597 427L570 441M307 291L320 284L333 291"/>
                                            <path d="M93 399L93 302M547 399L547 302" strokeDasharray="2 6"/>
                                        </g>

                                        <g>
                                            <path d="M76 438L320 564L564 438V447L320 573L76 447Z" fill="var(--core-light)" fillOpacity=".7"/>
                                            <path d="M320 312L564 438L320 564L76 438Z" fill={`url(#${heroVisualId}-base)`} stroke="var(--core-border)"/>
                                            <path d="M91 438L320 320L549 438" stroke="var(--core-highlight)" strokeWidth="2" strokeOpacity=".9"/>
                                            <path d="M320 337L516 438L320 539L124 438Z" stroke="var(--core-light)" strokeWidth=".8"/>
                                            <g stroke="var(--core-border)" strokeWidth="1" strokeLinecap="round" strokeLinejoin="round" opacity=".75">
                                                <path d="M320 388L264 417L208 388L159 413M293 402L237 431L201 431L173 446"/>
                                                <path d="M293 427L238 456L238 476L267 491M320 441L281 461L281 483L312 499"/>
                                                <path d="M347 427L402 456L441 436L475 453M347 402L403 431L450 407L479 422"/>
                                                <path d="M320 388L376 417L430 389L450 399M320 441L358 461L358 485L331 499"/>
                                                <path d="M251 384L285 366L320 384L355 366L389 384"/>
                                                <path d="M250 502L290 523L320 508L350 523L391 502"/>
                                            </g>
                                            <g fill="var(--core-tint)" stroke="var(--core-line)" strokeWidth="1.2">
                                                <circle cx="159" cy="413" r="3"/><circle cx="173" cy="446" r="3"/>
                                                <circle cx="267" cy="491" r="3"/><circle cx="312" cy="499" r="3"/>
                                                <circle cx="475" cy="453" r="3"/><circle cx="479" cy="422" r="3"/>
                                                <circle cx="450" cy="399" r="3"/><circle cx="331" cy="499" r="3"/>
                                            </g>
                                            <path className="core-signal" d="M159 413L208 388L264 417L320 388L376 417L430 389L450 399"
                                                  stroke="var(--core-primary)" strokeWidth="2.5" strokeLinecap="round" strokeDasharray="26 654"/>
                                            <path className="core-signal core-signal-secondary" d="M173 446L201 431L237 431L293 402L347 427L402 456L441 436L475 453"
                                                  stroke="var(--core-bright)" strokeWidth="2.2" strokeLinecap="round" strokeDasharray="20 660"/>
                                            <path d="M297 421L320 409L343 421L320 433Z" fill="var(--core-primary)" fillOpacity=".14" stroke="var(--core-border)"/>
                                        </g>

                                        <path d="M233 226L320 272L407 226V400L320 446L233 400Z" fill={`url(#${heroVisualId}-beam)`}/>
                                        <g stroke="var(--core-line)" strokeWidth=".8" strokeDasharray="3 7" opacity=".36">
                                            <path d="M233 229V393M407 229V393M320 275V428"/>
                                        </g>

                                        <g className="core-middle">
                                            <path d="M99 319L320 433L541 319V328L320 442L99 328Z" fill="var(--core-border)" fillOpacity=".38"/>
                                            <path d="M320 205L541 319L320 433L99 319Z" fill={`url(#${heroVisualId}-glass)`} stroke="var(--core-line)" strokeWidth="1.1"/>
                                            <path d="M100 319L320 205L540 319" stroke="var(--core-highlight)" strokeWidth="2"/>
                                            <path d="M122 319L320 217L518 319L320 421Z" stroke="var(--core-line)" strokeWidth=".7" opacity=".6"/>
                                            <path d="M155 319L320 234L485 319L320 404Z" stroke="var(--core-light)"/>
                                            <path d="M153 300L132 311L153 322M487 300L508 311L487 322M301 393L320 403L339 393" stroke="var(--core-line)" strokeWidth="1.5"/>
                                            <g stroke="var(--core-line)" strokeWidth="1.2">
                                                <path d="M240 285L211 300M251 291L222 306M262 297L233 312M273 303L244 318"/>
                                                <path d="M367 303L396 318M378 297L407 312M389 291L418 306M400 285L429 300"/>
                                                <path d="M286 354L269 363M297 360L280 369M343 354L360 363M332 360L349 369"/>
                                            </g>
                                            <path d="M232 294L320 339L408 294V318L320 365L232 318Z" fill="var(--core-line)" fillOpacity=".35"/>
                                            <path d="M232 282L320 327L408 282V302L320 348L232 302Z" fill="var(--core-primary)"/>
                                            <path d="M320 327L408 282V302L320 348Z" fill="var(--core-deep)"/>
                                            <path d="M320 236L408 282L320 327L232 282Z" fill={`url(#${heroVisualId}-chip)`} stroke="var(--core-light)" strokeWidth="1.3"/>
                                            <path d="M320 247L386 282L320 316L254 282Z" stroke="var(--core-light)" strokeWidth=".7" strokeOpacity=".55"/>
                                            <path d="M287 279L320 262L353 279L320 296Z" fill="var(--core-highlight)" fillOpacity=".93"/>
                                            <path d="M298 290L320 278L342 290L320 302Z" fill="var(--core-light)"/>
                                            <path d="M309 272L320 266L331 272L320 278Z" fill="var(--core-primary)"/>
                                            <path className="core-pulse" d="M276 282L320 259L364 282L320 305Z" stroke="var(--core-soft)" strokeWidth="2"/>
                                            <path d="M243 296L311 331M329 331L397 296" stroke="var(--core-border)" strokeWidth="1.4" opacity=".8"/>
                                            <g className="core-node" fill="var(--core-bright)">
                                                <ellipse cx="176" cy="325" rx="3.5" ry="2"/>
                                                <ellipse cx="464" cy="325" rx="3.5" ry="2"/>
                                                <ellipse cx="320" cy="391" rx="3.5" ry="2"/>
                                            </g>
                                        </g>

                                        <g className="core-top">
                                            <path d="M117 169L320 274L523 169V176L320 281L117 176Z" fill="var(--core-line)" fillOpacity=".45"/>
                                            <path d="M320 64L523 169L320 274L117 169Z" fill={`url(#${heroVisualId}-glass)`} stroke={`url(#${heroVisualId}-rim)`} strokeWidth="1.5"/>
                                            <path d="M134 169L320 73L506 169L320 265Z" stroke="var(--core-highlight)" strokeWidth="1.2" strokeOpacity=".92"/>
                                            <path d="M158 169L320 85L482 169L320 253Z" stroke="var(--core-border)" strokeWidth=".8" strokeOpacity=".65"/>
                                            <path d="M178 159L198 149M208 143L228 133M238 127L258 117M268 111L288 101" stroke="var(--core-line)" strokeWidth="2" strokeOpacity=".52"/>
                                            <path d="M252 169L320 134L388 169L320 204Z" fill="var(--core-tint)" fillOpacity=".45" stroke="var(--core-line)" strokeOpacity=".65"/>
                                            <path d="M268 169L320 142L372 169L320 196Z" stroke="var(--core-highlight)" strokeWidth="1.3"/>
                                            <path d="M294 169L320 155L346 169L320 183Z" fill="var(--core-line)" fillOpacity=".22"/>
                                            <path d="M304 169L320 161L336 169L320 177Z" fill="var(--core-highlight)" fillOpacity=".85"/>
                                            <g clipPath={`url(#${heroVisualId}-top-clip)`}>
                                                <g className="core-scan">
                                                    <path d="M62 76L320 210L578 76" stroke="var(--core-highlight)" strokeWidth="12" strokeOpacity=".18"/>
                                                    <path d="M62 76L320 210L578 76" stroke="var(--core-highlight)" strokeWidth="2"/>
                                                </g>
                                            </g>
                                            <path d="M135 169L160 182M295 260L320 273L345 260M480 191L505 178" stroke="var(--core-primary)" strokeWidth="2" strokeLinecap="round"/>
                                            <g fill="var(--core-tint)" stroke="var(--core-line)" strokeWidth="1">
                                                <path d="M207 191L214 187L221 191L214 195Z"/>
                                                <path d="M228 202L235 198L242 202L235 206Z"/>
                                                <path d="M249 213L256 209L263 213L256 217Z"/>
                                            </g>
                                        </g>
                                    </svg>
                                </div>
                                <div className={styles.artCaption}><span>IDEAS. CONNECTED.</span><span>构建无限可能</span></div>
                            </div>
                        </div>
                        <div className={cx(styles.container, styles.heroBottom)}>
                            <span>TECHNOLOGY, MADE HUMAN.</span>
                            <div><span><b>01</b>高效连接</span><span><b>02</b>安全守护</span><span><b>03</b>自由扩展</span></div>
                            <button onClick={() => scrollToSection("features")} aria-label="探索平台优势"><Icon type="ArrowDownOutlined"/></button>
                        </div>
                    </section>

                    <section id="features" className={styles.featuresSection}>
                        <div className={cx(styles.container, styles.section)}>
                            <div className={cx(styles.sectionHeading, "home-reveal")}>
                                <div><span className={styles.eyebrow}>01 / BUILT TO GO FURTHER</span><h2>强大的内核。<span>轻盈的体验。</span></h2></div>
                                <p>每一个顺畅的瞬间，背后都有可靠的技术。<br/>从快速开始，到从容进阶。</p>
                            </div>
                            <div className={styles.featureGrid}>
                                {features.map((feature, index) => <article key={feature.title} className="home-reveal" data-step={index < 2 ? index : index - 2} data-spotlight="card">
                                    <div className={styles.spotlight} aria-hidden="true"/>
                                    {index < 2 && <div className={styles.featureArt} aria-hidden="true"><HomeIcon type={feature.icon}/></div>}
                                    <div className={styles.featureTop}><HomeIcon type={feature.icon}/><span>0{index + 1}</span></div>
                                    <h3>{feature.title}</h3><p>{feature.description}</p>
                                </article>)}
                            </div>
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
                                <div key={product.icon} id="home-product-detail" className={styles.productDetail}>
                                    <div className={styles.productCopy}><span className={styles.eyebrow}>0{productIndex + 1} / CONNECTED SERVICES</span><h3>{product.title}</h3><p>{product.description}</p><ul>{product.features.map(feature => <li key={feature}><Icon type="CheckOutlined"/>{feature}</li>)}</ul><button className={styles.textButton} onClick={handleGetStarted}>开始使用<Icon type="ArrowRightOutlined"/></button></div>
                                    <div ref={productVisualRef} className={styles.productVisual} data-running={productVisualRunning} role="group" aria-labelledby={productVisualTitleId}>
                                        <div className={styles.productNetwork}>
                                            <svg className={styles.productConnections} viewBox="0 0 400 274" preserveAspectRatio="none" fill="none" aria-hidden="true">
                                                {productConnections.map((path, index) => <g key={path}>
                                                    <path d={path} stroke="var(--product-line)" strokeWidth="1" vectorEffect="non-scaling-stroke"/>
                                                    <path className="product-signal" d={path} pathLength="100" vectorEffect="non-scaling-stroke" style={{animationDelay: `${index * -1.35}s`}}/>
                                                </g>)}
                                                <g stroke="var(--product-line)" strokeWidth="1" vectorEffect="non-scaling-stroke">
                                                    <path d="M8 117v-8h8M384 109h8v8M8 157v8h8M384 165h8v-8"/>
                                                    <path d="M32 137h22M346 137h22" strokeDasharray="2 4"/>
                                                </g>
                                            </svg>
                                            {product.features.slice(0, 4).map((item, index) => <div className={styles.productNode} data-position={index} key={`${index}-${item}`} style={{gridColumn: index % 2 + 1, gridRow: index < 2 ? 1 : 3}}>
                                                <svg viewBox="0 0 20 20" fill="none" stroke="currentColor" strokeWidth="1.1" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
                                                    {index === 0 && <><rect x="3" y="3" width="14" height="14" rx="3"/><path d="M6 11h2l2-4 2 6 2-3"/></>}
                                                    {index === 1 && <><path d="m10 2 7 4v8l-7 4-7-4V6Z"/><path d="m3 6 7 4 7-4M10 10v8"/></>}
                                                    {index === 2 && <><rect x="3" y="3" width="5" height="5" rx="1"/><rect x="12" y="3" width="5" height="5" rx="1"/><rect x="3" y="12" width="5" height="5" rx="1"/><rect x="12" y="12" width="5" height="5" rx="1"/></>}
                                                    {index === 3 && <><circle cx="10" cy="10" r="6"/><circle cx="10" cy="10" r="2"/><path d="M10 1v3M10 16v3M1 10h3M16 10h3"/></>}
                                                </svg>
                                                <span>{item}</span>
                                            </div>)}
                                            <div className={styles.productCenter}>
                                                <div className={styles.productOrbit} aria-hidden="true"><i/></div>
                                                <div className={styles.productCore} aria-hidden="true">
                                                    <svg viewBox="0 0 40 40" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
                                                        {product.icon === "monitor" && <><rect x="5" y="7" width="30" height="23" rx="4"/><path d="M14 35h12M20 30v5M10 21h5l3-7 5 12 3-5h4"/></>}
                                                        {product.icon === "api" && <><rect x="4" y="7" width="32" height="26" rx="5"/><path d="m14 15-5 5 5 5m12-10 5 5-5 5m-4-13-4 16"/></>}
                                                        {product.icon === "safe" && <><path d="M20 3 33 8v11c0 8-5 14-13 18C12 33 7 27 7 19V8Z"/><path d="m13 20 5 5 9-10"/></>}
                                                        {product.icon === "ops" && <><path d="M11 9a14 14 0 0 1 22 9M29 31a14 14 0 0 1-22-9M33 10v8h-8M7 30v-8h8"/><rect x="15" y="15" width="10" height="10" rx="3"/><path d="M20 11v4M20 25v4M11 20h4M25 20h4"/></>}
                                                    </svg>
                                                </div>
                                                <span id={productVisualTitleId}>{product.title}</span>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </section>

                    <section className={cx(styles.container, styles.section)}>
                        <div className={styles.experienceLayout}>
                            <div className={cx(styles.manifesto, "home-reveal")}>
                                <span className={styles.eyebrow}>LESS FRICTION. MORE FOCUS.</span>
                                <h2>好的工具，<span>让你忘记工具。</span></h2>
                                <p>直观的体验、可靠的服务，与持续生长的能力。<br/>陪你把日常的小想法，变成下一步的大可能。</p>
                            </div>
                            <div className={styles.advantages}>
                                {advantages.map((item, index) => <article key={item.title} className="home-reveal" data-step={index}><span>0{index + 1}</span><h3>{item.title}</h3><p>{item.describe}</p></article>)}
                            </div>
                        </div>
                    </section>

                    <section id="testimonials" className={styles.testimonials}>
                        <div className={cx(styles.container, styles.section)}>
                            <div className={cx(styles.sectionHeading, "home-reveal")}>
                                <div><span className={styles.eyebrow}>IN GOOD COMPANY</span><h2>体验，由使用者诉说。</h2></div>
                                <p>好的体验，体现在每一次使用中。</p>
                            </div>
                            <div className={styles.testimonialGrid}>
                                {testimonials.map((item, index) => <article key={item.name} className="home-reveal" data-step={index} data-spotlight="card">
                                    <div className={styles.spotlight} aria-hidden="true"/>
                                    <span className={styles.quoteMark} aria-hidden="true">“</span>
                                    <blockquote>{item.content}</blockquote>
                                    <div className={styles.author}><span>{item.name.slice(0, 1)}</span><div><strong>{item.name}</strong><small>{item.role}</small></div></div>
                                </article>)}
                            </div>
                        </div>
                    </section>

                    <div className={styles.container}>
                        <section className={cx(styles.cta, "home-reveal")} data-spotlight="hero">
                            <div className={styles.spotlight} aria-hidden="true"/>
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
