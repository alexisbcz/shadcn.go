package site

import (
	_ "embed"
	"fmt"
	"net/http"
	"strings"

	"shadcn.go/components/icons/lucide"
	"shadcn.go/components/icons/tabler"
	"shadcn.go/components/lib"
	"shadcn.go/components/style"
	"shadcn.go/components/ui/avatar"
	"shadcn.go/components/ui/badge"
	"shadcn.go/components/ui/button"
	buttongroup "shadcn.go/components/ui/button-group"
	"shadcn.go/components/ui/card"
	"shadcn.go/components/ui/checkbox"
	"shadcn.go/components/ui/command"
	"shadcn.go/components/ui/direction"
	dm "shadcn.go/components/ui/dropdown-menu"
	"shadcn.go/components/ui/empty"
	"shadcn.go/components/ui/field"
	"shadcn.go/components/ui/input"
	ig "shadcn.go/components/ui/input-group"
	"shadcn.go/components/ui/item"
	"shadcn.go/components/ui/label"
	"shadcn.go/components/ui/popover"
	rg "shadcn.go/components/ui/radio-group"
	"shadcn.go/components/ui/runtime"
	selectui "shadcn.go/components/ui/select"
	"shadcn.go/components/ui/separator"
	"shadcn.go/components/ui/slider"
	"shadcn.go/components/ui/spinner"
	switchui "shadcn.go/components/ui/switch"
	"shadcn.go/components/ui/textarea"
	"shadcn.go/components/ui/tooltip"
	h "shadcn.go/libhtml"
)

// The RTL example: app/(app)/examples/rtl/components/*.tsx, built on
// styles/base-nova/ui-rtl (here: the components in nova, ui-rtl variant).
// DropdownMenu and Select take the direction as an option (they do not
// read direction.Provider yet). The reference places most of its popups
// with physical (LTR) alignment while the Radix popper flips start/end for
// rtl content, so those use the opposite align (noted where they do).
// The page renders the grid once per language (Arabic, Hebrew), each with
// its language selector; static/examples-rtl.js shows the selected one, as
// the LanguageProvider re-renders upstream.

//go:embed static/examples-rtl.js
var examplesRTLScript string

func init() { runtime.Module("example-rtl", examplesRTLScript) }

// ExampleRTL serves /examples/rtl (app/(app)/examples/rtl).
func (s *Site) ExampleRTL(r *http.Request) (h.Node, error) {
	return s.examplesLayoutMeta(r, Meta{Title: "RTL", PlainTitle: true, Description: "RTL example page with right-to-left language support.", Canonical: "/examples/rtl"},
		style.Scope(style.Nova, style.ScopeRTL(direction.Provider(direction.RTL,
			h.Div(h.Class("contents"), h.Data("rtl-example", ""),
				rtlComponents("ar"),
				rtlComponents("he"),
			),
		))),
	), nil
}

// rtlLanguageSelector is components/language-selector.tsx's
// LanguageSelector (base-nova select) as the RTL grid places it.
func rtlLanguageSelector(lang string) h.Node {
	return selectui.Select(selectui.DefaultValue(lang),
		selectui.Trigger(selectui.Value()).Size(selectui.SizeSm).
			Class("w-36 absolute -top-12 right-52 hidden h-8! data-[size=sm]:rounded-lg lg:flex").
			Attr("dir", "ltr").Data("name", "language-selector"),
		selectui.Content(h.Attr("dir", "ltr"), h.Class("data-open:animate-none data-closed:animate-none"),
			selectui.Group(
				selectui.Item("Arabic (العربية)").Value("ar"),
				selectui.Item("Hebrew (עברית)").Value("he"),
			),
		),
	)
}

func rtlComponents(lang string) h.Node {
	col := "flex flex-col gap-6 *:[div]:w-full *:[div]:max-w-full"
	sep := "إعدادات المظهر"
	if lang == "he" {
		sep = "הגדרות מראה"
	}
	grid := h.Div(h.Class("relative grid gap-8 p-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 xl:gap-6 2xl:gap-8"),
		h.Dir("rtl"), h.Data("lang", lang), h.Data("slot", "rtl-components"),
		rtlLanguageSelector(lang),
		h.Div(h.Class(col), rtlFieldDemo(lang)),
		h.Div(h.Class(col),
			rtlEmptyAvatarGroup(lang),
			rtlSpinnerBadge(lang),
			rtlButtonGroupInputGroup(lang),
			rtlFieldSlider(lang),
			rtlInputGroupDemo(lang),
		),
		h.Div(h.Class(col),
			rtlInputGroupButton(lang),
			rtlItemDemo(lang),
			field.Separator(h.Class("my-4"), sep),
			rtlAppearanceSettings(lang),
		),
		h.Div(h.Class("order-first flex flex-col gap-6 lg:hidden xl:order-last xl:flex *:[div]:w-full *:[div]:max-w-full"),
			rtlNotionPromptForm(lang),
			rtlButtonGroupDemo(lang),
			rtlFieldCheckbox(lang),
			h.Div(h.Class("flex justify-between gap-4"),
				rtlButtonGroupNested(lang),
				rtlButtonGroupPopover(lang),
			),
			rtlFieldHear(lang),
			rtlSpinnerEmpty(lang),
		),
	)
	if lang != "ar" {
		grid.Children(h.Hidden())
	}
	return grid
}

