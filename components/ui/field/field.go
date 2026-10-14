// Package field is the shadcn/ui Field: building blocks to compose
// accessible form fields (labels, controls, descriptions, errors) and
// groups of fields.
//
//	field.Set(
//		field.Legend("Profile"),
//		field.Group(
//			field.Field(
//				field.Label(h.For("name"), "Full name"),
//				input.Input(h.ID("name")),
//				field.Description("This appears on invoices and emails."),
//			),
//		),
//	)
package field

import (
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/internal/styled"
	"shadcn.go/components/ui/label"
	"shadcn.go/components/ui/separator"
	"shadcn.go/cva"
	h "shadcn.go/libhtml"
)

// Orientation lays out a Field.
type Orientation string

// Orientations.
const (
	Vertical   Orientation = "vertical"
	Horizontal Orientation = "horizontal"
	Responsive Orientation = "responsive"
)

// LegendVariant selects the size of a Legend.
type LegendVariant string

// Legend variants.
const (
	LegendLegend LegendVariant = "legend"
	LegendLabel  LegendVariant = "label"
)

var fieldVariants = cva.New(
	"group/field flex w-full gap-3 data-[invalid=true]:text-destructive",
	cva.Config{
		Variants: cva.Variants{
			"orientation": {
				"vertical": "flex-col [&>*]:w-full [&>.sr-only]:w-auto",
				"horizontal": "flex-row items-center " +
					"[&>[data-slot=field-label]]:flex-auto " +
					"has-[>[data-slot=field-content]]:items-start has-[>[data-slot=field-content]]:[&>[role=checkbox],[role=radio]]:mt-px",
				"responsive": "flex-col @md/field-group:flex-row @md/field-group:items-center [&>*]:w-full @md/field-group:[&>*]:w-auto [&>.sr-only]:w-auto " +
					"@md/field-group:[&>[data-slot=field-label]]:flex-auto " +
					"@md/field-group:has-[>[data-slot=field-content]]:items-start @md/field-group:has-[>[data-slot=field-content]]:[&>[role=checkbox],[role=radio]]:mt-px",
			},
		},
		Defaults: cva.Selection{"orientation": "vertical"},
	},
)

// Variants returns the classes of a Field with the given orientation plus
// className, like fieldVariants() upstream.
func Variants(o Orientation, className ...string) string {
	return lib.Cn(fieldVariants.Classes(cva.Selection{"orientation": string(o)}, className...))
}

func slot(tag, name string, classes ...string) *h.Element {
	return h.El(tag, h.Data("slot", name), h.Class(lib.Cn(classes...)))
}

var (
	setBase = slot("fieldset", "field-set",
		"flex flex-col gap-6",
		"has-[>[data-slot=checkbox-group]]:gap-3 has-[>[data-slot=radio-group]]:gap-3")
	groupBase = slot("div", "field-group",
		"group/field-group @container/field-group flex w-full flex-col gap-7 data-[slot=checkbox-group]:gap-3 [&>[data-slot=field-group]]:gap-4")
	contentBase = slot("div", "field-content",
		"group/field-content flex flex-1 flex-col gap-1.5 leading-snug")
	labelClass = lib.Cn(
		"group/field-label peer/field-label flex w-fit gap-2 leading-snug group-data-[disabled=true]/field:opacity-50",
		"has-[>[data-slot=field]]:w-full has-[>[data-slot=field]]:flex-col has-[>[data-slot=field]]:rounded-md has-[>[data-slot=field]]:border [&>*]:data-[slot=field]:p-4",
		"has-data-[state=checked]:border-primary has-data-[state=checked]:bg-primary/5 dark:has-data-[state=checked]:bg-primary/10",
	)
	titleBase = slot("div", "field-label",
		"flex w-fit items-center gap-2 text-sm leading-snug font-medium group-data-[disabled=true]/field:opacity-50")
	descriptionBase = slot("p", "field-description",
		"text-sm leading-normal font-normal text-muted-foreground group-has-[[data-orientation=horizontal]]/field:text-balance",
		"last:mt-0 nth-last-2:-mt-1 [[data-variant=legend]+&]:-mt-1.5",
		"[&>a]:underline [&>a]:underline-offset-4 [&>a:hover]:text-primary")
	separatorBase = slot("div", "field-separator",
		"relative -my-2 h-5 text-sm group-data-[variant=outline]/field-group:-mb-2")
	errorBase   = h.Div(h.Role("alert"), h.Data("slot", "field-error"), h.Class(lib.Cn("text-sm font-normal text-destructive")))
	legendClass = lib.Cn("mb-3 font-medium", "data-[variant=legend]:text-base", "data-[variant=label]:text-sm")
)

