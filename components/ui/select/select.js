// select.js: Radix Select behavior (@radix-ui/react-select 2.2.6) for the
// server-rendered markup of components/ui/select.
//
// Markup contract (rendered by the Go package):
//   - <template data-radix-select data-content-id="<id>" data-position=
//       "item-aligned"|"popper" [data-value="<v>"] [data-default-open]
//       [data-disabled]> holds the content as React renders it when open:
//       item-aligned: <div style="display:flex;..."> (the content wrapper)
//       popper:       the popper wrapper (see popper.js)
//     around the listbox, plus an optional <span data-radix-select-placeholder>
//     with the placeholder nodes of the Value.
//   - The trigger is the [role=combobox] with aria-controls="<id>".
//   - Markers removed on mount: data-radix-select-value="<id>" (SelectValue
//     span, inside the trigger), data-radix-select-value-children (the Value
//     has its own children), data-radix-select-bubble="<id>" (the native
//     <select> for forms, after the trigger), and inside the template:
//     data-radix-select-item="<value>", data-radix-select-text-value,
//     data-radix-select-item-indicator (ItemIndicator, kept only on the
//     selected item), data-radix-select-scroll-up / -down (scroll buttons,
//     inserted only while the viewport can scroll that way).
//
// Scripts can observe and drive a select:
//   template.addEventListener("shadcn:open-change", (e) => e.detail.open)   // cancelable
//   template.addEventListener("shadcn:value-change", (e) => e.detail.value)
//   import { getSelect } from "/_shadcn/select.js"
//   getSelect(contentIdOrTemplateOrTrigger)?.setValue("apple")   // .value, .open, .setOpen(b)

import { define, emit } from "./core.js"
import { instantiate, mount, takeMarked } from "./portal.js"
import { createDismissableLayer } from "./dismissable-layer.js"
import { createFocusScope } from "./focus-scope.js"
import { acquireFocusGuards } from "./focus-guards.js"
import { lockScroll } from "./scroll-lock.js"
import { hideOthers } from "./hide-others.js"
import { createPopper, openOnLoad } from "./popper.js"
import { getNonce } from "./style.js"

const OPEN_KEYS = [" ", "Enter", "ArrowUp", "ArrowDown"]
const SELECTION_KEYS = [" ", "Enter"]
const CONTENT_MARGIN = 10
const VIEWPORT_CSS =
  "[data-radix-select-viewport]{scrollbar-width:none;-ms-overflow-style:none;-webkit-overflow-scrolling:touch;}[data-radix-select-viewport]::-webkit-scrollbar{display:none}"

const selects = new Map() // content id -> controller

/** Returns the controller of a select by content id, template or trigger. */
export function getSelect(ref) {
  if (typeof ref === "string") return selects.get(ref)
  if (ref instanceof Element) {
    const id = ref.getAttribute("data-content-id") ?? ref.getAttribute("aria-controls")
    return id ? selects.get(id) : undefined
  }
  return undefined
}

const clamp = (value, [min, max]) => Math.min(max, Math.max(min, value))
const shouldShowPlaceholder = (value) => value === "" || value === undefined

// useTypeaheadSearch
function createTypeahead(onSearchChange) {
  let search = ""
  let timer = 0
  const update = (value) => {
    search = value
    window.clearTimeout(timer)
    if (value !== "") timer = window.setTimeout(() => update(""), 1000)
  }
  return {
    get search() {
      return search
    },
    handle(key) {
      const next = search + key
      onSearchChange(next)
      update(next)
    },
    reset() {
      search = ""
      window.clearTimeout(timer)
    },
    destroy() {
      window.clearTimeout(timer)
    },
  }
}

