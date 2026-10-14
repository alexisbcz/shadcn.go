package site

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/badge"
	"shadcn.go/components/ui/breadcrumb"
	"shadcn.go/components/ui/button"
	buttongroup "shadcn.go/components/ui/button-group"
	"shadcn.go/components/ui/calendar"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/checkbox"
	cb "shadcn.go/components/ui/combobox"
	dm "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/components/ui/empty"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/input"
	inputgroup "shadcn.go/components/ui/input-group"
	"shadcn.go/components/ui/item"
	"shadcn.go/components/ui/pagination"
	"shadcn.go/components/ui/popover"
	rg "shadcn.go/components/ui/radio-group"
	selectui "shadcn.go/components/ui/select"
	"shadcn.go/components/ui/separator"
	"shadcn.go/components/ui/slider"
	"shadcn.go/components/ui/table"
	"shadcn.go/components/ui/textarea"
	h "shadcn.go/libhtml"
)

// The Sera page's lazy previews (components/lazy-preview.tsx loads them
// once they come within 800px of the viewport): each renders into a
// <template> that static/sera.js swaps in for the placeholder.

func seraPreviews() map[string]func() h.Node {
	return map[string]func() h.Node{
		"articleDirectory":  seraArticleDirectory,
		"emptyState":        seraEmptyState,
		"editArticle":       seraEditArticle,
		"mediaLibrary":      seraMediaLibrary,
		"mediaLibraryTable": seraMediaLibraryTable,
	}
}

// seraBackLink is the preview headers' breadcrumb back link.
func seraBackLink(label, listClass, linkClass, iconClass string) h.Node {
	return breadcrumb.Breadcrumb(
		breadcrumb.List(h.If(listClass != "", h.Class(listClass)),
			breadcrumb.Item(
				breadcrumb.Link(h.Href("#"), h.Class(linkClass), lucide.ArrowLeft(h.Class(iconClass)), label),
			),
		),
	)
}

func seraHeading(text string) h.Node {
	return h.H1(h.Class("line-clamp-1 font-heading text-3xl tracking-wide uppercase md:text-3xl lg:text-4xl"), text)
}

func seraHeader(left []any, right ...any) h.Node {
	return h.Header(
		h.Div(h.Class("container flex flex-col items-center justify-center gap-(--gap) py-(--gap) sm:flex-row sm:justify-between"),
			h.Div(h.Class("flex flex-col gap-2 text-center sm:text-left"), h.Group(left...)),
			h.Group(right...),
		),
	)
}

func seraSection(header h.Node, body h.Node) h.Node {
	return h.Div(h.Class(seraPreviewClass),
		header,
		separator.Separator(h.Class("hidden sm:block")),
		body,
	)
}

func seraSearch(placeholder string, props ...any) h.Node {
	return inputgroup.InputGroup(h.Group(props...),
		inputgroup.Addon(lucide.Search()),
		inputgroup.Input(h.Placeholder(placeholder)),
	)
}

// --- article directory ----------------------------------------------------------

var seraArticles = []struct {
	title, words, author, issue, status, label string
	progress                                   int
}{
	{"The Future of Sustainable Architecture", "1.4k / 2.6k words", "Elena Rostova", "Summer 2024", "in-revision", "In revision", 45},
	{"Brutalism's Second Act", "2.1k / 2.5k words", "Marcus Chen", "Summer 2024", "final-edit", "Final edit", 90},
	{"The Typography of Public Spaces", "0.5k / 1.5k words", "Sarah Jenkins", "Autumn 2024", "drafting", "Drafting", 20},
	{"Rethinking Urban Canopies", "1.8k / 1.8k words", "David O'Connor", "Summer 2024", "published", "Published", 100},
	{"Light, Glass, and the Modern Museum", "1.2k / 2.0k words", "Amara Osei", "Autumn 2024", "in-revision", "In revision", 55},
	{"Concrete Utopias: Housing in the 21st Century", "3.0k / 3.0k words", "Tomás Herrera", "Summer 2024", "published", "Published", 100},
	{"Designing for Silence", "0.8k / 2.2k words", "Ingrid Solberg", "Winter 2024", "drafting", "Drafting", 30},
	{"The Invisible Infrastructure of Cities", "2.4k / 2.8k words", "James Whitfield", "Autumn 2024", "final-edit", "Final edit", 85},
}

