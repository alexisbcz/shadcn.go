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

// The docs previews (ui/apps/v4/examples/radix/sonner-*.tsx) that reuse a
// label with another handler, by the marker's data-example-sonner value.
const variants = {
  description: {
    "Show Toast": () => toast("Event has been created", { description: "Monday, January 3rd at 6:00pm" }),
  },
}
for (const [label, position] of [
  ["Top Left", "top-left"],
  ["Top Center", "top-center"],
  ["Top Right", "top-right"],
  ["Bottom Left", "bottom-left"],
  ["Bottom Center", "bottom-center"],
  ["Bottom Right", "bottom-right"],
]) {
  handlers[label] = () => toast("Event has been created", { position })
}

define("example-sonner", {
  // The marker template rendered inside each example: the buttons next to
  // it get the example's onClick handlers.
  selector: "template[data-example-sonner]",
  mount(tpl) {
    const root = tpl.parentElement
    if (!root) return
    const bound = []
    for (const button of root.querySelectorAll('[data-slot="button"]')) {
      const label = button.textContent.trim()
      const handler = variants[tpl.dataset.exampleSonner]?.[label] ?? handlers[label]
      if (!handler) continue
      button.addEventListener("click", handler)
      bound.push([button, handler])
    }
    return () => bound.forEach(([button, handler]) => button.removeEventListener("click", handler))
  },
})
