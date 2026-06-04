package main

import (
	"bytes"
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
