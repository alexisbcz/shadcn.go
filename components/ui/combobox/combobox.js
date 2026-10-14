// combobox.js: behavior of the Base UI Combobox (@base-ui/react 1.6.0,
// combobox/root/AriaCombobox and its parts, with the floating-ui-react
// hooks it composes: useListNavigation (virtual focus), useDismiss,
// useClick, useTypeahead, FloatingFocusManager, FloatingPortal and
// useAnchorPositioning) for the markup rendered by components/ui/combobox.
//
// Markup contract (markers are removed on mount, so the live DOM matches
// React's):
//   - every part carries data-combobox="<root id>" and data-combobox-part:
//     "hidden" (the hidden input, also carrying data-combobox-options, the
//     JSON {id, multiple, autoHighlight, disabled, items: [{v, l}], name,
//     inputInside, defaults}), "input", "trigger", "clear" (a template of the
//     clear button) and "clear-live" (the rendered one), "chips",
//     "chip-template", "content" (a template of the positioner, with
//     data-combobox-content='{"side","align","sideOffset","alignOffset",
//     "anchor"}');
//   - inside the content template, items carry data-combobox-item='{"v","l"}'
//     and their check indicator data-combobox-indicator; lists, groups,
//     labels, separators and the empty state data-combobox-list,
//     -group, -label, -separator and -empty;
//   - Value renders its text between <!--combobox-value:<id>--> and
//     <!--/combobox-value-->.
//
// Events, dispatched on the hidden input: "shadcn:value-change" ({ value }),
// "shadcn:open-change" ({ open }); on the list after each render:
// "shadcn:combobox-render" ({ groups }).

import { define, emit } from "./core.js"
import { autoUpdate, computePosition, flip, hide as nativeHide, limitShift, offset, platform, shift, size } from "./floating-ui-1.7.6.js"

const VISUALLY_HIDDEN_INPUT = "clip-path: inset(50%); overflow: hidden; white-space: nowrap; border: 0px; padding: 0px; width: 1px; height: 1px; margin: -1px; position: absolute;"
const VISUALLY_HIDDEN = "clip-path: inset(50%); overflow: hidden; white-space: nowrap; border: 0px; padding: 0px; width: 1px; height: 1px; margin: -1px; position: fixed; top: 0px; left: 0px;"
const OWNER_VISUALLY_HIDDEN = "clip-path: inset(50%); position: fixed; top: 0px; left: 0px;"
const LIVE_REGION_MARKER = "⁠"
const TYPEAHEAD_RESET_MS = 750

let idCounter = 0
const reactId = () => `_cb${(idCounter++).toString(36)}_`

// --- filter (internals/filter.mjs) ----------------------------------------------------------------

const collator = new Intl.Collator(undefined, { usage: "search", sensitivity: "base", ignorePunctuation: true })
function contains(text, query) {
  if (!query) return true
  for (let i = 0; i <= text.length - query.length; i += 1) {
    if (collator.compare(text.slice(i, i + query.length), query) === 0) return true
  }
  return false
}

// --- markOthers (floating-ui-react/utils/markOthers.mjs) ---------------------------------------------

const counters = { "aria-hidden": new WeakMap() }
let markerCounterMap = new WeakMap()
let uncontrolled = new WeakSet()
let lockCount = 0
const buildKeepSet = (targets) => {
  const keep = new Set()
  for (const target of targets) {
    let node = target
    while (node && !keep.has(node)) {
      keep.add(node)
      node = node.parentNode
    }
  }
  return keep
}
const collectOutside = (root, keep, stop) => {
  const outside = []
  const walk = (parent) => {
    if (!parent || stop.has(parent)) return
    for (const node of parent.children) {
      if (node.nodeName.toLowerCase() === "script") continue
      if (keep.has(node)) walk(node)
      else outside.push(node)
    }
  }
  walk(root)
  return outside
}
function markOthers(avoid, { ariaHidden = false, mark = true } = {}) {
  const body = document.body
  const avoidElements = avoid.filter((t) => body.contains(t))
  const markerTargets = mark ? collectOutside(body, buildKeepSet(avoidElements), new Set(avoidElements)) : []
  const hidden = []
  const marked = []
  if (ariaHidden) {
    const map = counters["aria-hidden"]
    const live = [...body.querySelectorAll("[aria-live]")]
    const control = avoidElements.concat(live)
    for (const node of collectOutside(body, buildKeepSet(control), new Set(control))) {
      const attr = node.getAttribute("aria-hidden")
      const already = attr !== null && attr !== "false"
      const count = (map.get(node) || 0) + 1
      map.set(node, count)
      hidden.push(node)
      if (count === 1 && already) uncontrolled.add(node)
      if (!already) node.setAttribute("aria-hidden", "true")
    }
  }
  for (const node of markerTargets) {
    const count = (markerCounterMap.get(node) || 0) + 1
    markerCounterMap.set(node, count)
    marked.push(node)
    if (count === 1) node.setAttribute("data-base-ui-inert", "")
  }
  lockCount += 1
  return () => {
    const map = counters["aria-hidden"]
    for (const el of hidden) {
      const count = (map.get(el) || 0) - 1
      map.set(el, count)
      if (!count) {
        if (!uncontrolled.has(el)) el.removeAttribute("aria-hidden")
        uncontrolled.delete(el)
      }
    }
    for (const el of marked) {
      const count = (markerCounterMap.get(el) || 0) - 1
      markerCounterMap.set(el, count)
      if (!count) el.removeAttribute("data-base-ui-inert")
    }
    lockCount -= 1
    if (!lockCount) {
      counters["aria-hidden"] = new WeakMap()
      uncontrolled = new WeakSet()
      markerCounterMap = new WeakMap()
    }
  }
}

// --- tabbable helpers (floating-ui-react/utils/tabbable.mjs) ------------------------------------

const CANDIDATE_SELECTOR = 'a[href],button,input,select,textarea,summary,details,iframe,object,embed,[tabindex],[contenteditable]:not([contenteditable="false"]),audio[controls],video[controls]'
const isCandidate = (el) => el != null && el.matches(CANDIDATE_SELECTOR) && (el.nodeName.toLowerCase() !== "input" || el.type !== "hidden")
function isFocusableElement(el) {
  if (!isCandidate(el) || !el.isConnected || el.matches(":disabled")) return false
  for (let cur = el; cur; cur = cur.parentElement) {
    const isAncestor = cur !== el
    if (cur.hasAttribute("inert") || cur.hasAttribute("hidden")) return false
    const st = getComputedStyle(cur)
    if (isAncestor ? st.display === "none" : st.display === "none" || st.visibility === "hidden" || st.visibility === "collapse" || !(el.offsetWidth || el.offsetHeight || el.getClientRects().length)) return false
  }
  return true
}
function tabbables(container = document.body) {
  const list = []
  const walk = (c) => {
    for (const child of c.children) {
      if (isCandidate(child)) list.push(child)
      walk(child)
    }
  }
  walk(container)
  return list.filter(isFocusableElement).filter((el) => el.tabIndex >= 0)
}
function getTabbableIn(container, dir) {
  const list = tabbables(container)
  if (list.length === 0) return undefined
  const index = list.indexOf(document.activeElement)
  const next = index === -1 ? (dir === 1 ? 0 : list.length - 1) : index + dir
  return list[next]
}
const getNextTabbable = (ref) => getTabbableIn(document.body, 1) || ref
const getPreviousTabbable = (ref) => getTabbableIn(document.body, -1) || ref
const isOutsideEvent = (event, container) => !event.relatedTarget || !containsEl(container || event.currentTarget, event.relatedTarget)
function disableFocusInside(container) {
  for (const el of tabbables(container)) {
    el.dataset.tabindex = el.getAttribute("tabindex") || ""
    el.setAttribute("tabindex", "-1")
  }
}
function enableFocusInside(container) {
  for (const el of container.querySelectorAll("[data-tabindex]")) {
    const t = el.dataset.tabindex
    delete el.dataset.tabindex
    if (t) el.setAttribute("tabindex", t)
    else el.removeAttribute("tabindex")
  }
}
const isTypeable = (el) => el instanceof HTMLInputElement || el instanceof HTMLTextAreaElement
const containsEl = (parent, child) => !!parent && !!child && (parent === child || parent.contains(child))

// --- the arrow middleware fork (floating-ui-react/middleware/arrow.mjs) --------------------------

