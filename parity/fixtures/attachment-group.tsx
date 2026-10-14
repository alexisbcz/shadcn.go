import { XIcon } from "lucide-react"

import {
  Attachment,
  AttachmentAction,
  AttachmentActions,
  AttachmentContent,
  AttachmentDescription,
  AttachmentGroup,
  AttachmentMedia,
  AttachmentTitle,
} from "@/registry/new-york-v4/ui/attachment"

// A 4x4 orange PNG: an image that loads without the network.
const pixel =
  "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAQAAAAECAIAAAAmkwkpAAAAEElEQVR4nGP4WSwGRwzEcQCBwhghZfPo6AAAAABJRU5ErkJggg=="

const states = ["idle", "uploading", "done"] as const

// Port fixture: AttachmentGroup (a scrolling, snapping, fading row) of
// vertical image attachments in each state (image media is dimmed unless
// idle or done).
export default function AttachmentGroupFixture() {
  return (
    <div className="flex w-full max-w-3xl flex-col gap-4 py-4">
      <AttachmentGroup className="w-full max-w-md">
        {states.map((state) => (
          <Attachment key={state} state={state} orientation="vertical">
            <AttachmentMedia variant="image">
              <img src={pixel} alt={state} />
            </AttachmentMedia>
            <AttachmentContent>
              <AttachmentTitle>{state}.png</AttachmentTitle>
              <AttachmentDescription>PNG · 820 KB</AttachmentDescription>
            </AttachmentContent>
            <AttachmentActions>
              <AttachmentAction aria-label={`Remove ${state}.png`}>
                <XIcon />
              </AttachmentAction>
            </AttachmentActions>
          </Attachment>
        ))}
        <Attachment state="error" orientation="vertical">
          <AttachmentMedia variant="image">
            <img src={pixel} alt="error" />
          </AttachmentMedia>
        </Attachment>
        <Attachment state="processing" orientation="vertical">
          <AttachmentMedia variant="image">
            <img src={pixel} alt="processing" />
          </AttachmentMedia>
        </Attachment>
      </AttachmentGroup>
    </div>
  )
}
