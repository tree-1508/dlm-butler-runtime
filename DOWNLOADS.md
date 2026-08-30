# Downloads

## DLD-051 unsigned bootstrap / pre-SignPath

Current public Windows bootstrap:

**[dld-051-v15.30.0-unsigned-bootstrap.1](https://github.com/tree-1508/dlm-butler-runtime/releases/tag/dld-051-v15.30.0-unsigned-bootstrap.1)**

`UNSIGNED BOOTSTRAP / PRE-SIGNPATH — NOT DLD-051 PASS — NOT FOR DLM INTEGRATION`

This is an immutable GitHub **pre-release** containing the exact one-file Windows runtime candidate and its compliance/provenance evidence. It is intentionally unsigned and exists as the released form to be evaluated for the SignPath Foundation Open Source Code Signing program.

It is not a DLD-051 PASS artifact and must not be integrated, auto-installed, or distributed by Doujin Library Manager.

## Functionality

DLM Butler Runtime is a DLM-maintained, hardened, visibly modified fork of `itchio/butler`. The DLD-051 candidate is an optional runtime building block for future explicitly authorized provider operations in Doujin Library Manager.

The current distributed runtime payload is exactly one file, `butler.exe`. The bootstrap publication itself does not authorize or perform provider login, owned-library access, provider payload acquisition, or DLD-052 execution.

## Exact release identity

- Repository: <https://github.com/tree-1508/dlm-butler-runtime>
- Release branch: `dld-051-v15.30.0`
- Exact commit: `616652180ebacedb44b70150b4f995c5f7d9be67`
- Release tag: `dld-051-v15.30.0-unsigned-bootstrap.1`
- Exact `butler.exe` SHA-256: `d18440209557268f1ec1377260f39d4cabe54ade60cea92649266f0ddc685113`
- Corresponding-source SHA-256: `4229b7fde6bea66ed45d991a65941f64d16c136934eb0a668cfacd90a7f6cc89`
- GitHub Actions evidence artifact ID: `9725570125`
- Evidence artifact/release ZIP SHA-256: `170970fd4ec8ec2aeea6fc51d23640f2f0c1d3a3027b25d234c2b1f5c91af138`
- Exact third-party Go module: `golang.org/x/sys v0.47.0`
- License gate: PASS
- Vulnerability gate: PASS
- Binary-mode `govulncheck` findings: 0

The release includes `butler.exe`, SHA-256 manifests, corresponding source, third-party notices, a license manifest, exact-binary and Windows-source CycloneDX SBOMs, vulnerability evidence, and provenance.

## Code signing policy

Free code signing provided by SignPath.io, certificate by SignPath Foundation.

See the **[Code signing policy](CODE_SIGNING_POLICY.md)**.

No binary is represented as SignPath-signed before Foundation acceptance and successful completion of the approved Trusted Build and Origin Verification signing path.

## Privacy

See the **[Privacy Policy](PRIVACY.md)**.

This program will not transfer information to other networked systems unless specifically requested by the user or by the software/person operating it for an explicit provider operation.

## License

Project source: MIT.

The exact distributed executable contains only `golang.org/x/sys v0.47.0` as a third-party Go module in addition to project and Go standard-library code. Exact license evidence is included with the release.
