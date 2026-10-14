// Package questionnaire ports the docs' questionnaire examples
// (ui/apps/v4/examples/radix/questionnaire-*.tsx), rendered in radix-nova.
// examples.js carries their React state: the submit toasts, the
// controlled and conditional flows, the custom progress, the shortcut
// selector and the schema validation.
package questionnaire

import (
	_ "embed"
	"net/http"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/dialog"
	nativeselect "shadcn.go/components/ui/native-select"
	q "shadcn.go/components/ui/questionnaire"
	"shadcn.go/components/ui/runtime"
	"shadcn.go/components/ui/sonner"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

//go:embed examples.js
var script string

func init() {
	runtime.Module("example-questionnaire", script)
	docs := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "questionnaire", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		}), Extra: toaster}
	}
	registry.Register(
		docs("questionnaire-animated", Animated),
		docs("questionnaire-card", Card),
		docs("questionnaire-conditional", Conditional),
		docs("questionnaire-controlled", Controlled),
		docs("questionnaire-demo", Demo),
		docs("questionnaire-dialog", Dialog),
		docs("questionnaire-freeform", Freeform),
		docs("questionnaire-multiple", Multiple),
		docs("questionnaire-navigation-state", NavigationState),
		docs("questionnaire-progress", Progress),
		docs("questionnaire-resume", Resume),
		docs("questionnaire-shortcuts", Shortcuts),
		docs("questionnaire-skip", Skip),
		docs("questionnaire-validation", Validation),
	)
}

// toaster is the <Toaster /> the docs mount in their layout.
func toaster(*http.Request) h.Node { return sonner.Toaster(sonner.ThemeLight) }

// example marks a questionnaire for examples.js.
func example(name string) h.Prop { return h.Data("example-questionnaire", name) }

func def(name string, required bool, values ...string) q.ItemDef {
	d := q.ItemDef{Name: name, Required: required}
	for _, v := range values {
		d.Choices = append(d.Choices, q.ChoiceDef{Value: v})
	}
	return d
}

func nodes(ns []h.Node) []any {
	out := make([]any, len(ns))
	for i, n := range ns {
		out[i] = n
	}
	return out
}

func choice(value string, children ...any) h.Node {
	return q.Choice(append([]any{q.Value(value)}, children...)...)
}

// item is a required single-choice item with a title and an optional
// description.
func item(name, title, description string, choices ...h.Node) h.Node {
	kids := []any{q.Name(name), q.Required(true), q.Title(title)}
	if description != "" {
		kids = append(kids, q.Description(description))
	}
	kids = append(kids, q.Choices(nodes(choices)...), q.Error())
	return q.Item(kids...)
}

func Demo() h.Node {
	return q.Questionnaire(example("demo"), h.Class("mx-auto max-w-md"), q.DefaultItem("direction"), q.ShortcutsLetters,
		q.Items(
			def("direction", true, "tool-calls", "approvals", "handoffs"),
			def("signals", false, "progress", "decisions", "risks", "next-step"),
			def("timing", true, "now", "next-cycle", "backlog"),
		),
		q.Progress(),
		q.Item(q.Name("direction"), q.Required(true),
			q.Title("What should the agent build next?"),
			q.Description("Choose a direction or describe another task."),
			q.Choices(
				choice("tool-calls",
					h.Span(h.Class("font-medium"), "Tool call timeline"),
					h.Span(h.Class("text-muted-foreground"), "Show what the agent ran and what came back."),
				),
				choice("approvals",
					h.Span(h.Class("font-medium"), "Approval checkpoints"),
					h.Span(h.Class("text-muted-foreground"), "Ask before sensitive or destructive actions."),
				),
				choice("handoffs",
					h.Span(h.Class("font-medium"), "Sub-agent handoffs"),
					h.Span(h.Class("text-muted-foreground"), "Make delegated work and results easier to follow."),
				),
				q.Input(h.Aria("label", "Another agent feature"), h.Placeholder("Describe another feature…")),
			),
			q.Error(),
		),
		q.Item(q.Name("signals"), q.Multiple(true),
			q.Title("What should every progress update include?"),
			q.Description("Select all that apply, or skip this question."),
			q.Choices(
				choice("progress", "Progress"),
				choice("decisions", "Decisions"),
				choice("risks", "Risks"),
				choice("next-step", "Next step"),
			),
			q.Error(),
		),
		item("timing", "When should work begin?", "Choose when the agent should begin the work.",
			choice("now", "Start now"),
			choice("next-cycle", "Next development cycle"),
			choice("backlog", "Add it to the backlog"),
		),
		q.Actions(q.Previous(), q.Skip(), q.Next("Next"), q.Submit("Save plan")),
	)
}

