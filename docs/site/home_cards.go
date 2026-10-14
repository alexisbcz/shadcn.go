package site

import (
	"fmt"
	"net/http"
	"strconv"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	alertdialog "shadcn.go/components/ui/alert-dialog"
	"shadcn.go/components/ui/badge"
	"shadcn.go/components/ui/breadcrumb"
	"shadcn.go/components/ui/button"
	buttongroup "shadcn.go/components/ui/button-group"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/checkbox"
	dropdownmenu "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/components/ui/empty"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/input"
	inputgroup "shadcn.go/components/ui/input-group"
	"shadcn.go/components/ui/item"
	"shadcn.go/components/ui/progress"
	radiogroup "shadcn.go/components/ui/radio-group"
	selectui "shadcn.go/components/ui/select"
	"shadcn.go/components/ui/separator"
	"shadcn.go/components/ui/sidebar"
	switchui "shadcn.go/components/ui/switch"
	"shadcn.go/components/ui/textarea"
	h "shadcn.go/libhtml"
)

// The cards demo of the home page: app/(app)/(root)/cards. The cards are
// written with styles/base-rhea; the Go components render them in Rhea. The
// skeleton rails (cards/skeleton) only show from 2200px wide and are not
// ported.

// hugeicon renders a hugeicon as HugeiconsIcon does with strokeWidth={2}.
func hugeicon(name string, props ...any) h.Node {
	class := ""
	var rest []any
	for _, p := range props {
		if s, ok := p.(string); ok {
			class = s
		} else {
			rest = append(rest, p)
		}
	}
	return h.El("svg", h.Attr("xmlns", "http://www.w3.org/2000/svg"), h.Attr("width", "24"), h.Attr("height", "24"),
		h.Attr("viewBox", "0 0 24 24"), h.Attr("fill", "none"), h.Attr("color", "currentColor"), h.Class(class),
		h.Attr("stroke-width", "2"), h.Attr("stroke", "currentColor"), h.Group(rest...),
		// Generated from @hugeicons/core-free-icons; never user input.
		h.UnsafeHTML(hugeicons[name]))
}

func cardsDemo(r *http.Request) h.Node {
	return h.Div(h.Data("slot", "demo"),
		h.Class("theme-neutral relative flex w-full max-w-none flex-col gap-(--gap) overflow-hidden bg-muted p-12 pb-0! [--gap:--spacing(8)] 3xl:[--gap:--spacing(8)] min-[1900px]:p-12 min-[1900px]:[--gap:--spacing(10)]! lg:p-6 lg:[--gap:--spacing(6)] dark:bg-background"),
		style.Scope(style.Rhea,
			h.Div(h.Class("relative z-10 mx-auto grid gap-(--gap) **:data-[slot=card]:w-full min-[1400px]:grid-cols-4! min-[1900px]:grid-cols-5! md:max-w-3xl md:grid-cols-2 lg:max-w-none lg:grid-cols-3 xl:max-w-[1600px] 2xl:max-w-[1900px]"),
				h.Div(h.Class("flex flex-col items-start gap-(--gap)"), uiElementsCard(), sidebarNavCard(), savingsTargetsCard()),
				h.Div(h.Class("hidden flex-col gap-(--gap) lg:flex"), contributionHistoryCard(), claimableBalanceCard(), dividendIncomeCard()),
				h.Div(h.Class("hidden flex-col gap-(--gap) min-[1400px]:flex"), newMilestoneCard(), payoutThresholdCard(), accountAccessCard()),
				h.Div(h.Class("hidden flex-col gap-(--gap) md:flex"),
					qrConnectCard(),
					h.Div(h.Class("**:[.text-center.text-xs]:hidden"), messageScrollerDemo(r)),
					paymentsCard(),
				),
				h.Div(h.Class("hidden flex-col gap-(--gap) min-[1900px]:flex"), emptyDistributeTrackCard(), analyticsCard(), notificationSettingsCard(), powerUsageCard()),
			),
		),
		h.Div(h.Class("absolute inset-x-0 top-0 z-1 h-120 bg-linear-to-b from-background via-muted to-transparent dark:hidden")),
		h.Div(h.Class("absolute inset-x-0 bottom-0 z-20 h-48 bg-linear-to-t from-background via-muted/80 to-transparent lg:h-80 xl:h-64 dark:via-background/80")),
	)
}

