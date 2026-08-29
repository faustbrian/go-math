#!/bin/sh
set -eu

matches=$(find . -type f -name '*.go' ! -name '*_test.go' \
	-exec grep -nEH '"C"|unsafe\.' {} + || true)
if [ -n "$matches" ]; then
	printf '%s\n' "$matches"
	printf 'unsafe or cgo use is forbidden\n' >&2
	exit 1
fi
matches=$(find . -type f -name '*.go' ! -name '*_test.go' \
	-exec grep -nEH 'math/rand' {} + || true)
if [ -n "$matches" ]; then
	printf '%s\n' "$matches"
	printf 'ambient pseudo-randomness is forbidden\n' >&2
	exit 1
fi
matches=$(find . -type f -name '*.go' ! -name '*_test.go' \
	-exec grep -nEH 'go[[:space:]]+([[:alnum:]_]+\.)?[[:alnum:]_]+\(' {} + || true)
if [ -n "$matches" ]; then
	printf '%s\n' "$matches"
	printf 'production goroutines are forbidden in this synchronous library\n' >&2
	exit 1
fi
matches=$(find . -type f -name '*.go' ! -name '*_test.go' \
	! -path './mathtest/*' \
	-exec grep -nEH '^type [[:alnum:]_]+ interface' {} + || true)
if [ -n "$matches" ]; then
	printf '%s\n' "$matches"
	printf 'cross-family production interfaces are forbidden\n' >&2
	exit 1
fi
go list -deps ./... >/dev/null
