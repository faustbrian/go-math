# Compatibility Policy

The repository contains one releasable module at its root and follows semantic
versioning. Root-module tags use `v<version>`, such as `v1.1.0`. This
repository has no nested Go modules or directory-prefixed module tags.

Before `v1`, minor releases MAY contain reviewed breaking changes, but every
break MUST be documented with migration guidance. Patch releases MUST remain
backward compatible. At and after `v1`, incompatible exported API or documented
behavior changes require a new major version.

Compatibility includes exported Go APIs, error classification, serialization,
protocol behavior, persistence schemas, environment variables, command output,
resource ownership, ordering, retry/idempotency semantics, and documented
defaults. A compile-compatible change can still be behaviorally breaking.

Specification-backed modules MUST NOT diverge from their declared standards.
Ambiguities require documented decisions and stable tests. Deprecated APIs
follow [`DEPRECATION.md`](DEPRECATION.md). Decimal interpretation and corpus
boundaries are governed by the
[specification decision register](docs/specification-decisions.md); changing a
resolved decision requires a compatibility review and a new retained digest.
