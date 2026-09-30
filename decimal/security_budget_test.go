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

func TestOverflowCoefficientRespectsIntermediateBudget(t *testing.T) {
	limits := gomath.DefaultLimits()
	limits.MaxIntermediateBits = 8
	operation := decimal.Context{
		Precision: 8, MinExponent: -2, MaxExponent: -1,
		Rounding: decimal.HalfEven, Limits: limits,
	}
	result, err := operation.Apply(context.Background(), decimal.New(1))
	if !errors.Is(err, gomath.ErrLimitExceeded) {
		t.Fatalf("overflow coefficient bits = %d, error = %v; want ErrLimitExceeded", result.Value.Coefficient().BitLen(), err)
	}
}

func TestUnderflowShortcutPreservesRoundingAndConditions(t *testing.T) {
	limits := gomath.DefaultLimits()
	limits.MaxIntermediateBits = 8
	for _, sign := range []int64{-1, 1} {
		value, err := decimal.FromBig(big.NewInt(sign), -10_000, limits)
		if err != nil {
			t.Fatal(err)
		}
		for _, mode := range []gomath.RoundingMode{
			gomath.RoundHalfEven, gomath.RoundHalfUp, gomath.RoundHalfDown,
			gomath.RoundDown, gomath.RoundUp, gomath.RoundCeiling, gomath.RoundFloor,
		} {
			operation := decimal.Context{Precision: 1, Rounding: mode, Limits: limits}
			result, err := operation.Apply(context.Background(), value)
			want := int64(0)
			if mode == gomath.RoundUp || mode == gomath.RoundCeiling && sign > 0 || mode == gomath.RoundFloor && sign < 0 {
				want = sign
			}
			conditions := gomath.ConditionSubnormal | gomath.ConditionRounded | gomath.ConditionInexact | gomath.ConditionUnderflow
			if err != nil || result.Value.Coefficient().Int64() != want || result.Value.Exponent() != 0 || result.Conditions != conditions {
				t.Fatalf("underflow sign=%d mode=%s: coefficient=%s exponent=%d conditions=%s error=%v", sign, mode, result.Value.Coefficient(), result.Value.Exponent(), result.Conditions, err)
			}
		}
	}
}

func TestJSONRetainsEscapedDigitBoundaryAndRejectedReceiver(t *testing.T) {
	input := []byte(`"` + strings.Repeat(`\u0039`, gomath.DefaultLimits().MaxInputDigits) + `"`)
	var value decimal.Decimal
	if err := value.UnmarshalJSON(input); err != nil || len(value.String()) != gomath.DefaultLimits().MaxInputDigits {
		t.Fatalf("escaped digit boundary error = %v", err)
	}
	value = decimal.New(42)
	input = []byte(`"` + strings.Repeat("9", 8<<20) + `"`)
	if err := value.UnmarshalJSON(input); !errors.Is(err, gomath.ErrLimitExceeded) || !value.Equal(decimal.New(42)) {
		t.Fatalf("rejected JSON receiver = %s, error = %v", value, err)
	}
}

func TestRoundedCarryRejectsCrossedExponentBudget(t *testing.T) {
	limits := gomath.DefaultLimits()
	limits.MaxExponentMagnitude = 2
	operation := decimal.Context{
		Precision: 1, MinExponent: -2, MaxExponent: 2,
		Rounding: decimal.HalfEven, Limits: limits,
	}
	for _, test := range []struct {
		exponent  int32
		wantLimit bool
	}{
		{0, false},
		{1, true},
	} {
		value, err := decimal.FromBig(big.NewInt(99), test.exponent, limits)
		if err != nil {
			t.Fatal(err)
		}
		result, err := operation.Apply(context.Background(), value)
		if errors.Is(err, gomath.ErrLimitExceeded) != test.wantLimit {
			t.Fatalf("rounded carry exponent %d: error = %v, want limit %t", test.exponent, err, test.wantLimit)
		}
		if !test.wantLimit && (err != nil || result.Value.String() != "100" || result.Conditions != gomath.ConditionRounded|gomath.ConditionInexact) {
			t.Fatalf("rounded carry at boundary = %s [%s], %v", result.Value, result.Conditions, err)
		}
	}
}

func TestQuotientRejectsDerivedExponentBeforeContextClamping(t *testing.T) {
	limits := gomath.DefaultLimits()
	limits.MaxExponentMagnitude = 2
	operation := decimal.Context{
		Precision: 1, MinExponent: -2, MaxExponent: 2,
		Rounding: decimal.HalfEven, Limits: limits,
	}
	for _, exponent := range []int32{-2, 2} {
		numerator, err := decimal.FromBig(big.NewInt(1), exponent, limits)
		if err != nil {
			t.Fatal(err)
		}
		denominator, err := decimal.FromBig(big.NewInt(1), -exponent, limits)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := operation.Quo(context.Background(), numerator, denominator); !errors.Is(err, gomath.ErrLimitExceeded) {
			t.Fatalf("quotient exponent %d: error = %v, want ErrLimitExceeded", 2*exponent, err)
		}
	}
}
