# claude-launcher Go Implementation Plan

## 1. Project Structure

```txt
claude-launcher/
├── cmd/
│   └── claude-launcher/
│       └── main.go              # Entry point
├── internal/
│   ├── config/
│   │   ├── config.go            # Configuration loading and management
│   │   └── config_test.go
│   ├── security/
│   │   ├── directory.go         # Directory check functionality
│   │   └── directory_test.go
│   ├── session/
│   │   ├── prompt.go            # Session continuation prompt
│   │   └── prompt_test.go
│   ├── launcher/
│   │   ├── launcher.go          # Claude Code launch process
│   │   └── launcher_test.go
│   └── ui/
│       ├── color.go             # Color output
│       └── messages.go          # Message definitions
├── go.mod
├── go.sum
├── README.md
├── tasks/
│   └── GoTasks.yml              # Build and test tasks
└── .gitignore
```

## 2. Package Design

### 2.1 `cmd/claude-launcher`

**Responsibility**: Application entry point

**Main Processing**:

- Parse command-line arguments
- Initialize and call each package
- Overall flow control

### 2.2 `internal/config`

**Responsibility**: Configuration loading and management

**Key Types**:

```go
type Config struct {
    AllowedDirs []string
}

type Loader interface {
    Load() (*Config, error)
}

// Load from environment variable
type EnvLoader struct{}

// Load from settings.json
type FileLoader struct {
    Path string
}

// Manage multiple Loaders with priority
type ChainLoader struct {
    Loaders []Loader
}
```

**Key Functions**:

- `LoadConfig() (*Config, error)`: Load configuration according to priority
- `ExpandPath(path string) (string, error)`: Expand `~` to home directory

### 2.3 `internal/security`

**Responsibility**: Directory access check

**Key Types**:

```go
type DirectoryChecker struct {
    AllowedDirs []string
}
```

**Key Functions**:

- `NewDirectoryChecker(allowedDirs []string) *DirectoryChecker`
- `IsAllowed(currentDir string) (bool, error)`: Check if current directory is allowed
- `ResolvePath(path string) (string, error)`: Resolve symlinks

### 2.4 `internal/session`

**Responsibility**: Confirm session continuation

**Key Types**:

```go
type Prompter interface {
    AskContinue() (bool, error)
}

type InteractivePrompter struct {
    Reader io.Reader
    Writer io.Writer
}
```

**Key Functions**:

- `AskContinue() (bool, error)`: Ask user about session continuation

### 2.5 `internal/launcher`

**Responsibility**: Launch Claude Code

**Key Types**:

```go
type Launcher struct {
    ClaudePath string
}

type LaunchOptions struct {
    Continue bool
    Args     []string
}
```

**Key Functions**:

- `Launch(opts LaunchOptions) error`: Launch Claude Code

### 2.6 `internal/ui`

**Responsibility**: User interface (color output, messages)

**Key Types**:

```go
type Color int

const (
    ColorRed Color = iota
    ColorGreen
    ColorYellow
)

type Printer struct {
    Writer     io.Writer
    ColorEnabled bool
}
```

**Key Functions**:

- `Print(color Color, format string, args ...interface{})`: Color output
- `ShowAllowedDirs(dirs []string)`: Display list of allowed directories
- `ShowAccessDenied(currentDir string, allowedDirs []string)`: Access denied message
- `ShowConfigError()`: Configuration error message

## 3. Implementation Steps

### Phase 1: Basic Structure

1. **Project Initialization**
   - `go mod init github.com/23prime/claude-launcher`
   - Create basic directory structure

2. **Configuration Loading**
   - Implement `internal/config` package
   - Load from environment variable
   - Load from JSON file
   - Create tests

3. **Directory Check**
   - Implement `internal/security` package
   - Path resolution and matching
   - Create tests

### Phase 2: User Interface

1. **UI Features**
   - Implement `internal/ui` package
   - Color output (consider using `fatih/color`)
   - Message templates

2. **Session Management**
   - Implement `internal/session` package
   - Interactive prompt
   - Create tests (using mocks)

### Phase 3: Launch Process and CLI

1. **Launch Functionality**
   - Implement `internal/launcher` package
   - Claude launch using `os/exec`
   - Create tests

