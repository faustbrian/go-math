# math

[![CI](https://github.com/faustbrian/go-math/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/faustbrian/go-math/actions/workflows/ci.yml)
[![CodeQL](https://img.shields.io/badge/CodeQL-required-blue)](https://github.com/faustbrian/go-math/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Mutation](https://img.shields.io/badge/mutation-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Documentation](https://img.shields.io/badge/docs-checked_in_CI-blue)](docs/)
[![Go Reference](https://pkg.go.dev/badge/github.com/faustbrian/go-math.svg)](https://pkg.go.dev/github.com/faustbrian/go-math)
[![Release](https://img.shields.io/github/v/release/faustbrian/go-math?sort=semver)](https://github.com/faustbrian/go-math/releases)
[![Go](https://img.shields.io/badge/go-1.26.6-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`math` is an immutable arbitrary-precision numeric foundation for Go. It
provides distinct APIs for signed integers, exact rationals, finite base-10
decimals, and explicitly inexact binary floats. Use ordinary Go numeric types
when their fixed width and machine arithmetic are sufficient.

```go
amount := decimal.MustParse("19.995")
result, err := amount.Quantize(
	context.Background(), 2, decimal.HalfEven, gomath.DefaultLimits(),
)
```

All constructors copy mutable `math/big` inputs. Operations return new values,
accessors return copies, JSON uses strings, and potentially expensive work is
bounded by `gomath.Limits`. Decimal contexts make precision, exponent range,
rounding, conditions, and traps explicit. No conversion passes through
`float64`.

See the [documentation index](docs/README.md),
[specification decisions](docs/specification-decisions.md),
[cookbook](docs/cookbook.md), and [verification guide](docs/verification.md).
The minimum supported toolchain is Go 1.26.6.
See the versioned [Golib ecosystem index](https://github.com/faustbrian/go-library-tools/blob/v1.3.0/docs/ecosystem/README.md)
and [package-family selection guidance](https://github.com/faustbrian/go-library-tools/blob/v1.3.0/docs/ecosystem/design-language.md#package-families-and-selection)
for the shared design language this module follows.

## Packages

- `integer`: exact signed integer arithmetic, roots, GCD/LCM, and unbiased
  injected-source random values.
- `rational`: normalized exact fractions and bounded decimal conversion.
- `decimal`: finite coefficient/exponent values, exact operations, and
  context-rounded operations with conditions.
- `bigfloat`: explicit precision and rounding around `math/big.Float`.
- `encoding`: deterministic versioned binary codecs.
- `mathtest`: reusable algebraic-law and round-trip assertions.

## Development

Run `make cohesion` for the repository-owned cohesion contract, `make check`
for blocking local gates, and `make check-all` to include the advisory NilAway
analysis. See [CHANGELOG.md](CHANGELOG.md) for releases and [SECURITY.md](SECURITY.md)
for vulnerability reporting.

Licensed under MIT.

## Documentation

Use the [documentation index](docs/README.md) for package-owned guides,
operational contracts, examples, and maintainer references.
