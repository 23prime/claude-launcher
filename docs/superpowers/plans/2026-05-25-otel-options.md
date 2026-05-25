# OTel Options Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add per-account and global OpenTelemetry environment variable configuration to claude-launcher,
injected into Claude Code's process with shell env vars taking highest priority.

**Architecture:** Three struct additions (`Config.OtelEnv`, `Account.OtelEnv`, `LaunchOptions.OtelEnv`)
plus a helper function `buildOtelEnv` in the launcher that skips keys already present in the shell.
`main.go` merges global and account OTel maps before passing them to the launcher.

**Tech Stack:** Go standard library only — `encoding/json`, `os`, `strings`.

---

## File Map

| File | Change |
| ---- | ------ |
| `internal/config/config.go` | Add `OtelEnv map[string]string` to `Config` and `configJSON` |
| `internal/config/config_test.go` | Tests for `otelEnv` JSON parsing |
| `internal/account/config.go` | Add `OtelEnv map[string]string` to `Account` and `accountJSON` |
| `internal/account/config_test.go` | Tests for per-account `otelEnv` JSON parsing |
| `internal/launcher/launcher.go` | Add `OtelEnv` to `LaunchOptions`; extract `buildOtelEnv` helper; call it in `Launch` |
| `internal/launcher/launcher_test.go` | Tests for `buildOtelEnv` (new file) |
| `cmd/claude-launcher/main.go` | Merge global + account `OtelEnv` before creating `LaunchOptions` |

---

### Task 1: Add `OtelEnv` to `Config`

**Files:**

- Modify: `internal/config/config.go`
- Modify: `internal/config/config_test.go`

- [ ] **Step 1: Write the failing test**

Append to `internal/config/config_test.go`:

```go
func TestFileLoaderOtelEnv(t *testing.T) {
 tmpDir := t.TempDir()
 testFile := filepath.Join(tmpDir, "config.json")

 tests := []struct {
  name            string
  jsonContent     string
  wantErr         bool
  expectedOtelEnv map[string]string
 }{
  {
   name: "with otelEnv",
   jsonContent: `{
    "allowedDirs": ["/home/user/projects"],
    "otelEnv": {
     "CLAUDE_CODE_ENABLE_TELEMETRY": "1",
     "OTEL_EXPORTER_OTLP_ENDPOINT": "http://localhost:4317"
    }
   }`,
   wantErr: false,
   expectedOtelEnv: map[string]string{
    "CLAUDE_CODE_ENABLE_TELEMETRY": "1",
    "OTEL_EXPORTER_OTLP_ENDPOINT":  "http://localhost:4317",
   },
  },
  {
   name: "without otelEnv",
   jsonContent: `{
    "allowedDirs": ["/home/user/projects"]
   }`,
   wantErr:         false,
   expectedOtelEnv: nil,
  },
 }

 for _, tt := range tests {
  t.Run(tt.name, func(t *testing.T) {
   if err := os.WriteFile(testFile, []byte(tt.jsonContent), 0o644); err != nil {
    t.Fatalf("failed to create test file: %v", err)
   }

   loader := &FileLoader{Path: testFile}
   cfg, err := loader.Load()

   if (err != nil) != tt.wantErr {
    t.Errorf("FileLoader.Load() error = %v, wantErr %v", err, tt.wantErr)
    return
   }

   if !tt.wantErr {
    if len(cfg.OtelEnv) != len(tt.expectedOtelEnv) {
     t.Errorf("OtelEnv length = %d, expected %d", len(cfg.OtelEnv), len(tt.expectedOtelEnv))
     return
    }
    for k, v := range tt.expectedOtelEnv {
     if cfg.OtelEnv[k] != v {
      t.Errorf("OtelEnv[%q] = %q, expected %q", k, cfg.OtelEnv[k], v)
     }
    }
   }
  })
 }
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test -v ./internal/config/... -run TestFileLoaderOtelEnv
```

Expected: FAIL — `cfg.OtelEnv` field does not exist yet.

