package message

import (
	"shadcn.go/components/style"
)

// The registry/bases/radix/ui/message.tsx classes, resolved per style. The
// parts render them inside style.Scope; outside, the new-york-v4 classes.
var file = style.NewFile("ui/message.tsx")

var (
	sGroup   = file.Cn("cn-message-group flex min-w-0 flex-col", style.Arg)
	sMessage = file.Cn("cn-message group/message relative flex w-full min-w-0 data-[align=end]:flex-row-reverse", style.Arg)
	sAvatar  = file.Cn("cn-message-avatar flex w-fit shrink-0 items-center justify-center self-end overflow-hidden rounded-full bg-muted", style.Arg)
	sContent = file.Cn("cn-message-content flex w-full min-w-0 flex-col wrap-break-word", style.Arg)
	sHeader  = file.Cn("cn-message-header flex max-w-full min-w-0 items-center", style.Arg)
	sFooter  = file.Cn("cn-message-footer flex max-w-full min-w-0 items-center group-data-[align=end]/message:justify-end", style.Arg)
)
