// sonner.js: sonner 2.0.7 (the toast library behind shadcn/ui's Sonner)
// for the server-rendered <Toaster /> of components/ui/sonner.
//
//   import { toast } from "/_shadcn/sonner.js"
//   toast("Event has been created", { description: "...", action: { label: "Undo", onClick() {} } })
//   toast.success("Saved")   // .info .warning .error .loading .message
//   toast.promise(fetch("/api"), { loading: "Saving...", success: (r) => "Saved", error: "Failed" })
//   toast.dismiss(id)        // or toast.dismiss() for all
//
// Titles, descriptions and action labels are strings or DOM nodes.
//
// The toaster renders and behaves like sonner's React components: the same
// markup and data attributes, the same stack layout (heights, offsets,
// front toast height, lift, scale), hover / focus / Alt+T expansion,
// pausing timers while expanded or interacting, swipe to dismiss in the
// position's directions, enter and exit transitions, promise toasts, and
// sonner's stylesheet (injected verbatim with the page's CSP nonce). Since
// those transitions start from whatever state the DOM was in when styles
// were last computed, state changes are committed in the same order and at
// the same moments as React does (synchronously for flushSync and discrete
// events, in a later task for timers and continuous events).
//
// Markup contract (rendered by the Go package): the <section> of the
// toaster contains a <template data-sonner-toaster ...> carrying the
// Toaster props as data attributes and the icons as children marked
// data-sonner-icon="success|info|warning|error|loading|close".

import { define } from "./core.js"
import { injectStyle } from "./style.js"

