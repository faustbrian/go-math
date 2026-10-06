package decimal

import (
	"context"
	"math"
	"math/big"
	"testing"

	gomath "github.com/faustbrian/go-math"
)

func TestRatioComparisonPreservesLargeExponent(t *testing.T) {
	for _, test := range []struct {
		name     string
		exponent int64
		want     int
	}{
		{"positive", 1 << 32, -1},
		{"negative", -(1 << 32), 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := compareRatioPower10(big.NewInt(1), big.NewInt(1), test.exponent); got != test.want {
				t.Fatalf("comparison with exponent %d = %d, want %d", test.exponent, got, test.want)
			}
		})
	}
}

func TestRatioComparisonBoundariesAndOwnership(t *testing.T) {
	for _, test := range []struct {
		name                   string
		numerator, denominator int64
		exponent               int64
		want                   int
	}{
		{"maximum exponent", 1, 1, math.MaxInt64, -1},
		{"minimum exponent", 1, 1, math.MinInt64, 1},
		{"zero equal", 0, 0, math.MinInt64, 0},
		{"zero numerator", 0, 1, math.MinInt64, -1},
		{"zero denominator", 1, 0, math.MaxInt64, 1},
		{"negative numerator", -1, 0, math.MaxInt64, -1},
		{"negative denominator", 0, -1, math.MinInt64, 1},
		{"unscaled equal", 12, 12, 0, 0},
		{"unscaled below", 11, 12, 0, -1},
		{"unscaled above", 13, 12, 0, 1},
		{"positive equal", 1200, 12, 2, 0},
		{"positive below", 1199, 12, 2, -1},
		{"positive above", 1201, 12, 2, 1},
		{"negative equal", 12, 1200, -2, 0},
		{"negative below", 12, 1201, -2, -1},
		{"negative above", 12, 1199, -2, 1},
		{"short left", 1, 10, 0, -1},
		{"short right", 10, 1, 0, 1},
		{"signed equal", -1200, -12, 2, 0},
		{"signed below", -1201, -12, 2, -1},
		{"signed above", -1199, -12, 2, 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			numerator, denominator := big.NewInt(test.numerator), big.NewInt(test.denominator)
			if got := compareRatioPower10(numerator, denominator, test.exponent); got != test.want {
				t.Fatalf("comparison = %d, want %d", got, test.want)
			}
			if numerator.Int64() != test.numerator || denominator.Int64() != test.denominator {
				t.Fatal("comparison mutated an operand")
			}
		})
	}
}

func TestRatioComparisonPublicQuotient(t *testing.T) {
	operation := Context{Precision: 3, MinExponent: -20, MaxExponent: 20, Rounding: HalfEven}
	for _, test := range []struct {
		name                   string
		numerator, denominator int64
		want                   string
		conditions             gomath.Condition
	}{
		{"equal ratio", 1, 1, "1", 0},
		{"negative equal ratio", -1, 1, "-1", 0},
		{"exact small ratio", 1, 8, "0.125", 0},
		{"exact large ratio", 10, 2, "5", 0},
		{"repeating", 1, 3, "0.333", gomath.ConditionRounded | gomath.ConditionInexact},
		{"negative repeating", -1, 3, "-0.333", gomath.ConditionRounded | gomath.ConditionInexact},
		{"exact differing digit counts", 999, 100, "9.99", 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			result, err := operation.Quo(context.Background(), New(test.numerator), New(test.denominator))
			if err != nil {
				t.Fatalf("Quo() error = %v", err)
			}
			if result.Value.String() != test.want || result.Conditions != test.conditions {
				t.Fatalf("Quo() = %s [%s], want %s [%s]", result.Value, result.Conditions, test.want, test.conditions)
			}
		})
	}
}
