# DLM Fork Notice

This repository is a security-hardened, publicly auditable fork of `itchio/butler` for optional use by
Doujin Library Manager (DLM).

- Upstream project: https://github.com/itchio/butler
- This fork is not official itch.io software and does not imply itch.io endorsement.
- DLD hardening removes the gops diagnostic listener from the integrated runtime candidate.
- DLM integrated use is stdio JSON-RPC only; no DLM-approved TCP/HTTP JSON-RPC transport.
- Direct DLM `LibraryRoot` writes are outside this runtime's approved integration contract.

See `CODE_SIGNING_POLICY.md`, `PRIVACY.md`, and `SECURITY.md`.