// rtlTr picks a component's translations.
func rtlTr(lang string, tr map[string]map[string]string) map[string]string { return tr[lang] }

// rtlID keeps ids unique across the two language grids.
func rtlID(lang, id string) string {
	if lang == "ar" {
		return id
	}
	return id + "-" + lang
}

// rtlNum is Intl.NumberFormat(locale, {useGrouping: false}) for the two
// locales: ar-SA uses Arabic-Indic digits, he-IL Latin ones.
func rtlNum(lang string, n, minDigits int) string {
	s := fmt.Sprintf("%0*d", minDigits, n)
	if lang != "ar" {
		return s
	}
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			r = '٠' + (r - '0')
		}
		b.WriteRune(r)
	}
	return b.String()
}

func rtlFieldDemo(lang string) h.Node {
	t := rtlTr(lang, map[string]map[string]string{
		"ar": {"paymentMethod": "طريقة الدفع", "secureEncrypted": "جميع المعاملات آمنة ومشفرة", "nameOnCard": "الاسم على البطاقة", "namePlaceholder": "أحمد محمد", "cardNumber": "رقم البطاقة", "cardDescription": "أدخل رقمك المكون من 16 رقمًا.", "cvv": "رمز الأمان", "month": "الشهر", "year": "السنة", "billingAddress": "عنوان الفواتير", "billingDescription": "عنوان الفواتير المرتبط بطريقة الدفع الخاصة بك", "sameAsShipping": "نفس عنوان الشحن", "comments": "تعليقات", "commentsPlaceholder": "أضف أي تعليقات إضافية", "submit": "إرسال", "cancel": "إلغاء"},
		"he": {"paymentMethod": "אמצעי תשלום", "secureEncrypted": "כל העסקאות מאובטחות ומוצפנות", "nameOnCard": "שם על הכרטיס", "namePlaceholder": "ישראל ישראלי", "cardNumber": "מספר כרטיס", "cardDescription": "הזן את המספר בן 16 הספרות שלך.", "cvv": "קוד אבטחה", "month": "חודש", "year": "שנה", "billingAddress": "כתובת לחיוב", "billingDescription": "כתובת החיוב המשויכת לאמצעי התשלום שלך", "sameAsShipping": "זהה לכתובת המשלוח", "comments": "הערות", "commentsPlaceholder": "הוסף הערות נוספות", "submit": "שלח", "cancel": "ביטול"},
	})
	id := func(s string) string { return rtlID(lang, s) }
	var months, years []any
	for i := 1; i <= 12; i++ {
		months = append(months, selectui.Item(rtlNum(lang, i, 2)).Value(fmt.Sprintf("%02d", i)))
	}
	for y := 2024; y < 2030; y++ {
		years = append(years, selectui.Item(rtlNum(lang, y, 1)).Value(fmt.Sprint(y)))
	}
	card := strings.Join([]string{rtlNum(lang, 1234, 1), rtlNum(lang, 5678, 1), rtlNum(lang, 9012, 1), rtlNum(lang, 3456, 1)}, " ")
	sel := func(sid, placeholder string, items []any) h.Node {
		return selectui.Select(selectui.DefaultValue(""), selectui.Dir("rtl"),
			selectui.Trigger(h.ID(sid), selectui.Value(selectui.Placeholder(placeholder))),
			selectui.Content(h.Data("lang", lang), h.Dir("rtl"), selectui.Group(items...)),
		)
	}
	return h.Div(h.Dir("rtl"), h.Class("w-full max-w-md rounded-lg border p-6"),
		h.Form(
			field.Group(
				field.Set(
					field.Legend(t["paymentMethod"]),
					field.Description(t["secureEncrypted"]),
					field.Group(
						field.Field(
							field.Label(h.For(id("rtl-card-name")), t["nameOnCard"]),
							input.Input(h.ID(id("rtl-card-name")), h.Placeholder(t["namePlaceholder"]), h.Required()),
						),
						h.Div(h.Class("grid grid-cols-3 gap-4"),
							field.Field(h.Class("col-span-2"),
								field.Label(h.For(id("rtl-card-number")), t["cardNumber"]),
								input.Input(h.ID(id("rtl-card-number")), h.Placeholder(card), h.Required()),
								field.Description(t["cardDescription"]),
							),
							field.Field(h.Class("col-span-1"),
								field.Label(h.For(id("rtl-cvv")), t["cvv"]),
								input.Input(h.ID(id("rtl-cvv")), h.Placeholder(rtlNum(lang, 123, 1)), h.Required()),
							),
						),
						h.Div(h.Class("grid grid-cols-2 gap-4"),
							field.Field(
								field.Label(h.For(id("rtl-exp-month")), t["month"]),
								sel(id("rtl-exp-month"), "MM", months),
							),
							field.Field(
								field.Label(h.For(id("rtl-exp-year")), t["year"]),
								sel(id("rtl-exp-year"), "YYYY", years),
							),
						),
					),
				),
				field.Separator(),
				field.Set(
					field.Legend(t["billingAddress"]),
					field.Description(t["billingDescription"]),
					field.Group(
						field.Field(field.Horizontal,
							checkbox.Checkbox(h.ID(id("rtl-same-as-shipping")), checkbox.Checked(true)),
							field.Label(h.For(id("rtl-same-as-shipping")), h.Class("font-normal"), t["sameAsShipping"]),
						),
					),
				),
				field.Separator(),
				field.Set(
					field.Group(
						field.Field(
							field.Label(h.For(id("rtl-comments")), t["comments"]),
							textarea.Textarea(h.ID(id("rtl-comments")), h.Placeholder(t["commentsPlaceholder"]), h.Class("resize-none")),
						),
					),
				),
				field.Field(field.Horizontal,
					button.Button(h.Type("submit"), t["submit"]),
					button.Button(button.VariantOutline, h.Type("button"), t["cancel"]),
				),
			),
		),
	)
}

