package classifier

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/amurru/lk/internal/domain"
	"github.com/amurru/lk/internal/kinds"
)

func TestClassifierUsesExtensionBeforeMagic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "note.PDF")
	if err := os.WriteFile(path, []byte("not a pdf"), 0o600); err != nil {
		t.Fatal(err)
	}

	c := Classifier{Registry: kinds.DefaultRegistry(), UseMagic: true}
	got, err := c.Classify(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Kind != "documents" || got.MatchedBy != domain.MatchSourceExtension {
		t.Fatalf("unexpected classification: %#v", got)
	}
}

func TestClassifierFallsBackToShebangMagic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "script")
	if err := os.WriteFile(path, []byte("#!/bin/sh\necho hi\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	c := Classifier{Registry: kinds.DefaultRegistry(), UseMagic: true}
	got, err := c.Classify(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Kind != "code" || got.MatchedBy != domain.MatchSourceMagic {
		t.Fatalf("unexpected classification: %#v", got)
	}
}

func TestClassifierNoMagicReturnsUnknown(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "script.bin")
	if err := os.WriteFile(path, []byte("#!/bin/sh\necho hi\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	c := Classifier{Registry: kinds.DefaultRegistry(), UseMagic: false}
	got, err := c.Classify(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Kind != "unknown" || got.MatchedBy != domain.MatchSourceUnknown {
		t.Fatalf("unexpected classification: %#v", got)
	}
}
