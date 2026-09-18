package recipescale

import "testing"

func TestParseQuantity(t *testing.T) {
	cases := []struct {
		in   string
		want float64
	}{
		{"2", 2},
		{"0", 0},
		{"1/2", 0.5},
		{"3/4", 0.75},
		{"1 1/2", 1.5},
		{"2 3/8", 2.375},
		{"  1/2  ", 0.5},
		{"-1 1/2", -1.5},
		{"0.25", 0.25},
	}
	for _, c := range cases {
		got, err := ParseQuantity(c.in)
		if err != nil {
			t.Errorf("ParseQuantity(%q) returned error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("ParseQuantity(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestParseQuantityErrors(t *testing.T) {
	cases := []string{
		"",
		"   ",
		"abc",
		"1/0",
		"1 2 3",
		"a/2",
		"1/b",
		"1 abc",
	}
	for _, in := range cases {
		if _, err := ParseQuantity(in); err == nil {
			t.Errorf("ParseQuantity(%q) expected an error, got nil", in)
		}
	}
}

func TestParseQuantityRange(t *testing.T) {
	cases := []struct {
		in       string
		wantLow  float64
		wantHigh float64
	}{
		{"2", 2, 2},
		{"1/2", 0.5, 0.5},
		{"1-2", 1, 2},
		{"1/2-3/4", 0.5, 0.75},
		{"1 1/2-2 1/4", 1.5, 2.25},
		{"  1-2  ", 1, 2},
		{"-1 1/2", -1.5, -1.5},
	}
	for _, c := range cases {
		low, high, err := ParseQuantityRange(c.in)
		if err != nil {
			t.Errorf("ParseQuantityRange(%q) returned error: %v", c.in, err)
			continue
		}
		if low != c.wantLow || high != c.wantHigh {
			t.Errorf("ParseQuantityRange(%q) = (%v, %v), want (%v, %v)", c.in, low, high, c.wantLow, c.wantHigh)
		}
	}
}

func TestParseQuantityRangeErrors(t *testing.T) {
	cases := []string{
		"",
		"2-1",
		"abc-2",
		"1-abc",
		"1/0-2",
	}
	for _, in := range cases {
		if _, _, err := ParseQuantityRange(in); err == nil {
			t.Errorf("ParseQuantityRange(%q) expected an error, got nil", in)
		}
	}
}

func TestFormatQuantity(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{0, "0"},
		{2, "2"},
		{0.5, "1/2"},
		{1.5, "1 1/2"},
		{0.375, "3/8"},
		{2.375, "2 3/8"},
		{0.125, "1/8"},
		{0.0625, "1/8"},   // rounds up to nearest eighth
		{0.05, "0"},       // rounds down to nearest eighth
		{-1.5, "-1 1/2"},
		{2.999, "3"},
	}
	for _, c := range cases {
		got := FormatQuantity(c.in)
		if got != c.want {
			t.Errorf("FormatQuantity(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestParseFormatRoundTrip(t *testing.T) {
	// Values already on eighths should round-trip through parse -> format
	// back to the same string, modulo whitespace.
	cases := []string{"1/2", "1 1/2", "3/4", "2 3/8", "3"}
	for _, in := range cases {
		q, err := ParseQuantity(in)
		if err != nil {
			t.Fatalf("ParseQuantity(%q) returned error: %v", in, err)
		}
		if got := FormatQuantity(q); got != in {
			t.Errorf("FormatQuantity(ParseQuantity(%q)) = %q, want %q", in, got, in)
		}
	}
}
