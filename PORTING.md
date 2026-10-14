# Porting guide

This is the contract for porting a shadcn/ui component to Go. The bar is
**exact parity** with upstream (shadcn/ui `new-york-v4`, React 19, radix-ui
1.4.3): every pixel in light and dark mode, the DOM after every interaction,
and every animation frame. The parity harness proves it; nothing is done until
it passes.

## Where things are

| What | Where |
| --- | --- |
| Upstream component source | `ui/apps/v4/registry/new-york-v4/ui/<name>.tsx` |
| Upstream examples | `ui/apps/v4/registry/new-york-v4/examples/<name>-*.tsx`, also `charts/`, `blocks/<name>/`, `internal/` |
| Upstream dependencies, exact lockfile versions | `ui/apps/v4/node_modules/` (pnpm). Radix primitives: `ui/node_modules/.pnpm/@radix-ui+react-<x>@*/node_modules/@radix-ui/react-<x>/dist/index.mjs` (readable source) |
| Go component | `components/ui/<name>/<name>.go` (one package per component) |
| Go component behavior (JS) | `components/ui/<name>/<name>.js`, registered with `runtime.Module` |
| Shared JS primitives (ported Radix internals) | `components/ui/runtime/js/*.js` |
| Go examples | `examples/<name>/<name>.go` (package per component), registered in `init` |
| Interaction scenarios | `parity/scenarios/<reference>.json` |
| Fixtures for components without upstream examples | `parity/fixtures/<name>.tsx` (React, using the upstream component unmodified) |
| Icons | `shadcn.go/components/icons/lucide` (all of lucide-react 0.474.0), `.../tabler` |
| Helpers | `shadcn.go/components/lib` (`lib.Cn`, `lib.AsChild`), `shadcn.go/cva`, `shadcn.go/libhtml` (`h`) |

Study `components/ui/button` (component), `components/ui/spinner` (`h.Spread`)
and `examples/button` (examples) first: they pass the harness and show every
convention below.

## Conventions

**API.** Each sub-component is a function in the component's package, named
after the upstream export without the component prefix, taking `...any`
children JSX-style:

```go
dialog.Dialog(
    dialog.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Open")),
    dialog.Content(
        dialog.Header(dialog.Title("Edit profile"), dialog.Description("...")),
        ...
    ),
)
```

The root function is named like the package (`card.Card`, `dialog.Dialog`).
Variants are typed string constants (`button.VariantOutline`); booleans and
other props are typed values accepted as children through `ChildAcceptor` and
also offered as builder methods. Use `h.Builder[T]` for anything that takes
attributes (see button). Keep upstream names: `Header`, `Footer`,
`Title`, `Description`, `Content`, `Trigger`, `Close`, `Item`, ...

**Markup.** Render exactly what React renders: same tags, same attributes
(`data-slot`, `data-state`, `aria-*`, `role`, ...), same class strings
(`lib.Cn` of the same inputs, in the same order), same text, same children
order. Where upstream spreads `{...props}` after its own props, use
`h.Spread(defaults, userElement, lib.Cn)`; for `asChild` use `h.Slot`. Radix
primitives render attributes you do not see in the shadcn source (e.g.
`data-orientation`, `aria-controls`, `type="button"`, `dir`, hidden
inputs, `style` with CSS variables): read the Radix source and the harness DOM
dumps. Server-render the initial state (including `data-state="closed"`,
`hidden`, etc.) exactly as React's first render does, so the page is correct
before JavaScript runs.

**Ids.** Radix generates ids (`useId`). The harness normalizes ids, so any
unique ids work, but references (`aria-controls`, `aria-labelledby`, `for`,
...) must point to the right element. Generate ids server-side with a
per-render counter or `crypto/rand`, never global mutable state shared across
requests without synchronization.

**Behavior.** Interactivity is vanilla JavaScript, no framework, no build
step: ES modules served by `components/ui/runtime`. A component registers its
module from Go:

```go
//go:embed dialog.js
var script string

func init() { runtime.Module("dialog", script) }
```

and the module attaches to server-rendered markup:

```js
import { define } from "./core.js"
define("dialog", { selector: '[data-slot="dialog"]', mount(root) { ...; return cleanup } })
```

Port the Radix behavior faithfully: state attributes and their timing, focus
management, keyboard map, dismissal (outside click, Escape), pointer events,
portal placement (content moved to `<body>` where Radix portals it), scroll
lock, `aria-hidden` on outside content, presence (unmount after the exit
animation ends), positioning (Radix popper on the vendored
`runtime/js/floating-ui.js`, same middleware, same CSS variables such as
`--radix-popper-transform-origin`). Shared mechanisms belong in
`runtime/js/` as primitives (`presence.js`, `dismissable-layer.js`,
`focus-scope.js`, `popper.js`, `roving-focus.js`, ...) so components reuse
one implementation. Dispatch `shadcn:<event>` CustomEvents (via `emit` in
core.js) for state changes applications may observe. The Go pages run under a
strict CSP: no inline scripts or `<style>` without the nonce, no `eval`.

