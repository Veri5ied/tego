package hooks

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/veri5ied/tego/internal/config"
)

type Installer struct {
	gitDir string
}

func NewInstaller(gitDir string) *Installer {
	return &Installer{gitDir: gitDir}
}

func (i *Installer) Install(cfg *config.Config) (int, error) {
	hooksDir := filepath.Join(i.gitDir, "hooks")
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return 0, fmt.Errorf("failed to create hooks directory: %v", err)
	}

	// Find tego binary path
	tegoPath, err := exec.LookPath("tego")
	if err != nil {
		// Fallback to current executable path
		tegoPath, err = os.Executable()
		if err != nil {
			return 0, fmt.Errorf("cannot locate tego binary: %v", err)
		}
	}

	installed := 0
	for hookName := range cfg.Hooks {
		hookPath := filepath.Join(hooksDir, hookName)

		hookScript := fmt.Sprintf(`#!/bin/sh
# Installed by Tego
# https://github.com/veri5ied/tego

%s run %s "$@"
`, tegoPath, hookName)

		err := os.WriteFile(hookPath, []byte(hookScript), 0755)
		if err != nil {
			fmt.Printf("⚠ Warning: Failed to install %s: %v\n", hookName, err)
			continue
		}
		fmt.Printf("✓ Installed %s\n", hookName)
		installed++
	}

	if installed == 0 {
		return 0, fmt.Errorf("no hooks were installed")
	}

	return installed, nil
}