- [ ] **Step 3: Add `OtelEnv` to `Config` and `configJSON`**

In `internal/config/config.go`, update the two structs and `FileLoader.Load()`:

```go
// Config represents the configuration for claude-launcher
type Config struct {
 AllowedDirs []string
 OtelEnv     map[string]string
}
```

```go
// configJSON represents the structure of the config file
type configJSON struct {
 AllowedDirs []string          `json:"allowedDirs"`
 OtelEnv     map[string]string `json:"otelEnv,omitempty"`
}
```

In `FileLoader.Load()`, replace the final `return`:

```go
 return &Config{
  AllowedDirs: expandedDirs,
  OtelEnv:     cfg.OtelEnv,
 }, nil
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test -v ./internal/config/... -run TestFileLoaderOtelEnv
```

Expected: PASS

- [ ] **Step 5: Run full config tests**

```bash
go test -v ./internal/config/...
```

Expected: All PASS

- [ ] **Step 6: Commit**

```bash
git add internal/config/config.go internal/config/config_test.go
git commit -m "feat: add OtelEnv field to Config"
```

---

### Task 2: Add `OtelEnv` to `Account`

**Files:**

- Modify: `internal/account/config.go`
- Modify: `internal/account/config_test.go`

- [ ] **Step 1: Write the failing test**

Append to `internal/account/config_test.go`:

```go
func TestFileLoaderAccountOtelEnv(t *testing.T) {
 tmpDir := t.TempDir()
 testFile := filepath.Join(tmpDir, "config.json")

 tests := []struct {
  name            string
  jsonContent     string
  wantErr         bool
  expectedOtelEnv map[string]string
 }{
  {
   name: "account with otelEnv",
   jsonContent: `{
    "accounts": [
     {
      "name": "Work",
      "configDir": "/home/user/.claude-work",
      "otelEnv": {
       "OTEL_SERVICE_NAME": "claude-work",
       "OTEL_RESOURCE_ATTRIBUTES": "team.id=platform"
      }
     }
    ]
   }`,
   wantErr: false,
   expectedOtelEnv: map[string]string{
    "OTEL_SERVICE_NAME":        "claude-work",
    "OTEL_RESOURCE_ATTRIBUTES": "team.id=platform",
   },
  },
  {
   name: "account without otelEnv",
   jsonContent: `{
    "accounts": [
     {"name": "Personal", "configDir": "/home/user/.claude-personal"}
    ]
   }`,
   wantErr:         false,
   expectedOtelEnv: nil,
  },
 }

 for _, tt := range tests {
  t.Run(tt.name, func(t *testing.T) {
   if err := os.WriteFile(testFile, []byte(tt.jsonContent), 0o644); err != nil {
    t.Fatalf("failed to create test file: %v", err)
   }

   loader := &FileLoader{Path: testFile}
   cfg, err := loader.Load()

   if (err != nil) != tt.wantErr {
    t.Errorf("FileLoader.Load() error = %v, wantErr %v", err, tt.wantErr)
    return
   }

   if !tt.wantErr {
    acc := cfg.Accounts[0]
    if len(acc.OtelEnv) != len(tt.expectedOtelEnv) {
     t.Errorf("OtelEnv length = %d, expected %d", len(acc.OtelEnv), len(tt.expectedOtelEnv))
     return
    }
    for k, v := range tt.expectedOtelEnv {
     if acc.OtelEnv[k] != v {
      t.Errorf("OtelEnv[%q] = %q, expected %q", k, acc.OtelEnv[k], v)
     }
    }
   }
  })
 }
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test -v ./internal/account/... -run TestFileLoaderAccountOtelEnv
```

Expected: FAIL — `acc.OtelEnv` field does not exist yet.

- [ ] **Step 3: Add `OtelEnv` to `Account` and `accountJSON`**

In `internal/account/config.go`, update the structs:

```go
// Account represents a Claude account configuration
type Account struct {
 Name      string
 ConfigDir string
 OtelEnv   map[string]string
}
```

