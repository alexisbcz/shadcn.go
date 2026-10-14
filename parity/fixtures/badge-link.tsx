import { ArrowUpRightIcon } from "lucide-react"

import { Badge } from "@/registry/new-york-v4/ui/badge"

// Port fixture: badges rendered as links with asChild, to exercise the
// [a&]:hover and focus-visible styles of every variant, plus aria-invalid.
export default function BadgeLink() {
  return (
    <div className="flex flex-wrap items-center gap-2">
      <Badge asChild>
        <a href="#">Default</a>
      </Badge>
      <Badge variant="secondary" asChild>
        <a href="#">Secondary</a>
      </Badge>
      <Badge variant="destructive" asChild>
        <a href="#">Destructive</a>
      </Badge>
      <Badge variant="outline" asChild>
        <a href="#">
          Outline <ArrowUpRightIcon />
        </a>
      </Badge>
      <Badge variant="ghost" asChild>
        <a href="#">Ghost</a>
      </Badge>
      <Badge variant="link" asChild>
        <a href="#">Link</a>
      </Badge>
      <Badge aria-invalid="true">Invalid</Badge>
    </div>
  )
}
