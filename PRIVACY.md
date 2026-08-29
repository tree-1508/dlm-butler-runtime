# Privacy Policy

DLM Butler Runtime is a command-line/provider runtime. It will not transfer information to networked
systems unless specifically requested by the user or by the software/person operating it for an explicit
provider operation.

When used by Doujin Library Manager in integrated mode:
- provider authentication is designed around Authorization Code + PKCE;
- DLM must not import provider passwords, browser sessions, or API keys as a fallback;
- provider credential state may be persisted by the upstream runtime in the dedicated local ProviderState database;
- production secret/debug logging is restricted by the DLM integration contract;
- provider operation traffic remains subject to DLM's policy gates.

This project is not an official itch.io client and does not imply itch.io endorsement.
