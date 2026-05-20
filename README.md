# arcmesh-pm

The package manager for MCP servers.

`apm` installs and manages MCP servers for AI clients like Claude Desktop, VS Code, Cursor, and Windsurf — with a single command and no infrastructure required.

---

## What is MCP?

MCP (Model Context Protocol) lets AI assistants connect to external tools — GitHub, Notion, databases, your own codebase, and more. Without a package manager, setting up an MCP server means: find the package, read the docs, install manually, write JSON config, find the right config file for your client. `apm` eliminates all of that.

---

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

Restart your terminal after installation for PATH changes to take effect.

---

## Quick start

```bash
# Find a server
apm search github

# Install it
apm install github

# Check what's installed
apm status
```

That's it. `apm install` fetches the manifest, runs the install, prompts for any required API tokens, and writes the config — then tells you to restart your client.

---

## Making your own repo AI-ready

`apm add` scaffolds a local MCP server for the current directory and registers it with your AI client. Your assistant gets `read_file`, `list_directory`, and `search_content` tools scoped to your repo.

```bash
cd my-project
apm add
```

This creates `.mcp/server.py` and `.mcp/config.json` in your project, and registers the server in your client config. Restart your client and the tools are available.

```bash
# Custom name
apm add --name my-project

# Target a specific client
apm add --client cursor
```

---

## Command reference

### Registry

| Command | Description |
|---|---|
| `apm search <query>` | Search for MCP servers across ArcMesh and the official MCP Registry |
| `apm list` | List all servers available in the ArcMesh registry |

```bash
apm search github
apm search database
apm list
```

---

### Install

| Command | Description |
|---|---|
| `apm install <server>` | Install a server and register it with your AI client config |
| `apm install <server> --client <client>` | Install for a specific client |
| `apm install <server> --name <alias>` | Install with a custom name in the config |
| `apm uninstall <server>` | Remove a server from your client config |
| `apm set-env <server> <KEY>` | Update an environment variable for an installed server |

```bash
apm install github
apm install notion --client cursor
apm install github --name gh
apm uninstall github
apm set-env github GITHUB_PERSONAL_ACCESS_TOKEN
```

During install, `apm` prompts interactively for any required environment variables (e.g. API tokens). Input is hidden for secrets. If you skip a required variable, you'll get a reminder with the exact command to set it later.

---

### Local server

| Command | Description |
|---|---|
| `apm add` | Scaffold a local MCP server for the current directory |
| `apm add --name <name>` | Use a custom server name |
| `apm add --client <client>` | Register with a specific client |

---

### Config

| Command | Description |
|---|---|
| `apm config path` | Print the path to your client config file |
| `apm config edit` | Open your client config in `$EDITOR` |

```bash
apm config path
apm config edit
```

---

### Clients

| Command | Description |
|---|---|
| `apm clients` | List supported clients and their config status |
| `apm status` | Show all installed servers across all configured clients |
| `apm status --client <client>` | Show installed servers for a specific client |
| `apm doctor` | Validate your MCP client config health |

```bash
apm clients
apm status
apm status --client vscode
apm doctor
```

`apm doctor` checks that your config file exists, is valid JSON, has the correct structure, and that no required environment variables are empty.

---

## Supported clients

| Client | Config file |
|---|---|
| Claude Desktop | `claude_desktop_config.json` |
| VS Code | `mcp.json` |
| Cursor | `mcp.json` |
| Windsurf | `mcp_config.json` |

When only one client is configured, `apm` selects it automatically. When multiple clients are configured, use `--client` to specify which one.

Valid values for `--client`: `claude-desktop`, `vscode`, `cursor`, `windsurf`.

---

## Registry

ArcMesh maintains a curated registry of verified MCP servers. When a server is not found there, `apm` automatically falls back to the [official MCP Registry](https://registry.modelcontextprotocol.io). Results are clearly labelled `[arcmesh]` or `[official]`.

Want to add a server to the ArcMesh registry? See [arcmesh-registry](https://github.com/arcmesh-labs/arcmesh-registry).

---

## WSL

On WSL, `apm` automatically detects your Windows AI client configs and wraps server entries with `wsl.exe` so Windows-hosted clients (Claude Desktop, VS Code, Windsurf) can reach WSL-hosted servers. No manual config needed.

---

## License

MIT