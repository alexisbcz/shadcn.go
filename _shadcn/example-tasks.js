// The tasks example's table (app/(app)/examples/tasks): TanStack Table's
// state and row models for the server-rendered page — the title filter
// (includesString), the status and priority faceted filters, sorting from
// the column headers' menus, column visibility, row selection and
// pagination — applied to the DOM the way React reconciles it (rows keyed
// by index, cells by column). Every row is prerendered: the rows past the
// first page wait in the template.
import { define } from "./core.js"
import { getCommand } from "./command.js"
import { setChecked } from "./menu.js"

const setAttr = (el, name, v) => el.getAttribute(name) !== v && el.setAttribute(name, v)

const COLUMNS = ["select", "id", "title", "status", "priority", "actions"]
const HIDEABLE = ["title", "status", "priority"]
const SORTABLE = { 2: "title", 3: "status", 4: "priority" } // header cell index -> column
const OPTIONS = {
  status: ["backlog", "todo", "in progress", "done", "canceled"],
  priority: ["low", "medium", "high"],
}
const LABELS = {
  backlog: "Backlog", todo: "Todo", "in progress": "In Progress", done: "Done", canceled: "Canceled",
  low: "Low", medium: "Medium", high: "High",
}

// TanStack Table's sortFn_alphanumeric and sortFn_text.
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
const alphanumeric = (a, b) => compareAlphanumeric(a.toLowerCase(), b.toLowerCase())
const text = (a, b) => {
  a = a.toLowerCase()
  b = b.toLowerCase()
  return a === b ? 0 : a > b ? 1 : -1
}

