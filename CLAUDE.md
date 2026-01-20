# claude-launcher Project

## Language Settings / 言語設定

**For Claude Code**: When working on this project, please follow these language guidelines:

- **Conversation responses**: Japanese or English (user's preference)
- **Everything else**: English only
  - Code (variables, functions, comments)
  - Documentation files (except this file)
  - Commit messages
  - Error messages
  - Log messages
  - Test names and descriptions

**Claude Code へ**: このプロジェクトでは、以下の言語ガイドラインに従ってください：

- **会話の応答**: 日本語または英語（ユーザーの好みに応じて）
- **それ以外のすべて**: 英語のみ
  - コード（変数、関数、コメント）
  - ドキュメントファイル（このファイルを除く）
  - コミットメッセージ
  - エラーメッセージ
  - ログメッセージ
  - テスト名と説明

## Project Overview

`claude-launcher` is a comprehensive launcher tool for Claude Code. It provides both safety and convenience through two main features:

1. **Directory Security**: Only allows Claude Code execution in permitted directories
2. **Session Management**: Choice between continuing previous session or starting fresh

## Current Status

### Completed

- ✅ Complete Bash implementation (`tmp/claude-launcher.sh`)
- ✅ Specification document (`docs/specification.md`)
- ✅ Go implementation plan (`docs/implementation-plan.md`)
- ✅ Go implementation (all phases complete)
  - Phase 1: Basic structure (config loading, directory check)
  - Phase 2: UI and session management
  - Phase 3: Launch process and CLI
  - Phase 4: Integration tests and release preparation

## Directory Structure

```txt
claude-launcher/
├── docs/
│   ├── specification.md         # Specification
│   └── implementation-plan.md   # Go implementation plan
├── tmp/
│   └── claude-launcher.sh       # Bash version (reference)
├── cmd/
│   └── claude-launcher/         # Go implementation entry point
├── internal/
│   ├── config/                  # Configuration loading
│   ├── security/                # Directory checking
│   ├── session/                 # Session management
│   ├── launcher/                # Claude Code execution
│   └── ui/                      # User interface
├── tasks/
│   └── GoTasks.yml              # Taskfile task definitions
├── README.md                    # Project README
└── CLAUDE.md                    # This file
```

## Important Documents

### 1. Specification (`docs/specification.md`)

Detailed specification analyzing the Bash version. Includes:

- Detailed core features
- Configuration methods
- Command-line interface
- Error messages and user feedback
- Execution flow diagram

### 2. Implementation Plan (`docs/implementation-plan.md`)

Go language implementation plan. Includes:

- Project structure
- Package design
- Type definitions and key functions
- Implementation steps (Phase 1-4)
- Testing strategy
- Build and deployment methods

## Development

### Build and Test

```bash
# Build
task go:build

# Run tests
task go:test

# Run all checks
task check

# Format and check
task fix-and-check
```

### Run

```bash
# Set configuration
export CLAUDE_SAFE_DIRS="$HOME/develop:$HOME/projects"

# Show help
./bin/claude-launcher --help

# Show allowed directories
./bin/claude-launcher --show-dirs

# Run in allowed directory
cd ~/develop
./bin/claude-launcher
```

## Configuration Examples

### Environment Variable

```bash
export CLAUDE_SAFE_DIRS="$HOME/develop:$HOME/projects"
```

### ~/.config/claude-launcher/config.json

```json
{
  "allowedDirs": [
    "/home/user/develop",
    "/home/user/projects"
  ]
}
```

## Usage Examples

### Basic Usage

```bash
# Check directory and prompt for session continuation
claude-launcher

# Show help
claude-launcher --help

# Show allowed directories
claude-launcher --show-dirs

# Pass arguments to Claude
claude-launcher --model opus
```

## Development Guidelines

### General Rule

**IMPORTANT**: Always run `task fix-and-check` after any changes and fix all errors before proceeding.

This applies to:

- Go code modifications
- Markdown file edits
- YAML/JSON configuration changes
- Any other file modifications

### Go Coding Style

- Format with `gofmt`
- Static analysis with `golangci-lint`
- Target 80%+ test coverage

### Markdown Style

- Code block language specifications:
  - Directory structures: `txt`
  - Shell examples: `sh`
  - Go code: `go`
  - JSON: `json`
  - YAML: `yml`

### Commit Messages

- Clear and concise descriptions
- Recommended prefixes:
  - `feat:` New feature
  - `fix:` Bug fix
  - `docs:` Documentation
  - `test:` Test addition/modification
  - `refactor:` Refactoring

### Testing

```bash
# Run all tests
task go:test

# With coverage
task go:test-cover

# Generate coverage report
task go:test-coverage-report
```

## References

- [Claude Code Official Documentation](https://claude.ai/claude-code)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
- Shell script version: `tmp/claude-launcher.sh`

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
