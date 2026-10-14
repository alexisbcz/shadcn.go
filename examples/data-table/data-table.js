// data-table-demo's behavior: TanStack Table's state and row models
// (filtering with includesString, sorting with the alphanumeric sort,
// column visibility, row selection, pagination) for the server-rendered
// table, rendered the way React reconciles it (rows keyed by id, cells by
// column).
import { define } from "./core.js"

// Sets an attribute only when it changes (a write repaints even when equal).
const setAttr = (el, name, v) => el.getAttribute(name) !== v && el.setAttribute(name, v)

const COLUMNS = ["select", "status", "email", "amount", "actions"]
const HIDEABLE = ["status", "email", "amount"]
const PAGE_SIZE = 10

// TanStack Table's sortFn_alphanumeric.
const reSplitAlphaNumeric = /([0-9]+)/gm
function compareAlphanumeric(aStr, bStr) {
  const a = aStr.split(reSplitAlphaNumeric).filter(Boolean)
  const b = bStr.split(reSplitAlphaNumeric).filter(Boolean)
  while (a.length && b.length) {
    const aa = a.shift()
    const bb = b.shift()
    const an = parseInt(aa, 10)
    const bn = parseInt(bb, 10)
    const combo = [an, bn].sort()
    if (isNaN(combo[0])) {
      if (aa > bb) return 1
      if (bb > aa) return -1
      continue
    }
    if (isNaN(combo[1])) return isNaN(an) ? -1 : 1
    if (an > bn) return 1
    if (bn > an) return -1
  }
  return a.length - b.length
}
const alphanumeric = (a, b) => compareAlphanumeric(String(a).toLowerCase(), String(b).toLowerCase())

/**
 * Places nodes as parent's children before ref, moving only the nodes
 * React's keyed reconciliation moves (an existing node stays unless it
 * comes before the last node placed without moving).
 */
function reconcile(parent, nodes, ref = null) {
  const current = [...parent.children].filter((n) => !ref || n !== ref)
  for (const n of current) if (!nodes.includes(n)) n.remove()
  const old = [...parent.children]
  let lastPlaced = -1
  const place = new Set()
  for (const n of nodes) {
    const i = old.indexOf(n)
    if (i < 0 || i < lastPlaced) place.add(n)
    else lastPlaced = i
  }
  let next = ref
  for (let i = nodes.length - 1; i >= 0; i--) {
    const n = nodes[i]
    if (place.has(n)) parent.insertBefore(n, next)
    next = n
  }
}