// sonner's injected stylesheet, verbatim.
const SONNER_CSS = "[data-sonner-toaster][dir=ltr],html[dir=ltr]{--toast-icon-margin-start:-3px;--toast-icon-margin-end:4px;--toast-svg-margin-start:-1px;--toast-svg-margin-end:0px;--toast-button-margin-start:auto;--toast-button-margin-end:0;--toast-close-button-start:0;--toast-close-button-end:unset;--toast-close-button-transform:translate(-35%, -35%)}[data-sonner-toaster][dir=rtl],html[dir=rtl]{--toast-icon-margin-start:4px;--toast-icon-margin-end:-3px;--toast-svg-margin-start:0px;--toast-svg-margin-end:-1px;--toast-button-margin-start:0;--toast-button-margin-end:auto;--toast-close-button-start:unset;--toast-close-button-end:0;--toast-close-button-transform:translate(35%, -35%)}[data-sonner-toaster]{position:fixed;width:var(--width);font-family:ui-sans-serif,system-ui,-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica Neue,Arial,Noto Sans,sans-serif,Apple Color Emoji,Segoe UI Emoji,Segoe UI Symbol,Noto Color Emoji;--gray1:hsl(0, 0%, 99%);--gray2:hsl(0, 0%, 97.3%);--gray3:hsl(0, 0%, 95.1%);--gray4:hsl(0, 0%, 93%);--gray5:hsl(0, 0%, 90.9%);--gray6:hsl(0, 0%, 88.7%);--gray7:hsl(0, 0%, 85.8%);--gray8:hsl(0, 0%, 78%);--gray9:hsl(0, 0%, 56.1%);--gray10:hsl(0, 0%, 52.3%);--gray11:hsl(0, 0%, 43.5%);--gray12:hsl(0, 0%, 9%);--border-radius:8px;box-sizing:border-box;padding:0;margin:0;list-style:none;outline:0;z-index:999999999;transition:transform .4s ease}@media (hover:none) and (pointer:coarse){[data-sonner-toaster][data-lifted=true]{transform:none}}[data-sonner-toaster][data-x-position=right]{right:var(--offset-right)}[data-sonner-toaster][data-x-position=left]{left:var(--offset-left)}[data-sonner-toaster][data-x-position=center]{left:50%;transform:translateX(-50%)}[data-sonner-toaster][data-y-position=top]{top:var(--offset-top)}[data-sonner-toaster][data-y-position=bottom]{bottom:var(--offset-bottom)}[data-sonner-toast]{--y:translateY(100%);--lift-amount:calc(var(--lift) * var(--gap));z-index:var(--z-index);position:absolute;opacity:0;transform:var(--y);touch-action:none;transition:transform .4s,opacity .4s,height .4s,box-shadow .2s;box-sizing:border-box;outline:0;overflow-wrap:anywhere}[data-sonner-toast][data-styled=true]{padding:16px;background:var(--normal-bg);border:1px solid var(--normal-border);color:var(--normal-text);border-radius:var(--border-radius);box-shadow:0 4px 12px rgba(0,0,0,.1);width:var(--width);font-size:13px;display:flex;align-items:center;gap:6px}[data-sonner-toast]:focus-visible{box-shadow:0 4px 12px rgba(0,0,0,.1),0 0 0 2px rgba(0,0,0,.2)}[data-sonner-toast][data-y-position=top]{top:0;--y:translateY(-100%);--lift:1;--lift-amount:calc(1 * var(--gap))}[data-sonner-toast][data-y-position=bottom]{bottom:0;--y:translateY(100%);--lift:-1;--lift-amount:calc(var(--lift) * var(--gap))}[data-sonner-toast][data-styled=true] [data-description]{font-weight:400;line-height:1.4;color:#3f3f3f}[data-rich-colors=true][data-sonner-toast][data-styled=true] [data-description]{color:inherit}[data-sonner-toaster][data-sonner-theme=dark] [data-description]{color:#e8e8e8}[data-sonner-toast][data-styled=true] [data-title]{font-weight:500;line-height:1.5;color:inherit}[data-sonner-toast][data-styled=true] [data-icon]{display:flex;height:16px;width:16px;position:relative;justify-content:flex-start;align-items:center;flex-shrink:0;margin-left:var(--toast-icon-margin-start);margin-right:var(--toast-icon-margin-end)}[data-sonner-toast][data-promise=true] [data-icon]>svg{opacity:0;transform:scale(.8);transform-origin:center;animation:sonner-fade-in .3s ease forwards}[data-sonner-toast][data-styled=true] [data-icon]>*{flex-shrink:0}[data-sonner-toast][data-styled=true] [data-icon] svg{margin-left:var(--toast-svg-margin-start);margin-right:var(--toast-svg-margin-end)}[data-sonner-toast][data-styled=true] [data-content]{display:flex;flex-direction:column;gap:2px}[data-sonner-toast][data-styled=true] [data-button]{border-radius:4px;padding-left:8px;padding-right:8px;height:24px;font-size:12px;color:var(--normal-bg);background:var(--normal-text);margin-left:var(--toast-button-margin-start);margin-right:var(--toast-button-margin-end);border:none;font-weight:500;cursor:pointer;outline:0;display:flex;align-items:center;flex-shrink:0;transition:opacity .4s,box-shadow .2s}[data-sonner-toast][data-styled=true] [data-button]:focus-visible{box-shadow:0 0 0 2px rgba(0,0,0,.4)}[data-sonner-toast][data-styled=true] [data-button]:first-of-type{margin-left:var(--toast-button-margin-start);margin-right:var(--toast-button-margin-end)}[data-sonner-toast][data-styled=true] [data-cancel]{color:var(--normal-text);background:rgba(0,0,0,.08)}[data-sonner-toaster][data-sonner-theme=dark] [data-sonner-toast][data-styled=true] [data-cancel]{background:rgba(255,255,255,.3)}[data-sonner-toast][data-styled=true] [data-close-button]{position:absolute;left:var(--toast-close-button-start);right:var(--toast-close-button-end);top:0;height:20px;width:20px;display:flex;justify-content:center;align-items:center;padding:0;color:var(--gray12);background:var(--normal-bg);border:1px solid var(--gray4);transform:var(--toast-close-button-transform);border-radius:50%;cursor:pointer;z-index:1;transition:opacity .1s,background .2s,border-color .2s}[data-sonner-toast][data-styled=true] [data-close-button]:focus-visible{box-shadow:0 4px 12px rgba(0,0,0,.1),0 0 0 2px rgba(0,0,0,.2)}[data-sonner-toast][data-styled=true] [data-disabled=true]{cursor:not-allowed}[data-sonner-toast][data-styled=true]:hover [data-close-button]:hover{background:var(--gray2);border-color:var(--gray5)}[data-sonner-toast][data-swiping=true]::before{content:'';position:absolute;left:-100%;right:-100%;height:100%;z-index:-1}[data-sonner-toast][data-y-position=top][data-swiping=true]::before{bottom:50%;transform:scaleY(3) translateY(50%)}[data-sonner-toast][data-y-position=bottom][data-swiping=true]::before{top:50%;transform:scaleY(3) translateY(-50%)}[data-sonner-toast][data-swiping=false][data-removed=true]::before{content:'';position:absolute;inset:0;transform:scaleY(2)}[data-sonner-toast][data-expanded=true]::after{content:'';position:absolute;left:0;height:calc(var(--gap) + 1px);bottom:100%;width:100%}[data-sonner-toast][data-mounted=true]{--y:translateY(0);opacity:1}[data-sonner-toast][data-expanded=false][data-front=false]{--scale:var(--toasts-before) * 0.05 + 1;--y:translateY(calc(var(--lift-amount) * var(--toasts-before))) scale(calc(-1 * var(--scale)));height:var(--front-toast-height)}[data-sonner-toast]>*{transition:opacity .4s}[data-sonner-toast][data-x-position=right]{right:0}[data-sonner-toast][data-x-position=left]{left:0}[data-sonner-toast][data-expanded=false][data-front=false][data-styled=true]>*{opacity:0}[data-sonner-toast][data-visible=false]{opacity:0;pointer-events:none}[data-sonner-toast][data-mounted=true][data-expanded=true]{--y:translateY(calc(var(--lift) * var(--offset)));height:var(--initial-height)}[data-sonner-toast][data-removed=true][data-front=true][data-swipe-out=false]{--y:translateY(calc(var(--lift) * -100%));opacity:0}[data-sonner-toast][data-removed=true][data-front=false][data-swipe-out=false][data-expanded=true]{--y:translateY(calc(var(--lift) * var(--offset) + var(--lift) * -100%));opacity:0}[data-sonner-toast][data-removed=true][data-front=false][data-swipe-out=false][data-expanded=false]{--y:translateY(40%);opacity:0;transition:transform .5s,opacity .2s}[data-sonner-toast][data-removed=true][data-front=false]::before{height:calc(var(--initial-height) + 20%)}[data-sonner-toast][data-swiping=true]{transform:var(--y) translateY(var(--swipe-amount-y,0)) translateX(var(--swipe-amount-x,0));transition:none}[data-sonner-toast][data-swiped=true]{user-select:none}[data-sonner-toast][data-swipe-out=true][data-y-position=bottom],[data-sonner-toast][data-swipe-out=true][data-y-position=top]{animation-duration:.2s;animation-timing-function:ease-out;animation-fill-mode:forwards}[data-sonner-toast][data-swipe-out=true][data-swipe-direction=left]{animation-name:swipe-out-left}[data-sonner-toast][data-swipe-out=true][data-swipe-direction=right]{animation-name:swipe-out-right}[data-sonner-toast][data-swipe-out=true][data-swipe-direction=up]{animation-name:swipe-out-up}[data-sonner-toast][data-swipe-out=true][data-swipe-direction=down]{animation-name:swipe-out-down}@keyframes swipe-out-left{from{transform:var(--y) translateX(var(--swipe-amount-x));opacity:1}to{transform:var(--y) translateX(calc(var(--swipe-amount-x) - 100%));opacity:0}}@keyframes swipe-out-right{from{transform:var(--y) translateX(var(--swipe-amount-x));opacity:1}to{transform:var(--y) translateX(calc(var(--swipe-amount-x) + 100%));opacity:0}}@keyframes swipe-out-up{from{transform:var(--y) translateY(var(--swipe-amount-y));opacity:1}to{transform:var(--y) translateY(calc(var(--swipe-amount-y) - 100%));opacity:0}}@keyframes swipe-out-down{from{transform:var(--y) translateY(var(--swipe-amount-y));opacity:1}to{transform:var(--y) translateY(calc(var(--swipe-amount-y) + 100%));opacity:0}}@media (max-width:600px){[data-sonner-toaster]{position:fixed;right:var(--mobile-offset-right);left:var(--mobile-offset-left);width:100%}[data-sonner-toaster][dir=rtl]{left:calc(var(--mobile-offset-left) * -1)}[data-sonner-toaster] [data-sonner-toast]{left:0;right:0;width:calc(100% - var(--mobile-offset-left) * 2)}[data-sonner-toaster][data-x-position=left]{left:var(--mobile-offset-left)}[data-sonner-toaster][data-y-position=bottom]{bottom:var(--mobile-offset-bottom)}[data-sonner-toaster][data-y-position=top]{top:var(--mobile-offset-top)}[data-sonner-toaster][data-x-position=center]{left:var(--mobile-offset-left);right:var(--mobile-offset-right);transform:none}}[data-sonner-toaster][data-sonner-theme=light]{--normal-bg:#fff;--normal-border:var(--gray4);--normal-text:var(--gray12);--success-bg:hsl(143, 85%, 96%);--success-border:hsl(145, 92%, 87%);--success-text:hsl(140, 100%, 27%);--info-bg:hsl(208, 100%, 97%);--info-border:hsl(221, 91%, 93%);--info-text:hsl(210, 92%, 45%);--warning-bg:hsl(49, 100%, 97%);--warning-border:hsl(49, 91%, 84%);--warning-text:hsl(31, 92%, 45%);--error-bg:hsl(359, 100%, 97%);--error-border:hsl(359, 100%, 94%);--error-text:hsl(360, 100%, 45%)}[data-sonner-toaster][data-sonner-theme=light] [data-sonner-toast][data-invert=true]{--normal-bg:#000;--normal-border:hsl(0, 0%, 20%);--normal-text:var(--gray1)}[data-sonner-toaster][data-sonner-theme=dark] [data-sonner-toast][data-invert=true]{--normal-bg:#fff;--normal-border:var(--gray3);--normal-text:var(--gray12)}[data-sonner-toaster][data-sonner-theme=dark]{--normal-bg:#000;--normal-bg-hover:hsl(0, 0%, 12%);--normal-border:hsl(0, 0%, 20%);--normal-border-hover:hsl(0, 0%, 25%);--normal-text:var(--gray1);--success-bg:hsl(150, 100%, 6%);--success-border:hsl(147, 100%, 12%);--success-text:hsl(150, 86%, 65%);--info-bg:hsl(215, 100%, 6%);--info-border:hsl(223, 43%, 17%);--info-text:hsl(216, 87%, 65%);--warning-bg:hsl(64, 100%, 6%);--warning-border:hsl(60, 100%, 9%);--warning-text:hsl(46, 87%, 65%);--error-bg:hsl(358, 76%, 10%);--error-border:hsl(357, 89%, 16%);--error-text:hsl(358, 100%, 81%)}[data-sonner-toaster][data-sonner-theme=dark] [data-sonner-toast] [data-close-button]{background:var(--normal-bg);border-color:var(--normal-border);color:var(--normal-text)}[data-sonner-toaster][data-sonner-theme=dark] [data-sonner-toast] [data-close-button]:hover{background:var(--normal-bg-hover);border-color:var(--normal-border-hover)}[data-rich-colors=true][data-sonner-toast][data-type=success]{background:var(--success-bg);border-color:var(--success-border);color:var(--success-text)}[data-rich-colors=true][data-sonner-toast][data-type=success] [data-close-button]{background:var(--success-bg);border-color:var(--success-border);color:var(--success-text)}[data-rich-colors=true][data-sonner-toast][data-type=info]{background:var(--info-bg);border-color:var(--info-border);color:var(--info-text)}[data-rich-colors=true][data-sonner-toast][data-type=info] [data-close-button]{background:var(--info-bg);border-color:var(--info-border);color:var(--info-text)}[data-rich-colors=true][data-sonner-toast][data-type=warning]{background:var(--warning-bg);border-color:var(--warning-border);color:var(--warning-text)}[data-rich-colors=true][data-sonner-toast][data-type=warning] [data-close-button]{background:var(--warning-bg);border-color:var(--warning-border);color:var(--warning-text)}[data-rich-colors=true][data-sonner-toast][data-type=error]{background:var(--error-bg);border-color:var(--error-border);color:var(--error-text)}[data-rich-colors=true][data-sonner-toast][data-type=error] [data-close-button]{background:var(--error-bg);border-color:var(--error-border);color:var(--error-text)}.sonner-loading-wrapper{--size:16px;height:var(--size);width:var(--size);position:absolute;inset:0;z-index:10}.sonner-loading-wrapper[data-visible=false]{transform-origin:center;animation:sonner-fade-out .2s ease forwards}.sonner-spinner{position:relative;top:50%;left:50%;height:var(--size);width:var(--size)}.sonner-loading-bar{animation:sonner-spin 1.2s linear infinite;background:var(--gray11);border-radius:6px;height:8%;left:-10%;position:absolute;top:-3.9%;width:24%}.sonner-loading-bar:first-child{animation-delay:-1.2s;transform:rotate(.0001deg) translate(146%)}.sonner-loading-bar:nth-child(2){animation-delay:-1.1s;transform:rotate(30deg) translate(146%)}.sonner-loading-bar:nth-child(3){animation-delay:-1s;transform:rotate(60deg) translate(146%)}.sonner-loading-bar:nth-child(4){animation-delay:-.9s;transform:rotate(90deg) translate(146%)}.sonner-loading-bar:nth-child(5){animation-delay:-.8s;transform:rotate(120deg) translate(146%)}.sonner-loading-bar:nth-child(6){animation-delay:-.7s;transform:rotate(150deg) translate(146%)}.sonner-loading-bar:nth-child(7){animation-delay:-.6s;transform:rotate(180deg) translate(146%)}.sonner-loading-bar:nth-child(8){animation-delay:-.5s;transform:rotate(210deg) translate(146%)}.sonner-loading-bar:nth-child(9){animation-delay:-.4s;transform:rotate(240deg) translate(146%)}.sonner-loading-bar:nth-child(10){animation-delay:-.3s;transform:rotate(270deg) translate(146%)}.sonner-loading-bar:nth-child(11){animation-delay:-.2s;transform:rotate(300deg) translate(146%)}.sonner-loading-bar:nth-child(12){animation-delay:-.1s;transform:rotate(330deg) translate(146%)}@keyframes sonner-fade-in{0%{opacity:0;transform:scale(.8)}100%{opacity:1;transform:scale(1)}}@keyframes sonner-fade-out{0%{opacity:1;transform:scale(1)}100%{opacity:0;transform:scale(.8)}}@keyframes sonner-spin{0%{opacity:1}100%{opacity:.15}}@media (prefers-reduced-motion){.sonner-loading-bar,[data-sonner-toast],[data-sonner-toast]>*{transition:none!important;animation:none!important}}.sonner-loader{position:absolute;top:50%;left:50%;transform:translate(-50%,-50%);transform-origin:center;transition:opacity .2s,transform .2s}.sonner-loader[data-visible=false]{opacity:0;transform:scale(.8) translate(-50%,-50%)}"

