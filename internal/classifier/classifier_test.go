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

func TestClassifierMagicELF(t *testing.T) {
	dir := t.TempDir()
	// Minimal ELF header
	data := []byte("\x7fELF\x02\x01\x01\x00" + "\x00\x00\x00\x00\x00\x00\x00\x00")
	path := filepath.Join(dir, "binary")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}

	c := Classifier{Registry: kinds.DefaultRegistry(), UseMagic: true}
	got, err := c.Classify(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Kind != "executables" || got.MatchedBy != domain.MatchSourceMagic {
		t.Fatalf("expected executables via magic, got %#v", got)
	}
}

func TestClassifierMagicFont(t *testing.T) {
	dir := t.TempDir()
	// Minimal TTF header (TrueType sfVersion)
	data := []byte("\x00\x01\x00\x00\x00\x00\x00\x00")
	path := filepath.Join(dir, "fontfile")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}

	c := Classifier{Registry: kinds.DefaultRegistry(), UseMagic: true}
	got, err := c.Classify(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Kind != "fonts" || got.MatchedBy != domain.MatchSourceMagic {
		t.Fatalf("expected fonts via magic, got %#v", got)
	}
}

func TestClassifierMagicISO(t *testing.T) {
	dir := t.TempDir()
	// ISO 9660 signature at offset 32769: "CD001" at byte 1 of sector 16
	data := make([]byte, 32768+6)
	copy(data[32769:], "CD001\x01")
	path := filepath.Join(dir, "diskimage")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}

	c := Classifier{Registry: kinds.DefaultRegistry(), UseMagic: true}
	got, err := c.Classify(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Kind != "diskimages" || got.MatchedBy != domain.MatchSourceMagic {
		t.Fatalf("expected diskimages via magic, got %#v", got)
	}
}

func TestClassifierMagicGLB(t *testing.T) {
	dir := t.TempDir()
	// Minimal glTF binary header
	data := []byte("glTF\x02\x00\x00\x00\x14\x00\x00\x00\x00\x00\x00\x00\x4a\x53\x4f\x4e")
	path := filepath.Join(dir, "model")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}

	c := Classifier{Registry: kinds.DefaultRegistry(), UseMagic: true}
	got, err := c.Classify(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Kind != "models" || got.MatchedBy != domain.MatchSourceMagic {
		t.Fatalf("expected models via magic, got %#v", got)
	}
}

func TestClassifierMagicSQLite(t *testing.T) {
	dir := t.TempDir()
	// SQLite header: "SQLite format 3\x00"
	data := []byte("SQLite format 3\x00")
	data = append(data, make([]byte, 512-len(data))...)
	path := filepath.Join(dir, "data")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}

	c := Classifier{Registry: kinds.DefaultRegistry(), UseMagic: true}
	got, err := c.Classify(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Kind != "databases" || got.MatchedBy != domain.MatchSourceMagic {
		t.Fatalf("expected databases via magic, got %#v", got)
	}
}