var seraStatusVariant = map[string]badge.Variant{
	"in-revision": badge.VariantOutline,
	"final-edit":  badge.VariantDefault,
	"drafting":    badge.VariantGhost,
	"published":   badge.VariantSecondary,
}

var seraStatusDot = map[string]string{
	"in-revision": "bg-amber-600/80",
	"final-edit":  "bg-foreground/90",
	"drafting":    "bg-muted-foreground/60",
	"published":   "bg-emerald-600/80",
}

// seraProgressValue is the base-sera Progress with only a value, laid out
// by className.
func seraProgressValue(label string, value int, class string) h.Node {
	v := fmt.Sprint(value)
	return h.Div(h.Aria("label", label), h.Aria("valuemax", "100"), h.Aria("valuemin", "0"),
		h.Aria("valuenow", v), h.Aria("valuetext", v+"%"), h.Class(lib.Cn("flex flex-wrap gap-3", class)), h.Data("progressing", ""),
		h.Data("slot", "progress"), h.Role("progressbar"),
		h.Span(h.Aria("hidden", "true"), h.Class("ml-auto text-sm text-muted-foreground tabular-nums"), h.Data("progressing", ""), h.Data("slot", "progress-value"), v+"%"),
		h.Div(h.Class("relative flex h-0.5 w-full items-center overflow-x-hidden rounded-none bg-muted"), h.Data("progressing", ""), h.Data("slot", "progress-track"),
			h.Div(h.Class("h-full bg-primary transition-all"), h.Data("progressing", ""), h.Data("slot", "progress-indicator"),
				h.Style("height: inherit; inset-inline-start: 0px; width: "+v+"%")),
		),
	)
}

func seraArticleDirectory() h.Node {
	var rows []any
	for _, a := range seraArticles {
		rows = append(rows, table.Row(
			table.Cell(
				h.Div(h.Class("flex flex-col gap-1"),
					h.P(h.Class("font-heading text-xl tracking-tight text-foreground"), a.title),
					h.P(h.Class("text-xs text-muted-foreground"), a.words),
				),
			),
			table.Cell(a.author),
			table.Cell(a.issue),
			table.Cell(badge.Badge(seraStatusVariant[a.status],
				h.Span(h.Class(lib.Cn("size-1.5 rounded-full", seraStatusDot[a.status]))), a.label)),
			table.Cell(seraProgressValue(fmt.Sprintf("%d%% complete", a.progress), a.progress, "flex flex-row-reverse items-center **:data-[slot=progress-track]:w-16")),
		))
	}
	pages := []any{
		pagination.Item(pagination.Link(h.Href("#"), button.SizeIconSm, h.Aria("label", "Previous page"), lucide.ChevronLeft(h.Class("cn-rtl-flip")))),
	}
	for i := 1; i <= 3; i++ {
		l := pagination.Link(h.Href("#"), button.SizeIconSm, fmt.Sprint(i))
		if i == 1 {
			l.Active(true)
		} else {
			l.Active(false)
		}
		pages = append(pages, pagination.Item(l))
	}
	pages = append(pages, pagination.Item(pagination.Link(h.Href("#"), button.SizeIconSm, h.Aria("label", "Next page"), lucide.ChevronRight(h.Class("cn-rtl-flip")))))
	return seraSection(
		seraHeader([]any{
			seraBackLink("Editorial Dashboard", "justify-center md:justify-start", "inline-flex items-center gap-1.5", "size-3"),
			seraHeading("Article Directory"),
		}, h.Div(
			buttongroup.ButtonGroup(h.Class("gap-2 sm:ml-auto md:gap-4"),
				button.Button(lucide.Plus(h.Data("icon", "inline-start")), "New Article"),
			),
		)),
		h.Div(h.Class("container py-(--gap)"),
			card.Card(
				card.Header(
					inputgroup.InputGroup(
						inputgroup.Addon(lucide.Search()),
						inputgroup.Input(h.Type("search"), h.Placeholder("Search articles...")),
					),
				),
				card.Content(
					table.Table(
						table.Header(table.Row(h.Class("hover:bg-transparent"),
							table.Head("Title"),
							table.Head(h.Class("w-[170px]"), "Author"),
							table.Head(h.Class("w-[150px]"), "Issue"),
							table.Head(h.Class("w-[180px]"), "Status"),
							table.Head(h.Class("w-[140px]"), "Progress"),
						)),
						table.Body(h.Group(rows...)),
					),
				),
				card.Footer(
					pagination.Pagination(pagination.Content(h.Group(pages...))),
				),
			),
		),
	)
}