const getAlignment = (placement) => placement.split("-")[1]
const getSide = (placement) => placement.split("-")[0]
const getSideAxis = (placement) => (["top", "bottom"].includes(getSide(placement)) ? "y" : "x")
const getAlignmentAxis = (placement) => (getSideAxis(placement) === "x" ? "y" : "x")
const clamp = (lo, v, hi) => Math.max(lo, Math.min(v, hi))
const baseArrow = (options) => ({
  name: "arrow",
  options,
  async fn(state) {
    const { x, y, placement, rects, platform: pf, elements, middlewareData } = state
    const { element, padding = 0 } = (typeof options === "function" ? options(state) : options) || {}
    if (element == null) return {}
    const pad = typeof padding === "number" ? { top: padding, right: padding, bottom: padding, left: padding } : { top: 0, right: 0, bottom: 0, left: 0, ...padding }
    const coords = { x, y }
    const axis = getAlignmentAxis(placement)
    const length = axis === "y" ? "height" : "width"
    const arrowDimensions = await pf.getDimensions(element)
    const isYAxis = axis === "y"
    const minProp = isYAxis ? "top" : "left"
    const maxProp = isYAxis ? "bottom" : "right"
    const clientProp = isYAxis ? "clientHeight" : "clientWidth"
    const endDiff = rects.reference[length] + rects.reference[axis] - coords[axis] - rects.floating[length]
    const startDiff = coords[axis] - rects.reference[axis]
    const clientSize = elements.floating[clientProp] || rects.floating[length]
    const centerToReference = endDiff / 2 - startDiff / 2
    const largestPossiblePadding = clientSize / 2 - arrowDimensions[length] / 2 - 1
    const minPadding = Math.min(pad[minProp], largestPossiblePadding)
    const maxPadding = Math.min(pad[maxProp], largestPossiblePadding)
    const min = minPadding
    const max = clientSize - arrowDimensions[length] - maxPadding
    const center = clientSize / 2 - arrowDimensions[length] / 2 + centerToReference
    const off = clamp(min, center, max)
    const shouldAddOffset = !middlewareData.arrow && getAlignment(placement) != null && center !== off && rects.reference[length] / 2 - (center < min ? minPadding : maxPadding) - arrowDimensions[length] / 2 < 0
    const alignmentOffset = shouldAddOffset ? (center < min ? center - min : center - max) : 0
    return {
      [axis]: coords[axis] + alignmentOffset,
      data: { [axis]: off, centerOffset: center - off - alignmentOffset, ...(shouldAddOffset && { alignmentOffset }) },
      reset: shouldAddOffset,
    }
  },
})
const nativeHideFn = nativeHide().fn
const hideMiddleware = {
  name: "hide",
  async fn(state) {
    const { width, height, x, y } = state.rects.reference
    const anchorHidden = width === 0 && height === 0 && x === 0 && y === 0
    const result = await nativeHideFn(state)
    return { data: { referenceHidden: result.data?.referenceHidden || anchorHidden } }
  },
}
const roundByDPR = (el, value) => {
  const dpr = (el.ownerDocument.defaultView || window).devicePixelRatio || 1
  return Math.round(value * dpr) / dpr
}

// --- the controller --------------------------------------------------------------------------------

