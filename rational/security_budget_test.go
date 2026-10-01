package rational_test

import (
	"errors"
	"runtime"
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

func TestDecimalExpansionRejectsRaisedScaleBeforePowerConstruction(t *testing.T) {
	value, err := rational.New(1, 1)
	if err != nil {
		t.Fatal(err)
	}
	limits := gomath.DefaultLimits()
	limits.MaxDecimalExpansion = 100_001
	limits.MaxIntermediateBits = 1_000
	runtime.GC()
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	text, _, err := value.Decimal(limits.MaxDecimalExpansion, gomath.RoundDown, limits)
	runtime.ReadMemStats(&after)
	if !errors.Is(err, gomath.ErrLimitExceeded) {
		t.Fatalf("Decimal(raised scale) = %q, %v; want ErrLimitExceeded", text, err)
	}
	if allocated := after.TotalAlloc - before.TotalAlloc; allocated > 16_384 {
		t.Fatalf("Decimal(raised scale) allocated %d bytes before rejection", allocated)
	}
}
