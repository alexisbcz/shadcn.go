// sidebar.js: SidebarProvider's state and behavior (shadcn/ui sidebar.tsx)
// for the server-rendered markup of components/ui/sidebar.
//
// Markup contract (rendered by the Go package):
//   - The provider's wrapper, div[data-slot="sidebar-wrapper"], carries
//     data-sidebar-default-open="true|false" (removed on mount).
//   - Every Sidebar that can collapse renders its desktop markup
//     ([data-slot="sidebar"] with data-state, data-collapsible, ...) holding
//     a <template data-sidebar-mobile data-sidebar-collapsible="<mode>"> with the
//     Sheet it becomes below the md breakpoint; the sheet's
//     [data-sidebar-mobile-body] receives the sidebar's children.
//   - MenuButton tooltips are tooltip templates marked
//     data-sidebar-tooltip: their content is hidden unless the sidebar is
//     collapsed on desktop.
//
// State, like the React context: open (desktop: expanded or collapsed,
// persisted in the sidebar_state cookie), openMobile (the mobile sheet) and
// isMobile (window.innerWidth < 768, tracked with matchMedia).
//
//   wrapper.addEventListener("shadcn:open-change", (e) => e.target === wrapper && e.detail.open) // cancelable
//   import { getSidebar } from "/_shadcn/sidebar.js"
//   getSidebar(wrapperOrDescendant)  // { state, open, setOpen, openMobile, setOpenMobile, isMobile, toggleSidebar }

import { define, emit } from "./core.js"
import { instantiate } from "./portal.js"
import { getDialog } from "./dialog.js"

const COOKIE_NAME = "sidebar_state"
const COOKIE_MAX_AGE = 60 * 60 * 24 * 7
const KEYBOARD_SHORTCUT = "b"
const MOBILE_BREAKPOINT = 768

const providers = new WeakMap() // wrapper -> controller

// Pristine copies of every sidebar's children, taken when this module is
// evaluated: the document is parsed but no behavior has mounted yet (the
// runtime starts after importing every module). React remounts a Sidebar's
// children whenever it switches between the desktop markup and the mobile
// Sheet, and the Sheet mounts them afresh on every open: copies of these
// are what gets mounted then.
const pristine = new WeakMap() // [data-slot="sidebar-inner"] -> Node[]
const snapshot = (inner) => [...inner.childNodes].map((n) => n.cloneNode(true))
for (const inner of document.querySelectorAll('[data-slot="sidebar-inner"]')) pristine.set(inner, snapshot(inner))
// Live mobile sheet contents -> controller (Radix portals them to <body>,
// outside the wrapper, while React context still reaches them).
const mobileContents = new Map()

/** Returns the sidebar controller owning el (a wrapper or any element in it, mobile sheet included). */
export function getSidebar(el) {
  if (!(el instanceof Element)) return undefined
  const wrapper = el.closest('[data-slot="sidebar-wrapper"]')
  if (wrapper && providers.has(wrapper)) return providers.get(wrapper)
  const mobile = el.closest('[data-slot="sidebar"][data-mobile="true"]')
  return mobile ? mobileContents.get(mobile) : undefined
}

