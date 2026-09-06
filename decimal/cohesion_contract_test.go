package decimal_test

import (
	"errors"
	"strings"
	"testing"

	gomath "github.com/faustbrian/go-math"
	"github.com/faustbrian/go-math/decimal"
)

func TestParseUsesFirstConclusiveRejection(t *testing.T) {
	t.Parallel()

	limits := gomath.DefaultLimits()
	limits.MaxInputDigits = 3
	options := decimal.ParseOptions{AllowExponent: true, Limits: limits}
	tests := []struct {
		name  string
		input string
		want  error
		text  string
	}{
		{name: "digit limit before later empty exponent", input: "1234e+", want: gomath.ErrLimitExceeded, text: "math: resource limit exceeded: decimal input digits"},
		{name: "early mantissa syntax before exponent range", input: "12x4e2147483648", want: gomath.ErrInvalidSyntax, text: "math: invalid syntax"},
		{name: "sign-only exponent", input: "1e+", want: gomath.ErrInvalidSyntax, text: "math: invalid syntax"},
		{name: "invalid exponent within token", input: "1e+00000000x", want: gomath.ErrInvalidSyntax, text: "math: invalid syntax"},
		{name: "twelfth exponent byte", input: "1e+0000000000x", want: gomath.ErrLimitExceeded, text: "math: resource limit exceeded: decimal exponent"},
		{name: "host int32 exponent range", input: "1e2147483648", want: gomath.ErrLimitExceeded, text: "math: resource limit exceeded: decimal exponent"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := decimal.ParseWithOptions(test.input, options)
			assertDecimalInputError(t, err, test.want, test.text)
		})
	}

	valid, err := decimal.ParseWithOptions("1e+0000000001", options)
	if err != nil || valid.String() != "10" {
		t.Fatalf("ParseWithOptions(valid exponent boundary) = %s, %v", valid, err)
	}
	options.AllowUnderscores = true
	for _, input := range []string{"1_", "1._2"} {
		_, err = decimal.ParseWithOptions(input, options)
		assertDecimalInputError(t, err, gomath.ErrInvalidSyntax, "math: invalid syntax")
	}
}

func TestParseLeadingZeroSyntaxPrecedesDigitBudget(t *testing.T) {
	t.Parallel()

	limits := gomath.DefaultLimits()
	limits.MaxInputDigits = 1
	_, err := decimal.ParseWithOptions("00", decimal.ParseOptions{Limits: limits})
	assertDecimalInputError(t, err, gomath.ErrInvalidSyntax, "math: invalid syntax")
}

func TestParseSharesDigitBudgetAcrossMantissa(t *testing.T) {
	t.Parallel()

	limits := gomath.DefaultLimits()
	limits.MaxInputDigits = 3
	options := decimal.ParseOptions{Limits: limits}
	_, err := decimal.ParseWithOptions("12.34", options)
	assertDecimalInputError(
		t, err, gomath.ErrLimitExceeded,
		"math: resource limit exceeded: decimal input digits",
	)
	value, err := decimal.ParseWithOptions("12.3", options)
	if err != nil || value.String() != "12.3" {
		t.Fatalf("ParseWithOptions(exact digit boundary) = %s, %v", value, err)
	}
	_, err = decimal.ParseWithOptions("12.x4", options)
	assertDecimalInputError(t, err, gomath.ErrInvalidSyntax, "math: invalid syntax")
}

func TestParseAppliesRawInputBoundBeforeTrimming(t *testing.T) {
	t.Parallel()

	limits := gomath.DefaultLimits()
	limits.MaxInputDigits = 3
	options := decimal.ParseOptions{AllowWhitespace: true, Limits: limits}
	boundary := strings.Repeat(" ", 69) + "1"
	value, err := decimal.ParseWithOptions(boundary, options)
	if err != nil || value.String() != "1" {
		t.Fatalf("ParseWithOptions(raw boundary) = %s, %v", value, err)
	}
	_, err = decimal.ParseWithOptions(" "+boundary, options)
	assertDecimalInputError(
		t, err, gomath.ErrLimitExceeded,
		"math: resource limit exceeded: decimal input bytes",
	)
}

func TestParseChecksConfiguredExponentMagnitude(t *testing.T) {
	t.Parallel()

	limits := gomath.DefaultLimits()
	limits.MaxExponentMagnitude = 3
	_, err := decimal.ParseWithOptions(
		"1e4", decimal.ParseOptions{AllowExponent: true, Limits: limits},
	)
	assertDecimalInputError(
		t, err, gomath.ErrLimitExceeded,
		"math: resource limit exceeded: decimal exponent",
	)
}

func TestParseRawBoundCalculationSaturates(t *testing.T) {
	t.Parallel()

	limits := gomath.DefaultLimits()
	limits.MaxInputDigits = int(^uint(0) >> 1)
	value, err := decimal.ParseWithOptions("1", decimal.ParseOptions{Limits: limits})
	if err != nil || value.String() != "1" {
		t.Fatalf("ParseWithOptions(maximum limit) = %s, %v", value, err)
	}
}

func assertDecimalInputError(t *testing.T, err, want error, text string) {
	t.Helper()
	if err == nil || !errors.Is(err, want) || err.Error() != text {
		t.Fatalf("error = %v, want %q matching %v", err, text, want)
	}
	for _, category := range []error{
		gomath.ErrInvalidArgument, gomath.ErrInvalidSyntax, gomath.ErrLimitExceeded,
	} {
		if category != want && errors.Is(err, category) {
			t.Fatalf("error %v unexpectedly matches %v", err, category)
		}
	}
}
