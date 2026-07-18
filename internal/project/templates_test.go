package project

import (
	"strings"
	"testing"
)

func TestTemplates_SaveListLoadDelete(t *testing.T) {
	dir := t.TempDir()
	s, _ := NewStorage(dir)

	p := DefaultProject()
	p.Content.Title = "Titre Modele"

	id, err := s.SaveTemplate("Ma Charte RSSI 2026 !", &p)
	if err != nil {
		t.Fatalf("SaveTemplate: %v", err)
	}
	if id != "ma-charte-rssi-2026" {
		t.Errorf("slug inattendu: %q", id)
	}

	// Unicité du slug
	id2, _ := s.SaveTemplate("Ma Charte RSSI 2026 !", &p)
	if id2 != "ma-charte-rssi-2026-2" {
		t.Errorf("suffixe d'unicite attendu, obtenu %q", id2)
	}

	list, err := s.ListTemplates()
	if err != nil || len(list) != 2 {
		t.Fatalf("ListTemplates: %v (%d)", err, len(list))
	}
	if list[0].Name != "Ma Charte RSSI 2026 !" {
		t.Errorf("nom: %q", list[0].Name)
	}

	loaded, err := s.LoadTemplate(id)
	if err != nil {
		t.Fatalf("LoadTemplate: %v", err)
	}
	if loaded.Content.Title != "Titre Modele" {
		t.Errorf("contenu perdu: %q", loaded.Content.Title)
	}
	if loaded.Metadata.Name != "Ma Charte RSSI 2026 !" {
		t.Errorf("nom du modele: %q", loaded.Metadata.Name)
	}

	if err := s.DeleteTemplate(id); err != nil {
		t.Fatalf("DeleteTemplate: %v", err)
	}
	if _, err := s.LoadTemplate(id); err == nil {
		t.Error("le modele supprime se charge encore")
	}

	// ID hostile : pas de traversée de chemin
	if err := s.DeleteTemplate("../projects/x"); err == nil ||
		!strings.Contains(err.Error(), "introuvable") {
		t.Errorf("id hostile non neutralise: %v", err)
	}
}

func TestTemplates_EmptyNameRejected(t *testing.T) {
	s, _ := NewStorage(t.TempDir())
	p := DefaultProject()
	if _, err := s.SaveTemplate("   ", &p); err == nil {
		t.Error("nom vide accepte")
	}
}
