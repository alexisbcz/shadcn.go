package site

import (
	"fmt"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/alert"
	alertdialog "shadcn.go/components/ui/alert-dialog"
	"shadcn.go/components/ui/avatar"
	"shadcn.go/components/ui/badge"
	"shadcn.go/components/ui/button"
	buttongroup "shadcn.go/components/ui/button-group"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/chart"
	"shadcn.go/components/ui/checkbox"
	"shadcn.go/components/ui/dialog"
	dropdownmenu "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/components/ui/empty"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/input"
	inputgroup "shadcn.go/components/ui/input-group"
	"shadcn.go/components/ui/item"
	"shadcn.go/components/ui/kbd"
	nativeselect "shadcn.go/components/ui/native-select"
	"shadcn.go/components/ui/progress"
	radiogroup "shadcn.go/components/ui/radio-group"
	selectui "shadcn.go/components/ui/select"
	"shadcn.go/components/ui/separator"
	"shadcn.go/components/ui/skeleton"
	"shadcn.go/components/ui/slider"
	switchui "shadcn.go/components/ui/switch"
	"shadcn.go/components/ui/table"
	"shadcn.go/components/ui/tabs"
	"shadcn.go/components/ui/textarea"
	togglegroup "shadcn.go/components/ui/toggle-group"
	"shadcn.go/components/ui/tooltip"
	h "shadcn.go/libhtml"
)

// The preview block of the create designer
// (registry/bases/<base>/blocks/preview): one card per file upstream.

func preview01(ds designSystem) h.Node {
	col := func(children ...any) h.Node { return previewColumn(false, children...) }
	if ds.Base != "radix" {
		col = previewColumnPx
	}
	return previewGrid(
		col(styleOverview(ds), typographySpecimen(ds), h.Div(h.Class("md:hidden"), uiElements()), codespacesCard(), invoiceCard()),
		col(iconPreviewGrid(), h.Div(h.Class("hidden w-full md:flex"), uiElements()), observabilityCard(), shippingAddress()),
		col(environmentVariables(), barChartCard(ds), inviteTeam(), activateAgentDialog()),
		col(skeletonLoading(), pieChartCard(), noTeamMembers(), reportBug(), contributors()),
		col(feedbackForm(), bookAppointment(), sleepReport(ds), githubProfile()),
		col(weeklyFitnessSummary(ds), fileUpload(), previewAnalyticsCard(), usageCard(), shortcuts()),
		col(anomalyAlert(), liveWaveformCard(), visitors(), contributionsActivity(), notFoundCard()),
	)
}

// fontLabels are the style overview's and typography specimen's names.
func (ds designSystem) fontNames() (body, heading string) {
	d := createDataOnce()
	if f, ok := d.font(ds.Font); ok {
		body = f.Title
	}
	if f, ok := d.font(ds.FontHeading); ok {
		heading = f.Title
	}
	return body, heading
}

func (ds designSystem) styleTitle() string {
	for _, s := range createDataOnce().Styles {
		if s.Name == ds.Style {
			return s.Title
		}
	}
	return ""
}

// --- column 1 ---

func styleOverview(ds designSystem) h.Node {
	body, heading := ds.fontNames()
	name := body
	if heading != "" && heading != body {
		name = heading
	}
	var swatches []any
	for _, v := range []string{"--background", "--foreground", "--primary", "--secondary", "--muted", "--accent", "--border", "--chart-1", "--chart-2", "--chart-3", "--chart-4", "--chart-5"} {
		swatches = append(swatches, h.Div(h.Class("flex flex-col flex-wrap items-center gap-2"),
			h.Div(h.Class("relative aspect-square w-full rounded-lg bg-(--color) after:absolute after:inset-0 after:rounded-lg after:border after:border-border after:mix-blend-darken dark:after:mix-blend-lighten style-sera:rounded-none style-sera:after:rounded-none"),
				h.Style("--color: var("+v+")")),
			h.Div(h.Class("hidden max-w-14 truncate font-mono text-[0.60rem] md:block style-lyra:max-w-10 style-mira:max-w-10"), v),
		))
	}
	return card.Card(
		card.Content(h.Class("flex flex-col gap-6 style-lyra:gap-4 style-mira:gap-4"),
			h.Div(h.Class("flex flex-col gap-1"),
				h.Div(h.Class("cn-font-heading text-2xl font-medium style-lyra:text-lg style-mira:text-lg style-sera:text-lg style-sera:font-semibold style-sera:tracking-wide style-sera:uppercase"),
					h.Data("create-style-overview", ""), ds.styleTitle()+" - "+name),
				h.Div(h.Class("line-clamp-2 text-base text-muted-foreground style-lyra:text-sm style-mira:text-sm style-sera:text-sm style-sera:leading-relaxed"),
					"Designers love packing quirky glyphs into test phrases. This is a preview of the typography styles."),
			),
			h.Div(h.Class("grid grid-cols-6 gap-3"), h.Group(swatches...)),
		),
	)
}

func typographySpecimen(ds designSystem) h.Node {
	body, heading := ds.fontNames()
	headingLabel := "Inherit"
	if heading != "" && heading != body {
		headingLabel = heading
	}
	if body == "" {
		body = "Default"
	}
	return card.Card(
		card.Content(h.Class("flex flex-col gap-2"),
			h.Div(h.Class("text-xs font-medium text-muted-foreground uppercase"), h.Data("create-typography-label", ""), headingLabel+" - "+body),
			h.P(h.Class("cn-font-heading text-2xl font-medium style-sera:text-lg style-sera:font-semibold style-sera:tracking-wider style-sera:uppercase"), "Designing with rhythm and hierarchy."),
			h.P(h.Class("text-sm leading-relaxed text-muted-foreground"), "A strong body style keeps long-form content readable and balances the visual weight of headings."),
			h.P(h.Class("text-sm leading-relaxed text-muted-foreground"), "Thoughtful spacing and cadence help paragraphs scan quickly without feeling dense."),
		),
		card.Footer(
			dialog.Dialog(
				dialog.Trigger(lib.AsChild, button.Button(button.VariantOutline, h.Class("w-full"), "Share Feedback")),
				dialog.Content(
					dialog.Header(
						dialog.Title("Share Feedback"),
						dialog.Description("Let us know how we can improve your experience."),
					),
					field.Group(
						h.Div(h.Class("grid grid-cols-2 gap-3"),
							field.Field(field.Label(h.For("feedback-name"), "Name"), input.Input(h.ID("feedback-name"), h.Placeholder("Your name"))),
							field.Field(field.Label(h.For("feedback-email"), "Email"), input.Input(h.ID("feedback-email"), h.Type("email"), h.Placeholder("you@example.com"))),
						),
						field.Field(field.Label(h.For("feedback-category"), "Category"),
							accountSelect("feedback-category", "general",
								[2]string{"general", "General"}, [2]string{"bug", "Bug Report"},
								[2]string{"feature", "Feature Request"}, [2]string{"improvement", "Improvement"}),
						),
						field.Field(field.Label(h.For("feedback-message"), "Message"),
							textarea.Textarea(h.ID("feedback-message"), h.Placeholder("Tell us what's on your mind..."), h.Class("min-h-24 resize-none"))),
					),
					dialog.Footer(
						dialog.Close(lib.AsChild, button.Button(button.VariantOutline, "Cancel")),
						button.Button("Submit"),
					),
				),
			),
		),
	)
}

