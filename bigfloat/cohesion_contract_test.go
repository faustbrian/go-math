package bigfloat_test

import (
	"errors"
	"testing"

	gomath "github.com/faustbrian/go-math"
	"github.com/faustbrian/go-math/bigfloat"
)

func TestParseDistinguishesInputLimitFromSyntax(t *testing.T) {
	t.Parallel()

	operation := testContext()
	operation.Limits.MaxInputDigits = 4
	_, err := bigfloat.Parse("12345", 10, operation)
	assertBigFloatInputError(
		t, err, gomath.ErrLimitExceeded,
		"math: resource limit exceeded: binary float input",
	)
	_, err = bigfloat.Parse("xxxx", 10, operation)
	assertBigFloatInputError(t, err, gomath.ErrInvalidSyntax, "math: invalid syntax")
	result, err := bigfloat.Parse("1.25", 10, operation)
	if err != nil || result.Value.String() != "1.25" {
		t.Fatalf("Parse(exact boundary) = %s, %v", result.Value, err)
	}
}

func TestParseRejectsInvalidBaseBeforeInspectingInput(t *testing.T) {
	t.Parallel()

	operation := testContext()
	operation.Limits.MaxInputDigits = 4
	for _, input := range []string{"", "x", "12345", "1"} {
		_, err := bigfloat.Parse(input, 8, operation)
		assertBigFloatInputError(
			t, err, gomath.ErrInvalidArgument,
			"math: invalid argument: float base",
		)
	}
}

func TestOperationsRejectUnsupportedRoundingExactly(t *testing.T) {
	t.Parallel()

	operation := testContext()
	operation.Rounding = gomath.RoundHalfDown
	_, err := bigfloat.NewInt64(1, operation)
	assertBigFloatInputError(
		t, err, gomath.ErrInvalidArgument,
		"math: invalid argument: binary rounding mode",
	)
}

func assertBigFloatInputError(t *testing.T, err, want error, text string) {
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