func uiElementsCard() h.Node {
	return card.Card(h.Class("w-full"),
		card.Content(h.Class("flex flex-col gap-6"),
			h.Div(h.Class("flex gap-2"),
				button.Button("Button ", hugeicon("ArrowRight02Icon", h.Data("icon", "inline-end"))),
				button.Button(button.VariantSecondary, "Secondary"),
				button.Button(button.VariantOutline, "Outline"),
			),
			field.Group(
				field.Field(
					inputgroup.InputGroup(
						inputgroup.Input(h.Placeholder("Name")),
						inputgroup.Addon(inputgroup.AlignInlineEnd, inputgroup.Text(hugeicon("Search01Icon"))),
					),
				),
				field.Field(h.Class("flex-1"), textarea.Textarea(h.Placeholder("Message"), h.Class("resize-none"))),
			),
			h.Div(h.Class("flex items-center gap-2"),
				h.Div(h.Class("flex gap-2"),
					badge.Badge("Badge"),
					badge.Badge(badge.VariantSecondary, "Secondary"),
					badge.Badge(badge.VariantOutline, h.Class("hidden 4xl:flex"), "Outline"),
				),
				radiogroup.RadioGroup(radiogroup.Value("apple"), h.Class("ml-auto flex w-fit gap-3"), h.Aria("label", "Fruit preference"),
					radiogroup.Item(radiogroup.Value("apple"), h.Aria("label", "Apple")),
					radiogroup.Item(radiogroup.Value("banana"), h.Aria("label", "Banana")),
				),
				h.Div(h.Class("flex gap-3"),
					checkbox.Checkbox(checkbox.Checked(true), h.Aria("label", "Enable email alerts")),
					checkbox.Checkbox(h.Class("hidden 4xl:flex"), h.Aria("label", "Enable push alerts")),
				),
				switchui.Switch(switchui.Checked(true), h.Class("flex 4xl:hidden"), h.Aria("label", "Enable compact notifications")),
			),
			h.Div(h.Class("flex items-center gap-4"),
				alertdialog.AlertDialog(
					alertdialog.Trigger(lib.AsChild, button.Button(button.VariantOutline,
						h.Span(h.Class("hidden md:flex style-sera:md:hidden"), "Alert Dialog"),
						h.Span(h.Class("flex md:hidden style-sera:md:flex"), "Dialog"),
					)),
					alertdialog.Content(alertdialog.SizeSm, h.Class("theme-neutral"),
						alertdialog.Header(
							alertdialog.Title("Allow accessory to connect?"),
							alertdialog.Description("Do you want to allow the USB accessory to connect to this device and your data?"),
						),
						alertdialog.Footer(
							alertdialog.Cancel("Don't allow"),
							alertdialog.Action("Allow"),
						),
					),
				),
				buttongroup.ButtonGroup(h.Class("ml-auto"),
					button.Button(button.VariantOutline,
						h.Span(h.Class("style-sera:hidden"), "Button Group"),
						h.Span(h.Class("hidden style-sera:block"), "Group"),
					),
					dropdownmenu.DropdownMenu(
						dropdownmenu.Trigger(lib.AsChild, button.Button(button.VariantOutline, button.SizeIcon, h.Aria("label", "Open quick actions"), hugeicon("ArrowUp01Icon"))),
						dropdownmenu.Content(dropdownmenu.AlignEnd, dropdownmenu.SideTop, h.Class("w-40"),
							dropdownmenu.Group(
								dropdownmenu.Label("Quick Actions"),
								dropdownmenu.Item("Mute Conversation"),
								dropdownmenu.Item("Mark as Read"),
								dropdownmenu.Item("Block User"),
							),
							dropdownmenu.Separator(),
							dropdownmenu.Group(dropdownmenu.Item(dropdownmenu.VariantDestructive, "Delete Conversation")),
						),
					),
				),
				switchui.Switch(switchui.Checked(true), h.Class("hidden 4xl:flex"), h.Aria("label", "Enable advanced setting")),
			),
		),
	)
}

type navEntry struct {
	icon, label string
	active      bool
}

