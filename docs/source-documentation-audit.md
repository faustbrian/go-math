# Source documentation audit

[Documentation index](README.md)

## Scope and acceptance

This audit began before source edits at remote-main base
`0ffddb5ac2a61ec3c26bd0058a27449e5b93f42a`.
The sole module is `github.com/faustbrian/go-math` (`go.mod`); the pinned
local toolchain is Go 1.27.2. This is a non-runtime source-documentation batch:
executable tokens, public signatures, build constraints, and tool directives
must remain unchanged. No release or numerical qualification is implied.

Reproduce the tracked inventory with `git ls-files '*.go'` and
`golib source-docs inventory --json`; compare the package manifest with
`packages.json`. All tracked Go files, not only files containing public APIs,
are in the table below. Public members and private invariant paths are examined
against tests, maintained numeric/security documentation, and the governed
specification decisions. Runtime suites are outside this prose-only batch.

## Initial gaps and contract risks

The initial objective inventory reports 69 findings: field contracts in both
contexts, results and parse options; shared condition/rounding constants;
decimal error aliases; law-helper interface methods; random-source error
methods; and the reserved diagnostic-limit deprecation paragraph. Existing
package comments are present but omit admission, ownership and failure boundaries.

Source-derived clarifications requiring substantive review include:

- Integer and rational values are exact; decimal precision is decimal digits,
  binary-float precision is bits. Low-precision binary construction can round.
- Finite decimal overflow returns a bounded finite value and conditions, not
  infinity; traps retain the value/conditions. Exact intermediate operations
  can exceed limits before context rounding. Rounding and quantization have
  different representation and condition consequences.
- Numeric equality does not imply decimal representation equality. Decimal
  text/JSON preserve numeric value, while binary decimal frames preserve coefficient
  and exponent. The law helper's round-trip assertion only calls `Equal`.
- Constructors/accessors copy mutable big-number storage; value operations do
  not mutate operands. Decoder receivers require exclusive mutation. Copies
  do not imply fresh operation budgets on formatting and accessor methods.
- Context admission is synchronous; most big-number primitives cannot be
  interrupted mid-call. Root search and rejection sampling have checkpoints;
  a caller-provided blocking reader cannot be preempted by cancellation.
- Parser raw envelopes, first-conclusive syntax/digit rejection, exponent
  token limits, output and intermediate budgets are distinct. Random-source
  diagnostics are fixed, but traversing the preserved caller cause is not
  a redaction guarantee for that cause.
- Binary frames have version/type/sign/length/canonicality checks. Re-encoding
  checks canonical representation; frame limits do not budget unbounded
  encoding of an already constructed value.

Authorities: [numeric model](numeric-model.md), [precision](precision.md),
[conditions](conditions.md), [serialization](serialization.md),
[security](security.md), and [specification decisions](specification-decisions.md).
Decision history, fixture hashes, mutation evidence, and archived records are
immutable and are not regenerated for prose changes.

## Package and API surface