func rtlEmptyAvatarGroup(lang string) h.Node {
	t := rtlTr(lang, map[string]map[string]string{
		"ar": {"title": "لا يوجد أعضاء في الفريق", "description": "قم بدعوة فريقك للتعاون في هذا المشروع.", "invite": "دعوة أعضاء"},
		"he": {"title": "אין חברי צוות", "description": "הזמן את הצוות שלך לשתף פעולה בפרויקט זה.", "invite": "הזמן חברים"},
	})
	av := func(src, alt, fb string) h.Node {
		return avatar.Avatar(avatar.Image(h.Src(src), h.Alt(alt)), avatar.Fallback(fb))
	}
	return empty.Empty(h.Class("flex-none border py-10"), h.Dir("rtl"),
		empty.Header(
			empty.Media(
				avatar.Group(h.Class("grayscale"),
					av("https://github.com/shadcn.png", "@shadcn", "CN"),
					av("https://github.com/maxleiter.png", "@maxleiter", "LR"),
					av("https://github.com/evilrabbit.png", "@evilrabbit", "ER"),
				),
			),
			empty.Title(t["title"]),
			empty.Description(t["description"]),
		),
		empty.Content(button.Button(button.SizeSm, lucide.Plus(), t["invite"])),
	)
}

func rtlSpinnerBadge(lang string) h.Node {
	t := rtlTr(lang, map[string]map[string]string{
		"ar": {"syncing": "جارٍ المزامنة", "updating": "جارٍ التحديث", "loading": "جارٍ التحميل"},
		"he": {"syncing": "מסנכרן", "updating": "מעדכן", "loading": "טוען"},
	})
	return h.Div(h.Dir("rtl"), h.Class("flex items-center gap-2"),
		badge.Badge(spinner.Spinner(), t["syncing"]),
		badge.Badge(badge.VariantSecondary, spinner.Spinner(), t["updating"]),
		badge.Badge(badge.VariantOutline, spinner.Spinner(), t["loading"]),
	)
}

func rtlButtonGroupInputGroup(lang string) h.Node {
	t := rtlTr(lang, map[string]map[string]string{
		"ar": {"add": "إضافة", "voicePlaceholder": "سجل وأرسل صوتًا...", "messagePlaceholder": "أرسل رسالة...", "voiceMode": "الوضع الصوتي"},
		"he": {"add": "הוסף", "voicePlaceholder": "הקלט ושלח אודיו...", "messagePlaceholder": "שלח הודעה...", "voiceMode": "מצב קולי"},
	})
	return buttongroup.ButtonGroup(h.Dir("rtl"),
		buttongroup.ButtonGroup(
			button.Button(button.VariantOutline, button.SizeIcon, h.Aria("label", t["add"]), lucide.Plus()),
		),
		buttongroup.ButtonGroup(h.Class("flex-1"),
			ig.InputGroup(
				ig.Input(h.Placeholder(t["messagePlaceholder"]),
					h.Data("rtl-voice-off", t["messagePlaceholder"]), h.Data("rtl-voice-on", t["voicePlaceholder"])),
				ig.Addon(ig.AlignInlineEnd,
					tooltip.Tooltip(
						tooltip.Trigger(lib.AsChild,
							ig.Button(ig.ButtonSizeIconXs, h.Data("active", "false"), h.Aria("pressed", "false"),
								h.Class("data-[active=true]:bg-primary data-[active=true]:text-primary-foreground"),
								h.Aria("label", t["voiceMode"]), h.Data("rtl-voice", ""),
								lucide.AudioLines()),
						),
						tooltip.Content(t["voiceMode"]),
					),
				),
			),
		),
	)
}

func rtlFieldSlider(lang string) h.Node {
	t := rtlTr(lang, map[string]map[string]string{
		"ar": {"locale": "ar-SA", "title": "نطاق السعر", "description": "حدد نطاق ميزانيتك", "currency": "﷼"},
		"he": {"locale": "he-IL", "title": "טווח מחירים", "description": "הגדר את טווח התקציב שלך", "currency": "₪"},
	})
	return field.Field(h.Dir("rtl"), h.Data("rtl-slider", t["locale"]),
		field.Title(t["title"]),
		field.Description(
			t["description"], " (", t["currency"],
			h.Span(h.Class("font-medium tabular-nums"), rtlNum(lang, 200, 1)),
			" ", "- ",
			h.Span(h.Class("font-medium tabular-nums"), rtlNum(lang, 800, 1)),
			").",
		),
		slider.Slider(slider.DefaultValue(200, 800), slider.Max(1000), slider.Min(0), slider.Step(10),
			h.Class("mt-2 w-full"), h.Aria("label", t["title"])),
	)
}

