import {
  Avatar,
  AvatarFallback,
} from "@/registry/new-york-v4/ui/avatar"
import {
  Bubble,
  BubbleContent,
  BubbleGroup,
  BubbleReactions,
} from "@/registry/new-york-v4/ui/bubble"
import { Marker, MarkerContent } from "@/registry/new-york-v4/ui/marker"
import {
  Message,
  MessageAvatar,
  MessageContent,
  MessageFooter,
  MessageGroup,
  MessageHeader,
} from "@/registry/new-york-v4/ui/message"

// Port fixture: Message aligned start (default and explicit) and end, with
// and without avatars, header and footer (the avatar moves up over a
// footer), ghost bubbles (header and footer lose their padding), grouped
// bubbles, MessageGroup and className overrides.
export default function MessageVariants() {
  return (
    <div className="grid w-full max-w-3xl grid-cols-2 gap-8 py-6">
      <MessageGroup>
        <Message>
          <MessageAvatar>
            <Avatar>
              <AvatarFallback>R</AvatarFallback>
            </Avatar>
          </MessageAvatar>
          <MessageContent>
            <MessageHeader>Rabbit</MessageHeader>
            <Bubble variant="muted">
              <BubbleContent>It&apos;s 4:55 PM. On a Friday.</BubbleContent>
            </Bubble>
          </MessageContent>
        </Message>
        <Message align="end">
          <MessageAvatar>
            <Avatar>
              <AvatarFallback>ME</AvatarFallback>
            </Avatar>
          </MessageAvatar>
          <MessageContent>
            <Bubble>
              <BubbleContent>It&apos;s a one-line change.</BubbleContent>
            </Bubble>
            <MessageFooter>Delivered</MessageFooter>
          </MessageContent>
        </Message>
        <Message align="start">
          <MessageAvatar className="bg-primary text-primary-foreground">
            AI
          </MessageAvatar>
          <MessageContent>
            <BubbleGroup>
              <Bubble variant="muted">
                <BubbleContent>It&apos;s always a one-line change.</BubbleContent>
              </Bubble>
              <Bubble variant="muted">
                <BubbleContent>Alright, let me take a look.</BubbleContent>
                <BubbleReactions aria-label="Reactions: thumbs up">
                  <span>👍</span>
                </BubbleReactions>
              </Bubble>
            </BubbleGroup>
            <MessageFooter>Seen</MessageFooter>
          </MessageContent>
        </Message>
      </MessageGroup>
      <MessageGroup className="gap-6">
        <Message>
          <MessageContent>
            <MessageHeader>Assistant</MessageHeader>
            <Bubble variant="ghost">
              <BubbleContent>
                Ghost bubbles take the full width, and the header and footer
                lose their padding to line up with the text.
              </BubbleContent>
            </Bubble>
            <MessageFooter>
              <div>
                Read <span className="font-normal">Yesterday</span>
              </div>
            </MessageFooter>
          </MessageContent>
        </Message>
        <Message align="end" className="text-base">
          <MessageContent className="gap-1">
            <MessageHeader>You</MessageHeader>
            <Bubble variant="secondary">
              <BubbleContent>Send the report to the team.</BubbleContent>
            </Bubble>
            <MessageFooter className="text-foreground">Sent</MessageFooter>
          </MessageContent>
        </Message>
        <Marker role="status">
          <MarkerContent className="shimmer">
            <span className="font-medium">Oliver</span> is typing...
          </MarkerContent>
        </Marker>
      </MessageGroup>
    </div>
  )
}
