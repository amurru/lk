package integration

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func buildLK(t *testing.T) (string, string) {
	t.Helper()
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	return root, filepath.Join(root, "tests", "integration")
}

func TestCLIListsImages(t *testing.T) {
	_, dir := buildLK(t)
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, "photo.jpg"), []byte("fake"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "note.txt"), []byte("hello"), 0o600); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "run", "../../cmd/lk", "images", tmp)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("command failed: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "photo.jpg") {
		t.Fatalf("expected photo.jpg in output, got %q", stdout.String())
	}
	if strings.Contains(stdout.String(), "note.txt") {
		t.Fatalf("did not expect note.txt in output, got %q", stdout.String())
	}
}

func TestCLIListsFonts(t *testing.T) {
	_, dir := buildLK(t)
	cmd := exec.Command("go", "run", "../../cmd/lk", "fonts", "testdata/fonts")
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("command failed: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "test.ttf") {
		t.Fatalf("expected test.ttf in output, got %q", stdout.String())
	}
}

func TestCLIListsExecutables(t *testing.T) {
	_, dir := buildLK(t)
	cmd := exec.Command("go", "run", "../../cmd/lk", "executables", "testdata/executables")
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("command failed: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "hello_elf") {
		t.Fatalf("expected hello_elf in output, got %q", stdout.String())
	}
}

func TestCLIListsModels(t *testing.T) {
	_, dir := buildLK(t)
	cmd := exec.Command("go", "run", "../../cmd/lk", "models", "testdata/models")
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("command failed: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "test.glb") {
		t.Fatalf("expected test.glb in output, got %q", stdout.String())
	}
}

func TestCLIListsDatabases(t *testing.T) {
	_, dir := buildLK(t)
	cmd := exec.Command("go", "run", "../../cmd/lk", "databases", "testdata/databases")
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("command failed: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "test.db") {
		t.Fatalf("expected test.db in output, got %q", stdout.String())
	}
}

func TestCLIListsDiskimages(t *testing.T) {
	_, dir := buildLK(t)
	cmd := exec.Command("go", "run", "../../cmd/lk", "diskimages", "testdata/diskimages")
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("command failed: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "test.iso") {
		t.Fatalf("expected test.iso in output, got %q", stdout.String())
	}
}
