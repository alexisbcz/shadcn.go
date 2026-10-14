// dashboard-01's behavior: the chart's time range (toggle group and select
// kept in sync, data filtered like upstream), and the data table's state
// as TanStack Table and dnd-kit keep it upstream: row selection, column
// visibility, pagination and drag-and-drop row reordering (mouse and
// keyboard sensors, vertical list sorting strategy, closest center
// collision detection, restricted to the vertical axis).
import { define } from "./core.js"
import { getSelect } from "./select.js"
import { toast } from "./sonner.js"

const setAttr = (el, name, v) => {
  if (v === null || v === undefined) el.removeAttribute(name)
  else if (el.getAttribute(name) !== v) el.setAttribute(name, v)
}

// --- chart ----------------------------------------------------------------------

function filterRange(data, timeRange) {
  return data.filter((item) => {
    const date = new Date(item.date)
    const referenceDate = new Date("2024-06-30")
    let daysToSubtract = 90
    if (timeRange === "30d") {
      daysToSubtract = 30
    } else if (timeRange === "7d") {
      daysToSubtract = 7
    }
    const startDate = new Date(referenceDate)
    startDate.setDate(startDate.getDate() - daysToSubtract)
    return date >= startDate
  })
}

define("example-dashboard-chart", {
  selector: "[data-dashboard-chart]",
  mount(card) {
    const action = card.querySelector('[data-slot="card-action"]')
    const group = action?.querySelector('[data-slot="toggle-group"]')
    const trigger = action?.querySelector('[data-slot="select-trigger"]')
    if (!group || !trigger) return
    const items = () => [...group.querySelectorAll('[data-slot="toggle-group-item"]')]
    let all = null
    let syncing = false
    const apply = (range) => {
      const chart = card.querySelector('[data-slot="chart"]')?.chart
      chart?.update((spec) => {
        all ??= spec.props.data
        spec.props.data = filterRange(all, range)
      })
    }
    // data-dashboard-chart names the initial range of examples/dashboard,
    // which has no isMobile effect; empty is dashboard-01 (90 days).
    const initial = card.getAttribute("data-dashboard-chart")
    let range = "90d"
    const renderToggles = () => {
      for (const it of items()) {
        const on = it.getAttribute("data-value") === range
        setAttr(it, "aria-checked", String(on))
        setAttr(it, "data-state", on ? "on" : "off")
      }
    }
    const renderSelect = () => {
      const s = getSelect(trigger)
      if (!s || s.value === range) return
      syncing = true
      try {
        s.setValue(range)
      } finally {
        syncing = false
      }
    }
    const setRange = (next) => {
      if (next === range) return
      range = next
      renderToggles()
      renderSelect()
      apply(range)
    }
    const onChange = (event) => {
      // Handled here: nothing above the card's action sees the change.
      event.stopPropagation()
      if (syncing) return
      if (group.contains(event.target)) setRange(event.detail?.pressedItems?.[0]?.getAttribute("data-value") ?? "")
      else setRange(event.detail?.value ?? "")
    }
    action.addEventListener("shadcn:value-change", onChange)
    if (initial) {
      range = initial
      apply(range)
      return () => action.removeEventListener("shadcn:value-change", onChange)
    }
    // useIsMobile: on a narrow viewport the range becomes the last 7 days.
    const mql = window.matchMedia("(max-width: 767px)")
    const onMedia = () => {
      if (window.innerWidth < 768) setRange("7d")
    }
    mql.addEventListener("change", onMedia)
    onMedia()
    return () => {
      action.removeEventListener("shadcn:value-change", onChange)
      mql.removeEventListener("change", onMedia)
    }
  },
})

// --- data table -------------------------------------------------------------------

const HIDEABLE = ["type", "status", "target", "limit", "reviewer"]

/** Places nodes as parent's children, in order, moving only what must move. */
function reconcile(parent, nodes) {
  for (const n of [...parent.children]) if (!nodes.includes(n)) n.remove()
  let prev = null
  for (const n of nodes) {
    const want = prev ? prev.nextElementSibling : parent.firstElementChild
    if (want !== n) parent.insertBefore(n, want)
    prev = n
  }
}

