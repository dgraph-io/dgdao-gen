package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveDefaults_SchemaSubdirPresent(t *testing.T) {
	cwd := t.TempDir()
	mustMkdir(t, filepath.Join(cwd, "schema"))

	cfg := resolveDefaults(cwd, defaults{})
	if cfg.SchemaDir != filepath.Join(cwd, "schema") {
		t.Fatalf("expected -schema-dir = CWD/schema, got %q", cfg.SchemaDir)
	}
	if cfg.EntityDir != cwd {
		t.Fatalf("expected -entity-dir = CWD (schema is in subdir), got %q", cfg.EntityDir)
	}
}

func TestResolveDefaults_SchemaLocal(t *testing.T) {
	cwd := t.TempDir() // no ./schema/ subdir

	cfg := resolveDefaults(cwd, defaults{})
	if cfg.SchemaDir != cwd {
		t.Fatalf("expected -schema-dir = CWD, got %q", cfg.SchemaDir)
	}
	expectedEntity := filepath.Join(cwd, "entity")
	if cfg.EntityDir != expectedEntity {
		t.Fatalf("expected -entity-dir = CWD/entity, got %q", cfg.EntityDir)
	}
}

func TestResolveDefaults_ExplicitSchemaDirEqualsCWD(t *testing.T) {
	cwd := t.TempDir()
	mustMkdir(t, filepath.Join(cwd, "schema")) // present but should be IGNORED since explicit flag given

	cfg := resolveDefaults(cwd, defaults{schemaDirExplicit: cwd})
	if cfg.SchemaDir != cwd {
		t.Fatalf("expected explicit -schema-dir to win, got %q", cfg.SchemaDir)
	}
	if cfg.EntityDir != filepath.Join(cwd, "entity") {
		t.Fatalf("explicit -schema-dir = CWD must trigger -entity-dir = CWD/entity, got %q", cfg.EntityDir)
	}
}

func TestResolveDefaults_ExplicitSchemaDirElsewhere(t *testing.T) {
	cwd := t.TempDir()
	mytypes := filepath.Join(cwd, "mytypes")
	mustMkdir(t, mytypes)

	cfg := resolveDefaults(cwd, defaults{schemaDirExplicit: mytypes})
	if cfg.SchemaDir != mytypes {
		t.Fatalf("expected explicit -schema-dir to win, got %q", cfg.SchemaDir)
	}
	if cfg.EntityDir != cwd {
		t.Fatalf("explicit -schema-dir != CWD must yield -entity-dir = CWD, got %q", cfg.EntityDir)
	}
}

func TestResolveDefaults_ClientDirsFollowParents(t *testing.T) {
	// When -schema-client-dir / -entity-client-dir are not explicitly set,
	// they default to the same paths as their parents.
	cwd := t.TempDir()
	mustMkdir(t, filepath.Join(cwd, "schema"))

	cfg := resolveDefaults(cwd, defaults{})
	if cfg.SchemaClientDir != cfg.SchemaDir {
		t.Fatalf("expected -schema-client-dir = -schema-dir by default, got %q vs %q", cfg.SchemaClientDir, cfg.SchemaDir)
	}
	if cfg.EntityClientDir != cfg.EntityDir {
		t.Fatalf("expected -entity-client-dir = -entity-dir by default, got %q vs %q", cfg.EntityClientDir, cfg.EntityDir)
	}
}

func TestResolveDefaults_ClientDirsExplicit(t *testing.T) {
	cwd := t.TempDir()
	mustMkdir(t, filepath.Join(cwd, "schema"))
	apiDir := filepath.Join(cwd, "api")
	mustMkdir(t, apiDir)

	cfg := resolveDefaults(cwd, defaults{
		schemaClientDirExplicit: apiDir,
		entityClientDirExplicit: apiDir,
	})
	if cfg.SchemaClientDir != apiDir {
		t.Fatalf("expected explicit -schema-client-dir to win, got %q", cfg.SchemaClientDir)
	}
	if cfg.EntityClientDir != apiDir {
		t.Fatalf("expected explicit -entity-client-dir to win, got %q", cfg.EntityClientDir)
	}
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
}

func TestResolveLicenseHeader_BothSetErrors(t *testing.T) {
	if _, err := resolveLicenseHeader("inline text", "/some/file"); err == nil {
		t.Fatal("expected an error when both -license-header and -license-header-file are set")
	}
}

func TestResolveLicenseHeader_NeitherSetReturnsEmpty(t *testing.T) {
	got, err := resolveLicenseHeader("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty header, got %q", got)
	}
}

func TestResolveLicenseHeader_InlineTextPassesThrough(t *testing.T) {
	got, err := resolveLicenseHeader("Copyright 2026 Example Corp.", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "Copyright 2026 Example Corp." {
		t.Fatalf("got %q, want inline text unchanged", got)
	}
}

func TestResolveLicenseHeader_ReadsFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "header.txt")
	want := "Copyright 2026 Example Corp.\nAll rights reserved.\n"
	if err := os.WriteFile(path, []byte(want), 0o644); err != nil {
		t.Fatalf("writing header file: %v", err)
	}

	got, err := resolveLicenseHeader("", path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestResolveLicenseHeader_MissingFileErrors(t *testing.T) {
	if _, err := resolveLicenseHeader("", filepath.Join(t.TempDir(), "missing.txt")); err == nil {
		t.Fatal("expected an error when -license-header-file does not exist")
	}
}