func uiElements() h.Node {
	return card.Card(h.Class("w-full"),
		card.Content(h.Class("flex flex-col gap-6"),
			h.Div(h.Class("flex flex-col gap-4"),
				h.Div(h.Class("flex flex-wrap gap-2"),
					button.Button("Button"),
					button.Button(button.VariantSecondary, "Secondary"),
					button.Button(button.VariantOutline, "Outline"),
					button.Button(button.VariantGhost, "Ghost"),
				),
				item.Item(item.VariantOutline,
					item.Content(
						item.Title("Two-factor authentication"),
						item.Description(h.Class("text-pretty xl:hidden 2xl:block"), "Verify via email or phone number."),
					),
					item.Actions(h.Class("hidden md:flex"), button.Button(button.SizeSm, button.VariantSecondary, "Enable")),
				),
			),
			slider.Slider(slider.DefaultValue(500), slider.Max(1000), slider.Min(0), slider.Step(10), h.Class("flex-1"), h.Aria("label", "Slider")),
			field.Group(
				field.Field(
					inputgroup.InputGroup(
						inputgroup.Input(h.Placeholder("Name")),
						inputgroup.Addon(inputgroup.AlignInlineEnd, inputgroup.Text(ip(icSearch, ""))),
					),
				),
				field.Field(h.Class("flex-1"), textarea.Textarea(h.Placeholder("Message"), h.Class("resize-none"))),
			),
			h.Div(h.Class("flex items-center gap-2"),
				h.Div(h.Class("flex gap-2"),
					badge.Badge("Badge"),
					badge.Badge(badge.VariantSecondary, "Secondary"),
					badge.Badge(badge.VariantOutline, "Outline"),
				),
				radiogroup.RadioGroup(radiogroup.Value("apple"), h.Class("ml-auto flex w-fit gap-3"),
					radiogroup.Item(radiogroup.Value("apple")),
					radiogroup.Item(radiogroup.Value("banana")),
				),
				h.Div(h.Class("flex gap-3"),
					checkbox.Checkbox(checkbox.Checked(true)),
					checkbox.Checkbox(),
				),
			),
			h.Div(h.Class("flex items-center gap-4"),
				alertdialog.AlertDialog(
					alertdialog.Trigger(lib.AsChild, button.Button(button.VariantOutline,
						h.Span(h.Class("hidden md:flex style-sera:md:hidden"), "Alert Dialog"),
						h.Span(h.Class("flex md:hidden style-sera:md:flex"), "Dialog"),
					)),
					alertdialog.Content(alertdialog.SizeSm,
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
				buttongroup.ButtonGroup(
					button.Button(button.VariantOutline,
						h.Span(h.Class("style-sera:hidden"), "Button Group"),
						h.Span(h.Class("hidden style-sera:block"), "Group"),
					),
					dropdownmenu.DropdownMenu(
						dropdownmenu.Trigger(lib.AsChild, button.Button(button.VariantOutline, button.SizeIcon,
							ip("ChevronUpIcon|IconChevronUp|ArrowUp01Icon|CaretUpIcon|RiArrowUpSLine", ""))),
						dropdownmenu.Content(dropdownmenu.AlignEnd, dropdownmenu.SideTop, h.Class("w-fit"),
							dropdownmenu.Group(
								dropdownmenu.Label("Quick Actions"),
								dropdownmenu.Item("Mute Conversation"),
								dropdownmenu.Item("Mark as Read"),
								dropdownmenu.Item("Block User"),
							),
							dropdownmenu.Separator(),
							dropdownmenu.Group(
								dropdownmenu.Label("Conversation"),
								dropdownmenu.Item("Share Conversation"),
								dropdownmenu.Item("Copy Conversation"),
								dropdownmenu.Item("Report Conversation"),
							),
							dropdownmenu.Separator(),
							dropdownmenu.Group(
								dropdownmenu.Item(dropdownmenu.VariantDestructive, "Delete Conversation"),
							),
						),
					),
				),
				switchui.Switch(switchui.Checked(true), h.Class("ml-auto")),
			),
		),
	)
}

const (
	icPlusSign = "PlusIcon|IconPlus|PlusSignIcon|PlusIcon|RiAddLine"
	icServer   = "ServerIcon|IconServer|ServerStackIcon|HardDrivesIcon|RiHardDriveLine"
	icInfo     = "InfoIcon|IconInfoCircle|AlertCircleIcon|InfoIcon|RiInformationLine"
	icCopy     = "CopyIcon|IconCopy|Copy01Icon|CopyIcon|RiFileCopyLine"
)

func codespacesCard() h.Node {
	cloneField := func(id, label, value string, description ...any) h.Node {
		return field.Field(h.Class("gap-2"),
			field.Label(h.For(id), h.Class("sr-only"), label),
			inputgroup.InputGroup(
				inputgroup.Addon(inputgroup.AlignInlineEnd,
					inputgroup.Button(button.VariantGhost, inputgroup.ButtonSizeIconXs, ip(icCopy, "")),
				),
				inputgroup.Input(h.ID(id), h.Value(value), h.Bool("readonly")),
			),
			field.Description(description...),
		)
	}
	return card.Card(
		card.Content(
			tabs.Tabs(tabs.Value("codespaces"),
				tabs.List(h.Class("w-full"),
					tabs.Trigger(tabs.Value("codespaces"), "Codespaces"),
					tabs.Trigger(tabs.Value("local"), "Local"),
				),
				tabs.Content(tabs.Value("codespaces"),
					item.Item(item.SizeSm, h.Class("px-1 pt-2"),
						item.Content(item.Title("Codespaces"), item.Description("Your workspaces in the cloud")),
						item.Actions(
							tooltip.Tooltip(
								tooltip.Trigger(lib.AsChild, button.Button(button.VariantGhost, button.SizeIconSm, ip(icPlusSign, ""))),
								tooltip.Content(tooltip.SideBottom, "Create a codespace on main"),
							),
							dropdownmenu.DropdownMenu(
								dropdownmenu.Trigger(lib.AsChild, button.Button(button.VariantGhost, button.SizeIconSm,
									ip("MoreHorizontalIcon|IconDots|MoreHorizontalCircle01Icon|DotsThreeOutlineIcon|RiMoreLine", ""))),
								dropdownmenu.Content(dropdownmenu.AlignEnd, h.Class("w-56"),
									dropdownmenu.Group(
										dropdownmenu.Item(ip(icPlusSign, ""), "New with options..."),
										dropdownmenu.Item(ip("ContainerIcon|IconBox|CubeIcon|CubeIcon|RiBox1Line", ""), "Configure container"),
										dropdownmenu.Item(ip("ZapIcon|IconBolt|ZapIcon|LightningIcon|RiFlashlightLine", ""), "Set up prebuilds"),
									),
									dropdownmenu.Separator(),
									dropdownmenu.Group(
										dropdownmenu.Item(ip(icServer, ""), "Manage codespaces"),
										dropdownmenu.Item(ip("ShareIcon|IconShare2|Share03Icon|ShareIcon|RiShareLine", ""), "Share deep link"),
										dropdownmenu.Item(ip(icInfo, ""), "What are codespaces?"),
									),
								),
							),
						),
					),
					separator.Separator(h.Class("-mx-2 my-2 w-auto!")),
					empty.Empty(h.Class("p-4"),
						empty.Header(
							empty.Media(empty.MediaIcon, ip(icServer, "")),
							empty.Title("No codespaces"),
							empty.Description("You don't have any codespaces with this repository checked out"),
						),
						empty.Content(
							button.Button(button.SizeSm, h.Data("create-codespace", ""), "Create Codespace"),
							h.A(h.Href("#learn-more"), h.Class("text-xs text-muted-foreground underline underline-offset-4"), "Learn more about codespaces"),
						),
					),
					separator.Separator(h.Class("-mx-2 my-2 w-auto!")),
					h.Div(h.Class("p-1.5 text-xs text-muted-foreground"),
						"Codespace usage for this repository is paid for by ", h.Span(h.Class("font-medium"), "shadcn"), "."),
				),
				tabs.Content(tabs.Value("local"),
					item.Item(item.SizeSm, h.Class("hidden p-0"),
						item.Content(item.Title(h.Class("gap-2"),
							ip("TerminalIcon|IconTerminal|ComputerTerminal01Icon|TerminalIcon|RiTerminalBoxLine", "size-4"), "Clone")),
						item.Actions(
							tooltip.Tooltip(
								tooltip.Trigger(lib.AsChild, button.Button(button.VariantGhost, button.SizeIcon, ip(icInfo, ""))),
								tooltip.Content(tooltip.SideLeft, "Which remote URL should I use?"),
							),
						),
					),
					tabs.Tabs(tabs.Value("https"),
						tabs.List(tabs.ListVariantLine, h.Class("w-full justify-start border-b *:[button]:flex-0"),
							tabs.Trigger(tabs.Value("https"), "HTTPS"),
							tabs.Trigger(tabs.Value("ssh"), "SSH"),
							tabs.Trigger(tabs.Value("cli"), "GitHub CLI"),
						),
						h.Div(h.Class("rounded-md border bg-muted/30 p-2"),
							tabs.Content(tabs.Value("https"), cloneField("https-url", "HTTPS URL", "https://github.com/shadcn-ui/ui.git", "Clone using the web URL.")),
							tabs.Content(tabs.Value("ssh"), cloneField("ssh-url", "SSH URL", "git@github.com:shadcn-ui/ui.git", "Use a password-protected SSH key.")),
							tabs.Content(tabs.Value("cli"), cloneField("cli-command", "CLI Command", "gh repo clone shadcn-ui/ui", "Work fast with our official CLI. ", h.A(h.Href("#learn-more"), "Learn more"))),
						),
					),
					separator.Separator(h.Class("-mx-2 my-2 w-auto!")),
					h.Div(h.Class("flex flex-col"),
						button.Button(button.VariantGhost, button.SizeSm, h.Class("justify-start gap-1.5"),
							ip("MonitorIcon|IconDeviceDesktop|ComputerIcon|MonitorIcon|RiComputerLine", "", "data-icon", "inline-start"), "Open with GitHub Desktop"),
						button.Button(button.VariantGhost, button.SizeSm, h.Class("justify-start gap-1.5"),
							ip("DownloadIcon|IconDownload|DownloadIcon|DownloadIcon|RiDownloadLine", "", "data-icon", "inline-start"), "Download ZIP"),
					),
				),
			),
		),
	)
}

// usd is Intl.NumberFormat("en-US", { style: "currency", currency: "USD" }).
func usd(v float64) string {
	s := fmt.Sprintf("%.2f", v)
	whole, frac := s[:len(s)-3], s[len(s)-3:]
	var out []byte
	for i, c := range whole {
		if i > 0 && (len(whole)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, byte(c))
	}
	return "$" + string(out) + frac
}

func invoiceCard() h.Node {
	items := []struct {
		name       string
		qty, price float64
	}{{"Design System License", 1, 499}, {"Priority Support", 12, 99}, {"Custom Components", 3, 250}}
	var rows []any
	subtotal := 0.0
	for _, x := range items {
		subtotal += x.qty * x.price
		rows = append(rows, table.Row(
			table.Cell(x.name),
			table.Cell(h.Class("text-right tabular-nums"), fmt.Sprint(x.qty)),
			table.Cell(h.Class("text-right tabular-nums"), usd(x.price)),
			table.Cell(h.Class("text-right tabular-nums"), usd(x.qty*x.price)),
		))
	}
	total := func(label, value string) h.Node {
		return table.Row(
			table.Cell(h.Attr("colspan", "3"), h.Class("text-right"), label),
			table.Cell(h.Class("text-right tabular-nums"), value),
		)
	}
	rows = append(rows, total("Subtotal", usd(subtotal)), total("Tax", "$0.00"), total("Total Due", usd(subtotal)))
	return card.Card(
		card.Header(
			card.Title("Invoice #INV-2847"),
			card.Description("Due March 30, 2026"),
			card.Action(badge.Badge(badge.VariantSecondary, "Pending")),
		),
		card.Content(
			table.Table(
				table.Header(table.Row(
					table.Head("Item"),
					table.Head(h.Class("text-right"), "Qty"),
					table.Head(h.Class("text-right"), "Rate"),
					table.Head(h.Class("text-right"), "Amount"),
				)),
				table.Body(rows...),
			),
		),
		card.Footer(
			button.Button(button.VariantOutline, button.SizeSm, "Download PDF"),
			button.Button(button.SizeSm, h.Class("ml-auto"), "Pay Now"),
		),
	)
}

// --- column 2 ---

func iconPreviewGrid() h.Node {
	icons := []string{
		icCopy,
		"CircleAlertIcon|IconExclamationCircle|AlertCircleIcon|WarningCircleIcon|RiErrorWarningLine",
		"TrashIcon|IconTrash|Delete02Icon|TrashIcon|RiDeleteBinLine",
		"ShareIcon|IconShare|Share03Icon|ShareIcon|RiShareLine",
		"ShoppingBagIcon|IconShoppingBag|ShoppingBag01Icon|BagIcon|RiShoppingBagLine",
		"MoreHorizontalIcon|IconDots|MoreHorizontalCircle01Icon|DotsThreeIcon|RiMoreLine",
		"Loader2Icon|IconLoader|Loading03Icon|SpinnerIcon|RiLoaderLine",
		icPlusSign,
		"MinusIcon|IconMinus|MinusSignIcon|MinusIcon|RiSubtractLine",
		"ArrowLeftIcon|IconArrowLeft|ArrowLeft02Icon|ArrowLeftIcon|RiArrowLeftLine",
		"ArrowRightIcon|IconArrowRight|ArrowRight02Icon|ArrowRightIcon|RiArrowRightLine",
		"CheckIcon|IconCheck|Tick02Icon|CheckIcon|RiCheckLine",
		"ChevronDownIcon|IconChevronDown|ArrowDown01Icon|CaretDownIcon|RiArrowDownSLine",
		icChevronRight,
		icSearch,
		"SettingsIcon|IconSettings|Settings01Icon|GearIcon|RiSettingsLine",
	}
	var cells []any
	for _, ic := range icons {
		cells = append(cells, h.Div(h.Class("flex size-8 items-center justify-center rounded-md ring ring-border style-sera:rounded-none *:[svg]:size-4"), ip(ic, "")))
	}
	return card.Card(card.Content(h.Div(h.Class("grid grid-cols-8 place-items-center gap-4"), h.Group(cells...))))
}

func observabilityCard() h.Node {
	return card.Card(h.Class("relative w-full max-w-md overflow-hidden pt-0"),
		h.Div(h.Class("absolute inset-0 z-30 aspect-video bg-primary opacity-50 mix-blend-color")),
		h.Img(h.UnsafeAttr("src", "https://images.unsplash.com/photo-1604076850742-4c7221f3101b?q=80&w=1887&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"),
			h.Alt("Photo by mymind on Unsplash"), h.Attr("title", "Photo by mymind on Unsplash"),
			h.Class("relative z-20 aspect-video w-full object-cover brightness-60 grayscale")),
		card.Header(
			card.Title("Observability Plus is replacing Monitoring"),
			card.Description("Switch to the improved way to explore your data, with natural language. Monitoring will no longer be available on the Pro plan in November, 2025"),
		),
		card.Footer(
			button.Button("Create Query ", ip(icPlusSign, "", "data-icon", "inline-end")),
			badge.Badge(badge.VariantSecondary, h.Class("ml-auto"), "Warning"),
		),
	)
}

func shippingAddress() h.Node {
	text := func(id, label, placeholder string) h.Node {
		return field.Field(field.Label(h.For(id), label), input.Input(h.ID(id), h.Placeholder(placeholder)))
	}
	return card.Card(
		card.Header(card.Title("Shipping Address"), card.Description("Where should we deliver?")),
		card.Content(
			field.Group(
				text("shipping-street", "Street address", "123 Main Street"),
				text("shipping-apt", "Apt / Suite", "Apt 4B"),
				field.Group(h.Class("grid grid-cols-2"),
					text("shipping-city", "City", "San Francisco"),
					field.Field(field.Label(h.For("shipping-state"), "State"),
						accountSelect("shipping-state", "CA", [2]string{"CA", "California"}, [2]string{"NY", "New York"}, [2]string{"TX", "Texas"})),
				),
				field.Group(h.Class("grid grid-cols-2"),
					text("shipping-zip", "ZIP Code", "94102"),
					field.Field(field.Label(h.For("shipping-country"), "Country"),
						accountSelect("shipping-country", "US", [2]string{"US", "United States"}, [2]string{"CA", "Canada"}, [2]string{"UK", "United Kingdom"})),
				),
				field.Field(field.Horizontal,
					checkbox.Checkbox(h.ID("shipping-save"), checkbox.Checked(true)),
					field.Label(h.For("shipping-save"), h.Class("font-normal"), "Save as default address"),
				),
			),
		),
		card.Footer(
			button.Button(button.VariantOutline, button.SizeSm, "Cancel"),
			button.Button(button.SizeSm, h.Class("ml-auto"), "Save Address"),
		),
	)
}

// --- column 3 ---

func environmentVariables() h.Node {
	var rows []any
	for _, e := range []struct {
		key    string
		masked bool
	}{{"DATABASE_URL", true}, {"NEXT_PUBLIC_API", false}, {"STRIPE_SECRET", true}} {
		v := "https://api.example.com"
		if e.masked {
			v = "••••••••"
		}
		rows = append(rows, h.Div(h.Class("flex items-center gap-2 rounded-md px-2.5 py-2 font-mono text-xs ring ring-border"),
			h.Span(h.Class("font-medium"), e.key),
			h.Span(h.Class("ml-auto text-muted-foreground"), v),
		))
	}
	return card.Card(
		card.Header(card.Title("Environment Variables"), card.Description("Production · 8 variables")),
		card.Content(h.Class("flex flex-col gap-2"), h.Group(rows...)),
		card.Footer(
			button.Button(button.VariantOutline, "Edit"),
			button.Button(h.Class("ml-auto"), "Deploy"),
		),
	)
}

func barChartCard(ds designSystem) h.Node {
	months := []string{"January", "February", "March", "April", "May", "June"}
	desktop := []float64{186, 305, 237, 73, 209, 214}
	mobile := []float64{80, 200, 120, 190, 130, 140}
	var rows []chart.Row
	dt, mt := 0.0, 0.0
	for i, m := range months {
		rows = append(rows, chart.Row{"month": m, "desktop": desktop[i], "mobile": mobile[i]})
		dt += desktop[i]
		mt += mobile[i]
	}
	delta := int((dt-mt)/mt*100 + 0.5)
	prefix := ""
	if delta > 0 {
		prefix = "+"
	}
	stat := func(label, value string) h.Node {
		return h.Div(h.Class("px-2 text-center"),
			h.Div(h.Class("text-[0.65rem] text-muted-foreground uppercase"), label),
			h.Div(h.Class("text-sm font-medium tabular-nums"), value),
		)
	}
	return card.Card(
		card.Header(
			card.Title(h.Class("text-lg"), "Traffic channels"),
			card.Description(h.Class("line-clamp-2 text-sm leading-snug"), "Monthly desktop and mobile traffic for the last six months—compare volume and mix across platforms and devices at a glance."),
		),
		card.Content(h.Class("flex flex-col gap-4 pt-0"),
			chart.Container(chart.Config{{Key: "desktop", Label: "Desktop", Color: "var(--chart-1)"}, {Key: "mobile", Label: "Mobile", Color: "var(--chart-2)"}},
				h.Class("max-h-[180px] w-full"),
				chart.BarChart(chart.AccessibilityLayer(), chart.Data(rows), chart.P("margin", map[string]any{"left": 0, "right": 0, "top": 8, "bottom": 0}),
					chart.CartesianGrid(chart.Vertical(false), chart.P("strokeDasharray", "3 3")),
					chart.XAxis(chart.DataKey("month"), chart.TickLine(false), chart.TickMargin(8), chart.AxisLine(false),
						chart.TickFormatter(chart.Fn{Name: "slice", Args: []any{0, 3}})),
					chart.Tooltip(chart.Cursor(false), chart.TooltipContent(chart.Indicator("dashed"))),
					chart.Legend(chart.LegendContent()),
					chart.Bar(chart.DataKey("desktop"), chart.Fill("var(--color-desktop)"), barRadius(rounded(ds), 6)),
					chart.Bar(chart.DataKey("mobile"), chart.Fill("var(--color-mobile)"), barRadius(rounded(ds), 6)),
				),
			),
			h.Div(h.Class("grid w-full grid-cols-3 divide-x divide-border/60"),
				stat("Desktop", thousands(dt)),
				stat("Mobile", thousands(mt)),
				stat("Mix Delta", fmt.Sprintf("%s%d%%", prefix, delta)),
			),
		),
		card.Footer(button.Button(h.Class("w-full"), "View report")),
	)
}

// thousands is Number.prototype.toLocaleString() for integers (en-US).
func thousands(v float64) string {
	s := fmt.Sprintf("%.0f", v)
	var out []byte
	for i := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, s[i])
	}
	return string(out)
}

