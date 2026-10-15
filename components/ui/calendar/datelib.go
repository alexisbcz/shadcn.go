package calendar

import (
	"fmt"
	"strings"
	"time"
)

// Dates are handled as Julian Day Numbers (whole days, no time zone), the
// representation calendar.js uses too. The conversions are those of
// date-fns-jalali (and the proleptic Gregorian calendar).

func div(a, b int) int { return a / b } // truncating, like ~~(a / b)
func mod(a, b int) int { return a - (a/b)*b }

func fdiv(a, b int) int {
	q := a / b
	if (a%b != 0) && ((a < 0) != (b < 0)) {
		q--
	}
	return q
}

func fmod(a, b int) int { return ((a % b) + b) % b }

func g2d(gy, gm, gd int) int {
	return div(1461*(gy+4800+div(gm-14, 12)), 4) +
		div(367*(gm-2-12*div(gm-14, 12)), 12) -
		div(3*div(gy+4900+div(gm-14, 12), 100), 4) +
		gd - 32075
}

func d2g(jdn int) (int, int, int) {
	L := jdn + 68569
	n := div(4*L, 146097)
	L = L - div(146097*n+3, 4)
	i := div(4000*(L+1), 1461001)
	L = L - div(1461*i, 4) + 31
	j := div(80*L, 2447)
	gd := L - div(2447*j, 80)
	L = div(j, 11)
	gm := j + 2 - 12*L
	gy := 100*(n-49) + i + L
	return gy, gm, gd
}

const persianEpoch = 1948320

var persianNumDays = [12]int{0, 31, 62, 93, 124, 155, 186, 216, 246, 276, 306, 336}

func j2d(jy, jm, jd int) int {
	month := jm - 1
	julianDay := persianEpoch - 1 + 365*(jy-1) + div(8*jy+21, 33)
	if month != 0 {
		julianDay += persianNumDays[month]
	}
	return julianDay + jd
}

func d2j(julianDay int) (int, int, int) {
	daysSinceEpoch := julianDay - persianEpoch
	year := 1 + div(33*daysSinceEpoch+3, 12053)
	dayOfYear := daysSinceEpoch - (365*(year-1) + div(8*year+21, 33))
	if dayOfYear < 0 {
		year--
		dayOfYear = daysSinceEpoch - (365*(year-1) + div(8*year+21, 33))
	}
	var month int
	if dayOfYear < 216 {
		month = div(dayOfYear, 31)
	} else {
		month = div(dayOfYear-6, 30)
	}
	return year, month + 1, dayOfYear - persianNumDays[month] + 1
}

func isLeapJalaliYear(jy int) bool {
	m := mod(25*jy+11, 33)
	return (m < 8 && m >= -1) || m <= -27
}

// system maps a JDN to (year, month 0-11, day) and back.
type system interface {
	ymd(j int) (int, int, int)
	jdn(y, m, d int) int
	dim(y, m int) int
}

type gregorian struct{}

func (gregorian) ymd(j int) (int, int, int) {
	y, m, d := d2g(j)
	return y, m - 1, d
}

func (gregorian) jdn(y, m, d int) int { return g2d(y+fdiv(m, 12), fmod(m, 12)+1, d) }

func (gregorian) dim(y, m int) int {
	ny, nm := y+fdiv(m, 12), fmod(m, 12)
	return g2d(ny+fdiv(nm+1, 12), fmod(nm+1, 12)+1, 1) - g2d(ny, nm+1, 1)
}

type jalali struct{}

func (jalali) ymd(j int) (int, int, int) {
	y, m, d := d2j(j)
	return y, m - 1, d
}

func (jalali) jdn(y, m, d int) int { return j2d(y+fdiv(m, 12), fmod(m, 12)+1, d) }

func (jalali) dim(y, m int) int {
	ny, nm := y+fdiv(m, 12), fmod(m, 12)
	switch {
	case nm < 6:
		return 31
	case nm < 11:
		return 30
	case isLeapJalaliYear(ny):
		return 30
	}
	return 29
}

// JDN returns the Julian Day Number of t's calendar date (in t's location).
func JDN(t time.Time) int {
	y, m, d := t.Date()
	return g2d(y, int(m), d)
}

// isoOf formats a JDN as a Gregorian "YYYY-MM-DD".
func isoOf(j int) string {
	y, m, d := d2g(j)
	return fmt.Sprintf("%04d-%02d-%02d", y, m, d)
}