let cssInjected = false
function ensureCSS() {
  if (cssInjected) return
  cssInjected = true
  injectStyle(SONNER_CSS)
}

const VISIBLE_TOASTS_AMOUNT = 3
const VIEWPORT_OFFSET = "24px"
const MOBILE_VIEWPORT_OFFSET = "16px"
const TOAST_LIFETIME = 4000
const TOAST_WIDTH = 356
const GAP = 14
const SWIPE_THRESHOLD = 45
const TIME_BEFORE_UNMOUNT = 200

const SVG_NS = "http://www.w3.org/2000/svg"

// --- React-like scheduling ------------------------------------------------------------------------
//
// Every toaster commits its state like a React root: updates are queued
// and rendered together. Their lane decides when: "sync" (flushSync, and
// updates made while a discrete event such as click or keydown is being
// dispatched) renders in a microtask; "default" (timers, animation frames,
// continuous events such as mousemove) renders in a later task.

const DISCRETE = new Set("cancel click close contextmenu copy cut auxclick dblclick dragend dragstart drop focusin focusout input invalid keydown keypress keyup mousedown mouseup paste pause play pointercancel pointerdown pointerup ratechange reset resize seeked submit toggle touchcancel touchend touchstart volumechange change selectionchange textInput compositionstart compositionend compositionupdate beforeblur afterblur beforeinput blur fullscreenchange focus hashchange popstate select selectstart visibilitychange".split(" "))

let laneOverride = null
const currentLane = () => {
  if (laneOverride) return laneOverride
  const type = window.event?.type
  return type && DISCRETE.has(type) ? "sync" : "default"
}
const withLane = (lane, fn) => {
  const prev = laneOverride
  laneOverride = lane
  try {
    return fn()
  } finally {
    laneOverride = prev
  }
}

const channel = new MessageChannel()
const taskQueue = []
channel.port1.onmessage = () => {
  const fn = taskQueue.shift()
  if (fn) withLane("default", fn)
}
function postTask(fn) {
  taskQueue.push(fn)
  channel.port2.postMessage(null)
}

// --- the observer (ToastState) ---------------------------------------------------------------------

let toastsCounter = 1

const isHttpResponse = (data) =>
  data && typeof data === "object" && "ok" in data && typeof data.ok === "boolean" && "status" in data && typeof data.status === "number"
const isNode = (v) => v instanceof Node

class Observer {
  constructor() {
    this.subscribers = []
    this.toasts = []
    this.dismissedToasts = new Set()
  }
  subscribe = (subscriber) => {
    this.subscribers.push(subscriber)
    return () => {
      const index = this.subscribers.indexOf(subscriber)
      this.subscribers.splice(index, 1)
    }
  }
  publish = (data) => {
    this.subscribers.forEach((subscriber) => subscriber(data))
  }
  addToast = (data) => {
    this.publish(data)
    this.toasts = [...this.toasts, data]
  }
  create = (data) => {
    const { message, ...rest } = data
    const id = typeof data?.id === "number" || data.id?.length > 0 ? data.id : toastsCounter++
    const alreadyExists = this.toasts.find((toast) => toast.id === id)
    const dismissible = data.dismissible === undefined ? true : data.dismissible
    if (this.dismissedToasts.has(id)) this.dismissedToasts.delete(id)
    if (alreadyExists) {
      this.toasts = this.toasts.map((toast) => {
        if (toast.id === id) {
          this.publish({ ...toast, ...data, id, title: message })
          return { ...toast, ...data, id, dismissible, title: message }
        }
        return toast
      })
    } else {
      this.addToast({ title: message, ...rest, dismissible, id })
    }
    return id
  }
  dismiss = (id) => {
    if (id) {
      this.dismissedToasts.add(id)
      requestAnimationFrame(() => this.subscribers.forEach((subscriber) => subscriber({ id, dismiss: true })))
    } else {
      this.toasts.forEach((toast) => {
        this.subscribers.forEach((subscriber) => subscriber({ id: toast.id, dismiss: true }))
      })
    }
    return id
  }
  message = (message, data) => this.create({ ...data, message })
  error = (message, data) => this.create({ ...data, message, type: "error" })
  success = (message, data) => this.create({ ...data, type: "success", message })
  info = (message, data) => this.create({ ...data, type: "info", message })
  warning = (message, data) => this.create({ ...data, type: "warning", message })
  loading = (message, data) => this.create({ ...data, type: "loading", message })
  promise = (promise, data) => {
    if (!data) return
    let id = undefined
    if (data.loading !== undefined) {
      id = this.create({
        ...data,
        promise,
        type: "loading",
        message: data.loading,
        description: typeof data.description !== "function" ? data.description : undefined,
      })
    }
    const p = Promise.resolve(promise instanceof Function ? promise() : promise)
    let shouldDismiss = id !== undefined
    let result
    const settle = async (type, value, pick) => {
      shouldDismiss = false
      const promiseData = typeof pick === "function" ? await pick(value) : pick
      const description = typeof data.description === "function" ? await data.description(value) : data.description
      const isExtendedResult = typeof promiseData === "object" && promiseData !== null && !isNode(promiseData)
      const toastSettings = isExtendedResult ? promiseData : { message: promiseData }
      this.create({ id, type, description, ...toastSettings })
    }
    const originalPromise = p
      .then(async (response) => {
        result = ["resolve", response]
        if (isNode(response)) {
          shouldDismiss = false
          this.create({ id, type: "default", message: response })
        } else if (isHttpResponse(response) && !response.ok) {
          await settle("error", `HTTP error! status: ${response.status}`, data.error)
        } else if (response instanceof Error) {
          await settle("error", response, data.error)
        } else if (data.success !== undefined) {
          await settle("success", response, data.success)
        }
      })
      .catch(async (error) => {
        result = ["reject", error]
        if (data.error !== undefined) await settle("error", error, data.error)
      })
      .finally(() => {
        if (shouldDismiss) {
          this.dismiss(id)
          id = undefined
        }
        data.finally?.call(data)
      })
    const unwrap = () =>
      new Promise((resolve, reject) =>
        originalPromise.then(() => (result[0] === "reject" ? reject(result[1]) : resolve(result[1]))).catch(reject),
      )
    if (typeof id !== "string" && typeof id !== "number") return { unwrap }
    return Object.assign(id, { unwrap })
  }
  custom = (jsx, data) => {
    const id = data?.id || toastsCounter++
    this.create({ jsx: jsx(id), id, ...data })
    return id
  }
  getActiveToasts = () => this.toasts.filter((toast) => !this.dismissedToasts.has(toast.id))
}