func inviteTeam() h.Node {
	var rows []any
	for _, x := range [][2]string{{"alex@example.com", "editor"}, {"sam@example.com", "viewer"}} {
		rows = append(rows, h.Div(h.Class("flex items-center gap-2"),
			input.Input(h.Value(x[0]), h.Class("flex-1")),
			selectui.Select(selectui.DefaultValue(x[1]),
				selectui.Trigger(h.Class("w-24"), selectui.Value()),
				selectui.Content(selectui.PositionPopper, selectui.AlignEnd, selectui.Group(
					selectui.Item(selectui.ItemValue("admin"), "Admin"),
					selectui.Item(selectui.ItemValue("editor"), "Editor"),
					selectui.Item(selectui.ItemValue("viewer"), "Viewer"),
				)),
			),
		))
	}
	return card.Card(
		card.Header(card.Title("Invite Team"), card.Description("Add members to your workspace")),
		card.Content(h.Class("flex flex-col gap-4"),
			h.Div(h.Class("flex flex-col gap-3"), h.Group(rows...)),
			button.Button(button.VariantOutline, ip(icPlusSign, "", "data-icon", "inline-start"), "Add another"),
			separator.Separator(),
			field.Field(
				field.Label(h.For("invite-link"), "Or share invite link"),
				inputgroup.InputGroup(
					inputgroup.Input(h.ID("invite-link"), h.Value("https://app.co/invite/x8f2k"), h.Bool("readonly")),
					inputgroup.Addon(inputgroup.AlignInlineEnd,
						inputgroup.Button(inputgroup.ButtonSizeIconXs, h.Aria("label", "Copy link"), ip(icCopy, "")),
					),
				),
			),
		),
		card.Footer(button.Button(h.Class("w-full"), "Send Invites")),
	)
}

