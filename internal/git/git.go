package git

import (
	"os"
	"path/filepath"
	"strings"
)

func FindGitDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		gitPath := filepath.Join(dir, ".git")
		if fileExists(gitPath) {
			info, err := os.Stat(gitPath)
			if err != nil {
				return ""
			}

			if info.IsDir() {
				return gitPath
			}

			content, err := os.ReadFile(gitPath)
			if err != nil {
				return ""
			}

			if strings.HasPrefix(string(content), "gitdir:") {
				gitDir := strings.TrimSpace(strings.TrimPrefix(string(content), "gitdir:"))
				if !filepath.IsAbs(gitDir) {
					gitDir = filepath.Join(dir, gitDir)
				}
				return gitDir
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return ""
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
