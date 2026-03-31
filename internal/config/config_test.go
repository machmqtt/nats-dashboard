package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadValid(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "config.yaml")
	os.WriteFile(p, []byte(`
listen: ":9090"
poll_interval: 10s
session_secret: "test-secret"
data_dir: "./testdata"
environments:
  - name: dev
    servers:
      - url: "http://localhost:8222"
`), 0o644)

	cfg, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Listen != ":9090" {
		t.Errorf("listen = %q, want :9090", cfg.Listen)
	}
	if cfg.PollInterval.Seconds() != 10 {
		t.Errorf("poll_interval = %v, want 10s", cfg.PollInterval)
	}
	if len(cfg.Environments) != 1 {
		t.Fatalf("environments = %d, want 1", len(cfg.Environments))
	}
	if cfg.Environments[0].Name != "dev" {
		t.Errorf("env name = %q, want dev", cfg.Environments[0].Name)
	}
}

func TestLoadMissingSecret(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "config.yaml")
	os.WriteFile(p, []byte(`
environments:
  - name: dev
    servers:
      - url: "http://localhost:8222"
`), 0o644)

	_, err := Load(p)
	if err == nil {
		t.Fatal("expected error for missing session_secret")
	}
}

func TestLoadNoEnvironments(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "config.yaml")
	os.WriteFile(p, []byte(`
session_secret: "test"
`), 0o644)

	_, err := Load(p)
	if err == nil {
		t.Fatal("expected error for no environments")
	}
}

func TestLoadDefaults(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "config.yaml")
	os.WriteFile(p, []byte(`
session_secret: "test"
environments:
  - name: dev
    servers:
      - url: "http://localhost:8222"
`), 0o644)

	cfg, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Listen != ":8080" {
		t.Errorf("default listen = %q, want :8080", cfg.Listen)
	}
	if cfg.DataDir != "./data" {
		t.Errorf("default data_dir = %q, want ./data", cfg.DataDir)
	}
}
