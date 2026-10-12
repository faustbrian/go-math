// Package bigfloat provides immutable arbitrary-precision binary floating-point
// values for inexact algorithms, not money or exact quantities. Context
// precision is measured in bits; conversion and arithmetic may round and report
// accuracy.
//
// Mutable math/big inputs and returned accessors are copied. Value operations
// can share operands concurrently. Infinity and signed zero are supported, but
// NaN-producing infinity arithmetic can panic; see Context and its methods.
// Text does not encode precision or rounding, and directed-rounding parses are
// not guaranteed to recover the rendered value. Use exact Rat or binary
// encoding when those distinctions matter.
package bigfloat

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	gomath "github.com/faustbrian/go-math"
)

// Context makes binary precision, rounding, traps and resource limits explicit.
// The zero Context is invalid. Context-taking arithmetic checks cancellation
// before admission, but the synchronous math/big primitive is not preemptible.
//
// Infinities are admitted, but NaN is not representable: Add of opposite
// infinities, Sub of same-signed infinities, Mul of zero and infinity, and Quo
// of two infinities panic with big.ErrNaN. Validate those operand combinations
// before calling.
type Context struct {
	// Precision counts significant binary bits; it must be positive and within
	// MaxPrecision and MaxIntermediateBits.
	Precision uint
	// Rounding selects a shared mode except RoundHalfDown, which is unsupported.
	// Zero selects nearest with ties to even.
	Rounding gomath.RoundingMode
	// Traps selects raised rounding or division-by-zero conditions that return
	// ErrTrappedCondition with the result intact. Domain errors are returned
	// directly.
	Traps gomath.Condition
	// Limits uses DefaultLimits only when the entire value is zero; otherwise
	// every field must pass Validate.
	Limits gomath.Limits
}

// Result contains the rounded Float, its accuracy relative to the exact result,
// and raised conditions. ErrTrappedCondition retains the result; domain
// failures may retain InvalidOperation without a usable computed value.
type Result struct {
	// Value is the rounded result, also retained when a condition is trapped.
	Value Float
	// Accuracy reports whether Value is below, equal to or above the exact result.
	Accuracy big.Accuracy
	// Conditions raises Rounded and Inexact for non-exact accuracy, plus the
	// explicit DivisionByZero or InvalidOperation branches.
	Conditions gomath.Condition
}

// Float is an immutable arbitrary-precision binary floating-point value. The
// zero value is numeric zero with no operation precision; contexts remain
// mandatory at construction and operation boundaries.
type Float struct {
	value     big.Float
	precision uint
	rounding  gomath.RoundingMode
}

// NewInt64 converts an integer at the context's binary precision. Low precision
// may round even an integer source; Accuracy and Conditions report that
// rounding.
func NewInt64(value int64, operation Context) (Result, error) {
	mode, limits, err := operation.validate()
	if err != nil {
		return Result{}, err
	}
	result := new(big.Float).SetPrec(operation.Precision).SetMode(mode).SetInt64(value)

	return operation.finish(result, result.Acc(), limits)
}

// FromBig admits a non-nil source's precision and exponent under Limits, then
// copies and rounds it to the explicit context. Rounding to a lower precision
// does not bypass source admission. Infinities and signed zeros are retained.
func FromBig(value *big.Float, operation Context) (Result, error) {
	if value == nil {
		return Result{}, fmt.Errorf("%w: nil big.Float", gomath.ErrInvalidArgument)
	}
	mode, limits, err := operation.validate()
	if err != nil {
		return Result{}, err
	}
	if err := checkFloatInputLimits(value, limits); err != nil {
		return Result{}, err
	}
	result := new(big.Float).SetPrec(operation.Precision).SetMode(mode).Set(value)

	return operation.finish(result, result.Acc(), limits)
}

// FromRat converts an exact non-nil rational at context binary precision and
// reports rounding. Both source components must fit MaxIntermediateBits.
func FromRat(value *big.Rat, operation Context) (Result, error) {
	if value == nil {
		return Result{}, fmt.Errorf("%w: nil big.Rat", gomath.ErrInvalidArgument)
	}
	mode, limits, err := operation.validate()
	if err != nil {
		return Result{}, err
	}
	if value.Num().BitLen() > limits.MaxIntermediateBits || value.Denom().BitLen() > limits.MaxIntermediateBits {
		return Result{}, fmt.Errorf("%w: rational source bits", gomath.ErrLimitExceeded)
	}
	result := new(big.Float).SetPrec(operation.Precision).SetMode(mode).SetRat(value)

	return operation.finish(result, result.Acc(), limits)
}