```go
// accountJSON represents the account structure in JSON
type accountJSON struct {
 Name      string            `json:"name"`
 ConfigDir string            `json:"configDir"`
 OtelEnv   map[string]string `json:"otelEnv,omitempty"`
}
```

In `FileLoader.Load()`, update the `accounts = append(...)` call:

```go
  accounts = append(accounts, Account{
   Name:      acc.Name,
   ConfigDir: expandedDir,
   OtelEnv:   acc.OtelEnv,
  })
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test -v ./internal/account/... -run TestFileLoaderAccountOtelEnv
```

Expected: PASS

- [ ] **Step 5: Run full account tests**

```bash
go test -v ./internal/account/...
```

Expected: All PASS

- [ ] **Step 6: Commit**

```bash
git add internal/account/config.go internal/account/config_test.go
git commit -m "feat: add OtelEnv field to Account"
```

---

### Task 3: Add OTel injection to `Launcher`

**Files:**

- Modify: `internal/launcher/launcher.go`
- Create: `internal/launcher/launcher_test.go`

- [ ] **Step 1: Write the failing test**

Create `internal/launcher/launcher_test.go`:

```go
package launcher

import (
 "testing"
)

func TestBuildOtelEnv(t *testing.T) {
 tests := []struct {
  name     string
  base     []string
  otelEnv  map[string]string
  wantKeys map[string]string
  skipKeys []string
 }{
  {
   name:    "injects new keys",
   base:    []string{"PATH=/usr/bin", "HOME=/home/user"},
   otelEnv: map[string]string{"OTEL_SERVICE_NAME": "claude"},
   wantKeys: map[string]string{
    "OTEL_SERVICE_NAME": "claude",
   },
  },
  {
   name:     "skips keys already in base (shell priority)",
   base:     []string{"OTEL_SERVICE_NAME=shell-value"},
   otelEnv:  map[string]string{"OTEL_SERVICE_NAME": "config-value"},
   skipKeys: []string{"OTEL_SERVICE_NAME=config-value"},
   wantKeys: map[string]string{
    "OTEL_SERVICE_NAME": "shell-value",
   },
  },
  {
   name:    "nil otelEnv returns base unchanged",
   base:    []string{"PATH=/usr/bin"},
   otelEnv: nil,
   wantKeys: map[string]string{
    "PATH": "/usr/bin",
   },
  },
  {
   name:    "empty otelEnv returns base unchanged",
   base:    []string{"PATH=/usr/bin"},
   otelEnv: map[string]string{},
   wantKeys: map[string]string{
    "PATH": "/usr/bin",
   },
  },
  {
   name: "injects multiple keys",
   base: []string{"PATH=/usr/bin"},
   otelEnv: map[string]string{
    "CLAUDE_CODE_ENABLE_TELEMETRY":  "1",
    "OTEL_EXPORTER_OTLP_ENDPOINT":   "http://localhost:4317",
    "OTEL_EXPORTER_OTLP_PROTOCOL":   "grpc",
   },
   wantKeys: map[string]string{
    "CLAUDE_CODE_ENABLE_TELEMETRY": "1",
    "OTEL_EXPORTER_OTLP_ENDPOINT":  "http://localhost:4317",
    "OTEL_EXPORTER_OTLP_PROTOCOL":  "grpc",
   },
  },
 }

 for _, tt := range tests {
  t.Run(tt.name, func(t *testing.T) {
   result := buildOtelEnv(tt.base, tt.otelEnv)

   // Build a lookup map from result for easy assertion
   got := make(map[string]string, len(result))
   for _, e := range result {
    parts := splitEnv(e)
    got[parts[0]] = parts[1]
   }

   for k, v := range tt.wantKeys {
    if got[k] != v {
     t.Errorf("env[%q] = %q, expected %q", k, got[k], v)
    }
   }

   for _, entry := range tt.skipKeys {
    for _, r := range result {
     if r == entry {
      t.Errorf("result contains %q, which should have been skipped", entry)
     }
    }
   }
  })
 }
}

// splitEnv splits "KEY=VALUE" into ["KEY", "VALUE"]
func splitEnv(e string) [2]string {
 for i := 0; i < len(e); i++ {
  if e[i] == '=' {
   return [2]string{e[:i], e[i+1:]}
  }
 }
 return [2]string{e, ""}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test -v ./internal/launcher/... -run TestBuildOtelEnv
```

