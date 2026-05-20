#!/bin/bash
# Integration test for `apm add` on Linux.
# Runs inside an Ubuntu 24.04 Docker container; exit 0 = all pass, 1 = any fail.
set -uo pipefail

PASS=0
FAIL=0
APM="$HOME/.local/bin/apm"
TEST_REPO="$HOME/test-repo"
CLAUDE_CONFIG="$HOME/.config/Claude/claude_desktop_config.json"

pass() { echo "  [PASS] $1"; ((PASS++)); }
fail() { echo "  [FAIL] $1"; ((FAIL++)); }
header() { echo; echo "──── $1 ────────────────────────────────────────"; }

# ── 1. Install apm ────────────────────────────────────────────────────────────

header "Install apm"
if curl -fsSL https://raw.githubusercontent.com/arcmesh-labs/arcmesh-pm/master/install.sh | sh; then
    if [ -x "$APM" ]; then
        pass "apm installed at $APM"
    else
        fail "install script succeeded but $APM not found"
        exit 1
    fi
else
    fail "install script exited non-zero"
    exit 1
fi

# ── 2. Create fake AI client config ──────────────────────────────────────────

header "Setup"
mkdir -p "$(dirname "$CLAUDE_CONFIG")"
echo '{}' > "$CLAUDE_CONFIG"
pass "claude_desktop_config.json created"

# ── 3. Run apm add ────────────────────────────────────────────────────────────

header "apm add"
mkdir -p "$TEST_REPO"
cd "$TEST_REPO"
if "$APM" add; then
    pass "apm add exited 0"
else
    fail "apm add exited non-zero"
    exit 1
fi

# ── 4. File existence ─────────────────────────────────────────────────────────

header "File existence"
SERVER_PY="$TEST_REPO/.mcp/server.py"
MCP_CONFIG="$TEST_REPO/.mcp/config.json"

[ -f "$SERVER_PY"    ] && pass ".mcp/server.py exists"    || fail ".mcp/server.py missing"
[ -f "$MCP_CONFIG"   ] && pass ".mcp/config.json exists"  || fail ".mcp/config.json missing"
[ -f "$CLAUDE_CONFIG" ] && pass "claude_desktop_config.json still present" || fail "claude_desktop_config.json missing after apm add"

# ── 5. server.py content ──────────────────────────────────────────────────────

header "server.py content"
for sym in BASE_DIR FastMCP read_file list_directory search_content; do
    grep -q "$sym" "$SERVER_PY" \
        && pass "server.py contains '$sym'" \
        || fail "server.py missing '$sym'"
done

# ── 6. .mcp/config.json structure ────────────────────────────────────────────

header ".mcp/config.json structure"
python3 - "$MCP_CONFIG" <<'PYEOF'
import json, sys

cfg = json.load(open(sys.argv[1]))
servers = cfg.get("servers", {})
errors = []

if not servers:
    errors.append("no 'servers' key")
else:
    name = next(iter(servers))
    if name != "test-repo":
        errors.append(f"server name is '{name}', want 'test-repo'")
    entry = servers.get(name, {})
    cmd = entry.get("command", "")
    if "python" not in cmd:
        errors.append(f"command is '{cmd}', want python or python3")
    args = entry.get("args", [])
    if not args or ".mcp/server.py" not in args[0]:
        errors.append(f"args do not point to .mcp/server.py: {args}")

if errors:
    print("  errors:", "; ".join(errors))
    sys.exit(1)

entry = servers[next(iter(servers))]
print(f"    command={entry['command']}  args={entry['args']}")
PYEOF

if [ $? -eq 0 ]; then
    pass ".mcp/config.json has correct command and server.py path"
else
    fail ".mcp/config.json has wrong structure"
fi

# ── 7. claude_desktop_config.json ────────────────────────────────────────────

header "claude_desktop_config.json"
python3 - "$CLAUDE_CONFIG" "$TEST_REPO" <<'PYEOF'
import json, sys

cfg  = json.load(open(sys.argv[1]))
repo = sys.argv[2]
mcp  = cfg.get("mcpServers", {})
errors = []

if "test-repo" not in mcp:
    errors.append(f"'test-repo' not in mcpServers (found: {list(mcp.keys())})")
else:
    entry = mcp["test-repo"]
    cmd   = entry.get("command", "")
    args  = entry.get("args", [])
    if cmd == "wsl.exe":
        # WSL: args[-1] is "python3 \"/path/to/.mcp/server.py\""
        if not args or ".mcp/server.py" not in args[-1]:
            errors.append(f"wsl.exe args[-1] does not contain .mcp/server.py: {args}")
    elif "python" in cmd:
        # Native Linux/macOS: args[0] is the absolute path to server.py
        if not args or ".mcp/server.py" not in args[0]:
            errors.append(f"args[0] does not point to .mcp/server.py: {args}")
        elif not args[0].startswith(repo):
            errors.append(f"server.py path '{args[0]}' is not under {repo}")
    else:
        errors.append(f"command is '{cmd}', want python, python3, or wsl.exe")

if errors:
    print("  errors:", "; ".join(errors))
    sys.exit(1)

entry = mcp["test-repo"]
print(f"    command={entry['command']}  args={entry['args']}")
PYEOF

if [ $? -eq 0 ]; then
    pass "claude_desktop_config.json has mcpServers['test-repo'] with correct path"
else
    fail "claude_desktop_config.json mcpServers check failed"
fi

# ── 8. MCP server stdio test ──────────────────────────────────────────────────

header "MCP server stdio (initialize)"
python3 - "$SERVER_PY" <<'PYEOF'
import asyncio, sys

SERVER_PY = sys.argv[1]

async def run():
    from mcp.client.stdio import stdio_client
    try:
        from mcp import ClientSession, StdioServerParameters
    except ImportError:
        from mcp.client.stdio import StdioServerParameters
        from mcp import ClientSession

    params = StdioServerParameters(command="python3", args=[SERVER_PY])
    async with stdio_client(params) as streams:
        read, write = streams
        async with ClientSession(read, write) as session:
            result = await session.initialize()
            print(f"    serverInfo.name    = {result.serverInfo.name}")
            print(f"    protocolVersion    = {result.protocolVersion}")

try:
    asyncio.run(run())
    sys.exit(0)
except Exception as exc:
    print(f"  error: {exc}", file=sys.stderr)
    sys.exit(1)
PYEOF

if [ $? -eq 0 ]; then
    pass "server responded to MCP initialize"
else
    fail "server did not respond to MCP initialize"
fi

# ── Summary ───────────────────────────────────────────────────────────────────

echo
echo "════════════════════════════════════════"
echo "  PASSED : $PASS"
echo "  FAILED : $FAIL"
echo "════════════════════════════════════════"

[ "$FAIL" -eq 0 ]