const animatedItemClass = "data-active:animate-in data-active:fade-in-0 data-active:slide-in-from-bottom-2 data-active:duration-300 motion-reduce:animate-none"

func Animated() h.Node {
	animated := func(name, title, description string, choices ...h.Node) h.Node {
		return q.Item(h.Class(animatedItemClass), q.Name(name), q.Required(true),
			q.Title(title), q.Description(description), q.Choices(nodes(choices)...), q.Error())
	}
	return q.Questionnaire(example("animated"), h.Class("mx-auto max-w-md"), q.DefaultItem("task"),
		q.Items(def("task", true), def("review", true), def("delivery", true)),
		q.Progress(),
		animated("task", "What should the agent do?", "Choose the task for this run.",
			choice("implement", "Implement the requested change"),
			choice("debug", "Debug the current behavior"),
			choice("review", "Review the implementation"),
		),
		animated("review", "How should the work be reviewed?", "Select the verification depth.",
			choice("targeted", "Targeted checks"),
			choice("complete", "Complete test suite"),
			choice("manual", "Tests and manual QA"),
		),
		animated("delivery", "How should the result be delivered?", "Choose the final handoff format.",
			choice("summary", "Concise summary"),
			choice("diff", "Summary and changed files"),
			choice("handoff", "Detailed review handoff"),
		),
		q.Actions(q.Previous(), q.Next("Next"), q.Submit("Save workflow")),
	)
}

func Card() h.Node {
	cardItem := func(name, titleID, title, description string, choices ...h.Node) h.Node {
		return q.Item(h.Aria("labelledby", titleID), q.Name(name), q.Required(true),
			card.Header(
				q.Title(h.ID(titleID), q.Render(card.Title()), title),
				q.Description(q.Render(card.Description()), description),
				card.Action(q.Progress()),
			),
			card.Content(q.Choices(nodes(choices)...), q.Error()),
		)
	}
	return q.Questionnaire(example("card"), h.Class("mx-auto max-w-md"), q.DefaultItem("task"), q.ShortcutsNumbers,
		q.Items(def("task", true, "fix", "refactor", "docs"), def("output", true, "summary", "files", "review")),
		card.Card(
			cardItem("task", "questionnaire-card-task-title", "What should the agent work on?", "Choose the task that should be handled next.",
				choice("fix", "Fix the failing tests"),
				choice("refactor", "Refactor the data layer"),
				choice("docs", "Update the integration guide"),
			),
			cardItem("output", "questionnaire-card-output-title", "What should the final handoff include?", "Pick the level of detail needed for review.",
				choice("summary", "Summary only"),
				choice("files", "Summary and changed files"),
				choice("review", "Full review handoff"),
			),
			card.Footer(q.Actions(h.Class("w-full"), q.Previous(), q.Next("Next"), q.Submit("Create task"))),
		),
	)
}

func Conditional() h.Node {
	return q.Questionnaire(example("conditional"), h.Class("mx-auto max-w-md"), q.DefaultItem("runtime"),
		q.Items(def("runtime", true), q.ItemDef{Name: "environment", Required: true, Disabled: true}, def("approval", true)),
		q.Progress(),
		item("runtime", "Where should the agent run?", "Cloud runs add an environment question to this flow.",
			q.Choice(q.Value("local"), q.DefaultChecked(true), "Local workspace"),
			choice("cloud", "Cloud workspace"),
		),
		q.Item(q.Name("environment"), q.Required(true), q.Disabled(true),
			q.Title("Which cloud environment should it use?"),
			q.Choices(
				choice("preview", "Preview"),
				choice("staging", "Staging"),
				choice("isolated", "Isolated sandbox"),
			),
			q.Error(),
		),
		item("approval", "When should the agent request approval?", "",
			choice("writes", "Before writing files"),
			choice("commands", "Before running commands"),
			choice("sensitive", "Only for sensitive actions"),
		),
		q.Actions(q.Previous(), q.Next("Next"), q.Submit("Save execution plan")),
	)
}