// --- empty state ----------------------------------------------------------------

func seraEmptyState() h.Node {
	stages := []struct{ label, description, dot string }{
		{"Drafting", "Start the writing process. Articles here are works in progress, visible only to editors and authors.", "bg-amber-600"},
		{"In Revision", "Content undergoing editorial review. Track changes and word counts as pieces take shape.", "bg-orange-700"},
		{"Final Edit", "The final polish before publication. Ensure all styling and factual checks are complete.", "bg-foreground"},
	}
	var cols []any
	for _, s := range stages {
		cols = append(cols, h.Div(h.Class("flex flex-col gap-2"),
			badge.Badge(h.Span(h.Aria("hidden", "true"), h.Class("size-1.5 rounded-full "+s.dot)), s.label),
			h.P(h.Class("text-xs leading-relaxed text-muted-foreground"), s.description),
		))
	}
	return seraSection(
		seraHeader([]any{
			seraBackLink("Editorial Dashboard", "", "flex items-center gap-1.5", "size-3.5"),
			seraHeading("Article Directory"),
		}, button.Button(h.Class("sm:ml-auto"), lucide.Plus(h.Data("icon", "inline-start")), "New Article")),
		h.Div(h.Class("container py-(--gap)"),
			card.Card(h.Class("py-24"),
				card.Content(h.Class("flex flex-col items-center gap-10"),
					empty.Empty(h.Class("min-h-96"),
						empty.Header(
							empty.Media(empty.MediaIcon, h.Class("size-14 rounded-full bg-muted/70 text-muted-foreground"), lucide.FileText(h.Class("size-5"))),
							empty.Title(h.Class("font-heading text-2xl tracking-normal normal-case"), "A Blank Canvas"),
							empty.Description("Your editorial directory is currently empty. Start building your publication's next issue by drafting the first piece."),
						),
						empty.Content(button.Button(lucide.Plus(h.Data("icon", "inline-start")), "Create first article")),
					),
					separator.Separator(h.Class("max-w-2xl")),
					h.Div(h.Class("grid w-full max-w-2xl grid-cols-1 gap-8 sm:grid-cols-3"), h.Group(cols...)),
				),
			),
		),
	)
}

// --- edit article ---------------------------------------------------------------

