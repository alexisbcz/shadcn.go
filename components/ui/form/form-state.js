// form-state.js: client-side form state and validation with the semantics
// of the form libraries the shadcn/ui form examples use, so a server-
// rendered form behaves exactly like its React counterpart:
//
//   - "rhf": react-hook-form 7.62 (mode "onSubmit", reValidateMode
//     "onChange", criteriaMode "firstError", shouldFocusError) with a
//     schema resolver: nothing is validated before the first submit; a
//     submit validates everything, keeps the first issue per field and
//     focuses the first invalid field (in registration order); after a
//     submit, every change revalidates the changed field.
//   - "tanstack": @tanstack/form-core 1.20 with a form-level onSubmit
//     schema validator: a submit marks every field touched and, if the form
//     can submit, validates; each field gets all of its issues. A change or
//     blur clears the field's submit errors (and the form's). A form with
//     remaining errors cannot submit (canSubmit) until they are cleared.
//   - "formisch": @formisch/react 0.4 (validate "submit", revalidate
//     "input"): a submit validates everything and focuses the first invalid
//     field; after a submit, every input revalidates the whole form; each
//     field gets all of its messages.
//
// A schema is a function (values) => issues, issues being
// [{ path: ["emails", 0, "address"], message }] in the order zod (v3) or
// valibot report them; the rules below help writing them.
//
//   import { createForm, rules, issuesOf } from "./form-state.js"
//   const form = createForm({
//     lib: "rhf",
//     defaults: { username: "" },
//     schema: (v) => issuesOf(["username"], v.username, [rules.minLength(3, "Too short.")]),
//     onSubmit: (values) => toast(...),
//     render: (form) => { ... patch the DOM from form.field("username") ... },
//   })
//   form.register("username", inputElement) // focus target, in field order
//   input.addEventListener("input", () => form.change("username", input.value))
//
// Field keys are dotted paths ("emails.0.address").

/** Returns the dotted key of a path. */
export const keyOf = (path) => (Array.isArray(path) ? path.join(".") : String(path))

const clone = (v) => (v === undefined ? v : JSON.parse(JSON.stringify(v)))

/** Reads a dotted key from an object. */
export function getIn(obj, key) {
  if (key === "") return obj
  let cur = obj
  for (const part of String(key).split(".")) {
    if (cur == null) return undefined
    cur = cur[part]
  }
  return cur
}

/** Writes a dotted key into an object (in place). */
export function setIn(obj, key, value) {
  const parts = String(key).split(".")
  let cur = obj
  for (let i = 0; i < parts.length - 1; i++) {
    const p = parts[i]
    if (cur[p] == null) cur[p] = /^\d+$/.test(parts[i + 1]) ? [] : {}
    cur = cur[p]
  }
  cur[parts[parts.length - 1]] = value
}

// --- schema helpers -------------------------------------------------------------------------------

/**
 * Rules: each returns a check (value) => message | null. They run in
 * order and every failing one is reported, like zod's string/array checks
 * followed by refinements and valibot's pipes (abortPipeEarly false).
 */
