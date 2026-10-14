import { CheckIcon, PlusIcon } from "lucide-react"

import {
  Avatar,
  AvatarBadge,
  AvatarFallback,
  AvatarGroup,
  AvatarGroupCount,
  AvatarImage,
} from "@/registry/new-york-v4/ui/avatar"

// A 4x4 orange PNG: an image that loads without the network.
const pixel =
  "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAQAAAAECAIAAAAmkwkpAAAAEElEQVR4nGP4WSwGRwzEcQCBwhghZfPo6AAAAABJRU5ErkJggg=="

// Port fixture: what the upstream example does not use: the sm and lg
// sizes, an image that loads from a data: URL, an image that fails
// (fallback stays), an image without src, fallbacks with delayMs (shown
// after the delay unless the image loaded first), and AvatarBadge,
// AvatarGroup and AvatarGroupCount.
export default function AvatarStates() {
  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center gap-4">
        <Avatar size="sm">
          <AvatarImage src={pixel} alt="sm" />
          <AvatarFallback>SM</AvatarFallback>
        </Avatar>
        <Avatar>
          <AvatarImage src={pixel} alt="default" />
          <AvatarFallback>DF</AvatarFallback>
        </Avatar>
        <Avatar size="lg">
          <AvatarImage src={pixel} alt="lg" />
          <AvatarFallback>LG</AvatarFallback>
        </Avatar>
      </div>
      <div className="flex items-center gap-4">
        <Avatar size="sm">
          <AvatarImage src="data:image/png;base64,AAAA" alt="broken" />
          <AvatarFallback>BR</AvatarFallback>
        </Avatar>
        <Avatar>
          <AvatarImage alt="no src" />
          <AvatarFallback>NS</AvatarFallback>
        </Avatar>
        <Avatar size="lg">
          <AvatarImage src="data:image/png;base64,AAAA" alt="delayed" />
          <AvatarFallback delayMs={600}>DL</AvatarFallback>
        </Avatar>
        <Avatar>
          <AvatarImage src={pixel} alt="fast" />
          <AvatarFallback delayMs={600}>FA</AvatarFallback>
        </Avatar>
        <Avatar>
          <AvatarFallback>FB</AvatarFallback>
        </Avatar>
      </div>
      <div className="flex items-center gap-4">
        <Avatar size="sm">
          <AvatarFallback>S</AvatarFallback>
          <AvatarBadge />
        </Avatar>
        <Avatar>
          <AvatarFallback>D</AvatarFallback>
          <AvatarBadge>
            <CheckIcon />
          </AvatarBadge>
        </Avatar>
        <Avatar size="lg">
          <AvatarFallback>L</AvatarFallback>
          <AvatarBadge className="bg-green-600">
            <CheckIcon />
          </AvatarBadge>
        </Avatar>
      </div>
      <AvatarGroup>
        <Avatar>
          <AvatarImage src={pixel} alt="one" />
          <AvatarFallback>A</AvatarFallback>
        </Avatar>
        <Avatar>
          <AvatarFallback>B</AvatarFallback>
        </Avatar>
        <AvatarGroupCount>+3</AvatarGroupCount>
      </AvatarGroup>
      <AvatarGroup>
        <Avatar size="lg">
          <AvatarFallback>A</AvatarFallback>
        </Avatar>
        <Avatar size="lg">
          <AvatarFallback>B</AvatarFallback>
        </Avatar>
        <AvatarGroupCount>
          <PlusIcon />
        </AvatarGroupCount>
      </AvatarGroup>
    </div>
  )
}