**Third-party libraries** (vaul, sonner, embla, cmdk, react-day-picker,
input-otp, react-resizable-panels, recharts, base-ui, ...): reproduce their
output and behavior. If the library has a framework-free core with the same
behavior (e.g. `embla-carousel`), you may vendor it into `runtime/js/` as a
single unmodified ESM bundle with its license header (see `floating-ui.js`
for how it was produced), at the exact version the upstream lockfile uses.
Otherwise port its logic to vanilla JS (or to Go, for rendering such as chart
SVG). Libraries that inject `<style>` at runtime: ship those rules in the
stylesheet or a nonced `<style>`, with the identical CSS.

**Security.** Everything goes through libhtml. Never use `UnsafeHTML`,
`UnsafeJS`, `UnsafeAttr` or `UnsafeURL` with data that could come from users.
Generated or vendored constants only.

**Isolation (many agents work concurrently).** Only edit files you own:
your component packages, your example packages, your `.js` modules, your
scenario and fixture files. Shared primitives in `runtime/js/` are owned by the
agent that created them; if you must change one, keep the change minimal and
backward compatible, and re-run the harness for every component that imports
it. Never edit `libhtml`, `tailwind`, `cmd/parity`, `parity/run.mjs`,
`parity/reference/build.mjs` or another agent's package; if you need a change
there, describe it in your final report. Do not run `git` commands that
modify the repository (no commits, checkout, stash, reset).

## The harness

```sh
node parity/reference/build.mjs            # after adding a fixture (others' builds are unaffected)
node parity/run.mjs <reference>...         # or --component <name>
node parity/run.mjs --component dialog --themes light --keep   # keep images of passing captures too
```

Many agents run it at once: always pass `--jobs 2` and run only the
references you are working on, not `--all`. It builds and starts `cmd/parity`, opens `/react/<name>` and `/go/<name>` in
headless Chromium (1024x768) and compares, per theme and per capture:

- **pixels**: full-page screenshots, 0 differing pixels allowed;
- **DOM**: normalized `<body>` (ids, attribute and class order, style
  serialization normalized; focus marked) must be identical;
- **animations**: `document.getAnimations()` definitions (keyframes,
  duration, delay, easing, ...) must be identical and frames sampled at 0,
  25, 50 and 75% of their duration must be pixel-identical (the document
  timeline is frozen with the DevTools protocol before each interaction).

Captures: `initial`, then automatic interactions (hover, Tab focus, click the
first trigger, Escape), then your scenario. Differences are written to
`parity/out/<reference>/<theme>/`: `*-react.png`, `*-go.png`, `*-diff.png`
(red = differing pixel), `*.dom.txt`, and `report.json`. Look at them.

**Scenarios.** For every interactive component write
`parity/scenarios/<reference>.json` covering every behavior upstream has:
each keyboard key in the Radix keyboard map, pointer interactions (hover,
click, outside click, drag where relevant), opening and closing (with
`"frames"` to compare animations), focus movement, disabled states, typing.

```json
[
  { "label": "open", "steps": [{ "do": "click", "target": "[data-slot=dialog-trigger]" }] },
  { "label": "tab-cycle", "steps": [{ "do": "press", "key": "Tab" }, { "do": "press", "key": "Tab" }] },
  { "label": "close-escape", "steps": [{ "do": "press", "key": "Escape" }] },
  { "label": "no-anim", "frames": false, "steps": [{ "do": "hover", "target": "button", "nth": 1 }] }
]
```

Step kinds: `hover`, `click`, `dblclick`, `focus`, `press` (key, optional
target), `type` (text), `mouse` (x, y, optional steps), `down` / `up` (mouse
button, optional `button`), `drag` (target, dx, dy, optional steps), `wait`
(ms), `scroll` (y, optional target), `eval` (js). Add `"force": true` to
hover/click/dblclick to skip Playwright's actionability checks (e.g. when an
overlay covers the target). Before every screenshot the harness forces a full
repaint and, after load, waits for page-load CSS transitions, so scenarios do
not need workarounds for either. Targets are CSS selectors (`nth` picks among matches); use
`data-slot` selectors, which both sides share. Captures accumulate: each
scenario continues from the previous state, unless it has `"reset": true`,
which reloads both pages first.

## Definition of done

For each assigned component:

1. `go vet ./...` and `go test` for your packages pass. Add unit tests for
   anything non-trivial (variants, prop handling, escaping).
2. Every upstream reference assigned to you (examples, charts, internal demos,
   blocks, plus fixtures you add for components without upstream examples) is
   ported and **passes `node parity/run.mjs` in light and dark**, with
   scenarios covering the component's full behavior.
3. Upstream references that use components owned by other agents which do not
   exist yet are listed as blocked in your report (port them later or leave
   them to the examples sweep).
4. Your final report lists: components done, references passing, references
   failing or blocked with the exact reason, shared files you touched and why,
   and anything another agent needs to know (e.g. a primitive's API).

If something cannot reach exact parity (for example an upstream bug that
depends on React internals), document it precisely in the report with the
harness output; do not loosen the harness.