- `bigfloat/bigfloat.go` (bigfloat): `Context`, `Context.Precision`, `Context.Rounding`, `Context.Traps`, `Context.Limits`, `Result`, `Result.Value`, `Result.Accuracy`, `Result.Conditions`, `Float`, `NewInt64`, `FromBig`, `FromRat`, `Parse`, `Float.Big`, `Float.Precision`, `Float.Rounding`, `Float.String`, `Float.IsInf`, `Float.Sign`, `Float.Signbit`, `Float.Cmp`, `Float.Equal`, `Float.Neg`, `Float.Abs`, `Float.Int`, `Float.Rat`, `Float.MarshalText`, `Float.MarshalJSON`, `Context.Add`, `Context.Sub`, `Context.Mul`, `Context.Quo`, `Context.Sqrt`.
- `contracts.go` (gomath): `ErrInvalidArgument`, `ErrInvalidSyntax`, `ErrLimitExceeded`, `ErrDivisionByZero`, `ErrDomain`, `ErrConversion`, `ErrOverflow`, `ErrUnderflow`, `ErrTrappedCondition`, `ErrRandomSource`, `Condition`, `ConditionRounded`, `ConditionInexact`, `ConditionOverflow`, `ConditionUnderflow`, `ConditionDivisionByZero`, `ConditionInvalidOperation`, `ConditionClamped`, `ConditionSubnormal`, `Condition.Has`, `Condition.String`, `RoundingMode`, `RoundHalfEven`, `RoundHalfUp`, `RoundHalfDown`, `RoundDown`, `RoundUp`, `RoundCeiling`, `RoundFloor`, `RoundingMode.Valid`, `RoundingMode.String`.
- `decimal/decimal.go` (decimal): `ErrInvalid`, `ErrDivisionByZero`, `ErrNonTerminating`, `ErrLimit`, `RoundingMode`, `HalfEven`, `HalfUp`, `HalfDown`, `Down`, `Up`, `Ceiling`, `Floor`, `ParseOptions`, `ParseOptions.AllowExponent`, `ParseOptions.AllowPlus`, `ParseOptions.AllowUnderscores`, `ParseOptions.AllowLeadingZeros`, `ParseOptions.AllowWhitespace`, `ParseOptions.AllowBarePoint`, `ParseOptions.Limits`, `Context`, `Context.Precision`, `Context.MinExponent`, `Context.MaxExponent`, `Context.Rounding`, `Context.Traps`, `Context.Limits`, `Result`, `Result.Value`, `Result.Conditions`, `Decimal`, `New`, `FromBig`, `Parse`, `ParseWithOptions`, `MustParse`, `Decimal.Coefficient`, `Decimal.Exponent`, `Decimal.Scale`, `Decimal.String`, `Decimal.BigRat`, `Decimal.Sign`, `Decimal.IsZero`, `Decimal.Neg`, `Decimal.Abs`, `Decimal.Cmp`, `Decimal.Equal`, `Decimal.SameRepresentation`, `Decimal.AddExact`, `Decimal.SubExact`, `Decimal.MulExact`, `Decimal.QuoExact`, `Context.Add`, `Context.Sub`, `Context.Mul`, `Context.Quo`, `Context.Apply`, `Decimal.Quantize`, `QuantizedQuo`, `Decimal.Clamp`, `Decimal.MarshalText`, `Decimal.UnmarshalText`, `Decimal.MarshalJSON`, `Decimal.UnmarshalJSON`.
- `encoding/binary.go` (encoding): `MarshalInteger`, `UnmarshalInteger`, `MarshalRational`, `UnmarshalRational`, `MarshalDecimal`, `UnmarshalDecimal`, `MarshalFloat`, `UnmarshalFloat`.
- `integer/integer.go` (integer): `Integer`, `ParseOptions`, `ParseOptions.Base`, `ParseOptions.AllowUnderscores`, `ParseOptions.AllowLeadingZeros`, `ParseOptions.AllowWhitespace`, `ParseOptions.RejectSign`, `ParseOptions.Limits`, `Zero`, `New`, `NewUnsigned`, `FromBig`, `Parse`, `Integer.Big`, `Integer.String`, `Integer.MarshalText`, `Integer.MarshalJSON`, `Integer.Sign`, `Integer.Cmp`, `Integer.Equal`, `Integer.Neg`, `Integer.Abs`, `Integer.Add`, `Integer.Sub`, `Integer.Mul`, `Integer.Quo`, `Integer.QuoRem`, `Integer.Mod`, `Integer.Pow`, `Integer.Root`, `Min`, `Max`, `Clamp`, `GCD`, `LCM`, `Random`, `randomSourceError.Error`, `randomSourceError.Unwrap`, `Integer.Int64`, `Integer.Uint64`.
- `limits.go` (gomath): `Limits`, `Limits.MaxInputDigits`, `Limits.MaxOutputDigits`, `Limits.MaxExponentMagnitude`, `Limits.MaxPrecision`, `Limits.MaxPowerExponent`, `Limits.MaxRootDegree`, `Limits.MaxRandomBits`, `Limits.MaxRandomAttempts`, `Limits.MaxIntermediateBits`, `Limits.MaxDecimalExpansion`, `Limits.MaxDiagnosticBytes`, `DefaultLimits`, `Limits.Validate`.
- `mathtest/laws.go` (mathtest): `T`, `T.Helper`, `T.Errorf`, `EqualValue`, `EqualValue.Equal`, `Operation`, `Commutative`, `Associative`, `Identity`, `RoundTrip`.
- `rational/rational.go` (rational): `Rational`, `Zero`, `New`, `NewChecked`, `FromBig`, `Parse`, `Rational.Big`, `Rational.Numerator`, `Rational.Denominator`, `Rational.String`, `Rational.MarshalText`, `Rational.MarshalJSON`, `Rational.Sign`, `Rational.Cmp`, `Rational.Equal`, `Rational.Neg`, `Rational.Abs`, `Rational.Add`, `Rational.Sub`, `Rational.Mul`, `Rational.Quo`, `Rational.Pow`, `Rational.Decimal`, `Min`, `Max`, `Clamp`.

