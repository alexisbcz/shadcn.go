// The create designer's dialogs (see create_dialogs.go): Reset
// (use-reset.tsx), Open Preset (open-preset.tsx, a drawer on mobile), Get
// Code (project-form.tsx) and the action menu (action-menu.tsx).
import { define } from "./core.js"
import { getDialog } from "./dialog.js"
import { getDrawer } from "./drawer.js"
import { decodePreset, encodePreset, isPresetCode, presetThemeCSS } from "./create.js"

const tplOf = (name, attr = "data-radix-dialog") => document.querySelector(`[data-create-dialog="${name}"] template[${attr}]`)
const isMobile = () => matchMedia("(max-width: 767px)").matches

// parsePresetInput: a code, or "--preset <code>".
function parsePresetInput(value) {
  const input = value.trim()
  if (!input) return null
  const preset = input.match(/^--preset\b\s+(.+)$/i)?.[1]?.trim() ?? input
  return isPresetCode(preset) ? preset : null
}

define("create-dialogs", {
  selector: "[data-create]",
  mount(root) {
    const api = () => root.createAPI
    const dialogs = {
      reset: () => getDialog(tplOf("reset")),
      "open-preset": () => (isMobile() ? getDrawer(tplOf("open-preset-mobile", "data-vaul-root")) : getDialog(tplOf("open-preset"))),
      "project-form": () => getDialog(tplOf("project-form")),
      "action-menu": () => getDialog(tplOf("action-menu")),
    }
    const isOpen = (name) => !!dialogs[name]()?.open
    const open = (name) => {
      api()?.close()
      // Let a closing menu restore focus first.
      setTimeout(() => {
        dialogs[name]()?.setOpen(true)
        if (name === "project-form") requestAnimationFrame(syncProjectForm)
      })
    }

    // --- Get Code ---
    let projectTab = "new-project"
    let applyMode = "full"
    let packageManager = "pnpm"
    try {
      packageManager = JSON.parse(localStorage.getItem("config") ?? "{}").packageManager || "pnpm"
    } catch {}
    const copied = {}
    const commands = (s) => {
      const code = encodePreset(s)
      const framework = s.template.replace("-monorepo", "")
      const flags = ` --preset ${code}${s.base !== "base" ? ` --base ${s.base}` : ""} --template ${framework}${s.template.endsWith("-monorepo") ? " --monorepo" : ""}${s.rtl ? " --rtl" : ""}${s.pointer ? " --pointer" : ""}`
      const only = applyMode === "theme" ? " --only theme" : applyMode === "font" ? " --only font" : ""
      // NEXT_PUBLIC_APP_URL is unset in the reference build: local dev
      // commands.
      return { command: `shadcn init${flags}`, apply: `shadcn apply --preset ${code}${only}`, theme: presetThemeCSS(s) }
    }
    const content = () => document.getElementById(tplOf("project-form")?.getAttribute("data-content-id"))
    function syncProjectForm() {
      const el = content()
      const s = api()?.state
      if (!el || !s) return
      const c = commands(s)
      for (const code of el.querySelectorAll("[data-create-command]")) {
        const [kind] = code.getAttribute("data-create-command").split(":")
        code.textContent = c[kind]
      }
      const css = el.querySelector("[data-create-theme-css]")
      if (css) css.textContent = c.theme
      for (const panel of el.querySelectorAll("[data-create-project-tab]")) panel.hidden = panel.getAttribute("data-create-project-tab") !== projectTab
      for (const [id, on, disabled] of [["pointer", s.pointer], ["rtl", s.rtl], ["monorepo", s.template.endsWith("-monorepo"), s.template.startsWith("laravel")]]) {
        const sw = el.querySelector(`[data-create-option="${id}"]`)
        if (!sw) continue
        sw.setAttribute("aria-checked", String(on))
        sw.dataset.state = on ? "checked" : "unchecked"
        for (const t of sw.querySelectorAll("[data-state]")) t.dataset.state = on ? "checked" : "unchecked"
        sw.disabled = !!disabled
        sw.closest("[data-slot=field]")?.toggleAttribute("data-disabled", !!disabled)
      }
      for (const [sel, value] of [["[data-create-template]", s.template.replace("-monorepo", "")], ["[data-create-base]", s.base], ["[data-create-apply-mode]", applyMode]]) {
        const group = el.querySelector(sel)
        for (const item of group?.querySelectorAll("[data-slot=radio-group-item]") ?? []) {
          const on = item.value === value || item.getAttribute("value") === value
          item.setAttribute("aria-checked", String(on))
          item.dataset.state = on ? "checked" : "unchecked"
          const ind = item.querySelector("[data-slot=radio-group-indicator]")
          if (ind) ind.hidden = !on
        }
      }
      for (const [kind, label] of [["command", "Copy Command"], ["apply", "Copy Command"], ["theme", "Copy Theme"]]) {
        for (const b of el.querySelectorAll(`[data-create-copy="${kind}"]`)) {
          if (b.hasAttribute("data-create-copy-label")) b.textContent = copied[kind] ? "Copied" : label
          else {
            const icon = document.querySelector(`template[data-create-copy-icon="${copied[kind] ? "tick" : "copy"}"]`)
            const svg = b.querySelector("svg")
            if (icon && svg) svg.replaceWith(icon.content.cloneNode(true))
          }
        }
      }
      // The package manager tabs follow the stored choice.
      for (const tabsRoot of el.querySelectorAll("[data-create-package-manager]")) {
        const trigger = [...tabsRoot.querySelectorAll("[data-slot=tabs-trigger]")].find((t) => t.textContent.trim() === packageManager)
        if (trigger && trigger.getAttribute("aria-selected") !== "true") trigger.dispatchEvent(new MouseEvent("mousedown", { bubbles: true, button: 0 }))
      }
    }

    const onCreateOpen = (e) => {
      const name = e.detail
      if (name === "action-menu") return open(isOpen("action-menu") ? null : "action-menu")
      open(name)
    }
    document.addEventListener("create:open", onCreateOpen)
    document.addEventListener("create:state", () => content() && syncProjectForm())

    const onClick = (e) => {
      const t = e.target instanceof Element ? e.target : null
      if (!t) return
      const opener = t.closest("[data-create-open]")
      if (opener) return open(opener.getAttribute("data-create-open"))
      if (t.closest("[data-create-reset-confirm]")) return reset()
      const copy = t.closest("[data-create-copy]")
      if (copy) {
        const kind = copy.getAttribute("data-create-copy")
        navigator.clipboard?.writeText(commands(api().state)[kind])
        for (const k in copied) delete copied[k]
        copied[kind] = true
        syncProjectForm()
        clearTimeout(copy._t)
        copy._t = setTimeout(() => {
          delete copied[kind]
          syncProjectForm()
        }, 2000)
        return
      }
    }
    document.addEventListener("click", onClick, true)

    const onValueChange = (e) => {
      const t = e.target instanceof Element ? e.target : null
      if (!t) return
      const s = api()?.state
      if (t.closest("[data-create-project-tabs]")) {
        projectTab = e.detail.pressedItems?.[0] ?? "new-project"
        syncProjectForm()
      } else if (t.closest("[data-create-template]")) {
        const v = e.detail.value
        api().commit({ template: v === "laravel" ? v : s.template.endsWith("-monorepo") ? `${v}-monorepo` : v })
      } else if (t.closest("[data-create-base]")) {
        api().commit({ base: e.detail.value })
      } else if (t.closest("[data-create-apply-mode]")) {
        applyMode = e.detail.value
        syncProjectForm()
      } else if (t.closest("[data-create-package-manager]") && packageManagers.includes(e.detail.value)) {
        packageManager = e.detail.value
        try {
          const cfg = JSON.parse(localStorage.getItem("config") ?? "{}")
          localStorage.setItem("config", JSON.stringify({ ...cfg, packageManager }))
        } catch {}
      }
    }
    const packageManagers = ["pnpm", "npm", "yarn", "bun"]
    document.addEventListener("shadcn:value-change", onValueChange)
    const onChecked = (e) => {
      const option = e.target instanceof Element && e.target.closest("[data-create-option]")
      if (!option) return
      const s = api().state
      const id = option.getAttribute("data-create-option")
      if (id === "pointer") api().commit({ pointer: e.detail.checked })
      if (id === "rtl") api().commit({ rtl: e.detail.checked })
      if (id === "monorepo") {
        const framework = s.template.replace("-monorepo", "")
        api().commit({ template: e.detail.checked && framework !== "laravel" ? `${framework}-monorepo` : framework })
      }
    }
    document.addEventListener("shadcn:checked-change", onChecked)

    // --- Open Preset ---
    const onInput = (e) => {
      const input = e.target instanceof Element && e.target.closest("[data-create-preset-input]")
      if (!input) return
      const ok = parsePresetInput(input.value)
      const invalid = input.value.trim().length > 0 && !ok
      input.setAttribute("aria-invalid", String(invalid))
      input.closest("[data-create-preset-field]")?.toggleAttribute("data-invalid", invalid)
      const submit = input.closest("form")?.querySelector("button[type=submit]")
      if (submit) submit.disabled = !ok
    }
    document.addEventListener("input", onInput)
    const onSubmit = (e) => {
      const form = e.target instanceof Element && e.target.closest("[data-create-preset-form]")
      if (!form) return
      e.preventDefault()
      const code = parsePresetInput(form.querySelector("[data-create-preset-input]").value)
      const decoded = code && decodePreset(code)
      if (!decoded) return
      decoded.chartColor = decoded.chartColor ?? decoded.theme
      api().commit(decoded)
      dialogs["open-preset"]()?.setOpen(false)
    }
    document.addEventListener("submit", onSubmit, true)

    // --- Reset ---
    function reset() {
      const s = api().state
      const config = api().config()
      const preset = config.presets.find((p) => p.base === s.base && p.style === s.style) ?? { baseColor: "neutral", theme: "neutral", chartColor: "neutral", iconLibrary: "lucide", font: "inter", fontHeading: "inherit", menuAccent: "subtle", menuColor: "default", radius: "default" }
      api().commit({
        baseColor: preset.baseColor, theme: preset.theme, chartColor: preset.chartColor ?? "neutral", iconLibrary: preset.iconLibrary,
        font: preset.font, fontHeading: preset.fontHeading, menuAccent: preset.menuAccent, menuColor: preset.menuColor, radius: preset.radius,
        template: "next",
      })
      dialogs.reset()?.setOpen(false)
    }
    const onKey = (e) => {
      // Shift+R again confirms the open reset dialog.
      if (e.key === "R" && e.shiftKey && !e.metaKey && !e.ctrlKey && isOpen("reset")) {
        e.preventDefault()
        e.stopImmediatePropagation()
        reset()
      }
    }
    document.addEventListener("keydown", onKey, true)

    // --- action menu ---
    const onSelect = (e) => {
      const item = e.target instanceof Element && e.target.closest("[data-create-navigate]")
      if (!item) return
      api().commit({ item: item.getAttribute("data-create-navigate") })
      dialogs["action-menu"]()?.setOpen(false)
    }
    document.addEventListener("shadcn:select", onSelect)
    document.addEventListener("create:state", (e) => {
      for (const scope of [tplOf("action-menu")?.content, document]) {
        for (const item of scope?.querySelectorAll("[data-create-navigate]") ?? []) item.dataset.checked = String(item.getAttribute("data-create-navigate") === e.detail.state.item)
      }
    })

    return () => {
      document.removeEventListener("create:open", onCreateOpen)
      document.removeEventListener("click", onClick, true)
      document.removeEventListener("shadcn:value-change", onValueChange)
      document.removeEventListener("shadcn:checked-change", onChecked)
      document.removeEventListener("input", onInput)
      document.removeEventListener("submit", onSubmit, true)
      document.removeEventListener("keydown", onKey, true)
      document.removeEventListener("shadcn:select", onSelect)
    }
  },
})
