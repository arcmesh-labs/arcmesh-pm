package install

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/arcmesh-labs/arcmesh-pm/internal/registry"
)

// RunInstall runs the install step for a manifest.
// npx and uvx entries only require the tool to exist — no separate install step.
func RunInstall(install registry.ManifestInstall) error {
	switch install.Type {
	case "pip":
		if install.Package == "" {
			return fmt.Errorf("manifest install.package is missing")
		}
		python, err := FindPython()
		if err != nil {
			return err
		}
		cmd := exec.Command(python, "-m", "pip", "install", install.Package)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("pip install failed: %w", err)
		}
	case "npx":
		if install.Package == "" {
			return fmt.Errorf("manifest install.package is missing")
		}
		if _, err := exec.LookPath("npx"); err != nil {
			return fmt.Errorf("'npx' not found. Install Node.js: https://nodejs.org")
		}
	case "uvx":
		if install.Package == "" {
			return fmt.Errorf("manifest install.package is missing")
		}
		if _, err := exec.LookPath("uvx"); err != nil {
			return fmt.Errorf("'uvx' not found. Install uv: https://docs.astral.sh/uv/getting-started/installation/")
		}
	default:
		return fmt.Errorf("unknown install type '%s'", install.Type)
	}
	return nil
}

func FindPython() (string, error) {
	candidates := []string{"python3", "python"}
	if runtime.GOOS == "windows" {
		candidates = []string{"python"}
	}
	for _, name := range candidates {
		p, err := exec.LookPath(name)
		if err != nil {
			continue
		}
		if runtime.GOOS == "windows" && strings.Contains(p, "WindowsApps") {
			continue
		}
		return p, nil
	}
	return "", fmt.Errorf("python not found. Install Python: https://python.org")
}
