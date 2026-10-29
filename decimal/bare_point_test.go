package decimal_test

import (
	"errors"
	"math/big"
	"regexp"
	"strings"
	"testing"

	gomath "github.com/faustbrian/go-math"
	"github.com/faustbrian/go-math/decimal"
)

const (
	needExponent = 1 << iota
	needPlus
	needUnderscores
	needLeadingZeros
	needWhitespace
	needBarePoint
	allGrammarFlags = 1<<6 - 1
)

func grammarOptions(flags int, limits gomath.Limits) decimal.ParseOptions {
	return decimal.ParseOptions{
		AllowExponent:     flags&needExponent != 0,
		AllowPlus:         flags&needPlus != 0,
		AllowUnderscores:  flags&needUnderscores != 0,
		AllowLeadingZeros: flags&needLeadingZeros != 0,
		AllowWhitespace:   flags&needWhitespace != 0,
		AllowBarePoint:    flags&needBarePoint != 0,
		Limits:            limits,
	}
}

func TestParseBarePointAcceptsOneSidedMantissaUnderEveryOptionCombination(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		needs    int
		value    string
		exponent int32
	}{
		{".5", needBarePoint, "0.5", -1},
		{"5.", needBarePoint, "5", 0},
		{"-.5", needBarePoint, "-0.5", -1},
		{"-5.", needBarePoint, "-5", 0},
		{"+.5", needBarePoint | needPlus, "0.5", -1},
		{"+5.", needBarePoint | needPlus, "5", 0},
		{"0.", needBarePoint, "0", 0},
		{".0", needBarePoint, "0.0", -1},
		{".50", needBarePoint, "0.50", -2},
		{".05", needBarePoint, "0.05", -2},
		{"5.e2", needBarePoint | needExponent, "500", 2},
		{".5e2", needBarePoint | needExponent, "50", 1},
		{"5.E-1", needBarePoint | needExponent, "0.5", -1},
		{"-.5e+1", needBarePoint | needExponent, "-5", 0},
		{"00.5", needLeadingZeros, "0.5", -1},
		{"0_0.5", needLeadingZeros | needUnderscores, "0.5", -1},
		{"00.", needBarePoint | needLeadingZeros, "0", 0},
		{"1_0.", needBarePoint | needUnderscores, "10", 0},
		{".5_5", needBarePoint | needUnderscores, "0.55", -2},
		{" .5 ", needBarePoint | needWhitespace, "0.5", -1},
		{"\t5.\n", needBarePoint | needWhitespace, "5", 0},
		{"+1_0.e+1", needBarePoint | needPlus | needUnderscores | needExponent, "100", 1},
		{" -00.5_0E1 ", needWhitespace | needLeadingZeros | needUnderscores | needExponent, "-5.0", -1},
		{" -01_0.E1 ", needBarePoint | needWhitespace | needLeadingZeros | needUnderscores | needExponent, "-100", 1},
		{"+00.5_0e-1", needPlus | needLeadingZeros | needUnderscores | needExponent, "0.050", -3},
		// Two-sided mantissas keep their meaning when the option is enabled.
		{"1.5", 0, "1.5", -1},
		{"-0.25", 0, "-0.25", -2},
		{"12", 0, "12", 0},
	}
	limits := gomath.DefaultLimits()
	for _, test := range tests {
		for flags := 0; flags <= allGrammarFlags; flags++ {
			got, err := decimal.ParseWithOptions(test.input, grammarOptions(flags, limits))
			if flags&test.needs != test.needs {
				if !errors.Is(err, decimal.ErrInvalid) {
					t.Fatalf("%q with flags %06b: got %s, %v; want ErrInvalid", test.input, flags, got, err)
				}
				continue
			}
			if err != nil {
				t.Fatalf("%q with flags %06b: %v", test.input, flags, err)
			}
			if got.String() != test.value || got.Exponent() != test.exponent || got.Scale() != -test.exponent {
				t.Fatalf("%q with flags %06b: got %s exponent %d scale %d; want %s exponent %d",
					test.input, flags, got, got.Exponent(), got.Scale(), test.value, test.exponent)
			}
		}
	}
}

func TestParseBarePointRejectsMantissaWithoutDigitsOrMisplacedSeparators(t *testing.T) {
	t.Parallel()

	limits := gomath.DefaultLimits()
	for _, input := range []string{
		".", "-.", "+.", " . ", ".e3", "-.e3", ".E3", "+.e+3", "..5", "5..", ".5.",
		"._5", "5_.", "_.5", "._", "5._", "_5.", ".5_", "5_.e3", "._5e3", ".5_e3",
		". 5", "5 .", ".5e", "5.e", "5.e+", "e5.", ".-5", "-+.5", "0x.5", "0x1p-2",
		"1/2", ".5/2", ".5.e3", "5.e3.", ". 5", "٥.",
	} {
		for flags := 0; flags <= allGrammarFlags; flags++ {
			if got, err := decimal.ParseWithOptions(input, grammarOptions(flags, limits)); !errors.Is(err, decimal.ErrInvalid) {
				t.Fatalf("%q with flags %06b: got %s, %v; want ErrInvalid", input, flags, got, err)
			}
		}
	}
}

