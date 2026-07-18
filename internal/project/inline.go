package project

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// InlineAssets remplace, dans une copie du projet, toutes les références
// d'assets locaux (/api/projects/{id}/assets/{fichier}) par des data-URI
// base64. Indispensable pour l'export HTML autonome : le fichier doit
// s'afficher sur n'importe quel PC, sans le serveur local.
//
// Les URL externes (http/https) et les data-URI déjà présents sont
// conservés tels quels. Retourne une copie : le projet stocké sur disque
// n'est jamais modifié (il continue de référencer les assets par URL,
// ce qui permet de les remplacer sans ré-éditer le projet).
func InlineAssets(p *Project, store *Storage) *Project {
	clone := *p // copie superficielle
	logos := make([]Logo, len(p.Elements.Logos))
	copy(logos, p.Elements.Logos)
	clone.Elements.Logos = logos

	for i := range clone.Elements.Logos {
		clone.Elements.Logos[i].File = inlineRef(clone.Elements.Logos[i].File, store)
		// Un champ URL pointant vers un asset local (cas limite) est aussi inliné.
		clone.Elements.Logos[i].URL = inlineRef(clone.Elements.Logos[i].URL, store)
	}
	return &clone
}

// inlineRef convertit une référence d'asset local en data-URI.
// Toute autre valeur (vide, data:, http(s), chemin inconnu) est retournée telle quelle.
func inlineRef(ref string, store *Storage) string {
	if ref == "" || strings.HasPrefix(ref, "data:") {
		return ref
	}

	projectID, filename, ok := parseAssetURL(ref)
	if !ok {
		return ref
	}

	path := store.GetAssetPath(projectID, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		// Asset introuvable : on laisse la référence (image cassée visible,
		// plutôt qu'une disparition silencieuse du logo).
		return ref
	}

	mime := detectMime(filename, data)
	return fmt.Sprintf("data:%s;base64,%s", mime, base64.StdEncoding.EncodeToString(data))
}

// parseAssetURL extrait (projectID, filename) d'une URL de la forme
// /api/projects/{id}/assets/{fichier}. Accepte aussi la forme absolue
// http://127.0.0.1:PORT/api/projects/... produite par d'anciens clients.
func parseAssetURL(ref string) (string, string, bool) {
	// Réduire une URL absolue locale à son chemin
	if strings.HasPrefix(ref, "http://") || strings.HasPrefix(ref, "https://") {
		idx := strings.Index(ref, "/api/projects/")
		if idx == -1 {
			return "", "", false
		}
		ref = ref[idx:]
	}

	const prefix = "/api/projects/"
	if !strings.HasPrefix(ref, prefix) {
		return "", "", false
	}
	rest := strings.TrimPrefix(ref, prefix)
	parts := strings.Split(rest, "/")
	// Attendu : {id} / assets / {fichier}
	if len(parts) != 3 || parts[1] != "assets" || parts[0] == "" || parts[2] == "" {
		return "", "", false
	}
	return parts[0], parts[2], true
}

// detectMime détermine le type MIME d'un asset : d'abord par extension
// (fiable pour nos formats d'image autorisés), sinon par sniffing du contenu.
func detectMime(filename string, data []byte) string {
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	}
	return http.DetectContentType(data)
}
