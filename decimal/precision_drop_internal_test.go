package decimal

import (
	"context"
	"errors"
	"math"
	"testing"

	gomath "github.com/faustbrian/go-math"
)

func TestPrecisionDropRejectsWrappedExponent(t *testing.T) {
	if _, err := checkedRoundedExponent(-math.MaxInt32, 1<<32, math.MaxInt32); !errors.Is(err, ErrLimit) {
		t.Fatalf("rounding drop 2^32 error = %v, want ErrLimit", err)
	}
}

func TestPrecisionDropAdmissionBoundaries(t *testing.T) {
	for _, test := range []struct {
		name     string
		exponent int32
		drop     int64
		maximum  int32
		want     int32
		refuse   bool
	}{
		{"exact upper", -math.MaxInt32, 2 * int64(math.MaxInt32), math.MaxInt32, math.MaxInt32, false},
		{"one over upper", -math.MaxInt32, 2*int64(math.MaxInt32) + 1, math.MaxInt32, 0, true},
		{"maximum signed drop", math.MaxInt32, math.MaxInt64, math.MaxInt32, 0, true},
		{"negative drop", 0, math.MinInt64, math.MaxInt32, 0, true},
		{"exact lower", -10, 0, 10, -10, false},
		{"below lower", math.MinInt32, 0, math.MaxInt32, 0, true},
		{"ordinary drop", -2, 3, 10, 1, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := checkedRoundedExponent(test.exponent, test.drop, test.maximum)
			if test.refuse {
				if !errors.Is(err, ErrLimit) {
					t.Fatalf("admission error = %v, want ErrLimit", err)
				}
			} else if err != nil || got != test.want {
				t.Fatalf("admission = %d, %v, want %d, nil", got, err, test.want)
			}
		})
	}
}

func TestPrecisionDropPublicApply(t *testing.T) {
	operation := Context{Precision: 3, MinExponent: -20, MaxExponent: 20, Rounding: HalfEven}
	for _, test := range []struct {
		input, want string
		exponent    int32
		conditions  gomath.Condition
	}{
		{"9016", "9020", 1, gomath.ConditionRounded | gomath.ConditionInexact},
		{"9999", "10000", 2, gomath.ConditionRounded | gomath.ConditionInexact},
		{"1230", "1230", 1, gomath.ConditionRounded},
		{"-1255", "-1260", 1, gomath.ConditionRounded | gomath.ConditionInexact},
	} {
		t.Run(test.input, func(t *testing.T) {
			input := MustParse(test.input)
			result, err := operation.Apply(context.Background(), input)
			if err != nil || result.Value.String() != test.want || result.Value.Exponent() != test.exponent || result.Conditions != test.conditions {
				t.Fatalf("Apply() = %s exp %d [%s], %v", result.Value, result.Value.Exponent(), result.Conditions, err)
			}
			result.Value.Coefficient().SetInt64(7)
			if input.String() != test.input || result.Value.String() != test.want {
				t.Fatal("rounding or coefficient access changed owned values")
			}
		})
	}
	limits := gomath.DefaultLimits()
	limits.MaxExponentMagnitude = 1
	bounded := Context{Precision: 1, MinExponent: -1, MaxExponent: 1, Rounding: HalfEven, Limits: limits}
	if _, err := bounded.Apply(context.Background(), New(9999)); !errors.Is(err, ErrLimit) {
		t.Fatalf("existing exponent limit error = %v, want ErrLimit", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := operation.Apply(ctx, New(9016)); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled Apply error = %v", err)
	}
	operation.Traps = gomath.ConditionInexact
	result, err := operation.Apply(context.Background(), New(9016))
	if !errors.Is(err, gomath.ErrTrappedCondition) || result.Value.String() != "9020" || result.Conditions != gomath.ConditionRounded|gomath.ConditionInexact {
		t.Fatalf("trapped Apply = %s [%s], %v", result.Value, result.Conditions, err)
	}
}

func TestPrecisionDropLargeSupportedPrecision(t *testing.T) {
	maximumInt := int(^uint(0) >> 1)
	if uint64(maximumInt) < math.MaxUint32 {
		t.Skip("active host bit budget cannot admit MaxUint32 precision")
	}
	limits := gomath.DefaultLimits()
	limits.MaxPrecision = math.MaxUint32
	limits.MaxIntermediateBits = maximumInt
	operation := Context{Precision: math.MaxUint32, MinExponent: -20, MaxExponent: 20, Rounding: HalfEven, Limits: limits}
	result, err := operation.Apply(context.Background(), New(12))
	if err != nil || result.Value.String() != "12" || result.Conditions != 0 {
		t.Fatalf("large-precision Apply = %s [%s], %v", result.Value, result.Conditions, err)
	}
}
