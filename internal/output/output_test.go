package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/amurru/lk/internal/domain"
)

func TestWriteJSON(t *testing.T) {
	var buf bytes.Buffer
	entries := []domain.FileEntry{{Path: "a.txt", Name: "a.txt", Kind: "documents", MatchedBy: domain.MatchSourceExtension, ModifiedAt: time.Unix(0, 0)}}
	if err := writeJSON(&buf, entries); err != nil {
		t.Fatal(err)
	}
	var decoded []domain.FileEntry
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &decoded); err != nil {
		t.Fatal(err)
	}
	if len(decoded) != 1 || decoded[0].Kind != "documents" {
		t.Fatalf("unexpected decoded output: %#v", decoded)
	}
}

func TestWriteTableEmptyMessage(t *testing.T) {
	var buf bytes.Buffer
	if err := writeTable(&buf, nil, "documents"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No documents files found.") {
		t.Fatalf("unexpected table output: %q", buf.String())
	}
}