func rtlInputGroupDemo(lang string) h.Node {
	t := rtlTr(lang, map[string]map[string]string{
		"ar": {"search": "بحث...", "results": "12 نتيجة", "tooltipContent": "هذا محتوى في تلميح.", "askSearchChat": "اسأل، ابحث أو تحدث...", "add": "إضافة", "auto": "تلقائي", "agent": "وكيل", "manual": "يدوي", "used": "52% مستخدم", "send": "إرسال"},
		"he": {"search": "חיפוש...", "results": "12 תוצאות", "tooltipContent": "זה תוכן בטולטיפ.", "askSearchChat": "שאל, חפש או שוחח...", "add": "הוסף", "auto": "אוטומטי", "agent": "סוכן", "manual": "ידני", "used": "52% בשימוש", "send": "שלח"},
	})
	return h.Div(h.Dir("rtl"), h.Class("grid w-full max-w-sm gap-6"),
		ig.InputGroup(
			ig.Input(h.Placeholder(t["search"])),
			ig.Addon(lucide.Search()),
			ig.Addon(ig.AlignInlineEnd, t["results"]),
		),
		ig.InputGroup(
			ig.Input(h.Placeholder("example.com")),
			ig.Addon(ig.AlignInlineEnd,
				tooltip.Tooltip(
					tooltip.Trigger(lib.AsChild,
						ig.Button(h.Class("rounded-full"), ig.ButtonSizeIconXs, h.Aria("label", t["add"]), tabler.IconInfoCircle())),
					tooltip.Content(t["tooltipContent"]),
				),
			),
		),
		ig.InputGroup(
			ig.Textarea(h.Placeholder(t["askSearchChat"])),
			ig.Addon(ig.AlignBlockEnd,
				ig.Button(button.VariantOutline, h.Class("rounded-full"), ig.ButtonSizeIconXs, h.Aria("label", t["add"]), tabler.IconPlus()),
				dm.DropdownMenu(dm.Dir("rtl"),
					dm.Trigger(lib.AsChild, ig.Button(button.VariantGhost, tabler.IconChevronDown())),
					// Upstream align is the opposite (see the file comment).
					dm.Content(dm.SideTop, dm.AlignEnd,
						dm.Item(t["auto"]),
						dm.Item(t["agent"]),
						dm.Item(t["manual"]),
					),
				),
				ig.Text(h.Class("ms-auto"), t["used"]),
				separator.Separator(separator.Vertical, h.Class("h-4!")),
				ig.Button(button.VariantDefault, h.Class("rounded-full"), ig.ButtonSizeIconXs,
					lucide.ArrowUp(), h.Span(h.Class("sr-only"), t["send"])),
			),
		),
		ig.InputGroup(
			ig.Input(h.Placeholder("shadcn")),
			ig.Addon(ig.AlignInlineEnd,
				h.Div(h.Class("flex size-4 items-center justify-center rounded-full bg-primary text-foreground"),
					tabler.IconCheck(h.Class("size-3 text-white"))),
			),
		),
	)
}

func rtlInputGroupButton(lang string) h.Node {
	t := rtlTr(lang, map[string]map[string]string{
		"ar": {"inputLabel": "السعر", "info": "معلومات", "priceInfo": "أدخل السعر بالريال السعودي.", "priceDescription": "سيتم تحويل السعر تلقائياً.", "favorite": "مفضل", "currency": "ر.س"},
		"he": {"inputLabel": "מחיר", "info": "מידע", "priceInfo": "הזן את המחיר בשקלים.", "priceDescription": "המחיר יומר אוטומטית.", "favorite": "מועדף", "currency": "₪"},
	})
	id := rtlID(lang, "input-secure-rtl")
	return h.Div(h.Dir("rtl"), h.Class("grid w-full max-w-sm gap-6"),
		label.Label(h.For(id), h.Class("sr-only"), t["inputLabel"]),
		ig.InputGroup(h.Class("[--radius:9999px]"),
			ig.Input(h.ID(id), h.Class("pr-0.5!")),
			ig.Addon(
				popover.Popover(
					popover.Trigger(lib.AsChild,
						ig.Button(button.VariantSecondary, ig.ButtonSizeIconXs, h.Aria("label", t["info"]), tabler.IconInfoCircle())),
					// Upstream align is the opposite (see the file comment).
					popover.Content(popover.AlignStart, popover.AlignOffset(10),
						h.Class("flex flex-col gap-1 rounded-xl text-sm"), h.Data("lang", lang), h.Dir("rtl"),
						h.P(h.Class("font-medium"), t["priceInfo"]),
						h.P(t["priceDescription"]),
					),
				),
			),
			ig.Addon(h.Class("text-muted-foreground"), t["currency"]),
			ig.Addon(ig.AlignInlineEnd,
				ig.Button(ig.ButtonSizeIconXs, h.Aria("label", t["favorite"]), h.Data("rtl-favorite", ""),
					tabler.IconStar(h.Data("favorite", "false"), h.Class("data-[favorite=true]:fill-primary data-[favorite=true]:stroke-primary"))),
			),
		),
	)
}

