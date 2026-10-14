package blocks

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/input"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Login and signup blocks (blocks/login-0N, blocks/signup-0N).

func init() {
	ex := func(name, component string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: component, Render: registry.Static(fn)}
	}
	registry.Register(
		ex("login-01", "login", Login01),
		ex("login-02", "login", Login02),
		ex("login-03", "login", Login03),
		ex("login-04", "login", Login04),
		ex("login-05", "login", Login05),
		ex("signup-01", "signup", Signup01),
		ex("signup-02", "signup", Signup02),
		ex("signup-03", "signup", Signup03),
		ex("signup-04", "signup", Signup04),
		ex("signup-05", "signup", Signup05),
	)
}

// authBrandIcon renders an inline brand <svg> as the blocks write it.
func authBrandIcon(path string) h.Node {
	return h.Svg(h.Xmlns("http://www.w3.org/2000/svg"), h.ViewBox("0 0 24 24"),
		h.El("path", h.D(path), h.Fill("currentColor")),
	)
}

// authOutline is an outline type="button" button.
func authOutline(children ...any) h.Node {
	return button.Button(append([]any{button.VariantOutline, h.Type("button")}, children...)...)
}

func authSubmit(label string) h.Node {
	return button.Button(h.Type("submit"), label)
}

// authInputField is <Field><FieldLabel/><Input/>...extra</Field>.
func authInputField(id, typ, label, placeholder string, extra ...any) h.Node {
	in := []any{h.ID(id), h.Type(typ)}
	if placeholder != "" {
		in = append(in, h.Placeholder(placeholder))
	}
	in = append(in, h.Required())
	return field.Field(append([]any{field.Label(h.For(id), label), input.Input(in...)}, extra...)...)
}

// authPasswordField is the password field with the "Forgot your password?" link.
func authPasswordField(linkClass string) h.Node {
	return field.Field(
		h.Div(h.Class("flex items-center"),
			field.Label(h.For("password"), "Password"),
			h.A(h.Href("#"), h.Class(linkClass), "Forgot your password?"),
		),
		input.Input(h.ID("password"), h.Type("password"), h.Required()),
	)
}

// authLogo is the "Acme Inc." link of the page layouts.
func authLogo(class string) h.Node {
	return h.A(h.Href("#"), h.Class(class),
		h.Div(h.Class("flex size-6 items-center justify-center rounded-md bg-primary text-primary-foreground"),
			lucide.GalleryVerticalEnd(h.Class("size-4")),
		),
		"Acme Inc.",
	)
}

// authImage is the placeholder cover image.
func authImage() h.Node {
	return h.Img(h.Src("/placeholder.svg"), h.Alt("Image"),
		h.Class("absolute inset-0 h-full w-full object-cover dark:brightness-[0.2] dark:grayscale"))
}

// authSplitPage is the two-column layout of login-02 and signup-02.
func authSplitPage(form h.Node) h.Node {
	return h.Div(h.Class("grid min-h-svh lg:grid-cols-2"),
		h.Div(h.Class("flex flex-col gap-4 p-6 md:p-10"),
			h.Div(h.Class("flex justify-center gap-2 md:justify-start"),
				authLogo("flex items-center gap-2 font-medium"),
			),
			h.Div(h.Class("flex flex-1 items-center justify-center"),
				h.Div(h.Class("w-full max-w-xs"), form),
			),
		),
		h.Div(h.Class("relative hidden bg-muted lg:block"), authImage()),
	)
}

// authTerms is the "By clicking continue..." footer.
func authTerms() h.Node {
	return field.Description(h.Class("px-6 text-center"),
		"By clicking continue, you agree to our ", h.A(h.Href("#"), "Terms of Service"), " ",
		"and ", h.A(h.Href("#"), "Privacy Policy"), ".",
	)
}

// authSwitch is "<question> <a href=#>link</a>" in a FieldDescription.
func authSwitch(class, question, link string) h.Node {
	args := []any{}
	if class != "" {
		args = append(args, h.Class(class))
	}
	return field.Description(append(args, question+" ", h.A(h.Href("#"), link))...)
}

// authPasswordPair is the side-by-side password / confirm password field.
func authPasswordPair() h.Node {
	return field.Field(
		field.Field(h.Class("grid grid-cols-2 gap-4"),
			field.Field(
				field.Label(h.For("password"), "Password"),
				input.Input(h.ID("password"), h.Type("password"), h.Required()),
			),
			field.Field(
				field.Label(h.For("confirm-password"), "Confirm Password"),
				input.Input(h.ID("confirm-password"), h.Type("password"), h.Required()),
			),
		),
		field.Description("Must be at least 8 characters long."),
	)
}

// authCenteredHeader is the logo + title header of login-05 and signup-05.
func authCenteredHeader(question, link string) h.Node {
	return h.Div(h.Class("flex flex-col items-center gap-2 text-center"),
		h.A(h.Href("#"), h.Class("flex flex-col items-center gap-2 font-medium"),
			h.Div(h.Class("flex size-8 items-center justify-center rounded-md"),
				lucide.GalleryVerticalEnd(h.Class("size-6")),
			),
			h.Span(h.Class("sr-only"), "Acme Inc."),
		),
		h.H1(h.Class("text-xl font-bold"), "Welcome to Acme Inc."),
		authSwitch("", question, link),
	)
}