func activateAgentDialog() h.Node {
	feature := func(content ...any) h.Node {
		return item.Item(item.SizeXs, h.Class("px-0"),
			item.Media(item.MediaIcon, h.Class("self-start"),
				ip("CheckCircle2Icon|IconCircleCheckFilled|CheckmarkCircle02Icon|CheckCircleIcon|RiCheckboxCircleLine", "size-5 fill-primary text-primary-foreground")),
			item.Content(item.Title(h.Class("inline leading-relaxed font-normal text-muted-foreground *:[strong]:font-medium *:[strong]:text-foreground"), h.Group(content...))),
		)
	}
	return card.Card(
		card.Header(
			card.Title("Ship faster & safer with Vercel Agent"),
			card.Description("Your use is subject to Vercel's ", h.A(h.Href("#"), "Public Beta Agreement"), " and ", h.A(h.Href("#"), "AI Product Terms"), "."),
		),
		card.Content(h.Class("flex flex-col gap-4"),
			item.Group(h.Class(itemGroupGap0),
				feature(h.Strong("Code reviews"), " with full codebase context to catch ", h.Strong("hard-to-find"), " bugs."),
				feature(h.Strong("Code suggestions"), " validated in sandboxes before you merge."),
				feature(h.Strong("Root-cause analysis"), " for production issues with deployment context. ",
					badge.Badge(badge.VariantSecondary, h.Class("bg-chart-1 text-chart-5"), "Requires Observability Plus")),
			),
			alert.Alert(alert.Description("Pro teams get $100 in Vercel Agent trial credit for 2 weeks after activation.")),
		),
		card.Footer(h.Class("justify-end gap-2"),
			button.Button(button.VariantOutline, h.Class("style-sera:hidden"), "Cancel"),
			button.Button(h.Class("style-sera:w-full"), "Enable with $100 credits"),
		),
	)
}

