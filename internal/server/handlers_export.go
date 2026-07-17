package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

	html := generateStandaloneHTML(p)

	filename := fmt.Sprintf("infographie_%s_%s.html",
		req.ProjectID, time.Now().Format("2006-01-02"))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Write([]byte(html))
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

// generateStandaloneHTML creates a self-contained HTML file with embedded data.
func generateStandaloneHTML(p *project.Project) string {
	projectJSON, _ := json.MarshalIndent(p, "    ", "  ")
	css := web.RenderCSS()
	html := web.BuildStaticInfographic(p)

	var sb strings.Builder
	sb.WriteString("<!DOCTYPE html>\n<html lang=\"fr\">\n<head>\n")
	sb.WriteString("  <meta charset=\"UTF-8\">\n")
	sb.WriteString(fmt.Sprintf("  <meta name=\"generator\" content=\"Generateur Infographie v2.0 (Go)\">\n"))
	sb.WriteString(fmt.Sprintf("  <title>%s</title>\n", p.Content.Title))
	sb.WriteString("  <style>\n")
	sb.WriteString(css)
	sb.WriteString("\n  </style>\n</head>\n<body>\n")
	sb.WriteString(html)
	sb.WriteString("\n  <script id=\"project-data\" type=\"application/json\">\n    ")
	sb.WriteString(string(projectJSON))
	sb.WriteString("\n  </script>\n</body>\n</html>")

	return sb.String()
}
