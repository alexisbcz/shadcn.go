package datatable

import "testing"

func TestUSD(t *testing.T) {
	for in, want := range map[float64]string{316: "$316.00", 1234567.5: "$1,234,567.50", -42.1: "-$42.10", 0: "$0.00"} {
		if got := usd(in); got != want {
			t.Errorf("usd(%v) = %q, want %q", in, got, want)
		}
	}
}
