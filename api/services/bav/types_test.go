package bav

import "testing"

func TestParseTriState(t *testing.T) {
	cases := []struct {
		name string
		in   any
		want TriState
	}{
		{"Y uppercase", "Y", TriYes},
		{"N uppercase", "N", TriNo},
		{"yes lowercase", "yes", TriYes},
		{"no lowercase", "no", TriNo},
		{"Yes mixed case", "Yes", TriYes},
		{"whitespace padded", "  y  ", TriYes},
		{"bool true", true, TriYes},
		{"bool false", false, TriNo},
		{"empty string", "", TriUnknown},
		{"literal unknown", "unknown", TriUnknown},
		{"nil", nil, TriUnknown},
		{"unrecognised string", "maybe", TriUnknown},
		{"already TriYes", TriYes, TriYes},
		{"already TriNo", TriNo, TriNo},
		{"already TriUnknown", TriUnknown, TriUnknown},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ParseTriState(tc.in); got != tc.want {
				t.Fatalf("ParseTriState(%v) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