func Controlled() h.Node {
	return h.Div(h.Class("relative mx-auto flex h-full w-full max-w-md flex-col"),
		h.P(h.Class("absolute end-0 top-0 text-sm text-muted-foreground"), h.Role("status"),
			"Current checkpoint: ", h.Comment(""), h.Span(h.Data("example-questionnaire-checkpoint", ""), "Change scope"),
		),
		q.Questionnaire(example("controlled"), h.Class("mt-auto"), q.ActiveItem("scope"),
			q.Items(def("scope", true), def("checks", true), def("output", true)),
			q.Progress(),
			item("scope", "What may the agent change?", "The host stores the active checkpoint while Questionnaire navigates.",
				choice("component", "Only the target component"),
				choice("tests", "Component and related tests"),
				choice("feature", "The complete feature area"),
			),
			item("checks", "Which verification level should it use?", "",
				choice("targeted", "Targeted tests"),
				choice("package", "Package tests and typecheck"),
				choice("full", "Full workspace verification"),
			),
			item("output", "What should the agent return when finished?", "",
				choice("summary", "Concise summary"),
				choice("diff", "Summary with changed files"),
				choice("handoff", "Detailed implementation handoff"),
			),
			q.Actions(q.Previous(), q.Next("Next"), q.Submit("Save workflow")),
		),
	)
}

func Dialog() h.Node {
	dialogItem := func(name, title, description string, choices ...h.Node) h.Node {
		return q.Item(q.Name(name), q.Required(true),
			dialog.Header(
				q.Progress(),
				q.Title(q.Render(dialog.Title()), title),
				q.Description(q.Render(dialog.Description()), description),
			),
			q.Choices(nodes(choices)...),
			q.Error(),
		)
	}
	return dialog.Dialog(
		dialog.Trigger(lib.AsChild, button.Button(button.VariantOutline, "Open clarification")),
		dialog.Content(
			q.Questionnaire(example("dialog"), q.DefaultItem("scope"),
				q.Items(def("scope", true), def("tests", true)),
				dialogItem("scope", "Which files are in scope?", "Choose how broadly the agent can update the workspace.",
					choice("component", "Component only"),
					choice("feature", "Complete feature directory"),
					choice("workspace", "Any related workspace file"),
				),
				dialogItem("tests", "How much verification is needed?", "Choose the checks the agent should run before handoff.",
					choice("targeted", "Targeted tests"),
					choice("package", "Package tests"),
					choice("full", "Full workspace verification"),
				),
				dialog.Footer(
					dialog.Close(lib.AsChild, button.Button(h.Type("button"), button.VariantOutline, "Cancel")),
					q.Actions(q.Previous(), q.Next("Next"), q.Submit("Send answer")),
				),
			),
		),
	)
}

func Freeform() h.Node {
	return q.Questionnaire(example("freeform"), h.Class("mx-auto max-w-md"), q.ShortcutsLetters,
		q.Items(def("approach", true, "incremental", "module", "rewrite")),
		q.Item(q.Name("approach"), q.Required(true),
			q.Title("How should the agent approach this refactor?"),
			q.Description("Choose a strategy or write a more specific instruction."),
			q.Choices(
				choice("incremental", "Make the smallest safe change"),
				choice("module", "Refactor one module at a time"),
				choice("rewrite", "Replace the implementation completely"),
				q.Input(h.Aria("label", "Another refactoring approach"), h.Placeholder("Describe another approach…")),
			),
			q.Error(),
		),
		q.Actions(q.Submit("Use this approach")),
	)
}

