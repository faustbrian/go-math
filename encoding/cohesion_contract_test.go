package encoding_test

import (
	"errors"
	"testing"

	gomath "github.com/faustbrian/go-math"
	"github.com/faustbrian/go-math/bigfloat"
	mathencoding "github.com/faustbrian/go-math/encoding"
)

type decodeCase struct {
	name string
	kind byte
	run  func([]byte, gomath.Limits) error
}

func TestBinaryDecodersClassifyShortEnvelopesAsSyntax(t *testing.T) {
	t.Parallel()

	limits := gomath.DefaultLimits()
	for _, decoder := range binaryDecoders() {
		decoder := decoder
		t.Run(decoder.name, func(t *testing.T) {
			t.Parallel()
			for length := 0; length < 4; length++ {
				err := decoder.run(make([]byte, length), limits)
				assertBinaryInputError(t, err, gomath.ErrInvalidSyntax, "math: invalid syntax")
			}
		})
	}
}

func TestBinaryDecodersHideNestedParserDiagnostics(t *testing.T) {
	t.Parallel()

	limits := gomath.DefaultLimits()
	for _, decoder := range binaryDecoders() {
		decoder := decoder
		t.Run(decoder.name, func(t *testing.T) {
			t.Parallel()
			data := []byte{'G', 'M', 1, decoder.kind}
			err := decoder.run(data, limits)
			assertBinaryInputError(t, err, gomath.ErrInvalidSyntax, "math: invalid syntax")
		})
	}
}

func TestBinaryDecodersApplyMaximumBeforeHeaderGrammar(t *testing.T) {
	t.Parallel()

	limits := gomath.DefaultLimits()
	limits.MaxIntermediateBits = 8
	for _, decoder := range binaryDecoders() {
		decoder := decoder
		t.Run(decoder.name, func(t *testing.T) {
			t.Parallel()
			atMaximum := make([]byte, 65)
			assertBinaryInputError(
				t, decoder.run(atMaximum, limits),
				gomath.ErrInvalidSyntax, "math: invalid syntax",
			)
			overMaximum := make([]byte, 66)
			assertBinaryInputError(
				t, decoder.run(overMaximum, limits),
				gomath.ErrLimitExceeded,
				"math: resource limit exceeded: binary payload size",
			)

			invalid := limits
			invalid.MaxInputDigits = 0
			for _, data := range [][]byte{nil, overMaximum} {
				err := decoder.run(data, invalid)
				if !errors.Is(err, gomath.ErrInvalidArgument) ||
					errors.Is(err, gomath.ErrInvalidSyntax) ||
					errors.Is(err, gomath.ErrLimitExceeded) {
					t.Fatalf("invalid limits precedence error = %v", err)
				}
			}
		})
	}
}

func TestUnmarshalFloatClassifiesInvalidRoundingAsSyntax(t *testing.T) {
	t.Parallel()

	limits := gomath.DefaultLimits()
	result, err := bigfloat.NewInt64(1, bigfloat.Context{
		Precision: 8,
		Rounding:  gomath.RoundHalfEven,
		Limits:    limits,
	})
	if err != nil {
		t.Fatal(err)
	}
	data, err := mathencoding.MarshalFloat(result.Value)
	if err != nil {
		t.Fatal(err)
	}
	for _, rounding := range []byte{byte(gomath.RoundHalfDown), 0xff} {
		malformed := append([]byte(nil), data...)
		malformed[4] = rounding
		_, err = mathencoding.UnmarshalFloat(malformed, limits)
		assertBinaryInputError(t, err, gomath.ErrInvalidSyntax, "math: invalid syntax")
	}
}

func TestFloatRoundTripPreservesEverySupportedRoundingMode(t *testing.T) {
	t.Parallel()

	limits := gomath.DefaultLimits()
	for _, rounding := range []gomath.RoundingMode{
		gomath.RoundHalfEven,
		gomath.RoundHalfUp,
		gomath.RoundDown,
		gomath.RoundUp,
		gomath.RoundCeiling,
		gomath.RoundFloor,
	} {
		result, err := bigfloat.NewInt64(1, bigfloat.Context{
			Precision: 8,
			Rounding:  rounding,
			Limits:    limits,
		})
		if err != nil {
			t.Fatalf("NewInt64(%s) error = %v", rounding, err)
		}
		data, err := mathencoding.MarshalFloat(result.Value)
		if err != nil {
			t.Fatalf("MarshalFloat(%s) error = %v", rounding, err)
		}
		decoded, err := mathencoding.UnmarshalFloat(data, limits)
		if err != nil {
			t.Fatalf("UnmarshalFloat(%s) error = %v", rounding, err)
		}
		if decoded.Rounding() != rounding || decoded.String() != result.Value.String() {
			t.Fatalf(
				"round trip %s = value %s, rounding %s",
				rounding, decoded.String(), decoded.Rounding(),
			)
		}
	}
}

func TestUnmarshalFloatPreservesDecodedValueLimit(t *testing.T) {
	t.Parallel()

	limits := gomath.DefaultLimits()
	result, err := bigfloat.NewInt64(1, bigfloat.Context{
		Precision: 8,
		Rounding:  gomath.RoundHalfEven,
		Limits:    limits,
	})
	if err != nil {
		t.Fatal(err)
	}
	data, err := mathencoding.MarshalFloat(result.Value)
	if err != nil {
		t.Fatal(err)
	}
	limits.MaxIntermediateBits = 1
	_, err = mathencoding.UnmarshalFloat(data, limits)
	assertBinaryInputError(
		t, err, gomath.ErrLimitExceeded,
		"math: resource limit exceeded: binary precision",
	)
}

func binaryDecoders() []decodeCase {
	return []decodeCase{
		{name: "integer", kind: 1, run: func(data []byte, limits gomath.Limits) error {
			_, err := mathencoding.UnmarshalInteger(data, limits)
			return err
		}},
		{name: "rational", kind: 2, run: func(data []byte, limits gomath.Limits) error {
			_, err := mathencoding.UnmarshalRational(data, limits)
			return err
		}},
		{name: "decimal", kind: 3, run: func(data []byte, limits gomath.Limits) error {
			_, err := mathencoding.UnmarshalDecimal(data, limits)
			return err
		}},
		{name: "float", kind: 4, run: func(data []byte, limits gomath.Limits) error {
			_, err := mathencoding.UnmarshalFloat(data, limits)
			return err
		}},
	}
}

func assertBinaryInputError(t *testing.T, err, want error, text string) {
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
