// command.js: cmdk 1.1.1 behavior for the server-rendered markup of
// components/ui/command (the shadcn Command).
//
// cmdk keeps its state in a store and lets React render it; this module
// keeps the same store and applies the same renders to the DOM:
//   - setState(key, value) mutates the store and requests a render;
//   - schedule(id, fn) is cmdk's useScheduleLayoutEffect: callbacks run
//     after the next render, in a Map keyed by id, iterated with
//     Map.forEach (callbacks scheduled while it runs are run by the same
//     iteration when their id is new), exactly like cmdk;
//   - commit() renders and runs the scheduled callbacks until the store is
//     stable, like React flushing a discrete event's synchronous updates.
// Rendering mirrors React's reconciliation: items that stop matching are
// removed and re-inserted before their next mounted sibling in source
// order, wherever cmdk's sorting moved that sibling.
//
// Markup contract (rendered by the Go package): [cmdk-root] with options
// data-cmdk-loop, data-cmdk-should-filter="false",
// data-cmdk-vim-bindings="false", data-cmdk-disable-pointer-selection;
// [cmdk-input], [cmdk-list] > [cmdk-list-sizer], [cmdk-group] (data-value,
// data-cmdk-force-mount), [cmdk-item] (data-value, data-cmdk-keywords
// JSON, data-cmdk-force-mount), [cmdk-separator]
// (data-cmdk-always-render), and the Empty element inside
// <template data-cmdk-empty-template>.
//
// Events: "shadcn:select" on an item (detail.value) when it is chosen,
// "shadcn:value-change" (detail.value) and "shadcn:search-change"
// (detail.search) on the root.

import { define, emit } from "./core.js"

// --- command-score (cmdk 1.1.1, chunk-NZJY6EH4) ------------------------------------------------