var (
	enMonths      = [12]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	enMonthsShort = [12]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	enDays        = [7]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	enDaysShort   = [7]string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"}
	faMonths      = [12]string{"فروردین", "اردیبهشت", "خرداد", "تیر", "مرداد", "شهریور", "مهر", "آبان", "آذر", "دی", "بهمن", "اسفند"}
	faDays        = [7]string{"یک‌شنبه", "دوشنبه", "سه‌شنبه", "چهارشنبه", "پنج‌شنبه", "جمعه", "شنبه"}
	faDaysShort   = [7]string{"1ش", "2ش", "3ش", "4ش", "5ش", "ج", "ش"}
	arMonths      = [12]string{"يناير", "فبراير", "مارس", "أبريل", "مايو", "يونيو", "يوليو", "أغسطس", "سبتمبر", "أكتوبر", "نوفمبر", "ديسمبر"}
	arDays        = [7]string{"الأحد", "الاثنين", "الثلاثاء", "الأربعاء", "الخميس", "الجمعة", "السبت"}
	arDaysShort   = [7]string{"أحد", "اثنين", "ثلاثاء", "أربعاء", "خميس", "جمعة", "سبت"}
)

var numeralDigits = map[Numerals]string{
	"arab":    "٠١٢٣٤٥٦٧٨٩",
	"arabext": "۰۱۲۳۴۵۶۷۸۹",
	"deva":    "०१२३४५६७८९",
	"beng":    "০১২৩৪৫৬৭৮৯",
}

func enOrdinal(n int) string {
	r100 := n % 100
	if r100 > 20 || r100 < 10 {
		switch r100 % 10 {
		case 1:
			return fmt.Sprintf("%dst", n)
		case 2:
			return fmt.Sprintf("%dnd", n)
		case 3:
			return fmt.Sprintf("%drd", n)
		}
	}
	return fmt.Sprintf("%dth", n)
}

// dateLib is react-day-picker's DateLib over JDNs, for one calendar system,
// locale and numbering system.
type dateLib struct {
	s            system
	fa           bool
	ar           bool // arSA names (Gregorian, Latin digits)
	digits       []rune
	weekStartsOn int
}

func newDateLib(persian bool, numerals Numerals, weekStartsOn *int) *dateLib {
	l := &dateLib{s: gregorian{}}
	if persian {
		l.s = jalali{}
		l.fa = true
		l.weekStartsOn = 6
	}
	if d, ok := numeralDigits[numerals]; ok {
		l.digits = []rune(d)
	}
	if weekStartsOn != nil {
		l.weekStartsOn = *weekStartsOn
	}
	return l
}

