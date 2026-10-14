// scroll-area.js: Radix ScrollArea (@radix-ui/react-scroll-area 1.2.10).
//
// The Go side renders the root, the <style> (when it has the CSP nonce),
// the viewport with its content wrapper, and every ScrollBar as a
// <template> (at its place in the tree: the vertical one after the
// viewport, a horizontal one wherever the caller put it). This module
// mounts the scrollbars exactly when Radix renders them, per type:
//
//   hover   while the pointer is over the root (hidden scrollHideDelay ms
//           after it leaves), and only while the content overflows (auto)
//   scroll  while the viewport scrolls in the scrollbar's direction (and
//           while the pointer is over the scrollbar), hidden
//           scrollHideDelay ms after scrolling stopped
//   auto    while the content overflows
//   always  always
//
// Each scrollbar measures the viewport, content and itself (ResizeObserver,
// next animation frame, 10ms debounce, like Radix) to size its thumb
// (--radix-scroll-area-thumb-width/height) and positions the thumb from the
// scroll position (translate3d). Dragging the thumb or pressing the track
// scrolls; the wheel over a scrollbar scrolls the viewport. When both
// scrollbars are mounted (types other than scroll), a corner fills the gap
// and its size is published as --radix-scroll-area-corner-width/height.
//
// Radix props with no DOM trace come as data attributes on the root and are
// removed: data-type, data-scroll-hide-delay.
import { define } from "./core.js"
import { createPresence } from "./presence.js"
import { getNonce } from "./style.js"

const VIEWPORT_CSS =
  "[data-radix-scroll-area-viewport]{scrollbar-width:none;-ms-overflow-style:none;-webkit-overflow-scrolling:touch;}[data-radix-scroll-area-viewport]::-webkit-scrollbar{display:none}"

const savedOptions = new WeakMap()
function takeOption(el, name) {
  let saved = savedOptions.get(el)
  if (!saved) savedOptions.set(el, (saved = {}))
  if (el.hasAttribute(name)) {
    saved[name] = el.getAttribute(name)
    el.removeAttribute(name)
  }
  return saved[name] ?? null
}

// --- Radix helpers ------------------------------------------------------------------------

const clamp = (value, [min, max]) => Math.min(max, Math.max(min, value))
const toInt = (value) => (value ? parseInt(value, 10) : 0)

function getThumbRatio(viewportSize, contentSize) {
  const ratio = viewportSize / contentSize
  return isNaN(ratio) ? 0 : ratio
}

function getThumbSize(sizes) {
  const ratio = getThumbRatio(sizes.viewport, sizes.content)
  const scrollbarPadding = sizes.scrollbar.paddingStart + sizes.scrollbar.paddingEnd
  const thumbSize = (sizes.scrollbar.size - scrollbarPadding) * ratio
  return Math.max(thumbSize, 18)
}

function getScrollPositionFromPointer(pointerPos, pointerOffset, sizes, dir = "ltr") {
  const thumbSizePx = getThumbSize(sizes)
  const thumbCenter = thumbSizePx / 2
  const offset = pointerOffset || thumbCenter
  const thumbOffsetFromEnd = thumbSizePx - offset
  const minPointerPos = sizes.scrollbar.paddingStart + offset
  const maxPointerPos = sizes.scrollbar.size - sizes.scrollbar.paddingEnd - thumbOffsetFromEnd
  const maxScrollPos = sizes.content - sizes.viewport
  const scrollRange = dir === "ltr" ? [0, maxScrollPos] : [maxScrollPos * -1, 0]
  const interpolate = linearScale([minPointerPos, maxPointerPos], scrollRange)
  return interpolate(pointerPos)
}

