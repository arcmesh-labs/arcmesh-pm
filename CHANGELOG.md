# Changelog

## v0.1.1 — 2026-05-20

### Added
- `apm add` — register the MCP server in the current directory with your AI client config; auto-detects WSL distro and Python venv, wraps with `wsl.exe` on WSL
- `install.sh` — one-liner installer for Linux and macOS

### Fixed
- `install.ps1` — replaced `RuntimeInformation.OSArchitecture` with `$env:PROCESSOR_ARCHITECTURE` for compatibility with PowerShell 5.1

### Docs
- Updated Windows installation guide in README

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