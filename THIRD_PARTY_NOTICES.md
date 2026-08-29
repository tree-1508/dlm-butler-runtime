# Third-party notices

The authoritative third-party notice set for a DLD Windows runtime is generated from the **exact built `butler.exe`**, not from the full `go.mod` graph.

For every qualification/release candidate, `.github/scripts/dld051-compliance.ps1` extracts the linked module set using `go version -m`, downloads the matching module sources, collects their `LICENSE`, `LICENCE`, `COPYING`, `NOTICE`, and `COPYRIGHT` evidence, and produces:

- `THIRD_PARTY_NOTICES.txt` — complete concatenated notice/license evidence for the exact linked module set;
- `LICENSE_MANIFEST.json` — module/version to collected-license-file mapping;
- `licenses/` — raw collected files;
- `EXACT_BINARY_MODULES.tsv` — exact linked module/version/checksum inventory.

The workflow fails when any linked module lacks discoverable license evidence, so an incomplete notice set cannot silently pass the DLD-051 release gate.

The exact executable includes GPL-family code through `github.com/itchio/dmcunrar-go`. Consequently every executable pre-release/release must also carry the generated `dld-051-corresponding-source.zip` and `SHA256SUMS.txt` in the same GitHub release. See `DLD_051_COMPLIANCE.md` for the full gate.