const tableStates = new WeakMap()

define("example-dashboard-table", {
  // Mounted on the outline panel (re-created by tabs.js when the tab is
  // selected again); the table state lives above it, like DataTable's
  // state above <Tabs> upstream.
  selector: "[data-dashboard-table] template[data-rows]",
  mount(tpl) {
    const tabsRoot = tpl.closest("[data-dashboard-table]")
    const root = tpl.parentElement
    const body = root.querySelector('[data-slot="table-body"]')
    const headRow = root.querySelector('[data-slot="table-header"] > tr')
    if (!tpl || !body || !headRow) return
    const cleanups = []
    const on = (el, type, fn, opts) => {
      el.addEventListener(type, fn, opts)
      cleanups.push(() => el.removeEventListener(type, fn, opts))
    }

    const protos = new Map()
    for (const tr of tpl.content.children) if (tr.tagName === "TR") protos.set(tr.getAttribute("data-row-id"), tr)
    const indicatorProto = tpl.content.querySelector('[data-proto="checked"] [data-slot="checkbox-indicator"]')
    const live = new Map()
    for (const tr of body.children) live.set(tr.getAttribute("data-row-id"), tr)
    const selectAll = headRow.querySelector('[data-slot="checkbox"]')
    const columnsTrigger = tabsRoot.querySelector("[data-columns-trigger]")
    const selectedText = root.querySelector("[data-selected-text]")
    const pageText = root.querySelector("[data-page-text]")
    const pageButtons = Object.fromEntries([...root.querySelectorAll("[data-page]")].map((b) => [b.getAttribute("data-page"), b]))
    const pageSizeTrigger = root.querySelector("#rows-per-page")

    const state = tableStates.get(tabsRoot) ?? {
      order: (tpl.getAttribute("data-rows") ?? "").split(" ").filter(Boolean),
      selection: new Set(),
      hidden: new Set(),
      pageIndex: 0,
      pageSize: 10,
    }
    tableStates.set(tabsRoot, state)
    const pageCount = () => Math.ceil(state.order.length / state.pageSize)
    const pageIds = () => state.order.slice(state.pageIndex * state.pageSize, (state.pageIndex + 1) * state.pageSize)
    const rowEl = (id) => {
      let tr = live.get(id)
      if (!tr) {
        tr = protos.get(id).cloneNode(true)
        live.set(id, tr)
      }
      return tr
    }

    // Checkboxes are fully controlled (TanStack state): their clicks are
    // taken before checkbox.js and their state rendered here, like React.
    const renderCheckbox = (el, want) => {
      if (!el) return
      setAttr(el, "aria-checked", want === "indeterminate" ? "mixed" : String(want === "checked"))
      setAttr(el, "data-state", want)
      let ind = el.querySelector(':scope > [data-slot="checkbox-indicator"]')
      if (want === "unchecked") ind?.remove()
      else {
        if (!ind && indicatorProto) {
          ind = indicatorProto.cloneNode(true)
          el.appendChild(ind)
        }
        if (ind) setAttr(ind, "data-state", want)
      }
    }

    const render = () => {
      const ids = pageIds()
      const rows = ids.map(rowEl)
      for (const [id, tr] of live) {
        if (!ids.includes(id)) {
          tr.remove()
          live.delete(id)
        }
      }
      reconcile(body, rows)
      for (const [i, tr] of rows.entries()) {
        const selected = state.selection.has(ids[i])
        setAttr(tr, "data-state", selected ? "selected" : "false")
        renderCheckbox(tr.querySelector('[data-col="select"] [data-slot="checkbox"]'), selected ? "checked" : "unchecked")
      }
      for (const cell of root.querySelectorAll("[data-col]")) {
        const hide = state.hidden.has(cell.getAttribute("data-col"))
        if (cell.hidden !== hide) cell.hidden = hide
      }
      const all = ids.length > 0 && ids.every((id) => state.selection.has(id))
      const some = ids.some((id) => state.selection.has(id))
      renderCheckbox(selectAll, all ? "checked" : some ? "indeterminate" : "unchecked")
      if (selectedText) selectedText.textContent = `${state.selection.size} of ${state.order.length} row(s) selected.`
      const sizeSelect = pageSizeTrigger && getSelect(pageSizeTrigger)
      if (sizeSelect && sizeSelect.value !== String(state.pageSize)) sizeSelect.setValue(String(state.pageSize))
      if (pageText) pageText.textContent = `Page ${state.pageIndex + 1} of ${pageCount()}`
      const canPrev = state.pageIndex > 0
      const canNext = state.pageIndex < pageCount() - 1
      for (const [k, b] of Object.entries(pageButtons)) {
        const disabled = k === "first" || k === "previous" ? !canPrev : !canNext
        if (b.disabled !== disabled) b.disabled = disabled
      }
    }

    // Selection.
    on(
      root,
      "click",
      (event) => {
        const el = event.target instanceof Element && event.target.closest('[data-slot="checkbox"]')
        if (!el || !body.parentElement.contains(el)) return
        event.preventDefault()
        if (el === selectAll) {
          const cur = el.getAttribute("data-state")
          const checked = cur === "indeterminate" ? true : cur !== "checked"
          for (const id of pageIds()) {
            if (checked) state.selection.add(id)
            else state.selection.delete(id)
          }
        } else {
          const id = el.closest("tr")?.getAttribute("data-row-id")
          if (!id) return
          if (state.selection.has(id)) state.selection.delete(id)
          else state.selection.add(id)
        }
        render()
      },
      true,
    )

    // Pagination.
    on(root, "click", (event) => {
      const b = event.target instanceof Element && event.target.closest("[data-page]")
      if (!b || b.disabled) return
      const k = b.getAttribute("data-page")
      if (k === "first") state.pageIndex = 0
      else if (k === "previous") state.pageIndex = Math.max(0, state.pageIndex - 1)
      else if (k === "next") state.pageIndex = Math.min(pageCount() - 1, state.pageIndex + 1)
      else if (k === "last") state.pageIndex = pageCount() - 1
      render()
    })
    on(root, "shadcn:value-change", (event) => {
      const t = event.target
      if (!(t instanceof Element) || !pageSizeTrigger) return
      const sel = getSelect(pageSizeTrigger)
      if (!sel || getSelect(t) !== sel) return
      const size = Math.max(1, Number(event.detail?.value) || 10)
      const top = state.pageSize * state.pageIndex
      state.pageSize = size
      state.pageIndex = Math.floor(top / size)
      render()
    })

    // Column visibility: the Customize Columns menu (portaled).
    on(document, "shadcn:checked-change", (event) => {
      const item = event.target
      if (!(item instanceof Element) || !item.matches('[data-slot="dropdown-menu-checkbox-item"]')) return
      const content = item.closest('[data-slot="dropdown-menu-content"]')
      if (!content || !columnsTrigger || content.id !== columnsTrigger.getAttribute("aria-controls")) return
      const id = item.textContent.trim()
      if (!HIDEABLE.includes(id)) return
      if (event.detail?.checked) state.hidden.delete(id)
      else state.hidden.add(id)
      queueMicrotask(render)
    })

    // Target and limit cells: submitting shows a promise toast.
    on(root, "submit", (event) => {
      const form = event.target
      if (!(form instanceof HTMLFormElement) || !form.hasAttribute("data-saving")) return
      event.preventDefault()
      toast.promise(new Promise((resolve) => setTimeout(resolve, 1000)), {
        loading: form.getAttribute("data-saving"),
        success: "Done",
        error: "Error",
      })
    })

    // --- drag and drop (dnd-kit) ---
    let drag = null
    const TRANSITION = "transform 200ms ease"
    const liveRegion = root.querySelector('[id^="DndLiveRegion"]')
    const announce = (text) => {
      if (liveRegion) liveRegion.textContent = text
    }
    let lastOver = null
    const center = (r) => ({ x: r.left + r.width / 2, y: r.top + r.height / 2 })
    const startDrag = (handle, mode, clientY) => {
      const tr = handle.closest("tr")
      const ids = pageIds()
      const rows = ids.map((id) => live.get(id))
      const activeIndex = rows.indexOf(tr)
      if (activeIndex < 0) return
      drag = {
        mode,
        handle,
        tr,
        ids,
        rows,
        activeIndex,
        overIndex: activeIndex,
        rects: rows.map((r) => r.getBoundingClientRect()),
        startY: clientY,
        dy: 0,
      }
      document.getSelection()?.removeAllRanges()
      lastOver = null
      announce(`Picked up draggable item ${ids[activeIndex]}.`)
      setAttr(tr, "data-dragging", "true")
      setAttr(handle, "aria-pressed", "true")
      layout()
    }
    const layout = () => {
      const d = drag
      const a = d.rects[d.activeIndex]
      const moved = { left: a.left, top: a.top + d.dy, width: a.width, height: a.height }
      const c = center(moved)
      let best = -1
      let bestDist = Infinity
      d.rects.forEach((r, i) => {
        const rc = center(r)
        const dist = Math.hypot(rc.x - c.x, rc.y - c.y)
        if (dist < bestDist) {
          bestDist = dist
          best = i
        }
      })
      d.overIndex = best < 0 ? d.activeIndex : best
      if (lastOver !== d.overIndex) {
        lastOver = d.overIndex
        announce(`Draggable item ${d.ids[d.activeIndex]} was moved over droppable area ${d.ids[d.overIndex]}.`)
      }
      d.rows.forEach((row, i) => {
        let y = 0
        if (i === d.activeIndex) y = d.dy
        else if (i > d.activeIndex && i <= d.overIndex) y = -a.height
        else if (i < d.activeIndex && i >= d.overIndex) y = a.height
        row.style.transform = `translate3d(0px, ${y}px, 0) scaleX(1) scaleY(1)`
        row.style.transition = i === d.activeIndex && d.mode === "mouse" ? "" : TRANSITION
      })
    }
    const endDrag = (drop) => {
      const d = drag
      drag = null
      if (!d) return
      for (const row of d.rows) {
        row.style.transform = ""
        row.style.transition = ""
        if (!row.getAttribute("style")) row.removeAttribute("style")
      }
      setAttr(d.tr, "data-dragging", "false")
      d.handle.removeAttribute("aria-pressed")
      const activeId = d.ids[d.activeIndex]
      announce(
        drop
          ? `Draggable item ${activeId} was dropped over droppable area ${d.ids[d.overIndex]}`
          : `Dragging was cancelled. Draggable item ${activeId} was dropped.`,
      )
      if (drop && d.overIndex !== d.activeIndex) {
        const from = state.order.indexOf(d.ids[d.activeIndex])
        const to = state.order.indexOf(d.ids[d.overIndex])
        const [moved] = state.order.splice(from, 1)
        state.order.splice(to, 0, moved)
        // TanStack resets the page index when the data changes.
        state.pageIndex = 0
        render()
      }
    }
    on(body, "mousedown", (event) => {
      const handle = event.button === 0 && event.target instanceof Element && event.target.closest("[data-drag-handle]")
      if (!handle || drag) return
      startDrag(handle, "mouse", event.clientY)
    })
    on(document, "mousemove", (event) => {
      if (!drag || drag.mode !== "mouse") return
      drag.dy = event.clientY - drag.startY
      layout()
    })
    on(document, "mouseup", () => {
      if (drag?.mode === "mouse") endDrag(true)
    })
    on(body, "keydown", (event) => {
      const handle = event.target instanceof Element && event.target.closest("[data-drag-handle]")
      if (drag?.mode === "keyboard") {
        if (event.key === " " || event.key === "Enter") {
          event.preventDefault()
          endDrag(true)
        } else if (event.key === "Escape") {
          event.preventDefault()
          endDrag(false)
        } else if (event.key === "ArrowDown" || event.key === "ArrowUp") {
          event.preventDefault()
          drag.dy += event.key === "ArrowDown" ? 25 : -25
          layout()
        }
        return
      }
      if (handle && (event.key === " " || event.key === "Enter") && !drag) {
        event.preventDefault()
        startDrag(handle, "keyboard", 0)
      }
    })

    render()
    return () => {
      cleanups.forEach((fn) => fn())
      if (drag) endDrag(false)
    }
  },
})
