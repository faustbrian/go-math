# Performance

Exact arithmetic cost scales with operand size. Decimal exponent alignment,
division, powers, and roots can allocate substantially, so limits are part of
the API rather than an operational afterthought. Reuse immutable values and
contexts freely, but do not raise limits for untrusted input.

Run the bounded benchmark suite:

```sh
go test -run '^$' -bench '^Benchmark' -benchtime=100ms -benchmem \
  . ./bigfloat ./decimal ./encoding ./integer ./mathtest ./rational
```

The comparisons use direct `math/big`, `apd`, and `shopspring/decimal`
operations. Benchmark semantics are aligned where possible; different rounding
or representation contracts are identified in benchmark names.

The suite covers bounded powers, roots, division, rational normalization,
decimal expansion, formatting, binary-float square roots, and encoding
conversion. See the [recorded baseline](benchmark-baseline.md). Allocation
budgets are deterministic gates; wall-clock results are recorded but are not
used as noisy shared-runner assertions.
