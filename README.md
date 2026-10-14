# shadcn.go

A port of [shadcn/ui](https://ui.shadcn.com) to Go. Components are plain Go
functions that render HTML on the server. No JavaScript framework is needed.
Two layers make up the project:

- **libhtml**: a library that builds HTML in Go with a JSX-like API. It
  escapes and sanitizes output by default and ships CSRF and
  security-header middleware.
- **components/ui**: the shadcn components, one Go package each, written
  with libhtml. Their client-side behavior is vanilla JavaScript served by
  `components/ui/runtime`.
- **tailwind**: Tailwind CSS v4.3.0 with Go only. `go generate` produces
  the exact stylesheet, with no Node.js required.

```go
import (
    "shadcn.go/components/icons/lucide"
    "shadcn.go/components/lib"
    "shadcn.go/components/ui/button"
    h "shadcn.go/libhtml"
)

h.Div(h.Class("flex gap-2"),
    button.Button(button.VariantOutline, "Button"),
    button.Button(button.VariantOutline, button.SizeIcon, h.Aria("label", "Submit"), lucide.ArrowUp()),
    button.Button(lib.AsChild, h.A(h.Href("/login"), "Login")),
)
```

The port is in progress. `PORTING.md` is the porting contract,
`parity/manifest.json` is the plan, and `node parity/status.mjs` shows
progress.

Run the demo. It shows every shadcn Button example, plus a CSRF form and an
XSS test:

```sh
go run ./cmd/demo -dev        # http://localhost:8080
```

## Layout

| Path | What it is |
| --- | --- |
| `libhtml/` | The HTML builder: elements, props, control flow, escaping, rendering |
| `libhtml/csrf/` | CSRF middleware with origin checks and signed, masked tokens |
| `libhtml/secure/` | Security headers with a per-request CSP nonce |
| `tailwind/`, `cmd/tailwind/` | Tailwind CSS v4.3.0 compiler, scanner and CLI, written in Go only |
| `twmerge/` | A port of tailwind-merge. It resolves Tailwind class conflicts |
| `cva/` | A port of class-variance-authority |
| `components/lib/` | `Cn` (shadcn's `cn()`) and `AsChild` |
| `components/ui/<name>/` | One package per shadcn component |
| `components/ui/runtime/` | The client runtime: vanilla JS replacing Radix UI, served with CSP nonces |
| `components/icons/` | Every Lucide icon (lucide-react 0.474.0) and the Tabler icons the registry uses, generated |
| `examples/` | Go ports of the upstream examples, charts and blocks |
| `parity/`, `cmd/parity/` | The parity harness: React reference vs Go, compared on pixels, DOM and animations |
| `cmd/demo/` | The demo server |
| `bench/` | Benchmarks against `html/template` and gomponents. It is a separate module, so the library has no dependencies |
| `scripts/jsref/` | Produces reference outputs from the real JavaScript libraries, for the fidelity tests |
| `ui/` | The upstream shadcn/ui clone, kept for reference. It is git-ignored |

## libhtml

### Children of any type

Each tag function takes `...any`. Anything can be a child:

```go
h.Div(
    h.Class("card"), h.ID("c1"),           // props become attributes
    "Hello, ", user.Name,                  // strings become escaped text
    42, 3.14, true,                        // numbers and bools become text
    h.Span("nested"),                      // nodes
    []h.Node{a, b},                        // slices are flattened
    h.Group(h.Class("x"), h.P("y")),       // fragments are flattened, and their props apply to the parent
    func() h.Node { return expensive() },  // lazy children are built at render time
    nil,                                   // nil is skipped, including a nil *Element
)
```

The same element can also be built with chained methods:

```go
h.Div().Class("card").ID("c1").Text("Hello, ", user.Name).Children(h.Span("nested"))
```

### Control flow

The helpers are generic. They work for nodes, props, strings and fragments:

```go
h.If(active, h.Class("active"))                // eager: gives the zero value (nothing) when false
h.Unless(loggedIn, h.A(h.Href("/login"), "Sign in"))
h.Iff(user != nil, func() h.Node { return h.Text(user.Name) })   // lazy: safe when user is nil
h.IfElse(ok, "yes", "no")
h.IfElseF(ok, func() h.Node { ... }, func() h.Node { ... })
h.Class(h.Ternary(open, "block", "hidden"))

h.Map(items, func(it Item) h.Node { return h.Li(it.Name) })
h.MapIdx(items, func(i int, it Item) h.Node { ... })
h.MapKV(m, func(k string, v int) h.Node { ... })               // in sorted key order
h.MapSeq(seq, fn)                                                 // iter.Seq
h.Range(n, func(i int) h.Node { ... })
h.Filter(items, keep)

h.Switch(status).
    Case("ok", h.Text("All good")).
    Case("warn", h.Class("text-amber-600")).
    Default(h.Text("Unknown"))
```

Use `If` and `IfElse` when both branches are cheap and safe to build. Go
evaluates their arguments before the call. Use `Iff` and `IfElseF` when a
branch dereferences something that may be nil.

### Components

`h.Builder[T]` gives a component every standard builder method, and each
method returns the component's own type, so calls chain without losing it.
Through `ChildAcceptor`, a component can also accept its own typed values as
children (for example `ui.ButtonVariantOutline`). `h.Slot` implements Radix's
`asChild`. See `components/ui/button.go` for the pattern that future
components will follow.

## Security model

Every protection is on by default. The only way to opt out is a type or
method named `Unsafe*`, so a reviewer can run `grep Unsafe` to see every
bypass.

**XSS (libhtml)**
- Text and attribute values are always escaped, and attributes are always
  double-quoted. NUL characters become U+FFFD.
- URL attributes (`href`, `src`, `action`, `formaction`, `xlink:href`,
  `srcset`, `hx-get`, …) accept only relative URLs and `http`, `https`,
  `mailto` and `tel`. The check matches how browsers parse URLs: they ignore
  tabs, newlines and control characters, so `java\tscript:` is caught.
  Other values become `about:invalid#libhtml-filtered`. An image `src`
  may also be a raster `data:image/...`. SVG data URLs are refused.
- Attributes that run JavaScript (`on*`, Alpine `x-*`/`@*`/`:*`, Vue `v-*`,
  `hx-on*`, `hx-vals`, hyperscript `_`) are refused when set from a plain
  string. Set them with `On` or `JSAttr`, which take `UnsafeJS`. `srcdoc`
  needs `UnsafeAttr`.
- `<script>` accepts only `UnsafeJS` content, and `<style>` only
  `UnsafeCSS`. Content that would close the element early (`</script`,
  `<!--`) is refused.
- `style` values that contain `expression(`, `javascript:`, `-moz-binding`,
  `@import`, backslash escapes or unsafe `url(...)` are filtered.
- Tag and attribute names are validated. `<plaintext>` and `<xmp>` are
  refused. Text in comments cannot close the comment.
- Links with `target="_blank"` get `rel="noopener noreferrer"`.
- `<meta http-equiv=refresh>` URLs, and SVG `<animate>`/`<set>` elements
  that target `href`, are sanitized.
- A refused element is reported as a `*SecurityError` from `Render`.
  `WriteHTTP` renders into a buffer first, so a failure never sends half a
  page.

**CSRF (`libhtml/csrf`)** uses two independent layers:
1. Go's `http.CrossOriginProtection`, which checks `Sec-Fetch-Site` and
   `Origin` against `Host`.
2. A random secret in a `__Host-` cookie that is HMAC-signed, `HttpOnly`,
   `Secure` and `SameSite=Lax`. Each form echoes it back through
   `csrf.Field(r)`, or scripts send it in the `X-CSRF-Token` header. The
   echoed token is masked with a new one-time pad on each render, as a
   BREACH defense, and compared in constant time. Tokens can optionally be
   bound to a session.

**Headers (`libhtml/secure`)** set a strict CSP with a new nonce per
request (`script-src 'nonce-…' 'strict-dynamic'`, `object-src 'none'`,
`base-uri 'none'`, `frame-ancestors 'none'`, `form-action 'self'`). They
also set `nosniff`, `X-Frame-Options`, `Referrer-Policy`, COOP and optional
HSTS. Use `secure.Script(r, …)` to emit a script that carries the nonce.
Under this CSP the browser also blocks inline `on*` handlers. Attach
listeners from a nonced script instead, as the demo does.

## Fidelity of class merging

`components/ui/button/testdata/button_ref.json` holds
`cn(buttonVariants({variant, size, className}))` for all 6 variants × 8
sizes × 6 `className` overrides, computed with the real
`class-variance-authority` and `tailwind-merge`. The Go output matches all
288 cases exactly. To regenerate the data:

```sh
cd scripts/jsref && npm install
node button.mjs ../../ui/apps/v4/registry/new-york-v4/ui/button.tsx > ../../components/ui/button/testdata/button_ref.json
```

## Tailwind, with Go only

`shadcn.go/tailwind` runs Tailwind's own compiler, the vendored `tailwindcss`
4.3.0 package, unmodified, inside goja, a JavaScript engine written in Go.
The CSS generated for any class is therefore exactly what Tailwind
generates. The surrounding steps are Go:

- source discovery: `@source`, `source(...)`, `.gitignore`, and automatic
  detection;
- class-candidate extraction, which reproduces the official Rust scanner's
  rules;
- an optional minifier.

`tw-animate-css` and `shadcn/tailwind.css` are embedded, so
`@import "tw-animate-css"` works without npm.

```go
//go:generate go run shadcn.go/cmd/tailwind -i styles/app.css -o static/app.css
```

The output is **byte-identical to the official CLI**
(`@tailwindcss/cli -i in.css -o out.css`, unminified). That holds for the
golden fixture in `tailwind/testdata`, for the demo, and for a stress build
over the 1,856 files of the upstream shadcn app (350 KB of CSS). Reaching
that required fixing two bugs in goja, both worked around in the prelude:
`Number(bigint)` truncates to 64 bits, and `Array.prototype.sort` treats `-0`
as "less than". Both changed Tailwind's rule order. `--minify` is a
conservative whitespace and comment minifier, not Lightning CSS, so it does
not add vendor prefixes or lower CSS nesting. The demo has a test that fails
when its committed stylesheet is stale.

## Parity with upstream

Every component is held to exact parity with shadcn/ui `new-york-v4`, and
`parity/` proves it:

- `parity/reference/build.mjs` bundles all 349 upstream references:
  examples, charts, blocks and internal demos. It resolves them against
  upstream's own lockfile-pinned dependencies (pnpm, frozen lockfile).
- `cmd/parity` serves `/react/<name>` and `/go/<name>` in an identical shell,
  with one stylesheet built by the Go Tailwind from both codebases. The Go
  pages run under the production CSP.
- `parity/run.mjs` drives both sides in headless Chromium, in light and dark
  mode. It compares full-page pixels (0 differences allowed) and the
  normalized DOM. It then replays interactions: automatic hover, focus, open
  and Escape, plus the per-reference scenarios in `parity/scenarios/`. For
  animations, it freezes the document timeline through the DevTools
  protocol, compares the animation definitions, and compares frames at 0,
  25, 50 and 75% of each animation's duration.

```sh
node parity/run.mjs --component button     # 13/13 pixel-, DOM- and animation-identical
node parity/status.mjs --todo
```

The harness is verification tooling only. It needs Node.js and the upstream
clone in `./ui`; applications built with shadcn.go need neither.

## Benchmarks

`cd bench && go test -bench . -benchmem` builds and renders a full page: a
document with a table of N rows, two links per row, escaped user data and
a conditional class. The table shows the median of 3 runs on an Apple
M-series machine with Go 1.27.

| Build + render | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: |
| libhtml, 10 rows | **9.2 µs** | 15.1 KB | **288** |
| html/template, 10 rows | 31.3 µs | **10.1 KB** | 454 |
| gomponents, 10 rows | 14.7 µs | 14.1 KB | 694 |
| libhtml, 1000 rows | **1.01 ms** | 1.34 MB | **25.8 k** |
| html/template, 1000 rows | 2.80 ms | **0.99 MB** | 46.7 k |
| gomponents, 1000 rows | 1.23 ms | 1.26 MB | 60.9 k |

| Render a prebuilt tree, 1000 rows | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: |
| libhtml | **361 µs** (~470 MB/s) | **0.3 KB** | **4** |
| gomponents | 780 µs (~185 MB/s) | 390 KB | 29.4 k |

| button.Button | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: |
| plain | 1.4 µs | 656 | 10 |
| with a `className` override (`rounded-full`) | 1.8 µs | 801 | 14 |
| `AsChild` link | 3.2 µs | 1.6 KB | 25 |

On these runs, libhtml was the fastest of the three and made the fewest
allocations. Its renderer copies text without allocating and reuses pooled
buffers. One cost remains: building a tree allocates about 6% more bytes than
gomponents and about 35% more than html/template, which parses the template
once and builds no tree. Most of that memory is short-lived garbage from tree
construction.

## Tests

```sh
go test ./...                       # libhtml, csrf, secure, tailwind, twmerge, cva, components
go generate ./...                   # stylesheets, example index
cd bench && go test ./...           # checks that all three renderers produce the same markup
```
