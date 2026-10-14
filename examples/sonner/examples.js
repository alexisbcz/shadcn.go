// The sonner-demo and sonner-types examples: their buttons' onClick
// handlers, calling sonner's toast() exactly like the React examples.
import { define } from "./core.js"
import { toast } from "./sonner.js"

const handlers = {
  "Show Toast": () =>
    toast("Event has been created", {
      description: "Sunday, December 03, 2023 at 9:00 AM",
      action: {
        label: "Undo",
        onClick: () => console.log("Undo"),
      },
    }),
  Default: () => toast("Event has been created"),
  Success: () => toast.success("Event has been created"),
  Info: () => toast.info("Be at the area 10 minutes before the event time"),
  Warning: () => toast.warning("Event start time cannot be earlier than 8am"),
  Error: () => toast.error("Event has not been created"),
  Promise: () => {
    toast.promise(() => new Promise((resolve) => setTimeout(() => resolve({ name: "Event" }), 2000)), {
      loading: "Loading...",
      success: (data) => `${data.name} has been created`,
      error: "Error",
    })
  },
  // sonner-options (parity fixture)
  Rich: () => toast.success("Rich success", { richColors: true, closeButton: true, description: "With a close button" }),
  Confirm: () =>
    toast("Are you sure?", {
      duration: Infinity,
      cancel: { label: "Cancel", onClick: () => {} },
      action: { label: "Confirm", onClick: () => {} },
    }),
  Top: () => toast.info("At the top", { position: "top-center" }),
  "Dismiss all": () => toast.dismiss(),
  Reject: () =>
    toast.promise(() => new Promise((_, reject) => setTimeout(() => reject(new Error("nope")), 1000)), {
      loading: "Saving...",
      success: "Saved",
      error: "Could not save",
    }),
  Sticky: () => toast.warning("Cannot be swiped", { dismissible: false, invert: true }),
}

define("example-sonner", {
  // Buttons of the example previews, on pages that have a toaster.
  selector: '[data-slot="button"]',
  mount(button) {
    if (!document.querySelector("template[data-sonner-toaster]")) return
    const handler = handlers[button.textContent.trim()]
    if (!handler) return
    button.addEventListener("click", handler)
    return () => button.removeEventListener("click", handler)
  },
})
