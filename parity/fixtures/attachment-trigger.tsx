import { CopyIcon, FileSearchIcon, FileTextIcon, XIcon } from "lucide-react"

import {
  Attachment,
  AttachmentAction,
  AttachmentActions,
  AttachmentContent,
  AttachmentDescription,
  AttachmentMedia,
  AttachmentTitle,
  AttachmentTrigger,
} from "@/registry/new-york-v4/ui/attachment"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/registry/new-york-v4/ui/dialog"

// Port fixture: AttachmentTrigger as a button (default type, and an
// explicit type), as a link (asChild), and as a dialog trigger, with
// actions that stay clickable above it.
export default function AttachmentTriggerFixture() {
  return (
    <div className="flex w-full max-w-sm flex-col gap-3 py-12">
      <Dialog>
        <Attachment className="w-full">
          <AttachmentMedia>
            <FileSearchIcon />
          </AttachmentMedia>
          <AttachmentContent>
            <AttachmentTitle>research-summary.pdf</AttachmentTitle>
            <AttachmentDescription>Open preview dialog</AttachmentDescription>
          </AttachmentContent>
          <AttachmentActions>
            <AttachmentAction aria-label="Copy link">
              <CopyIcon />
            </AttachmentAction>
            <AttachmentAction aria-label="Remove research-summary.pdf">
              <XIcon />
            </AttachmentAction>
          </AttachmentActions>
          <DialogTrigger asChild>
            <AttachmentTrigger aria-label="Preview research-summary.pdf" />
          </DialogTrigger>
        </Attachment>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>research-summary.pdf</DialogTitle>
            <DialogDescription>
              The attachment trigger fills the card and opens the dialog, while
              the actions stay independently clickable above it.
            </DialogDescription>
          </DialogHeader>
        </DialogContent>
      </Dialog>
      <Attachment size="sm" className="w-full">
        <AttachmentMedia>
          <FileTextIcon />
        </AttachmentMedia>
        <AttachmentContent>
          <AttachmentTitle>notes.txt</AttachmentTitle>
          <AttachmentDescription>Open as link</AttachmentDescription>
        </AttachmentContent>
        <AttachmentTrigger asChild>
          <a href="#notes" aria-label="Open notes.txt" />
        </AttachmentTrigger>
      </Attachment>
      <Attachment size="xs" state="idle" className="w-full">
        <AttachmentContent>
          <AttachmentTitle>draft.md</AttachmentTitle>
        </AttachmentContent>
        <AttachmentTrigger type="submit" className="rounded-lg" aria-label="Submit draft.md" />
      </Attachment>
    </div>
  )
}
