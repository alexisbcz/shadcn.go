// menu.js: port of Radix Menu (@radix-ui/react-menu 2.1.16), the primitive
// behind DropdownMenu, ContextMenu and Menubar: content mounting and
// presence, positioning, dismissal, focus management, keyboard navigation
// and typeahead, items (select on click, Enter/Space, pointer up after a
// press elsewhere), checkbox and radio items, submenus (open on hover after
// 100ms with a pointer grace area towards the submenu, on click or with
// ArrowRight/ArrowLeft in RTL; close with ArrowLeft, Escape closes the
// whole menu).
//
// Markup contract (rendered by components/ui/menu):
//   - The root content is a <template data-content-id data-trigger-id
//     data-dir [data-modal="false"] [data-default-open] [data-loop]>
//     holding the popper wrapper (see popper.js) with the content as React
//     renders it when open. The component module finds it and calls
//     createMenu.
//   - Submenus: inside a content, <template data-radix-menu-sub
//     data-content-id data-trigger-id [data-portal] [data-default-open]
//     [data-loop]> where React renders the SubContent; the sub trigger is
//     the item whose id is data-trigger-id.
//   - Markers removed when a content is cloned: data-radix-menu-checkbox
//     (checkbox item state key), data-radix-menu-radio-group (radio group
//     state key) with data-radix-menu-radio-value (its value; also on radio
//     items), data-radix-menu-controlled, data-radix-menu-indicator
//     (ItemIndicator, "force" when force-mounted), data-radix-menu-text-value.
//
// API (for component modules)
//
//   import { createMenu } from "./menu.js"
//   const menu = createMenu(template, {
//     modal: true, dir: "ltr",
//     anchor: () => triggerOrVirtualElement,   // popper anchor
//     layerParent: () => trigger,              // dismissable layer parent
//     onOpenChange(open) {},  // the menu asks to open/close (dismiss, select)
//     // MenuContent props, composed like Radix composes them:
//     onCloseAutoFocus(e), onOpenAutoFocus(e), onEscapeKeyDown(e),
//     onPointerDownOutside(e), onFocusOutside(e), onInteractOutside(e),
//     onEntryFocus(e), onKeyDown(e),
//   })
//   menu.setOpen(open)   // the component's state changed (controlled `open`)
//   menu.open / menu.content (the live content element, or null)
//   menu.destroy()
//
// Events dispatched on live elements: "shadcn:select" (item, cancelable:
// keeps the menu open), "shadcn:checked-change" (checkbox item, cancelable,
// detail.checked), "shadcn:value-change" (radio group, cancelable,
// detail.value). getMenuState / setChecked / setRadioValue let scripts read
// and change checked state.

import { emit } from "./core.js"
import { createPresence } from "./presence.js"
import { instantiate, mount, takeMarked } from "./portal.js"
import { createDismissableLayer } from "./dismissable-layer.js"
import { createFocusScope } from "./focus-scope.js"
import { acquireFocusGuards } from "./focus-guards.js"
import { lockScroll } from "./scroll-lock.js"
import { hideOthers } from "./hide-others.js"
import { createPopper } from "./popper.js"
import { createRovingFocusGroup } from "./roving-focus.js"
import { createCollection } from "./collection.js"

const SELECTION_KEYS = ["Enter", " "]
const FIRST_KEYS = ["ArrowDown", "PageUp", "Home"]
const LAST_KEYS = ["ArrowUp", "PageDown", "End"]
const FIRST_LAST_KEYS = [...FIRST_KEYS, ...LAST_KEYS]
const SUB_OPEN_KEYS = { ltr: [...SELECTION_KEYS, "ArrowRight"], rtl: [...SELECTION_KEYS, "ArrowLeft"] }
const SUB_CLOSE_KEYS = { ltr: ["ArrowLeft"], rtl: ["ArrowRight"] }

const M_CHECKBOX = "data-radix-menu-checkbox"
const M_RADIO_GROUP = "data-radix-menu-radio-group"
const M_RADIO_VALUE = "data-radix-menu-radio-value"
const M_CONTROLLED = "data-radix-menu-controlled"
const M_INDICATOR = "data-radix-menu-indicator"
const M_TEXT_VALUE = "data-radix-menu-text-value"

