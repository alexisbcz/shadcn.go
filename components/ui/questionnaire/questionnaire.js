// questionnaire.js: the behavior of Questionnaire, a port of the
// @shadcn/react/questionnaire primitive (use-questionnaire-root, -item,
// -choice and -input) onto the server-rendered markup.
//
// The form shows one item (fieldset) at a time. Next validates the item
// (an answer is required unless the item was skipped) and moves on; Skip
// leaves an optional item; Submit validates every item. Keyboard: Enter on
// an answer confirms it, Cmd/Ctrl+Enter confirms the item, ArrowUp/Down
// move between answers, ArrowLeft/Right between items, letter or number
// shortcuts pick choices.
//
//   import { getQuestionnaire } from "/_shadcn/questionnaire.js"
//   getQuestionnaire(form).setItem("timing")
//   form.addEventListener("shadcn:item-change", (e) => e.detail.item) // cancelable
//   item.addEventListener("shadcn:status-change", (e) => e.detail.status)
import { define, emit, attr } from "./core.js"

const instances = new WeakMap()

/** Returns the controller of a questionnaire form (or an element inside it). */
export function getQuestionnaire(el) {
  const form = el?.closest?.('[data-slot="questionnaire"]')
  return form ? instances.get(form) : undefined
}

const SHORTCUT_KEYS = {
  letters: Array.from({ length: 26 }, (_, i) => String.fromCharCode(65 + i)),
  numbers: Array.from({ length: 9 }, (_, i) => String(i + 1)),
}

const hasValue = (v) => String(v ?? "").trim().length > 0

function isTextEntryTarget(el) {
  if (el instanceof HTMLTextAreaElement || el instanceof HTMLSelectElement) return true
  if (el instanceof HTMLInputElement) return !["button", "checkbox", "radio", "reset", "submit"].includes(el.type)
  return el instanceof HTMLElement && el.isContentEditable
}

const isRadioTarget = (el) => el instanceof HTMLInputElement && el.type === "radio"