func TestParseBarePointSharesTheMantissaDigitBudget(t *testing.T) {
	t.Parallel()

	limits := gomath.DefaultLimits()
	limits.MaxInputDigits = 3
	options := decimal.ParseOptions{AllowBarePoint: true, Limits: limits}
	for _, input := range []string{".123", "123."} {
		if _, err := decimal.ParseWithOptions(input, options); err != nil {
			t.Fatalf("%q within budget: %v", input, err)
		}
	}
	for _, input := range []string{".1234", "1234."} {
		if _, err := decimal.ParseWithOptions(input, options); !errors.Is(err, decimal.ErrLimit) {
			t.Fatalf("%q over budget: %v; want ErrLimit", input, err)
		}
	}
	// The first conclusive rejection wins: the budget is crossed before the
	// malformed exponent is reached, and a misplaced separator precedes it.
	options.AllowExponent = true
	options.AllowUnderscores = true
	for input, want := range map[string]error{
		".1234e+": decimal.ErrLimit,
		"1234.e":  decimal.ErrLimit,
		"._1234":  decimal.ErrInvalid,
		".e1234":  decimal.ErrInvalid,
	} {
		if _, err := decimal.ParseWithOptions(input, options); !errors.Is(err, want) {
			t.Fatalf("%q: %v; want %v", input, err, want)
		}
	}
}

// FuzzParseBarePointOracle compares ParseWithOptions with an independent
// regular-expression grammar for every flag combination and with big.Rat for
// accepted values. A grammar match may fail only when the input could cross a
// configured limit.
func FuzzParseBarePointOracle(f *testing.F) {
	for _, seed := range []string{
		".5", "5.", "-.5", "+.5", "5.e3", ".5e-2", "00.5", "1_0.", "._5", "5_.",
		".", ".e3", " .50 ", "0x.5", "1/2", "-00.5_0E1",
	} {
		f.Add(seed, uint8(allGrammarFlags))
		f.Add(seed, uint8(allGrammarFlags&^needBarePoint))
	}
	for flags := 0; flags <= allGrammarFlags; flags++ {
		for _, seed := range []string{" +0_0.5e-1 ", "-.5", "5.e2", "05.", "1_0"} {
			f.Add(seed, uint8(flags))
		}
	}
	limits := gomath.DefaultLimits()
	limits.MaxInputDigits = 64
	limits.MaxOutputDigits = 256
	limits.MaxExponentMagnitude = 128
	grammars := make([]*regexp.Regexp, allGrammarFlags+1)
	for flags := range grammars {
		grammars[flags] = referenceGrammar(flags)
	}
	f.Fuzz(func(t *testing.T, input string, rawFlags uint8) {
		flags := int(rawFlags) & allGrammarFlags
		got, err := decimal.ParseWithOptions(input, grammarOptions(flags, limits))
		trimmed := input
		if flags&needWhitespace != 0 {
			trimmed = strings.TrimSpace(input)
		}
		if !grammars[flags].MatchString(trimmed) {
			if err == nil {
				t.Fatalf("%q with flags %06b accepted as %s outside the reference grammar", input, flags, got)
			}
			return
		}
		cleaned := strings.ReplaceAll(trimmed, "_", "")
		mantissa, exponentText, _ := strings.Cut(strings.ToLower(cleaned), "e")
		integerPart, fraction, _ := strings.Cut(strings.TrimLeft(mantissa, "+-"), ".")
		exponent := new(big.Int)
		if exponentText != "" {
			exponent.SetString(exponentText, 10)
		}
		final := new(big.Int).Sub(exponent, big.NewInt(int64(len(fraction))))
		mantissaDigits := len(integerPart) + len(fraction)
		// Each bound mirrors one documented ParseOptions limit; the output bound
		// over-approximates rendered digits so it never excuses a valid input.
		withinLimits := len(input) <= 2*limits.MaxInputDigits+64 &&
			mantissaDigits <= limits.MaxInputDigits && len(exponentText) <= 11 &&
			exponent.CmpAbs(big.NewInt(int64(limits.MaxExponentMagnitude))) <= 0 &&
			final.CmpAbs(big.NewInt(int64(limits.MaxExponentMagnitude))) <= 0 &&
			mantissaDigits+int(new(big.Int).Abs(final).Int64())+1 <= limits.MaxOutputDigits
		if err != nil {
			if withinLimits || !errors.Is(err, decimal.ErrLimit) {
				t.Fatalf("%q with flags %06b rejected inside the reference grammar: %v", input, flags, err)
			}
			return
		}
		want, ok := new(big.Rat).SetString(cleaned)
		if !ok {
			t.Fatalf("reference big.Rat rejected %q", cleaned)
		}
		if got.BigRat().Cmp(want) != 0 {
			t.Fatalf("%q with flags %06b: got %s, want %s", input, flags, got, want.RatString())
		}
		if !final.IsInt64() || int64(got.Exponent()) != final.Int64() {
			t.Fatalf("%q with flags %06b: exponent %d, want %s", input, flags, got.Exponent(), final)
		}
	})
}

// referenceGrammar spells the documented option grammar directly, independent
// of the byte scanner.
func referenceGrammar(flags int) *regexp.Regexp {
	digits := "[0-9]+"
	if flags&needUnderscores != 0 {
		digits = "[0-9](?:_?[0-9])*"
	}
	integer := "(?:0|[1-9]" + strings.TrimPrefix(digits, "[0-9]") + ")"
	if flags&needUnderscores == 0 {
		integer = "(?:0|[1-9][0-9]*)"
	}
	if flags&needLeadingZeros != 0 {
		integer = digits
	}
	sign := "-?"
	if flags&needPlus != 0 {
		sign = "[-+]?"
	}
	mantissa := integer + "(?:\\." + digits + ")?"
	if flags&needBarePoint != 0 {
		mantissa = "(?:" + integer + "(?:\\.(?:" + digits + ")?)?|\\." + digits + ")"
	}
	exponent := ""
	if flags&needExponent != 0 {
		exponent = "(?:[eE][-+]?[0-9]+)?"
	}

	return regexp.MustCompile("^" + sign + mantissa + exponent + "$")
}
