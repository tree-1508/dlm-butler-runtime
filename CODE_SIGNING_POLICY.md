# Code signing policy

Free code signing provided by SignPath.io, certificate by SignPath Foundation.

## Project

DLM Butler Runtime is a DLM-maintained, visibly modified GitHub fork of `itchio/butler` used only as an optional runtime for Doujin Library Manager.

- Repository: <https://github.com/tree-1508/dlm-butler-runtime>
- Downloads: [DOWNLOADS.md](DOWNLOADS.md)
- Privacy policy: [PRIVACY.md](PRIVACY.md)
- Upstream: <https://github.com/itchio/butler>

This project is not official itch.io software and does not imply itch.io endorsement.

The qualified Windows signing target is built from this fork's own source and build scripts. The exact candidate is a one-file `butler.exe` runtime. It does not distribute `7z.dll`, `c7zip.dll`, or other native runtime sidecars.

## Roles

The current project is maintained by one person, so the roles are truthfully assigned as follows:

- Authors / committers: `@tree-1508`
- Reviewers for changes from non-committers: `@tree-1508`
- Signing approver: `@tree-1508`

Every production signing request requires manual approval.

If SignPath Foundation requires another independent human for a particular responsibility, release signing remains blocked until that requirement is satisfied.

## Source, review, and build controls

- GitHub-hosted runners only for jobs leading to an Open Source signing request.
- SignPath Trusted Build System verification is mandatory for release signing.
- SignPath Origin Verification is mandatory for release signing.
- Signing is restricted to the approved public repository, protected release branch, exact commit, and build.
- The unsigned artifact must exist as a GitHub workflow artifact before submission for signing.
- Release signing uses the project's protected release source/commit; arbitrary runtime substitution and `LATEST` acquisition are prohibited.
- Third-party Actions used in the signing path must be reviewed and pinned to immutable full commit SHAs.
- Default workflow permissions remain read-only/minimum privilege.
- Signing/build workflows, compliance scripts, and `.signpath/` policy files are sensitive review paths.
- No signing private key or exportable Foundation certificate material is stored in GitHub.
- Signing credentials must not appear in source, logs, workflow artifacts, release assets, or Issues.

## Modified-upstream fork boundary

This project visibly uses GitHub's fork relationship to `itchio/butler`. The project signs only the modified runtime that it builds and maintains from its own public fork.

The DLD-051 signing candidate is dependency-minimized and differs from the historical upstream/full-butler package topology. Qualification and license decisions are made from the exact built candidate, not from unrelated optional upstream components.

## License and release evidence

The current exact distributed candidate:
- is one file: `butler.exe`;
- contains `golang.org/x/sys v0.47.0` as its only third-party Go module in addition to project and Go standard-library code;
- has exact-binary and Windows-source CycloneDX SBOMs;
- has `THIRD_PARTY_NOTICES.txt` and `LICENSE_MANIFEST.json`;
- has corresponding source and SHA-256 provenance;
- has license gate PASS;
- has vulnerability gate PASS;
- has binary-mode `govulncheck` finding count 0.

The evidence is published on the project's download/release page for the exact candidate.

## Privacy policy

This program will not transfer information to other networked systems unless specifically requested by the user or by the software/person operating it for an explicit provider operation.

See [PRIVACY.md](PRIVACY.md).

## Signing sequence

1. qualify the unsigned exe-only candidate;
2. close exact license/SBOM/vulnerability/governance gates;
3. publish the immutable unsigned bootstrap pre-release;
4. submit the standard SignPath Foundation Open Source Code Signing application;
5. wait for Foundation acceptance;
6. only after acceptance, configure the approved Trusted Build / Origin Verification / signing path;
7. perform fresh qualification of the first signed candidate.

No production SignPath GitHub App/token/signing request is started before Foundation acceptance.