const ToastState = new Observer()

const toastFunction = (message, data) => {
  const id = data?.id || toastsCounter++
  ToastState.addToast({ title: message, ...data, id })
  return id
}

/** sonner's toast(): see the file header. */
export const toast = Object.assign(
  toastFunction,
  {
    success: ToastState.success,
    info: ToastState.info,
    warning: ToastState.warning,
    error: ToastState.error,
    custom: ToastState.custom,
    message: ToastState.message,
    promise: ToastState.promise,
    dismiss: ToastState.dismiss,
    loading: ToastState.loading,
  },
  { getHistory: () => ToastState.toasts, getToasts: () => ToastState.getActiveToasts() },
)

// --- helpers --------------------------------------------------------------------------------------

const cn = (...classes) => classes.filter(Boolean).join(" ")

function getDefaultSwipeDirections(position) {
  const [y, x] = position.split("-")
  const directions = []
  if (y) directions.push(y)
  if (x) directions.push(x)
  return directions
}

function getDocumentDirection() {
  const dirAttribute = document.documentElement.getAttribute("dir")
  if (dirAttribute === "auto" || !dirAttribute) return window.getComputedStyle(document.documentElement).direction
  return dirAttribute
}

function assignOffset(defaultOffset, mobileOffset) {
  const styles = {}
  ;[defaultOffset, mobileOffset].forEach((offset, index) => {
    const isMobile = index === 1
    const prefix = isMobile ? "--mobile-offset" : "--offset"
    const defaultValue = isMobile ? MOBILE_VIEWPORT_OFFSET : VIEWPORT_OFFSET
    const assignAll = (offset) => {
      ;["top", "right", "bottom", "left"].forEach((key) => {
        styles[`${prefix}-${key}`] = typeof offset === "number" ? `${offset}px` : offset
      })
    }
    if (typeof offset === "number" || typeof offset === "string") assignAll(offset)
    else if (typeof offset === "object" && offset !== null) {
      ;["top", "right", "bottom", "left"].forEach((key) => {
        if (offset[key] === undefined) styles[`${prefix}-${key}`] = defaultValue
        else styles[`${prefix}-${key}`] = typeof offset[key] === "number" ? `${offset[key]}px` : offset[key]
      })
    } else assignAll(defaultValue)
  })
  return styles
}

/** React's rendering of a data-* attribute value: null/undefined removes it. */
function setAttr(el, name, value) {
  if (value === undefined || value === null) {
    if (el.hasAttribute(name)) el.removeAttribute(name)
    return
  }
  const v = String(value)
  if (el.getAttribute(name) !== v) el.setAttribute(name, v)
}

/** React's style prop diffing: sets the given properties, removes the ones no longer given. */
function setStyles(el, next, prevRef) {
  const prev = prevRef.style ?? {}
  for (const key of Object.keys(prev)) {
    if (!(key in next) || next[key] === undefined || next[key] === null) {
      if (key.startsWith("--")) el.style.removeProperty(key)
      else el.style[key] = ""
    }
  }
  for (const [key, value] of Object.entries(next)) {
    if (value === undefined || value === null) continue
    if (prev[key] === value && key in prev) continue
    if (key.startsWith("--")) el.style.setProperty(key, String(value))
    else el.style[key] = typeof value === "number" && value !== 0 ? `${value}px` : String(value)
  }
  prevRef.style = { ...next }
}

/** Renders a title, description or label: a string, a node, or a function returning one. */
function content(value) {
  const v = typeof value === "function" ? value() : value
  if (v === undefined || v === null || v === false) return []
  if (isNode(v)) return [v.cloneNode(true)]
  return [document.createTextNode(String(v))]
}
function sameContent(a, b) {
  return a === b
}

// --- the Toaster ------------------------------------------------------------------------------------

function readConfig(tpl) {
  const a = (n) => tpl.getAttribute(n)
  const json = (n) => (tpl.hasAttribute(n) ? JSON.parse(a(n)) : undefined)
  const bool = (n) => (tpl.hasAttribute(n) ? a(n) !== "false" : undefined)
  const num = (n) => (tpl.hasAttribute(n) ? Number(a(n)) : undefined)
  const icons = {}
  for (const el of tpl.content.querySelectorAll("[data-sonner-icon]")) {
    icons[el.getAttribute("data-sonner-icon")] = el
  }
  return {
    id: a("data-id") || undefined,
    invert: bool("data-invert"),
    position: a("data-position") || "bottom-right",
    hotkey: json("data-hotkey") ?? ["altKey", "KeyT"],
    expand: bool("data-expand"),
    closeButton: bool("data-close-button"),
    className: a("data-class") ?? undefined,
    offset: json("data-offset"),
    mobileOffset: json("data-mobile-offset"),
    theme: a("data-theme") || "light",
    richColors: bool("data-rich-colors"),
    duration: num("data-duration"),
    style: json("data-style"),
    visibleToasts: num("data-visible-toasts") ?? VISIBLE_TOASTS_AMOUNT,
    toastOptions: json("data-toast-options"),
    dir: a("data-dir") || undefined,
    gap: num("data-gap") ?? GAP,
    icons,
    containerAriaLabel: a("data-container-aria-label") || "Notifications",
    swipeDirections: json("data-swipe-directions"),
  }
}

function iconNode(icons, type) {
  const el = icons[type]
  if (!el) return null
  const c = el.cloneNode(true)
  c.removeAttribute("data-sonner-icon")
  return c
}

// sonner's built-in icons (used when the Toaster has none for a type).
const ASSETS = {
  success: ["0 0 20 20", "M10 18a8 8 0 100-16 8 8 0 000 16zm3.857-9.809a.75.75 0 00-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 10-1.06 1.061l2.5 2.5a.75.75 0 001.137-.089l4-5.5z"],
  warning: ["0 0 24 24", "M9.401 3.003c1.155-2 4.043-2 5.197 0l7.355 12.748c1.154 2-.29 4.5-2.599 4.5H4.645c-2.309 0-3.752-2.5-2.598-4.5L9.4 3.003zM12 8.25a.75.75 0 01.75.75v3.75a.75.75 0 01-1.5 0V9a.75.75 0 01.75-.75zm0 8.25a.75.75 0 100-1.5.75.75 0 000 1.5z"],
  info: ["0 0 20 20", "M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a.75.75 0 000 1.5h.253a.25.25 0 01.244.304l-.459 2.066A1.75 1.75 0 0010.747 15H11a.75.75 0 000-1.5h-.253a.25.25 0 01-.244-.304l.459-2.066A1.75 1.75 0 009.253 9H9z"],
  error: ["0 0 20 20", "M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-8-5a.75.75 0 01.75.75v4.5a.75.75 0 01-1.5 0v-4.5A.75.75 0 0110 5zm0 10a1 1 0 100-2 1 1 0 000 2z"],
}
function getAsset(type) {
  const a = ASSETS[type]
  if (!a) return null
  const svg = document.createElementNS(SVG_NS, "svg")
  for (const [k, v] of [["xmlns", SVG_NS], ["viewBox", a[0]], ["fill", "currentColor"], ["height", "20"], ["width", "20"]]) svg.setAttribute(k, v)
  const path = document.createElementNS(SVG_NS, "path")
  path.setAttribute("fill-rule", "evenodd")
  path.setAttribute("d", a[1])
  path.setAttribute("clip-rule", "evenodd")
  svg.appendChild(path)
  return svg
}
function closeIcon() {
  const svg = document.createElementNS(SVG_NS, "svg")
  for (const [k, v] of [["xmlns", SVG_NS], ["width", "12"], ["height", "12"], ["viewBox", "0 0 24 24"], ["fill", "none"], ["stroke", "currentColor"], ["stroke-width", "1.5"], ["stroke-linecap", "round"], ["stroke-linejoin", "round"]]) svg.setAttribute(k, v)
  for (const [x1, y1, x2, y2] of [["18", "6", "6", "18"], ["6", "6", "18", "18"]]) {
    const line = document.createElementNS(SVG_NS, "line")
    line.setAttribute("x1", x1)
    line.setAttribute("y1", y1)
    line.setAttribute("x2", x2)
    line.setAttribute("y2", y2)
    svg.appendChild(line)
  }
  return svg
}
function loaderBars(visible, className) {
  const wrap = document.createElement("div")
  wrap.className = ["sonner-loading-wrapper", className].filter(Boolean).join(" ")
  wrap.setAttribute("data-visible", String(visible))
  const spinner = document.createElement("div")
  spinner.className = "sonner-spinner"
  for (let i = 0; i < 12; i++) {
    const bar = document.createElement("div")
    bar.className = "sonner-loading-bar"
    spinner.appendChild(bar)
  }
  wrap.appendChild(spinner)
  return wrap
}

