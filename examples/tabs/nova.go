package tabs

import (
	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/tabs"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The docs previews (ui/apps/v4/examples/radix/tabs-*.tsx) whose content
// differs from the new-york-v4 example of the same name, in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "tabs", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("tabs-demo", NovaDemo),
		nova("tabs-disabled", NovaDisabled),
		nova("tabs-icons", NovaIcons),
		nova("tabs-line", NovaLine),
		nova("tabs-vertical", NovaVertical),
		nova("tabs-rtl", NovaRtl),
	)
}

// novaPanels renders the four tabs of tabs-demo: {value, label, title,
// description, content}.
func novaPanels(rows [][5]string, dir string, attrs ...any) h.Node {
	var d []any
	if dir != "" {
		d = []any{h.Dir(dir)}
	}
	var triggers, panels []any
	for _, r := range rows {
		triggers = append(triggers, tabs.Trigger(tabs.Value(r[0]), r[1]))
		panels = append(panels, tabs.Content(tabs.Value(r[0]),
			card.Card(append(d,
				card.Header(card.Title(r[2]), card.Description(r[3])),
				card.Content(h.Class("text-sm text-muted-foreground"), r[4]),
			)...),
		))
	}
	kids := append(attrs, tabs.List(append(d, triggers...)...))
	return tabs.Tabs(append(kids, panels...)...)
}

func NovaDemo() h.Node {
	return novaPanels([][5]string{
		{"overview", "Overview", "Overview", "View your key metrics and recent project activity. Track progress across all your active projects.", "You have 12 active projects and 3 pending tasks."},
		{"analytics", "Analytics", "Analytics", "Track performance and user engagement metrics. Monitor trends and identify growth opportunities.", "Page views are up 25% compared to last month."},
		{"reports", "Reports", "Reports", "Generate and download your detailed reports. Export data in multiple formats for analysis.", "You have 5 reports ready and available to export."},
		{"settings", "Settings", "Settings", "Manage your account preferences and options. Customize your experience to fit your needs.", "Configure notifications, security, and themes."},
	}, "", tabs.Value("overview"), h.Class("w-[400px]"))
}

func NovaRtl() h.Node {
	return novaPanels([][5]string{
		{"overview", "نظرة عامة", "نظرة عامة", "عرض مقاييسك الرئيسية وأنشطة المشروع الأخيرة. تتبع التقدم عبر جميع مشاريعك النشطة.", "لديك ١٢ مشروعًا نشطًا و٣ مهام معلقة."},
		{"analytics", "التحليلات", "التحليلات", "تتبع مقاييس الأداء ومشاركة المستخدمين. راقب الاتجاهات وحدد فرص النمو.", "زادت مشاهدات الصفحة بنسبة ٢٥٪ مقارنة بالشهر الماضي."},
		{"reports", "التقارير", "التقارير", "إنشاء وتنزيل تقاريرك التفصيلية. تصدير البيانات بتنسيقات متعددة للتحليل.", "لديك ٥ تقارير جاهزة ومتاحة للتصدير."},
		{"settings", "الإعدادات", "الإعدادات", "إدارة تفضيلات حسابك وخياراته. تخصيص تجربتك لتناسب احتياجاتك.", "تكوين الإشعارات والأمان والسمات."},
	}, "rtl", tabs.Value("overview"), h.Class("w-full max-w-sm"), h.Dir("rtl"))
}

func NovaDisabled() h.Node {
	return tabs.Tabs(tabs.Value("home"),
		tabs.List(
			tabs.Trigger(tabs.Value("home"), "Home"),
			tabs.Trigger(tabs.Value("settings"), h.Disabled(), "Disabled"),
		),
	)
}

func NovaIcons() h.Node {
	return tabs.Tabs(tabs.Value("preview"),
		tabs.List(
			tabs.Trigger(tabs.Value("preview"), lucide.AppWindow(), "Preview"),
			tabs.Trigger(tabs.Value("code"), lucide.Code(), "Code"),
		),
	)
}

func NovaLine() h.Node {
	return tabs.Tabs(tabs.Value("overview"),
		tabs.List(tabs.ListVariantLine,
			tabs.Trigger(tabs.Value("overview"), "Overview"),
			tabs.Trigger(tabs.Value("analytics"), "Analytics"),
			tabs.Trigger(tabs.Value("reports"), "Reports"),
		),
	)
}

func NovaVertical() h.Node {
	return tabs.Tabs(tabs.Value("account"), tabs.OrientationVertical,
		tabs.List(
			tabs.Trigger(tabs.Value("account"), "Account"),
			tabs.Trigger(tabs.Value("password"), "Password"),
			tabs.Trigger(tabs.Value("notifications"), "Notifications"),
		),
	)
}
