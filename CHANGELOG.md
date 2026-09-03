# Changelog

## Unreleased

### Changed

- Use the released `go-library-tools` v1.2.0 CLI and immutable merged
  workflow at `1f9629e5f27418600460b55a50a5b2fc81697fab` while preserving
  package policy and content-addressed verification evidence.
- Adopt the checksum-verified `go-library-tools` v1.3.0 CLI, schema-v2 cohesion
  metadata, and repository-local cohesion gate while retaining package-owned
  source and evidence.
- Make architecture verification portable to CI runners without `ripgrep` and
  restrict it to repository-owned Go sources.
- Remove the sibling-checkout compatibility command; reverse consumers now
  verify their own compatibility without requiring a monorepo layout.

### Documentation

- Govern the finite General Decimal Arithmetic semantic and corpus boundary
  through the [specification decision register](docs/specification-decisions.md),
  exact authority pins, executable conformance bindings, append-only decision
  history, and explicit maintained-peer gaps:
  - MATH-DEC-001 sha256:bf1c165c9f1dc38c6373deb5a104ab520a2392efd600b71e2cc75c5391303ff0
  - MATH-DEC-002 sha256:6b3149aea3a5fcb70b8bc58bc6618d1275209ee25489875c62e956b4a778d34d
  - MATH-DEC-003 sha256:58fa0f8a0d1d8af32d75659513da17ee3c44b02b086b6cc441f5e4e9c5890511
  - MATH-DEC-004 sha256:3c27045c7e9639a8fe7f8be4c3cad57195fca56521e5778a72cb88cec8321151
  - MATH-DEC-005 sha256:935afbb1f907f4da22212156dc6781175f3370941ac848156cbabd5e20301811
  - MATH-DEC-006 sha256:0bf5610e00555661d051678175d50f4a9d919cee4ce30b5da9838798ee7f67fc
  - MATH-DEC-007 sha256:0dc1ca93c10d999ee4038ce5d2dce1ce9575eb887c1960850b595e515925dba1
  - MATH-DEC-008 sha256:fddf149e8250739ca2f94cce735d77da5dcd6b642c27d76e629a53d7d3cb3c1e
  - MATH-DEC-009 sha256:e97d1a4031a30a04a3a520dd93275630f30114c147ececae10736b5539b3b4dc
  - MATH-DEC-010 sha256:bf3e572531b979dddd588313d43a85daa91b0af71ecc81d7319513704b197f1e
  - MATH-DEC-011 sha256:a1015e23bedf4dbf62a3ffbe2f0d43520ee6fa21abcfb4b12bf53935774f53c8
- Replace archived monorepo links and completed execution artifacts with a
  standalone, human-oriented documentation structure.
- Link the module to the immutable v1.3.0 Golib ecosystem guidance and expose
  the repository-local cohesion command in its contributor entry point.

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
