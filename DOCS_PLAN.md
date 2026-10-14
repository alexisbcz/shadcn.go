# Plan: port the shadcn/ui documentation site, pixel-perfect

Goal: a Go server that renders the upstream documentation site
(`ui/apps/v4`, ui.shadcn.com) so that every page is pixel-, DOM- and
animation-identical to the Next.js original, in light and dark, after every
interaction. The verification method is the same as for components: run the
original and the port side by side and diff.

## What the survey of `ui/apps/v4` found

**1. The docs do not show `new-york-v4`.** Every preview is
`<ComponentPreview styleName="radix-nova" name="button-demo" />`. Components
come from `styles/<base>-<style>/`, which `registry:build` generates by
taking `registry/bases/<base>/ui/*.tsx`, written with semantic `cn-*` class
tokens, and replacing each token with that style's utilities. The token
mapping comes from `registry/styles/style-<style>.css`, for example
`.cn-button-variant-default { @apply ... }`. The transform is
`transformStyle` in `packages/shadcn/src/styles/`.

- There are 3 bases: `radix` (Radix UI), `base` (Base UI) and `aria`
  (React Aria Components).
- There are 8 styles: nova, vega, lyra, maia, mira, luma, sera and rhea.
- Bases differ from `new-york-v4` in structure, not only in classes. For
  example, the Dialog close button is a `Button`, icons go through
  `IconPlaceholder`, which picks lucide, tabler, hugeicons, phosphor or
  remixicon, and some `data-slot`s differ.

**2. Docs examples are a separate tree:**
`apps/v4/examples/{radix,base,aria}/*.tsx`. It has 508 + 513 + 501 = 1,522
files. They import `@/styles/<base>-nova/ui/*`. Their content mostly mirrors
the `new-york-v4` examples.

**3. Content.** There are 316 MDX pages (fumadocs-mdx). 193 of them are
component pages, about 64 per base. The rest are installation, theming,
dark-mode, forms, registry, RTL and the changelog (57 pages).

- The pages use about 30 custom components. The most frequent are
  `ComponentPreview` (1,480 uses), `Steps`/`Step`, `CodeTabs` and
  `Tabs*`, `ComponentSource` (198), `Callout`, `LinkedCard`, `Questionnaire`
  and `ComponentsList`.
- Code is highlighted by `rehype-pretty-code` (shiki).

**4. The app has more than docs.** Its routes include `/`, `/docs/[[...slug]]`,
`/docs/changelog`, `/blocks/[...categories]`, `/charts/[type]`, `/colors`,
`/examples/{dashboard,tasks,playground,authentication,rtl}`, `/create`,
`/typeset`, `/sera`, and the iframe views
`/view/[style]/[name]`, `/preview/[base]/[name]` and
`/examples/[base]/[name]`.

## Decision: re-target the component port to `bases/radix` + styles

The components being ported now follow `new-york-v4`. The docs need
`bases/radix` resolved through a style. The expensive part, the vanilla-JS
behavior that replaces Radix and the other libraries, is identical for both,
because the same Radix primitives sit underneath. So:

1. **Port the style system to Go.** That means `createStyleMap` (parse
   `style-*.css` into token → classes) and `transformStyle`. A Go component
   is authored once with the `cn-*` tokens of `bases/radix`, and the style
   resolves the tokens: `button.Button(...)` renders radix-nova by default,
   and `style.Use("vega")`, or a per-request style, switches. Prove it with a
   differential test: every component source transformed by upstream's
   transform equals the Go output, for all 8 styles.
2. **Icons.** `IconPlaceholder` resolves to the configured library. Generate
   hugeicons, phosphor and remixicon packages, like lucide and tabler, at
   the exact versions in the lockfile.
3. **Keep `new-york-v4`.** The existing component ports and the 349-reference
   harness stay valid as the "new-york-v4" flavor. The Radix-based markup is
   factored so both flavors share the behavior modules, and only the markup
   functions differ where the structure differs. This decides whether a
   wave-1/2 component needs a structural variant or only a class mapping.
4. **The component harness gains the docs examples.** The reference build
   adds `apps/v4/examples/radix/*` (508), rendered with `styles/radix-nova`.
   That requires `registry:build` to have generated `styles/`, see D0.
   `cmd/parity` serves `/go/radix-nova/<name>`. Later, base and aria do the
   same.

## Status

- D0: done. `parity/docs/reference.sh` builds and serves the reference.
  `registry:build` takes about 1 minute and generates all 24 `styles/`;
  `next build` takes about 2 minutes. Pages render fully offline: no
  failed requests and no console errors on /docs/components/radix/button.
  The GitHub star count is baked in at build time, so it is deterministic.