func rtlItemDemo(lang string) h.Node {
	t := rtlTr(lang, map[string]map[string]string{
		"ar": {"twoFactor": "المصادقة الثنائية", "twoFactorDescription": "التحقق عبر البريد الإلكتروني أو رقم الهاتف.", "enable": "تفعيل", "verified": "تم التحقق من ملفك الشخصي."},
		"he": {"twoFactor": "אימות דו-שלבי", "twoFactorDescription": "אמת באמצעות אימייל או מספר טלפון.", "enable": "הפעל", "verified": "הפרופיל שלך אומת."},
	})
	return h.Div(h.Dir("rtl"), h.Class("flex w-full max-w-md flex-col gap-6"),
		item.Item(item.VariantOutline,
			item.Content(
				item.Title(t["twoFactor"]),
				item.Description(h.Class("text-pretty xl:hidden 2xl:block"), t["twoFactorDescription"]),
			),
			item.Actions(button.Button(button.SizeSm, t["enable"])),
		),
		item.Item(item.VariantOutline, item.SizeSm,
			item.Media(lucide.BadgeCheck(h.Class("size-5"))),
			item.Content(item.Title(t["verified"])),
			item.Actions(lucide.ChevronRight(h.Class("size-4 rtl:rotate-180"))),
		),
	)
}

func rtlAppearanceSettings(lang string) h.Node {
	t := rtlTr(lang, map[string]map[string]string{
		"ar": {"computeEnvironment": "بيئة الحوسبة", "computeDescription": "اختر بيئة الحوسبة لمجموعتك.", "kubernetes": "كوبرنيتس", "kubernetesDescription": "تشغيل أحمال عمل GPU على مجموعة مُهيأة بـ K8s. هذا هو الافتراضي.", "virtualMachine": "جهاز افتراضي", "vmDescription": "الوصول إلى مجموعة VM مُهيأة لتشغيل أحمال العمل. (قريبًا)", "numberOfGpus": "عدد وحدات GPU", "gpuDescription": "يمكنك إضافة المزيد لاحقًا.", "decrement": "إنقاص", "increment": "زيادة", "wallpaperTinting": "تلوين الخلفية", "wallpaperDescription": "السماح بتلوين الخلفية."},
		"he": {"computeEnvironment": "סביבת מחשוב", "computeDescription": "בחר את סביבת המחשוב לאשכול שלך.", "kubernetes": "קוברנטיס", "kubernetesDescription": "הפעל עומסי עבודה של GPU באשכול מוגדר K8s. זו ברירת המחדל.", "virtualMachine": "מכונה וירטואלית", "vmDescription": "גש לאשכול VM מוגדר להפעלת עומסי עבודה. (בקרוב)", "numberOfGpus": "מספר GPUs", "gpuDescription": "תוכל להוסיף עוד מאוחר יותר.", "decrement": "הפחת", "increment": "הגדל", "wallpaperTinting": "צביעת טפט", "wallpaperDescription": "אפשר לטפט להיצבע."},
	})
	id := func(s string) string { return rtlID(lang, s) }
	choice := func(cid, value, title, desc string) h.Node {
		return field.Label(h.For(cid),
			field.Field(field.Horizontal,
				field.Content(field.Title(title), field.Description(desc)),
				rg.Item(rg.Value(value), h.ID(cid), h.Aria("label", title)),
			),
		)
	}
	return h.Div(h.Dir("rtl"),
		field.Set(
			field.Group(
				field.Set(
					field.Legend(t["computeEnvironment"]),
					field.Description(t["computeDescription"]),
					rg.RadioGroup(rg.Value("kubernetes"),
						choice(id("rtl-kubernetes"), "kubernetes", t["kubernetes"], t["kubernetesDescription"]),
						choice(id("rtl-vm"), "vm", t["virtualMachine"], t["vmDescription"]),
					),
				),
				field.Separator(),
				field.Field(field.Horizontal,
					field.Content(
						field.Label(h.For(id("rtl-gpu-count")), t["numberOfGpus"]),
						field.Description(t["gpuDescription"]),
					),
					buttongroup.ButtonGroup(h.Data("rtl-gpu", ""),
						input.Input(h.ID(id("rtl-gpu-count")), h.Value("8"), h.Attr("size", "3"), h.Class("h-7 w-14! font-mono"), h.Attr("maxlength", "3")),
						button.Button(button.VariantOutline, button.SizeIconSm, h.Type("button"), h.Aria("label", t["decrement"]), h.Data("rtl-gpu-step", "-1"), rtlIconMinus()),
						button.Button(button.VariantOutline, button.SizeIconSm, h.Type("button"), h.Aria("label", t["increment"]), h.Data("rtl-gpu-step", "1"), tabler.IconPlus()),
					),
				),
				field.Separator(),
				field.Field(field.Horizontal,
					field.Content(
						field.Label(h.For(id("rtl-tinting")), t["wallpaperTinting"]),
						field.Description(t["wallpaperDescription"]),
					),
					switchui.Switch(h.ID(id("rtl-tinting")), switchui.Checked(true)),
				),
			),
		),
	)
}

