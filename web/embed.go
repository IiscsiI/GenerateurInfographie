package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"html"
	"io/fs"
	"strings"
)

//go:embed editor/*
var editorFiles embed.FS

//go:embed render/*
var renderFiles embed.FS

// EditorFS returns the filesystem for serving the editor frontend.
func EditorFS() fs.FS {
	sub, _ := fs.Sub(editorFiles, "editor")
	return sub
}

// RenderTemplate returns a complete HTML page that renders the infographic
// from injected JSON data. Used by chromedp for export capture AND as the
// basis of the standalone HTML export (single rendering engine).
func RenderTemplate(projectJSON string) string {
	renderHTML, _ := fs.ReadFile(renderFiles, "render/index.html")
	// "</" est neutralisé en "<\/" (échappement JSON/JS valide du solidus) :
	// aucune valeur utilisateur ne peut fermer prématurément le bloc <script>.
	safeJSON := strings.ReplaceAll(projectJSON, "</", "<\\/")
	// Replace the full placeholder including the fallback empty string
	return strings.Replace(string(renderHTML), "\"%%PROJECT_DATA%%\"", safeJSON, 1)
}

// StandaloneHTML builds the self-contained export file from the render
// template. The project JSON must already have its assets inlined (base64).
//
// The exported file:
//   - renders itself with the exact same JS engine as the /render page
//     (no divergence between preview, capture and export),
//   - embeds the project data in <script id="project-data"> for re-import,
//   - opens in the requested default orientation.
func StandaloneHTML(projectJSON, title, orientation string) string {
	page := RenderTemplate(projectJSON)

	// Titre du document (échappé : défense en profondeur)
	page = strings.Replace(page, "<title>Render</title>",
		fmt.Sprintf("<title>%s</title>", html.EscapeString(title)), 1)

	// Orientation par défaut à l'ouverture du fichier (surchargée par ?orientation=)
	if orientation == "landscape" {
		page = strings.Replace(page, "<script>",
			"<script>window.EXPORT_ORIENTATION='landscape';</script>\n    <script>", 1)
	}

	// Bloc de données ré-importable, à la fin du body.
	// Les caractères "</" sont neutralisés pour empêcher toute fermeture
	// prématurée de la balise script par le contenu utilisateur.
	safeJSON := strings.ReplaceAll(projectJSON, "</", "<\\/")
	dataBlock := "\n    <script id=\"project-data\" type=\"application/json\">\n" +
		safeJSON + "\n    </script>\n</body>"
	page = strings.Replace(page, "</body>", dataBlock, 1)

	return page
}

// MarshalProject serialise un projet en JSON indenté pour l'embarquer.
func MarshalProject(v interface{}) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
