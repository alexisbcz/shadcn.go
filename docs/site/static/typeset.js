// The typeset designer (app/(app)/(typeset)): the page's client state.
// The server renders the page for the URL's params; this keeps it in sync
// when they change: the search params (nuqs, history "replace"), the
// pickers (radio groups, labels, the "Aa" font), locks (use-locks.tsx),
// shuffle/reset (use-shuffle.tsx), undo/redo (use-history.tsx), hover
// previews (preview-override.tsx), the toolbar, the preview iframe
// (preview.tsx: postMessage of the params, reload when the item changes,
// commands forwarded by the iframe) and the docs (docs-panel.tsx: tabs,
// framework, package manager, copy buttons; the docs body is re-rendered
// by the server at /typeset/docs-content).
import { define } from "./core.js"
import { setRadioValue } from "./menu.js"

const KEYS = ["body", "heading", "mono", "scale", "measure", "flow", "leading", "item"]
const LOCKABLE = KEYS.filter((k) => k !== "item")
const PARAMS_MESSAGE = "typeset-params"
const COMMAND_MESSAGE = "typeset-command"
const OVERRIDE_DEBOUNCE_MS = 50

const isEditable = (t) =>
  (t instanceof HTMLElement && t.isContentEditable) ||
  t instanceof HTMLInputElement ||
  t instanceof HTMLTextAreaElement ||
  t instanceof HTMLSelectElement