function createProvider(wrapper) {
  let open = wrapper.getAttribute("data-sidebar-default-open") !== "false"
  wrapper.removeAttribute("data-sidebar-default-open")
  let openMobile = false
  let isMobile = false

  // Sidebars that can collapse: { root, inner, tpl, mode, mobile } where
  // mobile is the dialog template standing in for root while isMobile.
  const sidebars = [...wrapper.querySelectorAll('[data-slot="sidebar"] > template[data-sidebar-mobile]')]
    .filter((tpl) => tpl.parentElement.parentElement?.closest('[data-slot="sidebar-wrapper"]') === wrapper)
    .map((tpl) => ({
      root: tpl.parentElement,
      inner: tpl.parentElement.querySelector(':scope > [data-slot="sidebar-container"] > [data-sidebar="sidebar"]'),
      tpl,
      mode: tpl.getAttribute("data-sidebar-collapsible") || "offcanvas",
      mobile: null,
    }))
  for (const s of sidebars) s.pristine = pristine.get(s.inner) ?? snapshot(s.inner)
  const fresh = (s) => s.pristine.map((n) => n.cloneNode(true))

  const tooltipTemplates = () => [...wrapper.querySelectorAll("template[data-sidebar-tooltip]")]

  // TooltipContent hidden={state !== "collapsed" || isMobile}.
  function renderTooltips() {
    const hidden = open || isMobile
    for (const tpl of tooltipTemplates()) {
      for (const content of tpl.content.querySelectorAll('[data-slot="tooltip-content"]')) content.toggleAttribute("hidden", hidden)
      const id = tpl.getAttribute("data-content-id")
      const live = id ? document.getElementById(id)?.closest('[data-slot="tooltip-content"]') : null
      if (live) live.toggleAttribute("hidden", hidden)
    }
  }

  function render() {
    const state = open ? "expanded" : "collapsed"
    for (const s of sidebars) {
      s.root.setAttribute("data-state", state)
      s.root.setAttribute("data-collapsible", open ? "" : s.mode)
    }
    renderTooltips()
  }

  function setOpen(value) {
    const next = typeof value === "function" ? !!value(open) : !!value
    // This sets the cookie to keep the sidebar state (also when a
    // listener vetoes the change, like a controlled provider).
    document.cookie = `${COOKIE_NAME}=${next}; path=/; max-age=${COOKIE_MAX_AGE}`
    if (next === open) return
    if (!emit(wrapper, "shadcn:open-change", { open: next })) return
    open = next
    render()
  }

  function setOpenMobile(value) {
    const next = typeof value === "function" ? !!value(openMobile) : !!value
    if (next === openMobile) return
    openMobile = next
    for (const s of sidebars) {
      const id = s.mobile?.getAttribute("data-content-id")
      const dialog = id ? getDialog(id) : null
      if (dialog && dialog.open !== openMobile) setSheetOpen(dialog, openMobile)
    }
  }

  // The Sheet's FocusScope focuses its first tabbable as it mounts the
  // content, before the runtime mounts the behaviors of that fresh content
  // (in its MutationObserver callback, a microtask). React's handlers
  // exist from the first render: a MenuButton tooltip trigger focused that
  // way opens its tooltip (hidden on mobile, still in the DOM). Replay the
  // missed focus once the content's behaviors are mounted.
  function setSheetOpen(dialog, value) {
    const before = document.activeElement
    dialog.setOpen(value)
    if (!value) return
    const focused = document.activeElement
    if (focused === before || !focused?.closest('[data-slot="sidebar"][data-mobile="true"]')) return
    queueMicrotask(() => {
      if (document.activeElement === focused) focused.dispatchEvent(new FocusEvent("focus"))
    })
  }

  const toggleSidebar = () => (isMobile ? setOpenMobile((o) => !o) : setOpen((o) => !o))

  // --- mobile: the Sheet replaces the desktop sidebar ---------------------------------------

  function onMobileOpenChange(event) {
    if (event.target !== event.currentTarget || event.defaultPrevented) return
    // Sheet onOpenChange={setOpenMobile}
    queueMicrotask(() => {
      if (!event.defaultPrevented) openMobile = !!event.detail.open
    })
  }

  // isMobile switches between <Sheet> and the desktop markup: everything
  // inside the Sidebar remounts, with fresh state.
  function toMobile(s) {
    const [dialogTpl] = instantiate(s.tpl)
    const body = dialogTpl.content.querySelector("[data-sidebar-mobile-body]")
    body.removeAttribute("data-sidebar-mobile-body")
    body.append(...fresh(s))
    s.inner.replaceChildren()
    s.mobile = dialogTpl
    dialogTpl.addEventListener("shadcn:open-change", onMobileOpenChange)
    s.root.replaceWith(dialogTpl)
    if (openMobile)
      queueMicrotask(() => {
        const dialog = getDialog(dialogTpl.getAttribute("data-content-id"))
        if (dialog) setSheetOpen(dialog, true)
      })
  }

  function toDesktop(s) {
    const dialogTpl = s.mobile
    if (!dialogTpl) return
    s.inner.replaceChildren(...fresh(s))
    dialogTpl.removeEventListener("shadcn:open-change", onMobileOpenChange)
    dialogTpl.replaceWith(s.root)
    s.mobile = null
  }

  function setIsMobile(next) {
    if (next === isMobile) return
    isMobile = next
    for (const s of sidebars) (isMobile ? toMobile : toDesktop)(s)
    renderTooltips()
  }

  // Track the sheet content while it is mounted, for getSidebar().
  const observer = new MutationObserver(() => {
    for (const [content, owner] of mobileContents) if (owner === controller && !content.isConnected) mobileContents.delete(content)
    for (const s of sidebars) {
      const id = s.mobile?.getAttribute("data-content-id")
      const content = id ? document.getElementById(id) : null
      if (content) mobileContents.set(content, controller)
    }
  })
  observer.observe(document.body, { childList: true })

  // useIsMobile: window.innerWidth < 768, read when the media query
  // changes; React applies the state update in a scheduler task (a
  // MessageChannel message), batching the updates made meanwhile.
  const mql = window.matchMedia(`(max-width: ${MOBILE_BREAKPOINT - 1}px)`)
  const channel = new MessageChannel()
  let pending = null
  channel.port1.onmessage = () => {
    const next = pending
    pending = null
    if (next !== null) setIsMobile(next)
  }
  const onMediaChange = () => {
    if (pending === null) channel.port2.postMessage(null)
    pending = window.innerWidth < MOBILE_BREAKPOINT
  }
  mql.addEventListener("change", onMediaChange)

  const onKeyDown = (event) => {
    if (event.key === KEYBOARD_SHORTCUT && (event.metaKey || event.ctrlKey)) {
      event.preventDefault()
      toggleSidebar()
    }
  }
  window.addEventListener("keydown", onKeyDown)

  const controller = {
    get state() {
      return open ? "expanded" : "collapsed"
    },
    get open() {
      return open
    },
    setOpen,
    get openMobile() {
      return openMobile
    },
    setOpenMobile,
    get isMobile() {
      return isMobile
    },
    toggleSidebar,
    dispose() {
      mql.removeEventListener("change", onMediaChange)
      channel.port1.onmessage = null
      window.removeEventListener("keydown", onKeyDown)
      observer.disconnect()
      for (const s of sidebars) s.mobile?.removeEventListener("shadcn:open-change", onMobileOpenChange)
      providers.delete(wrapper)
    },
  }
  providers.set(wrapper, controller)
  // useIsMobile starts undefined (desktop markup) and its effect sets the
  // real value: React re-renders afterwards. Doing it later also lets the
  // runtime see the Sheet's template inserted (its behavior mounts then),
  // which it would miss while it is still scanning the page.
  pending = window.innerWidth < MOBILE_BREAKPOINT
  channel.port2.postMessage(null)
  return controller
}

