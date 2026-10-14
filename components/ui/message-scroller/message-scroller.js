// message-scroller.js: behavior of @shadcn/react/message-scroller (the
// headless MessageScroller under shadcn's message-scroller) for the
// server-rendered markup of components/ui/message-scroller.
//
// It owns the transcript viewport's scroll position: the opening position
// (defaultScrollPosition), following the end while autoScroll is on,
// anchoring a newly appended scrollAnchor item near the top (with the tail
// spacer), keeping the reading position when rows are prepended, the
// data-scrollable / data-autoscrolling / data-pending-scroll attributes and
// the scroll buttons' data-active / inert / tabindex. The logic is a line by
// line port of use-message-scroller-controller.ts, -commands.ts and
// geometry.ts; React's re-renders become direct DOM writes.
//
// Markup contract (rendered by the Go package, markers removed on mount):
//   - every part carries data-message-scroller="<provider id>" and
//     data-message-scroller-part="root" | "viewport" | "content" | "button";
//   - the viewport carries data-message-scroller-options='{"autoScroll":..,
//     "defaultScrollPosition":..,"scrollEdgeThreshold":..,
//     "scrollPreviousItemPeek":..,"scrollMargin":..,
//     "preserveScrollOnPrepend":..}';
//   - buttons carry data-message-scroller-behavior when not "smooth" and
//     data-direction.
// Rows are the content's element children except the tail spacer
// ([data-message-scroller-spacer]); rows with data-message-id can be
// targeted by scrollToMessage.
//
// Scripts can drive and observe a scroller:
//   import { getMessageScroller } from "/_shadcn/message-scroller.js"
//   const s = getMessageScroller(viewportOrRootOrContent)
//   s.scrollToEnd({ behavior: "smooth" }); s.scrollToStart(); s.scrollToMessage("m3", { align: "center" })
//   s.scrollable                 // { start, end }
//   const off = s.subscribeVisibility((v) => v.currentAnchorId, v.visibleMessageIds)
//   viewport.addEventListener("shadcn:scrollable-change", (e) => e.detail) // { start, end }

import { define, emit } from "./core.js"

const DEFAULT_SCROLL_EDGE_THRESHOLD = 8
const DEFAULT_SCROLL_PREVIOUS_ITEM_PEEK = 64
const DEFAULT_SCROLL_MARGIN = 0
const SCROLL_POSITION_EPSILON = 0.5
const AUTOSCROLLING_CLEAR_DELAY = 180
const USER_SCROLL_KEYS = new Set(["ArrowDown", "ArrowUp", "End", "Home", "PageDown", "PageUp", " "])
const EMPTY_SCROLLABLE = { start: false, end: false }
const EMPTY_VISIBILITY = { currentAnchorId: null, visibleMessageIds: [] }

const scrollers = new WeakMap() // part element -> controller

/** Returns the controller of the scroller a part (root, viewport, content, item, button) belongs to. */
export function getMessageScroller(el) {
  for (let node = el; node; node = node.parentElement) {
    const c = scrollers.get(node)
    if (c) return c
  }
  return undefined
}

// --- geometry.ts -------------------------------------------------------------------------------

function getItems(content, spacer) {
  return Array.from(content.children).filter((child) => child instanceof HTMLElement && child !== spacer)
}

function readCssPixel(value) {
  if (!value) return 0
  const number = Number.parseFloat(value)
  return Number.isFinite(number) ? number : 0
}

function getBlockPadding(element) {
  const style = window.getComputedStyle(element)
  return {
    end: readCssPixel(style.paddingBlockEnd || style.paddingBottom),
    start: readCssPixel(style.paddingBlockStart || style.paddingTop),
  }
}

function getContentBlockPadding(spacer) {
  const content = spacer?.parentElement
  if (!content) return { end: 0, start: 0 }
  return getBlockPadding(content)
}

function getFlexGap(element) {
  if (!element) return 0
  const style = window.getComputedStyle(element)
  const gap = style.rowGap === "normal" ? style.gap : style.rowGap
  return readCssPixel(gap)
}

