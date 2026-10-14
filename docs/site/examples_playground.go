package site

import (
	_ "embed"
	"net/http"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/lib"
	alertdialog "shadcn.go/components/ui/alert-dialog"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/command"
	"shadcn.go/components/ui/dialog"
	dm "shadcn.go/components/ui/dropdown-menu"
	hovercard "shadcn.go/components/ui/hover-card"
	"shadcn.go/components/ui/input"
	"shadcn.go/components/ui/label"
	"shadcn.go/components/ui/popover"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/components/ui/separator"
	"shadcn.go/components/ui/slider"
	switchui "shadcn.go/components/ui/switch"
	"shadcn.go/components/ui/tabs"
	"shadcn.go/components/ui/textarea"
	h "shadcn.go/libhtml"
)

// The playground example: app/(app)/examples/playground/{page.tsx,
// components/*, data/*}. examples-playground.js keeps the React state: the
// picked preset and model, the slider readouts and the preset actions'
// dialogs.

//go:embed static/examples-playground.js
var playgroundScript string

func init() { runtime.Module("examples-playground", playgroundScript) }

type playgroundModel struct{ name, description, strengths, typ string }

// playgroundModels is data/models.ts (the ids are not rendered).
var playgroundModels = []playgroundModel{
	{"text-davinci-003", "Most capable GPT-3 model. Can do any task the other models can do, often with higher quality, longer output and better instruction-following. Also supports inserting completions within text.", "Complex intent, cause and effect, creative generation, search, summarization for audience", "GPT-3"},
	{"text-curie-001", "Very capable, but faster and lower cost than Davinci.", "Language translation, complex classification, sentiment, summarization", "GPT-3"},
	{"text-babbage-001", "Capable of straightforward tasks, very fast, and lower cost.", "Moderate classification, semantic search", "GPT-3"},
	{"text-ada-001", "Capable of very simple tasks, usually the fastest model in the GPT-3 series, and lowest cost.", "Parsing text, simple classification, address correction, keywords", "GPT-3"},
	{"code-davinci-002", "Most capable Codex model. Particularly good at translating natural language to code. In addition to completing code, also supports inserting completions within code.", "", "Codex"},
	{"code-cushman-001", "Almost as capable as Davinci Codex, but slightly faster. This speed advantage may make it preferable for real-time applications.", "Real-time application where low-latency is preferable", "Codex"},
}

var playgroundModelTypes = []string{"GPT-3", "Codex"}

// playgroundPresets is data/presets.ts.
var playgroundPresets = []string{
	"Grammatical Standard English",
	"Summarize for a 2nd grader",
	"Text to command",
	"Q&A",
	"English to other languages",
	"Parse unstructured data",
	"Classification",
	"Natural language to Python",
	"Explain code",
	"Chat",
}

