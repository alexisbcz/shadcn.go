// Command parity serves every example twice, in an identical page shell:
//
//	/react/<name>   the upstream React example (parity/reference build)
//	/go/<name>      the Go port (examples/<component>)
//
// Both use the same stylesheet (parity/styles.css, compiled at startup by
// the Go Tailwind toolchain from the upstream sources and the Go sources),
// so any visual difference comes from the markup or behavior. The
// Playwright runner in parity/run.mjs drives it.
//
//	go run ./cmd/parity -addr 127.0.0.1:4600
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"shadcn.go/components/ui/runtime"
	_ "shadcn.go/examples/all"
	"shadcn.go/examples/registry"
	h "shadcn.go/libhtml"
	"shadcn.go/libhtml/secure"
	"shadcn.go/tailwind"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:4600", "listen address (port 0 picks a free one)")
	root := flag.String("root", ".", "repository root")
	flag.Parse()

	rootAbs, err := filepath.Abs(*root)
	if err != nil {
		log.Fatal(err)
	}
	start := time.Now()
	css, err := tailwind.Build(tailwind.Options{Input: filepath.Join(rootAbs, "parity/styles.css"), Dir: rootAbs})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("css: %d bytes from %d files in %s", len(css.CSS), len(css.Files), time.Since(start).Round(time.Millisecond))

	reference := filepath.Join(rootAbs, "parity/reference/dist")
	upstream := upstreamExamples(rootAbs)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /assets/app.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		_, _ = w.Write([]byte(css.CSS))
	})
	mux.Handle("GET /reference/", http.StripPrefix("/reference/", http.FileServer(http.Dir(reference))))
	mux.Handle("GET "+runtime.Prefix, runtime.Handler())

	// Go pages get the same strict CSP as production apps, so the port is
	// verified to work under it. The React reference gets none: upstream
	// libraries inject <style> elements (sonner, vaul) that it would block.
	strict := secure.Headers(secure.Config{})
	mux.Handle("GET /go/{name}", strict(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ex, ok := registry.Get(r.PathValue("name"))
		if !ok {
			http.Error(w, "example not ported", http.StatusNotFound)
			return
		}
		var extra h.Node
		if ex.Extra != nil {
			extra = ex.Extra(r)
		}
		page := shell(r, ex.Name+" · go", h.Group(runtime.Scripts(r)), ex.Render(r), extra)
		if err := h.WriteHTTP(w, http.StatusOK, page); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})))
	mux.HandleFunc("GET /react/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		if !slices.Contains(upstream, name) {
			http.Error(w, "no such upstream example", http.StatusNotFound)
			return
		}
		head := h.Script(h.Type("module"), h.Src("/reference/app.js"))
		page := shell(r, name+" · react", head, nil, nil)
		if err := h.WriteHTTP(w, http.StatusOK, page); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("GET /api/examples", func(w http.ResponseWriter, r *http.Request) {
		type item struct {
			Name      string `json:"name"`
			Component string `json:"component,omitempty"`
			Ported    bool   `json:"ported"`
		}
		var list []item
		for _, n := range upstream {
			ex, ok := registry.Get(n)
			list = append(list, item{Name: n, Component: ex.Component, Ported: ok})
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(list)
	})
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		_ = h.WriteHTTP(w, http.StatusOK, index(r, upstream))
	})

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	// The runner reads this line to find the port.
	fmt.Printf("parity listening on http://%s\n", ln.Addr())
	log.Fatal(http.Serve(ln, mux))
}

// shell is the page both sides render into. #preview mirrors the preview
// box of the upstream docs; #parity-extra holds page-level companions such
// as sonner's <Toaster />.
func shell(r *http.Request, title string, head, preview, extra h.Node) h.Node {
	theme := "light"
	if r.URL.Query().Get("theme") == "dark" {
		theme = "dark"
	}
	return h.Doctype(h.Html(h.Lang("en"), h.Class(theme), h.Style("color-scheme: "+theme),
		h.Head(
			h.Meta(h.Charset("utf-8")),
			h.Meta(h.Name("viewport"), h.Content("width=device-width, initial-scale=1")),
			h.Title(title),
			h.Link(h.Rel("stylesheet"), h.Href("/assets/app.css")),
			head,
		),
		h.Body(h.Class("bg-background text-foreground antialiased"),
			h.Div(h.ID("preview"), h.Data("parity-root", ""),
				h.Class("flex min-h-[350px] w-full items-center justify-center p-10"),
				preview,
			),
			h.Div(h.ID("parity-extra"), extra),
		),
	))
}

func index(r *http.Request, upstream []string) h.Node {
	ported := 0
	for _, n := range upstream {
		if _, ok := registry.Get(n); ok {
			ported++
		}
	}
	return shell(r, "parity", nil, h.Div(h.Class("w-full max-w-3xl"),
		h.H1(h.Class("mb-4 text-2xl font-semibold"), fmt.Sprintf("Examples: %d/%d ported", ported, len(upstream))),
		h.Ul(h.Class("grid grid-cols-2 gap-1 text-sm"), h.Map(upstream, func(n string) h.Node {
			_, ok := registry.Get(n)
			return h.Li(
				h.A(h.Class("underline"), h.Href("/react/"+n), n),
				" · ",
				h.IfElse(ok, h.Node(h.A(h.Class("underline"), h.Href("/go/"+n), "go")), h.Node(h.Span(h.Class("text-muted-foreground"), "todo"))),
			)
		})),
	), nil)
}

// upstreamExamples lists the upstream new-york-v4 example names.
func upstreamExamples(root string) []string {
	dir := filepath.Join(root, "ui/apps/v4/registry/new-york-v4/examples")
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("upstream examples: %v (clone shadcn-ui/ui into ./ui)", err)
	}
	var names []string
	for _, e := range entries {
		if n, ok := strings.CutSuffix(e.Name(), ".tsx"); ok {
			names = append(names, n)
		}
	}
	return names
}