## Tracked Go coverage

Every file below is included in the source/test review. Test names and assertions
supply rationale, not fresh execution claims. Public-package tests include
strict parsing, source/return alias isolation, numeric laws, context precedence,
rounding modes, condition/trap behavior, codec rejection and resource boundaries.
Differential tests use bounded overlapping APD, Shopspring, and math/big behavior;
they do not establish universal conformance. Decimal GDA tests explicitly bind
supported decisions and skip unsupported corpus semantics.

| File | Role | Variant |
| --- | --- | --- |
| `allocation_test.go` | test | ordinary |
| `benchmark_test.go` | test | ordinary |
| `bigfloat/aaa_mutation_contract_test.go` | test | ordinary |
| `bigfloat/bigfloat.go` | source | ordinary |
| `bigfloat/bigfloat_test.go` | test | ordinary |
| `bigfloat/cohesion_contract_test.go` | test | ordinary |
| `bigfloat/differential_test.go` | test | ordinary |
| `bigfloat/edge_test.go` | test | ordinary |
| `bigfloat/fuzz_test.go` | test | ordinary |
| `bigfloat/internal_edge_test.go` | test | ordinary |
| `bigfloat/leak_test.go` | test | ordinary |
| `contracts.go` | source | ordinary |
| `contracts_edge_test.go` | test | ordinary |
| `contracts_test.go` | test | ordinary |
| `decimal/bare_point_test.go` | test | ordinary |
| `decimal/cohesion_contract_test.go` | test | ordinary |
| `decimal/context_test.go` | test | ordinary |
| `decimal/decimal.go` | source | ordinary |
| `decimal/decimal_test.go` | test | ordinary |
| `decimal/differential_test.go` | test | ordinary |
| `decimal/edge_internal_test.go` | test | ordinary |
| `decimal/fuzz_test.go` | test | ordinary |
| `decimal/gda_test.go` | test | ordinary |
| `decimal/hardening_allocation_test.go` | test | build-constrained |
| `decimal/hardening_test.go` | test | ordinary |
| `decimal/leak_test.go` | test | ordinary |
| `decimal/mutation_contract_test.go` | test | ordinary |
| `decimal/precision_drop_internal_test.go` | test | ordinary |
| `decimal/quantize_test.go` | test | ordinary |
| `decimal/ratio_internal_test.go` | test | ordinary |
| `decimal/security_budget_test.go` | test | ordinary |
| `encoding/aliasing_test.go` | test | ordinary |
| `encoding/binary.go` | source | ordinary |
| `encoding/binary_test.go` | test | ordinary |
| `encoding/cohesion_contract_test.go` | test | ordinary |
| `encoding/edge_internal_test.go` | test | ordinary |
| `encoding/fuzz_test.go` | test | ordinary |
| `encoding/leak_test.go` | test | ordinary |
| `encoding/mutation_contract_test.go` | test | ordinary |
| `example_test.go` | test | ordinary |
| `integer/cohesion_contract_test.go` | test | ordinary |
| `integer/differential_test.go` | test | ordinary |
| `integer/edge_test.go` | test | ordinary |
| `integer/fuzz_test.go` | test | ordinary |
| `integer/hardening_allocation_test.go` | test | build-constrained |
| `integer/integer.go` | source | ordinary |
| `integer/integer_test.go` | test | ordinary |
| `integer/internal_edge_test.go` | test | ordinary |
| `integer/leak_test.go` | test | ordinary |
| `integer/mutation_contract_test.go` | test | ordinary |
| `integer/security_budget_test.go` | test | ordinary |
| `laws_test.go` | test | ordinary |
| `leak_test.go` | test | ordinary |
| `limits.go` | source | ordinary |
| `mathtest/edge_test.go` | test | ordinary |
| `mathtest/laws.go` | source | ordinary |
| `mathtest/laws_test.go` | test | ordinary |
| `mathtest/leak_test.go` | test | ordinary |
| `race_disabled_test.go` | test | build-constrained |
| `race_enabled_test.go` | test | build-constrained |
| `race_test.go` | test | ordinary |
| `rational/cohesion_contract_test.go` | test | ordinary |
| `rational/differential_test.go` | test | ordinary |
| `rational/edge_test.go` | test | ordinary |
| `rational/fuzz_test.go` | test | ordinary |
| `rational/hardening_allocation_test.go` | test | build-constrained |
| `rational/internal_edge_test.go` | test | ordinary |
| `rational/leak_test.go` | test | ordinary |
| `rational/mutation_contract_test.go` | test | ordinary |
| `rational/rational.go` | source | ordinary |
| `rational/rational_test.go` | test | ordinary |
| `rational/security_budget_test.go` | test | ordinary |
| `resource_test.go` | test | ordinary |