func sidebarSection(label, class string, entries ...navEntry) h.Node {
	var items []any
	for _, e := range entries {
		items = append(items, sidebar.MenuItem(sidebar.MenuButton(sidebar.IsActive(e.active), hugeicon(e.icon), e.label)))
	}
	return card.Card(h.Class(lib.Cn("w-full overflow-hidden rounded-3xl py-0", class)),
		sidebar.Provider(h.Class("min-h-0"),
			sidebar.Sidebar(sidebar.CollapsibleNone, h.Class("w-full bg-transparent"),
				sidebar.Content(h.Class("gap-0 overflow-hidden"),
					sidebar.Group(
						sidebar.GroupLabel(label),
						sidebar.GroupContent(sidebar.Menu(h.Class("gap-1"), h.Group(items...))),
					),
				),
			),
		),
	)
}

func sidebarNavCard() h.Node {
	return h.Div(h.Class("grid w-full grid-cols-2 gap-4 xl:gap-6"),
		sidebarSection("Overview", "xl:col-start-1 xl:row-start-2",
			navEntry{"Analytics01Icon", "Analytics", true},
			navEntry{"ArrowDataTransferHorizontalIcon", "Transactions", false},
			navEntry{"AnalyticsUpIcon", "Investments", false},
			navEntry{"BankIcon", "Accounts", false},
			navEntry{"PieChartIcon", "Spending", false},
		),
		sidebarSection("Planning", "xl:col-start-1 xl:row-start-1",
			navEntry{"File02Icon", "Documents", false},
			navEntry{"Wallet01Icon", "Budget", false},
			navEntry{"ChartBarLineIcon", "Reports", false},
			navEntry{"Target02Icon", "Goals", false},
			navEntry{"Calendar03Icon", "Calendar", false},
		),
		sidebarSection("Support", "flex xl:col-start-2 xl:row-start-1",
			navEntry{"HelpCircleIcon", "Help Center", false},
			navEntry{"BookOpen02Icon", "Docs", false},
			navEntry{"Message01Icon", "Contact Us", false},
			navEntry{"ActivityIcon", "Status", false},
			navEntry{"Globe02Icon", "Community", false},
		),
		sidebarSection("Account", "flex xl:col-start-2 xl:row-start-2",
			navEntry{"UserIcon", "Profile", false},
			navEntry{"CreditCardIcon", "Billing", true},
			navEntry{"Notification03Icon", "Notifications", false},
			navEntry{"ShieldIcon", "Security", false},
			navEntry{"PaintBoardIcon", "Appearance", false},
		),
	)
}

func savingsTargetsCard() h.Node {
	target := func(name, amount string, pct float64, label, achieved, saved string) h.Node {
		return item.Item(h.Role("listitem"), item.VariantMuted, h.Class("flex-col items-stretch"),
			item.Content(h.Class("gap-3"),
				item.Description(h.Class("cn-font-heading text-xs font-medium tracking-wider text-muted-foreground uppercase"), name),
				h.Span(h.Class("text-3xl font-semibold tabular-nums"), amount),
				progress.Progress(progress.Value(pct), h.Aria("label", label)),
			),
			item.Footer(
				h.Span(h.Class("text-sm text-muted-foreground"), achieved),
				h.Span(h.Class("text-sm font-medium tabular-nums"), saved),
			),
		)
	}
	return card.Card(
		card.Header(
			card.Title("Savings Targets"),
			card.Description("Active milestones for 2024 across your portfolio. Monitor how close you are to each savings goal."),
		),
		card.Content(
			item.Group(h.Class("gap-3"),
				target("Retirement", "$420,000", 65, "Retirement savings progress", "65% achieved", "$273,000"),
				target("Real Estate", "$85,000", 32, "Real estate savings progress", "32% achieved", "$27,200"),
			),
		),
		card.Footer(card.Description(h.Class("text-center"), "You have not met your targets for this year.")),
	)
}

