# AGENTS.md

This file provides guidance to AI coding agents when working with code in this repository.

## General agent rules

- When users ask questions, answer them instead of doing the work.

### Shell Rules

- Always use `rm -f` (never bare `rm`)
- Before running a series of `git` commands, confirm you are in the project root; if not, `cd` there first. Then run all subsequent `git` commands from that directory without the `-C` option.

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

- ✅ Specification document (`docs/specification.md`)
- ✅ Go implementation (all phases complete)

## Directory Structure

```txt
claude-launcher/
├── docs/
│   └── specification.md           # Specification
├── cmd/
│   └── claude-launcher/           # Go implementation entry point
├── internal/
│   ├── config/                    # Configuration loading
│   ├── security/                  # Directory checking
│   ├── session/                   # Session management
│   ├── launcher/                  # Claude Code execution
│   └── ui/                        # User interface
├── mise.toml                      # mise tool/task definitions
├── README.md                      # Project README
└── AGENTS.md                      # This file
```

## Important Documents

### Specification (`docs/specification.md`)

Detailed specification analyzing the Bash version. Includes:

- Detailed core features
- Configuration methods
- Command-line interface
- Error messages and user feedback
- Execution flow diagram

## Development

### Build and Test

```bash
# Build
mise run go-build

# Run tests
mise run go-test

# Run all checks
mise run check

# Format and check
mise run fix-and-check
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

**IMPORTANT**: Always run `mise run fix-and-check` after any changes and fix all errors before proceeding.

This applies to:

- Go code modifications
- Markdown file edits
- YAML/JSON configuration changes
- Any other file modifications

### Go Coding Style

- Format with `gofumpt` + `goimports` (via `golangci-lint run --fix`)
- Static analysis with `golangci-lint`
- Target 80%+ test coverage

#### Enabled Linters

| Linter | Category | Default |
| -------- | ---------- | --------- |
| govet | Bug detection | Yes |
| staticcheck | Static analysis | Yes |
| errcheck | Error handling | Yes |
| unused | Dead code | Yes |
| gosec | Security | No |
| bidichk | Security | No |
| errorlint | Error handling | No |
| bodyclose | Resource leak | No |
| unconvert | Dead code | No |
| usestdlibvars | Clean code | No |
| modernize | Modernization | No |

### Markdown Style

- Use dashes (`-`) for unordered lists
- Use asterisks (`*`) for emphasis and strong
- Code block language specifications:
  - Directory structures: `txt`
  - Shell examples: `sh` or `bash`
  - Go code: `go`
  - JSON: `json`
  - YAML: `yml` or `yaml`

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
mise run go-test

# With coverage
mise run go-test-cover

# Generate coverage report
mise run go-test-coverage-report
```

## References

- [Claude Code Official Documentation](https://code.claude.com/docs/en/overview)
- [Go Project Layout](https://github.com/golang-standards/project-layout)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
