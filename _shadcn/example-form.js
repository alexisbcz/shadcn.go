// The form-* examples' behavior: each example's form library (react-hook-
// form, TanStack Form, Formisch) and schema, ported on form-state.js, and
// the React render of the fields' state (data-invalid, aria-invalid,
// FieldError, counters) patched into the server-rendered markup.
import { define } from "./core.js"
import { toast } from "./sonner.js"
import { setProgress } from "./progress.js"
import { getSelect } from "./select.js"
import { bindForm } from "./form.js"
import { createForm, rules, issuesOf, setAttr, syncText, fieldError, slot, submittedToast } from "./form-state.js"

const examples = {}

define("example-form", {
  selector: "template[data-form-example]",
  mount(tpl) {
    const name = tpl.getAttribute("data-form-example")
    const root = tpl.previousElementSibling
    const setup = examples[name]
    if (!root || !setup) return
    const lib = name.split("-")[1]
    return setup(root, lib, name)
  },
})

// --- shared wiring --------------------------------------------------------------------------------

/** Listens to events; returns the cleanup. */
function listeners() {
  const list = []
  const on = (el, type, fn, opts) => {
    if (!el) return
    el.addEventListener(type, fn, opts)
    list.push(() => el.removeEventListener(type, fn, opts))
  }
  return { on, cleanup: () => list.forEach((fn) => fn()) }
}

/** The submit and reset wiring shared by every example. */
function wireForm(root, formEl, form, on) {
  on(formEl, "submit", (event) => {
    event.preventDefault()
    form.submit()
  })
  const reset = [...root.querySelectorAll('[data-slot="card-footer"] button[type="button"]')].find((b) => b.textContent.trim() === "Reset")
  on(reset, "click", () => form.reset())
}

/**
 * A text control bound to a field key the way each library's example binds
 * it (onChange, onBlur; Formisch's props also register the element).
 */
function bindText(form, key, control, on) {
  form.register(key, control)
  on(control, "input", () => form.change(key, control.value))
  on(control, "focusout", () => form.blur(key))
}

/** Renders a field's invalid state on a Field wrapper and its control. */
function renderInvalid(state, ...els) {
  for (const el of els) {
    if (!el) continue
    if (el.hasAttribute("data-slot") && (el.getAttribute("data-slot") === "field" || el.getAttribute("data-slot") === "field-set")) setAttr(el, "data-invalid", state.invalid)
    else setAttr(el, "aria-invalid", state.invalid)
  }
}

/**
 * Sets an element's text nodes, in order: React renders {a}{b} as two
 * text nodes (the server markup separates them with a comment).
 */
function setTexts(el, ...parts) {
  const texts = [...el.childNodes].filter((n) => n.nodeType === 3)
  parts.forEach((p, i) => {
    const v = String(p)
    if (texts[i]) {
      if (texts[i].data !== v) texts[i].data = v
    } else el.appendChild(document.createTextNode(v))
  })
}

const onSubmitToast = (values) => submittedToast(toast, values)

// --- schemas --------------------------------------------------------------------------------------

const email = (lib) => (lib === "formisch" ? rules.valibotEmail : rules.zodEmail)

// --- input, textarea --------------------------------------------------------------------------------

/** A single text field example: a Field with label, control, description and error last. */
function singleText(key, check) {
  return (root, lib) => {
    const { on, cleanup } = listeners()
    const formEl = root.querySelector("form")
    const control = formEl.querySelector("input, textarea")
    const fieldEl = control.closest('[data-slot="field"]')
    const error = slot(fieldEl, fieldEl.querySelector('[data-slot="field-description"]'))
    const form = createForm({
      lib,
      defaults: { [key]: "" },
      schema: (v) => issuesOf([key], v[key], check),
      onSubmit: onSubmitToast,
      render(form) {
        const st = form.field(key)
        syncText(control, form.value(key))
        renderInvalid(st, fieldEl, control)
        error.set(st.invalid ? fieldError(st.errors) : null)
      },
      fields: [key],
    })
    bindText(form, key, control, on)
    wireForm(root, formEl, form, on)
    return cleanup
  }
}

