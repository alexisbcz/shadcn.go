package site

import (
	"strings"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/accordion"
	"shadcn.go/components/ui/badge"
	"shadcn.go/components/ui/breadcrumb"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/calendar"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/chart"
	"shadcn.go/components/ui/checkbox"
	"shadcn.go/components/ui/combobox"
	dropdownmenu "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/components/ui/empty"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/input"
	inputgroup "shadcn.go/components/ui/input-group"
	"shadcn.go/components/ui/item"
	"shadcn.go/components/ui/label"
	nativeselect "shadcn.go/components/ui/native-select"
	"shadcn.go/components/ui/progress"
	radiogroup "shadcn.go/components/ui/radio-group"
	selectui "shadcn.go/components/ui/select"
	"shadcn.go/components/ui/separator"
	"shadcn.go/components/ui/sidebar"
	"shadcn.go/components/ui/skeleton"
	"shadcn.go/components/ui/slider"
	"shadcn.go/components/ui/spinner"
	switchui "shadcn.go/components/ui/switch"
	"shadcn.go/components/ui/table"
	"shadcn.go/components/ui/tabs"
	"shadcn.go/components/ui/textarea"
	togglegroup "shadcn.go/components/ui/toggle-group"
	h "shadcn.go/libhtml"
)

// The preview-02 block of the create designer
// (registry/bases/<base>/blocks/preview-02): one card per file upstream.
// The radix, base and aria sources render the same; the Go components
// render them in the design system's style.

// previewGrid is the 7-column grid the preview blocks share.
func previewGrid(columns ...h.Node) h.Node {
	cols := make([]any, len(columns))
	for i, c := range columns {
		cols[i] = c
	}
	return h.Div(h.Class("overflow-x-auto overflow-y-hidden bg-muted contain-[paint] [--gap:--spacing(4)] 3xl:[--gap:--spacing(12)] md:[--gap:--spacing(10)] dark:bg-background style-lyra:md:[--gap:--spacing(6)] style-mira:md:[--gap:--spacing(6)]"),
		h.Div(h.Class("flex w-full min-w-max justify-center"),
			h.Div(h.Class("grid w-[2400px] grid-cols-7 items-start gap-(--gap) bg-muted p-(--gap) md:w-[3000px] dark:bg-background style-lyra:md:w-[2600px] style-mira:md:w-[2600px] *:[div]:gap-(--gap)"),
				h.Data("slot", "capture-target"),
				h.Group(cols...),
			),
		),
	)
}

// previewColumn is a grid column; span2 for the double-width one.
func previewColumn(span2 bool, children ...any) h.Node {
	if span2 {
		return h.Div(h.Class("col-span-2 flex flex-col p-1 [contain-intrinsic-size:760px_1200px] [content-visibility:auto]"), h.Group(children...))
	}
	return h.Div(h.Class("flex flex-col p-1 [contain-intrinsic-size:380px_1200px] [content-visibility:auto]"), h.Group(children...))
}

// previewColumnPx is a column of the base and aria preview block (p-px).
func previewColumnPx(children ...any) h.Node {
	return h.Div(h.Class("flex flex-col p-px [contain-intrinsic-size:380px_1200px] [content-visibility:auto]"), h.Group(children...))
}

func preview02(ds designSystem) h.Node {
	return previewGrid(
		previewColumn(false, contributionHistory(ds), emptyDistributeTrack(), qrConnect(), dividendIncome(ds), indexInvesting(), syncingState()),
		previewColumn(false, payoutThreshold(), claimableBalance(), preferences(), savingsProgress(), kitchenIsland()),
		previewColumn(true, savingsTargets(), recentTransactions(ds),
			h.Div(h.Class("grid grid-cols-2 items-start gap-(--gap)"),
				h.Div(h.Class("flex flex-col gap-(--gap)"), sidebarNav(), faq(ds)),
				h.Div(h.Class("flex flex-col gap-(--gap)"), payments(), frontDoor()),
			),
			releaseCatalog(),
		),
		previewColumn(false, accountAccess(), cardOverview(), transferFunds(), coverArt(), loadingCard()),
		previewColumn(false, receivingMethod(), powerUsage(ds), emptyConnectBank(), upcomingPayments(), rollerShades()),
		previewColumn(false, stockPerformance(), emptyExploreCatalog(), newMilestone(), socialLinks(), notificationSettings()),
	)
}

// --- shared bits ---

const (
	icX            = "XIcon|IconX|Cancel01Icon|XIcon|RiCloseLine"
	icPlus         = "PlusIcon|IconPlus|Add01Icon|PlusIcon|RiAddLine"
	icSearch       = "SearchIcon|IconSearch|Search01Icon|MagnifyingGlassIcon|RiSearchLine"
	icChevronRight = "ChevronRightIcon|IconChevronRight|ArrowRight01Icon|CaretRightIcon|RiArrowRightSLine"
	icWallet       = "WalletIcon|IconWallet|Wallet01Icon|WalletIcon|RiWalletLine"
	icCreditCard   = "CreditCardIcon|IconCreditCard|CreditCardIcon|CreditCardIcon|RiBankCardLine"
)

// rounded is the cards' isRounded: no rounded bars in lyra and sera.
func rounded(ds designSystem) bool { return ds.Style != "lyra" && ds.Style != "sera" }

func closeAction() h.Node {
	return card.Action(button.Button(button.VariantGhost, button.SizeIconSm, h.Class("bg-muted"), ip(icX, "")))
}

// kv is a "flex items-center justify-between" label/value row.
func kv(label, value, valueClass string, labelClass ...string) h.Node {
	lc := "text-sm text-muted-foreground"
	if len(labelClass) > 0 {
		lc = labelClass[0]
	}
	return h.Div(h.Class("flex items-center justify-between"),
		h.Span(h.Class(lc), label),
		h.Span(h.Class(valueClass), value),
	)
}

func currencySelect(id string) h.Node {
	return selectui.Select(selectui.DefaultValue("usd"),
		selectui.Trigger(h.ID(id), h.Class("w-full"), selectui.Value()),
		selectui.Content(selectui.Group(
			selectui.Item(selectui.ItemValue("usd"), "USD — United States Dollar"),
			selectui.Item(selectui.ItemValue("eur"), "EUR — Euro"),
			selectui.Item(selectui.ItemValue("gbp"), "GBP — British Pound"),
			selectui.Item(selectui.ItemValue("jpy"), "JPY — Japanese Yen"),
		)),
	)
}

func barChart(config chart.Config, class string, rows []chart.Row, margin map[string]any, children ...*chart.Node) h.Node {
	args := []any{chart.Data(rows), chart.P("margin", margin)}
	for _, c := range children {
		args = append(args, c)
	}
	return chart.Container(config, h.Class(class), chart.BarChart(args...))
}

func barRadius(on bool, r float64) chart.Prop {
	if on {
		return chart.Radii(r, r, 0, 0)
	}
	return chart.Radius(0)
}

func noCursorTooltip(extra ...any) *chart.Node {
	return chart.Tooltip(chart.Cursor(false), chart.TooltipContent(append([]any{chart.HideLabel()}, extra...)...))
}

