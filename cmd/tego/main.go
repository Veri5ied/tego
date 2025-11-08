package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/veri5ied/tego/internal/config"
	"github.com/veri5ied/tego/internal/git"
	"github.com/veri5ied/tego/internal/hooks"
)

const version = "1.0.0"

var (
	successColor = color.New(color.FgGreen, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
	infoColor    = color.New(color.FgCyan)
	warnColor    = color.New(color.FgYellow)
)

func main() {
	if os.Getenv("TEGO_SKIP") == "1" || os.Getenv("TEGO_SKIP") == "true" {
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		initTego()
	case "install":
		installHooks()
	case "uninstall":
		uninstallHooks()
	case "run":
		if len(os.Args) < 3 {
			errorColor.Println("✗ Error: 'run' command requires a hook name")
			fmt.Println("Usage: tego run <hook-name>")
			os.Exit(1)
		}
		runHook(os.Args[2])
	case "list":
		listHooks()
	case "version", "-v", "--version":
		fmt.Printf("tego v%s\n", version)
	case "help", "-h", "--help":
		printHelp()
	default:
		errorColor.Printf("✗ Unknown command: %s\n\n", command)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	help := `
╔═══════════════════════════════════════════════════════════╗
║                        Tego v` + version + `                          ║
║         Minimal Git Hooks Manager for Any Language        ║
╚═══════════════════════════════════════════════════════════╝

Usage:
  tego init              Create a sample .tegorc.json config
  tego install           Install hooks to .git/hooks
  tego uninstall         Remove hooks from .git/hooks
  tego run <hook-name>   Run a specific hook manually
  tego list              List all configured hooks
  tego version           Show version
  tego help              Show this help

Examples:
  # Set up Tego in your project
  tego init
  tego install

  # Run a specific hook
  tego run pre-commit

  # Skip hooks for a single commit
  TEGO_SKIP=1 git commit -m "skip hooks"

Configuration:
  Tego reads from .tegorc.json, .tegorc, tego.json, tego.toml, 
  or tego.yaml in your project root.

JavaScript/TypeScript:
  npm install --save-dev tego-kit
  npx tego init
  npx tego install

Python:
  pip install tego
  tego init
  tego install

Learn more: https://github.com/veri5ied/tego
`
	fmt.Print(help)
}

func initTego() {
	if config.ConfigExists() {
		warnColor.Println("⚠ Config file already exists!")
		fmt.Println("\nFound: " + config.GetConfigPath())
		os.Exit(1)
	}

	sampleConfig := `{
  "hooks": {
    "pre-commit": "npm run lint && npm test",
    "commit-msg": "npx commitlint --edit $1",
    "pre-push": "npm run test"
  }
}
`
	err := os.WriteFile(".tegorc.json", []byte(sampleConfig), 0644)
	if err != nil {
		errorColor.Printf("✗ Error creating .tegorc.json: %v\n", err)
		os.Exit(1)
	}

	successColor.Println("✓ Created .tegorc.json")
	fmt.Println("\n" + strings.Repeat("─", 60))
	infoColor.Println("Next steps:")
	fmt.Println("  1. Edit .tegorc.json to configure your hooks")
	fmt.Println("  2. Run 'tego install' to activate the hooks")
	fmt.Println("  3. Try making a commit!")
	fmt.Println(strings.Repeat("─", 60))
}

func installHooks() {
	cfg, err := config.LoadConfig()
	if err != nil {
		errorColor.Printf("✗ Error: %v\n", err)
		fmt.Println("\n" + infoColor.Sprint("💡 Tip:") + " Run 'tego init' to create a config file")
		os.Exit(1)
	}

	gitDir := git.FindGitDir()
	if gitDir == "" {
		errorColor.Println("✗ Error: Not a git repository")
		fmt.Println("\n" + infoColor.Sprint("💡 Tip:") + " Run 'git init' first")
		os.Exit(1)
	}

	installer := hooks.NewInstaller(gitDir)
	installed, err := installer.Install(cfg)
	if err != nil {
		errorColor.Printf("✗ Installation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	successColor.Printf("✓ Installed %d hook(s) successfully\n", installed)
	fmt.Println("\n" + strings.Repeat("─", 60))
	infoColor.Println("Hooks are now active! 🎉")
	fmt.Println("Your configured hooks will run automatically on git operations.")
	fmt.Println("\nTo skip hooks: TEGO_SKIP=1 git commit -m \"message\"")
	fmt.Println(strings.Repeat("─", 60))
}

func uninstallHooks() {
	gitDir := git.FindGitDir()
	if gitDir == "" {
		errorColor.Println("✗ Error: Not a git repository")
		os.Exit(1)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		warnColor.Printf("⚠ Warning: %v\n", err)
		fmt.Println("Will try to remove common hooks...")
		cfg = &config.Config{
			Hooks: map[string]interface{}{
				"pre-commit":  nil,
				"commit-msg":  nil,
				"pre-push":    nil,
				"post-commit": nil,
				"post-merge":  nil,
				"pre-rebase":  nil,
			},
		}
	}

	uninstaller := hooks.NewUninstaller(gitDir)
	removed := uninstaller.Uninstall(cfg)

	fmt.Println()
	if removed == 0 {
		infoColor.Println("ℹ No Tego hooks found")
	} else {
		successColor.Printf("✓ Removed %d hook(s)\n", removed)
	}
}

func runHook(hookName string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		errorColor.Printf("✗ Error: %v\n", err)
		os.Exit(1)
	}

	runner := hooks.NewRunner()
	if err := runner.Run(hookName, cfg, os.Args[3:]); err != nil {
		errorColor.Printf("✗ Hook '%s' failed\n", hookName)
		os.Exit(1)
	}

	successColor.Printf("✓ Hook '%s' passed\n", hookName)
}

func listHooks() {
	cfg, err := config.LoadConfig()
	if err != nil {
		errorColor.Printf("✗ Error: %v\n", err)
		fmt.Println("\n" + infoColor.Sprint("💡 Tip:") + " Run 'tego init' to create a config file")
		os.Exit(1)
	}

	fmt.Println()
	infoColor.Println("Configured hooks:")
	fmt.Println(strings.Repeat("─", 60))

	if len(cfg.Hooks) == 0 {
		warnColor.Println("No hooks configured")
		return
	}

	for hookName, hookCmd := range cfg.Hooks {
		fmt.Printf("\n%s:\n", color.CyanString(hookName))
		switch v := hookCmd.(type) {
		case string:
			fmt.Printf("  → %s\n", v)
		case []interface{}:
			for i, cmd := range v {
				if cmdStr, ok := cmd.(string); ok {
					fmt.Printf("  %d. %s\n", i+1, cmdStr)
				}
			}
		}
	}
	fmt.Println()
}
