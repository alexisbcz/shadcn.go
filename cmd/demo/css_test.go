package main

import (
	"os"
	"testing"

	"shadcn.go/tailwind"
)

// TestStylesheetUpToDate fails when static/app.css no longer matches the
// classes used in the code. Fix it with: go generate ./cmd/demo
func TestStylesheetUpToDate(t *testing.T) {
	res, err := tailwind.Build(tailwind.Options{Input: "styles/app.css", Output: "static/app.css"})
	if err != nil {
		t.Fatal(err)
	}
	have, err := os.ReadFile("static/app.css")
	if err != nil {
		t.Fatal(err)
	}
	if string(have) != res.CSS {
		t.Fatal("static/app.css is stale: run go generate ./cmd/demo")
	}
}
