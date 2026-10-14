package calendar

import (
	"strings"
	"testing"
	"time"

	h "shadcn.go/libhtml"
)

func TestJalaliConversions(t *testing.T) {
	// 2025-06-12 is 22 Khordad 1404; 2025-03-21 is 1 Farvardin 1404.
	for _, c := range []struct {
		g       time.Time
		y, m, d int
	}{
		{time.Date(2025, 6, 12, 0, 0, 0, 0, time.UTC), 1404, 2, 22},
		{time.Date(2025, 3, 21, 0, 0, 0, 0, time.UTC), 1404, 0, 1},
		{time.Date(2025, 3, 20, 0, 0, 0, 0, time.UTC), 1403, 11, 30}, // 1403 is leap
	} {
		y, m, d := jalali{}.ymd(JDN(c.g))
		if y != c.y || m != c.m || d != c.d {
			t.Errorf("%s: got %d-%d-%d, want %d-%d-%d", c.g.Format(time.DateOnly), y, m, d, c.y, c.m, c.d)
		}
		if back := (jalali{}).jdn(y, m, d); back != JDN(c.g) {
			t.Errorf("round trip of %s: %d != %d", c.g.Format(time.DateOnly), back, JDN(c.g))
		}
	}
	if (jalali{}).dim(1403, 11) != 30 || (jalali{}).dim(1404, 11) != 29 {
		t.Error("Esfand length")
	}
}

func TestGregorianHelpers(t *testing.T) {
	l := newDateLib(false, "", nil)
	j := JDN(time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC))
	if got := isoOf(l.addMonths(j, 1)); got != "2024-02-29" {
		t.Errorf("addMonths clamps: %s", got)
	}
	if got := isoOf(l.addMonths(j, -13)); got != "2022-12-31" {
		t.Errorf("addMonths back: %s", got)
	}
	if got := l.fullDate(JDN(time.Date(2022, 11, 21, 0, 0, 0, 0, time.UTC))); got != "Monday, November 21st, 2022" {
		t.Errorf("PPPP: %s", got)
	}
	for n, want := range map[int]string{1: "1st", 2: "2nd", 3: "3rd", 4: "4th", 11: "11th", 12: "12th", 13: "13th", 21: "21st", 22: "22nd", 23: "23rd", 31: "31st"} {
		if got := enOrdinal(n); got != want {
			t.Errorf("ordinal %d: %s", n, got)
		}
	}
	if l.weekday(JDN(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))) != 6 {
		t.Error("2000-01-01 is a Saturday")
	}
}

func TestPersianFormats(t *testing.T) {
	l := newDateLib(true, "arabext", nil)
	j := JDN(time.Date(2025, 5, 17, 0, 0, 0, 0, time.UTC))
	if got := l.fullDate(j); got != "شنبه ۲۷-ام اردیبهشت ۱۴۰۴" {
		t.Errorf("PPPP: %s", got)
	}
	if got := l.isoDay(j); got != "۱۴۰۴-۰۲-۲۷" {
		t.Errorf("yyyy-MM-dd: %s", got)
	}
	if l.startOfWeek(j) != j {
		t.Error("Persian weeks start on Saturday")
	}
}

func TestRenderSingle(t *testing.T) {
	today := time.Date(2026, 10, 14, 0, 0, 0, 0, time.UTC)
	got, err := h.RenderString(Calendar(ModeSingle, Today(today), Selected(today), h.Class("rounded-md border")))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`data-slot="calendar"`,
		`data-mode="single"`,
		`rounded-md border`,
		`aria-label="October 2026"`,
		`<td class="`,
		`role="gridcell" data-day="2026-10-14" aria-selected="true" data-selected="true" data-today="true">`,
		`data-day="10/14/2026"`,
		`data-selected-single="true"`,
		`aria-label="Today, Wednesday, October 14th, 2026, selected"`,
		`data-month="2026-09" data-outside="true"`,
		`>October 2026</span>`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q", want)
		}
	}
	if strings.Count(got, `tabindex="0"`) != 1 {
		t.Error("exactly one day is the focus target")
	}
}

func TestRenderRangeAndDropdown(t *testing.T) {
	from := time.Date(2022, 1, 20, 0, 0, 0, 0, time.UTC)
	got, err := h.RenderString(Calendar(ModeRange, Today(from), DefaultMonth(from), SelectedRange(from, from.AddDate(0, 0, 20)), NumberOfMonths(2), CaptionDropdown))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`data-multiple-months="true"`,
		`aria-multiselectable="true"`,
		`data-range-start="true"`,
		`data-range-end="true"`,
		`data-range-middle="true"`,
		`data-range-start="false"`,
		`>January 2022</span>`,
		`>February 2022</span>`,
		`<option value="1922">1922</option>`,
		`aria-label="Choose the Month"`,
		`data-disabled="false"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q", want)
		}
	}
}

func TestDisabledMatchers(t *testing.T) {
	today := time.Date(2026, 10, 14, 0, 0, 0, 0, time.UTC)
	got, err := h.RenderString(Calendar(ModeSingle, Today(today), Disabled(Before(today), DayOfWeek{0})))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, `data-day="2026-10-13" data-disabled="true"`) {
		t.Error("days before today are disabled")
	}
	if !strings.Contains(got, `data-day="2026-10-18" data-disabled="true"`) {
		t.Error("Sundays are disabled")
	}
	if strings.Contains(got, `data-day="2026-10-15" data-disabled`) {
		t.Error("Thursday after today is enabled")
	}
}