## Variants, tooling and numerical assurance

- `race_enabled_test.go` / `race_disabled_test.go` select the race indicator.
  Integer, rational and decimal allocation-boundary tests are excluded from
  race builds; those intentionally different assumptions are preserved.
- All seven fuzz functions, thirteen benchmark functions, two executable Examples,
  differential/allocation/leak/race assertions and their oracles are reviewed.
  Four tracked saved fuzz regressions are also inspected:
  `bigfloat/testdata/fuzz/FuzzParse/e42a55ed93908e6f`,
  `decimal/testdata/fuzz/FuzzParseContextAndRoundTrip/ae405b7f3b9fc4fe`,
  `decimal/testdata/fuzz/FuzzParseContextAndRoundTrip/af1ecf1d2149bf0a`, and
  `integer/testdata/fuzz/FuzzParseAndArithmetic/14ffc53d6bb612b2`.
  They exercise low precision, zero/large exponents and digit boundaries.
  There is no generated Go source or generation directive.
  Existing security-analysis directives and specification bindings stay intact.
- The shared API baseline uses `apidiff` (`api/baseline.txt`); native
  `scripts/check-api-compat.sh` compares `go doc -short` with
  `compat/public-api.txt`. Prose must not change those public signatures.
- Selected verification is parsed AST/executable-token and directive equivalence,
  formatting, native documentation/examples and structure/API checks, rendered
  package contracts, local links, source-documentation inventory, changed-content
  secret inspection and the complete final diff. Compile-heavy work requires
  the coordinator's serialized slot and three disposable task-owned Go caches.
- Existing historical conformance, mutation and performance evidence is not
  recertified. No new runtime, race, fuzz, mutation or benchmark campaign is
  needed to verify unchanged executable behavior.

## Investigated rationale and unresolved questions

The source explains admission and output budgets separately, finite decimal
conditions versus error sentinels, exact-versus-rounded boundaries, alias copies,
and why canonical binary validation re-encodes. Those contracts are grounded in
implementation and boundary tests rather than inferred from function names.

