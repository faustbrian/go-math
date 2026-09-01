# Specification conformance matrix

The [source manifest](manifest.tsv) pins General Decimal Arithmetic 1.70, the
General Decimal Arithmetic Testcases sources, both 2.62 archives, and the
mutable release authority. The
[decision register](../docs/specification-decisions.md) is the canonical human
record. Machine bindings live in [decisions.json](decisions.json),
[conformance.json](conformance.json), and the append-only
[decision history](decision-history.json).

The arithmetic publication is version 1.70 dated 7 April 2009. The testcase
documentation page is version 2.44 dated 7 April 2009, while `testall.decTest`
and `testall0.decTest` inside the two archives declare corpus version 2.62 and
the archives are dated 19 April 2010. These independently versioned sources are
not collapsed into one synthetic release.

The package claims only the finite decimal behavior named below. It does not
implement NaN, infinity, signed-zero identity, or non-extended operand
pre-rounding and does not claim IEEE 754 or IEEE 854 conformance. Passing the
applicable official vectors is interoperability evidence for the named
operations, not broad standards certification.

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and
"OPTIONAL" in this document are to be interpreted as described in BCP 14
[RFC2119] [RFC8174] when, and only when, they appear in all capitals, as shown
here.

A source digest change MUST NOT silently change behavior. Review the changed
authority, decision records, compatibility, corpus accounting, peer evidence,
and executable evidence before updating a pin.

## Decision matrix

| Decision | Observable contract | Primary evidence | Independent status |
| --- | --- | --- | --- |
| MATH-DEC-001 | Finite decimal model with canonical numeric zero | value, parsing, text, JSON, and fuzz tests | no broad peer claim |
| MATH-DEC-002 | Strict parsing and separate numeric and representation equality | parser, representation, and fuzz tests | not assessed |
| MATH-DEC-003 | Separate exact and explicit context-rounded arithmetic | exact, context, division, and peer tests | maintained peer agreement on the recorded overlap |
| MATH-DEC-004 | Significant precision, adjusted exponent, conditions, and traps | context condition and trap tests | maintained peer agreement on shared Rounded and Inexact conditions |
| MATH-DEC-005 | Seven explicit rounding modes | tie tests and `rounding0.decTest` | official fixture agreement |
| MATH-DEC-006 | Classified exceptional-operation errors with retained conditions | trap, division, resource, and cancellation tests | not assessed |
| MATH-DEC-007 | Fractional scale maps to target exponent and quotient rounds once | quantize tests, `quantize0.decTest`, and peer tests | official fixture and maintained peer agreement |
| MATH-DEC-008 | Two pinned archives, LF normalization, 3,547 executed vectors, and 1,542 classified skips | six executed operation files and exact harness accounting | official fixture agreement |
| MATH-DEC-009 | Explicit resource limits and synchronous cancellation | resource, hostile operand, and fuzz tests | defensive package policy |
| MATH-DEC-010 | Canonical text, JSON strings, and versioned deterministic binary frames | text, JSON, binary, and decoder fuzz tests | package wire policy; peer gap |
| MATH-DEC-011 | 488 maintained-peer comparisons with explicit gaps | `TestDecimalDifferentialAgainstAPDAndShopspring` and `maintained-peers.json` | maintained peer agreement on the recorded overlap |

The GDA harness executes these exact records:

| File | Executed | Skipped |
| --- | ---: | ---: |
| `add.decTest` | 1,483 | 617 |
| `subtract.decTest` | 533 | 148 |
| `multiply.decTest` | 160 | 361 |
| `divide.decTest` | 357 | 274 |
| `quantize0.decTest` | 377 | 51 |
| `rounding0.decTest` | 637 | 91 |
| **Total** | **3,547** | **1,542** |

Run the offline structural check with `golib specification check`. Use
`golib specification check --online` to re-fetch all bounded authorities and
verify the reviewed content digests.

[RFC2119]: https://www.rfc-editor.org/rfc/rfc2119
[RFC8174]: https://www.rfc-editor.org/rfc/rfc8174