class Toaster {
  constructor(section, tpl) {
    ensureCSS()
    this.section = section
    this.tpl = tpl
    this.props = readConfig(tpl)
    const theme = this.props.theme
    this.state = {
      toasts: [],
      heights: [],
      expanded: false,
      interacting: false,
      actualTheme: theme !== "system" ? theme : window.matchMedia?.("(prefers-color-scheme: dark)").matches ? "dark" : "light",
    }
    this.queue = [] // pending updates
    this.syncScheduled = false
    this.defaultScheduled = false
    this.pendingPassive = null
    this.comps = new Map() // toast id -> Toast
    this.lists = new Map() // position -> { ol, styleRef }
    this.listRef = null
    this.effectListRef = null
    this.lastFocusedElement = null
    this.isFocusWithin = false
    this.isDocumentHidden = document.hidden
    this.prevToastsForExpand = this.state.toasts
    this.cleanups = []

    section.setAttribute("aria-label", `${this.props.containerAriaLabel} ${this.hotkeyLabel()}`)

    // ToastState subscription (the [toasts] effect re-subscribes on every
    // change; one subscription behaves the same).
    this.cleanups.push(
      ToastState.subscribe((t) => {
        if (t.dismiss) {
          requestAnimationFrame(() => {
            withLane("default", () => this.setToasts((toasts) => toasts.map((x) => (x.id === t.id ? { ...x, delete: true } : x))))
          })
          return
        }
        setTimeout(() => {
          this.flushSync(() => {
            this.setToasts((toasts) => {
              const i = toasts.findIndex((x) => x.id === t.id)
              if (i !== -1) return [...toasts.slice(0, i), { ...toasts[i], ...t }, ...toasts.slice(i + 1)]
              return [t, ...toasts]
            })
          })
        })
      }),
    )

    // Theme.
    if (theme === "system") {
      const mq = window.matchMedia("(prefers-color-scheme: dark)")
      const onChange = ({ matches }) => this.update(() => (this.state.actualTheme = matches ? "dark" : "light"))
      mq.addEventListener("change", onChange)
      this.cleanups.push(() => mq.removeEventListener("change", onChange))
    }

    // Hotkey and Escape.
    const onKeyDown = (event) => {
      const isHotkeyPressed = this.props.hotkey.every((key) => event[key] || event.code === key)
      if (isHotkeyPressed) {
        this.setExpanded(true)
        this.listRef?.focus()
      }
      if (event.code === "Escape" && (document.activeElement === this.listRef || this.listRef?.contains(document.activeElement))) {
        this.setExpanded(false)
      }
    }
    document.addEventListener("keydown", onKeyDown)
    this.cleanups.push(() => document.removeEventListener("keydown", onKeyDown))

    // useIsDocumentHidden (one per toast upstream; they all agree).
    const onVisibility = () => this.update(() => (this.isDocumentHidden = document.hidden))
    document.addEventListener("visibilitychange", onVisibility)
    this.cleanups.push(() => document.removeEventListener("visibilitychange", onVisibility))
  }

  hotkeyLabel() {
    return this.props.hotkey.join("+").replace(/Key/g, "").replace(/Digit/g, "")
  }

  // --- state updates ---

  update(fn, lane = currentLane()) {
    this.queue.push({ fn, lane })
    if (lane === "sync") {
      if (!this.syncScheduled) {
        this.syncScheduled = true
        queueMicrotask(() => {
          this.syncScheduled = false
          if (this.queue.some((u) => u.lane === "sync")) this.flush("sync")
        })
      }
    } else if (!this.defaultScheduled) {
      this.defaultScheduled = true
      postTask(() => {
        this.defaultScheduled = false
        if (this.queue.length) this.flush("default")
      })
    }
  }
  flushSync(fn) {
    withLane("sync", fn)
    this.flush("sync")
  }
  setToasts(fn) {
    this.update(() => (this.state.toasts = fn(this.state.toasts)))
  }
  setHeights(fn) {
    this.update(() => (this.state.heights = fn(this.state.heights)))
  }
  setExpanded(v) {
    this.update(() => (this.state.expanded = v))
  }
  setInteracting(v) {
    this.update(() => (this.state.interacting = v))
  }

  removeToast = (toastToRemove) => {
    this.setToasts((toasts) => {
      if (!toasts.find((t) => t.id === toastToRemove.id)?.delete) ToastState.dismiss(toastToRemove.id)
      return toasts.filter(({ id }) => id !== toastToRemove.id)
    })
  }

  // --- render ---

  // Renders the pending updates of a lane ("sync" renders only sync
  // updates; "default" renders everything pending), like React's
  // performSyncWorkOnRoot / performConcurrentWorkOnRoot.
  flush(lane) {
    if (this.pendingPassive) {
      const p = this.pendingPassive
      this.pendingPassive = null
      withLane("default", p)
    }
    const take = (l) => {
      if (l === "default") {
        const all = this.queue
        this.queue = []
        return all
      }
      const mine = this.queue.filter((u) => u.lane === "sync")
      this.queue = this.queue.filter((u) => u.lane !== "sync")
      return mine
    }
    let current = lane
    for (let guard = 0; guard < 50; guard++) {
      const updates = take(current)
      if (!updates.length) break
      for (const u of updates) u.fn()
      // ReactDOM turns its event system off while it mutates the DOM
      // (prepareForCommit): the blur Chromium fires when a focused node is
      // removed never reaches sonner's handlers.
      this.inCommit = true
      let effects
      try {
        effects = this.commit()
      } finally {
        this.inCommit = false
      }
      // Layout effects run during the commit; their updates are sync and
      // render before the browser paints (the loop below).
      withLane("sync", () => effects.layout.forEach((f) => f()))
      const passive = () => effects.passive.forEach((f) => f())
      if (current === "sync") withLane("default", passive)
      else {
        this.pendingPassive = passive
        postTask(() => {
          if (this.pendingPassive === passive) {
            this.pendingPassive = null
            passive()
          }
        })
      }
      if (!this.queue.some((u) => u.lane === "sync")) break
      current = "sync"
    }
  }

