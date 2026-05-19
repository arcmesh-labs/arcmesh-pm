package install

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"

	"github.com/arcmesh-labs/arcmesh-pm/internal/clients"
	"github.com/arcmesh-labs/arcmesh-pm/internal/registry"
)

var envWarnStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))

// PromptEnv prompts the user for each environment variable in envSpec.
// Secrets are read with terminal echo disabled. Returns the resolved values
// and a list of required keys the user left empty.
func PromptEnv(envSpec map[string]registry.EnvSpec, client string) (map[string]string, []string, error) {
	if len(envSpec) == 0 {
		return nil, nil, nil
	}

	for _, spec := range envSpec {
		if spec.Secret {
			fmt.Printf("%s Tokens are stored in plain text in %s. Keep this file private.\n",
				envWarnStyle.Render("⚠"), clients.ConfigFilename(client))
			break
		}
	}

	// Sort keys for deterministic prompt order (Go maps have no stable iteration order).
	keys := make([]string, 0, len(envSpec))
	for k := range envSpec {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	resolved := make(map[string]string, len(envSpec))
	var skipped []string

	for _, key := range keys {
		spec := envSpec[key]
		fmt.Printf("\n  %s\n", spec.Description)

		var value string
		var err error
		if spec.Secret {
			value, err = readSecret(fmt.Sprintf("  %s", key))
		} else {
			value, err = readLine(fmt.Sprintf("  %s: ", key))
		}
		if err != nil {
			return nil, nil, err
		}

		if value == "" && spec.Required {
			skipped = append(skipped, key)
		}
		resolved[key] = value
	}

	return resolved, skipped, nil
}

func readLine(prompt string) (string, error) {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimRight(scanner.Text(), "\r\n"), nil
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", nil
}

// ReadSecret reads a value with terminal echo disabled (exported for set-env command).
func ReadSecret(prompt string) (string, error) {
	return readSecret(prompt)
}

func readSecret(prompt string) (string, error) {
	fmt.Print(prompt + ": ")
	if term.IsTerminal(int(os.Stdin.Fd())) {
		b, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	// stdin is not a terminal (piped input) — fall back to visible read
	return readLine("")
}
