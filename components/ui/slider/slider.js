// slider.js: Radix Slider (@radix-ui/react-slider 1.3.6).
//
// - Pointer down on the slider (not on a thumb) moves the closest thumb to
//   the pointer and focuses it; moving the pointer (captured) drags it.
//   Pointer down on a thumb focuses it and drags it. Values snap to the
//   step, stay within [min, max] and keep minStepsBetweenThumbs steps
//   between thumbs (otherwise the change is ignored); values stay sorted,
//   and the moved value's thumb gets focus.
// - Keys (on the root, from a focused thumb): Home sets the first value to
//   min, End the last value to max, arrows step by one step (10 with
//   Shift), Page Up/Down by 10 steps. Which arrows decrease depends on
//   orientation, direction and inverted, as in Radix's BACK_KEYS.
// - The thumbs' wrapper and the range are positioned from the values, with
//   Radix's in-bounds offset from the measured thumb size (useSize).
// - Inside a form each thumb's hidden <input> (Radix SliderBubbleInput)
//   carries its value; an "input" event is dispatched on change. Outside a
//   form the inputs are removed.
// - "shadcn:value-change" (cancelable, detail: { value }) is dispatched on
//   the root before every change, "shadcn:value-commit" (detail: { value })
//   when a keyboard change or a drag ends with a new value.
//
// Radix props with no DOM trace are server-rendered as data attributes
// (data-step, data-min-steps-between-thumbs, data-inverted) and removed
// here, so the DOM matches React's.
import { define, emit } from "./core.js"
import { getDirection } from "./direction.js"
import { takeBubbleInput, setValue as bubbleValue } from "./bubble-input.js"

const PAGE_KEYS = ["PageUp", "PageDown"]
const ARROW_KEYS = ["ArrowUp", "ArrowDown", "ArrowLeft", "ArrowRight"]
const BACK_KEYS = {
  "from-left": ["Home", "PageDown", "ArrowDown", "ArrowLeft"],
  "from-right": ["Home", "PageDown", "ArrowDown", "ArrowRight"],
  "from-bottom": ["Home", "PageDown", "ArrowDown", "ArrowLeft"],
  "from-top": ["Home", "PageDown", "ArrowUp", "ArrowLeft"],
}

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