func contributionHistoryCard() h.Node {
	data := []struct {
		month  string
		amount float64
	}{{"Dec", 800}, {"Jan", 1100}, {"Feb", 900}, {"Mar", 1300}, {"Apr", 750}}
	var bars []any
	for i, d := range data {
		bars = append(bars, h.Div(h.Class("flex h-full flex-1 flex-col justify-end gap-2"),
			h.Div(h.Data("index", fmt.Sprint(i)),
				h.Class("data-[index=5]:bg-chart-6 min-h-2 rounded-lg data-[index=0]:bg-chart-1 data-[index=1]:bg-chart-2 data-[index=2]:bg-chart-3 data-[index=3]:bg-chart-4 data-[index=4]:bg-chart-5"),
				h.Style("height:"+jsNumber(d.amount/1300*100)+"%")),
			h.Span(h.Class("text-center text-xs text-muted-foreground"), d.month),
		))
	}
	stat := func(class, label, value, sub string) h.Node {
		return item.Item(item.VariantMuted, h.Class(lib.Cn("flex-col items-stretch", class)),
			item.Content(h.Class("gap-1"),
				item.Description(h.Class("text-xs font-medium tracking-wider text-muted-foreground uppercase"), label),
				h.Span(h.Class("cn-font-heading text-base font-semibold"), value),
				h.Span(h.Class("text-sm text-muted-foreground"), sub),
			),
		)
	}
	return card.Card(
		card.Header(card.Title("Contribution History"), card.Description("Last 6 months of activity")),
		card.Content(h.Div(h.Class("flex h-[200px] w-full items-end gap-3"), h.Role("img"), h.Aria("label", "Last 6 months of contribution activity"), h.Group(bars...))),
		card.Content(h.Div(h.Class("grid w-full grid-cols-1 gap-3 xl:grid-cols-2"),
			stat("", "Upcoming", "May 2024", "Scheduled"),
			stat("hidden xl:flex", "Savings Plan", "Accelerated", "Recurring"),
		)),
		card.Footer(button.Button(h.Class("w-full"), "View Full Report")),
	)
}

func claimableBalanceCard() h.Node {
	row := func(label, value, valueClass string) h.Node {
		return h.Div(h.Class("flex items-center justify-between"),
			h.Span(h.Class("text-sm text-muted-foreground"), label),
			h.Span(h.Class(valueClass+" tabular-nums"), value),
		)
	}
	return card.Card(
		card.Header(
			card.Description("Claimable Balance"),
			card.Title(h.Class("text-4xl tabular-nums"), "$1,211.29"),
			badge.Badge(badge.VariantOutline, h.Span(h.Class("size-2 rounded-full bg-yellow-500")), "Pending Setup"),
		),
		card.Content(h.Class("flex flex-1 flex-col justify-end"),
			item.Item(item.VariantMuted, h.Class("flex-col items-stretch"),
				item.Content(h.Class("gap-3"),
					row("Net Royalties", "$1,248.75", "text-sm font-medium"),
					row("Processing Fee", "-$37.46", "text-sm font-medium"),
					separator.Separator(),
					row("Total Ready to Claim", "$1,211.29 USD", "text-sm font-semibold"),
				),
			),
		),
		card.Footer(card.Description("Once your bank is connected, balances over $10.00 are automatically eligible for monthly distribution on the 15th of each month.")),
	)
}

func dismissButton(label string) h.Node {
	return card.Action(button.Button(button.VariantGhost, button.SizeIconSm, h.Class("bg-muted"), h.Aria("label", label), hugeicon("Cancel01Icon")))
}

func dividendIncomeCard() h.Node {
	holdings := []struct {
		name, shares string
		data         [4]float64
	}{
		{"Vanguard", "450 Shares", [4]float64{380, 420, 390, 652}},
		{"S&P 500 VOO", "112 Shares", [4]float64{180, 210, 320, 218}},
		{"Apple AAPL", "85 Shares", [4]float64{60, 70, 120, 90}},
		{"Realty Income", "320 Shares", [4]float64{240, 260, 280, 360}},
	}
	var items []any
	for _, hd := range holdings {
		max := 0.0
		for _, v := range hd.data {
			if v > max {
				max = v
			}
		}
		var bars []any
		for _, v := range hd.data {
			bars = append(bars, h.Div(h.Class("min-h-1 flex-1 rounded-t-sm bg-chart-2"), h.Style("height:"+jsNumber(v/max*100)+"%")))
		}
		items = append(items, item.Item(h.Role("listitem"), item.VariantMuted,
			item.Content(item.Title(hd.name), item.Description(hd.shares)),
			h.Div(h.Class("hidden h-8 w-24 items-end gap-1 md:flex"), h.Role("img"), h.Aria("label", hd.name+" quarterly dividends"), h.Group(bars...)),
		))
	}
	return card.Card(
		card.Header(
			card.Title("Q2 Dividend Income"),
			card.Description("Quarterly dividend payouts across your portfolio holdings."),
			dismissButton("Dismiss dividend income"),
		),
		card.Content(item.Group(h.Group(items...))),
	)
}