define("typeset", {
  selector: "[data-typeset]",
  mount(root) {
    const data = JSON.parse(root.dataset.typesetData)
    const { defaults, values, fonts, labels, icons } = data
    let params = JSON.parse(root.dataset.typesetParams)
    let override = null
    let overrideTimer = 0
    let lastSent = null
    let framework = "next"
    const locks = new Set()
    const iframe = root.querySelector("iframe[data-typeset-preview]")
    const isMobile = () => matchMedia("(max-width: 767px)").matches
    const findFont = (id) => fonts.find((f) => f.id === id)

    // Every element in the page and in its templates (closed menus, the
    // drawer): state is rendered into both.
    function* roots(node = document) {
      yield node
      for (const t of node.querySelectorAll("template")) yield* roots(t.content)
    }
    const all = (selector) => [...roots()].flatMap((r) => [...r.querySelectorAll(selector)])

    // --- params ------------------------------------------------------------

    const serialize = (base, p) => {
      const q = new URLSearchParams()
      for (const k of KEYS) if (p[k] !== defaults[k]) q.set(k, p[k])
      const s = q.toString()
      return s ? `${base}?${s}` : base
    }
    const coerce = (key, value) => (values[key].includes(value) ? value : null)

    function setParams(next) {
      const prev = params
      if (next === null) params = { ...defaults }
      else {
        params = { ...params }
        for (const [k, v] of Object.entries(next)) params[k] = v === null ? defaults[k] : v
      }
      const url = new URL(location.href)
      for (const k of KEYS) {
        if (params[k] === defaults[k]) url.searchParams.delete(k)
        else url.searchParams.set(k, params[k])
      }
      history.replaceState(history.state, "", url)
      render(prev)
      record()
    }

    function render(prev) {
      // Picker labels and radio groups.
      const body = findFont(params.body) ?? fonts[0]
      for (const key of LOCKABLE) {
        const v = params[key]
        let label
        if (key === "body" || key === "heading" || key === "mono") {
          const font = key === "heading" && v === "inherit" ? body : (findFont(v) ?? body)
          label = font.label
          for (const aa of all(`[data-typeset-aa="${key}"]`)) aa.style.fontFamily = font.value
        } else label = labels[key][v]
        for (const el of all(`[data-typeset-label="${key}"]`)) el.textContent = label
        for (const key2 of radioKeys(key)) setRadioValue(key2, v)
      }
      for (const el of all("[data-typeset-inherit-label]")) el.textContent = body.label
      // Toolbar.
      for (const pill of root.querySelectorAll("[data-typeset-item]")) {
        pill.dataset.active = String(pill.dataset.typesetItem === params.item)
      }
      const href = serialize(`/preview/typeset/${params.item}`, params)
      for (const a of root.querySelectorAll("[data-typeset-open]")) a.setAttribute("href", href)
      sendPreview(prev)
      // Docs.
      if (!prev || KEYS.some((k) => prev[k] !== params[k])) refreshDocs()
    }

    // The menu state key of the radio group of a picker (the marker
    // attribute of its content template).
    const radioKeyCache = new Map()
    function radioKeys(param) {
      if (!radioKeyCache.has(param)) {
        const keys = []
        for (const t of document.querySelectorAll("template")) {
          const g = t.content.querySelector(`[data-typeset-group="${param}"]`)
          const key = g?.getAttribute("data-radix-menu-radio-group")
          if (key) keys.push(key)
        }
        radioKeyCache.set(param, keys)
      }
      return radioKeyCache.get(param)
    }

    // --- preview iframe ----------------------------------------------------

    const merged = () => (override ? { ...params, ...override } : params)
    let currentItem = params.item
    function sendPreview() {
      if (!iframe) return
      const p = merged()
      if (p.item !== currentItem) {
        currentItem = p.item
        lastSent = null
        iframe.src = serialize(`/preview/typeset/${p.item}`, p)
        return
      }
      if (lastSent && KEYS.every((k) => lastSent[k] === p[k])) return
      post(p)
    }
    function post(p = merged()) {
      iframe?.contentWindow?.postMessage({ type: PARAMS_MESSAGE, data: p }, location.origin)
      lastSent = p
    }
    iframe?.addEventListener("load", () => post())

    function setOverride(next) {
      clearTimeout(overrideTimer)
      overrideTimer = setTimeout(() => {
        if (override && Object.keys(next).every((k) => override[k] === next[k]) && Object.keys(override).length === Object.keys(next).length) return
        override = next
        sendPreview()
      }, OVERRIDE_DEBOUNCE_MS)
    }
    function clearOverride() {
      clearTimeout(overrideTimer)
      if (!override) return
      override = null
      sendPreview()
    }

    // --- history (use-history.tsx) ------------------------------------------

    let entries = [JSON.stringify(params)]
    let index = 0
    let navigating = false
    function record() {
      const snapshot = JSON.stringify(params)
      if (navigating) {
        navigating = false
        return
      }
      if (snapshot === entries[index]) return
      entries = entries.slice(0, index + 1)
      entries.push(snapshot)
      index = entries.length - 1
      renderHistory()
    }
    function go(delta) {
      const next = index + delta
      if (next < 0 || next >= entries.length) return
      index = next
      navigating = true
      setParams(JSON.parse(entries[index]))
      renderHistory()
    }
    function renderHistory() {
      for (const [action, enabled] of [["undo", index > 0], ["redo", index < entries.length - 1]]) {
        for (const item of all(`[data-typeset-action="${action}"]`)) {
          if (enabled) {
            item.removeAttribute("data-disabled")
            item.removeAttribute("aria-disabled")
          } else {
            item.setAttribute("data-disabled", "")
            item.setAttribute("aria-disabled", "true")
          }
        }
      }
    }
    const isMac = /Mac|iPhone|iPad|iPod/.test(navigator.platform || navigator.userAgent)
    if (!isMac) {
      for (const el of all('[data-typeset-shortcut="undo"]')) el.textContent = "Ctrl+Z"
      for (const el of all('[data-typeset-shortcut="redo"]')) el.textContent = "Ctrl+Shift+Z"
    }

    // --- shuffle (use-shuffle.tsx) -----------------------------------------

    const randomItem = (items) => items[Math.floor(Math.random() * items.length)]
    function shuffle() {
      const bodyFonts = fonts.filter((f) => f.type !== "mono").map((f) => f.id)
      const monoFonts = fonts.filter((f) => f.type === "mono").map((f) => f.id)
      const next = {
        body: randomItem(bodyFonts),
        heading: randomItem(["inherit", ...bodyFonts]),
        mono: randomItem(monoFonts),
        scale: randomItem(values.scale),
        measure: randomItem(values.measure),
        leading: randomItem(values.leading),
        flow: randomItem(values.flow),
      }
      for (const k of locks) delete next[k]
      setParams(next)
    }
    const reset = () => setParams(null)
    const toggleTheme = () => document.querySelector('[data-site-action="toggle-theme"]')?.click()
    const commands = { shuffle, reset, undo: () => go(-1), redo: () => go(1), "toggle-theme": toggleTheme }

    // --- locks (lock-button.tsx) --------------------------------------------

    function toggleLock(param) {
      if (locks.has(param)) locks.delete(param)
      else locks.add(param)
      const locked = locks.has(param)
      for (const b of root.querySelectorAll(`[data-typeset-lock="${param}"]`)) {
        b.dataset.locked = String(locked)
        b.title = locked ? "Unlock" : "Lock"
        b.setAttribute("aria-label", locked ? "Unlock" : "Lock")
        const svg = b.querySelector("svg")
        if (svg) svg.innerHTML = locked ? icons.lock : icons.unlock
      }
    }

    // --- docs (docs-panel.tsx) ----------------------------------------------

    let docsRequest = 0
    async function refreshDocs() {
      const id = ++docsRequest
      const url = serialize("/typeset/docs-content", params)
      const html = await (await fetch(url + (url.includes("?") ? "&" : "?") + "framework=" + framework)).text()
      if (id !== docsRequest) return
      for (const body of all("[data-typeset-docs-body]")) {
        body.innerHTML = html
        const tabs = body.closest("[data-typeset-tabs]")
        if (tabs) selectTab(tabs, tabs.dataset.typesetTab ?? "docs")
        applyPackageManager(body)
      }
    }

    function selectTab(tabs, value) {
      tabs.dataset.typesetTab = value
      for (const t of tabs.querySelectorAll("[data-typeset-tab]")) {
        const active = t.dataset.typesetTab === value
        if (active) t.setAttribute("data-active", "")
        else t.removeAttribute("data-active")
        t.setAttribute("aria-selected", String(active))
        t.tabIndex = active ? 0 : -1
      }
      for (const p of tabs.querySelectorAll("[data-typeset-panel]")) p.hidden = p.dataset.typesetPanel !== value
    }

    const packageManager = () => {
      try {
        return JSON.parse(localStorage.getItem("config") ?? "{}").packageManager || "pnpm"
      } catch {
        return "pnpm"
      }
    }
    function applyPackageManager(scope = document) {
      const pm = packageManager()
      for (const b of scope.querySelectorAll("[data-typeset-pm]")) b.dataset.active = String(b.dataset.typesetPm === pm)
      for (const pre of scope.querySelectorAll("[data-typeset-pm-command]")) pre.hidden = pre.dataset.typesetPmCommand !== pm
    }
    for (const r of roots()) applyPackageManager(r)

    async function copy(button) {
      let text
      if (button.dataset.typesetCopy === "css") text = await (await fetch("/typeset.css")).text()
      else text = button.closest("[data-typeset-panel]")?.querySelector("pre")?.textContent ?? ""
      await navigator.clipboard?.writeText(text)
      const svg = button.querySelector("svg")
      if (!svg) return
      const prev = svg.outerHTML
      svg.outerHTML = `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-check" data-icon="inline-start"><path d="M20 6 9 17l-5-5"></path></svg>`
      setTimeout(() => {
        const now = button.querySelector("svg")
        if (now) now.outerHTML = prev
      }, 2000)
    }

    // --- events -------------------------------------------------------------

    const onClick = (e) => {
      const t = e.target instanceof Element ? e.target : null
      if (!t) return
      const pill = t.closest("[data-typeset-item]")
      if (pill && root.contains(pill)) return setParams({ item: pill.dataset.typesetItem })
      const lock = t.closest("[data-typeset-lock]")
      if (lock && root.contains(lock)) return toggleLock(lock.dataset.typesetLock)
      const button = t.closest("button[data-typeset-action]")
      if (button && root.contains(button)) return commands[button.dataset.typesetAction]?.()
      const tab = t.closest("[data-typeset-tab]")
      if (tab) return selectTab(tab.closest("[data-typeset-tabs]"), tab.dataset.typesetTab)
      const pm = t.closest("[data-typeset-pm]")
      if (pm) {
        let config = {}
        try {
          config = JSON.parse(localStorage.getItem("config") ?? "{}")
        } catch {}
        try {
          localStorage.setItem("config", JSON.stringify({ style: "new-york-v4", installationType: "cli", ...config, packageManager: pm.dataset.typesetPm }))
        } catch {}
        for (const r of roots()) applyPackageManager(r)
        return
      }
      const copyButton = t.closest("[data-typeset-copy]")
      if (copyButton) copy(copyButton)
    }
    // Tabs: arrow keys move between the triggers (Base UI activates on focus).
    const onKeyDown = (e) => {
      const t = e.target instanceof Element ? e.target : null
      const tab = t?.closest("[data-typeset-tab]")
      if (tab && (e.key === "ArrowLeft" || e.key === "ArrowRight")) {
        const list = [...tab.parentElement.querySelectorAll("[data-typeset-tab]")]
        const next = list[(list.indexOf(tab) + (e.key === "ArrowRight" ? 1 : list.length - 1)) % list.length]
        next.focus()
        selectTab(tab.closest("[data-typeset-tabs]"), next.dataset.typesetTab)
        e.preventDefault()
        return
      }
      if (e.altKey || isEditable(e.target)) return
      if (e.metaKey || e.ctrlKey) {
        const key = e.key.toLowerCase()
        if ((key === "z" && e.shiftKey) || (key === "y" && e.ctrlKey)) {
          e.preventDefault()
          go(1)
        } else if (key === "z") {
          e.preventDefault()
          go(-1)
        }
        return
      }
      if (e.key === "r" || e.key === "R") {
        e.preventDefault()
        if (e.shiftKey) reset()
        else shuffle()
      }
    }
    // Menus: radio groups (controlled: the params are the state), actions
    // and hover previews.
    const onValueChange = (e) => {
      const t = e.target instanceof Element ? e.target : null
      const group = t?.closest("[data-typeset-group]")
      if (group) {
        const key = group.dataset.typesetGroup
        const v = coerce(key, String(e.detail?.value))
        if (v !== null) setParams({ [key]: v })
        return
      }
      if (t?.closest("[data-typeset-framework]")) {
        framework = String(e.detail?.value ?? "next")
        refreshDocs()
      }
    }
    const onSelect = (e) => {
      const t = e.target instanceof Element ? e.target : null
      // Base UI radio items do not close the menu (closeOnClick false).
      if (t?.closest('[role="menuitemradio"]')?.closest("[data-typeset-picker]")) e.preventDefault()
      const item = t?.closest("[data-typeset-action]")
      if (item) commands[item.dataset.typesetAction]?.()
    }
    const onOpenChange = (e) => {
      if (e.detail && e.detail.open === false) clearOverride()
      // Pickers open to the right, or above (anchored on the customizer,
      // which the shift keeps them over) on mobile: the content template's
      // popper options are read when it opens.
      const tpl = e.target instanceof HTMLTemplateElement ? e.target : null
      const wrapper = e.detail?.open && tpl?.content.querySelector("[data-radix-popper]")
      if (!wrapper || !wrapper.querySelector("[data-typeset-picker]")) return
      const opts = JSON.parse(wrapper.getAttribute("data-radix-popper"))
      // Base UI's collision padding (5px) leaves --available-width 4px
      // short of the viewport on each side.
      Object.assign(opts, isMobile() ? { side: "top", align: "center", collisionPadding: 4 } : { side: "right", align: "start", collisionPadding: 0 })
      wrapper.setAttribute("data-radix-popper", JSON.stringify(opts))
      if (isMobile()) tpl.dataset.anchor = "[data-typeset-customizer]"
      else delete tpl.dataset.anchor
    }
    const onPreview = (e) => {
      if (isMobile()) return
      const t = e.target instanceof Element ? e.target : null
      const item = t?.closest("[data-slot=dropdown-menu-radio-item][data-value]")
      const picker = item?.closest("[data-typeset-picker]")
      if (!picker) return
      const key = picker.dataset.typesetPicker
      const v = coerce(key, item.dataset.value)
      if (v !== null) setOverride({ [key]: v })
    }
    const onLeave = (e) => {
      const t = e.target instanceof Element ? e.target : null
      if (t?.matches("[data-typeset-picker]")) clearOverride()
    }
    const onMessage = (e) => {
      if (!iframe || e.origin !== location.origin || e.source !== iframe.contentWindow) return
      if (e.data?.type === COMMAND_MESSAGE && typeof e.data.command === "string") commands[e.data.command]?.()
    }

    document.addEventListener("click", onClick)
    document.addEventListener("keydown", onKeyDown)
    document.addEventListener("shadcn:value-change", onValueChange)
    document.addEventListener("shadcn:select", onSelect)
    document.addEventListener("shadcn:open-change", onOpenChange)
    document.addEventListener("mousemove", onPreview)
    document.addEventListener("focusin", onPreview)
    document.addEventListener("mouseleave", onLeave, true)
    window.addEventListener("message", onMessage)
    return () => {
      document.removeEventListener("click", onClick)
      document.removeEventListener("keydown", onKeyDown)
      document.removeEventListener("shadcn:value-change", onValueChange)
      document.removeEventListener("shadcn:select", onSelect)
      document.removeEventListener("shadcn:open-change", onOpenChange)
      document.removeEventListener("mousemove", onPreview)
      document.removeEventListener("focusin", onPreview)
      document.removeEventListener("mouseleave", onLeave, true)
      window.removeEventListener("message", onMessage)
    }
  },
})
