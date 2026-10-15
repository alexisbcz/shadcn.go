// The create designer (app/(app)/(create)): the React state of the page.
//
// - the design system params, read from and written to the URL as
//   useDesignSystemSearchParams does (a preset code for the design system
//   keys, normalizeDesignSystemParams);
// - the pickers (components/picker.tsx on Base UI's Menu): open on click,
//   radio items commit, hovering an item previews it in the iframe
//   (PreviewOverrideProvider, 50ms trailing);
// - the preview iframe (components/preview.tsx): new src for a base or item
//   change, otherwise the params are posted to it (create-preview.js);
// - locks, shuffle (use-random.tsx), undo/redo (use-history.tsx), the main
//   menu, the keyboard shortcuts, copy preset, the welcome dialog.
// The dialogs (Get Code, Open Preset, the action menu, reset) are in
// create-dialogs.js.
import { define } from "./core.js"
import { getDialog } from "./dialog.js"

// --- preset codes (packages/registry/src/preset/preset.ts) ---

const THEMES = ["neutral", "stone", "zinc", "gray", "amber", "blue", "cyan", "emerald", "fuchsia", "green", "indigo", "lime", "orange", "pink", "purple", "red", "rose", "sky", "teal", "violet", "yellow", "mauve", "olive", "mist", "taupe"]
const FONTS = ["inter", "noto-sans", "nunito-sans", "figtree", "roboto", "raleway", "dm-sans", "public-sans", "outfit", "jetbrains-mono", "geist", "geist-mono", "lora", "merriweather", "playfair-display", "noto-serif", "roboto-slab", "oxanium", "manrope", "space-grotesk", "montserrat", "ibm-plex-sans", "source-sans-3", "instrument-sans", "eb-garamond", "instrument-serif"]
const FIELDS = [
  ["menuColor", ["default", "inverted", "default-translucent", "inverted-translucent"], 3],
  ["menuAccent", ["subtle", "bold"], 3],
  ["radius", ["default", "none", "small", "medium", "large"], 4],
  ["font", FONTS, 6],
  ["iconLibrary", ["lucide", "hugeicons", "tabler", "phosphor", "remixicon"], 6],
  ["theme", THEMES, 6],
  ["baseColor", ["neutral", "stone", "zinc", "gray", "mauve", "olive", "mist", "taupe"], 6],
  ["style", ["nova", "vega", "maia", "lyra", "mira", "luma", "sera", "rhea"], 6],
  ["chartColor", THEMES, 6],
  ["fontHeading", ["inherit", ...FONTS], 5],
]
const BASE62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const V1_CHART_COLOR_MAP = { neutral: "blue", stone: "lime", zinc: "amber", mauve: "emerald", olive: "violet", mist: "rose", taupe: "cyan" }

export function encodePreset(config) {
  let bits = 0
  let offset = 0
  for (const [key, values, n] of FIELDS) {
    const idx = values.indexOf(config[key])
    bits += (idx === -1 ? 0 : idx) * 2 ** offset
    offset += n
  }
  if (bits === 0) return "b0"
  let s = ""
  while (bits > 0) {
    s = BASE62[bits % 62] + s
    bits = Math.floor(bits / 62)
  }
  return "b" + s
}

export function isPresetCode(v) {
  if (!v || v.length < 2 || v.length > 10 || (v[0] !== "a" && v[0] !== "b")) return false
  for (let i = 1; i < v.length; i++) if (BASE62.indexOf(v[i]) === -1) return false
  return true
}

export function decodePreset(code) {
  if (!isPresetCode(code)) return null
  const fields = code[0] === "a" ? FIELDS.slice(0, 8) : FIELDS
  let bits = 0
  for (let i = 1; i < code.length; i++) bits = bits * 62 + BASE62.indexOf(code[i])
  const out = {}
  let offset = 0
  for (const [key, values, n] of fields) {
    const idx = Math.floor(bits / 2 ** offset) % 2 ** n
    out[key] = idx < values.length ? values[idx] : values[0]
    offset += n
  }
  if (code[0] === "a") out.fontHeading = "inherit"
  return out
}

// --- params (lib/search-params.ts) ---

