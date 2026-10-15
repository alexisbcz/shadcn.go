// Generates docs/site/data/create/config.json: the data of the /create
// designer (registry/config.ts, registry/styles.tsx, bases, base colors,
// themes, fonts, presets) and the icons of its IconPlaceholders (the create
// components and the preview/preview-02 cards), rendered for each icon
// library as the upstream renderers do.
//
// Usage (from the repo root):
//   ui/node_modules/.bin/tsx --tsconfig docs/site/gen/create.tsconfig.json docs/site/gen/create.tsx
import "./create-react-global"
import fs from "node:fs"
import path from "node:path"
import * as React from "react"
import { renderToStaticMarkup } from "react-dom/server"
import * as Lucide from "lucide-react"
import * as Tabler from "@tabler/icons-react"
import * as Huge from "@hugeicons/core-free-icons"
import { HugeiconsIcon } from "@hugeicons/react"
import * as Phosphor from "@phosphor-icons/react/dist/ssr"
import * as Remix from "@remixicon/react"
import QRCode from "react-qr-code"

import {
  BASE_COLORS,
  BASES,
  DEFAULT_CONFIG,
  iconLibraries,
  MENU_ACCENTS,
  MENU_COLORS,
  PRESETS,
  RADII,
  STYLES,
  THEMES,
} from "@/registry/config"
import { FONT_DEFINITIONS } from "@/lib/font-definitions"
import { TEMPLATES } from "@/app/(app)/(create)/lib/templates"

const root = path.resolve(__dirname, "../../..")
const app = path.join(root, "ui/apps/v4")

const markup = (el: React.ReactElement) => renderToStaticMarkup(el)

// IconPlaceholder usages.
const libs = ["lucide", "tabler", "hugeicons", "phosphor", "remixicon"] as const
const sources = [
  path.join(app, "app/(app)/(create)/components"),
  path.join(app, "registry/bases/radix/blocks/preview/cards"),
  path.join(app, "registry/bases/radix/blocks/preview-02/cards"),
]
const names: Record<string, Set<string>> = Object.fromEntries(libs.map((l) => [l, new Set<string>()]))
for (const dir of sources) {
  for (const f of fs.readdirSync(dir)) {
    if (!f.endsWith(".tsx")) continue
    const src = fs.readFileSync(path.join(dir, f), "utf8")
    for (const [tag] of src.matchAll(/<IconPlaceholder\b[\s\S]*?\/>/g)) {
      for (const l of libs) {
        const m = tag.match(new RegExp(`\\b${l}="([^"]+)"`))
        if (m) names[l].add(m[1])
      }
    }
    // Icon names in data ({ lucide: "CopyIcon", tabler: ... }).
    for (const l of libs) {
      for (const [, n] of src.matchAll(new RegExp(`\\b${l}: "([^"]+)"`, "g"))) names[l].add(n)
    }
  }
}
// The icon renderers get className "__CLS__"; the Go side substitutes it.
const render: Record<string, (n: string) => string | null> = {
  lucide: (n) => ((Lucide as any)[n] ? markup(React.createElement((Lucide as any)[n], { className: "__CLS__" })) : null),
  tabler: (n) => ((Tabler as any)[n] ? markup(React.createElement((Tabler as any)[n], { className: "__CLS__" })) : null),
  hugeicons: (n) => ((Huge as any)[n] ? markup(React.createElement(HugeiconsIcon, { icon: (Huge as any)[n], strokeWidth: 2, className: "__CLS__" })) : null),
  phosphor: (n) => ((Phosphor as any)[n] ? markup(React.createElement((Phosphor as any)[n], { className: "__CLS__" })) : null),
  remixicon: (n) => ((Remix as any)[n] ? markup(React.createElement((Remix as any)[n], { className: "__CLS__" })) : null),
}
const icons: Record<string, Record<string, string>> = {}
for (const l of libs) {
  icons[l] = {}
  for (const n of [...names[l]].sort()) {
    const svg = render[l](n)
    if (svg) icons[l][n] = svg
    else console.warn(`missing ${l} ${n}`)
  }
}

