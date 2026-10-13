# Conditions and traps

[Documentation index](README.md)

Decimal results report `Rounded`, `Inexact`, `Overflow`, `Underflow`,
`Subnormal`, and `DivisionByZero`. Conditions remain available even when a
matching result-condition trap returns an error. Decimal division by zero
returns `ErrDivisionByZero` directly, even when that condition is in the trap
mask. Finite decimal overflow returns a value at the maximum permitted magnitude
with conditions rather than infinity; resource limits may reject that result.
Callers must decide whether a rounded result is acceptable; financial code commonly traps inexact conversion at
domain boundaries and explicitly quantizes display or settlement amounts.

