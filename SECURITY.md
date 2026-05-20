# Security Policy

## Reporting a vulnerability

Please do not report security vulnerabilities through public GitHub issues.

Instead, open a [private security advisory](https://github.com/arcmesh-labs/arcmesh-pm/security/advisories/new) on GitHub. This lets us discuss and fix the issue before it is publicly disclosed.

Include as much detail as possible:
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

We will acknowledge your report within 48 hours and aim to release a fix within 14 days for confirmed vulnerabilities.

## Scope

Issues we consider in scope:
- Token handling and storage in client config files
- Registry manifest fetching and validation
- Install runner security (pip, npx, uvx)

Issues out of scope:
- Vulnerabilities in third-party MCP servers listed in the registry
- Issues requiring physical access to the machine