// CSS variables as ordered [name, value] pairs (the order of the upstream
// objects, which the generated theme CSS follows).
const pairs = (v: any) => Object.fromEntries(Object.entries(v ?? {}).map(([k, o]) => [k, Object.entries(o as any)]))

// The hugeicons the create components import (HugeiconsIcon, strokeWidth 2).
const hugeNames = new Set<string>()
for (const f of fs.readdirSync(sources[0])) {
  const src = fs.readFileSync(path.join(sources[0], f), "utf8")
  const m = src.match(/import \{([^}]*)\} from "@hugeicons\/core-free-icons"/)
  if (m) for (const n of m[1].split(",").map((x) => x.trim()).filter(Boolean)) hugeNames.add(n)
}
const createHugeicons = Object.fromEntries([...hugeNames].sort().map((n) => [n, markup(React.createElement(HugeiconsIcon, { icon: (Huge as any)[n], strokeWidth: 2, className: "__CLS__" }))]))

// icon-library-picker.tsx's logos: the object literal, compiled with the
// classic JSX runtime from a scratch module.
const pickerSrc = fs.readFileSync(path.join(sources[0], "icon-library-picker.tsx"), "utf8")
const logosSrc = pickerSrc.slice(pickerSrc.indexOf("const logos = {"), pickerSrc.indexOf("\n}\n", pickerSrc.indexOf("const logos = {")) + 2)
const tmp = path.join(__dirname, ".create-logos.tsx")
fs.writeFileSync(tmp, `/** @jsxRuntime classic */\n/** @jsx React.createElement */\n${logosSrc}\nexport { logos }\n`)
const { logos } = require(tmp)
fs.rmSync(tmp)
const iconLibraryLogos = Object.fromEntries(Object.entries(logos).map(([k, v]) => [k, markup(v as any)]))

const out = {
  default: DEFAULT_CONFIG,
  styles: STYLES.map((s: any) => ({ name: s.name, title: s.title, description: s.description, icon: s.icon ? markup(React.cloneElement(s.icon, { className: "size-4" })) : "" })),
  bases: BASES.map((b: any) => ({ name: b.name, title: b.title, description: b.description, logo: b.meta?.logo ?? "" })),
  baseColors: BASE_COLORS.map((c: any) => ({ name: c.name, title: c.title, cssVars: pairs(c.cssVars) })),
  themes: THEMES.map((t: any) => ({ name: t.name, title: t.title, cssVars: pairs(t.cssVars) })),
  iconLibraries: Object.values(iconLibraries).map((l: any) => ({ name: l.name, title: l.title })),
  menuAccents: MENU_ACCENTS,
  menuColors: MENU_COLORS,
  radii: RADII,
  fonts: FONT_DEFINITIONS,
  presets: PRESETS,
  icons,
  hugeicons: createHugeicons,
  iconLibraryLogos,
  // project-form.tsx
  templates: TEMPLATES,
  turborepoLogo: fs.readFileSync(path.join(sources[0], "project-form.tsx"), "utf8").match(/const TURBOREPO_LOGO =\s*'([^']*)'/)?.[1] ?? "",
  // qr-connect.tsx
  qrConnect: markup(React.createElement(QRCode as any, { value: "https://ledger.app/connect/jd-4829", size: 160, level: "M" })),
}
fs.mkdirSync(path.join(root, "docs/site/data/create"), { recursive: true })
fs.writeFileSync(path.join(root, "docs/site/data/create/config.json"), JSON.stringify(out, null, 1) + "\n")
console.log(`config.json: ${STYLES.length} styles, ${THEMES.length} themes, ${Object.values(icons).reduce((n, m) => n + Object.keys(m).length, 0)} icons`)
