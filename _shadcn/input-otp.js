// input-otp.js: behavior of shadcn/ui InputOTP, a port of the input-otp
// library (v1.4.2, ui/apps/v4/node_modules/input-otp/dist/index.mjs).
//
// One transparent <input data-input-otp> lies over the visual slots. The
// library keeps the caret on one character (selecting it), so typing
// replaces it; the slots render that selection as data-active and a fake
// caret. Everything here mirrors the library's React state and effects,
// including their timing (selectionchange handling, the 0/10/50 ms resyncs
// after each value or focus change, and the password-manager badge check
// that widens the input by 40px once it was focused).
//
// Events (bubbling, from the input): "shadcn:value-change" {value} after
// every change, "shadcn:complete" {value} when the last slot gets filled.
import { define, emit } from "./core.js"

const PWM_BADGE_MARGIN_RIGHT = 18
const PWM_BADGE_SPACE_WIDTH_PX = 40
const PWM_BADGE_SPACE_WIDTH = `${PWM_BADGE_SPACE_WIDTH_PX}px`
const PASSWORD_MANAGERS_SELECTORS = [
  "[data-lastpass-icon-root]",
  "com-1password-button",
  "[data-dashlanecreated]",
  '[style$="2147483647 !important;"]',
].join(",")

const CARET_WRAPPER_CLASS = "pointer-events-none absolute inset-0 flex items-center justify-center"
const CARET_CLASS = "h-4 w-px animate-caret-blink bg-foreground duration-1000"

// The rules input-otp injects into <head> once per document.
function injectStyle() {
  if (document.getElementById("input-otp-style")) return
  const style = document.createElement("style")
  style.id = "input-otp-style"
  const nonce = document.querySelector("script[nonce]")?.nonce
  if (nonce) style.nonce = nonce
  document.head.appendChild(style)
  if (!style.sheet) return
  const autofill =
    "background: transparent !important; color: transparent !important; border-color: transparent !important; opacity: 0 !important; box-shadow: none !important; -webkit-box-shadow: none !important; -webkit-text-fill-color: transparent !important;"
  const rules = [
    "[data-input-otp]::selection { background: transparent !important; color: transparent !important; }",
    `[data-input-otp]:autofill { ${autofill} }`,
    `[data-input-otp]:-webkit-autofill { ${autofill} }`,
    "@supports (-webkit-touch-callout: none) { [data-input-otp] { letter-spacing: -.6em !important; font-weight: 100 !important; font-stretch: ultra-condensed; font-optical-sizing: none !important; left: -1px !important; right: 1px !important; } }",
    "[data-input-otp] + * { pointer-events: all !important; }",
  ]
  for (const rule of rules) {
    try {
      style.sheet.insertRule(rule)
    } catch {
      console.error("input-otp could not insert CSS rule:", rule)
    }
  }
}

function timeouts(fn) {
  return [setTimeout(fn, 0), setTimeout(fn, 10), setTimeout(fn, 50)]
}