export const DESIGN_KEYS = ["style", "baseColor", "theme", "chartColor", "iconLibrary", "font", "fontHeading", "radius", "menuAccent", "menuColor"]
const DEFAULTS = {
  preset: "b0", base: "base", item: "preview-02", iconLibrary: "lucide", style: "nova", theme: "neutral", chartColor: "neutral",
  font: "inter", fontHeading: "inherit", baseColor: "neutral", menuAccent: "subtle", menuColor: "default", radius: "default",
  template: "next", rtl: false, pointer: false, size: 100, custom: false,
}

let config // the page's data-create-config

const isBaseColor = (name) => config.baseColors.some((c) => c.name === name)
export const themesForBaseColor = (baseColor) => config.themes.filter((t) => t.name === baseColor || !isBaseColor(t.name))
const isTranslucent = (menuColor) => menuColor === "default-translucent" || menuColor === "inverted-translucent"

export function normalize(p) {
  const r = { ...p }
  if (r.fontHeading === r.font) r.fontHeading = "inherit"
  const available = themesForBaseColor(r.baseColor)
  const ok = (n) => available.some((t) => t.name === n)
  if (!ok(r.theme) || !ok(r.chartColor)) {
    const fallback = available[0]?.name ?? r.baseColor
    if (!ok(r.theme)) r.theme = fallback
    if (!ok(r.chartColor)) r.chartColor = fallback
  }
  if (r.menuAccent === "bold" && isTranslucent(r.menuColor)) r.menuAccent = "subtle"
  return r
}

// readURL is resolvePresetParams over the location's search params.
function readURL() {
  const q = new URLSearchParams(location.search)
  const raw = { ...DEFAULTS }
  for (const k of Object.keys(DEFAULTS)) {
    const v = q.get(k)
    if (v == null) continue
    raw[k] = typeof DEFAULTS[k] === "boolean" ? v === "true" : typeof DEFAULTS[k] === "number" ? Number(v) : v
  }
  const decoded = raw.preset && decodePreset(raw.preset)
  if (decoded) {
    return normalize({
      ...raw,
      ...decoded,
      fontHeading: q.has("fontHeading") ? q.get("fontHeading") : decoded.fontHeading,
      chartColor: q.has("chartColor") ? q.get("chartColor") : (decoded.chartColor ?? V1_CHART_COLOR_MAP[decoded.theme] ?? decoded.theme),
    })
  }
  return normalize(raw)
}

// writeURL is buildPresetUrlUpdate through nuqs (clearOnDefault).
function urlFor(p) {
  const q = new URLSearchParams(location.search)
  for (const k of DESIGN_KEYS) q.delete(k)
  const code = encodePreset(p)
  const set = (k, v) => (v === DEFAULTS[k] || v == null ? q.delete(k) : q.set(k, String(v)))
  set("preset", code)
  for (const k of ["base", "item", "template", "rtl", "pointer", "size", "custom"]) set(k, p[k])
  const s = q.toString()
  return location.pathname + (s ? "?" + s : "") + location.hash
}

// --- theme (DesignSystemProvider, buildRegistryTheme) ---

// themeVars is buildRegistryTheme's light and dark variables.
export function themeVars(p) {
  const bc = config.baseColors.find((c) => c.name === p.baseColor)
  const th = config.themes.find((t) => t.name === p.theme)
  if (!bc || !th) return null
  const radius = p.style === "lyra" || p.style === "sera" ? "none" : p.style === "rhea" && p.radius === "large" ? "default" : p.radius
  const merge = (mode) => {
    const out = new Map()
    for (const [k, v] of bc.cssVars[mode] ?? []) out.set(k, v)
    for (const [k, v] of th.cssVars[mode] ?? []) out.set(k, v)
    const chart = config.themes.find((t) => t.name === p.chartColor)
    if (chart) {
      for (let i = 1; i <= 5; i++) {
        const v = (chart.cssVars[mode] ?? []).find(([k]) => k === `chart-${i}`)?.[1]
        if (v) out.set(`chart-${i}`, v)
      }
    }
    if (p.menuAccent === "bold") {
      out.set("accent", out.get("primary"))
      out.set("accent-foreground", out.get("primary-foreground"))
    }
    if (mode === "light" && radius !== "default") {
      const r = config.radii.find((x) => x.name === radius)
      if (r?.value) out.set("radius", r.value)
    }
    return out
  }
  return { light: merge("light"), dark: merge("dark"), radius }
}

