package mdx

// Questionnaire and its parts (QuestionnaireItem, QuestionnaireChoice, ...). Package content does not route these names to the registry yet; they appear in code samples. Registered here so the implementation has a home.
//
// Placeholder: content.Placeholder (an <mdx-component> element with the
// props). Replace with the real rendering.
func init() {
	Register("Questionnaire", placeholder("Questionnaire"))
}