// Parse uses math/big.ParseFloat grammar with base 0, 2, 10 or 16, rounding to
// the explicit binary precision. Base zero infers the base; surrounding
// whitespace is rejected. MaxInputDigits here bounds all input bytes, not only
// digits. Finite values and signed infinities are accepted, but NaN is not.
func Parse(text string, base int, operation Context) (Result, error) {
	mode, limits, err := operation.validate()
	if err != nil {
		return Result{}, err
	}
	if base != 0 && base != 2 && base != 10 && base != 16 {
		return Result{}, fmt.Errorf("%w: float base", gomath.ErrInvalidArgument)
	}
	if text == "" {
		return Result{}, gomath.ErrInvalidSyntax
	}
	if len(text) > limits.MaxInputDigits {
		return Result{}, fmt.Errorf("%w: binary float input", gomath.ErrLimitExceeded)
	}
	if strings.TrimSpace(text) != text {
		return Result{}, gomath.ErrInvalidSyntax
	}
	value, _, err := big.ParseFloat(text, base, operation.Precision, mode)
	if err != nil {
		return Result{}, gomath.ErrInvalidSyntax
	}

	return operation.finish(value, value.Acc(), limits)
}

// Big returns a mutable copy of f.
func (f Float) Big() *big.Float { return new(big.Float).Copy(&f.value) }

// Precision returns the construction precision in bits.
func (f Float) Precision() uint { return f.precision }

// Rounding returns the construction rounding mode.
func (f Float) Rounding() gomath.RoundingMode { return f.rounding }

// String renders decimal text with significant digits derived from f's binary
// precision. It retains infinity and the sign of zero, but does not encode
// precision or rounding policy. Re-parsing with directed rounding may change
// the value; use Rat for the exact finite value or the binary codec for stored
// policy.
func (f Float) String() string { return f.value.Text('g', roundTripDecimalDigits(f.precision)) }

// IsInf reports whether f is positive or negative infinity.
func (f Float) IsInf() bool { return f.value.IsInf() }

// Sign returns -1, 0, or +1.
func (f Float) Sign() int { return f.value.Sign() }

// Signbit reports whether f has a negative sign, including negative zero.
func (f Float) Signbit() bool { return f.value.Signbit() }

// Cmp compares f and other numerically.
func (f Float) Cmp(other Float) int { return f.value.Cmp(&other.value) }

// Equal reports numeric equality, ignoring precision, rounding policy and the
// distinction between positive and negative zero.
func (f Float) Equal(other Float) bool { return f.Cmp(other) == 0 }

// Neg returns -f while preserving representation policy.
func (f Float) Neg() Float {
	value := new(big.Float).Copy(&f.value)
	value.Neg(value)

	return fromBig(value, f.precision, f.rounding)
}

// Abs returns |f| while preserving representation policy.
func (f Float) Abs() Float {
	value := new(big.Float).Copy(&f.value)
	value.Abs(value)

	return fromBig(value, f.precision, f.rounding)
}

// Int returns an independent mutable integer only if f is finite and integral;
// otherwise it returns ErrConversion. It does not apply a fresh resource
// budget.
func (f Float) Int() (*big.Int, error) {
	if f.IsInf() {
		return nil, gomath.ErrConversion
	}
	value, accuracy := f.value.Int(nil)
	if accuracy != big.Exact {
		return nil, gomath.ErrConversion
	}

	return new(big.Int).Set(value), nil
}

// Rat returns the exact finite binary value as an independent mutable rational.
// Infinity returns ErrConversion; no fresh resource budget is applied.
func (f Float) Rat() (*big.Rat, error) {
	if f.IsInf() {
		return nil, gomath.ErrConversion
	}
	value, _ := f.value.Rat(nil)

	return new(big.Rat).Set(value), nil
}

// MarshalText returns String's decimal form in a new buffer; binary precision
// and rounding policy are not encoded.
func (f Float) MarshalText() ([]byte, error) { return []byte(f.String()), nil }

// MarshalJSON encodes String as a JSON string, avoiding JSON-number narrowing
// but not encoding binary precision or rounding policy.
func (f Float) MarshalJSON() ([]byte, error) { return json.Marshal(f.String()) }

