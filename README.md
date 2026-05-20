# arcmesh-pm

The package manager for MCP servers.

`apm` installs and manages MCP servers for AI clients like Claude Desktop, VS Code, Cursor, and Windsurf. It uses the ArcMesh curated registry as its primary source, with automatic fallback to the official MCP Registry for servers not yet in ArcMesh.

## Installation

**macOS / Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/arcmesh-labs/arcmesh-pm/master/install.sh | sh
```
No admin rights required. Installs to `~/.local/bin` and updates your shell PATH automatically.

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/arcmesh-labs/arcmesh-pm/master/install.ps1 | iex
```
After installation, restart your terminal for PATH changes to take effect.

## Usage

```bash
# Search for servers
apm search github
apm search database

# Install a server
apm install github
apm install notion --client cursor

# Manage installed servers
apm status
apm uninstall github
apm set-env github GITHUB_PERSONAL_ACCESS_TOKEN

# Inspect your setup
apm clients
apm doctor
apm config path
apm config edit
```

## Supported clients

| Client | Config file |
|---|---|
| Claude Desktop | `claude_desktop_config.json` |
| VS Code | `mcp.json` |
| Cursor | `mcp.json` |
| Windsurf | `mcp_config.json` |

## Registry

ArcMesh maintains a curated registry of verified MCP servers. When a server is not found there, `apm` automatically searches the [official MCP Registry](https://registry.modelcontextprotocol.io) as a fallback. Results are clearly labelled `[arcmesh]` or `[official]`.

## License

MIT