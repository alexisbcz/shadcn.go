import {
  CheckIcon,
  ClockIcon,
  FileTextIcon,
  FileWarningIcon,
  RefreshCwIcon,
  XIcon,
} from "lucide-react"

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

// Port fixture: every Attachment state, horizontal (with actions) and
// vertical, and the defaults (state done, size default, horizontal).
export default function AttachmentStates() {
  return (
    <div className="flex w-full max-w-2xl flex-col gap-6 py-6">
      <div className="flex flex-col gap-2">
        <Attachment state="idle" className="w-full">
          <AttachmentMedia>
            <ClockIcon />
          </AttachmentMedia>
          <AttachmentContent>
            <AttachmentTitle>selected-file.pdf</AttachmentTitle>
            <AttachmentDescription>Ready to upload</AttachmentDescription>
          </AttachmentContent>
          <AttachmentActions>
            <AttachmentAction aria-label="Remove selected-file.pdf">
              <XIcon />
            </AttachmentAction>
          </AttachmentActions>
        </Attachment>
        <Attachment state="uploading" className="w-full">
          <AttachmentMedia>
            <Spinner />
          </AttachmentMedia>
          <AttachmentContent>
            <AttachmentTitle>design-system.zip</AttachmentTitle>
            <AttachmentDescription>Uploading · 64%</AttachmentDescription>
          </AttachmentContent>
          <AttachmentActions>
            <AttachmentAction aria-label="Cancel upload">
              <XIcon />
            </AttachmentAction>
          </AttachmentActions>
        </Attachment>
        <Attachment state="processing" className="w-full">
          <AttachmentMedia>
            <FileTextIcon />
          </AttachmentMedia>
          <AttachmentContent>
            <AttachmentTitle>market-research.pdf</AttachmentTitle>
            <AttachmentDescription>Processing document</AttachmentDescription>
          </AttachmentContent>
        </Attachment>
        <Attachment state="error" className="w-full">
          <AttachmentMedia>
            <FileWarningIcon />
          </AttachmentMedia>
          <AttachmentContent>
            <AttachmentTitle>financial-model.xlsx</AttachmentTitle>
            <AttachmentDescription>Upload failed. Try again.</AttachmentDescription>
          </AttachmentContent>
          <AttachmentActions>
            <AttachmentAction aria-label="Retry upload">
              <RefreshCwIcon />
            </AttachmentAction>
            <AttachmentAction variant="destructive" aria-label="Remove financial-model.xlsx">
              <XIcon />
            </AttachmentAction>
          </AttachmentActions>
        </Attachment>
        <Attachment>
          <AttachmentMedia>
            <CheckIcon />
          </AttachmentMedia>
          <AttachmentContent>
            <AttachmentTitle>uploaded-report.pdf</AttachmentTitle>
            <AttachmentDescription>Uploaded · 1.8 MB</AttachmentDescription>
          </AttachmentContent>
          <AttachmentActions>
            <AttachmentAction size="icon-sm" aria-label="Remove uploaded-report.pdf">
              <XIcon />
            </AttachmentAction>
          </AttachmentActions>
        </Attachment>
      </div>
      <div className="flex flex-wrap gap-3">
        <Attachment state="idle" orientation="vertical">
          <AttachmentMedia>
            <ClockIcon />
          </AttachmentMedia>
          <AttachmentContent>
            <AttachmentTitle>selected.pdf</AttachmentTitle>
            <AttachmentDescription>Ready</AttachmentDescription>
          </AttachmentContent>
        </Attachment>
        <Attachment state="uploading" orientation="vertical">
          <AttachmentMedia>
            <Spinner />
          </AttachmentMedia>
          <AttachmentContent>
            <AttachmentTitle>design.zip</AttachmentTitle>
            <AttachmentDescription>64%</AttachmentDescription>
          </AttachmentContent>
        </Attachment>
        <Attachment state="processing" orientation="vertical">
          <AttachmentMedia>
            <FileTextIcon />
          </AttachmentMedia>
        </Attachment>
        <Attachment state="error" orientation="vertical">
          <AttachmentMedia>
            <FileWarningIcon />
          </AttachmentMedia>
          <AttachmentContent>
            <AttachmentTitle>model.xlsx</AttachmentTitle>
            <AttachmentDescription>Failed</AttachmentDescription>
          </AttachmentContent>
          <AttachmentActions>
            <AttachmentAction aria-label="Retry">
              <RefreshCwIcon />
            </AttachmentAction>
          </AttachmentActions>
        </Attachment>
        <Attachment state="done" orientation="vertical">
          <AttachmentMedia>
            <CheckIcon />
          </AttachmentMedia>
          <AttachmentContent>
            <AttachmentTitle>uploaded-report-with-a-long-name.pdf</AttachmentTitle>
            <AttachmentDescription>1.8 MB</AttachmentDescription>
          </AttachmentContent>
        </Attachment>
      </div>
    </div>
  )
}
