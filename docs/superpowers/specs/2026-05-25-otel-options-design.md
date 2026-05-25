# OTel Options Configuration Design

## Overview

Add OpenTelemetry (OTel) environment variable configuration to claude-launcher, supporting
both global and per-account settings. This allows users to configure telemetry for Claude Code
without setting environment variables in their shell.

## Goals

- Support all Claude Code OTel environment variables (`OTEL_*`, `CLAUDE_CODE_*`)
- Per-account OTel settings (different telemetry per account)
- Global OTel settings as fallback
- Shell environment variables take highest priority (no override)

## Priority Order

Shell env var > Account `otelEnv` > Global `otelEnv`

If a variable is already set in the shell, the config file value is ignored.

## Configuration Schema

### `~/.config/claude-launcher/config.json`

```json
{
  "allowedDirs": ["~/develop"],
  "otelEnv": {
    "CLAUDE_CODE_ENABLE_TELEMETRY": "1",
    "OTEL_METRICS_EXPORTER": "otlp",
    "OTEL_EXPORTER_OTLP_PROTOCOL": "grpc",
    "OTEL_EXPORTER_OTLP_ENDPOINT": "http://localhost:4317",
    "OTEL_EXPORTER_OTLP_HEADERS": "Authorization=Bearer your-token",
    "OTEL_SERVICE_NAME": "claude"
  },
  "accounts": [
    {
      "name": "Work",
      "configDir": "~/.claude-work",
      "otelEnv": {
        "OTEL_SERVICE_NAME": "claude-work",
        "OTEL_RESOURCE_ATTRIBUTES": "team.id=platform,cost_center=eng-123"
      }
    }
  ]
}
```

The `otelEnv` map accepts any key-value pair. All variables documented by Claude Code are
supported, including:

| Variable | Description |
| -------- | ----------- |
| `CLAUDE_CODE_ENABLE_TELEMETRY` | **Required** to enable telemetry (value: `1`) |
| `OTEL_METRICS_EXPORTER` | Metrics exporter (`otlp`, `prometheus`, `console`, `none`) |
| `OTEL_LOGS_EXPORTER` | Logs exporter (`otlp`, `console`, `none`) |
| `OTEL_EXPORTER_OTLP_PROTOCOL` | Protocol (`grpc`, `http/protobuf`, `http/json`) |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | OTLP collector endpoint |
| `OTEL_EXPORTER_OTLP_HEADERS` | Auth headers (e.g. `Authorization=Bearer token`) |
| `OTEL_METRIC_EXPORT_INTERVAL` | Metrics export interval in ms (default: `60000`) |
| `OTEL_LOGS_EXPORT_INTERVAL` | Logs export interval in ms (default: `5000`) |
| `OTEL_LOG_USER_PROMPTS` | Log user prompt content (`1` to enable) |
| `OTEL_LOG_TOOL_DETAILS` | Log tool parameters (`1` to enable) |
| `OTEL_LOG_TOOL_CONTENT` | Log tool input/output in span events (`1` to enable) |
| `OTEL_LOG_RAW_API_BODIES` | Log raw API bodies (`1` or `file:<dir>`) |
| `OTEL_METRICS_INCLUDE_SESSION_ID` | Include session ID in metrics (default: `true`) |
| `OTEL_METRICS_INCLUDE_VERSION` | Include app version in metrics (default: `false`) |
| `OTEL_METRICS_INCLUDE_ACCOUNT_UUID` | Include account UUID in metrics (default: `true`) |
| `CLAUDE_CODE_ENHANCED_TELEMETRY_BETA` | Enable distributed tracing beta (`1`) |
| `OTEL_TRACES_EXPORTER` | Traces exporter (`otlp`, `console`, `none`) |
| `OTEL_EXPORTER_OTLP_TRACES_PROTOCOL` | Traces-specific protocol |
| `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT` | Traces-specific endpoint |
| `OTEL_TRACES_EXPORT_INTERVAL` | Traces export interval in ms |
| `OTEL_RESOURCE_ATTRIBUTES` | Resource attributes (`key=value,key2=value2`) |
| `CLAUDE_CODE_OTEL_HEADERS_HELPER_DEBOUNCE_MS` | Dynamic headers refresh interval in ms |

## Code Changes

### `internal/config/config.go`

Add `OtelEnv` field to `Config` and `configJSON`:

```go
type Config struct {
    AllowedDirs []string
    OtelEnv     map[string]string
}

type configJSON struct {
    AllowedDirs []string          `json:"allowedDirs"`
    OtelEnv     map[string]string `json:"otelEnv,omitempty"`
}
```

### `internal/account/config.go`

Add `OtelEnv` field to `Account` and `accountJSON`:

```go
type Account struct {
    Name      string
    ConfigDir string
    OtelEnv   map[string]string
}

type accountJSON struct {
    Name      string            `json:"name"`
    ConfigDir string            `json:"configDir"`
    OtelEnv   map[string]string `json:"otelEnv,omitempty"`
}
```

### `internal/launcher/launcher.go`

Add `OtelEnv` to `LaunchOptions` and inject env vars with shell priority check:

```go
type LaunchOptions struct {
    Continue  bool
    Args      []string
    ConfigDir string
    OtelEnv   map[string]string
}
```

In `Launch()`, after setting `cmd.Env = os.Environ()`:

```go
existing := make(map[string]bool)
for _, e := range os.Environ() {
    key := strings.SplitN(e, "=", 2)[0]
    existing[key] = true
}
for k, v := range opts.OtelEnv {
    if !existing[k] {
        cmd.Env = append(cmd.Env, k+"="+v)
    }
}
```

### `cmd/claude-launcher/main.go`

Merge global and account OTel settings before launching:

```go
// Merge OTel env: global first, then account overrides
otelEnv := make(map[string]string)
for k, v := range cfg.OtelEnv {
    otelEnv[k] = v
}
if selectedAccount != nil {
    for k, v := range selectedAccount.OtelEnv {
        otelEnv[k] = v
    }
}

launchOpts := launcher.LaunchOptions{
    Continue:  shouldContinue,
    Args:      flag.Args(),
    ConfigDir: configDir,
    OtelEnv:   otelEnv,
}
```

## Execution Flow

```txt
Load config (global OtelEnv)
  ↓
Select account (account OtelEnv)
  ↓
Merge: global OtelEnv ← overridden by account OtelEnv
  ↓
Launch: inject merged OtelEnv, skip keys already in shell env
```

## Testing

- `internal/config/config_test.go`: `otelEnv` field parsing from JSON
- `internal/account/config_test.go`: account `otelEnv` field parsing
- `internal/launcher/launcher_test.go`: env var injection, shell priority check
- `cmd/claude-launcher/main_test.go` (if exists): global + account merge logic
