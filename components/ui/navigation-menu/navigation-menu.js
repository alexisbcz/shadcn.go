// navigation-menu.js: Radix NavigationMenu behavior
// (@radix-ui/react-navigation-menu 1.2.14) for the server-rendered markup of
// components/ui/navigation-menu.
//
// Markup contract (rendered by the Go package):
//   - <nav data-radix-navigation-menu='{"delayDuration":200,
//       "skipDelayDuration":300,"value":""}' data-viewport="true|false" ...>
//     (marker removed on mount).
//   - Items: <li data-radix-navigation-menu-item="<value>"> (marker removed)
//     holding a trigger (<button aria-controls="<content id>">) and/or links,
//     and the panel in <template data-radix-navigation-menu-content="<value>">
//     exactly as React renders it while open. Without a viewport the panel is
//     cloned into the item after the trigger; with one, into the viewport.
//   - Optional <template data-radix-navigation-menu-indicator> (the
//     Indicator, portalled into the list's track) and
//     <template data-radix-navigation-menu-viewport> (the Viewport, inside
//     its positioning wrapper).
//
// Scripts can observe and drive a menu:
//   nav.addEventListener("shadcn:value-change", (e) => e.detail.value)
//   import { getNavigationMenu } from "/_shadcn/navigation-menu.js"
//   getNavigationMenu(nav)?.setValue("<item value>")   // "" closes; .value
// Links dispatch Radix's cancelable "navigationMenu.linkSelect" event when
// clicked; preventing it keeps the menu open.

import { define, emit } from "./core.js"
import { createPresence } from "./presence.js"
import { instantiate } from "./portal.js"
import { createDismissableLayer } from "./dismissable-layer.js"

const ROOT_CONTENT_DISMISS = "navigationMenu.rootContentDismiss"
const LINK_SELECT = "navigationMenu.linkSelect"
const ARROW_KEYS = ["ArrowRight", "ArrowLeft", "ArrowUp", "ArrowDown"]
const VISUALLY_HIDDEN =
  "position: absolute; border: 0px; width: 1px; height: 1px; padding: 0px; margin: -1px; overflow: hidden; clip: rect(0px, 0px, 0px, 0px); white-space: nowrap; overflow-wrap: normal;"

const menus = new WeakMap() // nav -> controller

/** Returns the controller of a navigation menu by its <nav> (or an element inside it). */
export function getNavigationMenu(el) {
  const nav = el?.closest?.('nav[data-slot="navigation-menu"]') ?? el
  return menus.get(nav)
}

function getTabbableCandidates(container) {
  const nodes = []
  const walker = document.createTreeWalker(container, NodeFilter.SHOW_ELEMENT, {
    acceptNode: (node) => {
      const isHiddenInput = node.tagName === "INPUT" && node.type === "hidden"
      if (node.disabled || node.hidden || isHiddenInput) return NodeFilter.FILTER_SKIP
      return node.tabIndex >= 0 ? NodeFilter.FILTER_ACCEPT : NodeFilter.FILTER_SKIP
    },
  })
  while (walker.nextNode()) nodes.push(walker.currentNode)
  return nodes
}

function focusFirst(candidates) {
  const previouslyFocusedElement = document.activeElement
  return candidates.some((candidate) => {
    if (candidate === previouslyFocusedElement) return true
    candidate.focus()
    return document.activeElement !== previouslyFocusedElement
  })
}

function removeFromTabOrder(candidates) {
  candidates.forEach((candidate) => {
    candidate.dataset.tabindex = candidate.getAttribute("tabindex") || ""
    candidate.setAttribute("tabindex", "-1")
  })
  return () => {
    candidates.forEach((candidate) => {
      const prevTabIndex = candidate.dataset.tabindex
      candidate.setAttribute("tabindex", prevTabIndex)
    })
  }
}

// useResizeObserver: onResize runs in the animation frame after a resize.
function observeResize(element, onResize) {
  let rAF = 0
  const ro = new ResizeObserver(() => {
    cancelAnimationFrame(rAF)
    rAF = window.requestAnimationFrame(onResize)
  })
  ro.observe(element)
  return () => {
    window.cancelAnimationFrame(rAF)
    ro.unobserve(element)
  }
}

const whenMouse = (handler) => (event) => (event.pointerType === "mouse" ? handler(event) : undefined)

