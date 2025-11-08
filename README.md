# Tego

Tego is a minimal, declarative Git hooks manager that works with any programming language. It provides a simple, configuration-based approach to managing Git hooks without requiring manual script management or complex setup procedures.

## Overview

Tego is built as a single native binary written in Go, offering fast execution and zero runtime dependencies. Unlike language-specific tools, Tego can be used in projects regardless of their technology stack—JavaScript, Python, Ruby, Go, Rust, PHP, or any other language.

## Features

- **Universal Language Support** - Works with any programming language or framework
- **Simple Configuration** - Define hooks using JSON, TOML, or YAML
- **Native Performance** - Single binary with no runtime dependencies
- **Explicit Installation** - Manual hook installation with no magic scripts
- **Multiple Commands** - Run multiple commands per hook sequentially
- **Standard Git Hooks** - Supports all standard Git hook types

## Installation

### Download Binary

Download the appropriate binary for your platform from the [releases page](https://github.com/veri5ied/tego/releases).

**Linux/macOS:**

```bash
curl -L https://github.com/veri5ied/tego/releases/download/v1.0.0/tego-linux-amd64 -o tego
chmod +x tego
sudo mv tego /usr/local/bin/
```

**Windows:**
Download `tego-windows-amd64.exe` from releases and add to your PATH.

### Build from Source

```bash
git clone https://github.com/veri5ied/tego.git
cd tego
go build -o tego cmd/tego/main.go
```

## Quick Start

Initialize Tego in your project:

```bash
tego init
```

This creates a `.tegorc.json` file with sample configuration.

Edit the configuration file to define your hooks:

```json
{
  "hooks": {
    "pre-commit": "npm run lint && npm test",
    "commit-msg": "npx commitlint --edit $1",
    "pre-push": "npm run build"
  }
}
```

Install the hooks:

```bash
tego install
```

The hooks are now active and will run automatically during Git operations.

## Configuration

Tego looks for configuration files in the following order:

- `.tegorc.json`
- `.tegorc`
- `tego.json`
- `tego.toml`
- `tego.yaml`
- `tego.yml`

### JSON Format

```json
{
  "hooks": {
    "pre-commit": "npm test",
    "pre-push": "npm run lint"
  }
}
```

### TOML Format

```toml
[hooks]
pre-commit = "npm test"
pre-push = "npm run lint"
```

### YAML Format

```yaml
hooks:
  pre-commit: npm test
  pre-push: npm run lint
```

### Multiple Commands

You can define multiple commands for a single hook:

```json
{
  "hooks": {
    "pre-commit": ["npm run format", "npm run lint", "npm test"]
  }
}
```

Commands are executed sequentially. If any command fails, the Git operation is aborted.

## Commands

### init

Create a sample configuration file:

```bash
tego init
```

### install

Install configured hooks to `.git/hooks`:

```bash
tego install
```

### uninstall

Remove Tego-managed hooks:

```bash
tego uninstall
```

### list

Display all configured hooks:

```bash
tego list
```

### run

Manually execute a specific hook:

```bash
tego run pre-commit
```

### version

Display version information:

```bash
tego version
```

## Supported Hooks

Tego supports all standard Git hooks:

- `pre-commit` - Runs before a commit is created
- `prepare-commit-msg` - Runs before the commit message editor is opened
- `commit-msg` - Validates or modifies the commit message
- `post-commit` - Runs after a commit is created
- `pre-push` - Runs before pushing to a remote
- `pre-rebase` - Runs before a rebase operation
- `post-merge` - Runs after a merge operation
- `post-checkout` - Runs after checking out a branch

And many others supported by Git.

## Usage Examples

### JavaScript/TypeScript Project

**.tegorc.json:**

```json
{
  "hooks": {
    "pre-commit": ["npm run lint", "npm run format", "npm test"],
    "commit-msg": "npx commitlint --edit $1",
    "pre-push": "npm run build"
  }
}
```

### Python Project

**.tegorc.json:**

```json
{
  "hooks": {
    "pre-commit": "black . && pylint src/",
    "pre-push": "pytest"
  }
}
```

### Go Project

**.tegorc.json:**

```json
{
  "hooks": {
    "pre-commit": ["go fmt ./...", "go vet ./...", "go test ./..."],
    "pre-push": "go test -race ./..."
  }
}
```

### Monorepo

**.tegorc.json:**

```json
{
  "hooks": {
    "pre-commit": [
      "cd packages/frontend && npm run lint",
      "cd packages/backend && npm run lint",
      "npm test"
    ]
  }
}
```

## Skipping Hooks

To skip hooks for a specific Git operation, set the `TEGO_SKIP` environment variable:

```bash
TEGO_SKIP=1 git commit -m "skip hooks for this commit"
```

On Windows:

```cmd
set TEGO_SKIP=1 && git commit -m "skip hooks for this commit"
```

## How It Works

When you run `tego install`, Tego creates shell scripts in `.git/hooks/` that delegate to the `tego run` command. For example, the `pre-commit` hook becomes:

```bash
#!/bin/sh
# Installed by Tego
tego run pre-commit "$@"
```

When Git triggers the hook, Tego reads your configuration file and executes the defined commands. If any command fails, the Git operation is aborted.

## Team Collaboration

Commit your configuration file (`.tegorc.json`) to version control. Team members can then install hooks by running:

```bash
tego install
```

This ensures all team members use the same hooks and maintain consistent code quality standards.

## Design Philosophy

Tego is designed with the following principles:

- **Language Agnostic** - Should work seamlessly in any project, regardless of programming language
- **Explicit over Magic** - Hook installation is a deliberate action, not an automatic side effect
- **Configuration as Code** - Use structured configuration files (JSON, TOML, YAML) rather than executable scripts
- **Zero Dependencies** - No runtime requirements beyond Git itself
- **Performance First** - Native binary execution for fast hook processing

## Requirements

- Git 2.0 or higher
- Unix-like shell (Linux, macOS, WSL, Git Bash on Windows)

## Contributing

Contributions are welcome. Please open an issue or submit a pull request on GitHub.

## License

MIT License. See LICENSE file for details.

## Links

- [GitHub Repository](https://github.com/veri5ied/tego)
- [Issue Tracker](https://github.com/veri5ied/tego/issues)
- [Releases](https://github.com/veri5ied/tego/releases)

## Acknowledgments

Inspired by Husky and other Git hooks management tools, Tego aims to provide a universal solution that works across all programming languages and environments.
