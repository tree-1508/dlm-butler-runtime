# DLD-051 compliance and release gates

This file defines the release boundary for the DLD Butler Runtime used by Doujin Library Manager.
It applies to the `dld-051-v15.30.0` qualification line.

## Runtime topology

The Windows runtime payload is exactly one file: `butler.exe`.
`7z.dll`, `c7zip.dll`, and other native sidecars are not part of the DLD runtime payload.
Compliance material (SBOMs, notices, hashes, vulnerability evidence, and corresponding source) is distributed as separate release evidence and is not loaded as a runtime sidecar.

## Exact dependency truth

The distributed executable is the dependency source of truth. The qualification workflow reads the exact built executable with `go version -m` and records `EXACT_BINARY_MODULES.tsv`.
A dependency present only in `go.mod` is not treated as linked unless it is present in the binary metadata.

The current dependency-minimized clean runtime candidate links only `golang.org/x/sys v0.47.0` as a third-party Go module. Earlier full-butler qualification candidates that linked `github.com/itchio/dmcunrar-go` and other butler dependencies are historical and do not describe the current clean runtime topology.

## License closure

Every distributed DLD Windows executable must be accompanied, in the same release evidence set, by:

- `THIRD_PARTY_NOTICES.txt`, generated from license/copying/notice files for every module embedded in the exact executable;
- `LICENSE_MANIFEST.json` and the collected raw exact-binary license files;
- `dld-051-corresponding-source.zip`, containing the exact DLD runtime source plus vendored Go dependency sources needed to rebuild the executable;
- `SHA256SUMS.txt`, binding the executable and corresponding-source bundle.

The workflow fails if license evidence cannot be found for any exact linked module. No license may be silently classified or omitted.

The corresponding-source bundle is retained for reproducibility, auditability, and to satisfy any source-distribution obligations that apply to a future exact candidate. It must not be described as a GPL-specific requirement unless the exact candidate actually contains GPL-family code.

## SBOM closure

Every candidate build produces both:

1. `butler-exact-binary.cdx.json` from the exact `butler.exe`; and
2. `butler-windows-amd64-source.cdx.json` for the Windows/amd64 application under the same Windows build environment.

CycloneDX generation is pinned to a recorded tool version in build provenance.

## Vulnerability closure

The exact executable is scanned with a pinned `govulncheck` version in binary mode. The raw JSON report and finding count are retained as release evidence.
A non-zero finding count blocks DLD-051 closure until every finding has an explicit written disposition and any required repair has been rebuilt and rescanned.

## GitHub branch protection gate

Before an unsigned bootstrap pre-release is published, server-side protection for `dld-051-v15.30.0` must be enabled and verified. At minimum:

- force/non-fast-forward updates blocked;
- branch deletion blocked;
- changes merged through a pull request;
- required status checks include the DLD Windows unsigned trusted build, both DLD CI matrix jobs, and CodeQL analysis;
- required checks must be current for the exact release commit;
- high-or-higher CodeQL alerts in changed code block merge;
- no bypass actor may skip the release gate for a production signing request.

The project currently has one maintainer, so the policy must not invent an independent reviewer or configure an approval rule that makes truthful maintenance impossible. If SignPath Foundation requires an additional independent reviewer, signing remains NO-GO until that requirement is genuinely met.

## Release/signing sequence

1. exe-only feasibility
2. exact topology/license/SBOM/vulnerability closure
3. branch protection and policy closure
4. unsigned bootstrap GitHub pre-release
5. SignPath Foundation standard application
6. Foundation acceptance
7. SignPath Trusted Build + Origin Verification + signing
8. fresh DLD-051 Security/Windows qualification of the signed artifact

No provider credential, provider login, owned-library access, or real-provider traffic is authorized by this release process.
DLD-052 is outside this gate.