// --- column 4 ---

func skeletonLoading() h.Node {
	sk := func(class string) h.Node { return skeleton.Skeleton(h.Class(class)) }
	return card.Card(
		card.Content(h.Class("flex flex-col gap-4"),
			h.Div(h.Class("flex items-center gap-3"),
				sk("size-10 rounded-full"),
				h.Div(h.Class("flex flex-1 flex-col gap-2"), sk("h-4 w-3/4"), sk("h-3 w-1/2")),
			),
			h.Div(h.Class("flex flex-col gap-2"), sk("h-3 w-full"), sk("h-3 w-full"), sk("h-3 w-4/5")),
			h.Div(h.Class("flex gap-2"), sk("h-8 w-20"), sk("h-8 w-20")),
		),
	)
}

func pieChartCard() h.Node {
	rows := []chart.Row{
		{"browser": "chrome", "visitors": 275, "fill": "var(--color-chrome)"},
		{"browser": "safari", "visitors": 200, "fill": "var(--color-safari)"},
		{"browser": "firefox", "visitors": 287, "fill": "var(--color-firefox)"},
		{"browser": "edge", "visitors": 173, "fill": "var(--color-edge)"},
	}
	// The top browser (firefox, 287 of 935) and its share.
	total, top := 935.0, 287.0
	share := int(top/total*100 + 0.5)
	return card.Card(
		card.Header(h.Class("pb-0"),
			card.Title("Browser Share"),
			card.Description("January - June 2026"),
			card.Action(badge.Badge(badge.VariantOutline, "Firefox")),
		),
		card.Content(h.Class("pt-0"),
			chart.Container(chart.Config{
				{Key: "visitors", Label: "Visitors"},
				{Key: "chrome", Label: "Chrome", Color: "var(--chart-1)"},
				{Key: "safari", Label: "Safari", Color: "var(--chart-2)"},
				{Key: "firefox", Label: "Firefox", Color: "var(--chart-3)"},
				{Key: "edge", Label: "Edge", Color: "var(--chart-4)"},
			}, h.Class("mx-auto aspect-square max-h-[190px]"),
				chart.PieChart(
					noCursorTooltip(),
					chart.Pie(chart.Data(rows), chart.DataKey("visitors"), chart.NameKey("browser"), chart.InnerRadius(50), chart.StrokeWidth(5),
						chart.Label(chart.Content(chart.Fn{Name: "pieCenterLines", Args: []any{
							[]any{thousands(total), "fill-foreground text-2xl font-bold", -16},
							[]any{"Visitors", "fill-muted-foreground text-xs", 4},
						}})),
					),
					chart.Legend(chart.LegendContent(chart.NameKey("browser")), chart.ClassName("translate-y-2")),
				),
			),
		),
		card.Footer(h.Class("flex-col items-stretch gap-2"),
			h.Div(h.Class("flex items-center text-xs"),
				h.Span(h.Class("font-medium"), "Firefox"),
				h.Span(h.Class("ml-auto text-muted-foreground tabular-nums"), fmt.Sprintf("%d%%", share)),
			),
			progress.Progress(progress.Value(float64(share)), h.Class("**:data-[slot=progress-indicator]:bg-chart-3")),
		),
	)
}

