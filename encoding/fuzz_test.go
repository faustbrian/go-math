package encoding_test

import (
	"bytes"
	"errors"
	"testing"

	gomath "github.com/faustbrian/go-math"
	mathencoding "github.com/faustbrian/go-math/encoding"
)

func FuzzBinaryDecoders(f *testing.F) {
	f.Add([]byte("GM\x01\x01\x00\x00"))
	f.Add([]byte{})
	f.Fuzz(func(t *testing.T, input []byte) {
		limits := gomath.DefaultLimits()
		limits.MaxIntermediateBits = 4096
		if value, err := mathencoding.UnmarshalInteger(input, limits); err == nil {
			encoded, encodeErr := mathencoding.MarshalInteger(value)
			assertCanonical(t, input, encoded, encodeErr)
		} else {
			assertFuzzBinaryError(t, err)
		}
		if value, err := mathencoding.UnmarshalRational(input, limits); err == nil {
			encoded, encodeErr := mathencoding.MarshalRational(value)
			assertCanonical(t, input, encoded, encodeErr)
		} else {
			assertFuzzBinaryError(t, err)
		}
		if value, err := mathencoding.UnmarshalDecimal(input, limits); err == nil {
			encoded, encodeErr := mathencoding.MarshalDecimal(value)
			assertCanonical(t, input, encoded, encodeErr)
		} else {
			assertFuzzBinaryError(t, err)
		}
		if value, err := mathencoding.UnmarshalFloat(input, limits); err == nil {
			encoded, encodeErr := mathencoding.MarshalFloat(value)
			assertCanonical(t, input, encoded, encodeErr)
		} else {
			assertFuzzBinaryError(t, err)
		}
	})
}

func assertFuzzBinaryError(t *testing.T, err error) {
	t.Helper()
	if len(err.Error()) > 128 {
		t.Fatalf("unbounded decoder diagnostic of %d bytes", len(err.Error()))
	}
	matches := 0
	for _, category := range []error{
		gomath.ErrInvalidArgument, gomath.ErrInvalidSyntax, gomath.ErrLimitExceeded,
	} {
		if errors.Is(err, category) {
			matches++
		}
	}
	if matches != 1 {
		t.Fatalf("decoder error matches %d terminal categories: %v", matches, err)
	}
}

func assertCanonical(t *testing.T, input []byte, encoded []byte, err error) {
	t.Helper()
	if err != nil || !bytes.Equal(input, encoded) {
		t.Fatalf("accepted non-canonical input %x as %x: %v", input, encoded, err)
	}
}
