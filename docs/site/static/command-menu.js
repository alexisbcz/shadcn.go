// command-menu.js: the behavior of the docs search palette
// (components/command-menu.tsx), for the markup docs/site/search.go renders.
//
// The palette is a cmdk Command with a custom filter (an item matches when
// its value and keywords, joined, contain the search, case-insensitively;
// all matches score 1, so nothing is reordered) inside a Radix Dialog
// (dialog.js). This module ports the parts of cmdk 1.1.1 the palette uses:
// the first item is selected on mount and after every search change,
// ArrowUp/ArrowDown (Ctrl+N/J/P/K), Home/End, Meta+Arrow (first/last),
// Alt+Arrow (previous/next group) move the selection without wrapping,
// pointer movement selects, Enter or a click runs the item; groups without
// a match are hidden and the empty message shows when nothing matches.
//
// command-menu.tsx itself: ⌘K / Ctrl+K / "/" toggle the dialog (not from
// inputs), the footer shows the highlighted item's action and ⌘C payload,
// ⌘C copies it and closes, the search is sent (after 500 ms, then
// useDocsSearch's 100 ms) to /api/search and the de-duplicated results
// fill the "Search Results" group.

import { define } from "./core.js"
import { getDialog } from "./dialog.js"

const ITEM = "[cmdk-item]"
const GROUP = "[cmdk-group]"

// CommandMenu state that outlives the dialog content.
let selectedType = null // "page" | "component" | "color" | "style" | null
let copyPayload = ""

function packageManager() {
  try {
    const c = JSON.parse(localStorage.getItem("config") || "{}")
    return c.packageManager || "pnpm"
  } catch {
    return "pnpm"
  }
}

function trigger() {
  return document.querySelector("[data-command-menu-trigger]")
}

function dialog() {
  const t = trigger()
  return t ? getDialog(t) : undefined
}

function setOpen(open) {
  const d = dialog()
  if (d && d.open !== open) d.setOpen(open)
}

function copy(text) {
  try {
    navigator.clipboard?.writeText(text)
  } catch {}
}

// runCommand: close, then run.
function runCommand(fn) {
  setOpen(false)
  fn()
}

// router.push keeps the header, and Radix gives the trigger focus back when
// the dialog closes. A Go navigation reloads the page: the trigger gets the
// focus back on the next page, visibly so after a keyboard interaction.
const REFOCUS = "shadcn-command-menu-refocus"
let lastKeyboard = false
addEventListener("keydown", () => (lastKeyboard = true), true)
addEventListener("pointerdown", () => (lastKeyboard = false), true)

function navigate(href) {
  try {
    sessionStorage.setItem(REFOCUS, lastKeyboard ? "keyboard" : "pointer")
  } catch {}
  location.assign(href)
}

try {
  const how = sessionStorage.getItem(REFOCUS)
  if (how) {
    sessionStorage.removeItem(REFOCUS)
    // After load: scrolling to the URL fragment would reset the focus.
    const refocus = () => setTimeout(() => trigger()?.focus({ preventScroll: true, focusVisible: how === "keyboard" }))
    if (document.readyState === "complete") refocus()
    else addEventListener("load", refocus, { once: true })
  }
} catch {}

document.addEventListener("keydown", (e) => {
  if ((e.key === "k" && (e.metaKey || e.ctrlKey)) || e.key === "/") {
    const t = e.target
    if ((t instanceof HTMLElement && t.isContentEditable) || t instanceof HTMLInputElement || t instanceof HTMLTextAreaElement || t instanceof HTMLSelectElement) return
    e.preventDefault()
    const d = dialog()
    if (d) d.setOpen(!d.open)
  }
  if (e.key === "c" && (e.metaKey || e.ctrlKey)) {
    runCommand(() => {
      if (selectedType) copy(copyPayload)
    })
  }
})

