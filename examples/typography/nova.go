package typography

import (
	"shadcn.go/components/style"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

// The typography-rtl docs preview (ui/apps/v4/examples/radix/typography-rtl.tsx).
func init() {
	registry.Register(registry.Example{Name: "nova--typography-rtl", Component: "typography", Render: registry.Static(func() h.Node {
		return style.Scope(style.Nova, NovaRtl())
	})})
}

func NovaRtl() h.Node {
	cell := func(tag func(...any) *h.Element, bold bool, s string) h.Node {
		class := "border px-4 py-2 text-start"
		if bold {
			class += " font-bold"
		}
		return tag(h.Class(class), s)
	}
	row := func(tag func(...any) *h.Element, bold bool, a, b string) h.Node {
		return h.Tr(h.Class("m-0 border-t p-0 even:bg-muted"), cell(tag, bold, a), cell(tag, bold, b))
	}
	return h.Div(h.Dir("rtl"),
		h.H1(h.Class("scroll-m-20 text-4xl font-extrabold tracking-tight text-balance"), "فرض الضرائب على الضحك: سجلات ضريبة النكتة"),
		h.P(h.Class("text-xl leading-7 text-muted-foreground [&:not(:first-child)]:mt-6"),
			"في قديم الزمان، في أرض بعيدة، كان هناك ملك كسول جداً يقضي يومه كله مستلقياً على عرشه. في أحد الأيام، جاءه مستشاروه بمشكلة: المملكة كانت تنفد من المال."),
		h.H2(h.Class("mt-10 scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight transition-colors first:mt-0"), "خطة الملك"),
		h.P(h.Class(pClass), "فكر الملك طويلاً وبجد، وأخيراً توصل إلى", " ",
			h.A(h.Href("#"), h.Class("font-medium text-primary underline underline-offset-4"), "خطة عبقرية"),
			": سيفرض ضريبة على النكات في المملكة."),
		h.Blockquote(h.Class("mt-6 border-s-2 ps-6 italic"),
			`"في النهاية،" قال، "الجميع يستمتع بنكتة جيدة، لذا من العدل أن يدفعوا مقابل هذا الامتياز."`),
		h.H3(h.Class(h3Class), "ضريبة النكتة"),
		h.P(h.Class(pClass), "لم يكن رعايا الملك سعداء. تذمروا واشتكوا، لكن الملك كان حازماً:"),
		h.Ul(h.Class("my-6 ms-6 list-disc [&>li]:mt-2"),
			h.Li("المستوى الأول من التورية: 5 قطع ذهبية"),
			h.Li("المستوى الثاني من النكات: 10 قطع ذهبية"),
			h.Li("المستوى الثالث من النكات القصيرة: 20 قطعة ذهبية"),
		),
		h.P(h.Class(pClass), "نتيجة لذلك، توقف الناس عن رواية النكات، وغرقت المملكة في الكآبة. لكن كان هناك شخص واحد رفض أن تحبطه حماقة الملك: مهرج البلاط المسمى المازح."),
		h.H3(h.Class(h3Class), "ثورة المازح"),
		h.P(h.Class(pClass), "بدأ المازح يتسلل إلى القلعة في منتصف الليل ويترك النكات في كل مكان: تحت وسادة الملك، في حسائه، حتى في المرحاض الملكي. كان الملك غاضباً، لكنه لم يستطع إيقاف المازح."),
		h.P(h.Class(pClass), "وبعد ذلك، في يوم من الأيام، اكتشف سكان المملكة أن النكات التي تركها المازح كانت مضحكة جداً لدرجة أنهم لم يستطيعوا منع أنفسهم من الضحك. وبمجرد أن بدأوا بالضحك، لم يستطيعوا التوقف."),
		h.H3(h.Class(h3Class), "ثورة الشعب"),
		h.P(h.Class(pClass), "شعر سكان المملكة بالبهجة من الضحك، وبدأوا في رواية النكات والتورية مرة أخرى، وسرعان ما أصبحت المملكة بأكملها جزءاً من النكتة."),
		h.Div(h.Class("my-6 w-full overflow-y-auto"),
			h.Table(h.Class("w-full"),
				h.Thead(row(h.Th, true, "خزينة الملك", "سعادة الشعب")),
				h.Tbody(
					row(h.Td, false, "فارغة", "فائضة"),
					row(h.Td, false, "متواضعة", "راضٍ"),
					row(h.Td, false, "ممتلئة", "منتشٍ"),
				),
			),
		),
		h.P(h.Class(pClass), "الملك، عندما رأى مدى سعادة رعاياه، أدرك خطأ طرقه وألغى ضريبة النكتة. أُعلن المازح بطلاً، وعاشت المملكة في سعادة دائمة."),
		h.P(h.Class(pClass), "مغزى القصة هو: لا تستهن أبداً بقوة الضحك الجيد وكن دائماً حذراً من الأفكار السيئة."),
	)
}