const SCORE_CONTINUE_MATCH = 1
const SCORE_SPACE_WORD_JUMP = 0.9
const SCORE_NON_SPACE_WORD_JUMP = 0.8
const SCORE_CHARACTER_JUMP = 0.17
const SCORE_TRANSPOSITION = 0.1
const PENALTY_SKIPPED = 0.999
const PENALTY_CASE_MISMATCH = 0.9999
const PENALTY_NOT_COMPLETE = 0.99
const IS_GAP_REGEXP = /[\\/_+.#"@[({&]/
const COUNT_GAPS_REGEXP = /[\\/_+.#"@[({&]/g
const IS_SPACE_REGEXP = /[\s-]/
const COUNT_SPACE_REGEXP = /[\s-]/g

function commandScoreInner(string, abbreviation, lowerString, lowerAbbreviation, stringIndex, abbreviationIndex, memo) {
  if (abbreviationIndex === abbreviation.length) {
    if (stringIndex === string.length) return SCORE_CONTINUE_MATCH
    return PENALTY_NOT_COMPLETE
  }
  const memoizeKey = `${stringIndex},${abbreviationIndex}`
  if (memo[memoizeKey] !== undefined) return memo[memoizeKey]
  const abbreviationChar = lowerAbbreviation.charAt(abbreviationIndex)
  let index = lowerString.indexOf(abbreviationChar, stringIndex)
  let highScore = 0
  let score, transposedScore, wordBreaks, spaceBreaks
  while (index >= 0) {
    score = commandScoreInner(string, abbreviation, lowerString, lowerAbbreviation, index + 1, abbreviationIndex + 1, memo)
    if (score > highScore) {
      if (index === stringIndex) {
        score *= SCORE_CONTINUE_MATCH
      } else if (IS_GAP_REGEXP.test(string.charAt(index - 1))) {
        score *= SCORE_NON_SPACE_WORD_JUMP
        wordBreaks = string.slice(stringIndex, index - 1).match(COUNT_GAPS_REGEXP)
        if (wordBreaks && stringIndex > 0) score *= Math.pow(PENALTY_SKIPPED, wordBreaks.length)
      } else if (IS_SPACE_REGEXP.test(string.charAt(index - 1))) {
        score *= SCORE_SPACE_WORD_JUMP
        spaceBreaks = string.slice(stringIndex, index - 1).match(COUNT_SPACE_REGEXP)
        if (spaceBreaks && stringIndex > 0) score *= Math.pow(PENALTY_SKIPPED, spaceBreaks.length)
      } else {
        score *= SCORE_CHARACTER_JUMP
        if (stringIndex > 0) score *= Math.pow(PENALTY_SKIPPED, index - stringIndex)
      }
      if (string.charAt(index) !== abbreviation.charAt(abbreviationIndex)) score *= PENALTY_CASE_MISMATCH
    }
    if (
      (score < SCORE_TRANSPOSITION && lowerString.charAt(index - 1) === lowerAbbreviation.charAt(abbreviationIndex + 1)) ||
      (lowerAbbreviation.charAt(abbreviationIndex + 1) === lowerAbbreviation.charAt(abbreviationIndex) &&
        lowerString.charAt(index - 1) !== lowerAbbreviation.charAt(abbreviationIndex))
    ) {
      transposedScore = commandScoreInner(string, abbreviation, lowerString, lowerAbbreviation, index + 1, abbreviationIndex + 2, memo)
      if (transposedScore * SCORE_TRANSPOSITION > score) score = transposedScore * SCORE_TRANSPOSITION
    }
    if (score > highScore) highScore = score
    index = lowerString.indexOf(abbreviationChar, index + 1)
  }
  memo[memoizeKey] = highScore
  return highScore
}

const formatInput = (string) => string.toLowerCase().replace(COUNT_SPACE_REGEXP, " ")

/** cmdk's default filter: scores how well search matches value (0 = no match). */
export function commandScore(string, abbreviation, aliases) {
  string = aliases && aliases.length > 0 ? `${string + " " + aliases.join(" ")}` : string
  return commandScoreInner(string, abbreviation, formatInput(string), formatInput(abbreviation), 0, 0, {})
}

// --- Command -----------------------------------------------------------------------------------

const GROUP_SELECTOR = '[cmdk-group=""]'
const GROUP_ITEMS_SELECTOR = '[cmdk-group-items=""]'
const GROUP_HEADING_SELECTOR = '[cmdk-group-heading=""]'
const ITEM_SELECTOR = '[cmdk-item=""]'
const VALID_ITEM_SELECTOR = `${ITEM_SELECTOR}:not([aria-disabled="true"])`
const SELECT_EVENT = "cmdk-item-select"
const VALUE_ATTR = "data-value"

function findNextSibling(el, selector) {
  let sibling = el.nextElementSibling
  while (sibling) {
    if (sibling.matches(selector)) return sibling
    sibling = sibling.nextElementSibling
  }
}

function findPreviousSibling(el, selector) {
  let sibling = el.previousElementSibling
  while (sibling) {
    if (sibling.matches(selector)) return sibling
    sibling = sibling.previousElementSibling
  }
}

const roots = new WeakMap()

/** Returns the controller of a [cmdk-root] element. */
export function getCommand(root) {
  return roots.get(root)
}

function createCommand(root) {
  const own = (el) => el.closest("[cmdk-root]") === root
  const ownAll = (sel) => [...root.querySelectorAll(sel)].filter(own)
  const input = ownAll("[cmdk-input]")[0] ?? null
  const list = ownAll("[cmdk-list]")[0] ?? null
  const sizer = list?.querySelector(":scope > [cmdk-list-sizer]") ?? null
  const opts = {
    loop: root.hasAttribute("data-cmdk-loop"),
    shouldFilter: root.getAttribute("data-cmdk-should-filter") !== "false",
    vimBindings: root.getAttribute("data-cmdk-vim-bindings") !== "false",
    disablePointerSelection: root.hasAttribute("data-cmdk-disable-pointer-selection"),
  }

  // --- the React tree, as virtual children of each host container ---
  // node: { el, kind: "item"|"group"|"separator"|"empty"|"static", mounted, ... }
  const containers = new Map() // container element -> [node]
  const items = [] // item nodes, in source order
  const groups = [] // group nodes
  const separators = []
  let empty = null

  const scan = (container) => {
    const nodes = []
    for (const child of [...container.children]) {
      if (child.tagName === "TEMPLATE" && child.hasAttribute("data-cmdk-empty-template")) {
        const el = child.content.firstElementChild?.cloneNode(true)
        if (el) {
          empty = { el, kind: "empty", mounted: false, anchor: child }
          nodes.push(empty)
        }
        continue
      }
      if (child.hasAttribute("cmdk-item")) {
        const node = {
          el: child,
          kind: "item",
          mounted: true,
          id: child.id,
          value: child.getAttribute(VALUE_ATTR) ?? "",
          keywords: parseKeywords(child.getAttribute("data-cmdk-keywords")),
          disabled: child.getAttribute("aria-disabled") === "true",
          forceMount: child.hasAttribute("data-cmdk-force-mount"),
          group: null,
        }
        items.push(node)
        nodes.push(node)
        continue
      }
      if (child.hasAttribute("cmdk-group")) {
        const node = { el: child, kind: "group", mounted: true, forceMount: child.hasAttribute("data-cmdk-force-mount"), items: new Set() }
        groups.push(node)
        nodes.push(node)
        const inner = child.querySelector(":scope > [cmdk-group-items]")
        if (inner) {
          const before = items.length
          scan(inner)
          for (const it of items.slice(before)) it.group = it.group ?? node
        }
        continue
      }
      if (child.hasAttribute("cmdk-separator")) {
        const node = { el: child, kind: "separator", mounted: true, alwaysRender: child.hasAttribute("data-cmdk-always-render") }
        separators.push(node)
        nodes.push(node)
        continue
      }
      if (child.hasAttribute("cmdk-empty")) {
        empty = { el: child, kind: "empty", mounted: true }
        nodes.push(empty)
        continue
      }
      nodes.push({ el: child, kind: "static", mounted: true })
      // Items wrapped in other elements: their container is that element.
      if (child.querySelector("[cmdk-item],[cmdk-group],[cmdk-separator]")) scan(child)
    }
    containers.set(container, nodes)
  }
  if (sizer) scan(sizer)

  // --- the store (cmdk's state) ---
  const selected = items.find((n) => n.el.getAttribute("aria-selected") === "true")
  const state = {
    search: "",
    value: selected ? selected.value : "",
    selectedItemId: undefined,
    filtered: { count: 0, items: new Map(), groups: new Set() },
  }
  const registered = items.filter((n) => !n.forceMount) // allItems, in mount order
  const itemById = new Map(items.map((n) => [n.id, n]))
  const groupItems = new Map() // allGroups: group node -> Set(item id)
  for (const g of groups) groupItems.set(g, new Set())
  for (const n of registered) if (n.group) groupItems.get(n.group).add(n.id)
  for (const n of items) state.filtered.items.set(n.id, score(n.value, n.keywords))
  state.filtered.count = registered.length

  // --- scheduling and rendering ---
  let renderPending = false
  let scheduled = new Map()
  let scheduleDirty = false
  const schedule = (id, fn) => {
    scheduled.set(id, fn)
    scheduleDirty = true
  }
  const store = { emit: () => (renderPending = true) }

  function commit() {
    for (let guard = 0; guard < 50 && (renderPending || scheduleDirty); guard++) {
      if (renderPending) {
        renderPending = false
        render()
      }
      if (scheduleDirty) {
        scheduleDirty = false
        const map = scheduled
        map.forEach((fn) => fn())
        scheduled = new Map()
      }
    }
  }

  function setState(key, value, opts) {
    if (Object.is(state[key], value)) return
    state[key] = value
    if (key === "search") {
      filterItems()
      sort()
      schedule(1, selectFirstItem)
      // The first match is selected: show it, wherever the list had been
      // scrolled to before the search changed.
      if (list) list.scrollTop = 0
      emit(root, "shadcn:search-change", { search: value })
    } else if (key === "value") {
      const active = document.activeElement
      if (active && (active.hasAttribute("cmdk-input") || active.hasAttribute("cmdk-root"))) {
        const el = input ?? list
        el?.focus()
      }
      schedule(7, () => {
        state.selectedItemId = getSelectedItem()?.id
        store.emit()
      })
      if (!opts) schedule(5, scrollSelectedIntoView)
      emit(root, "shadcn:value-change", { value: value ?? "" })
    }
    store.emit()
  }

  const isVisible = (n) => {
    switch (n.kind) {
      case "item":
        if (n.forceMount || !opts.shouldFilter) return true
        return state.search ? state.filtered.items.get(n.id) > 0 : true
      case "separator":
        return n.alwaysRender || !state.search
      case "empty":
        return state.filtered.count === 0
      default:
        return true
    }
  }

  function render() {
    // Mutations, container by container (React: deletions, then
    // placements before the next stable host sibling).
    for (const [container, nodes] of containers) {
      const next = nodes.map(isVisible)
      nodes.forEach((n, i) => {
        if (n.mounted && !next[i]) {
          n.el.remove()
          n.mounted = false
          n.placed = false
        }
      })
      nodes.forEach((n, i) => {
        if (!n.mounted && next[i]) n.placed = true
      })
      nodes.forEach((n, i) => {
        if (!n.placed) return
        let before = null
        for (let j = i + 1; j < nodes.length; j++) {
          const s = nodes[j]
          if (s.mounted && !s.placed && s.el.parentNode === container) {
            before = s.el
            break
          }
        }
        container.insertBefore(n.el, before)
      })
      nodes.forEach((n) => {
        if (n.placed) {
          n.placed = false
          n.mounted = true
        }
      })
    }
    // Updates.
    for (const n of items) {
      const sel = !!state.value && state.value === n.value
      n.el.setAttribute("aria-selected", String(sel))
      n.el.setAttribute("data-selected", String(sel))
    }
    for (const g of groups) {
      const visible = g.forceMount || !opts.shouldFilter ? true : state.search ? state.filtered.groups.has(g) : true
      if (visible) g.el.removeAttribute("hidden")
      else g.el.setAttribute("hidden", "")
    }
    for (const el of [input, list]) {
      if (!el) continue
      if (state.selectedItemId === undefined) el.removeAttribute("aria-activedescendant")
      else el.setAttribute("aria-activedescendant", state.selectedItemId)
    }
    if (input) {
      // React keeps a controlled input's value attribute in sync.
      if (input.value !== state.search) input.value = state.search
      input.setAttribute("value", state.search)
    }
  }

  // --- cmdk's helpers ---
  function score(value, keywords) {
    return value ? commandScore(value, state.search, keywords) : 0
  }

  function sort() {
    if (!state.search || !opts.shouldFilter) return
    const scores = state.filtered.items
    const listInner = sizer
    getValidItems()
      .sort((a, b) => (scores.get(b.getAttribute("id")) ?? 0) - (scores.get(a.getAttribute("id")) ?? 0))
      .forEach((item) => {
        const group = item.closest(GROUP_ITEMS_SELECTOR)
        if (group) group.appendChild(item.parentElement === group ? item : item.closest(`${GROUP_ITEMS_SELECTOR} > *`))
        else listInner.appendChild(item.parentElement === listInner ? item : item.closest(`${GROUP_ITEMS_SELECTOR} > *`))
      })
    // cmdk then reorders groups by their best score, looking them up with
    // `[data-value="${encodeURIComponent(groupId)}"]`: groups carry their
    // heading as data-value, never their React id, so nothing moves.
  }

  function selectFirstItem() {
    const item = getValidItems().find((el) => el.getAttribute("aria-disabled") !== "true")
    const value = item?.getAttribute(VALUE_ATTR)
    setState("value", value || undefined)
  }

  function filterItems() {
    if (!state.search || !opts.shouldFilter) {
      state.filtered.count = registered.length
      return
    }
    state.filtered.groups = new Set()
    let itemCount = 0
    for (const n of registered) {
      const rank = score(n.value, n.keywords)
      state.filtered.items.set(n.id, rank)
      if (rank > 0) itemCount++
    }
    for (const [g, ids] of groupItems) {
      for (const id of ids) {
        if (state.filtered.items.get(id) > 0) {
          state.filtered.groups.add(g)
          break
        }
      }
    }
    state.filtered.count = itemCount
  }

  function scrollSelectedIntoView() {
    const item = getSelectedItem()
    if (item) {
      if (item.parentElement?.firstChild === item) {
        item.closest(GROUP_SELECTOR)?.querySelector(GROUP_HEADING_SELECTOR)?.scrollIntoView({ block: "nearest" })
      }
      item.scrollIntoView({ block: "nearest" })
    }
  }

  const getSelectedItem = () => sizer?.querySelector(`${ITEM_SELECTOR}[aria-selected="true"]`) ?? undefined
  const getValidItems = () => Array.from(sizer?.querySelectorAll(VALID_ITEM_SELECTOR) || [])

  function updateSelectedToIndex(index) {
    const item = getValidItems()[index]
    if (item) setState("value", item.getAttribute(VALUE_ATTR))
  }

  function updateSelectedByItem(change) {
    const selected = getSelectedItem()
    const all = getValidItems()
    const index = all.findIndex((item) => item === selected)
    let newSelected = all[index + change]
    if (opts.loop) {
      newSelected = index + change < 0 ? all[all.length - 1] : index + change === all.length ? all[0] : all[index + change]
    }
    if (newSelected) setState("value", newSelected.getAttribute(VALUE_ATTR))
  }

  function updateSelectedByGroup(change) {
    const selected = getSelectedItem()
    let group = selected?.closest(GROUP_SELECTOR)
    let item
    while (group && !item) {
      group = change > 0 ? findNextSibling(group, GROUP_SELECTOR) : findPreviousSibling(group, GROUP_SELECTOR)
      item = group?.querySelector(VALID_ITEM_SELECTOR)
    }
    if (item) setState("value", item.getAttribute(VALUE_ATTR))
    else updateSelectedByItem(change)
  }

  const last = () => updateSelectedToIndex(getValidItems().length - 1)
  const next = (e) => {
    e.preventDefault()
    if (e.metaKey) last()
    else if (e.altKey) updateSelectedByGroup(1)
    else updateSelectedByItem(1)
  }
  const prev = (e) => {
    e.preventDefault()
    if (e.metaKey) updateSelectedToIndex(0)
    else if (e.altKey) updateSelectedByGroup(-1)
    else updateSelectedByItem(-1)
  }

  // --- events ---
  const onKeyDown = (e) => {
    const isComposing = e.isComposing || e.keyCode === 229
    if (e.defaultPrevented || isComposing) return
    switch (e.key) {
      case "n":
      case "j":
        if (opts.vimBindings && e.ctrlKey) next(e)
        break
      case "ArrowDown":
        next(e)
        break
      case "p":
      case "k":
        if (opts.vimBindings && e.ctrlKey) prev(e)
        break
      case "ArrowUp":
        prev(e)
        break
      case "Home":
        e.preventDefault()
        updateSelectedToIndex(0)
        break
      case "End":
        e.preventDefault()
        last()
        break
      case "Enter": {
        e.preventDefault()
        const item = getSelectedItem()
        if (item) item.dispatchEvent(new Event(SELECT_EVENT))
      }
    }
    commit()
  }

  const itemOf = (target) => {
    const el = target instanceof Element ? target.closest(ITEM_SELECTOR) : null
    return el && own(el) ? itemById.get(el.id) : undefined
  }
  // Item: onPointerMove={disabled || disablePointerSelection ? undefined : select}
  const onPointerMove = (e) => {
    const n = itemOf(e.target)
    if (!n || n.disabled || opts.disablePointerSelection) return
    setState("value", n.value, true)
    commit()
  }
  // Item: onClick={disabled ? undefined : onSelect}
  const onSelect = (n) => {
    setState("value", n.value, true)
    emit(n.el, "shadcn:select", { value: n.value })
  }
  const onClick = (e) => {
    const n = itemOf(e.target)
    if (!n || n.disabled) return
    onSelect(n)
    commit()
  }
  const selectListeners = items
    .filter((n) => !n.disabled)
    .map((n) => {
      const fn = () => onSelect(n)
      n.el.addEventListener(SELECT_EVENT, fn)
      return [n.el, fn]
    })
  const onInput = () => {
    setState("search", input.value)
    commit()
  }

  root.addEventListener("keydown", onKeyDown)
  root.addEventListener("pointermove", onPointerMove)
  root.addEventListener("click", onClick)
  input?.addEventListener("input", onInput)

  // List: --cmdk-list-height follows the sizer's height.
  let raf = 0
  let ro = null
  if (sizer && list && typeof ResizeObserver !== "undefined") {
    ro = new ResizeObserver(() => {
      raf = requestAnimationFrame(() => {
        const height = sizer.offsetHeight
        list.style.setProperty("--cmdk-list-height", height.toFixed(1) + "px")
      })
    })
    ro.observe(sizer)
  }

  // Mount: cmdk scrolls the selected item into view.
  scrollSelectedIntoView()

  const controller = {
    get value() {
      return state.value ?? ""
    },
    get search() {
      return state.search
    },
    setValue(v) {
      setState("value", v)
      commit()
    },
    setSearch(s) {
      if (input) input.value = s
      setState("search", s)
      commit()
    },
  }
  roots.set(root, controller)

  return () => {
    root.removeEventListener("keydown", onKeyDown)
    root.removeEventListener("pointermove", onPointerMove)
    root.removeEventListener("click", onClick)
    input?.removeEventListener("input", onInput)
    for (const [el, fn] of selectListeners) el.removeEventListener(SELECT_EVENT, fn)
    cancelAnimationFrame(raf)
    ro?.disconnect()
    roots.delete(root)
  }
}

function parseKeywords(s) {
  if (!s) return []
  try {
    const v = JSON.parse(s)
    return Array.isArray(v) ? v.map((k) => String(k).trim()) : []
  } catch {
    return []
  }
}

define("command", {
  selector: "[cmdk-root]",
  mount(root) {
    return createCommand(root)
  },
})