// presetThemeCSS is project-form.tsx's formatThemeCss(buildThemeForPreset).
export function presetThemeCSS(p) {
  const v = themeVars(p)
  if (!v) return ""
  const r = config.radii.find((x) => x.name === v.radius)
  v.light.set("radius", v.radius === "default" ? (v.light.get("radius") ?? "0.625rem") : (r?.value ?? v.light.get("radius")))
  const rule = (sel, vars) => `${sel} {\n${[...vars].map(([k, val]) => `  --${k}: ${val};`).join("\n")}\n}`
  return [rule(":root", v.light), rule(".dark", v.dark)].join("\n\n")
}

export function themeCSS(p) {
  const v = themeVars(p)
  if (!v) return ""
  const merge = (mode) => v[mode]
  const rule = (sel, vars) => {
    const lines = [...vars].filter(([, v]) => v).map(([k, v]) => `  --${k}: ${v};`)
    return lines.length ? `${sel} {\n${lines.join("\n")}\n}\n` : `${sel} {}\n`
  }
  let css = rule(":root", merge("light")) + "\n" + rule(".dark", merge("dark"))
  if (p.pointer) css += '\n@layer base {\n  button:not(:disabled), [role="button"]:not(:disabled) {\n    cursor: pointer;\n  }\n}\n'
  return css
}

export function fontFamilies(p) {
  const sans = config.fonts[p.font]?.family ?? ""
  const heading = p.fontHeading !== "inherit" && p.fontHeading !== p.font ? (config.fonts[p.fontHeading]?.family ?? sans) : sans
  return { sans, heading }
}

// --- the designer ---

const LOCKABLE = DESIGN_KEYS
const isMobile = () => matchMedia("(max-width: 767px)").matches
const editable = (t) => t instanceof HTMLElement && (t.isContentEditable || t instanceof HTMLInputElement || t instanceof HTMLTextAreaElement || t instanceof HTMLSelectElement)

