package clients

import (
	"os"
	"path/filepath"
	"strings"
)

type ClientConfig struct {
	Name  string
	Paths []string
}

var ClientConfigs = map[string][]string{
	"claude-desktop": {
		filepath.Join(home(), "AppData", "Roaming", "Claude", "claude_desktop_config.json"),
		filepath.Join(home(), "AppData", "Local", "Packages", "AnthropicPBC.Claude_p3wjhkq2mfxz8", "LocalCache", "Roaming", "Claude", "claude_desktop_config.json"),
		filepath.Join(home(), "Library", "Application Support", "Claude", "claude_desktop_config.json"),
		filepath.Join(home(), ".config", "Claude", "claude_desktop_config.json"),
	},
	"vscode": {
		filepath.Join(home(), "AppData", "Roaming", "Code", "User", "mcp.json"),
		filepath.Join(home(), "Library", "Application Support", "Code", "User", "mcp.json"),
		filepath.Join(home(), ".config", "Code", "User", "mcp.json"),
	},
	"cursor": {
		filepath.Join(home(), ".cursor", "mcp.json"),
	},
	"windsurf": {
		filepath.Join(home(), ".codeium", "windsurf", "mcp_config.json"),
	},
}

// ClientOrder preserves display/iteration order (maps are unordered in Go).
var ClientOrder = []string{"claude-desktop", "vscode", "cursor", "windsurf"}

func home() string {
	h, _ := os.UserHomeDir()
	return h
}

// IsWSL reports whether the process is running inside Windows Subsystem for Linux.
func IsWSL() bool {
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(data)), "microsoft")
}

// WSLDistro returns the WSL distribution name from $WSL_DISTRO_NAME, falling
// back to "Ubuntu" when the variable is not set.
func WSLDistro() string {
	if d := os.Getenv("WSL_DISTRO_NAME"); d != "" {
		return d
	}
	return "Ubuntu"
}

// FindConfig returns (resolvedPath, found). resolvedPath is empty when not found.
func FindConfig(client string) (string, bool) {
	paths, ok := ClientConfigs[client]
	if !ok {
		return "", false
	}

	if IsWSL() {
		var glob, base string
		switch client {
		case "claude-desktop":
			base = "/mnt/c/Users"
			glob = "*/AppData/Roaming/Claude/claude_desktop_config.json"
		case "vscode":
			base = "/mnt/c/Users"
			glob = "*/AppData/Roaming/Code/User/mcp.json"
		case "cursor":
			base = "/mnt/c/Users"
			glob = "*/.cursor/mcp.json"
		case "windsurf":
			base = "/mnt/c/Users"
			glob = "*/.codeium/windsurf/mcp_config.json"
		}
		if glob != "" {
			matches, _ := filepath.Glob(filepath.Join(base, glob))
			if len(matches) > 0 {
				return matches[0], true
			}
		}
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p, true
		}
	}
	return "", false
}

// DisplayName returns the human-readable name for a client.
func DisplayName(client string) string {
	switch client {
	case "vscode":
		return "VS Code"
	case "cursor":
		return "Cursor"
	case "windsurf":
		return "Windsurf"
	default:
		return "Claude Desktop"
	}
}

// ConfigFilename returns the config filename for a client.
func ConfigFilename(client string) string {
	switch client {
	case "windsurf":
		return "mcp_config.json"
	case "vscode", "cursor":
		return "mcp.json"
	default:
		return "claude_desktop_config.json"
	}
}

// CheckedPaths returns all paths that would be checked for a client (for error messages).
func CheckedPaths(client string) []string {
	paths, ok := ClientConfigs[client]
	if !ok {
		return nil
	}
	result := make([]string, len(paths))
	copy(result, paths)
	if IsWSL() {
		switch client {
		case "claude-desktop":
			result = append(result, "/mnt/c/Users/*/AppData/Roaming/Claude/claude_desktop_config.json")
		case "vscode":
			result = append(result, "/mnt/c/Users/*/AppData/Roaming/Code/User/mcp.json")
		case "cursor":
			result = append(result, "/mnt/c/Users/*/.cursor/mcp.json")
		case "windsurf":
			result = append(result, "/mnt/c/Users/*/.codeium/windsurf/mcp_config.json")
		}
	}
	return result
}
