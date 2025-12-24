# claude-launcher Specification

## Overview

`claude-launcher` is a comprehensive launcher tool for Claude Code. It provides a safe and user-friendly environment for launching Claude Code by integrating directory security and session management features.

## Core Features

### 1. Directory Security Check

#### 1.1 Purpose

Ensures Claude Code can only be executed within allowed directories, preventing unintended execution in unauthorized locations.

#### 1.2 Behavior

- Checks if the current directory is in the allowed list
- Allows execution if the current directory is an allowed directory itself or a subdirectory of one
- Resolves symlinks to actual paths using `realpath`
- Expands tilde (`~`) to home directory

#### 1.3 Allowed List Configuration

Configuration is loaded in the following priority order:

**Priority 1: Environment Variable `CLAUDE_SAFE_DIRS`**

- Specify multiple directories separated by colons (`:`)
- Example: `export CLAUDE_SAFE_DIRS="$HOME/projects:$HOME/work"`

**Priority 2: `~/.claude/settings.json`**

- Specified in the `customConfig.allowedDirs` array
- JSON format
- Example:

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

#### 1.4 Configuration Parsing

- Uses `jq` for accurate parsing when available
- Falls back to `grep` when `jq` is unavailable
- Displays error message and exits if no configuration is found

### 2. Session Management

#### 2.1 Purpose

Allows users to choose whether to continue the previous session or start a new one.

#### 2.2 Behavior

- Displays confirmation prompt after passing directory check
- Prompt: `Continue previous Claude session? [Y/n] (default: y):`
- Default is to continue (Enter key or `y`/`Y` continues)
- `n` or `no` (case-insensitive) starts a new session

#### 2.3 Claude Code Launch

- **Continue**: `claude --continue [additional arguments...]`
- **New**: `claude [additional arguments...]`
- Command-line arguments are passed through to Claude

### 3. Command-Line Interface

#### 3.1 Basic Syntax

```bash
claude-launcher [OPTIONS] [CLAUDE_ARGUMENTS...]
```

#### 3.2 Options

| Option | Short | Description |
| -------- | ------- | ------------- |
| `--help` | `-h` | Display help message and exit |
| `--show-dirs` | `-l` | Display list of allowed directories and exit |

#### 3.3 Additional Arguments

- Arguments other than options are passed directly to Claude Code
- Example: `claude-launcher --model opus` → `claude --continue --model opus`

### 4. User Feedback

#### 4.1 Color Output

- **GREEN**: Success messages
  - `✓ Directory allowed`
  - `→ Continuing previous session...`
  - `→ Starting new session...`
- **RED**: Error messages
  - `✗ Access denied`
  - `Error: No allowed directories configured`
- **YELLOW**: Confirmation prompts
  - `Continue previous Claude session?`

#### 4.2 Error Messages

**Directory Access Denied:**

```txt
✗ Access denied

Current directory: /path/to/current/dir

Claude Code is not allowed to run in this directory.
Allowed directories:
  - /home/user/projects
  - /home/user/work
```

**Configuration Not Found:**

```txt
Error: No allowed directories configured

Please set allowed directories using one of these methods:

1. Environment variable (colon-separated):
   export CLAUDE_SAFE_DIRS="$HOME/projects:$HOME/work"

2. Edit ~/.claude/settings.json:
   {"customConfig": {"allowedDirs": ["/home/user/projects"]}}
```

### 5. Execution Flow

```txt
Start
  ↓
[--help/-h?] → Show help → Exit
  ↓ No
[--show-dirs/-l?] → Show directory list → Exit
  ↓ No
Load allowed directory configuration
  ↓
[Configuration exists?]
  ↓ No → Show error message → Exit
  ↓ Yes
Check current directory
  ↓
[Allowed?]
  ↓ No → Show access denied message → Exit
  ↓ Yes
Show "✓ Directory allowed"
  ↓
Display session continuation prompt
  ↓
[Continue?]
  ↓ Yes → Execute claude --continue [args...]
  ↓ No → Execute claude [args...]
  ↓
Exit
```

## Technical Requirements

### Dependencies

- **Required**: `bash`, `grep`, `pwd`
- **Recommended**: `realpath` (for symlink resolution)
- **Optional**: `jq` (for JSON parsing, fallback functionality available)

### Compatibility

- POSIX-compatible systems (Linux, macOS, WSL, etc.)
- Bash 4.0 or higher recommended

## Considerations for Go Implementation

### Parts Implementable with Standard Library

- File path operations: `path/filepath`
- JSON parsing: `encoding/json`
- Environment variable reading: `os.Getenv`
- Standard I/O: `fmt`, `bufio`
- Process execution: `os/exec`

### Additional Features to Consider

- Configuration file validation
- More detailed logging functionality
- Automatic configuration file generation
- Interactive configuration wizard
- Cross-platform color output support

### Platform-Specific Processing

- Windows support: Path separator, home directory retrieval
- Color output: Windows Console API support