  commit() {
    const effects = { layout: [], passive: [] }
    const p = this.props
    const s = this.state
    const filteredToasts = p.id ? s.toasts.filter((t) => t.toasterId === p.id) : s.toasts.filter((t) => !t.toasterId)
    const possiblePositions = Array.from(new Set([p.position].concat(filteredToasts.filter((t) => t.position).map((t) => t.position))))
    const alive = new Set()
    // Elements of toasts that leave in this commit. React commits deletions
    // before placements, so they never count as a sibling to insert before:
    // moving a remaining toast's <li> would cancel its CSS transitions.
    const rendered = new Set(filteredToasts.map((t) => t.id))
    const leaving = new Set()
    for (const [id, comp] of this.comps) if (!rendered.has(id)) leaving.add(comp.el)

    // listRef is one object ref shared by every <ol>: React sets it when an
    // <ol> mounts (layout phase) and clears it when any <ol> unmounts
    // (mutation phase, unconditionally); it is not set again on updates.
    let attachedList = null
    possiblePositions.forEach((position, index) => {
      if (!filteredToasts.length) return
      const [y, x] = position.split("-")
      let list = this.lists.get(position)
      if (!list) {
        const ol = document.createElement("ol")
        list = { ol, styleRef: {} }
        this.lists.set(position, list)
        this.bindList(ol)
        this.section.appendChild(ol)
        attachedList = ol
      }
      const ol = list.ol
      setAttr(ol, "dir", p.dir === "auto" || !p.dir ? getDocumentDirection() : p.dir)
      setAttr(ol, "tabindex", "-1")
      if (p.className !== undefined) setAttr(ol, "class", p.className)
      setAttr(ol, "data-sonner-toaster", "true")
      setAttr(ol, "data-sonner-theme", s.actualTheme)
      setAttr(ol, "data-y-position", y)
      setAttr(ol, "data-x-position", x)
      setStyles(
        ol,
        {
          "--front-toast-height": `${s.heights[0]?.height || 0}px`,
          "--width": `${TOAST_WIDTH}px`,
          "--gap": `${p.gap}px`,
          ...p.style,
          ...assignOffset(p.offset, p.mobileOffset),
        },
        list.styleRef,
      )

      const toastsHere = filteredToasts.filter((t) => (!t.position && index === 0) || t.position === position)
      let prevEl = null
      toastsHere.forEach((t, i) => {
        alive.add(t.id)
        let comp = this.comps.get(t.id)
        if (!comp) {
          comp = new Toast(this, t)
          this.comps.set(t.id, comp)
        }
        const props = {
          toast: t,
          index: i,
          toasts: filteredToasts.filter((x) => x.position == t.position),
          heights: s.heights.filter((h) => h.position == t.position),
          expanded: s.expanded,
          interacting: s.interacting,
          position,
          isDocumentHidden: this.isDocumentHidden,
        }
        comp.render(props, effects)
        // React inserts new items before their next sibling.
        let ref = prevEl ? prevEl.nextSibling : ol.firstChild
        while (ref && leaving.has(ref)) ref = ref.nextSibling
        if (comp.el !== ref) ol.insertBefore(comp.el, ref)
        prevEl = comp.el
      })
    })

    // Unmount toasts no longer rendered.
    for (const [id, comp] of this.comps) {
      if (!alive.has(id)) {
        this.comps.delete(id)
        comp.el.remove()
        effects.passive.unshift(() => comp.unmount())
      }
    }
    // Lists no longer rendered.
    for (const [position, list] of this.lists) {
      if (!filteredToasts.length || !possiblePositions.includes(position)) {
        list.ol.remove()
        this.lists.delete(position)
        this.listRef = null
      }
    }
    if (attachedList) this.listRef = attachedList
    // The Toaster effect on [listRef.current]: its cleanup (registered while
    // listRef.current was set) gives focus back to where it was before the
    // toasts took it.
    if (this.listRef !== this.effectListRef) {
      const hadCleanup = !!this.effectListRef
      this.effectListRef = this.listRef
      if (hadCleanup) {
        effects.passive.unshift(() => {
          if (this.lastFocusedElement) {
            this.lastFocusedElement.focus({ preventScroll: true })
            this.lastFocusedElement = null
            this.isFocusWithin = false
          }
        })
      }
    }

    // Toaster effect: collapse when one toast or none is left.
    if (this.prevToastsForExpand !== s.toasts) {
      this.prevToastsForExpand = s.toasts
      const toasts = s.toasts
      effects.passive.push(() => {
        if (toasts.length <= 1) this.setExpanded(false)
      })
    }
    return effects
  }

  bindList(ol) {
    const discreteTarget = (event) => event.target instanceof HTMLElement && event.target.dataset.dismissible === "false"
    ol.addEventListener("focusout", (event) => {
      if (this.inCommit) return
      if (this.isFocusWithin && !event.currentTarget.contains(event.relatedTarget)) {
        this.isFocusWithin = false
        if (this.lastFocusedElement) {
          this.lastFocusedElement.focus({ preventScroll: true })
          this.lastFocusedElement = null
        }
      }
    })
    ol.addEventListener("focusin", (event) => {
      if (this.inCommit) return
      if (discreteTarget(event)) return
      if (!this.isFocusWithin) {
        this.isFocusWithin = true
        this.lastFocusedElement = event.relatedTarget
      }
    })
    ol.addEventListener("mouseenter", () => this.setExpanded(true))
    ol.addEventListener("mousemove", () => this.setExpanded(true))
    // React's onMouseLeave comes from mouseout events (a native mouseleave
    // also fires when a hovered toast is removed under the pointer, which
    // React never sees).
    ol.addEventListener("mouseout", (event) => {
      if (ol.contains(event.relatedTarget)) return
      if (!this.state.interacting) this.setExpanded(false)
    })
    ol.addEventListener("dragend", () => this.setExpanded(false))
    ol.addEventListener("pointerdown", (event) => {
      if (discreteTarget(event)) return
      this.setInteracting(true)
    })
    ol.addEventListener("pointerup", () => this.setInteracting(false))
  }

  destroy() {
    this.cleanups.forEach((c) => c())
    for (const comp of this.comps.values()) comp.unmount()
  }
}

// --- a Toast ----------------------------------------------------------------------------------------

class Toast {
  constructor(toaster, t) {
    this.toaster = toaster
    this.id = t.id
    this.el = document.createElement("li")
    this.styleRef = {}
    this.state = {
      swipeDirection: null,
      swipeOutDirection: null,
      mounted: false,
      removed: false,
      swiping: false,
      swipeOut: false,
      isSwiped: false,
      offsetBeforeRemove: 0,
      initialHeight: 0,
    }
    this.props = null
    this.offset = 0
    this.remainingTime = t.duration || toaster.props.toastOptions?.duration || toaster.props.duration || TOAST_LIFETIME
    this.dragStartTime = null
    this.closeTimerStartTime = 0
    this.lastCloseTimerStartTime = 0
    this.pointerStart = null
    this.timeoutId = null
    this.timerCleanup = null
    this.heightCleanup = null
    this.deps = {} // previous effect dependencies
    this.view = null // rendered children slots
    this.bind()
  }

  set(patch) {
    this.toaster.update(() => Object.assign(this.state, patch))
  }

  get duration() {
    const tp = this.toaster.props
    return this.props.toast.duration || tp.toastOptions?.duration || tp.duration || TOAST_LIFETIME
  }

  deleteToast = () => {
    const t = this.props.toast
    const offset = this.offset
    this.set({ removed: true })
    this.toaster.update(() => (this.state.offsetBeforeRemove = offset))
    this.toaster.setHeights((h) => h.filter((height) => height.toastId !== t.id))
    setTimeout(() => this.toaster.removeToast(t), TIME_BEFORE_UNMOUNT)
  }

  render(props, effects) {
    const prev = this.props
    this.props = props
    const tp = this.toaster.props
    const opts = tp.toastOptions ?? {}
    const { toast: t, index, heights, toasts, expanded, position } = props
    const st = this.state
    const isFront = index === 0
    const isVisible = index + 1 <= tp.visibleToasts
    const toastType = t.type
    const dismissible = t.dismissible !== false
    const heightIndex = heights.findIndex((h) => h.toastId === t.id) || 0
    const toastsHeightBefore = heights.reduce((prevH, curr, i) => (i >= heightIndex ? prevH : prevH + curr.height), 0)
    this.offset = heightIndex * tp.gap + toastsHeightBefore
    const invert = t.invert || tp.invert
    const [y, x] = position.split("-")
    const classNames = opts.classNames
    const el = this.el

    setAttr(el, "tabindex", "0")
    setAttr(el, "class", cn(opts.className, t.className || "", classNames?.toast, t.classNames?.toast, classNames?.default, classNames?.[toastType], t.classNames?.[toastType]))
    setAttr(el, "data-sonner-toast", "")
    setAttr(el, "data-rich-colors", t.richColors ?? tp.richColors)
    setAttr(el, "data-styled", !(t.jsx || t.unstyled || opts.unstyled))
    setAttr(el, "data-mounted", st.mounted)
    setAttr(el, "data-promise", Boolean(t.promise))
    setAttr(el, "data-swiped", st.isSwiped)
    setAttr(el, "data-removed", st.removed)
    setAttr(el, "data-visible", isVisible)
    setAttr(el, "data-y-position", y)
    setAttr(el, "data-x-position", x)
    setAttr(el, "data-index", index)
    setAttr(el, "data-front", isFront)
    setAttr(el, "data-swiping", st.swiping)
    setAttr(el, "data-dismissible", dismissible)
    setAttr(el, "data-type", toastType)
    setAttr(el, "data-invert", invert)
    setAttr(el, "data-swipe-out", st.swipeOut)
    setAttr(el, "data-swipe-direction", st.swipeOutDirection)
    setAttr(el, "data-expanded", Boolean(expanded || (tp.expand && st.mounted)))
    setAttr(el, "data-testid", t.testId)
    setStyles(
      el,
      {
        "--index": index,
        "--toasts-before": index,
        "--z-index": toasts.length - index,
        "--offset": `${st.removed ? st.offsetBeforeRemove : this.offset}px`,
        "--initial-height": tp.expand ? "auto" : `${st.initialHeight}px`,
        ...opts.style,
        ...t.style,
      },
      this.styleRef,
    )
    this.renderChildren(t)

    // --- effects, in declaration order ---
    const d = this.deps
    const changed = (key, values) => {
      const old = d[key]
      d[key] = values
      return !old || old.length !== values.length || old.some((v, i) => !Object.is(v, values[i]))
    }
    if (changed("duration", [this.duration])) effects.passive.push(() => (this.remainingTime = this.duration))
    if (changed("mount", [])) effects.passive.push(() => this.set({ mounted: true }))
    if (changed("height", [t.id])) {
      effects.passive.push(() => {
        const height = this.el.getBoundingClientRect().height
        this.set({ initialHeight: height })
        this.toaster.setHeights((h) => [{ toastId: t.id, height, position: t.position }, ...h])
        this.heightCleanup = () => this.toaster.setHeights((h) => h.filter((x) => x.toastId !== t.id))
      })
    }
    if (changed("layout", [st.mounted, t.title, t.description, t.id, t.jsx, t.action, t.cancel])) {
      effects.layout.push(() => {
        if (!this.state.mounted) return
        const node = this.el
        const originalHeight = node.style.height
        node.style.height = "auto"
        const newHeight = node.getBoundingClientRect().height
        node.style.height = originalHeight
        this.set({ initialHeight: newHeight })
        this.toaster.setHeights((hs) => {
          const exists = hs.find((h) => h.toastId === t.id)
          if (!exists) return [{ toastId: t.id, height: newHeight, position: t.position }, ...hs]
          return hs.map((h) => (h.toastId === t.id ? { ...h, height: newHeight } : h))
        })
      })
    }
    if (changed("timer", [expanded, props.interacting, t, toastType, props.isDocumentHidden])) {
      effects.passive.push(() => {
        this.timerCleanup?.()
        this.timerCleanup = null
        if ((t.promise && toastType === "loading") || t.duration === Infinity || t.type === "loading") return
        let timeoutId
        const pauseTimer = () => {
          if (this.lastCloseTimerStartTime < this.closeTimerStartTime) {
            const elapsedTime = new Date().getTime() - this.closeTimerStartTime
            this.remainingTime = this.remainingTime - elapsedTime
          }
          this.lastCloseTimerStartTime = new Date().getTime()
        }
        const startTimer = () => {
          if (this.remainingTime === Infinity) return
          this.closeTimerStartTime = new Date().getTime()
          timeoutId = setTimeout(() => {
            t.onAutoClose?.call(t, t)
            this.deleteToast()
          }, this.remainingTime)
        }
        if (expanded || props.interacting || props.isDocumentHidden) pauseTimer()
        else startTimer()
        this.timerCleanup = () => clearTimeout(timeoutId)
      })
    }
    if (changed("delete", [t, t.delete])) {
      effects.passive.push(() => {
        if (t.delete) {
          this.deleteToast()
          t.onDismiss?.call(t, t)
        }
      })
    }
    void prev
  }

