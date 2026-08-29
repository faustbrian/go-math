# Changelog

## Unreleased

### Changed

- Replace the repository-local verification implementation with the pinned
  `go-library-tools` v1.0.4 CLI and reusable workflow while preserving package
  policy and content-addressed verification evidence.
- Make architecture verification portable to CI runners without `ripgrep`.
- Remove the sibling-checkout compatibility command; reverse consumers now
  verify their own compatibility without requiring a monorepo layout.

### Documentation

- Replace archived monorepo links and completed execution artifacts with a
  standalone, human-oriented documentation structure.

## 1.0.0 - 2026-08-25

### Changed

- Exclude intentional nested modules from root local-proxy archives so local,
  bootstrap, CI, and public module checksums describe the same source
  boundary.

- Track the pinned documentation-tool lockfile so clean CI checkouts install
  the exact validated cspell dependency.

- Reconcile standalone dependency checksums against deterministic current
  module archives so CI, local verification, and release consumers resolve
  identical content.

- Harden standalone documentation validation with deterministic spelling and
  link checks, package-specific documentation gates, and repository-local
  contributor guidance.

### Changed

- Publish the module from its standalone `github.com/faustbrian/go-math` identity while preserving its documented API and behavior.

### Documentation

- Link the package README to package-owned documentation.

- Add immutable arbitrary-precision integer, rational, decimal, and binary
  float families.
- Add shared limits, rounding, conditions, deterministic codecs, numeric laws,
  conformance vectors, fuzzing, mutation checks, and benchmarks.
- Add the MIT license for source distribution and dependency review.
- Use the canonical `math:` prefix for exported sentinel error messages after
  migration from the former standalone repository.
- Harden decimal exponent, scaling, normalization, and rounding boundaries so
  malformed limits and extreme intermediate values fail without runaway work.
- Reject malformed deterministic binary frames at exact payload, header, sign,
  length, and exponent boundaries.
- Harden integer parsing, arithmetic preflights, random range limits, and root
  search boundaries against overflow and runaway work.
- Harden rational construction, powers, decimal expansion, parsing, and
  rounding at exact resource boundaries.
- Reject decimal digit input as soon as its configured budget is exhausted and
  reject repeated separators without attacker-sized allocation.