func zeroMargin(top float64) map[string]any {
	return map[string]any{"top": top, "right": 0, "bottom": 0, "left": 0}
}

// --- column 1 ---

func contributionHistory(ds designSystem) h.Node {
	rows := []chart.Row{
		{"month": "Dec", "amount": 800}, {"month": "Jan", "amount": 1100}, {"month": "Feb", "amount": 900},
		{"month": "Mar", "amount": 1300}, {"month": "Apr", "amount": 750}, {"month": "May", "amount": 1400},
	}
	stat := func(label, value, sub string) h.Node {
		return item.Item(item.VariantMuted, h.Class("flex-col items-stretch"),
			item.Content(h.Class("gap-1"),
				item.Description(h.Class("text-xs font-medium tracking-wider text-muted-foreground uppercase"), label),
				h.Span(h.Class("cn-font-heading text-lg font-semibold"), value),
				h.Span(h.Class("text-sm text-muted-foreground"), sub),
			),
		)
	}
	return card.Card(
		card.Header(card.Title("Contribution History"), card.Description("Last 6 months of activity")),
		card.Content(
			barChart(chart.Config{{Key: "amount", Label: "Contribution", Color: "var(--chart-2)"}}, "h-[200px] w-full", rows,
				map[string]any{"left": 0, "right": 0, "top": 8, "bottom": 0},
				chart.XAxis(chart.DataKey("month"), chart.TickLine(false), chart.TickMargin(8), chart.AxisLine(false)),
				noCursorTooltip(chart.ClassName("min-w-40")),
				chart.Bar(chart.DataKey("amount"), chart.Fill("var(--color-amount)"), barRadius(rounded(ds), 6), chart.P("maxBarSize", 40)),
			),
		),
		card.Content(
			h.Div(h.Class("grid w-full grid-cols-1 gap-3 md:grid-cols-2"),
				stat("Upcoming", "May 25, 2024", "$1,000 scheduled"),
				stat("Auto-Save Plan", "Accelerated", "Recurring weekly"),
			),
		),
		card.Footer(button.Button(h.Class("w-full"), "View Full Report")),
	)
}

func emptyCard(icon, title, description, action string) h.Node {
	return card.Card(
		card.Content(
			empty.Empty(h.Class("p-4"),
				empty.Media(empty.MediaIcon, ip(icon, "")),
				empty.Header(empty.Title(title), empty.Description(description)),
				empty.Content(button.Button(action)),
			),
		),
	)
}

func emptyDistributeTrack() h.Node {
	return emptyCard(icPlus, "Distribute Track", "Upload your first master to start reaching listeners on Spotify, Apple Music, and more.", "Create Release")
}

func qrConnect() h.Node {
	return card.Card(
		card.Content(h.Class("flex justify-center pt-6"),
			// react-qr-code's output, rendered by gen/create.tsx.
			h.Div(h.Class("rounded-xl border bg-white p-4"), h.UnsafeHTML(createDataOnce().QRConnect)),
		),
		card.Header(h.Class("text-center"),
			card.Title("Scan to connect your mobile device"),
			card.Description("Open the Ledger mobile app and scan this code to link your device."),
		),
		card.Footer(button.Button(button.VariantSecondary, h.Class("w-full"), "Got it")),
	)
}

func dividendIncome(ds designSystem) h.Node {
	holdings := []struct {
		name, shares, amount string
		data                 [4]float64
	}{
		{"Vanguard VIG", "450 Shares", "$1,842.10", [4]float64{380, 420, 390, 652}},
		{"S&P 500 VOO", "112 Shares", "$928.40", [4]float64{180, 210, 320, 218}},
		{"Apple AAPL", "85 Shares", "$340.00", [4]float64{60, 70, 120, 90}},
		{"Realty Income", "320 Shares", "$1,139.50", [4]float64{240, 260, 280, 360}},
	}
	var items []any
	for _, hd := range holdings {
		var rows []chart.Row
		for i, v := range hd.data {
			rows = append(rows, chart.Row{"q": "Q" + string(rune('1'+i)), "value": v})
		}
		items = append(items, item.Item(item.VariantMuted,
			item.Content(item.Title(hd.name), item.Description(hd.shares)),
			barChart(chart.Config{{Key: "value", Label: "Dividend", Color: "var(--chart-2)"}}, "hidden h-8 w-24 md:block", rows, zeroMargin(0),
				noCursorTooltip(),
				chart.Bar(chart.DataKey("value"), chart.Fill("var(--color-value)"), barRadius(rounded(ds), 3)),
			),
			h.Span(h.Class("hidden text-sm font-semibold tabular-nums md:block"), hd.amount),
		))
	}
	return card.Card(
		card.Header(
			card.Title("Q2 Dividend Income"),
			card.Description("Quarterly dividend payouts across your portfolio holdings."),
			closeAction(),
		),
		card.Content(item.Group(items...)),
	)
}

func indexInvesting() h.Node {
	return card.Card(
		card.Header(card.Title("Dollar-Cost Averaging"), card.Description("A strategy for building wealth over time.")),
		card.Content(
			card.Description(h.Class("mt-3 text-sm leading-relaxed style-sera:mt-0"),
				h.A(h.Href("#"), h.Class("underline underline-offset-4 hover:text-primary"), "Over time"),
				", this smooths out the average cost of your investments. When prices drop, your fixed amount buys more shares. When prices rise, you buy fewer. The result is a lower average cost per share compared to lump-sum investing during volatile periods.",
			),
		),
	)
}

func syncingState() h.Node {
	return card.Card(
		card.Content(h.Class("p-0"),
			empty.Empty(h.Class("p-4"),
				empty.Header(
					empty.Media(empty.MediaIcon, spinner.Spinner()),
					empty.Title("Syncing your accounts"),
					empty.Description("We're pulling in your latest transactions. This usually takes a few seconds."),
				),
				empty.Content(button.Button(button.VariantOutline, "Cancel")),
			),
		),
	)
}

// --- column 2 ---

