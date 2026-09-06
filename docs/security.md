# Security and resource limits

Every parser and material operation accepts or derives validated `Limits` for
input digits, exponent magnitude, precision, output digits, powers,
intermediate bits, and random ranges. Long-running operations accept a context.
Applications handling untrusted input should lower the defaults to the smallest
domain-appropriate values and reject limit errors without echoing the input.
Integer quotient, quotient-with-remainder, and modulus operations require both
a context and limits so division cannot bypass that boundary.

Text parsers reject raw input above fixed overflow-safe bounds before grammar
work: Integer and Decimal use `2*MaxInputDigits + 64` bytes, and Rational uses
`2*MaxInputDigits + 3`. Decimal exponent tokens are limited to eleven bytes
after `e` or `E`; a valid eleven-byte prefix followed by a twelfth token byte
is a limit error. Binary decoders reject payloads above
`MaxIntermediateBits/8 + 64` bytes. At or below those outer bounds, owned
parsers stop at the first conclusive syntax or digit-limit rejection and do not
continue scanning rejected suffixes.

Integer, rational, decimal, and binary-float operations preflight every
operand against the active intermediate-size and exponent budgets before
performing arithmetic. Decimal power-of-ten alignment and exact-quotient
scaling are rejected before materializing an oversized coefficient.

The implementation uses no unsafe code, cgo, ambient randomness, background
goroutines, mutable globals, or hidden caches. Random integers require an
injected `io.Reader` and use unbiased rejection sampling with a bounded number
of attempts. Reader failures have the fixed public message
`math: random source failed`; their cause remains available through
`errors.Is` and `errors.As` without formatting cause text.

Binary-float construction and arithmetic reject source values, operands, and
results whose mantissa precision or exponent exceeds the active limits.
Rendered significant digits are also checked before a `Float` can escape an
operation, so later text and JSON encoding cannot bypass the output budget.

Context-aware operations validate nil, cancellation, and deadline state before
limits and operands. Iterative operations recheck at their documented loop or
read boundaries. Individual `math/big` primitives are synchronous and are not
claimed to be preemptible.
