// carousel.js: shadcn/ui Carousel behavior. Embla Carousel runs on the
// viewport exactly like useEmblaCarousel in the upstream component; this
// module adds what carousel.tsx adds around it: the keyboard map
// (ArrowLeft/ArrowRight, captured on the root), the Previous/Next buttons
// and their disabled state (canScrollPrev/canScrollNext, refreshed on
// "select" and "reInit"), and the autoplay plugin.
//
// The server renders the root with a data-carousel attribute holding the
// configuration ({ axis, opts, autoplay, stopOnHover }); it is removed on
// mount.
import { define, emit } from "./core.js"
import EmblaCarousel from "./embla-carousel.js"
import Autoplay from "./embla-carousel-autoplay.js"

const apis = new WeakMap() // root -> Embla API
const waiting = new WeakMap() // root -> callbacks waiting for the API

/**
 * Calls cb with the Embla API of a carousel root ([data-slot=carousel]) as
 * soon as it is available.
 * @param {Element} root
 * @param {(api: any) => void} cb
 */
export function onCarouselApi(root, cb) {
  const api = apis.get(root)
  if (api) {
    cb(api)
    return
  }
  if (!waiting.has(root)) waiting.set(root, [])
  waiting.get(root).push(cb)
}

function own(root, selector) {
  return [...root.querySelectorAll(selector)].filter((el) => el.closest('[data-slot="carousel"]') === root)
}

function readConfig(root) {
  const raw = root.getAttribute("data-carousel")
  root.removeAttribute("data-carousel")
  if (!raw) return {}
  try {
    return JSON.parse(raw)
  } catch {
    return {}
  }
}

define("carousel", {
  selector: '[data-slot="carousel"]',
  mount(root) {
    const cfg = readConfig(root)
    const viewport = own(root, '[data-slot="carousel-content"]')[0]
    if (!viewport) return
    const plugins = []
    let autoplay = null
    if (cfg.autoplay) {
      autoplay = Autoplay(cfg.autoplay)
      plugins.push(autoplay)
    }
    const api = EmblaCarousel(viewport, { ...(cfg.opts ?? {}), axis: cfg.axis === "y" ? "y" : "x" }, plugins)

    const prevs = own(root, '[data-slot="carousel-previous"]')
    const nexts = own(root, '[data-slot="carousel-next"]')
    const setDisabled = (buttons, disabled) => {
      for (const b of buttons) b.toggleAttribute("disabled", disabled)
    }
    const onSelect = (api) => {
      setDisabled(prevs, !api.canScrollPrev())
      setDisabled(nexts, !api.canScrollNext())
      emit(root, "shadcn:carousel-select", { api, index: api.selectedScrollSnap(), count: api.scrollSnapList().length })
    }
    const scrollPrev = () => api.scrollPrev()
    const scrollNext = () => api.scrollNext()
    const onKeyDown = (event) => {
      if (event.key === "ArrowLeft") {
        event.preventDefault()
        scrollPrev()
      } else if (event.key === "ArrowRight") {
        event.preventDefault()
        scrollNext()
      }
    }
    root.addEventListener("keydown", onKeyDown, true)
    prevs.forEach((b) => b.addEventListener("click", scrollPrev))
    nexts.forEach((b) => b.addEventListener("click", scrollNext))

    // onMouseEnter={plugin.stop} onMouseLeave={plugin.reset}
    const stop = () => autoplay?.stop()
    const reset = () => autoplay?.reset()
    if (cfg.stopOnHover) {
      root.addEventListener("mouseenter", stop)
      root.addEventListener("mouseleave", reset)
    }

    onSelect(api)
    api.on("reInit", onSelect)
    api.on("select", onSelect)

    apis.set(root, api)
    emit(root, "shadcn:carousel-init", { api })
    for (const cb of waiting.get(root) ?? []) cb(api)
    waiting.delete(root)

    return () => {
      root.removeEventListener("keydown", onKeyDown, true)
      root.removeEventListener("mouseenter", stop)
      root.removeEventListener("mouseleave", reset)
      prevs.forEach((b) => b.removeEventListener("click", scrollPrev))
      nexts.forEach((b) => b.removeEventListener("click", scrollNext))
      apis.delete(root)
      api.destroy()
    }
  },
})