func payoutThreshold() h.Node {
	return card.Card(
		card.Header(
			card.Title("Payout Threshold"),
			card.Description("Set the minimum balance required before a payout is triggered."),
			closeAction(),
		),
		card.Content(
			field.Group(
				field.Field(field.Label(h.For("preferred-currency"), "Preferred Currency"), currencySelect("preferred-currency")),
				field.Field(
					h.Div(h.Class("flex items-baseline justify-between"),
						field.Label(h.For("min-payout"), "Minimum Payout Amount"),
						h.Span(h.Class("text-2xl font-semibold tabular-nums"), h.Data("create-slider-value", "min-payout"), "$2500.00"),
					),
					slider.Slider(h.ID("min-payout"), slider.DefaultValue(2500), slider.Min(50), slider.Max(10000), slider.Step(50)),
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

func claimableBalance() h.Node {
	return card.Card(
		card.Header(
			card.Description("Claimable Balance"),
			card.Title(h.Class("text-5xl tabular-nums"), "$0.00"),
			badge.Badge(badge.VariantOutline, h.Span(h.Class("size-2 rounded-full bg-yellow-500")), "Pending Setup"),
		),
		card.Content(h.Class("flex flex-1 flex-col justify-end"),
			item.Item(item.VariantMuted, h.Class("flex-col items-stretch"),
				item.Content(h.Class("gap-3"),
					kv("Net Royalties", "$0.00", "text-sm font-medium tabular-nums"),
					kv("Processing Fee", "-$0.00", "text-sm font-medium tabular-nums"),
					separator.Separator(),
					kv("Total Ready to Claim", "$0.00 USD", "text-sm font-semibold tabular-nums"),
				),
			),
		),
		card.Footer(card.Description("Once your bank is connected, balances over $10.00 are automatically eligible for monthly distribution on the 15th of each month.")),
	)
}

func preferences() h.Node {
	toggle := func(id, label, description string) h.Node {
		return field.Field(field.Horizontal,
			field.Content(field.Label(h.For(id), label), field.Description(description)),
			switchui.Switch(h.ID(id), switchui.Checked(true)),
		)
	}
	return card.Card(
		card.Header(
			card.Title("Preferences"),
			card.Description("Manage your account settings and notifications."),
			closeAction(),
		),
		card.Content(
			field.Group(
				field.Field(field.Label(h.For("default-currency"), "Default Currency"), currencySelect("default-currency")),
				field.Separator(h.Class("-my-4 style-sera:hidden")),
				toggle("public-statistics", "Public Statistics", "Allow others to see your total stream count and listening activity"),
				field.Separator(h.Class("-my-4 style-sera:hidden")),
				toggle("email-notifications", "Email Notifications", "Monthly royalty reports and distribution updates"),
			),
		),
		card.Footer(
			button.Button(button.VariantOutline, "Reset"),
			button.Button(h.Class("ml-auto"), "Save Preferences"),
		),
	)
}

func savingsProgress() h.Node {
	config := chart.Config{
		{Key: "saved", Label: "Saved", Color: "var(--chart-2)"},
		{Key: "remaining", Label: "Remaining", Color: "var(--chart-1)"},
	}
	rows := []chart.Row{
		{"name": "saved", "value": 24000, "fill": "var(--color-saved)"},
		{"name": "remaining", "value": 6000, "fill": "var(--color-remaining)"},
	}
	line := func(label, value, class string) h.Node {
		return h.Div(h.Class("flex w-full items-center justify-between py-3"),
			h.Span(h.Class("text-sm text-muted-foreground"), label),
			h.Span(h.Class(class), value),
		)
	}
	return card.Card(
		card.Content(
			chart.Container(config, h.Class("mx-auto aspect-square max-h-[220px]"),
				chart.PieChart(
					noCursorTooltip(),
					chart.Pie(chart.Data(rows), chart.DataKey("value"), chart.NameKey("name"),
						chart.InnerRadius(70), chart.OuterRadius(95), chart.StrokeWidth(0), chart.P("startAngle", 90), chart.P("endAngle", -270),
						chart.Label(chart.Content(chart.Fn{Name: "pieCenterLines", Args: []any{
							[]any{"$24,000", "fill-foreground text-2xl font-bold", -12},
							[]any{"80% of $30,000", "fill-muted-foreground text-xs", 12},
						}})),
					),
				),
			),
		),
		card.Footer(h.Class("flex-col gap-0"),
			line("Projected Finish", "October 2024", "text-sm font-semibold"),
			separator.Separator(),
			line("Monthly Average", "$1,250", "text-sm font-semibold tabular-nums"),
			separator.Separator(),
			line("Top Contributor", "Auto-Transfer", "text-sm font-semibold"),
		),
	)
}

func kitchenIsland() h.Node {
	control := func(icon, title string, value float64, class string) h.Node {
		s := slider.Slider(slider.DefaultValue(value), slider.Max(100), h.Data("create-scene-slider", strings.ToLower(strings.ReplaceAll(title, " ", ""))))
		if class != "" {
			s.Children(h.Class(class))
		}
		return item.Item(item.SizeSm, item.VariantOutline,
			item.Media(item.MediaIcon, ip(icon, "")),
			item.Content(h.Class("flex-row items-center gap-3"), item.Title(h.Class("shrink-0"), title)),
			item.Actions(h.Class("flex-1"), s),
		)
	}
	return card.Card(h.Data("create-kitchen-island", ""),
		card.Header(
			card.Title("Kitchen Island"),
			card.Description("Hue Color Ambient"),
			card.Action(switchui.Switch(switchui.Checked(true))),
		),
		card.Content(h.Class("flex flex-col gap-4"),
			h.Div(h.Class("flex flex-col gap-2"),
				h.Span(h.Class("sr-only"), "Scenes"),
				togglegroup.ToggleGroup(togglegroup.TypeSingle, togglegroup.Value("cooking"), togglegroup.VariantOutline, togglegroup.Spacing(1), h.Class("flex-wrap"),
					togglegroup.Item(togglegroup.Value("cooking"), "Cooking"),
					togglegroup.Item(togglegroup.Value("dining"), "Dining"),
					togglegroup.Item(togglegroup.Value("nightlight"), "Nightlight"),
					togglegroup.Item(togglegroup.Value("focus"), "Focus"),
				),
			),
			item.Group(
				control("SunIcon|IconSun|Sun03Icon|SunIcon|RiSunLine", "Brightness", 90, "w-full"),
				control("ThermometerIcon|IconThermometer|ThermometerWarmIcon|ThermometerIcon|RiThermometerLine", "Color Temp", 70, ""),
				control("Volume2Icon|IconVolume|VolumeHighIcon|SpeakerHighIcon|RiVolumeUpLine", "Volume", 30, ""),
				control("TimerIcon|IconClock|Clock03Icon|TimerIcon|RiTimerLine", "Fade", 0, ""),
			),
		),
	)
}

// --- column 3 (double) ---

func savingsTargets() h.Node {
	target := func(label, amount string, value float64, achieved, saved string) h.Node {
		return item.Item(item.VariantMuted, h.Class("flex-col items-stretch"),
			item.Content(h.Class("gap-3"),
				item.Description(h.Class("cn-font-heading text-xs font-medium tracking-wider text-muted-foreground uppercase"), label),
				h.Span(h.Class("text-3xl font-semibold tabular-nums"), amount),
				progress.Progress(progress.Value(value)),
			),
			item.Footer(
				h.Span(h.Class("text-sm text-muted-foreground"), achieved),
				h.Span(h.Class("text-sm font-medium tabular-nums"), saved),
			),
		)
	}
	return h.Div(h.Class("grid grid-cols-2 gap-(--gap)"),
		card.Card(
			card.Header(
				card.Title("Savings Targets"),
				card.Description("Active milestones for 2024"),
				card.Action(button.Button(button.VariantOutline, button.SizeSm, "New Goal")),
			),
			card.Content(
				item.Group(h.Class("gap-3"),
					target("Retirement", "$420,000", 65, "65% achieved", "$273,000"),
					target("Real Estate", "$85,000", 32, "32% achieved", "$27,200"),
				),
			),
			card.Footer(card.Description(h.Class("text-center"), "You have not met your targets for this year.")),
		),
		card.Card(
			card.Header(card.Title("Buy Investment")),
			card.Content(h.Class("flex flex-1 flex-col gap-3"),
				field.Group(h.Class("flex-1"),
					field.Field(
						field.Label(h.For("invest-amount"), "Amount to Invest"),
						inputgroup.InputGroup(
							inputgroup.Addon(inputgroup.Text("$")),
							inputgroup.Input(h.ID("invest-amount"), h.Value("1,000.00")),
						),
					),
					field.Field(
						field.Label(h.For("invest-type"), "Order Type"),
						nativeselect.NativeSelect(h.ID("invest-type"),
							nativeselect.Option(h.Value("market"), h.Bool("selected"), "Market Order"),
							nativeselect.Option(h.Value("limit"), "Limit Order"),
							nativeselect.Option(h.Value("stop"), "Stop Order"),
						),
						field.Description("Market orders execute at the current price."),
					),
					h.Div(h.Class("flex flex-col gap-2"),
						kv("Estimated Shares", "1.95", "text-sm font-semibold tabular-nums"),
						kv("Buying Power", "$12,450.00", "text-sm font-semibold tabular-nums"),
					),
				),
			),
			card.Footer(h.Class("flex-col gap-3"),
				button.Button(h.Class("w-full"), "Review Order"),
				card.Description(h.Class("text-center"), "Trades are typically executed within minutes during market hours."),
			),
		),
	)
}

func recentTransactions(ds designSystem) h.Node {
	rows := []struct{ icon, name, category, when, amount, amountClass string }{
		{"CoffeeIcon|IconCoffee|CoffeeIcon|CoffeeIcon|RiCupLine", "Blue Bottle Coffee", "Food & Drink", "Today, 10:24 AM", "-$6.50", ""},
		{"ShoppingCartIcon|IconShoppingCart|ShoppingCart01Icon|ShoppingCartIcon|RiShoppingCartLine", "Whole Foods Market", "Groceries", "Yesterday", "-$142.30", ""},
		{icWallet, "Stripe Payout", "Income", "Oct 12", "+$4,200.00", "text-emerald-500"},
		{"CarIcon|IconCar|Car01Icon|CarIcon|RiCarLine", "Uber Technologies", "Transport", "Oct 11", "-$24.10", ""},
		{"TvIcon|IconDeviceTv|Tv01Icon|TelevisionIcon|RiTvLine", "Netflix Subscription", "Entertainment", "Oct 10", "-$19.99", ""},
	}
	var trs []any
	for _, r := range rows {
		trs = append(trs, table.Row(
			table.Cell(h.Class("w-10"),
				h.Div(h.Class("flex size-10 items-center justify-center rounded-lg bg-muted"), ip(r.icon, "size-4 shrink-0")),
			),
			table.Cell(
				h.Div(h.Class("flex flex-col"),
					h.Span(h.Class("font-medium"), r.name),
					h.Span(h.Class("text-sm text-muted-foreground"), r.category),
				),
			),
			table.Cell(h.Class("text-sm text-muted-foreground"), r.when),
			table.Cell(h.Class("text-right"), h.Span(h.Class(lib.Cn("text-sm font-semibold tabular-nums", r.amountClass)), r.amount)),
			table.Cell(h.Class("w-8"),
				dropdownmenu.DropdownMenu(
					dropdownmenu.Trigger(lib.AsChild, button.Button(button.VariantGhost, button.SizeIconSm,
						ip("MoreHorizontalIcon|IconDotsVertical|MoreVerticalCircle01Icon|DotsThreeIcon|RiMore2Line", ""))),
					dropdownmenu.Content(dropdownmenu.AlignEnd,
						dropdownmenu.Item("View details"),
						dropdownmenu.Item("Add note"),
						dropdownmenu.Item("Categorize"),
						dropdownmenu.Separator(),
						dropdownmenu.Item("Dispute"),
					),
				),
			),
		))
	}
	return card.Card(
		card.Header(
			card.Title("Recent Transactions"),
			card.Description("Your latest account activity."),
			card.Action(button.Button(button.VariantOutline, button.SizeSm, "View All")),
		),
		card.Content(table.Table(recentTransactionsHeader(ds), table.Body(trs...))),
	)
}

// recentTransactionsHeader is the aria source's table header (React Aria
// tables need one); the radix and base sources have none.
func recentTransactionsHeader(ds designSystem) h.Node {
	if ds.Base != "aria" {
		return nil
	}
	return table.Header(table.Row(table.Head(), table.Head("Merchant"), table.Head("Date"), table.Head("Amount"), table.Head()))
}

type createNavEntry struct {
	icon, label string
	active      bool
}

func createSidebarNavCard(groups ...struct {
	label   string
	entries []createNavEntry
}) h.Node {
	var content []any
	for i, g := range groups {
		if i > 0 {
			content = append(content, sidebar.Separator(h.Class("w-auto!")))
		}
		var items []any
		for _, e := range g.entries {
			items = append(items, sidebar.MenuItem(sidebar.MenuButton(sidebar.IsActive(e.active), ip(e.icon, ""), e.label)))
		}
		class := "pb-1"
		if i > 0 {
			class = "pt-1"
		}
		content = append(content, sidebar.Group(h.Class(class),
			sidebar.GroupLabel(g.label),
			sidebar.GroupContent(sidebar.Menu(items...)),
		))
	}
	return card.Card(h.Class("overflow-hidden py-0"),
		sidebar.Provider(h.Class("min-h-0"),
			sidebar.Sidebar(sidebar.CollapsibleNone, h.Class("w-full bg-transparent"),
				sidebar.Content(h.Class("gap-0"), h.Group(content...)),
			),
		),
	)
}

func sidebarNav() h.Node {
	type group = struct {
		label   string
		entries []createNavEntry
	}
	return h.Div(h.Class("grid grid-cols-2 items-start gap-6"),
		createSidebarNavCard(
			group{"Overview", []createNavEntry{
				{"LayoutDashboardIcon|IconLayoutDashboard|DashboardSquare01Icon|SquaresFourIcon|RiDashboardLine", "Dashboard", true},
				{"ArrowLeftRightIcon|IconArrowsLeftRight|ArrowDataTransferHorizontalIcon|ArrowsLeftRightIcon|RiArrowLeftRightLine", "Transactions", false},
				{"TrendingUpIcon|IconTrendingUp|AnalyticsUpIcon|TrendUpIcon|RiLineChartLine", "Investments", false},
				{"Building2Icon|IconBuildingBank|BankIcon|BankIcon|RiBankLine", "Accounts", false},
				{"PieChartIcon|IconChartPie|PieChartIcon|ChartPieIcon|RiPieChartLine", "Spending", false},
			}},
			group{"Planning", []createNavEntry{
				{"TargetIcon|IconTarget|Target02Icon|TargetIcon|RiFocus3Line", "Goals", false},
				{icWallet, "Budget", false},
				{"FileBarChartIcon|IconReportAnalytics|ChartBarLineIcon|ChartBarIcon|RiBarChartLine", "Reports", false},
				{"FileTextIcon|IconFileText|File02Icon|FileTextIcon|RiFileTextLine", "Documents", false},
			}},
		),
		createSidebarNavCard(
			group{"Account", []createNavEntry{
				{"UserIcon|IconUser|UserIcon|UserIcon|RiUserLine", "Profile", false},
				{icCreditCard, "Billing", true},
				{"BellIcon|IconBell|Notification03Icon|BellIcon|RiBellLine", "Notifications", false},
				{"ShieldIcon|IconShield|ShieldIcon|ShieldIcon|RiShieldLine", "Security", false},
				{"PaintbrushIcon|IconPalette|PaintBoardIcon|PaletteIcon|RiPaletteLine", "Appearance", false},
			}},
			group{"Support", []createNavEntry{
				{"CircleHelpIcon|IconHelp|HelpCircleIcon|QuestionIcon|RiQuestionLine", "Help Center", false},
				{"MessageSquareIcon|IconMessage|Message01Icon|ChatIcon|RiChat1Line", "Contact Us", false},
				{"BookOpenIcon|IconBook|BookOpen02Icon|BookOpenIcon|RiBookOpenLine", "Documentation", false},
				{"ActivityIcon|IconActivity|ActivityIcon|ActivityIcon|RiPulseLine", "Status", false},
			}},
		),
	)
}

func faq(ds designSystem) h.Node {
	type qa struct{ q, a string }
	list := func(questions []qa) h.Node {
		var items []any
		for i, x := range questions {
			v := "item-" + string(rune('0'+i))
			items = append(items, accordion.Item(accordion.Value(v), accordion.Trigger(x.q), accordion.Content(x.a)))
		}
		// The aria source's Accordion has no default expanded item.
		if ds.Base == "aria" {
			return accordion.Accordion(accordion.TypeSingle, accordion.Collapsible(true), h.Group(items...))
		}
		return accordion.Accordion(accordion.TypeSingle, accordion.Collapsible(true), accordion.Value("item-0"), h.Group(items...))
	}
	return card.Card(
		card.Content(
			tabs.Tabs(tabs.Value("general"),
				tabs.List(h.Class("w-full"),
					tabs.Trigger(tabs.Value("general"), h.Class("flex-1"), "General"),
					tabs.Trigger(tabs.Value("billing"), h.Class("flex-1"), "Billing"),
					tabs.Trigger(tabs.Value("goals"), h.Class("flex-1"), "Goals"),
				),
				tabs.Content(tabs.Value("general"), list([]qa{
					{"How secure is my financial data with Ledger?", "We use bank-level AES-256 encryption, SOC 2 Type II certified infrastructure, and never store your credentials. All connections use read-only access tokens. We are a SEC registered investment advisor."},
					{"How do I connect my bank or investment accounts?", "Go to Settings > Linked Accounts and search for your institution. We support over 12,000 banks and brokerages via Plaid and MX."},
					{"Can I export my data for tax purposes?", "Yes. Navigate to Reports > Tax Export to download a CSV or PDF summary of your transactions, dividends, and capital gains for any tax year."},
				})),
				tabs.Content(tabs.Value("billing"), list([]qa{
					{"What is the difference between Basic and Pro pricing tiers?", "Basic includes budgeting, goal tracking, and up to 3 linked accounts. Pro adds unlimited accounts, dividend tracking, portfolio analysis, and priority support."},
					{"How do I cancel my subscription?", "Go to Settings > Billing > Manage Plan and click Cancel. Your access continues until the end of your current billing period."},
					{"Do you offer a free trial?", "Yes. All new accounts start with a 14-day Pro trial. No credit card required."},
				})),
				tabs.Content(tabs.Value("goals"), list([]qa{
					{"How do I set up a custom financial goal?", "Click New Goal from the Savings Targets card. Choose a category, set a target amount and date, and we'll calculate the monthly contribution needed."},
					{"Can I track multiple goals at once?", "Yes. Pro accounts can track unlimited goals. Basic accounts support up to 3 active goals."},
					{"How are monthly contributions calculated?", "We divide the remaining amount by the number of months until your target date, adjusted for your current savings rate and any auto-transfer schedules."},
				})),
			),
		),
		card.Footer(
			button.Button(button.VariantOutline, h.Class("w-full"), "Contact Support"),
			button.Button(button.VariantLink, h.Class("w-full"), "Learn More"),
		),
	)
}

func payments() h.Node {
	link := func(icon, title, description string) h.Node {
		return item.Item(lib.AsChild, item.VariantMuted,
			h.A(h.Href("#"),
				item.Media(item.MediaIcon, ip(icon, "")),
				item.Content(item.Title(title), item.Description(description)),
				ip(icChevronRight, "size-4 shrink-0 text-muted-foreground"),
			),
		)
	}
	return card.Card(
		card.Header(h.Class("flex flex-col gap-3"),
			breadcrumb.Breadcrumb(
				breadcrumb.List(
					breadcrumb.Item(breadcrumb.Link(h.Href("#"), "Home")),
					breadcrumb.Separator(),
					breadcrumb.Item(
						dropdownmenu.DropdownMenu(
							dropdownmenu.Trigger(lib.AsChild, button.Button(button.SizeIconSm, button.VariantGhost,
								ip("MoreHorizontalIcon|IconDots|MoreHorizontalCircle01Icon|DotsThreeIcon|RiMoreLine", ""),
								h.Span(h.Class("sr-only"), "Account options"),
							)),
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
				link("GaugeIcon|IconGauge|Settings01Icon|GaugeIcon|RiDashboardLine", "Change transfer limit", "Adjust how much you can send from your balance."),
				link("CalendarIcon|IconCalendar|Calendar03Icon|CalendarIcon|RiCalendarLine", "Scheduled transfers", "Set up a transfer to send at a later date."),
				link("RepeatIcon|IconRepeat|RepeatIcon|RepeatIcon|RiRepeatLine", "Direct Debits", "Set up and manage regular payments."),
				link("RefreshCwIcon|IconRefresh|RepeatIcon|ArrowsClockwiseIcon|RiRefreshLine", "Recurring card payments", "Manage your repeated card transactions."),
			),
		),
	)
}

func frontDoor() h.Node {
	return card.Card(
		card.Header(
			card.Title("Front Door"),
			card.Description("Smart Lock Pro"),
			card.Action(
				h.Div(h.Class("flex items-center gap-1.5 text-sm text-muted-foreground"),
					"Locked",
					ip("LockIcon|IconLock|SquareLock02Icon|LockKeyIcon|RiLockLine", "size-4"),
				),
			),
		),
		card.Content(
			h.Div(h.Class("relative flex aspect-video items-center justify-center overflow-hidden rounded-lg bg-muted bg-[repeating-linear-gradient(45deg,transparent,transparent_10px,var(--border)_10px,var(--border)_11px)]"),
				badge.Badge(badge.VariantDestructive, h.Class("absolute top-2 right-2"), "Live"),
			),
		),
	)
}

func releaseCatalog() h.Node {
	holdings := []struct{ ticker, name, typ, added, shares, value string }{
		{"VOO", "Vanguard S&P 500 ETF", "ETF", "Jan 2021", "112", "$48,230.40"},
		{"VIG", "Vanguard Dividend Appreciation", "ETF", "Mar 2022", "450", "$26,033.79"},
		{"AAPL", "Apple Inc.", "Stock", "Nov 2020", "85", "$18,488.90"},
		{"O", "Realty Income Corp", "REIT", "Jun 2023", "320", "$15,136.59"},
	}
	var items []any
	for _, x := range holdings {
		items = append(items, item.Item(item.VariantMuted,
			item.Media(h.Div(h.Class("flex size-12 items-center justify-center rounded-lg border text-sm font-semibold"), x.ticker)),
			item.Content(
				item.Title(x.name),
				item.Description(h.Class("text-xs tracking-wider uppercase"), x.shares+" Shares · "+x.added),
			),
			h.Div(h.Class("flex shrink-0 items-center gap-6"),
				badge.Badge(badge.VariantOutline, x.typ),
				h.Div(h.Class("flex flex-col items-end gap-0.5"),
					h.Span(h.Class("text-xs tracking-wider text-muted-foreground uppercase"), "Value"),
					h.Span(h.Class("font-medium tabular-nums"), x.value),
				),
			),
		))
	}
	return card.Card(
		card.Header(
			h.Div(h.Class("flex items-center justify-between gap-3"),
				inputgroup.InputGroup(h.Class("max-w-sm"),
					inputgroup.Addon(ip(icSearch, "")),
					inputgroup.Input(h.Placeholder("Search holdings or tickers...")),
				),
				togglegroup.ToggleGroup(togglegroup.TypeSingle, togglegroup.Value("etfs"), togglegroup.VariantOutline, togglegroup.Spacing(1),
					togglegroup.Item(togglegroup.Value("stocks"), "Stocks"),
					togglegroup.Item(togglegroup.Value("etfs"), "ETFs"),
					togglegroup.Item(togglegroup.Value("reits"), "REITs"),
				),
			),
		),
		card.Content(item.Group(items...)),
	)
}

// --- column 4 ---

func accountAccess() h.Node {
	return card.Card(
		card.Header(card.Title("Account Access"), card.Description("Update your credentials or re-authenticate.")),
		card.Content(
			field.Group(
				field.Field(
					field.Label(h.For("email-address"), "Email Address"),
					input.Input(h.ID("email-address"), h.Type("email"), h.Value("artist@studio.inc")),
				),
				field.Field(
					h.Div(h.Class("flex items-center justify-between"),
						field.Label(h.For("current-password"), "Current Password"),
						h.A(h.Href("#"), h.Class("text-xs font-medium tracking-wider text-muted-foreground uppercase hover:text-foreground"), "Forgot?"),
					),
					input.Input(h.ID("current-password"), h.Type("password"), h.Value("password123")),
				),
			),
		),
		card.Footer(h.Class("flex-col gap-4"),
			button.Button(h.Class("w-full"),
				ip("LockKeyholeIcon|IconLock|SquareLock02Icon|LockKeyIcon|RiLockLine", ""),
				"Update Security",
			),
			item.Item(lib.AsChild, item.VariantMuted,
				h.A(h.Href("#"),
					item.Media(item.MediaIcon, ip("AlertCircleIcon|IconAlertCircle|AlertCircleIcon|WarningCircleIcon|RiErrorWarningLine", "text-destructive")),
					item.Content(item.Title("Danger Zone"), item.Description(h.Class("line-clamp-1"), "Archive account and remove catalog")),
					ip("ArrowRightIcon|IconArrowRight|ArrowRight01Icon|ArrowRightIcon|RiArrowRightLine", "size-4"),
				),
			),
		),
	)
}

func cardOverview() h.Node {
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	values := []float64{40, 55, 35, 60, 45, 50, 65, 40, 55, 70, 45, 80}
	var rows []chart.Row
	for i, m := range months {
		rows = append(rows, chart.Row{"month": m, "amount": values[i]})
	}
	return h.Div(h.Class("grid grid-cols-2 gap-3"),
		card.Card(
			card.Content(
				card.Description("Card Balance"),
				card.Title(h.Class("text-2xl tabular-nums"), "US$12.94"),
				card.Description(h.Class("tabular-nums"), "US$11,337.06 Available"),
			),
		),
		card.Card(h.Class("flex flex-col justify-between"),
			card.Content(h.Class("flex flex-1 flex-col justify-between"),
				h.Div(h.Class("flex flex-col gap-1"),
					card.Description("Payment Due"),
					card.Title(h.Class("text-2xl"), "1 Apr"),
				),
				button.Button(button.VariantOutline, button.SizeSm, h.Class("mt-3 w-full"), "Pay Early"),
			),
		),
		card.Card(h.Class("col-span-2"),
			card.Content(h.Class("flex flex-col gap-2"),
				h.Div(h.Class("flex items-center justify-between"),
					card.Description("Yearly Activity"),
					badge.Badge(badge.VariantSecondary, "+US$0.25 Daily Cash"),
				),
				barChart(chart.Config{{Key: "amount", Label: "Activity", Color: "var(--chart-2)"}}, "h-20 w-full", rows, zeroMargin(4),
					chart.XAxis(chart.DataKey("month"), chart.TickLine(false), chart.TickMargin(4), chart.AxisLine(false),
						chart.TickFormatter(chart.Fn{Name: "slice", Args: []any{0, 1}}), chart.ClassName("text-[10px]")),
					noCursorTooltip(),
					chart.Bar(chart.DataKey("amount"), chart.Fill("var(--color-amount)"), chart.Radii(3, 3, 0, 0)),
				),
			),
		),
	)
}

func accountSelect(id, value string, options ...[2]string) h.Node {
	var items []any
	for _, o := range options {
		items = append(items, selectui.Item(selectui.ItemValue(o[0]), o[1]))
	}
	return selectui.Select(selectui.DefaultValue(value),
		selectui.Trigger(h.ID(id), h.Class("w-full"), selectui.Value()),
		selectui.Content(selectui.Group(items...)),
	)
}

func transferFunds() h.Node {
	return card.Card(
		card.Header(
			card.Title("Transfer Funds"),
			card.Description("Move money between your connected accounts."),
			closeAction(),
		),
		card.Content(
			field.Group(
				field.Field(
					field.Label(h.For("transfer-amount"), "Amount to Transfer"),
					inputgroup.InputGroup(
						inputgroup.Addon(inputgroup.Text("$")),
						inputgroup.Input(h.ID("transfer-amount"), h.Value("1,200.00")),
					),
				),
				field.Field(
					field.Label(h.For("from-account"), "From Account"),
					accountSelect("from-account", "checking",
						[2]string{"checking", "Main Checking (··8402) — $12,450.00"},
						[2]string{"business", "Business (··7731) — $8,920.00"}),
				),
				field.Field(
					field.Label(h.For("to-account"), "To Account"),
					accountSelect("to-account", "savings",
						[2]string{"savings", "High Yield Savings (··1192) — $42,100.00"},
						[2]string{"investment", "Investment (··3349) — $18,200.00"}),
				),
				item.Item(item.VariantMuted, h.Class("flex-col items-stretch"),
					item.Content(h.Class("gap-3"),
						kv("Estimated arrival", "Today, Apr 14", "text-sm font-medium"),
						separator.Separator(),
						kv("Transaction fee", "$0.00", "text-sm font-medium tabular-nums"),
						separator.Separator(),
						kv("Total amount", "$1,200.00", "text-sm font-semibold tabular-nums", "text-sm font-medium"),
					),
				),
			),
		),
		card.Footer(button.Button(h.Class("w-full"), "Confirm Transfer")),
	)
}

func coverArt() h.Node {
	return card.Card(
		card.Content(h.Class("flex flex-col gap-3"),
			label.Label(h.For("cover-art"), h.Class("text-center text-xs font-normal tracking-wider text-muted-foreground uppercase"), "Cover Art"),
			item.Item(h.Class("aspect-square"), item.VariantOutline,
				h.Label(h.For("cover-art"), h.Class("flex size-full cursor-pointer items-center justify-center"),
					ip("ImageIcon|IconPhoto|Image01Icon|ImageIcon|RiImageLine", "size-10 text-muted-foreground/50"),
				),
			),
			h.Input(h.ID("cover-art"), h.Type("file"), h.Attr("accept", "image/jpeg,image/png"), h.Class("sr-only")),
		),
		card.Footer(h.Class("flex-col gap-2"),
			button.Button(button.VariantSecondary, h.Class("w-full"), lib.AsChild,
				h.Label(h.For("cover-art"), h.Class("cursor-pointer"), "Upload Artwork"),
			),
			card.Description(h.Class("text-center text-xs"), "Minimum 3000 × 3000px", h.Br(), "JPEG or PNG only"),
		),
	)
}

func loadingCard() h.Node {
	return card.Card(
		card.Header(
			skeleton.Skeleton(h.Class("h-5 w-32")),
			skeleton.Skeleton(h.Class("h-4 w-48")),
		),
		card.Content(h.Class("flex flex-col gap-4"),
			skeleton.Skeleton(h.Class("h-32 w-full rounded-lg")),
			h.Div(h.Class("flex flex-col gap-2"),
				skeleton.Skeleton(h.Class("h-4 w-full")),
				skeleton.Skeleton(h.Class("h-4 w-3/4")),
				skeleton.Skeleton(h.Class("h-4 w-1/2")),
			),
			h.Div(h.Class("flex gap-2"),
				skeleton.Skeleton(h.Class("h-9 flex-1 rounded-md")),
				skeleton.Skeleton(h.Class("h-9 flex-1 rounded-md")),
			),
		),
	)
}

// --- column 5 ---

func receivingMethod() h.Node {
	method := func(id, value, title, description, descClass string) h.Node {
		return field.Label(h.For(id),
			field.Field(field.Horizontal, h.Class("pb-2.5"),
				radiogroup.Item(radiogroup.Value(value), h.ID(id)),
				field.Content(field.Title(title), field.Description(h.Class(descClass), description)),
			),
		)
	}
	return card.Card(
		card.Header(
			card.Description("Payout Preferences"),
			card.Title("Receiving Method"),
			closeAction(),
		),
		card.Content(
			field.Group(
				field.Field(
					field.Label(h.For("account-holder"), "Account Holder Name"),
					input.Input(h.ID("account-holder"), h.Value("Synthetic Horizons Music LLC")),
				),
				field.Set(
					field.Legend(field.LegendLabel, "Receiving Method"),
					radiogroup.RadioGroup(radiogroup.Value("bank"), h.Class("grid grid-cols-1 items-start gap-3 md:grid-cols-2 style-sera:grid-cols-1"),
						method("method-bank", "bank", "Bank Transfer", "SWIFT / IBAN", ""),
						method("method-paypal", "paypal", "PayPal", "Instant Payout", "line-clamp-1"),
					),
				),
				field.Field(
					field.Label(h.For("iban"), "IBAN / Account Number"),
					input.Input(h.ID("iban"), h.Placeholder("DE89 3704 0044 ....")),
				),
			),
		),
		card.Footer(button.Button(h.Class("w-full"), h.Bool("disabled"), "Save Payout Settings")),
	)
}

func powerUsage(ds designSystem) h.Node {
	hours := []string{"6a", "8a", "10a", "12p", "2p", "4p", "6p", "8p"}
	usage := []float64{1.2, 2.8, 3.1, 2.4, 3.4, 2.9, 3.8, 3.2}
	var rows []chart.Row
	for i, x := range hours {
		rows = append(rows, chart.Row{"hour": x, "usage": usage[i]})
	}
	stat := func(label, value, class string) h.Node {
		return h.Div(h.Class("flex flex-col gap-0.5"),
			h.Span(h.Class("text-sm text-muted-foreground"), label),
			h.Span(h.Class(class), value),
		)
	}
	return card.Card(
		card.Header(card.Title("Power Usage"), card.Description("Whole Home")),
		card.Content(h.Class("flex flex-col gap-4"),
			barChart(chart.Config{{Key: "usage", Label: "Usage (kW)", Color: "var(--chart-2)"}}, "h-[140px] w-full", rows,
				map[string]any{"left": 0, "right": 0, "top": 4, "bottom": 0},
				chart.XAxis(chart.DataKey("hour"), chart.TickLine(false), chart.TickMargin(6), chart.AxisLine(false), chart.ClassName("text-xs")),
				noCursorTooltip(),
				chart.Bar(chart.DataKey("usage"), chart.Fill("var(--color-usage)"), barRadius(rounded(ds), 4)),
			),
			separator.Separator(),
			h.Div(h.Class("grid grid-cols-2 gap-4"),
				stat("Currently Using", "3.4 kW", "text-lg font-semibold tabular-nums"),
				stat("Solar Gen", "+1.2 kW", "text-lg font-semibold text-chart-1 tabular-nums"),
			),
		),
		card.Footer(h.Class("flex-col items-start gap-1"),
			h.Span(h.Class("text-sm text-muted-foreground"), "Battery Level"),
			h.Div(h.Class("flex w-full items-center gap-2"),
				progress.Progress(progress.Value(85), h.Class("flex-1")),
				h.Span(h.Class("text-sm font-medium tabular-nums"), "85%"),
			),
		),
	)
}

func emptyConnectBank() h.Node {
	return emptyCard(icCreditCard, "Connect Bank", "Link your payout method to receive monthly royalty distributions automatically.", "Set Up Payouts")
}

func upcomingPayments() h.Node {
	pay := func(title, date, amount string) h.Node {
		return item.Item(item.VariantMuted,
			item.Content(item.Title(title), item.Description(date)),
			badge.Badge(badge.VariantSecondary, amount),
		)
	}
	return card.Card(
		card.Header(card.Title("Upcoming Payments"), card.Description("Select a date to view scheduled payments.")),
		card.Content(h.Class("flex flex-col gap-4"),
			item.Item(item.VariantOutline, h.Class("justify-center"),
				calendar.Calendar(calendar.ModeSingle, calendar.Selected(lib.Now()),
					h.Class("w-full [--cell-size:--spacing(8)] md:[--cell-size:--spacing(10)] style-sera:md:[--cell-size:--spacing(9)]")),
			),
			item.Group(h.Class("w-full"),
				pay("Netflix Subscription", "Apr 15, 2024", "$19.99"),
				pay("Rent Payment", "Apr 1, 2024", "$2,400.00"),
				pay("Auto Insurance", "Apr 22, 2024", "$186.00"),
			),
		),
	)
}

func rollerShades() h.Node {
	return card.Card(h.Data("create-roller-shades", ""),
		card.Header(card.Title("Living Room"), card.Description("Roller Shades")),
		card.Content(h.Class("flex flex-col gap-4"),
			h.Div(h.Class("flex h-32 flex-col overflow-hidden rounded-lg border bg-muted"),
				h.Div(h.Class("bg-muted-foreground transition-all duration-300"), h.Style("height: 50%")),
			),
			h.Div(h.Class("flex items-center gap-3"),
				h.Span(h.Class("text-xs font-medium tracking-wider text-muted-foreground uppercase"), "Open"),
				slider.Slider(slider.DefaultValue(50), slider.Max(100), h.Class("flex-1")),
				h.Span(h.Class("text-xs font-medium tracking-wider text-muted-foreground uppercase"), "Close"),
			),
		),
		card.Footer(
			togglegroup.ToggleGroup(togglegroup.TypeSingle, togglegroup.Value("half"), togglegroup.VariantOutline, togglegroup.Spacing(1), h.Class("w-full"),
				togglegroup.Item(togglegroup.Value("open"), h.Class("flex-1"), "Open"),
				togglegroup.Item(togglegroup.Value("half"), h.Class("flex-1"), "Half"),
				togglegroup.Item(togglegroup.Value("closed"), h.Class("flex-1"), "Closed"),
			),
		),
	)
}

// --- column 6 ---

func stockPerformance() h.Node {
	var items []any
	for _, t := range []string{"VOO", "VIG", "AAPL", "MSFT", "GOOGL", "AMZN", "TSLA"} {
		items = append(items, combobox.Item(combobox.ItemValue(t), t))
	}
	rows := []chart.Row{
		{"month": "Jan", "price": 412}, {"month": "Feb", "price": 438}, {"month": "Mar", "price": 395},
		{"month": "Apr", "price": 450}, {"month": "May", "price": 420}, {"month": "Jun", "price": 462},
	}
	return card.Card(h.Data("create-stock-performance", ""),
		card.Header(card.Title("Stock Performance"), card.Description("6-month price history.")),
		card.Content(h.Class("flex flex-col gap-4"),
			field.Group(
				field.Field(
					field.Label(h.For("ticker-select"), "Ticker"),
					combobox.Combobox(combobox.DefaultValue("VOO"),
						combobox.Input(h.ID("ticker-select"), h.Placeholder("Search ticker...")),
						combobox.Content(
							combobox.Empty("No tickers found."),
							combobox.List(items...),
						),
					),
				),
			),
			separator.Separator(h.Class("style-sera:hidden")),
			chart.Container(chart.Config{{Key: "price", Label: "Price", Color: "var(--chart-1)"}}, h.Class("h-[200px] w-full"),
				chart.AreaChart(chart.AccessibilityLayer(), chart.Data(rows), chart.P("margin", map[string]any{"left": 0, "right": 0, "top": 8, "bottom": 0}),
					chart.Defs(
						chart.El("linearGradient", chart.P("id", "fillPrice"), chart.P("x1", "0"), chart.P("y1", "0"), chart.P("x2", "0"), chart.P("y2", "1"),
							chart.El("stop", chart.P("offset", "0%"), chart.P("stopColor", "var(--color-price)"), chart.P("stopOpacity", 0.3)),
							chart.El("stop", chart.P("offset", "100%"), chart.P("stopColor", "var(--color-price)"), chart.P("stopOpacity", 0.05)),
						),
					),
					chart.CartesianGrid(chart.Vertical(false), chart.P("strokeDasharray", "3 3")),
					noCursorTooltip(),
					chart.Area(chart.Type("monotone"), chart.DataKey("price"), chart.Stroke("var(--color-price)"), chart.StrokeWidth(2), chart.Fill("url(#fillPrice)")),
				),
			),
		),
	)
}

func emptyExploreCatalog() h.Node {
	return emptyCard("AudioLinesIcon|IconPlayerRecordFilled|AudioWave01Icon|RecordIcon|RiRecordCircleLine", "Explore Catalog", "Check your ISRC codes, metadata, and visual assets before going live.", "View Catalog")
}

func newMilestone() h.Node {
	return card.Card(
		card.Header(card.Title("Set a new milestone"), card.Description("Define your financial target and we'll help you pace your savings.")),
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

func socialLinks() h.Node {
	link := func(id, label, icon string, value, placeholder string) h.Node {
		in := inputgroup.Input(h.ID(id))
		if value != "" {
			in.Children(h.Value(value))
		} else {
			in.Children(h.Placeholder(placeholder))
		}
		return field.Field(
			field.Label(h.For(id), label),
			inputgroup.InputGroup(inputgroup.Addon(ip(icon, "")), in),
		)
	}
	return card.Card(
		card.Header(card.Title("Social Links")),
		card.Content(
			field.Group(
				link("spotify-url", "Spotify Artist URL", "CirclePlusIcon|IconCirclePlus|PlusSignCircleIcon|PlusCircleIcon|RiAddCircleLine", "spotify.com/artist/3j...2k", ""),
				link("instagram-handle", "Instagram Handle", "CameraIcon|IconCamera|Camera01Icon|CameraIcon|RiCameraLine", "@julianduryea_music", ""),
				link("soundcloud-url", "SoundCloud URL", "CloudIcon|IconCloud|CloudUploadIcon|CloudIcon|RiCloudLine", "", "soundcloud.com/username"),
				link("website-url", "Website", "GlobeIcon|IconWorld|Globe02Icon|GlobeIcon|RiGlobalLine", "", "https://yoursite.com"),
			),
		),
		card.Footer(h.Class("justify-end gap-2 style-sera:justify-center"),
			button.Button(button.VariantSecondary, h.Class("style-sera:flex-1"), "Discard"),
			button.Button(h.Class("style-sera:flex-1"), "Save Changes"),
		),
	)
}

func notificationSettings() h.Node {
	notifications := []struct {
		id, label, description string
		checked                bool
	}{
		{"transactions", "Transaction alerts", "Deposits, withdrawals, and transfers.", true},
		{"security", "Security alerts", "Login attempts and account changes.", true},
		{"goals", "Goal milestones", "Updates at 25%, 50%, 75%, and 100%.", false},
		{"market", "Market updates", "Daily portfolio summary and price alerts.", false},
	}
	fields := []any{
		field.Field(field.Horizontal,
			checkbox.Checkbox(h.ID("notify-all"), checkbox.Indeterminate, h.Data("create-notify-all", "")),
			field.Content(field.Label(h.For("notify-all"), "Select all")),
		),
	}
	for _, n := range notifications {
		fields = append(fields, field.Field(field.Horizontal,
			checkbox.Checkbox(h.ID("notify-"+n.id), checkbox.Checked(n.checked), h.Data("create-notify", n.id)),
			field.Content(field.Label(h.For("notify-"+n.id), n.label), field.Description(n.description)),
		))
	}
	return card.Card(h.Data("create-notifications", ""),
		card.Header(card.Title("Notifications"), card.Description("Choose what you want to be notified about.")),
		card.Content(field.Group(fields...)),
		card.Footer(button.Button(h.Class("w-full"), "Save Preferences")),
	)
}
