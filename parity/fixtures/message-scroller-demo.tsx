import { Bubble, BubbleContent } from "@/registry/new-york-v4/ui/bubble"
import { Message, MessageContent } from "@/registry/new-york-v4/ui/message"
import {
  MessageScroller,
  MessageScrollerButton,
  MessageScrollerContent,
  MessageScrollerItem,
  MessageScrollerProvider,
  MessageScrollerViewport,
} from "@/registry/new-york-v4/ui/message-scroller"

const messages = Array.from({ length: 12 }, (_, i) => ({
  id: `m${i + 1}`,
  user: i % 2 === 0,
  text:
    i % 3 === 0
      ? `Message ${i + 1}: a longer message that wraps onto a few lines so the transcript is taller than its viewport and needs to scroll.`
      : `Message ${i + 1}`,
}))

// Port fixture: a transcript taller than its viewport. It opens at the end
// (defaultScrollPosition "end"), with scroll-to-start and scroll-to-end
// buttons that show while there is more to see; user turns are scroll
// anchors.
export default function MessageScrollerDemo() {
  return (
    <div className="h-[420px] w-full max-w-md rounded-xl border">
      <MessageScrollerProvider>
        <MessageScroller>
          <MessageScrollerViewport>
            <MessageScrollerContent className="p-4">
              {messages.map((message) => (
                <MessageScrollerItem
                  key={message.id}
                  messageId={message.id}
                  scrollAnchor={message.user}
                >
                  <Message align={message.user ? "end" : "start"}>
                    <MessageContent>
                      <Bubble variant={message.user ? "default" : "muted"}>
                        <BubbleContent>{message.text}</BubbleContent>
                      </Bubble>
                    </MessageContent>
                  </Message>
                </MessageScrollerItem>
              ))}
            </MessageScrollerContent>
          </MessageScrollerViewport>
          <MessageScrollerButton direction="start" />
          <MessageScrollerButton />
        </MessageScroller>
      </MessageScrollerProvider>
    </div>
  )
}
