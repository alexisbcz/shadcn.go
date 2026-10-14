// Interaction scenarios of the documentation site chrome, replayed by
// parity/docs/run.mjs on both servers after the "initial" capture.
//
// A scenario applies to a page when its `when` selector matches a visible
// element on the REFERENCE page (both sides then replay the same steps; a
// step the candidate cannot perform is a failure), when the viewport is in
// `viewports` (default all) and the page group is in `groups` (default all).
//
// Fields:
//   label      capture name (unique per page)
//   when       CSS selector (Playwright syntax) that must be visible on the reference
//   viewports  subset of ["desktop", "laptop", "mobile"]
//   groups     subset of ["docs", "app", "view"]
//   reset      reload both pages before the scenario (default: continue from the previous state)
//   frames     false to skip the animation comparison (default: compare definitions and sample frames)
//   full       true for a full-page screenshot (default: the viewport, where the interaction is)
//   steps      see run.mjs "Steps"
//
// Captures accumulate like in parity/run.mjs: every scenario starts from the
// state the previous one left, unless it has `reset: true`.
//
// Per-page scenarios can be added in parity/docs/scenarios/<path with / as __>.json
// (same format); they run after these.

// Shared selectors. They only use what the markup itself carries
// (data-slot, roles, sr-only labels, text), so both sides can share them.
export const S = {
  themeToggle: 'header button:has(> .sr-only:text-is("Toggle theme"))',
  searchTrigger: 'header button[data-slot="dialog-trigger"]',
  searchInput: '[role="dialog"] [cmdk-input], [role="dialog"] input',
  mobileMenu: 'header button[data-slot="popover-trigger"]:has-text("Toggle Menu")',
  sidebarLink: '[data-slot="sidebar"] a[data-slot="sidebar-menu-button"]',
  baseTab: '[data-slot="docs"] div.not-typeset > a[data-active="false"]',
  codeTabsManual: '[data-slot="tabs-trigger"]:text-is("Manual")',
  codeTabsCommand: '[data-slot="tabs-trigger"]:text-is("Command")',
  pmNpm: '[data-slot="tabs-trigger"]:text-is("npm")',
  pmBun: '[data-slot="tabs-trigger"]:text-is("bun")',
  copyButton: 'button[data-slot="copy-button"]',
  viewCode: '[data-slot="component-preview"] button:text-is("View Code")',
  toc: 'a[data-depth][data-active]',
  copyPage: '[data-slot="docs"] .docs-nav button:has-text("Copy Page")',
  copyPageMenu: '[data-slot="docs"] .docs-nav [data-slot="dropdown-menu-trigger"]',
  copyPagePopover: '[data-slot="docs"] .docs-nav [data-slot="popover-trigger"]',
}

export const scenarios = [
  // Theme toggle (header button), then the "d" keyboard shortcut back.
  { label: "theme-toggle", when: S.themeToggle, steps: [{ do: "click", target: S.themeToggle }, { do: "settle" }], full: true },
  { label: "theme-shortcut-d", when: S.themeToggle, steps: [{ do: "mouse", x: 1, y: 1 }, { do: "press", key: "d" }, { do: "settle" }], full: true },

  // Search: open with a click, type a query, close; then open with Cmd+K.
  { label: "search-open-click", when: S.searchTrigger, steps: [{ do: "click", target: S.searchTrigger }, { do: "settle" }] },
  { label: "search-type", when: S.searchTrigger, frames: false, steps: [{ do: "type", text: "button" }, { do: "wait", ms: 700 }, { do: "settle" }] },
  { label: "search-arrow", when: S.searchTrigger, frames: false, steps: [{ do: "press", key: "ArrowDown" }, { do: "press", key: "ArrowDown" }, { do: "settle" }] },
  { label: "search-close", when: S.searchTrigger, steps: [{ do: "press", key: "Escape" }, { do: "settle" }] },
  { label: "search-cmdk", when: S.searchTrigger, steps: [{ do: "press", key: "ControlOrMeta+k" }, { do: "settle" }] },
  { label: "search-cmdk-close", when: S.searchTrigger, steps: [{ do: "press", key: "Escape" }, { do: "settle" }] },

  // Mobile menu (the header popover below lg).
  { label: "mobile-menu-open", when: S.mobileMenu, steps: [{ do: "click", target: S.mobileMenu }, { do: "settle" }] },
  { label: "mobile-menu-close", when: S.mobileMenu, steps: [{ do: "click", target: S.mobileMenu, force: true }, { do: "settle" }] },

  // Package-manager tabs and code tabs (the CLI / manual installation).
  { label: "pm-tab-npm", when: S.pmNpm, steps: [{ do: "click", target: S.pmNpm }, { do: "settle" }] },
  { label: "pm-tab-bun", when: S.pmBun, steps: [{ do: "click", target: S.pmBun }, { do: "settle" }] },
  { label: "code-tab-manual", when: S.codeTabsManual, steps: [{ do: "click", target: S.codeTabsManual }, { do: "settle" }] },
  { label: "code-tab-command", when: S.codeTabsCommand, steps: [{ do: "click", target: S.codeTabsCommand }, { do: "settle" }] },

  // Copy buttons: the clipboard text is captured and compared too.
  { label: "copy-button", when: S.copyButton, steps: [{ do: "hover", target: S.copyButton }, { do: "click", target: S.copyButton }, { do: "settle" }, { do: "clipboard" }] },

  // ComponentPreview: the preview/code toggle.
  { label: "view-code", when: S.viewCode, steps: [{ do: "click", target: S.viewCode }, { do: "settle" }] },

  // "Copy Page" and its menu (dropdown from sm, popover below).
  { label: "copy-page", when: S.copyPage, steps: [{ do: "click", target: S.copyPage }, { do: "settle" }, { do: "clipboard" }] },
  { label: "copy-page-menu", when: S.copyPageMenu, steps: [{ do: "click", target: S.copyPageMenu }, { do: "settle" }] },
  { label: "copy-page-menu-close", when: S.copyPageMenu, steps: [{ do: "press", key: "Escape" }, { do: "settle" }] },
  { label: "copy-page-popover", when: S.copyPagePopover, steps: [{ do: "click", target: S.copyPagePopover }, { do: "settle" }] },
  { label: "copy-page-popover-close", when: S.copyPagePopover, steps: [{ do: "press", key: "Escape" }, { do: "settle" }] },

  // TOC scroll-spy: scroll the document in steps; each step is a capture.
  ...[0.2, 0.4, 0.6, 0.8, 1].map((f) => ({ label: `scroll-${Math.round(f * 100)}`, when: S.toc, frames: false, steps: [{ do: "scrollTo", fraction: f }, { do: "settle" }] })),
  { label: "scroll-top", when: S.toc, frames: false, steps: [{ do: "scrollTo", fraction: 0 }, { do: "settle" }] },

  // Sidebar: hover a link, then navigate with it (client-side on Next).
  { label: "sidebar-hover", when: S.sidebarLink, frames: false, steps: [{ do: "hover", target: S.sidebarLink, nth: 1 }, { do: "settle" }] },
  { label: "sidebar-navigate", when: S.sidebarLink, steps: [{ do: "click", target: S.sidebarLink, nth: 1, nav: true }, { do: "settle" }] },

  // Base tabs (Base UI / React Aria / Radix UI): navigate to the next base.
  { label: "base-tab", when: S.baseTab, reset: true, steps: [{ do: "click", target: S.baseTab, nav: true }, { do: "settle" }] },
]
