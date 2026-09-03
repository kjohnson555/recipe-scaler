package recipescale

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// nearestEighth is the smallest fraction FormatQuantity will round to.
// Measuring cups and spoons come in eighths, so rounding any finer than
// that produces numbers no one can actually measure out.
const nearestEighth = 8

// ParseQuantity parses a recipe quantity such as "2", "1/2", or "1 1/2"
// into a decimal value.
func ParseQuantity(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty quantity")
	}

	parts := strings.Fields(s)
	switch len(parts) {
	case 1:
		return parseFractionPart(parts[0])
	case 2:
		whole, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid quantity %q: %w", s, err)
		}
		frac, err := parseFractionPart(parts[1])
		if err != nil {
			return 0, fmt.Errorf("invalid quantity %q: %w", s, err)
		}
		if whole < 0 {
			return whole - frac, nil
		}
		return whole + frac, nil
	default:
		return 0, fmt.Errorf("invalid quantity %q", s)
	}
}

func parseFractionPart(s string) (float64, error) {
	if i := strings.IndexByte(s, '/'); i >= 0 {
		num, err := strconv.Atoi(s[:i])
		if err != nil {
			return 0, fmt.Errorf("invalid fraction %q: %w", s, err)
		}
		denom, err := strconv.Atoi(s[i+1:])
		if err != nil {
			return 0, fmt.Errorf("invalid fraction %q: %w", s, err)
		}
		if denom == 0 {
			return 0, fmt.Errorf("invalid fraction %q: division by zero", s)
		}
		return float64(num) / float64(denom), nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid quantity %q: %w", s, err)
	}
	return v, nil
}

// FormatQuantity renders a decimal quantity as a whole number plus a
// simplified eighth fraction, e.g. 1.375 -> "1 3/8". Scaling a recipe
// almost never lands on a clean decimal, and "0.666666 cups" is useless
// at the counter.
func FormatQuantity(q float64) string {
	sign := ""
	if q < 0 {
		sign = "-"
		q = -q
	}

	rounded := math.Round(q*nearestEighth) / nearestEighth
	whole := math.Trunc(rounded)
	frac := rounded - whole

	if frac == 0 {
		return sign + strconv.FormatFloat(whole, 'f', -1, 64)
	}

	num := int(math.Round(frac * nearestEighth))
	den := nearestEighth
	if g := gcd(num, den); g > 1 {
		num /= g
		den /= g
	}

	if whole == 0 {
		return fmt.Sprintf("%s%d/%d", sign, num, den)
	}
	return fmt.Sprintf("%s%d %d/%d", sign, int(whole), num, den)
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