func newMilestoneCard() h.Node {
	return card.Card(
		card.Header(
			card.Title("Set a new milestone"),
			card.Description("Define your financial target and we'll help you pace your savings."),
		),
		card.Content(
			field.Group(
				field.Field(field.Label(h.For("goal-name"), "Goal Name"), input.Input(h.ID("goal-name"), h.Placeholder("e.g. New Car, Home Downpayment"))),
				h.Div(h.Class("grid grid-cols-2 gap-3"),
					field.Field(field.Label(h.For("target-amount"), "Target Amount"), input.Input(h.ID("target-amount"), h.Value("$15,000"))),
					field.Field(field.Label(h.For("target-date"), "Target Date"), input.Input(h.ID("target-date"), h.Value("Dec 2025"))),
				),
			),
		),
		card.Footer(h.Class("flex-col gap-2"),
			button.Button(h.Class("w-full"), "Create Goal"),
			button.Button(button.VariantOutline, h.Class("w-full"), "Cancel"),
		),
	)
}

func payoutThresholdCard() h.Node {
	currencies := [][2]string{
		{"USD — United States Dollar", "usd"},
		{"EUR — Euro", "eur"},
		{"GBP — British Pound", "gbp"},
		{"JPY — Japanese Yen", "jpy"},
	}
	var items []any
	for _, c := range currencies {
		items = append(items, selectui.Item(selectui.ItemValue(c[1]), c[0]))
	}
	return card.Card(
		card.Header(
			card.Title("Payout Threshold"),
			card.Description("Set the minimum balance required before a payout is triggered."),
			dismissButton("Dismiss payout threshold"),
		),
		card.Content(
			field.Group(
				field.Field(
					field.Label(h.For("preferred-currency"), "Preferred Currency"),
					selectui.Select(selectui.DefaultValue("usd"),
						selectui.Trigger(h.ID("preferred-currency"), h.Class("w-full"), selectui.Value()),
						selectui.Content(selectui.Group(h.Group(items...))),
					),
				),
				field.Field(
					h.Div(h.Class("flex items-baseline justify-between"),
						field.Label(h.ID("min-payout-label"), "Minimum Payout Amount"),
						h.Span(h.Class("text-2xl font-semibold tabular-nums"), "$2500.00"),
					),
					progress.Progress(progress.Value(25), h.Aria("labelledby", "min-payout-label"), h.Aria("valuetext", "$2,500 of $10,000")),
					h.Div(h.Class("flex items-center justify-between"),
						field.Description("$50 (MIN)"),
						field.Description("$10,000 (MAX)"),
					),
				),
				field.Field(
					field.Label(h.For("payout-notes"), "Notes"),
					textarea.Textarea(h.ID("payout-notes"), h.Placeholder("Add any notes for this payout configuration..."), h.Class("min-h-[100px]")),
				),
			),
		),
		card.Footer(button.Button(h.Class("w-full"), "Save Threshold")),
	)
}

func linkItem(icon, iconClass, title, description, descClass, arrowClass string) h.Node {
	return item.Item(item.VariantMuted, lib.AsChild,
		h.A(h.Href("#"),
			item.Media(item.MediaIcon, hugeicon(icon, iconClass)),
			item.Content(item.Title(title), item.Description(h.Class(descClass), description)),
			hugeicon("ArrowRight01Icon", arrowClass),
		),
	)
}

