package rtl

import "testing"

func TestClasses(t *testing.T) {
	for in, want := range map[string]string{
		"cn-rtl-flip ml-auto size-4":                        "rtl:rotate-180 ms-auto size-4",
		"data-inset:pl-7 focus:bg-accent":                   "data-inset:ps-7 focus:bg-accent",
		"absolute right-2 border-ring border-r":             "absolute end-2 border-ring border-e",
		"data-[side=left]:-right-1 -translate-x-1/2":        "data-[side=left]:-right-1 -translate-x-1/2 rtl:translate-x-1/2",
		"has-data-[slot=kbd]:pr-1.5 group-focus/x:text-red": "has-data-[slot=kbd]:pe-1.5 group-focus/x:text-red",
	} {
		if got := Classes(in); got != want {
			t.Errorf("Classes(%q) = %q, want %q", in, got, want)
		}
	}
}
