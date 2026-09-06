# Troubleshooting

`ErrInvalidSyntax` means the selected strict grammar rejected the input; enable
sign, exponent, underscore, whitespace, or leading-zero options deliberately.
`ErrLimitExceeded` means work was rejected before unsafe allocation; lower the
input size or raise a reviewed limit. `ErrConversion` from exact decimal
division means the rational expansion repeats. A trapped-condition error still
returns the decimal result and its conditions for diagnostics.

Classification uses `errors.Is`, never string matching. Invalid limits and
unsupported parser bases match `ErrInvalidArgument`; within valid options,
malformed input matches only `ErrInvalidSyntax`, while crossed raw, digit,
exponent, or binary-payload bounds match only `ErrLimitExceeded`. Diagnostics
are fixed, bounded, and never include rejected input or nested parser text.

`integer.Random` reader failures match `ErrRandomSource` and retain the exact
reader cause for ordinary `errors.Is`/`errors.As` traversal. Their public string
is always `math: random source failed`, so cause text is intentionally absent.

A nil context matches `ErrInvalidArgument`. Pre-cancellation returns the exact
`context.Canceled` sentinel, and an expired deadline returns the exact
`context.DeadlineExceeded` sentinel. Pure parsers and serializers do not take a
context; context-aware arithmetic is synchronous.

`Limits.MaxDiagnosticBytes` remains positive-validation input for source
compatibility but is a deprecated reserved no-op. It does not truncate stable
category messages.
