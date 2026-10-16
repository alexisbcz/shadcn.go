// select-autocomplete: React Aria's Select + Autocomplete over the markup of
// aria.go. The popover (a template) is mounted at the end of <body> below
// the trigger: 4px off, as wide as the trigger, its max-height leaving 12px
// to the viewport edge (flipped above when the space below is smaller than
// the content, under 100px and smaller than the space above). The search field filters
// the options ("contains", ignoring case and accents) and keeps its text
// while closed (Autocomplete wraps the popover); arrows move the virtual
// focus (data-focused), Enter or a click picks the option.

import { define } from "./core.js"

const OFFSET = 4
const PADDING = 12

const normalize = (s) => s.normalize("NFD").replace(/\p{Diacritic}/gu, "").toLocaleLowerCase()

define("example-select-autocomplete", {
  selector: "[data-select-autocomplete]",
  mount(root) {
    const trigger = root.querySelector('[data-slot="select-trigger"]')
    const value = trigger.querySelector('[data-slot="select-value"]')
    const tpl = root.querySelector("template[data-select-autocomplete-popover]")
    const check = root.querySelector("template[data-select-autocomplete-check]")
    let popup = null
    let query = ""
    let selected = null
    let focused = null

    const items = () => [...popup.querySelectorAll('[data-slot="select-item"]')]
    const visible = () => items().filter((i) => !i.hidden)

    function setFocused(item) {
      for (const i of items()) {
        i.removeAttribute("data-focused")
        i.removeAttribute("data-focus-visible")
      }
      focused = item
      if (item) {
        item.setAttribute("data-focused", "true")
        item.setAttribute("data-focus-visible", "true")
        item.scrollIntoView({ block: "nearest" })
      }
    }

    function filter() {
      const q = normalize(query)
      let shown = 0
      for (const i of items()) {
        const match = normalize(i.textContent).includes(q)
        i.hidden = !match
        if (match) shown++
      }
      const list = popup.querySelector('[data-slot="select-list"]')
      const group = popup.querySelector('[data-slot="select-group"]')
      group.hidden = shown === 0
      if (shown === 0) list.setAttribute("data-empty", "true")
      else list.removeAttribute("data-empty")
      const wrapper = popup.querySelector('[data-slot="select-input-wrapper"]')
      if (query) wrapper.removeAttribute("data-empty")
      else wrapper.setAttribute("data-empty", "true")
      if (focused?.hidden) setFocused(null)
    }

    function position() {
      const r = trigger.getBoundingClientRect()
      const vh = document.documentElement.clientHeight
      popup.style.position = "absolute"
      popup.style.zIndex = "100000"
      popup.style.setProperty("--trigger-anchor-point", "0px 0px")
      popup.style.setProperty("--trigger-width", `${r.width}px`)
      popup.style.left = `${Math.round(r.left + scrollX)}px`
      const below = vh - r.bottom - OFFSET - PADDING
      const above = r.top - OFFSET - PADDING
      popup.style.maxHeight = "none"
      const natural = popup.scrollHeight
      if (below < natural && below < 100 && above > below) {
        popup.dataset.placement = "top"
        const max = Math.round(above)
        popup.style.maxHeight = `${max}px`
        const h = Math.min(natural, max)
        popup.style.top = `${Math.round(r.top + scrollY - OFFSET - h)}px`
      } else {
        popup.dataset.placement = "bottom"
        popup.style.maxHeight = `${Math.round(below)}px`
        popup.style.top = `${Math.round(r.bottom + scrollY + OFFSET)}px`
      }
    }

    function open() {
      if (popup) return
      popup = tpl.content.firstElementChild.cloneNode(true)
      const input = popup.querySelector("input")
      input.value = query
      for (const i of items()) {
        const on = i.dataset.key === selected
        i.setAttribute("aria-selected", String(on))
        const ind = i.lastElementChild
        ind.replaceChildren()
        if (on) ind.append(check.content.firstElementChild.cloneNode(true))
      }
      document.body.append(popup)
      filter()
      position()
      popup.setAttribute("data-entering", "")
      requestAnimationFrame(() =>
        Promise.all(popup.getAnimations().map((a) => a.finished.catch(() => {}))).then(() => popup?.removeAttribute("data-entering")),
      )
      trigger.setAttribute("aria-expanded", "true")
      input.focus({ preventScroll: true })
      input.addEventListener("input", () => {
        query = input.value
        filter()
        setFocused(null)
      })
      input.addEventListener("keydown", (e) => {
        const list = visible()
        const idx = list.indexOf(focused)
        if (e.key === "ArrowDown") {
          e.preventDefault()
          setFocused(list[Math.min(idx + 1, list.length - 1)] ?? null)
        } else if (e.key === "ArrowUp") {
          e.preventDefault()
          setFocused(list[idx <= 0 ? 0 : idx - 1] ?? null)
        } else if (e.key === "Enter") {
          e.preventDefault()
          if (focused) pick(focused)
        } else if (e.key === "Escape") {
          e.preventDefault()
          close()
        }
      })
      popup.addEventListener("click", (e) => {
        const item = e.target.closest('[data-slot="select-item"]')
        if (item) pick(item)
      })
      popup.addEventListener("pointermove", (e) => {
        const item = e.target.closest('[data-slot="select-item"]')
        for (const i of items()) i.toggleAttribute("data-hovered", i === item)
      })
    }

    function close() {
      if (!popup) return
      const p = popup
      popup = null
      focused = null
      trigger.setAttribute("aria-expanded", "false")
      p.setAttribute("data-exiting", "")
      requestAnimationFrame(() => Promise.all(p.getAnimations().map((a) => a.finished.catch(() => {}))).then(() => p.remove()))
      trigger.focus({ preventScroll: true })
    }

    function pick(item) {
      selected = item.dataset.key
      value.removeAttribute("data-placeholder")
      value.replaceChildren(...[...item.children].map((c) => c.cloneNode(true)))
      value.lastElementChild?.replaceChildren()
      close()
    }

    const onTrigger = () => (popup ? close() : open())
    const onOutside = (e) => {
      if (popup && !popup.contains(e.target) && !trigger.contains(e.target)) close()
    }
    trigger.addEventListener("click", onTrigger)
    document.addEventListener("pointerdown", onOutside, true)
    return () => {
      trigger.removeEventListener("click", onTrigger)
      document.removeEventListener("pointerdown", onOutside, true)
      popup?.remove()
    }
  },
})