define("create", {
  selector: "[data-create]",
  mount(root) {
    config = JSON.parse(root.querySelector("[data-create-config]").textContent)
    let state = normalize({ ...DEFAULTS, ...JSON.parse(root.getAttribute("data-create")) })
    let override = null
    const locks = new Set()
    const iframe = root.querySelector("[data-create-iframe]")
    let iframeKey = state.base + state.item
    let sent = null

    // History (use-history.tsx): the preset codes, "" for the first entry.
    const initialPreset = new URLSearchParams(location.search).get("preset") ?? ""
    const history = { entries: [initialPreset], index: 0, max: 0 }

    // --- iframe sync ---
    const previewSrc = (p) => {
      const q = new URLSearchParams()
      const code = encodePreset(p)
      if (code !== "b0") q.set("preset", code)
      if (p.pointer) q.set("pointer", "true")
      const s = q.toString()
      return `/preview/${p.base}/${p.item}${s ? "?" + s : ""}`
    }
    const message = (p) => ({ type: "design-system-params", params: p, css: themeCSS(p), fonts: fontFamilies(p), fontTitles: { body: config.fonts[p.font]?.title, heading: config.fonts[p.fontHeading]?.title }, styleTitle: config.styles.find((s) => s.name === p.style)?.title })
    const sync = () => {
      const p = override ? normalize({ ...state, ...override }) : state
      if (state.base + state.item !== iframeKey) {
        iframeKey = state.base + state.item
        iframe.src = previewSrc(state)
        sent = null
        return
      }
      const key = JSON.stringify(p)
      if (key === sent || !iframe.contentWindow) return
      sent = key
      iframe.contentWindow.postMessage(message(p), location.origin)
    }
    iframe.addEventListener("load", () => {
      sent = null
      sync()
    })

    // --- UI ---
    const pickers = [...root.querySelectorAll("[data-create-picker]")]
    const pickerValue = (param) => {
      if (param === "menuColor") return state.menuColor
      if (param === "radius" && (state.style === "lyra" || state.style === "sera")) return "none"
      return state[param]
    }
    const titleOf = (param, value) => {
      switch (param) {
        case "style":
          return config.styles.find((s) => s.name === value)?.title
        case "baseColor":
          return config.baseColors.find((c) => c.name === value)?.title
        case "theme":
        case "chartColor":
          return config.themes.find((t) => t.name === value)?.title
        case "font":
          return config.fonts[value]?.title
        case "fontHeading":
          return value === "inherit" ? config.fonts[state.font]?.title : config.fonts[value]?.title
        case "iconLibrary":
          return config.iconLibraries.find((l) => l.name === value)?.title
        case "radius":
          return config.radii.find((r) => r.name === value)?.label
        case "menuColor":
          return config.menuColorLabels[value]
        case "menuAccent":
          return value === "bold" ? "Bold" : "Subtle"
        case "base":
          return config.bases.find((b) => b.name === value)?.title
      }
    }
    const render = () => {
      for (const el of pickers) {
        const param = el.getAttribute("data-create-picker")
        const value = pickerValue(param)
        el.querySelector("[data-create-picker-value]").textContent = titleOf(param, value) ?? ""
        const ind = el.querySelector(`template[data-create-indicator="${CSS.escape(param === "fontHeading" && value === "inherit" ? state.font : value)}"]`)
        if (ind) el.querySelector("[data-create-picker-indicator]").replaceChildren(ind.content.cloneNode(true))
        const trigger = el.querySelector("[data-create-menu-trigger]")
        if (param === "radius") trigger.disabled = state.style === "lyra" || state.style === "sera"
        const lock = el.querySelector("[data-create-lock]")
        if (lock) {
          const on = locks.has(param)
          lock.dataset.locked = String(on)
          lock.title = on ? "Unlock" : "Lock"
          lock.setAttribute("aria-label", on ? "Unlock" : "Lock")
          const tpl = root.querySelector(`template[data-create-lock-icon="${on ? "locked" : "unlocked"}"]`)
          if (tpl) lock.replaceChildren(tpl.content.cloneNode(true))
        }
      }
      const code = encodePreset(state)
      for (const el of document.querySelectorAll("[data-create-preset-label]")) {
        if (el.dataset.copied) continue
        el.textContent = `--preset ${code}`
        el.closest("button")?.setAttribute("title", `--preset ${code}`)
      }
      for (const a of document.querySelectorAll("[data-create-v0]")) {
        a.setAttribute("href", `undefined/chat/api/open?url=undefined/init/v0?preset=${code}&base=${state.base}&title=New ${state.base}-${state.style} project`)
      }
      const switcher = root.querySelector("[data-create-preview-switcher]")
      if (switcher) {
        switcher.hidden = !(state.item === "preview" || state.item.startsWith("preview-0"))
        for (const b of switcher.querySelectorAll("[data-create-item]")) b.dataset.active = String(b.dataset.createItem === state.item)
      }
      document.dispatchEvent(new CustomEvent("create:state", { detail: { state, history } }))
      if (open) syncMenu(open)
    }

    const commit = (updates, { replace = false, fromHistory = false } = {}) => {
      if (updates.menuAccent === "bold" && isTranslucent(updates.menuColor)) updates = { ...updates, menuAccent: "subtle" }
      const next = normalize({ ...state, ...updates })
      const changed = Object.keys(next).some((k) => next[k] !== state[k])
      if (!changed) return
      state = next
      const url = urlFor(state)
      if (url !== location.pathname + location.search + location.hash) {
        if (replace) window.history.replaceState(window.history.state, "", url)
        else window.history.pushState(window.history.state, "", url)
      }
      // use-history.tsx: a new preset becomes the next entry.
      const preset = new URLSearchParams(location.search).get("preset") ?? ""
      if (!fromHistory && preset !== history.entries[history.index]) {
        history.entries = history.entries.slice(0, history.index + 1)
        history.entries.push(preset)
        history.index = history.max = history.entries.length - 1
      }
      render()
      sync()
    }
    const goHistory = (delta) => {
      const next = history.index + delta
      if (next < 0 || next > history.max) return
      history.index = next
      const preset = history.entries[next]
      const decoded = preset ? decodePreset(preset) : decodePreset("b0")
      if (!decoded) return
      decoded.chartColor = decoded.chartColor ?? V1_CHART_COLOR_MAP[decoded.theme] ?? decoded.theme
      commit(decoded, { replace: true, fromHistory: true })
    }
    window.addEventListener("popstate", () => {
      state = readURL()
      render()
      sync()
    })

    // --- shuffle (use-random.tsx: a curated preset, locks kept) ---
    const SHUFFLE = ["b6sUj34d9", "b2tqYzpa88", "b1W4tDrk", "b1aIuQ2XC", "b7jsW1RxJ5", "b870VEw0in", "b3Zheoix4U", "b1x9M2c4aI", "b1W7jDEW", "b51GFh7y6", "b2fms620zo", "b1Q5GC", "buKEvLs", "b5rR41Mtnc", "b6tOz2I0x", "b2hNTREGRN", "bdIJ7Sq", "b6TqMNb5Wb", "bJIirQ", "b4aRK5K0fb", "b5HCiD38LI", "bdHjvCi", "b7QDHijUjj", "b4ZVZIPi9h", "b1W4bcno"]
    const shuffle = () => {
      const current = encodePreset(state)
      const codes = SHUFFLE.filter((c) => c !== current)
      const list = codes.length ? codes : SHUFFLE
      const decoded = decodePreset(list[Math.floor(Math.random() * list.length)])
      if (!decoded) return
      decoded.chartColor = decoded.chartColor ?? decoded.theme
      const next = {}
      for (const k of DESIGN_KEYS) next[k] = locks.has(k) ? state[k] : decoded[k]
      commit(next)
    }

    // --- menus (picker.tsx) ---
    let open = null // { name, trigger, portal, popup }
    const syncMenu = (m) => {
      const param = m.name
      for (const group of m.popup.querySelectorAll("[data-create-radio-group]")) {
        const g = group.getAttribute("data-create-radio-group")
        let value = pickerValue(g)
        if (g === "menuColor:color") value = state.menuColor.startsWith("inverted") ? "inverted" : "default"
        if (g === "menuColor:surface") value = isTranslucent(state.menuColor) ? "translucent" : "solid"
        for (const item of group.querySelectorAll("[role=menuitemradio]")) {
          const v = item.dataset.value
          const on = v === value
          item.setAttribute("aria-checked", String(on))
          item.toggleAttribute("data-checked", on)
          item.toggleAttribute("data-unchecked", !on)
          item.querySelector("[data-create-check]").hidden = !on
          // Availability that depends on the state.
          let disabled = item.hasAttribute("data-create-disabled")
          if (param === "menuAccent" && v === "bold") disabled = isTranslucent(state.menuColor)
          if (param === "radius" && v === "large") disabled = state.style === "rhea"
          if (g === "menuColor:color" && v === "inverted") disabled = document.documentElement.classList.contains("dark")
          item.toggleAttribute("data-disabled", disabled)
          if (disabled) item.setAttribute("aria-disabled", "true")
          else item.removeAttribute("aria-disabled")
          if (param === "theme" || param === "chartColor") item.hidden = isBaseColor(v) && v !== state.baseColor
          if (param === "fontHeading" && v === "inherit") item.lastChild.textContent = config.fonts[state.font]?.title ?? "Body font"
        }
      }
    }
    const position = (m) => {
      const pos = m.positioner
      const pad = 5
      const vw = document.documentElement.clientWidth
      const vh = window.innerHeight
      pos.style.setProperty("--available-width", `${vw - 2 * pad}px`)
      pos.style.setProperty("--available-height", `${vh - 2 * pad}px`)
      const t = m.trigger.getBoundingClientRect()
      const pw = m.popup.offsetWidth
      const ph = m.popup.offsetHeight
      let x, y, side, align
      if (isMobile() && m.name !== "main") {
        // anchor: the customizer, side top, align center.
        const a = root.querySelector("[data-create-customizer]").getBoundingClientRect()
        side = "top"
        align = "center"
        x = a.left + a.width / 2 - pw / 2
        y = a.top - 20 - ph
        x = Math.max(pad, Math.min(x, vw - pad - pw))
        pos.style.setProperty("--transform-origin", `${a.left + a.width / 2 - x}px ${ph + 20}px`)
      } else {
        side = "right"
        align = "start"
        const alignOffset = Number(m.tpl.getAttribute("data-align-offset") ?? 0)
        x = t.right + 20
        y = t.top + alignOffset
        y = Math.max(pad, Math.min(y, vh - pad - ph))
        pos.style.setProperty("--transform-origin", `-20px ${t.top + t.height / 2 - y}px`)
      }
      pos.style.transform = `translate(${Math.round(x + scrollX)}px, ${Math.round(y + scrollY)}px)`
      for (const el of [pos, m.popup]) {
        el.dataset.side = side
        el.dataset.align = align
      }
    }
    const close = () => {
      if (!open) return
      const m = open
      open = null
      m.trigger.setAttribute("aria-expanded", "false")
      m.trigger.removeAttribute("data-popup-open")
      m.popup.toggleAttribute("data-open", false)
      m.popup.toggleAttribute("data-closed", true)
      m.portal.remove()
      clearOverride()
    }
    const openMenu = (trigger) => {
      const name = trigger.getAttribute("data-create-menu-trigger")
      const tpl = trigger.parentElement.querySelector(`template[data-create-menu="${CSS.escape(name)}"]`)
      if (!tpl) return
      close()
      const portal = document.createElement("div")
      portal.setAttribute("data-create-portal", "")
      const positioner = document.createElement("div")
      positioner.className = "isolate z-50 outline-none"
      positioner.setAttribute("role", "presentation")
      positioner.style.cssText = "position: absolute; left: 0px; top: 0px;"
      positioner.append(tpl.content.cloneNode(true))
      const popup = positioner.firstElementChild
      popup.setAttribute("tabindex", "-1")
      popup.toggleAttribute("data-open", true)
      const backdrop = document.createElement("div")
      backdrop.className = "absolute inset-y-0 right-0 left-62 z-40 bg-transparent"
      portal.append(positioner, backdrop)
      document.body.append(portal)
      open = { name, trigger, tpl, portal, positioner, popup }
      syncMenu(open)
      position(open)
      trigger.setAttribute("aria-expanded", "true")
      trigger.toggleAttribute("data-popup-open", true)
      popup.focus({ preventScroll: true })
      popup.addEventListener("mouseleave", clearOverride)
    }

    // --- preview override (hover previews) ---
    let overrideTimer = 0
    const setOverride = (next) => {
      clearTimeout(overrideTimer)
      overrideTimer = setTimeout(() => {
        override = next
        sync()
      }, 50)
    }
    function clearOverride() {
      clearTimeout(overrideTimer)
      if (!override) return
      override = null
      sync()
    }
    const menuColorUpdate = (group, v) => {
      const color = group === "menuColor:color" ? v : state.menuColor.startsWith("inverted") ? "inverted" : "default"
      const translucent = group === "menuColor:surface" ? v === "translucent" : isTranslucent(state.menuColor)
      const menuColor = translucent ? `${color}-translucent` : color
      return translucent ? { menuColor, menuAccent: "subtle" } : { menuColor, menuAccent: lastSolidAccent }
    }
    let lastSolidAccent = state.menuAccent
    const updateFor = (item) => {
      const group = item.closest("[data-create-radio-group]").getAttribute("data-create-radio-group")
      const v = item.dataset.value
      if (group.startsWith("menuColor:")) return menuColorUpdate(group, v)
      return { [group]: v }
    }

    const actions = {
      navigate: () => document.dispatchEvent(new CustomEvent("create:open", { detail: "action-menu" })),
      "open-preset": () => document.dispatchEvent(new CustomEvent("create:open", { detail: "open-preset" })),
      shuffle,
      "toggle-theme": () => document.querySelector('[data-site-action="toggle-theme"]')?.click(),
      undo: () => goHistory(-1),
      redo: () => goHistory(1),
      reset: () => document.dispatchEvent(new CustomEvent("create:open", { detail: "reset" })),
      "copy-preset": (btn) => {
        const code = encodePreset(state)
        navigator.clipboard?.writeText(`--preset ${code}`)
        const label = btn.querySelector("[data-create-preset-label]")
        label.dataset.copied = "1"
        label.textContent = "Copied"
        btn.title = "Copied"
        clearTimeout(btn._t)
        btn._t = setTimeout(() => {
          delete label.dataset.copied
          render()
        }, 2000)
      },
    }

    const onClick = (e) => {
      const t = e.target instanceof Element ? e.target : null
      if (!t) return
      const trigger = t.closest("[data-create-menu-trigger]")
      if (trigger && root.contains(trigger)) {
        if (open?.trigger === trigger) close()
        else openMenu(trigger)
        return
      }
      if (open && open.portal.contains(t)) {
        const radio = t.closest("[role=menuitemradio]")
        if (radio && !radio.hasAttribute("data-disabled")) {
          const update = updateFor(radio)
          if (update.menuAccent && !isTranslucent(update.menuColor ?? state.menuColor)) lastSolidAccent = state.menuAccent
          override = null
          clearTimeout(overrideTimer)
          commit(update)
          if (isMobile()) close()
          return
        }
        const item = t.closest("[data-create-menu-action]")
        if (item && !item.hasAttribute("data-disabled")) {
          close()
          actions[item.getAttribute("data-create-menu-action")]?.()
          return
        }
        if (!t.closest("[role=menu]")) close()
        return
      }
      if (open && !t.closest("[data-create-menu-trigger]")) close()
      const lock = t.closest("[data-create-lock]")
      if (lock) {
        const param = lock.getAttribute("data-create-lock")
        if (locks.has(param)) locks.delete(param)
        else locks.add(param)
        render()
        return
      }
      const action = t.closest("[data-create-action]")
      if (action) {
        actions[action.getAttribute("data-create-action")]?.(action)
        return
      }
      const itemBtn = t.closest("[data-create-item]")
      if (itemBtn) commit({ item: itemBtn.dataset.createItem })
    }
    document.addEventListener("click", onClick)

    const onMove = (e) => {
      if (!open || !(e.target instanceof Element)) return
      const item = e.target.closest("[role=menuitemradio],[role=menuitem]")
      if (!item || !open.popup.contains(item)) return
      if (document.activeElement !== item) item.focus({ preventScroll: true })
      if (item.getAttribute("role") === "menuitemradio" && !isMobile() && open.name !== "base" && !item.hasAttribute("data-disabled")) {
        if (open.name === "radius" && (state.style === "lyra" || state.style === "sera")) return
        setOverride(updateFor(item))
      }
    }
    document.addEventListener("mousemove", onMove)

    const onKey = (e) => {
      if (open) {
        const items = [...open.popup.querySelectorAll("[role=menuitemradio]:not([hidden]):not([data-disabled]),[role=menuitem]:not([data-disabled])")]
        const i = items.indexOf(document.activeElement)
        if (e.key === "Escape") {
          e.preventDefault()
          const trigger = open.trigger
          close()
          trigger.focus()
          return
        }
        if (e.key === "ArrowDown" || e.key === "ArrowUp") {
          e.preventDefault()
          const next = items[(i + (e.key === "ArrowDown" ? 1 : -1) + items.length) % items.length] ?? items[0]
          next?.focus({ preventScroll: false })
          if (next?.getAttribute("role") === "menuitemradio" && !isMobile()) setOverride(updateFor(next))
          return
        }
        if ((e.key === "Enter" || e.key === " ") && i >= 0) {
          e.preventDefault()
          items[i].click()
          return
        }
        if (e.key === "Tab") close()
      }
      if (editable(e.target) || document.querySelector("[data-slot=dialog-overlay],[data-slot=alert-dialog-overlay]")) return
      const key = e.key.toLowerCase()
      if ((e.metaKey || e.ctrlKey) && ((key === "z" && e.shiftKey) || (key === "y" && e.ctrlKey))) {
        e.preventDefault()
        goHistory(1)
      } else if ((e.metaKey || e.ctrlKey) && key === "z") {
        e.preventDefault()
        goHistory(-1)
      } else if ((e.metaKey || e.ctrlKey) && (key === "p" || key === "k")) {
        e.preventDefault()
        actions.navigate()
      } else if (!e.metaKey && !e.ctrlKey && !e.altKey) {
        if (e.key === "r") {
          e.preventDefault()
          shuffle()
        } else if (e.key === "R" && e.shiftKey) {
          e.preventDefault()
          actions.reset()
        } else if (key === "o") {
          e.preventDefault()
          actions["open-preset"]()
        }
      }
    }
    document.addEventListener("keydown", onKey)

    // The preview's own shortcuts are forwarded here (keydown events).
    const onMessage = (e) => {
      if (e.origin !== location.origin || e.source !== iframe.contentWindow || !e.data || typeof e.data !== "object") return
      const fwd = e.data
      const isMac = /Mac|iPhone|iPad|iPod/.test(navigator.userAgent)
      const keys = {
        "cmd-k-forward": { key: fwd.key || "k", metaKey: isMac, ctrlKey: !isMac },
        "randomize-forward": { key: fwd.key || "r" },
        "reset-forward": { key: "R", shiftKey: true },
        "open-preset-forward": { key: fwd.key || "o" },
        "undo-forward": { key: "z", metaKey: isMac, ctrlKey: !isMac },
        "redo-forward": { key: "z", shiftKey: true, metaKey: isMac, ctrlKey: !isMac },
        "dark-mode-forward": { key: fwd.key || "d" },
      }[fwd.type]
      if (keys) document.dispatchEvent(new KeyboardEvent("keydown", { ...keys, bubbles: true, cancelable: true }))
    }
    window.addEventListener("message", onMessage)

    // The site theme changes the menu picker's availability (Inverted).
    const themeObserver = new MutationObserver(() => open && syncMenu(open))
    themeObserver.observe(document.documentElement, { attributes: true, attributeFilter: ["class"] })

    // Platform shortcuts in the main menu.
    const isMac = /Mac|iPhone|iPad|iPod/.test(navigator.platform || navigator.userAgent)
    if (isMac) {
      const tpl = root.querySelector('template[data-create-menu="main"]')
      for (const [from, to] of [["Ctrl+P", "⌘P"], ["Ctrl+Z", "⌘Z"], ["Ctrl+Shift+Z", "⇧⌘Z"]]) {
        const el = tpl?.content.querySelector(`[data-create-shortcut="${from}"]`)
        if (el) el.textContent = to
      }
    }
    // Undo/Redo availability in the main menu.
    document.addEventListener("create:state", () => {
      const tpl = root.querySelector('template[data-create-menu="main"]')
      for (const [action, ok] of [["undo", history.index > 0], ["redo", history.index < history.max]]) {
        for (const scope of [tpl?.content, open?.name === "main" ? open.popup : null]) {
          const el = scope?.querySelector(`[data-create-menu-action="${action}"]`)
          if (!el) continue
          el.toggleAttribute("data-disabled", !ok)
          if (ok) el.removeAttribute("aria-disabled")
          else el.setAttribute("aria-disabled", "true")
        }
      }
    })

    // Welcome dialog: first visit only.
    const welcome = () => getDialog(root.querySelector('[data-create-dialog="welcome"] template[data-radix-dialog]'))
    let welcomed = false
    try {
      welcomed = !!localStorage.getItem("shadcn-create-welcome-dialog")
    } catch {}
    if (!welcomed) {
      const show = () => {
        const d = welcome()
        if (!d) return setTimeout(show, 16)
        d.setOpen(true)
        const tpl = root.querySelector('[data-create-dialog="welcome"] template[data-radix-dialog]')
        tpl?.addEventListener("shadcn:open-change", (e) => {
          if (!e.detail.open) {
            try {
              localStorage.setItem("shadcn-create-welcome-dialog", "true")
            } catch {}
          }
        })
      }
      show()
    }

    // Exposed for create-dialogs.js.
    root.createAPI = { get state() { return state }, commit, config: () => config, close, locks, history, goHistory }
    render()
    sync()

    return () => {
      document.removeEventListener("click", onClick)
      document.removeEventListener("mousemove", onMove)
      document.removeEventListener("keydown", onKey)
      window.removeEventListener("message", onMessage)
      themeObserver.disconnect()
      close()
    }
  },
})
