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
				"CLAUDE_CODE_ENABLE_TELEMETRY": "1",
				"OTEL_EXPORTER_OTLP_ENDPOINT":  "http://localhost:4317",
				"OTEL_EXPORTER_OTLP_PROTOCOL":  "grpc",
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
