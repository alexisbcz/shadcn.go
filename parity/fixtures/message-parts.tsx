import {
  CopyIcon,
  FileTextIcon,
  RefreshCcwIcon,
  ThumbsDownIcon,
  ThumbsUpIcon,
} from "lucide-react"

import {
  Attachment,
  AttachmentContent,
  AttachmentDescription,
  AttachmentMedia,
  AttachmentTitle,
} from "@/registry/new-york-v4/ui/attachment"
import { Bubble, BubbleContent } from "@/registry/new-york-v4/ui/bubble"
import { Button } from "@/registry/new-york-v4/ui/button"
import {
  Message,
  MessageAvatar,
  MessageContent,
  MessageFooter,
  MessageGroup,
  MessageHeader,
} from "@/registry/new-york-v4/ui/message"

// Port fixture: MessageHeader and MessageFooter on both sides (padded, and
// flush against a ghost bubble), footer actions, an attachment in a
// message, a MessageGroup, an avatar with text, and className overrides.
export default function MessageParts() {
  return (
    <div className="grid w-full max-w-3xl grid-cols-2 gap-10 py-8">
      <div className="flex flex-col gap-8">
        <Message>
          <MessageContent>
            <MessageHeader>Olivia</MessageHeader>
            <Bubble variant="muted">
              <BubbleContent>I already checked the logs.</BubbleContent>
            </Bubble>
          </MessageContent>
        </Message>
        <Message align="end">
          <MessageContent>
            <MessageHeader>You</MessageHeader>
            <Bubble>
              <BubbleContent>
                Send the report to the team. Ping @shadcn if you need help.
              </BubbleContent>
            </Bubble>
            <MessageFooter>
              <div>
                Read <span className="font-normal">Yesterday</span>
              </div>
            </MessageFooter>
          </MessageContent>
        </Message>
        <Message>
          <MessageContent>
            <MessageHeader>Assistant</MessageHeader>
            <Bubble variant="ghost">
              <BubbleContent>
                The install failure is coming from the workspace package.
              </BubbleContent>
            </Bubble>
            <MessageFooter>
              <Button variant="ghost" size="icon" aria-label="Copy" title="Copy">
                <CopyIcon />
              </Button>
              <Button variant="ghost" size="icon" aria-label="Like" title="Like">
                <ThumbsUpIcon />
              </Button>
              <Button
                variant="ghost"
                size="icon"
                aria-label="Dislike"
                title="Dislike"
              >
                <ThumbsDownIcon />
              </Button>
            </MessageFooter>
          </MessageContent>
        </Message>
      </div>
      <div className="flex flex-col gap-8">
        <Message align="end">
          <MessageContent>
            <Bubble>
              <BubbleContent>Okay drop me a link. Taking a look...</BubbleContent>
            </Bubble>
            <MessageFooter className="gap-2">
              <span className="font-normal text-destructive">Failed to send</span>
              <Button variant="ghost" size="icon-xs" title="Retry" aria-label="Retry">
                <RefreshCcwIcon />
              </Button>
            </MessageFooter>
          </MessageContent>
        </Message>
        <Message align="end">
          <MessageAvatar className="size-8 text-xs font-medium">ME</MessageAvatar>
          <MessageContent>
            <Attachment>
              <AttachmentMedia>
                <FileTextIcon />
              </AttachmentMedia>
              <AttachmentContent>
                <AttachmentTitle>quarterly-report.pdf</AttachmentTitle>
                <AttachmentDescription>PDF · 2.4 MB</AttachmentDescription>
              </AttachmentContent>
            </Attachment>
            <Bubble>
              <BubbleContent>Here is the report.</BubbleContent>
            </Bubble>
          </MessageContent>
        </Message>
        <MessageGroup className="gap-1">
          <Message>
            <MessageAvatar className="size-8 text-xs font-medium">AI</MessageAvatar>
            <MessageContent>
              <Bubble variant="secondary">
                <BubbleContent>First message in a group.</BubbleContent>
              </Bubble>
            </MessageContent>
          </Message>
          <Message className="pl-10">
            <MessageContent className="gap-1">
              <Bubble variant="secondary">
                <BubbleContent>
                  A second message without an avatar, and a long text that
                  wraps: averyveryverylongwordwithoutanybreakopportunityatall.
                </BubbleContent>
              </Bubble>
            </MessageContent>
          </Message>
        </MessageGroup>
      </div>
    </div>
  )
}