function getThumbOffsetFromScroll(scrollPos, sizes, dir = "ltr") {
  const thumbSizePx = getThumbSize(sizes)
  const scrollbarPadding = sizes.scrollbar.paddingStart + sizes.scrollbar.paddingEnd
  const scrollbar = sizes.scrollbar.size - scrollbarPadding
  const maxScrollPos = sizes.content - sizes.viewport
  const maxThumbPos = scrollbar - thumbSizePx
  const scrollClampRange = dir === "ltr" ? [0, maxScrollPos] : [maxScrollPos * -1, 0]
  const scrollWithoutMomentum = clamp(scrollPos, scrollClampRange)
  const interpolate = linearScale([0, maxScrollPos], [0, maxThumbPos])
  return interpolate(scrollWithoutMomentum)
}

function linearScale(input, output) {
  return (value) => {
    if (input[0] === input[1] || output[0] === output[1]) return output[0]
    const ratio = (output[1] - output[0]) / (input[1] - input[0])
    return output[0] + ratio * (value - input[0])
  }
}

function isScrollingWithinScrollbarBounds(scrollPos, maxScrollPos) {
  return scrollPos > 0 && scrollPos < maxScrollPos
}

function addUnlinkedScrollListener(node, handler = () => {}) {
  let prevPosition = { left: node.scrollLeft, top: node.scrollTop }
  let rAF = 0
  ;(function loop() {
    const position = { left: node.scrollLeft, top: node.scrollTop }
    const isHorizontalScroll = prevPosition.left !== position.left
    const isVerticalScroll = prevPosition.top !== position.top
    if (isHorizontalScroll || isVerticalScroll) handler()
    prevPosition = position
    rAF = window.requestAnimationFrame(loop)
  })()
  return () => window.cancelAnimationFrame(rAF)
}

// useDebounceCallback
function debounce(fn, delay) {
  let timer = 0
  return {
    call() {
      window.clearTimeout(timer)
      timer = window.setTimeout(fn, delay)
    },
    cancel() {
      window.clearTimeout(timer)
    },
  }
}

