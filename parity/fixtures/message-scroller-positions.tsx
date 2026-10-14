import { ChevronsDownIcon } from "lucide-react"

import { Bubble, BubbleContent } from "@/registry/new-york-v4/ui/bubble"
import { Marker, MarkerContent } from "@/registry/new-york-v4/ui/marker"
import {
  MessageScroller,
  MessageScrollerButton,
  MessageScrollerContent,
  MessageScrollerItem,
  MessageScrollerProvider,
  MessageScrollerViewport,
} from "@/registry/new-york-v4/ui/message-scroller"

const rows = Array.from({ length: 10 }, (_, i) => i + 1)

function Rows({ anchors = [] as number[] }) {
  return (
    <>
      {rows.map((row) => (
        <MessageScrollerItem
          key={row}
          messageId={`row-${row}`}
          scrollAnchor={anchors.includes(row)}
        >
          {anchors.includes(row) ? (
            <Marker variant="separator">
              <MarkerContent>Turn {row}</MarkerContent>
            </Marker>
          ) : (
            <Bubble variant="secondary">
              <BubbleContent>Row {row} of the transcript.</BubbleContent>
            </Bubble>
          )}
        </MessageScrollerItem>
      ))}
    </>
  )
}

// Port fixture: the opening positions (start, last-anchor with a last turn
// taller than the viewport, end with autoScroll), the provider's scroll
// options, a button with a custom variant, size, behavior and content, a
// viewport and content with custom classes and props, rows without a
// messageId, and an empty scroller.
export default function MessageScrollerPositions() {
  return (
    <div className="grid w-full max-w-4xl grid-cols-2 gap-4">
      <div className="h-64 rounded-lg border">
        <MessageScrollerProvider defaultScrollPosition="start">
          <MessageScroller>
            <MessageScrollerViewport aria-label="Start">
              <MessageScrollerContent className="gap-3 p-3">
                <Rows />
              </MessageScrollerContent>
            </MessageScrollerViewport>
            <MessageScrollerButton variant="outline" size="sm" behavior="auto">
              <ChevronsDownIcon />
              Latest
            </MessageScrollerButton>
          </MessageScroller>
        </MessageScrollerProvider>
      </div>
      <div className="h-64 rounded-lg border">
        <MessageScrollerProvider
          defaultScrollPosition="last-anchor"
          scrollPreviousItemPeek={24}
          scrollMargin={4}
        >
          <MessageScroller>
            <MessageScrollerViewport>
              <MessageScrollerContent className="gap-3 p-3" spacerClassName="bg-muted">
                <Rows anchors={[3, 6]} />
              </MessageScrollerContent>
            </MessageScrollerViewport>
            <MessageScrollerButton direction="start" />
            <MessageScrollerButton />
          </MessageScroller>
        </MessageScrollerProvider>
      </div>
      <div className="h-64 rounded-lg border">
        <MessageScrollerProvider autoScroll scrollEdgeThreshold={16}>
          <MessageScroller className="bg-muted/30">
            <MessageScrollerViewport tabIndex={-1} preserveScrollOnPrepend={false}>
              <MessageScrollerContent className="gap-3 p-3" role="list">
                <div className="text-xs text-muted-foreground">
                  A row without messageId.
                </div>
                <Rows anchors={[9]} />
              </MessageScrollerContent>
            </MessageScrollerViewport>
            <MessageScrollerButton variant="default" />
          </MessageScroller>
        </MessageScrollerProvider>
      </div>
      <div className="h-64 rounded-lg border">
        <MessageScrollerProvider>
          <MessageScroller>
            <MessageScrollerViewport>
              <MessageScrollerContent className="p-3" />
            </MessageScrollerViewport>
            <MessageScrollerButton />
          </MessageScroller>
        </MessageScrollerProvider>
      </div>
    </div>
  )
}
