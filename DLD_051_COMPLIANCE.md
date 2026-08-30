# DLD-051 compliance and release gates

This file defines the release boundary for the DLD Butler Runtime used by Doujin Library Manager.
It applies to the `dld-051-v15.30.0` qualification line and the `dld_runtime` build mode.

## Controlling Security decision

DLD Issue #150 Security comment `5465627296` is the controlling license-closure decision for the current line.
It establishes a hard release/signing gate: the exact Windows `butler.exe` must have zero unresolved linked component/module license entries before final DLD-051 license closure, unsigned bootstrap publication, or SignPath Foundation application.

Acceptable license evidence for an exact linked module is limited to:

1. exact-version/revision `LICENSE`, `LICENCE`, `COPYING`, `NOTICE`, or equivalent governing terms in the module/repository;
2. durable authoritative public developer/module metadata that names the governing license and is tied to the exact source revision; or
3. documented fork/derivation provenance where the exact imported source retains an OSI-approved license and all material later additions/modifications are mapped to compatible terms.

Organization-level assumptions, scanner guesses, package-index heuristics, or private correspondence do not close the gate.
If acceptable evidence cannot be established, the default remediation is removal of the unresolved code from the exact distributed binary.

## Runtime topology

The Windows runtime payload is exactly one file: `butler.exe`.
`7z.dll`, `c7zip.dll`, and other native sidecars are not part of the DLD runtime payload.
Compliance material (SBOMs, notices, hashes, vulnerability evidence, and corresponding source) is distributed as separate release evidence and is not loaded as a runtime sidecar.

The DLD qualification build uses `-tags dld_runtime`. Builds without that tag preserve the broader upstream-compatible command/router surface; they are not the DLD release artifact.
See `DLD_051_RUNTIME_SURFACE.md` for the release surface contract.

## Exact dependency truth

The distributed executable is the dependency source of truth. The qualification workflow reads the exact built executable with `go version -m` and records `EXACT_BINARY_MODULES.tsv`.
A dependency present only in `go.mod` is not treated as linked unless it is present in the binary metadata.

## License closure

The exact DLD Windows executable may link code carrying obligations beyond the main repository MIT license, including GPL-family code through `github.com/itchio/dmcunrar-go`. Therefore an executable release must not be treated as MIT-only merely because the upstream butler repository is MIT-licensed.

For every exact linked module the compliance job first searches the exact module source for governing license/notice files. Where Security #150 Q3(B) allows authoritative public metadata instead, the evidence must be explicitly keyed by exact `module@version` in `.github/compliance/dld051-curated-license-evidence.json`. No organization-wide or version-floating rule is permitted.

Every distributed DLD Windows executable must be accompanied, in the same release, by:

- `THIRD_PARTY_NOTICES.txt`, generated for the exact linked module set;
- `LICENSE_MANIFEST.json`, including exact module/version/checksum and evidence provenance;
- collected raw license files plus `CURATED_LICENSE_EVIDENCE.json` when curated exact-revision evidence is used;
- `dld-051-corresponding-source.zip`, containing the exact repository source plus vendored Go dependency sources needed to rebuild the executable;
- `SHA256SUMS.txt`, binding the executable and corresponding-source bundle.

The workflow records `LICENSE_GATE_STATUS.txt`. Any unresolved exact linked module sets this status to `HOLD`; the final combined gate then fails closed after all other evidence has been generated.

## Current exact-license reduction work

The initial exact exe-only candidate contained 87 linked Go modules and six modules without discoverable root/shallow license files.
Security-authorized exact-revision public metadata closes four of those candidates at policy level when the exact versions remain unchanged:

- `github.com/itchio/arkive@v0.0.0-20260428180635-32e8e9c72151` — exact README declares BSD licensing and Go-derived provenance;
- `github.com/itchio/hades@v0.0.0-20260711210423-80ab837c55cd` — exact README declares MIT;
- `github.com/itchio/intact@v0.0.0-20200301161822-f8c4a3336c2a` — exact README declares MIT licensing;
- `github.com/itchio/kompress@v0.0.0-20200301155538-5c2eecce9e51` — exact README declares BSD licensing and Go-derived provenance.

A GitHub-hosted Windows qualification of the current narrow-runtime PR reduced the exact linked module count from 87 to 85. Its remaining unresolved license entries were `github.com/itchio/hush@v0.0.0-20260710191509-8882c242cb2b` and `github.com/itchio/screw@v0.0.0-20260221011136-e674b460b040`.
This is evidence of the current work state, not a release approval. The release branch must be rebuilt after merge and must reach zero unresolved entries before this gate is closed.

## SBOM closure

Every candidate build produces both:

1. `butler-exact-binary.cdx.json` from the exact `butler.exe`; and
2. `butler-windows-amd64-source.cdx.json` for the source/application dependency view.

The exact-binary SBOM is authoritative for the shipped binary. CycloneDX generation is pinned to a recorded tool version in build provenance.

## Vulnerability closure

The exact executable is scanned with a pinned `govulncheck` version in binary mode. The raw JSON report and finding count are retained as release evidence.
The workflow records `VULNERABILITY_GATE_STATUS.txt`. A non-zero finding count remains `HOLD` until every finding has an explicit written disposition and any required repair has been rebuilt and rescanned.

The narrow-runtime PR qualification produced zero `govulncheck` finding records. This means only that the pinned scanner reported no reachable Go vulnerability findings in that exact candidate; it is not a general claim that the program or all third-party software is vulnerability-free.

## Governance and workflow controls

- GitHub-hosted runners only for the DLD Windows qualification/signing path.
- Default workflow permissions are read-only/minimum privilege; CodeQL receives only the additional `security-events: write` permission required for analysis upload.
- Checkout credentials are not persisted.
- Third-party Actions are pinned to immutable commit SHAs.
- Go/Node/compliance tool versions are pinned and recorded.
- `CODEOWNERS` covers workflows, compliance scripts/evidence policy, release build logic, DLD runtime surface files, signing policy, and dependency manifests.
- No `pull_request_target` execution of untrusted code is permitted.

## GitHub branch protection gate

Before an unsigned bootstrap pre-release is published, server-side protection for `dld-051-v15.30.0` must be enabled and verified. At minimum:

- force pushes disabled;
- branch deletion disabled;
- required status checks include the DLD Windows unsigned trusted build and both DLD CI matrix jobs;
- CodeQL analysis is required when GitHub exposes it as an eligible status check;
- required checks must be current for the exact release commit;
- administrators/maintainers do not bypass the release gate for a production signing request.

The project currently has one maintainer, so the policy must not invent an independent reviewer or configure an approval rule that makes truthful maintenance impossible. If SignPath Foundation requires an additional independent reviewer, signing remains NO-GO until that requirement is genuinely met.

Server-side protection is not considered complete merely because this document and `CODEOWNERS` exist; the repository setting itself must be read back and verified.

## Release/signing sequence

1. exe-only feasibility;
2. exact topology/license/SBOM/vulnerability closure;
3. branch protection and governance closure;
4. unsigned bootstrap GitHub pre-release with the exact warning `UNSIGNED BOOTSTRAP / PRE-SIGNPATH — NOT DLD-051 PASS — NOT FOR DLM INTEGRATION`;
5. SignPath Foundation standard application;
6. Foundation acceptance;
7. GitHub Trusted Build + Origin Verification + signing;
8. fresh DLD-051 Security/Windows qualification of the signed artifact.

No provider credential, provider login, owned-library access, real-provider traffic, or real provider payload acquisition is authorized by this release process.
DLD-052 is outside this gate and remains NOT AUTHORIZED until a genuine DLD-051 PASS plus explicit Owner approval.
