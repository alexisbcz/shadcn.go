// The dropdown-menu-dialog example: "New File..." and "Share..." open the
// two dialogs (onSelect sets the example's React state).
import { getDialog } from "./dialog.js"

document.addEventListener("shadcn:select", (event) => {
  const item = event.target
  if (!(item instanceof Element) || item.getAttribute("data-slot") !== "dropdown-menu-item") return
  const text = item.textContent.trim()
  const index = text === "New File..." ? 0 : text === "Share..." ? 1 : -1
  if (index < 0) return
  const tpl = document.querySelectorAll("template[data-radix-dialog]")[index]
  if (tpl) getDialog(tpl)?.setOpen(true)
})