// Menu's isUsingKeyboardRef: a keydown anywhere sets it until the next
// pointer down or move. Every Menu root registers the same listeners, so
// one flag serves them all.
let isUsingKeyboard = false
const handlePointer = () => (isUsingKeyboard = false)
document.addEventListener(
  "keydown",
  () => {
    isUsingKeyboard = true
    document.addEventListener("pointerdown", handlePointer, { capture: true, once: true })
    document.addEventListener("pointermove", handlePointer, { capture: true, once: true })
  },
  { capture: true },
)

/** Whether the keyboard was used last (Radix's isUsingKeyboardRef). */
export const usingKeyboard = () => isUsingKeyboard

// Checked state outlives the content (it is the page's React state):
// checkbox key -> true | false | "indeterminate", radio group key -> value.
const checkedState = new Map()
const radioState = new Map()
// Live elements per key, for re-rendering state changes.
const liveCheckboxes = new Map() // key -> Set<{ item, indicator }>
const liveRadioGroups = new Map() // key -> Set<{ group, items: [{ item, value, indicator }] }>

/** Returns a checkbox item's state or a radio group's value by key. */
export function getMenuState(key) {
  return checkedState.has(key) ? checkedState.get(key) : radioState.get(key)
}

const getOpenState = (open) => (open ? "open" : "closed")
const isIndeterminate = (checked) => checked === "indeterminate"
const getCheckedState = (checked) => (isIndeterminate(checked) ? "indeterminate" : checked ? "checked" : "unchecked")
// React's synthetic events remember preventDefault() even on events that are
// not cancelable (pointerleave): track it on the event.
const PREVENTED = Symbol("menu.prevented")
const prevent = (event) => {
  event[PREVENTED] = true
  event.preventDefault()
}
const isPrevented = (event) => event.defaultPrevented || !!event[PREVENTED]
const whenMouse = (handler) => (event) => (event.pointerType === "mouse" ? handler(event) : undefined)

/** Sets a checkbox item's state (true, false or "indeterminate") by key. */
export function setChecked(key, checked) {
  checkedState.set(key, checked)
  for (const live of liveCheckboxes.get(key) ?? []) renderCheckbox(live, checked)
}

/** Sets a radio group's value by key. */
export function setRadioValue(key, value) {
  radioState.set(key, value)
  for (const live of liveRadioGroups.get(key) ?? []) {
    for (const r of live.items) renderIndicatorItem(r, r.value === value, String(r.value === value))
  }
}

function renderCheckbox(live, checked) {
  live.item.setAttribute("aria-checked", isIndeterminate(checked) ? "mixed" : String(!!checked))
  renderIndicatorItem(live, isIndeterminate(checked) || checked === true, null, getCheckedState(checked))
}

// Re-renders an item's data-state (and aria-checked when given) and its
// ItemIndicator's presence. (Presence: shadcn's indicators have no exit
// animation, so they unmount at once.)
function renderIndicatorItem(live, present, aria, stateValue) {
  const st = stateValue ?? (present ? "checked" : "unchecked")
  if (aria !== null && aria !== undefined) live.item.setAttribute("aria-checked", aria)
  live.item.setAttribute("data-state", st)
  const ind = live.indicator
  if (!ind) return
  ind.el.setAttribute("data-state", st)
  const show = present || ind.force
  const attached = ind.el.parentNode === ind.parent
  if (show && !attached) ind.parent.insertBefore(ind.el, ind.next && ind.next.parentNode === ind.parent ? ind.next : null)
  else if (!show && attached) ind.el.remove()
}

// Radix's focusFirst (menu flavor: no preventScroll).
function focusFirst(candidates) {
  const PREVIOUSLY_FOCUSED_ELEMENT = document.activeElement
  for (const candidate of candidates) {
    if (candidate === PREVIOUSLY_FOCUSED_ELEMENT) return
    candidate.focus()
    if (document.activeElement !== PREVIOUSLY_FOCUSED_ELEMENT) return
  }
}