2. **CLI Implementation**
   - Implement `cmd/claude-launcher/main.go`
   - Use `flag` package or `cobra`
   - Help message
   - `--show-dirs` option

### Phase 4: Integration and Release

1. **Integration Testing**
   - End-to-end tests
   - Error handling verification

2. **Documentation**
   - Update README.md
   - Add usage examples

3. **Build and Release**
    - Create Taskfile tasks
    - GoReleaser configuration (for cross-compilation)

## 4. Dependencies

### Recommended External Libraries

```go
require (
    github.com/fatih/color v1.16.0         // Color output
    github.com/spf13/cobra v1.8.0          // CLI framework (optional)
)
```

### Standard Library

- `encoding/json`: JSON parsing
- `os`: Environment variables, file operations
- `os/exec`: Process execution
- `path/filepath`: Path operations
- `flag`: Command-line argument parsing
- `bufio`: Standard input reading
- `fmt`: Formatted output
- `io`: Interfaces

## 5. Testing Strategy

### 5.1 Unit Tests

- Place `*_test.go` files in each package
- Test coverage goal: 80% or higher
- Utilize table-driven tests

### 5.2 Test Case Examples

**config package:**

- Successful load from environment variable
- Successful load from JSON file
- Priority when both exist
- Error when no configuration exists
- Handling invalid JSON

**security package:**

- Execution in allowed directory
- Execution in subdirectory
- Execution in non-allowed directory
- Symlink resolution
- Tilde expansion

**session package:**

- Continue with "y" input
- New session with "n" input
- Continue with Enter (default)

### 5.3 Integration Tests

- Use temporary directories
- Create mock Claude command
- End-to-end scenarios

## 6. Build and Deploy

### 6.1 Taskfile

```yaml
version: "3"

tasks:
  build:
    desc: Build the application
    cmds:
      - go build -o bin/claude-launcher ./cmd/claude-launcher

  test:
    desc: Run tests
    cmds:
      - go test -v ./...

  clean:
    desc: Clean build artifacts
    cmds:
      - rm -rf bin/

  install:
    desc: Install binary
    cmds:
      - go install ./cmd/claude-launcher
```

### 6.2 Cross-Compilation

Build for multiple platforms using GoReleaser:

- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

### 6.3 Installation Methods

#### Method 1: go install

```bash
go install github.com/23prime/claude-launcher/cmd/claude-launcher@latest
```

#### Method 2: Download Binary

Download platform-specific binary from GitHub Releases

#### Method 3: Build from Source

```bash
git clone https://github.com/23prime/claude-launcher
cd claude-launcher
task go:build
mv bin/claude-launcher ~/.local/bin/
```

## 7. Future Enhancements

### 7.1 Additional Feature Candidates

- **Configuration File Generation Command**
  - `claude-launcher init` to interactively generate configuration file

- **Detailed Logging**
  - `--verbose` flag for debug information
  - Log to file

- **Configuration Validation**
  - `claude-launcher validate` to check configuration validity

- **Auto-completion**
  - Generate completion scripts for Bash/Zsh

### 7.2 Performance Optimization

- Configuration file caching
- Parallel processing (multiple directory checks)

### 7.3 Security Enhancements

- Directory whitelist encryption option
- Audit logging functionality

## 8. Implementation Notes

### 8.1 Error Handling

- Add appropriate context to all errors
- User-friendly error messages
- Distinguish between recoverable and fatal errors

### 8.2 Cross-Platform Support

- Use `filepath.Separator` for path separators
- Use `os.UserHomeDir()` for home directory
- Use Windows-compatible library for color output

### 8.3 Backward Compatibility

- Maintain same configuration file format as shell script version
- Don't change environment variable names

### 8.4 Code Style

- Use `gofmt`, `golint`, `go vet`
- Static analysis with `golangci-lint`
- Clear function and variable names
- Appropriate comments (especially for public APIs)

## 9. Implementation Schedule Overview

Implementation will proceed in stages, with testing and verification at each phase:

1. **Phase 1**: Basic features (configuration loading, directory check)
2. **Phase 2**: UI and session management
3. **Phase 3**: Launch process and CLI integration
4. **Phase 4**: Testing complete and release preparation

Review points will be set at the end of each phase, with design adjustments as needed.