function findNextItem(items, search, currentItem) {
  const isRepeated = search.length > 1 && Array.from(search).every((char) => char === search[0])
  const normalizedSearch = isRepeated ? search[0] : search
  const currentItemIndex = currentItem ? items.indexOf(currentItem) : -1
  let wrappedItems = wrapArray(items, Math.max(currentItemIndex, 0))
  const excludeCurrentItem = normalizedSearch.length === 1
  if (excludeCurrentItem) wrappedItems = wrappedItems.filter((v) => v !== currentItem)
  const nextItem = wrappedItems.find((item) => item.textValue.toLowerCase().startsWith(normalizedSearch.toLowerCase()))
  return nextItem !== currentItem ? nextItem : undefined
}

function wrapArray(array, startIndex) {
  return array.map((_, index) => array[(startIndex + index) % array.length])
}

// Reads an item's data from server-rendered markup (template or live).
function readItem(el, root) {
  const value = el.getAttribute("data-radix-select-item")
  const textId = el.getAttribute("aria-labelledby")
  const textEl = (textId && [...el.querySelectorAll("[id]")].find((n) => n.id === textId)) || null
  const tv = el.getAttribute("data-radix-select-text-value")
  return {
    el,
    value,
    disabled: el.hasAttribute("data-disabled"),
    textValue: tv ?? (textEl?.textContent ?? "").trim(),
    textEl,
  }
}