func githubAvatar(user, fallback string, extra ...any) h.Node {
	return avatar.Avatar(h.Group(extra...),
		avatar.Image(h.Src("https://github.com/"+user+".png"), h.Alt(user)),
		avatar.Fallback(fallback),
	)
}

func noTeamMembers() h.Node {
	return card.Card(
		card.Content(
			empty.Empty(h.Class("h-56 border"),
				empty.Header(
					empty.Media(
						avatar.Group(h.Class("grayscale"),
							avatar.Avatar(avatar.SizeLg, avatar.Image(h.Src("https://github.com/shadcn.png"), h.Alt("@shadcn")), avatar.Fallback("CN")),
							avatar.Avatar(avatar.SizeLg, avatar.Image(h.Src("https://github.com/maxleiter.png"), h.Alt("@maxleiter")), avatar.Fallback("LR")),
							avatar.Avatar(avatar.SizeLg, avatar.Image(h.Src("https://github.com/evilrabbit.png"), h.Alt("@evilrabbit")), avatar.Fallback("ER")),
						),
					),
					empty.Title("No Team Members"),
					empty.Description("Invite your team to collaborate on this project."),
				),
				empty.Content(button.Button(button.SizeSm, "Invite Members")),
			),
		),
	)
}

func reportBug() h.Node {
	return card.Card(
		card.Header(card.Title("Report Bug"), card.Description("Help us fix issues faster.")),
		card.Content(
			field.Group(
				field.Field(field.Label(h.For("bug-title"), "Title"), input.Input(h.ID("bug-title"), h.Placeholder("Brief description of the issue"))),
				h.Div(h.Class("grid grid-cols-2 gap-3"),
					field.Field(field.Label(h.For("bug-severity"), "Severity"),
						accountSelect("bug-severity", "medium", [2]string{"critical", "Critical"}, [2]string{"high", "High"}, [2]string{"medium", "Medium"}, [2]string{"low", "Low"})),
					field.Field(field.Label(h.For("bug-component"), "Component"),
						accountSelect("bug-component", "dashboard", [2]string{"dashboard", "Dashboard"}, [2]string{"auth", "Auth"}, [2]string{"api", "API"}, [2]string{"billing", "Billing"})),
				),
				field.Field(field.Label(h.For("bug-steps"), "Steps to reproduce"),
					textarea.Textarea(h.ID("bug-steps"), h.Placeholder("1. Go to\n2. Click on\n3. Observe..."), h.Class("min-h-24 resize-none"))),
			),
		),
		card.Footer(
			field.Field(field.Horizontal, h.Class("justify-end style-sera:justify-center"),
				button.Button(button.VariantOutline, h.Class("style-sera:flex-1"), "Attach File"),
				button.Button(h.Class("style-sera:flex-1"), "Submit Bug"),
			),
		),
	)
}

func contributors() h.Node {
	var avatars []any
	for _, u := range []string{"shadcn", "vercel", "nextjs", "tailwindlabs", "typescript-lang", "eslint", "prettier", "babel", "webpack", "rollup", "parcel", "vite", "react", "vue", "angular", "solid"} {
		avatars = append(avatars, githubAvatar(u, u[:1], h.Class("grayscale")))
	}
	return card.Card(h.Class("max-w-sm"),
		card.Header(card.Title("Contributors ", badge.Badge(badge.VariantSecondary, "312"))),
		card.Content(h.Div(h.Class("flex flex-wrap gap-2"), h.Group(avatars...))),
		card.Footer(h.A(h.Href("#"), h.Class("text-sm underline underline-offset-4"), "+ 810 contributors")),
	)
}

// --- column 5 ---

func feedbackForm() h.Node {
	topics := [][2]string{
		{"", "Select a topic"}, {"ai", "AI"}, {"accounts-and-access-controls", "Accounts and Access Controls"},
		{"billing", "Billing"}, {"cdn", "CDN (Firewall, Caching)"}, {"ci-cd", "CI/CD (Builds, Deployments, Environment Variables)"},
		{"dashboard-interface", "Dashboard Interface (Navigation, UI Issues)"}, {"domains", "Domains"}, {"frameworks", "Frameworks"},
		{"marketplace-and-integrations", "Marketplace and Integrations"}, {"observability", "Observability (Observability, Logs, Monitoring)"},
		{"storage", "Storage"},
	}
	var options []any
	for _, t := range topics {
		options = append(options, nativeselect.Option(h.Value(t[0]), t[1]))
	}
	return card.Card(
		card.Content(
			h.Form(h.ID("feedback-form"),
				field.Group(
					field.Field(field.Label(h.For("topic"), "Topic"), nativeselect.NativeSelect(h.ID("topic"), h.Group(options...))),
					field.Field(field.Label(h.For("feedback"), "Feedback"), textarea.Textarea(h.ID("feedback"), h.Placeholder("Your feedback helps us improve..."))),
				),
			),
		),
		card.Footer(button.Button(h.Type("submit"), h.Attr("form", "feedback-form"), h.Class("style-sera:w-full"), "Submit")),
	)
}

func bookAppointment() h.Node {
	var slots []any
	for i, t := range []string{"9:00 AM", "10:30 AM", "11:00 AM", "1:30 PM"} {
		slots = append(slots, togglegroup.Item(togglegroup.Value(fmt.Sprintf("slot-%d", i)), t))
	}
	return card.Card(
		card.Header(card.Title("Book Appointment"), card.Description("Dr. Sarah Chen · Cardiology")),
		card.Content(h.Class("flex flex-col gap-4"),
			field.Group(
				field.Field(
					field.Label("Available on March 18, 2026"),
					togglegroup.ToggleGroup(togglegroup.TypeMultiple, togglegroup.Spacing(2), togglegroup.Value("slot-0"), h.Group(slots...)),
				),
			),
			alert.Alert(alert.Title("New patient?"), alert.Description("Please arrive 15 minutes early.")),
		),
		card.Footer(button.Button(h.Class("w-full"), "Book Appointment")),
	)
}

func sleepReport(ds designSystem) h.Node {
	type r struct {
		hour             string
		deep, light, rem float64
	}
	data := []r{{"10pm", 0, 30, 0}, {"11pm", 20, 10, 0}, {"12am", 40, 0, 10}, {"1am", 30, 5, 15}, {"2am", 10, 20, 30}, {"3am", 25, 10, 20}, {"4am", 15, 25, 10}, {"5am", 5, 35, 15}, {"6am", 0, 20, 25}}
	var rows []chart.Row
	for _, x := range data {
		rows = append(rows, chart.Row{"hour": x.hour, "deep": x.deep, "light": x.light, "rem": x.rem})
	}
	var stats []any
	for _, s := range [][2]string{{"Deep", "2h 10m"}, {"Light", "3h 48m"}, {"REM", "1h 26m"}, {"Score", "84"}} {
		stats = append(stats, h.Div(h.Class("text-center"),
			h.Div(h.Class("text-sm font-medium tabular-nums"), s[1]),
			h.Div(h.Class("text-xs text-muted-foreground"), s[0]),
		))
	}
	return card.Card(
		card.Header(card.Title("Sleep Report"), card.Description("Last night · 7h 24m")),
		card.Content(h.Class("flex flex-col gap-3"),
			chart.Container(chart.Config{
				{Key: "deep", Label: "Deep", Color: "var(--chart-1)"},
				{Key: "light", Label: "Light", Color: "var(--chart-2)"},
				{Key: "rem", Label: "REM", Color: "var(--chart-3)"},
			}, h.Class("h-32 w-full"),
				chart.BarChart(chart.AccessibilityLayer(), chart.Data(rows), chart.P("margin", zeroMargin(0)), chart.P("barSize", 16),
					chart.Bar(chart.DataKey("deep"), chart.StackID("a"), chart.Fill("var(--color-deep)"), chart.Radius(0)),
					chart.Bar(chart.DataKey("light"), chart.StackID("a"), chart.Fill("var(--color-light)"), chart.Radius(0)),
					chart.Bar(chart.DataKey("rem"), chart.StackID("a"), chart.Fill("var(--color-rem)"), chart.Radii(2, 2, 0, 0)),
				),
			),
			h.Div(h.Class("grid grid-cols-4 gap-2"), h.Group(stats...)),
		),
		card.Footer(
			badge.Badge(badge.VariantOutline, "Good"),
			button.Button(button.VariantOutline, button.SizeSm, h.Class("ml-auto"), "Details"),
		),
	)
}