export const rules = {
  minLength: (n, message) => (v) => (v.length < n ? message : null),
  maxLength: (n, message) => (v) => (v.length > n ? message : null),
  regex: (re, message) => (v) => {
    re.lastIndex = 0
    return re.test(v) ? null : message
  },
  // zod 3.25 z.string().email()
  zodEmail: (message) => (v) => (/^(?!\.)(?!.*\.\.)([A-Z0-9_'+\-\.]*)[A-Z0-9_+-]@([A-Z0-9][A-Z0-9\-]*\.)+[A-Z]{2,}$/i.test(v) ? null : message),
  // valibot 1.4 v.email()
  valibotEmail: (message) => (v) => (/^[\w+-]+(?:\.[\w+-]+)*@[\da-z]+(?:[.-][\da-z]+)*\.[a-z]{2,}$/iu.test(v) ? null : message),
  refine: (test, message) => (v) => (test(v) ? null : message),
}

/** Runs rules on a value and returns its issues at path. */
export function issuesOf(path, value, list) {
  const out = []
  for (const rule of list) {
    const message = rule(value)
    if (message != null) out.push({ path, message })
  }
  return out
}

// --- the form -------------------------------------------------------------------------------------

/**
 * Creates a form. Options:
 *   lib       "rhf" | "tanstack" | "formisch"
 *   defaults  the default (initial) values
 *   schema    (values) => issues
 *   onSubmit  (output, form) => void, called with the values when valid
 *   output    optional (values) => output, the parsed output passed to
 *             onSubmit (default: a copy of the values)
 *   render    (form) => void, called after every state change
 *   fields    optional list of the field keys, in the order the library
 *             walks them (TanStack: the mounted fields)
 */
export function createForm(opts) {
  const lib = opts.lib
  const order = [] // registration order: [key, element]
  const state = {
    values: clone(opts.defaults),
    submitted: false, // rhf, formisch: isSubmitted
    attempts: 0, // tanstack: submissionAttempts
    errors: {}, // rhf: key -> issue; formisch: key -> [message]; tanstack: key -> [issue]
    touched: {}, // tanstack
    formError: false, // tanstack: the form-level onSubmit error
  }
  let rendering = false

  const render = () => {
    if (rendering) return
    rendering = true
    try {
      opts.render?.(form)
    } finally {
      rendering = false
    }
  }

  const grouped = (issues) => {
    const map = {}
    for (const issue of issues) (map[keyOf(issue.path)] ??= []).push(issue)
    return map
  }

  const fieldKeys = () => (typeof opts.fields === "function" ? opts.fields() : opts.fields) ?? order.map(([k]) => k)
  // TanStack: when the form-level schema runs ("submit" or "blur").
  const tanstackOn = opts.tanstackOn ?? "submit"

  // TanStack: runs the form-level schema and gives every mounted field its
  // issues (validateSync), returns whether there were none.
  const tanstackValidate = () => {
    const map = grouped(opts.schema(state.values))
    const errors = {}
    for (const k of fieldKeys()) if (map[k]) errors[k] = map[k]
    state.errors = errors
    state.formError = Object.keys(map).length > 0
    return !state.formError
  }

  const focusFirstError = () => {
    if (lib === "rhf") {
      for (const k of fieldKeys()) {
        const el = order.find(([key]) => key === k)?.[1]
        if (state.errors[k] && el) {
          el.focus()
          return
        }
      }
    } else if (lib === "formisch") {
      for (const k of fieldKeys()) {
        if (!state.errors[k]) continue
        const el = order.find(([key]) => key === k)?.[1]
        if (el) {
          el.focus()
          return
        }
        // Formisch focuses the first invalid field only (if it has an element).
        return
      }
    }
  }

  const output = () => (opts.output ? opts.output(clone(state.values)) : clone(state.values))

  // Formisch: validateFormInput, every field's errors from scratch.
  const formischValidate = () => {
    const map = grouped(opts.schema(state.values))
    const errors = {}
    for (const [k, list] of Object.entries(map)) errors[k] = list.map((i) => i.message)
    state.errors = errors
    return Object.keys(map).length === 0
  }

  // react-hook-form: a field's error after a change (onChange with a resolver).
  const rhfRevalidate = (k) => {
    const map = grouped(opts.schema(state.values))
    // schemaErrorLookup: the error at the key, or at the nearest parent
    // that has one (field arrays).
    let name = k
    let issue = map[k]?.[0]
    if (!issue) {
      const parts = k.split(".")
      while (parts.length > 1) {
        parts.pop()
        const parent = parts.join(".")
        if (order.some(([key]) => key === parent)) break
        if (map[parent]) {
          name = parent
          issue = map[parent][0]
          break
        }
      }
    }
    if (issue) state.errors[name] = issue
    else delete state.errors[name]
  }

  const form = {
    lib,
    state,
    get values() {
      return state.values
    },
    value: (k) => getIn(state.values, k),

    /** Registers a field's focus target (react-hook-form's ref, formisch's props.ref). */
    register(k, el) {
      const i = order.findIndex(([key]) => key === k)
      if (i >= 0) order[i][1] = el
      else order.push([k, el])
    },
    unregister(k) {
      const i = order.findIndex(([key]) => key === k)
      if (i >= 0) order.splice(i, 1)
    },

    /** The field's display state: { invalid, errors: [{ message }] }. */
    field(k) {
      if (lib === "rhf") {
        const e = state.errors[k]
        return { invalid: !!e, errors: e ? [e] : [] }
      }
      if (lib === "tanstack") {
        const errs = state.errors[k] ?? []
        const invalid = !!state.touched[k] && errs.length > 0
        return { invalid, errors: errs, touched: !!state.touched[k], valid: errs.length === 0 }
      }
      const msgs = state.errors[k] ?? null
      return { invalid: msgs !== null, errors: msgs ? msgs.map((message) => ({ message })) : [] }
    },

    /** A field's value changed (user input). */
    change(k, value) {
      setIn(state.values, k, value)
      if (lib === "rhf") {
        if (state.submitted) rhfRevalidate(k)
      } else if (lib === "tanstack") {
        state.touched[k] = true
        if (tanstackOn === "submit") {
          delete state.errors[k]
          state.formError = false
        }
      } else if (state.submitted) {
        formischValidate()
      }
      render()
    },

    /** A field lost focus. */
    blur(k) {
      if (lib === "tanstack") {
        state.touched[k] = true
        if (tanstackOn === "blur") tanstackValidate()
        else {
          delete state.errors[k]
          state.formError = false
        }
        render()
      }
    },

    /** Replaces values outside of a field change (field arrays); revalidate: see each library. */
    setValues(fn) {
      fn(state.values)
      render()
    },

    /**
     * Removes the item at index of the array at key (react-hook-form's
     * remove, TanStack's removeValue, Formisch's remove): later items'
     * state shifts down.
     */
    arrayRemove(key, index) {
      const arr = getIn(state.values, key)
      const len = arr.length
      arr.splice(index, 1)
      const shift = (map) => {
        const out = {}
        for (const [k, v] of Object.entries(map)) {
          if (!k.startsWith(key + ".")) {
            out[k] = v
            continue
          }
          const rest = k.slice(key.length + 1)
          const dot = rest.indexOf(".")
          const i = Number(dot < 0 ? rest : rest.slice(0, dot))
          const tail = dot < 0 ? "" : rest.slice(dot)
          if (i < index) out[k] = v
          else if (i > index) out[`${key}.${i - 1}${tail}`] = v
          // Formisch keeps the last item's stale state (it is not walked).
          else if (lib === "formisch" && i === len - 1) out[k] = v
        }
        return out
      }
      state.errors = shift(state.errors)
      state.touched = shift(state.touched)
      if (lib === "tanstack") state.touched[key] = true
      if (lib === "formisch") {
        for (const k of Object.keys(state.errors)) {
          const m = k.slice(key.length + 1).match(/^(\d+)/)
          if (k.startsWith(key + ".") && m && Number(m[1]) >= len - 1) delete state.errors[k]
        }
        if (state.submitted) formischValidate()
      }
      render()
    },

    /** Appends item to the array at key (append, pushValue, insert). */
    arrayAppend(key, item) {
      const arr = getIn(state.values, key)
      arr.push(clone(item))
      const index = arr.length - 1
      for (const k of Object.keys(state.errors)) if (k.startsWith(`${key}.${index}.`) || k === `${key}.${index}`) delete state.errors[k]
      if (lib === "tanstack") state.touched[key] = true
      if (lib === "formisch" && state.submitted) formischValidate()
      render()
      return index
    },

    /** Submits: validates and calls onSubmit when valid. */
    submit() {
      if (lib === "rhf") {
        const map = grouped(opts.schema(state.values))
        const errors = {}
        for (const [k, list] of Object.entries(map)) errors[k] = list[0]
        state.errors = errors
        const valid = Object.keys(errors).length === 0
        // handleSubmit focuses the first invalid field before React renders
        // the errors.
        if (!valid) focusFirstError()
        if (valid) opts.onSubmit?.(output(), form)
        state.submitted = true
        render()
        return
      }
      if (lib === "tanstack") {
        state.attempts++
        for (const k of fieldKeys()) state.touched[k] = true
        const isValid = !state.formError && fieldKeys().every((k) => !(state.errors[k]?.length))
        if (!isValid) {
          render()
          return
        }
        const valid = tanstackValidate()
        render()
        if (valid) opts.onSubmit?.(output(), form)
        return
      }
      state.submitted = true
      const valid = formischValidate()
      // validateFormInput focuses inside its batch, before React renders.
      focusFirstError()
      render()
      if (valid) opts.onSubmit?.(output(), form)
    },

    /** Resets values and state to the defaults. */
    reset() {
      state.values = clone(opts.defaults)
      state.submitted = false
      state.attempts = 0
      state.errors = {}
      state.touched = {}
      state.formError = false
      opts.onReset?.(form)
      render()
    },

    render,
  }
  return form
}

// --- DOM helpers ----------------------------------------------------------------------------------

/** Sets an attribute to String(value) (React's rendering of booleans in data-* and aria-*). */
export function setAttr(el, name, value) {
  if (!el) return
  if (value === undefined || value === null) {
    if (el.hasAttribute(name)) el.removeAttribute(name)
    return
  }
  const v = String(value)
  if (el.getAttribute(name) !== v) el.setAttribute(name, v)
}

/**
 * Syncs a React-controlled text control to value: its value, and the value
 * attribute (input) or default value (textarea) that React keeps in sync.
 */
export function syncText(el, value) {
  if (!el) return
  const v = value == null ? "" : String(value)
  if (el.value !== v) el.value = v
  if (el.defaultValue !== v) el.defaultValue = v
}

/** Upstream FieldError's content for errors: deduplicated messages, a list if several. */
export function fieldError(errors, className) {
  if (!errors?.length) return null
  const unique = new Map()
  for (const e of errors) unique.set(e?.message, e)
  const list = [...unique.values()]
  const div = document.createElement("div")
  div.setAttribute("role", "alert")
  div.setAttribute("data-slot", "field-error")
  div.className = className ? `text-sm font-normal text-destructive ${className}` : "text-sm font-normal text-destructive"
  if (list.length === 1) {
    if (!list[0]?.message) return null
    div.textContent = list[0].message
  } else {
    const ul = document.createElement("ul")
    ul.className = "ml-4 flex list-disc flex-col gap-1"
    for (const e of list) {
      if (!e?.message) continue
      const li = document.createElement("li")
      li.textContent = e.message
      ul.appendChild(li)
    }
    div.appendChild(ul)
  }
  return div
}

/**
 * A place in the DOM where an element is conditionally rendered, like
 * {cond && <X />} in JSX: after `after` (an always-present sibling) inside
 * parent, or first if after is null.
 */
export function slot(parent, after = null) {
  let node = null
  return {
    get node() {
      return node
    },
    /** Renders next (an element or null). Equal content keeps the existing node. */
    set(next) {
      if (node && next && node.isEqualNode(next)) return
      if (node) node.remove()
      node = next
      if (!next) return
      const ref = after ? after.nextSibling : parent.firstChild
      parent.insertBefore(next, ref)
    },
  }
}

/**
 * Shows sonner's "You submitted the following values:" toast of the form
 * examples, its description a <pre><code> of the values as JSON.
 */
export function submittedToast(toast, values, title = "You submitted the following values:") {
  const pre = document.createElement("pre")
  pre.className = "mt-2 w-[320px] overflow-x-auto rounded-md bg-code p-4 text-code-foreground"
  const code = document.createElement("code")
  code.textContent = JSON.stringify(values, null, 2)
  pre.appendChild(code)
  toast(title, {
    description: pre,
    position: "bottom-right",
    classNames: { content: "flex flex-col gap-2" },
    style: { "--border-radius": "calc(var(--radius)  + 4px)" },
  })
}

/** Toggles a checkbox/switch (Radix) to checked by clicking it, as a user would. */
export function setChecked(control, checked) {
  if (!control) return
  const cur = control.getAttribute("aria-checked") === "true"
  if (cur !== !!checked) control.click()
}
