package gomath

import "fmt"

// Limits bounds admission and selected arithmetic/formatting work. Every field
// must be positive. It is copied by value and contains no mutable shared state.
//
// Each family applies the relevant fields; typed constructors and unbounded
// accessors such as Integer.String do not apply every parser/formatting limit.
// Start with DefaultLimits and adjust for the caller's workload. Very large
// trusted limits can still authorize expensive synchronous allocation and
// arithmetic.
type Limits struct {
	// MaxInputDigits bounds parser digits (per rational component); bigfloat.Parse
	// uses it as a byte ceiling for the entire input, including punctuation.
	MaxInputDigits int
	// MaxOutputDigits bounds decimal and binary-float rendered digits.
	// Rational.Decimal instead admits text length up to this value plus two bytes.
	// Integer.String and Rational.String do not apply this bound.
	MaxOutputDigits int
	// MaxExponentMagnitude bounds the absolute stored decimal exponent and binary
	// MantExp exponent, as well as admitted decimal scaling work.
	MaxExponentMagnitude int32
	// MaxPrecision bounds significant decimal digits or binary precision bits.
	MaxPrecision uint32
	// MaxPowerExponent bounds the magnitude of integer/rational power exponents.
	MaxPowerExponent uint64
	// MaxRootDegree bounds the integer root degree.
	MaxRootDegree uint32
	// MaxRandomBits bounds the bit length of the random interval size max-min.
	MaxRandomBits int
	// MaxRandomAttempts bounds sampling attempts before ErrLimitExceeded.
	MaxRandomAttempts int
	// MaxIntermediateBits bounds admitted operand/result bit lengths and selected
	// explicit intermediate work. It is not a byte-allocation quota.
	MaxIntermediateBits int
	// MaxDecimalExpansion bounds the nonnegative fractional scale in
	// Rational.Decimal.
	MaxDecimalExpansion int
	// MaxDiagnosticBytes is retained for source compatibility and must remain
	// positive for Validate, but does not configure diagnostic length.
	//
	// Deprecated: reserved no-op; package errors use fixed, input-free categories.
	// Keep the DefaultLimits value rather than using this field to truncate
	// errors.
	MaxDiagnosticBytes int
}

// DefaultLimits returns a fresh value of general-purpose defensive bounds. The
// bounds are package policy, not precision or conformance requirements.
func DefaultLimits() Limits {
	return Limits{
		MaxInputDigits:       100_000,
		MaxOutputDigits:      100_000,
		MaxExponentMagnitude: 1_000_000,
		MaxPrecision:         100_000,
		MaxPowerExponent:     1_000_000,
		MaxRootDegree:        1_000_000,
		MaxRandomBits:        1_000_000,
		MaxRandomAttempts:    128,
		MaxIntermediateBits:  4_000_000,
		MaxDecimalExpansion:  100_000,
		MaxDiagnosticBytes:   256,
	}
}

// Validate reports whether every limit is positive.
func (l Limits) Validate() error {
	if l.MaxInputDigits <= 0 {
		return fmt.Errorf("%w: MaxInputDigits must be positive", ErrInvalidArgument)
	}
	if l.MaxOutputDigits <= 0 {
		return fmt.Errorf("%w: MaxOutputDigits must be positive", ErrInvalidArgument)
	}
	if l.MaxExponentMagnitude <= 0 {
		return fmt.Errorf("%w: MaxExponentMagnitude must be positive", ErrInvalidArgument)
	}
	if l.MaxPrecision == 0 {
		return fmt.Errorf("%w: MaxPrecision must be positive", ErrInvalidArgument)
	}
	if l.MaxPowerExponent == 0 {
		return fmt.Errorf("%w: MaxPowerExponent must be positive", ErrInvalidArgument)
	}
	if l.MaxRootDegree == 0 {
		return fmt.Errorf("%w: MaxRootDegree must be positive", ErrInvalidArgument)
	}
	if l.MaxRandomBits <= 0 {
		return fmt.Errorf("%w: MaxRandomBits must be positive", ErrInvalidArgument)
	}
	if l.MaxRandomAttempts <= 0 {
		return fmt.Errorf("%w: MaxRandomAttempts must be positive", ErrInvalidArgument)
	}
	if l.MaxIntermediateBits <= 0 {
		return fmt.Errorf("%w: MaxIntermediateBits must be positive", ErrInvalidArgument)
	}
	if l.MaxDecimalExpansion <= 0 {
		return fmt.Errorf("%w: MaxDecimalExpansion must be positive", ErrInvalidArgument)
	}
	if l.MaxDiagnosticBytes <= 0 {
		return fmt.Errorf("%w: MaxDiagnosticBytes must be positive", ErrInvalidArgument)
	}

	return nil
}