// Add performs context-rounded addition. Opposite infinities panic with
// big.ErrNaN rather than returning an InvalidOperation result.
func (c Context) Add(ctx context.Context, left, right Float) (Result, error) {
	return c.binary(ctx, left, right, (*big.Float).Add)
}

// Sub performs context-rounded subtraction. Same-signed infinities panic with
// big.ErrNaN rather than returning an InvalidOperation result.
func (c Context) Sub(ctx context.Context, left, right Float) (Result, error) {
	return c.binary(ctx, left, right, (*big.Float).Sub)
}

// Mul performs context-rounded multiplication. Zero times infinity panics with
// big.ErrNaN rather than returning an InvalidOperation result.
func (c Context) Mul(ctx context.Context, left, right Float) (Result, error) {
	return c.binary(ctx, left, right, (*big.Float).Mul)
}

// Quo performs context-rounded division. A nonzero numerator over signed zero
// returns the corresponding infinity and DivisionByZero, or ErrTrappedCondition
// with that result when trapped. Zero over zero returns ErrDomain and
// InvalidOperation independently of Traps. Two infinities panic with
// big.ErrNaN.
func (c Context) Quo(ctx context.Context, numerator, denominator Float) (Result, error) {
	mode, limits, err := c.validateContext(ctx)
	if err != nil {
		return Result{}, err
	}
	if err := checkFloatInputLimits(&numerator.value, limits); err != nil {
		return Result{}, err
	}
	if err := checkFloatInputLimits(&denominator.value, limits); err != nil {
		return Result{}, err
	}
	if denominator.Sign() == 0 {
		if numerator.Sign() == 0 {
			result := Result{Conditions: gomath.ConditionInvalidOperation}

			return result, fmt.Errorf("%w: zero divided by zero", gomath.ErrDomain)
		}
		value := new(big.Float).SetPrec(c.Precision).SetMode(mode)
		value.SetInf(numerator.value.Signbit() != denominator.value.Signbit())
		result := Result{
			Value:      fromBig(value, c.Precision, c.Rounding),
			Accuracy:   big.Exact,
			Conditions: gomath.ConditionDivisionByZero,
		}
		if c.Traps.Has(gomath.ConditionDivisionByZero) {
			return result, fmt.Errorf("%w: division_by_zero", gomath.ErrTrappedCondition)
		}

		return result, nil
	}

	return c.binaryWithMode(numerator, denominator, mode, limits, (*big.Float).Quo)
}

// Sqrt performs a context-rounded square root. Negative values, including
// negative infinity, return ErrDomain with InvalidOperation independently of
// Traps. Negative zero retains its sign.
func (c Context) Sqrt(ctx context.Context, value Float) (Result, error) {
	mode, limits, err := c.validateContext(ctx)
	if err != nil {
		return Result{}, err
	}
	if err := checkFloatInputLimits(&value.value, limits); err != nil {
		return Result{}, err
	}
	if value.Sign() < 0 {
		result := Result{Conditions: gomath.ConditionInvalidOperation}

		return result, fmt.Errorf("%w: square root of negative float", gomath.ErrDomain)
	}
	result := new(big.Float).SetPrec(c.Precision).SetMode(mode).Sqrt(&value.value)

	return c.finish(result, result.Acc(), limits)
}

func (c Context) binary(ctx context.Context, left, right Float, operation func(*big.Float, *big.Float, *big.Float) *big.Float) (Result, error) {
	mode, limits, err := c.validateContext(ctx)
	if err != nil {
		return Result{}, err
	}

	return c.binaryWithMode(left, right, mode, limits, operation)
}

func (c Context) binaryWithMode(left, right Float, mode big.RoundingMode, limits gomath.Limits, operation func(*big.Float, *big.Float, *big.Float) *big.Float) (Result, error) {
	if err := checkFloatInputLimits(&left.value, limits); err != nil {
		return Result{}, err
	}
	if err := checkFloatInputLimits(&right.value, limits); err != nil {
		return Result{}, err
	}
	result := new(big.Float).SetPrec(c.Precision).SetMode(mode)
	operation(result, &left.value, &right.value)

	return c.finish(result, result.Acc(), limits)
}

