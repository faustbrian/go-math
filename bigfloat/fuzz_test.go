package bigfloat_test

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	gomath "github.com/faustbrian/go-math"
	"github.com/faustbrian/go-math/bigfloat"
)

func FuzzParse(f *testing.F) {
	for _, seed := range []string{
		"0", "-1.25", "0x1.8p+2", "+Inf",
		strings.Repeat("9", 256), strings.Repeat("9", 257),
	} {
		f.Add(seed, uint8(64))
	}
	f.Fuzz(func(t *testing.T, input string, precisionByte uint8) {
		precision := uint(precisionByte%64 + 1)
		bases := [...]int{0, 2, 10, 16}
		base := bases[precisionByte>>6]
		operation := bigfloat.Context{
			Precision: precision, Rounding: gomath.RoundHalfEven,
			Limits: gomath.DefaultLimits(),
		}
		operation.Limits.MaxInputDigits = 256
		result, err := bigfloat.Parse(input, base, operation)
		if err != nil {
			assertFuzzBigFloatError(t, err)
			return
		}
		if result.Value.Precision() != precision {
			t.Fatalf("precision = %d, want %d", result.Value.Precision(), precision)
		}
		roundTrip, err := bigfloat.Parse(result.Value.String(), 10, operation)
		if err != nil || !roundTrip.Value.Equal(result.Value) || roundTrip.Value.Signbit() != result.Value.Signbit() {
			t.Fatalf("round trip changed %s: %v", result.Value, err)
		}
		data, err := json.Marshal(result.Value)
		if err != nil || len(data) < 2 || data[0] != '"' {
			t.Fatalf("unsafe JSON encoding %q: %v", data, err)
		}
	})
}

func assertFuzzBigFloatError(t *testing.T, err error) {
	t.Helper()
	if len(err.Error()) > 128 {
		t.Fatalf("unbounded parser diagnostic of %d bytes", len(err.Error()))
	}
	matches := 0
	for _, category := range []error{gomath.ErrInvalidSyntax, gomath.ErrLimitExceeded} {
		if errors.Is(err, category) {
			matches++
		}
	}
	if matches != 1 {
		t.Fatalf("parser error matches %d terminal categories: %v", matches, err)
	}
}