define("slider", {
  selector: '[data-slot="slider"]',
  mount(root) {
    const cleanups = []
    const on = (el, type, fn, opts) => {
      el.addEventListener(type, fn, opts)
      cleanups.push(() => el.removeEventListener(type, fn, opts))
    }
    const own = (el) => el.closest('[data-slot="slider"]') === root

    const thumbs = [...root.querySelectorAll('[role="slider"]')].filter(own)
    if (!thumbs.length) return
    const range = [...root.querySelectorAll('[data-slot="slider-range"]')].find(own)
    const wrappers = thumbs.map((t) => t.parentElement)
    const inputs = thumbs.map((t) => {
      const next = t.nextElementSibling
      // Radix's form prop is on the inputs only (not on the thumbs).
      return next?.tagName === "INPUT" ? takeBubbleInput(t, next, next) : null
    })

    const min = Number(thumbs[0].getAttribute("aria-valuemin"))
    const max = Number(thumbs[0].getAttribute("aria-valuemax"))
    const step = Number(takeOption(root, "data-step") ?? 1)
    const minStepsBetweenThumbs = Number(takeOption(root, "data-min-steps-between-thumbs") ?? 0)
    const inverted = takeOption(root, "data-inverted") !== null
    const disabled = root.hasAttribute("data-disabled")
    const horizontal = root.getAttribute("data-orientation") !== "vertical"

    let values = []
    for (const t of thumbs) {
      if (!t.hasAttribute("aria-valuenow")) break
      values.push(Number(t.getAttribute("aria-valuenow")))
    }
    let valueIndexToChange = 0
    let valuesBeforeSlideStart = values
    let rect

    // Orientation context (SliderHorizontal / SliderVertical).
    const orientation = () => {
      if (horizontal) {
        const fromLeft = (getDirection(root, root.getAttribute("dir")) === "ltr") !== inverted
        return fromLeft
          ? { start: "left", end: "right", size: "width", direction: 1, slide: "from-left" }
          : { start: "right", end: "left", size: "width", direction: -1, slide: "from-right" }
      }
      return inverted
        ? { start: "top", end: "bottom", size: "height", direction: -1, slide: "from-top" }
        : { start: "bottom", end: "top", size: "height", direction: 1, slide: "from-bottom" }
    }

    const render = () => {
      const o = orientation()
      const percentages = values.map((v) => convertValueToPercentage(v, min, max))
      if (range) {
        range.style[o.start] = (values.length > 1 ? Math.min(...percentages) : 0) + "%"
        range.style[o.end] = 100 - Math.max(...percentages) + "%"
      }
      thumbs.forEach((thumb, i) => {
        const value = values[i]
        const percent = value === undefined ? 0 : convertValueToPercentage(value, min, max)
        if (value === undefined) {
          thumb.removeAttribute("aria-valuenow")
          thumb.style.display = "none"
        } else {
          thumb.setAttribute("aria-valuenow", String(value))
          thumb.style.display = ""
        }
        const label = getLabel(i, values.length)
        if (label) thumb.setAttribute("aria-label", label)
        // useSize: the thumb's offsetWidth/offsetHeight.
        const size = o.size === "width" ? thumb.offsetWidth : thumb.offsetHeight
        const offset = size ? getThumbInBoundsOffset(size, percent, o.direction) : 0
        wrappers[i].style[o.start] = `calc(${percent}% + ${offset}px)`
        const input = inputs[i]
        if (input && value !== undefined) input.defaultValue = String(value)
      })
    }
    render()
    if (typeof ResizeObserver !== "undefined") {
      const ro = new ResizeObserver(() => render())
      thumbs.forEach((t) => ro.observe(t, { box: "border-box" }))
      cleanups.push(() => ro.disconnect())
    }

    const setValues = (next) => {
      if (!emit(root, "shadcn:value-change", { value: next })) return false
      const prev = values
      values = next
      render()
      // useControllableState onChange: focus the moved thumb.
      thumbs[valueIndexToChange]?.focus()
      // Bubble inputs' effect.
      inputs.forEach((input, i) => {
        if (input && prev[i] !== values[i]) bubbleValue(input, values[i])
      })
      return true
    }

    const updateValues = (value, atIndex, { commit = false } = {}) => {
      const decimalCount = getDecimalCount(step)
      const snapToStep = roundValue(Math.round((value - min) / step) * step + min, decimalCount)
      const nextValue = clamp(snapToStep, [min, max])
      const nextValues = getNextSortedValues(values, nextValue, atIndex)
      if (!hasMinStepsBetweenValues(nextValues, minStepsBetweenThumbs * step)) return
      valueIndexToChange = nextValues.indexOf(nextValue)
      if (String(nextValues) === String(values)) return
      if (setValues(nextValues) && commit) emit(root, "shadcn:value-commit", { value: values })
    }

    const getValueFromPointer = (event) => {
      const r = rect || root.getBoundingClientRect()
      rect = r
      const o = orientation()
      if (horizontal) {
        const output = o.direction === 1 ? [min, max] : [max, min]
        return linearScale([0, r.width], output)(event.clientX - r.left)
      }
      const output = o.direction === 1 ? [max, min] : [min, max]
      return linearScale([0, r.height], output)(event.clientY - r.top)
    }

    on(root, "keydown", (event) => {
      if (event.key === "Home") {
        if (!disabled) updateValues(min, 0, { commit: true })
        event.preventDefault()
      } else if (event.key === "End") {
        if (!disabled) updateValues(max, values.length - 1, { commit: true })
        event.preventDefault()
      } else if (PAGE_KEYS.concat(ARROW_KEYS).includes(event.key)) {
        const isBackKey = BACK_KEYS[orientation().slide].includes(event.key)
        if (!disabled) {
          const isPageKey = PAGE_KEYS.includes(event.key)
          const isSkipKey = isPageKey || (event.shiftKey && ARROW_KEYS.includes(event.key))
          const multiplier = isSkipKey ? 10 : 1
          const atIndex = valueIndexToChange
          const value = values[atIndex]
          updateValues(value + step * multiplier * (isBackKey ? -1 : 1), atIndex, { commit: true })
        }
        event.preventDefault()
      }
    })

    on(root, "pointerdown", (event) => {
      if (!disabled) valuesBeforeSlideStart = values
      const target = event.target
      target.setPointerCapture(event.pointerId)
      event.preventDefault()
      if (thumbs.includes(target)) target.focus()
      else {
        const value = getValueFromPointer(event)
        if (!disabled) updateValues(value, getClosestValueIndex(values, value))
      }
    })
    on(root, "pointermove", (event) => {
      const target = event.target
      if (target.hasPointerCapture(event.pointerId)) {
        const value = getValueFromPointer(event)
        if (!disabled) updateValues(value, valueIndexToChange)
      }
    })
    on(root, "pointerup", (event) => {
      const target = event.target
      if (target.hasPointerCapture(event.pointerId)) {
        target.releasePointerCapture(event.pointerId)
        rect = undefined
        if (!disabled && values[valueIndexToChange] !== valuesBeforeSlideStart[valueIndexToChange]) {
          emit(root, "shadcn:value-commit", { value: values })
        }
      }
    })
    thumbs.forEach((thumb, i) => on(thumb, "focus", () => (valueIndexToChange = i)))

    return () => cleanups.forEach((fn) => fn())
  },
})

