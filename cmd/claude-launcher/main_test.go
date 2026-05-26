package main

import (
	"testing"

	"github.com/23prime/claude-launcher/internal/account"
	"github.com/23prime/claude-launcher/internal/config"
)

func TestBuildLaunchOtelEnv_NoOtel(t *testing.T) {
	cfg := &config.Config{
		OtelEnv: map[string]string{
			"OTEL_METRICS_EXPORTER":   "otlp",
			"CLAUDE_CODE_ENABLE_OTEL": "1",
		},
	}
	acc := &account.Account{
		Name:      "test",
		ConfigDir: "/tmp",
		OtelEnv:   map[string]string{"OTEL_SERVICE_NAME": "test"},
	}

	result := buildLaunchOtelEnv(cfg, acc, true)

	if len(result) != 0 {
		t.Errorf("expected empty map when noOtel=true, got %v", result)
	}
}

func TestBuildLaunchOtelEnv_WithOtel(t *testing.T) {
	cfg := &config.Config{
		OtelEnv: map[string]string{
			"OTEL_METRICS_EXPORTER": "otlp",
		},
	}
	acc := &account.Account{
		Name:      "test",
		ConfigDir: "/tmp",
		OtelEnv:   map[string]string{"OTEL_SERVICE_NAME": "test"},
	}

	result := buildLaunchOtelEnv(cfg, acc, false)

	if result["OTEL_METRICS_EXPORTER"] != "otlp" {
		t.Errorf("expected OTEL_METRICS_EXPORTER=otlp, got %v", result["OTEL_METRICS_EXPORTER"])
	}
	if result["OTEL_SERVICE_NAME"] != "test" {
		t.Errorf("expected OTEL_SERVICE_NAME=test, got %v", result["OTEL_SERVICE_NAME"])
	}
}

func TestBuildLaunchOtelEnv_NilAccount(t *testing.T) {
	cfg := &config.Config{
		OtelEnv: map[string]string{
			"OTEL_METRICS_EXPORTER": "otlp",
		},
	}

	result := buildLaunchOtelEnv(cfg, nil, false)

	if result["OTEL_METRICS_EXPORTER"] != "otlp" {
		t.Errorf("expected OTEL_METRICS_EXPORTER=otlp, got %v", result["OTEL_METRICS_EXPORTER"])
	}
}