func rtlNotionPromptForm(lang string) h.Node {
	t := rtlTr(lang, map[string]map[string]string{
		"ar": {"prompt": "الأمر", "placeholder": "اسأل، ابحث، أو أنشئ أي شيء...", "addContext": "أضف سياق", "mentionTooltip": "اذكر شخصًا أو صفحة أو تاريخًا", "attachFile": "إرفاق ملف", "selectModel": "اختر نموذج الذكاء الاصطناعي", "webSearch": "البحث على الويب", "appsIntegrations": "التطبيقات والتكاملات", "allSourcesAccess": "جميع المصادر التي يمكنني الوصول إليها", "findKnowledge": "ابحث أو استخدم المعرفة في...", "noKnowledgeFound": "لم يتم العثور على معرفة", "helpCenter": "مركز المساعدة", "connectApps": "ربط التطبيقات", "searchSourcesNote": "سنبحث فقط في المصادر المحددة هنا.", "send": "إرسال", "allSources": "جميع المصادر", "auto": "تلقائي", "workspace": "مساحة العمل"},
		"he": {"prompt": "פקודה", "placeholder": "שאל, חפש, או צור משהו...", "addContext": "הוסף הקשר", "mentionTooltip": "הזכר אדם, עמוד או תאריך", "attachFile": "צרף קובץ", "selectModel": "בחר מודל AI", "webSearch": "חיפוש באינטרנט", "appsIntegrations": "אפליקציות ואינטגרציות", "allSourcesAccess": "כל המקורות שיש לי גישה אליהם", "findKnowledge": "מצא או השתמש בידע ב...", "noKnowledgeFound": "לא נמצא ידע", "helpCenter": "מרכז עזרה", "connectApps": "חבר אפליקציות", "searchSourcesNote": "נחפש רק במקורות שנבחרו כאן.", "send": "שלח", "allSources": "כל המקורות", "auto": "אוטומטי", "workspace": "סביבת עבודה"},
	})
	id := func(s string) string { return rtlID(lang, s) }
	small := func(src, fb string) h.Node {
		return avatar.Avatar(h.Class("size-4"), avatar.Image(h.Src(src)), avatar.Fallback(fb))
	}
	var users []any
	for _, u := range []string{"shadcn", "maxleiter", "evilrabbit"} {
		users = append(users, command.Item(command.Value(u),
			small("https://github.com/"+u+".png", u[:1]), u, " ",
			h.Span(h.Class("text-muted-foreground"), "- ", t["workspace"]),
		))
	}
	// The mention popover and the model menu have no trigger upstream (their
	// open state is never set), so only their tooltips are rendered.
	return h.Div(h.Dir("rtl"),
		h.Form(
			field.Field(
				field.Label(h.For(id("rtl-notion-prompt")), h.Class("sr-only"), t["prompt"]),
				ig.InputGroup(
					ig.Textarea(h.ID(id("rtl-notion-prompt")), h.Placeholder(t["placeholder"])),
					ig.Addon(ig.AlignBlockStart,
						tooltip.Tooltip(
							tooltip.Trigger(lib.AsChild,
								ig.Button(button.VariantOutline, ig.ButtonSizeSm, h.Class("rounded-full transition-transform"),
									rtlIconAt(), " ", t["addContext"])),
							tooltip.Content(t["mentionTooltip"]),
						),
						h.Div(h.Class("-m-1.5 no-scrollbar flex gap-1 overflow-y-auto p-1.5")),
					),
					ig.Addon(ig.AlignBlockEnd, h.Class("gap-1"),
						tooltip.Tooltip(
							tooltip.Trigger(lib.AsChild,
								ig.Button(ig.ButtonSizeIconSm, h.Class("rounded-full"), h.Aria("label", t["attachFile"]), rtlIconPaperclip())),
							tooltip.Content(t["attachFile"]),
						),
						tooltip.Tooltip(
							tooltip.Trigger(lib.AsChild, ig.Button(ig.ButtonSizeSm, h.Class("rounded-full"), t["auto"])),
							tooltip.Content(t["selectModel"]),
						),
						dm.DropdownMenu(dm.Dir("rtl"),
							dm.Trigger(lib.AsChild, ig.Button(ig.ButtonSizeSm, h.Class("rounded-full"), rtlIconWorld(), " ", t["allSources"])),
							dm.Content(dm.SideTop, dm.AlignEnd, h.Class("w-72"), h.Dir("rtl"),
								dm.Group(
									dm.Item(lib.AsChild,
										h.Label(h.For(id("rtl-web-search")), rtlIconWorld(), " ", t["webSearch"], " ",
											switchui.Switch(h.ID(id("rtl-web-search")), h.Class("ms-auto"), switchui.Checked(true), switchui.SizeSm))),
								),
								dm.Separator(),
								dm.Group(
									dm.Item(lib.AsChild,
										h.Label(h.For(id("rtl-apps")), rtlIconApps(), " ", t["appsIntegrations"],
											switchui.Switch(h.ID(id("rtl-apps")), h.Class("ms-auto"), switchui.Checked(true), switchui.SizeSm))),
									dm.Item(rtlIconCircleDashedPlus(), " ", t["allSourcesAccess"]),
									dm.Sub(
										dm.SubTrigger(small("https://github.com/shadcn.png", "CN"), "shadcn"),
										dm.SubContent(h.Class("w-72 rounded-lg p-0"), h.Dir("rtl"),
											command.Command(
												command.Input(h.Placeholder(t["findKnowledge"]), h.AutoFocus()),
												command.List(
													command.Empty(t["noKnowledgeFound"]),
													command.Group(users...),
												),
											),
										),
									),
									dm.Item(rtlIconBook(), " ", t["helpCenter"]),
								),
								dm.Separator(),
								dm.Group(
									dm.Item(tabler.IconPlus(), " ", t["connectApps"]),
									dm.Label(h.Class("text-xs text-muted-foreground"), t["searchSourcesNote"]),
								),
							),
						),
						ig.Button(button.VariantDefault, ig.ButtonSizeIconSm, h.Class("ms-auto rounded-full"), h.Aria("label", t["send"]), rtlIconArrowUp()),
					),
				),
			),
		),
	)
}

