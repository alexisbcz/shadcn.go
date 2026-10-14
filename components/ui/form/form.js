// form.js: upstream form.tsx's client side (useFormField): renders a form
// item's error state, and binds a form of items to form-state.js with
// react-hook-form's semantics (zodResolver-style schema).
//
//   setItemError(item, "Username must be at least 2 characters.") // or null
//   bindForm(formElement, { defaults, schema, onSubmit, lib: "rhf" })
//
// An item's control is the element with data-slot="form-control"; its
// field name is the control's name attribute.
import { createForm, getIn } from "./form-state.js"

// A FormMessage's own children (shown when the field has no error), keyed
// by item, recorded before the first error replaces them.
const ownChildren = new WeakMap()

/** Renders a FormItem's error state like useFormField: label, control, message. */
export function setItemError(item, message) {
  const error = message ? String(message) : ""
  const label = item.querySelector('[data-slot="form-label"]')
  const control = item.querySelector('[data-slot="form-control"]')
  const set = (el, name, v) => el.getAttribute(name) !== v && el.setAttribute(name, v)
  if (label) set(label, "data-error", String(!!error))
  if (!control) return
  const id = control.id
  const descriptionId = `${id}-description`
  const messageId = `${id}-message`
  set(control, "aria-invalid", String(!!error))
  set(control, "aria-describedby", error ? `${descriptionId} ${messageId}` : descriptionId)
  let msg = item.querySelector('[data-slot="form-message"]')
  if (!ownChildren.has(item)) ownChildren.set(item, msg ? [...msg.childNodes].map((n) => n.cloneNode(true)) : [])
  const own = ownChildren.get(item)
  if (error) {
    if (!msg) {
      msg = document.createElement("p")
      msg.setAttribute("data-slot", "form-message")
      msg.id = messageId
      msg.className = "text-sm text-destructive"
      item.appendChild(msg)
    }
    if (msg.textContent !== error || msg.childNodes.length !== 1) msg.textContent = error
  } else if (msg) {
    if (own.length) {
      const same = msg.childNodes.length === own.length && own.every((n, i) => n.isEqualNode(msg.childNodes[i]))
      if (!same) msg.replaceChildren(...own.map((n) => n.cloneNode(true)))
    }
    else msg.remove()
  }
}

/**
 * Binds a form of FormItems: each item's control (by its name) is a field
 * of a form-state.js form ("rhf" semantics by default).
 */
export function bindForm(formEl, opts) {
  const items = [...formEl.querySelectorAll('[data-slot="form-item"]')]
  const fields = items
    .map((item) => ({ item, control: item.querySelector('[data-slot="form-control"]') }))
    .filter((f) => f.control?.getAttribute("name"))
  const form = createForm({
    lib: opts.lib ?? "rhf",
    defaults: opts.defaults,
    schema: opts.schema,
    onSubmit: opts.onSubmit,
    render(form) {
      for (const { item, control } of fields) {
        const key = control.getAttribute("name")
        const st = form.field(key)
        setItemError(item, st.invalid ? st.errors[0]?.message ?? "" : null)
        const v = getIn(form.values, key)
        if (typeof v === "string") {
          // React-controlled: the value, and the value attribute it syncs.
          if (control.value !== v) control.value = v
          if (control.defaultValue !== v) control.defaultValue = v
        }
      }
      opts.render?.(form)
    },
  })
  const cleanups = []
  for (const { control } of fields) {
    const key = control.getAttribute("name")
    form.register(key, control)
    const onInput = () => form.change(key, control.value)
    const onBlur = () => form.blur(key)
    control.addEventListener("input", onInput)
    control.addEventListener("focusout", onBlur)
    cleanups.push(() => {
      control.removeEventListener("input", onInput)
      control.removeEventListener("focusout", onBlur)
    })
  }
  const onSubmit = (event) => {
    event.preventDefault()
    form.submit()
  }
  formEl.addEventListener("submit", onSubmit)
  cleanups.push(() => formEl.removeEventListener("submit", onSubmit))
  form.destroy = () => cleanups.forEach((fn) => fn())
  return form
}
