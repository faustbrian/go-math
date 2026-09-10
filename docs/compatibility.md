# Compatibility

Decimal parsing, arithmetic, condition, quantize, serialization, corpus, and
peer-comparison boundaries are recorded in the
[specification decision register](specification-decisions.md). A changed
resolved decision is a compatibility event even when exported Go signatures do
not change.

The module requires Go 1.27.0 and uses the same toolchain in CI. Public API
changes are checked against `api/baseline.txt`. Binary encodings carry a version
byte; unknown versions fail. Text and JSON forms are canonical strings.

v1.1.0 narrows only hostile raw forms that previously bypassed configured work
bounds: arbitrarily long Integer whitespace, oversized Rational envelopes, and
arbitrarily zero-padded Decimal exponents. It also makes syntax and configured-
limit rejection deterministic at the first conclusive boundary. Consumers
should classify with `errors.Is`; exact strings are safe diagnostics, not the
classification API.

The root package identifier remains `gomath` as an intentional v1
compatibility exception that avoids collision with the standard-library
`math` package.

`money` and `measurement` should import only the numeric family they need
and keep currency, unit, locale, and presentation policy in their own domains.
