import { FileTextIcon, XIcon } from "lucide-react"

import {
  Attachment,
  AttachmentAction,
  AttachmentActions,
  AttachmentContent,
  AttachmentDescription,
  AttachmentMedia,
  AttachmentTitle,
} from "@/registry/new-york-v4/ui/attachment"
import { Spinner } from "@/registry/new-york-v4/ui/spinner"

// A 4x4 orange PNG: an image that loads without the network.
const pixel =
  "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAQAAAAECAIAAAAmkwkpAAAAEElEQVR4nGP4WSwGRwzEcQCBwhghZfPo6AAAAABJRU5ErkJggg=="

const sizes = ["default", "sm", "xs"] as const

// Port fixture: every size in both orientations, icon and image media,
// and an attachment without media and one with media only.
export default function AttachmentSizes() {
  return (
    <div className="flex w-full max-w-3xl flex-col gap-4 py-4">
      {sizes.map((size) => (
        <div key={size} className="flex flex-wrap items-start gap-3">
          <Attachment size={size}>
            <AttachmentMedia>
              <FileTextIcon />
            </AttachmentMedia>
            <AttachmentContent>
              <AttachmentTitle>{size} attachment</AttachmentTitle>
              <AttachmentDescription>PDF · 2.4 MB</AttachmentDescription>
            </AttachmentContent>
            <AttachmentActions>
              <AttachmentAction aria-label="Remove">
                <XIcon />
              </AttachmentAction>
            </AttachmentActions>
          </Attachment>
          <Attachment size={size}>
            <AttachmentContent>
              <AttachmentTitle>No media</AttachmentTitle>
            </AttachmentContent>
          </Attachment>
          <Attachment size={size}>
            <AttachmentMedia variant="image">
              <img src={pixel} alt="Pixel" />
            </AttachmentMedia>
          </Attachment>
          <Attachment size={size} orientation="vertical">
            <AttachmentMedia>
              <FileTextIcon />
            </AttachmentMedia>
            <AttachmentContent>
              <AttachmentTitle>{size}.pdf</AttachmentTitle>
              <AttachmentDescription>2.4 MB</AttachmentDescription>
            </AttachmentContent>
          </Attachment>
          <Attachment size={size} orientation="vertical" state="uploading">
            <AttachmentMedia>
              <Spinner />
            </AttachmentMedia>
          </Attachment>
        </div>
      ))}
    </div>
  )
}
