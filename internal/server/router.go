package server

import (
	"net/http"

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

	// Browser status
	mux.HandleFunc("GET /api/browser/status", s.handleBrowserStatus)

	// Shutdown (called from web UI quit button)
	mux.HandleFunc("POST /api/shutdown", s.handleShutdown)

	// Render page (used by chromedp for export)
	mux.HandleFunc("GET /render/{id}", s.handleRenderPage)

	// Static files (embedded frontend)
	mux.Handle("/", http.FileServer(http.FS(web.EditorFS())))

	return withCORS(mux)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}
		next.ServeHTTP(w, r)
	})
}
