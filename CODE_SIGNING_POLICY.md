# Code signing policy

Free code signing is intended to be provided by SignPath.io, certificate by SignPath Foundation, after the Foundation application is accepted.

## Project

DLM Butler Runtime is a visible modified-upstream fork of `itchio/butler` used only as an optional provider runtime for Doujin Library Manager. It is not an official itch.io client and does not imply itch.io endorsement.

The DLD release artifact is built with `-tags dld_runtime` and is intentionally narrower than a normal upstream-compatible butler build. See `DLD_051_RUNTIME_SURFACE.md`.

## Roles

- Authors / committers: `@tree-1508`
- Reviewers for changes from non-committers: `@tree-1508`
- Signing approver: `@tree-1508`

These roles describe the current one-maintainer project truthfully. No shared or fictitious identity is used.
If SignPath requires an additional independent human for this project, release signing remains NO-GO until satisfied.

## Build and signing controls

- GitHub-hosted runners only for jobs leading to an OSS signing request.
- SignPath Trusted Build System and Origin Verification are mandatory for release signing.
- The unsigned artifact and its exact compliance evidence must be stored as a GitHub workflow artifact before submission.
- Release signing must be bound to the approved repository/branch/commit/build.
- Server-side branch protection must be verified before the unsigned bootstrap pre-release and remain enabled through signing.
- Every production signing request requires manual approval.
- No `pull_request_target` execution of untrusted contributor code.
- Default workflow token permissions remain read-only/minimum privilege.
- Signing/build workflows, compliance scripts, exact-revision license evidence, DLD runtime surface files, and `.signpath/` policy files are sensitive review paths.
- Third-party Actions are pinned to immutable commit SHAs.
- Compliance tool versions are pinned and recorded in provenance.
- No unreviewed `LATEST` runtime/native dependency channel is permitted.
- No signing key or exportable certificate material is stored in GitHub.
- A pull-request merge-ref artifact is qualification evidence only. The artifact submitted for release/signing must be rebuilt from the exact protected release-branch commit so Origin Verification can bind the canonical branch/commit/build.

## License and release evidence

Security #150 comment `5465627296` requires zero unresolved exact linked license entries before unsigned bootstrap publication or SignPath Foundation application.

The runtime payload remains exe-only. SBOMs, third-party notices, hashes, vulnerability evidence, curated exact-revision public license metadata where permitted, and corresponding source are separate release evidence.
Because the exact executable may link GPL-family code, the generated corresponding-source bundle is mandatory for every executable distribution. See `DLD_051_COMPLIANCE.md`.

## Required sequence

1. unsigned exe-only qualification;
2. exact license/SBOM/vulnerability closure;
3. branch-protection/policy closure;
4. unsigned bootstrap pre-release;
5. SignPath Foundation standard application and acceptance;
6. Trusted Build + Origin Verification + signing;
7. fresh Security/Windows qualification of the signed artifact.

Signing must not be attempted before Foundation acceptance.

## Provider boundary

This code-signing/release work does not authorize provider credentials, provider login, owned-library access, real-provider traffic, or real provider payload acquisition.
It does not authorize DLD-052.

## Privacy

See `PRIVACY.md`.
