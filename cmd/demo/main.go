// Command demo serves every shadcn/ui Button example rendered with libhtml,
// behind the CSRF and security-header middleware.
//
//	go run ./cmd/demo -dev    # http://localhost:8080
package main

import (
	"crypto/rand"
	"embed"
	"encoding/hex"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"shadcn.go/libhtml/csrf"
	"shadcn.go/libhtml/secure"
)

// The stylesheet is generated with Go only (no Node.js) by the tailwind
// command, which runs Tailwind CSS v4 itself. Run "go generate" after
// changing classes; TestStylesheetUpToDate fails if you forget.
//
//go:generate go run shadcn.go/cmd/tailwind -i styles/app.css -o static/app.css

//go:embed static
var static embed.FS

func main() {
	addr := flag.String("addr", "localhost:8080", "listen address")
	dev := flag.Bool("dev", false, "development mode: allow the CSRF cookie over plain HTTP")
	flag.Parse()

	key := csrfKey()
	protect, err := csrf.New(csrf.Config{Key: key, Insecure: *dev})
	if err != nil {
		log.Fatal(err)
	}

	assets, _ := fs.Sub(static, "static")
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(assets)))
	mux.HandleFunc("GET /{$}", handleIndex)
	mux.HandleFunc("POST /subscribe", handleSubscribe)

	handler := secure.Headers(secure.Config{})(protect.Handler(mux))
	srv := &http.Server{
		Addr:              *addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		MaxHeaderBytes:    1 << 16,
	}
	log.Printf("listening on http://%s", *addr)
	log.Fatal(srv.ListenAndServe())
}

// csrfKey reads a hex key from SHADCN_DEMO_CSRF_KEY, or makes a random one
// (tokens then stop working when the process restarts).
func csrfKey() []byte {
	if s := os.Getenv("SHADCN_DEMO_CSRF_KEY"); s != "" {
		key, err := hex.DecodeString(s)
		if err != nil || len(key) < 32 {
			log.Fatal("SHADCN_DEMO_CSRF_KEY must be at least 32 bytes, hex-encoded")
		}
		return key
	}
	key := make([]byte, 32)
	_, _ = rand.Read(key)
	log.Print("SHADCN_DEMO_CSRF_KEY not set: using a random CSRF key for this run")
	return key
}