All 72 tracked files have been read, including the 64 test sources. The 226
exported declarations/members include two exported methods of the private
random-source error type. Member contracts are repaired at their owners rather
than copied into tests.
Unchanged private helpers retain meaningful existing rationale for saturated
bounds, checked exponent narrowing, virtual-zero ratio comparison, carry-bit
preflights and non-materialized underflow divisors.

Existing binary-float exceptional infinity combinations delegate to math/big
and can panic with big.ErrNaN: opposite infinities in Add, same-signed infinities
in Sub, zero times infinity in Mul, and two infinities in Quo. Callers must reject
these combinations; translating them into conditions is a separate behavior
change, not part of this batch. Directed-rounding text recovery and the binary
codec's rejection of a zero-precision Float are also explicit limits, not new
guarantees. The bounded Go 1.27.2 characterization confirmed all four infinity
panic cases.
At precision four with RoundCeiling, FromRat(1/3) stores 11/32, String emits
`0.344`, and parsing that text with the same context stores 3/8. Encoding then
decoding an unconstructed Float returns ErrInvalidArgument for binary precision.
NewInt64(3) at precision one and HalfEven returns 4 with Rounded and Inexact.
These are existing behaviors, not fixes or new round-trip/conformance promises.

Decimal.Clamp reports ErrInvalid/ErrInvalidSyntax for an inverted interval,
unlike the integer/rational ErrInvalidArgument policy. The comment preserves
that existing classification. No runtime fix is authorized by this batch.

## Final selected verification

All selected checks passed against the comment batch using Go 1.27.2,
`GOWORK=off`, bounded compilation parallelism and task-owned build/module/temp
caches:

- Pinned `gofmt -l` over all 72 tracked Go files: no findings.
- Native Go parser/AST comparison excluding comments and positions, scanner
  executable-token equality, and exact directive comparison against the recorded
  base: all 72 files unchanged in executable structure; 226 members inventoried.
- `make -f verification/package.mk architecture`: passed.
- `make -f verification/package.mk api`: passed; native go-doc-short snapshot
  and shared apidiff baseline remain untouched because signatures are unchanged.
- `make -f verification/package.mk docs`: passed. Local Markdown links/anchors,
  documented Make targets, module/package identities, public package docs and
  attached Example output oracles reconcile. Its native `-run '^Example'`
  executes the two root Examples; the other six packages have no Examples.
- `go doc -all` for all seven packages: rendered and inspected, especially the
  limits/options, decimal/binary contexts, trap results, codec and law-helper
  contracts. No raw output is retained as a new evidence framework.
- `golib source-docs check --json`: zero findings, down from the initial 69.
- `git diff --check`, build/generator/debt-marker inventory, local-link review,
  changed-content secret inspection and ordinary complete-diff review: passed.

The three owned Go caches, package-manager cache and newly installed
`node_modules` were removed immediately after the bounded run. Disposable
comparison/characterization helpers and rendered output are removed after their
results are captured. The detached committed worktree is retained for coordinator
delivery; no publication or module release is performed by this worker.

No documentation blocker remains. The exceptional infinity panic and
non-unconditional directed-rounding text recovery are follow-up behavior decisions
for the owner, not unresolved prose claims. Full arithmetic/GDA/peer/runtime,
race/leak, fuzz, mutation, benchmark, reverse-consumer and release campaigns were
not run for unchanged executable inputs. Historical qualification remains
historical; exact-source hosted CI and remote-main delivery belong to the
coordinator.

## Objective enforcement adoption

The local `source-docs-check` target is required by ordinary/full checks and CI.
Ordinary shared CI uses `source_documentation: check` with immutable Tools
v2.4.0 and matching authenticated binary/checksum declaration. Separately
selected qualification source identities and gates remain unchanged. Structural
checks do not replace this maintained semantic audit or native docs/examples.