func githubProfile() h.Node {
	return card.Card(h.Class("mx-auto w-full max-w-md"),
		card.Header(card.Title("Profile"), card.Description("Manage your profile information.")),
		card.Content(
			h.Form(h.ID("profile"),
				field.Group(
					field.Field(
						field.Label(h.For("name"), "Name"),
						input.Input(h.ID("name"), h.Placeholder("shadcn")),
						field.Description("Your name may appear around GitHub where you contribute or are mentioned. You can remove it at any time."),
					),
					field.Field(
						field.Label(h.For("email"), "Public Email"),
						nativeselect.NativeSelect(h.ID("email"),
							nativeselect.Option(h.Value("m@shadcn.com"), "m@shadcn.com"),
							nativeselect.Option(h.Value("m@gmail.com"), "m@gmail.com"),
						),
						field.Description("You can manage verified email addresses in your ", h.A(h.Href("#email-settings"), "email settings"), "."),
					),
					field.Field(
						field.Label(h.For("bio"), "Bio"),
						textarea.Textarea(h.ID("bio"), h.Placeholder("Tell us a little bit about yourself")),
						field.Description("You can ", h.Span("@mention"), " other users and organizations to link to them."),
					),
				),
			),
		),
		card.Footer(button.Button(h.Attr("form", "profile"), h.Class("style-sera:w-full"), "Save Profile")),
	)
}

// --- column 6 ---

func weeklyFitnessSummary(ds designSystem) h.Node {
	var days []any
	for _, d := range []struct {
		day  string
		load int
	}{{"M", 84}, {"T", 52}, {"W", 73}, {"T", 66}, {"F", 91}, {"S", 48}, {"S", 61}} {
		days = append(days, h.Div(h.Class("rounded-md p-1.5 text-center ring ring-border"),
			h.Div(h.Class("text-sm text-muted-foreground"), d.day),
			h.Div(h.Class("relative mt-1 h-16 overflow-hidden rounded-sm bg-muted"),
				h.Div(h.Class("absolute inset-x-0 bottom-0 rounded-sm bg-chart-3"), h.Style(fmt.Sprintf("height: %d%%", d.load))),
			),
		))
	}
	return card.Card(
		card.Header(card.Title("Weekly Fitness Summary"), card.Description("Calories and workout load by day")),
		card.Content(h.Class("flex flex-col gap-4"), h.Div(h.Class("grid grid-cols-7 gap-1.5"), h.Group(days...))),
		card.Footer(button.Button(h.Class("w-full"), "View details")),
	)
}

func fileUpload() h.Node {
	return card.Card(
		card.Header(card.Title("File Upload"), card.Description("Drag and drop or browse")),
		card.Content(
			empty.Empty(h.Class("border"),
				empty.Header(
					empty.Media(empty.MediaIcon, ip("UploadCloudIcon|IconCloudUpload|CloudUploadIcon|CloudArrowUpIcon|RiUploadCloudLine", "")),
					empty.Title("Upload files"),
					empty.Description("PNG, JPG, PDF up to 10MB"),
				),
				empty.Content(button.Button("Browse Files")),
			),
		),
	)
}

var sixMonths = []string{"January", "February", "March", "April", "May", "June"}

func previewAnalyticsCard() h.Node {
	var rows []chart.Row
	for i, v := range []float64{186, 305, 237, 73, 209, 214} {
		rows = append(rows, chart.Row{"month": sixMonths[i], "visitors": v})
	}
	return card.Card(card.SizeSm, h.Class("mx-auto w-full max-w-sm data-[size=sm]:pb-0"),
		card.Header(
			card.Title("Analytics"),
			card.Description("418.2K Visitors ", badge.Badge("+10%")),
			card.Action(button.Button(button.VariantOutline, button.SizeSm, "View Analytics")),
		),
		chart.Container(chart.Config{{Key: "visitors", Label: "Visitors", Color: "var(--chart-1)"}}, h.Class("aspect-[1/0.35]"),
			chart.AreaChart(chart.AccessibilityLayer(), chart.Data(rows), chart.P("margin", map[string]any{"left": 0, "right": 0}),
				chart.Tooltip(chart.Cursor(false), chart.TooltipContent(chart.Indicator("line"), chart.HideLabel()), chart.P("defaultIndex", 2)),
				chart.Area(chart.DataKey("visitors"), chart.Type("linear"), chart.Fill("var(--color-visitors)"), chart.FillOpacity(0.4), chart.Stroke("var(--color-visitors)")),
			),
		),
	)
}

// circularGauge is usage-card.tsx's CircularGauge.
func circularGauge(percentage float64) h.Node {
	c := 2 * 3.141592653589793 * 42.5
	stroke := percentage / 100 * c
	circle := func(class string, dash string) h.Node {
		return h.El("circle", h.Attr("cx", "50"), h.Attr("cy", "50"), h.Attr("r", "42.5"), h.Attr("stroke-width", "12"),
			h.Attr("stroke-dashoffset", "0"), h.Attr("stroke-linecap", "round"), h.Attr("stroke-linejoin", "round"),
			h.Class(class), h.Attr("stroke", "currentColor"), h.Style("stroke-dasharray: "+dash))
	}
	return h.El("svg", h.Aria("hidden", "true"), h.Attr("fill", "none"), h.Attr("height", "16"), h.Attr("width", "16"),
		h.Attr("stroke-width", "2"), h.Attr("viewBox", "0 0 100 100"), h.Class("-rotate-90"),
		circle("opacity-20", jsNum(c)+" "+jsNum(c)),
		circle("transition-all duration-300", jsNum(stroke)+" "+jsNum(c)),
	)
}

// jsNum formats a float as JavaScript's String(number) does for these values.
func jsNum(v float64) string {
	s := fmt.Sprintf("%.15g", v)
	if len(s) > 17 {
		s = fmt.Sprintf("%.14g", v)
	}
	return s
}

