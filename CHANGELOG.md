# Changelog

## v0.1.0 — 2026-05-19

### Added
- `apm install` — install MCP servers from ArcMesh registry, with automatic fallback to the official MCP Registry
- `apm uninstall` — remove an MCP server from your client config
- `apm search` — search across ArcMesh and official MCP Registry; results clearly labelled `[arcmesh]` or `[official]`
- `apm list` — list all available servers in the ArcMesh registry
- `apm status` — show all installed servers across configured clients
- `apm doctor` — validate MCP client config health
- `apm clients` — list supported clients and their config status
- `apm config edit` — open client config in `$EDITOR`
- `apm config path` — print path to client config file
- `apm set-env` — update an environment variable for an installed server
- Support for Claude Desktop, VS Code, Cursor, and Windsurf
- WSL support for Windows-hosted client configs
- Cross-platform release builds via goreleaser (Linux, macOS, Windows)