function getContentBottom({ content, spacer, viewport }) {
  const items = getItems(content, spacer)
  const padding = getBlockPadding(content)
  const viewportRect = viewport.getBoundingClientRect()
  const scrollTop = viewport.scrollTop
  let contentBottom = padding.start + padding.end
  for (const item of items) {
    const rect = item.getBoundingClientRect()
    contentBottom = Math.max(contentBottom, rect.bottom - viewportRect.top + scrollTop + padding.end)
  }
  return contentBottom
}

function getScrollable({ content, scrollEdgeThreshold, spacer, viewport }) {
  if (!viewport || !content) return EMPTY_SCROLLABLE
  const contentBottom = getContentBottom({ content, spacer, viewport })
  return {
    start: viewport.scrollTop > scrollEdgeThreshold,
    end: contentBottom - viewport.scrollTop - viewport.clientHeight > scrollEdgeThreshold,
  }
}

function getVisibilityState({ content, scrollMargin, scrollPreviousItemPeek, spacer, viewport, visibleMessageIds }) {
  if (!content || !viewport) return EMPTY_VISIBILITY
  const viewportRect = viewport.getBoundingClientRect()
  const lineTop = viewportRect.top + scrollMargin + scrollPreviousItemPeek
  const trackByLayout = typeof IntersectionObserver === "undefined"
  const visible = []
  let currentAnchorId = null
  for (const item of getItems(content, spacer)) {
    const messageId = item.dataset.messageId
    if (!messageId) continue
    const isAnchor = item.dataset.scrollAnchor === "true"
    const rect = isAnchor || trackByLayout ? item.getBoundingClientRect() : null
    const isVisible = trackByLayout && rect ? rect.bottom > lineTop && rect.top < viewportRect.bottom : visibleMessageIds.has(messageId)
    if (isVisible) visible.push(messageId)
    if (isAnchor && rect && rect.top <= lineTop + SCROLL_POSITION_EPSILON) currentAnchorId = messageId
  }
  if (visible.length === 0 && currentAnchorId === null) return EMPTY_VISIBILITY
  return { currentAnchorId, visibleMessageIds: visible }
}

function getNewScrollAnchor(items, previousItemCount) {
  for (let i = previousItemCount; i < items.length; i++) {
    if (items[i]?.dataset.scrollAnchor === "true") return items[i]
  }
  return null
}

function getUnanchoredScrollAnchor(items, handled) {
  for (const item of items) {
    if (item.dataset.scrollAnchor === "true" && !handled.has(item)) return item
  }
  return null
}

function hasMultipleNewScrollAnchors(items, previousItemCount) {
  let count = 0
  for (let i = previousItemCount; i < items.length; i++) {
    if (items[i]?.dataset.scrollAnchor !== "true") continue
    count += 1
    if (count > 1) return true
  }
  return false
}

function getLastScrollAnchor(items) {
  for (let i = items.length - 1; i >= 0; i--) {
    if (items[i]?.dataset.scrollAnchor === "true") return items[i]
  }
  return null
}

function getFirstVisibleMessageItem({ content, spacer, viewport }) {
  const viewportRect = viewport.getBoundingClientRect()
  for (const item of getItems(content, spacer)) {
    if (!item.dataset.messageId) continue
    const rect = item.getBoundingClientRect()
    if (rect.bottom > viewportRect.top && rect.top < viewportRect.bottom) return item
  }
  return null
}

function getElementTop(element, viewport) {
  return element.getBoundingClientRect().top - viewport.getBoundingClientRect().top + viewport.scrollTop
}

function getElementViewportTop(element, viewport) {
  return element.getBoundingClientRect().top - viewport.getBoundingClientRect().top
}

