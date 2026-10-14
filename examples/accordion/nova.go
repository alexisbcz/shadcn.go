package accordion

import (
	"shadcn.go/components/style"
	"shadcn.go/components/ui/accordion"
	"shadcn.go/components/ui/card"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The docs previews (ui/apps/v4/examples/radix/accordion-*.tsx) whose
// content differs from the new-york-v4 example of the same name, in
// radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "accordion", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("accordion-basic", NovaBasic),
		nova("accordion-borders", NovaBorders),
		nova("accordion-card", NovaCard),
		nova("accordion-demo", NovaDemo),
		nova("accordion-disabled", NovaDisabled),
		nova("accordion-multiple", NovaMultiple),
		nova("accordion-rtl", NovaRtl),
	)
}

// items renders {value, trigger, content} items, with extra item attrs.
func items(rows [][3]string, itemAttrs ...any) []any {
	var out []any
	for _, r := range rows {
		out = append(out, accordion.Item(append(append([]any{accordion.Value(r[0])}, itemAttrs...),
			accordion.Trigger(r[1]),
			accordion.Content(r[2]),
		)...))
	}
	return out
}

var basicItems = [][3]string{
	{"item-1", "How do I reset my password?", "Click on 'Forgot Password' on the login page, enter your email address, and we'll send you a link to reset your password. The link will expire in 24 hours."},
	{"item-2", "Can I change my subscription plan?", "Yes, you can upgrade or downgrade your plan at any time from your account settings. Changes will be reflected in your next billing cycle."},
	{"item-3", "What payment methods do you accept?", "We accept all major credit cards, PayPal, and bank transfers. All payments are processed securely through our payment partners."},
}

func NovaBasic() h.Node {
	return accordion.Accordion(append([]any{accordion.TypeSingle, accordion.Collapsible(true), accordion.Value("item-1"), h.Class("max-w-lg")}, items(basicItems)...)...)
}

func novaRtl() h.Node {
	return accordion.Accordion(append([]any{accordion.TypeSingle, accordion.Collapsible(true), accordion.Value("item-1"), h.Class("max-w-md")}, items([][3]string{
		{"item-1", "كيف يمكنني إعادة تعيين كلمة المرور؟", "انقر على 'نسيت كلمة المرور' في صفحة تسجيل الدخول، أدخل عنوان بريدك الإلكتروني، وسنرسل لك رابطًا لإعادة تعيين كلمة المرور. سينتهي صلاحية الرابط خلال 24 ساعة."},
		{"item-2", "هل يمكنني تغيير خطة الاشتراك الخاصة بي؟", "نعم، يمكنك ترقية أو تخفيض خطتك في أي وقت من إعدادات حسابك. ستظهر التغييرات في دورة الفوترة التالية."},
		{"item-3", "ما هي طرق الدفع التي تقبلونها؟", "نقبل جميع بطاقات الائتمان الرئيسية و PayPal والتحويلات المصرفية. تتم معالجة جميع المدفوعات بأمان من خلال شركاء الدفع لدينا."},
	})...)...)
}

func NovaBorders() h.Node {
	return accordion.Accordion(append([]any{accordion.TypeSingle, accordion.Collapsible(true), h.Class("max-w-lg rounded-lg border"), accordion.Value("billing")}, items([][3]string{
		{"billing", "How does billing work?", "We offer monthly and annual subscription plans. Billing is charged at the beginning of each cycle, and you can cancel anytime. All plans include automatic backups, 24/7 support, and unlimited team members."},
		{"security", "Is my data secure?", "Yes. We use end-to-end encryption, SOC 2 Type II compliance, and regular third-party security audits. All data is encrypted at rest and in transit using industry-standard protocols."},
		{"integration", "What integrations do you support?", "We integrate with 500+ popular tools including Slack, Zapier, Salesforce, HubSpot, and more. You can also build custom integrations using our REST API and webhooks."},
	}, h.Class("border-b px-4 last:border-b-0"))...)...)
}

func NovaCard() h.Node {
	return card.Card(h.Class("w-full max-w-sm"),
		card.Header(
			card.Title("Subscription & Billing"),
			card.Description("Common questions about your account, plans, payments and cancellations."),
		),
		card.Content(
			accordion.Accordion(append([]any{accordion.TypeSingle, accordion.Collapsible(true), accordion.Value("plans")}, items([][3]string{
				{"plans", "What subscription plans do you offer?", "We offer three subscription tiers: Starter ($9/month), Professional ($29/month), and Enterprise ($99/month). Each plan includes increasing storage limits, API access, priority support, and team collaboration features."},
				{"billing", "How does billing work?", "Billing occurs automatically at the start of each billing cycle. We accept all major credit cards, PayPal, and ACH transfers for enterprise customers. You'll receive an invoice via email after each payment."},
				{"cancel", "How do I cancel my subscription?", "You can cancel your subscription anytime from your account settings. There are no cancellation fees or penalties. Your access will continue until the end of your current billing period."},
			})...)...),
		),
	)
}

func NovaDemo() h.Node {
	return accordion.Accordion(append([]any{accordion.TypeSingle, accordion.Collapsible(true), accordion.Value("shipping"), h.Class("max-w-lg")}, items([][3]string{
		{"shipping", "What are your shipping options?", "We offer standard (5-7 days), express (2-3 days), and overnight shipping. Free shipping on international orders."},
		{"returns", "What is your return policy?", "Returns accepted within 30 days. Items must be unused and in original packaging. Refunds processed within 5-7 business days."},
		{"support", "How can I contact customer support?", "Reach us via email, live chat, or phone. We respond within 24 hours during business days."},
	})...)...)
}

func NovaDisabled() h.Node {
	return accordion.Accordion(accordion.TypeSingle, accordion.Collapsible(true), h.Class("w-full"),
		accordion.Item(accordion.Value("item-1"),
			accordion.Trigger("Can I access my account history?"),
			accordion.Content("Yes, you can view your complete account history including all transactions, plan changes, and support tickets in the Account History section of your dashboard."),
		),
		accordion.Item(accordion.Value("item-2"), accordion.Disabled(true),
			accordion.Trigger("Premium feature information"),
			accordion.Content("This section contains information about premium features. Upgrade your plan to access this content."),
		),
		accordion.Item(accordion.Value("item-3"),
			accordion.Trigger("How do I update my email address?"),
			accordion.Content("You can update your email address in your account settings. You'll receive a verification email at your new address to confirm the change."),
		),
	)
}

func NovaMultiple() h.Node {
	return accordion.Accordion(append([]any{accordion.TypeMultiple, h.Class("max-w-lg"), accordion.Value("notifications")}, items([][3]string{
		{"notifications", "Notification Settings", "Manage how you receive notifications. You can enable email alerts for updates or push notifications for mobile devices."},
		{"privacy", "Privacy & Security", "Control your privacy settings and security preferences. Enable two-factor authentication, manage connected devices, review active sessions, and configure data sharing preferences. You can also download your data or delete your account."},
		{"billing", "Billing & Subscription", "View your current plan, payment history, and upcoming invoices. Update your payment method, change your subscription tier, or cancel your subscription."},
	})...)...)
}

// NovaRtl uses the ui-rtl components (styles/radix-nova/ui-rtl).
func NovaRtl() h.Node { return style.ScopeRTL(novaRtl()) }
