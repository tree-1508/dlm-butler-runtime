# DLM Butler Runtime

> [!IMPORTANT]
> This repository is the DLM-maintained hardened, visibly modified fork of [`itchio/butler`](https://github.com/itchio/butler).
> It is intended only as an optional runtime for Doujin Library Manager (DLM). It is not official itch.io software and does not imply itch.io endorsement.

The current DLD-051 Windows signing candidate is a dependency-minimized, one-file `butler.exe` runtime built from the protected release branch and qualified with exact-binary/source SBOMs, license notices, vulnerability evidence, corresponding source, and SHA-256 provenance.

## Downloads

See **[Downloads and release provenance](DOWNLOADS.md)**.

The currently published bootstrap is intentionally unsigned and is labeled:

`UNSIGNED BOOTSTRAP / PRE-SIGNPATH — NOT DLD-051 PASS — NOT FOR DLM INTEGRATION`

It is public evidence for the SignPath Foundation application and must not be integrated or auto-installed by DLM.

## Code signing policy

Free code signing provided by SignPath.io, certificate by SignPath Foundation.

See the project **[Code signing policy](CODE_SIGNING_POLICY.md)** and **[Privacy Policy](PRIVACY.md)**.

Current project roles:
- Authors / committers: `@tree-1508`
- Reviewers for changes from non-committers: `@tree-1508`
- Signing approver: `@tree-1508`

Every production signing request requires manual approval. Open Source signing must use SignPath Trusted Build System verification and Origin Verification.

No release binary is represented as SignPath-signed unless and until the SignPath Foundation application is accepted and the approved signing path has completed.

## Current qualified release identity

- Release branch: `dld-051-v15.30.0`
- Exact candidate commit: `616652180ebacedb44b70150b4f995c5f7d9be67`
- Unsigned bootstrap tag: `dld-051-v15.30.0-unsigned-bootstrap.1`
- Exact `butler.exe` SHA-256: `d18440209557268f1ec1377260f39d4cabe54ade60cea92649266f0ddc685113`
- Exact corresponding-source SHA-256: `4229b7fde6bea66ed45d991a65941f64d16c136934eb0a668cfacd90a7f6cc89`
- Exact third-party Go module in the distributed executable: `golang.org/x/sys v0.47.0`
- License gate: PASS
- Vulnerability gate: PASS
- Binary-mode `govulncheck` findings: 0

## Build and provenance

The qualified Windows artifact is built on GitHub-hosted runners. Release signing is intended to use SignPath's GitHub Trusted Build System and Origin Verification so that the signed artifact is bound to the public repository, protected release branch, exact commit, and GitHub-hosted build.

The unsigned candidate and its compliance evidence are retained as GitHub Actions/release artifacts before any signing step.

## Privacy

This runtime does not transfer information to networked systems unless specifically requested by the user or by the software/person operating it for an explicit provider operation. See **[PRIVACY.md](PRIVACY.md)** for the full project privacy statement.

## Upstream

This repository is a GitHub fork of [`itchio/butler`](https://github.com/itchio/butler). Upstream documentation remains available at:
- <https://itch.io/docs/butler>
- <https://itchio.github.io/butler/>
- <https://itchio.github.io/butler/butlerd/>

DLM-specific hardening, qualification, signing policy, and release provenance are maintained in this fork.

## License

The project source is MIT licensed. See [`LICENSE`](LICENSE).

The exact distributed DLD-051 candidate is exe-only and contains only `golang.org/x/sys v0.47.0` as a third-party Go module in addition to project and Go standard-library code. Exact license evidence is published with each candidate.