func Multiple() h.Node {
	return q.Questionnaire(example("multiple"), h.Class("mx-auto max-w-md"), q.ShortcutsLetters,
		q.Items(def("context", true, "source", "tests", "docs", "history")),
		q.Item(q.Name("context"), q.Multiple(true), q.Required(true),
			q.Title("What context should the agent inspect?"),
			q.Description("Select every source that may affect the implementation."),
			q.Choices(
				choice("source", "Relevant source files"),
				choice("tests", "Existing tests"),
				choice("docs", "Architecture documentation"),
				choice("history", "Recent commit history"),
			),
			q.Error(),
		),
		q.Actions(q.Submit("Share context")),
	)
}

func NavigationState() h.Node {
	return q.Questionnaire(example("navigation-state"), h.Class("mx-auto max-w-md"), q.ActiveItem("permission"),
		q.Items(def("permission", true), def("verification", true)),
		q.Progress(),
		item("permission", "What may the agent modify?", "Next is intentionally disabled until an answer is selected.",
			choice("files", "Project files"),
			choice("tests", "Project files and tests"),
			choice("config", "Files, tests, and configuration"),
		),
		item("verification", "What must pass before completion?", "",
			choice("tests", "Tests"),
			choice("types", "Tests and types"),
			choice("all", "Tests, types, and visual QA"),
		),
		q.Actions(
			q.Previous(),
			q.Next(h.Class("data-[status=unanswered]:opacity-50"), h.Disabled(), button.VariantSecondary, "Next"),
			q.Submit(h.Disabled(), "Save permissions"),
		),
	)
}

func Progress() h.Node {
	bars := make([]any, 4)
	for i := range bars {
		class := "h-1.5 flex-1 rounded-full bg-muted"
		if i < 1 {
			class = "h-1.5 flex-1 rounded-full bg-primary"
		}
		bars[i] = h.Span(h.Class(class))
	}
	return q.Questionnaire(example("progress"), h.Class("mx-auto max-w-md"), q.DefaultItem("scope"),
		q.Items(def("scope", true), def("strategy", true), def("tests", true), def("delivery", true)),
		q.Progress(h.Class("w-full"),
			h.Div(h.Class("mb-2 flex gap-1.5"), h.Aria("hidden", "true"), h.Data("example-questionnaire-bars", ""), bars),
			h.Span("Checkpoint ", h.Comment(""), h.Span(h.Data("questionnaire-current", ""), "1"), h.Comment(""), " of ", h.Comment(""), h.Span(h.Data("questionnaire-total", ""), "4")),
		),
		item("scope", "How large is the change?", "",
			choice("small", "Small patch"),
			choice("medium", "Feature-sized change"),
			choice("large", "Cross-package change"),
		),
		item("strategy", "How should commits be organized?", "",
			choice("single", "Single commit"),
			choice("logical", "Logical commits"),
			choice("squash", "Squash before review"),
		),
		item("tests", "Which tests should run?", "",
			choice("targeted", "Targeted tests"),
			choice("package", "Package suite"),
			choice("workspace", "Full workspace"),
		),
		item("delivery", "How should the work be delivered?", "",
			choice("patch", "Patch only"),
			choice("commit", "Committed locally"),
			choice("branch", "Push a review branch"),
		),
		q.Actions(q.Previous(), q.Next("Next"), q.Submit("Finish plan")),
	)
}

func Resume() h.Node {
	return q.Questionnaire(example("resume"), h.Class("mx-auto max-w-md"), q.DefaultItem("verification"),
		q.Items(def("change", true), def("verification", true), def("notes", false)),
		q.Progress(),
		item("change", "What kind of migration is this?", "This answer was saved during the previous session.",
			q.Choice(q.Value("incremental"), q.DefaultChecked(true), "Incremental migration"),
			choice("cutover", "Single cutover"),
		),
		q.Item(q.Name("verification"), q.Multiple(true), q.Required(true),
			q.Title("How should the migration be verified?"),
			q.Description("These checks were selected during the previous session."),
			q.Choices(
				q.Choice(q.Value("tests"), q.DefaultChecked(true), "Run migration tests"),
				q.Choice(q.Value("typecheck"), q.DefaultChecked(true), "Run the typecheck"),
				choice("manual", "Perform a manual smoke test"),
			),
			q.Error(),
		),
		q.Item(q.Name("notes"),
			q.Title("Anything else the agent should remember?"),
			q.Description("This note was saved with the draft."),
			q.Input(h.Aria("label", "Saved migration note"), h.Value("Keep the existing public API stable.")),
		),
		q.Actions(
			button.Button(h.Type("reset"), button.VariantOutline, "Reset changes"),
			q.Previous(), q.Next("Next"), q.Submit("Update draft"),
		),
	)
}

