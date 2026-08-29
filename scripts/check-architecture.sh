#!/bin/sh
set -eu

production_go_files() {
	files=$(git ls-files --cached --others --exclude-standard -- '*.go')
	if [ -z "$files" ]; then
		return
	fi
	printf '%s\n' "$files" | while IFS= read -r file; do
		case "$file" in
		*_test.go | .golib-tooling/* | .verification/*) continue ;;
		esac
		printf '%s\n' "$file"
	done
}

scan_production_go() {
	pattern=$1
	excluded_prefix=${2-}
	files=$(production_go_files)
	if [ -z "$files" ]; then
		return
	fi
	printf '%s\n' "$files" | while IFS= read -r file; do
		if [ -n "$excluded_prefix" ]; then
			case "$file" in
			"$excluded_prefix"/*) continue ;;
			esac
		fi
		grep -nEH "$pattern" "$file" || true
	done
}

matches=$(scan_production_go '"C"|unsafe\.')
if [ -n "$matches" ]; then
	printf '%s\n' "$matches"
	printf 'unsafe or cgo use is forbidden\n' >&2
	exit 1
fi
matches=$(scan_production_go 'math/rand')
if [ -n "$matches" ]; then
	printf '%s\n' "$matches"
	printf 'ambient pseudo-randomness is forbidden\n' >&2
	exit 1
fi
matches=$(scan_production_go 'go[[:space:]]+([[:alnum:]_]+\.)?[[:alnum:]_]+\(')
if [ -n "$matches" ]; then
	printf '%s\n' "$matches"
	printf 'production goroutines are forbidden in this synchronous library\n' >&2
	exit 1
fi
matches=$(scan_production_go '^type [[:alnum:]_]+ interface' mathtest)
if [ -n "$matches" ]; then
	printf '%s\n' "$matches"
	printf 'cross-family production interfaces are forbidden\n' >&2
	exit 1
fi
go list -deps ./... >/dev/null
