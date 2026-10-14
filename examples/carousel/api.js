// The carousel-api example: "Slide {current} of {count}" from the Embla API,
// as four text nodes like React renders them.
import { define } from "./core.js"
import { onCarouselApi } from "./carousel.js"

define("example-carousel-api", {
  selector: "[data-example-carousel-api]",
  mount(root) {
    root.removeAttribute("data-example-carousel-api")
    const carousel = root.querySelector('[data-slot="carousel"]')
    const status = root.querySelector(":scope > .text-center")
    if (!carousel || !status) return
    let count = 0
    let current = 0
    const render = () => status.replaceChildren("Slide ", String(current), " of ", String(count))
    render()
    onCarouselApi(carousel, (api) => {
      count = api.scrollSnapList().length
      current = api.selectedScrollSnap() + 1
      render()
      api.on("select", () => {
        current = api.selectedScrollSnap() + 1
        render()
      })
    })
  },
})
