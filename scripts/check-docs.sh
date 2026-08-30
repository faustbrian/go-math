#!/bin/sh
set -eu

required='README.md CHANGELOG.md SECURITY.md docs/README.md docs/specification-decisions.md docs/numeric-model.md docs/precision.md docs/conditions.md docs/serialization.md docs/security.md docs/performance.md docs/benchmark-baseline.md docs/migration.md docs/cookbook.md docs/compatibility.md docs/faq.md docs/troubleshooting.md docs/verification.md specification/README.md specification/manifest.tsv specification/decisions.json specification/conformance.json specification/decision-history.json specification/monitoring.json specification/maintained-peers.json'
for file in $required; do
	test -s "$file" || { printf 'missing documentation: %s\n' "$file" >&2; exit 1; }
done
go test . -run '^Example' -count=1
