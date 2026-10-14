"use client"

// React from the upstream app (a bare "react" import would resolve to
// parity/node_modules and load a second copy of React).
import * as React from "../../ui/apps/v4/node_modules/react"

import { Bubble, BubbleContent } from "@/registry/new-york-v4/ui/bubble"
import { Button } from "@/registry/new-york-v4/ui/button"
import {
  MessageScroller,
  MessageScrollerButton,
  MessageScrollerContent,
  MessageScrollerItem,
  MessageScrollerProvider,
  MessageScrollerViewport,
  useMessageScroller,
  useMessageScrollerScrollable,
  useMessageScrollerVisibility,
} from "@/registry/new-york-v4/ui/message-scroller"

type Row = { id: number; user: boolean; long: boolean }

const initialRows: Row[] = Array.from({ length: 6 }, (_, i) => ({
  id: i + 1,
  user: i % 3 === 0,
  long: i % 2 === 1,
}))

function text(row: Row) {
  return row.long
    ? `Reply ${row.id}: a longer reply that wraps onto a few lines so new turns have to scroll into place.`
    : `${row.user ? "Question" : "Reply"} ${row.id}`
}

function Commands() {
  const { scrollToMessage, scrollToStart, scrollToEnd } = useMessageScroller()
  return (
    <>
      <Button size="sm" variant="outline" onClick={() => scrollToMessage("r2", { align: "center" })}>
        Jump to 2
      </Button>
      <Button size="sm" variant="outline" onClick={() => scrollToMessage("r4", { align: "end", scrollMargin: 8 })}>
        Jump to 4
      </Button>
      <Button size="sm" variant="outline" onClick={() => scrollToStart()}>
        Start
      </Button>
      <Button size="sm" variant="outline" onClick={() => scrollToEnd()}>
        End
      </Button>
    </>
  )
}

function Status() {
  const { currentAnchorId, visibleMessageIds } = useMessageScrollerVisibility()
  const { start, end } = useMessageScrollerScrollable()
  return (
    <p className="text-xs text-muted-foreground">
      Current: {currentAnchorId ?? "none"} · Visible: {visibleMessageIds.join(" ")} · Scrollable: {start ? "start" : "-"} {end ? "end" : "-"}
    </p>
  )
}

// Port fixture: a live transcript with autoScroll. Appending a reply
// follows the end, appending a user turn (scrollAnchor) anchors it near the
// top with the tail spacer, prepending older rows keeps the reading
// position, and the useMessageScroller commands, the visibility and the
// scrollable state are exercised.
export default function MessageScrollerDynamic() {
  const [rows, setRows] = React.useState(initialRows)
  const next = React.useRef(7)
  const older = React.useRef(0)

  return (
    <MessageScrollerProvider autoScroll>
      <div className="flex w-full max-w-xl flex-col gap-3">
        <div className="flex flex-wrap gap-2">
          <Button
            size="sm"
            onClick={() =>
              setRows((r) => [...r, { id: next.current++, user: false, long: true }])
            }
          >
            Add reply
          </Button>
          <Button
            size="sm"
            onClick={() =>
              setRows((r) => [...r, { id: next.current++, user: true, long: false }])
            }
          >
            Add turn
          </Button>
          <Button
            size="sm"
            onClick={() =>
              setRows((r) => [{ id: older.current--, user: false, long: false }, ...r])
            }
          >
            Prepend
          </Button>
          <Commands />
        </div>
        <Status />
        <div className="h-72 rounded-lg border">
          <MessageScroller>
            <MessageScrollerViewport>
              <MessageScrollerContent className="gap-4 p-4">
                {rows.map((row) => (
                  <MessageScrollerItem
                    key={row.id}
                    messageId={`r${row.id}`}
                    scrollAnchor={row.user}
                  >
                    <Bubble variant={row.user ? "default" : "muted"} align={row.user ? "end" : "start"}>
                      <BubbleContent>{text(row)}</BubbleContent>
                    </Bubble>
                  </MessageScrollerItem>
                ))}
              </MessageScrollerContent>
            </MessageScrollerViewport>
            <MessageScrollerButton />
          </MessageScroller>
        </div>
      </div>
    </MessageScrollerProvider>
  )
}
