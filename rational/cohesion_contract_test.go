package rational_test

import (
	"errors"
	"strings"
	"testing"

	gomath "github.com/faustbrian/go-math"
	"github.com/faustbrian/go-math/rational"
)

func TestParseDistinguishesSyntaxAndComponentLimits(t *testing.T) {
	t.Parallel()

	limits := gomath.DefaultLimits()
	limits.MaxInputDigits = 3
	tests := []struct {
		name  string
		input string
		want  error
		text  string
	}{
		{name: "numerator limit", input: "1234/1", want: gomath.ErrLimitExceeded, text: "math: resource limit exceeded: rational numerator digits"},
		{name: "denominator limit", input: "1/1234", want: gomath.ErrLimitExceeded, text: "math: resource limit exceeded: rational denominator digits"},
		{name: "limit before later numerator syntax", input: "1234x/1", want: gomath.ErrLimitExceeded, text: "math: resource limit exceeded: rational numerator digits"},
		{name: "limit before later denominator syntax", input: "1/1234x", want: gomath.ErrLimitExceeded, text: "math: resource limit exceeded: rational denominator digits"},
		{name: "early numerator syntax", input: "12x4/1", want: gomath.ErrInvalidSyntax, text: "math: invalid syntax"},
		{name: "early denominator syntax", input: "1/12x4", want: gomath.ErrInvalidSyntax, text: "math: invalid syntax"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := rational.Parse(test.input, limits)
			assertRationalInputError(t, err, test.want, test.text)
		})
	}

	value, err := rational.Parse("123/123", limits)
	if err != nil || value.String() != "1" {
		t.Fatalf("Parse(exact boundary) = %s, %v", value, err)
	}
}

func TestParseLeadingZeroSyntaxPrecedesDigitBudget(t *testing.T) {
	t.Parallel()

	limits := gomath.DefaultLimits()
	limits.MaxInputDigits = 1
	for _, input := range []string{"00/1", "1/00"} {
		_, err := rational.Parse(input, limits)
		assertRationalInputError(t, err, gomath.ErrInvalidSyntax, "math: invalid syntax")
	}
}

func TestParseAppliesRawInputBoundBeforeGrammar(t *testing.T) {
	t.Parallel()

	limits := gomath.DefaultLimits()
	limits.MaxInputDigits = 3
	value, err := rational.Parse("-123/-123", limits)
	if err != nil || value.String() != "1" {
		t.Fatalf("Parse(raw boundary) = %s, %v", value, err)
	}

	_, err = rational.Parse(strings.Repeat("1", 10), limits)
	assertRationalInputError(
		t, err, gomath.ErrLimitExceeded,
		"math: resource limit exceeded: rational input bytes",
	)
	_, err = rational.Parse("123/-123x", limits)
	assertRationalInputError(t, err, gomath.ErrInvalidSyntax, "math: invalid syntax")

	invalid := limits
	invalid.MaxInputDigits = 0
	_, err = rational.Parse(strings.Repeat("1", 10), invalid)
	if !errors.Is(err, gomath.ErrInvalidArgument) ||
		errors.Is(err, gomath.ErrInvalidSyntax) || errors.Is(err, gomath.ErrLimitExceeded) {
		t.Fatalf("invalid limits precedence error = %v", err)
	}
}

func TestParseRawBoundCalculationSaturates(t *testing.T) {
	t.Parallel()

	limits := gomath.DefaultLimits()
	limits.MaxInputDigits = int(^uint(0) >> 1)
	value, err := rational.Parse("1/1", limits)
	if err != nil || value.String() != "1" {
		t.Fatalf("Parse(maximum limit) = %s, %v", value, err)
	}
}

func assertRationalInputError(t *testing.T, err, want error, text string) {
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
