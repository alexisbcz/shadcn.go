package site

import (
	"net/http"
	"net/url"

	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/input"
	h "shadcn.go/libhtml"
)

// ExampleAuthentication serves /examples/authentication
// (app/(app)/examples/authentication/page.tsx).
func (s *Site) ExampleAuthentication(r *http.Request) (h.Node, error) {
	return s.examplesLayoutMeta(r, Meta{Title: "Authentication", PlainTitle: true, Description: "Authentication forms built using the components.", Canonical: "/examples/authentication"},
		h.Div(h.Class("md:hidden"),
			authExampleImage("/examples/authentication-light.png", "Authentication", "block dark:hidden"),
			authExampleImage("/examples/authentication-dark.png", "Authentication", "hidden dark:block"),
		),
		h.Div(h.Class("relative container hidden flex-1 shrink-0 items-center justify-center md:grid lg:max-w-none lg:grid-cols-2 lg:px-0"),
			button.Button(lib.AsChild, button.VariantGhost, h.Class("absolute top-4 right-4 md:top-8 md:right-8"),
				h.A(h.Href("/examples/authentication"), "Login")),
			h.Div(h.Class("relative hidden h-full flex-col p-10 text-primary lg:flex dark:border-r"),
				h.Div(h.Class("absolute inset-0 bg-primary/5")),
				h.Div(h.Class("relative z-20 flex items-center text-lg font-medium"),
					h.El("svg", h.Attr("xmlns", "http://www.w3.org/2000/svg"), h.Attr("viewBox", "0 0 24 24"), h.Attr("fill", "none"),
						h.Attr("stroke", "currentColor"), h.Attr("stroke-width", "2"), h.Attr("stroke-linecap", "round"),
						h.Attr("stroke-linejoin", "round"), h.Class("mr-2 h-6 w-6"),
						h.Path(h.Attr("d", "M15 6v12a3 3 0 1 0 3-3H6a3 3 0 1 0 3 3V6a3 3 0 1 0-3 3h12a3 3 0 1 0-3-3")),
					),
					"Acme Inc",
				),
				h.Div(h.Class("relative z-20 mt-auto"),
					h.Blockquote(h.Class("leading-normal text-balance"),
						"“This library has saved me countless hours of work and helped me deliver stunning designs to my clients faster than ever before.” - Sofia Davis"),
				),
			),
			h.Div(h.Class("flex items-center justify-center lg:h-[1000px] lg:p-8"),
				h.Div(h.Class("mx-auto flex w-full flex-col justify-center gap-6 sm:w-[350px]"),
					h.Div(h.Class("flex flex-col gap-2 text-center"),
						h.H1(h.Class("text-2xl font-semibold tracking-tight"), "Create an account"),
						h.P(h.Class("text-sm text-muted-foreground"), "Enter your email below to create your account"),
					),
					userAuthForm(),
					field.Description(h.Class("px-6 text-center"),
						"By clicking continue, you agree to our ",
						h.A(h.Href("/terms"), "Terms of Service"), " and ",
						h.A(h.Href("/privacy"), "Privacy Policy"), ".",
					),
				),
			),
		),
		h.Script(h.Type("module"), h.Src(staticPrefix+"examples-auth.js")),
	), nil
}

// userAuthForm is examples/authentication/components/user-auth-form.tsx;
// examples-auth.js plays the 3s loading state on submit.
func userAuthForm() h.Node {
	return h.Div(h.Class("grid gap-6"), h.Data("user-auth-form", ""),
		h.Form(
			field.Group(
				field.Field(
					field.Label(h.Class("sr-only"), h.For("email"), "Email"),
					input.Input(h.ID("email"), h.Placeholder("name@example.com"), h.Type("email"),
						h.Attr("autocapitalize", "none"), h.Attr("autocomplete", "email"), h.Attr("autocorrect", "off")),
				),
				field.Field(
					button.Button("Sign In with Email"),
				),
			),
		),
		field.Separator("Or continue with"),
		button.Button(button.VariantOutline, h.Type("button"),
			h.El("svg", h.Attr("viewBox", "0 0 438.549 438.549"), h.Class("mr-2 h-4 w-4"),
				h.Path(h.Attr("fill", "currentColor"), h.Attr("d", githubPath)),
			),
			" GitHub",
		),
	)
}

// authExampleImage is a priority next/image of a public image (1x/2x srcset at
// the optimizer's 1920 and 3840 widths).
func authExampleImage(src, alt, class string) h.Node {
	u := "/_next/image?url=" + url.QueryEscape(src)
	return h.Img(h.Alt(alt), h.Width("1280"), h.Height("843"), h.Attr("decoding", "async"), h.Data("nimg", "1"),
		h.Class(class), h.Style("color:transparent"),
		h.UnsafeAttr("srcset", u+"&w=1920&q=75 1x, "+u+"&w=3840&q=75 2x"),
		h.UnsafeAttr("src", u+"&w=3840&q=75"))
}