function createNavigationMenu(nav) {
  let opts = {}
  try {
    opts = JSON.parse(nav.getAttribute("data-radix-navigation-menu") || "{}")
  } catch {}
  nav.removeAttribute("data-radix-navigation-menu")
  const delayDuration = opts.delayDuration ?? 200
  const skipDelayDuration = opts.skipDelayDuration ?? 300
  const orientation = nav.getAttribute("data-orientation") === "vertical" ? "vertical" : "horizontal"
  const dir = nav.getAttribute("dir") === "rtl" ? "rtl" : "ltr"
  const ownedBy = (el) => el.closest('nav[data-slot="navigation-menu"]') === nav

  let value = opts.value ?? ""
  let previousValue = value
  let isOpenDelayed = true
  let openTimer = 0
  let closeTimer = 0
  let skipDelayTimer = 0
  const cleanups = []

  const list = [...nav.querySelectorAll('ul[data-slot="navigation-menu-list"]')].find(ownedBy) ?? null
  const indicatorTrack = list?.parentElement ?? null

  // --- items -------------------------------------------------------------------------------

  const items = [...nav.querySelectorAll("li[data-radix-navigation-menu-item]")].filter(ownedBy).map((li) => {
    const itemValue = li.getAttribute("data-radix-navigation-menu-item")
    li.removeAttribute("data-radix-navigation-menu-item")
    const template = [...li.querySelectorAll("template[data-radix-navigation-menu-content]")].find(
      (t) => t.getAttribute("data-radix-navigation-menu-content") === itemValue,
    )
    const contentId = template?.content.firstElementChild?.id
    const trigger =
      [...li.querySelectorAll('[data-slot="navigation-menu-trigger"]')].find((b) => b.closest("li") === li) ?? null
    return {
      li,
      value: itemValue,
      template: template ?? null,
      contentId,
      trigger,
      content: null, // the mounted panel: { el, presence, layer, prevMotion, cleanups }
      focusProxy: null,
      ownsSpan: null,
      wasEscapeClose: false,
      wasClickClose: false,
      hasPointerMoveOpened: false,
      restoreContentTabOrder: () => {},
    }
  })
  // NavigationMenu's Collection: the items with a trigger.
  const triggerItems = () => items.filter((it) => it.trigger)

  // --- viewport ------------------------------------------------------------------------------

  const viewportTemplate = [...nav.querySelectorAll("template[data-radix-navigation-menu-viewport]")].find(ownedBy) ?? null
  let viewport = null // { el, presence, cleanupSize, sizeContent }

  // --- value ---------------------------------------------------------------------------------

  function setValue(next) {
    if (next === value) return
    const isOpen = next !== ""
    if (isOpen) {
      window.clearTimeout(skipDelayTimer)
      if (skipDelayDuration > 0) isOpenDelayed = false
    } else {
      window.clearTimeout(skipDelayTimer)
      skipDelayTimer = window.setTimeout(() => (isOpenDelayed = true), skipDelayDuration)
    }
    previousValue = value
    value = next
    render()
    emit(nav, "shadcn:value-change", { value })
  }

  const startCloseTimer = () => {
    window.clearTimeout(closeTimer)
    closeTimer = window.setTimeout(() => setValue(""), 150)
  }
  const handleOpen = (itemValue) => {
    window.clearTimeout(closeTimer)
    setValue(itemValue)
  }
  const handleDelayedOpen = (itemValue) => {
    const isOpenItem = value === itemValue
    if (isOpenItem) {
      window.clearTimeout(closeTimer)
    } else {
      openTimer = window.setTimeout(() => {
        window.clearTimeout(closeTimer)
        setValue(itemValue)
      }, delayDuration)
    }
  }
  const onTriggerEnter = (itemValue) => {
    window.clearTimeout(openTimer)
    if (isOpenDelayed) handleDelayedOpen(itemValue)
    else handleOpen(itemValue)
  }
  const onTriggerLeave = () => {
    window.clearTimeout(openTimer)
    startCloseTimer()
  }
  const onContentEnter = () => window.clearTimeout(closeTimer)
  const onContentLeave = startCloseTimer
  const onItemSelect = (itemValue) => setValue(value === itemValue ? "" : itemValue)
  const onItemDismiss = () => setValue("")

  // --- content -------------------------------------------------------------------------------

  const handleContentEntry = (it, side = "start") => {
    const content = it.content?.el
    if (content) {
      it.restoreContentTabOrder()
      const candidates = getTabbableCandidates(content)
      if (candidates.length) focusFirst(side === "start" ? candidates : candidates.reverse())
    }
  }
  const handleContentExit = (it) => {
    const content = it.content?.el
    if (content) {
      const candidates = getTabbableCandidates(content)
      if (candidates.length) it.restoreContentTabOrder = removeFromTabOrder(candidates)
    }
  }

  function motionAttribute(it) {
    const values = triggerItems().map((item) => item.value)
    if (dir === "rtl") values.reverse()
    const index = values.indexOf(value)
    const prevIndex = values.indexOf(previousValue)
    const isSelected = it.value === value
    const wasSelected = prevIndex === values.indexOf(it.value)
    if (!isSelected && !wasSelected) return it.content.prevMotion
    let attribute = null
    if (index !== prevIndex) {
      if (isSelected && prevIndex !== -1) attribute = index > prevIndex ? "from-end" : "from-start"
      else if (wasSelected && index !== -1) attribute = index > prevIndex ? "to-start" : "to-end"
    }
    it.content.prevMotion = attribute
    return attribute
  }

  // NavigationMenuContentImpl: mounts a panel from its template.
  function mountContent(it, parent, before) {
    const [el] = instantiate(it.template)
    const state = { el, presence: null, layer: null, prevMotion: null, cleanups: [] }
    it.content = state
    parent.insertBefore(el, before ?? null)

    const onDismissEvent = () => {
      onItemDismiss()
      handleContentExit(it)
      if (el.contains(document.activeElement)) it.trigger?.focus()
    }
    el.addEventListener(ROOT_CONTENT_DISMISS, onDismissEvent)
    state.cleanups.push(() => el.removeEventListener(ROOT_CONTENT_DISMISS, onDismissEvent))

    state.layer = createDismissableLayer(el, {
      disableOutsidePointerEvents: false,
      onDismiss: () => {
        el.dispatchEvent(new Event(ROOT_CONTENT_DISMISS, { bubbles: true, cancelable: true }))
      },
      onFocusOutside: (event) => {
        handleContentExit(it)
        if (nav.contains(event.target)) event.preventDefault()
      },
      onPointerDownOutside: (event) => {
        const target = event.target
        const isTrigger = triggerItems().some((item) => item.trigger.contains(target))
        const isRootViewport = !!viewport?.el.contains(target)
        if (isTrigger || isRootViewport) event.preventDefault()
      },
      onEscapeKeyDown: () => {
        it.wasEscapeClose = true
      },
    })
    state.cleanups.push(() => state.layer.destroy())

    const onKeyDown = (event) => {
      if (event.defaultPrevented) return
      const isMetaKey = event.altKey || event.ctrlKey || event.metaKey
      const isTabKey = event.key === "Tab" && !isMetaKey
      if (isTabKey) {
        const candidates = getTabbableCandidates(el)
        const focusedElement = document.activeElement
        const index = candidates.findIndex((candidate) => candidate === focusedElement)
        const isMovingBackwards = event.shiftKey
        const nextCandidates = isMovingBackwards ? candidates.slice(0, index).reverse() : candidates.slice(index + 1, candidates.length)
        if (focusFirst(nextCandidates)) event.preventDefault()
        else it.focusProxy?.focus()
      }
    }
    el.addEventListener("keydown", onKeyDown)
    state.cleanups.push(() => el.removeEventListener("keydown", onKeyDown))

    if (!viewportTemplate) {
      const enter = () => onContentEnter()
      const leave = whenMouse(() => onContentLeave())
      el.addEventListener("pointerenter", enter)
      el.addEventListener("pointerleave", leave)
      state.cleanups.push(() => {
        el.removeEventListener("pointerenter", enter)
        el.removeEventListener("pointerleave", leave)
      })
    }

    state.presence = createPresence(el, {
      present: true,
      onExitComplete: () => unmountContent(it),
    })
    return state
  }

  function unmountContent(it) {
    const state = it.content
    if (!state) return
    state.presence.destroy()
    state.el.remove()
    for (const c of state.cleanups.reverse()) c()
    it.content = null
    if (viewport?.sizeContent === state.el) viewport.sizeContent = null
  }

  const setMotion = (el, motion) => {
    if (motion) el.setAttribute("data-motion", motion)
    else el.removeAttribute("data-motion")
  }

  // --- render --------------------------------------------------------------------------------

  function render() {
    const open = value !== ""
    // Viewport presence first: contents render inside it.
    if (viewportTemplate) renderViewport(open)

    for (const it of items) {
      if (!it.trigger) continue
      const isOpen = it.value === value
      it.trigger.setAttribute("data-state", isOpen ? "open" : "closed")
      it.trigger.setAttribute("aria-expanded", String(isOpen))
      if (isOpen) {
        if (!it.focusProxy) {
          const proxy = document.createElement("span")
          proxy.setAttribute("aria-hidden", "true")
          proxy.setAttribute("tabindex", "0")
          proxy.setAttribute("style", VISUALLY_HIDDEN)
          proxy.addEventListener("focus", (event) => {
            const content = it.content?.el
            const prevFocusedElement = event.relatedTarget
            const wasTriggerFocused = prevFocusedElement === it.trigger
            const wasFocusFromContent = content?.contains(prevFocusedElement)
            if (wasTriggerFocused || !wasFocusFromContent) handleContentEntry(it, wasTriggerFocused ? "start" : "end")
          })
          it.trigger.after(proxy)
          it.focusProxy = proxy
        }
        if (viewport && !it.ownsSpan && it.contentId) {
          const span = document.createElement("span")
          span.setAttribute("aria-owns", it.contentId)
          it.focusProxy.after(span)
          it.ownsSpan = span
        }
      } else {
        it.focusProxy?.remove()
        it.focusProxy = null
        it.ownsSpan?.remove()
        it.ownsSpan = null
      }
      if (!it.template) continue
      if (viewportTemplate) continue

      // Content in place (no viewport).
      if (isOpen) {
        if (!it.content) {
          const anchor = it.ownsSpan ?? it.focusProxy ?? it.trigger
          mountContent(it, anchor.parentNode, anchor.nextSibling)
          setMotion(it.content.el, motionAttribute(it))
        } else {
          const el = it.content.el
          el.setAttribute("data-state", "open")
          setMotion(el, motionAttribute(it))
          el.style.pointerEvents = ""
          if (!el.getAttribute("style")) el.removeAttribute("style")
          it.content.presence.setPresent(true)
        }
      } else if (it.content) {
        const el = it.content.el
        el.setAttribute("data-state", "closed")
        setMotion(el, motionAttribute(it))
        el.style.pointerEvents = "none"
        it.content.presence.setPresent(false)
      }
    }

    if (viewportTemplate && viewport) renderViewportContents(open)
    renderIndicator(open)
  }

  function renderViewport(open) {
    if (open) {
      if (!viewport) {
        const [el] = instantiate(viewportTemplate)
        viewportTemplate.parentNode.insertBefore(el, viewportTemplate)
        const state = { el, presence: null, cleanupSize: null, sizeContent: null, cleanups: [] }
        viewport = state
        const enter = () => onContentEnter()
        const leave = whenMouse(() => onContentLeave())
        el.addEventListener("pointerenter", enter)
        el.addEventListener("pointerleave", leave)
        state.presence = createPresence(el, {
          present: true,
          onExitComplete: () => {
            state.presence.destroy()
            for (const it of items) if (it.content && state.el.contains(it.content.el)) unmountContent(it)
            state.cleanupSize?.()
            el.remove()
            if (viewport === state) viewport = null
          },
        })
      } else {
        viewport.el.setAttribute("data-state", "open")
        viewport.el.style.pointerEvents = ""
        if (!viewport.el.getAttribute("style")) viewport.el.removeAttribute("style")
        viewport.presence.setPresent(true)
      }
    } else if (viewport) {
      viewport.el.setAttribute("data-state", "closed")
      viewport.el.style.pointerEvents = "none"
      viewport.presence.setPresent(false)
    }
  }

  function renderViewportContents(open) {
    const activeContentValue = open ? value : previousValue
    for (const it of items) {
      if (!it.template) continue
      const isActive = activeContentValue === it.value
      if (isActive) {
        if (!it.content) {
          // Keep the panels in item order.
          const later = items.slice(items.indexOf(it) + 1).find((o) => o.content && viewport.el.contains(o.content.el))
          mountContent(it, viewport.el, later?.content.el ?? null)
          setMotion(it.content.el, motionAttribute(it))
        } else {
          setMotion(it.content.el, motionAttribute(it))
          it.content.presence.setPresent(true)
        }
      } else if (it.content) {
        setMotion(it.content.el, motionAttribute(it))
        it.content.presence.setPresent(false)
      }
    }
    // The viewport takes the size of the active panel.
    const active = items.find((it) => it.value === activeContentValue)?.content?.el ?? null
    if (active && viewport.sizeContent !== active) {
      viewport.cleanupSize?.()
      viewport.sizeContent = active
      const vp = viewport
      vp.cleanupSize = observeResize(active, () => {
        if (!vp.sizeContent) return
        vp.el.style.setProperty("--radix-navigation-menu-viewport-width", vp.sizeContent.offsetWidth + "px")
        vp.el.style.setProperty("--radix-navigation-menu-viewport-height", vp.sizeContent.offsetHeight + "px")
      })
    }
  }

  // --- indicator -----------------------------------------------------------------------------

  const indicatorTemplate = [...nav.querySelectorAll("template[data-radix-navigation-menu-indicator]")].find(ownedBy) ?? null
  let indicator = null // { el, presence, activeTrigger, position, cleanups }

  function renderIndicator(open) {
    if (!indicatorTemplate || !indicatorTrack) return
    const isHorizontal = orientation === "horizontal"
    if (open) {
      if (!indicator) {
        const state = { el: null, presence: null, activeTrigger: null, position: null, cleanups: [] }
        indicator = state
        const handlePositionChange = () => {
          if (!state.activeTrigger || indicator !== state) return
          state.position = {
            size: isHorizontal ? state.activeTrigger.offsetWidth : state.activeTrigger.offsetHeight,
            offset: isHorizontal ? state.activeTrigger.offsetLeft : state.activeTrigger.offsetTop,
          }
          applyIndicator(state)
        }
        state.handlePositionChange = handlePositionChange
        state.cleanups.push(observeResize(indicatorTrack, handlePositionChange))
      } else if (indicator.el) {
        indicator.el.setAttribute("data-state", "visible")
        indicator.presence.setPresent(true)
      }
      const triggerNode = triggerItems().find((it) => it.value === value)?.trigger
      if (triggerNode && triggerNode !== indicator.activeTrigger) {
        indicator.activeTrigger = triggerNode
        indicator.cleanupTrigger?.()
        indicator.cleanupTrigger = observeResize(triggerNode, indicator.handlePositionChange)
      }
    } else if (indicator) {
      if (indicator.el) {
        indicator.el.setAttribute("data-state", "hidden")
        indicator.presence.setPresent(false)
      } else {
        destroyIndicator(indicator)
      }
    }
  }

  function applyIndicator(state) {
    if (!state.position) return
    const isHorizontal = orientation === "horizontal"
    if (!state.el) {
      const [el] = instantiate(indicatorTemplate)
      state.el = el
      indicatorTrack.appendChild(el)
      state.presence = createPresence(el, { present: true, onExitComplete: () => destroyIndicator(state) })
    }
    const s = state.el.style
    s.position = "absolute"
    if (isHorizontal) {
      s.left = "0px"
      s.width = state.position.size + "px"
      s.transform = `translateX(${state.position.offset}px)`
    } else {
      s.top = "0px"
      s.height = state.position.size + "px"
      s.transform = `translateY(${state.position.offset}px)`
    }
  }

  function destroyIndicator(state) {
    state.presence?.destroy()
    state.el?.remove()
    state.cleanupTrigger?.()
    for (const c of state.cleanups) c()
    if (indicator === state) indicator = null
  }

  // --- triggers --------------------------------------------------------------------------------

  for (const it of items) {
    const t = it.trigger
    if (!t) continue
    const disabled = () => t.disabled
    const onPointerEnter = () => {
      it.wasClickClose = false
      it.wasEscapeClose = false
    }
    const onPointerMove = whenMouse(() => {
      if (disabled() || it.wasClickClose || it.wasEscapeClose || it.hasPointerMoveOpened) return
      onTriggerEnter(it.value)
      it.hasPointerMoveOpened = true
    })
    const onPointerLeave = whenMouse(() => {
      if (disabled()) return
      onTriggerLeave()
      it.hasPointerMoveOpened = false
    })
    const onClick = (event) => {
      if (event.defaultPrevented) return
      const open = it.value === value
      onItemSelect(it.value)
      it.wasClickClose = open
    }
    const onKeyDown = (event) => {
      if (event.defaultPrevented) return
      const verticalEntryKey = dir === "rtl" ? "ArrowLeft" : "ArrowRight"
      const entryKey = { horizontal: "ArrowDown", vertical: verticalEntryKey }[orientation]
      if (it.value === value && event.key === entryKey) {
        handleContentEntry(it)
        event.preventDefault()
      }
    }
    t.addEventListener("pointerenter", onPointerEnter)
    t.addEventListener("pointermove", onPointerMove)
    t.addEventListener("pointerleave", onPointerLeave)
    t.addEventListener("click", onClick)
    t.addEventListener("keydown", onKeyDown)
    cleanups.push(() => {
      t.removeEventListener("pointerenter", onPointerEnter)
      t.removeEventListener("pointermove", onPointerMove)
      t.removeEventListener("pointerleave", onPointerLeave)
      t.removeEventListener("click", onClick)
      t.removeEventListener("keydown", onKeyDown)
    })
  }

  // --- focus groups (the list and each panel) --------------------------------------------------

  const groupOf = (el) => {
    for (let n = el.parentElement; n && n !== nav.parentElement; n = n.parentElement) {
      if (n === list) return n
      if (items.some((it) => it.content?.el === n)) return n
    }
    return null
  }
  const onGroupKeyDown = (event) => {
    if (event.defaultPrevented) return
    if (!["Home", "End", ...ARROW_KEYS].includes(event.key)) return
    const item = event.target instanceof Element ? event.target.closest("[data-radix-collection-item]") : null
    if (!item || !nav.contains(item)) return
    const group = groupOf(item)
    if (!group) return
    let candidateNodes = [...group.querySelectorAll("[data-radix-collection-item]")].filter((el) => groupOf(el) === group)
    const prevItemKey = dir === "rtl" ? "ArrowRight" : "ArrowLeft"
    const prevKeys = [prevItemKey, "ArrowUp", "End"]
    if (prevKeys.includes(event.key)) candidateNodes.reverse()
    if (ARROW_KEYS.includes(event.key)) {
      const currentIndex = candidateNodes.indexOf(item)
      candidateNodes = candidateNodes.slice(currentIndex + 1)
    }
    setTimeout(() => focusFirst(candidateNodes))
    event.preventDefault()
  }
  nav.addEventListener("keydown", onGroupKeyDown)
  cleanups.push(() => nav.removeEventListener("keydown", onGroupKeyDown))

  // NavigationMenuLink: select event, then dismiss the panel holding it.
  const onLinkClick = (event) => {
    const link = event.target instanceof Element ? event.target.closest('[data-slot="navigation-menu-link"]') : null
    if (!link || !nav.contains(link)) return
    const target = event.target
    const linkSelectEvent = new CustomEvent(LINK_SELECT, { bubbles: true, cancelable: true })
    target.dispatchEvent(linkSelectEvent)
    if (!linkSelectEvent.defaultPrevented && !event.metaKey) {
      target.dispatchEvent(new CustomEvent(ROOT_CONTENT_DISMISS, { bubbles: true, cancelable: true }))
    }
  }
  nav.addEventListener("click", onLinkClick)
  cleanups.push(() => nav.removeEventListener("click", onLinkClick))

  if (value !== "") render()

  const controller = {
    get value() {
      return value
    },
    setValue,
    dispose() {
      window.clearTimeout(openTimer)
      window.clearTimeout(closeTimer)
      window.clearTimeout(skipDelayTimer)
      for (const it of items) if (it.content) unmountContent(it)
      if (indicator) destroyIndicator(indicator)
      for (const c of cleanups) c()
      menus.delete(nav)
    },
  }
  menus.set(nav, controller)
  return controller
}

define("navigation-menu", {
  selector: "nav[data-radix-navigation-menu]",
  mount(nav) {
    const m = createNavigationMenu(nav)
    return () => m.dispose()
  },
})