define("example-data-table", {
  selector: "template[data-data-table]",
  mount(tpl) {
    const root = tpl.previousElementSibling
    if (!root) return
    const cleanups = []
    const on = (el, type, fn, opts) => {
      el.addEventListener(type, fn, opts)
      cleanups.push(() => el.removeEventListener(type, fn, opts))
    }
    const filter = root.querySelector('[data-slot="input"]')
    const columnsTrigger = root.querySelector('[data-slot="dropdown-menu-trigger"]')
    const headRow = root.querySelector('[data-slot="table-header"] > tr')
    const body = root.querySelector('[data-slot="table-body"]')
    const footerText = root.querySelector(".flex-1.text-sm")
    const noResults = tpl.content.querySelector("tr")
    const protos = {}
    for (const el of tpl.content.querySelectorAll("[data-proto]")) protos[el.getAttribute("data-proto")] = el.firstElementChild

    let headCells = [...headRow.children]
    const selectAll = headCells[0].querySelector('[data-slot="checkbox"]')
    const paymentIds = (tpl.getAttribute("data-ids") ?? "").split(" ")
    const rows = [...body.children].map((tr, index) => ({
      id: String(index),
      paymentId: paymentIds[index],
      actions: tr.children[4].querySelector('[data-slot="dropdown-menu-trigger"]'),
      tr,
      cells: [...tr.children],
      email: tr.children[2].textContent,
      checkbox: tr.children[0].querySelector('[data-slot="checkbox"]'),
    }))

    const state = { sorting: null, filter: "", visible: { status: true, email: true, amount: true }, selection: {} }
    let muted = false

    const filteredRows = () => {
      if (!state.filter) return rows
      const search = state.filter.toLowerCase()
      return rows.filter((r) => r.email.toLowerCase().includes(search))
    }
    const sortedRows = () => {
      const list = [...filteredRows()]
      if (!state.sorting) return list
      return list.sort((a, b) => {
        let s = alphanumeric(a.email, b.email)
        if (s !== 0) return state.sorting === "desc" ? -s : s
        return Number(a.id) - Number(b.id)
      })
    }
    const pageRows = () => sortedRows().slice(0, PAGE_SIZE)
    const isAllPageSelected = () => {
      const page = pageRows()
      return page.length > 0 && page.every((r) => state.selection[r.id])
    }
    const isSomePageSelected = () => !isAllPageSelected() && pageRows().some((r) => state.selection[r.id])

    const stateOf = (el) => el.getAttribute("data-state")
    /** Drives a row's Radix checkbox to checked/unchecked, by clicking it. */
    const setCheckbox = (el, want) => {
      if (stateOf(el) === want) return el
      muted = true
      try {
        el.click()
      } finally {
        muted = false
      }
      return el
    }

    // The header checkbox is fully controlled here (React's checked prop can
    // be "indeterminate", which no click reaches): its clicks are taken
    // before checkbox.js, and its state is rendered in place like React.
    const indicatorProto = protos.checked.querySelector('[data-slot="checkbox-indicator"]')
    const renderSelectAll = (want) => {
      const el = selectAll
      setAttr(el, "aria-checked", want === "indeterminate" ? "mixed" : String(want === "checked"))
      setAttr(el, "data-state", want)
      let ind = el.querySelector(':scope > [data-slot="checkbox-indicator"]')
      if (want === "unchecked") ind?.remove()
      else {
        if (!ind) {
          ind = indicatorProto.cloneNode(true)
          el.appendChild(ind)
        }
        setAttr(ind, "data-state", want)
      }
    }
    on(
      headRow,
      "click",
      (event) => {
        if (event.target instanceof Element && event.target.closest('[data-slot="checkbox"]') === selectAll) {
          event.preventDefault()
          const cur = stateOf(selectAll)
          const checked = cur === "indeterminate" ? true : cur !== "checked"
          for (const r of pageRows()) {
            if (checked) state.selection[r.id] = true
            else delete state.selection[r.id]
          }
          render()
        }
      },
      true,
    )

    const setTexts = (el, ...parts) => {
      const texts = [...el.childNodes].filter((n) => n.nodeType === 3)
      parts.forEach((p, i) => {
        if (p !== undefined && texts[i] && texts[i].data !== String(p)) texts[i].data = String(p)
      })
    }

    const visibleColumns = () => COLUMNS.filter((c) => !HIDEABLE.includes(c) || state.visible[c])

    const render = () => {
      const cols = visibleColumns()
      reconcile(headRow, cols.map((c) => headCells[COLUMNS.indexOf(c)]))
      const page = pageRows()
      if (page.length) {
        if (noResults.isConnected) noResults.remove()
        reconcile(body, page.map((r) => r.tr))
      } else {
        for (const r of rows) r.tr.remove()
        if (!noResults.isConnected) body.appendChild(noResults)
      }
      for (const r of page) {
        const selected = !!state.selection[r.id]
        setAttr(r.tr, "data-state", selected ? "selected" : "false")
        r.checkbox = setCheckbox(r.checkbox, selected ? "checked" : "unchecked")
        reconcile(r.tr, cols.map((c) => r.cells[COLUMNS.indexOf(c)]))
      }
      renderSelectAll(isAllPageSelected() ? "checked" : isSomePageSelected() ? "indeterminate" : "unchecked")
      const filtered = filteredRows()
      setTexts(footerText, filtered.filter((r) => state.selection[r.id]).length, undefined, undefined, filtered.length)
    }

    on(filter, "input", () => {
      state.filter = filter.value
      if (filter.defaultValue !== filter.value) filter.defaultValue = filter.value
      render()
    })
    on(headRow, "click", (event) => {
      const button = event.target instanceof Element && event.target.closest('[data-slot="button"]')
      if (button && headRow.contains(button)) {
        state.sorting = state.sorting === "asc" ? "desc" : "asc"
        render()
      }
    })
    on(root, "shadcn:checked-change", (event) => {
      if (muted) return
      const el = event.target
      if (!(el instanceof Element) || el.getAttribute("data-slot") !== "checkbox") return
      const checked = !!event.detail.checked
      // React renders after the click: let the checkbox update first.
      queueMicrotask(() => {
        const r = rows.find((x) => x.checkbox === el)
        if (!r) return
        if (checked) state.selection[r.id] = true
        else delete state.selection[r.id]
        render()
      })
    })
    // Column visibility: the Columns menu's checkbox items (portaled).
    on(document, "shadcn:checked-change", (event) => {
      const item = event.target
      if (!(item instanceof Element) || !item.matches('[data-slot="dropdown-menu-checkbox-item"]')) return
      const content = item.closest('[data-slot="dropdown-menu-content"]')
      if (!content || content.id !== columnsTrigger.getAttribute("aria-controls")) return
      const id = item.textContent.trim()
      if (!HIDEABLE.includes(id)) return
      state.visible[id] = !!event.detail.checked
      queueMicrotask(render)
    })
    // A row's "Copy payment ID" item (its menu is portaled): onClick
    // writes the row's payment id to the clipboard. The row is found by its
    // trigger's aria-controls, which the trigger has only while open.
    on(document, "click", (event) => {
      const item = event.target instanceof Element && event.target.closest('[data-slot="dropdown-menu-item"]')
      if (!item || item.textContent.trim() !== "Copy payment ID") return
      const content = item.closest('[data-slot="dropdown-menu-content"]')
      const r = content && rows.find((x) => x.actions?.getAttribute("aria-controls") === content.id)
      if (r) navigator.clipboard.writeText(r.paymentId)
    }, true) // capture: before the item's selection closes the menu
    return () => cleanups.forEach((fn) => fn())
  },
})