func Shortcuts() h.Node {
	return h.Div(h.Class("relative mx-auto flex h-full w-full max-w-md flex-col"),
		nativeselect.NativeSelect(h.Aria("label", "Shortcut style"), h.Class("absolute end-0 top-0"), h.Data("example-questionnaire-shortcuts", ""),
			nativeselect.Option(h.Value("none"), "No shortcuts"),
			nativeselect.Option(h.Value("letters"), h.Selected(), "Letters"),
			nativeselect.Option(h.Value("numbers"), "Numbers"),
		),
		q.Questionnaire(example("shortcuts"), h.Class("mt-auto"), q.ShortcutsLetters,
			q.Items(def("action", true, "inspect", "tests", "patch")),
			q.Item(q.Name("action"), q.Required(true),
				q.Title("What should the agent do next?"),
				q.Description("Use the displayed shortcut or navigate with the keyboard."),
				q.Choices(
					choice("inspect", "Inspect the implementation"),
					choice("tests", "Run the relevant tests"),
					choice("patch", "Prepare the patch"),
				),
				q.Error(),
			),
			q.Actions(q.Submit("Confirm action")),
		),
	)
}

func Skip() h.Node {
	return q.Questionnaire(example("skip"), h.Class("mx-auto max-w-md"), q.DefaultItem("task"),
		q.Items(def("task", true), def("constraints", false), def("review", true)),
		q.Progress(),
		item("task", "What kind of change is this?", "Choose the category that best describes the work.",
			choice("feature", "New feature"),
			choice("fix", "Bug fix"),
			choice("refactor", "Refactor"),
		),
		q.Item(q.Name("constraints"),
			q.Title("Are there any implementation constraints?"),
			q.Description("Answer if needed, or intentionally skip this question."),
			q.Choices(
				choice("no-dependencies", "Do not add dependencies"),
				choice("no-migrations", "Do not change the database"),
				choice("preserve-api", "Preserve the public API"),
				q.Input(h.Aria("label", "Another implementation constraint"), h.Placeholder("Describe another constraint…")),
			),
		),
		item("review", "How should the work be reviewed?", "Choose the checks the agent should complete before handoff.",
			choice("tests", "Run the test suite"),
			choice("diff", "Review the final diff"),
			choice("both", "Tests and diff review"),
		),
		q.Actions(q.Previous(), q.Skip(), q.Next("Next"), q.Submit("Submit brief")),
	)
}

func Validation() h.Node {
	progress := func() h.Node {
		return q.Progress(h.Class("min-w-0"),
			h.Span(h.Data("questionnaire-current", ""), "1"), h.Comment(""), " / ", h.Comment(""), h.Span(h.Data("questionnaire-total", ""), "2"),
		)
	}
	validationItem := func(name, title, description string, choices ...h.Node) h.Node {
		return q.Item(q.Name(name), q.Required(true),
			card.Header(q.Title(title), q.Description(description), card.Action(progress())),
			card.Content(q.Choices(nodes(choices)...), q.Error()),
		)
	}
	return q.Questionnaire(example("validation"), h.Class("mx-auto max-w-md"), q.ActiveItem("detail"),
		q.Items(def("detail", true), def("audience", true)),
		card.Card(h.Class("w-full"),
			validationItem("detail", "How much detail should the answer include?", "Choose the response depth.",
				choice("summary", "Concise summary"),
				choice("complete", "Complete answer"),
			),
			validationItem("audience", "Who will read the answer?", "Public answers require complete context.",
				choice("team", "My team"),
				choice("public", "Public audience"),
			),
			card.Footer(q.Actions(q.Previous(), q.Next("Next"), q.Submit("Validate answers"))),
		),
	)
}
