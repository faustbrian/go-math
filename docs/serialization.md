# Conversions and serialization

[Documentation index](README.md)

Integer, rational and decimal text/JSON preserve numeric values when decoded
under sufficient limits; they do not preserve every decimal coefficient/exponent
representation. JSON deliberately emits strings to prevent decoder precision
loss. Binary-float text requires an explicit decode precision/rounding policy;
re-parsing with directed rounding can change the rendered value. Binary codecs
are deterministic,
versioned, and preserve decimal representation and float precision/rounding.
Mutable `math/big` values are copied on input and output.

Narrowing conversions are explicit and fail unless exact. Locale-aware output
and SQL policy belong in consumer adapters because locale and database numeric
contracts are not numeric identity.

