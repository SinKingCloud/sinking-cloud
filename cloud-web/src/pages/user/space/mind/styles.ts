import {createStyles} from "antd-style";

const useStyles = createStyles(({css, token, isDarkMode}: any) => ({
    page: css`
        --mind-page-height: calc(100dvh - 177px);
        --mind-library-min-height: calc(var(--mind-page-height) - 63px);
        min-height: var(--mind-library-min-height);

        @supports not (height: 100dvh) {
            --mind-page-height: calc(100vh - 177px);
        }
    `,
    library: css`
        container-name: mind-library;
        container-type: inline-size;
        min-width: 0;
        min-height: inherit;
    `,
    libraryCard: css`
        min-height: var(--mind-library-min-height);
        overflow: hidden;
        border-radius: ${token.borderRadiusLG}px;
        background: ${token.colorBgContainer};

        > .ant-card-body {
            display: grid;
            grid-template-columns: 238px minmax(0, 1fr);
            grid-template-rows: auto auto minmax(0, 1fr);
            min-width: 0;
            min-height: var(--mind-library-min-height);
            padding: 0 !important;
        }

        @container mind-library (max-width: 780px) {
            > .ant-card-body {
                grid-template-columns: 1fr;
                grid-template-rows: auto auto auto minmax(0, 1fr);
            }
        }
    `,
    categoryPanel: css`
        position: relative;
        display: block;
        grid-column: 1;
        grid-row: 2 / span 2;
        width: auto;
        min-width: 0;
        min-height: 0;
        max-height: none;
        padding: 0;
        border-radius: 0 0 0 ${token.borderRadiusLG}px;
        background: ${token.colorBgContainer};
        align-self: stretch;
        overflow: hidden;

        @container mind-library (max-width: 780px) {
            grid-column: 1;
            grid-row: 2;
            width: 100%;
            max-height: none;
            border-radius: 0;
            align-self: stretch;
        }

        @media (max-width: 640px) {
            width: 100%;
            height: auto;
            max-height: none;
        }
    `,
    categoryInner: css`
        position: absolute;
        inset: 0;
        display: flex;
        flex-direction: column;
        min-width: 0;
        min-height: 0;
        padding: 10px 4px 10px 10px;

        @container mind-library (max-width: 780px) {
            position: static;
        }

        @media (max-width: 640px) {
            padding: 9px;
        }
    `,
    categoryLoading: css`
        display: flex;
        flex: 1;
        align-items: center;
        justify-content: center;
        min-width: 0;
        min-height: 0;

        @container mind-library (max-width: 780px) {
            min-height: 52px;
        }
    `,
    categoryRail: css`
        display: flex;
        flex: 1;
        flex-direction: column;
        min-width: 0;
        min-height: 0;

        @container mind-library (max-width: 780px) {
            flex-direction: row;
            align-items: center;
            gap: 6px;
            width: 100%;
        }
    `,
    categoryScrollButton: css`
        display: none;
        flex: 0 0 auto;
        align-items: center;
        justify-content: center;
        width: 28px;
        height: 34px;
        padding: 0;
        border: 0;
        border-radius: ${token.borderRadius}px;
        outline: none;
        background: ${token.colorFillQuaternary};
        color: ${token.colorTextSecondary};
        cursor: pointer;
        transition: background 0.16s ease, color 0.16s ease, opacity 0.16s ease;

        &:hover:not(:disabled),
        &:focus-visible:not(:disabled) {
            background: ${token.colorFillTertiary};
            color: ${token.colorPrimary};
        }

        &:focus-visible {
            box-shadow: 0 0 0 2px ${token.colorPrimaryBg};
        }

        &:disabled {
            opacity: 0.3;
            cursor: default;
        }

        .anticon {
            font-size: 11px;
        }

        @container mind-library (max-width: 780px) {
            display: inline-flex;
        }

        @media (max-width: 640px), (hover: none) and (pointer: coarse) {
            display: none;
        }
    `,
    categoryFooter: css`
        flex: none;
        padding-top: 8px;

        @container mind-library (max-width: 780px) {
            display: none;
        }
    `,
    categoryAddCard: css`
        display: flex;
        align-items: center;
        gap: 8px;
        width: 100%;
        min-width: 0;
        min-height: 38px;
        padding: 0 10px;
        border: 1px dashed ${token.colorBorder};
        border-radius: ${token.borderRadius}px;
        outline: none;
        appearance: none;
        background: transparent;
        color: ${token.colorTextTertiary};
        cursor: pointer;
        font: inherit;
        font-size: 12px;
        line-height: 20px;
        text-align: left;
        transition: border-color 0.16s ease, background 0.16s ease, color 0.16s ease;

        &.mobile {
            display: none;
        }

        &:hover,
        &:focus-visible {
            border-color: ${token.colorPrimaryBorder};
            background: ${token.colorFillQuaternary};
            color: ${token.colorPrimary};
        }

        &:focus-visible {
            box-shadow: 0 0 0 2px ${token.colorPrimaryBg};
        }

        .icon {
            display: inline-flex;
            align-items: center;
            justify-content: center;
            width: 16px;
            height: 16px;
            border-radius: ${token.borderRadiusSM}px;
            background: ${token.colorFillTertiary};
            font-size: 10px;
        }

        @container mind-library (max-width: 780px) {
            min-height: 34px;
            padding: 0 9px;

            &.mobile {
                display: flex;
                flex: 0 0 auto;
                width: auto;
                white-space: nowrap;
            }
        }
    `,
    categoryScroll: css`
        flex: 1;
        min-height: 0;
        overflow-y: scroll;
        padding: 0 2px 0 0;
        scrollbar-color: ${token.colorBorderSecondary} transparent;
        scrollbar-width: thin;

        &::-webkit-scrollbar {
            width: 6px;
        }

        &::-webkit-scrollbar-track {
            background: transparent;
        }

        &::-webkit-scrollbar-thumb {
            border-radius: ${token.borderRadiusSM}px;
            background: ${token.colorBorderSecondary};
        }

        &::-webkit-scrollbar-thumb:hover {
            background: ${token.colorTextQuaternary};
        }

        @container mind-library (max-width: 780px) {
            overflow-x: auto;
            overflow-y: hidden;
            padding: 0;
            scrollbar-width: none;

            &::-webkit-scrollbar {
                display: none;
            }
        }
    `,
    categoryGrid: css`
        display: flex;
        flex-direction: column;
        gap: 4px;

        @container mind-library (max-width: 780px) {
            flex-direction: row;
            gap: 6px;
            min-width: max-content;
        }
    `,
    categoryCard: css`
        position: relative;
        display: flex;
        align-items: center;
        min-width: 0;
        min-height: 40px;
        padding: 0 7px 0 11px;
        overflow: hidden;
        border: 0;
        border-radius: ${token.borderRadius}px;
        background: transparent;
        cursor: pointer;
        transition: border-color 0.16s ease, background 0.16s ease, color 0.16s ease;

        &:hover,
        &.active {
            background: ${token.colorFillQuaternary};
        }

        &.active {
            color: ${token.colorTextHeading};
            background: ${token.colorPrimaryBg};
        }

        .row {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 8px;
            width: 100%;
            min-width: 0;
        }

        .right {
            display: flex;
            align-items: center;
            gap: 3px;
            flex: none;
        }

        .name-wrap {
            display: inline-flex;
            align-items: center;
            gap: 8px;
            min-width: 0;
        }

        .color {
            flex: 0 0 auto;
            width: 7px;
            height: 7px;
            border-radius: 50%;
            background: color-mix(in srgb, var(--accent), ${token.colorTextTertiary} 8%);
            opacity: 0.62;
        }

        .action {
            display: inline-flex;
            margin-right: -3px;
            opacity: 0;
            transition: opacity 0.16s ease;
        }

        &:hover .action,
        &.active .action {
            opacity: 1;
        }

        .name {
            min-width: 0;
            overflow: hidden;
            color: ${token.colorTextSecondary};
            font-size: 13px;
            font-weight: 400;
            line-height: 20px;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        &.active .name {
            color: ${token.colorTextHeading};
            font-weight: 600;
        }

        &.active .color {
            background: var(--accent);
            box-shadow: 0 0 0 3px color-mix(in srgb, var(--accent), transparent 78%);
            opacity: 1;
        }

        .action .ant-btn {
            color: ${token.colorTextTertiary};
        }

        @container mind-library (max-width: 780px) {
            flex: 0 0 auto;
            min-height: 34px;
            padding: 0 8px 0 9px;
            background: ${token.colorFillQuaternary};

            .name-wrap {
                gap: 6px;
            }

            .color {
                width: 7px;
                height: 7px;
            }

            .action {
                opacity: 1;
            }
        }
    `,
    actionDropdown: css`
        && {
            min-width: 96px;
            padding: 0;
            border-radius: ${token.borderRadius}px;
        }

        && .ant-dropdown-menu {
            min-width: 96px;
            padding: 3px;
            border-radius: ${token.borderRadius}px;
            background: ${token.colorBgElevated};
        }

        && .ant-dropdown-menu-item {
            box-sizing: border-box;
            min-height: 30px;
            margin: 0;
            padding: 4px 9px;
            border-radius: ${token.borderRadiusSM}px;
            color: ${token.colorTextSecondary};
            font-size: 12px;
            font-weight: 400;
            line-height: 20px;
        }

        && .ant-dropdown-menu-title-content {
            font-size: 12px;
            line-height: 20px;
        }

        && .ant-dropdown-menu-item-icon {
            margin-inline-end: 6px;
            color: ${token.colorTextTertiary};
            font-size: 11px;
        }

        && .ant-dropdown-menu-item:hover {
            background: ${token.colorFillQuaternary};
            color: ${token.colorText};
        }

        && .ant-dropdown-menu-item-danger {
            color: ${token.colorError};
        }

        && .ant-dropdown-menu-item-danger .ant-dropdown-menu-item-icon {
            color: ${token.colorError};
        }

        && .ant-dropdown-menu-item-danger:hover,
        && .ant-dropdown-menu-item-danger:hover .ant-dropdown-menu-item-icon {
            color: ${token.colorTextLightSolid};
        }

        && .ant-dropdown-menu-item-divider {
            margin: 3px 6px;
            background: ${token.colorSplit};
        }
    `,
    content: css`
        display: contents;
        min-width: 0;

        > .ant-row {
            display: flex;
            flex-direction: column;
            flex-wrap: nowrap;
            grid-column: 2;
            grid-row: 3;
            height: 100%;
            row-gap: 0 !important;
            column-gap: 0 !important;
            margin: 0 !important;
            min-width: 0;
        }

        > .ant-row > .ant-col {
            flex: none;
            width: 100%;
            max-width: 100%;
            padding: 0 !important;
        }

        @container mind-library (max-width: 780px) {
            display: contents;

            > .ant-row {
                grid-column: 1;
                grid-row: 4;
            }
        }
    `,
    contentPanel: css`
        display: contents;

        @container mind-library (max-width: 780px) {
            display: contents;
        }
    `,
    documentDataCard: css`
        display: block;
        flex: 1;
        width: 100%;
        background: ${token.colorBgContainer};
        padding: 0 6px 24px;

        &.empty {
            display: flex;
            align-items: center;
            justify-content: center;
            min-height: 240px;
            padding: 24px;
        }

        @media (max-width: 640px) {
            padding: 0 11px 12px;
        }
    `,
    documentDataColumn: css`
        display: flex;
        flex: 1 1 auto !important;
        min-height: 0;
    `,
    documentPaginationCard: css`
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 12px 15px;
        background: transparent;

        .ant-pagination {
            width: 100%;
            justify-content: center;
        }

        .ant-pagination-total-text {
            color: ${token.colorTextTertiary};
            font-size: 12px;
        }

        .ant-pagination-options {
            margin-inline-start: 10px;
        }

        .ant-pagination-item,
        .ant-pagination-prev .ant-pagination-item-link,
        .ant-pagination-next .ant-pagination-item-link {
            min-width: 24px;
            height: 24px;
            margin-inline-end: 4px;
            border: 0;
            border-radius: ${token.borderRadiusSM}px;
            background: transparent;
            color: ${token.colorTextSecondary};
            font-size: 11px;
            line-height: 24px;
        }

        .ant-pagination-item a {
            color: inherit;
        }

        .ant-pagination-item:hover,
        .ant-pagination-prev:hover .ant-pagination-item-link,
        .ant-pagination-next:hover .ant-pagination-item-link {
            background: ${token.colorFillTertiary};
            color: ${token.colorTextSecondary};
        }

        .ant-pagination-item-active,
        .ant-pagination-item-active:hover {
            background: ${token.colorBgContainer};
            box-shadow: inset 0 0 0 1px ${token.colorPrimaryBorder};
            color: ${token.colorPrimary};
        }

        .ant-pagination-item-active a,
        .ant-pagination-item-active:hover a {
            color: ${token.colorPrimary};
        }

        .ant-pagination-options-size-changer.ant-select .ant-select-selector {
            height: 24px !important;
            padding-inline: 7px 22px !important;
            border: 0 !important;
            border-radius: ${token.borderRadiusSM}px !important;
            background: ${token.colorBgContainer} !important;
            box-shadow: inset 0 0 0 1px ${token.colorBorderSecondary} !important;
        }

        .ant-pagination-options-size-changer.ant-select,
        .ant-pagination-options-size-changer .ant-select-content,
        .ant-pagination-options-size-changer .ant-select-content-value,
        .ant-pagination-options-size-changer .ant-select-selection-item,
        .ant-pagination-options-quick-jumper {
            color: ${token.colorTextTertiary};
            font-size: 11px;
            line-height: 24px;
        }

        .ant-pagination-options-quick-jumper input {
            width: 32px;
            height: 24px;
            margin-inline: 4px;
            border-radius: ${token.borderRadiusSM}px;
            color: ${token.colorTextSecondary};
            font-size: 11px;
        }

        @container mind-library (max-width: 780px) {
            padding: 12px;
        }
    `,
    documentPaginationColumn: css`
        margin-top: auto;
    `,
    pageSizeSelect: css`
        && {
            height: 24px;
            color: ${token.colorTextTertiary} !important;
            font-size: 11px;
        }

        && .ant-select-selector {
            height: 24px !important;
            padding-inline: 7px 22px !important;
            border: 0 !important;
            border-radius: ${token.borderRadiusSM}px !important;
            background: ${token.colorBgContainer} !important;
            box-shadow: inset 0 0 0 1px ${token.colorBorderSecondary} !important;
        }

        && .ant-select-content,
        && .ant-select-content-value,
        && .ant-select-selection-item {
            color: ${token.colorTextTertiary} !important;
            font-size: 11px !important;
            font-weight: 400 !important;
            line-height: 24px !important;
        }

        && .ant-select-arrow {
            color: ${token.colorTextQuaternary};
            font-size: 9px;
        }
    `,
    pageSizeDropdown: css`
        && {
            padding: 2px;
            border-radius: ${token.borderRadius}px;
        }

        && .ant-select-item {
            min-height: 26px;
            padding: 4px 8px;
            border-radius: ${token.borderRadiusSM}px;
            color: ${token.colorTextSecondary};
            font-size: 11px;
            line-height: 18px;
        }

        && .ant-select-item-option-content {
            font-size: 11px !important;
            font-weight: 400;
        }

        && .ant-select-item-option-selected:not(.ant-select-item-option-disabled) {
            background: ${token.colorFillTertiary};
            color: ${token.colorTextSecondary};
            font-weight: 400;
        }
    `,
    searchCard: css`
        grid-column: 2;
        grid-row: 2;
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 16px;
        min-width: 0;
        padding: 16px 12px 12px;
        background: ${token.colorBgContainer};

        @media (max-width: 640px) {
            gap: 8px;
            padding: 12px 12px 10px;
        }

        @container mind-library (max-width: 780px) {
            grid-column: 1;
            grid-row: 3;
        }
    `,
    searchBox: css`
        flex: 0 1 300px;
        width: 300px;
        max-width: 100%;
        min-width: 0;

        &.ant-input-affix-wrapper {
            height: 36px;
            padding: 0 11px;
            border: 1px solid transparent;
            border-radius: ${token.borderRadiusSM}px;
            background: color-mix(in srgb, ${token.colorFillTertiary} 45%, ${token.colorFillQuaternary});
            box-shadow: none;
            transition: border-color 0.16s ease, background-color 0.16s ease, box-shadow 0.16s ease;
        }

        &.ant-input-affix-wrapper:hover {
            background: color-mix(in srgb, ${token.colorFillTertiary} 65%, ${token.colorFillQuaternary});
        }

        &.ant-input-affix-wrapper-focused {
            border-color: ${token.colorPrimaryBorder};
            background: ${token.colorBgContainer};
            box-shadow: 0 0 0 3px ${token.colorPrimaryBg};
        }

        .ant-input-prefix {
            margin-right: 8px;
            color: ${token.colorTextTertiary};
            font-size: 14px;
        }

        .ant-input {
            background: transparent;
            color: color-mix(in srgb, ${token.colorTextSecondary}, ${token.colorBgContainer} 24%);
            font-size: 11px;
        }

        .ant-input::placeholder {
            color: ${token.colorTextQuaternary};
            font-size: 11px;
        }

        @media (max-width: 820px) {
            flex-basis: 260px;
            width: 260px;
        }

        @media (max-width: 640px) {
            flex: 1 1 auto;
            width: auto;
        }
    `,
    sortTrigger: css`
        display: inline-flex;
        flex: none;
        align-items: center;
        justify-content: center;
        gap: 5px;
        height: 32px;
        padding: 0 8px;
        border: 1px solid transparent;
        border-radius: ${token.borderRadiusSM}px;
        outline: none;
        appearance: none;
        background: ${token.colorFillQuaternary};
        color: ${token.colorTextSecondary};
        cursor: pointer;
        font: inherit;
        line-height: 1;
        transition: background-color 0.16s ease, color 0.16s ease;

        .marker {
            color: ${token.colorTextQuaternary};
            font-size: 11px;
        }

        .value {
            color: ${token.colorTextTertiary};
            font-size: 11px;
            font-weight: 400;
            white-space: nowrap;
            transition: color 0.16s ease;
        }

        .arrow {
            color: ${token.colorTextQuaternary};
            font-size: 9px;
            transition: color 0.16s ease, transform 0.16s ease;
        }

        &:hover,
        &.ant-dropdown-open {
            background: ${token.colorFillSecondary};
        }

        &:hover .value,
        &.ant-dropdown-open .value {
            color: ${token.colorTextSecondary};
        }

        &.ant-dropdown-open .arrow {
            color: ${token.colorTextTertiary};
            transform: rotate(180deg);
        }

        &:focus-visible {
            border-color: ${token.colorPrimaryBorder};
            background: ${token.colorBgContainer};
            box-shadow: 0 0 0 3px ${token.colorPrimaryBg};
        }

        @media (max-width: 640px) {
            flex: 0 0 auto;
        }
    `,
    sortDropdown: css`
        && {
            min-width: 100px;
            padding: 0;
            border-radius: ${token.borderRadius}px;
        }

        && .ant-dropdown-menu {
            padding: 2px !important;
            border-radius: ${token.borderRadius}px !important;
            background: ${token.colorBgElevated};
            box-shadow: ${token.boxShadowSecondary};
        }

        && .ant-dropdown-menu .ant-dropdown-menu-item {
            min-height: 26px !important;
            margin: 0 !important;
            padding: 0 8px !important;
            border-radius: ${token.borderRadiusSM}px !important;
            color: ${token.colorTextTertiary};
            font-size: 11px !important;
            font-weight: 400;
            line-height: 26px !important;
        }

        && .ant-dropdown-menu .ant-dropdown-menu-item .ant-dropdown-menu-title-content {
            font-size: 11px !important;
            line-height: 26px !important;
        }

        && .ant-dropdown-menu .ant-dropdown-menu-item:hover {
            background: ${token.colorFillQuaternary};
            color: ${token.colorTextSecondary};
        }

        && .ant-dropdown-menu .ant-dropdown-menu-item-selected,
        && .ant-dropdown-menu .ant-dropdown-menu-item-selected:hover {
            background: color-mix(in srgb, ${token.colorPrimary}, ${token.colorBgElevated} 94%);
            color: ${token.colorTextSecondary};
        }

        && .ant-dropdown-menu .ant-dropdown-menu-item .ant-dropdown-menu-item-icon {
            margin-inline-end: 4px !important;
            color: ${token.colorTextQuaternary};
            font-size: 9px !important;
        }
    `,
    documentGrid: css`
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(208px, 1fr));
        gap: 14px;
        padding: 5px 6px 8px;

        @media (max-width: 640px) {
            grid-template-columns: 1fr;
            gap: 10px;
            padding: 3px 1px 6px;
        }
    `,
    createHero: css`
        grid-column: 1 / -1;
        grid-row: 1;
        position: relative;
        display: flex;
        align-items: center;
        justify-content: space-between;
        min-height: 128px;
        margin-bottom: 0;
        overflow: hidden;
        border: 0;
        border-radius: ${token.borderRadiusLG}px ${token.borderRadiusLG}px 0 0;
        background: ${isDarkMode
            ? "linear-gradient(118deg, rgba(22,119,255,0.1) 0%, rgba(255,255,255,0.035) 58%, rgba(114,46,209,0.08) 100%)"
            : "linear-gradient(118deg, #f7faff 0%, #eef4ff 58%, #f7f4ff 100%)"};
        box-shadow: none;

        &::before {
            position: absolute;
            inset: 0;
            background-image:
                linear-gradient(rgba(72, 132, 202, 0.055) 1px, transparent 1px),
                linear-gradient(90deg, rgba(72, 132, 202, 0.055) 1px, transparent 1px);
            background-size: 24px 24px;
            content: "";
            -webkit-mask-image: linear-gradient(90deg, transparent 16%, rgba(0, 0, 0, 0.18) 55%, #000 100%);
            mask-image: linear-gradient(90deg, transparent 16%, rgba(0, 0, 0, 0.18) 55%, #000 100%);
            pointer-events: none;
        }

        &::after {
            position: absolute;
            top: 50%;
            right: 20%;
            width: 150px;
            height: 150px;
            border-radius: 50%;
            background: color-mix(in srgb, ${token.colorPrimary}, transparent 90%);
            content: "";
            filter: blur(48px);
            pointer-events: none;
            transform: translateY(-50%);
        }

        .hero-copy {
            position: relative;
            z-index: 2;
            flex: 1;
            min-width: 0;
            max-width: 68%;
            padding: 18px 24px 19px;
        }

        .eyebrow {
            display: inline-flex;
            align-items: center;
            gap: 7px;
            color: color-mix(in srgb, ${token.colorPrimary}, ${token.colorTextTertiary} 62%);
            font-size: 10px;
            font-weight: 600;
            line-height: 17px;
            letter-spacing: 0.09em;
            text-transform: uppercase;
        }

        .status-dot {
            width: 5px;
            height: 5px;
            border-radius: 50%;
            background: color-mix(in srgb, ${token.colorPrimary}, #ffffff 28%);
            box-shadow: 0 0 0 4px color-mix(in srgb, ${token.colorPrimary}, transparent 91%);
        }

        .hero-title {
            margin-top: 3px;
            color: color-mix(in srgb, ${token.colorTextHeading}, transparent 8%);
            font-size: 19px;
            font-weight: 600;
            line-height: 27px;
            letter-spacing: 0.005em;
        }

        .hero-desc {
            max-width: 480px;
            margin-top: 2px;
            color: ${token.colorTextTertiary};
            font-size: 12px;
            line-height: 19px;
        }

        .create-button {
            position: relative;
            z-index: 3;
            display: inline-flex;
            align-items: center;
            justify-content: center;
            flex: none;
            box-sizing: border-box;
            width: auto;
            min-width: 88px;
            height: 36px;
            margin: 0 24px 0 0;
            padding: 0 17px;
            border: 0;
            border-radius: ${token.borderRadius}px;
            outline: none;
            appearance: none;
            background: ${token.colorPrimary};
            color: #ffffff;
            cursor: pointer;
            font: inherit;
            font-size: 13px;
            font-weight: 500;
            line-height: 1;
            white-space: nowrap;
            transform: none;
            transition: background-color 0.18s ease;
        }

        .create-button:hover,
        .create-button:focus-visible {
            background: ${token.colorPrimaryHover};
            color: #ffffff;
            transform: none;
        }

        .create-button:focus-visible {
            box-shadow: 0 0 0 3px ${token.colorPrimaryBg};
        }

        .create-button:active {
            background: ${token.colorPrimaryActive};
            transform: none;
        }

        .hero-visual {
            position: absolute;
            top: 0;
            right: 106px;
            z-index: 1;
            width: 38%;
            height: 100%;
            opacity: 0.34;
            pointer-events: none;
        }

        .hero-graphic {
            display: block;
            width: 100%;
            height: 100%;
            overflow: visible;
        }

        .hero-orbit {
            fill: none;
            stroke: color-mix(in srgb, ${token.colorPrimary}, transparent 82%);
            stroke-dasharray: 5 7;
            stroke-width: 1;
        }

        .hero-link {
            fill: none;
            stroke: color-mix(in srgb, ${token.colorPrimary}, transparent 62%);
            stroke-linecap: round;
            stroke-width: 1.3;
        }

        .hero-link.secondary {
            stroke: rgba(114, 46, 209, 0.3);
        }

        .hero-node {
            fill: ${isDarkMode ? "rgba(255,255,255,0.045)" : "rgba(255,255,255,0.58)"};
            stroke: color-mix(in srgb, ${token.colorPrimary}, transparent 58%);
            stroke-width: 1;
        }

        .hero-node.secondary {
            stroke: rgba(114, 46, 209, 0.34);
        }

        .hero-core {
            fill: color-mix(in srgb, ${token.colorPrimary}, transparent 42%);
            stroke: color-mix(in srgb, ${token.colorPrimary}, #ffffff 48%);
            stroke-width: 1;
        }

        .hero-dot {
            fill: color-mix(in srgb, ${token.colorPrimary}, #ffffff 22%);
        }

        @media (max-width: 720px) {
            min-height: 126px;

            .hero-copy {
                max-width: calc(100% - 118px);
                padding: 16px 18px;
            }

            .hero-title {
                font-size: 17px;
                line-height: 24px;
            }

            .hero-desc {
                display: -webkit-box;
                overflow: hidden;
                -webkit-line-clamp: 2;
                -webkit-box-orient: vertical;
            }

            .create-button {
                height: 36px;
                margin-right: 16px;
            }

            .hero-visual {
                right: 76px;
                width: 50%;
                opacity: 0.2;
            }
        }

        @media (max-width: 420px) {
            min-height: 132px;

            .hero-copy {
                max-width: calc(100% - 104px);
                padding: 15px 8px 15px 16px;
            }

            .eyebrow {
                letter-spacing: 0.06em;
            }

            .hero-title {
                font-size: 15px;
                line-height: 22px;
            }

            .hero-desc {
                font-size: 11px;
                line-height: 17px;
            }

            .create-button {
                margin-right: 12px;
                padding: 0 14px;
            }

            .hero-visual {
                right: 58px;
                width: 64%;
                opacity: 0.13;
            }
        }
    `,
    documentCard: css`
        position: relative;
        min-width: 0;
        min-height: 184px;
        overflow: hidden;
        border-radius: ${token.borderRadiusLG}px;
        background: ${token.colorBgContainer};
        cursor: pointer;
        transition: border-color 0.16s ease, transform 0.16s ease;

        &.ant-card > .ant-card-body {
            padding: 0 !important;
        }

        &:hover {
            border-color: ${token.colorPrimary};
        }

        .preview {
            position: relative;
            height: 104px;
            overflow: hidden;
            border-bottom: 1px solid ${token.colorBorderSecondary};
            background: ${isDarkMode ? "rgba(255,255,255,0.012)" : "#fbfcff"};
        }

        .mind-preview {
            display: block;
            width: 100%;
            height: 100%;
        }

        .preview-surface {
            fill: ${isDarkMode ? "rgba(255,255,255,0.012)" : "#fbfcff"};
        }

        .grid-line {
            fill: none;
            stroke: ${isDarkMode ? "rgba(255,255,255,0.035)" : "rgba(15, 23, 42, 0.035)"};
            stroke-width: 1;
        }

        .preview-link {
            fill: none;
            stroke: color-mix(in srgb, var(--accent), ${token.colorTextTertiary} 28%);
            stroke-linecap: round;
            stroke-width: 1.8;
        }

        .preview-root,
        .preview-node {
            fill: ${token.colorBgContainer};
            stroke: color-mix(in srgb, var(--accent), ${token.colorBorderSecondary} 46%);
            stroke-width: 1.2;
        }

        .preview-dot {
            fill: color-mix(in srgb, var(--accent), ${token.colorBgContainer} 18%);
        }

        .body {
            padding: 11px 13px 12px;
        }

        .head {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 8px;
            min-height: 24px;
        }

        .actions {
            display: flex;
            align-items: center;
            gap: 2px;
            flex: none;
            margin-right: 0;
        }

        .title {
            display: block;
            flex: 1 1 auto;
            min-width: 0;
            padding: 0;
            border: 0;
            background: transparent;
            color: ${token.colorTextHeading};
            cursor: pointer;
            font-size: 14px;
            font-weight: 600;
            line-height: 21px;
            text-align: left;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .title:hover {
            color: ${token.colorPrimary};
        }

        .desc {
            display: -webkit-box;
            min-height: 18px;
            margin-top: 3px;
            overflow: hidden;
            color: ${token.colorTextTertiary};
            font-size: 12px;
            line-height: 18px;
            -webkit-line-clamp: 1;
            -webkit-box-orient: vertical;
        }

        .meta {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 8px;
            margin-top: 7px;
            color: ${token.colorTextTertiary};
            font-size: 12px;
            line-height: 18px;
        }

        .meta span {
            min-width: 0;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        @media (max-width: 640px) {
            min-height: 0;
            border-color: ${token.colorBorderSecondary};

            > .ant-card-body {
                display: block;
            }

            .preview {
                height: 124px;
                min-height: 0;
                border-right: 0;
                border-bottom: 1px solid ${token.colorBorderSecondary};
            }

            .mind-preview {
                transform: scale(1.08);
            }

            .body {
                min-width: 0;
                padding: 12px 14px 13px;
            }

            .head {
                gap: 8px;
                min-height: 24px;
            }

            .title {
                font-size: 14px;
                line-height: 21px;
            }

            .desc {
                -webkit-line-clamp: 1;
                min-height: 18px;
                margin-top: 3px;
            }

            .meta {
                flex-wrap: nowrap;
                justify-content: space-between;
                gap: 12px;
                margin-top: 8px;
            }
        }

        @media (max-width: 420px) {
            .preview {
                height: 112px;
            }

            .body {
                padding: 11px 12px 12px;
            }
        }

    `,
    editor: css`
        position: relative;
        height: var(--mind-page-height);
        min-height: min(680px, var(--mind-page-height));
        overflow: hidden;
        border: 1px solid ${token.colorBorderSecondary};
        border-radius: ${token.borderRadiusLG}px;
        background: ${token.colorBgContainer};
    `,
    editorDrawer: css`
        .ant-drawer-content-wrapper {
            width: 100vw !important;
            max-width: 100vw;
        }

        .ant-drawer-content {
            height: 100dvh;
            overflow: hidden;
            border-radius: 0 !important;
        }

        .ant-drawer-body {
            height: 100%;
            padding: 0;
            overflow: hidden;
            border-radius: 0 !important;
        }

        @supports not (height: 100dvh) {
            .ant-drawer-content {
                height: 100vh;
            }
        }
    `,
    editorToolbar: css`
        position: absolute;
        top: 12px;
        left: 12px;
        z-index: 20;
        display: flex;
        align-items: center;
        min-width: 0;
        max-width: calc(100vw - 68px);

        .actions {
            display: flex;
            align-items: center;
            gap: 6px;
            min-width: 0;
            padding: 5px;
            border-radius: ${token.borderRadiusLG}px;
            background: ${isDarkMode ? "rgba(30,30,30,0.84)" : "rgba(255,255,255,0.9)"};
            box-shadow: ${isDarkMode ? "0 8px 22px rgba(0,0,0,0.3)" : "0 8px 22px rgba(15,23,42,0.08)"};
            backdrop-filter: blur(10px);
        }

        @media (max-width: 780px) {
            .actions {
                max-width: calc(100vw - 68px);
                overflow-x: auto;
                scrollbar-width: none;
            }

            .actions::-webkit-scrollbar {
                display: none;
            }
        }

        @media (max-width: 480px) {
            top: 8px;
            left: 8px;
            max-width: calc(100vw - 56px);

            .actions {
                gap: 2px;
                max-width: calc(100vw - 56px);
                padding: 4px;
            }
        }
    `,
    editorToolButton: css`
        display: inline-flex;
        align-items: center;
        justify-content: center;
        flex: 0 0 auto;
        width: 30px;
        height: 30px;
        padding: 0;
        border: 0;
        border-radius: ${token.borderRadius}px;
        outline: none;
        appearance: none;
        color: ${token.colorTextSecondary};
        background: transparent;
        cursor: pointer;
        transition: color 0.16s ease, background 0.16s ease, transform 0.16s ease;

        .anticon {
            font-size: 14px;
            line-height: 1;
        }

        &:hover,
        &:focus-visible {
            color: ${token.colorPrimary};
            background: ${isDarkMode ? "rgba(255,255,255,0.08)" : token.colorFillTertiary};
        }

        &:active {
            transform: scale(0.96);
        }

        @media (max-width: 480px) {
            width: 28px;
            height: 28px;

            .anticon {
                font-size: 13px;
            }
        }

    `,
    editorClose: css`
        position: absolute;
        top: 12px;
        right: 12px;
        z-index: 20;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        width: 32px;
        height: 32px;
        padding: 0;
        border: 0;
        border-radius: 50%;
        outline: none;
        appearance: none;
        color: ${isDarkMode ? "rgba(255,255,255,0.86)" : token.colorTextSecondary};
        background: ${isDarkMode ? "rgba(42,42,42,0.86)" : "rgba(255,255,255,0.86)"};
        box-shadow: ${isDarkMode ? "0 8px 22px rgba(0,0,0,0.32)" : "0 8px 22px rgba(15,23,42,0.08)"};
        backdrop-filter: blur(10px);
        cursor: pointer;
        transition: color 0.16s ease, background 0.16s ease, transform 0.16s ease;

        .anticon {
            font-size: 14px;
            line-height: 1;
        }

        &:hover,
        &:focus-visible {
            color: ${token.colorError};
            background: ${isDarkMode ? "rgba(255,77,79,0.16)" : token.colorErrorBg};
        }

        &:active {
            transform: scale(0.96);
        }

        &:disabled {
            opacity: 0.62;
            cursor: wait;
            transform: none;
        }

        @media (max-width: 480px) {
            top: 8px;
            right: 8px;
            width: 30px;
            height: 30px;
        }
    `,
    editorFullscreen: css`
        width: 100vw;
        height: 100dvh !important;
        min-height: 0 !important;
        border: 0;
        border-radius: 0;
        box-shadow: none;

        @supports not (height: 100dvh) {
            height: 100vh !important;
        }

        @media (max-width: 780px) {
            & [data-mind-toolbar="true"] {
                top: 62px;
                right: auto;
                left: 50%;
                justify-content: center;
                flex-wrap: nowrap;
                width: max-content;
                max-width: calc(100vw - 24px);
                overflow-x: auto;
                overflow-y: hidden;
                overscroll-behavior-x: contain;
                scrollbar-width: none;
                touch-action: pan-x;
                transform: translateX(-50%);
            }

            & [data-mind-toolbar="true"]::-webkit-scrollbar {
                display: none;
            }
        }

        @media (max-width: 480px) {
            & [data-mind-toolbar="true"] {
                top: 52px;
                max-width: calc(100vw - 16px);
            }
        }
    `,
    editorLoading: css`
        position: absolute;
        inset: 0;
        display: flex;
        align-items: center;
        justify-content: center;
        background: ${isDarkMode ? "rgba(255,255,255,0.01)" : "#fbfcfe"};
    `,
    canvas: css`
        position: absolute;
        inset: 0;
        width: 100%;
        height: 100%;
        min-width: 0;
        min-height: 0;
        background: ${isDarkMode ? "rgba(255,255,255,0.01)" : "#fbfcfe"};
    `,
    versionDrawer: css`
        .ant-drawer-content-wrapper {
            width: min(420px, 100vw) !important;
            max-width: 100vw;
            transition-property: transform;
        }

        .ant-drawer-section {
            overflow: hidden;
            border-radius: ${token.borderRadiusLG}px 0 0 ${token.borderRadiusLG}px;
        }

        .ant-drawer-header {
            min-height: 56px;
            padding: 14px 18px;
            border-bottom-color: ${token.colorBorderSecondary};
        }

        .ant-drawer-title {
            color: ${token.colorTextHeading};
            font-size: 15px;
            font-weight: 600;
        }

        .ant-drawer-body {
            display: flex;
            min-height: 0;
            padding: 0;
            overflow: hidden;
        }

        @media (max-width: 640px) {
            .ant-drawer-section {
                border-radius: 0;
            }

            .ant-drawer-header {
                padding-inline: 16px;
            }

            .ant-drawer-body {
                padding-bottom: env(safe-area-inset-bottom);
            }
        }
    `,
    versionList: css`
        flex: 1;
        min-height: 0;
        overflow: auto;
        overscroll-behavior: contain;
        background: ${token.colorBgContainer};

        .version-filter {
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 14px 18px 12px;
        }

        .version-filter .ant-segmented {
            width: auto;
            padding: 2px;
            background: ${token.colorFillTertiary};
        }

        .version-filter .ant-segmented-item-label {
            min-height: 28px;
            padding-inline: 12px;
            font-size: 12px;
            line-height: 28px;
        }

        .timeline {
            display: flex;
            flex-direction: column;
            margin: 0;
            padding: 2px 18px 18px;
            list-style: none;
        }

        .version {
            position: relative;
            display: grid;
            grid-template-columns: 24px minmax(0, 1fr);
            gap: 10px;
            min-width: 0;
            padding-bottom: 10px;
        }

        .version:last-child {
            padding-bottom: 0;
        }

        .version:last-child .line {
            display: none;
        }

        .version.current .dot {
            background: ${token.colorPrimary};
        }

        .axis {
            position: relative;
            display: flex;
            align-items: flex-start;
            align-self: stretch;
            justify-content: center;
            padding-top: 15px;
        }

        .dot {
            position: relative;
            z-index: 1;
            width: 9px;
            height: 9px;
            border: 2px solid ${token.colorPrimary};
            border-radius: 50%;
            background: ${token.colorBgContainer};
        }

        .line {
            position: absolute;
            top: 0;
            bottom: -12px;
            left: 50%;
            width: 1px;
            transform: translateX(-50%);
            background: ${token.colorBorderSecondary};
        }

        .body {
            position: relative;
            box-sizing: border-box;
            width: 100%;
            min-width: 0;
            padding: 10px 12px;
            border-radius: ${token.borderRadius}px;
            background: ${token.colorFillQuaternary};
            cursor: pointer;
            transition: background ${token.motionDurationSlow};
        }

        .body::before {
            position: absolute;
            top: 14px;
            left: -6px;
            width: 0;
            height: 0;
            border-top: 6px solid transparent;
            border-right: 6px solid ${token.colorFillQuaternary};
            border-bottom: 6px solid transparent;
            content: "";
            transition: border-right-color ${token.motionDurationSlow};
        }

        @media (hover: hover) and (pointer: fine) {
            .body:hover {
                background: ${token.colorFillSecondary};
            }

            .body:hover::before {
                border-right-color: ${token.colorFillSecondary};
            }
        }

        .head {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 8px;
            min-width: 0;
            min-height: 20px;
            margin-bottom: 6px;
        }

        .meta {
            display: flex;
            flex: 1 1 auto;
            align-items: center;
            flex-wrap: nowrap;
            gap: 7px;
            min-width: 0;
            overflow: hidden;
        }

        .snapshot-version {
            flex: none;
            padding-inline: 6px;
            border-radius: ${token.borderRadiusSM}px;
            background: ${token.colorPrimaryBg};
            color: ${token.colorPrimary};
            font-size: 11px;
            font-weight: 600;
            line-height: 20px;
        }

        .current-label {
            flex: none;
            color: ${token.colorPrimary};
            font-size: 11px;
            line-height: 20px;
            white-space: nowrap;
        }

        .snapshot-time {
            flex: 0 1 auto;
            min-width: 0;
            max-width: 45%;
            overflow: hidden;
            color: ${token.colorTextTertiary};
            font-size: 11px;
            line-height: 18px;
            font-variant-numeric: tabular-nums;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .snapshot-text {
            min-height: 18px;
            overflow: hidden;
            color: ${token.colorTextSecondary};
            font-size: 12px;
            line-height: 18px;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .snapshot-extra {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 8px;
            min-height: 28px;
            margin-top: 6px;
        }

        .snapshot-stats {
            display: flex;
            flex: 1 1 auto;
            align-items: center;
            flex-wrap: nowrap;
            gap: 5px 10px;
            min-width: 0;
            overflow: hidden;
            color: ${token.colorTextTertiary};
            font-size: 11px;
            line-height: 18px;
            white-space: nowrap;
        }

        .preview-button.ant-btn {
            flex: none;
            height: 28px;
            min-width: 48px;
            padding-inline: 4px;
            font-size: 12px;
        }

        .preview-button.ant-btn:focus-visible {
            outline: 1px solid ${token.colorPrimaryBorder};
            outline-offset: -2px;
            box-shadow: none;
        }

        .empty {
            display: flex;
            align-items: center;
            justify-content: center;
            min-height: 240px;
            padding: 24px;
        }

        @media (max-width: 640px) {
            .version-filter {
                padding-inline: 16px;
            }

            .timeline {
                padding-inline: 16px;
            }

            .version {
                grid-template-columns: 18px minmax(0, 1fr);
                gap: 8px;
            }
        }
    `,
    versionPreviewModal: css`
        .ant-modal {
            width: min(900px, calc(100vw - 32px)) !important;
            max-width: calc(100vw - 32px);
        }

        @media (max-width: 640px), (max-width: 960px) and (max-height: 520px) {
            .ant-modal-wrap {
                box-sizing: border-box;
                overflow: hidden;
                padding:
                    max(16px, env(safe-area-inset-top))
                    max(12px, env(safe-area-inset-right))
                    max(16px, env(safe-area-inset-bottom))
                    max(12px, env(safe-area-inset-left));
            }

            .ant-modal {
                top: 0;
                width: 100% !important;
                max-width: 900px;
                margin: 0 auto;
                padding: 0;
            }

            .ant-modal-container {
                display: flex;
                flex-direction: column;
                width: 100%;
                height: min(
                    620px,
                    calc(
                        100vh
                        - max(16px, env(safe-area-inset-top))
                        - max(16px, env(safe-area-inset-bottom))
                    )
                );
                height: min(
                    620px,
                    calc(
                        100dvh
                        - max(16px, env(safe-area-inset-top))
                        - max(16px, env(safe-area-inset-bottom))
                    )
                );
                max-height: 100%;
                min-height: 0;
                border-radius: ${token.borderRadiusLG}px;
                overflow: hidden;
            }

            .ant-modal-header {
                flex: none;
            }

            .ant-modal-body {
                display: flex;
                flex: 1;
                min-width: 0;
                min-height: 0;
                overflow: hidden;
            }

            .ant-modal-footer {
                flex: none;
            }
        }
    `,
    versionPreview: css`
        display: flex;
        flex-direction: column;
        width: 100%;
        min-width: 0;

        .preview-meta {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 16px;
            margin-bottom: 12px;
        }

        .preview-remark {
            min-width: 0;
            max-width: 100%;
            overflow: hidden;
            color: ${token.colorTextHeading};
            font-size: 13px;
            font-weight: 500;
            line-height: 20px;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .preview-details {
            display: flex;
            align-items: center;
            justify-content: flex-end;
            flex-wrap: wrap;
            gap: 4px 12px;
            color: ${token.colorTextTertiary};
            font-size: 11px;
            line-height: 20px;
            text-align: right;
        }

        .preview-canvas {
            height: clamp(300px, 56dvh, 520px);
            min-height: 0;
            overflow: hidden;
            border: 1px solid ${token.colorBorderSecondary};
            border-radius: ${token.borderRadiusLG}px;
            background: ${token.colorBgLayout};
        }

        @media (max-width: 640px), (max-width: 960px) and (max-height: 520px) {
            flex: 1;
            min-height: 0;

            .preview-meta {
                flex-direction: column;
                align-items: flex-start;
                gap: 5px;
                margin-bottom: 10px;
            }

            .preview-details {
                width: 100%;
                max-width: 100%;
                justify-content: flex-start;
                text-align: left;
            }

            .preview-canvas {
                flex: 1;
                height: auto;
                min-height: 0;
            }
        }

        @media (min-width: 480px) and (max-width: 960px) and (max-height: 520px) {
            .preview-meta {
                flex-direction: row;
                align-items: center;
                gap: 16px;
                margin-bottom: 8px;
            }

            .preview-details {
                width: auto;
                justify-content: flex-end;
                text-align: right;
            }
        }
    `,
    versionPreviewLoading: css`
        display: flex;
        align-items: center;
        justify-content: center;
        width: 100%;
        height: clamp(332px, calc(56dvh + 32px), 552px);

        @media (max-width: 640px), (max-width: 960px) and (max-height: 520px) {
            flex: 1;
            height: auto;
            min-height: 0;
        }
    `,
    mindForm: css`
        .ant-form-item {
            margin-bottom: 16px;
        }

        .ant-form-item:last-child {
            margin-bottom: 0;
        }

        .ant-form-item-label {
            padding-bottom: 5px;
        }

        .ant-form-item-label > label {
            color: ${token.colorTextSecondary};
            font-size: 12px;
        }

        .ant-input,
        .ant-select-selector {
            font-size: 13px;
        }

        textarea.ant-input {
            min-height: 88px;
            resize: none;
        }

        .share-expire-segmented {
            width: 100%;
        }

        .share-expire-segmented .ant-segmented-item-label {
            min-height: 28px;
            padding-inline: 4px;
            font-size: 12px;
            line-height: 28px;
        }
    `,
    shareDrawer: css`
        .ant-drawer-content-wrapper {
            width: min(420px, 100vw) !important;
            max-width: 100vw;
        }

        .ant-drawer-section {
            overflow: hidden;
            border-radius: ${token.borderRadiusLG}px 0 0 ${token.borderRadiusLG}px;
        }

        .ant-drawer-header {
            min-height: 64px;
            padding: 10px 16px;
            border-bottom-color: ${token.colorBorderSecondary};
        }

        .ant-drawer-title {
            min-width: 0;
        }

        .ant-drawer-extra {
            flex: none;
            margin-inline-start: 12px;
        }

        .share-create-button {
            gap: 4px;
            height: 28px;
            padding-inline: 8px;
            border-radius: ${token.borderRadius}px;
            background: ${token.colorPrimaryBg};
            color: ${token.colorPrimary};
            font-size: 12px;
            font-weight: 500;
        }

        .share-create-button:hover,
        .share-create-button:focus-visible {
            background: ${token.colorPrimaryBg} !important;
            color: ${token.colorPrimary} !important;
        }

        .share-create-button .ant-btn-icon {
            margin-inline-end: 0;
            font-size: 11px;
        }

        .share-title {
            display: flex;
            flex-direction: column;
            gap: 1px;
            min-width: 0;
        }

        .share-title-text {
            color: ${token.colorTextHeading};
            font-size: 15px;
            font-weight: 600;
            line-height: 21px;
        }

        .share-title-document {
            display: block;
            max-width: 100%;
            overflow: hidden;
            color: ${token.colorTextTertiary};
            font-size: 11px;
            font-weight: 400;
            line-height: 17px;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .ant-drawer-body {
            display: flex;
            min-height: 0;
            padding: 0;
            overflow: hidden;
        }

        @media (max-width: 640px) {
            .ant-drawer-section {
                border-radius: 0;
            }

            .ant-drawer-header {
                padding-inline: 14px;
            }
        }
    `,
    sharePanel: css`
        display: flex;
        flex: 1;
        flex-direction: column;
        min-width: 0;
        min-height: 0;
        background: ${token.colorBgContainer};

        .share-list {
            box-sizing: border-box;
            flex: 1;
            min-height: 0;
            padding: 10px 12px 16px;
            overflow: auto;
            overscroll-behavior: contain;
        }

        .share-items {
            display: flex;
            flex-direction: column;
            gap: 8px;
        }

        .share-item {
            min-width: 0;
            padding: 10px 12px;
            border-radius: ${token.borderRadius}px;
            background: ${token.colorFillQuaternary};
            transition: background ${token.motionDurationSlow};
        }

        @media (hover: hover) and (pointer: fine) {
            .share-item:hover {
                background: ${token.colorFillTertiary};
            }
        }

        .share-item-head {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 10px;
        }

        .share-code-line {
            display: flex;
            align-items: center;
            gap: 8px;
            min-width: 0;
        }

        .share-code {
            min-width: 0;
            overflow: hidden;
            color: ${token.colorTextHeading};
            font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
            font-size: 14px;
            font-weight: 600;
            line-height: 20px;
            letter-spacing: 0;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .share-state {
            display: inline-flex;
            flex: none;
            align-items: center;
            gap: 5px;
            color: ${token.colorTextTertiary};
            font-size: 11px;
            line-height: 17px;
            white-space: nowrap;
        }

        .share-state i {
            display: block;
            width: 6px;
            height: 6px;
            border-radius: 50%;
            background: currentColor;
        }

        .share-item[data-state="enabled"] .share-state {
            color: ${token.colorSuccess};
        }

        .share-item[data-state="expired"] .share-state {
            color: ${token.colorWarning};
        }

        .share-remark {
            min-height: 18px;
            margin-top: 4px;
            overflow: hidden;
            color: ${token.colorTextSecondary};
            font-size: 12px;
            line-height: 18px;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .share-remark.empty {
            color: ${token.colorTextTertiary};
        }

        .share-meta {
            display: flex;
            align-items: center;
            gap: 7px;
            min-width: 0;
            margin-top: 6px;
            color: ${token.colorTextTertiary};
            font-size: 11px;
            line-height: 17px;
        }

        .share-meta span {
            min-width: 0;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .share-meta span:first-child {
            flex: 1 1 auto;
        }

        .share-meta span:last-child {
            flex: 0 1 auto;
        }

        .share-meta > i {
            flex: none;
            width: 2px;
            height: 2px;
            border-radius: 50%;
            background: ${token.colorTextQuaternary};
        }

        .share-action-button {
            flex: none;
            width: 26px;
            min-width: 26px;
            height: 26px;
            padding: 0;
            color: ${token.colorTextTertiary};
        }

        .share-empty,
        .share-loading {
            display: flex;
            align-items: center;
            justify-content: center;
            min-height: 100%;
        }

        .share-load-more {
            display: flex;
            align-items: center;
            justify-content: center;
            min-height: 40px;
        }

        .share-sentinel {
            width: 100%;
            height: 1px;
        }

        @media (max-width: 640px) {
            .share-list {
                padding: 10px 10px max(14px, env(safe-area-inset-bottom));
            }

            .share-item {
                padding: 10px 11px;
            }
        }
    `,
}));

export default useStyles;
