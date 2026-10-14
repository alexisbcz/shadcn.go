// Command docs serves the Go port of the shadcn/ui documentation site
// (ui/apps/v4, ui.shadcn.com) from the repository root:
//
//	go run ./cmd/docs -addr :4100
//
// then open http://localhost:4100/docs/components/radix/button. It serves
// the reference's URLs: /, /docs/[[...slug]] for every page of
// ui/apps/v4/content/docs, /docs/changelog, and the static files of
// ui/apps/v4/public. The stylesheet is built at startup by the Go Tailwind
// package from docs/site/styles.css (the app's globals.css plus the Go
// sources) and cached in the user cache directory, so the first start takes
// a few seconds longer. Pages are parsed on first request (-preload parses
// them all at startup).
//
// The reference to compare against runs at http://localhost:4000 (see
// parity/docs/reference.sh); parity/docs/run.mjs diffs the two.
//
// The site itself is package shadcn.go/docs/site; its README-level comment
// lists the extension points (header, docs layout, search, MDX components).
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	_ "shadcn.go/examples/all"

	"shadcn.go/docs/site"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:4100", "listen address (port 0 picks a free one)")
	root := flag.String("root", ".", "repository root")
	preload := flag.Bool("preload", false, "parse every docs page at startup")
	flag.Parse()

	start := time.Now()
	s, err := site.New(site.Options{Root: *root, Preload: *preload})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("docs: %d pages loaded in %s", len(s.Docs), time.Since(start).Round(time.Millisecond))

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	// The harness reads this line to find the port.
	fmt.Printf("docs listening on http://%s\n", ln.Addr())
	log.Fatal(http.Serve(ln, s.Handler()))
}