func rtlButtonGroupDemo(lang string) h.Node {
	t := rtlTr(lang, map[string]map[string]string{
		"ar": {"goBack": "رجوع", "archive": "أرشفة", "report": "إبلاغ", "snooze": "تأجيل", "moreOptions": "خيارات أخرى", "markAsRead": "تحديد كمقروء", "addToCalendar": "إضافة إلى التقويم", "addToList": "إضافة إلى القائمة", "labelAs": "تصنيف كـ...", "personal": "شخصي", "work": "عمل", "other": "أخرى", "trash": "حذف"},
		"he": {"goBack": "חזור", "archive": "ארכיון", "report": "דווח", "snooze": "נודניק", "moreOptions": "אפשרויות נוספות", "markAsRead": "סמן כנקרא", "addToCalendar": "הוסף ליומן", "addToList": "הוסף לרשימה", "labelAs": "תייג כ...", "personal": "אישי", "work": "עבודה", "other": "אחר", "trash": "מחק"},
	})
	sm := []any{button.VariantOutline, button.SizeSm}
	return h.Div(h.Dir("rtl"),
		buttongroup.ButtonGroup(
			buttongroup.ButtonGroup(h.Class("hidden sm:flex"),
				button.Button(button.VariantOutline, button.SizeIconSm, h.Aria("label", t["goBack"]), lucide.ArrowLeft(h.Class("rtl:rotate-180"))),
			),
			buttongroup.ButtonGroup(
				button.Button(append(sm, t["archive"])...),
				button.Button(append(sm, t["report"])...),
			),
			buttongroup.ButtonGroup(
				button.Button(append(sm, t["snooze"])...),
				dm.DropdownMenu(dm.Dir("rtl"),
					dm.Trigger(lib.AsChild, button.Button(button.VariantOutline, button.SizeIconSm, h.Aria("label", t["moreOptions"]), lucide.MoreHorizontal())),
					// Upstream align is the opposite (see the file comment).
					dm.Content(dm.AlignEnd, h.Dir("rtl"), h.Data("lang", lang), h.Class("w-44"),
						dm.Group(
							dm.Item(lucide.MailCheck(), t["markAsRead"]),
							dm.Item(lucide.Archive(), t["archive"]),
						),
						dm.Separator(),
						dm.Group(
							dm.Item(lucide.Clock(), t["snooze"]),
							dm.Item(lucide.CalendarPlus(), t["addToCalendar"]),
							dm.Item(lucide.ListFilter(), t["addToList"]),
							dm.Sub(
								dm.SubTrigger(lucide.Tag(), t["labelAs"]),
								dm.Portal(dm.SubContent(h.Dir("rtl"), h.Data("lang", lang),
									dm.RadioGroup(dm.Value("personal"),
										dm.RadioItem(dm.Value("personal"), t["personal"]),
										dm.RadioItem(dm.Value("work"), t["work"]),
										dm.RadioItem(dm.Value("other"), t["other"]),
									),
								)),
							),
						),
						dm.Separator(),
						dm.Group(dm.Item(dm.VariantDestructive, lucide.Trash2(), t["trash"])),
					),
				),
			),
		),
	)
}

func rtlFieldCheckbox(lang string) h.Node {
	terms := "أوافق على الشروط والأحكام"
	if lang == "he" {
		terms = "אני מסכים לתנאים וההגבלות"
	}
	id := rtlID(lang, "checkbox-demo-rtl")
	return h.Div(h.Dir("rtl"),
		field.Label(h.For(id),
			field.Field(field.Horizontal,
				checkbox.Checkbox(h.ID(id), checkbox.Checked(true)),
				field.Label(h.For(id), h.Class("line-clamp-1"), terms),
			),
		),
	)
}

func rtlButtonGroupNested(lang string) h.Node {
	prev, next := "السابق", "التالي"
	if lang == "he" {
		prev, next = "הקודם", "הבא"
	}
	return buttongroup.ButtonGroup(h.Dir("rtl"),
		buttongroup.ButtonGroup(
			button.Button(button.VariantOutline, button.SizeSm, rtlNum(lang, 1, 1)),
			button.Button(button.VariantOutline, button.SizeSm, rtlNum(lang, 2, 1)),
			button.Button(button.VariantOutline, button.SizeSm, rtlNum(lang, 3, 1)),
		),
		buttongroup.ButtonGroup(
			button.Button(button.VariantOutline, button.SizeIconSm, h.Aria("label", prev), lucide.ArrowLeft(h.Class("rtl:rotate-180"))),
			button.Button(button.VariantOutline, button.SizeIconSm, h.Aria("label", next), lucide.ArrowRight(h.Class("rtl:rotate-180"))),
		),
	)
}

