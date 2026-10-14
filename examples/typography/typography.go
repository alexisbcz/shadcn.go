// Package typography ports the upstream typography-* examples: plain
// markup styled with utility classes (shadcn ships no typography component).
package typography

import (
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
)

func init() {
	ex := func(name string, fn func() h.Node) registry.Example {
		return registry.Example{Name: name, Component: "typography", Render: registry.Static(fn)}
	}
	registry.Register(
		ex("typography-blockquote", Blockquote),
		ex("typography-demo", Demo),
		ex("typography-h1", H1),
		ex("typography-h2", H2),
		ex("typography-h3", H3),
		ex("typography-h4", H4),
		ex("typography-inline-code", InlineCode),
		ex("typography-large", Large),
		ex("typography-lead", Lead),
		ex("typography-list", List),
		ex("typography-muted", Muted),
		ex("typography-p", P),
		ex("typography-small", Small),
		ex("typography-table", Table),
	)
}

const (
	pClass    = "leading-7 [&:not(:first-child)]:mt-6"
	h3Class   = "mt-8 scroll-m-20 text-2xl font-semibold tracking-tight"
	quoteText = `"After all," he said, "everyone enjoys a good joke, so it's only fair that they should pay for the privilege."`
)

func Blockquote() h.Node {
	return h.Blockquote(h.Class("mt-6 border-l-2 pl-6 italic"), quoteText)
}

func H1() h.Node {
	return h.H1(h.Class("scroll-m-20 text-center text-4xl font-extrabold tracking-tight text-balance"),
		"Taxing Laughter: The Joke Tax Chronicles")
}

func H2() h.Node {
	return h.H2(h.Class("scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0"),
		"The People of the Kingdom")
}

func H3() h.Node {
	return h.H3(h.Class("scroll-m-20 text-2xl font-semibold tracking-tight"), "The Joke Tax")
}

func H4() h.Node {
	return h.H4(h.Class("scroll-m-20 text-xl font-semibold tracking-tight"), "People stopped telling jokes")
}

func InlineCode() h.Node {
	return h.Code(h.Class("relative rounded bg-muted px-[0.3rem] py-[0.2rem] font-mono text-sm font-semibold"),
		"@radix-ui/react-alert-dialog")
}

func Large() h.Node { return h.Div(h.Class("text-lg font-semibold"), "Are you absolutely sure?") }

func Lead() h.Node {
	return h.P(h.Class("text-xl text-muted-foreground"),
		"A modal dialog that interrupts the user with important content and expects a response.")
}

func List() h.Node {
	return h.Ul(h.Class("my-6 ml-6 list-disc [&>li]:mt-2"),
		h.Li("1st level of puns: 5 gold coins"),
		h.Li("2nd level of jokes: 10 gold coins"),
		h.Li("3rd level of one-liners : 20 gold coins"),
	)
}

func Muted() h.Node {
	return h.P(h.Class("text-sm text-muted-foreground"), "Enter your email address.")
}

// P is the typography-p example: on its own in the preview it is the first
// child, so [&:not(:first-child)]:mt-6 does not apply.
func P() h.Node {
	return h.P(h.Class(pClass),
		"The king, seeing how much happier his subjects were, realized the error of his ways and repealed the joke tax.")
}

func Small() h.Node {
	return h.Small(h.Class("text-sm leading-none font-medium"), "Email address")
}

func Table() h.Node {
	const (
		trClass = "m-0 border-t p-0 even:bg-muted"
		thClass = "border px-4 py-2 text-left font-bold [&[align=center]]:text-center [&[align=right]]:text-right"
		tdClass = "border px-4 py-2 text-left [&[align=center]]:text-center [&[align=right]]:text-right"
	)
	row := func(a, b string) h.Node {
		return h.Tr(h.Class(trClass), h.Td(h.Class(tdClass), a), h.Td(h.Class(tdClass), b))
	}
	return h.Div(h.Class("my-6 w-full overflow-y-auto"),
		h.Table(h.Class("w-full"),
			h.Thead(
				h.Tr(h.Class(trClass),
					h.Th(h.Class(thClass), "King's Treasury"),
					h.Th(h.Class(thClass), "People's happiness"),
				),
			),
			h.Tbody(
				row("Empty", "Overflowing"),
				row("Modest", "Satisfied"),
				row("Full", "Ecstatic"),
			),
		),
	)
}

func Demo() h.Node {
	return h.Div(
		h.H1(h.Class("scroll-m-20 text-4xl font-extrabold tracking-tight text-balance"),
			"Taxing Laughter: The Joke Tax Chronicles"),
		h.P(h.Class("text-xl leading-7 text-muted-foreground [&:not(:first-child)]:mt-6"),
			"Once upon a time, in a far-off land, there was a very lazy king who spent all day lounging on his throne. One day, his advisors came to him with a problem: the kingdom was running out of money."),
		h.H2(h.Class("mt-10 scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight transition-colors first:mt-0"),
			"The King's Plan"),
		h.P(h.Class(pClass),
			"The king thought long and hard, and finally came up with", " ",
			h.A(h.Href("#"), h.Class("font-medium text-primary underline underline-offset-4"), "a brilliant plan"),
			": he would tax the jokes in the kingdom.",
		),
		h.Blockquote(h.Class("mt-6 border-l-2 pl-6 italic"), quoteText),
		h.H3(h.Class(h3Class), "The Joke Tax"),
		h.P(h.Class(pClass),
			"The king's subjects were not amused. They grumbled and complained, but the king was firm:"),
		List(),
		h.P(h.Class(pClass),
			"As a result, people stopped telling jokes, and the kingdom fell into a gloom. But there was one person who refused to let the king's foolishness get him down: a court jester named Jokester."),
		h.H3(h.Class(h3Class), "Jokester's Revolt"),
		h.P(h.Class(pClass),
			"Jokester began sneaking into the castle in the middle of the night and leaving jokes all over the place: under the king's pillow, in his soup, even in the royal toilet. The king was furious, but he couldn't seem to stop Jokester."),
		h.P(h.Class(pClass),
			"And then, one day, the people of the kingdom discovered that the jokes left by Jokester were so funny that they couldn't help but laugh. And once they started laughing, they couldn't stop."),
		h.H3(h.Class(h3Class), "The People's Rebellion"),
		h.P(h.Class(pClass),
			"The people of the kingdom, feeling uplifted by the laughter, started to tell jokes and puns again, and soon the entire kingdom was in on the joke."),
		Table(),
		h.P(h.Class(pClass),
			"The king, seeing how much happier his subjects were, realized the error of his ways and repealed the joke tax. Jokester was declared a hero, and the kingdom lived happily ever after."),
		h.P(h.Class(pClass),
			"The moral of the story is: never underestimate the power of a good laugh and always be careful of bad ideas."),
	)
}
