package form

import (
	"reflect"
	"strings"
	"testing"

	h "shadcn.go/libhtml"
)

func TestNextDemoAction(t *testing.T) {
	got := NextDemoAction(NextDemoValues{Title: "Bug", Description: "short"})
	want := map[string][]string{
		"title":       {"Bug title must be at least 5 characters."},
		"description": {"Description must be at least 20 characters."},
	}
	if got.Success || !reflect.DeepEqual(got.Errors, want) {
		t.Fatalf("got %+v", got)
	}
	ok := NextDemoAction(NextDemoValues{Title: "Login broken", Description: "Steps: open the app and tap login."})
	if !ok.Success || ok.Errors != nil || ok.Values != (NextDemoValues{}) {
		t.Fatalf("got %+v", ok)
	}
}

func TestNextDemoRendersState(t *testing.T) {
	s := h.MustRenderString(NextDemo(NextDemoAction(NextDemoValues{Title: "<b>", Description: "abc"}), true))
	for _, want := range []string{"&lt;b&gt;", `aria-invalid="true"`, `data-invalid="true"`, "Bug title must be at least 5 characters.", `disabled`, `>3<!---->/100 characters`} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %q", want)
		}
	}
}

func TestNextComplexAction(t *testing.T) {
	got := NextComplexAction(NextComplexValues{Plan: "", BillingPeriod: "", Addons: []string{"x", "y", "z", "w"}})
	want := map[string][]string{
		"plan":          {"Please select a subscription plan", "Invalid plan selection. Please choose Basic or Pro"},
		"billingPeriod": {"Please select a billing period"},
		"addons":        {"You can select up to 3 add-ons", "You selected an invalid add-on"},
	}
	if !reflect.DeepEqual(got.Errors, want) {
		t.Fatalf("got %+v", got.Errors)
	}
	if !NextComplexAction(NextComplexValues{Plan: "pro", BillingPeriod: "monthly", Addons: []string{"backup"}}).Success {
		t.Fatal("valid values rejected")
	}
}

func TestJSLen(t *testing.T) {
	if jsLen("a😀") != 3 {
		t.Fatal("UTF-16 length")
	}
}
