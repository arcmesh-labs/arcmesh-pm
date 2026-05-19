# arcmesh-pm

The package manager for MCP servers.

`apm` installs and manages MCP servers for AI clients like Claude Desktop, VS Code, Cursor, and Windsurf. It uses the ArcMesh curated registry as its primary source, with automatic fallback to the official MCP Registry for servers not yet in ArcMesh.

## Installation

Download the latest binary for your platform from [Releases](https://github.com/arcmesh-labs/arcmesh-pm/releases) and place it somewhere in your PATH.

**macOS / Linux:**
```bash
curl -L https://github.com/arcmesh-labs/arcmesh-pm/releases/latest/download/apm-linux-amd64.tar.gz | tar xz
mv apm /usr/local/bin/
```

**Windows:** Download the `.zip` from Releases and add the binary to your PATH.

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/arcmesh-labs/arcmesh-pm/main/install.ps1 | iex
```

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