// ExamplePlayground serves /examples/playground.
func (s *Site) ExamplePlayground(r *http.Request) (h.Node, error) {
	img := func(theme, class string) h.Node {
		u := "%2Fexamples%2Fplayground-" + theme + ".png"
		return h.Img(h.Alt("Playground"), h.Loading("lazy"), h.Width("1280"), h.Height("916"),
			h.Attr("decoding", "async"), h.Data("nimg", "1"), h.Class(class), h.Style("color:transparent"),
			h.UnsafeAttr("srcset", "/_next/image?url="+u+"&w=1920&q=75 1x, /_next/image?url="+u+"&w=3840&q=75 2x"),
			h.UnsafeAttr("src", "/_next/image?url="+u+"&w=3840&q=75"))
	}
	return s.examplesLayoutMeta(r, Meta{Title: "Playground", PlainTitle: true, Description: "The OpenAI Playground built using the components.", Canonical: "/examples/playground"},
		h.Div(h.Class("md:hidden"),
			img("light", "block dark:hidden"),
			img("dark", "hidden dark:block"),
		),
		h.Div(h.Class("hidden flex-1 flex-col md:flex"),
			h.Div(h.Class("container flex flex-col items-start justify-between gap-2 py-4 sm:flex-row sm:items-center sm:gap-0 md:h-16"),
				h.H2(h.Class("pl-0.5 text-lg font-semibold"), "Playground"),
				h.Div(h.Class("ml-auto flex w-full gap-2 sm:justify-end"),
					playgroundPresetSelector(),
					playgroundPresetSave(),
					h.Div(h.Class("hidden gap-2 md:flex"),
						playgroundCodeViewer(),
						playgroundPresetShare(),
					),
					playgroundPresetActions(),
				),
			),
			separator.Separator(),
			tabs.Tabs(tabs.Value("complete"), h.Class("flex flex-1 flex-col"),
				h.Div(h.Class("container flex flex-1 flex-col py-6"),
					h.Div(h.Class("grid flex-1 items-stretch gap-6 md:grid-cols-[1fr_200px]"),
						h.Div(h.Class("hidden flex-col gap-6 sm:flex md:order-2"),
							h.Div(h.Class("grid gap-3"),
								hovercard.HoverCard(hovercard.OpenDelay(200),
									hovercard.Trigger(lib.AsChild,
										h.Span(h.Class("text-sm leading-none font-medium peer-disabled:cursor-not-allowed peer-disabled:opacity-70"), "Mode"),
									),
									hovercard.Content(h.Class("w-[320px] text-sm"), hovercard.SideLeft,
										"Choose the interface that best suits your task. You can provide: a simple prompt to complete, starting and ending text to insert a completion within, or some text with instructions to edit it.",
									),
								),
								tabs.List(h.Class("grid w-full grid-cols-3"),
									tabs.Trigger(tabs.Value("complete"), h.Span(h.Class("sr-only"), "Complete"), modeIconComplete()),
									tabs.Trigger(tabs.Value("insert"), h.Span(h.Class("sr-only"), "Insert"), modeIconInsert()),
									tabs.Trigger(tabs.Value("edit"), h.Span(h.Class("sr-only"), "Edit"), modeIconEdit()),
								),
							),
							playgroundModelSelector(),
							playgroundSlider("temperature", "Temperature", "0.56", slider.Slider(slider.Max(1), slider.DefaultValue(0.56), slider.Step(0.1)),
								"Controls randomness: lowering results in less random completions. As the temperature approaches zero, the model will become deterministic and repetitive."),
							playgroundSlider("maxlength", "Maximum Length", "256", slider.Slider(slider.Max(4000), slider.DefaultValue(256), slider.Step(10)),
								"The maximum number of tokens to generate. Requests can use up to 2,048 or 4,000 tokens, shared between prompt and completion. The exact limit varies by model."),
							playgroundSlider("top-p", "Top P", "0.9", slider.Slider(slider.Max(1), slider.DefaultValue(0.9), slider.Step(0.1)),
								"Control diversity via nucleus sampling: 0.5 means half of all likelihood-weighted options are considered."),
						),
						h.Div(h.Class("flex flex-1 flex-col *:data-[slot=tab-content]:flex-1 md:order-1"),
							tabs.Content(tabs.Value("complete"), h.Class("mt-0 border-0 p-0"),
								h.Div(h.Class("flex h-full flex-col gap-4"),
									textarea.Textarea(h.Placeholder("Write a tagline for an ice cream shop"), h.Class("min-h-[400px] flex-1 p-4 md:min-h-[700px] lg:min-h-[700px]")),
									playgroundSubmitRow(),
								),
							),
							tabs.Content(tabs.Value("insert"), h.Class("mt-0 flex flex-col gap-4 border-0 p-0"),
								h.Div(h.Class("grid h-full grid-rows-2 gap-6 lg:grid-cols-2 lg:grid-rows-1"),
									textarea.Textarea(h.Placeholder("We're writing to [inset]. Congrats from OpenAI!"), h.Class("h-full min-h-[300px] p-4 lg:min-h-[700px] xl:min-h-[700px]")),
									h.Div(h.Class("rounded-md border bg-muted")),
								),
								playgroundSubmitRow(),
							),
							tabs.Content(tabs.Value("edit"), h.Class("mt-0 flex flex-col gap-4 border-0 p-0"),
								h.Div(h.Class("grid h-full gap-6 lg:grid-cols-2"),
									h.Div(h.Class("flex flex-col gap-4"),
										h.Div(h.Class("flex flex-1 flex-col gap-2"),
											label.Label(h.For("input"), h.Class("sr-only"), "Input"),
											textarea.Textarea(h.ID("input"), h.Placeholder("We is going to the market."), h.Class("flex-1 p-4 lg:min-h-[580px]")),
										),
										h.Div(h.Class("flex flex-col gap-2"),
											label.Label(h.For("instructions"), "Instructions"),
											textarea.Textarea(h.ID("instructions"), h.Placeholder("Fix the grammar.")),
										),
									),
									h.Div(h.Class("min-h-[400px] rounded-md border bg-muted lg:min-h-[700px]")),
								),
								playgroundSubmitRow(),
							),
						),
					),
				),
			),
		),
	), nil
}

