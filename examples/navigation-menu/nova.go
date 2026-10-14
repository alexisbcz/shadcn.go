package navigationmenu

import (
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	nm "shadcn.go/components/ui/navigation-menu"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// Docs previews of the navigation-menu page
// (ui/apps/v4/examples/radix/navigation-menu-*.tsx), rendered in radix-nova.
func init() {
	nova := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: "nova--" + name, Component: "navigation-menu", Render: registry.Static(func() h.Node {
			return style.Scope(style.Nova, fn())
		})}
	}
	registry.Register(
		nova("navigation-menu-demo", func() h.Node { return novaMenu(novaEnglish, "", nil) }),
		nova("navigation-menu-rtl", func() h.Node { return style.ScopeRTL(novaMenu(novaArabic, "rtl", h.Data("lang", "ar"))) }),
	)
}

type novaTexts struct {
	gettingStarted, introduction, introductionDesc, installation, installationDesc, typography, typographyDesc, components, docs string
	items                                                                                                                         [6][2]string
}

var novaEnglish = novaTexts{
	"Getting started", "Introduction", "Re-usable components built with Tailwind CSS.",
	"Installation", "How to install dependencies and structure your app.",
	"Typography", "Styles for headings, paragraphs, lists...etc", "Components", "Docs",
	[6][2]string{
		{components[0].title, components[0].description}, {components[1].title, components[1].description},
		{components[2].title, components[2].description}, {components[3].title, components[3].description},
		{components[4].title, components[4].description}, {components[5].title, components[5].description},
	},
}

var novaArabic = novaTexts{
	"البدء", "مقدمة", "مكونات قابلة لإعادة الاستخدام مبنية باستخدام Tailwind CSS.",
	"التثبيت", "كيفية تثبيت التبعيات وتنظيم تطبيقك.",
	"الطباعة", "أنماط للعناوين والفقرات والقوائم...إلخ", "المكونات", "الوثائق",
	[6][2]string{
		{"حوار التنبيه", "حوار نافذة يقطع المستخدم بمحتوى مهم ويتوقع استجابة."},
		{"بطاقة التحويم", "للمستخدمين المبصرين لمعاينة المحتوى المتاح خلف الرابط."},
		{"التقدم", "يعرض مؤشرًا يوضح تقدم إتمام المهمة، عادةً يتم عرضه كشريط تقدم."},
		{"منطقة التمرير", "يفصل المحتوى بصريًا أو دلاليًا."},
		{"التبويبات", "مجموعة من أقسام المحتوى المتعددة الطبقات—المعروفة بألواح التبويب—التي يتم عرضها واحدة في كل مرة."},
		{"تلميح", "نافذة منبثقة تعرض معلومات متعلقة بعنصر عندما يتلقى العنصر التركيز على لوحة المفاتيح أو عند تحويم الماوس فوقه."},
	},
}

func novaListItem(href, title, desc string) h.Node {
	return h.Li(
		nm.Link(lib.AsChild, h.A(h.Href(href),
			h.Div(h.Class("flex flex-col gap-1 text-sm"),
				h.Div(h.Class("leading-none font-medium"), title),
				h.Div(h.Class("line-clamp-2 text-muted-foreground"), desc),
			),
		)),
	)
}

// novaMenu is examples/radix/navigation-menu-demo.tsx, and with dir "rtl"
// navigation-menu-rtl.tsx in Arabic (lang is its data-lang attribute).
func novaMenu(t novaTexts, dir string, lang any) h.Node {
	var dirOpts []any
	if dir != "" {
		dirOpts = append(dirOpts, h.Dir(dir), lang)
	}
	comps := make([]any, len(t.items))
	for i, c := range t.items {
		comps[i] = novaListItem(components[i].href, c[0], c[1])
	}
	root := []any{}
	if dir != "" {
		root = append(root, nm.Dir(dir))
	}
	return nm.NavigationMenu(append(root,
		nm.List(
			nm.Item(
				nm.Trigger(t.gettingStarted),
				nm.Content(append(dirOpts,
					h.Ul(h.Class("w-96"),
						novaListItem("/docs", t.introduction, t.introductionDesc),
						novaListItem("/docs/installation", t.installation, t.installationDesc),
						novaListItem("/docs/primitives/typography", t.typography, t.typographyDesc),
					))...,
				),
			),
			nm.Item(h.Class("hidden md:flex"),
				nm.Trigger(t.components),
				nm.Content(append(dirOpts,
					h.Ul(append([]any{h.Class("grid w-[400px] gap-2 md:w-[500px] md:grid-cols-2 lg:w-[600px]")}, comps...)...))...,
				),
			),
			nm.Item(
				nm.Link(lib.AsChild, h.Class(nm.TriggerStyleFor(style.Nova)), h.If(dir != "", lang), h.A(h.Href("/docs"), t.docs)),
			),
		),
	)...)
}