func spread(base *h.Element, children []any) *h.Element {
	return h.Spread(base, h.El(base.Tag(), children...), lib.Cn)
}

// Set is a <fieldset> grouping related fields (upstream FieldSet).
func Set(children ...any) *styled.Node {
	return styledSpread(setBase, children, styled.Slot("fieldset", "field-set", sSet))
}

// Group stacks fields (upstream FieldGroup); it is the container that
// Responsive fields query.
func Group(children ...any) *styled.Node {
	return styledSpread(groupBase, children, styled.Slot("div", "field-group", sGroup))
}

// Content groups a field's label and description beside its control
// (upstream FieldContent).
func Content(children ...any) *styled.Node {
	return styledSpread(contentBase, children, styled.Slot("div", "field-content", sContent))
}

// Label is the field's <label> (upstream FieldLabel, a Label).
func Label(children ...any) *styled.Node {
	props := h.Label(children...)
	className := props.ClassName()
	alt := props.Clone()
	props.SetClass(lib.Cn(labelClass, className))
	return styled.New(label.Label(h.Data("slot", "field-label"), spreadProp(props)), func(st style.Style) h.Node {
		p := alt.Clone().SetClass(sLabel.Class(st, className))
		return label.Label(h.Data("slot", "field-label"), spreadProp(p))
	})
}

// Title is a label-styled <div> for fields whose label wraps the control
// (upstream FieldTitle).
func Title(children ...any) *styled.Node {
	return styledSpread(titleBase, children, styled.Slot("div", "field-label", sTitle))
}

// Description is the field's helper text (upstream FieldDescription).
func Description(children ...any) *styled.Node {
	return styledSpread(descriptionBase, children, styled.Slot("p", "field-description", sDescription))
}

// Separator divides field groups, optionally with inline content (upstream
// FieldSeparator).
func Separator(children ...any) *styled.Node {
	props := h.Div(children...)
	kids := props.ChildNodes()
	props = props.Clone().SetChildren()
	hasContent := len(kids) > 0
	base := separatorBase.Clone().Data("content", boolString(hasContent))
	base.Children(separator.Separator(h.Class("absolute inset-0 top-1/2")))
	if hasContent {
		base.Children(h.Span(h.Class("relative mx-auto block w-fit bg-background px-2 text-muted-foreground"), h.Data("slot", "field-separator-content"), kids))
	}
	return styled.New(h.Spread(base, props, lib.Cn), func(st style.Style) h.Node {
		base := h.Div(h.Data("slot", "field-separator"), h.Data("content", boolString(hasContent)), h.Class(sSeparator.Class(st, "")))
		base.Children(separator.Separator(h.Class(sSepLine.Class(st))))
		if hasContent {
			base.Children(h.Span(h.Class(sSepContent.Class(st)), h.Data("slot", "field-separator-content"), kids))
		}
		return h.Spread(base, props, lib.Cn)
	})
}

// Error is the field's validation message (upstream FieldError). It
// renders nothing without children; see Errors for a list of messages.
func Error(children ...any) h.Node {
	props := h.Div(children...)
	if len(props.ChildNodes()) == 0 {
		return h.Fragment{}
	}
	return styled.Spread(errorBase, props, styled.Slot("div", "field-error", sError, h.Role("alert")))
}