  renderChildren(t) {
    const tp = this.toaster.props
    const opts = tp.toastOptions ?? {}
    const classNames = opts.classNames
    const icons = tp.icons
    const toastType = t.type
    const disabled = toastType === "loading"
    const dismissible = t.dismissible !== false
    const closeButton = t.closeButton ?? opts.closeButton ?? tp.closeButton
    const v = (this.view ??= { close: null, icon: null, loader: null, loaderKind: null, iconEl: null, iconKind: null, content: null, title: null, titleValue: undefined, description: null, descValue: undefined, cancel: null, cancelLabel: undefined, action: null, actionLabel: undefined })
    const el = this.el
    const slots = []

    // Close button.
    if (closeButton && !t.jsx && toastType !== "loading") {
      if (!v.close) {
        v.close = document.createElement("button")
        v.close.appendChild(icons.close ? iconNode(icons, "close") : closeIcon())
        v.close.addEventListener("click", () => {
          const tt = this.props.toast
          if (tt.type === "loading" || tt.dismissible === false) return
          this.deleteToast()
          tt.onDismiss?.call(tt, tt)
        })
      }
      setAttr(v.close, "aria-label", opts.closeButtonAriaLabel ?? "Close toast")
      setAttr(v.close, "data-disabled", disabled)
      setAttr(v.close, "data-close-button", "true")
      setAttr(v.close, "class", cn(classNames?.closeButton, t.classNames?.closeButton))
      slots.push(v.close)
    } else v.close = null

    // Icon.
    const hasIcon = (toastType || t.icon || t.promise) && t.icon !== null && (icons[toastType] !== null || t.icon)
    if (hasIcon) {
      if (!v.icon) v.icon = document.createElement("div")
      setAttr(v.icon, "data-icon", "")
      setAttr(v.icon, "class", cn(classNames?.icon, t.classNames?.icon))
      const iconChildren = []
      // Loader slot.
      if (t.promise || (t.type === "loading" && !t.icon)) {
        if (t.icon) {
          if (v.loaderKind !== "custom") {
            v.loader = content(t.icon)[0] ?? null
            v.loaderKind = "custom"
          }
        } else if (icons.loading) {
          if (v.loaderKind !== "icon") {
            v.loader = document.createElement("div")
            v.loader.appendChild(iconNode(icons, "loading"))
            v.loaderKind = "icon"
          }
          setAttr(v.loader, "class", cn(classNames?.loader, t.classNames?.loader, "sonner-loader"))
          setAttr(v.loader, "data-visible", toastType === "loading")
        } else {
          if (v.loaderKind !== "bars") {
            v.loader = loaderBars(toastType === "loading", cn(classNames?.loader, t.classNames?.loader))
            v.loaderKind = "bars"
          }
          setAttr(v.loader, "data-visible", toastType === "loading")
        }
        if (v.loader) iconChildren.push(v.loader)
      } else {
        v.loader = null
        v.loaderKind = null
      }
      // Type icon slot.
      if (t.type !== "loading") {
        const kind = t.icon ? t.icon : icons[toastType] ? `icon:${toastType}` : `asset:${toastType}`
        if (v.iconKind !== kind) {
          v.iconEl = t.icon ? content(t.icon)[0] ?? null : icons[toastType] ? iconNode(icons, toastType) : getAsset(toastType)
          v.iconKind = kind
        }
        if (v.iconEl) iconChildren.push(v.iconEl)
      } else {
        v.iconEl = null
        v.iconKind = null
      }
      syncChildren(v.icon, iconChildren)
      slots.push(v.icon)
    } else {
      v.icon = null
      v.loader = null
      v.loaderKind = null
      v.iconEl = null
      v.iconKind = null
    }

    // Content.
    if (!v.content) {
      v.content = document.createElement("div")
      v.title = document.createElement("div")
      v.content.appendChild(v.title)
    }
    setAttr(v.content, "data-content", "")
    setAttr(v.content, "class", cn(classNames?.content, t.classNames?.content))
    setAttr(v.title, "data-title", "")
    setAttr(v.title, "class", cn(classNames?.title, t.classNames?.title))
    const titleValue = t.jsx ? t.jsx : t.title
    if (!sameContent(v.titleValue, titleValue) || v.titleValue === undefined) {
      if (typeof titleValue === "string" && v.title.childNodes.length === 1 && v.title.firstChild.nodeType === 3) v.title.firstChild.data = titleValue
      else v.title.replaceChildren(...content(titleValue))
      v.titleValue = titleValue
    }
    if (t.description) {
      if (!v.description) v.description = document.createElement("div")
      setAttr(v.description, "data-description", "")
      setAttr(v.description, "class", cn(opts.descriptionClassName, t.descriptionClassName || "", classNames?.description, t.classNames?.description))
      if (!sameContent(v.descValue, t.description)) {
        if (typeof t.description === "string" && v.description.childNodes.length === 1 && v.description.firstChild.nodeType === 3) v.description.firstChild.data = t.description
        else v.description.replaceChildren(...content(t.description))
        v.descValue = t.description
      }
      syncChildren(v.content, [v.title, v.description])
    } else {
      v.description = null
      v.descValue = undefined
      syncChildren(v.content, [v.title])
    }
    slots.push(v.content)

    // Cancel and action buttons.
    const button = (key, action, kind) => {
      if (isNode(action)) {
        if (v[key + "Node"] !== action) {
          v[key] = action.cloneNode(true)
          v[key + "Node"] = action
        }
        return v[key]
      }
      if (!action || action.label === undefined) {
        v[key] = null
        return null
      }
      if (!v[key] || v[key].tagName !== "BUTTON") {
        v[key] = document.createElement("button")
        v[key].addEventListener("click", (event) => this.onButtonClick(kind, event))
      }
      const b = v[key]
      setAttr(b, "data-button", "true")
      setAttr(b, kind === "cancel" ? "data-cancel" : "data-action", "true")
      const style = kind === "cancel" ? t.cancelButtonStyle || opts.cancelButtonStyle : t.actionButtonStyle || opts.actionButtonStyle
      if (style) for (const [k, val] of Object.entries(style)) b.style[k] = val
      setAttr(b, "class", cn(kind === "cancel" ? classNames?.cancelButton : classNames?.actionButton, kind === "cancel" ? t.classNames?.cancelButton : t.classNames?.actionButton))
      if (v[key + "Label"] !== action.label) {
        b.replaceChildren(...content(action.label))
        v[key + "Label"] = action.label
      }
      return b
    }
    const cancel = button("cancel", t.cancel, "cancel")
    if (cancel) slots.push(cancel)
    const action = button("action", t.action, "action")
    if (action) slots.push(action)

    syncChildren(el, slots)
  }

