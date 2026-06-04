package walker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/amurru/lk/internal/domain"
	"github.com/amurru/lk/internal/kinds"
)

func TestWalkIncludesRecursiveFilesAndSortsByName(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "b.txt"), []byte("hello"), 0o600); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "a.md"), []byte("# note"), 0o600); err != nil {
		t.Fatal(err)
	}

	w := New(kinds.DefaultRegistry(), true)
	entries, errs, err := w.Walk(dir, domain.ScanOptions{Recursive: true, SortBy: domain.SortByName})
	if err != nil {
		t.Fatal(err)
	}
	if len(errs) != 0 {
		t.Fatalf("expected no walk errors, got %d", len(errs))
	}
	if len(entries) != 2 || entries[0].Name != "a.md" || entries[1].Name != "b.txt" {
		t.Fatalf("unexpected entries: %#v", entries)
	}
}
