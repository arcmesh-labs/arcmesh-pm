package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	ltable "github.com/charmbracelet/lipgloss/table"

	"github.com/arcmesh-labs/arcmesh-pm/internal/clients"
)

// ── styles ────────────────────────────────────────────────────────────────────

var (
	errStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	warnStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	dimStyle     = lipgloss.NewStyle().Faint(true)
	boldStyle    = lipgloss.NewStyle().Bold(true)
	cyanStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true)
)

// ── terminal output ───────────────────────────────────────────────────────────

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s %s\n", errStyle.Render("Error:"), fmt.Sprintf(format, args...))
	os.Exit(1)
}

func verifiedBadge(v bool) string {
	if v {
		return successStyle.Render("✓ verified")
	}
	return dimStyle.Render("—")
}

func sourceBadge(source string) string {
	if source == "arcmesh" {
		return successStyle.Render("[arcmesh]")
	}
	return dimStyle.Render("[official]")
}

// serverType returns "npx" when command=="npx", otherwise "local".
func serverType(cfg map[string]interface{}) string {
	if cmd, _ := cfg["command"].(string); cmd == "npx" {
		return "npx"
	}
	return "local"
}

// envCell formats the env block of a server config entry for a table cell.
func envCell(cfg map[string]interface{}) string {
	env, _ := cfg["env"].(map[string]interface{})
	if len(env) == 0 {
		return dimStyle.Render("(no env)")
	}
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var parts []string
	for _, key := range keys {
		val, _ := env[key].(string)
		indicator := "✓"
		if strings.TrimSpace(val) == "" {
			indicator = "⚠"
		}
		parts = append(parts, indicator+" "+key)
	}
	return strings.Join(parts, "\n")
}

// ── table ─────────────────────────────────────────────────────────────────────

// newTable returns a lipgloss table with consistent border/padding styling.
func newTable() *ltable.Table {
	return ltable.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("240"))).
		BorderRow(true).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == ltable.HeaderRow {
				return lipgloss.NewStyle().Bold(true).Padding(0, 1)
			}
			return lipgloss.NewStyle().Padding(0, 1)
		})
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

// printTable prints an optional title then the table.
func printTable(title string, t *ltable.Table) {
	if title != "" {
		fmt.Println(boldStyle.Render(title))
	}
	fmt.Println(t.Render())
}

// ── client resolution ─────────────────────────────────────────────────────────

var validClients = map[string]bool{
	"claude-desktop": true,
	"vscode":         true,
	"cursor":         true,
	"windsurf":       true,
}

// resolveClientActive picks a client using the active→configured→default strategy
// used by install, uninstall, config, status, and doctor.
func resolveClientActive(flag string) (string, error) {
	if flag != "" {
		if !validClients[flag] {
			return "", fmt.Errorf("invalid client '%s'. Choose from: claude-desktop, vscode, cursor, windsurf", flag)
		}
		return flag, nil
	}
	active := clients.DetectActiveClients()
	if len(active) == 1 {
		return active[0], nil
	}
	if len(active) == 0 {
		configured := clients.DetectConfiguredClients()
		if len(configured) == 1 {
			return configured[0], nil
		}
		return "claude-desktop", nil
	}
	return "", fmt.Errorf("multiple clients configured. Specify one with --client [claude-desktop|vscode|cursor|windsurf]")
}

// resolveClientConfigured picks a client from configured clients only (used by set-env).
func resolveClientConfigured(flag string) (string, error) {
	if flag != "" {
		if !validClients[flag] {
			return "", fmt.Errorf("invalid client '%s'. Choose from: claude-desktop, vscode, cursor, windsurf", flag)
		}
		return flag, nil
	}
	configured := clients.DetectConfiguredClients()
	if len(configured) == 1 {
		return configured[0], nil
	}
	if len(configured) == 0 {
		return "claude-desktop", nil
	}
	return "", fmt.Errorf("multiple clients configured. Specify one with --client [claude-desktop|vscode|cursor|windsurf]")
}