func playgroundSubmitRow() h.Node {
	return h.Div(h.Class("flex items-center gap-2"),
		button.Button("Submit"),
		button.Button(button.VariantSecondary, h.Span(h.Class("sr-only"), "Show history"), lucide.RotateCcw()),
	)
}

// playgroundSlider is components/{temperature,maxlength,top-p}-selector.tsx.
func playgroundSlider(id, name, value string, sl *slider.Builder, help string) h.Node {
	return h.Div(h.Class("grid gap-2 pt-2"),
		hovercard.HoverCard(hovercard.OpenDelay(200),
			hovercard.Trigger(lib.AsChild,
				h.Div(h.Class("grid gap-4"),
					h.Div(h.Class("flex items-center justify-between"),
						label.Label(h.For(id), name),
						h.Span(h.Data("playground-slider-value", id),
							h.Class("w-12 rounded-md border border-transparent px-2 py-0.5 text-right text-sm text-muted-foreground hover:border-border"),
							value),
					),
					sl.Children(h.ID(id), h.Aria("label", name), h.Data("playground-slider", id)),
				),
			),
			hovercard.Content(hovercard.AlignStart, h.Class("w-[260px] text-sm"), hovercard.SideLeft, help),
		),
	)
}

// playgroundPresetSelector is components/preset-selector.tsx.
func playgroundPresetSelector() h.Node {
	return popover.Popover(
		popover.Trigger(lib.AsChild, button.Button(button.VariantOutline,
			h.Role("combobox"), h.Aria("label", "Load a preset..."), h.Aria("expanded", "false"),
			h.Class("flex-1 justify-between md:max-w-[200px] lg:max-w-[300px]"),
			h.Data("playground-picker", "Load a preset..."),
			"Load a preset...",
			lucide.ChevronsUpDown(h.Class("opacity-50")),
		)),
		popover.Content(h.Class("w-[300px] p-0"),
			command.Command(
				command.Input(h.Placeholder("Search presets...")),
				command.List(
					command.Empty("No presets found."),
					command.Group(command.Heading("Examples"),
						h.Map(playgroundPresets, func(p string) h.Node {
							return command.Item(h.Data("playground-pick", p), p, lucide.Check(h.Class("ml-auto opacity-0")))
						}),
					),
					command.Separator(),
					command.Group(command.Item("More examples")),
				),
			),
		),
	)
}

// playgroundModelSelector is components/model-selector.tsx. Its peek card
// (a HoverCard whose trigger is an empty element in the list) never opens,
// so it is not rendered.
func playgroundModelSelector() h.Node {
	var groups []any
	for _, t := range playgroundModelTypes {
		var items []any
		for i, m := range playgroundModels {
			if m.typ != t {
				continue
			}
			check := "ml-auto opacity-0"
			if i == 0 {
				check = "ml-auto opacity-100"
			}
			items = append(items, command.Item(h.Data("playground-pick", m.name),
				h.Class("data-[selected=true]:bg-primary data-[selected=true]:text-primary-foreground"),
				m.name, lucide.Check(h.Class(check))))
		}
		groups = append(groups, command.Group(command.Heading(t), h.Group(items...)))
	}
	return h.Div(h.Class("grid gap-3"),
		hovercard.HoverCard(hovercard.OpenDelay(200),
			hovercard.Trigger(lib.AsChild, label.Label(h.For("model"), "Model")),
			hovercard.Content(hovercard.AlignStart, h.Class("w-[260px] text-sm"), hovercard.SideLeft,
				"The model which will generate the completion. Some models are suitable for natural language tasks, others specialize in code. Learn more."),
		),
		popover.Popover(
			popover.Trigger(lib.AsChild, button.Button(button.VariantOutline,
				h.Role("combobox"), h.Aria("expanded", "false"), h.Aria("label", "Select a model"),
				h.Class("w-full justify-between"),
				h.Data("playground-picker", ""),
				playgroundModels[0].name,
				lucide.ChevronsUpDown(h.Class("text-muted-foreground")),
			)),
			popover.Content(popover.AlignEnd, h.Class("w-[250px] p-0"),
				command.Command(command.Loop(true),
					command.List(h.Class("h-(--cmdk-list-height) max-h-[400px]"),
						command.Input(h.Placeholder("Search Models...")),
						command.Empty("No Models found."),
						h.Group(groups...),
					),
				),
			),
		),
	)
}