function mount(root) {
  const input = root.querySelector("[cmdk-input]")
  const list = root.querySelector("[cmdk-list]")
  const label = root.querySelector("[cmdk-label]")
  const empty = root.querySelector("[cmdk-empty]")
  const results = root.querySelector("[data-command-menu-results]")
  const resultTpl = root.querySelector("template[data-command-menu-result]")
  const spinner = root.querySelector("[data-command-menu-spinner]")
  const content = root.closest('[role="dialog"]') || root.parentElement
  const footer = content.querySelector("[data-command-menu-footer]")

  // ids (React's useId in cmdk).
  const id = (p) => `${p}-${Math.random().toString(36).slice(2, 8)}`
  const listId = id("cmdk-list")
  const labelId = id("cmdk-label")
  const inputId = id("cmdk-input")
  list.id = listId
  label.id = labelId
  label.htmlFor = inputId
  input.id = inputId
  input.setAttribute("aria-controls", listId)
  input.setAttribute("aria-labelledby", labelId)
  for (const g of root.querySelectorAll(GROUP)) {
    const hd = g.querySelector(":scope > [cmdk-group-heading]")
    const items = g.querySelector(":scope > [cmdk-group-items]")
    if (hd && items) {
      hd.id = id("cmdk-group")
      items.setAttribute("aria-labelledby", hd.id)
    }
  }
  let itemSeq = 0
  const prepare = (el) => {
    el.id = `cmdk-item-${++itemSeq}`
    try {
      el._keywords = JSON.parse(el.getAttribute("data-cmdk-keywords") || "[]")
    } catch {
      el._keywords = []
    }
  }
  root.querySelectorAll(ITEM).forEach(prepare)

  let search = ""
  let value = ""

  const allItems = () => [...root.querySelectorAll(ITEM)]
  const validItems = () => allItems().filter((el) => !el.hidden && el.getAttribute("aria-disabled") !== "true")
  const selectedItem = () => allItems().find((el) => el.getAttribute("aria-selected") === "true")

  const matches = (el) => {
    if (!search) return true
    const v = el.getAttribute("data-value") || ""
    if (!v) return false
    return (v + " " + el._keywords.join(" ")).toLowerCase().includes(search.toLowerCase())
  }

  function filter() {
    let count = 0
    for (const el of allItems()) {
      const ok = matches(el)
      el.hidden = !ok
      if (ok) count++
    }
    for (const g of root.querySelectorAll(GROUP)) {
      if (g === results && !results.querySelector(ITEM)) {
        g.hidden = true
        continue
      }
      g.hidden = !!search && !g.querySelector(`${ITEM}:not([hidden])`)
    }
    empty.hidden = count !== 0
  }

  // The footer follows the highlighted item (CommandMenuItem's onHighlight).
  function renderFooter() {
    const action = footer.querySelector("[data-command-menu-action]")
    action.textContent =
      selectedType === "page" || selectedType === "component" ? "Go to Page" : selectedType === "color" ? "Copy OKLCH" : selectedType === "style" ? "Open in shadcn/create" : ""
    footer.querySelectorAll("[data-command-menu-copy]").forEach((el) => (el.hidden = !copyPayload))
    footer.querySelector("[data-command-menu-payload]").textContent = copyPayload
  }

  function highlight(el) {
    const kind = el.getAttribute("data-kind")
    if (!kind) return // search results have no onHighlight
    selectedType = kind
    if (kind === "component") copyPayload = `${packageManager()} dlx shadcn@latest add ${el.getAttribute("data-copy")}`
    else if (kind === "color") copyPayload = el.getAttribute("data-copy")
    else copyPayload = ""
    renderFooter()
  }

  function scrollIntoView(el) {
    const visibleSibs = [...el.parentElement.children].filter((c) => !c.hidden)
    if (visibleSibs[0] === el) el.closest(GROUP)?.querySelector("[cmdk-group-heading]")?.scrollIntoView({ block: "nearest" })
    el.scrollIntoView({ block: "nearest" })
  }

  function select(el, scroll = true) {
    const prev = selectedItem()
    if (prev === el) return
    if (prev) {
      prev.setAttribute("aria-selected", "false")
      prev.setAttribute("data-selected", "false")
    }
    value = el ? el.getAttribute("data-value") : ""
    if (el) {
      el.setAttribute("aria-selected", "true")
      el.setAttribute("data-selected", "true")
      input.setAttribute("aria-activedescendant", el.id)
      highlight(el)
      if (scroll) scrollIntoView(el)
    } else {
      input.removeAttribute("aria-activedescendant")
    }
  }

  const selectFirst = () => select(validItems()[0] ?? null)

  function byItem(change) {
    const items = validItems()
    const cur = selectedItem()
    const i = items.indexOf(cur)
    const next = items[i + change]
    if (next) select(next)
  }

  function byGroup(change) {
    const cur = selectedItem()
    let group = cur?.closest(GROUP)
    let item
    const groups = [...root.querySelectorAll(GROUP)].filter((g) => !g.hidden)
    while (group && !item) {
      group = groups[groups.indexOf(group) + change]
      item = group?.querySelector(`${ITEM}:not([hidden])`)
    }
    if (item) select(item)
    else byItem(change)
  }

  function run(el) {
    if (!el || el.getAttribute("aria-disabled") === "true") return
    select(el, false)
    const kind = el.getAttribute("data-kind")
    if (kind === "color") runCommand(() => copy(el.getAttribute("data-oklch")))
    else if (el.hasAttribute("data-href")) {
      const href = el.getAttribute("data-href")
      runCommand(() => navigate(href))
    }
  }

  // --- search (useDocsSearch, type "fetch") ---
  let debounce = 0
  let debounce2 = 0
  let reqSeq = 0
  let loading = false
  let data = "empty"
  let lastSearch = ""

  function setLoading(on) {
    loading = on
    spinner.hidden = !on
    empty.textContent = on ? "Searching..." : "No results found."
  }

  function renderResults() {
    const box = results.querySelector("[cmdk-group-items]")
    const keep = selectedItem()
    const keepValue = keep && results.contains(keep) ? keep.getAttribute("data-value") : null
    box.replaceChildren()
    if (lastSearch.trim() && Array.isArray(data)) {
      const unique = data.filter(
        (item, index, self) => !(item.type === "text" && item.content.trim().split(/\s+/).length <= 1) && index === self.findIndex((t) => t.content === item.content),
      )
      for (const item of unique) {
        const el = resultTpl.content.firstElementChild.cloneNode(true)
        el.setAttribute("data-type", item.type)
        el.setAttribute("data-value", `${item.content} ${item.type}`.trim())
        el.setAttribute("data-cmdk-keywords", JSON.stringify([item.content.trim()]))
        el.setAttribute("data-href", item.url)
        el.firstElementChild.textContent = item.content
        prepare(el)
        box.appendChild(el)
      }
    }
    filter()
    const sel = selectedItem()
    if (keepValue && !sel) {
      const again = [...box.children].find((el) => el.getAttribute("data-value") === keepValue && !el.hidden)
      if (again) select(again, false)
    }
    // cmdk: a newly mounted item selects the first item only when nothing is.
    if (!selectedItem() || selectedItem().hidden) selectFirst()
  }

  async function fetchResults(q) {
    const seq = ++reqSeq
    if (!q) {
      data = "empty"
      setLoading(false)
      renderResults()
      return
    }
    setLoading(true)
    try {
      const res = await fetch(`/api/search?query=${encodeURIComponent(q)}`)
      const json = await res.json()
      if (seq !== reqSeq) return
      data = json
    } catch {
      if (seq !== reqSeq) return
    }
    setLoading(false)
    renderResults()
  }

  function onSearchChange(v) {
    clearTimeout(debounce)
    debounce = setTimeout(() => {
      lastSearch = v
      clearTimeout(debounce2)
      debounce2 = setTimeout(() => fetchResults(v), 100)
    }, 500)
  }

  // --- events ---
  const onKeyDown = (e) => {
    if (e.defaultPrevented || e.isComposing || e.keyCode === 229) return
    const next = () => (e.metaKey ? select(validItems().at(-1) ?? null) : e.altKey ? byGroup(1) : byItem(1))
    const prev = () => (e.metaKey ? select(validItems()[0] ?? null) : e.altKey ? byGroup(-1) : byItem(-1))
    switch (e.key) {
      case "n":
      case "j":
        if (e.ctrlKey) {
          e.preventDefault()
          next()
        }
        break
      case "ArrowDown":
        e.preventDefault()
        next()
        break
      case "p":
      case "k":
        if (e.ctrlKey) {
          e.preventDefault()
          prev()
        }
        break
      case "ArrowUp":
        e.preventDefault()
        prev()
        break
      case "Home":
        e.preventDefault()
        select(validItems()[0] ?? null)
        break
      case "End":
        e.preventDefault()
        select(validItems().at(-1) ?? null)
        break
      case "Enter":
        e.preventDefault()
        run(selectedItem())
        break
    }
  }
  const onPointerMove = (e) => {
    const el = e.target instanceof Element ? e.target.closest(ITEM) : null
    if (el && root.contains(el) && el.getAttribute("aria-disabled") !== "true") select(el, false)
  }
  const onClick = (e) => {
    const el = e.target instanceof Element ? e.target.closest(ITEM) : null
    if (el && root.contains(el)) run(el)
  }
  const onInput = () => {
    search = input.value
    input.setAttribute("value", search)
    filter()
    selectFirst()
    onSearchChange(search)
  }

  root.addEventListener("keydown", onKeyDown)
  root.addEventListener("pointermove", onPointerMove)
  root.addEventListener("click", onClick)
  input.addEventListener("input", onInput)

  filter()
  renderFooter()
  selectFirst()

  return () => {
    clearTimeout(debounce)
    clearTimeout(debounce2)
    reqSeq++
    root.removeEventListener("keydown", onKeyDown)
    root.removeEventListener("pointermove", onPointerMove)
    root.removeEventListener("click", onClick)
    input.removeEventListener("input", onInput)
  }
}

define("command-menu", { selector: "[data-command-menu]", mount })
