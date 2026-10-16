// The @shadcn/helpers docs demos (ai-sdk-helper-demo, tanstack-ai-helper-demo):
// useChat over createChat's scripted transport. Send plays the next user
// turn, then the assistant reply streams word by word (splitTextDeltas,
// delayMs 20) after the scripted sleep; Reset restores the empty state.
import { define } from "./core.js"

const sleep = (ms, signal) =>
  new Promise((resolve) => {
    const t = setTimeout(resolve, ms)
    signal.addEventListener("abort", () => {
      clearTimeout(t)
      resolve()
    })
  })

const deltas = (text) => text.match(/\S+\s*/g) ?? [text]

const paragraphs = (box, text, cls) =>
  box.replaceChildren(
    ...text
      .split(/\n\s*\n/)
      .map((p) => p.trim())
      .filter(Boolean)
      .map((p) => {
        const el = document.createElement("p")
        el.className = cls
        el.textContent = p
        return el
      })
  )

define("example-helpers-chat", {
  selector: "[data-example-helpers-chat]",
  mount(root) {
    const script = JSON.parse(root.dataset.exampleHelpersChat)
    root.removeAttribute("data-example-helpers-chat")
    const proto = root.querySelector("template[data-example-helpers-chat-proto]").content
    const cardContent = root.querySelector('[data-slot="card-content"]')
    const emptyState = cardContent.firstElementChild
    const form = root.querySelector("form")
    const area = form.querySelector("textarea")
    const send = form.querySelector('button[type="submit"]')
    const reset = root.querySelector('[aria-label="Reset conversation"]')
    let turn = 0
    let status = "ready"
    let content = null
    let ctrl = new AbortController()

    const nextUser = () => {
      for (let i = turn; i < script.length; i++) if (script[i].role === "user") return script[i]
      return null
    }
    const sync = () => {
      const next = nextUser()
      area.value = next ? next.parts.map((p) => p.text ?? "").join("") : ""
      area.setAttribute("data-status", status)
      const busy = status === "submitted" || status === "streaming"
      send.disabled = !next || busy
      if (reset) reset.disabled = busy
      content?.setAttribute("aria-busy", String(busy))
    }
    const append = (el) => content.insertBefore(el, content.querySelector(":scope > [data-message-scroller-spacer]"))
    const row = (role, id) => {
      const el = proto.querySelector(`[data-role="${role}"]`).cloneNode(true)
      el.removeAttribute("data-role")
      el.setAttribute("data-message-id", id)
      return el
    }

    const play = async (signal) => {
      const user = script[turn++]
      if (!content) {
        const scroller = proto.querySelector('[data-slot="message-scroller"]').cloneNode(true)
        cardContent.replaceChildren(scroller)
        content = scroller.querySelector('[data-slot="message-scroller-content"]')
      }
      const u = row("user", `user-${turn}`)
      paragraphs(u.querySelector('[data-slot="bubble-content"]'), user.parts[0].text, "whitespace-pre-wrap")
      append(u)
      if (!matchMedia("(prefers-reduced-motion: reduce)").matches) {
        u.animate(
          [
            { opacity: 0, transform: "translateY(10px)" },
            { opacity: 1, transform: "none" },
          ],
          { duration: 260, easing: "cubic-bezier(0.16, 1, 0.3, 1)" }
        )
      }
      status = "submitted"
      sync()
      const reply = script[turn]
      if (!reply || reply.role !== "assistant") {
        status = "ready"
        sync()
        return
      }
      turn++
      sync()
      await sleep(reply.sleep, signal)
      if (signal.aborted) return
      status = "streaming"
      sync()
      const a = row("assistant", `assistant-${turn}`)
      const msgContent = a.querySelector('[data-slot="message-content"]')
      const bubbleProto = msgContent.firstElementChild
      bubbleProto.remove()
      append(a)
      for (const part of reply.parts) {
        if (part.type === "sleep") {
          await sleep(part.sleep, signal)
          if (signal.aborted) return
          continue
        }
        let box, cls
        if (part.type === "reasoning") {
          const r = proto.querySelector('[data-role="reasoning"]').cloneNode(true)
          r.removeAttribute("data-role")
          msgContent.append(r)
          box = r.lastElementChild
          cls = "whitespace-pre-wrap"
        } else {
          const b = bubbleProto.cloneNode(true)
          msgContent.append(b)
          box = b.querySelector('[data-slot="bubble-content"]')
          cls = "whitespace-pre-wrap"
        }
        let text = ""
        for (const d of deltas(part.text)) {
          await sleep(20, signal)
          if (signal.aborted) return
          text += d
          paragraphs(box, text, cls)
        }
      }
      status = "ready"
      sync()
    }

    const onSubmit = (e) => {
      e.preventDefault()
      if (!nextUser() || status === "submitted" || status === "streaming") return
      play(ctrl.signal)
    }
    const onReset = () => {
      ctrl.abort()
      ctrl = new AbortController()
      turn = 0
      status = "ready"
      content = null
      cardContent.replaceChildren(emptyState)
      sync()
    }
    form.addEventListener("submit", onSubmit)
    reset?.addEventListener("click", onReset)
    return () => {
      ctrl.abort()
      form.removeEventListener("submit", onSubmit)
      reset?.removeEventListener("click", onReset)
    }
  },
})