function getElementScrollTop({ align, element, scrollMargin, spacer, viewport }) {
  const elementTop = getElementTop(element, viewport)
  const elementHeight = element.getBoundingClientRect().height
  const contentPadding = getContentBlockPadding(spacer)
  if (align === "center") {
    const insetHeight = Math.max(0, viewport.clientHeight - contentPadding.start - contentPadding.end)
    return elementTop - contentPadding.start - (insetHeight - elementHeight) / 2 - scrollMargin
  }
  if (align === "end") {
    return elementTop - viewport.clientHeight + elementHeight + contentPadding.end + scrollMargin
  }
  if (align === "nearest") {
    const elementBottom = elementTop + elementHeight
    const viewportTop = viewport.scrollTop + contentPadding.start
    const viewportBottom = viewport.scrollTop + viewport.clientHeight - contentPadding.end
    if (elementTop >= viewportTop && elementBottom <= viewportBottom) return viewport.scrollTop
    if (elementTop < viewportTop) return elementTop - contentPadding.start - scrollMargin
    return elementBottom - viewport.clientHeight + contentPadding.end + scrollMargin
  }
  return elementTop - contentPadding.start - scrollMargin
}

function getTailSpacerHeight({ content, scrollTop, spacer, viewport }) {
  return scrollTop + viewport.clientHeight - getContentBottom({ content, spacer, viewport })
}

const getMaxScrollTop = (viewport) => Math.max(0, viewport.scrollHeight - viewport.clientHeight)

// unrenderRows makes the rows with content-visibility: auto measure like
// rows that were never rendered (their contain-intrinsic-size without the
// remembered size) and returns a function restoring them.
function unrenderRows(content, spacer) {
  if (!content) return () => {}
  const saved = []
  for (const row of getItems(content, spacer)) {
    const cs = getComputedStyle(row)
    if (cs.contentVisibility !== "auto") continue
    const style = row.getAttribute("style")
    saved.push([row, style])
    const strip = (v) => v.replace(/\bauto\s+/g, "")
    // Attribute writes, not CSSOM ones: Chromium can leave a lazily
    // serialized style="" behind a CSSOM write followed by removeAttribute.
    row.setAttribute(
      "style",
      `${style ? style + ";" : ""}content-visibility:hidden;contain-intrinsic-width:${strip(cs.containIntrinsicWidth)};contain-intrinsic-height:${strip(cs.containIntrinsicHeight)}`
    )
  }
  return () => {
    for (const [row, style] of saved) {
      if (style === null) row.removeAttribute("style")
      else row.setAttribute("style", style)
    }
  }
}

// --- controller --------------------------------------------------------------------------------

