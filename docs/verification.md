# Verification

`make check` runs formatting, module tidiness, vet, architecture policy, tests,
race detection, exact production coverage, fuzz smoke, mutation checks,
examples, consumer compatibility, benchmarks, docs, API compatibility, lint,
static analysis, vulnerability analysis, and advisory NilAway. `make ci` adds
the repository and cohesion contracts.

The test suite includes algebraic laws, direct `math/big` comparisons,
independent `apd` and `shopspring/decimal` comparisons, applicable General
Decimal Arithmetic vectors, serialization round trips, alias checks, fuzz
targets, and concurrent use. Official vector provenance is recorded under
`specification/gda`; the exact source pins, decision records, conformance
bindings, and monitored release authority are recorded in the
[specification matrix](../specification/README.md) and
[decision register](specification-decisions.md).

The GDA harness executes and exactly accounts for 3,547 applicable vectors and
1,542 explicit skips across addition, subtraction, multiplication, division,
quantization, and rounding files. Non-finite values, unsupported conditions,
and non-extended operand pre-rounding are outside the finite extended-decimal
contract. Any accounting drift fails the suite.

Mutation testing uses Gremlins v0.6.0 over every production package and requires
both 100% test efficacy and 100% mutant coverage. Provenance checks verify every
vendored vector checksum before the release gate proceeds.

The maintained-peer lane evaluates 61 operand pairs. APD v3.2.3 contributes
305 value-and-shared-condition comparisons across add, subtract, multiply,
divide, and quantize. shopspring/decimal v1.4.0 contributes 183 exact-value
comparisons across add, subtract, and multiply. The resulting 488 comparisons
do not cover parsing grammar, non-finite values, signed zero, exponent edge
conditions, traps, serialization, defensive limits, cancellation, or corpus
skip classification. The exact evidence and gaps are machine-readable in
[`maintained-peers.json`](../specification/maintained-peers.json).

`golib specification check` validates the complete offline governance record.
`golib specification check --online` additionally verifies every reviewed
authority digest against its bounded HTTPS source.
