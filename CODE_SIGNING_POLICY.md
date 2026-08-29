# Code signing policy

Free code signing provided by SignPath.io, certificate by SignPath Foundation.

## Project
DLM Butler Runtime is a visible modified-upstream fork of `itchio/butler` used only as an optional
provider runtime for Doujin Library Manager. It is not an official itch.io client and does not imply
itch.io endorsement.

## Roles
- Authors / committers: `@tree-1508`
- Reviewers for changes from non-committers: `@tree-1508`
- Signing approver: `@tree-1508`

These roles describe the current one-maintainer project truthfully. No shared or fictitious identity is used.
If SignPath requires an additional independent human for this project, release signing remains NO-GO until satisfied.

## Build and signing controls
- GitHub-hosted runners only for jobs leading to an OSS signing request.
- SignPath Trusted Build System and Origin Verification are mandatory for release signing.
- The unsigned artifact and its compliance evidence must be stored as a GitHub workflow artifact before submission.
- Release signing must be bound to the approved repository/branch/commit/build.
- Server-side branch protection must be verified before the unsigned bootstrap pre-release and remain enabled through signing.
- Every production signing request requires manual approval.
- No `pull_request_target` execution of untrusted contributor code.
- Default workflow token permissions remain read-only/minimum privilege.
- Signing/build workflows, compliance scripts, and `.signpath/` policy files are sensitive review paths.
- Third-party Actions are pinned to immutable commit SHAs.
- Compliance tool versions are pinned and recorded in provenance.
- No unreviewed `LATEST` runtime/native dependency channel is permitted.
- No signing key or exportable certificate material is stored in GitHub.

## License and release evidence
The runtime payload remains exe-only. SBOMs, third-party notices, hashes, vulnerability evidence, and corresponding source are separate release evidence.
Because the exact executable links GPL-family code, the generated corresponding-source bundle is mandatory for every executable distribution. See `DLD_051_COMPLIANCE.md`.

## Required sequence
1. unsigned exe-only qualification;
2. exact license/SBOM/vulnerability closure;
3. branch-protection/policy closure;
4. unsigned bootstrap pre-release;
5. SignPath Foundation standard application and acceptance;
6. Trusted Build + Origin Verification + signing;
7. fresh Security/Windows qualification of the signed artifact.

Signing must not be attempted before Foundation acceptance.

## Privacy
See `PRIVACY.md`.