func seraEditArticle() h.Node {
	ghost := func(label string, icon func(...any) *h.Element, class string) h.Node {
		return button.Button(button.VariantGhost, button.SizeIconSm, h.Aria("label", label), h.If(class != "", h.Class(class)), icon())
	}
	sep := func(class string) h.Node {
		return buttongroup.Separator(h.Class(lib.Cn("mx-2", class, "data-vertical:h-4 data-vertical:self-center")))
	}
	issues := [][2]string{
		{"Spring Issue 2024", "spring-2024"}, {"Summer Issue 2024", "summer-2024"},
		{"Autumn Issue 2024", "autumn-2024"}, {"Winter Issue 2024", "winter-2024"},
	}
	var issueItems []any
	for _, is := range issues {
		issueItems = append(issueItems, selectui.Item(selectui.ItemValue(is[1]), is[0]))
	}
	milestones := []struct {
		name     string
		complete bool
	}{
		{"Outline & Commissioning", true}, {"First Draft Submitted", true}, {"Review & Revisions", false},
		{"Final Copy Edit", false}, {"Art Direction & Layout", false},
	}
	var ms []any
	for _, m := range milestones {
		ms = append(ms, field.Field(field.Horizontal,
			checkbox.Checkbox(checkbox.Checked(m.complete), h.Name(m.name), h.ID(m.name)),
			field.Label(h.For(m.name), m.name),
		))
	}
	return seraSection(
		seraHeader([]any{
			seraBackLink("Back to articles", "", "flex items-center gap-1.5", "size-3.5"),
			seraHeading("EDIT ARTICLE"),
		}, buttongroup.ButtonGroup(h.Class("gap-2 md:gap-4"),
			badge.Badge(h.Title("2 minutes ago"), "Autosaved"),
			buttongroup.ButtonGroup(h.Class("gap-2 md:gap-4"),
				button.Button(button.VariantLink, "Preview", lucide.ExternalLink(h.Data("icon", "inline-end"))),
				button.Button("Submit Draft"),
			),
		)),
		h.Div(h.Class("container py-(--gap)"),
			h.Div(h.Class("grid grid-cols-1 items-start gap-6 xl:grid-cols-[minmax(0,1fr)_300px]"),
				h.Section(h.Class("flex flex-col border border-border/70 bg-background"),
					h.Div(h.Class("flex border-b p-2"),
						buttongroup.ButtonGroup(
							dm.DropdownMenu(
								dm.Trigger(lib.AsChild, button.Button(button.VariantGhost, button.SizeSm, "Normal Text", lucide.ChevronDown(h.Data("icon", "inline-end")))),
								dm.Content(
									dm.Item(lucide.Type(), "Normal Text"),
									dm.Item(lucide.Heading1(), "Heading 1"),
									dm.Item(lucide.Heading2(), "Heading 2"),
									dm.Item(lucide.Heading3(), "Heading 3"),
									dm.Item(lucide.List(), "Bullet List"),
									dm.Item(lucide.ListOrdered(), "Numbered List"),
								),
							),
							sep(""),
							ghost("Bold", lucide.Bold, ""),
							ghost("Italic", lucide.Italic, ""),
							ghost("Underline", lucide.Underline, ""),
							ghost("Strikethrough", lucide.Strikethrough, "hidden md:flex"),
							ghost("Code", lucide.CodeXml, "hidden md:flex"),
							sep("hidden md:flex"),
							ghost("Align Left", lucide.AlignLeft, "hidden md:flex"),
							ghost("Align Center", lucide.AlignCenter, "hidden md:flex"),
							ghost("Align Right", lucide.AlignRight, "hidden md:flex"),
							sep("hidden md:flex"),
							ghost("Link", lucide.Link, "hidden md:flex"),
							ghost("Image", lucide.Image, "hidden md:flex"),
						),
						buttongroup.ButtonGroup(h.Class("ml-auto"),
							ghost("Undo", lucide.Undo, ""),
							ghost("Redo", lucide.Redo, ""),
						),
					),
					h.Div(h.Class("mx-auto flex max-w-2xl flex-1 flex-col gap-8 px-10 py-10 leading-loose md:px-14 lg:py-18"),
						h.H1(h.Class("font-heading text-4xl leading-12 font-medium tracking-wide uppercase"), "The Future of Sustainable Architecture"),
						h.P("As cities continue to expand at an unprecedented rate, the architectural paradigm is shifting from mere expansion to sustainable integration. The concrete jungles of the 20th century are making way for structures that breathe, adapt, and give back to their environments."),
						h.P("Historically, urban development has been a zero-sum game with nature."),
						h.H2(h.Class("font-heading text-2xl tracking-wide uppercase"), "The Living Building Challenge"),
						h.P("Sterling's latest project in downtown Seattle is a testament to this new philosophy. \"We are no longer designing static structures,\" Sterling explained during a recent site visit. \"We are engineering localized ecosystems.\""),
						h.P("The building features a facade of responsive biomaterials that adjust their porosity based on humidity and temperature, significantly reducing the need for artificial climate control. Rainwater is not merely channeled away but captured, filtered through a series of integrated rooftop wetlands, and reused within the building's greywater system."),
						h.P(h.Class("text-sm text-muted-foreground"), "This shift requires more than just innovative materials; it demands a fundamental change in how we value space. Check with engineering team for specific stats."),
					),
				),
				h.Aside(h.Class("grid grid-cols-12 gap-(--gap) xl:flex xl:flex-col"),
					card.Card(h.Class("col-span-full md:col-span-6 lg:col-span-4"),
						card.Header(card.Title("Article Details")),
						card.Content(
							field.Group(
								field.Field(
									field.Label("Issue"),
									selectui.Select(selectui.DefaultValue("summer-2024"),
										selectui.Trigger(selectui.Value()),
										selectui.Content(selectui.Group(h.Group(issueItems...))),
									),
								),
								field.Field(
									field.Label("Author"),
									input.Input(h.Value("Elena Rostova")),
								),
							),
						),
					),
					card.Card(h.Class("col-span-full md:col-span-6 lg:col-span-4"),
						card.Header(card.Title("Publication Flow")),
						card.Content(
							field.Group(
								field.Set(
									field.Legend("Required Milestones"),
									field.Field(h.Group(ms...)),
								),
								field.Field(
									field.Label("Add note for editor"),
									textarea.Textarea(h.Placeholder("This article needs to be revised for clarity and accuracy.")),
								),
							),
						),
					),
					card.Card(h.Class("col-span-full lg:col-span-4"),
						card.Header(card.Title("Word Count")),
						card.Content(seraProgress("sera-word-count", "1,402 / 2,000 words", 70)),
					),
				),
			),
		),
	)
}

