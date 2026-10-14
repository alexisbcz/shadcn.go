// avatar.js: Radix Avatar (@radix-ui/react-avatar 1.1.10).
//
// Each avatar tracks its image's loading status ("idle" until the image
// is known, then "loading", "loaded" or "error"), loading the image with a
// detached `new Image()` like Radix's useImageLoadingStatus:
// - the image (from its <template data-slot="avatar-image-template">) is
//   in the DOM only while its status is "loaded";
// - each fallback is in the DOM while the status is not "loaded" and, for
//   fallbacks with a delay (<template data-slot="avatar-fallback-template"
//   data-delay>), once the delay has passed.
// "shadcn:loading-status-change" (detail: { status }) is dispatched on the
// avatar when an image's status changes.
import { define, emit } from "./core.js"

// A part is shown by inserting it before its anchor, an empty comment left
// where React would render it.
function part(node, anchorFrom) {
  const anchor = document.createComment("")
  anchorFrom.replaceWith(anchor)
  return {
    node,
    anchor,
    show(on) {
      if (on && !node.isConnected) anchor.parentNode?.insertBefore(node, anchor)
      else if (!on && node.isConnected) node.remove()
    },
  }
}

define("avatar", {
  selector: '[data-slot="avatar"]',
  mount(root) {
    const own = (el) => el.parentElement?.closest('[data-slot="avatar"]') === root
    const cleanups = []
    let status = "idle"

    const fallbacks = []
    for (const el of root.querySelectorAll('[data-slot="avatar-fallback"]')) {
      if (own(el)) fallbacks.push({ part: part(el, el), canRender: true })
    }
    for (const tpl of root.querySelectorAll('template[data-slot="avatar-fallback-template"]')) {
      if (!own(tpl)) continue
      const node = tpl.content.firstElementChild
      if (!node) continue
      const f = { part: part(node, tpl), canRender: false }
      fallbacks.push(f)
      const delay = Number(tpl.getAttribute("data-delay")) || 0
      const timer = window.setTimeout(() => {
        f.canRender = true
        render()
      }, delay)
      cleanups.push(() => window.clearTimeout(timer))
    }
    const render = () => {
      for (const img of images) img.part.show(img.status === "loaded")
      for (const f of fallbacks) f.part.show(f.canRender && status !== "loaded")
    }

    const images = []
    for (const tpl of root.querySelectorAll('template[data-slot="avatar-image-template"]')) {
      if (!own(tpl)) continue
      const node = tpl.content.firstElementChild
      if (!node) continue
      images.push({ part: part(node, tpl), status: "idle", node })
    }

    const setStatus = (img, next) => {
      if (img.status === next) return
      img.status = next
      if (next !== "idle") {
        status = next
        emit(root, "shadcn:loading-status-change", { status: next })
      }
      render()
    }

    for (const img of images) {
      const src = img.node.getAttribute("src")
      const image = new window.Image()
      const onLoad = () => setStatus(img, "loaded")
      const onError = () => setStatus(img, "error")
      const referrerPolicy = img.node.getAttribute("referrerpolicy")
      const crossOrigin = img.node.getAttribute("crossorigin")
      // Radix resolves the status at first render (setting image.src)
      // before attaching the listeners.
      let initial
      if (!src) initial = "error"
      else {
        image.src = src
        initial = image.complete && image.naturalWidth > 0 ? "loaded" : "loading"
      }
      image.addEventListener("load", onLoad)
      image.addEventListener("error", onError)
      if (referrerPolicy) image.referrerPolicy = referrerPolicy
      if (crossOrigin !== null) image.crossOrigin = crossOrigin
      cleanups.push(() => {
        image.removeEventListener("load", onLoad)
        image.removeEventListener("error", onError)
      })
      setStatus(img, initial)
    }
    render()

    return () => cleanups.forEach((fn) => fn())
  },
})
