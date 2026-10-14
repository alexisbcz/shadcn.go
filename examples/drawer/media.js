// The drawer-dialog example renders a Dialog when useMediaQuery(query)
// matches and a Drawer otherwise. The server renders the Dialog version
// (the elements before the template) and the Drawer version inside the
// template; this swaps them whenever the media query changes.
import { define } from "./core.js"

define("example-drawer-dialog", {
  selector: "template[data-example-media-swap]",
  mount(tpl) {
    const parent = tpl.parentElement
    if (!parent) return
    const query = window.matchMedia(tpl.getAttribute("data-example-media-swap"))
    // The matching version: the nodes rendered before the template.
    let matching = [...parent.childNodes].slice(0, [...parent.childNodes].indexOf(tpl))
    let other = null
    const show = (matches) => {
      if (matches) {
        if (!other) return
        other.forEach((n) => n.remove())
        other = null
        matching.forEach((n) => parent.insertBefore(n, tpl))
      } else {
        if (other) return
        matching.forEach((n) => n.remove())
        other = [...tpl.content.cloneNode(true).childNodes]
        other.forEach((n) => parent.insertBefore(n, tpl))
      }
    }
    const onChange = (e) => show(e.matches)
    query.addEventListener("change", onChange)
    show(query.matches)
    return () => query.removeEventListener("change", onChange)
  },
})
