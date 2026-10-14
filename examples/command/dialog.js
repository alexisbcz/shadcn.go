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
