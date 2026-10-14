// The sidebar-group-action example: onClick={() => toast("You clicked the
// group action!")}.
import { define } from "./core.js"
import { toast } from "./sonner.js"

define("example-sidebar-group-action", {
  selector: "[data-example-sidebar-group-action]",
  mount(el) {
    el.removeAttribute("data-example-sidebar-group-action")
    const onClick = () => toast("You clicked the group action!")
    el.addEventListener("click", onClick)
    return () => el.removeEventListener("click", onClick)
  },
})