func accountAccessCard() h.Node {
	return card.Card(
		card.Header(
			card.Title("Account Access"),
			card.Description("Update your credentials or re-authenticate."),
		),
		card.Content(
			field.Group(
				field.Field(field.Label(h.For("email-address"), "Email Address"),
					input.Input(h.ID("email-address"), h.Type("email"), h.Placeholder("artist@studio.inc"))),
				field.Field(
					h.Div(h.Class("flex items-center justify-between"),
						field.Label(h.For("current-password"), "Current Password"),
						h.A(h.Href("#"), h.Class("text-xs font-medium tracking-wider text-muted-foreground uppercase hover:text-foreground"), "Forgot?"),
					),
					input.Input(h.ID("current-password"), h.Type("password"), h.Placeholder("••••••••••••••••••••••••")),
				),
			),
		),
		card.Footer(h.Class("flex-col gap-4"),
			button.Button(h.Class("w-full"), hugeicon("SquareLock02Icon"), "Update Security"),
			linkItem("AlertCircleIcon", "text-destructive", "Danger Zone", "Archive account and remove catalog", "line-clamp-1", "size-4"),
		),
	)
}

var qrCells = []string{
	"111111100101101111111",
	"100000101001001000001",
	"101110101111101011101",
	"101110100100001011101",
	"101110101010101011101",
	"100000100111001000001",
	"111111101010101111111",
	"000000001101000000000",
	"101011111001111010110",
	"010100001110010101001",
	"111010111011101111010",
	"001101000101000010101",
	"110111101111010111011",
	"000000001001010001010",
	"111111101101111101001",
	"100000100010001001111",
	"101110101011101110100",
	"101110100110100010011",
	"101110101000111101110",
	"100000101101000011001",
	"111111101011101101111",
}

func qrConnectCard() h.Node {
	rects := []any{h.El("rect", h.Width("21"), h.Height("21"), h.Attr("fill", "white"))}
	for y, row := range qrCells {
		for x, c := range row {
			if c == '1' {
				rects = append(rects, h.El("rect", h.Attr("x", fmt.Sprint(x)), h.Attr("y", fmt.Sprint(y)), h.Width("1"), h.Height("1")))
			}
		}
	}
	return card.Card(
		card.Content(h.Class("flex justify-center pt-6"),
			h.Div(h.Class("rounded-xl border bg-white p-4"),
				h.El("svg", h.Attr("viewBox", "0 0 21 21"), h.Class("size-40 text-black"), h.Role("img"),
					h.Aria("label", "Connect device QR code"), h.Attr("shape-rendering", "crispEdges"), h.Group(rects...)),
			),
		),
		card.Header(h.Class("text-center"),
			card.Title("Scan to connect your mobile device"),
			card.Description(h.Class("text-balance"), "Open the Ledger mobile app and scan this code to link your device."),
		),
	)
}

func paymentsCard() h.Node {
	entry := func(icon, title, description string) h.Node {
		return h.Div(h.Role("listitem"), h.Class("w-full"),
			linkItem(icon, "", title, description, "", "size-4 shrink-0 text-muted-foreground"))
	}
	return card.Card(
		card.Header(h.Class("flex flex-col gap-3"),
			breadcrumb.Breadcrumb(
				breadcrumb.List(
					breadcrumb.Item(breadcrumb.Link(h.Href("#"), "Home")),
					breadcrumb.Separator(),
					breadcrumb.Item(
						dropdownmenu.DropdownMenu(
							dropdownmenu.Trigger(lib.AsChild, button.Button(button.SizeIconSm, button.VariantGhost, h.Aria("label", "Account options"),
								hugeicon("MoreHorizontalCircle01Icon"), h.Span(h.Class("sr-only"), "Account options"))),
							dropdownmenu.Content(dropdownmenu.AlignStart,
								dropdownmenu.Group(
									dropdownmenu.Item("Profile"),
									dropdownmenu.Item("Statements"),
									dropdownmenu.Item("Documents"),
								),
							),
						),
					),
					breadcrumb.Separator(),
					breadcrumb.Item(breadcrumb.Page("Payments")),
				),
			),
		),
		card.Content(
			item.Group(
				entry("Settings01Icon", "Change transfer limit", "Adjust how much you can send from your balance."),
				entry("Calendar03Icon", "Scheduled transfers", "Set up a transfer to send at a later date."),
				entry("RefreshIcon", "Recurring card payments", "Manage your repeated card transactions."),
			),
		),
	)
}

// jsNumber formats a float as JavaScript's Number#toString does for the
// values of the cards.
func jsNumber(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) }

