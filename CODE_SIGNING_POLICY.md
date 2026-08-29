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
- The unsigned artifact must be stored as a GitHub workflow artifact before submission.
- Release signing must be bound to the approved repository/branch/commit/build.
- Every production signing request requires manual approval.
- No `pull_request_target` execution of untrusted contributor code.
- Default workflow token permissions remain read-only/minimum privilege.
- Signing/build workflows and `.signpath/` policy files are sensitive review paths.
- Third-party Actions are pinned to immutable commit SHAs.
- No unreviewed `LATEST` runtime/native dependency channel is permitted.
- No signing key or exportable certificate material is stored in GitHub.

## Privacy
See `PRIVACY.md`.
