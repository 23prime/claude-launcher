# Claude Launcher

Comprehensive launcher for Claude Code with directory security and session management.

## Features

- **Directory Security**: Only allows Claude Code to run in pre-configured directories
- **Session Management**: Prompts to continue previous session or start fresh
- **Flexible Configuration**: Supports environment variables and JSON configuration file
- **Cross-platform**: Works on Linux, macOS, and other POSIX-compatible systems

## Installation

### From GitHub Releases (Recommended)

Download the latest release for your platform from the [releases page](https://github.com/23prime/claude-launcher/releases).

#### Automated Installation (Linux/macOS)

Install latest version.

```bash
curl -fsSL https://raw.githubusercontent.com/23prime/claude-launcher/main/install.sh | bash
```

Or with wget:

```bash
wget -qO- https://raw.githubusercontent.com/23prime/claude-launcher/main/install.sh | bash
```

#### Manual Installation

**Linux/macOS:**

```bash
# Download the appropriate binary for your platform
# Linux amd64
curl -LO https://github.com/23prime/claude-launcher/releases/latest/download/claude-launcher-linux-amd64.tar.gz
tar -xzf claude-launcher-linux-amd64.tar.gz

# macOS amd64 (Intel)
curl -LO https://github.com/23prime/claude-launcher/releases/latest/download/claude-launcher-darwin-amd64.tar.gz
tar -xzf claude-launcher-darwin-amd64.tar.gz

# macOS arm64 (Apple Silicon)
curl -LO https://github.com/23prime/claude-launcher/releases/latest/download/claude-launcher-darwin-arm64.tar.gz
tar -xzf claude-launcher-darwin-arm64.tar.gz

# Make executable and move to PATH
chmod +x claude-launcher-*
mv claude-launcher-* ~/.local/bin/claude-launcher
```

**Windows:**

Download the `.zip` file for your architecture from the releases page, extract it, and add the directory to your PATH.

### Using go install

```bash
go install github.com/23prime/claude-launcher/cmd/claude-launcher@latest
```

### From source

```bash
git clone https://github.com/23prime/claude-launcher
cd claude-launcher
task go:build
mv bin/claude-launcher ~/.local/bin/
```

## Configuration

Configure allowed directories using one of these methods:

### Method 1: Environment Variable (Priority 1)

```bash
export CLAUDE_SAFE_DIRS="$HOME/develop:$HOME/projects"
```

### Method 2: Settings File (Priority 2)

Edit `~/.claude/settings.json`:

```json
{
  "customConfig": {
    "allowedDirs": [
      "/home/user/develop",
      "/home/user/projects"
    ]
  }
}
```

## Usage

### Basic usage

```bash
# Check directory and launch Claude
claude-launcher

# Show help
claude-launcher --help

# Show configured directories
claude-launcher --show-dirs

# Pass arguments to Claude
claude-launcher --model opus
```

### Example session

```sh
$ cd ~/develop/myproject
$ claude-launcher
✓ Directory allowed

Continue previous Claude session?
  [Y/n] (default: y): y
→ Continuing previous session...
```

## Development

This project uses [Taskfile](https://taskfile.dev) and [mise](https://mise.jdx.dev) for development.

### Setup

```bash
# Install dependencies
mise install

# Run tests
task go:test

# Build
task go:build

# Run all checks
task check
```

### Available Tasks

```bash
# List all available tasks
task --list-all

# Go-specific tasks
task go:build              # Build the application
task go:test               # Run tests
task go:test-cover         # Run tests with coverage
task go:fmt                # Format code
task go:vet                # Run go vet
```

## Project Structure

```txt
claude-launcher/
├── cmd/
│   └── claude-launcher/   # Main application entry point
├── internal/
│   ├── config/            # Configuration loading
│   ├── security/          # Directory access checking
│   ├── session/           # Session continuation prompts
│   ├── launcher/          # Claude Code execution
│   └── ui/                # User interface (colors, messages)
├── docs/
│   ├── specification.md   # Detailed specification
│   └── implementation-plan.md  # Implementation plan
└── tasks/
    └── GoTasks.yml        # Go build tasks
```

## Documentation

- [Specification](docs/specification.md) - Detailed feature specification
- [Implementation Plan](docs/implementation-plan.md) - Go implementation guide
- [CLAUDE.md](CLAUDE.md) - Project overview for Claude Code

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