func rtlButtonGroupPopover(lang string) h.Node {
	t := rtlTr(lang, map[string]map[string]string{
		"ar": {"copilot": "المساعد", "openPopover": "فتح القائمة", "agentTasks": "مهام الوكيل", "placeholder": "صف مهمتك بلغة طبيعية.", "startTask": "ابدأ مهمة جديدة مع المساعد", "description": "صف مهمتك بلغة طبيعية. سيعمل المساعد في الخلفية ويفتح طلب سحب لمراجعتك."},
		"he": {"copilot": "עוזר", "openPopover": "פתח תפריט", "agentTasks": "משימות סוכן", "placeholder": "תאר את המשימה שלך בשפה טבעית.", "startTask": "התחל משימה חדשה עם העוזר", "description": "תאר את המשימה שלך בשפה טבעית. העוזר יעבוד ברקע ויפתח בקשת משיכה לבדיקתך."},
	})
	return buttongroup.ButtonGroup(h.Dir("rtl"),
		button.Button(button.VariantOutline, button.SizeSm, lucide.Bot(), " ", t["copilot"]),
		popover.Popover(
			popover.Trigger(lib.AsChild, button.Button(button.VariantOutline, button.SizeIconSm, h.Aria("label", t["openPopover"]), lucide.ChevronDown())),
			// Upstream align is the opposite (see the file comment).
			popover.Content(popover.AlignEnd, h.Dir("rtl"), h.Data("lang", lang), h.Class("p-0"),
				h.Div(h.Class("px-4 py-3"), h.Div(h.Class("text-sm font-medium"), t["agentTasks"])),
				separator.Separator(),
				h.Div(h.Class("p-4 text-sm *:[p:not(:last-child)]:mb-2"),
					textarea.Textarea(h.Placeholder(t["placeholder"]), h.Class("mb-4 resize-none")),
					h.P(h.Class("font-medium"), t["startTask"]),
					h.P(h.Class("text-muted-foreground"), t["description"]),
				),
			),
		),
	)
}

func rtlFieldHear(lang string) h.Node {
	t := rtlTr(lang, map[string]map[string]string{
		"ar": {"legend": "كيف سمعت عنا؟", "description": "اختر الخيار الذي يصف أفضل طريقة سمعت عنا من خلالها.", "social-media": "التواصل الاجتماعي", "search-engine": "البحث", "referral": "إحالة", "other": "أخرى"},
		"he": {"legend": "איך שמעת עלינו?", "description": "בחר את האפשרות שמתארת בצורה הטובה ביותר כיצד שמעת עלינו.", "social-media": "חברתיות", "search-engine": "חיפוש", "referral": "הפניה", "other": "אחר"},
	})
	var options []any
	for _, v := range []string{"social-media", "search-engine", "referral", "other"} {
		id := rtlID(lang, "rtl-"+v)
		options = append(options, field.Label(h.For(id), h.Class("w-fit!"),
			field.Field(field.Horizontal,
				h.Class("gap-1.5 overflow-hidden px-3! py-1.5! transition-all duration-100 ease-linear"),
				// Upstream's group-has-data-[state=checked]/field-label:px-2! never
				// matches the Base UI checkbox (data-checked), so it is left out.
				// data-checked: (Base UI) is data-[state=checked] on the Radix checkbox.
				checkbox.Checkbox(checkbox.Value(v), h.ID(id), checkbox.Checked(v == "social-media"),
					h.Class("-ms-6 translate-x-1 rounded-full transition-all duration-100 ease-linear data-[state=checked]:ms-0 data-[state=checked]:translate-x-0")),
				field.Title(t[v]),
			),
		))
	}
	return h.Div(h.Dir("rtl"),
		card.Card(h.Class("border-0 py-4 shadow-none"),
			card.Content(h.Class("px-4"),
				h.Form(
					field.Group(
						field.Set(h.Class("gap-4"),
							field.Legend(t["legend"]),
							field.Description(h.Class("line-clamp-1"), t["description"]),
							field.Group(h.Class("flex flex-row flex-wrap gap-2 [--radius:9999rem]"), h.Group(options...)),
						),
					),
				),
			),
		),
	)
}

func rtlSpinnerEmpty(lang string) h.Node {
	t := rtlTr(lang, map[string]map[string]string{
		"ar": {"title": "جارٍ معالجة طلبك", "description": "يرجى الانتظار بينما نعالج طلبك. لا تقم بتحديث الصفحة.", "cancel": "إلغاء"},
		"he": {"title": "מעבד את הבקשה שלך", "description": "אנא המתן בזמן שאנו מעבדים את בקשתך. אל תרענן את הדף.", "cancel": "ביטול"},
	})
	return empty.Empty(h.Class("w-full border md:p-6"), h.Dir("rtl"),
		empty.Header(
			empty.Media(empty.MediaIcon, spinner.Spinner()),
			empty.Title(t["title"]),
			empty.Description(t["description"]),
		),
		empty.Content(button.Button(button.VariantOutline, button.SizeSm, t["cancel"])),
	)
}
