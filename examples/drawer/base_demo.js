// drawer-demo (base/aria): Confirm Delivery Time toasts the picked slot
// (sonner's toast("Delivery time confirmed", { description })); the
// button is also a basedrawer Close, which closes the drawer.
import { toast } from "./sonner.js"

document.addEventListener("click", (e) => {
  const btn = e.target instanceof Element ? e.target.closest("[data-drawer-confirm-delivery]") : null
  if (!btn) return
  const popup = btn.closest('[data-slot="drawer-popup"]')
  const checked = popup?.querySelector('[data-delivery-label][aria-checked="true"], [data-delivery-label][data-state="checked"]')
  if (!checked) return
  toast("Delivery time confirmed", { description: checked.dataset.deliveryLabel })
})
