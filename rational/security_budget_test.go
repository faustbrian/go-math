package rational_test

import (
	"errors"
	"testing"

	gomath "github.com/faustbrian/go-math"
	"github.com/faustbrian/go-math/rational"
)

func TestDecimalExpansionRejectsOversizedPowerOfTen(t *testing.T) {
	value, err := rational.New(1, 1)
	if err != nil {
		t.Fatal(err)
	}
	limits := gomath.DefaultLimits()
	limits.MaxIntermediateBits = 31
	text, _, err := value.Decimal(10, gomath.RoundHalfEven, limits)
	if !errors.Is(err, gomath.ErrLimitExceeded) {
		t.Fatalf("Decimal() = %q, %v; want ErrLimitExceeded", text, err)
	}
}

func TestDecimalExpansionRejectsMultiplicationCarryOverBudget(t *testing.T) {
	value, err := rational.New(3, 1)
	if err != nil {
		t.Fatal(err)
	}
	limits := gomath.DefaultLimits()
	limits.MaxIntermediateBits = 8
	text, _, err := value.Decimal(2, gomath.RoundHalfEven, limits)
	if !errors.Is(err, gomath.ErrLimitExceeded) {
		t.Fatalf("Decimal() = %q, %v; want ErrLimitExceeded", text, err)
	}
}