function createCombobox(hidden) {
  const opts = JSON.parse(hidden.getAttribute("data-combobox-options") || "{}")
  const rootId = opts.id
  const multiple = !!opts.multiple
  const single = !multiple
  const autoHighlightMode = opts.autoHighlight ? "input-change" : false
  const comboboxDisabled = !!opts.disabled
  const items = (opts.items || []).map((it, index) => ({ value: it.v, label: it.l, index }))
  const parts = [...document.querySelectorAll(`[data-combobox="${CSS.escape(rootId)}"]`)].map((el) => [el, el.getAttribute("data-combobox-part")])
  const part = (name) => parts.filter(([, n]) => n === name).map(([el]) => el)
  for (const [p] of parts) {
    p.removeAttribute("data-combobox")
    p.removeAttribute("data-combobox-part")
  }
  hidden.removeAttribute("data-combobox-options")

  const contentTemplate = part("content")[0] ?? null
  const contentOpts = contentTemplate ? JSON.parse(contentTemplate.getAttribute("data-combobox-content") || "{}") : {}
  const inputInsidePopup = !!opts.inputInside
  let input = inputInsidePopup ? null : (part("input")[0] ?? null)
  const triggers = part("trigger")
  const clearTemplate = part("clear")[0] ?? null
  let clearEl = part("clear-live")[0] ?? null
  const chipsEl = part("chips")[0] ?? null
  const chipTemplate = part("chip-template")[0] ?? null
  const anchorName = contentOpts.anchor
  const anchorEl = anchorName ? document.querySelector(`[data-combobox-anchor="${CSS.escape(anchorName)}"]`) : null
  anchorEl?.removeAttribute("data-combobox-anchor")
  const inputLocalDisabled = (el) => {
    if (!el?.hasAttribute("data-combobox-disabled")) return false
    el.removeAttribute("data-combobox-disabled")
    return true
  }
  let inputDisabledProp = inputLocalDisabled(input)
  const clearDisabledProp = clearTemplate ? !!clearTemplate.content.querySelector("[data-combobox-disabled]") : false
  clearTemplate?.content.querySelector("[data-combobox-disabled]")?.removeAttribute("data-combobox-disabled")
  if (clearEl) inputLocalDisabled(clearEl)
  // Value text regions.
  const valueRegions = []
  {
    const walker = document.createTreeWalker(document.body, NodeFilter.SHOW_COMMENT)
    while (walker.nextNode()) {
      const c = walker.currentNode
      if (c.nodeValue === `combobox-value:${rootId}`) valueRegions.push(c)
    }
  }
  const chips = () => (chipsEl ? [...chipsEl.children].filter((c) => c.getAttribute("data-slot") === "combobox-chip") : [])

  // --- state ---
  const byValue = (v) => items.find((it) => it.value === v) ?? null
  let selected = multiple ? [] : null
  if (Array.isArray(opts.defaults)) {
    if (multiple) selected = opts.defaults.map(byValue).filter(Boolean)
    else selected = byValue(opts.defaults[0]) ?? null
  }
  const hasSelectedValue = () => (multiple ? selected.length > 0 : selected != null)
  let inputValue = input ? input.value : single && selected ? selected.label : ""
  let open = false
  let mounted = false
  let transitionStatus
  let activeIndex = null
  let selectedIndex = null
  let queryChangedAfterOpen = false
  let closeQuery = null
  let openMethod = null
  let popupSide = null
  let mountSide = null
  let highlightedChipIndex
  let keyboardActive = true
  let hadInputClear = false
  let pendingQueryHighlight = null
  let domReference = inputInsidePopup ? triggers[0] ?? null : input
  let openEvent = null
  const FLOATING_REASONS = new Set(["trigger-press", "outside-press", "escape-key", "focus-out", "cancel-open"])
  let lastHighlight = null

  // list navigation refs
  let indexRef = -1
  let keyRef = null
  let isPointerModality = true
  let focusItemOnOpenRef
  let previousMounted = false
  let previousOpen = false
  // dismiss / focus manager
  let insideReactTree = false
  let isPointerDown = false

  const focusItemOnOpenProp = () => (queryChangedAfterOpen ? false : "auto")
  focusItemOnOpenRef = focusItemOnOpenProp()

  // popup DOM (while mounted)
  let portal = null
  let positioner = null
  let popup = null
  let list = null
  let emptyEl = null
  let emptyChildren = []
  let listModel = null // { flat: [itemEntry] } or { groups: [{ el, label, entries, rest }] }
  let startDismiss = null
  let endDismiss = null
  let guards = null
  let popupInput = null
  let unmarkOthers = null
  let stopAutoUpdate = null
  let isPositioned = false
  let floatingStyles = null
  let popupStyleTouched = false
  let positionerStyleExtra = {}
  let listId = null
  let finishController = null

  // --- derived ---
  const query = () => closeQuery ?? (inputValue === "" ? "" : String(inputValue).trim())
  const selectedLabel = () => (single && selected ? selected.label : "")
  const filteredItems = () => {
    const q = query()
    const bypass = single && !queryChangedAfterOpen && q !== "" && selectedLabel() !== "" && selectedLabel().length === q.length && contains(selectedLabel(), q)
    const fq = bypass ? "" : q
    if (fq === "") return items
    return items.filter((it) => {
      if (single && !queryChangedAfterOpen) {
        const s = selectedLabel()
        if (s && contains(s, fq) && s.length === fq.length) return true
      }
      return contains(it.label, fq)
    })
  }
  const isSelected = (it) => (multiple ? selected.includes(it) : selected === it)

  // --- DOM helpers ---
  const setAttr = (el, name, value) => {
    if (!el) return
    if (value === null || value === undefined || value === false) el.removeAttribute(name)
    else el.setAttribute(name, value === true ? "" : String(value))
  }
  const setValue = (el, v) => {
    if (!el) return
    if (el.value !== v) el.value = v
    if (el.getAttribute("value") !== v) el.setAttribute("value", v)
  }

  // --- render (React's re-render, as direct DOM writes) ---
  function render() {
    const filtered = filteredItems()
    const listEmpty = filtered.length === 0
    const side = mounted && positioner ? popupSide : null
    const disabled = comboboxDisabled
    // inputs
    for (const el of [input, popupInput]) {
      if (!el) continue
      const inside = el === popupInput
      const inDisabled = disabled || (inside ? false : inputDisabledProp)
      setValue(el, inputValue)
      setAttr(el, "aria-expanded", open ? "true" : "false")
      setAttr(el, "aria-controls", open ? listId : null)
      setAttr(el, "aria-activedescendant", open && activeIndex != null && !(inputInsidePopup && el !== domReference && false) ? `${rootId}-${activeIndex}` : null)
      setAttr(el, "data-popup-open", open)
      setAttr(el, "data-pressed", open)
      setAttr(el, "data-popup-side", side)
      setAttr(el, "data-list-empty", listEmpty)
      setAttr(el, "data-disabled", comboboxDisabled ? true : null)
      void inDisabled
    }
    // triggers
    for (const t of triggers) {
      setAttr(t, "aria-expanded", open ? "true" : "false")
      setAttr(t, "aria-controls", open ? (inputInsidePopup ? `${rootId}-popup` : listId) : null)
      setAttr(t, "data-popup-open", open)
      setAttr(t, "data-pressed", open)
      setAttr(t, "data-popup-side", side)
      setAttr(t, "data-list-empty", listEmpty)
      setAttr(t, "data-placeholder", multiple ? null : !hasSelectedValue())
    }
    // value regions
    for (const start of valueRegions) {
      let text = ""
      if (multiple) text = selected.map((it) => it.label).join(", ")
      else if (selected) text = selected.label
      else text = start.__placeholder ?? (start.__placeholder = readRegion(start))
      writeRegion(start, text)
    }
    renderChips()
    renderClear()
    // hidden input
    setValue(hidden, multiple ? "" : selected ? selected.value : "")
    if (multiple && opts.name) renderHiddenInputs()
    // popup
    if (mounted && positioner) renderPopup(filtered, listEmpty)
    // useListNavigation: sync the navigation index with activeIndex.
    if (open && positioner) {
      if (activeIndex !== lastSyncedActive) {
        if (activeIndex == null) {
          if (selectedIndex == null && previousMounted) indexRef = -1
        } else if (activeIndex < filtered.length) {
          indexRef = activeIndex
          focusItem()
          forceScrollIntoView = false
        }
        lastSyncedActive = activeIndex
      }
    } else lastSyncedActive = undefined
  }
  let lastSyncedActive
  // useListNavigation's focusItem (virtual): scroll the item into view when
  // navigating with the keyboard or opening on the selected item.
  let forceSyncFocus = false
  let forceScrollIntoView = false
  function focusItem() {
    const force = forceScrollIntoView
    const run = () => {
      const el = currentItemElements()[indexRef]
      if (el && (force || !isPointerModality)) el.scrollIntoView?.({ block: "nearest", inline: "nearest" })
    }
    if (forceSyncFocus) run()
    else requestAnimationFrame(run)
  }

  function readRegion(start) {
    let text = ""
    for (let n = start.nextSibling; n && !(n.nodeType === 8 && n.nodeValue === "/combobox-value"); n = n.nextSibling) text += n.textContent
    return text
  }
  function writeRegion(start, text) {
    const nodes = []
    for (let n = start.nextSibling; n && !(n.nodeType === 8 && n.nodeValue === "/combobox-value"); n = n.nextSibling) nodes.push(n)
    if (nodes.length === 1 && nodes[0].nodeType === 3 && nodes[0].nodeValue === text) return
    if (nodes.length === 0 && text === "") return
    nodes.forEach((n) => n.remove())
    if (text !== "") start.after(document.createTextNode(text))
  }

  function renderChips() {
    if (!chipsEl) return
    setAttr(chipsEl, "role", multiple && selected.length > 0 ? "toolbar" : null)
    if (!chipTemplate) return
    const current = chips()
    const want = multiple ? selected : []
    // Keyed by item: keep existing chips, add new ones before the input.
    const byItem = new Map(current.map((c) => [c.__item, c]))
    for (const c of current) if (!want.includes(c.__item)) c.remove()
    let anchor = chipTemplate
    for (const it of want) {
      let c = byItem.get(it)
      if (!c || !c.isConnected) {
        c = chipTemplate.content.firstElementChild.cloneNode(true)
        c.__item = it
        c.prepend(document.createTextNode(it.label))
        wireChip(c)
      }
      if (anchor.nextSibling !== c) anchor.after(c)
      anchor = c
    }
  }

  function renderClear() {
    if (!clearTemplate) return
    const visible = hasSelectedValue()
    if (visible) {
      if (!clearEl || clearEl.__ending) {
        if (clearEl) clearEl.remove()
        clearEl = clearTemplate.content.firstElementChild.cloneNode(true)
        clearTemplate.before(clearEl)
        wireClear(clearEl)
        clearEl.setAttribute("data-starting-style", "")
        const el = clearEl
        requestAnimationFrame(() => el.removeAttribute("data-starting-style"))
      }
      setAttr(clearEl, "data-visible", true)
      setAttr(clearEl, "data-popup-open", open)
    } else if (clearEl && !clearEl.__ending) {
      const el = clearEl
      el.__ending = true
      setAttr(el, "data-visible", null)
      el.setAttribute("data-ending-style", "")
      afterAnimations(el, () => {
        el.remove()
        if (clearEl === el) clearEl = null
      })
    } else if (clearEl) {
      setAttr(clearEl, "data-popup-open", open)
    }
  }

  function renderHiddenInputs() {
    let n = hidden.nextElementSibling
    const existing = []
    while (n && n.tagName === "INPUT" && n.type === "hidden" && n.name === opts.name) {
      existing.push(n)
      n = n.nextElementSibling
    }
    existing.forEach((e) => e.remove())
    let at = hidden
    for (const it of selected) {
      const el = document.createElement("input")
      el.type = "hidden"
      el.name = opts.name
      el.value = it.value
      if (comboboxDisabled) el.disabled = true
      at.after(el)
      at = el
    }
  }

  function renderPopup(filtered, listEmpty) {
    const side = getSide(mountSide ?? sideParam())
    const physicalSide = popupSide ?? side
    const align = renderedAlign ?? (contentOpts.align === "center" ? "center" : contentOpts.align || "start")
    for (const el of [positioner, popup]) {
      setAttr(el, "data-open", open)
      setAttr(el, "data-closed", !open)
      setAttr(el, "data-side", physicalSide)
      setAttr(el, "data-align", align)
      setAttr(el, "data-anchor-hidden", anchorHidden)
      setAttr(el, "data-empty", listEmpty)
    }
    setAttr(popup, "data-starting-style", transitionStatus === "starting")
    setAttr(popup, "data-ending-style", transitionStatus === "ending")
    // positioner style
    const style = { ...(floatingStyles ?? { position: "fixed", left: "0px", top: "0px", opacity: "0" }) }
    if (!open) style["pointer-events"] = "none"
    if (transitionStatus === "starting") style.transition = "none"
    for (const prop of ["position", "left", "top", "transform", "opacity", "pointer-events", "transition", "will-change"]) {
      const v = style[prop]
      if (v === undefined) positioner.style.removeProperty(prop)
      else positioner.style.setProperty(prop, v)
    }
    // popup style: React leaves style="" after removing the mount transition.
    if (transitionStatus === "starting") {
      popup.style.setProperty("transition", "none")
      popupStyleTouched = true
    } else if (popupStyleTouched) popup.style.removeProperty("transition")
    // empty
    if (emptyEl) {
      const want = listEmpty ? emptyChildren : []
      const have = [...emptyEl.childNodes]
      if (have.length !== want.length || have.some((n, i) => n !== want[i])) emptyEl.replaceChildren(...want)
    }
    // list
    setAttr(list, "data-empty", listEmpty)
    setAttr(list, "aria-activedescendant", inputInsidePopup && open && activeIndex != null ? `${rootId}-${activeIndex}` : null)
    renderList(filtered)
  }

  let lastFilteredKey = null
  function renderList(filtered) {
    const set = new Set(filtered)
    let index = 0
    const decorate = (entry) => {
      const it = entry.item
      const el = entry.el
      const sel = isSelected(it)
      const i = index++
      setAttr(el, "id", `${rootId}-${i}`)
      setAttr(el, "aria-selected", sel ? "true" : "false")
      setAttr(el, "data-selected", sel)
      setAttr(el, "data-highlighted", activeIndex === i)
      setIndicator(entry, sel)
      return el
    }
    const key = filtered.map((it) => it.index).join(",")
    // Keyed reconciliation: leave nodes in place when they already match.
    const reconcile = (parent, want) => {
      const have = [...parent.childNodes]
      if (have.length === want.length && have.every((n, i) => n === want[i])) return
      parent.replaceChildren(...want)
    }
    if (listModel.flat) {
      const els = listModel.flat.filter((e) => set.has(e.item)).map(decorate)
      if (key !== lastFilteredKey) {
        reconcile(list, [...listModel.leading, ...els, ...listModel.trailing])
      }
    } else {
      const groupEls = []
      for (const g of listModel.groups) {
        const entries = g.entries.filter((e) => set.has(e.item))
        if (entries.length === 0) continue
        const els = entries.map(decorate)
        if (key !== lastFilteredKey) reconcile(g.el, [...g.leading, ...els, ...g.trailing])
        groupEls.push(g.el)
      }
      if (key !== lastFilteredKey) {
        reconcile(list, [...listModel.leading, ...groupEls, ...listModel.trailing])
        emit(list, "shadcn:combobox-render", { groups: groupEls })
      }
    }
    lastFilteredKey = key
  }

  function setIndicator(entry, sel) {
    const ind = entry.indicator
    if (!ind) return
    if (sel) {
      if (!ind.isConnected || ind.__ending) {
        ind.__ending = false
        ind.removeAttribute("data-ending-style")
        entry.el.append(ind)
        if (entry.__shown) {
          ind.setAttribute("data-starting-style", "")
          requestAnimationFrame(() => ind.removeAttribute("data-starting-style"))
        }
      }
      setAttr(ind, "data-selected", true)
      entry.__shown = true
    } else if (ind.isConnected && !ind.__ending) {
      ind.__ending = true
      setAttr(ind, "data-selected", null)
      ind.setAttribute("data-ending-style", "")
      afterAnimations(ind, () => {
        if (ind.__ending) {
          ind.remove()
          ind.__ending = false
          ind.removeAttribute("data-ending-style")
        }
      })
    }
  }

  function afterAnimations(el, fn) {
    requestAnimationFrame(() => {
      const exec = () =>
        Promise.all(el.getAnimations().map((a) => a.finished)).then(fn, () => {
          if (el.getAnimations().some((a) => a.pending || a.playState !== "finished")) exec()
        })
      exec()
    })
  }

  // --- indices (setIndices / onItemHighlighted) ---
  function setIndices({ activeIndex: a, selectedIndex: s } = {}) {
    if (a !== undefined) activeIndex = a
    if (s !== undefined) selectedIndex = s
    render()
  }

  // --- value / input setters ---
  function setInputValue(next, reason, event) {
    hadInputClear = reason === "input-clear"
    if (reason === "input-change" && event) {
      const inputType = event.inputType
      const typed = event.type === "compositionend" || (inputType != null && inputType !== "" && inputType !== "insertReplacementText")
      if (typed) {
        const hasQuery = next.trim() !== ""
        if (hasQuery) queryChangedAfterOpen = true
        pendingQueryHighlight = { hasQuery }
        if (hasQuery && autoHighlightMode && activeIndex == null) activeIndex = 0
      }
    }
    if (inputValue === next) {
      scheduleLayout()
      return
    }
    inputValue = next
    scheduleQueryEffect()
    scheduleLayout()
    render()
  }
  // useValueChanged(query): once the event's updates are committed, a new
  // non-empty query while open marks the query as changed after opening.
  let lastEffectQuery
  let queryEffectQueued = false
  function scheduleQueryEffect() {
    if (queryEffectQueued) return
    queryEffectQueued = true
    queueMicrotask(() => {
      queryEffectQueued = false
      const q = query()
      if (q === lastEffectQuery) return
      lastEffectQuery = q
      if (!open || q === "" || q === String(initialInputValue)) return
      if (!queryChangedAfterOpen) {
        queryChangedAfterOpen = true
        render()
      }
    })
  }
  const initialInputValue = inputValue
  lastEffectQuery = query()

  // The root's layout effect on [activeIndex, flatFilteredItems, open].
  // The root's layout effect runs once the event's updates are committed.
  let layoutQueued = false
  function scheduleLayout() {
    if (layoutQueued) return
    layoutQueued = true
    queueMicrotask(() => {
      layoutQueued = false
      const before = activeIndex
      flushQueryHighlight()
      if (activeIndex !== before) render()
    })
  }
  function flushQueryHighlight() {
    const pending = pendingQueryHighlight
    if (pending) {
      if (pending.hasQuery && autoHighlightMode) activeIndex = 0
      pendingQueryHighlight = null
    }
    if (!open) return
    const candidates = filteredItems()
    if (activeIndex == null) return
    if (activeIndex >= candidates.length) activeIndex = null
  }

  function setSelectedValue(next, reason) {
    const prev = selected
    // The query of the render the event hit (before the input is filled).
    const renderQuery = query()
    selected = next
    if (single && !inputInsidePopup) {
      // shouldFillInput, then useValueChanged(selectedValue).
      const label = next ? next.label : ""
      setInputValue(label, reason)
    }
    if (single && next != null && reason !== "input-change" && queryChangedAfterOpen) closeQuery = renderQuery
    const prevSelectedIndex = selectedIndex
    syncSelectedIndex()
    if (prev !== next) forceMount()
    if (selectedIndex !== prevSelectedIndex) resetTypeaheadOnChange()
    if (prev !== next) emit(hidden, "shadcn:value-change", { value: multiple ? next.map((it) => it.value) : next ? next.value : null })
    render()
  }

  function handleSelection(event, it) {
    let itemValue = it
    if (itemValue === undefined) {
      if (activeIndex === null) return
      itemValue = filteredItems()[activeIndex]
    }
    if (multiple) {
      const was = selected.includes(itemValue)
      const next = was ? selected.filter((x) => x !== itemValue) : [...selected, itemValue]
      setSelectedValue(next, "item-press")
      const wasFiltering = activeInput() ? activeInput().value.trim() !== "" : false
      if (!wasFiltering) return
      if (inputInsidePopup) setInputValue("", "input-clear")
      else setOpen(false, "item-press", event)
    } else {
      setSelectedValue(itemValue, "item-press")
      setOpen(false, "item-press", event)
    }
  }

  const activeInput = () => (inputInsidePopup ? popupInput : input)

  // --- open / close ---
  function setOpen(next, reason, event) {
    if (open === next) return
    if (!emit(hidden, "shadcn:open-change", { open: next, reason })) return
    if (next && multiple && inputInsidePopup && closeQuery !== null) {
      queryChangedAfterOpen = false
      closeQuery = null
      if (inputValue !== "") setInputValue("", "input-clear")
    }
    if (!next && queryChangedAfterOpen) {
      if (single) {
        closeQuery = query()
        if (query() === "") queryChangedAfterOpen = false
      } else {
        closeQuery = query()
        if (inputInsidePopup) activeIndex = null
        if (!inputInsidePopup) setInputValue("", "input-clear")
      }
    }
    open = next
    // The floating root's openEvent (FloatingRootStore.syncOpenEvent) only
    // follows changes made through floating-ui hooks (useClick, useDismiss);
    // the combobox's own changes (typing, items, list navigation) leave it
    // as it was, stale included.
    if (FLOATING_REASONS.has(reason)) openEvent = next ? (event ?? null) : null
    // useTypeahead's layout effect also runs on open changes.
    resetTypeaheadOnChange()
    onOpenChangeLocal(next, reason, event)
    if (next) {
      if (!mounted) {
        mounted = true
        transitionStatus = "starting"
        mountPopup()
      } else {
        transitionStatus = "starting"
        // Interrupted close: the popup is reused.
        finishController?.abort()
        if (startDismiss) activeInput()?.before(startDismiss)
        if (guards) contentTemplate.before(guards.outBefore, guards.owns, guards.outAfter)
      }
      const frame = requestAnimationFrame(() => {
        if (open && transitionStatus === "starting") {
          transitionStatus = undefined
          render()
        }
      })
      void frame
    } else {
      transitionStatus = "ending"
      openMethod = null
      // ComboboxInput renders its dismiss button only while open, and
      // FloatingPortal its outside guards.
      startDismiss?.remove()
      if (guards) {
        guards.outBefore.remove()
        guards.owns.remove()
        guards.outAfter.remove()
      }
    }
    scheduleLayout()
    render()
    afterOpenChange()
  }

  function sideParam() {
    return contentOpts.side || "bottom"
  }

  // --- portal / popup mount ---
  function mountPopup() {
    const nodes = [...contentTemplate.content.cloneNode(true).children]
    positioner = nodes[0]
    popup = positioner.querySelector('[data-slot="combobox-content"]')
    list = positioner.querySelector("[data-combobox-list]")
    emptyEl = positioner.querySelector("[data-combobox-empty]")
    list?.removeAttribute("data-combobox-list")
    emptyEl?.removeAttribute("data-combobox-empty")
    emptyChildren = emptyEl ? [...emptyEl.childNodes] : []
    if (emptyEl) emptyEl.replaceChildren()
    for (const s of positioner.querySelectorAll("[data-combobox-separator]")) {
      s.removeAttribute("data-combobox-separator")
      s.setAttribute("data-orientation", "horizontal")
      s.setAttribute("role", "separator")
      s.setAttribute("aria-orientation", "horizontal")
    }
    for (const l of positioner.querySelectorAll("[data-combobox-label]")) {
      l.removeAttribute("data-combobox-label")
      l.id = `base-ui-${reactId()}`
      const g = l.closest("[data-combobox-group]")
      if (g) g.setAttribute("aria-labelledby", l.id)
    }
    // Items.
    const entryOf = (el) => {
      const data = JSON.parse(el.getAttribute("data-combobox-item"))
      el.removeAttribute("data-combobox-item")
      const it = items.find((x) => x.value === data.v && x.label === data.l && !x.__claimed) ?? items.find((x) => x.value === data.v)
      if (it) it.__claimed = true
      const indicator = el.querySelector("[data-combobox-indicator]")
      indicator?.removeAttribute("data-combobox-indicator")
      indicator?.remove()
      wireItem(el, it)
      return { el, item: it, indicator }
    }
    const splitChildren = (container, isEntry) => {
      const leading = []
      const trailing = []
      let seen = false
      for (const c of [...container.childNodes]) {
        if (isEntry(c)) seen = true
        else (seen ? trailing : leading).push(c)
      }
      return { leading, trailing }
    }
    const groups = list ? [...list.querySelectorAll("[data-combobox-group]")] : []
    if (groups.length) {
      const { leading, trailing } = splitChildren(list, (c) => c.nodeType === 1 && c.hasAttribute("data-combobox-group"))
      listModel = {
        leading,
        trailing,
        groups: groups.map((g) => {
          g.setAttribute("role", "group")
          const itemEls = [...g.querySelectorAll("[data-combobox-item]")]
          const { leading: gl, trailing: gt } = splitChildren(g, (c) => c.nodeType === 1 && c.hasAttribute("data-combobox-item"))
          g.removeAttribute("data-combobox-group")
          return { el: g, leading: gl, trailing: gt, entries: itemEls.map(entryOf) }
        }),
      }
    } else if (list) {
      const itemEls = [...list.querySelectorAll("[data-combobox-item]")]
      const { leading, trailing } = splitChildren(list, (c) => c.nodeType === 1 && c.hasAttribute("data-combobox-item"))
      listModel = { leading, trailing, flat: itemEls.map(entryOf) }
    } else listModel = { flat: [], leading: [], trailing: [] }
    for (const it of items) delete it.__claimed
    lastFilteredKey = null

    // List props.
    listId = reactId()
    if (list) {
      list.setAttribute("tabindex", "-1")
      list.setAttribute("data-base-ui-focusable", "")
      list.setAttribute("aria-orientation", "vertical")
      list.setAttribute("role", "listbox")
      list.id = listId
      if (multiple) list.setAttribute("aria-multiselectable", "true")
      wireList(list)
    }
    // Popup props.
    popup.setAttribute("tabindex", "-1")
    popup.setAttribute("role", inputInsidePopup ? "dialog" : "presentation")
    if (inputInsidePopup) popup.id = `${rootId}-popup`
    popupStyleTouched = false
    // The input inside the popup.
    if (inputInsidePopup) {
      popupInput = popup.querySelector("[data-combobox-part=input]")
      if (popupInput) {
        popupInput.removeAttribute("data-combobox")
        popupInput.removeAttribute("data-combobox-part")
        inputLocalDisabled(popupInput)
        wireInput(popupInput)
        // setInputElement: an input inside the popup starts empty.
        inputValue = ""
      }
    }
    // Focus manager parts.
    const modal = !inputInsidePopup
    if (modal) {
      startDismiss = dismissButton()
      endDismiss = dismissButton()
      const host = activeInput()
      host?.before(startDismiss)
    } else {
      startDismiss = endDismiss = null
    }
    portal = document.createElement("div")
    portal.id = reactId()
    portal.setAttribute("data-base-ui-portal", "")
    if (!modal) {
      // FloatingFocusManager's inside guards around the popup.
      const before = focusGuard("inside")
      const after = focusGuard("inside")
      positioner.prepend(before)
      positioner.append(after)
      wireInsideGuards(before, after)
      // FloatingPortal's outside guards and aria-owns span, where the portal
      // sits in the React tree.
      const outBefore = focusGuard("outside")
      const owns = document.createElement("span")
      owns.setAttribute("aria-owns", portal.id)
      owns.setAttribute("style", OWNER_VISUALLY_HIDDEN)
      const outAfter = focusGuard("outside")
      contentTemplate.before(outBefore, owns, outAfter)
      guards = { before, after, outBefore, owns, outAfter }
      wireOutsideGuards(outBefore, outAfter)
      focusInsideDisabled = false
      wirePortalFocus(portal)
    }
    if (endDismiss) positioner.append(endDismiss)
    portal.append(positioner)
    document.body.append(portal)
    isPositioned = false
    floatingStyles = null
    mountSide = null
    renderedAlign = null
    anchorHidden = false
    popupSide = getSide(sideParam())
    startPositioning()
    // Empty's initial live-region marker.
    if (emptyEl && emptyChildren.length) {
      const pending = filteredItems().length === 0
      if (pending) {
        queueMicrotask(() => {
          const walker = document.createTreeWalker(emptyEl, NodeFilter.SHOW_TEXT)
          let last = null
          while (walker.nextNode()) if (walker.currentNode.nodeValue !== "") last = walker.currentNode
          if (!last) return
          const original = last.nodeValue
          const marked = original + LIVE_REGION_MARKER
          last.nodeValue = marked
          setTimeout(() => {
            if (last.nodeValue === marked) last.nodeValue = original
          }, 200)
        })
      }
    }
  }

  function unmountPopup() {
    stopAutoUpdate?.()
    stopAutoUpdate = null
    portal?.remove()
    startDismiss?.remove()
    if (guards) {
      guards.outBefore.remove()
      guards.owns.remove()
      guards.outAfter.remove()
    }
    portal = positioner = popup = list = emptyEl = null
    popupInput = null
    startDismiss = endDismiss = null
    guards = null
    listModel = null
  }

  function dismissButton() {
    const span = document.createElement("span")
    span.setAttribute("role", "button")
    span.setAttribute("aria-label", "Dismiss")
    span.setAttribute("style", VISUALLY_HIDDEN_INPUT)
    span.addEventListener("click", (e) => setOpen(false, "close-press", e))
    return span
  }

  function focusGuard(type) {
    const span = document.createElement("span")
    span.setAttribute("data-type", type)
    span.setAttribute("aria-hidden", "true")
    span.setAttribute("tabindex", "0")
    span.setAttribute("data-base-ui-focus-guard", "")
    span.setAttribute("style", VISUALLY_HIDDEN)
    return span
  }

  // FloatingFocusManager's inside guards (non-modal, in a portal).
  function wireInsideGuards(before, after) {
    before.addEventListener("focus", (event) => {
      preventReturnFocus = false
      if (isOutsideEvent(event, portal)) getNextTabbable(domReference)?.focus()
      else guards?.outBefore?.focus()
    })
    after.addEventListener("focus", (event) => {
      preventReturnFocus = true
      if (isOutsideEvent(event, portal)) getPreviousTabbable(domReference)?.focus()
      else guards?.outAfter?.focus()
    })
  }
  // FloatingPortal's outside guards.
  function wireOutsideGuards(outBefore, outAfter) {
    outBefore.addEventListener("focus", (event) => {
      if (isOutsideEvent(event, portal)) guards?.before?.focus()
      else getPreviousTabbable(domReference)?.focus()
    })
    outAfter.addEventListener("focus", (event) => {
      if (isOutsideEvent(event, portal)) guards?.after?.focus()
      else {
        getNextTabbable(domReference)?.focus()
        setOpen(false, "focus-out", event)
      }
    })
  }
  // FloatingPortal (non-modal): tabbables inside the portal are tabbable
  // only while focus is inside it.
  let focusInsideDisabled = false
  function wirePortalFocus(node) {
    const onFocus = (event) => {
      if (event.relatedTarget && isOutsideEvent(event, node)) {
        if (event.type === "focusin") {
          if (focusInsideDisabled) {
            enableFocusInside(node)
            focusInsideDisabled = false
          }
        } else {
          disableFocusInside(node)
          focusInsideDisabled = true
        }
      }
    }
    node.addEventListener("focusin", onFocus, true)
    node.addEventListener("focusout", onFocus, true)
  }

  // --- positioning (useAnchorPositioning) ---
  let renderedAlign = null
  let anchorHidden = false
  function referenceElement() {
    if (anchorEl) return anchorEl
    if (inputInsidePopup) return triggers[0] ?? null
    return input
  }
  function startPositioning() {
    const ref = referenceElement()
    if (!ref || !positioner) return
    const virtual = { getBoundingClientRect: () => ref.getBoundingClientRect(), getClientRects: () => ref.getClientRects(), contextElement: ref }
    const update = () => positionNow(virtual)
    stopAutoUpdate = autoUpdate(virtual, positioner, update, { elementResize: typeof ResizeObserver !== "undefined", layoutShift: typeof IntersectionObserver !== "undefined" })
  }
  function positionNow(reference) {
    const fl = positioner
    if (!fl) return
    const sp = sideParam()
    const side = mountSide || sp
    const align = contentOpts.align || "start"
    const placement = align === "center" ? side : `${side}-${align}`
    const bias = 1
    const pad = { top: 5 + (sp === "bottom" ? bias : 0), right: 5 + (sp === "left" ? bias : 0), bottom: 5 + (sp === "top" ? bias : 0), left: 5 + (sp === "right" ? bias : 0) }
    const common = { boundary: "clippingAncestors", padding: pad }
    const sideOffset = contentOpts.sideOffset ?? 6
    const alignOffset = contentOpts.alignOffset ?? 0
    const flipMw = flip({ ...common, padding: { top: pad.top + bias, right: pad.right + bias, bottom: pad.bottom + bias, left: pad.left + bias }, mainAxis: true, crossAxis: "alignment", fallbackAxisSideDirection: "none" })
    const shiftMw = shift(() => ({ ...common, mainAxis: true, crossAxis: false, limiter: limitShift(() => ({})) }))
    const fakeArrow = document.createElement("div")
    const middleware = [
      offset(() => ({ mainAxis: sideOffset, crossAxis: alignOffset, alignmentAxis: alignOffset })),
      ...(align === "center" ? [shiftMw, flipMw] : [flipMw, shiftMw]),
      size({
        ...common,
        apply({ elements, availableWidth, availableHeight, rects }) {
          if (!mounted) return
          const st = elements.floating.style
          st.setProperty("--available-width", `${availableWidth}px`)
          st.setProperty("--available-height", `${availableHeight}px`)
          const dpr = window.devicePixelRatio || 1
          const { x, y, width, height } = rects.reference
          st.setProperty("--anchor-width", `${(Math.round((x + width) * dpr) - Math.round(x * dpr)) / dpr}px`)
          st.setProperty("--anchor-height", `${(Math.round((y + height) * dpr) - Math.round(y * dpr)) / dpr}px`)
        },
      }),
      baseArrow(() => ({ element: fakeArrow, padding: 5, offsetParent: "floating" })),
      {
        name: "transformOrigin",
        fn(state) {
          const { elements, middlewareData, placement: rp, rects, y } = state
          const cs = getSide(rp)
          const axis = getSideAxis(rp)
          const arrowX = middlewareData.arrow?.x || 0
          const arrowY = middlewareData.arrow?.y || 0
          const transformX = arrowX + 0 / 2
          const transformY = arrowY + 0 / 2
          const shiftY = Math.abs(middlewareData.shift?.y || 0)
          const half = rects.reference.height / 2
          const overlapping = shiftY > sideOffset
          const adjacent = { top: `${transformX}px calc(100% + ${sideOffset}px)`, bottom: `${transformX}px ${-sideOffset}px`, left: `calc(100% + ${sideOffset}px) ${transformY}px`, right: `${-sideOffset}px ${transformY}px` }[cs]
          const overlap = `${transformX}px ${rects.reference.y + half - y}px`
          void axis
          void overlap
          void overlapping
          elements.floating.style.setProperty("--transform-origin", adjacent)
          return {}
        },
      },
      hideMiddleware,
    ]
    computePosition(reference, fl, { placement, strategy: "absolute", middleware }).then((data) => {
      if (positioner !== fl) return
      const x = roundByDPR(fl, data.x)
      const yy = roundByDPR(fl, data.y)
      floatingStyles = { position: "absolute", left: "0px", top: "0px", transform: `translate(${x}px, ${yy}px)` }
      const wasPositioned = isPositioned
      isPositioned = true
      popupSide = getSide(data.placement)
      renderedAlign = getAlignment(data.placement) || "center"
      anchorHidden = !!data.middlewareData.hide?.referenceHidden
      // lazyFlip: lock the rendered side once positioned.
      if (mounted) mountSide = popupSide
      void wasPositioned
      render()
    })
  }

  // --- focus manager / list navigation effects after open changes ---
  let preventReturnFocus = false
  let closeType = ""
  let elementFocusedBeforeOpen = null
  function onOpenChangeLocal(next, reason, event) {
    if (!next) closeType = event ? (event instanceof KeyboardEvent ? "keyboard" : event.pointerType || (event instanceof MouseEvent ? (event.detail === 0 ? "keyboard" : "mouse") : "keyboard")) : ""
    if (reason === "outside-press") preventReturnFocus = false
  }

  function afterOpenChange() {
    if (open) {
      if (!previousOpen) elementFocusedBeforeOpen = document.activeElement
      // markOthers (FloatingFocusManager effect, on open).
      unmarkOthers?.()
      const modal = !inputInsidePopup
      const untrapped = !inputInsidePopup && isTypeable(domReference)
      const insideEls = [positioner, guards?.before, guards?.after, guards?.outBefore, guards?.outAfter, startDismiss, endDismiss, untrapped ? domReference : null].filter(Boolean)
      const unA = markOthers(insideEls, { ariaHidden: modal || untrapped, mark: false })
      const unB = markOthers([positioner].filter(Boolean))
      unmarkOthers = () => {
        unB()
        unA()
      }
      // List navigation: sync the selected index, initial highlight.
      indexRef = selectedIndex ?? -1
      if (focusItemOnOpenRef && selectedIndex != null) {
        forceScrollIntoView = true
        onNavigate(null)
      }
      if (activeIndex == null && selectedIndex == null && (!previousOpen || !previousMounted) && focusItemOnOpenRef && (keyRef != null || (focusItemOnOpenRef === true && keyRef == null))) {
        const filtered = filteredItems()
        if (filtered.length) {
          indexRef = keyRef == null || isMainOrientationToEndKey(keyRef) ? 0 : filtered.length - 1
          keyRef = null
          onNavigate(null)
        }
      }
      // Initial focus (input inside the popup).
      if (inputInsidePopup && popupInput) {
        const target = popupInput
        requestAnimationFrame(() => {
          if (open && !containsEl(positioner, document.activeElement)) target.focus()
          else if (open && document.activeElement !== target && !containsEl(positioner, document.activeElement)) target.focus()
        })
      }
    } else {
      unmarkOthers?.()
      unmarkOthers = null
      forceSyncFocus = false
      keyRef = null
      focusItemOnOpenRef = focusItemOnOpenProp()
      // Retain the highlight while the popup animates out.
      const pop = popup
      finishController?.abort()
      const controller = (finishController = new AbortController())
      if (pop) {
        requestAnimationFrame(() => {
          const exec = () =>
            Promise.all(pop.getAnimations().map((a) => a.finished)).then(
              () => !controller.signal.aborted && handleUnmount(),
              () => {
                if (controller.signal.aborted) return
                if (pop.getAnimations().some((a) => a.pending || a.playState !== "finished")) exec()
              }
            )
          exec()
        })
      } else handleUnmount()
    }
    previousOpen = open
    previousMounted = mounted
  }

  function handleUnmount() {
    if (open) return
    const wasFocusInside = containsEl(positioner, document.activeElement)
    const activeEl = document.activeElement
    mounted = false
    transitionStatus = undefined
    queryChangedAfterOpen = false
    closeQuery = null
    activeIndex = null
    focusItemOnOpenRef = focusItemOnOpenProp()
    if (multiple && activeInput() && activeInput().value !== "" && !hadInputClear) setInputValue("", "input-clear")
    if (single) {
      if (inputInsidePopup) {
        if (inputValue !== "") setInputValue("", "input-clear")
      } else {
        const s = selected ? selected.label : ""
        if (input && input.value !== s) setInputValue(s, s === "" ? "input-clear" : "none")
      }
    }
    unmountPopup()
    previousMounted = false
    syncSelectedIndex()
    render()
    // Return focus (input inside the popup: back to the trigger).
    if (inputInsidePopup) {
      const returnEl = domReference?.isConnected ? domReference : elementFocusedBeforeOpen
      queueMicrotask(() => {
        const active = document.activeElement
        if (!preventReturnFocus && returnEl && (returnEl !== activeEl && activeEl !== document.body ? wasFocusInside : true)) returnEl.focus({ preventScroll: true })
        else if (!preventReturnFocus && returnEl && (active === document.body || !active)) returnEl.focus({ preventScroll: true })
        preventReturnFocus = false
      })
    }
  }

  function syncSelectedIndex() {
    if (open) return
    if (multiple) {
      const last = selected[selected.length - 1]
      selectedIndex = last ? items.indexOf(last) : null
    } else {
      selectedIndex = selected ? items.indexOf(selected) : null
    }
    if (selectedIndex === -1) selectedIndex = null
  }

  function onNavigate(event) {
    const next = indexRef === -1 ? null : indexRef
    if ((!event && !open) || transitionStatus === "ending") return
    if (event) keyboardActive = keyboardActive // type tracking only
    setIndices({ activeIndex: next })
  }

  const isMainOrientationKey = (key) => key === "ArrowUp" || key === "ArrowDown"
  const isMainOrientationToEndKey = (key) => key === "ArrowDown" || key === "Enter" || key === " " || key === ""

  // useListNavigation's commonOnKeyDown (virtual, vertical, loop, allowEscape).
  function commonOnKeyDown(event) {
    isPointerModality = false
    forceSyncFocus = true
    if (event.which === 229) return
    const filtered = filteredItems()
    const currentIndex = indexRef
    const minIndex = filtered.length ? 0 : -1
    const maxIndex = filtered.length - 1
    const loopFocus = true
    const allowEscape = loopFocus && !autoHighlightMode
    if (isMainOrientationKey(event.key)) {
      event.preventDefault()
      event.stopPropagation()
      if (event.key === "ArrowDown") {
        if (currentIndex >= maxIndex) {
          if (allowEscape && currentIndex !== filtered.length) indexRef = -1
          else indexRef = minIndex
        } else indexRef = currentIndex + 1
      } else if (currentIndex <= minIndex) {
        if (allowEscape && currentIndex !== -1) indexRef = filtered.length
        else indexRef = maxIndex
      } else indexRef = currentIndex - 1
      if (indexRef < 0 || indexRef >= filtered.length) indexRef = -1
      onNavigate(event)
    }
  }

  // --- wiring ---
  function wireInput(el) {
    el.addEventListener("pointerdown", (e) => {
      keyboardActive = false
      // useListNavigation's checkVirtualPointer; useClick pointer type.
      focusItemOnOpenRef = focusItemOnOpenProp()
      lastPointerType = e.pointerType
      if (!open) openMethod = e.pointerType || ""
    })
    el.addEventListener("pointermove", () => (keyboardActive = false))
    el.addEventListener("mousedown", (e) => {
      // useClick ("mousedown-only", toggle off).
      if (inputDisabled(el) || e.button !== 0) return
      setOpen(true, "input-press", e)
    })
    el.addEventListener("focus", () => {})
    el.addEventListener("input", (e) => {
      const value = el.value
      setInputValue(value, "input-change", e)
      const empty = value === ""
      if (empty && !inputInsidePopup) {
        if (single) setSelectedValue(null, "input-clear")
      }
      const trimmed = value.trim()
      const autofillLike = !e.inputType || e.inputType === "insertReplacementText"
      if (!inputDisabled(el) && trimmed && !autofillLike) {
        setOpen(true, "input-change", e)
        if (!autoHighlightMode) setIndices({ activeIndex: null, selectedIndex: null })
      }
      if (open && activeIndex !== null && !autoHighlightMode) setIndices({ activeIndex: null, selectedIndex: null })
      // listNav: the list changed while open.
      if (open && activeIndex == null && selectedIndex == null) indexRef = -1
      scheduleLayout()
      render()
    })
    el.addEventListener("keydown", (e) => onInputKeyDown(e, el))
  }
  let lastPointerType

  const inputDisabled = (el) => comboboxDisabled || (el === input && inputDisabledProp)

  function onInputKeyDown(event, el) {
    if (inputDisabled(el)) return
    // The merged handlers all see the state of the render the event hit.
    const openAtStart = open
    const afterOwnKeyDown = (e) => afterOwnKeyDownWith(e, openAtStart)
    // The input's own handler.
    if (!(event.ctrlKey || event.shiftKey || event.altKey || event.metaKey)) {
      keyboardActive = true
      if (event.key === "Home") {
        event.preventDefault()
        event.stopPropagation()
        el.setSelectionRange(0, 0)
        el.scrollLeft = 0
        return afterOwnKeyDown(event)
      }
      if (event.key === "End") {
        event.preventDefault()
        event.stopPropagation()
        el.setSelectionRange(el.value.length, el.value.length)
        el.scrollLeft = el.scrollWidth - el.clientWidth
        return afterOwnKeyDown(event)
      }
      if (!mounted && event.key === "Escape") {
        const isClear = multiple ? selected.length === 0 : selected === null
        setInputValue("", "escape-key")
        setSelectedValue(multiple ? [] : null, "escape-key")
        if (!isClear) event.stopPropagation()
        return afterOwnKeyDown(event)
      }
      if (chipsEl && event.key === "Backspace" && el.value === "" && highlightedChipIndex === undefined && multiple && selected.length > 0) {
        const count = chips().length
        const removal = count > 0 ? count - 1 : selected.length - 1
        activeIndex = null
        selectedIndex = null
        setSelectedValue(selected.filter((_, i) => i !== removal), "none")
        return afterOwnKeyDown(event)
      }
      if (chipsEl) {
        const had = highlightedChipIndex !== undefined
        const next = chipsKeyDown(event, el)
        highlightedChipIndex = next
        if (next !== undefined) chips()[next]?.focus()
        else if (had) activeInput()?.focus()
      }
      if (event.which !== 229 && event.key === "Enter" && open) {
        if (activeIndex === null) {
          setOpen(false, "none", event)
          return afterOwnKeyDown(event)
        }
        event.preventDefault()
        event.stopPropagation()
        const target = currentItemElements()[activeIndex]
        target?.click()
      }
    }
    afterOwnKeyDown(event)
  }

  // dismiss.reference (Escape), then listNavigation.reference.
  function afterOwnKeyDownWith(event, openAtStart) {
    if (event.key === "Escape" && openAtStart) {
      setOpen(false, "escape-key", event)
      event.preventDefault()
      event.stopPropagation()
      return
    }
    listNavKeyDown(event, openAtStart)
  }

  function listNavKeyDown(event, openAtStart) {
    isPointerModality = false
    if (openAtStart) {
      commonOnKeyDown(event)
      return
    }
    const isMain = isMainOrientationKey(event.key)
    const isNavigationKey = isMain || event.key === "Enter" || event.key.trim() === ""
    if (isNavigationKey) keyRef = event.key
    if (isMain) {
      if (selectedIndex != null) indexRef = selectedIndex
      event.preventDefault()
      event.stopPropagation()
      setOpen(true, "list-navigation", event)
    }
  }

  function chipsKeyDown(event, el) {
    let nextIndex
    const count = chips().length
    const prevKey = "ArrowLeft"
    const nextKey = "ArrowRight"
    if (highlightedChipIndex !== undefined) {
      if (event.key === prevKey) {
        event.preventDefault()
        nextIndex = highlightedChipIndex > 0 ? highlightedChipIndex - 1 : undefined
      } else if (event.key === nextKey) {
        event.preventDefault()
        nextIndex = highlightedChipIndex < count - 1 ? highlightedChipIndex + 1 : undefined
      } else if (event.key === "Backspace" || event.key === "Delete") {
        event.preventDefault()
        const computed = highlightedChipIndex >= selected.length - 1 ? selected.length - 2 : highlightedChipIndex
        nextIndex = computed >= 0 ? computed : undefined
        activeIndex = null
        selectedIndex = null
        render()
      }
      return nextIndex
    }
    if (event.key === prevKey && (el.selectionStart ?? 0) === 0 && selected.length > 0) {
      event.preventDefault()
      nextIndex = count > 0 ? count - 1 : undefined
    } else if (event.key === "Backspace" && el.value === "" && selected.length > 0) {
      activeIndex = null
      selectedIndex = null
      event.preventDefault()
      render()
    }
    return nextIndex
  }

  function currentItemElements() {
    if (!list) return []
    return [...list.querySelectorAll('[role="option"]')]
  }

  function wireItem(el, it) {
    let didPointerDown = false
    el.addEventListener("pointerdown", (e) => {
      didPointerDown = true
      e.preventDefault()
    }, true)
    el.addEventListener("mousedown", (e) => e.preventDefault())
    el.addEventListener("click", (e) => {
      if (comboboxDisabled) return
      handleSelection(e, it)
    })
    el.addEventListener("mouseup", (e) => {
      const started = didPointerDown
      didPointerDown = false
      const idx = currentItemElements().indexOf(el)
      if (comboboxDisabled || e.button !== 0 || started || activeIndex !== idx) return
      handleSelection(e, it)
    })
    el.addEventListener("mousemove", () => {
      forceSyncFocus = true
      forceScrollIntoView = false
      if (!open) return
      const idx = currentItemElements().indexOf(el)
      if (idx !== -1 && (indexRef !== idx || activeIndex !== idx)) {
        indexRef = idx
        onNavigate({})
      }
    })
    el.addEventListener("pointerleave", (e) => {
      if (!open || !isPointerModality || e.pointerType === "touch") return
      if (currentItemElements().includes(e.relatedTarget)) return
      indexRef = -1
      onNavigate({})
    })
  }

  function wireList(el) {
    el.addEventListener("keydown", (e) => {
      if (comboboxDisabled) return
      if (e.key === "Escape" && open) {
        setOpen(false, "escape-key", e)
        e.preventDefault()
        e.stopPropagation()
        return
      }
      if (e.key === "Enter") {
        if (activeIndex == null) return
        e.preventDefault()
        e.stopPropagation()
        currentItemElements()[activeIndex]?.click()
        return
      }
      commonOnKeyDown(e)
    })
    el.addEventListener("keydown", () => (keyboardActive = true), true)
    el.addEventListener("pointermove", () => (keyboardActive = false), true)
    el.addEventListener("pointermove", () => (isPointerModality = true))
  }

  function wireChip(chip) {
    chip.addEventListener("keydown", (event) => {
      if (comboboxDisabled) return
      const index = chips().indexOf(chip)
      let nextIndex = index
      if (event.key === "ArrowLeft") {
        event.preventDefault()
        nextIndex = index > 0 ? index - 1 : undefined
      } else if (event.key === "ArrowRight") {
        event.preventDefault()
        nextIndex = index < chips().length - 1 ? index + 1 : undefined
      } else if (event.key === "Backspace" || event.key === "Delete") {
        const computed = index >= selected.length - 1 ? selected.length - 2 : index
        nextIndex = computed >= 0 ? computed : undefined
        event.preventDefault()
        event.stopPropagation()
        activeIndex = null
        selectedIndex = null
        setSelectedValue(selected.filter((_, i) => i !== index), "none")
      } else if (event.key === "Enter" || event.key === " ") {
        event.preventDefault()
        event.stopPropagation()
        nextIndex = undefined
      } else if (event.key === "ArrowDown" || event.key === "ArrowUp") {
        event.preventDefault()
        event.stopPropagation()
        setOpen(true, "list-navigation", event)
        nextIndex = undefined
      } else if (event.key.length === 1 && !event.ctrlKey && !event.metaKey && !event.altKey) {
        nextIndex = undefined
      }
      highlightedChipIndex = nextIndex
      if (nextIndex === undefined) activeInput()?.focus()
      else chips()[nextIndex]?.focus()
    })
    const remove = chip.querySelector('[data-slot="combobox-chip-remove"]')
    if (remove) {
      const removeChip = (event) => {
        const index = chips().indexOf(chip)
        const removed = selected[index]
        if (activeIndex != null) {
          const ri = filteredItems().indexOf(removed)
          if (ri !== -1 && activeIndex === ri) activeIndex = null
        }
        setSelectedValue(selected.filter((_, i) => i !== index), "chip-remove-press")
        activeInput()?.focus()
        void event
      }
      remove.addEventListener("mousedown", (e) => e.preventDefault())
      remove.addEventListener("click", (e) => {
        if (comboboxDisabled) return
        removeChip(e)
        e.stopPropagation()
      })
      remove.addEventListener("keydown", (e) => {
        if (comboboxDisabled) return
        if (e.key === "Enter" || e.key === " ") {
          removeChip(e)
          e.preventDefault()
          e.stopPropagation()
        }
      })
    }
  }

  function wireClear(el) {
    el.addEventListener("mousedown", (e) => e.preventDefault())
    el.addEventListener("click", () => {
      if (comboboxDisabled || clearDisabledProp) return
      setInputValue("", "clear-press")
      setSelectedValue(multiple ? [] : null, "clear-press")
      activeIndex = null
      selectedIndex = null
      render()
      activeInput()?.focus()
    })
  }

  // useTypeahead on the trigger (single selection, closed). Its list is
  // labelsRef, which forceMount fills with the labels of the filtered
  // items when the trigger is focused or pressed, or the value changes.
  let typeaheadLabels = []
  let typeString = ""
  let typeTimer = 0
  let prevTypeIndex = -1
  let matchTypeIndex = null
  let typeaheadItems = []
  const forceMount = () => {
    typeaheadItems = filteredItems()
    typeaheadLabels = typeaheadItems.map((it) => it.label)
  }
  function typeaheadKeyDown(event) {
    if (open || comboboxDisabled) return
    const list = typeaheadLabels
    const matching = (string, startIndex = 0) => {
      if (list.length === 0) return -1
      const start = ((startIndex % list.length) + list.length) % list.length
      const lower = string.toLowerCase()
      for (let off = 0; off < list.length; off += 1) {
        const i = (start + off) % list.length
        if (list[i]?.toLowerCase().startsWith(lower)) return i
      }
      return -1
    }
    if (typeString.length > 0 && event.key === " ") {
      event.preventDefault()
      event.stopPropagation()
    }
    if (event.key.length !== 1 || event.ctrlKey || event.metaKey || event.altKey) return
    const isNew = typeString === ""
    if (isNew) prevTypeIndex = selectedIndex ?? activeIndex ?? -1
    const rapid = list.every((t) => (t ? t[0]?.toLowerCase() !== t[1]?.toLowerCase() : true))
    if (rapid && typeString === event.key) {
      typeString = ""
      prevTypeIndex = matchTypeIndex
    }
    typeString += event.key
    clearTimeout(typeTimer)
    typeTimer = setTimeout(() => {
      typeString = ""
      prevTypeIndex = matchTypeIndex
    }, TYPEAHEAD_RESET_MS)
    const prev = isNew ? selectedIndex ?? activeIndex ?? -1 : prevTypeIndex
    const index = matching(typeString, (prev ?? 0) + 1)
    if (index !== -1) {
      // The labels and the items they name come from the same render.
      const value = typeaheadItems[index]
      if (value !== undefined) setSelectedValue(value, "none")
      matchTypeIndex = index
    } else if (event.key !== " ") typeString = ""
  }
  function resetTypeaheadOnChange() {
    // useTypeahead's layout effect on [open, selectedIndex].
    if (!open && selectedIndex !== null) return
    clearTimeout(typeTimer)
    matchTypeIndex = null
    typeString = ""
  }

  // Trigger (ComboboxTrigger).
  function wireTrigger(t) {
    let pointerType
    t.addEventListener("pointerdown", (e) => {
      pointerType = e.pointerType
      if (!open) openMethod = e.pointerType || ""
    })
    t.addEventListener("pointerenter", (e) => (pointerType = e.pointerType))
    t.addEventListener("focus", () => {
      if (comboboxDisabled || t.disabled) return
      setTimeout(forceMount, 0)
    })
    t.addEventListener("blur", (e) => {
      if (containsEl(positioner, e.relatedTarget)) return
      clearTimeout(typeTimer)
      typeString = ""
      prevTypeIndex = matchTypeIndex
    })
    t.addEventListener("mousedown", (event) => {
      if (comboboxDisabled || t.disabled) return
      forceMount()
      // The trigger becomes the floating root's DOM reference for this
      // event only: the next render syncs it back to the input.
      if (pointerType !== "touch") {
        activeInput()?.focus()
        if (!inputInsidePopup) event.preventDefault()
      }
      if (event.button !== 0) return
      // useClick (event "mousedown", toggle).
      let nextOpen
      if (!open) nextOpen = true
      else if (openEvent && (openEvent.type === "click" || openEvent.type === "mousedown")) nextOpen = false
      else nextOpen = openEvent ? true : false
      if (!open) {
        const handleMouseUp = (mouseEvent) => {
          if (!inputInsidePopup) return
          const target = mouseEvent.target
          if (containsEl(t, target) || containsEl(positioner, target) || containsEl(list, target)) return
          const r = t.getBoundingClientRect()
          if (mouseEvent.clientX >= r.left - 2 && mouseEvent.clientX <= r.right + 2 && mouseEvent.clientY >= r.top - 2 && mouseEvent.clientY <= r.bottom + 2) return
          setOpen(false, "cancel-open", mouseEvent)
        }
        if (inputInsidePopup) document.addEventListener("mouseup", handleMouseUp, { once: true })
      }
      requestAnimationFrame(() => setOpen(nextOpen, "trigger-press", event))
    })
    t.addEventListener("keydown", (event) => {
      // useClick onKeyDown: a keyboard press is not a pointer press.
      pointerType = undefined
      if (comboboxDisabled || t.disabled) return
      if (inputInsidePopup && single) typeaheadKeyDown(event)
      if (event.key === "ArrowDown" || event.key === "ArrowUp") {
        event.preventDefault()
        event.stopPropagation()
        keyRef = event.key
        setOpen(true, "list-navigation", event)
        activeInput()?.focus()
      }
    })
    if (inputInsidePopup) {
      // Native button click (keyboard Enter/Space) opens through useClick's onClick.
      t.addEventListener("click", (event) => {
        if (comboboxDisabled || t.disabled) return
        if (pointerType) {
          pointerType = undefined
          return
        }
        if (!open) {
          openMethod = event.detail === 0 ? "keyboard" : ""
          setOpen(true, "trigger-press", event)
        } else if (!openEvent || ["click", "mousedown", "keydown", "keyup"].includes(openEvent.type)) {
          // useClick getNextOpen (stickIfOpen): a popup opened by a
          // click-like event closes on the next press; one opened by focus
          // or hover stays open.
          setOpen(false, "trigger-press", event)
        }
      })
    }
  }

  function wireChips(el) {
    el.addEventListener("mousedown", (event) => {
      if (event.target !== el && (event.target.closest("button,input,a,select,textarea,[role=button]"))) return
      event.preventDefault()
      if (comboboxDisabled) return
      activeInput()?.focus()
      setOpen(true, "input-press", event)
    })
  }

  // --- document listeners while open (useDismiss, focus manager) ---
  const onDocKeyDown = (event) => {
    if (!open || event.key !== "Escape") return
    setOpen(false, "escape-key", event)
    event.preventDefault()
  }
  const isWithinOwn = (target) => containsEl(positioner, target) || containsEl(domReference, target)
  const outsidePressAllowed = (target) =>
    !triggers.some((t) => containsEl(t, target)) && !containsEl(clearEl, target) && !containsEl(chipsEl, target)
  const onDocPointerDown = (event) => {
    if (!open) return
    const target = event.target
    if (event.pointerType === "touch") return
    if (isWithinOwn(target)) return
    if (insideReactTree) {
      insideReactTree = false
      return
    }
    // Clicks on the scrollbar of a scroll container.
    if (target instanceof HTMLElement) {
      const style = getComputedStyle(target)
      const canScrollY = /auto|scroll/.test(style.overflowY) && target.clientHeight > 0 && target.scrollHeight > target.clientHeight
      if (canScrollY && event.offsetX > target.clientWidth) return
    }
    if (!outsidePressAllowed(target)) return
    setOpen(false, "outside-press", event)
  }
  const onFocusOut = (event) => {
    const relatedTarget = event.relatedTarget
    const currentTarget = event.currentTarget
    queueMicrotask(() => {
      if (!open) return
      const untrapped = !inputInsidePopup && isTypeable(domReference)
      const modal = !inputInsidePopup
      const insideEls = [startDismiss, endDismiss]
      const isGuard = relatedTarget?.hasAttribute?.("data-base-ui-focus-guard") && guards && Object.values(guards).includes(relatedTarget)
      const moved = !(containsEl(domReference, relatedTarget) || containsEl(positioner, relatedTarget) || containsEl(relatedTarget, positioner) || containsEl(portal, relatedTarget) || insideEls.some((e) => e && containsEl(e, relatedTarget)) || triggers.includes(relatedTarget) || isGuard)
      void currentTarget
      if ((untrapped ? true : !modal) && relatedTarget && moved && !isPointerDown) {
        preventReturnFocus = true
        setOpen(false, "focus-out", event)
      }
    })
  }
  const onRefPointerDown = () => {
    isPointerDown = true
    setTimeout(() => (isPointerDown = false), 0)
  }

  // --- mount ---
  if (input) {
    wireInput(input)
    input.addEventListener("focusout", onFocusOut)
    input.addEventListener("pointerdown", onRefPointerDown)
  }
  for (const t of triggers) {
    wireTrigger(t)
    if (inputInsidePopup) {
      t.addEventListener("focusout", onFocusOut)
      t.addEventListener("pointerdown", onRefPointerDown)
    }
  }
  if (chipsEl) {
    wireChips(chipsEl)
    for (const c of chips()) {
      c.__item = selected[chips().indexOf(c)]
      wireChip(c)
    }
  }
  if (clearEl) wireClear(clearEl)
  document.addEventListener("keydown", onDocKeyDown)
  document.addEventListener("pointerdown", (e) => onDocPointerDown(e), true)
  // Track pointer down inside the floating tree (dismiss's markInsideReactTree).
  const markInside = (e) => {
    if (containsEl(positioner, e.target)) {
      insideReactTree = true
      setTimeout(() => (insideReactTree = false), 0)
    }
  }
  document.addEventListener("mousedown", markInside, true)
  // The popup's onFocus: focus inside the list goes back to the input.
  document.addEventListener("focusin", (event) => {
    if (!popup) return
    const target = event.target
    if (inputInsidePopup && openMethod !== "touch" && (containsEl(list, target) || target === popup)) popupInput?.focus()
    else if (!inputInsidePopup && (containsEl(list, target) || target === popup)) input?.focus()
  })
  // The popup's focusout (non-modal focus manager closes on focus out).
  document.addEventListener(
    "focusout",
    (event) => {
      if (positioner && containsEl(positioner, event.target)) onFocusOut(event)
    },
    true
  )
  syncSelectedIndex()
  render()

  return () => {
    unmountPopup()
    unmarkOthers?.()
    document.removeEventListener("keydown", onDocKeyDown)
  }
}

define("combobox", {
  selector: '[data-combobox-part="hidden"]',
  mount: createCombobox,
})