function createSelect(tpl) {
  const contentId = tpl.getAttribute("data-content-id")
  const position = tpl.getAttribute("data-position") === "popper" ? "popper" : "item-aligned"
  const rootDisabled = tpl.hasAttribute("data-disabled")
  let value = tpl.hasAttribute("data-value") ? tpl.getAttribute("data-value") : undefined
  let open = false
  let mounted = null
  let triggerPointerDownPos = null

  const esc = CSS.escape(contentId)
  const trigger = () => document.querySelector(`[role="combobox"][aria-controls="${esc}"]`)
  // The Select's direction (Radix context.dir), not the trigger's own dir
  // attribute, which a dir prop on SelectTrigger overrides.
  const dir = () => {
    const d = tpl.getAttribute("data-dir") ?? trigger()?.getAttribute("dir")
    return d === "rtl" ? "rtl" : "ltr"
  }

  // SelectValue
  const valueNode = takeMarked([...document.querySelectorAll(`[data-radix-select-value="${esc}"]`)], "data-radix-select-value")[0]?.el ?? null
  const valueNodeHasChildren = !!valueNode?.hasAttribute("data-radix-select-value-children")
  valueNode?.removeAttribute("data-radix-select-value-children")
  const valueChildren = valueNodeHasChildren ? [...valueNode.childNodes] : []
  const placeholderEl = tpl.content.querySelector("[data-radix-select-placeholder]")

  // Items as the closed content renders them (Radix keeps them in a
  // detached fragment): typeahead on the trigger and the value text.
  const templateItems = () => [...tpl.content.querySelectorAll("[data-radix-select-item]")].map((el) => readItem(el))

  // SelectBubbleInput: kept only for form controls.
  let bubble = null
  for (const { el } of takeMarked([...document.querySelectorAll(`select[data-radix-select-bubble="${esc}"]`)], "data-radix-select-bubble")) {
    const t = trigger()
    const isFormControl = t ? !!el.getAttribute("form") || !!t.closest("form") : true
    if (isFormControl) bubble = el
    else el.remove()
  }
  const onBubbleChange = (event) => {
    if (event.target.value !== (value ?? "")) setValue(event.target.value)
  }
  bubble?.addEventListener("change", onBubbleChange)

  const isDisabled = () => rootDisabled || !!trigger()?.disabled

  function renderValueNode() {
    if (!valueNode) return
    if (shouldShowPlaceholder(value)) {
      valueNode.replaceChildren(...(placeholderEl ? [...placeholderEl.cloneNode(true).childNodes] : []))
    } else if (valueNodeHasChildren) {
      valueNode.replaceChildren(...valueChildren)
    } else {
      const item = templateItems().find((it) => it.value === value)
      valueNode.replaceChildren(...(item?.textEl ? [...item.textEl.cloneNode(true).childNodes] : []))
    }
  }

  function setValue(next) {
    if (next === value) return
    value = next
    const t = trigger()
    if (t) {
      if (shouldShowPlaceholder(value)) t.setAttribute("data-placeholder", "")
      else t.removeAttribute("data-placeholder")
    }
    renderValueNode()
    if (bubble) {
      // React re-keys the native select when its options change; the value
      // effect sets the value and dispatches change.
      const empty = bubble.querySelector('option[value=""]')
      if (value === undefined && !empty) bubble.prepend(new Option("", ""))
      else if (value !== undefined && empty && !empty.textContent) empty.remove()
      const setter = Object.getOwnPropertyDescriptor(HTMLSelectElement.prototype, "value").set
      setter.call(bubble, value ?? "")
      bubble.dispatchEvent(new Event("change", { bubbles: true }))
    }
    mounted?.syncItems()
    emit(tpl, "shadcn:value-change", { value })
  }

  function setOpen(next) {
    next = !!next
    if (next === open) return
    if (!emit(tpl, "shadcn:open-change", { open: next })) return
    open = next
    const t = trigger()
    if (t) {
      t.setAttribute("aria-expanded", String(open))
      t.setAttribute("data-state", open ? "open" : "closed")
    }
    if (open) mountContent()
    else if (mounted) mounted.unmount()
  }

  // --- trigger --------------------------------------------------------------------------

  let triggerPointerType = "touch"
  const triggerTypeahead = createTypeahead((search) => {
    const enabledItems = templateItems().filter((item) => !item.disabled)
    const currentItem = enabledItems.find((item) => item.value === value)
    const nextItem = findNextItem(enabledItems, search, currentItem)
    if (nextItem !== undefined) setValue(nextItem.value)
  })
  const handleOpen = (pointerEvent) => {
    // Radix records the position after requesting the open state, but the
    // content (and its pointer-up effect) mounts later: record it first.
    if (pointerEvent) {
      triggerPointerDownPos = { x: Math.round(pointerEvent.pageX), y: Math.round(pointerEvent.pageY) }
    }
    if (!isDisabled()) {
      setOpen(true)
      triggerTypeahead.reset()
    }
  }
  const onTriggerClick = (event) => {
    if (event.defaultPrevented) return
    event.currentTarget.focus()
    if (triggerPointerType !== "mouse") handleOpen(event)
  }
  const onTriggerPointerDown = (event) => {
    if (event.defaultPrevented) return
    triggerPointerType = event.pointerType
    const target = event.target
    if (target.hasPointerCapture(event.pointerId)) target.releasePointerCapture(event.pointerId)
    if (event.button === 0 && event.ctrlKey === false && event.pointerType === "mouse") {
      handleOpen(event)
      event.preventDefault()
    }
  }
  const onTriggerKeyDown = (event) => {
    if (event.defaultPrevented) return
    const isTypingAhead = triggerTypeahead.search !== ""
    const isModifierKey = event.ctrlKey || event.altKey || event.metaKey
    if (!isModifierKey && event.key.length === 1) triggerTypeahead.handle(event.key)
    if (isTypingAhead && event.key === " ") return
    if (OPEN_KEYS.includes(event.key)) {
      handleOpen()
      event.preventDefault()
    }
  }
  const boundTrigger = trigger()
  boundTrigger?.addEventListener("click", onTriggerClick)
  boundTrigger?.addEventListener("pointerdown", onTriggerPointerDown)
  boundTrigger?.addEventListener("keydown", onTriggerKeyDown)

  // --- content --------------------------------------------------------------------------

  function mountContent() {
    const nodes = instantiate(tpl).filter((n) => !n.hasAttribute("data-radix-select-placeholder"))
    const wrapper = nodes[0]
    const content = wrapper.firstElementChild
    const viewport = content.querySelector("[data-radix-select-viewport]")
    const scrollUp = takeMarked(content, "data-radix-select-scroll-up")[0]?.el ?? null
    const scrollDown = takeMarked(content, "data-radix-select-scroll-down")[0]?.el ?? null
    scrollUp?.remove()
    scrollDown?.remove()
    const items = [...content.querySelectorAll("[data-radix-select-item]")].map((el) => {
      const it = readItem(el)
      el.removeAttribute("data-radix-select-item")
      el.removeAttribute("data-radix-select-text-value")
      it.indicator = takeMarked(el, "data-radix-select-item-indicator")[0]?.el ?? null
      it.indicatorParent = it.indicator?.parentNode ?? null
      it.focused = false
      it.pointerType = "touch"
      return it
    })

    // SelectViewport's <style>, with the page's CSP nonce.
    const style = document.createElement("style")
    const nonce = getNonce()
    if (nonce) style.setAttribute("nonce", nonce)
    style.textContent = VIEWPORT_CSS
    viewport.before(style)

    const state = { cleanups: [], destroyed: false }
    let selectedItem = null
    let selectedItemText = null
    let isPositioned = false
    let shouldExpandOnScroll = false
    let shouldReposition = true
    let canScrollUp = false
    let canScrollDown = false
    let autoScrollTimer = null

    const syncItems = () => {
      for (const it of items) {
        const isSelected = value === it.value
        it.el.setAttribute("data-state", isSelected ? "checked" : "unchecked")
        it.el.setAttribute("aria-selected", String(isSelected && it.focused))
        if (it.indicator) {
          if (isSelected && !it.indicator.parentNode) it.indicatorParent.appendChild(it.indicator)
          else if (!isSelected && it.indicator.parentNode) it.indicator.remove()
        }
      }
    }
    state.syncItems = syncItems
    syncItems()

    // itemRefCallback / itemTextRefCallback: the selected item, else the
    // first enabled one (the last of the two in document order).
    let firstValidFound = false
    for (const it of items) {
      const isFirstValidItem = !firstValidFound && !it.disabled
      const isSelectedItem = value !== undefined && value === it.value
      if (isSelectedItem || isFirstValidItem) {
        selectedItem = it.el
        selectedItemText = it.textEl
        if (isFirstValidItem) firstValidFound = true
      }
    }

    const unmountPortal = mount(nodes)
    mounted = state

    const t = trigger()
    const getItems = () => items.filter((it) => it.el.isConnected)

    const focusFirst = (candidates) => {
      const [firstItem, ...restItems] = getItems().map((item) => item.el)
      const [lastItem] = restItems.slice(-1)
      const PREVIOUSLY_FOCUSED_ELEMENT = document.activeElement
      for (const candidate of candidates) {
        if (candidate === PREVIOUSLY_FOCUSED_ELEMENT) return
        candidate?.scrollIntoView({ block: "nearest" })
        if (candidate === firstItem && viewport) viewport.scrollTop = 0
        if (candidate === lastItem && viewport) viewport.scrollTop = viewport.scrollHeight
        candidate?.focus()
        if (document.activeElement !== PREVIOUSLY_FOCUSED_ELEMENT) return
      }
    }
    const focusSelectedItem = () => focusFirst([selectedItem, content])
    const onItemLeave = () => content.focus()

    // --- items --------------------------------------------------------------------------

    const contentTypeahead = createTypeahead((search) => {
      const enabledItems = getItems().filter((item) => !item.disabled)
      const currentItem = enabledItems.find((item) => item.el === document.activeElement)
      const nextItem = findNextItem(enabledItems, search, currentItem)
      if (nextItem) setTimeout(() => nextItem.el.focus())
    })
    state.cleanups.push(() => contentTypeahead.destroy())

    for (const it of items) {
      const el = it.el
      const handleSelect = () => {
        if (!it.disabled) {
          setValue(it.value)
          setOpen(false)
        }
      }
      el.addEventListener("focus", () => {
        it.focused = true
        el.setAttribute("data-highlighted", "")
        el.setAttribute("aria-selected", String(value === it.value))
      })
      el.addEventListener("blur", () => {
        it.focused = false
        el.removeAttribute("data-highlighted")
        el.setAttribute("aria-selected", "false")
      })
      el.addEventListener("click", (event) => {
        if (event.defaultPrevented) return
        if (it.pointerType !== "mouse") handleSelect()
      })
      el.addEventListener("pointerup", (event) => {
        if (event.defaultPrevented) return
        if (it.pointerType === "mouse") handleSelect()
      })
      el.addEventListener("pointerdown", (event) => {
        if (event.defaultPrevented) return
        it.pointerType = event.pointerType
      })
      el.addEventListener("pointermove", (event) => {
        if (event.defaultPrevented) return
        it.pointerType = event.pointerType
        if (it.disabled) onItemLeave()
        else if (it.pointerType === "mouse") el.focus({ preventScroll: true })
      })
      el.addEventListener("pointerleave", (event) => {
        if (event.defaultPrevented) return
        if (el === document.activeElement) onItemLeave()
      })
      el.addEventListener("keydown", (event) => {
        if (event.defaultPrevented) return
        const isTypingAhead = contentTypeahead.search !== ""
        if (isTypingAhead && event.key === " ") return
        if (SELECTION_KEYS.includes(event.key)) handleSelect()
        if (event.key === " ") event.preventDefault()
      })
    }

    content.addEventListener("contextmenu", (event) => event.preventDefault())
    content.addEventListener("keydown", (event) => {
      // composeEventHandlers: the item's Space handling prevents ours.
      if (event.defaultPrevented || state.destroyed) return
      const isModifierKey = event.ctrlKey || event.altKey || event.metaKey
      if (event.key === "Tab") event.preventDefault()
      if (!isModifierKey && event.key.length === 1) contentTypeahead.handle(event.key)
      if (["ArrowUp", "ArrowDown", "Home", "End"].includes(event.key)) {
        const enabled = getItems().filter((item) => !item.disabled)
        let candidateNodes = enabled.map((item) => item.el)
        if (["ArrowUp", "End"].includes(event.key)) candidateNodes = candidateNodes.slice().reverse()
        if (["ArrowUp", "ArrowDown"].includes(event.key)) {
          const currentIndex = candidateNodes.indexOf(event.target)
          candidateNodes = candidateNodes.slice(currentIndex + 1)
        }
        setTimeout(() => focusFirst(candidateNodes))
        event.preventDefault()
      }
    })

    // SelectViewport onScroll: expand the item-aligned content while
    // scrolling, and the scroll buttons' scroll listeners.
    let prevScrollTop = 0
    viewport.addEventListener("scroll", () => {
      if (position === "item-aligned" && shouldExpandOnScroll) {
        const scrolledBy = Math.abs(prevScrollTop - viewport.scrollTop)
        if (scrolledBy > 0) {
          const availableHeight = window.innerHeight - CONTENT_MARGIN * 2
          const cssMinHeight = parseFloat(wrapper.style.minHeight)
          const cssHeight = parseFloat(wrapper.style.height)
          const prevHeight = Math.max(cssMinHeight, cssHeight)
          if (prevHeight < availableHeight) {
            const nextHeight = prevHeight + scrolledBy
            const clampedNextHeight = Math.min(availableHeight, nextHeight)
            const heightDiff = nextHeight - clampedNextHeight
            wrapper.style.height = clampedNextHeight + "px"
            if (wrapper.style.bottom === "0px") {
              viewport.scrollTop = heightDiff > 0 ? heightDiff : 0
              wrapper.style.justifyContent = "flex-end"
            }
          }
        }
      }
      prevScrollTop = viewport.scrollTop
      updateScrollButtons()
    })


    // Effects of the first commit (children first): DismissableLayer,
    // FocusScope, RemoveScroll, focus guards.
    const layer = createDismissableLayer(content, {
      disableOutsidePointerEvents: true,
      parent: t,
      onFocusOutside: (event) => event.preventDefault(),
      onDismiss: () => setOpen(false),
    })
    state.cleanups.push(() => layer.destroy())

    const focusScope = createFocusScope(content, {
      trapped: true,
      onMountAutoFocus: (event) => event.preventDefault(),
      onUnmountAutoFocus: (event) => {
        trigger()?.focus({ preventScroll: true })
        event.preventDefault()
      },
    })
    state.cleanups.push(() => focusScope.destroy())
    state.cleanups.push(lockScroll(content, { allowPinchZoom: true }))
    state.cleanups.push(acquireFocusGuards())

    // --- scroll buttons -----------------------------------------------------------------

    const clearAutoScrollTimer = () => {
      if (autoScrollTimer !== null) {
        window.clearInterval(autoScrollTimer)
        autoScrollTimer = null
      }
    }
    state.cleanups.push(clearAutoScrollTimer)
    const bindScrollButton = (btn, onAutoScroll) => {
      if (!btn) return
      const start = () => {
        if (autoScrollTimer === null) autoScrollTimer = window.setInterval(onAutoScroll, 50)
      }
      btn.addEventListener("pointerdown", start)
      btn.addEventListener("pointermove", () => {
        onItemLeave()
        start()
      })
      btn.addEventListener("pointerleave", clearAutoScrollTimer)
    }
    bindScrollButton(scrollUp, () => {
      if (viewport && selectedItem) viewport.scrollTop = viewport.scrollTop - selectedItem.offsetHeight
    })
    bindScrollButton(scrollDown, () => {
      if (viewport && selectedItem) viewport.scrollTop = viewport.scrollTop + selectedItem.offsetHeight
    })

    // Mounts / unmounts the scroll buttons like their canScroll state.
    const updateScrollButtons = () => {
      if (!isPositioned || state.destroyed) return
      const nextUp = viewport.scrollTop > 0
      const maxScroll = viewport.scrollHeight - viewport.clientHeight
      const nextDown = Math.ceil(viewport.scrollTop) < maxScroll
      const mountedButtons = []
      if (scrollUp && nextUp !== canScrollUp) {
        canScrollUp = nextUp
        if (nextUp) {
          content.insertBefore(scrollUp, content.firstChild)
          mountedButtons.push(scrollUp)
        } else {
          if (scrollUp.isConnected) clearAutoScrollTimer()
          scrollUp.remove()
        }
      }
      if (scrollDown && nextDown !== canScrollDown) {
        canScrollDown = nextDown
        if (nextDown) {
          content.appendChild(scrollDown)
          mountedButtons.push(scrollDown)
        } else {
          if (scrollDown.isConnected) clearAutoScrollTimer()
          scrollDown.remove()
        }
      }
      if (!mountedButtons.length) return
      // Ref callbacks (onScrollButtonChange), then the buttons' layout effect.
      if (position === "item-aligned" && shouldReposition) {
        positionItemAligned()
        focusSelectedItem()
        shouldReposition = false
      }
      for (const _ of mountedButtons) {
        const activeItem = getItems().find((item) => item.el === document.activeElement)
        activeItem?.el.scrollIntoView({ block: "nearest" })
      }
    }

    // --- positioning --------------------------------------------------------------------

    const onPlaced = () => {
      if (isPositioned) return
      isPositioned = true
      updateScrollButtons()
      focusSelectedItem()
    }

    function positionItemAligned() {
      const t = trigger()
      if (!(t && valueNode && wrapper && content && viewport && selectedItem && selectedItemText)) return false
      const contentWrapper = wrapper
      const triggerRect = t.getBoundingClientRect()
      const contentRect = content.getBoundingClientRect()
      const valueNodeRect = valueNode.getBoundingClientRect()
      const itemTextRect = selectedItemText.getBoundingClientRect()
      if (dir() !== "rtl") {
        const itemTextOffset = itemTextRect.left - contentRect.left
        const left = valueNodeRect.left - itemTextOffset
        const leftDelta = triggerRect.left - left
        const minContentWidth = triggerRect.width + leftDelta
        const contentWidth = Math.max(minContentWidth, contentRect.width)
        const rightEdge = window.innerWidth - CONTENT_MARGIN
        const clampedLeft = clamp(left, [CONTENT_MARGIN, Math.max(CONTENT_MARGIN, rightEdge - contentWidth)])
        contentWrapper.style.minWidth = minContentWidth + "px"
        contentWrapper.style.left = clampedLeft + "px"
      } else {
        const itemTextOffset = contentRect.right - itemTextRect.right
        const right = window.innerWidth - valueNodeRect.right - itemTextOffset
        const rightDelta = window.innerWidth - triggerRect.right - right
        const minContentWidth = triggerRect.width + rightDelta
        const contentWidth = Math.max(minContentWidth, contentRect.width)
        const leftEdge = window.innerWidth - CONTENT_MARGIN
        const clampedRight = clamp(right, [CONTENT_MARGIN, Math.max(CONTENT_MARGIN, leftEdge - contentWidth)])
        contentWrapper.style.minWidth = minContentWidth + "px"
        contentWrapper.style.right = clampedRight + "px"
      }
      const all = getItems()
      const availableHeight = window.innerHeight - CONTENT_MARGIN * 2
      const itemsHeight = viewport.scrollHeight
      const contentStyles = window.getComputedStyle(content)
      const contentBorderTopWidth = parseInt(contentStyles.borderTopWidth, 10)
      const contentPaddingTop = parseInt(contentStyles.paddingTop, 10)
      const contentBorderBottomWidth = parseInt(contentStyles.borderBottomWidth, 10)
      const contentPaddingBottom = parseInt(contentStyles.paddingBottom, 10)
      const fullContentHeight = contentBorderTopWidth + contentPaddingTop + itemsHeight + contentPaddingBottom + contentBorderBottomWidth
      const minContentHeight = Math.min(selectedItem.offsetHeight * 5, fullContentHeight)
      const viewportStyles = window.getComputedStyle(viewport)
      const viewportPaddingTop = parseInt(viewportStyles.paddingTop, 10)
      const viewportPaddingBottom = parseInt(viewportStyles.paddingBottom, 10)
      const topEdgeToTriggerMiddle = triggerRect.top + triggerRect.height / 2 - CONTENT_MARGIN
      const triggerMiddleToBottomEdge = availableHeight - topEdgeToTriggerMiddle
      const selectedItemHalfHeight = selectedItem.offsetHeight / 2
      const itemOffsetMiddle = selectedItem.offsetTop + selectedItemHalfHeight
      const contentTopToItemMiddle = contentBorderTopWidth + contentPaddingTop + itemOffsetMiddle
      const itemMiddleToContentBottom = fullContentHeight - contentTopToItemMiddle
      const willAlignWithoutTopOverflow = contentTopToItemMiddle <= topEdgeToTriggerMiddle
      if (willAlignWithoutTopOverflow) {
        const isLastItem = all.length > 0 && selectedItem === all[all.length - 1].el
        contentWrapper.style.bottom = "0px"
        const viewportOffsetBottom = content.clientHeight - viewport.offsetTop - viewport.offsetHeight
        const clampedTriggerMiddleToBottomEdge = Math.max(
          triggerMiddleToBottomEdge,
          selectedItemHalfHeight + (isLastItem ? viewportPaddingBottom : 0) + viewportOffsetBottom + contentBorderBottomWidth,
        )
        const height = contentTopToItemMiddle + clampedTriggerMiddleToBottomEdge
        contentWrapper.style.height = height + "px"
      } else {
        const isFirstItem = all.length > 0 && selectedItem === all[0].el
        contentWrapper.style.top = "0px"
        const clampedTopEdgeToTriggerMiddle = Math.max(
          topEdgeToTriggerMiddle,
          contentBorderTopWidth + viewport.offsetTop + (isFirstItem ? viewportPaddingTop : 0) + selectedItemHalfHeight,
        )
        const height = clampedTopEdgeToTriggerMiddle + itemMiddleToContentBottom
        contentWrapper.style.height = height + "px"
        viewport.scrollTop = contentTopToItemMiddle - topEdgeToTriggerMiddle + viewport.offsetTop
      }
      contentWrapper.style.margin = `${CONTENT_MARGIN}px 0`
      contentWrapper.style.minHeight = minContentHeight + "px"
      contentWrapper.style.maxHeight = availableHeight + "px"
      requestAnimationFrame(() => (shouldExpandOnScroll = true))
      return true
    }

    let popperInstance = null
    if (position === "popper") {
      popperInstance = createPopper(wrapper, { anchor: t, onPlaced })
      state.cleanups.push(() => popperInstance.destroy())
    } else {
      wrapper.style.zIndex = window.getComputedStyle(content).zIndex
      if (positionItemAligned()) onPlaced()
    }

    // SelectContentImpl's effects once content is known.
    state.cleanups.push(hideOthers(content))

    if (triggerPointerDownPos !== null) {
      let pointerMoveDelta = { x: 0, y: 0 }
      const handlePointerMove = (event) => {
        pointerMoveDelta = {
          x: Math.abs(Math.round(event.pageX) - (triggerPointerDownPos?.x ?? 0)),
          y: Math.abs(Math.round(event.pageY) - (triggerPointerDownPos?.y ?? 0)),
        }
      }
      const handlePointerUp = (event) => {
        if (pointerMoveDelta.x <= 10 && pointerMoveDelta.y <= 10) {
          event.preventDefault()
        } else if (!content.contains(event.target)) {
          setOpen(false)
        }
        document.removeEventListener("pointermove", handlePointerMove)
        triggerPointerDownPos = null
      }
      document.addEventListener("pointermove", handlePointerMove)
      document.addEventListener("pointerup", handlePointerUp, { capture: true, once: true })
      state.cleanups.push(() => {
        document.removeEventListener("pointermove", handlePointerMove)
        document.removeEventListener("pointerup", handlePointerUp, { capture: true })
      })
    }

    const close = () => setOpen(false)
    window.addEventListener("blur", close)
    window.addEventListener("resize", close)
    state.cleanups.push(() => {
      window.removeEventListener("blur", close)
      window.removeEventListener("resize", close)
    })

    state.unmount = () => {
      if (state.destroyed) return
      state.destroyed = true
      unmountPortal()
      for (const c of state.cleanups.reverse()) c()
      if (mounted === state) mounted = null
    }
  }

  const controller = {
    get open() {
      return open
    },
    get value() {
      return value
    },
    setOpen,
    setValue,
    dispose() {
      boundTrigger?.removeEventListener("click", onTriggerClick)
      boundTrigger?.removeEventListener("pointerdown", onTriggerPointerDown)
      boundTrigger?.removeEventListener("keydown", onTriggerKeyDown)
      bubble?.removeEventListener("change", onBubbleChange)
      triggerTypeahead.destroy()
      mounted?.unmount()
      if (selects.get(contentId) === controller) selects.delete(contentId)
    },
  }
  selects.set(contentId, controller)
  if (tpl.hasAttribute("data-default-open")) openOnLoad(tpl, () => setOpen(true))
  return controller
}

define("select", {
  selector: "template[data-radix-select]",
  mount(tpl) {
    const s = createSelect(tpl)
    return () => s.dispose()
  },
})
