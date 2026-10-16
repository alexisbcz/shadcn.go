// portal.js: port of Radix Portal (@radix-ui/react-portal 1.1.9) for
// server-rendered markup.
//
// Radix renders portalled content (Dialog, Popover, Menu, Select, Tooltip,
// ... content) into document.body, and only while it is open: closed
// content is not in the DOM at all. The Go components therefore render
// that content inside an inert <template> where React's tree would have
// it; the component behavior clones it into the container when it opens
// (fresh elements every time, exactly like React mounting them) and
// removes the clone once presence says it may leave.
//
// Templates carry `data-parity-ignore` so the parity harness, like the
// browser, treats them as absent.
//
// API
//
//   import { instantiate, mount, takeMarked } from "./portal.js"
//
//   const nodes = instantiate(template)   // element children of a fresh
//                                         // clone of template.content; form
//                                         // controls keep the <form> around
//                                         // the template (form attribute)
//   const unmount = mount(nodes, container = document.body)
//                                         // appends them, in order, at the
//                                         // end of container (createPortal)
//   unmount()                             // removes those still attached
//
//   // Server-rendered markers: the Go side may tag elements inside a
//   // template with marker attributes (e.g. data-radix-dialog-close) the
//   // behavior needs but React never renders. takeMarked finds them, removes
//   // the attribute and returns [{ el, value }], so the live DOM matches
//   // React's exactly.
//   for (const { el, value } of takeMarked(root, "data-radix-dialog-close")) ...

/**
 * Returns the element children of a fresh deep clone of template's content.
 * @param {HTMLTemplateElement} template
 * @returns {Element[]}
 */
export function instantiate(template) {
  const fragment = template.content.cloneNode(true)
  const nodes = [...fragment.children]
  keepForm(template, nodes)
  return nodes
}

const formControls = "button, input, select, textarea, output, fieldset, object"

/**
 * Keeps the form controls of portalled content in the <form> that wraps its
 * template (as in the upstream examples, where a form wraps DialogContent):
 * mounted in document.body they would leave it, so they get the HTML form
 * attribute pointing at it, and submit buttons submit it.
 * @param {HTMLTemplateElement} template
 * @param {Element[]} nodes
 */
function keepForm(template, nodes) {
  const form = template.closest?.("form")
  if (!form) return
  let id = form.id
  for (const n of nodes) {
    const controls = n.matches?.(formControls) ? [n, ...n.querySelectorAll(formControls)] : [...n.querySelectorAll(formControls)]
    for (const c of controls) {
      if (c.hasAttribute("form") || c.parentElement?.closest("form")) continue
      if (!id) {
        id = form.id = "form-" + Math.random().toString(36).slice(2, 10)
      }
      c.setAttribute("form", id)
    }
  }
}

/**
 * Appends nodes at the end of container (Radix Portal's default container is
 * document.body) and returns a function that removes them.
 * @param {Element[]} nodes
 * @param {Element} [container]
 */
export function mount(nodes, container = globalThis.document.body) {
  for (const n of nodes) container.appendChild(n)
  return () => {
    for (const n of nodes) n.remove()
  }
}

/**
 * Finds the elements in root (root included) carrying the marker attribute
 * name, removes the attribute and returns them with its value.
 * @param {Element | DocumentFragment | Element[]} root
 * @param {string} name
 * @returns {{ el: Element, value: string }[]}
 */
export function takeMarked(root, name) {
  const roots = Array.isArray(root) ? root : [root]
  const out = []
  for (const r of roots) {
    if (r.nodeType === 1 && r.hasAttribute(name)) out.push(r)
    out.push(...r.querySelectorAll(`[${name}]`))
  }
  return out.map((el) => {
    const value = el.getAttribute(name)
    el.removeAttribute(name)
    return { el, value }
  })
}
