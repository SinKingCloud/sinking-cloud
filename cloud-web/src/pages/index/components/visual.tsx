import React, {useEffect, useId, useRef, useState} from "react";
import {createStyles} from "antd-style";

const useStyles = createStyles(({css, token}) => ({
    visual: css`
        --core-tint: color-mix(in srgb, ${token.colorPrimary} 5%, white);
        --core-soft: color-mix(in srgb, ${token.colorPrimary} 14%, white);
        --core-light: color-mix(in srgb, ${token.colorPrimary} 28%, white);
        --core-border: color-mix(in srgb, ${token.colorPrimary} 42%, white);
        --core-line: color-mix(in srgb, ${token.colorPrimary} 58%, white);
        --core-primary: ${token.colorPrimary};
        --core-bright: ${token.colorPrimaryHover};
        --core-deep: ${token.colorPrimaryActive};
        position: relative;
        width: 100%;
        aspect-ratio: 16 / 15;
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
    `,
}));

export default function HomeVisual({className}: {className?: string}) {
    const {styles, cx} = useStyles();
    const containerRef = useRef<HTMLDivElement>(null);
    const [running, setRunning] = useState(false);
    const id = `home-core-${useId().replace(/:/g, "")}`;

    useEffect(() => {
        let visible = true;
        const update = () => setRunning(visible && !document.hidden);
        const observer = new IntersectionObserver(([entry]) => {
            visible = entry.isIntersecting;
            update();
        }, {threshold: 0});
        if (containerRef.current) observer.observe(containerRef.current);
        document.addEventListener("visibilitychange", update);
        update();
        return () => {
            observer.disconnect();
            document.removeEventListener("visibilitychange", update);
        };
    }, []);

    return <div ref={containerRef} className={cx(styles.visual, className)} data-running={running} aria-hidden="true">
        <svg viewBox="0 0 640 600" fill="none">
            <defs>
                <linearGradient id={`${id}-glass`} x1="146" y1="160" x2="447" y2="319" gradientUnits="userSpaceOnUse">
                    <stop stopColor="var(--core-tint)" stopOpacity=".95"/>
                    <stop offset=".46" stopColor="var(--core-soft)" stopOpacity=".72"/>
                    <stop offset="1" stopColor="var(--core-border)" stopOpacity=".5"/>
                </linearGradient>
                <linearGradient id={`${id}-base`} x1="167" y1="365" x2="457" y2="534" gradientUnits="userSpaceOnUse">
                    <stop stopColor="var(--core-tint)" stopOpacity=".92"/>
                    <stop offset=".5" stopColor="var(--core-soft)" stopOpacity=".95"/>
                    <stop offset="1" stopColor="var(--core-light)" stopOpacity=".75"/>
                </linearGradient>
                <linearGradient id={`${id}-rim`} x1="111" y1="149" x2="515" y2="232" gradientUnits="userSpaceOnUse">
                    <stop stopColor="var(--core-border)"/>
                    <stop offset=".46" stopColor="#fff"/>
                    <stop offset="1" stopColor="var(--core-primary)"/>
                </linearGradient>
                <linearGradient id={`${id}-chip`} x1="259" y1="251" x2="381" y2="305" gradientUnits="userSpaceOnUse">
                    <stop stopColor="var(--core-bright)"/>
                    <stop offset=".5" stopColor="var(--core-primary)"/>
                    <stop offset="1" stopColor="var(--core-deep)"/>
                </linearGradient>
                <linearGradient id={`${id}-beam`} x1="0" y1="229" x2="0" y2="426" gradientUnits="userSpaceOnUse">
                    <stop stopColor="var(--core-primary)" stopOpacity="0"/>
                    <stop offset=".45" stopColor="var(--core-bright)" stopOpacity=".14"/>
                    <stop offset="1" stopColor="var(--core-primary)" stopOpacity="0"/>
                </linearGradient>
                <radialGradient id={`${id}-shadow`}>
                    <stop stopColor="var(--core-primary)" stopOpacity=".18"/>
                    <stop offset="1" stopColor="var(--core-primary)" stopOpacity="0"/>
                </radialGradient>
                <radialGradient id={`${id}-light`}>
                    <stop stopColor="var(--core-light)" stopOpacity=".55"/>
                    <stop offset="1" stopColor="var(--core-soft)" stopOpacity="0"/>
                </radialGradient>
                <clipPath id={`${id}-top-clip`}>
                    <path d="M320 64L523 169L320 274L117 169Z"/>
                </clipPath>
            </defs>

            <ellipse cx="325" cy="517" rx="286" ry="72" fill={`url(#${id}-shadow)`}/>
            <ellipse cx="320" cy="272" rx="275" ry="255" fill={`url(#${id}-light)`}/>
            <g stroke="var(--core-border)" strokeWidth=".8" opacity=".45">
                <path d="M39 429L320 284L602 429L320 575Z" strokeDasharray="3 6"/>
                <path d="M70 413L43 427L70 441M570 413L597 427L570 441M307 291L320 284L333 291"/>
                <path d="M93 399L93 302M547 399L547 302" strokeDasharray="2 6"/>
            </g>

            <g>
                <path d="M76 438L320 564L564 438V447L320 573L76 447Z" fill="var(--core-light)" fillOpacity=".7"/>
                <path d="M320 312L564 438L320 564L76 438Z" fill={`url(#${id}-base)`} stroke="var(--core-border)"/>
                <path d="M91 438L320 320L549 438" stroke="#fff" strokeWidth="2" strokeOpacity=".9"/>
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

            <path d="M233 226L320 272L407 226V400L320 446L233 400Z" fill={`url(#${id}-beam)`}/>
            <g stroke="var(--core-line)" strokeWidth=".8" strokeDasharray="3 7" opacity=".36">
                <path d="M233 229V393M407 229V393M320 275V428"/>
            </g>

            <g className="core-middle">
                <path d="M99 319L320 433L541 319V328L320 442L99 328Z" fill="var(--core-border)" fillOpacity=".38"/>
                <path d="M320 205L541 319L320 433L99 319Z" fill={`url(#${id}-glass)`} stroke="var(--core-line)" strokeWidth="1.1"/>
                <path d="M100 319L320 205L540 319" stroke="#fff" strokeWidth="2"/>
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
                <path d="M320 236L408 282L320 327L232 282Z" fill={`url(#${id}-chip)`} stroke="var(--core-light)" strokeWidth="1.3"/>
                <path d="M320 247L386 282L320 316L254 282Z" stroke="var(--core-light)" strokeWidth=".7" strokeOpacity=".55"/>
                <path d="M287 279L320 262L353 279L320 296Z" fill="var(--core-tint)" fillOpacity=".93"/>
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
                <path d="M320 64L523 169L320 274L117 169Z" fill={`url(#${id}-glass)`} stroke={`url(#${id}-rim)`} strokeWidth="1.5"/>
                <path d="M134 169L320 73L506 169L320 265Z" stroke="#fff" strokeWidth="1.2" strokeOpacity=".92"/>
                <path d="M158 169L320 85L482 169L320 253Z" stroke="var(--core-border)" strokeWidth=".8" strokeOpacity=".65"/>
                <path d="M178 159L198 149M208 143L228 133M238 127L258 117M268 111L288 101" stroke="var(--core-line)" strokeWidth="2" strokeOpacity=".52"/>
                <path d="M252 169L320 134L388 169L320 204Z" fill="var(--core-tint)" fillOpacity=".45" stroke="var(--core-line)" strokeOpacity=".65"/>
                <path d="M268 169L320 142L372 169L320 196Z" stroke="#fff" strokeWidth="1.3"/>
                <path d="M294 169L320 155L346 169L320 183Z" fill="var(--core-line)" fillOpacity=".22"/>
                <path d="M304 169L320 161L336 169L320 177Z" fill="#fff" fillOpacity=".85"/>
                <g clipPath={`url(#${id}-top-clip)`}>
                    <g className="core-scan">
                        <path d="M62 76L320 210L578 76" stroke="#fff" strokeWidth="12" strokeOpacity=".18"/>
                        <path d="M62 76L320 210L578 76" stroke="#fff" strokeWidth="2"/>
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
    </div>;
}
