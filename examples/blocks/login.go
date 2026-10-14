package blocks

import (
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/field"
	h "shadcn.go/libhtml"
)

// Login01 ports blocks/login-01.
func Login01() h.Node {
	return h.Div(h.Class("flex min-h-svh w-full items-center justify-center p-6 md:p-10"),
		h.Div(h.Class("w-full max-w-sm"),
			h.Div(h.Class("flex flex-col gap-6"),
				card.Card(
					card.Header(
						card.Title("Login to your account"),
						card.Description("Enter your email below to login to your account"),
					),
					card.Content(
						h.Form(
							field.Group(
								authInputField("email", "email", "Email", "m@example.com"),
								authPasswordField("ml-auto inline-block text-sm underline-offset-4 hover:underline"),
								field.Field(
									authSubmit("Login"),
									authOutline("Login with Google"),
									authSwitch("text-center", "Don't have an account?", "Sign up"),
								),
							),
						),
					),
				),
			),
		),
	)
}

// Login02 ports blocks/login-02.
func Login02() h.Node {
	return authSplitPage(
		h.Form(h.Class("flex flex-col gap-6"),
			field.Group(
				h.Div(h.Class("flex flex-col items-center gap-1 text-center"),
					h.H1(h.Class("text-2xl font-bold"), "Login to your account"),
					h.P(h.Class("text-sm text-balance text-muted-foreground"), "Enter your email below to login to your account"),
				),
				authInputField("email", "email", "Email", "m@example.com"),
				authPasswordField("ml-auto text-sm underline-offset-4 hover:underline"),
				field.Field(authSubmit("Login")),
				field.Separator("Or continue with"),
				field.Field(
					authOutline(authBrandIcon(authGitHubPath), "Login with GitHub"),
					field.Description(h.Class("text-center"),
						"Don't have an account? ",
						h.A(h.Href("#"), h.Class("underline underline-offset-4"), "Sign up"),
					),
				),
			),
		),
	)
}

// Login03 ports blocks/login-03.
func Login03() h.Node {
	return h.Div(h.Class("flex min-h-svh flex-col items-center justify-center gap-6 bg-muted p-6 md:p-10"),
		h.Div(h.Class("flex w-full max-w-sm flex-col gap-6"),
			authLogo("flex items-center gap-2 self-center font-medium"),
			h.Div(h.Class("flex flex-col gap-6"),
				card.Card(
					card.Header(h.Class("text-center"),
						card.Title(h.Class("text-xl"), "Welcome back"),
						card.Description("Login with your Apple or Google account"),
					),
					card.Content(
						h.Form(
							field.Group(
								field.Field(
									authOutline(authBrandIcon(authApplePath), "Login with Apple"),
									authOutline(authBrandIcon(authGooglePath), "Login with Google"),
								),
								field.Separator(h.Class("*:data-[slot=field-separator-content]:bg-card"), "Or continue with"),
								authInputField("email", "email", "Email", "m@example.com"),
								authPasswordField("ml-auto text-sm underline-offset-4 hover:underline"),
								field.Field(
									authSubmit("Login"),
									authSwitch("text-center", "Don't have an account?", "Sign up"),
								),
							),
						),
					),
				),
				authTerms(),
			),
		),
	)
}

// authSocialRow is the three icon-only buttons of login-04 and signup-04.
func authSocialRow(verb string) h.Node {
	return field.Field(h.Class("grid grid-cols-3 gap-4"),
		authOutline(authBrandIcon(authApplePath), h.Span(h.Class("sr-only"), verb+" with Apple")),
		authOutline(authBrandIcon(authGooglePath), h.Span(h.Class("sr-only"), verb+" with Google")),
		authOutline(authBrandIcon(authMetaPath), h.Span(h.Class("sr-only"), verb+" with Meta")),
	)
}

// authImageCard is the card with a form and a cover image of login-04 and signup-04.
func authImageCard(form h.Node) h.Node {
	return h.Div(h.Class("flex min-h-svh flex-col items-center justify-center bg-muted p-6 md:p-10"),
		h.Div(h.Class("w-full max-w-sm md:max-w-4xl"),
			h.Div(h.Class("flex flex-col gap-6"),
				card.Card(h.Class("overflow-hidden p-0"),
					card.Content(h.Class("grid p-0 md:grid-cols-2"),
						form,
						h.Div(h.Class("relative hidden bg-muted md:block"), authImage()),
					),
				),
				authTerms(),
			),
		),
	)
}

// Login04 ports blocks/login-04.
func Login04() h.Node {
	return authImageCard(
		h.Form(h.Class("p-6 md:p-8"),
			field.Group(
				h.Div(h.Class("flex flex-col items-center gap-2 text-center"),
					h.H1(h.Class("text-2xl font-bold"), "Welcome back"),
					h.P(h.Class("text-balance text-muted-foreground"), "Login to your Acme Inc account"),
				),
				authInputField("email", "email", "Email", "m@example.com"),
				authPasswordField("ml-auto text-sm underline-offset-2 hover:underline"),
				field.Field(authSubmit("Login")),
				field.Separator(h.Class("*:data-[slot=field-separator-content]:bg-card"), "Or continue with"),
				authSocialRow("Login"),
				authSwitch("text-center", "Don't have an account?", "Sign up"),
			),
		),
	)
}

// authMinimalPage is the logo-only layout of login-05 and signup-05.
func authMinimalPage(question, link, submit string) h.Node {
	return h.Div(h.Class("flex min-h-svh flex-col items-center justify-center gap-6 bg-background p-6 md:p-10"),
		h.Div(h.Class("w-full max-w-sm"),
			h.Div(h.Class("flex flex-col gap-6"),
				h.Form(
					field.Group(
						authCenteredHeader(question, link),
						authInputField("email", "email", "Email", "m@example.com"),
						field.Field(authSubmit(submit)),
						field.Separator("Or"),
						field.Field(h.Class("grid gap-4 sm:grid-cols-2"),
							authOutline(authBrandIcon(authApplePath), "Continue with Apple"),
							authOutline(authBrandIcon(authGooglePath), "Continue with Google"),
						),
					),
				),
				authTerms(),
			),
		),
	)
}

// Login05 ports blocks/login-05.
func Login05() h.Node {
	return authMinimalPage("Don't have an account?", "Sign up", "Login")
}