func emptyDistributeTrackCard() h.Node {
	return card.Card(
		card.Content(
			empty.Empty(h.Class("p-4"),
				empty.Media(empty.MediaIcon, hugeicon("Add01Icon")),
				empty.Header(
					empty.Title("Distribute Track"),
					empty.Description("Upload your first master to start reaching listeners on Spotify, Apple Music, and more."),
				),
				empty.Content(button.Button("Create Release")),
			),
		),
	)
}

func analyticsCard() h.Node {
	return card.Card(card.SizeSm, h.Class("mx-auto w-full max-w-sm data-[size=sm]:pb-0"),
		card.Header(
			card.Title("Analytics"),
			card.Description("418.2K Visitors ", badge.Badge("+10%")),
			card.Action(button.Button(button.VariantOutline, button.SizeSm, "View Analytics")),
		),
		h.El("svg", h.Attr("viewBox", "0 0 100 86"), h.Attr("preserveAspectRatio", "none"), h.Class("aspect-[1/0.35] w-full text-chart-1"),
			h.Role("img"), h.Aria("label", "Visitor trend"),
			h.Path(h.Attr("d", "M0 52L18 40L36 46L54 70L72 50L100 49V86H0Z"), h.Attr("fill", "currentColor"), h.Attr("opacity", "0.28")),
			h.Path(h.Attr("d", "M0 52L18 40L36 46L54 70L72 50L100 49"), h.Attr("fill", "none"), h.Attr("stroke", "currentColor"),
				h.Attr("stroke-width", "1.5"), h.Attr("vector-effect", "non-scaling-stroke")),
		),
	)
}

func notificationSettingsCard() h.Node {
	notifications := []struct {
		id, label, description string
		checked                bool
	}{
		{"transactions", "Transaction alerts", "Deposits, withdrawals, and transfers.", true},
		{"security", "Security alerts", "Login attempts and account changes.", true},
		{"goals", "Goal milestones", "Updates at 25%, 50%, 75%, and 100%.", false},
		{"market", "Market updates", "Daily portfolio summary and price alerts.", false},
	}
	var fields []any
	for _, n := range notifications {
		fields = append(fields, field.Field(field.Horizontal,
			checkbox.Checkbox(h.ID("notify-"+n.id), checkbox.Checked(n.checked)),
			field.Content(
				field.Label(h.For("notify-"+n.id), n.label),
				field.Description(n.description),
			),
		))
	}
	return card.Card(
		card.Header(
			card.Title("Notifications"),
			card.Description("Choose which email and push alerts you want to receive."),
		),
		card.Content(field.Group(h.Group(fields...))),
		card.Footer(button.Button(h.Class("w-full"), "Save Preferences")),
	)
}

func powerUsageCard() h.Node {
	data := []struct {
		hour  string
		usage float64
	}{{"6a", 1.2}, {"8a", 2.8}, {"10a", 3.1}, {"12p", 2.4}, {"2p", 3.4}, {"4p", 2.9}, {"6p", 3.8}, {"8p", 3.2}}
	var bars []any
	for _, d := range data {
		bars = append(bars, h.Div(h.Class("flex h-full flex-1 flex-col justify-end gap-1.5"),
			h.Div(h.Class("min-h-2 rounded-t bg-chart-2"), h.Style("height:"+jsNumber(d.usage/3.8*100)+"%")),
			h.Span(h.Class("text-center text-xs text-muted-foreground"), d.hour),
		))
	}
	stat := func(label, value string) h.Node {
		return h.Div(h.Class("flex flex-col gap-0.5"),
			h.Span(h.Class("text-sm text-muted-foreground"), label),
			h.Span(h.Class("text-lg font-semibold tabular-nums"), value),
		)
	}
	return card.Card(
		card.Header(card.Title("Power Usage"), card.Description("Whole Home")),
		card.Content(h.Class("flex flex-col gap-4"),
			h.Div(h.Class("flex h-[140px] w-full items-end gap-2"), h.Role("img"), h.Aria("label", "Power usage by hour"), h.Group(bars...)),
			separator.Separator(),
			h.Div(h.Class("grid grid-cols-2 gap-4"), stat("Currently Using", "3.4 kW"), stat("Solar Gen", "+1.2 kW")),
		),
	)
}
