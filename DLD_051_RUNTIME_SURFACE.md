# DLD-051 least-privilege runtime surface

The DLD Windows release artifact is compiled with the Go build tag `dld_runtime`.
This build tag exists to make the signed DLD runtime materially smaller in attack surface and dependency reachability than a normal upstream-compatible butler build.

This document describes the intended **DLD-051 qualification surface only**. It is not authorization to start DLD-052 or to perform provider operations during qualification.

## CLI surface

The `dld_runtime` build registers only the `daemon` command.
Broader butler CLI commands remain available only in builds without `dld_runtime`; they are not part of the DLD release artifact.

## RPC methods intentionally registered

The narrow daemon router currently registers only:

- `Version.Get`
- `Meta.Shutdown`
- `Profile.LoginWithOAuthCode`
- `Profile.List`
- `Profile.Forget`
- `Fetch.ProfileGames`
- `Game.FindUploads`
- `Install.GetUploads`
- `Install.PlanUpload`
- `Install.Queue`
- `Install.Perform`
- `Install.Cancel`

Registration means the code path is present for later separately authorized integration. DLD-051 build/security qualification must not exercise provider credentials, provider login, owned-library access, or real provider traffic.

## Explicitly absent from the DLD router

The DLD router deliberately does not register:

- password login;
- direct API-key login;
- launch/execute methods;
- uninstall;
- prerequisite installer execution;
- publishing/upload-to-provider endpoints;
- search/update endpoints not required by the narrow runtime contract;
- diagnostic/test RPC endpoints;
- the broad `EnsureAllRequests` registration that would make accidental expansion fail open.

A future method cannot become reachable in the DLD build merely because upstream adds it to the global message catalog; it must be explicitly registered in `cmd/daemon/endpoints_dld_runtime.go` and pass the DLD security/release gates.

## Transport qualification

The DLD release path is stdio JSON-RPC only. The build qualification invokes the synthetic stdio transport integration test and does not use provider credentials or real provider traffic.

Any listener/process-tree controls required by DLD-051 Security remain separate mandatory gates. This file does not weaken earlier Security decisions concerning gops/listener removal, process supervision, restart/orphan handling, ProviderState ACLs, OAuth transaction controls, logging/redaction, or Storage Root confinement.

## Packaging

The DLD Windows runtime payload is exactly `butler.exe`. `7z.dll` and `c7zip.dll` are not included in the DLD runtime payload.

## Dependency minimization

Security #150 comment `5465627296` authorizes dependency/path minimization within this public OSS runtime when needed to remove exact linked modules whose governing license cannot be established under the accepted evidence standard.

The exact binary, not the source import graph, decides whether a module remains part of the distributed component set. Every candidate therefore regenerates `EXACT_BINARY_MODULES.tsv`, SBOMs, notices, corresponding source, and vulnerability evidence.

## Change control

The following paths are sensitive and owned by `@tree-1508` through CODEOWNERS:

- `commands_dld_runtime.go`
- `cmd/daemon/endpoints_dld_runtime.go`
- `release/`
- `.github/workflows/`
- `.github/scripts/`
- `.github/compliance/`
- `.signpath/`

The release branch must be server-side protected before unsigned bootstrap publication.

## Authorization boundary

Nothing in this file authorizes:

- provider credentials;
- provider login;
- owned-library access;
- real provider traffic;
- real provider payload acquisition;
- DLD-052 implementation/execution;
- private DLM product source writes.

DLD-052 remains NOT AUTHORIZED until genuine DLD-051 PASS plus explicit Owner approval.
