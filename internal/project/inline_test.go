package project

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// pngHeader : magic bytes d'un PNG minimal valide pour le sniffing MIME.
var pngHeader = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

func newTestStorage(t *testing.T) *Storage {
	t.Helper()
	dir := t.TempDir()
	s, err := NewStorage(dir)
	if err != nil {
		t.Fatalf("NewStorage: %v", err)
	}
	return s
}

func TestInlineAssets_ConvertsLocalAssetToDataURI(t *testing.T) {
	s := newTestStorage(t)

	// Déposer un asset factice
	assetDir := filepath.Join(s.ProjectDir, "monprojet", "assets")
	os.MkdirAll(assetDir, 0755)
	os.WriteFile(filepath.Join(assetDir, "logo.png"), pngHeader, 0644)

	p := DefaultProject()
	p.Elements.Logos = []Logo{
		{ID: 1, File: "/api/projects/monprojet/assets/logo.png", Position: "top-right", Size: 120},
	}

	out := InlineAssets(&p, s)

	if !strings.HasPrefix(out.Elements.Logos[0].File, "data:image/png;base64,") {
		t.Errorf("asset local non inliné en data-URI : %q", out.Elements.Logos[0].File)
	}
	// L'original ne doit jamais être modifié (le projet stocké garde ses URL)
	if !strings.HasPrefix(p.Elements.Logos[0].File, "/api/projects/") {
		t.Errorf("le projet original a été modifié : %q", p.Elements.Logos[0].File)
	}
}

func TestInlineAssets_LeavesExternalAndDataURIsUntouched(t *testing.T) {
	s := newTestStorage(t)
	p := DefaultProject()
	p.Elements.Logos = []Logo{
		{ID: 1, URL: "https://example.org/logo.png", Position: "top-left", Size: 100},
		{ID: 2, File: "data:image/png;base64,AAAA", Position: "bottom-right", Size: 80},
	}

	out := InlineAssets(&p, s)

	if out.Elements.Logos[0].URL != "https://example.org/logo.png" {
		t.Errorf("URL externe modifiée : %q", out.Elements.Logos[0].URL)
	}
	if out.Elements.Logos[1].File != "data:image/png;base64,AAAA" {
		t.Errorf("data-URI existant modifié : %q", out.Elements.Logos[1].File)
	}
}

func TestInlineAssets_MissingAssetKeepsReference(t *testing.T) {
	s := newTestStorage(t)
	p := DefaultProject()
	ref := "/api/projects/inconnu/assets/absent.png"
	p.Elements.Logos = []Logo{{ID: 1, File: ref, Position: "top-right", Size: 120}}

	out := InlineAssets(&p, s)

	if out.Elements.Logos[0].File != ref {
		t.Errorf("référence d'asset manquant altérée : %q", out.Elements.Logos[0].File)
	}
}

func TestParseAssetURL(t *testing.T) {
	cases := []struct {
		in       string
		id, file string
		ok       bool
	}{
		{"/api/projects/p1/assets/a.png", "p1", "a.png", true},
		{"http://127.0.0.1:8080/api/projects/p1/assets/a.png", "p1", "a.png", true},
		{"https://example.org/image.png", "", "", false},
		{"/api/projects/p1/other/a.png", "", "", false},
		{"/api/projects//assets/a.png", "", "", false},
		{"", "", "", false},
	}
	for _, c := range cases {
		id, file, ok := parseAssetURL(c.in)
		if ok != c.ok || id != c.id || file != c.file {
			t.Errorf("parseAssetURL(%q) = (%q,%q,%v), attendu (%q,%q,%v)",
				c.in, id, file, ok, c.id, c.file, c.ok)
		}
	}
}