const usernameRules = [
  rules.minLength(3, "Username must be at least 3 characters."),
  rules.maxLength(10, "Username must be at most 10 characters."),
  rules.regex(/^[a-zA-Z0-9_]+$/, "Username can only contain letters, numbers, and underscores."),
]
const aboutRules = [
  rules.minLength(10, "Please provide at least 10 characters."),
  rules.maxLength(200, "Please keep it under 200 characters."),
]
for (const lib of ["rhf", "tanstack", "formisch"]) {
  examples[`form-${lib}-input`] = singleText("username", usernameRules)
  examples[`form-${lib}-textarea`] = singleText("about", aboutRules)
}

// --- demo (bug report) ------------------------------------------------------------------------------

const titleRules = [
  rules.minLength(5, "Bug title must be at least 5 characters."),
  rules.maxLength(32, "Bug title must be at most 32 characters."),
]
const descriptionRules = [
  rules.minLength(20, "Description must be at least 20 characters."),
  rules.maxLength(100, "Description must be at most 100 characters."),
]

function demo(root, lib) {
  const { on, cleanup } = listeners()
  const formEl = root.querySelector("form")
  const [titleField, descField] = formEl.querySelectorAll('[data-slot="field"]')
  const title = titleField.querySelector("input")
  const desc = descField.querySelector("textarea")
  const counter = descField.querySelector('[data-slot="input-group-addon"] span')
  const titleError = slot(titleField, title)
  const descError = slot(descField, descField.querySelector('[data-slot="field-description"]'))
  const form = createForm({
    lib,
    defaults: { title: "", description: "" },
    schema: (v) => [...issuesOf(["title"], v.title, titleRules), ...issuesOf(["description"], v.description, descriptionRules)],
    onSubmit: onSubmitToast,
    fields: ["title", "description"],
    render(form) {
      const t = form.field("title")
      syncText(title, form.value("title"))
      renderInvalid(t, titleField, title)
      titleError.set(t.invalid ? fieldError(t.errors) : null)
      const d = form.field("description")
      syncText(desc, form.value("description"))
      renderInvalid(d, descField, desc)
      setTexts(counter, form.value("description").length, "/100 characters")
      descError.set(d.invalid ? fieldError(d.errors) : null)
    },
  })
  bindText(form, "title", title, on)
  bindText(form, "description", desc, on)
  wireForm(root, formEl, form, on)
  return cleanup
}
for (const lib of ["rhf", "tanstack", "formisch"]) examples[`form-${lib}-demo`] = demo

// --- switch, radio group, checkbox ------------------------------------------------------------------

// Controlled Radix controls: the vanilla components update themselves on
// user input; the form state drives them only on reset (programmatically,
// without reporting the change back).
let muted = false
let userChange = false

/** Binds a Radix checkbox/switch: onCheckedChange. */
function bindChecked(form, key, control, on, toValue = (checked) => checked) {
  on(control, "shadcn:checked-change", (event) => {
    if (muted) return
    const value = toValue(event.detail.checked, control)
    // React commits the controlled checkbox's new state and the form's
    // render together, the checkbox first (it is deeper in the tree): let
    // the component update itself, then render the form.
    queueMicrotask(() => {
      userChange = true
      try {
        form.change(key, value)
      } finally {
        userChange = false
      }
    })
  })
}

/** Drives a Radix checkbox/switch to checked (form reset). */
function syncChecked(control, checked) {
  if (userChange || !control) return
  if ((control.getAttribute("aria-checked") === "true") === !!checked) return
  muted = true
  try {
    control.click()
  } finally {
    muted = false
  }
}

/** Drives a Radix radio group to value (form reset); "" unchecks every item. */
function syncRadio(group, value) {
  if (userChange || !group) return
  const items = [...group.querySelectorAll('[data-slot="radio-group-item"]')]
  const current = items.find((i) => i.getAttribute("aria-checked") === "true")
  if ((current?.value ?? "") === value) return
  muted = true
  try {
    const target = items.find((i) => i.value === value)
    if (target) {
      target.click()
      return
    }
    // Radix RadioGroup with a value no item has: the checked item unchecks.
    if (!current) return
    current.setAttribute("aria-checked", "false")
    current.setAttribute("data-state", "unchecked")
    current.querySelector(':scope > [data-slot="radio-group-indicator"]')?.remove()
    const input = current.nextElementSibling
    if (input?.matches('input[type="radio"][aria-hidden="true"]')) {
      input.defaultChecked = false
      input.checked = false
      input.dispatchEvent(new Event("click", { bubbles: true }))
    }
  } finally {
    muted = false
  }
}

