// The command-dialog example: ⌘J / Ctrl+J toggles the palette.
import { getDialog } from "./dialog.js"

document.addEventListener("keydown", (e) => {
  if (e.key === "j" && (e.metaKey || e.ctrlKey)) {
    e.preventDefault()
    const tpl = [...document.querySelectorAll("template[data-radix-dialog]")].find((t) => t.content.querySelector("[cmdk-root]"))
    const d = tpl && getDialog(tpl)
    d?.toggle()
  }
})

// The docs previews' "Open Menu" buttons (onClick={() => setOpen(true)}):
// they open the CommandDialog rendered next to them.
document.addEventListener("click", (e) => {
  const button = e.target instanceof Element ? e.target.closest("[data-example-command-open]") : null
  if (!button) return
  const tpl = [...(button.parentElement?.querySelectorAll("template[data-radix-dialog]") ?? [])].find((t) => t.content.querySelector("[cmdk-root]"))
  const d = tpl && getDialog(tpl)
  d?.setOpen(true)
})
