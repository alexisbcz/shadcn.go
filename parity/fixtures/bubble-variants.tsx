import {
  Bubble,
  BubbleContent,
  BubbleGroup,
  BubbleReactions,
} from "@/registry/new-york-v4/ui/bubble"

const variants = [
  "default",
  "secondary",
  "muted",
  "tinted",
  "outline",
  "ghost",
  "destructive",
] as const

// Port fixture: every Bubble variant on both sides, the default variant and
// alignment, a BubbleGroup, BubbleReactions on every side and alignment
// (defaults included), and a className override on each part.
export default function BubbleVariants() {
  return (
    <div className="grid w-full max-w-3xl grid-cols-2 gap-x-8 gap-y-8 py-8">
      <div className="flex flex-col gap-8">
        {variants.map((variant) => (
          <Bubble key={variant} variant={variant}>
            <BubbleContent>The {variant} bubble, at the start.</BubbleContent>
          </Bubble>
        ))}
        <Bubble>
          <BubbleContent>Default variant and alignment.</BubbleContent>
        </Bubble>
      </div>
      <div className="flex flex-col gap-8">
        {variants.map((variant) => (
          <Bubble key={variant} variant={variant} align="end">
            <BubbleContent>
              The {variant} bubble, at the end, with a longer text that wraps
              onto a second line.
            </BubbleContent>
          </Bubble>
        ))}
      </div>
      <BubbleGroup className="col-span-2 gap-4">
        <Bubble variant="muted">
          <BubbleContent>Reactions default to the bottom end.</BubbleContent>
          <BubbleReactions role="img" aria-label="Reaction: thumbs up">
            <span>👍</span>
          </BubbleReactions>
        </Bubble>
        <Bubble variant="secondary" align="end">
          <BubbleContent>Bottom start reactions.</BubbleContent>
          <BubbleReactions align="start">
            <span>👀</span>
            <span>🚀</span>
            <span>+2</span>
          </BubbleReactions>
        </Bubble>
        <Bubble variant="tinted" className="mt-4">
          <BubbleContent>Top start reactions.</BubbleContent>
          <BubbleReactions side="top" align="start">
            <span>🎉</span>
          </BubbleReactions>
        </Bubble>
        <Bubble variant="outline" align="end" className="mt-4">
          <BubbleContent className="font-medium">Top end reactions.</BubbleContent>
          <BubbleReactions side="top" align="end" className="bg-accent">
            <span>🔥</span>
            <span>👏</span>
          </BubbleReactions>
        </Bubble>
      </BubbleGroup>
    </div>
  )
}
