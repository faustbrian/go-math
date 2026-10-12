# Contributing

## Before Editing

1. Read [`AGENTS.md`](AGENTS.md) and the affected package's documentation.
2. Run ordinary Go tests for the affected package before editing.
3. Identify owned dependencies and reverse dependants in `modules.json`.
4. Preserve unrelated work and generated/corpus provenance.

## Changes

Keep commits focused and conventional. Update every affected changelog with
the behavior and migration impact. Public API changes require compatibility
evidence and documentation. Specification behavior requires a decision record,
fixture coverage, and interoperability evidence.

New direct dependencies and dependency updates must follow the
[dependency governance policy](AGENTS.md#dependencies-and-supply-chain). Package-local
update bots are forbidden; the root policy owns every module and action update.

Specification-backed changes must follow the
[specification governance contract](AGENTS.md#design), update
the affected stable entries in the
[specification decision register](docs/specification-decisions.md), and complete
the Specification Decisions section of the pull request template. An
unresolved interpretation or stale source pin is release-blocking; peer
behavior cannot silently select policy. Superseded identifiers and all prior
digests remain in the register and history instead of being erased.

Selected mutation gates must finish with zero surviving viable mutants.

Do not fork shared verification workflows or add permanent replacements,
machine-specific paths, bypass flags, broad mutation exclusions or aggregate
quality metrics that hide a failing package.

## Verification

Use Go 1.27.2 for development and ordinary CI; the public minimum remains
Go 1.27.0. An ordinary checkout is sufficient for tests and a focused active
fuzz smoke, without sibling repositories or goal ledgers:

```bash
go test ./...
go test ./decimal -run '^$' -fuzz '^FuzzParseContextAndRoundTrip$' -fuzztime=2s
```

Select the native target and package affected by your change; ordinary Go tests
already exercise deterministic fuzz seeds. For shared checks, install the public
`golib` v2.4.0 release declared in `.golib.yaml` and verify its checksums:

```bash
make check
make ci
```

Both Make routes use bounded local checks. The module declares no services or
test tags. Optional mutation analysis can select one affected package:

```bash
golib mutation --module . --package ./decimal
```

Use `make check-full` for explicit full assurance with the original native fuzz
budgets, API, architecture, provenance, decimal conformance and maintained-peer
comparisons. Preserve the authority fixtures and explicit skip accounting.
`golib specification check` verifies governance offline; `--online` additionally
checks the bounded authoritative sources. Select directly affected numeric
consumers for contract changes. Explicit full and release qualification retain
their full selected gates and release rehearsal. The existing native-security
and root-release qualification dispatches retain their separate source-built
verifiers; neither changes the ordinary public verifier identity. Report failing
or unavailable checks without claiming release readiness. CI/configuration-only
changes do not require a library release.

## Selected Monthly Assurance

The owner-level [selected-assurance caller](.github/workflows/assurance.yml) is
configured for main only, on day 22 of each month at 07:43 UTC. It selects one
of the six declared root-module targets using the shared UTC-week rotation,
with 30 seconds of active exploration. The production packages are integer,
rational, decimal, bigfloat and encoding; the selected risk is bounded numeric
grammar, exponent/precision handling, terminal error categories and canonical
text/JSON/binary representations. Inline seeds and saved regressions remain
ordinary tests. A small rotation explores inputs beyond that fixed coverage
without repeating the full numerical/conformance campaign.

Base-free ordinary checks execute deterministic seeds, not active exploration;
changed-causal PR checks select affected active smoke. Manual assurance requires
an exact anchored selector from the root-module fuzz operation in `.golib.yaml`,
for example `^FuzzParseContextAndRoundTrip$`.
Empty, unknown and malformed selections refuse. The shared job bounds setup
and execution to ten minutes, active workers and CPU to two, owned memory to
2 GiB, native output to 1 MiB and private-cache lifetime. Failure, cancellation,
exhausted budgets and cleanup refusal remain attributable failures.

This lane does not run mutation or replace Required, the original six
10,000-iteration fuzz budgets, explicit `make check-full`, GDA vector/skip
accounting, maintained-peer conformance, architecture/provenance/API checks,
native-security diagnostics or relevant public-consumer/release qualification.
The existing daily ordinary schedule remains bounded and monitors sources
online; it is not the monthly exploration or a full qualification campaign.

## Adding A Module

Follow [repository structure policy](AGENTS.md#repository-structure). New modules
require an explicit purpose, ownership boundary, dependency review, package
catalog entry, full quality gates, documentation, changelog, license, security
policy, compatibility plan, and release dry-run.

## Source documentation

Use the authenticated Tools v2.4.0 release declared in `.golib.yaml`.
`make source-docs-check` checks all tracked Go documentation without compilation.
It is required by `make check`, `make check-full` and `make ci`; ordinary shared
CI selects the same strict check and cannot pass its required result if it fails.
Separately selected diagnostic and release qualifications keep their reviewed
source identities and scope. Semantic review and native docs remain necessary.