define("input-otp", {
  selector: "[data-input-otp-container]",
  mount(container) {
    const input = container.querySelector(":scope > div > input[data-input-otp]")
    if (!input) return
    const slots = [...container.querySelectorAll('[data-slot="input-otp-slot"]')].filter(
      (el) => el.closest("[data-input-otp-container]") === container,
    )
    const maxLength = Number(input.getAttribute("maxlength")) || 0
    const pattern = input.getAttribute("pattern")
    const regexp = pattern ? new RegExp(pattern) : null
    const isIOS = typeof window !== "undefined" && window.CSS?.supports?.("-webkit-touch-callout", "none")

    // React state.
    const state = {
      value: input.value.slice(0, maxLength || undefined),
      isFocused: false,
      isHovering: false,
      mirrorSelectionStart: null,
      mirrorSelectionEnd: null,
      hasPWMBadge: false,
      hasPWMBadgeSpace: false,
      done: false,
    }
    // usePrevious(value): the value of the previous render.
    let previousValue
    let renderedValue = state.value
    // inputMetadataRef.prev
    let prevSelection = [input.selectionStart, input.selectionEnd, input.selectionDirection]
    // React's input value tracker: synthetic "input" events with an
    // unchanged value do not reach onChange.
    let trackedValue = input.value

    // Effect dependencies, to rerun effects like React does.
    let lastSyncDeps = null
    let syncTimers = []
    let pwmDeps = null
    let pwmTimers = []
    let batching = 0
    let dirty = false

    // React renders updates made in event handlers before the next task
    // (sync lane) and updates made in timers in a later scheduler task
    // (default lane), possibly after a frame was painted: setState and
    // setStateLater mirror that.
    const setState = (patch) => {
      if (!apply(patch)) return
      if (batching) dirty = true
      else commit()
    }
    let laterScheduled = false
    const channel = new MessageChannel()
    channel.port1.onmessage = () => {
      laterScheduled = false
      if (input.isConnected) commit()
    }
    const setStateLater = (patch) => {
      if (!apply(patch) || laterScheduled) return
      laterScheduled = true
      channel.port2.postMessage(null)
    }
    const apply = (patch) => {
      let changed = false
      for (const k in patch) {
        if (!Object.is(state[k], patch[k])) {
          state[k] = patch[k]
          changed = true
        }
      }
      return changed
    }
    const batch = (fn) => {
      batching++
      try {
        fn()
      } finally {
        batching--
        if (!batching && dirty) {
          dirty = false
          commit()
        }
      }
    }

    const setValue = (v) => {
      emit(input, "shadcn:value-change", { value: v })
      setState({ value: v })
    }

    // --- render -------------------------------------------------------------------

    const willPushPWMBadge = () => state.hasPWMBadge && state.hasPWMBadgeSpace

    function render() {
      previousValue = renderedValue
      renderedValue = state.value
      const v = state.value
      // Controlled <input value>: React sets the property, then syncs the
      // value attribute.
      if (input.value !== v) {
        input.value = v
      }
      trackedValue = v
      if (input.defaultValue !== v) input.defaultValue = v
      if (v.length === 0) input.setAttribute("data-input-otp-placeholder-shown", "true")
      else input.removeAttribute("data-input-otp-placeholder-shown")
      const mss = state.mirrorSelectionStart
      const mse = state.mirrorSelectionEnd
      if (mss === null) input.removeAttribute("data-input-otp-mss")
      else input.setAttribute("data-input-otp-mss", String(mss))
      if (mse === null) input.removeAttribute("data-input-otp-mse")
      else input.setAttribute("data-input-otp-mse", String(mse))
      const push = willPushPWMBadge()
      input.style.width = push ? `calc(100% + ${PWM_BADGE_SPACE_WIDTH})` : "100%"
      input.style.clipPath = push ? `inset(0 ${PWM_BADGE_SPACE_WIDTH} 0 0)` : ""

      slots.forEach((slot, i) => {
        const isActive =
          state.isFocused &&
          mss !== null &&
          mse !== null &&
          ((mss === mse && i === mss) || (i >= mss && i < mse))
        const char = v[i] !== undefined ? v[i] : null
        const hasFakeCaret = isActive && char === null
        slot.setAttribute("data-active", String(isActive))
        renderSlot(slot, char, hasFakeCaret)
      })
    }

    function renderSlot(slot, char, hasFakeCaret) {
      let text = null
      let caret = null
      for (const n of [...slot.childNodes]) {
        if (n.nodeType === 3 && text === null) text = n
        else if (n.nodeType === 1 && caret === null && n.firstElementChild?.classList.contains("animate-caret-blink")) caret = n
        else n.remove()
      }
      if (char !== null) {
        if (text) {
          if (text.data !== char) text.data = char
        } else {
          text = document.createTextNode(char)
          slot.insertBefore(text, slot.firstChild)
        }
      } else if (text) {
        text.remove()
      }
      if (hasFakeCaret) {
        if (!caret) {
          caret = document.createElement("div")
          caret.className = CARET_WRAPPER_CLASS
          const line = document.createElement("div")
          line.className = CARET_CLASS
          caret.appendChild(line)
          slot.appendChild(caret)
        }
      } else if (caret) {
        caret.remove()
      }
    }

    // --- effects ------------------------------------------------------------------

    function commit() {
      render()
      runEffects()
    }

    function runEffects() {
      // [value, isFocused]: resync the selection a few times.
      const syncKey = state.value + "\u0000" + state.isFocused
      if (syncKey !== lastSyncDeps) {
        lastSyncDeps = syncKey
        syncTimers.push(...timeouts(syncSelection))
      }
      // onComplete, compared with the previous render's value.
      const prev = previousValue
      if (prev !== undefined && state.value !== prev && prev.length < maxLength && state.value.length === maxLength) {
        emit(input, "shadcn:complete", { value: state.value })
      }
      // Password manager badge: check once focused.
      const pwmKey = state.isFocused + "\u0000" + state.done
      if (pwmKey !== pwmDeps) {
        pwmDeps = pwmKey
        pwmTimers.forEach(clearTimeout)
        pwmTimers = []
        if (state.isFocused || document.activeElement === input) {
          pwmTimers = [
            setTimeout(checkPWMBadge, 0),
            setTimeout(checkPWMBadge, 2000),
            setTimeout(checkPWMBadge, 5000),
            setTimeout(() => setStateLater({ done: true }), 6000),
          ]
        }
      }
    }

    function syncSelection() {
      if (!input.isConnected) return
      input.dispatchEvent(new Event("input"))
      const s = input.selectionStart
      const e = input.selectionEnd
      const dir = input.selectionDirection
      if (s !== null && e !== null) {
        setStateLater({ mirrorSelectionStart: s, mirrorSelectionEnd: e })
        prevSelection = [s, e, dir]
      }
    }

    function checkPWMBadge() {
      if (state.done || !input.isConnected) return
      const rect = container.getBoundingClientRect()
      const x = rect.left + container.offsetWidth - PWM_BADGE_MARGIN_RIGHT
      const y = rect.top + container.offsetHeight / 2
      if (document.querySelectorAll(PASSWORD_MANAGERS_SELECTORS).length === 0 && document.elementFromPoint(x, y) === container) return
      setStateLater({ hasPWMBadge: true, done: true })
    }

    function trackPWMBadgeSpace() {
      const space = window.innerWidth - container.getBoundingClientRect().right
      setStateLater({ hasPWMBadgeSpace: space >= PWM_BADGE_SPACE_WIDTH_PX })
    }

    // --- handlers -----------------------------------------------------------------

    function onSelectionChange() {
      if (document.activeElement !== input) {
        batch(() => setState({ mirrorSelectionStart: null, mirrorSelectionEnd: null }))
        return
      }
      const selStart = input.selectionStart
      const selEnd = input.selectionEnd
      const selDir = input.selectionDirection
      const max = input.maxLength
      const val = input.value
      const prev = prevSelection
      let start = -1
      let end = -1
      let direction
      if (val.length !== 0 && selStart !== null && selEnd !== null) {
        const isSingleCaret = selStart === selEnd
        const isInsertMode = selStart === val.length && val.length < max
        if (isSingleCaret && !isInsertMode) {
          const c = selStart
          if (c === 0) {
            start = 0
            end = 1
            direction = "forward"
          } else if (c === max) {
            start = c - 1
            end = c
            direction = "backward"
          } else if (max > 1 && val.length > 1) {
            let offset = 0
            if (prev[0] !== null && prev[1] !== null) {
              direction = c < prev[1] ? "backward" : "forward"
              const wasPreviouslyInserting = prev[0] === prev[1] && prev[0] < max
              if (direction === "backward" && !wasPreviouslyInserting) offset = -1
            }
            start = offset + c
            end = offset + c + 1
          }
        }
        if (start !== -1 && end !== -1 && start !== end) input.setSelectionRange(start, end, direction)
      }
      const s = start !== -1 ? start : selStart
      const e = end !== -1 ? end : selEnd
      const dir = direction ?? selDir
      batch(() => setState({ mirrorSelectionStart: s, mirrorSelectionEnd: e }))
      prevSelection = [s, e, dir]
    }

    function onInput(event) {
      if (input.value === trackedValue) return
      batch(() => {
        const newValue = input.value.slice(0, maxLength)
        if (newValue.length > 0 && regexp && !regexp.test(newValue)) {
          event.preventDefault()
          return
        }
        if (typeof previousValue === "string" && newValue.length < previousValue.length) {
          document.dispatchEvent(new Event("selectionchange"))
        }
        setValue(newValue)
      })
      // React restores the controlled value after the event.
      if (input.value !== state.value) input.value = state.value
      trackedValue = state.value
    }

    function onFocus() {
      batch(() => {
        const start = Math.min(input.value.length, maxLength - 1)
        const end = input.value.length
        input.setSelectionRange(start, end)
        setState({ mirrorSelectionStart: start, mirrorSelectionEnd: end, isFocused: true })
      })
    }

    function onBlur() {
      setState({ isFocused: false })
    }

    function onPaste(event) {
      if (!isIOS || !event.clipboardData) return
      const content = event.clipboardData.getData("text/plain")
      event.preventDefault()
      const start = input.selectionStart
      const end = input.selectionEnd
      const v = state.value
      const newValue = (start !== end ? v.slice(0, start) + content + v.slice(end) : v.slice(0, start) + content + v.slice(start)).slice(0, maxLength)
      if (newValue.length > 0 && regexp && !regexp.test(newValue)) return
      batch(() => {
        input.value = newValue
        setValue(newValue)
        const s = Math.min(newValue.length, maxLength - 1)
        const e = newValue.length
        input.setSelectionRange(s, e)
        setState({ mirrorSelectionStart: s, mirrorSelectionEnd: e })
      })
    }

    const onMouseOver = () => setStateLater({ isHovering: true })
    const onMouseLeave = () => setStateLater({ isHovering: false })

    // --- mount --------------------------------------------------------------------

    input.addEventListener("input", onInput)
    input.addEventListener("focus", onFocus)
    input.addEventListener("blur", onBlur)
    input.addEventListener("paste", onPaste)
    input.addEventListener("mouseover", onMouseOver)
    input.addEventListener("mouseleave", onMouseLeave)

    // Effects run in declaration order: usePrevious, the mount effect, the
    // [value, isFocused] resync, then the password manager hook.
    batch(() => {
      if (state.value !== input.value) setValue(input.value)
      prevSelection = [input.selectionStart, input.selectionEnd, input.selectionDirection]
      document.addEventListener("selectionchange", onSelectionChange, { capture: true })
      onSelectionChange()
      if (document.activeElement === input) setState({ isFocused: true })
      injectStyle()
    })
    const updateRootHeight = () => container.style.setProperty("--root-height", `${input.clientHeight}px`)
    updateRootHeight()
    const resizeObserver = new ResizeObserver(updateRootHeight)
    resizeObserver.observe(input)

    runEffects()
    trackPWMBadgeSpace()
    const spaceInterval = setInterval(trackPWMBadgeSpace, 1000)

    return () => {
      document.removeEventListener("selectionchange", onSelectionChange, { capture: true })
      input.removeEventListener("input", onInput)
      input.removeEventListener("focus", onFocus)
      input.removeEventListener("blur", onBlur)
      input.removeEventListener("paste", onPaste)
      input.removeEventListener("mouseover", onMouseOver)
      input.removeEventListener("mouseleave", onMouseLeave)
      resizeObserver.disconnect()
      clearInterval(spaceInterval)
      syncTimers.forEach(clearTimeout)
      pwmTimers.forEach(clearTimeout)
    }
  },
})
