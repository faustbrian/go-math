# General Decimal Arithmetic vectors

The non-`0` operation files (`add.decTest`, `subtract.decTest`,
`multiply.decTest`, and `divide.decTest`) come from `dectest.zip`. The
`0`-suffixed files plus `quantize0.decTest` and `rounding0.decTest` come from
`dectest0.zip`. Both archives declare testcase version 2.62 and were published
on 19 April 2010 by Mike Cowlishaw at
<https://speleotrove.com/decimal/>.

Upstream archive members use CRLF line endings. Repository copies are
deterministically normalized to LF; their checked-in hashes match the upstream
members after removing each carriage return. No other content transformation
is applied. The archive and authority digests are pinned in
[`../manifest.tsv`](../manifest.tsv), and the normalized file digests below are
enforced locally so conformance tests do not depend on the network.

The upstream notices describe the vectors as experimental and explain that
passing them does not by itself establish standards conformance. The test
harness executes the finite operations supported by this package and records
unsupported grammar or exceptional-value cases separately.

Each corpus member retains its upstream notice: copyright Mike Cowlishaw
(1981-2010), with parts copyright IBM Corporation (1981-2008), all rights
reserved, offered as-is. The testcase documentation footer says it is
reproduced with permission from IBM. Neither pinned archive contains a separate
license file, and the publication page does not state a broader redistribution
license; the repository MIT license does not replace these upstream notices.
Any corpus update therefore requires an explicit copyright and redistribution
review in addition to technical provenance review.

The harness directly executes the non-`0` add, subtract, multiply, and divide
files plus `quantize0.decTest` and `rounding0.decTest`. The four `0`-suffixed
arithmetic files remain provenance-checked corpus members but are not directly
executed. Exact executed and skipped counts are recorded in the
[conformance matrix](../README.md) and decision MATH-DEC-008.

To update the corpus, download both pinned HTTPS archives, verify their archive
digests from `../manifest.tsv`, extract only the ten named members, normalize
CRLF to LF by removing carriage returns, refresh `SHA256SUMS`, review harness
applicability and exact accounting under MATH-DEC-008, and run the provenance,
specification, and focused conformance gates. Do not silently replace a source,
transformation, notice, vector, or skip classification.

SHA-256 checksums:

```text
2a0e43703b9c5f0caf9aaf6077cf6cd7441f1d12c383dfc92f9459fdd5b4600e  add.decTest
8b3d3ba081aa1e80d5cd21b622d68d66305ed407d982e3fe77dcddf3b8df862b  add0.decTest
b71d0a18016c4434c80ad7d6aa0ade89b4d10bca148d1ca5299f81af96243c98  divide.decTest
923f64762157dd12de87d57fe190ae907e892e63a510f2f8b9b52db9663ac220  divide0.decTest
e4c3cdc1c2d0e08349ea9ceae094c94b1b942473c8ebc1e8a42f0e41e3fae282  multiply.decTest
ea7d578a2a4cd632eed564ef94d038e6a1f76be22475507bcab24a054454f9ff  multiply0.decTest
aad6ad6d6f20da15d6e0b7ed2d1c6abc45ca3b9976c33b0c306c0f423c04213e  quantize0.decTest
20e76b459d0d2881f6d3891089078e06319d43ac30485cad6bf850315164592b  rounding0.decTest
6e7b01cb70f844906c5d15609df32e86d14d8faf232b184e78eb96d65ee0b684  subtract.decTest
2e3efc8af9113a1e892312ae38cd63def7534059b8d34e28053795afc5cb344b  subtract0.decTest
```