Expected: FAIL — `buildOtelEnv` is not defined.

- [ ] **Step 3: Implement `buildOtelEnv` and update `LaunchOptions`**

In `internal/launcher/launcher.go`, add `"strings"` to the import block and make these changes:

```go
import (
 "fmt"
 "os"
 "os/exec"
 "strings"
)
```

Add `OtelEnv` to `LaunchOptions`:

```go
// LaunchOptions contains options for launching Claude
type LaunchOptions struct {
 Continue  bool
 Args      []string
 ConfigDir string
 OtelEnv   map[string]string
}
```

Add the `buildOtelEnv` helper function (unexported, package-level):

```go
// buildOtelEnv merges otelEnv into base, skipping keys already present in base.
// Shell env vars (base) take highest priority.
func buildOtelEnv(base []string, otelEnv map[string]string) []string {
 if len(otelEnv) == 0 {
  return base
 }

 existing := make(map[string]bool, len(base))
 for _, e := range base {
  key := strings.SplitN(e, "=", 2)[0]
  existing[key] = true
 }

 result := make([]string, len(base), len(base)+len(otelEnv))
 copy(result, base)

 for k, v := range otelEnv {
  if !existing[k] {
   result = append(result, k+"="+v)
  }
 }

 return result
}
```

In `Launch()`, replace `cmd.Env = os.Environ()` and the existing `ConfigDir` block with:

```go
 cmd.Env = buildOtelEnv(os.Environ(), opts.OtelEnv)

 if opts.ConfigDir != "" {
  cmd.Env = append(cmd.Env, "CLAUDE_CONFIG_DIR="+opts.ConfigDir)
 }
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test -v ./internal/launcher/... -run TestBuildOtelEnv
```

Expected: PASS

- [ ] **Step 5: Run all tests**

```bash
go test -v ./...
```

Expected: All PASS

- [ ] **Step 6: Commit**

```bash
git add internal/launcher/launcher.go internal/launcher/launcher_test.go
git commit -m "feat: add OtelEnv injection to launcher"
```

---

### Task 4: Wire OTel merge in `main.go`

**Files:**

- Modify: `cmd/claude-launcher/main.go`

- [ ] **Step 1: Add OTel merge and pass to `LaunchOptions`**

In `cmd/claude-launcher/main.go`, find the section that builds `launchOpts` (around line 158) and add
the merge block immediately before it:

```go
 // Merge OTel env: global config first, account overrides second.
 // Shell env vars take priority at injection time (handled in launcher).
 otelEnv := make(map[string]string)
 for k, v := range cfg.OtelEnv {
  otelEnv[k] = v
 }
 if selectedAccount != nil {
  for k, v := range selectedAccount.OtelEnv {
   otelEnv[k] = v
  }
 }

 // Launch Claude
 l := launcher.NewLauncher()
 launchOpts := launcher.LaunchOptions{
  Continue:  shouldContinue,
  Args:      flag.Args(),
  ConfigDir: configDir,
  OtelEnv:   otelEnv,
 }
```

Remove the old standalone `// Launch Claude` comment and `launchOpts` block that was there before.

- [ ] **Step 2: Run all tests**

```bash
go test -v ./...
```

Expected: All PASS

- [ ] **Step 3: Run fix-and-check**

```bash
mise run fix-and-check
```

Expected: All checks pass.

- [ ] **Step 4: Commit**

```bash
git add cmd/claude-launcher/main.go
git commit -m "feat: wire OTel env merge in main"
```
