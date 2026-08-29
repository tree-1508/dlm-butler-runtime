# Security Policy

This repository is the public OSS trust domain for the optional DLM Butler Runtime only.
It must not contain Doujin Library Manager proprietary/private application source, provider credentials,
user data, ProviderState databases, signing private keys, PFX/P12 files, token PINs, or recovery material.

For vulnerabilities, prefer GitHub private vulnerability reporting / Security Advisories when available.
Do not post credentials or exploit-sensitive user data in public Issues.

Production integration is fail-closed: invalid runtime identity, signature, provenance, listener/process
qualification, or policy state must prevent integrated provider use.
