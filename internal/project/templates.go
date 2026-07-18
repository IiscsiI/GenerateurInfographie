package project

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ==========================================================================
// MODÈLES PERSONNALISÉS
// --------------------------------------------------------------------------
// Un modèle est un projet complet figé, stocké dans data/templates/<id>.json.
// À la sauvegarde, les logos référencés par fichier sont convertis en data
// URLs (InlineAssets) : le modèle reste valide même si le projet d'origine
// et ses assets sont supprimés plus tard.
// ==========================================================================

// TemplateInfo décrit un modèle pour la liste de la bibliothèque.
type TemplateInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// TemplatesDir retourne (et crée au besoin) le répertoire des modèles.
func (s *Storage) TemplatesDir() string {
	dir := filepath.Join(s.BaseDir, "templates")
	os.MkdirAll(dir, 0755)
	return dir
}

// ListTemplates énumère les modèles enregistrés, triés par nom.
func (s *Storage) ListTemplates() ([]TemplateInfo, error) {
	dir := s.TemplatesDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	out := []TemplateInfo{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		id := strings.TrimSuffix(e.Name(), ".json")
		name := id
		if data, err := os.ReadFile(filepath.Join(dir, e.Name())); err == nil {
			var p Project
			if json.Unmarshal(data, &p) == nil && p.Metadata.Name != "" {
				name = p.Metadata.Name
			}
		}
		out = append(out, TemplateInfo{ID: id, Name: name})
	}
	sort.Slice(out, func(i, j int) bool { return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name) })
	return out, nil
}

// SaveTemplate fige le projet donné comme modèle nommé et retourne son id.
// Les assets fichiers sont inlinés ; le contenu est sanitisé.
func (s *Storage) SaveTemplate(name string, p *Project) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", fmt.Errorf("nom de modele vide")
	}
	if len(name) > 80 {
		name = name[:80]
	}

	// Copie profonde via JSON pour ne pas modifier le projet appelant.
	raw, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	var tpl Project
	if err := json.Unmarshal(raw, &tpl); err != nil {
		return "", err
	}

	SanitizeProject(&tpl)
	inlined := InlineAssets(&tpl, s)
	inlined.Metadata.Name = name
	inlined.Version = "4.1"

	// Identifiant : slug du nom, suffixé si déjà pris.
	base := strings.Trim(slugRe.ReplaceAllString(strings.ToLower(name), "-"), "-")
	if base == "" {
		base = "modele"
	}
	if len(base) > 40 {
		base = base[:40]
	}
	dir := s.TemplatesDir()
	id := base
	for i := 2; ; i++ {
		if _, err := os.Stat(filepath.Join(dir, id+".json")); os.IsNotExist(err) {
			break
		}
		id = fmt.Sprintf("%s-%d", base, i)
	}

	data, err := json.MarshalIndent(inlined, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(dir, id+".json"), data, 0644); err != nil {
		return "", err
	}
	return id, nil
}

// LoadTemplate charge un modèle par id.
func (s *Storage) LoadTemplate(id string) (*Project, error) {
	id = sanitizeID(id)
	data, err := os.ReadFile(filepath.Join(s.TemplatesDir(), id+".json"))
	if err != nil {
		return nil, fmt.Errorf("modele introuvable: %s", id)
	}
	p, err := ImportFromJSON(data)
	if err != nil {
		return nil, err
	}
	SanitizeProject(p)
	return p, nil
}

// DeleteTemplate supprime un modèle par id.
func (s *Storage) DeleteTemplate(id string) error {
	id = sanitizeID(id)
	path := filepath.Join(s.TemplatesDir(), id+".json")
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("modele introuvable: %s", id)
	}
	return os.Remove(path)
}
