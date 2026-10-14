// The RTL example's state (app/(app)/examples/rtl/components): the
// language selector switches between the Arabic and Hebrew grids, and the
// examples' React state (price range, voice mode, favorite, GPU count).
import { define } from "./core.js"

define("example-rtl", {
  selector: "[data-rtl-example]",
  mount(root) {
    const onValue = (e) => {
      const target = e.target instanceof Element ? e.target : null
      if (!target) return
      // Language selector: each grid has one, showing that grid's language.
      // The select emits from its <template>, keyed to the trigger by id.
      const contentId = target.getAttribute("data-content-id")
      const trigger = contentId && root.querySelector(`[aria-controls="${CSS.escape(contentId)}"]`)
      if (trigger?.matches('[data-name="language-selector"]')) {
        const lang = e.detail?.value
        const current = target.closest('[data-slot="rtl-components"]')
        const next = root.querySelector(`:scope > [data-slot="rtl-components"][data-lang="${lang}"]`)
        if (!next || next === current) return
        current.hidden = true
        next.hidden = false
        return
      }
      // Price range: the description shows the range in the locale.
      const field = target.closest("[data-rtl-slider]")
      if (field && target.matches('[data-slot="slider"]') && Array.isArray(e.detail?.value)) {
        const fmt = new Intl.NumberFormat(field.dataset.rtlSlider)
        field.querySelectorAll(':scope > [data-slot="field-description"] > span.tabular-nums').forEach((span, i) => {
          const text = fmt.format(e.detail.value[i])
          if (span.textContent !== text) span.textContent = text
        })
      }
    }
    const onClick = (e) => {
      const target = e.target instanceof Element ? e.target : null
      if (!target) return
      // examples/button-group/voice.js and examples/input-group/button.js
      // match these buttons too: toggle only if they did not (onCapture
      // recorded the state before their handlers ran).
      const voice = target.closest("[data-rtl-voice]")
      if (voice) {
        setTimeout(() => {
          let on = voice.getAttribute("aria-pressed") === "true"
          if (String(on) === before.get(voice)) {
            on = !on
            voice.setAttribute("aria-pressed", String(on))
            voice.setAttribute("data-active", String(on))
          }
          const input = voice.closest('[data-slot="input-group"]')?.querySelector("[data-rtl-voice-on]")
          if (input) {
            input.placeholder = on ? input.dataset.rtlVoiceOn : input.dataset.rtlVoiceOff
            input.disabled = on
          }
        })
        return
      }
      const fav = target.closest("[data-rtl-favorite]")
      if (fav) {
        const svg = fav.querySelector("svg")
        setTimeout(() => {
          if (svg && (svg.dataset.favorite ?? "false") === before.get(fav)) svg.dataset.favorite = String(svg.dataset.favorite !== "true")
        })
        return
      }
      const step = target.closest("[data-rtl-gpu-step]")
      if (step) {
        const group = step.closest("[data-rtl-gpu]")
        const input = group?.querySelector("input")
        if (!input) return
        const n = Math.max(1, Math.min(99, (parseInt(input.value, 10) || 8) + Number(step.dataset.rtlGpuStep)))
        setGpu(group, n)
      }
    }
    const before = new WeakMap()
    const onCapture = (e) => {
      const target = e.target instanceof Element ? e.target : null
      const voice = target?.closest("[data-rtl-voice]")
      if (voice) before.set(voice, voice.getAttribute("aria-pressed"))
      const fav = target?.closest("[data-rtl-favorite]")
      if (fav) before.set(fav, fav.querySelector("svg")?.dataset.favorite ?? "false")
    }
    const setGpu = (group, n) => {
      const input = group.querySelector("input")
      input.value = String(n)
      group.querySelector('[data-rtl-gpu-step="-1"]').disabled = n <= 1
      group.querySelector('[data-rtl-gpu-step="1"]').disabled = n >= 99
    }
    const onInput = (e) => {
      const input = e.target instanceof HTMLInputElement ? e.target : null
      const group = input?.closest("[data-rtl-gpu]")
      if (!group) return
      const n = parseInt(input.value, 10)
      if (!isNaN(n) && n >= 1 && n <= 99) setGpu(group, n)
    }
    root.addEventListener("shadcn:value-change", onValue)
    root.addEventListener("click", onCapture, true)
    root.addEventListener("click", onClick)
    root.addEventListener("input", onInput)
    return () => {
      root.removeEventListener("shadcn:value-change", onValue)
      root.removeEventListener("click", onCapture, true)
      root.removeEventListener("click", onClick)
      root.removeEventListener("input", onInput)
    }
  },
})
