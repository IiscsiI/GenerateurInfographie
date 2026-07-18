package server

import (
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"infographic-generator/internal/browser"
	"infographic-generator/internal/project"
	"infographic-generator/web"
)

// Server holds the HTTP server dependencies.
type Server struct {
	store      *project.Storage
	browser    *browser.Browser
	baseURL    string
	OnShutdown func() // Called when /api/shutdown is requested

	// lastBeat : horodatage (UnixNano) du dernier battement de coeur reçu
	// de l'interface. 0 tant qu'aucune fenêtre ne s'est manifestée.
	lastBeat atomic.Int64
}

// LastHeartbeat retourne l'instant du dernier battement (zéro si jamais reçu).
func (s *Server) LastHeartbeat() time.Time {
	n := s.lastBeat.Load()
	if n == 0 {
		return time.Time{}
	}
	return time.Unix(0, n)
}

// New creates a new Server instance.
func New(store *project.Storage, br *browser.Browser, baseURL string) *Server {
	return &Server{
		store:   store,
		browser: br,
		baseURL: baseURL,
	}
}

// Router builds and returns the HTTP handler.
func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("GET /api/projects", s.handleListProjects)
	mux.HandleFunc("POST /api/projects", s.handleCreateProject)
	mux.HandleFunc("GET /api/projects/{id}", s.handleGetProject)
	mux.HandleFunc("PUT /api/projects/{id}", s.handleUpdateProject)
	mux.HandleFunc("DELETE /api/projects/{id}", s.handleDeleteProject)
	mux.HandleFunc("POST /api/projects/{id}/duplicate", s.handleDuplicateProject)

	// Assets
	mux.HandleFunc("POST /api/projects/{id}/assets", s.handleUploadAsset)
	mux.HandleFunc("GET /api/projects/{id}/assets/{filename}", s.handleGetAsset)

	// Export
	mux.HandleFunc("POST /api/export/png", s.handleExportPNG)
	mux.HandleFunc("POST /api/export/pdf", s.handleExportPDF)
	mux.HandleFunc("POST /api/export/html", s.handleExportHTML)

	// Import
	mux.HandleFunc("POST /api/import", s.handleImport)

	// Modeles personnalises
	mux.HandleFunc("GET /api/templates", s.handleListTemplates)
	mux.HandleFunc("POST /api/templates", s.handleSaveTemplate)
	mux.HandleFunc("GET /api/templates/{id}", s.handleGetTemplate)
	mux.HandleFunc("DELETE /api/templates/{id}", s.handleDeleteTemplate)

	// Browser status
	mux.HandleFunc("GET /api/browser/status", s.handleBrowserStatus)

	// Battement de coeur de l'interface : pilote le cycle de vie de
	// l'application (voir main). On ne peut PAS se fier à la fin du
	// processus navigateur : Edge délègue parfois à un frère et notre
	// processus se termine alors que la fenêtre vit toujours.
	mux.HandleFunc("POST /api/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		s.lastBeat.Store(time.Now().UnixNano())
		w.WriteHeader(http.StatusNoContent)
	})

	// Shutdown (called from web UI quit button)
	mux.HandleFunc("POST /api/shutdown", s.handleShutdown)

	// Render page (used by chromedp for export)
	mux.HandleFunc("GET /render/{id}", s.handleRenderPage)

	// Render frame (préviz interactive de l'éditeur : rendu via postMessage)
	mux.HandleFunc("GET /render-frame", s.handleRenderFrame)

	// Favicon : réponse vide propre (évite un 404 dans la console navigateur)
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// Static files (embedded frontend)
	mux.Handle("/", http.FileServer(http.FS(web.EditorFS())))

	return withSecurity(mux)
}

// withSecurity applique les protections du serveur local :
//
//  1. Vérification du header Host : seul 127.0.0.1 / localhost / [::1] est
//     accepté. Bloque les attaques par DNS rebinding, où un site malveillant
//     fait résoudre son propre domaine vers 127.0.0.1 pour contourner la
//     same-origin policy et atteindre l'API.
//
//  2. En-têtes de sécurité systématiques. Noter l'ABSENCE délibérée de CORS :
//     l'éditeur est servi par la même origine que l'API, aucun accès
//     cross-origin n'est légitime. (L'ancien "Access-Control-Allow-Origin: *"
//     permettait à n'importe quel site ouvert dans le navigateur de lire,
//     modifier ou supprimer les projets.)
func withSecurity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		if h, _, err := net.SplitHostPort(host); err == nil {
			host = h
		}
		host = strings.Trim(strings.ToLower(host), "[]")
		if host != "127.0.0.1" && host != "localhost" && host != "::1" {
			http.Error(w, "Forbidden host", http.StatusForbidden)
			return
		}

		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("Referrer-Policy", "no-referrer")

		next.ServeHTTP(w, r)
	})
}