// useResizeObserver: the callback runs in the animation frame after a resize.
function observeResize(element, onResize) {
  if (!element) return () => {}
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

// --- scrollbar (ScrollAreaScrollbarVisible + X/Y + Impl + Thumb) ---------------------------

class Bar {
  /**
   * @param {Area} area
   * @param {Slot} slot
   * @param {{ dataState?: string, onPointerEnter?: () => void, onPointerLeave?: () => void }} options
   */
  constructor(area, slot, { dataState, onPointerEnter, onPointerLeave } = {}) {
    this.area = area
    this.slot = slot
    this.horizontal = slot.orientation === "horizontal"
    this.sizes = { content: 0, viewport: 0, scrollbar: { size: 0, paddingStart: 0, paddingEnd: 0 } }
    this.pointerOffset = 0
    this.rect = null
    this.prevWebkitUserSelect = ""
    this.thumb = null
    this.unlinked = null
    this.cleanups = []
    this.exiting = false

    const el = slot.fragment.firstElementChild.cloneNode(true)
    this.thumbProto = el.querySelector('[data-slot="scroll-area-thumb"]')
    this.thumbProto?.remove()
    this.el = el
    if (dataState !== undefined) el.setAttribute("data-state", dataState)
    const { dir } = area
    el.style.position = "absolute"
    if (this.horizontal) {
      el.style.bottom = "0"
      el.style.left = dir === "rtl" ? "var(--radix-scroll-area-corner-width)" : "0"
      el.style.right = dir === "ltr" ? "var(--radix-scroll-area-corner-width)" : "0"
    } else {
      el.style.top = "0"
      if (dir === "ltr") el.style.right = "0"
      if (dir === "rtl") el.style.left = "0"
      el.style.bottom = "var(--radix-scroll-area-corner-height)"
    }
    // The user's style comes last (it was rendered on the template's element).
    const userStyle = slot.fragment.firstElementChild.getAttribute("style")
    this.renderThumbSize()
    if (userStyle) el.style.cssText += ";" + userStyle

    slot.anchor.before(el)
    this.presence = createPresence(el, { present: true, onExitComplete: () => this.onExit?.() })
    area.setScrollbar(slot.orientation, el)

    // ScrollAreaScrollbarX/Y: computed style (useEffect after mount).
    this.computedStyle = getComputedStyle(el)

    const on = (target, type, fn, opts) => {
      target.addEventListener(type, fn, opts)
      this.cleanups.push(() => target.removeEventListener(type, fn, opts))
    }
    if (onPointerEnter) on(el, "pointerenter", onPointerEnter)
    if (onPointerLeave) on(el, "pointerleave", onPointerLeave)

    // Impl: wheel, resize, drag.
    on(document, "wheel", (event) => {
      const element = event.target
      if (this.el.contains(element)) this.onWheelScroll(event, this.sizes.content - this.sizes.viewport)
    }, { passive: false })
    this.resize = debounce(() => this.onResize(), 10)
    this.cleanups.push(() => this.resize.cancel())
    this.cleanups.push(observeResize(el, () => this.resize.call()))
    this.cleanups.push(observeResize(area.content, () => this.resize.call()))

    on(el, "pointerdown", (event) => {
      const mainPointer = 0
      if (event.button === mainPointer) {
        const element = event.target
        element.setPointerCapture(event.pointerId)
        this.rect = this.el.getBoundingClientRect()
        this.prevWebkitUserSelect = document.body.style.webkitUserSelect
        document.body.style.webkitUserSelect = "none"
        if (area.viewport) area.viewport.style.scrollBehavior = "auto"
        this.handleDragScroll(event)
      }
    })
    on(el, "pointermove", (event) => this.handleDragScroll(event))
    on(el, "pointerup", (event) => {
      const element = event.target
      if (element.hasPointerCapture(event.pointerId)) element.releasePointerCapture(event.pointerId)
      document.body.style.webkitUserSelect = this.prevWebkitUserSelect
      if (area.viewport) area.viewport.style.scrollBehavior = ""
      this.rect = null
    })
    this.onThumbPositionChange()
  }

  setDataState(state) {
    this.el.setAttribute("data-state", state)
  }

  renderThumbSize() {
    const prop = this.horizontal ? "--radix-scroll-area-thumb-width" : "--radix-scroll-area-thumb-height"
    this.el.style.setProperty(prop, getThumbSize(this.sizes) + "px")
  }

  get hasThumb() {
    const thumbRatio = getThumbRatio(this.sizes.viewport, this.sizes.content)
    return Boolean(thumbRatio > 0 && thumbRatio < 1)
  }

  onResize() {
    const { viewport } = this.area
    if (!this.el || !viewport || !this.computedStyle) return
    const cs = this.computedStyle
    this.sizes = this.horizontal
      ? {
          content: viewport.scrollWidth,
          viewport: viewport.offsetWidth,
          scrollbar: { size: this.el.clientWidth, paddingStart: toInt(cs.paddingLeft), paddingEnd: toInt(cs.paddingRight) },
        }
      : {
          content: viewport.scrollHeight,
          viewport: viewport.offsetHeight,
          scrollbar: { size: this.el.clientHeight, paddingStart: toInt(cs.paddingTop), paddingEnd: toInt(cs.paddingBottom) },
        }
    // Re-render.
    this.renderThumbSize()
    this.updateThumbPresence()
    // useEffect(handleThumbPositionChange, [sizes])
    this.onThumbPositionChange()
  }

  // ScrollAreaThumb: Presence(present = hasThumb); no exit animation upstream.
  updateThumbPresence() {
    const has = this.hasThumb
    if (has && !this.thumb && this.thumbProto) this.mountThumb()
    else if (!has && this.thumb) {
      this.thumb.setAttribute("data-state", "hidden")
      const thumb = this.thumb
      const presence = this.thumbPresence
      presence.setPresent(false)
      if (presence.state === "unmounted" && thumb === this.thumb) this.unmountThumb()
    } else if (has && this.thumb && this.thumbPresence.state === "unmountSuspended") {
      this.thumb.setAttribute("data-state", "visible")
      this.thumbPresence.setPresent(true)
    }
  }

  mountThumb() {
    const thumb = this.thumbProto.cloneNode(true)
    thumb.setAttribute("data-state", "visible")
    const userStyle = thumb.getAttribute("style")
    thumb.style.cssText = "width: var(--radix-scroll-area-thumb-width); height: var(--radix-scroll-area-thumb-height);" + (userStyle ? ";" + userStyle : "")
    this.el.appendChild(thumb)
    this.thumb = thumb
    this.thumbPresence = createPresence(thumb, { present: true, onExitComplete: () => thumb === this.thumb && this.unmountThumb() })
    const viewport = this.area.viewport
    const onPointerDownCapture = (event) => {
      const t = event.target
      const thumbRect = t.getBoundingClientRect()
      const x = event.clientX - thumbRect.left
      const y = event.clientY - thumbRect.top
      this.pointerOffset = this.horizontal ? x : y
    }
    const onPointerUp = () => (this.pointerOffset = 0)
    const debounceScrollEnd = debounce(() => {
      if (this.unlinked) {
        this.unlinked()
        this.unlinked = null
      }
    }, 100)
    const handleScroll = () => {
      debounceScrollEnd.call()
      if (!this.unlinked) {
        this.unlinked = addUnlinkedScrollListener(viewport, () => this.onThumbPositionChange())
        this.onThumbPositionChange()
      }
    }
    thumb.addEventListener("pointerdown", onPointerDownCapture, true)
    thumb.addEventListener("pointerup", onPointerUp)
    viewport.addEventListener("scroll", handleScroll)
    this.thumbCleanup = () => {
      thumb.removeEventListener("pointerdown", onPointerDownCapture, true)
      thumb.removeEventListener("pointerup", onPointerUp)
      viewport.removeEventListener("scroll", handleScroll)
      debounceScrollEnd.cancel()
      this.thumbPresence.destroy()
    }
    this.onThumbPositionChange()
  }

  unmountThumb() {
    this.thumbCleanup?.()
    this.thumbCleanup = null
    this.thumb?.remove()
    this.thumb = null
  }

  onThumbPositionChange() {
    const { viewport, dir } = this.area
    if (!viewport || !this.thumb) return
    if (this.horizontal) {
      const offset = getThumbOffsetFromScroll(viewport.scrollLeft, this.sizes, dir)
      this.thumb.style.transform = `translate3d(${offset}px, 0, 0)`
    } else {
      const offset = getThumbOffsetFromScroll(viewport.scrollTop, this.sizes)
      this.thumb.style.transform = `translate3d(0, ${offset}px, 0)`
    }
  }

  onWheelScroll(event, maxScrollPos) {
    const { viewport } = this.area
    if (!viewport) return
    if (this.horizontal) {
      const scrollPos = viewport.scrollLeft + event.deltaX
      viewport.scrollLeft = scrollPos
      if (isScrollingWithinScrollbarBounds(scrollPos, maxScrollPos)) event.preventDefault()
    } else {
      const scrollPos = viewport.scrollTop + event.deltaY
      viewport.scrollTop = scrollPos
      if (isScrollingWithinScrollbarBounds(scrollPos, maxScrollPos)) event.preventDefault()
    }
  }

  handleDragScroll(event) {
    if (!this.rect) return
    const x = event.clientX - this.rect.left
    const y = event.clientY - this.rect.top
    const { viewport, dir } = this.area
    if (!viewport) return
    if (this.horizontal) viewport.scrollLeft = getScrollPositionFromPointer(x, this.pointerOffset, this.sizes, dir)
    else viewport.scrollTop = getScrollPositionFromPointer(y, this.pointerOffset, this.sizes)
  }

  // Leaves (presence): onDone runs once the exit animation, if any, ended.
  exit(dataState, onDone) {
    if (dataState !== undefined) this.setDataState(dataState)
    this.exiting = true
    this.onExit = () => {
      this.destroy()
      onDone()
    }
    this.presence.setPresent(false)
  }

  // Re-shown during its exit animation.
  cancelExit(dataState) {
    if (dataState !== undefined) this.setDataState(dataState)
    this.exiting = false
    this.onExit = null
    this.presence.setPresent(true)
  }

  destroy() {
    if (!this.el) return
    this.unlinked?.()
    this.unlinked = null
    this.unmountThumb()
    this.cleanups.forEach((fn) => fn())
    this.cleanups = []
    this.presence.destroy()
    this.el.remove()
    this.el = null
    this.area.setScrollbar(this.slot.orientation, null)
  }
}

// A Presence around a Bar: mounts it while present, lets it exit otherwise.
class Gate {
  constructor(create, onUnmounted) {
    this.create = create
    this.onUnmounted = onUnmounted
    this.bar = null
  }
  set(present, dataState) {
    if (present) {
      if (!this.bar) this.bar = this.create()
      else if (this.bar.exiting) this.bar.cancelExit(dataState)
      else if (dataState !== undefined) this.bar.setDataState(dataState)
    } else if (this.bar && !this.bar.exiting) {
      const bar = this.bar
      bar.exit(dataState, () => {
        if (this.bar === bar) this.bar = null
        this.onUnmounted?.()
      })
    }
  }
  destroy() {
    this.bar?.destroy()
    this.bar = null
  }
}

// ScrollAreaScrollbarAuto: visible while the content overflows.
class Auto {
  constructor(area, slot, { dataState } = {}) {
    this.area = area
    this.slot = slot
    this.visible = false
    this.dataStateOverride = dataState // ScrollAreaScrollbarHover's data-state
    this.gate = new Gate(() => new Bar(area, slot, { dataState: this.dataState() }))
    const handleResize = debounce(() => {
      const { viewport } = area
      if (viewport) {
        const isOverflowX = viewport.offsetWidth < viewport.scrollWidth
        const isOverflowY = viewport.offsetHeight < viewport.scrollHeight
        this.setVisible(slot.orientation === "horizontal" ? isOverflowX : isOverflowY)
      }
    }, 10)
    this.cleanups = [
      () => handleResize.cancel(),
      observeResize(area.viewport, () => handleResize.call()),
      observeResize(area.content, () => handleResize.call()),
    ]
  }
  dataState() {
    return this.dataStateOverride ?? (this.visible ? "visible" : "hidden")
  }
  setVisible(visible) {
    if (visible === this.visible) return
    this.visible = visible
    this.gate.set(visible, this.dataState())
  }
  setDataStateOverride(state) {
    this.dataStateOverride = state
    if (this.gate.bar) this.gate.bar.setDataState(this.dataState())
  }
  get bar() {
    return this.gate.bar
  }
  destroy() {
    this.cleanups.forEach((fn) => fn())
    this.gate.destroy()
  }
}

// ScrollAreaScrollbarHover: an Auto mounted while the pointer is over the root.
class Hover {
  constructor(area, slot) {
    this.visible = false
    this.auto = null
    this.exiting = false
    let hideTimer = 0
    const root = area.root
    const handlePointerEnter = () => {
      window.clearTimeout(hideTimer)
      this.setVisible(true)
    }
    const handlePointerLeave = () => {
      hideTimer = window.setTimeout(() => this.setVisible(false), area.scrollHideDelay)
    }
    root.addEventListener("pointerenter", handlePointerEnter)
    root.addEventListener("pointerleave", handlePointerLeave)
    this.cleanup = () => {
      window.clearTimeout(hideTimer)
      root.removeEventListener("pointerenter", handlePointerEnter)
      root.removeEventListener("pointerleave", handlePointerLeave)
    }
    this.area = area
    this.slot = slot
  }
  setVisible(visible) {
    if (visible === this.visible) return
    this.visible = visible
    const state = visible ? "visible" : "hidden"
    if (visible) {
      if (!this.auto) this.auto = new Auto(this.area, this.slot, { dataState: state })
      else {
        this.auto.setDataStateOverride(state)
        if (this.exiting) {
          this.exiting = false
          this.auto.bar?.cancelExit(state)
        }
      }
    } else if (this.auto) {
      const auto = this.auto
      auto.setDataStateOverride(state)
      const bar = auto.bar
      const done = () => {
        this.exiting = false
        auto.destroy()
        if (this.auto === auto) this.auto = null
      }
      if (bar && !bar.exiting) {
        this.exiting = true
        bar.exit(state, () => {
          auto.gate.bar = null
          done()
        })
      } else done()
    }
  }
  destroy() {
    this.cleanup()
    this.auto?.destroy()
  }
}

// ScrollAreaScrollbarScroll: shown while scrolling.
const SCROLL_MACHINE = {
  hidden: { SCROLL: "scrolling" },
  scrolling: { SCROLL_END: "idle", POINTER_ENTER: "interacting" },
  interacting: { SCROLL: "interacting", POINTER_LEAVE: "idle" },
  idle: { HIDE: "hidden", SCROLL: "scrolling", POINTER_ENTER: "interacting" },
}

class Scroll {
  constructor(area, slot) {
    this.area = area
    this.state = "hidden"
    this.hideTimer = 0
    const isHorizontal = slot.orientation === "horizontal"
    this.gate = new Gate(() =>
      new Bar(area, slot, {
        dataState: this.dataState(),
        onPointerEnter: () => this.send("POINTER_ENTER"),
        onPointerLeave: () => this.send("POINTER_LEAVE"),
      }),
    )
    const debounceScrollEnd = debounce(() => this.send("SCROLL_END"), 100)
    const viewport = area.viewport
    const scrollDirection = isHorizontal ? "scrollLeft" : "scrollTop"
    let prevScrollPos = viewport[scrollDirection]
    const handleScroll = () => {
      const scrollPos = viewport[scrollDirection]
      const hasScrollInDirectionChanged = prevScrollPos !== scrollPos
      if (hasScrollInDirectionChanged) {
        this.send("SCROLL")
        debounceScrollEnd.call()
      }
      prevScrollPos = scrollPos
    }
    viewport.addEventListener("scroll", handleScroll)
    this.cleanup = () => {
      viewport.removeEventListener("scroll", handleScroll)
      debounceScrollEnd.cancel()
      window.clearTimeout(this.hideTimer)
    }
  }
  dataState() {
    return this.state === "hidden" ? "hidden" : "visible"
  }
  send(event) {
    const next = SCROLL_MACHINE[this.state][event]
    if (!next || next === this.state) return
    this.state = next
    window.clearTimeout(this.hideTimer)
    if (next === "idle") this.hideTimer = window.setTimeout(() => this.send("HIDE"), this.area.scrollHideDelay)
    this.gate.set(next !== "hidden", this.dataState())
  }
  destroy() {
    this.cleanup()
    this.gate.destroy()
  }
}

// ScrollAreaCorner: shown while both scrollbars are mounted (not for type scroll).
class Corner {
  constructor(area) {
    this.area = area
    this.el = null
    this.width = 0
    this.height = 0
    this.offX = () => {}
    this.offY = () => {}
    this.observedX = null
    this.observedY = null
    this.active = false
  }
  update() {
    const { area } = this
    const hasCorner = area.type !== "scroll" && Boolean(area.scrollbarX && area.scrollbarY)
    if (!hasCorner) {
      if (this.active) {
        this.active = false
        this.offX()
        this.offY()
        this.observedX = this.observedY = null
        this.width = this.height = 0
        this.el?.remove()
        this.el = null
      }
      return
    }
    this.active = true
    if (this.observedX !== area.scrollbarX) {
      this.offX()
      this.observedX = area.scrollbarX
      this.offX = observeResize(area.scrollbarX, () => {
        const height = area.scrollbarX?.offsetHeight || 0
        area.setCornerHeight(height)
        this.height = height
        this.render()
      })
    }
    if (this.observedY !== area.scrollbarY) {
      this.offY()
      this.observedY = area.scrollbarY
      this.offY = observeResize(area.scrollbarY, () => {
        const width = area.scrollbarY?.offsetWidth || 0
        area.setCornerWidth(width)
        this.width = width
        this.render()
      })
    }
  }
  render() {
    const hasSize = Boolean(this.width && this.height)
    if (!hasSize) {
      this.el?.remove()
      this.el = null
      return
    }
    if (!this.el) {
      this.el = document.createElement("div")
      this.area.cornerAnchor.before(this.el)
    }
    const { dir } = this.area
    const s = this.el.style
    s.width = this.width + "px"
    s.height = this.height + "px"
    s.position = "absolute"
    if (dir === "ltr") s.right = "0"
    if (dir === "rtl") s.left = "0"
    s.bottom = "0"
  }
  destroy() {
    this.offX()
    this.offY()
    this.el?.remove()
  }
}

class Area {
  constructor(root) {
    this.root = root
    this.type = takeOption(root, "data-type") || "hover"
    const delay = takeOption(root, "data-scroll-hide-delay")
    this.scrollHideDelay = delay === null ? 600 : Number(delay)
    this.dir = root.getAttribute("dir") === "rtl" ? "rtl" : "ltr"
    this.viewport = [...root.children].find((el) => el.hasAttribute("data-radix-scroll-area-viewport")) ?? null
    this.content = this.viewport?.firstElementChild ?? null
    this.scrollbarX = null
    this.scrollbarY = null
    this.cleanups = []

    // The <style> needs the CSP nonce: insert it if the server could not.
    if (this.viewport && !(this.viewport.previousElementSibling?.tagName === "STYLE")) {
      const style = document.createElement("style")
      const nonce = getNonce()
      if (nonce) style.setAttribute("nonce", nonce)
      style.textContent = VIEWPORT_CSS
      this.viewport.before(style)
      this.cleanups.push(() => style.remove())
    }

    // Scrollbar templates of this area (not of nested areas).
    const templates = [...root.querySelectorAll('template[data-slot="scroll-area-scrollbar-template"]')].filter(
      (t) => t.parentElement?.closest('[data-slot="scroll-area"]') === root,
    )
    this.slots = templates.map((tpl) => {
      const anchor = document.createComment("")
      tpl.replaceWith(anchor)
      const orientation = tpl.content.firstElementChild?.getAttribute("data-orientation") === "horizontal" ? "horizontal" : "vertical"
      return { anchor, fragment: tpl.content, orientation }
    })
    this.cornerAnchor = document.createComment("")
    root.appendChild(this.cornerAnchor)
    this.corner = new Corner(this)

    this.scrollbars = this.slots.map((slot) => {
      switch (this.type) {
        case "scroll":
          return new Scroll(this, slot)
        case "auto":
          return new Auto(this, slot)
        case "always": {
          const gate = new Gate(() => new Bar(this, slot))
          gate.set(true)
          return gate
        }
        default:
          return new Hover(this, slot)
      }
    })
  }

  setScrollbar(orientation, el) {
    if (orientation === "horizontal") this.scrollbarX = el
    else this.scrollbarY = el
    this.corner?.update()
  }

  setCornerWidth(w) {
    this.root.style.setProperty("--radix-scroll-area-corner-width", w + "px")
  }

  setCornerHeight(hgt) {
    this.root.style.setProperty("--radix-scroll-area-corner-height", hgt + "px")
  }

  destroy() {
    this.scrollbars.forEach((s) => s?.destroy())
    this.corner.destroy()
    this.cleanups.forEach((fn) => fn())
  }
}

define("scroll-area", {
  selector: '[data-slot="scroll-area"]',
  mount(root) {
    const area = new Area(root)
    return () => area.destroy()
  },
})
