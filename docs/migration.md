# Migration

Choose one family from the domain invariant, not the source representation.
Replace integer wrappers with `integer.Integer`, exact ratios with `Rational`,
money-like base-10 values with `Decimal`, and only explicitly binary numerical
algorithms with `bigfloat.Float`. Parse persisted values as strings, define
limits centrally, and make every former implicit rounding point an explicit
context or quantization call.

During migration, compare serialized fixtures and run both implementations at
consumer boundaries. Do not convert legacy decimal values through `float64`.

When upgrading from v1.0.0 to v1.1.0, replace direct error equality and string
matching with `errors.Is`. Reader failures from `integer.Random` can now be
classified with `ErrRandomSource`, while their original causes remain
traversable. Do not depend on the old reader-cause text appearing in the public
message.

The only accepted-input reduction is for hostile raw forms beyond the
documented Integer, Rational, and Decimal bounds. Review callers that used
unbounded surrounding whitespace or zero-padded exponents. Pinning v1.0.0
restores that behavior, but also restores ambiguous rejection categories and
unsafe reader-cause disclosure and is therefore a security and observability
regression rather than an equivalent rollback.
