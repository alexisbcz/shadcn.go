package main

import (
	"net/http"
	"net/mail"
	"strings"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/icons/tabler"
	"shadcn.go/components/lib"
	"shadcn.go/components/ui/button"
	"shadcn.go/components/ui/spinner"
	h "shadcn.go/libhtml"
	"shadcn.go/libhtml/csrf"
	"shadcn.go/libhtml/secure"
)

// layout is the page shell shared by every page.
func layout(r *http.Request, title string, children ...any) h.Node {
	return h.Doctype(h.Html(h.Lang("en"),
		h.Head(
			h.Meta(h.Charset("utf-8")),
			h.Meta(h.Name("viewport"), h.Content("width=device-width, initial-scale=1")),
			h.Title(title),
			h.Link(h.Rel("stylesheet"), h.Href("/static/app.css")),
			// Applies the saved theme before first paint. Runs only because
			// it carries this request's CSP nonce.
			secure.Script(r, h.UnsafeJS(themeScript)),
		),
		h.Body(h.Class("bg-background text-foreground antialiased"),
			h.Main(h.Class("mx-auto flex max-w-3xl flex-col gap-10 px-6 py-12"), h.Group(children...)),
		),
	))
}

const themeScript = `(() => {
  const root = document.documentElement;
  const saved = localStorage.getItem("theme");
  if (saved === "dark" || (!saved && matchMedia("(prefers-color-scheme: dark)").matches)) root.classList.add("dark");
  document.addEventListener("click", (e) => {
    if (!e.target.closest("[data-theme-toggle]")) return;
    const dark = root.classList.toggle("dark");
    localStorage.setItem("theme", dark ? "dark" : "light");
  });
})();`

// example is one titled demo, like the previews on ui.shadcn.com.
func example(title, description string, children ...any) *h.Element {
	return h.Section(h.Class("flex flex-col gap-3"),
		h.H2(h.Class("text-lg font-semibold tracking-tight"), title),
		h.If(description != "", h.P(h.Class("text-sm text-muted-foreground"), description)),
		h.Div(h.Class("flex min-h-28 flex-wrap items-center justify-center gap-2 rounded-lg border p-8"), h.Group(children...)),
	)
}


func handleIndex(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	page := layout(r, "Button · shadcn.go",
		h.Header(h.Class("flex items-start justify-between gap-4"),
			h.Div(h.Class("flex flex-col gap-2"),
				h.H1(h.Class("text-3xl font-bold tracking-tight"), "Button"),
				h.P(h.Class("text-muted-foreground"), "Displays a button or a component that looks like a button. Rendered in Go with libhtml."),
			),
			button.Button(button.VariantOutline, button.SizeSm, h.Data("theme-toggle", ""), h.Type("button"), "Toggle theme"),
		),

		example("Demo", "",
			button.Button(button.VariantOutline, "Button"),
			button.Button(button.VariantOutline, button.SizeIcon, h.Aria("label", "Submit"), lucide.ArrowUp()),
		),
		example("Variants", "",
			h.Map(button.VariantValues, func(v button.Variant) h.Node {
				return button.Button(v, strings.ToUpper(string(v[:1]))+string(v[1:]))
			}),
		),
		example("Size", "",
			h.Div(h.Class("flex flex-col items-start gap-8 sm:flex-row"),
				h.Div(h.Class("flex items-start gap-2"),
					button.Button(button.SizeSm, button.VariantOutline, "Small"),
					button.Button(button.SizeIconSm, h.Aria("label", "Submit"), button.VariantOutline, lucide.ArrowUpRight()),
				),
				h.Div(h.Class("flex items-start gap-2"),
					button.Button(button.VariantOutline, "Default"),
					button.Button(button.SizeIcon, h.Aria("label", "Submit"), button.VariantOutline, lucide.ArrowUpRight()),
				),
				h.Div(h.Class("flex items-start gap-2"),
					button.Button(button.VariantOutline, button.SizeLg, "Large"),
					button.Button(button.SizeIconLg, h.Aria("label", "Submit"), button.VariantOutline, lucide.ArrowUpRight()),
				),
			),
		),
		example("Icon", "",
			button.Button(button.VariantOutline, button.SizeIcon, h.Aria("label", "Upgrade"), lucide.CircleFadingArrowUp()),
		),
		example("With Icon", "",
			button.Button(button.VariantOutline, button.SizeSm, tabler.IconGitBranch(), " New Branch"),
		),
		example("Rounded", "A className override: twmerge replaces rounded-md with rounded-full.",
			button.Button(button.VariantOutline, button.SizeIcon, h.Class("rounded-full"), h.Aria("label", "Up"), lucide.ArrowUp()),
		),
		example("Loading", "",
			button.Button(button.SizeSm, button.VariantOutline, h.Disabled(), spinner.Spinner(), "Submit"),
		),
		example("As Child", "The button's classes and attributes are merged into the link.",
			button.Button(lib.AsChild, h.A(h.Href("#login"), "Login")),
		),

		example("Security: CSRF", "This form carries a masked CSRF token. Remove it in the devtools and the POST is rejected with 403.",
			h.Form(h.Method("post"), h.Action("/subscribe"), h.Class("flex w-full max-w-sm items-center gap-2"),
				csrf.Field(r),
				h.Input(h.Type("email"), h.Name("email"), h.Required(), h.Placeholder("you@example.com"),
					h.Aria("label", "Email"),
					h.Class("h-9 w-full min-w-0 rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-xs outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50")),
				button.Button(h.Type("submit"), "Subscribe"),
			),
		),
		example("Security: XSS", "Whatever you put in ?name= is escaped. Try the link.",
			h.Div(h.Class("flex flex-col items-center gap-3 text-sm"),
				h.P("Hello, ", h.Code(h.Class("rounded bg-muted px-1 py-0.5"), h.IfElse(name != "", name, "stranger")), "!"),
				button.Button(lib.AsChild, button.VariantLink,
					h.A(h.Href(`/?name=<img src=x onerror=alert(1)><script>alert(1)</script>`), "Inject a script tag"),
				),
			),
		),
	)
	if err := h.WriteHTTP(w, http.StatusOK, page); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func handleSubscribe(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	_, err := mail.ParseAddress(email)
	valid := err == nil && len(email) <= 254

	page := layout(r, "Subscribed · shadcn.go",
		h.Div(h.Class("flex flex-col items-start gap-4"),
			h.IfElseF(valid,
				func() h.Node {
					return h.P(h.Class("text-lg"), "Subscribed ", h.Strong(email), ". The CSRF token checked out.")
				},
				func() h.Node {
					return h.P(h.Class("text-lg text-destructive"), "That email address is not valid.")
				},
			),
			button.Button(lib.AsChild, button.VariantOutline, h.A(h.Href("/"), "Back")),
		),
	)
	status := http.StatusOK
	if !valid {
		status = http.StatusUnprocessableEntity
	}
	if err := h.WriteHTTP(w, status, page); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
