// The input-otp-controlled example: shows the value typed so far.
import { define } from "./core.js"

define("example-input-otp-controlled", {
  selector: ".space-y-2:has(> [data-input-otp-container])",
  mount(root) {
    const text = root.querySelector(":scope > .text-center")
    if (!text) return
    const onChange = (e) => {
      const value = e.detail.value
      if (value === "") text.replaceChildren("Enter your one-time password.")
      else text.replaceChildren("You entered: ", value)
    }
    root.addEventListener("shadcn:value-change", onChange)
    return () => root.removeEventListener("shadcn:value-change", onChange)
  },
})