## Phases

**D0. The reference docs site, running locally.** The goal is a
deterministic reference.

- Run upstream's `registry:build`, which needs bun (`npx bun`) and the
  workspace packages built, then `next build && next start --port 4000`.
- Stub what is nondeterministic or remote: GitHub star count, analytics,
  "new" badges based on dates, remote images (cache them), current year.
  Use a small patch file applied to the clone, never committed upstream.
- Vendor the fonts the site uses (Geist via `next/font`) so both sides load
  the same files.

**D1. Page-level parity harness** (extends `parity/run.mjs`).

- Input: a URL list from the docs source tree (316 docs pages plus the app
  routes).
- Compare `http://localhost:4000<path>` with the Go server `<path>` at
  several viewports (1440×900 and 1024×768 for desktop, 390×844 for mobile),
  in light and dark.
- Captures: full-page pixels, the normalized DOM, and scenarios. The
  scenarios cover the sidebar, opening search (⌘K), code-tab switching, the
  copy button, the theme switch, preview/code toggles, TOC scroll-spy, and
  each embedded preview's own scenario, which is reused from the component
  harness.
- Iframe previews (`/view/...`) are compared as their own pages.

**D2. The style system and the re-targeted components.** Agent waves, as in
`parity/manifest.json`:

- Style map and transform in Go, with the differential test.
- Icon libraries.
- One group per component family re-targets the ported components to
  `bases/radix` tokens and passes the 508 radix-nova examples. Components
  not yet ported are ported directly in bases form, and waves 3-4 of the
  component manifest move to this target.

**D3. The content pipeline in Go.**

- MDX: the content uses a small JSX subset (component tags with literal
  props and nested children). Parse Markdown with goldmark and the JSX tags
  with a dedicated parser. Reproduce fumadocs' output exactly: heading ids
  (github-slugger), the TOC, link handling, and the `mdx-components.tsx`
  element mapping and classes. A differential test renders each MDX page
  with upstream's pipeline and the Go one and requires the same HTML.
- Code highlighting: run shiki in goja with its JavaScript regex engine
  (no WASM), as was done for Tailwind. The target is byte-identical tokens
  and colors, and the `rehype-pretty-code` markup (titles, line numbers,
  highlighted lines). The test: every code block of the 316 pages produces
  the same HTML.
- `ComponentSource` shows the upstream TSX source of the style, exactly as
  the original. A Go tab is possible later, as an addition outside parity
  mode.

**D4. Docs chrome.** These are built with the ported components:

- site header (nav, GitHub link, theme toggle, style/base switcher);
- docs sidebar and its mobile drawer;
- page header (title, description, "Copy page" menu, prev/next);
- TOC with scroll-spy;
- search command palette (fumadocs search index: build the same index
  in Go);
- footer, callouts, steps, linked cards, the component list and the
  Questionnaire;
- 404 and the redirects in `next.config.mjs`.

**D5. All docs pages pass.** One agent group per docs section, with a sweep
until every URL passes. The changelog's 57 pages are mostly content and can
be batched.

**D6. The other bases.** `base` (Base UI: different data attributes and
behaviors, such as `data-open` and its transitions) and `aria` (React Aria
Components), each with ~500 examples and its docs pages. These are new
behavior modules, organized like the Radix ones.

**D7. The other app pages:** blocks and charts galleries, colors, the
examples apps (dashboard, tasks, playground, authentication, rtl), typeset,
sera and create (the style/theme builder: the largest interactive page).

## Order and dependencies

Planned order:

1. D0 and D1: the reference must be ready before anything can be proven.
2. D2's style system and D3's pipeline: these can run in parallel, since
   they don't depend on each other.
3. D2's component re-targeting: this runs as part of component waves 3-4.
4. D4.
5. D5, then D6 and D7.

The component port waves in progress (`new-york-v4`) keep running. Their
behavior modules carry over unchanged.

## Risks

- **Size of upstream's `registry:build`.** It needs bun and builds several
  packages. If it can't run here, `styles/` can be generated with only the
  transform (Go or a tsx script) from `bases/` plus `style-*.css`.
- **Shiki in goja.** Its regex usage is heavier than Tailwind's. The
  fallback is to port the TextMate tokenizer onto Go's regexp2 with the same
  grammars and themes, checked against the same differential test.
- **fumadocs internals** (search index format, TOC, slugging) have to be
  reproduced from its source at the lockfile version.
- **Fonts and rendering.** Both sides must load identical font files. The
  harness already enforces a single Chromium build and viewport.
