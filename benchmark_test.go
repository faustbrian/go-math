package gomath_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/cockroachdb/apd/v3"
	gomath "github.com/faustbrian/go-math"
	"github.com/faustbrian/go-math/bigfloat"
	"github.com/faustbrian/go-math/decimal"
	mathencoding "github.com/faustbrian/go-math/encoding"
	"github.com/faustbrian/go-math/integer"
	"github.com/faustbrian/go-math/rational"
	shopspring "github.com/shopspring/decimal"
)

// BenchmarkIntegerAdd measures bounded immutable addition of two nine-digit
// integers, including admission and returned storage. Operand and limit
// construction precede b.Loop; the BigInt comparison omits defensive policy
// checks.
func BenchmarkIntegerAdd(b *testing.B) {
	left := integer.New(123456789)
	right := integer.New(987654321)
	limits := gomath.DefaultLimits()
	b.ReportAllocs()
	for b.Loop() {
		if _, err := left.Add(context.Background(), right, limits); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBigIntAdd measures raw math/big addition into a fresh destination
// for the same nine-digit inputs as BenchmarkIntegerAdd. Inputs are constructed
// before b.Loop; it does not perform the wrapper's context/limit admission.
func BenchmarkBigIntAdd(b *testing.B) {
	left := big.NewInt(123456789)
	right := big.NewInt(987654321)
	b.ReportAllocs()
	for b.Loop() {
		new(big.Int).Add(left, right)
	}
}

// BenchmarkDecimalAdd measures context-rounded addition of two pre-parsed
// 13-digit coefficients at precision 18. Input/context setup is outside b.Loop;
// APD uses matching inputs/context, while Shopspring has different precision
// policy.
func BenchmarkDecimalAdd(b *testing.B) {
	left := decimal.MustParse("123456789.1234")
	right := decimal.MustParse("987654321.9876")
	operation := decimal.Context{
		Precision: 18, MinExponent: -100, MaxExponent: 100,
		Rounding: decimal.HalfEven,
	}
	b.ReportAllocs()
	for b.Loop() {
		if _, err := operation.Add(context.Background(), left, right); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkAPDAdd measures APD addition into a fresh result at precision 18
// with the same pre-parsed operands as BenchmarkDecimalAdd. Setup is outside
// b.Loop; error/trap and defensive-resource policies differ between the
// libraries.
func BenchmarkAPDAdd(b *testing.B) {
	left, _, _ := apd.NewFromString("123456789.1234")
	right, _, _ := apd.NewFromString("987654321.9876")
	operation := apd.Context{
		Precision: 18, MinExponent: -100, MaxExponent: 100,
		Rounding: apd.RoundHalfEven,
	}
	b.ReportAllocs()
	for b.Loop() {
		if _, err := operation.Add(new(apd.Decimal), left, right); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkShopspringAdd measures Shopspring addition of the same pre-parsed
// decimal values as BenchmarkDecimalAdd. Setup is outside b.Loop; this
// exact-add boundary does not apply the other benchmarks' significant-digit
// context.
func BenchmarkShopspringAdd(b *testing.B) {
	left := shopspring.RequireFromString("123456789.1234")
	right := shopspring.RequireFromString("987654321.9876")
	b.ReportAllocs()
	for b.Loop() {
		_ = left.Add(right)
	}
}

// BenchmarkBoundedIntegerPower measures bounded 3^1024 with a new returned
// integer. Base and limits are prepared outside b.Loop; no unbounded peer
// operation or adversarial rejection is measured.
func BenchmarkBoundedIntegerPower(b *testing.B) {
	value := integer.New(3)
	limits := gomath.DefaultLimits()
	b.ReportAllocs()
	for b.Loop() {
		if _, err := value.Pow(context.Background(), 1_024, limits); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBoundedIntegerRoot measures bounded degree-17 root search of 2^4096.
// Source construction is outside b.Loop; each iteration includes admission and
// the search, not a math/big square-root comparison.
func BenchmarkBoundedIntegerRoot(b *testing.B) {
	value, err := integer.FromBig(new(big.Int).Lsh(big.NewInt(1), 4_096), gomath.DefaultLimits())
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for b.Loop() {
		if _, err := value.Root(context.Background(), 17, gomath.DefaultLimits()); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkDecimalDivision measures a quotient of pre-parsed 29- and 10-digit
// coefficients at decimal precision 34. Setup is outside b.Loop; context
// admission, scaling and rounding remain inside.
func BenchmarkDecimalDivision(b *testing.B) {
	numerator := decimal.MustParse("12345678901234567890.123456789")
	denominator := decimal.MustParse("7.000000001")
	operation := decimal.Context{
		Precision: 34, MinExponent: -1_000, MaxExponent: 1_000,
		Rounding: decimal.HalfEven,
	}
	b.ReportAllocs()
	for b.Loop() {
		if _, err := operation.Quo(context.Background(), numerator, denominator); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkDecimalFormatting measures expanded String allocation for a
// pre-parsed 40-digit coefficient with scale 20. Parsing is outside b.Loop; it
// is not an arithmetic or resource-admission benchmark.
func BenchmarkDecimalFormatting(b *testing.B) {
	value := decimal.MustParse("12345678901234567890.12345678901234567890")
	b.ReportAllocs()
	for b.Loop() {
		_ = value.String()
	}
}

// BenchmarkRationalNormalization measures copying/reduction of 22*2^2048 over
// 7*2^2048. Components and limits are constructed outside b.Loop; NewChecked
// admission and normalization are measured together.
func BenchmarkRationalNormalization(b *testing.B) {
	factor := new(big.Int).Lsh(big.NewInt(1), 2_048)
	numerator := new(big.Int).Mul(big.NewInt(22), factor)
	denominator := new(big.Int).Mul(big.NewInt(7), factor)
	limits := gomath.DefaultLimits()
	b.ReportAllocs()
	for b.Loop() {
		if _, err := rational.NewChecked(numerator, denominator, limits); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRationalDecimalExpansion measures 1/7 rendered to 1000 fractional
// places using half-even and explicit limits. The fraction/limits are prepared
// outside b.Loop; scaling, rounding and string construction remain inside.
func BenchmarkRationalDecimalExpansion(b *testing.B) {
	value, err := rational.New(1, 7)
	if err != nil {
		b.Fatal(err)
	}
	limits := gomath.DefaultLimits()
	b.ReportAllocs()
	for b.Loop() {
		if _, _, err := value.Decimal(1_000, gomath.RoundHalfEven, limits); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBigFloatSquareRootAndConversion measures a 4096-bit square root of
// two followed by exact Rat extraction. Context/value construction is outside
// b.Loop; the combined result is not a standalone square-root timing.
func BenchmarkBigFloatSquareRootAndConversion(b *testing.B) {
	operation := bigfloat.Context{
		Precision: 4_096, Rounding: gomath.RoundHalfEven,
		Limits: gomath.DefaultLimits(),
	}
	value, err := bigfloat.FromRat(big.NewRat(2, 1), operation)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for b.Loop() {
		result, err := operation.Sqrt(context.Background(), value.Value)
		if err != nil {
			b.Fatal(err)
		}
		if _, err := result.Value.Rat(); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBinaryConversionRoundTrip measures decimal frame encoding plus
// bounded decoding for a 40-digit coefficient at scale 20. Parsing/limits are
// outside b.Loop; both buffer construction and canonicality admission are
// included.
func BenchmarkBinaryConversionRoundTrip(b *testing.B) {
	limits := gomath.DefaultLimits()
	value := decimal.MustParse("12345678901234567890.12345678901234567890")
	b.ReportAllocs()
	for b.Loop() {
		data, err := mathencoding.MarshalDecimal(value)
		if err != nil {
			b.Fatal(err)
		}
		if _, err := mathencoding.UnmarshalDecimal(data, limits); err != nil {
			b.Fatal(err)
		}
	}
}
