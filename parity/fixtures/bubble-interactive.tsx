import { Bubble, BubbleContent, BubbleReactions } from "@/registry/new-york-v4/ui/bubble"
import { Button } from "@/registry/new-york-v4/ui/button"

const variants = [
  "default",
  "secondary",
  "muted",
  "tinted",
  "outline",
  "ghost",
  "destructive",
] as const

// Port fixture: BubbleContent asChild as a button and as a link in every
// variant (hover and focus-visible styles), and a button in reactions.
export default function BubbleInteractive() {
  return (
    <div className="flex w-full max-w-xl flex-col gap-4 py-8">
      {variants.map((variant) => (
        <div key={variant} className="flex justify-between gap-4">
          <Bubble variant={variant}>
            <BubbleContent asChild>
              <button>{variant} button</button>
            </BubbleContent>
          </Bubble>
          <Bubble variant={variant} align="end">
            <BubbleContent asChild>
              <a href="#bubble">{variant} link</a>
            </BubbleContent>
          </Bubble>
        </div>
      ))}
      <Bubble variant="destructive" className="mb-6">
        <BubbleContent>Are you sure I can run this command?</BubbleContent>
        <BubbleReactions>
          <Button variant="ghost" size="xs">
            Yes, run it
          </Button>
        </BubbleReactions>
      </Bubble>
    </div>
  )
}