// --- media library --------------------------------------------------------------

//go:embed data/sera/assets.json
var seraAssetsJSON []byte

type seraAsset struct {
	ID, Name, Date, Size, Type, Dimensions, Duration string
	UploadedBy                                       string `json:"uploadedBy"`
	Tags                                             []string
	UsedIn                                           []struct{ Title, Role string } `json:"usedIn"`
}

func seraAssets() []seraAsset {
	var a []seraAsset
	_ = json.Unmarshal(seraAssetsJSON, &a)
	return a
}

func seraAssetIcon(typ, class string) h.Node {
	switch typ {
	case "MP4":
		return lucide.Video(h.Class(class))
	case "PDF":
		return lucide.FileText(h.Class(class))
	}
	return lucide.Image(h.Class(class))
}

var seraTypeLabel = map[string]string{
	"JPEG": "Image / JPEG", "PNG": "Image / PNG", "WEBP": "Image / WEBP", "MP4": "Video / MP4", "PDF": "Document / PDF",
}

func seraMediaHeader() h.Node {
	return seraHeader([]any{
		seraBackLink("Asset management", "", "flex items-center gap-1.5", "size-3.5"),
		seraHeading("Media Library"),
	}, buttongroup.ButtonGroup(h.Class("gap-2 sm:ml-auto md:gap-4"),
		button.Button(button.VariantOutline, h.Class("bg-background hover:bg-background/80"), lucide.SlidersHorizontal(h.Data("icon", "inline-start")), "Filters"),
		button.Button(lucide.Upload(h.Data("icon", "inline-start")), "Upload Assets"),
	))
}

// seraPages is the media previews' pagination; text sets Previous/Next's
// text (nil: the default).
func seraPages(text *string) h.Node {
	prev, next := pagination.Previous(h.Href("#")), pagination.Next(h.Href("#"))
	if text != nil {
		prev.Children(pagination.Text(*text))
		next.Children(pagination.Text(*text))
	}
	return pagination.Pagination(pagination.Content(
		pagination.Item(prev),
		pagination.Item(pagination.Link(h.Href("#"), pagination.IsActive, "1")),
		pagination.Item(pagination.Link(h.Href("#"), "2")),
		pagination.Item(pagination.Link(h.Href("#"), "3")),
		pagination.Item(next),
	))
}

