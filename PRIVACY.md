# Privacy Policy

DLM Butler Runtime is a command-line/provider runtime. It will not transfer information to networked systems unless specifically requested by the user or by the software/person operating it for an explicit provider operation.

When later used by Doujin Library Manager in an authorized integrated mode:
- provider authentication is designed around Authorization Code + PKCE;
- DLM must not import provider passwords, browser sessions, or API keys as a fallback;
- provider credential state may be persisted by the upstream-derived runtime in a dedicated local ProviderState database;
- production secret/debug logging is restricted by the DLM integration contract;
- provider operation traffic remains subject to DLM's policy gates and explicit authorization.

The current unsigned bootstrap is not integrated or auto-installed by DLM and does not authorize provider credentials, real OAuth/login, owned-library access, real provider traffic, or DLD-052 real-provider execution.

This project is not an official itch.io client and does not imply itch.io endorsement.
