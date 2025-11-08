package hooks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/veri5ied/tego/internal/config"
)

type Uninstaller struct {
	gitDir string
}

func NewUninstaller(gitDir string) *Uninstaller {
	return &Uninstaller{gitDir: gitDir}
}

func (u *Uninstaller) Uninstall(cfg *config.Config) int {
	hooksDir := filepath.Join(u.gitDir, "hooks")
	removed := 0

	for hookName := range cfg.Hooks {
		hookPath := filepath.Join(hooksDir, hookName)

		if !fileExists(hookPath) {
			continue
		}

		// Read file content to verify it's a Tego hook
		content, err := os.ReadFile(hookPath)
		if err != nil {
			continue
		}

		if strings.Contains(string(content), "Installed by Tego") {
			if err := os.Remove(hookPath); err == nil {
				fmt.Printf("✓ Removed %s\n", hookName)
				removed++
			} else {
				fmt.Printf("⚠ Warning: Failed to remove %s: %v\n", hookName, err)
			}
		}
	}

	return removed
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