function clamp(value, [lo, hi]) {
  return Math.min(hi, Math.max(lo, value))
}

function getNextSortedValues(prevValues = [], nextValue, atIndex) {
  const nextValues = [...prevValues]
  nextValues[atIndex] = nextValue
  return nextValues.sort((a, b) => a - b)
}

function convertValueToPercentage(value, min, max) {
  const maxSteps = max - min
  const percentPerStep = 100 / maxSteps
  const percentage = percentPerStep * (value - min)
  return clamp(percentage, [0, 100])
}

function getLabel(index, totalValues) {
  if (totalValues > 2) return `Value ${index + 1} of ${totalValues}`
  if (totalValues === 2) return ["Minimum", "Maximum"][index]
  return undefined
}

function getClosestValueIndex(values, nextValue) {
  if (values.length === 1) return 0
  const distances = values.map((value) => Math.abs(value - nextValue))
  const closestDistance = Math.min(...distances)
  return distances.indexOf(closestDistance)
}

function getThumbInBoundsOffset(width, left, direction) {
  const halfWidth = width / 2
  const halfPercent = 50
  const offset = linearScale([0, halfPercent], [0, halfWidth])
  return (halfWidth - offset(left) * direction) * direction
}

function getStepsBetweenValues(values) {
  return values.slice(0, -1).map((value, index) => values[index + 1] - value)
}

function hasMinStepsBetweenValues(values, minStepsBetweenValues) {
  if (minStepsBetweenValues > 0) {
    const stepsBetweenValues = getStepsBetweenValues(values)
    const actualMinStepsBetweenValues = Math.min(...stepsBetweenValues)
    return actualMinStepsBetweenValues >= minStepsBetweenValues
  }
  return true
}

function linearScale(input, output) {
  return (value) => {
    if (input[0] === input[1] || output[0] === output[1]) return output[0]
    const ratio = (output[1] - output[0]) / (input[1] - input[0])
    return output[0] + ratio * (value - input[0])
  }
}

function getDecimalCount(value) {
  return (String(value).split(".")[1] || "").length
}

function roundValue(value, decimalCount) {
  const rounder = Math.pow(10, decimalCount)
  return Math.round(value * rounder) / rounder
}
