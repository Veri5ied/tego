package hooks

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/veri5ied/tego/internal/config"
)

type Runner struct{}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) Run(hookName string, cfg *config.Config, args []string) error {
	commands, err := cfg.GetCommands(hookName)
	if err != nil {
		// If Hooks are not configured - exit gracefully
		return nil
	}

	// Set environment variables for hooks
	os.Setenv("TEGO_HOOK_NAME", hookName)

	for i, cmdStr := range commands {
		for j, arg := range args {
			placeholder := fmt.Sprintf("$%d", j+1)
			cmdStr = strings.ReplaceAll(cmdStr, placeholder, arg)
		}

		cmdStr = strings.ReplaceAll(cmdStr, "$@", strings.Join(args, " "))

		if len(commands) > 1 {
			fmt.Printf("\n[%d/%d] Running: %s\n", i+1, len(commands), cmdStr)
		} else {
			fmt.Printf("Running: %s\n", cmdStr)
		}

		cmd := exec.Command("sh", "-c", cmdStr)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Env = os.Environ()

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("command failed: %v", err)
		}
	}

	return nil
}
