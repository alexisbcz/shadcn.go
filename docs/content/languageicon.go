package content

import h "shadcn.go/libhtml"

// LanguageIcon is getIconForLanguageExtension (ui/apps/v4/components/icons.tsx),
// for components that render code titles themselves (ComponentSource).
func LanguageIcon(language string) h.Node { return languageIcon(language) }