// playgroundPresetSave is components/preset-save.tsx.
func playgroundPresetSave() h.Node {
	return dialog.Dialog(
		dialog.Trigger(lib.AsChild, button.Button(button.VariantSecondary, "Save")),
		dialog.Content(
			dialog.Header(
				dialog.Title("Save preset"),
				dialog.Description("This will save the current playground state as a preset which you can access later or share with others."),
			),
			h.Div(h.Class("grid gap-6 py-4"),
				h.Div(h.Class("grid gap-3"),
					label.Label(h.For("name"), "Name"),
					input.Input(h.ID("name"), h.AutoFocus()),
				),
				h.Div(h.Class("grid gap-3"),
					label.Label(h.For("description"), "Description"),
					textarea.Textarea(h.ID("description")),
				),
			),
			dialog.Footer(button.Button(h.Type("submit"), "Save")),
		),
	)
}

// playgroundCodeViewer is components/code-viewer.tsx.
func playgroundCodeViewer() h.Node {
	sky := func(t string) h.Node { return h.Span(h.Class("text-sky-300"), t) }
	green := func(t string) h.Node { return h.Span(h.Class("text-green-300"), t) }
	amber := func(t string) h.Node { return h.Span(h.Class("text-amber-300"), t) }
	return dialog.Dialog(
		dialog.Trigger(lib.AsChild, button.Button(button.VariantSecondary, "View code")),
		dialog.Content(
			dialog.Header(
				dialog.Title("View code"),
				dialog.Description("You can use the following code to start integrating your current prompt and settings into your application."),
			),
			h.Div(h.Class("grid gap-4"),
				h.Div(h.Class("rounded-md bg-black p-6"),
					h.Pre(h.Code(h.Class("grid gap-1 text-sm text-white [&_span]:h-4"),
						h.Span(sky("import"), " os"),
						h.Span(sky("import"), " openai"),
						h.Span(),
						h.Span("openai.api_key = os.getenv(", green(`"OPENAI_API_KEY"`), ")"),
						h.Span(),
						h.Span("response = openai.Completion.create("),
						h.Span(" ", "model=", green(`"davinci"`), ","),
						h.Span(" ", "prompt=", amber(`""`), ","),
						h.Span(" ", "temperature=", amber("0.9"), ","),
						h.Span(" ", "max_tokens=", amber("5"), ","),
						h.Span(" ", "top_p=", amber("1"), ","),
						h.Span(" ", "frequency_penalty=", amber("0"), ","),
						h.Span(" ", "presence_penalty=", green("0"), ","),
						h.Span(")"),
					)),
				),
				h.Div(h.P(h.Class("text-sm text-muted-foreground"),
					"Your API Key can be found here. You should use environment variables or a secret management tool to expose your key to your applications.")),
			),
		),
	)
}

// playgroundPresetShare is components/preset-share.tsx.
func playgroundPresetShare() h.Node {
	return popover.Popover(
		popover.Trigger(lib.AsChild, button.Button(button.VariantSecondary, "Share")),
		popover.Content(popover.AlignEnd, h.Class("flex w-[520px] flex-col gap-4"),
			h.Div(h.Class("flex flex-col gap-1 text-center sm:text-left"),
				h.H3(h.Class("text-lg font-semibold"), "Share preset"),
				h.P(h.Class("text-sm text-muted-foreground"), "Anyone who has this link and an OpenAI account will be able to view this."),
			),
			h.Div(h.Class("relative flex-1"),
				label.Label(h.For("link"), h.Class("sr-only"), "Link"),
				input.Input(h.ID("link"), h.Value("https://platform.openai.com/playground/p/7bbKYQvsVkNmVb8NGcdUOLae?model=text-davinci-003"), h.ReadOnly(), h.Class("h-9 pr-10")),
				button.Button(h.Type("submit"), button.SizeIcon, button.VariantGhost, h.Class("absolute top-1 right-1 size-7"),
					h.Span(h.Class("sr-only"), "Copy"), lucide.Copy(h.Class("size-3.5"))),
			),
		),
	)
}