func usageCard() h.Node {
	items := []struct {
		name, value string
		pct         float64
	}{
		{"Edge Requests", "$1.83K", 67.34}, {"Fast Data Transfer", "$952.51", 52.18},
		{"Monitoring data points", "$901.20", 89.42}, {"Web Analytics Events", "$603.71", 45.67},
		{"ISR Writes", "524.52K / 2M", 26.23}, {"Function Duration", "5.11 GB Hrs / 1K GB Hrs", 5.11},
	}
	var list []any
	for _, x := range items {
		list = append(list, item.Item(item.SizeXs, lib.AsChild, h.Class("px-0 group-hover/item-group:bg-transparent"),
			h.A(h.Href("#"),
				item.Media(item.MediaIcon, h.Class("text-primary"), circularGauge(x.pct)),
				item.Content(h.Class("inline-block truncate"), item.Title(h.Class("inline"), x.name)),
				item.Actions(h.Span(h.Class("font-mono text-xs font-medium text-muted-foreground tabular-nums"), x.value)),
			),
		))
	}
	return card.Card(h.Class("w-full max-w-sm gap-4"),
		card.Header(card.Title(h.Class("px-1 text-sm"), "5 days remaining in cycle")),
		card.Content(item.Group(h.Class(itemGroupGap0), h.Group(list...))),
	)
}

func shortcuts() h.Node {
	var list []any
	for i, s := range []struct{ label, key string }{{"Search", "K"}, {"Quick Actions", "J"}, {"New File", "N"}, {"Save", "S"}, {"Toggle Sidebar", "B"}} {
		if i > 0 {
			list = append(list, item.Separator())
		}
		list = append(list, item.Item(item.VariantDefault, item.SizeXs, h.Class("border-0 px-0 py-0"),
			item.Header(
				item.Title(h.Class("font-normal"), s.label),
				item.Actions(h.Div(h.Class("flex gap-1"), kbd.Kbd("⌘"), kbd.Kbd(s.key))),
			),
		))
	}
	return card.Card(
		card.Content(
			h.Div(h.Class("flex flex-col gap-3"),
				h.Div(h.Class("text-sm font-medium"), "Shortcuts"),
				item.Group(h.Class("gap-2 text-muted-foreground"), h.Data("size", "xs"), h.Group(list...)),
			),
		),
	)
}

// --- column 7 ---

func anomalyAlert() h.Node {
	return card.Card(
		card.Content(
			empty.Empty(h.Class("h-48"),
				empty.Header(
					empty.Title("Get alerted for anomalies"),
					empty.Description("Automatically monitor your projects for anomalies and get notified."),
				),
				empty.Content(button.Button("Upgrade to Observability Plus")),
			),
		),
	)
}

// liveWaveformCard is live-waveform.tsx; create-preview.js draws the canvas
// (the processing animation, the microphone when listening).
func liveWaveformCard() h.Node {
	return card.Card(h.Data("create-waveform", ""),
		card.Header(
			card.Title("Live Audio Waveform"),
			card.Description("Real-time microphone input visualization with audio reactivity"),
		),
		card.Content(
			h.Div(h.Class("relative h-full w-full"), h.Style("height: 80px"), h.Aria("label", "Processing audio"), h.Role("img"),
				h.Data("create-waveform-canvas", ""),
				h.El("canvas", h.Class("block h-full w-full"), h.Aria("hidden", "true")),
			),
		),
		card.Footer(h.Class("gap-2"),
			button.Button(button.SizeSm, button.VariantOutline, h.Data("create-waveform-action", "listen"), "Start Listening"),
			button.Button(button.SizeSm, h.Data("create-waveform-action", "process"), "Stop Processing"),
			button.Button(button.SizeSm, button.VariantOutline, h.Data("create-waveform-action", "mode"), "Static"),
		),
	)
}

func visitors() h.Node {
	var rows []chart.Row
	for i, v := range []float64{186, 305, 237, 73, 209, 214} {
		rows = append(rows, chart.Row{"month": sixMonths[i], "desktop": v})
	}
	// (214 - 209) / 209, rounded.
	return card.Card(h.Class("pb-0"),
		card.Header(
			card.Title("Visitors"),
			card.Description("Last 6 months "),
			card.Action(badge.Badge(badge.VariantSecondary, "+2% vs last month")),
		),
		card.Content(h.Class("px-0"),
			chart.Container(chart.Config{{Key: "desktop", Label: "Desktop", Color: "var(--chart-1)"}}, h.Class("h-48 w-full"),
				chart.AreaChart(chart.AccessibilityLayer(), chart.Data(rows), chart.P("margin", map[string]any{"left": 0, "right": 0, "top": 6, "bottom": 0}),
					chart.XAxis(chart.DataKey("month"), chart.TickLine(false), chart.Hide(), chart.AxisLine(false), chart.TickMargin(6),
						chart.TickFormatter(chart.Fn{Name: "slice", Args: []any{0, 3}})),
					chart.Tooltip(chart.Cursor(false), chart.TooltipContent(chart.Indicator("line"))),
					chart.Area(chart.DataKey("desktop"), chart.Type("natural"), chart.Fill("var(--color-desktop)"), chart.FillOpacity(0.15), chart.Stroke("var(--color-desktop)"), chart.StrokeWidth(2)),
				),
			),
		),
	)
}

func contributionsActivity() h.Node {
	return card.Card(
		card.Header(card.Title("Contributions & Activity"), card.Description("Manage your contributions and activity visibility.")),
		card.Content(
			h.Form(h.ID("contributions-activity"),
				field.Group(
					field.Set(
						field.Legend(h.Class("sr-only"), "Contributions & activity"),
						field.Group(
							field.Field(field.Horizontal,
								checkbox.Checkbox(h.ID("activity-private-profile")),
								field.Content(
									field.Label(h.For("activity-private-profile"), "Make profile private and hide activity"),
									field.Description("Enabling this will hide your contributions and activity from your GitHub profile and from social features like followers, stars, feeds, leaderboards and releases."),
								),
							),
						),
					),
				),
			),
		),
		card.Footer(button.Button(h.Attr("form", "contributions-activity"), h.Class("style-sera:w-full"), "Save Changes")),
	)
}

func notFoundCard() h.Node {
	return card.Card(
		card.Content(
			empty.Empty(h.Class("h-72"),
				empty.Header(
					empty.Title("404 - Not Found"),
					empty.Description("The page you're looking for doesn't exist. Try searching for what you need below."),
				),
				empty.Content(
					inputgroup.InputGroup(h.Class("w-3/4"),
						inputgroup.Input(h.Placeholder("Try searching for pages...")),
						inputgroup.Addon(ip(icSearch, "")),
						inputgroup.Addon(inputgroup.AlignInlineEnd, kbd.Kbd("/")),
					),
					button.Button(button.VariantLink, "Go to homepage"),
				),
			),
		),
	)
}

// itemGroupGap0 is className="gap-0" on an ItemGroup of xs items. Upstream's
// preview styles the group's token in the base layer, so the gap-0 utility
// beats its has-data-[size=xs]:gap-2; the generated classes need the
// variant cleared too.
const itemGroupGap0 = "gap-0 has-data-[size=xs]:gap-0"