func seraMediaLibrary() h.Node {
	assets := seraAssets()
	var grid, details []any
	for i, a := range assets {
		selected := i == 0
		grid = append(grid, h.Button(h.Type("button"), h.Data("sera-asset", a.ID), h.Aria("pressed", boolString(selected)),
			h.Class("group flex flex-col gap-2.5 text-left outline-none focus-visible:[&>div:first-child]:ring-2 focus-visible:[&>div:first-child]:ring-ring"),
			h.Div(h.Class(lib.Cn("relative flex aspect-4/3 items-center justify-center bg-muted/60 ring-1 ring-border/70 transition-shadow ring-inset group-hover:ring-foreground/40",
				"group-aria-pressed:ring-2 group-aria-pressed:ring-foreground group-aria-pressed:group-hover:ring-foreground")),
				h.Div(h.Class("absolute top-2 left-2 hidden size-5 items-center justify-center bg-foreground text-background group-aria-pressed:flex"),
					lucide.Check(h.Class("size-3"))),
				badge.Badge(badge.VariantOutline, h.Class("absolute top-2 right-2 border bg-background px-2 py-1 text-[0.625rem]"), a.Type),
				h.If[any](a.Duration != "", badge.Badge(h.Class("absolute bottom-2 left-2 bg-foreground px-2 py-1 text-background"), a.Duration)),
				seraAssetIcon(a.Type, "size-7 text-muted-foreground/60"),
			),
			h.Div(h.Class("flex flex-col gap-0.5 px-0.5"),
				h.P(h.Class("line-clamp-1 text-sm font-medium"), a.Name),
				h.P(h.Class("text-[0.625rem] font-semibold tracking-wider text-muted-foreground uppercase"), a.Date+" · "+a.Size),
			),
		))
		details = append(details, seraAssetDetails(a, !selected))
	}
	return seraSection(seraMediaHeader(),
		h.Div(h.Class("container grid grid-cols-1 items-start gap-(--gap) py-(--gap) xl:grid-cols-[minmax(0,1fr)_320px]"),
			card.Card(
				card.Header(seraSearch("Search files, tags, or metadata...", h.Class("w-full"))),
				card.Content(h.Div(h.Class("grid grid-cols-2 gap-6 sm:grid-cols-3 lg:grid-cols-4"), h.Group(grid...))),
				card.Footer(h.Class("justify-center"), seraPages(nil)),
			),
			h.Div(h.Class("contents"), h.Group(details...)),
		),
	)
}

func seraAssetDetails(a seraAsset, hidden bool) h.Node {
	dt := func(text string) h.Node {
		return h.Dt(h.Class("text-[0.625rem] font-semibold tracking-widest text-muted-foreground uppercase"), text)
	}
	h3 := func(text string) h.Node {
		return h.H3(h.Class("text-[0.625rem] font-semibold tracking-widest text-muted-foreground uppercase"), text)
	}
	var tags, used []any
	for _, t := range a.Tags {
		tags = append(tags, badge.Badge(t))
	}
	for _, u := range a.UsedIn {
		used = append(used, item.Item(item.VariantOutline, item.Content(item.Title(u.Title), item.Description(u.Role))))
	}
	return card.Card(h.Data("sera-asset-details", a.ID), h.If(hidden, h.Bool("hidden")), h.Class("gap-0"),
		card.Content(h.Class("flex flex-col gap-6"),
			h.Div(h.Class("flex aspect-5/4 items-center justify-center bg-muted/60 text-muted-foreground/60 ring-1 ring-border/70 ring-inset"),
				seraAssetIcon(a.Type, "size-8")),
			h.H2(h.Class("line-clamp-2 font-heading text-xl tracking-wide"), a.Name),
			separator.Separator(),
			h.Dl(h.Class("flex flex-col gap-5 text-sm"),
				h.Div(h.Class("flex flex-col gap-1.5"), dt("Asset Type"), h.Dd(seraTypeLabel[a.Type])),
				h.Div(h.Class("grid grid-cols-2 gap-4"),
					h.Div(h.Class("flex flex-col gap-1.5"), dt("Dimensions"), h.Dd(a.Dimensions)),
					h.Div(h.Class("flex flex-col gap-1.5"), dt("File Size"), h.Dd(a.Size)),
				),
			),
			separator.Separator(),
			h.Div(h.Class("flex flex-col gap-3"),
				h.Div(h.Class("flex items-center justify-between"),
					h3("Tags"),
					button.Button(button.VariantGhost, button.SizeIconXs, h.Aria("label", "Add tag"), lucide.Plus()),
				),
				h.Div(h.Class("flex flex-wrap gap-x-4 gap-y-2"), h.Group(tags...)),
			),
			separator.Separator(),
			h.Div(h.Class("flex flex-col gap-3"),
				h3("Used In"),
				h.Div(h.Class("flex flex-col gap-2"), h.Group(used...)),
			),
		),
		card.Footer(h.Class("mt-6 border-t pt-6"),
			button.Button(h.Class("w-full"), lucide.Download(h.Data("icon", "inline-start")), "Download"),
		),
	)
}

