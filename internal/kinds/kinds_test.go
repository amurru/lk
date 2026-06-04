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

func TestNewKindsAliases(t *testing.T) {
	r := DefaultRegistry()

	aliasTests := []struct {
		alias   string
		want    string
	}{
		{"font", "fonts"},
		{"typefaces", "fonts"},
		{"data", "databases"},
		{"db", "databases"},
		{"bin", "executables"},
		{"binary", "executables"},
		{"exe", "executables"},
		{"3d", "models"},
		{"cad", "models"},
		{"disk", "diskimages"},
		{"iso", "diskimages"},
		{"img", "diskimages"},
	}

	for _, tt := range aliasTests {
		got, ok := r.CanonicalName(tt.alias)
		if !ok {
			t.Errorf("alias %q: expected to resolve, but got false", tt.alias)
			continue
		}
		if got != tt.want {
			t.Errorf("alias %q: expected %q, got %q", tt.alias, tt.want, got)
		}
	}
}

func TestNewKindsExtensionLookup(t *testing.T) {
	r := DefaultRegistry()

	extTests := []struct {
		ext  string
		kind string
	}{
		// fonts
		{".ttf", "fonts"},
		{".otf", "fonts"},
		{".woff", "fonts"},
		{".woff2", "fonts"},
		{".eot", "fonts"},
		// databases
		{".db", "databases"},
		{".sqlite", "databases"},
		{".sqlite3", "databases"},
		{".dbf", "databases"},
		{".mdb", "databases"},
		{".accdb", "databases"},
		// executables
		{".exe", "executables"},
		{".com", "executables"},
		{".dll", "executables"},
		{".so", "executables"},
		{".dylib", "executables"},
		{".bat", "executables"},
		{".cmd", "executables"},
		{".ps1", "executables"},
		// models
		{".obj", "models"},
		{".stl", "models"},
		{".glb", "models"},
		{".gltf", "models"},
		{".fbx", "models"},
		{".dae", "models"},
		{".blend", "models"},
		{".ply", "models"},
		{".u3d", "models"},
		// diskimages
		{".iso", "diskimages"},
		{".img", "diskimages"},
		{".dmg", "diskimages"},
		{".vmdk", "diskimages"},
		{".vhd", "diskimages"},
		{".wim", "diskimages"},
	}

	for _, tt := range extTests {
		got, ok := r.KindForExtension(tt.ext)
		if !ok {
			t.Errorf("extension %q: expected to resolve, but got false", tt.ext)
			continue
		}
		if got != tt.kind {
			t.Errorf("extension %q: expected kind %q, got %q", tt.ext, tt.kind, got)
		}
	}
}

func TestNewKindsMIMELookup(t *testing.T) {
	r := DefaultRegistry()

	mimeTests := []struct {
		mime string
		kind string
	}{
		// fonts
		{"font/ttf", "fonts"},
		{"font/otf", "fonts"},
		{"font/woff", "fonts"},
		{"font/woff2", "fonts"},
		{"application/vnd.ms-fontobject", "fonts"},
		// databases
		{"application/vnd.sqlite3", "databases"},
		{"application/x-sqlite3", "databases"},
		{"application/x-msaccess", "databases"},
		// executables
		{"application/x-executable", "executables"},
		{"application/x-dosexec", "executables"},
		{"application/x-sharedlib", "executables"},
		{"application/x-mach-binary", "executables"},
		// models
		{"model/gltf-binary", "models"},
		{"model/gltf+json", "models"},
		{"model/obj", "models"},
		{"model/stl", "models"},
		// diskimages
		{"application/x-iso9660-image", "diskimages"},
		{"application/x-apple-diskimage", "diskimages"},
	}

	for _, tt := range mimeTests {
		got, ok := r.KindForMIME(tt.mime)
		if !ok {
			t.Errorf("mime %q: expected to resolve, but got false", tt.mime)
			continue
		}
		if got != tt.kind {
			t.Errorf("mime %q: expected kind %q, got %q", tt.mime, tt.kind, got)
		}
	}
}

func TestNewKindsByName(t *testing.T) {
	r := DefaultRegistry()

	for _, name := range []string{"fonts", "databases", "executables", "models", "diskimages"} {
		kind, ok := r.KindByName(name)
		if !ok {
			t.Errorf("expected kind %q to exist in registry", name)
			continue
		}
		if len(kind.Extensions) == 0 {
			t.Errorf("kind %q: expected at least one extension", name)
		}
		if len(kind.MIMEs) == 0 {
			t.Errorf("kind %q: expected at least one MIME type", name)
		}
	}
}

func TestSHStaysInCode(t *testing.T) {
	r := DefaultRegistry()
	got, ok := r.KindForExtension(".sh")
	if !ok {
		t.Fatal("expected .sh to resolve")
	}
	if got != "code" {
		t.Fatalf("expected .sh to resolve to code, got %q", got)
	}
}