func (l *dateLib) num(s string) string {
	if l.digits == nil {
		return s
	}
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(l.digits[r-'0'])
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func (l *dateLib) weekday(j int) int { return fmod(j+1, 7) }

func (l *dateLib) startOfMonth(j int) int {
	y, m, _ := l.s.ymd(j)
	return l.s.jdn(y, m, 1)
}

func (l *dateLib) endOfMonth(j int) int {
	y, m, _ := l.s.ymd(j)
	return l.s.jdn(y, m, l.s.dim(y, m))
}

func (l *dateLib) startOfYear(j int) int {
	y, _, _ := l.s.ymd(j)
	return l.s.jdn(y, 0, 1)
}

func (l *dateLib) endOfYear(j int) int {
	y, _, _ := l.s.ymd(j)
	return l.s.jdn(y, 11, l.s.dim(y, 11))
}

func (l *dateLib) addMonths(j, n int) int {
	y, m, d := l.s.ymd(j)
	ny, nm := y+fdiv(m+n, 12), fmod(m+n, 12)
	return l.s.jdn(ny, nm, min(d, l.s.dim(ny, nm)))
}

func (l *dateLib) addYears(j, n int) int { return l.addMonths(j, 12*n) }

func (l *dateLib) getMonth(j int) int {
	_, m, _ := l.s.ymd(j)
	return m
}

func (l *dateLib) getYear(j int) int {
	y, _, _ := l.s.ymd(j)
	return y
}

func (l *dateLib) diffMonths(a, b int) int {
	ya, ma, _ := l.s.ymd(a)
	yb, mb, _ := l.s.ymd(b)
	return ya*12 + ma - (yb*12 + mb)
}

func (l *dateLib) isSameMonth(a, b int) bool { return l.diffMonths(a, b) == 0 }
func (l *dateLib) isSameYear(a, b int) bool  { return l.getYear(a) == l.getYear(b) }

func (l *dateLib) startOfWeek(j int) int { return j - fmod(l.weekday(j)-l.weekStartsOn, 7) }
func (l *dateLib) endOfWeek(j int) int   { return l.startOfWeek(j) + 6 }

func (l *dateLib) getWeek(j int) int {
	y := l.getYear(j)
	sow := func(yy int) int { return l.startOfWeek(l.s.jdn(yy, 0, 1)) }
	start := sow(y)
	if j >= sow(y+1) {
		start = sow(y + 1)
	} else if j < start {
		start = sow(y - 1)
	}
	return fdiv(l.startOfWeek(j)-start, 7) + 1
}

// Formats ("LLLL y", "PPPP", "cccccc", "cccc", "d", "yyyy-MM-dd", ...).

func (l *dateLib) caption(j int) string {
	y, m, _ := l.s.ymd(j)
	if l.fa {
		return l.num(fmt.Sprintf("%s %d", faMonths[m], y))
	}
	if l.ar {
		return l.num(fmt.Sprintf("%s %d", arMonths[m], y))
	}
	return l.num(fmt.Sprintf("%s %d", enMonths[m], y))
}

func (l *dateLib) fullDate(j int) string {
	y, m, d := l.s.ymd(j)
	w := l.weekday(j)
	if l.fa {
		return l.num(fmt.Sprintf("%s %d-ام %s %d", faDays[w], d, faMonths[m], y))
	}
	if l.ar {
		return l.num(fmt.Sprintf("%s, %s %d, %d", arDays[w], arMonths[m], d, y))
	}
	return l.num(fmt.Sprintf("%s, %s %s, %d", enDays[w], enMonths[m], enOrdinal(d), y))
}

func (l *dateLib) weekdayShort(j int) string {
	if l.fa {
		return l.num(faDaysShort[l.weekday(j)])
	}
	if l.ar {
		return l.num(arDaysShort[l.weekday(j)])
	}
	return l.num(enDaysShort[l.weekday(j)])
}

func (l *dateLib) weekdayLong(j int) string {
	if l.fa {
		return l.num(faDays[l.weekday(j)])
	}
	if l.ar {
		return l.num(arDays[l.weekday(j)])
	}
	return l.num(enDays[l.weekday(j)])
}

func (l *dateLib) day(j int) string {
	_, _, d := l.s.ymd(j)
	return l.num(fmt.Sprint(d))
}

func (l *dateLib) isoDay(j int) string {
	y, m, d := l.s.ymd(j)
	return l.num(fmt.Sprintf("%04d-%02d-%02d", y, m+1, d))
}

func (l *dateLib) isoMonth(j int) string {
	y, m, _ := l.s.ymd(j)
	return l.num(fmt.Sprintf("%04d-%02d", y, m+1))
}

func (l *dateLib) year(j int) string { return l.num(fmt.Sprintf("%04d", l.getYear(j))) }

// monthDropdown is shadcn's formatMonthDropdown:
// date.toLocaleString("default", { month: "short" }) in en-US.
func (l *dateLib) monthDropdown(j int) string {
	_, gm, _ := d2g(j)
	if l.ar {
		return arMonths[gm-1]
	}
	return enMonthsShort[gm-1]
}

// localeDate is date.toLocaleDateString() in en-US.
func (l *dateLib) localeDate(j int) string {
	gy, gm, gd := d2g(j)
	if l.ar {
		// ar-SA: Arabic-Indic digits, day first, RLM before each slash.
		return arabicDigits(fmt.Sprintf("%d\u200f/%d\u200f/%d", gd, gm, gy))
	}
	return fmt.Sprintf("%d/%d/%d", gm, gd, gy)
}

func (l *dateLib) weekNumber(n int) string {
	if n < 10 {
		return l.num(fmt.Sprintf("0%d", n))
	}
	return l.num(fmt.Sprint(n))
}

func arabicDigits(s string) string {
	d := []rune(numeralDigits["arab"])
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			r = d[r-'0']
		}
		b.WriteRune(r)
	}
	return b.String()
}