function bindRadio(form, key, group, on) {
  on(group, "shadcn:value-change", (event) => {
    if (muted) return
    userChange = true
    try {
      form.change(key, event.detail.value)
    } finally {
      userChange = false
    }
  })
}

const twoFactorMessage = "It is highly recommended to enable two-factor authentication."

function switchExample(root, lib) {
  const { on, cleanup } = listeners()
  const formEl = root.querySelector("form")
  const control = formEl.querySelector('[data-slot="switch"]')
  const fieldEl = control.closest('[data-slot="field"]')
  const content = fieldEl.querySelector('[data-slot="field-content"]')
  const error = slot(content, content.querySelector('[data-slot="field-description"]'))
  const form = createForm({
    lib,
    defaults: { twoFactor: false },
    schema: (v) => issuesOf(["twoFactor"], v.twoFactor, [rules.refine((x) => x === true, twoFactorMessage)]),
    onSubmit: onSubmitToast,
    fields: ["twoFactor"],
    render(form) {
      const st = form.field("twoFactor")
      syncChecked(control, form.value("twoFactor"))
      renderInvalid(st, fieldEl, control)
      error.set(st.invalid ? fieldError(st.errors) : null)
    },
  })
  bindChecked(form, "twoFactor", control, on)
  wireForm(root, formEl, form, on)
  return cleanup
}

const planMessage = "You must select a subscription plan to continue."

function radioExample(root, lib) {
  const { on, cleanup } = listeners()
  const formEl = root.querySelector("form")
  const group = formEl.querySelector('[data-slot="radio-group"]')
  const set = group.closest('[data-slot="field-set"]')
  const fields = [...group.querySelectorAll('[data-slot="field"]')]
  const items = [...group.querySelectorAll('[data-slot="radio-group-item"]')]
  const error = slot(set, group)
  const form = createForm({
    lib,
    defaults: { plan: "" },
    schema: (v) => issuesOf(["plan"], v.plan, [rules.minLength(1, planMessage)]),
    onSubmit: onSubmitToast,
    fields: ["plan"],
    render(form) {
      const st = form.field("plan")
      syncRadio(group, form.value("plan"))
      if (lib !== "tanstack") renderInvalid(st, set, group)
      renderInvalid(st, ...fields, ...items)
      error.set(st.invalid ? fieldError(st.errors) : null)
    },
  })
  bindRadio(form, "plan", group, on)
  wireForm(root, formEl, form, on)
  return cleanup
}

const taskIds = ["push", "email"]

function checkboxExample(root, lib) {
  const { on, cleanup } = listeners()
  const formEl = root.querySelector("form")
  const [responsesSet, tasksSet] = formEl.querySelectorAll('[data-slot="field-set"]')
  const responses = responsesSet.querySelector('[data-slot="checkbox"]')
  const responsesField = responses.closest('[data-slot="field"]')
  const responsesError = slot(responsesSet.parentElement, responsesSet)
  const taskFields = [...tasksSet.querySelectorAll('[data-slot="field"]')]
  const taskBoxes = taskFields.map((f) => f.querySelector('[data-slot="checkbox"]'))
  const tasksError = slot(tasksSet.parentElement, tasksSet)
  const form = createForm({
    lib,
    defaults: { responses: true, tasks: [] },
    schema: (v) => [
      ...issuesOf(["tasks"], v.tasks, [
        rules.minLength(1, "Please select at least one notification type."),
        rules.refine((value) => value.every((task) => taskIds.includes(task)), "Invalid notification type selected."),
      ]),
    ],
    onSubmit: onSubmitToast,
    fields: ["responses", "tasks"],
    render(form) {
      const r = form.field("responses")
      syncChecked(responses, form.value("responses"))
      renderInvalid(r, lib === "tanstack" ? responsesField : responsesSet)
      responsesError.set(r.invalid ? fieldError(r.errors) : null)
      const t = form.field("tasks")
      const value = form.value("tasks")
      taskBoxes.forEach((box, i) => syncChecked(box, value.includes(taskIds[i])))
      renderInvalid(t, tasksSet, ...taskFields, ...taskBoxes)
      tasksError.set(t.invalid ? fieldError(t.errors) : null)
    },
  })
  bindChecked(form, "responses", responses, on, (checked) => checked === true)
  taskBoxes.forEach((box, i) =>
    bindChecked(form, "tasks", box, on, (checked) => {
      const current = form.value("tasks")
      return checked ? [...current, taskIds[i]] : current.filter((v) => v !== taskIds[i])
    }),
  )
  wireForm(root, formEl, form, on)
  return cleanup
}

