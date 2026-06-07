package walker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/amurru/lk/internal/domain"
	"github.com/amurru/lk/internal/kinds"
)

func TestSortEntriesShuffle(t *testing.T) {
	entries := []domain.FileEntry{
		{Name: "a.txt", Path: "/a.txt"},
		{Name: "b.txt", Path: "/b.txt"},
		{Name: "c.txt", Path: "/c.txt"},
		{Name: "d.txt", Path: "/d.txt"},
		{Name: "e.txt", Path: "/e.txt"},
	}

	// Run shuffle multiple times and collect the resulting orders.
	orders := make(map[string]bool)
	for i := 0; i < 10; i++ {
		shuffled := make([]domain.FileEntry, len(entries))
		copy(shuffled, entries)
		sortEntries(shuffled, domain.SortByShuffle)
		key := ""
		for _, e := range shuffled {
			key += e.Name
		}
		orders[key] = true

		// Verify it's a valid permutation — same names, just reordered.
		names := make(map[string]int)
		for _, e := range shuffled {
			names[e.Name]++
		}
		for _, orig := range entries {
			if names[orig.Name] != 1 {
				t.Fatalf("shuffle produced invalid permutation: duplicate or missing %q", orig.Name)
			}
		}
	}

	// With 5 entries (120 permutations) and 10 shuffles, we expect at
	// least 2 distinct orderings with overwhelming probability.
	if len(orders) < 2 {
		t.Fatalf("expected at least 2 distinct orderings over 10 shuffles, got %d", len(orders))
	}
}

func TestSortEntriesShufflePreservesEntries(t *testing.T) {
	entries := []domain.FileEntry{
		{Name: "a.txt", Path: "/dir/a.txt", SizeBytes: 100},
		{Name: "b.txt", Path: "/dir/b.txt", SizeBytes: 200},
	}
	sortEntries(entries, domain.SortByShuffle)

	// After shuffle, both entries must still be present with all fields intact.
	names := map[string]domain.FileEntry{}
	for _, e := range entries {
		names[e.Name] = e
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 entries after shuffle, got %d", len(names))
	}
	if names["a.txt"].SizeBytes != 100 || names["b.txt"].SizeBytes != 200 {
		t.Fatal("shuffle corrupted entry fields")
	}
}

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
