package integration

import (
	"bytes"
	"fmt"
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

func TestCLIShuffleWithLimit(t *testing.T) {
	_, dir := buildLK(t)
	tmp := t.TempDir()
	// Create 10 document files.
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("doc%02d.txt", i)
		if err := os.WriteFile(filepath.Join(tmp, name), []byte("content"), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	// Run multiple times and collect subsets to verify randomness and count.
	subsets := make(map[string]bool)
	for i := 0; i < 5; i++ {
		cmd := exec.Command("go", "run", "../../cmd/lk", "documents", "--sort", "shuffle", "--limit", "3", tmp)
		cmd.Dir = dir
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			t.Fatalf("command failed: %v\nstderr: %s", err, stderr.String())
		}
		lines := strings.Fields(stdout.String())
		// In table format, header + separator + 3 data rows = at least 3 file names.
		fileCount := 0
		for _, line := range lines {
			if strings.Contains(line, "doc") && strings.Contains(line, ".txt") {
				fileCount++
			}
		}
		if fileCount != 3 {
			t.Fatalf("expected 3 document entries, got %d in output:\n%s", fileCount, stdout.String())
		}
		subsets[stdout.String()] = true
	}
	// With 10 files and limit 3, we expect different subsets across runs.
	if len(subsets) < 2 {
		t.Fatalf("expected at least 2 distinct outputs over 5 shuffle+limit runs, got %d", len(subsets))
	}
}

func TestCLIShuffleJSON(t *testing.T) {
	_, dir := buildLK(t)
	tmp := t.TempDir()
	for i := 0; i < 5; i++ {
		if err := os.WriteFile(filepath.Join(tmp, fmt.Sprintf("img%02d.png", i)), []byte("fake"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	cmd := exec.Command("go", "run", "../../cmd/lk", "images", "--sort", "shuffle", "-f", "json", tmp)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("command failed: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"kind"`) || !strings.Contains(stdout.String(), `"matched_by"`) {
		t.Fatalf("expected JSON with kind and matched_by fields, got:\n%s", stdout.String())
	}
}

func TestCLIShuffleXML(t *testing.T) {
	_, dir := buildLK(t)
	tmp := t.TempDir()
	for i := 0; i < 5; i++ {
		if err := os.WriteFile(filepath.Join(tmp, fmt.Sprintf("doc%02d.txt", i)), []byte("hello"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	cmd := exec.Command("go", "run", "../../cmd/lk", "documents", "--sort", "shuffle", "-f", "xml", tmp)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("command failed: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "<files>") || !strings.Contains(stdout.String(), "<kind>") {
		t.Fatalf("expected valid XML output, got:\n%s", stdout.String())
	}
}

func TestCLIShuffleSimple(t *testing.T) {
	_, dir := buildLK(t)
	tmp := t.TempDir()
	for i := 0; i < 5; i++ {
		if err := os.WriteFile(filepath.Join(tmp, fmt.Sprintf("vid%02d.mp4", i)), []byte("fake"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	cmd := exec.Command("go", "run", "../../cmd/lk", "videos", "--sort", "shuffle", "-f", "simple", tmp)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("command failed: %v\nstderr: %s", err, stderr.String())
	}
	lines := strings.Fields(stdout.String())
	if len(lines) != 5 {
		t.Fatalf("expected 5 simple-format lines, got %d:\n%s", len(lines), stdout.String())
	}
}

func TestCLIShuffleNull(t *testing.T) {
	_, dir := buildLK(t)
	tmp := t.TempDir()
	for i := 0; i < 5; i++ {
		if err := os.WriteFile(filepath.Join(tmp, fmt.Sprintf("audio%02d.mp3", i)), []byte("fake"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	cmd := exec.Command("go", "run", "../../cmd/lk", "audio", "--sort", "shuffle", "--null", tmp)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("command failed: %v\nstderr: %s", err, stderr.String())
	}
	// Count null-terminated entries.
	nullCount := strings.Count(stdout.String(), "\x00")
	if nullCount != 5 {
		t.Fatalf("expected 5 null-terminated entries, got %d", nullCount)
	}
}
