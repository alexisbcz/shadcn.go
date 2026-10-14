// Package form is the shadcn/ui Form: building blocks for accessible forms
// whose fields show validation state (upstream form.tsx, made for
// react-hook-form).
//
//	h.Form(h.ID("profile"),
//		form.Field("username",
//			form.Item(form.Error(errs["username"]),
//				form.Label("Username"),
//				form.Control(input.Input(h.Name("username"), h.Placeholder("shadcn"))),
//				form.Description("This is your public display name."),
//				form.Message(),
//			),
//		),
//	)
//
// Item generates the ids that tie its label, control, description and
// message together, and renders the field's error state (Error) exactly
// like upstream's useFormField: data-error on the label, aria-invalid and
// aria-describedby on the control, the message.
//
// On the client, form.js (module "form") updates that state, and
// form-state.js (module "form-state") provides the validation semantics of
// react-hook-form, TanStack Form and Formisch (see its header). With
// form.js's bindForm, a form of Items validates like upstream's
// react-hook-form + zodResolver forms:
//
//	import { bindForm } from "/_shadcn/form.js"
//	import { rules, issuesOf } from "/_shadcn/form-state.js"
//	bindForm(document.querySelector("form#profile"), {
//		defaults: { username: "" },
//		schema: (v) => issuesOf(["username"], v.username, [rules.minLength(2, "Username must be at least 2 characters.")]),
//		onSubmit: (values) => { ... },
//	})
package form

import (
	_ "embed"
	"fmt"
	"io"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/internal/uid"
	"shadcn.go/components/ui/label"
	"shadcn.go/components/ui/runtime"
	h "shadcn.go/libhtml"
)

//go:embed form.js
var script string

//go:embed form-state.js
var stateScript string

func init() {
	runtime.Module("form", script)
	runtime.Module("form-state", stateScript)
}

// Form is upstream's Form (react-hook-form's FormProvider): it renders its
// children only.
func Form(children ...any) h.Node { return h.Group(children...) }

// Field is upstream's FormField (a react-hook-form Controller): it renders
// its children only; name is the field's name, which the client binding
// reads from the control's name attribute.
func Field(name string, children ...any) h.Node {
	_ = name
	return h.Group(children...)
}

// ErrorMessage is the field's current validation error, a child of Item.
// An empty message means no error.
type ErrorMessage string

// Error returns the field's validation error state for Item.
func Error(message string) ErrorMessage { return ErrorMessage(message) }

// ItemBuilder is returned by Item.
type ItemBuilder struct {
	h.Builder[ItemBuilder]
	err ErrorMessage
	id  string
}

// Item is upstream's FormItem: a <div data-slot="form-item"> laying out a
// field. It accepts an ErrorMessage child.
func Item(children ...any) *ItemBuilder {
	b := &ItemBuilder{}
	b.Builder = h.NewBuilder(b, "div")
	return b.Children(children...)
}

// AcceptChild consumes the ErrorMessage.
func (b *ItemBuilder) AcceptChild(c any) bool {
	if e, ok := c.(ErrorMessage); ok {
		b.err = e
		return true
	}
	return false
}

// ID sets the id prefix of the item's parts (default: generated, like
// React's useId).
func (b *ItemBuilder) ID(id string) *ItemBuilder { b.id = id; return b }

var itemBase = h.Div(h.Data("slot", "form-item"), h.Class(lib.Cn("grid gap-2")))

// Element resolves the item and its parts.
func (b *ItemBuilder) Element() *h.Element {
	id := b.id
	if id == "" {
		id = uid.New("form")
	}
	el := h.Spread(itemBase, b.Elem(), lib.Cn)
	ids := partIDs{
		item:        id + "-form-item",
		description: id + "-form-item-description",
		message:     id + "-form-item-message",
	}
	el.SetChildren(resolveParts(el.ChildNodes(), ids, string(b.err))...)
	return el
}

// Render renders the item.
func (b *ItemBuilder) Render(w io.Writer) error { return b.Element().Render(w) }

type partIDs struct{ item, description, message string }

// part is a placeholder for Label, Control, Description and Message: the
// enclosing Item resolves it with the field's ids and error state.
type part struct {
	kind     string
	children []any
}

func (p *part) Render(w io.Writer) error {
	return fmt.Errorf("form: %s must be inside a form.Item", p.kind)
}

// Label is upstream's FormLabel: a Label for the item's control, red when
// the field has an error.
func Label(children ...any) h.Node { return &part{kind: "label", children: children} }

// Control is upstream's FormControl: it merges the field's id and ARIA
// state into its single element child (a Radix Slot).
func Control(child any) h.Node {
	return &part{kind: "control", children: []any{child}}
}

// Description is upstream's FormDescription.
func Description(children ...any) h.Node { return &part{kind: "description", children: children} }

// Message is upstream's FormMessage: the field's error message, or its
// children when the field has none; nothing if both are empty.
func Message(children ...any) h.Node { return &part{kind: "message", children: children} }

func boolString(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

func resolveParts(nodes []any, ids partIDs, err string) []any {
	out := make([]any, 0, len(nodes))
	for _, n := range nodes {
		switch v := n.(type) {
		case *part:
			out = append(out, resolvePart(v, ids, err))
		case *h.Element:
			v.SetChildren(resolveParts(v.ChildNodes(), ids, err)...)
			out = append(out, v)
		default:
			if el, ok := n.(interface{ Element() *h.Element }); ok {
				e := el.Element()
				e.SetChildren(resolveParts(e.ChildNodes(), ids, err)...)
				out = append(out, e)
				continue
			}
			out = append(out, n)
		}
	}
	return out
}

func resolvePart(p *part, ids partIDs, err string) h.Node {
	switch p.kind {
	case "label":
		return h.Spread(
			label.Label(h.Data("slot", "form-label"), h.Data("error", boolString(err != "")), h.Class("data-[error=true]:text-destructive"), h.For(ids.item)),
			h.Label(p.children...), lib.Cn)
	case "control":
		// Upstream's Slot hands these props to the child component, which
		// spreads them after its own: they win over the rendered child's
		// attributes (e.g. Input's data-slot="input").
		describedBy := ids.description
		if err != "" {
			describedBy += " " + ids.message
		}
		var target *h.Element
		for _, c := range h.El("div", p.children...).ChildNodes() {
			switch v := c.(type) {
			case *h.Element:
				target = v
			case interface{ Element() *h.Element }:
				target = v.Element()
			}
			if target != nil {
				break
			}
		}
		if target == nil {
			return h.Errored(fmt.Errorf("form: Control needs an element child"))
		}
		out := target.Clone()
		out.Attr("data-slot", "form-control")
		out.Attr("id", ids.item)
		out.Attr("aria-describedby", describedBy)
		out.Attr("aria-invalid", boolString(err != ""))
		return out
	case "description":
		return h.Spread(
			h.P(h.Data("slot", "form-description"), h.ID(ids.description), h.Class("text-sm text-muted-foreground")),
			h.P(p.children...), lib.Cn)
	case "message":
		props := h.P(p.children...)
		if err != "" {
			props.SetChildren(h.Text(err))
		} else if len(props.ChildNodes()) == 0 {
			return h.Fragment{}
		}
		return h.Spread(
			h.P(h.Data("slot", "form-message"), h.ID(ids.message), h.Class("text-sm text-destructive")),
			props, lib.Cn)
	}
	return h.Fragment{}
}