func (c Context) finish(value *big.Float, accuracy big.Accuracy, limits gomath.Limits) (Result, error) {
	if err := checkValueLimits(value, limits); err != nil {
		return Result{}, err
	}
	conditions := gomath.Condition(0)
	if accuracy != big.Exact {
		conditions = gomath.ConditionRounded | gomath.ConditionInexact
	}
	result := Result{
		Value:      fromBig(value, c.Precision, c.Rounding),
		Accuracy:   accuracy,
		Conditions: conditions,
	}
	if trapped := conditions & c.Traps; trapped != 0 {
		return result, fmt.Errorf("%w: %s", gomath.ErrTrappedCondition, trapped)
	}

	return result, nil
}

func (c Context) validate() (big.RoundingMode, gomath.Limits, error) {
	limits := c.Limits
	if limits == (gomath.Limits{}) {
		limits = gomath.DefaultLimits()
	}
	if err := limits.Validate(); err != nil {
		return 0, gomath.Limits{}, err
	}
	if c.Precision == 0 || c.Precision > uint(limits.MaxPrecision) {
		return 0, gomath.Limits{}, fmt.Errorf("%w: binary precision", gomath.ErrInvalidArgument)
	}
	if c.Precision > uint(limits.MaxIntermediateBits) {
		return 0, gomath.Limits{}, fmt.Errorf("%w: binary precision", gomath.ErrLimitExceeded)
	}
	mode, ok := roundingMode(c.Rounding)
	if !ok {
		return 0, gomath.Limits{}, fmt.Errorf("%w: binary rounding mode", gomath.ErrInvalidArgument)
	}

	return mode, limits, nil
}

func checkValueLimits(value *big.Float, limits gomath.Limits) error {
	if err := checkFloatInputLimits(value, limits); err != nil {
		return err
	}
	if value.IsInf() {
		return nil
	}
	exponent := value.MantExp(nil)
	if maximumRenderedDigits(value.Prec(), exponent) <= limits.MaxOutputDigits {
		return nil
	}
	if renderedDigits(value.Text('g', roundTripDecimalDigits(value.Prec()))) > limits.MaxOutputDigits {
		return fmt.Errorf("%w: binary output digits", gomath.ErrLimitExceeded)
	}

	return nil
}

func checkFloatInputLimits(value *big.Float, limits gomath.Limits) error {
	if value.Prec() > uint(limits.MaxIntermediateBits) {
		return fmt.Errorf("%w: binary value precision", gomath.ErrLimitExceeded)
	}
	if value.IsInf() {
		return nil
	}
	exponent := value.MantExp(nil)
	maximumExponent := int(limits.MaxExponentMagnitude)
	if exponent < -maximumExponent || exponent > maximumExponent {
		return fmt.Errorf("%w: binary exponent", gomath.ErrLimitExceeded)
	}

	return nil
}

func maximumRenderedDigits(precision uint, exponent int) int {
	return roundTripDecimalDigits(precision) + integerDigits(exponent)
}

func roundTripDecimalDigits(precision uint) int {
	return int((uint64(precision)*30_103+99_999)/100_000) + 1 // #nosec G115 -- Public construction and contexts bound precision by uint32 MaxPrecision; this digit estimate fits even signed 32-bit int.
}

func integerDigits(value int) int {
	var buffer [32]byte
	text := strconv.AppendInt(buffer[:0], int64(value), 10)
	if text[0] == '-' {
		return len(text) - 1
	}
	return len(text)
}

func renderedDigits(text string) int {
	digits := 0
	for _, character := range text {
		if character >= '0' && character <= '9' {
			digits++
		}
	}

	return digits
}

func (c Context) validateContext(ctx context.Context) (big.RoundingMode, gomath.Limits, error) {
	if ctx == nil {
		return 0, gomath.Limits{}, fmt.Errorf("%w: nil context", gomath.ErrInvalidArgument)
	}
	if err := ctx.Err(); err != nil {
		return 0, gomath.Limits{}, err
	}

	return c.validate()
}

func roundingMode(mode gomath.RoundingMode) (big.RoundingMode, bool) {
	switch mode {
	case gomath.RoundHalfEven:
		return big.ToNearestEven, true
	case gomath.RoundHalfUp:
		return big.ToNearestAway, true
	case gomath.RoundDown:
		return big.ToZero, true
	case gomath.RoundUp:
		return big.AwayFromZero, true
	case gomath.RoundCeiling:
		return big.ToPositiveInf, true
	case gomath.RoundFloor:
		return big.ToNegativeInf, true
	default:
		return 0, false
	}
}

func fromBig(value *big.Float, precision uint, rounding gomath.RoundingMode) Float {
	var result Float
	result.value.Copy(value)
	result.precision = precision
	result.rounding = rounding

	return result
}