/** Binds a Radix select (its template dispatches the value changes). */
function bindSelect(form, key, trigger, on) {
  on(trigger.parentElement, "shadcn:value-change", (event) => {
    if (muted || !(event.target instanceof HTMLTemplateElement)) return
    if (getSelect(trigger) && event.target !== document.querySelector(`template[data-content-id="${trigger.getAttribute("aria-controls")}"]`)) return
    userChange = true
    try {
      form.change(key, event.detail.value)
    } finally {
      userChange = false
    }
  })
}

/** Drives a Radix select to value (form reset). */
function syncSelect(trigger, value) {
  if (userChange) return
  const sel = getSelect(trigger)
  if (!sel || (sel.value ?? "") === value) return
  muted = true
  try {
    sel.setValue(value)
  } finally {
    muted = false
  }
}

const languageMessages = ["Please select your spoken language.", "Auto-detection is not allowed. Please select a specific language."]

function selectExample(root, lib) {
  const { on, cleanup } = listeners()
  const formEl = root.querySelector("form")
  const trigger = formEl.querySelector('[data-slot="select-trigger"]')
  const fieldEl = trigger.closest('[data-slot="field"]')
  const content = fieldEl.querySelector('[data-slot="field-content"]')
  const error = slot(content, content.querySelector('[data-slot="field-description"]'))
  const form = createForm({
    lib,
    defaults: { language: "" },
    schema: (v) =>
      issuesOf(["language"], v.language, [rules.minLength(1, languageMessages[0]), rules.refine((x) => x !== "auto", languageMessages[1])]),
    onSubmit: onSubmitToast,
    fields: ["language"],
    render(form) {
      const st = form.field("language")
      syncSelect(trigger, form.value("language"))
      renderInvalid(st, fieldEl, trigger)
      error.set(st.invalid ? fieldError(st.errors) : null)
    },
  })
  bindSelect(form, "language", trigger, on)
  wireForm(root, formEl, form, on)
  return cleanup
}

for (const lib of ["rhf", "tanstack", "formisch"]) {
  examples[`form-${lib}-select`] = selectExample
  examples[`form-${lib}-switch`] = switchExample
  examples[`form-${lib}-radiogroup`] = radioExample
  examples[`form-${lib}-checkbox`] = checkboxExample
}

// --- password ---------------------------------------------------------------------------------------