function wrapArray(array, startIndex) {
  return array.map((_, index) => array[(startIndex + index) % array.length])
}

function getNextMatch(values, search, currentMatch) {
  const isRepeated = search.length > 1 && Array.from(search).every((char) => char === search[0])
  const normalizedSearch = isRepeated ? search[0] : search
  const currentMatchIndex = currentMatch ? values.indexOf(currentMatch) : -1
  let wrappedValues = wrapArray(values, Math.max(currentMatchIndex, 0))
  const excludeCurrentMatch = normalizedSearch.length === 1
  if (excludeCurrentMatch) wrappedValues = wrappedValues.filter((v) => v !== currentMatch)
  const nextMatch = wrappedValues.find((value) => value.toLowerCase().startsWith(normalizedSearch.toLowerCase()))
  return nextMatch !== currentMatch ? nextMatch : undefined
}

function isPointInPolygon(point, polygon) {
  const { x, y } = point
  let inside = false
  for (let i = 0, j = polygon.length - 1; i < polygon.length; j = i++) {
    const xi = polygon[i].x
    const yi = polygon[i].y
    const xj = polygon[j].x
    const yj = polygon[j].y
    const intersect = yi > y !== yj > y && x < ((xj - xi) * (y - yi)) / (yj - yi) + xi
    if (intersect) inside = !inside
  }
  return inside
}

function isPointerInGraceArea(event, area) {
  if (!area) return false
  return isPointInPolygon({ x: event.clientX, y: event.clientY }, area)
}

// Runs handlers like composeEventHandlers: ours only if the user's did not
// prevent the default (unless check is false).
function compose(user, ours, check = true) {
  return (event) => {
    user?.(event)
    if (!check || !event.defaultPrevented) ours?.(event)
  }
}

/**
 * Creates a root menu (Radix Menu + MenuContent) for a content template.
 * @param {HTMLTemplateElement} template
 * @param {object} props see the file header
 */
export function createMenu(template, props = {}) {
  const root = {
    open: false,
    modal: props.modal ?? template.getAttribute("data-modal") !== "false",
    dir: props.dir ?? (template.getAttribute("data-dir") === "rtl" ? "rtl" : "ltr"),
    onOpenChange: (open) => props.onOpenChange?.(open),
    onClose: () => props.onOpenChange?.(false),
  }
  // MenuProvider of the root: { open, content }
  const menuCtx = { open: false, content: null, listeners: new Set() }
  let mounted = null

  const contentProps = props
  const setOpen = (open) => {
    open = !!open
    if (open === menuCtx.open) return
    menuCtx.open = open
    root.open = open
    if (open) {
      if (mounted) mounted.reshow()
      else mounted = mountContent({ template, root, menuCtx, isSub: false, contentProps, onUnmounted: () => (mounted = null) })
    } else if (mounted) {
      mounted.hide()
    }
    for (const l of menuCtx.listeners) l(open)
  }

  return {
    get open() {
      return menuCtx.open
    },
    get content() {
      return menuCtx.content
    },
    setOpen,
    destroy() {
      mounted?.destroy()
      mounted = null
    },
  }
}

