# Compatibility

Decimal parsing, arithmetic, condition, quantize, serialization, corpus, and
peer-comparison boundaries are recorded in the
[specification decision register](specification-decisions.md). A changed
resolved decision is a compatibility event even when exported Go signatures do
not change.

The module targets Go 1.26.6 as both its minimum and CI toolchain. Public API
changes are checked against `api/baseline.txt`. Binary encodings carry a version
byte; unknown versions fail. Text and JSON forms are canonical strings.

`money` and `measurement` should import only the numeric family they need
and keep currency, unit, locale, and presentation policy in their own domains.