function createMessageScroller(viewport) {
  const id = viewport.getAttribute("data-message-scroller")
  const parts = [...document.querySelectorAll(`[data-message-scroller="${CSS.escape(id)}"]`)]
  const part = (name) => parts.filter((p) => p.getAttribute("data-message-scroller-part") === name)
  let options = {}
  try {
    options = JSON.parse(viewport.getAttribute("data-message-scroller-options") || "{}")
  } catch {}
  const root = part("root")[0] ?? null
  const content = part("content")[0] ?? null
  const buttons = part("button")
  for (const p of parts) {
    p.removeAttribute("data-message-scroller")
    p.removeAttribute("data-message-scroller-part")
  }
  viewport.removeAttribute("data-message-scroller-options")
  const spacer = content?.querySelector(":scope > [data-message-scroller-spacer]") ?? null

  // Props (MessageScrollerProvider defaults).
  const autoScroll = !!options.autoScroll
  const defaultScrollPosition = options.defaultScrollPosition ?? "end"
  const scrollEdgeThreshold = options.scrollEdgeThreshold ?? DEFAULT_SCROLL_EDGE_THRESHOLD
  const scrollPreviousItemPeek = options.scrollPreviousItemPeek ?? DEFAULT_SCROLL_PREVIOUS_ITEM_PEEK
  const scrollMargin = options.scrollMargin ?? DEFAULT_SCROLL_MARGIN
  const preserveScrollOnPrepend = options.preserveScrollOnPrepend ?? true

  // Refs.
  let autoscrolling = false
  let autoscrollingTimeout = null
  let streamingTurn = null
  let defaultScrollPositionApplied = false
  let itemCount = 0
  let firstItem = null
  let lastScrollTop = 0
  let mode = autoScroll ? "following-bottom" : "free-scrolling"
  let pendingScrollToMessage = null
  let prependRestore = null
  let pendingScrollFrame = null
  let spacerGap = getFlexGap(spacer?.parentElement ?? null)
  let spacerHeight = 0
  let stateFrame = null
  let visibilityFrame = null
  let visibilityObserver = null
  const visibleMessageIds = new Set()
  const handledScrollAnchors = new WeakSet()
  const messageElements = new Map()

  // Stores.
  let pendingDefaultScroll = defaultScrollPosition === "end" || defaultScrollPosition === "last-anchor"
  let scrollable = EMPTY_SCROLLABLE
  let visibility = EMPTY_VISIBILITY
  const visibilityListeners = new Set()

  function setPendingDefaultScroll(next) {
    if (pendingDefaultScroll === next) return
    pendingDefaultScroll = next
    for (const el of [root, viewport]) el?.toggleAttribute("data-pending-scroll", next)
  }
  const clearPendingDefaultScroll = () => setPendingDefaultScroll(false)
  function markDefaultScrollPositionApplied() {
    defaultScrollPositionApplied = true
    clearPendingDefaultScroll()
  }

  // MessageScrollerButton re-renders: data-active, inert, tabIndex.
  function renderButtons() {
    for (const b of buttons) {
      const active = b.getAttribute("data-direction") === "start" ? scrollable.start : scrollable.end
      b.setAttribute("data-active", active ? "true" : "false")
      b.toggleAttribute("inert", !active)
      if (active) {
        if (b.hasAttribute("data-message-scroller-tabindex")) b.setAttribute("tabindex", b.getAttribute("data-message-scroller-tabindex"))
        else b.removeAttribute("tabindex")
      } else b.setAttribute("tabindex", "-1")
    }
  }

  function setScrollable(next) {
    if (scrollable.start === next.start && scrollable.end === next.end) return
    scrollable = next
    renderButtons()
    emit(viewport, "shadcn:scrollable-change", { ...next })
  }

  function setVisibility(next) {
    const same =
      visibility.currentAnchorId === next.currentAnchorId &&
      visibility.visibleMessageIds.length === next.visibleMessageIds.length &&
      visibility.visibleMessageIds.every((m, i) => m === next.visibleMessageIds[i])
    if (same) return
    visibility = next
    visibilityListeners.forEach((l) => l(visibility))
  }

  function writeStateAttributes(state) {
    const value = [state.start && "start", state.end && "end"].filter(Boolean).join(" ")
    for (const element of [root, viewport]) {
      if (!element) continue
      if (value) element.setAttribute("data-scrollable", value)
      else element.removeAttribute("data-scrollable")
      element.toggleAttribute("data-autoscrolling", autoscrolling)
    }
  }

  function reconcileFollowMode(next) {
    const scrollTop = viewport.scrollTop ?? 0
    const scrolledUp = scrollTop < lastScrollTop - SCROLL_POSITION_EPSILON
    lastScrollTop = scrollTop
    if (autoScroll && !next.end && mode !== "settling-jump" && mode !== "anchored-to-message") {
      mode = "following-bottom"
    } else if (mode === "following-bottom" && next.end && scrolledUp && !autoscrolling) {
      mode = "free-scrolling"
    }
  }

  function commitScrollState() {
    const next = getScrollable({ content, scrollEdgeThreshold, spacer, viewport })
    reconcileFollowMode(next)
    const published = mode === "following-bottom" ? { ...next, end: false } : next
    writeStateAttributes(published)
    setScrollable(published)
  }

  function scheduleStateCommit() {
    if (stateFrame !== null) return
    stateFrame = requestAnimationFrame(() => {
      stateFrame = null
      commitScrollState()
    })
  }

  function scheduleVisibilitySync() {
    if (visibilityListeners.size === 0) return
    if (visibilityFrame !== null) return
    visibilityFrame = requestAnimationFrame(() => {
      visibilityFrame = null
      if (visibilityListeners.size === 0) return
      setVisibility(getVisibilityState({ content, scrollMargin, scrollPreviousItemPeek, spacer, viewport, visibleMessageIds }))
    })
  }

  // --- use-message-scroller-commands.ts ---

  function setAutoScrolling(next) {
    if (autoscrollingTimeout !== null) {
      clearTimeout(autoscrollingTimeout)
      autoscrollingTimeout = null
    }
    if (autoscrolling !== next) {
      autoscrolling = next
      commitScrollState()
    }
    if (next) {
      autoscrollingTimeout = setTimeout(() => {
        autoscrollingTimeout = null
        autoscrolling = false
        commitScrollState()
      }, AUTOSCROLLING_CLEAR_DELAY)
    }
  }

  function setTailSpacerHeight(height) {
    if (!spacer) return
    const nextHeight = Math.max(0, Math.ceil(height))
    if (spacerHeight === nextHeight) return
    spacerHeight = nextHeight
    spacer.hidden = nextHeight === 0
    spacer.style.height = `${nextHeight}px`
    spacer.style.marginTop = nextHeight > 0 ? `${-spacerGap}px` : ""
    if (!spacer.getAttribute("style")) spacer.removeAttribute("style")
  }

  function scrollToPosition(scrollTop, { behavior = "auto", autoscrolling: auto = false } = {}) {
    const nextScrollTop = Math.max(0, scrollTop)
    if (Math.abs(viewport.scrollTop - nextScrollTop) <= SCROLL_POSITION_EPSILON) {
      viewport.scrollTop = nextScrollTop
      commitScrollState()
      return
    }
    if (auto) setAutoScrolling(true)
    viewport.scrollTo({ top: nextScrollTop, behavior })
    scheduleStateCommit()
  }

  function scrollToStart({ behavior = "auto" } = {}) {
    setTailSpacerHeight(0)
    streamingTurn = null
    mode = "free-scrolling"
    scrollToPosition(0, { behavior })
    scheduleVisibilitySync()
    return true
  }

  function scrollToEnd({ behavior = "auto" } = {}) {
    setTailSpacerHeight(0)
    streamingTurn = null
    mode = autoScroll ? "following-bottom" : "free-scrolling"
    scrollToPosition(getMaxScrollTop(viewport), { autoscrolling: true, behavior })
    scheduleVisibilitySync()
    return true
  }

  function scrollToElement(element, { align = "start", behavior = "auto", scrollMargin: margin = scrollMargin } = {}, { keepPreviousPeek = false } = {}) {
    if (!content || !content.contains(element)) return false
    const scrollTop = getElementScrollTop({
      align,
      element,
      scrollMargin: keepPreviousPeek ? margin + scrollPreviousItemPeek : margin,
      spacer,
      viewport,
    })
    setTailSpacerHeight(getTailSpacerHeight({ content, scrollTop, spacer, viewport }))
    prependRestore = { element, viewportTop: getElementViewportTop(element, viewport) }
    mode = keepPreviousPeek ? "anchored-to-message" : "settling-jump"
    streamingTurn = keepPreviousPeek ? element : null
    scrollToPosition(scrollTop, { behavior })
    scheduleVisibilitySync()
    return true
  }

  function reanchorToAnchoredMessage() {
    const element = streamingTurn
    if (!element || !element.isConnected || mode !== "anchored-to-message") return false
    return scrollToElement(element, { align: "start" }, { keepPreviousPeek: true })
  }

  function scrollToMessage(messageId, opts) {
    const element = messageElements.get(messageId)
    if (!element) {
      if (itemCount === 0) {
        pendingScrollToMessage = { messageId, options: opts }
        markDefaultScrollPositionApplied()
        return true
      }
      return false
    }
    markDefaultScrollPositionApplied()
    if (scrollToElement(element, opts)) {
      pendingScrollToMessage = null
      return true
    }
    pendingScrollToMessage = { messageId, options: opts }
    return true
  }

  function flushPendingScrollToMessage() {
    const pending = pendingScrollToMessage
    if (!pending) return false
    const element = messageElements.get(pending.messageId)
    if (!element) return false
    if (!scrollToElement(element, pending.options)) return false
    pendingScrollToMessage = null
    markDefaultScrollPositionApplied()
    return true
  }

  // --- use-message-scroller-controller.ts ---

  function restorePrependedAnchor() {
    const anchor = prependRestore
    if (!anchor || !anchor.element.isConnected) return false
    const delta = getElementViewportTop(anchor.element, viewport) - anchor.viewportTop
    if (Math.abs(delta) <= SCROLL_POSITION_EPSILON) return false
    viewport.scrollTop += delta
    anchor.viewportTop = getElementViewportTop(anchor.element, viewport)
    scheduleStateCommit()
    scheduleVisibilitySync()
    return true
  }

  function capturePrependAnchor() {
    if (!content) {
      prependRestore = null
      return
    }
    const anchor = getFirstVisibleMessageItem({ content, spacer, viewport })
    prependRestore = anchor ? { element: anchor, viewportTop: getElementViewportTop(anchor, viewport) } : null
  }

  function schedulePendingScrollToMessageFlush() {
    if (pendingScrollFrame !== null) return
    pendingScrollFrame = requestAnimationFrame(() => {
      pendingScrollFrame = null
      if (flushPendingScrollToMessage()) capturePrependAnchor()
    })
  }

  function applyDefaultScrollPosition() {
    if (!defaultScrollPosition || defaultScrollPositionApplied || itemCount === 0) return false
    let handled = false
    if (defaultScrollPosition === "last-anchor") {
      const anchor = content ? getLastScrollAnchor(getItems(content, spacer)) : null
      if (!content || !anchor) handled = scrollToEnd({ behavior: "auto" })
      else {
        const anchorTop = getElementTop(anchor, viewport)
        const contentBottom = getContentBottom({ content, spacer, viewport })
        const lastTurnFits = contentBottom - anchorTop <= viewport.clientHeight
        handled = lastTurnFits ? scrollToEnd({ behavior: "auto" }) : scrollToElement(anchor, { align: "start" }, { keepPreviousPeek: true })
      }
    } else {
      handled = defaultScrollPosition === "end" ? scrollToEnd({ behavior: "auto" }) : scrollToStart({ behavior: "auto" })
    }
    if (!handled) return false
    markDefaultScrollPositionApplied()
    return true
  }

  // MessageScrollerItem refs: registerMessage for every row with a messageId.
  function registerMessage(messageId, element, removedElement) {
    if (element) {
      messageElements.set(messageId, element)
      visibilityObserver?.observe(element)
      scheduleVisibilitySync()
      if (pendingScrollToMessage?.messageId === messageId) schedulePendingScrollToMessageFlush()
      return
    }
    if (removedElement && messageElements.get(messageId) === removedElement) {
      messageElements.delete(messageId)
      visibleMessageIds.delete(messageId)
      visibilityObserver?.unobserve(removedElement)
      scheduleVisibilitySync()
    }
  }

  let registered = new Map() // element -> messageId
  function syncRegistrations() {
    if (!content) return
    const next = new Map()
    for (const el of getItems(content, spacer)) {
      const mid = el.getAttribute("data-message-id")
      if (mid) next.set(el, mid)
    }
    for (const [el, mid] of registered) if (next.get(el) !== mid) registerMessage(mid, null, el)
    for (const [el, mid] of next) if (registered.get(el) !== mid) registerMessage(mid, el)
    registered = next
  }

  function handleContentChange() {
    if (!content) return
    const items = getItems(content, spacer)
    const previousItemCount = itemCount
    const previousFirstItem = firstItem
    itemCount = items.length
    firstItem = items[0] ?? null

    const reconcile = () => {
      if (flushPendingScrollToMessage()) return
      if (previousItemCount === 0) {
        if (applyDefaultScrollPosition()) return
        if (items.length > 0 && autoScroll && scrollToEnd({ behavior: "auto" })) return
        commitScrollState()
        scheduleVisibilitySync()
        return
      }
      const previousFirstItemIndex = previousFirstItem ? items.indexOf(previousFirstItem) : -1
      if (preserveScrollOnPrepend && previousFirstItemIndex > 0) {
        restorePrependedAnchor()
        return
      }
      if (items.length > previousItemCount) {
        const anchor = getNewScrollAnchor(items, previousItemCount)
        if (anchor) {
          if (autoScroll && mode === "following-bottom" && hasMultipleNewScrollAnchors(items, previousItemCount)) {
            scrollToEnd({ behavior: "auto" })
            return
          }
          scrollToElement(anchor, { align: "start" }, { keepPreviousPeek: true })
          handledScrollAnchors.add(anchor)
          return
        }
      }
      if (items.length === previousItemCount) {
        const anchor = getUnanchoredScrollAnchor(items, handledScrollAnchors)
        if (anchor) {
          scrollToElement(anchor, { align: "start" }, { keepPreviousPeek: true })
          handledScrollAnchors.add(anchor)
          return
        }
      }
      if (mode === "following-bottom" && autoScroll) scrollToEnd({ behavior: "auto" })
      else {
        commitScrollState()
        scheduleVisibilitySync()
      }
    }
    reconcile()
    capturePrependAnchor()
  }

  function handleResize() {
    if (mode === "following-bottom" && autoScroll) {
      scrollToEnd({ behavior: "auto" })
      return
    }
    const previousSpacerHeight = spacerHeight
    if (reanchorToAnchoredMessage()) {
      if (autoScroll && previousSpacerHeight > 0 && spacerHeight === 0) scrollToEnd({ behavior: "auto" })
      return
    }
    scheduleStateCommit()
    scheduleVisibilitySync()
  }

  function observeVisibility() {
    if (visibilityListeners.size === 0) return
    if (typeof IntersectionObserver === "undefined") {
      scheduleVisibilitySync()
      return
    }
    if (!visibilityObserver) {
      visibilityObserver = new IntersectionObserver(
        (entries) => {
          for (const entry of entries) {
            const messageId = entry.target.dataset.messageId
            if (!messageId) continue
            if (entry.isIntersecting) visibleMessageIds.add(messageId)
            else visibleMessageIds.delete(messageId)
          }
          scheduleVisibilitySync()
        },
        { root: viewport, rootMargin: `${-(scrollMargin + scrollPreviousItemPeek)}px 0px 0px 0px`, threshold: [0, 0.01, 0.5, 1] }
      )
    }
    messageElements.forEach((element) => visibilityObserver?.observe(element))
    scheduleVisibilitySync()
  }

  function unobserveVisibility() {
    if (visibilityFrame !== null) {
      cancelAnimationFrame(visibilityFrame)
      visibilityFrame = null
    }
    visibilityObserver?.disconnect()
    visibilityObserver = null
    visibleMessageIds.clear()
    setVisibility(EMPTY_VISIBILITY)
  }

  function userScrollIntent() {
    if (mode === "following-bottom" || mode === "anchored-to-message" || mode === "settling-jump") {
      streamingTurn = null
      mode = "free-scrolling"
    }
  }

  function syncAfterScroll() {
    commitScrollState()
    scheduleVisibilitySync()
    capturePrependAnchor()
  }

  // --- mount, in React's commit order ---

  const cleanups = []
  // Ref callbacks: setRootElement / setViewportElement mirror the (empty)
  // state attributes; items register.
  writeStateAttributes(scrollable)
  syncRegistrations()

  // Viewport handlers.
  const onScroll = () => syncAfterScroll()
  const onIntent = () => userScrollIntent()
  const onKeyDown = (e) => {
    if (USER_SCROLL_KEYS.has(e.key)) userScrollIntent()
  }
  viewport.addEventListener("scroll", onScroll)
  viewport.addEventListener("wheel", onIntent, { passive: true })
  viewport.addEventListener("touchmove", onIntent, { passive: true })
  viewport.addEventListener("keydown", onKeyDown)
  cleanups.push(() => {
    viewport.removeEventListener("scroll", onScroll)
    viewport.removeEventListener("wheel", onIntent)
    viewport.removeEventListener("touchmove", onIntent)
    viewport.removeEventListener("keydown", onKeyDown)
  })

  // Buttons.
  for (const b of buttons) {
    const behavior = b.getAttribute("data-message-scroller-behavior") || "smooth"
    b.removeAttribute("data-message-scroller-behavior")
    const onClick = (event) => {
      const active = b.getAttribute("data-active") === "true"
      if (!active || event.defaultPrevented) return
      b.blur()
      if (b.getAttribute("data-direction") === "start") scrollToStart({ behavior })
      else scrollToEnd({ behavior })
    }
    b.addEventListener("click", onClick)
    cleanups.push(() => b.removeEventListener("click", onClick))
  }

  // React runs the layout effects below on rows that were never rendered:
  // rows with content-visibility: auto (MessageScrollerItem) are still
  // skipped and measure as their contain-intrinsic-size. Rows the page has
  // already painted remember their real size, so they are measured as
  // never rendered while the effects run.
  const restoreRows = unrenderRows(content, spacer)

  // MessageScrollerContent's layout effect: handleContentChange, then a
  // MutationObserver on its rows.
  if (content) {
    handleContentChange()
    const observer = new MutationObserver(() => {
      syncRegistrations()
      handleContentChange()
    })
    observer.observe(content, { childList: true })
    cleanups.push(() => observer.disconnect())
  }

  // The provider's layout effects.
  if (!applyDefaultScrollPosition() && itemCount === 0) clearPendingDefaultScroll()
  if (autoScroll && mode === "following-bottom" && itemCount > 0) scrollToEnd({ behavior: "auto" })
  else commitScrollState()

  restoreRows()

  // Passive effects: resize observers on the viewport and the content.
  if (typeof ResizeObserver !== "undefined") {
    for (const el of [viewport, content]) {
      if (!el) continue
      let frame = 0
      const observer = new ResizeObserver(() => {
        cancelAnimationFrame(frame)
        frame = requestAnimationFrame(handleResize)
      })
      observer.observe(el)
      cleanups.push(() => {
        cancelAnimationFrame(frame)
        observer.disconnect()
      })
    }
  }

  cleanups.push(() => {
    if (stateFrame !== null) cancelAnimationFrame(stateFrame)
    if (visibilityFrame !== null) cancelAnimationFrame(visibilityFrame)
    if (autoscrollingTimeout !== null) clearTimeout(autoscrollingTimeout)
    if (pendingScrollFrame !== null) cancelAnimationFrame(pendingScrollFrame)
    visibilityObserver?.disconnect()
  })

  const controller = {
    scrollToEnd,
    scrollToStart,
    scrollToMessage,
    get scrollable() {
      return scrollable
    },
    get visibility() {
      return visibility
    },
    /** Subscribes to { currentAnchorId, visibleMessageIds } (tracking is lazy, like the hook). */
    subscribeVisibility(listener) {
      const wasEmpty = visibilityListeners.size === 0
      visibilityListeners.add(listener)
      if (wasEmpty) observeVisibility()
      return () => {
        visibilityListeners.delete(listener)
        if (visibilityListeners.size === 0) unobserveVisibility()
      }
    },
  }
  for (const el of [root, viewport, content, ...buttons]) if (el) scrollers.set(el, controller)
  return () => {
    cleanups.forEach((fn) => fn())
    for (const el of [root, viewport, content, ...buttons]) if (el) scrollers.delete(el)
  }
}

define("message-scroller", {
  selector: '[data-message-scroller-part="viewport"]',
  mount: createMessageScroller,
})
