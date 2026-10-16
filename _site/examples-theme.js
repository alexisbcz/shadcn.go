// Examples theme selector (components/theme-selector.tsx, theme-customizer.tsx
// CopyCodeButton/CustomizerCode, active-theme.tsx): the selected theme sets
// body's theme-<name> class and re-renders the copy-code dialog's CSS. The
// active theme lives in React state that survives client-side navigations;
// sessionStorage plays that part across the examples' full page loads.
import { getSelect } from "/shadcn.go/_shadcn/select.js"

const root = document.querySelector("[data-theme-selector]")
const themes = JSON.parse(document.querySelector("script[data-themes]").textContent)
const KEY = "examples-active-theme"
const nbsp = " "

function el(tag, attrs, ...children) {
  const e = document.createElement(tag)
  for (const [k, v] of Object.entries(attrs)) e.setAttribute(k, v)
  e.append(...children)
  return e
}
const ind = (color) => el("span", { class: "inline-block size-3 border border-border/50", style: `background-color:${color}` })

// lines mirrors themeCodeLines (examples_theme.go).
function lines(t, format) {
  const out = []
  const line = (cls, ...c) => out.push(el("span", { "data-line": "", class: cls }, ...c))
  const fg = "line text-code-foreground"
  if (format === "v3") {
    const v = (vars, k) => {
      const val = (vars.find((e) => e[0] === k) ?? [])[1] ?? ""
      line("line", nbsp.repeat(4) + "--" + k + ": ", ind(`hsl(${val})`), " " + val + ";")
    }
    const block = (sel, vars, radius) => {
      line("line", nbsp + nbsp + sel + " {")
      v(vars, "background")
      v(vars, "foreground")
      for (const p of ["card", "popover", "primary", "secondary", "muted", "accent", "destructive"]) {
        v(vars, p)
        v(vars, p + "-foreground")
      }
      for (const k of ["border", "input", "ring"]) v(vars, k)
      if (radius) line("line", nbsp.repeat(4) + "--radius: 0.5rem;")
      for (let i = 1; i <= 5; i++) v(vars, "chart-" + i)
      line("line", nbsp + nbsp + "}")
    }
    line("line", "@layer base {")
    block(":root", t.hsl.light, true)
    line("line", nbsp)
    block(".dark", t.hsl.dark, false)
    line("line", "}")
    return out
  }
  const oklch = format === "v4-oklch" && t.oklch
  const vars = oklch ? t.oklch : t.hsl
  const entry = ([k, v]) =>
    oklch
      ? line(fg, nbsp.repeat(3) + "--" + k + ": ", ind(v), " " + v + ";")
      : line(fg, nbsp.repeat(3) + "--" + k + ": ", ind(`hsl(${v})`), ` hsl(${v});`)
  line(fg, nbsp + ":root {")
  line(fg, nbsp.repeat(3) + "--radius: 0.65rem;")
  vars.light.forEach(entry)
  line(fg, nbsp + "}")
  line(fg, nbsp)
  line(fg, nbsp + ".dark {")
  vars.dark.forEach(entry)
  line(fg, nbsp + "}")
  return out
}

// getThemeCodeOKLCH, getThemeCodeHSLV4 and getThemeCode (the v3 template).
function code(t, format) {
  if (format === "v3") {
    const keys = ["background", "foreground", "card", "card-foreground", "popover", "popover-foreground", "primary", "primary-foreground", "secondary", "secondary-foreground", "muted", "muted-foreground", "accent", "accent-foreground", "destructive", "destructive-foreground", "border", "input", "ring"]
    const charts = [1, 2, 3, 4, 5].map((i) => "chart-" + i)
    const get = (vars, k) => (vars.find((e) => e[0] === k) ?? [])[1]
    const sect = (vars, radius) =>
      [...keys, ...(radius ? ["radius"] : []), ...charts]
        .map((k) => `    --${k}: ${k === "radius" ? "0.5rem" : get(vars, k)};`)
        .join("\n")
    return `\n@layer base {\n  :root {\n${sect(t.hsl.light, true)}\n  }\n\n  .dark {\n${sect(t.hsl.dark, false)}\n  }\n}\n`
  }
  const oklch = format === "v4-oklch"
  const vars = oklch ? t.oklch : t.hsl
  if (!vars) return ""
  const fmt = (e) => "  --" + e[0] + ": " + (oklch ? e[1] : "hsl(" + e[1] + ")") + ";"
  return ":root {\n  --radius: 0.65rem;\n" + vars.light.map(fmt).join("\n") + "\n}\n\n.dark {\n" + vars.dark.map(fmt).join("\n") + "\n}\n"
}

// all queries the document and the template contents (closed dialogs and
// drawers keep their content in templates until opened).
function all(sel, root = document) {
  const out = [...root.querySelectorAll(sel)]
  for (const t of root.querySelectorAll("template")) out.push(...all(sel, t.content))
  return out
}

let active = "default"
const themeOf = () => themes.find((t) => t.name === (active === "default" ? "neutral" : active)) ?? themes[0]

function apply(name) {
  active = name
  for (const c of [...document.body.classList]) if (c.startsWith("theme-")) document.body.classList.remove(c)
  document.body.classList.add(`theme-${name}`)
  if (name.endsWith("-scaled")) document.body.classList.add("theme-scaled")
  const t = themeOf()
  for (const s of all("[data-theme-name]")) s.textContent = t.name
  for (const c of all("[data-theme-code]")) c.replaceChildren(...lines(t, c.dataset.themeCode))
  try { sessionStorage.setItem(KEY, name) } catch {}
}

root?.addEventListener("shadcn:value-change", (e) => e.detail.value && apply(e.detail.value))

let stored = null
try { stored = sessionStorage.getItem(KEY) } catch {}
if (stored && stored !== "default" && document.referrer.startsWith(location.origin + "/shadcn.go/examples")) {
  getSelect(root.querySelector("#theme-selector"))?.setValue(stored)
  apply(stored)
}

const copyIcon = '<path d="M7 7m0 2.667a2.667 2.667 0 0 1 2.667 -2.667h8.666a2.667 2.667 0 0 1 2.667 2.667v8.666a2.667 2.667 0 0 1 -2.667 2.667h-8.666a2.667 2.667 0 0 1 -2.667 -2.667z"></path><path d="M4.012 16.737a2.005 2.005 0 0 1 -1.012 -1.737v-10c0 -1.1 .9 -2 2 -2h10c.75 0 1.158 .385 1.5 1"></path>'
document.addEventListener("click", (e) => {
  const b = e.target instanceof Element ? e.target.closest("[data-theme-copy]") : null
  if (!b) return
  navigator.clipboard?.writeText(code(themeOf(), b.dataset.themeCopy))
  const svg = b.querySelector("svg")
  svg.innerHTML = '<path d="M5 12l5 5l10 -10"></path>'
  svg.setAttribute("class", "tabler-icon tabler-icon-check ")
  clearTimeout(b._t)
  b._t = setTimeout(() => {
    svg.innerHTML = copyIcon
    svg.setAttribute("class", "tabler-icon tabler-icon-copy ")
  }, 2000)
})
