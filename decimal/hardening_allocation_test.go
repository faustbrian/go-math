//go:build !race

package decimal_test

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"testing"

	gomath "github.com/faustbrian/go-math"
	"github.com/faustbrian/go-math/decimal"
)

func TestJSONRejectsHostileInputWithoutScalingAllocations(t *testing.T) {
	input := []byte(`"` + strings.Repeat("9", 8<<20) + `"`)
	result := testing.Benchmark(func(benchmark *testing.B) {
		for benchmark.Loop() {
			var value decimal.Decimal
			if err := value.UnmarshalJSON(input); !errors.Is(err, gomath.ErrLimitExceeded) {
				panic("hostile JSON did not retain its limit identity")
			}
		}
	})
	if allocated := result.AllocedBytesPerOp(); allocated > 4096 {
		t.Fatalf("hostile JSON allocated %d bytes, want at most 4096", allocated)
	}
}

func TestUnderflowRoundingDoesNotMaterializeOversizedDivisor(t *testing.T) {
	limits := gomath.DefaultLimits()
	value, err := decimal.FromBig(big.NewInt(1), -10_000, limits)
	if err != nil {
		t.Fatal(err)
	}
	limits.MaxIntermediateBits = 8
	operation := decimal.Context{
		Precision: 1, MinExponent: 0, MaxExponent: 0,
		Rounding: decimal.HalfEven, Limits: limits,
	}
	result := testing.Benchmark(func(benchmark *testing.B) {
		for benchmark.Loop() {
			result, err := operation.Apply(context.Background(), value)
			if err != nil || !result.Value.IsZero() || !result.Conditions.Has(gomath.ConditionUnderflow) {
				panic("underflow did not retain zero and its condition")
			}
		}
	})
	if allocated := result.AllocedBytesPerOp(); allocated > 4096 {
		t.Fatalf("underflow rounding allocated %d bytes, want at most 4096", allocated)
	}
}

// Race instrumentation adds size-dependent allocations that do not exist in
// production, so allocation bounds are measured only by the normal test build.
func TestParserRejectsHostileInputWithoutScalingAllocations(t *testing.T) {
	limits := gomath.DefaultLimits()
	limits.MaxInputDigits = 8
	measureBytes := func(input string, want error) int64 {
		result := testing.Benchmark(func(benchmark *testing.B) {
			for benchmark.Loop() {
				if _, err := decimal.ParseWithOptions(input, decimal.ParseOptions{Limits: limits}); !errors.Is(err, want) {
					panic("hostile decimal did not retain its error identity")
				}
			}
		})

		return result.AllocedBytesPerOp()
	}
	boundaryDigits := measureBytes(strings.Repeat("9", limits.MaxInputDigits+1), gomath.ErrLimitExceeded)
	hostileDigits := measureBytes(strings.Repeat("9", 1<<20), gomath.ErrLimitExceeded)
	if hostileDigits > boundaryDigits+1024 {
		t.Fatalf(
			"attacker-sized input allocated %d bytes; near-boundary rejection allocated %d",
			hostileDigits, boundaryDigits,
		)
	}
	boundarySeparators := measureBytes("1..", decimal.ErrInvalid)
	hostileSeparators := measureBytes(strings.Repeat(".", 1<<20), gomath.ErrLimitExceeded)
	if hostileSeparators > boundarySeparators+1024 {
		t.Fatalf(
			"attacker-sized separators allocated %d bytes; near-boundary rejection allocated %d",
			hostileSeparators, boundarySeparators,
		)
	}
}
