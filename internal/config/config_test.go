package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadReadsMetricsPassword(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(configFile, []byte("metrics:\n  enabled: true\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("APP_METRICS_PASSWORD", "test-secret")
	cfg, err := Load(configFile)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Metrics.Password != "test-secret" {
		t.Fatalf("metrics password = %q, want environment variable value", cfg.Metrics.Password)
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestValidateRequiresMetricsCredentials(t *testing.T) {
	cfg := &Config{
		Server:   Server{Port: 8080},
		Database: Database{Host: "localhost", Name: "app"},
		Metrics:  Metrics{Enabled: true},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want missing metrics credentials error")
	}
}