// --- media library table --------------------------------------------------------

func seraMediaLibraryTable() h.Node {
	var rows []any
	for _, a := range seraAssets() {
		selected := a.ID == "1"
		rows = append(rows, table.Row(h.If(selected, h.Data("state", "selected")), h.Class("cursor-pointer"), h.Data("sera-select-row", ""),
			table.Cell(h.Class("pl-6"), checkbox.Checkbox(checkbox.Checked(selected), h.Aria("label", "Select "+a.Name))),
			table.Cell(
				h.Div(h.Class("relative flex aspect-4/3 w-16 items-center justify-center bg-muted/60 ring-1 ring-border/70 ring-inset"),
					h.If[any](a.Duration != "", h.Span(h.Class("absolute right-1 bottom-1 bg-foreground/90 px-1 text-[0.5rem] font-semibold tracking-wider text-background"), a.Duration)),
					seraAssetIcon(a.Type, "size-4 text-muted-foreground/60"),
				),
			),
			table.Cell(h.Class("text-sm font-medium text-foreground"), a.Name),
			table.Cell(badge.Badge(badge.VariantOutline, h.Class("border px-2 py-0.5 text-[0.625rem]"), a.Type)),
			table.Cell(h.Class("text-sm"), a.Dimensions),
			table.Cell(h.Class("text-sm"), a.Size),
			table.Cell(a.UploadedBy),
			table.Cell(h.Class("text-xs font-semibold tracking-wider text-muted-foreground uppercase"), a.Date),
			table.Cell(h.Class("pr-6 text-right"),
				dm.DropdownMenu(
					dm.Trigger(lib.AsChild, button.Button(button.VariantGhost, button.SizeIconXs, h.Aria("label", "Open actions for "+a.Name), lucide.EllipsisVertical())),
					dm.Content(dm.AlignEnd,
						dm.Item("Preview"),
						dm.Item("Download"),
						dm.Item(dm.VariantDestructive, "Delete"),
					),
				),
			),
		))
	}
	empty := ""
	return seraSection(seraMediaHeader(),
		h.Div(h.Class("container grid grid-cols-1 items-start gap-(--gap) py-(--gap) xl:grid-cols-[minmax(0,1fr)_320px]"),
			card.Card(
				card.Header(seraSearch("Search files, tags, or metadata...", h.Class("w-full"))),
				card.Content(h.Class("px-0 py-0"),
					table.Table(
						table.Header(table.Row(
							table.Head(h.Class("w-10 pl-6"), h.Aria("label", "Select")),
							table.Head(h.Class("w-20"), h.Aria("label", "Preview")),
							table.Head("Filename"), table.Head("Type"), table.Head("Dimensions"), table.Head("Size"),
							table.Head("Uploaded By"), table.Head("Date"),
							table.Head(h.Class("w-10 pr-6"), h.Aria("label", "Actions")),
						)),
						table.Body(h.Group(rows...)),
					),
				),
				card.Footer(h.Class("justify-center py-4"), seraPages(&empty)),
			),
			seraFilterLibrary(),
		),
	)
}

