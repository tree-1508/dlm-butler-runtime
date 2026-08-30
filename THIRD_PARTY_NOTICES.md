# Third-party notices

The authoritative third-party notice set for a DLD Windows runtime is generated from the **exact built `butler.exe`**, not from the full `go.mod` graph.

For every qualification/release candidate, `.github/scripts/dld051-compliance.ps1` extracts the linked module set using `go version -m`, downloads the matching exact module sources, and produces:

- `THIRD_PARTY_NOTICES.txt` — notice/license evidence for the exact linked module set;
- `LICENSE_MANIFEST.json` — exact module/version/checksum to license-evidence mapping;
- `licenses/` — raw collected governing files;
- `CURATED_LICENSE_EVIDENCE.json` — exact-revision public metadata used only under Security #150 comment `5465627296` Q3(B), when applicable;
- `EXACT_BINARY_MODULES.tsv` — exact linked module/version/checksum inventory;
- `LICENSE_GATE_STATUS.txt` and, when unresolved entries exist, `MISSING_LICENSE_EVIDENCE.txt`.

The normal evidence path is an exact module `LICENSE`, `LICENCE`, `COPYING`, `NOTICE`, or `COPYRIGHT` file. Curated metadata is permitted only when it is authoritative, public, durable, and keyed to the exact module revision; it must never infer a license merely from repository organization, package-index metadata, or scanner heuristics.

The workflow deliberately generates SBOM, vulnerability, notice, hash, and corresponding-source evidence before enforcing the combined release gate. If any exact linked module still lacks acceptable license evidence, `LICENSE_GATE_STATUS.txt` remains `HOLD` and the release gate fails closed.

The exact executable may include GPL-family code through `github.com/itchio/dmcunrar-go`. Consequently every executable pre-release/release must also carry the generated `dld-051-corresponding-source.zip` and `SHA256SUMS.txt` in the same GitHub release.

Current narrow-runtime qualification has reduced the original six missing-license candidates to two unresolved exact entries (`github.com/itchio/hush` and `github.com/itchio/screw`). That state is **not** release approval; the release branch must reach zero unresolved exact entries before unsigned bootstrap publication or SignPath Foundation application.

See `DLD_051_COMPLIANCE.md` for the full gate and `DLD_051_RUNTIME_SURFACE.md` for the DLD build surface.
