package card

import (
	_ "embed"

	"shadcn.go/components/style"
	"shadcn.go/components/ui/badge"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/input"
	"shadcn.go/components/ui/label"
	"shadcn.go/components/ui/runtime"
	togglegroup "shadcn.go/components/ui/toggle-group"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

//go:embed spacing.js
var spacingScript string

// Docs previews of the card page (ui/apps/v4/examples/radix/card-*.tsx), in
// radix-nova.
func init() {
	runtime.Module("example-card-spacing", spacingScript)
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "card", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("card-edge-to-edge", NovaEdgeToEdge),
		nova("card-image", NovaImage),
		nova("card-rtl", NovaRtl),
		nova("card-small", NovaSmall),
		nova("card-spacing", NovaSpacing),
	)
}

type loginText struct {
	title, description, signUp, email, placeholder, password, forgot, login, google string
}

var loginEn = loginText{"Login to your account", "Enter your email below to login to your account", "Sign Up", "Email", "m@example.com", "Password", "Forgot your password?", "Login", "Login with Google"}

func novaLogin(t loginText, idSuffix, forgotClass string, attrs ...any) h.Node {
	return card.Card(append(attrs,
		card.Header(
			card.Title(t.title),
			card.Description(t.description),
			card.Action(button.Button(button.VariantLink, t.signUp)),
		),
		card.Content(
			h.Form(
				h.Div(h.Class("flex flex-col gap-6"),
					h.Div(h.Class("grid gap-2"),
						label.Label(h.For("email"+idSuffix), t.email),
						input.Input(h.ID("email"+idSuffix), h.Type("email"), h.Placeholder(t.placeholder), h.Required()),
					),
					h.Div(h.Class("grid gap-2"),
						h.Div(h.Class("flex items-center"),
							label.Label(h.For("password"+idSuffix), t.password),
							h.A(h.Href("#"), h.Class(forgotClass+" inline-block text-sm underline-offset-4 hover:underline"), t.forgot),
						),
						input.Input(h.ID("password"+idSuffix), h.Type("password"), h.Required()),
					),
				),
			),
		),
		card.Footer(h.Class("flex-col gap-2"),
			button.Button(h.Type("submit"), h.Class("w-full"), t.login),
			button.Button(button.VariantOutline, h.Class("w-full"), t.google),
		),
	)...)
}

func NovaEdgeToEdge() h.Node {
	return card.Card(h.Class("mx-auto w-full max-w-sm"),
		card.Header(
			card.Title("Terms of Service"),
			card.Description("Review the terms before accepting the agreement."),
		),
		card.Content(h.Class("-mb-(--card-spacing)"),
			h.Div(h.Class("-mx-(--card-spacing) max-h-48 space-y-4 overflow-y-scroll border-t bg-muted/50 px-(--card-spacing) py-4 text-sm leading-relaxed"),
				h.P("These terms govern your use of the workspace, including access to shared documents, project files, and collaboration tools."),
				h.P("You are responsible for the content you upload and for ensuring that your team has the appropriate permissions to view or edit it."),
				h.P("We may update features or limits as the service evolves. When those changes materially affect your workflow, we will notify your workspace administrators."),
				h.P("By continuing, you agree to keep your account credentials secure and to follow your organization's acceptable use policies."),
			),
		),
		card.Footer(h.Class("justify-end gap-2"),
			button.Button(button.VariantOutline, "Decline"),
			button.Button("Accept"),
		),
	)
}

func NovaImage() h.Node {
	return card.Card(h.Class("relative mx-auto w-full max-w-sm pt-0"),
		h.Div(h.Class("absolute inset-0 z-30 aspect-video bg-black/35")),
		h.Img(h.Src("https://avatar.vercel.sh/shadcn1"), h.Alt("Event cover"),
			h.Class("relative z-20 aspect-video w-full object-cover brightness-60 grayscale dark:brightness-40")),
		card.Header(
			card.Action(badge.Badge(badge.VariantSecondary, "Featured")),
			card.Title("Design systems meetup"),
			card.Description("A practical talk on component APIs, accessibility, and shipping faster."),
		),
		card.Footer(button.Button(h.Class("w-full"), "View Event")),
	)
}

func NovaRtl() h.Node {
	t := loginText{"تسجيل الدخول إلى حسابك", "أدخل بريدك الإلكتروني أدناه لتسجيل الدخول إلى حسابك", "إنشاء حساب", "البريد الإلكتروني", "m@example.com", "كلمة المرور", "نسيت كلمة المرور؟", "تسجيل الدخول", "تسجيل الدخول باستخدام Google"}
	return style.ScopeRTL(novaLogin(t, "-rtl", "ms-auto", h.Class("w-full max-w-sm"), h.Dir("rtl")))
}

func NovaSmall() h.Node {
	return card.Card(card.SizeSm, h.Class("mx-auto w-full max-w-sm"),
		card.Header(
			card.Title("Small Card"),
			card.Description("This card uses the small size variant."),
		),
		card.Content(h.P(`The card component supports a size prop that can be set to "sm" for a more compact appearance.`)),
		card.Footer(button.Button(button.VariantOutline, button.SizeSm, h.Class("w-full"), "Action")),
	)
}

func NovaSpacing() h.Node {
	options := []struct{ class, label, value string }{
		{"[--card-spacing:--spacing(4)]", "16px", "4"},
		{"[--card-spacing:--spacing(5)]", "20px", "5"},
		{"[--card-spacing:--spacing(6)]", "24px", "6"},
		{"[--card-spacing:--spacing(8)]", "32px", "8"},
	}
	group := togglegroup.ToggleGroup(togglegroup.TypeSingle, togglegroup.Value("4"), togglegroup.VariantOutline, togglegroup.SizeSm, h.Class("justify-center"))
	for _, o := range options {
		group.Children(togglegroup.Item(togglegroup.Value(o.value), h.Data("spacing-class", o.class), o.label))
	}
	return h.Div(h.Class("mx-auto grid w-full max-w-sm gap-4"), h.Data("example-card-spacing", ""),
		group,
		novaLogin(loginEn, "-spacing", "ml-auto", h.Class(options[0].class)),
	)
}
