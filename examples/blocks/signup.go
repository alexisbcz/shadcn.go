package blocks

import (
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/field"
	h "shadcn.go/libhtml"
)

const authEmailHint = "We'll use this to contact you. We will not share your email with anyone else."

// Signup01 ports blocks/signup-01.
func Signup01() h.Node {
	return h.Div(h.Class("flex min-h-svh w-full items-center justify-center p-6 md:p-10"),
		h.Div(h.Class("w-full max-w-sm"),
			card.Card(
				card.Header(
					card.Title("Create an account"),
					card.Description("Enter your information below to create your account"),
				),
				card.Content(
					h.Form(
						field.Group(
							authInputField("name", "text", "Full Name", "John Doe"),
							authInputField("email", "email", "Email", "m@example.com", field.Description(authEmailHint)),
							authInputField("password", "password", "Password", "", field.Description("Must be at least 8 characters long.")),
							authInputField("confirm-password", "password", "Confirm Password", "", field.Description("Please confirm your password.")),
							field.Group(
								field.Field(
									authSubmit("Create Account"),
									authOutline("Sign up with Google"),
									authSwitch("px-6 text-center", "Already have an account?", "Sign in"),
								),
							),
						),
					),
				),
			),
		),
	)
}

// Signup02 ports blocks/signup-02.
func Signup02() h.Node {
	return authSplitPage(
		h.Form(h.Class("flex flex-col gap-6"),
			field.Group(
				h.Div(h.Class("flex flex-col items-center gap-1 text-center"),
					h.H1(h.Class("text-2xl font-bold"), "Create your account"),
					h.P(h.Class("text-sm text-balance text-muted-foreground"), "Fill in the form below to create your account"),
				),
				authInputField("name", "text", "Full Name", "John Doe"),
				authInputField("email", "email", "Email", "m@example.com", field.Description(authEmailHint)),
				authInputField("password", "password", "Password", "", field.Description("Must be at least 8 characters long.")),
				authInputField("confirm-password", "password", "Confirm Password", "", field.Description("Please confirm your password.")),
				field.Field(authSubmit("Create Account")),
				field.Separator("Or continue with"),
				field.Field(
					authOutline(authBrandIcon(authGitHubPath), "Sign up with GitHub"),
					authSwitch("px-6 text-center", "Already have an account?", "Sign in"),
				),
			),
		),
	)
}

// Signup03 ports blocks/signup-03.
func Signup03() h.Node {
	return h.Div(h.Class("flex min-h-svh flex-col items-center justify-center gap-6 bg-muted p-6 md:p-10"),
		h.Div(h.Class("flex w-full max-w-sm flex-col gap-6"),
			authLogo("flex items-center gap-2 self-center font-medium"),
			h.Div(h.Class("flex flex-col gap-6"),
				card.Card(
					card.Header(h.Class("text-center"),
						card.Title(h.Class("text-xl"), "Create your account"),
						card.Description("Enter your email below to create your account"),
					),
					card.Content(
						h.Form(
							field.Group(
								authInputField("name", "text", "Full Name", "John Doe"),
								authInputField("email", "email", "Email", "m@example.com"),
								authPasswordPair(),
								field.Field(
									authSubmit("Create Account"),
									authSwitch("text-center", "Already have an account?", "Sign in"),
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

// Signup04 ports blocks/signup-04.
func Signup04() h.Node {
	return authImageCard(
		h.Form(h.Class("p-6 md:p-8"),
			field.Group(
				h.Div(h.Class("flex flex-col items-center gap-2 text-center"),
					h.H1(h.Class("text-2xl font-bold"), "Create your account"),
					h.P(h.Class("text-sm text-balance text-muted-foreground"), "Enter your email below to create your account"),
				),
				authInputField("email", "email", "Email", "m@example.com", field.Description(authEmailHint)),
				authPasswordPair(),
				field.Field(authSubmit("Create Account")),
				field.Separator(h.Class("*:data-[slot=field-separator-content]:bg-card"), "Or continue with"),
				authSocialRow("Sign up"),
				authSwitch("text-center", "Already have an account?", "Sign in"),
			),
		),
	)
}

// Signup05 ports blocks/signup-05.
func Signup05() h.Node {
	return authMinimalPage("Already have an account?", "Sign in", "Create Account")
}
