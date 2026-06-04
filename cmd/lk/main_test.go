package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunListsKinds(t *testing.T) {
	var out bytes.Buffer
	if code := run([]string{"kinds"}, &out, &bytes.Buffer{}); code != 0 {
		t.Fatalf("unexpected exit code: %d", code)
	}
	if !strings.Contains(out.String(), "documents") {
		t.Fatalf("expected kinds output to include documents, got %q", out.String())
	}
}

func TestRunExitCode1NoResults(t *testing.T) {
	dir := t.TempDir()
	var out, errs bytes.Buffer
	code := run([]string{"images", dir}, &out, &errs)
	if code != 1 {
		t.Fatalf("expected exit code 1 for no results, got %d", code)
	}
}

func TestRunExitCode0WithResults(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "photo.jpg"), []byte("fake"), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errs bytes.Buffer
	code := run([]string{"images", dir}, &out, &errs)
	if code != 0 {
		t.Fatalf("expected exit code 0 with results, got %d", code)
	}
}

func TestRunNullTerminatesOutput(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "photo.jpg"), []byte("fake"), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errs bytes.Buffer
	code := run([]string{"images", "--null", dir}, &out, &errs)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(out.String(), "photo.jpg\x00") {
		t.Fatalf("expected null-terminated output, got %q", out.String())
	}
}

func TestRunExecSubstitutesPath(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "photo.jpg"), []byte("fake"), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errs bytes.Buffer
	code := run([]string{"images", dir, "--exec", "echo", "{}"}, &out, &errs)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d: %s", code, errs.String())
	}
	if !strings.Contains(out.String(), "photo.jpg") {
		t.Fatalf("expected exec output with photo.jpg, got %q", out.String())
	}
}

func TestRunExecZeroMatchesExits1(t *testing.T) {
	dir := t.TempDir()
	var out, errs bytes.Buffer
	code := run([]string{"images", dir, "--exec", "echo", "{}"}, &out, &errs)
	if code != 1 {
		t.Fatalf("expected exit code 1 for exec with no matches, got %d", code)
	}
}

func TestRunExecSuppressesNormalOutput(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "photo.jpg"), []byte("fake"), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errs bytes.Buffer
	code := run([]string{"images", dir, "--exec", "echo", "{}"}, &out, &errs)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if strings.Contains(out.String(), "matched_by") {
		t.Fatalf("expected no table output with --exec, got %q", out.String())
	}
}
