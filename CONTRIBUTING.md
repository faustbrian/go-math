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

Do not add package-local workflows, permanent replacements, machine-specific
paths, bypass flags, broad mutation exclusions, or aggregate quality metrics
that hide a failing package.

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
`golib` v2.1.0 release declared in `.golib.yaml` and verify its checksums:

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
consumers for contract changes. Scheduled and release qualification retain
their full selected gates and release rehearsal. The existing native-security
and root-release qualification dispatches retain their separate source-built
verifiers; neither changes the ordinary public verifier identity. Report failing
or unavailable checks without claiming release readiness. CI/configuration-only
changes do not require a library release.

## Adding A Module

Follow [repository structure policy](AGENTS.md#repository-structure). New modules
require an explicit purpose, ownership boundary, dependency review, package
catalog entry, full quality gates, documentation, changelog, license, security
policy, compatibility plan, and release dry-run.