func seraFilterLibrary() h.Node {
	fileTypes := []struct {
		id, label string
		checked   bool
	}{
		{"images", "Images (JPEG, PNG, WEBP)", true}, {"video", "Video (MP4, MOV)", false},
		{"documents", "Documents (PDF)", false}, {"audio", "Audio (MP3, WAV)", false},
	}
	var types []any
	for _, t := range fileTypes {
		types = append(types, field.Field(field.Horizontal,
			checkbox.Checkbox(h.ID(t.id), checkbox.Checked(t.checked)),
			field.Label(h.For(t.id), t.label),
		))
	}
	var dates []any
	for _, o := range [][2]string{{"any", "Any time"}, {"24h", "Past 24 hours"}, {"week", "Past week"}, {"month", "Past month"}} {
		dates = append(dates, field.Field(field.Horizontal,
			rg.Item(rg.Value(o[0]), h.ID("date-"+o[0])),
			field.Label(h.For("date-"+o[0]), o[1]),
		))
	}
	tags := []string{"architecture", "brutalism", "ceramics", "design-week", "editorial", "exterior", "film", "food",
		"furniture", "interior", "kyoto", "minimalism", "print", "sustainability", "summer-issue", "video"}
	var tagItems []any
	for _, t := range tags {
		tagItems = append(tagItems, cb.Item(cb.ItemValue(t), t))
	}
	now := lib.Now()
	from := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	to := from.AddDate(0, 0, 21)
	day := func(t time.Time) string { return fmt.Sprintf("%s %02d, %d", t.Format("Jan"), t.Day(), t.Year()) }
	anchor := cb.UseAnchor()
	return card.Card(
		card.Header(h.Class("border-b"), card.Title("Filter Library")),
		card.Content(
			field.Group(
				field.Set(field.Legend("File Type"), field.Field(h.Group(types...))),
				field.Set(field.Legend("Date Uploaded"), rg.RadioGroup(rg.Value("any"), h.Group(dates...))),
				field.Field(
					field.Label(h.For("custom-range"), "Custom Range"),
					popover.Popover(
						popover.Trigger(lib.AsChild, button.Button(button.VariantOutline, h.ID("custom-range"), h.Class("w-full"),
							lucide.Calendar(h.Data("icon", "inline-start")),
							day(from)+" – "+day(to))),
						popover.Content(h.Class("w-auto p-0"), popover.AlignEnd,
							calendar.Calendar(calendar.ModeRange, calendar.DefaultMonth(from), calendar.SelectedRange(from, to), calendar.NumberOfMonths(2)),
						),
					),
				),
				field.Set(field.Legend("File Size"),
					h.Div(h.Class("flex flex-col gap-3"),
						h.Div(h.Class("flex items-center justify-between text-xs font-medium tracking-wider text-muted-foreground uppercase"),
							h.Span("0 MB"), h.Span("500+ MB")),
						slider.Slider(slider.DefaultValue(0, 60), slider.Max(100), slider.Step(1)),
						h.Div(h.Class("flex items-center justify-between text-xs font-medium"),
							h.Span("Min: 0 MB"), h.Span("Max: 300 MB")),
					),
				),
				field.Set(field.Legend("Tags"),
					field.Field(
						cb.Combobox(cb.Multiple, cb.AutoHighlight, cb.DefaultValue("architecture"), cb.DefaultValue("brutalism"),
							cb.Chips(anchor, cb.Value(cb.Chip(), cb.ChipsInput(h.Placeholder("Filter by tag...")))),
							cb.Content(anchor, cb.Empty("No tags found."), cb.List(h.Group(tagItems...))),
						),
					),
				),
			),
		),
		card.Footer(h.Class("flex flex-col gap-2 border-t"),
			button.Button(h.Class("w-full"), "Apply Filters"),
			button.Button(button.VariantGhost, h.Class("w-full"), "Reset"),
		),
	)
}