// MenuContent / MenuSubContent mounting (Presence + MenuContentImpl).
function mountContent({ template, root, menuCtx, isSub, contentProps, sub, onUnmounted }) {
  const nodes = instantiate(template)
  const wrapper = nodes[0]
  const content = wrapper.firstElementChild
  const loop = template.hasAttribute("data-loop")
  const cleanups = []
  const state = { content, wrapper }

  // Markers inside this content (not inside nested templates).
  const textValues = new Map()
  for (const { el, value } of takeMarked(content, M_TEXT_VALUE)) textValues.set(el, value)
  const controlled = new Set(takeMarked(content, M_CONTROLLED).map((m) => m.el))
  const indicators = new Map() // item -> indicator record
  for (const { el, value } of takeMarked(content, M_INDICATOR)) {
    const item = el.closest('[role="menuitemcheckbox"], [role="menuitemradio"]')
    if (item) indicators.set(item, { el, parent: el.parentNode, next: el.nextSibling, force: value === "force" })
  }
  const checkboxes = takeMarked(content, M_CHECKBOX)
  const radioGroups = takeMarked(content, M_RADIO_GROUP)
  const radioValues = new Map(takeMarked(content, M_RADIO_VALUE).map(({ el, value }) => [el, value]))

  // Checkbox items: initial state from the server, then the page's state.
  const checkboxRecords = new Map() // item -> { key, live }
  for (const { el: item, value: key } of checkboxes) {
    if (!checkedState.has(key)) {
      const aria = item.getAttribute("aria-checked")
      checkedState.set(key, aria === "mixed" ? "indeterminate" : aria === "true")
    }
    const live = { item, indicator: indicators.get(item) }
    if (!liveCheckboxes.has(key)) liveCheckboxes.set(key, new Set())
    liveCheckboxes.get(key).add(live)
    cleanups.push(() => liveCheckboxes.get(key)?.delete(live))
    checkboxRecords.set(item, { key, live, controlled: controlled.has(item) })
    renderCheckbox(live, checkedState.get(key))
  }
  // Radio groups and their items.
  const radioRecords = new Map() // item -> { key, value, group, controlled }
  const groupKeys = new Map(radioGroups.map(({ el, value }) => [el, value]))
  const groupItems = new Map([...groupKeys.keys()].map((g) => [g, []]))
  for (const item of content.querySelectorAll('[role="menuitemradio"]')) {
    let group = item.parentElement
    while (group && group !== content && !groupKeys.has(group)) group = group.parentElement
    if (!group || !groupKeys.has(group)) continue
    const value = radioValues.get(item) ?? ""
    groupItems.get(group).push({ item, value, indicator: indicators.get(item) })
    radioRecords.set(item, { key: groupKeys.get(group), value, group, controlled: controlled.has(group) })
  }
  for (const [group, key] of groupKeys) {
    if (!radioState.has(key)) radioState.set(key, radioValues.get(group) ?? "")
    const live = { group, items: groupItems.get(group) }
    if (!liveRadioGroups.has(key)) liveRadioGroups.set(key, new Set())
    liveRadioGroups.get(key).add(live)
    cleanups.push(() => liveRadioGroups.get(key)?.delete(live))
    const current = radioState.get(key)
    for (const r of live.items) renderIndicatorItem(r, r.value === current, String(r.value === current))
  }
  // Indicators of items without state markers keep the server state.
  for (const [item, ind] of indicators) {
    if (checkboxRecords.has(item) || radioRecords.has(item)) continue
    const st = item.getAttribute("data-state")
    if (st !== "checked" && st !== "indeterminate" && !ind.force) ind.el.remove()
  }

  // Portal (root content and portalled subs) or in place (inline subs).
  let unmount
  if (!isSub || template.hasAttribute("data-portal")) unmount = mount(nodes)
  else {
    template.before(...nodes)
    unmount = () => nodes.forEach((n) => n.remove())
  }
  menuCtx.content = content

  // --- MenuContentImpl state ---
  let searchTimer = 0
  let search = ""
  let pointerGraceTimer = 0
  let pointerGraceIntent = null
  let pointerDir = "right"
  let lastPointerX = 0
  cleanups.push(() => window.clearTimeout(searchTimer))

  const collection = createCollection(content)
  cleanups.push(() => collection.destroy())
  const getItems = () =>
    collection.getItems().map((el) => ({
      element: el,
      disabled: el.hasAttribute("data-disabled"),
      textValue: textValues.has(el) ? textValues.get(el) : (el.textContent ?? "").trim(),
    }))

  const isPointerMovingToSubmenu = (event) => {
    const isMovingTowards = pointerDir === pointerGraceIntent?.side
    return isMovingTowards && isPointerInGraceArea(event, pointerGraceIntent?.area)
  }
  const contentCtx = {
    get search() {
      return search
    },
    onItemEnter(event) {
      if (isPointerMovingToSubmenu(event)) prevent(event)
    },
    onItemLeave(event) {
      if (isPointerMovingToSubmenu(event)) return
      content.focus()
      roving.setCurrentTabStop(null)
    },
    onTriggerLeave(event) {
      if (isPointerMovingToSubmenu(event)) prevent(event)
    },
    get pointerGraceTimer() {
      return pointerGraceTimer
    },
    set pointerGraceTimer(v) {
      pointerGraceTimer = v
    },
    onPointerGraceIntentChange(intent) {
      pointerGraceIntent = intent
    },
  }

  const handleTypeaheadSearch = (key) => {
    const next = search + key
    const items = getItems().filter((item) => !item.disabled)
    const currentItem = document.activeElement
    const currentMatch = items.find((item) => item.element === currentItem)?.textValue
    const values = items.map((item) => item.textValue)
    const nextMatch = getNextMatch(values, next, currentMatch)
    const newItem = items.find((item) => item.textValue === nextMatch)?.element
    const updateSearch = (value) => {
      search = value
      window.clearTimeout(searchTimer)
      if (value !== "") searchTimer = window.setTimeout(() => updateSearch(""), 1000)
    }
    updateSearch(next)
    if (newItem) setTimeout(() => newItem.focus())
  }

  // --- submenus of this content ---
  const subs = new Map() // sub trigger -> sub controller
  for (const tpl of content.querySelectorAll("template[data-radix-menu-sub]")) {
    if (tpl.parentElement?.closest("[data-radix-menu-content]") !== content) continue
    const trig = content.querySelector(`#${CSS.escape(tpl.getAttribute("data-trigger-id"))}`)
    const s = createSub({ template: tpl, trigger: trig, root, parentCtx: menuCtx, parentState: state, contentCtx })
    if (trig) subs.set(trig, s)
    cleanups.push(() => s.destroy())
  }

  contentCtx.wireItemImpl = wireItemImpl

  // --- items (listeners on the items run before the group's) ---
  for (const item of collection.getItems()) {
    const sub = subs.get(item)
    cleanups.push(sub ? sub.wireTrigger(item) : wireItem(item))
  }

  function wireItem(item) {
    const disabled = item.hasAttribute("data-disabled")
    let isPointerDown = false
    const handleSelect = () => {
      if (disabled) return
      const selected = emit(item, "shadcn:select", {})
      // CheckboxItem / RadioItem onSelect (composed without the default check).
      const cb = checkboxRecords.get(item)
      if (cb) {
        const checked = checkedState.get(cb.key)
        const next = isIndeterminate(checked) ? true : !checked
        if (emit(item, "shadcn:checked-change", { checked: next }) && !cb.controlled) setChecked(cb.key, next)
      }
      const rr = radioRecords.get(item)
      if (rr) {
        if (emit(rr.group, "shadcn:value-change", { value: rr.value }) && !rr.controlled) setRadioValue(rr.key, rr.value)
      }
      if (!selected) isPointerDown = false
      else root.onClose()
    }
    const onClick = (event) => {
      if (!event.defaultPrevented) handleSelect()
    }
    const onPointerDown = () => (isPointerDown = true)
    const onPointerUp = (event) => {
      if (event.defaultPrevented) return
      if (!isPointerDown) item.click()
    }
    const onKeyDown = (event) => {
      if (event.defaultPrevented) return
      const isTypingAhead = search !== ""
      if (disabled || (isTypingAhead && event.key === " ")) return
      if (SELECTION_KEYS.includes(event.key)) {
        item.click()
        event.preventDefault()
      }
    }
    item.addEventListener("click", onClick)
    item.addEventListener("pointerdown", onPointerDown)
    item.addEventListener("pointerup", onPointerUp)
    item.addEventListener("keydown", onKeyDown)
    const unimpl = wireItemImpl(item, disabled)
    return () => {
      item.removeEventListener("click", onClick)
      item.removeEventListener("pointerdown", onPointerDown)
      item.removeEventListener("pointerup", onPointerUp)
      item.removeEventListener("keydown", onKeyDown)
      unimpl()
    }
  }

  // MenuItemImpl: pointer focus, highlight.
  function wireItemImpl(item, disabled, before = {}) {
    const onPointerMove = (event) => {
      before.pointermove?.(event)
      if (isPrevented(event) || event.pointerType !== "mouse") return
      if (disabled) contentCtx.onItemLeave(event)
      else {
        contentCtx.onItemEnter(event)
        if (!isPrevented(event)) item.focus({ preventScroll: true })
      }
    }
    const onPointerLeave = (event) => {
      before.pointerleave?.(event)
      if (isPrevented(event) || event.pointerType !== "mouse") return
      contentCtx.onItemLeave(event)
    }
    const onFocus = () => item.setAttribute("data-highlighted", "")
    const onBlur = () => item.removeAttribute("data-highlighted")
    item.addEventListener("pointermove", onPointerMove)
    item.addEventListener("pointerleave", onPointerLeave)
    item.addEventListener("focus", onFocus)
    item.addEventListener("blur", onBlur)
    return () => {
      item.removeEventListener("pointermove", onPointerMove)
      item.removeEventListener("pointerleave", onPointerLeave)
      item.removeEventListener("focus", onFocus)
      item.removeEventListener("blur", onBlur)
    }
  }

  // --- RovingFocusGroup (asChild on the content) ---
  const roving = createRovingFocusGroup(content, {
    orientation: "vertical",
    dir: root.dir,
    loop,
    focusable: (item) => !item.hasAttribute("data-disabled"),
    onEntryFocus: compose(contentProps.onEntryFocus, (event) => {
      if (!isUsingKeyboard) event.preventDefault()
    }),
    preventScrollOnEntryFocus: true,
    groupTabIndex: -1,
  })
  cleanups.push(() => roving.destroy())

  // --- content handlers (after the items' and the group's) ---
  const ownKeyDown = (event) => {
    const target = event.target
    const isKeyDownInside = target.closest("[data-radix-menu-content]") === content
    const isModifierKey = event.ctrlKey || event.altKey || event.metaKey
    const isCharacterKey = event.key.length === 1
    if (isKeyDownInside) {
      if (event.key === "Tab") event.preventDefault()
      if (!isModifierKey && isCharacterKey) handleTypeaheadSearch(event.key)
    }
    if (event.target !== content) return
    if (!FIRST_LAST_KEYS.includes(event.key)) return
    event.preventDefault()
    const candidateNodes = getItems()
      .filter((item) => !item.disabled)
      .map((item) => item.element)
    if (LAST_KEYS.includes(event.key)) candidateNodes.reverse()
    focusFirst(candidateNodes)
  }
  const userKeyDown = isSub ? compose(contentProps.onKeyDown, sub.onKeyDown) : contentProps.onKeyDown
  const onKeyDown = compose(userKeyDown, ownKeyDown)
  // React bubbles a portalled submenu's events through its parent content.
  state.reactKeyDown = (event) => {
    onKeyDown(event)
    if (isSub && !content.parentElement?.closest("[data-radix-menu-content]")) sub.parentContentKeyDown?.(event)
  }
  const onKeyDownListener = (event) => state.reactKeyDown(event)
  const onBlur = (event) => {
    if (!content.contains(event.target)) {
      window.clearTimeout(searchTimer)
      search = ""
    }
  }
  const onPointerMove = whenMouse((event) => {
    const target = event.target
    const pointerXHasChanged = lastPointerX !== event.clientX
    if (content.contains(target) && pointerXHasChanged) {
      pointerDir = event.clientX > lastPointerX ? "right" : "left"
      lastPointerX = event.clientX
    }
  })
  content.addEventListener("keydown", onKeyDownListener)
  content.addEventListener("focusout", onBlur)
  content.addEventListener("pointermove", onPointerMove)
  cleanups.push(() => {
    content.removeEventListener("keydown", onKeyDownListener)
    content.removeEventListener("focusout", onBlur)
    content.removeEventListener("pointermove", onPointerMove)
  })

  // --- effects: PopperContent, DismissableLayer, FocusScope, RemoveScroll,
  // focus guards, hideOthers (modal root) ---
  const popper = createPopper(wrapper, { anchor: isSub ? sub.trigger : contentProps.anchor?.() })
  cleanups.push(() => popper.destroy())
  state.popper = popper

  const modal = !isSub && root.modal
  const layer = createDismissableLayer(content, {
    disableOutsidePointerEvents: modal && menuCtx.open,
    parent: isSub ? sub.trigger : contentProps.layerParent?.(),
    onEscapeKeyDown: isSub
      ? compose(contentProps.onEscapeKeyDown, (event) => {
          root.onClose()
          event.preventDefault()
        })
      : contentProps.onEscapeKeyDown,
    onPointerDownOutside: contentProps.onPointerDownOutside,
    onFocusOutside: isSub
      ? compose(contentProps.onFocusOutside, (event) => {
          if (event.target !== sub.trigger) sub.setOpen(false)
        })
      : modal
        ? compose(contentProps.onFocusOutside, (event) => event.preventDefault(), false)
        : contentProps.onFocusOutside,
    onInteractOutside: contentProps.onInteractOutside,
    onDismiss: isSub ? undefined : () => root.onOpenChange(false),
  })
  cleanups.push(() => layer.destroy())
  state.layer = layer

  const focusScope = createFocusScope(content, {
    trapped: modal && menuCtx.open,
    onMountAutoFocus: isSub
      ? (event) => {
          if (isUsingKeyboard) content.focus()
          event.preventDefault()
        }
      : compose(contentProps.onOpenAutoFocus, (event) => {
          event.preventDefault()
          content.focus({ preventScroll: true })
        }),
    onUnmountAutoFocus: isSub ? (event) => event.preventDefault() : (event) => contentProps.onCloseAutoFocus?.(event),
  })
  cleanups.push(() => focusScope.destroy())
  state.focusScope = focusScope

  if (modal) cleanups.push(lockScroll(content, { allowPinchZoom: true }))
  cleanups.push(acquireFocusGuards())
  if (modal) cleanups.push(hideOthers(content))

  // Default-open submenus open with their menu.
  for (const s of subs.values()) s.mountDefault()

  let destroyed = false
  const finish = () => {
    if (destroyed) return
    destroyed = true
    presence.destroy()
    unmount()
    for (const c of cleanups.reverse()) c()
    if (menuCtx.content === content) menuCtx.content = null
    onUnmounted?.()
  }
  const presence = createPresence(content, { present: true, onExitComplete: finish })

  state.hide = () => {
    content.setAttribute("data-state", "closed")
    if (modal) {
      focusScope.setTrapped(false)
      layer.setDisableOutsidePointerEvents(false)
    }
    presence.setPresent(false)
  }
  state.reshow = () => {
    content.setAttribute("data-state", "open")
    if (modal) {
      focusScope.setTrapped(true)
      layer.setDisableOutsidePointerEvents(true)
    }
    presence.setPresent(true)
  }
  state.destroy = finish
  return state
}

