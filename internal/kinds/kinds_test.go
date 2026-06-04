package kinds

import "testing"

func TestRegistryCanonicalNameAndLookups(t *testing.T) {
	r := DefaultRegistry()

	if got, ok := r.CanonicalName("doc"); !ok || got != "documents" {
		t.Fatalf("expected doc alias to resolve to documents, got %q %v", got, ok)
	}
	if got, ok := r.KindForExtension(".PDF"); !ok || got != "documents" {
		t.Fatalf("expected PDF extension to resolve to documents, got %q %v", got, ok)
	}
	if got, ok := r.KindForMIME("text/plain"); !ok || got != "documents" {
		t.Fatalf("expected text/plain to resolve to documents, got %q %v", got, ok)
	}
}