const passwordTests = [
  (v) => v.length >= 8,
  (v) => /[a-z]/.test(v),
  (v) => /[A-Z]/.test(v),
  (v) => /\d/.test(v),
  (v) => /[!@#$%^&*(),.?":{}|<>]/.test(v),
]

examples["form-rhf-password"] = (root, lib, name) => {
  const { on, cleanup } = listeners()
  const classes = JSON.parse(root.nextElementSibling?.getAttribute("data-classes") || "{}")
  const formEl = root.querySelector("form")
  const control = formEl.querySelector("input")
  const fieldEl = control.closest('[data-slot="field"]')
  const addonIcon = fieldEl.querySelector('[data-slot="input-group-addon"] svg')
  const meter = fieldEl.querySelector('[data-slot="progress"]')
  const meterBox = meter.parentElement
  const rows = [...meterBox.querySelectorAll(".space-y-1\\.5 > div")]
  const error = slot(fieldEl, meterBox)
  meter.__pct = 0
  const form = createForm({
    lib,
    defaults: { password: "" },
    schema: (v) =>
      issuesOf(["password"], v.password, [
        rules.minLength(8, "Password must be at least 8 characters"),
        rules.refine((val) => /[a-z]/.test(val), "Password must contain at least one lowercase letter"),
        rules.refine((val) => /[A-Z]/.test(val), "Password must contain at least one uppercase letter"),
        rules.refine((val) => /\d/.test(val), "Password must contain at least one number"),
        rules.refine((val) => /[!@#$%^&*(),.?":{}|<>]/.test(val), "Password must contain at least one special character"),
      ]),
    onSubmit: onSubmitToast,
    render(form) {
      const password = form.value("password") || ""
      const st = form.field("password")
      syncText(control, password)
      renderInvalid(st, fieldEl, control)
      const met = passwordTests.map((t) => t(password))
      const count = met.filter(Boolean).length
      const pct = (count / passwordTests.length) * 100
      const color = pct === 0 ? "bg-neutral-200" : pct <= 40 ? "bg-red-500" : pct <= 80 ? "bg-yellow-500" : "bg-green-500"
      setAttr(addonIcon, "class", count === passwordTests.length ? classes.addonMet : classes.addonUnmet)
      setAttr(meter, "class", classes[color])
      // Only on change: a write repaints even when the value is the same.
      if (meter.__pct !== pct) {
        meter.__pct = pct
        setProgress(meter, pct)
      }
      rows.forEach((row, i) => {
        setAttr(row.querySelector("svg"), "class", met[i] ? classes.requirementMet : classes.requirementUnmet)
        setAttr(row.querySelector("span"), "class", met[i] ? "text-foreground" : "text-muted-foreground")
      })
      error.set(st.invalid ? fieldError(st.errors) : null)
    },
  })
  bindText(form, "password", control, on)
  wireForm(root, formEl, form, on)
  return cleanup
}

// --- array (contact emails) -------------------------------------------------------------------------

function arrayExample(root, lib, name) {
  const { on, cleanup } = listeners()
  const tpl = root.nextElementSibling
  const itemProto = tpl.content.querySelector('[data-proto="item"] > [data-slot="field"]')
  const removeProto = tpl.content.querySelector('[data-proto="remove"] > [data-slot="input-group-addon"]')
  const formEl = root.querySelector("form")
  const group = formEl.querySelector('[data-slot="field-set"] > [data-slot="field-group"]')
  const addButton = group.querySelector(':scope > [data-slot="button"]')
  const keyed = lib !== "tanstack" // React keys: item ids (rhf, formisch) or indexes (tanstack)
  const emailRules = lib === "formisch" ? [rules.minLength(1, "Enter an email address."), rules.valibotEmail("Enter a valid email address.")] : [rules.zodEmail("Enter a valid email address.")]
  const nameOf = (i) => (lib === "tanstack" ? `emails[${i}].address` : lib === "formisch" ? JSON.stringify(["emails", i, "address"]) : `emails.${i}.address`)
  const keyOf = (i) => `emails.${i}.address`

  /** @type {{ el: HTMLElement, input: HTMLInputElement, inputGroup: HTMLElement, error: any, addon: HTMLElement | null }[]} */
  let items = []
  const adopt = (el) => {
    const input = el.querySelector("input")
    const inputGroup = el.querySelector('[data-slot="input-group"]')
    const item = { el, input, inputGroup, error: slot(el.querySelector('[data-slot="field-content"]'), inputGroup), addon: inputGroup.querySelector('[data-slot="input-group-addon"]') }
    on(input, "input", () => form.change(keyOf(items.indexOf(item)), input.value))
    on(input, "focusout", () => form.blur(keyOf(items.indexOf(item))))
    on(el, "click", (event) => {
      const button = event.target instanceof Element && event.target.closest('[data-slot="input-group-addon"] button')
      if (button && el.contains(button)) form.arrayRemove("emails", items.indexOf(item))
    })
    return item
  }
  const create = () => {
    const el = itemProto.cloneNode(true)
    return adopt(el)
  }
  items = [...group.querySelectorAll(':scope > [data-slot="field"]')].map(adopt)

  const syncItems = (count) => {
    while (items.length > count) {
      const removed = items.pop()
      removed.el.remove()
    }
    while (items.length < count) {
      const item = create()
      item.id = nextId++
      group.insertBefore(item.el, addButton)
      items.push(item)
    }
  }

  const form = createForm({
    lib,
    defaults: { emails: lib === "tanstack" ? [{ address: "" }] : [{ address: "" }, { address: "" }] },
    schema: (v) => {
      const out = []
      v.emails.forEach((e, i) => out.push(...issuesOf(["emails", i, "address"], e.address, emailRules)))
      out.push(...issuesOf(["emails"], v.emails, [rules.minLength(1, "Add at least one email address."), rules.maxLength(5, "You can add up to 5 email addresses.")]))
      return lib === "rhf" || lib === "tanstack" ? [...out.filter((i) => i.path.length === 1), ...out.filter((i) => i.path.length > 1)] : out
    },
    fields: () => form.values.emails.map((_, i) => keyOf(i)),
    tanstackOn: "blur",
    onSubmit: onSubmitToast,
    render(form) {
      const emails = form.value("emails")
      syncItems(emails.length)
      items.forEach((item, i) => {
        const st = form.field(keyOf(i))
        setAttr(item.input, "id", `${name}-email-${i}`)
        setAttr(item.input, "name", nameOf(i))
        syncText(item.input, emails[i].address)
        renderInvalid(st, item.el, item.input)
        if (emails.length > 1) {
          if (!item.addon) {
            item.addon = removeProto.cloneNode(true)
            item.inputGroup.appendChild(item.addon)
          }
          setAttr(item.addon.querySelector("button"), "aria-label", `Remove email ${i + 1}`)
        } else if (item.addon) {
          item.addon.remove()
          item.addon = null
        }
        item.error.set(st.invalid ? fieldError(st.errors) : null)
        form.register(keyOf(i), item.input)
      })
      if (emails.length >= 5) addButton.setAttribute("disabled", "")
      else addButton.removeAttribute("disabled")
    },
  })

  items.forEach((item, i) => form.register(keyOf(i), item.input))
  if (keyed) {
    // Keyed items: the removed item's own node leaves; a reset gives every
    // item a new key (react-hook-form) or its initial one (Formisch).
    const remove = form.arrayRemove
    form.arrayRemove = (key, index) => {
      const [item] = items.splice(index, 1)
      item.el.remove()
      remove(key, index)
    }
  }
  // Item keys: react-hook-form gives every item a new key on reset (all
  // nodes are re-created); Formisch resets to the initial items' keys, so
  // the nodes of initial items that were not removed are kept.
  let nextId = 0
  const initialIds = items.map(() => nextId++)
  items.forEach((item, i) => (item.id = initialIds[i]))
  const reset = form.reset
  form.reset = () => {
    if (lib === "rhf") syncItems(0)
    else if (lib === "formisch") {
      const kept = initialIds.map((id) => items.find((item) => item.id === id) ?? Object.assign(create(), { id }))
      for (const item of items) if (!kept.includes(item)) item.el.remove()
      // Like React's reconciliation: move or insert only what is out of place.
      let ref = addButton
      for (let i = kept.length - 1; i >= 0; i--) {
        const el = kept[i].el
        if (el.nextSibling !== ref || !el.isConnected) group.insertBefore(el, ref)
        ref = el
      }
      items = kept
    }
    reset()
  }

  on(addButton, "click", () => {
    const index = form.arrayAppend("emails", { address: "" })
    // react-hook-form's append focuses the new field.
    if (lib === "rhf") items[index]?.input.focus()
  })
  wireForm(root, formEl, form, on)
  return cleanup
}
for (const lib of ["rhf", "tanstack", "formisch"]) examples[`form-${lib}-array`] = arrayExample

// --- Next.js forms --------------------------------------------------------------------------------

// form-next-demo: the description counter (onChange). The form itself
// submits natively (see next.go).
examples["form-next-demo"] = (root) => {
  const { on, cleanup } = listeners()
  const desc = root.querySelector("textarea")
  const counter = root.querySelector('[data-slot="input-group-addon"] span')
  on(desc, "input", () => setTexts(counter, desc.value.length, "/100 characters"))
  return cleanup
}

// --- complex (subscription) -----------------------------------------------------------------------

const addonIds = ["analytics", "backup", "support"]

function complexExample(root, lib) {
  const { on, cleanup } = listeners()
  const formEl = root.querySelector("form")
  const tanstack = lib === "tanstack"
  const [planSet, addonsSet] = formEl.querySelectorAll('[data-slot="field-set"]')
  const group = planSet.querySelector('[data-slot="radio-group"]')
  const planFields = [...group.querySelectorAll('[data-slot="field"]')]
  const planItems = [...group.querySelectorAll('[data-slot="radio-group-item"]')]
  const planError = slot(planSet, group)
  const trigger = formEl.querySelector('[data-slot="select-trigger"]')
  const billingField = trigger.closest('[data-slot="field"]')
  const billingError = slot(billingField, billingField.querySelector(':scope > [data-slot="field-description"]'))
  const addonFields = [...addonsSet.querySelectorAll('[data-slot="field"]')]
  const addonBoxes = addonFields.map((f) => f.querySelector('[data-slot="checkbox"]'))
  const addonsError = slot(addonsSet, addonsSet.querySelector('[data-slot="checkbox-group"]'))
  const sw = formEl.querySelector('[data-slot="switch"]')
  const notifyField = sw.closest('[data-slot="field"]')
  const notifyError = slot(notifyField, sw.nextElementSibling?.matches('input[aria-hidden="true"]') ? sw.nextElementSibling : sw)

  const planRules = lib === "formisch"
    ? [rules.minLength(1, "Please select a subscription plan"), rules.refine((v) => v === "basic" || v === "pro", "Invalid plan selection. Please choose Basic or Pro")]
    : [rules.minLength(1, "Please select a subscription plan"), rules.refine((v) => v === "basic" || v === "pro", "Invalid plan selection. Please choose Basic or Pro")]
  const form = createForm({
    lib,
    defaults: { plan: "basic", billingPeriod: tanstack ? "monthly" : "", addons: [], emailNotifications: false },
    schema: (v) => [
      ...issuesOf(["plan"], v.plan, planRules),
      ...issuesOf(["billingPeriod"], v.billingPeriod, [rules.minLength(1, "Please select a billing period")]),
      ...issuesOf(["addons"], v.addons, [
        rules.minLength(1, "Please select at least one add-on"),
        rules.maxLength(3, "You can select up to 3 add-ons"),
        rules.refine((value) => value.every((a) => addonIds.includes(a)), "You selected an invalid add-on"),
      ]),
    ],
    onSubmit: onSubmitToast,
    fields: ["plan", "billingPeriod", "addons", "emailNotifications"],
    render(form) {
      const plan = form.field("plan")
      syncRadio(group, form.value("plan"))
      if (tanstack) renderInvalid(plan, ...planFields, ...planItems)
      else renderInvalid(plan, planSet, group)
      planError.set(plan.invalid ? fieldError(plan.errors) : null)

      const billing = form.field("billingPeriod")
      syncSelect(trigger, form.value("billingPeriod"))
      renderInvalid(billing, billingField)
      if (!tanstack) renderInvalid(billing, trigger)
      billingError.set(billing.invalid ? fieldError(billing.errors) : null)

      const ad = form.field("addons")
      const value = form.value("addons")
      addonBoxes.forEach((box, i) => syncChecked(box, value.includes(addonIds[i])))
      // React's commit order: each checkbox, then its Field.
      addonFields.forEach((f, i) => renderInvalid(ad, addonBoxes[i], f))
      addonsError.set(ad.invalid ? fieldError(ad.errors) : null)

      const n = form.field("emailNotifications")
      syncChecked(sw, form.value("emailNotifications"))
      renderInvalid(n, notifyField, sw)
      notifyError.set(n.invalid ? fieldError(n.errors) : null)
    },
  })
  bindRadio(form, "plan", group, on)
  bindSelect(form, "billingPeriod", trigger, on)
  addonBoxes.forEach((box, i) =>
    bindChecked(form, "addons", box, on, (checked) => {
      const current = form.value("addons")
      return checked ? [...current, addonIds[i]] : current.filter((v) => v !== addonIds[i])
    }),
  )
  bindChecked(form, "emailNotifications", sw, on)
  wireForm(root, formEl, form, on)
  return cleanup
}
for (const lib of ["rhf", "tanstack", "formisch"]) examples[`form-${lib}-complex`] = complexExample

// --- form-hook-demo (the Form component) -------------------------------------------------------------

examples["form-hook-demo"] = (root) => {
  const form = bindForm(root, {
    lib: "rhf",
    defaults: { username: "", bio: "" },
    schema: (v) => [
      ...issuesOf(["username"], v.username, [
        rules.minLength(2, "Username must be at least 2 characters."),
        rules.maxLength(10, "Username must be at most 10 characters."),
      ]),
      ...issuesOf(["bio"], v.bio, [rules.maxLength(20, "Bio must be at most 20 characters.")]),
    ],
    onSubmit: () => {},
  })
  return () => form.destroy()
}