define("questionnaire", {
  selector: '[data-slot="questionnaire"]',
  mount(form) {
    const own = (el) => el.closest('[data-slot="questionnaire"]') === form
    const ownIn = (root, sel) => [...root.querySelectorAll(sel)].filter(own)

    // --- model ------------------------------------------------------------------------------
    let shortcuts = form.dataset.shortcuts || null
    const order = form.dataset.questionnaireOrder?.split(" ").filter(Boolean) ?? null
    const items = ownIn(form, '[data-slot="questionnaire-item"]').map((el) => {
      const answers = [...el.querySelectorAll('[data-slot="questionnaire-choice-input"], [data-slot="questionnaire-input"]')]
        .filter((a) => a.closest('[data-slot="questionnaire-item"]') === el)
        .map((input) => {
          const choice = input.dataset.slot === "questionnaire-choice-input"
          const label = choice ? input.closest('[data-slot="questionnaire-choice"]') : null
          return {
            input,
            label,
            type: choice ? "choice" : "input",
            shortcut: label?.dataset.shortcut ?? null,
            ownDisabled: input.disabled,
            defaultSelected: choice ? input.checked : input.hasAttribute("data-filled"),
            wrapper: choice ? null : input.closest('[data-slot="questionnaire-input-wrapper"]'),
          }
        })
      const item = {
        el,
        name: el.dataset.questionnaireName,
        required: el.hasAttribute("data-required"),
        multiple: el.hasAttribute("data-multiple"),
        get disabled() {
          return el.disabled
        },
        answers,
        selected: new Set(answers.filter((a) => a.defaultSelected)),
        skipped: false,
        validationAttempted: false,
        status: el.dataset.status,
        descIds: [...el.querySelectorAll('[data-slot="questionnaire-description"]')].map((d) => d.id).filter(Boolean),
        errors: [...el.querySelectorAll('[data-slot="questionnaire-error"]')].filter((e) => e.closest('[data-slot="questionnaire-item"]') === el),
        userDescribedBy: null,
      }
      const describedBy = (el.getAttribute("aria-describedby") ?? "").split(" ").filter(Boolean)
      item.userDescribedBy = describedBy.filter((id) => !item.descIds.includes(id)).join(" ") || null
      const keys = (el.getAttribute("aria-keyshortcuts") ?? "").replace(/Meta\+Enter Control\+Enter|ArrowUp ArrowDown|ArrowLeft|ArrowRight/g, "").trim()
      item.userKeys = keys || null
      return item
    })
    const byName = new Map(items.map((i) => [i.name, i]))
    const progress = ownIn(form, '[data-slot="questionnaire-progress"]').map((el) => ({
      el,
      defaultLabel: /^Question \d+ of \d+$/.test(el.textContent.trim()) || el.textContent.trim() === "",
    }))
    const navs = ownIn(form, "[data-questionnaire-nav]")
    const initialActive = items.find((i) => i.el.hasAttribute("data-active"))?.name ?? null
    let active = initialActive

    const enabledAnswers = (item) => item.answers.filter((a) => !a.input.disabled)
    const logical = () => {
      const names = order ?? items.map((i) => i.name)
      return names.map((n) => byName.get(n)).filter((i) => i && !i.disabled)
    }
    const statusOf = (item) =>
      item.skipped ? "skipped" : enabledAnswers(item).some((a) => item.selected.has(a)) ? "answered" : "unanswered"
    const valid = (item) =>
      item.disabled || (statusOf(item) === "skipped" && !item.required) || (!item.el.hasAttribute("data-questionnaire-invalid") && statusOf(item) === "answered")
    // data-questionnaire-invalid on an item is upstream's invalid prop.
    const externallyInvalid = (item) => item.el.hasAttribute("data-questionnaire-invalid")
    const invalid = (item) =>
      !item.disabled &&
      !(statusOf(item) === "skipped" && !item.required) &&
      (externallyInvalid(item) || (item.validationAttempted && !valid(item)))
    const filled = (a) => (a.type === "choice" ? a.input.checked : a.input.hasAttribute("name") && hasValue(a.input.value))

    // --- rendering ------------------------------------------------------------------------------
    function render() {
      const list = logical()
      const index = list.findIndex((i) => i.name === active)
      const total = list.length
      const current = index < 0 ? 0 : index + 1
      const first = total > 0 && index === 0
      const last = total > 0 && index === total - 1
      const activeItem = index < 0 ? null : list[index]

      for (const el of [form, ...progress.map((p) => p.el)]) {
        attr(el, "data-current", total ? current : null)
        attr(el, "data-first", first)
        attr(el, "data-last", last)
        attr(el, "data-total", total || null)
      }
      for (const p of progress) {
        const label = total ? `Question ${current} of ${total}` : null
        attr(p.el, "aria-valuemax", total || null)
        attr(p.el, "aria-valuemin", total ? 1 : null)
        attr(p.el, "aria-valuenow", total ? current : null)
        attr(p.el, "aria-valuetext", label)
        if (p.defaultLabel && p.el.textContent !== (label ?? "")) p.el.textContent = label ?? ""
        p.el.querySelectorAll("[data-questionnaire-current]").forEach((el) => (el.textContent = String(current)))
        p.el.querySelectorAll("[data-questionnaire-total]").forEach((el) => (el.textContent = String(total)))
      }

      for (const item of items) {
        const isActive = !item.disabled && item.name === active
        const status = statusOf(item)
        const isInvalid = invalid(item)
        if (status !== item.status) {
          item.status = status
          emit(item.el, "shadcn:status-change", { status })
        }
        const el = item.el
        attr(el, "hidden", !isActive)
        attr(el, "inert", !isActive)
        attr(el, "data-active", isActive)
        attr(el, "data-disabled", el.disabled)
        attr(el, "data-status", status)
        attr(el, "data-invalid", isInvalid)
        attr(el, "aria-invalid", isInvalid ? "true" : null)
        const describedBy = [...item.descIds, ...(isInvalid ? item.errors.map((e) => e.id).filter(Boolean) : []), item.userDescribedBy].filter(Boolean)
        attr(el, "aria-describedby", describedBy.join(" ") || null)
        const keys = [
          item.userKeys,
          isActive ? "Meta+Enter Control+Enter" : null,
          isActive && enabledAnswers(item).length ? "ArrowUp ArrowDown" : null,
          isActive && !first ? "ArrowLeft" : null,
          isActive && !last && status !== "unanswered" ? "ArrowRight" : null,
        ].filter(Boolean)
        attr(el, "aria-keyshortcuts", keys.join(" ") || null)

        for (const a of item.answers) {
          const disabled = a.input.disabled
          const selected = item.selected.has(a)
          if (a.type === "choice") {
            const checked = status === "skipped" ? false : selected
            a.input.checked = checked
            if (status === "skipped") a.input.removeAttribute("name")
            else a.input.setAttribute("name", item.name)
            for (const t of [a.label, a.input]) {
              attr(t, "data-checked", checked)
              attr(t, "data-unchecked", !checked)
              attr(t, "data-invalid", isInvalid)
            }
            attr(a.input, "aria-invalid", isInvalid ? "true" : null)
            attr(a.input, "aria-keyshortcuts", [a.shortcut, !disabled && checked ? "Enter" : null].filter(Boolean).join(" ") || null)
          } else {
            const isFilled = hasValue(a.input.value)
            if (selected) {
              a.input.setAttribute("name", item.name)
              a.input.removeAttribute("form")
            } else {
              a.input.removeAttribute("name")
              a.input.setAttribute("form", "")
            }
            attr(a.input, "data-filled", isFilled)
            attr(a.input, "data-empty", !isFilled)
            attr(a.input, "data-invalid", isInvalid)
            attr(a.input, "aria-invalid", isInvalid ? "true" : null)
            attr(a.input, "aria-keyshortcuts", !disabled && isFilled && selected ? "Enter" : null)
          }
        }
        for (const e of item.errors) {
          attr(e, "hidden", !isInvalid)
          attr(e, "role", isInvalid ? "alert" : null)
          attr(e, "data-invalid", isInvalid)
        }
      }

      const activeRequired = activeItem ? activeItem.required : null
      for (const b of navs) {
        const kind = b.dataset.questionnaireNav
        const visible =
          kind === "previous" ? total > 1 && !first : kind === "skip" ? activeRequired === false : kind === "next" ? total > 1 && !last : total > 0 && last
        const shortcut = (kind === "next" || kind === "submit") && visible && !b.disabled ? "Enter" : null
        attr(b, "aria-hidden", visible ? null : "true")
        attr(b, "hidden", !visible)
        attr(b, "inert", !visible)
        attr(b, "data-hidden", !visible)
        attr(b, "data-visible", visible)
        attr(b, "aria-keyshortcuts", shortcut)
        attr(b, "data-shortcut", shortcut)
        attr(b, "data-status", activeItem ? statusOf(activeItem) : null)
        if (!visible) b.setAttribute("tabindex", "-1")
        else if (b.dataset.questionnaireTabindex !== undefined) b.setAttribute("tabindex", b.dataset.questionnaireTabindex)
        else b.removeAttribute("tabindex")
      }
    }

    // --- actions --------------------------------------------------------------------------------
    function focusItem(item, target) {
      if (!item) return
      if (target === "invalid") {
        const selectedInput = item.el.querySelector("input[data-filled][name]:not(:disabled)")
        const firstControl = item.el.querySelector("input:not([type=hidden]):not(:disabled), textarea:not(:disabled)")
        ;(selectedInput ?? firstControl ?? item.el).focus()
      } else {
        item.el.focus()
      }
    }

    function setItem(name, target = "item") {
      if (name === active) return
      if (!emit(form, "shadcn:item-change", { item: name })) return
      active = name
      render()
      focusItem(byName.get(name), target)
    }

    const validate = (item) => {
      item.validationAttempted = true
      render()
      return valid(item)
    }

    function activeIndex() {
      return logical().findIndex((i) => i.name === active)
    }

    function goPrevious() {
      const list = logical()
      const i = activeIndex()
      if (i <= 0) return
      setItem(list[i - 1].name)
    }

    function goNext() {
      const list = logical()
      const i = activeIndex()
      const item = list[i]
      if (!item || i >= list.length - 1) return
      if (!validate(item)) {
        focusItem(item, "invalid")
        return
      }
      setItem(list[i + 1].name)
    }

    function confirmCurrent() {
      const list = logical()
      const i = activeIndex()
      const item = list[i]
      if (!item) return
      if (!validate(item)) {
        focusItem(item, "invalid")
        return
      }
      if (i === list.length - 1) {
        form.requestSubmit()
        return
      }
      setItem(list[i + 1].name)
    }

    function skipCurrent() {
      const list = logical()
      const i = activeIndex()
      const item = list[i]
      if (!item || item.required) return
      item.selected.clear()
      item.skipped = true
      render()
      if (i < list.length - 1) {
        setItem(list[i + 1].name)
        return
      }
      queueMicrotask(() => form.requestSubmit())
    }

    function select(item, answer, selected) {
      if (!selected) item.selected.delete(answer)
      else if (!item.multiple) {
        item.selected.clear()
        item.selected.add(answer)
      } else item.selected.add(answer)
    }

    // --- events ---------------------------------------------------------------------------------
    const itemOf = (el) => items.find((i) => i.el === el.closest('[data-slot="questionnaire-item"]'))

    const onChange = (event) => {
      const item = itemOf(event.target)
      const answer = item?.answers.find((a) => a.input === event.target)
      if (!answer || answer.type !== "choice") return
      item.skipped = false
      select(item, answer, event.target.checked)
      render()
    }

    const onInput = (event) => {
      const item = itemOf(event.target)
      const answer = item?.answers.find((a) => a.input === event.target)
      if (!answer || answer.type !== "input") return
      item.skipped = false
      select(item, answer, hasValue(event.target.value))
      render()
    }

    const onClick = (event) => {
      const b = event.target.closest?.("[data-questionnaire-nav]")
      if (!b || !own(b) || event.defaultPrevented) return
      const kind = b.dataset.questionnaireNav
      if (kind === "previous") goPrevious()
      else if (kind === "skip") skipCurrent()
      else if (kind === "next") goNext()
    }

    // Capture: validation runs before the application's submit handlers.
    const onSubmit = (event) => {
      const list = logical()
      const firstInvalid = list.find((item) => {
        item.validationAttempted = true
        return !valid(item)
      })
      render()
      if (!firstInvalid) return
      event.preventDefault()
      event.stopImmediatePropagation()
      if (firstInvalid.name === active) focusItem(firstInvalid, "invalid")
      else setItem(firstInvalid.name, "invalid")
    }

    const onReset = (event) => {
      if (event.defaultPrevented) return
      for (const item of items) {
        item.validationAttempted = false
        item.skipped = false
        const defaults = item.answers.filter((a) => a.defaultSelected)
        item.selected = new Set(item.multiple ? defaults : defaults.slice(0, 1))
      }
      queueMicrotask(() => {
        const name = (initialActive && byName.get(initialActive) && !byName.get(initialActive).disabled ? initialActive : logical()[0]?.name) ?? null
        render()
        if (name) setItem(name)
      })
    }

    const onKeyDown = (event) => {
      const list = logical()
      const i = activeIndex()
      const item = list[i]
      if (event.defaultPrevented || event.isComposing || event.keyCode === 229 || !item || !(event.target instanceof Element)) return
      if (event.key === "Enter" && (event.metaKey || event.ctrlKey) && !event.altKey && !event.shiftKey) {
        event.preventDefault()
        if (!event.repeat) confirmCurrent()
        return
      }
      if (event.metaKey || event.ctrlKey || event.altKey) return
      const answers = enabledAnswers(item)
      if (event.key === "ArrowUp" || event.key === "ArrowDown") {
        const target = event.target
        const index = answers.findIndex((a) => a.input === target)
        const currentAnswer = index < 0 ? null : answers[index]
        const emptyNavigable =
          currentAnswer?.type === "input" &&
          ["email", "password", "search", "tel", "text", "url"].includes(target.type) &&
          !hasValue(target.value)
        const blocked = !answers.length || (isTextEntryTarget(target) && !emptyNavigable) || (index < 0 && target !== item.el)
        if (!blocked) {
          const dir = event.key === "ArrowDown" ? 1 : -1
          const next =
            index < 0
              ? (answers.find(filled) ?? (dir > 0 ? answers[0] : answers[answers.length - 1]))
              : answers[(index + dir + answers.length) % answers.length]
          if (next && next.input !== target && !(index >= 0 && isRadioTarget(target) && isRadioTarget(next.input))) {
            next.input.focus()
            if (next.type === "choice" && isRadioTarget(next.input)) next.input.click()
            event.preventDefault()
            return
          }
        }
      }
      if ((event.key === "ArrowLeft" || event.key === "ArrowRight") && !isTextEntryTarget(event.target) && !isRadioTarget(event.target)) {
        event.preventDefault()
        if (event.repeat) return
        if (event.key === "ArrowLeft") goPrevious()
        else if (statusOf(item) !== "unanswered") goNext()
        return
      }
      if (event.key === "Enter") {
        const answer = answers.find((a) => a.input === event.target)
        if (!answer) return
        event.preventDefault()
        if (!event.repeat && filled(answer)) confirmCurrent()
        return
      }
      if (!shortcuts || isTextEntryTarget(event.target)) return
      const key = shortcuts === "letters" ? event.key.toUpperCase() : event.key
      if (!SHORTCUT_KEYS[shortcuts]?.includes(key)) return
      const answer = answers.find((a) => a.type === "choice" && a.shortcut === key)
      if (!answer) return
      event.preventDefault()
      if (event.repeat) return
      answer.input.focus()
      answer.input.click()
    }

    for (const b of navs) if (b.hasAttribute("tabindex") && !b.hasAttribute("hidden")) b.dataset.questionnaireTabindex = b.getAttribute("tabindex")
    form.addEventListener("change", onChange)
    form.addEventListener("input", onInput)
    form.addEventListener("click", onClick)
    form.addEventListener("submit", onSubmit, true)
    form.addEventListener("reset", onReset)
    form.addEventListener("keydown", onKeyDown)
    // Assigns the shortcut keys of mode ("letters", "numbers" or null) to the
    // enabled choices, in order (upstream's shortcuts prop changing).
    function setShortcuts(mode) {
      shortcuts = mode || null
      attr(form, "data-shortcuts", shortcuts)
      for (const item of items) {
        const keys = SHORTCUT_KEYS[shortcuts] ?? []
        let i = 0
        for (const a of item.answers) {
          if (a.type !== "choice") continue
          a.shortcut = !a.input.disabled && i < keys.length ? keys[i++] : null
          const span = a.label.querySelector('[data-slot="questionnaire-choice-shortcut"]')
          for (const t of [a.label, a.input, span]) if (t) attr(t, "data-shortcut", a.shortcut)
          if (span) {
            span.textContent = a.shortcut ?? ""
            attr(span, "hidden", a.shortcut === null)
          }
        }
        item.el.querySelectorAll('[data-slot="questionnaire-choices"]').forEach((c) => attr(c, "data-shortcuts", shortcuts))
      }
      render()
    }

    instances.set(form, {
      setItem: (name) => setItem(name),
      setShortcuts,
      status: (name) => (byName.get(name) ? statusOf(byName.get(name)) : null),
      get item() {
        return active
      },
      render,
    })
    return () => {
      form.removeEventListener("change", onChange)
      form.removeEventListener("input", onInput)
      form.removeEventListener("click", onClick)
      form.removeEventListener("submit", onSubmit, true)
      form.removeEventListener("reset", onReset)
      form.removeEventListener("keydown", onKeyDown)
      instances.delete(form)
    }
  },
})