// SidebarTrigger and SidebarRail: onClick toggles the sidebar (after the
// user's own handlers).
document.addEventListener("click", (event) => {
  const el = event.target instanceof Element ? event.target.closest('[data-sidebar="trigger"], [data-sidebar="rail"]') : null
  if (!el) return
  getSidebar(el)?.toggleSidebar()
})

// A sidebar part rendered asChild over a CollapsibleTrigger (e.g.
// <SidebarGroupLabel asChild><CollapsibleTrigger>) keeps the part's
// data-slot, as in React, but collapsible.js recognizes its triggers by
// data-slot="collapsible-trigger". While such a click is dispatched, the
// trigger carries that data-slot so the collapsible handles it (and updates
// it with its other triggers); it is restored right after.
// TODO: drop once collapsible.js identifies triggers by aria-controls.
document.addEventListener(
  "click",
  (event) => {
    const el = event.target instanceof Element ? event.target.closest("[data-sidebar][aria-controls]") : null
    if (!el || el.getAttribute("data-slot") === "collapsible-trigger") return
    const root = el.parentElement?.closest('[data-slot="collapsible"]')
    const id = el.getAttribute("aria-controls")
    const content = root ? [...root.querySelectorAll('[data-slot="collapsible-content"]')].find((c) => c.id === id) : null
    if (!content || content.parentElement?.closest('[data-slot="collapsible"]') !== root) return
    const slot = el.getAttribute("data-slot")
    el.setAttribute("data-slot", "collapsible-trigger")
    let restored = false
    const restore = () => {
      if (restored) return
      restored = true
      if (slot === null) el.removeAttribute("data-slot")
      else el.setAttribute("data-slot", slot)
    }
    window.addEventListener("click", restore, { once: true })
    setTimeout(restore, 0)
  },
  true,
)

define("sidebar", {
  selector: '[data-slot="sidebar-wrapper"]',
  mount(wrapper) {
    const c = createProvider(wrapper)
    return () => c.dispose()
  },
})
