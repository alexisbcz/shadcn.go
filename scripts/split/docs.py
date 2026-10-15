import os, sys, shutil
S = sys.argv[1]
OUT = os.path.join(S, "split")
CVA_LICENSE = "/Users/alex/Work/shadcn.go/ui/node_modules/.pnpm/class-variance-authority@0.7.1/node_modules/class-variance-authority/LICENSE"

README = {
"libhtml": ("libhtml", "Build HTML in Go with a JSX-like API. It escapes and sanitizes output by default.",
'''package main

import (
	"fmt"

	h "github.com/alexisbcz/libhtml"
)

func main() {
	page := h.Div(h.Class("card"),
		h.H1("Hello"),
		h.P("Any value can be a child: ", 42, " items."),
	)
	html, _ := h.RenderString(page)
	fmt.Println(html)
}'''),
"cn": ("cn", "Merge Tailwind CSS classes like the cn() helper of shadcn/ui. Later classes win over conflicting earlier ones.",
'''package main

import (
	"fmt"

	"github.com/alexisbcz/cn"
)

func main() {
	fmt.Println(cn.Cn("px-2 py-1 bg-red-500", "p-3 bg-blue-500")) // p-3 bg-blue-500
}'''),
"cva": ("cva", "Build class names from variants, like class-variance-authority.",
'''package main

import (
	"fmt"

	"github.com/alexisbcz/cva"
)

var button = cva.New("inline-flex items-center", cva.Config{
	Variants: cva.Variants{
		"size": {"default": "h-9 px-4", "sm": "h-8 px-3"},
	},
	Defaults: cva.Selection{"size": "default"},
})

func main() {
	fmt.Println(button.Classes(cva.Selection{"size": "sm"})) // inline-flex items-center h-8 px-3
}'''),
"twmerge": ("twmerge", "Merge Tailwind CSS class lists, like tailwind-merge. Later classes win over conflicting earlier ones.",
'''package main

import (
	"fmt"

	"github.com/alexisbcz/twmerge"
)

func main() {
	fmt.Println(twmerge.Merge("rounded-md px-2", "rounded-full")) // px-2 rounded-full
}'''),
"tailwind": ("tailwind", "Tailwind CSS v4 in Go. Build stylesheets without Node.js.",
'''package main

import (
	"log"
	"os"

	"github.com/alexisbcz/tailwind"
)

func main() {
	res, err := tailwind.Build(tailwind.Options{Input: "styles/app.css", Minify: true})
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile("static/app.css", []byte(res.CSS), 0o644)
}'''),
"icons": ("icons", "Lucide, Tabler and Hugeicons icons as Go functions that render SVG.",
'''package main

import (
	"fmt"

	"github.com/alexisbcz/icons/hugeicons"
	"github.com/alexisbcz/icons/lucide"
	"github.com/alexisbcz/icons/tabler"
	h "github.com/alexisbcz/libhtml"
)

func main() {
	for _, icon := range []h.Node{
		lucide.Check(h.Class("size-4")),
		tabler.IconBell(),
		hugeicons.Add01Icon(),
	} {
		html, _ := h.RenderString(icon)
		fmt.Println(html)
	}
}'''),
"components": ("components", "The shadcn/ui components for Go. Each component renders HTML on the server, and plain JavaScript drives its client behavior.",
'''package main

import (
	"net/http"

	"github.com/alexisbcz/components/ui/button"
	"github.com/alexisbcz/components/ui/runtime"
	h "github.com/alexisbcz/libhtml"
)

func main() {
	http.Handle(runtime.Prefix, runtime.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page := h.Html(
			h.Head(runtime.Scripts(r)),
			h.Body(button.Button(button.VariantOutline, "Click me")),
		)
		h.WriteHTTP(w, http.StatusOK, page)
	})
	http.ListenAndServe(":8080", nil)
}'''),
}
EXTRA = {
"tailwind": "\nOr use the command line tool:\n\n```sh\ngo run github.com/alexisbcz/tailwind/cmd/tailwind -i styles/app.css -o static/app.css --minify\n```\n",
}

