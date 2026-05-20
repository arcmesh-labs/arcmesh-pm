# arcmesh-pm

The package manager for MCP servers.

`apm` installs and manages MCP servers for AI clients like Claude Desktop, VS Code, Cursor, and Windsurf — with a single command and no infrastructure required.

---

## What is MCP?

MCP (Model Context Protocol) lets AI assistants connect to external tools — GitHub, Notion, databases, your own codebase, and more. Without a package manager, setting up an MCP server means: find the package, read the docs, install manually, write JSON config, find the right config file for your client. `apm` eliminates all of that.

---

## Getting started

### Install

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

### First run

Run `apm` to see all available commands:

```
$ apm

Usage:
  apm [command]

  arcmesh-pm (apm) — the package manager for MCP servers.

Registry:
  search           Search for MCP servers across ArcMesh and official MCP Registry.
  list             List all available MCP servers in the ArcMesh registry.

Install:
  add              Scaffold a local MCP server for the current directory and register it with your AI client config.
  install          Install an MCP server and register it with your AI client config.
  uninstall        Remove an MCP server from your AI client config.
  set-env          Update an environment variable for an installed MCP server.

Config:
  config edit      Open your MCP client config file in $EDITOR.
  config path      Print the path to your MCP client config file.

Clients:
  clients          List all supported MCP clients and whether they are configured on this system.
  status           List all installed MCP servers across configured clients.
  doctor           Check MCP client config health.
```

---

### Make your project AI-ready

Go to any project and run `apm add`. If you have multiple AI clients installed and none have MCP configured yet, you'll be asked to pick one:

```
$ cd my-project
$ apm add

Select a client:
  1. claude-desktop
  2. vscode
  3. cursor
  4. windsurf
Choose [1-4]: 2

✓ my-project added successfully.
  Config written to: ~/.config/Code/User/mcp.json
  Server script:     ~/my-project/.mcp/server.py

Next steps:
  1. Restart VS Code
  2. Look for my-project in the MCP tools panel
```

Once you have one client configured with an active MCP server, `apm` remembers it as your default — no need to specify `--client` again. If you add more clients later, `apm` will ask again.

---

### Install a server from the registry

Connect your AI assistant to external services like GitHub, Notion, or Slack:

```
$ apm install github

Fetching manifest for github...
Installing github...
? Enter GITHUB_PERSONAL_ACCESS_TOKEN (required, secret): ********

✓ github installed successfully.
  Config written to: ~/.config/Code/User/mcp.json

Next steps:
  1. Restart VS Code
  2. Look for github in the MCP tools panel
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

When only one client is configured with an active MCP server, `apm` selects it automatically. When multiple are active, use `--client` to specify which one.

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