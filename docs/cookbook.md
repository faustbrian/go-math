# Cookbook

[Documentation index](README.md)

Parse an exact amount with `decimal.Parse`, then use `Quantize` for a fixed
fractional scale. Use `Context.Quo` when a quotient may repeat and inspect
`Result.Conditions`; use `QuoExact` when repetition must fail. Construct ratios
with `rational.New` and call `Decimal` with a scale and rounding mode for a
bounded expansion. Use `integer.Random` with `crypto/rand.Reader` in production
or a deterministic reader in tests.

To accept looser external text, enable only the `decimal.ParseOptions` fields
its grammar needs. A bare point such as `".5"`, `"5."` or `"5.e3"` needs
`AllowBarePoint` and keeps the scale of the digits written: `"5."` has scale 0
and `".50"` has scale 2. PHP `is_numeric` decimal strings also need
`AllowExponent`, `AllowPlus`, `AllowLeadingZeros` and `AllowWhitespace`.
The parser still differs from PHP in two ways. First, `AllowWhitespace` trims
all Unicode whitespace, while PHP trims only ASCII space, tab, newline,
carriage return, vertical tab and form feed, so exact matching needs an ASCII
pre-check. Second, the parser is bounded where PHP is not: the exponent token
is capped at a fixed eleven bytes, independent of `Limits`, and
`Limits.MaxInputDigits`, `Limits.MaxExponentMagnitude` and the raw-byte bound
of twice `MaxInputDigits` plus 64 reject long inputs that PHP would accept.

For shared configuration, construct a context once and pass it by value. Values
and contexts contain no mutable shared state exposed by the API.

