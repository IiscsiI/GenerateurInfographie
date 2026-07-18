package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"infographic-generator/internal/browser"
	"infographic-generator/internal/project"
	"infographic-generator/web"
)

// ExportRequest is the JSON body for export endpoints.
type ExportRequest struct {
	ProjectID   string  `json:"projectId"`
	Format      string  `json:"format"`      // A3, A4, A2
	Orientation string  `json:"orientation"` // portrait, landscape
	DPI         float64 `json:"dpi"`
}

func (s *Server) handleExportPNG(w http.ResponseWriter, r *http.Request) {
	if s.browser == nil || s.browser.Path == "" {
		s.jsonError(w, 503, "Aucun navigateur Chromium detecte. Export PNG impossible.")
		return
	}

	var req ExportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.jsonError(w, 400, "JSON invalide")
		return
	}

	if req.DPI == 0 {
		req.DPI = 300
	}
	if req.Format == "" {
		req.Format = "A3"
	}
	if req.Orientation == "" {
		req.Orientation = "portrait"
	}

	renderURL := fmt.Sprintf("%s/render/%s", s.baseURL, req.ProjectID)

	cfg := browser.ExportConfig{
		URL:         renderURL,
		Format:      req.Format,
		Orientation: req.Orientation,
		DPI:         req.DPI,
		BrowserPath: s.browser.Path,
		ProfileDir:  filepath.Join(s.store.BaseDir, "export-profiles"),
	}

	data, err := browser.ExportPNG(r.Context(), cfg)
	if err != nil {
		s.jsonError(w, 500, fmt.Sprintf("Erreur export PNG: %v", err))
		return
	}

	filename := fmt.Sprintf("infographie_%s_%s_%s_%ddpi_%s.png",
		req.ProjectID, req.Format, req.Orientation, int(req.DPI),
		time.Now().Format("2006-01-02"))

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.Write(data)
}

func (s *Server) handleExportPDF(w http.ResponseWriter, r *http.Request) {
	if s.browser == nil || s.browser.Path == "" {
		s.jsonError(w, 503, "Aucun navigateur Chromium detecte. Export PDF impossible.")
		return
	}

	var req ExportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.jsonError(w, 400, "JSON invalide")
		return
	}

	if req.Format == "" {
		req.Format = "A3"
	}
	if req.Orientation == "" {
		req.Orientation = "portrait"
	}

	renderURL := fmt.Sprintf("%s/render/%s", s.baseURL, req.ProjectID)

	cfg := browser.ExportConfig{
		URL:         renderURL,
		Format:      req.Format,
		Orientation: req.Orientation,
		BrowserPath: s.browser.Path,
		ProfileDir:  filepath.Join(s.store.BaseDir, "export-profiles"),
	}

	data, err := browser.ExportPDF(r.Context(), cfg)
	if err != nil {
		s.jsonError(w, 500, fmt.Sprintf("Erreur export PDF: %v", err))
		return
	}

	filename := fmt.Sprintf("infographie_%s_%s_%s_%s.pdf",
		req.ProjectID, req.Format, req.Orientation,
		time.Now().Format("2006-01-02"))

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.Write(data)
}

func (s *Server) handleExportHTML(w http.ResponseWriter, r *http.Request) {
	var req ExportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.jsonError(w, 400, "JSON invalide")
		return
	}

	p, err := s.store.Load(req.ProjectID)
	if err != nil {
		s.jsonError(w, 404, err.Error())
		return
	}

	// ÉTAPE CLÉ : inliner les assets locaux en base64.
	// Le fichier exporté doit s'afficher sur n'importe quel PC,
	// sans dépendre du serveur local ni de ses fichiers d'assets.
	inlined := project.InlineAssets(p, s.store)

	projectJSON, err := web.MarshalProject(inlined)
	if err != nil {
		s.jsonError(w, 500, "Erreur de serialisation du projet")
		return
	}

	// Le fichier autonome est généré à partir du MÊME template de rendu
	// que la page /render : un seul moteur, zéro divergence preview/export.
	html := web.StandaloneHTML(projectJSON, p.Content.Title, req.Orientation)

	filename := fmt.Sprintf("infographie_%s_%s.html",
		req.ProjectID, time.Now().Format("2006-01-02"))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Write([]byte(html))
}

// handleRenderFrame sert la page de rendu en mode "frame" (PROJECT = null) :
// le rendu est alors piloté par l'éditeur via postMessage. C'est la préviz
// interactive — même moteur que les captures et l'export, zéro divergence.
func (s *Server) handleRenderFrame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(web.RenderTemplate("null")))
}

// handleRenderPage serves a read-only render of the infographic for chromedp capture.
func (s *Server) handleRenderPage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	p, err := s.store.Load(id)
	if err != nil {
		http.Error(w, "Projet introuvable", 404)
		return
	}

	// Content is already sanitized at save/import time.
	// The render page also has its own client-side sanitizer as defense in depth.
	// Re-sanitizing here would double-escape already-safe tags.

	projectJSON, _ := json.Marshal(p)
	renderHTML := web.RenderTemplate(string(projectJSON))

	// Security headers on the render page. The render page runs in a headless
	// browser controlled by chromedp — it's inert relative to users. Primary
	// security comes from sanitization; CSP here is defense in depth only.
	// We deliberately don't set CSP to avoid interfering with chromedp internals.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Write([]byte(renderHTML))
}
