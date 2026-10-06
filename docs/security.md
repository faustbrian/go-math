# Security and resource limits

## Threat model

The model covers the v1 root module and its Integer, Rational, Decimal,
binary-float, and binary-codec packages. Protected assets are numeric identity,
rounding and condition semantics, application CPU and memory, and input
confidentiality. Callers own input transport, concurrency, deadlines, and
selection of trusted resource policy; numeric text, JSON strings, binary frames,
and operand values can be attacker-controlled. The library owns no network,
filesystem, database, process, credential, or background-worker boundary.
Numeric operations do not implement authentication or authorization.

## Owned boundaries

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

Decimal JSON decoding applies a raw-byte bound before decoding or copying the
string: at most six encoded bytes per parser-envelope byte, with bounded
delimiter slack. Fully escaped digit strings remain supported. This does not
bound the caller's prior allocation of the input `[]byte` or work performed by
an enclosing JSON decoder; applications must cap transport and document bytes.

Integer, rational, decimal, and binary-float operations preflight every
operand against the active intermediate-size and exponent budgets before
performing arithmetic. Decimal power-of-ten alignment and exact-quotient
scaling are rejected before materializing an oversized coefficient.
Derived decimal exponents are checked before narrowing to `int32`, and overflow
clamping preflights its power-of-ten coefficient. Rounding a coefficient by more
places than it contains derives zero or a directed unit without creating the
divisor. Integer root comparisons use budgeted repeated squaring and discard
powers above the candidate value. Rational decimal scaling uses overflow-safe
base-ten growth
accounting and checks the resulting coefficient; as with ordinary bounded
arithmetic, a temporary multiplication can include one carry bit beyond the
preflight estimate, but that bit cannot escape as an accepted result.

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

## Accepted residual risks

| Risk | Severity | Owner | Rationale | Mitigation | Review condition |
| --- | --- | --- | --- | --- | --- |
| An injected `io.Reader` can block cancellation inside `Integer.Random`. | Medium | Caller integrating the random source | The library cannot preempt an arbitrary blocking read without creating an unbounded goroutine lifetime. | Supply a bounded or cancellation-aware source and enforce its deadline outside the call. | Revisit if the reader contract changes or an owned asynchronous random API is introduced. |
| A synchronous `math/big` primitive can finish after cancellation is requested. | Medium | go-math maintainers | These primitives do not expose safe interruption points. | Preflight operands and intermediate budgets; observe cancellation between root candidates and random attempts. Each root power comparison has at most 32 budgeted squaring steps. | Revisit when Go exposes interruptible primitives or an affected algorithm is replaced. |
| Caller-selected `Limits` can authorize large CPU, memory, or rendered output. | Medium | Service or library integrator | Limits are explicit trusted policy rather than untrusted request data. | Cap configuration independently of requests and use the smallest domain-appropriate values. | Revisit if limit values become request-controlled or any default limit increases. |

## Release and verification boundary

The false-default `root_release_qualification` workflow selection runs the
complete root release dry-run with immutable source-built development tooling
at `c799f67c5f4ac074e24a0172fb41a96cce9dad8e`. It also requires ordinary CI;
it is separate from both stable `release_dry_run` and the diagnostic-only
`native_security_qualification`. Combined selections are refused. This route
does not publish a release or claim a public tooling version.

The four decimal equivalents at source `0fdbdd89ba6507fb140038a85888dd47576f3d43`
were independently reviewed by `cl_report` on 2026-10-06 for the bounded public
caller and admission domain. Review expires on 2027-01-02 or any source, caller,
admission, or verifier change. Historical run `37388182213` used the published
verifier `5eba124f842305aecef71fc439aea0f1056429b253fad3b2ed47468f614595a3`
against identical pre-refactor decimal production bytes and recorded four LIVED
equivalents and two loop-negation TIMED OUT failures. Those historical outcomes
remain unchanged; current-source qualification with verifier
`c28707fd4ce35dd228260de564107c1fb7726e42d4e9ea2bf9e988083e1f3b5b`
is pending.

Ratio comparison now computes its comparison ceiling once from the larger
actual digit count, not the exponent-adjusted virtual length. Inner digit reads
and numeric comparison semantics are unchanged. The
[supported inventory](../.verification/mutation/equivalent-inventory.json)
retains only the unchanged coordinates at lines 677 and 1180, with the new
source binding; the two former loop-boundary coordinates no longer exist.
The source change expires the previous whole-source review, so the current
binding and retained dispositions require independent review and new hosted
qualification. No new mutant disposition, timeout conversion, or release-pass
verdict is asserted.

The v1 security patch retains exported API names, numeric representations,
rounding modes, and error classification while rejecting crossed budgets and
wrapped exponents. Go 1.27 remains the supported floor. A source fix is not a
public release: release acceptance also requires required CI and security
checks, signed tags, a clean proxy/SumDB consumer, and adoption by direct owned
consumers. Dependency and secret scans do not establish arithmetic bounds;
focused hostile-input, allocation, exponent, rounding, and cancellation cases
exercise those contracts. See [the security policy](../SECURITY.md) for private
disclosure.