// playgroundPresetActions is components/preset-actions.tsx: the menu items
// open the dialogs (examples-playground.js).
func playgroundPresetActions() h.Node {
	return h.Group(
		dm.DropdownMenu(
			dm.Trigger(lib.AsChild, button.Button(button.VariantSecondary, button.SizeIcon,
				h.Span(h.Class("sr-only"), "Actions"), lucide.Ellipsis())),
			dm.Content(dm.AlignEnd,
				dm.Item(h.Data("playground-open", "filter"), "Content filter preferences"),
				dm.Separator(),
				dm.Item(h.Data("playground-open", "delete"), h.Class("text-red-600"), "Delete preset"),
			),
		),
		h.Div(h.Class("contents"), h.Data("playground-dialog", "filter"),
			dialog.Dialog(
				dialog.Content(
					dialog.Header(
						dialog.Title("Content filter preferences"),
						dialog.Description("The content filter flags text that may violate our content policy. It's powered by our moderation endpoint which is free to use to moderate your OpenAI API traffic. Learn more."),
					),
					h.Div(h.Class("py-6"),
						h.H4(h.Class("text-sm text-muted-foreground"), "Playground Warnings"),
						h.Div(h.Class("flex items-start justify-between gap-4 pt-3"),
							switchui.Switch(switchui.Name("show"), h.ID("show"), switchui.Checked(true)),
							label.Label(h.Class("grid gap-1 font-normal"), h.For("show"),
								h.Span(h.Class("font-semibold"), "Show a warning when content is flagged"),
								h.Span(h.Class("text-sm text-muted-foreground"), "A warning will be shown when sexual, hateful, violent or self-harm content is detected."),
							),
						),
					),
					dialog.Footer(dialog.Close(lib.AsChild, button.Button(button.VariantSecondary, "Close"))),
				),
			),
		),
		h.Div(h.Class("contents"), h.Data("playground-dialog", "delete"),
			alertdialog.AlertDialog(
				alertdialog.Content(
					alertdialog.Header(
						alertdialog.Title("Are you absolutely sure?"),
						alertdialog.Description("This action cannot be undone. This preset will no longer be accessible by you or others you've shared it with."),
					),
					alertdialog.Footer(
						alertdialog.Cancel("Cancel"),
						button.Button(button.VariantDestructive, h.Data("playground-delete", ""), "Delete"),
					),
				),
			),
		),
	)
}

func modeIcon(class string, children ...any) h.Node {
	attrs := []any{h.Attr("xmlns", "http://www.w3.org/2000/svg"), h.Attr("viewBox", "0 0 20 20"), h.Attr("fill", "none")}
	if class != "" {
		attrs = append(attrs, h.Class(class))
	}
	return h.El("svg", append(attrs, children...)...)
}

func modeRect(x, y, w string) h.Node {
	return h.El("rect", h.Attr("x", x), h.Attr("y", y), h.Attr("width", w), h.Attr("height", "2"), h.Attr("rx", "1"), h.Attr("fill", "currentColor"))
}

func modePath(d string, evenodd bool) h.Node {
	attrs := []any{}
	if evenodd {
		attrs = append(attrs, h.Attr("fill-rule", "evenodd"), h.Attr("clip-rule", "evenodd"))
	}
	return h.El("path", append(attrs, h.Attr("d", d), h.Attr("fill", "currentColor"))...)
}

func modeIconComplete() h.Node {
	return modeIcon("",
		modeRect("4", "3", "12"), modeRect("4", "7", "12"), modeRect("4", "11", "3"), modeRect("4", "15", "3"),
		modeRect("8.5", "11", "3"), modeRect("8.5", "15", "3"), modeRect("13", "11", "3"))
}

func modeIconInsert() h.Node {
	return modeIcon("h-5 w-5",
		modePath("M14.491 7.769a.888.888 0 0 1 .287.648.888.888 0 0 1-.287.648l-3.916 3.667a1.013 1.013 0 0 1-.692.268c-.26 0-.509-.097-.692-.268L5.275 9.065A.886.886 0 0 1 5 8.42a.889.889 0 0 1 .287-.64c.181-.17.427-.267.683-.269.257-.002.504.09.69.258L8.903 9.87V3.917c0-.243.103-.477.287-.649.183-.171.432-.268.692-.268.26 0 .509.097.692.268a.888.888 0 0 1 .287.649V9.87l2.245-2.102c.183-.172.432-.269.692-.269.26 0 .508.097.692.269Z", true),
		modeRect("4", "15", "3"), modeRect("8.5", "15", "3"), modeRect("13", "15", "3"))
}

func modeIconEdit() h.Node {
	return modeIcon("h-5 w-5",
		modeRect("4", "3", "12"), modeRect("4", "7", "12"), modeRect("4", "11", "3"), modeRect("4", "15", "4"),
		modeRect("8.5", "11", "3"),
		modePath("M17.154 11.346a1.182 1.182 0 0 0-1.671 0L11 15.829V17.5h1.671l4.483-4.483a1.182 1.182 0 0 0 0-1.671Z", false))
}