/** Places nodes as parent's children before ref, moving only what React moves. */
function reconcile(parent, nodes, ref = null) {
  for (const n of [...parent.children]) if (n !== ref && !nodes.includes(n)) n.remove()
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

define("example-tasks", {
  selector: "template[data-tasks-data]",
  mount(tpl) {
    const root = tpl.closest("[data-tasks-table]")
    if (!root) return
    const cleanups = []
    const on = (el, type, fn, opts) => {
      el.addEventListener(type, fn, opts)
      cleanups.push(() => el.removeEventListener(type, fn, opts))
    }
    const protos = {}
    for (const el of tpl.content.querySelectorAll("[data-proto]")) protos[el.getAttribute("data-proto")] = el.firstElementChild
    const proto = (name) => protos[name].cloneNode(true)
    const tasks = JSON.parse(tpl.getAttribute("data-tasks-data"))

    const filterInput = root.querySelector('[data-slot="input"]')
    const filterBar = filterInput.parentElement
    const headRow = root.querySelector('[data-slot="table-header"] > tr')
    const body = root.querySelector('[data-slot="table-body"]')
    const selectedText = root.querySelector("[data-tasks-selected]")
    const pageText = root.querySelector("[data-tasks-page]")
    const pageButtons = [...pageText.nextElementSibling.querySelectorAll('[data-slot="button"]')]
    const viewTrigger = root.querySelector("[data-tasks-view]")
    const headCells = [...headRow.children]
    const selectAll = headCells[0].querySelector('[data-slot="checkbox"]')
    const noResults = protos["no-results"]
    const sortTriggers = {}
    for (const [i, col] of Object.entries(SORTABLE)) sortTriggers[col] = headCells[i].querySelector('[data-slot="dropdown-menu-trigger"]')
    const filterTriggers = {}
    for (const el of root.querySelectorAll("[data-tasks-filter]")) filterTriggers[el.getAttribute("data-tasks-filter")] = el

    const trs = [...body.children, ...tpl.content.querySelectorAll("tbody:not([data-proto]) > tr")]
    const rows = tasks.map((task, index) => ({
      index,
      task,
      tr: trs[index],
      cells: [...trs[index].children],
      get checkbox() {
        return this.cells[0].querySelector('[data-slot="checkbox"]')
      },
    }))

    const state = {
      title: "",
      filters: { status: [], priority: [] },
      sorting: null, // { id, desc }
      visible: { title: true, status: true, priority: true },
      selection: new Set(),
      pageIndex: 0,
      pageSize: 25,
    }
    let muted = false

    // --- row models ---
    const passes = (r, except) => {
      if (except !== "title" && state.title && !r.task.title.toLowerCase().includes(state.title.toLowerCase())) return false
      for (const col of ["status", "priority"]) {
        if (col !== except && state.filters[col].length && !state.filters[col].includes(r.task[col])) return false
      }
      return true
    }
    const filteredRows = () => rows.filter((r) => passes(r))
    const sortedRows = () => {
      const list = filteredRows()
      if (!state.sorting) return list
      const { id, desc } = state.sorting
      // sortFn "auto": alphanumeric when one of the first 10 filtered values has digits.
      const fn = list.slice(0, 10).some((r) => r.task[id].split(reSplitAlphaNumeric).length > 1) ? alphanumeric : text
      return [...list].sort((a, b) => {
        const s = fn(a.task[id], b.task[id])
        if (s !== 0) return desc ? -s : s
        return a.index - b.index
      })
    }
    const pageCount = () => Math.ceil(filteredRows().length / state.pageSize)
    const pageRows = () => sortedRows().slice(state.pageIndex * state.pageSize, (state.pageIndex + 1) * state.pageSize)
    const facets = (col) => {
      const counts = new Map()
      for (const r of rows) if (passes(r, col)) counts.set(r.task[col], (counts.get(r.task[col]) ?? 0) + 1)
      return counts
    }

    // --- checkboxes ---
    const stateOf = (el) => el.getAttribute("data-state")
    const setCheckbox = (el, want) => {
      if (!el || stateOf(el) === want) return
      muted = true
      try {
        el.click()
      } finally {
        muted = false
      }
    }
    const indicatorProto = protos.checked.querySelector('[data-slot="checkbox-indicator"]')
    const renderSelectAll = (want) => {
      setAttr(selectAll, "aria-checked", want === "indeterminate" ? "mixed" : String(want === "checked"))
      setAttr(selectAll, "data-state", want)
      let ind = selectAll.querySelector(':scope > [data-slot="checkbox-indicator"]')
      if (want === "unchecked") ind?.remove()
      else {
        if (!ind) {
          ind = indicatorProto.cloneNode(true)
          selectAll.appendChild(ind)
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
            if (checked) state.selection.add(r.index)
            else state.selection.delete(r.index)
          }
          render()
        }
      },
      true,
    )
    on(root, "shadcn:checked-change", (event) => {
      if (muted) return
      const el = event.target
      if (!(el instanceof Element) || el.getAttribute("data-slot") !== "checkbox" || el === selectAll) return
      const checked = !!event.detail.checked
      queueMicrotask(() => {
        const r = rows.find((x) => x.checkbox === el)
        if (!r) return
        if (checked) state.selection.add(r.index)
        else state.selection.delete(r.index)
        render()
      })
    })

    const setTexts = (el, ...parts) => {
      const texts = [...el.childNodes].filter((n) => n.nodeType === 3)
      parts.forEach((p, i) => {
        if (p !== undefined && texts[i] && texts[i].data !== String(p)) texts[i].data = String(p)
      })
    }

    // --- toolbar: the faceted filters' triggers and lists, Reset ---
    let resetButton = null
    const popoverTemplate = (trigger) =>
      document.querySelector(`template[data-content-id="${CSS.escape(trigger.getAttribute("aria-controls") ?? "")}"]`) ??
      trigger.nextElementSibling
    const renderFilterTrigger = (col) => {
      const trigger = filterTriggers[col]
      const selected = state.filters[col]
      // Children: icon, title text, then the selection.
      while (trigger.childNodes.length > 2) trigger.lastChild.remove()
      if (!selected.length) return
      trigger.appendChild(proto("separator"))
      const count = proto("badge-count")
      count.textContent = String(selected.length)
      trigger.appendChild(count)
      const badges = proto("badges")
      if (selected.length > 2) {
        const b = proto("badge")
        b.textContent = `${selected.length} selected`
        badges.appendChild(b)
      } else {
        for (const v of OPTIONS[col]) {
          if (!selected.includes(v)) continue
          const b = proto("badge")
          b.textContent = LABELS[v]
          badges.appendChild(b)
        }
      }
      trigger.appendChild(badges)
    }
    const CHECK_ON = "flex size-4 items-center justify-center rounded-[4px] border border-primary bg-primary text-primary-foreground"
    const CHECK_OFF = "flex size-4 items-center justify-center rounded-[4px] border border-input [&_svg]:invisible"
    /** Renders a faceted filter's list (in the template, and live when open). */
    const renderFilterList = (col, list) => {
      if (!list) return false
      const counts = facets(col)
      for (const item of list.querySelectorAll("[data-tasks-option]")) {
        const v = item.getAttribute("data-tasks-option")
        setAttr(item.firstElementChild, "class", state.filters[col].includes(v) ? CHECK_ON : CHECK_OFF)
        let count = item.querySelector(":scope > .ml-auto")
        const n = counts.get(v) ?? 0
        if (!n) count?.remove()
        else {
          if (!count) {
            count = proto("count")
            item.appendChild(count)
          }
          if (count.textContent !== String(n)) count.textContent = String(n)
        }
      }
      // The "Clear filters" group, after a separator, while values are selected.
      const sizer = list.querySelector("[cmdk-list-sizer]")
      const clear = sizer.querySelector(":scope > [data-tasks-clear]")
      const want = state.filters[col].length > 0
      if (want === !!clear) return false
      if (clear) {
        clear.previousElementSibling?.remove()
        clear.remove()
      } else {
        const p = proto("clear")
        const sep = p.querySelector("[cmdk-separator]")
        const group = p.querySelector("[cmdk-group]")
        group.setAttribute("data-tasks-clear", "")
        sizer.append(sep, group)
      }
      return true
    }
    const renderFilter = (col) => {
      const trigger = filterTriggers[col]
      renderFilterTrigger(col)
      const tplEl = popoverTemplate(trigger)
      renderFilterList(col, tplEl?.content)
      const contentId = trigger.getAttribute("aria-controls")
      const live = contentId && document.getElementById(contentId)
      if (!live) return
      const cmdRoot = live.querySelector("[cmdk-root]")
      if (!cmdRoot) return
      const structural = renderFilterList(col, live)
      if (structural) {
        // cmdk registers its items on mount: remount the list's command,
        // keeping the search, the selected item and the focus.
        const cmd = getCommand(cmdRoot)
        const search = cmd?.search ?? ""
        const value = cmd?.value ?? ""
        const hadFocus = document.activeElement?.hasAttribute("cmdk-input") && cmdRoot.contains(document.activeElement)
        const fresh = tplEl.content.querySelector("[cmdk-root]").cloneNode(true)
        renderFilterList(col, fresh)
        cmdRoot.replaceWith(fresh)
        const restore = () => {
          const c = getCommand(fresh)
          if (!c) return queueMicrotask(restore)
          if (search) c.setSearch(search)
          if (value && fresh.querySelector(`[cmdk-item][data-value="${CSS.escape(value)}"]`)) c.setValue(value)
          if (hadFocus) fresh.querySelector("[cmdk-input]")?.focus()
        }
        queueMicrotask(restore)
      }
    }
    const renderReset = () => {
      const filtered = !!state.title || state.filters.status.length > 0 || state.filters.priority.length > 0
      if (filtered && !resetButton) {
        resetButton = proto("reset")
        filterBar.appendChild(resetButton)
      } else if (!filtered && resetButton) {
        resetButton.remove()
        resetButton = null
      }
    }
    on(filterBar, "click", (event) => {
      if (resetButton && event.target instanceof Element && resetButton.contains(event.target)) {
        state.title = ""
        state.filters.status = []
        state.filters.priority = []
        filterInput.value = ""
        filterInput.defaultValue = ""
        changed(true)
      }
    })
    on(document, "shadcn:select", (event) => {
      const item = event.target
      if (!(item instanceof Element) || !item.hasAttribute("cmdk-item")) return
      const content = item.closest('[data-slot="popover-content"]')
      if (!content) return
      const col = Object.keys(filterTriggers).find((c) => filterTriggers[c].getAttribute("aria-controls") === content.id)
      if (!col) return
      if (item.closest("[data-tasks-clear]")) state.filters[col] = []
      else {
        const v = item.getAttribute("data-tasks-option")
        const list = state.filters[col]
        state.filters[col] = list.includes(v) ? list.filter((x) => x !== v) : [...list, v]
      }
      changed(true)
    })
    on(filterInput, "input", () => {
      state.title = filterInput.value
      if (filterInput.defaultValue !== filterInput.value) filterInput.defaultValue = filterInput.value
      changed(true)
    })

    // --- column headers: sorting and hiding ---
    const menuTemplate = (trigger) => document.querySelector(`template[data-trigger-id="${CSS.escape(trigger.id)}"]`)
    const renderSortIcons = () => {
      for (const [col, trigger] of Object.entries(sortTriggers)) {
        const want = state.sorting?.id === col ? (state.sorting.desc ? "sort-desc" : "sort-asc") : "sort-none"
        const icon = trigger.querySelector(":scope > svg")
        if (icon?.getAttribute("data-tasks-icon") === want || (!icon?.hasAttribute("data-tasks-icon") && want === "sort-none")) continue
        const next = proto(want)
        next.setAttribute("data-tasks-icon", want)
        icon ? icon.replaceWith(next) : trigger.appendChild(next)
      }
    }
    on(
      document,
      "shadcn:select",
      (event) => {
        const item = event.target
        if (!(item instanceof Element) || item.getAttribute("data-slot") !== "dropdown-menu-item") return
        const content = item.closest('[data-slot="dropdown-menu-content"]')
        if (!content) return
        const col = Object.keys(sortTriggers).find((c) => menuTemplate(sortTriggers[c])?.getAttribute("data-content-id") === content.id)
        if (!col) return
        const action = item.textContent.trim()
        if (action === "Asc") state.sorting = { id: col, desc: false }
        else if (action === "Desc") state.sorting = { id: col, desc: true }
        else if (action === "Hide") {
          state.visible[col] = false
          syncViewMenu(col)
          queueMicrotask(render)
          return
        } else return
        queueMicrotask(() => changed(false))
      },
      true,
    )
    // The View menu's checkbox items (portaled).
    const viewKeys = () => {
      const keys = {}
      const t = menuTemplate(viewTrigger)
      for (const item of t?.content.querySelectorAll("[data-radix-menu-checkbox]") ?? []) {
        keys[item.textContent.trim()] = item.getAttribute("data-radix-menu-checkbox")
      }
      return keys
    }
    const syncViewMenu = (col) => {
      const key = viewKeys()[col]
      if (key) setChecked(key, state.visible[col])
    }
    on(document, "shadcn:checked-change", (event) => {
      const item = event.target
      if (!(item instanceof Element) || !item.matches('[data-slot="dropdown-menu-checkbox-item"]')) return
      const content = item.closest('[data-slot="dropdown-menu-content"]')
      if (!content || content.id !== menuTemplate(viewTrigger)?.getAttribute("data-content-id")) return
      const id = item.textContent.trim()
      if (!HIDEABLE.includes(id)) return
      state.visible[id] = !!event.detail.checked
      queueMicrotask(render)
    })

    // --- pagination ---
    on(root, "shadcn:value-change", (event) => {
      const t = event.target
      if (!(t instanceof Element) || !t.closest?.("[data-tasks-table]")) return
      const size = Number(event.detail.value)
      if (!size || size === state.pageSize) return
      const top = state.pageSize * state.pageIndex
      state.pageSize = size
      state.pageIndex = Math.floor(top / size)
      render()
    })
    const [first, prev, next, last] = pageButtons
    on(first, "click", () => ((state.pageIndex = 0), render()))
    on(prev, "click", () => ((state.pageIndex = Math.max(0, state.pageIndex - 1)), render()))
    on(next, "click", () => ((state.pageIndex = Math.min(pageCount() - 1, state.pageIndex + 1)), render()))
    on(last, "click", () => ((state.pageIndex = Math.max(0, pageCount() - 1)), render()))

    // --- render ---
    /** A filter or the sorting changed: the page index resets (autoResetPageIndex). */
    const changed = (filters) => {
      state.pageIndex = 0
      if (filters) {
        renderReset()
        for (const col of Object.keys(filterTriggers)) renderFilter(col)
      }
      render()
    }
    const visibleColumns = () => COLUMNS.filter((c) => !HIDEABLE.includes(c) || state.visible[c])
    const render = () => {
      const cols = visibleColumns()
      reconcile(headRow, cols.map((c) => headCells[COLUMNS.indexOf(c)]))
      renderSortIcons()
      const page = pageRows()
      if (page.length) {
        if (noResults.isConnected) noResults.remove()
        reconcile(body, page.map((r) => r.tr))
      } else {
        for (const r of rows) r.tr.remove()
        if (!noResults.isConnected) body.appendChild(noResults)
      }
      for (const r of page) {
        const selected = state.selection.has(r.index)
        setAttr(r.tr, "data-state", selected ? "selected" : "false")
        reconcile(r.tr, cols.map((c) => r.cells[COLUMNS.indexOf(c)]))
      }
      // Rows just inserted mount their checkboxes after this task.
      queueMicrotask(() => {
        for (const r of pageRows()) setCheckbox(r.checkbox, state.selection.has(r.index) ? "checked" : "unchecked")
      })
      const all = page.length > 0 && page.every((r) => state.selection.has(r.index))
      const some = !all && page.some((r) => state.selection.has(r.index))
      renderSelectAll(all ? "checked" : some ? "indeterminate" : "unchecked")
      const filtered = filteredRows()
      setTexts(selectedText, filtered.filter((r) => state.selection.has(r.index)).length, undefined, undefined, filtered.length)
      const count = pageCount()
      setTexts(pageText, undefined, state.pageIndex + 1, undefined, count)
      const canPrev = state.pageIndex > 0
      const canNext = state.pageIndex < count - 1
      first.disabled = prev.disabled = !canPrev
      next.disabled = last.disabled = !canNext
    }
    return () => cleanups.forEach((fn) => fn())
  },
})