NOTICES = {
"cn": [("cn 0.2.2", "MIT", "Copyright (c) 2026 shadcn")],
"cva": [("class-variance-authority 0.7.1", "Apache-2.0", "Joe Bell. The full license text follows")],
"twmerge": [("tailwind-merge", "MIT", "Copyright (c) 2021 Dany Castillo")],
"tailwind": [("Tailwind CSS 4.3.0, embedded in dist/", "MIT", "Copyright (c) Tailwind Labs, Inc. See dist/LICENSE")],
"icons": [
    ("Lucide (lucide-react 0.474.0)", "ISC", "Copyright (c) Lucide Contributors 2022. Portions copyright (c) Cole Bemis 2013-2022, as part of Feather (MIT)"),
    ("Tabler Icons (@tabler/icons-react 3.34.1)", "MIT", "Copyright (c) 2020-2024 Paweł Kuna"),
    ("Hugeicons (@hugeicons/core-free-icons 1.2.1)", "MIT", "Copyright (c) Hugeicons"),
],
"components": [
    ("shadcn/ui", "MIT", "Copyright (c) 2023 shadcn"),
    ("Radix UI primitives", "MIT", "Copyright (c) 2022 WorkOS"),
    ("Base UI (@base-ui/react)", "MIT", "Copyright (c) 2019 Material-UI SAS"),
    ("Floating UI (@floating-ui/dom)", "MIT", "Copyright (c) 2021-present Floating UI contributors"),
    ("Embla Carousel 8.5.2", "MIT", "Copyright (c) 2019-present David Jerleke"),
    ("Recharts 3.8.0", "MIT", "Copyright (c) 2015-present recharts"),
    ("d3-scale, d3-shape, d3-path, d3-array, d3-interpolate, d3-format, d3-color, internmap", "ISC", "Copyright 2010-2023 Mike Bostock"),
    ("Sonner 2.0.7", "MIT", "Copyright (c) 2023 Emil Kowalski"),
    ("Vaul", "MIT", "Copyright (c) 2023 Emil Kowalski"),
    ("cmdk", "MIT", "Copyright (c) 2022 Paco Coursey"),
    ("input-otp", "MIT", "Copyright (c) Guilherme Rodz"),
    ("react-resizable-panels", "MIT", "Copyright (c) 2018 Brian Vaughn"),
    ("react-hook-form", "MIT", "Copyright (c) 2019-present Beier (Bill) Luo"),
    ("react-day-picker", "MIT", "Copyright (c) 2014-2025 Giampaolo Bellavite and contributors"),
    ("aria-hidden, react-remove-scroll", "MIT", "Copyright (c) 2017 Anton Korzunov"),
    ("tw-animate-css", "MIT", "Copyright (c) 2025 Wombosvideo"),
],
}

for repo, (title, desc, code) in README.items():
    d = os.path.join(OUT, repo)
    install = f"go get github.com/alexisbcz/{repo}"
    if repo in ("icons", "components"):
        # No package at the module root: go get alone records no checksums
        # for the imported packages' dependencies.
        install += "\ngo mod tidy"
    text = f"# {title}\n\n{desc}\n\n## Install\n\n```sh\n{install}\n```\n\n## Use\n\n```go\n{code}\n```\n"
    text += EXTRA.get(repo, "")
    text += "\n## License\n\nMIT. See [LICENSE](LICENSE)"
    text += " and [NOTICE](NOTICE).\n" if repo in NOTICES else ".\n"
    with open(os.path.join(d, "README.md"), "w") as fh:
        fh.write(text)
    if repo in NOTICES:
        n = "This project includes code ported from or bundled with the following projects.\n\n"
        for name, lic, holder in NOTICES[repo]:
            n += f"- {name}: {lic} License. {holder}.\n"
        if repo == "cva":
            n += "\n" + open(CVA_LICENSE).read()
        with open(os.path.join(d, "NOTICE"), "w") as fh:
            fh.write(n)
print("ok")
