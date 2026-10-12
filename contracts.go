// Package gomath defines resource limits, rounding modes, conditions and errors
// shared by the distinct integer, rational, decimal and bigfloat packages.
//
// Choose a numeric family explicitly: integer and rational arithmetic is exact,
// decimal supports exact and rounded finite base-10 results, and bigfloat
// rounds to a binary precision. This package does not own currency, units or
// presentation policy. Inspect errors with errors.Is; conditions and traps are
// family-specific. Limits are caller policy, not a guarantee that arbitrary
// trusted settings are affordable or that synchronous math/big work is
// preemptible.
package gomath

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidArgument reports an invalid option or operation argument.
	ErrInvalidArgument = errors.New("math: invalid argument")
	// ErrInvalidSyntax reports input outside a type's strict grammar.
	ErrInvalidSyntax = errors.New("math: invalid syntax")
	// ErrLimitExceeded reports work rejected by an explicit resource bound.
	ErrLimitExceeded = errors.New("math: resource limit exceeded")
	// ErrDivisionByZero reports division by numeric zero.
	ErrDivisionByZero = errors.New("math: division by zero")
	// ErrDomain reports an operation outside its mathematical domain.
	ErrDomain = errors.New("math: domain error")
	// ErrConversion reports a conversion that cannot be exact.
	ErrConversion = errors.New("math: inexact conversion")
	// ErrOverflow is the shared overflow error category. Current decimal contexts
	// instead report ConditionOverflow with a finite value or a trap/limit error.
	ErrOverflow = errors.New("math: overflow")
	// ErrUnderflow is the shared underflow error category. Current decimal contexts
	// instead report ConditionUnderflow with a finite value or a trap/limit error.
	ErrUnderflow = errors.New("math: underflow")
	// ErrTrappedCondition reports a decimal or floating-point condition selected
	// by an operation context's trap mask.
	ErrTrappedCondition = errors.New("math: trapped condition")
)

// ErrRandomSource reports a failure from a caller-provided random source.
var ErrRandomSource = errors.New("math: random source failed")

// Condition is a bit set of arithmetic conditions. Raised flags and trap
// behavior are family-specific; a flag does not by itself mean an error was
// returned.
type Condition uint16

const (
	// ConditionRounded records discarded digits, including discarded zeros.
	ConditionRounded Condition = 1 << iota
	// ConditionInexact records a result that differs from the exact value.
	ConditionInexact
	// ConditionOverflow records a decimal adjusted exponent above its context
	// range.
	ConditionOverflow
	// ConditionUnderflow records an inexact subnormal decimal result.
	ConditionUnderflow
	// ConditionDivisionByZero records division by numeric zero.
	ConditionDivisionByZero
	// ConditionInvalidOperation records an undefined binary-float operation.
	ConditionInvalidOperation
	// ConditionClamped is the shared clamping flag; current numeric operations do
	// not raise it.
	ConditionClamped
	// ConditionSubnormal records a decimal adjusted exponent below the normal
	// range.
	ConditionSubnormal
)

// Has reports whether every bit in wanted is present. Has(0) is true.
func (c Condition) Has(wanted Condition) bool { return c&wanted == wanted }

// String returns stable comma-separated condition names, or "none" for zero.
// Unknown bits are retained as a hexadecimal suffix, not silently discarded.
func (c Condition) String() string {
	if c == 0 {
		return "none"
	}

	names := []struct {
		condition Condition
		name      string
	}{
		{ConditionRounded, "rounded"},
		{ConditionInexact, "inexact"},
		{ConditionOverflow, "overflow"},
		{ConditionUnderflow, "underflow"},
		{ConditionDivisionByZero, "division_by_zero"},
		{ConditionInvalidOperation, "invalid_operation"},
		{ConditionClamped, "clamped"},
		{ConditionSubnormal, "subnormal"},
	}
	result := make([]string, 0, len(names)+1)
	known := Condition(0)
	for _, candidate := range names {
		known |= candidate.condition
		if c.Has(candidate.condition) {
			result = append(result, candidate.name)
		}
	}
	if unknown := c &^ known; unknown != 0 {
		result = append(result, fmt.Sprintf("unknown(0x%x)", uint16(unknown)))
	}

	return strings.Join(result, ",")
}

// RoundingMode selects how a discarded nonzero remainder affects a result.
type RoundingMode uint8

const (
	// RoundHalfEven selects the nearest result, breaking ties toward an even
	// retained digit.
	RoundHalfEven RoundingMode = iota
	// RoundHalfUp selects the nearest result, breaking ties away from zero.
	RoundHalfUp
	// RoundHalfDown selects the nearest result, breaking ties toward zero.
	RoundHalfDown
	// RoundDown rounds toward zero.
	RoundDown
	// RoundUp rounds away from zero.
	RoundUp
	// RoundCeiling rounds toward positive infinity.
	RoundCeiling
	// RoundFloor rounds toward negative infinity.
	RoundFloor
)

// Valid reports whether r is one of the seven shared rounding modes.
// Binary-float contexts additionally reject RoundHalfDown.
func (r RoundingMode) Valid() bool { return r <= RoundFloor }

// String returns the stable configuration name of r.
func (r RoundingMode) String() string {
	names := [...]string{
		"half_even", "half_up", "half_down", "down", "up", "ceiling", "floor",
	}
	if !r.Valid() {
		return fmt.Sprintf("RoundingMode(%d)", r)
	}

	return names[r]
}
