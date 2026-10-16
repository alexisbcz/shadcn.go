// The questionnaire docs examples' React state: submit toasts, the
// controlled checkpoint, the conditional cloud item, the navigation state
// (Next disabled while unanswered), the custom progress bars, the shortcut
// selector and the schema validation.
import { define } from "./core.js"
import { toast } from "./sonner.js"
import { getQuestionnaire } from "./questionnaire.js"

const get = (data, name) => data.get(name) ?? "None"
const all = (data, name) => data.getAll(name).join(", ") || "None"

// Submit toasts: [title, description(FormData, form)].
const toasts = {
  demo: ["Agent plan saved", (d) => `Direction: ${get(d, "direction")} · Progress signals: ${all(d, "signals")} · Timing: ${get(d, "timing")}`],
  animated: ["Agent workflow saved", (d) => `Task: ${get(d, "task")} · Review: ${get(d, "review")} · Delivery: ${get(d, "delivery")}`],
  card: ["Agent task created", (d) => `Task: ${get(d, "task")} · Handoff: ${get(d, "output")}`],
  conditional: [
    "Execution plan saved",
    (d) => `Runtime: ${get(d, "runtime")} · Environment: ${d.get("environment") ?? "Not applicable"} · Approval: ${get(d, "approval")}`,
  ],
  controlled: ["Agent workflow configured", (d) => `Scope: ${get(d, "scope")} · Verification: ${get(d, "checks")} · Output: ${get(d, "output")}`],
  dialog: ["Clarification sent", (d) => `Scope: ${get(d, "scope")} · Verification: ${get(d, "tests")}`],
  freeform: ["Approach selected", (d) => `Approach: ${get(d, "approach")}`],
  multiple: ["Context selected", (d) => `Context: ${all(d, "context")}`],
  "navigation-state": ["Permissions saved", (d) => `Permission: ${get(d, "permission")} · Verification: ${get(d, "verification")}`],
  progress: [
    "Pull request plan ready",
    (d) => `Scope: ${get(d, "scope")} · Commits: ${get(d, "strategy")} · Tests: ${get(d, "tests")} · Delivery: ${get(d, "delivery")}`,
  ],
  resume: [
    "Draft updated",
    (d) => `Migration: ${get(d, "change")} · Verification: ${all(d, "verification")} · Notes: ${d.get("notes") || "None"}`,
  ],
  shortcuts: ["Next action selected", (d, form) => `Action: ${get(d, "action")} · Shortcuts: ${form.dataset.shortcuts || "none"}`],
  skip: [
    "Agent brief submitted",
    (d, form) => {
      const skipped = getQuestionnaire(form)?.status("constraints") === "skipped"
      return `Task: ${get(d, "task")} · Constraints: ${skipped ? "Skipped" : get(d, "constraints")} · Review: ${get(d, "review")}`
    },
  ],
}

const setups = {
  controlled(form) {
    const labels = { scope: "Change scope", checks: "Verification", output: "Final output" }
    const label = form.parentElement?.querySelector("[data-example-questionnaire-checkpoint]")
    const onChange = (e) => {
      if (e.target === form && label) label.textContent = labels[e.detail.item] ?? e.detail.item
    }
    form.addEventListener("shadcn:item-change", onChange)
    return () => form.removeEventListener("shadcn:item-change", onChange)
  },

  conditional(form) {
    const environment = form.querySelector('[data-questionnaire-name="environment"]')
    const onChange = (e) => {
      if (e.target.name !== "runtime") return
      environment.disabled = e.target.value !== "cloud"
      getQuestionnaire(form)?.render()
    }
    form.addEventListener("change", onChange)
    return () => form.removeEventListener("change", onChange)
  },

  "navigation-state"(form) {
    const buttons = [...form.querySelectorAll('[data-questionnaire-nav="next"], [data-questionnaire-nav="submit"]')]
    const sync = () => {
      const q = getQuestionnaire(form)
      if (!q) return
      const unanswered = q.status(q.item) === "unanswered"
      for (const b of buttons) b.disabled = unanswered
      q.render()
    }
    form.addEventListener("shadcn:status-change", sync)
    form.addEventListener("shadcn:item-change", () => queueMicrotask(sync))
    return () => form.removeEventListener("shadcn:status-change", sync)
  },

  progress(form) {
    const bars = [...form.querySelectorAll("[data-example-questionnaire-bars] > span")]
    const onChange = () =>
      queueMicrotask(() => {
        const current = Number(form.dataset.current || 0)
        bars.forEach((bar, i) => (bar.className = `h-1.5 flex-1 rounded-full ${i < current ? "bg-primary" : "bg-muted"}`))
      })
    form.addEventListener("shadcn:item-change", onChange)
    return () => form.removeEventListener("shadcn:item-change", onChange)
  },

  resume(form) {
    const onReset = () => toast("Saved answers restored")
    form.addEventListener("reset", onReset)
    return () => form.removeEventListener("reset", onReset)
  },

  shortcuts(form) {
    const select = form.parentElement?.querySelector("[data-example-questionnaire-shortcuts] select, select[data-example-questionnaire-shortcuts]")
    if (!select) return
    const onChange = () => getQuestionnaire(form)?.setShortcuts(select.value === "none" ? null : select.value)
    select.addEventListener("change", onChange)
    return () => select.removeEventListener("change", onChange)
  },

  validation(form) {
    const item = (name) => form.querySelector(`[data-questionnaire-name="${name}"]`)
    const errorOf = (name) => item(name).querySelector('[data-slot="questionnaire-error"]')
    const defaults = { detail: errorOf("detail").textContent, audience: errorOf("audience").textContent }
    const clear = (name) => {
      if (!item(name).hasAttribute("data-questionnaire-invalid")) return
      item(name).removeAttribute("data-questionnaire-invalid")
      errorOf(name).textContent = defaults[name]
      getQuestionnaire(form)?.render()
    }
    const onChange = (e) => {
      if (e.target.name === "detail" || e.target.name === "audience") clear(e.target.name)
    }
    const onSubmit = (e) => {
      const d = new FormData(form)
      if (d.get("audience") === "public" && d.get("detail") === "summary") {
        e.preventDefault()
        e.stopImmediatePropagation()
        item("detail").setAttribute("data-questionnaire-invalid", "")
        errorOf("detail").textContent = "Public answers need enough context. Choose a complete answer."
        const q = getQuestionnaire(form)
        q?.render()
        q?.setItem("detail")
      }
    }
    form.addEventListener("change", onChange)
    form.addEventListener("submit", onSubmit)
    return () => {
      form.removeEventListener("change", onChange)
      form.removeEventListener("submit", onSubmit)
    }
  },
}

define("example-questionnaire", {
  selector: "[data-example-questionnaire]",
  mount(form) {
    const name = form.dataset.exampleQuestionnaire
    const cleanup = setups[name]?.(form)
    const onSubmit = (e) => {
      if (e.defaultPrevented) return
      e.preventDefault()
      const [title, description] = toasts[name] ?? []
      if (!title) return
      if (name === "dialog") {
        const content = form.closest('[data-slot="dialog-content"]')
        content?.querySelector('[data-slot="dialog-close"]')?.click()
      }
      toast(title, { description: description(new FormData(form), form) })
    }
    form.addEventListener("submit", onSubmit)
    return () => {
      cleanup?.()
      form.removeEventListener("submit", onSubmit)
    }
  },
})
