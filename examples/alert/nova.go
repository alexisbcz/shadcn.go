package alert

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/alert"
	"shadcn.go/components/ui/button"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the alert page (ui/apps/v4/examples/radix/alert-*.tsx),
// in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "alert", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("alert-action", NovaAction),
		nova("alert-basic", NovaBasic),
		nova("alert-colors", NovaColors),
		nova("alert-demo", NovaDemo),
		nova("alert-destructive", NovaDestructive),
		nova("alert-rtl", NovaRtl),
	)
}

func NovaAction() h.Node {
	return alert.Alert(h.Class("max-w-md"),
		alert.Title("Dark mode is now available"),
		alert.Description("Enable it under your profile settings to get started."),
		alert.Action(button.Button(button.SizeXs, button.VariantDefault, "Enable")),
	)
}

func NovaBasic() h.Node {
	return alert.Alert(h.Class("max-w-md"),
		lucide.CircleCheck(),
		alert.Title("Account updated successfully"),
		alert.Description("Your profile information has been saved. Changes will be reflected immediately."),
	)
}

func NovaColors() h.Node {
	return alert.Alert(h.Class("max-w-md border-amber-200 bg-amber-50 text-amber-900 dark:border-amber-900 dark:bg-amber-950 dark:text-amber-50"),
		lucide.TriangleAlert(),
		alert.Title("Your subscription will expire in 3 days."),
		alert.Description("Renew now to avoid service interruption or upgrade to a paid plan to continue using the service."),
	)
}

func novaDemo(dir string, paymentTitle, paymentDescription, featureTitle, featureDescription string) h.Node {
	return h.Div(h.Class("grid w-full max-w-md items-start gap-4"), h.If(dir != "", h.Dir(dir)),
		alert.Alert(lucide.CircleCheck(), alert.Title(paymentTitle), alert.Description(paymentDescription)),
		alert.Alert(lucide.Info(), alert.Title(featureTitle), alert.Description(featureDescription)),
	)
}

func NovaDemo() h.Node {
	return novaDemo("", "Payment successful",
		"Your payment of $29.99 has been processed. A receipt has been sent to your email address.",
		"New feature available",
		"We've added dark mode support. You can enable it in your account settings.")
}

func NovaDestructive() h.Node {
	return alert.Alert(alert.VariantDestructive, h.Class("max-w-md"),
		lucide.CircleAlert(),
		alert.Title("Payment failed"),
		alert.Description("Your payment could not be processed. Please check your payment method and try again."),
	)
}

func NovaRtl() h.Node {
	return style.ScopeRTL(novaDemo("rtl", "تم الدفع بنجاح",
		"تمت معالجة دفعتك البالغة 29.99 دولارًا. تم إرسال إيصال إلى عنوان بريدك الإلكتروني.",
		"ميزة جديدة متاحة",
		"لقد أضفنا دعم الوضع الداكن. يمكنك تفعيله في إعدادات حسابك."))
}