// MenuSub + MenuSubTrigger + MenuSubContent.
function createSub({ template, trigger, root, parentCtx, parentState, contentCtx }) {
  const menuCtx = { open: false, content: null, listeners: new Set() }
  let mounted = null
  let openTimer = null
  const defaultOpen = template.hasAttribute("data-default-open")

  const clearOpenTimer = () => {
    if (openTimer) window.clearTimeout(openTimer)
    openTimer = null
  }

  const sub = {
    trigger,
    setOpen,
    onKeyDown: (event) => {
      // MenuSubContent's onKeyDown
      const content = menuCtx.content
      const isKeyDownInside = content?.contains(event.target)
      const isCloseKey = SUB_CLOSE_KEYS[root.dir].includes(event.key)
      if (isKeyDownInside && isCloseKey) {
        setOpen(false)
        trigger?.focus()
        event.preventDefault()
      }
    },
    // Keydowns in a portalled sub content bubble (React tree) to the
    // parent content's handlers.
    parentContentKeyDown: (event) => parentState.reactKeyDown?.(event),
  }

  function setOpen(open) {
    open = !!open
    if (open === menuCtx.open) return
    menuCtx.open = open
    if (trigger) {
      trigger.setAttribute("aria-expanded", String(open))
      trigger.setAttribute("data-state", getOpenState(open))
    }
    if (open) {
      if (mounted) mounted.reshow()
      else {
        mounted = mountContent({ template, root, menuCtx, isSub: true, contentProps: {}, sub, onUnmounted: () => (mounted = null) })
      }
    } else if (mounted) {
      mounted.hide()
    }
    for (const l of menuCtx.listeners) l(open)
  }

  // MenuSub effect: close when the parent menu closes.
  const onParentOpenChange = (open) => {
    if (open === false) setOpen(false)
  }
  parentCtx.listeners.add(onParentOpenChange)

  function wireTrigger(item) {
    const disabled = item.hasAttribute("data-disabled")
    const onClick = (event) => {
      if (disabled || event.defaultPrevented) return
      item.focus()
      if (!menuCtx.open) setOpen(true)
    }
    const before = {
      pointermove: (event) => {
        if (event.pointerType !== "mouse") return
        contentCtx.onItemEnter(event)
        if (isPrevented(event)) return
        if (!disabled && !menuCtx.open && !openTimer) {
          contentCtx.onPointerGraceIntentChange(null)
          openTimer = window.setTimeout(() => {
            setOpen(true)
            clearOpenTimer()
          }, 100)
        }
      },
      pointerleave: (event) => {
        if (event.pointerType !== "mouse") return
        clearOpenTimer()
        const contentRect = menuCtx.content?.getBoundingClientRect()
        if (contentRect) {
          const side = menuCtx.content?.dataset.side
          const rightSide = side === "right"
          const bleed = rightSide ? -5 : 5
          const contentNearEdge = contentRect[rightSide ? "left" : "right"]
          const contentFarEdge = contentRect[rightSide ? "right" : "left"]
          contentCtx.onPointerGraceIntentChange({
            area: [
              { x: event.clientX + bleed, y: event.clientY },
              { x: contentNearEdge, y: contentRect.top },
              { x: contentFarEdge, y: contentRect.top },
              { x: contentFarEdge, y: contentRect.bottom },
              { x: contentNearEdge, y: contentRect.bottom },
            ],
            side,
          })
          window.clearTimeout(contentCtx.pointerGraceTimer)
          contentCtx.pointerGraceTimer = window.setTimeout(() => contentCtx.onPointerGraceIntentChange(null), 300)
        } else {
          contentCtx.onTriggerLeave(event)
          if (isPrevented(event)) return
          contentCtx.onPointerGraceIntentChange(null)
        }
      },
    }
    const onKeyDown = (event) => {
      if (event.defaultPrevented) return
      const isTypingAhead = contentCtx.search !== ""
      if (disabled || (isTypingAhead && event.key === " ")) return
      if (SUB_OPEN_KEYS[root.dir].includes(event.key)) {
        // context.content is the content rendered before this update (none
        // when the submenu was closed).
        const content = menuCtx.content
        setOpen(true)
        content?.focus()
        event.preventDefault()
      }
    }
    item.addEventListener("click", onClick)
    item.addEventListener("keydown", onKeyDown)
    const unimpl = contentCtx.wireItemImpl(item, disabled, before)
    return () => {
      item.removeEventListener("click", onClick)
      item.removeEventListener("keydown", onKeyDown)
      unimpl()
    }
  }

  return {
    wireTrigger,
    mountDefault() {
      if (defaultOpen) setOpen(true)
    },
    get open() {
      return menuCtx.open
    },
    destroy() {
      clearOpenTimer()
      window.clearTimeout(contentCtx.pointerGraceTimer)
      contentCtx.onPointerGraceIntentChange(null)
      parentCtx.listeners.delete(onParentOpenChange)
      // MenuSub unmount: onOpenChange(false); its content unmounts with the
      // parent content at once.
      menuCtx.open = false
      if (mounted) mounted.destroy()
      mounted = null
    },
  }
}