  onButtonClick(kind, event) {
    const t = this.props.toast
    const a = kind === "cancel" ? t.cancel : t.action
    if (!a || a.label === undefined) return
    if (kind === "cancel") {
      if (t.dismissible === false) return
      a.onClick?.call(a, event)
      this.deleteToast()
      return
    }
    a.onClick?.call(a, event)
    if (event.defaultPrevented) return
    this.deleteToast()
  }

  bind() {
    const el = this.el
    el.addEventListener("dragend", () => {
      this.set({ swiping: false, swipeDirection: null })
      this.pointerStart = null
    })
    el.addEventListener("pointerdown", (event) => {
      const t = this.props.toast
      if (event.button === 2) return
      if (t.type === "loading" || t.dismissible === false) return
      this.dragStartTime = new Date()
      const offset = this.offset
      this.toaster.update(() => (this.state.offsetBeforeRemove = offset))
      event.target.setPointerCapture(event.pointerId)
      if (event.target.tagName === "BUTTON") return
      this.set({ swiping: true })
      this.pointerStart = { x: event.clientX, y: event.clientY }
    })
    el.addEventListener("pointerup", () => {
      const t = this.props.toast
      const st = this.state
      const dismissible = t.dismissible !== false
      if (st.swipeOut || !dismissible) return
      this.pointerStart = null
      const swipeAmountX = Number(el.style.getPropertyValue("--swipe-amount-x").replace("px", "") || 0)
      const swipeAmountY = Number(el.style.getPropertyValue("--swipe-amount-y").replace("px", "") || 0)
      const timeTaken = new Date().getTime() - this.dragStartTime?.getTime()
      const swipeAmount = st.swipeDirection === "x" ? swipeAmountX : swipeAmountY
      const velocity = Math.abs(swipeAmount) / timeTaken
      if (Math.abs(swipeAmount) >= SWIPE_THRESHOLD || velocity > 0.11) {
        const offset = this.offset
        this.toaster.update(() => (this.state.offsetBeforeRemove = offset))
        t.onDismiss?.call(t, t)
        if (st.swipeDirection === "x") this.set({ swipeOutDirection: swipeAmountX > 0 ? "right" : "left" })
        else this.set({ swipeOutDirection: swipeAmountY > 0 ? "down" : "up" })
        this.deleteToast()
        this.set({ swipeOut: true })
        return
      }
      el.style.setProperty("--swipe-amount-x", "0px")
      el.style.setProperty("--swipe-amount-y", "0px")
      this.set({ isSwiped: false, swiping: false, swipeDirection: null })
    })
    el.addEventListener("pointermove", (event) => {
      const t = this.props.toast
      const dismissible = t.dismissible !== false
      if (!this.pointerStart || !dismissible) return
      const isHighlighted = window.getSelection()?.toString().length > 0
      if (isHighlighted) return
      const yDelta = event.clientY - this.pointerStart.y
      const xDelta = event.clientX - this.pointerStart.x
      const swipeDirections = this.toaster.props.swipeDirections ?? getDefaultSwipeDirections(this.props.position)
      const swipeDirection = this.state.swipeDirection
      if (!swipeDirection && (Math.abs(xDelta) > 1 || Math.abs(yDelta) > 1)) {
        this.set({ swipeDirection: Math.abs(xDelta) > Math.abs(yDelta) ? "x" : "y" })
      }
      const swipeAmount = { x: 0, y: 0 }
      const getDampening = (delta) => 1 / (1.5 + Math.abs(delta) / 20)
      if (swipeDirection === "y") {
        if (swipeDirections.includes("top") || swipeDirections.includes("bottom")) {
          if ((swipeDirections.includes("top") && yDelta < 0) || (swipeDirections.includes("bottom") && yDelta > 0)) swipeAmount.y = yDelta
          else {
            const dampenedDelta = yDelta * getDampening(yDelta)
            swipeAmount.y = Math.abs(dampenedDelta) < Math.abs(yDelta) ? dampenedDelta : yDelta
          }
        }
      } else if (swipeDirection === "x") {
        if (swipeDirections.includes("left") || swipeDirections.includes("right")) {
          if ((swipeDirections.includes("left") && xDelta < 0) || (swipeDirections.includes("right") && xDelta > 0)) swipeAmount.x = xDelta
          else {
            const dampenedDelta = xDelta * getDampening(xDelta)
            swipeAmount.x = Math.abs(dampenedDelta) < Math.abs(xDelta) ? dampenedDelta : xDelta
          }
        }
      }
      if (Math.abs(swipeAmount.x) > 0 || Math.abs(swipeAmount.y) > 0) this.set({ isSwiped: true })
      el.style.setProperty("--swipe-amount-x", `${swipeAmount.x}px`)
      el.style.setProperty("--swipe-amount-y", `${swipeAmount.y}px`)
    })
  }

  unmount() {
    this.timerCleanup?.()
    this.timerCleanup = null
    this.heightCleanup?.()
    this.heightCleanup = null
  }
}

/** Makes parent's children exactly nodes, in order, moving as few as possible. */
function syncChildren(parent, nodes) {
  let i = 0
  for (const node of nodes) {
    const current = parent.childNodes[i]
    if (current !== node) parent.insertBefore(node, current ?? null)
    i++
  }
  while (parent.childNodes.length > nodes.length) parent.lastChild.remove()
}

// --- mounting ----------------------------------------------------------------------------------------

const toasters = new Set()

define("sonner", {
  selector: "template[data-sonner-toaster]",
  mount(tpl) {
    const section = tpl.parentElement
    if (!section) return
    const t = new Toaster(section, tpl)
    toasters.add(t)
    return () => {
      toasters.delete(t)
      t.destroy()
    }
  },
})

// Declarative triggers rendered by sonner.OnClick in Go: a click on an
// element carrying data-sonner-trigger shows the toast it describes.
document.addEventListener("click", (event) => {
  const el = event.target instanceof Element ? event.target.closest("[data-sonner-trigger]") : null
  if (!el) return
  let spec
  try {
    spec = JSON.parse(el.getAttribute("data-sonner-trigger"))
  } catch {
    return
  }
  showSpec(spec)
})

/**
 * Shows a toast described by a JSON-able spec:
 * { type?: "success"|..., message, description?, action?: label, cancel?: label,
 *   duration?, id?, promise?: { loading, success, error, ms, reject } }
 */
export function showSpec(spec) {
  const data = {}
  if (spec.description !== undefined) data.description = spec.description
  if (spec.duration !== undefined) data.duration = spec.duration
  if (spec.id !== undefined) data.id = spec.id
  if (spec.closeButton !== undefined) data.closeButton = spec.closeButton
  if (spec.richColors !== undefined) data.richColors = spec.richColors
  if (spec.action !== undefined) data.action = { label: spec.action, onClick: () => document.dispatchEvent(new CustomEvent("shadcn:toast-action", { detail: { id: spec.id, label: spec.action } })) }
  if (spec.cancel !== undefined) data.cancel = { label: spec.cancel, onClick: () => document.dispatchEvent(new CustomEvent("shadcn:toast-cancel", { detail: { id: spec.id, label: spec.cancel } })) }
  if (spec.promise) {
    const pr = spec.promise
    return toast.promise(() => new Promise((resolve, reject) => setTimeout(() => (pr.reject ? reject(new Error(pr.error ?? "Error")) : resolve(pr.value)), pr.ms ?? 0)), {
      ...data,
      loading: pr.loading,
      success: pr.success,
      error: pr.error,
    })
  }
  switch (spec.type) {
    case "success":
    case "info":
    case "warning":
    case "error":
    case "loading":
    case "message":
      return toast[spec.type](spec.message, data)
    default:
      return toast(spec.message, data)
  }
}