// Errors renders the field's error messages like upstream FieldError's
// errors prop: duplicates are removed (keeping the first position of each
// message, with the last occurrence winning, like a Map), a single message
// is rendered as text and several as a list. Empty messages are skipped in
// the list. attrs are applied to the <div>.
func Errors(messages []string, attrs ...any) h.Node {
	if len(messages) == 0 {
		return h.Fragment{}
	}
	var unique []string
	seen := map[string]bool{}
	for _, m := range messages {
		if !seen[m] {
			seen[m] = true
			unique = append(unique, m)
		}
	}
	var content any
	if len(unique) == 1 {
		if unique[0] == "" {
			return h.Fragment{}
		}
		content = unique[0]
	} else {
		ul := h.Ul(h.Class("ml-4 flex list-disc flex-col gap-1"))
		for _, m := range unique {
			if m != "" {
				ul.Children(h.Li(m))
			}
		}
		content = ul
	}
	return styled.Spread(errorBase, h.Div(append(attrs, content)...), styled.Slot("div", "field-error", sError, h.Role("alert")))
}

// Builder is returned by Field and Legend.
type Builder struct {
	h.Builder[Builder]
	legend      bool
	orientation Orientation
	variant     LegendVariant
}

// Field is one form field (<div role="group">). Children may include an
// Orientation (vertical by default).
func Field(children ...any) *Builder {
	b := &Builder{orientation: Vertical}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// Legend is the <legend> of a Set. Children may include a LegendVariant.
func Legend(children ...any) *Builder {
	b := &Builder{legend: true, variant: LegendLegend}
	b.Builder = h.NewBuilder(b, "legend")
	return b.Children(children...)
}

// AcceptChild consumes Orientation (Field) and LegendVariant (Legend).
func (b *Builder) AcceptChild(c any) bool {
	switch v := c.(type) {
	case Orientation:
		if b.legend {
			return false
		}
		b.Orientation(v)
	case LegendVariant:
		if !b.legend {
			return false
		}
		b.Variant(v)
	default:
		return false
	}
	return true
}

// Orientation sets a Field's orientation.
func (b *Builder) Orientation(o Orientation) *Builder {
	if o == "" {
		o = Vertical
	}
	b.orientation = o
	return b
}

// Variant sets a Legend's variant.
func (b *Builder) Variant(v LegendVariant) *Builder {
	if v == "" {
		v = LegendLegend
	}
	b.variant = v
	return b
}

// Element resolves the component to its element.
func (b *Builder) Element() *h.Element { return b.ElementIn(nil) }

// ElementIn resolves the component to its element as it renders inside w
// (in the style of an enclosing style.Scope).
func (b *Builder) ElementIn(w io.Writer) *h.Element {
	st, _ := style.FromWriter(w)
	if b.legend {
		class := legendClass
		if st != "" {
			class = sLegend.Class(st, "")
		}
		base := h.Legend(h.Data("slot", "field-legend"), h.Data("variant", string(b.variant)), h.Class(class))
		return h.Spread(base, b.Elem(), lib.Cn)
	}
	props := b.Elem().Clone()
	className := props.ClassName()
	props.SetClass("")
	class := Variants(b.orientation, className)
	if st != "" {
		class = fieldClass(st, b.orientation, className)
	}
	base := h.Div(h.Role("group"), h.Data("slot", "field"), h.Data("orientation", string(b.orientation)), h.Class(class))
	return h.Spread(base, props, lib.Cn)
}

// Render renders the component.
func (b *Builder) Render(w io.Writer) error {
	return b.ElementIn(w).Render(w)
}

func boolString(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

// spreadProp applies props' attributes, classes and children to the
// element it is passed to, like {...props}.
func spreadProp(props *h.Element) h.Prop {
	return func(e *h.Element) { *e = *h.Spread(e, props, nil) }
}
