# Third-party notices

The authoritative third-party notice set for a DLD Windows runtime is generated from the **exact built `butler.exe`**, not from the full repository or `go.mod` graph.

For every qualification/release candidate, `.github/scripts/dld051-compliance.ps1` extracts the linked module set using `go version -m`, downloads the matching module sources, collects their `LICENSE`, `LICENCE`, `COPYING`, `NOTICE`, and `COPYRIGHT` evidence, and produces:

- `THIRD_PARTY_NOTICES.txt` — concatenated notice/license evidence for the exact linked module set;
- `LICENSE_MANIFEST.json` — exact module/version to collected-license-file mapping;
- `EXACT_BINARY_LICENSES/` — raw collected license evidence included in the corresponding-source bundle;
- `EXACT_BINARY_MODULES.tsv` — exact linked module/version/checksum inventory;
- `dld-051-corresponding-source.zip` — exact DLD runtime source plus vendored dependency source used for the distributed executable;
- `SHA256SUMS.txt` — hashes binding the release/compliance evidence.

The workflow fails when any linked module lacks discoverable license evidence, so an incomplete notice set cannot silently pass the DLD-051 release gate.

For the current dependency-minimized clean runtime candidate, the exact binary inventory contains only `golang.org/x/sys v0.47.0` as a third-party Go module. Earlier full-butler qualification candidates that included `github.com/itchio/dmcunrar-go`, 7-Zip-related code paths, or other butler dependency families are historical and do not describe this clean runtime artifact.

Do not infer the notice/license obligations of a new candidate from this prose. The generated exact-binary inventory, license manifest, notices, SBOMs, and corresponding-source evidence for that candidate are controlling. See `DLD_051_COMPLIANCE.md` for the full gate.
