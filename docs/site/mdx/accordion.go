package mdx

// Accordion, AccordionItem, AccordionTrigger, AccordionContent: new-york-v4 accordion (Accordion gets className "not-typeset" from package content).
//
// Placeholder: content.Placeholder (an <mdx-component> element with the
// props). Replace with the real rendering.
func init() {
	Register("Accordion", placeholder("Accordion"))
	Register("AccordionItem", placeholder("AccordionItem"))
	Register("AccordionTrigger", placeholder("AccordionTrigger"))
	Register("AccordionContent", placeholder("AccordionContent"))
}
