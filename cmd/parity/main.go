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
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
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
	list := flag.Bool("list", false, "print the references as JSON (name, kind, file, component, ported) and exit")
	flag.Parse()

	rootAbs, err := filepath.Abs(*root)
	if err != nil {
		log.Fatal(err)
	}
	if *list {
		refs := loadReferences(filepath.Join(rootAbs, "parity/reference/dist"))
		out := make([]map[string]any, 0, len(refs))
		for _, r := range refs {
			ex, ok := registry.Get(r.Name)
			out = append(out, map[string]any{"name": r.Name, "kind": r.Kind, "file": r.File, "component": ex.Component, "ported": ok})
		}
		_ = json.NewEncoder(os.Stdout).Encode(out)
		return
	}
	start := time.Now()
	css, err := cachedCSS(rootAbs)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("css: %d bytes in %s", len(css), time.Since(start).Round(time.Millisecond))

	reference := filepath.Join(rootAbs, "parity/reference/dist")
	refs := loadReferences(reference)
	upstream := make([]string, 0, len(refs))
	for _, r := range refs {
		upstream = append(upstream, r.Name)
	}
	kindOf := func(name string) string {
		for _, r := range refs {
			if r.Name == name {
				return r.Kind
			}
		}
		return ""
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /assets/app.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		_, _ = w.Write([]byte(css))
	})
	mux.Handle("GET /reference/", http.StripPrefix("/reference/", http.FileServer(http.Dir(reference))))
	mux.Handle("GET "+runtime.Prefix, runtime.Handler())

	// Go pages get the same strict CSP as production apps, so the port is
	// verified to work under it. The React reference gets none: upstream
	// libraries inject <style> elements (sonner, vaul) that it would block.
	// Upstream examples load remote images (unsplash, avatars), so img-src
	// also allows https:; everything else is the production default.
	strict := secure.Headers(secure.Config{CSP: func(nonce string) string {
		return strings.Replace(secure.DefaultCSP(nonce), "img-src 'self' data:", "img-src 'self' data: https:", 1)
	}})
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
		page := shell(r, kindOf(ex.Name), ex.Name+" · go", h.Group(runtime.Scripts(r)), ex.Render(r), extra)
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
		page := shell(r, kindOf(name), name+" · react", head, nil, nil)
		if err := h.WriteHTTP(w, http.StatusOK, page); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("GET /api/examples", func(w http.ResponseWriter, r *http.Request) {
		type item struct {
			Name      string `json:"name"`
			Kind      string `json:"kind"`
			File      string `json:"file"`
			Component string `json:"component,omitempty"`
			Ported    bool   `json:"ported"`
		}
		var list []item
		for _, ref := range refs {
			ex, ok := registry.Get(ref.Name)
			list = append(list, item{Name: ref.Name, Kind: ref.Kind, File: ref.File, Component: ex.Component, Ported: ok})
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(list)
	})
	// Upstream examples reference the docs site's public assets
	// (/placeholder.svg, /avatars/*.jpg, ...): serve them to both sides, as
	// the real site does. More specific routes above take precedence.
	mux.Handle("GET /", http.FileServer(http.Dir(filepath.Join(rootAbs, "ui/apps/v4/public"))))
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
// box of the upstream docs (blocks are full pages and get no box);
// #parity-extra holds page-level companions such as sonner's <Toaster />.
func shell(r *http.Request, kind, title string, head, preview, extra h.Node) h.Node {
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
				h.If(kind != "block", h.Class("flex min-h-[350px] w-full items-center justify-center p-10")),
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
	return shell(r, "", "parity", nil, h.Div(h.Class("w-full max-w-3xl"),
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

// reference is one upstream example, chart, block, internal demo or
// fixture, as listed by the reference build.
type reference struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
	File string `json:"file"`
}

func loadReferences(dist string) []reference {
	b, err := os.ReadFile(filepath.Join(dist, "references.json"))
	if err != nil {
		log.Fatalf("%v: build the React reference first (node parity/reference/build.mjs)", err)
	}
	var refs []reference
	if err := json.Unmarshal(b, &refs); err != nil {
		log.Fatal(err)
	}
	return refs
}

// cachedCSS builds parity/styles.css, reusing the previous result while no
// source file changed. Building scans thousands of files and runs the
// Tailwind compiler; every harness run starts a server, so the cache saves
// seconds of CPU per run when many agents run the harness at once. The key
// covers path, size and modification time of every file under the
// stylesheet's @source directories and the stylesheet itself.
func cachedCSS(root string) (string, error) {
	sum := sha256.New()
	for _, dir := range []string{"parity/styles.css", "ui/apps/v4/registry/new-york-v4", "components", "examples", "cmd/parity", "parity/fixtures"} {
		_ = filepath.WalkDir(filepath.Join(root, dir), func(p string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			if info, err := d.Info(); err == nil {
				fmt.Fprintf(sum, "%s\x00%d\x00%d\n", p, info.Size(), info.ModTime().UnixNano())
			}
			return nil
		})
	}
	fmt.Fprintf(sum, "tailwind %s\n", tailwind.Version)
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}
	cacheFile := filepath.Join(cacheDir, "shadcn-parity", hex.EncodeToString(sum.Sum(nil))+".css")
	if b, err := os.ReadFile(cacheFile); err == nil {
		return string(b), nil
	}
	res, err := tailwind.Build(tailwind.Options{Input: filepath.Join(root, "parity/styles.css"), Dir: root})
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(cacheFile), 0o755); err == nil {
		tmp := cacheFile + fmt.Sprintf(".%d.tmp", os.Getpid())
		if os.WriteFile(tmp, []byte(res.CSS), 0o644) == nil {
			_ = os.Rename(tmp, cacheFile)
		}
	}
	return res.CSS, nil
}
